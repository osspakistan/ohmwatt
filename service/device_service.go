package service

import (
	"changeme/core/hardware"
	"changeme/core/power"
	"changeme/platform"
)

// DeviceService handles static hardware specifications and charger detection.
type DeviceService struct {
	backend platform.BatteryBackend
}

// NewDeviceService constructs a new DeviceService.
func NewDeviceService(backend platform.BatteryBackend) *DeviceService {
	return &DeviceService{backend: backend}
}

// GetDeviceSpecs returns static hardware identity for battery, laptop, and BIOS.
func (s *DeviceService) GetDeviceSpecs() (*hardware.DeviceSpecs, error) {
	return s.backend.GetDeviceSpecs()
}

// GetAdapterInfo returns AC adapter details and negotiated wattage.
func (s *DeviceService) GetAdapterInfo() (*power.AdapterInfo, error) {
	return s.backend.GetAdapterInfo()
}
