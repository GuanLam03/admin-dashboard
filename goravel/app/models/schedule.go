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

