package crons

import (
	"context"
	"log"
	"time"

	"github.com/robfig/cron/v3"
	"github.com/yourusername/cron-observer/backend/internal/repositories"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

// TaskFailureStatsCron calculates and stores task failure stats every 6 hours
type TaskFailureStatsCron struct {
	repo repositories.Repository
	cron *cron.Cron
}

// NewTaskFailureStatsCron creates a new TaskFailureStatsCron
func NewTaskFailureStatsCron(repo repositories.Repository) *TaskFailureStatsCron {
	c := cron.New(cron.WithSeconds())
	return &TaskFailureStatsCron{
		repo: repo,
		cron: c,
	}
}

// Start starts the cron and schedules the job
func (c *TaskFailureStatsCron) Start(ctx context.Context) {
	// Schedule job to run every 6 hours: "0 0 0,6,12,18 * * *" (at 00:00, 06:00, 12:00, 18:00)
	_, err := c.cron.AddFunc("0 0 0,6,12,18 * * *", func() {
		log.Println("[TaskFailureStatsCron] Starting scheduled calculation...")
		c.calculateAllStats(context.Background())
	})
	if err != nil {
		log.Printf("[TaskFailureStatsCron] Failed to schedule cron job: %v", err)
		return
	}

	// Run immediately on startup
	go func() {
		log.Println("[TaskFailureStatsCron] Running initial calculation...")
		c.calculateAllStats(context.Background())
	}()

	// Start the cron engine
	c.cron.Start()
	log.Println("[TaskFailureStatsCron] Started (runs every 6 hours at 00:00, 06:00, 12:00, 18:00)")

	// Wait for context cancellation
	<-ctx.Done()
	log.Println("[TaskFailureStatsCron] Context cancelled, stopping...")
	c.cron.Stop()
	log.Println("[TaskFailureStatsCron] Stopped")
}

// calculateAllStats calculates stats for all projects for today and yesterday
func (c *TaskFailureStatsCron) calculateAllStats(ctx context.Context) {
	// Get all projects
	projects, err := c.repo.GetAllProjects(ctx)
	if err != nil {
		log.Printf("[TaskFailureStatsCron] Failed to get projects: %v", err)
		return
	}

	// Calculate stats for today and yesterday
	today := time.Now().UTC()
	yesterday := today.AddDate(0, 0, -1)
	dates := []time.Time{today, yesterday}

	for _, project := range projects {
		for _, date := range dates {
			dateStr := date.Format("2006-01-02")
			if err := c.calculateStatsForProjectAndDate(ctx, project.ID, dateStr); err != nil {
				log.Printf("[TaskFailureStatsCron] Failed to calculate stats for project %s on date %s: %v", project.ID.Hex(), dateStr, err)
				// Continue with other projects/dates
			}
		}
	}

	log.Println("[TaskFailureStatsCron] Completed scheduled calculation")
}

// calculateStatsForProjectAndDate calculates and stores stats for a specific project and date
func (c *TaskFailureStatsCron) calculateStatsForProjectAndDate(ctx context.Context, projectID primitive.ObjectID, date string) error {
	// Calculate stats
	stats, err := c.repo.CalculateTaskFailureStats(ctx, projectID, date)
	if err != nil {
		return err
	}

	// Store stats
	if err := c.repo.StoreTaskFailureStats(ctx, stats); err != nil {
		return err
	}

	log.Printf("[TaskFailureStatsCron] Calculated and stored stats for project %s on date %s (total failures: %d)", projectID.Hex(), date, stats.Total)
	return nil
}
