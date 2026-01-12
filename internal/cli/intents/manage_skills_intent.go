package intents

import (
	"fmt"
	"sort"
	"strings"
	"time"

	"github.com/baphled/kariya/internal/cli/components"
	"github.com/baphled/kariya/internal/cli/forms"
	"github.com/baphled/kariya/internal/cli/models"
	"github.com/baphled/kariya/internal/cli/terminal"
	"github.com/baphled/kariya/internal/cli/themes"
	domain "github.com/baphled/kariya/internal/domain/career"
	career "github.com/baphled/kariya/internal/repository/career"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
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

	// detail view data
	eventCounts  map[string]int        // Skill ID -> event count
	lastUsedMap  map[string]time.Time  // Skill ID -> last used date
	skillEvents  []*domain.CareerEvent // Events for selected skill
	eventsLoaded bool                  // Whether events have been loaded

	// form for add/edit
	skillForm *models.SkillForm

	// filter and sort state
	filters             *SkillsFilters // Active filters
	filterMenuIndex     int            // Selected option in filter menu
	sortMenuIndex       int            // Selected option in sort menu
	availableCategories []string       // Categories extracted from skills for filter menu

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

	return &ManageSkillsIntent{
		BaseIntent:    baseIntent,
		context:       ctx,
		currentState:  SkillsStateList,
		skills:        []*domain.Skill{},
		selectedIndex: 0,
		filters:       &SkillsFilters{},
		active:        true,
	}
}

// Init initializes the intent and loads skills
func (i *ManageSkillsIntent) Init() tea.Cmd {
	i.active = true

	// Load skills asynchronously
	return func() tea.Msg {
		skills, err := i.context.SkillRepository.List(i.context.Ctx, nil)
		return SkillsLoadedMsg{
			Skills: skills,
			Error:  err,
		}
	}
}

// Update handles messages and state transitions
func (i *ManageSkillsIntent) Update(msg tea.Msg) tea.Cmd {
	if !i.active {
		return nil
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

	return view.Render()
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
	case SkillsStateDetail, SkillsStateDetailEvents:
		if i.selectedSkill != nil {
			breadcrumbs = append(breadcrumbs, i.selectedSkill.Name)
		}
		if i.currentState == SkillsStateDetailEvents {
			breadcrumbs = append(breadcrumbs, "Events")
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
			ThemedDetailViewFooter(theme),
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
	switch msg.String() {
	case "j", "down":
		if i.selectedIndex < len(i.skills)-1 {
			i.selectedIndex++
		}
		return nil

	case "k", "up":
		if i.selectedIndex > 0 {
			i.selectedIndex--
		}
		return nil

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
		// Open filter menu
		i.currentState = SkillsStateFilter
		i.filterMenuIndex = 0
		i.extractAvailableCategories()
		return nil

	case "s":
		// Open sort menu
		i.currentState = SkillsStateSort
		i.sortMenuIndex = 0
		return nil

	case "x":
		// Clear all filters
		if i.hasActiveFilters() {
			i.filters = &SkillsFilters{}
			return i.reloadSkills()
		}
		return nil

	case "esc":
		// Complete intent
		i.result = &IntentResult[*ManageSkillsResult]{
			Status: Cancelled,
			Data: &ManageSkillsResult{
				Action: "cancelled",
			},
		}
		i.active = false
		return nil
	}

	return nil
}

// handleFilterKeys handles key presses in filter menu
func (i *ManageSkillsIntent) handleFilterKeys(msg tea.KeyMsg) tea.Cmd {
	switch msg.String() {
	case "esc":
		i.currentState = SkillsStateList
		return nil

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
}

// handleSortKeys handles key presses in sort menu
func (i *ManageSkillsIntent) handleSortKeys(msg tea.KeyMsg) tea.Cmd {
	switch msg.String() {
	case "esc":
		i.currentState = SkillsStateList
		return nil

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

	case "n", "esc":
		// Cancel delete
		i.currentState = SkillsStateList
		return nil
	}

	return nil
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
		return i.renderEmptyState()
	}

	// Group skills by category
	grouped := i.groupSkillsByCategory()

	var sections []string

	// Sort categories alphabetically
	categories := make([]string, 0, len(grouped))
	for cat := range grouped {
		categories = append(categories, cat)
	}
	sort.Strings(categories)

	// Render each category
	for _, category := range categories {
		skills := grouped[category]
		sections = append(sections, i.renderCategory(category, skills))
	}

	content := lipgloss.JoinVertical(lipgloss.Left, sections...)
	return i.getCardStyle().Render(content)
}

func (i *ManageSkillsIntent) renderCategory(category string, skills []*domain.Skill) string {
	theme := i.Theme()
	categoryStyle := lipgloss.NewStyle().
		Foreground(theme.PrimaryColor()).
		Bold(true).
		MarginTop(1)

	var lines []string
	lines = append(lines, categoryStyle.Render("▸ "+category))

	for _, skill := range skills {
		lines = append(lines, i.renderSkillItem(skill))
	}

	return lipgloss.JoinVertical(lipgloss.Left, lines...)
}

func (i *ManageSkillsIntent) renderSkillItem(skill *domain.Skill) string {
	theme := i.Theme()

	// Check if this skill is selected
	isSelected := false
	for idx, s := range i.skills {
		if s.ID == skill.ID && idx == i.selectedIndex {
			isSelected = true
			break
		}
	}

	// Build skill line
	line := "  "
	if isSelected {
		line += "▶ "
	} else {
		line += "  "
	}

	line += skill.Name

	// Add level if present
	if skill.Level != "" {
		levelStyle := lipgloss.NewStyle().Foreground(theme.MutedColor())
		line += " " + levelStyle.Render(fmt.Sprintf("(%s)", skill.Level))
	}

	// Add years if present
	if skill.YearsUsed != nil {
		yearsStyle := lipgloss.NewStyle().Foreground(theme.SecondaryColor())
		yearText := fmt.Sprintf("%d yr", *skill.YearsUsed)
		if *skill.YearsUsed != 1 {
			yearText += "s"
		}
		line += " " + yearsStyle.Render(fmt.Sprintf("[%s]", yearText))
	}

	// Add event count if available
	if i.eventCounts != nil {
		if count, ok := i.eventCounts[skill.ID]; ok && count > 0 {
			countStyle := lipgloss.NewStyle().Foreground(theme.MutedColor())
			eventText := fmt.Sprintf("%d event", count)
			if count != 1 {
				eventText += "s"
			}
			line += " " + countStyle.Render(fmt.Sprintf("(%s)", eventText))
		}
	}

	// Apply selection styling
	if isSelected {
		selectedStyle := lipgloss.NewStyle().
			Foreground(theme.SuccessColor()).
			Bold(true)
		return selectedStyle.Render(line)
	}

	return line
}

func (i *ManageSkillsIntent) renderEmptyState() string {
	theme := i.Theme()

	emptyStyle := lipgloss.NewStyle().
		Foreground(theme.MutedColor()).
		Align(lipgloss.Center).
		MarginTop(3).
		MarginBottom(3)

	message := "No skills defined yet.\n\nPress 'n' to add your first skill."

	return i.getCardStyle().Render(emptyStyle.Render(message))
}

func (i *ManageSkillsIntent) groupSkillsByCategory() map[string][]*domain.Skill {
	grouped := make(map[string][]*domain.Skill)

	for _, skill := range i.skills {
		category := skill.Category
		if category == "" {
			category = "other"
		}
		grouped[category] = append(grouped[category], skill)
	}

	// Sort skills within each category by name
	for category := range grouped {
		sort.Slice(grouped[category], func(i, j int) bool {
			return strings.ToLower(grouped[category][i].Name) < strings.ToLower(grouped[category][j].Name)
		})
	}

	return grouped
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
	return nil
}

// handleDetailKeys handles key presses in detail view
func (i *ManageSkillsIntent) handleDetailKeys(msg tea.KeyMsg) tea.Cmd {
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

	case "esc":
		// Back to list
		i.currentState = SkillsStateList
		i.selectedSkill = nil
		return nil
	}

	return nil
}

// handleDetailEventsKeys handles key presses in detail events view
func (i *ManageSkillsIntent) handleDetailEventsKeys(msg tea.KeyMsg) tea.Cmd {
	switch msg.String() {
	case "esc":
		// Back to detail view
		i.currentState = SkillsStateDetail
		i.skillEvents = nil
		i.eventsLoaded = false
		return nil
	}

	return nil
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

// renderSkillEvents renders the events using this skill
func (i *ManageSkillsIntent) renderSkillEvents() string {
	theme := i.Theme()

	if !i.eventsLoaded {
		loadingStyle := lipgloss.NewStyle().
			Foreground(theme.SecondaryColor()).
			MarginTop(2)
		return i.getCardStyle().Render(loadingStyle.Render("Loading events..."))
	}

	if len(i.skillEvents) == 0 {
		emptyStyle := lipgloss.NewStyle().
			Foreground(theme.MutedColor()).
			MarginTop(2)
		return i.getCardStyle().Render(emptyStyle.Render("No events use this skill yet."))
	}

	var lines []string
	lines = append(lines, fmt.Sprintf("Events using '%s' (%d total):", i.selectedSkill.Name, len(i.skillEvents)))
	lines = append(lines, "")

	// Render each event
	for _, event := range i.skillEvents {
		dateStr := event.Date.Format("2006-01-02")
		eventText := event.Text
		if len(eventText) > 80 {
			eventText = eventText[:77] + "..."
		}

		dateStyle := lipgloss.NewStyle().Foreground(theme.SecondaryColor())
		line := dateStyle.Render(dateStr) + " " + eventText

		if event.Company != "" {
			companyStyle := lipgloss.NewStyle().Foreground(theme.MutedColor())
			line += " " + companyStyle.Render(fmt.Sprintf("(%s)", event.Company))
		}

		lines = append(lines, "  "+line)
	}

	content := lipgloss.JoinVertical(lipgloss.Left, lines...)
	return i.getCardStyle().Render(content)
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
