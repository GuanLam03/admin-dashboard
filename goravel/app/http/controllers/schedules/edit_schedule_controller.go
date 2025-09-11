package schedules

import (

	"time"
	"strconv"
    "goravel/app/models"
	"github.com/goravel/framework/contracts/http"
	"github.com/goravel/framework/facades"
	gcal "goravel/app/http/controllers/googleCalendar"

)

type EditScheduleController struct{}

func NewEditScheduleController() *EditScheduleController {
	return &EditScheduleController{}
}

// GET /schedules/:id
func (c *EditScheduleController) ShowSchedule(ctx http.Context) http.Response {
	idStr := ctx.Request().Route("id")
	id, err := strconv.Atoi(idStr)
	if err != nil {
		return ctx.Response().Json(400, http.Json{
			"error": "Invalid schedule ID",
		})
	}

	var schedule models.Schedule
	if err := facades.Orm().Query().Where("id", id).First(&schedule); err != nil {
		return ctx.Response().Json(404, http.Json{
			"error": "Schedule not found",
		})
	}

	return ctx.Response().Json(200, http.Json{
		"schedule": schedule,
	})
}


func (c *EditScheduleController) EditSchedule(ctx http.Context) http.Response {
	
	googleCal := gcal.NewGoogleCalendarController()
	var schedule models.Schedule
	idStr := ctx.Request().Route("id")
	id, err := strconv.Atoi(idStr)
	

	if err != nil {
		return ctx.Response().Json(400, http.Json{
			"error": "Invalid schedule ID",
		})
	}

	var input struct {
		Title      string `json:"title"`
		Recurrence string `json:"recurrence"`
		StartAt    string `json:"start_at"`
		EndAt      string `json:"end_at"`
		Status     string `json:"status"`
	}

	if err := ctx.Request().Bind(&input); err != nil {
		return ctx.Response().Json(400, http.Json{
			"error": "Invalid request body",
		})
	}

	
	if err := facades.Orm().Query().Where("id", id).First(&schedule); err != nil {
		return ctx.Response().Json(404, http.Json{
			"error": "Schedule not found",
		})
	}

	// Parse datetime-local format ("2025-09-04T10:30")
	timezone := facades.Config().GetString("app.timezone")
	loc, _ := time.LoadLocation(timezone)

	startAt, err := time.ParseInLocation("2006-01-02T15:04", input.StartAt, loc)
	if err != nil {
		return ctx.Response().Json(400, http.Json{
			"error": "Invalid start_at format (use: YYYY-MM-DDTHH:MM)",
		})
	}

	endAt, err := time.ParseInLocation("2006-01-02T15:04", input.EndAt, loc)
	if err != nil {
		return ctx.Response().Json(400, http.Json{
			"error": "Invalid end_at format (use: YYYY-MM-DDTHH:MM)",
		})
	}

	// Update fields
	schedule.Title = input.Title
	schedule.Recurrence = &input.Recurrence
	schedule.StartAt = startAt
	schedule.EndAt = endAt
	schedule.Status = input.Status


	 // Google Calendar sync logic
    if input.Status == "active" {
        if schedule.GoogleEventID == nil {
            // create new Google Calendar event
            eventID, err := googleCal.AddGoogleCalendar(schedule.Title, startAt, endAt, schedule.Recurrence, []string{})
            if err != nil {
                return ctx.Response().Json(500, map[string]string{"error": err.Error()})
            }
            schedule.GoogleEventID = &eventID
        } else {
            // update existing event
            err := googleCal.UpdateGoogleCalendarEvent(*schedule.GoogleEventID, schedule.Title, schedule.StartAt, schedule.EndAt, schedule.Recurrence)
            if err != nil {
                return ctx.Response().Json(500, map[string]string{"error": err.Error()})
            }
        }
    } else if input.Status == "inactive" {
        if schedule.GoogleEventID != nil {
            err = googleCal.DeleteGoogleCalendarEvent(*schedule.GoogleEventID) // delete from Google Calendar
			if err != nil {
                return ctx.Response().Json(500, map[string]string{"error": err.Error()})
            }
            schedule.GoogleEventID = nil
        }
    } else if input.Status == "removed" {
        if schedule.GoogleEventID != nil {
            err = googleCal.DeleteGoogleCalendarEvent(*schedule.GoogleEventID) // permanent remove from Google Calendar
			if err != nil {
                return ctx.Response().Json(500, map[string]string{"error": err.Error()})
            }
            schedule.GoogleEventID = nil
        }
        facades.Orm().Query().Delete(&schedule)
        return ctx.Response().Json(200, map[string]string{"message": "Schedule removed"})
    }

	if err := facades.Orm().Query().Save(&schedule); err != nil {
		return ctx.Response().Json(500, http.Json{
			"error": "Failed to update schedule",
		})
	}

	return ctx.Response().Json(200, http.Json{
		"message":  "Schedule updated successfully",
		"schedule": schedule,
	})
}
