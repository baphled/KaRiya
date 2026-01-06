package components

import (
	"strings"

	"github.com/charmbracelet/lipgloss"

	"github.com/baphled/kariya/internal/cli/styles"
)

// IntentHeader renders a consistent header for intent screens
type IntentHeader struct {
	appName     string
	breadcrumbs []Breadcrumb
	description string
	width       int
	height      int
}

// NewIntentHeader creates a new intent header
func NewIntentHeader(width int) *IntentHeader {
	return &IntentHeader{
		appName:     "KaRiya",
		breadcrumbs: []Breadcrumb{},
		description: "",
		width:       width,
		height:      4, // Fixed height for consistency
	}
}

// SetBreadcrumbs sets the breadcrumb trail
func (h *IntentHeader) SetBreadcrumbs(crumbs []Breadcrumb) *IntentHeader {
	h.breadcrumbs = crumbs
	return h
}

// SetDescription sets the intent description
func (h *IntentHeader) SetDescription(description string) *IntentHeader {
	h.description = description
	return h
}

// SetWidth sets the header width
func (h *IntentHeader) SetWidth(width int) *IntentHeader {
	h.width = width
	return h
}

// GetHeight returns the fixed height of the header
func (h *IntentHeader) GetHeight() int {
	return h.height
}

// View renders the intent header
func (h *IntentHeader) View() string {
	// Build the header content
	var lines []string

	// Line 1: App name and breadcrumbs
	appNameStyle := lipgloss.NewStyle().
		Foreground(styles.ColorAccentTeal).
		Bold(true)

	appNameText := appNameStyle.Render(h.appName)

	// Build breadcrumb trail
	var breadcrumbText string
	if len(h.breadcrumbs) > 0 {
		bar := NewBreadcrumbBar(h.width-20, false) // Reserve space for app name
		bar.SetCrumbs(h.breadcrumbs)
		breadcrumbText = bar.View()
	}

	// Combine app name and breadcrumbs with separator
	if breadcrumbText != "" {
		separator := lipgloss.NewStyle().
			Foreground(styles.ColorTextMuted).
			Render("  •  ")
		line1 := "  " + appNameText + separator + breadcrumbText
		lines = append(lines, line1)
	} else {
		lines = append(lines, "  "+appNameText)
	}

	// Line 2: Separator line
	separatorStyle := lipgloss.NewStyle().
		Foreground(styles.ColorBorder)

	// Create a separator that fits the width
	separatorWidth := h.width - 4
	if separatorWidth < 0 {
		separatorWidth = h.width
	}
	separatorLine := separatorStyle.Render(strings.Repeat("─", separatorWidth))
	lines = append(lines, separatorLine)

	// Line 3: Description
	if h.description != "" {
		descStyle := lipgloss.NewStyle().
			Foreground(styles.ColorTextSecondary).
			Italic(true)

		descText := "  " + descStyle.Render(h.description)
		lines = append(lines, descText)
	} else {
		lines = append(lines, "")
	}

	// Line 4: Empty line for spacing
	lines = append(lines, "")

	// Apply border
	content := strings.Join(lines, "\n")

	boxStyle := lipgloss.NewStyle().
		Border(lipgloss.RoundedBorder()).
		BorderForeground(styles.ColorBorder).
		Width(h.width-4).
		Padding(0, 1)

	return boxStyle.Render(content)
}

// ViewSimple renders a simpler header without the box (for smaller terminals)
func (h *IntentHeader) ViewSimple() string {
	var lines []string

	// App name and breadcrumbs
	appNameStyle := lipgloss.NewStyle().
		Foreground(styles.ColorAccentTeal).
		Bold(true)

	appNameText := appNameStyle.Render(h.appName)

	var breadcrumbText string
	if len(h.breadcrumbs) > 0 {
		bar := NewBreadcrumbBar(h.width-20, false)
		bar.SetCrumbs(h.breadcrumbs)
		breadcrumbText = bar.View()
	}

	if breadcrumbText != "" {
		separator := lipgloss.NewStyle().
			Foreground(styles.ColorTextMuted).
			Render(" • ")
		lines = append(lines, appNameText+separator+breadcrumbText)
	} else {
		lines = append(lines, appNameText)
	}

	// Description
	if h.description != "" {
		descStyle := lipgloss.NewStyle().
			Foreground(styles.ColorTextSecondary).
			Italic(true)
		lines = append(lines, descStyle.Render(h.description))
	}

	return strings.Join(lines, "\n")
}

// CreateIntentHeaderWithBreadcrumbs is a helper to create a header with breadcrumbs
func CreateIntentHeaderWithBreadcrumbs(width int, description string, trail ...struct {
	Label  string
	Intent string
}) *IntentHeader {
	header := NewIntentHeader(width)
	header.SetDescription(description)

	crumbs := CreateBreadcrumbTrail(trail...)
	header.SetBreadcrumbs(crumbs)

	return header
}
