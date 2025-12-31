package components

import (
	"fmt"
	"strings"

	"github.com/baphled/kariya/internal/cli/styles"
	"github.com/charmbracelet/lipgloss"
)

// ProgressIndicator displays progress for multi-step workflows
type ProgressIndicator struct {
	current          int
	total            int
	label            string
	status           string // "success", "warning", "error", "info", ""
	showProgressBar  bool
	showPercentage   bool
	width            int
	barWidth         int
}

// NewProgressIndicator creates a new progress indicator
func NewProgressIndicator(current, total int) ProgressIndicator {
	return ProgressIndicator{
		current:         current,
		total:           total,
		label:           "",
		status:          "",
		showProgressBar: false,
		showPercentage:  false,
		width:           80,
		barWidth:        40,
	}
}

// SetCurrent updates the current value
func (p *ProgressIndicator) SetCurrent(current int) {
	p.current = current
}

// SetTotal updates the total value
func (p *ProgressIndicator) SetTotal(total int) {
	p.total = total
}

// SetLabel sets the label text
func (p *ProgressIndicator) SetLabel(label string) {
	p.label = label
}

// SetStatus sets the status type for coloring
func (p *ProgressIndicator) SetStatus(status string) {
	p.status = status
}

// EnableProgressBar enables or disables the progress bar visualization
func (p *ProgressIndicator) EnableProgressBar(enabled bool) {
	p.showProgressBar = enabled
}

// ShowPercentage enables or disables percentage display
func (p *ProgressIndicator) ShowPercentage(show bool) {
	p.showPercentage = show
}

// SetWidth sets the width for rendering
func (p *ProgressIndicator) SetWidth(width int) {
	p.width = width
	p.barWidth = width / 2
}

// GetCurrent returns the current value
func (p ProgressIndicator) GetCurrent() int {
	return p.current
}

// GetTotal returns the total value
func (p ProgressIndicator) GetTotal() int {
	return p.total
}

// GetLabel returns the label text
func (p ProgressIndicator) GetLabel() string {
	return p.label
}

// GetStatus returns the status type
func (p ProgressIndicator) GetStatus() string {
	return p.status
}

// GetPercentage calculates the completion percentage
func (p ProgressIndicator) GetPercentage() float64 {
	if p.total == 0 {
		return 0
	}
	return float64(p.current) / float64(p.total) * 100
}

// IsComplete returns true if progress is at 100%
func (p ProgressIndicator) IsComplete() bool {
	return p.current >= p.total && p.total > 0
}

// View renders the progress indicator
func (p ProgressIndicator) View() string {
	var parts []string

	// Build step indicator text
	stepText := fmt.Sprintf("%d/%d", p.current, p.total)
	
	// Add label if present
	if p.label != "" {
		stepText = fmt.Sprintf("%s: %s", p.label, stepText)
	}

	// Add percentage if enabled
	if p.showPercentage {
		stepText = fmt.Sprintf("%s (%.0f%%)", stepText, p.GetPercentage())
	}

	// Apply status color
	style := p.getStatusStyle()
	parts = append(parts, style.Render(stepText))

	// Add progress bar if enabled
	if p.showProgressBar {
		bar := p.renderProgressBar()
		parts = append(parts, bar)
	}

	return strings.Join(parts, "\n")
}

// getStatusStyle returns the appropriate style based on status
func (p ProgressIndicator) getStatusStyle() lipgloss.Style {
	baseStyle := lipgloss.NewStyle().Bold(true)

	switch p.status {
	case "success":
		return baseStyle.Foreground(styles.ColorSuccess)
	case "error":
		return baseStyle.Foreground(styles.ColorError)
	case "warning":
		return baseStyle.Foreground(styles.ColorWarning)
	case "info":
		return baseStyle.Foreground(styles.ColorInfo)
	default:
		return baseStyle.Foreground(styles.ColorTextPrimary)
	}
}

// renderProgressBar renders a visual progress bar
func (p ProgressIndicator) renderProgressBar() string {
	if p.total == 0 {
		return ""
	}

	percentage := p.GetPercentage()
	filledWidth := int(float64(p.barWidth) * percentage / 100)
	emptyWidth := p.barWidth - filledWidth

	filled := strings.Repeat("█", filledWidth)
	empty := strings.Repeat("░", emptyWidth)

	filledStyle := lipgloss.NewStyle().Foreground(p.getBarColor())
	emptyStyle := lipgloss.NewStyle().Foreground(styles.ColorTextMuted)

	return filledStyle.Render(filled) + emptyStyle.Render(empty)
}

// getBarColor returns the bar color based on status
func (p ProgressIndicator) getBarColor() lipgloss.Color {
	switch p.status {
	case "success":
		return styles.ColorSuccess
	case "error":
		return styles.ColorError
	case "warning":
		return styles.ColorWarning
	case "info":
		return styles.ColorInfo
	default:
		return styles.ColorAccentTeal
	}
}

