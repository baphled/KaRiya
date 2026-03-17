package burst

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

// Edit provides a way to edit an existing burst.
// It shows a form with Name and Description fields pre-populated from the existing burst.
//
// Usage:
//
//	modal := modals.NewEdit(existingBurst, width, height)
//	cmd := modal.Init()
//	// In Update:
//	cmd, completed, burstData := modal.Update(msg)
//	if completed && burstData != nil {
//	    // User completed form - update burst
//	} else if !modal.IsVisible() {
//	    // User cancelled (Esc)
//	}
type Edit struct {
	form          forms.Form
	formData      *forms.BurstFormData
	originalBurst display.Burst
	visible       bool
	width         int
	height        int
}

// NewEdit creates a new edit burst modal with fields pre-populated from the existing burst.
// burst: the existing burst to edit
// width, height: terminal dimensions for responsive sizing
//
// Expected:
//   - burst must be a non-nil *career.Burst pointer.
//   - width must be a positive integer.
//   - height must be a positive integer.
//
// Returns:
//   - A fully initialized Edit ready for use.
//
// Side effects:
//   - Builds form with pre-populated data.
func NewEdit(burst display.Burst, width, height int) *Edit {
	formData := &forms.BurstFormData{
		Name:        burst.Name,
		Description: burst.Description,
	}

	modal := &Edit{
		formData:      formData,
		originalBurst: burst,
		visible:       true,
		width:         width,
		height:        height,
	}

	modal.buildForm()
	return modal
}

// buildForm creates the huh form with name, description, and confirm button.
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

	m.form = forms.NewBurstForm(m.formData, formWidth, formHeight)
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

// Update handles messages for the edit burst modal.
//
// Expected:
//   - msg must be a valid tea.Msg type.
//
// Returns:
//   - tea.Cmd: command to execute.
//   - bool: true if form completed successfully.
//   - *EditData: burst data if completed, nil otherwise.
//
// Side effects:
//   - May hide modal if completed or cancelled.
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
			burstData := &EditData{
				Name:        m.formData.Name,
				Description: m.formData.Description,
			}
			return nil, true, burstData
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
				Name:        m.formData.Name,
				Description: m.formData.Description,
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

// View renders the edit burst modal with proper chrome (border, background)
//
// Returns:
//   - A string value.
//
// Side effects:
//   - None.
func (m *Edit) View() string {
	if !m.visible || m.form == nil {
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
		Title("Edit Burst").
		Content(content).
		BorderColor(theme.PrimaryColor()).
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

// Hide hides the modal.
//
// Side effects:
//   - None.
func (m *Edit) Hide() {
	m.visible = false
}

// Show shows the modal.
//
// Side effects:
//   - None.
func (m *Edit) Show() {
	m.visible = true
}

// GetOriginalBurst returns the original burst being edited.
//
// Returns:
//   - A display.Burst value.
//
// Side effects:
//   - None.
func (m *Edit) GetOriginalBurst() display.Burst {
	return m.originalBurst
}

// EditData holds the data from editing a burst.
type EditData struct {
	Name        string
	Description string
}
