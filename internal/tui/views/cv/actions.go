package cv

import "github.com/baphled/kariya/internal/ui/display"

// ActionKey is a string-typed constant identifying a CV generation action.
type ActionKey string

// Action constants for CV generation navigation.
const (
	ActionPreview ActionKey = "cv.preview"
	ActionExport  ActionKey = "cv.export"
	ActionEdit    ActionKey = "cv.edit"
	ActionConfirm ActionKey = "cv.confirm"
	ActionHelp    ActionKey = "cv.help"
)

// Nav is the typed payload for CV generation navigation results.
type Nav struct {
	Action ActionKey
	CV     display.CVView
}
