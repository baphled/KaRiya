package intents

import (
	"context"
	"fmt"

	"github.com/baphled/kariya/internal/cli/components"
	tea "github.com/charmbracelet/bubbletea"
)

// ConfigureSystemIntent implements the Intent interface for system configuration
type ConfigureSystemIntent struct {
	// Embed BaseIntent for terminal awareness, logo, and state management
	*BaseIntent

	model *ConfigureSystemModel
}

// NewConfigureSystemIntent creates a new ConfigureSystem intent
func NewConfigureSystemIntent(ctx context.Context) (*ConfigureSystemIntent, error) {
	if ctx == nil {
		return nil, fmt.Errorf("context is required")
	}

	context := NewConfigureSystemContext()
	model := NewConfigureSystemModel(context)

	// Create BaseIntent for terminal awareness and state management
	base := NewBaseIntent()

	return &ConfigureSystemIntent{
		BaseIntent: base,
		model:      model,
	}, nil
}

// Init initializes the intent
func (c *ConfigureSystemIntent) Init() tea.Cmd {
	// Pass theme to model
	if theme := c.Theme(); theme != nil {
		c.model.SetTheme(theme)
	}
	return c.model.Init()
}

// Update handles messages
func (c *ConfigureSystemIntent) Update(msg tea.Msg) tea.Cmd {
	// Ensure theme stays in sync with model
	if theme := c.Theme(); theme != nil {
		c.model.SetTheme(theme)
	}
	return c.model.Update(msg)
}

// getStateName returns a human-readable name for the current state.
func (c *ConfigureSystemIntent) getStateName() string {
	switch c.model.state {
	case ConfigStateSelectDomain:
		return "Select Domain"
	case ConfigStateEditSettings:
		return "Edit Settings"
	case ConfigStateReviewChanges:
		return "Review Changes"
	case ConfigStateConfirm:
		return "Confirm"
	case ConfigStateSaving:
		return "Saving"
	case ConfigStateComplete:
		return "Complete"
	case ConfigStateFailed:
		return "Failed"
	default:
		return string(c.model.state)
	}
}

// getContextHelp returns context-aware help text for the current state.
func (c *ConfigureSystemIntent) getContextHelp() string {
	theme := c.Theme()

	switch c.model.state {
	case ConfigStateSelectDomain:
		return CombineThemedFooters(
			ThemedNavigationFooter(theme),
			ThemedGlobalBadges(theme),
		)
	case ConfigStateEditSettings:
		if c.model.editingValue {
			return CombineThemedFooters(
				ThemedCustomFooter(theme,
					components.NewKeyBadge("Type", "Edit"),
					components.ConfirmBadge(),
					components.CancelBadge(),
				),
				ThemedGlobalBadges(theme),
			)
		}
		return CombineThemedFooters(
			ThemedCustomFooter(theme,
				components.NavigateBadge(),
				components.SelectBadge(),
				components.NewKeyBadge("Enter", "Edit"),
				components.SaveBadge(),
			),
			ThemedGlobalBadges(theme),
		)
	case ConfigStateReviewChanges:
		return CombineThemedFooters(
			ThemedCustomFooter(theme,
				components.ConfirmBadge(),
				components.BackBadge(),
			),
			ThemedGlobalBadges(theme),
		)
	case ConfigStateConfirm:
		return CombineThemedFooters(
			ThemedCustomFooter(theme,
				components.NewKeyBadge("y/Enter", "Confirm"),
				components.NewKeyBadge("n/Esc", "Cancel"),
			),
			ThemedGlobalBadges(theme),
		)
	case ConfigStateSaving:
		return CombineThemedFooters(
			ThemedCustomFooter(theme,
				components.NewKeyBadge("...", "Please wait"),
			),
			ThemedGlobalBadges(theme),
		)
	case ConfigStateComplete:
		return CombineThemedFooters(
			ThemedCustomFooter(theme,
				components.NewKeyBadge("Enter", "Done"),
			),
			ThemedGlobalBadges(theme),
		)
	case ConfigStateFailed:
		return CombineThemedFooters(
			ThemedCustomFooter(theme,
				components.NewKeyBadge("r", "Retry"),
				components.CancelBadge(),
			),
			ThemedGlobalBadges(theme),
		)
	default:
		return ThemedGlobalBadges(theme)
	}
}

// View renders the current state using StandardView.
func (c *ConfigureSystemIntent) View() string {
	// Create standard view with breadcrumbs
	view := c.CreateViewWithBreadcrumbs("Main Menu", "Configure System", c.getStateName())

	// Get content from model
	content := c.model.View()
	view.WithContent(content)

	// Get context-aware help
	help := c.getContextHelp()
	view.WithHelp(help).WithFooterSeparator(true)

	return view.Render()
}

// Result returns the intent result
func (c *ConfigureSystemIntent) Result() *IntentResult[interface{}] {
	// Return nil when intent hasn't completed yet
	// Only return non-nil result when the intent has explicitly completed
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

// SetSelectedIndex sets the selected index with clamping (for testing)
func (c *ConfigureSystemIntent) SetSelectedIndex(index int) {
	c.model.SetSelectedIndex(index)
}

// GetSelectedIndex returns the selected index (for testing)
func (c *ConfigureSystemIntent) GetSelectedIndex() int {
	return c.model.GetSelectedIndex()
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

// GetDomainCount returns the number of domains (for testing)
func (c *ConfigureSystemIntent) GetDomainCount() int {
	return c.model.GetTotalItems()
}

// GetDomainPageSize returns the domain list page size (for testing)
func (c *ConfigureSystemIntent) GetDomainPageSize() int {
	return c.model.GetPageSize()
}
