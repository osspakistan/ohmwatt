<script setup lang="ts">
import { ArrowLeft01Icon, ArrowRight01Icon } from '@hugeicons/core-free-icons'
import { HugeiconsIcon } from '@hugeicons/vue'
import { Window } from '@wailsio/runtime'
import { ref } from 'vue'
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

// Vue Router History Stack Tracker
const historyStack = ref<string[]>([window.location.hash || '#/dashboard'])
const currentIndex = ref(0)

router.afterEach((to) => {
  const target = `#${to.fullPath}`
  if (historyStack.value[currentIndex.value] !== target) {
    // If user navigated forward to a new route, discard forward redo stack
    historyStack.value = historyStack.value.slice(0, currentIndex.value + 1)
    historyStack.value.push(target)
    currentIndex.value = historyStack.value.length - 1
  }
})

function canGoBack() {
  return currentIndex.value > 0
}

function canGoForward() {
  return currentIndex.value < historyStack.value.length - 1
}

function goBack() {
  if (canGoBack()) {
    currentIndex.value--
    router.back()
  }
}

function goForward() {
  if (canGoForward()) {
    currentIndex.value++
    router.forward()
  }
}

function handleMinimize() {
  try {
    Window.Minimise()
  } catch (e) {
    console.warn('Window.Minimise error', e)
  }
}

async function handleToggleMaximize() {
  try {
    isMaximized.value = !isMaximized.value
    await Window.ToggleMaximise()
  } catch (e) {
    console.warn('Window.ToggleMaximise error', e)
  }
}

function handleClose() {
  try {
    Window.Close()
  } catch (e) {
    console.warn('Window.Close error', e)
  }
}
</script>

<template>
  <header
    class="flex h-10 flex-none select-none items-center justify-between border-b border-border bg-card/60 px-3 backdrop-blur-md"
    style="--wails-draggable: drag;"
  >
    <!-- Brand Icon, Title & Navigation History Controls -->
    <div class="flex items-center gap-3" style="--wails-draggable: no-drag;">
      <!-- Functional Back / Forward Buttons -->
      <div class="flex items-center gap-0.5">
        <button
          type="button"
          :disabled="!canGoBack()"
          :class="[
            'flex h-7 w-7 items-center justify-center rounded transition',
            canGoBack()
              ? 'text-foreground hover:bg-muted hover:text-primary cursor-pointer'
              : 'text-muted-foreground/30 cursor-not-allowed',
          ]"
          title="Back (Alt+Left)"
          aria-label="Go Back"
          @click="goBack"
        >
          <HugeiconsIcon :icon="ArrowLeft01Icon" class="h-3.5 w-3.5" />
        </button>

        <button
          type="button"
          :disabled="!canGoForward()"
          :class="[
            'flex h-7 w-7 items-center justify-center rounded transition',
            canGoForward()
              ? 'text-foreground hover:bg-muted hover:text-primary cursor-pointer'
              : 'text-muted-foreground/30 cursor-not-allowed',
          ]"
          title="Forward (Alt+Right)"
          aria-label="Go Forward"
          @click="goForward"
        >
          <HugeiconsIcon :icon="ArrowRight01Icon" class="h-3.5 w-3.5" />
        </button>
      </div>

      <div class="h-3.5 w-[1px] bg-border" />

      <!-- Brand Identity matching Antlion Toolbar Spec -->
      <div class="flex items-center gap-2">
        <svg viewBox="0 0 48 48" class="h-5 w-5 flex-none text-foreground" aria-hidden="true">
          <rect x="12.5" y="9" width="23" height="35" rx="6.5" fill="none" stroke="currentColor" stroke-width="3.5"/>
          <rect x="19" y="3" width="10" height="5.5" rx="2" fill="currentColor"/>
          <rect x="18" y="25" width="12" height="13" rx="3" fill="#644a40"/>
        </svg>
        <span class="font-serif text-sm font-semibold tracking-tight text-foreground">{{
          title
        }}</span>
      </div>
    </div>

    <!-- Draggable Central Region -->
    <div
      class="flex-1 h-full mx-2"
      title="Drag to move. Double-click to maximize."
      @dblclick="handleToggleMaximize"
    />

    <!-- Window Action Controls: Minimize, Maximize/Restore, Close -->
    <div class="flex items-center gap-1" style="--wails-draggable: no-drag;">
      <!-- Minimize -->
      <button
        type="button"
        class="flex h-7 w-8 items-center justify-center rounded-md text-muted-foreground transition hover:bg-muted hover:text-foreground focus-visible:outline-2 focus-visible:outline-primary cursor-pointer"
        aria-label="Minimize"
        title="Minimize"
        @click="handleMinimize"
      >
        <svg width="12" height="12" viewBox="0 0 12 12" fill="none" stroke="currentColor" stroke-width="1.4" stroke-linecap="round">
          <path d="M2 6h8" />
        </svg>
      </button>

      <!-- Maximize / Restore -->
      <button
        type="button"
        class="flex h-7 w-8 items-center justify-center rounded-md text-muted-foreground transition hover:bg-muted hover:text-foreground focus-visible:outline-2 focus-visible:outline-primary cursor-pointer"
        :aria-label="isMaximized ? 'Restore' : 'Maximize'"
        :title="isMaximized ? 'Restore' : 'Maximize'"
        @click="handleToggleMaximize"
      >
        <svg v-if="!isMaximized" width="12" height="12" viewBox="0 0 12 12" fill="none" stroke="currentColor" stroke-width="1.4" stroke-linejoin="round">
          <rect x="2" y="2" width="8" height="8" rx="1.5" />
        </svg>
        <svg v-else width="12" height="12" viewBox="0 0 12 12" fill="none" stroke="currentColor" stroke-width="1.4" stroke-linejoin="round">
          <rect x="3.5" y="1.5" width="7" height="7" rx="1.5" />
          <path d="M1.5 4v5.5a1 1 0 001 1H8" />
        </svg>
      </button>

      <!-- Close Window -->
      <button
        type="button"
        class="flex h-7 w-8 items-center justify-center rounded-md text-muted-foreground transition hover:bg-[#e54d2e] hover:text-white focus-visible:outline-2 focus-visible:outline-primary cursor-pointer"
        aria-label="Close"
        title="Close"
        @click="handleClose"
      >
        <svg width="12" height="12" viewBox="0 0 12 12" fill="none" stroke="currentColor" stroke-width="1.4" stroke-linecap="round">
          <path d="M2.5 2.5l7 7M9.5 2.5l-7 7" />
        </svg>
      </button>
    </div>
  </header>
</template>
