package intents

import (
	"errors"

	"github.com/baphled/kariya/internal/cli/components"
	"github.com/baphled/kariya/internal/domain/career"
)

// GenerateCVState represents the state of the GenerateCV intent.
type GenerateCVState string

const (
	// GenerateCVStateSelectProfile - User selects or creates a profile.
	GenerateCVStateSelectProfile GenerateCVState = "select_profile"

	// GenerateCVStateSelectAudience - User selects target audience(s).
	GenerateCVStateSelectAudience GenerateCVState = "select_audience"

	// GenerateCVStatePreview - User previews the generated CV.
	GenerateCVStatePreview GenerateCVState = "preview"

	// GenerateCVStateReview - User reviews and edits the CV.
	GenerateCVStateReview GenerateCVState = "review"

	// GenerateCVStateConfirm - User confirms the CV generation.
	GenerateCVStateConfirm GenerateCVState = "confirm"
)

// GenerateCVContext is the input context passed to the GenerateCV intent.
type GenerateCVContext struct {
	// AvailableProfiles are the profiles the user can choose from.
	AvailableProfiles []*CVProfile

	// Events are the career events to use for CV generation.
	Events []*career.CareerEvent

	// Facts are the career facts to use for CV generation.
	Facts []*career.Fact

	// DefaultProfile is the profile to select by default.
	DefaultProfile *CVProfile
}

// Validate checks if the context is valid.
func (ctx *GenerateCVContext) Validate() error {
	if ctx == nil {
		return errors.New("GenerateCVContext cannot be nil")
	}

	if len(ctx.AvailableProfiles) == 0 {
		return errors.New("GenerateCVContext must have at least one available profile")
	}

	if len(ctx.Events) == 0 {
		return errors.New("GenerateCVContext must have at least one event")
	}

	return nil
}

// CVProfile represents a CV profile (combination of role and audience).
type CVProfile struct {
	ID             string
	Name           string
	TargetRole     string   // principal, staff, em, senior_ic
	TargetAudience []string // hiring_manager, recruiter, peer
	Description    string
}

// GenerateCVResult is the result data returned when the intent completes.
type GenerateCVResult struct {
	// GeneratedCV is the generated CV.
	GeneratedCV *career.CVView

	// SelectedProfile is the profile that was used.
	SelectedProfile *CVProfile

	// AcceptedFields tracks which fields were accepted.
	AcceptedFields map[string]bool
}

// GenerateCVModel represents the internal state of the GenerateCV intent.
type GenerateCVModel struct {
	// context is the input context.
	context *GenerateCVContext

	// currentState is the current state of the intent.
	currentState GenerateCVState

	// selectedProfile is the currently selected profile.
	selectedProfile *CVProfile

	// selectedAudiences are the selected target audiences.
	selectedAudiences []string

	// generatedCV is the generated CV.
	generatedCV *career.CVView

	// selectedIndex is the current selection index.
	selectedIndex int

	// profileTable is the table container for displaying profiles.
	profileTable *components.TableListContainer
}

// Custom message types for GenerateCV state transitions.

// ProfileSelectedMsg indicates the user selected a profile.
type ProfileSelectedMsg struct {
	Profile *CVProfile
	Index   int
}

// AudienceSelectedMsg indicates the user selected an audience.
type AudienceSelectedMsg struct {
	Audiences []string
}

// CVGeneratedMsg indicates the CV has been generated.
type CVGeneratedMsg struct {
	CV    *career.CVView
	Error error
}
