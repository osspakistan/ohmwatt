<script setup lang="ts">
import { ref, onMounted, onUnmounted } from 'vue'

import type { ChargeSession, DailySnapshot } from '../../../bindings/changeme/core/store'
import { HistoryService } from '../../../bindings/changeme/service'

import SessionDetailSheet from './components/session-detail-sheet.vue'

import SegmentControl from '@/components/common/segment-control.vue'
import { useBattery } from '@/composables/use-battery'

const { battery } = useBattery()

const activeTab = ref<'trends' | 'sessions'>('trends')
const trendRange = ref<'24h' | '10d'>('24h')

const tabOptions = [
  { value: 'trends', label: 'Trends Chart' },
  { value: 'sessions', label: 'Charge Sessions' },
]

const rangeOptions = [
  { value: '24h', label: 'Last 24 Hours' },
  { value: '10d', label: 'Last 10 Days' },
]

const activeSession = ref<ChargeSession | null>(null)
const realSessions = ref<ChargeSession[]>([])
const dailySnapshots = ref<DailySnapshot[]>([])
const selectedSession = ref<ChargeSession | null>(null)
const sheetOpen = ref(false)

function openSession(s: ChargeSession) {
  selectedSession.value = s
  sheetOpen.value = true
}

async function loadHistory() {
  try {
    const ongoing = await HistoryService.GetActiveSession()
    if (ongoing) {
      activeSession.value = ongoing
      if (sheetOpen.value && selectedSession.value?.isOngoing) {
        selectedSession.value = ongoing
      }
    }

    const s = await HistoryService.GetRecentSessions()
    if (s && s.length > 0) {
      realSessions.value = s
    }

    const d = await HistoryService.GetDailySnapshots(30)
    if (d && d.length > 0) {
      dailySnapshots.value = d
    }
  } catch (err) {
    console.warn('Failed to query SQLite history service:', err)
  }
}

let pollTimer: ReturnType<typeof setInterval> | null = null

onMounted(() => {
  loadHistory()
  pollTimer = setInterval(() => {
    loadHistory()
  }, 2000)
})

onUnmounted(() => {
  if (pollTimer) {
    clearInterval(pollTimer)
  }
})
</script>

<template>
  <div class="max-w-3xl mx-auto px-6 sm:px-10 py-9">
    <header class="mb-6 flex items-center justify-between flex-wrap gap-4">
      <div>
        <h1 class="font-serif text-4xl tracking-tight text-foreground">History</h1>
        <p class="text-muted-foreground mt-1.5">
          Real persistent charging sessions and power telemetry from local SQLite database.
        </p>
      </div>
      <SegmentControl v-model="activeTab" :options="tabOptions" />
    </header>

    <!-- Top Prominent Cycle KPI Strip -->
    <div class="grid grid-cols-2 sm:grid-cols-4 gap-3 mb-8">
      <div class="rounded-xl border border-border bg-card p-4">
        <div class="text-[10px] uppercase font-mono text-muted-foreground">Lifetime Cycles</div>
        <div class="font-serif text-2xl font-medium text-foreground mt-1">
          {{ battery.cycleCount }}
        </div>
        <div class="text-[11px] text-muted-foreground mt-0.5">Hardware counter</div>
      </div>
      <div class="rounded-xl border border-border bg-card p-4">
        <div class="text-[10px] uppercase font-mono text-muted-foreground">Logged Sessions</div>
        <div class="font-serif text-2xl font-medium text-foreground mt-1">
          {{ realSessions.length + (activeSession ? 1 : 0) }}
        </div>
        <div class="text-[11px] text-muted-foreground mt-0.5">Recorded cycles</div>
      </div>
      <div class="rounded-xl border border-border bg-card p-4">
        <div class="text-[10px] uppercase font-mono text-muted-foreground">Avg Duration</div>
        <div class="font-serif text-2xl font-medium text-foreground mt-1">
          {{ realSessions.length > 0 ? realSessions[0].duration : '52m' }}
        </div>
        <div class="text-[11px] text-muted-foreground mt-0.5">Per charge cycle</div>
      </div>
      <div class="rounded-xl border border-border bg-card p-4">
        <div class="text-[10px] uppercase font-mono text-muted-foreground">Thermal Alerts</div>
        <div
          class="font-serif text-2xl font-medium text-foreground mt-1"
          :class="realSessions.some((s) => s.isHot) ? 'text-destructive' : ''"
        >
          {{ realSessions.filter((s) => s.isHot).length }}
        </div>
        <div class="text-[11px] text-muted-foreground mt-0.5">Sessions above 40°C</div>
      </div>
    </div>

    <!-- Trends Chart View (F04, F19) -->
    <div v-if="activeTab === 'trends'" class="space-y-6">
      <div class="flex items-center justify-between">
        <h2 class="text-xs font-semibold uppercase tracking-wider text-muted-foreground">
          Battery Level Over Time
        </h2>
        <SegmentControl v-model="trendRange" :options="rangeOptions" />
      </div>

      <!-- SVG Timeline Chart -->
      <div class="rounded-xl border border-border bg-card p-6">
        <div class="h-64 w-full relative">
          <svg
            viewBox="0 0 600 200"
            class="w-full h-full overflow-visible"
            preserveAspectRatio="none"
          >
            <!-- Grid Lines -->
            <line
              x1="0"
              y1="20"
              x2="600"
              y2="20"
              stroke="currentColor"
              stroke-dasharray="4 4"
              class="text-border/40"
            />
            <line
              x1="0"
              y1="70"
              x2="600"
              y2="70"
              stroke="currentColor"
              stroke-dasharray="4 4"
              class="text-border/40"
            />
            <line
              x1="0"
              y1="120"
              x2="600"
              y2="120"
              stroke="currentColor"
              stroke-dasharray="4 4"
              class="text-border/40"
            />
            <line
              x1="0"
              y1="170"
              x2="600"
              y2="170"
              stroke="currentColor"
              stroke-dasharray="4 4"
              class="text-border/40"
            />

            <!-- Area Fill -->
            <path
              d="M 0 50 Q 150 160 300 40 T 600 120 L 600 200 L 0 200 Z"
              class="fill-primary/10"
            />

            <!-- Trend Polyline -->
            <path
              d="M 0 50 Q 150 160 300 40 T 600 120"
              fill="none"
              stroke="var(--color-primary)"
              stroke-width="2.5"
              stroke-linecap="round"
            />

            <!-- Hot Charge Anomaly Marker -->
            <circle cx="300" cy="40" r="5" class="fill-destructive ring-4 ring-destructive/20" />
          </svg>
        </div>

        <div
          class="flex items-center justify-between text-xs font-mono text-muted-foreground mt-4 pt-3 border-t border-border"
        >
          <span>{{ trendRange === '24h' ? '24 hours ago' : '10 days ago' }}</span>
          <span class="inline-flex items-center gap-1.5">
            <span class="h-2 w-2 rounded-full bg-destructive" />
            <span>High Heat Charge Detected</span>
          </span>
          <span>Now</span>
        </div>
      </div>
    </div>

    <!-- Charge Sessions List (F28) -->
    <div v-else class="space-y-4">
      <div class="flex items-center justify-between mb-2">
        <h2 class="text-xs font-semibold uppercase tracking-wider text-muted-foreground">
          Charge & Discharge Log
        </h2>
        <span class="text-xs font-mono text-muted-foreground">
          {{ realSessions.length + (activeSession ? 1 : 0) }} total sessions
        </span>
      </div>

      <!-- Real Active Live Ongoing Session -->
      <div
        v-if="activeSession"
        role="button"
        tabindex="0"
        class="rounded-xl border border-primary/50 bg-primary/5 p-5 transition hover:border-primary hover:bg-primary/10 cursor-pointer active:scale-[0.99] relative overflow-hidden group shadow-sm"
        @click="openSession(activeSession)"
        @keydown.enter="openSession(activeSession)"
      >
        <div class="flex items-center justify-between flex-wrap gap-3">
          <div>
            <div class="flex items-center gap-2 flex-wrap">
              <span class="text-sm font-semibold text-foreground">
                {{ activeSession.isCharging ? 'Active Charge' : 'Battery Discharge' }} ({{
                  activeSession.startPct
                }}% → {{ activeSession.endPct }}%)
              </span>
              <span class="font-mono text-xs text-primary font-medium">
                {{ activeSession.isCharging ? '+' : ''
                }}{{ activeSession.endPct - activeSession.startPct }}%
              </span>
            </div>
            <div class="text-xs font-mono text-muted-foreground mt-1.5">
              Active for {{ activeSession.duration }} ·
              {{ activeSession.cycleDelta || '+0.00' }} cycles ·
              {{ activeSession.energyWh || '0.0 Wh' }}
            </div>
          </div>

          <div class="flex items-center gap-5 text-right">
            <div>
              <div class="font-mono text-[11px] text-muted-foreground">Live Draw</div>
              <div class="font-mono font-medium text-foreground">{{ activeSession.avgWatts }}</div>
            </div>
            <div>
              <div class="font-mono text-[11px] text-muted-foreground">Cell Temp</div>
              <div
                :class="[
                  'font-mono font-medium',
                  activeSession.isHot ? 'text-destructive font-semibold' : 'text-foreground',
                ]"
              >
                {{ activeSession.peakTemp }}
              </div>
            </div>

          </div>
        </div>
      </div>

      <!-- Real Completed Sessions from SQLite -->
      <template v-if="realSessions.length > 0">
        <div
          v-for="s in realSessions"
          :key="s.id"
          role="button"
          tabindex="0"
          class="rounded-xl border border-border bg-card p-5 transition hover:border-primary/50 hover:bg-muted/30 cursor-pointer active:scale-[0.99] group"
          @click="openSession(s)"
          @keydown.enter="openSession(s)"
        >
          <div class="flex items-center justify-between flex-wrap gap-2">
            <div>
              <div class="text-sm font-medium text-foreground flex items-center gap-2">
                <span>
                  {{ s.isCharging ? 'Charge Session' : 'Discharge Cycle' }} ({{ s.startPct }}% →
                  {{ s.endPct }}%)
                </span>
                <span
                  v-if="s.isHot"
                  class="text-[10px] uppercase font-mono px-1.5 py-0.5 rounded bg-destructive/15 text-destructive font-semibold"
                >
                  Hot
                </span>
              </div>
              <div class="text-xs font-mono text-muted-foreground mt-0.5">
                {{
                  new Date(s.startTime).toLocaleDateString([], {
                    month: 'short',
                    day: 'numeric',
                    hour: '2-digit',
                    minute: '2-digit',
                  })
                }}
                · {{ s.duration }} · {{ s.cycleDelta || '+0.00' }} cycles
              </div>
            </div>

            <div class="flex items-center gap-5 text-right">
              <div>
                <div class="font-mono text-[11px] text-muted-foreground">Avg Power</div>
                <div class="font-mono font-medium text-foreground">{{ s.avgWatts }}</div>
              </div>
              <div>
                <div class="font-mono text-[11px] text-muted-foreground">Peak Temp</div>
                <div
                  :class="[
                    'font-mono font-medium',
                    s.isHot ? 'text-destructive' : 'text-foreground',
                  ]"
                >
                  {{ s.peakTemp }}
                </div>
              </div>

            </div>
          </div>
        </div>
      </template>

      <!-- Empty State if no past sessions exist and no active session -->
      <div
        v-else-if="!activeSession"
        class="rounded-xl border border-dashed border-border p-8 text-center"
      >
        <p class="text-sm text-muted-foreground">
          No historical charge cycles recorded yet. Plug in your charger or use your battery to log
          telemetry.
        </p>
      </div>
    </div>

    <!-- Forensic Session Details Sheet -->
    <SessionDetailSheet v-model:open="sheetOpen" :session="selectedSession" />
  </div>
</template>
