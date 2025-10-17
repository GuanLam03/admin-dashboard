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
