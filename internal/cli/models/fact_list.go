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
	currentPage       int
	pageSize          int
	totalCount        int
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
	deletionConfirm   *ConfirmationDialog
	deletingFactID    string
	deleteSuccessMsg  string
	deleteErrorMsg    string
	showDeleteMessage bool
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
		width:             80,
		height:            20,
		pageSize:          10,
		currentPage:       1,
		totalCount:        0,
		sortBy:            "date",
		sortOrder:         "desc",
		helpFooter:        components.NewHelpFooter("fact_list", 80),
		deletionConfirm:   nil,
		deletingFactID:    "",
		deleteSuccessMsg:  "",
		deleteErrorMsg:    "",
		showDeleteMessage: false,
		breadcrumbs:       []string{"Home", "Facts"},
		header:            components.NewHeader("Facts", 80),
		listContainer:     components.NewTableListContainer(t, "Facts", 80),
	}
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
		flm.totalCount = len(flm.filtered)
		flm.updateTableRows()
		flm.err = nil
	}
}

// updateTableRows updates the table with rows from filtered facts
func (flm *FactListModel) updateTableRows() {
	var rows []table.Row

	for i, fact := range flm.filtered {
		factText := fact.Text
		if len(factText) > 57 {
			factText = factText[:54] + "..."
		}

		roleIcon := getRoleFitIcon(fact.RoleFit)

		competencies := strings.Join(fact.CompetencyCategories, ", ")
		if len(competencies) > 22 {
			competencies = competencies[:19] + "..."
		}

		// Add focus indicator for selected row based on container's selectedIdx
		indicator := "  "
		if i == flm.listContainer.GetSelectedIdx() {
			indicator = "▶ "
		}

		row := table.Row{
			indicator + factText,
			competencies,
			roleIcon,
		}
		rows = append(rows, row)
	}

	flm.table.SetRows(rows)
	// Ensure cursor is within valid bounds
	if len(rows) > 0 {
		if flm.table.Cursor() >= len(rows) {
			flm.listContainer.SetSelectedIdx(len(rows) - 1)
		}
		if flm.table.Cursor() < 0 {
			flm.listContainer.SetSelectedIdx(0)
		}
	}
	// Re-sync the table's cursor with the container's selectedIdx after SetRows
	// SetRows may reset the table's cursor, so we need to explicitly set it
	flm.table.SetCursor(flm.listContainer.GetSelectedIdx())
}

// Init initializes the model
func (flm *FactListModel) Init() tea.Cmd {
	return nil
}

// Update handles messages
func (flm *FactListModel) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	if flm.deletionConfirm != nil {
		updatedDialog, cmd := flm.deletionConfirm.Update(msg)
		flm.deletionConfirm = updatedDialog

		if flm.deletionConfirm.IsConfirmed() {
			return flm.performFactDeletion()
		}

		if flm.deletionConfirm.IsCancelled() {
			flm.deletionConfirm = nil
			flm.deletingFactID = ""
			return flm, nil
		}

		return flm, cmd
	}

	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch msg.String() {
		case "up", "k":
			flm.listContainer.MoveUp(1)
			flm.updateTableRows()
		case "down", "j":
			flm.listContainer.MoveDown(1)
			flm.updateTableRows()
		case "pgup", "ctrl+b":
			if flm.listContainer.GetSelectedIdx() >= flm.pageSize {
				flm.listContainer.MoveUp(flm.pageSize)
			} else {
				flm.listContainer.MoveToFirst()
			}
			flm.updateTableRows()
		case "pgdn", "ctrl+f":
			rows := len(flm.table.Rows())
			if flm.listContainer.GetSelectedIdx()+flm.pageSize < rows {
				flm.listContainer.MoveDown(flm.pageSize)
			} else {
				flm.listContainer.MoveToLast()
			}
			flm.updateTableRows()
		case "home", "g":
			flm.listContainer.MoveToFirst()
			flm.updateTableRows()
		case "end", "G":
			flm.listContainer.MoveToLast()
			flm.updateTableRows()
		case "enter":
			if len(flm.filtered) > 0 {
				fact := flm.GetSelectedFact()
				if fact != nil {
					return flm, func() tea.Msg { return FactActionMenuMsg{Fact: fact} }
				}
			}
		case " ", "space":
			if len(flm.filtered) > 0 {
				fact := flm.GetSelectedFact()
				if fact != nil {
					flm.selectedFacts[fact.ID] = !flm.selectedFacts[fact.ID]
				}
			}
		case "x", "d":
			if len(flm.filtered) > 0 {
				fact := flm.GetSelectedFact()
				if fact != nil {
					flm.showDeleteConfirmation(fact)
				}
			}
		case "v":
			if len(flm.filtered) > 0 {
				fact := flm.GetSelectedFact()
				if fact != nil {
					return flm, func() tea.Msg { return FactActionSelectedMsg{Fact: fact, Action: FactActionView} }
				}
			}
		case "e":
			if len(flm.filtered) > 0 {
				fact := flm.GetSelectedFact()
				if fact != nil {
					return flm, func() tea.Msg { return FactActionSelectedMsg{Fact: fact, Action: FactActionEdit} }
				}
			}
		case "esc":
			if flm.showDeleteMessage {
				flm.showDeleteMessage = false
				flm.deleteSuccessMsg = ""
				flm.deleteErrorMsg = ""
			} else {
				return flm, func() tea.Msg { return BackMsg{} }
			}
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

// showDeleteConfirmation shows the deletion confirmation dialog
func (flm *FactListModel) showDeleteConfirmation(fact *career.Fact) {
	message := fmt.Sprintf("Delete fact: \"%s\"?\n\nThis action cannot be undone.", truncateText(fact.Text, 60))
	flm.deletionConfirm = NewConfirmationDialog("Delete Fact", message)
	flm.deletingFactID = fact.ID
}

// performFactDeletion performs the actual fact deletion
func (flm *FactListModel) performFactDeletion() (tea.Model, tea.Cmd) {
	if flm.deletingFactID == "" {
		flm.deletionConfirm = nil
		return flm, nil
	}

	err := flm.service.DeleteFact(flm.ctx, flm.deletingFactID)
	if err != nil {
		flm.deleteErrorMsg = fmt.Sprintf("Error deleting fact: %v", err)
		flm.showDeleteMessage = true
		flm.deletionConfirm = nil
		flm.deletingFactID = ""
		return flm, nil
	}

	flm.removeFactFromLists(flm.deletingFactID)

	flm.deleteSuccessMsg = "Fact deleted successfully"
	flm.showDeleteMessage = true
	flm.deletionConfirm = nil
	flm.deletingFactID = ""

	flm.totalCount = len(flm.filtered)
	flm.updateTableRows()

	if flm.totalCount == 0 {
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
	if flm.deletionConfirm != nil {
		return flm.deletionConfirm.View()
	}

	if flm.showDeleteMessage {
		return flm.renderDeleteMessage()
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
		SetHelpFooterKey("fact_list")

	// Set pagination info - always show pagination like list.go does
	startIdx := (flm.currentPage-1)*flm.pageSize + 1
	endIdx := startIdx + len(flm.getPageFacts()) - 1
	if flm.totalCount == 0 {
		startIdx = 0
		endIdx = 0
	}
	paginationText := fmt.Sprintf("Showing %d-%d of %d facts", startIdx, endIdx, flm.totalCount)
	flm.listContainer.SetPaginationInfo(paginationText)

	return flm.listContainer.Render()
}

// renderDeleteMessage renders the deletion success/error message
func (flm *FactListModel) renderDeleteMessage() string {
	var messageContent string
	if flm.deleteSuccessMsg != "" {
		messageContent = styles.SuccessBox.Render(flm.deleteSuccessMsg + "\n\nPress 'esc' to continue")
	} else if flm.deleteErrorMsg != "" {
		messageContent = styles.ErrorBox.Render(flm.deleteErrorMsg + "\n\nPress 'esc' to continue")
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
	if flm.totalCount == 0 {
		return []*career.Fact{}
	}

	startIdx := (flm.currentPage - 1) * flm.pageSize
	endIdx := startIdx + flm.pageSize

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
	flm.totalCount = len(flm.filtered)
	flm.currentPage = 1
	flm.listContainer.MoveToFirst()
	flm.updateTableRows()
}

// SetCompetencyFilter sets the competency filter
func (flm *FactListModel) SetCompetencyFilter(competency string) {
	flm.competencyFilter = competency
	flm.applyFiltersAndSort()
	flm.totalCount = len(flm.filtered)
	flm.currentPage = 1
	flm.listContainer.MoveToFirst()
	flm.updateTableRows()
}

// SetRoleFitFilter sets the role fit filter
func (flm *FactListModel) SetRoleFitFilter(roleFit career.RoleFit) {
	flm.roleFitFilter = roleFit
	flm.applyFiltersAndSort()
	flm.totalCount = len(flm.filtered)
	flm.currentPage = 1
	flm.listContainer.MoveToFirst()
	flm.updateTableRows()
}

// SetAudienceFilter sets the audience filter
func (flm *FactListModel) SetAudienceFilter(audience string) {
	flm.audienceFilter = audience
	flm.applyFiltersAndSort()
	flm.totalCount = len(flm.filtered)
	flm.currentPage = 1
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

// getTotalPages calculates the total number of pages
func (flm *FactListModel) getTotalPages() int {
	if flm.totalCount == 0 {
		return 1
	}
	pages := (flm.totalCount + flm.pageSize - 1) / flm.pageSize
	return pages
}

// SetBreadcrumbs sets the breadcrumb trail for navigation
func (flm *FactListModel) SetBreadcrumbs(crumbs []string) {
	flm.breadcrumbs = crumbs
	flm.header.SetBreadcrumbs(crumbs)
}

// Refresh reloads the facts from the service
func (flm *FactListModel) Refresh() {
	flm.loadFactsSync()
	flm.currentPage = 1
	flm.listContainer.MoveToFirst()
}

// nextPage moves to the next page of results
func (flm *FactListModel) nextPage() {
	totalPages := flm.getTotalPages()
	if flm.currentPage < totalPages {
		flm.currentPage++
		flm.listContainer.MoveToFirst()
		flm.updateTableRows()
	}
}

// prevPage moves to the previous page of results
func (flm *FactListModel) prevPage() {
	if flm.currentPage > 1 {
		flm.currentPage--
		flm.listContainer.MoveToFirst()
		flm.updateTableRows()
	}
}

// goToFirstItem moves to the first item
func (flm *FactListModel) goToFirstItem() {
	flm.currentPage = 1
	flm.listContainer.MoveToFirst()
}

// goToLastItem moves to the last item
func (flm *FactListModel) goToLastItem() {
	totalPages := flm.getTotalPages()
	flm.currentPage = totalPages
	pageFacts := flm.getPageFacts()
	if len(pageFacts) > 0 {
		flm.listContainer.SetSelectedIdx(len(pageFacts) - 1)
	}
}

// truncateText truncates text to a maximum length
func truncateText(text string, maxLen int) string {
	if len(text) > maxLen {
		return text[:maxLen-3] + "..."
	}
	return text
}
