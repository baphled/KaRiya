package intents

import (
	"context"
	"errors"
	"time"

	"github.com/baphled/kariya/internal/domain/career"
	"github.com/baphled/kariya/internal/service/career/cv"
	"github.com/charmbracelet/bubbles/viewport"
)

// GenerateCVState represents the state of the GenerateCV intent.
type GenerateCVState string

const (
	// GenerateCVStateSelectProfile - User selects or creates a profile.
	GenerateCVStateSelectProfile GenerateCVState = "select_profile"

	// GenerateCVStateSelectAudience - User selects target audience(s).
	GenerateCVStateSelectAudience GenerateCVState = "select_audience"

	// GenerateCVStateGenerating - CV is being generated.
	GenerateCVStateGenerating GenerateCVState = "generating"

	// GenerateCVStatePreview - User previews the generated CV.
	GenerateCVStatePreview GenerateCVState = "preview"

	// GenerateCVStateReview - User reviews and edits the CV.
	GenerateCVStateReview GenerateCVState = "review"

	// GenerateCVStateConfirm - User confirms the CV generation.
	GenerateCVStateConfirm GenerateCVState = "confirm"

	// GenerateCVStateExportSelectFormat - User selects export format.
	GenerateCVStateExportSelectFormat GenerateCVState = "export_select_format"

	// GenerateCVStateExportSelectLocation - User selects export location.
	GenerateCVStateExportSelectLocation GenerateCVState = "export_select_location"

	// GenerateCVStateExporting - CV is being exported.
	GenerateCVStateExporting GenerateCVState = "exporting"

	// GenerateCVStateExportComplete - Export is complete.
	GenerateCVStateExportComplete GenerateCVState = "export_complete"
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

	// CVGenerationService generates CVs from configurations
	CVGenerationService cv.CVGenerationService

	// DataProcessingService processes career data for CV generation
	DataProcessingService cv.DataProcessingService

	// EnhancedBulletGenerator generates enhanced CV bullets
	EnhancedBulletGenerator cv.EnhancedBulletGenerator

	// ExportService exports CVs to various formats
	ExportService *cv.ExportService

	// AppContext is the background context for operations
	AppContext context.Context
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

	// Services are optional - only required when actually generating a CV
	// Tests may create contexts without services

	return nil
}

// CVProfile represents a CV profile (combination of role and audience).
type CVProfile struct {
	ID             string
	Name           string
	TargetRole     string // principal, staff, em, senior_ic
	TargetAudience string // hiring_manager, recruiter, peer
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

	// ExportPath is the path where the CV was exported (if exported).
	ExportPath string

	// CVExportFormat is the format the CV was exported as.
	CVExportFormat string

	// ExportedAt is when the CV was exported.
	ExportedAt *time.Time
}

// GenerateCVModel represents the internal state of the GenerateCV intent.
type GenerateCVModel struct {
	// context is the input context.
	context *GenerateCVContext

	// currentState is the current state of the intent.
	currentState GenerateCVState

	// selectedProfile is the currently selected profile.
	selectedProfile *CVProfile

	// selectedAudience is the selected target audience.
	selectedAudience string

	// audienceIndex is the index for audience selection UI
	audienceIndex int

	// generatedCV is the generated CV.
	generatedCV *career.CVView

	// selectedIndex is the current selection index.
	selectedIndex int

	// generationError tracks any errors during CV generation.
	generationError error

	// isGenerating indicates if CV generation is in progress.
	isGenerating bool

	// previewViewport is the viewport for scrolling CV preview
	previewViewport viewport.Model

	// Export-related fields
	selectedExportFormat CVExportFormat
	selectedExportOption CVExportOption
	exportedPath         string
	exportError          error
	isExporting          bool
}

// Custom message types for GenerateCV state transitions.

// ProfileSelectedMsg indicates the user selected a profile.
type ProfileSelectedMsg struct {
	Profile *CVProfile
	Index   int
}

// AudienceSelectedMsg indicates the user selected an audience.
type AudienceSelectedMsg struct {
	Audience string
}

// CVGeneratedMsg indicates the CV has been generated.
type CVGeneratedMsg struct {
	CV    *career.CVView
	Error error
}

// CVGenerationStartedMsg indicates CV generation has started.
type CVGenerationStartedMsg struct{}

// CVGenerationCompleteMsg indicates CV generation is complete.
type CVGenerationCompleteMsg struct {
	CV    *career.CVView
	Error error
}

// Export-related types

// CVExportFormat defines the export format type
type CVExportFormat string

const (
	// CVExportFormatText exports CV as plain text
	CVExportFormatText CVExportFormat = "text"
	// CVExportFormatMarkdown exports CV as markdown
	CVExportFormatMarkdown CVExportFormat = "markdown"
	// CVExportFormatYAML exports CV as YAML
	CVExportFormatYAML CVExportFormat = "yaml"
)

// CVExportOption defines where to save the CV
type CVExportOption string

const (
	// CVExportOptionSaveToFile saves CV to file
	CVExportOptionSaveToFile CVExportOption = "save_to_file"
	// CVExportOptionClipboard copies CV to clipboard
	CVExportOptionClipboard CVExportOption = "clipboard"
	// CVExportOptionCancel cancels export
	CVExportOptionCancel CVExportOption = "cancel"
)

// CVExportFormatSelectedMsg indicates the user selected an export format
type CVExportFormatSelectedMsg struct {
	Format CVExportFormat
}

// CVExportOptionSelectedMsg indicates the user selected an export option
type CVExportOptionSelectedMsg struct {
	Option CVExportOption
}

// CVExportCompleteMsg indicates export is complete
type CVExportCompleteMsg struct {
	Path  string
	Error error
}
