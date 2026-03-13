package event

import (
	"github.com/baphled/kariya/internal/tui/forms"
	"github.com/baphled/kariya/internal/ui/display"
	"github.com/baphled/kariya/internal/ui/types"
	"github.com/baphled/kariya/internal/ui/uikit/widgets"
	tea "github.com/charmbracelet/bubbletea"
)

// Form provides a form for capturing career events.
type Form struct {
	widgets.BaseView
	strategy    types.CaptureStrategy
	formData    *forms.CaptureEventFormData
	formBuilder func(data *forms.CaptureEventFormData, width, height int) forms.Form
	form        forms.Form
	footer      string
}

// NewForm creates a new event form view with the specified strategy.
//
// Expected:
//   - event must be valid.
//   - capturestrategy must be valid.
//
// Returns:
//   - A fully initialized Form ready for use.
//
// Side effects:
//   - None.
func NewForm(evt display.Event, strategy types.CaptureStrategy) *Form {
	formData := captureEventFormDataFromDisplayEvent(evt)

	strategyStr := string(strategy)
	builder := func(data *forms.CaptureEventFormData, w, h int) forms.Form {
		return forms.NewCaptureEventForm(data, strategyStr, w, h)
	}

	view := &Form{
		strategy:    strategy,
		formData:    formData,
		formBuilder: builder,
		footer:      "Esc: Back  Tab: Next field  Enter: Select",
	}
	view.SetTerminalInfo(120, 40)
	view.rebuildForm()
	return view
}
func captureEventFormDataFromDisplayEvent(evt display.Event) *forms.CaptureEventFormData {
	if isEmptyDisplayEvent(evt) {
		return forms.NewCaptureEventFormData()
	}

	date := ""
	if !evt.Date.IsZero() {
		date = evt.Date.Format("2006-01-02")
	}

	return &forms.CaptureEventFormData{
		Text:       evt.Text,
		Date:       date,
		Company:    evt.Company,
		Project:    evt.Project,
		Tags:       append([]string(nil), evt.Tags...),
		Categories: append([]string(nil), evt.Categories...),
	}
}

func isEmptyDisplayEvent(evt display.Event) bool {
	return evt.ID == "" &&
		evt.Text == "" &&
		evt.Date.IsZero() &&
		evt.Company == "" &&
		evt.Project == "" &&
		len(evt.Tags) == 0 &&
		len(evt.Categories) == 0 &&
		len(evt.Skills) == 0 &&
		evt.CreatedAt.IsZero() &&
		evt.UpdatedAt.IsZero()
}

// Init returns nil (no async loading needed).
//
// Returns:
//   - A tea.Cmd value.
//
// Side effects:
//   - None.
func (v *Form) Init() tea.Cmd {
	return nil
}

// Update handles messages and returns result when form is complete or cancelled.
//
// Expected:
//   - msg is a valid tea.Msg (key press, window resize, or custom message).
//
// Returns:
//   - A tea.Cmd and widgets.ViewResult pair.
//
// Side effects:
//   - Updates internal view state based on message type.
func (v *Form) Update(msg tea.Msg) (tea.Cmd, widgets.ViewResult) {
	switch msg := msg.(type) {
	case tea.WindowSizeMsg:
		v.SetTerminalInfo(msg.Width, msg.Height)
		v.rebuildForm()
		return nil, nil

	case tea.KeyMsg:
		if msg.Type == tea.KeyEsc || msg.String() == "esc" {
			return nil, &widgets.CancelViewResult{}
		}

		if msg.Type == tea.KeyCtrlS {
			if v.formData != nil {
				v.formData.ConfirmSubmit()
				return nil, &widgets.SubmitViewResult{FormData: v.formData}
			}
			return nil, nil
		}

		var cmd tea.Cmd
		v.form, cmd = forms.Update(v.form, msg)
		if forms.IsCompleted(v.form) {
			return cmd, &widgets.SubmitViewResult{FormData: v.formData}
		}
		return cmd, nil
	}

	var cmd tea.Cmd
	v.form, cmd = forms.Update(v.form, msg)
	if forms.IsCompleted(v.form) {
		return cmd, &widgets.SubmitViewResult{FormData: v.formData}
	}
	return cmd, nil
}

// RenderContent returns the form view content string — no chrome.
//
// Returns:
//   - A string value.
//
// Side effects:
//   - None.
func (v *Form) RenderContent() string {
	if v.form == nil {
		return ""
	}
	return v.form.View()
}

// HelpText returns the footer help text for this view.
//
// Returns:
//   - A string value.
//
// Side effects:
//   - None.
func (v *Form) HelpText() string {
	return v.footer
}

// GetStrategy returns the current capture strategy.
//
// Returns:
//   - A types.CaptureStrategy value.
//
// Side effects:
//   - None.
func (v *Form) GetStrategy() types.CaptureStrategy {
	return v.strategy
}

// GetFormData returns the form data structure.
//
// Returns:
//   - A fully initialized forms.CaptureEventFormData ready for use.
//
// Side effects:
//   - None.
func (v *Form) GetFormData() *forms.CaptureEventFormData {
	return v.formData
}

// SetFooter updates the footer help text.
//
// Expected:
//   - Must be a valid string.
//
// Side effects:
//   - None.
func (v *Form) SetFooter(footer string) {
	v.footer = footer
}

func (v *Form) rebuildForm() {
	if v.formBuilder == nil {
		return
	}

	formWidth := v.GetTerminalWidth() - 4
	if formWidth > 80 {
		formWidth = 80
	}
	if formWidth < 20 {
		formWidth = 20
	}

	formHeight := forms.DefaultFormHeight(v.GetTerminalHeight())
	v.form = v.formBuilder(v.formData, formWidth, formHeight)
	if v.form != nil {
		v.form.Init()
	}
}
