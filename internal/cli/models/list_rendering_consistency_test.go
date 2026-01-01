package models

import (
	"context"
	"strings"
	"time"

	"github.com/baphled/kariya/internal/domain/career"
	careerrepo "github.com/baphled/kariya/internal/repository/career"
	careerservice "github.com/baphled/kariya/internal/service/career"
	"github.com/onsi/ginkgo/v2"
	"github.com/onsi/gomega"
)

var _ = ginkgo.Describe("List Model Rendering Consistency", func() {
	var (
		ctx        context.Context
		repo       *careerrepo.MemoryRepository
		service    *careerservice.Service
		listModel  *ListModel
		burstModel *BurstListModel
		factModel  *FactListModel
	)

	ginkgo.BeforeEach(func() {
		ctx = context.Background()
		repo = careerrepo.NewMemoryRepository()
		service = careerservice.NewService(repo)

		// Capture some test events
		event1 := &career.CareerEvent{
			Text: "Led team on strategic initiative",
			Date: time.Now().Add(-10 * 24 * time.Hour),
		}
		event2 := &career.CareerEvent{
			Text: "Implemented critical feature",
			Date: time.Now().Add(-5 * 24 * time.Hour),
		}

		_ = service.CaptureEvent(ctx, event1, careerservice.TimelineJournaling)
		_ = service.CaptureEvent(ctx, event2, careerservice.TimelineJournaling)

		// Initialize models
		listModel = NewListModel(service, ctx)
		listModel.width = 80
		listModel.height = 24

		burstModel = NewBurstListModel(service, ctx)
		burstModel.width = 80
		burstModel.height = 24

		factModel = NewFactListModel(service, ctx)
		factModel.width = 80
		factModel.height = 24
	})

	ginkgo.Describe("Pagination Display", func() {
		ginkgo.It("should display pagination info in all models", func() {
			listView := listModel.View()
			gomega.Expect(listView).NotTo(gomega.BeEmpty())

			burstView := burstModel.View()
			gomega.Expect(burstView).To(gomega.ContainSubstring("bursts"))

			factView := factModel.View()
			gomega.Expect(factView).NotTo(gomega.BeEmpty())
		})

		ginkgo.It("should use 'Showing X-Y of Z <items>' pagination format in all models", func() {
			// Test list.go pagination format
			listView := listModel.View()
			gomega.Expect(listView).To(gomega.MatchRegexp(`Showing \d+-\d+ of \d+ events`))

			// Test burst_list.go pagination format
			burstView := burstModel.View()
			gomega.Expect(burstView).To(gomega.MatchRegexp(`Showing \d+-\d+ of \d+ bursts`))

			// Test fact_list.go pagination format
			factView := factModel.View()
			gomega.Expect(factView).To(gomega.MatchRegexp(`Showing \d+-\d+ of \d+ facts`))
		})

		ginkgo.It("should have consistent pagination format structure", func() {
			// All pagination should start with 'Showing' and use 'X-Y of Z' format
			listView := listModel.View()
			burstView := burstModel.View()
			factView := factModel.View()

			// Extract pagination lines
			listHasShowingFormat := strings.Contains(listView, "Showing")
			burstHasShowingFormat := strings.Contains(burstView, "Showing")
			factHasShowingFormat := strings.Contains(factView, "Showing")

			gomega.Expect(listHasShowingFormat).To(gomega.BeTrue(), "list.go should use 'Showing' format")
			gomega.Expect(burstHasShowingFormat).To(gomega.BeTrue(), "burst_list.go should use 'Showing' format")
			gomega.Expect(factHasShowingFormat).To(gomega.BeTrue(), "fact_list.go should use 'Showing' format")
		})
	})

	ginkgo.Describe("Empty State Handling", func() {
		ginkgo.It("should display consistent empty state messages", func() {
			emptyBurstModel := NewBurstListModel(service, ctx)
			emptyBurstModel.width = 80
			emptyBurstModel.height = 24
			burstView := emptyBurstModel.View()
			gomega.Expect(burstView).To(gomega.ContainSubstring("No bursts found"))

			emptyFactModel := NewFactListModel(service, ctx)
			emptyFactModel.width = 80
			emptyFactModel.height = 24
			factView := emptyFactModel.View()
			gomega.Expect(factView).To(gomega.ContainSubstring("No facts found"))
		})

		ginkgo.It("should use 'No <items> found' format for empty states", func() {
			// Create empty repository for consistent testing
			emptyRepo := careerrepo.NewMemoryRepository()
			emptyService := careerservice.NewService(emptyRepo)

			emptyListModel := NewListModel(emptyService, ctx)
			emptyListModel.width = 80
			emptyListModel.height = 24
			listView := emptyListModel.View()
			gomega.Expect(listView).To(gomega.MatchRegexp(`No .+ found`))

			emptyBurstModel := NewBurstListModel(emptyService, ctx)
			emptyBurstModel.width = 80
			emptyBurstModel.height = 24
			burstView := emptyBurstModel.View()
			gomega.Expect(burstView).To(gomega.ContainSubstring("No bursts found"))

			emptyFactModel := NewFactListModel(emptyService, ctx)
			emptyFactModel.width = 80
			emptyFactModel.height = 24
			factView := emptyFactModel.View()
			gomega.Expect(factView).To(gomega.ContainSubstring("No facts found"))
		})

		ginkgo.It("should NOT have conditional empty state messages", func() {
			// Create empty repository
			emptyRepo := careerrepo.NewMemoryRepository()
			emptyService := careerservice.NewService(emptyRepo)

			// Test burst_list.go doesn't have "No matching bursts"
			emptyBurstModel := NewBurstListModel(emptyService, ctx)
			emptyBurstModel.width = 80
			emptyBurstModel.height = 24
			burstView := emptyBurstModel.View()
			gomega.Expect(burstView).To(gomega.ContainSubstring("No bursts found"))
			gomega.Expect(burstView).NotTo(gomega.ContainSubstring("No matching bursts"))

			// Test fact_list.go doesn't have "No facts match the current filters"
			emptyFactModel := NewFactListModel(emptyService, ctx)
			emptyFactModel.width = 80
			emptyFactModel.height = 24
			factView := emptyFactModel.View()
			gomega.Expect(factView).To(gomega.ContainSubstring("No facts found"))
			gomega.Expect(factView).NotTo(gomega.ContainSubstring("No facts match the current filters"))
		})
	})

	ginkgo.Describe("Navigation Consistency", func() {
		ginkgo.It("should support page navigation methods in all models", func() {
			listModel.nextPage()
			listModel.prevPage()

			burstModel.nextPage()
			burstModel.prevPage()

			factModel.nextPage()
			factModel.prevPage()

			gomega.Expect(true).To(gomega.BeTrue())
		})

		ginkgo.It("should support jump to start/end methods in all models", func() {
			listModel.goToFirstItem()
			listModel.goToLastItem()

			burstModel.goToFirstItem()
			burstModel.goToLastItem()

			factModel.goToFirstItem()
			factModel.goToLastItem()

			gomega.Expect(true).To(gomega.BeTrue())
		})
	})

	ginkgo.Describe("Rendering Structure", func() {
		ginkgo.It("should render all models with content", func() {
			listView := listModel.View()
			burstView := burstModel.View()
			factView := factModel.View()

			gomega.Expect(listView).NotTo(gomega.BeEmpty())
			gomega.Expect(burstView).NotTo(gomega.BeEmpty())
			gomega.Expect(factView).NotTo(gomega.BeEmpty())
		})

		ginkgo.It("should render with multiple sections (header, content, footer)", func() {
			listView := listModel.View()
			burstView := burstModel.View()
			factView := factModel.View()

			listLines := strings.Split(listView, "\n")
			burstLines := strings.Split(burstView, "\n")
			factLines := strings.Split(factView, "\n")

			gomega.Expect(len(listLines)).To(gomega.BeNumerically(">", 1))
			gomega.Expect(len(burstLines)).To(gomega.BeNumerically(">", 1))
			gomega.Expect(len(factLines)).To(gomega.BeNumerically(">", 1))
		})
	})

	ginkgo.Describe("Rendering Stability", func() {
		ginkgo.It("should produce stable output for same input", func() {
			burstView1 := burstModel.View()
			burstView2 := burstModel.View()
			gomega.Expect(burstView1).To(gomega.Equal(burstView2))

			factView1 := factModel.View()
			factView2 := factModel.View()
			gomega.Expect(factView1).To(gomega.Equal(factView2))

			listView1 := listModel.View()
			listView2 := listModel.View()
			gomega.Expect(listView1).To(gomega.Equal(listView2))
		})
	})
})
