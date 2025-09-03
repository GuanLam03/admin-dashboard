package middleware

import (
	"fmt"
    "github.com/casbin/casbin/v2"
    "github.com/goravel/framework/contracts/http"
    "github.com/goravel/framework/facades"
    "goravel/app/models"
)

func CasbinMiddleware() http.Middleware {
    return func(ctx http.Context) {
        var user models.User
        if err := facades.Auth(ctx).User(&user); err != nil {
			ctx.Response().String(http.StatusUnauthorized, "Unauthorized").Abort()
			return
        }

        enforcerAny, err := facades.App().Make("casbin")
		if err != nil {
			ctx.Response().String(http.StatusInternalServerError, "Failed to get Casbin enforcer").Abort()
			return
		}

		enforcer := enforcerAny.(*casbin.Enforcer)

        sub := fmt.Sprintf("%d", user.ID) 
        obj := ctx.Request().Path()      // URL
        act := ctx.Request().Method()    // HTTP method

        ok, _ := enforcer.Enforce(sub, obj, act)
        if !ok {
			 ctx.Response().String(http.StatusUnauthorized, fmt.Sprintf("User ID: %d is forbidden", sub)).Abort()
        }

		fmt.Println("CasbinMiddleware invoked")

		ctx.Request().Next()
    }
}
