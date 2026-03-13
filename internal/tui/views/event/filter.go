package event

import (
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
	Selections map[string][]string
	SortBy     string
	SortOrder  string
}

// Filter displays a configurable filter form for events.
// It implements widgets.View — the intent owns chrome (breadcrumbs, logo, footer).
type Filter struct {
	widgets.BaseView
	fields     []FilterFieldConfig
	sortFields []shared.SortFieldOption
	formData   *FilterFormData
	form       forms.Form
	selectVals [][]string
}

// NewFilter creates a new Filter view with given fields, sort options, and defaults.
//
// Expected:
//   - fields: slice of FilterFieldConfig describing multi-select filter fields.
//   - sortFields: slice of SortFieldOption for sort configuration.
//   - defaults: initial FilterFormData values (may be nil).
//
// Returns:
//   - A fully initialized Filter ready for use.
//
// Side effects:
//   - None.
func NewFilter(fields []FilterFieldConfig, sortFields []shared.SortFieldOption, defaults *FilterFormData) *Filter {
	formData := &FilterFormData{
		Selections: map[string][]string{},
	}
	if defaults != nil {
		formData.Selections = make(map[string][]string)
		for k, v := range defaults.Selections {
			formData.Selections[k] = append([]string{}, v...)
		}
		formData.SortBy = defaults.SortBy
		formData.SortOrder = defaults.SortOrder
	}
	v := &Filter{
		fields:     fields,
		sortFields: sortFields,
		formData:   formData,
	}
	v.rebuildForm()
	return v
}

// rebuildForm creates a new form instance with current dimensions.
func (v *Filter) rebuildForm() {
	v.selectVals = make([][]string, len(v.fields))
	for i, field := range v.fields {
		if existing, ok := v.formData.Selections[field.Key]; ok {
			v.selectVals[i] = append([]string{}, existing...)
		}
	}

	var formFields []forms.Field
	for i, field := range v.fields {
		if len(field.Options) == 0 {
			continue
		}
		multi := forms.NewMultiSelect(field.Key, field.Title, "", field.Options, 0).
			Value(&v.selectVals[i])
		formFields = append(formFields, multi)
	}

	if len(v.sortFields) > 0 {
		sortOptions := make([]forms.SelectOption, len(v.sortFields))
		for i, sf := range v.sortFields {
			sortOptions[i] = forms.SelectOption{Key: sf.Key, Value: sf.Label}
		}
		sortBy := forms.NewSelect("sortBy", "Sort By", "", sortOptions).
			Value(&v.formData.SortBy)
		formFields = append(formFields, sortBy)

		orderOptions := []forms.SelectOption{
			{Key: "asc", Value: "Ascending"},
			{Key: "desc", Value: "Descending"},
		}
		sortOrder := forms.NewSelect("sortOrder", "Sort Order", "", orderOptions).
			Value(&v.formData.SortOrder)
		formFields = append(formFields, sortOrder)
	}

	group := forms.NewGroup(formFields...)
	v.form = forms.NewFormWithDimensions(
		v.GetTerminalWidth(),
		v.GetTerminalHeight(),
		group,
	)
}

// syncSelections copies multi-select backing values into formData.
func (v *Filter) syncSelections() {
	for i, field := range v.fields {
		if i < len(v.selectVals) {
			v.formData.Selections[field.Key] = v.selectVals[i]
		}
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
				Selections: v.formData.Selections,
				SortBy:     v.formData.SortBy,
				SortOrder:  v.formData.SortOrder,
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
