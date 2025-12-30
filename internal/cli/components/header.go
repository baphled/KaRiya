package components

import (
	"strings"

	"github.com/charmbracelet/lipgloss"

	"github.com/baphled/kariya/internal/cli/styles"
)

// HeaderModel renders a consistent screen header with title, subtitle, and breadcrumbs
type HeaderModel struct {
	title       string
	subtitle    string
	breadcrumbs []string // Navigation path: ["Home", "List", "Event Details"]
	width       int
	height      int
	showBorder  bool
}

// NewHeader creates a new header with a title
func NewHeader(title string, width int) HeaderModel {
	return HeaderModel{
		title:       title,
		subtitle:    "",
		breadcrumbs: []string{},
		width:       width,
		height:      1,
		showBorder:  false,
	}
}

// SetSubtitle sets the subtitle
func (h *HeaderModel) SetSubtitle(subtitle string) {
	h.subtitle = subtitle
}

// SetBreadcrumbs sets the navigation breadcrumbs
func (h *HeaderModel) SetBreadcrumbs(crumbs []string) {
	h.breadcrumbs = crumbs
}

// AddBreadcrumb adds a breadcrumb to the end
func (h *HeaderModel) AddBreadcrumb(crumb string) {
	h.breadcrumbs = append(h.breadcrumbs, crumb)
}

// SetWidth sets the header width
func (h *HeaderModel) SetWidth(width int) {
	h.width = width
}

// SetHeight sets the header height
func (h *HeaderModel) SetHeight(height int) {
	h.height = height
}

// SetShowBorder sets whether to show a border
func (h *HeaderModel) SetShowBorder(show bool) {
	h.showBorder = show
}

// ClearBreadcrumbs clears all breadcrumbs
func (h *HeaderModel) ClearBreadcrumbs() {
	h.breadcrumbs = []string{}
}

// GetBreadcrumbs returns current breadcrumbs
func (h HeaderModel) GetBreadcrumbs() []string {
	return h.breadcrumbs
}

// GetTitle returns the title
func (h HeaderModel) GetTitle() string {
	return h.title
}

// GetSubtitle returns the subtitle
func (h HeaderModel) GetSubtitle() string {
	return h.subtitle
}

// View renders the header
func (h HeaderModel) View() string {
	if h.width <= 0 {
		return ""
	}

	var parts []string

	// Render breadcrumbs if present
	if len(h.breadcrumbs) > 0 {
		breadcrumbStr := h.renderBreadcrumbs()
		parts = append(parts, breadcrumbStr)
	}

	// Render title
	titleStr := h.renderTitle()
	parts = append(parts, titleStr)

	// Render subtitle if present
	if h.subtitle != "" {
		subtitleStr := h.renderSubtitle()
		parts = append(parts, subtitleStr)
	}

	content := strings.Join(parts, "\n")

	if h.showBorder {
		return styles.WithBorder(lipgloss.NewStyle()).Render(content)
	}

	return content
}

// renderBreadcrumbs renders the breadcrumb navigation
func (h HeaderModel) renderBreadcrumbs() string {
	// Join breadcrumbs with separator
	separator := styles.InfoHint.Render(" > ")
	styledCrumbs := make([]string, len(h.breadcrumbs))

	for i, crumb := range h.breadcrumbs {
		if i == len(h.breadcrumbs)-1 {
			// Last breadcrumb - highlight it
			styledCrumbs[i] = lipgloss.NewStyle().
				Foreground(styles.ColorAccentTeal).
				Bold(true).
				Render(crumb)
		} else {
			// Previous breadcrumbs - muted
			styledCrumbs[i] = lipgloss.NewStyle().
				Foreground(styles.ColorTextMuted).
				Render(crumb)
		}
	}

	breadcrumbStr := strings.Join(styledCrumbs, separator)

	// Add bottom margin
	return breadcrumbStr + "\n"
}

// renderTitle renders the main title
func (h HeaderModel) renderTitle() string {
	titleStyle := styles.HeaderMain

	// Truncate title if it's too long
	maxWidth := h.width - 4 // Account for potential padding
	if len(h.title) > maxWidth {
		title := h.title[:maxWidth-3] + "..."
		return titleStyle.Render(title)
	}

	return titleStyle.Render(h.title)
}

// renderSubtitle renders the subtitle
func (h HeaderModel) renderSubtitle() string {
	subtitleStyle := lipgloss.NewStyle().
		Foreground(styles.ColorTextSecondary).
		Italic(true).
		MarginTop(1)

	// Truncate subtitle if it's too long
	maxWidth := h.width - 4
	if len(h.subtitle) > maxWidth {
		subtitle := h.subtitle[:maxWidth-3] + "..."
		return subtitleStyle.Render(subtitle)
	}

	return subtitleStyle.Render(h.subtitle)
}

