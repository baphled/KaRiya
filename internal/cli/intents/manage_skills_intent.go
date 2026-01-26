package intents

import (
	"fmt"
	"strings"
	"time"

	"github.com/baphled/kariya/internal/cli/behaviors"
	"github.com/baphled/kariya/internal/cli/components"
	"github.com/baphled/kariya/internal/cli/forms"
	"github.com/baphled/kariya/internal/cli/models"
	"github.com/baphled/kariya/internal/cli/screens"
	skills_screens "github.com/baphled/kariya/internal/cli/screens/skills"
	"github.com/baphled/kariya/internal/cli/terminal"
	"github.com/baphled/kariya/internal/cli/themes"
	"github.com/baphled/kariya/internal/cli/uikit/containers"
	"github.com/baphled/kariya/internal/cli/uikit/feedback"
	"github.com/baphled/kariya/internal/cli/uikit/primitives"
	domain "github.com/baphled/kariya/internal/domain/career"
	career "github.com/baphled/kariya/internal/repository/career"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
)

// Ensure ManageSkillsIntent implements FilterBehavior interface
var _ behaviors.FilterBehavior = (*ManageSkillsIntent)(nil)

// Ensure ManageSkillsIntent implements ScreenResultHandler interface
var _ behaviors.ScreenResultHandler = (*ManageSkillsIntent)(nil)

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

	// tableBehavior provides type-safe table operations for skills
	tableBehavior *behaviors.TableBehavior[*domain.Skill]

	// detail view data
	eventCounts  map[string]int        // Skill ID -> event count
	lastUsedMap  map[string]time.Time  // Skill ID -> last used date
	skillEvents  []*domain.CareerEvent // Events for selected skill
	eventsLoaded bool                  // Whether events have been loaded

	// eventsTableBehavior provides type-safe table operations for skill events
	eventsTableBehavior   *behaviors.TableBehavior[*domain.CareerEvent]
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
	searchModal *components.SkillSearchModal

	// view/edit modals (modal overlay architecture like BrowseTimelineIntent)
	viewDetailModal  *components.ViewSkillDetailModal
	addEditModal     *components.SkillAddEditModal
	deleteModal      *feedback.ConfirmModal
	skillEventsModal *components.ViewSkillEventsModal // Modal to show events using a skill
	eventDetailModal *components.ViewEventDetailModal // Modal to show event details from events list

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
	Category   string
	Level      string
	MinEvents  int
	SearchText string
	SortBy     string
	SortOrder  string
}

// NewManageSkillsIntent creates a new ManageSkills intent

// skillRowFormatterWithCounts creates a row formatter that includes event counts
func skillRowFormatterWithCounts(eventCounts map[string]int) behaviors.RowFormatter[*domain.Skill] {
	return func(skill *domain.Skill, index int) []string {
		// Name
		name := skill.Name
		if len(name) > 22 {
			name = name[:22] + "..."
		}

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
		if eventCounts != nil {
			if count, ok := eventCounts[skill.ID]; ok {
				eventCount = fmt.Sprintf("%d", count)
			}
		}

		return []string{name, category, level, years, eventCount}
	}
}

// eventRowFormatter formats a career event for table display
func eventRowFormatter(event *domain.CareerEvent, index int) []string {
	// Date
	dateStr := event.Date.Format("2006-01-02")

	// Truncate text to first 47 chars (50 - 3 for "...")
	text := event.Text
	if len(text) > 47 {
		text = text[:47] + "..."
	}

	// Company
	company := event.Company
	if company == "" {
		company = "-"
	}

	return []string{dateStr, text, company}
}

// NewManageSkillsIntent creates a new ManageSkills intent
func NewManageSkillsIntent(ctx *ManageSkillsContext) *ManageSkillsIntent {
	baseIntent := NewBaseIntent()
	baseIntent.SetThemeManager(themes.NewThemeManager())

	// Create TableBehavior for skills list
	skillsColumns := []behaviors.ColumnDef{
		{Title: "Name", Width: 25},
		{Title: "Category", Width: 15},
		{Title: "Level", Width: 12},
		{Title: "Years", Width: 8},
		{Title: "Events", Width: 8},
	}

	// Create TableBehavior for skill events list
	eventsColumns := []behaviors.ColumnDef{
		{Title: "Date", Width: 12},
		{Title: "Event", Width: 50},
		{Title: "Company", Width: 20},
	}

	intent := &ManageSkillsIntent{
		BaseIntent:   baseIntent,
		context:      ctx,
		currentState: SkillsStateList,
		skills:       []*domain.Skill{},
		filters:      &SkillsFilters{},
		active:       true,
	}

	// Initialize TableBehaviors - use closure to capture eventCounts reference
	intent.tableBehavior = behaviors.NewTableBehavior[*domain.Skill](nil, skillsColumns, skillRowFormatterWithCounts(intent.eventCounts)).
		PageSize(15).
		PaginationPrefix("Skills").
		EmptyMessage("No skills found. Press 'n' to add a new skill.")

	intent.eventsTableBehavior = behaviors.NewTableBehavior[*domain.CareerEvent](nil, eventsColumns, eventRowFormatter).
		PageSize(15).
		PaginationPrefix("Events").
		EmptyMessage("No events found for this skill.")

	return intent
}

// Init initializes the intent and loads skills
func (i *ManageSkillsIntent) Init() tea.Cmd {
	i.active = true

	// Disable screen architecture by default (tests expect legacy mode)
	// TODO: Fix screen orchestration bugs before re-enabling
	i.useScreens = false

	// Apply theme to TableBehaviors if available
	if theme := i.Theme(); theme != nil {
		i.tableBehavior.SetTheme(theme)
		i.eventsTableBehavior.SetTheme(theme)
	}

	// Load skills asynchronously
	return func() tea.Msg {
		// Guard against nil repository (e.g., in tests without full context setup)
		if i.context == nil || i.context.SkillRepository == nil {
			return SkillsLoadedMsg{
				Skills: nil,
				Error:  nil,
			}
		}
		skills, err := i.context.SkillRepository.List(i.context.Ctx, nil)
		return SkillsLoadedMsg{
			Skills: skills,
			Error:  err,
		}
	}
}

// syncTableSelection syncs the TableBehavior selection with the intent's data
func (i *ManageSkillsIntent) syncTableSelection() {
	i.selectedIndex = i.tableBehavior.GetSelectedIndex()
	if selected := i.tableBehavior.GetSelectedItem(); selected != nil {
		i.selectedSkill = *selected
	} else {
		i.selectedSkill = nil
	}
}

// syncEventsTableSelection syncs the events TableBehavior selection with the intent's data
func (i *ManageSkillsIntent) syncEventsTableSelection() {
	i.eventsSelectedIndex = i.eventsTableBehavior.GetSelectedIndex()
	if selected := i.eventsTableBehavior.GetSelectedItem(); selected != nil {
		i.selectedEventFromList = *selected
	} else {
		i.selectedEventFromList = nil
	}
}

// refreshSkillsTable updates the TableBehavior with current skills
func (i *ManageSkillsIntent) refreshSkillsTable() {
	i.tableBehavior.SetItems(i.skills)
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
	}

	// 2. SECOND PRIORITY: Modal updates (if visible)
	// CRITICAL: Pass full tea.Msg (not tea.KeyMsg) to modals
	// This allows huh forms to process Tab/Enter correctly
	if i.searchModal != nil && i.searchModal.IsVisible() {
		return i.handleSearchModalUpdate(msg)
	}
	if i.filterModal != nil && i.filterModal.IsVisible() {
		return i.handleFilterModalUpdate(msg)
	}
	if i.sortModal != nil && i.sortModal.IsVisible() {
		return i.handleSortModalUpdate(msg)
	}

	// Handle view detail modal
	if i.viewDetailModal != nil && i.viewDetailModal.IsVisible() {
		return i.handleViewDetailModalUpdate(msg)
	}

	// Handle add/edit modal
	if i.addEditModal != nil && i.addEditModal.IsVisible() {
		return i.handleAddEditModalUpdate(msg)
	}

	// Handle delete confirmation modal
	if i.deleteModal != nil && i.deleteModal.IsVisible() {
		return i.handleDeleteModalUpdate(msg)
	}

	// Handle skill events modal (shows events using a skill)
	if i.skillEventsModal != nil && i.skillEventsModal.IsVisible() {
		return i.handleSkillEventsModalUpdate(msg)
	}

	// Handle event detail modal (shows details of an event from events list)
	if i.eventDetailModal != nil && i.eventDetailModal.IsVisible() {
		return i.handleEventDetailModalUpdate(msg)
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

	case SkillEventsForModalLoadedMsg:
		return i.handleSkillEventsForModalLoaded(msg)

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
	if i.searchModal != nil && i.searchModal.IsVisible() {
		return i.renderSearchModalOverlay(baseView)
	}
	if i.filterModal != nil && i.filterModal.IsVisible() {
		return i.renderFilterModalOverlay(baseView)
	}
	if i.sortModal != nil && i.sortModal.IsVisible() {
		return i.renderSortModalOverlay(baseView)
	}
	if i.viewDetailModal != nil && i.viewDetailModal.IsVisible() {
		return i.renderViewDetailModalOverlay(baseView)
	}
	if i.addEditModal != nil && i.addEditModal.IsVisible() {
		return i.renderAddEditModalOverlay(baseView)
	}
	if i.deleteModal != nil && i.deleteModal.IsVisible() {
		return i.renderDeleteModalOverlay(baseView)
	}
	if i.skillEventsModal != nil && i.skillEventsModal.IsVisible() {
		return i.renderSkillEventsModalOverlay(baseView)
	}
	if i.eventDetailModal != nil && i.eventDetailModal.IsVisible() {
		return i.renderEventDetailModalOverlay(baseView)
	}

	return baseView
}

// renderFilterModalOverlay renders the filter modal centered on the background.
func (i *ManageSkillsIntent) renderFilterModalOverlay(baseView string) string {
	return behaviors.RenderModalOverlay(i.filterModal, baseView)
}

// renderSortModalOverlay renders the sort modal centered on the background.
func (i *ManageSkillsIntent) renderSortModalOverlay(baseView string) string {
	return behaviors.RenderModalOverlay(i.sortModal, baseView)
}

// renderSearchModalOverlay renders the search modal centered on the background.
func (i *ManageSkillsIntent) renderSearchModalOverlay(baseView string) string {
	return behaviors.RenderModalOverlay(i.searchModal, baseView)
}

// renderViewDetailModalOverlay renders the view detail modal centered on the background.
func (i *ManageSkillsIntent) renderViewDetailModalOverlay(baseView string) string {
	return behaviors.RenderModalOverlay(i.viewDetailModal, baseView)
}

// renderAddEditModalOverlay renders the add/edit modal centered on the background.
func (i *ManageSkillsIntent) renderAddEditModalOverlay(baseView string) string {
	return behaviors.RenderModalOverlay(i.addEditModal, baseView)
}

// renderDeleteModalOverlay renders the delete confirmation modal centered on the background.
func (i *ManageSkillsIntent) renderDeleteModalOverlay(baseView string) string {
	return behaviors.RenderModalOverlay(i.deleteModal, baseView)
}

// renderSkillEventsModalOverlay renders the skill events modal centered on the background.
func (i *ManageSkillsIntent) renderSkillEventsModalOverlay(baseView string) string {
	return behaviors.RenderModalOverlay(i.skillEventsModal, baseView)
}

// renderEventDetailModalOverlay renders the event detail modal centered on the background.
func (i *ManageSkillsIntent) renderEventDetailModalOverlay(baseView string) string {
	return behaviors.RenderModalOverlay(i.eventDetailModal, baseView)
}

// handleFilterModalUpdate handles updates when filter modal is visible
// CRITICAL: Takes tea.Msg (not tea.KeyMsg) to allow huh forms to work correctly
func (i *ManageSkillsIntent) handleFilterModalUpdate(msg tea.Msg) tea.Cmd {
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
		// NOTE: Sort is handled by SkillSortModal separately

		// Reload skills with new filters
		return i.reloadSkills()
	}

	// Modal was closed without completion (Esc) or still being edited
	return cmd
}

// handleSortModalUpdate handles updates when sort modal is visible
// CRITICAL: Takes tea.Msg (not tea.KeyMsg) to allow huh forms to work correctly
func (i *ManageSkillsIntent) handleSortModalUpdate(msg tea.Msg) tea.Cmd {
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

// openSearchModal opens the search modal with current search text pre-populated
// PATTERN 12: Form Modal with Immediate Init
func (i *ManageSkillsIntent) openSearchModal() tea.Cmd {
	// Get terminal dimensions
	termInfo := i.GetTerminalInfo()
	width := 120
	height := 40
	if termInfo != nil {
		width = termInfo.Width
		height = termInfo.Height
	}

	// Get current search text
	searchText := ""
	if i.filters != nil {
		searchText = i.filters.SearchText
	}

	// Create search modal
	i.searchModal = components.NewSkillSearchModal(
		searchText,
		width,
		height,
	)

	// CRITICAL: Call Init() for immediate rendering
	return i.searchModal.Init()
}

// handleSearchModalUpdate handles updates when search modal is visible
// CRITICAL: Takes tea.Msg (not tea.KeyMsg) to allow huh forms to work correctly
func (i *ManageSkillsIntent) handleSearchModalUpdate(msg tea.Msg) tea.Cmd {
	cmd, applied, searchData := i.searchModal.Update(msg)

	if applied && searchData != nil {
		// User confirmed search - apply it
		if i.filters == nil {
			i.filters = &SkillsFilters{}
		}
		i.filters.SearchText = searchData.SearchText

		// Reload skills with new search
		return i.reloadSkills()
	}

	// Modal was closed without completion (Esc) or still being edited
	return cmd
}

// handleViewDetailModalUpdate handles updates when view detail modal is visible
func (i *ManageSkillsIntent) handleViewDetailModalUpdate(msg tea.Msg) tea.Cmd {
	_, cmd := i.viewDetailModal.Update(msg)

	if !i.viewDetailModal.IsVisible() {
		// Modal was closed - check for action
		action := i.viewDetailModal.GetAction()
		switch action {
		case "events":
			// Load events for this skill and show modal
			i.viewDetailModal = nil
			return i.loadEventsForSkillModal()
		case "edit":
			// Open add/edit modal for this skill
			return i.openAddEditModal(i.selectedSkill)
		case "delete":
			// Open delete confirmation modal
			return i.openDeleteModal(i.selectedSkill)
		}
		// Simple close - clear the modal
		i.viewDetailModal = nil
	}

	return cmd
}

// handleAddEditModalUpdate handles updates when add/edit modal is visible
func (i *ManageSkillsIntent) handleAddEditModalUpdate(msg tea.Msg) tea.Cmd {
	cmd, completed, skillData := i.addEditModal.Update(msg)

	if !i.addEditModal.IsVisible() {
		if completed && skillData != nil {
			// User completed form - save skill
			originalSkill := i.addEditModal.GetOriginalSkill()
			if originalSkill != nil {
				// Editing existing skill
				skill := skillData.ToSkill(originalSkill.ID)
				i.addEditModal = nil
				return i.updateSkill(skill)
			} else {
				// Creating new skill
				skill := skillData.ToSkill("")
				i.addEditModal = nil
				return i.createSkill(skill)
			}
		}
		// User cancelled - close modal
		i.addEditModal = nil
	}

	return cmd
}

// handleDeleteModalUpdate handles updates when delete confirmation modal is visible
func (i *ManageSkillsIntent) handleDeleteModalUpdate(msg tea.Msg) tea.Cmd {
	cmd, confirmed := i.deleteModal.Update(msg)

	if !i.deleteModal.IsVisible() {
		if confirmed && i.selectedSkill != nil {
			// User confirmed deletion
			skillID := i.selectedSkill.ID
			i.deleteModal = nil
			return func() tea.Msg {
				err := i.context.SkillRepository.Delete(i.context.Ctx, skillID)
				return SkillDeletedMsg{
					SkillID: skillID,
					Error:   err,
				}
			}
		}
		// User cancelled
		i.deleteModal = nil
	}

	return cmd
}

// handleSkillEventsModalUpdate handles updates when skill events modal is visible
func (i *ManageSkillsIntent) handleSkillEventsModalUpdate(msg tea.Msg) tea.Cmd {
	_, cmd := i.skillEventsModal.Update(msg)

	if !i.skillEventsModal.IsVisible() {
		// Check if user selected an event
		if i.skillEventsModal.HasSelection() {
			selectedEvent := i.skillEventsModal.GetSelectedEvent()
			i.skillEventsModal.ClearSelection()
			// Open event detail modal
			return i.openEventDetailModal(selectedEvent)
		}
		// Simple close - clear the modal
		i.skillEventsModal = nil
	}

	return cmd
}

// handleEventDetailModalUpdate handles updates when event detail modal is visible
func (i *ManageSkillsIntent) handleEventDetailModalUpdate(msg tea.Msg) tea.Cmd {
	_, cmd := i.eventDetailModal.Update(msg)

	if !i.eventDetailModal.IsVisible() {
		// Event detail modal was closed
		i.eventDetailModal = nil
		// Re-show the skill events modal if it exists
		if i.skillEventsModal != nil {
			i.skillEventsModal.Show()
		}
	}

	return cmd
}

// openSkillEventsModal opens the skill events modal for the selected skill
func (i *ManageSkillsIntent) openSkillEventsModal(events []*domain.CareerEvent) tea.Cmd {
	if i.selectedSkill == nil {
		return nil
	}

	// Get terminal dimensions
	termInfo := i.GetTerminalInfo()
	width, height := 120, 40
	if termInfo != nil && termInfo.Width > 0 && termInfo.Height > 0 {
		width, height = termInfo.Width, termInfo.Height
	}

	i.skillEventsModal = components.NewViewSkillEventsModal(
		i.selectedSkill.ID,
		i.selectedSkill.Name,
		events,
		i.Theme(),
	)
	i.skillEventsModal.SetDimensions(width, height)
	i.skillEventsModal.Show()

	return nil
}

// openEventDetailModal opens the event detail modal for a selected event
func (i *ManageSkillsIntent) openEventDetailModal(event *domain.CareerEvent) tea.Cmd {
	if event == nil {
		return nil
	}

	// Get terminal dimensions
	termInfo := i.GetTerminalInfo()
	width, height := 120, 40
	if termInfo != nil && termInfo.Width > 0 && termInfo.Height > 0 {
		width, height = termInfo.Width, termInfo.Height
	}

	i.eventDetailModal = components.NewViewEventDetailModal(event, i.Theme()).
		WithShowSkillsOption(false) // Hide "s: Skills" - we're already in skills context
	i.eventDetailModal.SetDimensions(width, height)
	i.eventDetailModal.Show()

	return nil
}

// openViewDetailModal opens the view detail modal for the selected skill
func (i *ManageSkillsIntent) openViewDetailModal() tea.Cmd {
	if len(i.skills) == 0 || i.selectedIndex >= len(i.skills) {
		return nil
	}

	skill := i.skills[i.selectedIndex]
	i.selectedSkill = skill

	// Get event count and last used for this skill
	eventCount := 0
	if i.eventCounts != nil {
		eventCount = i.eventCounts[skill.ID]
	}

	var lastUsed *time.Time
	if i.lastUsedMap != nil {
		if lu, ok := i.lastUsedMap[skill.ID]; ok {
			lastUsed = &lu
		}
	}

	// Get terminal dimensions
	termInfo := i.GetTerminalInfo()
	width, height := 120, 40
	if termInfo != nil && termInfo.Width > 0 && termInfo.Height > 0 {
		width, height = termInfo.Width, termInfo.Height
	}

	i.viewDetailModal = components.NewViewSkillDetailModal(skill, i.Theme(), eventCount, lastUsed)
	i.viewDetailModal.SetDimensions(width, height)
	i.viewDetailModal.Show()

	return nil
}

// openAddEditModal opens the add/edit modal for a skill
func (i *ManageSkillsIntent) openAddEditModal(skill *domain.Skill) tea.Cmd {
	// Get terminal dimensions
	termInfo := i.GetTerminalInfo()
	width, height := 120, 40
	if termInfo != nil && termInfo.Width > 0 && termInfo.Height > 0 {
		width, height = termInfo.Width, termInfo.Height
	}

	i.addEditModal = components.NewSkillAddEditModal(skill, width, height)
	return i.addEditModal.Init()
}

// openDeleteModal opens the delete confirmation modal for a skill
func (i *ManageSkillsIntent) openDeleteModal(skill *domain.Skill) tea.Cmd {
	if skill == nil {
		return nil
	}

	i.selectedSkill = skill
	skillName := skill.Name
	if len(skillName) > 50 {
		skillName = skillName[:47] + "..."
	}

	i.deleteModal = feedback.NewConfirmModal(
		"Delete Skill",
		fmt.Sprintf("Are you sure you want to delete '%s'?", skillName),
	).WithVariant(feedback.ConfirmDestructive)
	return i.deleteModal.Init()
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
		badges := []*primitives.Badge{
			primitives.HelpKeyBadge("Enter", "View details", theme),
			primitives.HelpKeyBadge("n", "New skill", theme),
			primitives.HelpKeyBadge("f", "Filter", theme),
			primitives.HelpKeyBadge("s", "Sort", theme),
		}
		if i.HasActiveFilters() {
			badges = append(badges, primitives.HelpKeyBadge("x", "Clear filters", theme))
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
				primitives.HelpKeyBadge("Enter", "View events", theme),
				primitives.EditBadge(theme),
				primitives.DeleteBadge(theme),
			),
			ThemedGlobalBadges(theme),
		)
	case SkillsStateDetailEvents:
		return CombineThemedFooters(
			ThemedListFooter(theme),
			ThemedCustomFooter(theme,
				primitives.HelpKeyBadge("Enter", "View details", theme),
				primitives.EditBadge(theme),
			),
			ThemedGlobalBadges(theme),
		)
	case SkillsStateDetailEventDetail:
		return CombineThemedFooters(
			ThemedDetailViewFooter(theme),
			ThemedCustomFooter(theme,
				primitives.EditBadge(theme),
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
				primitives.HelpKeyBadge("y", "Confirm", theme),
				primitives.HelpKeyBadge("n/Esc", "Cancel", theme),
			),
			ThemedGlobalBadges(theme),
		)
	case SkillsStateFilter:
		return CombineThemedFooters(
			ThemedListFooter(theme),
			ThemedCustomFooter(theme,
				primitives.HelpKeyBadge("Enter", "Apply", theme),
				primitives.HelpKeyBadge("u", "Used skills only", theme),
			),
			ThemedGlobalBadges(theme),
		)
	case SkillsStateSort:
		return CombineThemedFooters(
			ThemedListFooter(theme),
			ThemedCustomFooter(theme,
				primitives.HelpKeyBadge("Enter", "Apply", theme),
				primitives.HelpKeyBadge("e", "Most used", theme),
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

// HasVisibleDetailModal returns true if the detail modal is visible (for testing)
func (i *ManageSkillsIntent) HasVisibleDetailModal() bool {
	return i.viewDetailModal != nil && i.viewDetailModal.IsVisible()
}

// HasVisibleAddEditModal returns true if the add/edit modal is visible (for testing)
func (i *ManageSkillsIntent) HasVisibleAddEditModal() bool {
	return i.addEditModal != nil && i.addEditModal.IsVisible()
}

// HasVisibleDeleteModal returns true if the delete modal is visible (for testing)
func (i *ManageSkillsIntent) HasVisibleDeleteModal() bool {
	return i.deleteModal != nil && i.deleteModal.IsVisible()
}

// HasVisibleSkillEventsModal returns true if the skill events modal is visible (for testing)
func (i *ManageSkillsIntent) HasVisibleSkillEventsModal() bool {
	return i.skillEventsModal != nil && i.skillEventsModal.IsVisible()
}

// HasVisibleEventDetailModal returns true if the event detail modal is visible (for testing)
func (i *ManageSkillsIntent) HasVisibleEventDetailModal() bool {
	return i.eventDetailModal != nil && i.eventDetailModal.IsVisible()
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

	// Apply in-memory search filtering if search text is set
	if i.filters != nil && i.filters.SearchText != "" {
		i.skills = i.applySearchFilter(i.skills, i.filters.SearchText)
	}

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

	// Update TableBehavior with new skills data
	i.refreshSkillsTable()
	i.syncTableSelection()

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
			// Try TableBehavior navigation first (handles j/k, up/down, pgup/pgdn, home/end, g/G)
			if i.tableBehavior.HandleNavigation(msg.String()) {
				i.syncTableSelection()
				return nil
			}

			switch msg.String() {
			case "enter":
				// View skill detail - use modal overlay
				if len(i.skills) == 0 {
					return nil
				}
				return i.openViewDetailModal()

			case "n":
				// Add new skill - use modal overlay
				return i.openAddEditModal(nil)

			case "e":
				// Edit selected skill - use modal overlay
				if len(i.skills) == 0 {
					return nil
				}
				i.selectedSkill = i.skills[i.selectedIndex]
				return i.openAddEditModal(i.selectedSkill)

			case "d":
				// Delete selected skill - use modal overlay
				if len(i.skills) == 0 {
					return nil
				}
				return i.openDeleteModal(i.skills[i.selectedIndex])

			case "f":
				// PATTERN 12: Form Modal with Immediate Init
				// Open filter modal instead of menu state
				return i.openFilterModal()

			case "s":
				// PATTERN 12: Form Modal with Immediate Init
				// Open sort modal instead of menu state
				return i.openSortModal()

			case "/":
				// PATTERN 12: Form Modal with Immediate Init
				// Open search modal
				return i.openSearchModal()

			case "x":
				// Clear filters in FIFO order
				if i.HasActiveFilters() {
					i.ClearFilters()
					return i.RefreshData()
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

// hasActiveFilters returns true if any filters are active (private implementation)
func (i *ManageSkillsIntent) hasActiveFilters() bool {
	if i.filters == nil {
		return false
	}
	return i.filters.Category != "" ||
		i.filters.Level != "" ||
		i.filters.MinEvents > 0 ||
		i.filters.SortBy != "" ||
		i.filters.SearchText != ""
}

// HasActiveFilters returns true if any non-default filters are active.
// Implements FilterBehavior interface.
func (i *ManageSkillsIntent) HasActiveFilters() bool {
	return i.hasActiveFilters()
}

// ClearFilters resets filters in FIFO order (most recent filter first).
// Implements FilterBehavior interface.
func (i *ManageSkillsIntent) ClearFilters() {
	if i.filters == nil {
		return
	}

	// Clear in FIFO order: search → filter → sort
	// Search is most recent (most specific), sort is least recent (most general)
	if i.filters.SearchText != "" {
		i.filters.SearchText = ""
		return
	}

	if i.filters.Category != "" || i.filters.Level != "" || i.filters.MinEvents > 0 {
		i.filters.Category = ""
		i.filters.Level = ""
		i.filters.MinEvents = 0
		return
	}

	// Clear sort (least specific)
	i.filters.SortBy = ""
	i.filters.SortOrder = ""
}

// ApplyFilters applies current filter state to the data.
// Implements FilterBehavior interface.
func (i *ManageSkillsIntent) ApplyFilters() {
	// Apply search filter to current skills list
	if i.filters != nil && i.filters.SearchText != "" {
		i.skills = i.applySearchFilter(i.skills, i.filters.SearchText)
	}
}

// RefreshData reloads/refreshes the filtered data.
// Implements FilterBehavior interface.
func (i *ManageSkillsIntent) RefreshData() tea.Cmd {
	return i.reloadSkills()
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

// applySearchFilter applies in-memory search filtering to skills
func (i *ManageSkillsIntent) applySearchFilter(skills []*domain.Skill, searchText string) []*domain.Skill {
	if searchText == "" {
		return skills
	}

	// Case-insensitive search across name and category
	searchLower := strings.ToLower(searchText)
	filtered := make([]*domain.Skill, 0)
	for _, skill := range skills {
		if strings.Contains(strings.ToLower(skill.Name), searchLower) ||
			strings.Contains(strings.ToLower(skill.Category), searchLower) {
			filtered = append(filtered, skill)
		}
	}
	return filtered
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
	modal := feedback.NewWarningModal("Delete Skill", modalContent)

	// Get terminal dimensions
	width, height := 80, 24
	if termInfo := i.GetTerminalInfo(); termInfo != nil && termInfo.IsValid {
		width = termInfo.Width
		height = termInfo.Height
	}

	return modal.Render(width, height)
}

func (i *ManageSkillsIntent) renderSkillsList() string {
	// TableBehavior handles empty state and pagination internally
	return i.tableBehavior.Render()
}

// handleSkillEventsLoaded handles the SkillEventsLoadedMsg (state-based flow)
func (i *ManageSkillsIntent) handleSkillEventsLoaded(msg SkillEventsLoadedMsg) tea.Cmd {
	if msg.Error != nil {
		// Show error but stay in detail view
		return nil
	}

	i.skillEvents = msg.Events
	i.eventsLoaded = true
	i.eventsSelectedIndex = 0

	// Apply theme to events TableBehavior
	if theme := i.Theme(); theme != nil {
		i.eventsTableBehavior.SetTheme(theme)
	}

	// Update TableBehavior with events data
	i.eventsTableBehavior.SetItems(i.skillEvents)
	i.syncEventsTableSelection()

	return nil
}

// handleSkillEventsForModalLoaded handles the SkillEventsForModalLoadedMsg (modal flow)
func (i *ManageSkillsIntent) handleSkillEventsForModalLoaded(msg SkillEventsForModalLoadedMsg) tea.Cmd {
	if msg.Error != nil {
		// Show error - could display error modal here
		return nil
	}

	// Open the skill events modal with the loaded events
	return i.openSkillEventsModal(msg.Events)
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
			// Try TableBehavior navigation first (handles j/k, up/down, pgup/pgdn, home/end, g/G)
			if i.eventsTableBehavior.HandleNavigation(msg.String()) {
				i.syncEventsTableSelection()
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

// loadEventsForSkill loads events that use the selected skill (for state-based flow)
func (i *ManageSkillsIntent) loadEventsForSkill() tea.Cmd {
	return func() tea.Msg {
		events, err := i.context.SkillRepository.GetEventsUsingSkill(i.context.Ctx, i.selectedSkill.ID)
		return SkillEventsLoadedMsg{
			Events: events,
			Error:  err,
		}
	}
}

// loadEventsForSkillModal loads events and opens the skill events modal
func (i *ManageSkillsIntent) loadEventsForSkillModal() tea.Cmd {
	return func() tea.Msg {
		events, err := i.context.SkillRepository.GetEventsUsingSkill(i.context.Ctx, i.selectedSkill.ID)
		return SkillEventsForModalLoadedMsg{
			Events: events,
			Error:  err,
		}
	}
}

// renderSkillDetail renders the skill detail content using UIKit components
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

	// Timestamps using UIKit primitives
	lines = append(lines, "")
	lines = append(lines, labelStyle.Render("Created:")+primitives.Muted(skill.CreatedAt.Format("2006-01-02 15:04"), theme).Render())
	lines = append(lines, labelStyle.Render("Updated:")+primitives.Muted(skill.UpdatedAt.Format("2006-01-02 15:04"), theme).Render())

	content := lipgloss.JoinVertical(lipgloss.Left, lines...)
	// Use UIKit Box container for consistent card styling
	return containers.NewBox(theme).Content(content).Render()
}

// renderSkillEvents renders the events using this skill as a table
func (i *ManageSkillsIntent) renderSkillEvents() string {
	theme := i.Theme()

	if !i.eventsLoaded {
		// Use UIKit Box container with muted loading text
		loadingContent := primitives.Muted("Loading events...", theme).Render()
		return containers.NewBox(theme).Content(loadingContent).Render()
	}

	// TableBehavior handles empty state and pagination internally
	return i.eventsTableBehavior.Render()
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

	selectedStyle := lipgloss.NewStyle().
		Foreground(theme.SuccessColor()).
		Bold(true)

	normalStyle := lipgloss.NewStyle().
		Foreground(theme.ForegroundColor())

	mutedStyle := lipgloss.NewStyle().
		Foreground(theme.MutedColor())

	var lines []string
	lines = append(lines, primitives.Title("Filter Skills", theme).MarginBottom(1).Render())
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
	// Use UIKit Box container for consistent card styling
	return containers.NewBox(i.Theme()).Content(content).Render()
}

// renderSortMenu renders the sort menu
func (i *ManageSkillsIntent) renderSortMenu() string {
	theme := i.Theme()

	selectedStyle := lipgloss.NewStyle().
		Foreground(theme.SuccessColor()).
		Bold(true)

	normalStyle := lipgloss.NewStyle().
		Foreground(theme.ForegroundColor())

	var lines []string
	lines = append(lines, primitives.Title("Sort Skills", theme).MarginBottom(1).Render())
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
	// Use UIKit Box container for consistent card styling
	return containers.NewBox(i.Theme()).Content(content).Render()
}

// ListNavigator interface implementation (delegates to TableBehavior)

// GetTotalItems returns the total number of skills.
func (i *ManageSkillsIntent) GetTotalItems() int {
	return i.tableBehavior.Count()
}

// GetSelectedIndex returns the current selection index.
func (i *ManageSkillsIntent) GetSelectedIndex() int {
	return i.tableBehavior.GetSelectedIndex()
}

// SetSelectedIndex sets the selection index and updates the display.
func (i *ManageSkillsIntent) SetSelectedIndex(idx int) {
	i.tableBehavior.SetSelectedIndex(idx)
	i.syncTableSelection()
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
//
// Uses ScreenResultDispatcher pattern to eliminate repetitive type switching.
// ManageSkillsIntent implements ScreenResultHandler interface for compile-time safety.
func (i *ManageSkillsIntent) handleScreenResult(result screens.ScreenResult) tea.Cmd {
	return behaviors.NewScreenResultDispatcher(i).Dispatch(result)
}

// HandleNavigate handles navigation to a new screen.
// NavigateResult.Data() contains a map with "target" and "data" keys.
//
// Implements ScreenResultHandler interface.
func (i *ManageSkillsIntent) HandleNavigate(result *screens.NavigateResult) tea.Cmd {
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

// HandleCancel handles screen cancellation.
//
// Implements ScreenResultHandler interface.
func (i *ManageSkillsIntent) HandleCancel(result *screens.CancelResult) tea.Cmd {
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

// HandleSubmit handles form/confirm submission.
//
// Implements ScreenResultHandler interface.
func (i *ManageSkillsIntent) HandleSubmit(result *screens.SubmitResult) tea.Cmd {
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

// HandleError handles screen errors.
//
// Implements ScreenResultHandler interface.
func (i *ManageSkillsIntent) HandleError(result *screens.ErrorResult) tea.Cmd {
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
