// Package fact_management implements the FactManagement intent for managing career facts.
package fact_management

import (
	"fmt"

	"github.com/baphled/kariya/internal/cli/behaviors"
	"github.com/baphled/kariya/internal/cli/intents"
	"github.com/baphled/kariya/internal/cli/screens"
	"github.com/baphled/kariya/internal/cli/screens/fact"
	domain "github.com/baphled/kariya/internal/domain/career"
	tea "github.com/charmbracelet/bubbletea"
)

// NewIntent creates a new FactManagement intent.
func NewIntent(ctx *IntentContext) (*Intent, error) {
	if err := ctx.Validate(); err != nil {
		return nil, err
	}

	// Create column definitions for TableBehavior.
	columns := []behaviors.ColumnDef{
		{Title: "Fact", Width: 50},
		{Title: "Strength", Width: 15},
		{Title: "Categories", Width: 30},
	}

	// Create TableBehavior with type-safe generic.
	tableBehavior := behaviors.NewTableBehavior[*domain.Fact](nil, columns, factRowFormatter).
		PageSize(15).
		PaginationPrefix("Facts").
		EmptyMessage("No facts found. Press 'n' to create a new fact, 'r' to refresh, or 'q' to quit.")

	intent := &Intent{
		BaseIntent:    intents.NewBaseIntent(),
		context:       ctx,
		state:         StateList,
		active:        true,
		tableBehavior: tableBehavior,
		modalRegistry: intents.NewModalRegistry(),
	}

	return intent, nil
}

// factRowFormatter formats a Fact for display in the table.
func factRowFormatter(fact *domain.Fact, _ int) []string {
	text := truncate(fact.Text, 50)
	strength := fact.StrengthSignal
	if strength == "" {
		strength = "-"
	}
	categories := fmt.Sprintf("%v", fact.CompetencyCategories)
	if len(categories) > 30 {
		categories = categories[:27] + "..."
	}
	return []string{text, strength, categories}
}

// truncate shortens a string to the specified length.
func truncate(s string, maxLen int) string {
	if len(s) <= maxLen {
		return s
	}
	return s[:maxLen] + "..."
}

// Init initializes the intent.
func (i *Intent) Init() tea.Cmd {
	// Apply theme to TableBehavior if available.
	if theme := i.Theme(); theme != nil {
		i.tableBehavior.SetTheme(theme)
	}

	// Load facts from repository.
	if err := i.context.LoadFacts(); err != nil {
		i.result = &intents.IntentResult[*Result]{
			Status: intents.Failed,
			Error: &intents.IntentError{
				Code:    "LOAD_FACTS_FAILED",
				Message: "Failed to load facts",
				Cause:   err,
			},
		}
		return tea.Quit
	}

	// Set items in TableBehavior.
	i.tableBehavior.SetItems(i.context.Facts)

	// Create the list screen.
	i.listScreen = fact.NewListScreen(i.tableBehavior)
	i.transitionToScreen(i.listScreen)

	return nil
}

// Update processes messages.
func (i *Intent) Update(msg tea.Msg) tea.Cmd {
	if !i.active {
		return nil
	}

	// Handle modals first.
	if cmd := i.handleModalUpdates(msg); cmd != nil || i.hasActiveModal() {
		return cmd
	}

	// Handle state-specific updates.
	switch i.state {
	case StateList:
		return i.handleListState(msg)
	case StateView:
		return i.handleViewState(msg)
	case StateEditor:
		return i.handleEditorState(msg)
	case StateDeleteConfirm:
		return i.handleDeleteConfirmState(msg)
	case StateResults:
		return i.handleResultsState(msg)
	case StateCompleted:
		return nil
	}

	return nil
}

// View renders the intent.
func (i *Intent) View() string {
	if !i.active {
		return "FactManagement intent is not active"
	}

	// Create standard view with dynamic breadcrumbs.
	view := intents.CreateStandardViewWithBreadcrumbs(i.BaseIntent, i.getBreadcrumbs()...)

	// Handle form errors.
	if i.context.HasFormErrors() {
		var errorMessages []string
		for field, err := range i.context.FormErrors {
			errorMessages = append(errorMessages, fmt.Sprintf("%s: %s", field, err))
		}
		i.SetError(fmt.Errorf("validation errors:\n%s", fmt.Sprintf("%v", errorMessages)))
	}

	// Get content for current state.
	content := i.getStateContent()
	view.WithContent(content)

	// Get context-aware help.
	help := i.getContextHelp()
	view.WithHelp(help).WithFooterSeparator(true)

	baseView := view.Render()

	// Render modal overlay if any modal is visible.
	i.rebuildModalRegistry()
	return i.modalRegistry.RenderOverlay(baseView)
}

// Result returns the intent result.
func (i *Intent) Result() *intents.IntentResult[interface{}] {
	if i.result == nil {
		return nil
	}
	return &intents.IntentResult[interface{}]{
		Status: i.result.Status,
		Error:  i.result.Error,
	}
}

// ScreenResultHandler implementation.

// HandleCancel handles cancel results from screens.
func (i *Intent) HandleCancel(result *screens.CancelResult) tea.Cmd {
	i.setCancelled()
	return nil
}

// HandleNavigate handles navigation results from screens.
func (i *Intent) HandleNavigate(result *screens.NavigateResult) tea.Cmd {
	return nil
}

// HandleSubmit handles submit results from screens.
func (i *Intent) HandleSubmit(result *screens.SubmitResult) tea.Cmd {
	return nil
}

// HandleError handles error results from screens.
func (i *Intent) HandleError(result *screens.ErrorResult) tea.Cmd {
	i.showErrorModal("Error", result.Err.Error())
	return nil
}

// ListNavigator interface implementation.

// GetTotalItems returns the total number of items.
func (i *Intent) GetTotalItems() int {
	return len(i.context.Facts)
}

// GetSelectedIndex returns the currently selected index.
func (i *Intent) GetSelectedIndex() int {
	return i.context.SelectedFactIndex
}

// SetSelectedIndex sets the selected index.
func (i *Intent) SetSelectedIndex(idx int) {
	i.tableBehavior.SetSelectedIndex(idx)
	i.syncTableSelection()
}

// GetPageSize returns the page size.
func (i *Intent) GetPageSize() int {
	return 15
}
