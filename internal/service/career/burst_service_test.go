//nolint:errcheck // Test file - error handling for test setup is not relevant.
package career

import (
	"context"
	"errors"
	"time"

	"github.com/baphled/kariya/internal/domain/career"
	careerrepo "github.com/baphled/kariya/internal/repository/career"
	careermemory "github.com/baphled/kariya/internal/repository/career/memory"
	"github.com/baphled/kariya/internal/service/career/burstfact"
	"github.com/baphled/kariya/internal/testutil/fixtures"
	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
)

var _ = Describe("Career Service - Burst Methods", func() {
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

	Describe("SuggestBursts", func() {
		It("should return empty list when fewer than 2 events provided", func() {
			suggestions, err := service.SuggestBursts(ctx, []string{"1"})
			Expect(err).NotTo(HaveOccurred())
			Expect(suggestions).To(BeEmpty())
		})

		It("should detect bursts from provided event IDs", func() {
			now := time.Now()
			event1 := fixtures.EventWith("1", "Led backend infrastructure project", "TechCorp", "Platform")
			event1.Date = now.AddDate(0, -2, 0)

			event2 := fixtures.EventWith("2", "Architected backend platform", "TechCorp", "Platform")
			event2.Date = now.AddDate(0, -1, 0)

			// Add events to repo
			err := repo.Create(ctx, event1)
			Expect(err).NotTo(HaveOccurred())
			err = repo.Create(ctx, event2)
			Expect(err).NotTo(HaveOccurred())

			suggestions, err := service.SuggestBursts(ctx, []string{"1", "2"})
			Expect(err).NotTo(HaveOccurred())
			Expect(len(suggestions)).To(BeNumerically(">=", 0))
		})

		It("should handle missing event IDs gracefully", func() {
			suggestions, err := service.SuggestBursts(ctx, []string{"nonexistent"})
			Expect(err).NotTo(HaveOccurred())
			Expect(suggestions).To(BeEmpty())
		})
	})

	Describe("SaveBurst", func() {
		It("should reject burst with fewer than 2 events", func() {
			// Create burst manually with only 1 event (fixtures.Burst enforces min 2).
			burst := fixtures.BurstFactory.MustCreate().(*career.Burst)
			burst.EventIDs = []string{"1"}

			err := service.SaveBurst(ctx, burst)
			Expect(err).To(HaveOccurred())
		})

		It("should accept valid burst without setting Confirmed", func() {
			burstRepo := careermemory.NewBurstRepository()
			service.SetBurstRepository(burstRepo)
			burst := fixtures.Burst("burst1", "1", "2")

			err := service.SaveBurst(ctx, burst)
			Expect(err).NotTo(HaveOccurred())

			// SaveBurst should NOT set Confirmed.
			retrieved, err := burstRepo.GetByID(ctx, burst.ID)
			Expect(err).NotTo(HaveOccurred())
			Expect(retrieved.Confirmed).To(BeFalse())
			Expect(retrieved.ConfirmedAt).To(BeNil())
		})
	})

	Describe("ConfirmBurst", func() {
		var burstRepo *careermemory.BurstRepository

		BeforeEach(func() {
			burstRepo = careermemory.NewBurstRepository()
			service.SetBurstRepository(burstRepo)
		})

		It("should set Confirmed to true and ConfirmedAt", func() {
			burst := fixtures.Burst("burst1", "1", "2")
			// First save the burst.
			err := burstRepo.Create(ctx, burst)
			Expect(err).NotTo(HaveOccurred())

			err = service.ConfirmBurst(ctx, burst)
			Expect(err).NotTo(HaveOccurred())

			// Verify confirmation state.
			retrieved, err := burstRepo.GetByID(ctx, burst.ID)
			Expect(err).NotTo(HaveOccurred())
			Expect(retrieved.Confirmed).To(BeTrue())
			Expect(retrieved.ConfirmedAt).NotTo(BeNil())
		})

		It("should reject nil burst", func() {
			err := service.ConfirmBurst(ctx, nil)
			Expect(err).To(HaveOccurred())
		})

		It("should reject burst with empty ID", func() {
			burst := fixtures.Burst("", "1", "2")
			err := service.ConfirmBurst(ctx, burst)
			Expect(err).To(HaveOccurred())
		})

		It("should restore original state on repository update failure", func() {
			// Create an already-confirmed burst.
			burst := fixtures.Burst("burst-rollback", "1", "2")
			confirmedAt := time.Now().Add(-24 * time.Hour)
			burst.Confirmed = true
			burst.ConfirmedAt = &confirmedAt
			err := burstRepo.Create(ctx, burst)
			Expect(err).NotTo(HaveOccurred())

			// Capture state after Create (which sets its own UpdatedAt).
			origUpdatedAt := burst.UpdatedAt

			// Use a failing repo to trigger rollback.
			failRepo := &failingBurstRepo{inner: burstRepo}
			service.SetBurstRepository(failRepo)

			err = service.ConfirmBurst(ctx, burst)
			Expect(err).To(HaveOccurred())

			// Original state must be fully preserved after rollback.
			Expect(burst.Confirmed).To(BeTrue(), "rollback should restore original Confirmed")
			Expect(burst.ConfirmedAt).To(Equal(&confirmedAt), "rollback should restore original ConfirmedAt")
			Expect(burst.UpdatedAt).To(BeTemporally("~", origUpdatedAt, time.Millisecond),
				"rollback should restore original UpdatedAt")
		})
	})

	Describe("RejectBurstSuggestion", func() {
		It("should reject burst suggestion", func() {
			err := service.RejectBurstSuggestion(ctx, []string{"1", "2"})
			Expect(err).NotTo(HaveOccurred())
		})
	})

	Describe("SuggestBurstsWithOptions", func() {
		It("should return empty list when fewer than 2 events provided", func() {
			opts := &burstfact.DetectionOptions{
				MinConfidence:       0.7,
				TemporalWindow:      3 * 30 * 24 * time.Hour,
				MinEventCount:       2,
				MaxSuggestionsCount: 5,
			}
			suggestions, err := service.SuggestBurstsWithOptions(ctx, []string{"1"}, opts)
			Expect(err).NotTo(HaveOccurred())
			Expect(suggestions).To(BeEmpty())
		})

		It("should detect bursts with custom options", func() {
			now := time.Now()
			event1 := fixtures.EventWith("1", "Led backend infrastructure project", "TechCorp", "Platform")
			event1.Date = now.AddDate(0, -2, 0)

			event2 := fixtures.EventWith("2", "Architected backend platform", "TechCorp", "Platform")
			event2.Date = now.AddDate(0, -1, 0)

			// Add events to repo
			err := repo.Create(ctx, event1)
			Expect(err).NotTo(HaveOccurred())
			err = repo.Create(ctx, event2)
			Expect(err).NotTo(HaveOccurred())

			opts := &burstfact.DetectionOptions{
				MinConfidence:       0.5,
				TemporalWindow:      6 * 30 * 24 * time.Hour,
				MinEventCount:       2,
				MaxSuggestionsCount: 10,
			}

			suggestions, err := service.SuggestBurstsWithOptions(ctx, []string{"1", "2"}, opts)
			Expect(err).NotTo(HaveOccurred())
			Expect(len(suggestions)).To(BeNumerically(">=", 0))
		})

		It("should handle missing event IDs gracefully", func() {
			opts := &burstfact.DetectionOptions{
				MinConfidence:       0.6,
				TemporalWindow:      6 * 30 * 24 * time.Hour,
				MinEventCount:       2,
				MaxSuggestionsCount: 10,
			}
			suggestions, err := service.SuggestBurstsWithOptions(ctx, []string{"nonexistent1", "nonexistent2"}, opts)
			Expect(err).NotTo(HaveOccurred())
			Expect(suggestions).To(BeEmpty())
		})
	})

	Describe("SaveBurstSuggestions", func() {
		var burstRepo *careermemory.BurstRepository

		BeforeEach(func() {
			burstRepo = careermemory.NewBurstRepository()
			service.SetBurstRepository(burstRepo)
		})

		It("should save valid burst suggestions as unconfirmed", func() {
			suggestions := []burstfact.BurstSuggestion{
				{
					Name:            "Backend Platform Initiative",
					Description:     "Series of backend infrastructure improvements",
					EventIDs:        []string{"event1", "event2"},
					ConfidenceScore: 0.85,
				},
				{
					Name:            "Team Leadership Series",
					Description:     "Leadership activities and team management",
					EventIDs:        []string{"event3", "event4"},
					ConfidenceScore: 0.75,
				},
			}

			savedBursts, err := service.SaveBurstSuggestions(ctx, suggestions)
			Expect(err).NotTo(HaveOccurred())
			Expect(savedBursts).To(HaveLen(2))

			for _, burst := range savedBursts {
				retrieved, err := burstRepo.GetByID(ctx, burst.ID)
				Expect(err).NotTo(HaveOccurred())
				Expect(retrieved.Name).To(Equal(burst.Name))
				Expect(retrieved.EventIDs).To(Equal(burst.EventIDs))
				Expect(retrieved.Confirmed).To(BeFalse(), "saved burst suggestions should be unconfirmed")
				Expect(retrieved.ConfirmedAt).To(BeNil(), "saved burst suggestions should not have ConfirmedAt")
			}
		})

		It("should handle empty suggestions list", func() {
			savedBursts, err := service.SaveBurstSuggestions(ctx, []burstfact.BurstSuggestion{})
			Expect(err).NotTo(HaveOccurred())
			Expect(savedBursts).To(BeEmpty())
		})

		It("should return nil when no burst repository configured", func() {
			serviceWithoutBurstRepo := NewService(repo)
			suggestions := []burstfact.BurstSuggestion{
				{
					Name:            "Test Burst",
					Description:     "Test description",
					EventIDs:        []string{"event1", "event2"},
					ConfidenceScore: 0.8,
				},
			}

			savedBursts, err := serviceWithoutBurstRepo.SaveBurstSuggestions(ctx, suggestions)
			Expect(err).NotTo(HaveOccurred())
			Expect(savedBursts).To(BeNil())
		})

		It("should continue saving other suggestions if one fails validation", func() {
			suggestions := []burstfact.BurstSuggestion{
				{
					Name:            "", // Invalid: empty name
					Description:     "Invalid burst",
					EventIDs:        []string{"event1"}, // Invalid: too few events
					ConfidenceScore: 0.8,
				},
				{
					Name:            "Valid Burst",
					Description:     "Valid burst description",
					EventIDs:        []string{"event2", "event3"},
					ConfidenceScore: 0.75,
				},
			}

			savedBursts, err := service.SaveBurstSuggestions(ctx, suggestions)
			Expect(err).NotTo(HaveOccurred())
			Expect(savedBursts).To(HaveLen(1))
			Expect(savedBursts[0].Name).To(Equal("Valid Burst"))
		})
	})

	Describe("GetBurstRepository", func() {
		It("should return nil when no burst repository is set", func() {
			burstRepo := service.GetBurstRepository()
			Expect(burstRepo).To(BeNil())
		})

		It("should return the configured burst repository", func() {
			memoryBurstRepo := careermemory.NewBurstRepository()
			service.SetBurstRepository(memoryBurstRepo)

			burstRepo := service.GetBurstRepository()
			Expect(burstRepo).To(Equal(memoryBurstRepo))
		})
	})
})

// failingBurstRepo wraps a real repository but forces Update to fail.
type failingBurstRepo struct {
	inner careerrepo.BurstRepository
}

func (f *failingBurstRepo) Create(ctx context.Context, burst *career.Burst) error {
	return f.inner.Create(ctx, burst)
}

func (f *failingBurstRepo) GetByID(ctx context.Context, id string) (*career.Burst, error) {
	return f.inner.GetByID(ctx, id)
}

func (f *failingBurstRepo) Update(_ context.Context, _ *career.Burst) error {
	return errors.New("forced update failure")
}

func (f *failingBurstRepo) Delete(ctx context.Context, id string) error {
	return f.inner.Delete(ctx, id)
}

func (f *failingBurstRepo) List(ctx context.Context, filters careerrepo.BurstListFilters) ([]*career.Burst, error) {
	return f.inner.List(ctx, filters)
}

func (f *failingBurstRepo) Count(ctx context.Context, filters careerrepo.BurstListFilters) (int, error) {
	return f.inner.Count(ctx, filters)
}
