package configure

import (
	"fmt"

	"github.com/baphled/kariya/internal/cli/behaviors"
	"github.com/baphled/kariya/internal/cli/intents"
	"github.com/baphled/kariya/internal/cli/screens"
	configscreens "github.com/baphled/kariya/internal/cli/screens/configure"
	"github.com/baphled/kariya/internal/cli/uikit/feedback"
	"github.com/baphled/kariya/internal/cli/uikit/primitives"
	"github.com/baphled/kariya/internal/config"
	tea "github.com/charmbracelet/bubbletea"
)

// getTerminalDimensions returns current terminal dimensions with fallback defaults.
func (i *Intent) getTerminalDimensions() (width, height int) {
	width, height = behaviors.DefaultModalDimensions()
	if termInfo := i.GetTerminalInfo(); termInfo != nil && termInfo.Width > 0 && termInfo.Height > 0 {
		width, height = termInfo.Width, termInfo.Height
	}
	return
}

// openEditModal opens the edit settings modal for the selected domain.
func (i *Intent) openEditModal() {
	width, height := i.getTerminalDimensions()
	settings := i.settings[i.selectedDomain]
	i.editModal = configscreens.NewEditSettingsModal(i.selectedDomain, settings, width, height)
	if theme := i.Theme(); theme != nil {
		i.editModal.SetTheme(theme)
	}
}

// openReviewModal opens the review changes modal.
func (i *Intent) openReviewModal() {
	width, height := i.getTerminalDimensions()
	i.reviewModal = configscreens.NewReviewChangesModal(i.selectedDomain, i.pendingChanges, width, height)
	if theme := i.Theme(); theme != nil {
		i.reviewModal.SetTheme(theme)
	}
}

// openConfirmModal opens the confirmation modal.
func (i *Intent) openConfirmModal() {
	width, height := i.getTerminalDimensions()
	i.confirmModal = configscreens.NewConfirmModal(
		"Confirm Changes",
		"Are you sure you want to save these changes?",
		width, height,
	)
	if theme := i.Theme(); theme != nil {
		i.confirmModal.SetTheme(theme)
	}
}

// startSaving starts the configuration save process with async command and spinner.
func (i *Intent) startSaving() tea.Cmd {
	i.savingModal = feedback.NewLoadingModal("Saving configuration...", false)

	cfg := i.cfg
	selectedDomain := i.selectedDomain
	pendingChanges := make(map[string]interface{}, len(i.pendingChanges))
	for k, v := range i.pendingChanges {
		pendingChanges[k] = v
	}

	asyncCmd := func() tea.Msg {
		for key, value := range pendingChanges {
			if err := applyConfigChange(cfg, selectedDomain, key, value); err != nil {
				return ConfigErrorMsg{
					Error: &intents.IntentError{
						Code:    "apply_failed",
						Message: fmt.Sprintf("Failed to apply change: %s", err),
					},
				}
			}
		}

		if err := config.SaveConfig(cfg); err != nil {
			return ConfigErrorMsg{
				Error: &intents.IntentError{
					Code:    "save_failed",
					Message: fmt.Sprintf("Failed to save config: %s", err),
				},
			}
		}

		return ConfigCompleteMsg{
			Result: &SystemResult{
				Success: true,
				Domain:  selectedDomain,
				Changes: &ConfigurationChanges{
					Domain:   selectedDomain,
					Modified: pendingChanges,
				},
			},
		}
	}

	return tea.Batch(asyncCmd, i.savingModal.Init())
}

// getStateName returns a human-readable name for the current workflow phase.
func (i *Intent) getStateName() string {
	if i.resultModal != nil {
		if i.configResult != nil && i.configResult.Success {
			return "Complete"
		}
		return "Failed"
	}
	if i.savingModal != nil {
		return "Saving"
	}
	if i.confirmModal != nil && i.confirmModal.IsVisible() {
		return "Confirm"
	}
	if i.reviewModal != nil && i.reviewModal.IsVisible() {
		return "Review Changes"
	}
	if i.editModal != nil && i.editModal.IsVisible() {
		return "Edit Settings"
	}
	return "Select Domain"
}

// getContextHelp returns context-aware help text for the current state.
func (i *Intent) getContextHelp() string {
	theme := i.Theme()

	if i.editModal != nil && i.editModal.IsVisible() {
		return ""
	}
	if i.reviewModal != nil && i.reviewModal.IsVisible() {
		return ""
	}
	if i.confirmModal != nil && i.confirmModal.IsVisible() {
		return ""
	}
	if i.savingModal != nil {
		return primitives.RenderHelpFooter(theme,
			primitives.HelpKeyBadge("...", "Please wait", theme),
		)
	}
	if i.resultModal != nil {
		return primitives.RenderHelpFooter(theme,
			primitives.HelpKeyBadge("Enter", "Continue", theme),
		)
	}

	return intents.CombineThemedFooters(
		primitives.RenderHelpFooter(theme,
			primitives.NavigateBadge(theme),
			primitives.SelectBadge(theme),
			primitives.BackBadge(theme),
		),
		intents.ThemedGlobalBadges(theme),
	)
}

// setCancelled marks the intent as cancelled with no result.
func (i *Intent) setCancelled() {
	i.active = false
	i.configResult = nil
}

// clearAllModals removes all active modal references.
func (i *Intent) clearAllModals() {
	i.editModal = nil
	i.reviewModal = nil
	i.confirmModal = nil
	i.savingModal = nil
	i.resultModal = nil
}

// transitionToScreen sets the active screen reference.
func (i *Intent) transitionToScreen(screen screens.Screen) {
	i.activeScreen = screen
}

// RenderDomainContent produces the raw domain selection list.
//
// Returns: A string value.
// Side effects: None.
func (i *Intent) RenderDomainContent() string {
	return i.domainScreen.RenderContent()
}

// editModalAdapter wraps EditSettingsModal to satisfy the ModalRenderer interface.
type editModalAdapter struct {
	modal         *configscreens.EditSettingsModal
	width, height int
}

// Render produces the edit settings form overlay.
//
// Expected: int must be valid.
// Returns: A string value.
//
// Side effects: None.
func (a editModalAdapter) Render(_, _ int) string {
	return a.modal.Render(a.width, a.height)
}

// reviewModalAdapter wraps ReviewChangesModal to satisfy the ModalRenderer interface.
type reviewModalAdapter struct {
	modal         *configscreens.ReviewChangesModal
	width, height int
}

// Render produces the review-changes diff overlay.
//
// Expected: int must be valid.
// Returns: A string value.
//
// Side effects: None.
func (a reviewModalAdapter) Render(_, _ int) string {
	return a.modal.Render(a.width, a.height)
}

// confirmModalAdapter wraps ConfirmModal to satisfy the ModalRenderer interface.
type confirmModalAdapter struct {
	modal         *configscreens.ConfirmModal
	width, height int
}

// Render produces the confirmation prompt overlay.
//
// Expected: int must be valid.
// Returns: A string value.
//
// Side effects: None.
func (a confirmModalAdapter) Render(_, _ int) string {
	return a.modal.Render(a.width, a.height)
}
