// Package skillsmanagement implements the ManageSkills intent for managing user-defined skills.
package skillsmanagement

import (
	"github.com/baphled/kariya/internal/cli/behaviors"
	"github.com/baphled/kariya/internal/cli/intents"
	"github.com/baphled/kariya/internal/cli/screens"
	"github.com/baphled/kariya/internal/cli/screens/skills"
	"github.com/baphled/kariya/internal/cli/themes"
	domain "github.com/baphled/kariya/internal/domain/career"
	tea "github.com/charmbracelet/bubbletea"
)

// Ensure Intent implements FilterBehavior interface.
var _ behaviors.FilterBehavior = (*Intent)(nil)

// Ensure Intent implements ScreenResultHandler interface.
var _ behaviors.ScreenResultHandler = (*Intent)(nil)

// NewIntent creates a new ManageSkills intent.
func NewIntent(ctx *IntentContext) (*Intent, error) {
	if err := ctx.Validate(); err != nil {
		return nil, err
	}

	baseIntent := intents.NewBaseIntent()
	baseIntent.SetThemeManager(themes.NewThemeManager())

	// Create column definitions for skills TableBehavior.
	skillsColumns := []behaviors.ColumnDef{
		{Title: "Name", Width: 25},
		{Title: "Category", Width: 15},
		{Title: "Level", Width: 12},
		{Title: "Years", Width: 8},
		{Title: "Events", Width: 8},
	}

	// Create TableBehavior for skills list.
	tableBehavior := behaviors.NewTableBehavior[*domain.Skill](nil, skillsColumns, skillRowFormatterWithCounts(nil)).
		PageSize(15).
		PaginationPrefix("Skills").
		EmptyMessage("No skills found. Press 'a' to add a new skill.")

	intent := &Intent{
		BaseIntent:    baseIntent,
		context:       ctx,
		state:         StateList,
		skills:        []*domain.Skill{},
		active:        false,
		tableBehavior: tableBehavior,
		modalRegistry: intents.NewModalRegistry(),
	}

	return intent, nil
}

// GetTableBehavior returns the table behavior for skills.
func (i *Intent) GetTableBehavior() *behaviors.TableBehavior[*domain.Skill] {
	return i.tableBehavior
}

// Init initializes the intent and loads skills.
func (i *Intent) Init() tea.Cmd {
	i.active = true

	// Apply theme to TableBehavior if available.
	if theme := i.Theme(); theme != nil {
		i.tableBehavior.SetTheme(theme)
	}

	// Create and transition to the list screen.
	i.listScreen = skills.NewSkillsListScreen(i.skills)
	i.transitionToScreen(i.listScreen)

	// Load skills asynchronously.
	return func() tea.Msg {
		if i.context == nil || i.context.SkillRepository == nil {
			return SkillsLoadedMsg{
				Skills: []*domain.Skill{},
				Error:  ErrRepositoryNotAvailable,
			}
		}
		loadedSkills, err := i.context.LoadSkills()
		return SkillsLoadedMsg{
			Skills: loadedSkills,
			Error:  err,
		}
	}
}

// Update handles messages and state transitions.
func (i *Intent) Update(msg tea.Msg) tea.Cmd {
	if !i.active {
		return nil
	}

	// PATTERN 4: Global Key Interception - 3-tier priority.
	// 1. Check global keys FIRST (work everywhere, even in modals).
	if keyMsg, ok := msg.(tea.KeyMsg); ok {
		switch intents.HandleGlobalKeys(keyMsg) {
		case intents.KeyQuit:
			return tea.Quit
		case intents.KeyHelp:
			i.ToggleHelp()
			return nil
		}
	}

	// 2. SECOND PRIORITY: Modal updates (if visible).
	// Pass full tea.Msg (not tea.KeyMsg) to modals for huh forms to work.
	if i.searchModal != nil && i.searchModal.IsVisible() {
		return i.handleSearchModalUpdate(msg)
	}
	if i.filterModal != nil && i.filterModal.IsVisible() {
		return i.handleFilterModalUpdate(msg)
	}
	if i.sortModal != nil && i.sortModal.IsVisible() {
		return i.handleSortModalUpdate(msg)
	}
	if i.viewDetailModal != nil && i.viewDetailModal.IsVisible() {
		return i.handleViewDetailModalUpdate(msg)
	}
	if i.addEditModal != nil && i.addEditModal.IsVisible() {
		return i.handleAddEditModalUpdate(msg)
	}
	if i.deleteModal != nil && i.deleteModal.IsVisible() {
		return i.handleDeleteModalUpdate(msg)
	}
	if i.skillEventsModal != nil && i.skillEventsModal.IsVisible() {
		return i.handleSkillEventsModalUpdate(msg)
	}
	if i.eventDetailModal != nil && i.eventDetailModal.IsVisible() {
		return i.handleEventDetailModalUpdate(msg)
	}

	// 3. THIRD PRIORITY: Message and state handling.
	switch msg := msg.(type) {
	case SkillsLoadedMsg:
		return i.handleSkillsLoaded(msg)

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
		// Handle key shortcuts at intent level (filter/sort/search modal openers).
		if cmd := i.handleKeyShortcuts(msg); cmd != nil {
			return cmd
		}
	}

	// 4. Delegate to active screen for all other messages.
	if i.activeScreen != nil {
		cmd, result := i.activeScreen.Update(msg)
		if result != nil {
			return tea.Batch(cmd, i.handleScreenResult(result))
		}
		return cmd
	}

	return nil
}

// View renders the current state.
func (i *Intent) View() string {
	if !i.active {
		return ""
	}

	if i.activeScreen == nil {
		return "No active screen"
	}

	// Render based on screen type.
	switch screen := i.activeScreen.(type) {
	case *skills.SkillsListScreen:
		return i.renderListView(screen)
	default:
		return i.activeScreen.View()
	}
}

// renderListView renders the skills list view with modal overlays.
func (i *Intent) renderListView(screen *skills.SkillsListScreen) string {
	// Create standard view with breadcrumbs.
	view := i.CreateViewWithBreadcrumbs(i.getBreadcrumbs()...)

	// Get content from the screen.
	view.WithContent(screen.RenderContent())

	// Add context-aware help.
	help := i.getContextHelp()
	view.WithHelp(help).WithFooterSeparator(true)

	// PATTERN 1: Modal Overlay Rendering.
	// StandardView FIRST, modal overlay LAST (prevents misalignment).
	baseView := view.Render()

	// Rebuild registry and render any visible modal as overlay.
	i.rebuildModalRegistry()
	return i.modalRegistry.RenderOverlay(baseView)
}

// Result returns the final result of the intent.
func (i *Intent) Result() *intents.IntentResult[interface{}] {
	if i.result == nil {
		return nil
	}
	return &intents.IntentResult[interface{}]{
		Status:   i.result.Status,
		Data:     i.result.Data,
		Error:    i.result.Error,
		Metadata: i.result.Metadata,
	}
}

// SetCancelled marks the intent as cancelled.
func (i *Intent) SetCancelled() {
	i.result = &intents.IntentResult[*Result]{
		Status: intents.Cancelled,
		Data: &Result{
			Action: "cancelled",
		},
	}
	i.active = false
}

// FilterBehavior interface implementation.

// HasActiveFilters returns true if any non-default filters are active.
func (i *Intent) HasActiveFilters() bool {
	if i.context == nil || i.context.Filters == nil {
		return false
	}
	return i.context.Filters.HasActiveFilters()
}

// ClearFilters resets filters in FIFO order.
func (i *Intent) ClearFilters() {
	if i.context != nil && i.context.Filters != nil {
		i.context.Filters.Clear()
	}
}

// ApplyFilters applies current filter state to the data.
func (i *Intent) ApplyFilters() {
	// Filters are applied when loading skills.
}

// RefreshData reloads/refreshes the filtered data.
func (i *Intent) RefreshData() tea.Cmd {
	return i.reloadSkills()
}

// ScreenResultHandler interface implementation.

// HandleCancel handles cancel results from screens.
func (i *Intent) HandleCancel(_ *screens.CancelResult) tea.Cmd {
	i.state = StateList
	return nil
}

// HandleNavigate handles navigation results from screens.
func (i *Intent) HandleNavigate(result *screens.NavigateResult) tea.Cmd {
	if actionData, ok := result.ResultData.(map[string]interface{}); ok {
		return i.handleNavigateData(actionData)
	}
	return nil
}

// HandleSubmit handles submit results from screens.
func (i *Intent) HandleSubmit(_ *screens.SubmitResult) tea.Cmd {
	return nil
}

// HandleError handles error results from screens.
func (i *Intent) HandleError(_ *screens.ErrorResult) tea.Cmd {
	return nil
}

// reloadSkills returns a command to reload skills.
func (i *Intent) reloadSkills() tea.Cmd {
	return func() tea.Msg {
		skills, err := i.context.LoadSkills()
		return SkillsLoadedMsg{
			Skills: skills,
			Error:  err,
		}
	}
}

// transitionToScreen sets the active screen and updates state.
func (i *Intent) transitionToScreen(screen screens.Screen) {
	i.activeScreen = screen

	termInfo := i.GetTerminalInfo()
	if termInfo != nil && termInfo.Width > 0 && termInfo.Height > 0 {
		screen.SetTerminalInfo(termInfo.Width, termInfo.Height)
	}

	if theme := i.Theme(); theme != nil {
		screen.SetTheme(theme)
	}

	if logo := i.GetLogo(); logo != nil {
		screen.SetLogo(logo, i.GetLogoSpacing())
	}
}
