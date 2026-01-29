package sql

import (
	"context"
	stdsql "database/sql"
	"time"

	"github.com/baphled/kariya/internal/domain/career"
	career_repo "github.com/baphled/kariya/internal/repository/career"
	"github.com/baphled/kariya/internal/repository/models"
	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
	_ "modernc.org/sqlite"
)

func setupFactTestDB() *gorm.DB {
	sqlDB, err := stdsql.Open("sqlite", ":memory:")
	Expect(err).NotTo(HaveOccurred())

	db, err := gorm.Open(sqlite.New(sqlite.Config{
		Conn: sqlDB,
	}), &gorm.Config{})
	Expect(err).NotTo(HaveOccurred())

	err = db.AutoMigrate(&models.Fact{})
	Expect(err).NotTo(HaveOccurred())

	return db
}

var _ = Describe("Fact Repository", func() {
	var (
		repo *FactRepository
		db   *gorm.DB
		ctx  context.Context
	)

	BeforeEach(func() {
		db = setupFactTestDB()
		repo = NewFactRepository(db)
		ctx = context.Background()
	})

	Describe("Create", func() {
		It("creates a fact with generated ID", func() {
			fact := &career.Fact{
				Text:                 "Led team of 5 engineers",
				CompetencyCategories: []string{"leadership"},
				RoleFit:              career.RoleFitSeniorIC,
				AudienceRelevance:    []string{"technical"},
			}

			err := repo.Create(ctx, fact)

			Expect(err).NotTo(HaveOccurred())
			Expect(fact.ID).NotTo(BeEmpty())
			Expect(fact.CreatedAt).NotTo(BeZero())
		})

		It("creates a fact with provided ID", func() {
			fact := &career.Fact{
				ID:                   "custom-id",
				Text:                 "Built microservice",
				CompetencyCategories: []string{"technical"},
				RoleFit:              career.RoleFitSeniorIC,
				AudienceRelevance:    []string{"engineering"},
			}

			err := repo.Create(ctx, fact)

			Expect(err).NotTo(HaveOccurred())
			Expect(fact.ID).To(Equal("custom-id"))
		})
	})

	Describe("GetByID", func() {
		It("returns the fact", func() {
			fact := &career.Fact{
				Text:                 "Test fact",
				CompetencyCategories: []string{"cat1"},
				RoleFit:              career.RoleFitSeniorIC,
				AudienceRelevance:    []string{"aud1"},
			}
			Expect(repo.Create(ctx, fact)).To(Succeed())

			found, err := repo.GetByID(ctx, fact.ID)

			Expect(err).NotTo(HaveOccurred())
			Expect(found.Text).To(Equal("Test fact"))
			Expect(found.CompetencyCategories).To(ConsistOf("cat1"))
		})

		It("returns ErrFactNotFound for missing fact", func() {
			_, err := repo.GetByID(ctx, "nonexistent")

			Expect(err).To(Equal(career_repo.ErrFactNotFound))
		})
	})

	Describe("Update", func() {
		It("updates the fact", func() {
			fact := &career.Fact{
				Text:                 "Original",
				CompetencyCategories: []string{"cat"},
				RoleFit:              career.RoleFitSeniorIC,
				AudienceRelevance:    []string{"aud"},
			}
			Expect(repo.Create(ctx, fact)).To(Succeed())

			fact.Text = "Updated"
			err := repo.Update(ctx, fact)

			Expect(err).NotTo(HaveOccurred())

			found, _ := repo.GetByID(ctx, fact.ID)
			Expect(found.Text).To(Equal("Updated"))
		})

		It("returns ErrFactNotFound for missing fact", func() {
			fact := &career.Fact{
				ID:                   "nonexistent",
				Text:                 "Test",
				CompetencyCategories: []string{"cat"},
				RoleFit:              career.RoleFitSeniorIC,
				AudienceRelevance:    []string{"aud"},
			}

			err := repo.Update(ctx, fact)

			Expect(err).To(Equal(career_repo.ErrFactNotFound))
		})
	})

	Describe("Delete", func() {
		It("deletes the fact", func() {
			fact := &career.Fact{
				Text:                 "To delete",
				CompetencyCategories: []string{"cat"},
				RoleFit:              career.RoleFitSeniorIC,
				AudienceRelevance:    []string{"aud"},
			}
			Expect(repo.Create(ctx, fact)).To(Succeed())

			err := repo.Delete(ctx, fact.ID)

			Expect(err).NotTo(HaveOccurred())

			_, err = repo.GetByID(ctx, fact.ID)
			Expect(err).To(Equal(career_repo.ErrFactNotFound))
		})

		It("returns ErrFactNotFound for missing fact", func() {
			err := repo.Delete(ctx, "nonexistent")

			Expect(err).To(Equal(career_repo.ErrFactNotFound))
		})
	})

	Describe("List", func() {
		BeforeEach(func() {
			facts := []*career.Fact{
				{Text: "Fact 1", CompetencyCategories: []string{"leadership"}, RoleFit: career.RoleFitSeniorIC, AudienceRelevance: []string{"technical"}},
				{Text: "Fact 2", CompetencyCategories: []string{"technical"}, RoleFit: career.RoleFitEM, AudienceRelevance: []string{"hr"}},
				{Text: "Fact 3", CompetencyCategories: []string{"leadership", "technical"}, RoleFit: career.RoleFitSeniorIC, AudienceRelevance: []string{"technical", "hr"}},
			}
			for _, f := range facts {
				Expect(repo.Create(ctx, f)).To(Succeed())
				time.Sleep(10 * time.Millisecond) // Ensure different timestamps.
			}
		})

		It("returns all facts without filters", func() {
			facts, err := repo.List(ctx, career_repo.FactListFilters{})

			Expect(err).NotTo(HaveOccurred())
			Expect(facts).To(HaveLen(3))
		})

		It("filters by competency category", func() {
			facts, err := repo.List(ctx, career_repo.FactListFilters{CompetencyCategory: "leadership"})

			Expect(err).NotTo(HaveOccurred())
			Expect(facts).To(HaveLen(2))
		})

		It("filters by role fit", func() {
			facts, err := repo.List(ctx, career_repo.FactListFilters{RoleFit: string(career.RoleFitEM)})

			Expect(err).NotTo(HaveOccurred())
			Expect(facts).To(HaveLen(1))
			Expect(facts[0].Text).To(Equal("Fact 2"))
		})

		It("sorts by created_at descending by default", func() {
			facts, err := repo.List(ctx, career_repo.FactListFilters{})

			Expect(err).NotTo(HaveOccurred())
			Expect(facts[0].Text).To(Equal("Fact 3")) // Most recent first.
		})

		It("sorts by text ascending", func() {
			facts, err := repo.List(ctx, career_repo.FactListFilters{SortBy: "text", SortOrder: "asc"})

			Expect(err).NotTo(HaveOccurred())
			Expect(facts[0].Text).To(Equal("Fact 1"))
		})

		It("applies pagination", func() {
			facts, err := repo.List(ctx, career_repo.FactListFilters{Limit: 2})

			Expect(err).NotTo(HaveOccurred())
			Expect(facts).To(HaveLen(2))
		})
	})

	Describe("Count", func() {
		BeforeEach(func() {
			facts := []*career.Fact{
				{Text: "Fact 1", CompetencyCategories: []string{"leadership"}, RoleFit: career.RoleFitSeniorIC, AudienceRelevance: []string{"a"}},
				{Text: "Fact 2", CompetencyCategories: []string{"technical"}, RoleFit: career.RoleFitSeniorIC, AudienceRelevance: []string{"b"}},
			}
			for _, f := range facts {
				Expect(repo.Create(ctx, f)).To(Succeed())
			}
		})

		It("counts all facts", func() {
			count, err := repo.Count(ctx, career_repo.FactListFilters{})

			Expect(err).NotTo(HaveOccurred())
			Expect(count).To(Equal(2))
		})

		It("counts with filters", func() {
			count, err := repo.Count(ctx, career_repo.FactListFilters{CompetencyCategory: "leadership"})

			Expect(err).NotTo(HaveOccurred())
			Expect(count).To(Equal(1))
		})
	})

	Describe("GetBySourceEventID", func() {
		It("returns facts for event", func() {
			fact1 := &career.Fact{Text: "F1", SourceEventID: "event-1", CompetencyCategories: []string{"a"}, RoleFit: career.RoleFitSeniorIC, AudienceRelevance: []string{"b"}}
			fact2 := &career.Fact{Text: "F2", SourceEventID: "event-1", CompetencyCategories: []string{"a"}, RoleFit: career.RoleFitSeniorIC, AudienceRelevance: []string{"b"}}
			fact3 := &career.Fact{Text: "F3", SourceEventID: "event-2", CompetencyCategories: []string{"a"}, RoleFit: career.RoleFitSeniorIC, AudienceRelevance: []string{"b"}}
			Expect(repo.Create(ctx, fact1)).To(Succeed())
			Expect(repo.Create(ctx, fact2)).To(Succeed())
			Expect(repo.Create(ctx, fact3)).To(Succeed())

			facts, err := repo.GetBySourceEventID(ctx, "event-1")

			Expect(err).NotTo(HaveOccurred())
			Expect(facts).To(HaveLen(2))
		})
	})

	Describe("GetBySourceBurstID", func() {
		It("returns facts for burst", func() {
			fact1 := &career.Fact{Text: "F1", SourceBurstID: "burst-1", CompetencyCategories: []string{"a"}, RoleFit: career.RoleFitSeniorIC, AudienceRelevance: []string{"b"}}
			fact2 := &career.Fact{Text: "F2", SourceBurstID: "burst-2", CompetencyCategories: []string{"a"}, RoleFit: career.RoleFitSeniorIC, AudienceRelevance: []string{"b"}}
			Expect(repo.Create(ctx, fact1)).To(Succeed())
			Expect(repo.Create(ctx, fact2)).To(Succeed())

			facts, err := repo.GetBySourceBurstID(ctx, "burst-1")

			Expect(err).NotTo(HaveOccurred())
			Expect(facts).To(HaveLen(1))
			Expect(facts[0].Text).To(Equal("F1"))
		})
	})
})
