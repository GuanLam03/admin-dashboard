package adsCampaign

import (
	
	"github.com/goravel/framework/contracts/http"
	"github.com/goravel/framework/facades"
	
	"goravel/app/models"
)

type EditAdsCampaignController struct {
}

func NewEditAdsCampaignController() *EditAdsCampaignController {
	return &EditAdsCampaignController{}
}

func (e *EditAdsCampaignController) ShowAdsCampaign(ctx http.Context) http.Response {
	id := ctx.Request().Route("id")

	var adsCampaign models.AdsCampaign
	if err := facades.Orm().Query().Find(&adsCampaign, id); err != nil || adsCampaign.ID == 0 {
		return ctx.Response().Json(404, map[string]string{"error": "Ads Campaign not found"})
	}


	return ctx.Response().Json(200, map[string]any{
		"ads_campaign": adsCampaign,
	})
}

func (c *EditAdsCampaignController) EditAdsCampaign(ctx http.Context) http.Response {
	id := ctx.Request().Route("id")

	// exisiting data
	var adsCampaign models.AdsCampaign
	if err := facades.Orm().Query().Find(&adsCampaign, id); err != nil || adsCampaign.ID == 0 {
		return ctx.Response().Json(404, map[string]string{"error": "Ads Campaign not found"})
	}

	// edit data
	var editAdsCampaign models.AdsCampaign
	if err := ctx.Request().Bind(&editAdsCampaign); err != nil {
		return ctx.Response().Json(http.StatusBadRequest, map[string]any{
			"error": "Invalid request body",
		})
	}

	errResp, err := validateAdsCampaignInput(editAdsCampaign)
	if err != nil {
		return ctx.Response().Json(500, map[string]string{"error": err.Error()})
	}
	if errResp != nil {
		return ctx.Response().Json(422, errResp)
	}

	//update to existing data
	adsCampaign.Name = editAdsCampaign.Name
	adsCampaign.TargetUrl = editAdsCampaign.TargetUrl
	

	if err := facades.Orm().Query().Save(&adsCampaign); err != nil {
		return ctx.Response().Json(500, map[string]string{"error": "Failed to update Ads Campaign"})
	}

	return ctx.Response().Json(200, map[string]any{
		"message": "Ads Campaign updated successfully",
		"data":    adsCampaign,
	})
}
