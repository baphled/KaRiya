package intents

import (
	"fmt"
	"strings"
	"time"

	"github.com/baphled/kariya/internal/cli/behaviors"
	"github.com/baphled/kariya/internal/cli/themes"
	"github.com/baphled/kariya/internal/cli/uikit/primitives"
	domain "github.com/baphled/kariya/internal/domain/career"
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
	*BaseIntent

	// context is the input context passed to the intent.
	context *BurstManagementContext

	// state represents the current state of the intent.
	state *BurstManagementIntentModel

	// tableBehavior provides type-safe table operations for bursts
	tableBehavior *behaviors.TableBehavior[*domain.Burst]

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
	deleteError error
	editError   error
	editModal   *EditBurstModal

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

// burstRowFormatter formats a burst for table display
func burstRowFormatter(burst *domain.Burst, index int) []string {
	// Column 1: Name (truncate to 27 chars)
	nameStr := burst.Name
	if len(nameStr) > 27 {
		nameStr = nameStr[:27] + "..."
	}

	// Column 2: Description (truncated preview, max 32 chars)
	descStr := strings.TrimSpace(burst.Description)
	descStr = strings.ReplaceAll(descStr, "\n", " ")
	descStr = strings.ReplaceAll(descStr, "\r", " ")
	if descStr == "" {
		descStr = "-"
	} else if len(descStr) > 32 {
		descStr = descStr[:32] + "..."
	}

	// Column 3: Confirmed Status
	confirmedStr := "✗ No"
	if burst.Confirmed {
		confirmedStr = "✓ Yes"
	}

	// Column 4: Event Count
	eventCount := fmt.Sprintf("%d", len(burst.EventIDs))

	// Column 5: Created Date (YYYY-MM-DD)
	createdStr := burst.CreatedAt.Format("2006-01-02")

	return []string{nameStr, descStr, confirmedStr, eventCount, createdStr}
}

// NewBurstManagementIntent creates a new BurstManagement intent.
func NewBurstManagementIntent(context *BurstManagementContext) (*BurstManagementIntent, error) {
	// Validate the context.
	if err := context.Validate(); err != nil {
		return nil, err
	}

	// Define columns for TableBehavior
	columns := []behaviors.ColumnDef{
		{Title: "Name", Width: 30},
		{Title: "Description", Width: 35},
		{Title: "Confirmed", Width: 10},
		{Title: "Events", Width: 8},
		{Title: "Created", Width: 12},
	}

	// Create TableBehavior with type-safe generics
	tableBehavior := behaviors.NewTableBehavior[*domain.Burst](nil, columns, burstRowFormatter).
		PageSize(15).
		PaginationPrefix("Bursts").
		EmptyMessage("No bursts found.")

	intent := &BurstManagementIntent{
		BaseIntent: NewBaseIntent(),
		context:    context,
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
		tableBehavior: tableBehavior,
		active:        true,
	}

	return intent, nil
}

// syncTableSelection syncs the TableBehavior selection with the intent's data context
func (i *BurstManagementIntent) syncTableSelection() {
	i.state.selectedIndex = i.tableBehavior.GetSelectedIndex()
	if selected := i.tableBehavior.GetSelectedItem(); selected != nil {
		i.state.selectedBurst = *selected
	} else {
		i.state.selectedBurst = nil
	}
}

// Init is called when the intent is activated.
func (i *BurstManagementIntent) Init() tea.Cmd {
	// Apply theme to TableBehavior if available
	if theme := i.Theme(); theme != nil {
		i.tableBehavior.SetTheme(theme)
	}

	// Load bursts from repository (if repository available)
	// Error is acceptable - context.Bursts may already be populated (e.g., in tests)
	if err := i.context.LoadBursts(); err != nil {
		// LoadBursts may fail if no repository configured, but context.Bursts
		// may be pre-populated directly. Only log for debugging if needed.
		_ = err // Acknowledged: intentionally ignored, context.Bursts used as-is
	}

	// Initialize filtered bursts with the provided bursts.
	i.state.filteredBursts = i.context.Bursts
	if len(i.state.filteredBursts) > 0 {
		i.state.selectedBurst = i.state.filteredBursts[0]
	}

	// Set items on TableBehavior
	i.tableBehavior.SetItems(i.state.filteredBursts)
	return nil
}

// getTheme returns the theme or a default.
func (i *BurstManagementIntent) getTheme() themes.Theme {
	if theme := i.Theme(); theme != nil {
		return theme
	}
	return themes.NewDefaultTheme()
}

// getCardStyle returns a themed card style.
func (i *BurstManagementIntent) getCardStyle() lipgloss.Style {
	return i.getTheme().Styles().CardBase
}

// getInfoColor returns the info color from theme.
func (i *BurstManagementIntent) getInfoColor() lipgloss.Color {
	return i.getTheme().InfoColor()
}

// getWarningColor returns the warning color from theme.
func (i *BurstManagementIntent) getWarningColor() lipgloss.Color {
	return i.getTheme().WarningColor()
}

// getErrorColor returns the error color from theme.
func (i *BurstManagementIntent) getErrorColor() lipgloss.Color {
	return i.getTheme().ErrorColor()
}

// getPrimaryColor returns the primary text color from theme.
func (i *BurstManagementIntent) getPrimaryColor() lipgloss.Color {
	return i.getTheme().ForegroundColor()
}

// getBackgroundCardColor returns the card background color from theme.
func (i *BurstManagementIntent) getBackgroundCardColor() lipgloss.Color {
	return i.getTheme().Palette().BackgroundCard
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
		// Handle global keys first (q=quit, ?=help, esc=back)
		switch HandleGlobalKeys(msg) {
		case KeyQuit:
			return tea.Quit
		case KeyHelp:
			i.ToggleHelp()
			return nil
		case KeyBack:
			// At root state, back means cancel and return to main menu
			i.setCancelled()
			return nil
		}

		// Try TableBehavior navigation
		if i.tableBehavior.HandleNavigation(msg.String()) {
			i.syncTableSelection()
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

		case "n":
			// Create new burst
			i.context.StartNewBurst()
			i.state.selectedBurst = i.context.EditingBurst
			i.state.editModal = NewEditBurstModal(i.context.EditingBurst)
			i.state.currentState = BurstStateEdit
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
		// Handle global keys first (q=quit, ?=help, esc=back)
		switch HandleGlobalKeys(msg) {
		case KeyQuit:
			return tea.Quit
		case KeyHelp:
			i.ToggleHelp()
			return nil
		case KeyBack:
			// Go back to list.
			i.state.currentState = BurstStateList
			return nil
		}

		switch msg.String() {
		case "v":
			// View events linked to this burst
			i.state.currentState = BurstStateDetailEvents
			i.state.loadingEvents = true
			return i.loadEventsForBurst()

		case "f":
			// View facts extracted from this burst
			i.state.currentState = BurstStateDetailFacts
			i.state.loadingFacts = true
			return i.loadFactsForBurst()

		case "e":
			// Edit burst (standardized shortcut per TUI_STANDARDS.md)
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
		// Handle global keys first (q=quit, ?=help, esc=back)
		switch HandleGlobalKeys(msg) {
		case KeyQuit:
			return tea.Quit
		case KeyHelp:
			i.ToggleHelp()
			return nil
		case KeyBack:
			// Back to detail
			i.state.currentState = BurstStateDetail
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
		// Handle global keys first (q=quit, ?=help, esc=back)
		switch HandleGlobalKeys(msg) {
		case KeyQuit:
			return tea.Quit
		case KeyHelp:
			i.ToggleHelp()
			return nil
		case KeyBack:
			// Back to detail
			i.state.currentState = BurstStateDetail
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

	// Create the EditBurstModal with the selected burst
	i.state.editModal = NewEditBurstModal(i.state.selectedBurst)

	// Return the form's init command to initialize the huh form
	return i.state.editModal.form.Init()
}

// updateEditView handles the edit state using the EditBurstModal
func (i *BurstManagementIntent) updateEditView(msg tea.Msg) tea.Cmd {
	// If modal is not initialized, handle legacy behavior (fallback)
	if i.state.editModal == nil {
		switch msg := msg.(type) {
		case tea.KeyMsg:
			switch HandleGlobalKeys(msg) {
			case KeyQuit:
				return tea.Quit
			case KeyHelp:
				i.ToggleHelp()
				return nil
			case KeyBack:
				i.state.currentState = BurstStateDetail
				i.state.editError = nil
				return nil
			}
		}
		return nil
	}

	// Check global keys BEFORE delegating to modal
	// This ensures esc, q, ?, m keys work even when modal has focus
	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch HandleGlobalKeys(msg) {
		case KeyQuit:
			return tea.Quit
		case KeyHelp:
			i.ToggleHelp()
			return nil
		case KeyBack:
			// Close modal and return to appropriate state
			i.state.editModal = nil
			if i.context.IsNewBurst {
				i.context.CancelEdit()
				i.state.selectedBurst = nil
				i.state.currentState = BurstStateList
			} else {
				i.state.currentState = BurstStateDetail
			}
			i.state.editError = nil
			return nil
		}
		// Also handle 'm' for main menu explicitly
		if msg.String() == "m" {
			i.state.editModal = nil
			i.setCancelled()
			return nil
		}
	}

	// NOW delegate to the modal for form handling
	cmd := i.state.editModal.Update(msg)

	// Check if modal completed (form submitted or cancelled)
	if result := i.state.editModal.Result(); result != nil {
		if result.Accepted {
			// Apply changes from the modal to the selected burst
			i.state.selectedBurst.Name = result.Modified.Name
			i.state.selectedBurst.Description = result.Modified.Description

			// Determine if we're creating or updating
			var err error
			if i.context.IsNewBurst {
				// Creating a new burst
				err = i.context.CreateBurst(i.state.selectedBurst)
			} else {
				// Updating existing burst
				err = i.context.UpdateBurst(i.state.selectedBurst)
			}

			if err != nil {
				i.state.editError = err
				i.state.editModal = nil
				// Return to appropriate state based on whether it's a new burst
				if i.context.IsNewBurst {
					i.state.currentState = BurstStateList
				} else {
					i.state.currentState = BurstStateDetail
				}
				return nil
			}

			// Reload bursts list to reflect changes
			if err := i.context.LoadBursts(); err != nil {
				i.state.editError = err
				i.state.editModal = nil
				if i.context.IsNewBurst {
					i.state.currentState = BurstStateList
				} else {
					i.state.currentState = BurstStateDetail
				}
				return nil
			}
			i.state.filteredBursts = i.context.Bursts

			// For new burst, select it and go to detail view
			if i.context.IsNewBurst {
				// Find the newly created burst
				for _, b := range i.state.filteredBursts {
					if b.Name == i.state.selectedBurst.Name {
						i.state.selectedBurst = b
						break
					}
				}
				i.context.IsNewBurst = false
			}
		} else {
			// User cancelled - clear new burst state if applicable
			if i.context.IsNewBurst {
				i.context.CancelEdit()
				i.state.selectedBurst = nil
			}
		}

		// Clear modal and return to appropriate view
		i.state.editModal = nil
		if i.context.IsNewBurst || i.state.selectedBurst == nil {
			i.state.currentState = BurstStateList
			// Select first burst if available
			if len(i.state.filteredBursts) > 0 {
				i.state.selectedBurst = i.state.filteredBursts[0]
				i.state.selectedIndex = 0
			}
		} else {
			i.state.currentState = BurstStateDetail
		}
		return nil
	}

	return cmd
}

// updateDeleteConfirmView handles the delete confirmation state
func (i *BurstManagementIntent) updateDeleteConfirmView(msg tea.Msg) tea.Cmd {
	switch msg := msg.(type) {
	case tea.KeyMsg:
		// Handle global keys first (q=quit, ?=help, esc=back)
		switch HandleGlobalKeys(msg) {
		case KeyQuit:
			return tea.Quit
		case KeyHelp:
			i.ToggleHelp()
			return nil
		case KeyBack:
			// Cancel delete
			i.state.currentState = BurstStateDetail
			i.state.deleteError = nil
			return nil
		}

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
				if err := i.context.LoadBursts(); err != nil {
					i.state.deleteError = err
					return nil
				}
				i.state.filteredBursts = i.context.Bursts
				i.state.selectedBurst = nil
				i.state.selectedIndex = 0

				// Go back to list
				i.state.currentState = BurstStateList
				return nil
			}

		case "n":
			// Cancel delete
			i.state.currentState = BurstStateDetail
			i.state.deleteError = nil
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
		// Handle global quit key in all confirm sub-states
		switch HandleGlobalKeys(msg) {
		case KeyQuit:
			return tea.Quit
		case KeyHelp:
			i.ToggleHelp()
			return nil
		}

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
			}
		} else if i.state.extractionComplete {
			// Extraction complete, any key returns to detail
			// Any key returns to detail
			i.state.currentState = BurstStateDetail
			i.state.confirmError = nil
			i.state.extractionComplete = false
			i.state.extractedFactsCount = 0
			return nil
		} else {
			// Handle back (esc)
			if HandleGlobalKeys(msg) == KeyBack {
				// Go back to detail
				i.state.currentState = BurstStateDetail
				i.state.confirmError = nil
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
		// Handle global keys (q=quit, ?=help)
		// Note: esc doesn't go back during extraction
		switch HandleGlobalKeys(msg) {
		case KeyQuit:
			return tea.Quit
		case KeyHelp:
			i.ToggleHelp()
			return nil
		}
	}

	return nil
}

// applyFilters filters the bursts based on current filter state.
// View renders the intent's current state using StandardView.
func (i *BurstManagementIntent) View() string {
	if !i.active {
		return "BurstManagement intent is not active"
	}

	// Create standard view with dynamic breadcrumbs
	view := CreateStandardViewWithBreadcrumbs(i.BaseIntent, i.getBreadcrumbs()...)

	// Handle loading states with modals
	if i.state.loadingEvents {
		i.SetLoading("Loading burst events...")
	} else if i.state.loadingFacts {
		i.SetLoading("Loading burst facts...")
	} else if i.state.extractingFacts {
		message := "Extracting facts from burst events..."
		if i.state.extractedFactsCount > 0 {
			message = fmt.Sprintf("Extracted %d facts so far...", i.state.extractedFactsCount)
		}
		i.SetLoading(message)
	}

	// Handle errors
	if i.state.deleteError != nil {
		i.SetError(fmt.Errorf("failed to delete burst: %w", i.state.deleteError))
	} else if i.state.confirmError != nil {
		i.SetError(fmt.Errorf("failed to confirm burst: %w", i.state.confirmError))
	}

	// Get content for current state
	content := i.getStateContent()
	view.WithContent(content)

	// Get context-aware help
	help := i.getContextHelp()
	view.WithHelp(help).WithFooterSeparator(true)

	return view.Render()
}

// Helper methods for StandardView

func (i *BurstManagementIntent) getBreadcrumbs() []string {
	breadcrumbs := []string{"Main Menu", "Manage Bursts"}

	// Add burst name for detail states
	if i.state.selectedBurst != nil {
		switch i.state.currentState {
		case BurstStateDetail, BurstStateEdit, BurstStateDeleteConfirm, BurstStateConfirm:
			breadcrumbs = append(breadcrumbs, i.state.selectedBurst.Name)
		case BurstStateDetailEvents:
			breadcrumbs = append(breadcrumbs, i.state.selectedBurst.Name, "Events")
		case BurstStateDetailFacts:
			breadcrumbs = append(breadcrumbs, i.state.selectedBurst.Name, "Facts")
		case BurstStateExtractingFacts:
			breadcrumbs = append(breadcrumbs, i.state.selectedBurst.Name, "Extracting Facts")
		}
	}

	return breadcrumbs
}

func (i *BurstManagementIntent) getStateContent() string {
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

func (i *BurstManagementIntent) getContextHelp() string {
	theme := i.Theme()

	switch i.state.currentState {
	case BurstStateList:
		return CombineThemedFooters(
			ThemedListFooter(theme),
			ThemedCustomFooter(theme,
				primitives.HelpKeyBadge("Enter", "View details", theme),
				primitives.HelpKeyBadge("n", "New burst", theme),
			),
			ThemedGlobalBadges(theme),
		)
	case BurstStateDetail:
		return CombineThemedFooters(
			ThemedDetailViewFooter(theme),
			ThemedCustomFooter(theme,
				primitives.HelpKeyBadge("v", "View events", theme),
				primitives.HelpKeyBadge("f", "View facts", theme),
				primitives.EditBadge(theme),
				primitives.DeleteBadge(theme),
				primitives.HelpKeyBadge("c", "Confirm", theme),
			),
			ThemedGlobalBadges(theme),
		)
	case BurstStateDetailEvents, BurstStateDetailFacts:
		return CombineThemedFooters(
			ThemedDetailViewFooter(theme),
			ThemedGlobalBadges(theme),
		)
	case BurstStateEdit:
		return CombineThemedFooters(
			ThemedFormFooter(theme),
			ThemedCustomFooter(theme,
				primitives.SaveBadge(theme),
			),
			ThemedGlobalBadges(theme),
		)
	case BurstStateDeleteConfirm:
		return CombineThemedFooters(
			ThemedCustomFooter(theme,
				primitives.HelpKeyBadge("y/Enter", "Confirm deletion", theme),
				primitives.HelpKeyBadge("n/Esc", "Cancel", theme),
			),
			ThemedGlobalBadges(theme),
		)
	case BurstStateConfirm:
		if i.state.extractionComplete {
			return CombineThemedFooters(
				ThemedCustomFooter(theme,
					primitives.HelpKeyBadge("Enter", "Continue", theme),
					primitives.BackBadge(theme),
				),
				ThemedGlobalBadges(theme),
			)
		}
		return CombineThemedFooters(
			ThemedCustomFooter(theme,
				primitives.HelpKeyBadge("y/Enter", "Confirm burst", theme),
				primitives.HelpKeyBadge("n/Esc", "Cancel", theme),
			),
			ThemedGlobalBadges(theme),
		)
	case BurstStateExtractingFacts:
		return CombineThemedFooters(
			ThemedCustomFooter(theme,
				primitives.HelpKeyBadge("...", "Please wait", theme),
			),
			ThemedGlobalBadges(theme),
		)
	default:
		return ThemedGlobalBadges(theme)
	}
}

// viewList renders the burst list view with all bursts as a table.
func (i *BurstManagementIntent) viewList() string {
	return i.tableBehavior.Render()
}
func (i *BurstManagementIntent) viewDetail() string {
	if i.state.selectedBurst == nil {
		return "No burst selected."
	}

	var content strings.Builder

	// Title with confirmation status indicator
	title := "Burst Details"
	if i.state.selectedBurst.Confirmed {
		confirmedBadge := primitives.SuccessText("✓ Confirmed", i.Theme()).Bold().Render()
		title = "Burst Details " + confirmedBadge
	}
	content.WriteString("\n" + primitives.NewText(title, i.Theme()).Bold().Foreground(i.getPrimaryColor()).Render() + "\n\n")

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

	// Apply themed card styling
	card := i.getCardStyle().Render(content.String())

	// Footer now handled by StandardView
	return card
}

// viewDetailEvents renders the events view for a burst
func (i *BurstManagementIntent) viewDetailEvents() string {
	if i.state.loadingEvents {
		infoStyle := lipgloss.NewStyle().Foreground(i.getInfoColor())
		return infoStyle.Render("Loading events...")
	}

	if len(i.state.burstEvents) == 0 {
		errorStyle := lipgloss.NewStyle().Foreground(i.getErrorColor())
		return errorStyle.Render("No events found for this burst.")
	}

	var content strings.Builder
	header := primitives.NewText(
		fmt.Sprintf("Events in Burst: %s", i.state.selectedBurst.Name),
		i.Theme(),
	).Bold().Foreground(i.getPrimaryColor()).MarginBottom(1)
	content.WriteString(header.Render())
	content.WriteString("\n\n")

	for idx, event := range i.state.burstEvents {
		content.WriteString(fmt.Sprintf("%d. %s\n", idx+1, event.Text))
		content.WriteString(fmt.Sprintf("   Date: %s\n", event.Date.Format("2006-01-02")))
		if len(event.Tags) > 0 {
			content.WriteString(fmt.Sprintf("   Tags: %s\n", strings.Join(event.Tags, ", ")))
		}
		content.WriteString("\n")
	}

	// Footer now handled by StandardView
	return content.String()
}

// viewDetailFacts renders the facts view for a burst
func (i *BurstManagementIntent) viewDetailFacts() string {
	if i.state.loadingFacts {
		infoStyle := lipgloss.NewStyle().Foreground(i.getInfoColor())
		return infoStyle.Render("Loading facts...")
	}

	if len(i.state.burstFacts) == 0 {
		warningStyle := lipgloss.NewStyle().
			Foreground(i.getWarningColor()).
			MarginBottom(1)

		return warningStyle.Render("No facts extracted yet. Confirm the burst in detail view to extract facts.")
	}

	var content strings.Builder
	header := primitives.NewText(
		fmt.Sprintf("Facts from Burst: %s", i.state.selectedBurst.Name),
		i.Theme(),
	).Bold().Foreground(i.getPrimaryColor()).MarginBottom(1)
	content.WriteString(header.Render())
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

	// Footer now handled by StandardView
	return content.String()
}

// viewEdit renders the edit view for a burst using the EditBurstModal
func (i *BurstManagementIntent) viewEdit() string {
	if i.state.selectedBurst == nil {
		return "No burst selected."
	}

	// If modal is active, render just the form content (not the full modal container)
	// This avoids duplicate help text since StandardView provides context-aware help
	if i.state.editModal != nil {
		return i.state.editModal.GetContent()
	}

	// Fallback view if modal not initialized (should not happen normally)
	var content strings.Builder

	// Header
	content.WriteString(primitives.NewText("Edit Burst", i.Theme()).Bold().Foreground(i.getPrimaryColor()).MarginBottom(1).Render())
	content.WriteString("\n\n")

	// Show error if any
	if i.state.editError != nil {
		content.WriteString(primitives.ErrorText(fmt.Sprintf("Error: %s", i.state.editError), i.Theme()).MarginBottom(1).Render())
		content.WriteString("\n\n")
	}

	// Display current burst details
	content.WriteString(fmt.Sprintf("Name: %s\n", i.state.selectedBurst.Name))

	if i.state.selectedBurst.Description != "" {
		content.WriteString(fmt.Sprintf("Description: %s\n", i.state.selectedBurst.Description))
	}

	content.WriteString("\n")
	content.WriteString("Initializing edit form...")
	content.WriteString("\n")

	// Apply themed card styling
	card := i.getCardStyle().Render(content.String())

	return card
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
		Foreground(i.getErrorColor()).
		MarginBottom(1)

	content.WriteString(warningStyle.Render("⚠️  DELETE BURST - Are you sure?"))
	content.WriteString("\n\n")

	// Show error if any
	if i.state.deleteError != nil {
		errorStyle := lipgloss.NewStyle().
			Foreground(i.getErrorColor()).
			MarginBottom(1)
		content.WriteString(errorStyle.Render(fmt.Sprintf("Error: %s", i.state.deleteError)))
		content.WriteString("\n\n")
	}

	// Burst details
	infoStyle := lipgloss.NewStyle().
		Foreground(i.getPrimaryColor())

	content.WriteString(infoStyle.Render(fmt.Sprintf("Burst Name: %s", i.state.selectedBurst.Name)))
	content.WriteString("\n\n")

	// Warning messages
	warningTextStyle := lipgloss.NewStyle().
		Foreground(i.getWarningColor()).
		Bold(true)

	content.WriteString(warningTextStyle.Render("This will remove the burst grouping but NOT delete the events."))
	content.WriteString("\n")
	content.WriteString(warningTextStyle.Render("This action cannot be undone."))
	content.WriteString("\n")

	// Apply themed card styling with error border
	var cardStyle lipgloss.Style
	if theme := i.Theme(); theme != nil {
		// Use assignment instead of deprecated Copy()
		cardStyle = theme.Styles().CardBase.BorderForeground(theme.ErrorColor())
	} else {
		cardStyle = lipgloss.NewStyle().
			Padding(1, 2).
			BorderStyle(lipgloss.RoundedBorder()).
			BorderForeground(i.getErrorColor()).
			Background(i.getBackgroundCardColor()).
			Foreground(i.getPrimaryColor())
	}

	card := cardStyle.Render(content.String())

	// Footer now handled by StandardView
	return card
}

// viewConfirm renders the confirmation view
func (i *BurstManagementIntent) viewConfirm() string {
	if i.state.selectedBurst == nil {
		return "No burst selected."
	}

	var content strings.Builder

	// Header
	content.WriteString(primitives.SuccessText("Confirm Burst", i.Theme()).Bold().MarginBottom(1).Render())
	content.WriteString("\n\n")

	// Show error if any
	if i.state.confirmError != nil {
		content.WriteString(primitives.ErrorText(fmt.Sprintf("Error: %s", i.state.confirmError), i.Theme()).MarginBottom(1).Render())
		content.WriteString("\n\n")
	}

	// Burst details
	content.WriteString(primitives.NewText(fmt.Sprintf("Burst: %s", i.state.selectedBurst.Name), i.Theme()).Foreground(i.getPrimaryColor()).Render())
	content.WriteString("\n\n")

	// Show different messages based on state
	if i.state.showReextractPrompt {
		// Facts already exist
		content.WriteString(primitives.WarningText(fmt.Sprintf("This burst already has %d facts extracted.", i.state.existingFactsCount), i.Theme()).Bold().Render())
		content.WriteString("\n\n")
		content.WriteString(primitives.NewText("Do you want to extract more facts? New facts will be added to existing ones.", i.Theme()).Foreground(i.getPrimaryColor()).Render())
		content.WriteString("\n")
	} else if i.state.extractionComplete {
		// Extraction completed successfully
		content.WriteString(primitives.SuccessText(fmt.Sprintf("✓ Successfully extracted and saved %d facts!", i.state.extractedFactsCount), i.Theme()).Bold().Render())
		content.WriteString("\n\n")
		content.WriteString(primitives.NewText("Burst has been confirmed.", i.Theme()).Foreground(i.getPrimaryColor()).Render())
		content.WriteString("\n")
	} else {
		// About to start extraction
		content.WriteString(primitives.NewText("No facts found for this burst. Starting fact extraction...", i.Theme()).Foreground(i.getPrimaryColor()).Render())
		content.WriteString("\n")
	}

	// Apply themed card styling
	card := i.getCardStyle().Render(content.String())

	// Footer now handled by StandardView
	return card
}

// viewExtractingFacts renders the fact extraction progress view
func (i *BurstManagementIntent) viewExtractingFacts() string {
	if i.state.selectedBurst == nil {
		return "No burst selected."
	}

	var content strings.Builder

	// Header
	content.WriteString(primitives.NewText("Extracting Facts", i.Theme()).Bold().Foreground(i.getInfoColor()).MarginBottom(1).Render())
	content.WriteString("\n\n")

	// Progress indicator
	content.WriteString(primitives.NewText(fmt.Sprintf("Burst: %s", i.state.selectedBurst.Name), i.Theme()).Foreground(i.getPrimaryColor()).Render())
	content.WriteString("\n\n")

	content.WriteString(primitives.NewText("⏳ Extracting and saving facts...", i.Theme()).Bold().Foreground(i.getInfoColor()).Render())
	content.WriteString("\n\n")
	content.WriteString(primitives.NewText("Analyzing events and persisting facts to database.", i.Theme()).Foreground(i.getPrimaryColor()).Render())
	content.WriteString("\n\n")
	content.WriteString(primitives.NewText("This may take a few moments.", i.Theme()).Foreground(i.getPrimaryColor()).Render())
	content.WriteString("\n")

	// Apply themed card styling
	card := i.getCardStyle().Render(content.String())

	// Footer now handled by StandardView
	return card
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

// ListNavigator interface implementation (delegates to TableBehavior)
func (i *BurstManagementIntent) GetTotalItems() int {
	return i.tableBehavior.Count()
}

func (i *BurstManagementIntent) GetSelectedIndex() int {
	return i.tableBehavior.GetSelectedIndex()
}

func (i *BurstManagementIntent) SetSelectedIndex(idx int) {
	i.tableBehavior.SetSelectedIndex(idx)
	i.syncTableSelection()
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
		if err := i.context.LoadBursts(); err != nil {
			return BurstConfirmedMsg{Error: err}
		}
		i.state.filteredBursts = i.context.Bursts

		// Stay in confirm state to show success message
		i.state.currentState = BurstStateConfirm
		i.state.extractionComplete = true

		return BurstConfirmedMsg{Burst: i.state.selectedBurst}
	}
}

// SetTestModalResult sets the edit modal's result directly for testing purposes.
// This allows E2E tests to simulate modal completion without interacting with the huh form.
// Only use this method in tests.
func (i *BurstManagementIntent) SetTestModalResult(result *ModalEditResult[*domain.Burst]) {
	if i.state.editModal != nil {
		i.state.editModal.SetTestResult(result)
	}
}
