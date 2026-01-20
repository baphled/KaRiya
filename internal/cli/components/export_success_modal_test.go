package components_test

import (
	"testing"

	"github.com/baphled/kariya/internal/cli/components"
	"github.com/baphled/kariya/internal/cli/themes"
	tea "github.com/charmbracelet/bubbletea"
	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
)

func TestExportSuccessModal(t *testing.T) {
	RegisterFailHandler(Fail)
	RunSpecs(t, "Export Success Modal Suite")
}

var _ = Describe("ExportSuccessModal", func() {
	var modal *components.ExportSuccessModal

	BeforeEach(func() {
		modal = components.NewExportSuccessModal(
			"Career Events",
			"JSON",
			"File",
			"/home/user/.kariya/exports/events_20240120.json",
			1024, // size in bytes
		)
	})

	Describe("NewExportSuccessModal", func() {
		It("creates a visible modal", func() {
			Expect(modal.IsVisible()).To(BeTrue())
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

		It("stores file path", func() {
			Expect(modal.GetFilePath()).To(Equal("/home/user/.kariya/exports/events_20240120.json"))
		})

		It("stores size", func() {
			Expect(modal.GetSize()).To(Equal(int64(1024)))
		})
	})

	Describe("Update", func() {
		Context("when pressing enter", func() {
			It("closes the modal and signals done", func() {
				_, done := modal.Update(tea.KeyMsg{Type: tea.KeyEnter})
				Expect(done).To(BeTrue())
				Expect(modal.IsVisible()).To(BeFalse())
			})
		})

		Context("when pressing escape", func() {
			It("closes the modal and signals done", func() {
				_, done := modal.Update(tea.KeyMsg{Type: tea.KeyEsc})
				Expect(done).To(BeTrue())
				Expect(modal.IsVisible()).To(BeFalse())
			})
		})

		Context("when pressing other keys", func() {
			It("does not close the modal", func() {
				_, done := modal.Update(tea.KeyMsg{Type: tea.KeyRunes, Runes: []rune{'x'}})
				Expect(done).To(BeFalse())
				Expect(modal.IsVisible()).To(BeTrue())
			})
		})

		Context("when modal is not visible", func() {
			It("ignores input", func() {
				modal.Hide()
				_, done := modal.Update(tea.KeyMsg{Type: tea.KeyEnter})
				Expect(done).To(BeFalse())
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

		It("shows success title", func() {
			view := modal.View()
			Expect(view).To(ContainSubstring("Export Complete"))
		})

		It("shows checkmark or success indicator", func() {
			view := modal.View()
			Expect(view).To(SatisfyAny(
				ContainSubstring("✓"),
				ContainSubstring("Success"),
				ContainSubstring("Complete"),
			))
		})

		It("shows artifact type", func() {
			view := modal.View()
			Expect(view).To(ContainSubstring("Career Events"))
		})

		It("shows file path for file destination", func() {
			view := modal.View()
			Expect(view).To(ContainSubstring("events_20240120.json"))
		})

		It("shows done instructions", func() {
			view := modal.View()
			Expect(view).To(SatisfyAny(
				ContainSubstring("Done"),
				ContainSubstring("Continue"),
				ContainSubstring("Enter"),
			))
		})
	})

	Describe("View for clipboard destination", func() {
		BeforeEach(func() {
			modal = components.NewExportSuccessModal(
				"Career Events",
				"JSON",
				"Clipboard",
				"clipboard",
				512,
			)
		})

		It("shows clipboard message", func() {
			view := modal.View()
			Expect(view).To(SatisfyAny(
				ContainSubstring("Clipboard"),
				ContainSubstring("clipboard"),
				ContainSubstring("copied"),
			))
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
			view := modal.View()
			Expect(view).NotTo(BeEmpty())
		})
	})

	Describe("Init", func() {
		It("returns nil", func() {
			cmd := modal.Init()
			Expect(cmd).To(BeNil())
		})
	})
})
