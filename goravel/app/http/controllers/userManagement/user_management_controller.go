package userManagement

import (
	"github.com/goravel/framework/contracts/http"
	"goravel/app/models"
    "github.com/goravel/framework/facades"
    "github.com/casbin/casbin/v2"
)

type UserManagementController struct {
	// Dependent services
}

func NewUserManagementController() *UserManagementController {
	return &UserManagementController{
		// Inject services
	}
}

type CasbinRule struct {
    Ptype string `gorm:"column:ptype"`
    V0    string `gorm:"column:v0"`
    V1    string `gorm:"column:v1"`
}

func (r *UserManagementController) Index(ctx http.Context) http.Response {
	return nil
}	


func (r *UserManagementController) ShowUserRole(ctx http.Context) http.Response {
	var roles []models.Role
	if err := facades.Orm().Query().Find(&roles); err != nil {
		return ctx.Response().Json(http.StatusInternalServerError, http.Json{"error":  err.Error()})
	}
	return ctx.Response().Json(200,http.Json{"message": roles})

}


func (r *UserManagementController) AssignRole(ctx http.Context) http.Response {
    id := ctx.Request().Route("id")
    roleID := ctx.Request().Input("role_id")

    if id == "" || roleID == "" {
        return ctx.Response().Json(422, map[string]interface{}{
            "error": "User ID and Role ID are required",
        })
    }

    // Get the Casbin enforcer instance (assuming it's bound in the container)
    enforcerAny, err := facades.App().Make("casbin")
		if err != nil {
			return ctx.Response().Json(500, "Failed to get Casbin enforcer")
			
		}

	e := enforcerAny.(*casbin.Enforcer)

    // Delete all existing roles for the user
    if _, err := e.DeleteRolesForUser(id); err != nil {
        // log.Printf("Error deleting roles for user %s: %v", id, err)
        return ctx.Response().Json(500, map[string]interface{}{
            "error": "Failed to remove old roles",
        })
    }

    // Assign the new role
    if _, err := e.AddRoleForUser(id, roleID); err != nil {
        return ctx.Response().Json(500, map[string]interface{}{
            "error": "Failed to assign new role",
        })
    }

    return ctx.Response().Json(200, map[string]interface{}{
        "message": "Role assigned successfully",
        "user_id": id,
        "role_id": roleID,
    })
}


