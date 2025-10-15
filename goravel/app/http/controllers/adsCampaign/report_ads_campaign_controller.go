package adsCampaign

import (
	// "strconv"
	"github.com/goravel/framework/contracts/http"
	"github.com/goravel/framework/facades"
	// "goravel/app/models"
)

type ReportAdsCampaignController struct{}

func NewReportAdsCampaignController() *ReportAdsCampaignController {
	return &ReportAdsCampaignController{}
}

// func (r *ReportAdsCampaignController) ShowReportAdsCampaign(ctx http.Context) http.Response {
// 	idStr := ctx.Request().Route("campaign_id")
// 	campaignID, _ := strconv.Atoi(idStr)

// 	var results []map[string]interface{}

// 	err := facades.Orm().Query().
// 		Table("ads_logs").
// 		Select("ads_logs.id,ads_logs.clicked_url,ads_logs.ads_log_detail_id,ads_logs.ads_campaign_id,ads_log_details.ip,ads_log_details.country").
// 		Join("inner join ads_log_details on ads_logs.ads_log_detail_id = ads_log_details.id").
// 		Where("ads_logs.ads_campaign_id", campaignID).
// 		Get(&results)

// 	if err != nil {
// 		return ctx.Response().Json(http.StatusInternalServerError, http.Json{
// 			"error": err.Error(),
// 		})
// 	}

// 	return ctx.Response().Json(http.StatusOK, http.Json{
// 		"data": results,
// 	})
// }



// func (r *ReportAdsCampaignController) ShowReportAdsCampaign(ctx http.Context) http.Response {
// 	campaignID := ctx.Request().Route("campaign_id") // Get campaign ID from route

// 	var results []map[string]interface{}

// 	err := facades.Orm().Query().
// 		Table("ads_logs").
// 		Select("ads_logs.*, ads_log_details.ip, ads_log_details.country"). // select what you need
// 		Join("inner join ads_log_details on ads_logs.ads_log_detail_id = ads_log_details.id").
// 		Where("ads_logs.ads_campaign_id", campaignID).
// 		Get(&results)

// 	if err != nil {
// 		return ctx.Response().Json(http.StatusInternalServerError, http.Json{
// 			"error": err.Error(),
// 		})
// 	}

// 	return ctx.Response().Json(http.StatusOK, http.Json{
// 		"data": results,
// 	})
// }


// func (r *ReportAdsCampaignController) ShowReportAdsCampaign(ctx http.Context) http.Response {
// 	// campaignID := ctx.Request().Route("campaign_id") // Get campaign ID from route

// 	var results []map[string]interface{}

// 	err := facades.Orm().Query().
// 		Table("ads_log_details").
// 		Select("id, ads_log_details.country").
// 		// Join("inner join ads_log_details on ads_logs.ads_log_detail_id = ads_log_details.id").
// 		// Where("ads_logs.ads_campaign_id", campaignID).
// 		Where("ads_log_details.country", "Malaysia"). 
// 		Get(&results)

// 	if err != nil {
// 		return ctx.Response().Json(http.StatusInternalServerError, http.Json{
// 			"error": err.Error(),
// 		})
// 	}

// 	return ctx.Response().Json(http.StatusOK, http.Json{
// 		"data": results,
// 	})
// }


func (r *ReportAdsCampaignController) ShowReportAdsCampaign(ctx http.Context) http.Response {
	campaignID := ctx.Request().Route("campaign_id") // Get campaign ID from route

	

	num ,err := facades.Orm().Query().
		Table("ads_event_logs AS el").
		Join("inner join ads_logs AS al on el.ads_log_id = al.id").
		Where("al.ads_campaign_id", campaignID).
		Where("el.event_name", "CONTENT_VIEW").
		Count()

	if err != nil {
		return ctx.Response().Json(http.StatusInternalServerError, http.Json{
			"error": err.Error(),
		})
	}

	return ctx.Response().Json(http.StatusOK, http.Json{
		"data": num,
	})
}
