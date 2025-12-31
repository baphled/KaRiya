package career

import (
	"context"
	"testing"
	"time"

	"github.com/baphled/kariya/internal/domain/career"
	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
)

func TestBurstRepository(t *testing.T) {
	RegisterFailHandler(Fail)
	RunSpecs(t, "Burst Repository Suite")
}

var _ = Describe("MemoryBurstRepository", func() {
	var (
		repo *MemoryBurstRepository
		ctx  context.Context
	)

	BeforeEach(func() {
		repo = NewMemoryBurstRepository()
		ctx = context.Background()
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
			Expect(burst.CreatedAt).NotTo(BeZero())
			Expect(burst.UpdatedAt).NotTo(BeZero())
		})

		It("should generate a unique ID if not provided", func() {
			burst1 := &career.Burst{
				Name:     "Burst 1",
				EventIDs: []string{"event-1", "event-2"},
			}
			burst2 := &career.Burst{
				Name:     "Burst 2",
				EventIDs: []string{"event-3", "event-4"},
			}

			err1 := repo.Create(ctx, burst1)
			err2 := repo.Create(ctx, burst2)

			Expect(err1).NotTo(HaveOccurred())
			Expect(err2).NotTo(HaveOccurred())
			Expect(burst1.ID).NotTo(Equal(burst2.ID))
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
			Expect(err2).To(Equal(ErrDuplicateBurst))
		})

		It("should reject invalid bursts", func() {
			burst := &career.Burst{
				ID:       "burst-123",
				Name:     "", // Empty name is invalid
				EventIDs: []string{"event-1", "event-2"},
			}

			err := repo.Create(ctx, burst)
			Expect(err).To(HaveOccurred())
		})
	})

	Context("when retrieving a burst", func() {
		var burstID string

		BeforeEach(func() {
			burst := &career.Burst{
				Name:     "Platform Migration",
				EventIDs: []string{"event-1", "event-2"},
			}
			repo.Create(ctx, burst)
			burstID = burst.ID
		})

		It("should retrieve a burst by ID", func() {
			burst, err := repo.GetByID(ctx, burstID)
			Expect(err).NotTo(HaveOccurred())
			Expect(burst).NotTo(BeNil())
			Expect(burst.ID).To(Equal(burstID))
			Expect(burst.Name).To(Equal("Platform Migration"))
		})

		It("should return error for non-existent burst", func() {
			burst, err := repo.GetByID(ctx, "non-existent-id")
			Expect(err).To(HaveOccurred())
			Expect(err).To(Equal(ErrBurstNotFound))
			Expect(burst).To(BeNil())
		})
	})

	Context("when updating a burst", func() {
		var burstID string

		BeforeEach(func() {
			burst := &career.Burst{
				Name:     "Original Name",
				EventIDs: []string{"event-1", "event-2"},
			}
			repo.Create(ctx, burst)
			burstID = burst.ID
		})

		It("should update a burst successfully", func() {
			burst, _ := repo.GetByID(ctx, burstID)
			originalUpdatedAt := burst.UpdatedAt

			// Simulate a small delay to ensure UpdatedAt changes
			time.Sleep(10 * time.Millisecond)

			burst.Name = "Updated Name"
			burst.Description = "New description"

			err := repo.Update(ctx, burst)
			Expect(err).NotTo(HaveOccurred())

			updated, _ := repo.GetByID(ctx, burstID)
			Expect(updated.Name).To(Equal("Updated Name"))
			Expect(updated.Description).To(Equal("New description"))
			Expect(updated.UpdatedAt.After(originalUpdatedAt)).To(BeTrue())
		})

		It("should reject update of non-existent burst", func() {
			burst := &career.Burst{
				ID:       "non-existent-id",
				Name:     "Test",
				EventIDs: []string{"event-1", "event-2"},
			}

			err := repo.Update(ctx, burst)
			Expect(err).To(HaveOccurred())
			Expect(err).To(Equal(ErrBurstNotFound))
		})

		It("should reject invalid burst updates", func() {
			burst, _ := repo.GetByID(ctx, burstID)
			burst.Name = "" // Make it invalid

			err := repo.Update(ctx, burst)
			Expect(err).To(HaveOccurred())
		})
	})

	Context("when deleting a burst", func() {
		var burstID string

		BeforeEach(func() {
			burst := &career.Burst{
				Name:     "Platform Migration",
				EventIDs: []string{"event-1", "event-2"},
			}
			repo.Create(ctx, burst)
			burstID = burst.ID
		})

		It("should delete a burst successfully", func() {
			err := repo.Delete(ctx, burstID)
			Expect(err).NotTo(HaveOccurred())

			_, err = repo.GetByID(ctx, burstID)
			Expect(err).To(Equal(ErrBurstNotFound))
		})

		It("should reject deletion of non-existent burst", func() {
			err := repo.Delete(ctx, "non-existent-id")
			Expect(err).To(HaveOccurred())
			Expect(err).To(Equal(ErrBurstNotFound))
		})
	})

	Context("when listing bursts", func() {
		BeforeEach(func() {
			for i := 1; i <= 5; i++ {
				burst := &career.Burst{
					Name:              "Burst " + string(rune(48+i)),
					EventIDs:          []string{"event-1", "event-2"},
					CompetencyFocus:   "technical",
					Description:       "Test burst",
				}
				repo.Create(ctx, burst)
			}
		})

		It("should list all bursts", func() {
			bursts, err := repo.List(ctx, BurstListFilters{})
			Expect(err).NotTo(HaveOccurred())
			Expect(bursts).To(HaveLen(5))
		})

		It("should filter bursts by competency focus", func() {
			// Create a burst with different competency
			burst := &career.Burst{
				Name:            "Leadership Burst",
				EventIDs:        []string{"event-1", "event-2"},
				CompetencyFocus: "leadership",
			}
			repo.Create(ctx, burst)

			bursts, err := repo.List(ctx, BurstListFilters{
				CompetencyFocus: "technical",
			})
			Expect(err).NotTo(HaveOccurred())
			Expect(bursts).To(HaveLen(5))

			bursts, err = repo.List(ctx, BurstListFilters{
				CompetencyFocus: "leadership",
			})
			Expect(err).NotTo(HaveOccurred())
			Expect(bursts).To(HaveLen(1))
		})

		It("should sort bursts by name ascending", func() {
			bursts, err := repo.List(ctx, BurstListFilters{
				SortBy:    "name",
				SortOrder: "asc",
			})
			Expect(err).NotTo(HaveOccurred())
			Expect(bursts).To(HaveLen(5))

			for i := 0; i < len(bursts)-1; i++ {
				Expect(bursts[i].Name <= bursts[i+1].Name).To(BeTrue())
			}
		})

		It("should sort bursts by name descending", func() {
			bursts, err := repo.List(ctx, BurstListFilters{
				SortBy:    "name",
				SortOrder: "desc",
			})
			Expect(err).NotTo(HaveOccurred())
			Expect(bursts).To(HaveLen(5))

			for i := 0; i < len(bursts)-1; i++ {
				Expect(bursts[i].Name >= bursts[i+1].Name).To(BeTrue())
			}
		})

		It("should sort bursts by created_at by default", func() {
			bursts, err := repo.List(ctx, BurstListFilters{})
			Expect(err).NotTo(HaveOccurred())
			Expect(bursts).To(HaveLen(5))

			// Default is descending by created_at (most recent first)
			for i := 0; i < len(bursts)-1; i++ {
				Expect(bursts[i].CreatedAt.After(bursts[i+1].CreatedAt) || bursts[i].CreatedAt.Equal(bursts[i+1].CreatedAt)).To(BeTrue())
			}
		})

		It("should apply pagination with limit and offset", func() {
			bursts, err := repo.List(ctx, BurstListFilters{
				Limit:  2,
				Offset: 0,
			})
			Expect(err).NotTo(HaveOccurred())
			Expect(bursts).To(HaveLen(2))

			bursts, err = repo.List(ctx, BurstListFilters{
				Limit:  2,
				Offset: 2,
			})
			Expect(err).NotTo(HaveOccurred())
			Expect(bursts).To(HaveLen(2))

			bursts, err = repo.List(ctx, BurstListFilters{
				Limit:  2,
				Offset: 4,
			})
			Expect(err).NotTo(HaveOccurred())
			Expect(bursts).To(HaveLen(1))
		})

		It("should return empty list for offset beyond total count", func() {
			bursts, err := repo.List(ctx, BurstListFilters{
				Limit:  10,
				Offset: 100,
			})
			Expect(err).NotTo(HaveOccurred())
			Expect(bursts).To(HaveLen(0))
		})

		It("should apply date range filter", func() {
			now := time.Now()
			startDate := now.Add(-1 * time.Hour)
			endDate := now.Add(1 * time.Hour)

			bursts, err := repo.List(ctx, BurstListFilters{
				StartDate: &startDate,
				EndDate:   &endDate,
			})
			Expect(err).NotTo(HaveOccurred())
			Expect(bursts).To(HaveLen(5))

			// Filter for bursts created in the past (before now)
			futureDate := now.Add(1 * time.Hour)
			bursts, err = repo.List(ctx, BurstListFilters{
				EndDate: &futureDate,
			})
			Expect(err).NotTo(HaveOccurred())
			Expect(bursts).To(HaveLen(5))
		})
	})

	Context("when counting bursts", func() {
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

		It("should count all bursts", func() {
			count, err := repo.Count(ctx, BurstListFilters{})
			Expect(err).NotTo(HaveOccurred())
			Expect(count).To(Equal(3))
		})

		It("should count bursts with filter", func() {
			burst := &career.Burst{
				Name:            "Leadership Burst",
				EventIDs:        []string{"event-1", "event-2"},
				CompetencyFocus: "leadership",
			}
			repo.Create(ctx, burst)

			count, err := repo.Count(ctx, BurstListFilters{
				CompetencyFocus: "technical",
			})
			Expect(err).NotTo(HaveOccurred())
			Expect(count).To(Equal(3))

			count, err = repo.Count(ctx, BurstListFilters{
				CompetencyFocus: "leadership",
			})
			Expect(err).NotTo(HaveOccurred())
			Expect(count).To(Equal(1))
		})
	})
})

