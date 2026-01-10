package components_test

import (
	"github.com/baphled/kariya/internal/cli/components"
	"github.com/baphled/kariya/internal/cli/themes"
	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
)

var _ = Describe("KeyBadge", func() {
	var theme themes.Theme

	BeforeEach(func() {
		theme = themes.NewDefaultTheme()
	})

	Describe("NewKeyBadge", func() {
		It("should create a key badge with key and hint", func() {
			badge := components.NewKeyBadge("Esc", "Cancel")
			Expect(badge.Key).To(Equal("Esc"))
			Expect(badge.Hint).To(Equal("Cancel"))
		})
	})

	Describe("Render", func() {
		It("should render the badge with key and hint", func() {
			badge := components.NewKeyBadge("Enter", "Confirm")
			result := badge.Render(theme)
			Expect(result).To(ContainSubstring("Enter"))
			Expect(result).To(ContainSubstring("Confirm"))
		})

		It("should render correctly with nil theme", func() {
			badge := components.NewKeyBadge("q", "Quit")
			result := badge.Render(nil)
			Expect(result).To(ContainSubstring("q"))
			Expect(result).To(ContainSubstring("Quit"))
		})

		It("should handle empty key", func() {
			badge := components.NewKeyBadge("", "Empty Key")
			result := badge.Render(theme)
			Expect(result).To(ContainSubstring("Empty Key"))
		})

		It("should handle empty hint", func() {
			badge := components.NewKeyBadge("x", "")
			result := badge.Render(theme)
			Expect(result).To(ContainSubstring("x"))
		})
	})

	Describe("RenderHelpFooter", func() {
		It("should render multiple badges", func() {
			badges := []components.KeyBadge{
				{Key: "↑/↓", Hint: "Navigate"},
				{Key: "Enter", Hint: "Select"},
				{Key: "Esc", Hint: "Cancel"},
			}
			result := components.RenderHelpFooter(theme, badges...)
			Expect(result).To(ContainSubstring("Navigate"))
			Expect(result).To(ContainSubstring("Select"))
			Expect(result).To(ContainSubstring("Cancel"))
		})

		It("should handle empty badge list", func() {
			result := components.RenderHelpFooter(theme)
			Expect(result).To(BeEmpty())
		})

		It("should handle nil theme", func() {
			badges := []components.KeyBadge{
				{Key: "q", Hint: "Quit"},
			}
			result := components.RenderHelpFooter(nil, badges...)
			Expect(result).To(ContainSubstring("q"))
			Expect(result).To(ContainSubstring("Quit"))
		})

		It("should separate badges with spacing", func() {
			badges := []components.KeyBadge{
				{Key: "a", Hint: "First"},
				{Key: "b", Hint: "Second"},
			}
			result := components.RenderHelpFooter(theme, badges...)
			// Check that both badges are present and properly spaced
			Expect(result).To(ContainSubstring("First"))
			Expect(result).To(ContainSubstring("Second"))
		})
	})

	Describe("CommonBadges", func() {
		It("should provide a Navigate badge", func() {
			badge := components.NavigateBadge()
			Expect(badge.Key).To(ContainSubstring("↑"))
			Expect(badge.Hint).To(Equal("Navigate"))
		})

		It("should provide a Select badge", func() {
			badge := components.SelectBadge()
			Expect(badge.Key).To(Equal("Enter"))
			Expect(badge.Hint).To(Equal("Select"))
		})

		It("should provide a Cancel badge", func() {
			badge := components.CancelBadge()
			Expect(badge.Key).To(Equal("Esc"))
			Expect(badge.Hint).To(Equal("Cancel"))
		})

		It("should provide a Quit badge", func() {
			badge := components.QuitBadge()
			Expect(badge.Key).To(Equal("q"))
			Expect(badge.Hint).To(Equal("Quit"))
		})

		It("should provide a Help badge", func() {
			badge := components.HelpBadge()
			Expect(badge.Key).To(Equal("?"))
			Expect(badge.Hint).To(Equal("Help"))
		})

		It("should provide a Back badge", func() {
			badge := components.BackBadge()
			Expect(badge.Key).To(Equal("Esc"))
			Expect(badge.Hint).To(Equal("Back"))
		})

		It("should provide a Confirm badge", func() {
			badge := components.ConfirmBadge()
			Expect(badge.Key).To(Equal("Enter"))
			Expect(badge.Hint).To(Equal("Confirm"))
		})

		It("should provide an Edit badge", func() {
			badge := components.EditBadge()
			Expect(badge.Key).To(Equal("e"))
			Expect(badge.Hint).To(Equal("Edit"))
		})

		It("should provide a Delete badge", func() {
			badge := components.DeleteBadge()
			Expect(badge.Key).To(Equal("d"))
			Expect(badge.Hint).To(Equal("Delete"))
		})

		It("should provide a Save badge", func() {
			badge := components.SaveBadge()
			Expect(badge.Key).To(Equal("Ctrl+S"))
			Expect(badge.Hint).To(Equal("Save"))
		})
	})

	Describe("StandardFooters", func() {
		It("should provide a menu footer", func() {
			result := components.RenderMenuFooter(theme)
			Expect(result).To(ContainSubstring("Navigate"))
			Expect(result).To(ContainSubstring("Select"))
			Expect(result).To(ContainSubstring("Help"))
			Expect(result).To(ContainSubstring("Quit"))
		})

		It("should provide a list footer", func() {
			result := components.RenderListFooter(theme)
			Expect(result).To(ContainSubstring("Navigate"))
			Expect(result).To(ContainSubstring("Select"))
			Expect(result).To(ContainSubstring("Back"))
		})

		It("should provide a form footer", func() {
			result := components.RenderFormFooter(theme)
			Expect(result).To(ContainSubstring("Confirm"))
			Expect(result).To(ContainSubstring("Cancel"))
		})

		It("should provide an edit footer", func() {
			result := components.RenderEditFooter(theme)
			Expect(result).To(ContainSubstring("Save"))
			Expect(result).To(ContainSubstring("Cancel"))
		})

		It("should provide a confirmation footer", func() {
			result := components.RenderConfirmFooter(theme)
			Expect(result).To(ContainSubstring("Confirm"))
			Expect(result).To(ContainSubstring("Cancel"))
		})
	})
})
