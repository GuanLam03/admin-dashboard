package googleDocument

import (
	"errors"
	"fmt"
	"regexp"
	// "strings"

	"github.com/goravel/framework/contracts/http"
	"github.com/goravel/framework/facades"
	

	"goravel/app/models"
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

	status, errResp, err := validateGoogleDocumentInput(data)
	if err != nil {
		return ctx.Response().Json(500, map[string]string{"error": err.Error()})
	}
	if errResp != nil {
		return ctx.Response().Json(422, errResp)
	}

	link, err := simplifyGoogleLink(data["original_link"].(string))
	if err != nil {
		return ctx.Response().Json(422, map[string]any{
			"error": err.Error(),
		})
	}

	doc := models.GoogleDocument{
		Name:         data["name"].(string),
		OriginalLink: data["original_link"].(string),
		Link:         link,
		Status:       status,
	}

	if err := facades.Orm().Query().Create(&doc); err != nil {
		return ctx.Response().Json(500, map[string]string{"error": err.Error()})
	}

	return ctx.Response().Json(200, doc)
}


func validateGoogleDocumentInput(data map[string]interface{}) (string, map[string]interface{}, error) {
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
			"error": "invalid status",
		}, nil
	}

	return status, nil, nil
}




func simplifyGoogleLink(original string) (string, error) {
    re := regexp.MustCompile(`https:\/\/(docs|drive)\.google\.com\/([a-zA-Z]+)\/d\/([a-zA-Z0-9_-]+)`)
    match := re.FindStringSubmatch(original)
    if len(match) != 4 {
        return "", errors.New("Google file link format is invalid")
    }

    domain, typ, id := match[1], match[2], match[3]
    previewLink := fmt.Sprintf("https://%s.google.com/%s/d/%s/preview", domain, typ, id)
	fmt.Println(previewLink);
    resp, err := facades.Http().
        WithHeader("Accept", "text/html").
        Get(previewLink)
    if err != nil || resp.Status() != 200 {
        return "", errors.New("Google file link not accessible or not public")
    }

    return previewLink, nil
}



// func (r *AddGoogleDocumentController) AddGoogleDocument(ctx http.Context) http.Response {
// 	//ctx.Request() unable to call twice times
// 	// Collect request data
// 	data := map[string]interface{}{
// 		"name":          ctx.Request().Input("name"),
// 		"original_link": ctx.Request().Input("original_link"),
// 		// "link":          ctx.Request().Input("link"),
// 		"status":        ctx.Request().Input("status"),
// 	}

// 	// 1. Validate data
// 	validator, err := facades.Validation().Make(data, models.GoogleDocumentRules)
// 	if err != nil {
// 		return ctx.Response().Json(http.StatusInternalServerError, map[string]string{"error": err.Error()})
// 	}
// 	if validator.Fails() {
// 		return ctx.Response().Json(http.StatusBadRequest, map[string]interface{}{
// 			"errors": validator.Errors().All(),
// 		})
// 	}

// 	// 2. Validate status against whitelist
// 	status, ok := data["status"].(string)
// 	if !ok || models.GoogleDocumentStatusMap[status] == "" {
// 		return ctx.Response().Json(http.StatusBadRequest, map[string]string{"error": "invalid status"})
// 	}

// 	link, err := simplifyGoogleLink(data["original_link"].(string))
// 	if err != nil {
// 		return ctx.Response().Json(http.StatusBadRequest, map[string]any{
// 			"error": err.Error(),
// 		})
// 	}

// 	// 3. Build model
// 	doc := models.GoogleDocument{
// 		Name:         data["name"].(string),
// 		OriginalLink: data["original_link"].(string),
// 		Link:         link,
// 		Status:       status,
// 	}

// 	// 4. Save to DB
// 	if err := facades.Orm().Query().Create(&doc); err != nil {
// 		return ctx.Response().Json(http.StatusInternalServerError, map[string]string{"error": err.Error()})
// 	}

// 	// 5. Return created resource
// 	return ctx.Response().Json(http.StatusCreated, doc)
// }
