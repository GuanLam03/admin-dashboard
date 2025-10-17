package migrations

import (
	"github.com/goravel/framework/contracts/database/schema"
	"github.com/goravel/framework/facades"
)

type M20251009193011CreateAdsLogDetailsTable struct{}

// Signature The unique signature for the migration.
func (r *M20251009193011CreateAdsLogDetailsTable) Signature() string {
	return "20251009193011_create_ads_log_details_table"
}

// Up Run the migrations.
func (r *M20251009193011CreateAdsLogDetailsTable) Up() error {
	if !facades.Schema().HasTable("ads_log_details") {
		return facades.Schema().Create("ads_log_details", func(table schema.Blueprint) {
			table.ID()
			table.String("ip").Nullable()
			table.String("country").Nullable()
			table.String("region").Nullable()
			table.String("city").Nullable()
			table.String("user_agent").Nullable()
			table.String("device_type").Nullable()    
			table.String("device_name").Nullable()     
			table.String("os_name").Nullable()         
			table.String("os_version").Nullable()      
			table.String("browser_name").Nullable()    
			table.String("browser_version").Nullable() 
			table.String("referrer").Nullable()
			table.TimestampsTz()
		})
	}

	return nil
}

// Down Reverse the migrations.
func (r *M20251009193011CreateAdsLogDetailsTable) Down() error {
 	return facades.Schema().DropIfExists("ads_log_details")
}
