package hardware

// DeviceSpecs captures static hardware identity for laptop, battery, and BIOS.
type DeviceSpecs struct {
	// Battery Identity
	BatteryManufacturer string `json:"batteryManufacturer"`
	BatteryModel        string `json:"batteryModel"`
	BatterySerial       string `json:"batterySerial"`
	BatteryTechnology   string `json:"batteryTechnology"`
	DesignCapacityWh    string `json:"designCapacityWh"`
	FullChargeWh        string `json:"fullChargeWh"`
	DesignVoltage       string `json:"designVoltage"`

	// Laptop DMI Identity
	LaptopManufacturer string `json:"laptopManufacturer"`
	LaptopModel        string `json:"laptopModel"`
	LaptopSerial       string `json:"laptopSerial"`
	BiosVersion        string `json:"biosVersion"`
	TotalMemoryRAM     string `json:"totalMemoryRAM"`
	KernelVersion      string `json:"kernelVersion"`
}
