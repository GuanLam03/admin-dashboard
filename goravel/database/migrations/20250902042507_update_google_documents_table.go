package migrations
import (
	"github.com/goravel/framework/contracts/database/schema"
	"github.com/goravel/framework/facades"
)
type M20250902042507UpdateGoogleDocumentsTable struct{}

// Signature The unique signature for the migration.
func (r *M20250902042507UpdateGoogleDocumentsTable) Signature() string {
	return "20250902042507_update_google_documents_table"
}

// Up Run the migrations.
func (r *M20250902042507UpdateGoogleDocumentsTable) Up() error {
	return facades.Schema().Table("google_documents", func(table schema.Blueprint) {
		// Change 'link' column to TEXT
		table.Text("link").Change()

		// Add 'original_link' column as TEXT after 'name'
		table.Text("original_link").After("name")
	})
}

// Down Reverse the migrations.
func (r *M20250902042507UpdateGoogleDocumentsTable) Down() error {
	return facades.Schema().Table("google_documents", func(table schema.Blueprint) {
		// Revert 'link' back to string (you can adjust the length if needed)
		table.String("link", 255).Change()

		// Drop the 'original_link' column
		table.DropColumn("original_link")
	})
}
