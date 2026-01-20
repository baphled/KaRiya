// Package primitives provides foundational UI components for the UIKit.
// All primitives are theme-aware and follow a fluent builder pattern.
package primitives

import (
	"github.com/baphled/kariya/internal/cli/uikit/theme"
	"github.com/charmbracelet/lipgloss"
)

// TextStyle defines the semantic style of text.
type TextStyle int

const (
	// TextBody is the default text style for regular content.
	TextBody TextStyle = iota
	// TextTitle is for large, prominent titles.
	TextTitle
	// TextSubtitle is for section subtitles.
	TextSubtitle
	// TextMuted is for disabled or de-emphasized text.
	TextMuted
	// TextError is for error messages.
	TextError
	// TextSuccess is for success messages.
	TextSuccess
	// TextWarning is for warning messages.
	TextWarning
)

// Text is a theme-aware text component with fluent API.
// It supports semantic styles, bold formatting, width constraints, and alignment.
//
// Example:
//
//	title := primitives.Title("Welcome", theme).Bold().Align(lipgloss.Center)
//	error := primitives.ErrorText("Failed", theme).Width(40)
type Text struct {
	theme.Aware
	content   string
	textStyle TextStyle
	bold      bool
	width     int
	align     lipgloss.Position
}

// NewText creates a new text component with the given content and theme.
// If theme is nil, the default theme is used.
func NewText(content string, th theme.Theme) *Text {
	t := &Text{
		content:   content,
		textStyle: TextBody,
		bold:      false,
		width:     0, // 0 means no width constraint
	}
	if th != nil {
		t.SetTheme(th)
	}
	return t
}

// Style sets the semantic style of the text.
// Returns the text for method chaining.
func (t *Text) Style(style TextStyle) *Text {
	t.textStyle = style
	return t
}

// Bold makes the text bold.
// Returns the text for method chaining.
func (t *Text) Bold() *Text {
	t.bold = true
	return t
}

// Width sets the maximum width of the text.
// Text will be wrapped if it exceeds this width.
// Returns the text for method chaining.
func (t *Text) Width(w int) *Text {
	t.width = w
	return t
}

// Align sets the text alignment (lipgloss.Left, lipgloss.Center, lipgloss.Right).
// Returns the text for method chaining.
func (t *Text) Align(align lipgloss.Position) *Text {
	t.align = align
	return t
}

// Render returns the styled text as a string.
func (t *Text) Render() string {
	style := t.buildStyle()
	return style.Render(t.content)
}

// buildStyle creates a lipgloss style based on the text configuration.
func (t *Text) buildStyle() lipgloss.Style {
	style := lipgloss.NewStyle()

	// Apply semantic color based on text style
	switch t.textStyle {
	case TextTitle:
		style = style.Foreground(t.PrimaryColor()).Bold(true)
	case TextSubtitle:
		style = style.Foreground(t.SecondaryColor())
	case TextBody:
		style = style.Foreground(t.Theme().ForegroundColor())
	case TextMuted:
		style = style.Foreground(t.MutedColor())
	case TextError:
		style = style.Foreground(t.ErrorColor())
	case TextSuccess:
		style = style.Foreground(t.SuccessColor())
	case TextWarning:
		style = style.Foreground(t.WarningColor())
	}

	// Apply bold if requested
	if t.bold {
		style = style.Bold(true)
	}

	// Apply width constraint if set
	if t.width > 0 {
		style = style.Width(t.width)
	}

	// Apply alignment if width is set (alignment requires width)
	if t.width > 0 && t.align != 0 {
		style = style.Align(t.align)
	}

	return style
}

// Convenience constructors for common text styles

// Title creates a title-styled text component.
func Title(content string, th theme.Theme) *Text {
	return NewText(content, th).Style(TextTitle)
}

// Subtitle creates a subtitle-styled text component.
func Subtitle(content string, th theme.Theme) *Text {
	return NewText(content, th).Style(TextSubtitle)
}

// Body creates a body-styled text component.
func Body(content string, th theme.Theme) *Text {
	return NewText(content, th).Style(TextBody)
}

// Muted creates a muted-styled text component.
func Muted(content string, th theme.Theme) *Text {
	return NewText(content, th).Style(TextMuted)
}

// ErrorText creates an error-styled text component.
func ErrorText(content string, th theme.Theme) *Text {
	return NewText(content, th).Style(TextError)
}

// SuccessText creates a success-styled text component.
func SuccessText(content string, th theme.Theme) *Text {
	return NewText(content, th).Style(TextSuccess)
}

// WarningText creates a warning-styled text component.
func WarningText(content string, th theme.Theme) *Text {
	return NewText(content, th).Style(TextWarning)
}
