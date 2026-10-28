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
	screenshotSvc := service.NewScreenshotService()

	app := application.New(application.Options{
		Name:        "Antlion",
		Description: "Calm, precision battery monitor & power management for Linux",
		Services: []application.Service{
			application.NewService(batterySvc),
			application.NewService(deviceSvc),
			application.NewService(hibernateSvc),
			application.NewService(historySvc),
			application.NewService(screenshotSvc),
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

	go runSessionRecorder(ctx, linuxBackend, db, historySvc, app)

	err = app.Run()
	if err != nil {
		log.Fatal(err)
	}
}

// runSessionRecorder logs real charging and discharging sessions into SQLite.
func runSessionRecorder(ctx context.Context, backend *linux.LinuxBackend, db *sqlite.DB, historySvc *service.HistoryService, app *application.App) {
	ticker := time.NewTicker(2 * time.Second)
	defer ticker.Stop()

	type activeState struct {
		isCharging bool
		startPct   int
		startTime  time.Time
		peakTemp   float64
		peakWatts  float64
		wattSum    float64
		voltSum    float64
		ampSum     float64
		sampleCt   int
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

			// Track stress conditions
			if db != nil && st.TemperatureC > 38.0 {
				_ = db.IncrementStressHours(2.0 / 3600.0) // 2 seconds in hours
			}

			// Initialize first cycle
			if current == nil {
				current = &activeState{
					isCharging: st.IsCharging,
					startPct:   st.Percentage,
					startTime:  time.Now(),
					peakTemp:   st.TemperatureC,
					peakWatts:  st.PowerWatts,
					wattSum:    st.PowerWatts,
					voltSum:    st.VoltageVolts,
					ampSum:     st.CurrentAmperes,
					sampleCt:   1,
				}
			} else {
				// Accumulate metrics
				current.wattSum += st.PowerWatts
				current.voltSum += st.VoltageVolts
				current.ampSum += st.CurrentAmperes
				current.sampleCt++
				if st.TemperatureC > current.peakTemp {
					current.peakTemp = st.TemperatureC
				}
				if st.PowerWatts > current.peakWatts {
					current.peakWatts = st.PowerWatts
				}
			}

			// Keep live ongoing session object updated
			if historySvc != nil && current != nil {
				elapsed := time.Since(current.startTime)
				durStr := fmt.Sprintf("%dm", int(elapsed.Minutes()))
				if elapsed.Minutes() < 1 {
					durStr = "< 1m"
				} else if elapsed.Hours() >= 1 {
					durStr = fmt.Sprintf("%dh %02dm", int(elapsed.Hours()), int(elapsed.Minutes())%60)
				}

				avgW := current.wattSum / float64(current.sampleCt)
				avgV := current.voltSum / float64(current.sampleCt)
				avgA := current.ampSum / float64(current.sampleCt)
				deltaPct := st.Percentage - current.startPct
				cycleDeltaStr := fmt.Sprintf("%+.2f", float64(deltaPct)/100.0)
				energyWh := fmt.Sprintf("%.1f Wh", avgW*elapsed.Hours())

				chargerType := "Internal Battery Rail"
				if st.IsCharging {
					chargerType = "USB-PD 65W AC Adapter"
				}

				thermalStatus := "Nominal (<38°C)"
				if current.peakTemp >= 40.0 {
					thermalStatus = "Thermal Stress (>40°C)"
				} else if current.peakTemp >= 38.0 {
					thermalStatus = "Elevated (38-40°C)"
				}

				liveSess := &store.ChargeSession{
					ID:              "live-active-session",
					StartTime:       current.startTime,
					EndTime:         time.Now(),
					Duration:        durStr,
					StartPct:        current.startPct,
					EndPct:          st.Percentage,
					CycleDelta:      cycleDeltaStr,
					AvgWatts:        fmt.Sprintf("%.1f W", avgW),
					PeakWatts:       fmt.Sprintf("%.1f W", current.peakWatts),
					AvgVolts:        fmt.Sprintf("%.2f V", avgV),
					AvgAmps:         fmt.Sprintf("%.2f A", avgA),
					PeakTemp:        fmt.Sprintf("%.1f °C", current.peakTemp),
					EnergyWh:        energyWh,
					ChargerType:     chargerType,
					ThermalStatus:   thermalStatus,
					HealthAtSession: st.HealthPercentage,
					IsHot:           current.peakTemp >= 40.0,
					IsCharging:      current.isCharging,
					IsOngoing:       true,
				}
				historySvc.SetActiveSession(liveSess)
			}

			// Detect deep discharge crossing (<15%)
			if db != nil && current != nil && st.Percentage <= 15 && current.startPct > 15 {
				_ = db.IncrementDeepDischarges()
			}

			// State transition (plugged in OR unplugged)
			if current != nil && current.isCharging != st.IsCharging {
				endTime := time.Now()
				durationMins := int(endTime.Sub(current.startTime).Minutes())
				if db != nil && durationMins >= 1 { // Only persist sessions lasting >= 1 min
					avgW := current.wattSum / float64(current.sampleCt)
					avgV := current.voltSum / float64(current.sampleCt)
					avgA := current.ampSum / float64(current.sampleCt)
					durStr := fmt.Sprintf("%dm", durationMins)
					if durationMins >= 60 {
						durStr = fmt.Sprintf("%dh %02dm", durationMins/60, durationMins%60)
					}

					deltaPct := st.Percentage - current.startPct
					cycleDeltaStr := fmt.Sprintf("%+.2f", float64(deltaPct)/100.0)
					hours := endTime.Sub(current.startTime).Hours()
					energyWh := fmt.Sprintf("%.1f Wh", avgW*hours)

					chargerType := "Internal Battery Rail"
					if current.isCharging {
						chargerType = "USB-PD 65W AC Adapter"
					}

					thermalStatus := "Nominal (<38°C)"
					if current.peakTemp >= 40.0 {
						thermalStatus = "Thermal Stress (>40°C)"
					} else if current.peakTemp >= 38.0 {
						thermalStatus = "Elevated (38-40°C)"
					}

					sess := &store.ChargeSession{
						ID:              fmt.Sprintf("sess-%d", time.Now().Unix()),
						StartTime:       current.startTime,
						EndTime:         endTime,
						Duration:        durStr,
						StartPct:        current.startPct,
						EndPct:          st.Percentage,
						CycleDelta:      cycleDeltaStr,
						AvgWatts:        fmt.Sprintf("%.1f W", avgW),
						PeakWatts:       fmt.Sprintf("%.1f W", current.peakWatts),
						AvgVolts:        fmt.Sprintf("%.2f V", avgV),
						AvgAmps:         fmt.Sprintf("%.2f A", avgA),
						PeakTemp:        fmt.Sprintf("%.1f °C", current.peakTemp),
						EnergyWh:        energyWh,
						ChargerType:     chargerType,
						ThermalStatus:   thermalStatus,
						HealthAtSession: st.HealthPercentage,
						IsHot:           current.peakTemp >= 40.0,
						IsCharging:      current.isCharging,
						IsOngoing:       false,
					}
					_ = db.RecordSession(sess)
				}

				// Reset for new state
				current = &activeState{
					isCharging: st.IsCharging,
					startPct:   st.Percentage,
					startTime:  time.Now(),
					peakTemp:   st.TemperatureC,
					peakWatts:  st.PowerWatts,
					wattSum:    st.PowerWatts,
					voltSum:    st.VoltageVolts,
					ampSum:     st.CurrentAmperes,
					sampleCt:   1,
				}
			}
		}
	}
}
