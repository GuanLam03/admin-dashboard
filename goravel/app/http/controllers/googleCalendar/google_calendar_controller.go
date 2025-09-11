package googleCalendar

import (
	"context"
	"errors"
	"fmt"
	"os"
	"time"

"golang.org/x/oauth2"
"encoding/json"
	"github.com/goravel/framework/facades"
	"golang.org/x/oauth2/google"
	"google.golang.org/api/calendar/v3"
	"google.golang.org/api/option"
)

type GoogleCalendarController struct{}

func NewGoogleCalendarController() *GoogleCalendarController {
	return &GoogleCalendarController{}
}


func (g *GoogleCalendarController) getService() (*calendar.Service, string, error) {
	credentialsFile := "storage/credentials.json"
	b, err := os.ReadFile(credentialsFile)
	if err != nil {
		return nil, "", fmt.Errorf("unable to read client secret file: %w", err)
	}

	// Config with Calendar scope
	config, err := google.ConfigFromJSON(b, calendar.CalendarScope)
	if err != nil {
		return nil, "", fmt.Errorf("unable to parse client secret file: %w", err)
	}

	// Read previously saved token.json
	tokFile := "storage/token.json"
	tok, err := os.ReadFile(tokFile)
	if err != nil {
		return nil, "", fmt.Errorf("token.json not found, please run OAuth flow first")
	}

	var token oauth2.Token
	if err := json.Unmarshal(tok, &token); err != nil {
		return nil, "", fmt.Errorf("unable to parse token.json: %w", err)
	}

	client := config.Client(context.Background(), &token)
	srv, err := calendar.NewService(context.Background(), option.WithHTTPClient(client))
	if err != nil {
		return nil, "", fmt.Errorf("unable to create calendar service: %w", err)
	}

	// Get calendar ID from .env (or just use "primary")
	calendarIDAny := facades.Config().Env("GOOGLE_CALENDAR_ID", "")
	calendarID, ok := calendarIDAny.(string)
	if !ok || calendarID == "" {
		return nil, "", errors.New("invalid or missing GOOGLE_CALENDAR_ID in .env")
	}

	return srv, calendarID, nil
}


func (g *GoogleCalendarController) AddGoogleCalendar(
    title string,
    start, end time.Time,
    recurrence *string,
    attendees []string, // can be empty
) (string, error) {
    srv, calendarID, err := g.getService()
    if err != nil {
        return "", err
    }

    // Build attendees only if provided
    var eventAttendees []*calendar.EventAttendee
    if len(attendees) > 0 {
        for _, email := range attendees {
            eventAttendees = append(eventAttendees, &calendar.EventAttendee{Email: email})
        }
    }

    event := &calendar.Event{
        Summary: title,
        Start: &calendar.EventDateTime{
            DateTime: start.Format(time.RFC3339),
            TimeZone: "Asia/Kuala_Lumpur",
        },
        End: &calendar.EventDateTime{
            DateTime: end.Format(time.RFC3339),
            TimeZone: "Asia/Kuala_Lumpur",
        },
    }

    // Only add attendees if not empty
    if len(eventAttendees) > 0 {
        event.Attendees = eventAttendees
    }

    // recurrence rules
    if recurrence != nil {
        switch *recurrence {
        case "daily":
            event.Recurrence = []string{"RRULE:FREQ=DAILY"}
        case "weekly":
            event.Recurrence = []string{"RRULE:FREQ=WEEKLY"}
        case "monthly":
            event.Recurrence = []string{"RRULE:FREQ=MONTHLY"}
        }
    }

    created, err := srv.Events.Insert(calendarID, event).Do()
    if err != nil {
        return "", err
    }

    return created.Id, nil
}


// Update event
func (g *GoogleCalendarController) UpdateGoogleCalendarEvent(eventID string, title string, start, end time.Time, recurrence *string) error {
	srv, calendarID, err := g.getService()
	if err != nil {
		return err
	}

	event := &calendar.Event{
		Summary: title,
		Start: &calendar.EventDateTime{
			DateTime: start.Format(time.RFC3339),
			TimeZone: "Asia/Kuala_Lumpur",
		},
		End: &calendar.EventDateTime{
			DateTime: end.Format(time.RFC3339),
			TimeZone: "Asia/Kuala_Lumpur",
		},
	}

	if recurrence != nil {
		switch *recurrence {
		case "daily":
			event.Recurrence = []string{"RRULE:FREQ=DAILY"}
		case "weekly":
			event.Recurrence = []string{"RRULE:FREQ=WEEKLY"}
		case "monthly":
			event.Recurrence = []string{"RRULE:FREQ=MONTHLY"}
		}
	}

	_, err = srv.Events.Update(calendarID, eventID, event).Do()
	return err
}

// Delete event
func (g *GoogleCalendarController) DeleteGoogleCalendarEvent(eventID string) error {
	srv, calendarID, err := g.getService()
	if err != nil {
		return err
	}

	err = srv.Events.Delete(calendarID, eventID).Do()
	if err != nil {
		// Return the error to the caller
		return fmt.Errorf("failed to delete event with ID %s: %v", eventID, err)
	}

	return nil
}

