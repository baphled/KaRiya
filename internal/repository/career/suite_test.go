package career

import (
	"context"
	"sync"
	"testing"
	"time"

	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"

	"github.com/baphled/kariya/internal/domain/career"
)

// TestCareerRepositorySuite runs the entire repository test suite
func TestCareerRepositorySuite(t *testing.T) {
	RegisterFailHandler(Fail)
	RunSpecs(t, "Career Repository Suite")
}

var _ = Describe("Repository Test Suite", func() {
	// Memory Repository Integration Tests
	Context("Memory Repository Integration", func() {
		var (
			repo *MemoryRepository
			ctx  context.Context
		)

		BeforeEach(func() {
			repo = NewMemoryRepository()
			ctx = context.Background()
		})

		Describe("Concurrent Event Creation", func() {
			It("should create events concurrently in a thread-safe manner", func() {
				// Number of concurrent events to create
				eventCount := 100
				var wg sync.WaitGroup
				var mu sync.Mutex
				createdEvents := make(map[string]bool)

				wg.Add(eventCount)
				for i := 0; i < eventCount; i++ {
					go func() {
						defer wg.Done()

						event := &career.CareerEvent{
							Text:    "Test Event",
							Date:    time.Now(),
							Tags:    []string{"project"},
							Company: "Test Company",
						}
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
				Expect(len(createdEvents)).To(Equal(count), "Expected all unique events to be created")
			})
		})
	})
})
