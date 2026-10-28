<script setup lang="ts">
import { computed, ref, onMounted } from 'vue'

import { BatteryService } from '../../../bindings/changeme/service'

import { useBattery } from '@/composables/use-battery'

const { battery } = useBattery()

const currentHealth = computed(() => battery.value.healthPercentage)
const cycleCount = computed(() => battery.value.cycleCount)
const designWh = computed(() => ((battery.value.designCapacity * 11.55) / 1000).toFixed(1))
const currentWh = computed(() => ((battery.value.fullCapacity * 11.55) / 1000).toFixed(1))

const lostPercent = computed(() => Math.max(0, 100 - currentHealth.value))
const hatchedHeight = computed(() => Math.max(4, Math.round((lostPercent.value / 100) * 142)))
const retainedHeight = computed(() => 142 - hatchedHeight.value)
const retainedY = computed(() => 14 + hatchedHeight.value)

const forecastDate = ref('Calculating...')
const hoursAtHighStress = ref(0)
const deepDischarges = ref(0)

onMounted(async () => {
  try {
    const h = await BatteryService.GetBatteryHealth()
    if (h) {
      if (h.estimatedDate80) forecastDate.value = h.estimatedDate80
      hoursAtHighStress.value = h.hoursAtHighStress ?? 0
      deepDischarges.value = h.deepDischarges ?? 0
    }
  } catch (err) {
    console.debug('Failed to query battery health:', err)
  }
})
</script>

<template>
  <div class="max-w-3xl mx-auto px-6 sm:px-10 py-9">
    <!-- Header -->
    <header class="mb-8">
      <h1 class="font-serif text-4xl tracking-tight text-foreground">Health</h1>
      <p class="text-muted-foreground mt-1.5">
        Your battery holds {{ currentHealth }}% of what it held when new.
      </p>
    </header>

    <!-- Side-by-side Battery Capacity Comparison SVG -->
    <div
      class="rounded-xl border border-border bg-card p-6 flex flex-col items-center justify-center"
    >
      <svg
        viewBox="0 0 300 215"
        class="w-full max-w-[340px] text-foreground"
        role="img"
        aria-label="Battery capacity comparison"
      >
        <defs>
          <pattern
            id="hatched"
            width="8"
            height="8"
            patternUnits="userSpaceOnUse"
            patternTransform="rotate(45)"
          >
            <line
              x1="0"
              y1="0"
              x2="0"
              y2="8"
              stroke="var(--destructive)"
              stroke-width="2.5"
              opacity="0.6"
            />
          </pattern>
        </defs>

        <!-- Left Battery: When New (100%) -->
        <g transform="translate(30, 0)">
          <rect x="38" y="2" width="24" height="9" rx="3" fill="currentColor" />
          <rect
            x="2"
            y="10"
            width="96"
            height="150"
            rx="20"
            fill="none"
            stroke="currentColor"
            stroke-width="3"
          />
          <rect x="6" y="14" width="88" height="142" rx="16" fill="var(--primary)" opacity="0.8" />
          <text
            x="50"
            y="184"
            text-anchor="middle"
            font-size="14"
            fill="var(--foreground)"
            font-weight="600"
          >
            When new
          </text>
          <text x="50" y="202" text-anchor="middle" font-size="12" fill="var(--muted-foreground)">
            {{ designWh }} Wh
          </text>
        </g>

        <!-- Right Battery: Today (dynamic lost % hatched) -->
        <g transform="translate(170, 0)">
          <rect x="38" y="2" width="24" height="9" rx="3" fill="currentColor" />
          <rect
            x="2"
            y="10"
            width="96"
            height="150"
            rx="20"
            fill="none"
            stroke="currentColor"
            stroke-width="3"
          />
          <!-- Lost capacity hatched top block -->
          <rect x="6" y="14" width="88" :height="hatchedHeight" rx="3" fill="url(#hatched)" />
          <!-- Retained capacity fill -->
          <rect
            x="6"
            :y="retainedY"
            width="88"
            :height="retainedHeight"
            rx="16"
            fill="var(--primary)"
          />
          <text
            x="50"
            y="184"
            text-anchor="middle"
            font-size="14"
            fill="var(--foreground)"
            font-weight="600"
          >
            Today
          </text>
          <text x="50" y="202" text-anchor="middle" font-size="12" fill="var(--muted-foreground)">
            {{ currentWh }} Wh ({{ currentHealth }}%)
          </text>
        </g>
      </svg>
    </div>

    <!-- Health Stats & Decay Milestone -->
    <div class="mt-8 grid gap-4 sm:grid-cols-2">
      <div class="rounded-lg border border-border bg-card p-4">
        <div class="text-xs uppercase tracking-wider text-muted-foreground">Condition Rating</div>
        <div class="font-serif text-2xl font-medium text-foreground mt-1">
          {{ currentHealth >= 80 ? 'Normal · Grade A' : 'Service Recommended' }}
        </div>
        <p class="text-xs text-muted-foreground mt-1">
          Capacity decay is within healthy manufacturer parameters for {{ cycleCount }} cycles.
        </p>
      </div>

      <div class="rounded-lg border border-border bg-card p-4">
        <div class="text-xs uppercase tracking-wider text-muted-foreground">
          80% Health Projection
        </div>
        <div class="font-serif text-2xl font-medium text-foreground mt-1">{{ forecastDate }}</div>
        <p class="text-xs text-muted-foreground mt-1">
          Based on degradation rate over completed hardware charge cycles.
        </p>
      </div>
    </div>

    <!-- Stress Factors & Time-at-Charge (F32) -->
    <div class="mt-8 border-t border-border pt-6">
      <h2 class="text-sm font-semibold tracking-wide uppercase text-muted-foreground mb-4">
        Battery Aging Stress Factors
      </h2>
      <div class="space-y-3">
        <div class="flex items-center justify-between border-b border-border py-2.5">
          <div>
            <div class="text-sm font-medium">Time held above 90% while plugged in</div>
            <div class="text-xs text-muted-foreground">
              High voltage stress accelerates chemical decay
            </div>
          </div>
          <span class="font-mono text-sm text-foreground">{{ hoursAtHighStress }} hours</span>
        </div>

        <div class="flex items-center justify-between border-b border-border py-2.5">
          <div>
            <div class="text-sm font-medium">Deep discharges below 10%</div>
            <div class="text-xs text-muted-foreground">
              Full drain cycles stress electrode materials
            </div>
          </div>
          <span class="font-mono text-sm text-foreground">{{ deepDischarges }} times</span>
        </div>
      </div>
    </div>
  </div>
</template>
