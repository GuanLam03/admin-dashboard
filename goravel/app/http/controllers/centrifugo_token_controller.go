package controllers

import (
	"time"
	"strconv"
	"github.com/golang-jwt/jwt/v5"
	"github.com/goravel/framework/contracts/http"
	"github.com/goravel/framework/facades"
	"goravel/app/messages"
	"goravel/app/models"

)

var centrifugoSecret = "6NYpez5jZ0cdNPs54c9BIQDZMqVzyzn5xp6p70aPm8WvXz6safUs7WPGhD4VJesQj216V42p1lbX3BGnDcg2fg" // SAME AS config.json

type CentrifugoTokenController struct{}


func NewCentrifugoTokenController() *CentrifugoTokenController {
	return &CentrifugoTokenController{}
}

func (c *CentrifugoTokenController) Generate(ctx http.Context) http.Response {
	var user models.User
	if err := facades.Auth(ctx).User(&user); err != nil {
		return ctx.Response().Json(401, http.Json{"error": messages.GetError("unauthorized")})
	}
	// Convert UID to string
    userID := strconv.Itoa(int(user.ID))

    // Load HMAC secret key from config
    hmacSecret := facades.Config().Env("hmac_secret_key").(string)

    if hmacSecret == "" {
        return ctx.Response().Json(500, http.Json{
            "error": "missing centrifugo hmac secret",
        })
    }

    // Create JWT
    claims := jwt.MapClaims{
		"sub": userID,
		"exp": time.Now().Add(1 * time.Hour).Unix(),
		// "info": map[string]interface{}{
		// 	"presence": true, 
		// },
	}


    tokenObj := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
    tokenString, err := tokenObj.SignedString([]byte(hmacSecret))

    if err != nil {
        return ctx.Response().Json(500, http.Json{
            "error": "token generation failed",
        })
    }

    return ctx.Response().Json(200, http.Json{
        "token": tokenString,
    })
}
