package components

import (
	"github.com/charmbracelet/bubbles/textinput"
	tea "github.com/charmbracelet/bubbletea"
)

// InputField represents a single interactive input field in a form
type InputField struct {
	Label       string
	Input       textinput.Model
	Required    bool
	Validator   func(string) error
	Error       string
	Placeholder string
}

// Form is a generic form component with multiple fields
type Form struct {
	fields     []InputField
	focusIndex int
	submitted  bool
	width      int
}

// NewForm creates a new generic form
func NewForm() *Form {
	return &Form{
		fields:     make([]InputField, 0),
		focusIndex: 0,
		submitted:  false,
		width:      80,
	}
}

// AddField adds a new field to the form
func (f *Form) AddField(label string, placeholder string, required bool, validator func(string) error) *Form {
	input := textinput.New()
	input.Placeholder = placeholder
	input.Width = 60

	// Focus first field
	if len(f.fields) == 0 {
		input.Focus()
	}

	field := InputField{
		Label:       label,
		Input:       input,
		Required:    required,
		Validator:   validator,
		Placeholder: placeholder,
	}

	f.fields = append(f.fields, field)
	return f
}

// SetWidth sets the form width
func (f *Form) SetWidth(width int) {
	f.width = width
	// Update all input widths
	for i := range f.fields {
		f.fields[i].Input.Width = width - 20 // Leave margin for labels
	}
}

// Init initializes the form
func (f *Form) Init() tea.Cmd {
	return textinput.Blink
}

// Update handles messages
func (f *Form) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch msg.String() {
		case "tab", "down", "j":
			return f.nextField(), nil
		case "shift+tab", "up", "k":
			return f.prevField(), nil
		case "enter":
			if f.focusIndex == len(f.fields) {
				// On submit button
				return f.submit(), nil
			}
			return f.nextField(), nil
		}
	}

	// Update focused field
	if f.focusIndex < len(f.fields) {
		var cmd tea.Cmd
		f.fields[f.focusIndex].Input, cmd = f.fields[f.focusIndex].Input.Update(msg)
		return f, cmd
	}

	return f, nil
}

// nextField moves focus to the next field
func (f *Form) nextField() *Form {
	if f.focusIndex < len(f.fields) {
		f.fields[f.focusIndex].Input.Blur()
	}
	f.focusIndex++
	if f.focusIndex >= len(f.fields)+1 { // +1 for submit button
		f.focusIndex = 0
	}
	if f.focusIndex < len(f.fields) {
		f.fields[f.focusIndex].Input.Focus()
	}
	return f
}

// prevField moves focus to the previous field
func (f *Form) prevField() *Form {
	if f.focusIndex < len(f.fields) {
		f.fields[f.focusIndex].Input.Blur()
	}
	f.focusIndex--
	if f.focusIndex < 0 {
		f.focusIndex = len(f.fields) // Submit button
	}
	if f.focusIndex < len(f.fields) {
		f.fields[f.focusIndex].Input.Focus()
	}
	return f
}

// submit validates and submits the form
func (f *Form) submit() *Form {
	// Validate all fields
	hasErrors := false
	for i := range f.fields {
		field := &f.fields[i]
		value := field.Input.Value()

		// Check required
		if field.Required && value == "" {
			field.Error = "This field is required"
			hasErrors = true
			continue
		}

		// Run custom validator
		if field.Validator != nil {
			if err := field.Validator(value); err != nil {
				field.Error = err.Error()
				hasErrors = true
				continue
			}
		}

		// Clear error if valid
		field.Error = ""
	}

	if !hasErrors {
		f.submitted = true
	}

	return f
}

// GetValue returns the value of a field by index
func (f *Form) GetValue(index int) string {
	if index < 0 || index >= len(f.fields) {
		return ""
	}
	return f.fields[index].Input.Value()
}

// GetValues returns all field values as a map (label -> value)
func (f *Form) GetValues() map[string]string {
	values := make(map[string]string)
	for _, field := range f.fields {
		values[field.Label] = field.Input.Value()
	}
	return values
}

// IsSubmitted returns whether the form was successfully submitted
func (f *Form) IsSubmitted() bool {
	return f.submitted
}

// Reset clears all field values and errors
func (f *Form) Reset() {
	f.submitted = false
	f.focusIndex = 0
	for i := range f.fields {
		f.fields[i].Input.SetValue("")
		f.fields[i].Error = ""
	}
	if len(f.fields) > 0 {
		f.fields[0].Input.Focus()
	}
}

// View renders the form (basic implementation, can be overridden)
func (f *Form) View() string {
	// This is intentionally simple - intents can use FormFieldContainer
	// or their own rendering logic
	return "Use FormFieldContainer for rendering fields"
}

// GetFields returns all fields (for custom rendering)
func (f *Form) GetFields() []InputField {
	return f.fields
}

// GetFocusIndex returns the currently focused field index
func (f *Form) GetFocusIndex() int {
	return f.focusIndex
}

// IsSubmitButtonFocused returns true if submit button is focused
func (f *Form) IsSubmitButtonFocused() bool {
	return f.focusIndex == len(f.fields)
}
