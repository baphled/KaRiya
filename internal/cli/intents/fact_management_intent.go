package intents

import (
	"fmt"

	"github.com/baphled/kariya/internal/cli/components"
	"github.com/baphled/kariya/internal/cli/navigation"
	"github.com/baphled/kariya/internal/cli/themes"
	"github.com/charmbracelet/bubbles/table"
	tea "github.com/charmbracelet/bubbletea"
)

type FactManagementModel struct {
	*BaseIntent
	data          *FactManagementContext
	table         *table.Model
	listContainer *components.TableListContainer
	navHandler    *navigation.ListNavigationHandler
	result        *IntentResult[*FactManagementResult]
	active        bool
}

func NewFactManagementIntent(data *FactManagementContext) *FactManagementModel {
	// Create table model for facts
	columns := []table.Column{
		{Title: "Fact", Width: 50},
		{Title: "Strength", Width: 15},
		{Title: "Categories", Width: 30},
	}

	t := table.New(
		table.WithColumns(columns),
		table.WithRows([]table.Row{}),
		table.WithFocused(true),
		table.WithHeight(15),
		table.WithWidth(100),
	)

	// Apply default styles initially - theme styles will be applied in Init()
	t.SetStyles(table.DefaultStyles())

	model := &FactManagementModel{
		BaseIntent:    NewBaseIntent(),
		data:          data,
		table:         &t,
		listContainer: components.NewTableListContainer(t, "Manage Facts", 100),
		result:        nil,
		active:        false,
	}
	model.navHandler = navigation.NewListNavigationHandler(model)
	return model
}

func (m *FactManagementModel) Init() tea.Cmd {
	// Apply themed table styles if theme is available
	if theme := m.Theme(); theme != nil {
		m.table.SetStyles(themes.NewThemedTableStyles(theme))
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

	m.updateTableRows()
	return nil
}

// updateTableRows updates the table rows based on facts
func (m *FactManagementModel) updateTableRows() {
	pageSize := 15
	total := len(m.data.Facts)

	// Determine which page current selection is on
	page := 0
	if pageSize > 0 && m.data.SelectedFactIndex >= 0 {
		page = m.data.SelectedFactIndex / pageSize
	}

	start := page * pageSize
	end := start + pageSize
	if end > total {
		end = total
	}

	pageFacts := m.data.Facts[start:end]

	rows := make([]table.Row, 0, len(pageFacts))
	for idx, fact := range pageFacts {
		realIdx := start + idx
		text := truncate(fact.Text, 50)

		// Use navigation handler to format row text with indicator
		text = m.navHandler.FormatRowText(realIdx, text)

		strength := fact.StrengthSignal
		if strength == "" {
			strength = "-"
		}
		categories := fmt.Sprintf("%v", fact.CompetencyCategories)
		if len(categories) > 30 {
			categories = categories[:27] + "..."
		}
		rows = append(rows, table.Row{text, strength, categories})
	}

	m.table.SetRows(rows)

	// Calculate relative cursor for current page
	relativeCursor := m.data.SelectedFactIndex - start

	// Set table cursor (for visual highlighting)
	m.table.SetCursor(relativeCursor)

	// Sync container's index to match (critical for rendering)
	m.listContainer.SetSelectedIdx(relativeCursor)

	// Update container with modified table
	m.listContainer.SetTable(*m.table)
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
			factName := fmt.Sprintf("Fact #%s", m.data.SelectedFact.ID[:8])
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
				components.EditBadge(),
				components.DeleteBadge(),
				components.NewKeyBadge("n", "New"),
				components.NewKeyBadge("r", "Refresh"),
			),
			ThemedGlobalBadges(theme),
		)
	case FactViewState:
		return CombineThemedFooters(
			ThemedDetailViewFooter(theme),
			ThemedCustomFooter(theme,
				components.EditBadge(),
				components.DeleteBadge(),
			),
			ThemedGlobalBadges(theme),
		)
	case FactEditorState:
		return CombineThemedFooters(
			ThemedFormFooter(theme),
			ThemedCustomFooter(theme,
				components.SaveBadge(),
			),
			ThemedGlobalBadges(theme),
		)
	case FactDeleteConfirmState:
		return CombineThemedFooters(
			ThemedCustomFooter(theme,
				components.NewKeyBadge("y/Enter", "Confirm"),
				components.NewKeyBadge("n/Esc", "Cancel"),
			),
			ThemedGlobalBadges(theme),
		)
	case FactResultsState:
		return CombineThemedFooters(
			ThemedCustomFooter(theme,
				components.BackBadge(),
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
	m.data.SelectedFactIndex = idx
	if idx >= 0 && idx < len(m.data.Facts) {
		m.data.SelectedFact = m.data.Facts[idx]
	}
	m.updateTableRows()
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

		// Try navigation handler
		if m.navHandler.HandleKey(msg.String()) {
			return nil
		}

		switch msg.String() {

		case "enter", " ":
			if m.data.SelectedFact != nil {
				m.data.CurrentState = FactViewState
			}
			return nil

		case "n":
			m.data.StartNewFact()
			m.data.CurrentState = FactEditorState

			return nil

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
				m.updateTableRows()
			}

			return nil

		case "tab", "right":
			if m.data.CurrentPage < (m.data.TotalFacts / m.data.PageSize) {
				m.data.CurrentPage++
				m.updateTableRows()
			}

			return nil

		case "shift+tab", "left":
			if m.data.CurrentPage > 0 {
				m.data.CurrentPage--
				m.updateTableRows()
			}

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
				m.data.CurrentState = FactEditorState
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
	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch msg.String() {
		case "ctrl+s":
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
				m.data.CurrentState = FactListState
				m.data.CancelEdit()
				m.updateTableRows()
			}

		case "esc":
			m.data.CancelEdit()
			if m.data.IsNewFact {
				m.data.CurrentState = FactListState
			} else {
				m.data.CurrentState = FactViewState
			}
		}
	}
	return nil
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
					m.updateTableRows()
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
	if len(m.data.Facts) == 0 {
		m.listContainer.SetEmptyStateMessage("No facts found. Press 'n' to create a new fact, 'r' to refresh, or 'q' to quit.")
		return m.listContainer.Render()
	}

	// Ensure table rows are synchronized with current state
	m.updateTableRows()

	// Build pagination info with page number indicator
	pageSize := 15
	totalItems := len(m.data.Facts)
	currentPage := (m.data.SelectedFactIndex / pageSize) + 1
	totalPages := (totalItems + pageSize - 1) / pageSize
	paginationInfo := fmt.Sprintf("Total: %d facts | Page %d of %d", totalItems, currentPage, totalPages)
	m.listContainer.SetPaginationInfo(paginationInfo)

	return m.listContainer.Render()
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
