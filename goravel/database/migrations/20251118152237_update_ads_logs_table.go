package migrations

import "github.com/goravel/framework/facades"

type M20251118152237UpdateAdsLogsTable struct{}

// Signature The unique signature for the migration.
func (r *M20251118152237UpdateAdsLogsTable) Signature() string {
	return "20251118152237_update_ads_logs_table"
}

// Up Run the migrations.
func (r *M20251118152237UpdateAdsLogsTable) Up() error {
	sql := `
ALTER TABLE ads_logs
  ADD INDEX idx_al_log_detail_campaign (ads_log_detail_id, ads_campaign_id),
  ADD INDEX idx_al_ads_log_detail_id (ads_log_detail_id),
  ALGORITHM=INPLACE,
  LOCK=NONE;
`
	return facades.DB().Statement(sql)
}

// Down Reverse the migrations.
func (r *M20251118152237UpdateAdsLogsTable) Down() error {
	sql := `
ALTER TABLE ads_logs
  DROP INDEX idx_al_log_detail_campaign,
  DROP INDEX idx_al_ads_log_detail_id;
`
	return facades.DB().Statement(sql)
}
