package components_test

import (
	"github.com/baphled/kariya/internal/cli/components"
	"github.com/baphled/kariya/internal/cli/themes"
	tea "github.com/charmbracelet/bubbletea"
	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
)

var _ = Describe("ExportConfirmModal", func() {
	var modal *components.ExportConfirmModal

	BeforeEach(func() {
		modal = components.NewExportConfirmModal(
			"Career Events",
			"JSON",
			"File",
		)
	})

	Describe("NewExportConfirmModal", func() {
		It("creates a visible modal", func() {
			Expect(modal.IsVisible()).To(BeTrue())
		})

		It("is not confirmed by default", func() {
			Expect(modal.WasConfirmed()).To(BeFalse())
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

	Describe("Update", func() {
		Context("when pressing y", func() {
			It("confirms the export", func() {
				_, confirmed := modal.Update(tea.KeyMsg{Type: tea.KeyRunes, Runes: []rune{'y'}})
				Expect(confirmed).To(BeTrue())
				Expect(modal.WasConfirmed()).To(BeTrue())
				Expect(modal.IsVisible()).To(BeFalse())
			})
		})

		Context("when pressing Y", func() {
			It("confirms the export", func() {
				_, confirmed := modal.Update(tea.KeyMsg{Type: tea.KeyRunes, Runes: []rune{'Y'}})
				Expect(confirmed).To(BeTrue())
				Expect(modal.WasConfirmed()).To(BeTrue())
			})
		})

		Context("when pressing enter", func() {
			It("confirms the export", func() {
				_, confirmed := modal.Update(tea.KeyMsg{Type: tea.KeyEnter})
				Expect(confirmed).To(BeTrue())
				Expect(modal.WasConfirmed()).To(BeTrue())
			})
		})

		Context("when pressing n", func() {
			It("cancels the export", func() {
				_, confirmed := modal.Update(tea.KeyMsg{Type: tea.KeyRunes, Runes: []rune{'n'}})
				Expect(confirmed).To(BeFalse())
				Expect(modal.WasConfirmed()).To(BeFalse())
				Expect(modal.IsVisible()).To(BeFalse())
			})
		})

		Context("when pressing N", func() {
			It("cancels the export", func() {
				_, confirmed := modal.Update(tea.KeyMsg{Type: tea.KeyRunes, Runes: []rune{'N'}})
				Expect(confirmed).To(BeFalse())
				Expect(modal.WasConfirmed()).To(BeFalse())
			})
		})

		Context("when pressing escape", func() {
			It("cancels the export", func() {
				_, confirmed := modal.Update(tea.KeyMsg{Type: tea.KeyEsc})
				Expect(confirmed).To(BeFalse())
				Expect(modal.WasConfirmed()).To(BeFalse())
				Expect(modal.IsVisible()).To(BeFalse())
			})
		})

		Context("when pressing other keys", func() {
			It("does not confirm or cancel", func() {
				_, confirmed := modal.Update(tea.KeyMsg{Type: tea.KeyRunes, Runes: []rune{'x'}})
				Expect(confirmed).To(BeFalse())
				Expect(modal.IsVisible()).To(BeTrue())
			})
		})

		Context("when modal is not visible", func() {
			It("ignores input", func() {
				modal.Hide()
				_, confirmed := modal.Update(tea.KeyMsg{Type: tea.KeyRunes, Runes: []rune{'y'}})
				Expect(confirmed).To(BeFalse())
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

		It("shows export confirmation title", func() {
			view := modal.View()
			Expect(view).To(ContainSubstring("Confirm Export"))
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

		It("shows confirmation instructions", func() {
			view := modal.View()
			Expect(view).To(ContainSubstring("Confirm"))
			Expect(view).To(ContainSubstring("Cancel"))
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

		It("resets confirmed state when showing", func() {
			// First confirm
			modal.Update(tea.KeyMsg{Type: tea.KeyEnter})
			Expect(modal.WasConfirmed()).To(BeTrue())
			// Then show again
			modal.Show()
			Expect(modal.WasConfirmed()).To(BeFalse())
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

	Describe("SetDimensions", func() {
		It("updates dimensions", func() {
			modal.SetDimensions(150, 60)
			w, h := modal.GetDimensions()
			Expect(w).To(Equal(150))
			Expect(h).To(Equal(60))
		})
	})

	Describe("Init", func() {
		It("returns nil", func() {
			cmd := modal.Init()
			Expect(cmd).To(BeNil())
		})
	})
})
