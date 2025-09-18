package seeders

import (
	"goravel/app/models"
	"github.com/goravel/framework/facades"
)

// GoogleDocumentSeeder is a struct for seeding google document data.
type GoogleDocumentSeeder struct{}

// Signature returns the name of the seeder.
func (s *GoogleDocumentSeeder) Signature() string {
	return "GoogleDocumentSeeder"
}

// Run executes the seeder logic.
func (s *GoogleDocumentSeeder) Run() error {
	
	document := models.GoogleDocument{
		Name:          "Datatable js",
		OriginalLink:  "https://docs.google.com/document/d/1gSSAzjU9JdJYfKRPkHyOFQY9iHAsaEvC4f9nPnu5gRM/edit?usp=sharing",
		Link:          "https://docs.google.com/document/d/1gSSAzjU9JdJYfKRPkHyOFQY9iHAsaEvC4f9nPnu5gRM/preview",
		Status:        "active",
	}

	if err := facades.Orm().Query().Create(&document); err != nil {
		return err 
	}

	return nil 
}
