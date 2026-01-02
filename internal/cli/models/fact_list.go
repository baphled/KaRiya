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

// FactListModel displays a list of facts with filtering, sorting, and selection using bubbles table
type FactListModel struct {
	*BaseStandardModel
	facts             []*career.Fact
	filtered          []*career.Fact
	service           *careerservice.Service
	ctx               context.Context
	table             table.Model
	pagination        *PaginationHelper
	width             int
	height            int
	competencyFilter  string
	roleFitFilter     career.RoleFit
	audienceFilter    string
	sortBy            string
	sortOrder         string
	selectedFacts     map[string]bool
	submitted         bool
	cancelled         bool
	err               error
	helpFooter        components.HelpFooterModel
	deletionState     *ListDeletionState
	navigationKeyHandler *ListNavigationKeyHandler
	breadcrumbs       []string
	header            components.HeaderModel
	listContainer     *components.TableListContainer
}

// NewFactListModel creates a new fact list model
func NewFactListModel(service *careerservice.Service, ctx context.Context) *FactListModel {
	columns := []table.Column{
		{Title: "Fact", Width: 60},
		{Title: "Competencies", Width: 25},
		{Title: "Role Fit", Width: 12},
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

	m := &FactListModel{
		BaseStandardModel: NewBaseStandardModel(),
		service:           service,
		ctx:               ctx,
		facts:             []*career.Fact{},
		filtered:          []*career.Fact{},
		table:             t,
		selectedFacts:     make(map[string]bool),
		pagination:        NewPaginationHelper(10),
		width:             80,
		height:            20,
		sortBy:            "date",
		sortOrder:         "desc",
		helpFooter:        components.NewHelpFooter("fact_list", 80),
		deletionState:     NewListDeletionState(),
		breadcrumbs:       []string{"Home", "Facts"},
		header:            components.NewHeader("Facts", 80),
		listContainer:     components.NewTableListContainer(t, "Facts", 80),
	}
	// Create navigation key handler with callbacks
	m.navigationKeyHandler = NewListNavigationKeyHandler(m)
	m.loadFactsSync()
	return m
}

// loadFactsSync loads facts synchronously from the service
func (flm *FactListModel) loadFactsSync() {
	factRepo := flm.service.GetFactRepository()
	if factRepo == nil {
		flm.facts = []*career.Fact{}
		flm.err = nil
		return
	}

	facts, err := factRepo.List(flm.ctx, careerrepo.FactListFilters{Limit: 1000})
	if err != nil {
		flm.err = err
	} else {
		flm.facts = facts
		flm.applyFiltersAndSort()
		flm.pagination.SetTotalCount(len(flm.filtered))
		flm.updateTableRows()
		flm.err = nil
	}
}

// updateTableRows updates the table with rows from filtered facts (current page only)
func (flm *FactListModel) updateTableRows() {
	var rows []table.Row
	pageFacts := flm.getPageFacts()

	// Get the selected index within the current page
	selectedIdx := flm.listContainer.GetSelectedIdx()

	// Ensure cursor is within valid bounds BEFORE creating rows
	if len(pageFacts) > 0 {
		if selectedIdx >= len(pageFacts) {
			selectedIdx = len(pageFacts) - 1
		}
		if selectedIdx < 0 {
			selectedIdx = 0
		}
	} else {
		selectedIdx = 0
	}

	for i, fact := range pageFacts {
		factText := fact.Text
		if len(factText) > 60 {
			factText = factText[:57] + "..."
		}

		// Add focus indicator for the selected row, or spaces for alignment
		if i == selectedIdx {
			factText = "▶ " + factText
		} else {
			factText = "  " + factText
		}

		roleIcon := getRoleFitIcon(fact.RoleFit)

		competencies := strings.Join(fact.CompetencyCategories, ", ")
		if len(competencies) > 22 {
			competencies = competencies[:19] + "..."
		}

		row := table.Row{
			factText,
			competencies,
			roleIcon,
		}
		rows = append(rows, row)
	}

	flm.table.SetRows(rows)
	// Sync the container and table cursor with the validated selectedIdx
	flm.listContainer.SetSelectedIdx(selectedIdx)
	flm.table.SetCursor(selectedIdx)
}

// Init initializes the model
func (flm *FactListModel) Init() tea.Cmd {
	return nil
}

// Update handles messages
func (flm *FactListModel) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	// Handle deletion confirmation if active
	if flm.deletionState.IsConfirming() {
		cmd := flm.deletionState.UpdateConfirmation(msg)

		if flm.deletionState.IsConfirmed() {
			return flm.performFactDeletion()
		}

		if flm.deletionState.IsCancelled() {
			flm.deletionState.Clear()
			return flm, nil
		}

		return flm, cmd
	}

	// Handle deletion message display
	if flm.deletionState.ShowMessage {
		switch msg := msg.(type) {
		case tea.KeyMsg:
			if msg.String() == "esc" {
				flm.deletionState.Clear()
				return flm, nil
			}
		}
	}

	switch msg := msg.(type) {
	case tea.KeyMsg:
		// Let navigation key handler process navigation and item action keys
		if cmd := flm.navigationKeyHandler.HandleNavigationKey(msg.String()); cmd != nil {
			return flm, cmd
		}

		// Handle screen navigation keys that aren't list-specific
		switch msg.String() {
		case "esc":
			return flm, func() tea.Msg { return BackMsg{} }
		case "q", "ctrl+c":
			return flm, func() tea.Msg { return QuitMsg{} }
		}
	case tea.WindowSizeMsg:
		flm.width = msg.Width
		flm.height = msg.Height
		flm.table.SetWidth(flm.width)
		flm.table.SetHeight(flm.height - 10)
		flm.listContainer.SetDimensions(flm.width, flm.height)
		// Recalculate rows with new dimensions to ensure emojis render properly
		flm.updateTableRows()
	}

	var cmd tea.Cmd
	flm.table, cmd = flm.table.Update(msg)
	return flm, cmd
}

// performFactDeletion performs the actual fact deletion
func (flm *FactListModel) performFactDeletion() (tea.Model, tea.Cmd) {
	if flm.deletionState.DeletingItemID == "" {
		flm.deletionState.Clear()
		return flm, nil
	}

	err := flm.service.DeleteFact(flm.ctx, flm.deletionState.DeletingItemID)
	if err != nil {
		flm.deletionState.SetErrorMsg(fmt.Sprintf("Error deleting fact: %v", err))
		return flm, nil
	}

	flm.removeFactFromLists(flm.deletionState.DeletingItemID)

	flm.deletionState.SetSuccessMsg("Fact deleted successfully")

	flm.pagination.SetTotalCount(len(flm.filtered))
	flm.updateTableRows()

	if flm.pagination.GetTotalCount() == 0 {
		flm.listContainer.MoveToFirst()
	} else {
		if flm.listContainer.GetSelectedIdx() >= len(flm.table.Rows()) && flm.listContainer.GetSelectedIdx() > 0 {
			flm.listContainer.SetSelectedIdx(flm.listContainer.GetSelectedIdx() - 1)
		}
	}

	return flm, nil
}

// removeFactFromLists removes a fact from both the facts and filtered lists
func (flm *FactListModel) removeFactFromLists(factID string) {
	for i, fact := range flm.facts {
		if fact.ID == factID {
			flm.facts = append(flm.facts[:i], flm.facts[i+1:]...)
			break
		}
	}

	for i, fact := range flm.filtered {
		if fact.ID == factID {
			flm.filtered = append(flm.filtered[:i], flm.filtered[i+1:]...)
			break
		}
	}

	delete(flm.selectedFacts, factID)
}

// View renders the fact list
func (flm *FactListModel) View() string {
	// Render deletion confirmation dialog
	if flm.deletionState.IsConfirming() {
		return flm.deletionState.ConfirmationDialog.View()
	}

	// Render deletion message (success or error)
	if flm.deletionState.ShowMessage {
		return flm.renderDeletionMessage()
	}

	if flm.err != nil {
		errorMsg := fmt.Sprintf("Error loading facts: %v\n\nPress 'r' to retry or 'esc' to cancel", flm.err)
		flm.listContainer.SetErrorMessage(errorMsg).SetDimensions(flm.width, flm.height)
		return flm.listContainer.Render()
	}

	// Update list container with current state
	flm.listContainer.SetTable(flm.table).
		SetDimensions(flm.width, flm.height).
		SetEmptyStateMessage("No facts found").
		SetHelpFooterKey("fact_list").SetBreadcrumbs(flm.breadcrumbs)

	// Set pagination info
	paginationText := flm.pagination.GetPaginationInfo("facts")
	flm.listContainer.SetPaginationInfo(paginationText)

	return flm.listContainer.Render()
}

// renderDeletionMessage renders the deletion success/error message
func (flm *FactListModel) renderDeletionMessage() string {
	var messageContent string
	if flm.deletionState.SuccessMsg != "" {
		messageContent = styles.SuccessBox.Render(flm.deletionState.SuccessMsg + "\n\nPress 'esc' to continue")
	} else if flm.deletionState.ErrorMsg != "" {
		messageContent = styles.ErrorBox.Render(flm.deletionState.ErrorMsg + "\n\nPress 'esc' to continue")
	}

	headerView := flm.header.View()
	footerView := components.NewFooter(flm.width).View()

	return lipgloss.JoinVertical(
		lipgloss.Left,
		headerView,
		"",
		messageContent,
		"",
		footerView,
	)
}

// getPageFacts returns the facts for the current page
func (flm *FactListModel) getPageFacts() []*career.Fact {
	if flm.pagination.GetTotalCount() == 0 {
		return []*career.Fact{}
	}

	startIdx := flm.pagination.GetPageStartIndex()
	endIdx := flm.pagination.GetPageEndIndex()

	if startIdx >= len(flm.filtered) {
		return []*career.Fact{}
	}

	if endIdx > len(flm.filtered) {
		endIdx = len(flm.filtered)
	}

	return flm.filtered[startIdx:endIdx]
}

// SetFacts sets the facts to display
func (flm *FactListModel) SetFacts(facts []*career.Fact) {
	flm.facts = facts
	flm.applyFiltersAndSort()
	flm.pagination.SetTotalCount(len(flm.filtered)).GoToFirstPage()
	flm.listContainer.MoveToFirst()
	flm.updateTableRows()
}

// SetCompetencyFilter sets the competency filter
func (flm *FactListModel) SetCompetencyFilter(competency string) {
	flm.competencyFilter = competency
	flm.applyFiltersAndSort()
	flm.pagination.SetTotalCount(len(flm.filtered)).GoToFirstPage()
	flm.listContainer.MoveToFirst()
	flm.updateTableRows()
}

// SetRoleFitFilter sets the role fit filter
func (flm *FactListModel) SetRoleFitFilter(roleFit career.RoleFit) {
	flm.roleFitFilter = roleFit
	flm.applyFiltersAndSort()
	flm.pagination.SetTotalCount(len(flm.filtered)).GoToFirstPage()
	flm.listContainer.MoveToFirst()
	flm.updateTableRows()
}

// SetAudienceFilter sets the audience filter
func (flm *FactListModel) SetAudienceFilter(audience string) {
	flm.audienceFilter = audience
	flm.applyFiltersAndSort()
	flm.pagination.SetTotalCount(len(flm.filtered)).GoToFirstPage()
	flm.listContainer.MoveToFirst()
	flm.updateTableRows()
}

// SetSort sets the sort order
func (flm *FactListModel) SetSort(sortBy, sortOrder string) {
	flm.sortBy = sortBy
	flm.sortOrder = sortOrder
	flm.applyFiltersAndSort()
	flm.updateTableRows()
}

// applyFiltersAndSort applies filters and sorts the facts
func (flm *FactListModel) applyFiltersAndSort() {
	flm.filtered = flm.filterFacts()
	flm.sortFacts(flm.filtered)
}

// filterFacts filters facts based on current filters
func (flm *FactListModel) filterFacts() []*career.Fact {
	var filtered []*career.Fact

	for _, fact := range flm.facts {
		if flm.competencyFilter != "" {
			hasCompetency := false
			for _, comp := range fact.CompetencyCategories {
				if strings.EqualFold(comp, flm.competencyFilter) {
					hasCompetency = true
					break
				}
			}
			if !hasCompetency {
				continue
			}
		}

		if flm.roleFitFilter != "" && fact.RoleFit != flm.roleFitFilter {
			continue
		}

		if flm.audienceFilter != "" {
			hasAudience := false
			for _, aud := range fact.AudienceRelevance {
				if strings.EqualFold(aud, flm.audienceFilter) {
					hasAudience = true
					break
				}
			}
			if !hasAudience {
				continue
			}
		}

		filtered = append(filtered, fact)
	}

	return filtered
}

// sortFacts sorts facts based on current sort settings
func (flm *FactListModel) sortFacts(facts []*career.Fact) {
	sort.SliceStable(facts, func(i, j int) bool {
		var less bool

		switch flm.sortBy {
		case "date":
			less = facts[i].CreatedAt.Before(facts[j].CreatedAt)
		case "relevance":
			compScore := len(facts[i].CompetencyCategories) - len(facts[j].CompetencyCategories)
			if compScore != 0 {
				less = compScore < 0
			} else {
				audScore := len(facts[i].AudienceRelevance) - len(facts[j].AudienceRelevance)
				less = audScore < 0
			}
		default:
			less = facts[i].CreatedAt.Before(facts[j].CreatedAt)
		}

		if flm.sortOrder == "asc" {
			return less
		}
		return !less
	})
}

// GetSelectedFact returns the currently selected fact
func (flm *FactListModel) GetSelectedFact() *career.Fact {
	cursor := flm.listContainer.GetSelectedIdx()
	if cursor >= 0 && cursor < len(flm.filtered) {
		return flm.filtered[cursor]
	}
	return nil
}

// GetSelectedFacts returns all selected facts
func (flm *FactListModel) GetSelectedFacts() []*career.Fact {
	var selected []*career.Fact
	for _, fact := range flm.filtered {
		if flm.selectedFacts[fact.ID] {
			selected = append(selected, fact)
		}
	}
	return selected
}

// ClearSelection clears all selections
func (flm *FactListModel) ClearSelection() {
	flm.selectedFacts = make(map[string]bool)
}

// IsSubmitted returns true if a fact was selected
func (flm *FactListModel) IsSubmitted() bool {
	return flm.submitted
}

// IsCancelled returns true if the user cancelled
func (flm *FactListModel) IsCancelled() bool {
	return flm.cancelled
}

// GetError returns the last error
func (flm *FactListModel) GetError() error {
	return flm.err
}

// GetFacts returns the currently filtered facts
func (flm *FactListModel) GetFacts() []*career.Fact {
	return flm.filtered
}

// GetSelectedIdx returns the currently selected index
func (flm *FactListModel) GetSelectedIdx() int {
	return flm.listContainer.GetSelectedIdx()
}

// SetBreadcrumbs sets the breadcrumb trail for navigation
func (flm *FactListModel) SetBreadcrumbs(crumbs []string) {
	flm.breadcrumbs = crumbs
	flm.header.SetBreadcrumbs(crumbs)
}

// Refresh reloads the facts from the service
func (flm *FactListModel) Refresh() {
	flm.loadFactsSync()
	flm.pagination.GoToFirstPage()
	flm.listContainer.MoveToFirst()
}

// ListItemCallbacks implementation for FactListModel
// These methods implement the ListItemCallbacks interface required by ListNavigationKeyHandler

// OnView sends a message to view the selected fact
func (flm *FactListModel) OnView(item interface{}) tea.Cmd {
	if fact, ok := item.(*career.Fact); ok {
		return func() tea.Msg { return FactActionSelectedMsg{Fact: fact, Action: FactActionView} }
	}
	return nil
}

// OnEdit sends a message to edit the selected fact
func (flm *FactListModel) OnEdit(item interface{}) tea.Cmd {
	if fact, ok := item.(*career.Fact); ok {
		return func() tea.Msg { return FactActionSelectedMsg{Fact: fact, Action: FactActionEdit} }
	}
	return nil
}

// OnDelete initiates deletion of the selected fact
func (flm *FactListModel) OnDelete(item interface{}) {
	if fact, ok := item.(*career.Fact); ok {
		flm.deletionState.ShowConfirmation("Fact", fact.Text)
		flm.deletionState.DeletingItemID = fact.ID
	}
}

// OnToggleSelection toggles the selection state of an item
func (flm *FactListModel) OnToggleSelection(item interface{}) {
	if fact, ok := item.(*career.Fact); ok {
		flm.selectedFacts[fact.ID] = !flm.selectedFacts[fact.ID]
	}
}

// HasSelectedItem returns the currently selected fact
func (flm *FactListModel) HasSelectedItem() interface{} {
	return flm.GetSelectedFact()
}

// MoveUp moves the selection up by count items
func (flm *FactListModel) MoveUp(count int) {
	flm.listContainer.MoveUp(count)
}

// MoveDown moves the selection down by count items
func (flm *FactListModel) MoveDown(count int) {
	flm.listContainer.MoveDown(count)
}

// MoveToFirst moves selection to the first item
func (flm *FactListModel) MoveToFirst() {
	flm.listContainer.MoveToFirst()
}

// MoveToLast moves selection to the last item
func (flm *FactListModel) MoveToLast() {
	flm.listContainer.MoveToLast()
}

// UpdateDisplay refreshes the table display
func (flm *FactListModel) UpdateDisplay() {
	flm.updateTableRows()
}

// GetRowCount returns the number of visible rows
func (flm *FactListModel) GetRowCount() int {
	return len(flm.table.Rows())
}

// GetCurrentIndex returns the current selection index
func (flm *FactListModel) GetCurrentIndex() int {
	return flm.listContainer.GetSelectedIdx()
}

// GetPageSize returns the page size for pagination
func (flm *FactListModel) GetPageSize() int {
	return flm.pagination.GetPageSize()
}

// nextPage moves to the next page of results
func (flm *FactListModel) nextPage() {
	flm.pagination.NextPage()
	flm.listContainer.MoveToFirst()
	flm.updateTableRows()
}

// prevPage moves to the previous page of results
func (flm *FactListModel) prevPage() {
	flm.pagination.PrevPage()
	flm.listContainer.MoveToFirst()
	flm.updateTableRows()
}

// goToFirstItem moves to the first item
func (flm *FactListModel) goToFirstItem() {
	flm.pagination.GoToFirstPage()
	flm.listContainer.MoveToFirst()
}

// goToLastItem moves to the last item
func (flm *FactListModel) goToLastItem() {
	flm.pagination.GoToLastPage()
	pageFacts := flm.getPageFacts()
	if len(pageFacts) > 0 {
		flm.listContainer.SetSelectedIdx(len(pageFacts) - 1)
	}
}
