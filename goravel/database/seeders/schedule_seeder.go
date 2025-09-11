package seeders

import (
	"fmt"
	"math/rand"
	"time"

	"github.com/goravel/framework/facades"
)

type ScheduleSeeder struct{}

func (s *ScheduleSeeder) Signature() string {
	return "ScheduleSeeder"
}

func (s *ScheduleSeeder) Run() error {
	// recurrences := []string{"daily", "weekly", "monthly"}
	statuses := []string{"active", "inactive"}
	titles := []string{
		"Dev Meeting", "Manager Sync", "Database Backup", "Server Restart",
		"Code Review", "Product Planning", "QA Testing", "Security Patch",
	}

	var rows []map[string]any
	for i := 1; i <= 20; i++ {
		start := time.Now().Add(time.Duration(rand.Intn(240)) * time.Hour)
		end := start.Add(time.Hour)

		rows = append(rows, map[string]any{
			"title":           fmt.Sprintf("%s #%d", titles[rand.Intn(len(titles))], i),
			"recurrence":      nil,
			"start_at":        start,
			"end_at":          end,
			"status":          statuses[rand.Intn(len(statuses))],
			"google_event_id": nil,
			"created_at":      time.Now(),
			"updated_at":      time.Now(),
		})
	}

	_, err := facades.DB().Table("schedules").Insert(rows)
	return err
}
