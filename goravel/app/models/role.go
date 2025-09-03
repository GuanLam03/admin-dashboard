package models

import "github.com/goravel/framework/database/orm"

type Role struct {
	orm.Model
	Name string `gorm:"unique;not null"`
}


