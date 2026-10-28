<script setup lang="ts">
import { ref, onMounted, computed } from 'vue'

import type { DeviceSpecs } from '../../../bindings/changeme/core/hardware'
import { DeviceService } from '../../../bindings/changeme/service'

import { useBattery } from '@/composables/use-battery'

const { battery } = useBattery()
const specs = ref<DeviceSpecs | null>(null)

const includeHistory = ref(true)
const includeDeviceSpecs = ref(true)
const includeSerials = ref(false) // Masked by default (DEC-17)

const currentDate = computed(() =>
  new Date().toLocaleDateString(undefined, {
    month: 'short',
    day: 'numeric',
    year: 'numeric',
  }),
)

const designMwh = computed(() => {
  const mwh = Math.round(battery.value.designCapacity * 11.55)
  return mwh.toLocaleString()
})

onMounted(async () => {
  try {
    const s = await DeviceService.GetDeviceSpecs()
    if (s) {
      specs.value = s
    }
  } catch (err) {
    console.debug('Failed to query device specs in report:', err)
  }
})
</script>

<template>
  <div class="max-w-3xl mx-auto px-6 sm:px-10 py-9">
    <header class="mb-8 flex items-center justify-between flex-wrap gap-4">
      <div>
        <h1 class="font-serif text-4xl tracking-tight text-foreground">Health Reports</h1>
        <p class="text-muted-foreground mt-1.5">
          Export inspection certificates for resale, warranty, or audit.
        </p>
      </div>

      <div class="flex items-center gap-2">
        <button
          type="button"
          class="rounded-lg bg-primary px-4 py-2 text-xs font-medium text-primary-foreground transition hover:opacity-90 shadow-sm"
        >
          Export PDF
        </button>
        <button
          type="button"
          class="rounded-lg border border-border bg-card px-4 py-2 text-xs font-medium text-foreground transition hover:bg-muted"
        >
          Export CSV
        </button>
      </div>
    </header>

    <div class="grid gap-6 md:grid-cols-[220px_1fr]">
      <!-- Report Inclusions Panel -->
      <div class="space-y-4">
        <div class="rounded-xl border border-border bg-card p-4 space-y-3">
          <div class="text-xs uppercase font-semibold text-muted-foreground tracking-wider mb-2">
            Inclusions
          </div>

          <label class="flex items-center gap-2 text-xs text-foreground cursor-pointer">
            <input
              v-model="includeHistory"
              type="checkbox"
              class="rounded border-border accent-primary"
            />
            <span>Charge history & cycles</span>
          </label>

          <label class="flex items-center gap-2 text-xs text-foreground cursor-pointer">
            <input
              v-model="includeDeviceSpecs"
              type="checkbox"
              class="rounded border-border accent-primary"
            />
            <span>Laptop & adapter details</span>
          </label>

          <label class="flex items-center gap-2 text-xs text-foreground cursor-pointer">
            <input
              v-model="includeSerials"
              type="checkbox"
              class="rounded border-border accent-primary"
            />
            <span>Include serial numbers</span>
          </label>
        </div>

        <div class="rounded-lg bg-muted/40 p-3 text-xs text-muted-foreground">
          🔒 Serials are omitted by default to ensure privacy during online marketplace resale
          sharing.
        </div>
      </div>

      <!-- Live Printable Document Preview Card -->
      <div class="rounded-xl border border-border bg-card p-6 shadow-sm font-sans space-y-6">
        <div class="flex items-center justify-between border-b border-border pb-4">
          <div class="flex items-center gap-2">
            <img src="/antlion-mark.svg" class="h-5 w-5" alt="Logo" />
            <span class="font-serif text-lg font-semibold"
              >Antlion Battery Inspection Certificate</span
            >
          </div>
          <span class="font-mono text-xs text-muted-foreground">Generated: {{ currentDate }}</span>
        </div>

        <!-- Certificate Summary -->
        <div class="grid grid-cols-2 gap-4">
          <div>
            <div class="text-xs text-muted-foreground">Machine</div>
            <div class="font-medium text-sm text-foreground">
              {{ specs?.laptopModel || 'Linux Laptop' }}
            </div>
          </div>
          <div>
            <div class="text-xs text-muted-foreground">Health Condition</div>
            <div class="font-medium text-sm text-primary">
              {{ battery.healthPercentage >= 80 ? 'Normal · Grade A' : 'Service Recommended' }} ({{
                battery.healthPercentage
              }}%)
            </div>
          </div>
          <div>
            <div class="text-xs text-muted-foreground">Cycle Count</div>
            <div class="font-mono text-sm text-foreground">{{ battery.cycleCount }} cycles</div>
          </div>
          <div>
            <div class="text-xs text-muted-foreground">Design Capacity</div>
            <div class="font-mono text-sm text-foreground">{{ designMwh }} mWh</div>
          </div>
        </div>

        <div
          class="rounded-lg bg-muted/60 p-4 text-xs font-mono text-muted-foreground leading-relaxed"
        >
          This battery operates at {{ battery.healthPercentage }}% retained capacity relative to OEM
          factory specifications. Physical fuel gauge logs {{ battery.cycleCount }} lifetime cycles.
          Chemical degradation rate is measured across historical sysfs telemetry.
        </div>
      </div>
    </div>
  </div>
</template>
