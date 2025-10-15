// package seeders

// import (
// 	"encoding/json"
// 	"fmt"
// 	"math/rand"
// 	"time"

// 	"goravel/app/models"
// 	"github.com/goravel/framework/facades"
// 	"github.com/goravel/framework/contracts/database/orm"
// 	"gorm.io/datatypes"
// )

// type AdsTrackingSeeder struct{}

// func (s *AdsTrackingSeeder) Signature() string {
// 	return "AdsTrackingSeeder"
// }

// func (s *AdsTrackingSeeder) Run() error {
// 	rand.Seed(time.Now().UnixNano())

// 	campaignID := uint(22)
// 	totalRecords := 1000000

// 	eventNames := []string{
// 		"CONTENT_VIEW",
// 		"ADD_TO_CART",
// 		"PURCHASE",
// 		"SUBSCRIBE",
// 	}

// 	deviceTypes := []string{"mobile", "desktop", "tablet"}
// 	deviceNames := []string{"iPhone 14", "Galaxy S23", "MacBook Pro", "Windows PC"}
// 	osNames := []string{"iOS", "Android", "Windows", "macOS"}
// 	browsers := []string{"Chrome", "Safari", "Firefox", "Edge"}
// 	countries := []string{"Malaysia", "Singapore", "Thailand", "Vietnam", "Philippines"}
// 	cities := []string{"Kuala Lumpur", "Petaling Jaya", "Johor Bahru", "Penang", "Ipoh"}

// 	facades.Log().Infof("🚀 Starting seeder for %d ads logs (campaign %d)", totalRecords, campaignID)

// 	return facades.Orm().Transaction(func(tx orm.Query) error {
// 		for i := 0; i < totalRecords; i++ {
// 			// Ads Log
// 			adsLog := models.AdsLog{
// 				AdsCampaignId: campaignID,
// 				ClickedUrl:    fmt.Sprintf("https://example.com/product/%d", rand.Intn(10000)),
// 			}
// 			if err := tx.Create(&adsLog); err != nil {
// 				return err
// 			}

// 			// Ads Log Detail
// 			ip := fmt.Sprintf("192.168.%d.%d", rand.Intn(255), rand.Intn(255))
// 			userAgent := "Mozilla/5.0 (Windows NT 10.0; Win64; x64)"

// 			detail := models.AdsLogDetail{
// 				Ip:             &ip,
// 				Country:        ptr(countries[rand.Intn(len(countries))]),
// 				Region:         ptr("Selangor"),
// 				City:           ptr(cities[rand.Intn(len(cities))]),
// 				UserAgent:      &userAgent,
// 				Referrer:       ptr("https://google.com"),
// 				DeviceType:     ptr(deviceTypes[rand.Intn(len(deviceTypes))]),
// 				DeviceName:     ptr(deviceNames[rand.Intn(len(deviceNames))]),
// 				OsName:         ptr(osNames[rand.Intn(len(osNames))]),
// 				OsVersion:      ptr(fmt.Sprintf("%d.%d", rand.Intn(15)+1, rand.Intn(10))),
// 				BrowserName:    ptr(browsers[rand.Intn(len(browsers))]),
// 				BrowserVersion: ptr(fmt.Sprintf("%d.0.%d", rand.Intn(120)+60, rand.Intn(100))),
// 			}

// 			if err := tx.Create(&detail); err != nil {
// 				return err
// 			}

// 			// Link detail to log
// 			adsLog.AdsLogDetailId = detail.ID
// 			if err := tx.Save(&adsLog); err != nil {
// 				return err
// 			}

// 			// Event Logs
// 			numEvents := rand.Intn(3) + 1 // 1–3 events
// 			for j := 0; j < numEvents; j++ {
// 				eventName := eventNames[rand.Intn(len(eventNames))]
// 				dataMap := map[string]interface{}{
// 					"content_id":   fmt.Sprintf("P-%d", rand.Intn(99999)),
// 					"content_type": "product",
// 					"currency":     "USD",
// 					"value":        fmt.Sprintf("%.2f", rand.Float64()*100),
// 					"quantity":     rand.Intn(5) + 1,
// 					"price":        fmt.Sprintf("%.2f", rand.Float64()*50),
// 				}
// 				dataJSON, _ := json.Marshal(dataMap)

// 				event := models.AdsEventLog{
// 					AdsLogId:  adsLog.ID,
// 					EventName: eventName,
// 					Data:      datatypes.JSON(dataJSON),
// 				}
// 				if err := tx.Create(&event); err != nil {
// 					return err
// 				}
// 			}

// 			// Log progress every 10,000 rows
// 			if (i+1)%10000 == 0 {
// 				facades.Log().Infof("✅ Inserted %d/%d ads logs", i+1, totalRecords)
// 			}
// 		}
// 		return nil
// 	})
// }

// func ptr[T any](v T) *T {
// 	return &v
// }


package seeders

import (
	"encoding/json"
	"fmt"
	"math/rand"
	"os"
	"time"
)

type AdsTrackingSeeder struct{}

func (s *AdsTrackingSeeder) Signature() string {
	return "AdsTrackingSeeder"
}

func (s *AdsTrackingSeeder) Run() error {
	rand.Seed(time.Now().UnixNano())

	campaignID := 22
	totalRecords := 1000000

	eventNames := []string{"CONTENT_VIEW", "ADD_TO_CART", "PURCHASE", "SUBSCRIBE"}
	deviceTypes := []string{"mobile", "desktop", "tablet"}
	deviceNames := []string{"iPhone 14", "Galaxy S23", "MacBook Pro", "Windows PC"}
	osNames := []string{"iOS", "Android", "Windows", "macOS"}
	browsers := []string{"Chrome", "Safari", "Firefox", "Edge"}
	countries := []string{"Malaysia", "Singapore", "Thailand", "Vietnam", "Philippines"}
	cities := []string{"Kuala Lumpur", "Petaling Jaya", "Johor Bahru", "Penang", "Ipoh"}

	file, err := os.Create("ads_tracking_seed.sql")
	if err != nil {
		return err
	}
	defer file.Close()

	write := func(format string, a ...interface{}) {
		file.WriteString(fmt.Sprintf(format, a...))
	}

	write("-- SQL Data Seeder for Ads Tracking\n")
	write("START TRANSACTION;\n\n")

	logID := 1
	detailID := 1
	eventID := 1

	for i := 0; i < totalRecords; i++ {
		ip := fmt.Sprintf("192.168.%d.%d", rand.Intn(255), rand.Intn(255))
		country := countries[rand.Intn(len(countries))]
		city := cities[rand.Intn(len(cities))]
		deviceType := deviceTypes[rand.Intn(len(deviceTypes))]
		deviceName := deviceNames[rand.Intn(len(deviceNames))]
		osName := osNames[rand.Intn(len(osNames))]
		browser := browsers[rand.Intn(len(browsers))]

		// --- ads_log_details
		write("INSERT INTO ads_log_details (id, ip, country, region, city, user_agent, referrer, device_type, device_name, os_name, os_version, browser_name, browser_version) VALUES ")
		write("(%d, '%s', '%s', 'Selangor', '%s', 'Mozilla/5.0 (Windows NT 10.0; Win64; x64)', 'https://google.com', '%s', '%s', '%s', '%d.%d', '%s', '%d.0.%d');\n",
			detailID, ip, country, city, deviceType, deviceName, osName,
			rand.Intn(15)+1, rand.Intn(10), browser, rand.Intn(120)+60, rand.Intn(100))

		// --- ads_logs
		write("INSERT INTO ads_logs (id, ads_campaign_id, clicked_url, ads_log_detail_id) VALUES ")
		write("(%d, %d, 'https://example.com/product/%d', %d);\n",
			logID, campaignID, rand.Intn(10000), detailID)

		// --- ads_event_logs (1–3 per log)
		numEvents := rand.Intn(3) + 1
		for j := 0; j < numEvents; j++ {
			eventName := eventNames[rand.Intn(len(eventNames))]
			data := map[string]interface{}{
				"content_id":   fmt.Sprintf("P-%d", rand.Intn(99999)),
				"content_type": "product",
				"currency":     "USD",
				"value":        fmt.Sprintf("%.2f", rand.Float64()*100),
				"quantity":     rand.Intn(5) + 1,
				"price":        fmt.Sprintf("%.2f", rand.Float64()*50),
			}
			dataJSON, _ := json.Marshal(data)

			write("INSERT INTO ads_event_logs (id, ads_log_id, event_name, data) VALUES ")
			write("(%d, %d, '%s', '%s');\n", eventID, logID, eventName, escapeQuotes(string(dataJSON)))
			eventID++
		}

		logID++
		detailID++

		if (i+1)%10000 == 0 {
			fmt.Printf("✅ Generated %d/%d records\n", i+1, totalRecords)
		}
	}

	write("\nCOMMIT;\n")
	fmt.Println("🎉 SQL file generated: ads_tracking_seed.sql")
	return nil
}

func escapeQuotes(s string) string {
	return fmt.Sprintf("%s", string([]rune(s)))
}
