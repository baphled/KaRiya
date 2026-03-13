//nolint:revive,nolintlint // "shared" is intentional — cross-domain reusable views
package shared

import (
	tea "github.com/charmbracelet/bubbletea"

	"github.com/baphled/kariya/internal/tui/forms"
	"github.com/baphled/kariya/internal/ui/uikit/primitives"
	"github.com/baphled/kariya/internal/ui/uikit/theme"
	"github.com/baphled/kariya/internal/ui/uikit/widgets"
)

// Compile-time check that SearchView implements the widgets.View interface.
var _ widgets.View = (*SearchView)(nil)

// SearchFormData holds the form field values for searching.
type SearchFormData struct {
	SearchText string
}

// SearchView displays a search text input form.
// It implements widgets.View — the intent owns chrome (breadcrumbs, logo, footer).
type SearchView struct {
	widgets.BaseView
	form        forms.Form
	formData    *SearchFormData
	placeholder string
}

// NewSearchView creates a new SearchView with the given placeholder text.
//
// Expected:
//   - Must be a valid string.
//
// Returns:
//   - A fully initialized SearchView ready for use.
//
// Side effects:
//   - None.
func NewSearchView(placeholder string) *SearchView {
	formData := &SearchFormData{}
	v := &SearchView{
		formData:    formData,
		placeholder: placeholder,
	}
	v.rebuildForm()
	return v
}

// rebuildForm creates a new form instance with current dimensions.
func (v *SearchView) rebuildForm() {
	input := forms.NewInput(forms.FieldConfig{
		Key:         "search",
		Title:       "Search",
		Placeholder: v.placeholder,
	}).Value(&v.formData.SearchText)

	group := forms.NewGroup(input)
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
func (v *SearchView) Init() tea.Cmd {
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
func (v *SearchView) Update(msg tea.Msg) (tea.Cmd, widgets.ViewResult) {
	form, cmd, result := HandleFormUpdate(
		v.form,
		msg,
		func(width, height int) {
			v.SetTerminalInfo(width, height)
			v.rebuildForm()
		},
		func() interface{} {
			return SearchFormData{SearchText: v.formData.SearchText}
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
func (v *SearchView) RenderContent() string {
	return v.form.View()
}

// HelpText returns contextual key binding hints for the search form.
//
// Returns:
//   - A string value.
//
// Side effects:
//   - None.
func (v *SearchView) HelpText() string {
	var th theme.Theme
	if t, ok := v.GetTheme().(theme.Theme); ok && t != nil {
		th = t
	}
	return primitives.RenderHelpFooter(th,
		primitives.NextFieldBadge(th),
		primitives.SubmitBadge(th),
		primitives.CancelBadge(th),
	)
}
