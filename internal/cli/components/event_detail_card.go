package components

import (
	"fmt"
	"strings"

	"github.com/baphled/kariya/internal/cli/themes"
	"github.com/baphled/kariya/internal/domain/career"
	"github.com/charmbracelet/lipgloss"
)

// EventDetailCard renders a career event in a styled card format.
// This is a reusable component that can be used across different intents.
type EventDetailCard struct {
	event *career.Event
	theme themes.Theme
}

// NewEventDetailCard creates a new EventDetailCard component.
//
// Expected:
//   - event must be valid.
//   - th must be a valid theme instance (can be nil).
//
// Returns:
//   - A fully initialized EventDetailCard ready for use.
//
// Side effects:
//   - None.
func NewEventDetailCard(event *career.Event, theme themes.Theme) *EventDetailCard {
	return &EventDetailCard{
		event: event,
		theme: theme,
	}
}

// Render renders the event detail card as a string.
//
// Returns:
//   - A string value.
//
// Side effects:
//   - None.
func (c *EventDetailCard) Render() string {
	if c.event == nil {
		return "No event selected."
	}

	var content strings.Builder
	content.WriteString("\nEvent Details\n\n")

	// Event header
	fmt.Fprintf(&content, "Date: %s\n", c.event.Date.Format("2006-01-02"))

	if c.event.Company != "" {
		fmt.Fprintf(&content, "Company: %s\n", c.event.Company)
	}
	if c.event.Project != "" {
		fmt.Fprintf(&content, "Project: %s\n", c.event.Project)
	}

	fmt.Fprintf(&content, "\nText:\n%s\n", c.event.Text)

	// Tags and categories
	if len(c.event.Tags) > 0 {
		fmt.Fprintf(&content, "\nTags: %s\n", strings.Join(c.event.Tags, ", "))
	}
	if len(c.event.Categories) > 0 {
		fmt.Fprintf(&content, "Categories: %s\n", strings.Join(c.event.Categories, ", "))
	}

	// Skills (just show count, as we only have IDs)
	if len(c.event.Skills) > 0 {
		fmt.Fprintf(&content, "Skills: %d associated\n", len(c.event.Skills))
	}

	// Apply themed card styling
	var cardStyle lipgloss.Style
	if c.theme != nil {
		cardStyle = c.theme.Styles().CardBase
	} else {
		// Fallback if no theme available
		cardStyle = lipgloss.NewStyle().
			Padding(1, 2).
			BorderStyle(lipgloss.RoundedBorder())
	}

	return cardStyle.Render(content.String())
}

// RenderEventDetailCard is a helper function that creates and renders an event detail card.
//
// Expected:
//   - event must be valid.
//   - th must be a valid theme instance (can be nil).
//
// Returns:
//   - A string value.
//
// Side effects:
//   - None.
func RenderEventDetailCard(event *career.Event, theme themes.Theme) string {
	card := NewEventDetailCard(event, theme)
	return card.Render()
}
