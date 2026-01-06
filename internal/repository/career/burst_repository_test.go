package career_test

import (
	"context"
	"time"

	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"

	"github.com/baphled/kariya/internal/domain/career"
	repo "github.com/baphled/kariya/internal/repository/career"
)

var _ = Describe("BurstRepository", func() {
	var (
		repository repo.BurstRepository
		ctx        context.Context
	)

	BeforeEach(func() {
		repository = repo.NewMemoryBurstRepository()
		ctx = context.Background()
	})

	Describe("Create", func() {
		It("should create a new burst", func() {
			burst := &career.Burst{
				Name:        "Platform Migration",
				Description: "Migrated entire platform to microservices",
				EventIDs:    []string{"event-1", "event-2"},
			}

			err := repository.Create(ctx, burst)
			Expect(err).ToNot(HaveOccurred())
			Expect(burst.ID).ToNot(BeEmpty())
			Expect(burst.CreatedAt).ToNot(BeZero())
			Expect(burst.UpdatedAt).ToNot(BeZero())
		})

		It("should return error for duplicate burst ID", func() {
			burst := &career.Burst{
				ID:          "burst-1",
				Name:        "Platform Migration",
				Description: "Migrated entire platform to microservices",
				EventIDs:    []string{"event-1", "event-2"},
			}

			err := repository.Create(ctx, burst)
			Expect(err).ToNot(HaveOccurred())

			// Try to create again with same ID
			duplicate := &career.Burst{
				ID:          "burst-1",
				Name:        "Different Name",
				Description: "Different description",
				EventIDs:    []string{"event-3", "event-4"},
			}

			err = repository.Create(ctx, duplicate)
			Expect(err).To(MatchError(repo.ErrDuplicateBurst))
		})

		It("should generate unique ID if not provided", func() {
			burst1 := &career.Burst{
				Name:     "Burst 1",
				EventIDs: []string{"event-1", "event-2"},
			}

			burst2 := &career.Burst{
				Name:     "Burst 2",
				EventIDs: []string{"event-3", "event-4"},
			}

			err1 := repository.Create(ctx, burst1)
			err2 := repository.Create(ctx, burst2)

			Expect(err1).ToNot(HaveOccurred())
			Expect(err2).ToNot(HaveOccurred())
			Expect(burst1.ID).ToNot(BeEmpty())
			Expect(burst2.ID).ToNot(BeEmpty())
			Expect(burst1.ID).ToNot(Equal(burst2.ID))
		})
	})

	Describe("GetByID", func() {
		It("should retrieve an existing burst", func() {
			burst := &career.Burst{
				Name:        "Platform Migration",
				Description: "Migrated entire platform to microservices",
				EventIDs:    []string{"event-1", "event-2"},
			}

			err := repository.Create(ctx, burst)
			Expect(err).ToNot(HaveOccurred())

			retrieved, err := repository.GetByID(ctx, burst.ID)
			Expect(err).ToNot(HaveOccurred())
			Expect(retrieved.ID).To(Equal(burst.ID))
			Expect(retrieved.Name).To(Equal(burst.Name))
			Expect(retrieved.Description).To(Equal(burst.Description))
			Expect(retrieved.EventIDs).To(Equal(burst.EventIDs))

		})

		It("should return error for non-existent burst", func() {
			_, err := repository.GetByID(ctx, "non-existent-id")
			Expect(err).To(MatchError(repo.ErrBurstNotFound))
		})
	})

	Describe("Update", func() {
		It("should update an existing burst", func() {
			burst := &career.Burst{
				Name:        "Platform Migration",
				Description: "Initial description",
				EventIDs:    []string{"event-1", "event-2"},
			}

			err := repository.Create(ctx, burst)
			Expect(err).ToNot(HaveOccurred())

			originalUpdatedAt := burst.UpdatedAt
			time.Sleep(10 * time.Millisecond) // Ensure timestamp difference

			burst.Name = "Updated Platform Migration"
			burst.Description = "Updated description"
			burst.EventIDs = []string{"event-1", "event-2", "event-3"}

			err = repository.Update(ctx, burst)
			Expect(err).ToNot(HaveOccurred())
			Expect(burst.UpdatedAt.After(originalUpdatedAt)).To(BeTrue())

			retrieved, err := repository.GetByID(ctx, burst.ID)
			Expect(err).ToNot(HaveOccurred())
			Expect(retrieved.Name).To(Equal("Updated Platform Migration"))
			Expect(retrieved.Description).To(Equal("Updated description"))
			Expect(retrieved.EventIDs).To(HaveLen(3))
		})

		It("should return error for non-existent burst", func() {
			burst := &career.Burst{
				ID:       "non-existent-id",
				Name:     "Platform Migration",
				EventIDs: []string{"event-1", "event-2"},
			}

			err := repository.Update(ctx, burst)
			Expect(err).To(MatchError(repo.ErrBurstNotFound))
		})
	})

	Describe("Delete", func() {
		It("should delete an existing burst", func() {
			burst := &career.Burst{
				Name:        "Platform Migration",
				Description: "Migrated entire platform to microservices",
				EventIDs:    []string{"event-1", "event-2"},
			}

			err := repository.Create(ctx, burst)
			Expect(err).ToNot(HaveOccurred())

			err = repository.Delete(ctx, burst.ID)
			Expect(err).ToNot(HaveOccurred())

			_, err = repository.GetByID(ctx, burst.ID)
			Expect(err).To(MatchError(repo.ErrBurstNotFound))
		})

		It("should return error for non-existent burst", func() {
			err := repository.Delete(ctx, "non-existent-id")
			Expect(err).To(MatchError(repo.ErrBurstNotFound))
		})
	})

	Describe("List", func() {
		BeforeEach(func() {
			// Create test bursts
			bursts := []*career.Burst{
				{
					Name:        "Platform Migration",
					Description: "Migrated to microservices",
					EventIDs:    []string{"event-1", "event-2"},
				},
				{
					Name:        "Team Leadership",
					Description: "Led cross-functional team",
					EventIDs:    []string{"event-3", "event-4", "event-5"},
				},
				{
					Name:        "Product Launch",
					Description: "Launched new product feature",
					EventIDs:    []string{"event-6", "event-7"},
				},
			}

			for _, burst := range bursts {
				err := repository.Create(ctx, burst)
				Expect(err).ToNot(HaveOccurred())
			}
		})

		It("should list all bursts with no filters", func() {
			bursts, err := repository.List(ctx, repo.BurstListFilters{})
			Expect(err).ToNot(HaveOccurred())
			Expect(bursts).To(HaveLen(3))
		})

		It("should sort bursts by name ascending", func() {
			filters := repo.BurstListFilters{
				SortBy:    "name",
				SortOrder: "asc",
			}

			bursts, err := repository.List(ctx, filters)
			Expect(err).ToNot(HaveOccurred())
			Expect(bursts).To(HaveLen(3))
			Expect(bursts[0].Name).To(Equal("Platform Migration"))
			Expect(bursts[1].Name).To(Equal("Product Launch"))
			Expect(bursts[2].Name).To(Equal("Team Leadership"))
		})

		It("should sort bursts by name descending", func() {
			filters := repo.BurstListFilters{
				SortBy:    "name",
				SortOrder: "desc",
			}

			bursts, err := repository.List(ctx, filters)
			Expect(err).ToNot(HaveOccurred())
			Expect(bursts).To(HaveLen(3))
			Expect(bursts[0].Name).To(Equal("Team Leadership"))
			Expect(bursts[1].Name).To(Equal("Product Launch"))
			Expect(bursts[2].Name).To(Equal("Platform Migration"))
		})

		It("should sort bursts by event count descending", func() {
			filters := repo.BurstListFilters{
				SortBy:    "event_count",
				SortOrder: "desc",
			}

			bursts, err := repository.List(ctx, filters)
			Expect(err).ToNot(HaveOccurred())
			Expect(bursts).To(HaveLen(3))
			Expect(bursts[0].EventIDs).To(HaveLen(3)) // Team Leadership
			Expect(bursts[1].EventIDs).To(HaveLen(2)) // Platform Migration or Product Launch
			Expect(bursts[2].EventIDs).To(HaveLen(2)) // Platform Migration or Product Launch
		})

		It("should apply pagination with limit", func() {
			filters := repo.BurstListFilters{
				Limit:     2,
				SortBy:    "name",
				SortOrder: "asc",
			}

			bursts, err := repository.List(ctx, filters)
			Expect(err).ToNot(HaveOccurred())
			Expect(bursts).To(HaveLen(2))
		})

		It("should apply pagination with offset", func() {
			filters := repo.BurstListFilters{
				Offset:    1,
				SortBy:    "name",
				SortOrder: "asc",
			}

			bursts, err := repository.List(ctx, filters)
			Expect(err).ToNot(HaveOccurred())
			Expect(bursts).To(HaveLen(2))
			Expect(bursts[0].Name).To(Equal("Product Launch"))
		})

		It("should apply pagination with both limit and offset", func() {
			filters := repo.BurstListFilters{
				Offset:    1,
				Limit:     1,
				SortBy:    "name",
				SortOrder: "asc",
			}

			bursts, err := repository.List(ctx, filters)
			Expect(err).ToNot(HaveOccurred())
			Expect(bursts).To(HaveLen(1))
			Expect(bursts[0].Name).To(Equal("Product Launch"))
		})
	})

	Describe("Count", func() {
		BeforeEach(func() {
			// Create test bursts
			bursts := []*career.Burst{
				{
					Name:     "Platform Migration",
					EventIDs: []string{"event-1", "event-2"},
				},
				{
					Name:     "Team Leadership",
					EventIDs: []string{"event-3", "event-4"},
				},
				{
					Name:     "API Design",
					EventIDs: []string{"event-5", "event-6"},
				},
			}

			for _, burst := range bursts {
				err := repository.Create(ctx, burst)
				Expect(err).ToNot(HaveOccurred())
			}
		})

		It("should count all bursts with no filters", func() {
			count, err := repository.Count(ctx, repo.BurstListFilters{})
			Expect(err).ToNot(HaveOccurred())
			Expect(count).To(Equal(3))
		})

	})
})
