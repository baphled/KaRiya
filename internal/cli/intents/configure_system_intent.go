package intents

import (
	"context"
	"fmt"

	tea "github.com/charmbracelet/bubbletea"
)

// ConfigureSystemIntent implements the Intent interface for system configuration
type ConfigureSystemIntent struct {
	model *ConfigureSystemModel
}

// NewConfigureSystemIntent creates a new ConfigureSystem intent
func NewConfigureSystemIntent(ctx context.Context) (*ConfigureSystemIntent, error) {
	if ctx == nil {
		return nil, fmt.Errorf("context is required")
	}

	context := NewConfigureSystemContext()
	model := NewConfigureSystemModel(context)

	return &ConfigureSystemIntent{
		model: model,
	}, nil
}

// Init initializes the intent
func (c *ConfigureSystemIntent) Init() tea.Cmd {
	return c.model.Init()
}

// Update handles messages
func (c *ConfigureSystemIntent) Update(msg tea.Msg) tea.Cmd {
	return c.model.Update(msg)
}

// View renders the current state
func (c *ConfigureSystemIntent) View() string {
	return c.model.View()
}

// Result returns the intent result
func (c *ConfigureSystemIntent) Result() *IntentResult[interface{}] {
	return c.model.Result()
}

// GetState returns the current state (for testing)
func (c *ConfigureSystemIntent) GetState() ConfigurationState {
	return c.model.state
}

// GetDomain returns the current domain (for testing)
func (c *ConfigureSystemIntent) GetDomain() ConfigurationDomain {
	return c.model.domain
}

// GetChanges returns the current changes (for testing)
func (c *ConfigureSystemIntent) GetChanges() *ConfigurationChanges {
	return c.model.changes
}

// GetResult returns the configuration result (for testing)
func (c *ConfigureSystemIntent) GetResult() *ConfigureSystemResult {
	return c.model.result
}

// SetSelectedIndex sets the selected index (for testing)
func (c *ConfigureSystemIntent) SetSelectedIndex(index int) {
	c.model.selectedIndex = index
}

// GetSelectedIndex returns the selected index (for testing)
func (c *ConfigureSystemIntent) GetSelectedIndex() int {
	return c.model.selectedIndex
}

// SetState sets the state (for testing)
func (c *ConfigureSystemIntent) SetState(state ConfigurationState) {
	c.model.state = state
}

// SetDomain sets the domain (for testing)
func (c *ConfigureSystemIntent) SetDomain(domain ConfigurationDomain) {
	c.model.domain = domain
}

// IsActive returns whether the intent is active
func (c *ConfigureSystemIntent) IsActive() bool {
	return c.model.active
}
