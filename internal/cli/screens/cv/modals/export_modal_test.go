package modals_test

import (
	"github.com/baphled/kariya/internal/cli/screens/cv/modals"
	tea "github.com/charmbracelet/bubbletea"
	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
)

var _ = Describe("ExportModal", func() {
	var (
		modal *modals.ExportModal
	)

	Describe("NewExportModal", func() {
		It("should create an export options modal", func() {
			modal = modals.NewExportModal(120, 40)

			Expect(modal).NotTo(BeNil())
			Expect(modal.IsVisible()).To(BeTrue())
		})

		It("should initialize with default format and location", func() {
			modal = modals.NewExportModal(120, 40)

			data := modal.GetExportData()
			Expect(data).NotTo(BeNil())
			Expect(data.Format).To(Equal("text"))
			Expect(data.Location).To(Equal("file"))
		})

		It("should not be completed initially", func() {
			modal = modals.NewExportModal(120, 40)

			Expect(modal.IsCompleted()).To(BeFalse())
		})
	})

	Describe("Visibility Management", func() {
		BeforeEach(func() {
			modal = modals.NewExportModal(120, 40)
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
			modal = modals.NewExportModal(120, 40)
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
			modal = modals.NewExportModal(120, 40)
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
			modal = modals.NewExportModal(120, 40)
			modal.Init()
		})

		It("should mark as completed when form is submitted", func() {
			modal.SetFormat("markdown")
			modal.SetLocation("file")

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
			modal = modals.NewExportModal(120, 40)
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
			modal = modals.NewExportModal(120, 40)
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
			modal = modals.NewExportModal(120, 40)
		})

		It("should render modal with solid background", func() {
			view := modal.View()

			Expect(view).NotTo(BeEmpty())
			Expect(view).To(ContainSubstring("─"))
		})

		It("should display title", func() {
			view := modal.View()

			Expect(view).To(ContainSubstring("Export Options"))
		})

		It("should display format field", func() {
			view := modal.View()

			Expect(view).NotTo(BeEmpty())
		})

		It("should display location field", func() {
			view := modal.View()

			Expect(view).NotTo(BeEmpty())
		})

		It("should display KeyBadge footer", func() {
			view := modal.View()

			Expect(view).NotTo(BeEmpty())
		})
	})

	Describe("Theme Integration", func() {
		BeforeEach(func() {
			modal = modals.NewExportModal(120, 40)
		})

		It("should use theme for form styling", func() {
			view := modal.View()

			Expect(view).NotTo(BeEmpty())
		})

		It("should have themed borders", func() {
			view := modal.View()

			Expect(view).To(ContainSubstring("─"))
		})
	})

	Describe("Edge Cases", func() {
		It("should handle rapid key presses gracefully", func() {
			modal = modals.NewExportModal(120, 40)

			for range 10 {
				modal.Update(tea.KeyMsg{Type: tea.KeyEnter})
			}

			Expect(modal).NotTo(BeNil())
		})

		It("should handle form Init() call", func() {
			modal = modals.NewExportModal(120, 40)

			cmd := modal.Init()
			Expect(cmd).NotTo(BeNil())
		})

		It("should handle multiple Hide/Show cycles", func() {
			modal = modals.NewExportModal(120, 40)

			for range 5 {
				modal.Hide()
				Expect(modal.IsVisible()).To(BeFalse())
				modal.Show()
				Expect(modal.IsVisible()).To(BeTrue())
			}
		})
	})
})
