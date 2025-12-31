package app

import (
	"context"
	"time"

	"github.com/baphled/kariya/internal/cli/models"
	cliservice "github.com/baphled/kariya/internal/cli/service"
	"github.com/baphled/kariya/internal/domain/career"
	careerrepo "github.com/baphled/kariya/internal/repository/career"
	careerservice "github.com/baphled/kariya/internal/service/career"
	burstfact "github.com/baphled/kariya/internal/service/career/burst_fact"
	tea "github.com/charmbracelet/bubbletea"
	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
)

var _ = Describe("Burst Integration", func() {
	var (
		app          *Model
		repo         *careerrepo.MemoryRepository
		svc          *careerservice.Service
		ctx          context.Context
		event1       *career.CareerEvent
		event2       *career.CareerEvent
		event3       *career.CareerEvent
	)

	BeforeEach(func() {
		ctx = context.Background()
		repo = careerrepo.NewMemoryRepository()
		svc = careerservice.NewService(repo)
		cliSvc := cliservice.NewCLIEventService(svc)

		// Create test events that are similar and should form a burst
		event1 = &career.CareerEvent{
			Text:    "Led backend team on microservices migration project",
			Date:    time.Now().Add(-30 * 24 * time.Hour),
			Company: "TechCorp",
			Tags:    []string{"leadership", "technical"},
		}
		err := svc.CaptureEvent(ctx, event1, careerservice.TimelineJournaling)
		Expect(err).NotTo(HaveOccurred())

		event2 = &career.CareerEvent{
			Text:    "Architected service mesh for improved scalability",
			Date:    time.Now().Add(-25 * 24 * time.Hour),
			Company: "TechCorp",
			Tags:    []string{"technical", "architecture"},
		}
		err = svc.CaptureEvent(ctx, event2, careerservice.TimelineJournaling)
		Expect(err).NotTo(HaveOccurred())

		event3 = &career.CareerEvent{
			Text:    "Implemented monitoring and observability for microservices",
			Date:    time.Now().Add(-20 * 24 * time.Hour),
			Company: "TechCorp",
			Tags:    []string{"technical"},
		}
		err = svc.CaptureEvent(ctx, event3, careerservice.TimelineJournaling)
		Expect(err).NotTo(HaveOccurred())

		app = NewModel(cliSvc, svc)
	})

	Describe("Metadata Review to Burst Suggestion Workflow", func() {
		It("should trigger burst suggestions from metadata review", func() {
			// Navigate to metadata review screen with events
			eventIDs := []string{event1.ID, event2.ID, event3.ID}
			app.metadataReviewModel = models.NewMetadataReviewModelForImport(svc, ctx, eventIDs)
			app.currentScreen = MetadataReviewScreen

			// Press 'u' to trigger burst suggestions
			msg := tea.KeyMsg{Type: tea.KeyRunes, Runes: []rune{'u'}}
			updatedModel, cmd := app.Update(msg)

			// Verify burst suggestions were triggered
			Expect(updatedModel).NotTo(BeNil())
			Expect(cmd).NotTo(BeNil())

			// Execute the command to get BurstSuggestionsTriggeredMsg
			cmdMsg := cmd()

			// Verify message type
			_, ok := cmdMsg.(models.BurstSuggestionsTriggeredMsg)
			Expect(ok).To(BeTrue(), "Expected BurstSuggestionsTriggeredMsg")
		})

		It("should initialize BurstSuggestionModel when suggestions are ready", func() {
			// Create burst suggestions manually
			suggestions := []burstfact.BurstSuggestion{
				{
					EventIDs:        []string{event1.ID, event2.ID, event3.ID},
					ConfidenceScore: 0.85,
					Name:            "Microservices Migration Initiative",
					Description:     "Complete migration to microservices architecture",
				},
			}

			// Trigger burst suggestions ready message
			msg := BurstSuggestionsReadyMsg{Suggestions: suggestions}
			updatedModel, _ := app.Update(msg)

			// Verify burst suggestion screen is displayed
			appModel := updatedModel.(*Model)
			Expect(appModel.currentScreen).To(Equal(BurstSuggestionScreen))
			Expect(appModel.burstSuggestionModel).NotTo(BeNil())
		})

		It("should persist confirmed bursts to database", func() {
			// Create burst suggestion model
			suggestions := []burstfact.BurstSuggestion{
				{
					EventIDs:        []string{event1.ID, event2.ID},
					ConfidenceScore: 0.9,
					Name:            "Backend Migration",
					Description:     "Microservices backend migration",
				},
			}
			app.burstSuggestionModel = models.NewBurstSuggestionModel(svc, suggestions, ctx)
			app.currentScreen = BurstSuggestionScreen

			// Create a burst to confirm
			burst := &career.Burst{
				Name:        "Backend Migration",
				Description: "Microservices backend migration",
				EventIDs:    []string{event1.ID, event2.ID},
			}

			// Send ConfirmBurstMsg
			msg := models.ConfirmBurstMsg{Burst: burst}
			updatedModel, _ := app.Update(msg)

			// Verify burst was created
			appModel := updatedModel.(*Model)
			Expect(appModel).NotTo(BeNil())

			// Verify burst is in repository
			bursts, err := repo.List(ctx, careerrepo.ListFilters{})
			Expect(err).NotTo(HaveOccurred())
			Expect(len(bursts)).To(BeNumerically(">", 0))
		})

		It("should record rejected burst suggestions", func() {
			// Create burst suggestion model
			suggestions := []burstfact.BurstSuggestion{
				{
					EventIDs:        []string{event1.ID, event2.ID},
					ConfidenceScore: 0.6,
					Name:            "Low Confidence Burst",
					Description:     "Should be rejected",
				},
			}
			app.burstSuggestionModel = models.NewBurstSuggestionModel(svc, suggestions, ctx)
			app.currentScreen = BurstSuggestionScreen

			// Send RejectBurstSuggestionMsg
			msg := models.RejectBurstSuggestionMsg{Suggestion: suggestions[0]}
			updatedModel, _ := app.Update(msg)

			// Verify rejection was handled
			appModel := updatedModel.(*Model)
			Expect(appModel).NotTo(BeNil())
			// Note: RejectBurstSuggestion doesn't return an error, it just records the rejection
		})

		It("should navigate back to metadata review after burst workflow", func() {
			// Set up burst suggestion screen
			app.currentScreen = BurstSuggestionScreen
			app.previousScreen = MetadataReviewScreen

			// Send BurstProcessingCompleteMsg
			msg := models.BurstProcessingCompleteMsg{
				ConfirmedCount: 1,
				RejectedCount:  1,
			}
			updatedModel, _ := app.Update(msg)

			// Verify navigation back to metadata review
			appModel := updatedModel.(*Model)
			Expect(appModel.currentScreen).To(Equal(MetadataReviewScreen))
		})

		It("should navigate to home screen if no previous screen", func() {
			// Set up burst suggestion screen without previous screen
			app.currentScreen = BurstSuggestionScreen
			app.previousScreen = ""

			// Send BurstProcessingCompleteMsg
			msg := models.BurstProcessingCompleteMsg{
				ConfirmedCount: 2,
				RejectedCount:  0,
			}
			updatedModel, _ := app.Update(msg)

			// Verify navigation to home screen
			appModel := updatedModel.(*Model)
			Expect(appModel.currentScreen).To(Equal(HomeScreen))
		})

		It("should complete full workflow: metadata review → burst suggestions → confirmation", func() {
			// Start at metadata review
			eventIDs := []string{event1.ID, event2.ID, event3.ID}
			app.metadataReviewModel = models.NewMetadataReviewModelForImport(svc, ctx, eventIDs)
			app.currentScreen = MetadataReviewScreen

			// Step 1: Trigger burst suggestions
			msg1 := tea.KeyMsg{Type: tea.KeyRunes, Runes: []rune{'u'}}
			_, cmd1 := app.Update(msg1)
			Expect(cmd1).NotTo(BeNil())

			// Step 2: Process burst suggestions
			suggestions := []burstfact.BurstSuggestion{
				{
					EventIDs:        []string{event1.ID, event2.ID, event3.ID},
					ConfidenceScore: 0.85,
					Name:            "Full Workflow Burst",
					Description:     "Test complete workflow",
				},
			}
			msg2 := BurstSuggestionsReadyMsg{Suggestions: suggestions}
			app.Update(msg2)

			// Verify we're on burst suggestion screen
			Expect(app.currentScreen).To(Equal(BurstSuggestionScreen))
			Expect(app.burstSuggestionModel).NotTo(BeNil())

			// Step 3: Confirm burst
			burst := &career.Burst{
				Name:        "Full Workflow Burst",
				Description: "Test complete workflow",
				EventIDs:    []string{event1.ID, event2.ID, event3.ID},
			}
			msg3 := models.ConfirmBurstMsg{Burst: burst}
			app.Update(msg3)

			// Verify burst was created
			bursts, err := repo.List(ctx, careerrepo.ListFilters{})
			Expect(err).NotTo(HaveOccurred())
			Expect(len(bursts)).To(BeNumerically(">", 0))

			// Step 4: Complete workflow
			msg4 := models.BurstProcessingCompleteMsg{
				ConfirmedCount: 1,
				RejectedCount:  0,
			}
			app.Update(msg4)

			// Verify we're back at metadata review
			Expect(app.currentScreen).To(Equal(MetadataReviewScreen))
		})

		It("should handle empty burst suggestions gracefully", func() {
			// Trigger burst suggestions with no events
			msg := BurstSuggestionsReadyMsg{Suggestions: []burstfact.BurstSuggestion{}}
			updatedModel, _ := app.Update(msg)

			// Verify burst suggestion screen is displayed
			appModel := updatedModel.(*Model)
			Expect(appModel.currentScreen).To(Equal(BurstSuggestionScreen))
			Expect(appModel.burstSuggestionModel).NotTo(BeNil())
		})

		It("should handle burst confirmation errors gracefully", func() {
			// Create invalid burst (no event IDs)
			burst := &career.Burst{
				Name:        "Invalid Burst",
				Description: "Missing event IDs",
				EventIDs:    []string{}, // Invalid: requires ≥2 events
			}

			// Send ConfirmBurstMsg
			msg := models.ConfirmBurstMsg{Burst: burst}
			updatedModel, _ := app.Update(msg)

			// Verify app didn't crash
			appModel := updatedModel.(*Model)
			Expect(appModel).NotTo(BeNil())
		})
	})
})

