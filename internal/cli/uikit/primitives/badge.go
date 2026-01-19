package primitives

import (
	"github.com/baphled/kariya/internal/cli/uikit/theme"
	"github.com/charmbracelet/lipgloss"
)

// BadgeVariant defines the visual style and semantic meaning of a badge.
type BadgeVariant int

const (
	// BadgeDefault is for general purpose badges.
	BadgeDefault BadgeVariant = iota
	// BadgeKey is for keyboard shortcuts (e.g., "[Tab] next").
	BadgeKey
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

// Preset badge constructors for common keyboard shortcuts

// NavigateBadge creates a badge for navigation keys (arrows and vim-style j/k).
func NavigateBadge(th theme.Theme) *Badge {
	return KeyBadge("↑↓", "Navigate", th)
}

// SelectBadge creates a badge for selecting items.
func SelectBadge(th theme.Theme) *Badge {
	return KeyBadge("Enter", "Select", th)
}

// CancelBadge creates a badge for canceling.
func CancelBadge(th theme.Theme) *Badge {
	return KeyBadge("Esc", "Cancel", th)
}

// BackBadge creates a badge for going back.
func BackBadge(th theme.Theme) *Badge {
	return KeyBadge("Esc", "Back", th)
}

// ConfirmBadge creates a badge for confirming.
func ConfirmBadge(th theme.Theme) *Badge {
	return KeyBadge("Enter", "Confirm", th)
}

// QuitBadge creates a badge for quitting.
func QuitBadge(th theme.Theme) *Badge {
	return KeyBadge("q", "Quit", th)
}

// HelpBadge creates a badge for showing help.
func HelpBadge(th theme.Theme) *Badge {
	return KeyBadge("?", "Help", th)
}

// SkipBadge creates a badge for skipping.
func SkipBadge(th theme.Theme) *Badge {
	return KeyBadge("Ctrl+S", "Skip", th)
}

// RenderHelpFooter renders multiple badges as a help footer.
// Badges are separated by spaces for readability.
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
	result := ""
	for i, part := range parts {
		if i > 0 {
			result += separator
		}
		result += part
	}

	return result
}
