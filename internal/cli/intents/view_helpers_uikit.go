package intents

import (
	"github.com/baphled/kariya/internal/cli/themes"
	"github.com/baphled/kariya/internal/cli/uikit/feedback"
	"github.com/baphled/kariya/internal/cli/uikit/layout"
	"github.com/baphled/kariya/internal/cli/uikit/primitives"
)

// =============================================================================
// UIKit View Helpers
// =============================================================================
// These functions provide UIKit-based alternatives to the legacy view helpers.
// Use these when migrating intents from components to UIKit.
//
// Migration guide:
//   - CreateStandardView() -> CreateScreenLayout()
//   - CreateStandardViewWithBreadcrumbs() -> CreateScreenLayoutWithBreadcrumbs()
//   - ThemedCustomFooter(theme, badges...) -> UIKitCustomFooter(theme, badges...)
//   - components.NewKeyBadge("k", "hint") -> primitives.HelpKeyBadge("k", "hint", theme)

// CreateScreenLayout creates a UIKit ScreenLayout with logo and automatic state modals.
// This is the UIKit equivalent of CreateStandardView.
//
// The view is configured with:
// - Terminal info from BaseIntent
// - Logo (if available) with configured spacing
// - Automatic modal display based on intent state (error, loading, progress, success)
// - Full width content rendering
//
// Example usage:
//
//	func (i *MyIntent) View() string {
//	    view := CreateScreenLayout(i.BaseIntent)
//	    view.WithContent(i.renderContent())
//	    view.WithHelp(UIKitCustomFooter(theme, primitives.QuitBadge(theme)))
//	    return view.Render()
//	}
func CreateScreenLayout(b *BaseIntent) *layout.ScreenLayout {
	view := layout.NewScreenLayout(b.GetTerminalInfo())

	// Configure logo if available
	if logo := b.GetLogo(); logo != nil {
		// LogoModel satisfies layout.LogoRenderer (both have ViewStatic and SetWidth)
		view.WithLogo(logo, b.GetLogoSpacing())
	}

	// Apply state modals automatically (priority-based)
	applyUIKitStateModals(view, b)

	view.SetUseFullWidth(true)

	// Apply theme if available
	if theme := b.Theme(); theme != nil {
		view.WithTheme(theme)
	}

	return view
}

// CreateScreenLayoutWithBreadcrumbs creates a UIKit ScreenLayout with breadcrumb navigation.
// This is the UIKit equivalent of CreateStandardViewWithBreadcrumbs.
//
// Breadcrumbs are displayed in the header area and show the navigation path.
//
// Example usage:
//
//	view := CreateScreenLayoutWithBreadcrumbs(i.BaseIntent, "Main Menu", "Settings", "Display")
func CreateScreenLayoutWithBreadcrumbs(b *BaseIntent, crumbs ...string) *layout.ScreenLayout {
	view := CreateScreenLayout(b)
	view.WithBreadcrumbs(crumbs...)
	return view
}

// applyUIKitStateModals applies UIKit modals based on BaseIntent state.
// Only the highest priority modal is shown:
// 1. Error (most critical, with bell)
// 2. Loading (ongoing operation)
// 3. Progress (specific progress tracking)
// 4. Success (least critical, auto-dismiss)
func applyUIKitStateModals(view *layout.ScreenLayout, base *BaseIntent) {
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

// =============================================================================
// UIKit Modal Helper Functions
// =============================================================================
// These provide UIKit-based modal helpers for direct modal control.

// ShowUIKitErrorModal creates and attaches a UIKit error modal to the layout.
func ShowUIKitErrorModal(view *layout.ScreenLayout, err error) *layout.ScreenLayout {
	if err == nil {
		return view
	}
	title := extractErrorTitle(err)
	modal := feedback.NewErrorModal(title, err.Error())
	view.ShowModalOverlay(modal)
	return view
}

// ShowUIKitLoadingModal creates and attaches a UIKit loading modal to the layout.
func ShowUIKitLoadingModal(view *layout.ScreenLayout, message string, cancellable bool) *layout.ScreenLayout {
	modal := feedback.NewLoadingModal(message, cancellable)
	view.ShowModalOverlay(modal)
	return view
}

// ShowUIKitProgressModal creates and attaches a UIKit progress modal to the layout.
func ShowUIKitProgressModal(view *layout.ScreenLayout, title, message string, progress float64) *layout.ScreenLayout {
	modal := feedback.NewProgressModal(title, message, progress)
	view.ShowModalOverlay(modal)
	return view
}

// ShowUIKitSuccessModal creates and attaches a UIKit success modal to the layout.
func ShowUIKitSuccessModal(view *layout.ScreenLayout, message string) *layout.ScreenLayout {
	modal := feedback.NewSuccessModal(message)
	view.ShowModalOverlay(modal)
	return view
}

// =============================================================================
// UIKit Footer Helper Functions
// =============================================================================
// These use primitives.Badge for styled, theme-aware help footers.

// UIKitNavigationFooter returns styled navigation shortcuts using UIKit badges.
func UIKitNavigationFooter(theme themes.Theme) string {
	return primitives.RenderHelpFooter(theme,
		primitives.NavigateBadge(theme),
		primitives.SelectBadge(theme),
		primitives.BackBadge(theme),
	)
}

// UIKitFormFooter returns styled form navigation shortcuts using UIKit badges.
func UIKitFormFooter(theme themes.Theme) string {
	return primitives.RenderHelpFooter(theme,
		primitives.NextBadge(theme),
		primitives.PrevBadge(theme),
		primitives.SubmitBadge(theme),
		primitives.CancelBadge(theme),
	)
}

// UIKitListFooter returns styled list view shortcuts using UIKit badges.
func UIKitListFooter(theme themes.Theme) string {
	return primitives.RenderHelpFooter(theme,
		primitives.NavigateBadge(theme),
		primitives.SelectBadge(theme),
		primitives.SearchBadge(theme),
		primitives.BackBadge(theme),
	)
}

// UIKitDetailViewFooter returns styled detail view shortcuts using UIKit badges.
func UIKitDetailViewFooter(theme themes.Theme) string {
	return primitives.RenderHelpFooter(theme,
		primitives.HelpKeyBadge("↑/↓", "Scroll", theme),
		primitives.BackBadge(theme),
	)
}

// UIKitBrowseFooter returns styled browse view shortcuts using UIKit badges.
func UIKitBrowseFooter(theme themes.Theme) string {
	return primitives.RenderBrowseFooter(theme)
}

// UIKitConfirmFooter returns styled confirmation shortcuts using UIKit badges.
func UIKitConfirmFooter(theme themes.Theme) string {
	return primitives.RenderConfirmFooter(theme)
}

// UIKitEditFooter returns styled edit shortcuts using UIKit badges.
func UIKitEditFooter(theme themes.Theme) string {
	return primitives.RenderEditFooter(theme)
}

// UIKitExportFooter returns styled export shortcuts using UIKit badges.
func UIKitExportFooter(theme themes.Theme) string {
	return primitives.RenderExportFooter(theme)
}

// UIKitMenuFooter returns styled menu shortcuts using UIKit badges.
func UIKitMenuFooter(theme themes.Theme) string {
	return primitives.RenderMenuFooter(theme)
}

// UIKitCustomFooter creates a custom themed footer from UIKit badges.
// This is the UIKit equivalent of ThemedCustomFooter.
//
// Example:
//
//	footer := UIKitCustomFooter(theme,
//	    primitives.NavigateBadge(theme),
//	    primitives.HelpKeyBadge("f", "Filter", theme),
//	    primitives.HelpKeyBadge("Enter", "View Details", theme),
//	    primitives.QuitBadge(theme),
//	)
func UIKitCustomFooter(theme themes.Theme, badges ...*primitives.Badge) string {
	return primitives.RenderHelpFooter(theme, badges...)
}

// UIKitGlobalBadges returns the standard global badges (Quit, Menu).
func UIKitGlobalBadges(theme themes.Theme) string {
	return primitives.RenderHelpFooter(theme,
		primitives.QuitBadge(theme),
		primitives.HelpKeyBadge("m", "Main Menu", theme),
	)
}
