package themes

import (
	"github.com/charmbracelet/glamour"
)

// NewGlamourStyleName returns the appropriate glamour style name based on the theme.
// Returns "dark" for dark themes and "light" for light themes.
func NewGlamourStyleName(theme Theme) string {
	if theme == nil {
		return "dark"
	}

	if theme.IsDark() {
		return "dark"
	}
	return "light"
}

// RenderMarkdown renders markdown content with theme-aware styling.
// Uses Glamour for rendering with the appropriate color scheme.
func RenderMarkdown(theme Theme, content string, width int) (string, error) {
	if content == "" {
		return "", nil
	}

	styleName := NewGlamourStyleName(theme)

	renderer, err := glamour.NewTermRenderer(
		glamour.WithStylePath(styleName),
		glamour.WithWordWrap(width),
	)
	if err != nil {
		return "", err
	}

	return renderer.Render(content)
}

// RenderCVPreview renders CV content as formatted markdown for preview.
// This is a convenience wrapper for RenderMarkdown specifically for CV content.
func RenderCVPreview(theme Theme, cvContent string, width int) (string, error) {
	return RenderMarkdown(theme, cvContent, width)
}

// MarkdownRenderer is a helper struct for rendering markdown with consistent theme styling.
type MarkdownRenderer struct {
	theme    Theme
	width    int
	renderer *glamour.TermRenderer
}

// NewMarkdownRenderer creates a new MarkdownRenderer with the given theme and width.
func NewMarkdownRenderer(theme Theme, width int) (*MarkdownRenderer, error) {
	styleName := NewGlamourStyleName(theme)

	renderer, err := glamour.NewTermRenderer(
		glamour.WithStylePath(styleName),
		glamour.WithWordWrap(width),
	)
	if err != nil {
		return nil, err
	}

	return &MarkdownRenderer{
		theme:    theme,
		width:    width,
		renderer: renderer,
	}, nil
}

// Render renders markdown content.
func (mr *MarkdownRenderer) Render(content string) (string, error) {
	if content == "" {
		return "", nil
	}
	return mr.renderer.Render(content)
}

// SetWidth updates the renderer width.
// Note: This creates a new internal renderer.
func (mr *MarkdownRenderer) SetWidth(width int) error {
	mr.width = width

	styleName := NewGlamourStyleName(mr.theme)
	renderer, err := glamour.NewTermRenderer(
		glamour.WithStylePath(styleName),
		glamour.WithWordWrap(width),
	)
	if err != nil {
		return err
	}

	mr.renderer = renderer
	return nil
}

// SetTheme updates the renderer theme.
// Note: This creates a new internal renderer.
func (mr *MarkdownRenderer) SetTheme(theme Theme) error {
	mr.theme = theme

	styleName := NewGlamourStyleName(theme)
	renderer, err := glamour.NewTermRenderer(
		glamour.WithStylePath(styleName),
		glamour.WithWordWrap(mr.width),
	)
	if err != nil {
		return err
	}

	mr.renderer = renderer
	return nil
}
