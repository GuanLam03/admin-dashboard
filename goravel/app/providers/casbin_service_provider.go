package providers

import (
	"log"
	"fmt"
	"github.com/casbin/casbin/v2"
	"github.com/casbin/casbin/v2/util"

	gormadapter "github.com/casbin/gorm-adapter/v3"
	"github.com/goravel/framework/contracts/foundation"
	"github.com/goravel/framework/facades"
)

const CasbinBinding = "casbin"

type CasbinServiceProvider struct{}

// Register binds the Casbin enforcer singleton to the Goravel service container
func (c *CasbinServiceProvider) Register(app foundation.Application) {
	//DB_DSN is docker else use local databse
	dbDSN := facades.Config().Env("DB_DSN", "")

	if dbDSN == "" {
		dbName := facades.Config().Env("DB_DATABASE", "").(string)
		dbDSN = fmt.Sprintf("root:@tcp(127.0.0.1:3306)/%s?charset=utf8mb4&parseTime=True", dbName)
	}

	app.Singleton(CasbinBinding, func(app foundation.Application) (any, error) {
		// 2. Initialize the Gorm adapter (will auto-create "casbin_rule" table if not exists)
		adapter, err := gormadapter.NewAdapter("mysql",dbDSN.(string),true)

		if err != nil {
			log.Fatalf("failed to create Casbin adapter: %v", err)
		}

		// 3. Initialize the Casbin enforcer with your RBAC model config
		enforcer, err := casbin.NewEnforcer("config/casbin_model.conf", adapter)
		if err != nil {
			log.Fatalf("failed to create Casbin enforcer: %v", err)
		}
		enforcer.AddFunction("keyMatch2",util.KeyMatch2Func)
		enforcer.EnableAutoSave(true)
		enforcer.EnableAutoBuildRoleLinks(true)

		// 4. Load policy rules from DB
		if err := enforcer.LoadPolicy(); err != nil {
			log.Fatalf("failed to load Casbin policies: %v", err)
		}

		return enforcer, nil
	})
}

// Boot runs after all providers are registered (optional)
func (c *CasbinServiceProvider) Boot(app foundation.Application) {
	 _, _ = app.Make(CasbinBinding)
}
