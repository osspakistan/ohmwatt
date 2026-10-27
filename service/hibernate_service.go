package service

import (
	"changeme/core/hibernate"
	"changeme/platform"
)

// HibernateService handles system swap, RAM, and hibernation checks.
type HibernateService struct {
	backend platform.BatteryBackend
}

// NewHibernateService constructs a new HibernateService.
func NewHibernateService(backend platform.BatteryBackend) *HibernateService {
	return &HibernateService{backend: backend}
}

// GetDiagnostic returns the hibernation readiness analysis and checklist.
func (s *HibernateService) GetDiagnostic() (*hibernate.Diagnostic, error) {
	return s.backend.GetHibernateDiagnostic()
}
