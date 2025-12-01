// app/http/controllers/document_controller.go
package docuements

import (
	"fmt"
	"path/filepath"
	"strconv"
	"goravel/app/websocket"
	"github.com/goravel/framework/contracts/http"
	"github.com/goravel/framework/facades"
	"goravel/app/messages"
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
		return ctx.Response().Json(422, http.Json{"error": messages.GetError("validation_failed")})
	}
	for _, file := range files {
		filename := file.GetClientOriginalName()
		if _, err := file.StoreAs("uploads", filename); err != nil {
			return ctx.Response().Json(500, http.Json{"error": messages.GetError("internal_error")})
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
				"error": messages.GetError("internal_error"),
			})
		}


		// websocket notification
		var user models.User
		facades.Auth(ctx).User(&user)
		websocket.PublishToCentrifugo("document_notifications", map[string]interface{}{
			"user_id": user.ID, 
			"message": fmt.Sprintf("New document added: %s, Please refresh the page", doc.Filename),
		})
	}

	


	return ctx.Response().Json(200, http.Json{"message": messages.GetSuccess("files_uploaded")})
}

// List
func (c *DocumentController) Index(ctx http.Context) http.Response {
	var docs []models.Document
	if err := facades.Orm().Query().Get(&docs); err != nil {
		return ctx.Response().Json(500, http.Json{"error": messages.GetError("internal_error")})
	}

	return ctx.Response().Json(200, http.Json{"documents": docs})
}

// Download file
func (c *DocumentController) Download(ctx http.Context) http.Response {
	idStr := ctx.Request().Route("id")
	id, err := strconv.Atoi(idStr)
	var doc models.Document

	if err != nil {
		return ctx.Response().Json(422, http.Json{"error": messages.GetError("validation_failed")})
	}

	if err := facades.Orm().Query().Where("id", id).First(&doc); err != nil {
		return ctx.Response().Json(404, http.Json{"error": messages.GetError("document_not_found")})
	}

	// Double-check file exists
	if !doc.Exists() {
		return ctx.Response().Json(404, http.Json{"error": messages.GetError("document_not_found")})
	}

	// Return file for download
	return ctx.Response().Download(doc.FullPath(), doc.Filename) // build full absolute path
}
