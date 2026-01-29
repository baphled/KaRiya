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
func NewEventDetailCard(event *career.Event, theme themes.Theme) *EventDetailCard {
	return &EventDetailCard{
		event: event,
		theme: theme,
	}
}

// Render renders the event detail card as a string.
func (c *EventDetailCard) Render() string {
	if c.event == nil {
		return "No event selected."
	}

	var content strings.Builder
	content.WriteString("\nEvent Details\n\n")

	// Event header
	content.WriteString(fmt.Sprintf("Date: %s\n", c.event.Date.Format("2006-01-02")))

	if c.event.Company != "" {
		content.WriteString(fmt.Sprintf("Company: %s\n", c.event.Company))
	}
	if c.event.Project != "" {
		content.WriteString(fmt.Sprintf("Project: %s\n", c.event.Project))
	}

	content.WriteString(fmt.Sprintf("\nText:\n%s\n", c.event.Text))

	// Tags and categories
	if len(c.event.Tags) > 0 {
		content.WriteString(fmt.Sprintf("\nTags: %s\n", strings.Join(c.event.Tags, ", ")))
	}
	if len(c.event.Categories) > 0 {
		content.WriteString(fmt.Sprintf("Categories: %s\n", strings.Join(c.event.Categories, ", ")))
	}

	// Skills (just show count, as we only have IDs)
	if len(c.event.Skills) > 0 {
		content.WriteString(fmt.Sprintf("Skills: %d associated\n", len(c.event.Skills)))
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
// This is a convenience function for quick usage without creating a struct instance.
func RenderEventDetailCard(event *career.Event, theme themes.Theme) string {
	card := NewEventDetailCard(event, theme)
	return card.Render()
}
