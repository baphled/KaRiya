package components

import (
	"strings"

	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
)

var _ = Describe("Modal Overlay System", func() {
	Describe("RenderOverlay", func() {
		It("renders modal content over background", func() {
			background := "This is the background content\nLine 2 of background\nLine 3 of background"
			modalContent := "Modal Content"

			output := RenderOverlay(background, modalContent, 80, 24)

			Expect(output).NotTo(BeEmpty())
			Expect(output).To(ContainSubstring("Modal Content"))
		})

		It("handles dimmed background", func() {
			background := "Background text that should be dimmed"
			modalContent := "Modal"

			output := RenderOverlay(background, modalContent, 80, 24)

			Expect(output).NotTo(BeEmpty())
		})

		It("handles small terminal gracefully", func() {
			background := "Background"
			modalContent := "A fairly long modal content that might not fit well in a small terminal"

			output := RenderOverlay(background, modalContent, 40, 10)

			Expect(output).NotTo(BeEmpty())
		})

		It("handles large terminal", func() {
			background := strings.Repeat("Background line\n", 50)
			modalContent := "Modal content"

			output := RenderOverlay(background, modalContent, 200, 60)

			Expect(output).NotTo(BeEmpty())
		})

		It("renders modal with empty background", func() {
			background := ""
			modalContent := "Modal content"

			output := RenderOverlay(background, modalContent, 80, 24)

			Expect(output).To(ContainSubstring("Modal content"))
		})

		It("handles empty modal content gracefully", func() {
			background := "Background content"
			modalContent := ""

			output := RenderOverlay(background, modalContent, 80, 24)

			Expect(output).NotTo(BeEmpty())
		})

		It("renders multiline modal", func() {
			background := "Background"
			modalContent := "Line 1\nLine 2\nLine 3\nLine 4"

			output := RenderOverlay(background, modalContent, 80, 24)

			Expect(output).NotTo(BeEmpty())
		})
	})

	Describe("OverlayModal", func() {
		Describe("NewOverlayModal", func() {
			It("creates a modal with the given title and content", func() {
				modal := NewOverlayModal("Test Title", "Test Content")

				Expect(modal).NotTo(BeNil())
				Expect(modal.Title).To(Equal("Test Title"))
				Expect(modal.Content).To(Equal("Test Content"))
			})

			It("sets default width within bounds", func() {
				modal := NewOverlayModal("Title", "Content")

				Expect(modal.Width).To(BeNumerically(">=", MinOverlayWidth))
				Expect(modal.Width).To(BeNumerically("<=", MaxOverlayWidth))
			})
		})

		Describe("SetWidth", func() {
			It("sets valid width", func() {
				modal := NewOverlayModal("Title", "Content")

				modal.SetWidth(80)

				Expect(modal.Width).To(Equal(80))
			})

			It("clamps width below minimum to 40", func() {
				modal := NewOverlayModal("Title", "Content")

				modal.SetWidth(20)

				Expect(modal.Width).To(Equal(MinOverlayWidth))
			})

			It("clamps width above maximum to 120", func() {
				modal := NewOverlayModal("Title", "Content")

				modal.SetWidth(150)

				Expect(modal.Width).To(Equal(MaxOverlayWidth))
			})
		})

		Describe("SetFooter", func() {
			It("sets the footer text", func() {
				modal := NewOverlayModal("Title", "Content")

				modal.SetFooter("Press Enter to confirm | Esc to cancel")

				Expect(modal.Footer).To(Equal("Press Enter to confirm | Esc to cancel"))
			})
		})

		Describe("RenderCentered", func() {
			It("renders the modal centered over background", func() {
				modal := NewOverlayModal("Test", "Modal Content")
				background := strings.Repeat("X", 80*24)

				output := modal.RenderCentered(background, 80, 24)

				Expect(output).NotTo(BeEmpty())
				Expect(output).To(ContainSubstring("Modal Content"))
			})

			It("includes footer in rendered output", func() {
				modal := NewOverlayModal("Title", "Content")
				modal.SetFooter("Enter: Confirm | Esc: Cancel")

				output := modal.RenderCentered("Background", 80, 24)

				Expect(output).To(ContainSubstring("Enter: Confirm"))
			})
		})
	})

	Describe("DimContent", func() {
		It("returns empty string for empty content", func() {
			dimmed := DimContent("")

			Expect(dimmed).To(BeEmpty())
		})

		It("returns non-empty string for single line", func() {
			dimmed := DimContent("Single line")

			Expect(dimmed).NotTo(BeEmpty())
		})

		It("preserves line count for multiline content", func() {
			original := "Line 1\nLine 2\nLine 3"
			dimmed := DimContent(original)

			Expect(dimmed).NotTo(BeEmpty())

			originalLines := strings.Count(original, "\n")
			dimmedLines := strings.Count(dimmed, "\n")
			Expect(dimmedLines).To(Equal(originalLines))
		})
	})
})
