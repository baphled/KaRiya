//nolint:errcheck // Test file - error handling for test setup is not relevant.
package memory

import (
	"context"
	"fmt"
	"math/rand"
	"sync"
	"time"

	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"

	"github.com/baphled/kariya/internal/domain/career"
	career_repo "github.com/baphled/kariya/internal/repository/career"
	"github.com/baphled/kariya/internal/testutil/fixtures"
)

var _ = Describe("EventRepository", func() {
	var (
		repo *EventRepository
		ctx  context.Context
	)

	BeforeEach(func() {
		repo = NewEventRepository()
		ctx = context.Background()
	})

	Describe("Create", func() {
		It("should successfully create an event", func() {
			event := fixtures.EventFactory.MustCreate().(*career.Event)
			event.ID = "" // Clear to test auto-generation

			err := repo.Create(ctx, event)
			Expect(err).NotTo(HaveOccurred())
			Expect(event.ID).NotTo(BeEmpty())
			Expect(event.CreatedAt).NotTo(BeZero())
			Expect(event.UpdatedAt).NotTo(BeZero())
		})

		It("should prevent duplicate event creation", func() {
			event := fixtures.Event("fixed-id")

			// First creation should succeed
			err := repo.Create(ctx, event)
			Expect(err).NotTo(HaveOccurred())

			// Second creation with same ID should fail
			err = repo.Create(ctx, event)
			Expect(err).To(HaveOccurred())
			Expect(err).To(Equal(career_repo.ErrDuplicateEvent))
		})
	})

	Describe("GetByID", func() {
		It("should retrieve existing event", func() {
			event := fixtures.EventFactory.MustCreate().(*career.Event)
			event.ID = "" // Clear to test auto-generation

			err := repo.Create(ctx, event)
			Expect(err).NotTo(HaveOccurred())

			retrievedEvent, err := repo.GetByID(ctx, event.ID)
			Expect(err).NotTo(HaveOccurred())
			Expect(retrievedEvent.ID).To(Equal(event.ID))
			Expect(retrievedEvent.Text).To(Equal(event.Text))
		})

		It("should return error for non-existent event", func() {
			_, err := repo.GetByID(ctx, "non-existent-id")
			Expect(err).To(HaveOccurred())
			Expect(err).To(Equal(career_repo.ErrEventNotFound))
		})
	})

	Describe("Update", func() {
		It("should successfully update an event", func() {
			event := fixtures.EventFactory.MustCreate().(*career.Event)
			event.ID = "" // Clear to test auto-generation

			err := repo.Create(ctx, event)
			Expect(err).NotTo(HaveOccurred())

			// Small delay to ensure timestamps differ on fast systems
			time.Sleep(10 * time.Millisecond)

			// Update event
			event.Text = "Updated Career Event Description"
			err = repo.Update(ctx, event)
			Expect(err).NotTo(HaveOccurred())

			updatedEvent, err := repo.GetByID(ctx, event.ID)
			Expect(err).NotTo(HaveOccurred())
			Expect(updatedEvent.Text).To(Equal("Updated Career Event Description"))
			Expect(updatedEvent.UpdatedAt).To(BeTemporally(">", updatedEvent.CreatedAt))
		})

		It("should prevent update of non-existent event", func() {
			event := fixtures.EventFactory.MustCreate().(*career.Event)
			// Don't create it - just try to update

			err := repo.Update(ctx, event)
			Expect(err).To(HaveOccurred())
			Expect(err).To(Equal(career_repo.ErrEventNotFound))
		})
	})

	Describe("Delete", func() {
		It("should successfully delete an event", func() {
			event := fixtures.EventFactory.MustCreate().(*career.Event)
			event.ID = "" // Clear to test auto-generation

			err := repo.Create(ctx, event)
			Expect(err).NotTo(HaveOccurred())

			err = repo.Delete(ctx, event.ID)
			Expect(err).NotTo(HaveOccurred())

			_, err = repo.GetByID(ctx, event.ID)
			Expect(err).To(HaveOccurred())
			Expect(err).To(Equal(career_repo.ErrEventNotFound))
		})

		It("should prevent deletion of non-existent event", func() {
			err := repo.Delete(ctx, "non-existent-id")
			Expect(err).To(HaveOccurred())
			Expect(err).To(Equal(career_repo.ErrEventNotFound))
		})
	})

	Describe("List", func() {
		It("should list events with no filters", func() {
			for i := 0; i < 5; i++ {
				event := fixtures.EventFactory.MustCreate().(*career.Event)
				event.ID = "" // Clear to test auto-generation
				err := repo.Create(ctx, event)
				Expect(err).NotTo(HaveOccurred())
			}

			events, err := repo.List(ctx, career_repo.EventListFilters{
				Limit: 10,
			})
			Expect(err).NotTo(HaveOccurred())
			Expect(events).To(HaveLen(5))
		})

		It("should list events with tag filter", func() {
			event1 := fixtures.EventFactory.MustCreate().(*career.Event)
			event1.ID = ""
			event1.Tags = []string{"project"}

			event2 := fixtures.EventFactory.MustCreate().(*career.Event)
			event2.ID = ""
			event2.Tags = []string{"technical"}

			err := repo.Create(ctx, event1)
			Expect(err).NotTo(HaveOccurred())
			err = repo.Create(ctx, event2)
			Expect(err).NotTo(HaveOccurred())

			events, err := repo.List(ctx, career_repo.EventListFilters{
				Tags:  []string{"project"},
				Limit: 10,
			})
			Expect(err).NotTo(HaveOccurred())
			Expect(events).To(HaveLen(1))
			Expect(events[0].Tags[0]).To(Equal("project"))
		})

		It("should filter events by date range", func() {
			baseDate := time.Now()
			for i := 0; i < 10; i++ {
				event := &career.Event{
					Text:    fmt.Sprintf("Event %d", i),
					Date:    baseDate.AddDate(0, 0, -i*30),
					Tags:    []string{"project"},
					Company: "Test Company",
				}
				err := repo.Create(ctx, event)
				Expect(err).NotTo(HaveOccurred())
			}

			oneMonthAgo := baseDate.AddDate(0, 0, -30)
			events, err := repo.List(ctx, career_repo.EventListFilters{
				StartDate: &oneMonthAgo,
				Limit:     10,
			})
			Expect(err).NotTo(HaveOccurred())
			Expect(len(events)).To(BeNumerically(">", 0))
			Expect(len(events)).To(BeNumerically("<=", 10))
		})
	})

	Describe("Count", func() {
		It("should count events with no filters", func() {
			for i := 0; i < 5; i++ {
				event := fixtures.EventFactory.MustCreate().(*career.Event)
				event.ID = "" // Clear to test auto-generation
				err := repo.Create(ctx, event)
				Expect(err).NotTo(HaveOccurred())
			}

			count, err := repo.Count(ctx, career_repo.EventListFilters{})
			Expect(err).NotTo(HaveOccurred())
			Expect(count).To(Equal(5))
		})

		It("should count events with tag filter", func() {
			event1 := fixtures.EventFactory.MustCreate().(*career.Event)
			event1.ID = ""
			event1.Tags = []string{"project"}

			event2 := fixtures.EventFactory.MustCreate().(*career.Event)
			event2.ID = ""
			event2.Tags = []string{"technical"}

			err := repo.Create(ctx, event1)
			Expect(err).NotTo(HaveOccurred())
			err = repo.Create(ctx, event2)
			Expect(err).NotTo(HaveOccurred())

			count, err := repo.Count(ctx, career_repo.EventListFilters{
				Tags: []string{"project"},
			})
			Expect(err).NotTo(HaveOccurred())
			Expect(count).To(Equal(1))
		})
	})

	Describe("Company and Project Fields", func() {
		It("should persist both fields correctly", func() {
			event := &career.Event{
				Text:    "Developed new microservice architecture",
				Date:    time.Now().Add(-24 * time.Hour),
				Company: "TechCorp Inc.",
				Project: "Platform Modernization",
				Tags:    []string{"technical"},
			}

			err := repo.Create(ctx, event)
			Expect(err).NotTo(HaveOccurred())

			retrieved, err := repo.GetByID(ctx, event.ID)
			Expect(err).NotTo(HaveOccurred())
			Expect(retrieved.Company).To(Equal("TechCorp Inc."))
			Expect(retrieved.Project).To(Equal("Platform Modernization"))
		})

		It("should handle empty company and project", func() {
			event := &career.Event{
				Text:    "Simple event without company or project",
				Date:    time.Now().Add(-24 * time.Hour),
				Company: "",
				Project: "",
				Tags:    []string{"technical"},
			}

			err := repo.Create(ctx, event)
			Expect(err).NotTo(HaveOccurred())

			retrieved, err := repo.GetByID(ctx, event.ID)
			Expect(err).NotTo(HaveOccurred())
			Expect(retrieved.Company).To(Equal(""))
			Expect(retrieved.Project).To(Equal(""))
		})

		It("should update company and project correctly", func() {
			event := &career.Event{
				Text:    "Initial event",
				Date:    time.Now().Add(-24 * time.Hour),
				Company: "OldCorp",
				Project: "Old Project",
				Tags:    []string{"technical"},
			}

			err := repo.Create(ctx, event)
			Expect(err).NotTo(HaveOccurred())

			event.Company = "NewCorp Ltd."
			event.Project = "Cloud Migration"
			err = repo.Update(ctx, event)
			Expect(err).NotTo(HaveOccurred())

			updated, err := repo.GetByID(ctx, event.ID)
			Expect(err).NotTo(HaveOccurred())
			Expect(updated.Company).To(Equal("NewCorp Ltd."))
			Expect(updated.Project).To(Equal("Cloud Migration"))
		})

		It("should list events with company and project", func() {
			events := []*career.Event{
				{
					Text:    "Event 1",
					Date:    time.Now().Add(-48 * time.Hour),
					Company: "Company A",
					Project: "Project Alpha",
					Tags:    []string{"technical"},
				},
				{
					Text:    "Event 2",
					Date:    time.Now().Add(-24 * time.Hour),
					Company: "Company B",
					Project: "Project Beta",
					Tags:    []string{"leadership"},
				},
			}

			for _, e := range events {
				err := repo.Create(ctx, e)
				Expect(err).NotTo(HaveOccurred())
			}

			listed, err := repo.List(ctx, career_repo.EventListFilters{
				Limit:     10,
				SortBy:    "date",
				SortOrder: "desc",
			})
			Expect(err).NotTo(HaveOccurred())
			Expect(listed).To(HaveLen(2))
			Expect(listed[0].Company).To(Equal("Company B"))
			Expect(listed[0].Project).To(Equal("Project Beta"))
			Expect(listed[1].Company).To(Equal("Company A"))
			Expect(listed[1].Project).To(Equal("Project Alpha"))
		})
	})

	Describe("Concurrent Access", func() {
		//nolint:gosec // G404: math/rand is fine for test data generation
		It("should handle concurrent event creations safely", func() {
			eventCount := 100
			var wg sync.WaitGroup
			var mu sync.Mutex
			createdEvents := make(map[string]bool)

			wg.Add(eventCount)
			for i := 0; i < eventCount; i++ {
				go func() {
					defer wg.Done()

					event := &career.Event{
						Text:    fmt.Sprintf("Test Event %d", rand.Intn(10000)),
						Date:    time.Now().Add(time.Duration(rand.Intn(365)) * -24 * time.Hour),
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

			count, err := repo.Count(ctx, career_repo.EventListFilters{})
			Expect(err).NotTo(HaveOccurred())
			Expect(count).To(Equal(len(createdEvents)))
		})
	})
})
