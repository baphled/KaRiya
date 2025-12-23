package career

import (
	"context"
	"errors"
	"os"
	"path/filepath"
	"time"

	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"

	"github.com/baphled/kariya/internal/domain/career"
)

var _ = Describe("SQLite Repository", func() {
	var (
		repo    *SQLiteRepository
		ctx     context.Context
		tempDir string
		dbPath  string
	)

	BeforeEach(func() {
		var err error
		// Create a temporary directory for test database
		tempDir, err = os.MkdirTemp("", "kariya-sqlite-test-")
		Expect(err).NotTo(HaveOccurred())

		// Create database path
		dbPath = filepath.Join(tempDir, "test_events.db")

		// Create SQLite repository
		repo, err = NewSQLiteRepository(dbPath)
		Expect(err).NotTo(HaveOccurred())

		ctx = context.Background()
	})

	AfterEach(func() {
		// Clean up: close repository and remove temporary files
		if repo != nil {
			repo.Close()
		}
		os.RemoveAll(tempDir)
	})

	// Helper function to generate test career event
	createSQLiteTestEvent := func() *career.CareerEvent {
		return &career.CareerEvent{
			Text:    "Test Career Event",
			Date:    time.Now(),
			Tags:    []string{"project", "technical"},
			Company: "Test Company",
		}
	}

	Describe("Create", func() {
		It("should create an event and retrieve it successfully", func() {
			event := createSQLiteTestEvent()

			// Create event
			err := repo.Create(ctx, event)
			Expect(err).NotTo(HaveOccurred())

			// Retrieve event
			retrievedEvent, err := repo.GetByID(ctx, event.ID)
			Expect(err).NotTo(HaveOccurred())

			// Validate retrieved event
			Expect(retrievedEvent.Text).To(Equal(event.Text))
			Expect(retrievedEvent.Date.Unix()).To(Equal(event.Date.Unix()))
			Expect(retrievedEvent.Tags).To(Equal(event.Tags))
			Expect(retrievedEvent.Company).To(Equal(event.Company))
		})

		It("should prevent duplicate event creation", func() {
			event := createSQLiteTestEvent()
			event.ID = "fixed-id"

			// First creation should succeed
			err := repo.Create(ctx, event)
			Expect(err).NotTo(HaveOccurred())

			// Second creation with same ID should fail
			err = repo.Create(ctx, event)
			Expect(err).To(HaveOccurred())
			Expect(errors.Is(err, ErrDuplicateEvent)).To(BeTrue())
		})
	})

	Describe("Update", func() {
		It("should update event successfully", func() {
			event := createSQLiteTestEvent()

			// Create initial event
			err := repo.Create(ctx, event)
			Expect(err).NotTo(HaveOccurred())

			// Update event
			event.Text = "Updated Test Career Event"
			event.Tags = []string{"leadership", "project"}
			err = repo.Update(ctx, event)
			Expect(err).NotTo(HaveOccurred())

			// Retrieve updated event
			retrievedEvent, err := repo.GetByID(ctx, event.ID)
			Expect(err).NotTo(HaveOccurred())

			// Validate updates
			Expect(retrievedEvent.Text).To(Equal("Updated Test Career Event"))
			Expect(retrievedEvent.Tags).To(Equal([]string{"leadership", "project"}))
		})

		It("should return error when updating non-existent event", func() {
			event := createSQLiteTestEvent()
			event.ID = "non-existent-id"

			// Update non-existent event should fail
			err := repo.Update(ctx, event)
			Expect(err).To(HaveOccurred())
			Expect(errors.Is(err, ErrEventNotFound)).To(BeTrue())
		})
	})

	Describe("Delete", func() {
		It("should delete event successfully", func() {
			event := createSQLiteTestEvent()

			// Create event
			err := repo.Create(ctx, event)
			Expect(err).NotTo(HaveOccurred())

			// Delete event
			err = repo.Delete(ctx, event.ID)
			Expect(err).NotTo(HaveOccurred())

			// Retrieve should fail
			_, err = repo.GetByID(ctx, event.ID)
			Expect(err).To(HaveOccurred())
			Expect(errors.Is(err, ErrEventNotFound)).To(BeTrue())
		})

		It("should return error when deleting non-existent event", func() {
			// Delete non-existent event should fail
			err := repo.Delete(ctx, "non-existent-id")
			Expect(err).To(HaveOccurred())
			Expect(errors.Is(err, ErrEventNotFound)).To(BeTrue())
		})
	})

	Describe("List", func() {
		BeforeEach(func() {
			// Create multiple test events
			events := []*career.CareerEvent{
				{
					Text:    "Event 1",
					Date:    time.Now().AddDate(0, 0, -30),
					Tags:    []string{"project", "technical"},
					Company: "Company A",
				},
				{
					Text:    "Event 2",
					Date:    time.Now(),
					Tags:    []string{"leadership", "project"},
					Company: "Company B",
				},
			}

			for _, event := range events {
				err := repo.Create(ctx, event)
				Expect(err).NotTo(HaveOccurred())
			}
		})

		It("should filter events by tag", func() {
			// List events with filters
			listedEvents, err := repo.List(ctx, ListFilters{
				Tags: []string{"project"},
			})
			Expect(err).NotTo(HaveOccurred())

			// Validate number of returned events
			Expect(listedEvents).To(HaveLen(2),
				"Should return both events with project tag")
		})

		It("should filter events by date range", func() {
			startDate := time.Now().AddDate(0, 0, -15)

			// List events with filters
			listedEvents, err := repo.List(ctx, ListFilters{
				StartDate: &startDate,
			})
			Expect(err).NotTo(HaveOccurred())

			// Validate number of returned events
			Expect(listedEvents).To(HaveLen(1),
				"Should return only recent event within date range")
		})
	})

	Describe("Count", func() {
		BeforeEach(func() {
			// Create multiple test events
			events := []*career.CareerEvent{
				{
					Text:    "Event 1",
					Date:    time.Now().AddDate(0, 0, -30),
					Tags:    []string{"project", "technical"},
					Company: "Company A",
				},
				{
					Text:    "Event 2",
					Date:    time.Now(),
					Tags:    []string{"leadership", "project"},
					Company: "Company B",
				},
			}

			for _, event := range events {
				err := repo.Create(ctx, event)
				Expect(err).NotTo(HaveOccurred())
			}
		})

		It("should count all events", func() {
			// Count events with filters
			count, err := repo.Count(ctx, ListFilters{})
			Expect(err).NotTo(HaveOccurred())

			// Validate count of events
			Expect(count).To(Equal(2),
				"Should count all events")
		})

		It("should count events by tag", func() {
			// Count events with filters
			count, err := repo.Count(ctx, ListFilters{
				Tags: []string{"project"},
			})
			Expect(err).NotTo(HaveOccurred())

			// Validate count of events
			Expect(count).To(Equal(2),
				"Should count both events with project tag")
		})

		It("should count events by date range", func() {
			startDate := time.Now().AddDate(0, 0, -15)

			// Count events with filters
			count, err := repo.Count(ctx, ListFilters{
				StartDate: &startDate,
			})
			Expect(err).NotTo(HaveOccurred())

			// Validate count of events
			Expect(count).To(Equal(1),
				"Should count only event within recent date range")
		})
	})
})
