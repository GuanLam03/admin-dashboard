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

	now := time.Now()

	for i := 0; i < totalRecords; i++ {
		createdAt := now.Add(time.Duration(i) * time.Second).Format("2006-01-02 15:04:05")
		updatedAt := createdAt

		ip := fmt.Sprintf("192.168.%d.%d", rand.Intn(255), rand.Intn(255))
		country := countries[rand.Intn(len(countries))]
		city := cities[rand.Intn(len(cities))]
		deviceType := deviceTypes[rand.Intn(len(deviceTypes))]
		deviceName := deviceNames[rand.Intn(len(deviceNames))]
		osName := osNames[rand.Intn(len(osNames))]
		browser := browsers[rand.Intn(len(browsers))]

		// --- ads_log_details
		write("INSERT INTO ads_log_details (id, ip, country, region, city, user_agent, referrer, device_type, device_name, os_name, os_version, browser_name, browser_version, ads_campaign_id, clicked_url, created_at, updated_at) VALUES ")
		write("(%d, '%s', '%s', 'Selangor', '%s', 'Mozilla/5.0 (Windows NT 10.0; Win64; x64)', 'https://google.com', '%s', '%s', '%s', '%d.%d', '%s', '%d.0.%d', %d, '%s', '%s', '%s');\n",
			detailID, ip, country, city, deviceType, deviceName, osName,
			rand.Intn(15)+1, rand.Intn(10), browser, rand.Intn(120)+60, rand.Intn(100),
			campaignID, fmt.Sprintf("https://example.com/product/%d", rand.Intn(10000)),
			createdAt, updatedAt)

		// --- ads_logs
		write("INSERT INTO ads_logs (id, ads_campaign_id, ads_log_detail_id, created_at, updated_at) VALUES ")
		write("(%d, %d, %d, '%s', '%s');\n",
			logID, campaignID, detailID, createdAt, updatedAt)

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

			write("INSERT INTO ads_event_logs (id, ads_log_id, event_name, data, created_at, updated_at) VALUES ")
			write("(%d, %d, '%s', '%s', '%s', '%s');\n", eventID, logID, eventName, escapeQuotes(string(dataJSON)), createdAt, updatedAt)
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

	return s
}
