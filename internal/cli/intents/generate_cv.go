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
	"github.com/charmbracelet/bubbles/viewport"
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
	// GenerateCVStateSelectProfile - User selects or creates a profile.
	GenerateCVStateSelectProfile GenerateCVState = "select_profile"

	// GenerateCVStateSelectAudience - User selects target audience(s).
	GenerateCVStateSelectAudience GenerateCVState = "select_audience"

	// GenerateCVStateExtractingTechnologies - Extracting technologies from user skills.
	GenerateCVStateExtractingTechnologies GenerateCVState = "extracting_technologies"

	// GenerateCVStateSelectTechnologyFocus - User selects technology focus (Language Agnostic/Generalist/Specialist).
	GenerateCVStateSelectTechnologyFocus GenerateCVState = "select_technology_focus"

	// GenerateCVStateSelectTechnologies - User selects specific technologies (for Generalist/Specialist).
	GenerateCVStateSelectTechnologies GenerateCVState = "select_technologies"

	// GenerateCVStateSelectFocusArea - User selects focus area (Backend/Frontend/Fullstack/DevOps).
	GenerateCVStateSelectFocusArea GenerateCVState = "select_focus_area"

	// GenerateCVStateSelectSkillsConfig - User configures skills section format and limit.
	GenerateCVStateSelectSkillsConfig GenerateCVState = "select_skills_config"

	// GenerateCVStateSelectLengthFormat - User selects CV length format.
	GenerateCVStateSelectLengthFormat GenerateCVState = "select_length_format"

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

	// NEW: Wizard-based workflow states (Phase 5 - Task 43).
	// These 5 states replace the 17 states above when useWizardFlow=true.

	// CVStateConfiguring - User configures CV via wizard modal
	// Replaces: SelectProfile, SelectAudience, SelectTechnologyFocus,
	//           SelectTechnologies, SelectFocusArea, SelectSkillsConfig
	CVStateConfiguring GenerateCVState = "configuring"

	// CVStateExtracting - Extracting technologies from user skills.
	// Replaces: ExtractingTechnologies.
	CVStateExtracting GenerateCVState = "extracting"

	// CVStateGenerating - CV is being generated.
	// Replaces: Generating.
	CVStateGenerating GenerateCVState = "generating"

	// CVStateReview - User reviews CV metadata and statistics via ReviewScreen.
	// Shows metadata, section summary, and bullet counts before full preview.
	CVStateReview GenerateCVState = "review"

	// CVStatePreview - User previews full CV content via CVPreviewScreen.
	// Shows scrollable full CV content with all bullets.
	CVStatePreview GenerateCVState = "preview"

	// CVStateExporting - User exports CV via export modal.
	// Replaces: ExportSelectFormat, ExportSelectLocation, Exporting, ExportComplete.
	CVStateExporting GenerateCVState = "exporting"
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

	// Technology extraction fields (NEW)
	extractedTechnologies []*ExtractedTechnology
	technologiesAvailable bool
	focusAreaSuggestion   *FocusAreaSuggestion

	// Technology Focus selection fields (NEW)
	selectedTechnologyFocus cv.TechnologyFocus
	technologyFocusIndex    int

	// Technology selection fields (NEW - for Generalist/Specialist)
	selectedTechnologies []string
	technologyCursor     int
	technologySelected   map[int]bool

	// Focus area selection fields (NEW)
	selectedFocusArea cv.FocusArea
	focusAreaCursor   int

	// Skills configuration fields (NEW - Phase 11 UI)
	selectedSkillsFormat string
	selectedSkillsLimit  int
	skillsConfigCursor   int

	// CV Length selection (Phase 8 - Task 43 wizard integration)
	selectedCVLength string

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

// Technology-related message types (NEW)

// TechnologiesExtractedMsg indicates technologies have been extracted from user skills.
type TechnologiesExtractedMsg struct {
	Technologies []*ExtractedTechnology
	Suggestion   *FocusAreaSuggestion
	Error        error
}

// TechnologyFocusSelectedMsg indicates the user selected a technology focus.
type TechnologyFocusSelectedMsg struct {
	Focus cv.TechnologyFocus
}

// TechnologiesSelectedMsg indicates the user selected specific technologies (Generalist/Specialist).
type TechnologiesSelectedMsg struct {
	Technologies []string
}

// FocusAreaSelectedMsg indicates the user selected a focus area.
type FocusAreaSelectedMsg struct {
	Area cv.FocusArea
}

// LengthFormatSelectedMsg indicates the user selected a length format.
type LengthFormatSelectedMsg struct {
	Length cv.LengthFormat
}

// Wizard-based workflow messages (Phase 5 - Task 43)

// WizardCompleteMsg indicates the configuration wizard was completed.
// Contains all configuration data from the 3-step wizard (WHO → TECH → FORMAT).
type WizardCompleteMsg struct {
	// Step 1: WHO
	ProfileID string
	Audience  string

	// Step 2: TECH
	TechFocus    string
	Technologies []string
	FocusArea    string

	// Step 3: FORMAT
	SkillsFormat string
	SkillsLimit  int
	CVLength     string

	// Pre-extracted data (optional, from async extraction)
	TechnologiesExtracted []*ExtractedTechnology
}

// Export-related types

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

// CVExportFormatSelectedMsg indicates the user selected an export format.
type CVExportFormatSelectedMsg struct {
	Format CVExportFormat
}

// CVExportOptionSelectedMsg indicates the user selected an export option.
type CVExportOptionSelectedMsg struct {
	Option CVExportOption
}

// CVExportCompleteMsg indicates export is complete.
type CVExportCompleteMsg struct {
	Path  string
	Error error
}
