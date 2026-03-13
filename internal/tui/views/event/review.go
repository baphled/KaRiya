package event

import (
	"fmt"
	"strings"

	"github.com/baphled/kariya/internal/ui/display"
	"github.com/baphled/kariya/internal/ui/themes"
	"github.com/baphled/kariya/internal/ui/uikit/primitives"
	"github.com/baphled/kariya/internal/ui/uikit/theme"
	"github.com/baphled/kariya/internal/ui/uikit/widgets"
	"github.com/charmbracelet/bubbles/viewport"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
)

// ReviewResult holds the data returned when review is submitted.
type ReviewResult struct {
	Event  display.Event
	Bursts []display.Burst
	Facts  []display.Fact
	Skills []display.SkillSuggestion
}

// Review displays captured event details with inferred bursts and facts.
type Review struct {
	widgets.BaseView
	event          display.Event
	bursts         []display.Burst
	facts          []display.Fact
	skills         []display.SkillSuggestion
	acceptedSkills []display.SkillSuggestion
	acceptedFacts  []display.Fact
	acceptedBursts []display.Burst
	viewport       viewport.Model
}

// maxContentWidth is the maximum width for scrollable content within the viewport.
const maxContentWidth = 80

// NewReview creates a new Review.
//
// Expected:
//   - event must be valid.
//   - burst must be valid.
//   - fact must be valid.
//   - skillsuggestion must be valid.
//
// Returns:
//   - A fully initialized Review ready for use.
//
// Side effects:
//   - None.
func NewReview(evt display.Event, bursts []display.Burst, facts []display.Fact, skills []display.SkillSuggestion) *Review {
	vp := viewport.New(120, 30)
	return &Review{
		event:    evt,
		bursts:   append([]display.Burst(nil), bursts...),
		facts:    append([]display.Fact(nil), facts...),
		skills:   append([]display.SkillSuggestion(nil), skills...),
		viewport: vp,
	}
}

// Init returns nil (no async loading needed).
//
// Returns:
//   - A tea.Cmd value.
//
// Side effects:
//   - None.
func (v *Review) Init() tea.Cmd {
	return nil
}

// Update handles messages and returns a ViewResult on user action.
//
// Expected:
//   - msg is a valid tea.Msg (key press, window resize, or custom message).
//
// Returns:
//   - A tea.Cmd and widgets.ViewResult pair.
//
// Side effects:
//   - Updates internal view state based on message type.
func (v *Review) Update(msg tea.Msg) (tea.Cmd, widgets.ViewResult) {
	if sizeMsg, ok := msg.(tea.WindowSizeMsg); ok {
		const reservedLines = 11
		height := sizeMsg.Height - reservedLines
		if height < 1 {
			height = 1
		}
		v.viewport.Width = sizeMsg.Width
		v.viewport.Height = height
		v.SetTerminalInfo(sizeMsg.Width, sizeMsg.Height)
		return nil, nil
	}

	if keyMsg, ok := msg.(tea.KeyMsg); ok {
		switch keyMsg.Type {
		case tea.KeyEnter:
			return nil, &widgets.SubmitViewResult{
				FormData: ReviewResult{
					Event:  v.event,
					Bursts: v.acceptedBursts,
					Facts:  v.acceptedFacts,
					Skills: v.acceptedSkills,
				},
			}

		case tea.KeyEsc:
			return nil, &widgets.CancelViewResult{}

		case tea.KeyUp:
			v.viewport.ScrollUp(1)
			return nil, nil

		case tea.KeyDown:
			v.viewport.ScrollDown(1)
			return nil, nil

		case tea.KeyPgUp:
			v.viewport.HalfPageUp()
			return nil, nil

		case tea.KeyPgDown:
			v.viewport.HalfPageDown()
			return nil, nil

		case tea.KeyRunes:
			return v.handleRuneKey(keyMsg)
		}
	}

	return nil, nil
}

// RenderContent returns the view content string — no chrome.
//
// Returns:
//   - A string value.
//
// Side effects:
//   - None.
func (v *Review) RenderContent() string {
	content := v.renderContent()
	v.viewport.SetContent(v.centreContent(content))
	return v.viewport.View()
}

// HelpText returns the rendered key binding footer for this view.
//
// Returns:
//   - A string value.
//
// Side effects:
//   - None.
func (v *Review) HelpText() string {
	th := v.resolveThemesTheme()

	badges := []*primitives.Badge{
		primitives.ConfirmBadge(th),
		primitives.HelpKeyBadge("e", "Edit metadata", th),
		primitives.HelpKeyBadge("b", "Edit bursts", th),
		primitives.HelpKeyBadge("f", "Edit facts", th),
	}

	if len(v.skills) > 0 {
		badges = append(badges, primitives.HelpKeyBadge("s", "Edit skills", th))
	}

	badges = append(badges,
		primitives.HelpKeyBadge("↑↓", "Scroll", th),
		primitives.BackBadge(th),
		primitives.QuitBadge(th),
	)

	return primitives.RenderHelpFooter(th, badges...)
}

// SetAcceptedSkills updates the accepted skills list.
//
// Expected:
//   - skillsuggestion must be valid.
//
// Side effects:
//   - None.
func (v *Review) SetAcceptedSkills(skills []display.SkillSuggestion) {
	v.acceptedSkills = skills
}

// SetAcceptedFacts updates the accepted facts list.
//
// Expected:
//   - fact must be valid.
//
// Side effects:
//   - None.
func (v *Review) SetAcceptedFacts(facts []display.Fact) {
	v.acceptedFacts = facts
}

// SetAcceptedBursts updates the accepted bursts list.
//
// Expected:
//   - burst must be valid.
//
// Side effects:
//   - None.
func (v *Review) SetAcceptedBursts(bursts []display.Burst) {
	v.acceptedBursts = bursts
}

// GetSuggestedSkills returns the list of inferred skill suggestions.
//
// Returns:
//   - A []display.SkillSuggestion value.
//
// Side effects:
//   - None.
func (v *Review) GetSuggestedSkills() []display.SkillSuggestion {
	return v.skills
}

// SetSuggestedSkills updates the inferred skill suggestions.
//
// Expected:
//   - skillsuggestion must be valid.
//
// Side effects:
//   - None.
func (v *Review) SetSuggestedSkills(skills []display.SkillSuggestion) {
	v.skills = skills
}

// SetSuggestedBursts updates the inferred burst suggestions.
//
// Expected:
//   - burst must be valid.
//
// Side effects:
//   - None.
func (v *Review) SetSuggestedBursts(bursts []display.Burst) {
	v.bursts = bursts
}

// SetSuggestedFacts updates the inferred fact suggestions.
//
// Expected:
//   - fact must be valid.
//
// Side effects:
//   - None.
func (v *Review) SetSuggestedFacts(facts []display.Fact) {
	v.facts = facts
}

func (v *Review) handleRuneKey(keyMsg tea.KeyMsg) (tea.Cmd, widgets.ViewResult) {
	switch keyMsg.String() {
	case "e":
		return nil, &widgets.NavigateViewResult{ResultData: "edit_metadata"}
	case "b":
		return nil, &widgets.NavigateViewResult{ResultData: "suggest_bursts"}
	case "f":
		return nil, &widgets.NavigateViewResult{ResultData: "suggest_facts"}
	case "s":
		return nil, &widgets.NavigateViewResult{ResultData: "suggest_skills"}
	case "k":
		v.viewport.ScrollUp(1)
		return nil, nil
	case "j":
		v.viewport.ScrollDown(1)
		return nil, nil
	}
	return nil, nil
}

func (v *Review) centreContent(content string) string {
	contentWidth := maxContentWidth
	if v.viewport.Width > 0 && v.viewport.Width < contentWidth {
		contentWidth = v.viewport.Width
	}
	constrained := lipgloss.NewStyle().
		Width(contentWidth).
		MaxWidth(contentWidth).
		Render(content)
	if v.viewport.Width <= 0 {
		return constrained
	}
	return lipgloss.PlaceHorizontal(v.viewport.Width, lipgloss.Center, constrained)
}

func (v *Review) renderContent() string {
	th := v.resolveTheme()

	var parts []string
	parts = append(parts, v.renderEventDetails(th))
	parts = append(parts, v.renderBursts(th))
	parts = append(parts, v.renderFacts(th))
	parts = append(parts, v.renderSkills(th))

	return primitives.JoinVertical(primitives.AlignLeft, parts...)
}

func (v *Review) renderEventDetails(th theme.Theme) string {
	dv := widgets.NewDetailView(th).
		Title("Event Details")

	if isEmptyReviewEvent(v.event) {
		dv.Field("Status", "No event data")
		return dv.Render()
	}

	dv.Field("Text", v.event.Text).
		Field("Date", v.event.Date.Format("2006-01-02")).
		FieldIf("Company", v.event.Company).
		FieldIf("Project", v.event.Project)

	return dv.Render()
}

func (v *Review) renderBursts(th theme.Theme) string {
	var b strings.Builder
	b.WriteString(primitives.Subtitle("Bursts", th).MarginBottom(1).
		MarginTop(1).
		Render())
	b.WriteString("\n")

	if len(v.bursts) == 0 {
		b.WriteString(primitives.Muted("  No bursts detected", th).Render())
		b.WriteString("\n")
		return b.String()
	}

	for i := range v.bursts {
		burst := v.bursts[i]
		accepted := v.isBurstAccepted(burst)
		indicator := v.renderIndicator(accepted, th)
		label := primitives.Body(fmt.Sprintf("  %d. %s", i+1, burst.Name), th).Render()
		b.WriteString(primitives.JoinHorizontal(primitives.AlignLeft, indicator, " ", label))
		b.WriteString("\n")
	}

	return b.String()
}

func (v *Review) renderFacts(th theme.Theme) string {
	var b strings.Builder
	b.WriteString(primitives.Subtitle("Facts", th).
		MarginTop(1).
		Render())
	b.WriteString("\n")

	if len(v.facts) == 0 {
		b.WriteString(primitives.Muted("  No facts detected", th).Render())
		b.WriteString("\n")
		return b.String()
	}

	for i := range v.facts {
		fact := v.facts[i]
		accepted := v.isFactAccepted(fact)
		indicator := v.renderIndicator(accepted, th)
		label := primitives.Body(fmt.Sprintf("  %d. %s", i+1, fact.Text), th).Render()
		b.WriteString(primitives.JoinHorizontal(primitives.AlignLeft, indicator, " ", label))
		b.WriteString("\n")
	}

	return b.String()
}

func (v *Review) renderSkills(th theme.Theme) string {
	var b strings.Builder
	b.WriteString(primitives.Subtitle("Skills", th).
		MarginTop(1).
		Render())
	b.WriteString("\n")

	if len(v.skills) == 0 {
		b.WriteString(primitives.Muted("  No skills detected", th).Render())
		b.WriteString("\n")
		return b.String()
	}

	for i, skill := range v.skills {
		accepted := v.isSkillAccepted(skill)
		indicator := v.renderIndicator(accepted, th)
		confidence := fmt.Sprintf("%.0f%%", skill.Confidence*100)
		label := primitives.Body(fmt.Sprintf("  %d. %s (%s) %s", i+1, skill.Name, skill.Category, confidence), th).Render()
		b.WriteString(primitives.JoinHorizontal(primitives.AlignLeft, indicator, " ", label))
		b.WriteString("\n")
	}

	return b.String()
}

func (v *Review) isBurstAccepted(burst display.Burst) bool {
	for idx := range v.acceptedBursts {
		acceptedBurst := v.acceptedBursts[idx]
		if acceptedBurst.ID == burst.ID || acceptedBurst.Name == burst.Name {
			return true
		}
	}
	return false
}

func (v *Review) isFactAccepted(fact display.Fact) bool {
	for idx := range v.acceptedFacts {
		acceptedFact := v.acceptedFacts[idx]
		if acceptedFact.ID == fact.ID || acceptedFact.Text == fact.Text {
			return true
		}
	}
	return false
}

func (v *Review) isSkillAccepted(skill display.SkillSuggestion) bool {
	for idx := range v.acceptedSkills {
		acceptedSkill := v.acceptedSkills[idx]
		if acceptedSkill.Name == skill.Name {
			return true
		}
	}
	return false
}

func isEmptyReviewEvent(evt display.Event) bool {
	return evt.ID == "" &&
		evt.Text == "" &&
		evt.Date.IsZero() &&
		evt.Company == "" &&
		evt.Project == "" &&
		len(evt.Tags) == 0 &&
		len(evt.Categories) == 0 &&
		len(evt.Skills) == 0 &&
		evt.CreatedAt.IsZero() &&
		evt.UpdatedAt.IsZero()
}

func (v *Review) renderIndicator(accepted bool, th theme.Theme) string {
	if accepted {
		return primitives.SuccessText("●", th).Render()
	}
	return primitives.Muted("○", th).Render()
}

func (v *Review) resolveTheme() theme.Theme {
	if screenTheme := v.GetTheme(); screenTheme != nil {
		if th, ok := screenTheme.(theme.Theme); ok {
			return th
		}
	}
	return theme.Default()
}

func (v *Review) resolveThemesTheme() themes.Theme {
	if screenTheme := v.GetTheme(); screenTheme != nil {
		if th, ok := screenTheme.(themes.Theme); ok {
			return th
		}
	}
	return themes.NewDefaultTheme()
}
