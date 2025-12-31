package models

import (
	"fmt"
	"strings"


	"github.com/baphled/kariya/internal/cli/styles"
	"github.com/baphled/kariya/internal/domain/career"
	"github.com/charmbracelet/lipgloss"
)

// FactCard displays a single fact with all its details
type FactCard struct {
	fact  *career.Fact
	width int
}

// NewFactCard creates a new fact card component
func NewFactCard(fact *career.Fact) *FactCard {
	return &FactCard{
		fact:  fact,
		width: 80,
	}
}

// SetWidth sets the width for rendering
func (fc *FactCard) SetWidth(width int) {
	if width < 20 {
		width = 20
	}
	fc.width = width
}

// Render returns the rendered fact card
func (fc *FactCard) Render() string {
	if fc.fact == nil {
		return ""
	}

	var parts []string

	// Fact text (main content)
	parts = append(parts, fc.renderText())

	// Competencies (badges)
	parts = append(parts, fc.renderCompetencies())

	// Role fit with icon
	parts = append(parts, fc.renderRoleFit())

	// Audience relevance
	parts = append(parts, fc.renderAudience())

	// Strength signal
	if fc.fact.StrengthSignal != "" {
		parts = append(parts, fc.renderStrengthSignal())
	}

	// Source reference
	parts = append(parts, fc.renderSource())

	// Metadata (dates)
	parts = append(parts, fc.renderMetadata())

	return lipgloss.JoinVertical(
		lipgloss.Left,
		parts...,
	)
}

// renderText renders the fact text
func (fc *FactCard) renderText() string {
	textStyle := lipgloss.NewStyle().
		Foreground(styles.ColorTextPrimary).
		Width(fc.width).
		Padding(0)

	// Wrap text to fit width
	wrapped := wrapText(fc.fact.Text, fc.width-2)
	return textStyle.Render(wrapped)
}

// renderCompetencies renders competency badges
func (fc *FactCard) renderCompetencies() string {
	if len(fc.fact.CompetencyCategories) == 0 {
		return ""
	}

	var badges []string
	for _, comp := range fc.fact.CompetencyCategories {
		badge := styles.TagBase.Render(comp)
		badges = append(badges, badge)
	}

	badgesLine := lipgloss.JoinHorizontal(lipgloss.Left, badges...)
	return lipgloss.NewStyle().
		MarginTop(1).
		Render(badgesLine)
}

// renderRoleFit renders the role fit with icon
func (fc *FactCard) renderRoleFit() string {
	icon := getRoleFitIcon(fc.fact.RoleFit)
	label := fmt.Sprintf("%s  %s", icon, fc.fact.RoleFit)

	style := lipgloss.NewStyle().
		Foreground(styles.ColorAccentTeal).
		MarginTop(1)

	return style.Render(label)
}

// renderAudience renders the audience relevance
func (fc *FactCard) renderAudience() string {
	if len(fc.fact.AudienceRelevance) == 0 {
		return ""
	}

	audiences := strings.Join(fc.fact.AudienceRelevance, ", ")
	label := fmt.Sprintf("👥 Audience: %s", audiences)

	style := lipgloss.NewStyle().
		Foreground(styles.ColorTextSecondary).
		MarginTop(1)

	return style.Render(label)
}

// renderStrengthSignal renders the strength signal
func (fc *FactCard) renderStrengthSignal() string {
	if fc.fact.StrengthSignal == "" {
		return ""
	}

	label := fmt.Sprintf("💪 Strength: %s", fc.fact.StrengthSignal)

	style := lipgloss.NewStyle().
		Foreground(styles.ColorAccentGreen).
		MarginTop(1)

	return style.Render(label)
}

// renderSource renders the source reference (event or burst ID)
func (fc *FactCard) renderSource() string {
	var source string
	if fc.fact.SourceEventID != "" {
		source = fmt.Sprintf("📋 From Event: %s", fc.fact.SourceEventID)
	} else if fc.fact.SourceBurstID != "" {
		source = fmt.Sprintf("🎯 From Burst: %s", fc.fact.SourceBurstID)
	} else {
		source = "📌 Source: Unknown"
	}

	style := lipgloss.NewStyle().
		Foreground(styles.ColorTextMuted).
		MarginTop(1)

	return style.Render(source)
}

// renderMetadata renders creation date
func (fc *FactCard) renderMetadata() string {
	created := fc.fact.CreatedAt.Format("2006-01-02")
	label := fmt.Sprintf("Created: %s", created)

	style := lipgloss.NewStyle().
		Foreground(styles.ColorTextMuted).
		Italic(true).
		MarginTop(1)

	return style.Render(label)
}

// getRoleFitIcon returns an icon for the role fit
func getRoleFitIcon(roleFit career.RoleFit) string {
	switch roleFit {
	case career.RoleFitPrincipal:
		return "👑"
	case career.RoleFitEM:
		return "👔"
	case career.RoleFitStaff:
		return "⭐"
	case career.RoleFitSeniorIC:
		return "🎖️"
	default:
		return "✨"
	}
}

// wrapText wraps text to a specified width
func wrapText(text string, width int) string {
	if width < 10 {
		width = 10
	}

	words := strings.Fields(text)
	var lines []string
	var currentLine []string

	for _, word := range words {
		testLine := append(currentLine, word)
		testText := strings.Join(testLine, " ")
		if len(testText) > width {
			if len(currentLine) > 0 {
				lines = append(lines, strings.Join(currentLine, " "))
				currentLine = []string{word}
			} else {
				lines = append(lines, word)
				currentLine = []string{}
			}
		} else {
			currentLine = testLine
		}
	}

	if len(currentLine) > 0 {
		lines = append(lines, strings.Join(currentLine, " "))
	}

	return strings.Join(lines, "\n")
}
