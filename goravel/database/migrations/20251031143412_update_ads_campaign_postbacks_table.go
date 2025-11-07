package migrations

import (
	"github.com/goravel/framework/contracts/database/schema"
	"github.com/goravel/framework/facades"
)

type M20251031143412UpdateAdsCampaignPostbacksTable struct{}

// Signature The unique signature for the migration.
func (r *M20251031143412UpdateAdsCampaignPostbacksTable) Signature() string {
	return "20251031143412_update_ads_campaign_postbacks_table"
}

// Up Run the migrations.
func (r *M20251031143412UpdateAdsCampaignPostbacksTable) Up() error {
	return facades.Schema().Table("ads_campaign_postbacks", func(table schema.Blueprint) {

		table.Boolean("include_click_params").After("postback_url").Default(false)
	})
}

// Down Reverse the migrations.
func (r *M20251031143412UpdateAdsCampaignPostbacksTable) Down() error {
	return facades.Schema().Table("ads_campaign_postbacks", func(table schema.Blueprint) {
		table.DropColumn("include_click_params")
	})
}
