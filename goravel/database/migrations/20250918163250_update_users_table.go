package migrations

import (
	"github.com/goravel/framework/contracts/database/schema"
	"github.com/goravel/framework/facades"
)
type M20250918163250UpdateUsersTable struct{}

// Signature The unique signature for the migration.
func (r *M20250918163250UpdateUsersTable) Signature() string {
	return "20250918163250_update_users_table"
}

// Up Run the migrations.
func (r *M20250918163250UpdateUsersTable) Up() error {
	return facades.Schema().Table("users", func(table schema.Blueprint) {
		
		table.Boolean("two_factor_enabled").Default(false)
		table.String("two_factor_secret", 512).Nullable()

	})

	
}

// Down Reverse the migrations.
func (r *M20250918163250UpdateUsersTable) Down() error {
	return facades.Schema().Table("users", func(table schema.Blueprint) {

        table.DropColumn("two_factor_enabled")
		table.DropColumn("two_factor_secret")

    })
}
