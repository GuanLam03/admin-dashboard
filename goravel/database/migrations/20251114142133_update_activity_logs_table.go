package migrations

import (
	"github.com/goravel/framework/contracts/database/schema"
	"github.com/goravel/framework/facades"
)


type M20251114142133UpdateActivityLogsTable struct{}

// Signature The unique signature for the migration.
func (r *M20251114142133UpdateActivityLogsTable) Signature() string {
	return "20251114142133_update_activity_logs_table"
}

// Up Run the migrations.
func (r *M20251114142133UpdateActivityLogsTable) Up() error {
	return facades.Schema().Table("activity_logs", func(table schema.Blueprint) {

		table.Timestamp("start_at", 3).Nullable().After("description")
		table.Timestamp("end_at", 3).Nullable().After("start_at")


	})

}

// Down Reverse the migrations.
func (r *M20251114142133UpdateActivityLogsTable) Down() error {
	return facades.Schema().Table("activity_logs", func(table schema.Blueprint) {
		table.DropColumn("start_at")
		table.DropColumn("end_at")

	})
}
