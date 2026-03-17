package configure

import (
	"fmt"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"

	"github.com/baphled/kariya/internal/tui/forms"
	"github.com/baphled/kariya/internal/ui/configtypes"
	"github.com/baphled/kariya/internal/ui/themes"
	"github.com/baphled/kariya/internal/ui/uikit/containers"
	"github.com/baphled/kariya/internal/ui/uikit/primitives"
)

const sectionListWidth = 16

// SettingsChanges holds the changed configuration settings.
//
// The Changes field maps setting keys to their new values.
// Values can be bool, int, string, or []string depending on the setting type.
type SettingsChanges struct {
	Changes map[string]interface{}
}

// Settings is a split-panel modal overlay that displays a section
// list on the left and a settings form on the right. It replaces the entire
// multi-screen configure flow with a single overlay on the main menu.
type Settings struct {
	settings       map[configtypes.ConfigurationDomain][]*configtypes.ConfigurationSetting
	originalValues map[configtypes.ConfigurationDomain]map[string]string

	domains     []configtypes.ConfigurationDomain
	selectedIdx int

	formData   map[configtypes.ConfigurationDomain]*forms.ConfigureSettingsFormData
	activeForm forms.Form

	width  int
	height int

	theme themes.Theme

	completed bool
	cancelled bool
}

// NewSettings creates a new split-panel configure settings modal.
//
// Expected:
//   - settings must be a valid map of configuration domains to settings.
//   - width must be a positive integer.
//   - height must be a positive integer.
//
// Returns:
//   - A fully initialized Settings ready for use.
//
// Side effects:
//   - Initialises form data from settings.
func NewSettings(
	settings map[configtypes.ConfigurationDomain][]*configtypes.ConfigurationSetting,
	width, height int,
) *Settings {
	allDomains := []configtypes.ConfigurationDomain{
		configtypes.DomainSystem,
		configtypes.DomainProfile,
		configtypes.DomainExport,
		configtypes.DomainUI,
	}

	var domains []configtypes.ConfigurationDomain
	for _, d := range allDomains {
		if domainSettings, ok := settings[d]; ok && len(domainSettings) > 0 {
			domains = append(domains, d)
		}
	}

	formDataMap := make(map[configtypes.ConfigurationDomain]*forms.ConfigureSettingsFormData)
	originalValues := make(map[configtypes.ConfigurationDomain]map[string]string)

	for _, domain := range domains {
		domainSettings := settings[domain]
		formDataMap[domain] = forms.NewConfigureSettingsFormData(domainSettings)

		origMap := make(map[string]string)
		for _, s := range domainSettings {
			origMap[s.Key] = fmt.Sprintf("%v", s.Value)
		}
		originalValues[domain] = origMap
	}

	modal := &Settings{
		settings:       settings,
		originalValues: originalValues,
		domains:        domains,
		selectedIdx:    0,
		formData:       formDataMap,
		width:          width,
		height:         height,
		theme:          themes.NewDefaultTheme(),
	}

	modal.rebuildActiveForm()
	return modal
}

func (m *Settings) rebuildActiveForm() {
	if len(m.domains) == 0 {
		m.activeForm = nil
		return
	}

	domain := m.domains[m.selectedIdx]
	domainSettings := m.settings[domain]
	data := m.formData[domain]

	if len(domainSettings) == 0 || data == nil {
		m.activeForm = nil
		return
	}

	formWidth, formHeight := m.formDimensions()
	m.activeForm = forms.NewConfigureSettingsForm(
		domainSettings, data.Values, data.BoolValues, formWidth, formHeight,
	)

	if m.activeForm != nil {
		m.activeForm = forms.WithCustomSubmitKey(m.activeForm, "ctrl+s")
		m.activeForm.Init()
	}
}

func (m *Settings) formDimensions() (int, int) {
	modalWidth := m.width - 12
	if modalWidth > 100 {
		modalWidth = 100
	}
	if modalWidth < 80 {
		modalWidth = 80
	}

	modalHeight := int(float64(m.height) * 0.7)
	if modalHeight > 30 {
		modalHeight = 30
	}
	if modalHeight < 20 {
		modalHeight = 20
	}

	innerWidth := modalWidth - 4
	formWidth := innerWidth - sectionListWidth
	if formWidth < 30 {
		formWidth = 30
	}

	formHeight := forms.ModalFormHeight(modalHeight)
	return formWidth, formHeight
}

// Init initialises the modal.
//
// Returns:
//   - A tea.Cmd value.
//
// Side effects:
//   - None.
func (m *Settings) Init() tea.Cmd {
	if m.activeForm != nil {
		return m.activeForm.Init()
	}
	return nil
}

// Update handles messages for the modal.
//
// Expected:
//   - msg must be a valid tea.Msg type.
//
// Returns:
//   - A tea.Cmd value.
//
// Side effects:
//   - May update selected section or form state.
func (m *Settings) Update(msg tea.Msg) tea.Cmd {
	keyMsg, ok := msg.(tea.KeyMsg)
	if !ok {
		if m.activeForm != nil {
			var cmd tea.Cmd
			m.activeForm, cmd = forms.Update(m.activeForm, msg)
			return cmd
		}
		return nil
	}

	switch keyMsg.Type {
	case tea.KeyDown:
		return m.navigateDown()
	case tea.KeyUp:
		return m.navigateUp()
	case tea.KeyCtrlS:
		m.completed = true
		return nil
	case tea.KeyEsc:
		m.cancelled = true
		return nil
	case tea.KeyTab, tea.KeyShiftTab:
		if m.activeForm != nil {
			var cmd tea.Cmd
			m.activeForm, cmd = forms.Update(m.activeForm, msg)
			return cmd
		}
		return nil
	case tea.KeyRunes:
		if !forms.IsTextInputFocused(m.activeForm) {
			switch keyMsg.String() {
			case "j":
				return m.navigateDown()
			case "k":
				return m.navigateUp()
			}
		}
	}

	if m.activeForm != nil {
		var cmd tea.Cmd
		m.activeForm, cmd = forms.Update(m.activeForm, msg)
		return cmd
	}

	return nil
}

func (m *Settings) navigateDown() tea.Cmd {
	if m.selectedIdx < len(m.domains)-1 {
		m.selectedIdx++
		m.rebuildActiveForm()
		if m.activeForm != nil {
			return m.activeForm.Init()
		}
	}
	return nil
}

func (m *Settings) navigateUp() tea.Cmd {
	if m.selectedIdx > 0 {
		m.selectedIdx--
		m.rebuildActiveForm()
		if m.activeForm != nil {
			return m.activeForm.Init()
		}
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
func (m *Settings) View() string {
	leftPanel := m.renderSectionList()
	rightPanel := m.renderFormPanel()

	body := lipgloss.JoinHorizontal(lipgloss.Top, leftPanel, rightPanel)

	footer := primitives.RenderHelpFooter(m.theme,
		primitives.NavigateVimBadge(m.theme),
		primitives.NextFieldBadge(m.theme),
		primitives.SaveBadge(m.theme),
		primitives.CancelBadge(m.theme),
	)

	content := lipgloss.JoinVertical(lipgloss.Left, body, "", footer)

	modalWidth := m.width - 12
	if modalWidth > 100 {
		modalWidth = 100
	}
	if modalWidth < 80 {
		modalWidth = 80
	}

	box := containers.NewBox(m.theme).
		Title("⚙  Configure System").
		Content(content).
		Width(modalWidth)

	return box.Render()
}

// Render renders the modal at the specified dimensions.
//
// Expected:
//   - width must be a positive integer.
//   - height must be a positive integer.
//
// Returns:
//   - A string value.
//
// Side effects:
//   - None.
func (m *Settings) Render(width, height int) string {
	m.width = width
	m.height = height
	m.rebuildActiveForm()
	return m.View()
}

func (m *Settings) renderSectionList() string {
	var lines []string
	for i, domain := range m.domains {
		label := formatDomainLabel(domain)
		if i == m.selectedIdx {
			lines = append(lines, primitives.Title("▶ "+label, m.theme).Render())
		} else {
			lines = append(lines, primitives.Body("  "+label, m.theme).Render())
		}
	}

	listContent := lipgloss.JoinVertical(lipgloss.Left, lines...)
	panelStyle := lipgloss.NewStyle().Width(sectionListWidth)
	return panelStyle.Render(listContent)
}

func (m *Settings) renderFormPanel() string {
	if m.activeForm == nil {
		return primitives.Muted("No settings available.", m.theme).Render()
	}
	return m.activeForm.View()
}

// IsCompleted returns whether the modal was completed (saved).
//
// Returns:
//   - A bool value.
//
// Side effects:
//   - None.
func (m *Settings) IsCompleted() bool {
	return m.completed
}

// IsCancelled returns whether the modal was cancelled.
//
// Returns:
//   - A bool value.
//
// Side effects:
//   - None.
func (m *Settings) IsCancelled() bool {
	return m.cancelled
}

// GetChanges returns the changed settings across all domains.
//
// Returns:
//   - A SettingsChanges value containing all changed settings.
//
// Side effects:
//   - None.
func (m *Settings) GetChanges() *SettingsChanges {
	allChanges := make(map[string]interface{})
	for domain, data := range m.formData {
		domainSettings := m.settings[domain]
		originals := m.originalValues[domain]
		changes := forms.GetConfigureSettingChanges(domainSettings, data, originals)
		for k, v := range changes {
			allChanges[k] = v
		}
	}
	return &SettingsChanges{Changes: allChanges}
}

// SetTheme updates the theme.
//
// Expected:
//   - th must be a valid theme instance (can be nil).
//
// Side effects:
//   - None.
func (m *Settings) SetTheme(theme themes.Theme) {
	m.theme = theme
}

// SetDimensions updates the modal dimensions.
//
// Expected:
//   - int must be valid.
//
// Side effects:
//   - None.
func (m *Settings) SetDimensions(width, height int) {
	if m.width != width || m.height != height {
		m.width = width
		m.height = height
		m.rebuildActiveForm()
	}
}

func formatDomainLabel(domain configtypes.ConfigurationDomain) string {
	switch domain {
	case configtypes.DomainSystem:
		return "System"
	case configtypes.DomainProfile:
		return "Profile"
	case configtypes.DomainExport:
		return "Export"
	case configtypes.DomainUI:
		return "UI"
	default:
		return string(domain)
	}
}
