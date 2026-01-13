package components

import (
	"time"

	"github.com/baphled/kariya/internal/cli/forms"
	"github.com/baphled/kariya/internal/domain/career"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/huh"
)

// QuickAddEventModal provides a quick way to add a new career event.
// It shows a minimal form with just Text, Date, and Company fields.
//
// Usage:
//
//	modal := components.NewQuickAddEventModal(width, height)
//	cmd := modal.Init()
//	// In Update:
//	cmd, completed, eventData := modal.Update(msg)
//	if completed && eventData != nil {
//	    // User completed form - create event
//	    newEvent := eventData.ToCareerEvent()
//	} else if !modal.IsVisible() {
//	    // User cancelled (Esc)
//	}
type QuickAddEventModal struct {
	form     *huh.Form
	formData *forms.CaptureEventFormData
	visible  bool
	width    int
	height   int
}

// NewQuickAddEventModal creates a new quick add event modal.
// width, height: terminal dimensions for responsive sizing
func NewQuickAddEventModal(width, height int) *QuickAddEventModal {
	formData := forms.NewCaptureEventFormData()
	// Default date to today
	formData.Date = time.Now().Format("2006-01-02")

	modal := &QuickAddEventModal{
		formData: formData,
		visible:  true,
		width:    width,
		height:   height,
	}

	modal.buildForm()
	return modal
}

// buildForm creates the huh form with confirm button
func (m *QuickAddEventModal) buildForm() {
	// Calculate form dimensions
	modalWidth := m.width - 10
	if modalWidth > 80 {
		modalWidth = 80
	}
	if modalWidth < 40 {
		modalWidth = 40
	}

	formHeight := 18
	if m.height < 24 {
		formHeight = m.height - 6
	}

	// Use standard form with confirm button (Submit/Cancel)
	m.form = forms.NewCaptureEventForm(m.formData, "quick", modalWidth, formHeight)
}

// Init initializes the modal and its form.
func (m *QuickAddEventModal) Init() tea.Cmd {
	if m.form == nil {
		return nil
	}
	return m.form.Init()
}

// Update handles messages for the quick add event modal.
//
// Returns:
//   - tea.Cmd: command to execute
//   - bool: true if form completed successfully
//   - *QuickAddEventData: event data if completed, nil otherwise
func (m *QuickAddEventModal) Update(msg tea.Msg) (tea.Cmd, bool, *QuickAddEventData) {
	if !m.visible {
		return nil, false, nil
	}

	switch msg := msg.(type) {
	case tea.WindowSizeMsg:
		m.width = msg.Width
		m.height = msg.Height
		// Rebuild form with new dimensions
		m.buildForm()
		return m.form.Init(), false, nil

	case tea.KeyMsg:
		switch msg.String() {
		case "esc":
			// Close modal without saving
			m.visible = false
			return nil, false, nil
		}
	}

	// Update form
	form, cmd := m.form.Update(msg)
	m.form = form.(*huh.Form)

	// Check if form is complete AND user confirmed submission
	if m.form.State == huh.StateCompleted {
		m.visible = false
		// Only return data if user confirmed (pressed Submit, not Cancel)
		if m.formData.SubmitConfirmed {
			eventData := &QuickAddEventData{
				Text:    m.formData.Text,
				Date:    m.formData.Date,
				Company: m.formData.Company,
			}
			return cmd, true, eventData
		}
		// User cancelled - close without returning data
		return cmd, false, nil
	}

	return cmd, false, nil
}

// View renders the quick add event modal
func (m *QuickAddEventModal) View() string {
	if !m.visible {
		return ""
	}
	return m.form.View()
}

// IsVisible returns whether the modal is currently visible
func (m *QuickAddEventModal) IsVisible() bool {
	return m.visible
}

// Show makes the modal visible
func (m *QuickAddEventModal) Show() {
	m.visible = true
}

// Hide hides the modal
func (m *QuickAddEventModal) Hide() {
	m.visible = false
}

// QuickAddEventData holds the minimal data needed to create a new event.
type QuickAddEventData struct {
	Text    string
	Date    string
	Company string
}

// ToCareerEvent converts the form data to a CareerEvent domain object.
// Returns a new CareerEvent ready to be saved.
func (d *QuickAddEventData) ToCareerEvent() *career.CareerEvent {
	// Parse date or default to today
	eventDate, err := forms.ParseDateString(d.Date)
	if err != nil {
		eventDate = time.Now()
	}

	return &career.CareerEvent{
		Text:       d.Text,
		Date:       eventDate,
		Company:    d.Company,
		Project:    "",
		Tags:       []string{},
		Categories: []string{},
		Skills:     []string{},
		CreatedAt:  time.Now(),
		UpdatedAt:  time.Now(),
	}
}
