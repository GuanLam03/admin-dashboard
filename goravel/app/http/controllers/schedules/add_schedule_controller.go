package schedules

import (
	// "fmt"
	"strconv"
	"time"

	gcal "goravel/app/http/controllers/googleCalendar"
	"goravel/app/models"

	"github.com/goravel/framework/contracts/http"
	"github.com/goravel/framework/facades"
	"goravel/app/messages"
)

type AddScheduleController struct{}

func NewAddScheduleController() *AddScheduleController {
	return &AddScheduleController{}
}

// Request struct (includes notify_roles)
type ScheduleRequest struct {
	Title       string  `json:"title"`
	Recurrence  *string `json:"recurrence"`
	StartAt     string  `json:"start_at"`
	EndAt       string  `json:"end_at"`
	Status      string  `json:"status"`
	NotifyRoles []int   `json:"notify_roles"`
}

func (r *AddScheduleController) AddSchedule(ctx http.Context) http.Response {
	var req ScheduleRequest
	if err := ctx.Request().Bind(&req); err != nil {
		return ctx.Response().Json(400, http.Json{"error": messages.GetError("validation.invalid_request")})
	}

	validator, err := facades.Validation().Make(ctx.Request().All(), models.ScheduleRules)
	if err != nil {
		return ctx.Response().Json(500, http.Json{"error": messages.GetError("validation.internal_error")})
	}
	if validator.Fails() {
		return ctx.Response().Json(422, http.Json{"errors": validator.Errors().All()})
	}

	//Parse times
	timezone := facades.Config().GetString("app.timezone")
	loc, _ := time.LoadLocation(timezone)

	startAt, err := time.ParseInLocation("2006-01-02T15:04", req.StartAt, loc)
	if err != nil {
		return ctx.Response().Json(400, http.Json{"error": messages.GetError("validation.invalid_start_at")})
	}
	endAt, err := time.ParseInLocation("2006-01-02T15:04", req.EndAt, loc)
	if err != nil {
		return ctx.Response().Json(400, http.Json{"error": messages.GetError("validation.invalid_end_at")})
	}

	// STEP 1: Collect emails by role
	var emails []string
	if len(req.NotifyRoles) > 0 {
		var casbinRules []struct {
			V0 string
			V1 string
		}
		if err := facades.Orm().Query().
			Table("casbin_rule").
			Where("ptype = ?", "g").
			WhereIn("v1", toAny(req.NotifyRoles)).
			Get(&casbinRules); err == nil {

			var userIDs []any
			for _, cr := range casbinRules {
				if uid, _ := strconv.Atoi(cr.V0); uid > 0 {
					userIDs = append(userIDs, uid)
				}
			}

			if len(userIDs) > 0 {
				var users []struct {
					Email string
				}
				_ = facades.Orm().Query().
					Table("users").
					WhereIn("id", userIDs).
					Select("email").
					Get(&users)

				for _, u := range users {
					if u.Email != "" {
						emails = append(emails, u.Email)
					}
				}
			}
		}
	}

	//STEP 2: Google Calendar
	googleCal := gcal.NewGoogleCalendarController()
	eventID, err := googleCal.AddGoogleCalendar(req.Title, startAt, endAt, req.Recurrence, emails)
	if err != nil {
		facades.Log().Errorf("Failed to insert Google Calendar event: %v", err)
		return ctx.Response().Json(500, http.Json{"error": messages.GetError("validation.google_insert_failed")})
	}

	// STEP 3: Save in DB
	schedule := models.Schedule{
		Title:         req.Title,
		Recurrence:    req.Recurrence,
		StartAt:       startAt,
		EndAt:         endAt,
		Status:        req.Status,
		GoogleEventID: &eventID,
	}

	if err := facades.Orm().Query().Create(&schedule); err != nil {
		_ = googleCal.DeleteGoogleCalendarEvent(eventID) // rollback
		return ctx.Response().Json(500, http.Json{"error": messages.GetError("validation.schedule_create_failed")})
	}

	return ctx.Response().Json(201, http.Json{
		"message": messages.GetSuccess("schedule_created"),
		"data":    schedule,
	})
}

func toAny[T any](arr []T) []any {
	out := make([]any, len(arr))
	for i, v := range arr {
		out[i] = v
	}
	return out
}
