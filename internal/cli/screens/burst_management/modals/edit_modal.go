package modals

import (
	"github.com/baphled/kariya/internal/cli/forms"
	"github.com/baphled/kariya/internal/cli/themes"
	"github.com/baphled/kariya/internal/cli/uikit/containers"
	"github.com/baphled/kariya/internal/cli/uikit/primitives"
	"github.com/baphled/kariya/internal/domain/career"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
)

// EditBurstModal provides a way to edit an existing burst.
// It shows a form with Name and Description fields pre-populated from the existing burst.
//
// Usage:
//
//	modal := modals.NewEditBurstModal(existingBurst, width, height)
//	cmd := modal.Init()
//	// In Update:
//	cmd, completed, burstData := modal.Update(msg)
//	if completed && burstData != nil {
//	    // User completed form - update burst
//	} else if !modal.IsVisible() {
//	    // User cancelled (Esc)
//	}
type EditBurstModal struct {
	form          forms.Form
	formData      *forms.BurstFormData
	originalBurst *career.Burst
	visible       bool
	width         int
	height        int
}

// NewEditBurstModal creates a new edit burst modal with fields pre-populated from the existing burst.
// burst: the existing burst to edit
// width, height: terminal dimensions for responsive sizing
//
// Expected:
//   - burst must be a non-nil *career.Burst pointer.
//   - width must be a positive integer.
//   - height must be a positive integer.
//
// Returns:
//   - A fully initialized EditBurstModal ready for use.
//
// Side effects:
//   - Builds form with pre-populated data.
func NewEditBurstModal(burst *career.Burst, width, height int) *EditBurstModal {
	// Pre-populate form data from existing burst using forms package helper.
	formData := forms.GetBurstFormData(burst)

	modal := &EditBurstModal{
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
func (m *EditBurstModal) buildForm() {
	// Calculate form width.
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
func (m *EditBurstModal) Init() tea.Cmd {
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
//   - *EditBurstData: burst data if completed, nil otherwise.
//
// Side effects:
//   - May hide modal if completed or cancelled.
//   - May rebuild form on window resize.
func (m *EditBurstModal) Update(msg tea.Msg) (tea.Cmd, bool, *EditBurstData) {
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
			burstData := &EditBurstData{
				Name:        m.formData.Name,
				Description: m.formData.Description,
			}
			return nil, true, burstData
		}
	}

	// Update form using forms package helper.
	var cmd tea.Cmd
	m.form, cmd = forms.Update(m.form, msg)

	// Check if form is complete AND user confirmed submission.
	if forms.IsCompleted(m.form) {
		m.visible = false
		// Only return data if user confirmed (pressed Submit, not Cancel).
		if m.formData.SubmitConfirmed {
			burstData := &EditBurstData{
				Name:        m.formData.Name,
				Description: m.formData.Description,
			}
			return cmd, true, burstData
		}
		// User cancelled - close without returning data.
		return cmd, false, nil
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
func (m *EditBurstModal) View() string {
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
func (m *EditBurstModal) IsVisible() bool {
	return m.visible
}

// Hide hides the modal.
//
// Side effects:
//   - None.
func (m *EditBurstModal) Hide() {
	m.visible = false
}

// Show shows the modal.
//
// Side effects:
//   - None.
func (m *EditBurstModal) Show() {
	m.visible = true
}

// GetOriginalBurst returns the original burst being edited.
//
// Returns:
//   - A fully initialized career.Burst ready for use.
//
// Side effects:
//   - None.
func (m *EditBurstModal) GetOriginalBurst() *career.Burst {
	return m.originalBurst
}

// EditBurstData holds the data from editing a burst.
type EditBurstData struct {
	Name        string
	Description string
}
