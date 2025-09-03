package models

import "github.com/goravel/framework/database/orm"

type UserRole struct {
	orm.Model
	UserID uint
	RoleID uint
}
