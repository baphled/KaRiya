package models

import (
	"fmt"
	"strings"
	"time"

	"github.com/baphled/kariya/internal/cli/styles"
	"github.com/baphled/kariya/internal/domain/career"
	"github.com/charmbracelet/lipgloss"
)

// BurstCard displays a single burst with all its details
type BurstCard struct {
	burst *career.Burst
	width int
}

// NewBurstCard creates a new burst card component
func NewBurstCard(burst *career.Burst) *BurstCard {
	return &BurstCard{
		burst: burst,
		width: 80,
	}
}

// SetWidth sets the width for rendering
func (bc *BurstCard) SetWidth(width int) {
	if width < 20 {
		width = 20
	}
	bc.width = width
}

// Render returns the rendered burst card
func (bc *BurstCard) Render() string {
	if bc.burst == nil {
		return ""
	}

	var parts []string

	// Burst name (main content)
	parts = append(parts, bc.renderName())

	// Description if available
	if bc.burst.Description != "" {
		parts = append(parts, bc.renderDescription())
	}

	// Competency focus
	if bc.burst.CompetencyFocus != "" {
		parts = append(parts, bc.renderCompetencyFocus())
	}

	// Event IDs count
	parts = append(parts, bc.renderEventCount())

	// Metadata (dates)
	parts = append(parts, bc.renderMetadata())

	return lipgloss.JoinVertical(
		lipgloss.Left,
		parts...,
	)
}

// renderName renders the burst name
func (bc *BurstCard) renderName() string {
	textStyle := lipgloss.NewStyle().
		Foreground(styles.ColorTextPrimary).
		Bold(true).
		Width(bc.width).
		Padding(0)

	// Wrap text to fit width
	wrapped := wrapText(bc.burst.Name, bc.width-2)
	return textStyle.Render(wrapped)
}

// renderDescription renders the burst description
func (bc *BurstCard) renderDescription() string {
	textStyle := lipgloss.NewStyle().
		Foreground(styles.ColorTextSecondary).
		Width(bc.width).
		MarginTop(1).
		Padding(0)

	// Wrap text to fit width
	wrapped := wrapText(bc.burst.Description, bc.width-2)
	return textStyle.Render(wrapped)
}

// renderCompetencyFocus renders the competency focus
func (bc *BurstCard) renderCompetencyFocus() string {
	label := fmt.Sprintf("🎯 Focus: %s", bc.burst.CompetencyFocus)

	style := lipgloss.NewStyle().
		Foreground(styles.ColorAccentTeal).
		MarginTop(1)

	return style.Render(label)
}

// renderEventCount renders the count of associated events
func (bc *BurstCard) renderEventCount() string {
	count := len(bc.burst.EventIDs)
	eventWord := "event"
	if count != 1 {
		eventWord = "events"
	}

	label := fmt.Sprintf("📊 Contains %d %s", count, eventWord)

	style := lipgloss.NewStyle().
		Foreground(styles.ColorTextSecondary).
		MarginTop(1)

	return style.Render(label)
}

// renderMetadata renders the creation and update dates
func (bc *BurstCard) renderMetadata() string {
	var metadata []string

	if !bc.burst.CreatedAt.IsZero() {
		createdStr := bc.formatDate(bc.burst.CreatedAt)
		metadata = append(metadata, fmt.Sprintf("Created: %s", createdStr))
	}

	if !bc.burst.UpdatedAt.IsZero() && bc.burst.UpdatedAt.After(bc.burst.CreatedAt) {
		updatedStr := bc.formatDate(bc.burst.UpdatedAt)
		metadata = append(metadata, fmt.Sprintf("Updated: %s", updatedStr))
	}

	if len(metadata) == 0 {
		return ""
	}

	style := lipgloss.NewStyle().
		Foreground(styles.ColorTextSecondary).
		MarginTop(1)

	return style.Render(strings.Join(metadata, " | "))
}

// formatDate formats a time.Time to a readable string
func (bc *BurstCard) formatDate(t time.Time) string {
	return t.Format("Jan 2, 2006")
}

