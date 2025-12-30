package app

import (
	"context"
	"os"
	"path/filepath"
	"time"

	"github.com/baphled/kariya/internal/cli/service"
	careerrepo "github.com/baphled/kariya/internal/repository/career"
	careerservice "github.com/baphled/kariya/internal/service/career"
	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
)

var _ = Describe("SQLite Persistence Integration", func() {
	var (
		tempDBPath string
		repo       careerrepo.Repository
		svc        *careerservice.Service
		cliService *service.CLIEventService
		ctx        context.Context
	)

	BeforeEach(func() {
		ctx = context.Background()
		// Create temporary SQLite database
		tempDir := os.TempDir()
		tempDBPath = filepath.Join(tempDir, "test_kariya.db")

		var err error
		repo, err = careerrepo.NewSQLiteRepository(tempDBPath)
		Expect(err).ToNot(HaveOccurred())

		svc = careerservice.NewService(repo)
		cliService = service.NewCLIEventService(svc)
	})

	AfterEach(func() {
		// Clean up temporary database
		if tempDBPath != "" {
			os.Remove(tempDBPath)
		}
	})

	Describe("Event Persistence to SQLite", func() {
		It("should persist event to SQLite and retrieve it", func() {
			eventText := "Implemented critical feature for production system"
			eventDate := time.Now().Add(-24 * time.Hour)

			// Capture event
			err := cliService.CaptureEvent(
				ctx,
				eventText,
				eventDate,
				careerservice.ManualEntry,
				service.WithCompany("TechCorp"),
				service.WithTags([]string{"technical", "product"}),
			)
			Expect(err).ToNot(HaveOccurred())

			// Verify event persisted by listing
			events, err := cliService.ListEvents(ctx, nil)
			Expect(err).ToNot(HaveOccurred())
			Expect(events).To(HaveLen(1))
			Expect(events[0].Text).To(Equal(eventText))
			Expect(events[0].Company).To(Equal("TechCorp"))
			Expect(events[0].Tags).To(ContainElements("technical", "product"))
		})

		It("should persist multiple events and retrieve them", func() {
			events := []struct {
				text    string
				company string
			}{
				{"Event 1", "Company A"},
				{"Event 2", "Company B"},
				{"Event 3", "Company C"},
			}

			// Capture multiple events
			for _, e := range events {
				err := cliService.CaptureEvent(
					ctx,
					e.text,
					time.Now().Add(-24*time.Hour),
					careerservice.ManualEntry,
					service.WithCompany(e.company),
				)
				Expect(err).ToNot(HaveOccurred())
			}

			// Verify all events persisted
			listedEvents, err := cliService.ListEvents(ctx, nil)
			Expect(err).ToNot(HaveOccurred())
			Expect(listedEvents).To(HaveLen(3))
		})

		It("should survive repository reinitialization", func() {
			eventText := "Critical production fix"

			// Capture event
			err := cliService.CaptureEvent(
				ctx,
				eventText,
				time.Now(),
				careerservice.ManualEntry,
			)
			Expect(err).ToNot(HaveOccurred())

			// Reinitialize repository (simulates app restart)
			repo, err = careerrepo.NewSQLiteRepository(tempDBPath)
			Expect(err).ToNot(HaveOccurred())
			svc = careerservice.NewService(repo)
			cliService = service.NewCLIEventService(svc)

			// Verify event still exists
			events, err := cliService.ListEvents(ctx, nil)
			Expect(err).ToNot(HaveOccurred())
			Expect(events).To(HaveLen(1))
			Expect(events[0].Text).To(Equal(eventText))
		})
	})
})
