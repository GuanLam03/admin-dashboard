package seeders

import (
	"fmt"
	"math/rand"
	"time"

	"github.com/goravel/framework/facades"
	"goravel/app/models"
)

type ScheduleSeeder struct{}

func (s *ScheduleSeeder) Signature() string {
	return "ScheduleSeeder"
}

func (s *ScheduleSeeder) Run() error {
	statuses := []string{"active", "inactive"}
	titles := []string{
		"Dev Meeting", "Manager Sync", "Database Backup", "Server Restart",
		"Code Review", "Product Planning", "QA Testing", "Security Patch",
	}

	for i := 1; i <= 20; i++ {
		start := time.Now().Add(time.Duration(rand.Intn(240)) * time.Hour)
		end := start.Add(time.Hour)

		
		schedule := models.Schedule{
			Title:      fmt.Sprintf("%s #%d", titles[rand.Intn(len(titles))], i),
			Recurrence: nil,
			StartAt:    start,
			EndAt:      end,
			Status:     statuses[rand.Intn(len(statuses))],
			GoogleEventID: nil,
		}

		if err := facades.Orm().Query().Create(&schedule); err != nil {
			return err 
		}
	}

	return nil
}