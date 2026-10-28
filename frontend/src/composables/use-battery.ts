import { ref, onMounted, onUnmounted } from "vue";
import * as BatteryService from "../../bindings/changeme/batteryservice.js";
import { Events } from "@wailsio/runtime";

export interface BatteryTelemetry {
  percentage: number;
  status: string;
  isCharging: boolean;
  isAcConnected: boolean;
  powerWatts: number;
  voltageVolts: number;
  healthPercentage: number;
  cycleCount: number;
  designCapacity: number;
  currentCapacity: number;
  fullCapacity: number;
  temperatureC: number;
  timeToEmptyMins: number;
  timeToFullMins: number;
  manufacturer: string;
  modelName: string;
  serialNumber: string;
  technology: string;
  lastUpdated: string;
}

const battery = ref<BatteryTelemetry>({
  percentage: 88,
  status: "Discharging",
  isCharging: false,
  isAcConnected: false,
  powerWatts: 14.2,
  voltageVolts: 11.8,
  healthPercentage: 80,
  cycleCount: 1210,
  designCapacity: 4850,
  currentCapacity: 3415,
  fullCapacity: 3872,
  temperatureC: 29.5,
  timeToEmptyMins: 210,
  timeToFullMins: 0,
  manufacturer: "Hewlett-Packard",
  modelName: "Primary",
  serialNumber: "05494 2021/02/19",
  technology: "Li-ion",
  lastUpdated: "Just now",
});

const isConnected = ref(true);

export function useBattery() {
  async function refresh() {
    try {
      const state = await BatteryService.GetBatteryState();
      if (state) {
        battery.value = {
          percentage: state.percentage ?? 88,
          status: state.status ?? "Discharging",
          isCharging: state.isCharging ?? false,
          isAcConnected: state.isAcConnected ?? false,
          powerWatts: Number((state.powerWatts ?? 0).toFixed(1)),
          voltageVolts: Number((state.voltageVolts ?? 0).toFixed(2)),
          healthPercentage: state.healthPercentage ?? 100,
          cycleCount: state.cycleCount ?? 0,
          designCapacity: state.designCapacity ?? 4850,
          currentCapacity: state.currentCapacity ?? 3415,
          fullCapacity: state.fullCapacity ?? 3872,
          temperatureC: Number((state.temperatureC ?? 29.5).toFixed(1)),
          timeToEmptyMins: state.timeToEmptyMins ?? 0,
          timeToFullMins: state.timeToFullMins ?? 0,
          manufacturer: state.manufacturer || "Primary Battery",
          modelName: state.modelName || "Li-ion",
          serialNumber: state.serialNumber || "05494",
          technology: state.technology || "Li-ion",
          lastUpdated: state.lastUpdated || "Now",
        };
      }
    } catch (err) {
      console.warn("Failed to query BatteryService, running in preview mode:", err);
    }
  }

  let unlisten: (() => void) | null = null;

  onMounted(async () => {
    await refresh();
    try {
      unlisten = Events.On("battery:updated", (event: any) => {
        if (event && event.data) {
          const s = event.data;
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
            lastUpdated: s.lastUpdated ?? "Now",
          };
        }
      });
    } catch (e) {
      // preview environment
    }
  });

  onUnmounted(() => {
    if (unlisten) {
      unlisten();
    }
  });

  return {
    battery,
    isConnected,
    refresh,
    setChargeCeiling: BatteryService.SetChargeCeiling,
  };
}
