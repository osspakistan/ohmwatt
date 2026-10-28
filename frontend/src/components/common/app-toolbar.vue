<script setup lang="ts">
import { ArrowLeft01Icon, ArrowRight01Icon } from '@hugeicons/core-free-icons'
import { HugeiconsIcon } from '@hugeicons/vue'
import { computed } from 'vue'
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
</script>

<template>
  <header
    class="flex h-10 flex-none select-none items-center justify-between border-b border-border bg-card/60 px-3 backdrop-blur-md"
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
  </header>
</template>
