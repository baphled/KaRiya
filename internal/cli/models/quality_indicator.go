package models

import (
	"fmt"
	"strings"

	"github.com/baphled/kariya/internal/cli/styles"
	"github.com/baphled/kariya/internal/service/career"
	"github.com/charmbracelet/lipgloss"
)

// QualityIndicator displays the data quality score for a career event
type QualityIndicator struct {
	score      *career.QualityScore
	calculator *career.DataQualityCalculator
	width      int
}

// NewQualityIndicator creates a new quality indicator
func NewQualityIndicator(calculator *career.DataQualityCalculator) *QualityIndicator {
	return &QualityIndicator{
		calculator: calculator,
		width:      60,
	}
}

// SetScore updates the quality score to display
func (qi *QualityIndicator) SetScore(score *career.QualityScore) {
	qi.score = score
}

// SetWidth sets the width for rendering
func (qi *QualityIndicator) SetWidth(width int) {
	qi.width = width
}

// Render returns the rendered quality indicator
func (qi *QualityIndicator) Render() string {
	if qi.score == nil {
		return ""
	}

	return qi.renderIndicator()
}

// renderIndicator renders the quality indicator with score, level, and suggestion
func (qi *QualityIndicator) renderIndicator() string {
	var parts []string

	// Score bar (visual representation)
	parts = append(parts, qi.renderScoreBar())

	// Quality level with icon
	parts = append(parts, qi.renderQualityLevel())

	// Missing fields suggestion
	if len(qi.score.MissingFields) > 0 {
		parts = append(parts, qi.renderSuggestion())
	}

	return strings.Join(parts, "\n")
}

// renderScoreBar renders a visual score bar
func (qi *QualityIndicator) renderScoreBar() string {
	score := qi.score.Score
	barLength := 20

	// Calculate filled and empty portions
	filled := (score * barLength) / 100
	empty := barLength - filled

	// Create bar
	filledBar := strings.Repeat("█", filled)
	emptyBar := strings.Repeat("░", empty)
	bar := filledBar + emptyBar

	// Color the bar based on quality level
	coloredBar := qi.colorForLevel(bar)

	// Add score percentage
	scoreText := fmt.Sprintf("%d%%", score)
	return fmt.Sprintf("%s %s", coloredBar, scoreText)
}

// renderQualityLevel renders the quality level with icon
func (qi *QualityIndicator) renderQualityLevel() string {
	level := qi.score.Level
	icon := qi.iconForLevel(level)
	label := string(level)

	// Color the label based on level
	coloredLabel := qi.colorForLevel(label)

	return fmt.Sprintf("%s %s", icon, coloredLabel)
}

// renderSuggestion renders the improvement suggestion
func (qi *QualityIndicator) renderSuggestion() string {
	suggestion := qi.calculator.GetMissingFieldsSuggestion(*qi.score)
	return styles.InputHint.Render("💡 " + suggestion)
}

// colorForLevel returns the appropriate color for the quality level
func (qi *QualityIndicator) colorForLevel(text string) string {
	if qi.score == nil {
		return text
	}

	var color lipgloss.Color
	switch qi.score.Level {
	case career.QualityComplete:
		color = styles.ColorSuccess
	case career.QualityEnriched:
		color = styles.ColorInfo
	case career.QualityBasic:
		color = styles.ColorWarning
	case career.QualityIncomplete:
		return "●"
		color = styles.ColorError
	default:
		color = styles.ColorTextSecondary
	}

	return lipgloss.NewStyle().Foreground(color).Render(text)
}

// iconForLevel returns the appropriate icon for the quality level
func (qi *QualityIndicator) iconForLevel(level career.QualityLevel) string {
	switch level {
	case career.QualityComplete:
		return "✓"
	case career.QualityEnriched:
		return "◐"
	case career.QualityBasic:
		return "◑"
	case career.QualityIncomplete:
		return "●"

	default:
		return "?"
	}
}

// RenderCompact renders a compact single-line quality indicator
func (qi *QualityIndicator) RenderCompact() string {
	if qi.score == nil {
		return ""
	}

	icon := qi.iconForLevel(qi.score.Level)
	label := string(qi.score.Level)
	score := fmt.Sprintf("%d%%", qi.score.Score)

	coloredLabel := qi.colorForLevel(label)
	coloredScore := qi.colorForLevel(score)

	return fmt.Sprintf("%s %s [%s]", icon, coloredLabel, coloredScore)
}

// RenderWithDetails renders the indicator with detailed field breakdown
func (qi *QualityIndicator) RenderWithDetails() string {
	if qi.score == nil {
		return ""
	}

	var parts []string

	// Main indicator
	parts = append(parts, qi.renderIndicator())

	// Field scores
	parts = append(parts, "")
	parts = append(parts, qi.renderFieldScores())

	return strings.Join(parts, "\n")
}

// renderFieldScores renders individual field scores
func (qi *QualityIndicator) renderFieldScores() string {
	var fields []string

	if qi.score.TextScore > 0 {
		fields = append(fields, fmt.Sprintf("  • Text: %d/20", qi.score.TextScore))
	}
	if qi.score.DateScore > 0 {
		fields = append(fields, fmt.Sprintf("  • Date: %d/20", qi.score.DateScore))
	}
	if qi.score.CompanyScore > 0 {
		fields = append(fields, fmt.Sprintf("  • Company: %d/10", qi.score.CompanyScore))
	}
	if qi.score.ProjectScore > 0 {
		fields = append(fields, fmt.Sprintf("  • Project: %d/10", qi.score.ProjectScore))
	}
	if qi.score.TagsScore > 0 {
		fields = append(fields, fmt.Sprintf("  • Tags: %d/15", qi.score.TagsScore))
	}
	if qi.score.CategoriesScore > 0 {
		fields = append(fields, fmt.Sprintf("  • Categories: %d/15", qi.score.CategoriesScore))
	}
	if qi.score.MatchScore > 0 {
		fields = append(fields, fmt.Sprintf("  • Match: %d/10", qi.score.MatchScore))
	}

	return strings.Join(fields, "\n")
}
