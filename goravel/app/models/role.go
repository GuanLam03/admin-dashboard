package models

import "github.com/goravel/framework/database/orm"

type Role struct {
	orm.Model
	Name string `gorm:"unique;not null"`
}

var RoleErrorMessage = map[string]string{
	// General errors
	"internal_error":    "Something went wrong. Please try again later.",
	"invalid_request":   "Invalid request body. Please check your JSON format.",
	"not_found":         "Role not found.",
	"validation_failed": "Missing or invalid input data.",


	// CRUD errors
	"create_failed":     "Failed to create role. Please try again.",
	"update_failed":     "Failed to update role. Please try again.",
	"delete_failed":     "Failed to delete role. Please try again.",

	"assign_failed":     "Failed to assign role. Please try again.",

	// Casbin errors
	"casbin_not_initialized": "Casbin is not initialized.",
	"casbin_cast_failed":     "Failed to cast Casbin enforcer.",
}
