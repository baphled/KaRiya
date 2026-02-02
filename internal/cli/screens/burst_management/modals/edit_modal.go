package modals

import (
	"github.com/baphled/kariya/internal/cli/forms"
	"github.com/baphled/kariya/internal/cli/themes"
	"github.com/baphled/kariya/internal/cli/uikit/containers"
	"github.com/baphled/kariya/internal/domain/career"
	tea "github.com/charmbracelet/bubbletea"
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

	// Use natural height - overlay will handle positioning.
	formHeight := 0
	if m.height > 0 {
		formHeight = m.height - 10 // Reserve space for modal chrome
		if formHeight < 15 {
			formHeight = 15
		}
	}

	// Use the forms package burst editor form with dimensions.
	m.form = forms.NewBurstEditorFormWithDataAndDimensions(m.formData, modalWidth, formHeight)
}

// Init initializes the modal and its form.
//
// Returns:
//   - tea.Cmd: command to execute, or nil if form is nil.
//
// Side effects:
//   - Initializes the underlying form.
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
		switch msg.String() {
		case "esc":
			// Close modal without saving.
			m.visible = false
			return nil, false, nil
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
// for overlay compositing.
//
// Returns:
//   - string: the rendered modal view, or empty string if not visible.
//
// Side effects:
//   - None.
func (m *EditBurstModal) View() string {
	if !m.visible || m.form == nil {
		return ""
	}

	// Get theme or use default.
	theme := themes.NewDefaultTheme()

	formView := m.form.View()

	// Wrap in Box with solid background for overlay.
	return containers.NewBox(theme).
		Title("Edit Burst").
		Content(formView).
		BorderColor(theme.PrimaryColor()).
		Background(theme.BackgroundColor()). // REQUIRED for modal overlays
		Render()
}

// IsVisible returns whether the modal is currently visible.
//
// Returns:
//   - bool: true if modal is visible.
//
// Side effects:
//   - None.
func (m *EditBurstModal) IsVisible() bool {
	return m.visible
}

// Hide hides the modal.
//
// Side effects:
//   - Sets visible flag to false.
func (m *EditBurstModal) Hide() {
	m.visible = false
}

// Show shows the modal.
//
// Side effects:
//   - Sets visible flag to true.
func (m *EditBurstModal) Show() {
	m.visible = true
}

// GetOriginalBurst returns the original burst being edited.
//
// Returns:
//   - *career.Burst: the original burst.
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
