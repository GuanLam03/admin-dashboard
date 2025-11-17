package role

import (
	"fmt"
	"github.com/casbin/casbin/v2"
	"github.com/goravel/framework/contracts/http"
	"github.com/goravel/framework/facades"
	"github.com/spf13/cast"
	"goravel/app/messages"
	"goravel/app/models"
	"goravel/app/permissions"
	"strconv"
)

type RoleController struct{}

func NewRoleController() *RoleController {
	return &RoleController{}
}

// GET /roles
func (r *RoleController) Index(ctx http.Context) http.Response {
	var roles []models.Role
	if err := facades.Orm().Query().Find(&roles); err != nil {
		return ctx.Response().Json(500, http.Json{"error": messages.GetError("validation.internal_error")})
	}

	return ctx.Response().Json(200, http.Json{"message": roles})

}

// POST /roles
func (r *RoleController) Store(ctx http.Context) http.Response {
	name := ctx.Request().Input("name")

	role := models.Role{Name: name}
	if err := facades.Orm().Query().Create(&role); err != nil {
		return ctx.Response().Json(http.StatusInternalServerError, http.Json{"error": messages.GetError("validation.internal_error")})
	}
	return ctx.Response().Json(200, role)
}

// POST /users/:id/roles
func (r *RoleController) AssignToUser(ctx http.Context) http.Response {
	userID := ctx.Request().Route("id")
	roleID := ctx.Request().Input("role_id")

	// userRole := models.UserRole{}
	// create relation
	if err := facades.Orm().Query().Create(&models.UserRole{
		UserID: cast.ToUint(userID),
		RoleID: cast.ToUint(roleID),
	}); err != nil {
		return ctx.Response().Json(500, http.Json{
			"error": messages.GetError("validation.internal_error"),
		})
	}

	return ctx.Response().Json(201, http.Json{
		"message": messages.GetSuccess("role_assigned"),
	})
}

func (r *RoleController) Show(ctx http.Context) http.Response {
	idStr := ctx.Request().Route("id")
	roleID, err := strconv.Atoi(idStr)
	if err != nil {
		return ctx.Response().Json(422, map[string]interface{}{
			"error": messages.GetError("validation.invalid_request"),
		})
	}

	var role models.Role
	if err := facades.Orm().Query().Where("id = ?", roleID).First(&role); err != nil {
		return ctx.Response().Json(404, map[string]interface{}{
			"error": messages.GetError("validation.role_not_found"),
		})
	}

	// Get Casbin policies for this role
	enforcerAny, err := facades.App().Make("casbin")
	if err != nil {
		return ctx.Response().Json(500, map[string]interface{}{
			"error": messages.GetError("validation.casbin_not_initialized"),
		})
	}

	enforcer, ok := enforcerAny.(*casbin.Enforcer)
	if !ok {
		return ctx.Response().Json(500, map[string]interface{}{
			"error": messages.GetError("validation.casbin_cast_failed"),
		})
	}

	policies, _ := enforcer.GetFilteredPolicy(0, strconv.Itoa(int(roleID)))

	// Convert policies to permission keys (optional, map back to your PermissionMap)
	permissionsMap := []string{}
	for _, p := range policies {
		object := p[1] // object/path
		action := p[2] // method
		permKey := permissions.PermissionObjectActionToKey(object, action)
		if permKey != "" {
			permissionsMap = append(permissionsMap, permKey)
		}
	}

	return ctx.Response().Json(200, map[string]interface{}{
		"role":        role,
		"permissions": permissionsMap,
	})
}

// Update role permissions by ID
func (r *RoleController) UpdatePermissions(ctx http.Context) http.Response {
	obj := ctx.Request().Path()   // URL
	act := ctx.Request().Method() // HTTP method
	fmt.Println("Hello: ", obj, act)
	idStr := ctx.Request().Route("id")
	id, err := strconv.Atoi(idStr)
	if err != nil {
		return ctx.Response().Json(422, map[string]interface{}{
			"error": messages.GetError("validation.invalid_request"),
		})
	}

	var body struct {
		Name        string   `json:"name"`
		Permissions []string `json:"permissions"`
	}

	if err := ctx.Request().Bind(&body); err != nil {
		return ctx.Response().Json(422, map[string]interface{}{
			"error": messages.GetError("validation.invalid_request"),
		})
	}

	// Fetch role
	var role models.Role
	if err := facades.Orm().Query().Where("id", id).First(&role); err != nil {
		return ctx.Response().Json(404, map[string]interface{}{
			"error": messages.GetError("validation.role_not_found"),
		})
	}

	// Update role name
	role.Name = body.Name
	if _, err := facades.Orm().Query().Where("id", id).Update(&role); err != nil {
		return ctx.Response().Json(500, map[string]interface{}{
			"error": messages.GetError("validation.role_update_failed"),
		})
	}

	// --- Casbin: delete existing policies for this role ---
	enforcerAny, err := facades.App().Make("casbin")
	if err != nil {
		return ctx.Response().Json(500, map[string]interface{}{
			"error": messages.GetError("validation.casbin_not_initialized"),
		})
	}

	enforcer, ok := enforcerAny.(*casbin.Enforcer)
	if !ok {
		return ctx.Response().Json(500, map[string]interface{}{
			"error": messages.GetError("validation.casbin_cast_failed"),
		})
	}

	roleID := int(role.ID)
	// Remove all existing policies for this role
	_, _ = enforcer.DeletePermissionsForUser(strconv.Itoa(roleID))

	// Insert new policies
	for _, permKey := range body.Permissions {
		object, action := permissions.PermissionKeyToObjectAction(permKey)
		if object != "" && action != "" {
			_, _ = enforcer.AddPolicy(strconv.Itoa(roleID), object, action)
		}
	}

	return ctx.Response().Json(200, map[string]interface{}{
		"message": messages.GetSuccess("role_permissions_updated"),
		"role":    role,
	})
}
