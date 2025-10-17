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
		return ctx.Response().Json(http.StatusInternalServerError, map[string]any{
			"error": err.Error(),
		})
	}

	return ctx.Response().Json(http.StatusOK, templates)
}


// POST /gmail/templates
func (c *TemplateController) AddTemplate(ctx http.Context) http.Response {
	
	var template models.GmailTemplate
	if err := ctx.Request().Bind(&template); err != nil {
		return ctx.Response().Json(http.StatusBadRequest, map[string]any{
			"error": "Invalid request body",
		})
	}

	if err := facades.Orm().Query().Create(&template); err != nil {
		return ctx.Response().Json(http.StatusInternalServerError, map[string]any{
			"error": err.Error(),
		})
	}
	return ctx.Response().Json(http.StatusOK, template)
}

// PUT /gmail/templates/{id}
func (c *TemplateController) EditTemplate(ctx http.Context) http.Response {
	id := ctx.Request().Route("id")

	var existingTemplate models.GmailTemplate
	if err := facades.Orm().Query().Where("id", id).First(&existingTemplate); err != nil {
		return ctx.Response().Json(http.StatusNotFound, map[string]any{
			"error": "Template not found",
		})
	}

	var template models.GmailTemplate
	if err := ctx.Request().Bind(&template); err != nil {
		return ctx.Response().Json(http.StatusBadRequest, map[string]any{
			"error": "Invalid request body",
		})
	}

	existingTemplate.Team = template.Team
	existingTemplate.Name = template.Name
	existingTemplate.Content = template.Content

	if err := facades.Orm().Query().Save(&existingTemplate); err != nil {
		return ctx.Response().Json(http.StatusInternalServerError, map[string]any{
			"error": err.Error(),
		})
	}
	return ctx.Response().Json(http.StatusOK, existingTemplate)
}

// DELETE /gmail/templates/{id}
func (c *TemplateController) RemoveTemplate(ctx http.Context) http.Response {
	id := ctx.Request().Route("id")

	if _,err := facades.Orm().Query().Where("id", id).Delete(&models.GmailTemplate{}); err != nil {
		return ctx.Response().Json(http.StatusInternalServerError, map[string]any{
			"error": err.Error(),
		})
	}
	return ctx.Response().Json(http.StatusOK, map[string]string{
		"message": "Template deleted successfully",
	})
}
