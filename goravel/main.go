package main

import (
	"os"
	"os/signal"
	"syscall"

	"github.com/goravel/framework/facades"

	"goravel/bootstrap"
)

func main() {
	// This bootstraps the framework and gets it ready for use.
	bootstrap.Boot()

	// Create a channel to listen for OS signals
	quit := make(chan os.Signal)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)

	// Start http server by facades.Route().
	go func() {
		if err := facades.Route().Run(); err != nil {
			facades.Log().Errorf("Route Run error: %v", err)
		}
	}()

	// Start Scheduler //06/11/2025
	// you can manually run the task scheduling by typing "go run . artisan schedule:run" or uncomment below code it will auto run when "go run .""
	// go func() {
	// 	facades.Log().Info("Starting Goravel scheduler...")
	// 	facades.Schedule().Run();
	// }()


	// Listen for the OS signal
	go func() {
		<-quit

		// Shutdown scheduler //06/11/2025
		if err := facades.Schedule().Shutdown(); err != nil {
			facades.Log().Errorf("Schedule Shutdown error: %v", err)
		}

		if err := facades.Route().Shutdown(); err != nil {
			facades.Log().Errorf("Route Shutdown error: %v", err)
		}

		os.Exit(0)
	}()

	select {}
}
