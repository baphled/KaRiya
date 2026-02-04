package modals

import (
	"time"

	"github.com/baphled/kariya/internal/cli/forms"
	"github.com/baphled/kariya/internal/cli/themes"
	"github.com/baphled/kariya/internal/cli/uikit/containers"
	"github.com/baphled/kariya/internal/cli/uikit/primitives"
	"github.com/baphled/kariya/internal/domain/career"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/huh"
	"github.com/charmbracelet/lipgloss"
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
//
// Expected:
//   - int must be valid.
//
// Returns:
//   - A fully initialized QuickAddModal ready for use.
//
// Side effects:
//   - None.
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

// buildForm creates the huh form with confirm button using LayoutStack.
func (m *QuickAddModal) buildForm() {
	modalWidth := m.width - 10
	if modalWidth > 80 {
		modalWidth = 80
	}
	if modalWidth < 40 {
		modalWidth = 40
	}

	formWidth := forms.ModalFormWidth(modalWidth)
	formHeight := forms.ModalFormHeight(m.height)

	m.form = forms.NewCaptureEventForm(m.formData, "quick", formWidth, formHeight)
}

// Init initializes the modal and its form.
//
// Returns:
//   - A tea.Cmd value.
//
// Side effects:
//   - None.
func (m *QuickAddModal) Init() tea.Cmd {
	if m.form == nil {
		return nil
	}
	return m.form.Init()
}

// Update handles messages for the quick add event modal.
//
// Expected:
//   - msg must be a valid tea.Msg type.
//
// Returns:
//   - tea.Cmd: command to execute.
//   - bool: true if form completed successfully.
//   - *QuickAddData: event data if completed, nil otherwise.
//
// Side effects:
//   - May hide modal on completion or cancellation.
//   - May rebuild form on window resize.
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
		switch msg.Type {
		case tea.KeyEsc:
			m.visible = false
			return nil, false, nil

		case tea.KeyCtrlS:
			m.visible = false
			m.formData.SubmitConfirmed = true
			eventData := &QuickAddData{
				Text: m.formData.Text,
				Date: m.formData.Date,
			}
			return nil, true, eventData
		}
	}

	// Update form.
	form, cmd := m.form.Update(msg)
	if f, ok := form.(*huh.Form); ok {
		m.form = f
	}

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
//
// Returns:
//   - A string value.
//
// Side effects:
//   - None.
func (m *QuickAddModal) View() string {
	if !m.visible {
		return ""
	}

	theme := themes.NewDefaultTheme()

	helpFooter := primitives.RenderHelpFooter(theme,
		primitives.NextFieldBadge(theme),
		primitives.HelpKeyBadge("Ctrl+S", "Submit", theme),
		primitives.CancelBadge(theme),
	)

	content := lipgloss.JoinVertical(lipgloss.Left,
		m.form.View(),
		"",
		helpFooter,
	)

	return containers.NewBox(theme).
		Content(content).
		Padding(2).
		Background(theme.BackgroundColor()).
		Render()
}

// IsVisible returns whether the modal is currently visible.
//
// Returns:
//   - A bool value.
//
// Side effects:
//   - None.
func (m *QuickAddModal) IsVisible() bool {
	return m.visible
}

// Show makes the modal visible.
//
// Side effects:
//   - None.
func (m *QuickAddModal) Show() {
	m.visible = true
}

// Hide hides the modal.
//
// Side effects:
//   - None.
func (m *QuickAddModal) Hide() {
	m.visible = false
}

// QuickAddData holds the minimal data needed to create a new event.
// Quick events only capture text and date - other metadata can be added via Edit.
type QuickAddData struct {
	Text string
	Date string
}

// ToCareerEvent converts the form data to a Event domain object.
//
// Returns:
//   - A fully initialized career.Event ready for use.
//
// Side effects:
//   - None.
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
