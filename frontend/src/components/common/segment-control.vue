<script setup lang="ts">
export interface SegmentOption {
  value: string
  label: string
}

defineProps<{
  options: SegmentOption[]
  modelValue: string
}>()

const emit = defineEmits<{
  'update:modelValue': [val: string]
}>()
</script>

<template>
  <div class="inline-flex rounded-lg border border-border bg-muted p-0.5 text-xs">
    <button
      v-for="opt in options"
      :key="opt.value"
      type="button"
      :class="[
        'rounded-md px-3 py-1.5 transition font-medium',
        modelValue === opt.value
          ? 'bg-background text-foreground shadow-sm'
          : 'text-muted-foreground hover:text-foreground',
      ]"
      @click="emit('update:modelValue', opt.value)"
    >
      {{ opt.label }}
    </button>
  </div>
</template>
