package alert

import (
	"context"
	"fmt"
	"log"
	"time"

	"github.com/yourusername/cron-observer/backend/internal/events"
	"github.com/yourusername/cron-observer/backend/internal/gmail"
	"github.com/yourusername/cron-observer/backend/internal/models"
	"github.com/yourusername/cron-observer/backend/internal/repositories"
)

// Service handles alert notifications for execution failures
type Service struct {
	repo        repositories.Repository
	eventBus    *events.EventBus
	gmailSender gmail.Sender
}

// NewService creates a new alert service
func NewService(repo repositories.Repository, eventBus *events.EventBus, gmailSender gmail.Sender) *Service {
	return &Service{
		repo:        repo,
		eventBus:    eventBus,
		gmailSender: gmailSender,
	}
}

// Start starts the alert service and begins listening for execution failed events
func (s *Service) Start(ctx context.Context) {
	executionFailedCh := s.eventBus.Subscribe(events.ExecutionFailed)

	go func() {
		for {
			select {
			case <-ctx.Done():
				log.Println("[AlertService] Context cancelled, stopping")
				return
			case event, ok := <-executionFailedCh:
				if !ok {
					log.Println("[AlertService] ExecutionFailed channel closed")
					return
				}
				s.handleExecutionFailed(event)
			}
		}
	}()

	log.Println("[AlertService] Started and listening for execution failed events")
}

// handleExecutionFailed processes an execution failed event and sends alerts
func (s *Service) handleExecutionFailed(event events.Event) {
	payload, ok := event.Payload.(events.ExecutionFailedPayload)
	if !ok {
		log.Printf("[AlertService] Invalid payload for ExecutionFailed event")
		return
	}

	// Get project from task's ProjectID
	ctx := context.Background()
	project, err := s.repo.GetProjectByID(ctx, payload.Task.ProjectID)
	if err != nil {
		log.Printf("[AlertService] Failed to get project %s: %v", payload.Task.ProjectID.Hex(), err)
		return
	}

	// Check if Gmail sender is available
	if s.gmailSender == nil {
		log.Printf("[AlertService] Gmail sender not configured, skipping alert for task %s", payload.Task.UUID)
		return
	}

	// Collect email addresses from project_users
	var recipients []string
	for _, projectUser := range project.ProjectUsers {
		if projectUser.Email != "" {
			recipients = append(recipients, projectUser.Email)
		}
	}

	// If no project users, skip sending alert
	if len(recipients) == 0 {
		log.Printf("[AlertService] No project users found for project %s, skipping alert", project.Name)
		return
	}

	// Format execution time
	executionTime := payload.Execution.StartedAt.Format(time.RFC3339)
	if payload.Execution.EndedAt != nil {
		executionTime = payload.Execution.EndedAt.Format(time.RFC3339)
	}

	// Build email subject and body
	subject := fmt.Sprintf("Task Execution Failed: %s", payload.Task.Name)
	body := s.buildEmailBody(payload, project, executionTime)

	// Send email to all project users
	msg := gmail.EmailMessage{
		To:      recipients,
		Subject: subject,
		Body:    body,
	}

	if err := s.gmailSender.Send(msg); err != nil {
		log.Printf("[AlertService] Failed to send alert email for task %s: %v", payload.Task.UUID, err)
		return
	}

	log.Printf("[AlertService] Successfully sent alert email to %d recipients for failed task %s", len(recipients), payload.Task.UUID)
}

// buildEmailBody creates the HTML email body for the alert
func (s *Service) buildEmailBody(payload events.ExecutionFailedPayload, project *models.Project, executionTime string) string {
	errorMsg := "No error message available"
	if payload.Execution.Error != "" {
		errorMsg = payload.Execution.Error
	}

	html := fmt.Sprintf(`
<!DOCTYPE html>
<html>
<head>
	<style>
		body { font-family: Arial, sans-serif; line-height: 1.6; color: #333; }
		.container { max-width: 600px; margin: 0 auto; padding: 20px; }
		.header { background-color: #dc3545; color: white; padding: 20px; border-radius: 5px 5px 0 0; }
		.content { background-color: #f8f9fa; padding: 20px; border: 1px solid #dee2e6; border-top: none; }
		.detail-row { margin: 10px 0; }
		.label { font-weight: bold; color: #495057; }
		.value { color: #212529; }
		.error-box { background-color: #fff3cd; border: 1px solid #ffc107; border-radius: 4px; padding: 15px; margin: 15px 0; }
		.footer { margin-top: 20px; padding-top: 20px; border-top: 1px solid #dee2e6; font-size: 12px; color: #6c757d; }
	</style>
</head>
<body>
	<div class="container">
		<div class="header">
			<h2 style="margin: 0;">⚠️ Task Execution Failed</h2>
		</div>
		<div class="content">
			<div class="detail-row">
				<span class="label">Project:</span>
				<span class="value">%s</span>
			</div>
			<div class="detail-row">
				<span class="label">Task Name:</span>
				<span class="value">%s</span>
			</div>
			<div class="detail-row">
				<span class="label">Task UUID:</span>
				<span class="value">%s</span>
			</div>
			<div class="detail-row">
				<span class="label">Execution UUID:</span>
				<span class="value">%s</span>
			</div>
			<div class="detail-row">
				<span class="label">Execution Time:</span>
				<span class="value">%s</span>
			</div>
			<div class="error-box">
				<strong>Error Message:</strong><br>
				%s
			</div>
		</div>
		<div class="footer">
			<p>This is an automated alert from Cron Observer. Please check the task execution logs for more details.</p>
		</div>
	</div>
</body>
</html>
`,
		project.Name,
		payload.Task.Name,
		payload.Task.UUID,
		payload.Execution.UUID,
		executionTime,
		errorMsg,
	)

	return html
}
