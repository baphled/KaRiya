package generatecv

import (
	"time"

	tea "github.com/charmbracelet/bubbletea"

	"github.com/baphled/kariya/internal/domain/career"
	cvsvc "github.com/baphled/kariya/internal/service/career/cv"
	"github.com/baphled/kariya/internal/tui/intents"
	"github.com/baphled/kariya/internal/ui/uikit/widgets"
)

// NewIntent constructs a GenerateCV intent initialised with the given
// context, selecting a default profile when one is available.
//
// Expected:
//   - ctx must pass Validate (non-nil, at least one profile and one event).
//
// Returns:
//   - A ready-to-activate intent and nil error on success.
//   - Nil intent and a validation error when context is invalid.
//
// Side effects:
//   - None.
func NewIntent(ctx *IntentValidator) (*Intent, error) {
	if err := ctx.Validate(); err != nil {
		return nil, err
	}

	selectedProfile := ctx.DefaultProfile
	if selectedProfile == nil && len(ctx.AvailableProfiles) > 0 {
		selectedProfile = ctx.AvailableProfiles[0]
	}

	baseIntent := intents.NewBaseIntent()

	return &Intent{
		BaseIntent:      baseIntent,
		context:         ctx,
		state:           StateConfiguring,
		selectedProfile: selectedProfile,
		active:          true,
	}, nil
}

// Init prepares the intent for its first render cycle by creating the
//
// Returns:
//   - A tea.Cmd value.
//
// Side effects:
//   - None.
func (i *Intent) Init() tea.Cmd {
	if i.context.DefaultProfile != nil {
		i.selectedProfile = i.context.DefaultProfile
	}

	return i.initWizardFlow()
}

// Update advances the intent state machine by processing a single Bubble Tea
// message, delegating to the wizard flow or active screen handlers.
//
// Expected:
//   - msg must be a valid tea.Msg (key press, window resize, or async result).
//
// Returns:
//   - A tea.Cmd for follow-up work (async generation, quit, etc.), or nil.
//
// Side effects:
//   - Mutates internal state, active screens, and modal visibility.
func (i *Intent) Update(msg tea.Msg) tea.Cmd {
	if !i.active {
		return nil
	}

	return i.updateWizardFlow(msg)
}

// View produces the terminal UI string for the intent's current state.
//
// Returns:
//   - A string value.
//
// Side effects:
//   - None.
func (i *Intent) View() string {
	if !i.active {
		return "GenerateCV intent is not active"
	}

	return i.wizardView()
}

// Result retrieves the outcome of the intent after it becomes inactive.
//
// Returns:
//   - A fully initialized intents.IntentResult[interface{}] ready for use.
//
// Side effects:
//   - None.
func (i *Intent) Result() *intents.IntentResult[interface{}] {
	if i.result == nil {
		return nil
	}

	return &intents.IntentResult[interface{}]{
		Status:   i.result.Status,
		Data:     i.result.Data,
		Error:    i.result.Error,
		Metadata: i.result.Metadata,
	}
}

// setCompleted marks the intent as completed with success.
func (i *Intent) setCompleted() {
	now := time.Now()
	var exportedAt *time.Time
	if i.exportedPath != "" {
		exportedAt = &now
	}
	i.result = &intents.IntentResult[*Result]{
		Status: intents.Completed,
		Data: &Result{
			GeneratedCV:     i.generatedCV,
			SelectedProfile: i.selectedProfile,
			AcceptedFields:  make(map[string]bool),
			ExportPath:      i.exportedPath,
			CVExportFormat:  string(i.selectedExportFormat),
			ExportedAt:      exportedAt,
		},
		Metadata: map[string]interface{}{
			"profile":         i.selectedProfile.ID,
			"audience":        i.selectedAudience,
			"timestamp":       now,
			"event_count":     len(i.context.Events),
			"fact_count":      len(i.context.Facts),
			"export_format":   string(i.selectedExportFormat),
			"export_location": i.exportedPath,
		},
	}
	i.active = false
}

// setCancelled marks the intent as cancelled by the user.
func (i *Intent) setCancelled() {
	i.result = &intents.IntentResult[*Result]{
		Status: intents.Cancelled,
	}
	i.active = false
}

// GetTestContext exposes the intent context for test assertions.
//
// Returns:
//   - A fully initialized IntentValidator ready for use.
//
// Side effects:
//   - None.
func (i *Intent) GetTestContext() *IntentValidator {
	return i.context
}

// GetState exposes the current state for test assertions.
//
// Returns:
//   - A State value.
//
// Side effects:
//   - None.
func (i *Intent) GetState() State {
	return i.state
}

// GetSelectedProfile exposes the selected profile for test assertions.
//
// Returns:
//   - A fully initialized CVProfile ready for use.
//
// Side effects:
//   - None.
func (i *Intent) GetSelectedProfile() *CVProfile {
	return i.selectedProfile
}

// GetSelectedAudience exposes the selected audience for test assertions.
//
// Returns:
//   - A string value.
//
// Side effects:
//   - None.
func (i *Intent) GetSelectedAudience() string {
	return i.selectedAudience
}

// GetGeneratedCV exposes the generated CV for test assertions.
//
// Returns:
//   - A fully initialized career.CVView ready for use.
//
// Side effects:
//   - None.
func (i *Intent) GetGeneratedCV() *career.CVView {
	return i.generatedCV
}

// GetReviewView exposes the review view for test assertions.
//
// Returns:
//   - A widgets.View value.
//
// Side effects:
//   - None.
func (i *Intent) GetReviewView() widgets.View {
	return i.reviewView
}

// GetPreviewView exposes the preview view for test assertions.
//
// Returns:
//   - A widgets.View value.
//
// Side effects:
//   - None.
func (i *Intent) GetPreviewView() widgets.View {
	return i.previewView
}

// GetExtractedTechnologies exposes extracted technologies for test assertions.
//
// Returns:
//   - A []*ExtractedTechnology value.
//
// Side effects:
//   - None.
func (i *Intent) GetExtractedTechnologies() []*ExtractedTechnology {
	return i.extractedTechnologies
}

// GetSelectedExportFormat exposes the export format for test assertions.
//
// Returns:
//   - A ExportFormat value.
//
// Side effects:
//   - None.
func (i *Intent) GetSelectedExportFormat() ExportFormat {
	return i.selectedExportFormat
}

// GetSelectedExportOption exposes the export option for test assertions.
//
// Returns:
//   - A ExportOption value.
//
// Side effects:
//   - None.
func (i *Intent) GetSelectedExportOption() ExportOption {
	return i.selectedExportOption
}

// GetExportedPath exposes the exported path for test assertions.
//
// Returns:
//   - A string value.
//
// Side effects:
//   - None.
func (i *Intent) GetExportedPath() string {
	return i.exportedPath
}

// GetExportError exposes the export error for test assertions.
//
// Returns:
//   - A error value.
//
// Side effects:
//   - None.
func (i *Intent) GetExportError() error {
	return i.exportError
}

// GetIsExporting exposes the exporting flag for test assertions.
//
// Returns:
//   - A bool value.
//
// Side effects:
//   - None.
func (i *Intent) GetIsExporting() bool {
	return i.isExporting
}

// GetSelectedTechnologyFocus exposes the technology focus for test assertions.
//
// Returns:
//   - A cvsvc.TechnologyFocus value.
//
// Side effects:
//   - None.
func (i *Intent) GetSelectedTechnologyFocus() cvsvc.TechnologyFocus {
	return i.selectedTechnologyFocus
}

// GetSelectedTechnologies exposes the selected technologies for test assertions.
//
// Returns:
//   - A []string value.
//
// Side effects:
//   - None.
func (i *Intent) GetSelectedTechnologies() []string {
	return i.selectedTechnologies
}

// GetSelectedFocusArea exposes the focus area for test assertions.
//
// Returns:
//   - A cvsvc.FocusArea value.
//
// Side effects:
//   - None.
func (i *Intent) GetSelectedFocusArea() cvsvc.FocusArea {
	return i.selectedFocusArea
}

// GetSelectedSkillsFormat exposes the skills format for test assertions.
//
// Returns:
//   - A string value.
//
// Side effects:
//   - None.
func (i *Intent) GetSelectedSkillsFormat() string {
	return i.selectedSkillsFormat
}

// GetSelectedSkillsLimit exposes the skills limit for test assertions.
//
// Returns:
//   - A int value.
//
// Side effects:
//   - None.
func (i *Intent) GetSelectedSkillsLimit() int {
	return i.selectedSkillsLimit
}

// GetSelectedCVLength exposes the CV length for test assertions.
//
// Returns:
//   - A string value.
//
// Side effects:
//   - None.
func (i *Intent) GetSelectedCVLength() string {
	return i.selectedCVLength
}

// SetStateForTest sets the state for test purposes.
//
// Expected:
//   - state must be valid.
//
// Side effects:
//   - None.
func (i *Intent) SetStateForTest(state State) {
	i.state = state
}

// SetSelectedExportFormatForTest sets the export format for test purposes.
//
// Expected:
//   - exportformat must be valid.
//
// Side effects:
//   - None.
func (i *Intent) SetSelectedExportFormatForTest(format ExportFormat) {
	i.selectedExportFormat = format
}

// SetSelectedExportOptionForTest sets the export option for test purposes.
//
// Expected:
//   - exportoption must be valid.
//
// Side effects:
//   - None.
func (i *Intent) SetSelectedExportOptionForTest(option ExportOption) {
	i.selectedExportOption = option
}

// SetIsExportingForTest sets the exporting flag for test purposes.
//
// Expected:
//   - bool must be valid.
//
// Side effects:
//   - None.
func (i *Intent) SetIsExportingForTest(exporting bool) {
	i.isExporting = exporting
}

// InvokeExportCVAsyncForTest exposes exportCVAsync for test purposes.
//
// Returns:
//   - A tea.Cmd value.
//
// Side effects:
//   - None.
func (i *Intent) InvokeExportCVAsyncForTest() tea.Cmd {
	return i.exportCVAsync()
}
