package migrations

import (
	"github.com/goravel/framework/contracts/database/schema"
	"github.com/goravel/framework/facades"
)

type M20251016142201UpdateAdsLogsTable struct{}

// Signature The unique signature for the migration.
func (r *M20251016142201UpdateAdsLogsTable) Signature() string {
	return "20251016142201_update_ads_logs_table"
}

// Up Run the migrations.
func (r *M20251016142201UpdateAdsLogsTable) Up() error {
	return facades.Schema().Table("ads_logs", func(table schema.Blueprint) {
		table.DropColumn("clicked_url")
	})
}

// Down Reverse the migrations.
func (r *M20251016142201UpdateAdsLogsTable) Down() error {
	return facades.Schema().Table("ads_logs", func(table schema.Blueprint) {
		table.String("clicked_url").After("ads_log_detail_id")
	})
}
