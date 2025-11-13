package adsCampaign

import (
	"strconv"
	"fmt"
	"sync"
	"github.com/goravel/framework/contracts/http"
	"github.com/goravel/framework/facades"
	"goravel/app/models"
		"github.com/goravel/framework/contracts/database/orm"

)

type ReportAdsCampaignController struct{}

func NewReportAdsCampaignController() *ReportAdsCampaignController {
	return &ReportAdsCampaignController{}
}




// func (r *ReportAdsCampaignController) ShowReportAdsCampaign(ctx http.Context) http.Response {
// 	campaignID := ctx.Request().Route("campaign_id")

// 	// Total Clicks
// 	totalClicks, err := facades.Orm().Query().
// 		Table("ads_log_details").
// 		Where("ads_campaign_id", campaignID).
// 		Count()

// 	if err != nil {
// 		return ctx.Response().Json(http.StatusInternalServerError, http.Json{
// 			"error": "Failed to get total clicks: " + err.Error(),
// 		})
// 	}

// 	// Country Stats (clicks per country)
// 	// var countryStats []struct {
// 	// 	Country string
// 	// 	Count   int
// 	// }

// 	// err = facades.Orm().Query().
// 	// 	Table("ads_log_details").
// 	// 	Select("country, COUNT(id) AS count").
// 	// 	Where("ads_campaign_id", campaignID).
// 	// 	Group("country").
// 	// 	Get(&countryStats)

// 	var countryStats []struct {
// 		Country string
// 		Count   int
// 		TotalRevenue float64
// 	}

// 	err := facades.Orm().Query().
// 		Table("ads_logs AS al").
// 		Select("country, COUNT(al.id) AS count, SUM(ael.value_extracted) AS total_revenue").
// 		Join("INNER JOIN ads_log_details AS ald ON al.ads_log_detail_id = ald.id").
// 		Join("LEFT JOIN ads_event_logs AS ael ON ael.ads_log_id = al.id AND ael.event_name IN ('PURCHASE', 'SUBSCRIBE')").
// 		Where("al.ads_campaign_id", campaignID).
// 		Group("country").
// 		Get(&countryStats)

	

// 	if err != nil {
// 		return ctx.Response().Json(http.StatusInternalServerError, http.Json{
// 			"error": "Failed to get country stats: " + err.Error(),
// 		})
// 	}


// 	// Transform into map[string]int
// 	// countryCount := map[string]int{}
// 	// for _, row := range countryStats {
// 	// 	countryCount[row.Country] = row.Count
// 	// }

	

// 	var conversions struct {
// 		TotalConversions int
// 		TotalRevenue    float64
// 	}

// 	// Total Conversions
// 	 err = facades.Orm().Query().
// 		Table("ads_event_logs AS ael").
// 		Select("COUNT(DISTINCT ael.ads_log_id) AS total_conversions, SUM(value_extracted) AS total_revenue").
// 		Join("inner join ads_logs AS al ON ael.ads_log_id = al.id").
// 		Where("al.ads_campaign_id", campaignID).
// 		WhereIn("ael.event_name",[]any{"PURCHASE", "SUBSCRIBE"}).
// 		Get(&conversions)
		

// 	if err != nil {
// 		return ctx.Response().Json(http.StatusInternalServerError, http.Json{
// 			"error": "Failed to get total conversions: " + err.Error(),
// 		})
// 	}
	

// 	return ctx.Response().Json(http.StatusOK, http.Json{
// 		"summary": http.Json{
// 			"total_clicks":      totalClicks,
// 			"total_conversions": conversions.TotalConversions,
// 			"total_revenue":     conversions.TotalRevenue, 
// 		},
// 		"country_stats": countryStats,

// 	})
// }


func (r *ReportAdsCampaignController) ShowReportAdsCampaign(ctx http.Context) http.Response {
	campaignID := ctx.Request().Route("campaign_id")

	var (
		totalClicks int64
		countryStats []struct {
			Country string
			Count   int
			TotalRevenue float64
		}
		conversions struct {
			TotalConversions int
			TotalRevenue     float64
		}

		wg      sync.WaitGroup
		errOnce sync.Once
		mu      sync.Mutex
		errResp error
	)

	// Define helper for error handling
	setErr := func(err error, msg string) {
		errOnce.Do(func() {
			errResp = fmt.Errorf("%s: %w", msg, err)
		})
	}

	wg.Add(3)

	// Query 1: Total Clicks
	go func() {
		defer wg.Done()
		count, err := facades.Orm().Query().
			Table("ads_log_details").
			Where("ads_campaign_id", campaignID).
			Count()
		if err != nil {
			setErr(err, "Failed to get total clicks")
			return
		}
		mu.Lock()
		totalClicks = count
		mu.Unlock()
	}()

	// Query 2: Country Stats
	go func() {
		defer wg.Done()
		var stats []struct {
			Country string
			Count   int
			TotalRevenue float64
		}

		// err := facades.Orm().Query().
		// 	Table("ads_log_details").
		// 	Select("country, COUNT(id) AS count").
		// 	Where("ads_campaign_id", campaignID).
		// 	Group("country").
		// 	Get(&stats)
		err := facades.Orm().Query().
			Table("ads_logs AS al").
			Select("country, COUNT(DISTINCT al.id) AS count, SUM(ael.value_extracted) AS total_revenue").
			Join("INNER JOIN ads_log_details AS ald ON al.ads_log_detail_id = ald.id").
			Join("LEFT JOIN ads_event_logs AS ael ON ael.ads_log_id = al.id AND ael.event_name IN ('PURCHASE', 'SUBSCRIBE')").
			Where("al.ads_campaign_id", campaignID).
			Group("country").
			Get(&stats)
		if err != nil {
			setErr(err, "Failed to get country stats")
			return
		}
		mu.Lock()
		countryStats = stats
		mu.Unlock()
	}()

	// Query 3: Conversions
	go func() {
		defer wg.Done()
		var conv struct {
			TotalConversions int
			TotalRevenue     float64
		}
		

		err := facades.Orm().Query().
			Table("ads_logs AS al").
			Select("COUNT(DISTINCT ael.ads_log_id) AS total_conversions, SUM(value_extracted) AS total_revenue").
			Join("INNER JOIN ads_event_logs AS ael ON al.id = ael.ads_log_id").
			Where("al.ads_campaign_id", campaignID).
			WhereIn("ael.event_name", []any{"PURCHASE", "SUBSCRIBE"}).
			Get(&conv)

		if err != nil {
			setErr(err, "Failed to get conversions")
			return
		}
		mu.Lock()
		conversions = conv
		mu.Unlock()
	}()

	// Wait for all queries
	wg.Wait()

	if errResp != nil {
		return ctx.Response().Json(http.StatusInternalServerError, http.Json{
			"error": facades.Lang(ctx).Get("validation.internal_error"),
		})
	}


	return ctx.Response().Json(http.StatusOK, http.Json{
		"summary": http.Json{
			"total_clicks":      totalClicks,
			"total_conversions": conversions.TotalConversions,
			"total_revenue":     conversions.TotalRevenue,
		},
		"country_stats": countryStats,
	})
}





func (r *ReportAdsCampaignController) ShowReportAdsLogDetailsCampaign(ctx http.Context) http.Response {
	campaignID := ctx.Request().Route("campaign_id")
	length := ctx.Request().Query("length", "10")
	start := ctx.Request().Query("start", "0")
	draw := ctx.Request().Query("draw", "1")
	searchValue := ctx.Request().Query("search[value]", "")

	lengthInt, _ := strconv.Atoi(length)
	startInt, _ := strconv.Atoi(start)

	// Base query
	query := facades.Orm().Query().
		Table("ads_log_details").
		Where("ads_campaign_id", campaignID)

	// Reusable filter
	query = ApplyFilters(query, ctx, FilterConfig{
		LikeFields:  []string{"ip", "country"},
		ExactFields: []string{}, // add if needed
		DateField:   "created_at",
	})

	// DataTables search
	if searchValue != "" {
		query = query.Where("MATCH(ip, country, region, city, user_agent, referrer, device_type, device_name, os_name, browser_name) AGAINST (?)", searchValue)
	}

	query = query.OrderBy("id", "asc").Offset(startInt).Limit(lengthInt)

	var adsLogDetail []*models.AdsLogDetail
	if err := query.Get(&adsLogDetail); err != nil {
		return ctx.Response().Json(http.StatusInternalServerError, http.Json{
			"error": facades.Lang(ctx).Get("validation.internal_error"),
		})
	}

	// --- Count total rows
	// totalRows, _ := facades.Orm().Query().
	// 	Table("ads_log_details").
	// 	Where("ads_campaign_id", campaignID).
	// 	Count()

	// --- Count filtered
	filteredQuery := facades.Orm().Query().
		Table("ads_log_details").
		Where("ads_campaign_id", campaignID)

	filteredQuery = ApplyFilters(filteredQuery, ctx, FilterConfig{
		LikeFields:  []string{"ip", "country"},
		DateField:   "created_at",
	})
	if searchValue != "" {
		filteredQuery = filteredQuery.Where("MATCH(ip, country, region, city, user_agent, referrer, device_type, device_name, os_name, browser_name) AGAINST (?)", searchValue)
	}
	totalFiltered, _ := filteredQuery.Count()

	return ctx.Response().Json(http.StatusOK, http.Json{
		"draw":            draw,
		"recordsTotal":    totalFiltered,
		"recordsFiltered": totalFiltered,
		"data":            adsLogDetail,
	})
}


type FilterConfig struct {
	LikeFields   []string
	ExactFields  []string
	DateField    string
}

func ApplyFilters(query orm.Query, ctx http.Context, config FilterConfig) orm.Query {
	fdate := ctx.Request().Query("fdate", "")
	tdate := ctx.Request().Query("tdate", "")

	// LIKE fields
	for _, field := range config.LikeFields {
		val := ctx.Request().Query(field, "")
		if val != "" {
			query = query.Where(field+" LIKE ?", val+"%")
		}
	}

	// Exact match fields
	for _, field := range config.ExactFields {
		val := ctx.Request().Query(field, "")
		if val != "" {
			query = query.Where(field, val)
		}
	}


	// Date range
	if fdate != "" && tdate != "" {
		query = query.WhereBetween(config.DateField, fdate+" 00:00:00", tdate+" 23:59:59")
	} else if fdate != "" {
		query = query.Where(config.DateField, ">=", fdate+" 00:00:00")
	} else if tdate != "" {
		query = query.Where(config.DateField, "<=", tdate+" 23:59:59")
	}

	return query
}




func (r *ReportAdsCampaignController) ShowReportAdsFilterCampaign(ctx http.Context) http.Response {
	campaignID := ctx.Request().Route("campaign_id")
	filterType := ctx.Request().Query("type", "country")

	// Optional: date filter
	fdate := ctx.Request().Query("fdate")
	tdate := ctx.Request().Query("tdate")

	type ReportResult struct {
		GroupName    string `json:"group_name"`
		Clicks       int64  `json:"clicks"`
		UniqueClicks int64  `json:"unique_clicks"`
		Conversions  int64  `json:"conversions"`
	}

	// Allowed filter columns
	allowedColumns := map[string]string{
		"country":      "ald.country",
		"city":         "ald.city",
		"os_name":      "ald.os_name",
		"os_version":   "ald.os_version",
		"device_type":  "ald.device_type",
		"device_name":  "ald.device_name",
		"browser_name": "ald.browser_name",
		"region":       "ald.region",
		"ip":           "ald.ip",
		"event_name":   "ael.event_name",
		"date":         "CAST(DATE(ald.created_at) AS CHAR)",
		"month":        "MONTHNAME(ald.created_at)",
		"hour_of_day":  "CONCAT(DATE_FORMAT(ald.created_at, '%l:00 %p'), ' - ', DATE_FORMAT(ald.created_at, '%l:59 %p'))",
		"day_of_week":  "DAYNAME(ald.created_at)",
	}

	groupColumn, ok := allowedColumns[filterType]
	if !ok {
		return ctx.Response().Json(http.StatusBadRequest, http.Json{
			"error": facades.Lang(ctx).Get("validation.validation_failed"),
		})
	}

	// Get campaign start and end dates
	type DateRange struct {
		StartDate string `json:"start_date"`
		EndDate   string `json:"end_date"`
	}
	var dateRange DateRange
	err := facades.Orm().Query().
		Table("ads_log_details AS ald").
		Where("ald.ads_campaign_id", campaignID).
		Select("MIN(ald.created_at) AS start_date, MAX(ald.created_at) AS end_date").
		Get(&dateRange)

	if err != nil {
		return ctx.Response().Json(http.StatusInternalServerError, http.Json{
			"error":  facades.Lang(ctx).Get("validation.internal_error"),
		})
	}

	startDate := dateRange.StartDate
	endDate := dateRange.EndDate


	// Base query
	query := facades.Orm().Query().
		Table("ads_log_details AS ald").
		Join("INNER JOIN ads_logs AS al ON ald.id = al.ads_log_detail_id").
		Join("LEFT JOIN ads_event_logs AS ael ON ael.ads_log_id = al.id").
		Select(groupColumn + " AS group_name, " +
			"COUNT(DISTINCT ald.id) AS clicks, " +
			"COUNT(DISTINCT CONCAT(ald.ip, '-', ald.user_agent)) AS unique_clicks, " +
			"COUNT(DISTINCT CASE WHEN ael.event_name IN ('PURCHASE','SUBSCRIBE') THEN al.id END) AS conversions").
		Where("ald.ads_campaign_id", campaignID).
		GroupBy(groupColumn).
		OrderByDesc("clicks")

	if fdate != "" {
		query = query.Where("ald.created_at >= ?", fdate)
		startDate = fdate
	}
	if tdate != "" {
		query = query.Where("ald.created_at <= ?", tdate)
		endDate = tdate
	}


	var results []ReportResult
	if err := query.Get(&results); err != nil {
		return ctx.Response().Json(http.StatusInternalServerError, http.Json{
			"error": facades.Lang(ctx).Get("validation.internal_error"),
		})
	}


	// Return results + start/end date for frontend
	return ctx.Response().Json(http.StatusOK, http.Json{
		"data":       results,
		"start_date": startDate,
		"end_date":   endDate,
	})
}
