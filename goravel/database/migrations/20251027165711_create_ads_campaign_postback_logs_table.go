package migrations

import (
	"github.com/goravel/framework/contracts/database/schema"
	"github.com/goravel/framework/facades"
)

type M20251027165711CreateAdsCampaignPostbackLogsTable struct{}

// Signature The unique signature for the migration.
func (r *M20251027165711CreateAdsCampaignPostbackLogsTable) Signature() string {
	return "20251027165711_create_ads_campaign_postback_logs_table"
}

// Up Run the migrations.
func (r *M20251027165711CreateAdsCampaignPostbackLogsTable) Up() error {
	if !facades.Schema().HasTable("ads_campaign_postback_logs") {
		return facades.Schema().Create("ads_campaign_postback_logs", func(table schema.Blueprint) {
			table.ID()
			table.UnsignedBigInteger("ads_event_log_id").Nullable()
			table.UnsignedBigInteger("ads_campaign_postback_id")
			table.Text("url")
			table.String("request_method")
			table.Text("request_body").Nullable()
			table.Integer("response_status").Nullable()
			table.Text("response_body").Nullable()
			table.Text("error_message").Nullable()
			table.TimestampsTz()
		})
	}

	return nil
}

// Down Reverse the migrations.
func (r *M20251027165711CreateAdsCampaignPostbackLogsTable) Down() error {
 	return facades.Schema().DropIfExists("ads_campaign_postback_logs")
}
