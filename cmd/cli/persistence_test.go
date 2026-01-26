//nolint:errcheck // Test file - error handling for test setup is not relevant.
package main

import (
	"context"
	"path/filepath"
	"time"

	"github.com/baphled/kariya/internal/cli/service"
	careerrepo "github.com/baphled/kariya/internal/repository/career"
	careerservice "github.com/baphled/kariya/internal/service/career"
	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
)

var _ = Describe("Data Persistence", func() {
	Context("when capturing events with SQLite", func() {
		It("should persist events to SQLite database", func() {
			// Create a temporary database
			tmpDir := GinkgoT().TempDir()
			dbPath := filepath.Join(tmpDir, "test-events.db")

			// Create SQLite repository
			repo, err := careerrepo.NewSQLiteRepository(dbPath)
			Expect(err).ToNot(HaveOccurred())
			DeferCleanup(func() {
				repo.Close()
			})

			// Create service
			svc := careerservice.NewService(repo)
			cliSvc := service.NewCLIEventService(svc)

			// Capture an event
			ctx := context.Background()
			eventText := "Led implementation of critical feature"
			eventDate := time.Now().Add(-24 * time.Hour)

			err = cliSvc.CaptureEvent(
				ctx,
				eventText,
				eventDate,
				careerservice.ManualEntry,
				service.WithCompany("TechCorp"),
				service.WithTags([]string{"technical", "leadership"}),
			)
			Expect(err).ToNot(HaveOccurred())

			// Verify event was persisted to database
			events, err := svc.ListEvents(ctx, careerrepo.ListFilters{})
			Expect(err).ToNot(HaveOccurred())
			Expect(events).To(HaveLen(1), "Event should be persisted to database")
			Expect(events[0].Text).To(Equal(eventText))
			Expect(events[0].Company).To(Equal("TechCorp"))
			Expect(events[0].Tags).To(ContainElements("technical", "leadership"))
		})

		It("should survive application restart with SQLite", func() {
			// Create a temporary database
			tmpDir := GinkgoT().TempDir()
			dbPath := filepath.Join(tmpDir, "persistent-events.db")

			// Create first repository instance and add event
			repo1, err := careerrepo.NewSQLiteRepository(dbPath)
			Expect(err).ToNot(HaveOccurred())
			DeferCleanup(func() {
				repo1.Close()
			})

			svc1 := careerservice.NewService(repo1)
			cliSvc1 := service.NewCLIEventService(svc1)

			ctx := context.Background()
			eventText := "Architected microservices platform"
			eventDate := time.Now().Add(-7 * 24 * time.Hour)

			err = cliSvc1.CaptureEvent(
				ctx,
				eventText,
				eventDate,
				careerservice.CVBackfill,
				service.WithCompany("StartupCo"),
			)
			Expect(err).ToNot(HaveOccurred())

			// Verify event was saved in first instance
			events1, err := svc1.ListEvents(ctx, careerrepo.ListFilters{})
			Expect(err).ToNot(HaveOccurred())
			Expect(events1).To(HaveLen(1))

			// Create second repository instance pointing to same database
			repo2, err := careerrepo.NewSQLiteRepository(dbPath)
			Expect(err).ToNot(HaveOccurred())
			DeferCleanup(func() {
				repo2.Close()
			})

			svc2 := careerservice.NewService(repo2)

			// Verify event persists across instances
			events2, err := svc2.ListEvents(ctx, careerrepo.ListFilters{})
			Expect(err).ToNot(HaveOccurred())
			Expect(events2).To(HaveLen(1), "Event should persist across database connections")
			Expect(events2[0].Text).To(Equal(eventText))
			Expect(events2[0].Company).To(Equal("StartupCo"))
		})

		It("should create database at custom path", func() {
			// Create a temporary database
			tmpDir := GinkgoT().TempDir()
			dbPath := filepath.Join(tmpDir, "custom-location.db")

			// Create SQLite repository at custom path
			repo, err := careerrepo.NewSQLiteRepository(dbPath)
			Expect(err).ToNot(HaveOccurred())
			DeferCleanup(func() {
				repo.Close()
			})

			// Verify database file was created
			Expect(dbPath).To(BeAnExistingFile())

			svc := careerservice.NewService(repo)
			cliSvc := service.NewCLIEventService(svc)

			ctx := context.Background()
			err = cliSvc.CaptureEvent(
				ctx,
				"Test event at custom path",
				time.Now().Add(-1*time.Hour),
				careerservice.ManualEntry,
			)
			Expect(err).ToNot(HaveOccurred())

			// Verify persistence
			events, err := svc.ListEvents(ctx, careerrepo.ListFilters{})
			Expect(err).ToNot(HaveOccurred())
			Expect(events).To(HaveLen(1))
		})
	})
})
