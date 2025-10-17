package adsLogs

import (
	"strings"

	"github.com/goravel/framework/contracts/http"
	"github.com/goravel/framework/facades"

	"goravel/app/models"
)

type AdsLogController struct{}

func NewAdsLogController() *AdsLogController {
	return &AdsLogController{}
}

func (a *AdsLogController) ListAdsLogs(ctx http.Context) http.Response {
	adsLogs, err := a.filter(ctx)
	if err != nil {
		return ctx.Response().Json(500, map[string]any{"error": err.Error()})
	}

	return ctx.Response().Json(200, map[string]any{
		"ads_logs": adsLogs,
	})
}

func (a *AdsLogController) filter(ctx http.Context) ([]models.AdsLog, error) {
	query := facades.Orm().Query().Model(&models.AdsLog{})

	// --- Filters ---
	if ip := ctx.Request().Query("ip"); ip != "" {
		ip = strings.Trim(ip, "\"'")
		query = query.Where("ip LIKE ?", "%"+ip+"%")
	}

	if country := ctx.Request().Query("country"); country != "" {
		country = strings.Trim(country, "\"'")
		query = query.Where("country LIKE ?", "%"+country+"%")
	}

	if converted := ctx.Request().Query("converted"); converted != "" {
		if converted == "1" {
			query = query.Where("converted = ?", 1)
		} else if converted == "0" {
			query = query.Where("converted = ?", 0)
		}
	}

	var adsLogs []models.AdsLog
	if err := query.Order("id DESC").Get(&adsLogs); err != nil {
		return nil, err
	}

	return adsLogs, nil


}
