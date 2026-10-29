<div align="center">

<img src="docs/assets/ohmwatt-icon.png" alt="OhmWatt Logo" width="96" height="96" />

# OhmWatt

**A clean, accurate battery and power monitor for Linux laptops.**  
No guesswork. No cloud tracking. 100% direct hardware readings.

[![Go Version](https://img.shields.io/badge/Go-1.23+-00ADD8?style=flat-square&logo=go)](https://golang.org)
[![Wails v3](https://img.shields.io/badge/Wails-v3_Alpha-DF0000?style=flat-square)](https://v3.wails.io)
[![Vue 3](https://img.shields.io/badge/Vue-3.5+-4FC08D?style=flat-square&logo=vuedotjs)](https://vuejs.org)
[![Tailwind v4](https://img.shields.io/badge/Tailwind-v4.0-38B2AC?style=flat-square&logo=tailwind-css)](https://tailwindcss.com)
[![SQLite](https://img.shields.io/badge/SQLite-Pure--Go_WAL-003B57?style=flat-square&logo=sqlite)](https://modernc.org/sqlite)
[![Zero Telemetry](https://img.shields.io/badge/Privacy-100%25_Offline-success?style=flat-square)](#privacy--offline-guarantee)
[![License](https://img.shields.io/badge/License-MIT-blue?style=flat-square)](LICENSE)

[Visual Showcase](#visual-showcase) • [Problem & Solution](#problem--solution) • [Status Matrix](#honest-status-matrix-what-works-vs-what-does-not) • [Tech Stack](#technology-stack) • [Project Tree](#project-structure) • [Running & Building](#running--building-with-wails) • [Roadmap](#future-roadmap)

</div>

---

## Visual Showcase

OhmWatt reads hardware registers directly from Linux kernel sysfs nodes, rendering live battery levels, ACPI cycle degradation, and real-time power draw:

<div align="center">
  <img src="docs/assets/ohmwatt-dashboard-preview.png" alt="OhmWatt Dashboard Preview" width="92%" />
  <p><em>Live 2x Retina capture of OhmWatt running on Linux (HP EliteBook 850 G7) with ACPI cycle odometer and live 48-slot hourly telemetry strip.</em></p>
  <br />
  <img src="docs/assets/ohmwatt-health-preview.png" alt="OhmWatt Battery Health Preview" width="92%" />
  <p><em>Battery Health analytics showing true milliwatt-hour capacity, wear delta, high-temperature stress hours, and deep discharge counts.</em></p>
</div>

---

## Problem & Solution

### The Problem

1. **Operating System Inaccuracy**: Most desktop battery meters provide a coarse percentage and an erratic "time remaining" estimate that swings wildly between several hours and thirty minutes under brief CPU spikes.
2. **Resource Overhead**: Many third-party battery utilities run on heavy Chromium or Electron runtimes that use 300 MB to 500 MB of RAM and 5% to 10% background CPU, draining the very battery they monitor.
3. **Synthetic Calculations**: Standard battery tools often compute linear degradation curves and fail on multi-battery hardware (such as ThinkPads with internal `BAT0` and removable `BAT1`).
4. **Transient State Loss**: Most tools discard session history on exit. When the laptop reboots or wakes from sleep, previous charge history is gone.
5. **Telemetry & Tracking**: Several battery monitoring utilities bundle analytics SDKs, crash reporters, and cloud sync services that phone home with system identifiers and usage logs.

### The Solution

- **Direct Kernel Telemetry**: Reads raw voltage (`voltage_now`), micro-amperes (`current_now`), power draw (`power_now`), micro-watt-hours (`energy_now` / `charge_now`), and ACPI cycle counts directly from `/sys/class/power_supply/`.
- **100% Real Hardware Data**: Real battery numbers only. No simulated curves, no fake cycle numbers, and no placeholder historical sessions. The UI adapts automatically if your laptop has one or two batteries.
- **Persistent SQLite Ledger**: Every charge and discharge session, temperature threshold event, and daily degradation snapshot is saved locally to SQLite (`~/.local/share/ohmwatt/ohmwatt.db`) using Write-Ahead Logging (WAL).
- **Lightweight Desktop Footprint**: Built with Go and Wails v3 Alpha using the system WebKitGTK renderer. Memory footprint sits around 15 to 25 MB RAM with negligible background CPU usage.
- **In-App 2x Retina Snapshots**: Built-in screenshot engine captures clean 2x Retina PNGs of any view directly to `~/Downloads/` with screen name and timestamp.
- **100% Offline & Private**: Zero network calls. No telemetry. No analytics. No cloud accounts.

---

## Honest Status Matrix: What Works vs What Does Not

To be completely transparent about the current state of OhmWatt, here is what is **actually implemented and functional today** versus what is **under active construction or planned**:

| Capability | Status | Implementation Details |
| :--- | :---: | :--- |
| **Linux Kernel Telemetry** | **LIVE** | Full support for `BAT0`, `BAT1`, `AC` via `/sys/class/power_supply/`. |
| **Live Fluid Vessel & Dashboard** | **LIVE** | 60 FPS CSS liquid vessel reflecting real battery percentage and charge state. |
| **Hourly Telemetry Strip** | **LIVE** | 48 real hourly slots recording live battery levels and charge history. |
| **Pure-Go SQLite Storage** | **LIVE** | `modernc.org/sqlite` database with WAL mode for sessions, metrics, and snapshots. |
| **Session History Ledger** | **LIVE** | Records start/end time, duration, start/end %, avg/peak watts, peak temperature. |
| **Session Detail Sheet** | **LIVE** | Slide-out Radix sheet displaying granular electrical forensics for any session. |
| **Battery Health Analytics** | **LIVE** | Real design vs full capacity (mWh), true wear delta, ACPI cycle odometer. |
| **Thermal & Deep Discharge Tracking** | **LIVE** | Cumulative high-stress hours (>38°C) and deep discharge events (<15%). |
| **Hardware Polyfill (Multi-Battery)** | **LIVE** | Automatically hides `BAT1` vessel when machine has single battery. |
| **In-App Retina Snapshot (Method A)** | **LIVE** | Exports 2x DPI lossless PNG to `~/Downloads/` with screen name and timestamp. |
| **Theme System** | **LIVE** | Dark mode, Light mode, and Caffeine theme with localStorage persistence. |
| **Frameless Window & Navigation** | **LIVE** | Custom draggable titlebar, navigation rail, and directional page transitions. |
| **Linux Packaging & CI Automation** | **LIVE** | GitHub Actions workflow packaging `.deb`, `.rpm`, `.pkg.tar.zst`, `.tar.gz`, and `.AppImage`. |
| **System Tray Background Daemon** | **IN PROGRESS** | Window minimize-to-tray and background polling daemon are being wired up. |
| **Charge Threshold Limiter** | **PLANNED** | Setting charge limits (e.g. stop at 80%) requires PolicyKit / root daemon helper. |
| **Windows & macOS Native Drivers** | **PLANNED** | Platform interface exists; Win32 / IOKit backend implementations are upcoming. |

---

## Technology Stack

```
┌────────────────────────────────────────────────────────┐
│                   OhmWatt Frontend                     │
│  Vue 3.5 (Composition API) • TypeScript • Vite 6       │
│  Tailwind CSS v4 • Radix Vue • @fontsource/instrument  │
├────────────────────────────────────────────────────────┤
│                   IPC Transport Layer                  │
│       Wails v3 Alpha (RPC Bindings & Native Events)    │
├────────────────────────────────────────────────────────┤
│                    Go Core Backend                     │
│  BatteryService • DeviceService • HistoryService       │
│  ScreenshotService • HibernateService                  │
├──────────────────────┬─────────────────────────────────┤
│  Persistent Storage  │       Platform Drivers          │
│  Pure-Go SQLite WAL  │  Linux: /sys/class/power_supply │
│  (~/.local/share/...) │  (Windows/macOS stubs in place) │
└──────────────────────┴─────────────────────────────────┘
```

### Backend
- **Go 1.23+**: High-concurrency, memory-safe backend engine.
- **Wails v3 Alpha**: Lightweight desktop framework coupling Go with native system webviews (WebKitGTK on Linux).
- **SQLite (Pure-Go)**: `modernc.org/sqlite` database engine compiling directly without CGO dependencies.

### Frontend
- **Vue 3.5+**: Composition API (`<script setup lang="ts">`), reactive state management, KeepAlive view caching.
- **TypeScript 5.8+**: Strict type safety across all components, composables, and Wails RPC bindings.
- **Tailwind CSS v4**: Ultra-fast utility styling engine with CSS variable design tokens.
- **Radix Vue**: Accessible, headless primitives (Dialog, Sheet, Tooltip, Select, Switch).
- **Icons**: Hugeicons, Lucide Vue Next, custom OhmWatt SVG iconography.
- **Typography**: `@fontsource/instrument-sans` for clear, readable text.

### Tooling & Code Quality
- **Bun**: Fast JavaScript package manager and test runner.
- **Oxlint & Oxfmt**: High-speed Rust-based linter and formatter.
- **ESLint**: Flat configuration (`eslint.config.ts`) for strict TypeScript and Vue validation.
- **vue-tsc**: Strict Vue template type checker.

---

## Project Structure

A comprehensive nested view of the repository:

```
ohmwatt/
├── AGENTS.md                          # Engineering directives (Zero-Mock policy, quality gates)
├── README.md                          # Project documentation and architecture guide
├── Taskfile.yml                       # Wails build, dev, package, and run automation tasks
├── go.mod                             # Go module definition and dependencies
├── go.sum                             # Cryptographic checksums for Go dependencies
├── main.go                            # Application bootstrap, Wails window setup & services
│
├── core/                              # Pure domain models (zero framework dependencies)
│   ├── battery/
│   │   └── battery.go                 # BatteryState, BatteryHealth, Capacity domain structs
│   ├── hardware/
│   │   └── specs.go                   # DeviceSpecs, CPU, RAM, DMI motherboard models
│   ├── hibernate/
│   │   └── readiness.go               # Swap partition size & hibernate readiness models
│   ├── power/
│   │   └── adapter.go                 # ACAdapterInfo, wattage, charging status models
│   └── store/
│       └── models.go                  # ChargeSession, DailySnapshot, StressMetrics models
│
├── platform/                          # Hardware & OS abstraction layer
│   ├── platform.go                    # BatteryDriver and HardwareDriver interfaces
│   └── linux/
│       └── sysfs.go                   # Real Linux kernel sysfs driver (/sys/class/power_supply)
│
├── service/                           # Wails RPC services exposed to frontend
│   ├── battery_service.go             # Live battery state, health telemetry, events
│   ├── device_service.go              # DMI hardware specs, CPU cores, adapter info
│   ├── hibernate_service.go           # Linux swap file / partition hibernate check
│   ├── history_service.go             # Active session tracking, SQLite query RPCs
│   └── screenshot_service.go          # Saves 2x Retina PNGs to ~/Downloads/
│
├── storage/                           # Persistent storage engine
│   └── sqlite/
│       └── db.go                      # Pure-Go SQLite initialization, WAL mode, queries
│
├── docs/                              # Project technical documentation & assets
│   ├── ARCHITECTURE.md                # In-depth architectural blueprint
│   └── assets/
│       ├── ohmwatt-dashboard-preview.png  # Retina dashboard screenshot
│       └── ohmwatt-health-preview.png     # Retina health screen screenshot
│
├── build/                             # Wails platform packaging configurations
│   ├── config.yml                     # Wails v3 build configuration
│   ├── appicon.png                    # App master icon
│   ├── Taskfile.yml                   # Common build tasks
│   ├── linux/                         # Linux desktop integration & packaging
│   │   ├── ohmwatt.desktop            # Linux XDG desktop entry
│   │   ├── Taskfile.yml               # Linux build tasks
│   │   ├── appimage/build.sh          # AppImage generation script
│   │   └── nfpm/nfpm.yaml             # Linux packages definition (.deb, .rpm, .pkg.tar.zst)
│   ├── windows/                       # Windows MSIX, NSIS, manifest configs
│   ├── darwin/                        # macOS plist, icons, dmg configs
│   ├── android/                       # Android wrapper configs
│   └── ios/                           # iOS wrapper configs
│
└── frontend/                          # Vue 3 + TypeScript single-page application
    ├── package.json                   # Frontend scripts, dependencies, quality gates
    ├── bun.lock                       # Bun lockfile
    ├── vite.config.ts                 # Vite bundler configuration
    ├── tsconfig.json                  # TypeScript compiler options
    ├── eslint.config.ts               # ESLint flat config
    ├── .oxlintrc.json                 # Oxlint linter rules
    ├── .oxfmtrc.json                  # Oxfmt code formatter rules
    ├── index.html                     # Web entry point
    │
    ├── public/                        # Static public branding assets
    │   ├── ohmwatt-mark.svg           # Primary OhmWatt vector logo
    │   ├── ohmwatt-app-icon.svg       # Desktop app icon
    │   └── favicon.ico                # Browser / window favicon
    │
    ├── bindings/                      # Auto-generated Wails v3 TypeScript bindings
    │   └── changeme/
    │       ├── core/                  # TypeScript interfaces mirroring Go domain structs
    │       └── service/               # TypeScript RPC client wrappers for Go services
    │
    └── src/
        ├── main.ts                    # Vue application entry point
        ├── app.vue                    # Root component with page transition animations
        ├── style.css                  # Tailwind CSS v4 directives & theme variables
        │
        ├── layouts/
        │   └── app-layout.vue         # Main layout: toolbar, sidebar rail, view container
        │
        ├── router/
        │   └── index.ts               # Vue Router configuration with history stack
        │
        ├── composables/
        │   ├── use-battery.ts         # Reactive hook binding frontend to BatteryService RPC
        │   ├── use-theme.ts           # Theme manager (Dark, Light, Caffeine)
        │   └── use-screenshot.ts      # Retina screenshot trigger calling ScreenshotService
        │
        ├── components/
        │   ├── common/                # Shared application components (branding, header)
        │   └── ui/                    # Headless Radix UI components (sheet, button, etc.)
        │
        ├── features/                  # Screen views
        │   ├── dashboard/             # Dashboard screen with liquid vessel & hour strip
        │   │   ├── dashboard-view.vue
        │   │   └── components/liquid-vessel.vue
        │   ├── health/                # Battery health degradation & stress forensics
        │   │   └── health-view.vue
        │   ├── history/               # Session ledger with slide-out detail sheet
        │   │   └── history-view.vue
        │   │   └── components/session-detail-sheet.vue
        │   ├── device/                # Real DMI hardware & power adapter specs
        │   │   └── device-view.vue
        │   ├── reports/               # Power analytics & exportable reports
        │   │   └── reports-view.vue
        │   ├── hibernate/             # Linux swap readiness assessment
        │   │   └── hibernate-view.vue
        │   ├── alerts/                # Thermal & battery threshold warnings
        │   │   └── alerts-view.vue
        │   └── settings/              # Display preferences, theme, telemetry interval
        │       └── settings-view.vue
        │
        └── lib/
            └── utils.ts               # Class merging (clsx + tailwind-merge)
```

---

## Running & Building with Wails

OhmWatt is built on [Wails v3](https://v3.wails.io). Use the official `wails3` CLI commands to develop, build, and package the application:

### 1. Prerequisites

- **Linux**: Kernel 5.4+ (Ubuntu 22.04+, Debian 12+, Fedora 38+, Arch Linux)
- **Go**: 1.23 or newer (`go version`)
- **Wails v3 CLI**: Install or update the CLI:
  ```bash
  go install github.com/wailsapp/wails/v3/cmd/wails3@latest
  ```
- **Bun** (recommended) or **Node.js 20+**: (`bun --version` or `node -v`)
- **WebKitGTK Development Headers**:
  - Debian / Ubuntu:
    ```bash
    sudo apt update
    sudo apt install -y libwebkit2gtk-4.1-dev libgtk-3-dev build-essential
    ```
  - Fedora:
    ```bash
    sudo dnf install -y webkit2gtk4.1-devel gtk3-devel
    ```
  - Arch Linux:
    ```bash
    sudo pacman -S webkit2gtk-4.1 gtk3
    ```

### 2. Verify Environment

Run Wails doctor to check toolchain dependencies:

```bash
wails3 doctor
```

### 3. Install Frontend Dependencies

```bash
cd frontend
bun install
cd ..
```

### 4. Run in Development Mode

Start the application with live hot-reloading for both Go backend services and Vue 3 frontend:

```bash
wails3 dev
```

### 5. Build Production Executable

Compile the production desktop binary. This automatically regenerates TypeScript bindings, builds the Vite production bundle, and outputs the native executable to `bin/ohmwatt`:

```bash
wails3 build
```

### 6. Run the Built Binary

```bash
./bin/ohmwatt
```

*(or via Wails task runner: `wails3 task run`)*

### 7. Frontend Tooling & Commands

All frontend commands run from the `frontend/` directory using `bun` (or `npm`):

```bash
cd frontend
```

| Command | Description |
| :--- | :--- |
| `bun run check` | Runs the full pre-commit pipeline: formatting, lint autofixing, and type checking. Must pass with 0 errors. |
| `bun run fix` | Combines `format` and `lint:fix` into a single autofix pass. |
| `bun run format` | Formats all files under `src/` using ultra-fast Rust-based `oxfmt`. |
| `bun run lint` | Runs ESLint across the codebase using disk caching (`eslint . --cache`). |
| `bun run lint:fix` | Runs `oxlint . --fix` followed by `eslint . --fix --cache` to resolve stylistic and rule violations. |
| `bun run type` | Validates TypeScript and Vue 3 template types without emitting output (`vue-tsc --noEmit`). |
| `bun run dev` | Launches standalone Vite dev server with hot module replacement (HMR). |
| `bun run build` | Runs `check` quality gates and compiles the production frontend assets into `frontend/dist/`. |
| `bun run shadcn` | Adds official shadcn-vue primitives to `components/ui/` (e.g. `bun run shadcn button`). |

---

## Deployment & Release Management

OhmWatt ships automated releases across five Linux package formats via GitHub Actions on every tag push (`v*`):

| Format | Target Distribution | Installation Command |
| :--- | :--- | :--- |
| **`.deb`** | Debian 12+, Ubuntu 22.04/24.04, Pop!_OS, Mint | `sudo dpkg -i ohmwatt.deb` |
| **`.rpm`** | Fedora 38+, RHEL 9+, openSUSE | `sudo dnf install ohmwatt.rpm` |
| **`.pkg.tar.zst`** | Arch Linux, Manjaro, EndeavourOS | `sudo pacman -U ohmwatt.pkg.tar.zst` |
| **`.tar.gz`** | Any Linux Distro / Rootless | `./install.sh --user` (no root required) |
| **`.AppImage`** | Universal Linux Standalone | `chmod +x ohmwatt-x86_64.AppImage && ./ohmwatt-x86_64.AppImage` |

Pre-built binaries, cryptographic SHA256 checksums, and changelogs are published on [GitHub Releases](https://github.com/osspakistan/ohmwatt/releases).

---

## Contributing Guidelines

Contributions are welcome from anyone who wants to improve OhmWatt. All changes must follow our core guidelines:

### 1. Real Hardware Data Only
OhmWatt shows real hardware measurements. **Never hardcode fake numbers**, mock dates, static cycle counts (e.g. `212`), fake health values (e.g. `91%`), or synthetic seed sessions. If a laptop has a single battery (`BAT0`), the UI automatically hides the secondary battery vessel. Never show mock or placeholder data to users.

### 2. Strict Quality Gates
Before submitting any pull request, ensure:
1. `cd frontend && bun run check` passes with **0 warnings and 0 errors**.
2. `wails3 build` compiles cleanly with **0 errors**.

### 3. Git Commit Convention
The repository enforces atomic, conventional commits:
- `feat(scope): ...` for new features
- `fix(scope): ...` for bug fixes
- `refactor(scope): ...` for internal structural improvements
- `docs(scope): ...` for documentation updates

---

## Future Roadmap

- [x] **Phase 1: Release CI/CD Pipeline**
  - Implement GitHub Actions workflow building `.deb`, `.rpm`, `.pkg.tar.zst`, `.tar.gz`, and `.AppImage` releases.
  - SHA256 checksum generation and automated GitHub Release publication.
- [ ] **Phase 2: System Tray & Background Monitoring**
  - Minimize-to-tray background daemon using Wails v3 system tray APIs.
  - Background power loss and thermal alert notifications.
- [ ] **Phase 3: Hardware Charge Limiting**
  - PolicyKit (`pkexec`) daemon integration to safely write to vendor charge limit registers:
    - ThinkPad: `/sys/bus/platform/drivers/ideapad_acpi/.../charge_control_end_threshold`
    - ASUS: `/sys/class/power_supply/BAT0/charge_control_end_threshold`
    - Dell: `smbios-battery-ctl` integration
- [ ] **Phase 4: Cross-Platform Expansion**
  - Native Windows backend via Win32 `GetSystemPowerStatus` and `IOCTL_BATTERY_QUERY_INFORMATION`.
  - Native macOS backend via `IOKit.ps` and `IOPMPowerSource`.

---

## Privacy & Offline Guarantee

Hardware telemetry belongs exclusively to the machine owner:
- **Zero Network Calls**: OhmWatt makes zero outgoing HTTP, DNS, or socket calls.
- **Zero Analytics SDKs**: No Google Analytics, no PostHog, no Sentry, no crash collectors.
- **Local SQLite Storage**: All recorded telemetry stays inside `~/.local/share/ohmwatt/ohmwatt.db`.

---

## License

OhmWatt is open-source software licensed under the [MIT License](LICENSE).
