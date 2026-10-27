package hibernate

// CheckItem represents a single step in the zero-drain hibernation readiness checklist.
type CheckItem struct {
	Title   string `json:"title"`
	Passed  bool   `json:"passed"`
	Detail  string `json:"detail"`
	Command string `json:"command,omitempty"`
}

// Diagnostic represents the system hibernation readiness analysis.
type Diagnostic struct {
	IsReady        bool        `json:"isReady"`
	RamSizeGiB     float64     `json:"ramSizeGiB"`
	SwapSizeGiB    float64     `json:"swapSizeGiB"`
	HasSufficientSwap bool     `json:"hasSufficientSwap"`
	ResumeConfigured  bool     `json:"resumeConfigured"`
	Checks         []CheckItem `json:"checks"`
}
