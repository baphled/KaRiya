package intents

import (
	"errors"
	"fmt"
	"strings"
	"time"

	"github.com/baphled/kariya/internal/cli/components"
)

// CreateStandardView creates a standardized view with logo and automatic state modals.
// The view is configured with:
// - Terminal info from BaseIntent
// - Logo (if available) with configured spacing
// - Automatic modal display based on intent state (error, loading, progress, success)
// - Full width content rendering
//
// Example usage:
//
//	func (i *MyIntent) View() string {
//	    view := CreateStandardView(i.BaseIntent)
//	    view.WithContent(i.renderContent())
//	    view.WithHelp("↑/k Up  ↓/j Down  Enter Select  Esc Back")
//	    return view.Render()
//	}
func CreateStandardView(b *BaseIntent) *components.StandardView {
	view := components.NewStandardView(b.GetTerminalInfo())

	// Configure logo if available
	if logo := b.GetLogo(); logo != nil {
		view.WithLogo(logo, b.GetLogoSpacing())
	}

	// Apply state modals automatically (priority-based)
	applyStateModals(view, b)

	view.SetUseFullWidth(true)
	return view
}

// CreateStandardViewWithBreadcrumbs creates a standardized view with breadcrumb navigation.
// Breadcrumbs are displayed in the header area and show the navigation path.
//
// Example usage:
//
//	view := CreateStandardViewWithBreadcrumbs(i.BaseIntent, "Main Menu", "Settings", "Display")
func CreateStandardViewWithBreadcrumbs(b *BaseIntent, crumbs ...string) *components.StandardView {
	view := CreateStandardView(b)
	view.WithBreadcrumbs(crumbs...)
	return view
}

// applyStateModals applies modals based on BaseIntent state using priority-based display.
// Only the highest priority modal is shown:
// 1. Error (most critical, with bell)
// 2. Loading (ongoing operation)
// 3. Progress (specific progress tracking)
// 4. Success (least critical, auto-dismiss)
func applyStateModals(view *components.StandardView, base *BaseIntent) {
	// Only show the highest priority modal
	if base.HasError() {
		title := extractErrorTitle(base.GetError())
		modal := components.NewErrorModal(title, base.GetError().Error())
		view.ShowModalOverlay(modal)
	} else if base.IsLoading() {
		modal := components.NewLoadingModal(base.GetLoadingMessage(), true)
		view.ShowModalOverlay(modal)
	} else if base.IsProgressEnabled() {
		title, message, value := base.GetProgress()
		modal := components.NewProgressModal(title, message, value)
		view.ShowModalOverlay(modal)
	} else if base.ShouldShowSuccess() {
		modal := components.NewSuccessModal(base.GetSuccessMessage())
		view.ShowModalOverlay(modal)
	}
}

// extractErrorTitle extracts a meaningful title from an error.
// Attempts to extract context from the error message or type.
// Falls back to "Error" if no specific title can be extracted.
func extractErrorTitle(err error) string {
	if err == nil {
		return "Error"
	}

	errMsg := err.Error()

	// Check for common error prefixes that indicate type
	prefixes := map[string]string{
		"validation":   "Validation Error",
		"database":     "Database Error",
		"network":      "Network Error",
		"permission":   "Permission Denied",
		"not found":    "Not Found",
		"timeout":      "Timeout",
		"unauthorized": "Unauthorized",
		"invalid":      "Invalid Input",
		"failed to":    "Operation Failed",
		"unable to":    "Operation Failed",
		"cannot":       "Operation Failed",
	}

	lowerMsg := strings.ToLower(errMsg)
	for prefix, title := range prefixes {
		if strings.Contains(lowerMsg, prefix) {
			return title
		}
	}

	// Check for wrapped errors
	var unwrapped error = err
	for unwrapped != nil {
		if msg := unwrapped.Error(); msg != errMsg {
			// Try to extract from unwrapped error
			lowerUnwrapped := strings.ToLower(msg)
			for prefix, title := range prefixes {
				if strings.Contains(lowerUnwrapped, prefix) {
					return title
				}
			}
		}
		unwrapped = errors.Unwrap(unwrapped)
	}

	// Default title
	return "Error"
}

// Manual Modal Helper Functions
// These are for cases where intents need direct modal control outside of state management

// ShowErrorModal creates and attaches an error modal to the view.
// The modal is configured with a bell alert and is cancellable.
func ShowErrorModal(view *components.StandardView, err error) *components.StandardView {
	if err == nil {
		return view
	}
	title := extractErrorTitle(err)
	modal := components.NewErrorModal(title, err.Error())
	view.ShowModalOverlay(modal)
	return view
}

// ShowLoadingModal creates and attaches a loading modal to the view.
// The modal displays a spinner and optional loading message.
func ShowLoadingModal(view *components.StandardView, message string, cancellable bool) *components.StandardView {
	modal := components.NewLoadingModal(message, cancellable)
	view.ShowModalOverlay(modal)
	return view
}

// ShowProgressModal creates and attaches a progress modal to the view.
// The progress value should be between 0.0 and 1.0.
func ShowProgressModal(view *components.StandardView, title, message string, progress float64) *components.StandardView {
	modal := components.NewProgressModal(title, message, progress)
	view.ShowModalOverlay(modal)
	return view
}

// ShowSuccessModal creates and attaches a success modal to the view.
// The modal auto-dismisses after 3 seconds.
func ShowSuccessModal(view *components.StandardView, message string) *components.StandardView {
	modal := components.NewSuccessModal(message)
	view.ShowModalOverlay(modal)
	return view
}

// Footer Helper Functions
// These generate standardized help text based on TUI_STANDARDS.md keyboard shortcuts

// StandardHelpFooter formats a map of shortcuts into help text.
// Keys are the keyboard shortcuts, values are the action descriptions.
//
// Example:
//
//	shortcuts := map[string]string{
//	    "↑/k": "Up",
//	    "↓/j": "Down",
//	    "Enter": "Select",
//	}
//	footer := StandardHelpFooter(shortcuts)
func StandardHelpFooter(shortcuts map[string]string) string {
	if len(shortcuts) == 0 {
		return ""
	}

	var parts []string
	for key, action := range shortcuts {
		parts = append(parts, fmt.Sprintf("%s %s", key, action))
	}
	return strings.Join(parts, "  ")
}

// NavigationFooter returns standard navigation shortcuts.
// Used for list views, menu selections, and browsing.
func NavigationFooter() string {
	return "↑/k Up  ↓/j Down  Enter Select  Esc Back"
}

// FormFooter returns standard form navigation shortcuts.
// Used for form inputs and field navigation.
func FormFooter() string {
	return "Tab Next  Shift+Tab Previous  Enter Submit  Esc Cancel"
}

// ListFooter returns standard list view shortcuts including search.
// Used for lists with search and scroll capabilities.
func ListFooter() string {
	return "↑/k Up  ↓/j Down  Enter Select  / Search  g Top  G Bottom  Esc Back"
}

// DetailViewFooter returns standard detail view shortcuts.
// Used for viewing detailed content with scrolling.
func DetailViewFooter() string {
	return "↑/k Scroll Up  ↓/j Scroll Down  Esc Back"
}

// ModalFooter returns modal-specific action shortcuts.
// Actions are custom strings like "Enter Confirm", "Esc Cancel".
func ModalFooter(actions []string) string {
	if len(actions) == 0 {
		return "Esc Close"
	}
	return strings.Join(actions, "  ")
}

// CombineFooters combines multiple footer strings with a separator.
// Useful for combining standard shortcuts with intent-specific actions.
//
// Example:
//
//	footer := CombineFooters(NavigationFooter(), "q Quit", "m Main Menu")
//	// Result: "↑/k Up  ↓/j Down  Enter Select  Esc Back  |  q Quit  |  m Main Menu"
func CombineFooters(footers ...string) string {
	var nonEmpty []string
	for _, footer := range footers {
		if strings.TrimSpace(footer) != "" {
			nonEmpty = append(nonEmpty, footer)
		}
	}
	if len(nonEmpty) == 0 {
		return ""
	}
	return strings.Join(nonEmpty, "  |  ")
}

// UpdateLoadingRotator updates a loading message rotator and returns a tick command.
// This is a convenience function for intents that use LoadingMessageRotator.
func UpdateLoadingRotator(rotator *components.LoadingMessageRotator) {
	if rotator != nil {
		rotator.Rotate()
	}
}

// TickEvery returns a command that sends a tick message at the specified interval.
// Useful for animating loading spinners and rotating messages.
func TickEvery(d time.Duration) func() time.Duration {
	return func() time.Duration {
		return d
	}
}
