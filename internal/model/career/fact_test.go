package career

import (
	"time"

	domain "github.com/baphled/kariya/internal/domain/career"
	"github.com/baphled/kariya/internal/testutil/fixtures"
	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
)

var _ = Describe("Fact Model", func() {
	Describe("TableName", func() {
		It("returns facts", func() {
			Expect(Fact{}.TableName()).To(Equal("facts"))
		})
	})

	Describe("ToDomain", func() {
		It("converts all fields to domain Fact", func() {
			now := time.Now()
			model := &Fact{
				ID:                   "fact-1",
				Text:                 "Led critical migration",
				CompetencyCategories: StringSlice{"technical", "leadership"},
				RoleFit:              "staff",
				AudienceRelevance:    StringSlice{"hiring_manager", "peer"},
				StrengthSignal:       "high",
				SourceEventID:        "evt-src",
				SourceBurstID:        "burst-src",
				CreatedAt:            now,
				UpdatedAt:            now,
			}

			result := model.ToDomain()

			Expect(result.ID).To(Equal("fact-1"))
			Expect(result.Text).To(Equal("Led critical migration"))
			Expect(result.CompetencyCategories).To(Equal([]string{"technical", "leadership"}))
			Expect(result.RoleFit).To(Equal(domain.RoleFit("staff")))
			Expect(result.AudienceRelevance).To(Equal([]string{"hiring_manager", "peer"}))
			Expect(result.StrengthSignal).To(Equal("high"))
			Expect(result.SourceEventID).To(Equal("evt-src"))
			Expect(result.SourceBurstID).To(Equal("burst-src"))
			Expect(result.CreatedAt).To(Equal(now))
			Expect(result.UpdatedAt).To(Equal(now))
		})

		It("handles nil slice fields", func() {
			model := &Fact{
				ID:      "fact-2",
				Text:    "Minimal fact",
				RoleFit: "senior_ic",
			}

			result := model.ToDomain()

			Expect(result.CompetencyCategories).To(BeNil())
			Expect(result.AudienceRelevance).To(BeNil())
		})
	})

	Describe("FactFromDomain", func() {
		It("converts all fields from domain Fact", func() {
			domainFact := fixtures.FactWithCategories("fact-d1", "Architected platform", "evt-1", []string{"technical"}, []string{"peer"})
			domainFact.RoleFit = domain.RoleFitPrincipal
			domainFact.StrengthSignal = "leadership capability"

			model := FactFromDomain(domainFact)

			Expect(model.ID).To(Equal("fact-d1"))
			Expect(model.Text).To(Equal("Architected platform"))
			Expect([]string(model.CompetencyCategories)).To(Equal([]string{"technical"}))
			Expect(model.RoleFit).To(Equal(string(domain.RoleFitPrincipal)))
			Expect([]string(model.AudienceRelevance)).To(Equal([]string{"peer"}))
			Expect(model.StrengthSignal).To(Equal("leadership capability"))
			Expect(model.SourceEventID).To(Equal("evt-1"))
			Expect(model.CreatedAt).To(Equal(domainFact.CreatedAt))
			Expect(model.UpdatedAt).To(Equal(domainFact.UpdatedAt))
		})

		It("handles nil slice fields", func() {
			domainFact := fixtures.FactWith("fact-d2", "Minimal fact")
			domainFact.RoleFit = domain.RoleFitStaff

			model := FactFromDomain(domainFact)

			Expect(model.CompetencyCategories).To(BeNil())
			Expect(model.AudienceRelevance).To(BeNil())
		})
	})

	Describe("ToDomain-FromDomain round trip", func() {
		It("preserves all Fact data", func() {
			original := fixtures.FactWithCategories("fact-rt", "Round trip fact", "evt-rt", []string{"technical", "mentoring"}, []string{"hiring_manager", "recruiter"})
			original.RoleFit = domain.RoleFitEM
			original.StrengthSignal = "team growth"
			original.SourceBurstID = "burst-rt"
			original.CreatedAt = original.CreatedAt.Truncate(time.Second)
			original.UpdatedAt = original.UpdatedAt.Truncate(time.Second)

			model := FactFromDomain(original)
			restored := model.ToDomain()

			Expect(restored.ID).To(Equal(original.ID))
			Expect(restored.Text).To(Equal(original.Text))
			Expect(restored.CompetencyCategories).To(Equal(original.CompetencyCategories))
			Expect(restored.RoleFit).To(Equal(original.RoleFit))
			Expect(restored.AudienceRelevance).To(Equal(original.AudienceRelevance))
			Expect(restored.StrengthSignal).To(Equal(original.StrengthSignal))
			Expect(restored.SourceEventID).To(Equal(original.SourceEventID))
			Expect(restored.SourceBurstID).To(Equal(original.SourceBurstID))
			Expect(restored.CreatedAt).To(Equal(original.CreatedAt))
			Expect(restored.UpdatedAt).To(Equal(original.UpdatedAt))
		})
	})
})
