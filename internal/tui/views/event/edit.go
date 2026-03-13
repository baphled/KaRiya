package event

import (
	"github.com/baphled/kariya/internal/tui/forms"
	"github.com/baphled/kariya/internal/tui/views/shared"
	"github.com/baphled/kariya/internal/ui/display"
	"github.com/baphled/kariya/internal/ui/themes"
	"github.com/baphled/kariya/internal/ui/uikit/containers"
	"github.com/baphled/kariya/internal/ui/uikit/primitives"
	"github.com/baphled/kariya/internal/ui/uikit/widgets"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
)

// Edit provides a way to edit an existing career event.
// It shows a full form with all fields pre-populated from the existing event.
//
// Usage:
//
//	view := eventviews.NewEdit(existingEvent, width, height)
//	cmd := modal.Init()
//	// In Update:
//	cmd, completed, eventData := modal.Update(msg)
//	if completed && eventData != nil {
//	    // User completed form - pass EditData to a domain conversion function
//	} else if !modal.IsVisible() {
//	    // User cancelled (Esc)
//	}
type Edit struct {
	form          forms.Form
	formData      *forms.CaptureEventFormData
	originalEvent display.Event
	visible       bool
	width         int
	height        int
}

// NewEdit creates a new edit event modal with fields pre-populated from the existing event.
//
// Expected:
//   - event must be valid.
//   - int must be valid.
//
// Returns:
//   - A fully initialized Edit ready for use.
//
// Side effects:
//   - None.
func NewEdit(event display.Event, width, height int) *Edit {
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

	modal := &Edit{
		formData:      formData,
		originalEvent: event,
		visible:       true,
		width:         width,
		height:        height,
	}

	modal.buildForm()
	return modal
}

// buildForm creates the huh form with confirm button using LayoutStack.
func (m *Edit) buildForm() {
	modalWidth := m.width - 10
	if modalWidth > 90 {
		modalWidth = 90
	}
	if modalWidth < 50 {
		modalWidth = 50
	}

	formWidth := forms.ModalFormWidth(modalWidth)
	formHeight := forms.ModalFormHeight(m.height)

	m.form = forms.NewCaptureEventForm(m.formData, "manual", formWidth, formHeight)
}

// Init initializes the modal and its form.
//
// Returns:
//   - A tea.Cmd value.
//
// Side effects:
//   - None.
func (m *Edit) Init() tea.Cmd {
	if m.form == nil {
		return nil
	}
	return m.form.Init()
}

// Update handles messages for the edit event modal.
//
// Expected:
//   - msg must be a valid tea.Msg type.
//
// Returns:
//   - tea.Cmd: command to execute.
//   - bool: true if form completed successfully.
//   - *EditData: event data if completed, nil otherwise.
//
// Side effects:
//   - May hide modal on completion or cancellation.
//   - May rebuild form on window resize.
func (m *Edit) Update(msg tea.Msg) (tea.Cmd, bool, *EditData) {
	if !m.visible {
		return nil, false, nil
	}

	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch msg.Type {
		case tea.KeyEsc:
			m.visible = false
			return nil, false, nil

		case tea.KeyCtrlS:
			m.visible = false
			m.formData.SubmitConfirmed = true
			eventData := &EditData{
				Text:       m.formData.Text,
				Date:       m.formData.Date,
				Company:    m.formData.Company,
				Project:    m.formData.Project,
				Tags:       m.formData.Tags,
				Categories: m.formData.Categories,
			}
			return nil, true, eventData
		}
	}

	var cmd tea.Cmd
	var result widgets.ViewResult
	m.form, cmd, result = shared.HandleFormUpdate(m.form, msg,
		func(w, h int) {
			m.width = w
			m.height = h
			m.buildForm()
		},
		func() interface{} {
			m.visible = false
			return &EditData{
				Text:       m.formData.Text,
				Date:       m.formData.Date,
				Company:    m.formData.Company,
				Project:    m.formData.Project,
				Tags:       m.formData.Tags,
				Categories: m.formData.Categories,
			}
		},
	)

	if result != nil {
		switch result.Type() {
		case widgets.ResultSubmit:
			if sr, ok := result.(*widgets.SubmitViewResult); ok {
				if data, ok := sr.FormData.(*EditData); ok {
					return cmd, true, data
				}
			}
		case widgets.ResultCancel:
			return cmd, false, nil
		}
	}

	return cmd, false, nil
}

// View renders the edit event modal with proper chrome (border, background)
//
// Returns:
//   - A string value.
//
// Side effects:
//   - None.
func (m *Edit) View() string {
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
func (m *Edit) IsVisible() bool {
	return m.visible
}

// Show makes the modal visible.
//
// Side effects:
//   - None.
func (m *Edit) Show() {
	m.visible = true
}

// Hide hides the modal.
//
// Side effects:
//   - None.
func (m *Edit) Hide() {
	m.visible = false
}

// GetOriginalEvent returns the original event being edited (useful for comparison).
//
// Returns:
//   - A display.Event value.
//
// Side effects:
//   - None.
func (m *Edit) GetOriginalEvent() display.Event {
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
