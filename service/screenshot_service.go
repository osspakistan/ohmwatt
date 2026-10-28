package service

import (
	"encoding/base64"
	"fmt"
	"os"
	"path/filepath"
	"strings"
	"time"
)

// ScreenshotService handles writing in-app high-resolution screenshots to disk.
type ScreenshotService struct{}

// NewScreenshotService constructs a new ScreenshotService.
func NewScreenshotService() *ScreenshotService {
	return &ScreenshotService{}
}

// SaveScreenshot writes a base64 PNG data URL to ~/.local/share/antlion/screenshots/.
func (s *ScreenshotService) SaveScreenshot(dataURL string) (string, error) {
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
	dir := filepath.Join(home, ".local", "share", "antlion", "screenshots")
	if err := os.MkdirAll(dir, 0755); err != nil {
		return "", fmt.Errorf("failed to create screenshots directory: %w", err)
	}

	filename := fmt.Sprintf("antlion-snapshot-%s.png", time.Now().Format("2006-01-02-150405"))
	outPath := filepath.Join(dir, filename)
	if err := os.WriteFile(outPath, bytes, 0644); err != nil {
		return "", fmt.Errorf("failed to write screenshot file: %w", err)
	}

	return outPath, nil
}
