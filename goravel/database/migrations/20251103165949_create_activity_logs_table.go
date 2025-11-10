package migrations

import (
	"github.com/goravel/framework/contracts/database/schema"
	"github.com/goravel/framework/facades"
)

type M20251103165949CreateActivityLogsTable struct{}

// Signature The unique signature for the migration.
func (r *M20251103165949CreateActivityLogsTable) Signature() string {
	return "20251103165949_create_activity_logs_table"
}

// Up Run the migrations.
func (r *M20251103165949CreateActivityLogsTable) Up() error {
	if !facades.Schema().HasTable("activity_logs") {
		return facades.Schema().Create("activity_logs", func(table schema.Blueprint) {
			table.ID()
			table.UnsignedBigInteger("causer_id").Nullable()
			table.String("causer_type").Nullable()
			table.Text("properties").Nullable() // array as text
			table.Text("url").Nullable()
			table.Text("route").Nullable()
			table.Text("input").Nullable() // JSON as text
			table.String("log_name").Nullable()
			table.String("description").Nullable()
			table.TimestampsTz()
		})
		
	}

	return nil
}

// Down Reverse the migrations.
func (r *M20251103165949CreateActivityLogsTable) Down() error {
 	return facades.Schema().DropIfExists("activity_logs")
}
