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

// BurstListModel represents the burst list display screen using table-based rendering
type BurstListModel struct {
	*BaseStandardModel
	bursts            []*career.Burst
	filtered          []*career.Burst
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
	selectedBursts    map[string]bool
	expandedIndices   map[int]bool
	err               error
	helpFooter        components.HelpFooterModel
	deletionConfirm   *ConfirmationDialog
	deletingBurstID   string
	deleteSuccessMsg  string
	deleteErrorMsg    string
	showDeleteMessage bool
	breadcrumbs       []string
	header            components.HeaderModel
	listContainer     *components.TableListContainer
}

// NewBurstListModel creates a new burst list model
func NewBurstListModel(svc *careerservice.Service, ctx context.Context) *BurstListModel {
	columns := []table.Column{
		{Title: "Burst", Width: 40},
		{Title: "Events", Width: 10},
		{Title: "Focus", Width: 25},
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

	m := &BurstListModel{
		BaseStandardModel: NewBaseStandardModel(),
		service:           svc,
		ctx:               ctx,
		bursts:            []*career.Burst{},
		filtered:          []*career.Burst{},
		table:             t,
		width:             80,
		height:            20,
		pageSize:          10,
		currentPage:       1,
		totalCount:        0,
		sortBy:            "date",
		sortOrder:         "desc",
		selectedBursts:    make(map[string]bool),
		expandedIndices:   make(map[int]bool),
		deletionConfirm:   nil,
		deletingBurstID:   "",
		deleteSuccessMsg:  "",
		deleteErrorMsg:    "",
		showDeleteMessage: false,
		helpFooter:        components.NewHelpFooter("burst_list", 80),
		breadcrumbs:       []string{"Home", "Bursts"},
		header:            components.NewHeader("💥 Bursts", 80),
		listContainer:     components.NewTableListContainer(t, "Bursts", 80),
	}
	m.loadBurstsSync()
	return m
}

// loadBurstsSync loads bursts synchronously from the service
func (m *BurstListModel) loadBurstsSync() {
	burstRepo := m.service.GetBurstRepository()
	if burstRepo == nil {
		m.bursts = []*career.Burst{}
		m.err = nil
		return
	}

	bursts, err := burstRepo.List(m.ctx, careerrepo.BurstListFilters{Limit: 1000})
	if err != nil {
		m.err = err
	} else {
		m.bursts = bursts
		m.applyFiltersAndSort()
		m.totalCount = len(m.filtered)
		m.updateTableRows()
		m.err = nil
	}
}

// updateTableRows updates the table with rows from filtered bursts
func (m *BurstListModel) updateTableRows() {
	var rows []table.Row

	for i, burst := range m.filtered {
		burstText := burst.Name
		if len(burstText) > 37 {
			burstText = burstText[:34] + "..."
		}

		eventCount := fmt.Sprintf("%d", len(burst.EventIDs))

		focus := burst.CompetencyFocus
		if len(focus) > 22 {
			focus = focus[:19] + "..."
		}

		// Add focus indicator for selected row based on container's selectedIdx
		indicator := "  "
		if i == m.listContainer.GetSelectedIdx() {
			indicator = "▶ "
		}

		row := table.Row{
			indicator + burstText,
			eventCount,
			focus,
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
	// SetRows may reset the table's cursor, so we need to explicitly set it
	m.table.SetCursor(m.listContainer.GetSelectedIdx())
}

// Init initializes the model
func (m *BurstListModel) Init() tea.Cmd {
	return nil
}

// Update handles messages
func (m *BurstListModel) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	if m.deletionConfirm != nil {
		updatedDialog, cmd := m.deletionConfirm.Update(msg)
		m.deletionConfirm = updatedDialog

		if m.deletionConfirm.IsConfirmed() {
			return m.performBurstDeletion()
		}

		if m.deletionConfirm.IsCancelled() {
			m.deletionConfirm = nil
			m.deletingBurstID = ""
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
				burst := m.GetSelectedBurst()
				if burst != nil {
					return m, func() tea.Msg { return BurstActionSelectedMsg{Burst: burst, Action: BurstActionView} }
				}
			}
			return m, nil
		case " ", "space":
			if len(m.filtered) > 0 {
				burst := m.GetSelectedBurst()
				if burst != nil {
					m.selectedBursts[burst.ID] = !m.selectedBursts[burst.ID]
				}
			}
			return m, nil
		case "x", "d":
			if len(m.filtered) > 0 {
				burst := m.GetSelectedBurst()
				if burst != nil {
					m.showDeleteConfirmation(burst)
				}
			}
			return m, nil
		case "v":
			if len(m.filtered) > 0 {
				burst := m.GetSelectedBurst()
				if burst != nil {
					return m, func() tea.Msg { return BurstActionSelectedMsg{Burst: burst, Action: BurstActionView} }
				}
			}
			return m, nil
		case "e":
			if len(m.filtered) > 0 {
				burst := m.GetSelectedBurst()
				if burst != nil {
					return m, func() tea.Msg { return BurstActionSelectedMsg{Burst: burst, Action: BurstActionEdit} }
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
func (m *BurstListModel) showDeleteConfirmation(burst *career.Burst) {
	message := fmt.Sprintf("Delete burst: \"%s\"?\n\nThis action cannot be undone.", truncateText(burst.Name, 60))
	m.deletionConfirm = NewConfirmationDialog("Delete Burst", message)
	m.deletingBurstID = burst.ID
}

// performBurstDeletion performs the actual burst deletion
func (m *BurstListModel) performBurstDeletion() (tea.Model, tea.Cmd) {
	if m.deletingBurstID == "" {
		m.deletionConfirm = nil
		return m, nil
	}

	err := m.service.DeleteBurst(m.ctx, m.deletingBurstID)
	if err != nil {
		m.deleteErrorMsg = fmt.Sprintf("Error deleting burst: %v", err)
		m.showDeleteMessage = true
		m.deletionConfirm = nil
		m.deletingBurstID = ""
		return m, nil
	}

	m.removeBurstFromLists(m.deletingBurstID)

	m.deleteSuccessMsg = "Burst deleted successfully"
	m.showDeleteMessage = true
	m.deletionConfirm = nil
	m.deletingBurstID = ""

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

// removeBurstFromLists removes a burst from both the bursts and filtered lists
func (m *BurstListModel) removeBurstFromLists(burstID string) {
	for i, burst := range m.bursts {
		if burst.ID == burstID {
			m.bursts = append(m.bursts[:i], m.bursts[i+1:]...)
			break
		}
	}

	for i, burst := range m.filtered {
		if burst.ID == burstID {
			m.filtered = append(m.filtered[:i], m.filtered[i+1:]...)
			break
		}
	}

	delete(m.selectedBursts, burstID)
}

// View renders the burst list
func (m *BurstListModel) View() string {
	if m.deletionConfirm != nil {
		return m.deletionConfirm.View()
	}

	if m.showDeleteMessage {
		return m.renderDeleteMessage()
	}

	if m.err != nil {
		errorMsg := fmt.Sprintf("Error loading bursts: %v\n\nPress 'r' to retry or 'esc' to cancel", m.err)
		m.listContainer.SetErrorMessage(errorMsg).SetDimensions(m.width, m.height)
		return m.listContainer.Render()
	}

	// Update list container with current state
	m.listContainer.SetTable(m.table).
		SetDimensions(m.width, m.height).
		SetEmptyStateMessage("No bursts found").
		SetHelpFooterKey("burst_list")

	// Set pagination info
	startIdx := (m.currentPage-1)*m.pageSize + 1
	endIdx := startIdx + len(m.getPageBursts()) - 1
	if m.totalCount == 0 {
		startIdx = 0
		endIdx = 0
	}
	paginationText := fmt.Sprintf("Showing %d-%d of %d bursts", startIdx, endIdx, m.totalCount)
	m.listContainer.SetPaginationInfo(paginationText)

	return m.listContainer.Render()
}

// renderDeleteMessage renders the deletion success/error message
func (m *BurstListModel) renderDeleteMessage() string {
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

// getPageBursts returns the bursts for the current page
func (m *BurstListModel) getPageBursts() []*career.Burst {
	if m.totalCount == 0 {
		return []*career.Burst{}
	}

	startIdx := (m.currentPage - 1) * m.pageSize
	endIdx := startIdx + m.pageSize

	if startIdx >= len(m.filtered) {
		return []*career.Burst{}
	}

	if endIdx > len(m.filtered) {
		endIdx = len(m.filtered)
	}

	return m.filtered[startIdx:endIdx]
}

// SetBursts sets the bursts to display
func (m *BurstListModel) SetBursts(bursts []*career.Burst) {
	m.bursts = bursts
	m.applyFiltersAndSort()
	m.totalCount = len(m.filtered)
	m.currentPage = 1
	m.listContainer.MoveToFirst()
	m.expandedIndices = make(map[int]bool)
	m.updateTableRows()
}

// SetCompetencyFilter sets the competency focus filter
func (m *BurstListModel) SetCompetencyFilter(competency string) {
	m.competencyFilter = competency
	m.applyFiltersAndSort()
	m.totalCount = len(m.filtered)
	m.currentPage = 1
	m.listContainer.MoveToFirst()
	m.expandedIndices = make(map[int]bool)
	m.updateTableRows()
}

// SetSort sets the sort order
func (m *BurstListModel) SetSort(sortBy, sortOrder string) {
	m.sortBy = sortBy
	m.sortOrder = sortOrder
	m.applyFiltersAndSort()
	m.updateTableRows()
}

// applyFiltersAndSort applies filters and sorts the bursts
func (m *BurstListModel) applyFiltersAndSort() {
	m.filtered = m.filterBursts()
	m.sortBursts(m.filtered)
}

// filterBursts filters bursts based on current filters
func (m *BurstListModel) filterBursts() []*career.Burst {
	var filtered []*career.Burst

	for _, burst := range m.bursts {
		if m.competencyFilter != "" {
			if !strings.EqualFold(burst.CompetencyFocus, m.competencyFilter) {
				continue
			}
		}

		filtered = append(filtered, burst)
	}

	return filtered
}

// sortBursts sorts bursts based on current sort settings
func (m *BurstListModel) sortBursts(bursts []*career.Burst) {
	sort.SliceStable(bursts, func(i, j int) bool {
		var less bool

		switch m.sortBy {
		case "name":
			less = bursts[i].Name < bursts[j].Name
		case "event_count":
			less = len(bursts[i].EventIDs) < len(bursts[j].EventIDs)
		case "date":
			fallthrough
		default:
			less = bursts[i].CreatedAt.Before(bursts[j].CreatedAt)
		}

		if m.sortOrder == "asc" {
			return less
		}
		return !less
	})
}

// GetSelectedBurst returns the currently selected burst
func (m *BurstListModel) GetSelectedBurst() *career.Burst {
	cursor := m.listContainer.GetSelectedIdx()
	if cursor >= 0 && cursor < len(m.filtered) {
		return m.filtered[cursor]
	}
	return nil
}

// IsExpandedAt returns true if the burst at the given index is expanded
func (m *BurstListModel) IsExpandedAt(idx int) bool {
	return m.expandedIndices[idx]
}

// SetBreadcrumbs sets the breadcrumb trail for navigation
func (m *BurstListModel) SetBreadcrumbs(crumbs []string) {
	m.breadcrumbs = crumbs
	m.header.SetBreadcrumbs(crumbs)
}

// Refresh reloads the bursts from the service
func (m *BurstListModel) Refresh() {
	m.loadBurstsSync()
	m.currentPage = 1
	m.expandedIndices = make(map[int]bool)
	m.listContainer.MoveToFirst()
}

// GetError returns the last error
func (m *BurstListModel) GetError() error {
	return m.err
}

// GetBursts returns the currently filtered bursts
func (m *BurstListModel) GetBursts() []*career.Burst {
	return m.filtered
}

// GetSelectedIdx returns the currently selected index
func (m *BurstListModel) GetSelectedIdx() int {
	return m.listContainer.GetSelectedIdx()
}

// getTotalPages calculates the total number of pages
func (m *BurstListModel) getTotalPages() int {
	if m.totalCount == 0 {
		return 1
	}
	pages := (m.totalCount + m.pageSize - 1) / m.pageSize
	return pages
}

// nextPage moves to the next page of results
func (m *BurstListModel) nextPage() {
	totalPages := m.getTotalPages()
	if m.currentPage < totalPages {
		m.currentPage++
		m.listContainer.MoveToFirst()
		m.updateTableRows()
	}
}

// prevPage moves to the previous page of results
func (m *BurstListModel) prevPage() {
	if m.currentPage > 1 {
		m.currentPage--
		m.listContainer.MoveToFirst()
		m.updateTableRows()
	}
}

// goToFirstItem moves to the first item
func (m *BurstListModel) goToFirstItem() {
	m.currentPage = 1
	m.listContainer.MoveToFirst()
}

// goToLastItem moves to the last item
func (m *BurstListModel) goToLastItem() {
	totalPages := m.getTotalPages()
	m.currentPage = totalPages
	pageBursts := m.getPageBursts()
	if len(pageBursts) > 0 {
		m.listContainer.SetSelectedIdx(len(pageBursts) - 1)
	}
}
