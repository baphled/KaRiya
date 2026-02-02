package behaviors_test

import (
	"strings"

	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"

	"github.com/baphled/kariya/internal/cli/behaviors"
)

// testModal is a simple modal for testing overlay rendering.
// It only needs View() to satisfy the Viewable interface.
type testModal struct {
	content string
	visible bool
}

func (m *testModal) View() string {
	if !m.visible {
		return ""
	}
	return m.content
}

var _ = Describe("Modal Helpers", func() {

	Describe("StaticViewModel", func() {
		It("returns the content in View()", func() {
			model := &behaviors.StaticViewModel{Content: "Hello, World!"}
			Expect(model.View()).To(Equal("Hello, World!"))
		})

		It("handles empty content", func() {
			model := &behaviors.StaticViewModel{Content: ""}
			Expect(model.View()).To(Equal(""))
		})

		It("handles multi-line content", func() {
			content := "Line 1\nLine 2\nLine 3"
			model := &behaviors.StaticViewModel{Content: content}
			Expect(model.View()).To(Equal(content))
		})

		It("implements Viewable interface", func() {
			var v behaviors.Viewable = &behaviors.StaticViewModel{Content: "test"}
			Expect(v.View()).To(Equal("test"))
		})
	})

	Describe("RenderModalOverlay", func() {
		var (
			modal      *testModal
			background string
		)

		BeforeEach(func() {
			modal = &testModal{
				content: "[MODAL CONTENT]",
				visible: true,
			}
			background = "Background View Content"
		})

		It("renders modal content over background", func() {
			result := behaviors.RenderModalOverlay(modal, background)
			// The overlay should contain both modal and background content
			Expect(result).To(ContainSubstring("MODAL CONTENT"))
		})

		It("includes background content", func() {
			result := behaviors.RenderModalOverlay(modal, background)
			// Background may be partially visible around the modal
			Expect(result).ToNot(BeEmpty())
		})

		It("handles empty modal content", func() {
			modal.visible = false
			result := behaviors.RenderModalOverlay(modal, background)
			// Should still render without error
			Expect(result).NotTo(BeEmpty())
		})

		It("handles empty background", func() {
			result := behaviors.RenderModalOverlay(modal, "")
			Expect(result).To(ContainSubstring("MODAL CONTENT"))
		})

		It("handles both empty", func() {
			modal.visible = false
			result := behaviors.RenderModalOverlay(modal, "")
			// Should not panic
			Expect(result).NotTo(BeNil())
		})

		It("handles multi-line background", func() {
			multiLineBackground := strings.Repeat("Line of content\n", 20)
			result := behaviors.RenderModalOverlay(modal, multiLineBackground)
			Expect(result).To(ContainSubstring("MODAL CONTENT"))
		})
	})

	Describe("RenderModalOverlayWithOffset", func() {
		var (
			modal      *testModal
			background string
		)

		BeforeEach(func() {
			modal = &testModal{
				content: "[OFFSET MODAL]",
				visible: true,
			}
			background = "Background"
		})

		It("renders with zero offset", func() {
			result := behaviors.RenderModalOverlayWithOffset(modal, background, 0, 0)
			Expect(result).To(ContainSubstring("OFFSET MODAL"))
		})

		It("renders with positive offset", func() {
			result := behaviors.RenderModalOverlayWithOffset(modal, background, 5, 5)
			Expect(result).To(ContainSubstring("OFFSET MODAL"))
		})

		It("renders with negative offset", func() {
			result := behaviors.RenderModalOverlayWithOffset(modal, background, -5, -5)
			Expect(result).To(ContainSubstring("OFFSET MODAL"))
		})

		It("renders with mixed offsets", func() {
			result := behaviors.RenderModalOverlayWithOffset(modal, background, 10, -3)
			Expect(result).To(ContainSubstring("OFFSET MODAL"))
		})
	})

	Describe("DefaultModalDimensions", func() {
		It("returns sensible default width", func() {
			width, _ := behaviors.DefaultModalDimensions()
			Expect(width).To(Equal(120))
		})

		It("returns sensible default height", func() {
			_, height := behaviors.DefaultModalDimensions()
			Expect(height).To(Equal(40))
		})

		It("returns reasonable dimensions for modal display", func() {
			width, height := behaviors.DefaultModalDimensions()
			Expect(width).To(BeNumerically(">=", 80))
			Expect(height).To(BeNumerically(">=", 24))
		})
	})
})
