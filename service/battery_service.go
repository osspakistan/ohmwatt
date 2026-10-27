package service

import (
	"changeme/core/battery"
	"changeme/platform"
)

// BatteryService handles live battery state queries and threshold controls.
type BatteryService struct {
	backend platform.BatteryBackend
}

// NewBatteryService constructs a new BatteryService.
func NewBatteryService(backend platform.BatteryBackend) *BatteryService {
	return &BatteryService{backend: backend}
}

// GetBatteryState returns the live telemetry state.
func (s *BatteryService) GetBatteryState() (*battery.State, error) {
	return s.backend.GetLiveState()
}

// GetBatteryHealth returns capacity retention, cycle count, and grade.
func (s *BatteryService) GetBatteryHealth() (*battery.Health, error) {
	return s.backend.GetBatteryHealth()
}

// SetChargeCeiling sets the hardware charging limit percentage.
func (s *BatteryService) SetChargeCeiling(ceiling int) (bool, string) {
	return s.backend.SetChargeCeiling(ceiling)
}
