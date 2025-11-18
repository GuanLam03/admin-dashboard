package googleDocument

import (
	"errors"
	"fmt"
	"regexp"
	// "strings"

	"github.com/goravel/framework/contracts/http"
	"github.com/goravel/framework/facades"
	

	"goravel/app/models"
	"goravel/app/messages"
)

type AddGoogleDocumentController struct {
}

func NewAddGoogleDocumentController() *AddGoogleDocumentController {
	return &AddGoogleDocumentController{}
}

// POST /admin/excel
func (r *AddGoogleDocumentController) AddGoogleDocument(ctx http.Context) http.Response {
	data := map[string]interface{}{
		"name":          ctx.Request().Input("name"),
		"original_link": ctx.Request().Input("original_link"),
		"status":        ctx.Request().Input("status"),
	}

	status, errResp, err := validateGoogleDocumentInput(ctx,data)
	if err != nil {
		return ctx.Response().Json(500, map[string]string{"error": messages.GetError("internal_error")})
	}
	if errResp != nil {
		return ctx.Response().Json(422, errResp)
	}

	link, err := simplifyGoogleLink(ctx, data["original_link"].(string))
	if err != nil {
		return ctx.Response().Json(422, map[string]any{
			"error": messages.GetError("invalid_link_format"),
		})
	}

	doc := models.GoogleDocument{
		Name:         data["name"].(string),
		OriginalLink: data["original_link"].(string),
		Link:         link,
		Status:       status,
	}

	if err := facades.Orm().Query().Create(&doc); err != nil {
		return ctx.Response().Json(500, map[string]string{"error": messages.GetError("google_document_create_failed")})
	}

	return ctx.Response().Json(200, doc)
}


func validateGoogleDocumentInput(ctx http.Context,data map[string]interface{}) (string, map[string]interface{}, error) {
	// Run validation rules
	validator, err := facades.Validation().Make(data, models.GoogleDocumentRules)
	if err != nil {
		return "", nil, fmt.Errorf("validation error: %v", err)
	}
	if validator.Fails() {
		return "", map[string]interface{}{
			"errors": validator.Errors().All(),
		}, nil
	}

	// Check status against whitelist
	status, ok := data["status"].(string)
	if !ok || models.GoogleDocumentStatusMap[status] == "" {
		return "", map[string]interface{}{
			"error": messages.GetError("invalid_status"),
		}, nil
	}

	return status, nil, nil
}




func simplifyGoogleLink(ctx http.Context,original string) (string, error) {
    re := regexp.MustCompile(`https:\/\/(docs|drive)\.google\.com\/([a-zA-Z]+)\/d\/([a-zA-Z0-9_-]+)`)
    match := re.FindStringSubmatch(original)
    if len(match) != 4 {
        return "", errors.New(messages.GetError("invalid_link_format"))
    }

    domain, typ, id := match[1], match[2], match[3]
    previewLink := fmt.Sprintf("https://%s.google.com/%s/d/%s/preview", domain, typ, id)
	// fmt.Println(previewLink);
    resp, err := facades.Http().
        WithHeader("Accept", "text/html").
        Get(previewLink)
    if err != nil || resp.Status() != 200 {
        return "", errors.New(messages.GetError("link_not_accessible"))
    }

    return previewLink, nil
}


