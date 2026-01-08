package intents

import (
	"fmt"
	"time"

	"github.com/baphled/kariya/internal/cli/components"
	tea "github.com/charmbracelet/bubbletea"
)

// ExportArtifactIntent implements the Intent interface for artifact export
type ExportArtifactIntent struct {
	// Embed BaseIntent for terminal awareness, logo, and state management
	*BaseIntent

	model *ExportArtifactModel

	// loadingRotator rotates through export-specific loading messages
	loadingRotator *components.LoadingMessageRotator
}

// NewExportArtifactIntent creates a new ExportArtifact intent
func NewExportArtifactIntent(context *ExportArtifactContext) (*ExportArtifactIntent, error) {
	if context == nil {
		return nil, fmt.Errorf("context is required")
	}

	model := NewExportArtifactModel(context)

	// Create BaseIntent for terminal awareness and state management
	base := NewBaseIntent()

	// Create loading message rotator with export-specific messages
	loadingRotator := components.NewLoadingMessageRotator([]string{
		"📦 Preparing export...",
		"🔍 Gathering data...",
		"✨ Formatting output...",
		"💾 Writing file...",
		"✅ Export complete!",
	}, 2*time.Second)

	return &ExportArtifactIntent{
		BaseIntent:     base,
		model:          model,
		loadingRotator: loadingRotator,
	}, nil
}

// Init initializes the intent
func (e *ExportArtifactIntent) Init() tea.Cmd {
	return e.model.Init()
}

// Update handles messages
func (e *ExportArtifactIntent) Update(msg tea.Msg) tea.Cmd {
	return e.model.Update(msg)
}

// getStateName returns a human-readable name for the current state.
func (e *ExportArtifactIntent) getStateName() string {
	switch e.model.state {
	case ExportStateSelectType:
		return "Select Type"
	case ExportStateSelectFormat:
		return "Select Format"
	case ExportStateSelectDest:
		return "Select Destination"
	case ExportStateConfigure:
		return "Configure"
	case ExportStatePreview:
		return "Preview"
	case ExportStateConfirm:
		return "Confirm"
	case ExportStateInProgress:
		return "Exporting"
	case ExportStateComplete:
		return "Complete"
	case ExportStateFailed:
		return "Failed"
	default:
		return string(e.model.state)
	}
}

// getContextHelp returns context-aware help text for the current state.
func (e *ExportArtifactIntent) getContextHelp() string {
	base := "q Quit  m Main Menu"

	switch e.model.state {
	case ExportStateSelectType:
		return CombineFooters(NavigationFooter(), base)
	case ExportStateSelectFormat:
		return CombineFooters(NavigationFooter(), base)
	case ExportStateSelectDest:
		return CombineFooters(NavigationFooter(), base)
	case ExportStateConfigure:
		return CombineFooters(FormFooter(), base)
	case ExportStatePreview:
		return CombineFooters(DetailViewFooter(), "Enter Continue", base)
	case ExportStateConfirm:
		return CombineFooters("y/Enter Confirm  n/Esc Cancel", base)
	case ExportStateInProgress:
		return CombineFooters("Please wait...", base)
	case ExportStateComplete:
		return CombineFooters("Enter Continue", base)
	case ExportStateFailed:
		return CombineFooters("Enter Retry  Esc Cancel", base)
	default:
		return base
	}
}

// View renders the current state using StandardView.
func (e *ExportArtifactIntent) View() string {
	// Create standard view with breadcrumbs
	view := e.CreateViewWithBreadcrumbs("Main Menu", "Export Artifact", e.getStateName())

	// Get content from model
	content := e.model.View()
	view.WithContent(content)

	// Get context-aware help
	help := e.getContextHelp()
	view.WithHelp(help).WithFooterSeparator(true)

	return view.Render()
}

// Result returns the intent result
func (e *ExportArtifactIntent) Result() *IntentResult[interface{}] {
	return e.model.Result()
}

// GetState returns the current state (for testing)
func (e *ExportArtifactIntent) GetState() ExportState {
	return e.model.state
}

// GetConfig returns the current export configuration (for testing)
func (e *ExportArtifactIntent) GetConfig() *ExportConfiguration {
	return e.model.config
}

// GetResult returns the export result (for testing)
func (e *ExportArtifactIntent) GetResult() *ExportArtifactResult {
	return e.model.result
}

// SetSelectedIndex sets the selected index (for testing)
func (e *ExportArtifactIntent) SetSelectedIndex(index int) {
	e.model.selectedIndex = index
}

// GetSelectedIndex returns the selected index (for testing)
func (e *ExportArtifactIntent) GetSelectedIndex() int {
	return e.model.selectedIndex
}

// SetState sets the state (for testing)
func (e *ExportArtifactIntent) SetState(state ExportState) {
	e.model.state = state
}

// SetConfig sets the configuration (for testing)
func (e *ExportArtifactIntent) SetConfig(config *ExportConfiguration) {
	e.model.config = config
}

// IsActive returns whether the intent is active
func (e *ExportArtifactIntent) IsActive() bool {
	return e.model.active
}
