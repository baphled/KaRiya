package intents

import (
	"fmt"
	"sort"
	"strings"

	"github.com/baphled/kariya/internal/cli/components"
	"github.com/baphled/kariya/internal/cli/forms"
	"github.com/baphled/kariya/internal/cli/themes"
	domain "github.com/baphled/kariya/internal/domain/career"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/huh"
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

	// form for add/edit
	form *huh.Form

	// active indicates whether this intent is currently active
	active bool

	// result is the final result of the intent
	result *IntentResult[*ManageSkillsResult]
}

// NewManageSkillsIntent creates a new ManageSkills intent
func NewManageSkillsIntent(ctx *ManageSkillsContext) *ManageSkillsIntent {
	baseIntent := NewBaseIntent()
	baseIntent.SetThemeManager(themes.NewThemeManager())

	// Initialize logo
	logo := components.NewASCIILogo(false, 80)
	baseIntent.SetLogo(logo)

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

	case SkillFormCompleteMsg:
		return i.handleFormComplete(msg)

	case SkillCreatedMsg:
		return i.handleSkillCreated(msg)

	case SkillUpdatedMsg:
		return i.handleSkillUpdated(msg)

	case SkillDeletedMsg:
		return i.handleSkillDeleted(msg)

	case tea.KeyMsg:
		// If we have a form active, handle it specially
		if i.form != nil {
			// Check for Esc key to cancel form
			if msg.Type == tea.KeyEsc {
				i.handleFormCancel()
				return nil
			}

			// Update form with key messages
			form, cmd := i.form.Update(msg)
			if f, ok := form.(*huh.Form); ok {
				i.form = f

				// Check if form completed
				if i.form.State == huh.StateCompleted {
					return i.handleFormSubmit()
				} else if i.form.State == huh.StateAborted {
					return i.handleFormCancel()
				}
			}
			return cmd
		}
		return i.handleKeyPress(msg)

	case tea.WindowSizeMsg:
		// BaseIntent doesn't have UpdateTerminalSize, just store in terminal info
		return nil
	}

	// If we have a form active, let it handle other messages
	if i.form != nil {
		form, cmd := i.form.Update(msg)
		if f, ok := form.(*huh.Form); ok {
			i.form = f

			// Check if form completed
			if i.form.State == huh.StateCompleted {
				return i.handleFormSubmit()
			} else if i.form.State == huh.StateAborted {
				return i.handleFormCancel()
			}
		}
		return cmd
	}

	return nil
}

// View renders the current state
func (i *ManageSkillsIntent) View() string {
	if !i.active {
		return ""
	}

	switch i.currentState {
	case SkillsStateList:
		return i.viewList()
	case SkillsStateAdd, SkillsStateEdit:
		return i.viewForm()
	case SkillsStateDelete:
		return i.viewDeleteConfirm()
	default:
		return "Unknown state"
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
	return nil
}

func (i *ManageSkillsIntent) handleFormComplete(msg SkillFormCompleteMsg) tea.Cmd {
	if msg.Cancelled {
		i.currentState = SkillsStateList
		i.form = nil
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
	i.form = nil

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
	i.form = nil

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

	case "n":
		// Add new skill
		i.currentState = SkillsStateAdd
		i.form = forms.NewSkillForm(nil, nil)
		return i.form.Init()

	case "e":
		// Edit selected skill
		if len(i.skills) == 0 {
			return nil
		}
		i.currentState = SkillsStateEdit
		i.form = forms.NewSkillForm(i.skills[i.selectedIndex], nil)
		return i.form.Init()

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

func (i *ManageSkillsIntent) handleFormSubmit() tea.Cmd {
	// Extract form data
	var skill *domain.Skill

	if i.currentState == SkillsStateEdit {
		// Editing existing skill
		skill = i.skills[i.selectedIndex]
	} else {
		// Creating new skill
		skill = &domain.Skill{}
	}

	// Apply form data
	formData := &forms.SkillFormData{
		Name:      i.form.GetString("name"),
		Category:  i.form.GetString("category"),
		Level:     i.form.GetString("level"),
		YearsUsed: i.form.GetString("years"),
	}

	err := forms.ApplySkillFormData(skill, formData)
	if err != nil {
		// Stay in form state with error
		return nil
	}

	// Save skill based on state
	if i.currentState == SkillsStateAdd {
		return i.createSkill(skill)
	} else if i.currentState == SkillsStateEdit {
		return i.updateSkill(skill)
	}

	return nil
}

func (i *ManageSkillsIntent) handleFormCancel() tea.Cmd {
	i.currentState = SkillsStateList
	i.form = nil
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

func (i *ManageSkillsIntent) viewList() string {
	content := i.renderSkillsList()
	help := i.renderListHelp()

	view := CreateStandardViewWithBreadcrumbs(i.BaseIntent, "Skills")
	view.WithContent(content)
	view.WithHelp(help)
	return view.Render()
}

func (i *ManageSkillsIntent) viewForm() string {
	if i.form == nil {
		return "Form not initialized"
	}

	title := "Add Skill"
	if i.currentState == SkillsStateEdit {
		title = "Edit Skill"
	}

	content := i.form.View()

	view := CreateStandardViewWithBreadcrumbs(i.BaseIntent, "Skills", title)
	view.WithContent(content)
	return view.Render()
}

func (i *ManageSkillsIntent) viewDeleteConfirm() string {
	if len(i.skills) == 0 {
		return "No skill selected"
	}

	skill := i.skills[i.selectedIndex]

	modalContent := fmt.Sprintf("Are you sure you want to delete '%s'?\n\nThis action cannot be undone.\n\ny:confirm • n/Esc:cancel", skill.Name)
	modal := components.NewWarningModal("Delete Skill", modalContent)

	view := CreateStandardViewWithBreadcrumbs(i.BaseIntent, "Skills", "Delete")
	view.WithContent(i.renderSkillsList())
	view.ShowModalOverlay(modal)
	return view.Render()
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

	return lipgloss.JoinVertical(lipgloss.Left, sections...)
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

	return emptyStyle.Render(message)
}

func (i *ManageSkillsIntent) renderListHelp() string {
	return "j/k:navigate • n:add • e:edit • d:delete • Esc:back"
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
