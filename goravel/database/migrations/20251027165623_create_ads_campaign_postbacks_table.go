package migrations

import (
	"github.com/goravel/framework/contracts/database/schema"
	"github.com/goravel/framework/facades"
)

type M20251027165623CreateAdsCampaignPostbacksTable struct{}

// Signature The unique signature for the migration.
func (r *M20251027165623CreateAdsCampaignPostbacksTable) Signature() string {
	return "20251027165623_create_ads_campaign_postbacks_table"
}

// Up Run the migrations.
func (r *M20251027165623CreateAdsCampaignPostbacksTable) Up() error {
	if !facades.Schema().HasTable("ads_campaign_postbacks") {
		return facades.Schema().Create("ads_campaign_postbacks", func(table schema.Blueprint) {
			table.ID()
			table.UnsignedBigInteger("ads_campaign_id")
			table.Enum("event_name",[]any{
				"ADD_PAYMENT_INFO",
				"ADD_TO_CART",
				"BUTTON_CLICK",
				"PURCHASE",
				"CONTENT_VIEW",
				"DOWNLOAD",
				"FORM_SUBMIT",
				"INITIATED_CHECKOUT",
				"CONTACT",
				"PLACE_ORDER",
				"SEARCH",
				"COMPLETE_REGISTRATION",
				"ADD_TO_WISHLIST",
				"SUBSCRIBE",
				"FIRST_DEPOSIT",
				"FIRST_DAY_PURCHASE",
			})
			table.Text("postback_url")
			table.TimestampsTz()
		})
	}

	return nil
}

// Down Reverse the migrations.
func (r *M20251027165623CreateAdsCampaignPostbacksTable) Down() error {
 	return facades.Schema().DropIfExists("ads_campaign_postbacks")
}
