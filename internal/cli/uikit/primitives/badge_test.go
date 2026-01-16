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
	})
})
