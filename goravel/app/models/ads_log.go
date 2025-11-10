package models

import (
	"github.com/goravel/framework/database/orm"
)

type AdsLog struct {
	orm.Model

	AdsCampaignId     uint `json:"ads_campaign_id"` 
	AdsLogDetailId      uint   `json:"ads_log_detail_id"`
	
}


var AdsLogErrorMessage = map[string]string{
	"not_found":         "Ads log not found.",
	"create_failed":     "Failed to create ads log.",
	"validation_failed": "Invalid input. Please check the fields and try again",
	"invalid_request":   "Invalid request body. Please check your JSON format.",
	"internal_error":    "Something went wrong. Please try again later.",
}