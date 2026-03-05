package sql

import (
	"context"
	stdsql "database/sql"
	"time"

	"github.com/baphled/kariya/internal/domain/career"
	career_repo "github.com/baphled/kariya/internal/repository/career"
	"github.com/baphled/kariya/internal/repository/models"
	"github.com/baphled/kariya/internal/testutil/fixtures"
	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
)

var _ = Describe("Fact Repository", func() {
	var (
		repo *FactRepository
		tx   *gorm.DB
		ctx  context.Context
	)

	BeforeEach(func() {
		Expect(sharedGormDB).NotTo(BeNil(), "shared DB not initialized - BeforeSuite not run")
		tx = sharedGormDB.Begin()
		DeferCleanup(func() { tx.Rollback() })
		repo = NewFactRepository(tx)
		ctx = context.Background()
	})

	Describe("Create", func() {
		It("creates a fact with generated ID", func() {
			fact := fixtures.FactWithCategories("", "Led team of 5 engineers", "", []string{"leadership"}, []string{"technical"})
			fact.RoleFit = career.RoleFitSeniorIC

			err := repo.Create(ctx, fact)

			Expect(err).NotTo(HaveOccurred())
			Expect(fact.ID).NotTo(BeEmpty())
			Expect(fact.CreatedAt).NotTo(BeZero())
		})

		It("creates a fact with provided ID", func() {
			fact := fixtures.FactWithCategories("custom-id", "Built microservice", "", []string{"technical"}, []string{"engineering"})
			fact.RoleFit = career.RoleFitSeniorIC

			err := repo.Create(ctx, fact)

			Expect(err).NotTo(HaveOccurred())
			Expect(fact.ID).To(Equal("custom-id"))
		})
	})

	Describe("GetByID", func() {
		It("returns the fact", func() {
			fact := fixtures.FactWithCategories("", "Test fact", "", []string{"cat1"}, []string{"aud1"})
			fact.RoleFit = career.RoleFitSeniorIC
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
			fact := fixtures.FactWithCategories("", "Original", "", []string{"cat"}, []string{"aud"})
			fact.RoleFit = career.RoleFitSeniorIC
			Expect(repo.Create(ctx, fact)).To(Succeed())

			fact.Text = "Updated"
			err := repo.Update(ctx, fact)

			Expect(err).NotTo(HaveOccurred())

			found, _ := repo.GetByID(ctx, fact.ID)
			Expect(found.Text).To(Equal("Updated"))
		})

		It("returns ErrFactNotFound for missing fact", func() {
			fact := fixtures.FactWithCategories("nonexistent", "Test", "", []string{"cat"}, []string{"aud"})
			fact.RoleFit = career.RoleFitSeniorIC

			err := repo.Update(ctx, fact)

			Expect(err).To(Equal(career_repo.ErrFactNotFound))
		})
	})

	Describe("Delete", func() {
		It("deletes the fact", func() {
			fact := fixtures.FactWithCategories("", "To delete", "", []string{"cat"}, []string{"aud"})
			fact.RoleFit = career.RoleFitSeniorIC
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
			f1 := fixtures.FactWithCategories("", "Fact 1", "", []string{"leadership"}, []string{"technical"})
			f1.RoleFit = career.RoleFitSeniorIC
			f2 := fixtures.FactWithCategories("", "Fact 2", "", []string{"technical"}, []string{"hr"})
			f2.RoleFit = career.RoleFitEM
			f3 := fixtures.FactWithCategories("", "Fact 3", "", []string{"leadership", "technical"}, []string{"technical", "hr"})
			f3.RoleFit = career.RoleFitSeniorIC
			for _, f := range []*career.Fact{f1, f2, f3} {
				Expect(repo.Create(ctx, f)).To(Succeed())
				time.Sleep(10 * time.Millisecond) // Ensure different timestamps.
			}
		})

		It("returns all facts without filters", func() {
			facts, err := repo.List(ctx, *fixtures.FactListFilters())

			Expect(err).NotTo(HaveOccurred())
			Expect(facts).To(HaveLen(3))
		})

		It("filters by competency category", func() {
			facts, err := repo.List(ctx, *fixtures.FactListFiltersWithCategory("leadership"))

			Expect(err).NotTo(HaveOccurred())
			Expect(facts).To(HaveLen(2))
		})

		It("filters by role fit", func() {
			facts, err := repo.List(ctx, *fixtures.FactListFiltersWithRole(string(career.RoleFitEM)))

			Expect(err).NotTo(HaveOccurred())
			Expect(facts).To(HaveLen(1))
			Expect(facts[0].Text).To(Equal("Fact 2"))
		})

		It("sorts by created_at descending by default", func() {
			facts, err := repo.List(ctx, *fixtures.FactListFilters())

			Expect(err).NotTo(HaveOccurred())
			Expect(facts[0].Text).To(Equal("Fact 3")) // Most recent first.
		})

		It("sorts by text ascending", func() {
			facts, err := repo.List(ctx, *fixtures.FactListFiltersWithSort("text", "asc"))

			Expect(err).NotTo(HaveOccurred())
			Expect(facts[0].Text).To(Equal("Fact 1"))
		})

		It("applies pagination", func() {
			facts, err := repo.List(ctx, *fixtures.FactListFiltersWithLimit(0, 2))

			Expect(err).NotTo(HaveOccurred())
			Expect(facts).To(HaveLen(2))
		})
	})

	Describe("Count", func() {
		BeforeEach(func() {
			f1 := fixtures.FactWithCategories("", "Fact 1", "", []string{"leadership"}, []string{"a"})
			f1.RoleFit = career.RoleFitSeniorIC
			f2 := fixtures.FactWithCategories("", "Fact 2", "", []string{"technical"}, []string{"b"})
			f2.RoleFit = career.RoleFitSeniorIC
			for _, f := range []*career.Fact{f1, f2} {
				Expect(repo.Create(ctx, f)).To(Succeed())
			}
		})

		It("counts all facts", func() {
			count, err := repo.Count(ctx, *fixtures.FactListFilters())

			Expect(err).NotTo(HaveOccurred())
			Expect(count).To(Equal(2))
		})

		It("counts with filters", func() {
			count, err := repo.Count(ctx, *fixtures.FactListFiltersWithCategory("leadership"))

			Expect(err).NotTo(HaveOccurred())
			Expect(count).To(Equal(1))
		})
	})

	Describe("GetBySourceEventID", func() {
		It("returns facts for event", func() {
			fact1 := fixtures.Fact("", "event-1")
			fact1.Text = "F1"
			fact2 := fixtures.Fact("", "event-1")
			fact2.Text = "F2"
			fact3 := fixtures.Fact("", "event-2")
			fact3.Text = "F3"
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
			fact1 := fixtures.FactFromBurst("", "burst-1")
			fact1.Text = "F1"
			fact2 := fixtures.FactFromBurst("", "burst-2")
			fact2.Text = "F2"
			Expect(repo.Create(ctx, fact1)).To(Succeed())
			Expect(repo.Create(ctx, fact2)).To(Succeed())

			facts, err := repo.GetBySourceBurstID(ctx, "burst-1")

			Expect(err).NotTo(HaveOccurred())
			Expect(facts).To(HaveLen(1))
			Expect(facts[0].Text).To(Equal("F1"))
		})
	})

	Describe("Audience filter", func() {
		It("filters by audience relevance", func() {
			f1 := fixtures.FactWithCategories("", "Fact A", "", []string{"cat"}, []string{"technical"})
			f1.RoleFit = career.RoleFitSeniorIC
			f2 := fixtures.FactWithCategories("", "Fact B", "", []string{"cat"}, []string{"hr"})
			f2.RoleFit = career.RoleFitSeniorIC
			Expect(repo.Create(ctx, f1)).To(Succeed())
			Expect(repo.Create(ctx, f2)).To(Succeed())

			facts, err := repo.List(ctx, *fixtures.FactListFiltersWithAudience("technical"))

			Expect(err).NotTo(HaveOccurred())
			Expect(facts).To(HaveLen(1))
			Expect(facts[0].Text).To(Equal("Fact A"))
		})
	})

	Describe("Date range filter", func() {
		It("filters by start and end date", func() {
			now := time.Now()
			f1 := fixtures.FactWithCategories("", "Old fact", "", []string{"cat"}, []string{"aud"})
			f1.RoleFit = career.RoleFitSeniorIC
			Expect(repo.Create(ctx, f1)).To(Succeed())
			tx.Model(&models.Fact{}).Where("id = ?", f1.ID).Update("created_at", now.Add(-72*time.Hour))

			f2 := fixtures.FactWithCategories("", "New fact", "", []string{"cat"}, []string{"aud"})
			f2.RoleFit = career.RoleFitSeniorIC
			Expect(repo.Create(ctx, f2)).To(Succeed())

			start := now.Add(-24 * time.Hour)
			filters := fixtures.FactListFiltersWithDateRange(&start, nil)
			facts, err := repo.List(ctx, *filters)

			Expect(err).NotTo(HaveOccurred())
			Expect(facts).To(HaveLen(1))
			Expect(facts[0].Text).To(Equal("New fact"))
		})

		It("filters by end date only", func() {
			now := time.Now()
			f1 := fixtures.FactWithCategories("", "Old fact", "", []string{"cat"}, []string{"aud"})
			f1.RoleFit = career.RoleFitSeniorIC
			Expect(repo.Create(ctx, f1)).To(Succeed())
			tx.Model(&models.Fact{}).Where("id = ?", f1.ID).Update("created_at", now.Add(-72*time.Hour))

			f2 := fixtures.FactWithCategories("", "New fact", "", []string{"cat"}, []string{"aud"})
			f2.RoleFit = career.RoleFitSeniorIC
			Expect(repo.Create(ctx, f2)).To(Succeed())

			end := now.Add(-24 * time.Hour)
			filters := fixtures.FactListFiltersWithDateRange(nil, &end)
			facts, err := repo.List(ctx, *filters)

			Expect(err).NotTo(HaveOccurred())
			Expect(facts).To(HaveLen(1))
			Expect(facts[0].Text).To(Equal("Old fact"))
		})
	})

	Describe("Pagination with offset", func() {
		It("applies offset correctly", func() {
			for i := range 5 {
				f := fixtures.FactWithCategories("", "Fact "+string(rune('A'+i)), "", []string{"cat"}, []string{"aud"})
				f.RoleFit = career.RoleFitSeniorIC
				Expect(repo.Create(ctx, f)).To(Succeed())
				time.Sleep(5 * time.Millisecond)
			}

			facts, err := repo.List(ctx, *fixtures.FactListFiltersWithLimit(2, 2))

			Expect(err).NotTo(HaveOccurred())
			Expect(facts).To(HaveLen(2))
		})
	})

	Describe("Error handling with closed DB", func() {
		var closedRepo *FactRepository

		BeforeEach(func() {
			closedDB, err := stdsql.Open("sqlite", ":memory:")
			Expect(err).NotTo(HaveOccurred())
			gormDB, err := gorm.Open(sqlite.New(sqlite.Config{Conn: closedDB}), &gorm.Config{})
			Expect(err).NotTo(HaveOccurred())
			closedDB.Close()
			closedRepo = NewFactRepository(gormDB)
		})

		It("returns error on GetByID with closed DB", func() {
			_, err := closedRepo.GetByID(ctx, "id")
			Expect(err).To(HaveOccurred())
		})

		It("returns error on Update with closed DB", func() {
			f := fixtures.FactWithCategories("id", "Test", "", []string{"cat"}, []string{"aud"})
			f.RoleFit = career.RoleFitSeniorIC
			err := closedRepo.Update(ctx, f)
			Expect(err).To(HaveOccurred())
		})

		It("returns error on List with closed DB", func() {
			_, err := closedRepo.List(ctx, *fixtures.FactListFilters())
			Expect(err).To(HaveOccurred())
		})

		It("returns error on Count with closed DB", func() {
			_, err := closedRepo.Count(ctx, *fixtures.FactListFilters())
			Expect(err).To(HaveOccurred())
		})

		It("returns error on GetBySourceEventID with closed DB", func() {
			_, err := closedRepo.GetBySourceEventID(ctx, "e1")
			Expect(err).To(HaveOccurred())
		})

		It("returns error on GetBySourceBurstID with closed DB", func() {
			_, err := closedRepo.GetBySourceBurstID(ctx, "b1")
			Expect(err).To(HaveOccurred())
		})
	})
})
