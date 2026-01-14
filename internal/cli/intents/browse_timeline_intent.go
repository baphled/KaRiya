package intents

import (
	"context"
	"fmt"
	"sort"
	"strings"
	"time"

	"github.com/baphled/kariya/internal/cli/components"
	"github.com/baphled/kariya/internal/cli/navigation"
	"github.com/baphled/kariya/internal/cli/themes"
	"github.com/baphled/kariya/internal/domain/career"
	"github.com/charmbracelet/bubbles/table"
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

// RequestEditEventMsg requests that the app route to CaptureEvent intent for editing.
// This is sent to the app router which will activate CaptureEvent with PreviousEvent set.
type RequestEditEventMsg struct {
	Event *career.CareerEvent
}

// BrowseTimelineIntent implements the Intent interface for browsing career events.
// It owns the complete lifecycle of timeline browsing, including:
// - Displaying a filtered and sorted timeline of events
// - Selecting and viewing event details
// - Returning the selected event or cancelling
type BrowseTimelineIntent struct {
	// Embed BaseIntent for terminal awareness, logo, and state management
	*BaseIntent

	// context is the input context passed to the intent.
	context *BrowseTimelineContext

	// state represents the current state of the intent.
	state *BrowseTimelineModel

	// table is the table model for displaying events
	table *table.Model

	// listContainer provides table-based list UI
	listContainer *components.TableListContainer

	// navHandler handles list navigation
	navHandler *navigation.ListNavigationHandler

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

	// Create table model for events
	columns := []table.Column{
		{Title: "Date", Width: 12},
		{Title: "Event", Width: 50},
		{Title: "Company", Width: 20},
	}

	t := table.New(
		table.WithColumns(columns),
		table.WithRows([]table.Row{}),
		table.WithFocused(true),
		table.WithHeight(15),
		table.WithWidth(100),
	)

	// Apply default styles initially - theme styles will be applied in Init()
	// when the theme manager is available via BaseIntent
	t.SetStyles(table.DefaultStyles())

	// Create BaseIntent for terminal awareness and state management
	base := NewBaseIntent()

	intent := &BrowseTimelineIntent{
		BaseIntent: base,
		context:    context,
		state: &BrowseTimelineModel{
			context:        context,
			currentState:   BrowseStateTimeline,
			filteredEvents: context.Events,
			selectedIndex:  0,
			filters:        context.InitialFilters,
			selectedFacts:  make([]*career.Fact, 0),
			viewedEvents:   make([]*career.CareerEvent, 0),
		},
		table:         &t,
		listContainer: components.NewTableListContainer(t, "Browse Timeline", 100),
		active:        true,
	}

	// Initialize navigation handler
	intent.navHandler = navigation.NewListNavigationHandler(intent)

	return intent, nil
}

// Init is called when the intent is activated.
func (i *BrowseTimelineIntent) Init() tea.Cmd {
	// Apply themed table styles if theme is available
	if theme := i.Theme(); theme != nil {
		i.table.SetStyles(themes.NewThemedTableStyles(theme))
	}

	// Apply initial filters and sorting to the provided events.
	i.applyFilters()

	if len(i.state.filteredEvents) > 0 {
		i.state.selectedEvent = i.state.filteredEvents[0]
	}
	i.updateTableRows()
	return nil
}

// updateTableRows updates the table rows based on filtered events
func (i *BrowseTimelineIntent) updateTableRows() {
	pageSize := 15
	total := len(i.state.filteredEvents)

	// Determine which page current selection is on
	page := 0
	if pageSize > 0 && i.state.selectedIndex >= 0 {
		page = i.state.selectedIndex / pageSize
	}

	start := page * pageSize
	end := start + pageSize
	if end > total {
		end = total
	}

	pageEvents := i.state.filteredEvents[start:end]

	rows := make([]table.Row, 0, len(pageEvents))
	for idx, event := range pageEvents {
		realIdx := start + idx

		// Use centralized indicator formatting
		dateStr := i.navHandler.FormatRowText(realIdx, event.Date.Format("2006-01-02"))

		// Truncate text to first 50 chars
		text := event.Text
		if len(text) > 50 {
			text = text[:50] + "..."
		}
		company := event.Company
		if company == "" {
			company = "-"
		}
		rows = append(rows, table.Row{dateStr, text, company})
	}

	i.table.SetRows(rows)

	// Calculate relative cursor position for this page
	relativeCursor := 0
	if i.state.selectedIndex >= start && i.state.selectedIndex < end {
		relativeCursor = i.state.selectedIndex - start
	}

	// Set table cursor to relative position within the page
	// This makes the table highlight the correct row with its Selected style
	i.table.SetCursor(relativeCursor)

	// Sync the container's selectedIdx to match our relative cursor
	// This ensures SetTable() will push the correct cursor position to the table
	i.listContainer.SetSelectedIdx(relativeCursor)

	// Update the container with the modified table
	// This will call validateAndSyncIdx() which will use our relative cursor position
	i.listContainer.SetTable(*i.table)
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

	case BrowseStateDeleteConfirm:
		return i.updateDeleteConfirm(msg)
	}

	return nil
}

// updateTimelineView handles messages while viewing the timeline.
func (i *BrowseTimelineIntent) updateTimelineView(msg tea.Msg) tea.Cmd {
	switch msg := msg.(type) {
	case tea.KeyMsg:
		// Use MessageInterceptor for global keys (quit, help, back)
		return NewMessageInterceptor().
			OnQuit(StandardQuitHandler()).
			OnHelp(StandardHelpHandler(i.BaseIntent)).
			OnBack(func() tea.Cmd {
				// At root state, back means cancel and return to main menu
				i.setCancelled()
				return nil
			}).
			InterceptOr(msg, func() tea.Cmd {
				// Try list navigation handler
				if i.navHandler.HandleKey(msg.String()) {
					return nil
				}

				// Handle intent-specific keys
				switch msg.String() {
				case "enter":
					// Select current event and move to detail view.
					if len(i.state.filteredEvents) > 0 {
						if i.state.selectedIndex < len(i.state.filteredEvents) {
							i.state.selectedEvent = i.state.filteredEvents[i.state.selectedIndex]
							i.state.viewedEvents = append(i.state.viewedEvents, i.state.selectedEvent)
							i.state.currentState = BrowseStateEventDetail
						}
					}
					return nil
				}
				return nil
			})

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
		i.updateTableRows()
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
		// Use MessageInterceptor for global keys (quit, help, back)
		return NewMessageInterceptor().
			OnQuit(StandardQuitHandler()).
			OnHelp(StandardHelpHandler(i.BaseIntent)).
			OnBack(func() tea.Cmd {
				// Go back to timeline view
				i.state.currentState = BrowseStateTimeline
				return nil
			}).
			InterceptOr(msg, func() tea.Cmd {
				// Handle intent-specific keys
				switch msg.String() {
				case "enter":
					// Confirm selection and return event.
					i.setCompleted()
					return nil

				case "e":
					// Edit event - send message to app to route to CaptureEvent intent
					if i.state.selectedEvent != nil {
						// Send RequestEditEventMsg which app router will handle
						return func() tea.Msg {
							return RequestEditEventMsg{Event: i.state.selectedEvent}
						}
					}
					return nil

				case "d":
					// Delete event - go to confirmation
					if i.state.selectedEvent != nil && i.context.CLIEventService != nil {
						i.state.currentState = BrowseStateDeleteConfirm
						i.state.deleteError = nil
					}
					return nil
				}
				return nil
			})
	}

	return nil
}

// updateDeleteConfirm handles delete confirmation.
func (i *BrowseTimelineIntent) updateDeleteConfirm(msg tea.Msg) tea.Cmd {
	switch msg := msg.(type) {
	case tea.KeyMsg:
		// Use MessageInterceptor for consistent escape key handling
		return NewMessageInterceptor().
			OnQuit(StandardQuitHandler()).
			OnHelp(StandardHelpHandler(i.BaseIntent)).
			OnBack(func() tea.Cmd {
				// Cancel delete and return to event detail
				i.state.currentState = BrowseStateEventDetail
				i.state.deleteError = nil
				return nil
			}).
			InterceptOr(msg, func() tea.Cmd {
				// Handle confirmation keys
				switch msg.String() {
				case "y", "Y":
					// Confirm delete
					if i.state.selectedEvent != nil && i.context.CLIEventService != nil {
						err := i.context.CLIEventService.DeleteEvent(
							i.getContext(),
							i.state.selectedEvent.ID,
						)
						if err != nil {
							i.state.deleteError = err
							return nil
						}

						// Remove from filtered events
						i.removeEventFromList(i.state.selectedEvent.ID)

						// Clear selection and go back to timeline
						i.state.selectedEvent = nil
						i.state.currentState = BrowseStateTimeline

						// Select first event if available
						if len(i.state.filteredEvents) > 0 {
							i.state.selectedIndex = 0
							i.state.selectedEvent = i.state.filteredEvents[0]
						}
					}
					return nil

				case "n", "N":
					// Also cancel with 'n' key
					i.state.currentState = BrowseStateEventDetail
					i.state.deleteError = nil
					return nil
				}
				return nil
			})
	}

	return nil
}

// removeEventFromList removes an event from both context.Events and filteredEvents.
func (i *BrowseTimelineIntent) removeEventFromList(eventID string) {
	// Remove from context.Events
	for idx, evt := range i.context.Events {
		if evt.ID == eventID {
			i.context.Events = append(i.context.Events[:idx], i.context.Events[idx+1:]...)
			break
		}
	}

	// Remove from filteredEvents
	for idx, evt := range i.state.filteredEvents {
		if evt.ID == eventID {
			i.state.filteredEvents = append(i.state.filteredEvents[:idx], i.state.filteredEvents[idx+1:]...)
			break
		}
	}

	// Update table rows
	i.updateTableRows()
}

// getContext returns a context for service calls.
func (i *BrowseTimelineIntent) getContext() context.Context {
	return context.Background()
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
			if filtered[a].Date.Equal(filtered[b].Date) {
				// Secondary sort by CreatedAt for same-date events
				if i.state.filters.SortOrder == "asc" {
					return filtered[a].CreatedAt.Before(filtered[b].CreatedAt)
				}
				return filtered[a].CreatedAt.After(filtered[b].CreatedAt)
			}
			if i.state.filters.SortOrder == "asc" {
				return filtered[a].Date.Before(filtered[b].Date)
			}
			return filtered[a].Date.After(filtered[b].Date)
		}
	})

	i.state.filteredEvents = filtered
}

// getStateName returns a human-readable name for the current state.
func (i *BrowseTimelineIntent) getStateName() string {
	switch i.state.currentState {
	case BrowseStateTimeline:
		return "Timeline"
	case BrowseStateEventDetail:
		return "Event Detail"
	case BrowseStateDeleteConfirm:
		return "Delete Event"
	default:
		return string(i.state.currentState)
	}
}

// getStateContent returns the content for the current state.
func (i *BrowseTimelineIntent) getStateContent() string {
	switch i.state.currentState {
	case BrowseStateTimeline:
		return i.viewTimeline()
	case BrowseStateEventDetail:
		return i.viewEventDetail()
	case BrowseStateDeleteConfirm:
		return i.viewDeleteConfirm()
	default:
		return ""
	}
}

// getContextHelp returns context-aware help text for the current state.
func (i *BrowseTimelineIntent) getContextHelp() string {
	theme := i.Theme()

	switch i.state.currentState {
	case BrowseStateTimeline:
		return CombineThemedFooters(
			ThemedCustomFooter(theme,
				components.NavigateBadge(),
				components.SelectBadge(),
				components.FilterBadge(),
				components.NewKeyBadge("Enter", "View Details"),
			),
			ThemedGlobalBadges(theme),
		)
	case BrowseStateEventDetail:
		// Add edit and delete badges if service is available
		badges := []components.KeyBadge{
			components.NewKeyBadge("Enter", "Select"),
		}
		if i.context.CLIEventService != nil {
			badges = append(badges,
				components.NewKeyBadge("e", "Edit"),
				components.NewKeyBadge("d", "Delete"),
			)
		}
		return CombineThemedFooters(
			ThemedCustomFooter(theme, badges...),
			ThemedGlobalBadges(theme),
		)
	case BrowseStateDeleteConfirm:
		return CombineThemedFooters(
			ThemedCustomFooter(theme,
				components.NewKeyBadge("y", "Confirm Delete"),
				components.NewKeyBadge("n/Esc", "Cancel"),
			),
			ThemedGlobalBadges(theme),
		)
	default:
		return ThemedGlobalBadges(theme)
	}
}

// View renders the intent's current state using StandardView.
func (i *BrowseTimelineIntent) View() string {
	// Create standard view with breadcrumbs
	view := i.CreateViewWithBreadcrumbs("Main Menu", "Browse Timeline", i.getStateName())

	// Get content for current state
	content := i.getStateContent()
	view.WithContent(content)

	// Get context-aware help
	help := i.getContextHelp()
	view.WithHelp(help).WithFooterSeparator(true)

	return view.Render()
}

// viewTimeline renders the timeline view with all events as a table.
func (i *BrowseTimelineIntent) viewTimeline() string {
	if len(i.state.filteredEvents) == 0 {
		i.listContainer.SetEmptyStateMessage("No events found.")
		return i.listContainer.Render()
	}

	// Ensure table rows are synchronized with current state
	i.updateTableRows()

	// Build pagination info with page number indicator
	pageSize := 15
	totalItems := len(i.state.filteredEvents)
	currentPage := (i.state.selectedIndex / pageSize) + 1
	totalPages := (totalItems + pageSize - 1) / pageSize
	paginationInfo := fmt.Sprintf("Events: %d | Page %d of %d", totalItems, currentPage, totalPages)
	i.listContainer.SetPaginationInfo(paginationInfo)

	return i.listContainer.Render()
}

// viewEventDetail renders the event detail view using the reusable component.
func (i *BrowseTimelineIntent) viewEventDetail() string {
	return components.RenderEventDetailCard(i.state.selectedEvent, i.Theme())
}

// viewDeleteConfirm renders the delete confirmation dialog.
func (i *BrowseTimelineIntent) viewDeleteConfirm() string {
	if i.state.selectedEvent == nil {
		return "No event selected."
	}

	var content strings.Builder

	// Warning header
	warningStyle := lipgloss.NewStyle().
		Foreground(lipgloss.Color("#F38BA8")). // Catppuccin Red
		Bold(true)

	content.WriteString(warningStyle.Render("⚠️  Delete Event?"))
	content.WriteString("\n\n")

	// Event summary
	content.WriteString(fmt.Sprintf("Date: %s\n", i.state.selectedEvent.Date.Format("2006-01-02")))
	if i.state.selectedEvent.Company != "" {
		content.WriteString(fmt.Sprintf("Company: %s\n", i.state.selectedEvent.Company))
	}

	// Truncate text for display
	text := i.state.selectedEvent.Text
	if len(text) > 100 {
		text = text[:100] + "..."
	}
	content.WriteString(fmt.Sprintf("\nText: %s\n", text))

	// Error message if delete failed
	if i.state.deleteError != nil {
		errorStyle := lipgloss.NewStyle().
			Foreground(lipgloss.Color("#F38BA8")).
			MarginTop(1)
		content.WriteString("\n")
		content.WriteString(errorStyle.Render(fmt.Sprintf("Error: %s", i.state.deleteError.Error())))
	}

	// Confirmation prompt
	content.WriteString("\n\nThis action cannot be undone.\n")

	// Apply card styling
	var cardStyle lipgloss.Style
	if theme := i.Theme(); theme != nil {
		cardStyle = theme.Styles().CardBase.
			BorderForeground(lipgloss.Color("#F38BA8")) // Red border for warning
	} else {
		cardStyle = lipgloss.NewStyle().
			Padding(1, 2).
			BorderStyle(lipgloss.RoundedBorder()).
			BorderForeground(lipgloss.Color("#F38BA8"))
	}

	return cardStyle.Render(content.String())
}

// Result returns the final result of the intent.
func (i *BrowseTimelineIntent) Result() *IntentResult[interface{}] {
	if i.result == nil {
		return nil
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

// ListNavigator interface implementation

// GetTotalItems returns the total number of filtered events.
func (i *BrowseTimelineIntent) GetTotalItems() int {
	return len(i.state.filteredEvents)
}

// GetSelectedIndex returns the current selection index.
func (i *BrowseTimelineIntent) GetSelectedIndex() int {
	return i.state.selectedIndex
}

// SetSelectedIndex sets the selection index and updates the display.
// This is the single source of truth for selection state.
func (i *BrowseTimelineIntent) SetSelectedIndex(idx int) {
	// Validate and set index
	if idx < 0 {
		idx = 0
	}
	if idx >= len(i.state.filteredEvents) {
		idx = len(i.state.filteredEvents) - 1
	}
	if idx < 0 {
		idx = 0 // Handle empty list
	}

	i.state.selectedIndex = idx

	// Update selected event
	if idx >= 0 && idx < len(i.state.filteredEvents) {
		i.state.selectedEvent = i.state.filteredEvents[idx]
	}

	// Update table display
	i.updateTableRows()
}

// GetPageSize returns the page size for pagination.
func (i *BrowseTimelineIntent) GetPageSize() int {
	return 15
}
