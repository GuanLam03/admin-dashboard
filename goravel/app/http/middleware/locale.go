package middleware

import (
	"github.com/goravel/framework/contracts/http"
	"github.com/goravel/framework/facades"
	"goravel/app/helpers/system"
	"fmt"
)

func Locale() http.Middleware {
	return func(ctx http.Context) {
		lang := ctx.Request().Header("Accept-Language")

        if _, exists := system.Languages[lang]; !exists {
            lang = facades.Config().GetString("app.locale")
        }
		
        facades.App().SetLocale(ctx, lang)
		
		fmt.Println("current lang:",facades.App().CurrentLocale(ctx))
        ctx.Request().Next()
	}
}
