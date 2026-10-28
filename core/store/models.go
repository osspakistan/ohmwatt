package store

import "time"

// ChargeSession represents a logged charging or discharging cycle.
type ChargeSession struct {
	ID              string    `json:"id"`
	StartTime       time.Time `json:"startTime"`
	EndTime         time.Time `json:"endTime"`
	Duration        string    `json:"duration"`
	StartPct        int       `json:"startPct"`
	EndPct          int       `json:"endPct"`
	CycleDelta      string    `json:"cycleDelta"`      // e.g. "+0.35"
	AvgWatts        string    `json:"avgWatts"`        // e.g. "26.4 W"
	PeakWatts       string    `json:"peakWatts"`       // e.g. "45.0 W"
	AvgVolts        string    `json:"avgVolts"`        // e.g. "12.86 V"
	AvgAmps         string    `json:"avgAmps"`         // e.g. "2.05 A"
	PeakTemp        string    `json:"peakTemp"`        // e.g. "32.5 °C"
	EnergyWh        string    `json:"energyWh"`        // e.g. "18.2 Wh"
	ChargerType     string    `json:"chargerType"`     // e.g. "USB Power Delivery 65W" or "Internal Battery Rail"
	ThermalStatus   string    `json:"thermalStatus"`   // e.g. "Nominal (<38°C)"
	HealthAtSession int       `json:"healthAtSession"` // e.g. 80
	IsHot           bool      `json:"isHot"`           // true if temp exceeded 40°C
	IsCharging      bool      `json:"isCharging"`      // true if charged, false if discharged
	IsOngoing       bool      `json:"isOngoing"`       // true if currently active/running
}

// DailySnapshot represents a daily record of battery health and cycle count.
type DailySnapshot struct {
	Date             string  `json:"date"` // YYYY-MM-DD
	HealthPercentage int     `json:"healthPercentage"`
	CapacityWh       float64 `json:"capacityWh"`
	CycleCount       int     `json:"cycleCount"`
}

// Store defines persistence contracts for battery history.
type Store interface {
	RecordSession(session *ChargeSession) error
	GetRecentSessions(limit int) ([]ChargeSession, error)
	RecordDailySnapshot(snapshot *DailySnapshot) error
	GetSnapshots(days int) ([]DailySnapshot, error)
}
