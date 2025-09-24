package migrations

import (
	"github.com/goravel/framework/contracts/database/schema"
	"github.com/goravel/framework/facades"
)

type M20250922151403CreateGmailAccountsTable struct{}

// Signature The unique signature for the migration.
func (r *M20250922151403CreateGmailAccountsTable) Signature() string {
	return "20250922151403_create_gmail_accounts_table"
}

// Up Run the migrations.
func (r *M20250922151403CreateGmailAccountsTable) Up() error {
	if !facades.Schema().HasTable("gmail_accounts") {
		return facades.Schema().Create("gmail_accounts", func(table schema.Blueprint) {
			table.ID()
			table.String("email", 255)
			table.Text("access_token")
			table.Text("refresh_token").Nullable()
			table.Timestamp("expiry").Nullable()
			table.TimestampsTz()
		})
	}

	return nil
}

// Down Reverse the migrations.
func (r *M20250922151403CreateGmailAccountsTable) Down() error {
 	return facades.Schema().DropIfExists("gmail_accounts")
}
