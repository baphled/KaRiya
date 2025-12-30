package models

import (
	"context"
	"fmt"
	"strings"

	"github.com/baphled/kariya/internal/cli/styles"
	"github.com/baphled/kariya/internal/domain/career"
	careerrepo "github.com/baphled/kariya/internal/repository/career"
	careerservice "github.com/baphled/kariya/internal/service/career"
	"github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
)

// ListModel represents the event list screen
type ListModel struct {
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
}

// NewListModel creates a new list model
func NewListModel(svc *careerservice.Service, ctx context.Context) *ListModel {
	model := &ListModel{
		service:     svc,
		ctx:         ctx,
		pageSize:    10,
		currentPage: 1,
		selectedIdx: 0,
		filterModel: NewFilterModel(),
		searchModel: NewSearchModel(),
		sortModel:   NewSortModel(),
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
	case tea.KeyMsg:
		switch msg.String() {
		case "backspace":
			// Signal back navigation to parent
			return m, func() tea.Msg { return BackMsg{} }
		case "ctrl+c", "q", "esc":
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
	case tea.WindowSizeMsg:
		m.width = msg.Width
		m.height = msg.Height
	}

	return m, nil
}

// View renders the model
func (m *ListModel) View() string {
	// Title
	title := styles.HeaderMain.Render("Career Events")

	// Form content
	var content []string

	// Error handling
	if m.err != nil {
		content = append(content, styles.ErrorText.Render(fmt.Sprintf("Error loading events: %v", m.err)))
	}

	// Empty state
	if len(m.events) == 0 {
		content = append(content, styles.InfoText.Render("No events found. Start capturing your career journey!"))
	} else {
		// Events
		for i, event := range m.events {
			// Determine styling based on selection
			itemStyle := styles.ListItem
			if i == m.selectedIdx {
				itemStyle = styles.ListItemSelected
			}

			// Truncate text to 100 chars
			text := event.Text
			if len(text) > 100 {
				text = text[:97] + "..."
			}

			// Marker for selected item
			marker := "  "
			if i == m.selectedIdx {
				marker = "▶ "
			}

			// Render event text
			content = append(content,
				itemStyle.Render(marker + text),
			)

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
			content = append(content, detailLine, "")
		}

		// Pagination info
		totalPages := m.getTotalPages()
		paginationText := fmt.Sprintf("Page %d of %d (%d total events)", m.currentPage, totalPages, m.totalCount)
		content = append(content,
			styles.ListItem.Render(strings.Repeat("─", 40)),
			styles.ListItem.Foreground(styles.ColorTextSecondary).Render(paginationText),
		)

		// Instructions
		instructions := "↑/↓ or j/k: Navigate items | PgUp/PgDn or Ctrl+F/B: Change pages | g/G: First/Last | Enter: View details"
		content = append(content,
			styles.InfoText.Render(instructions),
		)
	}

	// Combine all content
	listContent := lipgloss.JoinVertical(
		lipgloss.Left,
		content...,
	)

	// Wrap in a card
	listCard := styles.CardBase.
		Render(listContent)

	// Combine all sections
	fullContent := lipgloss.JoinVertical(
		lipgloss.Left,
		title,
		"",
		listCard,
	)

	return fullContent
}

// nextItem moves to the next item in the current page
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
