package seeders

import (
	"fmt"
	"strconv"
	"goravel/app/models"
	"github.com/goravel/framework/facades"
    "github.com/casbin/casbin/v2"
)

// UserRoleSeeder is a struct for seeding user roles.
type UserRoleSeeder struct{}

// Signature returns the name of the seeder.
func (s *UserRoleSeeder) Signature() string {
	return "UserRoleSeeder"
}

// Run executes the seeder logic.
func (s *UserRoleSeeder) Run() error {
	// Map of users and their corresponding role IDs
	userRoleMap := map[string]string{
		"manager@gmail.com": "super admin",
		"admin@gmail.com":   "reader",  
		"developer@gmail.com": "developer", 
	
	}

	// Get the Casbin enforcer instance
	enforcerAny, err := facades.App().Make("casbin")
	if err != nil {
		return fmt.Errorf("Failed to get Casbin enforcer: %v", err)
	}
	e := enforcerAny.(*casbin.Enforcer)

	// Loop through the map and assign roles to users by ID
	for userEmail, roleName := range userRoleMap {
		// Find the user by email
		user := models.User{}
		if err := facades.Orm().Query().Where("email", userEmail).First(&user); err != nil {
			return fmt.Errorf("User not found: %v", userEmail)
		}

		// Check if the role exists by ID
		role := models.Role{}
		if err := facades.Orm().Query().Where("name", roleName).First(&role); err != nil {
			return fmt.Errorf("Role with ID %s not found", roleName)
		}

		// Assign the corresponding role to the user using Casbin enforcer
		if _, err := e.AddRoleForUser(strconv.Itoa(int(user.ID)), strconv.Itoa(int(role.ID))); err != nil {
			return fmt.Errorf("Failed to assign role ID %d to user %s: %v", roleName, userEmail, err)
		}
	}

	return nil
}
