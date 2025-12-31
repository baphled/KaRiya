package models

import (
	"context"
	"fmt"
	"strings"

	"github.com/baphled/kariya/internal/cli/components"
	"github.com/baphled/kariya/internal/cli/service"
	"github.com/baphled/kariya/internal/cli/styles"
	"github.com/baphled/kariya/internal/domain/career"
	careerservice "github.com/baphled/kariya/internal/service/career"
	tea "github.com/charmbracelet/bubbletea"
)

// BulkOperationsSummary represents the result of bulk operations
type BulkOperationsSummary struct {
	EventsAffected int
	FieldsUpdated  []string
	Errors         []string
	AppliedCount   int
	SkippedCount   int
}

// BulkOperationsModel manages bulk metadata editing for multiple events
type BulkOperationsModel struct {
	events              []*career.CareerEvent
	service             *careerservice.Service
	cliService          *service.CLIEventService
	ctx                 context.Context
	selected            map[int]bool
	focusIndex          int
	width               int
	height              int
	bulkCompany         string
	bulkProject         string
	bulkTags            []string
	bulkCategories      []string
	applyIfEmptyCompany bool
	applyIfEmptyProject bool
	submitted           bool
	cancelled           bool
	summary             *BulkOperationsSummary
	fieldOrigins        map[string]map[string]bool // eventID -> field -> isFromCSV

	helpFooter components.HelpFooterModel // Help footer
}

// NewBulkOperationsModel creates a new bulk operations model
func NewBulkOperationsModel(
	events []*career.CareerEvent,
	service *careerservice.Service,
	cliService *service.CLIEventService,
	ctx context.Context,
) *BulkOperationsModel {
	return &BulkOperationsModel{
		events:       events,
		service:      service,
		cliService:   cliService,
		ctx:          ctx,
		selected:     make(map[int]bool),
		focusIndex:   0,
		width:        80,
		height:       24,
		fieldOrigins: make(map[string]map[string]bool),
		helpFooter:   components.NewHelpFooter("bulk_operations", 80),
	}
}

// Init initializes the model
func (m *BulkOperationsModel) Init() tea.Cmd {
	return nil
}

// Update handles input and updates the model
func (m *BulkOperationsModel) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.WindowSizeMsg:
		m.width = msg.Width
		m.height = msg.Height
		m.helpFooter.SetWidth(msg.Width)
		return m, nil
	case tea.KeyMsg:
		switch msg.Type {
		case tea.KeyUp:
			m.MoveUp()
		case tea.KeyDown:
			m.MoveDown()
		case tea.KeySpace:
			m.ToggleSelection(m.focusIndex)
		case tea.KeyEsc:
			m.Cancel()
			return m, nil
		case tea.KeyRunes:
			for _, r := range msg.Runes {
				switch r {
				case 'a':
					m.SelectAll()
				case 'd':
					m.DeselectAll()
				}
			}
		}

	}
	return m, nil
}

// View renders the bulk operations interface
func (m *BulkOperationsModel) View() string {
	if len(m.events) == 0 {
		return styles.InfoBox.Render("No events available for bulk operations")
	}

	var b strings.Builder

	// Header
	header := components.NewHeader("Bulk Operations", m.width)
	b.WriteString(header.View())
	b.WriteString("\n\n")

	// Selection summary
	selectedCount := m.GetSelectedCount()
	b.WriteString(fmt.Sprintf("Selected: %d/%d events\n", selectedCount, len(m.events)))
	b.WriteString("\n")

	// Event list with checkboxes
	for i, event := range m.events {
		checkbox := "[ ]"
		if m.selected[i] {
			checkbox = "[x]"
		}

		prefix := "  "
		if i == m.focusIndex {
			prefix = "▶ "
		}

		eventText := event.Text
		if len(eventText) > 50 {
			eventText = eventText[:47] + "..."
		}

		line := fmt.Sprintf("%s%s %s", prefix, checkbox, eventText)
		if i == m.focusIndex {
			b.WriteString(styles.ListItemSelected.Render(line))
		} else {
			b.WriteString(line)
		}
		b.WriteString("\n")
	}

	// Bulk edit fields
	if m.bulkCompany != "" || m.bulkProject != "" || len(m.bulkTags) > 0 {
		b.WriteString("\n")
		b.WriteString(styles.HeaderSubsection.Render("Bulk Updates"))
		b.WriteString("\n")

		if m.bulkCompany != "" {
			b.WriteString(fmt.Sprintf("Company: %s\n", m.bulkCompany))
		}
		if m.bulkProject != "" {
			b.WriteString(fmt.Sprintf("Project: %s\n", m.bulkProject))
		}
		if len(m.bulkTags) > 0 {
			b.WriteString(fmt.Sprintf("Tags: %s\n", strings.Join(m.bulkTags, ", ")))
		}
	}

	// Keyboard shortcuts
	b.WriteString("\n")
	// Help footer
	m.helpFooter.SetWidth(80)
	b.WriteString(m.helpFooter.View())

	return b.String()
}

// GetSelectedCount returns the number of selected events
func (m *BulkOperationsModel) GetSelectedCount() int {
	count := 0
	for _, selected := range m.selected {
		if selected {
			count++
		}
	}
	return count
}

// GetEventCount returns the total number of events
func (m *BulkOperationsModel) GetEventCount() int {
	return len(m.events)
}

// ToggleSelection toggles selection state for an event
func (m *BulkOperationsModel) ToggleSelection(index int) {
	if index >= 0 && index < len(m.events) {
		m.selected[index] = !m.selected[index]
	}
}

// SelectAll selects all events
func (m *BulkOperationsModel) SelectAll() {
	for i := 0; i < len(m.events); i++ {
		m.selected[i] = true
	}
}

// DeselectAll deselects all events
func (m *BulkOperationsModel) DeselectAll() {
	m.selected = make(map[int]bool)
}

// MoveUp moves focus up one event
func (m *BulkOperationsModel) MoveUp() {
	if m.focusIndex > 0 {
		m.focusIndex--
	}
}

// MoveDown moves focus down one event
func (m *BulkOperationsModel) MoveDown() {
	if m.focusIndex < len(m.events)-1 {
		m.focusIndex++
	}
}

// GetFocusIndex returns the current focus index
func (m *BulkOperationsModel) GetFocusIndex() int {
	return m.focusIndex
}

// SetBulkCompany sets the bulk company value
func (m *BulkOperationsModel) SetBulkCompany(company string) {
	m.bulkCompany = company
}

// GetBulkCompany returns the bulk company value
func (m *BulkOperationsModel) GetBulkCompany() string {
	return m.bulkCompany
}

// SetBulkProject sets the bulk project value
func (m *BulkOperationsModel) SetBulkProject(project string) {
	m.bulkProject = project
}

// GetBulkProject returns the bulk project value
func (m *BulkOperationsModel) GetBulkProject() string {
	return m.bulkProject
}

// SetBulkTags sets the bulk tags
func (m *BulkOperationsModel) SetBulkTags(tags []string) {
	m.bulkTags = tags
}

// GetBulkTags returns the bulk tags
func (m *BulkOperationsModel) GetBulkTags() []string {
	return m.bulkTags
}

// SetBulkCategories sets the bulk categories
func (m *BulkOperationsModel) SetBulkCategories(categories []string) {
	m.bulkCategories = categories
}

// GetBulkCategories returns the bulk categories
func (m *BulkOperationsModel) GetBulkCategories() []string {
	return m.bulkCategories
}

// SetApplyIfEmptyCompany sets the apply-if-empty flag for company
func (m *BulkOperationsModel) SetApplyIfEmptyCompany(apply bool) {
	m.applyIfEmptyCompany = apply
}

// GetApplyIfEmptyCompany returns the apply-if-empty flag for company
func (m *BulkOperationsModel) GetApplyIfEmptyCompany() bool {
	return m.applyIfEmptyCompany
}

// SetApplyIfEmptyProject sets the apply-if-empty flag for project
func (m *BulkOperationsModel) SetApplyIfEmptyProject(apply bool) {
	m.applyIfEmptyProject = apply
}

// GetApplyIfEmptyProject returns the apply-if-empty flag for project
func (m *BulkOperationsModel) GetApplyIfEmptyProject() bool {
	return m.applyIfEmptyProject
}

// GetPreview generates a preview of bulk changes
func (m *BulkOperationsModel) GetPreview() string {
	var b strings.Builder

	selectedCount := m.GetSelectedCount()
	if selectedCount == 0 {
		return "No events selected"
	}

	b.WriteString(fmt.Sprintf("Will affect %d events:\n", selectedCount))

	if m.bulkCompany != "" {
		b.WriteString(fmt.Sprintf("- company: %s\n", m.bulkCompany))
	}
	if m.bulkProject != "" {
		b.WriteString(fmt.Sprintf("- project: %s\n", m.bulkProject))
	}
	if len(m.bulkTags) > 0 {
		b.WriteString(fmt.Sprintf("- tags: %s\n", strings.Join(m.bulkTags, ", ")))
	}
	if len(m.bulkCategories) > 0 {
		b.WriteString(fmt.Sprintf("- categories: %s\n", strings.Join(m.bulkCategories, ", ")))
	}

	return b.String()
}

// Submit submits the bulk operation
func (m *BulkOperationsModel) Submit() {
	m.submitted = true
	m.cancelled = false
	m.summary = &BulkOperationsSummary{
		EventsAffected: m.GetSelectedCount(),
		FieldsUpdated:  []string{},
		Errors:         []string{},
	}

	if m.bulkCompany != "" {
		m.summary.FieldsUpdated = append(m.summary.FieldsUpdated, "company")
	}
	if m.bulkProject != "" {
		m.summary.FieldsUpdated = append(m.summary.FieldsUpdated, "project")
	}
	if len(m.bulkTags) > 0 {
		m.summary.FieldsUpdated = append(m.summary.FieldsUpdated, "tags")
	}
	if len(m.bulkCategories) > 0 {
		m.summary.FieldsUpdated = append(m.summary.FieldsUpdated, "categories")
	}
}

// IsSubmitted returns whether the operation was submitted
func (m *BulkOperationsModel) IsSubmitted() bool {
	return m.submitted
}

// GetSummary returns the operation summary
func (m *BulkOperationsModel) GetSummary() *BulkOperationsSummary {
	return m.summary
}

// Cancel cancels the bulk operation
func (m *BulkOperationsModel) Cancel() {
	m.cancelled = true
	m.submitted = false
	m.bulkCompany = ""
	m.bulkProject = ""
	m.bulkTags = nil
	m.bulkCategories = nil
	m.applyIfEmptyCompany = false
	m.applyIfEmptyProject = false
}

// IsCancelled returns whether the operation was cancelled
func (m *BulkOperationsModel) IsCancelled() bool {
	return m.cancelled
}

// SetSize sets the width and height of the model
func (m *BulkOperationsModel) SetSize(width, height int) {
	m.width = width
	m.height = height
}

// GetWidth returns the width
func (m *BulkOperationsModel) GetWidth() int {
	return m.width
}

// GetHeight returns the height
func (m *BulkOperationsModel) GetHeight() int {
	return m.height
}

// ChangesApplied returns whether changes were applied
func (m *BulkOperationsModel) ChangesApplied() bool {
	return m.submitted && !m.cancelled
}

// WasCancelled returns whether the operation was cancelled
func (m *BulkOperationsModel) WasCancelled() bool {
	return m.cancelled
}

// SetFieldOrigins sets the field origins for an event
func (m *BulkOperationsModel) SetFieldOrigins(eventID string, origins map[string]bool) {
	if m.fieldOrigins == nil {
		m.fieldOrigins = make(map[string]map[string]bool)
	}
	m.fieldOrigins[eventID] = origins
}

// GetFieldOrigins returns the field origins for an event
func (m *BulkOperationsModel) GetFieldOrigins(eventID string) map[string]bool {
	if m.fieldOrigins == nil {
		return nil
	}
	return m.fieldOrigins[eventID]
}

// CanUpdateField checks if a field can be updated for an event
// Returns true if field is from default (not from CSV), false if from CSV
func (m *BulkOperationsModel) CanUpdateField(eventID, field string) bool {
	origins := m.GetFieldOrigins(eventID)
	if origins == nil {
		// If no origin info, allow update
		return true
	}
	fromCSV, exists := origins[field]
	// Only allow updating fields that are NOT from CSV (are defaults)
	return exists && !fromCSV
}

// GetDefaultFields returns the list of default fields for an event
func (m *BulkOperationsModel) GetDefaultFields(eventID string) []string {
	origins := m.GetFieldOrigins(eventID)
	if origins == nil {
		return []string{}
	}
	var defaults []string
	for field, fromCSV := range origins {
		if !fromCSV {
			defaults = append(defaults, field)
		}
	}
	return defaults
}

// GetUpdatableFields returns fields that can be updated for an event
func (m *BulkOperationsModel) GetUpdatableFields(eventID string) []string {
	return m.GetDefaultFields(eventID)
}
