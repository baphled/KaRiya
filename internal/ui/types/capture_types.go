package types

// CaptureStrategy defines how the event should be captured.
type CaptureStrategy string

const (
	// StrategyQuick captures only required fields (event text), date defaults to today.
	StrategyQuick CaptureStrategy = "quick"

	// StrategyManual shows all fields with optional field toggle.
	StrategyManual CaptureStrategy = "manual"
)
