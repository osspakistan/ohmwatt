<script setup lang="ts">
import { computed, ref } from 'vue'

const props = withDefaults(
  defineProps<{
    value: string
    maskLength?: number
  }>(),
  {
    maskLength: 8,
  },
)

const revealed = ref(false)

const maskedText = computed(() => {
  if (revealed.value) return props.value
  const visibleEnding = props.value.length > 4 ? props.value.slice(-4) : props.value
  return '•••• ' + visibleEnding
})
</script>

<template>
  <div class="inline-flex items-center gap-2 font-mono text-sm">
    <span class="tabular-nums text-muted-foreground">{{ maskedText }}</span>
    <button
      type="button"
      class="text-xs text-primary underline underline-offset-2 hover:opacity-80 transition"
      @click="revealed = !revealed"
    >
      {{ revealed ? 'Hide' : 'Show' }}
    </button>
  </div>
</template>
