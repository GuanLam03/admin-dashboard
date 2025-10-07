package migrations

import (
	"github.com/goravel/framework/contracts/database/schema"
	"github.com/goravel/framework/facades"
)

type M20251002161658CreateAdsCampaignsTable struct{}

// Signature The unique signature for the migration.
func (r *M20251002161658CreateAdsCampaignsTable) Signature() string {
	return "20251002161658_create_ads_campaigns_table"
}

// Up Run the migrations.
func (r *M20251002161658CreateAdsCampaignsTable) Up() error {
	if !facades.Schema().HasTable("ads_campaigns") {
		return facades.Schema().Create("ads_campaigns", func(table schema.Blueprint) {
			table.ID()
			table.String("name")                
			table.String("target_url")            
			table.String("code")    
			table.TimestampsTz()

			// unique constraint
			table.Unique("code")
		})
	}

	return nil
}

// Down Reverse the migrations.
func (r *M20251002161658CreateAdsCampaignsTable) Down() error {
 	return facades.Schema().DropIfExists("ads_campaigns")
}
