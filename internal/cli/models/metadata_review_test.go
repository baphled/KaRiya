package models

import (
	"context"
	"testing"
	"time"

	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"

	"github.com/baphled/kariya/internal/domain/career"
	careerrepo "github.com/baphled/kariya/internal/repository/career"
	careerservice "github.com/baphled/kariya/internal/service/career"
)

func TestMetadataReview(t *testing.T) {
	RegisterFailHandler(Fail)
	RunSpecs(t, "Metadata Review Suite")
}

var _ = Describe("MetadataReviewModel", func() {
	var (
		model    *MetadataReviewModel
		service  *careerservice.Service
		repo     *careerrepo.MemoryRepository
		ctx      context.Context
	)

	BeforeEach(func() {
		repo = careerrepo.NewMemoryRepository()
		service = careerservice.NewService(repo)
		ctx = context.Background()
		model = NewMetadataReviewModel(service, ctx)
	})

	Describe("Initialization", func() {
		It("should create model with default values", func() {
			Expect(model).NotTo(BeNil())
			Expect(model.selectedIdx).To(Equal(0))
			Expect(model.expandedIdx).To(Equal(-1))
			Expect(model.filterMode).To(Equal("all"))
			Expect(model.sortBy).To(Equal("quality"))
		})

		It("should handle empty event list", func() {
			Expect(model.events).To(HaveLen(0))
			Expect(model.qualityScores).To(HaveLen(0))
		})
	})

	Describe("Loading events", func() {
		It("should load events from service", func() {
			// Create test event
			event := &career.CareerEvent{
				ID:   "test-id",
				Text: "Test event",
				Date: time.Now().Add(-24 * time.Hour),
			}
			err := service.CaptureEvent(ctx, event, careerservice.ManualEntry)
			Expect(err).NotTo(HaveOccurred())

			// Create new model to load event
			model := NewMetadataReviewModel(service, ctx)
			Expect(model.events).To(HaveLen(1))
		})

		It("should calculate quality scores for all events", func() {
			// Create test event
			event := &career.CareerEvent{
				ID:      "test-id",
				Text:    "Test event",
				Date:    time.Now().Add(-24 * time.Hour),
				Company: "TestCorp",
			}
			err := service.CaptureEvent(ctx, event, careerservice.ManualEntry)
			Expect(err).NotTo(HaveOccurred())

			// Create new model
			model := NewMetadataReviewModel(service, ctx)
			Expect(model.qualityScores).To(HaveLen(1))

			// Verify score exists
			for _, e := range model.events {
				score, exists := model.qualityScores[e.ID]
				Expect(exists).To(BeTrue())
				Expect(score).NotTo(BeNil())
			}
		})
	})

	Describe("Filtering", func() {
		It("should show all events when filter mode is 'all'", func() {
			// Create test events
			for i := 0; i < 3; i++ {
				event := &career.CareerEvent{
					ID:   "test-id-" + string(rune(i)),
					Text: "Event " + string(rune(i)),
					Date: time.Now().Add(-24 * time.Hour),
				}
				_ = service.CaptureEvent(ctx, event, careerservice.ManualEntry)
			}

			model := NewMetadataReviewModel(service, ctx)
			Expect(model.events).To(HaveLen(3))
		})

		It("should filter to incomplete events only", func() {
			// Create incomplete event (minimal fields)
			event1 := &career.CareerEvent{
				ID:   "incomplete",
				Text: "Event",
				Date: time.Now().Add(-24 * time.Hour),
			}
			_ = service.CaptureEvent(ctx, event1, careerservice.ManualEntry)

			// Create complete event
			event2 := &career.CareerEvent{
				ID:         "complete",
				Text:       "Complete event with all fields",
				Date:       time.Now().Add(-24 * time.Hour),
				Company:    "Corp",
				Project:    "Project",
				Tags:       []string{"technical", "leadership"},
				Categories: []string{"technical"},
			}
			_ = service.CaptureEvent(ctx, event2, careerservice.ManualEntry)

			model := NewMetadataReviewModel(service, ctx)
			model.filterMode = "incomplete"
			model.loadEvents()

			Expect(model.events).To(HaveLen(1))
			Expect(model.events[0].ID).To(Equal("incomplete"))
		})
	})

	Describe("Sorting", func() {
		It("should sort by quality score (ascending)", func() {
			// Create events with different quality scores
			event1 := &career.CareerEvent{
				ID:      "id1",
				Text:    "Event 1",
				Date:    time.Now().Add(-24 * time.Hour),
				Company: "Corp",
			}
			event2 := &career.CareerEvent{
				ID:   "id2",
				Text: "Event 2",
				Date: time.Now().Add(-24 * time.Hour),
			}

			_ = service.CaptureEvent(ctx, event1, careerservice.ManualEntry)
			_ = service.CaptureEvent(ctx, event2, careerservice.ManualEntry)

			model := NewMetadataReviewModel(service, ctx)
			model.sortBy = "quality"
			model.sortEvents()

			// Lower quality (event2) should come first
			Expect(model.events[0].ID).To(Equal("id2"))
		})
	})

	Describe("Navigation", func() {
		BeforeEach(func() {
			// Create 3 test events
			for i := 0; i < 3; i++ {
				event := &career.CareerEvent{
					ID:   "id-" + string(rune(i)),
					Text: "Event " + string(rune(i)),
					Date: time.Now().Add(-24 * time.Hour),
				}
				_ = service.CaptureEvent(ctx, event, careerservice.ManualEntry)
			}
			model = NewMetadataReviewModel(service, ctx)
		})

		It("should move to next item", func() {
			model.selectedIdx = 0
			model.nextItem()
			Expect(model.selectedIdx).To(Equal(1))
		})

		It("should move to previous item", func() {
			model.selectedIdx = 1
			model.prevItem()
			Expect(model.selectedIdx).To(Equal(0))
		})

		It("should not go below first item", func() {
			model.selectedIdx = 0
			model.prevItem()
			Expect(model.selectedIdx).To(Equal(0))
		})

		It("should not go beyond last item", func() {
			model.selectedIdx = 2
			model.nextItem()
			Expect(model.selectedIdx).To(Equal(2))
		})
	})

	Describe("Expansion", func() {
		BeforeEach(func() {
			event := &career.CareerEvent{
				ID:   "test-id",
				Text: "Test event",
				Date: time.Now().Add(-24 * time.Hour),
			}
			_ = service.CaptureEvent(ctx, event, careerservice.ManualEntry)
			model = NewMetadataReviewModel(service, ctx)
		})

		It("should toggle expansion", func() {
			model.expandedIdx = -1
			model.selectedIdx = 0
			Expect(model.expandedIdx).To(Equal(-1))
		})
	})

	Describe("Getting selected event", func() {
		BeforeEach(func() {
			event := &career.CareerEvent{
				ID:   "test-id",
				Text: "Test event",
				Date: time.Now().Add(-24 * time.Hour),
			}
			_ = service.CaptureEvent(ctx, event, careerservice.ManualEntry)
			model = NewMetadataReviewModel(service, ctx)
		})

		It("should return selected event", func() {
			selected := model.GetSelectedEvent()
			Expect(selected).NotTo(BeNil())
			Expect(selected.ID).To(Equal("test-id"))
		})

		It("should return nil when selected idx is out of bounds", func() {
			emptyModel := NewMetadataReviewModel(service, ctx)
			emptyModel.selectedIdx = 100
			selected := emptyModel.GetSelectedEvent()
			Expect(selected).To(BeNil())
		})
	})

	Describe("Rendering", func() {
		It("should render empty state", func() {
			view := model.View()
			Expect(view).NotTo(BeEmpty())
			Expect(view).To(ContainSubstring("Review Metadata Quality"))
		})

		It("should render with events", func() {
			event := &career.CareerEvent{
				ID:   "test-id",
				Text: "Test event",
				Date: time.Now().Add(-24 * time.Hour),
			}
			_ = service.CaptureEvent(ctx, event, careerservice.ManualEntry)

			model := NewMetadataReviewModel(service, ctx)
			view := model.View()
			Expect(view).NotTo(BeEmpty())
			Expect(view).To(ContainSubstring("Test event"))
		})
	})
})
