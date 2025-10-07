package adsCampaign

import (
	"github.com/goravel/framework/contracts/http"
	"github.com/goravel/framework/facades"
	"goravel/app/models"
)

type ReportAdsCampaignController struct{}

func NewReportAdsCampaignController() *ReportAdsCampaignController {
	return &ReportAdsCampaignController{}
}

func (r *ReportAdsCampaignController) ShowReportAdsCampaign(ctx http.Context) http.Response {
	campaignID := ctx.Request().Route("campaign_id")

	var records []models.AdsLog
	if err := facades.Orm().Query().Where("ads_campaign_id", campaignID).Get(&records); err != nil {
		return ctx.Response().Json(http.StatusInternalServerError, http.Json{
			"error": err.Error(),
		})
	}

	totalClicks := len(records)
	totalConversions := 0
	totalRevenue := 0.0
	countryCount := map[string]int{}

	for _, record := range records {
		// ✅ handle conversion count and revenue safely
		if record.Converted {
			totalConversions++
			if record.Value != nil {
				totalRevenue += *record.Value
			}
		}

		// ✅ handle nil country safely
		if record.Country != nil && *record.Country != "" {
			country := *record.Country
			countryCount[country]++
		}
	}

	return ctx.Response().Json(http.StatusOK, http.Json{
		"summary": http.Json{
			"total_clicks":      totalClicks,
			"total_conversions": totalConversions,
			"total_revenue":     totalRevenue,
		},
		"country_stats": countryCount,
		"data":          records,
	})
}
