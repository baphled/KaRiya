// Package browse_timeline implements the BrowseTimeline intent for browsing career events.
package browse_timeline

import (
	"context"
	"fmt"
	"sort"

	"github.com/baphled/kariya/internal/cli/behaviors"
	"github.com/baphled/kariya/internal/cli/components"
	"github.com/baphled/kariya/internal/cli/intents"
	"github.com/baphled/kariya/internal/cli/screens"
	"github.com/baphled/kariya/internal/cli/screens/timeline"
	"github.com/baphled/kariya/internal/cli/uikit/feedback"
	"github.com/baphled/kariya/internal/cli/uikit/primitives"
	"github.com/baphled/kariya/internal/domain/career"
	tea "github.com/charmbracelet/bubbletea"
)

// Ensure Intent implements ScreenResultHandler interface.
var _ intents.ScreenResultHandler = (*Intent)(nil)

// Intent implements the Intent interface for browsing career events.
// It owns the complete lifecycle of timeline browsing, including:
// - Displaying a filtered and sorted timeline of events
// - Selecting and viewing event details
// - Returning the selected event or cancelling
type Intent struct {
	// Embed BaseIntent for terminal awareness, logo, and state management.
	*intents.BaseIntent

	// context is the input context passed to the intent.
	context *IntentContext

	// state tracks the current state of the intent (typed enum).
	state State

	// active indicates whether this intent is currently active.
	active bool

	// result is the final result of the intent (set when complete).
	result *intents.IntentResult[*Result]

	// --- Flattened state fields ---

	// filteredEvents are the events after applying current filters.
	filteredEvents []*career.CareerEvent

	// selectedIndex is the index of the currently selected event.
	selectedIndex int

	// filters is the current filter and sort state.
	filters *Filters

	// filterStack tracks active filters in FIFO order for progressive clearing.
	filterStack *intents.FilterStack

	// selectedEvent is the event currently being viewed.
	selectedEvent *career.CareerEvent

	// selectedFacts are facts selected from the event.
	selectedFacts []*career.Fact

	// viewedEvents tracks events viewed during the session.
	viewedEvents []*career.CareerEvent

	// deleteError stores any error from delete operation.
	deleteError error

	// --- Screen Orchestration ---

	// listScreen is the timeline event list screen.
	listScreen *timeline.TimelineEventListScreen

	// activeScreen holds the current screen being displayed.
	activeScreen screens.Screen

	// filterModal holds the filter modal (shown over the list).
	filterModal *components.FilterModalModel

	// searchModal holds the search modal for text search.
	searchModal *components.EventSearchModal

	// sortModal holds the sort modal for sorting events.
	sortModal *components.EventSortModal

	// deleteModal holds the delete confirmation modal (shown over the list).
	deleteModal *feedback.ConfirmModal

	// quickAddModal holds the quick add event modal (shown over the list).
	quickAddModal *components.QuickAddEventModal

	// editModal holds the edit event modal (shown over the list).
	editModal *components.EditEventModal

	// viewDetailModal holds the event detail viewer modal (shown over the list).
	viewDetailModal *components.ViewEventDetailModal

	// viewSkillsModal holds the skills viewer modal (shown over event detail).
	viewSkillsModal *components.ViewEventSkillsModal

	// errorModal holds the error modal (shown when operations fail).
	errorModal *feedback.Modal
}

// NewIntent creates a new BrowseTimeline intent.
func NewIntent(ctx *IntentContext) (*Intent, error) {
	// Validate the context.
	if err := ctx.Validate(); err != nil {
		return nil, err
	}

	// Create BaseIntent for terminal awareness and state management.
	base := intents.NewBaseIntent()

	intent := &Intent{
		BaseIntent:     base,
		context:        ctx,
		state:          StateTimeline,
		filteredEvents: ctx.Events,
		selectedIndex:  0,
		filters:        ctx.InitialFilters,
		filterStack:    intents.NewFilterStack(),
		selectedFacts:  make([]*career.Fact, 0),
		viewedEvents:   make([]*career.CareerEvent, 0),
		active:         true,
	}

	return intent, nil
}

// Init is called when the intent is activated.
func (i *Intent) Init() tea.Cmd {
	// Apply initial filters and sorting to the provided events.
	i.applyFilters()

	if len(i.filteredEvents) > 0 {
		i.selectedEvent = i.filteredEvents[0]
	}

	// Initialize with timeline list screen.
	i.state = StateTimeline
	i.transitionToScreen(timeline.NewTimelineEventListScreen(i.filteredEvents))
	return nil
}

// Update processes a message in the intent.
func (i *Intent) Update(msg tea.Msg) tea.Cmd {
	if !i.active {
		return nil
	}

	// PATTERN 4: Modal Priority Handling
	// Modals are checked in priority order: error -> search -> filter -> sort -> other modals
	// Only ONE modal can be active at a time

	// If error modal is visible, handle it first (highest priority).
	if i.errorModal != nil {
		if keyMsg, ok := msg.(tea.KeyMsg); ok && keyMsg.Type == tea.KeyEsc {
			i.errorModal = nil
			return nil
		}
		return nil
	}

	// Handle search modal if visible.
	if i.searchModal != nil && i.searchModal.IsVisible() {
		cmd, applied, searchData := i.searchModal.Update(msg)
		if applied && searchData != nil {
			i.filters.SearchText = searchData.SearchText
			if searchData.SearchText != "" {
				i.filterStack.Push(intents.FilterLayerSearch)
			}
			i.applyFilters()
			i.transitionToScreen(timeline.NewTimelineEventListScreen(i.filteredEvents))
		}
		return cmd
	}

	// Handle filter modal if visible.
	if i.filterModal != nil && i.filterModal.IsVisible() {
		cmd, applied, filterData := i.filterModal.Update(msg)
		if applied && filterData != nil {
			i.filters.Companies = filterData.Companies
			i.filters.Categories = filterData.Categories
			i.filters.Projects = filterData.Projects
			i.filters.SortBy = filterData.SortBy
			i.filters.SortOrder = filterData.SortOrder
			if len(filterData.Companies) > 0 {
				i.filterStack.Push(intents.FilterLayerCompany)
			}
			if len(filterData.Categories) > 0 {
				i.filterStack.Push(intents.FilterLayerCategory)
			}
			if len(filterData.Projects) > 0 {
				i.filterStack.Push(intents.FilterLayerProject)
			}
			i.applyFilters()
			i.transitionToScreen(timeline.NewTimelineEventListScreen(i.filteredEvents))
		}
		return cmd
	}

	// Handle sort modal if visible.
	if i.sortModal != nil && i.sortModal.IsVisible() {
		cmd, applied, sortData := i.sortModal.Update(msg)
		if applied && sortData != nil {
			i.filters.SortBy = sortData.SortBy
			i.filters.SortOrder = sortData.SortOrder
			if sortData.SortBy != "date" || sortData.SortOrder != "desc" {
				i.filterStack.Push(intents.FilterLayerSort)
			}
			i.applyFilters()
			i.transitionToScreen(timeline.NewTimelineEventListScreen(i.filteredEvents))
		}
		return cmd
	}

	// Handle quick add modal if visible.
	if i.quickAddModal != nil && i.quickAddModal.IsVisible() {
		cmd, completed, eventData := i.quickAddModal.Update(msg)
		if !i.quickAddModal.IsVisible() {
			if completed && eventData != nil {
				ctx := i.getContext()
				newEvent := eventData.ToCareerEvent()
				captureErr := i.context.CLIEventService.CaptureEvent(
					ctx,
					newEvent.Text,
					newEvent.Date,
					"manual",
				)
				if captureErr != nil {
					i.ShowErrorModal("Quick Add Failed", captureErr.Error())
					return cmd
				}
				refreshedEvents, listErr := i.context.CLIEventService.ListEvents(ctx, nil)
				if listErr != nil {
					i.ShowErrorModal("Refresh Failed", listErr.Error())
					return cmd
				}
				i.context.Events = refreshedEvents
				i.applyFilters()
				i.transitionToScreen(timeline.NewTimelineEventListScreen(i.filteredEvents))
			}
			i.quickAddModal = nil
		}
		return cmd
	}

	// Handle edit modal if visible.
	if i.editModal != nil && i.editModal.IsVisible() {
		cmd, completed, eventData := i.editModal.Update(msg)
		if !i.editModal.IsVisible() {
			if completed && eventData != nil {
				ctx := i.getContext()
				originalEvent := i.editModal.GetOriginalEvent()
				updatedEvent := eventData.ToCareerEvent(
					originalEvent.ID,
					originalEvent.CreatedAt,
					originalEvent.UpdatedAt,
				)
				if err := i.context.CLIEventService.UpdateEventMetadata(ctx, updatedEvent); err != nil {
					i.ShowErrorModal("Edit Failed", err.Error())
					return cmd
				}
				for idx, evt := range i.context.Events {
					if evt.ID == updatedEvent.ID {
						i.context.Events[idx] = updatedEvent
						break
					}
				}
				i.applyFilters()
				i.transitionToScreen(timeline.NewTimelineEventListScreen(i.filteredEvents))
			}
			i.editModal = nil
		}
		return cmd
	}

	// Handle delete modal if visible.
	if i.deleteModal != nil && i.deleteModal.IsVisible() {
		cmd, confirmed := i.deleteModal.Update(msg)
		if !i.deleteModal.IsVisible() {
			if confirmed && i.selectedEvent != nil {
				ctx := i.getContext()
				if err := i.context.CLIEventService.DeleteEvent(ctx, i.selectedEvent.ID); err != nil {
					i.ShowErrorModal("Delete Failed", err.Error())
					i.deleteError = err
					return cmd
				}
				deletedID := i.selectedEvent.ID
				newEvents := make([]*career.CareerEvent, 0, len(i.context.Events)-1)
				for _, evt := range i.context.Events {
					if evt.ID != deletedID {
						newEvents = append(newEvents, evt)
					}
				}
				i.context.Events = newEvents
				i.applyFilters()
				i.transitionToScreen(timeline.NewTimelineEventListScreen(i.filteredEvents))
			}
			i.selectedEvent = nil
		}
		return cmd
	}

	// Handle view skills modal if visible.
	if i.viewSkillsModal != nil && i.viewSkillsModal.IsVisible() {
		_, cmd := i.viewSkillsModal.Update(msg)
		if !i.viewSkillsModal.IsVisible() {
			i.viewSkillsModal = nil
		}
		return cmd
	}

	// Handle view detail modal if visible.
	if i.viewDetailModal != nil && i.viewDetailModal.IsVisible() {
		if keyMsg, ok := msg.(tea.KeyMsg); ok && keyMsg.String() == "s" {
			return i.showSkillsForCurrentEvent()
		}
		_, cmd := i.viewDetailModal.Update(msg)
		if !i.viewDetailModal.IsVisible() {
			i.viewDetailModal = nil
		}
		return cmd
	}

	// Handle key messages for shortcuts.
	if keyMsg, ok := msg.(tea.KeyMsg); ok {
		switch keyMsg.String() {
		case "/":
			// Open search modal.
			return i.openSearchModal()
		case "f":
			// Open filter modal.
			return i.openFilterModal()
		case "s":
			// Open sort modal.
			return i.openSortModal()
		case "x":
			// Clear filters progressively.
			if i.HasActiveFilters() {
				i.ClearFilters()
				i.applyFilters()
				i.transitionToScreen(timeline.NewTimelineEventListScreen(i.filteredEvents))
				return nil
			}
		}
	}

	// Delegate to active screen.
	cmd, result := i.activeScreen.Update(msg)
	if result != nil {
		return tea.Batch(cmd, i.handleScreenResult(result))
	}
	return cmd
}

// View renders the current state of the intent.
func (i *Intent) View() string {
	if i.activeScreen == nil {
		return "No active screen"
	}

	// Type-assert to screens with RenderContent method.
	switch screen := i.activeScreen.(type) {
	case *timeline.TimelineEventListScreen:
		// Create StandardView with table content.
		view := i.CreateViewWithBreadcrumbs("Main Menu", "Browse Timeline", i.getStateName())
		view.WithContent(screen.RenderContent())
		view.WithHelp(i.getContextHelp())

		// Render the complete view FIRST.
		baseView := view.Render()

		// Modal Overlay Rendering (direct checks, priority order).
		// Error modal (highest priority) - uses special adapter rendering.
		if i.errorModal != nil {
			width, height := i.getTerminalDimensions()
			adapter := intents.NewErrorModalAdapter(i.errorModal, width, height, i.Theme())
			return adapter.RenderOverlay(baseView)
		}

		// Search modal.
		if i.searchModal != nil && i.searchModal.IsVisible() {
			return behaviors.RenderModalOverlay(i.searchModal, baseView)
		}

		// Filter modal.
		if i.filterModal != nil && i.filterModal.IsVisible() {
			return behaviors.RenderModalOverlay(i.filterModal, baseView)
		}

		// Sort modal.
		if i.sortModal != nil && i.sortModal.IsVisible() {
			return behaviors.RenderModalOverlay(i.sortModal, baseView)
		}

		// Quick add modal.
		if i.quickAddModal != nil && i.quickAddModal.IsVisible() {
			return behaviors.RenderModalOverlay(i.quickAddModal, baseView)
		}

		// Edit modal.
		if i.editModal != nil && i.editModal.IsVisible() {
			return behaviors.RenderModalOverlay(i.editModal, baseView)
		}

		// Delete modal.
		if i.deleteModal != nil && i.deleteModal.IsVisible() {
			return behaviors.RenderModalOverlay(i.deleteModal, baseView)
		}

		// View skills modal.
		if i.viewSkillsModal != nil && i.viewSkillsModal.IsVisible() {
			return behaviors.RenderModalOverlay(i.viewSkillsModal, baseView)
		}

		// View detail modal.
		if i.viewDetailModal != nil && i.viewDetailModal.IsVisible() {
			return behaviors.RenderModalOverlay(i.viewDetailModal, baseView)
		}

		return baseView

	case *timeline.EventDeleteConfirmScreen:
		return screen.View()

	default:
		return i.activeScreen.View()
	}
}

// Result returns the intent result.
func (i *Intent) Result() *intents.IntentResult[interface{}] {
	if i.result == nil {
		return nil
	}

	return &intents.IntentResult[interface{}]{
		Status:   i.result.Status,
		Data:     i.result.Data,
		Error:    i.result.Error,
		Metadata: i.result.Metadata,
	}
}

// setCancelled marks the intent as cancelled.
func (i *Intent) setCancelled() {
	i.result = &intents.IntentResult[*Result]{
		Status: intents.Cancelled,
	}
	i.active = false
}

// getContext returns a context for service calls.
func (i *Intent) getContext() context.Context {
	return context.Background()
}

// applyFilters filters the events based on current filter state.
func (i *Intent) applyFilters() {
	filtered := make([]*career.CareerEvent, 0)

	for _, evt := range i.context.Events {
		// Apply text search (if specified).
		if i.filters.SearchText != "" {
			categoriesStr := ""
			if len(evt.Categories) > 0 {
				categoriesStr = evt.Categories[0]
			}
			if !intents.SearchableFields(i.filters.SearchText, evt.Text, evt.Company, categoriesStr, evt.Project) {
				continue
			}
		}

		// Apply tag filters.
		if len(i.filters.Tags) > 0 {
			hasTag := false
			for _, filterTag := range i.filters.Tags {
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

		// Apply company filter.
		if len(i.filters.Companies) > 0 {
			hasCompany := false
			for _, filterCompany := range i.filters.Companies {
				if evt.Company == filterCompany {
					hasCompany = true
					break
				}
			}
			if !hasCompany {
				continue
			}
		}

		// Apply category filter.
		if len(i.filters.Categories) > 0 {
			hasCategory := false
			for _, filterCat := range i.filters.Categories {
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

		// Apply project filter.
		if len(i.filters.Projects) > 0 {
			hasProject := false
			for _, filterProject := range i.filters.Projects {
				if evt.Project == filterProject {
					hasProject = true
					break
				}
			}
			if !hasProject {
				continue
			}
		}

		// Event passes all filters.
		filtered = append(filtered, evt)
	}

	// Apply sorting.
	switch i.filters.SortBy {
	case "date":
		sort.Slice(filtered, func(a, b int) bool {
			if i.filters.SortOrder == "asc" {
				return filtered[a].Date.Before(filtered[b].Date)
			}
			return filtered[a].Date.After(filtered[b].Date)
		})
	case "text":
		sort.Slice(filtered, func(a, b int) bool {
			if i.filters.SortOrder == "asc" {
				return filtered[a].Text < filtered[b].Text
			}
			return filtered[a].Text > filtered[b].Text
		})
	}

	i.filteredEvents = filtered
}

// HasActiveFilters returns true if any non-default filters are active.
func (i *Intent) HasActiveFilters() bool {
	f := i.filters
	if f == nil {
		return false
	}

	return f.SearchText != "" ||
		len(f.Tags) > 0 ||
		len(f.Companies) > 0 ||
		len(f.Categories) > 0 ||
		len(f.Projects) > 0 ||
		f.DateFrom != "" ||
		f.DateTo != "" ||
		(f.SortBy != "" && f.SortBy != "date") ||
		(f.SortOrder != "" && f.SortOrder != "desc")
}

// HasVisibleSkillsModal returns true if the skills modal is currently visible.
func (i *Intent) HasVisibleSkillsModal() bool {
	return i.viewSkillsModal != nil && i.viewSkillsModal.IsVisible()
}

// HasVisibleErrorModal returns true if the error modal is currently visible.
func (i *Intent) HasVisibleErrorModal() bool {
	return i.errorModal != nil
}

// ShowErrorModal displays an error modal with the given title and message.
func (i *Intent) ShowErrorModal(title, message string) {
	i.errorModal = feedback.NewErrorModal(title, message)
}

// ClearFilters clears filters in FIFO order (most recent first).
func (i *Intent) ClearFilters() {
	if i.filterStack == nil || i.filterStack.IsEmpty() {
		i.clearAllFilters()
		return
	}

	layer := i.filterStack.Pop()

	switch layer {
	case intents.FilterLayerSearch:
		i.filters.SearchText = ""
	case intents.FilterLayerCompany:
		i.filters.Companies = []string{}
	case intents.FilterLayerCategory:
		i.filters.Categories = []string{}
	case intents.FilterLayerProject:
		i.filters.Projects = []string{}
	case intents.FilterLayerTags:
		i.filters.Tags = []string{}
	case intents.FilterLayerSort:
		i.filters.SortBy = "date"
		i.filters.SortOrder = "desc"
	}

	if !i.HasActiveFilters() {
		i.filterStack.Clear()
	}
}

// clearAllFilters resets all filters to default state.
func (i *Intent) clearAllFilters() {
	i.filters = &Filters{
		SearchText: "",
		Tags:       []string{},
		Companies:  []string{},
		Categories: []string{},
		Projects:   []string{},
		DateFrom:   "",
		DateTo:     "",
		SortBy:     "date",
		SortOrder:  "desc",
	}
	if i.filterStack != nil {
		i.filterStack.Clear()
	}
}

// ApplyFilters applies current filter state.
func (i *Intent) ApplyFilters() {
	i.applyFilters()
}

// RefreshData reloads/refreshes the filtered data.
func (i *Intent) RefreshData() tea.Cmd {
	i.applyFilters()
	i.transitionToScreen(timeline.NewTimelineEventListScreen(i.filteredEvents))
	return nil
}

// removeEventFromList removes an event from both the full list and filtered list.
func (i *Intent) removeEventFromList(eventID string) {
	for idx, evt := range i.context.Events {
		if evt.ID == eventID {
			i.context.Events = append(i.context.Events[:idx], i.context.Events[idx+1:]...)
			break
		}
	}

	for idx, evt := range i.filteredEvents {
		if evt.ID == eventID {
			i.filteredEvents = append(i.filteredEvents[:idx], i.filteredEvents[idx+1:]...)
			break
		}
	}

	i.transitionToScreen(timeline.NewTimelineEventListScreen(i.filteredEvents))
}

// transitionToScreen sets the active screen and updates state.
func (i *Intent) transitionToScreen(screen screens.Screen) {
	i.activeScreen = screen

	termInfo := i.GetTerminalInfo()
	if termInfo != nil && termInfo.Width > 0 && termInfo.Height > 0 {
		screen.SetTerminalInfo(termInfo.Width, termInfo.Height)
	}

	if theme := i.Theme(); theme != nil {
		screen.SetTheme(theme)
	}

	if logo := i.GetLogo(); logo != nil {
		screen.SetLogo(logo, i.GetLogoSpacing())
	}
}

// handleScreenResult processes a screen result and determines next action.
func (i *Intent) handleScreenResult(result interface{}) tea.Cmd {
	if result == nil {
		return nil
	}

	screenResult, ok := result.(screens.ScreenResult)
	if !ok {
		return nil
	}

	return intents.NewScreenResultDispatcher(i).Dispatch(screenResult)
}

// HandleCancel handles screen cancellation (back/escape).
func (i *Intent) HandleCancel(result *screens.CancelResult) tea.Cmd {
	switch i.state {
	case StateTimeline:
		i.setCancelled()
		return nil

	case StateDeleteConfirm:
		i.state = StateTimeline
		i.transitionToScreen(timeline.NewTimelineEventListScreen(i.filteredEvents))
		return nil

	default:
		i.setCancelled()
		return nil
	}
}

// HandleNavigate handles screen navigation results.
func (i *Intent) HandleNavigate(result *screens.NavigateResult) tea.Cmd {
	if actionData, ok := result.ResultData.(map[string]interface{}); ok {
		action, _ := actionData["action"].(string)
		switch action {
		case "add":
			termInfo := i.GetTerminalInfo()
			width, height := 120, 40
			if termInfo != nil && termInfo.Width > 0 && termInfo.Height > 0 {
				width, height = termInfo.Width, termInfo.Height
			}
			i.quickAddModal = components.NewQuickAddEventModal(width, height)
			return i.quickAddModal.Init()

		case "edit":
			if event, ok := actionData["event"].(*career.CareerEvent); ok {
				termInfo := i.GetTerminalInfo()
				width, height := 120, 40
				if termInfo != nil && termInfo.Width > 0 && termInfo.Height > 0 {
					width, height = termInfo.Width, termInfo.Height
				}
				i.editModal = components.NewEditEventModal(event, width, height)
				return i.editModal.Init()
			}
			return nil

		case "delete":
			if event, ok := actionData["event"].(*career.CareerEvent); ok {
				eventText := event.Text
				if len(eventText) > 50 {
					eventText = eventText[:47] + "..."
				}
				i.deleteModal = feedback.NewConfirmModal(
					"Delete Event",
					fmt.Sprintf("Are you sure you want to delete '%s'?", eventText),
				).WithVariant(feedback.ConfirmDestructive)
				i.selectedEvent = event
				return i.deleteModal.Init()
			}
			return nil

		case "filter":
			termInfo := i.GetTerminalInfo()
			width, height := 120, 40
			if termInfo != nil && termInfo.Width > 0 && termInfo.Height > 0 {
				width, height = termInfo.Width, termInfo.Height
			}
			currentFilters := &components.TimelineFilters{
				SearchText: i.filters.SearchText,
				Tags:       i.filters.Tags,
				Companies:  i.filters.Companies,
				Categories: i.filters.Categories,
				Projects:   i.filters.Projects,
				SortBy:     i.filters.SortBy,
				SortOrder:  i.filters.SortOrder,
			}
			i.filterModal = components.NewFilterModal(
				i.context.Events,
				currentFilters,
				width,
				height,
			)
			return i.filterModal.Init()

		default:
			return nil
		}
	}

	if confirmed, ok := result.ResultData.(bool); ok {
		return i.handleDeleteConfirmation(confirmed)
	}

	if event, ok := result.ResultData.(*career.CareerEvent); ok {
		i.selectedEvent = event
		i.viewedEvents = append(i.viewedEvents, event)
		termInfo := i.GetTerminalInfo()
		width, height := 120, 40
		if termInfo != nil && termInfo.Width > 0 && termInfo.Height > 0 {
			width, height = termInfo.Width, termInfo.Height
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
func (i *Intent) handleDeleteConfirmation(confirmed bool) tea.Cmd {
	if !confirmed {
		i.state = StateTimeline
		i.transitionToScreen(timeline.NewTimelineEventListScreen(i.filteredEvents))
		return nil
	}

	if i.selectedEvent == nil {
		i.state = StateTimeline
		i.transitionToScreen(timeline.NewTimelineEventListScreen(i.filteredEvents))
		return nil
	}

	if i.context.CLIEventService != nil {
		eventID := i.selectedEvent.ID
		err := i.context.CLIEventService.DeleteEvent(i.getContext(), eventID)
		if err != nil {
			i.deleteError = err
			i.state = StateTimeline
			i.transitionToScreen(timeline.NewTimelineEventListScreen(i.filteredEvents))
			return nil
		}

		i.removeEventFromList(eventID)
	}

	i.state = StateTimeline
	i.transitionToScreen(timeline.NewTimelineEventListScreen(i.filteredEvents))
	return nil
}

// HandleSubmit handles form submission results.
func (i *Intent) HandleSubmit(result *screens.SubmitResult) tea.Cmd {
	return nil
}

// HandleError handles error results from screens.
func (i *Intent) HandleError(result *screens.ErrorResult) tea.Cmd {
	i.deleteError = result.Err
	return nil
}

// getStateName returns a human-readable name for the current state.
func (i *Intent) getStateName() string {
	switch i.state {
	case StateTimeline:
		return "Timeline"
	case StateDeleteConfirm:
		return "Delete Confirmation"
	default:
		return "Unknown"
	}
}

// getTerminalDimensions returns current terminal dimensions with fallback defaults.
func (i *Intent) getTerminalDimensions() (width, height int) {
	width, height = behaviors.DefaultModalDimensions()
	if termInfo := i.GetTerminalInfo(); termInfo != nil && termInfo.Width > 0 && termInfo.Height > 0 {
		width, height = termInfo.Width, termInfo.Height
	}
	return
}

// showSkillsForCurrentEvent loads and displays skills for the currently selected event.
func (i *Intent) showSkillsForCurrentEvent() tea.Cmd {
	if i.selectedEvent == nil {
		return nil
	}

	var skills []*career.Skill
	if i.context.CLIEventService != nil {
		ctx := i.getContext()
		var err error
		skills, err = i.context.CLIEventService.GetSkillsForEvent(ctx, i.selectedEvent.ID)
		if err != nil {
			skills = []*career.Skill{}
		}
	}

	termInfo := i.GetTerminalInfo()
	width, height := 120, 40
	if termInfo != nil && termInfo.Width > 0 && termInfo.Height > 0 {
		width, height = termInfo.Width, termInfo.Height
	}

	theme := i.Theme()
	i.viewSkillsModal = components.NewViewEventSkillsModal(i.selectedEvent.ID, skills, theme)
	i.viewSkillsModal.SetDimensions(width, height)
	i.viewSkillsModal.Show()

	return nil
}

// openSearchModal creates and initializes the search modal.
func (i *Intent) openSearchModal() tea.Cmd {
	termInfo := i.GetTerminalInfo()
	width, height := 120, 40
	if termInfo != nil && termInfo.Width > 0 && termInfo.Height > 0 {
		width, height = termInfo.Width, termInfo.Height
	}

	i.searchModal = components.NewEventSearchModal(
		i.filters.SearchText,
		width,
		height,
	)
	return i.searchModal.Init()
}

// openSortModal creates and initializes the sort modal.
func (i *Intent) openSortModal() tea.Cmd {
	termInfo := i.GetTerminalInfo()
	width, height := 120, 40
	if termInfo != nil && termInfo.Width > 0 && termInfo.Height > 0 {
		width, height = termInfo.Width, termInfo.Height
	}

	current := &components.EventSortConfig{
		SortBy:    i.filters.SortBy,
		SortOrder: i.filters.SortOrder,
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
func (i *Intent) openFilterModal() tea.Cmd {
	termInfo := i.GetTerminalInfo()
	width, height := 120, 40
	if termInfo != nil && termInfo.Width > 0 && termInfo.Height > 0 {
		width, height = termInfo.Width, termInfo.Height
	}

	currentFilters := &components.TimelineFilters{
		SearchText: i.filters.SearchText,
		Tags:       i.filters.Tags,
		Companies:  i.filters.Companies,
		Categories: i.filters.Categories,
		Projects:   i.filters.Projects,
		SortBy:     i.filters.SortBy,
		SortOrder:  i.filters.SortOrder,
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
func (i *Intent) getContextHelp() string {
	theme := i.Theme()

	switch i.state {
	case StateTimeline:
		badges := []*primitives.Badge{
			primitives.NavigateBadge(theme),
			primitives.HelpKeyBadge("Enter", "View Details", theme),
			primitives.AddBadge(theme),
			primitives.EditBadge(theme),
			primitives.DeleteBadge(theme),
			primitives.SearchBadge(theme),
			primitives.FilterBadge(theme),
			primitives.HelpKeyBadge("s", "Sort", theme),
		}

		if i.HasActiveFilters() {
			badges = append(badges, primitives.HelpKeyBadge("x", "Clear filters", theme))
		}

		badges = append(badges, primitives.BackBadge(theme))

		return intents.CombineThemedFooters(
			intents.ThemedCustomFooter(theme, badges...),
			intents.ThemedGlobalBadges(theme),
		)

	case StateDeleteConfirm:
		return ""

	default:
		return intents.ThemedGlobalBadges(theme)
	}
}
