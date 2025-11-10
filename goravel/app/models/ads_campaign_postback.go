package models

import "github.com/goravel/framework/database/orm"

type AdsCampaignPostback struct {
	orm.Model
	AdsCampaignId uint       `json:"ads_campaign_id"` 
	EventName string       `json:"event_name"`  // Enum 
	PostbackUrl string `json:"postback_url"`
	IncludeClickParams bool   `json:"include_click_params"` // default is false

}

var AdsCampaignPostbackRules = map[string]string{
	"adsCampaignId": "numeric",
	"eventName":	"required|string",
	"postbackUrl":  "required|string",
	"includeClickParams": "bool",
    
}

