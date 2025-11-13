package models

import "github.com/goravel/framework/database/orm"

type GoogleDocument struct {
    orm.Model

    Name   string `gorm:"type:varchar(255)" json:"name"`
    OriginalLink   string `gorm:"type:text" json:"original_link"`
    Link   string `gorm:"type:text" json:"link",omitempty`
    Status string `gorm:"type:varchar(50)" json:"status"`
}

// Validation rules
var GoogleDocumentRules = map[string]string{
    "name":          "required|string",
    "original_link": "required|string",
    "link":          "string",
    "status":        "required|string",
}


var GoogleDocumentStatusMap = map[string]string{
    "active":   "active",
    "inactive": "inactive",
    "removed":  "removed",
}

var GoogleDocumentErrorMessage = map[string]string{
	// General errors
	"internal_error":    "Something went wrong. Please try again later.",
	"validation_failed": "Invalid input. Please check the fields and try again.",
	"invalid_request":   "Invalid request body. Please check your JSON format.",
	"not_found":         "Google document not found.",

	// CRUD errors
	"create_failed": "Failed to create Google document. Please try again.",
	"update_failed": "Failed to update Google document. Please try again.",
	"delete_failed": "Failed to remove Google document. Please try again.",

	// Link errors
	"invalid_link_format":  "Google file link format is invalid.",
	"link_not_accessible":  "Google file link is not accessible or not public.",

	// Status errors
	"invalid_status": "Invalid status provided.",
}