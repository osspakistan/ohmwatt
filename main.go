package main

import (
	"embed"
	"log"
	"time"

	"github.com/wailsapp/wails/v3/pkg/application"
)

//go:embed all:frontend/dist
var assets embed.FS

func init() {
	application.RegisterEvent[string]("time")
	application.RegisterEvent[*BatteryState]("battery:updated")
}

func main() {
	batterySvc := NewBatteryService()

	app := application.New(application.Options{
		Name:        "Antlion",
		Description: "Quiet, calm battery monitor & power management for Linux",
		Services: []application.Service{
			application.NewService(batterySvc),
		},
		Assets: application.AssetOptions{
			Handler: application.AssetFileServerFS(assets),
		},
		Mac: application.MacOptions{
			ApplicationShouldTerminateAfterLastWindowClosed: true,
		},
	})

	// Main Application Window: Frameless with custom window chrome & Caffeine theme
	app.Window.NewWithOptions(application.WebviewWindowOptions{
		Title:     "Antlion — Battery Management",
		Width:     1060,
		Height:    740,
		MinWidth:  860,
		MinHeight: 600,
		Frameless: true,
		BackgroundColour: application.NewRGB(24, 20, 19),
		URL: "/",
	})

	// Telemetry ticker: Emit real-time battery readings every 2 seconds
	go func() {
		for {
			state, err := batterySvc.GetBatteryState()
			if err == nil && state != nil {
				app.Event.Emit("battery:updated", state)
			}
			time.Sleep(2 * time.Second)
		}
	}()

	err := app.Run()
	if err != nil {
		log.Fatal(err)
	}
}
