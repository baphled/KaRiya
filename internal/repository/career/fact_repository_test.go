package career_test

import (
	"context"
	"time"

	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"

	"github.com/baphled/kariya/internal/domain/career"
	repo "github.com/baphled/kariya/internal/repository/career"
)

var _ = Describe("FactRepository", func() {
	var (
		repository repo.FactRepository
		ctx        context.Context
	)

	BeforeEach(func() {
		repository = repo.NewMemoryFactRepository()
		ctx = context.Background()
	})

	Describe("Create", func() {
		It("should create a new fact", func() {
			fact := &career.Fact{
				Text:                "Led migration of platform to microservices architecture",
				CompetencyCategories: []string{"technical", "leadership"},
				RoleFit:             career.RoleFitStaff,
				AudienceRelevance:   []string{"hiring_manager", "peer"},
				StrengthSignal:      "leadership",
				SourceEventID:       "event-1",
			}

			err := repository.Create(ctx, fact)
			Expect(err).ToNot(HaveOccurred())
			Expect(fact.ID).ToNot(BeEmpty())
			Expect(fact.CreatedAt).ToNot(BeZero())
			Expect(fact.UpdatedAt).ToNot(BeZero())
		})

		It("should return error for duplicate fact ID", func() {
			fact := &career.Fact{
				ID:                  "fact-1",
				Text:                "Led migration of platform to microservices architecture",
				CompetencyCategories: []string{"technical"},
				RoleFit:             career.RoleFitStaff,
				AudienceRelevance:   []string{"hiring_manager"},
				StrengthSignal:      "leadership",
				SourceEventID:       "event-1",
			}

			err := repository.Create(ctx, fact)
			Expect(err).ToNot(HaveOccurred())

			// Try to create again with same ID
			duplicate := &career.Fact{
				ID:                  "fact-1",
				Text:                "Different fact text",
				CompetencyCategories: []string{"leadership"},
				RoleFit:             career.RoleFitEM,
				AudienceRelevance:   []string{"recruiter"},
				StrengthSignal:      "management",
				SourceEventID:       "event-2",
			}

			err = repository.Create(ctx, duplicate)
			Expect(err).To(MatchError(repo.ErrDuplicateFact))
		})

		It("should generate unique ID if not provided", func() {
			fact1 := &career.Fact{
				Text:                "Fact 1",
				CompetencyCategories: []string{"technical"},
				RoleFit:             career.RoleFitStaff,
				AudienceRelevance:   []string{"peer"},
				StrengthSignal:      "technical",
				SourceEventID:       "event-1",
			}

			fact2 := &career.Fact{
				Text:                "Fact 2",
				CompetencyCategories: []string{"leadership"},
				RoleFit:             career.RoleFitEM,
				AudienceRelevance:   []string{"hiring_manager"},
				StrengthSignal:      "leadership",
				SourceEventID:       "event-2",
			}

			err1 := repository.Create(ctx, fact1)
			err2 := repository.Create(ctx, fact2)

			Expect(err1).ToNot(HaveOccurred())
			Expect(err2).ToNot(HaveOccurred())
			Expect(fact1.ID).ToNot(BeEmpty())
			Expect(fact2.ID).ToNot(BeEmpty())
			Expect(fact1.ID).ToNot(Equal(fact2.ID))
		})

		It("should create fact with burst source", func() {
			fact := &career.Fact{
				Text:                "Coordinated team efforts across multiple projects",
				CompetencyCategories: []string{"leadership"},
				RoleFit:             career.RoleFitEM,
				AudienceRelevance:   []string{"hiring_manager"},
				StrengthSignal:      "management",
				SourceBurstID:       "burst-1",
			}

			err := repository.Create(ctx, fact)
			Expect(err).ToNot(HaveOccurred())
			Expect(fact.SourceBurstID).To(Equal("burst-1"))
			Expect(fact.SourceEventID).To(BeEmpty())
		})
	})

	Describe("GetByID", func() {
		It("should retrieve an existing fact", func() {
			fact := &career.Fact{
				Text:                "Led migration of platform to microservices architecture",
				CompetencyCategories: []string{"technical", "leadership"},
				RoleFit:             career.RoleFitStaff,
				AudienceRelevance:   []string{"hiring_manager", "peer"},
				StrengthSignal:      "leadership",
				SourceEventID:       "event-1",
			}

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
			Expect(err).To(MatchError(repo.ErrFactNotFound))
		})
	})

	Describe("Update", func() {
		It("should update an existing fact", func() {
			fact := &career.Fact{
				Text:                "Led migration of platform to microservices architecture",
				CompetencyCategories: []string{"technical"},
				RoleFit:             career.RoleFitStaff,
				AudienceRelevance:   []string{"peer"},
				StrengthSignal:      "technical",
				SourceEventID:       "event-1",
			}

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
			fact := &career.Fact{
				ID:                  "non-existent-id",
				Text:                "Some fact text",
				CompetencyCategories: []string{"technical"},
				RoleFit:             career.RoleFitStaff,
				AudienceRelevance:   []string{"peer"},
				StrengthSignal:      "technical",
				SourceEventID:       "event-1",
			}

			err := repository.Update(ctx, fact)
			Expect(err).To(MatchError(repo.ErrFactNotFound))
		})
	})

	Describe("Delete", func() {
		It("should delete an existing fact", func() {
			fact := &career.Fact{
				Text:                "Led migration of platform to microservices architecture",
				CompetencyCategories: []string{"technical"},
				RoleFit:             career.RoleFitStaff,
				AudienceRelevance:   []string{"peer"},
				StrengthSignal:      "technical",
				SourceEventID:       "event-1",
			}

			err := repository.Create(ctx, fact)
			Expect(err).ToNot(HaveOccurred())

			err = repository.Delete(ctx, fact.ID)
			Expect(err).ToNot(HaveOccurred())

			_, err = repository.GetByID(ctx, fact.ID)
			Expect(err).To(MatchError(repo.ErrFactNotFound))
		})

		It("should return error for non-existent fact", func() {
			err := repository.Delete(ctx, "non-existent-id")
			Expect(err).To(MatchError(repo.ErrFactNotFound))
		})
	})

	Describe("List", func() {
		BeforeEach(func() {
			// Create test facts
			facts := []*career.Fact{
				{
					Text:                "Led migration to microservices",
					CompetencyCategories: []string{"technical", "leadership"},
					RoleFit:             career.RoleFitStaff,
					AudienceRelevance:   []string{"hiring_manager", "peer"},
					StrengthSignal:      "leadership",
					SourceEventID:       "event-1",
				},
				{
					Text:                "Mentored junior engineers",
					CompetencyCategories: []string{"mentoring"},
					RoleFit:             career.RoleFitSeniorIC,
					AudienceRelevance:   []string{"peer"},
					StrengthSignal:      "mentoring",
					SourceEventID:       "event-2",
				},
				{
					Text:                "Designed product roadmap",
					CompetencyCategories: []string{"product"},
					RoleFit:             career.RoleFitEM,
					AudienceRelevance:   []string{"hiring_manager"},
					StrengthSignal:      "product vision",
					SourceEventID:       "event-3",
				},
			}

			for _, fact := range facts {
				err := repository.Create(ctx, fact)
				Expect(err).ToNot(HaveOccurred())
			}
		})

		It("should list all facts with no filters", func() {
			facts, err := repository.List(ctx, repo.FactListFilters{})
			Expect(err).ToNot(HaveOccurred())
			Expect(facts).To(HaveLen(3))
		})

		It("should filter facts by competency category", func() {
			filters := repo.FactListFilters{
				CompetencyCategory: "mentoring",
			}

			facts, err := repository.List(ctx, filters)
			Expect(err).ToNot(HaveOccurred())
			Expect(facts).To(HaveLen(1))
			Expect(facts[0].CompetencyCategories).To(ContainElement("mentoring"))
		})

		It("should filter facts by role fit", func() {
			filters := repo.FactListFilters{
				RoleFit: string(career.RoleFitEM),
			}

			facts, err := repository.List(ctx, filters)
			Expect(err).ToNot(HaveOccurred())
			Expect(facts).To(HaveLen(1))
			Expect(facts[0].RoleFit).To(Equal(career.RoleFitEM))
		})

		It("should filter facts by audience relevance", func() {
			filters := repo.FactListFilters{
				AudienceRelevance: "peer",
			}

			facts, err := repository.List(ctx, filters)
			Expect(err).ToNot(HaveOccurred())
			Expect(facts).To(HaveLen(2)) // Two facts have "peer" audience
		})

		It("should sort facts by creation date descending", func() {
			filters := repo.FactListFilters{
				SortBy:    "created_at",
				SortOrder: "desc",
			}

			facts, err := repository.List(ctx, filters)
			Expect(err).ToNot(HaveOccurred())
			Expect(facts).To(HaveLen(3))
			// Most recent first
			for i := 0; i < len(facts)-1; i++ {
				Expect(facts[i].CreatedAt.After(facts[i+1].CreatedAt) || facts[i].CreatedAt.Equal(facts[i+1].CreatedAt)).To(BeTrue())
			}
		})

		It("should apply pagination with limit", func() {
			filters := repo.FactListFilters{
				Limit: 2,
			}

			facts, err := repository.List(ctx, filters)
			Expect(err).ToNot(HaveOccurred())
			Expect(facts).To(HaveLen(2))
		})

		It("should apply pagination with offset", func() {
			filters := repo.FactListFilters{
				Offset: 1,
			}

			facts, err := repository.List(ctx, filters)
			Expect(err).ToNot(HaveOccurred())
			Expect(facts).To(HaveLen(2))
		})

		It("should apply pagination with both limit and offset", func() {
			filters := repo.FactListFilters{
				Offset: 1,
				Limit:  1,
			}

			facts, err := repository.List(ctx, filters)
			Expect(err).ToNot(HaveOccurred())
			Expect(facts).To(HaveLen(1))
		})
	})

	Describe("Count", func() {
		BeforeEach(func() {
			// Create test facts
			facts := []*career.Fact{
				{
					Text:                "Technical fact 1",
					CompetencyCategories: []string{"technical"},
					RoleFit:             career.RoleFitStaff,
					AudienceRelevance:   []string{"peer"},
					StrengthSignal:      "technical",
					SourceEventID:       "event-1",
				},
				{
					Text:                "Technical fact 2",
					CompetencyCategories: []string{"technical"},
					RoleFit:             career.RoleFitStaff,
					AudienceRelevance:   []string{"peer"},
					StrengthSignal:      "technical",
					SourceEventID:       "event-2",
				},
				{
					Text:                "Leadership fact",
					CompetencyCategories: []string{"leadership"},
					RoleFit:             career.RoleFitEM,
					AudienceRelevance:   []string{"hiring_manager"},
					StrengthSignal:      "leadership",
					SourceEventID:       "event-3",
				},
			}

			for _, fact := range facts {
				err := repository.Create(ctx, fact)
				Expect(err).ToNot(HaveOccurred())
			}
		})

		It("should count all facts with no filters", func() {
			count, err := repository.Count(ctx, repo.FactListFilters{})
			Expect(err).ToNot(HaveOccurred())
			Expect(count).To(Equal(3))
		})

		It("should count facts matching competency filter", func() {
			filters := repo.FactListFilters{
				CompetencyCategory: "technical",
			}

			count, err := repository.Count(ctx, filters)
			Expect(err).ToNot(HaveOccurred())
			Expect(count).To(Equal(2))
		})

		It("should count facts matching role fit filter", func() {
			filters := repo.FactListFilters{
				RoleFit: string(career.RoleFitEM),
			}

			count, err := repository.Count(ctx, filters)
			Expect(err).ToNot(HaveOccurred())
			Expect(count).To(Equal(1))
		})

		It("should return zero for no matches", func() {
			filters := repo.FactListFilters{
				CompetencyCategory: "non-existent",
			}

			count, err := repository.Count(ctx, filters)
			Expect(err).ToNot(HaveOccurred())
			Expect(count).To(Equal(0))
		})
	})

	Describe("GetBySourceEventID", func() {
		BeforeEach(func() {
			// Create test facts with different sources
			facts := []*career.Fact{
				{
					Text:                "Fact from event 1",
					CompetencyCategories: []string{"technical"},
					RoleFit:             career.RoleFitStaff,
					AudienceRelevance:   []string{"peer"},
					StrengthSignal:      "technical",
					SourceEventID:       "event-1",
				},
				{
					Text:                "Another fact from event 1",
					CompetencyCategories: []string{"leadership"},
					RoleFit:             career.RoleFitEM,
					AudienceRelevance:   []string{"hiring_manager"},
					StrengthSignal:      "leadership",
					SourceEventID:       "event-1",
				},
				{
					Text:                "Fact from event 2",
					CompetencyCategories: []string{"product"},
					RoleFit:             career.RoleFitSeniorIC,
					AudienceRelevance:   []string{"peer"},
					StrengthSignal:      "product",
					SourceEventID:       "event-2",
				},
			}

			for _, fact := range facts {
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
			// Create test facts with different burst sources
			facts := []*career.Fact{
				{
					Text:                "Fact from burst 1",
					CompetencyCategories: []string{"technical"},
					RoleFit:             career.RoleFitStaff,
					AudienceRelevance:   []string{"peer"},
					StrengthSignal:      "technical",
					SourceBurstID:       "burst-1",
				},
				{
					Text:                "Another fact from burst 1",
					CompetencyCategories: []string{"leadership"},
					RoleFit:             career.RoleFitEM,
					AudienceRelevance:   []string{"hiring_manager"},
					StrengthSignal:      "leadership",
					SourceBurstID:       "burst-1",
				},
				{
					Text:                "Fact from burst 2",
					CompetencyCategories: []string{"product"},
					RoleFit:             career.RoleFitSeniorIC,
					AudienceRelevance:   []string{"peer"},
					StrengthSignal:      "product",
					SourceBurstID:       "burst-2",
				},
			}

			for _, fact := range facts {
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

