package capture

import (
	"fmt"
	"strings"

	"github.com/baphled/kariya/internal/cli/screens"
	"github.com/baphled/kariya/internal/cli/screens/base"
	"github.com/baphled/kariya/internal/cli/themes"
	"github.com/baphled/kariya/internal/cli/uikit/primitives"
	"github.com/baphled/kariya/internal/cli/uikit/theme"
	"github.com/baphled/kariya/internal/cli/uikit/widgets"
	"github.com/baphled/kariya/internal/domain/career"
	tea "github.com/charmbracelet/bubbletea"
)

// EventReviewScreen displays captured event details with inferred bursts and facts.
//
// This screen presents:
// - Event details (text, date, company, project, etc.)
// - Inferred bursts (if any)
// - Inferred facts (if any)
//
// Users can:
// - Confirm and proceed to submission (Enter)
// - Navigate to edit screens (e=metadata, b=bursts, f=facts)
// - Cancel and return to form (Esc)
//
// Keyboard Shortcuts:
// - Enter: Confirm and proceed to submission
// - e: Edit event metadata
// - b: Edit bursts
// - f: Edit facts
// - Esc: Cancel and return to form
//
// Related:
// - internal/cli/intents/captureevent/types.go (ReviewInferredEventState)
// - tasks/tasks-42-tui-architecture-refactor.md (Phase 1: CaptureEvent Migration).
type EventReviewScreen struct {
	*base.Screen

	// event being reviewed.
	event *career.Event

	// bursts inferred from event.
	bursts []*career.Burst

	// facts inferred from event.
	facts []*career.Fact

	// breadcrumbs for the view header.
	breadcrumbs []string
}

// NewEventReviewScreen creates a new EventReviewScreen.
//
// Parameters:
//   - breadcrumbs: Breadcrumb trail for header (e.g., ["Main Menu", "Capture Event", "Review"])
//   - event: Captured event to review
//   - bursts: Inferred bursts (may be nil or empty)
//   - facts: Inferred facts (may be nil or empty)
//
// Expected:
//   - breadcrumbs must be a valid slice of strings.
//   - event must be a valid *career.Event.
//   - bursts must be a valid slice of *career.Burst (can be nil or empty).
//   - facts must be a valid slice of *career.Fact (can be nil or empty).
//
// Returns:
//   - A fully initialized EventReviewScreen ready for use.
//
// Side effects:
//   - None.
func NewEventReviewScreen(
	breadcrumbs []string,
	event *career.Event,
	bursts []*career.Burst,
	facts []*career.Fact,
) *EventReviewScreen {
	return &EventReviewScreen{
		Screen:      base.NewBaseScreen(),
		event:       event,
		bursts:      bursts,
		facts:       facts,
		breadcrumbs: breadcrumbs,
	}
}

// Update implements the Screen interface.
//
// Handles:
// - Enter -> returns SubmitResult with event, bursts, facts
// - e/b/f -> returns NavigateResult with edit action
// - Esc -> returns CancelResult
// - WindowSizeMsg -> updates dimensions.
//
// Expected:
//   - msg must be a valid tea.Msg type.
//
// Returns:
//   - A tea.Cmd value.
//   - A screens.ScreenResult value.
//
// Side effects:
//   - May return SubmitResult on Enter.
//   - May return NavigateResult on e/b/f keys.
//   - May return CancelResult on Escape.
func (s *EventReviewScreen) Update(msg tea.Msg) (tea.Cmd, screens.ScreenResult) {
	if cmd := s.HandleWindowSizeMsg(msg); cmd != nil {
		return cmd, nil
	}

	if keyMsg, ok := msg.(tea.KeyMsg); ok {
		switch keyMsg.Type {
		case tea.KeyEnter:
			return nil, &screens.SubmitResult{
				FormData: map[string]interface{}{
					"event":  s.event,
					"bursts": s.bursts,
					"facts":  s.facts,
				},
			}

		case tea.KeyEsc:
			return nil, &screens.CancelResult{}

		case tea.KeyRunes:
			return s.handleRuneKey(keyMsg)
		}
	}

	return nil, nil
}

// handleRuneKey dispatches single-character key presses to edit actions.
func (s *EventReviewScreen) handleRuneKey(keyMsg tea.KeyMsg) (tea.Cmd, screens.ScreenResult) {
	switch keyMsg.String() {
	case "e":
		return nil, &screens.NavigateResult{
			ResultData: "edit_metadata",
		}

	case "b":
		return nil, &screens.NavigateResult{
			ResultData: "edit_bursts",
		}

	case "f":
		return nil, &screens.NavigateResult{
			ResultData: "edit_facts",
		}
	}

	return nil, nil
}

// View implements the Screen interface.
//
// Returns:
//   - A string value.
//
// Side effects:
//   - None.
func (s *EventReviewScreen) View() string {
	content := s.renderContent()
	footer := s.renderFooter()

	return s.CreateView(s.breadcrumbs, content, footer)
}

// renderContent renders the event details with bursts and facts
// using the UIKit DetailView widget and primitives.
func (s *EventReviewScreen) renderContent() string {
	th := s.resolveTheme()

	var parts []string

	// Title.
	parts = append(parts, primitives.Title("Review Enrichment Results", th).
		MarginBottom(1).
		Render())

	// Event details section.
	parts = append(parts, s.renderEventDetails(th))

	// Inferred bursts section.
	parts = append(parts, s.renderBursts(th))

	// Inferred facts section.
	parts = append(parts, s.renderFacts(th))

	return primitives.JoinVertical(primitives.AlignLeft, parts...)
}

// renderEventDetails renders the event metadata using a DetailView widget.
func (s *EventReviewScreen) renderEventDetails(th theme.Theme) string {
	dv := widgets.NewDetailView(th).
		Title("Event Details")

	if s.event == nil {
		dv.Field("Status", "No event data")
		return dv.Render()
	}

	dv.Field("Text", s.event.Text).
		Field("Date", s.event.Date.Format("2006-01-02")).
		FieldIf("Company", s.event.Company).
		FieldIf("Project", s.event.Project)

	return dv.Render()
}

// renderBursts renders the inferred bursts list using UIKit primitives.
func (s *EventReviewScreen) renderBursts(th theme.Theme) string {
	var b strings.Builder

	b.WriteString(primitives.Subtitle("Inferred Bursts", th).
		MarginTop(1).
		Render())
	b.WriteString("\n")

	if len(s.bursts) == 0 {
		b.WriteString(primitives.Muted("  No bursts detected", th).Render())
		b.WriteString("\n")
		return b.String()
	}

	for i, burst := range s.bursts {
		b.WriteString(primitives.Body(fmt.Sprintf("  %d. %s", i+1, burst.Name), th).Render())
		b.WriteString("\n")
		if burst.Description != "" {
			b.WriteString(primitives.Muted("     "+burst.Description, th).Render())
			b.WriteString("\n")
		}
	}

	return b.String()
}

// renderFacts renders the inferred facts list using UIKit primitives.
func (s *EventReviewScreen) renderFacts(th theme.Theme) string {
	var b strings.Builder

	b.WriteString(primitives.Subtitle("Inferred Facts", th).
		MarginTop(1).
		Render())
	b.WriteString("\n")

	if len(s.facts) == 0 {
		b.WriteString(primitives.Muted("  No facts detected", th).Render())
		b.WriteString("\n")
		return b.String()
	}

	for i, fact := range s.facts {
		b.WriteString(primitives.Body(fmt.Sprintf("  %d. %s", i+1, fact.Text), th).Render())
		b.WriteString("\n")
	}

	return b.String()
}

// renderFooter renders footer with action shortcuts using UIKit badge primitives.
// Only shows applicable actions based on available data.
func (s *EventReviewScreen) renderFooter() string {
	th := s.resolveThemesTheme()

	badges := []*primitives.Badge{
		primitives.ConfirmBadge(th),
		primitives.HelpKeyBadge("e", "Edit metadata", th),
	}

	if len(s.bursts) > 0 {
		badges = append(badges, primitives.HelpKeyBadge("b", "Edit bursts", th))
	}

	if len(s.facts) > 0 {
		badges = append(badges, primitives.HelpKeyBadge("f", "Edit facts", th))
	}

	badges = append(badges,
		primitives.BackBadge(th),
		primitives.QuitBadge(th),
	)

	return primitives.RenderHelpFooter(th, badges...)
}

// resolveTheme returns the screen's theme or a default if none is set.
func (s *EventReviewScreen) resolveTheme() theme.Theme {
	if screenTheme := s.Theme(); screenTheme != nil {
		if th, ok := screenTheme.(theme.Theme); ok {
			return th
		}
	}
	return theme.Default()
}

// resolveThemesTheme returns the screen's theme as themes.Theme for badge rendering.
func (s *EventReviewScreen) resolveThemesTheme() themes.Theme {
	if screenTheme := s.Theme(); screenTheme != nil {
		if th, ok := screenTheme.(themes.Theme); ok {
			return th
		}
	}
	return themes.NewDefaultTheme()
}
