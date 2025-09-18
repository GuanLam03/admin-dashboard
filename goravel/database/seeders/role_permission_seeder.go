package seeders

import (
	"fmt"
	"strconv"

	"goravel/app/models"
	per "goravel/app/permissions"
	"github.com/casbin/casbin/v2"
	"github.com/goravel/framework/facades"
)

type RolePermissionSeeder struct{}

func (s *RolePermissionSeeder) Signature() string {
	return "RolePermissionSeeder"
}

func (s *RolePermissionSeeder) Run() error {
	// Define the role-permission mapping
	rolePermissionMap := map[string][]string{
		"super admin": {"roles.index", "roles.create", "roles.edit"},
		"reader":      {"roles.index"},
		"developer":   {"roles.index", "roles.create"},
	}

	// Get Casbin enforcer
	enforcerAny, err := facades.App().Make("casbin")
	if err != nil {
		return fmt.Errorf("failed to get Casbin enforcer: %v", err)
	}
	enforcer := enforcerAny.(*casbin.Enforcer)

	for roleName, permissions := range rolePermissionMap {
		// Find the role by name
		role := models.Role{}
		if err := facades.Orm().Query().Where("name", roleName).First(&role); err != nil {
			return fmt.Errorf("role '%s' not found: %v", roleName, err)
		}

		roleIDStr := strconv.Itoa(int(role.ID))

		// Assign each permission to the role
		for _, permKey := range permissions {
			object, action := per.PermissionKeyToObjectAction(permKey)
			if object != "" && action != "" {
				
				_, err := enforcer.AddPolicy(roleIDStr, object, action)
				if err != nil {
					return fmt.Errorf("failed to add policy for role %s (%d): %v", roleName, role.ID, err)
				}
			}
		}
	}

	return nil
}
