package components_test

import (
	"github.com/baphled/kariya/internal/cli/components"
	"github.com/baphled/kariya/internal/cli/themes"
	tea "github.com/charmbracelet/bubbletea"
	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
)

var _ = Describe("ExportErrorModal", func() {
	var modal *components.ExportErrorModal

	BeforeEach(func() {
		modal = components.NewExportErrorModal(
			"export_failed",
			"Failed to write file: permission denied",
		)
	})

	Describe("NewExportErrorModal", func() {
		It("creates a visible modal", func() {
			Expect(modal.IsVisible()).To(BeTrue())
		})

		It("stores error code", func() {
			Expect(modal.GetErrorCode()).To(Equal("export_failed"))
		})

		It("stores error message", func() {
			Expect(modal.GetErrorMessage()).To(Equal("Failed to write file: permission denied"))
		})

		It("is not retried by default", func() {
			Expect(modal.WantsRetry()).To(BeFalse())
		})
	})

	Describe("Update", func() {
		Context("when pressing r", func() {
			It("signals retry and closes modal", func() {
				_, result := modal.Update(tea.KeyMsg{Type: tea.KeyRunes, Runes: []rune{'r'}})
				Expect(result).To(Equal(components.ErrorResultRetry))
				Expect(modal.WantsRetry()).To(BeTrue())
				Expect(modal.IsVisible()).To(BeFalse())
			})
		})

		Context("when pressing R", func() {
			It("signals retry and closes modal", func() {
				_, result := modal.Update(tea.KeyMsg{Type: tea.KeyRunes, Runes: []rune{'R'}})
				Expect(result).To(Equal(components.ErrorResultRetry))
				Expect(modal.WantsRetry()).To(BeTrue())
			})
		})

		Context("when pressing escape", func() {
			It("signals cancel and closes modal", func() {
				_, result := modal.Update(tea.KeyMsg{Type: tea.KeyEsc})
				Expect(result).To(Equal(components.ErrorResultCancel))
				Expect(modal.WantsRetry()).To(BeFalse())
				Expect(modal.IsVisible()).To(BeFalse())
			})
		})

		Context("when pressing q", func() {
			It("signals cancel and closes modal", func() {
				_, result := modal.Update(tea.KeyMsg{Type: tea.KeyRunes, Runes: []rune{'q'}})
				Expect(result).To(Equal(components.ErrorResultCancel))
				Expect(modal.WantsRetry()).To(BeFalse())
			})
		})

		Context("when pressing other keys", func() {
			It("does not close the modal", func() {
				_, result := modal.Update(tea.KeyMsg{Type: tea.KeyRunes, Runes: []rune{'x'}})
				Expect(result).To(Equal(components.ErrorResultNone))
				Expect(modal.IsVisible()).To(BeTrue())
			})
		})

		Context("when modal is not visible", func() {
			It("ignores input", func() {
				modal.Hide()
				_, result := modal.Update(tea.KeyMsg{Type: tea.KeyRunes, Runes: []rune{'r'}})
				Expect(result).To(Equal(components.ErrorResultNone))
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

		It("shows error title", func() {
			view := modal.View()
			Expect(view).To(SatisfyAny(
				ContainSubstring("Export Failed"),
				ContainSubstring("Error"),
			))
		})

		It("shows error message", func() {
			view := modal.View()
			Expect(view).To(ContainSubstring("permission denied"))
		})

		It("shows retry option", func() {
			view := modal.View()
			Expect(view).To(ContainSubstring("Retry"))
		})

		It("shows cancel option", func() {
			view := modal.View()
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

		It("resets retry state when showing", func() {
			// First trigger retry
			modal.Update(tea.KeyMsg{Type: tea.KeyRunes, Runes: []rune{'r'}})
			Expect(modal.WantsRetry()).To(BeTrue())
			// Then show again
			modal.Show()
			Expect(modal.WantsRetry()).To(BeFalse())
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
