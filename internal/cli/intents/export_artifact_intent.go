package intents

import (
	"fmt"

	"github.com/baphled/kariya/internal/cli/uikit/primitives"
	tea "github.com/charmbracelet/bubbletea"
)

// ExportArtifactIntent implements the Intent interface for artifact export
type ExportArtifactIntent struct {
	// Embed BaseIntent for terminal awareness, logo, and state management
	*BaseIntent

	model *ExportArtifactModel
}

// NewExportArtifactIntent creates a new ExportArtifact intent
func NewExportArtifactIntent(context *ExportArtifactContext) (*ExportArtifactIntent, error) {
	if context == nil {
		return nil, fmt.Errorf("context is required")
	}

	model := NewExportArtifactModel(context)

	// Create BaseIntent for terminal awareness and state management
	base := NewBaseIntent()

	return &ExportArtifactIntent{
		BaseIntent: base,
		model:      model,
	}, nil
}

// Init initializes the intent
func (e *ExportArtifactIntent) Init() tea.Cmd {
	return e.model.Init()
}

// Update handles messages
func (e *ExportArtifactIntent) Update(msg tea.Msg) tea.Cmd {
	// Handle help modal toggle at intent level before delegating to model
	if keyMsg, ok := msg.(tea.KeyMsg); ok {
		if HandleGlobalKeys(keyMsg) == KeyHelp {
			e.ToggleHelp()
			return nil
		}
	}
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
	theme := e.Theme()

	switch e.model.state {
	case ExportStateSelectType:
		return CombineThemedFooters(
			ThemedNavigationFooter(theme),
			ThemedGlobalBadges(theme),
		)
	case ExportStateSelectFormat:
		return CombineThemedFooters(
			ThemedNavigationFooter(theme),
			ThemedGlobalBadges(theme),
		)
	case ExportStateSelectDest:
		return CombineThemedFooters(
			ThemedNavigationFooter(theme),
			ThemedGlobalBadges(theme),
		)
	case ExportStateConfigure:
		return CombineThemedFooters(
			ThemedFormFooter(theme),
			ThemedGlobalBadges(theme),
		)
	case ExportStatePreview:
		return CombineThemedFooters(
			ThemedDetailViewFooter(theme),
			ThemedCustomFooter(theme,
				primitives.HelpKeyBadge("Enter", "Continue", theme),
			),
			ThemedGlobalBadges(theme),
		)
	case ExportStateConfirm:
		return CombineThemedFooters(
			ThemedCustomFooter(theme,
				primitives.HelpKeyBadge("y/Enter", "Confirm", theme),
				primitives.HelpKeyBadge("n/Esc", "Cancel", theme),
			),
			ThemedGlobalBadges(theme),
		)
	case ExportStateInProgress:
		return CombineThemedFooters(
			ThemedCustomFooter(theme,
				primitives.HelpKeyBadge("...", "Please wait", theme),
			),
			ThemedGlobalBadges(theme),
		)
	case ExportStateComplete:
		return CombineThemedFooters(
			ThemedCustomFooter(theme,
				primitives.HelpKeyBadge("Enter", "Continue", theme),
			),
			ThemedGlobalBadges(theme),
		)
	case ExportStateFailed:
		return CombineThemedFooters(
			ThemedCustomFooter(theme,
				primitives.HelpKeyBadge("Enter", "Retry", theme),
				primitives.CancelBadge(theme),
			),
			ThemedGlobalBadges(theme),
		)
	default:
		return ThemedGlobalBadges(theme)
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
