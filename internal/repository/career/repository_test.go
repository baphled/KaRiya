package career

import (
	"context"

	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"

	"github.com/baphled/kariya/internal/domain/career"
	"github.com/baphled/kariya/internal/testutil/fixtures"
)

var _ = Describe("Memory Repository", func() {
	var (
		repo *MemoryRepository
		ctx  context.Context
	)

	BeforeEach(func() {
		repo = NewMemoryRepository()
		ctx = context.Background()
	})

	Describe("Create", func() {
		It("should successfully create an event", func() {
			event := fixtures.EventFactory.MustCreate().(*career.CareerEvent)
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
			Expect(err).To(Equal(ErrDuplicateEvent))
		})
	})

	Describe("GetByID", func() {
		It("should retrieve existing event", func() {
			event := fixtures.EventFactory.MustCreate().(*career.CareerEvent)
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
			Expect(err).To(Equal(ErrEventNotFound))
		})
	})

	Describe("Update", func() {
		It("should successfully update an event", func() {
			event := fixtures.EventFactory.MustCreate().(*career.CareerEvent)
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
			event := fixtures.EventFactory.MustCreate().(*career.CareerEvent)
			// Don't create it - just try to update

			err := repo.Update(ctx, event)
			Expect(err).To(HaveOccurred())
			Expect(err).To(Equal(ErrEventNotFound))
		})
	})

	Describe("Delete", func() {
		It("should successfully delete an event", func() {
			event := fixtures.EventFactory.MustCreate().(*career.CareerEvent)
			event.ID = "" // Clear to test auto-generation

			err := repo.Create(ctx, event)
			Expect(err).NotTo(HaveOccurred())

			err = repo.Delete(ctx, event.ID)
			Expect(err).NotTo(HaveOccurred())

			_, err = repo.GetByID(ctx, event.ID)
			Expect(err).To(HaveOccurred())
			Expect(err).To(Equal(ErrEventNotFound))
		})

		It("should prevent deletion of non-existent event", func() {
			err := repo.Delete(ctx, "non-existent-id")
			Expect(err).To(HaveOccurred())
			Expect(err).To(Equal(ErrEventNotFound))
		})
	})

	Describe("List", func() {
		It("should list events with no filters", func() {
			// Create multiple events using factory
			for i := 0; i < 5; i++ {
				event := fixtures.EventFactory.MustCreate().(*career.CareerEvent)
				event.ID = "" // Clear to test auto-generation
				err := repo.Create(ctx, event)
				Expect(err).NotTo(HaveOccurred())
			}

			events, err := repo.List(ctx, ListFilters{
				Limit: 10,
			})
			Expect(err).NotTo(HaveOccurred())
			Expect(events).To(HaveLen(5))
		})

		It("should list events with tag filter", func() {
			// Create events with different tags
			event1 := fixtures.EventFactory.MustCreate().(*career.CareerEvent)
			event1.ID = "" // Clear to test auto-generation
			event1.Tags = []string{"project"}

			event2 := fixtures.EventFactory.MustCreate().(*career.CareerEvent)
			event2.ID = "" // Clear to test auto-generation
			event2.Tags = []string{"technical"}

			err := repo.Create(ctx, event1)
			Expect(err).NotTo(HaveOccurred())
			err = repo.Create(ctx, event2)
			Expect(err).NotTo(HaveOccurred())

			events, err := repo.List(ctx, ListFilters{
				Tags:  []string{"project"},
				Limit: 10,
			})
			Expect(err).NotTo(HaveOccurred())
			Expect(events).To(HaveLen(1))
			Expect(events[0].Tags[0]).To(Equal("project"))
		})
	})

	Describe("Count", func() {
		It("should count events with no filters", func() {
			// Create multiple events using factory
			for i := 0; i < 5; i++ {
				event := fixtures.EventFactory.MustCreate().(*career.CareerEvent)
				event.ID = "" // Clear to test auto-generation
				err := repo.Create(ctx, event)
				Expect(err).NotTo(HaveOccurred())
			}

			count, err := repo.Count(ctx, ListFilters{})
			Expect(err).NotTo(HaveOccurred())
			Expect(count).To(Equal(5))
		})

		It("should count events with tag filter", func() {
			// Create events with different tags
			event1 := fixtures.EventFactory.MustCreate().(*career.CareerEvent)
			event1.ID = "" // Clear to test auto-generation
			event1.Tags = []string{"project"}

			event2 := fixtures.EventFactory.MustCreate().(*career.CareerEvent)
			event2.ID = "" // Clear to test auto-generation
			event2.Tags = []string{"technical"}

			err := repo.Create(ctx, event1)
			Expect(err).NotTo(HaveOccurred())
			err = repo.Create(ctx, event2)
			Expect(err).NotTo(HaveOccurred())

			count, err := repo.Count(ctx, ListFilters{
				Tags: []string{"project"},
			})
			Expect(err).NotTo(HaveOccurred())
			Expect(count).To(Equal(1))
		})
	})
})
