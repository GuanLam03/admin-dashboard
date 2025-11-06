package migrations

import (
	"github.com/goravel/framework/contracts/database/schema"
	"github.com/goravel/framework/facades"
)

type M20251027170307UpdateAdsLogDetailsTable struct{}

// Signature The unique signature for the migration.
func (r *M20251027170307UpdateAdsLogDetailsTable) Signature() string {
	return "20251027170307_update_ads_log_details_table"
}

// Up Run the migrations.
func (r *M20251027170307UpdateAdsLogDetailsTable) Up() error {
	return facades.Schema().Table("ads_log_details", func(table schema.Blueprint) {
		table.DropColumn("clicked_url")
		table.Text("clicked_url").After("referrer")
	})
}

// Down Reverse the migrations.
func (r *M20251027170307UpdateAdsLogDetailsTable) Down() error {
	return facades.Schema().Table("ads_log_details", func(table schema.Blueprint) {
		table.DropColumn("clicked_url")
		table.String("clicked_url").After("referrer")
	})
}
