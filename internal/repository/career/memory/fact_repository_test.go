//nolint:errcheck // Test file - error handling for test setup is not relevant.
package memory

import (
	"context"
	"time"

	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"

	"github.com/baphled/kariya/internal/domain/career"
	career_repo "github.com/baphled/kariya/internal/repository/career"
	"github.com/baphled/kariya/internal/testutil/fixtures"
)

var _ = Describe("FactRepository", func() {
	var (
		repository career_repo.FactRepository
		ctx        context.Context
	)

	BeforeEach(func() {
		repository = NewFactRepository()
		ctx = context.Background()
	})

	Describe("Create", func() {
		It("should create a new fact", func() {
			fact := fixtures.FactFactory.MustCreate().(*career.Fact)
			fact.ID = "" // Clear to test auto-generation

			err := repository.Create(ctx, fact)
			Expect(err).ToNot(HaveOccurred())
			Expect(fact.ID).ToNot(BeEmpty())
			Expect(fact.CreatedAt).ToNot(BeZero())
			Expect(fact.UpdatedAt).ToNot(BeZero())
		})

		It("should return error for duplicate fact ID", func() {
			fact := fixtures.Fact("fact-1", "event-1")

			err := repository.Create(ctx, fact)
			Expect(err).ToNot(HaveOccurred())

			// Try to create again with same ID
			duplicate := fixtures.Fact("fact-1", "event-2")

			err = repository.Create(ctx, duplicate)
			Expect(err).To(MatchError(career_repo.ErrDuplicateFact))
		})

		It("should generate unique ID if not provided", func() {
			fact1 := fixtures.FactFactory.MustCreate().(*career.Fact)
			fact1.ID = "" // Clear to test auto-generation

			fact2 := fixtures.FactFactory.MustCreate().(*career.Fact)
			fact2.ID = "" // Clear to test auto-generation

			err1 := repository.Create(ctx, fact1)
			err2 := repository.Create(ctx, fact2)

			Expect(err1).ToNot(HaveOccurred())
			Expect(err2).ToNot(HaveOccurred())
			Expect(fact1.ID).ToNot(BeEmpty())
			Expect(fact2.ID).ToNot(BeEmpty())
			Expect(fact1.ID).ToNot(Equal(fact2.ID))
		})

		It("should create fact with burst source", func() {
			fact := fixtures.FactFromBurst("", "burst-1")

			err := repository.Create(ctx, fact)
			Expect(err).ToNot(HaveOccurred())
			Expect(fact.SourceBurstID).To(Equal("burst-1"))
			Expect(fact.SourceEventID).To(BeEmpty())
		})
	})

	Describe("GetByID", func() {
		It("should retrieve an existing fact", func() {
			fact := fixtures.FactFactory.MustCreate().(*career.Fact)
			fact.ID = "" // Clear to test auto-generation

			err := repository.Create(ctx, fact)
			Expect(err).ToNot(HaveOccurred())

			retrieved, err := repository.GetByID(ctx, fact.ID)
			Expect(err).ToNot(HaveOccurred())
			Expect(retrieved.ID).To(Equal(fact.ID))
			Expect(retrieved.Text).To(Equal(fact.Text))
			Expect(retrieved.CompetencyCategories).To(Equal(fact.CompetencyCategories))
			Expect(retrieved.RoleFit).To(Equal(fact.RoleFit))
			Expect(retrieved.AudienceRelevance).To(Equal(fact.AudienceRelevance))
			Expect(retrieved.StrengthSignal).To(Equal(fact.StrengthSignal))
			Expect(retrieved.SourceEventID).To(Equal(fact.SourceEventID))
		})

		It("should return error for non-existent fact", func() {
			_, err := repository.GetByID(ctx, "non-existent-id")
			Expect(err).To(MatchError(career_repo.ErrFactNotFound))
		})
	})

	Describe("Update", func() {
		It("should update an existing fact", func() {
			fact := fixtures.Fact("", "event-1")

			err := repository.Create(ctx, fact)
			Expect(err).ToNot(HaveOccurred())

			originalUpdatedAt := fact.UpdatedAt
			time.Sleep(10 * time.Millisecond) // Ensure timestamp difference

			fact.Text = "Updated: Led migration and mentored team"
			fact.CompetencyCategories = []string{"technical", "leadership", "mentoring"}
			fact.RoleFit = career.RoleFitPrincipal
			fact.AudienceRelevance = []string{"hiring_manager", "peer"}
			fact.StrengthSignal = "leadership"

			err = repository.Update(ctx, fact)
			Expect(err).ToNot(HaveOccurred())
			Expect(fact.UpdatedAt.After(originalUpdatedAt)).To(BeTrue())

			retrieved, err := repository.GetByID(ctx, fact.ID)
			Expect(err).ToNot(HaveOccurred())
			Expect(retrieved.Text).To(ContainSubstring("Updated"))
			Expect(retrieved.CompetencyCategories).To(HaveLen(3))
			Expect(retrieved.RoleFit).To(Equal(career.RoleFitPrincipal))
			Expect(retrieved.AudienceRelevance).To(HaveLen(2))
			Expect(retrieved.StrengthSignal).To(Equal("leadership"))
		})

		It("should return error for non-existent fact", func() {
			fact := fixtures.Fact("non-existent-id", "event-1")

			err := repository.Update(ctx, fact)
			Expect(err).To(MatchError(career_repo.ErrFactNotFound))
		})
	})

	Describe("Delete", func() {
		It("should delete an existing fact", func() {
			fact := fixtures.FactFactory.MustCreate().(*career.Fact)
			fact.ID = "" // Clear to test auto-generation

			err := repository.Create(ctx, fact)
			Expect(err).ToNot(HaveOccurred())

			err = repository.Delete(ctx, fact.ID)
			Expect(err).ToNot(HaveOccurred())

			_, err = repository.GetByID(ctx, fact.ID)
			Expect(err).To(MatchError(career_repo.ErrFactNotFound))
		})

		It("should return error for non-existent fact", func() {
			err := repository.Delete(ctx, "non-existent-id")
			Expect(err).To(MatchError(career_repo.ErrFactNotFound))
		})
	})

	Describe("List", func() {
		BeforeEach(func() {
			// Create test facts with specific attributes for filtering tests
			fact1 := fixtures.Fact("", "event-1")
			fact1.Text = "Led migration to microservices"
			fact1.CompetencyCategories = []string{"technical", "leadership"}
			fact1.RoleFit = career.RoleFitStaff
			fact1.AudienceRelevance = []string{"hiring_manager", "peer"}
			fact1.StrengthSignal = "leadership"

			fact2 := fixtures.Fact("", "event-2")
			fact2.Text = "Mentored junior engineers"
			fact2.CompetencyCategories = []string{"mentoring"}
			fact2.RoleFit = career.RoleFitSeniorIC
			fact2.AudienceRelevance = []string{"peer"}
			fact2.StrengthSignal = "mentoring"

			fact3 := fixtures.Fact("", "event-3")
			fact3.Text = "Designed product roadmap"
			fact3.CompetencyCategories = []string{"product"}
			fact3.RoleFit = career.RoleFitEM
			fact3.AudienceRelevance = []string{"hiring_manager"}
			fact3.StrengthSignal = "product vision"

			for _, fact := range []*career.Fact{fact1, fact2, fact3} {
				err := repository.Create(ctx, fact)
				Expect(err).ToNot(HaveOccurred())
			}
		})

		It("should list all facts with no filters", func() {
			facts, err := repository.List(ctx, *fixtures.FactListFilters())
			Expect(err).ToNot(HaveOccurred())
			Expect(facts).To(HaveLen(3))
		})

		It("should filter facts by competency category", func() {
			filters := fixtures.FactListFiltersWithCategory("mentoring")

			facts, err := repository.List(ctx, *filters)
			Expect(err).ToNot(HaveOccurred())
			Expect(facts).To(HaveLen(1))
			Expect(facts[0].CompetencyCategories).To(ContainElement("mentoring"))
		})

		It("should filter facts by role fit", func() {
			filters := fixtures.FactListFiltersWithRole(string(career.RoleFitEM))

			facts, err := repository.List(ctx, *filters)
			Expect(err).ToNot(HaveOccurred())
			Expect(facts).To(HaveLen(1))
			Expect(facts[0].RoleFit).To(Equal(career.RoleFitEM))
		})

		It("should filter facts by audience relevance", func() {
			filters := fixtures.FactListFiltersWithAudience("peer")

			facts, err := repository.List(ctx, *filters)
			Expect(err).ToNot(HaveOccurred())
			Expect(facts).To(HaveLen(2)) // Two facts have "peer" audience
		})

		It("should sort facts by creation date descending", func() {
			filters := fixtures.FactListFiltersWithSort("created_at", "desc")

			facts, err := repository.List(ctx, *filters)
			Expect(err).ToNot(HaveOccurred())
			Expect(facts).To(HaveLen(3))
			for i := range len(facts) - 1 {
				Expect(facts[i].CreatedAt.After(facts[i+1].CreatedAt) || facts[i].CreatedAt.Equal(facts[i+1].CreatedAt)).To(BeTrue())
			}
		})

		It("should apply pagination with limit", func() {
			filters := fixtures.FactListFiltersWithLimit(0, 2)

			facts, err := repository.List(ctx, *filters)
			Expect(err).ToNot(HaveOccurred())
			Expect(facts).To(HaveLen(2))
		})

		It("should apply pagination with offset", func() {
			filters := fixtures.FactListFiltersWithLimit(1, 0)

			facts, err := repository.List(ctx, *filters)
			Expect(err).ToNot(HaveOccurred())
			Expect(facts).To(HaveLen(2))
		})

		It("should apply pagination with both limit and offset", func() {
			filters := fixtures.FactListFiltersWithLimit(1, 1)

			facts, err := repository.List(ctx, *filters)
			Expect(err).ToNot(HaveOccurred())
			Expect(facts).To(HaveLen(1))
		})
	})

	Describe("Count", func() {
		BeforeEach(func() {
			// Create test facts with specific competencies for counting tests
			fact1 := fixtures.Fact("", "event-1")
			fact1.CompetencyCategories = []string{"technical"}
			fact1.RoleFit = career.RoleFitStaff

			fact2 := fixtures.Fact("", "event-2")
			fact2.CompetencyCategories = []string{"technical"}
			fact2.RoleFit = career.RoleFitStaff

			fact3 := fixtures.Fact("", "event-3")
			fact3.CompetencyCategories = []string{"leadership"}
			fact3.RoleFit = career.RoleFitEM

			for _, fact := range []*career.Fact{fact1, fact2, fact3} {
				err := repository.Create(ctx, fact)
				Expect(err).ToNot(HaveOccurred())
			}
		})

		It("should count all facts with no filters", func() {
			count, err := repository.Count(ctx, *fixtures.FactListFilters())
			Expect(err).ToNot(HaveOccurred())
			Expect(count).To(Equal(3))
		})

		It("should count facts matching competency filter", func() {
			filters := fixtures.FactListFiltersWithCategory("technical")

			count, err := repository.Count(ctx, *filters)
			Expect(err).ToNot(HaveOccurred())
			Expect(count).To(Equal(2))
		})

		It("should count facts matching role fit filter", func() {
			filters := fixtures.FactListFiltersWithRole(string(career.RoleFitEM))

			count, err := repository.Count(ctx, *filters)
			Expect(err).ToNot(HaveOccurred())
			Expect(count).To(Equal(1))
		})

		It("should return zero for no matches", func() {
			filters := fixtures.FactListFiltersWithCategory("non-existent")

			count, err := repository.Count(ctx, *filters)
			Expect(err).ToNot(HaveOccurred())
			Expect(count).To(Equal(0))
		})
	})

	Describe("GetBySourceEventID", func() {
		BeforeEach(func() {
			// Create test facts with specific source event IDs
			fact1 := fixtures.Fact("", "event-1")
			fact2 := fixtures.Fact("", "event-1") // Same event
			fact3 := fixtures.Fact("", "event-2") // Different event

			for _, fact := range []*career.Fact{fact1, fact2, fact3} {
				err := repository.Create(ctx, fact)
				Expect(err).ToNot(HaveOccurred())
			}
		})

		It("should retrieve facts for a specific event", func() {
			facts, err := repository.GetBySourceEventID(ctx, "event-1")
			Expect(err).ToNot(HaveOccurred())
			Expect(facts).To(HaveLen(2))
			for _, fact := range facts {
				Expect(fact.SourceEventID).To(Equal("event-1"))
			}
		})

		It("should return empty list for event with no facts", func() {
			facts, err := repository.GetBySourceEventID(ctx, "non-existent-event")
			Expect(err).ToNot(HaveOccurred())
			Expect(facts).To(BeEmpty())
		})
	})

	Describe("GetBySourceBurstID", func() {
		BeforeEach(func() {
			// Create test facts with specific source burst IDs
			fact1 := fixtures.FactFromBurst("", "burst-1")
			fact2 := fixtures.FactFromBurst("", "burst-1") // Same burst
			fact3 := fixtures.FactFromBurst("", "burst-2") // Different burst

			for _, fact := range []*career.Fact{fact1, fact2, fact3} {
				err := repository.Create(ctx, fact)
				Expect(err).ToNot(HaveOccurred())
			}
		})

		It("should retrieve facts for a specific burst", func() {
			facts, err := repository.GetBySourceBurstID(ctx, "burst-1")
			Expect(err).ToNot(HaveOccurred())
			Expect(facts).To(HaveLen(2))
			for _, fact := range facts {
				Expect(fact.SourceBurstID).To(Equal("burst-1"))
			}
		})

		It("should return empty list for burst with no facts", func() {
			facts, err := repository.GetBySourceBurstID(ctx, "non-existent-burst")
			Expect(err).ToNot(HaveOccurred())
			Expect(facts).To(BeEmpty())
		})
	})
})
