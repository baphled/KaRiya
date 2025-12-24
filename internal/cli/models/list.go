package models

import (
	"context"
	"fmt"
	"strings"

	"github.com/baphled/kariya/internal/domain/career"
	careerrepo "github.com/baphled/kariya/internal/repository/career"
	careerservice "github.com/baphled/kariya/internal/service/career"
	"github.com/charmbracelet/bubbletea"
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
		case "up":
			m.prevItem()
		case "down":
			m.nextItem()
		case "pgup":
			m.prevPage()
		case "pgdn":
			m.nextPage()
		}
	case tea.WindowSizeMsg:
		m.width = msg.Width
		m.height = msg.Height
	}

	return m, nil
}

// View renders the model
func (m *ListModel) View() string {
	if m.err != nil {
		return fmt.Sprintf("Error loading events: %v", m.err)
	}

	if len(m.events) == 0 {
		return "No events found. Start capturing your career journey!"
	}

	var sb strings.Builder

	// Title
	sb.WriteString("Career Events\n")
	sb.WriteString(strings.Repeat("─", 40) + "\n\n")

	// Events
	for i, event := range m.events {
		marker := "  "
		if i == m.selectedIdx {
			marker = "▶ "
		}

		// Truncate text to 100 chars
		text := event.Text
		if len(text) > 100 {
			text = text[:97] + "..."
		}

		sb.WriteString(marker)
		sb.WriteString(text)
		sb.WriteString("\n")

		// Date and company on next line
		dateStr := event.Date.Format("2006-01-02")
		if event.Company != "" {
			sb.WriteString(fmt.Sprintf("   %s | %s\n", dateStr, event.Company))
		} else {
			sb.WriteString(fmt.Sprintf("   %s\n", dateStr))
		}
		sb.WriteString("\n")
	}

	// Pagination info
	sb.WriteString(strings.Repeat("─", 40) + "\n")
	totalPages := m.getTotalPages()
	sb.WriteString(fmt.Sprintf("Page %d of %d (%d total events)\n", m.currentPage, totalPages, m.totalCount))

	return sb.String()
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

// getTotalPages calculates the total number of pages
func (m *ListModel) getTotalPages() int {
	if m.totalCount == 0 {
		return 1
	}
	pages := (m.totalCount + m.pageSize - 1) / m.pageSize
	return pages
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
