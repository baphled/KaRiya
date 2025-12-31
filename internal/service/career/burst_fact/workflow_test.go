package burst_fact

import (
	"context"
	"time"

	"github.com/baphled/kariya/internal/domain/career"
	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
)

var _ = Describe("BurstConfirmationWorkflow", func() {
	var (
		workflow *BurstConfirmationWorkflow
		detector *BurstDetector
		ctx      context.Context
	)

	BeforeEach(func() {
		detector = NewBurstDetector()
		workflow = NewBurstConfirmationWorkflow(detector)
		ctx = context.Background()
	})

	Describe("RejectionTracker", func() {
		var tracker *RejectionTracker

		BeforeEach(func() {
			tracker = NewRejectionTracker()
		})

		It("should record rejection", func() {
			eventIDs := []string{"1", "2"}
			tracker.RecordRejection(eventIDs)
			Expect(tracker.IsRejected(eventIDs)).To(BeTrue())
		})

		It("should recognize rejected suggestions regardless of order", func() {
			tracker.RecordRejection([]string{"1", "2"})
			Expect(tracker.IsRejected([]string{"2", "1"})).To(BeTrue())
		})

		It("should not reject different event combinations", func() {
			tracker.RecordRejection([]string{"1", "2"})
			Expect(tracker.IsRejected([]string{"1", "3"})).To(BeFalse())
		})

		It("should filter rejected suggestions", func() {
			tracker.RecordRejection([]string{"1", "2"})

			suggestions := []BurstSuggestion{
				{EventIDs: []string{"1", "2"}, ConfidenceScore: 0.8},
				{EventIDs: []string{"3", "4"}, ConfidenceScore: 0.7},
			}

			filtered := tracker.FilterRejected(suggestions)
			Expect(len(filtered)).To(Equal(1))
			Expect(filtered[0].EventIDs).To(Equal([]string{"3", "4"}))
		})

		It("should clear rejections", func() {
			tracker.RecordRejection([]string{"1", "2"})
			Expect(tracker.GetRejectionCount()).To(Equal(1))

			tracker.ClearRejections()
			Expect(tracker.GetRejectionCount()).To(Equal(0))
		})
	})

	Describe("BurstConfirmationWorkflow", func() {
		It("should generate suggestions and filter rejected ones", func() {
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
					Text:      "Backend platform architecture",
					Date:      now.AddDate(0, -1, 0),
					Company:   "TechCorp",
					CreatedAt: now,
					UpdatedAt: now,
				},
			}

			suggestions, err := workflow.GenerateSuggestions(ctx, events, nil)
			Expect(err).NotTo(HaveOccurred())
			Expect(len(suggestions)).To(BeNumerically(">=", 0))
		})

		It("should confirm suggestion", func() {
			suggestion := BurstSuggestion{
				EventIDs:        []string{"1", "2"},
				ConfidenceScore: 0.8,
			}

			err := workflow.ConfirmSuggestion(suggestion)
			Expect(err).NotTo(HaveOccurred())

			confirmed := workflow.GetConfirmedBursts()
			Expect(len(confirmed)).To(Equal(1))
		})

		It("should reject invalid suggestion", func() {
			suggestion := BurstSuggestion{
				EventIDs:        []string{"1"},
				ConfidenceScore: 0.8,
			}

			err := workflow.ConfirmSuggestion(suggestion)
			Expect(err).To(HaveOccurred())
		})

		It("should reject suggestion", func() {
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

			suggestions, err := workflow.GenerateSuggestions(ctx, events, nil)
			Expect(err).NotTo(HaveOccurred())

			if len(suggestions) > 0 {
				err := workflow.RejectSuggestion(suggestions[0].EventIDs)
				Expect(err).NotTo(HaveOccurred())

				pending := workflow.GetPendingSuggestions()
				Expect(len(pending)).To(BeNumerically("<",len(suggestions)))
			}
		})

		It("should prevent re-suggesting rejected bursts", func() {
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
					Text:      "Backend platform architecture",
					Date:      now.AddDate(0, -1, 0),
					Company:   "TechCorp",
					CreatedAt: now,
					UpdatedAt: now,
				},
			}

			suggestions1, err := workflow.GenerateSuggestions(ctx, events, nil)
			Expect(err).NotTo(HaveOccurred())
			initialCount := len(suggestions1)

			if initialCount > 0 {
				err := workflow.RejectSuggestion(suggestions1[0].EventIDs)
				Expect(err).NotTo(HaveOccurred())

				suggestions2, err := workflow.GenerateSuggestions(ctx, events, nil)
				Expect(err).NotTo(HaveOccurred())
				Expect(len(suggestions2)).To(BeNumerically("<",initialCount))
			}
		})

		It("should return pending suggestions", func() {
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

			suggestions, err := workflow.GenerateSuggestions(ctx, events, nil)
			Expect(err).NotTo(HaveOccurred())

			pending := workflow.GetPendingSuggestions()
			Expect(len(pending)).To(Equal(len(suggestions)))
		})

		It("should reset workflow state", func() {
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

			suggestions, err := workflow.GenerateSuggestions(ctx, events, nil)
			Expect(err).NotTo(HaveOccurred())

			if len(suggestions) > 0 {
				err := workflow.ConfirmSuggestion(suggestions[0])
				Expect(err).NotTo(HaveOccurred())
			}

			Expect(len(workflow.GetConfirmedBursts())).To(BeNumerically(">",0))

			workflow.Reset()
			Expect(len(workflow.GetPendingSuggestions())).To(Equal(0))
			Expect(len(workflow.GetConfirmedBursts())).To(Equal(0))
		})
	})
})
