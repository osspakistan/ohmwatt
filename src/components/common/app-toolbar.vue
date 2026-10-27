<script setup lang="ts">
import { getCurrentWindow } from '@tauri-apps/api/window'
import { ref } from 'vue'

withDefaults(
  defineProps<{
    title?: string
  }>(),
  {
    title: 'Antlion',
  },
)

const isMaximized = ref(false)

async function handleMinimize() {
  try {
    const appWindow = getCurrentWindow()
    await appWindow.minimize()
  } catch {
    // Non-tauri browser preview fallback
  }
}

async function handleToggleMaximize() {
  try {
    const appWindow = getCurrentWindow()
    await appWindow.toggleMaximize()
    isMaximized.value = await appWindow.isMaximized()
  } catch {
    isMaximized.value = !isMaximized.value
  }
}

async function handleClose() {
  try {
    const appWindow = getCurrentWindow()
    await appWindow.close()
  } catch {
    // Non-tauri browser preview fallback
  }
}

function handleStartDragging() {
  try {
    const appWindow = getCurrentWindow()
    appWindow.startDragging()
  } catch {
    // Non-tauri browser preview fallback
  }
}
</script>

<template>
  <header
    class="flex h-10 flex-none select-none items-center border-b border-border bg-card/60 pl-3 pr-2 backdrop-blur-md"
  >
    <!-- Brand Icon & Title -->
    <div class="flex items-center gap-2">
      <img src="/antlion-mark.svg" alt="Antlion Mark" class="h-5 w-5" />
      <span class="font-serif text-sm font-semibold tracking-tight text-foreground">{{
        title
      }}</span>
    </div>

    <!-- Draggable Region -->
    <div
      class="h-full flex-1 cursor-default"
      title="Drag to move. Double-click to maximize."
      @mousedown="handleStartDragging"
      @dblclick="handleToggleMaximize"
    />

    <!-- Window Action Buttons -->
    <div class="flex items-center gap-1">
      <button
        class="flex h-7 w-8 items-center justify-center rounded text-muted-foreground transition hover:bg-muted hover:text-foreground focus-visible:outline-2 focus-visible:outline-primary"
        aria-label="Minimize"
        @click="handleMinimize"
      >
        <svg
          width="12"
          height="12"
          viewBox="0 0 12 12"
          fill="none"
          stroke="currentColor"
          stroke-width="1.4"
          stroke-linecap="round"
        >
          <path d="M2 6h8" />
        </svg>
      </button>

      <button
        class="flex h-7 w-8 items-center justify-center rounded text-muted-foreground transition hover:bg-muted hover:text-foreground focus-visible:outline-2 focus-visible:outline-primary"
        aria-label="Maximize"
        @click="handleToggleMaximize"
      >
        <svg
          v-if="isMaximized"
          width="12"
          height="12"
          viewBox="0 0 12 12"
          fill="none"
          stroke="currentColor"
          stroke-width="1.4"
          stroke-linejoin="round"
        >
          <rect x="3.5" y="1.5" width="7" height="7" rx="1.5" />
          <path d="M1.5 4v5.5a1 1 0 001 1H8" />
        </svg>
        <svg
          v-else
          width="12"
          height="12"
          viewBox="0 0 12 12"
          fill="none"
          stroke="currentColor"
          stroke-width="1.4"
          stroke-linejoin="round"
        >
          <rect x="2" y="2" width="8" height="8" rx="1.5" />
        </svg>
      </button>

      <button
        class="group flex h-7 w-8 items-center justify-center rounded text-muted-foreground transition hover:bg-destructive hover:text-destructive-foreground focus-visible:outline-2 focus-visible:outline-primary"
        aria-label="Close"
        @click="handleClose"
      >
        <svg
          width="12"
          height="12"
          viewBox="0 0 12 12"
          fill="none"
          stroke="currentColor"
          stroke-width="1.4"
          stroke-linecap="round"
        >
          <path d="M2.5 2.5l7 7M9.5 2.5l-7 7" />
        </svg>
      </button>
    </div>
  </header>
</template>
