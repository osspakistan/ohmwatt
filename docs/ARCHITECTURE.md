# Antlion Architecture Guide — Go + Wails v3 + Vue 3

Architectural blueprint for **Antlion Desktop Battery Monitor & Power Management Tool**.

---

## 1. High-Level Mental Model & Stack

Antlion follows a **two-tier architecture**:
1. **Frontend (Presentation & UI)**: Vue 3 (Composition API, `<script setup>`), Vite, Tailwind CSS v4, shadcn-vue primitives (`reka-ui`), Hugeicons, bundled `@fontsource/instrument-sans`, and directional desktop transitions.
2. **Backend (Core Hardware Engine)**: Wails v3 + Go desktop application shell, pure-Go SQLite local database, and strictly isolated platform hardware drivers.

```
┌────────────────────────────────────────────────────────────────────────┐
│                   Vue 3 Frontend (frontend/src/)                       │
│   Features: Dashboard, Health, History, Device, Hibernate, Alerts      │
└───────────────────────────────────▲────────────────────────────────────┘
                                    │ Wails v3 RPC & Events
┌───────────────────────────────────▼────────────────────────────────────┐
│                   Go Backend (antlion-app/)                            │
│   service/ ──> core/ (models, rules, db) ──> platform/ (drivers)       │
└───────────────────────────────────▲────────────────────────────────────┘
                                    │
          ┌─────────────────────────┼─────────────────────────┐
          │                         │                         │
┌─────────▼──────────────┐ ┌────────▼─────────┐ ┌─────────────▼──────────┐
│ platform/linux/        │ │ platform/darwin/ │ │ platform/windows/      │
│  ├── sysfs             │ │  ├── iokit (cgo) │ │  ├── win32             │
│  ├── dbus (godbus)     │ └──────────────────┘ └────────────────────────┘
│  └── distro/debian/    │
└────────────────────────┘
```

---

## 2. Directory Layout & Boundaries

```
antlion-app/
├── docs/                               <-- Architecture and technical references
│   └── ARCHITECTURE.md                 <-- This document
│
├── frontend/                           <-- Vue 3 Presentation Layer (ALREADY DEVELOPED)
│   ├── bindings/                       <-- Wails 3 auto-generated TypeScript services (DO NOT EDIT)
│   ├── public/                         <-- Retina marks, app icons, favicons
│   ├── src/
│   │   ├── assets/                     <-- SVGs, vector brand marks
│   │   ├── components/
│   │   │   ├── ui/                     <-- shadcn-vue primitives (button, card, dialog, tabs)
│   │   │   └── common/                 <-- Global widgets (app-toolbar, theme-toggle, status-pill, masked-field)
│   │   ├── composables/                <-- Vue bridges (use-battery.ts, use-theme.ts)
│   │   ├── features/                   <-- Feature-driven domain modules (DRY)
│   │   │   ├── dashboard/              <-- F01: Liquid vessel SVG animation, live telemetry
│   │   │   │   └── components/         <-- liquid-vessel.vue
│   │   │   ├── health/                 <-- F02, F04, F32: Health score, capacity decay, stress breakdown
│   │   │   ├── history/                <-- F04, F28: 24h/10d charts & charge sessions log
│   │   │   ├── device/                 <-- F27, F29, F30: Battery, Adapter, Laptop DMI tabs
│   │   │   ├── hibernate/              <-- F26-guide: Diagnostic checkmarks & copy commands
│   │   │   ├── alerts/                 <-- F11-lite, F17: 90% unplug slider, thermal warnings
│   │   │   ├── reports/                <-- F21: Live printable certificate preview
│   │   │   └── settings/               <-- Preferences, unit toggles (°C/°F), mock hardware switch
│   │   ├── layouts/                    <-- app-layout.vue (slim navigation rail + directional router-view)
│   │   ├── router/                     <-- index.ts (Eager-loaded routes + kinetic depth transition)
│   │   ├── style.css                   <-- Caffeine tokens + @fontsource/instrument-sans
│   │   └── main.ts
│   ├── package.json
│   ├── tsconfig.json
│   └── vite.config.ts
│
├── core/                               <-- Pure Domain Logic (Zero UI, Zero OS ties)
│   ├── battery/                        # Battery state machine, validation & health models
│   │   ├── battery.go                  # Core Battery entity & domain logic
│   │   ├── health.go                   # Health degradation & cycle decay algorithms
│   │   └── anomalies.go                # Anomaly detector (hot charges, vampire drain)
│   ├── hardware/                       # Host identity, DMI, CPU, Memory contracts
│   │   └── specs.go                    # Device identity models (Manufacturer, Serial, BIOS)
│   ├── power/                          # Adapter, USB-PD, Charger negotiation models
│   │   └── adapter.go                  # Power delivery models & underpowered warnings
│   ├── hibernate/                      # Sleep/Hibernate diagnostic rules engine
│   │   └── readiness.go                # Swap size vs RAM calculation & risk checks
│   └── store/                          # Persistence domain interface
│       ├── store.go                    # Snapshot & Session repository contracts
│       └── models.go                   # Daily snapshot & charging session entities
│
├── platform/                           # Hardware / OS Driver Implementations
│   ├── platform.go                     # Driver factory & registry
│   ├── linux/                          # Linux Subsystem Driver
│   │   ├── sysfs.go                    # /sys/class/power_supply reader
│   │   ├── dmi.go                      # /sys/class/dmi/id host reader
│   │   ├── dbus_upower.go              # UPower D-Bus listener (interrupt-driven, 0% CPU)
│   │   ├── ucsi.go                     # USB-PD /sys/class/typec negotiation
│   │   └── distro/                     # Distribution-specific system logic
│   │       ├── distro.go               # DistroHandler interface
│   │       └── debian/                 # Debian / Ubuntu Implementation
│   │           ├── grub.go             # /etc/default/grub parser
│   │           ├── swap.go             # filefrag physical resume offset calculation
│   │           └── initramfs.go        # /etc/initramfs-tools recipes
│   ├── darwin/                         # macOS Subsystem Driver (IOKit AppleSmartBattery)
│   │   └── iokit_cgo.go
│   └── windows/                        # Windows Subsystem Driver (IOCTL / WMI)
│       └── win32.go
│
├── storage/                            # Infrastructure & Database Engine
│   ├── sqlite/                         # SQLite implementation (pure Go, zero Cgo)
│   │   ├── db.go                       # Connection & SQLite pragma setup
│   │   ├── snapshots.go                # Daily health snapshot repository
│   │   ├── sessions.go                 # Charge/discharge session logging
│   │   └── migrations/                 # Embedded SQL schema migrations
│   └── config/                         # App settings (JSON/TOML in ~/.config/antlion)
│       └── config.go
│
├── service/                            # Wails Application Services (bound to frontend TS)
│   ├── battery_service.go              # Primary telemetry RPC (Now & Health screens)
│   ├── device_service.go               # Hardware specs RPC (Device screen)
│   ├── history_service.go              # Historical data RPC (History screen charts)
│   ├── hibernate_service.go            # Diagnostics RPC (Hibernate screen checklist)
│   ├── alert_service.go                # Threshold & heat monitoring daemon
│   └── report_service.go               # Certificate generation & PDF/CSV exporter
│
├── build/                              # Packaging, icons, platform config
├── main.go                             # App bootstrap, DI container, Wails window
├── go.mod
└── go.sum
```

---

## 3. Separation of Concerns: DRY vs Isolation

### DRY (Don't Repeat Yourself)
1. **RPC Interfaces**: Go structs in `service/` generate strongly-typed TypeScript clients into `frontend/bindings/changeme/` via `wails3 generate bindings`.
2. **Business Rules**: The 90% unplug threshold logic, alarm re-arming rules, and historical database logging live entirely in Go `core/`.
3. **UI Themes & Tokens**: The `caffeine` design system tokens are defined once in `frontend/src/style.css` and consumed uniformly across all views.

### Strictly Isolated Platform Logic
1. **Linux Hardware Code**:
   - `platform/linux/sysfs.go` reads standard Linux kernel interfaces (`/sys/class/power_supply/BAT*`).
   - `platform/linux/distro/debian/` encapsulates Debian/Ubuntu specific bootloader, GRUB resume offsets, and initramfs files. It never leaks into general Linux code.
2. **Windows & macOS Code**:
   - `platform/darwin/` and `platform/windows/` reside in separate modules isolated by Go build tags (`//go:build darwin`, `//go:build windows`).

---

## 4. Hardware Backend Contract (`platform/platform.go`)

```go
package platform

import (
	"context"
	"antlion-app/core/battery"
	"antlion-app/core/hardware"
	"antlion-app/core/hibernate"
	"antlion-app/core/power"
)

type BatteryBackend interface {
	// Read live operational state (watts, volts, current, percentage, charging state)
	GetSnapshot() (*battery.Snapshot, error)

	// Calculate capacity retention and health grade
	GetHealth() (*battery.Health, error)

	// Extract static hardware info (maker, model, chemistry, masked serial, laptop DMI)
	GetDeviceSpecs() (*hardware.DeviceSpecs, error)

	// Inspect connected AC adapter / USB-PD power delivery
	GetAdapterInfo() (*power.AdapterInfo, error)

	// Distro-aware hibernation diagnostic
	GetHibernationDiagnostic() (*hibernate.Diagnostic, error)

	// Interrupt-driven power event stream (D-Bus UPower on Linux)
	WatchEvents(ctx context.Context) (<-chan battery.PowerEvent, error)
}
```

---

## 5. Event-Driven Performance & Zero Polling

- When running minimized or in the background, Antlion **does not poll**.
- On Linux, it listens to D-Bus signals from UPower over `godbus`. When the power cable is plugged/unplugged or capacity changes, the OS notifies Antlion.
- When the user opens and focuses the dashboard window, real-time telemetry updates every **2 seconds** for smooth liquid wave animations without waking the CPU unnecessarily.
