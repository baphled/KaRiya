package intents

import (
	"context"
	"fmt"
	"sort"

	"github.com/baphled/kariya/internal/cli/components"
	"github.com/baphled/kariya/internal/cli/screens"
	"github.com/baphled/kariya/internal/cli/screens/timeline"
	"github.com/baphled/kariya/internal/domain/career"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/rmhubbert/bubbletea-overlay"
)

// Custom message types for BrowseTimeline state transitions.

// EventSelectedMsg indicates the user selected an event.
type EventSelectedMsg struct {
	Event *career.CareerEvent
	Index int
}

// staticViewModel is a simple tea.Model that just returns static content.
// Used as background for bubbletea-overlay compositing.
type staticViewModel struct {
	content string
}

func (m *staticViewModel) Init() tea.Cmd {
	return nil
}

func (m *staticViewModel) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	return m, nil
}

func (m *staticViewModel) View() string {
	return m.content
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

// RequestAddEventMsg requests that the app route to CaptureEvent intent for adding a new event.
// This is sent to the app router which will activate CaptureEvent for a new event.
type RequestAddEventMsg struct{}

// EventDeletedMsg notifies that an event was successfully deleted.
// This is sent back to the intent after a delete operation completes.
type EventDeletedMsg struct {
	EventID string
}

// Ensure BrowseTimelineIntent implements ScreenResultHandler interface
var _ ScreenResultHandler = (*BrowseTimelineIntent)(nil)

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

	// active indicates whether this intent is currently active.
	active bool

	// result is the final result of the intent (set when complete).
	result *IntentResult[*BrowseTimelineResult]

	// --- Screen Orchestration (Phase 4.2 - Complete) ---
	// activeScreen holds the current screen being displayed.
	activeScreen screens.Screen

	// filterModal holds the filter modal (shown over the list)
	filterModal *components.FilterModalModel

	// searchModal holds the search modal for text search
	searchModal *components.EventSearchModal

	// sortModal holds the sort modal for sorting events
	sortModal *components.EventSortModal

	// deleteModal holds the delete confirmation modal (shown over the list)
	deleteModal *components.DeleteConfirmModal

	// quickAddModal holds the quick add event modal (shown over the list) - Phase 4 UX Issue 3A
	quickAddModal *components.QuickAddEventModal

	// editModal holds the edit event modal (shown over the list) - Phase 4 UX Issue 3B
	editModal *components.EditEventModal

	// viewDetailModal holds the event detail viewer modal (shown over the list)
	viewDetailModal *components.ViewEventDetailModal
}

// NewBrowseTimelineIntent creates a new BrowseTimeline intent.
func NewBrowseTimelineIntent(context *BrowseTimelineContext) (*BrowseTimelineIntent, error) {
	// Validate the context.
	if err := context.Validate(); err != nil {
		return nil, err
	}

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
			filterStack:    NewFilterStack(),
			selectedFacts:  make([]*career.Fact, 0),
			viewedEvents:   make([]*career.CareerEvent, 0),
		},
		active: true,
	}

	return intent, nil
}

// Init is called when the intent is activated.
func (i *BrowseTimelineIntent) Init() tea.Cmd {
	// Apply initial filters and sorting to the provided events.
	i.applyFilters()

	if len(i.state.filteredEvents) > 0 {
		i.state.selectedEvent = i.state.filteredEvents[0]
	}

	// Initialize with timeline list screen
	i.state.currentState = BrowseStateTimeline
	i.transitionToScreen(timeline.NewTimelineEventListScreen(i.state.filteredEvents))
	return nil
}

// Update processes a message in the intent.
func (i *BrowseTimelineIntent) Update(msg tea.Msg) tea.Cmd {
	if !i.active {
		return nil
	}

	// PATTERN 4: Modal Priority Handling
	// Modals are checked in priority order: search → filter → sort → other modals
	// Only ONE modal can be active at a time

	// If search modal is visible, handle it first
	if i.searchModal != nil && i.searchModal.IsVisible() {
		cmd, applied, searchData := i.searchModal.Update(msg)
		if applied && searchData != nil {
			// Search was applied - update filters and refresh list
			i.state.filters.SearchText = searchData.SearchText
			// Track search filter in stack for FIFO clearing
			if searchData.SearchText != "" {
				i.state.filterStack.Push(FilterLayerSearch)
			}
			i.applyFilters()
			i.transitionToScreen(timeline.NewTimelineEventListScreen(i.state.filteredEvents))
		}
		return cmd
	}

	// If filter modal is visible, handle it next
	if i.filterModal != nil && i.filterModal.IsVisible() {
		cmd, applied, filterData := i.filterModal.Update(msg)
		if applied && filterData != nil {
			// Filters were applied - update and refresh list
			i.state.filters.Companies = filterData.Companies
			i.state.filters.Categories = filterData.Categories
			i.state.filters.Projects = filterData.Projects
			i.state.filters.SortBy = filterData.SortBy
			i.state.filters.SortOrder = filterData.SortOrder

			// Track filters in stack for FIFO clearing
			if len(filterData.Companies) > 0 {
				i.state.filterStack.Push(FilterLayerCompany)
			}
			if len(filterData.Categories) > 0 {
				i.state.filterStack.Push(FilterLayerCategory)
			}
			if len(filterData.Projects) > 0 {
				i.state.filterStack.Push(FilterLayerProject)
			}

			i.applyFilters()
			i.transitionToScreen(timeline.NewTimelineEventListScreen(i.state.filteredEvents))
		}
		return cmd
	}

	// If sort modal is visible, handle it next
	if i.sortModal != nil && i.sortModal.IsVisible() {
		cmd, applied, sortData := i.sortModal.Update(msg)
		if applied && sortData != nil {
			// Sort was applied - update filters and refresh list
			i.state.filters.SortBy = sortData.SortBy
			i.state.filters.SortOrder = sortData.SortOrder

			// Track sort in stack if it's not default
			if sortData.SortBy != "date" || sortData.SortOrder != "desc" {
				i.state.filterStack.Push(FilterLayerSort)
			}

			i.applyFilters()
			i.transitionToScreen(timeline.NewTimelineEventListScreen(i.state.filteredEvents))
		}
		return cmd
	}

	// If quick add modal is visible, handle it next (Phase 4 UX Issue 3A)
	if i.quickAddModal != nil && i.quickAddModal.IsVisible() {
		cmd, completed, eventData := i.quickAddModal.Update(msg)
		if !i.quickAddModal.IsVisible() {
			// Modal closed
			if completed && eventData != nil {
				// User completed form - save the new event
				ctx := context.Background()
				newEvent := eventData.ToCareerEvent()

				// Save using CaptureEvent API (text, date, mode, options)
				captureErr := i.context.CLIEventService.CaptureEvent(
					ctx,
					newEvent.Text,
					newEvent.Date,
					"manual", // EventCaptureMode.Manual
				)
				if captureErr != nil {
					// TODO: Show error modal (Phase 4 Issue 3)
					return cmd
				}

				// Reload events list to get the newly created event with ID
				refreshedEvents, listErr := i.context.CLIEventService.ListEvents(ctx, nil)
				if listErr != nil {
					// TODO: Show error modal
					return cmd
				}
				i.context.Events = refreshedEvents
				i.applyFilters()
				i.transitionToScreen(timeline.NewTimelineEventListScreen(i.state.filteredEvents))
			}
			// Clear the modal
			i.quickAddModal = nil
		}
		return cmd
	}

	// If edit modal is visible, handle it next (Phase 4 UX Issue 3B)
	if i.editModal != nil && i.editModal.IsVisible() {
		cmd, completed, eventData := i.editModal.Update(msg)
		if !i.editModal.IsVisible() {
			// Modal closed
			if completed && eventData != nil {
				// User completed form - update the event
				ctx := context.Background()
				originalEvent := i.editModal.GetOriginalEvent()
				updatedEvent := eventData.ToCareerEvent(
					originalEvent.ID,
					originalEvent.CreatedAt,
					originalEvent.UpdatedAt,
				)
				// Use UpdateEventMetadata since we have the full event object
				if err := i.context.CLIEventService.UpdateEventMetadata(ctx, updatedEvent); err != nil {
					// TODO: Show error modal (Phase 4 Issue 3)
					return cmd
				}
				// Update the event in the list
				for idx, evt := range i.context.Events {
					if evt.ID == updatedEvent.ID {
						i.context.Events[idx] = updatedEvent
						break
					}
				}
				i.applyFilters()
				i.transitionToScreen(timeline.NewTimelineEventListScreen(i.state.filteredEvents))
			}
			// Clear the modal
			i.editModal = nil
		}
		return cmd
	}

	// If delete modal is visible, handle it next (Phase 4 UX Issue 2)
	if i.deleteModal != nil && i.deleteModal.IsVisible() {
		cmd, confirmed := i.deleteModal.Update(msg)
		if !i.deleteModal.IsVisible() {
			// Modal closed
			if confirmed && i.state.selectedEvent != nil {
				// User confirmed - delete the event
				ctx := context.Background()
				if err := i.context.CLIEventService.DeleteEvent(ctx, i.state.selectedEvent.ID); err != nil {
					// TODO: Show error modal (Phase 4 Issue 3)
					i.state.deleteError = err
					return cmd
				}
				// Remove deleted event from the list
				deletedID := i.state.selectedEvent.ID
				newEvents := make([]*career.CareerEvent, 0, len(i.context.Events)-1)
				for _, evt := range i.context.Events {
					if evt.ID != deletedID {
						newEvents = append(newEvents, evt)
					}
				}
				i.context.Events = newEvents
				i.applyFilters()
				i.transitionToScreen(timeline.NewTimelineEventListScreen(i.state.filteredEvents))
			}
			// Clear selected event
			i.state.selectedEvent = nil
		}
		return cmd
	}

	// If view detail modal is visible, handle it next
	if i.viewDetailModal != nil && i.viewDetailModal.IsVisible() {
		_, cmd := i.viewDetailModal.Update(msg)
		if !i.viewDetailModal.IsVisible() {
			// Modal closed - clear the modal reference
			i.viewDetailModal = nil
		}
		return cmd
	}

	// Handle global keys BEFORE delegating to screen
	// This ensures q (quit), ? (help), etc. are always processed first
	if keyMsg, ok := msg.(tea.KeyMsg); ok {
		switch HandleGlobalKeys(keyMsg) {
		case KeyQuit:
			// 'q' pressed - quit the application
			return tea.Quit
		case KeyHelp:
			// '?' pressed - toggle help modal
			i.helpModal.Toggle()
			return nil
		}
		// KeyBack (esc) is handled by the screen as CancelResult

		// Handle intent-specific shortcuts (f, s, /, x) on timeline list screen
		if i.state.currentState == BrowseStateTimeline {
			switch keyMsg.String() {
			case "f":
				// Open filter modal
				return i.openFilterModal()
			case "s":
				// Open sort modal
				return i.openSortModal()
			case "/":
				// Open search modal
				return i.openSearchModal()
			case "x":
				// Clear filters (FIFO order - most recent first)
				if i.HasActiveFilters() {
					i.ClearFilters()
					return i.RefreshData()
				}
			}
		}
	}

	// Delegate to active screen
	if i.activeScreen != nil {
		cmd, result := i.activeScreen.Update(msg)
		if result != nil {
			screenCmd := i.handleScreenResult(result)
			if screenCmd != nil {
				return tea.Batch(cmd, screenCmd)
			}
		}
		return cmd
	}

	return nil
}

// View renders the intent's current state using StandardView.
func (i *BrowseTimelineIntent) View() string {
	if i.activeScreen == nil {
		return "No active screen"
	}

	// Type-assert to screens with RenderContent method
	switch screen := i.activeScreen.(type) {
	case *timeline.TimelineEventListScreen:
		// Create StandardView with table content
		view := i.CreateViewWithBreadcrumbs("Main Menu", "Browse Timeline", i.getStateName())
		view.WithContent(screen.RenderContent())
		view.WithHelp(i.getContextHelp())

		// Render the complete view FIRST
		baseView := view.Render()

		// PATTERN 1: Modal Overlay Rendering
		// Check modals in priority order (only ONE can be active at a time)

		// If search modal is visible, overlay it on the COMPLETE rendered view
		if i.searchModal != nil && i.searchModal.IsVisible() {
			return i.renderSearchModalOverlay(baseView)
		}

		// If filter modal is visible, overlay it on the COMPLETE rendered view
		if i.filterModal != nil && i.filterModal.IsVisible() {
			return i.renderFilterModalOverlay(baseView)
		}

		// If sort modal is visible, overlay it on the COMPLETE rendered view
		if i.sortModal != nil && i.sortModal.IsVisible() {
			return i.renderSortModalOverlay(baseView)
		}

		// If quick add modal is visible, overlay it on the COMPLETE rendered view
		if i.quickAddModal != nil && i.quickAddModal.IsVisible() {
			return i.renderQuickAddModalOverlay(baseView)
		}

		// If edit modal is visible, overlay it on the COMPLETE rendered view
		if i.editModal != nil && i.editModal.IsVisible() {
			return i.renderEditModalOverlay(baseView)
		}

		// If delete modal is visible, overlay it on the COMPLETE rendered view
		if i.deleteModal != nil && i.deleteModal.IsVisible() {
			return i.renderDeleteModalOverlay(baseView)
		}

		// If view detail modal is visible, overlay it on the COMPLETE rendered view
		if i.viewDetailModal != nil && i.viewDetailModal.IsVisible() {
			return i.renderViewDetailModalOverlay(baseView)
		}

		return baseView

	case *timeline.TimelineEventDetailScreen:
		// LEGACY: Event detail is now shown as a modal, not a full screen
		// This case is kept for backward compatibility but should not be reached
		// Event detail: Use RenderContent and add themed footer
		view := i.CreateViewWithBreadcrumbs("Main Menu", "Browse Timeline", i.getStateName())
		view.WithContent(screen.RenderContent())
		view.WithHelp(i.getContextHelp())
		return view.Render()

	case *timeline.EventDeleteConfirmScreen:
		// Delete confirmation: Use full View() (has its own footer)
		return screen.View()

	default:
		// Fallback: use full View() for unknown screens
		return i.activeScreen.View()
	}
}

// Result returns the final result of the intent.
func (i *BrowseTimelineIntent) Result() *IntentResult[interface{}] {
	if i.result == nil {
		return nil
	}

	// Type-erase the result for the Intent interface
	return &IntentResult[interface{}]{
		Status:   i.result.Status,
		Data:     i.result.Data,
		Error:    i.result.Error,
		Metadata: i.result.Metadata,
	}
}

// setCancelled marks the intent as cancelled by the user.
func (i *BrowseTimelineIntent) setCancelled() {
	i.result = &IntentResult[*BrowseTimelineResult]{
		Status: Cancelled,
	}
	i.active = false
}

// getContext returns a context for service calls.
func (i *BrowseTimelineIntent) getContext() context.Context {
	return context.Background()
}

// applyFilters filters the events based on current filter state.
func (i *BrowseTimelineIntent) applyFilters() {
	filtered := make([]*career.CareerEvent, 0)

	for _, evt := range i.context.Events {
		// Apply text search (if specified)
		// Search across text, company, categories, and project fields
		if i.state.filters.SearchText != "" {
			categoriesStr := ""
			if len(evt.Categories) > 0 {
				categoriesStr = evt.Categories[0] // Primary category for search
			}
			if !SearchableFields(i.state.filters.SearchText, evt.Text, evt.Company, categoriesStr, evt.Project) {
				continue // Skip events that don't match search
			}
		}

		// Apply tag filters
		if len(i.state.filters.Tags) > 0 {
			hasTag := false
			for _, filterTag := range i.state.filters.Tags {
				for _, evtTag := range evt.Tags {
					if evtTag == filterTag {
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

		// Apply company filter
		if len(i.state.filters.Companies) > 0 {
			hasCompany := false
			for _, filterCompany := range i.state.filters.Companies {
				if evt.Company == filterCompany {
					hasCompany = true
					break
				}
			}
			if !hasCompany {
				continue
			}
		}

		// Apply category filter
		if len(i.state.filters.Categories) > 0 {
			hasCategory := false
			for _, filterCat := range i.state.filters.Categories {
				for _, evtCat := range evt.Categories {
					if evtCat == filterCat {
						hasCategory = true
						break
					}
				}
				if hasCategory {
					break
				}
			}
			if !hasCategory {
				continue
			}
		}

		// Apply project filter
		if len(i.state.filters.Projects) > 0 {
			hasProject := false
			for _, filterProject := range i.state.filters.Projects {
				if evt.Project == filterProject {
					hasProject = true
					break
				}
			}
			if !hasProject {
				continue
			}
		}

		// Event passes all filters
		filtered = append(filtered, evt)
	}

	// Apply sorting
	switch i.state.filters.SortBy {
	case "date":
		sort.Slice(filtered, func(a, b int) bool {
			if i.state.filters.SortOrder == "asc" {
				return filtered[a].Date.Before(filtered[b].Date)
			}
			return filtered[a].Date.After(filtered[b].Date)
		})
	case "text":
		sort.Slice(filtered, func(a, b int) bool {
			if i.state.filters.SortOrder == "asc" {
				return filtered[a].Text < filtered[b].Text
			}
			return filtered[a].Text > filtered[b].Text
		})
	}

	i.state.filteredEvents = filtered
}

// ============================================================================
// FilterBehavior Implementation (Source of Truth: ManageSkills)
// ============================================================================

// HasActiveFilters returns true if any non-default filters are active.
// Implements FilterBehavior interface.
func (i *BrowseTimelineIntent) HasActiveFilters() bool {
	f := i.state.filters
	if f == nil {
		return false
	}

	// Check all filter types
	return f.SearchText != "" ||
		len(f.Tags) > 0 ||
		len(f.Companies) > 0 ||
		len(f.Categories) > 0 ||
		len(f.Projects) > 0 ||
		f.DateFrom != "" ||
		f.DateTo != "" ||
		(f.SortBy != "" && f.SortBy != "date") || // date is default
		(f.SortOrder != "" && f.SortOrder != "desc") // desc is default
}

// ClearFilters clears filters in FIFO order (most recent first).
// Pressing 'x' multiple times progressively removes filters.
// Implements FilterBehavior interface.
func (i *BrowseTimelineIntent) ClearFilters() {
	// If no filter stack or empty stack, clear everything
	if i.state.filterStack == nil || i.state.filterStack.IsEmpty() {
		i.clearAllFilters()
		return
	}

	// Pop most recent filter layer
	layer := i.state.filterStack.Pop()

	// Clear the specific filter layer
	switch layer {
	case FilterLayerSearch:
		i.state.filters.SearchText = ""

	case FilterLayerCompany:
		i.state.filters.Companies = []string{}

	case FilterLayerCategory:
		i.state.filters.Categories = []string{}

	case FilterLayerProject:
		i.state.filters.Projects = []string{}

	case FilterLayerTags:
		i.state.filters.Tags = []string{}

	case FilterLayerSort:
		// Reset to default sort
		i.state.filters.SortBy = "date"
		i.state.filters.SortOrder = "desc"
	}

	// If no more filters active, clear the entire stack
	if !i.HasActiveFilters() {
		i.state.filterStack.Clear()
	}
}

// clearAllFilters resets all filters to default state.
func (i *BrowseTimelineIntent) clearAllFilters() {
	i.state.filters = &TimelineFilters{
		SearchText: "",
		Tags:       []string{},
		Companies:  []string{},
		Categories: []string{},
		Projects:   []string{},
		DateFrom:   "",
		DateTo:     "",
		SortBy:     "date", // Default sort
		SortOrder:  "desc", // Default order (newest first)
	}
	if i.state.filterStack != nil {
		i.state.filterStack.Clear()
	}
}

// ApplyFilters applies current filter state (alias for applyFilters).
// Implements FilterBehavior interface.
func (i *BrowseTimelineIntent) ApplyFilters() {
	i.applyFilters()
}

// RefreshData reloads/refreshes the filtered data.
// For BrowseTimeline, this re-applies filters to the current event list.
// Implements FilterBehavior interface.
func (i *BrowseTimelineIntent) RefreshData() tea.Cmd {
	i.applyFilters()
	i.transitionToScreen(timeline.NewTimelineEventListScreen(i.state.filteredEvents))
	return nil
}

// removeEventFromList removes an event from both the full list and filtered list.
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

	// Update active screen if it's a list screen
	i.transitionToScreen(timeline.NewTimelineEventListScreen(i.state.filteredEvents))
}

// ============================================================================
// Screen Orchestration (Phase 4.2 - Complete)
// ============================================================================

// EnableScreens is a no-op for backward compatibility.
// Screens are now the default and only architecture.
func (i *BrowseTimelineIntent) EnableScreens() {
	// No-op: screens are always enabled
}

// transitionToScreen sets the active screen and updates state.
func (i *BrowseTimelineIntent) transitionToScreen(screen screens.Screen) {
	i.activeScreen = screen

	// Set terminal dimensions
	termInfo := i.GetTerminalInfo()
	if termInfo != nil && termInfo.Width > 0 && termInfo.Height > 0 {
		screen.SetTerminalInfo(termInfo.Width, termInfo.Height)
	}

	// Set theme if available
	if theme := i.Theme(); theme != nil {
		screen.SetTheme(theme)
	}

	// Set logo if available
	if logo := i.GetLogo(); logo != nil {
		screen.SetLogo(logo, i.GetLogoSpacing())
	}
}

// handleScreenResult processes a screen result and determines next action.
//
// Uses ScreenResultDispatcher pattern to eliminate repetitive type switching.
// BrowseTimelineIntent implements ScreenResultHandler interface for compile-time safety.
func (i *BrowseTimelineIntent) handleScreenResult(result interface{}) tea.Cmd {
	// Handle nil result
	if result == nil {
		return nil
	}

	// Cast to ScreenResult (BrowseTimeline accepts interface{} for backward compatibility)
	screenResult, ok := result.(screens.ScreenResult)
	if !ok {
		return nil
	}

	return NewScreenResultDispatcher(i).Dispatch(screenResult)
}

// HandleCancel handles screen cancellation (back/escape).
//
// Implements ScreenResultHandler interface.
func (i *BrowseTimelineIntent) HandleCancel(result *screens.CancelResult) tea.Cmd {
	switch i.state.currentState {
	case BrowseStateTimeline:
		// At root state, cancel means exit intent
		i.setCancelled()
		return nil

	case BrowseStateEventDetail:
		// LEGACY: Event detail is now a modal, not a separate state
		// This case is kept for backward compatibility but should not be reached
		// Return to timeline list
		i.state.currentState = BrowseStateTimeline
		i.transitionToScreen(timeline.NewTimelineEventListScreen(i.state.filteredEvents))
		return nil

	case BrowseStateDeleteConfirm:
		// Escape from delete confirmation - return to timeline list
		i.state.currentState = BrowseStateTimeline
		i.transitionToScreen(timeline.NewTimelineEventListScreen(i.state.filteredEvents))
		return nil

	default:
		// Fallback: cancel intent
		i.setCancelled()
		return nil
	}
}

// HandleNavigate handles screen navigation results.
//
// Implements ScreenResultHandler interface.
func (i *BrowseTimelineIntent) HandleNavigate(result *screens.NavigateResult) tea.Cmd {
	// Check if it's an action (map) or event selection
	if actionData, ok := result.ResultData.(map[string]interface{}); ok {
		action, _ := actionData["action"].(string)
		switch action {
		case "add":
			// Show quick add modal (Phase 4 UX Issue 3A)
			termInfo := i.GetTerminalInfo()
			width := 120
			height := 40
			if termInfo != nil {
				width = termInfo.Width
				height = termInfo.Height
			}
			i.quickAddModal = components.NewQuickAddEventModal(width, height)
			return i.quickAddModal.Init()
		case "edit":
			// Get the event from the action data
			if event, ok := actionData["event"].(*career.CareerEvent); ok {
				// Show edit modal (Phase 4 UX Issue 3B)
				termInfo := i.GetTerminalInfo()
				width := 120
				height := 40
				if termInfo != nil {
					width = termInfo.Width
					height = termInfo.Height
				}
				i.editModal = components.NewEditEventModal(event, width, height)
				return i.editModal.Init()
			}
			return nil
		case "delete":
			// Get the event from the action data
			if event, ok := actionData["event"].(*career.CareerEvent); ok {
				// Show delete confirmation modal (Phase 4 UX Issue 2)
				// Modal overlay instead of screen transition - preserves context
				eventText := event.Text
				if len(eventText) > 50 {
					eventText = eventText[:47] + "..."
				}
				i.deleteModal = components.NewDeleteConfirmModal(
					event.Text,
					"Delete Event",
					fmt.Sprintf("Are you sure you want to delete '%s'?", eventText),
				)
				// Store event for deletion if confirmed
				i.state.selectedEvent = event
				return i.deleteModal.Init()
			}
			return nil
		case "filter":
			// Show filter modal over the current list
			termInfo := i.GetTerminalInfo()
			width := 120
			height := 40
			if termInfo != nil {
				width = termInfo.Width
				height = termInfo.Height
			}
			// Convert TimelineFilters to components.TimelineFilters
			currentFilters := &components.TimelineFilters{
				SearchText: i.state.filters.SearchText,
				Tags:       i.state.filters.Tags,
				Companies:  i.state.filters.Companies,
				Categories: i.state.filters.Categories,
				Projects:   i.state.filters.Projects,
				SortBy:     i.state.filters.SortBy,
				SortOrder:  i.state.filters.SortOrder,
			}
			i.filterModal = components.NewFilterModal(
				i.context.Events,
				currentFilters,
				width,
				height,
			)
			// Initialize the modal's form to start its lifecycle
			return i.filterModal.Init()
		default:
			return nil
		}
	}

	// Check if it's a delete confirmation (bool result)
	if confirmed, ok := result.ResultData.(bool); ok {
		return i.handleDeleteConfirmation(confirmed)
	}

	// Event selection - show detail modal instead of transitioning to screen
	if event, ok := result.ResultData.(*career.CareerEvent); ok {
		i.state.selectedEvent = event
		i.state.viewedEvents = append(i.state.viewedEvents, event)
		// Show view detail modal (replaces full-screen detail view)
		termInfo := i.GetTerminalInfo()
		width := 120
		height := 40
		if termInfo != nil {
			width = termInfo.Width
			height = termInfo.Height
		}
		theme := i.Theme()
		i.viewDetailModal = components.NewViewEventDetailModal(event, theme)
		i.viewDetailModal.SetDimensions(width, height)
		i.viewDetailModal.Show()
		return nil
	}

	return nil
}

// handleDeleteConfirmation handles the result of the delete confirmation dialog.
func (i *BrowseTimelineIntent) handleDeleteConfirmation(confirmed bool) tea.Cmd {
	if !confirmed {
		// User cancelled - return to timeline list
		i.state.currentState = BrowseStateTimeline
		i.transitionToScreen(timeline.NewTimelineEventListScreen(i.state.filteredEvents))
		return nil
	}

	// User confirmed deletion - delete the event
	if i.state.selectedEvent == nil {
		// No event selected, return to timeline
		i.state.currentState = BrowseStateTimeline
		i.transitionToScreen(timeline.NewTimelineEventListScreen(i.state.filteredEvents))
		return nil
	}

	// Delete the event using the CLI service
	if i.context.CLIEventService != nil {
		eventID := i.state.selectedEvent.ID
		err := i.context.CLIEventService.DeleteEvent(i.getContext(), eventID)
		if err != nil {
			// Store error and return to timeline
			i.state.deleteError = err
			i.state.currentState = BrowseStateTimeline
			i.transitionToScreen(timeline.NewTimelineEventListScreen(i.state.filteredEvents))
			return nil
		}

		// Remove event from lists
		i.removeEventFromList(eventID)
	}

	// Return to timeline list
	i.state.currentState = BrowseStateTimeline
	i.transitionToScreen(timeline.NewTimelineEventListScreen(i.state.filteredEvents))
	return nil
}

// handleSubmitResult handles form submissions (not used in timeline).
// HandleSubmit handles form submission results.
//
// Implements ScreenResultHandler interface.
func (i *BrowseTimelineIntent) HandleSubmit(result *screens.SubmitResult) tea.Cmd {
	// Timeline doesn't have forms, but include for completeness
	return nil
}

// HandleError handles error results from screens.
//
// Implements ScreenResultHandler interface.
func (i *BrowseTimelineIntent) HandleError(result *screens.ErrorResult) tea.Cmd {
	// Store error and return to previous state
	i.state.deleteError = result.Err
	return nil
}

// getStateName returns a human-readable name for the current state.
func (i *BrowseTimelineIntent) getStateName() string {
	switch i.state.currentState {
	case BrowseStateTimeline:
		return "Timeline"
	case BrowseStateEventDetail:
		return "Event Details"
	default:
		return "Unknown"
	}
}

// renderFilterModalOverlay renders the filter modal overlay using bubbletea-overlay.
// The modal is automatically positioned and composited onto the background.
func (i *BrowseTimelineIntent) renderFilterModalOverlay(background string) string {
	// Create a simple background model that just returns the rendered view
	bgModel := &staticViewModel{content: background}

	// Use bubbletea-overlay to composite the form onto the background
	// Position at Center/Center with a small upward offset to avoid footer
	overlayModel := overlay.New(
		i.filterModal,  // Foreground: the form modal
		bgModel,        // Background: the rendered timeline view
		overlay.Center, // X position
		overlay.Center, // Y position
		0,              // X offset
		-2,             // Y offset (move up 2 lines to avoid footer)
	)

	return overlayModel.View()
}

// renderQuickAddModalOverlay renders the quick add event modal using bubbletea-overlay.
func (i *BrowseTimelineIntent) renderQuickAddModalOverlay(background string) string {
	// Create a simple background model that just returns the rendered view
	bgModel := &staticViewModel{content: background}

	// Use bubbletea-overlay to composite the form onto the background
	overlayModel := overlay.New(
		i.quickAddModal, // Foreground: the form modal
		bgModel,         // Background: the rendered timeline view
		overlay.Center,  // X position
		overlay.Center,  // Y position
		0,               // X offset
		-2,              // Y offset (move up 2 lines to avoid footer)
	)

	return overlayModel.View()
}

// renderEditModalOverlay renders the edit event modal using bubbletea-overlay.
func (i *BrowseTimelineIntent) renderEditModalOverlay(background string) string {
	// Create a simple background model that just returns the rendered view
	bgModel := &staticViewModel{content: background}

	// Use bubbletea-overlay to composite the form onto the background
	overlayModel := overlay.New(
		i.editModal,    // Foreground: the form modal
		bgModel,        // Background: the rendered timeline view
		overlay.Center, // X position
		overlay.Center, // Y position
		0,              // X offset
		-2,             // Y offset (move up 2 lines to avoid footer)
	)

	return overlayModel.View()
}

// renderDeleteModalOverlay renders the delete confirmation modal using bubbletea-overlay.
func (i *BrowseTimelineIntent) renderDeleteModalOverlay(background string) string {
	// Create a simple background model that just returns the rendered view
	bgModel := &staticViewModel{content: background}

	// Use bubbletea-overlay to composite the delete modal onto the background
	overlayModel := overlay.New(
		i.deleteModal,  // Foreground: the delete confirmation modal
		bgModel,        // Background: the rendered timeline view
		overlay.Center, // X position
		overlay.Center, // Y position
		0,              // X offset
		-2,             // Y offset (move up 2 lines to avoid footer)
	)

	return overlayModel.View()
}

// renderSearchModalOverlay renders the search modal using bubbletea-overlay.
func (i *BrowseTimelineIntent) renderSearchModalOverlay(background string) string {
	// Use the modal's RenderOverlay method for consistent rendering
	return i.searchModal.RenderOverlay(background)
}

// renderSortModalOverlay renders the sort modal using bubbletea-overlay.
func (i *BrowseTimelineIntent) renderSortModalOverlay(background string) string {
	// Use the modal's RenderOverlay method for consistent rendering
	return i.sortModal.RenderOverlay(background)
}

// renderViewDetailModalOverlay renders the event detail modal using bubbletea-overlay.
func (i *BrowseTimelineIntent) renderViewDetailModalOverlay(background string) string {
	// Pre-render the modal content to a static string
	// This follows the same pattern as EventSearchModal.RenderOverlay
	modalContent := &staticViewModel{content: i.viewDetailModal.View()}
	bgModel := &staticViewModel{content: background}

	// Use bubbletea-overlay to composite the detail modal onto the background
	overlayModel := overlay.New(
		modalContent,   // Foreground: pre-rendered modal content
		bgModel,        // Background: the rendered timeline view
		overlay.Center, // X position
		overlay.Center, // Y position
		0,              // X offset
		-2,             // Y offset (move up 2 lines to avoid footer)
	)

	return overlayModel.View()
}

// openSearchModal creates and initializes the search modal.
func (i *BrowseTimelineIntent) openSearchModal() tea.Cmd {
	termInfo := i.GetTerminalInfo()
	width, height := 120, 40
	if termInfo != nil {
		width, height = termInfo.Width, termInfo.Height
	}

	i.searchModal = components.NewEventSearchModal(
		i.state.filters.SearchText,
		width,
		height,
	)
	return i.searchModal.Init()
}

// openSortModal creates and initializes the sort modal.
func (i *BrowseTimelineIntent) openSortModal() tea.Cmd {
	termInfo := i.GetTerminalInfo()
	width, height := 120, 40
	if termInfo != nil {
		width, height = termInfo.Width, termInfo.Height
	}

	current := &components.EventSortConfig{
		SortBy:    i.state.filters.SortBy,
		SortOrder: i.state.filters.SortOrder,
	}

	i.sortModal = components.NewEventSortModal(
		i.context.Events,
		current,
		width,
		height,
	)
	return i.sortModal.Init()
}

// openFilterModal creates and initializes the filter modal.
func (i *BrowseTimelineIntent) openFilterModal() tea.Cmd {
	termInfo := i.GetTerminalInfo()
	width, height := 120, 40
	if termInfo != nil {
		width, height = termInfo.Width, termInfo.Height
	}

	currentFilters := &components.TimelineFilters{
		SearchText: i.state.filters.SearchText,
		Tags:       i.state.filters.Tags,
		Companies:  i.state.filters.Companies,
		Categories: i.state.filters.Categories,
		Projects:   i.state.filters.Projects,
		SortBy:     i.state.filters.SortBy,
		SortOrder:  i.state.filters.SortOrder,
	}

	i.filterModal = components.NewFilterModal(
		i.context.Events,
		currentFilters,
		width,
		height,
	)
	return i.filterModal.Init()
}

// getContextHelp returns themed keyboard shortcuts for the current state.
// This follows the legacy pattern of using KeyBadge components for consistent styling.
func (i *BrowseTimelineIntent) getContextHelp() string {
	theme := i.Theme()

	switch i.state.currentState {
	case BrowseStateTimeline:
		// Timeline list footer: Navigate, View Details, Add, Edit, Delete, Search, Filter, Sort, Clear (conditional), Back + Global shortcuts
		badges := []components.KeyBadge{
			components.NavigateBadge(),                      // ↑/↓: Navigate
			components.NewKeyBadge("Enter", "View Details"), // Enter: View Details
			components.AddBadge(),                           // a: Add
			components.EditBadge(),                          // e: Edit
			components.DeleteBadge(),                        // d: Delete
			components.SearchBadge(),                        // /: Search
			components.FilterBadge(),                        // f: Filter
			components.NewKeyBadge("s", "Sort"),             // s: Sort
		}

		// Conditionally add "Clear filters" badge when filters are active
		if i.HasActiveFilters() {
			badges = append(badges, components.NewKeyBadge("x", "Clear filters"))
		}

		badges = append(badges, components.BackBadge()) // Esc: Back

		return CombineThemedFooters(
			ThemedCustomFooter(theme, badges...),
			ThemedGlobalBadges(theme), // q: Quit, m: Main Menu
		)
	case BrowseStateEventDetail:
		// LEGACY: Event detail is now a modal with its own footer
		// This case is kept for backward compatibility but should not be reached
		// Event detail footer: Edit, Delete, Back + Global shortcuts
		return CombineThemedFooters(
			ThemedCustomFooter(theme,
				components.EditBadge(),   // e: Edit
				components.DeleteBadge(), // d: Delete
				components.BackBadge(),   // Esc: Back
			),
			ThemedGlobalBadges(theme), // q: Quit, m: Main Menu
		)
	case BrowseStateDeleteConfirm:
		// Delete confirmation uses BaseConfirmScreen's built-in footer
		// y/n: Choose, Enter: Confirm, ←→/hl: Toggle, Esc: Cancel
		// No need to override it
		return ""
	default:
		// Fallback: just global shortcuts
		return ThemedGlobalBadges(theme)
	}
}
