// Package cv provides CV-related screens for the GenerateCV workflow.
package cv

import (
	"fmt"
	"sort"
	"strings"

	"github.com/charmbracelet/bubbles/viewport"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"

	"github.com/baphled/kariya/internal/config"
	"github.com/baphled/kariya/internal/ui/display"
	"github.com/baphled/kariya/internal/ui/themes"
	"github.com/baphled/kariya/internal/ui/uikit/widgets"
)

// Preview displays the full CV content in a scrollable viewport.
// This screen shows the actual CV content with bullet points and allows
// the user to scroll through the entire document.
type Preview struct {
	widgets.BaseView

	cv            display.CVView
	profileConfig *config.ProfileConfig
	viewport      viewport.Model
	ready         bool
	width         int
	height        int
}

// NewPreview creates a Preview view for displaying full CV content.
//
// Expected:
//   - cv must be a valid CVView pointer.
//   - profileConfig may be nil for default profile.
//
// Returns:
//   - A fully initialized Preview ready for use.
//
// Side effects:
//   - None.
func NewPreview(cv display.CVView, profileConfig *config.ProfileConfig) *Preview {
	return &Preview{
		cv:            cv,
		profileConfig: profileConfig,
		width:         80,
		height:        24,
		ready:         false,
	}
}

// Init initializes the screen.
//
// Returns:
//   - A tea.Cmd value.
//
// Side effects:
//   - None.
func (v *Preview) Init() tea.Cmd {
	return nil
}

// Update handles messages.
//
// Expected:
//   - msg must be a valid tea.Msg type.
//
// Returns:
//   - tea.Cmd: command to execute.
//   - widgets.ViewResult: result indicating user action.
//
// Side effects:
//   - May update viewport dimensions.
//   - May return CancelViewResult or NavigateViewResult.
func (v *Preview) Update(msg tea.Msg) (tea.Cmd, widgets.ViewResult) {
	var cmd tea.Cmd

	switch msg := msg.(type) {
	case tea.WindowSizeMsg:
		handleViewportWindowResize(v, msg, &v.width, &v.height, &v.ready)
		return nil, nil

	case tea.KeyMsg:
		switch msg.String() {
		case "esc":
			return nil, &widgets.CancelViewResult{}

		case "enter", "y":
			return nil, &widgets.NavigateViewResult{
				ResultData: Nav{Action: ActionConfirm, CV: v.cv},
			}

		case "e":
			return nil, &widgets.NavigateViewResult{
				ResultData: Nav{Action: ActionEdit},
			}

		case "x":
			return nil, &widgets.NavigateViewResult{
				ResultData: Nav{Action: ActionExport},
			}

		case "g":
			v.viewport.GotoTop()
			return nil, nil

		case "G":
			v.viewport.GotoBottom()
			return nil, nil

		case "up", "k", "down", "j", "pgup", "pgdown", "ctrl+u", "ctrl+d":
			if v.ready {
				v.viewport, cmd = v.viewport.Update(msg)
				return cmd, nil
			}
		}
	}

	return nil, nil
}

// RenderContent renders the screen with full CV content in a scrollable viewport.
//
// Returns:
//   - A string value.
//
// Side effects:
//   - None.
func (v *Preview) RenderContent() string {
	theme := v.getTheme()
	var b strings.Builder

	titleStyle := lipgloss.NewStyle().
		Bold(true).
		Foreground(theme.AccentColor())

	b.WriteString(titleStyle.Render("📄 CV Preview"))
	b.WriteString("\n")
	b.WriteString(strings.Repeat("═", MinInt(60, v.width-4)))
	b.WriteString("\n\n")

	if isEmptyCVView(v.cv) {
		b.WriteString("No CV data available\n")
		b.WriteString("\n")
		b.WriteString(v.HelpText())
		return b.String()
	}

	content := v.renderCVContent()

	if !v.ready {
		viewportHeight := v.height - 7
		if viewportHeight < 5 {
			viewportHeight = 5
		}
		viewportWidth := v.width - 4
		if viewportWidth < 40 {
			viewportWidth = 40
		}

		v.viewport = viewport.New(viewportWidth, viewportHeight)
		v.viewport.SetContent(content)
		v.ready = true
	}

	b.WriteString(v.viewport.View())
	b.WriteString("\n")

	b.WriteString(v.HelpText())

	return b.String()
}

type bulletWithScore struct {
	bullet display.CVBullet
	score  float64
}

// renderCVContent renders the full CV content for the viewport.
func (v *Preview) renderCVContent() string {
	if isEmptyCVView(v.cv) {
		return "No CV data available"
	}

	var b strings.Builder

	contentWidth := v.width - 8
	if contentWidth < 40 {
		contentWidth = 40
	}

	theme := v.getTheme()
	sectionTitleStyle := lipgloss.NewStyle().
		Bold(true).
		Foreground(theme.AccentColor())

	headerStyle := lipgloss.NewStyle().
		Bold(true)

	mutedStyle := lipgloss.NewStyle().
		Foreground(theme.SecondaryColor())

	bulletStyle := lipgloss.NewStyle().
		Foreground(theme.ForegroundColor())

	b.WriteString(v.renderPersonalDetails(contentWidth))
	b.WriteString("\n")

	b.WriteString(v.renderHighlights(contentWidth))
	b.WriteString("\n")

	if len(v.cv.Sections) == 0 {
		b.WriteString("No sections generated yet\n")
		return b.String()
	}

	for idx, section := range v.cv.Sections {
		if idx > 0 {
			b.WriteString("\n")
		}
		b.WriteString(v.renderSection(section, contentWidth, sectionTitleStyle, headerStyle, bulletStyle))
	}

	contentLines := strings.Count(b.String(), "\n")
	if contentLines > v.height-7 {
		b.WriteString("\n")
		b.WriteString(mutedStyle.Render(fmt.Sprintf("(%d lines)", contentLines)))
	}

	return b.String()
}

// renderSection renders a single CV section with its content groups.
func (v *Preview) renderSection(section display.CVSection, contentWidth int, sectionTitleStyle, headerStyle, bulletStyle lipgloss.Style) string {
	var b strings.Builder

	b.WriteString(sectionTitleStyle.Render("## " + section.Title))
	b.WriteString("\n")
	b.WriteString(strings.Repeat("─", len(section.Title)+3))
	b.WriteString("\n")

	if section.SectionType == "summary" && section.Summary != "" {
		wrapped := WordWrap(section.Summary, contentWidth)
		b.WriteString(wrapped)
		b.WriteString("\n")
		return b.String()
	}

	for _, group := range section.Content {
		b.WriteString(v.renderContentGroup(group, contentWidth, headerStyle, bulletStyle))
	}

	return b.String()
}

// renderContentGroup renders a single content group with its header and bullets.
func (v *Preview) renderContentGroup(group display.SectionContentGroup, contentWidth int, headerStyle, bulletStyle lipgloss.Style) string {
	var b strings.Builder
	if headerLine := renderGroupHeader(group, headerStyle); headerLine != "" {
		b.WriteString(headerLine)
		b.WriteString("\n")
	}
	writeGroupBullets(&b, group, contentWidth, bulletStyle)
	return b.String()
}

// renderGroupHeader renders the content group header line.
//
// Expected:
//   - group must be valid.
//   - style must be valid.
//
// Returns:
//   - A string value.
//
// Side effects:
//   - None.
func renderGroupHeader(group display.SectionContentGroup, style lipgloss.Style) string {
	if group.Header == "" {
		return ""
	}
	headerText := buildGroupHeaderText(group)
	return style.Render("  " + headerText)
}

// buildGroupHeaderText builds the header text with optional dates.
//
// Expected:
//   - group must be valid.
//
// Returns:
//   - A string value.
//
// Side effects:
//   - None.
func buildGroupHeaderText(group display.SectionContentGroup) string {
	if group.StartDate == "" || group.EndDate == "" {
		return group.Header
	}
	if group.StartDate == group.EndDate {
		return fmt.Sprintf("%s (%s)", group.Header, group.StartDate)
	}
	return fmt.Sprintf("%s (%s - %s)", group.Header, group.StartDate, group.EndDate)
}

// writeGroupBullets writes bullets to the builder with wrapping.
//
// Expected:
//   - builder must be valid.
//   - group must be valid.
//   - int must be valid.
//   - style must be valid.
//
// Side effects:
//   - Writes to builder.
func writeGroupBullets(builder *strings.Builder, group display.SectionContentGroup, contentWidth int, style lipgloss.Style) {
	bulletWidth := contentWidth - 6
	for bulletIdx := range group.Bullets {
		writeBulletLines(builder, group.Bullets[bulletIdx], bulletWidth, style)
	}
	if len(group.Bullets) > 0 {
		builder.WriteString("\n")
	}
}

// writeBulletLines writes a wrapped bullet with indentation.
//
// Expected:
//   - builder must be valid.
//   - bullet must be valid.
//   - int must be valid.
//   - style must be valid.
//
// Side effects:
//   - Writes to builder.
func writeBulletLines(builder *strings.Builder, bullet display.CVBullet, bulletWidth int, style lipgloss.Style) {
	wrapped := WordWrap(bullet.Text, bulletWidth)
	lines := strings.Split(wrapped, "\n")
	for i, line := range lines {
		builder.WriteString(style.Render(bulletPrefix(i) + line))
		builder.WriteString("\n")
	}
}

// bulletPrefix returns the prefix for a bullet line.
//
// Expected:
//   - int must be valid.
//
// Returns:
//   - A string value.
//
// Side effects:
//   - None.
func bulletPrefix(lineIndex int) string {
	if lineIndex == 0 {
		return "    • "
	}
	return "      "
}

// renderPersonalDetails renders the personal details header.
func (v *Preview) renderPersonalDetails(width int) string {
	profile := NarrativeProfileFromConfig(v.profileConfig)

	var b strings.Builder
	theme := v.getTheme()

	nameStyle := lipgloss.NewStyle().
		Bold(true).
		Foreground(theme.WarningColor())

	roleStyle := lipgloss.NewStyle().
		Foreground(theme.ForegroundColor())

	contactStyle := lipgloss.NewStyle().
		Foreground(theme.SecondaryColor())

	b.WriteString(nameStyle.Render(profile.Name))
	b.WriteString("\n")

	b.WriteString(roleStyle.Render(profile.Role))
	b.WriteString("\n")

	contactLine := fmt.Sprintf("%s  |  %s  |  %s",
		profile.Location,
		profile.Email,
		profile.GitHub,
	)
	b.WriteString(contactStyle.Render(contactLine))
	b.WriteString("\n")

	b.WriteString(strings.Repeat("═", MinInt(width, 60)))
	b.WriteString("\n")

	return b.String()
}

// collectScoredBullets collects all bullets with audience-weighted scores.
func (v *Preview) collectScoredBullets() []bulletWithScore {
	var allBullets []bulletWithScore
	for _, section := range v.cv.Sections {
		for _, group := range section.Content {
			for bulletIdx := range group.Bullets {
				bullet := group.Bullets[bulletIdx]
				score := v.scoreBullet(bullet)
				allBullets = append(allBullets, bulletWithScore{bullet: bullet, score: score})
			}
		}
	}
	return allBullets
}

// scoreBullet calculates the audience-weighted score for a bullet.
func (v *Preview) scoreBullet(bullet display.CVBullet) float64 {
	if v.cv.TargetAudience != "" && !strings.EqualFold(v.cv.TargetAudience, "master") {
		if audienceScore, ok := bullet.AudienceRelevance[strings.ToLower(v.cv.TargetAudience)]; ok {
			return audienceScore * bullet.Confidence
		}
	}
	return bullet.Confidence
}

// renderHighlights renders the Key Highlights section with top-scored bullets.
func (v *Preview) renderHighlights(width int) string {
	if isEmptyCVView(v.cv) || len(v.cv.Sections) == 0 {
		return ""
	}

	allBullets := v.collectScoredBullets()
	if len(allBullets) == 0 {
		return ""
	}

	sort.Slice(allBullets, func(i, j int) bool {
		return allBullets[i].score > allBullets[j].score
	})

	topCount := 5
	if v.profileConfig != nil && v.profileConfig.MaxHighlights > 0 {
		topCount = v.profileConfig.MaxHighlights
	}
	if len(allBullets) < topCount {
		topCount = len(allBullets)
	}
	topBullets := allBullets[:topCount]

	var b strings.Builder
	theme := v.getTheme()

	sectionTitleStyle := lipgloss.NewStyle().
		Bold(true).
		Foreground(theme.AccentColor())

	highlightStyle := lipgloss.NewStyle().
		Foreground(theme.WarningColor())

	b.WriteString(sectionTitleStyle.Render("## Key Highlights"))
	b.WriteString("\n")
	b.WriteString(strings.Repeat("─", 16))
	b.WriteString("\n")

	bulletWidth := width - 6
	for idx := range topBullets {
		bws := topBullets[idx]
		text := bws.bullet.Text
		if text == "" {
			text = bws.bullet.EnhancedText
		}
		wrapped := WordWrap(text, bulletWidth)
		lines := strings.Split(wrapped, "\n")
		for i, line := range lines {
			if i == 0 {
				b.WriteString(highlightStyle.Render("  ★ " + line))
			} else {
				b.WriteString(highlightStyle.Render("    " + line))
			}
			b.WriteString("\n")
		}
	}

	return b.String()
}

// HelpText provides contextual help text for the preview view.
//
// Returns:
//   - A string value.
//
// Side effects:
//   - None.
func (v *Preview) HelpText() string {
	theme := v.getTheme()
	footerStyle := lipgloss.NewStyle().
		Foreground(theme.SecondaryColor())

	var footer strings.Builder
	footer.WriteString(strings.Repeat("─", MinInt(60, v.width-4)))
	footer.WriteString("\n")

	if v.ready && v.viewport.TotalLineCount() > v.viewport.Height {
		pct := int(v.viewport.ScrollPercent() * 100)
		footer.WriteString(footerStyle.Render(
			fmt.Sprintf("↑↓/jk: scroll  g/G: top/bottom  [%d%%]  enter/y: confirm  x: export  esc: back", pct)))
	} else {
		footer.WriteString(footerStyle.Render("↑↓/jk: scroll  enter/y: confirm  x: export  e: edit  esc: back"))
	}

	return footer.String()
}

// GetCV returns the CV data.
//
// Returns:
//   - A display.CVView value.
//
// Side effects:
//   - None.
func (v *Preview) GetCV() display.CVView {
	return v.cv
}

// getTheme returns the theme from BaseView or a default theme.
func (v *Preview) getTheme() themes.Theme {
	if t, ok := v.GetTheme().(themes.Theme); ok && t != nil {
		return t
	}
	return themes.NewDefaultTheme()
}
