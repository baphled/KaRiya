package intents

import (
	"fmt"

	"github.com/baphled/kariya/internal/cli/behaviors"
	"github.com/baphled/kariya/internal/cli/uikit/primitives"
	domain "github.com/baphled/kariya/internal/domain/career"
	tea "github.com/charmbracelet/bubbletea"
)

type FactManagementModel struct {
	*BaseIntent
	data          *FactManagementContext
	tableBehavior *behaviors.TableBehavior[*domain.Fact]
	result        *IntentResult[*FactManagementResult]
	active        bool
	editModal     *EditFactModal
}

// factRowFormatter formats a Fact for display in the table
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

func NewFactManagementIntent(data *FactManagementContext) *FactManagementModel {
	// Create column definitions for TableBehavior
	columns := []behaviors.ColumnDef{
		{Title: "Fact", Width: 50},
		{Title: "Strength", Width: 15},
		{Title: "Categories", Width: 30},
	}

	// Create TableBehavior with type-safe generic
	tableBehavior := behaviors.NewTableBehavior[*domain.Fact](nil, columns, factRowFormatter).
		PageSize(15).
		PaginationPrefix("Facts").
		EmptyMessage("No facts found. Press 'n' to create a new fact, 'r' to refresh, or 'q' to quit.")

	model := &FactManagementModel{
		BaseIntent:    NewBaseIntent(),
		data:          data,
		tableBehavior: tableBehavior,
		result:        nil,
		active:        false,
	}
	return model
}

func (m *FactManagementModel) Init() tea.Cmd {
	// Apply theme to TableBehavior if available
	if theme := m.Theme(); theme != nil {
		m.tableBehavior.SetTheme(theme)
	}

	// Mark intent as active
	m.active = true

	// context already set in data
	m.data.CurrentState = FactListState

	if err := m.data.LoadFacts(); err != nil {
		m.result = &IntentResult[*FactManagementResult]{
			Status: Failed,
			Error: &IntentError{
				Code:    "LOAD_FACTS_FAILED",
				Message: "Failed to load facts",
				Cause:   err,
			},
		}
		return tea.Quit
	}

	// Set items in TableBehavior (replaces updateTableRows)
	m.tableBehavior.SetItems(m.data.Facts)
	return nil
}

// syncTableSelection syncs the TableBehavior selection with the data context
func (m *FactManagementModel) syncTableSelection() {
	// Update data context from TableBehavior
	m.data.SelectedFactIndex = m.tableBehavior.GetSelectedIndex()
	if selected := m.tableBehavior.GetSelectedItem(); selected != nil {
		m.data.SelectedFact = *selected
	} else {
		m.data.SelectedFact = nil
	}
}

func (m *FactManagementModel) Update(msg tea.Msg) tea.Cmd {
	switch m.data.CurrentState {
	case FactListState:
		return m.handleListState(msg)
	case FactViewState:
		return m.handleViewState(msg)
	case FactEditorState:
		return m.handleEditorState(msg)
	case FactDeleteConfirmState:
		return m.handleDeleteConfirmState(msg)
	case FactResultsState:
		return m.handleResultsState(msg)
	case FactCompletedState:
		// Do not send tea.Quit, just return nil and let router manage exit.
		return nil
	}
	return nil
}

func (m *FactManagementModel) View() string {
	if !m.active {
		return "FactManagement intent is not active"
	}

	// Create standard view with dynamic breadcrumbs
	view := CreateStandardViewWithBreadcrumbs(m.BaseIntent, m.getBreadcrumbs()...)

	// Handle form errors
	if m.data.HasFormErrors() {
		var errorMessages []string
		for field, err := range m.data.FormErrors {
			errorMessages = append(errorMessages, fmt.Sprintf("%s: %s", field, err))
		}
		m.SetError(fmt.Errorf("validation errors:\n%s", fmt.Sprintf("%v", errorMessages)))
	}

	// Get content for current state
	content := m.getStateContent()
	view.WithContent(content)

	// Get context-aware help
	help := m.getContextHelp()
	view.WithHelp(help).WithFooterSeparator(true)

	return view.Render()
}

// Helper methods for StandardView

func (m *FactManagementModel) getBreadcrumbs() []string {
	breadcrumbs := []string{"Main Menu", "Manage Facts"}

	switch m.data.CurrentState {
	case FactViewState, FactEditorState, FactDeleteConfirmState:
		if m.data.SelectedFact != nil {
			factID := m.data.SelectedFact.ID
			if len(factID) > 8 {
				factID = factID[:8]
			}
			factName := fmt.Sprintf("Fact #%s", factID)
			breadcrumbs = append(breadcrumbs, factName)
		}
	case FactResultsState:
		breadcrumbs = append(breadcrumbs, "Results")
	}

	return breadcrumbs
}

func (m *FactManagementModel) getStateContent() string {
	switch m.data.CurrentState {
	case FactListState:
		// Table is self-contained, just render it
		return m.viewList()
	case FactViewState:
		return m.getViewFactContent()
	case FactEditorState:
		return m.getEditorContent()
	case FactDeleteConfirmState:
		return m.getDeleteConfirmContent()
	case FactResultsState:
		return m.getResultsContent()
	case FactCompletedState:
		return "✅ Fact management completed"
	}
	return "Unknown state"
}

func (m *FactManagementModel) getContextHelp() string {
	theme := m.Theme()

	switch m.data.CurrentState {
	case FactListState:
		return CombineThemedFooters(
			ThemedListFooter(theme),
			ThemedCustomFooter(theme,
				primitives.EditBadge(theme),
				primitives.DeleteBadge(theme),
				primitives.HelpKeyBadge("n", "New", theme),
				primitives.HelpKeyBadge("r", "Refresh", theme),
			),
			ThemedGlobalBadges(theme),
		)
	case FactViewState:
		return CombineThemedFooters(
			ThemedDetailViewFooter(theme),
			ThemedCustomFooter(theme,
				primitives.EditBadge(theme),
				primitives.DeleteBadge(theme),
			),
			ThemedGlobalBadges(theme),
		)
	case FactEditorState:
		return CombineThemedFooters(
			ThemedFormFooter(theme),
			ThemedCustomFooter(theme,
				primitives.SaveBadge(theme),
			),
			ThemedGlobalBadges(theme),
		)
	case FactDeleteConfirmState:
		return CombineThemedFooters(
			ThemedCustomFooter(theme,
				primitives.HelpKeyBadge("y/Enter", "Confirm", theme),
				primitives.HelpKeyBadge("n/Esc", "Cancel", theme),
			),
			ThemedGlobalBadges(theme),
		)
	case FactResultsState:
		return CombineThemedFooters(
			ThemedCustomFooter(theme,
				primitives.BackBadge(theme),
			),
			ThemedGlobalBadges(theme),
		)
	default:
		return ThemedGlobalBadges(theme)
	}
}

func (m *FactManagementModel) Result() *IntentResult[interface{}] {
	if m.result == nil {
		return nil
	}
	return &IntentResult[interface{}]{
		Status: m.result.Status,
		Error:  m.result.Error,
	}
}

// ListNavigator interface implementation
func (m *FactManagementModel) GetTotalItems() int {
	return len(m.data.Facts)
}

func (m *FactManagementModel) GetSelectedIndex() int {
	return m.data.SelectedFactIndex
}

func (m *FactManagementModel) SetSelectedIndex(idx int) {
	m.tableBehavior.SetSelectedIndex(idx)
	m.syncTableSelection()
}

func (m *FactManagementModel) GetPageSize() int {
	return 15
}

func (m *FactManagementModel) handleListState(msg tea.Msg) tea.Cmd {
	switch msg := msg.(type) {
	case tea.KeyMsg:
		// Handle global keys first (q=quit, ?=help, esc=back)
		switch HandleGlobalKeys(msg) {
		case KeyQuit:
			return tea.Quit
		case KeyHelp:
			m.ToggleHelp()
			return nil
		case KeyBack:
			// At root state, back means cancel and return to main menu
			m.data.CurrentState = FactCompletedState
			m.result = &IntentResult[*FactManagementResult]{
				Status: Cancelled,
				Data: &FactManagementResult{
					Action: "none",
					Facts:  m.data.Facts,
				},
			}
			return nil
		}

		// Try TableBehavior navigation handler
		if m.tableBehavior.HandleNavigation(msg.String()) {
			m.syncTableSelection()
			return nil
		}

		switch msg.String() {
		case "enter", " ":
			m.syncTableSelection()
			if m.data.SelectedFact != nil {
				m.data.CurrentState = FactViewState
			}
			return nil

		case "n":
			m.data.StartNewFact()
			m.editModal = NewEditFactModal(m.data.EditingFact)
			m.data.CurrentState = FactEditorState
			// Return form init command to properly initialize the huh form
			return m.editModal.form.Init()

		case "r":
			if err := m.data.LoadFacts(); err != nil {
				m.result = &IntentResult[*FactManagementResult]{
					Status: Failed,
					Error: &IntentError{
						Code:    "RELOAD_FAILED",
						Message: "Failed to reload facts",
						Cause:   err,
					},
				}
			} else {
				m.tableBehavior.SetItems(m.data.Facts)
				m.syncTableSelection()
			}

			return nil

		case "tab", "right":
			// Page navigation handled by TableBehavior navigation
			return nil

		case "shift+tab", "left":
			// Page navigation handled by TableBehavior navigation
			return nil
		}
	}
	return nil
}

func (m *FactManagementModel) handleViewState(msg tea.Msg) tea.Cmd {
	switch msg := msg.(type) {
	case tea.KeyMsg:
		// Handle global keys first (q=quit, ?=help, esc=back)
		switch HandleGlobalKeys(msg) {
		case KeyQuit:
			return tea.Quit
		case KeyHelp:
			m.ToggleHelp()
			return nil
		case KeyBack:
			// Go back to list state
			m.data.CurrentState = FactListState
			m.data.SelectedFact = nil
			return nil
		}

		switch msg.String() {
		case "e":
			if m.data.SelectedFact != nil {
				m.data.StartEditFact(m.data.SelectedFact)
				m.editModal = NewEditFactModal(m.data.EditingFact)
				m.data.CurrentState = FactEditorState
				// Return form init command to properly initialize the huh form
				return m.editModal.form.Init()
			}

		case "d":
			if m.data.SelectedFact != nil {
				m.data.FactToDelete = m.data.SelectedFact
				m.data.CurrentState = FactDeleteConfirmState
			}
		}
	}
	return nil
}

func (m *FactManagementModel) handleEditorState(msg tea.Msg) tea.Cmd {
	// Handle global keys FIRST (before delegating to modal or legacy handling)
	// This ensures esc, q, ?, m keys work even when modal has focus
	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch HandleGlobalKeys(msg) {
		case KeyQuit:
			return tea.Quit
		case KeyHelp:
			m.ToggleHelp()
			return nil
		case KeyBack:
			// Close modal if active
			if m.editModal != nil {
				m.editModal = nil
			}
			m.data.CancelEdit()
			if m.data.IsNewFact {
				m.data.CurrentState = FactListState
			} else {
				m.data.CurrentState = FactViewState
			}
			return nil
		}
	}

	// If modal is not initialized, handle legacy behavior (fallback)
	if m.editModal == nil {
		return nil
	}

	// NOW delegate to the modal for form handling
	cmd := m.editModal.Update(msg)

	// Check if modal completed (form submitted or cancelled)
	if result := m.editModal.Result(); result != nil {
		if result.Accepted {
			// Apply changes from the modal to the editing fact
			m.data.EditingFact.Text = result.Modified.Text
			m.data.EditingFact.CompetencyCategories = result.Modified.CompetencyCategories
			m.data.EditingFact.RoleFit = result.Modified.RoleFit
			m.data.EditingFact.AudienceRelevance = result.Modified.AudienceRelevance
			m.data.EditingFact.StrengthSignal = result.Modified.StrengthSignal

			// Save the fact
			if err := m.data.SaveEdit(); err != nil {
				m.data.SetFormError("general", fmt.Sprintf("Save failed: %v", err))
			} else {
				action := "updated"
				if m.data.IsNewFact {
					action = "created"
				}
				m.result = &IntentResult[*FactManagementResult]{
					Status: Completed,
					Data: &FactManagementResult{
						Action:  action,
						Fact:    m.data.EditingFact,
						Facts:   m.data.Facts,
						Message: fmt.Sprintf("Fact %s successfully", action),
					},
				}
				// Refresh table after save
				m.tableBehavior.SetItems(m.data.Facts)
				m.syncTableSelection()
			}
		}

		// Clear modal and return to appropriate state
		m.editModal = nil
		m.data.CancelEdit()
		if m.data.IsNewFact {
			m.data.CurrentState = FactListState
		} else {
			m.data.CurrentState = FactViewState
		}
		return nil
	}

	return cmd
}

func (m *FactManagementModel) handleDeleteConfirmState(msg tea.Msg) tea.Cmd {
	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch msg.String() {
		case "y":
			if m.data.FactToDelete != nil {
				if err := m.data.DeleteFact(m.data.FactToDelete.ID); err != nil {
					m.result = &IntentResult[*FactManagementResult]{
						Status: Failed,
						Error: &IntentError{
							Code:    "DELETE_FAILED",
							Message: "Failed to delete fact",
							Cause:   err,
						},
					}
				} else {
					m.result = &IntentResult[*FactManagementResult]{
						Status: Completed,
						Data: &FactManagementResult{
							Action:  "deleted",
							Fact:    m.data.FactToDelete,
							Facts:   m.data.Facts,
							Message: "Fact deleted successfully",
						},
					}
					// Refresh table after delete
					m.tableBehavior.SetItems(m.data.Facts)
					m.syncTableSelection()
				}
				m.data.FactToDelete = nil
				m.data.CurrentState = FactListState
			}

		case "n", "esc":
			m.data.FactToDelete = nil
			m.data.CurrentState = FactViewState
		}
	}
	return nil
}

func (m *FactManagementModel) handleResultsState(msg tea.Msg) tea.Cmd {
	switch msg := msg.(type) {
	case tea.KeyMsg:
		// Handle global keys first (q=quit, ?=help, esc=back)
		switch HandleGlobalKeys(msg) {
		case KeyQuit:
			return tea.Quit
		case KeyHelp:
			m.ToggleHelp()
			return nil
		case KeyBack:
			// Go back to list state
			m.data.CurrentState = FactListState
			return nil
		}
	}
	return nil
}

// View rendering methods

func (m *FactManagementModel) viewList() string {
	// TableBehavior handles empty state, pagination, and rendering
	return m.tableBehavior.Render()
}

// viewFact removed - unused wrapper method

func (m *FactManagementModel) getViewFactContent() string {
	if m.data.SelectedFact == nil {
		return "No fact selected"
	}

	fact := m.data.SelectedFact
	var content string
	content += "📝 Fact Details\n\n"
	content += fmt.Sprintf("Text: %s\n\n", fact.Text)
	content += fmt.Sprintf("Categories: %v\n", fact.CompetencyCategories)
	content += fmt.Sprintf("Strength Signal: %s\n", fact.StrengthSignal)
	content += fmt.Sprintf("Role Fit: %v\n", fact.RoleFit)
	content += fmt.Sprintf("Audience Relevance: %v\n", fact.AudienceRelevance)
	content += fmt.Sprintf("\nCreated: %s\n", fact.CreatedAt.Format("2006-01-02 15:04"))

	return content
}

// viewEditor removed - unused wrapper method

func (m *FactManagementModel) getEditorContent() string {
	// If modal is available, render just the form content (not full modal container)
	// StandardView already provides the layout structure
	if m.editModal != nil {
		return m.editModal.GetContent()
	}

	// Fallback for legacy behavior
	var content string
	content += "✏️  Edit Fact\n\n"

	if m.data.EditingFact != nil {
		content += fmt.Sprintf("Text: %s\n\n", m.data.EditingFact.Text)
		content += fmt.Sprintf("Categories: %v\n", m.data.EditingFact.CompetencyCategories)
		content += fmt.Sprintf("Strength Signal: %s\n", m.data.EditingFact.StrengthSignal)
		content += fmt.Sprintf("Role Fit: %v\n\n", m.data.EditingFact.RoleFit)

		// Errors are shown in modal via getContextHelp
		if !m.data.HasFormErrors() {
			content += "Make your changes and press Ctrl+S to save.\n"
		}
	} else {
		content += "No fact loaded for editing.\n"
	}

	return content
}

// viewDeleteConfirm removed - unused wrapper method

func (m *FactManagementModel) getDeleteConfirmContent() string {
	if m.data.FactToDelete == nil {
		return "No fact to delete"
	}

	var content string
	content += "⚠️  Confirm Deletion\n\n"
	content += "Delete this fact?\n\n"
	content += fmt.Sprintf("Text: %s\n\n", truncate(m.data.FactToDelete.Text, 100))
	content += "⚠️  Warning: This action cannot be undone.\n"

	return content
}

// viewResults removed - unused wrapper method

func (m *FactManagementModel) getResultsContent() string {
	var content string
	content += "📊 Fact Statistics\n\n"
	content += fmt.Sprintf("Total facts in database: %d\n", m.data.TotalFacts)

	if len(m.data.Facts) > 0 {
		content += fmt.Sprintf("Currently viewing: %d facts\n", len(m.data.Facts))
	}

	return content
}

func truncate(s string, maxLen int) string {
	if len(s) <= maxLen {
		return s
	}
	return s[:maxLen] + "..."
}
