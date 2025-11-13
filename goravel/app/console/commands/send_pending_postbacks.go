package commands

import (
	"github.com/goravel/framework/contracts/console"
	"github.com/goravel/framework/contracts/console/command"
	"goravel/app/http/controllers/adsTracking"
)

type SendPendingPostbacks struct {
}

// Signature The name and signature of the console command.
func (r *SendPendingPostbacks) Signature() string {
	return "app:send-pending-postbacks"
}

// Description The console command description.
func (r *SendPendingPostbacks) Description() string {
	return "Process the pending postbacks"
}

// Extend The console command extend.
func (r *SendPendingPostbacks) Extend() command.Extend {
	return command.Extend{Category: "app"}
}

// Handle Execute the console command.
func (r *SendPendingPostbacks) Handle(ctx console.Context) error {
	controller := adsTracking.NewAdsTrackingCampaignPostbackController()
	controller.ProcessPendingPostbacks()
	return nil
}
