package app

import (
	"context"
	"time"

	"github.com/baphled/kariya/internal/cli/service"
	"github.com/baphled/kariya/internal/domain/career"
	careerrepo "github.com/baphled/kariya/internal/repository/career"
	careerservice "github.com/baphled/kariya/internal/service/career"
	tea "github.com/charmbracelet/bubbletea"
	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
)

var _ = Describe("End-to-End Integration Tests", func() {
	var (
		repo       careerrepo.Repository
		svc        *careerservice.Service
		cliService *service.CLIEventService
		model      *Model
		ctx        context.Context
	)

	BeforeEach(func() {
		ctx = context.Background()
		repo = careerrepo.NewMemoryRepository()
		svc = careerservice.NewService(repo)
		cliService = service.NewCLIEventService(svc)
		model = NewModel(cliService, svc)
	})

	Describe("Complete Event Capture Workflow", func() {
		Context("when user completes full event capture to persistence", func() {
			It("should capture event through form and persist to repository", func() {
				// Navigate to capture screen
				model.currentScreen = CaptureScreen

				// Simulate filling form fields
				eventText := "Led cross-functional team to deliver critical migration project"
				eventDate := time.Now().Add(-24 * time.Hour)
				
				// Submit event through CLI service
				err := cliService.CaptureEvent(
					ctx,
					eventText,
					eventDate,
					careerservice.ManualEntry,
					service.WithCompany("TechCorp Inc."),
					service.WithProject("Platform Migration"),
					service.WithTags([]string{"leadership", "technical"}),
				)
				Expect(err).ToNot(HaveOccurred())

				// Verify event was persisted to repository
				events, err := repo.List(ctx, careerrepo.ListFilters{
					Limit: 10,
				})
				Expect(err).ToNot(HaveOccurred())
				Expect(events).To(HaveLen(1))
				Expect(events[0].Text).To(Equal(eventText))
				Expect(events[0].Company).To(Equal("TechCorp Inc."))
				Expect(events[0].Project).To(Equal("Platform Migration"))
				Expect(events[0].Tags).To(ContainElements("leadership", "technical"))
				Expect(events[0].Date.Format("2006-01-02")).To(Equal(eventDate.Format("2006-01-02")))
			})

			It("should handle Timeline Journaling mode with date constraints", func() {
				eventText := "Fixed critical production bug affecting 10k users"
						// Use 1 hour ago to ensure we're well within the 30-day window
						eventDate := time.Now().Add(-1 * time.Hour)
			
						err := cliService.CaptureEvent(
							ctx,
							eventText,
							eventDate,
							careerservice.TimelineJournaling,
							service.WithTags([]string{"technical"}),
						)
						Expect(err).ToNot(HaveOccurred())

						// Verify event persisted - use service layer to query
						events, err := svc.ListEvents(ctx, careerrepo.ListFilters{})
				Expect(err).ToNot(HaveOccurred())
				Expect(events).To(HaveLen(1))
				Expect(events[0].Text).To(Equal(eventText))
			})

			It("should handle CV Backfill mode with historical dates", func() {
				eventText := "Architected and delivered microservices platform"
				eventDate := time.Now().Add(-365 * 24 * time.Hour) // 1 year ago

				err := cliService.CaptureEvent(
					ctx,
					eventText,
					eventDate,
					careerservice.CVBackfill,
					service.WithCompany("OldCorp"),
					service.WithTags([]string{"technical", "leadership"}),
				)
				Expect(err).ToNot(HaveOccurred())

				// Verify event persisted
				events, err := repo.List(ctx, careerrepo.ListFilters{})
				Expect(err).ToNot(HaveOccurred())
				Expect(events).To(HaveLen(1))
				Expect(events[0].Text).To(Equal(eventText))
				Expect(events[0].Company).To(Equal("OldCorp"))
			})
		})

		Context("when user navigates through complete capture workflow", func() {
			It("should navigate from home to capture to success screen", func() {
				// Start at home
				Expect(model.currentScreen).To(Equal(HomeScreen))

				// Navigate to capture
				msg := tea.KeyMsg{Type: tea.KeyRunes, Runes: []rune{'c'}}
				newModel, _ := model.Update(msg)
				model = newModel.(*Model)
				Expect(model.currentScreen).To(Equal(CaptureScreen))

				// Simulate successful form submission
				eventText := "Completed important milestone"
				eventDate := time.Now()
				err := cliService.CaptureEvent(
					ctx,
					eventText,
					eventDate,
					careerservice.ManualEntry,
				)
				Expect(err).ToNot(HaveOccurred())

				// Verify navigation to success screen happens through FormSubmittedMsg
				submitMsg := FormSubmittedMsg{
					Event: &career.CareerEvent{
						ID:      "test-id",
						Text:    eventText,
						Date:    eventDate,
						Tags:    []string{},
						Company: "",
						Project: "",
					},
				}
				newModel, _ = model.Update(submitMsg)
				model = newModel.(*Model)
				Expect(model.currentScreen).To(Equal(SuccessScreen))
			})
		})

		Context("when retrieving events after capture", func() {
			It("should list all captured events", func() {
				// Capture multiple events
				events := []struct {
					text    string
					company string
					tags    []string
				}{
					{"Event 1", "Company A", []string{"technical"}},
					{"Event 2", "Company B", []string{"leadership"}},
					{"Event 3", "Company A", []string{"technical", "product"}},
				}

				for _, e := range events {
					err := cliService.CaptureEvent(
						ctx,
						e.text,
						time.Now().Add(-24*time.Hour),
						careerservice.ManualEntry,
						service.WithCompany(e.company),
						service.WithTags(e.tags),
					)
					Expect(err).ToNot(HaveOccurred())
				}

				// List all events
				listedEvents, err := cliService.ListEvents(ctx, nil)
				Expect(err).ToNot(HaveOccurred())
				Expect(listedEvents).To(HaveLen(3))
			})

			It("should filter events by tags", func() {
				// Capture events with different tags
				err := cliService.CaptureEvent(
					ctx,
					"Technical event",
					time.Now(),
					careerservice.ManualEntry,
					service.WithTags([]string{"technical"}),
				)
				Expect(err).ToNot(HaveOccurred())

				err = cliService.CaptureEvent(
					ctx,
					"Leadership event",
					time.Now(),
					careerservice.ManualEntry,
					service.WithTags([]string{"leadership"}),
				)
				Expect(err).ToNot(HaveOccurred())

				// Filter by technical tag
				filters := &careerrepo.ListFilters{
					Tags: []string{"technical"},
				}
				filteredEvents, err := cliService.ListEvents(ctx, filters)
				Expect(err).ToNot(HaveOccurred())
				Expect(filteredEvents).To(HaveLen(1))
				Expect(filteredEvents[0].Tags).To(ContainElement("technical"))
			})
		})
	})
})

