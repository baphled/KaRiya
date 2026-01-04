package models

import (
	"context"
	"fmt"
	"sort"
	"strings"

	"github.com/baphled/kariya/internal/cli/components"
	"github.com/baphled/kariya/internal/cli/styles"
	"github.com/baphled/kariya/internal/domain/career"
	careerrepo "github.com/baphled/kariya/internal/repository/career"
	careerservice "github.com/baphled/kariya/internal/service/career"
	"github.com/charmbracelet/bubbles/table"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
)

// ListModel represents the event list display screen using table-based rendering
// This is the refactored version using ListDeletionState and ListNavigationKeyHandler patterns
type ListModel struct {
	*BaseStandardModel
	events               []*career.CareerEvent
	filtered             []*career.CareerEvent
	service              *careerservice.Service
	ctx                  context.Context
	table                table.Model
	pagination           *PaginationHelper
	width                int
	height               int
	competencyFilter     string
	sortBy               string
	sortOrder            string
	selectedEvents       map[string]bool
	expandedIndices      map[int]bool
	err                  error
	helpFooter           components.HelpFooterModel
	deletionState        *ListDeletionState
	navigationKeyHandler *ListNavigationKeyHandler
	breadcrumbs          []string
	header               components.HeaderModel
	listContainer        *components.TableListContainer
}

// NewListModel creates a new refactored list model
func NewListModel(svc *careerservice.Service, ctx context.Context) *ListModel {
	columns := []table.Column{
		{Title: "Event", Width: 50},
		{Title: "Company", Width: 20},
		{Title: "Date", Width: 12},
	}

	t := table.New(
		table.WithColumns(columns),
		table.WithRows([]table.Row{}),
		table.WithFocused(true),
		table.WithHeight(15),
		table.WithWidth(100),
	)

	s := table.DefaultStyles()
	s.Header = s.Header.
		Foreground(styles.ColorAccentTeal).
		Bold(true).
		BorderStyle(lipgloss.NormalBorder()).
		BorderBottom(true).
		BorderForeground(styles.ColorAccentTeal)
	s.Selected = s.Selected.
		Foreground(styles.ColorAccentTeal).
		Background(styles.ColorBackground).
		Bold(true)
	t.SetStyles(s)

	m := &ListModel{
		BaseStandardModel: NewBaseStandardModel(),
		service:           svc,
		ctx:               ctx,
		events:            []*career.CareerEvent{},
		filtered:          []*career.CareerEvent{},
		table:             t,
		pagination:        NewPaginationHelper(10),
		width:             80,
		height:            20,
		sortBy:            "date",
		sortOrder:         "desc",
		selectedEvents:    make(map[string]bool),
		expandedIndices:   make(map[int]bool),
		deletionState:     NewListDeletionState(),
		helpFooter:        components.NewHelpFooter("list", 80),
		breadcrumbs:       []string{"Home", "Events"},
		header:            components.NewHeader("📝 Career Events", 80),
		listContainer:     components.NewTableListContainer(t, "📝 Career Events", 80),
	}
	// Create navigation key handler with callbacks
	m.navigationKeyHandler = NewListNavigationKeyHandler(m)
	m.loadEventsSync()
	return m
}

// loadEventsSync loads events synchronously from the service
func (m *ListModel) loadEventsSync() {
	// First pass: load with no pagination to apply filters and sort
	events, err := m.service.ListEvents(m.ctx, careerrepo.ListFilters{Limit: 1000})
	if err != nil {
		m.err = err
		return
	}

	m.events = events
	m.applyFiltersAndSort()
	m.pagination.SetTotalCount(len(m.filtered))
	m.updateTableRows()
	m.err = nil
}

// updateTableRows updates the table with rows from filtered events (current page only)
func (m *ListModel) updateTableRows() {
	var rows []table.Row
	pageEvents := m.getPageEvents()

	// Get the selected index within the current page
	selectedIdx := m.listContainer.GetSelectedIdx()

	// Ensure cursor is within valid bounds BEFORE creating rows
	if len(pageEvents) > 0 {
		if selectedIdx >= len(pageEvents) {
			selectedIdx = len(pageEvents) - 1
		}
		if selectedIdx < 0 {
			selectedIdx = 0
		}
	} else {
		selectedIdx = 0
	}

	for i, event := range pageEvents {
		eventText := event.Text
		if len(eventText) > 50 {
			eventText = eventText[:47] + "..."
		}

		// Add focus indicator for the selected row, or spaces for alignment
		if i == selectedIdx {
			eventText = "▶ " + eventText
		} else {
			eventText = "  " + eventText
		}

		company := event.Company
		if len(company) > 17 {
			company = company[:14] + "..."
		}

		dateStr := event.Date.Format("2006-01-02")

		row := table.Row{
			eventText,
			company,
			dateStr,
		}
		rows = append(rows, row)
	}

	m.table.SetRows(rows)
	// Sync the container and table cursor with the validated selectedIdx
	m.listContainer.SetSelectedIdx(selectedIdx)
	m.table.SetCursor(selectedIdx)
}

// Init initializes the model
func (m *ListModel) Init() tea.Cmd {
	return nil
}

// Update handles messages
func (m *ListModel) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	// Handle deletion confirmation if active
	if m.deletionState.IsConfirming() {
		cmd := m.deletionState.UpdateConfirmation(msg)

		if m.deletionState.IsConfirmed() {
			return m.performEventDeletion()
		}

		if m.deletionState.IsCancelled() {
			m.deletionState.Clear()
			return m, nil
		}

		return m, cmd
	}

	// Handle deletion message display
	if m.deletionState.ShowMessage {
		switch msg := msg.(type) {
		case tea.KeyMsg:
			if msg.String() == "esc" {
				m.deletionState.Clear()
				return m, nil
			}
		}
	}

	switch msg := msg.(type) {
	case tea.KeyMsg:
		// Let navigation key handler process navigation and item action keys
		if cmd := m.navigationKeyHandler.HandleNavigationKey(msg.String()); cmd != nil {
			return m, cmd
		}

		// Handle screen navigation keys that aren't list-specific
		switch msg.String() {
		case "esc":
			return m, func() tea.Msg { return BackMsg{} }
		case "q", "ctrl+c":
			return m, func() tea.Msg { return QuitMsg{} }
		}
	case tea.WindowSizeMsg:
		m.width = msg.Width
		m.height = msg.Height
		m.table.SetWidth(m.width)
		m.table.SetHeight(m.height - 10)
		m.listContainer.SetDimensions(m.width, m.height)
		m.updateTableRows()
		return m, nil
	}

	var cmd tea.Cmd
	m.table, cmd = m.table.Update(msg)
	return m, cmd
}

// performEventDeletion performs the actual event deletion
func (m *ListModel) performEventDeletion() (tea.Model, tea.Cmd) {
	if m.deletionState.DeletingItemID == "" {
		m.deletionState.Clear()
		return m, nil
	}

	err := m.service.DeleteEvent(m.ctx, m.deletionState.DeletingItemID)
	if err != nil {
		m.deletionState.SetErrorMsg(fmt.Sprintf("Error deleting event: %v", err))
		return m, nil
	}

	m.removeEventFromLists(m.deletionState.DeletingItemID)

	m.deletionState.SetSuccessMsg("Event deleted successfully")

	m.pagination.SetTotalCount(len(m.filtered))
	m.updateTableRows()

	if m.pagination.GetTotalCount() == 0 {
		m.listContainer.MoveToFirst()
	} else {
		if m.listContainer.GetSelectedIdx() >= len(m.table.Rows()) && m.listContainer.GetSelectedIdx() > 0 {
			m.listContainer.SetSelectedIdx(m.listContainer.GetSelectedIdx() - 1)
		}
	}

	return m, nil
}

// removeEventFromLists removes an event from both the events and filtered lists
func (m *ListModel) removeEventFromLists(eventID string) {
	for i, event := range m.events {
		if event.ID == eventID {
			m.events = append(m.events[:i], m.events[i+1:]...)
			break
		}
	}

	for i, event := range m.filtered {
		if event.ID == eventID {
			m.filtered = append(m.filtered[:i], m.filtered[i+1:]...)
			break
		}
	}

	delete(m.selectedEvents, eventID)
}

// View renders the event list
func (m *ListModel) View() string {
	// Render deletion confirmation dialog
	if m.deletionState.IsConfirming() {
		return m.deletionState.ConfirmationDialog.View()
	}

	// Render deletion message (success or error)
	if m.deletionState.ShowMessage {
		return m.renderDeletionMessage()
	}

	if m.err != nil {
		errorMsg := fmt.Sprintf("Error loading events: %v\n\nPress 'r' to retry or 'esc' to cancel", m.err)
		m.listContainer.SetErrorMessage(errorMsg).SetDimensions(m.width, m.height)
		return m.listContainer.Render()
	}

	// Update list container with current state
	m.listContainer.SetTable(m.table).
		SetDimensions(m.width, m.height).
		SetEmptyStateMessage("No events found").
		SetHelpFooterKey("list").
		SetBreadcrumbs(m.breadcrumbs)

	// Set pagination info
	paginationText := m.pagination.GetPaginationInfo("events")
	m.listContainer.SetPaginationInfo(paginationText)

	return m.listContainer.Render()
}

// renderDeletionMessage renders the deletion success/error message
func (m *ListModel) renderDeletionMessage() string {
	var messageContent string
	if m.deletionState.SuccessMsg != "" {
		messageContent = styles.SuccessBox.Render(m.deletionState.SuccessMsg + "\n\nPress 'esc' to continue")
	} else if m.deletionState.ErrorMsg != "" {
		messageContent = styles.ErrorBox.Render(m.deletionState.ErrorMsg + "\n\nPress 'esc' to continue")
	}

	headerView := m.header.View()
	footerView := components.NewFooter(m.width).View()

	return lipgloss.JoinVertical(
		lipgloss.Left,
		headerView,
		"",
		messageContent,
		"",
		footerView,
	)
}

// getPageEvents returns the events for the current page
func (m *ListModel) getPageEvents() []*career.CareerEvent {
	if m.pagination.GetTotalCount() == 0 {
		return []*career.CareerEvent{}
	}

	startIdx := m.pagination.GetPageStartIndex()
	endIdx := m.pagination.GetPageEndIndex()

	if startIdx >= len(m.filtered) {
		return []*career.CareerEvent{}
	}

	if endIdx > len(m.filtered) {
		endIdx = len(m.filtered)
	}

	return m.filtered[startIdx:endIdx]
}

// SetEvents sets the events to display
func (m *ListModel) SetEvents(events []*career.CareerEvent) {
	m.events = events
	m.applyFiltersAndSort()
	m.pagination.SetTotalCount(len(m.filtered)).GoToFirstPage()
	m.listContainer.MoveToFirst()
	m.expandedIndices = make(map[int]bool)
	m.updateTableRows()
}

// SetCompetencyFilter sets the competency focus filter
func (m *ListModel) SetCompetencyFilter(competency string) {
	m.competencyFilter = competency
	m.applyFiltersAndSort()
	m.pagination.SetTotalCount(len(m.filtered)).GoToFirstPage()
	m.listContainer.MoveToFirst()
	m.expandedIndices = make(map[int]bool)
	m.updateTableRows()
}

// SetSort sets the sort order
func (m *ListModel) SetSort(sortBy, sortOrder string) {
	m.sortBy = sortBy
	m.sortOrder = sortOrder
	m.applyFiltersAndSort()
	m.updateTableRows()
}

// applyFiltersAndSort applies filters and sorts the events
func (m *ListModel) applyFiltersAndSort() {
	m.filtered = m.filterEvents()
	m.sortEvents(m.filtered)
}

// filterEvents filters events based on current filters
func (m *ListModel) filterEvents() []*career.CareerEvent {
	var filtered []*career.CareerEvent

	for _, event := range m.events {
		if m.competencyFilter != "" {
			// Filter could be applied here based on event attributes
			// For now, just include all events if no filter is set
			if !strings.Contains(strings.ToLower(event.Text), strings.ToLower(m.competencyFilter)) {
				continue
			}
		}

		filtered = append(filtered, event)
	}

	return filtered
}

// sortEvents sorts events based on current sort settings
func (m *ListModel) sortEvents(events []*career.CareerEvent) {
	sort.SliceStable(events, func(i, j int) bool {
		var less bool

		switch m.sortBy {
		case "text":
			less = events[i].Text < events[j].Text
		case "company":
			less = events[i].Company < events[j].Company
		case "date":
			fallthrough
		default:
			less = events[i].Date.Before(events[j].Date)
		}

		if m.sortOrder == "asc" {
			return less
		}
		return !less
	})
}

// GetSelectedEvent returns the currently selected event
func (m *ListModel) GetSelectedEvent() *career.CareerEvent {
	cursor := m.listContainer.GetSelectedIdx()
	if cursor >= 0 && cursor < len(m.filtered) {
		return m.filtered[cursor]
	}
	return nil
}

// GetSelectedEvents returns all selected events
func (m *ListModel) GetSelectedEvents() []*career.CareerEvent {
	var selected []*career.CareerEvent
	for _, event := range m.filtered {
		if m.selectedEvents[event.ID] {
			selected = append(selected, event)
		}
	}
	return selected
}

// IsExpandedAt returns true if the event at the given index is expanded
func (m *ListModel) IsExpandedAt(idx int) bool {
	return m.expandedIndices[idx]
}

// SetBreadcrumbs sets the breadcrumb trail for navigation
func (m *ListModel) SetBreadcrumbs(crumbs []string) {
	m.breadcrumbs = crumbs
	m.header.SetBreadcrumbs(crumbs)
}

// Refresh reloads the events from the service
func (m *ListModel) Refresh() {
	m.loadEventsSync()
	m.pagination.GoToFirstPage()
	m.expandedIndices = make(map[int]bool)
	m.listContainer.MoveToFirst()
}

// GetError returns the last error
func (m *ListModel) GetError() error {
	return m.err
}

// GetEvents returns the currently filtered events
func (m *ListModel) GetEvents() []*career.CareerEvent {
	return m.filtered
}

// GetSelectedIdx returns the currently selected index
func (m *ListModel) GetSelectedIdx() int {
	return m.listContainer.GetSelectedIdx()
}

// ListItemCallbacks implementation for ListModel
// These methods implement the ListItemCallbacks interface required by ListNavigationKeyHandler

// OnView sends a message to view the selected event
func (m *ListModel) OnView(item interface{}) tea.Cmd {
	if event, ok := item.(*career.CareerEvent); ok {
		return func() tea.Msg { return EventActionSelectedMsg{Event: event, Action: EventActionView} }
	}
	return nil
}

// OnEdit sends a message to edit the selected event
func (m *ListModel) OnEdit(item interface{}) tea.Cmd {
	if event, ok := item.(*career.CareerEvent); ok {
		return func() tea.Msg { return EventActionSelectedMsg{Event: event, Action: EventActionEdit} }
	}
	return nil
}

// OnDelete initiates deletion of the selected event
func (m *ListModel) OnDelete(item interface{}) {
	if event, ok := item.(*career.CareerEvent); ok {
		m.deletionState.ShowConfirmation("Event", event.Text)
		m.deletionState.DeletingItemID = event.ID
	}
}

// OnToggleSelection toggles the selection state of an item
func (m *ListModel) OnToggleSelection(item interface{}) {
	if event, ok := item.(*career.CareerEvent); ok {
		m.selectedEvents[event.ID] = !m.selectedEvents[event.ID]
	}
}

// HasSelectedItem returns the currently selected event
func (m *ListModel) HasSelectedItem() interface{} {
	return m.GetSelectedEvent()
}

// MoveUp moves the selection up by count items
func (m *ListModel) MoveUp(count int) {
	m.listContainer.MoveUp(count)
}

// MoveDown moves the selection down by count items
func (m *ListModel) MoveDown(count int) {
	m.listContainer.MoveDown(count)
}

// MoveToFirst moves selection to the first item
func (m *ListModel) MoveToFirst() {
	m.listContainer.MoveToFirst()
}

// MoveToLast moves selection to the last item
func (m *ListModel) MoveToLast() {
	m.listContainer.MoveToLast()
}

// UpdateDisplay refreshes the table display
func (m *ListModel) UpdateDisplay() {
	m.updateTableRows()
}

// GetRowCount returns the number of visible rows
func (m *ListModel) GetRowCount() int {
	return len(m.table.Rows())
}

// GetCurrentIndex returns the current selection index
func (m *ListModel) GetCurrentIndex() int {
	return m.listContainer.GetSelectedIdx()
}

// GetPageSize returns the page size for pagination
func (m *ListModel) GetPageSize() int {
	return m.pagination.GetPageSize()
}

// nextPage moves to the next page of results
func (m *ListModel) nextPage() {
	m.pagination.NextPage()
	m.listContainer.MoveToFirst()
	m.updateTableRows()
}

// prevPage moves to the previous page of results
func (m *ListModel) prevPage() {
	m.pagination.PrevPage()
	m.listContainer.MoveToFirst()
	m.updateTableRows()
}

// goToFirstItem moves to the first item
func (m *ListModel) goToFirstItem() {
	m.pagination.GoToFirstPage()
	m.listContainer.MoveToFirst()
}

// goToLastItem moves to the last item
func (m *ListModel) goToLastItem() {
	m.pagination.GoToLastPage()
	pageEvents := m.getPageEvents()
	if len(pageEvents) > 0 {
		m.listContainer.SetSelectedIdx(len(pageEvents) - 1)
	}
}

// Accessor methods for backward compatibility with tests

// getTotalPages calculates the total number of pages
func (m *ListModel) getTotalPages() int {
	return m.pagination.GetTotalPages()
}

// pageSize getter for tests
func (m *ListModel) getPageSize() int {
	return m.pagination.GetPageSize()
}

// currentPage getter for tests (note: using reflection or adding public field would be cleaner)
// These are added as convenience methods for tests
// Tests should ideally be refactored to use public API only

// NextPage implements ListItemCallbacks interface - moves to the next page
func (m *ListModel) NextPage() {
	m.nextPage()
}

// PrevPage implements ListItemCallbacks interface - moves to the previous page
func (m *ListModel) PrevPage() {
	m.prevPage()
}

// GoToFirstPage implements ListItemCallbacks interface - goes to the first page
func (m *ListModel) GoToFirstPage() {
	m.goToFirstItem()
	m.updateTableRows()
}

// GoToLastPage implements ListItemCallbacks interface - goes to the last page
func (m *ListModel) GoToLastPage() {
	m.goToLastItem()
	m.updateTableRows()
}
