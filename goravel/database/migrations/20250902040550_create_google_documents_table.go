package migrations

import (
	"github.com/goravel/framework/contracts/database/schema"
	"github.com/goravel/framework/facades"
)

type M20250902040550CreateGoogleDocumentsTable struct{}

// Signature The unique signature for the migration.
func (r *M20250902040550CreateGoogleDocumentsTable) Signature() string {
	return "20250902040550_create_google_documents_table"
}

// Up Run the migrations.
func (r *M20250902040550CreateGoogleDocumentsTable) Up() error {
	if !facades.Schema().HasTable("google_documents") {
		return facades.Schema().Create("google_documents", func(table schema.Blueprint) {
			table.ID()
			table.String("name", 50)
			table.String("link", 255)
			table.String("status", 20).Default("active")
			table.TimestampsTz()
		})
	}

	return nil
}

// Down Reverse the migrations.
func (r *M20250902040550CreateGoogleDocumentsTable) Down() error {
 	return facades.Schema().DropIfExists("google_documents")
}
