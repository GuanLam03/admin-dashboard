package googleDocument

import (
	"github.com/goravel/framework/contracts/http"
	"github.com/goravel/framework/facades"

	"goravel/app/models"
)

type EditGoogleDocumentController struct{}

func NewEditGoogleDocumentController() *EditGoogleDocumentController {
	return &EditGoogleDocumentController{}
}

// GET /google-documents/:id
func (c *EditGoogleDocumentController) ShowGoogleDocument(ctx http.Context) http.Response {
	id := ctx.Request().Route("id")

	// Find document by ID
	var doc models.GoogleDocument
	if err := facades.Orm().Query().Find(&doc, id); err != nil || doc.ID == 0 {
		return ctx.Response().Json(404, map[string]string{"error":facades.Lang(ctx).Get("validation.google_document_not_found")})
	}

	// Remove "removed" status from allowed list (like Laravel did)
	status := models.GoogleDocumentStatusMap
	delete(status, "removed")

	return ctx.Response().Json(200, map[string]any{
		"document": doc,
		"status":   status,
	})
}


// PUT /google-documents/:id
func (c *EditGoogleDocumentController) EditGoogleDocument(ctx http.Context) http.Response {
	id := ctx.Request().Route("id")

	// Find document by ID
	var doc models.GoogleDocument
	if err := facades.Orm().Query().Find(&doc, id); err != nil || doc.ID == 0 {
		return ctx.Response().Json(404, map[string]string{"error": facades.Lang(ctx).Get("validation.google_document_not_found")})
	}

	// Collect request data
	data := map[string]interface{}{
		"name":          ctx.Request().Input("name"),
		"original_link": ctx.Request().Input("original_link"),
		"status":        ctx.Request().Input("status"),
	}

	// Validate
	status, errResp, err := validateGoogleDocumentInput(ctx,data)
	if err != nil {
		return ctx.Response().Json(500, map[string]string{"error": facades.Lang(ctx).Get("validation.internal_error")})
	}
	if errResp != nil {
		return ctx.Response().Json(422, errResp)
	}

	// Simplify Google link if provided
	link := doc.Link
	if data["original_link"] != nil && data["original_link"].(string) != "" {
		newLink, err := simplifyGoogleLink(ctx, data["original_link"].(string))
		if err != nil {
			return ctx.Response().Json(422, map[string]any{
				"error": err.Error(),
			})
		}
		link = newLink
	}

	// Update document
	doc.Name = data["name"].(string)
	doc.OriginalLink = data["original_link"].(string)
	doc.Link = link
	doc.Status = status

	if err := facades.Orm().Query().Save(&doc); err != nil {
		return ctx.Response().Json(500, map[string]string{"error": facades.Lang(ctx).Get("validation.google_document_update_failed")})
	}

	return ctx.Response().Json(200, map[string]any{
		"message": "Google document updated successfully",
		"data":    doc,
	})
}

//validateGoogleDocumentInput and simplifyGoogleLink are in add_google_document_controller.go !!!


func (c *EditGoogleDocumentController) RemoveGoogleDocument(ctx http.Context) http.Response {
	id := ctx.Request().Route("id")

	// Find document by ID
	var doc models.GoogleDocument
	if err := facades.Orm().Query().Find(&doc, id); err != nil || doc.ID == 0 {
		return ctx.Response().Json(404, map[string]string{"error": facades.Lang(ctx).Get("validation.google_document_not_found")})
	}

	// Mark as removed instead of deleting
	doc.Status = models.GoogleDocumentStatusMap["removed"]

	if err := facades.Orm().Query().Save(&doc); err != nil {
		return ctx.Response().Json(500, map[string]string{"error": facades.Lang(ctx).Get("validation.google_document_delete_failed")})
	}

	return ctx.Response().Json(200, map[string]any{
		"message": "Google document removed successfully",
		"data":    doc,
	})
}
