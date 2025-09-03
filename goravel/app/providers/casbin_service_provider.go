package providers

import (
	"log"

	"github.com/casbin/casbin/v2"
	"github.com/casbin/casbin/v2/util"

	gormadapter "github.com/casbin/gorm-adapter/v3"
	"github.com/goravel/framework/contracts/foundation"
	// "github.com/goravel/framework/facades"
)

const CasbinBinding = "casbin"

type CasbinServiceProvider struct{}

// Register binds the Casbin enforcer singleton to the Goravel service container
func (c *CasbinServiceProvider) Register(app foundation.Application) {
	app.Singleton(CasbinBinding, func(app foundation.Application) (any, error) {
		// 2. Initialize the Gorm adapter (will auto-create "casbin_rule" table if not exists)
		adapter, err := gormadapter.NewAdapter("mysql", "root:@tcp(127.0.0.1:3306)/testgoravel?charset=utf8mb4&parseTime=True",true)

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
func (c *CasbinServiceProvider) Boot(app foundation.Application) {}
