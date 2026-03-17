package configure

import (
	"github.com/baphled/kariya/internal/tui/intents"
)

// ConfigCompleteMsg represents successful configuration save.
type ConfigCompleteMsg struct {
	Result *SystemResult
}

// ConfigErrorMsg represents an error during configuration save.
type ConfigErrorMsg struct {
	Error *intents.IntentError
}
