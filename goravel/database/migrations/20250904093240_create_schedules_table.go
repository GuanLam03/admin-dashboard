package migrations

import (
	"github.com/goravel/framework/contracts/database/schema"
	"github.com/goravel/framework/facades"
)

type M20250904093240CreateSchedulesTable struct{}

// Signature The unique signature for the migration.
func (r *M20250904093240CreateSchedulesTable) Signature() string {
	return "20250904093240_create_schedules_table"
}

// Up Run the migrations.
func (r *M20250904093240CreateSchedulesTable) Up() error {
	if !facades.Schema().HasTable("schedules") {
		return facades.Schema().Create("schedules", func(table schema.Blueprint) {
			table.ID()
			table.String("title")                          // event title, e.g. "Server restart"
			table.String("recurrence").Nullable()          // daily, weekly, monthly, yearly (NULL = one-time)
			table.DateTimeTz("start_at")                  // full datetime (date + time, timezone-aware)
			table.DateTimeTz("end_at")                    // full datetime
			table.String("status").Default("active")       // active / inactive
			table.String("google_event_id").Nullable()     // to sync with Google Calendar
			table.TimestampsTz()                           // created_at & updated_at
		})

	}

	return nil
}

// Down Reverse the migrations.
func (r *M20250904093240CreateSchedulesTable) Down() error {
 	return facades.Schema().DropIfExists("schedules")
}
