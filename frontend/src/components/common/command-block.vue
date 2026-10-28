<script setup lang="ts">
import { ref } from 'vue'

const props = defineProps<{
  command: string
  label?: string
}>()

const copied = ref(false)

async function copyToClipboard() {
  try {
    await navigator.clipboard.writeText(props.command)
    copied.value = true
    setTimeout(() => {
      copied.value = false
    }, 2000)
  } catch {
    // clipboard fallback
  }
}
</script>

<template>
  <div class="relative overflow-hidden rounded-lg border border-border bg-card">
    <div
      v-if="label"
      class="flex items-center justify-between border-b border-border bg-muted/40 px-3 py-1.5 text-xs text-muted-foreground"
    >
      <span>{{ label }}</span>
      <span v-if="copied" class="text-xs font-medium text-primary">Copied!</span>
    </div>

    <div class="flex items-center justify-between p-3">
      <code class="font-mono text-xs text-foreground overflow-x-auto pr-3">{{ command }}</code>
      <button
        type="button"
        class="flex-none rounded border border-border bg-background px-2.5 py-1 text-xs text-muted-foreground transition hover:bg-muted hover:text-foreground"
        @click="copyToClipboard"
      >
        {{ copied ? 'Copied' : 'Copy' }}
      </button>
    </div>
  </div>
</template>
