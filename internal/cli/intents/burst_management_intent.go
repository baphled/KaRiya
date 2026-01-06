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

// BurstEventsLoadedMsg is sent when events for a burst are loaded
type BurstEventsLoadedMsg struct {
	Events []*domain.CareerEvent
	Error  error
}

// BurstFactsLoadedMsg is sent when facts for a burst are loaded
type BurstFactsLoadedMsg struct {
	Facts []*domain.Fact
	Error error
}

// BurstEditCompleteMsg is sent when burst editing is complete
type BurstEditCompleteMsg struct {
	Burst     *domain.Burst
	Cancelled bool
	Error     error
}

// BurstDeletedMsg is sent when a burst is deleted
type BurstDeletedMsg struct {
	BurstID string
	Error   error
}

// BurstConfirmedMsg is sent when a burst is confirmed
type BurstConfirmedMsg struct {
	Burst *domain.Burst
	Error error
}

// FactExtractionCompleteMsg is sent when fact extraction is complete
type FactExtractionCompleteMsg struct {
	Facts []*domain.Fact
	Error error
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

	// Events and facts for current burst
	burstEvents []*domain.CareerEvent
	burstFacts  []*domain.Fact

	// Loading states
	loadingEvents bool
	loadingFacts  bool

	// Edit and delete state
	burstEditor tea.Model
	deleteError error
	editError   error

	// Confirmation state
	confirmError        error
	extractingFacts     bool
	extractionComplete  bool
	extractedFactsCount int
	existingFactsCount  int
	showReextractPrompt bool

	// Filter and sort state
	searchText string
	sortBy     string
	sortOrder  string
}

// State constants for BurstManagement intent.
const (
	BurstStateList            = "list"
	BurstStateDetail          = "detail"
	BurstStateDetailEvents    = "detail_events"
	BurstStateDetailFacts     = "detail_facts"
	BurstStateEdit            = "edit"
	BurstStateDeleteConfirm   = "delete_confirm"
	BurstStateConfirm         = "confirm"
	BurstStateExtractingFacts = "extracting_facts"
)

// NewBurstManagementIntent creates a new BurstManagement intent.
func NewBurstManagementIntent(context *BurstManagementContext) (*BurstManagementIntent, error) {
	// Validate the context.
	if err := context.Validate(); err != nil {
		return nil, err
	}

	// Create table model for bursts with enhanced columns
	// Total column width: 30 + 35 + 10 + 8 + 12 = 95 chars
	columns := []table.Column{
		{Title: "Name", Width: 30},
		{Title: "Description", Width: 35},
		{Title: "Confirmed", Width: 10},
		{Title: "Events", Width: 8},
		{Title: "Created", Width: 12},
	}

	t := table.New(
		table.WithColumns(columns),
		table.WithRows([]table.Row{}),
		table.WithFocused(true),
		table.WithHeight(15),
		table.WithWidth(105), // Adjusted for 5 columns
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
			context:        context,
			currentState:   BurstStateList,
			filteredBursts: make([]*domain.Burst, 0),
			selectedIndex:  0,
			selectedBurst:  nil,
			viewedBursts:   make([]*domain.Burst, 0),
			searchText:     "",
			sortBy:         "name",
			sortOrder:      "asc",
		},
		table:         &t,
		listContainer: components.NewTableListContainer(t, "Manage Bursts", 105),
		active:        true,
	}
	intent.navHandler = navigation.NewListNavigationHandler(intent)

	return intent, nil
}

// formatConfirmedStatus returns a plain text confirmed status string.
// Note: BubbleTea table doesn't support Lipgloss-styled cells, so we use plain text.
func (i *BurstManagementIntent) formatConfirmedStatus(confirmed bool) string {
	if confirmed {
		return "✓ Yes"
	}
	return "✗ No"
}

// formatDescription returns a truncated description preview (max 35 chars).
// Note: BubbleTea table doesn't support Lipgloss-styled cells, so we use plain text.
func (i *BurstManagementIntent) formatDescription(description string) string {
	desc := strings.TrimSpace(description)
	// Remove newlines and carriage returns
	desc = strings.ReplaceAll(desc, "\n", " ")
	desc = strings.ReplaceAll(desc, "\r", " ")

	if desc == "" {
		return "-"
	}

	maxLen := 32 // 35 - 3 for "..."
	if len(desc) > maxLen {
		return desc[:maxLen] + "..."
	}

	return desc
}

// formatCreatedDate returns a formatted date string (YYYY-MM-DD).
// Note: BubbleTea table doesn't support Lipgloss-styled cells, so we use plain text.
func (i *BurstManagementIntent) formatCreatedDate(createdAt time.Time) string {
	return createdAt.Format("2006-01-02")
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

		// Column 1: Name (truncate to 27 chars for focus indicator, total width 30)
		nameStr := burst.Name
		if len(nameStr) > 27 {
			nameStr = nameStr[:27] + "..."
		}
		// Add focus indicator via navigation handler
		nameStr = i.navHandler.FormatRowText(realIdx, nameStr)

		// Column 2: Description (truncated preview, max 35 chars)
		descStr := i.formatDescription(burst.Description)

		// Column 3: Confirmed Status (icon + colored text)
		confirmedStr := i.formatConfirmedStatus(burst.Confirmed)

		// Column 4: Event Count
		eventCount := fmt.Sprintf("%d", len(burst.EventIDs))

		// Column 5: Created Date (YYYY-MM-DD)
		createdStr := i.formatCreatedDate(burst.CreatedAt)

		rows = append(rows, table.Row{
			nameStr,
			descStr,
			confirmedStr,
			eventCount,
			createdStr,
		})
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

	case BurstStateDetailEvents:
		return i.updateDetailEventsView(msg)

	case BurstStateDetailFacts:
		return i.updateDetailFactsView(msg)

	case BurstStateEdit:
		return i.updateEditView(msg)

	case BurstStateDeleteConfirm:
		return i.updateDeleteConfirmView(msg)

	case BurstStateConfirm:
		return i.updateConfirmView(msg)

	case BurstStateExtractingFacts:
		return i.updateExtractingFactsView(msg)
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
		case "e":
			// View events
			i.state.currentState = BurstStateDetailEvents
			i.state.loadingEvents = true
			return i.loadEventsForBurst()

		case "f":
			// View facts
			i.state.currentState = BurstStateDetailFacts
			i.state.loadingFacts = true
			return i.loadFactsForBurst()

		case "x":
			// Edit burst
			i.state.currentState = BurstStateEdit
			return i.initBurstEditor()

		case "d":
			// Delete burst
			i.state.currentState = BurstStateDeleteConfirm
			return nil

		case "c":
			// Confirm burst (extract facts if needed)
			i.state.currentState = BurstStateConfirm
			return i.checkForExistingFacts()

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

// updateDetailEventsView handles messages while viewing events in a burst
func (i *BurstManagementIntent) updateDetailEventsView(msg tea.Msg) tea.Cmd {
	switch msg := msg.(type) {
	case BurstEventsLoadedMsg:
		i.state.loadingEvents = false
		if msg.Error == nil {
			i.state.burstEvents = msg.Events
		}
		return nil

	case tea.KeyMsg:
		switch msg.String() {
		case "esc":
			// Back to detail
			i.state.currentState = BurstStateDetail
			return nil

		case "q", "ctrl+c":
			i.setCancelled()
			return nil
		}
	}

	return nil
}

// updateDetailFactsView handles messages while viewing facts from a burst
func (i *BurstManagementIntent) updateDetailFactsView(msg tea.Msg) tea.Cmd {
	switch msg := msg.(type) {
	case BurstFactsLoadedMsg:
		i.state.loadingFacts = false
		if msg.Error == nil {
			i.state.burstFacts = msg.Facts
		}
		return nil

	case tea.KeyMsg:
		switch msg.String() {
		case "esc":
			// Back to detail
			i.state.currentState = BurstStateDetail
			return nil

		case "q", "ctrl+c":
			i.setCancelled()
			return nil
		}
	}

	return nil
}

// initBurstEditor initializes the burst editor with the current burst
func (i *BurstManagementIntent) initBurstEditor() tea.Cmd {
	if i.state.selectedBurst == nil {
		return nil
	}

	// Import the models package to access BurstEditorModel
	// Note: We'll create a simple inline editor instead of importing models
	// to avoid circular dependencies
	return nil
}

// updateEditView handles the edit state
func (i *BurstManagementIntent) updateEditView(msg tea.Msg) tea.Cmd {
	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch msg.String() {
		case "esc":
			// Cancel edit and go back to detail
			i.state.currentState = BurstStateDetail
			i.state.editError = nil
			return nil

		case "ctrl+s":
			// Save changes
			if i.state.selectedBurst != nil {
				err := i.context.UpdateBurst(i.state.selectedBurst)
				if err != nil {
					i.state.editError = err
					return nil
				}

				// Reload bursts
				i.context.LoadBursts()
				i.state.filteredBursts = i.context.Bursts

				// Go back to detail view
				i.state.currentState = BurstStateDetail
				return nil
			}
		}
	}

	return nil
}

// updateDeleteConfirmView handles the delete confirmation state
func (i *BurstManagementIntent) updateDeleteConfirmView(msg tea.Msg) tea.Cmd {
	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch msg.String() {
		case "y":
			// Confirm delete
			if i.state.selectedBurst != nil {
				err := i.context.DeleteBurst(i.state.selectedBurst.ID)
				if err != nil {
					i.state.deleteError = err
					return nil
				}

				// Reload bursts
				i.context.LoadBursts()
				i.state.filteredBursts = i.context.Bursts
				i.state.selectedBurst = nil
				i.state.selectedIndex = 0

				// Go back to list
				i.state.currentState = BurstStateList
				return nil
			}

		case "n", "esc":
			// Cancel delete
			i.state.currentState = BurstStateDetail
			i.state.deleteError = nil
			return nil

		case "q", "ctrl+c":
			i.setCancelled()
			return nil
		}
	}

	return nil
}

// updateConfirmView handles the confirmation state
func (i *BurstManagementIntent) updateConfirmView(msg tea.Msg) tea.Cmd {
	switch msg := msg.(type) {
	case BurstFactsLoadedMsg:
		// Facts loaded, check count
		if msg.Error != nil {
			i.state.confirmError = msg.Error
			return nil
		}

		i.state.existingFactsCount = len(msg.Facts)

		if i.state.existingFactsCount > 0 {
			// Facts already exist, show re-extract prompt
			i.state.showReextractPrompt = true
			return nil
		}

		// No facts exist, start extraction
		i.state.currentState = BurstStateExtractingFacts
		i.state.extractingFacts = true
		return i.extractFacts()

	case BurstConfirmedMsg:
		// Handle confirmation completion (burst confirmed without extraction)
		if msg.Error != nil {
			i.state.confirmError = msg.Error
			return nil
		}

		// Success! extractionComplete should already be set by confirmBurstOnly()
		// State should already be BurstStateConfirm
		// Message properly handled, success view will be shown
		return nil

	case tea.KeyMsg:
		if i.state.showReextractPrompt {
			switch msg.String() {
			case "y":
				// Re-extract facts
				i.state.showReextractPrompt = false
				i.state.currentState = BurstStateExtractingFacts
				i.state.extractingFacts = true
				return i.extractFacts()

			case "n", "esc":
				// Don't re-extract, just mark as confirmed
				return i.confirmBurstOnly()

			case "q", "ctrl+c":
				i.setCancelled()
				return nil
			}
		} else if i.state.extractionComplete {
			// Extraction complete, any key returns to detail
			switch msg.String() {
			case "q", "ctrl+c":
				i.setCancelled()
				return nil
			default:
				// Any other key returns to detail
				i.state.currentState = BurstStateDetail
				i.state.confirmError = nil
				i.state.extractionComplete = false
				i.state.extractedFactsCount = 0
				return nil
			}
		} else {
			switch msg.String() {
			case "esc":
				// Go back to detail
				i.state.currentState = BurstStateDetail
				i.state.confirmError = nil
				return nil

			case "q", "ctrl+c":
				i.setCancelled()
				return nil
			}
		}
	}

	return nil
}

// updateExtractingFactsView handles the fact extraction state
func (i *BurstManagementIntent) updateExtractingFactsView(msg tea.Msg) tea.Cmd {
	switch msg := msg.(type) {
	case FactExtractionCompleteMsg:
		i.state.extractingFacts = false

		if msg.Error != nil {
			i.state.confirmError = msg.Error
			i.state.currentState = BurstStateConfirm
			return nil
		}

		i.state.extractedFactsCount = len(msg.Facts)
		i.state.extractionComplete = true

		// Now confirm the burst
		return i.confirmBurstOnly()

	case tea.KeyMsg:
		switch msg.String() {
		case "q", "ctrl+c":
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

	case BurstStateDetailEvents:
		return i.viewDetailEvents()

	case BurstStateDetailFacts:
		return i.viewDetailFacts()

	case BurstStateEdit:
		return i.viewEdit()

	case BurstStateDeleteConfirm:
		return i.viewDeleteConfirm()

	case BurstStateConfirm:
		return i.viewConfirm()

	case BurstStateExtractingFacts:
		return i.viewExtractingFacts()
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

	// Title with confirmation status indicator
	titleStyle := lipgloss.NewStyle().
		Bold(true).
		Foreground(styles.ColorTextPrimary)

	title := "Burst Details"
	if i.state.selectedBurst.Confirmed {
		confirmedStyle := lipgloss.NewStyle().
			Foreground(styles.ColorSuccess).
			Bold(true)
		title = "Burst Details " + confirmedStyle.Render("✓ Confirmed")
	}
	content.WriteString("\n" + titleStyle.Render(title) + "\n\n")

	// Burst header.
	content.WriteString(fmt.Sprintf("Name: %s\n", i.state.selectedBurst.Name))

	if i.state.selectedBurst.Description != "" {
		content.WriteString(fmt.Sprintf("Description: %s\n", i.state.selectedBurst.Description))
	}

	content.WriteString(fmt.Sprintf("Events: %d\n", len(i.state.selectedBurst.EventIDs)))

	// Show confirmation details if confirmed
	if i.state.selectedBurst.Confirmed && i.state.selectedBurst.ConfirmedAt != nil {
		content.WriteString(fmt.Sprintf("Confirmed: %s\n", i.state.selectedBurst.ConfirmedAt.Format("2006-01-02 15:04")))
	}

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

	footer := footerStyle.Render("e=events, f=facts, c=confirm burst, x=edit, d=delete, Enter=select, Esc=back, q=cancel")

	return lipgloss.JoinVertical(lipgloss.Left, card, footer)
}

// viewDetailEvents renders the events view for a burst
func (i *BurstManagementIntent) viewDetailEvents() string {
	if i.state.loadingEvents {
		infoStyle := lipgloss.NewStyle().Foreground(styles.ColorInfo)
		return infoStyle.Render("Loading events...")
	}

	if len(i.state.burstEvents) == 0 {
		errorStyle := lipgloss.NewStyle().Foreground(styles.ColorError)
		return errorStyle.Render("No events found for this burst.")
	}

	var content strings.Builder
	headerStyle := lipgloss.NewStyle().
		Bold(true).
		Foreground(styles.ColorTextPrimary).
		MarginBottom(1)

	content.WriteString(headerStyle.Render(
		fmt.Sprintf("Events in Burst: %s", i.state.selectedBurst.Name),
	))
	content.WriteString("\n\n")

	for idx, event := range i.state.burstEvents {
		content.WriteString(fmt.Sprintf("%d. %s\n", idx+1, event.Text))
		content.WriteString(fmt.Sprintf("   Date: %s\n", event.Date.Format("2006-01-02")))
		if len(event.Tags) > 0 {
			content.WriteString(fmt.Sprintf("   Tags: %s\n", strings.Join(event.Tags, ", ")))
		}
		content.WriteString("\n")
	}

	footerStyle := lipgloss.NewStyle().
		Foreground(styles.ColorTextSecondary).
		MarginTop(1)

	footer := footerStyle.Render("Esc=back, q=quit")
	content.WriteString(footer)

	return content.String()
}

// viewDetailFacts renders the facts view for a burst
func (i *BurstManagementIntent) viewDetailFacts() string {
	if i.state.loadingFacts {
		infoStyle := lipgloss.NewStyle().Foreground(styles.ColorInfo)
		return infoStyle.Render("Loading facts...")
	}

	if len(i.state.burstFacts) == 0 {
		warningStyle := lipgloss.NewStyle().
			Foreground(styles.ColorWarning).
			MarginBottom(1)

		return warningStyle.Render("No facts extracted yet. Confirm the burst in detail view to extract facts.")
	}

	var content strings.Builder
	headerStyle := lipgloss.NewStyle().
		Bold(true).
		Foreground(styles.ColorTextPrimary).
		MarginBottom(1)

	content.WriteString(headerStyle.Render(
		fmt.Sprintf("Facts from Burst: %s", i.state.selectedBurst.Name),
	))
	content.WriteString("\n\n")

	for idx, fact := range i.state.burstFacts {
		// Display fact text
		content.WriteString(fmt.Sprintf("%d. %s\n", idx+1, fact.Text))

		// Display competency categories
		if len(fact.CompetencyCategories) > 0 {
			content.WriteString(fmt.Sprintf("   Categories: %s\n", strings.Join(fact.CompetencyCategories, ", ")))
		}

		// Display strength signal
		if fact.StrengthSignal != "" {
			content.WriteString(fmt.Sprintf("   Strength: %s\n", fact.StrengthSignal))
		}

		content.WriteString("\n")
	}

	footerStyle := lipgloss.NewStyle().
		Foreground(styles.ColorTextSecondary).
		MarginTop(1)

	footer := footerStyle.Render("Esc=back, q=quit")
	content.WriteString(footer)

	return content.String()
}

// viewEdit renders the edit view for a burst
func (i *BurstManagementIntent) viewEdit() string {
	if i.state.selectedBurst == nil {
		return "No burst selected."
	}

	var content strings.Builder

	// Header
	headerStyle := lipgloss.NewStyle().
		Bold(true).
		Foreground(styles.ColorTextPrimary).
		MarginBottom(1)

	content.WriteString(headerStyle.Render("Edit Burst"))
	content.WriteString("\n\n")

	// Show error if any
	if i.state.editError != nil {
		errorStyle := lipgloss.NewStyle().
			Foreground(styles.ColorError).
			MarginBottom(1)
		content.WriteString(errorStyle.Render(fmt.Sprintf("Error: %s", i.state.editError)))
		content.WriteString("\n\n")
	}

	// Display current burst details
	content.WriteString(fmt.Sprintf("Name: %s\n", i.state.selectedBurst.Name))

	if i.state.selectedBurst.Description != "" {
		content.WriteString(fmt.Sprintf("Description: %s\n", i.state.selectedBurst.Description))
	}

	// Note: For now, this is a simple view showing current values
	// A full implementation would use text inputs for editing
	noteStyle := lipgloss.NewStyle().
		Foreground(styles.ColorTextSecondary).
		Italic(true).
		MarginTop(1)

	content.WriteString("\n")
	content.WriteString(noteStyle.Render("Note: Full edit functionality coming soon."))
	content.WriteString("\n")

	// Apply card styling
	cardStyle := lipgloss.NewStyle().
		Padding(1, 2).
		BorderStyle(lipgloss.RoundedBorder()).
		BorderForeground(styles.ColorBorder).
		Background(styles.ColorBackgroundCard).
		Foreground(styles.ColorTextPrimary)

	card := cardStyle.Render(content.String())

	// Footer with instructions
	footerStyle := lipgloss.NewStyle().
		Foreground(styles.ColorTextSecondary).
		MarginTop(1)

	footer := footerStyle.Render("Ctrl+S=save (placeholder), Esc=cancel")

	return lipgloss.JoinVertical(lipgloss.Left, card, footer)
}

// viewDeleteConfirm renders the delete confirmation view
func (i *BurstManagementIntent) viewDeleteConfirm() string {
	if i.state.selectedBurst == nil {
		return "No burst selected."
	}

	var content strings.Builder

	// Warning header
	warningStyle := lipgloss.NewStyle().
		Bold(true).
		Foreground(styles.ColorError).
		MarginBottom(1)

	content.WriteString(warningStyle.Render("⚠️  DELETE BURST - Are you sure?"))
	content.WriteString("\n\n")

	// Show error if any
	if i.state.deleteError != nil {
		errorStyle := lipgloss.NewStyle().
			Foreground(styles.ColorError).
			MarginBottom(1)
		content.WriteString(errorStyle.Render(fmt.Sprintf("Error: %s", i.state.deleteError)))
		content.WriteString("\n\n")
	}

	// Burst details
	infoStyle := lipgloss.NewStyle().
		Foreground(styles.ColorTextPrimary)

	content.WriteString(infoStyle.Render(fmt.Sprintf("Burst Name: %s", i.state.selectedBurst.Name)))
	content.WriteString("\n\n")

	// Warning messages
	warningTextStyle := lipgloss.NewStyle().
		Foreground(styles.ColorWarning).
		Bold(true)

	content.WriteString(warningTextStyle.Render("This will remove the burst grouping but NOT delete the events."))
	content.WriteString("\n")
	content.WriteString(warningTextStyle.Render("This action cannot be undone."))
	content.WriteString("\n")

	// Apply card styling
	cardStyle := lipgloss.NewStyle().
		Padding(1, 2).
		BorderStyle(lipgloss.RoundedBorder()).
		BorderForeground(styles.ColorError).
		Background(styles.ColorBackgroundCard).
		Foreground(styles.ColorTextPrimary)

	card := cardStyle.Render(content.String())

	// Footer with instructions
	footerStyle := lipgloss.NewStyle().
		Foreground(styles.ColorTextSecondary).
		MarginTop(1)

	footer := footerStyle.Render("y=confirm delete, n/Esc=cancel")

	return lipgloss.JoinVertical(lipgloss.Left, card, footer)
}

// viewConfirm renders the confirmation view
func (i *BurstManagementIntent) viewConfirm() string {
	if i.state.selectedBurst == nil {
		return "No burst selected."
	}

	var content strings.Builder

	// Header
	headerStyle := lipgloss.NewStyle().
		Bold(true).
		Foreground(styles.ColorSuccess).
		MarginBottom(1)

	content.WriteString(headerStyle.Render("Confirm Burst"))
	content.WriteString("\n\n")

	// Show error if any
	if i.state.confirmError != nil {
		errorStyle := lipgloss.NewStyle().
			Foreground(styles.ColorError).
			MarginBottom(1)
		content.WriteString(errorStyle.Render(fmt.Sprintf("Error: %s", i.state.confirmError)))
		content.WriteString("\n\n")
	}

	// Burst details
	infoStyle := lipgloss.NewStyle().
		Foreground(styles.ColorTextPrimary)

	content.WriteString(infoStyle.Render(fmt.Sprintf("Burst: %s", i.state.selectedBurst.Name)))
	content.WriteString("\n\n")

	// Show different messages based on state
	if i.state.showReextractPrompt {
		// Facts already exist
		warningStyle := lipgloss.NewStyle().
			Foreground(styles.ColorWarning).
			Bold(true)

		content.WriteString(warningStyle.Render(fmt.Sprintf("This burst already has %d facts extracted.", i.state.existingFactsCount)))
		content.WriteString("\n\n")
		content.WriteString(infoStyle.Render("Do you want to extract more facts? New facts will be added to existing ones."))
		content.WriteString("\n")
	} else if i.state.extractionComplete {
		// Extraction completed successfully
		successStyle := lipgloss.NewStyle().
			Foreground(styles.ColorSuccess).
			Bold(true)

		content.WriteString(successStyle.Render(fmt.Sprintf("✓ Successfully extracted and saved %d facts!", i.state.extractedFactsCount)))
		content.WriteString("\n\n")
		content.WriteString(infoStyle.Render("Burst has been confirmed."))
		content.WriteString("\n")
	} else {
		// About to start extraction
		content.WriteString(infoStyle.Render("No facts found for this burst. Starting fact extraction..."))
		content.WriteString("\n")
	}

	// Apply card styling
	cardStyle := lipgloss.NewStyle().
		Padding(1, 2).
		BorderStyle(lipgloss.RoundedBorder()).
		BorderForeground(styles.ColorBorder).
		Background(styles.ColorBackgroundCard).
		Foreground(styles.ColorTextPrimary)

	card := cardStyle.Render(content.String())

	// Footer with instructions
	footerStyle := lipgloss.NewStyle().
		Foreground(styles.ColorTextSecondary).
		MarginTop(1)

	var footer string
	if i.state.showReextractPrompt {
		footer = footerStyle.Render("y=re-extract facts, n=skip re-extraction, Esc=cancel")
	} else if i.state.extractionComplete {
		footer = footerStyle.Render("Press any key to continue...")
	} else {
		footer = footerStyle.Render("Esc=cancel")
	}

	return lipgloss.JoinVertical(lipgloss.Left, card, footer)
}

// viewExtractingFacts renders the fact extraction progress view
func (i *BurstManagementIntent) viewExtractingFacts() string {
	if i.state.selectedBurst == nil {
		return "No burst selected."
	}

	var content strings.Builder

	// Header
	headerStyle := lipgloss.NewStyle().
		Bold(true).
		Foreground(styles.ColorInfo).
		MarginBottom(1)

	content.WriteString(headerStyle.Render("Extracting Facts"))
	content.WriteString("\n\n")

	// Progress indicator
	infoStyle := lipgloss.NewStyle().
		Foreground(styles.ColorTextPrimary)

	content.WriteString(infoStyle.Render(fmt.Sprintf("Burst: %s", i.state.selectedBurst.Name)))
	content.WriteString("\n\n")

	progressStyle := lipgloss.NewStyle().
		Foreground(styles.ColorInfo).
		Bold(true)

	content.WriteString(progressStyle.Render("⏳ Extracting and saving facts..."))
	content.WriteString("\n\n")
	content.WriteString(infoStyle.Render("Analyzing events and persisting facts to database."))
	content.WriteString("\n\n")
	content.WriteString(infoStyle.Render("This may take a few moments."))
	content.WriteString("\n")

	// Apply card styling
	cardStyle := lipgloss.NewStyle().
		Padding(1, 2).
		BorderStyle(lipgloss.RoundedBorder()).
		BorderForeground(styles.ColorBorder).
		Background(styles.ColorBackgroundCard).
		Foreground(styles.ColorTextPrimary)

	card := cardStyle.Render(content.String())

	// Footer with instructions
	footerStyle := lipgloss.NewStyle().
		Foreground(styles.ColorTextSecondary).
		MarginTop(1)

	footer := footerStyle.Render("Please wait... (q to cancel)")

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

// loadEventsForBurst loads events for the selected burst
func (i *BurstManagementIntent) loadEventsForBurst() tea.Cmd {
	return func() tea.Msg {
		if i.state.selectedBurst == nil {
			return BurstEventsLoadedMsg{Error: fmt.Errorf("no burst selected")}
		}

		// Use service to get events by IDs
		events := make([]*domain.CareerEvent, 0, len(i.state.selectedBurst.EventIDs))
		for _, eventID := range i.state.selectedBurst.EventIDs {
			event, err := i.context.Service.GetEventByID(i.context.Context, eventID)
			if err != nil {
				// Skip missing events, but continue loading others
				continue
			}
			events = append(events, event)
		}

		return BurstEventsLoadedMsg{Events: events}
	}
}

// loadFactsForBurst loads facts for the selected burst
func (i *BurstManagementIntent) loadFactsForBurst() tea.Cmd {
	return func() tea.Msg {
		if i.state.selectedBurst == nil {
			return BurstFactsLoadedMsg{Error: fmt.Errorf("no burst selected")}
		}

		facts, err := i.context.Service.GetFactsBySourceBurstID(
			i.context.Context,
			i.state.selectedBurst.ID,
		)
		if err != nil {
			return BurstFactsLoadedMsg{Error: err}
		}

		return BurstFactsLoadedMsg{Facts: facts}
	}
}

// checkForExistingFacts checks if facts already exist for the current burst
func (i *BurstManagementIntent) checkForExistingFacts() tea.Cmd {
	return i.loadFactsForBurst()
}

// extractFacts extracts facts from the burst events and persists them to the database
func (i *BurstManagementIntent) extractFacts() tea.Cmd {
	return func() tea.Msg {
		if i.state.selectedBurst == nil {
			return FactExtractionCompleteMsg{Error: fmt.Errorf("no burst selected")}
		}

		// Extract facts from burst using the service
		facts, err := i.context.Service.ExtractFactsFromBurst(
			i.context.Context,
			i.state.selectedBurst,
		)
		if err != nil {
			return FactExtractionCompleteMsg{Error: err}
		}

		// Persist each extracted fact to the database
		savedFacts := make([]*domain.Fact, 0, len(facts))
		var saveErrors []error

		for idx := range facts {
			fact := &facts[idx]

			// Set source burst ID (linking fact to this burst)
			fact.SourceBurstID = i.state.selectedBurst.ID

			// Save fact to repository
			if err := i.context.Service.SaveFact(i.context.Context, fact); err != nil {
				// Collect errors but continue saving other facts
				saveErrors = append(saveErrors, fmt.Errorf("failed to save fact %d: %w", idx, err))
				continue
			}

			savedFacts = append(savedFacts, fact)
		}

		// If all facts failed to save, return error
		if len(savedFacts) == 0 && len(facts) > 0 {
			return FactExtractionCompleteMsg{
				Error: fmt.Errorf("failed to save any facts: %v", saveErrors),
			}
		}

		// Return saved facts (partial success is OK)
		return FactExtractionCompleteMsg{Facts: savedFacts}
	}
}

// confirmBurstOnly marks the burst as confirmed without extracting facts
func (i *BurstManagementIntent) confirmBurstOnly() tea.Cmd {
	return func() tea.Msg {
		if i.state.selectedBurst == nil {
			return BurstConfirmedMsg{Error: fmt.Errorf("no burst selected")}
		}

		// Mark burst as confirmed
		now := time.Now()
		i.state.selectedBurst.Confirmed = true
		i.state.selectedBurst.ConfirmedAt = &now

		// Update in repository
		err := i.context.UpdateBurst(i.state.selectedBurst)
		if err != nil {
			return BurstConfirmedMsg{Error: err}
		}

		// Reload bursts
		i.context.LoadBursts()
		i.state.filteredBursts = i.context.Bursts

		// Stay in confirm state to show success message
		i.state.currentState = BurstStateConfirm
		i.state.extractionComplete = true

		return BurstConfirmedMsg{Burst: i.state.selectedBurst}
	}
}
