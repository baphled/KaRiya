package generatecv

import (
	"time"

	"github.com/baphled/kariya/internal/cli/intents"
	"github.com/baphled/kariya/internal/cli/screens"
	"github.com/baphled/kariya/internal/domain/career"
	cvsvc "github.com/baphled/kariya/internal/service/career/cv"
	tea "github.com/charmbracelet/bubbletea"
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
func NewIntent(ctx *IntentContext) (*Intent, error) {
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
// configuration wizard modal and showing it.
//
// Returns:
//   - A tea.Cmd that initialises the wizard modal.
//
// Side effects:
//   - Creates and shows the wizard modal.
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
//   - A rendered string for the terminal.
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
//   - A fully initialized IntentResult, or nil if the intent is still active.
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
	i.result = &intents.IntentResult[*Result]{
		Status: intents.Completed,
		Data: &Result{
			GeneratedCV:     i.generatedCV,
			SelectedProfile: i.selectedProfile,
			AcceptedFields:  make(map[string]bool),
		},
		Metadata: map[string]interface{}{
			"profile":     i.selectedProfile.ID,
			"audience":    i.selectedAudience,
			"timestamp":   time.Now(),
			"event_count": len(i.context.Events),
			"fact_count":  len(i.context.Facts),
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
//   - The IntentContext used to create this intent.
//
// Side effects:
//   - None.
func (i *Intent) GetTestContext() *IntentContext {
	return i.context
}

// GetState exposes the current state for test assertions.
//
// Returns:
//   - The current State value.
//
// Side effects:
//   - None.
func (i *Intent) GetState() State {
	return i.state
}

// GetSelectedProfile exposes the selected profile for test assertions.
//
// Returns:
//   - The selected CVProfile pointer.
//
// Side effects:
//   - None.
func (i *Intent) GetSelectedProfile() *CVProfile {
	return i.selectedProfile
}

// GetSelectedAudience exposes the selected audience for test assertions.
//
// Returns:
//   - The selected audience string.
//
// Side effects:
//   - None.
func (i *Intent) GetSelectedAudience() string {
	return i.selectedAudience
}

// GetGeneratedCV exposes the generated CV for test assertions.
//
// Returns:
//   - The generated CVView pointer.
//
// Side effects:
//   - None.
func (i *Intent) GetGeneratedCV() *career.CVView {
	return i.generatedCV
}

// GetReviewScreen exposes the review screen for test assertions.
//
// Returns:
//   - The ReviewScreen pointer.
//
// Side effects:
//   - None.
func (i *Intent) GetReviewScreen() screens.Screen {
	return i.reviewScreen
}

// GetPreviewScreen exposes the preview screen for test assertions.
//
// Returns:
//   - The CVPreviewScreen pointer.
//
// Side effects:
//   - None.
func (i *Intent) GetPreviewScreen() screens.Screen {
	return i.previewScreen
}

// GetExtractedTechnologies exposes extracted technologies for test assertions.
//
// Returns:
//   - The slice of ExtractedTechnology pointers.
//
// Side effects:
//   - None.
func (i *Intent) GetExtractedTechnologies() []*ExtractedTechnology {
	return i.extractedTechnologies
}

// GetSelectedExportFormat exposes the export format for test assertions.
//
// Returns:
//   - The selected ExportFormat.
//
// Side effects:
//   - None.
func (i *Intent) GetSelectedExportFormat() ExportFormat {
	return i.selectedExportFormat
}

// GetSelectedExportOption exposes the export option for test assertions.
//
// Returns:
//   - The selected ExportOption.
//
// Side effects:
//   - None.
func (i *Intent) GetSelectedExportOption() ExportOption {
	return i.selectedExportOption
}

// GetExportedPath exposes the exported path for test assertions.
//
// Returns:
//   - The exported path string.
//
// Side effects:
//   - None.
func (i *Intent) GetExportedPath() string {
	return i.exportedPath
}

// GetExportError exposes the export error for test assertions.
//
// Returns:
//   - The export error.
//
// Side effects:
//   - None.
func (i *Intent) GetExportError() error {
	return i.exportError
}

// GetIsExporting exposes the exporting flag for test assertions.
//
// Returns:
//   - The isExporting boolean.
//
// Side effects:
//   - None.
func (i *Intent) GetIsExporting() bool {
	return i.isExporting
}

// GetSelectedTechnologyFocus exposes the technology focus for test assertions.
//
// Returns:
//   - The selected TechnologyFocus.
//
// Side effects:
//   - None.
func (i *Intent) GetSelectedTechnologyFocus() cvsvc.TechnologyFocus {
	return i.selectedTechnologyFocus
}

// GetSelectedTechnologies exposes the selected technologies for test assertions.
//
// Returns:
//   - The slice of selected technology names.
//
// Side effects:
//   - None.
func (i *Intent) GetSelectedTechnologies() []string {
	return i.selectedTechnologies
}

// GetSelectedFocusArea exposes the focus area for test assertions.
//
// Returns:
//   - The selected FocusArea.
//
// Side effects:
//   - None.
func (i *Intent) GetSelectedFocusArea() cvsvc.FocusArea {
	return i.selectedFocusArea
}

// GetSelectedSkillsFormat exposes the skills format for test assertions.
//
// Returns:
//   - The selected skills format string.
//
// Side effects:
//   - None.
func (i *Intent) GetSelectedSkillsFormat() string {
	return i.selectedSkillsFormat
}

// GetSelectedSkillsLimit exposes the skills limit for test assertions.
//
// Returns:
//   - The selected skills limit.
//
// Side effects:
//   - None.
func (i *Intent) GetSelectedSkillsLimit() int {
	return i.selectedSkillsLimit
}

// GetSelectedCVLength exposes the CV length for test assertions.
//
// Returns:
//   - The selected CV length string.
//
// Side effects:
//   - None.
func (i *Intent) GetSelectedCVLength() string {
	return i.selectedCVLength
}

// SetStateForTest sets the state for test purposes.
//
// Expected:
//   - state: The state to set.
//
// Side effects:
//   - Mutates the state field.
func (i *Intent) SetStateForTest(state State) {
	i.state = state
}

// SetSelectedExportFormatForTest sets the export format for test purposes.
//
// Expected:
//   - format: The export format to set.
//
// Side effects:
//   - Mutates the selectedExportFormat field.
func (i *Intent) SetSelectedExportFormatForTest(format ExportFormat) {
	i.selectedExportFormat = format
}

// SetSelectedExportOptionForTest sets the export option for test purposes.
//
// Expected:
//   - option: The export option to set.
//
// Side effects:
//   - Mutates the selectedExportOption field.
func (i *Intent) SetSelectedExportOptionForTest(option ExportOption) {
	i.selectedExportOption = option
}

// SetIsExportingForTest sets the exporting flag for test purposes.
//
// Expected:
//   - exporting: Whether exporting is in progress.
//
// Side effects:
//   - Mutates the isExporting field.
func (i *Intent) SetIsExportingForTest(exporting bool) {
	i.isExporting = exporting
}
