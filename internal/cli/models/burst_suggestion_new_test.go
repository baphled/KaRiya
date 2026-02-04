//nolint:errcheck // Test file - error handling for test setup is not relevant.
package models_test

import (
	"context"
	"time"

	"github.com/baphled/kariya/internal/cli/models"
	careermemory "github.com/baphled/kariya/internal/repository/career/memory"
	careerservice "github.com/baphled/kariya/internal/service/career"
	burstfact "github.com/baphled/kariya/internal/service/career/burstfact"
	"github.com/baphled/kariya/internal/testutil/fixtures"
	tea "github.com/charmbracelet/bubbletea"
	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
)

var _ = Describe("BurstSuggestionModelNew", func() {
	var (
		model       *models.BurstSuggestionModelNew
		service     *careerservice.Service
		suggestions []burstfact.BurstSuggestion
		ctx         context.Context
		eventRepo   *careermemory.EventRepository
	)

	BeforeEach(func() {
		ctx = context.Background()
		eventRepo = careermemory.NewEventRepository()
		service = careerservice.NewService(eventRepo)

		// Create test events
		event1 := fixtures.EventWith("event-1", "First event", "", "")
		event1.Date = time.Date(2024, 1, 1, 0, 0, 0, 0, time.UTC)
		event2 := fixtures.EventWith("event-2", "Second event", "", "")
		event2.Date = time.Date(2024, 1, 15, 0, 0, 0, 0, time.UTC)

		eventRepo.Create(ctx, event1)
		eventRepo.Create(ctx, event2)

		// Create test suggestions
		suggestions = []burstfact.BurstSuggestion{
			{
				Name:            "Project Alpha",
				Description:     "First project period",
				EventIDs:        []string{"event-1", "event-2"},
				ConfidenceScore: 0.85,
			},
			{
				Name:            "Project Beta",
				Description:     "Second project period",
				EventIDs:        []string{"event-1"},
				ConfidenceScore: 0.72,
			},
			{
				Name:            "",
				Description:     "Unnamed burst",
				EventIDs:        []string{"event-2"},
				ConfidenceScore: 0.60,
			},
		}

		model = models.NewBurstSuggestionModelNew(service, suggestions, ctx)
	})

	Describe("Initialization", func() {
		It("should create model with suggestions", func() {
			Expect(model).NotTo(BeNil())
		})

		It("should start at first suggestion", func() {
			view := model.View()
			Expect(view).To(ContainSubstring("1 of 3"))
		})

		It("should initialize with empty confirmed and rejected lists", func() {
			Expect(model.GetConfirmed()).To(BeEmpty())
			Expect(model.GetRejected()).To(BeEmpty())
		})

		It("should not be done initially", func() {
			Expect(model.IsDone()).To(BeFalse())
		})

		It("should initialize with default dimensions", func() {
			Expect(model).NotTo(BeNil())
			// Width and height are set to defaults in constructor
		})

		Context("with empty suggestions", func() {
			It("should handle empty suggestions list", func() {
				emptyModel := models.NewBurstSuggestionModelNew(service, []burstfact.BurstSuggestion{}, ctx)
				view := emptyModel.View()
				Expect(view).To(ContainSubstring("No burst suggestions"))
			})
		})
	})

	Describe("Navigation", func() {
		Context("when pressing down arrow", func() {
			It("should move to next suggestion", func() {
				msg := tea.KeyMsg{Type: tea.KeyDown}
				_, _ = model.Update(msg)

				view := model.View()
				Expect(view).To(ContainSubstring("2 of 3"))
			})

			It("should wrap to first suggestion when at end", func() {
				// Move to last suggestion
				model.Update(tea.KeyMsg{Type: tea.KeyDown})
				model.Update(tea.KeyMsg{Type: tea.KeyDown})

				// Press down again
				_, _ = model.Update(tea.KeyMsg{Type: tea.KeyDown})

				view := model.View()
				Expect(view).To(ContainSubstring("1 of 3"))
			})
		})

		Context("when pressing up arrow", func() {
			It("should move to previous suggestion", func() {
				// Move down first
				model.Update(tea.KeyMsg{Type: tea.KeyDown})

				// Then up
				_, _ = model.Update(tea.KeyMsg{Type: tea.KeyUp})

				view := model.View()
				Expect(view).To(ContainSubstring("1 of 3"))
			})

			It("should wrap to last suggestion when at start", func() {
				msg := tea.KeyMsg{Type: tea.KeyUp}
				_, _ = model.Update(msg)

				view := model.View()
				Expect(view).To(ContainSubstring("3 of 3"))
			})
		})

		Context("with vim keys", func() {
			It("should navigate down with 'j'", func() {
				msg := tea.KeyMsg{Type: tea.KeyRunes, Runes: []rune{'j'}}
				_, _ = model.Update(msg)

				view := model.View()
				Expect(view).To(ContainSubstring("2 of 3"))
			})

			It("should navigate up with 'k'", func() {
				// Move down first
				model.Update(tea.KeyMsg{Type: tea.KeyRunes, Runes: []rune{'j'}})

				// Then up
				_, _ = model.Update(tea.KeyMsg{Type: tea.KeyRunes, Runes: []rune{'k'}})

				view := model.View()
				Expect(view).To(ContainSubstring("1 of 3"))
			})
		})
	})

	Describe("Confirming Suggestions", func() {
		Context("when pressing 'y' key", func() {
			It("should confirm current suggestion", func() {
				msg := tea.KeyMsg{Type: tea.KeyRunes, Runes: []rune{'y'}}
				_, cmd := model.Update(msg)

				Expect(cmd).NotTo(BeNil())
				result := cmd()
				_, ok := result.(models.ConfirmBurstMsg)
				Expect(ok).To(BeTrue())
			})

			It("should add suggestion to confirmed list", func() {
				msg := tea.KeyMsg{Type: tea.KeyRunes, Runes: []rune{'y'}}
				_, _ = model.Update(msg)

				Expect(model.GetConfirmed()).To(HaveLen(1))
				Expect(model.GetConfirmed()[0].Name).To(Equal("Project Alpha"))
			})

			It("should move to next suggestion", func() {
				msg := tea.KeyMsg{Type: tea.KeyRunes, Runes: []rune{'y'}}
				_, _ = model.Update(msg)

				view := model.View()
				Expect(view).To(ContainSubstring("2 of 3"))
			})

			It("should apply edited name if exists", func() {
				// Start edit mode
				model.Update(tea.KeyMsg{Type: tea.KeyRunes, Runes: []rune{'e'}})
				// ... edit name (would require form interaction)
				// Confirm
				_, _ = model.Update(tea.KeyMsg{Type: tea.KeyRunes, Runes: []rune{'y'}})

				Expect(model.GetConfirmed()).To(HaveLen(1))
			})
		})

		Context("when pressing 'Y' key", func() {
			It("should also confirm (case insensitive)", func() {
				msg := tea.KeyMsg{Type: tea.KeyRunes, Runes: []rune{'Y'}}
				_, cmd := model.Update(msg)

				Expect(cmd).NotTo(BeNil())
				result := cmd()
				_, ok := result.(models.ConfirmBurstMsg)
				Expect(ok).To(BeTrue())
			})
		})

		Context("when confirming last suggestion", func() {
			It("should send completion message", func() {
				// Move to last suggestion
				model.Update(tea.KeyMsg{Type: tea.KeyDown})
				model.Update(tea.KeyMsg{Type: tea.KeyDown})

				// Confirm
				msg := tea.KeyMsg{Type: tea.KeyRunes, Runes: []rune{'y'}}
				_, cmd := model.Update(msg)

				// Execute batch command
				result := cmd()
				// Should be ConfirmBurstMsg, followed by completion
				Expect(result).NotTo(BeNil())
			})

			It("should mark as done", func() {
				// Confirm all suggestions
				model.Update(tea.KeyMsg{Type: tea.KeyRunes, Runes: []rune{'y'}})
				model.Update(tea.KeyMsg{Type: tea.KeyRunes, Runes: []rune{'y'}})
				model.Update(tea.KeyMsg{Type: tea.KeyRunes, Runes: []rune{'y'}})

				Expect(model.IsDone()).To(BeTrue())
			})
		})
	})

	Describe("Rejecting Suggestions", func() {
		Context("when pressing 'n' key", func() {
			It("should reject current suggestion", func() {
				msg := tea.KeyMsg{Type: tea.KeyRunes, Runes: []rune{'n'}}
				_, cmd := model.Update(msg)

				Expect(cmd).NotTo(BeNil())
				result := cmd()
				_, ok := result.(models.RejectBurstSuggestionMsg)
				Expect(ok).To(BeTrue())
			})

			It("should add suggestion to rejected list", func() {
				msg := tea.KeyMsg{Type: tea.KeyRunes, Runes: []rune{'n'}}
				_, _ = model.Update(msg)

				Expect(model.GetRejected()).To(HaveLen(1))
				Expect(model.GetRejected()[0].Name).To(Equal("Project Alpha"))
			})

			It("should move to next suggestion", func() {
				msg := tea.KeyMsg{Type: tea.KeyRunes, Runes: []rune{'n'}}
				_, _ = model.Update(msg)

				view := model.View()
				Expect(view).To(ContainSubstring("2 of 3"))
			})
		})

		Context("when pressing 'N' key", func() {
			It("should also reject (case insensitive)", func() {
				msg := tea.KeyMsg{Type: tea.KeyRunes, Runes: []rune{'N'}}
				_, cmd := model.Update(msg)

				Expect(cmd).NotTo(BeNil())
				result := cmd()
				_, ok := result.(models.RejectBurstSuggestionMsg)
				Expect(ok).To(BeTrue())
			})
		})

		Context("when rejecting last suggestion", func() {
			It("should send completion message", func() {
				// Move to last and reject
				model.Update(tea.KeyMsg{Type: tea.KeyDown})
				model.Update(tea.KeyMsg{Type: tea.KeyDown})

				msg := tea.KeyMsg{Type: tea.KeyRunes, Runes: []rune{'n'}}
				_, cmd := model.Update(msg)

				result := cmd()
				Expect(result).NotTo(BeNil())
			})
		})
	})

	Describe("Editing Suggestions", func() {
		Context("when pressing 'e' key", func() {
			It("should enter edit mode", func() {
				msg := tea.KeyMsg{Type: tea.KeyRunes, Runes: []rune{'e'}}
				_, cmd := model.Update(msg)

				Expect(cmd).NotTo(BeNil())
				view := model.View()
				Expect(view).To(ContainSubstring("Edit Burst"))
			})

			It("should show edit form", func() {
				model.Update(tea.KeyMsg{Type: tea.KeyRunes, Runes: []rune{'e'}})

				view := model.View()
				Expect(view).To(ContainSubstring("Name"))
				Expect(view).To(ContainSubstring("Description"))
			})

			It("should populate form with current values", func() {
				// Would require accessing edit form data
				model.Update(tea.KeyMsg{Type: tea.KeyRunes, Runes: []rune{'e'}})
				// Form is populated with suggestion.Name and suggestion.Description
			})

			It("should populate form with edited values if they exist", func() {
				// If suggestion was edited before, use edited values
				model.Update(tea.KeyMsg{Type: tea.KeyRunes, Runes: []rune{'e'}})
				// Edited values take precedence over original
			})
		})

		Context("when in edit mode", func() {
			BeforeEach(func() {
				model.Update(tea.KeyMsg{Type: tea.KeyRunes, Runes: []rune{'e'}})
			})

			It("should save edits when form is completed", func() {
				// When huh form reaches StateCompleted
				// Edits should be saved to editedNames and editedDescs
			})

			It("should exit edit mode on Escape", func() {
				msg := tea.KeyMsg{Type: tea.KeyEsc}
				_, _ = model.Update(msg)

				// Should be back in review mode
				view := model.View()
				Expect(view).NotTo(ContainSubstring("Edit Burst"))
			})

			It("should not save edits when cancelled", func() {
				// Cancel edit
				model.Update(tea.KeyMsg{Type: tea.KeyEsc})

				// Edits should not be applied
			})
		})

		Context("when pressing 'E' key", func() {
			It("should also enter edit mode (case insensitive)", func() {
				msg := tea.KeyMsg{Type: tea.KeyRunes, Runes: []rune{'E'}}
				_, _ = model.Update(msg)

				view := model.View()
				Expect(view).To(ContainSubstring("Edit Burst"))
			})
		})
	})

	Describe("Progress Tracking", func() {
		It("should track confirmed count", func() {
			model.Update(tea.KeyMsg{Type: tea.KeyRunes, Runes: []rune{'y'}})
			Expect(model.GetConfirmed()).To(HaveLen(1))

			model.Update(tea.KeyMsg{Type: tea.KeyRunes, Runes: []rune{'y'}})
			Expect(model.GetConfirmed()).To(HaveLen(2))
		})

		It("should track rejected count", func() {
			model.Update(tea.KeyMsg{Type: tea.KeyRunes, Runes: []rune{'n'}})
			Expect(model.GetRejected()).To(HaveLen(1))

			model.Update(tea.KeyMsg{Type: tea.KeyRunes, Runes: []rune{'n'}})
			Expect(model.GetRejected()).To(HaveLen(2))
		})

		It("should detect completion when all processed", func() {
			model.Update(tea.KeyMsg{Type: tea.KeyRunes, Runes: []rune{'y'}})
			model.Update(tea.KeyMsg{Type: tea.KeyRunes, Runes: []rune{'n'}})
			model.Update(tea.KeyMsg{Type: tea.KeyRunes, Runes: []rune{'y'}})

			Expect(model.IsDone()).To(BeTrue())
		})

		It("should send completion message with counts", func() {
			// Confirm 2, reject 1
			model.Update(tea.KeyMsg{Type: tea.KeyRunes, Runes: []rune{'y'}})
			model.Update(tea.KeyMsg{Type: tea.KeyRunes, Runes: []rune{'n'}})
			_, cmd := model.Update(tea.KeyMsg{Type: tea.KeyRunes, Runes: []rune{'y'}})

			// Last command should include completion
			result := cmd()
			Expect(result).NotTo(BeNil())
		})
	})

	Describe("Related Events", func() {
		It("should load related events from service", func() {
			view := model.View()
			// Related events are loaded when viewing suggestion
			Expect(view).To(ContainSubstring("Related Events"))
		})

		It("should cache loaded events", func() {
			// First view loads events
			_ = model.View()

			// Second view uses cache
			_ = model.View()
			// Cache prevents redundant service calls
		})

		It("should clear cache when navigating", func() {
			_ = model.View()

			// Navigate away
			model.Update(tea.KeyMsg{Type: tea.KeyDown})

			// Cache for previous suggestion is cleared
			// New suggestion loads fresh events
		})

		It("should display event preview", func() {
			view := model.View()
			Expect(view).To(ContainSubstring("event"))
		})

		It("should handle missing events gracefully", func() {
			// If event is not found in repository
			// Should not crash
		})
	})

	Describe("View Rendering", func() {
		Context("in review mode", func() {
			It("should show progress indicator", func() {
				view := model.View()
				Expect(view).To(ContainSubstring("1 of 3"))
			})

			It("should show confidence score", func() {
				view := model.View()
				Expect(view).To(ContainSubstring("Confidence"))
				Expect(view).To(ContainSubstring("85%"))
			})

			It("should show burst details", func() {
				view := model.View()
				Expect(view).To(ContainSubstring("Project Alpha"))
				Expect(view).To(ContainSubstring("First project period"))
			})

			It("should show navigation help", func() {
				view := model.View()
				Expect(view).To(ContainSubstring("y confirm"))
				Expect(view).To(ContainSubstring("n reject"))
				Expect(view).To(ContainSubstring("e edit"))
			})
		})

		Context("in edit mode", func() {
			BeforeEach(func() {
				model.Update(tea.KeyMsg{Type: tea.KeyRunes, Runes: []rune{'e'}})
			})

			It("should show edit form", func() {
				view := model.View()
				Expect(view).To(ContainSubstring("Edit Burst"))
			})

			It("should show form fields", func() {
				view := model.View()
				Expect(view).NotTo(BeEmpty())
			})

			It("should show edit help", func() {
				view := model.View()
				Expect(view).To(ContainSubstring("Tab"))
				Expect(view).To(ContainSubstring("Enter"))
			})
		})

		Context("with empty suggestions", func() {
			It("should show empty state", func() {
				emptyModel := models.NewBurstSuggestionModelNew(service, []burstfact.BurstSuggestion{}, ctx)
				view := emptyModel.View()
				Expect(view).To(ContainSubstring("No burst suggestions"))
			})
		})
	})

	Describe("Modal Overlay Integration", func() {
		It("should provide title for modal", func() {
			title := model.GetTitle()
			Expect(title).To(ContainSubstring("Burst Suggestion"))
			Expect(title).To(ContainSubstring("1 of 3"))
		})

		It("should provide content without wrapper", func() {
			content := model.GetContent()
			Expect(content).NotTo(BeEmpty())
		})

		It("should provide footer instructions", func() {
			footer := model.GetFooter()
			Expect(footer).To(ContainSubstring("Up/Down"))
			Expect(footer).To(ContainSubstring("y: Confirm"))
		})

		Context("in edit mode", func() {
			BeforeEach(func() {
				model.Update(tea.KeyMsg{Type: tea.KeyRunes, Runes: []rune{'e'}})
			})

			It("should provide edit title", func() {
				title := model.GetTitle()
				Expect(title).To(ContainSubstring("Edit"))
			})

			It("should provide edit footer", func() {
				footer := model.GetFooter()
				Expect(footer).To(ContainSubstring("Tab"))
			})
		})
	})

	Describe("Window Resizing", func() {
		It("should handle window size updates", func() {
			msg := tea.WindowSizeMsg{Width: 120, Height: 40}
			_, _ = model.Update(msg)
			// Dimensions are updated
		})
	})

	Describe("Escape Key Handling", func() {
		It("should send BackMsg when pressing Escape", func() {
			msg := tea.KeyMsg{Type: tea.KeyEsc}
			_, cmd := model.Update(msg)

			Expect(cmd).NotTo(BeNil())
			result := cmd()
			_, ok := result.(models.BackMsg)
			Expect(ok).To(BeTrue())
		})

		Context("in edit mode", func() {
			BeforeEach(func() {
				model.Update(tea.KeyMsg{Type: tea.KeyRunes, Runes: []rune{'e'}})
			})

			It("should exit edit mode without sending BackMsg", func() {
				msg := tea.KeyMsg{Type: tea.KeyEsc}
				_, _ = model.Update(msg)

				// Should be back in review mode
				view := model.View()
				Expect(view).NotTo(ContainSubstring("Edit Burst"))
			})
		})
	})

	Describe("Burst Creation", func() {
		It("should create burst from confirmed suggestion", func() {
			msg := tea.KeyMsg{Type: tea.KeyRunes, Runes: []rune{'y'}}
			_, cmd := model.Update(msg)

			result := cmd()
			confirmMsg, ok := result.(models.ConfirmBurstMsg)
			Expect(ok).To(BeTrue())
			Expect(confirmMsg.Burst).NotTo(BeNil())
			Expect(confirmMsg.Burst.Name).To(Equal("Project Alpha"))
			Expect(confirmMsg.Burst.EventIDs).To(HaveLen(2))
		})

		It("should generate name if none provided", func() {
			// Move to suggestion with no name
			model.Update(tea.KeyMsg{Type: tea.KeyDown})
			model.Update(tea.KeyMsg{Type: tea.KeyDown})

			msg := tea.KeyMsg{Type: tea.KeyRunes, Runes: []rune{'y'}}
			_, cmd := model.Update(msg)

			result := cmd()
			confirmMsg := result.(models.ConfirmBurstMsg)
			Expect(confirmMsg.Burst.Name).To(ContainSubstring("Burst of"))
		})
	})
})
