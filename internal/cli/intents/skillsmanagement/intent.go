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

// NewIntent constructs a fully configured ManageSkills intent with table behavior and modal registry.
//
// Expected: ctx must pass Validate with a non-nil Ctx and SkillRepository.
//
// Returns: the initialized intent in StateList, or an error if context validation fails.
//
// Side effects: allocates a BaseIntent, ThemeManager, TableBehavior, and ModalRegistry.
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

// GetTableBehavior provides access to the underlying table behavior for testing and screen integration.
//
// Returns: the TableBehavior instance managing the skills list.
//
// Side effects: None.
func (i *Intent) GetTableBehavior() *behaviors.TableBehavior[*domain.Skill] {
	return i.tableBehavior
}

// Init activates the intent, applies theming, creates the list screen, and triggers async skill loading.
//
// Returns: a tea.Cmd that asynchronously loads skills from the repository.
//
// Side effects: sets the intent to active, applies theme to TableBehavior, creates and transitions to the list screen.
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

// Update processes incoming messages through a 3-tier priority system: global keys, modals, then screen delegation.
//
// Expected: msg is a valid tea.Msg; the intent should be active.
//
// Returns: a tea.Cmd representing the next action, or nil if no action is needed.
//
// Side effects: may transition state, update modals, or delegate to the active screen.
func (i *Intent) Update(msg tea.Msg) tea.Cmd {
	if !i.active {
		return nil
	}

	// 1. Async result messages FIRST (must bypass modal blocking).
	switch msg := msg.(type) {
	case SkillSuggestionsLoadedMsg:
		return i.handleSkillSuggestionsLoaded(msg)
	case SkillsCreatedMsg:
		return i.handleSkillsCreatedFromInference(msg)
	case SkillsLoadedMsg:
		return i.handleSkillsLoaded(msg)
	default:
	}

	// 2. Global keys (work everywhere, even in modals).
	if keyMsg, ok := msg.(tea.KeyMsg); ok {
		switch intents.HandleGlobalKeys(keyMsg) {
		case intents.KeyQuit:
			return tea.Quit
		case intents.KeyHelp:
			i.ToggleHelp()
			return nil
		}
	}

	// 3. Loading/feedback modal updates (highest modal priority).
	if cmd := i.handleFeedbackModalUpdate(msg); cmd != nil {
		return cmd
	}
	if cmd := i.handleLoadingModalUpdate(msg); cmd != nil {
		return cmd
	}
	if cmd := i.handleSuggestionEventsModalUpdate(msg); cmd != nil {
		return cmd
	}
	if cmd := i.handleSkillSuggestionModalUpdate(msg); cmd != nil {
		return cmd
	}

	// 4. Form/view modal updates (if visible).
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

	// 5. Message and state handling.
	switch msg := msg.(type) {
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

// View renders the active screen with any visible modal overlays applied.
//
// Returns: the rendered string for the current intent state, or an empty string if inactive.
//
// Side effects: None.
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

// Result provides the final outcome of the intent for the caller to inspect after completion.
//
// Returns: the intent result with status and data, or nil if the intent has not completed.
//
// Side effects: None.
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

// SetCancelled terminates the intent with a cancelled status, preventing further updates.
//
// Side effects: sets the result to Cancelled status and deactivates the intent.
func (i *Intent) SetCancelled() {
	i.result = &intents.IntentResult[*Result]{
		Status: intents.Cancelled,
		Data: &Result{
			Action: "cancelled",
		},
	}
	i.active = false
}

// GetTestContext returns the intent context for testing purposes.
// This method is only for testing and should not be used in production code.
//
// Returns:
//   - The IntentContext instance.
//
// Side effects:
//   - None.
func (i *Intent) GetTestContext() *IntentContext {
	return i.context
}

// FilterBehavior interface implementation.

// HasActiveFilters checks whether the current filter state differs from the default, enabling UI indicators.
//
// Returns: true if any category, level, event count, search, or sort filter is set.
//
// Side effects: None.
func (i *Intent) HasActiveFilters() bool {
	if i.context == nil || i.context.Filters == nil {
		return false
	}
	return i.context.Filters.HasActiveFilters()
}

// ClearFilters progressively resets the most recently applied filter in FIFO order.
//
// Side effects: modifies the context Filters, clearing search first, then category/level, then sort.
func (i *Intent) ClearFilters() {
	if i.context != nil && i.context.Filters != nil {
		i.context.Filters.Clear()
	}
}

// ApplyFilters satisfies the FilterBehavior interface; actual filtering is handled during skill loading.
//
// Side effects: None.
func (i *Intent) ApplyFilters() {
	// Filters are applied when loading skills.
}

// RefreshData triggers an asynchronous reload of the skills list with current filters applied.
//
// Returns: a tea.Cmd that fetches skills from the repository.
//
// Side effects: initiates an async repository call that will produce a SkillsLoadedMsg.
func (i *Intent) RefreshData() tea.Cmd {
	return i.reloadSkills()
}

// ScreenResultHandler interface implementation.

// HandleCancel processes a screen cancellation by returning the intent to the list state.
//
// Expected: result is a valid CancelResult from an active screen.
//
// Returns: nil; no follow-up command is needed.
//
// Side effects: transitions the intent state back to StateList.
func (i *Intent) HandleCancel(_ *screens.CancelResult) tea.Cmd {
	i.state = StateList
	return nil
}

// HandleNavigate dispatches navigation actions such as view, edit, and delete from screen results.
//
// Expected: result contains ResultData castable to map[string]interface{} with action-specific keys.
//
// Returns: a tea.Cmd for the dispatched action, or nil if the data format is unexpected.
//
// Side effects: may open modals or trigger state transitions depending on the navigation action.
func (i *Intent) HandleNavigate(result *screens.NavigateResult) tea.Cmd {
	if actionData, ok := result.ResultData.(map[string]interface{}); ok {
		return i.handleNavigateData(actionData)
	}
	return nil
}

// HandleSubmit processes submission results from screens; currently a no-op for this intent.
//
// Expected: result is a valid SubmitResult from an active screen.
//
// Returns: nil; no follow-up command is needed.
//
// Side effects: None.
func (i *Intent) HandleSubmit(_ *screens.SubmitResult) tea.Cmd {
	return nil
}

// HandleError processes error results from screens; currently a no-op for this intent.
//
// Expected: result is a valid ErrorResult from an active screen.
//
// Returns: nil; no follow-up command is needed.
//
// Side effects: None.
func (i *Intent) HandleError(_ *screens.ErrorResult) tea.Cmd {
	return nil
}

// reloadSkills returns a command to reload skills.
func (i *Intent) reloadSkills() tea.Cmd {
	return func() tea.Msg {
		skillList, err := i.context.LoadSkills()
		return SkillsLoadedMsg{
			Skills: skillList,
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
