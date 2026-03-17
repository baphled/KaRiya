// Package cv provides presentational view components for the CV generation domain.
package cv

import (
	"fmt"
	"sort"
	"strconv"
	"strings"

	"github.com/charmbracelet/bubbles/viewport"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"

	"github.com/baphled/kariya/internal/config"
	"github.com/baphled/kariya/internal/ui/display"
	"github.com/baphled/kariya/internal/ui/themes"
	"github.com/baphled/kariya/internal/ui/uikit/primitives"
	"github.com/baphled/kariya/internal/ui/uikit/widgets"
)

const indentedPairFmt = "  %s %s\n"

// Review displays a summary review of the generated CV.
// This view shows metadata, statistics, and section overview before
// allowing the user to view the full CV preview with scrolling.
//
// Expected:
//   - None.
//
// Returns:
//   - A Review view instance.
//
// Side effects:
//   - None.
type Review struct {
	widgets.BaseView

	cv            display.CVView
	profileConfig *config.ProfileConfig
	summary       *GenerationSummary
	viewport      viewport.Model
	ready         bool
	width         int
	height        int
}

// NewReview creates a Review view for displaying CV metadata and statistics.
//
// Expected:
//   - cv must be a valid CVView pointer.
//   - profileConfig may be nil for default profile.
//   - summary may be nil for simple review mode.
//
// Returns:
//   - A fully initialized Review ready for use.
//
// Side effects:
//   - None.
func NewReview(cv display.CVView, profileConfig *config.ProfileConfig, summary *GenerationSummary) *Review {
	return &Review{
		cv:            cv,
		profileConfig: profileConfig,
		summary:       summary,
		width:         80,
		height:        24,
		ready:         false,
	}
}

// Init initializes the view.
//
// Returns:
//   - A tea.Cmd value.
//
// Side effects:
//   - None.
func (v *Review) Init() tea.Cmd {
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
func (v *Review) Update(msg tea.Msg) (tea.Cmd, widgets.ViewResult) {
	var cmd tea.Cmd

	switch msg := msg.(type) {
	case tea.WindowSizeMsg:
		handleViewportWindowResize(v, msg, &v.width, &v.height, &v.ready)
		return nil, nil

	case tea.KeyMsg:
		switch msg.String() {
		case "esc":
			return nil, &widgets.CancelViewResult{}

		case "enter", "p":
			return nil, &widgets.NavigateViewResult{
				ResultData: Nav{Action: ActionPreview},
			}

		case "x":
			return nil, &widgets.NavigateViewResult{
				ResultData: Nav{Action: ActionExport},
			}

		case "e":
			return nil, &widgets.NavigateViewResult{
				ResultData: Nav{Action: ActionEdit},
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

// RenderContent renders the review view with CV metadata and section summary.
//
// Expected:
//   - None.
//
// Returns:
//   - A string value.
//
// Side effects:
//   - None.
func (v *Review) RenderContent() string {
	theme := v.getTheme()
	var b strings.Builder

	titleStyle := lipgloss.NewStyle().
		Bold(true).
		Foreground(theme.AccentColor())

	title := "📋 CV Review"
	if v.summary != nil {
		title = "📋 CV Configuration Review"
	}
	b.WriteString(titleStyle.Render(title))
	b.WriteString("\n")
	b.WriteString(strings.Repeat("═", MinInt(60, v.width-4)))
	b.WriteString("\n\n")

	if isEmptyCVView(v.cv) {
		b.WriteString("No CV data available\n")
		b.WriteString("\n")
		b.WriteString(strings.Repeat("─", MinInt(60, v.width-4)))
		b.WriteString("\n")
		b.WriteString("esc: back")
		return b.String()
	}

	content := v.renderContent()

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

	return b.String()
}

// renderContent renders the scrollable content.
func (v *Review) renderContent() string {
	if v.summary != nil {
		return v.renderSummaryContent()
	}
	var b strings.Builder

	b.WriteString(v.renderPersonalDetails())
	b.WriteString(v.renderCVDetails())
	b.WriteString(v.renderStatistics())
	b.WriteString(v.renderHighlights())
	b.WriteString(v.renderSectionsList())

	return b.String()
}

func (v *Review) renderSummaryContent() string {
	var b strings.Builder

	b.WriteString(v.renderPersonalDetails())
	b.WriteString(v.renderCVDetails())
	b.WriteString(v.renderSummarySection())
	b.WriteString(v.renderGenerationSettings())
	b.WriteString(v.renderStatistics())

	return b.String()
}

func (v *Review) renderSummarySection() string {
	summaryText := v.getSummaryText()
	if summaryText == "" {
		return ""
	}

	theme := v.getTheme()
	sectionTitleStyle := lipgloss.NewStyle().Bold(true).Foreground(theme.AccentColor())

	var b strings.Builder
	b.WriteString(sectionTitleStyle.Render("📝 Summary"))
	b.WriteString("\n")
	b.WriteString(strings.Repeat("─", MinInt(60, v.width-4)))
	b.WriteString("\n")

	for _, line := range strings.Split(summaryText, "\n") {
		line = strings.TrimSpace(line)
		if line == "" {
			b.WriteString("\n")
			continue
		}
		wrapped := WordWrap(line, 56)
		for _, wl := range strings.Split(wrapped, "\n") {
			b.WriteString("  " + wl + "\n")
		}
	}
	b.WriteString("\n")

	return b.String()
}

func (v *Review) renderGenerationSettings() string {
	theme := v.getTheme()
	sectionTitleStyle := lipgloss.NewStyle().Bold(true).Foreground(theme.AccentColor())
	labelStyle := lipgloss.NewStyle().Foreground(theme.SecondaryColor())
	valueStyle := lipgloss.NewStyle().Foreground(theme.ForegroundColor())

	var b strings.Builder
	b.WriteString(sectionTitleStyle.Render("🎯 Generation Settings"))
	b.WriteString("\n")
	b.WriteString(strings.Repeat("─", MinInt(60, v.width-4)))
	b.WriteString("\n")

	if v.summary.SelectedProfile != nil && v.summary.SelectedProfile.Name != "" {
		profileName := v.summary.SelectedProfile.Name
		fmt.Fprintf(&b, indentedPairFmt, labelStyle.Render("Profile:"), valueStyle.Render(profileName))
	}
	if v.summary.SelectedAudience != "" {
		fmt.Fprintf(&b, indentedPairFmt, labelStyle.Render("Audience:"), valueStyle.Render(v.summary.SelectedAudience))
	}
	if v.summary.TechnologyFocus != "" {
		fmt.Fprintf(&b, indentedPairFmt, labelStyle.Render("Tech Focus:"), valueStyle.Render(v.summary.TechnologyFocus))
	}
	if len(v.summary.Technologies) > 0 {
		techText := strings.Join(v.summary.Technologies, ", ")
		fmt.Fprintf(&b, indentedPairFmt, labelStyle.Render("Technologies:"), valueStyle.Render(techText))
	}
	if v.summary.FocusArea != "" {
		fmt.Fprintf(&b, indentedPairFmt, labelStyle.Render("Focus Area:"), valueStyle.Render(v.summary.FocusArea))
	}
	if v.summary.SkillsFormat != "" {
		skillsText := v.summary.SkillsFormat
		if v.summary.SkillsLimit > 0 {
			skillsText = fmt.Sprintf("%s (%d max)", skillsText, v.summary.SkillsLimit)
		}
		fmt.Fprintf(&b, indentedPairFmt, labelStyle.Render("Skills Format:"), valueStyle.Render(skillsText))
	}
	if v.summary.CVLength != "" {
		fmt.Fprintf(&b, indentedPairFmt, labelStyle.Render("CV Length:"), valueStyle.Render(v.summary.CVLength))
	}
	b.WriteString("\n")

	return b.String()
}

func (v *Review) renderPersonalDetails() string {
	theme := v.getTheme()
	sectionTitleStyle := lipgloss.NewStyle().Bold(true).Foreground(theme.AccentColor())
	labelStyle := lipgloss.NewStyle().Foreground(theme.SecondaryColor())
	valueStyle := lipgloss.NewStyle().Foreground(theme.ForegroundColor())

	profile := NarrativeProfileFromConfig(v.profileConfig)

	var b strings.Builder
	b.WriteString(sectionTitleStyle.Render("👤 Personal Details"))
	b.WriteString("\n")
	b.WriteString(strings.Repeat("─", MinInt(60, v.width-4)))
	b.WriteString("\n")
	fmt.Fprintf(&b, indentedPairFmt, labelStyle.Render("Name:"), valueStyle.Render(profile.Name))
	fmt.Fprintf(&b, indentedPairFmt, labelStyle.Render("Email:"), valueStyle.Render(profile.Email))
	fmt.Fprintf(&b, indentedPairFmt, labelStyle.Render("Location:"), valueStyle.Render(profile.Location))
	b.WriteString("\n")

	return b.String()
}

func (v *Review) renderCVDetails() string {
	theme := v.getTheme()
	sectionTitleStyle := lipgloss.NewStyle().Bold(true).Foreground(theme.AccentColor())
	labelStyle := lipgloss.NewStyle().Foreground(theme.SecondaryColor())
	valueStyle := lipgloss.NewStyle().Foreground(theme.ForegroundColor())

	var b strings.Builder
	b.WriteString(sectionTitleStyle.Render("📄 CV Details"))
	b.WriteString("\n")
	b.WriteString(strings.Repeat("─", MinInt(60, v.width-4)))
	b.WriteString("\n")
	fmt.Fprintf(&b, indentedPairFmt, labelStyle.Render("CV Name:"), valueStyle.Render(v.cv.Name))
	fmt.Fprintf(&b, indentedPairFmt, labelStyle.Render("Role:"), valueStyle.Render(v.cv.TargetRole))
	fmt.Fprintf(&b, indentedPairFmt, labelStyle.Render("Audience:"), valueStyle.Render(v.cv.TargetAudience))
	b.WriteString("\n")

	return b.String()
}

func (v *Review) renderStatistics() string {
	theme := v.getTheme()
	sectionTitleStyle := lipgloss.NewStyle().Bold(true).Foreground(theme.AccentColor())
	labelStyle := lipgloss.NewStyle().Foreground(theme.SecondaryColor())
	valueStyle := lipgloss.NewStyle().Foreground(theme.ForegroundColor())

	var b strings.Builder
	b.WriteString(sectionTitleStyle.Render("📊 Statistics"))
	b.WriteString("\n")
	b.WriteString(strings.Repeat("─", MinInt(60, v.width-4)))
	b.WriteString("\n")
	fmt.Fprintf(&b, indentedPairFmt, labelStyle.Render("Source Events:"), valueStyle.Render(strconv.Itoa(v.cv.SourceEventCount)))
	fmt.Fprintf(&b, indentedPairFmt, labelStyle.Render("Source Facts:"), valueStyle.Render(strconv.Itoa(v.cv.SourceFactCount)))
	fmt.Fprintf(&b, indentedPairFmt, labelStyle.Render("Sections:"), valueStyle.Render(strconv.Itoa(len(v.cv.Sections))))

	totalBullets := v.countTotalBullets()
	fmt.Fprintf(&b, indentedPairFmt, labelStyle.Render("Total Bullets:"), valueStyle.Render(strconv.Itoa(totalBullets)))
	b.WriteString("\n")

	return b.String()
}

func (v *Review) countTotalBullets() int {
	total := 0
	for _, section := range v.cv.Sections {
		for _, group := range section.Content {
			total += len(group.Bullets)
		}
	}
	return total
}

func (v *Review) scoreBullet(bullet display.CVBullet) float64 {
	if v.cv.TargetAudience != "" && !strings.EqualFold(v.cv.TargetAudience, "master") {
		audienceRelevance := 0.0
		if bullet.AudienceRelevance != nil {
			if relevance, ok := bullet.AudienceRelevance[v.cv.TargetAudience]; ok {
				audienceRelevance = relevance
			}
		}
		return audienceRelevance * bullet.Confidence
	}
	return bullet.Confidence
}

func (v *Review) renderHighlights() string {
	theme := v.getTheme()
	sectionTitleStyle := lipgloss.NewStyle().Bold(true).Foreground(theme.AccentColor())
	valueStyle := lipgloss.NewStyle().Foreground(theme.ForegroundColor())

	type scoredBullet struct {
		bullet display.CVBullet
		score  float64
	}

	var allBullets []scoredBullet
	for _, section := range v.cv.Sections {
		for _, group := range section.Content {
			for bulletIdx := range group.Bullets {
				bullet := group.Bullets[bulletIdx]
				score := v.scoreBullet(bullet)
				allBullets = append(allBullets, scoredBullet{bullet: bullet, score: score})
			}
		}
	}

	sort.Slice(allBullets, func(i, j int) bool {
		return allBullets[i].score > allBullets[j].score
	})

	topCount := 10
	if len(allBullets) < topCount {
		topCount = len(allBullets)
	}

	if topCount == 0 {
		return ""
	}

	var b strings.Builder
	b.WriteString(sectionTitleStyle.Render("✨ Top Highlights"))
	b.WriteString("\n")
	b.WriteString(strings.Repeat("─", 60))
	b.WriteString("\n")

	for i := range topCount {
		bullet := allBullets[i].bullet
		text := bullet.Text
		if text == "" {
			text = bullet.EnhancedText
		}
		fmt.Fprintf(&b, "  • %s\n", valueStyle.Render(text))
	}
	b.WriteString("\n")

	return b.String()
}

func (v *Review) renderSectionsList() string {
	theme := v.getTheme()
	sectionTitleStyle := lipgloss.NewStyle().Bold(true).Foreground(theme.AccentColor())
	labelStyle := lipgloss.NewStyle().Foreground(theme.SecondaryColor())
	valueStyle := lipgloss.NewStyle().Foreground(theme.ForegroundColor())

	if len(v.cv.Sections) == 0 {
		return "\n  0 sections generated\n"
	}

	var b strings.Builder
	b.WriteString(sectionTitleStyle.Render("📑 Sections"))
	b.WriteString("\n")
	b.WriteString(strings.Repeat("─", MinInt(60, v.width-4)))
	b.WriteString("\n")

	for _, section := range v.cv.Sections {
		sectionBullets := v.countSectionBullets(section)
		bulletText := "bullets"
		if sectionBullets == 1 {
			bulletText = "bullet"
		}

		typeIndicator := "  •"
		if section.SectionType == "summary" {
			typeIndicator = "  ✎"
		}

		fmt.Fprintf(&b, "%s %s %s\n",
			typeIndicator,
			valueStyle.Render(section.Title),
			labelStyle.Render(fmt.Sprintf("(%d %s)", sectionBullets, bulletText)))
	}

	return b.String()
}

func (v *Review) countSectionBullets(section display.CVSection) int {
	count := 0
	for _, group := range section.Content {
		count += len(group.Bullets)
	}
	return count
}

// HelpText renders the help footer with scroll indicator.
//
// Expected:
//   - None.
//
// Returns:
//   - A string value.
//
// Side effects:
//   - None.
func (v *Review) HelpText() string {
	th := v.getTheme()
	badges := []*primitives.Badge{
		primitives.HelpKeyBadge("↑↓/jk", "Scroll", th),
		primitives.HelpKeyBadge("Enter/p", "Preview", th),
		primitives.HelpKeyBadge("x", "Export", th),
		primitives.HelpKeyBadge("e", "Edit", th),
		primitives.HelpKeyBadge("Esc", "Back", th),
	}
	if v.ready {
		pct := int(v.viewport.ScrollPercent() * 100)
		badges = append([]*primitives.Badge{primitives.HelpKeyBadge(fmt.Sprintf("%d%%", pct), "", th)}, badges...)
	}
	return primitives.RenderHelpFooter(th, badges...)
}

// GetCV returns the CV data.
//
// Returns:
//   - A display.CVView value.
//
// Side effects:
//   - None.
func (v *Review) GetCV() display.CVView {
	return v.cv
}

// getTheme returns the theme from BaseView or a default theme.
func (v *Review) getTheme() themes.Theme {
	if t, ok := v.GetTheme().(themes.Theme); ok && t != nil {
		return t
	}
	return themes.NewDefaultTheme()
}

// getSummaryText extracts the summary text from the CV sections.
func (v *Review) getSummaryText() string {
	if isEmptyCVView(v.cv) {
		return ""
	}
	for _, section := range v.cv.Sections {
		if section.SectionType == "summary" && section.Summary != "" {
			return section.Summary
		}
	}
	return ""
}

func isEmptyCVView(cv display.CVView) bool {
	return cv.ID == "" &&
		cv.Name == "" &&
		cv.TargetRole == "" &&
		cv.TargetAudience == "" &&
		len(cv.Sections) == 0 &&
		cv.GeneratedAt.IsZero() &&
		cv.SourceEventCount == 0 &&
		cv.SourceFactCount == 0
}
