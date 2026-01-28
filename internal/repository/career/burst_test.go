package career

import (
	"context"
	"database/sql"
	"time"

	"github.com/baphled/kariya/internal/domain/career"
	"github.com/baphled/kariya/internal/repository/models"
	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
	_ "modernc.org/sqlite"
)

func setupBurstTestDB() *gorm.DB {
	sqlDB, err := sql.Open("sqlite", ":memory:")
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
		repo *Burst
		db   *gorm.DB
		ctx  context.Context
	)

	BeforeEach(func() {
		db = setupBurstTestDB()
		repo = NewBurst(db)
		ctx = context.Background()
	})

	Describe("Create", func() {
		It("creates a burst with generated ID", func() {
			burst := &career.Burst{
				Name:        "Q4 Project Burst",
				Description: "Events from Q4 project",
				EventIDs:    []string{"event-1", "event-2"},
			}

			err := repo.Create(ctx, burst)

			Expect(err).NotTo(HaveOccurred())
			Expect(burst.ID).NotTo(BeEmpty())
			Expect(burst.CreatedAt).NotTo(BeZero())
		})

		It("creates a burst with provided ID", func() {
			burst := &career.Burst{
				ID:       "custom-id",
				Name:     "Test Burst",
				EventIDs: []string{"event-1"},
			}

			err := repo.Create(ctx, burst)

			Expect(err).NotTo(HaveOccurred())
			Expect(burst.ID).To(Equal("custom-id"))
		})
	})

	Describe("GetByID", func() {
		It("returns the burst", func() {
			burst := &career.Burst{
				Name:        "Test Burst",
				Description: "Test desc",
				EventIDs:    []string{"e1", "e2"},
			}
			Expect(repo.Create(ctx, burst)).To(Succeed())

			found, err := repo.GetByID(ctx, burst.ID)

			Expect(err).NotTo(HaveOccurred())
			Expect(found.Name).To(Equal("Test Burst"))
			Expect(found.EventIDs).To(ConsistOf("e1", "e2"))
		})

		It("returns ErrBurstNotFound for missing burst", func() {
			_, err := repo.GetByID(ctx, "nonexistent")

			Expect(err).To(Equal(ErrBurstNotFound))
		})
	})

	Describe("Update", func() {
		It("updates the burst", func() {
			burst := &career.Burst{
				Name:     "Original",
				EventIDs: []string{"e1"},
			}
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
			burst := &career.Burst{
				ID:       "nonexistent",
				Name:     "Test",
				EventIDs: []string{"e1"},
			}

			err := repo.Update(ctx, burst)

			Expect(err).To(Equal(ErrBurstNotFound))
		})
	})

	Describe("Delete", func() {
		It("deletes the burst", func() {
			burst := &career.Burst{
				Name:     "To delete",
				EventIDs: []string{"e1"},
			}
			Expect(repo.Create(ctx, burst)).To(Succeed())

			err := repo.Delete(ctx, burst.ID)

			Expect(err).NotTo(HaveOccurred())

			_, err = repo.GetByID(ctx, burst.ID)
			Expect(err).To(Equal(ErrBurstNotFound))
		})

		It("returns ErrBurstNotFound for missing burst", func() {
			err := repo.Delete(ctx, "nonexistent")

			Expect(err).To(Equal(ErrBurstNotFound))
		})
	})

	Describe("List", func() {
		BeforeEach(func() {
			bursts := []*career.Burst{
				{Name: "Burst A", EventIDs: []string{"e1"}},
				{Name: "Burst B", EventIDs: []string{"e1", "e2"}},
				{Name: "Burst C", EventIDs: []string{"e1", "e2", "e3"}},
			}
			for _, b := range bursts {
				Expect(repo.Create(ctx, b)).To(Succeed())
				time.Sleep(10 * time.Millisecond) // Ensure different timestamps.
			}
		})

		It("returns all bursts without filters", func() {
			bursts, err := repo.List(ctx, BurstListFilters{})

			Expect(err).NotTo(HaveOccurred())
			Expect(bursts).To(HaveLen(3))
		})

		It("sorts by created_at descending by default", func() {
			bursts, err := repo.List(ctx, BurstListFilters{})

			Expect(err).NotTo(HaveOccurred())
			Expect(bursts[0].Name).To(Equal("Burst C")) // Most recent first.
		})

		It("sorts by name ascending", func() {
			bursts, err := repo.List(ctx, BurstListFilters{SortBy: "name", SortOrder: "asc"})

			Expect(err).NotTo(HaveOccurred())
			Expect(bursts[0].Name).To(Equal("Burst A"))
		})

		It("applies pagination", func() {
			bursts, err := repo.List(ctx, BurstListFilters{Limit: 2})

			Expect(err).NotTo(HaveOccurred())
			Expect(bursts).To(HaveLen(2))
		})
	})

	Describe("Count", func() {
		BeforeEach(func() {
			bursts := []*career.Burst{
				{Name: "Burst 1", EventIDs: []string{"e1"}},
				{Name: "Burst 2", EventIDs: []string{"e2"}},
			}
			for _, b := range bursts {
				Expect(repo.Create(ctx, b)).To(Succeed())
			}
		})

		It("counts all bursts", func() {
			count, err := repo.Count(ctx, BurstListFilters{})

			Expect(err).NotTo(HaveOccurred())
			Expect(count).To(Equal(2))
		})
	})

	Describe("Confirmed bursts", func() {
		It("stores and retrieves confirmed status", func() {
			now := time.Now()
			burst := &career.Burst{
				Name:        "Confirmed Burst",
				EventIDs:    []string{"e1"},
				Confirmed:   true,
				ConfirmedAt: &now,
			}
			Expect(repo.Create(ctx, burst)).To(Succeed())

			found, err := repo.GetByID(ctx, burst.ID)

			Expect(err).NotTo(HaveOccurred())
			Expect(found.Confirmed).To(BeTrue())
			Expect(found.ConfirmedAt).NotTo(BeNil())
		})
	})
})
