package sqlite

import (
	"database/sql"
	"fmt"
	"os"
	"path/filepath"
	"time"

	"changeme/core/store"

	_ "modernc.org/sqlite"
)

// DB wraps the SQLite database connection.
type DB struct {
	conn *sql.DB
}

// NewDB initializes or opens the SQLite database at ~/.local/share/ohmwatt/ohmwatt.db.
func NewDB() (*DB, error) {
	home, err := os.UserHomeDir()
	if err != nil {
		home = "."
	}

	dataDir := filepath.Join(home, ".local", "share", "ohmwatt")
	if err := os.MkdirAll(dataDir, 0755); err != nil {
		return nil, fmt.Errorf("failed to create data dir: %w", err)
	}

	dbPath := filepath.Join(dataDir, "ohmwatt.db")

	// Backwards compatibility: migrate legacy ~/.local/share/antlion/antlion.db if present
	legacyDb := filepath.Join(home, ".local", "share", "antlion", "antlion.db")
	if _, err := os.Stat(dbPath); os.IsNotExist(err) {
		if data, err := os.ReadFile(legacyDb); err == nil {
			_ = os.WriteFile(dbPath, data, 0644)
		}
	}

	conn, err := sql.Open("sqlite", dbPath)
	if err != nil {
		return nil, fmt.Errorf("failed to open sqlite database: %w", err)
	}

	// Optimize SQLite performance for continuous background logging
	pragmas := []string{
		"PRAGMA journal_mode = WAL;",
		"PRAGMA synchronous = NORMAL;",
		"PRAGMA busy_timeout = 5000;",
	}
	for _, p := range pragmas {
		if _, err := conn.Exec(p); err != nil {
			return nil, fmt.Errorf("failed to execute pragma: %w", err)
		}
	}

	db := &DB{conn: conn}
	if err := db.migrate(); err != nil {
		return nil, fmt.Errorf("failed to run migrations: %w", err)
	}

	return db, nil
}

// migrate creates necessary tables for historical tracking.
func (db *DB) migrate() error {
	schema := `
	CREATE TABLE IF NOT EXISTS charge_sessions (
		id TEXT PRIMARY KEY,
		start_time DATETIME NOT NULL,
		end_time DATETIME NOT NULL,
		duration TEXT NOT NULL,
		start_pct INTEGER NOT NULL,
		end_pct INTEGER NOT NULL,
		avg_watts REAL NOT NULL,
		peak_temp REAL NOT NULL,
		is_hot BOOLEAN NOT NULL,
		is_charging BOOLEAN NOT NULL,
		cycle_delta TEXT DEFAULT '+0.00',
		peak_watts TEXT DEFAULT '0.0 W',
		avg_volts TEXT DEFAULT '12.0 V',
		avg_amps TEXT DEFAULT '1.0 A',
		energy_wh TEXT DEFAULT '0.0 Wh',
		charger_type TEXT DEFAULT '',
		thermal_status TEXT DEFAULT 'Nominal',
		health_at_session INTEGER DEFAULT 80
	);

	CREATE TABLE IF NOT EXISTS daily_snapshots (
		date TEXT PRIMARY KEY,
		health_percentage INTEGER NOT NULL,
		capacity_wh REAL NOT NULL,
		cycle_count INTEGER NOT NULL,
		created_at DATETIME DEFAULT CURRENT_TIMESTAMP
	);

	CREATE TABLE IF NOT EXISTS stress_metrics (
		id INTEGER PRIMARY KEY CHECK (id = 1),
		hours_high_stress REAL DEFAULT 0,
		deep_discharges INTEGER DEFAULT 0,
		last_evaluated DATETIME DEFAULT CURRENT_TIMESTAMP
	);

	INSERT OR IGNORE INTO stress_metrics (id, hours_high_stress, deep_discharges) VALUES (1, 0, 0);

	CREATE INDEX IF NOT EXISTS idx_sessions_time ON charge_sessions(start_time DESC);
	`
	if _, err := db.conn.Exec(schema); err != nil {
		return err
	}

	// Schema evolution for existing tables
	alters := []string{
		"ALTER TABLE charge_sessions ADD COLUMN cycle_delta TEXT DEFAULT '+0.00';",
		"ALTER TABLE charge_sessions ADD COLUMN peak_watts TEXT DEFAULT '0.0 W';",
		"ALTER TABLE charge_sessions ADD COLUMN avg_volts TEXT DEFAULT '12.0 V';",
		"ALTER TABLE charge_sessions ADD COLUMN avg_amps TEXT DEFAULT '1.0 A';",
		"ALTER TABLE charge_sessions ADD COLUMN energy_wh TEXT DEFAULT '0.0 Wh';",
		"ALTER TABLE charge_sessions ADD COLUMN charger_type TEXT DEFAULT '';",
		"ALTER TABLE charge_sessions ADD COLUMN thermal_status TEXT DEFAULT 'Nominal';",
		"ALTER TABLE charge_sessions ADD COLUMN health_at_session INTEGER DEFAULT 80;",
	}
	for _, alter := range alters {
		_, _ = db.conn.Exec(alter)
	}

	return nil
}

// RecordSession inserts a completed charge or discharge session into SQLite.
func (db *DB) RecordSession(session *store.ChargeSession) error {
	query := `
	INSERT OR REPLACE INTO charge_sessions 
	(id, start_time, end_time, duration, start_pct, end_pct, avg_watts, peak_temp, is_hot, is_charging,
	 cycle_delta, peak_watts, avg_volts, avg_amps, energy_wh, charger_type, thermal_status, health_at_session)
	VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?);
	`
	_, err := db.conn.Exec(
		query,
		session.ID,
		session.StartTime,
		session.EndTime,
		session.Duration,
		session.StartPct,
		session.EndPct,
		session.AvgWatts,
		session.PeakTemp,
		session.IsHot,
		session.IsCharging,
		session.CycleDelta,
		session.PeakWatts,
		session.AvgVolts,
		session.AvgAmps,
		session.EnergyWh,
		session.ChargerType,
		session.ThermalStatus,
		session.HealthAtSession,
	)
	return err
}

// GetRecentSessions returns chronological charge/discharge sessions.
func (db *DB) GetRecentSessions(limit int) ([]store.ChargeSession, error) {
	if limit <= 0 {
		limit = 20
	}
	query := `
	SELECT id, start_time, end_time, duration, start_pct, end_pct, avg_watts, peak_temp, is_hot, is_charging,
	       COALESCE(cycle_delta, '+0.00'), COALESCE(peak_watts, '0.0 W'), COALESCE(avg_volts, '12.0 V'),
	       COALESCE(avg_amps, '1.0 A'), COALESCE(energy_wh, '0.0 Wh'), COALESCE(charger_type, ''),
	       COALESCE(thermal_status, 'Nominal'), COALESCE(health_at_session, 80)
	FROM charge_sessions
	ORDER BY start_time DESC
	LIMIT ?;
	`
	rows, err := db.conn.Query(query, limit)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var sessions []store.ChargeSession
	for rows.Next() {
		var s store.ChargeSession
		var startTime, endTime time.Time
		if err := rows.Scan(
			&s.ID,
			&startTime,
			&endTime,
			&s.Duration,
			&s.StartPct,
			&s.EndPct,
			&s.AvgWatts,
			&s.PeakTemp,
			&s.IsHot,
			&s.IsCharging,
			&s.CycleDelta,
			&s.PeakWatts,
			&s.AvgVolts,
			&s.AvgAmps,
			&s.EnergyWh,
			&s.ChargerType,
			&s.ThermalStatus,
			&s.HealthAtSession,
		); err != nil {
			return nil, err
		}
		s.StartTime = startTime
		s.EndTime = endTime
		sessions = append(sessions, s)
	}

	return sessions, nil
}

// RecordDailySnapshot records capacity and cycle count for today.
func (db *DB) RecordDailySnapshot(snapshot *store.DailySnapshot) error {
	query := `
	INSERT OR REPLACE INTO daily_snapshots (date, health_percentage, capacity_wh, cycle_count)
	VALUES (?, ?, ?, ?);
	`
	_, err := db.conn.Exec(query, snapshot.Date, snapshot.HealthPercentage, snapshot.CapacityWh, snapshot.CycleCount)
	return err
}

// GetSnapshots retrieves past daily snapshots for degradation curves.
func (db *DB) GetSnapshots(days int) ([]store.DailySnapshot, error) {
	query := `
	SELECT date, health_percentage, capacity_wh, cycle_count
	FROM daily_snapshots
	ORDER BY date ASC
	LIMIT ?;
	`
	rows, err := db.conn.Query(query, days)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var snapshots []store.DailySnapshot
	for rows.Next() {
		var s store.DailySnapshot
		if err := rows.Scan(&s.Date, &s.HealthPercentage, &s.CapacityWh, &s.CycleCount); err != nil {
			return nil, err
		}
		snapshots = append(snapshots, s)
	}
	return snapshots, nil
}

// GetStressMetrics retrieves cumulative high stress hours and deep discharge counts.
func (db *DB) GetStressMetrics() (hoursHighStress float64, deepDischarges int, err error) {
	query := `SELECT hours_high_stress, deep_discharges FROM stress_metrics WHERE id = 1;`
	err = db.conn.QueryRow(query).Scan(&hoursHighStress, &deepDischarges)
	return
}

// IncrementStressHours adds elapsed high-heat hours (>38°C) to the stress metrics.
func (db *DB) IncrementStressHours(hours float64) error {
	query := `UPDATE stress_metrics SET hours_high_stress = hours_high_stress + ?, last_evaluated = CURRENT_TIMESTAMP WHERE id = 1;`
	_, err := db.conn.Exec(query, hours)
	return err
}

// IncrementDeepDischarges logs an event when the battery dropped below 15%.
func (db *DB) IncrementDeepDischarges() error {
	query := `UPDATE stress_metrics SET deep_discharges = deep_discharges + 1, last_evaluated = CURRENT_TIMESTAMP WHERE id = 1;`
	_, err := db.conn.Exec(query)
	return err
}
