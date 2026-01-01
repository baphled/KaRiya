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
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
)

// BurstsLoadedMsg is sent when bursts have been loaded from repository
type BurstsLoadedMsg struct {
	Bursts []*career.Burst
	Err    error
}

// BurstListModel represents the burst list display screen
type BurstListModel struct {
	*BaseStandardModel
	service         *careerservice.Service
	ctx             context.Context
	bursts          []*career.Burst
	selectedIdx     int
	expandedIndices map[int]bool
	filterBy        string
	sortBy          string
	width           int
	height          int
}

// NewBurstListModel creates a new burst list model
func NewBurstListModel(svc *careerservice.Service, ctx context.Context) *BurstListModel {
	return &BurstListModel{
		BaseStandardModel: NewBaseStandardModel(),
		service:           svc,
		ctx:               ctx,
		bursts:            []*career.Burst{},
		selectedIdx:       0,
		expandedIndices:   make(map[int]bool),
		filterBy:          "",
		sortBy:            "date",
		width:             80,
		height:            24,
	}
}

// Init initializes the model
func (m *BurstListModel) Init() tea.Cmd {
	return m.loadBursts()
}

// loadBursts loads bursts from the service
func (m *BurstListModel) loadBursts() tea.Cmd {
	return func() tea.Msg {
		burstRepo := m.service.GetBurstRepository()
		if burstRepo == nil {
			return BurstsLoadedMsg{Bursts: []*career.Burst{}, Err: fmt.Errorf("burst repository not configured")}
		}

		bursts, err := burstRepo.List(m.ctx, careerrepo.BurstListFilters{Limit: 1000})
		return BurstsLoadedMsg{Bursts: bursts, Err: err}
	}
}

// Update handles messages
func (m *BurstListModel) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.WindowSizeMsg:
		m.width = msg.Width
		m.height = msg.Height
		return m, nil

	case BurstsLoadedMsg:
		if msg.Err == nil {
			m.bursts = msg.Bursts
		}
		return m, nil

	case tea.KeyMsg:
		switch msg.String() {
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
		case " ", "space":
			m.expandedIndices[m.selectedIdx] = !m.expandedIndices[m.selectedIdx]
		case "esc":
			return m, func() tea.Msg { return BackMsg{} }
		case "q", "ctrl+c":
			return m, func() tea.Msg { return QuitMsg{} }
		case "enter":
			m.expandedIndices[m.selectedIdx] = !m.expandedIndices[m.selectedIdx]
		}
	}

	return m, nil
}

// View renders the burst list
func (m *BurstListModel) View() string {
	// Render items using list.go pattern
	items := m.renderListItems()

	// Simplified empty state message - matching list.go
	emptyStateMessage := "No bursts found"

	// Create pagination info - matching list.go format: "Showing X-Y of Z bursts"
	displayedBursts := m.getDisplayedBursts()
	startIdx := 1
	endIdx := len(displayedBursts)
	if endIdx > m.height-5 {
		endIdx = m.height - 5
	}
	paginationInfo := fmt.Sprintf("Showing %d-%d of %d bursts", startIdx, endIdx, len(m.bursts))

	listContainer := components.NewListContainer().
		SetItems(items).
		SetEmptyStateMessage(emptyStateMessage).
		SetPaginationInfo(paginationInfo)

	listContent := listContainer.Render()

	headerView := components.NewHeader("💥 Bursts", m.width).View()
	footerView := components.NewFooter(m.width).View()

	screenContent := lipgloss.JoinVertical(
		lipgloss.Left,
		listContent,
	)

	screenContainer := components.NewScreenContainer(screenContent).
		WithPaddingMode(components.PaddingNormal)

	fullContent := lipgloss.JoinVertical(
		lipgloss.Left,
		headerView,
		"",
		screenContainer.Render(),
		"",
		footerView,
	)

	return fullContent
}

// renderListItems renders burst items following list.go's pattern
func (m *BurstListModel) renderListItems() []string {
	var items []string

	displayedBursts := m.getDisplayedBursts()
	maxDisplay := m.height - 5
	if len(displayedBursts) < maxDisplay {
		maxDisplay = len(displayedBursts)
	}

	for i := 0; i < maxDisplay; i++ {
		burstIdx := displayedBursts[i]
		burst := m.bursts[burstIdx]

		// Marker for selected item (matching list.go pattern)
		marker := "  "
		if i == m.selectedIdx {
			marker = "▶ "
		}

		// Truncate burst name to 100 chars (matching list.go pattern)
		text := burst.Name
		if len(text) > 100 {
			text = text[:97] + "..."
		}

		// Determine styling based on selection (matching list.go pattern)
		itemStyle := styles.ListItem
		if i == m.selectedIdx {
			itemStyle = styles.ListItemSelected
		}

		// Render burst name
		burstText := itemStyle.Render(marker + text)

		// Render count, competency, and date as detail line
		count := fmt.Sprintf("%d events", len(burst.EventIDs))
		competency := truncateString(burst.CompetencyFocus, 20)
		date := burst.CreatedAt.Format("2006-01-02")
		detailLine := fmt.Sprintf("   %s | %s | %s",
			styles.ListItem.Foreground(styles.ColorTextSecondary).Render(count),
			styles.ListItem.Foreground(styles.ColorTextMuted).Render(competency),
			styles.ListItem.Foreground(styles.ColorTextSecondary).Render(date),
		)

		items = append(items, burstText, detailLine, "")

		// Add expanded events if this burst is expanded
		if m.expandedIndices[burstIdx] {
			items = append(items, m.renderExpandedEvents(burst))
		}
	}

	return items
}

// renderExpandedEvents renders the event IDs for an expanded burst
func (m *BurstListModel) renderExpandedEvents(burst *career.Burst) string {
	var sb strings.Builder
	sb.WriteString("  ")
	sb.WriteString(styles.InfoHint.Render("Events: "))

	for i, eventID := range burst.EventIDs {
		if i > 0 {
			sb.WriteString(", ")
		}
		truncated := truncateString(eventID, 8)
		sb.WriteString(styles.InfoText.Render(truncated))
	}

	sb.WriteString("\n")
	return sb.String()
}

// getDisplayedBursts returns the filtered and sorted bursts as indices
func (m *BurstListModel) getDisplayedBursts() []int {
	var displayed []int

	for i, burst := range m.bursts {
		if m.filterBy == "" || burst.CompetencyFocus == m.filterBy {
			displayed = append(displayed, i)
		}
	}

	m.sortBursts(displayed)

	return displayed
}

// sortBursts sorts the burst indices based on current sort order
func (m *BurstListModel) sortBursts(indices []int) {
	switch m.sortBy {
	case "name":
		sort.Slice(indices, func(i, j int) bool {
			return m.bursts[indices[i]].Name < m.bursts[indices[j]].Name
		})

	case "event_count":
		sort.Slice(indices, func(i, j int) bool {
			return len(m.bursts[indices[i]].EventIDs) < len(m.bursts[indices[j]].EventIDs)
		})

	case "date":
		fallthrough
	default:
		sort.Slice(indices, func(i, j int) bool {
			return m.bursts[indices[i]].CreatedAt.Before(m.bursts[indices[j]].CreatedAt)
		})
	}
}

// truncateString truncates a string to a maximum width, adding "..." if truncated
func truncateString(s string, maxWidth int) string {
	if len(s) <= maxWidth {
		return s
	}
	if maxWidth <= 3 {
		return "..."
	}
	return s[:maxWidth-3] + "..."
}

// SetBursts sets the bursts to display
func (m *BurstListModel) SetBursts(bursts []*career.Burst) {
	m.bursts = bursts
	m.selectedIdx = 0
	m.expandedIndices = make(map[int]bool)
}

// GetSelectedBurst returns the currently selected burst
func (m *BurstListModel) GetSelectedBurst() *career.Burst {
	displayedBursts := m.getDisplayedBursts()
	if len(displayedBursts) == 0 || m.selectedIdx >= len(displayedBursts) {
		return nil
	}
	return m.bursts[displayedBursts[m.selectedIdx]]
}

// SetFilter sets the competency focus filter
func (m *BurstListModel) SetFilter(competency string) {
	m.filterBy = competency
	m.selectedIdx = 0
}

// SetSort sets the sort order
func (m *BurstListModel) SetSort(sortBy string) {
	m.sortBy = sortBy
}

// nextItem moves to the next item in the list
func (m *BurstListModel) nextItem() {
	displayedBursts := m.getDisplayedBursts()
	if len(displayedBursts) > 0 {
		m.selectedIdx = (m.selectedIdx + 1) % len(displayedBursts)
	}
}

// prevItem moves to the previous item in the list
func (m *BurstListModel) prevItem() {
	displayedBursts := m.getDisplayedBursts()
	if len(displayedBursts) > 0 {
		m.selectedIdx = (m.selectedIdx - 1 + len(displayedBursts)) % len(displayedBursts)
	}
}

// nextPage moves to the next page
func (m *BurstListModel) nextPage() {
	displayedBursts := m.getDisplayedBursts()
	pageSize := m.height - 5
	newIdx := m.selectedIdx + pageSize
	if newIdx >= len(displayedBursts) {
		newIdx = len(displayedBursts) - 1
	}
	m.selectedIdx = newIdx
}

// prevPage moves to the previous page
func (m *BurstListModel) prevPage() {
	pageSize := m.height - 5
	newIdx := m.selectedIdx - pageSize
	if newIdx < 0 {
		newIdx = 0
	}
	m.selectedIdx = newIdx
}

// goToFirstItem moves to the first item in the list
func (m *BurstListModel) goToFirstItem() {
	m.selectedIdx = 0
}

// goToLastItem moves to the last item in the list
func (m *BurstListModel) goToLastItem() {
	displayedBursts := m.getDisplayedBursts()
	if len(displayedBursts) > 0 {
		m.selectedIdx = len(displayedBursts) - 1
	}
}
