package adsCampaign

import (
	
	"github.com/goravel/framework/contracts/http"
	"github.com/goravel/framework/facades"
	"github.com/goravel/framework/contracts/database/orm"
	"goravel/app/models"
	// "goravel/app/helpers"
	"goravel/app/messages"


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
		return ctx.Response().Json(404, map[string]string{"error": messages.GetError("validation.ads_campaign_not_found")})
	}

	status := models.AdsCampaignStatusMap
	delete(status, "removed")

	// Select only id, event_name, and postback_url
	var campaignPostbacks []struct {
		ID          uint   `json:"id"`
		EventName   string `json:"event_name"`
		PostbackUrl string `json:"postback_url"`
		IncludeClickParams bool   `json:"include_click_params"`
	}

	if err := facades.Orm().Query().
		Table("ads_campaign_postbacks").
		Select("id", "event_name", "postback_url","include_click_params").
		Where("ads_campaign_id", adsCampaign.ID).
		Get(&campaignPostbacks); err != nil {
		return ctx.Response().Json(500, map[string]string{"error":messages.GetError("validation.internal_error")})
	}

	return ctx.Response().Json(200, map[string]any{
		"ads_campaign": adsCampaign,
		"status" : status,
		"ads_campaign_postbacks": campaignPostbacks,
	})
}


func (c *EditAdsCampaignController) EditAdsCampaign(ctx http.Context) http.Response {
	id := ctx.Request().Route("id")

	var (
		adsCampaign models.AdsCampaign
		input struct {
			Name           string                       `json:"name"`
			TargetUrl      string                       `json:"target_url"`
			Status         string                       `json:"status"`
			PostbackEvents []models.AdsCampaignPostback  `json:"postback_events"`
		}
		existingIDs []any
	)

	// Find existing campaign
	if err := facades.Orm().Query().Find(&adsCampaign, id); err != nil || adsCampaign.ID == 0 {
		return ctx.Response().Json(404, map[string]string{"error": messages.GetError("validation.ads_campaign_not_found")})
	}

	// Bind request body
	if err := ctx.Request().Bind(&input); err != nil {
		return ctx.Response().Json(400, map[string]any{"error": messages.GetError("validation.invalid_request")})
	}

	// Validate campaign fields
	errResp, err := validateAdsCampaignInput(models.AdsCampaign{
		Name:      input.Name,
		TargetUrl: input.TargetUrl,
		Status:    input.Status,
	})


	if err != nil {
		return ctx.Response().Json(500, map[string]string{"error": messages.GetError("validation.internal_error")})
	}
	if errResp != nil {
		return ctx.Response().Json(422, errResp)
	}

	var postbackModels []models.AdsCampaignPostback
	for _, pb := range input.PostbackEvents {
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
			return ctx.Response().Json(http.StatusInternalServerError, map[string]string{"error": messages.GetError("validation.internal_error")})
		}
		if errResp != nil {
			return ctx.Response().Json(http.StatusUnprocessableEntity, errResp)
		}
	}


	if err := facades.Orm().Transaction(func(tx  orm.Query) error {

		// Update campaign fields
		adsCampaign.Name = input.Name
		adsCampaign.TargetUrl = input.TargetUrl
		adsCampaign.Status = input.Status
		if err := tx.Save(&adsCampaign); err != nil {
			return err
		}

		// Handle Postback Events (create / update)
		for _, pb := range input.PostbackEvents {
			if pb.ID > 0 {
				// Update existing
				var existing models.AdsCampaignPostback
				if err := tx.Find(&existing, pb.ID); err == nil && existing.ID > 0 {
					existing.EventName = pb.EventName
					existing.PostbackUrl = pb.PostbackUrl
					existing.IncludeClickParams = pb.IncludeClickParams
					if err := tx.Save(&existing); err != nil {
						return err
					}
					existingIDs = append(existingIDs, existing.ID)
				}
			} else {
				// Create new
				newPB := models.AdsCampaignPostback{
					AdsCampaignId:      adsCampaign.ID,
					EventName:          pb.EventName,
					PostbackUrl:        pb.PostbackUrl,
					IncludeClickParams: pb.IncludeClickParams,
				}
				if err := tx.Create(&newPB); err != nil {
					return err
				}
				existingIDs = append(existingIDs, newPB.ID)
			}
		}

		// Delete old postbacks not in request
		query := tx.Where("ads_campaign_id", adsCampaign.ID)
		if len(existingIDs) > 0 {
			query = query.WhereNotIn("id", existingIDs)  // goravel orm whereNotIn needs []any
		}
		if _,err := query.Delete(&models.AdsCampaignPostback{}); err != nil {
			return err
		}

		return nil
	}); err != nil {
		return ctx.Response().Json(500, map[string]any{
			"error":   messages.GetError("validation.internal_error"),
			"details": err.Error(),
		})
	}

	
	// helpers.Activity().
	// 	CausedBy(ctx).
	// 	InLog("Ads Campaign").
	// 	OnUrl(ctx.Request().FullUrl()).
	// 	OnRoute(ctx.Request().Method() + " " + ctx.Request().Path()).
	// 	WithInput(ctx.Request().All()).
	// 	Log("Edit ads campaign")
	
	return ctx.Response().Json(200, map[string]any{
		"message": "Ads Campaign updated successfully",
		"data":    adsCampaign,
	})
}

