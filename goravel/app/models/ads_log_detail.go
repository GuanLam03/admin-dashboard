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

	ClickedUrl   string `json:"clicked_url"`
	AdsCampaignId   uint `json:"ads_campaign_id"`

}





var AdsLogDetailErrorMessage = map[string]string{
	"not_found":         "Ads log detail not found.",
	"create_failed":     "Failed to create the ads log detail.",
	"validation_failed": "Invalid input. Please check the fields and try again.",
	"invalid_request":   "Invalid request body. Please check your JSON format.",
	"internal_error":    "Something went wrong. Please try again later.",
}
