package main

import (
	"bufio"
	"fmt"
	"math"
	"os"
	"path/filepath"
	"strconv"
	"strings"
	"time"
)

// BatteryState represents the unified battery telemetry returned to the frontend.
type BatteryState struct {
	Percentage       int     `json:"percentage"`
	Status           string  `json:"status"` // "Charging", "Discharging", "Full", "Not charging"
	IsCharging       bool    `json:"isCharging"`
	IsAcConnected    bool    `json:"isAcConnected"`
	PowerWatts       float64 `json:"powerWatts"`
	VoltageVolts     float64 `json:"voltageVolts"`
	HealthPercentage int     `json:"healthPercentage"`
	CycleCount       int     `json:"cycleCount"`
	DesignCapacity   int     `json:"designCapacity"` // in mAh
	CurrentCapacity  int     `json:"currentCapacity"` // in mAh
	FullCapacity     int     `json:"fullCapacity"` // in mAh
	TemperatureC     float64 `json:"temperatureC"`
	TimeToEmptyMins  int     `json:"timeToEmptyMins"`
	TimeToFullMins   int     `json:"timeToFullMins"`
	Manufacturer     string  `json:"manufacturer"`
	ModelName        string  `json:"modelName"`
	SerialNumber     string  `json:"serialNumber"`
	Technology       string  `json:"technology"`
	LastUpdated      string  `json:"lastUpdated"`
}

// BatteryService provides native battery queries and operations to Wails.
type BatteryService struct{}

// NewBatteryService creates a new battery service instance.
func NewBatteryService() *BatteryService {
	return &BatteryService{}
}

// GetBatteryState reads real-time sysfs battery metrics from /sys/class/power_supply.
func (s *BatteryService) GetBatteryState() (*BatteryState, error) {
	state := &BatteryState{
		Percentage:       100,
		Status:           "Unknown",
		HealthPercentage: 100,
		TemperatureC:     29.5,
		LastUpdated:      time.Now().Format("15:04:05"),
	}

	// 1. Detect AC status
	acOnline := readSysfsInt("/sys/class/power_supply/AC/online")
	state.IsAcConnected = acOnline == 1

	// 2. Discover battery path (BAT0, BAT1, etc.)
	batPath := "/sys/class/power_supply/BAT0"
	if _, err := os.Stat(batPath); os.IsNotExist(err) {
		matches, err := filepath.Glob("/sys/class/power_supply/BAT*")
		if err == nil && len(matches) > 0 {
			batPath = matches[0]
		} else {
			// Fallback mock if running in a container or desktop without battery
			return state, nil
		}
	}

	ueventPath := filepath.Join(batPath, "uevent")
	file, err := os.Open(ueventPath)
	if err != nil {
		return state, nil
	}
	defer file.Close()

	scanner := bufio.NewScanner(file)
	props := make(map[string]string)
	for scanner.Scan() {
		line := scanner.Text()
		if idx := strings.Index(line, "="); idx != -1 {
			k := line[:idx]
			v := line[idx+1:]
			props[k] = v
		}
	}

	state.Status = props["POWER_SUPPLY_STATUS"]
	if state.Status == "" {
		state.Status = "Discharging"
	}
	state.IsCharging = state.Status == "Charging"
	state.Manufacturer = props["POWER_SUPPLY_MANUFACTURER"]
	state.ModelName = props["POWER_SUPPLY_MODEL_NAME"]
	state.SerialNumber = props["POWER_SUPPLY_SERIAL_NUMBER"]
	state.Technology = props["POWER_SUPPLY_TECHNOLOGY"]

	if capStr, ok := props["POWER_SUPPLY_CAPACITY"]; ok {
		if c, err := strconv.Atoi(capStr); err == nil {
			state.Percentage = c
		}
	}

	if cyclesStr, ok := props["POWER_SUPPLY_CYCLE_COUNT"]; ok {
		if cy, err := strconv.Atoi(cyclesStr); err == nil {
			state.CycleCount = cy
		}
	}

	// Microvolts to Volts
	voltageMicro := parseNum(props["POWER_SUPPLY_VOLTAGE_NOW"])
	if voltageMicro > 0 {
		state.VoltageVolts = float64(voltageMicro) / 1_000_000.0
	}

	// Charge design, full, now in microamp-hours
	chargeDesign := parseNum(props["POWER_SUPPLY_CHARGE_FULL_DESIGN"])
	chargeFull := parseNum(props["POWER_SUPPLY_CHARGE_FULL"])
	chargeNow := parseNum(props["POWER_SUPPLY_CHARGE_NOW"])

	// Energy fallback if microamp-hours not present (some vendors report ENERGY_*)
	if chargeDesign == 0 {
		chargeDesign = parseNum(props["POWER_SUPPLY_ENERGY_FULL_DESIGN"]) / 1000
		chargeFull = parseNum(props["POWER_SUPPLY_ENERGY_FULL"]) / 1000
		chargeNow = parseNum(props["POWER_SUPPLY_ENERGY_NOW"]) / 1000
	}

	if chargeDesign > 0 {
		state.DesignCapacity = int(chargeDesign / 1000)
	}
	if chargeFull > 0 {
		state.FullCapacity = int(chargeFull / 1000)
	}
	if chargeNow > 0 {
		state.CurrentCapacity = int(chargeNow / 1000)
	}

	// Calculate battery health percentage
	if chargeDesign > 0 && chargeFull > 0 {
		health := (float64(chargeFull) / float64(chargeDesign)) * 100.0
		state.HealthPercentage = int(math.Round(health))
		if state.HealthPercentage > 100 {
			state.HealthPercentage = 100
		}
	}

	// Current / Power draw
	currentNowMicro := parseNum(props["POWER_SUPPLY_CURRENT_NOW"])
	powerNowMicro := parseNum(props["POWER_SUPPLY_POWER_NOW"])

	if powerNowMicro > 0 {
		state.PowerWatts = float64(powerNowMicro) / 1_000_000.0
	} else if currentNowMicro > 0 && state.VoltageVolts > 0 {
		state.PowerWatts = (float64(currentNowMicro) / 1_000_000.0) * state.VoltageVolts
	}

	// Estimate time to empty / full
	if currentNowMicro > 0 {
		if state.Status == "Discharging" && chargeNow > 0 {
			hours := float64(chargeNow) / float64(currentNowMicro)
			state.TimeToEmptyMins = int(hours * 60)
		} else if state.Status == "Charging" && chargeFull > chargeNow {
			hours := float64(chargeFull-chargeNow) / float64(currentNowMicro)
			state.TimeToFullMins = int(hours * 60)
		}
	}

	return state, nil
}

// SetChargeCeiling writes a charge threshold (e.g. 80%) to sysfs if supported.
func (s *BatteryService) SetChargeCeiling(ceiling int) (bool, string) {
	if ceiling < 50 || ceiling > 100 {
		return false, "Ceiling percentage must be between 50 and 100"
	}

	thresholdPaths := []string{
		"/sys/class/power_supply/BAT0/charge_control_end_threshold",
		"/sys/class/power_supply/BAT1/charge_control_end_threshold",
	}

	for _, p := range thresholdPaths {
		if _, err := os.Stat(p); err == nil {
			err := os.WriteFile(p, []byte(strconv.Itoa(ceiling)+"\n"), 0644)
			if err != nil {
				return false, fmt.Sprintf("Permission required: run 'echo %d | sudo tee %s'", ceiling, p)
			}
			return true, fmt.Sprintf("Charge ceiling set to %d%% successfully", ceiling)
		}
	}

	return false, "Hardware charge limiters not supported by ACPI driver on this laptop."
}

func readSysfsInt(path string) int {
	data, err := os.ReadFile(path)
	if err != nil {
		return 0
	}
	val, err := strconv.Atoi(strings.TrimSpace(string(data)))
	if err != nil {
		return 0
	}
	return val
}

func parseNum(str string) int64 {
	if str == "" {
		return 0
	}
	n, _ := strconv.ParseInt(strings.TrimSpace(str), 10, 64)
	return n
}
