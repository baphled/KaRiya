package modals

import (
	"time"

	"github.com/baphled/kariya/internal/cli/forms"
	"github.com/baphled/kariya/internal/cli/themes"
	"github.com/baphled/kariya/internal/cli/uikit/containers"
	"github.com/baphled/kariya/internal/domain/career"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/huh"
)

// EditModal provides a way to edit an existing career event.
// It shows a full form with all fields pre-populated from the existing event.
//
// Usage:
//
//	modal := modals.NewEditModal(existingEvent, width, height)
//	cmd := modal.Init()
//	// In Update:
//	cmd, completed, eventData := modal.Update(msg)
//	if completed && eventData != nil {
//	    // User completed form - update event
//	    updatedEvent := eventData.ToCareerEvent(existingEvent.ID)
//	} else if !modal.IsVisible() {
//	    // User cancelled (Esc)
//	}
type EditModal struct {
	form          *huh.Form
	formData      *forms.CaptureEventFormData
	originalEvent *career.Event
	visible       bool
	width         int
	height        int
}

// NewEditModal creates a new edit event modal with fields pre-populated from the existing event.
// event: the existing event to edit
// width, height: terminal dimensions for responsive sizing
func NewEditModal(event *career.Event, width, height int) *EditModal {
	// Pre-populate form data from existing event.
	formData := &forms.CaptureEventFormData{
		Text:            event.Text,
		Date:            event.Date.Format("2006-01-02"),
		Company:         event.Company,
		Project:         event.Project,
		Tags:            event.Tags,
		Categories:      event.Categories,
		SubmitConfirmed: false,
	}

	modal := &EditModal{
		formData:      formData,
		originalEvent: event,
		visible:       true,
		width:         width,
		height:        height,
	}

	modal.buildForm()
	return modal
}

// buildForm creates the huh form with confirm button.
func (m *EditModal) buildForm() {
	// Calculate form width.
	modalWidth := m.width - 10
	if modalWidth > 90 {
		modalWidth = 90
	}
	if modalWidth < 50 {
		modalWidth = 50
	}

	// Let Huh use natural height.
	formHeight := 0

	// Use standard form with confirm button (Submit/Cancel).
	m.form = forms.NewCaptureEventForm(m.formData, "manual", modalWidth, formHeight)
}

// Init initializes the modal and its form.
func (m *EditModal) Init() tea.Cmd {
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
//   - *EditData: event data if completed, nil otherwise
func (m *EditModal) Update(msg tea.Msg) (tea.Cmd, bool, *EditData) {
	if !m.visible {
		return nil, false, nil
	}

	switch msg := msg.(type) {
	case tea.WindowSizeMsg:
		m.width = msg.Width
		m.height = msg.Height
		// Rebuild form with new dimensions.
		m.buildForm()
		return m.form.Init(), false, nil

	case tea.KeyMsg:
		switch msg.String() {
		case "esc":
			// Close modal without saving.
			m.visible = false
			return nil, false, nil
		}
	}

	// Update form.
	form, cmd := m.form.Update(msg)
	m.form = form.(*huh.Form)

	// Check if form is complete AND user confirmed submission.
	if m.form.State == huh.StateCompleted {
		m.visible = false
		// Only return data if user confirmed (pressed Submit, not Cancel).
		if m.formData.SubmitConfirmed {
			eventData := &EditData{
				Text:       m.formData.Text,
				Date:       m.formData.Date,
				Company:    m.formData.Company,
				Project:    m.formData.Project,
				Tags:       m.formData.Tags,
				Categories: m.formData.Categories,
			}
			return cmd, true, eventData
		}
		// User cancelled - close without returning data.
		return cmd, false, nil
	}

	return cmd, false, nil
}

// View renders the edit event modal with proper chrome (border, background)
// for overlay compositing. The chrome provides a solid background so the modal
// doesn't show the background layer through.
func (m *EditModal) View() string {
	if !m.visible {
		return ""
	}

	// Wrap the form in a styled box with solid background using UIKit.
	theme := themes.NewDefaultTheme()
	return containers.NewBox(theme).
		Content(m.form.View()).
		Padding(2).
		Background(theme.BackgroundColor()).
		Render()
}

// IsVisible returns whether the modal is currently visible.
func (m *EditModal) IsVisible() bool {
	return m.visible
}

// Show makes the modal visible.
func (m *EditModal) Show() {
	m.visible = true
}

// Hide hides the modal.
func (m *EditModal) Hide() {
	m.visible = false
}

// GetOriginalEvent returns the original event being edited (useful for comparison).
func (m *EditModal) GetOriginalEvent() *career.Event {
	return m.originalEvent
}

// EditData holds the full data from editing an event.
type EditData struct {
	Text       string
	Date       string
	Company    string
	Project    string
	Tags       []string
	Categories []string
}

// ToCareerEvent converts the form data to a Event domain object.
// eventID: the ID of the event being updated (preserved from original)
// Returns an updated Event ready to be saved.
func (d *EditData) ToCareerEvent(eventID string, createdAt, updatedAt interface{}) *career.Event {
	// Parse date or default to original.
	eventDate, err := forms.ParseDateString(d.Date)
	if err != nil {
		// If parse fails, this should have been caught by validation.
		// But as a fallback, try "today" or use current time.
		eventDate, err = forms.ParseDateString("today")
		if err != nil {
			eventDate = time.Now()
		}
	}

	event := &career.Event{
		ID:         eventID,
		Text:       d.Text,
		Date:       eventDate,
		Company:    d.Company,
		Project:    d.Project,
		Tags:       d.Tags,
		Categories: d.Categories,
		Skills:     []string{},
	}

	// Preserve CreatedAt, update UpdatedAt.
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
