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

	// active indicates whether this intent is currently active
	active bool

	// result is the final result of the intent
	result *IntentResult[*ManageSkillsResult]
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
		return CombineThemedFooters(
			ThemedListFooter(theme),
			ThemedCustomFooter(theme,
				components.NewKeyBadge("Enter", "View details"),
				components.NewKeyBadge("n", "New skill"),
			),
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
