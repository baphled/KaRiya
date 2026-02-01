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

var _ = Describe("BurstRepository", func() {
	var (
		repository career_repo.BurstRepository
		ctx        context.Context
	)

	BeforeEach(func() {
		repository = NewBurstRepository()
		ctx = context.Background()
	})

	Describe("Create", func() {
		It("should create a new burst", func() {
			burst := fixtures.BurstFactory.MustCreate().(*career.Burst)
			burst.ID = "" // Clear ID to test auto-generation

			err := repository.Create(ctx, burst)
			Expect(err).ToNot(HaveOccurred())
			Expect(burst.ID).ToNot(BeEmpty())
			Expect(burst.CreatedAt).ToNot(BeZero())
			Expect(burst.UpdatedAt).ToNot(BeZero())
		})

		It("should return error for duplicate burst ID", func() {
			burst := fixtures.Burst("burst-1", "event-1", "event-2")

			err := repository.Create(ctx, burst)
			Expect(err).ToNot(HaveOccurred())

			// Try to create again with same ID
			duplicate := fixtures.Burst("burst-1", "event-3", "event-4")

			err = repository.Create(ctx, duplicate)
			Expect(err).To(MatchError(career_repo.ErrDuplicateBurst))
		})

		It("should generate unique ID if not provided", func() {
			burst1 := fixtures.BurstFactory.MustCreate().(*career.Burst)
			burst1.ID = "" // Clear ID to test auto-generation

			burst2 := fixtures.BurstFactory.MustCreate().(*career.Burst)
			burst2.ID = "" // Clear ID to test auto-generation

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
			burst := fixtures.BurstFactory.MustCreate().(*career.Burst)
			burst.ID = "" // Clear to test auto-generation

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
			Expect(err).To(MatchError(career_repo.ErrBurstNotFound))
		})
	})

	Describe("Update", func() {
		It("should update an existing burst", func() {
			burst := fixtures.Burst("", "event-1", "event-2")

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
			burst := fixtures.Burst("non-existent-id", "event-1", "event-2")

			err := repository.Update(ctx, burst)
			Expect(err).To(MatchError(career_repo.ErrBurstNotFound))
		})
	})

	Describe("Delete", func() {
		It("should delete an existing burst", func() {
			burst := fixtures.BurstFactory.MustCreate().(*career.Burst)
			burst.ID = "" // Clear to test auto-generation

			err := repository.Create(ctx, burst)
			Expect(err).ToNot(HaveOccurred())

			err = repository.Delete(ctx, burst.ID)
			Expect(err).ToNot(HaveOccurred())

			_, err = repository.GetByID(ctx, burst.ID)
			Expect(err).To(MatchError(career_repo.ErrBurstNotFound))
		})

		It("should return error for non-existent burst", func() {
			err := repository.Delete(ctx, "non-existent-id")
			Expect(err).To(MatchError(career_repo.ErrBurstNotFound))
		})
	})

	Describe("List", func() {
		BeforeEach(func() {
			// Create test bursts with specific names for sorting tests
			burst1 := fixtures.Burst("", "event-1", "event-2")
			burst1.Name = "Platform Migration"
			burst1.Description = "Migrated to microservices"

			burst2 := fixtures.Burst("", "event-3", "event-4", "event-5")
			burst2.Name = "Team Leadership"
			burst2.Description = "Led cross-functional team"

			burst3 := fixtures.Burst("", "event-6", "event-7")
			burst3.Name = "Product Launch"
			burst3.Description = "Launched new product feature"

			for _, burst := range []*career.Burst{burst1, burst2, burst3} {
				err := repository.Create(ctx, burst)
				Expect(err).ToNot(HaveOccurred())
			}
		})

		It("should list all bursts with no filters", func() {
			bursts, err := repository.List(ctx, *fixtures.BurstListFilters())
			Expect(err).ToNot(HaveOccurred())
			Expect(bursts).To(HaveLen(3))
		})

		It("should sort bursts by name ascending", func() {
			filters := fixtures.BurstListFiltersWithSort("name", "asc")

			bursts, err := repository.List(ctx, *filters)
			Expect(err).ToNot(HaveOccurred())
			Expect(bursts).To(HaveLen(3))
			Expect(bursts[0].Name).To(Equal("Platform Migration"))
			Expect(bursts[1].Name).To(Equal("Product Launch"))
			Expect(bursts[2].Name).To(Equal("Team Leadership"))
		})

		It("should sort bursts by name descending", func() {
			filters := fixtures.BurstListFiltersWithSort("name", "desc")

			bursts, err := repository.List(ctx, *filters)
			Expect(err).ToNot(HaveOccurred())
			Expect(bursts).To(HaveLen(3))
			Expect(bursts[0].Name).To(Equal("Team Leadership"))
			Expect(bursts[1].Name).To(Equal("Product Launch"))
			Expect(bursts[2].Name).To(Equal("Platform Migration"))
		})

		It("should sort bursts by event count descending", func() {
			filters := fixtures.BurstListFiltersWithSort("event_count", "desc")

			bursts, err := repository.List(ctx, *filters)
			Expect(err).ToNot(HaveOccurred())
			Expect(bursts).To(HaveLen(3))
			Expect(bursts[0].EventIDs).To(HaveLen(3)) // Team Leadership
			Expect(bursts[1].EventIDs).To(HaveLen(2)) // Platform Migration or Product Launch
			Expect(bursts[2].EventIDs).To(HaveLen(2)) // Platform Migration or Product Launch
		})

		It("should apply pagination with limit", func() {
			filters := fixtures.BurstListFiltersWithSort("name", "asc")
			filters.Limit = 2

			bursts, err := repository.List(ctx, *filters)
			Expect(err).ToNot(HaveOccurred())
			Expect(bursts).To(HaveLen(2))
		})

		It("should apply pagination with offset", func() {
			filters := fixtures.BurstListFiltersWithSort("name", "asc")
			filters.Offset = 1

			bursts, err := repository.List(ctx, *filters)
			Expect(err).ToNot(HaveOccurred())
			Expect(bursts).To(HaveLen(2))
			Expect(bursts[0].Name).To(Equal("Product Launch"))
		})

		It("should apply pagination with both limit and offset", func() {
			filters := fixtures.BurstListFiltersWithSort("name", "asc")
			filters.Offset = 1
			filters.Limit = 1

			bursts, err := repository.List(ctx, *filters)
			Expect(err).ToNot(HaveOccurred())
			Expect(bursts).To(HaveLen(1))
			Expect(bursts[0].Name).To(Equal("Product Launch"))
		})
	})

	Describe("Count", func() {
		BeforeEach(func() {
			// Create test bursts using factory
			for i := 0; i < 3; i++ {
				burst := fixtures.BurstFactory.MustCreate().(*career.Burst)
				burst.ID = "" // Clear to test auto-generation
				err := repository.Create(ctx, burst)
				Expect(err).ToNot(HaveOccurred())
			}
		})

		It("should count all bursts with no filters", func() {
			count, err := repository.Count(ctx, *fixtures.BurstListFilters())
			Expect(err).ToNot(HaveOccurred())
			Expect(count).To(Equal(3))
		})
	})
})
