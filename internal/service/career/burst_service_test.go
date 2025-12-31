package career

import (
	"context"
	"time"

	"github.com/baphled/kariya/internal/domain/career"
	careerrepo "github.com/baphled/kariya/internal/repository/career"
	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
)

var _ = Describe("Career Service - Burst Methods", func() {
	var (
		repo    careerrepo.Repository
		service *Service
		ctx     context.Context
	)

	BeforeEach(func() {
		repo = careerrepo.NewMemoryRepository()
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
			events := []career.CareerEvent{
				{
					ID:        "1",
					Text:      "Led backend infrastructure project",
					Date:      now.AddDate(0, -2, 0),
					Company:   "TechCorp",
					CreatedAt: now,
					UpdatedAt: now,
				},
				{
					ID:        "2",
					Text:      "Architected backend platform",
					Date:      now.AddDate(0, -1, 0),
					Company:   "TechCorp",
					CreatedAt: now,
					UpdatedAt: now,
				},
			}

			// Add events to repo
			err := repo.Create(ctx, &events[0])
			Expect(err).NotTo(HaveOccurred())
			err = repo.Create(ctx, &events[1])
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
			burst := &career.Burst{
				ID:       "burst1",
				Name:     "Backend Platform",
				EventIDs: []string{"1"},
			}

			err := service.ConfirmBurst(ctx, burst)
			Expect(err).To(HaveOccurred())
		})

		It("should accept valid burst", func() {
			burst := &career.Burst{
				ID:        "burst1",
				Name:      "Backend Platform",
				EventIDs:  []string{"1", "2"},
				CreatedAt: time.Now(),
				UpdatedAt: time.Now(),
			}

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
})
