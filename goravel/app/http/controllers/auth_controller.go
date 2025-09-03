package controllers

import (
	// "fmt"
    // "time"
  
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


func (a *AuthController) Login(ctx http.Context) http.Response {
    email := ctx.Request().Input("email")
    password := ctx.Request().Input("password")

    var user models.User
    if err := facades.Orm().Query().Where("email", email).First(&user); err != nil {
        return ctx.Response().Json(401, http.Json{"error": "invalid credentials"})
    }

    // Check password
    if !facades.Hash().Check(password, user.Password) {
        return ctx.Response().Json(401, http.Json{"error": "invalid credentials"})
    }

    // Generate JWT
    token, err := facades.Auth(ctx).Login(&user)
    if err != nil {
        return ctx.Response().Json(500, http.Json{"error": err.Error()})
    }

    cookie := "jwt_token=" + token + "; Path=/; HttpOnly; Secure; SameSite=None; Max-Age=86400"
    ctx.Response().Header("Set-Cookie", cookie)  // using header as goravel cookie does not support sameSite (https://www.goravel.dev/the-basics/response.html#attach-header)

    return ctx.Response().Json(200, http.Json{"message": "login success"})
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
        "created_at": user.CreatedAt,
    }

    return ctx.Response().Json(200, http.Json{
        "user": response,
    })
}



