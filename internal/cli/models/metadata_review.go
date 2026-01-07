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

// BulkOperationsMsg is sent to open bulk operations for selected events
type BulkOperationsMsg struct {
	Events []*career.CareerEvent
}

// BurstSuggestionsTriggeredMsg is sent when burst suggestions should be reviewed
type BurstSuggestionsTriggeredMsg struct {
	EventIDs []string
}

// MetadataReviewModel represents the metadata review screen
type MetadataReviewModel struct {
	*BaseStandardModel
	service              *careerservice.Service
	calculator           *careerservice.DataQualityCalculator
	ctx                  context.Context
	events               []*career.CareerEvent
	filtered             []*career.CareerEvent
	qualityScores        map[string]*careerservice.QualityScore
	table                table.Model
	pagination           *PaginationHelper
	width                int
	height               int
	err                  error
	expandedIdx          int                            // Index of expanded event (-1 if none)
	filterMode           string                         // "all", "incomplete"
	sortBy               string                         // "date", "company", "quality"
	importedEventIDs     map[string]bool                // IDs of recently imported events
	isImportReview       bool                           // True if reviewing only imported events
	fieldOrigins         map[string]map[string]bool     // eventID -> field -> isFromCSV
	parsingWarnings      map[string][]string            // eventID -> warnings
	duplicateStatus      map[string]string              // eventID -> original event ID (empty if not duplicate)
	helpFooter           components.HelpFooterModel     // Help footer
	header               components.HeaderModel         // Header component
	breadcrumbs          []string                       // Navigation breadcrumb trail
	deletionState        *ListDeletionState             // Deletion state management
	navigationKeyHandler *ListNavigationKeyHandler      // Navigation key handler
	listContainer        *components.TableListContainer // Table list container
	selectedEvents       map[string]bool                // Selected events for bulk operations
}

// NewMetadataReviewModel creates a new metadata review model
func NewMetadataReviewModel(svc *careerservice.Service, ctx context.Context) *MetadataReviewModel {
	calculator := careerservice.NewDataQualityCalculator()

	columns := []table.Column{
		{Title: "Event", Width: 50},
		{Title: "Quality", Width: 15},
		{Title: "Company", Width: 20},
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

	model := &MetadataReviewModel{
		BaseStandardModel: NewBaseStandardModel(),
		service:           svc,
		calculator:        calculator,
		ctx:               ctx,
		events:            []*career.CareerEvent{},
		filtered:          []*career.CareerEvent{},
		table:             t,
		pagination:        NewPaginationHelper(10),
		expandedIdx:       -1,
		filterMode:        "all",
		sortBy:            "quality",
		qualityScores:     make(map[string]*careerservice.QualityScore),
		fieldOrigins:      make(map[string]map[string]bool),
		parsingWarnings:   make(map[string][]string),
		duplicateStatus:   make(map[string]string),
		width:             80,
		height:            20,
		helpFooter:        components.NewHelpFooter("metadata_review", 80),
		header:            components.NewHeader("Metadata Review", 80),
		breadcrumbs:       []string{"Home", "Metadata Review"},
		deletionState:     NewListDeletionState(),
		listContainer:     components.NewTableListContainer(t, "Metadata Review", 80),
		selectedEvents:    make(map[string]bool),
	}

	// Create navigation key handler with callbacks
	model.navigationKeyHandler = NewListNavigationKeyHandler(model)

	// Load events
	model.loadEvents()

	return model
}

// NewMetadataReviewModelForImport creates a metadata review model for imported events
func NewMetadataReviewModelForImport(svc *careerservice.Service, ctx context.Context, importedEventIDs []string) *MetadataReviewModel {
	calculator := careerservice.NewDataQualityCalculator()

	// Convert slice to map for O(1) lookup
	importedMap := make(map[string]bool)
	for _, id := range importedEventIDs {
		importedMap[id] = true
	}

	columns := []table.Column{
		{Title: "Event", Width: 50},
		{Title: "Quality", Width: 15},
		{Title: "Company", Width: 20},
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

	model := &MetadataReviewModel{
		BaseStandardModel: NewBaseStandardModel(),
		service:           svc,
		calculator:        calculator,
		ctx:               ctx,
		events:            []*career.CareerEvent{},
		filtered:          []*career.CareerEvent{},
		table:             t,
		pagination:        NewPaginationHelper(10),
		expandedIdx:       -1,
		filterMode:        "all",
		sortBy:            "quality",
		qualityScores:     make(map[string]*careerservice.QualityScore),
		importedEventIDs:  importedMap,
		isImportReview:    true,
		fieldOrigins:      make(map[string]map[string]bool),
		parsingWarnings:   make(map[string][]string),
		duplicateStatus:   make(map[string]string),
		width:             80,
		height:            20,
		helpFooter:        components.NewHelpFooter("metadata_review", 80),
		header:            components.NewHeader("Metadata Review (Imported)", 80),
		breadcrumbs:       []string{"Home", "Metadata Review"},
		deletionState:     NewListDeletionState(),
		listContainer:     components.NewTableListContainer(t, "Metadata Review (Imported)", 80),
		selectedEvents:    make(map[string]bool),
	}

	// Create navigation key handler with callbacks
	model.navigationKeyHandler = NewListNavigationKeyHandler(model)

	// Load events
	model.loadEvents()

	return model
}

// loadEvents loads events from the service and applies filters and sorting
func (m *MetadataReviewModel) loadEvents() {
	filters := careerrepo.ListFilters{
		SortBy:    "date",
		SortOrder: "desc",
		Limit:     1000,
	}

	events, err := m.service.ListEvents(m.ctx, filters)
	if err != nil {
		m.err = err
		m.events = []*career.CareerEvent{}
		return
	}

	// Calculate quality scores for all events
	m.qualityScores = make(map[string]*careerservice.QualityScore)
	for _, event := range events {
		score := m.calculator.CalculateQuality(event)
		m.qualityScores[event.ID] = &score
	}

	m.events = events
	m.applyFiltersAndSort()
	m.pagination.SetTotalCount(len(m.filtered))
	m.updateTableRows()
	m.err = nil
}

// applyFiltersAndSort applies filters and sorting to events
func (m *MetadataReviewModel) applyFiltersAndSort() {
	m.filtered = m.filterEvents()
	m.sortEvents()
}

// filterEvents filters events based on current filter mode
func (m *MetadataReviewModel) filterEvents() []*career.CareerEvent {
	var filtered []*career.CareerEvent

	// Filter by imported events if in import review mode
	if m.isImportReview && len(m.importedEventIDs) > 0 {
		for _, event := range m.events {
			if m.importedEventIDs[event.ID] {
				filtered = append(filtered, event)
			}
		}
	} else {
		filtered = m.events
	}

	// Apply quality filter if set
	if m.filterMode == "incomplete" {
		incompleteFiltered := make([]*career.CareerEvent, 0)
		for _, event := range filtered {
			score := m.qualityScores[event.ID]
			if score != nil && (score.Level == careerservice.QualityIncomplete || score.Level == careerservice.QualityBasic) {
				incompleteFiltered = append(incompleteFiltered, event)
			}
		}
		return incompleteFiltered
	}

	return filtered
}

// sortEvents sorts events based on current sort mode
func (m *MetadataReviewModel) sortEvents() {
	switch m.sortBy {
	case "quality":
		// Sort by quality score (ascending - worst first)
		sort.SliceStable(m.filtered, func(i, j int) bool {
			scoreI := m.qualityScores[m.filtered[i].ID]
			scoreJ := m.qualityScores[m.filtered[j].ID]
			if scoreI != nil && scoreJ != nil {
				return scoreI.Score < scoreJ.Score
			}
			return false
		})
	case "company":
		// Sort by company name
		sort.SliceStable(m.filtered, func(i, j int) bool {
			return m.filtered[i].Company < m.filtered[j].Company
		})
	case "date":
		// Sort by date
		sort.SliceStable(m.filtered, func(i, j int) bool {
			return m.filtered[i].Date.Before(m.filtered[j].Date)
		})
	}
}

// updateTableRows updates the table with rows from filtered events (current page only)
func (m *MetadataReviewModel) updateTableRows() {
	var rows []table.Row
	pageEvents := m.getPageEvents()

	// Get the selected index within the current page
	selectedIdx := m.listContainer.GetSelectedIdx()

	// Ensure cursor is within valid bounds BEFORE creating rows
	if len(pageEvents) > 0 {
		if selectedIdx >= len(pageEvents) {
			selectedIdx = len(pageEvents) - 1
		}
		if selectedIdx < 0 {
			selectedIdx = 0
		}
	} else {
		selectedIdx = 0
	}

	for i, event := range pageEvents {
		eventText := event.Text
		if len(eventText) > 50 {
			eventText = eventText[:47] + "..."
		}

		// Add focus indicator for the selected row, or spaces for alignment
		if i == selectedIdx {
			eventText = "▶ " + eventText
		} else {
			eventText = "  " + eventText
		}

		// Get quality score and render it
		score := m.qualityScores[event.ID]
		qualityStr := "N/A"
		if score != nil {
			qualityStr = fmt.Sprintf("%s (%d%%)", score.Level, score.Score)
		}

		company := event.Company
		if len(company) > 17 {
			company = company[:14] + "..."
		}

		row := table.Row{
			eventText,
			qualityStr,
			company,
		}
		rows = append(rows, row)
	}

	m.table.SetRows(rows)
	// Sync the container and table cursor with the validated selectedIdx
	m.listContainer.SetSelectedIdx(selectedIdx)
	m.table.SetCursor(selectedIdx)
}

// getPageEvents returns the events for the current page
func (m *MetadataReviewModel) getPageEvents() []*career.CareerEvent {
	if m.pagination.GetTotalCount() == 0 {
		return []*career.CareerEvent{}
	}

	startIdx := m.pagination.GetPageStartIndex()
	endIdx := m.pagination.GetPageEndIndex()

	if startIdx >= len(m.filtered) {
		return []*career.CareerEvent{}
	}

	if endIdx > len(m.filtered) {
		endIdx = len(m.filtered)
	}

	return m.filtered[startIdx:endIdx]
}

// Init initializes the model
func (m *MetadataReviewModel) Init() tea.Cmd {
	return nil
}

// Update handles messages
func (m *MetadataReviewModel) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	// Handle deletion confirmation if active
	if m.deletionState.IsConfirming() {
		cmd := m.deletionState.UpdateConfirmation(msg)

		if m.deletionState.IsConfirmed() {
			return m.performEventDeletion()
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
		case "f":
			// Toggle filter mode
			if m.filterMode == "all" {
				m.filterMode = "incomplete"
			} else {
				m.filterMode = "all"
			}
			m.applyFiltersAndSort()
			m.pagination.SetTotalCount(len(m.filtered)).GoToFirstPage()
			m.listContainer.MoveToFirst()
			m.updateTableRows()
		case "s":
			// Cycle through sort modes
			switch m.sortBy {
			case "quality":
				m.sortBy = "date"
			case "date":
				m.sortBy = "company"
			case "company":
				m.sortBy = "quality"
			}
			m.applyFiltersAndSort()
			m.updateTableRows()
		case "b":
			// Trigger bulk operations on all events
			return m, func() tea.Msg {
				return BulkOperationsMsg{Events: m.filtered}
			}
		case "u":
			// Trigger burst suggestions for all events
			if len(m.filtered) >= 2 {
				eventIDs := make([]string, len(m.filtered))
				for i, event := range m.filtered {
					eventIDs[i] = event.ID
				}
				return m, func() tea.Msg {
					return BurstSuggestionsTriggeredMsg{EventIDs: eventIDs}
				}
			}
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
		m.helpFooter.SetWidth(msg.Width)
		m.header.SetWidth(msg.Width)
		m.updateTableRows()
		return m, nil
	}

	var cmd tea.Cmd
	m.table, cmd = m.table.Update(msg)
	return m, cmd
}

// performEventDeletion performs the actual event deletion
func (m *MetadataReviewModel) performEventDeletion() (tea.Model, tea.Cmd) {
	if m.deletionState.DeletingItemID == "" {
		m.deletionState.Clear()
		return m, nil
	}

	err := m.service.DeleteEvent(m.ctx, m.deletionState.DeletingItemID)
	if err != nil {
		m.deletionState.SetErrorMsg(fmt.Sprintf("Error deleting event: %v", err))
		return m, nil
	}

	m.removeEventFromLists(m.deletionState.DeletingItemID)

	m.deletionState.SetSuccessMsg("Event deleted successfully")

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

// removeEventFromLists removes an event from both the events and filtered lists
func (m *MetadataReviewModel) removeEventFromLists(eventID string) {
	for i, event := range m.events {
		if event.ID == eventID {
			m.events = append(m.events[:i], m.events[i+1:]...)
			break
		}
	}

	for i, event := range m.filtered {
		if event.ID == eventID {
			m.filtered = append(m.filtered[:i], m.filtered[i+1:]...)
			break
		}
	}

	delete(m.selectedEvents, eventID)
	delete(m.qualityScores, eventID)
}

// View renders the model
func (m *MetadataReviewModel) View() string {
	// Render deletion confirmation dialog
	if m.deletionState.IsConfirming() {
		return m.deletionState.ConfirmationDialog.View()
	}

	// Render deletion message (success or error)
	if m.deletionState.ShowMessage {
		return m.renderDeletionMessage()
	}

	if m.err != nil {
		errorMsg := fmt.Sprintf("Error loading events: %v\n\nPress 'r' to retry or 'esc' to cancel", m.err)
		m.listContainer.SetErrorMessage(errorMsg).SetDimensions(m.width, m.height)
		return m.listContainer.Render()
	}

	// Update list container with current state
	m.listContainer.SetTable(m.table).
		SetDimensions(m.width, m.height).
		SetEmptyStateMessage("No events found")

	// Set pagination info
	paginationText := m.pagination.GetPaginationInfo("events")
	m.listContainer.SetPaginationInfo(paginationText)

	return m.listContainer.Render()
}

// renderDeletionMessage renders the deletion success/error message
func (m *MetadataReviewModel) renderDeletionMessage() string {
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

// ListItemCallbacks implementation for MetadataReviewModel
// These methods implement the ListItemCallbacks interface required by ListNavigationKeyHandler

// OnView sends a message to view the selected event
func (m *MetadataReviewModel) OnView(item interface{}) tea.Cmd {
	if event, ok := item.(*career.CareerEvent); ok {
		return func() tea.Msg { return EventActionSelectedMsg{Event: event, Action: EventActionView} }
	}
	return nil
}

// OnEdit sends a message to edit the selected event
func (m *MetadataReviewModel) OnEdit(item interface{}) tea.Cmd {
	if event, ok := item.(*career.CareerEvent); ok {
		return func() tea.Msg { return EditEventMsg{Event: event} }
	}
	return nil
}

// OnDelete initiates deletion of the selected event
func (m *MetadataReviewModel) OnDelete(item interface{}) {
	if event, ok := item.(*career.CareerEvent); ok {
		m.deletionState.ShowConfirmation("Event", event.Text)
		m.deletionState.DeletingItemID = event.ID
	}
}

// OnToggleSelection toggles the selection state of an item
func (m *MetadataReviewModel) OnToggleSelection(item interface{}) {
	if event, ok := item.(*career.CareerEvent); ok {
		m.selectedEvents[event.ID] = !m.selectedEvents[event.ID]
	}
}

// HasSelectedItem returns the currently selected event
func (m *MetadataReviewModel) HasSelectedItem() interface{} {
	return m.GetSelectedEvent()
}

// MoveUp moves the selection up by count items
func (m *MetadataReviewModel) MoveUp(count int) {
	m.listContainer.MoveUp(count)
}

// MoveDown moves the selection down by count items
func (m *MetadataReviewModel) MoveDown(count int) {
	m.listContainer.MoveDown(count)
}

// MoveToFirst moves selection to the first item
func (m *MetadataReviewModel) MoveToFirst() {
	m.listContainer.MoveToFirst()
}

// MoveToLast moves selection to the last item
func (m *MetadataReviewModel) MoveToLast() {
	m.listContainer.MoveToLast()
}

// UpdateDisplay refreshes the table display
func (m *MetadataReviewModel) UpdateDisplay() {
	m.updateTableRows()
}

// GetRowCount returns the number of visible rows
func (m *MetadataReviewModel) GetRowCount() int {
	return len(m.table.Rows())
}

// GetCurrentIndex returns the current selection index
func (m *MetadataReviewModel) GetCurrentIndex() int {
	return m.listContainer.GetSelectedIdx()
}

// GetPageSize returns the page size for pagination
func (m *MetadataReviewModel) GetPageSize() int {
	return m.pagination.GetPageSize()
}

// nextPage moves to the next page of results
func (m *MetadataReviewModel) nextPage() {
	m.pagination.NextPage()
	m.listContainer.MoveToFirst()
	m.updateTableRows()
}

// prevPage moves to the previous page of results
func (m *MetadataReviewModel) prevPage() {
	m.pagination.PrevPage()
	m.listContainer.MoveToFirst()
	m.updateTableRows()
}

// goToFirstItem moves to the first item
func (m *MetadataReviewModel) goToFirstItem() {
	m.pagination.GoToFirstPage()
	m.listContainer.MoveToFirst()
}

// goToLastItem moves to the last item
func (m *MetadataReviewModel) goToLastItem() {
	m.pagination.GoToLastPage()
	pageEvents := m.getPageEvents()
	if len(pageEvents) > 0 {
		m.listContainer.SetSelectedIdx(len(pageEvents) - 1)
	}
}

// Accessor methods for data and state

// GetSelectedEvent returns the currently selected event
func (m *MetadataReviewModel) GetSelectedEvent() *career.CareerEvent {
	cursor := m.listContainer.GetSelectedIdx()
	if cursor >= 0 && cursor < len(m.filtered) {
		return m.filtered[cursor]
	}
	return nil
}

// GetEvents returns the currently filtered events
func (m *MetadataReviewModel) GetEvents() []*career.CareerEvent {
	return m.filtered
}

// Refresh reloads events from service
func (m *MetadataReviewModel) Refresh() {
	m.loadEvents()
	if m.listContainer.GetSelectedIdx() >= len(m.filtered) {
		m.listContainer.SetSelectedIdx(0)
	}
	if m.listContainer.GetSelectedIdx() < 0 {
		m.listContainer.SetSelectedIdx(0)
	}
}

// SetFieldOrigins sets the field origins for an imported event
func (m *MetadataReviewModel) SetFieldOrigins(eventID string, origins map[string]bool) {
	if m.fieldOrigins == nil {
		m.fieldOrigins = make(map[string]map[string]bool)
	}
	m.fieldOrigins[eventID] = origins
}

// GetFieldOrigins returns the field origins for an event
func (m *MetadataReviewModel) GetFieldOrigins(eventID string) map[string]bool {
	if m.fieldOrigins == nil {
		return nil
	}
	return m.fieldOrigins[eventID]
}

// IsFieldFromCSV returns whether a field came from CSV (true) or is a default (false)
func (m *MetadataReviewModel) IsFieldFromCSV(eventID, field string) bool {
	origins := m.GetFieldOrigins(eventID)
	if origins == nil {
		return false // Default assumption if no origin info
	}
	fromCSV, exists := origins[field]
	return exists && fromCSV
}

// GetDefaultFields returns the list of fields that are defaults for an event
func (m *MetadataReviewModel) GetDefaultFields(eventID string) []string {
	origins := m.GetFieldOrigins(eventID)
	if origins == nil {
		return []string{}
	}
	var defaults []string
	for field, fromCSV := range origins {
		if !fromCSV {
			defaults = append(defaults, field)
		}
	}
	return defaults
}

// HasDefaultFields checks if an event has any default fields
func (m *MetadataReviewModel) HasDefaultFields(eventID string) bool {
	return len(m.GetDefaultFields(eventID)) > 0
}

// SetParsingWarnings sets the parsing warnings for an imported event
func (m *MetadataReviewModel) SetParsingWarnings(eventID string, warnings []string) {
	if m.parsingWarnings == nil {
		m.parsingWarnings = make(map[string][]string)
	}
	m.parsingWarnings[eventID] = warnings
}

// GetParsingWarnings returns the parsing warnings for an event
func (m *MetadataReviewModel) GetParsingWarnings(eventID string) []string {
	if m.parsingWarnings == nil {
		return []string{}
	}
	warnings, exists := m.parsingWarnings[eventID]
	if !exists {
		return []string{}
	}
	return warnings
}

// HasParsingWarnings checks if an event has parsing warnings
func (m *MetadataReviewModel) HasParsingWarnings(eventID string) bool {
	return len(m.GetParsingWarnings(eventID)) > 0
}

// SetDuplicateStatus sets the duplicate status for an event
func (m *MetadataReviewModel) SetDuplicateStatus(eventID string, isDuplicate bool, originalEventID string) {
	if m.duplicateStatus == nil {
		m.duplicateStatus = make(map[string]string)
	}
	if isDuplicate {
		m.duplicateStatus[eventID] = originalEventID
	} else {
		m.duplicateStatus[eventID] = ""
	}
}

// GetDuplicateStatus returns whether an event is a duplicate and its original event ID
func (m *MetadataReviewModel) GetDuplicateStatus(eventID string) (bool, string) {
	if m.duplicateStatus == nil {
		return false, ""
	}
	originalID, exists := m.duplicateStatus[eventID]
	if !exists || originalID == "" {
		return false, ""
	}
	return true, originalID
}

// IsDuplicate checks if an event is a duplicate
func (m *MetadataReviewModel) IsDuplicate(eventID string) bool {
	isDuplicate, _ := m.GetDuplicateStatus(eventID)
	return isDuplicate
}

// GetOriginalEventID returns the original event ID if this is a duplicate
func (m *MetadataReviewModel) GetOriginalEventID(eventID string) string {
	_, originalID := m.GetDuplicateStatus(eventID)
	return originalID
}

// SetBreadcrumbs sets breadcrumb trail for display in header
func (m *MetadataReviewModel) SetBreadcrumbs(crumbs []string) {
	m.breadcrumbs = crumbs
	// Note: header.SetBreadcrumbs removed - breadcrumbs now handled by StandardView
}

// NextPage implements ListItemCallbacks interface - moves to the next page
func (m *MetadataReviewModel) NextPage() {
	m.nextPage()
}

// PrevPage implements ListItemCallbacks interface - moves to the previous page
func (m *MetadataReviewModel) PrevPage() {
	m.prevPage()
}

// GoToFirstPage implements ListItemCallbacks interface - goes to the first page
func (m *MetadataReviewModel) GoToFirstPage() {
	m.goToFirstItem()
	m.updateTableRows()
}

// GoToLastPage implements ListItemCallbacks interface - goes to the last page
func (m *MetadataReviewModel) GoToLastPage() {
	m.goToLastItem()
	m.updateTableRows()
}
