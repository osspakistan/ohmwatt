package service

import (
	"time"
	"changeme/core/store"
)

// HistoryService serves charging session logs and degradation trends.
type HistoryService struct{}

// NewHistoryService constructs a new HistoryService.
func NewHistoryService() *HistoryService {
	return &HistoryService{}
}

// GetRecentSessions returns chronological charge/discharge sessions.
func (s *HistoryService) GetRecentSessions() ([]store.ChargeSession, error) {
	now := time.Now()
	return []store.ChargeSession{
		{
			ID:         "sess-01",
			StartTime:  now.Add(-2 * time.Hour),
			EndTime:    now.Add(-40 * time.Minute),
			Duration:   "1h 20m",
			StartPct:   32,
			EndPct:     88,
			AvgWatts:   "38.5 W",
			PeakTemp:   "34.2 °C",
			IsHot:      false,
			IsCharging: true,
		},
		{
			ID:         "sess-02",
			StartTime:  now.Add(-6 * time.Hour),
			EndTime:    now.Add(-2 * time.Hour),
			Duration:   "4h 00m",
			StartPct:   91,
			EndPct:     32,
			AvgWatts:   "12.4 W",
			PeakTemp:   "31.0 °C",
			IsHot:      false,
			IsCharging: false,
		},
		{
			ID:         "sess-03",
			StartTime:  now.Add(-24 * time.Hour),
			EndTime:    now.Add(-22 * time.Hour),
			Duration:   "1h 45m",
			StartPct:   18,
			EndPct:     91,
			AvgWatts:   "42.1 W",
			PeakTemp:   "41.5 °C",
			IsHot:      true,
			IsCharging: true,
		},
	}, nil
}
