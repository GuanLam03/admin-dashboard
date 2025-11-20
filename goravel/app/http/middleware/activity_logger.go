package middleware

import (
	"encoding/json"
	"time"

	ua "github.com/mssola/useragent"
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

		

		ip := ctx.Request().Ip()
		rawUA := ctx.Request().Header("User-Agent")
		referrer := ctx.Request().Header("Referer")

		// Parse User-Agent deeply
		uaParser := ua.New(rawUA)

		browserName, browserVersion := uaParser.Browser()

		type MetaData struct {
			Ip            string `json:"ip"`
			RawUserAgent  string `json:"raw_user_agent"`
			Browser       string `json:"browser"`
			BrowserName   string `json:"browser_name"`
			BrowserVer    string `json:"browser_version"`
			OS            string `json:"os"`
			Platform      string `json:"platform"`
			Mobile        bool   `json:"mobile"`
			Bot           bool   `json:"bot"`
			Referrer      string `json:"referrer"`
		}

		meta := MetaData{
			Ip:           ip,
			RawUserAgent: rawUA,
			Browser:      browserName + " " + browserVersion,
			BrowserName:  browserName,
			BrowserVer:   browserVersion,
			OS:           uaParser.OS(),
			Platform:     uaParser.Platform(),
			Mobile:       uaParser.Mobile(),
			Bot:          uaParser.Bot(),
			Referrer:     referrer,
		}

		metaJSON, _ := json.Marshal(meta)

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
			RequestMeta: metaJSON,
			StartAt:     &startTime,
			EndAt:       &endTime,
		})
	}
}
