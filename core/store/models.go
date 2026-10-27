package store

import "time"

// ChargeSession represents a logged charging or discharging cycle.
type ChargeSession struct {
	ID         string    `json:"id"`
	StartTime  time.Time `json:"startTime"`
	EndTime    time.Time `json:"endTime"`
	Duration   string    `json:"duration"`
	StartPct   int       `json:"startPct"`
	EndPct     int       `json:"endPct"`
	AvgWatts   string    `json:"avgWatts"`
	PeakTemp   string    `json:"peakTemp"`
	IsHot      bool      `json:"isHot"`      // true if temp exceeded 40°C
	IsCharging bool      `json:"isCharging"` // true if charged, false if discharged
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
