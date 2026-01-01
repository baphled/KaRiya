package components_test

import (
	"testing"

	"github.com/baphled/kariya/internal/cli/components"
)

func TestFormContainerBasicRendering(t *testing.T) {
	fc := components.NewFormContainer().
		SetWidth(80).
		AddField(components.FormField{
			Label: "Test Label",
			Input: "Test Input",
		})

	output := fc.Render()
	if output == "" {
		t.Errorf("expected non-empty output, got empty string")
	}

	if !contains(output, "Test Label") {
		t.Errorf("expected label in output")
	}
}

func TestFormContainerResponsiveLayout(t *testing.T) {
	tests := []struct {
		name           string
		width          int
		expectedLayout components.FormLayout
	}{
		{
			name:           "very narrow terminal uses single column",
			width:          40,
			expectedLayout: components.SingleColumn,
		},
		{
			name:           "wide terminal uses two column",
			width:          120,
			expectedLayout: components.TwoColumn,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			fc := components.NewFormContainer().
				SetWidth(tt.width).
				SetLayout(components.Responsive).
				AddField(components.FormField{
					Label: "Field 1",
					Input: "Input 1",
				}).
				AddField(components.FormField{
					Label: "Field 2",
					Input: "Input 2",
				})

			layout := fc.GetOptimalLayout()
			if layout != tt.expectedLayout {
				t.Errorf("expected layout %v, got %v", tt.expectedLayout, layout)
			}
		})
	}
}

func TestFormContainerSingleColumnLayout(t *testing.T) {
	fc := components.NewFormContainer().
		SetWidth(80).
		SetLayout(components.SingleColumn).
		AddField(components.FormField{
			Label: "Field 1",
			Input: "Input 1",
		}).
		AddField(components.FormField{
			Label: "Field 2",
			Input: "Input 2",
		})

	output := fc.Render()

	if !contains(output, "Field 1") {
		t.Errorf("expected 'Field 1' in output")
	}
	if !contains(output, "Field 2") {
		t.Errorf("expected 'Field 2' in output")
	}

	layout := fc.GetOptimalLayout()
	if layout != components.SingleColumn {
		t.Errorf("expected SingleColumn layout, got %v", layout)
	}
}

func TestFormContainerTwoColumnLayout(t *testing.T) {
	fc := components.NewFormContainer().
		SetWidth(120).
		SetLayout(components.TwoColumn).
		AddField(components.FormField{
			Label:     "Field 1",
			Input:     "Input 1",
			FullWidth: false,
		}).
		AddField(components.FormField{
			Label:     "Field 2",
			Input:     "Input 2",
			FullWidth: false,
		})

	output := fc.Render()

	if !contains(output, "Field 1") {
		t.Errorf("expected 'Field 1' in output")
	}
	if !contains(output, "Field 2") {
		t.Errorf("expected 'Field 2' in output")
	}

	layout := fc.GetOptimalLayout()
	if layout != components.TwoColumn {
		t.Errorf("expected TwoColumn layout, got %v", layout)
	}
}

func TestFormContainerFullWidthFields(t *testing.T) {
	fc := components.NewFormContainer().
		SetWidth(120).
		SetLayout(components.TwoColumn).
		AddField(components.FormField{
			Label:     "Regular Field 1",
			Input:     "Input 1",
			FullWidth: false,
		}).
		AddField(components.FormField{
			Label:     "Full Width Field",
			Input:     "Full Width Input",
			FullWidth: true,
		})

	output := fc.Render()

	if !contains(output, "Regular Field 1") {
		t.Errorf("expected 'Regular Field 1' in output")
	}
	if !contains(output, "Full Width Field") {
		t.Errorf("expected 'Full Width Field' in output")
	}
}

func TestFormContainerIsWideLayout(t *testing.T) {
	tests := []struct {
		name         string
		width        int
		expectedWide bool
	}{
		{
			name:         "narrow width is not wide",
			width:        40,
			expectedWide: false,
		},
		{
			name:         "wide width is wide",
			width:        120,
			expectedWide: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			fc := components.NewFormContainer().
				SetWidth(tt.width).
				SetLayout(components.Responsive).
				AddField(components.FormField{
					Label: "Test",
					Input: "Input",
				})

			isWide := fc.IsWideLayout()
			if isWide != tt.expectedWide {
				t.Errorf("expected IsWideLayout() = %v, got %v", tt.expectedWide, isWide)
			}
		})
	}
}

func TestFormContainerEmptyFields(t *testing.T) {
	fc := components.NewFormContainer().
		SetWidth(80)

	output := fc.Render()
	if output != "" {
		t.Errorf("expected empty output for empty form, got: %s", output)
	}
}

// Helper function
func contains(s, substr string) bool {
	for i := 0; i <= len(s)-len(substr); i++ {
		if s[i:i+len(substr)] == substr {
			return true
		}
	}
	return false
}
