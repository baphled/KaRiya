package intents

import (
	"fmt"
	"sort"
	"strings"
	"time"

	"github.com/baphled/kariya/internal/cli/components"
	"github.com/baphled/kariya/internal/cli/navigation"
	"github.com/baphled/kariya/internal/cli/styles"
	domain "github.com/baphled/kariya/internal/domain/career"
	"github.com/charmbracelet/bubbles/table"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
)

// Custom message types for BurstManagement state transitions.

// BurstSelectedMsg indicates the user selected a burst.
type BurstSelectedMsg struct {
	Burst *domain.Burst
	Index int
}

// BurstManagementIntent implements the Intent interface for managing bursts.
// It owns the complete lifecycle of burst management, including:
// - Displaying a list of bursts
// - Selecting and viewing burst details
// - Returning the selected burst or cancelling
type BurstManagementIntent struct {
	// context is the input context passed to the intent.
	context *BurstManagementContext

	// state represents the current state of the intent.
	state *BurstManagementIntentModel

	// table is the table model for displaying bursts
	table *table.Model

	// listContainer provides table-based list UI
	listContainer *components.TableListContainer

	// navHandler centralizes navigation logic
	navHandler *navigation.ListNavigationHandler

	// active indicates whether this intent is currently active.
	active bool

	// result is the final result of the intent (set when complete).
	result *IntentResult[*BurstManagementResult]
}

// BurstManagementIntentModel represents the local state of the BurstManagement intent.
// This is the ONLY mutable state owned by the intent.
type BurstManagementIntentModel struct {
	// context is the input context passed to the intent.
	context *BurstManagementContext

	// currentState tracks which view is active.
	currentState string // BurstStateList, BurstStateDetail

	// filteredBursts are the bursts after applying current filters.
	filteredBursts []*domain.Burst

	// selectedIndex is the index of the currently selected burst.
	selectedIndex int

	// selectedBurst is the burst currently being viewed.
	selectedBurst *domain.Burst

	// viewedBursts tracks bursts viewed during the session.
	viewedBursts []*domain.Burst

	// Filter and sort state
	searchText       string
	filterCompetency string
	sortBy           string
	sortOrder        string
}

// State constants for BurstManagement intent.
const (
	BurstStateList   = "list"
	BurstStateDetail = "detail"
)

// NewBurstManagementIntent creates a new BurstManagement intent.
func NewBurstManagementIntent(context *BurstManagementContext) (*BurstManagementIntent, error) {
	// Validate the context.
	if err := context.Validate(); err != nil {
		return nil, err
	}

	// Create table model for bursts
	columns := []table.Column{
		{Title: "Name", Width: 30},
		{Title: "Competency", Width: 25},
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

	intent := &BurstManagementIntent{
		context: context,
		state: &BurstManagementIntentModel{
			context:          context,
			currentState:     BurstStateList,
			filteredBursts:   make([]*domain.Burst, 0),
			selectedIndex:    0,
			selectedBurst:    nil,
			viewedBursts:     make([]*domain.Burst, 0),
			searchText:       "",
			filterCompetency: "",
			sortBy:           "name",
			sortOrder:        "asc",
		},
		table:         &t,
		listContainer: components.NewTableListContainer(t, "Manage Bursts", 100),
		active:        true,
	}
	intent.navHandler = navigation.NewListNavigationHandler(intent)

	return intent, nil
}

// Init is called when the intent is activated.
func (i *BurstManagementIntent) Init() tea.Cmd {
	// Load bursts from repository
	_ = i.context.LoadBursts()

	// Initialize filtered bursts with the provided bursts.
	i.state.filteredBursts = i.context.Bursts
	if len(i.state.filteredBursts) > 0 {
		i.state.selectedBurst = i.state.filteredBursts[0]
	}
	i.updateTableRows()
	return nil
}

// updateTableRows updates the table rows based on filtered bursts, paginated
func (i *BurstManagementIntent) updateTableRows() {
	pageSize := 15
	total := len(i.state.filteredBursts)

	// Determine which page current selection is on
	page := 0
	if pageSize > 0 && i.state.selectedIndex >= 0 {
		page = i.state.selectedIndex / pageSize
	}

	start := page * pageSize
	end := start + pageSize
	if end > total {
		end = total
	}

	pageBursts := i.state.filteredBursts[start:end]

	rows := make([]table.Row, 0, len(pageBursts))
	for idx, burst := range pageBursts {
		realIdx := start + idx
		nameStr := burst.Name

		// Use navigation handler to format row text with indicator
		nameStr = i.navHandler.FormatRowText(realIdx, nameStr)

		competency := burst.CompetencyFocus
		if competency == "" {
			competency = "-"
		}

		eventCount := fmt.Sprintf("%d", len(burst.EventIDs))
		rows = append(rows, table.Row{nameStr, competency, eventCount})
	}

	i.table.SetRows(rows)

	// Calculate relative cursor for current page
	relativeCursor := i.state.selectedIndex - start

	// Set table cursor (for visual highlighting)
	i.table.SetCursor(relativeCursor)

	// Sync container's index to match (critical for rendering)
	i.listContainer.SetSelectedIdx(relativeCursor)

	// Update container with modified table
	i.listContainer.SetTable(*i.table)
}

// Update processes a message in the intent.
func (i *BurstManagementIntent) Update(msg tea.Msg) tea.Cmd {
	if !i.active {
		return nil
	}

	switch i.state.currentState {
	case BurstStateList:
		return i.updateListView(msg)

	case BurstStateDetail:
		return i.updateDetailView(msg)
	}

	return nil
}

// updateListView handles messages while viewing the burst list.
func (i *BurstManagementIntent) updateListView(msg tea.Msg) tea.Cmd {
	switch msg := msg.(type) {
	case tea.KeyMsg:
		// Try navigation handler first
		if i.navHandler.HandleKey(msg.String()) {
			return nil
		}

		switch msg.String() {
		case "enter":
			// Select current burst and move to detail view.
			if len(i.state.filteredBursts) > 0 && i.state.selectedBurst != nil {
				i.state.viewedBursts = append(i.state.viewedBursts, i.state.selectedBurst)
				i.state.currentState = BurstStateDetail
			}
			return nil

		case "q", "ctrl+c":
			i.setCancelled()
			return nil
		case "esc":
			i.setCancelled()
			return nil
		}

	case BurstSelectedMsg:
		i.state.selectedBurst = msg.Burst
		i.state.selectedIndex = msg.Index
		i.state.viewedBursts = append(i.state.viewedBursts, msg.Burst)
		i.state.currentState = BurstStateDetail
		return nil
	}

	return nil
}

// updateDetailView handles messages while viewing burst details.
func (i *BurstManagementIntent) updateDetailView(msg tea.Msg) tea.Cmd {
	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch msg.String() {
		case "enter":
			// Confirm selection and return burst.
			i.setCompleted()
			return nil

		case "esc":
			// Go back to list.
			i.state.currentState = BurstStateList
			return nil

		case "q", "ctrl+c":
			// Cancel.
			i.setCancelled()
			return nil
		}
	}

	return nil
}

// applyFilters filters the bursts based on current filter state.
func (i *BurstManagementIntent) applyFilters() {
	filtered := make([]*domain.Burst, 0)

	for _, burst := range i.context.Bursts {
		// Apply search text filter.
		if i.state.searchText != "" {
			if !strings.Contains(strings.ToLower(burst.Name), strings.ToLower(i.state.searchText)) &&
				!strings.Contains(strings.ToLower(burst.Description), strings.ToLower(i.state.searchText)) {
				continue
			}
		}

		// Apply competency filter.
		if i.state.filterCompetency != "" {
			if burst.CompetencyFocus != i.state.filterCompetency {
				continue
			}
		}

		filtered = append(filtered, burst)
	}

	// Apply sorting.
	sort.Slice(filtered, func(a, b int) bool {
		switch i.state.sortBy {
		case "name":
			if i.state.sortOrder == "asc" {
				return filtered[a].Name < filtered[b].Name
			}
			return filtered[a].Name > filtered[b].Name

		case "competency":
			if i.state.sortOrder == "asc" {
				return filtered[a].CompetencyFocus < filtered[b].CompetencyFocus
			}
			return filtered[a].CompetencyFocus > filtered[b].CompetencyFocus

		default: // date
			if i.state.sortOrder == "asc" {
				return filtered[a].CreatedAt.Before(filtered[b].CreatedAt)
			}
			return filtered[a].CreatedAt.After(filtered[b].CreatedAt)
		}
	})

	i.state.filteredBursts = filtered
}

// View renders the intent's current state.
func (i *BurstManagementIntent) View() string {
	switch i.state.currentState {
	case BurstStateList:
		return i.viewList()

	case BurstStateDetail:
		return i.viewDetail()
	}

	return ""
}

// viewList renders the burst list view with all bursts as a table.
// viewList renders the burst list view with all bursts as a table.
func (i *BurstManagementIntent) viewList() string {
	if len(i.state.filteredBursts) == 0 {
		i.listContainer.SetEmptyStateMessage("No bursts found.")
		return i.listContainer.Render()
	}

	// Ensure table rows are synchronized with current state
	i.updateTableRows()

	// Build pagination info with page number indicator
	pageSize := 15
	totalItems := len(i.state.filteredBursts)
	currentPage := (i.state.selectedIndex / pageSize) + 1
	totalPages := (totalItems + pageSize - 1) / pageSize
	paginationInfo := fmt.Sprintf("Bursts: %d | Page %d of %d", totalItems, currentPage, totalPages)
	i.listContainer.SetPaginationInfo(paginationInfo)

	// Set breadcrumbs if needed
	i.listContainer.SetBreadcrumbs([]string{"Home", "Bursts"})

	// Set help footer
	i.listContainer.SetHelpFooterKey("burst_management")

	return i.listContainer.Render()
}
func (i *BurstManagementIntent) viewDetail() string {
	if i.state.selectedBurst == nil {
		return "No burst selected."
	}

	var content strings.Builder
	content.WriteString("\nBurst Details\n\n")

	// Burst header.
	content.WriteString(fmt.Sprintf("Name: %s\n", i.state.selectedBurst.Name))

	if i.state.selectedBurst.Description != "" {
		content.WriteString(fmt.Sprintf("Description: %s\n", i.state.selectedBurst.Description))
	}
	if i.state.selectedBurst.CompetencyFocus != "" {
		content.WriteString(fmt.Sprintf("Competency Focus: %s\n", i.state.selectedBurst.CompetencyFocus))
	}

	content.WriteString(fmt.Sprintf("Events: %d\n", len(i.state.selectedBurst.EventIDs)))
	content.WriteString(fmt.Sprintf("Created: %s\n", i.state.selectedBurst.CreatedAt.Format("2006-01-02")))
	content.WriteString(fmt.Sprintf("Updated: %s\n", i.state.selectedBurst.UpdatedAt.Format("2006-01-02")))

	// Apply card styling.
	cardStyle := lipgloss.NewStyle().
		Padding(1, 2).
		BorderStyle(lipgloss.RoundedBorder()).
		BorderForeground(styles.ColorBorder).
		Background(styles.ColorBackgroundCard).
		Foreground(styles.ColorTextPrimary)

	card := cardStyle.Render(content.String())

	// Add footer with instructions.
	footerStyle := lipgloss.NewStyle().
		Foreground(styles.ColorTextSecondary).
		MarginTop(1)

	footer := footerStyle.Render("Enter to confirm, Esc to go back, q to cancel")

	return lipgloss.JoinVertical(lipgloss.Left, card, footer)
}

// Result returns the final result of the intent.
func (i *BurstManagementIntent) Result() *IntentResult[interface{}] {
	if i.result == nil {
		return nil
	}

	return &IntentResult[interface{}]{
		Status:   i.result.Status,
		Data:     i.result.Data,
		Error:    i.result.Error,
		Metadata: i.result.Metadata,
	}
}

// ListNavigator interface implementation
func (i *BurstManagementIntent) GetTotalItems() int {
	return len(i.state.filteredBursts)
}

func (i *BurstManagementIntent) GetSelectedIndex() int {
	return i.state.selectedIndex
}

func (i *BurstManagementIntent) SetSelectedIndex(idx int) {
	i.state.selectedIndex = idx
	if idx >= 0 && idx < len(i.state.filteredBursts) {
		i.state.selectedBurst = i.state.filteredBursts[idx]
	}
	i.updateTableRows()
}

func (i *BurstManagementIntent) GetPageSize() int {
	return 15
}

// Helper methods for result management.

func (i *BurstManagementIntent) setCompleted() {
	i.result = &IntentResult[*BurstManagementResult]{
		Status: Completed,
		Data: &BurstManagementResult{
			Action: "selected",
			Burst:  i.state.selectedBurst,
			Bursts: i.state.filteredBursts,
		},
		Metadata: map[string]interface{}{
			"selected_index": i.state.selectedIndex,
			"viewed_count":   len(i.state.viewedBursts),
			"timestamp":      time.Now(),
		},
	}
	i.active = false
}

func (i *BurstManagementIntent) setCancelled() {
	i.result = &IntentResult[*BurstManagementResult]{
		Status: Cancelled,
	}
	i.active = false
}

func (i *BurstManagementIntent) setFailed(code, message string, cause error) {
	i.result = &IntentResult[*BurstManagementResult]{
		Status: Failed,
		Error: &IntentError{
			Code:    code,
			Message: message,
			Cause:   cause,
		},
	}
	i.active = false
}
