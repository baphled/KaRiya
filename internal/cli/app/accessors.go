package app

import (
	"github.com/baphled/kariya/internal/cli/intents"
)

// SetInitialScreen sets the initial screen to display.
// This must be called before Init() to take effect.
// When set to ListScreen, the app will start in browse_timeline intent.
func (m *Model) SetInitialScreen(screen Screen) {
	m.initialScreen = screen
}

// SetInitialCaptureMode sets the initial capture mode for CaptureEvent intent.
// This must be called before Init() to take effect.
// When set, the app will start in capture_event intent with the specified mode.
// Valid modes: "manual", "burst", "csv".
func (m *Model) SetInitialCaptureMode(mode string) {
	m.initialCaptureMode = mode
}

// GetState returns the current application state.
func (m *Model) GetState() AppState {
	return m.state
}

// GetActiveIntent returns the currently active intent.
func (m *Model) GetActiveIntent() intents.Intent {
	return m.intentRouter.GetActiveIntent()
}
