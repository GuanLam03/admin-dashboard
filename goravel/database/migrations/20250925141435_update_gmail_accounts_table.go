package migrations

import (
	"github.com/goravel/framework/contracts/database/schema"
	"github.com/goravel/framework/facades"
)

type M20250925141435UpdateGmailAccountsTable struct{}

// Signature The unique signature for the migration.
func (r *M20250925141435UpdateGmailAccountsTable) Signature() string {
	return "20250925141435_update_gmail_accounts_table"
}

// Up Run the migrations.
func (r *M20250925141435UpdateGmailAccountsTable) Up() error {
	return facades.Schema().Table("gmail_accounts", func(table schema.Blueprint) {
		table.String("team",50).Nullable().After("email")
	})
}

// Down Reverse the migrations.
func (r *M20250925141435UpdateGmailAccountsTable) Down() error {
	return facades.Schema().Table("gmail_accounts", func(table schema.Blueprint) {
        table.DropColumn("team")

    })
}
