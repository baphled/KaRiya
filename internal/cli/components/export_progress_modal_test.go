package components_test

import (
	"fmt"
	"strings"
	"testing"

	"github.com/baphled/kariya/internal/cli/components"
	"github.com/baphled/kariya/internal/cli/themes"
	tea "github.com/charmbracelet/bubbletea"
	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
)

func TestExportProgressModal(t *testing.T) {
	RegisterFailHandler(Fail)
	RunSpecs(t, "Export Progress Modal Suite")
}

var _ = Describe("ExportProgressModal", func() {
	var modal *components.ExportProgressModal

	BeforeEach(func() {
		modal = components.NewExportProgressModal(
			"Career Events",
			"JSON",
			"File",
			100, // width
			40,  // height
		)
	})

	Describe("NewExportProgressModal", func() {
		It("creates a visible modal", func() {
			Expect(modal.IsVisible()).To(BeTrue())
		})

		It("is not completed by default", func() {
			Expect(modal.IsCompleted()).To(BeFalse())
		})

		It("is not cancelled by default", func() {
			Expect(modal.IsCancelled()).To(BeFalse())
		})

		It("stores artifact type", func() {
			Expect(modal.GetArtifactType()).To(Equal("Career Events"))
		})

		It("stores format", func() {
			Expect(modal.GetFormat()).To(Equal("JSON"))
		})

		It("stores destination", func() {
			Expect(modal.GetDestination()).To(Equal("File"))
		})
	})

	Describe("Init", func() {
		It("returns a tick command for spinner animation", func() {
			cmd := modal.Init()
			Expect(cmd).NotTo(BeNil())
		})
	})

	Describe("Update", func() {
		Context("when receiving SpinnerTickMsg", func() {
			It("advances the spinner", func() {
				initialFrame := modal.GetSpinnerFrame()
				modal.Update(components.SpinnerTickMsg{})
				newFrame := modal.GetSpinnerFrame()
				Expect(newFrame).To(Equal((initialFrame + 1) % 10))
			})

			It("returns another tick command", func() {
				cmd := modal.Update(components.SpinnerTickMsg{})
				Expect(cmd).NotTo(BeNil())
			})

			It("does not tick after completion", func() {
				modal.Complete()
				initialFrame := modal.GetSpinnerFrame()
				modal.Update(components.SpinnerTickMsg{})
				Expect(modal.GetSpinnerFrame()).To(Equal(initialFrame))
			})
		})

		Context("when pressing escape", func() {
			It("cancels the operation", func() {
				modal.Update(tea.KeyMsg{Type: tea.KeyEsc})
				Expect(modal.IsCancelled()).To(BeTrue())
				Expect(modal.IsVisible()).To(BeFalse())
			})
		})

		Context("when pressing other keys", func() {
			It("does not cancel", func() {
				modal.Update(tea.KeyMsg{Type: tea.KeyRunes, Runes: []rune{'x'}})
				Expect(modal.IsCancelled()).To(BeFalse())
				Expect(modal.IsVisible()).To(BeTrue())
			})
		})

		Context("when modal is not visible", func() {
			It("ignores input", func() {
				modal.Hide()
				modal.Update(tea.KeyMsg{Type: tea.KeyEsc})
				Expect(modal.IsCancelled()).To(BeFalse())
			})
		})

		Context("when receiving window size message", func() {
			It("updates dimensions", func() {
				modal.Update(tea.WindowSizeMsg{Width: 200, Height: 50})
				w, h := modal.GetDimensions()
				Expect(w).To(Equal(200))
				Expect(h).To(Equal(50))
			})
		})
	})

	Describe("View", func() {
		It("returns empty string when not visible", func() {
			modal.Hide()
			Expect(modal.View()).To(BeEmpty())
		})

		It("shows exporting title", func() {
			view := modal.View()
			Expect(view).To(ContainSubstring("Exporting"))
		})

		It("shows artifact type", func() {
			view := modal.View()
			Expect(view).To(ContainSubstring("Career Events"))
		})

		It("shows format", func() {
			view := modal.View()
			Expect(view).To(ContainSubstring("JSON"))
		})

		It("shows destination", func() {
			view := modal.View()
			Expect(view).To(ContainSubstring("File"))
		})

		It("shows spinner character", func() {
			view := modal.View()
			// Check for one of the spinner frames
			hasSpinner := false
			for _, frame := range []string{"⠋", "⠙", "⠹", "⠸", "⠼", "⠴", "⠦", "⠧", "⠇", "⠏"} {
				if strings.Contains(view, frame) {
					hasSpinner = true
					break
				}
			}
			Expect(hasSpinner).To(BeTrue())
		})

		It("shows cancel instructions", func() {
			view := modal.View()
			Expect(view).To(ContainSubstring("Cancel"))
		})
	})

	Describe("Complete", func() {
		It("marks operation as completed", func() {
			modal.Complete()
			Expect(modal.IsCompleted()).To(BeTrue())
		})

		It("hides the modal", func() {
			modal.Complete()
			Expect(modal.IsVisible()).To(BeFalse())
		})
	})

	Describe("SetError", func() {
		It("marks operation as completed", func() {
			modal.SetError(fmt.Errorf("test error"))
			Expect(modal.IsCompleted()).To(BeTrue())
		})

		It("hides the modal", func() {
			modal.SetError(fmt.Errorf("test error"))
			Expect(modal.IsVisible()).To(BeFalse())
		})

		It("stores the error", func() {
			err := fmt.Errorf("test error")
			modal.SetError(err)
			Expect(modal.GetError()).To(Equal(err))
		})
	})

	Describe("Show/Hide", func() {
		It("can show the modal", func() {
			modal.Hide()
			modal.Show()
			Expect(modal.IsVisible()).To(BeTrue())
		})

		It("can hide the modal", func() {
			modal.Hide()
			Expect(modal.IsVisible()).To(BeFalse())
		})
	})

	Describe("SetTheme", func() {
		It("accepts a theme", func() {
			theme := themes.NewDefaultTheme()
			modal.SetTheme(theme)
			// No panic means success - theme is used in View()
			view := modal.View()
			Expect(view).NotTo(BeEmpty())
		})
	})
})
