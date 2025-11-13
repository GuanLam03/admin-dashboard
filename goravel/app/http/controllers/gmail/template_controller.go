package gmail

import (
	"goravel/app/models"
	"github.com/goravel/framework/contracts/http"
	"github.com/goravel/framework/facades"
)

type TemplateController struct{}

func NewTemplateController() *TemplateController {
	return &TemplateController{}
}



// GET /gmail/templates
func (c *TemplateController) ShowTemplates(ctx http.Context) http.Response {
	team := ctx.Request().Query("team") // optional filter

	var query = facades.Orm().Query()
	if team != "" {
		query = query.Where("team = ?", team)
	}

	var templates []models.GmailTemplate
	if err := query.Find(&templates); err != nil {
		facades.Log().Errorf("Failed to fetch Gmail templates: %v", err)
		return ctx.Response().Json(http.StatusInternalServerError, map[string]any{
			"error": facades.Lang(ctx).Get("validation.internal_error"),
		})
	}

	return ctx.Response().Json(http.StatusOK, templates)
}


// POST /gmail/templates
func (c *TemplateController) AddTemplate(ctx http.Context) http.Response {
	
	var template models.GmailTemplate
	if err := ctx.Request().Bind(&template); err != nil {
		facades.Log().Warningf("Invalid template input: %v", err)
		return ctx.Response().Json(http.StatusBadRequest, map[string]any{
			"error": facades.Lang(ctx).Get("validation.invalid_request"),
		})
	}

	if err := facades.Orm().Query().Create(&template); err != nil {
		facades.Log().Errorf("Failed to create Gmail template: %v", err)
		return ctx.Response().Json(http.StatusInternalServerError, map[string]any{
			"error": facades.Lang(ctx).Get("validation.gmail_template_create_failed"),
		})
	}
	return ctx.Response().Json(http.StatusOK, template)
}

// PUT /gmail/templates/{id}
func (c *TemplateController) EditTemplate(ctx http.Context) http.Response {
	id := ctx.Request().Route("id")

	var existingTemplate models.GmailTemplate
	if err := facades.Orm().Query().Where("id", id).First(&existingTemplate); err != nil {
		facades.Log().Warningf("Gmail template not found (id=%s): %v", id, err)
		return ctx.Response().Json(http.StatusNotFound, map[string]any{
			"error": facades.Lang(ctx).Get("validation.gmail_template_not_found"),
		})
	}

	var template models.GmailTemplate
	if err := ctx.Request().Bind(&template); err != nil {
		facades.Log().Warningf("Invalid update body: %v", err)
		return ctx.Response().Json(http.StatusBadRequest, map[string]any{
			"error": facades.Lang(ctx).Get("validation.invalid_request"),
		})
	}

	existingTemplate.Team = template.Team
	existingTemplate.Name = template.Name
	existingTemplate.Content = template.Content

	if err := facades.Orm().Query().Save(&existingTemplate); err != nil {
		facades.Log().Errorf("Failed to update Gmail template (id=%s): %v", id, err)
		return ctx.Response().Json(http.StatusInternalServerError, map[string]any{
			"error": facades.Lang(ctx).Get("validation.internal_error"),
		})
	}
	return ctx.Response().Json(http.StatusOK, existingTemplate)
}

// DELETE /gmail/templates/{id}
func (c *TemplateController) RemoveTemplate(ctx http.Context) http.Response {
	id := ctx.Request().Route("id")

	if _,err := facades.Orm().Query().Where("id", id).Delete(&models.GmailTemplate{}); err != nil {
		facades.Log().Errorf("Failed to delete Gmail template (id=%s): %v", id, err)
		return ctx.Response().Json(http.StatusInternalServerError, map[string]any{
			"error":facades.Lang(ctx).Get("validation.internal_error"),
		})
	}
	return ctx.Response().Json(http.StatusOK, map[string]string{
		"message": "Template deleted successfully",
	})
}
