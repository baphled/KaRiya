package intents

import (
	"context"
	"errors"
	"time"

	"github.com/baphled/kariya/internal/cli/screens"
	"github.com/baphled/kariya/internal/cli/types"
	"github.com/baphled/kariya/internal/config"
	"github.com/baphled/kariya/internal/domain/career"
	careerRepo "github.com/baphled/kariya/internal/repository/career"
	"github.com/baphled/kariya/internal/service/career/cv"
	"github.com/baphled/kariya/internal/service/career/technology"
)

// ExtractedTechnology re-exports technology.ExtractedTechnology so that the intents
// package can reference extracted skill data without importing the technology service
// package directly. Each ExtractedTechnology holds a skill's identity, category, and
// the set of career events that reference it, enabling the CV generation wizard to
// offer technology-focus choices based on the user's actual skill distribution.
type ExtractedTechnology = technology.ExtractedTechnology

// FocusAreaSuggestion re-exports technology.FocusAreaSuggestion so that intents can
// consume AI-generated focus area recommendations without a direct dependency on the
// technology service package. A FocusAreaSuggestion pairs a recommended focus area
// (e.g., Backend, Frontend) with a confidence score and the category-count evidence
// that supports the recommendation.
type FocusAreaSuggestion = technology.FocusAreaSuggestion

// GenerateCVState represents the state of the GenerateCV intent.
type GenerateCVState string

const (
	// CVStateConfiguring - User configures CV via wizard modal.
	CVStateConfiguring GenerateCVState = "configuring"

	// CVStateExtracting - Extracting technologies from user skills.
	CVStateExtracting GenerateCVState = "extracting"

	// CVStateGenerating - CV is being generated.
	CVStateGenerating GenerateCVState = "generating"

	// CVStateReview - User reviews CV metadata and statistics via ReviewScreen.
	CVStateReview GenerateCVState = "review"

	// CVStatePreview - User previews full CV content via CVPreviewScreen.
	CVStatePreview GenerateCVState = "preview"

	// CVStateExporting - User selects export format and location via export modal.
	CVStateExporting GenerateCVState = "exporting"

	// CVStateExportSelectLocation - User selects export location.
	CVStateExportSelectLocation GenerateCVState = "export_select_location"

	// CVStateExportComplete - Export is complete.
	CVStateExportComplete GenerateCVState = "export_complete"
)

// GenerateCVContext is the input context passed to the GenerateCV intent.
type GenerateCVContext struct {
	// AvailableProfiles are the profiles the user can choose from.
	AvailableProfiles []*CVProfile

	// Events are the career events to use for CV generation.
	Events []*career.Event

	// Facts are the career facts to use for CV generation.
	Facts []*career.Fact

	// DefaultProfile is the profile to select by default.
	DefaultProfile *CVProfile

	// CVGenerationService generates CVs from configurations
	CVGenerationService cv.CVGenerationService

	// DataProcessingService processes career data for CV generation
	DataProcessingService cv.DataProcessingService

	// BulletGenerator generates CV bullets
	BulletGenerator cv.BulletGenerator

	// ExportService exports CVs to various formats
	ExportService *cv.ExportService

	// SkillRepository provides access to user skills (for technology extraction)
	SkillRepository careerRepo.SkillRepository

	// EventRepository provides access to career events (for technology extraction)
	EventRepository careerRepo.EventRepository

	// ProfileConfig is the user's profile configuration for narrative CVs
	ProfileConfig *config.ProfileConfig

	// AppContext is the background context for operations
	AppContext context.Context

	// ReviewScreenFactory creates a ReviewScreen (avoids import cycle)
	// Signature: func(cv *career.CVView) screens.Screen
	ReviewScreenFactory func(cv *career.CVView) screens.Screen

	// PreviewScreenFactory creates a CVPreviewScreen (avoids import cycle)
	// Signature: func(cv *career.CVView) screens.Screen
	PreviewScreenFactory func(cv *career.CVView) screens.Screen
}

// Validate ensures the context carries the minimum data needed to start the
// CV generation workflow.
//
// Expected:
//   - ctx must be non-nil.
//   - ctx.AvailableProfiles must contain at least one profile.
//   - ctx.Events must contain at least one event.
//
// Returns:
//   - Nil when all preconditions are met.
//   - A descriptive error when any required field is missing or empty.
//
// Side effects:
//   - None.
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

// CVProfile is a type alias for the shared types.CVProfile.
// This allows intents to use CVProfile directly while the actual type
// is defined in types package for sharing with screens.
type CVProfile = types.CVProfile

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
// It tracks the selected profile, audience, technology/focus preferences,
// skills configuration, and export progress throughout the wizard workflow.
type GenerateCVModel struct {
	context      *GenerateCVContext
	currentState GenerateCVState

	selectedProfile  *CVProfile
	selectedAudience string
	generatedCV      *career.CVView

	extractedTechnologies []*ExtractedTechnology
	focusAreaSuggestion   *FocusAreaSuggestion

	selectedTechnologyFocus cv.TechnologyFocus
	selectedTechnologies    []string
	selectedFocusArea       cv.FocusArea

	selectedSkillsFormat string
	selectedSkillsLimit  int
	selectedCVLength     string

	selectedExportFormat CVExportFormat
	selectedExportOption CVExportOption
	exportedPath         string
	exportError          error
	isExporting          bool
}

// CVGenerationCompleteMsg indicates CV generation is complete.
type CVGenerationCompleteMsg struct {
	CV    *career.CVView
	Error error
}

// TechnologiesExtractedMsg indicates technologies have been extracted from user skills.
type TechnologiesExtractedMsg struct {
	Technologies []*ExtractedTechnology
	Suggestion   *FocusAreaSuggestion
	Error        error
}

// WizardCompleteMsg indicates the configuration wizard was completed.
// Contains all configuration data from the 3-step wizard: WHO (ProfileID,
// Audience), TECH (TechFocus, Technologies, FocusArea), and FORMAT
// (SkillsFormat, SkillsLimit, CVLength). TechnologiesExtracted carries
// pre-extracted skill data when available from async extraction.
type WizardCompleteMsg struct {
	ProfileID string
	Audience  string

	TechFocus    string
	Technologies []string
	FocusArea    string

	SkillsFormat          string
	SkillsLimit           int
	CVLength              string
	TechnologiesExtracted []*ExtractedTechnology
}

// CVExportFormat defines the export format type.
type CVExportFormat string

const (
	// CVExportFormatText exports CV as plain text.
	CVExportFormatText CVExportFormat = "text"
	// CVExportFormatMarkdown exports CV as markdown.
	CVExportFormatMarkdown CVExportFormat = "markdown"
	// CVExportFormatYAML exports CV as YAML.
	CVExportFormatYAML CVExportFormat = "yaml"
)

// CVExportOption defines where to save the CV.
type CVExportOption string

const (
	// CVExportOptionSaveToFile saves CV to file.
	CVExportOptionSaveToFile CVExportOption = "save_to_file"
	// CVExportOptionClipboard copies CV to clipboard.
	CVExportOptionClipboard CVExportOption = "clipboard"
	// CVExportOptionCancel cancels export.
	CVExportOptionCancel CVExportOption = "cancel"
)

// CVExportCompleteMsg indicates export is complete.
type CVExportCompleteMsg struct {
	Path  string
	Error error
}
