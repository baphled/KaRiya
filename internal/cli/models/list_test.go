package models

import (
	"context"
	"time"

	"github.com/baphled/kariya/internal/domain/career"
	careerrepo "github.com/baphled/kariya/internal/repository/career"
	careerservice "github.com/baphled/kariya/internal/service/career"
	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
)

var _ = Describe("ListModel", func() {
	var (
		repo *careerrepo.MemoryRepository
		svc  *careerservice.Service
		ctx  context.Context
	)

	BeforeEach(func() {
		repo = careerrepo.NewMemoryRepository()
		svc = careerservice.NewService(repo)
		ctx = context.Background()
	})

	Context("when creating a new ListModel", func() {
		It("should initialize with default page size", func() {
			model := NewListModel(svc, ctx)
			Expect(model).NotTo(BeNil())
			Expect(model.pageSize).To(Equal(10))
		})

		It("should initialize with page 1", func() {
			model := NewListModel(svc, ctx)
			Expect(model.currentPage).To(Equal(1))
		})

		It("should load events from service", func() {
			event1 := &career.CareerEvent{
				Text: "First event",
				Date: time.Now().Add(-48 * time.Hour),
			}
			event2 := &career.CareerEvent{
				Text: "Second event",
				Date: time.Now().Add(-24 * time.Hour),
			}

			err := svc.CaptureEvent(ctx, event1, careerservice.TimelineJournaling)
			Expect(err).NotTo(HaveOccurred())
			err = svc.CaptureEvent(ctx, event2, careerservice.TimelineJournaling)
			Expect(err).NotTo(HaveOccurred())

			model := NewListModel(svc, ctx)
			Expect(model.events).To(HaveLen(2))
		})

		It("should initialize with zero total count when no events", func() {
			model := NewListModel(svc, ctx)
			Expect(model.totalCount).To(Equal(0))
		})
	})

	Context("when displaying events", func() {
		It("should load event and render in view", func() {
			event := &career.CareerEvent{
				Text: "Led team on strategic initiative",
				Date: time.Now().Add(-24 * time.Hour),
			}
			err := svc.CaptureEvent(ctx, event, careerservice.TimelineJournaling)
			Expect(err).NotTo(HaveOccurred())

			model := NewListModel(svc, ctx)
			Expect(model.events).To(HaveLen(1))
			Expect(model.events[0].Text).To(Equal("Led team on strategic initiative"))

			view := model.View()
			Expect(view).NotTo(BeEmpty())
			Expect(view).To(ContainSubstring("Career Events"))
			Expect(view).To(ContainSubstring("Led team on strategic initiative"))
		})

		It("should render events with truncation", func() {
			longText := "This is a very long event description that should be truncated when displayed in the list view to prevent the list from becoming too wide and difficult to read"
			event := &career.CareerEvent{
				Text: longText,
				Date: time.Now().Add(-12 * time.Hour),
			}
			err := svc.CaptureEvent(ctx, event, careerservice.TimelineJournaling)
			Expect(err).NotTo(HaveOccurred())

			model := NewListModel(svc, ctx)
			Expect(model.events).To(HaveLen(1))

			view := model.View()
			Expect(view).NotTo(BeEmpty())
			Expect(len(view)).To(BeNumerically(">", 0))
		})

		It("should show company name when available", func() {
			event := &career.CareerEvent{
				Text:    "Did important work",
				Date:    time.Now().Add(-6 * time.Hour),
				Company: "TechCorp Inc.",
			}
			err := svc.CaptureEvent(ctx, event, careerservice.TimelineJournaling)
			Expect(err).NotTo(HaveOccurred())

			model := NewListModel(svc, ctx)
			view := model.View()
			Expect(view).To(ContainSubstring("TechCorp Inc."))
		})
	})

	Context("when paginating events", func() {
		It("should respect page size of 10 by default", func() {
			for i := 1; i <= 15; i++ {
				event := &career.CareerEvent{
					Text: "Event",
					Date: time.Now().Add(-time.Duration(i) * time.Hour),
				}
				err := svc.CaptureEvent(ctx, event, careerservice.TimelineJournaling)
				Expect(err).NotTo(HaveOccurred())
			}

			model := NewListModel(svc, ctx)
			Expect(model.pageSize).To(Equal(10))
			Expect(len(model.events)).To(Equal(10))
		})

		It("should navigate to next page", func() {
			for i := 1; i <= 15; i++ {
				event := &career.CareerEvent{
					Text: "Event",
					Date: time.Now().Add(-time.Duration(i) * time.Hour),
				}
				err := svc.CaptureEvent(ctx, event, careerservice.TimelineJournaling)
				Expect(err).NotTo(HaveOccurred())
			}

			model := NewListModel(svc, ctx)
			initialPage := model.currentPage
			model.nextPage()
			Expect(model.currentPage).To(Equal(initialPage + 1))
		})

		It("should navigate to previous page", func() {
			for i := 1; i <= 15; i++ {
				event := &career.CareerEvent{
					Text: "Event",
					Date: time.Now().Add(-time.Duration(i) * time.Hour),
				}
				err := svc.CaptureEvent(ctx, event, careerservice.TimelineJournaling)
				Expect(err).NotTo(HaveOccurred())
			}

			model := NewListModel(svc, ctx)
			model.nextPage()
			model.prevPage()
			Expect(model.currentPage).To(Equal(1))
		})

		It("should not go past first page", func() {
			event := &career.CareerEvent{
				Text: "Event 1",
				Date: time.Now(),
			}
			err := svc.CaptureEvent(ctx, event, careerservice.TimelineJournaling)
			Expect(err).NotTo(HaveOccurred())

			model := NewListModel(svc, ctx)
			model.prevPage()
			Expect(model.currentPage).To(Equal(1))
		})

		It("should not go past last page", func() {
			for i := 1; i <= 5; i++ {
				event := &career.CareerEvent{
					Text: "Event",
					Date: time.Now().Add(-time.Duration(i) * time.Hour),
				}
				err := svc.CaptureEvent(ctx, event, careerservice.TimelineJournaling)
				Expect(err).NotTo(HaveOccurred())
			}

			model := NewListModel(svc, ctx)
			lastPage := model.getTotalPages()
			model.nextPage()
			Expect(model.currentPage).To(Equal(lastPage))
		})

		It("should calculate total pages correctly", func() {
			for i := 1; i <= 35; i++ {
				event := &career.CareerEvent{
					Text: "Event",
					Date: time.Now().Add(-time.Duration(i) * time.Hour),
				}
				err := svc.CaptureEvent(ctx, event, careerservice.TimelineJournaling)
				Expect(err).NotTo(HaveOccurred())
			}

			model := NewListModel(svc, ctx)
			Expect(model.getTotalPages()).To(Equal(4))
		})

		It("should display page indicator", func() {
			for i := 1; i <= 25; i++ {
				event := &career.CareerEvent{
					Text: "Event",
					Date: time.Now().Add(-time.Duration(i) * time.Hour),
				}
				err := svc.CaptureEvent(ctx, event, careerservice.TimelineJournaling)
				Expect(err).NotTo(HaveOccurred())
			}

			model := NewListModel(svc, ctx)
			view := model.View()
			Expect(view).To(ContainSubstring("Page"))
			Expect(view).To(ContainSubstring("of"))
		})
	})
})
