package battery

import "time"

// State represents the live telemetry readings of a battery.
type State struct {
	Percentage       int       `json:"percentage"`
	Status           string    `json:"status"` // "Charging", "Discharging", "Full", "Not charging"
	IsCharging       bool      `json:"isCharging"`
	IsAcConnected    bool      `json:"isAcConnected"`
	PowerWatts       float64   `json:"powerWatts"`
	VoltageVolts     float64   `json:"voltageVolts"`
	CurrentAmperes   float64   `json:"currentAmperes"`
	HealthPercentage int       `json:"healthPercentage"`
	CycleCount       int       `json:"cycleCount"`
	DesignCapacity   int       `json:"designCapacity"`  // mAh
	CurrentCapacity  int       `json:"currentCapacity"` // mAh
	FullCapacity     int       `json:"fullCapacity"`    // mAh
	TemperatureC     float64   `json:"temperatureC"`
	TimeToEmptyMins  int       `json:"timeToEmptyMins"`
	TimeToFullMins   int       `json:"timeToFullMins"`
	LastUpdated      time.Time `json:"lastUpdated"`
}

// Health represents the long-term degradation and wear analysis.
type Health struct {
	HealthPercentage int     `json:"healthPercentage"`
	Grade            string  `json:"grade"` // "Grade A", "Grade B", etc.
	Condition        string  `json:"condition"` // "Normal", "Replace Soon", "Damaged"
	DesignCapacityWh float64 `json:"designCapacityWh"`
	FullCapacityWh   float64 `json:"fullCapacityWh"`
	CycleCount       int     `json:"cycleCount"`
	EstimatedDate80  string  `json:"estimatedDate80"` // e.g. "Mar 2028"
	HoursAtHighStress int    `json:"hoursAtHighStress"`
	DeepDischarges    int    `json:"deepDischarges"`
}

// PowerEvent represents an interrupt-driven power change from the OS (e.g. UPower).
type PowerEvent struct {
	Type        string    `json:"type"` // "ac_plugged", "ac_unplugged", "capacity_changed"
	Percentage  int       `json:"percentage"`
	Timestamp   time.Time `json:"timestamp"`
}
