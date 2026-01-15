package primitives_test

import (
	"testing"

	"github.com/baphled/kariya/internal/cli/uikit/primitives"
	"github.com/baphled/kariya/internal/cli/uikit/theme"
	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
)

func TestPrimitives(t *testing.T) {
	RegisterFailHandler(Fail)
	RunSpecs(t, "UIKit Primitives Suite")
}

var _ = Describe("Text", func() {
	var th theme.Theme

	BeforeEach(func() {
		th = theme.Default()
	})

	Describe("NewText", func() {
		It("should create text with content", func() {
			text := primitives.NewText("Hello World", th)
			Expect(text).NotTo(BeNil())
			rendered := text.Render()
			Expect(rendered).To(ContainSubstring("Hello World"))
		})

		It("should accept nil theme and use default", func() {
			text := primitives.NewText("Hello", nil)
			Expect(text).NotTo(BeNil())
			rendered := text.Render()
			Expect(rendered).To(ContainSubstring("Hello"))
		})
	})

	Describe("Fluent API", func() {
		Describe("Style", func() {
			It("should set TextTitle style", func() {
				text := primitives.NewText("Title", th).Style(primitives.TextTitle)
				rendered := text.Render()
				Expect(rendered).To(ContainSubstring("Title"))
			})

			It("should set TextSubtitle style", func() {
				text := primitives.NewText("Subtitle", th).Style(primitives.TextSubtitle)
				rendered := text.Render()
				Expect(rendered).To(ContainSubstring("Subtitle"))
			})

			It("should set TextBody style", func() {
				text := primitives.NewText("Body", th).Style(primitives.TextBody)
				rendered := text.Render()
				Expect(rendered).To(ContainSubstring("Body"))
			})

			It("should set TextMuted style", func() {
				text := primitives.NewText("Muted", th).Style(primitives.TextMuted)
				rendered := text.Render()
				Expect(rendered).To(ContainSubstring("Muted"))
			})

			It("should set TextError style", func() {
				text := primitives.NewText("Error", th).Style(primitives.TextError)
				rendered := text.Render()
				Expect(rendered).To(ContainSubstring("Error"))
			})

			It("should set TextSuccess style", func() {
				text := primitives.NewText("Success", th).Style(primitives.TextSuccess)
				rendered := text.Render()
				Expect(rendered).To(ContainSubstring("Success"))
			})

			It("should set TextWarning style", func() {
				text := primitives.NewText("Warning", th).Style(primitives.TextWarning)
				rendered := text.Render()
				Expect(rendered).To(ContainSubstring("Warning"))
			})
		})

		Describe("Bold", func() {
			It("should make text bold", func() {
				text := primitives.NewText("Bold Text", th).Bold()
				rendered := text.Render()
				Expect(rendered).To(ContainSubstring("Bold Text"))
			})

			It("should chain with other methods", func() {
				text := primitives.NewText("Bold Title", th).Style(primitives.TextTitle).Bold()
				rendered := text.Render()
				Expect(rendered).To(ContainSubstring("Bold Title"))
			})
		})

		Describe("Width", func() {
			It("should constrain output width", func() {
				text := primitives.NewText("This is a very long text that should be constrained", th).Width(20)
				rendered := text.Render()
				// Split by newlines to check each line width
				lines := splitLines(rendered)
				for _, line := range lines {
					// Use lipgloss.Width to measure actual width (ignoring ANSI codes)
					Expect(visualWidth(line)).To(BeNumerically("<=", 20))
				}
			})

			It("should chain with other methods", func() {
				text := primitives.NewText("Short", th).Bold().Width(10)
				rendered := text.Render()
				Expect(rendered).To(ContainSubstring("Short"))
			})
		})
	})

	Describe("Convenience Constructors", func() {
		Describe("Title", func() {
			It("should create title text", func() {
				text := primitives.Title("Page Title", th)
				rendered := text.Render()
				Expect(rendered).To(ContainSubstring("Page Title"))
			})
		})

		Describe("Subtitle", func() {
			It("should create subtitle text", func() {
				text := primitives.Subtitle("Section Subtitle", th)
				rendered := text.Render()
				Expect(rendered).To(ContainSubstring("Section Subtitle"))
			})
		})

		Describe("Body", func() {
			It("should create body text", func() {
				text := primitives.Body("Regular text", th)
				rendered := text.Render()
				Expect(rendered).To(ContainSubstring("Regular text"))
			})
		})

		Describe("Muted", func() {
			It("should create muted text", func() {
				text := primitives.Muted("Disabled text", th)
				rendered := text.Render()
				Expect(rendered).To(ContainSubstring("Disabled text"))
			})
		})

		Describe("ErrorText", func() {
			It("should create error text", func() {
				text := primitives.ErrorText("Error message", th)
				rendered := text.Render()
				Expect(rendered).To(ContainSubstring("Error message"))
			})
		})

		Describe("SuccessText", func() {
			It("should create success text", func() {
				text := primitives.SuccessText("Success message", th)
				rendered := text.Render()
				Expect(rendered).To(ContainSubstring("Success message"))
			})
		})

		Describe("WarningText", func() {
			It("should create warning text", func() {
				text := primitives.WarningText("Warning message", th)
				rendered := text.Render()
				Expect(rendered).To(ContainSubstring("Warning message"))
			})
		})
	})
})

// Helper functions for testing
func splitLines(s string) []string {
	var lines []string
	current := ""
	for _, ch := range s {
		if ch == '\n' {
			lines = append(lines, current)
			current = ""
		} else {
			current += string(ch)
		}
	}
	if current != "" {
		lines = append(lines, current)
	}
	return lines
}

func visualWidth(s string) int {
	// Simple visual width calculation (ignoring ANSI codes)
	// For production, use lipgloss.Width() but this is sufficient for tests
	width := 0
	inEscape := false
	for _, ch := range s {
		if ch == '\x1b' {
			inEscape = true
			continue
		}
		if inEscape {
			if ch == 'm' {
				inEscape = false
			}
			continue
		}
		width++
	}
	return width
}
