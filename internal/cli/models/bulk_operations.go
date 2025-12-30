package models

import (
	"context"
	"fmt"
	"strings"

	"github.com/baphled/kariya/internal/cli/service"
	"github.com/baphled/kariya/internal/cli/styles"
	"github.com/baphled/kariya/internal/domain/career"
	careerservice "github.com/baphled/kariya/internal/service/career"
	tea "github.com/charmbracelet/bubbletea"
)

// BulkOperationsModel represents the bulk operations screen state
type BulkOperationsModel struct {
	events           []*career.CareerEvent
	service          *careerservice.Service
	cliService       *service.CLIEventService
	ctx              context.Context
	selectedIndices  map[int]bool // Map of event index to selection state
	selectedEventIdx int          // Currently focused event index
	width            int
	height           int

	// Edit mode state
	inEditMode          bool
	bulkCompanyInput    string
	bulkProjectInput    string
	bulkTagsInput       []string
	bulkCategoriesInput []string
	bulkEditFieldIdx    int

	// Preview and confirmation state
	showPreview      bool
	showConfirmation bool
	appliedChanges   bool
	cancelled        bool

	// Undo state
	previousState map[int]*career.CareerEvent // Map of event index to previous state

	err error
}

// NewBulkOperationsModel creates a new bulk operations model
func NewBulkOperationsModel(events []*career.CareerEvent, svc *careerservice.Service, cliSvc *service.CLIEventService, ctx context.Context) *BulkOperationsModel {
	return &BulkOperationsModel{
		events:           events,
		service:          svc,
		cliService:       cliSvc,
		ctx:              ctx,
		selectedIndices:  make(map[int]bool),
		previousState:    make(map[int]*career.CareerEvent),
		selectedEventIdx: 0,
		inEditMode:       false,
		showPreview:      false,
		showConfirmation: false,
		appliedChanges:   false,
		cancelled:        false,
	}
}

// Init initializes the model
func (m *BulkOperationsModel) Init() tea.Cmd {
	return nil
}

// Update handles messages
func (m *BulkOperationsModel) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.KeyMsg:
		return m.handleKeyMsg(msg)
	case tea.WindowSizeMsg:
		m.width = msg.Width
		m.height = msg.Height
		return m, nil
	}
	return m, nil
}

// handleKeyMsg handles keyboard input
func (m *BulkOperationsModel) handleKeyMsg(msg tea.KeyMsg) (tea.Model, tea.Cmd) {
	switch msg.Type {
	case tea.KeySpace:
		m.ToggleSelection(m.selectedEventIdx)
		return m, nil
	case tea.KeyUp:
		m.selectedEventIdx = (m.selectedEventIdx - 1 + len(m.events)) % len(m.events)
		return m, nil
	case tea.KeyDown:
		m.selectedEventIdx = (m.selectedEventIdx + 1) % len(m.events)
		return m, nil
	case tea.KeyRunes:
		for _, r := range msg.Runes {
			switch r {
			case 'a':
				m.SelectAll()
			case 'd':
				m.DeselectAll()
			case 'e':
				if m.GetSelectedCount() > 0 {
					m.EnterEditMode()
				}
			case 'u':
				m.Undo()
			case 'y':
				if m.showConfirmation {
					m.ApplyChanges()
					m.showConfirmation = false
				}
			case 'n':
				if m.showConfirmation {
					m.showConfirmation = false
				}
			}
		}
		return m, nil
	case tea.KeyEnter:
		if m.inEditMode && !m.showPreview {
			m.ShowPreview()
		} else if m.showPreview && !m.showConfirmation {
			m.ShowConfirmation()
		}
		return m, nil
	case tea.KeyEsc:
		if m.inEditMode {
			m.ExitEditMode()
		} else if m.showPreview {
			m.showPreview = false
		} else if m.showConfirmation {
			m.showConfirmation = false
		} else {
			m.cancelled = true
		}
		return m, nil
	}
	return m, nil
}

// View renders the model
func (m *BulkOperationsModel) View() string {
	if len(m.events) == 0 {
		return styles.InfoBox.Render("No events available for bulk operations")
	}

	if m.showConfirmation {
		return m.renderConfirmation()
	}

	if m.showPreview {
		return m.renderPreview()
	}

	if m.inEditMode {
		return m.renderEditMode()
	}

	return m.renderSelectionMode()
}

// renderSelectionMode renders the event selection view
func (m *BulkOperationsModel) renderSelectionMode() string {
	var sb strings.Builder

	// Header
	header := fmt.Sprintf("Bulk Operations - Select Events (%d/%d selected)", m.GetSelectedCount(), len(m.events))
	sb.WriteString(styles.HeaderSection.Render(header))
	sb.WriteString("\n\n")

	// Events list
	for i, event := range m.events {
		checkbox := "[ ]"
		if m.selectedIndices[i] {
			checkbox = "[x]"
		}

		isSelected := i == m.selectedEventIdx
		eventText := fmt.Sprintf("%s %s", checkbox, truncateText(event.Text, 40))

		if isSelected {
			sb.WriteString(styles.ListItemSelected.Render(eventText))
		} else {
			sb.WriteString(styles.ListItem.Render(eventText))
		}
		sb.WriteString("\n")
	}

	// Instructions
	sb.WriteString("\n")
	sb.WriteString(styles.InfoHint.Render("Space: toggle | a: select all | d: deselect all | e: edit selected | q: quit"))

	return sb.String()
}

// renderEditMode renders the bulk edit form
func (m *BulkOperationsModel) renderEditMode() string {
	var sb strings.Builder

	sb.WriteString(styles.HeaderSection.Render("Bulk Edit Selected Events"))
	sb.WriteString("\n\n")

	// Display selected events
	sb.WriteString(styles.HeaderSubsection.Render("Selected Events:"))
	sb.WriteString("\n")
	for i, event := range m.events {
		if m.selectedIndices[i] {
			sb.WriteString(fmt.Sprintf("  • %s\n", truncateText(event.Text, 50)))
		}
	}

	sb.WriteString("\n")
	sb.WriteString(styles.HeaderSubsection.Render("Edit Fields:"))
	sb.WriteString("\n")

	// Company field
	companyLabel := "Company:"
	if m.bulkEditFieldIdx == 0 {
		companyLabel = "> " + companyLabel
	}
	sb.WriteString(fmt.Sprintf("%s %s\n", companyLabel, m.bulkCompanyInput))

	// Project field
	projectLabel := "Project:"
	if m.bulkEditFieldIdx == 1 {
		projectLabel = "> " + projectLabel
	}
	sb.WriteString(fmt.Sprintf("%s %s\n", projectLabel, m.bulkProjectInput))

	sb.WriteString("\n")
	sb.WriteString(styles.InfoHint.Render("Enter: preview changes | Esc: cancel | Tab: next field"))

	return sb.String()
}

// renderPreview renders the preview of bulk changes
func (m *BulkOperationsModel) renderPreview() string {
	var sb strings.Builder

	sb.WriteString(styles.HeaderSection.Render("Preview Bulk Changes"))
	sb.WriteString("\n\n")

	sb.WriteString(styles.HeaderSubsection.Render("Changes to apply:"))
	sb.WriteString("\n")

	if m.bulkCompanyInput != "" {
		sb.WriteString(fmt.Sprintf("  Company: %s\n", m.bulkCompanyInput))
	}
	if m.bulkProjectInput != "" {
		sb.WriteString(fmt.Sprintf("  Project: %s\n", m.bulkProjectInput))
	}

	sb.WriteString("\n")
	sb.WriteString(styles.HeaderSubsection.Render("Affected Events:"))
	sb.WriteString("\n")

	for i, event := range m.events {
		if m.selectedIndices[i] {
			sb.WriteString(fmt.Sprintf("  • %s\n", truncateText(event.Text, 50)))
		}
	}

	sb.WriteString("\n")
	sb.WriteString(styles.InfoHint.Render("Enter: confirm | Esc: cancel"))

	return sb.String()
}

// renderConfirmation renders the confirmation dialog
func (m *BulkOperationsModel) renderConfirmation() string {
	var sb strings.Builder

	sb.WriteString(styles.WarningBox.Render("Confirm Bulk Changes"))
	sb.WriteString("\n\n")

	sb.WriteString(fmt.Sprintf("Apply changes to %d event(s)?\n\n", m.GetSelectedCount()))
	sb.WriteString(styles.InfoHint.Render("y: confirm | n: cancel"))

	return sb.String()
}

// Selection methods

// ToggleSelection toggles selection for an event
func (m *BulkOperationsModel) ToggleSelection(idx int) {
	if idx >= 0 && idx < len(m.events) {
		m.selectedIndices[idx] = !m.selectedIndices[idx]
	}
}

// SelectAll selects all events
func (m *BulkOperationsModel) SelectAll() {
	for i := 0; i < len(m.events); i++ {
		m.selectedIndices[i] = true
	}
}

// DeselectAll deselects all events
func (m *BulkOperationsModel) DeselectAll() {
	m.selectedIndices = make(map[int]bool)
}

// GetSelectedCount returns the number of selected events
func (m *BulkOperationsModel) GetSelectedCount() int {
	count := 0
	for _, selected := range m.selectedIndices {
		if selected {
			count++
		}
	}
	return count
}

// GetSelectedEventIDs returns IDs of selected events
func (m *BulkOperationsModel) GetSelectedEventIDs() []string {
	var ids []string
	for i, selected := range m.selectedIndices {
		if selected && i < len(m.events) {
			ids = append(ids, m.events[i].ID)
		}
	}
	return ids
}

// Edit mode methods

// EnterEditMode enters bulk edit mode
func (m *BulkOperationsModel) EnterEditMode() {
	m.inEditMode = true
	m.bulkEditFieldIdx = 0
	// Save current state for undo
	for i, selected := range m.selectedIndices {
		if selected && i < len(m.events) {
			eventCopy := *m.events[i]
			m.previousState[i] = &eventCopy
		}
	}
}

// ExitEditMode exits bulk edit mode
func (m *BulkOperationsModel) ExitEditMode() {
	m.inEditMode = false
	m.bulkCompanyInput = ""
	m.bulkProjectInput = ""
	m.bulkTagsInput = []string{}
	m.bulkCategoriesInput = []string{}
	m.bulkEditFieldIdx = 0
}

// IsInEditMode returns whether model is in edit mode
func (m *BulkOperationsModel) IsInEditMode() bool {
	return m.inEditMode
}

// SetBulkCompany sets the bulk company input
func (m *BulkOperationsModel) SetBulkCompany(company string) {
	m.bulkCompanyInput = company
}

// GetBulkCompany gets the bulk company input
func (m *BulkOperationsModel) GetBulkCompany() string {
	return m.bulkCompanyInput
}

// SetBulkProject sets the bulk project input
func (m *BulkOperationsModel) SetBulkProject(project string) {
	m.bulkProjectInput = project
}

// GetBulkProject gets the bulk project input
func (m *BulkOperationsModel) GetBulkProject() string {
	return m.bulkProjectInput
}

// Preview and confirmation methods

// ShowPreview shows the preview
func (m *BulkOperationsModel) ShowPreview() {
	m.showPreview = true
}

// IsShowingPreview returns whether preview is shown
func (m *BulkOperationsModel) IsShowingPreview() bool {
	return m.showPreview
}

// ShowConfirmation shows the confirmation dialog
func (m *BulkOperationsModel) ShowConfirmation() {
	m.showConfirmation = true
}

// IsShowingConfirmation returns whether confirmation is shown
func (m *BulkOperationsModel) IsShowingConfirmation() bool {
	return m.showConfirmation
}

// ApplyChanges applies the bulk changes
func (m *BulkOperationsModel) ApplyChanges() {
	for i, selected := range m.selectedIndices {
		if selected && i < len(m.events) {
			event := m.events[i]
			if m.bulkCompanyInput != "" {
				event.Company = m.bulkCompanyInput
			}
			if m.bulkProjectInput != "" {
				event.Project = m.bulkProjectInput
			}
		}
	}
	m.appliedChanges = true
	m.ExitEditMode()
}

// Undo reverts to previous state
func (m *BulkOperationsModel) Undo() {
	for i, prevState := range m.previousState {
		if i < len(m.events) && prevState != nil {
			m.events[i] = prevState
		}
	}
	m.previousState = make(map[int]*career.CareerEvent)
}

// State query methods

// ChangesApplied returns whether changes were applied
func (m *BulkOperationsModel) ChangesApplied() bool {
	return m.appliedChanges
}

// MarkChangesApplied marks changes as applied
func (m *BulkOperationsModel) MarkChangesApplied() {
	m.appliedChanges = true
}

// WasCancelled returns whether operation was cancelled
func (m *BulkOperationsModel) WasCancelled() bool {
	return m.cancelled
}

// Cancel cancels the operation
func (m *BulkOperationsModel) Cancel() {
	m.cancelled = true
}

// Navigation methods

// GetCurrentIndex returns the currently focused event index
func (m *BulkOperationsModel) GetCurrentIndex() int {
	return m.selectedEventIdx
}

// SetCurrentIndex sets the currently focused event index
func (m *BulkOperationsModel) SetCurrentIndex(idx int) {
	if idx >= 0 && idx < len(m.events) {
		m.selectedEventIdx = idx
	}
}

// GetEvents returns the events list
func (m *BulkOperationsModel) GetEvents() []*career.CareerEvent {
	return m.events
}

// SetSize sets the terminal size
func (m *BulkOperationsModel) SetSize(width, height int) {
	m.width = width
	m.height = height
}

// Helper function to truncate text
func truncateText(text string, maxLen int) string {
	if len(text) <= maxLen {
		return text
	}
	return text[:maxLen-3] + "..."
}
