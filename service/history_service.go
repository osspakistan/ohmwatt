package service

import (
	"sync"

	"changeme/core/store"
	"changeme/storage/sqlite"
)

// HistoryService serves charging session logs and degradation trends from SQLite.
type HistoryService struct {
	db            *sqlite.DB
	mu            sync.RWMutex
	activeSession *store.ChargeSession
}

// NewHistoryService constructs a new HistoryService backed by SQLite.
func NewHistoryService(db *sqlite.DB) *HistoryService {
	return &HistoryService{db: db}
}

// SetActiveSession updates the thread-safe in-progress session.
func (s *HistoryService) SetActiveSession(session *store.ChargeSession) {
	s.mu.Lock()
	defer s.mu.Unlock()
	s.activeSession = session
}

// GetActiveSession returns the live session currently underway.
func (s *HistoryService) GetActiveSession() (*store.ChargeSession, error) {
	s.mu.RLock()
	defer s.mu.RUnlock()
	if s.activeSession == nil {
		return nil, nil
	}
	cp := *s.activeSession
	return &cp, nil
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
