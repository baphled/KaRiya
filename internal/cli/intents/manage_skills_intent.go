package intents

import (
	"fmt"
	"sort"
	"time"

	"github.com/baphled/kariya/internal/cli/components"
	"github.com/baphled/kariya/internal/cli/forms"
	"github.com/baphled/kariya/internal/cli/models"
	"github.com/baphled/kariya/internal/cli/navigation"
	"github.com/baphled/kariya/internal/cli/screens"
	skills_screens "github.com/baphled/kariya/internal/cli/screens/skills"
	"github.com/baphled/kariya/internal/cli/terminal"
	"github.com/baphled/kariya/internal/cli/themes"
	domain "github.com/baphled/kariya/internal/domain/career"
	career "github.com/baphled/kariya/internal/repository/career"
	"github.com/charmbracelet/bubbles/table"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
	overlay "github.com/rmhubbert/bubbletea-overlay"
)

// ManageSkillsIntent implements the Intent interface for managing user-defined skills.
type ManageSkillsIntent struct {
	*BaseIntent

	// context is the input context passed to the intent
	context *ManageSkillsContext

	// state tracks the current state
	currentState SkillsState

	// skills data
	skills        []*domain.Skill
	selectedIndex int
	selectedSkill *domain.Skill // Selected skill for detail view

	// table components for skills list view
	table         *table.Model
	listContainer *components.TableListContainer
	navHandler    *navigation.ListNavigationHandler

	// detail view data
	eventCounts  map[string]int        // Skill ID -> event count
	lastUsedMap  map[string]time.Time  // Skill ID -> last used date
	skillEvents  []*domain.CareerEvent // Events for selected skill
	eventsLoaded bool                  // Whether events have been loaded

	// events table components for skill events view
	eventsTable           *table.Model
	eventsListContainer   *components.TableListContainer
	eventsNavHandler      *navigation.ListNavigationHandler
	eventsSelectedIndex   int                 // Selection index for events list
	selectedEventFromList *domain.CareerEvent // Selected event for detail view from events list

	// form for add/edit
	skillForm *models.SkillForm

	// filter and sort state
	filters             *SkillsFilters // Active filters
	filterMenuIndex     int            // Selected option in filter menu
	sortMenuIndex       int            // Selected option in sort menu
	availableCategories []string       // Categories extracted from skills for filter menu

	// modals (new architecture with bubbletea-overlay)
	filterModal *components.SkillFilterModal
	sortModal   *components.SkillSortModal

	// screen orchestration (new architecture)
	activeScreen screens.Screen // Currently active screen (when using screen architecture)
	useScreens   bool           // Whether to use screen-based architecture (opt-in, default: false)

	// active indicates whether this intent is currently active
	active bool

	// result is the final result of the intent
	result *IntentResult[*ManageSkillsResult]
}

// SkillsFilters holds the active filter and sort state
type SkillsFilters struct {
	Category  string
	Level     string
	MinEvents int
	SortBy    string
	SortOrder string
}

// NewManageSkillsIntent creates a new ManageSkills intent
func NewManageSkillsIntent(ctx *ManageSkillsContext) *ManageSkillsIntent {
	baseIntent := NewBaseIntent()
	baseIntent.SetThemeManager(themes.NewThemeManager())

	// Create table model for skills list
	skillsColumns := []table.Column{
		{Title: "Name", Width: 25},
		{Title: "Category", Width: 15},
		{Title: "Level", Width: 12},
		{Title: "Years", Width: 8},
		{Title: "Events", Width: 8},
	}

	skillsTable := table.New(
		table.WithColumns(skillsColumns),
		table.WithRows([]table.Row{}),
		table.WithFocused(true),
		table.WithHeight(15),
		table.WithWidth(100),
	)

	// Apply default styles initially - theme styles will be applied in Init()
	skillsTable.SetStyles(table.DefaultStyles())

	// Create table model for skill events list
	eventsColumns := []table.Column{
		{Title: "Date", Width: 12},
		{Title: "Event", Width: 50},
		{Title: "Company", Width: 20},
	}

	eventsTable := table.New(
		table.WithColumns(eventsColumns),
		table.WithRows([]table.Row{}),
		table.WithFocused(true),
		table.WithHeight(15),
		table.WithWidth(100),
	)
	eventsTable.SetStyles(table.DefaultStyles())

	intent := &ManageSkillsIntent{
		BaseIntent:          baseIntent,
		context:             ctx,
		currentState:        SkillsStateList,
		skills:              []*domain.Skill{},
		selectedIndex:       0,
		filters:             &SkillsFilters{},
		active:              true,
		table:               &skillsTable,
		listContainer:       components.NewTableListContainer(skillsTable, "Manage Skills", 100),
		eventsTable:         &eventsTable,
		eventsListContainer: components.NewTableListContainer(eventsTable, "Skill Events", 100),
	}

	// Initialize navigation handler for skills list
	intent.navHandler = navigation.NewListNavigationHandler(intent)

	// Initialize navigation handler for events list using wrapper
	intent.eventsNavHandler = navigation.NewListNavigationHandler(&skillEventsNavigator{intent: intent})

	return intent
}

// skillEventsNavigator wraps ManageSkillsIntent to implement ListNavigator for the events list
type skillEventsNavigator struct {
	intent *ManageSkillsIntent
}

func (n *skillEventsNavigator) GetTotalItems() int {
	return len(n.intent.skillEvents)
}

func (n *skillEventsNavigator) GetSelectedIndex() int {
	return n.intent.eventsSelectedIndex
}

func (n *skillEventsNavigator) SetSelectedIndex(idx int) {
	// Validate and set index
	if idx < 0 {
		idx = 0
	}
	if idx >= len(n.intent.skillEvents) {
		idx = len(n.intent.skillEvents) - 1
	}
	if idx < 0 {
		idx = 0 // Handle empty list
	}

	n.intent.eventsSelectedIndex = idx

	// Update table display
	n.intent.updateEventsTableRows()
}

func (n *skillEventsNavigator) GetPageSize() int {
	return 15
}

// Init initializes the intent and loads skills
func (i *ManageSkillsIntent) Init() tea.Cmd {
	i.active = true

	// Disable screen architecture by default (tests expect legacy mode)
	// TODO: Fix screen orchestration bugs before re-enabling
	i.useScreens = false

	// Apply themed table styles if theme is available (for legacy fallback states)
	if theme := i.Theme(); theme != nil {
		i.table.SetStyles(themes.NewThemedTableStyles(theme))
	}

	// Load skills asynchronously
	return func() tea.Msg {
		skills, err := i.context.SkillRepository.List(i.context.Ctx, nil)
		return SkillsLoadedMsg{
			Skills: skills,
			Error:  err,
		}
	}
}

// updateTableRows updates the table rows based on skills
func (i *ManageSkillsIntent) updateTableRows() {
	pageSize := 15
	total := len(i.skills)

	// Determine which page current selection is on
	page := 0
	if pageSize > 0 && i.selectedIndex >= 0 {
		page = i.selectedIndex / pageSize
	}

	start := page * pageSize
	end := start + pageSize
	if end > total {
		end = total
	}

	// Handle empty list
	if total == 0 {
		i.table.SetRows([]table.Row{})
		i.listContainer.SetTable(*i.table)
		return
	}

	pageSkills := i.skills[start:end]

	rows := make([]table.Row, 0, len(pageSkills))
	for idx, skill := range pageSkills {
		realIdx := start + idx

		// Use centralized indicator formatting
		name := i.navHandler.FormatRowText(realIdx, skill.Name)

		// Category
		category := skill.Category
		if category == "" {
			category = "-"
		}

		// Level
		level := skill.Level
		if level == "" {
			level = "-"
		}

		// Years
		years := "-"
		if skill.YearsUsed != nil {
			years = fmt.Sprintf("%d", *skill.YearsUsed)
		}

		// Event count
		eventCount := "-"
		if i.eventCounts != nil {
			if count, ok := i.eventCounts[skill.ID]; ok {
				eventCount = fmt.Sprintf("%d", count)
			}
		}

		rows = append(rows, table.Row{name, category, level, years, eventCount})
	}

	i.table.SetRows(rows)

	// Calculate relative cursor position for this page
	relativeCursor := 0
	if i.selectedIndex >= start && i.selectedIndex < end {
		relativeCursor = i.selectedIndex - start
	}

	// Set table cursor to relative position within the page
	i.table.SetCursor(relativeCursor)

	// Sync the container's selectedIdx to match our relative cursor
	i.listContainer.SetSelectedIdx(relativeCursor)

	// Update the container with the modified table
	i.listContainer.SetTable(*i.table)
}

// updateEventsTableRows updates the events table rows based on skillEvents
func (i *ManageSkillsIntent) updateEventsTableRows() {
	pageSize := 15
	total := len(i.skillEvents)

	// Determine which page current selection is on
	page := 0
	if pageSize > 0 && i.eventsSelectedIndex >= 0 {
		page = i.eventsSelectedIndex / pageSize
	}

	start := page * pageSize
	end := start + pageSize
	if end > total {
		end = total
	}

	// Handle empty list
	if total == 0 {
		i.eventsTable.SetRows([]table.Row{})
		i.eventsListContainer.SetTable(*i.eventsTable)
		return
	}

	pageEvents := i.skillEvents[start:end]

	rows := make([]table.Row, 0, len(pageEvents))
	for idx, event := range pageEvents {
		realIdx := start + idx

		// Use centralized indicator formatting
		dateStr := i.eventsNavHandler.FormatRowText(realIdx, event.Date.Format("2006-01-02"))

		// Truncate text to first 50 chars
		text := event.Text
		if len(text) > 50 {
			text = text[:50] + "..."
		}

		// Company
		company := event.Company
		if company == "" {
			company = "-"
		}

		rows = append(rows, table.Row{dateStr, text, company})
	}

	i.eventsTable.SetRows(rows)

	// Calculate relative cursor position for this page
	relativeCursor := 0
	if i.eventsSelectedIndex >= start && i.eventsSelectedIndex < end {
		relativeCursor = i.eventsSelectedIndex - start
	}

	// Set table cursor to relative position within the page
	i.eventsTable.SetCursor(relativeCursor)

	// Sync the container's selectedIdx to match our relative cursor
	i.eventsListContainer.SetSelectedIdx(relativeCursor)

	// Update the container with the modified table
	i.eventsListContainer.SetTable(*i.eventsTable)
}

// Update handles messages and state transitions
func (i *ManageSkillsIntent) Update(msg tea.Msg) tea.Cmd {
	if !i.active {
		return nil
	}

	// PATTERN 4: Global Key Interception - 3-tier priority
	// 1. Check global keys FIRST (work everywhere, even in modals)
	if keyMsg, ok := msg.(tea.KeyMsg); ok {
		switch HandleGlobalKeys(keyMsg) {
		case KeyQuit:
			return tea.Quit
		case KeyHelp:
			i.helpModal.Toggle()
			return nil
		}

		// 2. SECOND PRIORITY: Modal updates (if visible)
		if i.filterModal != nil && i.filterModal.IsVisible() {
			return i.handleFilterModalUpdate(keyMsg)
		}
		if i.sortModal != nil && i.sortModal.IsVisible() {
			return i.handleSortModalUpdate(keyMsg)
		}
	}

	// Screen orchestration: delegate to active screen if present
	if i.useScreens && i.activeScreen != nil {
		// Handle global keys FIRST, even in screen mode
		if keyMsg, ok := msg.(tea.KeyMsg); ok {
			switch HandleGlobalKeys(keyMsg) {
			case KeyQuit:
				return tea.Quit
			case KeyHelp:
				i.helpModal.Toggle()
				return nil
			}
		}

		// Handle window size messages for screen
		if wsMsg, ok := msg.(tea.WindowSizeMsg); ok {
			i.activeScreen.SetTerminalInfo(wsMsg.Width, wsMsg.Height)
		}

		cmd, result := i.activeScreen.Update(msg)
		if result != nil {
			return i.handleScreenResult(result)
		}
		return cmd
	}

	switch msg := msg.(type) {
	case SkillsLoadedMsg:
		return i.handleSkillsLoaded(msg)

	case models.SkillFormCompleteMsg:
		return i.handleFormComplete(msg)

	case SkillFormCompleteMsg:
		// Handle the legacy message type for backward compatibility with tests
		return i.handleLegacyFormComplete(msg)

	case SkillCreatedMsg:
		return i.handleSkillCreated(msg)

	case SkillUpdatedMsg:
		return i.handleSkillUpdated(msg)

	case SkillDeletedMsg:
		return i.handleSkillDeleted(msg)

	case SkillEventsLoadedMsg:
		return i.handleSkillEventsLoaded(msg)

	case tea.KeyMsg:
		// Global keys already handled above at top of function

		// If we have a form active, forward key messages to it
		if i.skillForm != nil {
			// Check for Esc key to cancel form
			if msg.Type == tea.KeyEsc {
				return i.handleFormCancel()
			}

			// Forward to skillForm
			_, cmd := i.skillForm.Update(msg)
			return cmd
		}
		return i.handleKeyPress(msg)

	case tea.WindowSizeMsg:
		// Update terminal info
		termInfo := terminal.NewInfo()
		termInfo.Width = msg.Width
		termInfo.Height = msg.Height
		termInfo.IsValid = true
		i.UpdateTerminalInfo(termInfo)

		// Forward to skillForm if active - it handles its own dimensions
		if i.skillForm != nil {
			_, cmd := i.skillForm.Update(msg)
			return cmd
		}
		return nil
	}

	// If we have a form active, let it handle other messages
	if i.skillForm != nil {
		_, cmd := i.skillForm.Update(msg)
		return cmd
	}

	return nil
}

// View renders the current state
func (i *ManageSkillsIntent) View() string {
	if !i.active {
		return ""
	}

	// Screen orchestration: delegate to active screen if present
	if i.useScreens && i.activeScreen != nil {
		return i.activeScreen.View()
	}

	// Create standard view with breadcrumbs
	view := CreateStandardViewWithBreadcrumbs(i.BaseIntent, i.getBreadcrumbs()...)

	// Reduce logo spacing on small terminals to maximize form visibility
	if info := i.GetTerminalInfo(); info != nil && info.Height < 30 {
		if logo := i.GetLogo(); logo != nil {
			view.WithLogo(logo, 0) // No spacing above logo for small terminals
		}
	}

	// Get content for current state
	content := i.getStateContent()
	view.WithContent(content)

	// Add context-aware help
	help := i.getContextHelp()
	view.WithHelp(help).WithFooterSeparator(true)

	// PATTERN 1: Modal Overlay Rendering
	// StandardView FIRST, modal overlay LAST (prevents misalignment)
	baseView := view.Render()

	// Overlay modals as final step
	if i.filterModal != nil && i.filterModal.IsVisible() {
		return i.renderFilterModalOverlay(baseView)
	}
	if i.sortModal != nil && i.sortModal.IsVisible() {
		return i.renderSortModalOverlay(baseView)
	}

	return baseView
}

// renderFilterModalOverlay renders the filter modal over the base view
func (i *ManageSkillsIntent) renderFilterModalOverlay(baseView string) string {
	bgModel := &staticViewModel{content: baseView}
	overlayModel := overlay.New(
		i.filterModal,  // Foreground: the filter modal
		bgModel,        // Background: the rendered view
		overlay.Center, // X position
		overlay.Center, // Y position
		0,              // X offset
		-2,             // Y offset (move up 2 lines to avoid footer)
	)
	return overlayModel.View()
}

// renderSortModalOverlay renders the sort modal over the base view
func (i *ManageSkillsIntent) renderSortModalOverlay(baseView string) string {
	bgModel := &staticViewModel{content: baseView}
	overlayModel := overlay.New(
		i.sortModal,    // Foreground: the sort modal
		bgModel,        // Background: the rendered view
		overlay.Center, // X position
		overlay.Center, // Y position
		0,              // X offset
		-2,             // Y offset (move up 2 lines to avoid footer)
	)
	return overlayModel.View()
}

// handleFilterModalUpdate handles updates when filter modal is visible
func (i *ManageSkillsIntent) handleFilterModalUpdate(msg tea.KeyMsg) tea.Cmd {
	cmd, applied, filterData := i.filterModal.Update(msg)

	if applied && filterData != nil {
		// User confirmed filters - convert to SkillFilters and apply
		newFilters := i.filterModal.ToSkillFilters()

		// Update internal filter state
		if i.filters == nil {
			i.filters = &SkillsFilters{}
		}
		// Map SkillFilters to internal SkillsFilters format
		if len(newFilters.Categories) > 0 {
			i.filters.Category = newFilters.Categories[0] // Use first category for now
		} else {
			i.filters.Category = ""
		}
		if len(newFilters.Levels) > 0 {
			i.filters.Level = newFilters.Levels[0] // Use first level for now
		} else {
			i.filters.Level = ""
		}
		i.filters.MinEvents = newFilters.MinYears // Map years to events for now
		i.filters.SortBy = newFilters.SortBy
		i.filters.SortOrder = newFilters.SortOrder

		// Reload skills with new filters
		return i.reloadSkills()
	}

	// Modal was closed without completion (Esc) or still being edited
	return cmd
}

// handleSortModalUpdate handles updates when sort modal is visible
func (i *ManageSkillsIntent) handleSortModalUpdate(msg tea.KeyMsg) tea.Cmd {
	cmd, applied, sortData := i.sortModal.Update(msg)

	if applied && sortData != nil {
		// User confirmed sort - apply it
		sortConfig := i.sortModal.ToSkillSortConfig()

		// Update internal filter state
		if i.filters == nil {
			i.filters = &SkillsFilters{}
		}
		i.filters.SortBy = sortConfig.SortBy
		i.filters.SortOrder = sortConfig.SortOrder

		// Reload skills with new sort
		return i.reloadSkills()
	}

	// Modal was closed without completion (Esc) or still being edited
	return cmd
}

// openFilterModal opens the filter modal with current filters pre-populated
// PATTERN 12: Form Modal with Immediate Init
func (i *ManageSkillsIntent) openFilterModal() tea.Cmd {
	// Get terminal dimensions
	termInfo := i.GetTerminalInfo()
	width := 120
	height := 40
	if termInfo != nil {
		width = termInfo.Width
		height = termInfo.Height
	}

	// Build current filters for pre-population
	var currentFilters *components.SkillFilters
	if i.filters != nil {
		currentFilters = &components.SkillFilters{
			Categories: []string{},
			Levels:     []string{},
			MinYears:   i.filters.MinEvents, // Map events to years for now
			MaxYears:   0,
			SearchText: "",
			SortBy:     i.filters.SortBy,
			SortOrder:  i.filters.SortOrder,
		}
		if i.filters.Category != "" {
			currentFilters.Categories = []string{i.filters.Category}
		}
		if i.filters.Level != "" {
			currentFilters.Levels = []string{i.filters.Level}
		}
	}

	// Create filter modal
	i.filterModal = components.NewSkillFilterModal(
		i.skills,
		currentFilters,
		width,
		height,
	)

	// CRITICAL: Call Init() for immediate rendering
	return i.filterModal.Init()
}

// openSortModal opens the sort modal with current sort config pre-populated
// PATTERN 12: Form Modal with Immediate Init
func (i *ManageSkillsIntent) openSortModal() tea.Cmd {
	// Get terminal dimensions
	termInfo := i.GetTerminalInfo()
	width := 120
	height := 40
	if termInfo != nil {
		width = termInfo.Width
		height = termInfo.Height
	}

	// Build current sort config for pre-population
	var currentSort *components.SkillSortConfig
	if i.filters != nil {
		currentSort = &components.SkillSortConfig{
			SortBy:    i.filters.SortBy,
			SortOrder: i.filters.SortOrder,
		}
	}

	// Create sort modal
	i.sortModal = components.NewSkillSortModal(
		i.skills,
		currentSort,
		width,
		height,
	)

	// CRITICAL: Call Init() for immediate rendering
	return i.sortModal.Init()
}

// getStateContent returns the content for the current state
func (i *ManageSkillsIntent) getStateContent() string {
	switch i.currentState {
	case SkillsStateList:
		return i.renderSkillsList()
	case SkillsStateDetail:
		return i.renderSkillDetail()
	case SkillsStateDetailEvents:
		return i.renderSkillEvents()
	case SkillsStateDetailEventDetail:
		return i.renderEventDetail()
	case SkillsStateAdd, SkillsStateEdit:
		return i.renderForm()
	case SkillsStateDelete:
		return i.renderDeleteConfirm()
	case SkillsStateFilter:
		return i.renderFilterMenu()
	case SkillsStateSort:
		return i.renderSortMenu()
	default:
		return "Unknown state"
	}
}

// getBreadcrumbs returns breadcrumbs for the current state
func (i *ManageSkillsIntent) getBreadcrumbs() []string {
	breadcrumbs := []string{"Skills"}

	switch i.currentState {
	case SkillsStateDetail, SkillsStateDetailEvents, SkillsStateDetailEventDetail:
		if i.selectedSkill != nil {
			breadcrumbs = append(breadcrumbs, i.selectedSkill.Name)
		}
		if i.currentState == SkillsStateDetailEvents {
			breadcrumbs = append(breadcrumbs, "Events")
		}
		if i.currentState == SkillsStateDetailEventDetail {
			breadcrumbs = append(breadcrumbs, "Events", "Detail")
		}
	case SkillsStateAdd:
		breadcrumbs = append(breadcrumbs, "Add")
	case SkillsStateEdit:
		breadcrumbs = append(breadcrumbs, "Edit")
	case SkillsStateDelete:
		breadcrumbs = append(breadcrumbs, "Delete")
	}

	return breadcrumbs
}

// getContextHelp returns help text for the current state
func (i *ManageSkillsIntent) getContextHelp() string {
	theme := i.Theme()

	switch i.currentState {
	case SkillsStateList:
		badges := []components.KeyBadge{
			components.NewKeyBadge("Enter", "View details"),
			components.NewKeyBadge("n", "New skill"),
			components.NewKeyBadge("f", "Filter"),
			components.NewKeyBadge("s", "Sort"),
		}
		if i.hasActiveFilters() {
			badges = append(badges, components.NewKeyBadge("x", "Clear filters"))
		}
		return CombineThemedFooters(
			ThemedListFooter(theme),
			ThemedCustomFooter(theme, badges...),
			ThemedGlobalBadges(theme),
		)
	case SkillsStateDetail:
		return CombineThemedFooters(
			ThemedDetailViewFooter(theme),
			ThemedCustomFooter(theme,
				components.NewKeyBadge("Enter", "View events"),
				components.EditBadge(),
				components.DeleteBadge(),
			),
			ThemedGlobalBadges(theme),
		)
	case SkillsStateDetailEvents:
		return CombineThemedFooters(
			ThemedListFooter(theme),
			ThemedCustomFooter(theme,
				components.NewKeyBadge("Enter", "View details"),
				components.EditBadge(),
			),
			ThemedGlobalBadges(theme),
		)
	case SkillsStateDetailEventDetail:
		return CombineThemedFooters(
			ThemedDetailViewFooter(theme),
			ThemedCustomFooter(theme,
				components.EditBadge(),
			),
			ThemedGlobalBadges(theme),
		)
	case SkillsStateAdd, SkillsStateEdit:
		return CombineThemedFooters(
			ThemedFormFooter(theme),
			ThemedGlobalBadges(theme),
		)
	case SkillsStateDelete:
		return CombineThemedFooters(
			ThemedCustomFooter(theme,
				components.NewKeyBadge("y", "Confirm"),
				components.NewKeyBadge("n/Esc", "Cancel"),
			),
			ThemedGlobalBadges(theme),
		)
	case SkillsStateFilter:
		return CombineThemedFooters(
			ThemedListFooter(theme),
			ThemedCustomFooter(theme,
				components.NewKeyBadge("Enter", "Apply"),
				components.NewKeyBadge("u", "Used skills only"),
			),
			ThemedGlobalBadges(theme),
		)
	case SkillsStateSort:
		return CombineThemedFooters(
			ThemedListFooter(theme),
			ThemedCustomFooter(theme,
				components.NewKeyBadge("Enter", "Apply"),
				components.NewKeyBadge("e", "Most used"),
			),
			ThemedGlobalBadges(theme),
		)
	default:
		return ThemedGlobalBadges(theme)
	}
}

// Result returns the final result of the intent
func (i *ManageSkillsIntent) Result() *IntentResult[interface{}] {
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

// State returns the current state (for testing)
func (i *ManageSkillsIntent) State() SkillsState {
	return i.currentState
}

// SelectedIndex returns the currently selected index (for testing)
func (i *ManageSkillsIntent) SelectedIndex() int {
	return i.selectedIndex
}

// ActiveFilters returns the current filter/sort state (for testing)
func (i *ManageSkillsIntent) ActiveFilters() *SkillsFilters {
	if i.filters == nil {
		i.filters = &SkillsFilters{}
	}
	return i.filters
}

// Message handlers

func (i *ManageSkillsIntent) handleSkillsLoaded(msg SkillsLoadedMsg) tea.Cmd {
	if msg.Error != nil {
		i.result = &IntentResult[*ManageSkillsResult]{
			Status: Failed,
			Error:  &IntentError{Code: "LOAD_FAILED", Message: "Failed to load skills", Cause: msg.Error},
		}
		i.active = false
		return nil
	}

	i.skills = msg.Skills
	i.selectedIndex = 0

	// Load event counts for displaying in list view
	eventCounts, err := i.context.SkillRepository.GetEventCountsForSkills(i.context.Ctx)
	if err == nil {
		i.eventCounts = eventCounts
	} else {
		// If loading fails, use empty map
		i.eventCounts = make(map[string]int)
	}

	// If using screens, transition to list screen
	if i.useScreens {
		return i.transitionToListScreen()
	}

	// Legacy: Update table rows with new skills data
	i.updateTableRows()

	return nil
}

func (i *ManageSkillsIntent) handleFormComplete(msg models.SkillFormCompleteMsg) tea.Cmd {
	if msg.Cancelled {
		i.currentState = SkillsStateList
		i.skillForm = nil
		return nil
	}

	// Extract form data and create/update skill
	formData := msg.Data
	if formData == nil {
		return nil
	}

	// Check if user cancelled via the submit button
	if !formData.SubmitConfirmed {
		i.currentState = SkillsStateList
		i.skillForm = nil
		return nil
	}

	// Create skill from form data
	var skill *domain.Skill
	if i.currentState == SkillsStateEdit && len(i.skills) > 0 {
		// Editing existing skill - preserve ID
		skill = i.skills[i.selectedIndex]
	} else {
		// Creating new skill
		skill = &domain.Skill{}
	}

	// Apply form data to skill
	forms.ApplySkillFormData(skill, formData)

	// Save skill based on state
	if i.currentState == SkillsStateAdd {
		return i.createSkill(skill)
	} else if i.currentState == SkillsStateEdit {
		return i.updateSkill(skill)
	}

	return nil
}

// handleLegacyFormComplete handles the legacy SkillFormCompleteMsg (used by tests)
func (i *ManageSkillsIntent) handleLegacyFormComplete(msg SkillFormCompleteMsg) tea.Cmd {
	if msg.Cancelled {
		i.currentState = SkillsStateList
		i.skillForm = nil
		return nil
	}

	if msg.Error != nil {
		// Stay in form state with error
		return nil
	}

	// Save skill based on state
	if i.currentState == SkillsStateAdd {
		return i.createSkill(msg.Skill)
	} else if i.currentState == SkillsStateEdit {
		return i.updateSkill(msg.Skill)
	}

	return nil
}

func (i *ManageSkillsIntent) handleSkillCreated(msg SkillCreatedMsg) tea.Cmd {
	if msg.Error != nil {
		i.result = &IntentResult[*ManageSkillsResult]{
			Status: Failed,
			Error:  &IntentError{Code: "CREATE_FAILED", Message: "Failed to create skill", Cause: msg.Error},
		}
		return nil
	}

	// Transition back to list and reload
	i.currentState = SkillsStateList
	i.skillForm = nil

	return i.Init()
}

func (i *ManageSkillsIntent) handleSkillUpdated(msg SkillUpdatedMsg) tea.Cmd {
	if msg.Error != nil {
		i.result = &IntentResult[*ManageSkillsResult]{
			Status: Failed,
			Error:  &IntentError{Code: "UPDATE_FAILED", Message: "Failed to update skill", Cause: msg.Error},
		}
		return nil
	}

	// Transition back to list and reload
	i.currentState = SkillsStateList
	i.skillForm = nil

	return i.Init()
}

func (i *ManageSkillsIntent) handleSkillDeleted(msg SkillDeletedMsg) tea.Cmd {
	if msg.Error != nil {
		i.result = &IntentResult[*ManageSkillsResult]{
			Status: Failed,
			Error:  &IntentError{Code: "DELETE_FAILED", Message: "Failed to delete skill", Cause: msg.Error},
		}
		return nil
	}

	// Transition back to list and reload
	i.currentState = SkillsStateList

	return i.Init()
}

func (i *ManageSkillsIntent) handleKeyPress(msg tea.KeyMsg) tea.Cmd {
	switch i.currentState {
	case SkillsStateList:
		return i.handleListKeys(msg)
	case SkillsStateDetail:
		return i.handleDetailKeys(msg)
	case SkillsStateDetailEvents:
		return i.handleDetailEventsKeys(msg)
	case SkillsStateDetailEventDetail:
		return i.handleEventDetailKeys(msg)
	case SkillsStateDelete:
		return i.handleDeleteKeys(msg)
	case SkillsStateFilter:
		return i.handleFilterKeys(msg)
	case SkillsStateSort:
		return i.handleSortKeys(msg)
	default:
		return nil
	}
}

func (i *ManageSkillsIntent) handleListKeys(msg tea.KeyMsg) tea.Cmd {
	// Handle global keys first using MessageInterceptor
	return NewMessageInterceptor().
		OnQuit(StandardQuitHandler()).
		OnHelp(StandardHelpHandler(i.BaseIntent)).
		OnBack(func() tea.Cmd {
			// At root state, back means cancel and return to main menu
			i.setCancelled()
			return nil
		}).
		InterceptOr(msg, func() tea.Cmd {
			// Try list navigation handler first (handles j/k, up/down, pgup/pgdn, home/end, g/G)
			if i.navHandler.HandleKey(msg.String()) {
				return nil
			}

			switch msg.String() {
			case "enter":
				// View skill detail
				if len(i.skills) == 0 {
					return nil
				}
				i.selectedSkill = i.skills[i.selectedIndex]
				i.currentState = SkillsStateDetail
				return i.loadDetailData()

			case "n":
				// Add new skill
				i.currentState = SkillsStateAdd
				i.skillForm = models.NewSkillForm()
				return i.skillForm.Init()

			case "e":
				// Edit selected skill
				if len(i.skills) == 0 {
					return nil
				}
				i.currentState = SkillsStateEdit
				i.skillForm = models.NewSkillFormWithData(i.skills[i.selectedIndex])
				return i.skillForm.Init()

			case "d":
				// Delete selected skill
				if len(i.skills) == 0 {
					return nil
				}
				i.currentState = SkillsStateDelete
				return nil

			case "f":
				// PATTERN 12: Form Modal with Immediate Init
				// Open filter modal instead of menu state
				return i.openFilterModal()

			case "s":
				// PATTERN 12: Form Modal with Immediate Init
				// Open sort modal instead of menu state
				return i.openSortModal()

			case "x":
				// Clear all filters
				if i.hasActiveFilters() {
					i.filters = &SkillsFilters{}
					return i.reloadSkills()
				}
				return nil
			}

			return nil
		})
}

// setCancelled marks the intent as cancelled and returns to main menu
func (i *ManageSkillsIntent) setCancelled() {
	i.result = &IntentResult[*ManageSkillsResult]{
		Status: Cancelled,
		Data: &ManageSkillsResult{
			Action: "cancelled",
		},
	}
	i.active = false
}

// handleFilterKeys handles key presses in filter menu
func (i *ManageSkillsIntent) handleFilterKeys(msg tea.KeyMsg) tea.Cmd {
	// Handle global keys first using MessageInterceptor
	return NewMessageInterceptor().
		OnQuit(StandardQuitHandler()).
		OnHelp(StandardHelpHandler(i.BaseIntent)).
		OnBack(func() tea.Cmd {
			i.currentState = SkillsStateList
			return nil
		}).
		InterceptOr(msg, func() tea.Cmd {
			switch msg.String() {
			case "j", "down":
				i.filterMenuIndex++
				// Wrap around: 0=All Categories, 1..n=categories, n+1=All Levels, n+2..m=levels, m+1=Used skills only
				maxIndex := len(i.availableCategories) + 6 // Categories + "All" + Levels + "All" + "Used only"
				if i.filterMenuIndex >= maxIndex {
					i.filterMenuIndex = 0
				}
				return nil

			case "k", "up":
				i.filterMenuIndex--
				maxIndex := len(i.availableCategories) + 6
				if i.filterMenuIndex < 0 {
					i.filterMenuIndex = maxIndex - 1
				}
				return nil

			case "enter":
				// Apply selected filter
				return i.applyFilterSelection()

			case "u":
				// Quick shortcut for "Used skills only"
				i.filters.MinEvents = 1
				i.currentState = SkillsStateList
				return i.reloadSkills()
			}

			return nil
		})
}

// handleSortKeys handles key presses in sort menu
func (i *ManageSkillsIntent) handleSortKeys(msg tea.KeyMsg) tea.Cmd {
	// Handle global keys first using MessageInterceptor
	return NewMessageInterceptor().
		OnQuit(StandardQuitHandler()).
		OnHelp(StandardHelpHandler(i.BaseIntent)).
		OnBack(func() tea.Cmd {
			i.currentState = SkillsStateList
			return nil
		}).
		InterceptOr(msg, func() tea.Cmd {
			switch msg.String() {
			case "j", "down":
				i.sortMenuIndex++
				if i.sortMenuIndex > 5 { // 6 sort options
					i.sortMenuIndex = 0
				}
				return nil

			case "k", "up":
				i.sortMenuIndex--
				if i.sortMenuIndex < 0 {
					i.sortMenuIndex = 5
				}
				return nil

			case "enter":
				return i.applySortSelection()

			case "e":
				// Quick shortcut for "Most used" (events desc)
				i.filters.SortBy = "events"
				i.filters.SortOrder = "desc"
				i.currentState = SkillsStateList
				return i.reloadSkills()
			}

			return nil
		})
}

// extractAvailableCategories extracts unique categories from loaded skills
func (i *ManageSkillsIntent) extractAvailableCategories() {
	categorySet := make(map[string]bool)
	for _, skill := range i.skills {
		if skill.Category != "" {
			categorySet[skill.Category] = true
		}
	}

	i.availableCategories = make([]string, 0, len(categorySet))
	for cat := range categorySet {
		i.availableCategories = append(i.availableCategories, cat)
	}
	sort.Strings(i.availableCategories)
}

// applyFilterSelection applies the currently selected filter option
func (i *ManageSkillsIntent) applyFilterSelection() tea.Cmd {
	// Menu structure:
	// 0: All Categories
	// 1..n: Specific categories
	// n+1: All Levels
	// n+2..n+5: Specific levels (beginner, intermediate, advanced, expert)
	// n+6: Used skills only

	numCategories := len(i.availableCategories)
	levels := []string{"beginner", "intermediate", "advanced", "expert"}

	if i.filterMenuIndex == 0 {
		// All Categories
		i.filters.Category = ""
	} else if i.filterMenuIndex <= numCategories {
		// Specific category
		i.filters.Category = i.availableCategories[i.filterMenuIndex-1]
	} else if i.filterMenuIndex == numCategories+1 {
		// All Levels
		i.filters.Level = ""
	} else if i.filterMenuIndex <= numCategories+1+len(levels) {
		// Specific level
		levelIndex := i.filterMenuIndex - numCategories - 2
		if levelIndex >= 0 && levelIndex < len(levels) {
			i.filters.Level = levels[levelIndex]
		}
	} else {
		// Used skills only
		i.filters.MinEvents = 1
	}

	i.currentState = SkillsStateList
	return i.reloadSkills()
}

// applySortSelection applies the currently selected sort option
func (i *ManageSkillsIntent) applySortSelection() tea.Cmd {
	// Sort options:
	// 0: Name A-Z
	// 1: Name Z-A
	// 2: Most used (events desc)
	// 3: Least used (events asc)
	// 4: Category A-Z
	// 5: Category Z-A

	switch i.sortMenuIndex {
	case 0:
		i.filters.SortBy = "name"
		i.filters.SortOrder = "asc"
	case 1:
		i.filters.SortBy = "name"
		i.filters.SortOrder = "desc"
	case 2:
		i.filters.SortBy = "events"
		i.filters.SortOrder = "desc"
	case 3:
		i.filters.SortBy = "events"
		i.filters.SortOrder = "asc"
	case 4:
		i.filters.SortBy = "category"
		i.filters.SortOrder = "asc"
	case 5:
		i.filters.SortBy = "category"
		i.filters.SortOrder = "desc"
	}

	i.currentState = SkillsStateList
	return i.reloadSkills()
}

// hasActiveFilters returns true if any filters are active
func (i *ManageSkillsIntent) hasActiveFilters() bool {
	return i.filters != nil && (i.filters.Category != "" || i.filters.Level != "" || i.filters.MinEvents > 0 || i.filters.SortBy != "")
}

// reloadSkills reloads skills with current filters
func (i *ManageSkillsIntent) reloadSkills() tea.Cmd {
	return func() tea.Msg {
		// Convert internal filters to repository filters
		var repoFilters *career.SkillFilters
		if i.filters != nil {
			repoFilters = &career.SkillFilters{
				Category:  i.filters.Category,
				Level:     i.filters.Level,
				MinEvents: i.filters.MinEvents,
				SortBy:    i.filters.SortBy,
				SortOrder: i.filters.SortOrder,
			}
		}

		skills, err := i.context.SkillRepository.List(i.context.Ctx, repoFilters)
		return SkillsLoadedMsg{
			Skills: skills,
			Error:  err,
		}
	}
}

func (i *ManageSkillsIntent) handleDeleteKeys(msg tea.KeyMsg) tea.Cmd {
	// Handle global keys first using MessageInterceptor
	return NewMessageInterceptor().
		OnQuit(StandardQuitHandler()).
		OnHelp(StandardHelpHandler(i.BaseIntent)).
		OnBack(func() tea.Cmd {
			// Cancel delete
			i.currentState = SkillsStateList
			return nil
		}).
		InterceptOr(msg, func() tea.Cmd {
			switch msg.String() {
			case "y":
				// Confirm delete
				if len(i.skills) == 0 {
					return nil
				}

				skillToDelete := i.skills[i.selectedIndex]
				return func() tea.Msg {
					err := i.context.SkillRepository.Delete(i.context.Ctx, skillToDelete.ID)
					return SkillDeletedMsg{
						SkillID: skillToDelete.ID,
						Error:   err,
					}
				}

			case "n":
				// Cancel delete
				i.currentState = SkillsStateList
				return nil
			}

			return nil
		})
}

func (i *ManageSkillsIntent) handleFormCancel() tea.Cmd {
	i.currentState = SkillsStateList
	i.skillForm = nil
	return nil
}

func (i *ManageSkillsIntent) createSkill(skill *domain.Skill) tea.Cmd {
	return func() tea.Msg {
		err := i.context.SkillRepository.Create(i.context.Ctx, skill)
		return SkillCreatedMsg{
			Skill: skill,
			Error: err,
		}
	}
}

func (i *ManageSkillsIntent) updateSkill(skill *domain.Skill) tea.Cmd {
	return func() tea.Msg {
		err := i.context.SkillRepository.Update(i.context.Ctx, skill)
		return SkillUpdatedMsg{
			Skill: skill,
			Error: err,
		}
	}
}

// View renderers

func (i *ManageSkillsIntent) renderForm() string {
	if i.skillForm == nil {
		return "Form not initialized"
	}

	return i.skillForm.View()
}

func (i *ManageSkillsIntent) renderDeleteConfirm() string {
	if len(i.skills) == 0 {
		return "No skill selected"
	}

	skill := i.skills[i.selectedIndex]

	modalContent := fmt.Sprintf("Are you sure you want to delete '%s'?\n\nThis action cannot be undone.", skill.Name)
	modal := components.NewWarningModal("Delete Skill", modalContent)

	// Get terminal dimensions
	width, height := 80, 24
	if termInfo := i.GetTerminalInfo(); termInfo != nil && termInfo.IsValid {
		width = termInfo.Width
		height = termInfo.Height
	}

	return modal.Render(width, height)
}

// getCardStyle returns a themed card style for consistent content presentation.
func (i *ManageSkillsIntent) getCardStyle() lipgloss.Style {
	if theme := i.Theme(); theme != nil {
		return theme.Styles().CardBase
	}
	// Fallback to default styling
	return lipgloss.NewStyle().
		Padding(1, 2).
		BorderStyle(lipgloss.RoundedBorder()).
		BorderForeground(lipgloss.Color("#585B70"))
}

func (i *ManageSkillsIntent) renderSkillsList() string {
	if len(i.skills) == 0 {
		i.listContainer.SetEmptyStateMessage("No skills defined yet.\n\nPress 'n' to add your first skill.")
		return i.listContainer.Render()
	}

	// Ensure table rows are synchronized with current state
	i.updateTableRows()

	// Build pagination info with page number indicator
	pageSize := 15
	totalItems := len(i.skills)
	currentPage := (i.selectedIndex / pageSize) + 1
	totalPages := (totalItems + pageSize - 1) / pageSize
	paginationInfo := fmt.Sprintf("Skills: %d | Page %d of %d", totalItems, currentPage, totalPages)
	i.listContainer.SetPaginationInfo(paginationInfo)

	return i.listContainer.Render()
}

// loadDetailData loads event counts and last used dates for detail view
func (i *ManageSkillsIntent) loadDetailData() tea.Cmd {
	// Load synchronously since we need this data immediately
	eventCounts, err := i.context.SkillRepository.GetEventCountsForSkills(i.context.Ctx)
	if err != nil {
		// If loading fails, use empty maps
		i.eventCounts = make(map[string]int)
	} else {
		i.eventCounts = eventCounts
	}

	lastUsedMap, err := i.context.SkillRepository.GetLastUsedForSkills(i.context.Ctx)
	if err != nil {
		// If loading fails, use empty map
		i.lastUsedMap = make(map[string]time.Time)
	} else {
		i.lastUsedMap = lastUsedMap
	}

	return nil
}

// handleSkillEventsLoaded handles the SkillEventsLoadedMsg
func (i *ManageSkillsIntent) handleSkillEventsLoaded(msg SkillEventsLoadedMsg) tea.Cmd {
	if msg.Error != nil {
		// Show error but stay in detail view
		return nil
	}

	i.skillEvents = msg.Events
	i.eventsLoaded = true
	i.eventsSelectedIndex = 0

	// Apply themed table styles for events table
	if theme := i.Theme(); theme != nil {
		i.eventsTable.SetStyles(themes.NewThemedTableStyles(theme))
	}

	// Update table rows with new events data
	i.updateEventsTableRows()

	return nil
}

// handleDetailKeys handles key presses in detail view
func (i *ManageSkillsIntent) handleDetailKeys(msg tea.KeyMsg) tea.Cmd {
	// Handle global keys first using MessageInterceptor
	return NewMessageInterceptor().
		OnQuit(StandardQuitHandler()).
		OnHelp(StandardHelpHandler(i.BaseIntent)).
		OnBack(func() tea.Cmd {
			// Go back to list
			i.currentState = SkillsStateList
			i.selectedSkill = nil
			return nil
		}).
		InterceptOr(msg, func() tea.Cmd {
			switch msg.String() {
			case "enter":
				// View events using this skill
				i.currentState = SkillsStateDetailEvents
				i.eventsLoaded = false
				return i.loadEventsForSkill()

			case "e":
				// Edit this skill
				i.currentState = SkillsStateEdit
				i.skillForm = models.NewSkillFormWithData(i.selectedSkill)
				return i.skillForm.Init()

			case "d":
				// Delete this skill
				i.currentState = SkillsStateDelete
				return nil
			}

			return nil
		})
}

// handleDetailEventsKeys handles key presses in detail events view
func (i *ManageSkillsIntent) handleDetailEventsKeys(msg tea.KeyMsg) tea.Cmd {
	// Handle global keys first using MessageInterceptor
	return NewMessageInterceptor().
		OnQuit(StandardQuitHandler()).
		OnHelp(StandardHelpHandler(i.BaseIntent)).
		OnBack(func() tea.Cmd {
			// Back to skill detail view
			i.currentState = SkillsStateDetail
			i.skillEvents = nil
			i.eventsLoaded = false
			i.eventsSelectedIndex = 0
			return nil
		}).
		InterceptOr(msg, func() tea.Cmd {
			// Try list navigation handler first (handles j/k, up/down, pgup/pgdn, home/end, g/G)
			if i.eventsNavHandler.HandleKey(msg.String()) {
				return nil
			}

			switch msg.String() {
			case "enter":
				// View event details
				if len(i.skillEvents) > 0 && i.eventsSelectedIndex >= 0 && i.eventsSelectedIndex < len(i.skillEvents) {
					i.selectedEventFromList = i.skillEvents[i.eventsSelectedIndex]
					i.currentState = SkillsStateDetailEventDetail
				}
				return nil

			case "e":
				// Edit selected event - send to app router to open CaptureEvent intent
				if len(i.skillEvents) > 0 && i.eventsSelectedIndex >= 0 && i.eventsSelectedIndex < len(i.skillEvents) {
					selectedEvent := i.skillEvents[i.eventsSelectedIndex]
					return func() tea.Msg {
						return RequestEditEventMsg{Event: selectedEvent}
					}
				}
				return nil
			}

			return nil
		})
}

// loadEventsForSkill loads events that use the selected skill
func (i *ManageSkillsIntent) loadEventsForSkill() tea.Cmd {
	return func() tea.Msg {
		events, err := i.context.SkillRepository.GetEventsUsingSkill(i.context.Ctx, i.selectedSkill.ID)
		return SkillEventsLoadedMsg{
			Events: events,
			Error:  err,
		}
	}
}

// renderSkillDetail renders the skill detail content
func (i *ManageSkillsIntent) renderSkillDetail() string {
	skill := i.selectedSkill
	theme := i.Theme()

	labelStyle := lipgloss.NewStyle().
		Foreground(theme.MutedColor()).
		Width(15)

	valueStyle := lipgloss.NewStyle().
		Foreground(theme.PrimaryColor()).
		Bold(true)

	var lines []string

	// Name
	lines = append(lines, labelStyle.Render("Name:")+valueStyle.Render(skill.Name))

	// Category
	lines = append(lines, labelStyle.Render("Category:")+valueStyle.Render(skill.Category))

	// Level (if set)
	if skill.Level != "" {
		lines = append(lines, labelStyle.Render("Level:")+valueStyle.Render(skill.Level))
	}

	// Years Used (if set)
	if skill.YearsUsed != nil {
		yearText := fmt.Sprintf("%d year", *skill.YearsUsed)
		if *skill.YearsUsed != 1 {
			yearText += "s"
		}
		lines = append(lines, labelStyle.Render("Years Used:")+valueStyle.Render(yearText))
	}

	// Event Count
	eventCount := 0
	if i.eventCounts != nil {
		eventCount = i.eventCounts[skill.ID]
	}
	lines = append(lines, labelStyle.Render("Event Count:")+valueStyle.Render(fmt.Sprintf("%d", eventCount)))

	// Last Used (if available)
	if i.lastUsedMap != nil {
		if lastUsed, ok := i.lastUsedMap[skill.ID]; ok {
			lines = append(lines, labelStyle.Render("Last Used:")+valueStyle.Render(lastUsed.Format("2006-01-02")))
		}
	}

	// Timestamps
	lines = append(lines, "")
	lines = append(lines, labelStyle.Render("Created:")+lipgloss.NewStyle().Foreground(theme.MutedColor()).Render(skill.CreatedAt.Format("2006-01-02 15:04")))
	lines = append(lines, labelStyle.Render("Updated:")+lipgloss.NewStyle().Foreground(theme.MutedColor()).Render(skill.UpdatedAt.Format("2006-01-02 15:04")))

	content := lipgloss.JoinVertical(lipgloss.Left, lines...)
	return i.getCardStyle().Render(content)
}

// renderSkillEvents renders the events using this skill as a table
func (i *ManageSkillsIntent) renderSkillEvents() string {
	theme := i.Theme()

	if !i.eventsLoaded {
		loadingStyle := lipgloss.NewStyle().
			Foreground(theme.SecondaryColor()).
			MarginTop(2)
		return i.getCardStyle().Render(loadingStyle.Render("Loading events..."))
	}

	if len(i.skillEvents) == 0 {
		i.eventsListContainer.SetEmptyStateMessage("No events use this skill yet.")
		return i.eventsListContainer.Render()
	}

	// Ensure table rows are synchronized with current state
	i.updateEventsTableRows()

	// Build pagination info with page number indicator
	pageSize := 15
	totalItems := len(i.skillEvents)
	currentPage := (i.eventsSelectedIndex / pageSize) + 1
	totalPages := (totalItems + pageSize - 1) / pageSize
	paginationInfo := fmt.Sprintf("Events: %d | Page %d of %d", totalItems, currentPage, totalPages)
	i.eventsListContainer.SetPaginationInfo(paginationInfo)

	return i.eventsListContainer.Render()
}

// renderEventDetail renders a single event's details using the reusable component
func (i *ManageSkillsIntent) renderEventDetail() string {
	return components.RenderEventDetailCard(i.selectedEventFromList, i.Theme())
}

// handleEventDetailKeys handles key presses in event detail view
func (i *ManageSkillsIntent) handleEventDetailKeys(msg tea.KeyMsg) tea.Cmd {
	// Handle global keys first using MessageInterceptor
	return NewMessageInterceptor().
		OnQuit(StandardQuitHandler()).
		OnHelp(StandardHelpHandler(i.BaseIntent)).
		OnBack(func() tea.Cmd {
			// Back to events list
			i.currentState = SkillsStateDetailEvents
			i.selectedEventFromList = nil
			return nil
		}).
		InterceptOr(msg, func() tea.Cmd {
			switch msg.String() {
			case "e":
				// Edit event - send to app router to open CaptureEvent intent
				if i.selectedEventFromList != nil {
					return func() tea.Msg {
						return RequestEditEventMsg{Event: i.selectedEventFromList}
					}
				}
				return nil
			}

			return nil
		})
}

// renderFilterMenu renders the filter menu
func (i *ManageSkillsIntent) renderFilterMenu() string {
	theme := i.Theme()

	titleStyle := lipgloss.NewStyle().
		Foreground(theme.PrimaryColor()).
		Bold(true).
		MarginBottom(1)

	selectedStyle := lipgloss.NewStyle().
		Foreground(theme.SuccessColor()).
		Bold(true)

	normalStyle := lipgloss.NewStyle().
		Foreground(theme.ForegroundColor())

	mutedStyle := lipgloss.NewStyle().
		Foreground(theme.MutedColor())

	var lines []string
	lines = append(lines, titleStyle.Render("Filter Skills"))
	lines = append(lines, "")

	// Category section
	lines = append(lines, mutedStyle.Render("Category:"))
	currentIndex := 0

	// All categories option
	marker := "  "
	if i.filterMenuIndex == currentIndex {
		marker = "▶ "
	}
	text := "All Categories"
	if i.filters.Category == "" {
		text += " ✓"
	}
	if i.filterMenuIndex == currentIndex {
		lines = append(lines, selectedStyle.Render(marker+text))
	} else {
		lines = append(lines, normalStyle.Render(marker+text))
	}
	currentIndex++

	// Specific categories
	for _, cat := range i.availableCategories {
		marker = "  "
		if i.filterMenuIndex == currentIndex {
			marker = "▶ "
		}
		text := cat
		if i.filters.Category == cat {
			text += " ✓"
		}
		if i.filterMenuIndex == currentIndex {
			lines = append(lines, selectedStyle.Render(marker+text))
		} else {
			lines = append(lines, normalStyle.Render(marker+text))
		}
		currentIndex++
	}

	lines = append(lines, "")
	lines = append(lines, mutedStyle.Render("Level:"))

	// All levels option
	marker = "  "
	if i.filterMenuIndex == currentIndex {
		marker = "▶ "
	}
	text = "All Levels"
	if i.filters.Level == "" {
		text += " ✓"
	}
	if i.filterMenuIndex == currentIndex {
		lines = append(lines, selectedStyle.Render(marker+text))
	} else {
		lines = append(lines, normalStyle.Render(marker+text))
	}
	currentIndex++

	// Specific levels
	levels := []string{"beginner", "intermediate", "advanced", "expert"}
	for _, lvl := range levels {
		marker = "  "
		if i.filterMenuIndex == currentIndex {
			marker = "▶ "
		}
		text := lvl
		if i.filters.Level == lvl {
			text += " ✓"
		}
		if i.filterMenuIndex == currentIndex {
			lines = append(lines, selectedStyle.Render(marker+text))
		} else {
			lines = append(lines, normalStyle.Render(marker+text))
		}
		currentIndex++
	}

	lines = append(lines, "")
	lines = append(lines, mutedStyle.Render("Usage:"))

	// Used skills only option
	marker = "  "
	if i.filterMenuIndex == currentIndex {
		marker = "▶ "
	}
	text = "Used skills only"
	if i.filters.MinEvents > 0 {
		text += " ✓"
	}
	if i.filterMenuIndex == currentIndex {
		lines = append(lines, selectedStyle.Render(marker+text))
	} else {
		lines = append(lines, normalStyle.Render(marker+text))
	}

	content := lipgloss.JoinVertical(lipgloss.Left, lines...)
	return i.getCardStyle().Render(content)
}

// renderSortMenu renders the sort menu
func (i *ManageSkillsIntent) renderSortMenu() string {
	theme := i.Theme()

	titleStyle := lipgloss.NewStyle().
		Foreground(theme.PrimaryColor()).
		Bold(true).
		MarginBottom(1)

	selectedStyle := lipgloss.NewStyle().
		Foreground(theme.SuccessColor()).
		Bold(true)

	normalStyle := lipgloss.NewStyle().
		Foreground(theme.ForegroundColor())

	var lines []string
	lines = append(lines, titleStyle.Render("Sort Skills"))
	lines = append(lines, "")

	sortOptions := []struct {
		label     string
		sortBy    string
		sortOrder string
	}{
		{"Name (A-Z)", "name", "asc"},
		{"Name (Z-A)", "name", "desc"},
		{"Most used (Event count)", "events", "desc"},
		{"Least used (Event count)", "events", "asc"},
		{"Category (A-Z)", "category", "asc"},
		{"Category (Z-A)", "category", "desc"},
	}

	for idx, opt := range sortOptions {
		marker := "  "
		if i.sortMenuIndex == idx {
			marker = "▶ "
		}
		text := opt.label
		if i.filters.SortBy == opt.sortBy && i.filters.SortOrder == opt.sortOrder {
			text += " ✓"
		}
		if i.sortMenuIndex == idx {
			lines = append(lines, selectedStyle.Render(marker+text))
		} else {
			lines = append(lines, normalStyle.Render(marker+text))
		}
	}

	content := lipgloss.JoinVertical(lipgloss.Left, lines...)
	return i.getCardStyle().Render(content)
}

// ListNavigator interface implementation

// GetTotalItems returns the total number of skills.
func (i *ManageSkillsIntent) GetTotalItems() int {
	return len(i.skills)
}

// GetSelectedIndex returns the current selection index.
func (i *ManageSkillsIntent) GetSelectedIndex() int {
	return i.selectedIndex
}

// SetSelectedIndex sets the selection index and updates the display.
func (i *ManageSkillsIntent) SetSelectedIndex(idx int) {
	// Validate and set index
	if idx < 0 {
		idx = 0
	}
	if idx >= len(i.skills) {
		idx = len(i.skills) - 1
	}
	if idx < 0 {
		idx = 0 // Handle empty list
	}

	i.selectedIndex = idx

	// Update selected skill
	if idx >= 0 && idx < len(i.skills) {
		i.selectedSkill = i.skills[idx]
	}

	// Update table display
	i.updateTableRows()
}

// GetPageSize returns the page size for pagination.
func (i *ManageSkillsIntent) GetPageSize() int {
	return 15
}

// ============================================================================
// Screen Orchestration Methods (New Architecture)
// ============================================================================

// EnableScreens enables the screen-based architecture for this intent.
// This is an opt-in method to maintain backward compatibility during migration.
func (i *ManageSkillsIntent) EnableScreens() {
	i.useScreens = true
}

// handleScreenResult processes results from screen updates.
// This is the central hub for all screen-to-intent communication.
func (i *ManageSkillsIntent) handleScreenResult(result screens.ScreenResult) tea.Cmd {
	switch r := result.(type) {
	case *screens.NavigateResult:
		return i.handleNavigateResult(r)
	case *screens.CancelResult:
		return i.handleCancelResult(r)
	case *screens.SubmitResult:
		return i.handleSubmitResult(r)
	case *screens.ErrorResult:
		return i.handleErrorResult(r)
	default:
		// Unknown result type - treat as error
		i.result = &IntentResult[*ManageSkillsResult]{
			Status: Failed,
			Error: &IntentError{
				Code:    "UNKNOWN_RESULT",
				Message: fmt.Sprintf("unknown screen result type: %T", result),
			},
		}
		i.active = false
		return nil
	}
}

// handleNavigateResult handles navigation to a new screen.
// NavigateResult.Data() contains a map with "target" and "data" keys.
func (i *ManageSkillsIntent) handleNavigateResult(result *screens.NavigateResult) tea.Cmd {
	data := result.Data()

	// Extract target and data from navigation result
	dataMap, ok := data.(map[string]interface{})
	if !ok {
		return i.handleErrorInternal(fmt.Errorf("invalid navigation data: expected map, got %T", data))
	}

	target, _ := dataMap["target"].(string)
	skillData := dataMap["data"]

	switch target {
	case "detail":
		// Navigate to skill detail screen
		if skill, ok := skillData.(*domain.Skill); ok {
			i.selectedSkill = skill
			return i.transitionToDetailScreen()
		}
		return i.handleErrorInternal(fmt.Errorf("invalid data for detail screen: expected *domain.Skill, got %T", skillData))

	case "add":
		// Navigate to add skill form
		return i.transitionToFormScreen(nil)

	case "edit":
		// Navigate to edit skill form
		if skill, ok := skillData.(*domain.Skill); ok {
			i.selectedSkill = skill
			return i.transitionToFormScreen(skill)
		}
		return i.handleErrorInternal(fmt.Errorf("invalid data for edit screen: expected *domain.Skill, got %T", skillData))

	case "delete":
		// Navigate to delete confirmation
		if skill, ok := skillData.(*domain.Skill); ok {
			i.selectedSkill = skill
			return i.transitionToDeleteScreen(skill)
		}
		return i.handleErrorInternal(fmt.Errorf("invalid data for delete screen: expected *domain.Skill, got %T", skillData))

	case "list":
		// Navigate back to list (after delete/cancel)
		return i.transitionToListScreen()

	default:
		return i.handleErrorInternal(fmt.Errorf("unknown navigation target: %s", target))
	}
}

// handleCancelResult handles screen cancellation.
func (i *ManageSkillsIntent) handleCancelResult(result *screens.CancelResult) tea.Cmd {
	// Check if we should return to previous screen or exit intent
	if i.currentState == SkillsStateList {
		// Root state - cancel the entire intent
		i.result = &IntentResult[*ManageSkillsResult]{
			Status: Cancelled,
			Data: &ManageSkillsResult{
				Action: "cancelled",
			},
		}
		i.active = false
		return nil
	}

	// Return to list screen
	return i.transitionToListScreen()
}

// handleSubmitResult handles form/confirm submission.
func (i *ManageSkillsIntent) handleSubmitResult(result *screens.SubmitResult) tea.Cmd {
	data := result.Data()

	switch i.currentState {
	case SkillsStateAdd, SkillsStateEdit:
		// Form submitted - save skill
		if skill, ok := data.(*domain.Skill); ok {
			if i.currentState == SkillsStateAdd {
				return i.createSkill(skill)
			}
			return i.updateSkill(skill)
		}
		return i.handleErrorInternal(fmt.Errorf("invalid form data: expected *domain.Skill, got %T", data))

	case SkillsStateDelete:
		// Delete confirmed
		if skill, ok := data.(*domain.Skill); ok {
			return func() tea.Msg {
				err := i.context.SkillRepository.Delete(i.context.Ctx, skill.ID)
				return SkillDeletedMsg{
					SkillID: skill.ID,
					Error:   err,
				}
			}
		}
		return i.handleErrorInternal(fmt.Errorf("invalid delete data: expected *domain.Skill, got %T", data))

	default:
		return i.handleErrorInternal(fmt.Errorf("unexpected submit in state: %s", i.currentState))
	}
}

// handleErrorResult handles screen errors.
func (i *ManageSkillsIntent) handleErrorResult(result *screens.ErrorResult) tea.Cmd {
	// Screen encountered an error - propagate to intent
	data := result.Data()
	if err, ok := data.(error); ok {
		return i.handleErrorInternal(err)
	}
	return i.handleErrorInternal(fmt.Errorf("screen error: %v", data))
}

// applyIntentContextToScreen applies terminal info, theme, and logo to a screen.
// This ensures consistent setup across all screen transitions.
func (i *ManageSkillsIntent) applyIntentContextToScreen(screen screens.Screen) {
	// Set terminal dimensions
	if termInfo := i.GetTerminalInfo(); termInfo != nil {
		screen.SetTerminalInfo(termInfo.Width, termInfo.Height)
	}

	// Set theme
	if theme := i.Theme(); theme != nil {
		screen.SetTheme(theme)
	}

	// Set logo
	if logo := i.GetLogo(); logo != nil {
		screen.SetLogo(logo, i.GetLogoSpacing())
	}
}

// transitionToListScreen transitions to the skills list screen.
func (i *ManageSkillsIntent) transitionToListScreen() tea.Cmd {
	i.currentState = SkillsStateList

	// Create list screen
	listScreen := NewSkillsListScreenFromIntent(i.skills, i.GetThemeManager())

	// Set event counts if available (CRITICAL for event count column)
	if skillScreen, ok := listScreen.(*skills_screens.SkillsListScreen); ok {
		if i.eventCounts != nil {
			skillScreen.SetEventCounts(i.eventCounts)
		}
	}

	// Apply intent context (terminal, theme, logo)
	i.applyIntentContextToScreen(listScreen)

	i.activeScreen = listScreen
	return nil
}

// transitionToDetailScreen transitions to the skill detail screen.
func (i *ManageSkillsIntent) transitionToDetailScreen() tea.Cmd {
	i.currentState = SkillsStateDetail

	detailScreen := NewSkillDetailScreenFromIntent(i.selectedSkill, i.GetThemeManager())

	// Apply intent context (terminal, theme, logo)
	i.applyIntentContextToScreen(detailScreen)

	i.activeScreen = detailScreen
	return nil
}

// transitionToFormScreen transitions to the add/edit form screen.
func (i *ManageSkillsIntent) transitionToFormScreen(skill *domain.Skill) tea.Cmd {
	if skill == nil {
		i.currentState = SkillsStateAdd
	} else {
		i.currentState = SkillsStateEdit
	}

	formScreen := NewSkillFormScreenFromIntent(skill, i.GetThemeManager())

	// Apply intent context (terminal, theme, logo)
	i.applyIntentContextToScreen(formScreen)

	i.activeScreen = formScreen
	return nil
}

// transitionToDeleteScreen transitions to the delete confirmation screen.
func (i *ManageSkillsIntent) transitionToDeleteScreen(skill *domain.Skill) tea.Cmd {
	i.currentState = SkillsStateDelete

	deleteScreen := NewSkillDeleteConfirmScreenFromIntent(skill, i.GetThemeManager())

	// Apply intent context (terminal, theme, logo)
	i.applyIntentContextToScreen(deleteScreen)

	i.activeScreen = deleteScreen
	return nil
}

// handleErrorInternal handles errors by transitioning to error state.
func (i *ManageSkillsIntent) handleErrorInternal(err error) tea.Cmd {
	i.result = &IntentResult[*ManageSkillsResult]{
		Status: Failed,
		Error: &IntentError{
			Code:    "SCREEN_ERROR",
			Message: err.Error(),
			Cause:   err,
		},
	}
	// Don't deactivate - allow retry by returning to list
	return i.transitionToListScreen()
}

// ============================================================================
// Screen Constructor Wrappers (Avoid Import Cycles)
// ============================================================================

// NewSkillsListScreenFromIntent creates a SkillsListScreen from intent context.
// This wrapper avoids import cycles between intents and screens/skills packages.
func NewSkillsListScreenFromIntent(skills []*domain.Skill, themeManager interface{}) screens.Screen {
	return skills_screens.NewSkillsListScreen(skills)
}

// NewSkillDetailScreenFromIntent creates a SkillDetailScreen from intent context.
func NewSkillDetailScreenFromIntent(skill *domain.Skill, themeManager interface{}) screens.Screen {
	return skills_screens.NewSkillDetailScreen(skill)
}

// NewSkillFormScreenFromIntent creates a SkillFormScreen from intent context.
func NewSkillFormScreenFromIntent(skill *domain.Skill, themeManager interface{}) screens.Screen {
	return skills_screens.NewSkillFormScreen(skill)
}

// NewSkillDeleteConfirmScreenFromIntent creates a SkillDeleteConfirmScreen from intent context.
func NewSkillDeleteConfirmScreenFromIntent(skill *domain.Skill, themeManager interface{}) screens.Screen {
	return skills_screens.NewSkillDeleteConfirmScreen(skill)
}
