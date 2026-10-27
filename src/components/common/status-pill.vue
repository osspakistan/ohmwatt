<script setup lang="ts">
import { computed } from 'vue'

const props = withDefaults(
  defineProps<{
    status: 'charging' | 'battery' | 'full' | 'idle'
  }>(),
  {
    status: 'battery',
  },
)

const config = computed(() => {
  switch (props.status) {
    case 'charging':
      return {
        label: 'Plugged in, charging',
        classes: 'bg-secondary/30 text-secondary-foreground border-secondary/50',
        dot: 'bg-primary animate-pulse',
      }
    case 'full':
      return {
        label: 'Plugged in, full',
        classes: 'bg-secondary/30 text-secondary-foreground border-secondary/50',
        dot: 'bg-primary',
      }
    case 'idle':
      return {
        label: 'Plugged in, not charging',
        classes: 'bg-muted text-foreground border-border',
        dot: 'bg-muted-foreground',
      }
    case 'battery':
    default:
      return {
        label: 'On battery',
        classes: 'bg-muted/80 text-foreground border-border',
        dot: 'bg-muted-foreground',
      }
  }
})
</script>

<template>
  <span
    :class="[
      'inline-flex items-center gap-1.5 rounded-full border px-3 py-1 text-xs font-medium transition-colors',
      config.classes,
    ]"
  >
    <i :class="['h-1.5 w-1.5 rounded-full', config.dot]" />
    {{ config.label }}
  </span>
</template>
