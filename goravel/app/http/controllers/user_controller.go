package controllers

import (
	"goravel/app/messages"
	"goravel/app/models"
	"strconv"

	"github.com/goravel/framework/contracts/http"
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

// Index returns all users with their role
func (r *UserController) Index(ctx http.Context) http.Response {
	var users []models.User

	// Fetch all users
	if err := facades.Orm().Query().Find(&users); err != nil {
		return ctx.Response().Json(500, map[string]any{
			"error": messages.GetError("validation.internal_error"),
		})
	}

	// Fetch roles from casbin_rule (ptype=g means user→role assignment)
	type CasbinRule struct {
		V0 string // userID
		V1 string // roleID
	}
	var rules []CasbinRule
	if err := facades.Orm().Query().Table("casbin_rule").Where("ptype = ?", "g").Find(&rules); err != nil {
		return ctx.Response().Json(500, map[string]any{
			"error": messages.GetError("validation.internal_error"),
		})
	}

	// Fetch all roles to map roleID → roleName
	type Role struct {
		ID   uint
		Name string
	}
	var roles []Role
	if err := facades.Orm().Query().Table("roles").Find(&roles); err != nil {
		return ctx.Response().Json(500, map[string]any{
			"error": messages.GetError("validation.internal_error"),
		})
	}
	roleMap := make(map[string]string)
	for _, role := range roles {
		roleMap[strconv.Itoa(int(role.ID))] = role.Name
	}

	// Map userID → roleID from casbin_rule
	userRoles := make(map[string]string)
	for _, rule := range rules {
		userRoles[rule.V0] = rule.V1
	}

	// Build response with role name
	response := []map[string]any{}
	for _, u := range users {
		roleID := userRoles[strconv.Itoa(int(u.ID))]
		roleName := roleMap[roleID]

		response = append(response, map[string]any{
			"id":         u.ID,
			"name":       u.Name,
			"email":      u.Email,
			"created_at": u.CreatedAt,
			"role":       roleName,
		})
	}

	return ctx.Response().Json(200, map[string]any{
		"users": response,
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
		return ctx.Response().Json(401, http.Json{"error": messages.GetError("validation.unauthorized")})
	}

	// Update name
	user.Name = name

	// If new password is provided, validate current password first
	if newPassword != "" {
		// 1. Check current password
		if !facades.Hash().Check(currentPassword, user.Password) {
			return ctx.Response().Json(400, http.Json{"error": messages.GetError("validation.current_password_incorrect")})
		}

		// 2. Match confirm
		if newPassword != confirmPassword {
			return ctx.Response().Json(400, http.Json{"error": messages.GetError("validation.password_mismatch")})
		}

		// 3. Hash and save new password
		hashed, err := facades.Hash().Make(newPassword)
		if err != nil {
			return ctx.Response().Json(500, http.Json{"error": messages.GetError("validation.internal_error")})
		}
		user.Password = hashed
	}

	// Save changes
	if err := facades.Orm().Query().Save(&user); err != nil {
		return ctx.Response().Json(500, http.Json{"error": messages.GetError("validation.internal_error")})
	}

	return ctx.Response().Json(200, http.Json{
		"message": messages.GetSuccess("profile_updated"),
		"user": map[string]any{
			"id":    user.ID,
			"name":  user.Name,
			"email": user.Email,
		},
	})
}
