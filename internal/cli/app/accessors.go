package app

import "github.com/baphled/kariya/internal/cli/intents"

// SetInitialScreen sets the screen to navigate to after Init.
//
// Expected:
//   - screen must be valid.
//
// Side effects:
//   - None.
func (m *Model) SetInitialScreen(screen Screen) {
	m.initialScreen = screen
}

// SetInitialCaptureMode sets the capture mode to use when navigating directly
//
// Expected:
//   - Must be a valid string.
//
// Side effects:
//   - None.
func (m *Model) SetInitialCaptureMode(mode string) {
	m.initialCaptureMode = mode
}

// GetState returns the current application state.
//
// Returns:
//   - A State value.
//
// Side effects:
//   - None.
func (m *Model) GetState() State {
	return m.state
}

// GetActiveIntent returns the currently active intent from the router.
//
// Returns:
//   - A intents.Intent value.
//
// Side effects:
//   - None.
func (m *Model) GetActiveIntent() intents.Intent {
	return m.intentRouter.GetActiveIntent()
}

// GetIntentRouter returns the intent router (for testing purposes).
//
// Returns:
//   - A fully initialized intents.DefaultIntentRouter ready for use.
//
// Side effects:
//   - None.
func (m *Model) GetIntentRouter() *intents.DefaultIntentRouter {
	return m.intentRouter
}

// SetStateForTesting allows tests to set the app state directly.
//
// Expected:
//   - state must be valid.
//
// Side effects:
//   - None.
func (m *Model) SetStateForTesting(state State) {
	m.state = state
}
