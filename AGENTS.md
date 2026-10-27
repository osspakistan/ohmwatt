# antlion-app

A desktop application built with Tauri, Vue 3, TypeScript, and Vite. Combines a Rust backend (via Tauri) with a Vue 3 frontend styled with Tailwind CSS v4 and shadcn-vue components.

## Build & Test

- Install dependencies: `bun install`
- Start development server: `bun run dev`
- Start Tauri dev (desktop app): `bun run tauri dev`
- Build frontend only: `bun run build`
- Build Tauri app: `bun run tauri build`
- Preview build: `bun run preview`
- Type check: `bun run type`
- Lint: `bun run lint`
- Lint (fix): `bun run lint:fix`
- Format: `bun run format`
- Fix all (format + lint): `bun run fix`
- Full check (fix + lint + type): `bun run check`
- Add shadcn-vue component: `bun run shadcn <component>`

## Project Layout

```
├─ src/                → Vue 3 frontend source
│  ├─ components/      → Vue components
│  │  └─ ui/           → shadcn-vue UI components (auto-generated)
│  ├─ lib/             → Utilities (cn helper)
│  ├─ assets/          → Static assets
│  ├─ app.vue          → Root component
│  ├─ main.ts          → App entry point
│  └─ style.css        → Global styles (Tailwind CSS v4)
├─ src-tauri/          → Rust backend (Tauri)
│  ├─ src/             → Rust source (lib.rs, main.rs)
│  ├─ capabilities/    → Tauri permissions/capabilities
│  ├─ icons/           → App icons
│  ├─ Cargo.toml       → Rust dependencies
│  └─ tauri.conf.json  → Tauri configuration
├─ public/             → Public static assets
└─ index.html          → Vite entry HTML
```

## Architecture Overview

**Frontend:** Vue 3 with Composition API (`<script setup>`), TypeScript, Vite 8, Tailwind CSS v4. Uses shadcn-vue (new-york style) with Reka UI primitives and Hugeicons. Path alias `@/` maps to `./src/`. Tauri API calls via `@tauri-apps/api/core` `invoke()`.

**Backend:** Rust via Tauri 2. Commands are defined with `#[tauri::command]` in `src-tauri/src/lib.rs` and registered via `invoke_handler`. Uses `tauri-plugin-opener` for opening URLs/files.

**Communication:** Frontend calls Rust functions using `invoke('command_name', { args })`. Data serialization via serde/serde_json.

## Development Patterns & Constraints

### Coding Style
- **TypeScript:** Strict mode enabled, ES2020 target, bundler module resolution
- **Vue:** Composition API with `<script setup lang="ts">`
- **Imports:** Use path alias `@/` for src imports (configured in tsconfig.json and vite)
- **Formatting:** oxfmt for formatting (`.oxfmtrc.json`)
- **Linting:** Dual linting with oxlint + ESLint (eslint.config.ts), run via `bun run lint`
- **Components:** shadcn-vue pattern - UI components in `src/components/ui/`, organized by component with index exports

### Styling
- **Framework:** Tailwind CSS v4 via `@tailwindcss/vite`
- **Components:** shadcn-vue with style `new-york`, base color `taupe`, CSS variables enabled
- **Icons:** Hugeicons via `@hugeicons/vue` and `@hugeicons/core-free-icons`
- **Animations:** tw-animate-css
- **Font:** Instrument Sans (loaded from Google Fonts)

### Tauri Patterns
- **Commands:** Define in Rust with `#[tauri::command]`, register in `generate_handler![]`
- **Plugins:** Use Tauri plugins (e.g., `tauri-plugin-opener` initialized in `run()`)
- **Frontend API:** Import from `@tauri-apps/api/core` for `invoke()`

## Security

- **Tauri sandboxing:** Commands must be explicitly registered via `invoke_handler`
- **CSP/Permissions:** Configure via `src-tauri/capabilities/` and `tauri.conf.json`
- **Frontend:** Standard web security practices, Tauri provides OS-level API access controls

## Git Workflows

- **Main branch:** `main` (per repo context)
- **Current branch:** `master`
- **Commits:** Follow conventional style as seen in recent history

## Gotchas

- **Tauri dev server:** Vite runs on port 1420 (strictPort), HMR on 1421 if using Tauri dev host
- **Rust lib naming:** `antlion_app_lib` with crate-type `staticlib/cdylib/rlib` (Windows compatibility)
- **Path resolution:** Uses `tsconfigPaths` plugin in Vite (native resolution)
- **Dual linting:** Project uses both oxlint and ESLint - `lint:fix` runs both
- **Build command:** `bun run build` runs type checking (`check`) before building
