package harness

import (
	"strings"

	"github.com/baphled/kariya/internal/tui/app"
)

// GetView returns the current view output.
//
// Returns:
//   - A string value.
//
// Side effects:
//   - None.
func (e *TestEnv) GetView() string {
	return e.Model.View()
}

// AssertViewContains checks that the view contains the given substring.
//
// Expected:
//   - Must be a valid string.
//
// Returns:
//   - A fully initialized TestEnv ready for use.
//
// Side effects:
//   - None.
func (e *TestEnv) AssertViewContains(substr string) *TestEnv {
	e.T.Helper()

	view := e.GetView()
	if !strings.Contains(view, substr) {
		e.T.Errorf("expected view to contain %q, but it doesn't.\nView:\n%s", substr, view)
	}

	return e
}

// AssertViewNotContains checks that the view does NOT contain the given substring.
//
// Expected:
//   - Must be a valid string.
//
// Returns:
//   - A fully initialized TestEnv ready for use.
//
// Side effects:
//   - None.
func (e *TestEnv) AssertViewNotContains(substr string) *TestEnv {
	e.T.Helper()

	view := e.GetView()
	if strings.Contains(view, substr) {
		e.T.Errorf("expected view NOT to contain %q, but it does.\nView:\n%s", substr, view)
	}

	return e
}

// AssertViewContainsAny checks that the view contains at least one of the given substrings.
//
// Expected:
//   - Must be a valid string.
//
// Returns:
//   - A fully initialized TestEnv ready for use.
//
// Side effects:
//   - None.
func (e *TestEnv) AssertViewContainsAny(substrs ...string) *TestEnv {
	e.T.Helper()

	view := e.GetView()
	for _, substr := range substrs {
		if strings.Contains(view, substr) {
			return e
		}
	}

	e.T.Errorf("expected view to contain one of %v, but none found.\nView:\n%s", substrs, view)
	return e
}

// GetMenuItems returns the list of menu items from the model.
//
// Returns:
//   - A []app.MenuItem value.
//
// Side effects:
//   - None.
func (e *TestEnv) GetMenuItems() []app.MenuItem {
	return e.Model.GetMenuItems()
}

// IsInMenuState checks if the application is currently showing the main menu.
//
// Returns:
//   - A bool value.
//
// Side effects:
//   - None.
func (e *TestEnv) IsInMenuState() bool {
	view := e.GetView()
	// Menu shows the tagline and menu items
	return strings.Contains(view, "Career Event Management System") &&
		strings.Contains(view, "Capture Event")
}

// IsInOnboardingState checks if the application is currently showing the onboarding wizard.
//
// Returns:
//   - A bool value.
//
// Side effects:
//   - None.
func (e *TestEnv) IsInOnboardingState() bool {
	// Onboarding is now a separate program that runs before the main app.
	// The main app is never in an "onboarding state".
	return false
}

// SkipOnboarding is deprecated - onboarding is now skipped by default.
//
// Returns:
//   - A fully initialized TestEnv ready for use.
//
// Side effects:
//   - None.
func (e *TestEnv) SkipOnboarding() *TestEnv {
	// No-op - onboarding is handled by bootstrap before app creation.
	return e
}
