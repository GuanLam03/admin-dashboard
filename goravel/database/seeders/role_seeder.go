package seeders

import (
	"goravel/app/models"
	"github.com/goravel/framework/facades"
)

// RoleSeeder is a struct for seeding role data.
type RoleSeeder struct{}

// Signature returns the name of the seeder.
func (s *RoleSeeder) Signature() string {
	return "RoleSeeder"
}

// Run executes the seeder logic.
func (s *RoleSeeder) Run() error {

	roles := []models.Role{
		{Name: "reader"},
		{Name: "developer"},
		{Name: "super admin"},
	}

	for _, role := range roles {
		if err := facades.Orm().Query().Create(&role); err != nil {
			return err 
		}
	}

	return nil
}
