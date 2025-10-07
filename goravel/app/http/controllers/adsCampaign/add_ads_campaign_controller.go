package adsCampaign

import (
	"fmt"
	"math/rand"
	"time"
	"github.com/goravel/framework/contracts/http"
	"github.com/goravel/framework/facades"
	
	"goravel/app/models"
)

type AddAdsCampaignController struct {
}

func NewAddAdsCampaignController() *AddAdsCampaignController {
	return &AddAdsCampaignController{}
}


func (a *AddAdsCampaignController) AddAdsCampaign(ctx http.Context) http.Response {

	basedUrl := facades.Config().Env("APP_URL", "")
	port := facades.Config().Env("APP_PORT", "")

	var adsCampaign models.AdsCampaign
	if err := ctx.Request().Bind(&adsCampaign); err != nil {
		return ctx.Response().Json(http.StatusBadRequest, map[string]any{
			"error": "Invalid request body",
		})
	}

	errResp, err := validateAdsCampaignInput(adsCampaign)
	if err != nil {
		return ctx.Response().Json(500, map[string]string{"error": err.Error()})
	}
	if errResp != nil {
		return ctx.Response().Json(422, errResp)
	}

	// Generate unique code for campaign
	code, err := generateUniqueCode()
	if err != nil {
		return ctx.Response().Json(http.StatusInternalServerError, map[string]any{
			"error": "Failed to generate campaign code",
		})
	}
	adsCampaign.Code = code

	trackingLink := fmt.Sprintf("%s:%s/%s", basedUrl, port, adsCampaign.Code)
	adsCampaign.TrackingLink = &trackingLink

	postbackLink := fmt.Sprintf("%s:%s/postback/?log_id={logId}&product_id={productId}&value={value}", basedUrl, port)
	adsCampaign.PostbackLink = &postbackLink


	if err:= facades.Orm().Query().Create(&adsCampaign); err != nil{
		return ctx.Response().Json(http.StatusInternalServerError,map[string]any{
			"error": err.Error(),
		})
	}

	res := map[string]string{
		"tracking_link" : trackingLink,
		"postback_link" : postbackLink,
	}



	return ctx.Response().Json(http.StatusOK,res)

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
