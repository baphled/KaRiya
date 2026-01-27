package capture

import (
	"fmt"
	"strings"

	"github.com/baphled/kariya/internal/cli/screens"
	"github.com/baphled/kariya/internal/cli/screens/base"
	"github.com/baphled/kariya/internal/domain/career"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
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
// - internal/cli/intents/capture_event.go (ReviewInferredEventState)
// - tasks/tasks-42-tui-architecture-refactor.md (Phase 1: CaptureEvent Migration)
type EventReviewScreen struct {
	*base.BaseScreen

	// event being reviewed
	event *career.CareerEvent

	// bursts inferred from event
	bursts []*career.Burst

	// facts inferred from event
	facts []*career.Fact

	// breadcrumbs for the view header
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
// Returns a EventReviewScreen displaying event details with bursts and facts.
func NewEventReviewScreen(
	breadcrumbs []string,
	event *career.CareerEvent,
	bursts []*career.Burst,
	facts []*career.Fact,
) *EventReviewScreen {
	return &EventReviewScreen{
		BaseScreen:  base.NewBaseScreen(),
		event:       event,
		bursts:      bursts,
		facts:       facts,
		breadcrumbs: breadcrumbs,
	}
}

// Update implements the Screen interface.
//
// Handles:
// - Enter → returns SubmitResult with event, bursts, facts
// - e/b/f → returns NavigateResult with edit action
// - Esc → returns CancelResult
// - WindowSizeMsg → updates dimensions
func (s *EventReviewScreen) Update(msg tea.Msg) (tea.Cmd, screens.ScreenResult) {
	// Handle window size via BaseScreen
	if cmd := s.BaseScreen.HandleWindowSizeMsg(msg); cmd != nil {
		return cmd, nil
	}

	// Handle key messages
	if keyMsg, ok := msg.(tea.KeyMsg); ok {
		switch keyMsg.String() {
		case "enter":
			// Confirm review - return SubmitResult with all data
			return nil, &screens.SubmitResult{
				FormData: map[string]interface{}{
					"event":  s.event,
					"bursts": s.bursts,
					"facts":  s.facts,
				},
			}

		case "e":
			// Navigate to edit metadata
			return nil, &screens.NavigateResult{
				ResultData: "edit_metadata",
			}

		case "b":
			// Navigate to edit bursts
			return nil, &screens.NavigateResult{
				ResultData: "edit_bursts",
			}

		case "f":
			// Navigate to edit facts
			return nil, &screens.NavigateResult{
				ResultData: "edit_facts",
			}

		case "esc":
			// Cancel and return to form
			return nil, &screens.CancelResult{}
		}
	}

	return nil, nil
}

// View implements the Screen interface.
//
// Renders the review screen with event details, bursts, and facts.
func (s *EventReviewScreen) View() string {
	content := s.renderContent()
	footer := s.renderFooter()

	return s.CreateView(s.breadcrumbs, content, footer)
}

// renderContent renders the event details with bursts and facts.
func (s *EventReviewScreen) renderContent() string {
	var b strings.Builder

	b.WriteString("\n")

	// Title with themed styling
	title := "Review Enrichment Results"
	theme := s.Theme()
	if theme != nil {
		// Use card style for consistent theming
		cardStyle := s.getCardStyle()
		b.WriteString(cardStyle.Render(title) + "\n\n")
	} else {
		// Fallback: plain text
		b.WriteString(title + "\n")
		b.WriteString("═════════════════════\n\n")
	}

	// Event details
	s.renderEventDetails(&b)
	b.WriteString("\n")

	// Inferred bursts
	s.renderBursts(&b)
	b.WriteString("\n")

	// Inferred facts
	s.renderFacts(&b)

	return b.String()
}

// getCardStyle returns a themed card style for content sections.
func (s *EventReviewScreen) getCardStyle() lipgloss.Style {
	theme := s.Theme()
	if theme != nil {
		// Try to cast to Theme interface
		if t, ok := theme.(interface {
			PrimaryColor() lipgloss.Color
			BorderColor() lipgloss.Color
		}); ok {
			return lipgloss.NewStyle().
				Foreground(t.PrimaryColor()).
				Bold(true)
		}
	}
	return lipgloss.NewStyle()
}

// renderEventDetails renders the event metadata.
func (s *EventReviewScreen) renderEventDetails(b *strings.Builder) {
	sectionStyle := s.getSectionHeaderStyle()
	dimStyle := s.getDimStyle()

	b.WriteString(sectionStyle.Render("Event Details:") + "\n")
	b.WriteString(dimStyle.Render("─────────────") + "\n")

	if s.event == nil {
		b.WriteString("  No event data\n")
		return
	}

	b.WriteString(fmt.Sprintf("  %s\n", s.event.Text))
	b.WriteString(fmt.Sprintf("  Date: %s\n", s.event.Date.Format("2006-01-02")))

	if s.event.Company != "" {
		b.WriteString(fmt.Sprintf("  Company: %s\n", s.event.Company))
	}
	if s.event.Project != "" {
		b.WriteString(fmt.Sprintf("  Project: %s\n", s.event.Project))
	}
}

// getSectionHeaderStyle returns a themed style for section headers.
func (s *EventReviewScreen) getSectionHeaderStyle() lipgloss.Style {
	theme := s.Theme()
	if theme != nil {
		if t, ok := theme.(interface{ SecondaryColor() lipgloss.Color }); ok {
			return lipgloss.NewStyle().
				Foreground(t.SecondaryColor()).
				Bold(true)
		}
	}
	return lipgloss.NewStyle()
}

// getDimStyle returns a themed style for dimmed text.
func (s *EventReviewScreen) getDimStyle() lipgloss.Style {
	theme := s.Theme()
	if theme != nil {
		if t, ok := theme.(interface{ MutedColor() lipgloss.Color }); ok {
			return lipgloss.NewStyle().
				Foreground(t.MutedColor())
		}
	}
	return lipgloss.NewStyle()
}

// renderBursts renders the inferred bursts list.
func (s *EventReviewScreen) renderBursts(b *strings.Builder) {
	sectionStyle := s.getSectionHeaderStyle()
	dimStyle := s.getDimStyle()

	b.WriteString(sectionStyle.Render("Inferred Bursts:") + "\n")
	b.WriteString(dimStyle.Render("────────────────") + "\n")

	if len(s.bursts) == 0 {
		b.WriteString(dimStyle.Render("  No bursts detected") + "\n")
		return
	}

	for i, burst := range s.bursts {
		fmt.Fprintf(b, "  %d. %s\n", i+1, burst.Name)
		if burst.Description != "" {
			fmt.Fprintf(b, "     %s\n", burst.Description)
		}
	}
}

// renderFacts renders the inferred facts list.
func (s *EventReviewScreen) renderFacts(b *strings.Builder) {
	sectionStyle := s.getSectionHeaderStyle()
	dimStyle := s.getDimStyle()

	b.WriteString(sectionStyle.Render("Inferred Facts:") + "\n")
	b.WriteString(dimStyle.Render("───────────────") + "\n")

	if len(s.facts) == 0 {
		b.WriteString(dimStyle.Render("  No facts detected") + "\n")
		return
	}

	for i, fact := range s.facts {
		b.WriteString(fmt.Sprintf("  %d. %s\n", i+1, fact.Text))
	}
}

// renderFooter renders footer with action shortcuts.
// Only shows applicable actions based on available data.
func (s *EventReviewScreen) renderFooter() string {
	actions := []string{"Enter: Confirm", "e: Edit metadata"}

	// Only show burst editing if bursts exist
	if len(s.bursts) > 0 {
		actions = append(actions, "b: Edit bursts")
	}

	// Only show fact editing if facts exist
	if len(s.facts) > 0 {
		actions = append(actions, "f: Edit facts")
	}

	actions = append(actions, "Esc: Back", "q: Quit")

	return strings.Join(actions, "  ")
}
