package burst_fact

import (
	"context"
	"time"

	"github.com/baphled/kariya/internal/domain/career"
	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
)

var _ = Describe("BurstDetector", func() {
	var (
		detector *BurstDetector
		ctx      context.Context
	)

	BeforeEach(func() {
		detector = NewBurstDetector()
		ctx = context.Background()
	})

	Describe("DetectBursts", func() {
		It("should return empty list for fewer than minimum events", func() {
			event := career.CareerEvent{
				ID:        "1",
				Text:      "Single event",
				Date:      time.Now(),
				CreatedAt: time.Now(),
				UpdatedAt: time.Now(),
			}

			suggestions, err := detector.DetectBursts(ctx, []career.CareerEvent{event}, nil)
			Expect(err).NotTo(HaveOccurred())
			Expect(suggestions).To(BeEmpty())
		})

		It("should build similarity matrix for event pairs", func() {
			now := time.Now()
			events := []career.CareerEvent{
				{
					ID:        "1",
					Text:      "Backend work",
					Date:      now,
					CreatedAt: now,
					UpdatedAt: now,
				},
				{
					ID:        "2",
					Text:      "Backend platform",
					Date:      now.AddDate(0, -1, 0),
					CreatedAt: now,
					UpdatedAt: now,
				},
			}

			matrix := detector.buildSimilarityMatrix(events)
			Expect(matrix).NotTo(BeNil())
			Expect(matrix["1"]).NotTo(BeNil())
			Expect(len(matrix)).To(Equal(2))
		})

		It("should find clusters of similar events", func() {
			now := time.Now()
			events := []career.CareerEvent{
				{
					ID:        "1",
					Text:      "Backend infrastructure work",
					Date:      now,
					CreatedAt: now,
					UpdatedAt: now,
				},
				{
					ID:        "2",
					Text:      "Backend platform architecture",
					Date:      now.AddDate(0, -1, 0),
					CreatedAt: now,
					UpdatedAt: now,
				},
			}

			matrix := detector.buildSimilarityMatrix(events)
			clusters := detector.findClusters(events, matrix, 0.3)
			Expect(len(clusters)).To(BeNumerically(">", 0))
		})

		It("should convert cluster to suggestion with confidence score", func() {
			now := time.Now()
			events := []career.CareerEvent{
				{
					ID:        "1",
					Text:      "Backend work",
					Date:      now,
					CreatedAt: now,
					UpdatedAt: now,
				},
				{
					ID:        "2",
					Text:      "Backend platform",
					Date:      now.AddDate(0, -1, 0),
					CreatedAt: now,
					UpdatedAt: now,
				},
			}

			matrix := detector.buildSimilarityMatrix(events)
			suggestion := detector.clusterToSuggestion(events, matrix)

			Expect(suggestion.EventIDs).To(Equal([]string{"1", "2"}))
			Expect(suggestion.ConfidenceScore).To(BeNumerically(">=", 0.0))
			Expect(suggestion.ConfidenceScore).To(BeNumerically("<=", 1.0))
		})

		It("should validate burst suggestions", func() {
			validSuggestion := BurstSuggestion{
				EventIDs:        []string{"1", "2"},
				ConfidenceScore: 0.8,
			}

			err := detector.ValidateSuggestion(validSuggestion)
			Expect(err).NotTo(HaveOccurred())
		})

		It("should respect minimum confidence threshold", func() {
			now := time.Now()
			events := []career.CareerEvent{
				{
					ID:        "1",
					Text:      "Event A",
					Date:      now,
					CreatedAt: now,
					UpdatedAt: now,
				},
				{
					ID:        "2",
					Text:      "Event B",
					Date:      now.AddDate(0, -1, 0),
					CreatedAt: now,
					UpdatedAt: now,
				},
			}

			opts := &DetectionOptions{MinConfidence: 0.99}
			suggestions, err := detector.DetectBursts(ctx, events, opts)
			Expect(err).NotTo(HaveOccurred())
			Expect(len(suggestions)).To(Equal(0))
		})

		It("should limit results to MaxSuggestionsCount", func() {
			now := time.Now()
			events := []career.CareerEvent{}
			for i := 0; i < 15; i++ {
				events = append(events, career.CareerEvent{
					ID:        string(rune(48 + i)),
					Text:      "Backend infrastructure work",
					Date:      now.AddDate(0, -int(i/2), 0),
					Company:   "TechCorp",
					CreatedAt: now,
					UpdatedAt: now,
				})
			}

			opts := &DetectionOptions{MaxSuggestionsCount: 5}
			suggestions, err := detector.DetectBursts(ctx, events, opts)
			Expect(err).NotTo(HaveOccurred())
			Expect(len(suggestions)).To(BeNumerically("<=", 5))
		})

		It("should handle empty event list", func() {
			suggestions, err := detector.DetectBursts(ctx, []career.CareerEvent{}, nil)
			Expect(err).NotTo(HaveOccurred())
			Expect(suggestions).To(BeEmpty())
		})

		It("should sort suggestions by confidence", func() {
			now := time.Now()
			events := []career.CareerEvent{
				{
					ID:        "1",
					Text:      "Backend infrastructure work",
					Date:      now.AddDate(0, -2, 0),
					Company:   "TechCorp",
					CreatedAt: now,
					UpdatedAt: now,
				},
				{
					ID:        "2",
					Text:      "Backend platform",
					Date:      now.AddDate(0, -1, 0),
					Company:   "TechCorp",
					CreatedAt: now,
					UpdatedAt: now,
				},
				{
					ID:        "3",
					Text:      "Unrelated event",
					Date:      now,
					Company:   "OtherCorp",
					CreatedAt: now,
					UpdatedAt: now,
				},
			}

			suggestions, err := detector.DetectBursts(ctx, events, nil)
			Expect(err).NotTo(HaveOccurred())

			for i := 0; i < len(suggestions)-1; i++ {
				Expect(suggestions[i].ConfidenceScore).To(BeNumerically(">=",
					suggestions[i+1].ConfidenceScore,
				))
			}
		})
	})

	Describe("ValidateSuggestion", func() {
		It("should reject suggestion with fewer than 2 events", func() {
			suggestion := BurstSuggestion{
				EventIDs:        []string{"1"},
				ConfidenceScore: 0.8,
			}

			err := detector.ValidateSuggestion(suggestion)
			Expect(err).To(HaveOccurred())
			Expect(err.Error()).To(ContainSubstring("at least 2 events"))
		})

		It("should reject suggestion with invalid confidence", func() {
			suggestion := BurstSuggestion{
				EventIDs:        []string{"1", "2"},
				ConfidenceScore: 1.5,
			}

			err := detector.ValidateSuggestion(suggestion)
			Expect(err).To(HaveOccurred())
			Expect(err.Error()).To(ContainSubstring("between 0.0 and 1.0"))
		})

		It("should reject suggestion with duplicate event IDs", func() {
			suggestion := BurstSuggestion{
				EventIDs:        []string{"1", "1", "2"},
				ConfidenceScore: 0.8,
			}

			err := detector.ValidateSuggestion(suggestion)
			Expect(err).To(HaveOccurred())
			Expect(err.Error()).To(ContainSubstring("duplicate"))
		})

		It("should accept valid suggestion", func() {
			suggestion := BurstSuggestion{
				EventIDs:        []string{"1", "2"},
				ConfidenceScore: 0.8,
			}

			err := detector.ValidateSuggestion(suggestion)
			Expect(err).NotTo(HaveOccurred())
		})
	})
})
