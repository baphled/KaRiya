package capture_event

import (
	"github.com/baphled/kariya/internal/cli/intents"
	"github.com/baphled/kariya/internal/cli/service"
	"github.com/baphled/kariya/internal/cli/types"
	"github.com/baphled/kariya/internal/domain/career"
	careerservice "github.com/baphled/kariya/internal/service/career"
)

// CaptureStrategy defines the event capture approach and controls which form fields
// are presented during the capture workflow. This type alias re-exports
// types.CaptureStrategy so that the capture_event package can reference it without
// requiring callers to import the types package.
type CaptureStrategy = types.CaptureStrategy

const (
	// StrategyQuick captures only required fields (event text), date defaults to today.
	StrategyQuick = types.StrategyQuick

	// StrategyManual shows all fields with optional field toggle.
	StrategyManual = types.StrategyManual
)

// IntentContext is the minimal context passed to the CaptureEvent intent.
// It contains only what's necessary to start the intent.
type IntentContext struct {
	// CaptureStrategy determines how the event is captured (manual, quick, enriched).
	CaptureStrategy string

	// PreviousEvent is an optional existing event to edit (nil for new capture).
	PreviousEvent *career.Event

	// Metadata is the initial metadata for the event (may be empty).
	Metadata map[string]string

	// CLIEventService is the service for CLI operations (for form submission).
	CLIEventService *service.CLIEventService

	// CareerService is the domain service for enrichment operations.
	// Used for suggesting bursts and extracting facts in enriched mode.
	CareerService *careerservice.Service
}

// Validate ensures the context is complete.
func (c *IntentContext) Validate() error {
	if c.CaptureStrategy == "" {
		return intents.NewFailedResult[*IntentContext](
			"invalid_context",
			"CaptureStrategy must not be empty",
			nil,
		).Error
	}
	if c.Metadata == nil {
		c.Metadata = make(map[string]string)
	}
	return nil
}
