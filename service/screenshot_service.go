package service

import (
	"encoding/base64"
	"fmt"
	"os"
	"path/filepath"
	"regexp"
	"strings"
	"time"
)

// ScreenshotService handles writing in-app high-resolution screenshots to disk.
type ScreenshotService struct{}

// NewScreenshotService constructs a new ScreenshotService.
func NewScreenshotService() *ScreenshotService {
	return &ScreenshotService{}
}

// detectCountryCode resolves the system 2-letter country code (e.g. PK, US, GB).
func detectCountryCode() string {
	// 1. Inspect /etc/localtime symlink target
	if target, err := os.Readlink("/etc/localtime"); err == nil {
		if strings.Contains(target, "Karachi") {
			return "PK"
		}
		parts := strings.Split(target, "zoneinfo/")
		if len(parts) > 1 {
			tz := parts[1]
			switch {
			case strings.HasPrefix(tz, "Asia/Karachi"):
				return "PK"
			case strings.HasPrefix(tz, "America/New_York"), strings.HasPrefix(tz, "America/Los_Angeles"), strings.HasPrefix(tz, "America/Chicago"):
				return "US"
			case strings.HasPrefix(tz, "Europe/London"):
				return "GB"
			case strings.HasPrefix(tz, "Europe/Berlin"), strings.HasPrefix(tz, "Europe/Paris"):
				return "EU"
			case strings.HasPrefix(tz, "Asia/Dubai"):
				return "AE"
			case strings.HasPrefix(tz, "Asia/Tokyo"):
				return "JP"
			}
		}
	}

	// 2. Fallback to locale country
	for _, env := range []string{"LC_ALL", "LC_TIME", "LANG"} {
		v := os.Getenv(env)
		if strings.Contains(v, "_") {
			parts := strings.Split(v, "_")
			if len(parts) > 1 {
				sub := strings.Split(parts[1], ".")[0]
				if len(sub) == 2 {
					return strings.ToUpper(sub)
				}
			}
		}
	}

	// 3. Fallback to timezone abbreviation
	zone, _ := time.Now().Zone()
	if strings.HasPrefix(zone, "PK") {
		return "PK"
	}

	return "PK"
}

// SaveScreenshot writes a base64 PNG data URL to ~/Downloads/
// with the format: ohmwatt-[screenName]-[YYYY-MM-DD-hh-mm-ss-AM/PM]-[CountryCode].png
func (s *ScreenshotService) SaveScreenshot(screenName string, dataURL string) (string, error) {
	raw := dataURL
	if idx := strings.Index(dataURL, ","); idx != -1 {
		raw = dataURL[idx+1:]
	}

	bytes, err := base64.StdEncoding.DecodeString(raw)
	if err != nil {
		return "", fmt.Errorf("failed to decode base64 screenshot: %w", err)
	}

	home, err := os.UserHomeDir()
	if err != nil {
		home = "."
	}

	// Destination directory: ~/Downloads
	dir := filepath.Join(home, "Downloads")
	if err := os.MkdirAll(dir, 0755); err != nil {
		dir = filepath.Join(home, ".local", "share", "ohmwatt", "screenshots")
		_ = os.MkdirAll(dir, 0755)
	}

	// Normalize screen name
	cleanScreen := strings.ToLower(strings.TrimSpace(screenName))
	reg := regexp.MustCompile("[^a-z0-9_-]+")
	cleanScreen = reg.ReplaceAllString(cleanScreen, "")
	if cleanScreen == "" {
		cleanScreen = "dashboard"
	}

	now := time.Now()
	timestamp := now.Format("2006-01-02-03-04-05-PM")
	countryCode := detectCountryCode()

	filename := fmt.Sprintf("ohmwatt-%s-%s-%s.png", cleanScreen, timestamp, countryCode)
	outPath := filepath.Join(dir, filename)

	if err := os.WriteFile(outPath, bytes, 0644); err != nil {
		return "", fmt.Errorf("failed to write screenshot file: %w", err)
	}

	// Also mirror to ~/.local/share/ohmwatt/screenshots/ for historical persistence
	archiveDir := filepath.Join(home, ".local", "share", "ohmwatt", "screenshots")
	if err := os.MkdirAll(archiveDir, 0755); err == nil {
		_ = os.WriteFile(filepath.Join(archiveDir, filename), bytes, 0644)
	}

	return outPath, nil
}
