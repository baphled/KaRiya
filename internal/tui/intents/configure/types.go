package configure

import (
	"github.com/baphled/kariya/internal/config"
	"github.com/baphled/kariya/internal/tui/intents"
	configviews "github.com/baphled/kariya/internal/tui/views/configure"
	"github.com/baphled/kariya/internal/ui/uikit/feedback"
)

// Intent orchestrates the configuration workflow using the activeView + ModalRegistry pattern.
type Intent struct {
	*intents.BaseIntent

	context *IntentValidator

	state  ConfigurationState
	active bool
	result *intents.IntentResult[interface{}]

	cfg      *config.Config
	settings map[ConfigurationDomain][]*ConfigurationSetting

	selectedDomain ConfigurationDomain
	pendingChanges map[string]interface{}

	settingsModal *configviews.Settings

	savingModal   *feedback.Modal
	resultModal   *feedback.Modal
	modalRegistry *intents.ModalRegistry

	configResult *SystemResult
}

// SettingsFromConfig extracts settings from config for display and editing.
//
// Expected: config must be a valid configuration object.
// Returns: A map[ConfigurationDomain][]*ConfigurationSetting value.
//
// Side effects: None.
func SettingsFromConfig(cfg *config.Config) map[ConfigurationDomain][]*ConfigurationSetting {
	return settingsFromConfig(cfg)
}

// SetContext replaces the intent's context, allowing reconfiguration.
//
// Expected:
//   - intentcontext must be valid.
//
// Side effects:
//   - None.
func (i *Intent) SetContext(ctx *IntentValidator) {
	i.context = ctx
}

// GetContext provides access to the intent's configuration and dependencies.
//
// Returns:
//   - A fully initialized IntentValidator ready for use.
//
// Side effects:
//   - None.
func (i *Intent) GetContext() *IntentValidator {
	return i.context
}

// GetState exposes the current workflow state for testing and screen orchestration.
//
// Returns:
//   - A ConfigurationState value.
//
// Side effects:
//   - None.
func (i *Intent) GetState() ConfigurationState {
	return i.state
}

// SetActive controls whether the intent processes messages and renders views.
//
// Expected:
//   - bool must be valid.
//
// Side effects:
//   - None.
func (i *Intent) SetActive(active bool) {
	i.active = active
}

// IsActive indicates whether the intent is currently processing messages and rendering.
//
// Returns:
//   - A bool value.
//
// Side effects:
//   - None.
func (i *Intent) IsActive() bool {
	return i.active
}

// GetConfig provides access to the loaded configuration for testing.
//
// Returns:
//   - A fully initialized config.Config ready for use.
//
// Side effects:
//   - None.
func (i *Intent) GetConfig() *config.Config {
	return i.cfg
}

// GetSettings provides access to the settings map for testing.
//
// Returns:
//   - A map[ConfigurationDomain][]*ConfigurationSetting value.
//
// Side effects:
//   - None.
func (i *Intent) GetSettings() map[ConfigurationDomain][]*ConfigurationSetting {
	return i.settings
}

// GetSettingsForDomain returns settings for a specific configuration domain.
//
// Expected: config must be a valid configuration object.
// Returns: A []*ConfigurationSetting value.
//
// Side effects: None.
func (i *Intent) GetSettingsForDomain(domain ConfigurationDomain) []*ConfigurationSetting {
	return i.settings[domain]
}

// SetDomain sets the currently selected configuration domain.
//
// Expected:
//   - config must be a valid configuration object.
//
// Side effects:
//   - None.
func (i *Intent) SetDomain(domain ConfigurationDomain) {
	i.selectedDomain = domain
}

// GetDomain returns the currently selected configuration domain.
//
// Returns:
//   - A ConfigurationDomain value.
//
// Side effects:
//   - None.
func (i *Intent) GetDomain() ConfigurationDomain {
	return i.selectedDomain
}

// GetResult returns the internal SystemResult for testing save outcomes.
//
// Returns:
//   - A fully initialized SystemResult ready for use.
//
// Side effects:
//   - None.
func (i *Intent) GetResult() *SystemResult {
	return i.configResult
}

// GetSavingModal returns the saving modal for testing.
//
// Returns:
//   - A fully initialized feedback.Modal ready for use.
//
// Side effects:
//   - None.
func (i *Intent) GetSavingModal() *feedback.Modal {
	return i.savingModal
}

// GetResultModal returns the result modal for testing.
//
// Returns:
//   - A fully initialized feedback.Modal ready for use.
//
// Side effects:
//   - None.
func (i *Intent) GetResultModal() *feedback.Modal {
	return i.resultModal
}
