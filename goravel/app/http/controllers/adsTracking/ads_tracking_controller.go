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

	"io"
	httpRaw "net/http"
	"net/url"
	"regexp"
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
			ClickedUrl:  clickedUrl,
			AdsCampaignId:  campaign.ID,

		}

		if err := tx.Create(&logDetail); err != nil {
			return err
		}

		log = models.AdsLog{
			AdsCampaignId:  campaign.ID,
			AdsLogDetailId: logDetail.ID,
			
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

	clientUrl := fmt.Sprintf("%s?ads_log_id=%d", campaign.TargetUrl, log.ID)
	return ctx.Response().Redirect(http.StatusFound, clientUrl)
}


func (a *AdsTrackingController) PostBackAdsTracking(ctx http.Context) http.Response {
	var req struct {
		EventName string                 `json:"event_name"`
		AdsLogId  string                 `json:"ads_log_id"`
		Data      map[string]interface{} `json:"data"`
	}

	if err := ctx.Request().Bind(&req); err != nil {
		return ctx.Response().Json(http.StatusBadRequest, map[string]string{
			"code":        "400",
			"status_name": "bad_request",
		})
	}

	// Validation error
	if errors, err := validatePostBackAdsTrackingInput(req); err != nil || errors != nil {
		return ctx.Response().Json(http.StatusBadRequest, map[string]string{
			"code":        "400",
			"status_name": "bad_request",
		})
	}

	// Parse ads_log_id
	adsLogID, err := strconv.ParseUint(req.AdsLogId, 10, 64)
	if err != nil {
		return ctx.Response().Json(http.StatusBadRequest, map[string]string{
			"code":        "400",
			"status_name": "bad_request",
		})
	}

	// Check if AdsLog exists
	exists, err := facades.Orm().Query().
		Model(&models.AdsLog{}).
		Where("id", adsLogID).
		Exists()
	if err != nil {
		return ctx.Response().Json(http.StatusInternalServerError, map[string]string{
			"code":        "500",
			"status_name": "internal_server_error",
		})
	}
	if !exists {
		return ctx.Response().Json(http.StatusNotFound, map[string]string{
			"code":        "404",
			"status_name": "not_found",
		})
	}

	// Create event log
	jsonData, _ := json.Marshal(req.Data)
	eventLog := models.AdsEventLog{
		AdsLogId:  uint(adsLogID),
		EventName: req.EventName,
		Data:      datatypes.JSON(jsonData),
	}
	if err := facades.Orm().Query().Create(&eventLog); err != nil {
		return ctx.Response().Json(http.StatusInternalServerError, map[string]string{
			"code":        "500",
			"status_name": "internal_server_error",
		})
	}
	// Trigger background postback
	go sendCampaignPostback(eventLog)

	return ctx.Response().Json(http.StatusOK, map[string]string{
		"code":        "200",
		"status_name": "successful",
	})
}



func validatePostBackAdsTrackingInput(req interface{}) (map[string]interface{}, error) {
	// Convert input struct > map for validator
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

	// Check allowed "data" fields
	allowed := map[string]bool{}
	for _, k := range models.AllowedEventDataFields {
		allowed[k] = true
	}

	if data, ok := payload["data"].(map[string]interface{}); ok {
		for key := range data {
			if !allowed[key] {
				// facades.Log().Errorf("invalid field '%s' in data", key)
				return map[string]interface{}{
					"errors": fmt.Sprintf("invalid field '%s' in data", key),
				}, nil
			}
		}
	}

	return nil, nil
}


func sendCampaignPostback(eventLog models.AdsEventLog) {
	// Get campaign ID from ads_logs
	var campaignID uint
	if err := facades.Orm().Query().
		Table("ads_logs").
		Select("ads_campaign_id").
		Where("id", eventLog.AdsLogId).
		Get(&campaignID); 
		err != nil {
			facades.Log().Errorf("Failed to get campaign ID for ads_log_id=%d: %v", eventLog.AdsLogId, err)
			return
		}

	// Find postback url 
	var postback models.AdsCampaignPostback
	if err := facades.Orm().Query().
		Where("ads_campaign_id", campaignID).
		Where("event_name", eventLog.EventName).
		First(&postback);
		err != nil {
			facades.Log().Errorf("Postback query error (campaign_id=%d, event=%s): %v", campaignID, eventLog.EventName, err)
			return
		}

	if postback.ID == 0 {
		facades.Log().Infof("No matching postback found (campaign_id=%d, event=%s) — skipping", campaignID, eventLog.EventName)
		return
	}

	// Parse event data
	var data map[string]interface{}
	_ = json.Unmarshal(eventLog.Data, &data)

	// Build placeholders dynamically from allowed data fields
	placeholders := make(map[string]string)
	for _, key := range models.AllowedEventDataFields {
		if val, ok := data[key]; ok {
			placeholders[key] = fmt.Sprintf("%v", val)
		} else {
			placeholders[key] = ""
		}
	}

	// Get click_id from ads_log_details
	var adsLogDetail models.AdsLogDetail
	if err := facades.Orm().Query().
		Table("ads_log_details AS ald").
		Join("INNER JOIN ads_logs AS al ON ald.id = al.ads_log_detail_id").
		Where("ald.id", eventLog.AdsLogId).
		First(&adsLogDetail); 
		err != nil {
			facades.Log().Errorf("Failed to get AdsLogDetail for ads_log_id=%d: %v", eventLog.AdsLogId, err)
			return
		}

	clickID := extractClickID(adsLogDetail.ClickedUrl)
	// Add click_id as placeholder
	placeholders["click_id"] = clickID

	finalURL := replacePlaceholders(postback.PostbackUrl, placeholders)

	// Send HTTP request (Get by default)
	resp, err := httpRaw.Get(finalURL)
	fmt.Println("res",resp)
	var (
		statusCode  *int
		respBodyStr *string
		errorMsg    *string
	)

	if err != nil {
		errStr := err.Error()
		errorMsg = &errStr
	} else {
		defer resp.Body.Close()
		body, _ := io.ReadAll(resp.Body)
		bodyStr := string(body)
		respBodyStr = &bodyStr
		code := resp.StatusCode
		statusCode = &code
	}

	// Log result
	postbackLog := models.AdsCampaignPostbackLog{
		AdsEventLogId:         &eventLog.ID,
		AdsCampaignPostbackId: postback.ID,
		Url:                   finalURL,
		RequestMethod:         "Get",
		ResponseStatus:        statusCode,
		ResponseBody:          respBodyStr,
		ErrorMessage:          errorMsg,
	}
	if err := facades.Orm().Query().Create(&postbackLog); err != nil {
		facades.Log().Errorf("Failed to save postback log (campaign_id=%d, event=%s): %v", campaignID, eventLog.EventName, err)
	}
}

//Replace placeholders dynamically
func replacePlaceholders(urlTemplate string, values map[string]string) string {
	re := regexp.MustCompile(`\{([a-zA-Z0-9_]+)\}`)
	return re.ReplaceAllStringFunc(urlTemplate, func(match string) string {
		key := re.FindStringSubmatch(match)[1]
		if val, ok := values[key]; ok {
			return url.QueryEscape(val)
		}
		return match
	})
}

func extractClickID(clickedURL string) string {
	u, err := url.Parse(clickedURL)
	if err != nil {
		return ""
	}

	query := u.Query()
	candidateKeys := []string{
		"cid", "click_id", "clickid",
		"fbclid", "gclid", "ttclid", "msclkid", "twclid", "gbraid", "wbraid",
	}

	for _, key := range candidateKeys {
		if val := query.Get(key); val != "" {
			return val
		}
	}

	return ""
}


func getString(data map[string]interface{}, key string) string {
	if val, ok := data[key]; ok {
		if s, ok := val.(string); ok {
			return s
		}
	}
	return ""
}