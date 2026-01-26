package intents

import (
	"context"
	"fmt"

	"github.com/baphled/kariya/internal/cli/configtypes"
	"github.com/baphled/kariya/internal/cli/screens"
	"github.com/baphled/kariya/internal/cli/screens/configure"
	"github.com/baphled/kariya/internal/cli/uikit/feedback"
	"github.com/baphled/kariya/internal/cli/uikit/primitives"
	"github.com/baphled/kariya/internal/config"
	tea "github.com/charmbracelet/bubbletea"
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

	// Configuration data loaded from file
	cfg      *config.Config
	settings map[configtypes.ConfigurationDomain][]*configtypes.ConfigurationSetting

	// State
	selectedDomain configtypes.ConfigurationDomain
	pendingChanges map[string]interface{}
	active         bool
	saving         bool // Prevents concurrent save operations
	result         *ConfigureSystemResult
}

// NewConfigureSystemIntent creates a new ConfigureSystem intent.
func NewConfigureSystemIntent(ctx context.Context) (*ConfigureSystemIntent, error) {
	if ctx == nil {
		return nil, fmt.Errorf("context is required")
	}

	// Load configuration from file (or use defaults)
	cfg, err := config.LoadConfig()
	if err != nil {
		// Fall back to defaults if config file doesn't exist
		cfg = config.DefaultConfig()
	}

	// Convert config to settings
	settings := settingsFromConfig(cfg)

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
		cfg:          cfg,
		settings:     settings,
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

	// Handle async save completion messages globally (regardless of current modal state)
	// These messages should clear all other modals and show the result
	switch msg := msg.(type) {
	case ConfigCompleteMsg:
		// Clear all modals and reset saving state
		c.editModal = nil
		c.reviewModal = nil
		c.confirmModal = nil
		c.savingModal = nil
		c.saving = false
		// Update settings from the saved config
		c.settings = settingsFromConfig(c.cfg)
		// Set result
		c.result = msg.Result
		c.resultModal = feedback.NewSuccessModal("Configuration saved!")
		return nil
	case ConfigErrorMsg:
		// Clear all modals and reset saving state
		c.editModal = nil
		c.reviewModal = nil
		c.confirmModal = nil
		c.savingModal = nil
		c.saving = false
		// Set error result
		c.resultModal = feedback.NewErrorModal("Save Failed", msg.Error.Message)
		return nil
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

	case feedback.ModalSpinnerTickMsg:
		// Forward tick to loading modal to advance spinner.
		return c.savingModal.Update(msg)

	default:
		// For any other message, ensure spinner tick is running.
		return c.savingModal.Init()
	}
}

// updateResultModal handles updates to the result modal (success/error).
func (c *ConfigureSystemIntent) updateResultModal(msg tea.Msg) tea.Cmd {
	if keyMsg, ok := msg.(tea.KeyMsg); ok {
		// Any key dismisses the result modal
		switch keyMsg.String() {
		case "enter", "esc", "q", " ":
			c.resultModal = nil
			// If error, go back to edit
			if c.result == nil || !c.result.Success {
				c.openEditModal()
				return c.editModal.Init()
			}
			// Success - complete the intent
			c.active = false
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
// Uses a saving flag to prevent concurrent save operations.
func (c *ConfigureSystemIntent) startSaving() tea.Cmd {
	// Prevent concurrent save operations
	if c.saving {
		return nil
	}
	c.saving = true
	c.savingModal = feedback.NewLoadingModal("Saving configuration...", false)

	// Capture values to avoid race conditions with the goroutine
	cfg := c.cfg
	selectedDomain := c.selectedDomain
	pendingChanges := make(map[string]interface{}, len(c.pendingChanges))
	for k, v := range c.pendingChanges {
		pendingChanges[k] = v
	}

	asyncCmd := func() tea.Msg {
		// Apply changes to config
		for key, value := range pendingChanges {
			if err := applyConfigChange(cfg, selectedDomain, key, value); err != nil {
				return ConfigErrorMsg{
					Error: &IntentError{
						Code:    "apply_failed",
						Message: fmt.Sprintf("Failed to apply change: %s", err),
					},
				}
			}
		}

		// Save config to file
		if err := config.SaveConfig(cfg); err != nil {
			return ConfigErrorMsg{
				Error: &IntentError{
					Code:    "save_failed",
					Message: fmt.Sprintf("Failed to save config: %s", err),
				},
			}
		}

		return ConfigCompleteMsg{
			Result: &ConfigureSystemResult{
				Success: true,
				Domain:  selectedDomain,
				Changes: &ConfigurationChanges{
					Domain:   selectedDomain,
					Modified: pendingChanges,
				},
			},
		}
	}

	// Batch async save with spinner init to start animation immediately.
	return tea.Batch(asyncCmd, c.savingModal.Init())
}

// getDimensions returns the current terminal dimensions.
func (c *ConfigureSystemIntent) getDimensions() (int, int) {
	width, height := 120, 40
	if termInfo := c.GetTerminalInfo(); termInfo != nil && termInfo.Width > 0 && termInfo.Height > 0 {
		width = termInfo.Width
		height = termInfo.Height
	}
	return width, height
}

// getSettingsForDomain returns the settings for a given domain.
func (c *ConfigureSystemIntent) getSettingsForDomain(domain configtypes.ConfigurationDomain) []*configtypes.ConfigurationSetting {
	// Return settings loaded from config file
	if c.settings != nil {
		if domainSettings, ok := c.settings[domain]; ok {
			return domainSettings
		}
	}
	return nil
}

// View renders the current state.
func (c *ConfigureSystemIntent) View() string {
	width, height := c.getDimensions()

	// Create base view with breadcrumbs
	view := c.CreateViewWithBreadcrumbs("Main Menu", "Configure System", c.getStateName())

	// Render domain screen RAW content (not the full View which has its own layout)
	// This avoids double-wrapping in StandardView
	view.WithContent(c.domainScreen.RenderContent())
	view.WithHelp(c.getContextHelp()).WithFooterSeparator(true)

	// Use ScreenLayout's ShowModalOverlay for proper modal handling
	// The ScreenLayout handles dimensions and centering automatically
	if c.editModal != nil && c.editModal.IsVisible() {
		view.ShowModalOverlay(editModalAdapter{c.editModal, width, height})
	} else if c.reviewModal != nil && c.reviewModal.IsVisible() {
		view.ShowModalOverlay(reviewModalAdapter{c.reviewModal, width, height})
	} else if c.confirmModal != nil && c.confirmModal.IsVisible() {
		view.ShowModalOverlay(confirmModalAdapter{c.confirmModal, width, height})
	} else if c.savingModal != nil {
		view.ShowModalOverlay(c.savingModal)
	} else if c.resultModal != nil {
		view.ShowModalOverlay(c.resultModal)
	}

	return view.Render()
}

// Modal adapters to satisfy ModalRenderer interface
// These wrap the configure modals to provide the Render(width, height) signature

type editModalAdapter struct {
	modal         *configure.EditSettingsModal
	width, height int
}

func (a editModalAdapter) Render(_, _ int) string {
	return a.modal.Render(a.width, a.height)
}

type reviewModalAdapter struct {
	modal         *configure.ReviewChangesModal
	width, height int
}

func (a reviewModalAdapter) Render(_, _ int) string {
	return a.modal.Render(a.width, a.height)
}

type confirmModalAdapter struct {
	modal         *configure.ConfirmModal
	width, height int
}

func (a confirmModalAdapter) Render(_, _ int) string {
	return a.modal.Render(a.width, a.height)
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

// Testing helpers

// GetSavingModal returns the saving modal (for testing).
func (c *ConfigureSystemIntent) GetSavingModal() *feedback.Modal {
	return c.savingModal
}

// GetResultModal returns the result modal (for testing).
func (c *ConfigureSystemIntent) GetResultModal() *feedback.Modal {
	return c.resultModal
}

// RenderDomainContent returns the domain screen content (for testing).
func (c *ConfigureSystemIntent) RenderDomainContent() string {
	return c.domainScreen.RenderContent()
}
