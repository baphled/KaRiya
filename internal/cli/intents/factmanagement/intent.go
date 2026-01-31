// Package factmanagement implements the FactManagement intent for managing career facts.
package factmanagement

import (
	"fmt"

	"github.com/baphled/kariya/internal/cli/behaviors"
	"github.com/baphled/kariya/internal/cli/intents"
	domain "github.com/baphled/kariya/internal/domain/career"
	tea "github.com/charmbracelet/bubbletea"
)

// NewIntent constructs a fully wired FactManagement intent with table behavior and default state.
//
// Expected: ctx must be non-nil and pass Validate. Passing a nil ctx will cause a nil-pointer panic.
//
// Returns: the initialised intent and nil error on success, or nil and an error if context validation fails.
//
// Side effects: None.
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
		active:        false,
		tableBehavior: tableBehavior,
	}

	return intent, nil
}

// truncate shortens a string to the specified length.
func truncate(s string, maxLen int) string {
	if len(s) <= maxLen {
		return s
	}
	return s[:maxLen] + "..."
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

// Init activates the intent, applies theming, and loads facts from the repository for initial display.
//
// Returns: tea.Quit if fact loading fails, or nil on success.
//
// Side effects: Sets active to true, applies theme to table behavior,
// loads facts into context, and populates the table. On failure,
// populates the result with an error status.
func (i *Intent) Init() tea.Cmd {
	// Apply theme to TableBehavior if available.
	if theme := i.Theme(); theme != nil {
		i.tableBehavior.SetTheme(theme)
	}

	// Mark intent as active.
	i.active = true

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

	return nil
}

// Update dispatches incoming messages to the handler for the current intent state.
//
// Expected: msg must be a valid tea.Msg. The intent should be active before receiving updates.
//
// Returns: a tea.Cmd if the handler produces one, or nil.
//
// Side effects: Delegates to state-specific handlers which may mutate intent and context state.
func (i *Intent) Update(msg tea.Msg) tea.Cmd {
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
		// Do not send tea.Quit, just return nil and let router manage exit.
		return nil
	}
	return nil
}

// View composes the full visual output for the current intent state, including breadcrumbs, content, and help.
//
// Returns: the rendered string for terminal display, or an inactive message if the intent has not been initialised.
//
// Side effects: May set an error on BaseIntent when form validation errors are present.
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

	return view.Render()
}

// Result provides the outcome of the intent lifecycle for the router to inspect after completion.
//
// Returns: the intent result with status and optional error, or nil if the intent has not completed.
//
// Side effects: None.
func (i *Intent) Result() *intents.IntentResult[interface{}] {
	if i.result == nil {
		return nil
	}
	return &intents.IntentResult[interface{}]{
		Status: i.result.Status,
		Error:  i.result.Error,
	}
}

// ListNavigator interface implementation.

// GetTotalItems reports the size of the full fact collection for pagination calculations.
//
// Returns: the number of facts currently held in the context.
//
// Side effects: None.
func (i *Intent) GetTotalItems() int {
	return len(i.context.Facts)
}

// GetSelectedIndex provides the zero-based index of the currently highlighted fact in the list.
//
// Returns: the selected index, or -1 if no fact is selected.
//
// Side effects: None.
func (i *Intent) GetSelectedIndex() int {
	return i.context.SelectedFactIndex
}

// SetSelectedIndex moves the selection cursor to the given index, synchronising both table behavior and context.
//
// Expected: idx must be a valid zero-based index within the total items range.
//
// Side effects: Updates the table behavior selection and syncs the context selection state.
func (i *Intent) SetSelectedIndex(idx int) {
	i.tableBehavior.SetSelectedIndex(idx)
	i.syncTableSelection()
}

// GetPageSize provides the fixed number of facts displayed per page for the ListNavigator contract.
//
// Returns: the page size (15).
//
// Side effects: None.
func (i *Intent) GetPageSize() int {
	return 15
}
