package intents

import (
	"context"
	"sort"

	"github.com/baphled/kariya/internal/cli/screens"
	"github.com/baphled/kariya/internal/cli/screens/timeline"
	"github.com/baphled/kariya/internal/domain/career"
	tea "github.com/charmbracelet/bubbletea"
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

	// active indicates whether this intent is currently active.
	active bool

	// result is the final result of the intent (set when complete).
	result *IntentResult[*BrowseTimelineResult]

	// --- Screen Orchestration (Phase 4.2 - Complete) ---
	// activeScreen holds the current screen being displayed.
	activeScreen screens.Screen
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
	// Delegate to active screen
	if i.activeScreen != nil {
		return i.activeScreen.View()
	}

	// Fallback (should not happen)
	return "No active screen"
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

// setCompleted marks the intent as successfully completed.
func (i *BrowseTimelineIntent) setCompleted() {
	i.result = &IntentResult[*BrowseTimelineResult]{
		Status: Completed,
		Data: &BrowseTimelineResult{
			SelectedEvent: i.state.selectedEvent,
			FinalFilters:  i.state.filters,
			ViewedEvents:  i.state.viewedEvents,
			SelectedFacts: i.state.selectedFacts,
		},
	}
	i.active = false
}

// setCancelled marks the intent as cancelled by the user.
func (i *BrowseTimelineIntent) setCancelled() {
	i.result = &IntentResult[*BrowseTimelineResult]{
		Status: Cancelled,
	}
	i.active = false
}

// setFailed marks the intent as failed with an error.
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

// getContext returns a context for service calls.
func (i *BrowseTimelineIntent) getContext() context.Context {
	return context.Background()
}

// applyFilters filters the events based on current filter state.
func (i *BrowseTimelineIntent) applyFilters() {
	filtered := make([]*career.CareerEvent, 0)

	for _, evt := range i.context.Events {
		// Apply text search (if specified)
		if i.state.filters.SearchText != "" {
			// Simple case-insensitive contains match
			// TODO: More sophisticated search (for now, skip text search)
			// Note: When implemented, should check if event text contains search text
			// and continue if it doesn't match
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
func (i *BrowseTimelineIntent) handleScreenResult(result interface{}) tea.Cmd {
	if result == nil {
		return nil
	}

	switch r := result.(type) {
	case *screens.CancelResult:
		return i.handleCancelResult()
	case *screens.NavigateResult:
		return i.handleNavigateResult(r)
	case *screens.SubmitResult:
		return i.handleSubmitResult(r)
	case *screens.ErrorResult:
		return i.handleErrorResult(r)
	default:
		return nil
	}
}

// handleCancelResult handles screen cancellation (back/escape).
func (i *BrowseTimelineIntent) handleCancelResult() tea.Cmd {
	switch i.state.currentState {
	case BrowseStateTimeline:
		// At root state, cancel means exit intent
		i.setCancelled()
		return nil

	case BrowseStateEventDetail:
		// Return to timeline list
		i.state.currentState = BrowseStateTimeline
		i.transitionToScreen(timeline.NewTimelineEventListScreen(i.state.filteredEvents))
		return nil

	default:
		// Fallback: cancel intent
		i.setCancelled()
		return nil
	}
}

// handleNavigateResult handles screen navigation results.
func (i *BrowseTimelineIntent) handleNavigateResult(result *screens.NavigateResult) tea.Cmd {
	// Check if it's an action (map) or event selection
	if actionData, ok := result.ResultData.(map[string]interface{}); ok {
		action, _ := actionData["action"].(string)
		switch action {
		case "add":
			// TODO: Handle add event action (route to CaptureEvent intent)
			return nil
		case "edit":
			// TODO: Handle edit event action
			return nil
		case "delete":
			// TODO: Handle delete event action
			return nil
		default:
			return nil
		}
	}

	// Event selection - navigate to detail view
	if event, ok := result.ResultData.(*career.CareerEvent); ok {
		i.state.selectedEvent = event
		i.state.viewedEvents = append(i.state.viewedEvents, event)
		i.state.currentState = BrowseStateEventDetail
		i.transitionToScreen(timeline.NewTimelineEventDetailScreen(event))
		return nil
	}

	return nil
}

// handleSubmitResult handles form submissions (not used in timeline).
func (i *BrowseTimelineIntent) handleSubmitResult(result *screens.SubmitResult) tea.Cmd {
	// Timeline doesn't have forms, but include for completeness
	return nil
}

// handleErrorResult handles error results from screens.
func (i *BrowseTimelineIntent) handleErrorResult(result *screens.ErrorResult) tea.Cmd {
	// Store error and return to previous state
	i.state.deleteError = result.Err
	return nil
}
