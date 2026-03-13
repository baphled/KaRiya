package configure

import (
	"github.com/baphled/kariya/internal/tui/intents"
)

// SystemResult contains the result of configuration.
type SystemResult struct {
	Success bool
	Domain  ConfigurationDomain
	Changes *ConfigurationChanges
	Error   *intents.IntentError
}
