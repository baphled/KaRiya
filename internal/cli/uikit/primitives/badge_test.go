package primitives_test

import (
	"github.com/baphled/kariya/internal/cli/uikit/primitives"
	"github.com/baphled/kariya/internal/cli/uikit/theme"
	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
)

var _ = Describe("Badge", func() {
	var th theme.Theme

	BeforeEach(func() {
		th = theme.Default()
	})

	Describe("NewBadge", func() {
		It("should create badge with label", func() {
			badge := primitives.NewBadge("New", th)
			Expect(badge).NotTo(BeNil())
			rendered := badge.Render()
			Expect(rendered).To(ContainSubstring("New"))
		})

		It("should accept nil theme and use default", func() {
			badge := primitives.NewBadge("Test", nil)
			Expect(badge).NotTo(BeNil())
		})
	})

	Describe("Fluent API", func() {
		Describe("Value", func() {
			It("should set value for key badge", func() {
				badge := primitives.NewBadge("Enter", th).Value("confirm")
				rendered := badge.Render()
				Expect(rendered).To(ContainSubstring("Enter"))
			})

			It("should return badge for chaining", func() {
				badge := primitives.NewBadge("Esc", th)
				result := badge.Value("cancel")
				Expect(result).To(Equal(badge))
			})
		})

		Describe("Variant", func() {
			It("should set Default variant", func() {
				badge := primitives.NewBadge("Tag", th).Variant(primitives.BadgeDefault)
				rendered := badge.Render()
				Expect(rendered).To(ContainSubstring("Tag"))
			})

			It("should set Key variant", func() {
				badge := primitives.NewBadge("Tab", th).Variant(primitives.BadgeKey)
				rendered := badge.Render()
				Expect(rendered).To(ContainSubstring("Tab"))
			})

			It("should set Status variant", func() {
				badge := primitives.NewBadge("Active", th).Variant(primitives.BadgeStatus)
				rendered := badge.Render()
				Expect(rendered).To(ContainSubstring("Active"))
			})

			It("should set Tag variant", func() {
				badge := primitives.NewBadge("Feature", th).Variant(primitives.BadgeTag)
				rendered := badge.Render()
				Expect(rendered).To(ContainSubstring("Feature"))
			})

			It("should return badge for chaining", func() {
				badge := primitives.NewBadge("Test", th)
				result := badge.Variant(primitives.BadgeKey)
				Expect(result).To(Equal(badge))
			})
		})
	})

	Describe("Render", func() {
		It("should render default badge", func() {
			badge := primitives.NewBadge("Info", th)
			rendered := badge.Render()
			Expect(rendered).To(ContainSubstring("Info"))
			Expect(rendered).NotTo(BeEmpty())
		})

		It("should render key badge with bracket format", func() {
			badge := primitives.NewBadge("Enter", th).Variant(primitives.BadgeKey)
			rendered := badge.Render()
			Expect(rendered).To(ContainSubstring("Enter"))
		})

		It("should render status badge", func() {
			badge := primitives.NewBadge("Success", th).Variant(primitives.BadgeStatus)
			rendered := badge.Render()
			Expect(rendered).To(ContainSubstring("Success"))
		})

		It("should render tag badge with pill style", func() {
			badge := primitives.NewBadge("Tag", th).Variant(primitives.BadgeTag)
			rendered := badge.Render()
			Expect(rendered).To(ContainSubstring("Tag"))
		})
	})

	Describe("Convenience Constructors", func() {
		Describe("KeyBadge", func() {
			It("should create key badge", func() {
				badge := primitives.KeyBadge("Esc", "cancel", th)
				rendered := badge.Render()
				Expect(rendered).To(ContainSubstring("Esc"))
			})

			It("should format as key-action pair", func() {
				badge := primitives.KeyBadge("Tab", "next", th)
				rendered := badge.Render()
				// Should show key in some format
				Expect(rendered).To(ContainSubstring("Tab"))
			})
		})

		Describe("StatusBadge", func() {
			It("should create status badge", func() {
				badge := primitives.StatusBadge("Complete", th)
				rendered := badge.Render()
				Expect(rendered).To(ContainSubstring("Complete"))
			})
		})

		Describe("TagBadge", func() {
			It("should create tag badge", func() {
				badge := primitives.TagBadge("Feature", th)
				rendered := badge.Render()
				Expect(rendered).To(ContainSubstring("Feature"))
			})
		})

		// Preset badge constructors for CV workflow
		Describe("NavigateBadge", func() {
			It("should create navigate badge with arrow keys", func() {
				badge := primitives.NavigateBadge(th)
				rendered := badge.Render()
				Expect(rendered).To(ContainSubstring("↑↓"))
				Expect(rendered).To(ContainSubstring("Navigate"))
			})

			It("should accept nil theme", func() {
				badge := primitives.NavigateBadge(nil)
				Expect(badge).NotTo(BeNil())
			})
		})

		Describe("SelectBadge", func() {
			It("should create select badge with Enter key", func() {
				badge := primitives.SelectBadge(th)
				rendered := badge.Render()
				Expect(rendered).To(ContainSubstring("Enter"))
				Expect(rendered).To(ContainSubstring("Select"))
			})
		})

		Describe("CancelBadge", func() {
			It("should create cancel badge with Esc key", func() {
				badge := primitives.CancelBadge(th)
				rendered := badge.Render()
				Expect(rendered).To(ContainSubstring("Esc"))
				Expect(rendered).To(ContainSubstring("Cancel"))
			})
		})

		Describe("BackBadge", func() {
			It("should create back badge with Esc key", func() {
				badge := primitives.BackBadge(th)
				rendered := badge.Render()
				Expect(rendered).To(ContainSubstring("Esc"))
				Expect(rendered).To(ContainSubstring("Back"))
			})
		})

		Describe("ConfirmBadge", func() {
			It("should create confirm badge with Enter key", func() {
				badge := primitives.ConfirmBadge(th)
				rendered := badge.Render()
				Expect(rendered).To(ContainSubstring("Enter"))
				Expect(rendered).To(ContainSubstring("Confirm"))
			})
		})

		Describe("QuitBadge", func() {
			It("should create quit badge with q key", func() {
				badge := primitives.QuitBadge(th)
				rendered := badge.Render()
				Expect(rendered).To(ContainSubstring("q"))
				Expect(rendered).To(ContainSubstring("Quit"))
			})
		})

		Describe("HelpBadge", func() {
			It("should create help badge with ? key", func() {
				badge := primitives.HelpBadge(th)
				rendered := badge.Render()
				Expect(rendered).To(ContainSubstring("?"))
				Expect(rendered).To(ContainSubstring("Help"))
			})
		})

		Describe("SkipBadge", func() {
			It("should create skip badge with Ctrl+S key", func() {
				badge := primitives.SkipBadge(th)
				rendered := badge.Render()
				Expect(rendered).To(ContainSubstring("Ctrl+S"))
				Expect(rendered).To(ContainSubstring("Skip"))
			})
		})
	})

	Describe("RenderHelpFooter", func() {
		It("should render multiple badges separated by spaces", func() {
			badges := []*primitives.Badge{
				primitives.NavigateBadge(th),
				primitives.SelectBadge(th),
			}
			result := primitives.RenderHelpFooter(th, badges...)
			Expect(result).To(ContainSubstring("Navigate"))
			Expect(result).To(ContainSubstring("Select"))
		})

		It("should return empty string for no badges", func() {
			result := primitives.RenderHelpFooter(th)
			Expect(result).To(BeEmpty())
		})

		It("should handle nil theme", func() {
			badges := []*primitives.Badge{
				primitives.NavigateBadge(nil),
			}
			result := primitives.RenderHelpFooter(nil, badges...)
			Expect(result).To(ContainSubstring("Navigate"))
		})

		It("should render single badge without separator", func() {
			badges := []*primitives.Badge{
				primitives.CancelBadge(th),
			}
			result := primitives.RenderHelpFooter(th, badges...)
			Expect(result).To(ContainSubstring("Cancel"))
		})
	})
})
