import { Events } from '@wailsio/runtime'
import { ref, onMounted, onUnmounted } from 'vue'

import type { State } from '../../bindings/changeme/core/battery'
import { BatteryService } from '../../bindings/changeme/service'

export interface BatteryTelemetry {
  percentage: number
  status: string
  isCharging: boolean
  isAcConnected: boolean
  powerWatts: number
  voltageVolts: number
  currentAmperes: number
  healthPercentage: number
  cycleCount: number
  designCapacity: number
  currentCapacity: number
  fullCapacity: number
  temperatureC: number
  timeToEmptyMins: number
  timeToFullMins: number
  lastUpdated: string
}

const battery = ref<BatteryTelemetry>({
  percentage: 88,
  status: 'Discharging',
  isCharging: false,
  isAcConnected: false,
  powerWatts: 14.2,
  voltageVolts: 11.8,
  currentAmperes: 1.2,
  healthPercentage: 80,
  cycleCount: 1210,
  designCapacity: 4850,
  currentCapacity: 3415,
  fullCapacity: 3872,
  temperatureC: 29.5,
  timeToEmptyMins: 210,
  timeToFullMins: 0,
  lastUpdated: 'Just now',
})

const isConnected = ref(true)

export function useBattery() {
  async function refresh() {
    try {
      const state: State | null = await BatteryService.GetBatteryState()
      if (state) {
        battery.value = {
          percentage: state.percentage ?? 88,
          status: state.status ?? 'Discharging',
          isCharging: state.isCharging ?? false,
          isAcConnected: state.isAcConnected ?? false,
          powerWatts: Number((state.powerWatts ?? 0).toFixed(1)),
          voltageVolts: Number((state.voltageVolts ?? 0).toFixed(2)),
          currentAmperes: Number((state.currentAmperes ?? 0).toFixed(2)),
          healthPercentage: state.healthPercentage ?? 100,
          cycleCount: state.cycleCount ?? 0,
          designCapacity: state.designCapacity ?? 4850,
          currentCapacity: state.currentCapacity ?? 3415,
          fullCapacity: state.fullCapacity ?? 3872,
          temperatureC: Number((state.temperatureC ?? 29.5).toFixed(1)),
          timeToEmptyMins: state.timeToEmptyMins ?? 0,
          timeToFullMins: state.timeToFullMins ?? 0,
          lastUpdated: new Date().toLocaleTimeString([], {
            hour: '2-digit',
            minute: '2-digit',
            second: '2-digit',
          }),
        }
      }
    } catch (err) {
      console.warn('useBattery query warning:', err)
    }
  }

  let unlisten: (() => void) | null = null

  onMounted(async () => {
    await refresh()
    try {
      unlisten = Events.On('battery:updated', (event: { data: State | null }) => {
        if (event && event.data) {
          const s = event.data
          battery.value = {
            ...battery.value,
            percentage: s.percentage ?? battery.value.percentage,
            status: s.status ?? battery.value.status,
            isCharging: s.isCharging ?? battery.value.isCharging,
            isAcConnected: s.isAcConnected ?? battery.value.isAcConnected,
            powerWatts: Number((s.powerWatts ?? battery.value.powerWatts).toFixed(1)),
            voltageVolts: Number((s.voltageVolts ?? battery.value.voltageVolts).toFixed(2)),
            healthPercentage: s.healthPercentage ?? battery.value.healthPercentage,
            cycleCount: s.cycleCount ?? battery.value.cycleCount,
            lastUpdated: new Date().toLocaleTimeString([], {
              hour: '2-digit',
              minute: '2-digit',
              second: '2-digit',
            }),
          }
        }
      })
    } catch {
      // preview environment
    }
  })

  onUnmounted(() => {
    if (unlisten) {
      unlisten()
    }
  })

  return {
    battery,
    isConnected,
    refresh,
    getHealth: BatteryService.GetBatteryHealth,
    setChargeCeiling: BatteryService.SetChargeCeiling,
  }
}
