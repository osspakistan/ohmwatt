<script setup lang="ts">
import { computed, ref } from 'vue'
import { useRouter } from 'vue-router'

import LiquidVessel from './components/liquid-vessel.vue'

import StatusPill from '@/components/common/status-pill.vue'

const router = useRouter()

// Reactive state mirroring prototype values
const percentage = ref(78)
const isCharging = ref(true)
const headsUpThreshold = ref(90)
const batteryCount = ref(1) // Toggle 1 or 2
const voltage = ref('12.9 V')
const current = ref('2.51 A')
const temperature = ref('34 °C')
const chargerWattage = ref('65 W')
const isUnderpowered = ref(true) // 65W charger on 90W laptop

const sentence = computed(() => {
  if (isCharging.value) {
    return 'About 48 minutes until your 90% heads-up.'
  }
  return 'About 5 hours 12 minutes of battery remaining.'
})

// Hour strip data: 48 half-hour slots
const hourSlots = Array.from({ length: 48 }, (_, i) => {
  const h = i / 2
  const level = h < 14 ? 100 - h * 5.7 : 20 + (h - 14) * 5.8
  const charging = h >= 14
  return {
    height: Math.max(8, Math.min(100, Math.round(level))),
    charging,
  }
})
</script>

<template>
  <div class="max-w-4xl mx-auto px-6 sm:px-10 py-8">
    <div class="grid md:grid-cols-[250px_1fr] gap-10 items-start">
      <!-- Left Hero: Animated Liquid Vessel -->
      <div class="flex flex-col items-center">
        <!-- Single or Dual Vessel -->
        <div v-if="batteryCount === 2" class="grid grid-cols-2 gap-3">
          <LiquidVessel
            :percentage="percentage"
            :heads-up-threshold="headsUpThreshold"
            :is-charging="isCharging"
            :small="true"
            battery-label="Internal"
            id="v1"
          />
          <LiquidVessel
            :percentage="64"
            :is-charging="isCharging"
            :small="true"
            battery-label="Removable"
            id="v2"
          />
        </div>
        <LiquidVessel
          v-else
          :percentage="percentage"
          :heads-up-threshold="headsUpThreshold"
          :is-charging="isCharging"
          id="main"
        />
      </div>

      <!-- Right Details & Live Telemetry -->
      <div class="pt-2">
        <StatusPill :status="isCharging ? 'charging' : 'battery'" />

        <!-- Big Lora Serif Number Readout -->
        <div class="flex items-baseline gap-2 mt-4">
          <span class="font-serif text-8xl font-light leading-none tracking-tight text-foreground">
            {{ percentage }}
          </span>
          <span class="font-serif text-4xl text-muted-foreground">%</span>
        </div>

        <p class="text-base sm:text-lg mt-3 text-muted-foreground max-w-md">
          {{ sentence }}
        </p>

        <!-- 4-Grid Telemetry Strip -->
        <div class="grid grid-cols-2 sm:grid-cols-4 gap-y-5 gap-x-4 mt-8">
          <div class="border-l border-border pl-3">
            <div class="text-xs text-muted-foreground uppercase tracking-wider">Voltage</div>
            <div class="font-mono text-lg font-medium text-foreground mt-0.5">{{ voltage }}</div>
          </div>
          <div class="border-l border-border pl-3">
            <div class="text-xs text-muted-foreground uppercase tracking-wider">Current</div>
            <div class="font-mono text-lg font-medium text-foreground mt-0.5">{{ current }}</div>
          </div>
          <div class="border-l border-border pl-3">
            <div class="text-xs text-muted-foreground uppercase tracking-wider">Temperature</div>
            <div class="font-mono text-lg font-medium text-foreground mt-0.5">
              {{ temperature }}
            </div>
          </div>
          <div class="border-l border-border pl-3">
            <div class="text-xs text-muted-foreground uppercase tracking-wider">Charger</div>
            <div class="font-mono text-lg font-medium text-foreground mt-0.5">
              {{ chargerWattage }}
            </div>
          </div>
        </div>

        <!-- Underpowered Charger Warning Hint -->
        <div
          v-if="isUnderpowered && isCharging"
          class="mt-4 rounded-md border border-destructive/20 bg-destructive/10 px-3 py-2 text-xs text-destructive flex items-center gap-2"
        >
          <span>⚠</span>
          <span
            >Your charger is 65 W and this laptop is rated for 90 W. Charging may be slow under
            load.</span
          >
        </div>

        <!-- Hour-by-Hour Battery Distribution Strip -->
        <div class="mt-8 border-t border-border pt-6">
          <h2 class="text-xs uppercase tracking-wider font-semibold text-muted-foreground mb-3">
            Today, hour by hour
          </h2>
          <div
            class="flex items-end gap-[2px] h-16 w-full"
            role="img"
            aria-label="Battery level over the last 24 hours"
          >
            <div
              v-for="(slot, idx) in hourSlots"
              :key="idx"
              :class="[
                'flex-1 rounded-sm transition-all',
                slot.charging ? 'bg-primary' : 'bg-muted-foreground',
              ]"
              :style="{
                height: `${slot.height}%`,
                opacity: slot.charging ? 1 : 0.45,
              }"
            />
          </div>

          <!-- Time Legend Labels -->
          <div class="flex justify-between text-[11px] font-mono text-muted-foreground mt-2">
            <span>12 AM</span>
            <span>6 AM</span>
            <span>12 PM</span>
            <span>6 PM</span>
            <span>Now</span>
          </div>

          <div class="flex items-center gap-4 text-xs text-muted-foreground mt-3">
            <span class="flex items-center gap-1.5">
              <i class="h-2.5 w-2.5 rounded-sm bg-muted-foreground opacity-45" />
              On battery
            </span>
            <span class="flex items-center gap-1.5">
              <i class="h-2.5 w-2.5 rounded-sm bg-primary" />
              Charging
            </span>
          </div>
        </div>

        <!-- Quick Jump Links to Health & Device -->
        <div class="mt-8 space-y-2 border-t border-border pt-4">
          <button
            type="button"
            class="w-full text-left py-2 flex items-baseline justify-between hover:text-primary transition group"
            @click="router.push('/health')"
          >
            <div>
              <span class="font-serif text-2xl font-medium">91%</span>
              <span class="text-sm text-muted-foreground ml-2"
                >of original capacity · 212 cycles</span
              >
            </div>
            <span
              class="text-xs uppercase font-mono tracking-wider text-muted-foreground group-hover:text-primary"
            >
              View Health →
            </span>
          </button>

          <button
            type="button"
            class="w-full text-left py-2 flex items-baseline justify-between hover:text-primary transition group border-t border-border/60 pt-3"
            @click="router.push('/device')"
          >
            <span class="text-sm text-muted-foreground"> HP EliteBook 850 G7 · SMP battery </span>
            <span
              class="text-xs uppercase font-mono tracking-wider text-muted-foreground group-hover:text-primary"
            >
              Device Specs →
            </span>
          </button>
        </div>
      </div>
    </div>
  </div>
</template>
