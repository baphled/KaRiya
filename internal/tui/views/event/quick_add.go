package event

import (
	"time"

	"github.com/baphled/kariya/internal/tui/forms"
	"github.com/baphled/kariya/internal/tui/views/shared"
	"github.com/baphled/kariya/internal/ui/themes"
	"github.com/baphled/kariya/internal/ui/uikit/containers"
	"github.com/baphled/kariya/internal/ui/uikit/primitives"
	"github.com/baphled/kariya/internal/ui/uikit/widgets"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
)

// QuickAdd provides a quick way to add a new career event.
// It shows a minimal form with just Text, Date, and Company fields.
//
// Usage:
//
//	view := eventviews.NewQuickAdd(width, height)
//	cmd := modal.Init()
//	// In Update:
//	cmd, completed, eventData := modal.Update(msg)
//	if completed && eventData != nil {
//	    // User completed form - pass QuickAddData to a domain conversion function
//	} else if !modal.IsVisible() {
//	    // User cancelled (Esc)
//	}
type QuickAdd struct {
	form     forms.Form
	formData *forms.CaptureEventFormData
	visible  bool
	width    int
	height   int
}

// NewQuickAdd creates a new quick add event modal.
//
// Expected:
//   - int must be valid.
//
// Returns:
//   - A fully initialized QuickAdd ready for use.
//
// Side effects:
//   - None.
func NewQuickAdd(width, height int) *QuickAdd {
	formData := forms.NewCaptureEventFormData()
	// Default date to today.
	formData.Date = time.Now().Format("2006-01-02")

	modal := &QuickAdd{
		formData: formData,
		visible:  true,
		width:    width,
		height:   height,
	}

	modal.buildForm()
	return modal
}

// buildForm creates the huh form with confirm button using LayoutStack.
func (m *QuickAdd) buildForm() {
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
func (m *QuickAdd) Init() tea.Cmd {
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
func (m *QuickAdd) Update(msg tea.Msg) (tea.Cmd, bool, *QuickAddData) {
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
			eventData := &QuickAddData{
				Text: m.formData.Text,
				Date: m.formData.Date,
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
			return &QuickAddData{
				Text: m.formData.Text,
				Date: m.formData.Date,
			}
		},
	)

	if result != nil {
		switch result.Type() {
		case widgets.ResultSubmit:
			if sr, ok := result.(*widgets.SubmitViewResult); ok {
				if data, ok := sr.FormData.(*QuickAddData); ok {
					return cmd, true, data
				}
			}
		case widgets.ResultCancel:
			return cmd, false, nil
		}
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
func (m *QuickAdd) View() string {
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
func (m *QuickAdd) IsVisible() bool {
	return m.visible
}

// Show makes the modal visible.
//
// Side effects:
//   - None.
func (m *QuickAdd) Show() {
	m.visible = true
}

// Hide hides the modal.
//
// Side effects:
//   - None.
func (m *QuickAdd) Hide() {
	m.visible = false
}

// QuickAddData holds the minimal data needed to create a new event.
// Quick events only capture text and date - other metadata can be added via Edit.
type QuickAddData struct {
	Text string
	Date string
}
