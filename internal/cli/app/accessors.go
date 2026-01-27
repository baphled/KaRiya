package app

import "github.com/baphled/kariya/internal/cli/intents"

// SetInitialScreen sets the screen to navigate to after Init.
// Must be called before Init.
func (m *Model) SetInitialScreen(screen Screen) {
	m.initialScreen = screen
}

// SetInitialCaptureMode sets the capture mode to use when navigating directly
// to capture_event. Must be called before Init.
func (m *Model) SetInitialCaptureMode(mode string) {
	m.initialCaptureMode = mode
}

// GetState returns the current application state.
func (m *Model) GetState() AppState {
	return m.state
}

// GetActiveIntent returns the currently active intent from the router.
func (m *Model) GetActiveIntent() intents.Intent {
	return m.intentRouter.GetActiveIntent()
}

// GetIntentRouter returns the intent router (for testing purposes).
func (m *Model) GetIntentRouter() *intents.DefaultIntentRouter {
	return m.intentRouter
}

// SetStateForTesting allows tests to set the app state directly.
// This should only be used in tests to trigger edge cases.
func (m *Model) SetStateForTesting(state AppState) {
	m.state = state
}
