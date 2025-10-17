package models

import (
	"github.com/goravel/framework/database/orm"
)

type AdsLogDetail struct {
	orm.Model

	Ip              *string `json:"ip"`
	Country         *string `json:"country"`
	Region          *string `json:"region"`
	City            *string `json:"city"`
	UserAgent       *string `json:"user_agent"`
	Referrer        *string `json:"referrer"`
	
	DeviceType      *string `json:"device_type"`      // e.g., "mobile", "desktop", "tablet"
	DeviceName      *string `json:"device_name"`      // e.g., "iPhone", "Windows PC"
	OsName          *string `json:"os_name"`          // e.g., "iOS", "Android", "Windows"
	OsVersion       *string `json:"os_version"`       // e.g., "14.4", "10"
	BrowserName     *string `json:"browser_name"`     // e.g., "Chrome", "Safari"
	BrowserVersion  *string `json:"browser_version"`  // e.g., "117.0.0.0"
}
