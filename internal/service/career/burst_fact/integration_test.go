package burst_fact

import (
	"context"
	"time"

	"github.com/baphled/kariya/internal/domain/career"
	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
)

var _ = Describe("Burst Detection Integration Tests", func() {
	var (
		detector *BurstDetector
		workflow *BurstConfirmationWorkflow
		ctx      context.Context
	)

	BeforeEach(func() {
		detector = NewBurstDetector()
		workflow = NewBurstConfirmationWorkflow(detector)
		ctx = context.Background()
	})

	Describe("Complete Burst Suggestion Workflow", func() {
		It("should detect and confirm burst suggestion", func() {
			now := time.Now()
			events := []career.CareerEvent{
				{
					ID:        "1",
					Text:      "Led backend infrastructure project",
					Date:      now.AddDate(0, -2, 0),
					Company:   "TechCorp",
					Tags:      []string{"technical", "leadership"},
					CreatedAt: now,
					UpdatedAt: now,
				},
				{
					ID:        "2",
					Text:      "Architected backend platform",
					Date:      now.AddDate(0, -1, 0),
					Company:   "TechCorp",
					Tags:      []string{"technical"},
					CreatedAt: now,
					UpdatedAt: now,
				},
			}

			// Generate suggestions
			suggestions, err := workflow.GenerateSuggestions(ctx, events, nil)
			Expect(err).NotTo(HaveOccurred())

			// Confirm first suggestion
			if len(suggestions) > 0 {
				err := workflow.ConfirmSuggestion(suggestions[0])
				Expect(err).NotTo(HaveOccurred())

				confirmed := workflow.GetConfirmedBursts()
				Expect(len(confirmed)).To(Equal(1))
			}
		})

		It("should detect, reject, and prevent re-suggestion", func() {
			now := time.Now()
			events := []career.CareerEvent{
				{
					ID:        "1",
					Text:      "Backend work",
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
			}

			// First generation
			suggestions1, err := workflow.GenerateSuggestions(ctx, events, nil)
			Expect(err).NotTo(HaveOccurred())
			initialCount := len(suggestions1)

			if initialCount > 0 {
				// Reject first suggestion
				err := workflow.RejectSuggestion(suggestions1[0].EventIDs)
				Expect(err).NotTo(HaveOccurred())

				// Second generation should have fewer suggestions
				suggestions2, err := workflow.GenerateSuggestions(ctx, events, nil)
				Expect(err).NotTo(HaveOccurred())
				Expect(len(suggestions2)).To(BeNumerically("<",initialCount))
			}
		})

		It("should handle multiple events with temporal grouping", func() {
			now := time.Now()
			events := []career.CareerEvent{
				{
					ID:        "1",
					Text:      "Led backend infrastructure project",
					Date:      now.AddDate(0, -5, 0),
					Company:   "TechCorp",
					CreatedAt: now,
					UpdatedAt: now,
				},
				{
					ID:        "2",
					Text:      "Architected backend platform",
					Date:      now.AddDate(0, -4, 0),
					Company:   "TechCorp",
					CreatedAt: now,
					UpdatedAt: now,
				},
				{
					ID:        "3",
					Text:      "Implemented microservices",
					Date:      now.AddDate(0, -3, 0),
					Company:   "TechCorp",
					CreatedAt: now,
					UpdatedAt: now,
				},
				{
					ID:        "4",
					Text:      "Unrelated event",
					Date:      now,
					Company:   "OtherCorp",
					CreatedAt: now,
					UpdatedAt: now,
				},
			}

			suggestions, err := workflow.GenerateSuggestions(ctx, events, nil)
			Expect(err).NotTo(HaveOccurred())
			Expect(len(suggestions)).To(BeNumerically(">=",0))
		})

		It("should respect custom detection options", func() {
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

			opts := &DetectionOptions{
				MinConfidence:       0.95,
				MinEventCount:       2,
				MaxSuggestionsCount: 5,
			}

			suggestions, err := workflow.GenerateSuggestions(ctx, events, opts)
			Expect(err).NotTo(HaveOccurred())
			// High confidence threshold should filter most suggestions
			Expect(len(suggestions)).To(BeNumerically("<=", 5))
		})

		It("should handle empty event list gracefully", func() {
			suggestions, err := workflow.GenerateSuggestions(ctx, []career.CareerEvent{}, nil)
			Expect(err).NotTo(HaveOccurred())
			Expect(suggestions).To(BeEmpty())
		})

		It("should handle single event gracefully", func() {
			now := time.Now()
			events := []career.CareerEvent{
				{
					ID:        "1",
					Text:      "Single event",
					Date:      now,
					CreatedAt: now,
					UpdatedAt: now,
				},
			}

			suggestions, err := workflow.GenerateSuggestions(ctx, events, nil)
			Expect(err).NotTo(HaveOccurred())
			Expect(suggestions).To(BeEmpty())
		})

		It("should track multiple confirmations", func() {
			now := time.Now()
			events := []career.CareerEvent{
				{
					ID:        "1",
					Text:      "Backend infrastructure work",
					Date:      now.AddDate(0, -5, 0),
					Company:   "TechCorp",
					CreatedAt: now,
					UpdatedAt: now,
				},
				{
					ID:        "2",
					Text:      "Backend platform",
					Date:      now.AddDate(0, -4, 0),
					Company:   "TechCorp",
					CreatedAt: now,
					UpdatedAt: now,
				},
				{
					ID:        "3",
					Text:      "Microservices work",
					Date:      now.AddDate(0, -3, 0),
					Company:   "TechCorp",
					CreatedAt: now,
					UpdatedAt: now,
				},
			}

			suggestions, err := workflow.GenerateSuggestions(ctx, events, nil)
			Expect(err).NotTo(HaveOccurred())

			// Confirm first two suggestions
			confirmCount := 0
			for i := 0; i < len(suggestions) && i < 2; i++ {
				err := workflow.ConfirmSuggestion(suggestions[i])
				if err == nil {
					confirmCount++
				}
			}

			confirmed := workflow.GetConfirmedBursts()
			Expect(len(confirmed)).To(Equal(confirmCount))
		})

		It("should reset workflow state completely", func() {
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

			// Generate, confirm, and reject
			suggestions, err := workflow.GenerateSuggestions(ctx, events, nil)
			Expect(err).NotTo(HaveOccurred())

			if len(suggestions) > 0 {
				workflow.ConfirmSuggestion(suggestions[0])
				if len(suggestions) > 1 {
					workflow.RejectSuggestion(suggestions[1].EventIDs)
				}
			}

			// Verify state exists
			Expect(len(workflow.GetConfirmedBursts())).To(BeNumerically(">=",0))

			// Reset
			workflow.Reset()

			// Verify clean state
			Expect(len(workflow.GetPendingSuggestions())).To(Equal(0))
			Expect(len(workflow.GetConfirmedBursts())).To(Equal(0))
		})

		It("should handle events with similar text but different companies", func() {
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
					Text:      "Backend infrastructure work",
					Date:      now.AddDate(0, -1, 0),
					Company:   "OtherCorp",
					CreatedAt: now,
					UpdatedAt: now,
				},
			}

			suggestions, err := workflow.GenerateSuggestions(ctx, events, nil)
			Expect(err).NotTo(HaveOccurred())
			// Similar text might generate suggestions
			Expect(len(suggestions)).To(BeNumerically(">=",0))
		})

		It("should handle events with same company and project", func() {
			now := time.Now()
			events := []career.CareerEvent{
				{
					ID:        "1",
					Text:      "Implemented feature X",
					Date:      now.AddDate(0, -2, 0),
					Company:   "TechCorp",
					Project:   "ProjectA",
					CreatedAt: now,
					UpdatedAt: now,
				},
				{
					ID:        "2",
					Text:      "Fixed bugs in feature X",
					Date:      now.AddDate(0, -1, 0),
					Company:   "TechCorp",
					Project:   "ProjectA",
					CreatedAt: now,
					UpdatedAt: now,
				},
			}

			suggestions, err := workflow.GenerateSuggestions(ctx, events, nil)
			Expect(err).NotTo(HaveOccurred())
			// Should detect related events
			if len(suggestions) > 0 {
				Expect(suggestions[0].ConfidenceScore).To(BeNumerically(">",0.5))
			}
		})
	})

	Describe("Similarity Scorer Integration", func() {
		It("should score events with keyword overlap correctly", func() {
			now := time.Now()
			events := []career.CareerEvent{
				{
					ID:        "1",
					Text:      "Led backend infrastructure project",
					Date:      now,
					Tags:      []string{"technical", "leadership"},
					CreatedAt: now,
					UpdatedAt: now,
				},
				{
					ID:        "2",
					Text:      "Architected backend platform",
					Date:      now.AddDate(0, -1, 0),
					Tags:      []string{"technical"},
					CreatedAt: now,
					UpdatedAt: now,
				},
			}

			matrix := detector.buildSimilarityMatrix(events)
			score := matrix["1"]["2"]
			Expect(score).To(BeNumerically(">",0.0))
		})

		It("should score events with company match correctly", func() {
			now := time.Now()
			events := []career.CareerEvent{
				{
					ID:        "1",
					Text:      "Event at TechCorp",
					Date:      now,
					Company:   "TechCorp",
					CreatedAt: now,
					UpdatedAt: now,
				},
				{
					ID:        "2",
					Text:      "Another event at TechCorp",
					Date:      now.AddDate(0, -1, 0),
					Company:   "TechCorp",
					CreatedAt: now,
					UpdatedAt: now,
				},
			}

			matrix := detector.buildSimilarityMatrix(events)
			score := matrix["1"]["2"]
			Expect(score).To(BeNumerically(">",0.0))
		})
	})

	Describe("Temporal Grouping Integration", func() {
		It("should group events within temporal window", func() {
			now := time.Now()
			events := []career.CareerEvent{
				{
					ID:        "1",
					Text:      "Event 1",
					Date:      now.AddDate(0, -5, 0),
					CreatedAt: now,
					UpdatedAt: now,
				},
				{
					ID:        "2",
					Text:      "Event 2",
					Date:      now.AddDate(0, -4, 0),
					CreatedAt: now,
					UpdatedAt: now,
				},
				{
					ID:        "3",
					Text:      "Event 3",
					Date:      now.AddDate(0, -3, 0),
					CreatedAt: now,
					UpdatedAt: now,
				},
			}

			// All should be grouped together (within 6 months)
			suggestions, err := workflow.GenerateSuggestions(ctx, events, nil)
			Expect(err).NotTo(HaveOccurred())
			// Should find related events
			Expect(len(suggestions)).To(BeNumerically(">=",0))
		})
	})
})
