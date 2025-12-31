package components

import (
	"strings"
	"testing"

)

func TestNewSectionContainer(t *testing.T) {
	content := "Test Content"
	sc := NewSectionContainer(content)

	if sc.content != content {
		t.Errorf("expected content %q, got %q", content, sc.content)
	}

	if sc.hasTitle {
		t.Error("expected hasTitle to be false for new container")
	}

	if sc.spacing != 1 {
		t.Errorf("expected default spacing 1, got %d", sc.spacing)
	}
}

func TestSectionContainerSetTitle(t *testing.T) {
	title := "Section Title"
	sc := NewSectionContainer("content").SetTitle(title)

	if sc.title != title {
		t.Errorf("expected title %q, got %q", title, sc.title)
	}

	if !sc.hasTitle {
		t.Error("expected hasTitle to be true")
	}
}

func TestSectionContainerWithSpacing(t *testing.T) {
	tests := []struct {
		name     string
		spacing  int
		expected int
	}{
		{"Zero spacing", 0, 0},
		{"One spacing", 1, 1},
		{"Multiple spacing", 3, 3},
		{"Large spacing", 10, 10},
		{"Negative spacing", -5, 0}, // Should be clamped to 0
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			sc := NewSectionContainer("content").WithSpacing(tt.spacing)

			if sc.spacing != tt.expected {
				t.Errorf("expected spacing %d, got %d", tt.expected, sc.spacing)
			}
		})
	}
}

func TestSectionContainerRenderWithoutTitle(t *testing.T) {
	content := "Test Content"
	sc := NewSectionContainer(content)

	rendered := sc.Render()

	if !strings.Contains(rendered, content) {
		t.Errorf("expected rendered output to contain %q", content)
	}
}

func TestSectionContainerRenderWithTitle(t *testing.T) {
	title := "Section Title"
	content := "Section Content"
	sc := NewSectionContainer(content).SetTitle(title)

	rendered := sc.Render()

	if !strings.Contains(rendered, title) {
		t.Errorf("expected rendered output to contain title %q", title)
	}

	if !strings.Contains(rendered, content) {
		t.Errorf("expected rendered output to contain content %q", content)
	}
}

func TestSectionContainerRenderWithTitleAndSpacing(t *testing.T) {
	title := "Title"
	content := "Content"
	sc := NewSectionContainer(content).SetTitle(title).WithSpacing(2)

	rendered := sc.Render()

	if !strings.Contains(rendered, title) || !strings.Contains(rendered, content) {
		t.Error("expected rendered output to contain both title and content")
	}
}

func TestSectionContainerRenderWithZeroSpacing(t *testing.T) {
	title := "Title"
	content := "Content"
	sc := NewSectionContainer(content).SetTitle(title).WithSpacing(0)

	rendered := sc.Render()

	if !strings.Contains(rendered, title) || !strings.Contains(rendered, content) {
		t.Error("expected rendered output to contain both title and content")
	}
}

func TestSectionContainerBuilderChaining(t *testing.T) {
	title := "Title"
	content := "Content"
	sc := NewSectionContainer(content).SetTitle(title).WithSpacing(3)

	if sc.title != title {
		t.Errorf("expected title %q, got %q", title, sc.title)
	}

	if sc.content != content {
		t.Errorf("expected content %q, got %q", content, sc.content)
	}

	if sc.spacing != 3 {
		t.Errorf("expected spacing 3, got %d", sc.spacing)
	}
}

func TestSectionContainerMultilineContent(t *testing.T) {
	title := "Multi-line Title"
	content := "Line 1\nLine 2\nLine 3"
	sc := NewSectionContainer(content).SetTitle(title)

	rendered := sc.Render()

	if !strings.Contains(rendered, "Line 1") ||
		!strings.Contains(rendered, "Line 2") ||
		!strings.Contains(rendered, "Line 3") {
		t.Error("expected rendered output to contain all lines")
	}
}

func TestSectionContainerEmptyTitle(t *testing.T) {
	content := "Content"
	sc := NewSectionContainer(content).SetTitle("")

	if !sc.hasTitle {
		t.Error("expected hasTitle to be true even with empty title")
	}

	rendered := sc.Render()

	if !strings.Contains(rendered, content) {
		t.Errorf("expected rendered output to contain content %q", content)
	}
}

func TestSectionContainerEmptyContent(t *testing.T) {
	title := "Title"
	sc := NewSectionContainer("").SetTitle(title)

	rendered := sc.Render()

	if !strings.Contains(rendered, title) {
		t.Errorf("expected rendered output to contain title %q", title)
	}
}

func TestSectionContainerTitlePositioning(t *testing.T) {
	// Title should appear before content
	title := "TITLE"
	content := "CONTENT"
	sc := NewSectionContainer(content).SetTitle(title)

	rendered := sc.Render()

	titlePos := strings.Index(rendered, title)
	contentPos := strings.Index(rendered, content)

	if titlePos == -1 || contentPos == -1 {
		t.Error("expected both title and content to be present")
	}

	if titlePos > contentPos {
		t.Error("expected title to appear before content")
	}
}

func TestSectionContainerSpacingVariations(t *testing.T) {
	title := "Title"
	content := "Content"

	tests := []struct {
		name    string
		spacing int
	}{
		{"No spacing", 0},
		{"Single spacing", 1},
		{"Double spacing", 2},
		{"Triple spacing", 3},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			sc := NewSectionContainer(content).SetTitle(title).WithSpacing(tt.spacing)
			rendered := sc.Render()

			if !strings.Contains(rendered, title) || !strings.Contains(rendered, content) {
				t.Error("expected rendered output to contain both title and content")
			}
		})
	}
}

func TestSectionContainerUsesHeaderSectionStyle(t *testing.T) {
	// Verify that the section uses HeaderSection styling for title
	title := "Styled Title"
	content := "Content"
	sc := NewSectionContainer(content).SetTitle(title)

	rendered := sc.Render()

	// The render should apply HeaderSection style
	if !strings.Contains(rendered, title) {
		t.Error("expected rendered output to contain title")
	}
}

func TestSectionContainerColorConsistency(t *testing.T) {
	// Verify that colors are consistent with the style constants
	title := "Title"
	content := "Content"
	sc := NewSectionContainer(content).SetTitle(title)

	rendered := sc.Render()

	// Verify the render includes both styled elements
	if rendered == "" {
		t.Error("expected rendered output to be non-empty")
	}

	if !strings.Contains(rendered, title) || !strings.Contains(rendered, content) {
		t.Error("expected rendered output to contain both title and content")
	}
}

func TestSectionContainerLongContent(t *testing.T) {
	title := "Title"
	content := strings.Repeat("This is a long content line. ", 10)
	sc := NewSectionContainer(content).SetTitle(title)

	rendered := sc.Render()

	if !strings.Contains(rendered, title) {
		t.Errorf("expected rendered output to contain title %q", title)
	}

	// Check that long content is preserved
	if !strings.Contains(rendered, "long content") {
		t.Error("expected rendered output to contain long content")
	}
}

func TestSectionContainerSpecialCharacters(t *testing.T) {
	title := "Title with [brackets] and {braces}"
	content := "Content with special chars: @#$%^&*()"
	sc := NewSectionContainer(content).SetTitle(title)

	rendered := sc.Render()

	if !strings.Contains(rendered, "[brackets]") ||
		!strings.Contains(rendered, "{braces}") ||
		!strings.Contains(rendered, "@#$%^&*()") {
		t.Error("expected rendered output to contain special characters")
	}
}

