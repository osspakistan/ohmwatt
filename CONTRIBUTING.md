# Contributing to OhmWatt

OhmWatt is a clean, reliable battery monitor for Linux laptops, built with Go, Wails v3, Vue 3, and SQLite. Contributions that improve battery accuracy, hardware support, UI polish, and performance are welcome.

---

## 1. Core Guidelines

All contributions must follow three principles:

### A. Real Hardware Numbers Only
OhmWatt is made for real laptop users, not a demo with fake data.
- **Never commit fake numbers**: Do not hardcode static battery percentages (e.g. `91%`), synthetic cycle counts (e.g. `212`), dummy wattage, or simulated charge curves.
- **Never mock hardware registers**: All battery data must come from real kernel sysfs entries (`/sys/class/power_supply/`), ACPI registers, or genuine platform APIs.
- **Graceful hardware adaptation**: If a laptop has a single battery (`BAT0`), automatically collapse and hide dual-battery vessels (`BAT1`). Never render placeholder `(mock)` labels to end users.
- **No synthetic seed databases**: Never commit pre-filled test databases or fake historical sessions.

### B. Local-First & Zero Telemetry
- OhmWatt makes **zero outbound network calls**.
- Never introduce analytics SDKs, tracking pixels, crash collectors, or third-party telemetry services.
- All session logs and health metrics must remain strictly in the local SQLite database (`~/.local/share/ohmwatt/ohmwatt.db`).

### C. Performance & Resource Footprint
- Keep the runtime memory footprint under **~25 MB** RAM.
- Background CPU consumption must remain near **0%** during idle monitoring.

---

## 2. Development Setup

### Prerequisites
- **Linux**: Kernel 5.4+ (Ubuntu 22.04+, Debian 12+, Fedora 38+, Arch Linux)
- **Go**: Version 1.23 or newer (`go version`)
- **Bun**: Latest version (`bun --version`) or Node.js 20+
- **Wails v3 CLI**:
  ```bash
  go install github.com/wailsapp/wails/v3/cmd/wails3@latest
  ```
- **WebKitGTK & GTK Development Headers**:
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

### Initial Setup
```bash
git clone https://github.com/osspakistan/ohmwatt.git
cd ohmwatt

# Install frontend dependencies
cd frontend
bun install
cd ..
```

### Running in Development Mode
Start the live development environment with hot-reloading for both Go backend services and Vue 3 frontend components:

```bash
wails3 dev
```

---

## 3. Strict Quality Gates

Prior to committing or submitting a pull request, all automated quality gates must pass with **0 warnings and 0 errors**:

### 1. Frontend Quality Gate
Runs code formatting (`oxfmt`), lint autofixes (`oxlint` + `eslint`), and strict template type validation (`vue-tsc --noEmit`):

```bash
cd frontend
bun run check
cd ..
```

### 2. Backend Compilation Gate
Verifies Go compilation, CGO bindings, and Wails desktop packaging:

```bash
wails3 build
```

Ensure the compiled executable generates cleanly in `bin/ohmwatt`.

---

## 4. Git Commit Convention

Commit history must remain clean, atomic, and informative. Follow the [Conventional Commits](https://www.conventionalcommits.org/) specification:

```
<type>(<scope>): <short imperative summary>
```

### Allowed Types
- `feat`: A new user-facing feature or hardware driver capability.
- `fix`: A bug fix or driver correction.
- `refactor`: Code restructuring with no behavioral changes.
- `perf`: Performance optimization reducing memory or CPU usage.
- `docs`: Documentation, README, or architecture updates.
- `chore`: Packaging scripts, Taskfile updates, or dependency bumps.

### Rules
- One logical change per commit. Avoid monolithic commits that combine refactoring, documentation, and feature code.
- Write commit summaries in lowercase imperative mood (e.g. `feat(sysfs): read battery cycle count from sysfs` instead of `Added cycle counts`).

---

## 5. Pull Request Protocol

1. Fork the repository and create a descriptive branch:
   ```bash
   git checkout -b feat/charge-threshold-limiter
   ```
2. Verify hardware behavior against real Linux sysfs nodes.
3. Ensure both `bun run check` and `wails3 build` succeed locally with 0 errors.
4. Keep git history linear. Rebase onto `origin/main` before opening a pull request:
   ```bash
   git fetch origin
   git rebase origin/main
   ```
5. Submit the pull request to `osspakistan/ohmwatt` with clear test notes detailing the Linux distribution, kernel version, and laptop model used for verification.

---

## 6. Release Management Guide

OhmWatt release cycles adhere to [Semantic Versioning 2.0.0](https://semver.org/) (`vMAJOR.MINOR.PATCH`) and [Keep a Changelog 1.1.0](https://keepachangelog.com/).

### 6.1 Versioning Rules
- **MAJOR (`X.0.0`)**: Incompatible architectural overhauls, database schema migrations without backward compatibility, or breaking driver interface shifts.
- **MINOR (`0.X.0`)**: New features (such as background system tray monitoring, vendor charge limiters, or new forensic sheets) added without breaking existing features.
- **PATCH (`0.0.X`)**: Bug fixes, sensor calibrations, performance improvements, and non-breaking UI adjustments.

### 6.2 Changelog Maintenance (`CHANGELOG.md`)
Every release must document all user-facing changes under standard subheadings:
- `Added`: for new features, drivers, or capabilities.
- `Changed`: for changes in existing functionality.
- `Deprecated`: for soon-to-be removed features.
- `Removed`: for features removed in this version.
- `Fixed`: for any bug fixes.
- `Security`: for vulnerability or privacy enhancements.

### 6.3 Tagging & Release Execution
1. Move items from `## [Unreleased]` into a new version block:
   ```markdown
   ## [0.2.0] - 2026-11-15
   ### Added
   - System tray daemon with background polling...
   ```
2. Commit the changelog and release metadata:
   ```bash
   git commit -am "chore(release): prepare v0.2.0"
   ```
3. Create an annotated Git tag:
   ```bash
   git tag -a v0.2.0 -m "Release v0.2.0: System tray daemon and charge limiters"
   ```
4. Push the branch and tag to GitHub:
   ```bash
   git push origin main --tags
   ```

### 6.4 Automated Release Pipeline
Pushing a `v*` tag triggers the GitHub Actions workflow in [`.github/workflows/release.yml`](.github/workflows/release.yml), which automatically:
1. Validates `bun run check` with 0 errors.
2. Compiles production binaries and packages all 5 Linux targets:
   - **`ohmwatt.deb`**: Debian / Ubuntu package via `wails3 task linux:create:deb`.
   - **`ohmwatt.rpm`**: Fedora / RHEL package via `wails3 task linux:create:rpm`.
   - **`ohmwatt.pkg.tar.zst`**: Arch Linux package via `wails3 task linux:create:aur`.
   - **`ohmwatt-linux-amd64.tar.gz`**: Universal tarball with rootless `install.sh` via `wails3 task linux:create:tarball`.
   - **`ohmwatt-x86_64.AppImage`**: Standalone executable bundling WebKitGTK via `wails3 task linux:create:appimage`.
3. Computes cryptographic SHA256 checksums (`SHA256SUMS.txt`).
4. Generates categorized release notes using [`.github/release.yml`](.github/release.yml).
5. Publishes artifacts directly to [GitHub Releases](https://github.com/osspakistan/ohmwatt/releases).
