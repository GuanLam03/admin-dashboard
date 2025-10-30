package adsCampaign

import (
	"fmt"
	"math/rand"
	"time"
	"github.com/goravel/framework/contracts/http"
	"github.com/goravel/framework/facades"
	// "encoding/json"
	"goravel/app/models"
	"github.com/goravel/framework/contracts/database/orm"
)

type AddAdsCampaignController struct {
}

func NewAddAdsCampaignController() *AddAdsCampaignController {
	return &AddAdsCampaignController{}
}


func (a *AddAdsCampaignController) AddAdsCampaign(ctx http.Context) http.Response {
	basedUrl := facades.Config().Env("APP_URL", "")
	port := facades.Config().Env("APP_PORT", "")

	type PostbackEvent struct {
		EventName string `json:"event_name"`
		URL       string `json:"url"`
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
			"error": "Invalid request body",
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
		return ctx.Response().Json(http.StatusInternalServerError, map[string]string{"error": err.Error()})
	}
	if errResp != nil {
		return ctx.Response().Json(http.StatusUnprocessableEntity, errResp)
	}

	// Generate unique code
	code, err := generateUniqueCode()
	if err != nil {
		return ctx.Response().Json(http.StatusInternalServerError, map[string]any{
			"error": "Failed to generate campaign code",
		})
	}
	adsCampaign.Code = code

	// Set tracking and postback links
	trackingLink := fmt.Sprintf("%s:%s/%s", basedUrl, port, adsCampaign.Code)
	adsCampaign.TrackingLink = &trackingLink

	postbackLink := fmt.Sprintf("%s:%s/postback/", basedUrl, port)
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
					PostbackUrl:   e.URL,
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


func validateAdsCampaignInput(data models.AdsCampaign) (map[string]interface{}, error) {
	validator, err := facades.Validation().Make(data, models.AdsCampaignRules)
	if err != nil {
		return nil, fmt.Errorf("validation error: %v", err)
	}
	if validator.Fails() {
		return map[string]interface{}{
			"errors": validator.Errors().All(),
		}, nil
	}

	return  nil, nil
}


func generateUniqueCode() (string, error) {
	const letters = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789"
	rand.Seed(time.Now().UnixNano())

	for {
		codeBytes := make([]byte, 8)
		for i := range codeBytes {
			codeBytes[i] = letters[rand.Intn(len(letters))]
		}
		code := string(codeBytes)

		// check if exists
		count, err := facades.Orm().Query().Table("ads_campaigns").Where("code", code).Count()
		if err != nil {
			return "", err
		}

		if count == 0 {
			return code, nil
		}
		// else retry
	}
}



func (a *AddAdsCampaignController) ShowSupportParameter(ctx http.Context) http.Response {
	return ctx.Response().Json(http.StatusOK, map[string]any{
		"support_parameter":models.AllowedEventDataFields,
	})
}
