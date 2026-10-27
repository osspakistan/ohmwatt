<script setup lang="ts">
import { ArrowLeft01Icon, ArrowRight01Icon } from '@hugeicons/core-free-icons'
import { HugeiconsIcon } from '@hugeicons/vue'
import { getCurrentWindow } from '@tauri-apps/api/window'
import { computed, ref } from 'vue'
import { useRouter } from 'vue-router'

withDefaults(
  defineProps<{
    title?: string
  }>(),
  {
    title: 'Antlion',
  },
)

const router = useRouter()
const isMaximized = ref(false)

const canGoBack = computed(() => {
  if (typeof window !== 'undefined' && window.history) {
    return window.history.state?.back !== null
  }
  return false
})

const canGoForward = computed(() => {
  if (typeof window !== 'undefined' && window.history) {
    return window.history.state?.forward !== null
  }
  return false
})

function goBack() {
  router.back()
}

function goForward() {
  router.forward()
}

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
    <!-- Brand Icon, Title & Navigation History Controls -->
    <div class="flex items-center gap-3">
      <!-- Back / Forward Buttons -->
      <div class="flex items-center gap-0.5">
        <button
          type="button"
          :disabled="!canGoBack"
          :class="[
            'flex h-7 w-7 items-center justify-center rounded transition',
            canGoBack
              ? 'text-foreground hover:bg-muted hover:text-primary cursor-pointer'
              : 'text-muted-foreground/40 cursor-not-allowed',
          ]"
          title="Back (Alt+Left)"
          aria-label="Go Back"
          @click="goBack"
        >
          <HugeiconsIcon :icon="ArrowLeft01Icon" class="h-3.5 w-3.5" />
        </button>

        <button
          type="button"
          :disabled="!canGoForward"
          :class="[
            'flex h-7 w-7 items-center justify-center rounded transition',
            canGoForward
              ? 'text-foreground hover:bg-muted hover:text-primary cursor-pointer'
              : 'text-muted-foreground/40 cursor-not-allowed',
          ]"
          title="Forward (Alt+Right)"
          aria-label="Go Forward"
          @click="goForward"
        >
          <HugeiconsIcon :icon="ArrowRight01Icon" class="h-3.5 w-3.5" />
        </button>
      </div>

      <div class="h-3.5 w-[1px] bg-border" />

      <!-- Brand Identity -->
      <div class="flex items-center gap-2">
        <img src="/antlion-mark.svg" alt="Antlion Mark" class="h-4.5 w-4.5" />
        <span class="font-serif text-sm font-semibold tracking-tight text-foreground">{{
          title
        }}</span>
      </div>
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
