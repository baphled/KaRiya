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
	displayedBursts := m.getDisplayedBursts()

	var items []string
	if len(displayedBursts) > 0 {
		for i, burstIdx := range displayedBursts {
			if i >= m.height-5 {
				break
			}

			burst := m.bursts[burstIdx]
			items = append(items, m.renderBurstRow(burst, i == m.selectedIdx))

			if m.expandedIndices[burstIdx] {
				items = append(items, m.renderExpandedEvents(burst))
			}
		}
	}

	var emptyStateMessage string
	if m.filterBy != "" {
		emptyStateMessage = "No matching bursts"
	} else {
		emptyStateMessage = "No bursts found"
	}

	// Create pagination info - standard format for all list models
	paginationInfo := fmt.Sprintf("Showing %d of %d bursts", len(displayedBursts), len(m.bursts))

	listContainer := components.NewListContainer().
		SetItems(items).
		SetEmptyStateMessage(emptyStateMessage).
		SetPaginationInfo(paginationInfo)

	listContent := listContainer.Render()

	headerView := components.NewHeader("💥 Bursts", m.width).View()
	footerView := components.NewFooter(m.width).View()

	screenContent := lipgloss.JoinVertical(
		lipgloss.Left,
		m.renderHeader(),
		"",
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

// renderHeader renders the column headers
func (m *BurstListModel) renderHeader() string {
	nameCol := lipgloss.NewStyle().Width(30).Render("Name")
	countCol := lipgloss.NewStyle().Width(8).Render("Events")
	compCol := lipgloss.NewStyle().Width(15).Render("Competency")
	dateCol := lipgloss.NewStyle().Width(12).Render("Created")

	header := nameCol + " " + countCol + " " + compCol + " " + dateCol
	return styles.HeaderSection.Render(header)
}

// renderBurstRow renders a single burst row
func (m *BurstListModel) renderBurstRow(burst *career.Burst, isSelected bool) string {
	name := truncateString(burst.Name, 28)
	count := fmt.Sprintf("%d", len(burst.EventIDs))
	competency := truncateString(burst.CompetencyFocus, 13)
	date := burst.CreatedAt.Format("2006-01-02")

	nameCol := lipgloss.NewStyle().Width(30).Render(name)
	countCol := lipgloss.NewStyle().Width(8).Render(count)
	compCol := lipgloss.NewStyle().Width(15).Render(competency)
	dateCol := lipgloss.NewStyle().Width(12).Render(date)

	row := nameCol + " " + countCol + " " + compCol + " " + dateCol

	if isSelected {
		row = styles.ListItemSelected.Render(row)
	} else {
		row = styles.ListItem.Render(row)
	}

	return row + "\n"
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
