package career

import (
	"context"
	"database/sql"
	"testing"

	"github.com/baphled/kariya/internal/domain/career"
	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
	_ "modernc.org/sqlite"
)

func TestSQLiteBurstRepository(t *testing.T) {
	RegisterFailHandler(Fail)
	RunSpecs(t, "SQLite Burst Repository Suite")
}

var _ = Describe("SQLiteBurstRepository", func() {
	var (
		db   *sql.DB
		repo *SQLiteBurstRepository
		ctx  context.Context
	)

	BeforeEach(func() {
		var err error
		db, err = sql.Open("sqlite", ":memory:")
		Expect(err).NotTo(HaveOccurred())

		repo, err = NewSQLiteBurstRepository(db)
		Expect(err).NotTo(HaveOccurred())

		ctx = context.Background()
	})

	AfterEach(func() {
		db.Close()
	})

	Context("when creating a burst", func() {
		It("should create a burst successfully", func() {
			burst := &career.Burst{
				Name:     "Platform Migration",
				EventIDs: []string{"event-1", "event-2"},
			}

			err := repo.Create(ctx, burst)
			Expect(err).NotTo(HaveOccurred())
			Expect(burst.ID).NotTo(BeEmpty())
		})

		It("should reject duplicate burst IDs", func() {
			burst1 := &career.Burst{
				ID:       "burst-123",
				Name:     "Burst 1",
				EventIDs: []string{"event-1", "event-2"},
			}
			burst2 := &career.Burst{
				ID:       "burst-123",
				Name:     "Burst 2",
				EventIDs: []string{"event-3", "event-4"},
			}

			err1 := repo.Create(ctx, burst1)
			err2 := repo.Create(ctx, burst2)

			Expect(err1).NotTo(HaveOccurred())
			Expect(err2).To(HaveOccurred())
		})
	})

	Context("when retrieving a burst", func() {
		var burstID string

		BeforeEach(func() {
			burst := &career.Burst{
				Name:     "Test Burst",
				EventIDs: []string{"event-1", "event-2"},
			}
			repo.Create(ctx, burst)
			burstID = burst.ID
		})

		It("should retrieve a burst by ID", func() {
			burst, err := repo.GetByID(ctx, burstID)
			Expect(err).NotTo(HaveOccurred())
			Expect(burst.ID).To(Equal(burstID))
			Expect(burst.Name).To(Equal("Test Burst"))
			Expect(burst.EventIDs).To(HaveLen(2))
		})

		It("should return error for non-existent burst", func() {
			_, err := repo.GetByID(ctx, "non-existent")
			Expect(err).To(HaveOccurred())
		})
	})

	Context("when updating a burst", func() {
		var burstID string

		BeforeEach(func() {
			burst := &career.Burst{
				Name:     "Original",
				EventIDs: []string{"event-1", "event-2"},
			}
			repo.Create(ctx, burst)
			burstID = burst.ID
		})

		It("should update a burst successfully", func() {
			burst, _ := repo.GetByID(ctx, burstID)
			burst.Name = "Updated"
			burst.Description = "New description"

			err := repo.Update(ctx, burst)
			Expect(err).NotTo(HaveOccurred())

			updated, _ := repo.GetByID(ctx, burstID)
			Expect(updated.Name).To(Equal("Updated"))
			Expect(updated.Description).To(Equal("New description"))
		})
	})

	Context("when deleting a burst", func() {
		var burstID string

		BeforeEach(func() {
			burst := &career.Burst{
				Name:     "To Delete",
				EventIDs: []string{"event-1", "event-2"},
			}
			repo.Create(ctx, burst)
			burstID = burst.ID
		})

		It("should delete a burst successfully", func() {
			err := repo.Delete(ctx, burstID)
			Expect(err).NotTo(HaveOccurred())

			_, err = repo.GetByID(ctx, burstID)
			Expect(err).To(HaveOccurred())
		})
	})

	Context("when listing bursts", func() {
		BeforeEach(func() {
			for i := 1; i <= 3; i++ {
				burst := &career.Burst{
					Name:            "Burst " + string(rune(48+i)),
					EventIDs:        []string{"event-1", "event-2"},
					CompetencyFocus: "technical",
				}
				repo.Create(ctx, burst)
			}
		})

		It("should list all bursts", func() {
			bursts, err := repo.List(ctx, BurstListFilters{})
			Expect(err).NotTo(HaveOccurred())
			Expect(bursts).To(HaveLen(3))
		})

		It("should filter by competency focus", func() {
			burst := &career.Burst{
				Name:            "Leadership",
				EventIDs:        []string{"event-1", "event-2"},
				CompetencyFocus: "leadership",
			}
			repo.Create(ctx, burst)

			bursts, _ := repo.List(ctx, BurstListFilters{
				CompetencyFocus: "technical",
			})
			Expect(bursts).To(HaveLen(3))

			bursts, _ = repo.List(ctx, BurstListFilters{
				CompetencyFocus: "leadership",
			})
			Expect(bursts).To(HaveLen(1))
		})

		It("should apply pagination", func() {
			bursts, _ := repo.List(ctx, BurstListFilters{
				Limit:  1,
				Offset: 0,
			})
			Expect(bursts).To(HaveLen(1))

			bursts, _ = repo.List(ctx, BurstListFilters{
				Limit:  1,
				Offset: 1,
			})
			Expect(bursts).To(HaveLen(1))
		})
	})

	Context("when counting bursts", func() {
		BeforeEach(func() {
			for i := 1; i <= 3; i++ {
				burst := &career.Burst{
					Name:     "Burst " + string(rune(48+i)),
					EventIDs: []string{"event-1", "event-2"},
				}
				repo.Create(ctx, burst)
			}
		})

		It("should count all bursts", func() {
			count, err := repo.Count(ctx, BurstListFilters{})
			Expect(err).NotTo(HaveOccurred())
			Expect(count).To(Equal(3))
		})
	})
})

