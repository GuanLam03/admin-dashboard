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
		// defer func() {
		// 	if r := recover(); r != nil {
		// 		facades.Log().Error(fmt.Sprintf("Recovered from panic: %v\nStack trace:\n%s", r, string(debug.Stack())))
		// 		ctx.Response().String(http.StatusInternalServerError, "Internal Server Error").Abort()
		// 	}
		// }()

		// token := ctx.Request().Header("Authorization", "")
		// if token == "" {
		// 	ctx.Response().String(http.StatusUnauthorized, "Missing token").Abort()
		// 	return
		// }

		// parts := strings.Fields(token)
		// if len(parts) != 2 || strings.ToLower(parts[0]) != "bearer" {
		// 	ctx.Response().String(http.StatusUnauthorized, "Invalid Authorization header").Abort()
		// 	return
		// }

		// tokenString := parts[1]

		// // Use context-aware auth
		// auth := facades.Auth(ctx)
		// if auth == nil {
		// 	ctx.Response().String(http.StatusInternalServerError, "Auth facade is nil").Abort()
		// 	return
		// }

		// // Use .Parse safely
		// user, err := auth.Parse(tokenString)
		// if err != nil {
		// 	ctx.Response().String(http.StatusUnauthorized, "Invalid token: "+err.Error()).Abort()
		// 	return
		// }

		// fmt.Println("Authenticated user:", user)

		token := ctx.Request().Cookie("jwt_token")
		if token == "" {
			ctx.Response().String(http.StatusUnauthorized, "Unauthorized").Abort()
			return
		}

		// Parse & validate token
		user, err := facades.Auth(ctx).Parse(token)
		if err != nil {
			ctx.Response().String(http.StatusUnauthorized, "Invalid token").Abort()
			return
		}

		// Attach user to context
		ctx.WithValue("user", user)
		ctx.Request().Next()
	}
}
