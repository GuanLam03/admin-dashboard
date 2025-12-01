package adsCampaign

import (
	"fmt"
	"crypto/rand"
    "math/big"
	"github.com/goravel/framework/contracts/http"
	"github.com/goravel/framework/facades"
	// "encoding/json"
	"goravel/app/models"
	"github.com/goravel/framework/contracts/database/orm"
	"goravel/app/messages"
)

type AddAdsCampaignController struct {
}

func NewAddAdsCampaignController() *AddAdsCampaignController {
	return &AddAdsCampaignController{}
}

func buildURL(base string, port string, path string) string {
	// If production → no port
	if facades.Config().Env("APP_ENV") == "production" {
		return fmt.Sprintf("%s/%s", base, path)
	}

	// Development → keep port
	if port != "" {
		return fmt.Sprintf("%s:%s/%s", base, port, path)
	}

	// No port defined → fallback
	return fmt.Sprintf("%s/%s", base, path)
}


func (a *AddAdsCampaignController) AddAdsCampaign(ctx http.Context) http.Response {
	basedUrl := facades.Config().Env("APP_URL", "").(string)
	port := facades.Config().Env("APP_PORT", "").(string)

	type PostbackEvent struct {
		EventName string `json:"event_name"`
		PostbackUrl       string `json:"url"`
		IncludeClickParams bool   `json:"include_click_params"` 
	}

	type AddAdsCampaignRequest struct {
		Name            string          `json:"name"`
		PostbackEnabled bool            `json:"postback_enabled"`
		PostbackEvents  []PostbackEvent `json:"postback_events"`
		TargetUrl       string          `json:"target_url"`
	}

	// Bind request JSON
	var req AddAdsCampaignRequest
	if err := ctx.Request().Bind(&req); err != nil {
		return ctx.Response().Json(http.StatusBadRequest, map[string]any{
			"error": messages.GetError("invalid_request"),
		})
	}

	// Map request to AdsCampaign model
	var adsCampaign models.AdsCampaign
	adsCampaign.Name = req.Name
	adsCampaign.TargetUrl = req.TargetUrl
	adsCampaign.Status = models.AdsCampaignStatusMap["active"]

	// Validate campaign
	errResp, err := validateAdsCampaignInput(adsCampaign)
	if err != nil {
		return ctx.Response().Json(http.StatusInternalServerError, map[string]any{"error": messages.GetError("internal_error")})
	}
	if errResp != nil {
		return ctx.Response().Json(http.StatusUnprocessableEntity, errResp)
	}

	// Convert PostbackEvents to model type for validation
	var postbackModels []models.AdsCampaignPostback
	for _, pb := range req.PostbackEvents {
		postbackModels = append(postbackModels, models.AdsCampaignPostback{
			EventName:          pb.EventName,
			PostbackUrl:        pb.PostbackUrl,
			IncludeClickParams: pb.IncludeClickParams,
		})
	}

	// Validate postback events
	if len(postbackModels) > 0 {
		errResp, err = ValidateAdsCampaignPostbackInput(postbackModels)
		if err != nil {
			return ctx.Response().Json(http.StatusInternalServerError, map[string]string{"error": messages.GetError("internal_error")})
		}
		if errResp != nil {
			return ctx.Response().Json(http.StatusUnprocessableEntity, errResp)
		}
	}



	// Generate unique code
	code, err := generateUniqueCode()
	if err != nil {
		return ctx.Response().Json(http.StatusInternalServerError, map[string]any{
			"error": messages.GetError("internal_error"),
		})
	}
	adsCampaign.Code = code

	// Set tracking and postback links
	// Build URLs based on environment
	trackingLink := buildURL(basedUrl, port, adsCampaign.Code)
	postbackLink := buildURL(basedUrl, port, "postback/")

		
	adsCampaign.TrackingLink = &trackingLink
	adsCampaign.PostbackLink = &postbackLink

	// Transaction: create campaign and optional postbacks
	err = facades.Orm().Transaction(func(tx orm.Query) error {
		if err := tx.Create(&adsCampaign); err != nil {
			return err
		}

		if req.PostbackEnabled {
			for _, e := range req.PostbackEvents {
				postback := models.AdsCampaignPostback{
					AdsCampaignId: adsCampaign.ID,
					EventName:     e.EventName,
					PostbackUrl:   e.PostbackUrl,
					IncludeClickParams: e.IncludeClickParams,
				}
				if err := tx.Create(&postback); err != nil {
					return err
				}
			}
		}

		return nil
	})
	if err != nil {
		return ctx.Response().Json(http.StatusInternalServerError, map[string]string{"error": err.Error()})
	}

	// Return response
	return ctx.Response().Json(http.StatusOK, map[string]string{
		"status_name":   "successful",
		"tracking_link": trackingLink,
		"postback_link": postbackLink,
	})
}


func validateAdsCampaignInput(data models.AdsCampaign) (map[string]any, error) {
	validator, err := facades.Validation().Make(data, models.AdsCampaignRules)
	if err != nil {
		return nil, fmt.Errorf("validation error: %v", err)
	}
	if validator.Fails() {
		return map[string]any{
			"errors": validator.Errors().All(),
		}, nil
	}

	return  nil, nil
}

func ValidateAdsCampaignPostbackInput(inputs []models.AdsCampaignPostback) (map[string]any, error) {
	for i, pb := range inputs {
		validator, err := facades.Validation().Make(pb, models.AdsCampaignPostbackRules)
		if err != nil {
			return nil, fmt.Errorf("validation error: %v", err)
		}
		if validator.Fails() {
			return map[string]any{
				"index":  i,
				"errors": validator.Errors().All(),
			}, nil
		}
	}
	return nil, nil
}



func generateUniqueCode() (string, error) {
	const letters = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789"
	length := 16;
    for {
        code := make([]byte, length)
        for i := range code {
            n, err := rand.Int(rand.Reader, big.NewInt(int64(len(letters))))
            if err != nil {
                return "", err
            }
            code[i] = letters[n.Int64()]
        }

        codeStr := string(code)

        count, err := facades.Orm().Query().Table("ads_campaigns").Where("code", codeStr).Count()
        if err != nil {
            return "", err
        }
        if count == 0 {
            return codeStr, nil
        }
        // else retry
    }
}




func (a *AddAdsCampaignController) ShowSupportParameter(ctx http.Context) http.Response {
	return ctx.Response().Json(http.StatusOK, map[string]any{
		"support_parameter":models.AllowedEventDataFields,
	})
}
