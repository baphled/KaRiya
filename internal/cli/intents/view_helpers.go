package intents

import (
	"errors"
	"fmt"
	"strings"
	"time"

	"github.com/charmbracelet/bubbles/key"
	tea "github.com/charmbracelet/bubbletea"

	"github.com/baphled/kariya/internal/cli/components"
<<<<<<< HEAD
	"github.com/baphled/kariya/internal/cli/themes"
=======
	"github.com/baphled/kariya/internal/cli/navigation"
>>>>>>> 88eaef9 (feat(intents): standardize BrowseTimeline key handling with HandleGlobalKeys)
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

<<<<<<< HEAD
// =============================================================================
// Theme-Aware KeyBadge Footer Functions
// =============================================================================
// These functions use the KeyBadge component for styled, consistent help footers.
// They accept a theme parameter and return professionally styled keyboard shortcuts.

// ThemedNavigationFooter returns styled navigation shortcuts using KeyBadge.
// Used for list views, menu selections, and browsing.
func ThemedNavigationFooter(theme themes.Theme) string {
	return components.RenderHelpFooter(theme,
		components.NavigateBadge(),
		components.SelectBadge(),
		components.BackBadge(),
	)
}

// ThemedFormFooter returns styled form navigation shortcuts using KeyBadge.
// Used for form inputs and field navigation.
func ThemedFormFooter(theme themes.Theme) string {
	return components.RenderHelpFooter(theme,
		components.NextBadge(),
		components.PrevBadge(),
		components.SubmitBadge(),
		components.CancelBadge(),
	)
}

// ThemedListFooter returns styled list view shortcuts including search using KeyBadge.
// Used for lists with search and scroll capabilities.
func ThemedListFooter(theme themes.Theme) string {
	return components.RenderHelpFooter(theme,
		components.NavigateBadge(),
		components.SelectBadge(),
		components.SearchBadge(),
		components.BackBadge(),
	)
}

// ThemedDetailViewFooter returns styled detail view shortcuts using KeyBadge.
// Used for viewing detailed content with scrolling.
func ThemedDetailViewFooter(theme themes.Theme) string {
	return components.RenderHelpFooter(theme,
		components.NewKeyBadge("↑/↓", "Scroll"),
		components.BackBadge(),
	)
}

// ThemedBrowseFooter returns styled browse view shortcuts using KeyBadge.
// Used for browsing lists with edit and delete capabilities.
func ThemedBrowseFooter(theme themes.Theme) string {
	return components.RenderBrowseFooter(theme)
}

// ThemedConfirmFooter returns styled confirmation shortcuts using KeyBadge.
// Used for confirmation dialogs.
func ThemedConfirmFooter(theme themes.Theme) string {
	return components.RenderConfirmFooter(theme)
}

// ThemedEditFooter returns styled edit shortcuts using KeyBadge.
// Used for edit views.
func ThemedEditFooter(theme themes.Theme) string {
	return components.RenderEditFooter(theme)
}

// ThemedExportFooter returns styled export shortcuts using KeyBadge.
// Used for export views.
func ThemedExportFooter(theme themes.Theme) string {
	return components.RenderExportFooter(theme)
}

// ThemedMenuFooter returns styled menu shortcuts using KeyBadge.
// Used for main menus.
func ThemedMenuFooter(theme themes.Theme) string {
	return components.RenderMenuFooter(theme)
}

// ThemedCustomFooter creates a custom themed footer from KeyBadges.
// Use this when standard footers don't match the required shortcuts.
//
// Example:
//
//	footer := ThemedCustomFooter(theme,
//	    components.NavigateBadge(),
//	    components.NewKeyBadge("f", "Filter"),
//	    components.NewKeyBadge("Enter", "View Details"),
//	    components.QuitBadge(),
//	)
func ThemedCustomFooter(theme themes.Theme, badges ...components.KeyBadge) string {
	return components.RenderHelpFooter(theme, badges...)
}

// ThemedGlobalBadges returns the standard global badges (Quit, Main Menu).
// Can be appended to other footers for consistency.
func ThemedGlobalBadges(theme themes.Theme) string {
	return components.RenderHelpFooter(theme,
		components.QuitBadge(),
		components.NewKeyBadge("m", "Main Menu"),
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
=======
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
	case key.Matches(msg, globalKeys.Quit):
		return KeyQuit
	case key.Matches(msg, globalKeys.Help):
		return KeyHelp
	case key.Matches(msg, globalKeys.Back):
		return KeyBack
	}

	return KeyNotHandled
}

// HandleListKeys checks if a key message matches any list navigation shortcuts.
// Returns true if the key was handled by the ListNavigationHandler.
// This is a convenience wrapper that ensures consistent list navigation.
//
// Example usage:
//
//	if keyMsg, ok := msg.(tea.KeyMsg); ok {
//	    if HandleListKeys(keyMsg, i.navHandler) {
//	        return nil
//	    }
//	}
func HandleListKeys(msg tea.KeyMsg, handler *navigation.ListNavigationHandler) bool {
	if handler == nil {
		return false
	}
	return handler.HandleKey(msg.String())
}

// GetCombinedKeyMap returns a combined keymap for help display.
// Includes global keys and optionally list or form keys.
func GetCombinedKeyMap(includeList, includeForm bool) navigation.CombinedKeyMap {
	global := navigation.DefaultGlobalKeyMap()
	list := navigation.ListKeyMap{}
	form := navigation.FormKeyMap{}

	if includeList {
		list = navigation.DefaultListKeyMap()
	}
	if includeForm {
		form = navigation.DefaultFormKeyMap()
	}

	return navigation.CombinedKeyMap{
		Global: global,
		List:   list,
		Form:   form,
	}
}

// GlobalFooter returns the standard global shortcuts footer.
// Per docs/KEYBOARD_REFERENCE.md: q=quit, ?=help, Esc=back
func GlobalFooter() string {
	return "q Quit  ? Help  Esc Back"
>>>>>>> 88eaef9 (feat(intents): standardize BrowseTimeline key handling with HandleGlobalKeys)
}
