package models

import (
	"context"
	"fmt"

	"github.com/baphled/kariya/internal/cli/components"
	"github.com/baphled/kariya/internal/cli/styles"
	"github.com/baphled/kariya/internal/domain/career"
	careerrepo "github.com/baphled/kariya/internal/repository/career"
	careerservice "github.com/baphled/kariya/internal/service/career"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
)

// ListModel represents the event list screen
type ListModel struct {
	*BaseStandardModel
	service     *careerservice.Service
	ctx         context.Context
	events      []*career.CareerEvent
	totalCount  int
	currentPage int
	pageSize    int
	selectedIdx int
	width       int
	height      int
	err         error
	filterModel *FilterModel
	searchModel *SearchModel
	sortModel   *SortModel
	helpFooter  components.HelpFooterModel // Help footer
	header      components.HeaderModel     // Header component
	footer      components.FooterModel     // Footer component
	breadcrumbs []string                   // Navigation breadcrumb trail
}

// NewListModel creates a new list model
func NewListModel(svc *careerservice.Service, ctx context.Context) *ListModel {
	model := &ListModel{
		BaseStandardModel: NewBaseStandardModel(),
		service:           svc,
		ctx:               ctx,
		pageSize:          10,
		currentPage:       1,
		selectedIdx:       0,
		filterModel:       NewFilterModel(),
		searchModel:       NewSearchModel(),
		sortModel:         NewSortModel(),
		helpFooter:        components.NewHelpFooter("list", 80),
		header:            components.NewHeader("Career Events", 80),
		footer:            components.NewFooter(80),
	}

	// Load events
	model.loadEvents()

	return model
}

// loadEvents loads events from the service
func (m *ListModel) loadEvents() {
	filters := careerrepo.ListFilters{
		SortBy:    "date",
		SortOrder: "desc",
		Limit:     m.pageSize,
		Offset:    (m.currentPage - 1) * m.pageSize,
	}

	events, err := m.service.ListEvents(m.ctx, filters)
	if err != nil {
		m.err = err
		m.events = []*career.CareerEvent{}
		return
	}

	m.events = events

	// Get total count for pagination calculation
	totalCount, err := m.service.CountEvents(m.ctx, careerrepo.ListFilters{})
	if err != nil {
		m.totalCount = 0
		return
	}

	m.totalCount = totalCount
}

// Init initializes the model
func (m *ListModel) Init() tea.Cmd {
	return nil
}

// Update handles messages
func (m *ListModel) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.WindowSizeMsg:
		m.width = msg.Width
		m.height = msg.Height
		m.header.SetWidth(msg.Width)
		m.footer.SetWidth(msg.Width)
		m.helpFooter.SetWidth(msg.Width)
		return m, nil
	case tea.KeyMsg:
		switch msg.String() {
		case "esc":
			// Signal back navigation to parent
			return m, func() tea.Msg { return BackMsg{} }
		case "ctrl+c", "q":
			// Signal quit to parent
			return m, func() tea.Msg { return QuitMsg{} }
		case "up", "k":
			m.prevItem()
		case "down", "j":
			m.nextItem()
		case "pgup", "ctrl+b":
			m.prevPage()
		case "pgdn", "ctrl+f":
			m.nextPage()
		case "home", "g":
			m.goToFirstItem()
		case "end", "G":
			m.goToLastItem()
		case "enter":
			return m, m.viewSelectedEvent()
		}
	}

	return m, nil
}

// View renders the model
func (m *ListModel) View() string {
	var content []string

	// Error handling
	if m.err != nil {
		errorMsg := fmt.Sprintf("Error loading events: %v\n\nPress 'r' to retry or 'esc' to cancel", m.err)
		content = append(content, styles.ErrorBox.Render(errorMsg))
	}

	// Render list using ListContainer
	listContent := m.renderListWithContainer()
	content = append(content, listContent)

	// Help footer with keyboard shortcuts
	m.helpFooter.SetWidth(styles.MaxWidth(80))
	helpFooterContent := m.helpFooter.View()
	content = append(content, helpFooterContent)

	// Combine all content
	fullListContent := lipgloss.JoinVertical(
		lipgloss.Left,
		content...,
	)

	// Wrap in a card
	listCard := styles.CardBase.
		Render(fullListContent)

	// Use header and footer components
	headerView := m.header.View()
	footerView := m.footer.View()

	// Combine all sections
	fullContent := lipgloss.JoinVertical(
		lipgloss.Left,
		headerView,
		"",
		listCard,
		"",
		footerView,
	)

	return fullContent
}

// nextItem moves to the next item in the current page
// renderListItems renders all events as formatted strings for display
func (m *ListModel) renderListItems() []string {
	var items []string

	for i, event := range m.events {
		// Marker for selected item
		marker := "  "
		if i == m.selectedIdx {
			marker = "▶ "
		}

		// Truncate text to 100 chars
		text := event.Text
		if len(text) > 100 {
			text = text[:97] + "..."
		}

		// Determine styling based on selection
		itemStyle := styles.ListItem
		if i == m.selectedIdx {
			itemStyle = styles.ListItemSelected
		}

		// Render event text
		eventText := itemStyle.Render(marker + text)

		// Render date and company
		dateStr := event.Date.Format("2006-01-02")
		var detailLine string
		if event.Company != "" || event.Project != "" {
			detailLine = fmt.Sprintf("   %s | %s",
				styles.ListItem.Foreground(styles.ColorTextSecondary).Render(dateStr),
				styles.ListItem.Foreground(styles.ColorTextMuted).Render(event.Company),
			)
		} else {
			detailLine = styles.ListItem.Foreground(styles.ColorTextSecondary).Render(fmt.Sprintf("   %s", dateStr))
		}

		items = append(items, eventText, detailLine, "")
	}

	return items
}

// renderListWithContainer renders the list using ListContainer
func (m *ListModel) renderListWithContainer() string {
	// Render items
	items := m.renderListItems()

	// Create pagination info
	startIdx := (m.currentPage-1)*m.pageSize + 1
	endIdx := startIdx + len(m.events) - 1
	if len(m.events) == 0 {
		startIdx = 0
		endIdx = 0
	}
	paginationText := fmt.Sprintf("Showing %d-%d of %d events", startIdx, endIdx, m.totalCount)

	// Render using ListContainer
	listContainer := components.NewListContainer().
		SetItems(items).
		SetEmptyStateMessage("No events found").
		SetPaginationInfo(paginationText)

	return listContainer.Render()
}

func (m *ListModel) nextItem() {
	if m.selectedIdx < len(m.events)-1 {
		m.selectedIdx++
	}
}

// prevItem moves to the previous item in the current page
func (m *ListModel) prevItem() {
	if m.selectedIdx > 0 {
		m.selectedIdx--
	}
}

// nextPage moves to the next page
func (m *ListModel) nextPage() {
	totalPages := m.getTotalPages()
	if m.currentPage < totalPages {
		m.currentPage++
		m.selectedIdx = 0
		m.loadEvents()
	}
}

// prevPage moves to the previous page
func (m *ListModel) prevPage() {
	if m.currentPage > 1 {
		m.currentPage--
		m.selectedIdx = 0
		m.loadEvents()
	}
}

// goToFirstItem moves the selection to the first item on the current page
func (m *ListModel) goToFirstItem() {
	m.selectedIdx = 0
}

// goToLastItem moves the selection to the last item on the current page
func (m *ListModel) goToLastItem() {
	if len(m.events) > 0 {
		m.selectedIdx = len(m.events) - 1
	}
}

// getTotalPages calculates the total number of pages
func (m *ListModel) getTotalPages() int {
	if m.totalCount == 0 {
		return 1
	}
	pages := (m.totalCount + m.pageSize - 1) / m.pageSize
	return pages
}

// viewSelectedEvent returns a command to show the action menu for the selected event
func (m *ListModel) viewSelectedEvent() tea.Cmd {
	if len(m.events) == 0 || m.selectedIdx < 0 || m.selectedIdx >= len(m.events) {
		return nil
	}

	selectedEvent := m.events[m.selectedIdx]
	return func() tea.Msg {
		// Instead of directly viewing the event, send a message to show the action menu
		return EventActionMenuMsg{Event: selectedEvent}
	}
}

// applyFilters applies the filter model's filters to the event list
func (m *ListModel) applyFilters() {
	filters := m.filterModel.ToListFilters()
	filters.SortBy = "date"
	filters.SortOrder = "desc"
	filters.Limit = m.pageSize
	filters.Offset = (m.currentPage - 1) * m.pageSize

	events, err := m.service.ListEvents(m.ctx, filters)
	if err != nil {
		m.err = err
		m.events = []*career.CareerEvent{}
		return
	}

	m.events = events
	m.selectedIdx = 0

	// Get total count with filters for pagination calculation
	totalCount, err := m.service.CountEvents(m.ctx, filters)
	if err != nil {
		m.totalCount = 0
		return
	}

	m.totalCount = totalCount
	m.currentPage = 1
}

// GetFilterModel returns the filter model for this list
func (m *ListModel) GetFilterModel() *FilterModel {
	return m.filterModel
}

// GetSearchModel returns the search model for this list
func (m *ListModel) GetSearchModel() *SearchModel {
	return m.searchModel
}

// applySearch applies the search model's search query to the event list
func (m *ListModel) applySearch() {
	// First apply filter, then apply search to filtered results
	filters := m.filterModel.ToListFilters()
	filters.SortBy = "date"
	filters.SortOrder = "desc"
	filters.Limit = m.pageSize
	filters.Offset = (m.currentPage - 1) * m.pageSize

	events, err := m.service.ListEvents(m.ctx, filters)
	if err != nil {
		m.err = err
		m.events = []*career.CareerEvent{}
		return
	}

	// Apply search filter to results
	if m.searchModel.IsActive() {
		var searchResults []*career.CareerEvent
		for _, event := range events {
			if m.searchModel.Matches(event) {
				searchResults = append(searchResults, event)
			}
		}
		m.events = searchResults
	} else {
		m.events = events
	}

	m.selectedIdx = 0
	m.currentPage = 1

	// Get total count with all filters and search for pagination
	if m.searchModel.IsActive() {
		// For search, we need to count matches differently
		totalEvents, err := m.service.ListEvents(m.ctx, filters)
		if err != nil {
			m.totalCount = len(m.events)
		} else {
			// Count matches in all results
			matchCount := 0
			for _, event := range totalEvents {
				if m.searchModel.Matches(event) {
					matchCount++
				}
			}
			m.totalCount = matchCount
		}
	} else {
		// Use standard count with filters
		totalCount, err := m.service.CountEvents(m.ctx, filters)
		if err != nil {
			m.totalCount = 0
			return
		}
		m.totalCount = totalCount
	}
}

// FilterAndSearch applies both filters and search to the event list
func (m *ListModel) FilterAndSearch() {
	m.applySearch()
}

// GetSortModel returns the sort model for this list
func (m *ListModel) GetSortModel() *SortModel {
	return m.sortModel
}

// GetSelectedEvent returns the currently selected event
func (m *ListModel) GetSelectedEvent() *career.CareerEvent {
	if len(m.events) == 0 || m.selectedIdx < 0 || m.selectedIdx >= len(m.events) {
		return nil
	}
	return m.events[m.selectedIdx]
}

// SetBreadcrumbs sets breadcrumb trail for display in header
func (m *ListModel) SetBreadcrumbs(crumbs []string) {
	m.breadcrumbs = crumbs
	m.header.SetBreadcrumbs(crumbs)
}
