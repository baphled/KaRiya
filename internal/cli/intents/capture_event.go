package intents

import (
	"github.com/baphled/kariya/internal/cli/models"
	"github.com/baphled/kariya/internal/cli/service"
	"github.com/baphled/kariya/internal/domain/career"
	careerservice "github.com/baphled/kariya/internal/service/career"
)

// CaptureStrategy defines how the event should be captured.
type CaptureStrategy string

const (
	// StrategyQuick captures only required fields (event text), date defaults to today
	StrategyQuick CaptureStrategy = "quick"

	// StrategyManual shows all fields with optional field toggle
	StrategyManual CaptureStrategy = "manual"
)

// CaptureEventContext is the minimal context passed to CaptureEvent intent.
// It contains only what's necessary to start the intent.
type CaptureEventContext struct {
	// CaptureStrategy determines how the event is captured (manual, quick, enriched).
	CaptureStrategy string

	// PreviousEvent is an optional existing event to edit (nil for new capture).
	PreviousEvent *career.CareerEvent

	// Metadata is the initial metadata for the event (may be empty).
	Metadata map[string]string

	// CLIEventService is the service for CLI operations (for form submission).
	CLIEventService *service.CLIEventService

	// CareerService is the domain service for enrichment operations.
	// Used for suggesting bursts and extracting facts in enriched mode.
	CareerService *careerservice.Service
}

// Validate ensures the context is complete.
func (c *CaptureEventContext) Validate() error {
	if c.CaptureStrategy == "" {
		return NewFailedResult[*CaptureEventContext](
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

// CaptureEventResult is the output of a successful CaptureEvent intent.
type CaptureEventResult struct {
	// Event is the captured or edited event.
	Event *career.CareerEvent

	// Bursts are the inferred bursts from the event (may be empty).
	Bursts []*career.Burst

	// Facts are the inferred facts from the event (may be empty).
	Facts []*career.Fact

	// AcceptedFields tracks which fields were accepted by the user.
	// Used for partial results to understand what was rejected.
	AcceptedFields map[string]bool

	// RejectedFields tracks which inferred fields were rejected.
	RejectedFields map[string]string // field -> rejection reason
}

// ReviewInferredEventState represents the state of the ReviewInferredEvent sub-flow.
type ReviewInferredEventState struct {
	// Event being reviewed.
	Event *career.CareerEvent

	// InferredBursts are the bursts suggested by enrichment.
	InferredBursts []*career.Burst

	// InferredFacts are the facts suggested by enrichment.
	InferredFacts []*career.Fact

	// EditingMode indicates which sub-flow is active (EditMetadata, EditBursts, EditFacts, or empty for review).
	EditingMode string

	// EditingIndex is the index of the item being edited (for bursts/facts).
	EditingIndex int

	// AcceptedBursts tracks user-accepted bursts.
	AcceptedBursts []*career.Burst

	// AcceptedFacts tracks user-accepted facts.
	AcceptedFacts []*career.Fact

	// RejectedItems tracks items rejected by the user.
	RejectedItems map[string]string // item_id -> rejection reason
}

// EditingModes for ReviewInferredEvent sub-flows.
const (
	EditingModeNone     = ""
	EditingModeMetadata = "metadata"
	EditingModeBursts   = "bursts"
	EditingModeFacts    = "facts"
)

// CaptureEventModel represents the local state of the CaptureEvent intent.
// This is the ONLY mutable state owned by the intent.
type CaptureEventModel struct {
	// context is the input context passed to the intent.
	context *CaptureEventContext

	// currentState tracks which sub-flow is active.
	currentState string // ChooseCaptureStrategy, CaptureForm, ReviewInferredEvent, Submit

	// captureForm is the form for capturing event details.
	captureForm *models.FormModel

	// reviewState is the state of the ReviewInferredEvent sub-flow.
	reviewState *ReviewInferredEventState

	// result is the final result to be returned.
	result *CaptureEventResult

	// error tracks any errors during the intent.
	error *IntentError

	// selectedStrategyIndex tracks the selected strategy in the choose strategy screen (0=Quick, 1=Manual)
	selectedStrategyIndex int

	// strategy is the chosen capture strategy
	strategy CaptureStrategy

	// showOptionalFields tracks whether optional fields are visible (for manual mode toggle persistence)
	showOptionalFields bool
}

// CaptureEventStates for navigation.
const (
	CaptureStateChooseStrategy = "choose_strategy"
	CaptureStateForm           = "form"
	CaptureStateReview         = "review"
	CaptureStateSubmit         = "submit"
)
