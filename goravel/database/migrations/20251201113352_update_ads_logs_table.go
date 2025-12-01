package migrations


import "github.com/goravel/framework/facades"

type M20251201113352UpdateAdsLogsTable struct{}

// Signature The unique signature for the migration.
func (r *M20251201113352UpdateAdsLogsTable) Signature() string {
	return "20251201113352_update_ads_logs_table"
}

// Up Run the migrations.
func (r *M20251201113352UpdateAdsLogsTable) Up() error {
	sql := `
		ALTER TABLE ads_logs
		ADD INDEX idx_ads_campaign_id (ads_campaign_id);
	`
	return facades.DB().Statement(sql)
}

// Down Reverse the migrations.
func (r *M20251201113352UpdateAdsLogsTable) Down() error {
	sql := `
		ALTER TABLE ads_logs
		DROP INDEX idx_ads_campaign_id;

	`	
	return facades.DB().Statement(sql)
}
