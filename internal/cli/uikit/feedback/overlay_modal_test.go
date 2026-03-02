package feedback_test

import (
	"strings"

	"github.com/baphled/kariya/internal/cli/themes"
	"github.com/baphled/kariya/internal/cli/uikit/feedback"
	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
)

var _ = Describe("OverlayModal", func() {
	var modal *feedback.OverlayModal

	BeforeEach(func() {
		modal = feedback.NewOverlayModal("Test Title", "Test Content")
	})

	Describe("WithTheme", func() {
		It("returns the modal for chaining", func() {
			theme := themes.NewDefaultTheme()
			result := modal.WithTheme(theme)
			Expect(result).To(BeIdenticalTo(modal))
		})

		It("renders correctly after setting theme", func() {
			theme := themes.NewDefaultTheme()
			modal.WithTheme(theme)
			background := strings.Repeat("Background line\n", 20)
			output := modal.RenderCentered(background, 80, 24)
			Expect(output).NotTo(BeEmpty())
		})
	})

	Describe("RenderCentered", func() {
		It("renders modal over dimmed background", func() {
			background := strings.Repeat("Background content\n", 20)
			output := modal.RenderCentered(background, 80, 24)
			Expect(output).NotTo(BeEmpty())
			Expect(output).To(ContainSubstring("Test Title"))
		})

		It("includes content in the overlay", func() {
			background := strings.Repeat("Line of text\n", 20)
			output := modal.RenderCentered(background, 100, 30)
			Expect(output).To(ContainSubstring("Test Content"))
		})

		It("includes footer when set", func() {
			modal.SetFooter("Press Esc to close")
			background := strings.Repeat("BG\n", 20)
			output := modal.RenderCentered(background, 80, 24)
			Expect(output).To(ContainSubstring("Press Esc to close"))
		})

		It("renders with nil theme falling back to default", func() {
			modal.WithTheme(nil)
			background := strings.Repeat("Content\n", 20)
			output := modal.RenderCentered(background, 80, 24)
			Expect(output).NotTo(BeEmpty())
		})

		It("renders with only title and no content", func() {
			emptyContentModal := feedback.NewOverlayModal("Title Only", "")
			background := strings.Repeat("BG\n", 20)
			output := emptyContentModal.RenderCentered(background, 80, 24)
			Expect(output).To(ContainSubstring("Title Only"))
		})

		It("renders with only content and no title", func() {
			noTitleModal := feedback.NewOverlayModal("", "Content Only")
			background := strings.Repeat("BG\n", 20)
			output := noTitleModal.RenderCentered(background, 80, 24)
			Expect(output).To(ContainSubstring("Content Only"))
		})
	})
})

var _ = Describe("Modal WithTheme", func() {
	It("sets the theme and returns the modal for chaining", func() {
		modal := feedback.NewErrorModal("Error", "Something failed")
		theme := themes.NewDefaultTheme()
		result := modal.WithTheme(theme)
		Expect(result).To(BeIdenticalTo(modal))
	})

	It("renders correctly with a custom theme", func() {
		modal := feedback.NewLoadingModal("Processing...", false)
		theme := themes.NewDefaultTheme()
		modal.WithTheme(theme)
		output := modal.Render(80, 24)
		Expect(output).NotTo(BeEmpty())
	})
})

var _ = Describe("Modal RotateMessage without rotator", func() {
	It("returns the original message when no rotator is set", func() {
		modal := feedback.NewErrorModal("Error", "Original message")
		result := modal.RotateMessage()
		Expect(result).To(Equal("Original message"))
	})
})
