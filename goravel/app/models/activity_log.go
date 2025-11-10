package models

import "github.com/goravel/framework/database/orm"

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
}

