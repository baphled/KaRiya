package timeline

import (
	"fmt"
	"strings"

	"github.com/baphled/kariya/internal/cli/screens"
	"github.com/baphled/kariya/internal/cli/screens/base"
	"github.com/baphled/kariya/internal/domain/career"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
)

// State constant for state matrix tracking (REQUIRED)
const TimelineEventDetailState = "timeline_event_detail"

// TimelineEventDetailScreen displays detailed information about a career event.
//
// This screen provides:
// - Detailed view of event date, text, company, project
// - Display of tags, categories, and associated skills
// - Actions: Edit (e), Delete (d)
// - Back navigation with Escape, q, or backspace
//
// Usage:
//
//	screen := timeline.NewTimelineEventDetailScreen(event)
//	cmd, result := screen.Update(msg)
//	if result != nil && result.Type() == screens.ResultNavigate {
//	    navResult := result.(*screens.NavigateResult)
//	    if action, ok := navResult.ResultData.(map[string]interface{}); ok {
//	        // Handle action (edit, delete)
//	    }
//	}
//
// Related:
// - BaseScreen provides the foundation
// - docs/TUI_DEVELOPER_GUIDE.md (Screen patterns)
// - docs/TUI_STANDARDS.md (Keyboard shortcuts)
type TimelineEventDetailScreen struct {
	*base.BaseScreen
	event *career.CareerEvent
}

// NewTimelineEventDetailScreen creates a new event detail screen.
//
// The screen:
// - Displays all event fields in a formatted layout
// - Shows optional fields only when present
// - Provides edit and delete actions
// - Supports back navigation
//
// Parameters:
//   - event: The career event to display
func NewTimelineEventDetailScreen(event *career.CareerEvent) *TimelineEventDetailScreen {
	return &TimelineEventDetailScreen{
		BaseScreen: base.NewBaseScreen(),
		event:      event,
	}
}

// Update handles messages and actions.
func (s *TimelineEventDetailScreen) Update(msg tea.Msg) (tea.Cmd, screens.ScreenResult) {
	switch msg := msg.(type) {
	case tea.WindowSizeMsg:
		s.SetTerminalInfo(msg.Width, msg.Height)
		return nil, nil

	case tea.KeyMsg:
		switch msg.String() {
		case "esc", "q", "backspace":
			// Back to event list
			return nil, &screens.CancelResult{}

		case "e":
			// Edit event
			return nil, &screens.NavigateResult{
				ResultData: map[string]interface{}{
					"action": "edit",
					"event":  s.event,
				},
			}

		case "d":
			// Delete event
			return nil, &screens.NavigateResult{
				ResultData: map[string]interface{}{
					"action": "delete",
					"event":  s.event,
				},
			}
		}
	}

	return nil, nil
}

// View renders the event detail screen using StandardView.
func (s *TimelineEventDetailScreen) View() string {
	var b strings.Builder

	// Title
	titleStyle := lipgloss.NewStyle().Bold(true).Foreground(lipgloss.Color("12"))
	b.WriteString(titleStyle.Render("Event Details"))
	b.WriteString("\n\n")

	// Date
	labelStyle := lipgloss.NewStyle().Bold(true).Foreground(lipgloss.Color("10"))
	b.WriteString(labelStyle.Render("Date: "))
	b.WriteString(s.event.Date.Format("2006-01-02 (Monday)"))
	b.WriteString("\n\n")

	// Company (if present)
	if s.event.Company != "" {
		b.WriteString(labelStyle.Render("Company: "))
		b.WriteString(s.event.Company)
		b.WriteString("\n\n")
	}

	// Project (if present)
	if s.event.Project != "" {
		b.WriteString(labelStyle.Render("Project: "))
		b.WriteString(s.event.Project)
		b.WriteString("\n\n")
	}

	// Event Text
	b.WriteString(labelStyle.Render("Description:"))
	b.WriteString("\n")
	b.WriteString(s.wrapText(s.event.Text, s.Width()-4))
	b.WriteString("\n\n")

	// Tags (if present)
	if len(s.event.Tags) > 0 {
		b.WriteString(labelStyle.Render("Tags: "))
		tagStyle := lipgloss.NewStyle().
			Foreground(lipgloss.Color("14")).
			Padding(0, 1).
			Background(lipgloss.Color("8"))
		for i, tag := range s.event.Tags {
			if i > 0 {
				b.WriteString(" ")
			}
			b.WriteString(tagStyle.Render(tag))
		}
		b.WriteString("\n\n")
	}

	// Categories (if present)
	if len(s.event.Categories) > 0 {
		b.WriteString(labelStyle.Render("Categories: "))
		for i, cat := range s.event.Categories {
			if i > 0 {
				b.WriteString(", ")
			}
			b.WriteString(cat)
		}
		b.WriteString("\n\n")
	}

	// Skills (if present)
	if len(s.event.Skills) > 0 {
		b.WriteString(labelStyle.Render(fmt.Sprintf("Skills: ")))
		skillCount := len(s.event.Skills)
		if skillCount == 1 {
			b.WriteString("1 skill associated")
		} else {
			b.WriteString(fmt.Sprintf("%d skills associated", skillCount))
		}
		b.WriteString("\n\n")
	}

	// Metadata
	b.WriteString(labelStyle.Render("Created: "))
	b.WriteString(s.event.CreatedAt.Format("2006-01-02 15:04"))
	b.WriteString("  ")
	b.WriteString(labelStyle.Render("Updated: "))
	b.WriteString(s.event.UpdatedAt.Format("2006-01-02 15:04"))

	content := b.String()

	// Footer with actions
	footer := "e: Edit  d: Delete  Esc/q/Backspace: Back to timeline"

	// Use BaseScreen's CreateView helper for StandardView integration
	breadcrumbs := []string{"Main Menu", "Timeline", "Event Details"}
	return s.CreateView(breadcrumbs, content, footer)
}

// wrapText wraps text to fit within the specified width.
func (s *TimelineEventDetailScreen) wrapText(text string, width int) string {
	if width < 20 {
		width = 20
	}

	words := strings.Fields(text)
	if len(words) == 0 {
		return text
	}

	var lines []string
	var currentLine strings.Builder

	for _, word := range words {
		// Check if adding this word would exceed width
		if currentLine.Len() > 0 && currentLine.Len()+1+len(word) > width {
			// Start new line
			lines = append(lines, currentLine.String())
			currentLine.Reset()
		}

		// Add word to current line
		if currentLine.Len() > 0 {
			currentLine.WriteString(" ")
		}
		currentLine.WriteString(word)
	}

	// Add last line
	if currentLine.Len() > 0 {
		lines = append(lines, currentLine.String())
	}

	return strings.Join(lines, "\n")
}

// GetEvent returns the event being displayed.
func (s *TimelineEventDetailScreen) GetEvent() *career.CareerEvent {
	return s.event
}
