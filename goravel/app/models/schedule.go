package models

import (
	"time"

	"github.com/goravel/framework/database/orm"
)

type Schedule struct {
	orm.Model

	Title         string     `gorm:"type:varchar(50)" json:"title"`
	Recurrence    *string    `gorm:"type:varchar(50)" json:"recurrence"`  // nullable
	StartAt       time.Time  `json:"start_at"`
	EndAt         time.Time  `json:"end_at"`
	Status        string     `gorm:"type:varchar(50)" json:"status"` 
	GoogleEventID *string    `gorm:"type:varchar(255)" json:"google_event_id"` // nullable


}

var ScheduleRules = map[string]string{
    "title":          "required|string|max_len:50", 
    "recurrence":     "string", 
    "start_at":       "required|date",  
    "end_at":         "required|date", 
    "status":         "required|string", 
    "google_event_id": "string", 
}

var ScheduleStatusMap = map[string]string{
    "active":   "active",
    "inactive": "inactive",
    "removed":  "removed",
}

var ScheduleErrorMessage = map[string]string{
	// General errors
	"internal_error":    "Something went wrong. Please try again later.",
	"validation_failed": "Invalid input. Please check the fields and try again.",
	"invalid_request":   "Invalid request body. Please check your JSON format.",
	"not_found":         "Schedule not found.",

	// CRUD errors
	"create_failed":     "Failed to create schedule. Please try again.",
	"update_failed":     "Failed to update schedule. Please try again.",
	"delete_failed":     "Failed to delete schedule. Please try again.",

	// Date/time errors
	"invalid_start_at":  "Invalid start_at format. Use YYYY-MM-DDTHH:MM",
	"invalid_end_at":    "Invalid end_at format. Use YYYY-MM-DDTHH:MM",

	// Google Calendar errors
	"google_insert_failed": "Failed to create Google Calendar event.",
	"google_update_failed": "Failed to update Google Calendar event.",
	"google_delete_failed": "Failed to delete Google Calendar event.",

	

}