package migrations

import (
	"github.com/goravel/framework/contracts/database/schema"
	"github.com/goravel/framework/facades"
)

type M20251009152316CreateAdsEventLogsTable struct{}

// Signature The unique signature for the migration.
func (r *M20251009152316CreateAdsEventLogsTable) Signature() string {
	return "20251009152316_create_ads_event_logs_table"
}

// Up Run the migrations.
func (r *M20251009152316CreateAdsEventLogsTable) Up() error {
	if !facades.Schema().HasTable("ads_event_logs") {
		return facades.Schema().Create("ads_event_logs", func(table schema.Blueprint) {
			table.ID()
			table.UnsignedBigInteger("ads_log_id")
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
			table.Json("data")
			table.TimestampsTz()
		})
	}

	return nil
}

// Down Reverse the migrations.
func (r *M20251009152316CreateAdsEventLogsTable) Down() error {
 	return facades.Schema().DropIfExists("ads_event_logs")
}
