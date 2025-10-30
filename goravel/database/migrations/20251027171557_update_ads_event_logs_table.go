package migrations

import (
	"github.com/goravel/framework/facades"
)
type M20251027171557UpdateAdsEventLogsTable struct{}

// Signature The unique signature for the migration.
func (r *M20251027171557UpdateAdsEventLogsTable) Signature() string {
	return "20251027171557_update_ads_event_logs_table"
}

// Up Run the migrations.
func (r *M20251027171557UpdateAdsEventLogsTable) Up() error {
	sql := `
		ALTER TABLE ads_event_logs
		ADD COLUMN value_extracted DECIMAL(10,2)
		GENERATED ALWAYS AS (CAST(JSON_UNQUOTE(JSON_EXTRACT(data, '$.value')) AS DECIMAL(10,2))) STORED
	`
	return facades.DB().Statement(sql)
}
// Down Reverse the migrations.
func (r *M20251027171557UpdateAdsEventLogsTable) Down() error {
	sql := `ALTER TABLE ads_event_logs DROP COLUMN IF EXISTS value_extracted`
	return facades.DB().Statement(sql)
}
