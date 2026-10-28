<script setup lang="ts">
import type { ChargeSession } from '../../../../bindings/changeme/core/store'

import {
  Sheet,
  SheetContent,
  SheetHeader,
  SheetTitle,
  SheetDescription,
} from '@/components/ui/sheet'

defineProps<{
  open: boolean
  session: ChargeSession | null
}>()

const emit = defineEmits<{
  (e: 'update:open', val: boolean): void
}>()

function calculateCycleDelta(start: number, end: number): string {
  const diff = Math.abs(end - start)
  return (diff / 100).toFixed(2)
}
</script>

<template>
  <Sheet :open="open" @update:open="emit('update:open', $event)">
    <SheetContent class="w-full sm:max-w-md bg-card border-l border-border p-6 overflow-y-auto">
      <SheetHeader v-if="session" class="space-y-1 text-left">
        <div class="inline-flex items-center gap-2">
          <span
            :class="[
              'text-[10px] uppercase font-mono font-semibold px-2 py-0.5 rounded tracking-wide',
              session.isCharging
                ? 'bg-primary/15 text-primary'
                : 'bg-muted-foreground/15 text-muted-foreground',
            ]"
          >
            {{ session.isCharging ? 'Charge Session' : 'Discharge Cycle' }}
          </span>
          <span
            v-if="session.isHot"
            class="text-[10px] uppercase font-mono font-semibold px-2 py-0.5 rounded bg-destructive/15 text-destructive tracking-wide"
          >
            Hot Cell Warning
          </span>
        </div>
        <SheetTitle class="font-serif text-2xl font-normal text-foreground">
          {{ session.startPct }}% → {{ session.endPct }}%
        </SheetTitle>
        <SheetDescription class="text-xs text-muted-foreground">
          {{
            new Date(session.startTime).toLocaleDateString([], {
              weekday: 'short',
              month: 'short',
              day: 'numeric',
              year: 'numeric',
            })
          }}
          ·
          {{
            new Date(session.startTime).toLocaleTimeString([], {
              hour: '2-digit',
              minute: '2-digit',
            })
          }}
          to
          {{
            new Date(session.endTime).toLocaleTimeString([], { hour: '2-digit', minute: '2-digit' })
          }}
        </SheetDescription>
      </SheetHeader>

      <div v-if="session" class="mt-6 space-y-6">
        <!-- 3-KPI Metrics Grid -->
        <div
          class="grid grid-cols-3 gap-2 rounded-lg border border-border bg-muted/30 p-3 text-center"
        >
          <div>
            <div class="text-[10px] uppercase font-mono text-muted-foreground">Duration</div>
            <div class="font-mono text-sm font-medium text-foreground mt-0.5">
              {{ session.duration }}
            </div>
          </div>
          <div>
            <div class="text-[10px] uppercase font-mono text-muted-foreground">Cycle Fraction</div>
            <div class="font-mono text-sm font-medium text-foreground mt-0.5">
              +{{ calculateCycleDelta(session.startPct, session.endPct) }}
            </div>
          </div>
          <div>
            <div class="text-[10px] uppercase font-mono text-muted-foreground">Peak Temp</div>
            <div
              :class="[
                'font-mono text-sm font-medium mt-0.5',
                session.isHot ? 'text-destructive' : 'text-foreground',
              ]"
            >
              {{ session.peakTemp }}
            </div>
          </div>
        </div>

        <!-- Telemetry Details Breakdown -->
        <div class="space-y-3">
          <h4 class="text-xs uppercase font-mono tracking-wider text-muted-foreground">
            Session Diagnostics
          </h4>
          <div class="divide-y divide-border rounded-lg border border-border bg-card p-4 text-xs">
            <div class="flex items-center justify-between py-2 first:pt-0">
              <span class="text-muted-foreground">Average Power Draw</span>
              <span class="font-mono font-medium text-foreground">{{ session.avgWatts }}</span>
            </div>
            <div class="flex items-center justify-between py-2">
              <span class="text-muted-foreground">Thermal Condition</span>
              <span
                class="font-mono font-medium"
                :class="session.isHot ? 'text-destructive' : 'text-foreground'"
              >
                {{ session.isHot ? 'High Stress (>40°C)' : 'Nominal (<38°C)' }}
              </span>
            </div>
            <div class="flex items-center justify-between py-2">
              <span class="text-muted-foreground">Energy Delivery Mode</span>
              <span class="font-mono font-medium text-foreground">
                {{ session.isCharging ? 'AC / USB-PD Ingest' : 'Battery Active Drain' }}
              </span>
            </div>
            <div class="flex items-center justify-between py-2 last:pb-0">
              <span class="text-muted-foreground">Hardware Accumulation</span>
              <span class="font-mono font-medium text-foreground"> Logged to SQLite Database </span>
            </div>
          </div>
        </div>

        <!-- Contextual Recommendation -->
        <div
          class="rounded-lg border border-border bg-muted/20 p-4 text-xs text-muted-foreground leading-relaxed"
        >
          <p v-if="session.isHot" class="text-destructive">
            ⚠ This session exceeded 40°C. Consider avoiding heavy compilation or gaming while
            charging on soft surfaces to protect cell longevity.
          </p>
          <p v-else>
            ✓ Optimal charge session. Keeping temperature below 38°C preserves lithium-ion
            electrolyte and minimizes irreversible capacity loss.
          </p>
        </div>
      </div>
    </SheetContent>
  </Sheet>
</template>
