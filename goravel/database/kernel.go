package database

import (
	"github.com/goravel/framework/contracts/database/schema"
	"github.com/goravel/framework/contracts/database/seeder"
	"goravel/database/migrations"
	"goravel/database/seeders"
)

type Kernel struct {
}

func (kernel Kernel) Migrations() []schema.Migration {
	return []schema.Migration{
		&migrations.M20210101000001CreateUsersTable{},
		&migrations.M20210101000002CreateJobsTable{},
		&migrations.M20250825132604CreateRolesTable{},
		&migrations.M20250825133209CreateUserRolesTable{},
		&migrations.M20250828064700CreateDocumentsTable{},
		&migrations.M20250902040550CreateGoogleDocumentsTable{},
		&migrations.M20250902042507UpdateGoogleDocumentsTable{},
		&migrations.M20250904093240CreateSchedulesTable{},
		&migrations.M20250918163250UpdateUsersTable{},
		&migrations.M20250922151403CreateGmailAccountsTable{},
		&migrations.M20250925141435UpdateGmailAccountsTable{},
		&migrations.M20250930104236CreateGmailTemplatesTable{},
		&migrations.M20251002161658CreateAdsCampaignsTable{},
		&migrations.M20251002205238CreateAdsLogsTable{},
		&migrations.M20251006115649UpdateAdsCampaignsTable{},
		&migrations.M20251009152316CreateAdsEventLogsTable{},
		&migrations.M20251009171322UpdateAdsLogsTable{},
		&migrations.M20251009193011CreateAdsLogDetailsTable{},
		&migrations.M20251016142201UpdateAdsLogsTable{},
		&migrations.M20251016142454UpdateAdsLogDetailsTable{},
		&migrations.M20251027165623CreateAdsCampaignPostbacksTable{},
		&migrations.M20251027165711CreateAdsCampaignPostbackLogsTable{},
		&migrations.M20251027170307UpdateAdsLogDetailsTable{},
		&migrations.M20251027171557UpdateAdsEventLogsTable{},
		&migrations.M20251031143412UpdateAdsCampaignPostbacksTable{},
		&migrations.M20251103165949CreateActivityLogsTable{},
		&migrations.M20251106124040UpdateAdsCampaignPostbackLogsTable{},
		&migrations.M20251114142133UpdateActivityLogsTable{},
		&migrations.M20251118152211UpdateAdsLogDetailsTable{},
		&migrations.M20251118152237UpdateAdsLogsTable{},
		&migrations.M20251118152258UpdateAdsEventLogsTable{},
		&migrations.M20251120145235UpdateActivityLogTable{},
		&migrations.M20251201113352UpdateAdsLogsTable{},
	}
}
func (kernel Kernel) Seeders() []seeder.Seeder {
	return []seeder.Seeder{
		&seeders.DatabaseSeeder{},
		&seeders.ScheduleSeeder{},
		&seeders.UserSeeder{},
		&seeders.RoleSeeder{},
		&seeders.GoogleDocumentSeeder{},
		&seeders.UserRoleSeeder{},
		&seeders.RolePermissionSeeder{},
		// &seeders.AdsTrackingSeeder{},
	}
}
