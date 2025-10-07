package adsTracking

import (
	"fmt"
	"encoding/json"
	"github.com/goravel/framework/contracts/http"
	"github.com/goravel/framework/facades"
	"strconv"
	"goravel/app/models"
)

type AdsTrackingController struct {
}

func NewAdsTrackingController() *AdsTrackingController {
	return &AdsTrackingController{}
}


// GeoIP API response structure
type geoResponse struct {
	Country    string `json:"country"`
	RegionName string `json:"regionName"`
	City       string `json:"city"`
}

// Helper to get geo info from IP using ip-api.com
func getGeoInfo(ip string) (*geoResponse, error) {
	resp, err := facades.Http().Get(fmt.Sprintf("http://ip-api.com/json/%s", ip))
	if err != nil {
		return nil, err
	}

	body, err := resp.Body() 
	if err != nil {
		return nil, err
	}

	var data geoResponse
	if err := json.Unmarshal([]byte(body), &data); err != nil {
		return nil, err
	}

	return &data, nil
}

func (a *AdsTrackingController) Track(ctx http.Context) http.Response {
	code := ctx.Request().Route("code")

	// Find the campaign
	var campaign models.AdsCampaign
	if err := facades.Orm().Query().Where("code", code).First(&campaign); err != nil || campaign.ID == 0 {
		return ctx.Response().Json(http.StatusNotFound, map[string]string{
			"error": "Campaign not found",
		})
	}


	ip := ctx.Request().Ip()
	userAgent := ctx.Request().Header("User-Agent")
	referrer := ctx.Request().Header("Referer")

	// Get geolocation info from IP
	var country, region, city *string
	if geo, err := getGeoInfo(ip); err == nil {
		country = &geo.Country
		region = &geo.RegionName
		city = &geo.City
	}


	log := models.AdsLog{
		AdsCampaignId: campaign.ID,
		Ip:            &ip,
		Country:       country,
		Region:        region,
		City:          city,
		UserAgent:     &userAgent,
		Referrer:      &referrer,
		Converted:     false,
	}

	if err := facades.Orm().Query().Create(&log); err != nil {
		return ctx.Response().Json(http.StatusInternalServerError, map[string]string{
			"error": "Failed to log ad tracking",
		})
	}

	clientUrl := fmt.Sprintf("%s?log_id=%d",campaign.TargetUrl,log.ID)

	return ctx.Response().Redirect(http.StatusFound, clientUrl)
}

func (a *AdsTrackingController) PostBackAdsTracking(ctx http.Context) http.Response {
	logId := ctx.Request().Query("log_id")
	productId := ctx.Request().Query("product_id")
	value := ctx.Request().Query("value")

	if logId == "" {
		return ctx.Response().Json(http.StatusBadRequest, http.Json{
			"error": "log_id is required",
		})
	}

	// Find the existing log record
	var adsLog models.AdsLog
	if err := facades.Orm().Query().Find(&adsLog, logId); err != nil {
		return ctx.Response().Json(http.StatusNotFound, http.Json{
			"error": "Log not found",
		})
	}

	// Parse value 
	var parsedValue *float64
	if value != "" {
		v, err := strconv.ParseFloat(value, 64)
		if err == nil {
			parsedValue = &v
		}
	}

	// Parse productId to uint
	var parsedProductID *uint
	if productId != "" {
		if pid, err := strconv.ParseUint(productId, 10, 64); err == nil {
			pidUint := uint(pid)
			parsedProductID = &pidUint
		}
	}

	adsLog.Converted = true
	adsLog.ClientProductId = parsedProductID
	adsLog.Value = parsedValue

	if err := facades.Orm().Query().Save(&adsLog); err != nil {
		return ctx.Response().Json(http.StatusInternalServerError, http.Json{
			"error": "Failed to update log: " + err.Error(),
		})
	}

	return ctx.Response().Json(http.StatusOK, http.Json{
		"status":  "success",
		"message": "Log updated successfully",
		"data":    adsLog,
	})
}
