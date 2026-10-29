# Changelog

All notable changes to Antlion will be documented in this file.

The format is based on [Keep a Changelog](https://keepachangelog.com/en/1.1.0/),
and this project adheres to [Semantic Versioning](https://semver.org/spec/v2.0.0.html).

---

## [Unreleased]

### Planned
- Background system tray daemon minimizing window to system tray.
- Hardware charge threshold limiter (e.g. 80% stop limit) via PolicyKit daemon.
- Cross-platform hardware drivers for Windows (Win32) and macOS (IOKit).

---

## [0.1.0] - 2026-10-29

### Initial Linux Release

First production-ready Linux release of Antlion, a calm, precision battery and power management desktop app built with Go, Wails v3, Vue 3, and SQLite.

### Added
- **Direct Linux Kernel Telemetry**: Real-time reading of battery registers (`BAT0`, `BAT1`, `AC`) directly from `/sys/class/power_supply/` (voltage, micro-amperes, wattage, capacity, cycle count, temperature).
- **Physical Fluid Vessel Dashboard**: 60 FPS CSS-animated fluid level representation reflecting real state-of-charge with dynamic charging ripples.
- **48-Slot Hourly Telemetry Strip**: Real historical hourly battery levels showing charge cycles and drain rates.
- **Pure-Go SQLite Persistence**: Local-first storage (`~/.local/share/antlion/antlion.db`) using WAL mode for zero-corruption session recording across system reboots.
- **Forensic Session Ledger**: Chronological charge and discharge history logging start/end levels, elapsed duration, peak wattage, and peak temperature.
- **Session Detail Sheet**: Granular electrical slide-out inspector detailing volts, amps, energy Wh, and thermal status for every recorded session.
- **Cycle Odometer & Degradation Curve**: ACPI cycle counter calculation, true Wh wear delta, and daily capacity degradation snapshots.
- **Thermal & Deep Discharge Tracking**: Real tracking of high-heat stress hours (>38°C) and deep discharge events (<15%).
- **Multi-Battery Dynamic Polyfill**: Automatically discovers dual-battery systems (e.g. ThinkPad) and gracefully collapses to single-battery layout when secondary battery is absent.
- **In-App 2x Retina Snapshot Engine**: Built-in screenshot mechanism saving 2x DPI lossless PNG captures of any view to `~/Downloads/` with screen name and timestamp.
- **Desktop Navigation & Theming**: Frameless window with custom draggable toolbar, directional view transitions, and persistent Light / Dark / Caffeine theme modes.
- **Linux Distribution Packages**:
  - Debian / Ubuntu `.deb` package (`wails3 task linux:create:deb`)
  - Fedora / RHEL `.rpm` package (`wails3 task linux:create:rpm`)
  - Arch Linux `.pkg.tar.zst` package (`wails3 task linux:create:aur`)
  - Universal `.tar.gz` standalone archive with rootless `install.sh` (`wails3 task linux:create:tarball`)
  - Universal portable `.AppImage` bundle (`wails3 task linux:create:appimage`)
- **High-DPI Desktop Branding**: Embedded multi-resolution icons (128x128, 512x512, and scalable vector SVG) and XDG `.desktop` launcher.

### Changed
- Migrated core desktop shell architecture from Tauri v2 to Go + Wails v3 Alpha to eliminate WebKit rendering crashes, reduce memory footprint to ~15–25MB, and enable direct kernel CGO bindings.
- Established strict Zero-Mock Engineering Policy across all frontend views and backend services.

### Security
- 100% offline architecture: zero outbound network requests, zero analytics SDKs, zero third-party telemetry pings.
