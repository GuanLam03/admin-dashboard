package models

import "github.com/goravel/framework/database/orm"

type AdsCampaignPostbackLog struct {
	orm.Model

	AdsEventLogId         *uint   `json:"ads_event_log_id"`          
	AdsCampaignPostbackId uint   `json:"ads_campaign_postback_id"`  
	Url                   string `json:"url"`                       
	RequestMethod         string `json:"request_method"`            
	RequestBody           *string `json:"request_body"`             
	ResponseStatus        *int    `json:"response_status"`           
	ResponseBody          *string `json:"response_body"`             
	ErrorMessage          *string `json:"error_message"`  
	Status 				  string  `json:"status"`         
}

var AdsPostbackLogRules = map[string]string{
	"ads_campaign_postback_id": "required|numeric",
	"url":                      "required|url",
	"request_method":           "required|string|in:GET,POST,PUT,DELETE",
	"request_body":             "string",
	"response_status":          "numeric",
	"response_body":            "string",
	"error_message":            "string",
	"status": 					"string",
}

var AdsCampaignPostbackLogStatusMap = map[string]string{
    "pending":   "pending",    
    "successful": "successful",
    "failed":    "failed",    
}
