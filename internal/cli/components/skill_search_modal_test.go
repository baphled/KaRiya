package components_test

import (
	tea "github.com/charmbracelet/bubbletea"
	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"

	"github.com/baphled/kariya/internal/cli/components"
)

var _ = Describe("SkillSearchModal", func() {
	var (
		modal  *components.SkillSearchModal
		width  int
		height int
	)

	BeforeEach(func() {
		width = 80
		height = 24
		modal = components.NewSkillSearchModal("", width, height)
	})

	Describe("NewSkillSearchModal", func() {
		It("should create a new search modal", func() {
			Expect(modal).NotTo(BeNil())
		})

		It("should initialize as not visible", func() {
			Expect(modal.IsVisible()).To(BeFalse())
		})

		It("should accept pre-populated search text", func() {
			modal = components.NewSkillSearchModal("test search", width, height)
			Expect(modal.GetSearchText()).To(Equal("test search"))
		})
	})

	Describe("Init", func() {
		It("should make modal visible", func() {
			modal.Init()
			Expect(modal.IsVisible()).To(BeTrue())
		})

		It("should return a command", func() {
			cmd := modal.Init()
			Expect(cmd).NotTo(BeNil())
		})
	})

	Describe("Update", func() {
		BeforeEach(func() {
			modal.Init()
		})

		Context("when Esc is pressed", func() {
			It("should hide modal and return applied=false", func() {
				cmd, applied, data := modal.Update(tea.KeyMsg{Type: tea.KeyEsc})
				Expect(cmd).To(BeNil())
				Expect(applied).To(BeFalse())
				Expect(data).To(BeNil())
				Expect(modal.IsVisible()).To(BeFalse())
			})
		})

		It("should forward other keys to form", func() {
			cmd, applied, data := modal.Update(tea.KeyMsg{Type: tea.KeyRunes, Runes: []rune{'a'}})
			// Form should handle the key
			Expect(applied).To(BeFalse())
			Expect(data).To(BeNil())
			// Modal should still be visible
			Expect(modal.IsVisible()).To(BeTrue())
			// Cmd may or may not be nil depending on form state
			_ = cmd
		})
	})

	Describe("View", func() {
		Context("when not visible", func() {
			It("should return empty string", func() {
				Expect(modal.View()).To(BeEmpty())
			})
		})

		Context("when visible", func() {
			BeforeEach(func() {
				modal.Init()
			})

			It("should return non-empty view", func() {
				view := modal.View()
				Expect(view).NotTo(BeEmpty())
			})

			It("should contain search field", func() {
				view := modal.View()
				Expect(view).To(ContainSubstring("Search"))
			})
		})
	})

	Describe("Visibility Methods", func() {
		It("should toggle visibility with Show/Hide", func() {
			Expect(modal.IsVisible()).To(BeFalse())

			modal.Show()
			Expect(modal.IsVisible()).To(BeTrue())

			modal.Hide()
			Expect(modal.IsVisible()).To(BeFalse())
		})
	})

	Describe("SetSize", func() {
		It("should update modal dimensions", func() {
			newWidth := 120
			newHeight := 40
			modal.SetSize(newWidth, newHeight)
			// Modal should still be functional after resize
			modal.Init()
			Expect(modal.IsVisible()).To(BeTrue())
		})
	})

	Describe("RenderOverlay", func() {
		var baseView string

		BeforeEach(func() {
			baseView = "Base View Content\nLine 2\nLine 3"
		})

		Context("when not visible", func() {
			It("should return base view unchanged", func() {
				result := modal.RenderOverlay(baseView)
				Expect(result).To(Equal(baseView))
			})
		})

		Context("when visible", func() {
			BeforeEach(func() {
				modal.Init()
			})

			It("should return overlay view", func() {
				result := modal.RenderOverlay(baseView)
				Expect(result).NotTo(Equal(baseView))
				// Should contain overlay rendering
				Expect(result).NotTo(BeEmpty())
			})
		})
	})

	Describe("GetSearchText", func() {
		It("should return current search text", func() {
			modal = components.NewSkillSearchModal("initial", width, height)
			Expect(modal.GetSearchText()).To(Equal("initial"))
		})

		It("should return empty string when not set", func() {
			Expect(modal.GetSearchText()).To(BeEmpty())
		})
	})
})
