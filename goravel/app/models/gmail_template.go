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
