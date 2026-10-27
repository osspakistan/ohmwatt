# Antlion Architecture Guide — Tauri v2 + Rust + Vue 3

Architectural blueprint for **Antlion Desktop Battery Monitor & Power Management Tool**.

---

## 1. High-Level Mental Model & Stack

Antlion follows a **two-tier architecture**:
1. **Frontend (Presentation & UI)**: Vue 3 (Composition API, `<script setup>`), Vite, Tailwind CSS v4, shadcn-vue primitives (`reka-ui`), Hugeicons, and `@vueuse/core`.
2. **Backend (Core Hardware Engine)**: Tauri v2 + Rust desktop daemon, SQLite local database, and isolated platform hardware backends.

```
┌────────────────────────────────────────────────────────────────────────┐
│                   Vue 3 Frontend (src/)                                │
│   Features: Dashboard, Health, History, Device, Hibernate, Alerts      │
└───────────────────────────────────▲────────────────────────────────────┘
                                    │ Tauri v2 IPC (`invoke`, `listen`)
┌───────────────────────────────────▼────────────────────────────────────┐
│                   Rust Backend (src-tauri/src/)                        │
│   commands/ ──> core/ (db, engine, models) ──> platform/ (traits)      │
└───────────────────────────────────▲────────────────────────────────────┘
                                    │
          ┌─────────────────────────┼─────────────────────────┐
          │                         │                         │
┌─────────▼──────────────┐ ┌────────▼─────────┐ ┌─────────────▼──────────┐
│ platform/linux/        │ │ platform/macos/  │ │ platform/windows/      │
│  ├── sysfs             │ │  ├── iokit       │ │  ├── win32             │
│  ├── dbus (zbus)       │ └──────────────────┘ └────────────────────────┘
│  └── distro/debian/    │
└────────────────────────┘
```

---

## 2. Directory Layout & Boundaries

```
antlion-app/
├── docs/                               <-- Architecture and technical references
├── src/                                <-- Vue 3 Presentation Layer
│   ├── assets/                         <-- SVGs, brand marks
│   ├── components/
│   │   ├── ui/                         <-- shadcn-vue primitives (button, card, dialog, tabs)
│   │   └── common/                     <-- Global shared widgets (theme-toggle, status-pill)
│   ├── features/                       <-- Feature-driven domain modules (DRY)
│   │   ├── dashboard/                  <-- F01: Fluid battery vessel, live telemetry
│   │   │   ├── components/             <-- Vessel.vue, MetricsGrid.vue
│   │   │   ├── composables/            <-- use-dashboard.ts
│   │   │   └── index.ts
│   │   ├── health/                     <-- F02, F04, F32: Health score, stress breakdown
│   │   ├── history/                    <-- F04, F28: 24h/10d charts & charge sessions log
│   │   ├── device/                     <-- F27, F29, F30: Battery, Adapter, Laptop tabs
│   │   ├── hibernate/                  <-- F26-guide: Diagnostic checkmarks & copy commands
│   │   ├── alerts/                     <-- F11-lite, F17: 90% unplug slider, thermal warnings
│   │   └── reports/                    <-- F21: Live PDF/CSV export preview
│   ├── composables/                    <-- Global reusable composables (use-battery, use-ipc)
│   ├── stores/                         <-- Reactive application stores
│   ├── types/                          <-- TypeScript IPC interfaces (matches Rust models)
│   ├── app.vue                         <-- Main layout shell with slim icon navigation rail
│   └── main.ts
│
└── src-tauri/                          <-- Rust Hardware & System Layer
    ├── Cargo.toml
    ├── tauri.conf.json
    └── src/
        ├── lib.rs                      <-- Tauri plugin setup & command registry
        ├── main.rs                     <-- Desktop binary entry point
        ├── commands/                   <-- Thin IPC handler layer
        │   ├── battery.rs              <-- get_battery_snapshot, get_health
        │   ├── device.rs               <-- get_device_info, get_adapter_info
        │   ├── hibernate.rs            <-- get_hibernation_diagnostic
        │   ├── alerts.rs               <-- set_heads_up_threshold
        │   └── mod.rs
        ├── core/                       <-- DRY: Reusable System Engine
        │   ├── models/                 <-- Serializable telemetry structs (BatterySnapshot, etc.)
        │   ├── db/                     <-- SQLite database initialization & migrations
        │   ├── engine/                 <-- 90% alarm trigger logic & anomaly detection
        │   └── mod.rs
        └── platform/                   <-- HARDWARE ISOLATION LAYER
            ├── traits.rs               <-- Universal `BatteryBackend` trait
            ├── mod.rs                  <-- Platform factory selector (`cfg(target_os)`)
            │
            ├── linux/                  <-- Linux Hardware Implementation
            │   ├── mod.rs              <-- LinuxBackend struct implementing BatteryBackend
            │   ├── sysfs.rs            <-- Direct reading of /sys/class/power_supply/BAT*
            │   ├── dbus.rs             <-- Zero-polling UPower event listener (zbus)
            │   ├── typec.rs            <-- USB-PD adapter detection
            │   └── distro/             <-- Distribution-specific system logic
            │       ├── mod.rs          <-- DistroHandler trait (Debian, Arch, Fedora)
            │       ├── debian/         <-- 🎯 DEBIAN / UBUNTU IMPLEMENTATION
            │       │   ├── mod.rs      <-- DebianHandler implementation
            │       │   ├── grub.rs     <-- /etc/default/grub parser
            │       │   ├── swap.rs     <-- filefrag swap physical offset calculation
            │       │   └── initramfs.rs<-- /etc/initramfs-tools recipes
            │       ├── arch/           <-- Future placeholder
            │       └── fedora/         <-- Future placeholder
            │
            ├── macos/                  <-- macOS Implementation (IOKit AppleSmartBattery)
            └── windows/                <-- Windows Implementation (WMI / IOCTL)
```

---

## 3. Separation of Concerns: DRY vs Isolation

### DRY (Don't Repeat Yourself)
1. **IPC Interfaces**: Shared data contracts (`BatterySnapshot`, `AdapterInfo`, `SessionLog`) are declared once in Rust (`core/models/`) and mirrored in TypeScript (`src/types/`).
2. **Business Rules**: The 90% unplug threshold logic, alarm re-arming rules, and historical database logging live entirely in Rust `core/engine/`.
3. **UI Themes & Tokens**: The `caffeine` design system tokens are defined once in CSS and consumed uniformly across all views.

### Strictly Isolated Platform Logic
1. **Linux Hardware Code**:
   - `platform/linux/sysfs.rs` reads standard Linux kernel interfaces (`/sys/class/power_supply/BAT*`).
   - `platform/linux/distro/debian/` encapsulates Debian/Ubuntu specific bootloader, GRUB resume offsets, and initramfs files. It never leaks into general Linux code.
2. **Windows & macOS Code**:
   - `platform/macos/` and `platform/windows/` reside in separate modules isolated by `#[cfg(target_os = "...")]` compilation flags.

---

## 4. Hardware Backend Contract (`platform/traits.rs`)

```rust
use crate::core::models::{
    AdapterInfo, BatteryIdentity, BatterySnapshot, HealthSnapshot, HibernationReport,
};
use anyhow::Result;

pub trait BatteryBackend: Send + Sync {
    /// Read live operational state (watts, volts, current, percentage, charging state)
    fn get_snapshot(&self) -> Result<BatterySnapshot>;

    /// Calculate capacity retention and health grade
    fn get_health(&self) -> Result<HealthSnapshot>;

    /// Extract static hardware info (maker, model, chemistry, masked serial)
    fn get_identity(&self) -> Result<BatteryIdentity>;

    /// Inspect connected AC adapter / USB-PD power delivery
    fn get_adapter(&self) -> Result<AdapterInfo>;

    /// Distro-aware hibernation diagnostic
    fn get_hibernation_diagnostic(&self) -> Result<HibernationReport>;
}
```

---

## 5. Event-Driven Performance & Zero Polling

- When running minimized or in the system tray, Antlion **does not poll**.
- On Linux, it listens to D-Bus signals from UPower over `zbus`. When the power cable is plugged/unplugged or capacity changes, the OS notifies Antlion.
- When the user opens and focuses the dashboard window, polling runs at **2 Hz (500ms)** for smooth liquid wave animations.
