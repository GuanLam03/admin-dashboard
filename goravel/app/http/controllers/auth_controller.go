package controllers

import (
	// "time"

	"fmt"
	"goravel/app/messages"
	"goravel/app/models"
	"goravel/app/permissions"

	"github.com/casbin/casbin/v2"
	"github.com/goravel/framework/contracts/http"
	"github.com/goravel/framework/facades"
	"github.com/pquerna/otp/totp"
)

type AuthController struct{}

func NewAuthController() *AuthController {
	return &AuthController{}
}

// Register endpoint
func (a *AuthController) Register(ctx http.Context) http.Response {
	var user models.User
	if err := ctx.Request().Bind(&user); err != nil {
		return ctx.Response().Json(400, http.Json{"error": models.UserErrorMessage["invalid_request"]})
	}

	// Check if the email already exists in the database
	var existingUser models.User
	facades.Orm().Query().Where("email = ?", user.Email).First(&existingUser)

	if existingUser.ID != 0 {
		return ctx.Response().Json(400, http.Json{"error": models.UserErrorMessage["email_exists"]})
	}

	// Hash the password
	hashed, err := facades.Hash().Make(user.Password)
	if err != nil {
		return ctx.Response().Json(500, http.Json{"error": models.UserErrorMessage["internal_error"]})
	}
	user.Password = hashed

	if err := facades.Orm().Query().Create(&user); err != nil {
		return ctx.Response().Json(500, http.Json{"error": models.UserErrorMessage["create_failed"]})
	}

	return ctx.Response().Json(200, http.Json{"message": messages.GetSuccess("user_registered")})
}

// with google authenticator (new)
func (c *AuthController) Login(ctx http.Context) http.Response {
	email := ctx.Request().Input("email")
	password := ctx.Request().Input("password")

	var user models.User
	if err := facades.Orm().Query().Where("email", email).First(&user); err != nil {
		return ctx.Response().Json(401, http.Json{"error": models.UserErrorMessage["invalid_credentials"]})
	}

	if !facades.Hash().Check(password, user.Password) {
		return ctx.Response().Json(401, http.Json{"error": models.UserErrorMessage["invalid_credentials"]})
	}

	// If 2FA is enabled, don’t log in yet
	if user.TwoFactorEnabled {
		return ctx.Response().Json(200, http.Json{
			"message":        messages.GetSuccess("twofa_required"),
			"twofa_required": true,
			"user_id":        user.ID,
		})
	}

	// Normal login (no 2FA)
	token, err := facades.Auth(ctx).Login(&user)
	if err != nil {
		return ctx.Response().Json(500, http.Json{"error": models.UserErrorMessage["internal_error"]})
	}

	cookie := "jwt_token=" + token + "; Path=/; HttpOnly; Secure; SameSite=None; Max-Age=86400"
	ctx.Response().Header("Set-Cookie", cookie)

	return ctx.Response().Json(200, http.Json{"message": messages.GetSuccess("login_success"), "token": token})
}

func (a *AuthController) Logout(ctx http.Context) http.Response {

	expired := "jwt_token=; Path=/; HttpOnly; Secure; SameSite=None; Max-Age=-1"
	ctx.Response().Header("Set-Cookie", expired)
	return ctx.Response().Json(200, http.Json{"message": messages.GetSuccess("logged_out")})
}

// func (a *AuthController) Profile(ctx http.Context) http.Response {
//     var user models.User

//     // Get the auth guard
//     facades.Auth(ctx).User(&user);

//     response := map[string]any{
//         "name":       user.Name,
//         "email":      user.Email,
//         "two_factor_enabled": user.TwoFactorEnabled,
//         "created_at": user.CreatedAt,
//     }

//     return ctx.Response().Json(200, http.Json{
//         "user": response,
//     })
// }

func (a *AuthController) Profile(ctx http.Context) http.Response {
	var user models.User
	facades.Auth(ctx).User(&user)

	enforcerAny, err := facades.App().Make("casbin")
	if err != nil {
		return ctx.Response().Json(500, models.UserErrorMessage["internal_error"])

	}

	enforcer := enforcerAny.(*casbin.Enforcer)

	// Step 1: get roles for this user
	roles, _ := enforcer.GetRolesForUser(fmt.Sprint(user.ID)) // g(userId, roleId)

	// Step 2: collect permissions for each role
	permissionsList := []string{}
	for _, role := range roles {
		policies, _ := enforcer.GetFilteredPolicy(0, role) // p(role, obj, act)
		for _, policy := range policies {
			if len(policy) >= 3 {
				key := permissions.PermissionObjectActionToKey(policy[1], policy[2])
				if key != "" {
					permissionsList = append(permissionsList, key)
				}
			}
		}
	}

	response := map[string]any{
		"name":               user.Name,
		"email":              user.Email,
		"two_factor_enabled": user.TwoFactorEnabled,
		"created_at":         user.CreatedAt,
		"roles":              roles,
		"permissions":        permissionsList,
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
		return ctx.Response().Json(404, http.Json{"error": models.UserErrorMessage["not_found"]})
	}

	if !user.TwoFactorEnabled {
		return ctx.Response().Json(400, http.Json{"error": models.TwofaErrorMessage["not_enabled"]})
	}

	decryptedSecret, err := facades.Crypt().DecryptString(user.TwoFactorSecret)
	if err != nil {
		return ctx.Response().Json(500, http.Json{"error": models.TwofaErrorMessage["decrypt_failed"]})
	}

	// Validate TOTP code from Google Authenticator
	if totp.Validate(code, decryptedSecret) {
		token, err := facades.Auth(ctx).Login(&user)
		if err != nil {
			return ctx.Response().Json(500, http.Json{"error": models.TwofaErrorMessage["internal_error"]})
		}

		cookie := "jwt_token=" + token + "; Path=/; HttpOnly; Secure; SameSite=None; Max-Age=86400"
		ctx.Response().Header("Set-Cookie", cookie)

		return ctx.Response().Json(200, http.Json{"message": messages.GetSuccess("login_success"), "token": token})
	}

	return ctx.Response().Json(400, http.Json{"error": models.TwofaErrorMessage["invalid_code"]})
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
