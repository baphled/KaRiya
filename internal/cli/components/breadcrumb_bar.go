package components

import (
	"strings"

	"github.com/charmbracelet/lipgloss"

	"github.com/baphled/kariya/internal/cli/styles"
)

// Breadcrumb represents a single breadcrumb in the navigation trail
type Breadcrumb struct {
	Label  string
	Icon   string
	Intent string // For future navigation support
}

// BreadcrumbBar renders a navigation breadcrumb trail with icons
type BreadcrumbBar struct {
	crumbs   []Breadcrumb
	width    int
	boxed    bool
	showIcon bool
}

// Icon constants for different intents
const (
	IconHome         = "🏠"
	IconCaptureEvent = "✏️"
	IconTimeline     = "📅"
	IconGenerateCV   = "📄"
	IconExport       = "💾"
	IconConfigure    = "⚙️"
	IconBursts       = "📊"
	IconFacts        = "💡"
	IconImport       = "📥"
	IconMetadata     = "🏷️"
	IconBulkOps      = "📦"
	IconSearch       = "🔍"
	IconFilter       = "🔎"
	IconSort         = "↕️"
	IconEdit         = "✎"
	IconDelete       = "🗑️"
	IconSave         = "💾"
	IconCancel       = "❌"
	IconConfirm      = "✓"
)

// NewBreadcrumbBar creates a new breadcrumb bar
func NewBreadcrumbBar(width int, boxed bool) *BreadcrumbBar {
	return &BreadcrumbBar{
		crumbs:   []Breadcrumb{},
		width:    width,
		boxed:    boxed,
		showIcon: true,
	}
}

// SetCrumbs sets the breadcrumb trail
func (b *BreadcrumbBar) SetCrumbs(crumbs []Breadcrumb) *BreadcrumbBar {
	b.crumbs = crumbs
	return b
}

// AddCrumb adds a breadcrumb to the trail
func (b *BreadcrumbBar) AddCrumb(crumb Breadcrumb) *BreadcrumbBar {
	b.crumbs = append(b.crumbs, crumb)
	return b
}

// SetWidth sets the width for the breadcrumb bar
func (b *BreadcrumbBar) SetWidth(width int) *BreadcrumbBar {
	b.width = width
	return b
}

// SetBoxed sets whether to show a box around the breadcrumbs
func (b *BreadcrumbBar) SetBoxed(boxed bool) *BreadcrumbBar {
	b.boxed = boxed
	return b
}

// ShowIcons controls icon visibility
func (b *BreadcrumbBar) ShowIcons(show bool) *BreadcrumbBar {
	b.showIcon = show
	return b
}

// View renders the breadcrumb bar
func (b *BreadcrumbBar) View() string {
	if len(b.crumbs) == 0 {
		return ""
	}

	// Build breadcrumb string
	separator := "  ▸  "
	var parts []string

	for i, crumb := range b.crumbs {
		var part string

		// Add icon if enabled
		if b.showIcon && crumb.Icon != "" {
			part = crumb.Icon + " " + crumb.Label
		} else {
			part = crumb.Label
		}

		// Style based on position
		if i == len(b.crumbs)-1 {
			// Last breadcrumb - current location (bold and teal)
			styledPart := lipgloss.NewStyle().
				Foreground(styles.ColorAccentTeal).
				Bold(true).
				Render(part)
			parts = append(parts, styledPart)
		} else {
			// Previous breadcrumbs - muted
			styledPart := lipgloss.NewStyle().
				Foreground(styles.ColorTextMuted).
				Render(part)
			parts = append(parts, styledPart)
		}
	}

	breadcrumbStr := strings.Join(parts, separator)

	// Check if it fits in the width
	visualWidth := lipgloss.Width(breadcrumbStr)
	if visualWidth > b.width-4 && len(b.crumbs) > 2 {
		// Too long - show only first and last with ellipsis
		breadcrumbStr = b.renderTruncated()
	}

	// Add box if enabled
	if b.boxed {
		boxStyle := lipgloss.NewStyle().
			Border(lipgloss.RoundedBorder()).
			BorderForeground(styles.ColorBorder).
			Padding(0, 1)

		return boxStyle.Render(breadcrumbStr)
	}

	return breadcrumbStr
}

// renderTruncated renders a truncated version of breadcrumbs
func (b *BreadcrumbBar) renderTruncated() string {
	if len(b.crumbs) < 2 {
		return b.View()
	}

	first := b.crumbs[0]
	last := b.crumbs[len(b.crumbs)-1]

	// Build truncated string
	var firstPart, lastPart string

	// First crumb
	if b.showIcon && first.Icon != "" {
		firstPart = first.Icon + " " + first.Label
	} else {
		firstPart = first.Label
	}
	firstPart = lipgloss.NewStyle().
		Foreground(styles.ColorTextMuted).
		Render(firstPart)

	// Last crumb
	if b.showIcon && last.Icon != "" {
		lastPart = last.Icon + " " + last.Label
	} else {
		lastPart = last.Label
	}
	lastPart = lipgloss.NewStyle().
		Foreground(styles.ColorAccentTeal).
		Bold(true).
		Render(lastPart)

	ellipsis := lipgloss.NewStyle().
		Foreground(styles.ColorTextMuted).
		Render("...")

	return firstPart + "  ▸  " + ellipsis + "  ▸  " + lastPart
}

// GetIconForIntent returns the appropriate icon for a given intent name
func GetIconForIntent(intent string) string {
	iconMap := map[string]string{
		"home":             IconHome,
		"capture_event":    IconCaptureEvent,
		"browse_timeline":  IconTimeline,
		"generate_cv":      IconGenerateCV,
		"export_artifact":  IconExport,
		"configure_system": IconConfigure,
		"burst_management": IconBursts,
		"fact_management":  IconFacts,
		"import_wizard":    IconImport,
		"metadata_editor":  IconMetadata,
		"bulk_operations":  IconBulkOps,
		"search":           IconSearch,
		"filter":           IconFilter,
		"sort":             IconSort,
		"edit":             IconEdit,
		"delete":           IconDelete,
		"save":             IconSave,
		"cancel":           IconCancel,
		"confirm":          IconConfirm,
	}

	if icon, ok := iconMap[intent]; ok {
		return icon
	}
	return ""
}

// CreateBreadcrumbTrail is a helper to create a breadcrumb trail
func CreateBreadcrumbTrail(trail ...struct {
	Label  string
	Intent string
}) []Breadcrumb {
	crumbs := make([]Breadcrumb, len(trail))
	for i, item := range trail {
		crumbs[i] = Breadcrumb{
			Label:  item.Label,
			Icon:   GetIconForIntent(item.Intent),
			Intent: item.Intent,
		}
	}
	return crumbs
}
