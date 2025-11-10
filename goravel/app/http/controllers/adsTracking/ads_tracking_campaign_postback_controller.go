package adsTracking

import (
	"fmt"
	"io"
	httpRaw "net/http"
	"time"
	"encoding/json"

	// "github.com/goravel/framework/contracts/http"
	"github.com/goravel/framework/facades"

	"goravel/app/models"

)

type AdsTrackingCampaignPostbackController struct{}

func NewAdsTrackingCampaignPostbackController() *AdsTrackingCampaignPostbackController {
	return &AdsTrackingCampaignPostbackController{}
}

/*
Postback Processing Test Cases:

	Test Case 1: Successful 
		- HTTP status code: 200
		- Response body: empty OR contains normal data (not containing "error")
		- Result:
			status = "successful"
			error_message = null

	Test Case 2: Failed (error in JSON)
		- HTTP status code: 200
		- Response body: JSON containing an error {"error": "Some error message"}
		- Result:
			status = "failed"
			error_message = "Some error message"

	Test Case 3: Failed (HTTP error)
		- HTTP status code: not 200 (e.g., 404, 500)
		- Response body: any content
		- Result:
			status = "failed"
			error_message = null

*/


// ProcessPendingPostbacks handles triggering all pending postbacks
func (c *AdsTrackingCampaignPostbackController) ProcessPendingPostbacks() {
	var pendingLogs []models.AdsCampaignPostbackLog
	successCount := 0
	failCount := 0

	// Get all pending postbacks
	if err := facades.Orm().Query().
		Where("status", models.AdsCampaignPostbackLogStatusMap["pending"]).
		Find(&pendingLogs); err != nil {
		fmt.Println("[Postback] Failed to query pending postbacks:", err)
		return
	}

	if len(pendingLogs) == 0 {
		fmt.Println("[Postback] No pending postbacks found")
		return
	}

	for _, log := range pendingLogs {
		resp, err := httpRaw.Get(log.Url)
		if err != nil {
			errStr := err.Error()
			c.updatePostback(log.ID, map[string]any{
				"error_message": errStr,
				"status":        models.AdsCampaignPostbackLogStatusMap["failed"],
			})
			failCount++
			continue
		}

		body, _ := io.ReadAll(resp.Body)
		_ = resp.Body.Close()

		bodyStr := string(body)
		code := resp.StatusCode
		var status string
		var errMsg *string

		// Try to detect JSON error
		if len(body) > 0 {
			var jsonResp map[string]any
			if json.Unmarshal(body, &jsonResp) == nil {
				if errVal, ok := jsonResp["error"]; ok && errVal != nil { // 'ok' is true/false
					errText := fmt.Sprintf("%v", errVal)
					errMsg = &errText
				}
			}
		}


		if errMsg == nil && code == 200 {
			status = models.AdsCampaignPostbackLogStatusMap["successful"]
		} else {
			status = models.AdsCampaignPostbackLogStatusMap["failed"]
		}

		updateData := map[string]any{
			"response_status": code,
			"response_body":   bodyStr,
			"status":          status,
		}
		if errMsg != nil {
			updateData["error_message"] = *errMsg
		}

		c.updatePostback(log.ID, updateData)

		if status == models.AdsCampaignPostbackLogStatusMap["successful"] {
			successCount++
		} else {
			failCount++
		}
	}

	fmt.Printf("[Postback] Completed — %d success, %d failed\n", successCount, failCount)
}



func (c *AdsTrackingCampaignPostbackController) updatePostback(id uint, updates map[string]any) {
	updates["updated_at"] = time.Now()
	if _,err := facades.Orm().Query().
		Table("ads_campaign_postback_logs").
		Where("id", id).
		Update(updates); err != nil {
		facades.Log().Errorf("Failed to update postback log (id=%d): %v", id, err)
	}
}

