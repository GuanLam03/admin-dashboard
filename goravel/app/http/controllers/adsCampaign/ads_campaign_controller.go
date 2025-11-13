package adsCampaign

import (
	"strings"
	"time"
	"github.com/goravel/framework/contracts/http"
	"github.com/goravel/framework/facades"
	
	"goravel/app/models"
)

type AdsCampaignController struct {
}

func NewAdsCampaignController() *AdsCampaignController {
	return &AdsCampaignController{}
}

func (a *AdsCampaignController) ListAdsCampaigns(ctx http.Context) http.Response {

	adsCampaigns, err := a.filter(ctx)
	if err != nil {
		return ctx.Response().Json(500, map[string]string{"error":facades.Lang(ctx).Get("validation.internal_error")})
	}

	var status  = models.AdsCampaignStatusMap

	return ctx.Response().Json(200, map[string]any{
		"ads_campaigns": adsCampaigns,
		"status" : status,
	})
}

func (a *AdsCampaignController) filter(ctx http.Context) ([]models.AdsCampaign, error) {
	query := facades.Orm().Query()

	if name := ctx.Request().Query("name"); name != "" {
		name = strings.Trim(name, "\"'")
		query = query.Where("name LIKE ?", "%"+name+"%")
	}

	if targetUrl := ctx.Request().Query("target_url"); targetUrl != "" {
		query = query.Where("target_url LIKE ?", "%" + targetUrl + "%")
	}

	if status := ctx.Request().Query("status"); status != "" {
		query = query.Where("status", status)
	}

	if fdate := ctx.Request().Query("fdate"); fdate != "" {
		if _, err := time.Parse("2006-01-02", fdate); err == nil {
			query = query.Where("created_at >= ?", fdate)
		}
	}
	if tdate := ctx.Request().Query("tdate"); tdate != "" {
		if _, err := time.Parse("2006-01-02", tdate); err == nil {
			query = query.Where("created_at <= ?", tdate+" 23:59:59")
		}
	}

	// Fetch from DB
	var adsCampaigns []models.AdsCampaign
	if err := query.Get(&adsCampaigns); err != nil {
		return nil, err
	}

	return adsCampaigns, nil
}
