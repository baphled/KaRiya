package configure

import (
	"github.com/baphled/kariya/internal/cli/intents"
)

// SystemResult contains the result of configuration.
type SystemResult struct {
	Success bool
	Domain  ConfigurationDomain
	Changes *ConfigurationChanges
	Error   *intents.IntentError
}
