package components

import (
	"fmt"

	"github.com/baphled/kariya/internal/cli/themes"
	"github.com/baphled/kariya/internal/cli/uikit/containers"
	"github.com/baphled/kariya/internal/cli/uikit/primitives"
	"github.com/baphled/kariya/internal/domain/career"
	"github.com/charmbracelet/bubbles/table"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
)

// ViewSkillEventsModal displays events that use a skill in a modal overlay with a data table.
// This modal uses bubbles/table for consistent table display and navigation.
//
// Features:
// - Data table display with Date, Event, and Company columns
// - Standard table navigation (j/k, up/down, page up/down)
// - Press Enter to select an event for detailed view
// - Close with Escape, backspace, or 'q'
//
// Usage:
//
//	modal := NewViewSkillEventsModal(skillID, skillName, events, theme)
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
type ViewSkillEventsModal struct {
	skillID       string
	skillName     string
	events        []*career.CareerEvent
	theme         themes.Theme
	visible       bool
	width         int
	height        int
	table         table.Model
	selectedEvent *career.CareerEvent // Set when user selects an event
}

// NewViewSkillEventsModal creates a new events modal for a skill.
func NewViewSkillEventsModal(skillID, skillName string, events []*career.CareerEvent, theme themes.Theme) *ViewSkillEventsModal {
	// Filter out nil events
	filteredEvents := make([]*career.CareerEvent, 0, len(events))
	for _, e := range events {
		if e != nil {
			filteredEvents = append(filteredEvents, e)
		}
	}

	m := &ViewSkillEventsModal{
		skillID:   skillID,
		skillName: skillName,
		events:    filteredEvents,
		theme:     theme,
		visible:   false,
		width:     100,
		height:    24,
	}
	m.initTable()
	return m
}

// initTable initializes the bubbles/table with columns, rows, and styling
func (m *ViewSkillEventsModal) initTable() {
	// Nil theme guard
	theme := m.theme
	if theme == nil {
		theme = themes.NewDefaultTheme()
	}

	columns := []table.Column{
		{Title: "Date", Width: 12},
		{Title: "Event", Width: 45},
		{Title: "Company", Width: 20},
	}

	rows := m.buildRows()

	// Calculate table height (modal height - borders - padding - title - footer - spacing)
	tableHeight := m.height - 14
	if tableHeight < 5 {
		tableHeight = 5
	}
	if tableHeight > 12 {
		tableHeight = 12
	}

	t := table.New(
		table.WithColumns(columns),
		table.WithRows(rows),
		table.WithFocused(true),
		table.WithHeight(tableHeight),
	)

	// Apply theme styling consistent with other tables in the app
	s := table.DefaultStyles()
	s.Header = s.Header.
		BorderStyle(lipgloss.NormalBorder()).
		BorderForeground(theme.BorderColor()).
		BorderBottom(true).
		Bold(true).
		Foreground(theme.PrimaryColor())
	s.Selected = s.Selected.
		Foreground(theme.PrimaryColor()).
		Background(theme.AccentColor()).
		Bold(true)
	s.Cell = s.Cell.
		Foreground(theme.SecondaryColor())

	t.SetStyles(s)
	m.table = t
}

// buildRows converts events to table rows
func (m *ViewSkillEventsModal) buildRows() []table.Row {
	rows := make([]table.Row, 0, len(m.events))
	for _, event := range m.events {
		if event == nil {
			continue
		}

		// Date
		dateStr := "-"
		if !event.Date.IsZero() {
			dateStr = event.Date.Format("2006-01-02")
		}

		// Event text (truncated to fit column)
		text := event.Text
		if text == "" {
			text = "(No description)"
		}
		if len(text) > 42 {
			text = text[:42] + "..."
		}

		// Company
		company := event.Company
		if company == "" {
			company = "-"
		}

		rows = append(rows, table.Row{dateStr, text, company})
	}
	return rows
}

// Init initializes the modal (implements tea.Model for bubbletea-overlay).
func (m *ViewSkillEventsModal) Init() tea.Cmd {
	return nil
}

// Update handles keyboard input and window sizing.
func (m *ViewSkillEventsModal) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	if !m.visible {
		return m, nil
	}

	var cmd tea.Cmd

	switch msg := msg.(type) {
	case tea.WindowSizeMsg:
		m.width = msg.Width
		m.height = msg.Height
		m.initTable() // Rebuild table with new dimensions
		return m, nil

	case tea.KeyMsg:
		switch msg.String() {
		case "esc", "backspace", "q":
			// Close modal without selection
			m.Hide()
			return m, nil

		case "enter":
			// Select the current event
			idx := m.table.Cursor()
			if len(m.events) > 0 && idx >= 0 && idx < len(m.events) {
				m.selectedEvent = m.events[idx]
				m.Hide()
			}
			return m, nil
		}
	}

	// Forward other messages to table for navigation (up/down/j/k/pgup/pgdn/home/end/g/G)
	m.table, cmd = m.table.Update(msg)
	return m, cmd
}

// View renders the modal content with data table.
func (m *ViewSkillEventsModal) View() string {
	if !m.visible {
		return ""
	}

	// Nil theme guard
	theme := m.theme
	if theme == nil {
		theme = themes.NewDefaultTheme()
	}

	// Calculate modal dimensions
	// We need space for: border(2) + padding(2) + title(1) + blank(1) + table + blank(1) + footer(1)
	// Minimum chrome = 8 lines
	maxModalHeight := 24
	terminalMaxHeight := int(float64(m.height) * 0.85) // Use 85% of terminal height
	if terminalMaxHeight > maxModalHeight {
		maxModalHeight = terminalMaxHeight
	}
	if maxModalHeight < 16 {
		maxModalHeight = 16 // Minimum usable height
	}

	modalWidth := m.width - 8 // Leave margins
	if modalWidth > 95 {
		modalWidth = 95 // Max width for table
	}
	if modalWidth < 60 {
		modalWidth = 60 // Minimum usable width
	}

	// Build title using UIKit
	title := primitives.Title(fmt.Sprintf("Events using \"%s\" (%d)", m.skillName, len(m.events)), theme).Render()

	// Build content - either table or empty message
	var content string
	if len(m.events) == 0 {
		// Use UIKit Text with margin for empty state
		content = primitives.Muted("No events use this skill.", theme).
			Italic().
			MarginTop(2).
			MarginBottom(2).
			Render()
	} else {
		content = m.table.View()
	}

	// Build footer with pagination info using UIKit
	var footerText string
	if len(m.events) > 0 {
		footerText = fmt.Sprintf("Enter: View Details | ↑↓/j/k: Navigate | Esc: Close  [%d/%d]",
			m.table.Cursor()+1, len(m.events))
	} else {
		footerText = "Esc: Close"
	}
	footer := primitives.Muted(footerText, theme).Render()

	// Build modal content
	modalContent := lipgloss.JoinVertical(lipgloss.Left, title, "", content, "", footer)

	// Wrap in styled box with solid background using UIKit (with MaxHeight)
	return containers.NewBox(theme).
		Content(modalContent).
		Width(modalWidth).
		MaxHeight(maxModalHeight).
		Padding(2).
		Background(theme.BackgroundColor()).
		Render()
}

// SetDimensions updates the modal's available dimensions.
func (m *ViewSkillEventsModal) SetDimensions(width, height int) {
	m.width = width
	m.height = height
	m.initTable() // Rebuild table with new dimensions
}

// Show makes the modal visible.
func (m *ViewSkillEventsModal) Show() {
	m.visible = true
	m.selectedEvent = nil
	m.table.SetCursor(0) // Reset selection
}

// Hide hides the modal.
func (m *ViewSkillEventsModal) Hide() {
	m.visible = false
}

// IsVisible returns whether the modal is currently visible.
func (m *ViewSkillEventsModal) IsVisible() bool {
	return m.visible
}

// GetSkillID returns the skill ID this modal is showing events for.
func (m *ViewSkillEventsModal) GetSkillID() string {
	return m.skillID
}

// GetSkillName returns the skill name this modal is showing events for.
func (m *ViewSkillEventsModal) GetSkillName() string {
	return m.skillName
}

// SetEvents updates the events being displayed.
func (m *ViewSkillEventsModal) SetEvents(events []*career.CareerEvent) {
	// Filter out nil events
	filteredEvents := make([]*career.CareerEvent, 0, len(events))
	for _, e := range events {
		if e != nil {
			filteredEvents = append(filteredEvents, e)
		}
	}

	m.events = filteredEvents
	m.selectedEvent = nil
	m.initTable() // Rebuild table with new data
}

// HasSelection returns true if the user selected an event.
func (m *ViewSkillEventsModal) HasSelection() bool {
	return m.selectedEvent != nil
}

// GetSelectedEvent returns the selected event (nil if none selected).
func (m *ViewSkillEventsModal) GetSelectedEvent() *career.CareerEvent {
	return m.selectedEvent
}

// ClearSelection clears any previous selection.
func (m *ViewSkillEventsModal) ClearSelection() {
	m.selectedEvent = nil
}

// GetSelectedIndex returns the current selection index.
func (m *ViewSkillEventsModal) GetSelectedIndex() int {
	return m.table.Cursor()
}
