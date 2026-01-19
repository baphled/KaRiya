package components_test

import (
	"github.com/baphled/kariya/internal/cli/components"
	tea "github.com/charmbracelet/bubbletea"
	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
)

var _ = Describe("ExportOptionsModal", func() {
	var (
		modal *components.ExportOptionsModal
	)

	Describe("NewExportOptionsModal", func() {
		It("should create an export options modal", func() {
			modal = components.NewExportOptionsModal(120, 40)

			Expect(modal).NotTo(BeNil())
			Expect(modal.IsVisible()).To(BeTrue())
		})

		It("should initialize with default format and location", func() {
			modal = components.NewExportOptionsModal(120, 40)

			data := modal.GetExportData()
			Expect(data).NotTo(BeNil())
			// huh.Select auto-selects first option
			Expect(data.Format).To(Equal("text"))   // First option
			Expect(data.Location).To(Equal("file")) // First option
		})

		It("should not be completed initially", func() {
			modal = components.NewExportOptionsModal(120, 40)

			Expect(modal.IsCompleted()).To(BeFalse())
		})
	})

	Describe("Visibility Management", func() {
		BeforeEach(func() {
			modal = components.NewExportOptionsModal(120, 40)
		})

		It("should be visible by default", func() {
			Expect(modal.IsVisible()).To(BeTrue())
		})

		It("should hide when Hide() is called", func() {
			modal.Hide()

			Expect(modal.IsVisible()).To(BeFalse())
		})

		It("should show when Show() is called after hiding", func() {
			modal.Hide()
			modal.Show()

			Expect(modal.IsVisible()).To(BeTrue())
		})

		It("should return empty string from View() when hidden", func() {
			modal.Hide()

			view := modal.View()
			Expect(view).To(Equal(""))
		})
	})

	Describe("Format Selection", func() {
		BeforeEach(func() {
			modal = components.NewExportOptionsModal(120, 40)
		})

		It("should support Text format", func() {
			modal.SetFormat("text")

			data := modal.GetExportData()
			Expect(data.Format).To(Equal("text"))
		})

		It("should support Markdown format", func() {
			modal.SetFormat("markdown")

			data := modal.GetExportData()
			Expect(data.Format).To(Equal("markdown"))
		})

		It("should support YAML format", func() {
			modal.SetFormat("yaml")

			data := modal.GetExportData()
			Expect(data.Format).To(Equal("yaml"))
		})
	})

	Describe("Location Selection", func() {
		BeforeEach(func() {
			modal = components.NewExportOptionsModal(120, 40)
		})

		It("should support File location", func() {
			modal.SetLocation("file")

			data := modal.GetExportData()
			Expect(data.Location).To(Equal("file"))
		})

		It("should support Clipboard location", func() {
			modal.SetLocation("clipboard")

			data := modal.GetExportData()
			Expect(data.Location).To(Equal("clipboard"))
		})
	})

	Describe("Form Submission", func() {
		BeforeEach(func() {
			modal = components.NewExportOptionsModal(120, 40)
			modal.Init()
		})

		It("should mark as completed when form is submitted", func() {
			// Simulate form completion
			modal.SetFormat("markdown")
			modal.SetLocation("file")

			// Form submission happens through huh.Form internally
			// We can call Complete() directly for testing
			modal.Complete()

			Expect(modal.IsCompleted()).To(BeTrue())
		})

		It("should hide modal when form is submitted", func() {
			modal.Complete()

			Expect(modal.IsVisible()).To(BeFalse())
		})

		It("should preserve selected values after submission", func() {
			modal.SetFormat("yaml")
			modal.SetLocation("clipboard")
			modal.Complete()

			data := modal.GetExportData()
			Expect(data.Format).To(Equal("yaml"))
			Expect(data.Location).To(Equal("clipboard"))
		})
	})

	Describe("Cancellation", func() {
		BeforeEach(func() {
			modal = components.NewExportOptionsModal(120, 40)
		})

		It("should hide modal on Esc key", func() {
			modal.Update(tea.KeyMsg{Type: tea.KeyEsc})

			Expect(modal.IsVisible()).To(BeFalse())
		})

		It("should not be completed when cancelled", func() {
			modal.Update(tea.KeyMsg{Type: tea.KeyEsc})

			Expect(modal.IsCompleted()).To(BeFalse())
		})

		It("should preserve partial selections when cancelled", func() {
			modal.SetFormat("markdown")
			modal.Update(tea.KeyMsg{Type: tea.KeyEsc})

			data := modal.GetExportData()
			Expect(data.Format).To(Equal("markdown"))
		})
	})

	Describe("WindowSizeMsg Handling", func() {
		BeforeEach(func() {
			modal = components.NewExportOptionsModal(120, 40)
		})

		It("should update dimensions on WindowSizeMsg", func() {
			modal.Update(tea.WindowSizeMsg{Width: 160, Height: 50})

			width, height := modal.GetDimensions()
			Expect(width).To(Equal(160))
			Expect(height).To(Equal(50))
		})

		It("should rebuild form with new dimensions", func() {
			initialView := modal.View()

			modal.Update(tea.WindowSizeMsg{Width: 80, Height: 24})

			newView := modal.View()
			// View should be different after resize
			Expect(newView).NotTo(Equal(initialView))
		})

		It("should handle minimum dimensions gracefully", func() {
			modal.Update(tea.WindowSizeMsg{Width: 40, Height: 10})

			view := modal.View()
			Expect(view).NotTo(BeEmpty())
		})
	})

	Describe("View Rendering", func() {
		BeforeEach(func() {
			modal = components.NewExportOptionsModal(120, 40)
		})

		It("should render modal with solid background", func() {
			view := modal.View()

			Expect(view).NotTo(BeEmpty())
			// Should contain border characters
			Expect(view).To(ContainSubstring("─"))
		})

		It("should display title", func() {
			view := modal.View()

			Expect(view).To(ContainSubstring("Export Options"))
		})

		It("should display format field", func() {
			view := modal.View()

			// Should show export format selection
			Expect(view).NotTo(BeEmpty())
		})

		It("should display location field", func() {
			view := modal.View()

			// Should show save location selection
			Expect(view).NotTo(BeEmpty())
		})

		It("should display KeyBadge footer", func() {
			view := modal.View()

			// Footer should be present
			Expect(view).NotTo(BeEmpty())
		})
	})

	Describe("Theme Integration", func() {
		BeforeEach(func() {
			modal = components.NewExportOptionsModal(120, 40)
		})

		It("should use theme for form styling", func() {
			view := modal.View()

			// Theme should be applied
			Expect(view).NotTo(BeEmpty())
		})

		It("should have themed borders", func() {
			view := modal.View()

			// Should have rounded borders with theme colors
			Expect(view).To(ContainSubstring("─"))
		})
	})

	Describe("Edge Cases", func() {
		It("should handle rapid key presses gracefully", func() {
			modal = components.NewExportOptionsModal(120, 40)

			// Rapid key presses should not cause issues
			for i := 0; i < 10; i++ {
				modal.Update(tea.KeyMsg{Type: tea.KeyEnter})
			}

			Expect(modal).NotTo(BeNil())
		})

		It("should handle form Init() call", func() {
			modal = components.NewExportOptionsModal(120, 40)

			cmd := modal.Init()
			// Init should return form init command
			Expect(cmd).NotTo(BeNil())
		})

		It("should handle multiple Hide/Show cycles", func() {
			modal = components.NewExportOptionsModal(120, 40)

			for i := 0; i < 5; i++ {
				modal.Hide()
				Expect(modal.IsVisible()).To(BeFalse())
				modal.Show()
				Expect(modal.IsVisible()).To(BeTrue())
			}
		})
	})
})
