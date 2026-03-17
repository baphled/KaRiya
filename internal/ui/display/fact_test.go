package display_test

import (
	"time"

	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"

	"github.com/baphled/kariya/internal/domain/career"
	"github.com/baphled/kariya/internal/testutil/fixtures"
	"github.com/baphled/kariya/internal/ui/display"
)

var _ = Describe("Fact display type", func() {
	Describe("FactFromDomain", func() {
		It("converts a fully populated fact", func() {
			createdAt := time.Date(2024, time.April, 1, 12, 0, 0, 0, time.UTC)
			updatedAt := time.Date(2024, time.April, 2, 12, 0, 0, 0, time.UTC)

			fact := fixtures.Fact("fact-1", "event-1")
			fact.Text = "Improved deployment reliability across services"
			fact.CompetencyCategories = []string{"technical", "leadership"}
			fact.AudienceRelevance = []string{"hiring-manager", "technical-peer"}
			fact.StrengthSignal = "strong"
			fact.SourceBurstID = "burst-1"
			fact.CreatedAt = createdAt
			fact.UpdatedAt = updatedAt

			result := display.FactFromDomain(fact)

			Expect(result).To(Equal(display.Fact{
				ID:                   "fact-1",
				Text:                 "Improved deployment reliability across services",
				CompetencyCategories: []string{"technical", "leadership"},
				RoleFit:              "staff",
				AudienceRelevance:    []string{"hiring-manager", "technical-peer"},
				StrengthSignal:       "strong",
				SourceEventID:        "event-1",
				SourceBurstID:        "burst-1",
				CreatedAt:            createdAt,
				UpdatedAt:            updatedAt,
			}))
		})

		It("handles nil input gracefully", func() {
			Expect(display.FactFromDomain(nil)).To(Equal(display.Fact{}))
		})
	})

	Describe("FactsFromDomain", func() {
		It("converts a slice of facts", func() {
			f1 := fixtures.Fact("fact-1", "event-1")
			f1.Text = "First fact"
			f2 := fixtures.Fact("fact-2", "event-2")
			f2.Text = "Second fact"
			facts := []*career.Fact{f1, f2}

			result := display.FactsFromDomain(facts)

			Expect(result).To(HaveLen(2))
			Expect(result[0].ID).To(Equal("fact-1"))
			Expect(result[0].Text).To(Equal("First fact"))
			Expect(result[1].ID).To(Equal("fact-2"))
			Expect(result[1].Text).To(Equal("Second fact"))
		})

		It("returns nil for nil input", func() {
			Expect(display.FactsFromDomain(nil)).To(BeNil())
		})

		It("returns empty slice for empty input", func() {
			Expect(display.FactsFromDomain([]*career.Fact{})).To(Equal([]display.Fact{}))
		})
	})
})
