package controllers

import (
	"github.com/goravel/framework/contracts/http"
	"goravel/app/permissions"
)

type PermissionController struct{}

// NewPermissionController is the constructor for PermissionController
func NewPermissionController() *PermissionController {
	return &PermissionController{}
}

// Index returns all permissions (for frontend display)
type Permission struct {
    Key   string `json:"key"`
    Label string `json:"label"`
}

func (p *PermissionController) Index(ctx http.Context) http.Response {
    perms := []Permission{}
    for k, v := range permissions.Permissions {
        perms = append(perms, Permission{Key: k, Label: v.Label})
    }

    return ctx.Response().Json(200, map[string]interface{}{
        "permissions": perms,
    })
}

