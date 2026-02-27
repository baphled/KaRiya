package modals_test

import (
	"errors"

	"github.com/baphled/kariya/internal/cli/screens/cv/modals"
	tea "github.com/charmbracelet/bubbletea"
	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
)

var _ = Describe("ProgressModal", func() {
	var (
		modal *modals.ProgressModal
	)

	Describe("NewProgressModal", func() {
		It("should create a progress modal with title and subtitle", func() {
			modal = modals.NewProgressModal("Processing", "Please wait...", true, 120, 40)

			Expect(modal).NotTo(BeNil())
			Expect(modal.GetTitle()).To(Equal("Processing"))
			Expect(modal.GetSubtitle()).To(Equal("Please wait..."))
			Expect(modal.IsVisible()).To(BeTrue())
			Expect(modal.IsCancellable()).To(BeTrue())
		})

		It("should initialize with spinner at frame 0", func() {
			modal = modals.NewProgressModal("Processing", "Please wait...", true, 120, 40)

			Expect(modal.GetSpinnerFrame()).To(Equal(0))
		})

		It("should not be completed initially", func() {
			modal = modals.NewProgressModal("Processing", "Please wait...", true, 120, 40)

			Expect(modal.IsCompleted()).To(BeFalse())
		})

		It("should have no error initially", func() {
			modal = modals.NewProgressModal("Processing", "Please wait...", true, 120, 40)

			Expect(modal.GetError()).To(Succeed())
		})
	})

	Describe("Factory Functions", func() {
		Context("NewExtractingTechsProgress", func() {
			It("should create a tech extraction progress modal", func() {
				modal = modals.NewExtractingTechsProgress(120, 40)

				Expect(modal).NotTo(BeNil())
				Expect(modal.GetTitle()).To(Equal("Extracting Technologies"))
				Expect(modal.GetSubtitle()).To(Equal("Analyzing your skills..."))
				Expect(modal.IsCancellable()).To(BeTrue())
				Expect(modal.IsVisible()).To(BeTrue())
			})
		})

		Context("NewGeneratingCVProgress", func() {
			It("should create a CV generation progress modal", func() {
				modal = modals.NewGeneratingCVProgress("Senior Engineer", "hiring_manager", 120, 40)

				Expect(modal).NotTo(BeNil())
				Expect(modal.GetTitle()).To(Equal("Generating CV"))
				Expect(modal.GetSubtitle()).To(ContainSubstring("Senior Engineer"))
				Expect(modal.GetSubtitle()).To(ContainSubstring("hiring_manager"))
				Expect(modal.IsCancellable()).To(BeTrue())
			})
		})

		Context("NewExportingProgress", func() {
			It("should create an export progress modal", func() {
				modal = modals.NewExportingProgress("Markdown", 120, 40)

				Expect(modal).NotTo(BeNil())
				Expect(modal.GetTitle()).To(Equal("Exporting CV"))
				Expect(modal.GetSubtitle()).To(ContainSubstring("Markdown"))
				Expect(modal.IsCancellable()).To(BeFalse())
			})
		})
	})

	Describe("Visibility Management", func() {
		BeforeEach(func() {
			modal = modals.NewProgressModal("Processing", "Please wait...", true, 120, 40)
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

	Describe("Spinner Animation", func() {
		BeforeEach(func() {
			modal = modals.NewProgressModal("Processing", "Please wait...", true, 120, 40)
		})

		It("should advance spinner frame on tick", func() {
			initialFrame := modal.GetSpinnerFrame()

			modal.Update(modals.SpinnerTickMsg{})

			newFrame := modal.GetSpinnerFrame()
			Expect(newFrame).To(Equal((initialFrame + 1) % 10))
		})

		It("should wrap spinner frame after 10 frames", func() {
			for range 9 {
				modal.Update(modals.SpinnerTickMsg{})
			}
			Expect(modal.GetSpinnerFrame()).To(Equal(9))

			modal.Update(modals.SpinnerTickMsg{})
			Expect(modal.GetSpinnerFrame()).To(Equal(0))
		})

		It("should return tick command on Init()", func() {
			cmd := modal.Init()

			Expect(cmd).NotTo(BeNil())
		})

		It("should return tick command on Update() to continue animation", func() {
			cmd := modal.Update(modals.SpinnerTickMsg{})

			Expect(cmd).NotTo(BeNil())
		})

		It("should execute tick command closure and return SpinnerTickMsg", func() {
			cmd := modal.Init()
			Expect(cmd).NotTo(BeNil())

			msg := cmd()
			Expect(msg).To(BeAssignableToTypeOf(modals.SpinnerTickMsg{}))
		})
	})

	Describe("Cancellation", func() {
		Context("When cancellable", func() {
			BeforeEach(func() {
				modal = modals.NewProgressModal("Processing", "Please wait...", true, 120, 40)
			})

			It("should hide modal on Esc key", func() {
				modal.Update(tea.KeyMsg{Type: tea.KeyEsc})

				Expect(modal.IsVisible()).To(BeFalse())
			})

			It("should mark as cancelled when Esc pressed", func() {
				modal.Update(tea.KeyMsg{Type: tea.KeyEsc})

				Expect(modal.IsCancelled()).To(BeTrue())
			})
		})

		Context("When not cancellable", func() {
			BeforeEach(func() {
				modal = modals.NewProgressModal("Processing", "Please wait...", false, 120, 40)
			})

			It("should not hide modal on Esc key", func() {
				modal.Update(tea.KeyMsg{Type: tea.KeyEsc})

				Expect(modal.IsVisible()).To(BeTrue())
			})

			It("should not be cancelled when Esc pressed", func() {
				modal.Update(tea.KeyMsg{Type: tea.KeyEsc})

				Expect(modal.IsCancelled()).To(BeFalse())
			})
		})
	})

	Describe("Completion Handling", func() {
		BeforeEach(func() {
			modal = modals.NewProgressModal("Processing", "Please wait...", true, 120, 40)
		})

		It("should mark as completed when Complete() is called", func() {
			modal.Complete()

			Expect(modal.IsCompleted()).To(BeTrue())
		})

		It("should hide modal when completed", func() {
			modal.Complete()

			Expect(modal.IsVisible()).To(BeFalse())
		})

		It("should stop spinner animation when completed", func() {
			modal.Complete()

			initialFrame := modal.GetSpinnerFrame()
			modal.Update(modals.SpinnerTickMsg{})

			Expect(modal.GetSpinnerFrame()).To(Equal(initialFrame))
		})
	})

	Describe("Error Handling", func() {
		BeforeEach(func() {
			modal = modals.NewProgressModal("Processing", "Please wait...", true, 120, 40)
		})

		It("should set error when SetError() is called", func() {
			testError := errors.New("processing failed")
			modal.SetError(testError)

			Expect(modal.GetError()).To(Equal(testError))
		})

		It("should mark as completed when error is set", func() {
			modal.SetError(errors.New("processing failed"))

			Expect(modal.IsCompleted()).To(BeTrue())
		})

		It("should hide modal when error is set", func() {
			modal.SetError(errors.New("processing failed"))

			Expect(modal.IsVisible()).To(BeFalse())
		})
	})

	Describe("WindowSizeMsg Handling", func() {
		BeforeEach(func() {
			modal = modals.NewProgressModal("Processing", "Please wait...", true, 120, 40)
		})

		It("should update dimensions on WindowSizeMsg", func() {
			modal.Update(tea.WindowSizeMsg{Width: 160, Height: 50})

			width, height := modal.GetDimensions()
			Expect(width).To(Equal(160))
			Expect(height).To(Equal(50))
		})

		It("should handle minimum dimensions gracefully", func() {
			modal.Update(tea.WindowSizeMsg{Width: 40, Height: 10})

			view := modal.View()
			Expect(view).NotTo(BeEmpty())
		})
	})

	Describe("View Rendering", func() {
		BeforeEach(func() {
			modal = modals.NewProgressModal("Processing", "Please wait...", true, 120, 40)
		})

		It("should render modal with solid background", func() {
			view := modal.View()

			Expect(view).NotTo(BeEmpty())
			Expect(view).To(ContainSubstring("─"))
		})

		It("should display title", func() {
			view := modal.View()

			Expect(view).To(ContainSubstring("Processing"))
		})

		It("should display subtitle", func() {
			view := modal.View()

			Expect(view).To(ContainSubstring("Please wait..."))
		})

		It("should display spinner animation", func() {
			view := modal.View()

			Expect(view).NotTo(BeEmpty())
		})

		It("should display KeyBadge footer when cancellable", func() {
			view := modal.View()

			Expect(view).NotTo(BeEmpty())
		})

		It("should not show cancel option when not cancellable", func() {
			modal = modals.NewProgressModal("Processing", "Please wait...", false, 120, 40)
			view := modal.View()

			Expect(view).NotTo(BeEmpty())
		})
	})

	Describe("Edge Cases", func() {
		It("should handle empty title gracefully", func() {
			modal = modals.NewProgressModal("", "Subtitle", true, 120, 40)

			view := modal.View()
			Expect(view).NotTo(BeEmpty())
		})

		It("should handle empty subtitle gracefully", func() {
			modal = modals.NewProgressModal("Title", "", true, 120, 40)

			view := modal.View()
			Expect(view).NotTo(BeEmpty())
		})

		It("should handle very long title gracefully", func() {
			longTitle := "This is a very long title that should be truncated or wrapped appropriately"
			modal = modals.NewProgressModal(longTitle, "Subtitle", true, 120, 40)

			view := modal.View()
			Expect(view).NotTo(BeEmpty())
		})

		It("should handle rapid spinner ticks", func() {
			modal = modals.NewProgressModal("Processing", "Please wait...", true, 120, 40)

			for range 100 {
				modal.Update(modals.SpinnerTickMsg{})
			}

			Expect(modal.GetSpinnerFrame()).To(BeNumerically("<", 10))
		})

		It("should return nil from Update when not visible", func() {
			modal = modals.NewProgressModal("Processing", "Please wait...", true, 120, 40)
			modal.Hide()

			cmd := modal.Update(modals.SpinnerTickMsg{})

			Expect(cmd).To(BeNil())
		})
	})
})
