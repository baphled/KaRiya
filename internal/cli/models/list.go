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
type ListModel struct {
	*BaseStandardModel
	events            []*career.CareerEvent
	filtered          []*career.CareerEvent
	service           *careerservice.Service
	ctx               context.Context
	table             table.Model
	currentPage       int
	pageSize          int
	totalCount        int
	width             int
	height            int
	competencyFilter  string
	sortBy            string
	sortOrder         string
	selectedEvents    map[string]bool
	expandedIndices   map[int]bool
	err               error
	helpFooter        components.HelpFooterModel
	deletionConfirm   *ConfirmationDialog
	deletingEventID   string
	deleteSuccessMsg  string
	deleteErrorMsg    string
	showDeleteMessage bool
	breadcrumbs       []string
	header            components.HeaderModel
	listContainer     *components.TableListContainer
}

// NewListModel creates a new list model
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
		width:             80,
		height:            20,
		pageSize:          10,
		currentPage:       1,
		totalCount:        0,
		sortBy:            "date",
		sortOrder:         "desc",
		selectedEvents:    make(map[string]bool),
		expandedIndices:   make(map[int]bool),
		deletionConfirm:   nil,
		deletingEventID:   "",
		deleteSuccessMsg:  "",
		deleteErrorMsg:    "",
		showDeleteMessage: false,
		helpFooter:        components.NewHelpFooter("list", 80),
		breadcrumbs:       []string{"Home", "Events"},
		header:            components.NewHeader("📝 Career Events", 80),
		listContainer:     components.NewTableListContainer(t, "Events", 80),
	}
	m.loadEventsSync()
	return m
}

// loadEventsSync loads events synchronously from the service
func (m *ListModel) loadEventsSync() {
	events, err := m.service.ListEvents(m.ctx, careerrepo.ListFilters{Limit: m.pageSize})
	if err != nil {
		m.err = err
	} else {
		m.events = events
		m.applyFiltersAndSort()
		m.totalCount = len(m.filtered)
		m.updateTableRows()
		m.err = nil
	}
}

// updateTableRows updates the table with rows from filtered events
func (m *ListModel) updateTableRows() {
	var rows []table.Row

	for i, event := range m.filtered {
		eventText := event.Text
		if len(eventText) > 47 {
			eventText = eventText[:44] + "..."
		}

		company := event.Company
		if len(company) > 17 {
			company = company[:14] + "..."
		}

		dateStr := event.Date.Format("2006-01-02")

		// Add indicator for selected row based on container's selectedIdx
		indicator := "  "
		if i == m.listContainer.GetSelectedIdx() {
			indicator = "▶ "
		}

		row := table.Row{
			indicator + eventText,
			company,
			dateStr,
		}
		rows = append(rows, row)
	}

	m.table.SetRows(rows)
	// Ensure cursor is within valid bounds
	if len(rows) > 0 {
		if m.table.Cursor() >= len(rows) {
			m.listContainer.SetSelectedIdx(len(rows) - 1)
		}
		if m.table.Cursor() < 0 {
			m.listContainer.SetSelectedIdx(0)
		}
	}
	// Re-sync the table's cursor with the container's selectedIdx after SetRows
	m.table.SetCursor(m.listContainer.GetSelectedIdx())
}

// Init initializes the model
func (m *ListModel) Init() tea.Cmd {
	return nil
}

// Update handles messages
func (m *ListModel) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	if m.deletionConfirm != nil {
		updatedDialog, cmd := m.deletionConfirm.Update(msg)
		m.deletionConfirm = updatedDialog

		if m.deletionConfirm.IsConfirmed() {
			return m.performEventDeletion()
		}

		if m.deletionConfirm.IsCancelled() {
			m.deletionConfirm = nil
			m.deletingEventID = ""
			return m, nil
		}

		return m, cmd
	}

	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch msg.String() {
		case "up", "k":
			m.listContainer.MoveUp(1)
			m.updateTableRows()
			return m, nil
		case "down", "j":
			m.listContainer.MoveDown(1)
			m.updateTableRows()
			return m, nil
		case "pgup", "ctrl+b":
			if m.listContainer.GetSelectedIdx() >= m.pageSize {
				m.listContainer.MoveUp(m.pageSize)
			} else {
				m.listContainer.MoveToFirst()
			}
			m.updateTableRows()
			return m, nil
		case "pgdn", "ctrl+f":
			rows := len(m.table.Rows())
			if m.listContainer.GetSelectedIdx()+m.pageSize < rows {
				m.listContainer.MoveDown(m.pageSize)
			} else {
				m.listContainer.MoveToLast()
			}
			m.updateTableRows()
			return m, nil
		case "home", "g":
			m.listContainer.MoveToFirst()
			m.updateTableRows()
			return m, nil
		case "end", "G":
			m.listContainer.MoveToLast()
			m.updateTableRows()
			return m, nil
		case "enter":
			if len(m.filtered) > 0 {
				event := m.GetSelectedEvent()
				if event != nil {
					return m, func() tea.Msg { return EventActionSelectedMsg{Event: event, Action: EventActionView} }
				}
			}
			return m, nil
		case " ", "space":
			if len(m.filtered) > 0 {
				event := m.GetSelectedEvent()
				if event != nil {
					m.selectedEvents[event.ID] = !m.selectedEvents[event.ID]
				}
			}
			return m, nil
		case "x", "d":
			if len(m.filtered) > 0 {
				event := m.GetSelectedEvent()
				if event != nil {
					m.showDeleteConfirmation(event)
				}
			}
			return m, nil
		case "v":
			if len(m.filtered) > 0 {
				event := m.GetSelectedEvent()
				if event != nil {
					return m, func() tea.Msg { return EventActionSelectedMsg{Event: event, Action: EventActionView} }
				}
			}
			return m, nil
		case "e":
			if len(m.filtered) > 0 {
				event := m.GetSelectedEvent()
				if event != nil {
					return m, func() tea.Msg { return EventActionSelectedMsg{Event: event, Action: EventActionEdit} }
				}
			}
			return m, nil
		case "esc":
			if m.showDeleteMessage {
				m.showDeleteMessage = false
				m.deleteSuccessMsg = ""
				m.deleteErrorMsg = ""
				return m, nil
			} else {
				return m, func() tea.Msg { return BackMsg{} }
			}
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

// showDeleteConfirmation shows the deletion confirmation dialog
func (m *ListModel) showDeleteConfirmation(event *career.CareerEvent) {
	message := fmt.Sprintf("Delete event: \"%s\"?\n\nThis action cannot be undone.", truncateText(event.Text, 60))
	m.deletionConfirm = NewConfirmationDialog("Delete Event", message)
	m.deletingEventID = event.ID
}

// performEventDeletion performs the actual event deletion
func (m *ListModel) performEventDeletion() (tea.Model, tea.Cmd) {
	if m.deletingEventID == "" {
		m.deletionConfirm = nil
		return m, nil
	}

	err := m.service.DeleteEvent(m.ctx, m.deletingEventID)
	if err != nil {
		m.deleteErrorMsg = fmt.Sprintf("Error deleting event: %v", err)
		m.showDeleteMessage = true
		m.deletionConfirm = nil
		m.deletingEventID = ""
		return m, nil
	}

	m.removeEventFromLists(m.deletingEventID)

	m.deleteSuccessMsg = "Event deleted successfully"
	m.showDeleteMessage = true
	m.deletionConfirm = nil
	m.deletingEventID = ""

	m.totalCount = len(m.filtered)
	m.updateTableRows()

	if m.totalCount == 0 {
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
	if m.deletionConfirm != nil {
		return m.deletionConfirm.View()
	}

	if m.showDeleteMessage {
		return m.renderDeleteMessage()
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
		SetHelpFooterKey("list")

	// Set pagination info
	startIdx := (m.currentPage-1)*m.pageSize + 1
	endIdx := startIdx + len(m.getPageEvents()) - 1
	if m.totalCount == 0 {
		startIdx = 0
		endIdx = 0
	}
	paginationText := fmt.Sprintf("Showing %d-%d of %d events", startIdx, endIdx, m.totalCount)
	m.listContainer.SetPaginationInfo(paginationText)

	return m.listContainer.Render()
}

// renderDeleteMessage renders the deletion success/error message
func (m *ListModel) renderDeleteMessage() string {
	var messageContent string
	if m.deleteSuccessMsg != "" {
		messageContent = styles.SuccessBox.Render(m.deleteSuccessMsg + "\n\nPress 'esc' to continue")
	} else if m.deleteErrorMsg != "" {
		messageContent = styles.ErrorBox.Render(m.deleteErrorMsg + "\n\nPress 'esc' to continue")
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
	if m.totalCount == 0 {
		return []*career.CareerEvent{}
	}

	startIdx := (m.currentPage - 1) * m.pageSize
	endIdx := startIdx + m.pageSize

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
	m.totalCount = len(m.filtered)
	m.currentPage = 1
	m.listContainer.MoveToFirst()
	m.expandedIndices = make(map[int]bool)
	m.updateTableRows()
}

// SetCompetencyFilter sets the competency focus filter
func (m *ListModel) SetCompetencyFilter(competency string) {
	m.competencyFilter = competency
	m.applyFiltersAndSort()
	m.totalCount = len(m.filtered)
	m.currentPage = 1
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
	m.currentPage = 1
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

// getTotalPages calculates the total number of pages
func (m *ListModel) getTotalPages() int {
	if m.totalCount == 0 {
		return 1
	}
	pages := (m.totalCount + m.pageSize - 1) / m.pageSize
	return pages
}

// nextPage moves to the next page of results
func (m *ListModel) nextPage() {
	totalPages := m.getTotalPages()
	if m.currentPage < totalPages {
		m.currentPage++
		m.listContainer.MoveToFirst()
		m.updateTableRows()
	}
}

// prevPage moves to the previous page of results
func (m *ListModel) prevPage() {
	if m.currentPage > 1 {
		m.currentPage--
		m.listContainer.MoveToFirst()
		m.updateTableRows()
	}
}

// goToFirstItem moves to the first item
func (m *ListModel) goToFirstItem() {
	m.currentPage = 1
	m.listContainer.MoveToFirst()
}

// goToLastItem moves to the last item
func (m *ListModel) goToLastItem() {
	totalPages := m.getTotalPages()
	m.currentPage = totalPages
	pageEvents := m.getPageEvents()
	if len(pageEvents) > 0 {
		m.listContainer.SetSelectedIdx(len(pageEvents) - 1)
	}
}
