package configure

import (
	"errors"
	"fmt"
	"strconv"

	"github.com/baphled/kariya/internal/cli/configtypes"
	"github.com/baphled/kariya/internal/cli/screens"
	"github.com/baphled/kariya/internal/cli/terminal"
	"github.com/baphled/kariya/internal/cli/themes"
	"github.com/baphled/kariya/internal/cli/uikit/layout"
	"github.com/baphled/kariya/internal/cli/uikit/primitives"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/huh"
)

// EditSettingsState identifies the settings editing screen in the state
// matrix. On this screen the user sees an interactive huh form listing
// every configuration setting for the selected domain. String and integer
// settings appear as text inputs, booleans as Yes/No confirm fields, and
// enumerations as select dropdowns. Tab moves between fields, Ctrl+S saves
// immediately, and Escape cancels. On completion the screen computes a diff
// of changed values and returns it as the submit result.
const EditSettingsState = "edit_settings"

// SettingsFormData holds the form values for configuration settings.
type SettingsFormData struct {
	// Values holds the current values keyed by setting key.
	// We use pointers to strings to satisfy huh form requirements.
	Values map[string]*string

	// BoolValues holds boolean values separately since huh.Confirm needs *bool.
	BoolValues map[string]*bool

	// SubmitConfirmed is set to true when form is submitted.
	SubmitConfirmed bool
}

// EditSettingsScreen allows users to edit configuration settings for a domain.
type EditSettingsScreen struct {
	domain   configtypes.ConfigurationDomain
	settings []*configtypes.ConfigurationSetting
	formData *SettingsFormData
	form     *huh.Form

	termInfo *terminal.Info
	theme    themes.Theme
	logo     layout.LogoRenderer

	originalValues map[string]string
}

// NewEditSettingsScreen creates a new edit settings screen.
//
// Expected:
//   - config must be a valid configuration object.
//   - config must be a valid configuration object.
//
// Returns:
//   - A fully initialized EditSettingsScreen ready for use.
//
// Side effects:
//   - None.
func NewEditSettingsScreen(domain configtypes.ConfigurationDomain, settings []*configtypes.ConfigurationSetting) *EditSettingsScreen {
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

	screen := &EditSettingsScreen{
		domain:         domain,
		settings:       settings,
		formData:       formData,
		originalValues: originalValues,
		termInfo:       &terminal.Info{Width: 120, Height: 40},
		theme:          themes.NewDefaultTheme(),
	}

	screen.rebuildForm()
	return screen
}

// rebuildForm creates the huh form based on current settings.
func (s *EditSettingsScreen) rebuildForm() {
	if len(s.settings) == 0 {
		s.form = nil
		return
	}

	var fields []huh.Field
	for _, setting := range s.settings {
		field := s.createFieldForSetting(setting)
		if field != nil {
			fields = append(fields, field)
		}
	}

	group := huh.NewGroup(fields...)

	huhTheme := themes.GenerateHuhTheme(s.theme)
	formHeight := s.termInfo.Height - 20
	if formHeight < 10 {
		formHeight = 10
	}

	s.form = huh.NewForm(group).
		WithTheme(huhTheme).
		WithWidth(s.termInfo.Width - 10).
		WithHeight(formHeight)
}

// createFieldForSetting creates the appropriate huh field for a setting.
func (s *EditSettingsScreen) createFieldForSetting(setting *configtypes.ConfigurationSetting) huh.Field {
	switch setting.Type {
	case "string":
		return huh.NewInput().
			Key(setting.Key).
			Title(setting.Label).
			Description(setting.Description).
			Value(s.formData.Values[setting.Key])

	case "int":
		return huh.NewInput().
			Key(setting.Key).
			Title(setting.Label).
			Description(setting.Description).
			Value(s.formData.Values[setting.Key]).
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
			Value(s.formData.BoolValues[setting.Key]).
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
			Value(s.formData.Values[setting.Key])

	default:
		return huh.NewInput().
			Key(setting.Key).
			Title(setting.Label).
			Description(setting.Description).
			Value(s.formData.Values[setting.Key])
	}
}

// Screen interface implementation

// Init initializes the screen.
//
// Returns:
//   - A tea.Cmd value.
//
// Side effects:
//   - None.
func (s *EditSettingsScreen) Init() tea.Cmd {
	if s.form != nil {
		return s.form.Init()
	}
	return nil
}

// Update handles messages.
//
// Expected:
//   - msg must be a valid tea.Msg type.
//
// Returns:
//   - tea.Cmd: command to execute.
//   - screens.ScreenResult: result indicating form state.
//
// Side effects:
//   - May rebuild form on window resize.
//   - May return CancelResult or SubmitResult.
func (s *EditSettingsScreen) Update(msg tea.Msg) (tea.Cmd, screens.ScreenResult) {
	switch msg := msg.(type) {
	case tea.WindowSizeMsg:
		s.termInfo = &terminal.Info{Width: msg.Width, Height: msg.Height}
		s.rebuildForm()
		return nil, nil

	case tea.KeyMsg:
		switch msg.Type {
		case tea.KeyEsc:
			return nil, &screens.CancelResult{}

		case tea.KeyCtrlS:
			s.formData.SubmitConfirmed = true
			return nil, &screens.SubmitResult{FormData: s.GetChanges()}
		}
	}

	if s.form != nil {
		model, cmd := s.form.Update(msg)
		if f, ok := model.(*huh.Form); ok {
			s.form = f
		}

		if s.form.State == huh.StateCompleted {
			s.formData.SubmitConfirmed = true
			return cmd, &screens.SubmitResult{FormData: s.GetChanges()}
		}

		return cmd, nil
	}

	return nil, nil
}

// View renders the screen using UIKit layout.
//
// Returns:
//   - A string value.
//
// Side effects:
//   - None.
func (s *EditSettingsScreen) View() string {
	domainLabel := formatDomainLabel(s.domain)

	var formContent string
	if s.form != nil {
		formContent = s.form.View()
	} else if len(s.settings) == 0 {
		formContent = primitives.Muted("No settings available for this domain.", s.theme).Render()
	}

	helpFooter := primitives.RenderHelpFooter(s.theme,
		primitives.NextFieldBadge(s.theme),
		primitives.SaveBadge(s.theme),
		primitives.CancelBadge(s.theme),
	)

	screenLayout := layout.NewScreenLayout(s.termInfo).
		WithTheme(s.theme).
		WithBreadcrumbs("Main Menu", "Configure System", domainLabel).
		WithContent(formContent).
		WithHelp(helpFooter).
		WithFooterSeparator(true)

	if s.logo != nil {
		screenLayout = screenLayout.WithLogo(s.logo, 2)
	}

	return screenLayout.Render()
}

// SetTerminalInfo updates terminal dimensions.
//
// Expected:
//   - int must be valid.
//
// Side effects:
//   - None.
func (s *EditSettingsScreen) SetTerminalInfo(width, height int) {
	s.termInfo = &terminal.Info{Width: width, Height: height}
	s.rebuildForm()
}

// SetTheme updates the theme.
//
// Expected:
//   - interface{} must be valid.
//
// Side effects:
//   - None.
func (s *EditSettingsScreen) SetTheme(theme interface{}) {
	if t, ok := theme.(themes.Theme); ok {
		s.theme = t
		s.rebuildForm()
	}
}

// SetLogo updates the logo.
//
// Expected:
//   - interface{} must be valid.
//   - int must be valid.
//
// Side effects:
//   - None.
func (s *EditSettingsScreen) SetLogo(logo interface{}, _ int) {
	if l, ok := logo.(layout.LogoRenderer); ok {
		s.logo = l
	}
}

// GetChanges returns the map of changed settings.
//
// Returns:
//   - A map[string]interface{} value.
//
// Side effects:
//   - None.
func (s *EditSettingsScreen) GetChanges() map[string]interface{} {
	changes := make(map[string]interface{})

	for _, setting := range s.settings {
		var newStrVal string

		if setting.Type == "bool" {
			if boolPtr, ok := s.formData.BoolValues[setting.Key]; ok && boolPtr != nil {
				newStrVal = strconv.FormatBool(*boolPtr)
			}
		} else {
			if strPtr, ok := s.formData.Values[setting.Key]; ok && strPtr != nil {
				newStrVal = *strPtr
			}
		}

		originalVal := s.originalValues[setting.Key]

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

// GetFormData returns the current form data.
//
// Returns:
//   - A fully initialized SettingsFormData ready for use.
//
// Side effects:
//   - None.
func (s *EditSettingsScreen) GetFormData() *SettingsFormData {
	return s.formData
}
