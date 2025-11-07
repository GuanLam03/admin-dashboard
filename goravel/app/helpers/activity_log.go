package helpers

import (
	"encoding/json"
	"github.com/goravel/framework/facades"
	"github.com/goravel/framework/contracts/http"
	"goravel/app/models"
)

type ActivityLogger struct {
	causerId    uint
	causerType  string
	properties  string
	url         string
	route       string
	input       string
	logName     string
	description string
}

// Entry point
func Activity() *ActivityLogger {
	return &ActivityLogger{
		properties: "",
		input:      "",
	}
}
func (a *ActivityLogger) CausedBy(ctx http.Context, id ...uint) *ActivityLogger {
	// If id is passed explicitly, use it
	if len(id) > 0 && id[0] > 0 {
		a.causerId = id[0]
		a.causerType = ""
		return a
	}

	// Otherwise try to get from authenticated user
	var user models.User
	if err := facades.Auth(ctx).User(&user); err == nil {
		a.causerId = user.ID
		a.causerType = ""
	}

	return a
}


func (a *ActivityLogger) WithProperties(props any) *ActivityLogger {
	jsonBytes, _ := json.Marshal(props)
	a.properties = string(jsonBytes)
	return a
}

func (a *ActivityLogger) OnUrl(url string) *ActivityLogger {
	a.url = url
	return a
}

func (a *ActivityLogger) OnRoute(route string) *ActivityLogger {
	a.route = route
	return a
}

func (a *ActivityLogger) WithInput(input any) *ActivityLogger {
	jsonBytes, _ := json.Marshal(input)
	a.input = string(jsonBytes)
	return a
}


func (a *ActivityLogger) InLog(name string) *ActivityLogger {
	a.logName = name
	return a
}

func (a *ActivityLogger) Log(description string) error {
	a.description = description


	log := models.ActivityLog{
		CauserId:    a.causerId,
		CauserType:  a.causerType,
		Properties:  a.properties,
		Url:         a.url,
		Route:       a.route,
		Input:       a.input,
		LogName:     a.logName,
		Description: a.description,
	}

	return facades.Orm().Query().Create(&log)
}
