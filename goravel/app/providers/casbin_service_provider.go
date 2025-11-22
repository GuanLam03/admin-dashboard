package providers

import (
    "fmt"
    "log"

    "github.com/casbin/casbin/v2"
    "github.com/casbin/casbin/v2/util"
    gormadapter "github.com/casbin/gorm-adapter/v3"

    "github.com/goravel/framework/contracts/foundation"
    "github.com/goravel/framework/facades"
)

const CasbinBinding = "casbin"

type CasbinServiceProvider struct{}

func (c *CasbinServiceProvider) Register(app foundation.Application) {
    // Use .Env to retrieve environment variable, with default if not set
    dbDSN := facades.Config().Env("DB_DSN", "").(string)
	
    if dbDSN == "" {
        host := facades.Config().Env("DB_HOST", "127.0.0.1")
        port := facades.Config().Env("DB_PORT", "3306")
        user := facades.Config().Env("DB_USERNAME", "root")
        pass := facades.Config().Env("DB_PASSWORD", "")
        name := facades.Config().Env("DB_DATABASE", "")

        dbDSN = fmt.Sprintf(
            "%s:%s@tcp(%s:%s)/%s?charset=utf8mb4&parseTime=True&loc=Local",
            user, pass, host, port, name,
        )
    }

    app.Singleton(CasbinBinding, func(app foundation.Application) (any, error) {
        adapter, err := gormadapter.NewAdapter("mysql", dbDSN, true)
        if err != nil {
            log.Fatalf("failed to create Casbin adapter: %v", err)
        }

        enforcer, err := casbin.NewEnforcer("config/casbin_model.conf", adapter)
        if err != nil {
            log.Fatalf("failed to create Casbin enforcer: %v", err)
        }

        enforcer.AddFunction("keyMatch2", util.KeyMatch2Func)
        enforcer.EnableAutoSave(true)
        enforcer.EnableAutoBuildRoleLinks(true)

        if err := enforcer.LoadPolicy(); err != nil {
            log.Fatalf("failed to load Casbin policies: %v", err)
        }

        return enforcer, nil
    })
}

func (c *CasbinServiceProvider) Boot(app foundation.Application) {
    _, _ = app.Make(CasbinBinding)
}
