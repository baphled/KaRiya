package intents

import (
	"context"
	"fmt"
	"strings"

	"github.com/baphled/kariya/internal/cli/screens"
	"github.com/baphled/kariya/internal/cli/uikit/feedback"
	"github.com/baphled/kariya/internal/cli/uikit/primitives"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
)

// ConfigureSystemIntent implements the Intent interface for system configuration.
//
// Screen Orchestration:
// This intent supports gradual migration to screen-based architecture via EnableScreens().
// When screens are enabled, Update/View delegate to activeScreen.
// When screens are disabled (default), the legacy model handles Update/View.
//
// Related:
// - internal/cli/screens/configure/ (ConfigureSystem screens)
// - docs/TUI_DEVELOPER_GUIDE.md (Screen patterns)
type ConfigureSystemIntent struct {
	// Embed BaseIntent for terminal awareness, logo, and state management
	*BaseIntent

	model *ConfigureSystemModel

	// --- Screen Orchestration ---
	// activeScreen holds the current screen being displayed.
	// When non-nil and useScreens is true, Update/View delegate to this screen.
	activeScreen screens.Screen

	// useScreens controls whether to use the new screen-based architecture.
	// When true, screens handle Update/View. When false, model handles them.
	// This allows gradual migration without breaking existing functionality.
	useScreens bool

	// --- Modal Overlays ---
	// savingModal is displayed during the saving state to provide visual feedback.
	savingModal *feedback.Modal
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

	// Handle save completion/error messages to update modal
	switch msg := msg.(type) {
	case ConfigCompleteMsg:
		// Save completed successfully - show success modal briefly then clear
		c.savingModal = feedback.NewSuccessModal("Configuration saved!")
		// Let model handle the state transition
		return c.model.Update(msg)

	case ConfigErrorMsg:
		// Save failed - show error modal
		c.savingModal = feedback.NewErrorModal("Save Failed", msg.Error.Message)
		// Let model handle the state transition
		return c.model.Update(msg)

	case tea.KeyMsg:
		// Handle help modal toggle at intent level before delegating to model
		if HandleGlobalKeys(msg) == KeyHelp {
			c.ToggleHelp()
			return nil
		}

		// If success modal is showing, any key dismisses it
		if c.savingModal != nil && c.model.state == ConfigStateComplete {
			c.savingModal = nil
		}
	}

	// Delegate to active screen when using screen-based architecture
	if c.useScreens && c.activeScreen != nil {
		cmd, result := c.activeScreen.Update(msg)
		if result != nil {
			screenCmd := c.handleScreenResult(result)
			if screenCmd != nil {
				return tea.Batch(cmd, screenCmd)
			}
		}
		return cmd
	}

	// Legacy: delegate to model
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
					primitives.HelpKeyBadge("Type", "Edit", theme),
					primitives.ConfirmBadge(theme),
					primitives.CancelBadge(theme),
				),
				ThemedGlobalBadges(theme),
			)
		}
		return CombineThemedFooters(
			ThemedCustomFooter(theme,
				primitives.NavigateBadge(theme),
				primitives.SelectBadge(theme),
				primitives.HelpKeyBadge("Enter", "Edit", theme),
				primitives.SaveBadge(theme),
			),
			ThemedGlobalBadges(theme),
		)
	case ConfigStateReviewChanges:
		return CombineThemedFooters(
			ThemedCustomFooter(theme,
				primitives.ConfirmBadge(theme),
				primitives.BackBadge(theme),
			),
			ThemedGlobalBadges(theme),
		)
	case ConfigStateConfirm:
		return CombineThemedFooters(
			ThemedCustomFooter(theme,
				primitives.HelpKeyBadge("y/Enter", "Confirm", theme),
				primitives.HelpKeyBadge("n/Esc", "Cancel", theme),
			),
			ThemedGlobalBadges(theme),
		)
	case ConfigStateSaving:
		return CombineThemedFooters(
			ThemedCustomFooter(theme,
				primitives.HelpKeyBadge("...", "Please wait", theme),
			),
			ThemedGlobalBadges(theme),
		)
	case ConfigStateComplete:
		return CombineThemedFooters(
			ThemedCustomFooter(theme,
				primitives.HelpKeyBadge("Enter", "Done", theme),
			),
			ThemedGlobalBadges(theme),
		)
	case ConfigStateFailed:
		return CombineThemedFooters(
			ThemedCustomFooter(theme,
				primitives.HelpKeyBadge("r", "Retry", theme),
				primitives.CancelBadge(theme),
			),
			ThemedGlobalBadges(theme),
		)
	default:
		return ThemedGlobalBadges(theme)
	}
}

// View renders the current state using StandardView.
func (c *ConfigureSystemIntent) View() string {
	// Get terminal dimensions for modal rendering
	width, height := 120, 40 // Defaults
	if termInfo := c.GetTerminalInfo(); termInfo != nil {
		width = termInfo.Width
		height = termInfo.Height
	}

	// Delegate to active screen when using screen-based architecture
	if c.useScreens && c.activeScreen != nil {
		// Use StandardView with screen content
		view := c.CreateViewWithBreadcrumbs("Main Menu", "Configure System", c.getStateName())
		view.WithContent(c.activeScreen.View())
		view.WithHelp(c.getContextHelp()).WithFooterSeparator(true)
		baseView := view.Render()

		// Overlay saving modal if visible
		if c.savingModal != nil {
			modalContent := c.savingModal.Render(width, height)
			return c.overlayModal(baseView, modalContent, width, height)
		}

		return baseView
	}

	// Legacy: Create standard view with breadcrumbs from model
	view := c.CreateViewWithBreadcrumbs("Main Menu", "Configure System", c.getStateName())

	// Get content from model
	content := c.model.View()
	view.WithContent(content)

	// Get context-aware help
	help := c.getContextHelp()
	view.WithHelp(help).WithFooterSeparator(true)

	baseView := view.Render()

	// Overlay saving modal if visible
	if c.savingModal != nil {
		modalContent := c.savingModal.Render(width, height)
		return c.overlayModal(baseView, modalContent, width, height)
	}

	return baseView
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

// =============================================================================
// Screen Orchestration
// =============================================================================

// EnableScreens enables the screen-based architecture for this intent.
// When enabled, Update/View delegate to the active screen.
// This is opt-in to allow gradual migration.
func (c *ConfigureSystemIntent) EnableScreens() {
	c.useScreens = true
}

// IsUsingScreens returns whether the intent is using screen-based architecture.
func (c *ConfigureSystemIntent) IsUsingScreens() bool {
	return c.useScreens
}

// handleScreenResult processes a screen result and determines next action.
func (c *ConfigureSystemIntent) handleScreenResult(result screens.ScreenResult) tea.Cmd {
	if result == nil {
		return nil
	}

	switch result.Type() {
	case screens.ResultCancel:
		return c.HandleCancel(result.(*screens.CancelResult))
	case screens.ResultNavigate:
		return c.HandleNavigate(result.(*screens.NavigateResult))
	case screens.ResultSubmit:
		return c.HandleSubmit(result.(*screens.SubmitResult))
	case screens.ResultError:
		return c.HandleError(result.(*screens.ErrorResult))
	}

	return nil
}

// HandleCancel handles screen cancellation (back/escape).
// Implements ScreenResultHandler interface.
func (c *ConfigureSystemIntent) HandleCancel(result *screens.CancelResult) tea.Cmd {
	// Check if main_menu was requested
	if result.Metadata()["main_menu"] == true {
		c.setCancelled()
		return nil
	}

	// Navigate back based on current state
	switch c.model.state {
	case ConfigStateSelectDomain:
		// At root state, cancel means exit intent
		c.setCancelled()
		return nil

	case ConfigStateEditSettings:
		// Go back to domain selection
		c.model.state = ConfigStateSelectDomain
		c.activeScreen = nil // Will use model's view
		return nil

	case ConfigStateReviewChanges:
		// Go back to edit settings
		c.model.state = ConfigStateEditSettings
		c.activeScreen = nil
		return nil

	case ConfigStateConfirm:
		// Go back to review changes
		c.model.state = ConfigStateReviewChanges
		c.activeScreen = nil
		return nil

	case ConfigStateComplete, ConfigStateFailed:
		// Complete/failed - exit intent
		c.setCancelled()
		return nil

	default:
		c.setCancelled()
		return nil
	}
}

// HandleNavigate handles screen navigation results.
// Implements ScreenResultHandler interface.
func (c *ConfigureSystemIntent) HandleNavigate(result *screens.NavigateResult) tea.Cmd {
	data := result.Data()

	switch c.model.state {
	case ConfigStateReviewChanges:
		// User wants to confirm
		if data == "confirm" {
			c.model.state = ConfigStateConfirm
			c.activeScreen = nil
		}
		return nil

	case ConfigStateConfirm:
		// User confirmed or declined
		if confirmed, ok := data.(bool); ok {
			if confirmed {
				// Start saving
				c.model.state = ConfigStateSaving
				c.activeScreen = nil
				return c.saveConfiguration()
			}
			// User declined, go back to review
			c.model.state = ConfigStateReviewChanges
			c.activeScreen = nil
		}
		return nil

	case ConfigStateFailed:
		// Retry
		if data == "retry" {
			c.model.state = ConfigStateSaving
			c.activeScreen = nil
			return c.saveConfiguration()
		}
		return nil

	default:
		return nil
	}
}

// HandleSubmit handles screen submit results.
// Implements ScreenResultHandler interface.
func (c *ConfigureSystemIntent) HandleSubmit(result *screens.SubmitResult) tea.Cmd {
	// Check if main_menu was requested
	if result.Metadata()["main_menu"] == true {
		c.setCancelled()
		return nil
	}

	switch c.model.state {
	case ConfigStateEditSettings:
		// Form submitted, get changes and go to review
		if changes, ok := result.Data().(map[string]interface{}); ok {
			c.model.changes = &ConfigurationChanges{
				Domain:   c.model.domain,
				Original: make(map[string]interface{}),
				Modified: changes,
			}
		}
		c.model.state = ConfigStateReviewChanges
		c.activeScreen = nil
		return nil

	case ConfigStateComplete:
		// User dismissed success screen
		c.setCompleted(&IntentResult[interface{}]{
			Status: Completed,
			Data:   c.model.result,
		})
		return nil

	default:
		return nil
	}
}

// HandleError handles screen error results.
// Implements ScreenResultHandler interface.
func (c *ConfigureSystemIntent) HandleError(result *screens.ErrorResult) tea.Cmd {
	// Store error and transition to failed state
	c.model.error = &IntentError{
		Message: result.Message,
	}
	c.model.state = ConfigStateFailed
	c.activeScreen = nil
	return nil
}

// setCancelled marks the intent as cancelled.
func (c *ConfigureSystemIntent) setCancelled() {
	c.model.active = false
	c.model.result = &ConfigureSystemResult{
		Success: false,
		Error:   &IntentError{Message: "cancelled by user"},
	}
}

// setCompleted marks the intent as completed with the given result.
func (c *ConfigureSystemIntent) setCompleted(result *IntentResult[interface{}]) {
	c.model.active = false
}

// saveConfiguration starts the configuration save process.
func (c *ConfigureSystemIntent) saveConfiguration() tea.Cmd {
	// Show loading modal during save
	c.savingModal = feedback.NewLoadingModal("Saving configuration...", false)

	// Return a command that will simulate saving (in real implementation, this would save to disk)
	return func() tea.Msg {
		// Simulate save success
		return ConfigCompleteMsg{
			Result: &ConfigureSystemResult{
				Success: true,
				Domain:  c.model.domain,
				Changes: c.model.changes,
			},
		}
	}
}

// overlayModal overlays modal content on top of background content (centered).
// This follows the StandardView modal overlay pattern for consistent modal rendering.
func (c *ConfigureSystemIntent) overlayModal(background, modal string, width, height int) string {
	bgLines := strings.Split(background, "\n")
	modalLines := strings.Split(modal, "\n")

	// Calculate vertical position to center modal
	bgHeight := len(bgLines)
	modalHeight := len(modalLines)
	startLine := (bgHeight - modalHeight) / 2
	if startLine < 0 {
		startLine = 0
	}

	// Overlay modal lines onto background
	result := make([]string, len(bgLines))
	copy(result, bgLines)

	for i, modalLine := range modalLines {
		lineIndex := startLine + i
		if lineIndex >= 0 && lineIndex < len(result) {
			// Center modal line horizontally
			centeredModalLine := lipgloss.PlaceHorizontal(width, lipgloss.Center, modalLine)
			result[lineIndex] = centeredModalLine
		}
	}

	return strings.Join(result, "\n")
}
