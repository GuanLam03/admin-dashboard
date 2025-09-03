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

