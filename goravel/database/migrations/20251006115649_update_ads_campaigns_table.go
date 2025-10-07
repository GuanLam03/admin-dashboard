package migrations

import (
	"github.com/goravel/framework/contracts/database/schema"
	"github.com/goravel/framework/facades"
)

type M20251006115649UpdateAdsCampaignsTable struct{}

// Signature The unique signature for the migration.
func (r *M20251006115649UpdateAdsCampaignsTable) Signature() string {
	return "20251006115649_update_ads_campaigns_table"
}

// Up Run the migrations.
func (r *M20251006115649UpdateAdsCampaignsTable) Up() error {
	return facades.Schema().Table("ads_campaigns", func(table schema.Blueprint) {
		table.Text("tracking_link").Nullable().After("code")
		table.Text("postback_link").Nullable().After("tracking_link")

	})
}

// Down Reverse the migrations.
func (r *M20251006115649UpdateAdsCampaignsTable) Down() error {
	return facades.Schema().Table("ads_campaigns", func(table schema.Blueprint) {
		table.DropColumn("tracking_link")
		table.DropColumn("postback_link")
	})
}
