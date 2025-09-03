package migrations

import (
	"github.com/goravel/framework/contracts/database/schema"
	"github.com/goravel/framework/facades"
)

type M20250828064700CreateDocumentsTable struct{}

// Signature The unique signature for the migration.
func (r *M20250828064700CreateDocumentsTable) Signature() string {
	return "20250828064700_create_documents_table"
}

// Up Run the migrations.
func (r *M20250828064700CreateDocumentsTable) Up() error {
	if !facades.Schema().HasTable("documents") {
		return facades.Schema().Create("documents", func(table schema.Blueprint) {
			table.ID()
			table.String("filename", 255)
			table.String("path", 255)
			table.TimestampsTz()
		})
	}

	return nil
}

// Down Reverse the migrations.
func (r *M20250828064700CreateDocumentsTable) Down() error {
 	return facades.Schema().DropIfExists("documents")
}
