package skill

import (
	"encoding/json"
	"fmt"

	"github.com/baphled/kariya/internal/ui/behaviors"
	"github.com/baphled/kariya/internal/ui/display"
	"github.com/baphled/kariya/internal/ui/themes"
	"github.com/baphled/kariya/internal/ui/uikit/containers"
	"github.com/baphled/kariya/internal/ui/uikit/primitives"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
)

// maxEventTextLength defines the maximum length for event text before truncation.
// This ensures event descriptions fit within the table column width.
const maxEventTextLength = 42

// Events displays events that use a skill in a modal overlay with a data table.
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
//	modal := NewEvents(skillID, skillName, events, theme)
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
type Events struct {
	skillID       string
	skillName     string
	events        []display.Event
	theme         themes.Theme
	visible       bool
	width         int
	height        int
	table         *behaviors.TableBehavior[display.Event]
	selectedEvent *display.Event
}

// eventRowFormatter formats an event for table display.
func eventRowFormatter(event display.Event, _ int) []string {
	dateStr := "-"
	if !event.Date.IsZero() {
		dateStr = event.Date.Format("2006-01-02")
	}

	text := event.Text
	if text == "" {
		text = "(No description)"
	}
	if len(text) > maxEventTextLength {
		text = text[:maxEventTextLength] + "..."
	}

	company := event.Company
	if company == "" {
		company = "-"
	}

	return []string{dateStr, text, company}
}

func filterDisplayEvents(events []display.Event) []display.Event {
	filteredEvents := make([]display.Event, 0, len(events))
	for i := range events {
		event := events[i]
		if isEmptyDisplayEvent(event) {
			continue
		}
		filteredEvents = append(filteredEvents, event)
	}
	return filteredEvents
}

func isEmptyDisplayEvent(event display.Event) bool {
	return event.ID == "" &&
		event.Text == "" &&
		event.Date.IsZero() &&
		event.Company == "" &&
		event.Project == "" &&
		len(event.Tags) == 0 &&
		len(event.Categories) == 0 &&
		len(event.Skills) == 0
}

func normalizeDisplayEvents(rawEvents any) []display.Event {
	if rawEvents == nil {
		return nil
	}

	if events, ok := rawEvents.([]display.Event); ok {
		return filterDisplayEvents(append([]display.Event(nil), events...))
	}

	payload, err := json.Marshal(rawEvents)
	if err != nil {
		return nil
	}

	var events []display.Event
	if json.Unmarshal(payload, &events) != nil {
		return nil
	}

	return filterDisplayEvents(events)
}

// NewEvents creates a new events modal for a skill.
//
// Expected:
//   - Must be a valid string.
//   - event must be valid.
//   - th must be a valid theme instance (can be nil).
//
// Returns:
//   - A fully initialized Events ready for use.
//
// Side effects:
//   - None.
func NewEvents(skillID, skillName string, events any, theme themes.Theme) *Events {
	filteredEvents := normalizeDisplayEvents(events)

	if theme == nil {
		theme = themes.NewDefaultTheme()
	}

	columns := []behaviors.ColumnDef{
		{Title: "Date", Width: 12},
		{Title: "Event", Width: 45},
		{Title: "Company", Width: 20},
	}

	tableBehavior := behaviors.NewTableBehavior(theme, columns, eventRowFormatter).
		PageSize(12).
		EmptyMessage("No events use this skill.").
		HidePagination()

	tableBehavior.SetItems(filteredEvents)

	m := &Events{
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
//
// Returns:
//   - A tea.Cmd value.
//
// Side effects:
//   - None.
func (m *Events) Init() tea.Cmd {
	return nil
}

// Update handles keyboard input and window sizing.
//
// Expected:
//   - msg must be a valid tea.Msg type.
//
// Returns:
//   - tea.Model: the updated model.
//   - tea.Cmd: command to execute.
//
// Side effects:
//   - May update dimensions.
//   - May hide modal.
//   - May set selected event.
func (m *Events) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
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
		// Handle special keys.
		switch msg.Type {
		case tea.KeyEsc, tea.KeyBackspace:
			// Close modal without selection.
			m.Hide()
			return m, nil

		case tea.KeyEnter:
			if selected := m.table.GetSelectedItem(); selected != nil {
				event := *selected
				m.selectedEvent = &event
				m.Hide()
			}
			return m, nil
		}

		// Handle character keys.
		if msg.String() == "q" {
			m.Hide()
			return m, nil
		}

		// Handle table navigation.
		if m.table.HandleNavigation(msg.String()) {
			return m, nil
		}
	}

	return m, nil
}

// updateTableDimensions updates the table dimensions based on modal size.
func (m *Events) updateTableDimensions() {
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
//
// Returns:
//   - A string value.
//
// Side effects:
//   - None.
func (m *Events) View() string {
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
//
// Expected:
//   - int must be valid.
//
// Side effects:
//   - None.
func (m *Events) SetDimensions(width, height int) {
	m.width = width
	m.height = height
	m.updateTableDimensions()
}

// Show makes the modal visible.
//
// Side effects:
//   - None.
func (m *Events) Show() {
	m.visible = true
	m.selectedEvent = nil
	m.table.SetSelectedIndex(0)
}

// Hide hides the modal.
//
// Side effects:
//   - None.
func (m *Events) Hide() {
	m.visible = false
}

// IsVisible returns whether the modal is currently visible.
//
// Returns:
//   - A bool value.
//
// Side effects:
//   - None.
func (m *Events) IsVisible() bool {
	return m.visible
}

// GetSkillID returns the skill ID this modal is showing events for.
//
// Returns:
//   - A string value.
//
// Side effects:
//   - None.
func (m *Events) GetSkillID() string {
	return m.skillID
}

// GetSkillName returns the skill name this modal is showing events for.
//
// Returns:
//   - A string value.
//
// Side effects:
//   - None.
func (m *Events) GetSkillName() string {
	return m.skillName
}

// SetEvents updates the events being displayed.
//
// Expected:
//   - event must be valid.
//
// Side effects:
//   - None.
func (m *Events) SetEvents(events []display.Event) {
	filteredEvents := filterDisplayEvents(append([]display.Event(nil), events...))
	m.events = filteredEvents
	m.selectedEvent = nil
	m.table.SetItems(filteredEvents)
	m.table.SetSelectedIndex(0)
}

// HasSelection returns true if the user selected an event.
//
// Returns:
//   - A bool value.
//
// Side effects:
//   - None.
func (m *Events) HasSelection() bool {
	return m.selectedEvent != nil
}

// GetSelectedEvent returns the selected event (nil if none selected).
//
// Returns:
//   - A fully initialized display.Event ready for use.
//
// Side effects:
//   - None.
func (m *Events) GetSelectedEvent() *display.Event {
	return m.selectedEvent
}

// ClearSelection clears any previous selection.
//
// Side effects:
//   - None.
func (m *Events) ClearSelection() {
	m.selectedEvent = nil
}

// GetSelectedIndex returns the current selection index.
//
// Returns:
//   - A int value.
//
// Side effects:
//   - None.
func (m *Events) GetSelectedIndex() int {
	return m.table.GetSelectedIndex()
}
