package components

import (
	"strings"
	"testing"
)

func TestNewFormFieldContainer(t *testing.T) {
	ffc := NewFormFieldContainer()

	if ffc.isFocused {
		t.Error("expected isFocused to be false by default")
	}

	if ffc.hasLabel || ffc.hasInput || ffc.hasError || ffc.hasHint {
		t.Error("expected new container to have no sections by default")
	}
}

func TestFormFieldContainerSetLabel(t *testing.T) {
	label := "Email"
	ffc := NewFormFieldContainer().SetLabel(label)

	if ffc.label != label {
		t.Errorf("expected label %q, got %q", label, ffc.label)
	}

	if !ffc.hasLabel {
		t.Error("expected hasLabel to be true")
	}
}

func TestFormFieldContainerSetInput(t *testing.T) {
	input := "user@example.com"
	ffc := NewFormFieldContainer().SetInput(input)

	if ffc.input != input {
		t.Errorf("expected input %q, got %q", input, ffc.input)
	}

	if !ffc.hasInput {
		t.Error("expected hasInput to be true")
	}
}

func TestFormFieldContainerSetError(t *testing.T) {
	error := "Invalid email format"
	ffc := NewFormFieldContainer().SetError(error)

	if ffc.error != error {
		t.Errorf("expected error %q, got %q", error, ffc.error)
	}

	if !ffc.hasError {
		t.Error("expected hasError to be true")
	}
}

func TestFormFieldContainerSetHint(t *testing.T) {
	hint := "Enter a valid email address"
	ffc := NewFormFieldContainer().SetHint(hint)

	if ffc.hint != hint {
		t.Errorf("expected hint %q, got %q", hint, ffc.hint)
	}

	if !ffc.hasHint {
		t.Error("expected hasHint to be true")
	}
}

func TestFormFieldContainerSetFocused(t *testing.T) {
	tests := []struct {
		name     string
		focused  bool
		expected bool
	}{
		{"Focused", true, true},
		{"Not focused", false, false},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ffc := NewFormFieldContainer().SetFocused(tt.focused)

			if ffc.isFocused != tt.expected {
				t.Errorf("expected isFocused %v, got %v", tt.expected, ffc.isFocused)
			}
		})
	}
}

func TestFormFieldContainerRenderWithLabel(t *testing.T) {
	label := "Username"
	ffc := NewFormFieldContainer().SetLabel(label)

	rendered := ffc.Render()

	if !strings.Contains(rendered, label) {
		t.Errorf("expected rendered output to contain label %q", label)
	}
}

func TestFormFieldContainerRenderWithInput(t *testing.T) {
	input := "john_doe"
	ffc := NewFormFieldContainer().SetInput(input)

	rendered := ffc.Render()

	if !strings.Contains(rendered, input) {
		t.Errorf("expected rendered output to contain input %q", input)
	}
}

func TestFormFieldContainerRenderWithError(t *testing.T) {
	error := "This field is required"
	ffc := NewFormFieldContainer().SetError(error)

	rendered := ffc.Render()

	if !strings.Contains(rendered, error) {
		t.Errorf("expected rendered output to contain error %q", error)
	}
}

func TestFormFieldContainerRenderWithHint(t *testing.T) {
	hint := "Use lowercase letters and numbers"
	ffc := NewFormFieldContainer().SetHint(hint)

	rendered := ffc.Render()

	if !strings.Contains(rendered, hint) {
		t.Errorf("expected rendered output to contain hint %q", hint)
	}
}

func TestFormFieldContainerRenderAllSections(t *testing.T) {
	label := "Email"
	input := "user@example.com"
	error := "Invalid format"
	hint := "example@domain.com"

	ffc := NewFormFieldContainer().
		SetLabel(label).
		SetInput(input).
		SetError(error).
		SetHint(hint)

	rendered := ffc.Render()

	if !strings.Contains(rendered, label) {
		t.Errorf("expected rendered output to contain label %q", label)
	}

	if !strings.Contains(rendered, input) {
		t.Errorf("expected rendered output to contain input %q", input)
	}

	if !strings.Contains(rendered, error) {
		t.Errorf("expected rendered output to contain error %q", error)
	}

	if !strings.Contains(rendered, hint) {
		t.Errorf("expected rendered output to contain hint %q", hint)
	}
}

func TestFormFieldContainerBuilderChaining(t *testing.T) {
	label := "Field"
	input := "value"
	error := "Error"
	hint := "Hint"

	ffc := NewFormFieldContainer().
		SetLabel(label).
		SetInput(input).
		SetError(error).
		SetHint(hint).
		SetFocused(true)

	if ffc.label != label || ffc.input != input || ffc.error != error || ffc.hint != hint {
		t.Error("expected all sections to be set after chaining")
	}

	if !ffc.isFocused {
		t.Error("expected isFocused to be true")
	}
}

func TestFormFieldContainerRenderWithoutError(t *testing.T) {
	label := "Username"
	input := "john_doe"
	hint := "3-20 characters"

	ffc := NewFormFieldContainer().
		SetLabel(label).
		SetInput(input).
		SetHint(hint)

	rendered := ffc.Render()

	if !strings.Contains(rendered, label) || !strings.Contains(rendered, input) || !strings.Contains(rendered, hint) {
		t.Error("expected rendered output to contain all sections")
	}
}

func TestFormFieldContainerRenderFocusedState(t *testing.T) {
	label := "Password"
	input := "••••••••"

	ffc := NewFormFieldContainer().
		SetLabel(label).
		SetInput(input).
		SetFocused(true)

	rendered := ffc.Render()

	if !strings.Contains(rendered, label) || !strings.Contains(rendered, input) {
		t.Error("expected rendered output to contain label and input")
	}
}

func TestFormFieldContainerRenderErrorState(t *testing.T) {
	label := "Email"
	input := "invalid"
	error := "Invalid email"

	ffc := NewFormFieldContainer().
		SetLabel(label).
		SetInput(input).
		SetError(error).
		SetFocused(true)

	rendered := ffc.Render()

	if !strings.Contains(rendered, label) || !strings.Contains(rendered, input) || !strings.Contains(rendered, error) {
		t.Error("expected rendered output to contain all sections")
	}
}

func TestFormFieldContainerSectionOrdering(t *testing.T) {
	// Sections should appear in order: label, input, error, hint
	label := "LABEL"
	input := "INPUT"
	error := "ERROR"
	hint := "HINT"

	ffc := NewFormFieldContainer().
		SetLabel(label).
		SetInput(input).
		SetError(error).
		SetHint(hint)

	rendered := ffc.Render()

	labelPos := strings.Index(rendered, label)
	inputPos := strings.Index(rendered, input)
	errorPos := strings.Index(rendered, error)
	hintPos := strings.Index(rendered, hint)

	if labelPos == -1 || inputPos == -1 || errorPos == -1 || hintPos == -1 {
		t.Error("expected all sections to be present")
	}

	if labelPos > inputPos || inputPos > errorPos || errorPos > hintPos {
		t.Error("expected sections to appear in order: label, input, error, hint")
	}
}

func TestFormFieldContainerEmptyLabel(t *testing.T) {
	ffc := NewFormFieldContainer().SetLabel("")

	if !ffc.hasLabel {
		t.Error("expected hasLabel to be true even with empty label")
	}

	rendered := ffc.Render()

	if rendered == "" {
		t.Error("expected rendered output to be non-empty")
	}
}

func TestFormFieldContainerEmptyInput(t *testing.T) {
	label := "Field"
	ffc := NewFormFieldContainer().SetLabel(label).SetInput("")

	rendered := ffc.Render()

	if !strings.Contains(rendered, label) {
		t.Errorf("expected rendered output to contain label %q", label)
	}
}

func TestFormFieldContainerMultilineInput(t *testing.T) {
	label := "Description"
	input := "Line 1\nLine 2\nLine 3"

	ffc := NewFormFieldContainer().SetLabel(label).SetInput(input)

	rendered := ffc.Render()

	if !strings.Contains(rendered, "Line 1") || !strings.Contains(rendered, "Line 2") || !strings.Contains(rendered, "Line 3") {
		t.Error("expected rendered output to contain all input lines")
	}
}

func TestFormFieldContainerComplexError(t *testing.T) {
	label := "Password"
	input := "password"
	error := "Password must contain uppercase, lowercase, number, and special character"

	ffc := NewFormFieldContainer().
		SetLabel(label).
		SetInput(input).
		SetError(error)

	rendered := ffc.Render()

	if !strings.Contains(rendered, "uppercase") || !strings.Contains(rendered, "special character") {
		t.Error("expected rendered output to contain full error message")
	}
}

func TestFormFieldContainerWithoutLabel(t *testing.T) {
	input := "value"
	ffc := NewFormFieldContainer().SetInput(input)

	rendered := ffc.Render()

	if !strings.Contains(rendered, input) {
		t.Errorf("expected rendered output to contain input %q", input)
	}
}

func TestFormFieldContainerToggleFocused(t *testing.T) {
	label := "Field"
	input := "value"

	ffc := NewFormFieldContainer().
		SetLabel(label).
		SetInput(input).
		SetFocused(true)

	if !ffc.isFocused {
		t.Error("expected isFocused to be true")
	}

	ffc.SetFocused(false)

	if ffc.isFocused {
		t.Error("expected isFocused to be false after toggling")
	}
}

func TestFormFieldContainerRenderConsistency(t *testing.T) {
	label := "Email"
	input := "test@example.com"

	ffc := NewFormFieldContainer().SetLabel(label).SetInput(input)

	rendered1 := ffc.Render()
	rendered2 := ffc.Render()

	if rendered1 != rendered2 {
		t.Error("expected consistent rendering")
	}
}
