package middleware

import (
	"github.com/goravel/framework/contracts/http"
	"github.com/goravel/framework/facades"

	// "fmt"
	// "strings"
	// "runtime/debug"

	
)

func Auth() http.Middleware {
	return func(ctx http.Context) {

		token := ctx.Request().Cookie("jwt_token")
		if token == "" {
			ctx.Response().String(http.StatusUnauthorized, "login Unauthorized").Abort()
			return
		}

		// Parse & validate token
		_, err := facades.Auth(ctx).Parse(token)
		if err != nil {
			ctx.Response().String(http.StatusUnauthorized, "Invalid token").Abort()
			return
		}

		// Attach user to context
		// ctx.WithValue("user", user)
		// fmt.Println("User: ",user)
		ctx.Request().Next()
	}
}
