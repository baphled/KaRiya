package burst_fact

import (
	"context"

	"github.com/baphled/kariya/internal/domain/career"
	"github.com/baphled/kariya/internal/testutil/fixtures"
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
			event := fixtures.EventValWith("1", "Single event", "", "")

			suggestions, err := detector.DetectBursts(ctx, []career.CareerEvent{event}, nil)
			Expect(err).NotTo(HaveOccurred())
			Expect(suggestions).To(BeEmpty())
		})

		It("should build similarity matrix for event pairs", func() {
			events := []career.CareerEvent{
				fixtures.EventValWith("1", "Backend work", "", ""),
				fixtures.EventValWith("2", "Backend platform", "", ""),
			}

			matrix := detector.buildSimilarityMatrix(events)
			Expect(matrix).NotTo(BeNil())
			Expect(matrix["1"]).NotTo(BeNil())
			Expect(len(matrix)).To(Equal(2))
		})

		It("should find clusters of similar events", func() {
			events := []career.CareerEvent{
				fixtures.EventValWith("1", "Backend infrastructure work", "", ""),
				fixtures.EventValWith("2", "Backend platform architecture", "", ""),
			}

			matrix := detector.buildSimilarityMatrix(events)
			clusters := detector.findClusters(events, matrix, 0.3)
			Expect(len(clusters)).To(BeNumerically(">", 0))
		})

		It("should convert cluster to suggestion with confidence score", func() {
			events := []career.CareerEvent{
				fixtures.EventValWith("1", "Backend work", "", ""),
				fixtures.EventValWith("2", "Backend platform", "", ""),
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
			events := []career.CareerEvent{
				fixtures.EventValWith("1", "Event A", "", ""),
				fixtures.EventValWith("2", "Event B", "", ""),
			}

			opts := &DetectionOptions{MinConfidence: 0.99}
			suggestions, err := detector.DetectBursts(ctx, events, opts)
			Expect(err).NotTo(HaveOccurred())
			Expect(len(suggestions)).To(Equal(0))
		})

		It("should limit results to MaxSuggestionsCount", func() {
			events := make([]career.CareerEvent, 15)
			for i := 0; i < 15; i++ {
				events[i] = fixtures.EventValWith(string(rune(48+i)), "Backend infrastructure work", "TechCorp", "")
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
			events := []career.CareerEvent{
				fixtures.EventValWith("1", "Backend infrastructure work", "TechCorp", ""),
				fixtures.EventValWith("2", "Backend platform", "TechCorp", ""),
				fixtures.EventValWith("3", "Unrelated event", "OtherCorp", ""),
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
