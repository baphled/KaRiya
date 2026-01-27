package primitives

import (
	"strings"

	"github.com/baphled/kariya/internal/cli/uikit/theme"
	"github.com/charmbracelet/lipgloss"
)

// BadgeVariant defines the visual style and semantic meaning of a badge.
type BadgeVariant int

const (
	// BadgeDefault is for general purpose badges.
	BadgeDefault BadgeVariant = iota
	// BadgeKey is for keyboard shortcuts compact style (e.g., "[Esc → Back]").
	BadgeKey
	// BadgeHelpKey is for help footer shortcuts (e.g., "[Esc] Back" - two-part styling).
	BadgeHelpKey
	// BadgeStatus is for status indicators (e.g., "Active", "Pending").
	BadgeStatus
	// BadgeTag is for tags and labels (e.g., "Feature", "Bug").
	BadgeTag
)

// Badge is a theme-aware badge component for displaying small pieces of information.
// Badges support multiple variants for different use cases.
//
// Example:
//
//	keyBadge := primitives.KeyBadge("Esc", "cancel", theme)
//	statusBadge := primitives.StatusBadge("Active", theme)
//	tagBadge := primitives.TagBadge("Feature", theme)
type Badge struct {
	theme.Aware
	label   string
	value   string
	variant BadgeVariant
}

// NewBadge creates a new badge with the given label and theme.
// If theme is nil, the default theme is used.
// The badge defaults to Default variant.
func NewBadge(label string, th theme.Theme) *Badge {
	b := &Badge{
		label:   label,
		value:   "",
		variant: BadgeDefault,
	}
	if th != nil {
		b.SetTheme(th)
	}
	return b
}

// Value sets the value for the badge (used for key badges).
// Returns the badge for method chaining.
func (b *Badge) Value(value string) *Badge {
	b.value = value
	return b
}

// Variant sets the badge variant (Default, Key, Status, Tag).
// Returns the badge for method chaining.
func (b *Badge) Variant(v BadgeVariant) *Badge {
	b.variant = v
	return b
}

// Render returns the styled badge as a string.
func (b *Badge) Render() string {
	// BadgeHelpKey has special two-part rendering
	if b.variant == BadgeHelpKey {
		return b.renderHelpKey()
	}

	style := b.buildStyle()

	// Format content based on variant
	var content string
	switch b.variant {
	case BadgeKey:
		// Key badge: [key] or [key → action]
		if b.value != "" {
			content = "[" + b.label + " → " + b.value + "]"
		} else {
			content = "[" + b.label + "]"
		}
	case BadgeStatus, BadgeTag, BadgeDefault:
		content = b.label
	}

	return style.Render(content)
}

// renderHelpKey renders a two-part help key badge: "[Key] Hint"
func (b *Badge) renderHelpKey() string {
	th := b.Theme()

	var keyStyle, hintStyle lipgloss.Style
	if th != nil {
		palette := th.Palette()
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

	return keyStyle.Render(b.label) + " " + hintStyle.Render(b.value)
}

// buildStyle creates a lipgloss style based on the badge configuration.
func (b *Badge) buildStyle() lipgloss.Style {
	style := lipgloss.NewStyle()

	switch b.variant {
	case BadgeKey:
		// Key badges: Accent color, bold, no background
		style = style.
			Foreground(b.AccentColor()).
			Bold(true)

	case BadgeStatus:
		// Status badges: Colored text with subtle background
		style = style.
			Foreground(b.SuccessColor()).
			Background(b.Theme().BackgroundColor()).
			Padding(0, 1).
			Border(lipgloss.RoundedBorder()).
			BorderForeground(b.SuccessColor())

	case BadgeTag:
		// Tag badges: Pill-style with background
		style = style.
			Foreground(b.Theme().BackgroundColor()).
			Background(b.SecondaryColor()).
			Padding(0, 2).
			Border(lipgloss.RoundedBorder(), false).
			Bold(false)

	case BadgeDefault:
		// Default badges: Muted style
		style = style.
			Foreground(b.MutedColor()).
			Padding(0, 1)
	}

	return style
}

// Convenience constructors for common badge types

// KeyBadge creates a key badge formatted as "[key → action]".
// If action is empty, formats as "[key]".
func KeyBadge(key, action string, th theme.Theme) *Badge {
	return NewBadge(key, th).Value(action).Variant(BadgeKey)
}

// StatusBadge creates a status badge with success styling.
func StatusBadge(status string, th theme.Theme) *Badge {
	return NewBadge(status, th).Variant(BadgeStatus)
}

// TagBadge creates a tag badge with pill styling.
func TagBadge(tag string, th theme.Theme) *Badge {
	return NewBadge(tag, th).Variant(BadgeTag)
}

// HelpKeyBadge creates a help footer key badge with two-part styling: "[Key] Hint".
// This is the standard format for keyboard shortcuts in help footers.
func HelpKeyBadge(key, hint string, th theme.Theme) *Badge {
	return NewBadge(key, th).Value(hint).Variant(BadgeHelpKey)
}

// =============================================================================
// Common Help Key Badge Constructors
// =============================================================================
// These create pre-configured badges for common keyboard shortcuts.
// Use these for consistency across the application.

// NavigateBadge returns a badge for navigation keys (arrows and vim-style j/k).
func NavigateBadge(th theme.Theme) *Badge {
	return HelpKeyBadge("↑↓/jk", "Navigate", th)
}

// NavigateHorizontalBadge returns a badge for horizontal navigation.
func NavigateHorizontalBadge(th theme.Theme) *Badge {
	return HelpKeyBadge("←/→/hl", "Navigate", th)
}

// NavigateVimBadge returns a badge for vim-style navigation.
func NavigateVimBadge(th theme.Theme) *Badge {
	return HelpKeyBadge("j/k", "Navigate", th)
}

// SelectBadge returns a badge for selecting items.
func SelectBadge(th theme.Theme) *Badge {
	return HelpKeyBadge("Enter", "Select", th)
}

// CancelBadge returns a badge for canceling.
func CancelBadge(th theme.Theme) *Badge {
	return HelpKeyBadge("Esc", "Cancel", th)
}

// QuitBadge returns a badge for quitting.
func QuitBadge(th theme.Theme) *Badge {
	return HelpKeyBadge("q", "Quit", th)
}

// HelpBadge returns a badge for showing help.
func HelpBadge(th theme.Theme) *Badge {
	return HelpKeyBadge("?", "Help", th)
}

// BackBadge returns a badge for going back.
func BackBadge(th theme.Theme) *Badge {
	return HelpKeyBadge("Esc", "Back", th)
}

// ConfirmBadge returns a badge for confirming.
func ConfirmBadge(th theme.Theme) *Badge {
	return HelpKeyBadge("Enter", "Confirm", th)
}

// AddBadge returns a badge for adding.
func AddBadge(th theme.Theme) *Badge {
	return HelpKeyBadge("a", "Add", th)
}

// EditBadge returns a badge for editing.
func EditBadge(th theme.Theme) *Badge {
	return HelpKeyBadge("e", "Edit", th)
}

// DeleteBadge returns a badge for deleting.
func DeleteBadge(th theme.Theme) *Badge {
	return HelpKeyBadge("d", "Delete", th)
}

// SaveBadge returns a badge for saving.
func SaveBadge(th theme.Theme) *Badge {
	return HelpKeyBadge("Ctrl+S", "Save", th)
}

// SkipBadge returns a badge for skipping.
func SkipBadge(th theme.Theme) *Badge {
	return HelpKeyBadge("Ctrl+S", "Skip", th)
}

// SubmitBadge returns a badge for submitting.
func SubmitBadge(th theme.Theme) *Badge {
	return HelpKeyBadge("Enter", "Submit", th)
}

// NextBadge returns a badge for going to the next item.
func NextBadge(th theme.Theme) *Badge {
	return HelpKeyBadge("Tab", "Next", th)
}

// NextFieldBadge returns a badge for going to the next form field.
func NextFieldBadge(th theme.Theme) *Badge {
	return HelpKeyBadge("Tab", "Next field", th)
}

// PrevBadge returns a badge for going to the previous item.
func PrevBadge(th theme.Theme) *Badge {
	return HelpKeyBadge("Shift+Tab", "Previous", th)
}

// ApplyBadge returns a badge for applying changes.
func ApplyBadge(th theme.Theme) *Badge {
	return HelpKeyBadge("Enter", "Apply", th)
}

// SearchBadge returns a badge for searching.
func SearchBadge(th theme.Theme) *Badge {
	return HelpKeyBadge("/", "Search", th)
}

// FilterBadge returns a badge for filtering.
func FilterBadge(th theme.Theme) *Badge {
	return HelpKeyBadge("f", "Filter", th)
}

// YesBadge returns a badge for yes/confirm shortcut.
func YesBadge(th theme.Theme) *Badge {
	return HelpKeyBadge("y", "Yes", th)
}

// NoBadge returns a badge for no/cancel shortcut.
func NoBadge(th theme.Theme) *Badge {
	return HelpKeyBadge("n", "No", th)
}

// ToggleBadge returns a badge for toggle selection (left/right).
func ToggleBadge(th theme.Theme) *Badge {
	return HelpKeyBadge("←→/hl", "Toggle", th)
}

// RetryBadge returns a badge for retry action (r key only).
func RetryBadge(th theme.Theme) *Badge {
	return HelpKeyBadge("r", "Retry", th)
}

// RetryEnterBadge returns a badge for retry action (Enter or r key).
func RetryEnterBadge(th theme.Theme) *Badge {
	return HelpKeyBadge("Enter/r", "Retry", th)
}

// ContinueBadge returns a badge for continue action.
func ContinueBadge(th theme.Theme) *Badge {
	return HelpKeyBadge("Enter", "Continue", th)
}

// =============================================================================
// Help Footer Rendering
// =============================================================================

// RenderHelpFooter renders multiple badges as a help footer.
// Badges are separated by styled spaces for readability.
func RenderHelpFooter(th theme.Theme, badges ...*Badge) string {
	if len(badges) == 0 {
		return ""
	}

	var parts []string
	for _, badge := range badges {
		parts = append(parts, badge.Render())
	}

	// Join with spacing between badges
	separator := "  "
	if th != nil {
		palette := th.Palette()
		separatorStyle := lipgloss.NewStyle().
			Foreground(palette.ForegroundMuted)
		separator = separatorStyle.Render("  ")
	}

	return strings.Join(parts, separator)
}

// =============================================================================
// Standard Footer Presets
// =============================================================================

// RenderMenuFooter renders a standard menu footer.
func RenderMenuFooter(th theme.Theme) string {
	return RenderHelpFooter(th,
		NavigateBadge(th),
		SelectBadge(th),
		HelpBadge(th),
		QuitBadge(th),
	)
}

// RenderListFooter renders a standard list footer.
func RenderListFooter(th theme.Theme) string {
	return RenderHelpFooter(th,
		NavigateBadge(th),
		SelectBadge(th),
		BackBadge(th),
	)
}

// RenderFormFooter renders a standard form footer.
func RenderFormFooter(th theme.Theme) string {
	return RenderHelpFooter(th,
		ConfirmBadge(th),
		CancelBadge(th),
	)
}

// RenderEditFooter renders a standard edit footer.
func RenderEditFooter(th theme.Theme) string {
	return RenderHelpFooter(th,
		SaveBadge(th),
		CancelBadge(th),
	)
}

// RenderConfirmFooter renders a standard confirmation footer.
func RenderConfirmFooter(th theme.Theme) string {
	return RenderHelpFooter(th,
		ConfirmBadge(th),
		CancelBadge(th),
	)
}

// RenderBrowseFooter renders a footer for browse views.
func RenderBrowseFooter(th theme.Theme) string {
	return RenderHelpFooter(th,
		NavigateBadge(th),
		SelectBadge(th),
		EditBadge(th),
		DeleteBadge(th),
		BackBadge(th),
	)
}

// RenderExportFooter renders a footer for export views.
func RenderExportFooter(th theme.Theme) string {
	return RenderHelpFooter(th,
		NavigateBadge(th),
		SelectBadge(th),
		ConfirmBadge(th),
		BackBadge(th),
	)
}
