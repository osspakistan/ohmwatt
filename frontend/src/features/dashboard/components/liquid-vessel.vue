<script setup lang="ts">
import { computed } from 'vue'

const props = withDefaults(
  defineProps<{
    percentage: number
    headsUpThreshold?: number
    isCharging?: boolean
    small?: boolean
    batteryLabel?: string
    id?: string
  }>(),
  {
    headsUpThreshold: 90,
    isCharging: false,
    small: false,
    batteryLabel: '',
    id: 'a',
  },
)

// Wave Path calculation matching prototype
const wavePath = computed(() => {
  return 'M0 25 C40 10, 80 40, 120 25 C160 10, 200 40, 240 25 L240 500 L0 500 Z'
})

// Y translation: 100% -> 20px, 0% -> 370px
const translateY = computed(() => {
  return (370 - (350 * props.percentage) / 100).toFixed(1)
})

const thresholdY = computed(() => {
  return (370 - (350 * (props.headsUpThreshold || 90)) / 100).toFixed(1)
})
</script>

<template>
  <div class="flex flex-col items-center">
    <svg
      viewBox="0 0 240 384"
      :class="['w-full text-foreground', small ? 'max-w-[130px]' : 'max-w-[240px]']"
      role="img"
      :aria-label="`Battery at ${percentage} percent`"
    >
      <defs>
        <clipPath :id="`vessel-clip-${id}`">
          <rect x="20" y="20" width="160" height="350" rx="38" />
        </clipPath>
      </defs>

      <!-- Battery Terminal Cap -->
      <rect x="80" y="4" width="40" height="14" rx="5" fill="currentColor" />

      <!-- Clipped Vessel Interior -->
      <g :clip-path="`url(#vessel-clip-${id})`">
        <!-- Background Chamber -->
        <rect x="20" y="20" width="160" height="350" fill="var(--muted)" />

        <!-- Animated Liquid Mass -->
        <g
          :style="{
            transform: `translateY(${translateY}px)`,
            transition: 'transform 1.2s cubic-bezier(0.34, 1.56, 0.64, 1)',
          }"
        >
          <!-- Back Wave -->
          <path
            :d="wavePath"
            fill="var(--primary)"
            class="opacity-60 transition-transform duration-1000"
            :style="{ transform: `translateX(${-20}px)` }"
          />

          <!-- Front Wave -->
          <path :d="wavePath" fill="var(--primary)" class="opacity-95" />

          <!-- Rising Bubbles when Charging -->
          <g v-if="isCharging" class="animate-pulse">
            <circle cx="60" cy="80" r="3" fill="var(--primary-foreground)" opacity="0.4" />
            <circle cx="100" cy="140" r="4.5" fill="var(--primary-foreground)" opacity="0.3" />
            <circle cx="130" cy="90" r="2.5" fill="var(--primary-foreground)" opacity="0.5" />
            <circle cx="80" cy="180" r="3.5" fill="var(--primary-foreground)" opacity="0.4" />
          </g>
        </g>
      </g>

      <!-- Outer Shell Border -->
      <rect
        x="20"
        y="20"
        width="160"
        height="350"
        rx="38"
        fill="none"
        stroke="currentColor"
        stroke-width="3.5"
      />

      <!-- Dashed 90% Heads-Up Line -->
      <line
        v-if="headsUpThreshold"
        x1="24"
        x2="176"
        :y1="thresholdY"
        :y2="thresholdY"
        stroke="var(--secondary)"
        stroke-dasharray="4 3"
        stroke-width="2"
      />
    </svg>

    <!-- Sub-label for dual battery setups -->
    <span v-if="batteryLabel" class="mt-2 text-xs font-mono text-muted-foreground">
      {{ batteryLabel }} {{ percentage }}%
    </span>
  </div>
</template>
