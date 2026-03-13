//nolint:revive,nolintlint // "shared" is intentional — cross-domain reusable views
package shared

import (
	tea "github.com/charmbracelet/bubbletea"

	"github.com/baphled/kariya/internal/tui/forms"
	"github.com/baphled/kariya/internal/ui/uikit/primitives"
	"github.com/baphled/kariya/internal/ui/uikit/theme"
	"github.com/baphled/kariya/internal/ui/uikit/widgets"
)

// Compile-time check that SortView implements the widgets.View interface.
var _ widgets.View = (*SortView)(nil)

// SortFieldOption represents a selectable sort field.
type SortFieldOption struct {
	Key   string
	Label string
}

// SortFormData holds the form field values for sorting.
type SortFormData struct {
	SortBy    string
	SortOrder string
}

// SortView displays a sort configuration form with sort field and sort order selects.
// It implements widgets.View — the intent owns chrome (breadcrumbs, logo, footer).
type SortView struct {
	widgets.BaseView
	form     forms.Form
	formData *SortFormData
	fields   []SortFieldOption
}

// NewSortView creates a new SortView with the given field options and defaults.
//
// Expected:
//   - fields must contain at least one SortFieldOption.
//   - defaultSortBy must match a Key in fields.
//   - defaultSortOrder must be "asc" or "desc".
//
// Returns:
//   - A fully initialized SortView ready for use.
//
// Side effects:
//   - None.
func NewSortView(fields []SortFieldOption, defaultSortBy, defaultSortOrder string) *SortView {
	formData := &SortFormData{
		SortBy:    defaultSortBy,
		SortOrder: defaultSortOrder,
	}
	v := &SortView{
		formData: formData,
		fields:   fields,
	}
	v.rebuildForm()
	return v
}

func (v *SortView) rebuildForm() {
	fieldOptions := make([]forms.SelectOption, len(v.fields))
	for i, f := range v.fields {
		fieldOptions[i] = forms.SelectOption{Key: f.Key, Value: f.Label}
	}

	orderOptions := []forms.SelectOption{
		{Key: "asc", Value: "Ascending"},
		{Key: "desc", Value: "Descending"},
	}

	sortBySelect := forms.NewSelect("sortBy", "Sort By", "", fieldOptions).
		Value(&v.formData.SortBy)
	sortOrderSelect := forms.NewSelect("sortOrder", "Sort Order", "", orderOptions).
		Value(&v.formData.SortOrder)

	group := forms.NewGroup(sortBySelect, sortOrderSelect)
	v.form = forms.NewFormWithDimensions(
		v.GetTerminalWidth(),
		v.GetTerminalHeight(),
		group,
	)
}

// Init returns the form's initial command.
//
// Returns:
//   - A tea.Cmd value.
//
// Side effects:
//   - None.
func (v *SortView) Init() tea.Cmd {
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
func (v *SortView) Update(msg tea.Msg) (tea.Cmd, widgets.ViewResult) {
	form, cmd, result := HandleFormUpdate(
		v.form,
		msg,
		func(width, height int) {
			v.SetTerminalInfo(width, height)
			v.rebuildForm()
		},
		func() interface{} {
			return SortFormData{SortBy: v.formData.SortBy, SortOrder: v.formData.SortOrder}
		},
	)
	v.form = form
	return cmd, result
}

// RenderContent returns the form content string — no chrome.
//
// Returns:
//   - A string value.
//
// Side effects:
//   - None.
func (v *SortView) RenderContent() string {
	return v.form.View()
}

// HelpText returns contextual key binding hints for the sort form.
//
// Returns:
//   - A string value.
//
// Side effects:
//   - None.
func (v *SortView) HelpText() string {
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
