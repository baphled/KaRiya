package components

import (
	"github.com/baphled/kariya/internal/cli/forms"
	"github.com/baphled/kariya/internal/domain/career"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/huh"
)

// EditEventModal provides a way to edit an existing career event.
// It shows a full form with all fields pre-populated from the existing event.
//
// Usage:
//
//	modal := components.NewEditEventModal(existingEvent, width, height)
//	cmd := modal.Init()
//	// In Update:
//	cmd, completed, eventData := modal.Update(msg)
//	if completed && eventData != nil {
//	    // User completed form - update event
//	    updatedEvent := eventData.ToCareerEvent(existingEvent.ID)
//	} else if !modal.IsVisible() {
//	    // User cancelled (Esc)
//	}
type EditEventModal struct {
	form          *huh.Form
	formData      *forms.CaptureEventFormData
	originalEvent *career.CareerEvent
	visible       bool
	width         int
	height        int
}

// NewEditEventModal creates a new edit event modal with fields pre-populated from the existing event.
// event: the existing event to edit
// width, height: terminal dimensions for responsive sizing
func NewEditEventModal(event *career.CareerEvent, width, height int) *EditEventModal {
	// Pre-populate form data from existing event
	formData := &forms.CaptureEventFormData{
		Text:            event.Text,
		Date:            event.Date.Format("2006-01-02"),
		Company:         event.Company,
		Project:         event.Project,
		Tags:            event.Tags,
		Categories:      event.Categories,
		SubmitConfirmed: false,
	}

	modal := &EditEventModal{
		formData:      formData,
		originalEvent: event,
		visible:       true,
		width:         width,
		height:        height,
	}

	modal.buildForm()
	return modal
}

// buildForm creates the huh form with confirm button
func (m *EditEventModal) buildForm() {
	// Calculate form dimensions
	modalWidth := m.width - 10
	if modalWidth > 90 {
		modalWidth = 90
	}
	if modalWidth < 50 {
		modalWidth = 50
	}

	// Calculate maximum form height to fit within terminal
	// Total overhead:
	// - Logo: 6 lines
	// - Logo spacing: 2 lines
	// - Footer: 4 lines
	// - Modal chrome (borders, padding, title, footer): 8 lines
	// - Safety margins: 4 lines
	// Total: 24 lines overhead
	const overhead = 24
	maxFormHeight := m.height - overhead
	if maxFormHeight < 12 {
		maxFormHeight = 12 // Minimum usable height
	}

	// Edit form has more fields (text, date, company, project, tags, categories)
	// Ideally needs ~25 lines, but must fit within terminal constraints and will scroll if needed
	formHeight := 25
	if formHeight > maxFormHeight {
		formHeight = maxFormHeight
	}

	// Use standard form with confirm button (Submit/Cancel)
	// The form will be scrollable if content exceeds formHeight
	m.form = forms.NewCaptureEventForm(m.formData, "manual", modalWidth, formHeight)
}

// Init initializes the modal and its form.
func (m *EditEventModal) Init() tea.Cmd {
	if m.form == nil {
		return nil
	}
	return m.form.Init()
}

// Update handles messages for the edit event modal.
//
// Returns:
//   - tea.Cmd: command to execute
//   - bool: true if form completed successfully
//   - *EditEventData: event data if completed, nil otherwise
func (m *EditEventModal) Update(msg tea.Msg) (tea.Cmd, bool, *EditEventData) {
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
			eventData := &EditEventData{
				Text:       m.formData.Text,
				Date:       m.formData.Date,
				Company:    m.formData.Company,
				Project:    m.formData.Project,
				Tags:       m.formData.Tags,
				Categories: m.formData.Categories,
			}
			return cmd, true, eventData
		}
		// User cancelled - close without returning data
		return cmd, false, nil
	}

	return cmd, false, nil
}

// View renders the edit event modal
func (m *EditEventModal) View() string {
	if !m.visible {
		return ""
	}
	return m.form.View()
}

// IsVisible returns whether the modal is currently visible
func (m *EditEventModal) IsVisible() bool {
	return m.visible
}

// Show makes the modal visible
func (m *EditEventModal) Show() {
	m.visible = true
}

// Hide hides the modal
func (m *EditEventModal) Hide() {
	m.visible = false
}

// GetOriginalEvent returns the original event being edited (useful for comparison)
func (m *EditEventModal) GetOriginalEvent() *career.CareerEvent {
	return m.originalEvent
}

// EditEventData holds the full data from editing an event.
type EditEventData struct {
	Text       string
	Date       string
	Company    string
	Project    string
	Tags       []string
	Categories []string
}

// ToCareerEvent converts the form data to a CareerEvent domain object.
// eventID: the ID of the event being updated (preserved from original)
// Returns an updated CareerEvent ready to be saved.
func (d *EditEventData) ToCareerEvent(eventID string, createdAt, updatedAt interface{}) *career.CareerEvent {
	// Parse date or default to original
	eventDate, err := forms.ParseDateString(d.Date)
	if err != nil {
		// If parse fails, this should have been caught by validation
		// But as a fallback, use the provided date string as-is and try standard parse
		eventDate, _ = forms.ParseDateString("today")
	}

	event := &career.CareerEvent{
		ID:         eventID,
		Text:       d.Text,
		Date:       eventDate,
		Company:    d.Company,
		Project:    d.Project,
		Tags:       d.Tags,
		Categories: d.Categories,
		Skills:     []string{}, // Skills managed separately
	}

	// Preserve CreatedAt, update UpdatedAt
	if ca, ok := createdAt.(interface{ Format(string) string }); ok {
		if parsed, err := forms.ParseDateString(ca.Format("2006-01-02")); err == nil {
			event.CreatedAt = parsed
		}
	}
	if ua, ok := updatedAt.(interface{ Format(string) string }); ok {
		if parsed, err := forms.ParseDateString(ua.Format("2006-01-02")); err == nil {
			event.UpdatedAt = parsed
		}
	}

	return event
}
