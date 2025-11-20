package migrations

import "github.com/goravel/framework/facades"

type M20251118152258UpdateAdsEventLogsTable struct{}

// Signature The unique signature for the migration.
func (r *M20251118152258UpdateAdsEventLogsTable) Signature() string {
	return "20251118152258_update_ads_event_logs_table"
}

// Up Run the migrations.
func (r *M20251118152258UpdateAdsEventLogsTable) Up() error {
	sql := `
ALTER TABLE ads_event_logs
  ADD INDEX idx_ael_log_event (ads_log_id, event_name),
  ADD INDEX idx_ael_event_log_value (event_name, ads_log_id, value_extracted),
  ALGORITHM=INPLACE,
  LOCK=NONE;
`
	return facades.DB().Statement(sql)
}

// Down Reverse the migrations.
func (r *M20251118152258UpdateAdsEventLogsTable) Down() error {
	sql := `
ALTER TABLE ads_event_logs
  DROP INDEX idx_ael_log_event,
  DROP INDEX idx_ael_event_log_value;
`
	return facades.DB().Statement(sql)
}
