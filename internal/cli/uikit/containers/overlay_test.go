package containers_test

import (
	"strings"

	"github.com/baphled/kariya/internal/cli/uikit/containers"
	"github.com/charmbracelet/lipgloss"
	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
)

var _ = Describe("Overlay", func() {
	var overlay *containers.Overlay

	Describe("Construction", func() {
		It("should create overlay with dimensions", func() {
			overlay = containers.NewOverlay(80, 24)
			Expect(overlay).NotTo(BeNil())
		})
	})

	Describe("Configuration", func() {
		BeforeEach(func() {
			overlay = containers.NewOverlay(80, 24)
		})

		It("should set content to center", func() {
			result := overlay.Content("Hello World")
			Expect(result).To(Equal(overlay)) // Check chaining
		})

		It("should enable background dimming", func() {
			result := overlay.Dimmed()
			Expect(result).To(Equal(overlay)) // Check chaining
		})

		It("should set custom dim character", func() {
			result := overlay.DimmedWith('·')
			Expect(result).To(Equal(overlay)) // Check chaining
		})

		It("should support chaining", func() {
			result := overlay.
				Content("Test").
				Dimmed().
				DimmedWith('░')

			Expect(result).To(Equal(overlay))
		})
	})

	Describe("Rendering", func() {
		BeforeEach(func() {
			overlay = containers.NewOverlay(40, 10)
		})

		It("should center content horizontally", func() {
			overlay.Content("CENTER")
			rendered := overlay.Render()

			lines := strings.Split(rendered, "\n")
			// Find line with content
			for _, line := range lines {
				if strings.Contains(line, "CENTER") {
					// Content should be roughly in middle
					plainLine := stripAllANSI(line)
					centerPos := strings.Index(plainLine, "CENTER")
					lineWidth := len(plainLine)

					// Should be in the middle third of the line
					Expect(centerPos).To(BeNumerically(">", lineWidth/3))
					Expect(centerPos).To(BeNumerically("<", 2*lineWidth/3))
					break
				}
			}
		})

		It("should center content vertically", func() {
			overlay.Content("MIDDLE")
			rendered := overlay.Render()

			lines := strings.Split(rendered, "\n")

			// Find which line has the content
			contentLine := -1
			for i, line := range lines {
				if strings.Contains(line, "MIDDLE") {
					contentLine = i
					break
				}
			}

			// Should be in middle third of total lines
			Expect(contentLine).To(BeNumerically(">", len(lines)/3))
			Expect(contentLine).To(BeNumerically("<", 2*len(lines)/3))
		})

		It("should dim background when enabled", func() {
			overlay.Dimmed().Content("Text")
			rendered := overlay.Render()

			// Dimmed background should have repeated characters
			lines := strings.Split(rendered, "\n")

			// Count lines with dimming characters (should be all of them)
			dimmedLines := 0
			for _, line := range lines {
				plainLine := stripAllANSI(line)
				// Dimmed lines have repeated dot characters
				if strings.Count(plainLine, "·") > 3 {
					dimmedLines++
				}
			}

			Expect(dimmedLines).To(BeNumerically(">", 0))
		})

		It("should use custom dim character", func() {
			overlay.DimmedWith('░').Content("Test")
			rendered := overlay.Render()

			// Should have custom dim character
			Expect(rendered).To(ContainSubstring("░"))
		})

		It("should respect dimensions", func() {
			overlay.Content("Test")
			rendered := overlay.Render()

			lines := strings.Split(rendered, "\n")

			// Should have roughly the specified number of lines
			Expect(len(lines)).To(BeNumerically(">=", 8)) // Close to 10
			Expect(len(lines)).To(BeNumerically("<=", 12))

			// Each line should be close to specified width
			if len(lines) > 0 {
				width := lipgloss.Width(lines[0])
				Expect(width).To(BeNumerically(">=", 35)) // Close to 40
				Expect(width).To(BeNumerically("<=", 45))
			}
		})
	})
})

// Helper to strip all ANSI codes including CSI sequences
func stripAllANSI(s string) string {
	result := ""
	inEscape := false
	for i := 0; i < len(s); i++ {
		if s[i] == '\x1b' {
			inEscape = true
			continue
		}
		if inEscape {
			if s[i] == 'm' || (s[i] >= '@' && s[i] <= '~') {
				inEscape = false
			}
			continue
		}
		result += string(s[i])
	}
	return result
}
