package migrations
import (
	"github.com/goravel/framework/contracts/database/schema"
	"github.com/goravel/framework/facades"
)

type M20251009171322UpdateAdsLogsTable struct{}

// Signature The unique signature for the migration.
func (r *M20251009171322UpdateAdsLogsTable) Signature() string {
	return "20251009171322_update_ads_logs_table"
}

// Up Run the migrations.
func (r *M20251009171322UpdateAdsLogsTable) Up() error {
	return facades.Schema().Table("ads_logs", func(table schema.Blueprint) {
		table.DropColumn("ip")
		table.DropColumn("country")
		table.DropColumn("region")
		table.DropColumn("city")
		table.DropColumn("user_agent")
		table.DropColumn("referrer")
		table.DropColumn("converted")
		table.DropColumn("client_product_id")
		table.DropColumn("value")

		table.UnsignedBigInteger("ads_log_detail_id").Nullable().After("ads_campaign_id")
		table.String("clicked_url").After("ads_log_detail_id")

	})
}

// Down Reverse the migrations.
func (r *M20251009171322UpdateAdsLogsTable) Down() error {
	return facades.Schema().Table("ads_logs", func(table schema.Blueprint) {
		table.String("ip").Nullable()
		table.String("country").Nullable()
		table.String("region").Nullable()
		table.String("city").Nullable()
		table.String("user_agent").Nullable()
		table.String("referrer").Nullable()
		table.Boolean("converted").Default(false)
		table.UnsignedBigInteger("client_product_id").Nullable()
		table.Decimal("value").Nullable()


		table.DropColumn("ads_log_detail_id")
		table.DropColumn("clicked_url")
	})
}
