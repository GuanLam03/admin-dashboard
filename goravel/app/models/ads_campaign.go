package models

import "github.com/goravel/framework/database/orm"

type AdsCampaign struct {
	orm.Model
	Name string `json:"name"`
	TargetUrl string `json:"target_url"`
	Code string `json:"code"`
	TrackingLink *string `json:"tracking_link"`
	PostbackLink *string `json:"postback_link"`


}

var AdsCampaignRules = map[string]string{
    "name":          "required|string",
    "targetUrl":     "required|string",
    "code":          "string",
    
}