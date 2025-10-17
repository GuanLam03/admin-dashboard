package models

import (
	"github.com/goravel/framework/database/orm"
)

type AdsLog struct {
	orm.Model

	AdsCampaignId     uint `json:"ads_campaign_id"` 
	AdsLogDetailId      uint   `json:"ads_log_detail_id"`
	ClickedUrl   string `json:"clicked_url"`
}
