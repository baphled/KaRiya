package intents

import (
	"context"
	"fmt"
	"strings"

	"github.com/baphled/kariya/internal/cli/configtypes"
	"github.com/baphled/kariya/internal/cli/screens"
	"github.com/baphled/kariya/internal/cli/screens/configure"
	"github.com/baphled/kariya/internal/cli/uikit/feedback"
	"github.com/baphled/kariya/internal/cli/uikit/primitives"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
)

// ConfigureSystemIntent implements the Intent interface for system configuration.
//
// Architecture: Modal Sequence Flow
// - DomainSelectScreen is the base screen (always visible)
// - Modals overlay in sequence: Edit → Review → Confirm → Saving → Result
//
// Flow:
// 1. User selects a domain from DomainSelectScreen
// 2. EditSettingsModal opens (form for editing settings)
// 3. User submits → ReviewChangesModal opens (shows diff)
// 4. User confirms → ConfirmModal opens ("Are you sure?")
// 5. User confirms → SavingModal shows (spinner)
// 6. Save completes → SuccessModal or ErrorModal shows
type ConfigureSystemIntent struct {
	*BaseIntent

	// Base screen - always visible
	domainScreen *configure.DomainSelectScreen

	// Modal overlays - only one visible at a time
	editModal    *configure.EditSettingsModal
	reviewModal  *configure.ReviewChangesModal
	confirmModal *configure.ConfirmModal
	savingModal  *feedback.Modal
	resultModal  *feedback.Modal

	// State
	selectedDomain configtypes.ConfigurationDomain
	pendingChanges map[string]interface{}
	active         bool
	result         *ConfigureSystemResult
}

// NewConfigureSystemIntent creates a new ConfigureSystem intent.
func NewConfigureSystemIntent(ctx context.Context) (*ConfigureSystemIntent, error) {
	if ctx == nil {
		return nil, fmt.Errorf("context is required")
	}

	// Available configuration domains
	domains := []configtypes.ConfigurationDomain{
		configtypes.DomainSystem,
		configtypes.DomainProfile,
		configtypes.DomainExport,
		configtypes.DomainUI,
	}

	intent := &ConfigureSystemIntent{
		BaseIntent:   NewBaseIntent(),
		domainScreen: configure.NewDomainSelectScreen(domains),
		active:       true,
	}

	return intent, nil
}

// Init initializes the intent.
func (c *ConfigureSystemIntent) Init() tea.Cmd {
	// Set theme on domain screen
	if theme := c.Theme(); theme != nil {
		c.domainScreen.SetTheme(theme)
	}

	// Set terminal info
	if termInfo := c.GetTerminalInfo(); termInfo != nil {
		c.domainScreen.SetTerminalInfo(termInfo.Width, termInfo.Height)
	}

	return c.domainScreen.Init()
}

// Update handles messages.
func (c *ConfigureSystemIntent) Update(msg tea.Msg) tea.Cmd {
	// Handle window size for all components
	if wsMsg, ok := msg.(tea.WindowSizeMsg); ok {
		c.domainScreen.SetTerminalInfo(wsMsg.Width, wsMsg.Height)
		if c.editModal != nil {
			c.editModal.Update(msg)
		}
		if c.reviewModal != nil {
			c.reviewModal.Update(msg)
		}
		if c.confirmModal != nil {
			c.confirmModal.Update(msg)
		}
	}

	// Handle key messages
	if keyMsg, ok := msg.(tea.KeyMsg); ok {
		// Global help toggle
		if HandleGlobalKeys(keyMsg) == KeyHelp {
			c.ToggleHelp()
			return nil
		}

		// Global quit
		if keyMsg.String() == "q" && c.editModal == nil && c.reviewModal == nil && c.confirmModal == nil {
			c.setCancelled()
			return nil
		}
	}

	// Route to active modal (in priority order)
	if c.resultModal != nil {
		return c.updateResultModal(msg)
	}
	if c.savingModal != nil {
		return c.updateSavingModal(msg)
	}
	if c.confirmModal != nil && c.confirmModal.IsVisible() {
		return c.updateConfirmModal(msg)
	}
	if c.reviewModal != nil && c.reviewModal.IsVisible() {
		return c.updateReviewModal(msg)
	}
	if c.editModal != nil && c.editModal.IsVisible() {
		return c.updateEditModal(msg)
	}

	// No modal - route to domain screen
	return c.updateDomainScreen(msg)
}

// updateDomainScreen handles updates to the domain selection screen.
func (c *ConfigureSystemIntent) updateDomainScreen(msg tea.Msg) tea.Cmd {
	cmd, result := c.domainScreen.Update(msg)

	if result != nil {
		switch r := result.(type) {
		case *screens.NavigateResult:
			// Domain selected - open edit modal
			if domain, ok := r.Data().(configtypes.ConfigurationDomain); ok {
				c.selectedDomain = domain
				c.openEditModal()
				return c.editModal.Init()
			}

		case *screens.CancelResult:
			// Cancel at root - exit intent
			c.setCancelled()
			return nil
		}
	}

	return cmd
}

// updateEditModal handles updates to the edit settings modal.
func (c *ConfigureSystemIntent) updateEditModal(msg tea.Msg) tea.Cmd {
	cmd := c.editModal.Update(msg)

	if !c.editModal.IsVisible() {
		if c.editModal.IsCompleted() {
			// Form submitted - get changes and open review modal
			c.pendingChanges = c.editModal.GetChanges()
			c.editModal = nil
			c.openReviewModal()
			return nil
		}
		if c.editModal.IsCancelled() {
			// Cancelled - close modal, stay on domain screen
			c.editModal = nil
			return nil
		}
	}

	return cmd
}

// updateReviewModal handles updates to the review changes modal.
func (c *ConfigureSystemIntent) updateReviewModal(msg tea.Msg) tea.Cmd {
	cmd := c.reviewModal.Update(msg)

	if !c.reviewModal.IsVisible() {
		if c.reviewModal.IsConfirmed() {
			// Confirmed - open confirm modal
			c.reviewModal = nil
			c.openConfirmModal()
			return nil
		}
		if c.reviewModal.IsCancelled() {
			// Cancelled - go back to edit
			c.reviewModal = nil
			c.openEditModal()
			return c.editModal.Init()
		}
	}

	return cmd
}

// updateConfirmModal handles updates to the confirm modal.
func (c *ConfigureSystemIntent) updateConfirmModal(msg tea.Msg) tea.Cmd {
	cmd := c.confirmModal.Update(msg)

	if !c.confirmModal.IsVisible() {
		if c.confirmModal.IsConfirmed() {
			// Confirmed - start saving
			c.confirmModal = nil
			return c.startSaving()
		}
		if c.confirmModal.IsCancelled() {
			// Cancelled - go back to review
			c.confirmModal = nil
			c.openReviewModal()
			return nil
		}
	}

	return cmd
}

// updateSavingModal handles updates during saving.
func (c *ConfigureSystemIntent) updateSavingModal(msg tea.Msg) tea.Cmd {
	switch msg := msg.(type) {
	case ConfigCompleteMsg:
		// Save completed successfully
		c.savingModal = nil
		c.result = msg.Result
		c.resultModal = feedback.NewSuccessModal("Configuration saved!")
		return nil

	case ConfigErrorMsg:
		// Save failed
		c.savingModal = nil
		c.resultModal = feedback.NewErrorModal("Save Failed", msg.Error.Message)
		return nil
	}

	return nil
}

// updateResultModal handles updates to the result modal (success/error).
func (c *ConfigureSystemIntent) updateResultModal(msg tea.Msg) tea.Cmd {
	if keyMsg, ok := msg.(tea.KeyMsg); ok {
		// Any key dismisses the result modal
		switch keyMsg.String() {
		case "enter", "esc", "q", " ":
			c.resultModal = nil
			// If success, complete the intent
			if c.result != nil && c.result.Success {
				c.active = false
			} else {
				// Error - go back to edit
				c.openEditModal()
				return c.editModal.Init()
			}
		}
	}

	return nil
}

// openEditModal opens the edit settings modal for the selected domain.
func (c *ConfigureSystemIntent) openEditModal() {
	width, height := c.getDimensions()
	settings := c.getSettingsForDomain(c.selectedDomain)
	c.editModal = configure.NewEditSettingsModal(c.selectedDomain, settings, width, height)
	if theme := c.Theme(); theme != nil {
		c.editModal.SetTheme(theme)
	}
}

// openReviewModal opens the review changes modal.
func (c *ConfigureSystemIntent) openReviewModal() {
	width, height := c.getDimensions()
	c.reviewModal = configure.NewReviewChangesModal(c.selectedDomain, c.pendingChanges, width, height)
	if theme := c.Theme(); theme != nil {
		c.reviewModal.SetTheme(theme)
	}
}

// openConfirmModal opens the confirmation modal.
func (c *ConfigureSystemIntent) openConfirmModal() {
	width, height := c.getDimensions()
	c.confirmModal = configure.NewConfirmModal(
		"Confirm Changes",
		"Are you sure you want to save these changes?",
		width, height,
	)
	if theme := c.Theme(); theme != nil {
		c.confirmModal.SetTheme(theme)
	}
}

// startSaving starts the configuration save process.
func (c *ConfigureSystemIntent) startSaving() tea.Cmd {
	c.savingModal = feedback.NewLoadingModal("Saving configuration...", false)

	return func() tea.Msg {
		// Simulate save (in real implementation, this would save to disk/database)
		return ConfigCompleteMsg{
			Result: &ConfigureSystemResult{
				Success: true,
				Domain:  c.selectedDomain,
				Changes: &ConfigurationChanges{
					Domain:   c.selectedDomain,
					Modified: c.pendingChanges,
				},
			},
		}
	}
}

// getDimensions returns the current terminal dimensions.
func (c *ConfigureSystemIntent) getDimensions() (int, int) {
	width, height := 120, 40
	if termInfo := c.GetTerminalInfo(); termInfo != nil {
		width = termInfo.Width
		height = termInfo.Height
	}
	return width, height
}

// getSettingsForDomain returns the settings for a given domain.
func (c *ConfigureSystemIntent) getSettingsForDomain(domain configtypes.ConfigurationDomain) []*configtypes.ConfigurationSetting {
	// Return sample settings based on domain
	switch domain {
	case configtypes.DomainSystem:
		return []*configtypes.ConfigurationSetting{
			{Key: "auto_save", Label: "Auto Save", Type: "bool", Value: true, Description: "Automatically save changes"},
			{Key: "backup_count", Label: "Backup Count", Type: "int", Value: 5, Description: "Number of backups to keep"},
		}
	case configtypes.DomainProfile:
		return []*configtypes.ConfigurationSetting{
			{Key: "display_name", Label: "Display Name", Type: "string", Value: "User", Description: "Your display name"},
			{Key: "email", Label: "Email", Type: "string", Value: "", Description: "Your email address"},
		}
	case configtypes.DomainExport:
		return []*configtypes.ConfigurationSetting{
			{Key: "default_format", Label: "Default Format", Type: "select", Value: "markdown", Options: []string{"markdown", "json", "yaml"}, Description: "Default export format"},
			{Key: "include_metadata", Label: "Include Metadata", Type: "bool", Value: true, Description: "Include metadata in exports"},
		}
	case configtypes.DomainUI:
		return []*configtypes.ConfigurationSetting{
			{Key: "theme", Label: "Theme", Type: "select", Value: "default", Options: []string{"default", "dark", "light"}, Description: "UI theme"},
			{Key: "show_tips", Label: "Show Tips", Type: "bool", Value: true, Description: "Show helpful tips"},
		}
	default:
		return nil
	}
}

// View renders the current state.
func (c *ConfigureSystemIntent) View() string {
	width, height := c.getDimensions()

	// Create base view with breadcrumbs
	view := c.CreateViewWithBreadcrumbs("Main Menu", "Configure System", c.getStateName())

	// Render domain screen content
	view.WithContent(c.domainScreen.View())
	view.WithHelp(c.getContextHelp()).WithFooterSeparator(true)

	baseView := view.Render()

	// Overlay modals in priority order (last one rendered on top)
	if c.editModal != nil && c.editModal.IsVisible() {
		modalContent := c.editModal.Render(width, height)
		return c.overlayModal(baseView, modalContent, width, height)
	}

	if c.reviewModal != nil && c.reviewModal.IsVisible() {
		modalContent := c.reviewModal.Render(width, height)
		return c.overlayModal(baseView, modalContent, width, height)
	}

	if c.confirmModal != nil && c.confirmModal.IsVisible() {
		modalContent := c.confirmModal.Render(width, height)
		return c.overlayModal(baseView, modalContent, width, height)
	}

	if c.savingModal != nil {
		modalContent := c.savingModal.Render(width, height)
		return c.overlayModal(baseView, modalContent, width, height)
	}

	if c.resultModal != nil {
		modalContent := c.resultModal.Render(width, height)
		return c.overlayModal(baseView, modalContent, width, height)
	}

	return baseView
}

// getStateName returns a human-readable name for the current state.
func (c *ConfigureSystemIntent) getStateName() string {
	if c.resultModal != nil {
		if c.result != nil && c.result.Success {
			return "Complete"
		}
		return "Failed"
	}
	if c.savingModal != nil {
		return "Saving"
	}
	if c.confirmModal != nil && c.confirmModal.IsVisible() {
		return "Confirm"
	}
	if c.reviewModal != nil && c.reviewModal.IsVisible() {
		return "Review Changes"
	}
	if c.editModal != nil && c.editModal.IsVisible() {
		return "Edit Settings"
	}
	return "Select Domain"
}

// getContextHelp returns context-aware help text for the current state.
func (c *ConfigureSystemIntent) getContextHelp() string {
	theme := c.Theme()

	// Modal help is rendered inside the modal itself
	if c.editModal != nil && c.editModal.IsVisible() {
		return ""
	}
	if c.reviewModal != nil && c.reviewModal.IsVisible() {
		return ""
	}
	if c.confirmModal != nil && c.confirmModal.IsVisible() {
		return ""
	}
	if c.savingModal != nil {
		return primitives.RenderHelpFooter(theme,
			primitives.HelpKeyBadge("...", "Please wait", theme),
		)
	}
	if c.resultModal != nil {
		return primitives.RenderHelpFooter(theme,
			primitives.HelpKeyBadge("Enter", "Continue", theme),
		)
	}

	// Domain selection help
	return CombineThemedFooters(
		ThemedNavigationFooter(theme),
		ThemedGlobalBadges(theme),
	)
}

// overlayModal overlays modal content on top of background content (centered).
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
			centeredModalLine := lipgloss.PlaceHorizontal(width, lipgloss.Center, modalLine)
			result[lineIndex] = centeredModalLine
		}
	}

	return strings.Join(result, "\n")
}

// Result returns the intent result.
func (c *ConfigureSystemIntent) Result() *IntentResult[interface{}] {
	if !c.active && c.result != nil {
		return &IntentResult[interface{}]{
			Status: Completed,
			Data:   c.result,
		}
	}
	if !c.active {
		return &IntentResult[interface{}]{
			Status: Cancelled,
		}
	}
	return nil
}

// setCancelled marks the intent as cancelled.
func (c *ConfigureSystemIntent) setCancelled() {
	c.active = false
	c.result = nil
}

// IsActive returns whether the intent is active.
func (c *ConfigureSystemIntent) IsActive() bool {
	return c.active
}

// GetState returns the current state (for testing).
func (c *ConfigureSystemIntent) GetState() ConfigurationState {
	if c.resultModal != nil {
		if c.result != nil && c.result.Success {
			return ConfigStateComplete
		}
		return ConfigStateFailed
	}
	if c.savingModal != nil {
		return ConfigStateSaving
	}
	if c.confirmModal != nil && c.confirmModal.IsVisible() {
		return ConfigStateConfirm
	}
	if c.reviewModal != nil && c.reviewModal.IsVisible() {
		return ConfigStateReviewChanges
	}
	if c.editModal != nil && c.editModal.IsVisible() {
		return ConfigStateEditSettings
	}
	return ConfigStateSelectDomain
}

// GetDomain returns the selected domain (for testing).
func (c *ConfigureSystemIntent) GetDomain() ConfigurationDomain {
	return c.selectedDomain
}

// GetChanges returns the pending changes (for testing).
func (c *ConfigureSystemIntent) GetChanges() *ConfigurationChanges {
	if c.pendingChanges == nil {
		return nil
	}
	return &ConfigurationChanges{
		Domain:   c.selectedDomain,
		Modified: c.pendingChanges,
	}
}

// GetResult returns the configuration result (for testing).
func (c *ConfigureSystemIntent) GetResult() *ConfigureSystemResult {
	return c.result
}

// SetState sets the state (for testing).
func (c *ConfigureSystemIntent) SetState(state ConfigurationState) {
	// Clear all modals first
	c.editModal = nil
	c.reviewModal = nil
	c.confirmModal = nil
	c.savingModal = nil
	c.resultModal = nil

	// Set appropriate modal based on state
	switch state {
	case ConfigStateEditSettings:
		if c.selectedDomain != "" {
			c.openEditModal()
		}
	case ConfigStateReviewChanges:
		if c.pendingChanges != nil {
			c.openReviewModal()
		}
	case ConfigStateConfirm:
		c.openConfirmModal()
	case ConfigStateSaving:
		c.savingModal = feedback.NewLoadingModal("Saving...", false)
	case ConfigStateComplete:
		c.result = &ConfigureSystemResult{Success: true}
		c.resultModal = feedback.NewSuccessModal("Complete")
	case ConfigStateFailed:
		c.resultModal = feedback.NewErrorModal("Failed", "Error occurred")
	}
}

// SetDomain sets the domain (for testing).
func (c *ConfigureSystemIntent) SetDomain(domain ConfigurationDomain) {
	c.selectedDomain = domain
}
