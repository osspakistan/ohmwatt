package power

// AdapterInfo represents the active AC charger and USB-PD negotiation state.
type AdapterInfo struct {
	IsConnected      bool    `json:"isConnected"`
	Type             string  `json:"type"`             // "USB Power Delivery", "AC Adapter", "Mains"
	NegotiatedWatts  float64 `json:"negotiatedWatts"`  // e.g. 65.0 W
	CurrentDrawWatts float64 `json:"currentDrawWatts"` // e.g. 32.4 W
	VoltageVolts     float64 `json:"voltageVolts"`     // e.g. 20.0 V
	Amperes          float64 `json:"amperes"`          // e.g. 3.25 A
	IsUnderpowered   bool    `json:"isUnderpowered"`   // true if charger delivers less than laptop requires
}
