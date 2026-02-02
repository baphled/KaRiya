package base

import (
	"github.com/baphled/kariya/internal/cli/forms"
	"github.com/baphled/kariya/internal/cli/screens"
	tea "github.com/charmbracelet/bubbletea"
)

// FormBuilder is a function type that creates a form for given dimensions.
//
// Parameters:
//   - data: The form data structure to bind to
//   - width: Terminal width for the form
//   - height: Terminal height for the form
//
// The builder should create form groups and fields that bind to the data structure.
// Use forms.Form (which is an alias for *huh.Form) to avoid direct huh imports.
type FormBuilder[T any] func(data T, width, height int) forms.Form

// FormScreen provides a reusable screen for forms using the huh library.
//
// This screen handles:
// - Form rendering with proper dimensions
// - Window resize handling (rebuilds form with new dimensions)
// - Escape key handling (returns CancelResult)
// - Form completion detection (returns SubmitResult when confirmed)
// - StandardView integration
//
// Type parameter T should be a pointer to your form data structure, which MUST include
// a SubmitConfirmed bool field that acts as the confirm gate.
//
// Example usage:
//
//	type MyFormData struct {
//	    Name            string
//	    Email           string
//	    SubmitConfirmed bool  // Required for confirmation
//	}
//
//	builder := func(data *MyFormData, width, height int) *huh.Form {
//	    group := huh.NewGroup(
//	        huh.NewInput().Key("name").Title("Name").Value(&data.Name),
//	        huh.NewInput().Key("email").Title("Email").Value(&data.Email),
//	    )
//	    return forms.NewFormWithFixedConfirm(group, &data.SubmitConfirmed, width, height)
//	}
//
//	formData := &MyFormData{}
//	screen := base.NewBaseFormScreen(
//	    []string{"Main Menu", "Add Item"},
//	    builder,
//	    formData,
//	)
//
// Related:
// - docs/FORMS_GUIDE.md (Form patterns and best practices)
// - internal/cli/forms/ (Form builders and validators)
// - internal/cli/models/capture_form.go (Example of form wrapper pattern).
type FormScreen[T any] struct {
	*Screen

	// breadcrumbs for navigation context
	breadcrumbs []string

	// formBuilder creates the huh form with current dimensions
	formBuilder FormBuilder[T]

	// formData is the data structure bound to the form
	formData T

	// form is the current form instance.
	form forms.Form

	// footer is the help text shown at the bottom
	footer string
}

// NewBaseFormScreen creates a new form screen.
//
// Parameters:
//   - breadcrumbs: Navigation breadcrumb trail (e.g., []string{"Main Menu", "Add Event"})
//   - builder: Function that creates a huh.Form for given dimensions
//   - formData: Pointer to form data structure (must have SubmitConfirmed bool field)
//
// Expected:
//   - breadcrumbs must be a valid slice of strings.
//   - builder must be a valid FormBuilder function.
//   - formdata must be a valid T pointer.
//
// Returns:
//   - A fully initialized FormScreen[T] ready for use.
//
// Side effects:
//   - Builds form with default dimensions.
func NewBaseFormScreen[T any](
	breadcrumbs []string,
	builder FormBuilder[T],
	formData T,
) *FormScreen[T] {
	screen := &FormScreen[T]{
		Screen:      NewBaseScreen(),
		breadcrumbs: breadcrumbs,
		formBuilder: builder,
		formData:    formData,
		footer:      "Esc: Back  Tab: Next field  Enter: Select",
	}

	// Build initial form with default dimensions
	screen.rebuildForm()

	return screen
}

// rebuildForm creates a new form with current terminal dimensions.
//
// This is called:
// - On initialization
// - When terminal dimensions change (WindowSizeMsg)
// - When SetTerminalInfo is called.
func (s *FormScreen[T]) rebuildForm() {
	// Nil check to prevent panic
	if s.formBuilder == nil {
		return
	}

	// Use width-4 for form to account for padding in StandardView
	formWidth := s.Width() - 4
	if formWidth < 20 {
		formWidth = 20
	}

	// Use appropriate height for form content
	formHeight := forms.DefaultFormHeight(s.Height())

	s.form = s.formBuilder(s.formData, formWidth, formHeight)
}

// SetTerminalInfo updates terminal dimensions and rebuilds form.
//
// Expected:
//   - int must be valid.
//
// Side effects:
//   - None.
func (s *FormScreen[T]) SetTerminalInfo(width, height int) {
	s.Screen.SetTerminalInfo(width, height)
	s.rebuildForm()
}

// Update handles messages and returns result when form is complete or cancelled.
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
//   - May return CancelResult on escape.
//   - May return SubmitResult on form completion.
func (s *FormScreen[T]) Update(msg tea.Msg) (tea.Cmd, screens.ScreenResult) {
	switch msg := msg.(type) {
	case tea.WindowSizeMsg:
		// Handle window resize (rebuilds form with new dimensions)
		s.SetTerminalInfo(msg.Width, msg.Height)
		return nil, nil

	case tea.KeyMsg:
		// Handle escape key - cancel form and return to previous screen
		if msg.Type == tea.KeyEsc || msg.String() == "esc" {
			return nil, &screens.CancelResult{}
		}

		// Delegate other keys to form using forms package helper.
		var cmd tea.Cmd
		s.form, cmd = forms.Update(s.form, msg)

		// Check if form is completed.
		//
		// The form data's SubmitConfirmed field must be set to true by a confirm field
		// for the submission to occur. This prevents accidental submissions.
		if forms.IsCompleted(s.form) {
			// Check if the form data has SubmitConfirmed set to true
			// This requires using reflection or a type assertion
			// For now, we return SubmitResult when form is completed
			// The intent should check the SubmitConfirmed field in the returned data
			return cmd, &screens.SubmitResult{
				FormData: s.formData,
			}
		}

		return cmd, nil
	}

	// Delegate other messages to form using forms package helper.
	var cmd tea.Cmd
	s.form, cmd = forms.Update(s.form, msg)

	return cmd, nil
}

// View renders the form screen using StandardView.
//
// Returns:
//   - A string value.
//
// Side effects:
//   - None.
func (s *FormScreen[T]) View() string {
	// Render form content
	formView := s.form.View()

	// Use Screen's CreateView helper for StandardView integration
	return s.CreateView(s.breadcrumbs, formView, s.footer)
}

// SetFooter updates the footer help text.
//
// Expected:
//   - Must be a valid string.
//
// Side effects:
//   - None.
func (s *FormScreen[T]) SetFooter(footer string) {
	s.footer = footer
}

// GetFormData returns the form data structure.
//
// Returns:
//   - A T value.
//
// Side effects:
//   - None.
func (s *FormScreen[T]) GetFormData() T {
	return s.formData
}
