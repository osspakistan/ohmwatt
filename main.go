package main

import (
	"context"
	"embed"
	"log"

	"changeme/core/battery"
	"changeme/platform/linux"
	"changeme/service"

	"github.com/wailsapp/wails/v3/pkg/application"
)

//go:embed all:frontend/dist
var assets embed.FS

func init() {
	application.RegisterEvent[string]("time")
	application.RegisterEvent[*battery.State]("battery:updated")
}

func main() {
	// Initialize Platform Hardware Backend
	linuxBackend := linux.NewLinuxBackend()

	// Initialize Application Services
	batterySvc := service.NewBatteryService(linuxBackend)
	deviceSvc := service.NewDeviceService(linuxBackend)
	hibernateSvc := service.NewHibernateService(linuxBackend)
	historySvc := service.NewHistoryService()

	app := application.New(application.Options{
		Name:        "Antlion",
		Description: "Calm, precision battery monitor & power management for Linux",
		Services: []application.Service{
			application.NewService(batterySvc),
			application.NewService(deviceSvc),
			application.NewService(hibernateSvc),
			application.NewService(historySvc),
		},
		Assets: application.AssetOptions{
			Handler: application.AssetFileServerFS(assets),
		},
		Mac: application.MacOptions{
			ApplicationShouldTerminateAfterLastWindowClosed: true,
		},
	})

	// Main Application Window: Frameless with Caffeine dark background #111111
	app.Window.NewWithOptions(application.WebviewWindowOptions{
		Title:     "Antlion — Battery Management",
		Width:     1060,
		Height:    740,
		MinWidth:  860,
		MinHeight: 600,
		Frameless: true,
		BackgroundColour: application.NewRGB(17, 17, 17),
		URL: "/",
	})

	// Background Hardware Event Watcher: Zero-polling interrupt loop
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	events, err := linuxBackend.WatchEvents(ctx)
	if err == nil {
		go func() {
			for range events {
				st, err := linuxBackend.GetLiveState()
				if err == nil && st != nil {
					app.Event.Emit("battery:updated", st)
				}
			}
		}()
	}

	err = app.Run()
	if err != nil {
		log.Fatal(err)
	}
}
