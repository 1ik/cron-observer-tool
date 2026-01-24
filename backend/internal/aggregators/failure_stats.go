package aggregators

import (
	"context"
	"log"
	"time"

	"github.com/yourusername/cron-observer/backend/internal/events"
	"github.com/yourusername/cron-observer/backend/internal/repositories"
)

type FailureStatsAggregator struct {
	repo     repositories.Repository
	eventBus *events.EventBus
}

func NewFailureStatsAggregator(repo repositories.Repository, eventBus *events.EventBus) *FailureStatsAggregator {
	return &FailureStatsAggregator{
		repo:     repo,
		eventBus: eventBus,
	}
}

func (a *FailureStatsAggregator) Start(ctx context.Context) {
	executionFailedCh := a.eventBus.Subscribe(events.ExecutionFailed)

	go func() {
		for {
			select {
			case <-ctx.Done():
				log.Println("FailureStatsAggregator context cancelled, stopping")
				return
			case event, ok := <-executionFailedCh:
				if !ok {
					log.Println("ExecutionFailed channel closed")
					return
				}
				a.handleExecutionFailed(event)
			}
		}
	}()
}

func (a *FailureStatsAggregator) handleExecutionFailed(event events.Event) {
	payload, ok := event.Payload.(events.ExecutionFailedPayload)
	if !ok {
		log.Printf("Invalid payload for ExecutionFailed event")
		return
	}

	// Extract date from execution (use ended_at if available, else started_at)
	var date time.Time
	if payload.Execution.EndedAt != nil {
		date = *payload.Execution.EndedAt
	} else {
		date = payload.Execution.StartedAt
	}

	// Format as YYYY-MM-DD (in UTC)
	dateStr := date.UTC().Format("2006-01-02")

	// Increment stat for the project and date
	ctx := context.Background()
	if err := a.repo.IncrementFailureStat(ctx, payload.Task.ProjectID, dateStr); err != nil {
		log.Printf("Failed to increment failure stat: %v", err)
	}
}

