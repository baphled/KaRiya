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
)

var _ = Describe("Career Service Integration Tests", func() {
	var (
		service *Service
		ctx     context.Context
		tempDir string
		dbPath  string
		sqlRepo *repo.SQLiteRepository
	)

	BeforeEach(func() {
		var err error
		// Create a temporary directory for the test database
		tempDir, err = os.MkdirTemp("", "kariya-service-test-")
		Expect(err).NotTo(HaveOccurred())

		// Create SQLite repository
		dbPath = filepath.Join(tempDir, "test_events.db")
		sqlRepo, err = repo.NewSQLiteRepository(dbPath)
		Expect(err).NotTo(HaveOccurred())

		// Create service with the repository
		service = NewService(sqlRepo)

		ctx = context.Background()
	})

	AfterEach(func() {
		// Clean up: close repository and remove temporary files
		if sqlRepo != nil {
			sqlRepo.Close()
		}
		os.RemoveAll(tempDir)
	})

	Describe("Event Capture Workflow", func() {
		It("should capture events with TimelineJournaling mode", func() {
			event := &career.CareerEvent{
				Text:    "Developed a high-performance backend service",
				Date:    time.Now().AddDate(0, 0, -10), // 10 days ago
				Tags:    []string{"technical", "project"},
				Company: "Test Company",
			}

			err := service.CaptureEvent(ctx, event, TimelineJournaling)
			Expect(err).NotTo(HaveOccurred())

			retrievedEvent, err := service.GetEventByID(ctx, event.ID)
			Expect(err).NotTo(HaveOccurred())
			Expect(retrievedEvent.Text).To(Equal(event.Text))
			Expect(retrievedEvent.Date).To(Equal(event.Date))
			Expect(retrievedEvent.Tags).To(Equal(event.Tags))
			Expect(retrievedEvent.Company).To(Equal(event.Company))
		})

		It("should capture events with CVBackfill mode", func() {
			event := &career.CareerEvent{
				Text:    "Led major system redesign",
				Date:    time.Now().AddDate(-2, 0, 0), // 2 years ago
				Tags:    []string{"leadership", "technical"},
				Company: "Old Company",
			}

			err := service.CaptureEvent(ctx, event, CVBackfill)
			Expect(err).NotTo(HaveOccurred())

			retrievedEvent, err := service.GetEventByID(ctx, event.ID)
			Expect(err).NotTo(HaveOccurred())
			Expect(retrievedEvent.Text).To(Equal(event.Text))
			Expect(retrievedEvent.Company).To(Equal(event.Company))
		})

		It("should capture events with ManualEntry mode", func() {
			event := &career.CareerEvent{
				Text:    "Mentored junior developers on best practices",
				Date:    time.Now(),
				Tags:    []string{"mentoring"},
				Company: "Current Company",
			}

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
			event := &career.CareerEvent{
				Text:    "Developed a scalable microservices architecture using Go",
				Date:    time.Now(),
				Tags:    []string{"technical"},
				Company: "Tech Corp",
			}

			err := service.CaptureEvent(ctx, event, ManualEntry)
			Expect(err).NotTo(HaveOccurred())

			retrievedEvent, err := service.GetEventByID(ctx, event.ID)
			Expect(err).NotTo(HaveOccurred())
			Expect(retrievedEvent.Tags).To(Equal([]string{"technical"}))
		})

		It("should preserve tags for leadership event", func() {
			event := &career.CareerEvent{
				Text:    "Led a cross-functional team to deliver a critical project",
				Date:    time.Now(),
				Tags:    []string{"leadership"},
				Company: "Tech Corp",
			}

			err := service.CaptureEvent(ctx, event, ManualEntry)
			Expect(err).NotTo(HaveOccurred())

			retrievedEvent, err := service.GetEventByID(ctx, event.ID)
			Expect(err).NotTo(HaveOccurred())
			Expect(retrievedEvent.Tags).To(Equal([]string{"leadership"}))
		})
	})

	Describe("Event Filtering", func() {
		BeforeEach(func() {
			// Create multiple test events
			testEvents := []*career.CareerEvent{
				{
					Text:    "Developed backend service in Go",
					Date:    time.Now().AddDate(0, 0, -30),
					Tags:    []string{"technical", "project"},
					Company: "Company A",
				},
				{
					Text:    "Led team strategy meeting",
					Date:    time.Now().AddDate(0, 0, -15),
					Tags:    []string{"leadership"},
					Company: "Company B",
				},
				{
					Text:    "Consulted on cloud migration",
					Date:    time.Now().AddDate(0, 0, -45),
					Tags:    []string{"consulting"},
					Company: "Company C",
				},
			}

			// Capture all test events
			for _, event := range testEvents {
				err := service.CaptureEvent(ctx, event, ManualEntry)
				Expect(err).NotTo(HaveOccurred())
			}
		})

		It("should filter events by tag", func() {
			events, err := service.ListEvents(ctx, repo.ListFilters{
				Tags: []string{"technical"},
			})
			Expect(err).NotTo(HaveOccurred())
			Expect(events).To(HaveLen(1))
			Expect(events[0].Tags).To(ContainElement("technical"))
		})

		It("should filter events by date range", func() {
			startDate := time.Now().AddDate(0, 0, -20)
			endDate := time.Now()

			events, err := service.ListEvents(ctx, repo.ListFilters{
				StartDate: &startDate,
				EndDate:   &endDate,
			})
			Expect(err).NotTo(HaveOccurred())
			Expect(events).To(HaveLen(1))
			Expect(events[0].Text).To(ContainSubstring("team strategy"))
		})

		It("should return all events when no filters applied", func() {
			events, err := service.ListEvents(ctx, repo.ListFilters{})
			Expect(err).NotTo(HaveOccurred())
			Expect(events).To(HaveLen(3))
		})
	})

	Describe("Event Management Operations", func() {
		It("should update an existing event", func() {
			// Create initial event
			event := &career.CareerEvent{
				Text:    "Initial event description",
				Date:    time.Now(),
				Tags:    []string{"project"},
				Company: "Test Company",
			}

			err := service.CaptureEvent(ctx, event, ManualEntry)
			Expect(err).NotTo(HaveOccurred())

			// Update event
			event.Text = "Updated event description"
			err = service.UpdateEvent(ctx, event)
			Expect(err).NotTo(HaveOccurred())

			// Verify update
			retrievedEvent, err := service.GetEventByID(ctx, event.ID)
			Expect(err).NotTo(HaveOccurred())
			Expect(retrievedEvent.Text).To(Equal("Updated event description"))
		})

		It("should delete an existing event", func() {
			// Create event
			event := &career.CareerEvent{
				Text:    "Event to be deleted",
				Date:    time.Now(),
				Tags:    []string{"project"},
				Company: "Test Company",
			}

			err := service.CaptureEvent(ctx, event, ManualEntry)
			Expect(err).NotTo(HaveOccurred())

			// Delete event
			err = service.DeleteEvent(ctx, event.ID)
			Expect(err).NotTo(HaveOccurred())

			// Verify deletion
			_, err = service.GetEventByID(ctx, event.ID)
			Expect(err).To(HaveOccurred())
		})

		It("should retrieve event by ID", func() {
			event := &career.CareerEvent{
				Text:    "Test event for retrieval",
				Date:    time.Now(),
				Tags:    []string{"technical"},
				Company: "Test Company",
			}

			err := service.CaptureEvent(ctx, event, ManualEntry)
			Expect(err).NotTo(HaveOccurred())

			retrievedEvent, err := service.GetEventByID(ctx, event.ID)
			Expect(err).NotTo(HaveOccurred())
			Expect(retrievedEvent.ID).To(Equal(event.ID))
			Expect(retrievedEvent.Text).To(Equal(event.Text))
		})

		It("should count events with filters", func() {
			// Create multiple events
			events := []*career.CareerEvent{
				{
					Text:    "Technical event",
					Date:    time.Now(),
					Tags:    []string{"technical"},
					Company: "Company A",
				},
				{
					Text:    "Leadership event",
					Date:    time.Now(),
					Tags:    []string{"leadership"},
					Company: "Company B",
				},
			}

			for _, event := range events {
				err := service.CaptureEvent(ctx, event, ManualEntry)
				Expect(err).NotTo(HaveOccurred())
			}

			count, err := service.CountEvents(ctx, repo.ListFilters{
				Tags: []string{"technical"},
			})
			Expect(err).NotTo(HaveOccurred())
			Expect(count).To(Equal(1))
		})
	})
})
