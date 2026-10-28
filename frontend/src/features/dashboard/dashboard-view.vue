<script setup lang="ts">
import { computed, ref, onMounted } from 'vue'
import { useRouter } from 'vue-router'

import { DeviceService, BatteryService } from '../../../bindings/changeme/service'

import LiquidVessel from './components/liquid-vessel.vue'

import StatusPill from '@/components/common/status-pill.vue'
import { useBattery } from '@/composables/use-battery'

const router = useRouter()
const { battery } = useBattery()

// Reactive state mirroring prototype values and native battery readings
const percentage = computed(() => battery.value.percentage)
const isCharging = computed(() => battery.value.isCharging)
const headsUpThreshold = ref(90)
const batteryCount = ref(1) // Automatically detected from hardware (1 or 2)
const voltage = computed(() => `${battery.value.voltageVolts} V`)
const current = computed(
  () => `${(battery.value.powerWatts / (battery.value.voltageVolts || 1)).toFixed(2)} A`,
)
const temperature = computed(() => `${battery.value.temperatureC} °C`)
const chargerWattage = computed(() =>
  battery.value.isAcConnected
    ? `${battery.value.powerWatts > 0 ? battery.value.powerWatts : 65} W`
    : 'Disconnected',
)
const isUnderpowered = ref(false)

const deviceLabel = ref('Detecting device...')
const batteryLabel = ref('Hardware battery')

onMounted(async () => {
  try {
    const s = await DeviceService.GetDeviceSpecs()
    if (s) {
      if (s.laptopModel) {
        deviceLabel.value = s.laptopModel.replace(/ Notebook PC$/i, '')
      }
      if (s.batteryManufacturer) {
        batteryLabel.value = `${s.batteryManufacturer} battery`
      }
    }
  } catch (err) {
    console.debug('Failed to query device specs on dashboard:', err)
  }

  try {
    const count = await BatteryService.GetBatteryCount()
    if (count && count > 0) {
      batteryCount.value = count
    }
  } catch (err) {
    console.debug('Failed to query battery count:', err)
  }
})

const sentence = computed(() => {
  if (isCharging.value) {
    if (battery.value.timeToFullMins > 0) {
      return `About ${battery.value.timeToFullMins} minutes until full.`
    }
    return 'Charging connected. Approaching 90% heads-up.'
  }
  if (battery.value.timeToEmptyMins > 0) {
    const hrs = Math.floor(battery.value.timeToEmptyMins / 60)
    const mins = battery.value.timeToEmptyMins % 60
    return `About ${hrs > 0 ? hrs + ' hours ' : ''}${mins} minutes of battery remaining.`
  }
  return `Discharging at ${battery.value.powerWatts} W.`
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
      <!-- Left Column: Animated Liquid Vessel + Quick Navigation Links -->
      <div class="flex flex-col items-center">
        <!-- Single or Dual Vessel -->
        <div v-if="batteryCount === 2" class="grid grid-cols-2 gap-3">
          <LiquidVessel
            :percentage="percentage"
            :heads-up-threshold="headsUpThreshold"
            :is-charging="isCharging"
            small
            battery-label="Internal"
          />
          <LiquidVessel
            :percentage="64"
            :heads-up-threshold="headsUpThreshold"
            :is-charging="false"
            small
            battery-label="Removable"
          />
        </div>
        <LiquidVessel
          v-else
          :percentage="percentage"
          :heads-up-threshold="headsUpThreshold"
          :is-charging="isCharging"
        />

        <!-- Quick Jump Links (Under battery vessel) -->
        <div class="w-full mt-6 space-y-3 border-t border-border pt-4">
          <button
            type="button"
            class="w-full text-left py-2 hover:opacity-80 transition group cursor-pointer block"
            @click="router.push('/health')"
          >
            <div class="flex items-baseline justify-between">
              <span class="font-serif text-2xl font-medium text-foreground leading-none">
                {{ battery.healthPercentage }}%
              </span>
              <span
                class="text-[10px] uppercase font-mono tracking-wider text-muted-foreground group-hover:text-primary transition"
              >
                View Health →
              </span>
            </div>
            <p class="text-xs text-muted-foreground mt-1.5 leading-snug">
              of original capacity · {{ battery.cycleCount }} cycles
            </p>
          </button>

          <button
            type="button"
            class="w-full text-left py-2 border-t border-border/60 pt-3 hover:opacity-80 transition group cursor-pointer block"
            @click="router.push('/device')"
          >
            <div class="flex items-baseline justify-between">
              <span
                class="text-xs font-medium text-foreground truncate max-w-[130px]"
                :title="deviceLabel"
              >
                {{ deviceLabel }}
              </span>
              <span
                class="text-[10px] uppercase font-mono tracking-wider text-muted-foreground group-hover:text-primary transition shrink-0"
              >
                Device Specs →
              </span>
            </div>
            <p class="text-[11px] text-muted-foreground font-mono mt-1">
              {{ batteryLabel }}
            </p>
          </button>
        </div>
      </div>

      <!-- Right Panel: Diagnostics, Metric Tiles, and Hour Strip -->
      <div class="space-y-6">
        <!-- Top Status Pill & Dynamic Sentence -->
        <div>
          <StatusPill
            :status="isCharging ? (percentage >= 100 ? 'full' : 'charging') : 'battery'"
          />

          <!-- Big Serif Percentage -->
          <div class="flex items-baseline gap-2 mt-3">
            <span
              class="font-serif text-8xl font-light leading-none tracking-tight text-foreground"
            >
              {{ percentage }}
            </span>
            <span class="font-serif text-4xl text-muted-foreground">%</span>
          </div>

          <p class="text-sm text-muted-foreground mt-2 leading-relaxed">
            {{ sentence }}
          </p>
        </div>

        <!-- Live Metric Grid: 4 columns when charging, 3 columns when on battery -->
        <div
          class="grid grid-cols-2 gap-3 pt-2"
          :class="isCharging ? 'sm:grid-cols-4' : 'sm:grid-cols-3'"
        >
          <div class="rounded-xl border border-border bg-card p-3">
            <div class="text-[11px] font-mono uppercase tracking-wider text-muted-foreground">
              Voltage
            </div>
            <div class="font-mono text-lg font-medium text-foreground mt-0.5">{{ voltage }}</div>
          </div>
          <div class="rounded-xl border border-border bg-card p-3">
            <div class="text-[11px] font-mono uppercase tracking-wider text-muted-foreground">
              Current
            </div>
            <div class="font-mono text-lg font-medium text-foreground mt-0.5">{{ current }}</div>
          </div>
          <div class="rounded-xl border border-border bg-card p-3">
            <div class="text-[11px] font-mono uppercase tracking-wider text-muted-foreground">
              Temperature
            </div>
            <div class="font-mono text-lg font-medium text-foreground mt-0.5">
              {{ temperature }}
            </div>
          </div>
          <div v-if="isCharging" class="rounded-xl border border-border bg-card p-3">
            <div class="text-[11px] font-mono uppercase tracking-wider text-muted-foreground">
              Charger
            </div>
            <div class="font-mono text-lg font-medium text-foreground mt-0.5">
              {{ chargerWattage }}
            </div>
          </div>
        </div>

        <!-- Underpowered Warning Chip -->
        <div
          v-if="isUnderpowered && isCharging"
          class="rounded-lg border border-destructive/30 bg-destructive/10 p-3 text-xs text-destructive flex items-center gap-2"
        >
          <span>⚠</span>
          <span
            >Slow charge: 65 W adapter detected on a 90 W laptop. Battery may discharge under heavy
            load.</span
          >
        </div>

        <!-- Prominent Cycle Odometer & Health Banner -->
        <div
          role="button"
          tabindex="0"
          class="rounded-xl border border-border bg-card p-4 transition hover:border-primary/50 hover:bg-muted/20 cursor-pointer active:scale-[0.99] flex items-center justify-between flex-wrap gap-3"
          @click="router.push('/history')"
          @keydown.enter="router.push('/history')"
        >
          <div class="space-y-0.5">
            <div class="text-[10px] uppercase font-mono tracking-wider text-muted-foreground">
              Battery Lifetime Odometer
            </div>
            <div class="flex items-baseline gap-2">
              <span class="font-serif text-2xl font-medium text-foreground">{{
                battery.cycleCount
              }}</span>
              <span class="text-xs text-muted-foreground font-mono">cycles completed</span>
              <span class="text-xs text-primary font-medium ml-1"
                >· {{ battery.healthPercentage }}% Health (Grade A)</span
              >
            </div>
          </div>
          <div
            class="text-xs uppercase font-mono tracking-wider text-muted-foreground hover:text-primary flex items-center gap-1"
          >
            <span>Cycle Ledger</span>
            <span>→</span>
          </div>
        </div>

        <!-- Hour-by-Hour 24h Bar Strip (F01, F22) -->
        <div class="rounded-xl border border-border bg-card p-4 pt-5">
          <div class="flex items-center justify-between mb-3">
            <h2 class="text-xs font-semibold uppercase tracking-wider text-muted-foreground">
              Today, hour by hour
            </h2>
            <span class="text-[11px] font-mono text-muted-foreground">48 slots (30m each)</span>
          </div>

          <!-- 48-slot timeline bar chart -->
          <div class="flex items-end gap-1 h-16 pt-2">
            <div
              v-for="(slot, i) in hourSlots"
              :key="i"
              class="flex-1 rounded-xs transition-all hover:opacity-80"
              :class="slot.charging ? 'bg-primary' : 'bg-muted-foreground/35'"
              :style="{ height: `${slot.height}%` }"
              :title="`Slot ${Math.floor(i / 2)}:${i % 2 === 0 ? '00' : '30'} — ${slot.height}% (${slot.charging ? 'Charging' : 'On Battery'})`"
            />
          </div>

          <!-- Timeline Time Labels -->
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
      </div>
    </div>
  </div>
</template>
