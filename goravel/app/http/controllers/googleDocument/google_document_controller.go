package googleDocument

import (
	"time"
	"strings"
	"github.com/goravel/framework/contracts/http"
	"github.com/goravel/framework/facades"

	"goravel/app/models"
	"goravel/app/messages"
)

type GoogleDocumentController struct{}

func NewGoogleDocumentController() *GoogleDocumentController {
	return &GoogleDocumentController{}
}

// GET /google-documents
// List documents with optional filters
func (c *GoogleDocumentController) ListGoogleDocuments(ctx http.Context) http.Response {
	// Apply filters if query exists, otherwise just get all
	documents, err := c.filter(ctx)
	if err != nil {
		return ctx.Response().Json(500, map[string]string{"error": messages.GetError("validation.internal_error"),})
	}

	return ctx.Response().Json(200, map[string]any{
		"documents": documents,
		"status":    models.GoogleDocumentStatusMap,
	})
}

// filter applies query filters, similar to Laravel
func (c *GoogleDocumentController) filter(ctx http.Context) ([]models.GoogleDocument, error) {
	query := facades.Orm().Query()

	// Filter: name
	if name := ctx.Request().Query("name"); name != "" {
		name = strings.Trim(name, "\"'")
		query = query.Where("name LIKE ?", "%"+name+"%")
	}

	// Filter: status
	if status := ctx.Request().Query("status"); status != "" {
		status = strings.Trim(status, "\"'")
		query = query.Where("status = ?", status)
	}

	// Filter: date range
	if fdate := ctx.Request().Query("fdate"); fdate != "" {
		if _, err := time.Parse("2006-01-02", fdate); err == nil {
			query = query.Where("created_at >= ?", fdate)
		}
	}
	if tdate := ctx.Request().Query("tdate"); tdate != "" {
		if _, err := time.Parse("2006-01-02", tdate); err == nil {
			query = query.Where("created_at <= ?", tdate+" 23:59:59")
		}
	}

	// Fetch from DB
	var documents []models.GoogleDocument
	if err := query.Get(&documents); err != nil {
		return nil, err
	}

	return documents, nil
}

// GET /google-documents/:id
// Show a single Google document
func (c *GoogleDocumentController) ShowGoogleDocument(ctx http.Context) http.Response {
	id := ctx.Request().Route("id")

	var doc models.GoogleDocument
	if err := facades.Orm().Query().Find(&doc, id); err != nil || doc.ID == 0 {
		return ctx.Response().Json(404, map[string]string{"error": messages.GetError("validation.google_document_not_found"),})
	}

	return ctx.Response().Json(200, map[string]any{
		"document": doc,
	})
}
