<script setup lang="ts">
import { ref } from 'vue'

const includeHistory = ref(true)
const includeDeviceSpecs = ref(true)
const includeSerials = ref(false) // Masked by default (DEC-17)
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
          <span class="font-mono text-xs text-muted-foreground">Generated: Oct 27, 2026</span>
        </div>

        <!-- Certificate Summary -->
        <div class="grid grid-cols-2 gap-4">
          <div>
            <div class="text-xs text-muted-foreground">Machine</div>
            <div class="font-medium text-sm text-foreground">HP EliteBook 850 G7</div>
          </div>
          <div>
            <div class="text-xs text-muted-foreground">Health Condition</div>
            <div class="font-medium text-sm text-primary">Normal · Grade A (91%)</div>
          </div>
          <div>
            <div class="text-xs text-muted-foreground">Cycle Count</div>
            <div class="font-mono text-sm text-foreground">212 cycles</div>
          </div>
          <div>
            <div class="text-xs text-muted-foreground">Design Capacity</div>
            <div class="font-mono text-sm text-foreground">53,010 mWh</div>
          </div>
        </div>

        <div
          class="rounded-lg bg-muted/60 p-4 text-xs font-mono text-muted-foreground leading-relaxed"
        >
          This battery operates at 91.0% retained capacity relative to OEM factory specifications.
          Chemical degradation rate is measured at ~0.04% per cycle. No cell imbalances or thermal
          throttling anomalies detected.
        </div>
      </div>
    </div>
  </div>
</template>
