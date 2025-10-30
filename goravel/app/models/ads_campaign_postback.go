package models

import "github.com/goravel/framework/database/orm"

type AdsCampaignPostback struct {
	orm.Model
	AdsCampaignId uint       `json:"ads_campaign_id"` 
	EventName string       `json:"event_name"`  // Enum 
	PostbackUrl string `json:"postback_url"`


}

var AdsCampaignPostbackRules = map[string]string{
	"adsCampaignId": "required|numeric",
	"eventName":	"required|string",
	"postbackUrl":  "required|string",
    
}

