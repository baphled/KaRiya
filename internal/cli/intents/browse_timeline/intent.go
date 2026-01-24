// Package browse_timeline implements the BrowseTimeline intent for browsing career events.
package browse_timeline

import (
	"context"
	"fmt"

	"github.com/baphled/kariya/internal/cli/intents"
	"github.com/baphled/kariya/internal/cli/screens"
	"github.com/baphled/kariya/internal/cli/screens/timeline"
	"github.com/baphled/kariya/internal/cli/uikit/feedback"
	"github.com/baphled/kariya/internal/domain/career"
	tea "github.com/charmbracelet/bubbletea"
)

// Ensure Intent implements ScreenResultHandler interface.
var _ behaviors.ScreenResultHandler = (*Intent)(nil)

// NewIntent creates a new BrowseTimeline intent.
func NewIntent(ctx *IntentContext) (*Intent, error) {
	if err := ctx.Validate(); err != nil {
		return nil, err
	}

	base := intents.NewBaseIntent()

	intent := &Intent{
		BaseIntent:     base,
		context:        ctx,
		state:          StateTimeline,
		filteredEvents: ctx.Events,
		selectedIndex:  0,
		filters:        ctx.InitialFilters,
		filterStack:    behaviors.NewFilterStack(),
		selectedFacts:  make([]*career.Fact, 0),
		viewedEvents:   make([]*career.CareerEvent, 0),
		active:         true,
		modalRegistry:  intents.NewModalRegistry(),
	}

	return intent, nil
}

// Init is called when the intent is activated.
func (i *Intent) Init() tea.Cmd {
	i.applyFilters()

	if len(i.filteredEvents) > 0 {
		i.selectedEvent = i.filteredEvents[0]
	}

	i.state = StateTimeline
	i.transitionToScreen(timeline.NewTimelineEventListScreen(i.filteredEvents))
	return nil
}

// Update processes a message in the intent.
func (i *Intent) Update(msg tea.Msg) tea.Cmd {
	if !i.active {
		return nil
	}

	// Handle modals in priority order.
	if cmd := i.handleModalUpdates(msg); cmd != nil || i.hasActiveModal() {
		return cmd
	}

	// Handle key shortcuts.
	if keyMsg, ok := msg.(tea.KeyMsg); ok {
		if cmd := i.handleKeyShortcuts(keyMsg); cmd != nil {
			return cmd
		}
	}

	// Delegate to active screen.
	cmd, result := i.activeScreen.Update(msg)
	if result != nil {
		return tea.Batch(cmd, i.handleScreenResult(result))
	}
	return cmd
}

// noopCmd is a sentinel command to indicate a message was consumed.
// This prevents the message from propagating to the screen after a modal closes.
func noopCmd() tea.Msg { return nil }

// modalHandler defines the interface for handling a modal's update cycle.
type modalHandler struct {
	isActive func() bool
	update   func(tea.Msg) tea.Cmd
	isClosed func() bool
}

// handleModalUpdates handles updates for all modals in priority order.
// Returns a command (possibly noopCmd) if a modal consumed the message.
func (i *Intent) handleModalUpdates(msg tea.Msg) tea.Cmd {
	// Error modal has special handling (highest priority).
	if i.errorModal != nil {
		if keyMsg, ok := msg.(tea.KeyMsg); ok && keyMsg.Type == tea.KeyEsc {
			i.errorModal = nil
			return noopCmd
		}
		return noopCmd
	}

	// Define modal handlers in priority order.
	handlers := []modalHandler{
		{
			isActive: func() bool { return i.searchModal != nil && i.searchModal.IsVisible() },
			update:   i.updateSearchModal,
			isClosed: func() bool { return i.searchModal == nil || !i.searchModal.IsVisible() },
		},
		{
			isActive: func() bool { return i.filterModal != nil && i.filterModal.IsVisible() },
			update:   i.updateFilterModal,
			isClosed: func() bool { return i.filterModal == nil || !i.filterModal.IsVisible() },
		},
		{
			isActive: func() bool { return i.sortModal != nil && i.sortModal.IsVisible() },
			update:   i.updateSortModal,
			isClosed: func() bool { return i.sortModal == nil || !i.sortModal.IsVisible() },
		},
		{
			isActive: func() bool { return i.quickAddModal != nil && i.quickAddModal.IsVisible() },
			update:   i.updateQuickAddModal,
			isClosed: func() bool { return i.quickAddModal == nil || !i.quickAddModal.IsVisible() },
		},
		{
			isActive: func() bool { return i.editModal != nil && i.editModal.IsVisible() },
			update:   i.updateEditModal,
			isClosed: func() bool { return i.editModal == nil || !i.editModal.IsVisible() },
		},
		{
			isActive: func() bool { return i.deleteModal != nil && i.deleteModal.IsVisible() },
			update:   i.updateDeleteModal,
			isClosed: func() bool { return i.deleteModal == nil || !i.deleteModal.IsVisible() },
		},
		{
			isActive: func() bool { return i.viewSkillsModal != nil && i.viewSkillsModal.IsVisible() },
			update:   i.updateViewSkillsModal,
			isClosed: func() bool { return i.viewSkillsModal == nil || !i.viewSkillsModal.IsVisible() },
		},
		{
			isActive: func() bool { return i.viewDetailModal != nil && i.viewDetailModal.IsVisible() },
			update:   i.updateViewDetailModal,
			isClosed: func() bool { return i.viewDetailModal == nil || !i.viewDetailModal.IsVisible() },
		},
	}

	// Process the first active modal.
	for _, h := range handlers {
		if h.isActive() {
			cmd := h.update(msg)
			if h.isClosed() && cmd == nil {
				return noopCmd
			}
			return cmd
		}
	}

	return nil
}

// hasActiveModal returns true if any modal is currently visible.
func (i *Intent) hasActiveModal() bool {
	i.rebuildModalRegistry()
	return i.modalRegistry.HasVisibleModal()
}

// updateSearchModal handles search modal updates.
func (i *Intent) updateSearchModal(msg tea.Msg) tea.Cmd {
	cmd, applied, searchData := i.searchModal.Update(msg)
	if applied && searchData != nil {
		i.filters.SearchText = searchData.SearchText
		if searchData.SearchText != "" {
			i.filterStack.Push(behaviors.FilterLayerSearch)
		}
		i.RefreshData()
	}
	return cmd
}

// updateFilterModal handles filter modal updates.
func (i *Intent) updateFilterModal(msg tea.Msg) tea.Cmd {
	cmd, applied, filterData := i.filterModal.Update(msg)
	if applied && filterData != nil {
		i.filters.Companies = filterData.Companies
		i.filters.Categories = filterData.Categories
		i.filters.Projects = filterData.Projects
		i.filters.SortBy = filterData.SortBy
		i.filters.SortOrder = filterData.SortOrder
		if len(filterData.Companies) > 0 {
			i.filterStack.Push(behaviors.FilterLayerCompany)
		}
		if len(filterData.Categories) > 0 {
			i.filterStack.Push(behaviors.FilterLayerCategory)
		}
		if len(filterData.Projects) > 0 {
			i.filterStack.Push(behaviors.FilterLayerProject)
		}
		i.RefreshData()
	}
	return cmd
}

// updateSortModal handles sort modal updates.
func (i *Intent) updateSortModal(msg tea.Msg) tea.Cmd {
	cmd, applied, sortData := i.sortModal.Update(msg)
	if applied && sortData != nil {
		i.filters.SortBy = sortData.SortBy
		i.filters.SortOrder = sortData.SortOrder
		if sortData.SortBy != "date" || sortData.SortOrder != "desc" {
			i.filterStack.Push(behaviors.FilterLayerSort)
		}
		i.RefreshData()
	}
	return cmd
}

// updateQuickAddModal handles quick add modal updates.
func (i *Intent) updateQuickAddModal(msg tea.Msg) tea.Cmd {
	cmd, completed, eventData := i.quickAddModal.Update(msg)
	if !i.quickAddModal.IsVisible() {
		if completed && eventData != nil {
			ctx := i.getContext()
			newEvent := eventData.ToCareerEvent()
			if err := i.context.CLIEventService.CaptureEvent(ctx, newEvent.Text, newEvent.Date, "manual"); err != nil {
				i.ShowErrorModal("Quick Add Failed", err.Error())
				return cmd
			}
			refreshedEvents, listErr := i.context.CLIEventService.ListEvents(ctx, nil)
			if listErr != nil {
				i.ShowErrorModal("Refresh Failed", listErr.Error())
				return cmd
			}
			i.context.Events = refreshedEvents
			i.RefreshData()
		}
		i.quickAddModal = nil
	}
	return cmd
}

// updateEditModal handles edit modal updates.
func (i *Intent) updateEditModal(msg tea.Msg) tea.Cmd {
	cmd, completed, eventData := i.editModal.Update(msg)
	if !i.editModal.IsVisible() {
		if completed && eventData != nil {
			ctx := i.getContext()
			originalEvent := i.editModal.GetOriginalEvent()
			updatedEvent := eventData.ToCareerEvent(originalEvent.ID, originalEvent.CreatedAt, originalEvent.UpdatedAt)
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
			i.RefreshData()
		}
		i.editModal = nil
	}
	return cmd
}

// updateDeleteModal handles delete modal updates.
func (i *Intent) updateDeleteModal(msg tea.Msg) tea.Cmd {
	cmd, confirmed := i.deleteModal.Update(msg)
	if !i.deleteModal.IsVisible() {
		if confirmed && i.selectedEvent != nil {
			ctx := i.getContext()
			if err := i.context.CLIEventService.DeleteEvent(ctx, i.selectedEvent.ID); err != nil {
				i.ShowErrorModal("Delete Failed", err.Error())
				i.deleteError = err
				i.deleteModal = nil
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
			i.RefreshData()
		}
		i.selectedEvent = nil
		i.deleteModal = nil
	}
	return cmd
}

// updateViewSkillsModal handles view skills modal updates.
func (i *Intent) updateViewSkillsModal(msg tea.Msg) tea.Cmd {
	_, cmd := i.viewSkillsModal.Update(msg)
	if !i.viewSkillsModal.IsVisible() {
		i.viewSkillsModal = nil
	}
	return cmd
}

// updateViewDetailModal handles view detail modal updates.
func (i *Intent) updateViewDetailModal(msg tea.Msg) tea.Cmd {
	if keyMsg, ok := msg.(tea.KeyMsg); ok {
		switch keyMsg.String() {
		case "s":
			return i.showSkillsForCurrentEvent()
		case "e":
			if i.selectedEvent != nil {
				i.viewDetailModal.Hide()
				i.viewDetailModal = nil
				return i.openEditModalForEvent(i.selectedEvent)
			}
		case "d":
			if i.selectedEvent != nil {
				eventText := i.selectedEvent.Text
				if len(eventText) > 50 {
					eventText = eventText[:47] + "..."
				}
				i.viewDetailModal.Hide()
				i.viewDetailModal = nil
				i.deleteModal = feedback.NewConfirmModal(
					"Delete Event",
					fmt.Sprintf("Are you sure you want to delete '%s'?", eventText),
				).WithVariant(feedback.ConfirmDestructive)
				return i.deleteModal.Init()
			}
		}
	}
	_, cmd := i.viewDetailModal.Update(msg)
	if !i.viewDetailModal.IsVisible() {
		i.viewDetailModal = nil
	}
	return cmd
}

// handleKeyShortcuts handles keyboard shortcuts when no modal is active.
func (i *Intent) handleKeyShortcuts(keyMsg tea.KeyMsg) tea.Cmd {
	switch keyMsg.String() {
	case "/":
		return i.openSearchModal()
	case "f":
		return i.openFilterModal()
	case "s":
		return i.openSortModal()
	case "x":
		if i.HasActiveFilters() {
			i.ClearFilters()
			i.RefreshData()
			return nil
		}
	}
	return nil
}

// View renders the current state of the intent.
func (i *Intent) View() string {
	if i.activeScreen == nil {
		return "No active screen"
	}

	switch screen := i.activeScreen.(type) {
	case *timeline.TimelineEventListScreen:
		return i.renderTimelineView(screen)
	case *timeline.EventDeleteConfirmScreen:
		return screen.View()
	default:
		return i.activeScreen.View()
	}
}

// renderTimelineView renders the timeline list view with modal overlays.
func (i *Intent) renderTimelineView(screen *timeline.TimelineEventListScreen) string {
	view := i.CreateViewWithBreadcrumbs("Main Menu", "Browse Timeline", i.getStateName())
	view.WithContent(screen.RenderContent())
	view.WithHelp(i.getContextHelp())
	baseView := view.Render()

	// Rebuild registry and render any visible modal as overlay.
	i.rebuildModalRegistry()
	return i.modalRegistry.RenderOverlay(baseView)
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
