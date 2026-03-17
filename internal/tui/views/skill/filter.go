package skill

import (
	"errors"
	"strconv"

	tea "github.com/charmbracelet/bubbletea"

	"github.com/baphled/kariya/internal/tui/forms"
	"github.com/baphled/kariya/internal/tui/views/shared"
	"github.com/baphled/kariya/internal/ui/uikit/primitives"
	"github.com/baphled/kariya/internal/ui/uikit/theme"
	"github.com/baphled/kariya/internal/ui/uikit/widgets"
)

// Compile-time check that Filter implements the widgets.View interface.
var _ widgets.View = (*Filter)(nil)

// FilterFieldConfig describes a multi-select filter field.
type FilterFieldConfig struct {
	Key     string
	Title   string
	Options []forms.SelectOption
}

// FilterFormData holds the submitted filter form values.
type FilterFormData struct {
	Selections  map[string][]string
	MinYearsStr string
	MaxYearsStr string
}

// ValidateYearsInput validates years input for the filter form.
// Returns nil for empty input; error if non-numeric or negative.
//
// Expected:
//   - s is a string to validate.
//
// Returns:
//   - error if validation fails, nil otherwise.
//
// Side effects:
//   - None.
func ValidateYearsInput(s string) error {
	if s == "" {
		return nil
	}
	years, err := strconv.Atoi(s)
	if err != nil {
		return errors.New("must be a number")
	}
	if years < 0 {
		return errors.New("must be positive")
	}
	return nil
}

// Filter displays a filter form with configurable multi-select fields and
// years range inputs. It implements widgets.View — the intent owns chrome
// (breadcrumbs, logo, footer).
type Filter struct {
	widgets.BaseView
	form          forms.Form
	formData      *FilterFormData
	fields        []FilterFieldConfig
	selectionPtrs map[string]*[]string
}

// NewFilter creates a new Filter view with the given field configs and defaults.
//
// Expected:
//   - fields describes the configurable multi-select fields (may be nil).
//   - defaults may be nil for fresh filter state.
//
// Returns:
//   - A fully initialized Filter ready for use.
//
// Side effects:
//   - None.
func NewFilter(fields []FilterFieldConfig, defaults *FilterFormData) *Filter {
	formData := &FilterFormData{
		Selections:  make(map[string][]string),
		MinYearsStr: "",
		MaxYearsStr: "",
	}
	if defaults != nil {
		for k, v := range defaults.Selections {
			formData.Selections[k] = append([]string{}, v...)
		}
		formData.MinYearsStr = defaults.MinYearsStr
		formData.MaxYearsStr = defaults.MaxYearsStr
	}
	v := &Filter{
		formData: formData,
		fields:   fields,
	}
	v.rebuildForm()
	return v
}

func (v *Filter) rebuildForm() {
	var groupFields []forms.Field
	selectionPtrs := make(map[string]*[]string)
	for _, field := range v.fields {
		if len(field.Options) == 0 {
			continue
		}
		if _, ok := v.formData.Selections[field.Key]; !ok {
			v.formData.Selections[field.Key] = []string{}
		}
		ptr := new([]string)
		*ptr = append([]string{}, v.formData.Selections[field.Key]...)
		selectionPtrs[field.Key] = ptr
		groupFields = append(groupFields, forms.NewMultiSelect(
			field.Key,
			field.Title,
			"",
			field.Options,
			0,
		).Value(ptr))
	}
	v.selectionPtrs = selectionPtrs
	groupFields = append(groupFields,
		forms.NewInput(forms.FieldConfig{
			Key:         "min_years",
			Title:       "Minimum Years",
			Placeholder: "0",
			Validate:    ValidateYearsInput,
		}).Value(&v.formData.MinYearsStr),
		forms.NewInput(forms.FieldConfig{
			Key:         "max_years",
			Title:       "Maximum Years",
			Placeholder: "100",
			Validate:    ValidateYearsInput,
		}).Value(&v.formData.MaxYearsStr),
	)
	v.form = forms.NewFormWithDimensions(
		v.GetTerminalWidth(),
		v.GetTerminalHeight(),
		forms.NewGroup(groupFields...),
	)
}

func (v *Filter) syncSelections() {
	for k, ptr := range v.selectionPtrs {
		v.formData.Selections[k] = append([]string{}, *ptr...)
	}
}

// Init returns the form's initial command.
//
// Returns:
//   - A tea.Cmd value.
//
// Side effects:
//   - None.
func (v *Filter) Init() tea.Cmd {
	return v.form.Init()
}

// Update handles tea.Msg and returns (cmd, ViewResult).
// Returns nil ViewResult for internal state changes.
// Returns non-nil ViewResult when the view's purpose is complete.
//
// Expected:
//   - msg is a valid tea.Msg (KeyMsg, WindowSizeMsg, or internal form message).
//
// Returns:
//   - A tea.Cmd and a widgets.ViewResult (nil while in progress, non-nil on submit/cancel).
//
// Side effects:
//   - May rebuild the internal form on WindowSizeMsg.
func (v *Filter) Update(msg tea.Msg) (tea.Cmd, widgets.ViewResult) {
	var cmd tea.Cmd
	var result widgets.ViewResult
	v.form, cmd, result = shared.HandleFormUpdate(v.form, msg,
		func(w, h int) {
			v.SetTerminalInfo(w, h)
			v.rebuildForm()
		},
		func() interface{} {
			v.syncSelections()
			return FilterFormData{
				Selections:  v.formData.Selections,
				MinYearsStr: v.formData.MinYearsStr,
				MaxYearsStr: v.formData.MaxYearsStr,
			}
		},
	)
	return cmd, result
}

// RenderContent returns the form content string — no chrome.
//
// Returns:
//   - A string value.
//
// Side effects:
//   - None.
func (v *Filter) RenderContent() string {
	return v.form.View()
}

// HelpText returns contextual key binding hints for the filter form.
//
// Returns:
//   - A string value.
//
// Side effects:
//   - None.
func (v *Filter) HelpText() string {
	var th theme.Theme
	if t, ok := v.GetTheme().(theme.Theme); ok && t != nil {
		th = t
	}
	return primitives.RenderHelpFooter(th,
		primitives.NextFieldBadge(th),
		primitives.ApplyBadge(th),
		primitives.CancelBadge(th),
	)
}
