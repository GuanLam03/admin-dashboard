package migrations

import (
	"github.com/goravel/framework/contracts/database/schema"
	"github.com/goravel/framework/facades"
)

type M20250930104236CreateGmailTemplatesTable struct{}

// Signature The unique signature for the migration.
func (r *M20250930104236CreateGmailTemplatesTable) Signature() string {
	return "20250930104236_create_gmail_templates_table"
}

// Up Run the migrations.
func (r *M20250930104236CreateGmailTemplatesTable) Up() error {
	if !facades.Schema().HasTable("gmail_templates") {
		return facades.Schema().Create("gmail_templates", func(table schema.Blueprint) {
			table.ID()
			table.String("team", 50) 
			table.String("name", 150)     
			table.Text("content") 
			table.TimestampsTz()

			
		})
	}

	return nil
}

// Down Reverse the migrations.
func (r *M20250930104236CreateGmailTemplatesTable) Down() error {
 	return facades.Schema().DropIfExists("gmail_templates")
}
