package intents

import (
	"fmt"
	"sort"
	"strings"
	"time"

	"github.com/baphled/kariya/internal/cli/styles"
	"github.com/baphled/kariya/internal/domain/career"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
)

// Custom message types for BrowseTimeline state transitions.

// EventSelectedMsg indicates the user selected an event.
type EventSelectedMsg struct {
	Event *career.CareerEvent
	Index int
}

// FilterChangedMsg indicates the filters have changed.
type FilterChangedMsg struct {
	Filters *TimelineFilters
}

// BrowseTimelineIntent implements the Intent interface for browsing career events.
// It owns the complete lifecycle of timeline browsing, including:
// - Displaying a filtered and sorted timeline of events
// - Selecting and viewing event details
// - Returning the selected event or cancelling
type BrowseTimelineIntent struct {
	// context is the input context passed to the intent.
	context *BrowseTimelineContext

	// state represents the current state of the intent.
	state *BrowseTimelineModel

	// active indicates whether this intent is currently active.
	active bool

	// result is the final result of the intent (set when complete).
	result *IntentResult[*BrowseTimelineResult]
}

// NewBrowseTimelineIntent creates a new BrowseTimeline intent.
func NewBrowseTimelineIntent(context *BrowseTimelineContext) (*BrowseTimelineIntent, error) {
	// Validate the context.
	if err := context.Validate(); err != nil {
		return nil, err
	}

	return &BrowseTimelineIntent{
		context: context,
		state: &BrowseTimelineModel{
			context:        context,
			currentState:   BrowseStateTimeline,
			filteredEvents: context.Events,
			selectedIndex:  0,
			filters:        context.InitialFilters,
			selectedFacts:  make([]*career.Fact, 0),
			viewedEvents:   make([]*career.CareerEvent, 0),
		},
		active: true,
	}, nil
}

// Init is called when the intent is activated.
func (i *BrowseTimelineIntent) Init() tea.Cmd {
	// Initialize filtered events with the provided events.
	i.state.filteredEvents = i.context.Events
	if len(i.state.filteredEvents) > 0 {
		i.state.selectedEvent = i.state.filteredEvents[0]
	}
	return nil
}

// Update processes a message in the intent.
func (i *BrowseTimelineIntent) Update(msg tea.Msg) tea.Cmd {
	if !i.active {
		return nil
	}

	switch i.state.currentState {
	case BrowseStateTimeline:
		return i.updateTimelineView(msg)

	case BrowseStateEventDetail:
		return i.updateEventDetail(msg)
	}

	return nil
}

// updateTimelineView handles messages while viewing the timeline.
func (i *BrowseTimelineIntent) updateTimelineView(msg tea.Msg) tea.Cmd {
	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch msg.String() {
		case "enter":
			// Select current event and move to detail view.
			if len(i.state.filteredEvents) > 0 && i.state.selectedIndex < len(i.state.filteredEvents) {
				i.state.selectedEvent = i.state.filteredEvents[i.state.selectedIndex]
				i.state.viewedEvents = append(i.state.viewedEvents, i.state.selectedEvent)
				i.state.currentState = BrowseStateEventDetail
			}
			return nil

		case "up", "k":
			// Move selection up.
			if i.state.selectedIndex > 0 {
				i.state.selectedIndex--
				if len(i.state.filteredEvents) > 0 {
					i.state.selectedEvent = i.state.filteredEvents[i.state.selectedIndex]
				}
			}
			return nil

		case "down", "j":
			// Move selection down.
			if i.state.selectedIndex < len(i.state.filteredEvents)-1 {
				i.state.selectedIndex++
				if len(i.state.filteredEvents) > 0 {
					i.state.selectedEvent = i.state.filteredEvents[i.state.selectedIndex]
				}
			}
			return nil

		case "q", "ctrl+c":
			// Cancel without selection.
			i.setCancelled()
			return nil

		case "esc":
			// Go back (no-op at timeline view).
			i.setCancelled()
			return nil
		}

	case EventSelectedMsg:
		// Event was selected (possibly by router or other component).
		i.state.selectedEvent = msg.Event
		i.state.selectedIndex = msg.Index
		i.state.viewedEvents = append(i.state.viewedEvents, msg.Event)
		i.state.currentState = BrowseStateEventDetail
		return nil

	case FilterChangedMsg:
		// Filters have changed - re-filter events.
		i.state.filters = msg.Filters
		i.applyFilters()
		i.state.selectedIndex = 0
		if len(i.state.filteredEvents) > 0 {
			i.state.selectedEvent = i.state.filteredEvents[0]
		}
		return nil
	}

	return nil
}

// updateEventDetail handles messages while viewing event details.
func (i *BrowseTimelineIntent) updateEventDetail(msg tea.Msg) tea.Cmd {
	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch msg.String() {
		case "enter":
			// Confirm selection and return event.
			i.setCompleted()
			return nil

		case "esc":
			// Go back to timeline.
			i.state.currentState = BrowseStateTimeline
			return nil

		case "q", "ctrl+c":
			// Cancel.
			i.setCancelled()
			return nil
		}
	}

	return nil
}

// applyFilters filters the events based on current filter state.
func (i *BrowseTimelineIntent) applyFilters() {
	filtered := make([]*career.CareerEvent, 0)

	for _, event := range i.context.Events {
		// Apply search text filter.
		if i.state.filters.SearchText != "" {
			if !strings.Contains(strings.ToLower(event.Text), strings.ToLower(i.state.filters.SearchText)) {
				continue
			}
		}

		// Apply tag filter.
		if len(i.state.filters.Tags) > 0 {
			hasTag := false
			for _, tag := range i.state.filters.Tags {
				for _, eventTag := range event.Tags {
					if eventTag == tag {
						hasTag = true
						break
					}
				}
				if hasTag {
					break
				}
			}
			if !hasTag {
				continue
			}
		}

		filtered = append(filtered, event)
	}

	// Apply sorting.
	sort.Slice(filtered, func(a, b int) bool {
		switch i.state.filters.SortBy {
		case "text":
			if i.state.filters.SortOrder == "asc" {
				return filtered[a].Text < filtered[b].Text
			}
			return filtered[a].Text > filtered[b].Text

		case "relevance":
			// For now, relevance is the same as date (most recent first).
			if i.state.filters.SortOrder == "asc" {
				return filtered[a].CreatedAt.Before(filtered[b].CreatedAt)
			}
			return filtered[a].CreatedAt.After(filtered[b].CreatedAt)

		default: // date
			if i.state.filters.SortOrder == "asc" {
				return filtered[a].Date.Before(filtered[b].Date)
			}
			return filtered[a].Date.After(filtered[b].Date)
		}
	})

	i.state.filteredEvents = filtered
}

// View renders the intent's current state.
func (i *BrowseTimelineIntent) View() string {
	switch i.state.currentState {
	case BrowseStateTimeline:
		return i.viewTimeline()

	case BrowseStateEventDetail:
		return i.viewEventDetail()
	}

	return ""
}

// viewTimeline renders the timeline view with all events.
func (i *BrowseTimelineIntent) viewTimeline() string {
	var content strings.Builder
	content.WriteString("\nBrowse Timeline\n\n")

	if len(i.state.filteredEvents) == 0 {
		content.WriteString("No events found.\n")
	} else {
		// Show current filter info.
		if i.state.filters.SearchText != "" {
			content.WriteString(fmt.Sprintf("Search: %s\n", i.state.filters.SearchText))
		}
		if len(i.state.filters.Tags) > 0 {
			content.WriteString(fmt.Sprintf("Tags: %s\n", strings.Join(i.state.filters.Tags, ", ")))
		}
		content.WriteString("\n")

		// Show event list.
		for idx, event := range i.state.filteredEvents {
			prefix := "  "
			if idx == i.state.selectedIndex {
				prefix = "> "
			}

			dateStr := event.Date.Format("2006-01-02")
			// Truncate text to first 50 chars
			text := event.Text
			if len(text) > 50 {
				text = text[:50] + "..."
			}
			content.WriteString(fmt.Sprintf("%s[%s] %s\n", prefix, dateStr, text))
		}
	}

	// Apply card styling.
	cardStyle := lipgloss.NewStyle().
		Padding(1, 2).
		BorderStyle(lipgloss.RoundedBorder()).
		BorderForeground(styles.ColorBorder).
		Background(styles.ColorBackgroundCard).
		Foreground(styles.ColorTextPrimary)

	card := cardStyle.Render(content.String())

	// Add footer with instructions.
	footerStyle := lipgloss.NewStyle().
		Foreground(styles.ColorTextSecondary).
		MarginTop(1)

	footer := footerStyle.Render("↑/↓ or k/j to navigate, Enter to select, q to cancel")

	return lipgloss.JoinVertical(lipgloss.Left, card, footer)
}

// viewEventDetail renders the event detail view.
func (i *BrowseTimelineIntent) viewEventDetail() string {
	if i.state.selectedEvent == nil {
		return "No event selected."
	}

	var content strings.Builder
	content.WriteString("\nEvent Details\n\n")

	// Event header.
	content.WriteString(fmt.Sprintf("Date: %s\n", i.state.selectedEvent.Date.Format("2006-01-02")))

	if i.state.selectedEvent.Company != "" {
		content.WriteString(fmt.Sprintf("Company: %s\n", i.state.selectedEvent.Company))
	}
	if i.state.selectedEvent.Project != "" {
		content.WriteString(fmt.Sprintf("Project: %s\n", i.state.selectedEvent.Project))
	}

	content.WriteString(fmt.Sprintf("\nText:\n%s\n", i.state.selectedEvent.Text))

	// Tags and categories.
	if len(i.state.selectedEvent.Tags) > 0 {
		content.WriteString(fmt.Sprintf("\nTags: %s\n", strings.Join(i.state.selectedEvent.Tags, ", ")))
	}
	if len(i.state.selectedEvent.Categories) > 0 {
		content.WriteString(fmt.Sprintf("Categories: %s\n", strings.Join(i.state.selectedEvent.Categories, ", ")))
	}

	// Apply card styling.
	cardStyle := lipgloss.NewStyle().
		Padding(1, 2).
		BorderStyle(lipgloss.RoundedBorder()).
		BorderForeground(styles.ColorBorder).
		Background(styles.ColorBackgroundCard).
		Foreground(styles.ColorTextPrimary)

	card := cardStyle.Render(content.String())

	// Add footer with instructions.
	footerStyle := lipgloss.NewStyle().
		Foreground(styles.ColorTextSecondary).
		MarginTop(1)

	footer := footerStyle.Render("Enter to confirm, Esc to go back, q to cancel")

	return lipgloss.JoinVertical(lipgloss.Left, card, footer)
}

// Result returns the final result of the intent.
func (i *BrowseTimelineIntent) Result() *IntentResult[interface{}] {
	if i.result == nil {
		return &IntentResult[interface{}]{
			Status: Cancelled,
		}
	}

	return &IntentResult[interface{}]{
		Status:   i.result.Status,
		Data:     i.result.Data,
		Error:    i.result.Error,
		Metadata: i.result.Metadata,
	}
}

// Helper methods for result management.

func (i *BrowseTimelineIntent) setCompleted() {
	i.result = &IntentResult[*BrowseTimelineResult]{
		Status: Completed,
		Data: &BrowseTimelineResult{
			SelectedEvent: i.state.selectedEvent,
			FinalFilters:  i.state.filters,
			ViewedEvents:  i.state.viewedEvents,
			SelectedFacts: i.state.selectedFacts,
		},
		Metadata: map[string]interface{}{
			"selected_index": i.state.selectedIndex,
			"viewed_count":   len(i.state.viewedEvents),
			"timestamp":      time.Now(),
		},
	}
	i.active = false
}

func (i *BrowseTimelineIntent) setCancelled() {
	i.result = &IntentResult[*BrowseTimelineResult]{
		Status: Cancelled,
	}
	i.active = false
}

func (i *BrowseTimelineIntent) setFailed(code, message string, cause error) {
	i.result = &IntentResult[*BrowseTimelineResult]{
		Status: Failed,
		Error: &IntentError{
			Code:    code,
			Message: message,
			Cause:   cause,
		},
	}
	i.active = false
}
