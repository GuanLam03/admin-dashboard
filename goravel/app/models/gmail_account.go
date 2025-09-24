package models

import (
	"time"

	"github.com/goravel/framework/database/orm"
)

type GmailAccount struct {
	orm.Model

	Email        string     `gorm:"type:varchar(255);unique" json:"email"`
	AccessToken  string     `gorm:"type:text" json:"access_token"`
	RefreshToken string     `gorm:"type:text" json:"refresh_token,omitempty"`
	Expiry       *time.Time `json:"expiry,omitempty"`
	
}

// ✅ Validation rules
var GmailAccountRules = map[string]string{
	"email":         "required|email",
	"access_token":  "required|string",
	"refresh_token": "string",
	"expiry":        "date",

}


