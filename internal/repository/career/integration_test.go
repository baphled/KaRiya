package career

import (
	"context"
	"fmt"
	"math/rand"
	"sync"
	"time"

	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"

	"github.com/baphled/kariya/internal/domain/career"
)

var _ = Describe("Career Event Repository Integration Tests", func() {
	var (
		repo *MemoryRepository
		ctx  context.Context
	)

	BeforeEach(func() {
		repo = NewMemoryRepository()
		ctx = context.Background()
	})

	// Helper function to generate unique test events
	generateTestEvent := func() *career.CareerEvent {
		return &career.CareerEvent{
			Text:    fmt.Sprintf("Test Event %d", rand.Intn(10000)),
			Date:    time.Now().Add(time.Duration(rand.Intn(365)) * -24 * time.Hour),
			Tags:    []string{"project"},
			Company: "Test Company",
		}
	}

	Describe("Concurrent Event Creation", func() {
		It("should handle multiple concurrent event creations safely", func() {
			// Number of concurrent events to create
			eventCount := 100
			var wg sync.WaitGroup
			var mu sync.Mutex
			createdEvents := make(map[string]bool)

			wg.Add(eventCount)
			for i := 0; i < eventCount; i++ {
				go func() {
					defer wg.Done()

					event := generateTestEvent()
					err := repo.Create(ctx, event)

					mu.Lock()
					defer mu.Unlock()

					if err == nil {
						createdEvents[event.ID] = true
					}
				}()
			}

			wg.Wait()

			// Verify all events were created successfully
			count, err := repo.Count(ctx, ListFilters{})
			Expect(err).NotTo(HaveOccurred())
			Expect(count).To(Equal(len(createdEvents)),
				"All unique events should be created")
		})
	})

	Describe("Complex Event Filtering", func() {
		BeforeEach(func() {
			// Create multiple events with different characteristics
			baseDate := time.Now()
			for i := 0; i < 10; i++ {
				event := generateTestEvent()
				event.Date = baseDate.AddDate(0, 0, -i*30)
				event.Tags = []string{"project"}

				err := repo.Create(ctx, event)
				Expect(err).NotTo(HaveOccurred())
			}
		})

		It("should filter events by tag", func() {
			filteredByTag, err := repo.List(ctx, ListFilters{
				Tags:  []string{"project"},
				Limit: 10,
			})
			Expect(err).NotTo(HaveOccurred())
			Expect(filteredByTag).To(HaveLen(10),
				"All events should have project tag")
		})

		It("should filter events by date range", func() {
			baseDate := time.Now()
			oneMonthAgo := baseDate.AddDate(0, 0, -30)

			filteredByDate, err := repo.List(ctx, ListFilters{
				StartDate: &oneMonthAgo,
				Limit:     10,
			})
			Expect(err).NotTo(HaveOccurred())
			Expect(len(filteredByDate)).To(BeNumerically(">", 0),
				"Should return at least one event within date range")
			Expect(len(filteredByDate)).To(BeNumerically("<=", 10),
				"Should respect limit parameter")
		})
	})

	Describe("Error Scenarios", func() {
		It("should return error on duplicate event creation", func() {
			event := generateTestEvent()
			event.ID = "fixed-id"

			// First creation should succeed
			err := repo.Create(ctx, event)
			Expect(err).NotTo(HaveOccurred())

			// Second creation with same ID should fail
			err = repo.Create(ctx, event)
			Expect(err).To(HaveOccurred())
			Expect(err).To(Equal(ErrDuplicateEvent))
		})

		It("should return error when updating non-existent event", func() {
			nonExistentEvent := generateTestEvent()
			nonExistentEvent.ID = "non-existent-id"

			err := repo.Update(ctx, nonExistentEvent)
			Expect(err).To(HaveOccurred())
			Expect(err).To(Equal(ErrEventNotFound))
		})

		It("should return error when deleting non-existent event", func() {
			err := repo.Delete(ctx, "non-existent-id")
			Expect(err).To(HaveOccurred())
			Expect(err).To(Equal(ErrEventNotFound))
		})
	})
})
