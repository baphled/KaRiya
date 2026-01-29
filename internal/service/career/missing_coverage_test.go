package career

import (
	"context"
	"time"

	"github.com/baphled/kariya/internal/domain/career"
	careerrepo "github.com/baphled/kariya/internal/repository/career"
	careermemory "github.com/baphled/kariya/internal/repository/career/memory"
	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
)

var _ = Describe("Career Service - Missing Coverage", func() {
	var (
		repo    careerrepo.EventRepository
		service *Service
		ctx     context.Context
	)

	BeforeEach(func() {
		repo = careermemory.NewEventRepository()
		service = NewService(repo)
		ctx = context.Background()
	})

	Describe("DeleteBurst", func() {
		var (
			burstRepo *careermemory.BurstRepository
			testBurst *career.Burst
		)

		BeforeEach(func() {
			burstRepo = careermemory.NewBurstRepository()
			service.SetBurstRepository(burstRepo)

			// Create a test burst
			now := time.Now()
			testBurst = &career.Burst{
				ID:          "burst-test-1",
				Name:        "Test Burst",
				Description: "A test burst for deletion",
				EventIDs:    []string{"event-1", "event-2"},
				Confirmed:   false,
				CreatedAt:   now,
				UpdatedAt:   now,
			}

			err := burstRepo.Create(ctx, testBurst)
			Expect(err).NotTo(HaveOccurred())
		})

		It("should delete an existing burst", func() {
			err := service.DeleteBurst(ctx, testBurst.ID)
			Expect(err).NotTo(HaveOccurred())

			// Verify burst was deleted
			_, err = burstRepo.GetByID(ctx, testBurst.ID)
			Expect(err).To(HaveOccurred())
		})

		It("should return error when burst ID is empty", func() {
			err := service.DeleteBurst(ctx, "")
			Expect(err).To(HaveOccurred())
			Expect(err.Error()).To(ContainSubstring("burst ID cannot be empty"))
		})

		It("should return error when burst does not exist", func() {
			err := service.DeleteBurst(ctx, "nonexistent-burst-id")
			Expect(err).To(HaveOccurred())
			Expect(err.Error()).To(ContainSubstring("failed to delete burst"))
		})

		It("should return ErrBurstRepositoryNotConfigured when repository is nil", func() {
			serviceWithoutBurst := NewService(repo)
			err := serviceWithoutBurst.DeleteBurst(ctx, "burst-1")
			Expect(err).To(Equal(ErrBurstRepositoryNotConfigured))
		})

		It("should handle concurrent deletions gracefully", func() {
			// First deletion should succeed
			err1 := service.DeleteBurst(ctx, testBurst.ID)
			Expect(err1).NotTo(HaveOccurred())

			// Second deletion should fail (already deleted)
			err2 := service.DeleteBurst(ctx, testBurst.ID)
			Expect(err2).To(HaveOccurred())
		})

		It("should not affect other bursts when deleting one", func() {
			// Create another burst
			now := time.Now()
			burst2 := &career.Burst{
				ID:          "burst-test-2",
				Name:        "Second Burst",
				Description: "Another burst",
				EventIDs:    []string{"event-3", "event-4"}, // Need at least 2 events
				Confirmed:   false,
				CreatedAt:   now,
				UpdatedAt:   now,
			}
			err := burstRepo.Create(ctx, burst2)
			Expect(err).NotTo(HaveOccurred())

			// Delete first burst
			err = service.DeleteBurst(ctx, testBurst.ID)
			Expect(err).NotTo(HaveOccurred())

			// Verify second burst still exists
			retrieved, err := burstRepo.GetByID(ctx, burst2.ID)
			Expect(err).NotTo(HaveOccurred())
			Expect(retrieved.ID).To(Equal(burst2.ID))
			Expect(retrieved.Name).To(Equal("Second Burst"))
		})
	})

	Describe("GetEventRepository", func() {
		It("should return the event repository", func() {
			eventRepo := service.GetEventRepository()
			Expect(eventRepo).NotTo(BeNil())
			Expect(eventRepo).To(Equal(repo))
		})

		It("should return same repository instance on multiple calls", func() {
			repo1 := service.GetEventRepository()
			repo2 := service.GetEventRepository()
			Expect(repo1).To(BeIdenticalTo(repo2))
		})

		It("should allow using returned repository for operations", func() {
			eventRepo := service.GetEventRepository()

			// Create an event using returned repository
			event := &career.CareerEvent{
				ID:        "test-event-1",
				Text:      "Test event via GetEventRepository",
				Date:      time.Now(),
				Company:   "TestCo",
				CreatedAt: time.Now(),
				UpdatedAt: time.Now(),
			}

			err := eventRepo.Create(ctx, event)
			Expect(err).NotTo(HaveOccurred())

			// Retrieve via service to verify
			retrieved, err := service.GetEventByID(ctx, event.ID)
			Expect(err).NotTo(HaveOccurred())
			Expect(retrieved.ID).To(Equal(event.ID))
			Expect(retrieved.Text).To(Equal(event.Text))
		})
	})

	Describe("GetFactsBySourceEventID", func() {
		var factRepo *careermemory.FactRepository

		BeforeEach(func() {
			factRepo = careermemory.NewFactRepository()
			service.SetFactRepository(factRepo)
		})

		It("should return empty slice when repository is nil", func() {
			serviceWithoutFact := NewService(repo)
			facts, err := serviceWithoutFact.GetFactsBySourceEventID(ctx, "event-1")
			Expect(err).NotTo(HaveOccurred())
			Expect(facts).To(BeEmpty())
		})

		It("should return error for empty event ID", func() {
			facts, err := service.GetFactsBySourceEventID(ctx, "")
			Expect(err).To(HaveOccurred())
			Expect(err.Error()).To(ContainSubstring("event ID cannot be empty"))
			Expect(facts).To(BeNil())
		})

		It("should return empty list for event with no facts", func() {
			facts, err := service.GetFactsBySourceEventID(ctx, "event-no-facts")
			Expect(err).NotTo(HaveOccurred())
			Expect(facts).To(BeEmpty())
		})
	})

	Describe("GetFactsBySourceBurstID", func() {
		var factRepo *careermemory.FactRepository

		BeforeEach(func() {
			factRepo = careermemory.NewFactRepository()
			service.SetFactRepository(factRepo)
		})

		It("should return empty slice when repository is nil", func() {
			serviceWithoutFact := NewService(repo)
			facts, err := serviceWithoutFact.GetFactsBySourceBurstID(ctx, "burst-1")
			Expect(err).NotTo(HaveOccurred())
			Expect(facts).To(BeEmpty())
		})

		It("should return error for empty burst ID", func() {
			facts, err := service.GetFactsBySourceBurstID(ctx, "")
			Expect(err).To(HaveOccurred())
			Expect(err.Error()).To(ContainSubstring("burst ID cannot be empty"))
			Expect(facts).To(BeNil())
		})

		It("should return empty list for burst with no facts", func() {
			facts, err := service.GetFactsBySourceBurstID(ctx, "burst-no-facts")
			Expect(err).NotTo(HaveOccurred())
			Expect(facts).To(BeEmpty())
		})
	})
})
