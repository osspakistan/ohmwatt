# OhmWatt Developer & Agent Operational Directives

## 1. Tooling Directives
- **Search Command**: **ALWAYS USE `ripgrep` (`rg`)** for code and text searches. Never use standard `grep`, `find -exec grep`, or recursive grep commands that freeze on `node_modules`, `.git`, or build outputs.
- **Git Commit Protocol**: Commit per task (atomic, granular commits). Do not bundle unrelated changes together.

## 2. Zero Hardcoded / Mock Data Policy
- **100% Real Hardware Data**: OhmWatt is an accurate power monitor for everyday users. Never hardcode fake battery numbers, mock dates, static cycle counts (e.g., `212`), fake health values (e.g., `91%`), or synthetic DB seed sessions.
- **Live Bindings Integration**:
  - `Health` view must bind directly to `useBattery()` and `BatteryService.GetBatteryHealth()`.
  - `Device` view must bind dynamically to `DeviceService.GetDeviceSpecs()` and `DeviceService.GetAdapterInfo()`.
  - `History` view must bind to `HistoryService.GetActiveSession()` and SQLite `GetRecentSessions()`. If there are no sessions yet, show a clean empty state, never hardcoded fake cards.
  - `Reports` view must bind to live device specs and real battery state.
  - `Hour Strip` on Dashboard must record or render real hourly slots, not math approximations.
- **Hardware Polyfill**: If a machine lacks dual batteries (`BAT1`), automatically hide dual-battery UI. Never show `(mock)` labels to end users.

## 3. Technology Stack & Parity
- **Backend**: Go + Wails v3 (`core/`, `platform/`, `storage/sqlite/`, `service/`).
- **Database**: Pure-Go SQLite (`~/.local/share/ohmwatt/ohmwatt.db`) with WAL mode.
- **Frontend**: Vue 3 + TypeScript + Tailwind CSS v4 + `@fontsource/instrument-sans`.
- **Quality Gates**: Every change must pass `bun run check` (`oxfmt`, `oxlint`, `eslint`, `vue-tsc`) and `go build -o bin/ohmwatt .`.
