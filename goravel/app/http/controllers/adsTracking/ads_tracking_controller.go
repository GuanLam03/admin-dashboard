package adsTracking

import (
	"fmt"
	"encoding/json"
	"github.com/goravel/framework/contracts/http"
	"github.com/goravel/framework/facades"
	"github.com/goravel/framework/contracts/database/orm"
	"strconv"
	"goravel/app/models"
	"github.com/mileusna/useragent"
	"gorm.io/datatypes"
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

func getDeviceType(ua useragent.UserAgent) string {
    if ua.Mobile {
        return "mobile"
    }
    if ua.Tablet {
        return "tablet"
    }
    return "desktop"
}



func (a *AdsTrackingController) Track(ctx http.Context) http.Response {
	code := ctx.Request().Route("code")
	clickedUrl := ctx.Request().FullUrl()
	
	// Find campaign
	var campaign models.AdsCampaign
	if err := facades.Orm().Query().Where("code", code).First(&campaign); err != nil || campaign.ID == 0 {
		return ctx.Response().Json(http.StatusNotFound, map[string]string{
			"error": "Campaign not found",
		})
	}

	ip := ctx.Request().Ip()
	userAgent := ctx.Request().Header("User-Agent")
	referrer := ctx.Request().Header("Referer")

	ua := useragent.Parse(userAgent)
	deviceType := getDeviceType(ua)
	deviceName := ua.Device
	osName := ua.OS
	osVersion := ua.OSVersion
	browserName := ua.Name
	browserVersion := ua.Version


	var country, region, city *string
	if geo, err := getGeoInfo(ip); err == nil {
		country = &geo.Country
		region = &geo.RegionName
		city = &geo.City
	}

	var logDetail models.AdsLogDetail
	var log models.AdsLog


	err := facades.Orm().Transaction(func(tx orm.Query) error {
		logDetail = models.AdsLogDetail{
			Ip:             &ip,
			Country:        country,
			Region:         region,
			City:           city,
			UserAgent:      &userAgent,
			Referrer:       &referrer,
			DeviceType:     &deviceType,
			DeviceName:     &deviceName,
			OsName:         &osName,
			OsVersion:      &osVersion,
			BrowserName:    &browserName,
			BrowserVersion: &browserVersion,
		}

		if err := tx.Create(&logDetail); err != nil {
			return err
		}

		log = models.AdsLog{
			AdsCampaignId:  campaign.ID,
			AdsLogDetailId: logDetail.ID,
			ClickedUrl:  clickedUrl,
		}

		if err := tx.Create(&log); err != nil {
			return err
		}

		return nil
	})

	if err != nil {
		return ctx.Response().Json(http.StatusInternalServerError, map[string]string{
			"error": err.Error(),
		})
	}

	clientUrl := fmt.Sprintf("%s?log_id=%d", campaign.TargetUrl, log.ID)
	return ctx.Response().Redirect(http.StatusFound, clientUrl)
}


func (a *AdsTrackingController) PostBackAdsTracking(ctx http.Context) http.Response {
	var req struct {
		EventName string          `json:"event_name"`
		AdsLogId  string          `json:"ads_log_id"`
		Data      map[string]interface{} `json:"data"`
	}

	if err := ctx.Request().Bind(&req); err != nil {
		return ctx.Response().Json(http.StatusBadRequest, map[string]string{
			"error": "Invalid request body",
		})
	}

	if errors, err := validatePostBackAdsTrackingInput(req); err != nil {
		return ctx.Response().Json(http.StatusInternalServerError, map[string]string{
			"error": err.Error(),
		})
	} else if errors != nil {
		return ctx.Response().Json(http.StatusBadRequest, errors)
	}


	adsLogID, _ := strconv.ParseUint(req.AdsLogId, 10, 64)

	// Check if AdsLog exists
	exists, err := facades.Orm().Query().
		Model(&models.AdsLog{}).
		Where("id", adsLogID).
		Exists()
	if err != nil {
		return ctx.Response().Json(http.StatusInternalServerError, map[string]string{
			"error": "Database error",
		})
	}
	if !exists {
		return ctx.Response().Json(http.StatusNotFound, map[string]string{
			"error": "Ads Log not found",
		})
	}
	jsonData, _ := json.Marshal(req.Data)

	eventLog := models.AdsEventLog{
		AdsLogId:  uint(adsLogID),
		EventName: req.EventName,
		Data:      datatypes.JSON(jsonData),
	}

	if err := facades.Orm().Query().Create(&eventLog); err != nil {
		return ctx.Response().Json(http.StatusInternalServerError, map[string]string{
			"error": err.Error(),
		})
	}

	return ctx.Response().NoContent(http.StatusOK)
}


func validatePostBackAdsTrackingInput(req interface{}) (map[string]interface{}, error) {
	// Convert input struct → map for validator
	payload := map[string]any{}
	bytes, _ := json.Marshal(req)
	if err := json.Unmarshal(bytes, &payload); err != nil {
		return nil, fmt.Errorf("failed to parse input: %v", err)
	}

	validator, err := facades.Validation().Make(payload, models.AdsEventLogRules)
	if err != nil {
		return nil, fmt.Errorf("validation setup error: %v", err)
	}

	if validator.Fails() {
		return map[string]interface{}{
			"errors": validator.Errors().All(),
		}, nil
	}


	return nil, nil
}
