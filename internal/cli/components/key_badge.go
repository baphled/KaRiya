package components

import (
	"strings"

	"github.com/baphled/kariya/internal/cli/themes"
	"github.com/charmbracelet/lipgloss"
)

// KeyBadge represents a styled keyboard shortcut badge.
// This is a thin helper, not a full Bubble Tea model.
type KeyBadge struct {
	Key  string
	Hint string
}

// NewKeyBadge creates a new KeyBadge with the given key and hint.
func NewKeyBadge(key, hint string) KeyBadge {
	return KeyBadge{
		Key:  key,
		Hint: hint,
	}
}

// Render renders the key badge with theme-aware styling.
func (kb KeyBadge) Render(theme themes.Theme) string {
	var keyStyle, hintStyle lipgloss.Style

	if theme != nil {
		palette := theme.Palette()
		keyStyle = lipgloss.NewStyle().
			Background(palette.BackgroundAlt).
			Foreground(palette.Primary).
			Padding(0, 1).
			Bold(true)

		hintStyle = lipgloss.NewStyle().
			Foreground(palette.ForegroundDim)
	} else {
		// Fallback styling when no theme is provided
		keyStyle = lipgloss.NewStyle().
			Padding(0, 1).
			Bold(true)

		hintStyle = lipgloss.NewStyle()
	}

	return keyStyle.Render(kb.Key) + " " + hintStyle.Render(kb.Hint)
}

// RenderHelpFooter renders multiple key badges as a help footer.
// Badges are separated by spaces for readability.
func RenderHelpFooter(theme themes.Theme, badges ...KeyBadge) string {
	if len(badges) == 0 {
		return ""
	}

	var parts []string
	for _, badge := range badges {
		parts = append(parts, badge.Render(theme))
	}

	// Join with spacing between badges
	separator := "  "
	if theme != nil {
		palette := theme.Palette()
		separatorStyle := lipgloss.NewStyle().
			Foreground(palette.ForegroundMuted)
		separator = separatorStyle.Render("  ")
	}

	return strings.Join(parts, separator)
}

// Common badge constructors for consistency across the application

// NavigateBadge returns a badge for navigation keys (arrows and vim-style j/k).
func NavigateBadge() KeyBadge {
	return KeyBadge{Key: "↑↓/jk", Hint: "Navigate"}
}

// NavigateHorizontalBadge returns a badge for horizontal navigation.
func NavigateHorizontalBadge() KeyBadge {
	return KeyBadge{Key: "←/→", Hint: "Navigate"}
}

// NavigateVimBadge returns a badge for vim-style navigation.
func NavigateVimBadge() KeyBadge {
	return KeyBadge{Key: "j/k", Hint: "Navigate"}
}

// SelectBadge returns a badge for selecting items.
func SelectBadge() KeyBadge {
	return KeyBadge{Key: "Enter", Hint: "Select"}
}

// CancelBadge returns a badge for canceling.
func CancelBadge() KeyBadge {
	return KeyBadge{Key: "Esc", Hint: "Cancel"}
}

// QuitBadge returns a badge for quitting.
func QuitBadge() KeyBadge {
	return KeyBadge{Key: "q", Hint: "Quit"}
}

// HelpBadge returns a badge for showing help.
func HelpBadge() KeyBadge {
	return KeyBadge{Key: "?", Hint: "Help"}
}

// BackBadge returns a badge for going back.
func BackBadge() KeyBadge {
	return KeyBadge{Key: "Esc", Hint: "Back"}
}

// ConfirmBadge returns a badge for confirming.
func ConfirmBadge() KeyBadge {
	return KeyBadge{Key: "Enter", Hint: "Confirm"}
}

// AddBadge returns a badge for adding.
func AddBadge() KeyBadge {
	return KeyBadge{Key: "a", Hint: "Add"}
}

// EditBadge returns a badge for editing.
func EditBadge() KeyBadge {
	return KeyBadge{Key: "e", Hint: "Edit"}
}

// DeleteBadge returns a badge for deleting.
func DeleteBadge() KeyBadge {
	return KeyBadge{Key: "d", Hint: "Delete"}
}

// SaveBadge returns a badge for saving.
func SaveBadge() KeyBadge {
	return KeyBadge{Key: "Ctrl+S", Hint: "Save"}
}

// SubmitBadge returns a badge for submitting.
func SubmitBadge() KeyBadge {
	return KeyBadge{Key: "Enter", Hint: "Submit"}
}

// NextBadge returns a badge for going to the next item.
func NextBadge() KeyBadge {
	return KeyBadge{Key: "Tab", Hint: "Next"}
}

// PrevBadge returns a badge for going to the previous item.
func PrevBadge() KeyBadge {
	return KeyBadge{Key: "Shift+Tab", Hint: "Previous"}
}

// SearchBadge returns a badge for searching.
func SearchBadge() KeyBadge {
	return KeyBadge{Key: "/", Hint: "Search"}
}

// FilterBadge returns a badge for filtering.
func FilterBadge() KeyBadge {
	return KeyBadge{Key: "f", Hint: "Filter"}
}

// Standard footer helpers for common UI patterns

// RenderMenuFooter renders a standard menu footer.
func RenderMenuFooter(theme themes.Theme) string {
	return RenderHelpFooter(theme,
		NavigateBadge(),
		SelectBadge(),
		HelpBadge(),
		QuitBadge(),
	)
}

// RenderListFooter renders a standard list footer.
func RenderListFooter(theme themes.Theme) string {
	return RenderHelpFooter(theme,
		NavigateBadge(),
		SelectBadge(),
		BackBadge(),
	)
}

// RenderFormFooter renders a standard form footer.
func RenderFormFooter(theme themes.Theme) string {
	return RenderHelpFooter(theme,
		ConfirmBadge(),
		CancelBadge(),
	)
}

// RenderEditFooter renders a standard edit footer.
func RenderEditFooter(theme themes.Theme) string {
	return RenderHelpFooter(theme,
		SaveBadge(),
		CancelBadge(),
	)
}

// RenderConfirmFooter renders a standard confirmation footer.
func RenderConfirmFooter(theme themes.Theme) string {
	return RenderHelpFooter(theme,
		ConfirmBadge(),
		CancelBadge(),
	)
}

// RenderBrowseFooter renders a footer for browse views.
func RenderBrowseFooter(theme themes.Theme) string {
	return RenderHelpFooter(theme,
		NavigateBadge(),
		SelectBadge(),
		EditBadge(),
		DeleteBadge(),
		BackBadge(),
	)
}

// RenderExportFooter renders a footer for export views.
func RenderExportFooter(theme themes.Theme) string {
	return RenderHelpFooter(theme,
		NavigateBadge(),
		SelectBadge(),
		ConfirmBadge(),
		BackBadge(),
	)
}
