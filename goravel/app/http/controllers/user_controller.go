package controllers

import (
	"github.com/goravel/framework/contracts/http"
	"goravel/app/models"
    "github.com/goravel/framework/facades"
    
)

type UserController struct {
	// Dependent services
}

func NewUserController() *UserController {
	return &UserController{
		// Inject services
	}
}

func (r *UserController) Show(ctx http.Context) http.Response {
	return ctx.Response().Success().Json(http.Json{
		"Hello": "Goravel",
	})
}

func (r *UserController) Edit(ctx http.Context) http.Response {
    name := ctx.Request().Input("name")
    currentPassword := ctx.Request().Input("currentPassword")
    newPassword := ctx.Request().Input("newPassword")
    confirmPassword := ctx.Request().Input("confirmPassword")

    validator, _ := ctx.Request().Validate(map[string]string{
        "name": "required|max_len:10",
    })

    if validator.Fails() {
        return ctx.Response().Json(422, http.Json{"error": validator.Errors().One()})
    }

    // Get logged-in user
    var user models.User
    if err := facades.Auth(ctx).User(&user); err != nil {
        return ctx.Response().Json(401, http.Json{"error": "Unauthenticated"})
    }

    // Update name
    user.Name = name

    // If new password is provided, validate current password first
    if newPassword != "" {
        // 1. Check current password
        if !facades.Hash().Check(currentPassword, user.Password) {
            return ctx.Response().Json(400, http.Json{"error": "Current password is incorrect"})
        }

        // 2. Match confirm
        if newPassword != confirmPassword {
            return ctx.Response().Json(400, http.Json{"error": "New passwords do not match"})
        }

        // 3. Hash and save new password
        hashed, err := facades.Hash().Make(newPassword)
        if err != nil {
            return ctx.Response().Json(500, http.Json{"error": "Failed to hash password"})
        }
        user.Password = hashed
    }

    // Save changes
    if err := facades.Orm().Query().Save(&user); err != nil {
        return ctx.Response().Json(500, http.Json{"error": err.Error()})
    }

    return ctx.Response().Json(200, http.Json{
        "message": "Profile updated successfully",
        "user": map[string]any{
            "id":    user.ID,
            "name":  user.Name,
            "email": user.Email,
        },
    })
}
