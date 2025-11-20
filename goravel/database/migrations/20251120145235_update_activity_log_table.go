package migrations

import (
	"github.com/goravel/framework/contracts/database/schema"
	"github.com/goravel/framework/facades"
)



type M20251120145235UpdateActivityLogTable struct{}

// Signature The unique signature for the migration.
func (r *M20251120145235UpdateActivityLogTable) Signature() string {
	return "20251120145235_update_activity_log_table"
}

// Up Run the migrations.
func (r *M20251120145235UpdateActivityLogTable) Up() error {
	return facades.Schema().Table("activity_logs", func(table schema.Blueprint) {

		table.Json("request_meta").Nullable().After("description")

	})
}

// Down Reverse the migrations.
func (r *M20251120145235UpdateActivityLogTable) Down() error {
	return facades.Schema().Table("activity_logs", func(table schema.Blueprint) {
		table.DropColumn("request_meta")


	})
}
