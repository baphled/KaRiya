package configure

import (
	"fmt"

	"github.com/baphled/kariya/internal/config"
	"github.com/baphled/kariya/internal/tui/intents"
	configviews "github.com/baphled/kariya/internal/tui/views/configure"
	"github.com/baphled/kariya/internal/ui/behaviors"
	"github.com/baphled/kariya/internal/ui/uikit/feedback"
	"github.com/baphled/kariya/internal/ui/uikit/primitives"
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

// openSettingsModal opens the unified settings modal.
func (i *Intent) openSettingsModal() {
	width, height := i.getTerminalDimensions()
	i.settingsModal = configviews.NewSettings(i.settings, width, height)
	if theme := i.Theme(); theme != nil {
		i.settingsModal.SetTheme(theme)
	}
	i.rebuildModalRegistry()
}

// startSaving starts the configuration save process with async command and spinner.
func (i *Intent) startSaving() tea.Cmd {
	i.savingModal = feedback.NewLoadingModal("Saving configuration...", false)
	i.rebuildModalRegistry()

	cfg := i.cfg
	settings := i.settings
	pendingChanges := make(map[string]interface{}, len(i.pendingChanges))
	for k, v := range i.pendingChanges {
		pendingChanges[k] = v
	}

	asyncCmd := func() tea.Msg {
		if err := ApplyChanges(cfg, settings, pendingChanges); err != nil {
			return ConfigErrorMsg{
				Error: &intents.IntentError{
					Code:    "apply_failed",
					Message: fmt.Sprintf("Failed to apply change: %s", err),
				},
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
				Domain:  i.selectedDomain,
				Changes: &ConfigurationChanges{
					Domain:   i.selectedDomain,
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
	if i.settingsModal != nil {
		return "Configure"
	}
	return "Select Domain"
}

// getContextHelp returns context-aware help text for the current state.
func (i *Intent) getContextHelp() string {
	theme := i.Theme()

	if i.settingsModal != nil {
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
	i.settingsModal = nil
	i.savingModal = nil
	i.resultModal = nil
}

// settingsModalAdapter wraps SettingsModal to satisfy the ModalRenderer interface.
type settingsModalAdapter struct {
	modal         *configviews.Settings
	width, height int
}

// Render produces the settings modal overlay.
//
// Expected: int must be valid.
// Returns: A string value.
//
// Side effects: None.
func (a settingsModalAdapter) Render(_, _ int) string {
	return a.modal.Render(a.width, a.height)
}

// rebuildModalRegistry registers feedback modals (savingModal, resultModal) with the registry.
// SettingsModal uses custom interface so it's handled separately via individual handlers.
//
// Priority order (highest to lowest):
// 1. savingModal (loading spinner) - handled via registry
// 2. resultModal (success/error modal) - handled via registry
//
// Side effects:
//   - Clears and rebuilds the modalRegistry
func (i *Intent) rebuildModalRegistry() {
	i.modalRegistry.Clear()

	width, height := i.getTerminalDimensions()

	// Register savingModal as an error modal
	if i.savingModal != nil {
		i.modalRegistry.Register(intents.NewErrorModalAdapter(
			i.savingModal,
			width,
			height,
			i.Theme(),
		))
	}

	// Register resultModal as an error modal
	if i.resultModal != nil {
		i.modalRegistry.Register(intents.NewErrorModalAdapter(
			i.resultModal,
			width,
			height,
			i.Theme(),
		))
	}
}
