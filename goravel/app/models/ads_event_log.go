package models



import (
	"github.com/goravel/framework/database/orm"
	"gorm.io/datatypes"
)



type AdsEventLog struct {
	orm.Model

	AdsLogId uint       `json:"ads_log_id"` 
	EventName string       `json:"event_name"`  // Enum 
	Data      datatypes.JSON     `json:"data"`       
}


var AdsEventLogRules = map[string]string{
	"ads_log_id": "required|numeric",
	"event_name": "required|string",
	"data":       "required",
}

var AllowedEventDataFields = []string{
	"click_id",
	"content_id",
	"content_type",
	"value",
	"currency",
	"price",
}


var AdsEventLogErrorMessage = map[string]string{
	"not_found":         "Ads event log not found.",
	"create_failed":     "Failed to create ads event log.",
	"validation_failed": "Invalid input. Please check the fields and try again",
	"invalid_request":   "Invalid request body. Please check your JSON format.",
	"internal_error":    "Something went wrong. Please try again later.",
}