package googleCalendar

import (
	"fmt"
	"encoding/json"
	"context"
	"golang.org/x/oauth2"
	"golang.org/x/oauth2/google"
	"github.com/goravel/framework/facades"
    "github.com/goravel/framework/contracts/http"
	"google.golang.org/api/calendar/v3"
	"goravel/app/models"

)
type GoogleCalendarAccountController struct{}

func NewGoogleCalendarAccountController() *GoogleCalendarAccountController {
	return &GoogleCalendarAccountController{}
}

// GET /google/schedule/account
func (r *GoogleCalendarAccountController) ShowGoogleAccount(ctx http.Context) http.Response {
	var account models.GmailAccount
	
	err := facades.Orm().Query().
		Table("gmail_accounts").
		Where("team", "schedule").
		First(&account)

	// If record not found
	if err != nil || account.ID == 0{
		return ctx.Response().Json(404, http.Json{
			"error": "Gmail Not found",
		})
	}

	return ctx.Response().Json(200, http.Json{
		"account": account,
	})
}


func (r *GoogleCalendarAccountController) AuthURL(ctx http.Context) (http.Response) {
	clientID := facades.Config().Env("GOOGLE_CALENDAR_CLIENT_ID", "").(string)
	clientSecret := facades.Config().Env("GOOGLE_CALENDAR_CLIENT_SECRET", "").(string)
	redirectURI := facades.Config().Env("GOOGLE_CALENDAR_REDIRECT_URI", "").(string)

	if clientID == "" || clientSecret == "" || redirectURI == "" {
		return ctx.Response().Json(500, http.Json{
			"error": "Missing Google credentials in .env",
		})
	}

	config := &oauth2.Config{
		ClientID:     clientID,
		ClientSecret: clientSecret,
		RedirectURL:  redirectURI,
		Scopes: []string{
			calendar.CalendarScope,
			"https://www.googleapis.com/auth/userinfo.email",
		},
		Endpoint: google.Endpoint,
	}

	authURL := config.AuthCodeURL("state-token", oauth2.AccessTypeOffline, oauth2.ApprovalForce)

	return ctx.Response().Json(200, http.Json{
		"url": authURL,
	})
}



func (r *GoogleCalendarAccountController) Callback(ctx http.Context) http.Response {
	code := ctx.Request().Query("code")
	facades.Log().Info("Google OAuth callback received", map[string]any{
		"code": code,
	})

	if code == "" {
		facades.Log().Error("Missing code parameter in OAuth callback")
		return ctx.Response().Json(400, http.Json{
			"error": "Missing code parameter",
		})
	}

	clientIDAny := facades.Config().Env("GOOGLE_CALENDAR_CLIENT_ID", "")
	clientSecretAny := facades.Config().Env("GOOGLE_CALENDAR_CLIENT_SECRET", "")
	redirectURIAny := facades.Config().Env("GOOGLE_CALENDAR_REDIRECT_URI", "")

	clientID, _ := clientIDAny.(string)
	clientSecret, _ := clientSecretAny.(string)
	redirectURI, _ := redirectURIAny.(string)

	facades.Log().Info("Loaded Google credentials from .env", map[string]any{
		"clientID": clientID,
		"redirectURI": redirectURI,
	})

	if clientID == "" || clientSecret == "" || redirectURI == "" {
		facades.Log().Error("Missing Google credentials in .env")
		return ctx.Response().Json(500, http.Json{
			"error": "Missing Google Calendar credentials in .env",
		})
	}

	config := &oauth2.Config{
		ClientID:     clientID,
		ClientSecret: clientSecret,
		RedirectURL:  redirectURI,
		Scopes:       []string{calendar.CalendarScope, "https://www.googleapis.com/auth/userinfo.email"},
		Endpoint:     google.Endpoint,
	}

	// Exchange code for token
	token, err := config.Exchange(ctx, code)
	if err != nil {
		facades.Log().Error("Failed to exchange code for token", map[string]any{"error": err.Error()})
		return ctx.Response().Json(500, http.Json{
			"error": fmt.Sprintf("Failed to exchange token: %v", err),
		})
	}
	facades.Log().Info("Token exchanged successfully", map[string]any{
		"access_token": token.AccessToken,
		"expiry":       token.Expiry,
	})

	// Get user email via Google API
	client := config.Client(context.Background(), token)
	userInfoResp, err := client.Get("https://www.googleapis.com/oauth2/v2/userinfo")
	if err != nil {
		facades.Log().Error("Failed to get user info", map[string]any{"error": err.Error()})
		return ctx.Response().Json(500, http.Json{
			"error": "Failed to get user info",
		})
	}
	defer userInfoResp.Body.Close()

	var userInfo map[string]any
	if err := json.NewDecoder(userInfoResp.Body).Decode(&userInfo); err != nil {
		facades.Log().Error("Failed to parse user info", map[string]any{"error": err.Error()})
		return ctx.Response().Json(500, http.Json{
			"error": "Failed to parse user info",
		})
	}
	facades.Log().Info("User info received from Google", userInfo)

	emailAny, ok := userInfo["email"]
	if !ok || emailAny == nil {
		facades.Log().Error("Email is missing in user info", userInfo)
		return ctx.Response().Json(500, http.Json{
			"error": "Google user info missing email",
		})
	}
	email := emailAny.(string)

	data := map[string]any{
		"email":         email,
		"access_token":  token.AccessToken,
		"refresh_token": token.RefreshToken,
		"expiry":        token.Expiry,
		"team":          "schedule",
	}

	// Upsert record
	var existing map[string]any
	facades.Orm().Query().Table("gmail_accounts").Where("team", "schedule").First(&existing)
	if existing != nil {
		facades.Orm().Query().Table("gmail_accounts").Where("team", "schedule").Update(data)
		
	} else {
		facades.Orm().Query().Table("gmail_accounts").Create(data)
		
	}

	return ctx.Response().Json(200, http.Json{
		"message": "Google account connected",
		"email":   email,
	})
}


