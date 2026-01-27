package intents

import (
	"errors"
	"strings"

	"github.com/charmbracelet/bubbles/key"
	tea "github.com/charmbracelet/bubbletea"

	"github.com/baphled/kariya/internal/cli/navigation"
	"github.com/baphled/kariya/internal/cli/themes"
	"github.com/baphled/kariya/internal/cli/uikit/feedback"
	"github.com/baphled/kariya/internal/cli/uikit/layout"
	"github.com/baphled/kariya/internal/cli/uikit/primitives"
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
func CreateStandardView(b *BaseIntent) *layout.ScreenLayout {
	view := layout.NewScreenLayout(b.GetTerminalInfo())

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
func CreateStandardViewWithBreadcrumbs(b *BaseIntent, crumbs ...string) *layout.ScreenLayout {
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
func applyStateModals(view *layout.ScreenLayout, base *BaseIntent) {
	// Only show the highest priority modal
	if base.HasError() {
		title := extractErrorTitle(base.GetError())
		modal := feedback.NewErrorModal(title, base.GetError().Error())
		view.ShowModalOverlay(modal)
	} else if base.IsLoading() {
		modal := feedback.NewLoadingModal(base.GetLoadingMessage(), true)
		view.ShowModalOverlay(modal)
	} else if base.IsProgressEnabled() {
		title, message, value := base.GetProgress()
		modal := feedback.NewProgressModal(title, message, value)
		view.ShowModalOverlay(modal)
	} else if base.ShouldShowSuccess() {
		modal := feedback.NewSuccessModal(base.GetSuccessMessage())
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
// =============================================================================
// Theme-Aware KeyBadge Footer Functions
// =============================================================================
// These functions use UIKit primitives for styled, consistent help footers.
// They accept a theme parameter and return professionally styled keyboard shortcuts.

// ThemedNavigationFooter returns styled navigation shortcuts.
// Used for list views, menu selections, and browsing.
func ThemedNavigationFooter(theme themes.Theme) string {
	return primitives.RenderHelpFooter(theme,
		primitives.NavigateBadge(theme),
		primitives.SelectBadge(theme),
		primitives.BackBadge(theme),
	)
}

// ThemedFormFooter returns styled form navigation shortcuts.
// Used for form inputs and field navigation.
func ThemedFormFooter(theme themes.Theme) string {
	return primitives.RenderHelpFooter(theme,
		primitives.NextBadge(theme),
		primitives.PrevBadge(theme),
		primitives.SubmitBadge(theme),
		primitives.CancelBadge(theme),
	)
}

// ThemedListFooter returns styled list view shortcuts including search.
// Used for lists with search and scroll capabilities.
func ThemedListFooter(theme themes.Theme) string {
	return primitives.RenderHelpFooter(theme,
		primitives.NavigateBadge(theme),
		primitives.SelectBadge(theme),
		primitives.SearchBadge(theme),
		primitives.BackBadge(theme),
	)
}

// ThemedDetailViewFooter returns styled detail view shortcuts.
// Used for viewing detailed content with scrolling.
func ThemedDetailViewFooter(theme themes.Theme) string {
	return primitives.RenderHelpFooter(theme,
		primitives.HelpKeyBadge("↑/↓", "Scroll", theme),
		primitives.BackBadge(theme),
	)
}

// ThemedCustomFooter creates a custom themed footer from badges.
// Use this when standard footers don't match the required shortcuts.
//
// Example:
//
//	footer := ThemedCustomFooter(theme,
//	    primitives.NavigateBadge(theme),
//	    primitives.HelpKeyBadge("f", "Filter", theme),
//	    primitives.HelpKeyBadge("Enter", "View Details", theme),
//	    primitives.QuitBadge(theme),
//	)
func ThemedCustomFooter(theme themes.Theme, badges ...*primitives.Badge) string {
	return primitives.RenderHelpFooter(theme, badges...)
}

// ThemedGlobalBadges returns the standard global badges (Quit, Main Menu).
// Can be appended to other footers for consistency.
func ThemedGlobalBadges(theme themes.Theme) string {
	return primitives.RenderHelpFooter(theme,
		primitives.QuitBadge(theme),
		primitives.MenuBadge(theme),
	)
}

// CombineThemedFooters combines multiple themed footer strings.
// Unlike CombineFooters, this doesn't add separators as KeyBadges
// have their own visual separation.
func CombineThemedFooters(footers ...string) string {
	var nonEmpty []string
	for _, footer := range footers {
		if strings.TrimSpace(footer) != "" {
			nonEmpty = append(nonEmpty, footer)
		}
	}
	if len(nonEmpty) == 0 {
		return ""
	}
	return strings.Join(nonEmpty, "  ")
}

// ============================================================================
// Global Key Handling
// ============================================================================
// These helpers provide standardized key handling across all intents,
// implementing the keyboard shortcuts defined in docs/KEYBOARD_REFERENCE.md

// GlobalKeyResult represents the result of handling a global key
type GlobalKeyResult int

const (
	// KeyNotHandled indicates the key was not a global key
	KeyNotHandled GlobalKeyResult = iota
	// KeyQuit indicates the user wants to quit the application
	KeyQuit
	// KeyHelp indicates the user wants to see help
	KeyHelp
	// KeyBack indicates the user wants to go back
	KeyBack
)

// HandleGlobalKeys checks if a key message matches any global shortcuts.
// Returns the type of global key matched, or KeyNotHandled if no match.
// This function should be called at the beginning of each intent's Update method.
//
// Per docs/KEYBOARD_REFERENCE.md:
//   - q/ctrl+c: Quit application
//   - ?: Show context-sensitive help
//   - esc: Go back / Cancel
//
// Example usage:
//
//	func (i *MyIntent) Update(msg tea.Msg) tea.Cmd {
//	    if keyMsg, ok := msg.(tea.KeyMsg); ok {
//	        switch HandleGlobalKeys(keyMsg) {
//	        case KeyQuit:
//	            return tea.Quit
//	        case KeyHelp:
//	            i.helpModal.Toggle()
//	            return nil
//	        case KeyBack:
//	            return i.handleBack()
//	        }
//	    }
//	    // ... handle intent-specific keys
//	}
func HandleGlobalKeys(msg tea.KeyMsg) GlobalKeyResult {
	globalKeys := navigation.DefaultGlobalKeyMap()

	switch {
	// Note: Quit (q) is intentionally NOT handled here.
	// Users should only be able to quit from the main menu to prevent accidental exits.
	case key.Matches(msg, globalKeys.Help):
		return KeyHelp
	case key.Matches(msg, globalKeys.Back):
		return KeyBack
	}

	return KeyNotHandled
}

// MessageInterceptor provides a middleware layer for handling global keys before delegation.
// This ensures escape, quit, and other global keys are always processed first,
// preventing sub-components (forms, modals) from consuming them.
//
// Usage:
//
//	func (i *Intent) updateWithForm(msg tea.Msg) tea.Cmd {
//	    interceptor := NewMessageInterceptor()
//	    return interceptor.
//	        OnBack(func() tea.Cmd {
//	            i.state = previousState
//	            return nil
//	        }).
//	        OnQuit(func() tea.Cmd {
//	            return tea.Quit
//	        }).
//	        OnHelp(func() tea.Cmd {
//	            i.ToggleHelp()
//	            return nil
//	        }).
//	        InterceptOr(msg, func() tea.Cmd {
//	            // Only called if no global keys matched
//	            return i.formModel.Update(msg)
//	        })
//	}
type MessageInterceptor struct {
	backHandler GlobalKeyHandler
	quitHandler GlobalKeyHandler
	helpHandler GlobalKeyHandler
}

// GlobalKeyHandler is a function that handles a global key event.
type GlobalKeyHandler func() tea.Cmd

// NewMessageInterceptor creates a new message interceptor with no handlers.
// Use the OnBack, OnQuit, and OnHelp methods to configure behavior.
func NewMessageInterceptor() *MessageInterceptor {
	return &MessageInterceptor{}
}

// OnBack sets the handler for escape key (back navigation).
// This handler is called when the user presses Escape.
func (m *MessageInterceptor) OnBack(handler GlobalKeyHandler) *MessageInterceptor {
	m.backHandler = handler
	return m
}

// OnQuit sets the handler for quit key (q or Ctrl+C).
// This handler is called when the user wants to quit the application.
func (m *MessageInterceptor) OnQuit(handler GlobalKeyHandler) *MessageInterceptor {
	m.quitHandler = handler
	return m
}

// OnHelp sets the handler for help key (?).
// This handler is called when the user requests help.
func (m *MessageInterceptor) OnHelp(handler GlobalKeyHandler) *MessageInterceptor {
	m.helpHandler = handler
	return m
}

// InterceptOr checks for global keys and calls the appropriate handler.
// If no global key is matched, it calls the fallback function.
// This ensures global keys are always processed before sub-component delegation.
//
// Returns:
//   - tea.Cmd from the matched global key handler, OR
//   - tea.Cmd from the fallback function if no global keys matched
func (m *MessageInterceptor) InterceptOr(msg tea.Msg, fallback func() tea.Cmd) tea.Cmd {
	// Check if this is a key message
	keyMsg, ok := msg.(tea.KeyMsg)
	if !ok {
		// Not a key message, call fallback
		return fallback()
	}

	// Check for global keys
	result := HandleGlobalKeys(keyMsg)

	switch result {
	case KeyBack:
		if m.backHandler != nil {
			return m.backHandler()
		}
	case KeyQuit:
		if m.quitHandler != nil {
			return m.quitHandler()
		}
	case KeyHelp:
		if m.helpHandler != nil {
			return m.helpHandler()
		}
	}

	// No global key matched or no handler set, call fallback
	return fallback()
}

// StandardQuitHandler returns a GlobalKeyHandler that quits the application.
// This is the standard behavior for the quit key (q or Ctrl+C).
//
// Example usage:
//
//	interceptor.OnQuit(StandardQuitHandler())
func StandardQuitHandler() GlobalKeyHandler {
	return func() tea.Cmd {
		return tea.Quit
	}
}

// StandardHelpHandler creates a GlobalKeyHandler that toggles the help modal
// on a BaseIntent. This is the standard behavior for the help key (?).
//
// Example usage:
//
//	interceptor.OnHelp(StandardHelpHandler(i.BaseIntent))
func StandardHelpHandler(intent *BaseIntent) GlobalKeyHandler {
	return func() tea.Cmd {
		intent.ToggleHelp()
		return nil
	}
}
