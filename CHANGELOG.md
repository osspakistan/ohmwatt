# Changelog

All notable changes to OhmWatt are documented in this file.

The format is based on [Keep a Changelog](https://keepachangelog.com/en/1.1.0/),
and this project adheres to [Semantic Versioning](https://semver.org/spec/v2.0.0.html).

---

## [Unreleased]

### Planned
- Background system tray daemon to minimize the window to the tray.
- Hardware charge threshold limiter (such as an 80% stop limit) via PolicyKit daemon.
- Cross-platform hardware drivers for Windows (Win32) and macOS (IOKit).

---

## [0.1.0] - 2026-10-29

### Initial Linux Release

First production release of OhmWatt, a precision battery and power management desktop app built with Go, Wails v3, Vue 3, and SQLite.

### Added
- **Direct Linux Kernel Telemetry**: Reads battery registers (`BAT0`, `BAT1`, `AC`) directly from `/sys/class/power_supply/` (voltage, micro-amperes, wattage, capacity, cycle count, and temperature).
- **Fluid Battery Vessel Dashboard**: Animated fluid gauge showing the current battery percentage and charge state with a live charging pulse.
- **48-Slot Hourly Telemetry Strip**: Displays 48 hours of recorded battery levels to show drain rates and charging cycles over time.
- **Pure-Go SQLite Persistence**: Stores session history in `~/.local/share/ohmwatt/ohmwatt.db` using WAL mode for safe writes without database corruption across reboots. Automatically migrates legacy database if present.
- **Session History Ledger**: Records charge and discharge sessions with start and end levels, duration, average wattage, peak wattage, and peak temperature.
- **Session Detail Sheet**: Slide-out panel showing electrical forensics (volts, amps, watt-hours, and temperature ranges) for any recorded session.
- **Cycle Count and Health Tracking**: Tracks ACPI cycle counts, design capacity versus actual full capacity, and daily battery wear.
- **Thermal and Deep Discharge Tracking**: Logs cumulative hours spent above 38°C and counts deep discharge events below 15%.
- **Multi-Battery Support**: Detects dual-battery laptops (such as ThinkPads) and automatically switches to a single-battery view when no secondary battery is installed.
- **In-App 2x Retina Snapshots**: Exports clean 2x DPI PNG screenshots directly to `~/Downloads/` with screen name and timestamp.
- **Desktop Window and Theming**: Draggable header bar, smooth page transitions, and Light, Dark, and Caffeine theme options.
- **Linux Distribution Deliverables**:
  - Debian and Ubuntu `.deb` package (`wails3 task linux:create:deb`)
  - Fedora and RHEL `.rpm` package (`wails3 task linux:create:rpm`)
  - Arch Linux `.pkg.tar.zst` package (`wails3 task linux:create:aur`)
  - Universal `.tar.gz` standalone archive with rootless `install.sh` (`wails3 task linux:create:tarball`)
  - Portable `.AppImage` bundle (`wails3 task linux:create:appimage`)
- **Desktop Assets**: High-resolution icons (128x128, 512x512, and scalable vector SVG) with an XDG `.desktop` launcher.

### Changed
- Switched desktop shell from Tauri v2 to Go and Wails v3 Alpha. Memory footprint dropped to roughly 15 to 25 MB RAM, WebKit crashes stopped, and hardware access runs through direct Go code instead of IPC bridges.
- Enforced zero-mock engineering policy across all views and services. Telemetry uses only real kernel values.

### Security
- 100% offline architecture: zero outbound network requests, zero analytics SDKs, and zero third-party telemetry pings.
