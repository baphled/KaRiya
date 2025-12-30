package service_test

import (
	"context"

	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"

	"github.com/baphled/kariya/internal/cli/service"
	careerdom "github.com/baphled/kariya/internal/domain/career"
	careerrepo "github.com/baphled/kariya/internal/repository/career"
	careerservice "github.com/baphled/kariya/internal/service/career"
)


var _ = Describe("CLI Event Service - Bulk Operations", func() {
	var (
		repo       *careerrepo.MemoryRepository
		svc        *careerservice.Service
		cliSvc     *service.CLIEventService
		ctx        context.Context
	)

	BeforeEach(func() {
		repo = careerrepo.NewMemoryRepository()
		svc = careerservice.NewService(repo)
		cliSvc = service.NewCLIEventService(svc)
		ctx = context.Background()
	})

	Context("BulkUpdateMetadata", func() {
		It("should return error when event IDs is nil", func() {
			var eventIDs []string
			eventIDs = nil

			summary, err := cliSvc.BulkUpdateMetadata(ctx, eventIDs, &service.BulkMetadataUpdate{})

			Expect(err).To(HaveOccurred())
			Expect(summary).To(BeNil())
		})

		It("should return error when event IDs is empty", func() {
			eventIDs := []string{}

			summary, err := cliSvc.BulkUpdateMetadata(ctx, eventIDs, &service.BulkMetadataUpdate{})

			Expect(err).To(HaveOccurred())
			Expect(summary).To(BeNil())
		})

		It("should return error in summary when event ID does not exist", func() {
			// Create one valid event
			validEvent := &careerdom.CareerEvent{
				Text: "Valid event",
			}
			err := svc.CaptureEvent(ctx, validEvent, careerservice.ManualEntry)
			Expect(err).NotTo(HaveOccurred())

			eventIDs := []string{validEvent.ID, "nonexistent-id"}
			update := &service.BulkMetadataUpdate{
				Company: "NewCorp",
			}

			summary, err := cliSvc.BulkUpdateMetadata(ctx, eventIDs, update)

			Expect(err).NotTo(HaveOccurred())
			Expect(summary).NotTo(BeNil())
			Expect(summary.Errors).To(HaveLen(1))
		})

		It("should successfully update company on all valid events", func() {
			// Create two events
			event1 := &careerdom.CareerEvent{
				Text: "Event 1",
			}
			event2 := &careerdom.CareerEvent{
				Text: "Event 2",
			}
			err := svc.CaptureEvent(ctx, event1, careerservice.ManualEntry)
			Expect(err).NotTo(HaveOccurred())
			err = svc.CaptureEvent(ctx, event2, careerservice.ManualEntry)
			Expect(err).NotTo(HaveOccurred())

			eventIDs := []string{event1.ID, event2.ID}
			update := &service.BulkMetadataUpdate{
				Company: "TechCorp",
			}

			summary, err := cliSvc.BulkUpdateMetadata(ctx, eventIDs, update)

			Expect(err).NotTo(HaveOccurred())
			Expect(summary).NotTo(BeNil())
			Expect(summary.AppliedCount).To(Equal(2))
			Expect(summary.FieldsUpdated).To(ContainElement("company"))

			// Verify updates persisted
			updated1, _ := svc.GetEventByID(ctx, event1.ID)
			Expect(updated1.Company).To(Equal("TechCorp"))

			updated2, _ := svc.GetEventByID(ctx, event2.ID)
			Expect(updated2.Company).To(Equal("TechCorp"))
		})

		It("should skip non-empty fields when applyIfEmpty is true", func() {
			// Create event with existing company
			event := &careerdom.CareerEvent{
				Text:    "Event with company",
				Company: "ExistingCorp",
			}
			err := svc.CaptureEvent(ctx, event, careerservice.ManualEntry)
			Expect(err).NotTo(HaveOccurred())

			eventIDs := []string{event.ID}
			update := &service.BulkMetadataUpdate{
				Company:              "NewCorp",
				ApplyIfEmptyCompany:  true,
				Project:              "NewProject",
				ApplyIfEmptyProject:  true,
			}

			summary, err := cliSvc.BulkUpdateMetadata(ctx, eventIDs, update)

			Expect(err).NotTo(HaveOccurred())
			Expect(summary).NotTo(BeNil())

			// Company should not be updated (already has value)
			updated, _ := svc.GetEventByID(ctx, event.ID)
			Expect(updated.Company).To(Equal("ExistingCorp"))

			// Project should be updated (was empty)
			Expect(updated.Project).To(Equal("NewProject"))
		})

		It("should validate all events before applying any changes", func() {
			// Create two valid events
			event1 := &careerdom.CareerEvent{
				Text: "Event 1",
			}
			event2 := &careerdom.CareerEvent{
				Text: "Event 2",
			}
			err := svc.CaptureEvent(ctx, event1, careerservice.ManualEntry)
			Expect(err).NotTo(HaveOccurred())
			err = svc.CaptureEvent(ctx, event2, careerservice.ManualEntry)
			Expect(err).NotTo(HaveOccurred())

			// Mix valid and invalid event IDs
			eventIDs := []string{event1.ID, "nonexistent-id", event2.ID}
			update := &service.BulkMetadataUpdate{
				Company: "TechCorp",
			}

			summary, err := cliSvc.BulkUpdateMetadata(ctx, eventIDs, update)

			// Operation should not fail completely, but report errors
			Expect(err).NotTo(HaveOccurred())
			Expect(summary).NotTo(BeNil())
			Expect(summary.Errors).To(HaveLen(1))
		})

		It("should return summary with count, fields updated, and errors", func() {
			event := &careerdom.CareerEvent{
				Text: "Test event",
			}
			err := svc.CaptureEvent(ctx, event, careerservice.ManualEntry)
			Expect(err).NotTo(HaveOccurred())

			eventIDs := []string{event.ID}
			update := &service.BulkMetadataUpdate{
				Company: "TechCorp",
				Project: "ProjectX",
				Tags:    []string{"technical", "leadership"},
			}

			summary, err := cliSvc.BulkUpdateMetadata(ctx, eventIDs, update)

			Expect(err).NotTo(HaveOccurred())
			Expect(summary).NotTo(BeNil())
			Expect(summary.EventsAffected).To(Equal(1))
			Expect(summary.AppliedCount).To(Equal(1))
			Expect(summary.SkippedCount).To(Equal(0))
			Expect(summary.FieldsUpdated).To(ContainElement("company"))
			Expect(summary.FieldsUpdated).To(ContainElement("project"))
			Expect(summary.FieldsUpdated).To(ContainElement("tags"))
			Expect(summary.Errors).To(BeEmpty())
		})

		It("should return error for invalid tags in bulk update", func() {
			event := &careerdom.CareerEvent{
				Text: "Test event",
			}
			err := svc.CaptureEvent(ctx, event, careerservice.ManualEntry)
			Expect(err).NotTo(HaveOccurred())

			eventIDs := []string{event.ID}
			update := &service.BulkMetadataUpdate{
				Tags: []string{"invalid-tag-not-in-allowed"},
			}

			summary, err := cliSvc.BulkUpdateMetadata(ctx, eventIDs, update)

			// Should have error in summary
			Expect(summary).NotTo(BeNil())
			Expect(summary.Errors).NotTo(BeEmpty())
		})

		It("should handle multiple events with mixed success and failure", func() {
			// Create three events
			event1 := &careerdom.CareerEvent{Text: "Event 1"}
			event2 := &careerdom.CareerEvent{Text: "Event 2"}
			event3 := &careerdom.CareerEvent{Text: "Event 3"}

			err := svc.CaptureEvent(ctx, event1, careerservice.ManualEntry)
			Expect(err).NotTo(HaveOccurred())
			err = svc.CaptureEvent(ctx, event2, careerservice.ManualEntry)
			Expect(err).NotTo(HaveOccurred())
			err = svc.CaptureEvent(ctx, event3, careerservice.ManualEntry)
			Expect(err).NotTo(HaveOccurred())

			// Mix valid and invalid IDs
			eventIDs := []string{event1.ID, "invalid-id-1", event2.ID, "invalid-id-2", event3.ID}
			update := &service.BulkMetadataUpdate{
				Company: "BulkCorp",
			}

			summary, err := cliSvc.BulkUpdateMetadata(ctx, eventIDs, update)

			Expect(err).NotTo(HaveOccurred())
			Expect(summary).NotTo(BeNil())
			Expect(summary.AppliedCount).To(Equal(3))
			Expect(summary.Errors).To(HaveLen(2))
		})
	})
})

