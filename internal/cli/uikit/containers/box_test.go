package containers_test

import (
	"testing"

	"github.com/baphled/kariya/internal/cli/themes"
	"github.com/baphled/kariya/internal/cli/uikit/containers"
	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
)

func TestContainers(t *testing.T) {
	RegisterFailHandler(Fail)
	RunSpecs(t, "Containers Suite")
}

var _ = Describe("Box", func() {
	var (
		theme themes.Theme
		box   *containers.Box
	)

	BeforeEach(func() {
		theme = themes.NewDefaultTheme()
	})

	Describe("Construction", func() {
		It("should create box with theme", func() {
			box = containers.NewBox(theme)
			Expect(box).NotTo(BeNil())
		})
	})

	Describe("Configuration", func() {
		BeforeEach(func() {
			box = containers.NewBox(theme)
		})

		It("should set content", func() {
			result := box.Content("Hello World")
			Expect(result).To(Equal(box)) // Check chaining
		})

		It("should set title", func() {
			result := box.Title("My Box")
			Expect(result).To(Equal(box)) // Check chaining
		})

		It("should set variant to Default", func() {
			result := box.Variant(containers.BoxDefault)
			Expect(result).To(Equal(box))
		})

		It("should set variant to Emphasized", func() {
			result := box.Variant(containers.BoxEmphasized)
			Expect(result).To(Equal(box))
		})

		It("should set variant to Destructive", func() {
			result := box.Variant(containers.BoxDestructive)
			Expect(result).To(Equal(box))
		})

		It("should set variant to Subtle", func() {
			result := box.Variant(containers.BoxSubtle)
			Expect(result).To(Equal(box))
		})

		It("should set width (0 = auto)", func() {
			result := box.Width(50)
			Expect(result).To(Equal(box))
		})

		It("should set height (0 = auto)", func() {
			result := box.Height(20)
			Expect(result).To(Equal(box))
		})

		It("should set padding", func() {
			result := box.Padding(2)
			Expect(result).To(Equal(box))
		})

		It("should enable shadow", func() {
			result := box.WithShadow()
			Expect(result).To(Equal(box))
		})

		It("should support chaining", func() {
			result := box.
				Title("Test").
				Content("Content").
				Variant(containers.BoxEmphasized).
				Width(60).
				Height(10).
				Padding(1).
				WithShadow()

			Expect(result).To(Equal(box))
		})
	})

	Describe("Rendering", func() {
		BeforeEach(func() {
			box = containers.NewBox(theme)
		})

		It("should produce bordered box", func() {
			box.Content("Test content")
			rendered := box.Render()

			Expect(rendered).NotTo(BeEmpty())
			Expect(rendered).To(ContainSubstring("Test content"))
		})

		It("should render title when set", func() {
			box.Title("My Title").Content("Content")
			rendered := box.Render()

			Expect(rendered).To(ContainSubstring("My Title"))
		})

		It("should use default border for Default variant", func() {
			box.Variant(containers.BoxDefault).Content("Content")
			rendered := box.Render()

			// Should have border characters
			Expect(rendered).To(MatchRegexp(`[─│┌┐└┘]`))
		})

		It("should use thick border for Emphasized variant", func() {
			box.Variant(containers.BoxEmphasized).Content("Important")
			rendered := box.Render()

			// Emphasized should have thick border
			Expect(rendered).To(MatchRegexp(`[━┃┏┓┗┛]`))
		})

		It("should use error color for Destructive variant", func() {
			box.Variant(containers.BoxDestructive).Content("Warning")
			rendered := box.Render()

			// Should contain ANSI escape codes for color (destructive uses red)
			Expect(rendered).To(ContainSubstring("\x1b["))
		})

		It("should render shadow when enabled", func() {
			box.WithShadow().Content("Shadowed")
			rendered := box.Render()

			// Shadow adds extra characters/spacing
			Expect(len(rendered)).To(BeNumerically(">", len("Shadowed")+10))
		})

		It("should respect width constraint", func() {
			box.Width(30).Content("This is a test")
			rendered := box.Render()

			// Rendered width should be around 30 (accounting for borders and ANSI codes)
			lines := splitLines(rendered)
			if len(lines) > 0 {
				// Strip ANSI codes for width measurement
				cleanLine := stripANSI(lines[0])
				Expect(len(cleanLine)).To(BeNumerically("<=", 35)) // Width + borders
			}
		})

		It("should apply padding", func() {
			box.Padding(2).Content("Padded")
			rendered := box.Render()

			// Padded content should have more whitespace
			Expect(rendered).To(ContainSubstring("Padded"))
			Expect(len(rendered)).To(BeNumerically(">", len("Padded")+20))
		})
	})
})

// Helper to split rendered output into lines
func splitLines(s string) []string {
	lines := []string{}
	current := ""
	for _, r := range s {
		if r == '\n' {
			lines = append(lines, current)
			current = ""
		} else {
			current += string(r)
		}
	}
	if current != "" {
		lines = append(lines, current)
	}
	return lines
}

// Helper to strip ANSI escape codes
func stripANSI(s string) string {
	result := ""
	inEscape := false
	for _, r := range s {
		if r == '\x1b' {
			inEscape = true
			continue
		}
		if inEscape {
			if r == 'm' {
				inEscape = false
			}
			continue
		}
		result += string(r)
	}
	return result
}
