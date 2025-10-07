package migrations

import (
	"github.com/goravel/framework/contracts/database/schema"
	"github.com/goravel/framework/facades"
)

type M20251002205238CreateAdsLogsTable struct{}

// Signature The unique signature for the migration.
func (r *M20251002205238CreateAdsLogsTable) Signature() string {
	return "20251002205238_create_ads_logs_table"
}

// Up Run the migrations.
func (r *M20251002205238CreateAdsLogsTable) Up() error {
	if !facades.Schema().HasTable("ads_logs") {
		return facades.Schema().Create("ads_logs", func(table schema.Blueprint) {
			table.ID()
			table.UnsignedBigInteger("ads_campaign_id")    
			
			table.String("ip").Nullable()
			table.String("country").Nullable()
			table.String("region").Nullable()
			table.String("city").Nullable()
			table.String("user_agent").Nullable()
			table.String("referrer").Nullable()
		
			table.Boolean("converted").Default(false)
			table.UnsignedBigInteger("client_product_id").Nullable()  
			table.Decimal("value").Nullable() 
			table.TimestampsTz()
		})
	}

	return nil
}

// Down Reverse the migrations.
func (r *M20251002205238CreateAdsLogsTable) Down() error {
 	return facades.Schema().DropIfExists("ads_logs")
}
