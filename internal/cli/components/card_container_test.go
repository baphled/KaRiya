package components

import (
	"strings"
	"testing"

	"github.com/baphled/kariya/internal/cli/styles"
)

func TestNewCardContainer(t *testing.T) {
	cc := NewCardContainer()

	if cc.backgroundColor != styles.ColorBackgroundCard {
		t.Errorf("expected backgroundColor %v, got %v", styles.ColorBackgroundCard, cc.backgroundColor)
	}

	if cc.borderColor != styles.ColorBorder {
		t.Errorf("expected borderColor %v, got %v", styles.ColorBorder, cc.borderColor)
	}

	if cc.hasHeader || cc.hasBody || cc.hasFooter {
		t.Error("expected new card to have no sections by default")
	}
}

func TestCardContainerSetHeader(t *testing.T) {
	header := "Card Header"
	cc := NewCardContainer().SetHeader(header)

	if cc.header != header {
		t.Errorf("expected header %q, got %q", header, cc.header)
	}

	if !cc.hasHeader {
		t.Error("expected hasHeader to be true")
	}
}

func TestCardContainerSetBody(t *testing.T) {
	body := "Card Body"
	cc := NewCardContainer().SetBody(body)

	if cc.body != body {
		t.Errorf("expected body %q, got %q", body, cc.body)
	}

	if !cc.hasBody {
		t.Error("expected hasBody to be true")
	}
}

func TestCardContainerSetFooter(t *testing.T) {
	footer := "Card Footer"
	cc := NewCardContainer().SetFooter(footer)

	if cc.footer != footer {
		t.Errorf("expected footer %q, got %q", footer, cc.footer)
	}

	if !cc.hasFooter {
		t.Error("expected hasFooter to be true")
	}
}

func TestCardContainerWithBackgroundColor(t *testing.T) {
	customColor := styles.ColorBackgroundAlt
	cc := NewCardContainer().WithBackgroundColor(customColor)

	if cc.backgroundColor != customColor {
		t.Errorf("expected backgroundColor %v, got %v", customColor, cc.backgroundColor)
	}
}

func TestCardContainerWithBorderColor(t *testing.T) {
	customColor := styles.ColorBorderActive
	cc := NewCardContainer().WithBorderColor(customColor)

	if cc.borderColor != customColor {
		t.Errorf("expected borderColor %v, got %v", customColor, cc.borderColor)
	}
}

func TestCardContainerRenderEmpty(t *testing.T) {
	cc := NewCardContainer()
	rendered := cc.Render()

	// Should still render a card, even with no sections
	if rendered == "" {
		t.Error("expected rendered output to be non-empty")
	}
}

func TestCardContainerRenderWithHeader(t *testing.T) {
	header := "Test Header"
	cc := NewCardContainer().SetHeader(header)

	rendered := cc.Render()

	if !strings.Contains(rendered, header) {
		t.Errorf("expected rendered output to contain header %q", header)
	}
}

func TestCardContainerRenderWithBody(t *testing.T) {
	body := "Test Body Content"
	cc := NewCardContainer().SetBody(body)

	rendered := cc.Render()

	if !strings.Contains(rendered, body) {
		t.Errorf("expected rendered output to contain body %q", body)
	}
}

func TestCardContainerRenderWithFooter(t *testing.T) {
	footer := "Test Footer"
	cc := NewCardContainer().SetFooter(footer)

	rendered := cc.Render()

	if !strings.Contains(rendered, footer) {
		t.Errorf("expected rendered output to contain footer %q", footer)
	}
}

func TestCardContainerRenderAllSections(t *testing.T) {
	header := "Header"
	body := "Body"
	footer := "Footer"

	cc := NewCardContainer().
		SetHeader(header).
		SetBody(body).
		SetFooter(footer)

	rendered := cc.Render()

	if !strings.Contains(rendered, header) {
		t.Errorf("expected rendered output to contain header %q", header)
	}

	if !strings.Contains(rendered, body) {
		t.Errorf("expected rendered output to contain body %q", body)
	}

	if !strings.Contains(rendered, footer) {
		t.Errorf("expected rendered output to contain footer %q", footer)
	}
}

func TestCardContainerBuilderChaining(t *testing.T) {
	header := "Header"
	body := "Body"
	footer := "Footer"

	cc := NewCardContainer().
		SetHeader(header).
		SetBody(body).
		SetFooter(footer).
		WithBackgroundColor(styles.ColorBackgroundAlt).
		WithBorderColor(styles.ColorBorderActive)

	if cc.header != header || cc.body != body || cc.footer != footer {
		t.Error("expected all sections to be set after chaining")
	}

	if cc.backgroundColor != styles.ColorBackgroundAlt {
		t.Error("expected background color to be set")
	}

	if cc.borderColor != styles.ColorBorderActive {
		t.Error("expected border color to be set")
	}
}

func TestCardContainerRenderWithCustomColors(t *testing.T) {
	body := "Colored Card"
	cc := NewCardContainer().
		SetBody(body).
		WithBackgroundColor(styles.ColorBackgroundAlt).
		WithBorderColor(styles.ColorBorderActive)

	rendered := cc.Render()

	if !strings.Contains(rendered, body) {
		t.Errorf("expected rendered output to contain %q", body)
	}
}

func TestCardContainerMultilineContent(t *testing.T) {
	header := "Multi-line\nHeader"
	body := "Multi-line\nBody"
	footer := "Multi-line\nFooter"

	cc := NewCardContainer().
		SetHeader(header).
		SetBody(body).
		SetFooter(footer)

	rendered := cc.Render()

	// Verify all lines are present
	if !strings.Contains(rendered, "Header") || !strings.Contains(rendered, "Body") ||
		!strings.Contains(rendered, "Footer") {
		t.Error("expected rendered output to contain all multiline content")
	}
}

func TestCardContainerRenderDestructiveStyle(t *testing.T) {
	body := "Destructive Action"
	cc := NewCardContainer().
		SetBody(body).
		WithBorderColor(styles.ColorBorderError)

	rendered := cc.Render()

	if !strings.Contains(rendered, body) {
		t.Errorf("expected rendered output to contain %q", body)
	}
}

func TestCardContainerRenderHeaderOnly(t *testing.T) {
	header := "Header Only"
	cc := NewCardContainer().SetHeader(header)

	rendered := cc.Render()

	if !strings.Contains(rendered, header) {
		t.Errorf("expected rendered output to contain header %q", header)
	}

	// Body and footer should not be in the output
	if strings.Contains(rendered, "Body") || strings.Contains(rendered, "Footer") {
		t.Error("expected rendered output to not contain body or footer")
	}
}

func TestCardContainerRenderBodyOnly(t *testing.T) {
	body := "Body Only"
	cc := NewCardContainer().SetBody(body)

	rendered := cc.Render()

	if !strings.Contains(rendered, body) {
		t.Errorf("expected rendered output to contain body %q", body)
	}
}

func TestCardContainerRenderFooterOnly(t *testing.T) {
	footer := "Footer Only"
	cc := NewCardContainer().SetFooter(footer)

	rendered := cc.Render()

	if !strings.Contains(rendered, footer) {
		t.Errorf("expected rendered output to contain footer %q", footer)
	}
}

func TestCardContainerRenderHeaderAndBody(t *testing.T) {
	header := "Header"
	body := "Body"
	cc := NewCardContainer().SetHeader(header).SetBody(body)

	rendered := cc.Render()

	if !strings.Contains(rendered, header) || !strings.Contains(rendered, body) {
		t.Error("expected rendered output to contain both header and body")
	}

	if strings.Contains(rendered, "Footer") {
		t.Error("expected rendered output to not contain footer")
	}
}

func TestCardContainerEmptyContent(t *testing.T) {
	cc := NewCardContainer().
		SetHeader("").
		SetBody("").
		SetFooter("")

	rendered := cc.Render()

	// Should still render a card
	if rendered == "" {
		t.Error("expected rendered output to be non-empty")
	}
}

func TestCardContainerBorderAndBackgroundStyling(t *testing.T) {
	// Verify that the card applies border and background styling
	cc := NewCardContainer().SetBody("Test")

	// The render should apply lipgloss styles
	rendered := cc.Render()

	if rendered == "" {
		t.Error("expected rendered output to be non-empty")
	}

	// Verify the body content is present
	if !strings.Contains(rendered, "Test") {
		t.Error("expected rendered output to contain body content")
	}
}

func TestCardContainerSectionOrdering(t *testing.T) {
	// Sections should appear in order: header, body, footer
	header := "HEADER"
	body := "BODY"
	footer := "FOOTER"

	cc := NewCardContainer().
		SetHeader(header).
		SetBody(body).
		SetFooter(footer)

	rendered := cc.Render()

	// Find positions of each section
	headerPos := strings.Index(rendered, header)
	bodyPos := strings.Index(rendered, body)
	footerPos := strings.Index(rendered, footer)

	if headerPos == -1 || bodyPos == -1 || footerPos == -1 {
		t.Error("expected all sections to be present")
	}

	if headerPos > bodyPos || bodyPos > footerPos {
		t.Error("expected sections to appear in order: header, body, footer")
	}
}

