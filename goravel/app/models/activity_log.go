package models

import (
	"time"
	"github.com/goravel/framework/database/orm"
	"gorm.io/datatypes"
)

type ActivityLog struct {
	orm.Model
	CauserId    uint   `json:"causer_id"`
	CauserType  string `json:"causer_type"`
	Properties  string `json:"properties"`
	Url 		string 	`json:url`
	Route		string 	`json:route`
	Input       string `json:"input"`
	LogName     string `json:"log_name"`
	Description string `json:"description"`
	RequestMeta datatypes.JSON `json:"request_meta"`     

	StartAt     *time.Time `json:"start_at" gorm:"type:timestamp(3)"`
	EndAt       *time.Time `json:"end_at" gorm:"type:timestamp(3)"`
}


// Ip              *string `json:"ip"`
// 	Country         *string `json:"country"`
// 	Region          *string `json:"region"`
// 	City            *string `json:"city"`
// 	UserAgent       *string `json:"user_agent"`
// 	Referrer        *string `json:"referrer"`
	
// 	DeviceType      *string `json:"device_type"`      // e.g., "mobile", "desktop", "tablet"
// 	DeviceName      *string `json:"device_name"`      // e.g., "iPhone", "Windows PC"
// 	OsName          *string `json:"os_name"`          // e.g., "iOS", "Android", "Windows"
// 	OsVersion       *string `json:"os_version"`       // e.g., "14.4", "10"
// 	BrowserName     *string `json:"browser_name"`     // e.g., "Chrome", "Safari"
// 	BrowserVersion  *string `json:"browser_version"`  // e.g., "117.0.0.0"
