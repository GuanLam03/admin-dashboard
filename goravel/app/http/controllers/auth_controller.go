package controllers

import (
	// "fmt"
    // "time"

    "github.com/pquerna/otp/totp"
    "goravel/app/models"
    "github.com/goravel/framework/facades"
    "github.com/goravel/framework/contracts/http"
)

type AuthController struct{}

func NewAuthController() *AuthController {
    return &AuthController{}
}

// Register endpoint
func (a *AuthController) Register(ctx http.Context) http.Response {
    var user models.User
    if err := ctx.Request().Bind(&user); err != nil {
        return ctx.Response().Json(400, http.Json{"error": err.Error()})
    }


    // Check if the email already exists in the database
    var existingUser models.User
    facades.Orm().Query().Where("email = ?", user.Email).First(&existingUser)
    
    if existingUser.ID != 0 {
     return ctx.Response().Json(400, http.Json{"error": "email already taken"})
    }
    

    // Hash the password
    hashed, err := facades.Hash().Make(user.Password)
    if err != nil {
        return ctx.Response().Json(500, http.Json{"error": "failed to hash password"})
    }
    user.Password = hashed

    if err := facades.Orm().Query().Create(&user); err != nil {
        return ctx.Response().Json(500, http.Json{"error": err.Error()})
    }

    return ctx.Response().Json(200, http.Json{"message": "user registered successfully"})
}



// with google authenticator (new)
func (c *AuthController) Login(ctx http.Context) http.Response {
    email := ctx.Request().Input("email")
    password := ctx.Request().Input("password")

    var user models.User
    if err := facades.Orm().Query().Where("email", email).First(&user); err != nil {
        return ctx.Response().Json(401, http.Json{"error": "invalid credentials"})
    }

    if !facades.Hash().Check(password, user.Password) {
        return ctx.Response().Json(401, http.Json{"error": "invalid credentials"})
    }

    // If 2FA is enabled, don’t log in yet
    if user.TwoFactorEnabled {
        return ctx.Response().Json(200, http.Json{
            "message":      "2FA required",
            "twofa_required": true,
            "user_id":      user.ID,
        })
    }

    // Normal login (no 2FA)
    token, err := facades.Auth(ctx).Login(&user)
    if err != nil {
        return ctx.Response().Json(500, http.Json{"error": err.Error()})
    }

    cookie := "jwt_token=" + token + "; Path=/; HttpOnly; Secure; SameSite=None; Max-Age=86400"
    ctx.Response().Header("Set-Cookie", cookie)

    return ctx.Response().Json(200, http.Json{"message": "login success", "token": token})
}




func (a *AuthController) Logout(ctx http.Context) http.Response {
    
    expired := "jwt_token=; Path=/; HttpOnly; Secure; SameSite=None; Max-Age=-1"
    ctx.Response().Header("Set-Cookie", expired)
    return ctx.Response().Json(200, http.Json{"message": "logged out"})
}


func (a *AuthController) Profile(ctx http.Context) http.Response {
    var user models.User

    // Get the auth guard
    facades.Auth(ctx).User(&user);
	
    response := map[string]any{
        "name":       user.Name,
        "email":      user.Email,
        "two_factor_enabled": user.TwoFactorEnabled,
        "created_at": user.CreatedAt,
    }

    return ctx.Response().Json(200, http.Json{
        "user": response,
    })
}




func (c *AuthController) VerifyTwoFA(ctx http.Context) http.Response {
    userId := ctx.Request().Input("user_id")
    code := ctx.Request().Input("code")

    var user models.User
    if err := facades.Orm().Query().Find(&user, userId); err != nil {
        return ctx.Response().Json(404, http.Json{"error": "user not found"})
    }

    if !user.TwoFactorEnabled {
        return ctx.Response().Json(400, http.Json{"error": "2FA not enabled"})
    }

    decryptedSecret, err := facades.Crypt().DecryptString(user.TwoFactorSecret)
    if err != nil {
        return ctx.Response().Json(500, http.Json{"error": "failed to decrypt 2FA secret"})
    }

    // Validate TOTP code from Google Authenticator
    if totp.Validate(code, decryptedSecret) {
        token, err := facades.Auth(ctx).Login(&user)
        if err != nil {
            return ctx.Response().Json(500, http.Json{"error": err.Error()})
        }

        cookie := "jwt_token=" + token + "; Path=/; HttpOnly; Secure; SameSite=None; Max-Age=86400"
        ctx.Response().Header("Set-Cookie", cookie)

        return ctx.Response().Json(200, http.Json{"message": "login success", "token": token})
    }

    return ctx.Response().Json(400, http.Json{"error": "invalid 2FA code"})
}




// func (a *AuthController) Login(ctx http.Context) http.Response {
//     email := ctx.Request().Input("email")
//     password := ctx.Request().Input("password")

//     var user models.User
//     if err := facades.Orm().Query().Where("email", email).First(&user); err != nil {
//         return ctx.Response().Json(401, http.Json{"error": "invalid credentials"})
//     }

//     // Check password
//     if !facades.Hash().Check(password, user.Password) {
//         return ctx.Response().Json(401, http.Json{"error": "invalid credentials"})
//     }

//     // Generate JWT
//     token, err := facades.Auth(ctx).Login(&user)
//     if err != nil {
//         return ctx.Response().Json(500, http.Json{"error": err.Error()})
//     }

//     cookie := "jwt_token=" + token + "; Path=/; HttpOnly; Secure; SameSite=None; Max-Age=86400"
//     ctx.Response().Header("Set-Cookie", cookie)  // using header as goravel cookie does not support sameSite (https://www.goravel.dev/the-basics/response.html#attach-header)

//     return ctx.Response().Json(200, http.Json{"message": "login success"})
// }
