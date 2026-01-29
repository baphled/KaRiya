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

// QuickAddModal provides a quick way to add a new career event.
// It shows a minimal form with just Text, Date, and Company fields.
//
// Usage:
//
//	modal := modals.NewQuickAddModal(width, height)
//	cmd := modal.Init()
//	// In Update:
//	cmd, completed, eventData := modal.Update(msg)
//	if completed && eventData != nil {
//	    // User completed form - create event
//	    newEvent := eventData.ToCareerEvent()
//	} else if !modal.IsVisible() {
//	    // User cancelled (Esc)
//	}
type QuickAddModal struct {
	form     *huh.Form
	formData *forms.CaptureEventFormData
	visible  bool
	width    int
	height   int
}

// NewQuickAddModal creates a new quick add event modal.
// width, height: terminal dimensions for responsive sizing
func NewQuickAddModal(width, height int) *QuickAddModal {
	formData := forms.NewCaptureEventFormData()
	// Default date to today.
	formData.Date = time.Now().Format("2006-01-02")

	modal := &QuickAddModal{
		formData: formData,
		visible:  true,
		width:    width,
		height:   height,
	}

	modal.buildForm()
	return modal
}

// buildForm creates the huh form with confirm button.
func (m *QuickAddModal) buildForm() {
	// Calculate form width.
	modalWidth := m.width - 10
	if modalWidth > 80 {
		modalWidth = 80
	}
	if modalWidth < 40 {
		modalWidth = 40
	}

	// Let Huh use natural height - bubbletea-overlay will handle positioning.
	formHeight := 0

	// Use standard form with confirm button (Submit/Cancel).
	m.form = forms.NewCaptureEventForm(m.formData, "quick", modalWidth, formHeight)
}

// Init initializes the modal and its form.
func (m *QuickAddModal) Init() tea.Cmd {
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
//   - *QuickAddData: event data if completed, nil otherwise
func (m *QuickAddModal) Update(msg tea.Msg) (tea.Cmd, bool, *QuickAddData) {
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
			eventData := &QuickAddData{
				Text: m.formData.Text,
				Date: m.formData.Date,
			}
			return cmd, true, eventData
		}
		// User cancelled - close without returning data.
		return cmd, false, nil
	}

	return cmd, false, nil
}

// View renders the quick add event modal with proper chrome (border, background)
// for overlay compositing. The chrome provides a solid background so the modal
// doesn't show the background layer through.
func (m *QuickAddModal) View() string {
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
func (m *QuickAddModal) IsVisible() bool {
	return m.visible
}

// Show makes the modal visible.
func (m *QuickAddModal) Show() {
	m.visible = true
}

// Hide hides the modal.
func (m *QuickAddModal) Hide() {
	m.visible = false
}

// QuickAddData holds the minimal data needed to create a new event.
// Quick events only capture text and date - other metadata can be added via Edit.
type QuickAddData struct {
	Text string
	Date string
}

// ToCareerEvent converts the form data to an Event domain object.
// Returns a new Event ready to be saved.
func (d *QuickAddData) ToCareerEvent() *career.Event {
	// Parse date or default to today.
	eventDate, err := forms.ParseDateString(d.Date)
	if err != nil {
		eventDate = time.Now()
	}

	return &career.Event{
		Text:       d.Text,
		Date:       eventDate,
		Company:    "",
		Project:    "",
		Tags:       []string{},
		Categories: []string{},
		Skills:     []string{},
		CreatedAt:  time.Now(),
		UpdatedAt:  time.Now(),
	}
}
