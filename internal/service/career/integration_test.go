package career

import (
	"context"
	"os"
	"path/filepath"
	"time"

	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"

	"github.com/baphled/kariya/internal/domain/career"
	repo "github.com/baphled/kariya/internal/repository/career"
	careersql "github.com/baphled/kariya/internal/repository/career/sql"
	"github.com/baphled/kariya/internal/testutil/fixtures"
)

var _ = Describe("Career Service Integration Tests", func() {
	var (
		service *Service
		ctx     context.Context
		tempDir string
		repos   *repo.Repositories
	)

	BeforeEach(func() {
		var err error
		tempDir, err = os.MkdirTemp("", "kariya-service-test-")
		Expect(err).NotTo(HaveOccurred())

		dbPath := filepath.Join(tempDir, "test_events.db")
		repos, err = careersql.NewRepositoriesFromPath(dbPath)
		Expect(err).NotTo(HaveOccurred())

		service = NewService(repos.Event)

		ctx = context.Background()
	})

	AfterEach(func() {
		if repos != nil {
			repos.Close()
		}
		os.RemoveAll(tempDir)
	})

	Describe("Event Capture Workflow", func() {
		It("should capture events with TimelineJournaling mode", func() {
			event := fixtures.EventWith("", "Developed a high-performance backend service", "Test Company", "")
			event.Date = time.Now().AddDate(0, 0, -10)
			event.Tags = []string{"technical", "project"}

			err := service.CaptureEvent(ctx, event, TimelineJournaling)
			Expect(err).NotTo(HaveOccurred())

			retrievedEvent, err := service.GetEventByID(ctx, event.ID)
			Expect(err).NotTo(HaveOccurred())
			Expect(retrievedEvent.Text).To(Equal(event.Text))
			Expect(retrievedEvent.Date).To(BeTemporally("~", event.Date, time.Second))
			Expect(retrievedEvent.Tags).To(Equal(event.Tags))
			Expect(retrievedEvent.Company).To(Equal(event.Company))
		})

		It("should capture events with CVBackfill mode", func() {
			event := fixtures.EventWith("", "Led major system redesign", "Old Company", "")
			event.Date = time.Now().AddDate(-2, 0, 0)
			event.Tags = []string{"leadership", "technical"}

			err := service.CaptureEvent(ctx, event, CVBackfill)
			Expect(err).NotTo(HaveOccurred())

			retrievedEvent, err := service.GetEventByID(ctx, event.ID)
			Expect(err).NotTo(HaveOccurred())
			Expect(retrievedEvent.Text).To(Equal(event.Text))
			Expect(retrievedEvent.Company).To(Equal(event.Company))
		})

		It("should capture events with ManualEntry mode", func() {
			event := fixtures.EventWith("", "Mentored junior developers on best practices", "Current Company", "")
			event.Date = time.Now()
			event.Tags = []string{"mentoring"}

			err := service.CaptureEvent(ctx, event, ManualEntry)
			Expect(err).NotTo(HaveOccurred())

			retrievedEvent, err := service.GetEventByID(ctx, event.ID)
			Expect(err).NotTo(HaveOccurred())
			Expect(retrievedEvent.Text).To(Equal(event.Text))
			Expect(retrievedEvent.Tags).To(Equal(event.Tags))
		})
	})

	Describe("Event Tag Handling", func() {
		It("should preserve tags for technical development event", func() {
			event := fixtures.EventWith("", "Developed a scalable microservices architecture using Go", "Tech Corp", "")
			event.Date = time.Now()
			event.Tags = []string{"technical"}

			err := service.CaptureEvent(ctx, event, ManualEntry)
			Expect(err).NotTo(HaveOccurred())

			retrievedEvent, err := service.GetEventByID(ctx, event.ID)
			Expect(err).NotTo(HaveOccurred())
			Expect(retrievedEvent.Tags).To(Equal([]string{"technical"}))
		})

		It("should preserve tags for leadership event", func() {
			event := fixtures.EventWith("", "Led a cross-functional team to deliver a critical project", "Tech Corp", "")
			event.Date = time.Now()
			event.Tags = []string{"leadership"}

			err := service.CaptureEvent(ctx, event, ManualEntry)
			Expect(err).NotTo(HaveOccurred())

			retrievedEvent, err := service.GetEventByID(ctx, event.ID)
			Expect(err).NotTo(HaveOccurred())
			Expect(retrievedEvent.Tags).To(Equal([]string{"leadership"}))
		})
	})

	Describe("Event Filtering", func() {
		BeforeEach(func() {
			e1 := fixtures.EventWith("", "Developed backend service in Go", "Company A", "")
			e1.Date = time.Now().AddDate(0, 0, -30)
			e1.Tags = []string{"technical", "project"}
			e2 := fixtures.EventWith("", "Led team strategy meeting", "Company B", "")
			e2.Date = time.Now().AddDate(0, 0, -15)
			e2.Tags = []string{"leadership"}
			e3 := fixtures.EventWith("", "Consulted on cloud migration", "Company C", "")
			e3.Date = time.Now().AddDate(0, 0, -45)
			e3.Tags = []string{"consulting"}

			for _, event := range []*career.Event{e1, e2, e3} {
				err := service.CaptureEvent(ctx, event, ManualEntry)
				Expect(err).NotTo(HaveOccurred())
			}
		})

		It("should filter events by tag", func() {
			events, err := service.ListEvents(ctx, *fixtures.EventListFiltersWithTags([]string{"technical"}))
			Expect(err).NotTo(HaveOccurred())
			Expect(events).To(HaveLen(1))
			Expect(events[0].Tags).To(ContainElement("technical"))
		})

		It("should filter events by date range", func() {
			startDate := time.Now().AddDate(0, 0, -20)
			endDate := time.Now()

			events, err := service.ListEvents(ctx, *fixtures.EventListFiltersWithDateRange(&startDate, &endDate))
			Expect(err).NotTo(HaveOccurred())
			Expect(events).To(HaveLen(1))
			Expect(events[0].Text).To(ContainSubstring("team strategy"))
		})

		It("should return all events when no filters applied", func() {
			events, err := service.ListEvents(ctx, *fixtures.EventListFilters())
			Expect(err).NotTo(HaveOccurred())
			Expect(events).To(HaveLen(3))
		})
	})

	Describe("Event Management Operations", func() {
		It("should update an existing event", func() {
			event := fixtures.EventWith("", "Initial event description", "Test Company", "")
			event.Date = time.Now()
			event.Tags = []string{"project"}

			err := service.CaptureEvent(ctx, event, ManualEntry)
			Expect(err).NotTo(HaveOccurred())

			event.Text = "Updated event description"
			err = service.UpdateEvent(ctx, event)
			Expect(err).NotTo(HaveOccurred())

			retrievedEvent, err := service.GetEventByID(ctx, event.ID)
			Expect(err).NotTo(HaveOccurred())
			Expect(retrievedEvent.Text).To(Equal("Updated event description"))
		})

		It("should delete an existing event", func() {
			event := fixtures.EventWith("", "Event to be deleted", "Test Company", "")
			event.Date = time.Now()
			event.Tags = []string{"project"}

			err := service.CaptureEvent(ctx, event, ManualEntry)
			Expect(err).NotTo(HaveOccurred())

			err = service.DeleteEvent(ctx, event.ID)
			Expect(err).NotTo(HaveOccurred())

			_, err = service.GetEventByID(ctx, event.ID)
			Expect(err).To(HaveOccurred())
		})

		It("should retrieve event by ID", func() {
			event := fixtures.EventWith("", "Test event for retrieval", "Test Company", "")
			event.Date = time.Now()
			event.Tags = []string{"technical"}

			err := service.CaptureEvent(ctx, event, ManualEntry)
			Expect(err).NotTo(HaveOccurred())

			retrievedEvent, err := service.GetEventByID(ctx, event.ID)
			Expect(err).NotTo(HaveOccurred())
			Expect(retrievedEvent.ID).To(Equal(event.ID))
			Expect(retrievedEvent.Text).To(Equal(event.Text))
		})

		It("should count events with filters", func() {
			e1 := fixtures.EventWith("", "Technical event", "Company A", "")
			e1.Date = time.Now()
			e1.Tags = []string{"technical"}
			e2 := fixtures.EventWith("", "Leadership event", "Company B", "")
			e2.Date = time.Now()
			e2.Tags = []string{"leadership"}

			for _, event := range []*career.Event{e1, e2} {
				err := service.CaptureEvent(ctx, event, ManualEntry)
				Expect(err).NotTo(HaveOccurred())
			}

			count, err := service.CountEvents(ctx, *fixtures.EventListFiltersWithTags([]string{"technical"}))
			Expect(err).NotTo(HaveOccurred())
			Expect(count).To(Equal(1))
		})
	})
})
