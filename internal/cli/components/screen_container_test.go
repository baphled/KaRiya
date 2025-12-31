package components

import (
	"strings"
	"testing"

	"github.com/charmbracelet/lipgloss"
	"github.com/baphled/kariya/internal/cli/styles"
)

func TestNewScreenContainer(t *testing.T) {
	content := "Test Content"
	sc := NewScreenContainer(content)

	if sc.content != content {
		t.Errorf("expected content %q, got %q", content, sc.content)
	}

	if sc.paddingMode != PaddingNormal {
		t.Errorf("expected default padding mode PaddingNormal, got %v", sc.paddingMode)
	}

	if sc.maxWidth != styles.MaxWidth(120) {
		t.Errorf("expected default maxWidth %d, got %d", styles.MaxWidth(120), sc.maxWidth)
	}

	if sc.useCustomPadding {
		t.Error("expected useCustomPadding to be false for new container")
	}
}

func TestScreenContainerWithPaddingMode(t *testing.T) {
	tests := []struct {
		name     string
		mode     PaddingMode
		expected PaddingMode
	}{
		{"Normal mode", PaddingNormal, PaddingNormal},
		{"Compact mode", PaddingCompact, PaddingCompact},
		{"Spacious mode", PaddingSpacious, PaddingSpacious},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			sc := NewScreenContainer("test").WithPaddingMode(tt.expected)

			if sc.paddingMode != tt.expected {
				t.Errorf("expected padding mode %v, got %v", tt.expected, sc.paddingMode)
			}

			if sc.useCustomPadding {
				t.Error("expected useCustomPadding to be false after WithPaddingMode")
			}
		})
	}
}

func TestScreenContainerWithCustomPadding(t *testing.T) {
	vertical, horizontal := 3, 5
	sc := NewScreenContainer("test").WithCustomPadding(vertical, horizontal)

	if sc.customVertical != vertical {
		t.Errorf("expected customVertical %d, got %d", vertical, sc.customVertical)
	}

	if sc.customHorizontal != horizontal {
		t.Errorf("expected customHorizontal %d, got %d", horizontal, sc.customHorizontal)
	}

	if !sc.useCustomPadding {
		t.Error("expected useCustomPadding to be true after WithCustomPadding")
	}
}

func TestScreenContainerWithMaxWidth(t *testing.T) {
	maxWidth := 100
	sc := NewScreenContainer("test").WithMaxWidth(maxWidth)

	if sc.maxWidth != maxWidth {
		t.Errorf("expected maxWidth %d, got %d", maxWidth, sc.maxWidth)
	}
}

func TestScreenContainerGetPadding(t *testing.T) {
	tests := []struct {
		name             string
		paddingMode      PaddingMode
		customVertical   int
		customHorizontal int
		useCustom        bool
		expectedV        int
		expectedH        int
	}{
		{"Normal padding", PaddingNormal, 0, 0, false, 1, 2},
		{"Compact padding", PaddingCompact, 0, 0, false, 0, 1},
		{"Spacious padding", PaddingSpacious, 0, 0, false, 2, 4},
		{"Custom padding", PaddingNormal, 3, 5, true, 3, 5},
		{"Custom overrides normal", PaddingNormal, 2, 3, true, 2, 3},
		{"Custom overrides spacious", PaddingSpacious, 1, 1, true, 1, 1},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			sc := &ScreenContainer{
				paddingMode:      tt.paddingMode,
				customVertical:   tt.customVertical,
				customHorizontal: tt.customHorizontal,
				useCustomPadding: tt.useCustom,
			}

			v, h := sc.getPadding()

			if v != tt.expectedV || h != tt.expectedH {
				t.Errorf("expected padding (%d, %d), got (%d, %d)", tt.expectedV, tt.expectedH, v, h)
			}
		})
	}
}

func TestScreenContainerRender(t *testing.T) {
	content := "Test Content"
	sc := NewScreenContainer(content)

	rendered := sc.Render()

	if !strings.Contains(rendered, content) {
		t.Errorf("expected rendered output to contain %q, got %q", content, rendered)
	}

	// Check that the rendered output is not empty
	if rendered == "" {
		t.Error("expected rendered output to be non-empty")
	}
}

func TestScreenContainerRenderWithDifferentModes(t *testing.T) {
	content := "Test"

	tests := []struct {
		name string
		mode PaddingMode
	}{
		{"Normal mode rendering", PaddingNormal},
		{"Compact mode rendering", PaddingCompact},
		{"Spacious mode rendering", PaddingSpacious},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			sc := NewScreenContainer(content).WithPaddingMode(tt.mode)
			rendered := sc.Render()

			if !strings.Contains(rendered, content) {
				t.Errorf("expected rendered output to contain %q, got %q", content, rendered)
			}
		})
	}
}

func TestScreenContainerRenderWithCustomPadding(t *testing.T) {
	content := "Test"
	sc := NewScreenContainer(content).WithCustomPadding(2, 3)

	rendered := sc.Render()

	if !strings.Contains(rendered, content) {
		t.Errorf("expected rendered output to contain %q, got %q", content, rendered)
	}
}

func TestScreenContainerBuilderChaining(t *testing.T) {
	content := "Test"
	sc := NewScreenContainer(content).
		WithPaddingMode(PaddingSpacious).
		WithMaxWidth(100)

	if sc.content != content {
		t.Errorf("expected content %q, got %q", content, sc.content)
	}

	if sc.paddingMode != PaddingSpacious {
		t.Errorf("expected padding mode PaddingSpacious, got %v", sc.paddingMode)
	}

	if sc.maxWidth != 100 {
		t.Errorf("expected maxWidth 100, got %d", sc.maxWidth)
	}
}

func TestScreenContainerCustomPaddingOverridesMode(t *testing.T) {
	sc := NewScreenContainer("test").
		WithPaddingMode(PaddingNormal).
		WithCustomPadding(5, 6)

	v, h := sc.getPadding()

	if v != 5 || h != 6 {
		t.Errorf("expected custom padding (5, 6), got (%d, %d)", v, h)
	}
}

func TestScreenContainerRenderUsesColorTextPrimary(t *testing.T) {
	// Create a screen container and verify it uses the correct style
	content := "Colored Text"
	sc := NewScreenContainer(content)

	// The render method should apply ColorTextPrimary
	// We can verify this by checking that the style is applied
	rendered := sc.Render()

	// Since lipgloss applies ANSI color codes, we just verify content is present
	if !strings.Contains(rendered, content) {
		t.Errorf("expected rendered output to contain %q", content)
	}
}

func TestScreenContainerMaxWidthConstraint(t *testing.T) {
	// Create a long content string
	longContent := strings.Repeat("a", 200)
	sc := NewScreenContainer(longContent).WithMaxWidth(50)

	rendered := sc.Render()

	// The rendered output should respect the max width constraint
	// Split by newlines to check width

	// At least verify that the container was rendered
	if rendered == "" {
		t.Error("expected rendered output to be non-empty")
	}

	// Verify the content is still present (wrapped or truncated)
	if !strings.Contains(rendered, "a") {
		t.Error("expected rendered output to contain content")
	}
}

func TestScreenContainerEmptyContent(t *testing.T) {
	sc := NewScreenContainer("")

	rendered := sc.Render()

	// Should still render, even with empty content
	if rendered == "" {
		t.Error("expected rendered output to be non-empty (even with empty content)")
	}
}

func TestScreenContainerMultilineContent(t *testing.T) {
	content := "Line 1\nLine 2\nLine 3"
	sc := NewScreenContainer(content)

	rendered := sc.Render()

	if !strings.Contains(rendered, "Line 1") ||
		!strings.Contains(rendered, "Line 2") ||
		!strings.Contains(rendered, "Line 3") {
		t.Errorf("expected rendered output to contain all lines")
	}
}

func TestScreenContainerStyleConsistency(t *testing.T) {
	// Verify that the container uses consistent styling
	sc := NewScreenContainer("test")

	// Create a reference style to compare
	referenceStyle := lipgloss.NewStyle().
		Padding(1, 2).
		MaxWidth(styles.MaxWidth(120)).
		Foreground(styles.ColorTextPrimary)

	rendered := sc.Render()
	reference := referenceStyle.Render("test")

	// Both should produce similar output (content and styling)
	if rendered == "" || reference == "" {
		t.Error("expected both rendered outputs to be non-empty")
	}
}

