package migrations

import (
	"github.com/goravel/framework/contracts/database/schema"
	"github.com/goravel/framework/facades"
)

type M20251016142454UpdateAdsLogDetailsTable struct{}

// Signature The unique signature for the migration.
func (r *M20251016142454UpdateAdsLogDetailsTable) Signature() string {
	return "20251016142454_update_ads_log_details_table"
}

// Up Run the migrations.
func (r *M20251016142454UpdateAdsLogDetailsTable) Up() error {
	return facades.Schema().Table("ads_log_details", func(table schema.Blueprint) {
		table.UnsignedBigInteger("ads_campaign_id").After("id") 
		table.String("clicked_url").After("referrer")
	})
	
}

// Down Reverse the migrations.
func (r *M20251016142454UpdateAdsLogDetailsTable) Down() error {
	return facades.Schema().Table("ads_log_details", func(table schema.Blueprint) {
		table.DropColumn("ads_campaign_id")
		table.DropColumn("clicked_url")
	})

	
}
