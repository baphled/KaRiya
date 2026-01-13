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

// buildForm creates the huh form with all fields (Text, Date, Company, Project, Tags, Categories)
func (m *EditEventModal) buildForm() {
	// Full edit: all fields available
	fieldsGroup := huh.NewGroup(
		forms.NewText(forms.FieldConfig{
			Key:         "text",
			Title:       "Event Description",
			Description: "What did you accomplish? (required, 10-2000 characters)",
			Placeholder: "Enter event description...",
			CharLimit:   2000,
			Validate:    forms.Compose(forms.Required, forms.MinLength(10), forms.MaxLength(2000)),
		}).Value(&m.formData.Text).Lines(4),

		forms.NewInput(forms.FieldConfig{
			Key:         "date",
			Title:       "Date",
			Description: "YYYY-MM-DD or 'today'",
			Placeholder: "2006-01-02",
			Validate:    forms.DateFormat,
		}).Value(&m.formData.Date),

		forms.NewInput(forms.FieldConfig{
			Key:         "company",
			Title:       "Company",
			Description: "Company or organization name (optional)",
			Placeholder: "Optional",
			Validate:    nil,
		}).Value(&m.formData.Company),

		forms.NewInput(forms.FieldConfig{
			Key:         "project",
			Title:       "Project",
			Description: "Project name (optional)",
			Placeholder: "Optional",
			Validate:    nil,
		}).Value(&m.formData.Project),

		forms.NewMultiSelect(
			"tags",
			"Tags",
			"Select applicable tags (use Space to select, Enter to confirm)",
			[]forms.SelectOption{
				{Key: "achievement", Value: "achievement"},
				{Key: "learning", Value: "learning"},
				{Key: "collaboration", Value: "collaboration"},
				{Key: "leadership", Value: "leadership"},
				{Key: "technical", Value: "technical"},
				{Key: "business", Value: "business"},
				{Key: "innovation", Value: "innovation"},
				{Key: "problem-solving", Value: "problem-solving"},
			},
			0, // no limit
		).Value(&m.formData.Tags),

		forms.NewMultiSelect(
			"categories",
			"Categories",
			"Select applicable categories (use Space to select, Enter to confirm)",
			[]forms.SelectOption{
				{Key: "development", Value: "development"},
				{Key: "design", Value: "design"},
				{Key: "management", Value: "management"},
				{Key: "research", Value: "research"},
				{Key: "testing", Value: "testing"},
				{Key: "documentation", Value: "documentation"},
				{Key: "deployment", Value: "deployment"},
				{Key: "maintenance", Value: "maintenance"},
			},
			0, // no limit
		).Value(&m.formData.Categories),
	)

	// Calculate form dimensions
	modalWidth := m.width - 10
	if modalWidth > 90 {
		modalWidth = 90
	}
	if modalWidth < 50 {
		modalWidth = 50
	}

	formHeight := 24
	if m.height < 30 {
		formHeight = m.height - 6
	}

	m.form = forms.NewFormWithFixedConfirm(fieldsGroup, &m.formData.SubmitConfirmed, modalWidth, formHeight)
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

	// Check if form is complete AND user confirmed (not cancelled)
	if m.form.State == huh.StateCompleted {
		m.visible = false
		// Check if user confirmed submission (not cancelled)
		if m.formData.SubmitConfirmed {
			// Return event data
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
		// User cancelled - close modal without returning data
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
