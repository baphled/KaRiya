package configure

import (
	"errors"
	"fmt"
	"strconv"

	"github.com/baphled/kariya/internal/cli/configtypes"
	"github.com/baphled/kariya/internal/cli/forms"
	"github.com/baphled/kariya/internal/cli/themes"
	"github.com/baphled/kariya/internal/cli/uikit/containers"
	"github.com/baphled/kariya/internal/cli/uikit/primitives"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/huh"
	"github.com/charmbracelet/lipgloss"
)

// EditSettingsModal is a modal overlay for editing configuration settings.
type EditSettingsModal struct {
	domain   configtypes.ConfigurationDomain
	settings []*configtypes.ConfigurationSetting
	formData *SettingsFormData
	form     *huh.Form

	// Dimensions
	width  int
	height int

	// Theme
	theme themes.Theme

	// Original values for change detection
	originalValues map[string]string

	// Modal state
	visible   bool
	completed bool
	cancelled bool
}

// NewEditSettingsModal creates a new edit settings modal.
//
// Expected:
//   - domain must be a valid ConfigurationDomain.
//   - settings must be a valid slice of ConfigurationSetting.
//   - width must be a positive integer.
//   - height must be a positive integer.
//
// Returns:
//   - A fully initialized EditSettingsModal ready for use.
//
// Side effects:
//   - Initializes form data from settings.
func NewEditSettingsModal(
	domain configtypes.ConfigurationDomain,
	settings []*configtypes.ConfigurationSetting,
	width, height int,
) *EditSettingsModal {
	// Initialize form data from settings
	formData := &SettingsFormData{
		Values:          make(map[string]*string),
		BoolValues:      make(map[string]*bool),
		SubmitConfirmed: false,
	}

	originalValues := make(map[string]string)

	for _, setting := range settings {
		strVal := fmt.Sprintf("%v", setting.Value)
		originalValues[setting.Key] = strVal

		if setting.Type == "bool" {
			boolVal := strVal == "true"
			formData.BoolValues[setting.Key] = &boolVal
		} else {
			formData.Values[setting.Key] = &strVal
		}
	}

	modal := &EditSettingsModal{
		domain:         domain,
		settings:       settings,
		formData:       formData,
		originalValues: originalValues,
		width:          width,
		height:         height,
		theme:          themes.NewDefaultTheme(),
		visible:        true,
	}

	modal.rebuildForm()
	return modal
}

// rebuildForm creates the huh form based on current settings.
func (m *EditSettingsModal) rebuildForm() {
	if len(m.settings) == 0 {
		m.form = nil
		return
	}

	var fields []huh.Field
	for _, setting := range m.settings {
		field := m.createFieldForSetting(setting)
		if field != nil {
			fields = append(fields, field)
		}
	}

	// Add submit button as the last field to enable form completion
	submitValue := true
	fields = append(fields, huh.NewConfirm().
		Key("submit").
		Title("Submit Settings").
		Affirmative("Submit").
		Negative("Cancel").
		Value(&submitValue))

	group := huh.NewGroup(fields...)

	modalWidth := m.width - 10
	if modalWidth > 90 {
		modalWidth = 90
	}
	if modalWidth < 50 {
		modalWidth = 50
	}

	formWidth := forms.ModalFormWidth(modalWidth)
	formHeight := forms.ModalFormHeight(m.height)

	group = group.WithHeight(formHeight)

	huhTheme := themes.GenerateHuhTheme(m.theme)
	m.form = huh.NewForm(group).
		WithTheme(huhTheme).
		WithWidth(formWidth).
		WithHeight(formHeight)
}

// createFieldForSetting creates the appropriate huh field for a setting.
func (m *EditSettingsModal) createFieldForSetting(setting *configtypes.ConfigurationSetting) huh.Field {
	switch setting.Type {
	case "string":
		return huh.NewInput().
			Key(setting.Key).
			Title(setting.Label).
			Description(setting.Description).
			Value(m.formData.Values[setting.Key])

	case "int":
		return huh.NewInput().
			Key(setting.Key).
			Title(setting.Label).
			Description(setting.Description).
			Value(m.formData.Values[setting.Key]).
			Validate(func(val string) error {
				if val == "" {
					return nil
				}
				_, err := strconv.Atoi(val)
				if err != nil {
					return errors.New("must be a number")
				}
				return nil
			})

	case "bool":
		return huh.NewConfirm().
			Key(setting.Key).
			Title(setting.Label).
			Description(setting.Description).
			Value(m.formData.BoolValues[setting.Key]).
			Affirmative("Yes").
			Negative("No")

	case "select":
		if len(setting.Options) == 0 {
			return nil
		}
		options := make([]huh.Option[string], len(setting.Options))
		for i, opt := range setting.Options {
			options[i] = huh.NewOption(opt, opt)
		}
		return huh.NewSelect[string]().
			Key(setting.Key).
			Title(setting.Label).
			Description(setting.Description).
			Options(options...).
			Value(m.formData.Values[setting.Key])

	default:
		return huh.NewInput().
			Key(setting.Key).
			Title(setting.Label).
			Description(setting.Description).
			Value(m.formData.Values[setting.Key])
	}
}

// Init initializes the modal.
//
// Returns:
//   - A tea.Cmd value.
//
// Side effects:
//   - None.
func (m *EditSettingsModal) Init() tea.Cmd {
	if m.form != nil {
		return m.form.Init()
	}
	return nil
}

// Update handles messages.
//
// Expected:
//   - msg must be valid.
//
// Returns:
//   - A tea.Cmd value.
//
// Side effects:
//   - None.
func (m *EditSettingsModal) Update(msg tea.Msg) tea.Cmd {
	if !m.visible {
		return nil
	}

	switch msg := msg.(type) {
	case tea.WindowSizeMsg:
		m.width = msg.Width
		m.height = msg.Height
		m.rebuildForm()
		return nil

	case tea.KeyMsg:
		switch msg.Type {
		case tea.KeyEsc:
			m.cancelled = true
			m.visible = false
			return nil

		case tea.KeyCtrlS:
			m.formData.SubmitConfirmed = true
			m.completed = true
			m.visible = false
			return nil
		}
	}

	// Delegate to form
	if m.form != nil {
		model, cmd := m.form.Update(msg)
		if f, ok := model.(*huh.Form); ok {
			m.form = f
		}

		// Check if form completed
		if m.form.State == huh.StateCompleted {
			m.formData.SubmitConfirmed = true
			m.completed = true
			m.visible = false
		}

		return cmd
	}

	return nil
}

// View renders the modal content.
//
// Returns:
//   - A string value.
//
// Side effects:
//   - None.
func (m *EditSettingsModal) View() string {
	if !m.visible {
		return ""
	}

	domainLabel := formatDomainLabel(m.domain)
	title := fmt.Sprintf("Edit %s Settings", domainLabel)

	// Build form content
	var formContent string
	if m.form != nil {
		formContent = m.form.View()
	} else if len(m.settings) == 0 {
		formContent = primitives.Muted("No settings available for this domain.", m.theme).Render()
	}

	// Build help footer
	helpFooter := primitives.RenderHelpFooter(m.theme,
		primitives.NextFieldBadge(m.theme),
		primitives.SaveBadge(m.theme),
		primitives.CancelBadge(m.theme),
	)

	// Combine content
	content := lipgloss.JoinVertical(lipgloss.Left,
		formContent,
		"",
		helpFooter,
	)

	// Render in a box
	box := containers.NewBox(m.theme).
		Title(title).
		Content(content).
		Width(m.width - 10).
		Background(m.theme.BackgroundColor())

	return box.Render()
}

// Render renders the modal at the specified dimensions.
//
// Expected:
//   - int must be valid.
//
// Returns:
//   - A string value.
//
// Side effects:
//   - None.
func (m *EditSettingsModal) Render(width, height int) string {
	m.width = width
	m.height = height
	return m.View()
}

// IsVisible returns whether the modal is visible.
//
// Returns:
//   - A bool value.
//
// Side effects:
//   - None.
func (m *EditSettingsModal) IsVisible() bool {
	return m.visible
}

// IsCompleted returns whether the form was completed (submitted).
//
// Returns:
//   - A bool value.
//
// Side effects:
//   - None.
func (m *EditSettingsModal) IsCompleted() bool {
	return m.completed
}

// IsCancelled returns whether the modal was cancelled.
//
// Returns:
//   - A bool value.
//
// Side effects:
//   - None.
func (m *EditSettingsModal) IsCancelled() bool {
	return m.cancelled
}

// GetChanges returns the map of changed settings.
//
// Returns:
//   - A map[string]interface{} value.
//
// Side effects:
//   - None.
func (m *EditSettingsModal) GetChanges() map[string]interface{} {
	changes := make(map[string]interface{})

	for _, setting := range m.settings {
		var newStrVal string

		if setting.Type == "bool" {
			if boolPtr, ok := m.formData.BoolValues[setting.Key]; ok && boolPtr != nil {
				newStrVal = strconv.FormatBool(*boolPtr)
			}
		} else {
			if strPtr, ok := m.formData.Values[setting.Key]; ok && strPtr != nil {
				newStrVal = *strPtr
			}
		}

		originalVal := m.originalValues[setting.Key]

		if newStrVal != originalVal {
			switch setting.Type {
			case "int":
				if val, err := strconv.Atoi(newStrVal); err == nil {
					changes[setting.Key] = val
				}
			case "bool":
				changes[setting.Key] = (newStrVal == "true")
			default:
				changes[setting.Key] = newStrVal
			}
		}
	}

	return changes
}

// SetTheme updates the theme.
//
// Expected:
//   - th must be a valid theme instance (can be nil).
//
// Side effects:
//   - None.
func (m *EditSettingsModal) SetTheme(theme themes.Theme) {
	m.theme = theme
	m.rebuildForm()
}

// Show makes the modal visible.
//
// Side effects:
//   - None.
func (m *EditSettingsModal) Show() {
	m.visible = true
	m.completed = false
	m.cancelled = false
}

// Hide hides the modal.
//
// Side effects:
//   - None.
func (m *EditSettingsModal) Hide() {
	m.visible = false
}

// GetFormData returns the current form data.
//
// Returns:
//   - A fully initialized SettingsFormData ready for use.
//
// Side effects:
//   - None.
func (m *EditSettingsModal) GetFormData() *SettingsFormData {
	return m.formData
}
