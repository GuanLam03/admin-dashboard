package schedules

import (
	"time"
	// "fmt"

    "goravel/app/models"
	"github.com/goravel/framework/contracts/http"
	"github.com/goravel/framework/facades"
)

type ScheduleController struct{}

func NewScheduleController() *ScheduleController {
	return &ScheduleController{}
}

// GET /schedules
func (r *ScheduleController) ShowSchedule(ctx http.Context) http.Response {
	var schedules []models.Schedule

	query := facades.Orm().Query()

	if title := ctx.Request().Query("title"); title != "" {
		query = query.Where("title LIKE ?", "%"+title+"%")
	}
	if status := ctx.Request().Query("status"); status != "" {
		query = query.Where("status", status)
	}

	if fdate := ctx.Request().Query("fdate"); fdate != "" {
		if _, err := time.Parse("2006-01-02", fdate); err == nil {
			query = query.Where("start_at >= ?", fdate)
		}
	}

	if tdate := ctx.Request().Query("tdate"); tdate != "" {
		if _, err := time.Parse("2006-01-02", tdate); err == nil {
			query = query.Where("end_at <= ?", tdate+" 23:59:59")
		}
	}

	if err := query.Order("start_at asc").Find(&schedules); err != nil {
		return ctx.Response().Json(500, http.Json{
			"message": "Failed to fetch schedules",
			"error":   err.Error(),
		})
	}

	return ctx.Response().Json(200, http.Json{
		"message": schedules,
	})
	
}

