<script setup lang="ts">
import { ref } from 'vue'

import SegmentControl from '@/components/common/segment-control.vue'

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

// Charge Sessions demo data (F28)
const sessions = [
  {
    date: 'Today, 2:14 PM',
    duration: '48m',
    startPct: 22,
    endPct: 78,
    avgWatts: '26.4 W',
    peakTemp: '36 °C',
    charger: '65 W USB-PD',
    hot: false,
  },
  {
    date: 'Yesterday, 8:30 PM',
    duration: '1h 12m',
    startPct: 15,
    endPct: 90,
    avgWatts: '31.2 W',
    peakTemp: '42 °C',
    charger: '90 W OEM',
    hot: true,
  },
  {
    date: 'Oct 25, 11:05 AM',
    duration: '35m',
    startPct: 40,
    endPct: 75,
    avgWatts: '24.1 W',
    peakTemp: '33 °C',
    charger: '65 W USB-PD',
    hot: false,
  },
]
</script>

<template>
  <div class="max-w-3xl mx-auto px-6 sm:px-10 py-9">
    <header class="mb-8 flex items-center justify-between flex-wrap gap-4">
      <div>
        <h1 class="font-serif text-4xl tracking-tight text-foreground">History</h1>
        <p class="text-muted-foreground mt-1.5">
          Historical battery drain, charge cycles, and anomalies.
        </p>
      </div>
      <SegmentControl v-model="activeTab" :options="tabOptions" />
    </header>

    <!-- Trends Tab -->
    <div v-if="activeTab === 'trends'" class="space-y-6">
      <div class="flex items-center justify-between">
        <h2 class="text-sm font-semibold uppercase tracking-wider text-muted-foreground">
          Battery Level Over Time
        </h2>
        <SegmentControl v-model="trendRange" :options="rangeOptions" />
      </div>

      <!-- Trend Area Chart SVG -->
      <div class="rounded-xl border border-border bg-card p-6">
        <div class="relative w-full h-44">
          <svg
            viewBox="0 0 600 150"
            preserveAspectRatio="none"
            class="w-full h-full text-primary"
            role="img"
          >
            <!-- Grid lines -->
            <line x1="0" x2="600" y1="0" y2="0" stroke="var(--border)" stroke-dasharray="2 2" />
            <line x1="0" x2="600" y1="75" y2="75" stroke="var(--border)" stroke-dasharray="2 2" />
            <line x1="0" x2="600" y1="150" y2="150" stroke="var(--border)" />

            <!-- Dashed 80% line -->
            <line
              x1="0"
              x2="600"
              y1="30"
              y2="30"
              stroke="var(--destructive)"
              stroke-dasharray="4 4"
              opacity="0.6"
            />

            <!-- Area fill -->
            <path
              d="M0 40 L60 55 L120 45 L180 90 L240 85 L300 25 L360 30 L420 110 L480 95 L540 60 L600 35 L600 150 L0 150 Z"
              fill="currentColor"
              opacity="0.12"
            />
            <!-- Trend line -->
            <path
              d="M0 40 L60 55 L120 45 L180 90 L240 85 L300 25 L360 30 L420 110 L480 95 L540 60 L600 35"
              fill="none"
              stroke="currentColor"
              stroke-width="2.5"
            />
          </svg>

          <!-- Anomaly Dot Flag (F19) -->
          <div
            class="absolute top-[73%] left-[70%] h-3.5 w-3.5 -translate-x-1/2 -translate-y-1/2 rounded-full border-2 border-background bg-destructive"
            title="Sep 14: Drain rate 2x higher than usual workload baseline"
          />
        </div>

        <div
          class="flex items-center justify-between text-xs font-mono text-muted-foreground mt-4 pt-3 border-t border-border"
        >
          <span>{{ trendRange === '24h' ? '24 hours ago' : '10 days ago' }}</span>
          <span class="text-destructive flex items-center gap-1 font-sans">
            <span class="inline-block h-2 w-2 rounded-full bg-destructive" />
            Anomaly detected (rapid drain spike)
          </span>
          <span>Now</span>
        </div>
      </div>
    </div>

    <!-- Sessions Tab (F28) -->
    <div v-else class="space-y-4">
      <div class="rounded-xl border border-border bg-card overflow-hidden">
        <div
          class="grid grid-cols-[1.5fr_1fr_1fr_1fr] border-b border-border bg-muted/40 px-4 py-2.5 text-xs font-semibold text-muted-foreground uppercase tracking-wider"
        >
          <span>Charge Session</span>
          <span>Duration</span>
          <span>Power & Temp</span>
          <span>Adapter</span>
        </div>

        <div class="divide-y divide-border">
          <div
            v-for="(s, idx) in sessions"
            :key="idx"
            class="grid grid-cols-[1.5fr_1fr_1fr_1fr] items-center px-4 py-3.5 text-sm transition hover:bg-muted/30"
          >
            <div>
              <div class="font-medium text-foreground">{{ s.date }}</div>
              <div class="text-xs font-mono text-muted-foreground mt-0.5">
                {{ s.startPct }}% ➔ {{ s.endPct }}%
              </div>
            </div>

            <div class="font-mono text-xs text-muted-foreground">
              {{ s.duration }}
            </div>

            <div class="text-xs">
              <div class="font-mono font-medium">{{ s.avgWatts }}</div>
              <div class="flex items-center gap-1 mt-0.5">
                <span class="font-mono text-muted-foreground">{{ s.peakTemp }}</span>
                <span
                  v-if="s.hot"
                  class="rounded bg-destructive/10 px-1 py-0.2 text-[10px] font-semibold text-destructive"
                >
                  Hot
                </span>
              </div>
            </div>

            <div class="text-xs text-muted-foreground font-mono">
              {{ s.charger }}
            </div>
          </div>
        </div>
      </div>
    </div>
  </div>
</template>
