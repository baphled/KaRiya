package fixtures_test

import (
	"github.com/baphled/kariya/internal/domain/career"
	"github.com/baphled/kariya/internal/testutil/fixtures"
	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
)

var _ = Describe("FactFactory", func() {
	BeforeEach(func() {
		fixtures.SetSeed(42)
	})

	Describe("MustCreate", func() {
		It("should create a fact with all required fields", func() {
			fact := fixtures.FactFactory.MustCreate().(*career.Fact)
			Expect(fact).NotTo(BeNil())
			Expect(fact.ID).NotTo(BeEmpty())
			Expect(fact.Text).NotTo(BeEmpty())
			Expect(fact.CompetencyCategories).NotTo(BeEmpty())
			Expect(fact.RoleFit).NotTo(BeEmpty())
			Expect(fact.AudienceRelevance).NotTo(BeEmpty())
			Expect(fact.StrengthSignal).NotTo(BeEmpty())
			Expect(fact.SourceEventID).NotTo(BeEmpty())
			Expect(fact.CreatedAt).NotTo(BeZero())
			Expect(fact.UpdatedAt).NotTo(BeZero())
		})

		It("should create facts with unique sequential IDs", func() {
			fact1 := fixtures.FactFactory.MustCreate().(*career.Fact)
			fact2 := fixtures.FactFactory.MustCreate().(*career.Fact)
			Expect(fact1.ID).NotTo(Equal(fact2.ID))
		})

		It("should create facts with valid role fit values", func() {
			fact := fixtures.FactFactory.MustCreate().(*career.Fact)
			validRoleFits := []career.RoleFit{
				career.RoleFitStaff,
				career.RoleFitSeniorIC,
				career.RoleFitPrincipal,
				career.RoleFitEM,
			}
			Expect(validRoleFits).To(ContainElement(fact.RoleFit))
		})
	})
})

var _ = Describe("Fact Quick Helper", func() {
	Describe("Fact", func() {
		It("should create a minimal valid fact with given ID and source", func() {
			fact := fixtures.Fact("fact-1", "evt-1")
			Expect(fact.ID).To(Equal("fact-1"))
			Expect(fact.SourceEventID).To(Equal("evt-1"))
			Expect(fact.Text).To(ContainSubstring("fact-1"))
			Expect(fact.CompetencyCategories).To(Equal([]string{"technical"}))
			Expect(fact.RoleFit).To(Equal(career.RoleFitStaff))
			Expect(fact.AudienceRelevance).To(Equal([]string{"hiring_manager"}))
			Expect(fact.StrengthSignal).To(Equal("high"))
		})
	})

	Describe("FactFromBurst", func() {
		It("should create a fact linked to a burst", func() {
			fact := fixtures.FactFromBurst("fact-1", "burst-1")
			Expect(fact.SourceBurstID).To(Equal("burst-1"))
			Expect(fact.SourceEventID).To(BeEmpty())
		})
	})
})

var _ = Describe("Facts Batch Helper", func() {
	BeforeEach(func() {
		fixtures.SetSeed(42)
	})

	Describe("Facts", func() {
		It("should create n facts", func() {
			events := fixtures.Events(3)
			facts := fixtures.Facts(5, events)
			Expect(facts).To(HaveLen(5))
		})

		It("should link facts to provided events", func() {
			events := fixtures.Events(2)
			facts := fixtures.Facts(4, events)
			// Facts should cycle through available events
			Expect(facts[0].SourceEventID).To(Equal(events[0].ID))
			Expect(facts[1].SourceEventID).To(Equal(events[1].ID))
			Expect(facts[2].SourceEventID).To(Equal(events[0].ID))
		})
	})
})
