package intents

import (
	"context"
	"fmt"

	tea "github.com/charmbracelet/bubbletea"
)

// ExportArtifactIntent implements the Intent interface for artifact export
type ExportArtifactIntent struct {
	model *ExportArtifactModel
}

// NewExportArtifactIntent creates a new ExportArtifact intent
func NewExportArtifactIntent(ctx context.Context) (*ExportArtifactIntent, error) {
	if ctx == nil {
		return nil, fmt.Errorf("context is required")
	}

	context := NewExportArtifactContext()
	model := NewExportArtifactModel(context)

	return &ExportArtifactIntent{
		model: model,
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

// View renders the current state
func (e *ExportArtifactIntent) View() string {
	return e.model.View()
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
