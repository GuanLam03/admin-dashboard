package models

import (
	"github.com/goravel/framework/database/orm"
)

type AdsLog struct {
	orm.Model

	AdsCampaignId     uint    `json:"ads_campaign_id"` 
	Ip                *string `json:"ip"`              // Nullable fields use pointers
	Country           *string `json:"country"`
	Region            *string `json:"region"`
	City              *string `json:"city"`
	UserAgent         *string `json:"user_agent"`
	Referrer          *string `json:"referrer"`
	Converted         bool    `json:"converted"`
	ClientProductId   *uint   `json:"client_product_id"`
	Value             *float64 `json:"value"`           
}
