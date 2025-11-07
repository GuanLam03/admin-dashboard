package migrations

import (
	"github.com/goravel/framework/contracts/database/schema"
	"github.com/goravel/framework/facades"
)


type M20251106124040UpdateAdsCampaignPostbackLogsTable struct{}

// Signature The unique signature for the migration.
func (r *M20251106124040UpdateAdsCampaignPostbackLogsTable) Signature() string {
	return "20251106124040_update_ads_campaign_postback_logs_table"
}

// Up Run the migrations.
func (r *M20251106124040UpdateAdsCampaignPostbackLogsTable) Up() error {
	return facades.Schema().Table("ads_campaign_postback_logs", func(table schema.Blueprint) {

		table.String("status",50).After("error_message")
	})
}

// Down Reverse the migrations.
func (r *M20251106124040UpdateAdsCampaignPostbackLogsTable) Down() error {
	return facades.Schema().Table("ads_campaign_postback_logs", func(table schema.Blueprint) {
		table.DropColumn("status")
	})
}
