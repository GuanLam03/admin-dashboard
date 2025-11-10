package models

import (
	"time"

	"github.com/goravel/framework/database/orm"
)

type GmailAccount struct {
	orm.Model

	Email        string     `gorm:"type:varchar(255);unique" json:"email"`
	Team        *string    `json:"team,omitempty"`
	AccessToken  string     `gorm:"type:text" json:"access_token"`
	RefreshToken string     `gorm:"type:text" json:"refresh_token,omitempty"`
	Expiry       *time.Time `json:"expiry,omitempty"`
	
}

// Validation rules
var GmailAccountRules = map[string]string{
	"email":         "required|email",
	"access_token":  "required|string",
	"refresh_token": "string",
	"expiry":        "date",

}


var GmailAccountTeams = map[string]string{
	"technical":"Technical Support",
	"support":"Customer Support",
	"info":"Info",
}

var GmailAccountErrorMessage = map[string]string{
	// General Errors
	"internal_error":     "Something went wrong on our end. Please try again later.",
	"invalid_request":    "Your request could not be processed. Please check your input and try again.",
	"validation_failed":  "Some fields are invalid. Please review your input and try again.",

	// Authentication & Authorization
	"auth_failed":        "Failed to connect with Google. Please reauthorize your Gmail account.",
	"token_expired":      "Your Gmail session has expired. Please reconnect your account.",
	"permission_denied":  "You don’t have permission to access this Gmail account.",

	// Gmail Account Management
	"not_found":          "The Gmail account you’re trying to access was not found.",
	"already_linked":     "This Gmail account is already linked to another department.",
	"create_failed":      "Unable to link your Gmail account. Please try again.",
	"update_failed":      "Failed to update your Gmail account. Please try again later.",
	"delete_failed":      "Failed to remove this Gmail account. Please try again.",
	"fetch_failed":       "Unable to fetch Gmail account information. Please refresh and try again.",

	// Gmail Data Operations
	"list_failed":        "Unable to load your Gmail messages right now. Please try again later.",
	"read_failed":        "Unable to open this email thread. Please try again later.",
	"thread_not_found":   "The selected email conversation could not be found.",
	"star_toggle_failed": "We couldn’t update the star status for this conversation. Please retry.",
	"save_draft_failed":  "We couldn’t save your email draft. Please try again.",
	"send_failed":        "Failed to send. Please try again later.",
	"forward_failed":     "Failed to forward message. Please try again later.",

	
}
