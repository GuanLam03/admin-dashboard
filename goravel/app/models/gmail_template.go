package models

import (
	"github.com/goravel/framework/database/orm"
)


type GmailTemplate struct {
	orm.Model
	Team      string    `json:"team"`
	Name      string    `json:"name"`
	Content   string    `json:"content"`
	
}


var GmailTemplateErrorMessage = map[string]string{
	"not_found":         "Gmail template not found.",
	"create_failed":     "Failed to create the gmail template.",
	"validation_failed": "Invalid input. Please check the fields and try again.",
	"invalid_request":   "Invalid request body. Please check your JSON format.",
	"internal_error":    "Something went wrong. Please try again later.",
}