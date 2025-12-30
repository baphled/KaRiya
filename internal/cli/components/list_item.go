package components

import (
	"strings"

	"github.com/charmbracelet/lipgloss"

	"github.com/baphled/kariya/internal/cli/styles"
)

// ListItemModel renders a consistent list item with title, subtitle, metadata, and status
type ListItemModel struct {
	title       string
	subtitle    string
	metadata    map[string]string // Additional metadata fields (e.g., "date" -> "2025-12-30")
	isSelected  bool
	isFocused   bool
	maxWidth    int
	showStatus  bool
	statusText  string // Optional status indicator
	statusColor lipgloss.Color
}

// NewListItem creates a new list item with title
func NewListItem(title string, maxWidth int) ListItemModel {
	return ListItemModel{
		title:       title,
		subtitle:    "",
		metadata:    make(map[string]string),
		isSelected:  false,
		isFocused:   false,
		maxWidth:    maxWidth,
		showStatus:  false,
		statusText:  "",
		statusColor: styles.ColorSuccess,
	}
}

// SetSubtitle sets the subtitle
func (li *ListItemModel) SetSubtitle(subtitle string) {
	li.subtitle = subtitle
}

// SetMetadata sets metadata key-value pairs
func (li *ListItemModel) SetMetadata(key, value string) {
	li.metadata[key] = value
}

// GetMetadata returns a metadata value by key
func (li ListItemModel) GetMetadata(key string) string {
	return li.metadata[key]
}

// GetAllMetadata returns all metadata
func (li ListItemModel) GetAllMetadata() map[string]string {
	return li.metadata
}

// SetSelected sets whether the item is selected
func (li *ListItemModel) SetSelected(selected bool) {
	li.isSelected = selected
}

// SetFocused sets whether the item is focused
func (li *ListItemModel) SetFocused(focused bool) {
	li.isFocused = focused
}

// IsSelected returns whether the item is selected
func (li ListItemModel) IsSelected() bool {
	return li.isSelected
}

// IsFocused returns whether the item is focused
func (li ListItemModel) IsFocused() bool {
	return li.isFocused
}

// SetStatus sets the status indicator
func (li *ListItemModel) SetStatus(text string, color lipgloss.Color) {
	li.statusText = text
	li.statusColor = color
	li.showStatus = true
}

// ClearStatus clears the status indicator
func (li *ListItemModel) ClearStatus() {
	li.statusText = ""
	li.showStatus = false
}

// SetMaxWidth sets the maximum width for the item
func (li *ListItemModel) SetMaxWidth(maxWidth int) {
	li.maxWidth = maxWidth
}

// GetTitle returns the title
func (li ListItemModel) GetTitle() string {
	return li.title
}

// GetSubtitle returns the subtitle
func (li ListItemModel) GetSubtitle() string {
	return li.subtitle
}

// View renders the list item
func (li ListItemModel) View() string {
	if li.maxWidth <= 0 {
		return ""
	}

	var parts []string

	// Render title line
	titleLine := li.renderTitleLine()
	parts = append(parts, titleLine)

	// Render subtitle if present
	if li.subtitle != "" {
		subtitleLine := li.renderSubtitle()
		parts = append(parts, subtitleLine)
	}

	// Render metadata if present
	if len(li.metadata) > 0 {
		metadataLine := li.renderMetadata()
		if metadataLine != "" {
			parts = append(parts, metadataLine)
		}
	}

	// Render status if present
	if li.showStatus && li.statusText != "" {
		statusLine := li.renderStatus()
		parts = append(parts, statusLine)
	}

	content := strings.Join(parts, "\n")
	return content
}

// renderTitleLine renders the title with optional selection indicator and status
func (li ListItemModel) renderTitleLine() string {
	var titleStr string

	// Add selection indicator
	indicator := ""
	if li.isSelected {
		indicator = "✓ "
	} else if li.isFocused {
		indicator = "► "
	} else {
		indicator = "  "
	}

	// Style title based on selection/focus
	titleStyle := lipgloss.NewStyle()
	if li.isSelected {
		titleStyle = titleStyle.
			Foreground(styles.ColorAccentTeal).
			Bold(true)
	} else if li.isFocused {
		titleStyle = titleStyle.
			Foreground(styles.ColorAccentTeal)
	} else {
		titleStyle = titleStyle.
			Foreground(styles.ColorTextPrimary)
	}

	// Truncate title if needed (account for indicator)
	title := li.title
	availableWidth := li.maxWidth - 2 // 2 for indicator
	if len(title) > availableWidth {
		title = title[:availableWidth-3] + "..."
	}

	titleStr = indicator + titleStyle.Render(title)

	// Add status indicator on the right if present
	if li.showStatus && li.statusText != "" {
		statusStyle := lipgloss.NewStyle().Foreground(li.statusColor)
		statusIndicator := statusStyle.Render(li.statusText)
		titleStr += " " + statusIndicator
	}

	return titleStr
}

// renderSubtitle renders the subtitle with muted color
func (li ListItemModel) renderSubtitle() string {
	subtitleStyle := lipgloss.NewStyle().
		Foreground(styles.ColorTextMuted).
		Italic(true)

	subtitle := li.subtitle
	// Account for indent
	availableWidth := li.maxWidth - 2
	if len(subtitle) > availableWidth {
		subtitle = subtitle[:availableWidth-3] + "..."
	}

	return "  " + subtitleStyle.Render(subtitle)
}

// renderMetadata renders metadata fields in a compact format
func (li ListItemModel) renderMetadata() string {
	if len(li.metadata) == 0 {
		return ""
	}

	// Format metadata as: "key1: value1  key2: value2"
	var parts []string
	for key, value := range li.metadata {
		parts = append(parts, key+": "+value)
	}

	metadataStr := strings.Join(parts, "  ")
	metadataStyle := lipgloss.NewStyle().
		Foreground(styles.ColorTextMuted).
		Faint(true)

	// Truncate if needed
	availableWidth := li.maxWidth - 2
	if len(metadataStr) > availableWidth {
		metadataStr = metadataStr[:availableWidth-3] + "..."
	}

	return "  " + metadataStyle.Render(metadataStr)
}

// renderStatus renders the status indicator
func (li ListItemModel) renderStatus() string {
	statusStyle := lipgloss.NewStyle().Foreground(li.statusColor)
	return "  " + statusStyle.Render("["+li.statusText+"]")
}

