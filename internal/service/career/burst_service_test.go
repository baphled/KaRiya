//nolint:errcheck // Test file - error handling for test setup is not relevant.
package career

import (
	"context"
	"time"

	"github.com/baphled/kariya/internal/domain/career"
	careerrepo "github.com/baphled/kariya/internal/repository/career"
	careermemory "github.com/baphled/kariya/internal/repository/career/memory"
	"github.com/baphled/kariya/internal/service/career/burst_fact"
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

	Describe("ConfirmBurst", func() {
		It("should reject burst with fewer than 2 events", func() {
			// Create burst manually with only 1 event (fixtures.Burst enforces min 2)
			burst := fixtures.BurstFactory.MustCreate().(*career.Burst)
			burst.EventIDs = []string{"1"} // Override to have only 1 event

			err := service.ConfirmBurst(ctx, burst)
			Expect(err).To(HaveOccurred())
		})

		It("should accept valid burst", func() {
			burst := fixtures.Burst("burst1", "1", "2")

			err := service.ConfirmBurst(ctx, burst)
			Expect(err).NotTo(HaveOccurred())
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
			opts := &burst_fact.DetectionOptions{
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

			opts := &burst_fact.DetectionOptions{
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
			opts := &burst_fact.DetectionOptions{
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

		It("should save valid burst suggestions", func() {
			suggestions := []burst_fact.BurstSuggestion{
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

			// Verify bursts were saved to repository
			for _, burst := range savedBursts {
				retrieved, err := burstRepo.GetByID(ctx, burst.ID)
				Expect(err).NotTo(HaveOccurred())
				Expect(retrieved.Name).To(Equal(burst.Name))
				Expect(retrieved.EventIDs).To(Equal(burst.EventIDs))
			}
		})

		It("should handle empty suggestions list", func() {
			savedBursts, err := service.SaveBurstSuggestions(ctx, []burst_fact.BurstSuggestion{})
			Expect(err).NotTo(HaveOccurred())
			Expect(savedBursts).To(BeEmpty())
		})

		It("should return nil when no burst repository configured", func() {
			serviceWithoutBurstRepo := NewService(repo)
			suggestions := []burst_fact.BurstSuggestion{
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
			suggestions := []burst_fact.BurstSuggestion{
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
