package middleware

import (
	"encoding/json"
	// "fmt"
	"time"
	"github.com/goravel/framework/contracts/http"
	"github.com/goravel/framework/facades"
	"goravel/app/models"
)

func ActivityLogger() http.Middleware {
	return func(ctx http.Context) {

		// --------------------------------------------------
		// START TIME (before controller executes)
        // --------------------------------------------------
		startTime := time.Now()

		method := ctx.Request().Method()
		path := ctx.Request().Path()
		url := ctx.Request().FullUrl()
		inputData := ctx.Request().All()
		bodyJSON, _ := json.Marshal(inputData)

		// Skip read-only requests
		if method == "GET" {
			ctx.Request().Next()
			return
		}

		// --------------------------------------------------
		// Execute the controller
		ctx.Request().Next()
		// --------------------------------------------------

		statusCode := ctx.Response().Origin().Status()

		// Only log successful responses (status codes 200–299)
		if statusCode < 200 || statusCode >= 300 {
			return
		}

		// Get authenticated user
		var user models.User
		_ = facades.Auth(ctx).User(&user)

		// --------------------------------------------------
		// END TIME (after controller returns)
        // --------------------------------------------------
		endTime := time.Now()

		// Save activity log
		facades.Orm().Query().Create(&models.ActivityLog{
			CauserId:    user.ID,
			CauserType:  "",
			Properties:  "",
			Url:         url,
			Route:       method + " " + path,
			Input:       string(bodyJSON),
			LogName:     "",
			Description: "",
			StartAt:     &startTime,
			EndAt:       &endTime,
		})
	}
}
