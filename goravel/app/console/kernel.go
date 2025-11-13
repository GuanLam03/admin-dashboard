package console

import (
	"github.com/goravel/framework/contracts/console"
	"github.com/goravel/framework/contracts/schedule"
	"github.com/goravel/framework/facades"
	"goravel/app/console/commands"
	// "goravel/app/http/controllers/adsTracking"
)

type Kernel struct {
}

// you can manually run the task scheduling by typing "go run . artisan schedule:run"
func (kernel Kernel) Schedule() []schedule.Event {
	return []schedule.Event{
		// facades.Schedule().Call(func() {
		// 	controller := adsTracking.NewAdsTrackingCampaignPostbackController()
		// 	controller.ProcessPendingPostbacks()
		// }).EverySecond().SkipIfStillRunning(),
		facades.Schedule().Command("app:send-pending-postbacks").EverySecond().SkipIfStillRunning(),
	}
}
func (kernel Kernel) Commands() []console.Command {
	return []console.Command{
		&commands.SendPendingPostbacks{},
	}
}
