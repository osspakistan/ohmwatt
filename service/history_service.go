package service

import (
	"changeme/core/store"
	"changeme/storage/sqlite"
)

// HistoryService serves charging session logs and degradation trends from SQLite.
type HistoryService struct {
	db *sqlite.DB
}

// NewHistoryService constructs a new HistoryService backed by SQLite.
func NewHistoryService(db *sqlite.DB) *HistoryService {
	return &HistoryService{db: db}
}

// GetRecentSessions returns chronological charge/discharge sessions from SQLite.
func (s *HistoryService) GetRecentSessions() ([]store.ChargeSession, error) {
	if s.db == nil {
		return []store.ChargeSession{}, nil
	}
	sessions, err := s.db.GetRecentSessions(30)
	if err != nil {
		return nil, err
	}
	if len(sessions) == 0 {
		return []store.ChargeSession{}, nil
	}
	return sessions, nil
}

// GetDailySnapshots returns daily snapshots for health degradation charts.
func (s *HistoryService) GetDailySnapshots(days int) ([]store.DailySnapshot, error) {
	if s.db == nil {
		return []store.DailySnapshot{}, nil
	}
	return s.db.GetSnapshots(days)
}
