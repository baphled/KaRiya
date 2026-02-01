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
	_ "modernc.org/sqlite"
)

func setupBurstTestDB() *gorm.DB {
	sqlDB, err := stdsql.Open("sqlite", ":memory:")
	Expect(err).NotTo(HaveOccurred())

	db, err := gorm.Open(sqlite.New(sqlite.Config{
		Conn: sqlDB,
	}), &gorm.Config{})
	Expect(err).NotTo(HaveOccurred())

	err = db.AutoMigrate(&models.Burst{})
	Expect(err).NotTo(HaveOccurred())

	return db
}

var _ = Describe("Burst Repository", func() {
	var (
		repo *BurstRepository
		db   *gorm.DB
		ctx  context.Context
	)

	BeforeEach(func() {
		db = setupBurstTestDB()
		repo = NewBurstRepository(db)
		ctx = context.Background()
	})

	Describe("Create", func() {
		It("creates a burst with generated ID", func() {
			burst := fixtures.Burst("", "event-1", "event-2")
			burst.Name = "Q4 Project Burst"
			burst.Description = "Events from Q4 project"

			err := repo.Create(ctx, burst)

			Expect(err).NotTo(HaveOccurred())
			Expect(burst.ID).NotTo(BeEmpty())
			Expect(burst.CreatedAt).NotTo(BeZero())
		})

		It("creates a burst with provided ID", func() {
			burst := fixtures.Burst("custom-id", "event-1")

			err := repo.Create(ctx, burst)

			Expect(err).NotTo(HaveOccurred())
			Expect(burst.ID).To(Equal("custom-id"))
		})
	})

	Describe("GetByID", func() {
		It("returns the burst", func() {
			burst := fixtures.Burst("", "e1", "e2")
			burst.Name = "Test Burst"
			burst.Description = "Test desc"
			Expect(repo.Create(ctx, burst)).To(Succeed())

			found, err := repo.GetByID(ctx, burst.ID)

			Expect(err).NotTo(HaveOccurred())
			Expect(found.Name).To(Equal("Test Burst"))
			Expect(found.EventIDs).To(ConsistOf("e1", "e2"))
		})

		It("returns ErrBurstNotFound for missing burst", func() {
			_, err := repo.GetByID(ctx, "nonexistent")

			Expect(err).To(Equal(career_repo.ErrBurstNotFound))
		})
	})

	Describe("Update", func() {
		It("updates the burst", func() {
			burst := fixtures.Burst("", "e1")
			burst.Name = "Original"
			Expect(repo.Create(ctx, burst)).To(Succeed())

			burst.Name = "Updated"
			burst.EventIDs = []string{"e1", "e2"}
			err := repo.Update(ctx, burst)

			Expect(err).NotTo(HaveOccurred())

			found, _ := repo.GetByID(ctx, burst.ID)
			Expect(found.Name).To(Equal("Updated"))
			Expect(found.EventIDs).To(ConsistOf("e1", "e2"))
		})

		It("returns ErrBurstNotFound for missing burst", func() {
			burst := fixtures.Burst("nonexistent", "e1")

			err := repo.Update(ctx, burst)

			Expect(err).To(Equal(career_repo.ErrBurstNotFound))
		})
	})

	Describe("Delete", func() {
		It("deletes the burst", func() {
			burst := fixtures.Burst("", "e1")
			burst.Name = "To delete"
			Expect(repo.Create(ctx, burst)).To(Succeed())

			err := repo.Delete(ctx, burst.ID)

			Expect(err).NotTo(HaveOccurred())

			_, err = repo.GetByID(ctx, burst.ID)
			Expect(err).To(Equal(career_repo.ErrBurstNotFound))
		})

		It("returns ErrBurstNotFound for missing burst", func() {
			err := repo.Delete(ctx, "nonexistent")

			Expect(err).To(Equal(career_repo.ErrBurstNotFound))
		})
	})

	Describe("List", func() {
		BeforeEach(func() {
			b1 := fixtures.Burst("", "e1")
			b1.Name = "Burst A"
			b2 := fixtures.Burst("", "e1", "e2")
			b2.Name = "Burst B"
			b3 := fixtures.Burst("", "e1", "e2", "e3")
			b3.Name = "Burst C"
			for _, b := range []*career.Burst{b1, b2, b3} {
				Expect(repo.Create(ctx, b)).To(Succeed())
				time.Sleep(10 * time.Millisecond) // Ensure different timestamps.
			}
		})

		It("returns all bursts without filters", func() {
			bursts, err := repo.List(ctx, *fixtures.BurstListFilters())

			Expect(err).NotTo(HaveOccurred())
			Expect(bursts).To(HaveLen(3))
		})

		It("sorts by created_at descending by default", func() {
			bursts, err := repo.List(ctx, *fixtures.BurstListFilters())

			Expect(err).NotTo(HaveOccurred())
			Expect(bursts[0].Name).To(Equal("Burst C")) // Most recent first.
		})

		It("sorts by name ascending", func() {
			bursts, err := repo.List(ctx, *fixtures.BurstListFiltersWithSort("name", "asc"))

			Expect(err).NotTo(HaveOccurred())
			Expect(bursts[0].Name).To(Equal("Burst A"))
		})

		It("applies pagination", func() {
			bursts, err := repo.List(ctx, *fixtures.BurstListFiltersWithLimit(0, 2))

			Expect(err).NotTo(HaveOccurred())
			Expect(bursts).To(HaveLen(2))
		})
	})

	Describe("Count", func() {
		BeforeEach(func() {
			b1 := fixtures.Burst("", "e1")
			b1.Name = "Burst 1"
			b2 := fixtures.Burst("", "e2")
			b2.Name = "Burst 2"
			for _, b := range []*career.Burst{b1, b2} {
				Expect(repo.Create(ctx, b)).To(Succeed())
			}
		})

		It("counts all bursts", func() {
			count, err := repo.Count(ctx, *fixtures.BurstListFilters())

			Expect(err).NotTo(HaveOccurred())
			Expect(count).To(Equal(2))
		})
	})

	Describe("Confirmed bursts", func() {
		It("stores and retrieves confirmed status", func() {
			burst := fixtures.BurstConfirmed("", "e1")
			burst.Name = "Confirmed Burst"
			Expect(repo.Create(ctx, burst)).To(Succeed())

			found, err := repo.GetByID(ctx, burst.ID)

			Expect(err).NotTo(HaveOccurred())
			Expect(found.Confirmed).To(BeTrue())
			Expect(found.ConfirmedAt).NotTo(BeNil())
		})
	})
})
