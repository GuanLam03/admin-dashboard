// app/http/controllers/document_controller.go
package docuements

import (
	// "fmt"
	"strconv"
	"path/filepath"

	"github.com/goravel/framework/contracts/http"
	"github.com/goravel/framework/facades"
	"goravel/app/models"

)

type DocumentController struct{}

func NewDocumentController() *DocumentController {
	return &DocumentController{}
}

// Upload
func (c *DocumentController) Store(ctx http.Context) http.Response {
	files, err := ctx.Request().Files("files")
	var uploadedPaths []string
	// fmt.Println(files);

	if err != nil || len(files) == 0 {
		return ctx.Response().Json(422, http.Json{"error": facades.Lang(ctx).Get("validation.validation_failed")})
	}

	
	for _, file := range files {
		filename := file.GetClientOriginalName()
		if _, err := file.StoreAs("uploads", filename); err != nil {
			return ctx.Response().Json(500, http.Json{"error": facades.Lang(ctx).Get("validation.internal_error")})
		}

		savePath := filepath.Join("uploads", filename)

		// Track uploaded file path in case we need to delete it later
		uploadedPaths = append(uploadedPaths, savePath)
		

		doc := models.Document{Filename: filename, Path: savePath}
		if err := facades.Orm().Query().Create(&doc); err != nil {
			// DB insert failed – clean up uploaded files
			for _, path := range uploadedPaths {
				facades.Storage().Delete(path)
			}

			return ctx.Response().Json(500, http.Json{
				"error": facades.Lang(ctx).Get("validation.internal_error"),
			})
		}
	}

	return ctx.Response().Json(200, http.Json{"message": "Files uploaded"})
}


// List
func (c *DocumentController) Index(ctx http.Context) http.Response {
	var docs []models.Document
	if err := facades.Orm().Query().Get(&docs); err != nil {
		return ctx.Response().Json(500, http.Json{"error": facades.Lang(ctx).Get("validation.internal_error")})
	}

	return ctx.Response().Json(200, http.Json{"documents": docs})
}

// Download file
func (c *DocumentController) Download(ctx http.Context) http.Response {
	idStr := ctx.Request().Route("id")
	id, err := strconv.Atoi(idStr)
	var doc models.Document

	if err != nil {
		return ctx.Response().Json(422, http.Json{"error": facades.Lang(ctx).Get("validation.validation_failed")})
	}

	if err := facades.Orm().Query().Where("id", id).First(&doc); err != nil {
		return ctx.Response().Json(404, http.Json{"error":facades.Lang(ctx).Get("validation.document_not_found"),})
	}

	// Double-check file exists
	if !doc.Exists() {
		return ctx.Response().Json(404, http.Json{"error": facades.Lang(ctx).Get("validation.document_not_found"),})
	}

	// Return file for download
	return ctx.Response().Download(doc.FullPath() , doc.Filename) // build full absolute path
}