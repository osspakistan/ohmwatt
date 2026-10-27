package platform

import (
	"context"
	"changeme/core/battery"
	"changeme/core/hardware"
	"changeme/core/hibernate"
	"changeme/core/power"
)

// BatteryBackend is the polymorphic hardware contract implemented per operating system.
type BatteryBackend interface {
	// GetLiveState reads real-time voltage, current, percentage, and wattage.
	GetLiveState() (*battery.State, error)

	// GetBatteryHealth computes capacity retention, cycle degradation, and health grade.
	GetBatteryHealth() (*battery.Health, error)

	// GetDeviceSpecs extracts static hardware identity (laptop model, DMI, battery serial).
	GetDeviceSpecs() (*hardware.DeviceSpecs, error)

	// GetAdapterInfo inspects connected AC adapter and USB-PD negotiation.
	GetAdapterInfo() (*power.AdapterInfo, error)

	// GetHibernateDiagnostic checks kernel swap, RAM, and resume partition.
	GetHibernateDiagnostic() (*hibernate.Diagnostic, error)

	// SetChargeCeiling writes hardware charging thresholds if supported.
	SetChargeCeiling(ceiling int) (bool, string)

	// WatchEvents streams interrupt-driven power events (UPower D-Bus).
	WatchEvents(ctx context.Context) (<-chan battery.PowerEvent, error)
}
