package intents

import (
	"fmt"
	"sort"
	"strings"
	"time"

	"github.com/baphled/kariya/internal/cli/components"
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

	return intent, nil
}

// Init is called when the intent is activated.
func (i *BurstManagementIntent) Init() tea.Cmd {
	// Initialize filtered bursts with the provided bursts.
	i.state.filteredBursts = i.context.Bursts
	if len(i.state.filteredBursts) > 0 {
		i.state.selectedBurst = i.state.filteredBursts[0]
	}
	i.updateTableRows()
	return nil
}

// updateTableRows updates the table rows based on filtered bursts
func (i *BurstManagementIntent) updateTableRows() {
	rows := make([]table.Row, 0, len(i.state.filteredBursts))
	for idx, burst := range i.state.filteredBursts {
		nameStr := burst.Name

		// Add visual indicator for selected row
		if idx == i.state.selectedIndex {
			nameStr = "> " + nameStr
		} else {
			nameStr = "  " + nameStr
		}

		competency := burst.CompetencyFocus
		if competency == "" {
			competency = "-"
		}

		eventCount := fmt.Sprintf("%d", len(burst.EventIDs))
		rows = append(rows, table.Row{nameStr, competency, eventCount})
	}
	i.table.SetRows(rows)
	if i.state.selectedIndex < len(rows) {
		i.table.SetCursor(i.state.selectedIndex)
	} else if len(rows) > 0 {
		i.state.selectedIndex = 0
		i.table.SetCursor(0)
	}
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
		switch msg.String() {
		case "enter":
			// Select current burst and move to detail view.
			if len(i.state.filteredBursts) > 0 {
				i.state.selectedIndex = i.table.Cursor()
				if i.state.selectedIndex < len(i.state.filteredBursts) {
					i.state.selectedBurst = i.state.filteredBursts[i.state.selectedIndex]
					i.state.viewedBursts = append(i.state.viewedBursts, i.state.selectedBurst)
					i.state.currentState = BurstStateDetail
				}
			}
			return nil

		case "up", "k":
			// Move selection up.
			cursor := i.table.Cursor()
			if cursor > 0 {
				i.table.SetCursor(cursor - 1)
				i.state.selectedIndex = cursor - 1
				if len(i.state.filteredBursts) > 0 {
					i.state.selectedBurst = i.state.filteredBursts[i.state.selectedIndex]
				}
			}
			i.updateTableRows()
			return nil

		case "down", "j":
			// Move selection down.
			cursor := i.table.Cursor()
			if cursor < len(i.state.filteredBursts)-1 {
				i.table.SetCursor(cursor + 1)
				i.state.selectedIndex = cursor + 1
				if len(i.state.filteredBursts) > 0 {
					i.state.selectedBurst = i.state.filteredBursts[i.state.selectedIndex]
				}
			}
			i.updateTableRows()
			return nil

		case "q", "ctrl+c":
			// Cancel without selection.
			i.setCancelled()
			return nil

		case "esc":
			// Go back (no-op at list view).
			i.setCancelled()
			return nil
		}

	case BurstSelectedMsg:
		// Burst was selected (possibly by router or other component).
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
func (i *BurstManagementIntent) viewList() string {
	if len(i.state.filteredBursts) == 0 {
		i.listContainer.SetEmptyStateMessage("No bursts found.")
		return i.listContainer.Render()
	}

	// Build pagination info
	paginationInfo := fmt.Sprintf("Bursts: %d", len(i.state.filteredBursts))
	i.listContainer.SetPaginationInfo(paginationInfo)

	// Set breadcrumbs if needed
	i.listContainer.SetBreadcrumbs([]string{"Home", "Bursts"})

	// Set help footer
	i.listContainer.SetHelpFooterKey("burst_management")

	return i.listContainer.Render()
}

// viewDetail renders the burst detail view.
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
