package captureevent

import (
	"github.com/baphled/kariya/internal/cli/service"
	"github.com/baphled/kariya/internal/domain/career"
	careerservice "github.com/baphled/kariya/internal/service/career"
	"github.com/baphled/kariya/internal/service/career/skillinference"
	"github.com/baphled/kariya/internal/tui/intents"
	"github.com/baphled/kariya/internal/ui/types"
)

// CaptureStrategy defines the event capture approach and controls which form
// fields are presented during the capture workflow.
//
// This type alias re-exports types.CaptureStrategy so callers do not need to
// import the types package directly.
type CaptureStrategy = types.CaptureStrategy

const (
	// StrategyQuick captures only the event text; date defaults to today.
	StrategyQuick = types.StrategyQuick

	// StrategyManual presents all fields with an optional field toggle.
	StrategyManual = types.StrategyManual
)

// IntentValidator holds the input parameters required to start the CaptureEvent intent.
//
// Expected:
//   - CaptureStrategy must not be empty
//   - CLIEventService should be non-nil for form submission
//   - CareerService should be non-nil for enrichment operations
//
// PreviousEvent may be nil (new capture) or set (edit existing event).
type IntentValidator struct {
	// CaptureStrategy determines how the event is captured (e.g. "quick", "manual").
	CaptureStrategy string

	// PreviousEvent is the existing event to edit, or nil for a new capture.
	PreviousEvent *career.Event

	// Metadata holds initial metadata key-value pairs for the event.
	Metadata map[string]string

	// CLIEventService provides CLI-level operations needed by the capture form.
	CLIEventService *service.CLIEventService

	// CareerService provides domain operations for enrichment (burst detection,
	// fact extraction) during the review phase.
	CareerService *careerservice.Service

	// SkillInferenceService provides skill detection capabilities during review.
	SkillInferenceService skillinference.SkillInferenceService
}

// Validate checks that all required fields are present.
//
// Expected:
//   - CaptureStrategy must not be empty
//
// Returns:
//   - nil on success
//   - error if CaptureStrategy is empty
//
// Side effects:
//   - Initialises Metadata to an empty map if nil.
func (c *IntentValidator) Validate() error {
	if c.CaptureStrategy == "" {
		return intents.NewFailedResult[*IntentValidator](
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
