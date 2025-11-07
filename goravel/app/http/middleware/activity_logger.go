package middleware

import (
	"encoding/json"
	// "fmt"
	"github.com/goravel/framework/contracts/http"
	"github.com/goravel/framework/facades"
	"goravel/app/models"
)

func ActivityLogger() http.Middleware {
	return func(ctx http.Context) {
		method := ctx.Request().Method()
		path := ctx.Request().Path()
		url := ctx.Request().FullUrl()
		inputData := ctx.Request().All()
		// fmt.Println("input: ",inputData) // input:  map[id:46 name:test purchase postback_events:[map[event_name:PURCHASE id:18 include_click_params:true postback_url:http://localhost:8080/ffoff222] map[event_name:FORM_SUBMIT id:19 include_click_params:true postback_url:https://webhook.site/cf8aae5c-4078-46a2-b0dd-ba3ae23db740?mcid={click_id}&value={value}&content_id={content_id}]] postback_link:http://localhost:3000/postback/ status:active target_url:http://localhost:8080/ tracking_link:http://localhost:3000/ENym2FWD]
		bodyJSON, _ := json.Marshal(inputData)

		// Skip read-only requests
		if method == "GET" {
			ctx.Request().Next()
			return
		}

		// Execute the controller
		ctx.Request().Next()

	
		statusCode := ctx.Response().Origin().Status()

		// Only log successful responses (status codes 200–299)
		if statusCode < 200 || statusCode >= 300 {
			return
		}

		// Get authenticated user
		var user models.User
		_ = facades.Auth(ctx).User(&user)

		facades.Orm().Query().Create(&models.ActivityLog{
			CauserId:    user.ID,
			CauserType:  "",
			Properties:  "",
			Url:         url,
			Route:       method + " " + path,
			Input :  	string(bodyJSON),
			LogName:     "",
			Description: "",
		})
	}
}
