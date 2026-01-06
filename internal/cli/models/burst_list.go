package models

import (
	"context"
	"fmt"
	"sort"

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
	bursts               []*career.Burst
	filtered             []*career.Burst
	service              *careerservice.Service
	ctx                  context.Context
	table                table.Model
	pagination           *PaginationHelper
	width                int
	height               int
	sortBy               string
	sortOrder            string
	selectedBursts       map[string]bool
	expandedIndices      map[int]bool
	err                  error
	helpFooter           components.HelpFooterModel
	deletionState        *ListDeletionState
	navigationKeyHandler *ListNavigationKeyHandler
	breadcrumbs          []string
	header               components.HeaderModel
	listContainer        *components.TableListContainer
}

// NewBurstListModel creates a new burst list model
func NewBurstListModel(svc *careerservice.Service, ctx context.Context) *BurstListModel {
	columns := []table.Column{
		{Title: "Burst", Width: 40},
		{Title: "Events", Width: 10},
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
		pagination:        NewPaginationHelper(10),
		width:             80,
		height:            20,
		sortBy:            "date",
		sortOrder:         "desc",
		selectedBursts:    make(map[string]bool),
		expandedIndices:   make(map[int]bool),
		deletionState:     NewListDeletionState(),
		helpFooter:        components.NewHelpFooter("burst_list", 80),
		breadcrumbs:       []string{"Home", "Bursts"},
		header:            components.NewHeader("💥 Bursts", 80),
		listContainer:     components.NewTableListContainer(t, "Bursts", 80),
	}
	// Create navigation key handler with callbacks
	m.navigationKeyHandler = NewListNavigationKeyHandler(m)
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
		m.pagination.SetTotalCount(len(m.filtered))
		m.updateTableRows()
		m.err = nil
	}
}

// updateTableRows updates the table with rows from filtered bursts (current page only)
func (m *BurstListModel) updateTableRows() {
	var rows []table.Row
	pageBursts := m.getPageBursts()

	// Get the selected index within the current page
	selectedIdx := m.listContainer.GetSelectedIdx()

	// Ensure cursor is within valid bounds BEFORE creating rows
	if len(pageBursts) > 0 {
		if selectedIdx >= len(pageBursts) {
			selectedIdx = len(pageBursts) - 1
		}
		if selectedIdx < 0 {
			selectedIdx = 0
		}
	} else {
		selectedIdx = 0
	}

	for i, burst := range pageBursts {
		burstText := burst.Name
		if len(burstText) > 40 {
			burstText = burstText[:37] + "..."
		}

		// Add focus indicator for the selected row, or spaces for alignment
		if i == selectedIdx {
			burstText = "▶ " + burstText
		} else {
			burstText = "  " + burstText
		}

		eventCount := fmt.Sprintf("%d", len(burst.EventIDs))

		row := table.Row{
			burstText,
			eventCount,
		}
		rows = append(rows, row)
	}

	m.table.SetRows(rows)
	// Sync the container and table cursor with the validated selectedIdx
	m.listContainer.SetSelectedIdx(selectedIdx)
	m.table.SetCursor(selectedIdx)
}

// Init initializes the model
func (m *BurstListModel) Init() tea.Cmd {
	return nil
}

// Update handles messages
func (m *BurstListModel) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	// Handle deletion confirmation if active
	if m.deletionState.IsConfirming() {
		cmd := m.deletionState.UpdateConfirmation(msg)

		if m.deletionState.IsConfirmed() {
			return m.performBurstDeletion()
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

// performBurstDeletion performs the actual burst deletion
func (m *BurstListModel) performBurstDeletion() (tea.Model, tea.Cmd) {
	if m.deletionState.DeletingItemID == "" {
		m.deletionState.Clear()
		return m, nil
	}

	err := m.service.DeleteBurst(m.ctx, m.deletionState.DeletingItemID)
	if err != nil {
		m.deletionState.SetErrorMsg(fmt.Sprintf("Error deleting burst: %v", err))
		return m, nil
	}

	m.removeBurstFromLists(m.deletionState.DeletingItemID)

	m.deletionState.SetSuccessMsg("Burst deleted successfully")

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
	// Render deletion confirmation dialog
	if m.deletionState.IsConfirming() {
		return m.deletionState.ConfirmationDialog.View()
	}

	// Render deletion message (success or error)
	if m.deletionState.ShowMessage {
		return m.renderDeletionMessage()
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
		SetHelpFooterKey("burst_list").SetBreadcrumbs(m.breadcrumbs)

	// Set pagination info
	paginationText := m.pagination.GetPaginationInfo("bursts")
	m.listContainer.SetPaginationInfo(paginationText)

	return m.listContainer.Render()
}

// renderDeletionMessage renders the deletion success/error message
func (m *BurstListModel) renderDeletionMessage() string {
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

// getPageBursts returns the bursts for the current page
func (m *BurstListModel) getPageBursts() []*career.Burst {
	if m.pagination.GetTotalCount() == 0 {
		return []*career.Burst{}
	}

	startIdx := m.pagination.GetPageStartIndex()
	endIdx := m.pagination.GetPageEndIndex()

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
	m.pagination.SetTotalCount(len(m.filtered)).GoToFirstPage()
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
	m.pagination.GoToFirstPage()
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

// ListItemCallbacks implementation for BurstListModel
// These methods implement the ListItemCallbacks interface required by ListNavigationKeyHandler

// OnView sends a message to view the selected burst
func (m *BurstListModel) OnView(item interface{}) tea.Cmd {
	if burst, ok := item.(*career.Burst); ok {
		return func() tea.Msg { return BurstActionSelectedMsg{Burst: burst, Action: BurstActionView} }
	}
	return nil
}

// OnEdit sends a message to edit the selected burst
func (m *BurstListModel) OnEdit(item interface{}) tea.Cmd {
	if burst, ok := item.(*career.Burst); ok {
		return func() tea.Msg { return BurstActionSelectedMsg{Burst: burst, Action: BurstActionEdit} }
	}
	return nil
}

// OnDelete initiates deletion of the selected burst
func (m *BurstListModel) OnDelete(item interface{}) {
	if burst, ok := item.(*career.Burst); ok {
		m.deletionState.ShowConfirmation("Burst", burst.Name)
		m.deletionState.DeletingItemID = burst.ID
	}
}

// OnToggleSelection toggles the selection state of an item
func (m *BurstListModel) OnToggleSelection(item interface{}) {
	if burst, ok := item.(*career.Burst); ok {
		m.selectedBursts[burst.ID] = !m.selectedBursts[burst.ID]
	}
}

// HasSelectedItem returns the currently selected burst
func (m *BurstListModel) HasSelectedItem() interface{} {
	return m.GetSelectedBurst()
}

// MoveUp moves the selection up by count items
func (m *BurstListModel) MoveUp(count int) {
	m.listContainer.MoveUp(count)
}

// MoveDown moves the selection down by count items
func (m *BurstListModel) MoveDown(count int) {
	m.listContainer.MoveDown(count)
}

// MoveToFirst moves selection to the first item
func (m *BurstListModel) MoveToFirst() {
	m.listContainer.MoveToFirst()
}

// MoveToLast moves selection to the last item
func (m *BurstListModel) MoveToLast() {
	m.listContainer.MoveToLast()
}

// UpdateDisplay refreshes the table display
func (m *BurstListModel) UpdateDisplay() {
	m.updateTableRows()
}

// GetRowCount returns the number of visible rows
func (m *BurstListModel) GetRowCount() int {
	return len(m.table.Rows())
}

// GetCurrentIndex returns the current selection index
func (m *BurstListModel) GetCurrentIndex() int {
	return m.listContainer.GetSelectedIdx()
}

// GetPageSize returns the page size for pagination
func (m *BurstListModel) GetPageSize() int {
	return m.pagination.GetPageSize()
}

// nextPage moves to the next page of results
func (m *BurstListModel) nextPage() {
	m.pagination.NextPage()
	m.listContainer.MoveToFirst()
	m.updateTableRows()
}

// prevPage moves to the previous page of results
func (m *BurstListModel) prevPage() {
	m.pagination.PrevPage()
	m.listContainer.MoveToFirst()
	m.updateTableRows()
}

// goToFirstItem moves to the first item
func (m *BurstListModel) goToFirstItem() {
	m.pagination.GoToFirstPage()
	m.listContainer.MoveToFirst()
}

// goToLastItem moves to the last item
func (m *BurstListModel) goToLastItem() {
	m.pagination.GoToLastPage()
	pageBursts := m.getPageBursts()
	if len(pageBursts) > 0 {
		m.listContainer.SetSelectedIdx(len(pageBursts) - 1)
	}
}

// NextPage implements ListItemCallbacks interface - moves to the next page
func (m *BurstListModel) NextPage() {
	m.nextPage()
}

// PrevPage implements ListItemCallbacks interface - moves to the previous page
func (m *BurstListModel) PrevPage() {
	m.prevPage()
}

// GoToFirstPage implements ListItemCallbacks interface - goes to the first page
func (m *BurstListModel) GoToFirstPage() {
	m.goToFirstItem()
	m.updateTableRows()
}

// GoToLastPage implements ListItemCallbacks interface - goes to the last page
func (m *BurstListModel) GoToLastPage() {
	m.goToLastItem()
	m.updateTableRows()
}
