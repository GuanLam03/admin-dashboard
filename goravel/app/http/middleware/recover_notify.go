package middleware

import (
	"fmt"
	"runtime/debug"
	 "time"

	"github.com/goravel/framework/contracts/http"
	"github.com/goravel/framework/facades"
	"github.com/goravel/framework/mail"
)

func RecoverNotify() http.Middleware {
	return func(ctx http.Context) {
		defer func() {
			if r := recover(); r != nil {
				
				method := ctx.Request().Method()
				path := ctx.Request().Path()
				ua := ctx.Request().Header("User-Agent")

				body := fmt.Sprintf(
					`<h2>Goravel Panic Recovered [%s]</h2>
					 <p><b>Error:</b> %v</p>
					 <p><b>Request:</b> %s %s</p>
					 <p><b>User-Agent:</b> %s</p>
					 <pre style="white-space:pre-wrap">%s</pre>`,
					time.Now().Format("2025-09-02 15:04:05"), r, method, path, ua, string(debug.Stack()),
				)

				_ = facades.Mail().
					To([]string{"hgledgetech@gmail.com"}). 
					Subject("Goravel Panic").
					Content(mail.Html(body)).
					Send()

				ctx.Response().
					Json(500, map[string]string{"error": "Internal Server Error"}).
					Abort()
			}
		}()

		ctx.Request().Next()
	}
}
