package modals

import (
	"fmt"

	"github.com/baphled/kariya/internal/cli/behaviors"
	"github.com/baphled/kariya/internal/cli/themes"
	"github.com/baphled/kariya/internal/cli/uikit/containers"
	"github.com/baphled/kariya/internal/cli/uikit/primitives"
	"github.com/baphled/kariya/internal/domain/career"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
)

// maxEventTextLength defines the maximum length for event text before truncation.
// This ensures event descriptions fit within the table column width.
const maxEventTextLength = 42

// EventsModal displays events that use a skill in a modal overlay with a data table.
// This modal uses TableBehavior for consistent table display and navigation.
//
// Features:
// - Data table display with Date, Event, and Company columns
// - Standard table navigation (j/k, up/down, page up/down)
// - Press Enter to select an event for detailed view
// - Close with Escape, backspace, or 'q'
//
// Usage:
//
//	modal := NewEventsModal(skillID, skillName, events, theme)
//	modal.SetDimensions(width, height)
//	modal.Show()
//
//	// In Update:
//	model, cmd := modal.Update(msg)
//
//	// Check for selection:
//	if modal.HasSelection() {
//	    selectedEvent := modal.GetSelectedEvent()
//	    // Show event detail modal
//	}
//
//	// In View:
//	if modal.IsVisible() {
//	    return renderModalOverlay(modal, background)
//	}
type EventsModal struct {
	skillID       string
	skillName     string
	events        []*career.CareerEvent
	theme         themes.Theme
	visible       bool
	width         int
	height        int
	table         *behaviors.TableBehavior[*career.CareerEvent]
	selectedEvent *career.CareerEvent // Set when user selects an event.
}

// eventRowFormatter formats a career event for table display.
func eventRowFormatter(event *career.CareerEvent, _ int) []string {
	if event == nil {
		return []string{"-", "(No event)", "-"}
	}

	// Date.
	dateStr := "-"
	if !event.Date.IsZero() {
		dateStr = event.Date.Format("2006-01-02")
	}

	// Event text (truncated to fit column).
	text := event.Text
	if text == "" {
		text = "(No description)"
	}
	if len(text) > maxEventTextLength {
		text = text[:maxEventTextLength] + "..."
	}

	// Company.
	company := event.Company
	if company == "" {
		company = "-"
	}

	return []string{dateStr, text, company}
}

// NewEventsModal creates a new events modal for a skill.
func NewEventsModal(skillID, skillName string, events []*career.CareerEvent, theme themes.Theme) *EventsModal {
	// Filter out nil events.
	filteredEvents := make([]*career.CareerEvent, 0, len(events))
	for _, e := range events {
		if e != nil {
			filteredEvents = append(filteredEvents, e)
		}
	}

	// Nil theme guard.
	if theme == nil {
		theme = themes.NewDefaultTheme()
	}

	// Define table columns.
	columns := []behaviors.ColumnDef{
		{Title: "Date", Width: 12},
		{Title: "Event", Width: 45},
		{Title: "Company", Width: 20},
	}

	// Create table behavior.
	tableBehavior := behaviors.NewTableBehavior(theme, columns, eventRowFormatter).
		PageSize(12).
		EmptyMessage("No events use this skill.").
		HidePagination()

	tableBehavior.SetItems(filteredEvents)

	m := &EventsModal{
		skillID:   skillID,
		skillName: skillName,
		events:    filteredEvents,
		theme:     theme,
		visible:   false,
		width:     100,
		height:    24,
		table:     tableBehavior,
	}
	return m
}

// Init initializes the modal (implements tea.Model for bubbletea-overlay).
func (m *EventsModal) Init() tea.Cmd {
	return nil
}

// Update handles keyboard input and window sizing.
func (m *EventsModal) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	if !m.visible {
		return m, nil
	}

	switch msg := msg.(type) {
	case tea.WindowSizeMsg:
		m.width = msg.Width
		m.height = msg.Height
		m.updateTableDimensions()
		return m, nil

	case tea.KeyMsg:
		keyStr := msg.String()

		switch keyStr {
		case "esc", "backspace", "q":
			// Close modal without selection.
			m.Hide()
			return m, nil

		case "enter":
			// Select the current event.
			if selected := m.table.GetSelectedItem(); selected != nil {
				m.selectedEvent = *selected
				m.Hide()
			}
			return m, nil
		}

		// Handle table navigation.
		if m.table.HandleNavigation(keyStr) {
			return m, nil
		}
	}

	return m, nil
}

// updateTableDimensions updates the table dimensions based on modal size.
func (m *EventsModal) updateTableDimensions() {
	tableHeight := m.height - 14
	if tableHeight < 5 {
		tableHeight = 5
	}
	if tableHeight > 12 {
		tableHeight = 12
	}
	m.table.PageSize(tableHeight)
}

// View renders the modal content with data table.
func (m *EventsModal) View() string {
	if !m.visible {
		return ""
	}

	// Nil theme guard.
	theme := m.theme
	if theme == nil {
		theme = themes.NewDefaultTheme()
	}

	// Calculate modal dimensions.
	maxModalHeight := 24
	terminalMaxHeight := int(float64(m.height) * 0.85)
	if terminalMaxHeight > maxModalHeight {
		maxModalHeight = terminalMaxHeight
	}
	if maxModalHeight < 16 {
		maxModalHeight = 16
	}

	modalWidth := m.width - 8
	if modalWidth > 95 {
		modalWidth = 95
	}
	if modalWidth < 60 {
		modalWidth = 60
	}

	// Build title using UIKit.
	title := primitives.Title(fmt.Sprintf("Events using %q (%d)", m.skillName, len(m.events)), theme).Render()

	// Build content - table or empty message.
	var content string
	if m.table.IsEmpty() {
		content = primitives.Muted("No events use this skill.", theme).
			Italic().
			MarginTop(2).
			MarginBottom(2).
			Render()
	} else {
		content = m.table.Render()
	}

	// Build footer with navigation info.
	var footerText string
	if !m.table.IsEmpty() {
		footerText = fmt.Sprintf("Enter: View Details | ↑↓/j/k: Navigate | Esc: Close  [%d/%d]",
			m.table.GetSelectedIndex()+1, m.table.Count())
	} else {
		footerText = "Esc: Close"
	}
	footer := primitives.Muted(footerText, theme).Render()

	// Build modal content.
	modalContent := lipgloss.JoinVertical(lipgloss.Left, title, "", content, "", footer)

	// Wrap in styled box with solid background using UIKit.
	return containers.NewBox(theme).
		Content(modalContent).
		Width(modalWidth).
		MaxHeight(maxModalHeight).
		Padding(2).
		Background(theme.BackgroundColor()).
		Render()
}

// SetDimensions updates the modal's available dimensions.
func (m *EventsModal) SetDimensions(width, height int) {
	m.width = width
	m.height = height
	m.updateTableDimensions()
}

// Show makes the modal visible.
func (m *EventsModal) Show() {
	m.visible = true
	m.selectedEvent = nil
	m.table.SetSelectedIndex(0)
}

// Hide hides the modal.
func (m *EventsModal) Hide() {
	m.visible = false
}

// IsVisible returns whether the modal is currently visible.
func (m *EventsModal) IsVisible() bool {
	return m.visible
}

// GetSkillID returns the skill ID this modal is showing events for.
func (m *EventsModal) GetSkillID() string {
	return m.skillID
}

// GetSkillName returns the skill name this modal is showing events for.
func (m *EventsModal) GetSkillName() string {
	return m.skillName
}

// SetEvents updates the events being displayed.
func (m *EventsModal) SetEvents(events []*career.CareerEvent) {
	// Filter out nil events.
	filteredEvents := make([]*career.CareerEvent, 0, len(events))
	for _, e := range events {
		if e != nil {
			filteredEvents = append(filteredEvents, e)
		}
	}

	m.events = filteredEvents
	m.selectedEvent = nil
	m.table.SetItems(filteredEvents)
}

// HasSelection returns true if the user selected an event.
func (m *EventsModal) HasSelection() bool {
	return m.selectedEvent != nil
}

// GetSelectedEvent returns the selected event (nil if none selected).
func (m *EventsModal) GetSelectedEvent() *career.CareerEvent {
	return m.selectedEvent
}

// ClearSelection clears any previous selection.
func (m *EventsModal) ClearSelection() {
	m.selectedEvent = nil
}

// GetSelectedIndex returns the current selection index.
func (m *EventsModal) GetSelectedIndex() int {
	return m.table.GetSelectedIndex()
}
