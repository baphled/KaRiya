// Package skillsmanagement implements the ManageSkills intent for managing user-defined skills.
package skillsmanagement

import (
	"fmt"

	domain "github.com/baphled/kariya/internal/domain/career"
	"github.com/baphled/kariya/internal/tui/intents"
	skillview "github.com/baphled/kariya/internal/tui/views/skill"
	"github.com/baphled/kariya/internal/ui/behaviors"
	"github.com/baphled/kariya/internal/ui/display"
	"github.com/baphled/kariya/internal/ui/themes"
	"github.com/baphled/kariya/internal/ui/uikit/widgets"
	tea "github.com/charmbracelet/bubbletea"
)

// Ensure Intent implements FilterBehavior interface.
var _ behaviors.FilterBehavior = (*Intent)(nil)

// transitionToView sets the active view for the intent.
func (i *Intent) transitionToView(view widgets.View) {
	i.activeView = view
}

// handleViewResult processes the result returned by a view.

// NewIntent constructs a fully configured ManageSkills intent with table behavior and modal registry.
//
// Expected: ctx must pass Validate with a non-nil Ctx and SkillRepository.
//
// Returns: the initialized intent in StateList, or an error if context validation fails.
//
// Side effects: allocates a BaseIntent, ThemeManager, TableBehavior, and ModalRegistry.
func NewIntent(ctx *IntentValidator) (*Intent, error) {
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

// GetTableBehavior returns the skills table behavior for testing and integration.
//
// Returns:
//   - A fully initialized behaviors.TableBehavior[*domain.Skill] ready for use.
//
// Side effects:
//   - None.
func (i *Intent) GetTableBehavior() *behaviors.TableBehavior[*domain.Skill] {
	return i.tableBehavior
}

// Init activates the intent, applies theming, creates the list screen, and triggers async skill loading.
//
// Returns:
//   - A tea.Cmd value.
//
// Side effects:
//   - None.
func (i *Intent) Init() tea.Cmd {
	i.active = true

	// Apply theme to TableBehavior if available.
	if theme := i.Theme(); theme != nil {
		i.tableBehavior.SetTheme(theme)
	}

	// Create and transition to SkillListView.
	// Pass skills and eventCounts (may be empty initially).
	view := skillview.NewListView(display.SkillsFromDomain(i.skills), i.eventCounts)
	i.transitionToView(view)

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

// Update processes incoming messages through a 3-tier priority system: global keys, modals, then view delegation.
//
// Expected: msg is a valid tea.Msg; the intent should be active.
//
// Returns: a tea.Cmd representing the next action, or nil if no action is needed.
//
// Side effects: may transition state, update modals, or delegate to the active view.
func (i *Intent) Update(msg tea.Msg) tea.Cmd {
	if !i.active {
		return nil
	}

	if cmd, handled := i.handleAsyncUpdates(msg); handled {
		return cmd
	}
	if cmd, handled := i.handleGlobalKeyUpdates(msg); handled {
		return cmd
	}
	if cmd, handled := i.handlePriorityModalUpdates(msg); handled {
		return cmd
	}
	if cmd, handled := i.handleFormViewModalUpdate(msg); handled {
		return cmd
	}
	if cmd, handled := i.handleIntentMessageUpdates(msg); handled {
		return cmd
	}
	return i.delegateToActiveView(msg)
}

// handleAsyncUpdates processes async result messages that bypass modal blocking.
//
// Expected: msg is a valid tea.Msg.
//
// Returns: a tea.Cmd representing the next action, or nil if not handled.
//
// Side effects: may update intent state from async results.
func (i *Intent) handleAsyncUpdates(msg tea.Msg) (tea.Cmd, bool) {
	switch msg := msg.(type) {
	case SkillSuggestionsLoadedMsg:
		i.handleSkillSuggestionsLoaded(msg)
		return nil, true
	case SkillsCreatedMsg:
		return i.handleSkillsCreatedFromInference(msg), true
	case SkillsLoadedMsg:
		i.handleSkillsLoaded(msg)
		return nil, true
	default:
		return nil, false
	}
}

// handleGlobalKeyUpdates processes global key commands.
//
// Expected: msg is a valid tea.Msg.
//
// Returns: a tea.Cmd representing the next action, or nil if not handled.
//
// Side effects: may toggle help or quit.
func (i *Intent) handleGlobalKeyUpdates(msg tea.Msg) (tea.Cmd, bool) {
	keyMsg, ok := msg.(tea.KeyMsg)
	if !ok {
		return nil, false
	}
	switch intents.HandleGlobalKeys(keyMsg) {
	case intents.KeyQuit:
		return tea.Quit, true
	case intents.KeyHelp:
		i.ToggleHelp()
		return nil, true
	default:
		return nil, false
	}
}

// handlePriorityModalUpdates processes feedback and suggestion modals in priority order.
//
// Expected: msg is a valid tea.Msg.
//
// Returns: a tea.Cmd representing the next action, or nil if not handled.
//
// Side effects: may update modal state.
func (i *Intent) handlePriorityModalUpdates(msg tea.Msg) (tea.Cmd, bool) {
	if cmd := i.handleFeedbackModalUpdate(msg); cmd != nil {
		return cmd, true
	}
	if cmd := i.handleLoadingModalUpdate(msg); cmd != nil {
		return cmd, true
	}
	if cmd := i.handleSuggestionEventsModalUpdate(msg); cmd != nil {
		return cmd, true
	}
	if cmd := i.handleSkillSuggestionModalUpdate(msg); cmd != nil {
		return cmd, true
	}
	return nil, false
}

// handleIntentMessageUpdates processes intent-specific messages.
//
// Expected: msg is a valid tea.Msg.
//
// Returns: a tea.Cmd representing the next action, or nil if not handled.
//
// Side effects: may update intent state or modals.
func (i *Intent) handleIntentMessageUpdates(msg tea.Msg) (tea.Cmd, bool) {
	switch msg := msg.(type) {
	case SkillCreatedMsg:
		return i.handleSkillCreated(msg), true
	case SkillUpdatedMsg:
		return i.handleSkillUpdated(msg), true
	case SkillDeletedMsg:
		return i.handleSkillDeleted(msg), true
	case SkillEventsLoadedMsg:
		i.handleSkillEventsLoaded(msg)
		return nil, true
	case SkillEventsForModalLoadedMsg:
		return i.handleSkillEventsForModalLoaded(msg), true
	case tea.KeyMsg:
		if cmd := i.handleKeyShortcuts(msg); cmd != nil {
			return cmd, true
		}
		return nil, false
	default:
		return nil, false
	}
}

// delegateToActiveView delegates unhandled messages to the active view.
//
// Expected: msg is a valid tea.Msg.
//
// Returns: a tea.Cmd representing the next action, or nil if no view is active.
//
// Side effects: may update the active view state.
func (i *Intent) delegateToActiveView(msg tea.Msg) tea.Cmd {
	if i.activeView == nil {
		return nil
	}
	cmd, result := i.activeView.Update(msg)
	if result != nil {
		if result.Type() == widgets.ResultCancel {
			i.SetCancelled()
			return cmd
		}
		return tea.Batch(cmd, i.handleViewResult(result))
	}
	return cmd
}

// handleFormViewModalUpdate routes messages to the appropriate form or view modal.
// Returns a tea.Cmd if a modal consumed the message, or nil if no modal is active.
func (i *Intent) handleFormViewModalUpdate(msg tea.Msg) (tea.Cmd, bool) {
	cmd, consumed := runFirstVisibleModalUpdate(msg,
		func() bool { return i.searchAdapter != nil && i.searchAdapter.IsVisible() },
		i.handleSearchModalUpdate,
		func() bool { return i.filterAdapter != nil && i.filterAdapter.IsVisible() },
		i.handleFilterModalUpdate,
		func() bool { return i.sortAdapter != nil && i.sortAdapter.IsVisible() },
		i.handleSortModalUpdate,
		func() bool { return i.viewDetailModal != nil && i.viewDetailModal.IsVisible() },
		i.handleViewDetailModalUpdate,
		func() bool { return i.addEditModal != nil && i.addEditModal.IsVisible() },
		i.handleAddEditModalUpdate,
		func() bool { return i.deleteModal != nil && i.deleteModal.IsVisible() },
		i.handleDeleteModalUpdate,
		func() bool { return i.skillEventsModal != nil && i.skillEventsModal.IsVisible() },
		i.handleSkillEventsModalUpdate,
		func() bool { return i.eventDetailModal != nil && i.eventDetailModal.IsVisible() },
		i.handleEventDetailModalUpdate,
	)
	return cmd, consumed
}

type modalVisibility struct {
	isVisible func() bool
	update    func(tea.Msg) tea.Cmd
}

// runFirstVisibleModalUpdate routes a message to the first visible modal.
//
// Returns: (cmd, consumed) where consumed is true if a modal handled the message.
func runFirstVisibleModalUpdate(msg tea.Msg, entries ...interface{}) (tea.Cmd, bool) {
	modalEntries := buildModalEntries(entries...)
	for _, entry := range modalEntries {
		if entry.isVisible() {
			return entry.update(msg), true
		}
	}
	return nil, false
}

// buildModalEntries constructs modal visibility entries from paired inputs.
//
// Expected:
//   - entries must contain alternating visibility and update funcs.
//
// Returns:
//   - A slice of modalVisibility values.
//
// Side effects:
//   - None.
func buildModalEntries(entries ...interface{}) []modalVisibility {
	if len(entries) == 0 {
		return nil
	}
	modalEntries := make([]modalVisibility, 0, len(entries)/2)
	for i := 0; i+1 < len(entries); i += 2 {
		isVisible, ok := entries[i].(func() bool)
		if !ok {
			continue
		}
		update, ok := entries[i+1].(func(tea.Msg) tea.Cmd)
		if !ok {
			continue
		}
		modalEntries = append(modalEntries, modalVisibility{isVisible: isVisible, update: update})
	}
	return modalEntries
}

// handleViewResult processes a view result with typed action dispatch (ADR pattern).
// Extracts Nav payload and switches on action type.
func (i *Intent) handleViewResult(res widgets.ViewResult) tea.Cmd {
	if res.Type() != widgets.ResultNavigate {
		return nil
	}
	payload, ok := res.Data().(skillview.Nav)
	if !ok {
		return func() tea.Msg {
			return fmt.Errorf("invalid skill action payload type: %T", res.Data())
		}
	}

	switch payload.Action {
	case skillview.ActionEdit:
		i.selectedSkill = i.findSkillByID(payload.Skill.ID)
		return i.openAddEditModal(i.selectedSkill)
	case skillview.ActionDelete:
		i.selectedSkill = i.findSkillByID(payload.Skill.ID)
		return i.openDeleteModal(i.selectedSkill)
	case skillview.ActionView:
		i.selectedSkill = i.findSkillByID(payload.Skill.ID)
		return i.openViewDetailModal()
	case skillview.ActionAdd:
		return i.openAddEditModal(nil)
	case skillview.ActionFilter:
		return i.openFilterModal()
	case skillview.ActionSearch:
		return i.openSearchModal()
	case skillview.ActionSort:
		return i.openSortModal()
	case skillview.ActionInfer:
		return i.startSkillInference()
	default:
		return func() tea.Msg {
			return fmt.Errorf("unknown skill action: %s", payload.Action)
		}
	}
}

// View renders the active view with any visible modal overlays applied.
//
// Returns:
//   - A string value.
//
// Side effects:
//   - None.
func (i *Intent) View() string {
	if !i.active {
		return ""
	}
	if i.activeView == nil {
		return "No active view"
	}
	view := i.CreateViewWithBreadcrumbs(i.getBreadcrumbs()...)
	view.WithContent(i.activeView.RenderContent())
	view.WithHelp(i.activeView.HelpText()).WithFooterSeparator(true)
	baseView := view.Render()
	i.rebuildModalRegistry()
	return i.modalRegistry.RenderOverlay(baseView)
}

// Result provides the final outcome of the intent for the caller to inspect after completion.
//
// Returns:
//   - A fully initialized intents.IntentResult[interface{}] ready for use.
//
// Side effects:
//   - None.
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
// Side effects:
//   - None.
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
//
// Returns:
//   - A fully initialized IntentValidator ready for use.
//
// Side effects:
//   - None.
func (i *Intent) GetTestContext() *IntentValidator {
	return i.context
}

// FilterBehavior interface implementation.

// HasActiveFilters checks whether the current filter state differs from the default, enabling UI indicators.
//
// Returns:
//   - A bool value.
//
// Side effects:
//   - None.
func (i *Intent) HasActiveFilters() bool {
	if i.context == nil || i.context.Filters == nil {
		return false
	}
	return i.context.Filters.HasActiveFilters()
}

// ClearFilters progressively resets the most recently applied filter in FIFO order.
//
// Side effects:
//   - None.
func (i *Intent) ClearFilters() {
	if i.context != nil && i.context.Filters != nil {
		i.context.Filters.Clear()
	}
}

// ApplyFilters satisfies the FilterBehavior interface; actual filtering is handled during skill loading.
//
// Side effects:
//   - None.
func (i *Intent) ApplyFilters() {
	// Filters are applied when loading skills.
}

// RefreshData triggers an asynchronous reload of the skills list with current filters applied.
//
// Returns:
//   - A tea.Cmd value.
//
// Side effects:
//   - None.
func (i *Intent) RefreshData() tea.Cmd {
	return i.reloadSkills()
}

// ViewResult handler interface implementation.

// HandleCancel processes a view cancellation by returning the intent to the list state.
//
// Expected: result is a valid ViewResult from an active view.
//
// Returns: nil; no follow-up command is needed.
//
// Side effects: transitions the intent state back to StateList.
func (i *Intent) HandleCancel(_ widgets.ViewResult) tea.Cmd {
	i.state = StateList
	return nil
}

// HandleNavigate dispatches navigation actions such as view, edit, and delete from view results.
//
// Expected: result is a NavigateViewResult with ResultData castable to map[string]interface{}.
//
// Returns: a tea.Cmd for the dispatched action, or nil if the data format is unexpected.
//
// Side effects: may open modals or trigger state transitions depending on the navigation action.
func (i *Intent) HandleNavigate(result widgets.ViewResult) tea.Cmd {
	nav, ok := result.(*widgets.NavigateViewResult)
	if !ok {
		return nil
	}
	if actionData, ok := nav.ResultData.(map[string]interface{}); ok {
		return i.handleNavigateData(actionData)
	}
	return nil
}

// HandleSubmit processes submission results from views; currently a no-op for this intent.
//
// Expected: result is a valid ViewResult from an active view.
//
// Returns: nil; no follow-up command is needed.
//
// Side effects: None.
func (i *Intent) HandleSubmit(_ widgets.ViewResult) tea.Cmd {
	return nil
}

// HandleError processes error results from views; currently a no-op for this intent.
//
// Expected: result is a valid ViewResult from an active view.
//
// Returns: nil; no follow-up command is needed.
//
// Side effects: None.
func (i *Intent) HandleError(_ widgets.ViewResult) tea.Cmd {
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
