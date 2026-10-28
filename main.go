package main

import (
	"context"
	"embed"
	"fmt"
	"log"
	"time"

	"changeme/core/battery"
	"changeme/core/store"
	"changeme/platform/linux"
	"changeme/service"
	"changeme/storage/sqlite"

	"github.com/wailsapp/wails/v3/pkg/application"
)

//go:embed all:frontend/dist
var assets embed.FS

func init() {
	application.RegisterEvent[string]("time")
	application.RegisterEvent[*battery.State]("battery:updated")
}

func main() {
	// Initialize Persistent SQLite Database (~/.local/share/antlion/antlion.db)
	db, err := sqlite.NewDB()
	if err != nil {
		log.Printf("Warning: failed to initialize SQLite database: %v", err)
	}

	// Initialize Platform Hardware Backend
	linuxBackend := linux.NewLinuxBackend()

	// Initial Seed: Record today's baseline snapshot if empty
	if db != nil {
		if st, err := linuxBackend.GetLiveState(); err == nil && st != nil {
			today := time.Now().Format("2006-01-02")
			_ = db.RecordDailySnapshot(&store.DailySnapshot{
				Date:             today,
				HealthPercentage: st.HealthPercentage,
				CapacityWh:       float64(st.FullCapacity) * 11.55 / 1000.0,
				CycleCount:       st.CycleCount,
			})
		}
	}

	// Initialize Application Services with SQLite persistence
	batterySvc := service.NewBatteryService(linuxBackend)
	deviceSvc := service.NewDeviceService(linuxBackend)
	hibernateSvc := service.NewHibernateService(linuxBackend)
	historySvc := service.NewHistoryService(db)

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

	// Background Hardware Event Watcher: Real Charge Session Tracking
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	go runSessionRecorder(ctx, linuxBackend, db, app)

	err = app.Run()
	if err != nil {
		log.Fatal(err)
	}
}

// runSessionRecorder logs real charging and discharging sessions into SQLite.
func runSessionRecorder(ctx context.Context, backend *linux.LinuxBackend, db *sqlite.DB, app *application.App) {
	ticker := time.NewTicker(2 * time.Second)
	defer ticker.Stop()

	type activeState struct {
		isCharging bool
		startPct   int
		startTime  time.Time
		peakTemp   float64
		wattSum    float64
		wattCount  int
	}

	var current *activeState

	for {
		select {
		case <-ctx.Done():
			return
		case <-ticker.C:
			st, err := backend.GetLiveState()
			if err != nil || st == nil {
				continue
			}

			// Broadcast live telemetry to UI
			app.Event.Emit("battery:updated", st)

			if db == nil {
				continue
			}

			// Track stress conditions
			if st.TemperatureC > 38.0 {
				_ = db.IncrementStressHours(2.0 / 3600.0) // 2 seconds in hours
			}

			// Initialize first cycle
			if current == nil {
				current = &activeState{
					isCharging: st.IsCharging,
					startPct:   st.Percentage,
					startTime:  time.Now(),
					peakTemp:   st.TemperatureC,
					wattSum:    st.PowerWatts,
					wattCount:  1,
				}
				continue
			}

			// Accumulate metrics
			current.wattSum += st.PowerWatts
			current.wattCount++
			if st.TemperatureC > current.peakTemp {
				current.peakTemp = st.TemperatureC
			}

			// Detect deep discharge crossing (<15%)
			if st.Percentage <= 15 && current.startPct > 15 {
				_ = db.IncrementDeepDischarges()
			}

			// State transition (plugged in OR unplugged) or percentage delta > 5%
			if current.isCharging != st.IsCharging {
				endTime := time.Now()
				durationMins := int(endTime.Sub(current.startTime).Minutes())
				if durationMins >= 1 { // Only log meaningful sessions (> 1 min)
					avgW := current.wattSum / float64(current.wattCount)
					durStr := fmt.Sprintf("%dm", durationMins)
					if durationMins >= 60 {
						durStr = fmt.Sprintf("%dh %02dm", durationMins/60, durationMins%60)
					}

					sess := &store.ChargeSession{
						ID:         fmt.Sprintf("sess-%d", time.Now().Unix()),
						StartTime:  current.startTime,
						EndTime:    endTime,
						Duration:   durStr,
						StartPct:   current.startPct,
						EndPct:     st.Percentage,
						AvgWatts:   fmt.Sprintf("%.1f W", avgW),
						PeakTemp:   fmt.Sprintf("%.1f °C", current.peakTemp),
						IsHot:      current.peakTemp >= 40.0,
						IsCharging: current.isCharging,
					}
					_ = db.RecordSession(sess)
				}

				// Reset for new state
				current = &activeState{
					isCharging: st.IsCharging,
					startPct:   st.Percentage,
					startTime:  time.Now(),
					peakTemp:   st.TemperatureC,
					wattSum:    st.PowerWatts,
					wattCount:  1,
				}
			}
		}
	}
}
