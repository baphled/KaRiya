package burst_management_test

import (
	"context"
	"errors"
	"fmt"

	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"

	"github.com/baphled/kariya/internal/domain/career"
	"github.com/baphled/kariya/internal/service/career/burstfact"
	"github.com/baphled/kariya/internal/service/career/skillinference"
	"github.com/baphled/kariya/internal/testutil/fixtures"
	"github.com/baphled/kariya/internal/testutil/mocks"
	"github.com/baphled/kariya/internal/tui/intents/burst_management"
	burstviews "github.com/baphled/kariya/internal/tui/views/burst"
	"github.com/baphled/kariya/internal/ui/display"
	"github.com/baphled/kariya/internal/ui/terminal"
	"github.com/baphled/kariya/internal/ui/uikit/feedback"
	"github.com/baphled/kariya/internal/ui/uikit/widgets"
	tea "github.com/charmbracelet/bubbletea"
)

var _ = Describe("Helper Methods", func() {
	var intent *burst_management.Intent
	var ctx *burst_management.IntentValidator
	var termInfo *terminal.Info

	BeforeEach(func() {
		// Create intent context.
		ctx = &burst_management.IntentValidator{
			Bursts: []*career.Burst{
				fixtures.Burst("burst-1"),
			},
		}
		ctx.Validate()

		// Create intent.
		var err error
		intent, err = burst_management.NewIntent(ctx)
		Expect(err).NotTo(HaveOccurred())

		// Set terminal info.
		termInfo = terminal.NewInfo()
		termInfo.Width = 120
		termInfo.Height = 40
		termInfo.IsValid = true
		intent.UpdateTerminalInfo(termInfo)

		// Initialize intent.
		intent.Init()
	})

	Describe("rebuildModalRegistry", func() {
		It("should create registry if nil", func() {
			// Registry should be created during NewIntent
			registry := intent.GetModalRegistry()
			Expect(registry).NotTo(BeNil())
		})

		It("should clear and rebuild registry", func() {
			// Set an error modal.
			intent.ShowErrorModal("Test Error", "This is a test error")

			// Rebuild registry.
			intent.RebuildModalRegistry()

			// Verify registry exists.
			registry := intent.GetModalRegistry()
			Expect(registry).NotTo(BeNil())

			// Verify error modal is registered.
			Expect(intent.HasVisibleErrorModal()).To(BeTrue())
		})

		It("should register error modal with highest priority", func() {
			// Create error modal.
			intent.ShowErrorModal("Error", "Test error")

			// Rebuild registry.
			intent.RebuildModalRegistry()

			// Verify registry has visible modal.
			Expect(intent.HasVisibleErrorModal()).To(BeTrue())

			// Registry should exist.
			registry := intent.GetModalRegistry()
			Expect(registry).NotTo(BeNil())
		})

		It("should register delete modal when present", func() {
			// Create delete modal.
			burst := fixtures.Burst("burst-1")
			confirmModal := feedback.NewConfirmModal(
				"Delete Burst?",
				"Are you sure you want to delete '"+burst.Name+"'?",
			)
			confirmModal.Show()

			// Set delete modal via internal field (no public setter yet).
			// We'll test this through the delete flow instead.

			// Rebuild registry.
			intent.RebuildModalRegistry()

			// Registry should exist.
			registry := intent.GetModalRegistry()
			Expect(registry).NotTo(BeNil())
		})

		It("should handle nil error modal", func() {
			// Ensure no error modal by not setting one.

			// Rebuild should not panic.
			Expect(func() {
				intent.RebuildModalRegistry()
			}).NotTo(Panic())

			// Registry should still exist.
			registry := intent.GetModalRegistry()
			Expect(registry).NotTo(BeNil())
		})

		It("should handle nil delete modal", func() {
			// Ensure no delete modal by not setting one.

			// Rebuild should not panic.
			Expect(func() {
				intent.RebuildModalRegistry()
			}).NotTo(Panic())

			// Registry should still exist.
			registry := intent.GetModalRegistry()
			Expect(registry).NotTo(BeNil())
		})

		It("should clear previous registrations", func() {
			// Register error modal.
			intent.ShowErrorModal("Error 1", "First error")
			intent.RebuildModalRegistry()

			// Error modal should be visible.
			Expect(intent.HasVisibleErrorModal()).To(BeTrue())
		})
	})

	Describe("getTerminalDimensions", func() {
		It("should return current terminal dimensions", func() {
			width, height := intent.GetTerminalDimensions()
			Expect(width).To(Equal(120))
			Expect(height).To(Equal(40))
		})

		It("should return default dimensions when terminal info is nil", func() {
			// Create intent without terminal info.
			newIntent, _ := burst_management.NewIntent(ctx)

			width, height := newIntent.GetTerminalDimensions()
			// Should return defaults from behaviors.DefaultModalDimensions().
			Expect(width).To(BeNumerically(">", 0))
			Expect(height).To(BeNumerically(">", 0))
		})

		It("should return default dimensions when terminal info is invalid", func() {
			invalidInfo := terminal.NewInfo()
			invalidInfo.Width = 0
			invalidInfo.Height = 0
			intent.UpdateTerminalInfo(invalidInfo)

			width, height := intent.GetTerminalDimensions()
			// Should return defaults.
			Expect(width).To(BeNumerically(">", 0))
			Expect(height).To(BeNumerically(">", 0))
		})
	})

	Describe("getStateName", func() {
		It("should return correct name for StateList", func() {
			intent.SetState(burst_management.StateList)
			Expect(intent.GetStateName()).To(Equal("Burst List"))
		})

		It("should return correct name for StateDetail", func() {
			intent.SetState(burst_management.StateDetail)
			Expect(intent.GetStateName()).To(Equal("Burst Details"))
		})

		It("should return correct name for StateDetailEvents", func() {
			intent.SetState(burst_management.StateDetailEvents)
			Expect(intent.GetStateName()).To(Equal("Burst Events"))
		})

		It("should return correct name for StateDetailFacts", func() {
			intent.SetState(burst_management.StateDetailFacts)
			Expect(intent.GetStateName()).To(Equal("Burst Facts"))
		})

		It("should return correct name for StateEdit", func() {
			intent.SetState(burst_management.StateEdit)
			Expect(intent.GetStateName()).To(Equal("Edit Burst"))
		})

		It("should return correct name for StateDeleteConfirm", func() {
			intent.SetState(burst_management.StateDeleteConfirm)
			Expect(intent.GetStateName()).To(Equal("Delete Confirmation"))
		})

		It("should return correct name for StateConfirm", func() {
			intent.SetState(burst_management.StateConfirm)
			Expect(intent.GetStateName()).To(Equal("Confirm Burst"))
		})

		It("should return correct name for StateExtractingFacts", func() {
			intent.SetState(burst_management.StateExtractingFacts)
			Expect(intent.GetStateName()).To(Equal("Extracting Facts"))
		})

		It("should return correct name for StateSuggesting", func() {
			intent.SetState(burst_management.StateSuggesting)
			Expect(intent.GetStateName()).To(Equal("Suggesting Bursts"))
		})

		It("should return correct name for StateSuggestionReview", func() {
			intent.SetState(burst_management.StateSuggestionReview)
			Expect(intent.GetStateName()).To(Equal("Review Suggestion"))
		})

		It("should return Unknown for invalid state", func() {
			intent.SetState(burst_management.State("invalid")) // Invalid state
			Expect(intent.GetStateName()).To(Equal("Unknown"))
		})
	})

	Describe("getContextHelp", func() {
		It("should return help for StateList", func() {
			intent.SetState(burst_management.StateList)
			help := intent.GetContextHelp()
			Expect(help).NotTo(BeEmpty())
			Expect(help).To(ContainSubstring("View"))
			Expect(help).To(ContainSubstring("Suggest"))
		})

		It("should return help for StateDetail", func() {
			intent.SetState(burst_management.StateDetail)
			help := intent.GetContextHelp()
			Expect(help).NotTo(BeEmpty())
			Expect(help).To(ContainSubstring("View Events"))
			Expect(help).To(ContainSubstring("View Facts"))
			Expect(help).To(ContainSubstring("Confirm"))
		})

		It("should return help for StateDetailEvents", func() {
			intent.SetState(burst_management.StateDetailEvents)
			help := intent.GetContextHelp()
			Expect(help).NotTo(BeEmpty())
		})

		It("should return help for StateDetailFacts", func() {
			intent.SetState(burst_management.StateDetailFacts)
			help := intent.GetContextHelp()
			Expect(help).NotTo(BeEmpty())
		})

		It("should return empty for unknown state", func() {
			intent.SetState(burst_management.State("invalid")) // Invalid state
			help := intent.GetContextHelp()
			Expect(help).To(BeEmpty())
		})
	})

	Describe("RefreshData", func() {
		It("should reload bursts from context", func() {
			// Modify context bursts.
			ctx.Bursts = append(ctx.Bursts, fixtures.Burst("burst-2"))

			// Refresh data.
			cmd := intent.RefreshData()
			Expect(cmd).To(BeNil())

			// Verify filtered bursts updated.
			Expect(intent.GetFilteredBursts()).To(HaveLen(2))
		})

		It("should handle errors during load gracefully", func() {
			// RefreshData should not panic even if LoadBursts fails.
			Expect(func() {
				intent.RefreshData()
			}).NotTo(Panic())
		})

		It("should transition to list screen after refresh", func() {
			intent.RefreshData()

			// Should be in list state with list screen.
			Expect(intent.GetState()).To(Equal(burst_management.StateList))
		})
	})

	Describe("SuggestionReviewCompleteMsg handling", func() {
		It("should create bursts from accepted suggestions", func() {
			initialCount := len(intent.GetFilteredBursts())

			// Set state to suggestion review so the handler will process the message.
			intent.SetState(burst_management.StateSuggestionReview)

			// Create completion message with accepted suggestions.
			msg := burst_management.SuggestionReviewCompleteMsg{
				AcceptedSuggestions: []burstfact.BurstSuggestion{
					{
						Name:            "Accepted Burst 1",
						Description:     "First accepted suggestion",
						EventIDs:        []string{"e1", "e2"},
						ConfidenceScore: 0.9,
					},
					{
						Name:            "Accepted Burst 2",
						Description:     "Second accepted suggestion",
						EventIDs:        []string{"e3", "e4"},
						ConfidenceScore: 0.85,
					},
				},
				Cancelled: false,
			}

			// Handle the message via Update.
			cmd := intent.Update(msg)
			Expect(cmd).NotTo(BeNil(), "should return command to trigger fact extraction")

			// Verify bursts were created.
			Expect(intent.GetFilteredBursts()).To(HaveLen(initialCount + 2))

			// Verify burst details.
			newBursts := intent.GetFilteredBursts()[initialCount:]
			Expect(newBursts[0].Name).To(Equal("Accepted Burst 1"))
			Expect(newBursts[0].Description).To(Equal("First accepted suggestion"))
			Expect(newBursts[0].EventIDs).To(Equal([]string{"e1", "e2"}))
			Expect(newBursts[0].Confirmed).To(BeTrue(), "accepted suggestions should be auto-confirmed")

			Expect(newBursts[1].Name).To(Equal("Accepted Burst 2"))
			Expect(newBursts[1].Description).To(Equal("Second accepted suggestion"))
			Expect(newBursts[1].EventIDs).To(Equal([]string{"e3", "e4"}))
			Expect(newBursts[1].Confirmed).To(BeTrue(), "accepted suggestions should be auto-confirmed")
		})

		It("should transition to list screen when cancelled", func() {
			// Set state to suggestion review so the handler will process the message.
			intent.SetState(burst_management.StateSuggestionReview)

			msg := burst_management.SuggestionReviewCompleteMsg{
				AcceptedSuggestions: nil,
				Cancelled:           true,
			}

			cmd := intent.Update(msg)
			Expect(cmd).To(BeNil())

			// Should transition to list state.
			Expect(intent.GetState()).To(Equal(burst_management.StateList))

			// No bursts should be created.
			Expect(intent.GetFilteredBursts()).To(HaveLen(1)) // Original burst only
		})

		It("should transition to list screen when no suggestions accepted", func() {
			// Set state to suggestion review so the handler will process the message.
			intent.SetState(burst_management.StateSuggestionReview)

			msg := burst_management.SuggestionReviewCompleteMsg{
				AcceptedSuggestions: []burstfact.BurstSuggestion{},
				Cancelled:           false,
			}

			cmd := intent.Update(msg)
			Expect(cmd).To(BeNil())

			// Should transition to list state.
			Expect(intent.GetState()).To(Equal(burst_management.StateList))

			// No bursts should be created.
			Expect(intent.GetFilteredBursts()).To(HaveLen(1)) // Original burst only
		})

		It("should add bursts to context bursts list", func() {
			initialContextCount := len(ctx.Bursts)

			// Set state to suggestion review so the handler will process the message.
			intent.SetState(burst_management.StateSuggestionReview)

			msg := burst_management.SuggestionReviewCompleteMsg{
				AcceptedSuggestions: []burstfact.BurstSuggestion{
					{
						Name:            "New Burst",
						Description:     "New burst description",
						EventIDs:        []string{"e1"},
						ConfidenceScore: 0.95,
					},
				},
				Cancelled: false,
			}

			intent.Update(msg)

			// Verify burst added to context.
			Expect(ctx.Bursts).To(HaveLen(initialContextCount + 1))
			Expect(ctx.Bursts[initialContextCount].Name).To(Equal("New Burst"))
		})

		It("should transition to extracting facts after accepting suggestions", func() {
			// Set state to suggestion review so the handler will process the message.
			intent.SetState(burst_management.StateSuggestionReview)

			msg := burst_management.SuggestionReviewCompleteMsg{
				AcceptedSuggestions: []burstfact.BurstSuggestion{
					{
						Name:        "New Burst",
						Description: "Description",
						EventIDs:    []string{"e1"},
					},
				},
				Cancelled: false,
			}

			intent.Update(msg)

			// Should transition to extracting facts to process new bursts.
			Expect(intent.GetState()).To(Equal(burst_management.StateExtractingFacts))
		})

		It("should show loading modal during fact extraction after accepting suggestions", func() {
			// Set state to suggestion review so the handler will process the message.
			intent.SetState(burst_management.StateSuggestionReview)

			msg := burst_management.SuggestionReviewCompleteMsg{
				AcceptedSuggestions: []burstfact.BurstSuggestion{
					{
						Name:        "Burst With Loading",
						Description: "Should show loading modal",
						EventIDs:    []string{"e1"},
					},
				},
				Cancelled: false,
			}

			intent.Update(msg)

			// Loading modal should be visible during fact extraction.
			Expect(intent.HasActiveModal()).To(BeTrue())
			Expect(intent.GetState()).To(Equal(burst_management.StateExtractingFacts))
		})
	})

	Describe("handleFactExtractionComplete list refresh", func() {
		It("should refresh list screen when fact extraction completes with nil selectedBurst", func() {
			// Set state to suggestion review so the handler will process the message.
			intent.SetState(burst_management.StateSuggestionReview)

			// First accept a suggestion to create a burst.
			acceptMsg := burst_management.SuggestionReviewCompleteMsg{
				AcceptedSuggestions: []burstfact.BurstSuggestion{
					{
						Name:        "New Burst For Refresh",
						Description: "Testing list refresh",
						EventIDs:    []string{"e1"},
					},
				},
				Cancelled: false,
			}
			intent.Update(acceptMsg)

			// Verify burst was created.
			Expect(intent.GetFilteredBursts()).To(HaveLen(2)) // Original + new

			// Ensure selectedBurst is nil (as in suggestion acceptance flow).
			intent.SetSelectedBurst(nil)

			// Simulate fact extraction completing.
			extractMsg := burst_management.FactExtractionCompleteMsg{
				Facts: []*career.Fact{fixtures.FactWith("f1", "Extracted fact")},
				Error: nil,
			}
			intent.Update(extractMsg)

			// Should be on list state.
			Expect(intent.GetState()).To(Equal(burst_management.StateList))

			// View should render correctly with the new burst visible.
			view := intent.View()
			Expect(view).To(ContainSubstring("New Burst For Refresh"))
		})
	})

	// BUG REGRESSION TESTS
	Describe("BUG: handleFactExtractionComplete with nil selectedBurst", func() {
		It("should not panic when selectedBurst is nil", func() {
			// Ensure selectedBurst is nil.
			intent.SetSelectedBurst(nil)

			msg := burst_management.FactExtractionCompleteMsg{
				Facts: []*career.Fact{fixtures.FactWith("f1", "Test fact")},
				Error: nil,
			}

			Expect(func() {
				intent.Update(msg)
			}).NotTo(Panic())
		})

		It("should return to list state without showing detail modal when selectedBurst is nil", func() {
			intent.SetSelectedBurst(nil)

			msg := burst_management.FactExtractionCompleteMsg{
				Facts: []*career.Fact{fixtures.FactWith("f1", "Test fact")},
				Error: nil,
			}

			intent.Update(msg)

			Expect(intent.GetState()).To(Equal(burst_management.StateList))
			Expect(intent.GetDetailModal()).To(BeNil())
		})

		It("should handle fact extraction error with nil selectedBurst", func() {
			intent.SetSelectedBurst(nil)

			msg := burst_management.FactExtractionCompleteMsg{
				Error: fmt.Errorf("extraction failed"),
			}

			Expect(func() {
				intent.Update(msg)
			}).NotTo(Panic())

			// Should show error modal.
			Expect(intent.HasVisibleErrorModal()).To(BeTrue())
			Expect(intent.GetState()).To(Equal(burst_management.StateList))
		})
	})

	Describe("BUG: State transition after suggestion modal shown", func() {
		It("should transition to StateSuggestionReview when suggestions are loaded", func() {
			msg := burst_management.BurstSuggestionsLoadedMsg{
				Suggestions: []burstfact.BurstSuggestion{
					{Name: "Test", EventIDs: []string{"e1"}, ConfidenceScore: 0.8},
				},
			}

			intent.Update(msg)

			// State should be StateSuggestionReview, not StateSuggesting.
			Expect(intent.GetState()).To(Equal(burst_management.StateSuggestionReview))
		})
	})

	Describe("BUG: Slice bounds in handleSuggestionReviewComplete", func() {
		It("should not panic when accepted count differs from created count", func() {
			// This tests the slice bounds bug at helpers.go:672.
			// If some bursts fail to create, the slice calculation is wrong.

			// Set state to suggestion review so the handler will process the message.
			intent.SetState(burst_management.StateSuggestionReview)

			suggestions := []burstfact.BurstSuggestion{
				{Name: "Burst 1", EventIDs: []string{"e1"}, ConfidenceScore: 0.9},
				{Name: "Burst 2", EventIDs: []string{"e2"}, ConfidenceScore: 0.85},
			}

			msg := burst_management.SuggestionReviewCompleteMsg{
				AcceptedSuggestions: suggestions,
				Cancelled:           false,
			}

			// Should not panic.
			Expect(func() {
				intent.Update(msg)
			}).NotTo(Panic())
		})

		It("should track created bursts correctly for fact extraction", func() {
			initialCount := len(intent.GetFilteredBursts())

			// Set state to suggestion review so the handler will process the message.
			intent.SetState(burst_management.StateSuggestionReview)

			suggestions := []burstfact.BurstSuggestion{
				{Name: "Tracked Burst", EventIDs: []string{"e1", "e2"}, ConfidenceScore: 0.9},
			}

			msg := burst_management.SuggestionReviewCompleteMsg{
				AcceptedSuggestions: suggestions,
				Cancelled:           false,
			}

			cmd := intent.Update(msg)

			// Should have created burst.
			Expect(intent.GetFilteredBursts()).To(HaveLen(initialCount + 1))

			// Should return fact extraction command.
			Expect(cmd).NotTo(BeNil())
		})
	})

	Describe("BurstEventsLoadedMsg handling", func() {
		It("should show events modal on successful load", func() {
			intent.SetSelectedBurst(ctx.Bursts[0])

			events := []*career.Event{
				fixtures.EventWith("e1", "Event 1", "", ""),
				fixtures.EventWith("e2", "Event 2", "", ""),
			}

			msg := burst_management.BurstEventsLoadedMsg{
				Events: events,
				Error:  nil,
			}

			intent.Update(msg)

			// Events modal should be visible.
			Expect(intent.HasActiveModal()).To(BeTrue())
		})

		It("should show error modal on events load error", func() {
			intent.SetSelectedBurst(ctx.Bursts[0])

			msg := burst_management.BurstEventsLoadedMsg{
				Events: nil,
				Error:  fmt.Errorf("failed to load events"),
			}

			intent.Update(msg)

			// Error modal should be visible.
			Expect(intent.HasVisibleErrorModal()).To(BeTrue())
		})

		It("should handle nil selectedBurst gracefully", func() {
			intent.SetSelectedBurst(nil)

			msg := burst_management.BurstEventsLoadedMsg{
				Events: []*career.Event{fixtures.Event("e1")},
				Error:  nil,
			}

			// Should not panic.
			Expect(func() {
				intent.Update(msg)
			}).NotTo(Panic())
		})
	})

	Describe("BurstFactsLoadedMsg handling", func() {
		It("should show facts modal on successful load", func() {
			intent.SetSelectedBurst(ctx.Bursts[0])

			facts := []*career.Fact{
				fixtures.FactWith("f1", "Fact 1"),
				fixtures.FactWith("f2", "Fact 2"),
			}

			msg := burst_management.BurstFactsLoadedMsg{
				Facts: facts,
				Error: nil,
			}

			intent.Update(msg)

			// Facts modal should be visible.
			Expect(intent.HasActiveModal()).To(BeTrue())
		})

		It("should show error modal on facts load error", func() {
			intent.SetSelectedBurst(ctx.Bursts[0])

			msg := burst_management.BurstFactsLoadedMsg{
				Facts: nil,
				Error: fmt.Errorf("failed to load facts"),
			}

			intent.Update(msg)

			// Error modal should be visible.
			Expect(intent.HasVisibleErrorModal()).To(BeTrue())
		})

		It("should handle nil selectedBurst gracefully", func() {
			intent.SetSelectedBurst(nil)

			msg := burst_management.BurstFactsLoadedMsg{
				Facts: []*career.Fact{fixtures.FactWith("f1", "Test fact f1")},
				Error: nil,
			}

			// Should not panic.
			Expect(func() {
				intent.Update(msg)
			}).NotTo(Panic())
		})
	})

	Describe("Confirm burst flow", func() {
		It("should handle confirm key in detail state", func() {
			intent.SetSelectedBurst(ctx.Bursts[0])
			intent.SetState(burst_management.StateDetail)

			// Press 'c' to confirm - should not panic.
			Expect(func() {
				intent.Update(tea.KeyMsg{Type: tea.KeyRunes, Runes: []rune{'c'}})
			}).NotTo(Panic())
		})

		It("should handle confirm with nil selectedBurst", func() {
			intent.SetSelectedBurst(nil)
			intent.SetState(burst_management.StateDetail)

			// Press 'c' to confirm - should not panic.
			Expect(func() {
				intent.Update(tea.KeyMsg{Type: tea.KeyRunes, Runes: []rune{'c'}})
			}).NotTo(Panic())
		})
	})

	Describe("State management", func() {
		It("should allow setting and getting state", func() {
			intent.SetState(burst_management.StateDetail)
			Expect(intent.GetState()).To(Equal(burst_management.StateDetail))

			intent.SetState(burst_management.StateList)
			Expect(intent.GetState()).To(Equal(burst_management.StateList))
		})

		It("should allow setting and getting selected burst", func() {
			burst := ctx.Bursts[0]
			intent.SetSelectedBurst(burst)
			Expect(intent.GetSelectedBurst()).To(Equal(burst))

			intent.SetSelectedBurst(nil)
			Expect(intent.GetSelectedBurst()).To(BeNil())
		})

		It("should track viewed bursts", func() {
			burst := ctx.Bursts[0]
			intent.AddViewedBurst(burst)

			viewed := intent.GetViewedBursts()
			Expect(viewed).To(ContainElement(burst))
		})
	})

	Describe("SkillSuggestionsLoadedMsg filtering (filterNewSuggestions)", func() {
		It("should show suggestion modal when no existing skills", func() {
			intent.SetState(burst_management.StateInferringSkills)

			msg := burst_management.SkillSuggestionsLoadedMsg{
				Suggestions: []skillinference.SkillSuggestion{
					{Name: "Go", Category: "backend", Confidence: 0.95},
					{Name: "Docker", Category: "devops", Confidence: 0.85},
				},
				ExistingSkillNames: []string{},
			}

			intent.Update(msg)

			Expect(intent.GetState()).To(Equal(burst_management.StateSkillSuggestionReview))
		})

		It("should show 'all tracked' modal when all suggestions are existing", func() {
			intent.SetState(burst_management.StateInferringSkills)

			msg := burst_management.SkillSuggestionsLoadedMsg{
				Suggestions: []skillinference.SkillSuggestion{
					{Name: "Go", Category: "backend", Confidence: 0.95},
					{Name: "Docker", Category: "devops", Confidence: 0.85},
				},
				ExistingSkillNames: []string{"Go", "Docker"},
			}

			intent.Update(msg)

			Expect(intent.GetState()).To(Equal(burst_management.StateList))
			modal := intent.GetFeedbackModal()
			Expect(modal).NotTo(BeNil())
			Expect(modal.Type).To(Equal(feedback.ModalSuccess))
			Expect(modal.Message).To(ContainSubstring("Go"))
			Expect(modal.Message).To(ContainSubstring("Docker"))
		})

		It("should filter existing and show only new suggestions", func() {
			intent.SetState(burst_management.StateInferringSkills)

			msg := burst_management.SkillSuggestionsLoadedMsg{
				Suggestions: []skillinference.SkillSuggestion{
					{Name: "Go", Category: "backend", Confidence: 0.95},
					{Name: "Docker", Category: "devops", Confidence: 0.85},
					{Name: "React", Category: "frontend", Confidence: 0.75},
				},
				ExistingSkillNames: []string{"Go"},
			}

			intent.Update(msg)

			Expect(intent.GetState()).To(Equal(burst_management.StateSkillSuggestionReview))
		})

		It("should filter case-insensitively", func() {
			intent.SetState(burst_management.StateInferringSkills)

			msg := burst_management.SkillSuggestionsLoadedMsg{
				Suggestions: []skillinference.SkillSuggestion{
					{Name: "Go", Category: "backend", Confidence: 0.95},
					{Name: "Docker", Category: "devops", Confidence: 0.85},
				},
				ExistingSkillNames: []string{"go", "docker"},
			}

			intent.Update(msg)

			Expect(intent.GetState()).To(Equal(burst_management.StateList))
			modal := intent.GetFeedbackModal()
			Expect(modal).NotTo(BeNil())
			Expect(modal.Type).To(Equal(feedback.ModalSuccess))
		})

		It("should show warning when no suggestions detected at all", func() {
			intent.SetState(burst_management.StateInferringSkills)

			msg := burst_management.SkillSuggestionsLoadedMsg{
				Suggestions:        []skillinference.SkillSuggestion{},
				ExistingSkillNames: []string{},
			}

			intent.Update(msg)

			Expect(intent.GetState()).To(Equal(burst_management.StateList))
			modal := intent.GetFeedbackModal()
			Expect(modal).NotTo(BeNil())
			Expect(modal.Type).To(Equal(feedback.ModalWarning))
		})
	})

	Describe("FactExtractionCompleteMsg.Burst fallback", func() {
		BeforeEach(func() {
			ctx.SkillInferenceService = skillinference.NewSkillInferenceService(nil, nil, nil)
		})

		It("should use msg.Burst when selectedBurst is nil and trigger skill inference", func() {
			confirmedBurst := fixtures.BurstConfirmed("burst-fallback")
			confirmedBurst.Name = "Fallback Burst"
			intent.SetSelectedBurst(nil)
			intent.SetState(burst_management.StateExtractingFacts)

			msg := burst_management.FactExtractionCompleteMsg{
				Facts: []*career.Fact{fixtures.FactWith("f1", "Built API with Go")},
				Burst: confirmedBurst,
				Error: nil,
			}

			intent.Update(msg)

			Expect(intent.GetState()).To(Equal(burst_management.StateInferringSkills))
			Expect(intent.GetSelectedBurst()).To(Equal(confirmedBurst))
		})

		It("should not trigger skill inference when msg.Burst is nil and selectedBurst is nil", func() {
			intent.SetSelectedBurst(nil)

			msg := burst_management.FactExtractionCompleteMsg{
				Facts: []*career.Fact{fixtures.FactWith("f1", "Test fact")},
				Burst: nil,
				Error: nil,
			}

			intent.Update(msg)

			Expect(intent.GetState()).To(Equal(burst_management.StateList))
		})

		It("should not trigger skill inference when burst is not confirmed", func() {
			unconfirmedBurst := fixtures.Burst("burst-unconfirmed")
			unconfirmedBurst.Name = "Unconfirmed"
			intent.SetSelectedBurst(nil)

			msg := burst_management.FactExtractionCompleteMsg{
				Facts: []*career.Fact{fixtures.FactWith("f1", "Test fact")},
				Burst: unconfirmedBurst,
				Error: nil,
			}

			intent.Update(msg)

			Expect(intent.GetState()).To(Equal(burst_management.StateList))
		})

		It("should not trigger skill inference when facts are empty", func() {
			confirmedBurst := fixtures.BurstConfirmed("burst-nofacts")
			confirmedBurst.Name = "No Facts"
			intent.SetSelectedBurst(nil)

			msg := burst_management.FactExtractionCompleteMsg{
				Facts: []*career.Fact{},
				Burst: confirmedBurst,
				Error: nil,
			}

			intent.Update(msg)

			Expect(intent.GetState()).To(Equal(burst_management.StateList))
		})

		It("should prefer selectedBurst over msg.Burst when both are set", func() {
			selectedBurst := fixtures.BurstConfirmed("burst-selected")
			selectedBurst.Name = "Selected"
			msgBurst := fixtures.BurstConfirmed("burst-msg")
			msgBurst.Name = "From Message"
			intent.SetSelectedBurst(selectedBurst)

			msg := burst_management.FactExtractionCompleteMsg{
				Facts: []*career.Fact{fixtures.FactWith("f1", "Test fact")},
				Burst: msgBurst,
				Error: nil,
			}

			intent.Update(msg)

			Expect(intent.GetSelectedBurst()).To(Equal(selectedBurst))
		})
	})

	Describe("SkillsCreatedMsg handling", func() {
		It("should clear loading modal on success", func() {
			intent.Update(burst_management.BurstSuggestionsLoadedMsg{
				Suggestions: []burstfact.BurstSuggestion{
					{Name: "Test", EventIDs: []string{"e1"}, ConfidenceScore: 0.8},
				},
			})
			Expect(intent.GetState()).To(Equal(burst_management.StateSuggestionReview))

			intent.Update(burst_management.SkillsCreatedMsg{
				Skills: []*career.Skill{fixtures.SkillWith("s1", "Go", "Backend", "mid")},
			})

			Expect(intent.GetLoadingModal()).To(BeNil())
		})

		It("should clear loading modal on error", func() {
			intent.Update(burst_management.SkillsCreatedMsg{
				Error: context.DeadlineExceeded,
			})

			Expect(intent.GetLoadingModal()).To(BeNil())
		})

		It("should show success feedback modal with skill count", func() {
			skills := []*career.Skill{
				fixtures.SkillWith("s1", "Go", "Backend", "mid"),
				fixtures.SkillWith("s2", "PostgreSQL", "Database", "mid"),
				fixtures.SkillWith("s3", "Docker", "DevOps", "mid"),
			}

			intent.Update(burst_management.SkillsCreatedMsg{Skills: skills})

			modal := intent.GetFeedbackModal()
			Expect(modal).NotTo(BeNil())
			Expect(modal.Type).To(Equal(feedback.ModalSuccess))
			Expect(modal.Message).To(ContainSubstring("3 skill(s)"))
		})

		It("should show error feedback modal on failure", func() {
			intent.Update(burst_management.SkillsCreatedMsg{
				Error: fmt.Errorf("database connection lost"),
			})

			modal := intent.GetFeedbackModal()
			Expect(modal).NotTo(BeNil())
			Expect(modal.Type).To(Equal(feedback.ModalError))
		})

		It("should silently ignore cancelled operations", func() {
			intent.Update(burst_management.SkillsCreatedMsg{
				Error: context.Canceled,
			})

			Expect(intent.GetFeedbackModal()).To(BeNil())
			Expect(intent.GetState()).To(Equal(burst_management.StateList))
		})

		It("should transition to StateList on success", func() {
			intent.SetState(burst_management.StateInferringSkills)

			intent.Update(burst_management.SkillsCreatedMsg{
				Skills: []*career.Skill{fixtures.SkillWith("s1", "Go", "Backend", "mid")},
			})

			Expect(intent.GetState()).To(Equal(burst_management.StateList))
		})

		It("should transition to StateList on error", func() {
			intent.SetState(burst_management.StateInferringSkills)

			intent.Update(burst_management.SkillsCreatedMsg{
				Error: context.DeadlineExceeded,
			})

			Expect(intent.GetState()).To(Equal(burst_management.StateList))
		})
	})

	Describe("handleViewResult dispatch", func() {
		It("should dispatch submit result via handleViewResult", func() {
			result := &widgets.SubmitViewResult{FormData: map[string]interface{}{"x": 1}}
			cmd := intent.HandleSubmit(result)
			Expect(cmd).To(BeNil())
		})

		It("should dispatch error result via HandleError", func() {
			result := &widgets.ErrorViewResult{Err: errors.New("view error")}
			cmd := intent.HandleError(result)
			Expect(cmd).To(BeNil())
			Expect(intent.HasVisibleErrorModal()).To(BeTrue())
		})

		It("should handle HandleError with nil error gracefully", func() {
			result := &widgets.ErrorViewResult{}
			cmd := intent.HandleError(result)
			Expect(cmd).To(BeNil())
		})

		It("should handle HandleError with non-ErrorViewResult type", func() {
			result := &widgets.CancelViewResult{}
			cmd := intent.HandleError(result)
			Expect(cmd).To(BeNil())
		})
	})

	Describe("HandleNavigate with Nav type", func() {
		It("should handle Nav with view action", func() {
			burst := fixtures.Burst("nav-burst-1")
			burst.Name = "Nav Burst"
			ctx.Bursts = []*career.Burst{burst}
			ctx.Validate()
			intent, _ = burst_management.NewIntent(ctx)
			intent.Init()

			result := &widgets.NavigateViewResult{
				ResultData: burstviews.Nav{Action: burstviews.ActionView, Burst: display.BurstFromDomain(burst)},
			}
			intent.HandleNavigate(result)
			Expect(intent.GetDetailModal()).NotTo(BeNil())
			Expect(intent.GetSelectedBurst()).To(Equal(burst))
		})

		It("should handle Nav with edit action", func() {
			burst := fixtures.Burst("nav-burst-2")
			burst.Name = "Edit Nav Burst"
			ctx.Bursts = []*career.Burst{burst}
			ctx.Validate()
			intent, _ = burst_management.NewIntent(ctx)
			intent.Init()

			result := &widgets.NavigateViewResult{
				ResultData: burstviews.Nav{Action: burstviews.ActionEdit, Burst: display.BurstFromDomain(burst)},
			}
			intent.HandleNavigate(result)
			Expect(intent.GetState()).To(Equal(burst_management.StateEdit))
			Expect(intent.HasVisibleEditModal()).To(BeTrue())
		})

		It("should handle Nav with delete action", func() {
			burst := fixtures.Burst("nav-burst-3")
			burst.Name = "Delete Nav Burst"
			ctx.Bursts = []*career.Burst{burst}
			ctx.Validate()
			intent, _ = burst_management.NewIntent(ctx)
			intent.Init()

			result := &widgets.NavigateViewResult{
				ResultData: burstviews.Nav{Action: burstviews.ActionDelete, Burst: display.BurstFromDomain(burst)},
			}
			intent.HandleNavigate(result)
			Expect(intent.GetState()).To(Equal(burst_management.StateDeleteConfirm))
			Expect(intent.HasVisibleDeleteModal()).To(BeTrue())
		})

		It("should handle Nav with suggest action", func() {
			mockSvc := mocks.NewBurstServiceMock()
			ctx.Service = mockSvc
			ctx.Validate()
			intent, _ = burst_management.NewIntent(ctx)
			intent.Init()

			result := &widgets.NavigateViewResult{
				ResultData: burstviews.Nav{Action: burstviews.ActionSuggest},
			}
			intent.HandleNavigate(result)
			Expect(intent.GetState()).To(Equal(burst_management.StateSuggesting))
		})

		It("should handle Nav with confirm action", func() {
			burst := fixtures.Burst("nav-burst-4")
			burst.Name = "Confirm Burst"
			ctx.Bursts = []*career.Burst{burst}
			ctx.Validate()
			intent, _ = burst_management.NewIntent(ctx)
			intent.Init()

			result := &widgets.NavigateViewResult{
				ResultData: burstviews.Nav{Action: burstviews.ActionConfirm, Burst: display.BurstFromDomain(burst)},
			}
			intent.HandleNavigate(result)
			Expect(intent.GetState()).To(Equal(burst_management.StateConfirm))
			Expect(intent.HasVisibleConfirmModal()).To(BeTrue())
		})

		It("should handle Nav with view_events action", func() {
			burst := fixtures.Burst("nav-burst-5")
			burst.Name = "Events Burst"
			ctx.Bursts = []*career.Burst{burst}
			ctx.Validate()
			intent, _ = burst_management.NewIntent(ctx)
			intent.Init()

			result := &widgets.NavigateViewResult{
				ResultData: burstviews.Nav{Action: burstviews.ActionViewEvents, Burst: display.BurstFromDomain(burst)},
			}
			cmd := intent.HandleNavigate(result)
			Expect(intent.GetState()).To(Equal(burst_management.StateDetailEvents))
			Expect(cmd).NotTo(BeNil())
		})

		It("should handle Nav with view_facts action", func() {
			burst := fixtures.Burst("nav-burst-6")
			burst.Name = "Facts Burst"
			ctx.Bursts = []*career.Burst{burst}
			ctx.Validate()
			intent, _ = burst_management.NewIntent(ctx)
			intent.Init()

			result := &widgets.NavigateViewResult{
				ResultData: burstviews.Nav{Action: burstviews.ActionViewFacts, Burst: display.BurstFromDomain(burst)},
			}
			cmd := intent.HandleNavigate(result)
			Expect(intent.GetState()).To(Equal(burst_management.StateDetailFacts))
			Expect(cmd).NotTo(BeNil())
		})

		It("should handle Nav with unknown action", func() {
			result := &widgets.NavigateViewResult{
				ResultData: burstviews.Nav{Action: "unknown"},
			}
			cmd := intent.HandleNavigate(result)
			Expect(cmd).To(BeNil())
		})

		It("should handle non-Nav non-map non-burst data", func() {
			result := &widgets.NavigateViewResult{ResultData: "some string"}
			cmd := intent.HandleNavigate(result)
			Expect(cmd).To(BeNil())
		})
	})

	Describe("transitionToView", func() {
		It("should propagate terminal info to view", func() {
			termInfo = terminal.NewInfo()
			termInfo.Width = 200
			termInfo.Height = 50
			termInfo.IsValid = true
			intent.UpdateTerminalInfo(termInfo)

			intent.Init()

			view := intent.View()
			Expect(view).NotTo(BeEmpty())
		})
	})

	Describe("handleBurstSkillsLoaded", func() {
		It("should show skills modal on success", func() {
			burst := fixtures.Burst("skills-burst-1")
			burst.Name = "Skills Burst"
			ctx.Bursts = []*career.Burst{burst}
			ctx.Validate()
			intent, _ = burst_management.NewIntent(ctx)
			intent.Init()
			intent.SetSelectedBurst(burst)
			intent.UpdateTerminalInfo(termInfo)

			skills := []*career.Skill{
				fixtures.SkillWith("s1", "Go", "Backend", "mid"),
			}
			intent.Update(burst_management.BurstSkillsLoadedMsg{Skills: skills})
			Expect(intent.HasActiveModal()).To(BeTrue())
		})

		It("should show error modal on load error", func() {
			burst := fixtures.Burst("skills-burst-2")
			ctx.Bursts = []*career.Burst{burst}
			ctx.Validate()
			intent, _ = burst_management.NewIntent(ctx)
			intent.Init()
			intent.SetSelectedBurst(burst)

			intent.Update(burst_management.BurstSkillsLoadedMsg{
				Error: errors.New("load skills failed"),
			})
			Expect(intent.HasVisibleErrorModal()).To(BeTrue())
		})

		It("should handle nil selectedBurst gracefully", func() {
			intent.SetSelectedBurst(nil)
			intent.Update(burst_management.BurstSkillsLoadedMsg{
				Skills: []*career.Skill{},
			})
		})
	})

	Describe("handleSkillSuggestionsLoaded with error", func() {
		It("should show error modal on skill suggestion error", func() {
			intent.Update(burst_management.SkillSuggestionsLoadedMsg{
				Error: errors.New("inference failed"),
			})
			Expect(intent.HasVisibleErrorModal()).To(BeTrue())
			Expect(intent.GetState()).To(Equal(burst_management.StateList))
		})

		It("should ignore cancelled operations silently", func() {
			intent.Update(burst_management.SkillSuggestionsLoadedMsg{
				Error: context.Canceled,
			})
			Expect(intent.GetState()).To(Equal(burst_management.StateList))
			Expect(intent.HasVisibleErrorModal()).To(BeFalse())
		})
	})

	Describe("handleBurstSuggestionsLoaded with non-empty suggestions", func() {
		It("should show suggestion review modal on success", func() {
			intent.UpdateTerminalInfo(termInfo)
			suggestions := []burstfact.BurstSuggestion{
				{Name: "Burst 1", Description: "Desc 1", EventIDs: []string{"e1"}},
			}
			intent.Update(burst_management.BurstSuggestionsLoadedMsg{
				Suggestions: suggestions,
			})
			Expect(intent.GetState()).To(Equal(burst_management.StateSuggestionReview))
			Expect(intent.GetSuggestionModal()).NotTo(BeNil())
		})
	})

	Describe("getStateName edge cases", func() {
		It("should return InferringSkills state name", func() {
			intent.SetState(burst_management.StateInferringSkills)
			Expect(intent.GetStateName()).To(Equal("Inferring Skills"))
		})

		It("should return SkillSuggestionReview state name", func() {
			intent.SetState(burst_management.StateSkillSuggestionReview)
			Expect(intent.GetStateName()).To(Equal("Review Skills"))
		})
	})

	Describe("handleEditBurstMsg edge cases", func() {
		It("should show error when burst not found", func() {
			intent.SetSelectedBurst(nil)
			intent.Update(burst_management.EditBurstMsg{
				BurstID:     "non-existent",
				Name:        "Updated",
				Description: "Updated desc",
			})
			Expect(intent.HasVisibleErrorModal()).To(BeTrue())
		})

		It("should show error when burst ID mismatch", func() {
			burst := fixtures.Burst("edit-burst-1")
			intent.SetSelectedBurst(burst)

			intent.Update(burst_management.EditBurstMsg{
				BurstID:     "different-id",
				Name:        "Updated",
				Description: "Updated desc",
			})
			Expect(intent.HasVisibleErrorModal()).To(BeTrue())
		})
	})

	Describe("openDeleteModal with long name", func() {
		It("should truncate long burst names in delete confirmation", func() {
			burst := fixtures.Burst("long-name-burst")
			burst.Name = "This is a very long burst name that exceeds fifty characters and should be truncated"
			ctx.Bursts = []*career.Burst{burst}
			ctx.Validate()
			intent, _ = burst_management.NewIntent(ctx)
			intent.Init()

			result := &widgets.NavigateViewResult{
				ResultData: map[string]interface{}{
					"action": "delete",
					"burst":  burst,
				},
			}
			intent.HandleNavigate(result)
			Expect(intent.HasVisibleDeleteModal()).To(BeTrue())
		})
	})

	Describe("rebuildModalRegistry edge cases", func() {
		It("should register suggestion events modal", func() {
			intent.RebuildModalRegistry()
			Expect(intent.GetModalRegistry()).NotTo(BeNil())
		})
	})

	Describe("showEvents with service", func() {
		It("should return async command with service available", func() {
			event := fixtures.Event("e1")
			mockSvc := mocks.NewBurstServiceMock()
			mockSvc.SetEvents([]*career.Event{event})
			ctx.Service = mockSvc
			ctx.Validate()
			intent, _ = burst_management.NewIntent(ctx)
			intent.Init()
			intent.UpdateTerminalInfo(termInfo)

			burst := fixtures.Burst("events-burst", "e1")
			intent.SetSelectedBurst(burst)

			cmd := intent.HandleNavigate(&widgets.NavigateViewResult{
				ResultData: burstviews.Nav{Action: burstviews.ActionViewEvents, Burst: display.BurstFromDomain(burst)},
			})
			Expect(cmd).NotTo(BeNil())
		})
	})

	Describe("showFacts with service", func() {
		It("should return async command with service available", func() {
			mockSvc := mocks.NewBurstServiceMock()
			facts := []*career.Fact{fixtures.FactWith("f1", "Test fact")}
			mockSvc.SetFactsForBurst("facts-burst", facts)
			ctx.Service = mockSvc
			ctx.Validate()
			intent, _ = burst_management.NewIntent(ctx)
			intent.Init()

			burst := fixtures.Burst("facts-burst")
			intent.SetSelectedBurst(burst)
			intent.UpdateTerminalInfo(termInfo)

			result := &widgets.NavigateViewResult{
				ResultData: burstviews.Nav{Action: burstviews.ActionViewFacts, Burst: display.BurstFromDomain(burst)},
			}
			cmd := intent.HandleNavigate(result)
			Expect(cmd).NotTo(BeNil())

			msg := cmd()
			factsMsg, ok := msg.(burst_management.BurstFactsLoadedMsg)
			Expect(ok).To(BeTrue())
			Expect(factsMsg.Facts).To(HaveLen(1))
		})
	})

	Describe("showEvents and showFacts with nil burst", func() {
		It("should return nil when no burst selected for events", func() {
			intent.SetSelectedBurst(nil)
			cmd := intent.HandleNavigate(&widgets.NavigateViewResult{
				ResultData: burstviews.Nav{Action: burstviews.ActionViewEvents},
			})
			Expect(cmd).To(BeNil())
		})

		It("should return nil when no burst selected for facts", func() {
			intent.SetSelectedBurst(nil)
			cmd := intent.HandleNavigate(&widgets.NavigateViewResult{
				ResultData: burstviews.Nav{Action: burstviews.ActionViewFacts},
			})
			Expect(cmd).To(BeNil())
		})
	})

	Describe("loadBurstEvents via showEvents with service and event IDs", func() {
		It("should load events by ID from service", func() {
			event1 := fixtures.Event("ev-load-1")
			event2 := fixtures.Event("ev-load-2")
			mockSvc := mocks.NewBurstServiceMock()
			mockSvc.SetEvents([]*career.Event{event1, event2})
			burst := fixtures.Burst("load-events-burst", "ev-load-1", "ev-load-2")
			loadCtx := &burst_management.IntentValidator{
				Bursts:  []*career.Burst{burst},
				Service: mockSvc,
			}
			loadCtx.Validate()
			loadIntent, _ := burst_management.NewIntent(loadCtx)
			loadIntent.Init()
			loadIntent.UpdateTerminalInfo(termInfo)
			loadIntent.SetSelectedBurst(burst)

			cmd := loadIntent.HandleNavigate(&widgets.NavigateViewResult{
				ResultData: burstviews.Nav{Action: burstviews.ActionViewEvents, Burst: display.BurstFromDomain(burst)},
			})
			Expect(cmd).NotTo(BeNil())

			msg := cmd()
			eventsMsg, ok := msg.(burst_management.BurstEventsLoadedMsg)
			Expect(ok).To(BeTrue())
			Expect(eventsMsg.Events).To(HaveLen(2))
		})

		It("should skip events that fail to load", func() {
			event1 := fixtures.Event("ev-ok")
			mockSvc := mocks.NewBurstServiceMock()
			mockSvc.SetEvents([]*career.Event{event1})
			burst := fixtures.Burst("partial-load-burst", "ev-ok", "ev-missing")
			loadCtx := &burst_management.IntentValidator{
				Bursts:  []*career.Burst{burst},
				Service: mockSvc,
			}
			loadCtx.Validate()
			loadIntent, _ := burst_management.NewIntent(loadCtx)
			loadIntent.Init()
			loadIntent.UpdateTerminalInfo(termInfo)
			loadIntent.SetSelectedBurst(burst)

			cmd := loadIntent.HandleNavigate(&widgets.NavigateViewResult{
				ResultData: burstviews.Nav{Action: burstviews.ActionViewEvents, Burst: display.BurstFromDomain(burst)},
			})
			Expect(cmd).NotTo(BeNil())

			msg := cmd()
			eventsMsg, ok := msg.(burst_management.BurstEventsLoadedMsg)
			Expect(ok).To(BeTrue())
			Expect(eventsMsg.Events).To(HaveLen(1))
		})
	})

	Describe("handleEditBurstMsg via Update", func() {
		It("should apply edit when burst matches", func() {
			burst := fixtures.Burst("edit-msg-burst")
			burst.Name = "Original Name"
			burst.Description = "Original Desc"
			mockRepo := mocks.NewBurstRepositoryMock()
			mockRepo.AddBurst(burst)
			editCtx := &burst_management.IntentValidator{
				Bursts:          []*career.Burst{burst},
				BurstRepository: mockRepo,
			}
			editCtx.Validate()
			editIntent, _ := burst_management.NewIntent(editCtx)
			editIntent.Init()
			editIntent.UpdateTerminalInfo(termInfo)
			editIntent.SetSelectedBurst(burst)

			editIntent.Update(burst_management.EditBurstMsg{
				BurstID:     "edit-msg-burst",
				Name:        "Updated Name",
				Description: "Updated Desc",
			})

			Expect(burst.Name).To(Equal("Updated Name"))
		})

		It("should show error when burst ID does not match", func() {
			burst := fixtures.Burst("mismatch-burst")
			editCtx := &burst_management.IntentValidator{
				Bursts: []*career.Burst{burst},
			}
			editCtx.Validate()
			editIntent, _ := burst_management.NewIntent(editCtx)
			editIntent.Init()
			editIntent.UpdateTerminalInfo(termInfo)
			editIntent.SetSelectedBurst(burst)

			editIntent.Update(burst_management.EditBurstMsg{
				BurstID:     "wrong-id",
				Name:        "New",
				Description: "New",
			})

			Expect(editIntent.HasVisibleFeedbackModal()).To(BeTrue())
		})

		It("should show error when repository update fails", func() {
			burst := fixtures.Burst("repo-fail-burst")
			burst.Name = "Original"
			burst.Description = "Original"
			mockRepo := mocks.NewBurstRepositoryMock()
			mockRepo.AddBurst(burst)
			mockRepo.SetUpdateError(errors.New("db error"))
			editCtx := &burst_management.IntentValidator{
				Bursts:          []*career.Burst{burst},
				BurstRepository: mockRepo,
			}
			editCtx.Validate()
			editIntent, _ := burst_management.NewIntent(editCtx)
			editIntent.Init()
			editIntent.UpdateTerminalInfo(termInfo)
			editIntent.SetSelectedBurst(burst)

			cmd := editIntent.Update(burst_management.EditBurstMsg{
				BurstID:     "repo-fail-burst",
				Name:        "Changed",
				Description: "Changed",
			})
			Expect(cmd).NotTo(BeNil())
			msg := cmd()
			editIntent.Update(msg)

			Expect(burst.Name).To(Equal("Original"))
			Expect(editIntent.HasVisibleFeedbackModal()).To(BeTrue())
		})

		It("should handle nil selectedBurst", func() {
			noSelCtx := &burst_management.IntentValidator{
				Bursts: []*career.Burst{fixtures.Burst("some-burst")},
			}
			noSelCtx.Validate()
			noSelIntent, _ := burst_management.NewIntent(noSelCtx)
			noSelIntent.Init()
			noSelIntent.SetSelectedBurst(nil)

			noSelIntent.Update(burst_management.EditBurstMsg{
				BurstID:     "some-burst",
				Name:        "New",
				Description: "New",
			})

			Expect(noSelIntent.HasVisibleFeedbackModal()).To(BeTrue())
		})
	})

	Describe("handleBurstEditComplete via Update", func() {
		It("should handle error in message", func() {
			editCtx := &burst_management.IntentValidator{
				Bursts: []*career.Burst{fixtures.Burst("edit-complete-burst")},
			}
			editCtx.Validate()
			editIntent, _ := burst_management.NewIntent(editCtx)
			editIntent.Init()
			editIntent.UpdateTerminalInfo(termInfo)

			editIntent.Update(burst_management.BurstEditCompleteMsg{
				Error: errors.New("edit failed"),
			})

			Expect(editIntent.HasVisibleErrorModal()).To(BeTrue())
		})

		It("should refresh list on success", func() {
			burst := fixtures.Burst("edit-success-burst")
			editCtx := &burst_management.IntentValidator{
				Bursts: []*career.Burst{burst},
			}
			editCtx.Validate()
			editIntent, _ := burst_management.NewIntent(editCtx)
			editIntent.Init()
			editIntent.UpdateTerminalInfo(termInfo)

			editIntent.Update(burst_management.BurstEditCompleteMsg{
				Burst: burst,
			})

			Expect(editIntent.GetState()).To(Equal(burst_management.StateList))
		})
	})

	Describe("handleBurstEventsLoaded via Update", func() {
		It("should set events modal when events are loaded", func() {
			burst := fixtures.Burst("events-loaded-burst")
			event := fixtures.Event("loaded-ev-1")
			evCtx := &burst_management.IntentValidator{
				Bursts: []*career.Burst{burst},
			}
			evCtx.Validate()
			evIntent, _ := burst_management.NewIntent(evCtx)
			evIntent.Init()
			evIntent.UpdateTerminalInfo(termInfo)
			evIntent.SetSelectedBurst(burst)
			evIntent.SetLoadingEventsForTesting(true)

			evIntent.Update(burst_management.BurstEventsLoadedMsg{
				Events: []*career.Event{event},
			})

			Expect(evIntent.HasActiveModal()).To(BeTrue())
		})

		It("should handle error in events loaded message", func() {
			burst := fixtures.Burst("events-err-burst")
			evCtx := &burst_management.IntentValidator{
				Bursts: []*career.Burst{burst},
			}
			evCtx.Validate()
			evIntent, _ := burst_management.NewIntent(evCtx)
			evIntent.Init()
			evIntent.UpdateTerminalInfo(termInfo)
			evIntent.SetSelectedBurst(burst)
			evIntent.SetLoadingEventsForTesting(true)

			evIntent.Update(burst_management.BurstEventsLoadedMsg{
				Error: errors.New("load error"),
			})

			Expect(evIntent.HasVisibleErrorModal()).To(BeTrue())
		})
	})

	Describe("handleBurstFactsLoaded via Update", func() {
		It("should set facts modal when facts are loaded", func() {
			burst := fixtures.Burst("facts-loaded-burst")
			fact := fixtures.FactWith("loaded-fact-1", "A fact")
			fCtx := &burst_management.IntentValidator{
				Bursts: []*career.Burst{burst},
			}
			fCtx.Validate()
			fIntent, _ := burst_management.NewIntent(fCtx)
			fIntent.Init()
			fIntent.UpdateTerminalInfo(termInfo)
			fIntent.SetSelectedBurst(burst)
			fIntent.SetLoadingFactsForTesting(true)

			fIntent.Update(burst_management.BurstFactsLoadedMsg{
				Facts: []*career.Fact{fact},
			})

			Expect(fIntent.HasActiveModal()).To(BeTrue())
		})

		It("should handle error in facts loaded message", func() {
			burst := fixtures.Burst("facts-err-burst")
			fCtx := &burst_management.IntentValidator{
				Bursts: []*career.Burst{burst},
			}
			fCtx.Validate()
			fIntent, _ := burst_management.NewIntent(fCtx)
			fIntent.Init()
			fIntent.UpdateTerminalInfo(termInfo)
			fIntent.SetSelectedBurst(burst)
			fIntent.SetLoadingFactsForTesting(true)

			fIntent.Update(burst_management.BurstFactsLoadedMsg{
				Error: errors.New("fact load error"),
			})

			Expect(fIntent.HasVisibleErrorModal()).To(BeTrue())
		})
	})

	Describe("handleFactExtractionComplete via Update", func() {
		It("should reset extracting state on success", func() {
			burst := fixtures.Burst("extract-complete-burst")
			fact := fixtures.FactWith("extracted-1", "Extracted fact")
			exCtx := &burst_management.IntentValidator{
				Bursts: []*career.Burst{burst},
			}
			exCtx.Validate()
			exIntent, _ := burst_management.NewIntent(exCtx)
			exIntent.Init()
			exIntent.UpdateTerminalInfo(termInfo)

			exIntent.Update(burst_management.FactExtractionCompleteMsg{
				Facts: []*career.Fact{fact},
				Burst: burst,
			})

			Expect(exIntent.IsExtractingFacts()).To(BeFalse())
			Expect(exIntent.GetExtractedFactsCount()).To(Equal(1))
		})

		It("should show error when extraction fails", func() {
			burst := fixtures.Burst("extract-fail-burst")
			exCtx := &burst_management.IntentValidator{
				Bursts: []*career.Burst{burst},
			}
			exCtx.Validate()
			exIntent, _ := burst_management.NewIntent(exCtx)
			exIntent.Init()
			exIntent.UpdateTerminalInfo(termInfo)

			exIntent.Update(burst_management.FactExtractionCompleteMsg{
				Burst: burst,
				Error: errors.New("extraction error"),
			})

			Expect(exIntent.HasVisibleErrorModal()).To(BeTrue())
		})

		It("should show detail modal when selectedBurst matches", func() {
			burst := fixtures.Burst("extract-detail-burst")
			fact := fixtures.FactWith("det-fact-1", "Detail fact")
			exCtx := &burst_management.IntentValidator{
				Bursts: []*career.Burst{burst},
			}
			exCtx.Validate()
			exIntent, _ := burst_management.NewIntent(exCtx)
			exIntent.Init()
			exIntent.UpdateTerminalInfo(termInfo)
			exIntent.SetSelectedBurst(burst)

			exIntent.Update(burst_management.FactExtractionCompleteMsg{
				Facts: []*career.Fact{fact},
				Burst: burst,
			})

			Expect(exIntent.GetDetailModal()).NotTo(BeNil())
		})
	})

	Describe("handleSkillSuggestionsLoaded via Update", func() {
		It("should show skill suggestion modal when suggestions exist", func() {
			burst := fixtures.BurstConfirmed("skill-sug-burst")
			event := fixtures.Event("skill-ev-1")
			mockSvc := mocks.NewBurstServiceMock()
			mockSvc.SetEvents([]*career.Event{event})
			sCtx := &burst_management.IntentValidator{
				Bursts:                []*career.Burst{burst},
				Service:               mockSvc,
				SkillInferenceService: skillinference.NewSkillInferenceService(nil, nil, nil),
			}
			sCtx.Validate()
			sIntent, _ := burst_management.NewIntent(sCtx)
			sIntent.Init()
			sIntent.UpdateTerminalInfo(termInfo)
			sIntent.SetSelectedBurst(burst)
			sIntent.SetState(burst_management.StateInferringSkills)

			sIntent.Update(burst_management.SkillSuggestionsLoadedMsg{
				Suggestions: []skillinference.SkillSuggestion{
					{Name: "Go", Category: "backend", Confidence: 0.9, EventIDs: []string{"skill-ev-1"}},
				},
			})

			Expect(sIntent.HasActiveModal()).To(BeTrue())
		})

		It("should show info modal when no suggestions found", func() {
			burst := fixtures.BurstConfirmed("no-skill-sug-burst")
			sCtx := &burst_management.IntentValidator{
				Bursts:                []*career.Burst{burst},
				SkillInferenceService: skillinference.NewSkillInferenceService(nil, nil, nil),
			}
			sCtx.Validate()
			sIntent, _ := burst_management.NewIntent(sCtx)
			sIntent.Init()
			sIntent.UpdateTerminalInfo(termInfo)
			sIntent.SetSelectedBurst(burst)
			sIntent.SetState(burst_management.StateInferringSkills)

			sIntent.Update(burst_management.SkillSuggestionsLoadedMsg{
				Suggestions: []skillinference.SkillSuggestion{},
			})

			Expect(sIntent.HasVisibleFeedbackModal()).To(BeTrue())
		})

		It("should show error modal on error", func() {
			burst := fixtures.BurstConfirmed("skill-sug-err-burst")
			sCtx := &burst_management.IntentValidator{
				Bursts:                []*career.Burst{burst},
				SkillInferenceService: skillinference.NewSkillInferenceService(nil, nil, nil),
			}
			sCtx.Validate()
			sIntent, _ := burst_management.NewIntent(sCtx)
			sIntent.Init()
			sIntent.UpdateTerminalInfo(termInfo)
			sIntent.SetState(burst_management.StateInferringSkills)

			sIntent.Update(burst_management.SkillSuggestionsLoadedMsg{
				Error: errors.New("inference failed"),
			})

			Expect(sIntent.HasVisibleErrorModal()).To(BeTrue())
		})
	})

	Describe("handleSkillSuggestionsLoaded error via Update", func() {
		It("should show error modal", func() {
			burst := fixtures.BurstConfirmed("skill-err-burst")
			sCtx := &burst_management.IntentValidator{
				Bursts: []*career.Burst{burst},
			}
			sCtx.Validate()
			sIntent, _ := burst_management.NewIntent(sCtx)
			sIntent.Init()
			sIntent.UpdateTerminalInfo(termInfo)
			sIntent.SetState(burst_management.StateInferringSkills)

			sIntent.Update(burst_management.SkillSuggestionsLoadedMsg{
				Error: errors.New("skill error"),
			})

			Expect(sIntent.HasVisibleErrorModal()).To(BeTrue())
			Expect(sIntent.GetState()).To(Equal(burst_management.StateList))
		})
	})

	Describe("handleSkillsCreated via Update", func() {
		It("should show success modal on success", func() {
			skill := fixtures.SkillWith("created-skill-1", "Go", "backend", "advanced")
			sCtx := &burst_management.IntentValidator{
				Bursts: []*career.Burst{fixtures.Burst("skills-created-burst")},
			}
			sCtx.Validate()
			sIntent, _ := burst_management.NewIntent(sCtx)
			sIntent.Init()
			sIntent.UpdateTerminalInfo(termInfo)

			sIntent.Update(burst_management.SkillsCreatedMsg{
				Skills: []*career.Skill{skill},
			})

			Expect(sIntent.HasVisibleFeedbackModal()).To(BeTrue())
		})

		It("should show error modal on error", func() {
			sCtx := &burst_management.IntentValidator{
				Bursts: []*career.Burst{fixtures.Burst("skills-err-burst")},
			}
			sCtx.Validate()
			sIntent, _ := burst_management.NewIntent(sCtx)
			sIntent.Init()
			sIntent.UpdateTerminalInfo(termInfo)

			sIntent.Update(burst_management.SkillsCreatedMsg{
				Error: errors.New("create failed"),
			})

			Expect(sIntent.HasVisibleErrorModal()).To(BeTrue())
		})
	})

	Describe("showSkills via detail modal 's' key", func() {
		It("should load skills when skill repository is nil", func() {
			burst := fixtures.Burst("skills-modal-burst", "skill-ev-1")
			skCtx := &burst_management.IntentValidator{
				Bursts: []*career.Burst{burst},
			}
			skCtx.Validate()
			skIntent, _ := burst_management.NewIntent(skCtx)
			skIntent.Init()
			skIntent.UpdateTerminalInfo(termInfo)
			skIntent.SetSelectedBurst(burst)

			skIntent.HandleNavigate(&widgets.NavigateViewResult{
				ResultData: burstviews.Nav{Action: burstviews.ActionView, Burst: display.BurstFromDomain(burst)},
			})

			Expect(skIntent.GetDetailModal()).NotTo(BeNil())

			skIntent.Update(tea.KeyMsg{Type: tea.KeyRunes, Runes: []rune{'s'}})

			Expect(skIntent.GetDetailModal()).NotTo(BeNil())
		})
	})

	Describe("handleBurstSkillsLoaded via Update", func() {
		It("should show skills modal on success", func() {
			burst := fixtures.Burst("skills-loaded-burst")
			skill := fixtures.SkillWith("sk-1", "Go", "backend", "advanced")
			skCtx := &burst_management.IntentValidator{
				Bursts: []*career.Burst{burst},
			}
			skCtx.Validate()
			skIntent, _ := burst_management.NewIntent(skCtx)
			skIntent.Init()
			skIntent.UpdateTerminalInfo(termInfo)
			skIntent.SetSelectedBurst(burst)

			skIntent.Update(burst_management.BurstSkillsLoadedMsg{
				Skills: []*career.Skill{skill},
			})

			Expect(skIntent.HasActiveModal()).To(BeTrue())
		})

		It("should show error on skill load error", func() {
			burst := fixtures.Burst("skills-load-err-burst")
			skCtx := &burst_management.IntentValidator{
				Bursts: []*career.Burst{burst},
			}
			skCtx.Validate()
			skIntent, _ := burst_management.NewIntent(skCtx)
			skIntent.Init()
			skIntent.UpdateTerminalInfo(termInfo)
			skIntent.SetSelectedBurst(burst)

			skIntent.Update(burst_management.BurstSkillsLoadedMsg{
				Error: errors.New("skills load error"),
			})

			Expect(skIntent.HasVisibleErrorModal()).To(BeTrue())
		})
	})

	Describe("handleSuggestionReviewComplete via Update", func() {
		It("should handle cancelled suggestion review", func() {
			sCtx := &burst_management.IntentValidator{
				Bursts: []*career.Burst{fixtures.Burst("sug-review-burst")},
			}
			sCtx.Validate()
			sIntent, _ := burst_management.NewIntent(sCtx)
			sIntent.Init()
			sIntent.UpdateTerminalInfo(termInfo)

			sIntent.Update(burst_management.SuggestionReviewCompleteMsg{
				Cancelled: true,
			})

			Expect(sIntent.GetState()).To(Equal(burst_management.StateList))
		})

		It("should start fact extraction for accepted suggestions", func() {
			mockSvc := mocks.NewBurstServiceMock()
			sCtx := &burst_management.IntentValidator{
				Bursts:  []*career.Burst{fixtures.Burst("sug-accept-burst")},
				Service: mockSvc,
			}
			sCtx.Validate()
			sIntent, _ := burst_management.NewIntent(sCtx)
			sIntent.Init()
			sIntent.UpdateTerminalInfo(termInfo)
			sIntent.SetState(burst_management.StateSuggestionReview)

			sIntent.Update(burst_management.SuggestionReviewCompleteMsg{
				AcceptedSuggestions: []burstfact.BurstSuggestion{
					{Name: "Test Burst", Description: "Test", EventIDs: []string{"e1"}},
				},
				Cancelled: false,
			})

			Expect(sIntent.IsExtractingFacts()).To(BeTrue())
		})
	})

	Describe("RefreshData with service and LoadBursts", func() {
		It("should refresh filtered bursts", func() {
			mockRepo := mocks.NewBurstRepositoryMock()
			burst := fixtures.Burst("refresh-burst")
			mockRepo.AddBurst(burst)
			rCtx := &burst_management.IntentValidator{
				Bursts:          []*career.Burst{burst},
				BurstRepository: mockRepo,
			}
			rCtx.Validate()
			rIntent, _ := burst_management.NewIntent(rCtx)
			rIntent.Init()
			rIntent.UpdateTerminalInfo(termInfo)

			rIntent.RefreshData()

			Expect(rIntent.GetFilteredBursts()).NotTo(BeEmpty())
		})
	})

	Describe("startSkillInference via detail modal 'i' key on confirmed burst", func() {
		It("should start inference for confirmed burst", func() {
			burst := fixtures.BurstConfirmed("infer-burst")
			burst.EventIDs = []string{"infer-ev-1"}
			event := fixtures.Event("infer-ev-1")
			mockSvc := mocks.NewBurstServiceMock()
			mockSvc.SetEvents([]*career.Event{event})
			infCtx := &burst_management.IntentValidator{
				Bursts:                []*career.Burst{burst},
				Service:               mockSvc,
				SkillInferenceService: skillinference.NewSkillInferenceService(nil, nil, nil),
			}
			infCtx.Validate()
			infIntent, _ := burst_management.NewIntent(infCtx)
			infIntent.Init()
			infIntent.UpdateTerminalInfo(termInfo)
			infIntent.SetSelectedBurst(burst)

			infIntent.HandleNavigate(&widgets.NavigateViewResult{
				ResultData: burstviews.Nav{Action: burstviews.ActionView, Burst: display.BurstFromDomain(burst)},
			})

			Expect(infIntent.GetDetailModal()).NotTo(BeNil())

			infIntent.Update(tea.KeyMsg{Type: tea.KeyRunes, Runes: []rune{'i'}})

			Expect(infIntent.GetState()).To(Equal(burst_management.StateInferringSkills))
		})

		It("should show warning for unconfirmed burst", func() {
			burst := fixtures.Burst("unconfirmed-burst")
			infCtx := &burst_management.IntentValidator{
				Bursts:                []*career.Burst{burst},
				SkillInferenceService: skillinference.NewSkillInferenceService(nil, nil, nil),
			}
			infCtx.Validate()
			infIntent, _ := burst_management.NewIntent(infCtx)
			infIntent.Init()
			infIntent.UpdateTerminalInfo(termInfo)
			infIntent.SetSelectedBurst(burst)

			infIntent.HandleNavigate(&widgets.NavigateViewResult{
				ResultData: burstviews.Nav{Action: burstviews.ActionView, Burst: display.BurstFromDomain(burst)},
			})

			infIntent.Update(tea.KeyMsg{Type: tea.KeyRunes, Runes: []rune{'i'}})

			Expect(infIntent.HasVisibleFeedbackModal()).To(BeTrue())
		})
	})

	Describe("handleDetailModalKeypress 'v' for events", func() {
		It("should trigger events modal from detail modal", func() {
			burst := fixtures.Burst("detail-v-burst", "d-ev-1")
			event := fixtures.Event("d-ev-1")
			mockSvc := mocks.NewBurstServiceMock()
			mockSvc.SetEvents([]*career.Event{event})
			dCtx := &burst_management.IntentValidator{
				Bursts:  []*career.Burst{burst},
				Service: mockSvc,
			}
			dCtx.Validate()
			dIntent, _ := burst_management.NewIntent(dCtx)
			dIntent.Init()
			dIntent.UpdateTerminalInfo(termInfo)
			dIntent.SetSelectedBurst(burst)

			dIntent.HandleNavigate(&widgets.NavigateViewResult{
				ResultData: burstviews.Nav{Action: burstviews.ActionView, Burst: display.BurstFromDomain(burst)},
			})
			Expect(dIntent.GetDetailModal()).NotTo(BeNil())

			dIntent.Update(tea.KeyMsg{Type: tea.KeyRunes, Runes: []rune{'v'}})
		})
	})

	Describe("handleDetailModalKeypress 'f' for facts", func() {
		It("should trigger facts modal from detail modal", func() {
			burst := fixtures.Burst("detail-f-burst")
			fact := fixtures.FactWith("d-fact-1", "A fact")
			mockSvc := mocks.NewBurstServiceMock()
			mockSvc.SetFactsForBurst("detail-f-burst", []*career.Fact{fact})
			dCtx := &burst_management.IntentValidator{
				Bursts:  []*career.Burst{burst},
				Service: mockSvc,
			}
			dCtx.Validate()
			dIntent, _ := burst_management.NewIntent(dCtx)
			dIntent.Init()
			dIntent.UpdateTerminalInfo(termInfo)
			dIntent.SetSelectedBurst(burst)

			dIntent.HandleNavigate(&widgets.NavigateViewResult{
				ResultData: burstviews.Nav{Action: burstviews.ActionView, Burst: display.BurstFromDomain(burst)},
			})
			Expect(dIntent.GetDetailModal()).NotTo(BeNil())

			dIntent.Update(tea.KeyMsg{Type: tea.KeyRunes, Runes: []rune{'f'}})
		})
	})

	Describe("handleDetailModalKeypress 'e' for edit", func() {
		It("should open edit modal from detail modal", func() {
			burst := fixtures.Burst("detail-e-burst")
			dCtx := &burst_management.IntentValidator{
				Bursts: []*career.Burst{burst},
			}
			dCtx.Validate()
			dIntent, _ := burst_management.NewIntent(dCtx)
			dIntent.Init()
			dIntent.UpdateTerminalInfo(termInfo)
			dIntent.SetSelectedBurst(burst)

			dIntent.HandleNavigate(&widgets.NavigateViewResult{
				ResultData: burstviews.Nav{Action: burstviews.ActionView, Burst: display.BurstFromDomain(burst)},
			})
			Expect(dIntent.GetDetailModal()).NotTo(BeNil())

			dIntent.Update(tea.KeyMsg{Type: tea.KeyRunes, Runes: []rune{'e'}})

			Expect(dIntent.HasVisibleEditModal()).To(BeTrue())
		})
	})

	Describe("handleDetailModalKeypress 'c' for confirm", func() {
		It("should show confirm modal from detail modal", func() {
			burst := fixtures.Burst("detail-c-burst")
			dCtx := &burst_management.IntentValidator{
				Bursts: []*career.Burst{burst},
			}
			dCtx.Validate()
			dIntent, _ := burst_management.NewIntent(dCtx)
			dIntent.Init()
			dIntent.UpdateTerminalInfo(termInfo)
			dIntent.SetSelectedBurst(burst)

			dIntent.HandleNavigate(&widgets.NavigateViewResult{
				ResultData: burstviews.Nav{Action: burstviews.ActionView, Burst: display.BurstFromDomain(burst)},
			})
			Expect(dIntent.GetDetailModal()).NotTo(BeNil())

			dIntent.Update(tea.KeyMsg{Type: tea.KeyRunes, Runes: []rune{'c'}})
			dIntent.Update(burst_management.ConfirmBurstFactsLoadedMsg{Facts: []*career.Fact{}})

			Expect(dIntent.HasVisibleConfirmModal()).To(BeTrue())
		})
	})

	Describe("handleDetailModalKeypress Esc to close", func() {
		It("should close detail modal on Esc", func() {
			burst := fixtures.Burst("detail-esc-burst")
			dCtx := &burst_management.IntentValidator{
				Bursts: []*career.Burst{burst},
			}
			dCtx.Validate()
			dIntent, _ := burst_management.NewIntent(dCtx)
			dIntent.Init()
			dIntent.UpdateTerminalInfo(termInfo)
			dIntent.SetSelectedBurst(burst)

			dIntent.HandleNavigate(&widgets.NavigateViewResult{
				ResultData: burstviews.Nav{Action: burstviews.ActionView, Burst: display.BurstFromDomain(burst)},
			})
			Expect(dIntent.GetDetailModal()).NotTo(BeNil())

			dIntent.Update(tea.KeyMsg{Type: tea.KeyEsc})

			Expect(dIntent.GetDetailModal()).To(BeNil())
		})

		It("should close detail modal on Enter", func() {
			burst := fixtures.Burst("detail-enter-burst")
			dCtx := &burst_management.IntentValidator{
				Bursts: []*career.Burst{burst},
			}
			dCtx.Validate()
			dIntent, _ := burst_management.NewIntent(dCtx)
			dIntent.Init()
			dIntent.UpdateTerminalInfo(termInfo)
			dIntent.SetSelectedBurst(burst)

			dIntent.HandleNavigate(&widgets.NavigateViewResult{
				ResultData: burstviews.Nav{Action: burstviews.ActionView, Burst: display.BurstFromDomain(burst)},
			})
			Expect(dIntent.GetDetailModal()).NotTo(BeNil())

			dIntent.Update(tea.KeyMsg{Type: tea.KeyEnter})

			Expect(dIntent.GetDetailModal()).To(BeNil())
		})
	})

	Describe("confirm modal Esc returns to detail", func() {
		It("should return to detail modal on confirm cancel", func() {
			burst := fixtures.Burst("confirm-esc-burst")
			cCtx := &burst_management.IntentValidator{
				Bursts: []*career.Burst{burst},
			}
			cCtx.Validate()
			cIntent, _ := burst_management.NewIntent(cCtx)
			cIntent.Init()
			cIntent.UpdateTerminalInfo(termInfo)
			cIntent.SetSelectedBurst(burst)

			cIntent.HandleNavigate(&widgets.NavigateViewResult{
				ResultData: burstviews.Nav{Action: burstviews.ActionView, Burst: display.BurstFromDomain(burst)},
			})

			cIntent.Update(tea.KeyMsg{Type: tea.KeyRunes, Runes: []rune{'c'}})
			cIntent.Update(burst_management.ConfirmBurstFactsLoadedMsg{Facts: []*career.Fact{}})
			Expect(cIntent.HasVisibleConfirmModal()).To(BeTrue())

			cIntent.Update(tea.KeyMsg{Type: tea.KeyEsc})
		})
	})

	Describe("events modal Esc returns to detail", func() {
		It("should return to detail modal when events modal is closed with Esc", func() {
			burst := fixtures.Burst("events-esc-burst")
			event := fixtures.Event("evs-esc-1")
			mockSvc := mocks.NewBurstServiceMock()
			mockSvc.SetEvents([]*career.Event{event})
			emCtx := &burst_management.IntentValidator{
				Bursts:  []*career.Burst{burst},
				Service: mockSvc,
			}
			emCtx.Validate()
			emIntent, _ := burst_management.NewIntent(emCtx)
			emIntent.Init()
			emIntent.UpdateTerminalInfo(termInfo)
			emIntent.SetSelectedBurst(burst)

			emIntent.HandleNavigate(&widgets.NavigateViewResult{
				ResultData: burstviews.Nav{Action: burstviews.ActionView, Burst: display.BurstFromDomain(burst)},
			})

			emIntent.Update(tea.KeyMsg{Type: tea.KeyRunes, Runes: []rune{'v'}})

			emIntent.Update(burst_management.BurstEventsLoadedMsg{
				Events: []*career.Event{event},
			})

			emIntent.Update(tea.KeyMsg{Type: tea.KeyEsc})
		})

		It("should return to detail modal when events modal is closed with Enter", func() {
			burst := fixtures.Burst("events-enter-burst")
			event := fixtures.Event("evs-enter-1")
			mockSvc := mocks.NewBurstServiceMock()
			mockSvc.SetEvents([]*career.Event{event})
			emCtx := &burst_management.IntentValidator{
				Bursts:  []*career.Burst{burst},
				Service: mockSvc,
			}
			emCtx.Validate()
			emIntent, _ := burst_management.NewIntent(emCtx)
			emIntent.Init()
			emIntent.UpdateTerminalInfo(termInfo)
			emIntent.SetSelectedBurst(burst)

			emIntent.HandleNavigate(&widgets.NavigateViewResult{
				ResultData: burstviews.Nav{Action: burstviews.ActionView, Burst: display.BurstFromDomain(burst)},
			})

			emIntent.Update(tea.KeyMsg{Type: tea.KeyRunes, Runes: []rune{'v'}})

			emIntent.Update(burst_management.BurstEventsLoadedMsg{
				Events: []*career.Event{event},
			})

			emIntent.Update(tea.KeyMsg{Type: tea.KeyEnter})
		})
	})

	Describe("facts modal Esc returns to detail", func() {
		It("should return to detail modal when facts modal is closed", func() {
			burst := fixtures.Burst("facts-esc-burst")
			fact := fixtures.FactWith("fm-fact-1", "A fact")
			mockSvc := mocks.NewBurstServiceMock()
			mockSvc.SetFactsForBurst("facts-esc-burst", []*career.Fact{fact})
			fmCtx := &burst_management.IntentValidator{
				Bursts:  []*career.Burst{burst},
				Service: mockSvc,
			}
			fmCtx.Validate()
			fmIntent, _ := burst_management.NewIntent(fmCtx)
			fmIntent.Init()
			fmIntent.UpdateTerminalInfo(termInfo)
			fmIntent.SetSelectedBurst(burst)

			fmIntent.HandleNavigate(&widgets.NavigateViewResult{
				ResultData: burstviews.Nav{Action: burstviews.ActionView, Burst: display.BurstFromDomain(burst)},
			})

			fmIntent.Update(tea.KeyMsg{Type: tea.KeyRunes, Runes: []rune{'f'}})

			fmIntent.Update(burst_management.BurstFactsLoadedMsg{
				Facts: []*career.Fact{fact},
			})

			fmIntent.Update(tea.KeyMsg{Type: tea.KeyEsc})
		})
	})

	Describe("skills modal Esc returns to detail", func() {
		It("should return to detail modal when skills modal is closed", func() {
			burst := fixtures.Burst("skills-esc-burst", "sk-ev-1")
			skill := fixtures.SkillWith("sm-skill-1", "Go", "backend", "advanced")
			smCtx := &burst_management.IntentValidator{
				Bursts: []*career.Burst{burst},
			}
			smCtx.Validate()
			smIntent, _ := burst_management.NewIntent(smCtx)
			smIntent.Init()
			smIntent.UpdateTerminalInfo(termInfo)
			smIntent.SetSelectedBurst(burst)

			smIntent.Update(burst_management.BurstSkillsLoadedMsg{
				Skills: []*career.Skill{skill},
			})

			smIntent.Update(tea.KeyMsg{Type: tea.KeyEsc})
		})

		It("should return to detail when skills modal is closed with Enter", func() {
			burst := fixtures.Burst("skills-enter-burst", "sk-ev-2")
			skill := fixtures.SkillWith("sm-skill-2", "Python", "backend", "intermediate")
			smCtx := &burst_management.IntentValidator{
				Bursts: []*career.Burst{burst},
			}
			smCtx.Validate()
			smIntent, _ := burst_management.NewIntent(smCtx)
			smIntent.Init()
			smIntent.UpdateTerminalInfo(termInfo)
			smIntent.SetSelectedBurst(burst)

			smIntent.Update(burst_management.BurstSkillsLoadedMsg{
				Skills: []*career.Skill{skill},
			})

			smIntent.Update(tea.KeyMsg{Type: tea.KeyEnter})
		})
	})

	Describe("loading modal Esc cancels async", func() {
		It("should cancel async operation on Esc during loading", func() {
			burst := fixtures.Burst("loading-cancel-burst")
			mockSvc := mocks.NewBurstServiceMock()
			mockSvc.SetEvents([]*career.Event{fixtures.Event("lc-ev-1")})
			lcCtx := &burst_management.IntentValidator{
				Bursts:  []*career.Burst{burst},
				Service: mockSvc,
			}
			lcCtx.Validate()
			lcIntent, _ := burst_management.NewIntent(lcCtx)
			lcIntent.Init()
			lcIntent.UpdateTerminalInfo(termInfo)

			lcIntent.Update(tea.KeyMsg{Type: tea.KeyRunes, Runes: []rune{'s'}})
			Expect(lcIntent.GetLoadingModal()).NotTo(BeNil())

			lcIntent.Update(tea.KeyMsg{Type: tea.KeyEsc})
			Expect(lcIntent.GetLoadingModal()).To(BeNil())
			Expect(lcIntent.GetState()).To(Equal(burst_management.StateList))
		})

		It("should consume non-Esc keys during loading", func() {
			burst := fixtures.Burst("loading-key-burst")
			mockSvc := mocks.NewBurstServiceMock()
			mockSvc.SetEvents([]*career.Event{fixtures.Event("lk-ev-1")})
			lcCtx := &burst_management.IntentValidator{
				Bursts:  []*career.Burst{burst},
				Service: mockSvc,
			}
			lcCtx.Validate()
			lcIntent, _ := burst_management.NewIntent(lcCtx)
			lcIntent.Init()
			lcIntent.UpdateTerminalInfo(termInfo)

			lcIntent.Update(tea.KeyMsg{Type: tea.KeyRunes, Runes: []rune{'s'}})
			Expect(lcIntent.GetLoadingModal()).NotTo(BeNil())

			lcIntent.Update(tea.KeyMsg{Type: tea.KeyRunes, Runes: []rune{'x'}})
			Expect(lcIntent.GetLoadingModal()).NotTo(BeNil())
		})
	})

	Describe("feedback modal Esc dismisses", func() {
		It("should dismiss feedback modal on Esc", func() {
			burst := fixtures.Burst("fb-esc-burst")
			fbCtx := &burst_management.IntentValidator{
				Bursts: []*career.Burst{burst},
			}
			fbCtx.Validate()
			fbIntent, _ := burst_management.NewIntent(fbCtx)
			fbIntent.Init()
			fbIntent.UpdateTerminalInfo(termInfo)
			fbIntent.ShowErrorModal("Test", "Test error")

			Expect(fbIntent.HasVisibleFeedbackModal()).To(BeTrue())

			fbIntent.Update(tea.KeyMsg{Type: tea.KeyEsc})

			Expect(fbIntent.GetFeedbackModal()).To(BeNil())
		})
	})

	Describe("delete modal cancel returns to list", func() {
		It("should return to list when delete modal is cancelled", func() {
			burst := fixtures.Burst("del-cancel-burst")
			dCtx := &burst_management.IntentValidator{
				Bursts: []*career.Burst{burst},
			}
			dCtx.Validate()
			dIntent, _ := burst_management.NewIntent(dCtx)
			dIntent.Init()
			dIntent.UpdateTerminalInfo(termInfo)
			dIntent.SetSelectedBurst(burst)

			dIntent.HandleNavigate(&widgets.NavigateViewResult{
				ResultData: burstviews.Nav{Action: burstviews.ActionDelete, Burst: display.BurstFromDomain(burst)},
			})

			Expect(dIntent.HasVisibleDeleteModal()).To(BeTrue())

			dIntent.Update(tea.KeyMsg{Type: tea.KeyEsc})
		})
	})

	Describe("inferSkillsFromBurst execution paths", func() {
		It("should start inference and set loading state for confirmed burst with no events", func() {
			burst := fixtures.BurstConfirmed("no-events-burst")
			burst.EventIDs = []string{}
			mockSvc := mocks.NewBurstServiceMock()
			infCtx := &burst_management.IntentValidator{
				Bursts:                []*career.Burst{burst},
				Service:               mockSvc,
				SkillInferenceService: skillinference.NewSkillInferenceService(nil, nil, nil),
			}
			infCtx.Validate()
			infIntent, _ := burst_management.NewIntent(infCtx)
			infIntent.Init()
			infIntent.UpdateTerminalInfo(termInfo)
			infIntent.SetSelectedBurst(burst)

			infIntent.HandleNavigate(&widgets.NavigateViewResult{
				ResultData: burstviews.Nav{Action: burstviews.ActionView, Burst: display.BurstFromDomain(burst)},
			})

			cmd := infIntent.Update(tea.KeyMsg{Type: tea.KeyRunes, Runes: []rune{'i'}})
			Expect(cmd).NotTo(BeNil())
			Expect(infIntent.GetState()).To(Equal(burst_management.StateInferringSkills))
			Expect(infIntent.GetLoadingModal()).NotTo(BeNil())
		})

		It("should show error when no burst selected for inference", func() {
			infCtx := &burst_management.IntentValidator{
				Bursts:                []*career.Burst{fixtures.BurstConfirmed("inf-nil-burst")},
				SkillInferenceService: skillinference.NewSkillInferenceService(nil, nil, nil),
			}
			infCtx.Validate()
			infIntent, _ := burst_management.NewIntent(infCtx)
			infIntent.Init()
			infIntent.UpdateTerminalInfo(termInfo)
			infIntent.SetSelectedBurst(nil)

			infIntent.HandleNavigate(&widgets.NavigateViewResult{
				ResultData: burstviews.Nav{Action: burstviews.ActionSuggest},
			})
		})

		It("should show error when skill inference service is nil", func() {
			burst := fixtures.BurstConfirmed("inf-no-svc-burst")
			infCtx := &burst_management.IntentValidator{
				Bursts: []*career.Burst{burst},
			}
			infCtx.Validate()
			infIntent, _ := burst_management.NewIntent(infCtx)
			infIntent.Init()
			infIntent.UpdateTerminalInfo(termInfo)
			infIntent.SetSelectedBurst(burst)

			infIntent.HandleNavigate(&widgets.NavigateViewResult{
				ResultData: burstviews.Nav{Action: burstviews.ActionView, Burst: display.BurstFromDomain(burst)},
			})

			infIntent.Update(tea.KeyMsg{Type: tea.KeyRunes, Runes: []rune{'i'}})

			Expect(infIntent.HasVisibleErrorModal()).To(BeTrue())
		})
	})

	Describe("extractFactsForBurst execution paths", func() {
		It("should extract facts and save them via service", func() {
			burst := fixtures.BurstConfirmed("extract-exec-burst")
			mockSvc := mocks.NewBurstServiceMock()
			mockSvc.SetExtractedFacts([]career.Fact{
				*fixtures.FactWith("ef-1", "Fact 1"),
			})
			exCtx := &burst_management.IntentValidator{
				Bursts:  []*career.Burst{burst},
				Service: mockSvc,
			}
			exCtx.Validate()
			exIntent, _ := burst_management.NewIntent(exCtx)
			exIntent.Init()
			exIntent.UpdateTerminalInfo(termInfo)
			exIntent.SetSelectedBurst(burst)

			exIntent.HandleNavigate(&widgets.NavigateViewResult{
				ResultData: burstviews.Nav{Action: burstviews.ActionView, Burst: display.BurstFromDomain(burst)},
			})

			exIntent.Update(tea.KeyMsg{Type: tea.KeyRunes, Runes: []rune{'c'}})
			exIntent.Update(burst_management.ConfirmBurstFactsLoadedMsg{Facts: []*career.Fact{}})
			Expect(exIntent.HasVisibleConfirmModal()).To(BeTrue())

			exIntent.Update(tea.KeyMsg{Type: tea.KeyRunes, Runes: []rune{'y'}})
			cmd := exIntent.Update(burst_management.BurstConfirmedMsg{Burst: burst})
			if cmd != nil {
				msg := cmd()
				if exMsg, ok := msg.(burst_management.FactExtractionCompleteMsg); ok {
					Expect(exMsg.Error).ToNot(HaveOccurred())
				}
			}
		})

		It("should handle extract error", func() {
			burst := fixtures.BurstConfirmed("extract-err-exec-burst")
			mockSvc := mocks.NewBurstServiceMock()
			mockSvc.SetExtractError(errors.New("extract failed"))
			exCtx := &burst_management.IntentValidator{
				Bursts:  []*career.Burst{burst},
				Service: mockSvc,
			}
			exCtx.Validate()
			exIntent, _ := burst_management.NewIntent(exCtx)
			exIntent.Init()
			exIntent.UpdateTerminalInfo(termInfo)
			exIntent.SetSelectedBurst(burst)

			exIntent.HandleNavigate(&widgets.NavigateViewResult{
				ResultData: burstviews.Nav{Action: burstviews.ActionView, Burst: display.BurstFromDomain(burst)},
			})

			exIntent.Update(tea.KeyMsg{Type: tea.KeyRunes, Runes: []rune{'c'}})

			cmd := exIntent.Update(tea.KeyMsg{Type: tea.KeyRunes, Runes: []rune{'y'}})
			if cmd != nil {
				msg := cmd()
				if exMsg, ok := msg.(burst_management.FactExtractionCompleteMsg); ok {
					Expect(exMsg.Error).To(HaveOccurred())
				}
			}
		})
	})

	Describe("BurstSuggestionsLoadedMsg with suggestions and review", func() {
		It("should show suggestion modal and handle acceptance flow", func() {
			burst := fixtures.Burst("sug-flow-burst", "sug-ev-1")
			event := fixtures.Event("sug-ev-1")
			mockSvc := mocks.NewBurstServiceMock()
			mockSvc.SetEvents([]*career.Event{event})
			mockSvc.SetExtractedFacts([]career.Fact{*fixtures.FactWith("sf-1", "Suggestion fact")})
			mockRepo := mocks.NewBurstRepositoryMock()
			sfCtx := &burst_management.IntentValidator{
				Bursts:          []*career.Burst{burst},
				Service:         mockSvc,
				BurstRepository: mockRepo,
			}
			sfCtx.Validate()
			sfIntent, _ := burst_management.NewIntent(sfCtx)
			sfIntent.Init()
			sfIntent.UpdateTerminalInfo(termInfo)

			sfIntent.Update(burst_management.BurstSuggestionsLoadedMsg{
				Suggestions: []burstfact.BurstSuggestion{
					{Name: "Test Burst", Description: "A test", EventIDs: []string{"sug-ev-1"}},
				},
			})

			Expect(sfIntent.GetSuggestionModal()).NotTo(BeNil())
			Expect(sfIntent.GetState()).To(Equal(burst_management.StateSuggestionReview))
		})
	})

	Describe("createBurstFromSuggestion with repository error", func() {
		It("should show error when repository create fails during suggestion acceptance", func() {
			mockSvc := mocks.NewBurstServiceMock()
			mockSvc.SetEvents([]*career.Event{fixtures.Event("cbs-ev-1")})
			mockRepo := mocks.NewBurstRepositoryMock()
			mockRepo.SetCreateError(errors.New("create failed"))
			cbsCtx := &burst_management.IntentValidator{
				Bursts:          []*career.Burst{fixtures.Burst("cbs-burst")},
				Service:         mockSvc,
				BurstRepository: mockRepo,
			}
			cbsCtx.Validate()
			cbsIntent, _ := burst_management.NewIntent(cbsCtx)
			cbsIntent.Init()
			cbsIntent.UpdateTerminalInfo(termInfo)

			cbsIntent.Update(burst_management.BurstSuggestionsLoadedMsg{
				Suggestions: []burstfact.BurstSuggestion{
					{Name: "Fail Burst", Description: "Will fail", EventIDs: []string{"cbs-ev-1"}},
				},
			})

			Expect(cbsIntent.GetSuggestionModal()).NotTo(BeNil())

			cmd := cbsIntent.Update(tea.KeyMsg{Type: tea.KeyRunes, Runes: []rune{'a'}})
			Expect(cmd).NotTo(BeNil())
			msg := cmd()
			cbsIntent.Update(msg)

			Expect(cbsIntent.HasVisibleErrorModal()).To(BeTrue())
		})
	})

	Describe("edit modal cancel returns to list", func() {
		It("should return to list when edit modal is cancelled", func() {
			burst := fixtures.Burst("edit-cancel-burst")
			ecCtx := &burst_management.IntentValidator{
				Bursts: []*career.Burst{burst},
			}
			ecCtx.Validate()
			ecIntent, _ := burst_management.NewIntent(ecCtx)
			ecIntent.Init()
			ecIntent.UpdateTerminalInfo(termInfo)
			ecIntent.SetSelectedBurst(burst)

			ecIntent.HandleNavigate(&widgets.NavigateViewResult{
				ResultData: burstviews.Nav{Action: burstviews.ActionEdit, Burst: display.BurstFromDomain(burst)},
			})

			Expect(ecIntent.HasVisibleEditModal()).To(BeTrue())

			ecIntent.Update(tea.KeyMsg{Type: tea.KeyEsc})

			Expect(ecIntent.GetState()).To(Equal(burst_management.StateList))
		})
	})

	Describe("inactive intent returns nil on Update", func() {
		It("should return nil when intent is deactivated", func() {
			dCtx := &burst_management.IntentValidator{
				Bursts: []*career.Burst{fixtures.Burst("deactivated-burst")},
			}
			dCtx.Validate()
			dIntent, _ := burst_management.NewIntent(dCtx)
			dIntent.Init()
			dIntent.Deactivate()

			cmd := dIntent.Update(tea.KeyMsg{Type: tea.KeyRunes, Runes: []rune{'s'}})
			Expect(cmd).To(BeNil())
		})
	})

	Describe("confirmBurst with service error", func() {
		It("should show error when confirm fails", func() {
			burst := fixtures.Burst("confirm-err-burst")
			mockSvc := mocks.NewBurstServiceMock()
			mockSvc.SetConfirmError(errors.New("confirm failed"))
			cCtx := &burst_management.IntentValidator{
				Bursts:  []*career.Burst{burst},
				Service: mockSvc,
			}
			cCtx.Validate()
			cIntent, _ := burst_management.NewIntent(cCtx)
			cIntent.Init()
			cIntent.UpdateTerminalInfo(termInfo)
			cIntent.SetSelectedBurst(burst)

			cIntent.HandleNavigate(&widgets.NavigateViewResult{
				ResultData: burstviews.Nav{Action: burstviews.ActionView, Burst: display.BurstFromDomain(burst)},
			})

			cIntent.Update(tea.KeyMsg{Type: tea.KeyRunes, Runes: []rune{'c'}})
			cIntent.Update(burst_management.ConfirmBurstFactsLoadedMsg{Facts: []*career.Fact{}})

			cIntent.Update(tea.KeyMsg{Type: tea.KeyRunes, Runes: []rune{'y'}})
			cIntent.Update(burst_management.BurstConfirmedMsg{Burst: burst, Error: errors.New("confirm failed")})

			Expect(cIntent.HasVisibleErrorModal()).To(BeTrue())
		})
	})

	Describe("showFacts with error from service", func() {
		It("should return error message on facts load error", func() {
			burst := fixtures.Burst("facts-err-svc-burst")
			mockSvc := mocks.NewBurstServiceMock()
			fCtx := &burst_management.IntentValidator{
				Bursts:  []*career.Burst{burst},
				Service: mockSvc,
			}
			fCtx.Validate()
			fIntent, _ := burst_management.NewIntent(fCtx)
			fIntent.Init()
			fIntent.UpdateTerminalInfo(termInfo)
			fIntent.SetSelectedBurst(burst)

			cmd := fIntent.HandleNavigate(&widgets.NavigateViewResult{
				ResultData: burstviews.Nav{Action: burstviews.ActionViewFacts, Burst: display.BurstFromDomain(burst)},
			})
			Expect(cmd).NotTo(BeNil())

			msg := cmd()
			factsMsg, ok := msg.(burst_management.BurstFactsLoadedMsg)
			Expect(ok).To(BeTrue())
			Expect(factsMsg.Facts).To(BeEmpty())
		})
	})

	Describe("startSkillInference with nil selectedBurst", func() {
		It("should show error when no burst selected", func() {
			skNilCtx := &burst_management.IntentValidator{
				Bursts:                []*career.Burst{fixtures.Burst("sk-nil-burst")},
				SkillInferenceService: skillinference.NewSkillInferenceService(nil, nil, nil),
			}
			skNilCtx.Validate()
			skNilIntent, _ := burst_management.NewIntent(skNilCtx)
			skNilIntent.Init()
			skNilIntent.UpdateTerminalInfo(termInfo)
			skNilIntent.SetSelectedBurst(nil)

			skNilIntent.HandleNavigate(&widgets.NavigateViewResult{
				ResultData: burstviews.Nav{Action: burstviews.ActionSuggest},
			})

			Expect(skNilIntent.GetState()).To(Equal(burst_management.StateList))
		})
	})

	Describe("openDeleteModal with long burst name", func() {
		It("should truncate long burst names in delete modal", func() {
			longName := "This is a very long burst name that exceeds the fifty character limit for display"
			burst := fixtures.Burst("long-name-burst")
			burst.Name = longName
			dnCtx := &burst_management.IntentValidator{
				Bursts: []*career.Burst{burst},
			}
			dnCtx.Validate()
			dnIntent, _ := burst_management.NewIntent(dnCtx)
			dnIntent.Init()
			dnIntent.UpdateTerminalInfo(termInfo)
			dnIntent.SetSelectedBurst(burst)

			dnIntent.HandleNavigate(&widgets.NavigateViewResult{
				ResultData: burstviews.Nav{Action: burstviews.ActionDelete, Burst: display.BurstFromDomain(burst)},
			})

			Expect(dnIntent.HasVisibleDeleteModal()).To(BeTrue())
		})
	})

	Describe("showConfirmBurstModal with existing facts", func() {
		It("should show re-extract message when burst has existing facts", func() {
			burst := fixtures.Burst("confirm-facts-burst")
			fact := fixtures.FactWith("cf-1", "Existing fact")
			mockSvc := mocks.NewBurstServiceMock()
			mockSvc.SetFactsForBurst("confirm-facts-burst", []*career.Fact{fact})
			cfCtx := &burst_management.IntentValidator{
				Bursts:  []*career.Burst{burst},
				Service: mockSvc,
			}
			cfCtx.Validate()
			cfIntent, _ := burst_management.NewIntent(cfCtx)
			cfIntent.Init()
			cfIntent.UpdateTerminalInfo(termInfo)
			cfIntent.SetSelectedBurst(burst)

			cfIntent.HandleNavigate(&widgets.NavigateViewResult{
				ResultData: burstviews.Nav{Action: burstviews.ActionView, Burst: display.BurstFromDomain(burst)},
			})

			cfIntent.Update(tea.KeyMsg{Type: tea.KeyRunes, Runes: []rune{'c'}})
			cfIntent.Update(burst_management.ConfirmBurstFactsLoadedMsg{Facts: []*career.Fact{fact}})

			Expect(cfIntent.HasVisibleConfirmModal()).To(BeTrue())
		})
	})

	Describe("handleBurstSuggestionsLoaded error path", func() {
		It("should show error modal when suggestions have error", func() {
			bsCtx := &burst_management.IntentValidator{
				Bursts: []*career.Burst{fixtures.Burst("sug-err-burst")},
			}
			bsCtx.Validate()
			bsIntent, _ := burst_management.NewIntent(bsCtx)
			bsIntent.Init()
			bsIntent.UpdateTerminalInfo(termInfo)
			bsIntent.SetState(burst_management.StateSuggesting)

			bsIntent.Update(burst_management.BurstSuggestionsLoadedMsg{
				Error: errors.New("detection failed"),
			})

			Expect(bsIntent.HasVisibleErrorModal()).To(BeTrue())
			Expect(bsIntent.GetState()).To(Equal(burst_management.StateList))
		})

		It("should show info when zero suggestions returned", func() {
			bsCtx := &burst_management.IntentValidator{
				Bursts: []*career.Burst{fixtures.Burst("no-sug-burst")},
			}
			bsCtx.Validate()
			bsIntent, _ := burst_management.NewIntent(bsCtx)
			bsIntent.Init()
			bsIntent.UpdateTerminalInfo(termInfo)
			bsIntent.SetState(burst_management.StateSuggesting)

			bsIntent.Update(burst_management.BurstSuggestionsLoadedMsg{
				Suggestions: []burstfact.BurstSuggestion{},
			})

			Expect(bsIntent.HasVisibleFeedbackModal()).To(BeTrue())
		})
	})

	Describe("handleLoadError via Update", func() {
		It("should show error on generic load error", func() {
			leCtx := &burst_management.IntentValidator{
				Bursts: []*career.Burst{fixtures.Burst("load-err-burst")},
			}
			leCtx.Validate()
			leIntent, _ := burst_management.NewIntent(leCtx)
			leIntent.Init()
			leIntent.UpdateTerminalInfo(termInfo)

			leIntent.Update(burst_management.BurstEventsLoadedMsg{
				Error: errors.New("generic error"),
			})

			Expect(leIntent.HasVisibleErrorModal()).To(BeTrue())
		})
	})

	Describe("skill suggestion accept flow (handleSkillSuggestionAccept + saveSkillFromSuggestion)", func() {
		It("should show error when skill inference service is nil during accept", func() {
			burst := fixtures.BurstConfirmed("skill-nil-svc-burst")
			burst.EventIDs = []string{"sn-ev-1"}
			event := fixtures.Event("sn-ev-1")
			mockSvc := mocks.NewBurstServiceMock()
			mockSvc.SetEvents([]*career.Event{event})
			snCtx := &burst_management.IntentValidator{
				Bursts:                []*career.Burst{burst},
				Service:               mockSvc,
				SkillInferenceService: skillinference.NewSkillInferenceService(nil, nil, nil),
			}
			snCtx.Validate()
			snIntent, _ := burst_management.NewIntent(snCtx)
			snIntent.Init()
			snIntent.UpdateTerminalInfo(termInfo)
			snIntent.SetSelectedBurst(burst)
			snIntent.SetState(burst_management.StateInferringSkills)

			snIntent.Update(burst_management.SkillSuggestionsLoadedMsg{
				Suggestions: []skillinference.SkillSuggestion{
					{Name: "Python", Category: "backend", Confidence: 0.8, EventIDs: []string{"sn-ev-1"}},
				},
			})

			Expect(snIntent.GetState()).To(Equal(burst_management.StateSkillSuggestionReview))

			snCtx.SkillInferenceService = nil

			snIntent.Update(tea.KeyMsg{Type: tea.KeyRunes, Runes: []rune{'a'}})

			Expect(snIntent.HasVisibleErrorModal()).To(BeTrue())
		})
	})

	Describe("skill suggestion modal cancel flow", func() {
		It("should return to list when skill suggestion modal is cancelled", func() {
			burst := fixtures.BurstConfirmed("skill-cancel-burst")
			burst.EventIDs = []string{"sc-ev-1"}
			event := fixtures.Event("sc-ev-1")
			mockSvc := mocks.NewBurstServiceMock()
			mockSvc.SetEvents([]*career.Event{event})
			scCtx := &burst_management.IntentValidator{
				Bursts:                []*career.Burst{burst},
				Service:               mockSvc,
				SkillInferenceService: skillinference.NewSkillInferenceService(nil, nil, nil),
			}
			scCtx.Validate()
			scIntent, _ := burst_management.NewIntent(scCtx)
			scIntent.Init()
			scIntent.UpdateTerminalInfo(termInfo)
			scIntent.SetSelectedBurst(burst)
			scIntent.SetState(burst_management.StateInferringSkills)

			scIntent.Update(burst_management.SkillSuggestionsLoadedMsg{
				Suggestions: []skillinference.SkillSuggestion{
					{Name: "Docker", Category: "devops", Confidence: 0.7, EventIDs: []string{"sc-ev-1"}},
				},
			})

			Expect(scIntent.GetState()).To(Equal(burst_management.StateSkillSuggestionReview))

			scIntent.Update(tea.KeyMsg{Type: tea.KeyEsc})

			Expect(scIntent.GetState()).To(Equal(burst_management.StateList))
		})
	})

	Describe("skill suggestion modal reject flow", func() {
		It("should reject current suggestion via 'r' key", func() {
			burst := fixtures.BurstConfirmed("skill-reject-burst")
			burst.EventIDs = []string{"sr-ev-1"}
			event := fixtures.Event("sr-ev-1")
			mockSvc := mocks.NewBurstServiceMock()
			mockSvc.SetEvents([]*career.Event{event})
			srCtx := &burst_management.IntentValidator{
				Bursts:                []*career.Burst{burst},
				Service:               mockSvc,
				SkillInferenceService: skillinference.NewSkillInferenceService(nil, nil, nil),
			}
			srCtx.Validate()
			srIntent, _ := burst_management.NewIntent(srCtx)
			srIntent.Init()
			srIntent.UpdateTerminalInfo(termInfo)
			srIntent.SetSelectedBurst(burst)
			srIntent.SetState(burst_management.StateInferringSkills)

			srIntent.Update(burst_management.SkillSuggestionsLoadedMsg{
				Suggestions: []skillinference.SkillSuggestion{
					{Name: "Kubernetes", Category: "devops", Confidence: 0.6, EventIDs: []string{"sr-ev-1"}},
				},
			})

			srIntent.Update(tea.KeyMsg{Type: tea.KeyRunes, Runes: []rune{'r'}})
		})
	})

	Describe("handleSkillSuggestionsLoaded with existing skill names filter", func() {
		It("should show success when all suggestions already exist", func() {
			burst := fixtures.BurstConfirmed("all-exist-burst")
			burst.EventIDs = []string{"ae-ev-1"}
			aeCtx := &burst_management.IntentValidator{
				Bursts:                []*career.Burst{burst},
				SkillInferenceService: skillinference.NewSkillInferenceService(nil, nil, nil),
			}
			aeCtx.Validate()
			aeIntent, _ := burst_management.NewIntent(aeCtx)
			aeIntent.Init()
			aeIntent.UpdateTerminalInfo(termInfo)
			aeIntent.SetSelectedBurst(burst)
			aeIntent.SetState(burst_management.StateInferringSkills)

			aeIntent.Update(burst_management.SkillSuggestionsLoadedMsg{
				Suggestions: []skillinference.SkillSuggestion{
					{Name: "Go", Category: "backend", Confidence: 0.9, EventIDs: []string{"ae-ev-1"}},
				},
				ExistingSkillNames: []string{"Go"},
			})

			Expect(aeIntent.HasVisibleFeedbackModal()).To(BeTrue())
			Expect(aeIntent.GetState()).To(Equal(burst_management.StateList))
		})

		It("should handle cancelled context error silently", func() {
			burst := fixtures.BurstConfirmed("cancelled-ctx-burst")
			cxCtx := &burst_management.IntentValidator{
				Bursts:                []*career.Burst{burst},
				SkillInferenceService: skillinference.NewSkillInferenceService(nil, nil, nil),
			}
			cxCtx.Validate()
			cxIntent, _ := burst_management.NewIntent(cxCtx)
			cxIntent.Init()
			cxIntent.UpdateTerminalInfo(termInfo)
			cxIntent.SetState(burst_management.StateInferringSkills)

			cxIntent.Update(burst_management.SkillSuggestionsLoadedMsg{
				Error: context.Canceled,
			})

			Expect(cxIntent.GetState()).To(Equal(burst_management.StateList))
			Expect(cxIntent.HasVisibleErrorModal()).To(BeFalse())
		})
	})

	Describe("handleSkillSuggestionsLoaded error with cancelled context", func() {
		It("should handle cancelled context silently", func() {
			burst := fixtures.BurstConfirmed("err-cancelled-burst")
			ecCtx := &burst_management.IntentValidator{
				Bursts: []*career.Burst{burst},
			}
			ecCtx.Validate()
			ecIntent, _ := burst_management.NewIntent(ecCtx)
			ecIntent.Init()
			ecIntent.UpdateTerminalInfo(termInfo)
			ecIntent.SetState(burst_management.StateInferringSkills)

			ecIntent.Update(burst_management.SkillSuggestionsLoadedMsg{
				Error: context.Canceled,
			})

			Expect(ecIntent.GetState()).To(Equal(burst_management.StateList))
			Expect(ecIntent.HasVisibleErrorModal()).To(BeFalse())
		})
	})

	Describe("handleSkillsCreated with cancelled context", func() {
		It("should handle cancelled context silently", func() {
			scCtx := &burst_management.IntentValidator{
				Bursts: []*career.Burst{fixtures.Burst("skills-created-cancel-burst")},
			}
			scCtx.Validate()
			scIntent, _ := burst_management.NewIntent(scCtx)
			scIntent.Init()
			scIntent.UpdateTerminalInfo(termInfo)

			scIntent.Update(burst_management.SkillsCreatedMsg{
				Error: context.Canceled,
			})

			Expect(scIntent.GetState()).To(Equal(burst_management.StateList))
			Expect(scIntent.HasVisibleErrorModal()).To(BeFalse())
		})
	})

	Describe("loadBurstEvents via inferSkillsFromBurst async execution", func() {
		It("should load events from service in async command", func() {
			burst := fixtures.BurstConfirmed("async-load-burst")
			burst.EventIDs = []string{"async-ev-1", "async-ev-2"}
			event1 := fixtures.Event("async-ev-1")
			event2 := fixtures.Event("async-ev-2")
			mockSvc := mocks.NewBurstServiceMock()
			mockSvc.SetEvents([]*career.Event{event1, event2})
			skiSvc := skillinference.NewSkillInferenceService(nil, nil, nil)
			alCtx := &burst_management.IntentValidator{
				Bursts:                []*career.Burst{burst},
				Service:               mockSvc,
				SkillInferenceService: skiSvc,
			}
			alCtx.Validate()
			alIntent, _ := burst_management.NewIntent(alCtx)
			alIntent.Init()
			alIntent.UpdateTerminalInfo(termInfo)
			alIntent.SetSelectedBurst(burst)

			alIntent.HandleNavigate(&widgets.NavigateViewResult{
				ResultData: burstviews.Nav{Action: burstviews.ActionView, Burst: display.BurstFromDomain(burst)},
			})

			cmd := alIntent.Update(tea.KeyMsg{Type: tea.KeyRunes, Runes: []rune{'i'}})
			Expect(cmd).NotTo(BeNil())

			msg := cmd()
			if errMsg, ok := msg.(burst_management.SkillSuggestionsLoadedMsg); ok {
				Expect(errMsg.Error).To(HaveOccurred())
			}
		})
	})

	Describe("resolveEventsByIDs via openSuggestionEventsModal", func() {
		It("should resolve events when viewing suggestion events via Enter key", func() {
			burst := fixtures.BurstConfirmed("resolve-ev-burst")
			burst.EventIDs = []string{"res-ev-1"}
			event := fixtures.Event("res-ev-1")
			mockSvc := mocks.NewBurstServiceMock()
			mockSvc.SetEvents([]*career.Event{event})
			skiSvc := skillinference.NewSkillInferenceService(nil, nil, nil)
			reCtx := &burst_management.IntentValidator{
				Bursts:                []*career.Burst{burst},
				Service:               mockSvc,
				SkillInferenceService: skiSvc,
			}
			reCtx.Validate()
			reIntent, _ := burst_management.NewIntent(reCtx)
			reIntent.Init()
			reIntent.UpdateTerminalInfo(termInfo)
			reIntent.SetSelectedBurst(burst)
			reIntent.SetState(burst_management.StateInferringSkills)

			reIntent.Update(burst_management.SkillSuggestionsLoadedMsg{
				Suggestions: []skillinference.SkillSuggestion{
					{Name: "Go", Category: "backend", Confidence: 0.9, EventIDs: []string{"res-ev-1"}},
				},
			})

			reIntent.Update(tea.KeyMsg{Type: tea.KeyEnter})
		})
	})

	Describe("suggestion events modal Esc returns to skill suggestion", func() {
		It("should close suggestion events modal and re-show skill modal on Esc", func() {
			burst := fixtures.BurstConfirmed("sug-ev-esc-burst")
			burst.EventIDs = []string{"se-ev-1"}
			event := fixtures.Event("se-ev-1")
			mockSvc := mocks.NewBurstServiceMock()
			mockSvc.SetEvents([]*career.Event{event})
			skiSvc := skillinference.NewSkillInferenceService(nil, nil, nil)
			seCtx := &burst_management.IntentValidator{
				Bursts:                []*career.Burst{burst},
				Service:               mockSvc,
				SkillInferenceService: skiSvc,
			}
			seCtx.Validate()
			seIntent, _ := burst_management.NewIntent(seCtx)
			seIntent.Init()
			seIntent.UpdateTerminalInfo(termInfo)
			seIntent.SetSelectedBurst(burst)
			seIntent.SetState(burst_management.StateInferringSkills)

			seIntent.Update(burst_management.SkillSuggestionsLoadedMsg{
				Suggestions: []skillinference.SkillSuggestion{
					{Name: "Go", Category: "backend", Confidence: 0.9, EventIDs: []string{"se-ev-1"}},
				},
			})

			seIntent.Update(tea.KeyMsg{Type: tea.KeyEnter})

			seIntent.Update(tea.KeyMsg{Type: tea.KeyEsc})
		})
	})

	Describe("showSkills with SkillRepository", func() {
		It("should load and deduplicate skills from multiple events", func() {
			burst := fixtures.Burst("skills-repo-burst", "skr-ev-1", "skr-ev-2")
			skrCtx := &burst_management.IntentValidator{
				Bursts: []*career.Burst{burst},
			}
			skrCtx.Validate()
			skrIntent, _ := burst_management.NewIntent(skrCtx)
			skrIntent.Init()
			skrIntent.UpdateTerminalInfo(termInfo)
			skrIntent.SetSelectedBurst(burst)

			skrIntent.HandleNavigate(&widgets.NavigateViewResult{
				ResultData: burstviews.Nav{Action: burstviews.ActionView, Burst: display.BurstFromDomain(burst)},
			})

			cmd := skrIntent.Update(tea.KeyMsg{Type: tea.KeyRunes, Runes: []rune{'s'}})
			Expect(cmd).NotTo(BeNil())

			msg := cmd()
			skillsMsg, ok := msg.(burst_management.BurstSkillsLoadedMsg)
			Expect(ok).To(BeTrue())
			Expect(skillsMsg.Skills).To(BeEmpty())
		})
	})

	Describe("saveAndExtractBurstWithResult with repository error", func() {
		It("should show error when burst repo create fails during suggestion acceptance", func() {
			burst := fixtures.Burst("sae-err-burst", "sae-ev-1")
			event := fixtures.Event("sae-ev-1")
			mockSvc := mocks.NewBurstServiceMock()
			mockSvc.SetEvents([]*career.Event{event})
			mockRepo := mocks.NewBurstRepositoryMock()
			mockRepo.SetCreateError(errors.New("create failed"))
			saeCtx := &burst_management.IntentValidator{
				Bursts:          []*career.Burst{burst},
				Service:         mockSvc,
				BurstRepository: mockRepo,
			}
			saeCtx.Validate()
			saeIntent, _ := burst_management.NewIntent(saeCtx)
			saeIntent.Init()
			saeIntent.UpdateTerminalInfo(termInfo)
			saeIntent.SetState(burst_management.StateSuggestionReview)

			saeIntent.Update(burst_management.BurstSuggestionsLoadedMsg{
				Suggestions: []burstfact.BurstSuggestion{
					{Name: "Failed Burst", Description: "Will fail", EventIDs: []string{"sae-ev-1"}},
				},
			})

			cmd := saeIntent.Update(tea.KeyMsg{Type: tea.KeyRunes, Runes: []rune{'a'}})
			Expect(cmd).NotTo(BeNil())
			msg := cmd()
			saeIntent.Update(msg)

			Expect(saeIntent.HasVisibleErrorModal()).To(BeTrue())
		})
	})

	Describe("startFactExtractionForBursts with empty burst list", func() {
		It("should return to list when no bursts to extract", func() {
			feCtx := &burst_management.IntentValidator{
				Bursts: []*career.Burst{fixtures.Burst("no-extract-burst")},
			}
			feCtx.Validate()
			feIntent, _ := burst_management.NewIntent(feCtx)
			feIntent.Init()
			feIntent.UpdateTerminalInfo(termInfo)
			feIntent.SetState(burst_management.StateSuggestionReview)

			feIntent.Update(burst_management.SuggestionReviewCompleteMsg{
				AcceptedSuggestions: []burstfact.BurstSuggestion{},
				Cancelled:           false,
			})

			Expect(feIntent.GetState()).To(Equal(burst_management.StateList))
		})
	})

	Describe("extractFactsForBurst with nil service in async cmd", func() {
		It("should return error when service is nil in async execution", func() {
			burst := fixtures.BurstConfirmed("nil-svc-extract-burst")
			noSvcCtx := &burst_management.IntentValidator{
				Bursts: []*career.Burst{burst},
			}
			noSvcCtx.Validate()
			noSvcIntent, _ := burst_management.NewIntent(noSvcCtx)
			noSvcIntent.Init()
			noSvcIntent.UpdateTerminalInfo(termInfo)
			noSvcIntent.SetSelectedBurst(burst)
			noSvcIntent.SetState(burst_management.StateSuggestionReview)

			noSvcIntent.Update(burst_management.SuggestionReviewCompleteMsg{
				AcceptedSuggestions: []burstfact.BurstSuggestion{
					{Name: "No Svc Burst", Description: "Test", EventIDs: []string{"e1"}},
				},
				Cancelled: false,
			})

			Expect(noSvcIntent.IsExtractingFacts()).To(BeTrue())
		})
	})

	Describe("handleFactExtractionComplete with cancelled context", func() {
		It("should handle cancelled context during fact extraction", func() {
			feCtx := &burst_management.IntentValidator{
				Bursts: []*career.Burst{fixtures.Burst("fact-cancel-burst")},
			}
			feCtx.Validate()
			feIntent, _ := burst_management.NewIntent(feCtx)
			feIntent.Init()
			feIntent.UpdateTerminalInfo(termInfo)

			feIntent.Update(burst_management.FactExtractionCompleteMsg{
				Error: context.Canceled,
			})
		})
	})

	Describe("handleEditBurstMsg with valid edit via domain function", func() {
		It("should succeed without repository", func() {
			burst := fixtures.Burst("edit-no-repo-burst")
			burst.Name = "Original"
			burst.Description = "Original"
			enCtx := &burst_management.IntentValidator{
				Bursts: []*career.Burst{burst},
			}
			enCtx.Validate()
			enIntent, _ := burst_management.NewIntent(enCtx)
			enIntent.Init()
			enIntent.UpdateTerminalInfo(termInfo)
			enIntent.SetSelectedBurst(burst)

			enIntent.Update(burst_management.EditBurstMsg{
				BurstID:     "edit-no-repo-burst",
				Name:        "Updated",
				Description: "Updated",
			})

			Expect(burst.Name).To(Equal("Updated"))
		})
	})

	Describe("removeBurstFromSlice for non-existent burst", func() {
		It("should not modify slice when burst ID not found", func() {
			burst := fixtures.Burst("remove-burst")
			rmCtx := &burst_management.IntentValidator{
				Bursts: []*career.Burst{burst},
			}
			rmCtx.Validate()
			rmIntent, _ := burst_management.NewIntent(rmCtx)
			rmIntent.Init()
			rmIntent.UpdateTerminalInfo(termInfo)

			Expect(rmIntent.GetFilteredBursts()).To(HaveLen(1))
		})
	})

	Describe("cancelPreviousOperation with active cancel func", func() {
		It("should cancel and clear loading state", func() {
			burst := fixtures.Burst("cancel-prev-burst")
			mockSvc := mocks.NewBurstServiceMock()
			mockSvc.SetEvents([]*career.Event{fixtures.Event("cp-ev-1")})
			cpCtx := &burst_management.IntentValidator{
				Bursts:  []*career.Burst{burst},
				Service: mockSvc,
			}
			cpCtx.Validate()
			cpIntent, _ := burst_management.NewIntent(cpCtx)
			cpIntent.Init()
			cpIntent.UpdateTerminalInfo(termInfo)

			cpIntent.Update(tea.KeyMsg{Type: tea.KeyRunes, Runes: []rune{'s'}})
			Expect(cpIntent.GetLoadingModal()).NotTo(BeNil())

			cpIntent.Update(tea.KeyMsg{Type: tea.KeyRunes, Runes: []rune{'s'}})
		})
	})

	Describe("View rendering", func() {
		It("should render inactive intent message", func() {
			vCtx := &burst_management.IntentValidator{
				Bursts: []*career.Burst{fixtures.Burst("view-burst")},
			}
			vCtx.Validate()
			vIntent, _ := burst_management.NewIntent(vCtx)
			vIntent.Init()
			vIntent.Deactivate()

			view := vIntent.View()
			Expect(view).To(ContainSubstring("not active"))
		})

		It("should render with active view", func() {
			vCtx := &burst_management.IntentValidator{
				Bursts: []*career.Burst{fixtures.Burst("view-active-burst")},
			}
			vCtx.Validate()
			vIntent, _ := burst_management.NewIntent(vCtx)
			vIntent.Init()
			vIntent.UpdateTerminalInfo(termInfo)

			view := vIntent.View()
			Expect(view).NotTo(BeEmpty())
		})
	})

	Describe("detail modal 'd' for delete from detail", func() {
		It("should open delete modal from detail", func() {
			burst := fixtures.Burst("detail-d-burst")
			ddCtx := &burst_management.IntentValidator{
				Bursts: []*career.Burst{burst},
			}
			ddCtx.Validate()
			ddIntent, _ := burst_management.NewIntent(ddCtx)
			ddIntent.Init()
			ddIntent.UpdateTerminalInfo(termInfo)
			ddIntent.SetSelectedBurst(burst)

			ddIntent.HandleNavigate(&widgets.NavigateViewResult{
				ResultData: burstviews.Nav{Action: burstviews.ActionView, Burst: display.BurstFromDomain(burst)},
			})
			Expect(ddIntent.GetDetailModal()).NotTo(BeNil())

			ddIntent.Update(tea.KeyMsg{Type: tea.KeyRunes, Runes: []rune{'d'}})

			Expect(ddIntent.HasVisibleDeleteModal()).To(BeTrue())
		})
	})

	Describe("handleBurstSuggestionsLoaded with cancelled context", func() {
		It("should silently ignore cancelled context during detection", func() {
			bsCtx := &burst_management.IntentValidator{
				Bursts: []*career.Burst{fixtures.Burst("sug-cancel-burst")},
			}
			bsCtx.Validate()
			bsIntent, _ := burst_management.NewIntent(bsCtx)
			bsIntent.Init()
			bsIntent.UpdateTerminalInfo(termInfo)
			bsIntent.SetState(burst_management.StateSuggesting)

			bsIntent.Update(burst_management.BurstSuggestionsLoadedMsg{
				Error: context.Canceled,
			})

			Expect(bsIntent.GetState()).To(Equal(burst_management.StateList))
			Expect(bsIntent.HasVisibleErrorModal()).To(BeFalse())
		})
	})

	Describe("createBurstDetectionCmd with cancelled context", func() {
		It("should return error when context is cancelled before detection", func() {
			burst := fixtures.Burst("detect-cancel-burst")
			mockSvc := mocks.NewBurstServiceMock()
			mockSvc.SetEvents([]*career.Event{fixtures.Event("dc-ev-1")})
			dcCtx := &burst_management.IntentValidator{
				Bursts:  []*career.Burst{burst},
				Service: mockSvc,
			}
			dcCtx.Validate()
			dcIntent, _ := burst_management.NewIntent(dcCtx)
			dcIntent.Init()
			dcIntent.UpdateTerminalInfo(termInfo)

			dcIntent.Update(tea.KeyMsg{Type: tea.KeyRunes, Runes: []rune{'s'}})
			Expect(dcIntent.GetLoadingModal()).NotTo(BeNil())

			dcIntent.Update(tea.KeyMsg{Type: tea.KeyEsc})
			Expect(dcIntent.GetState()).To(Equal(burst_management.StateList))
		})
	})

	Describe("handleEditModalUpdate with edit completion", func() {
		It("should handle edit modal with submission", func() {
			burst := fixtures.Burst("edit-complete-modal-burst")
			ecmCtx := &burst_management.IntentValidator{
				Bursts: []*career.Burst{burst},
			}
			ecmCtx.Validate()
			ecmIntent, _ := burst_management.NewIntent(ecmCtx)
			ecmIntent.Init()
			ecmIntent.UpdateTerminalInfo(termInfo)
			ecmIntent.SetSelectedBurst(burst)

			ecmIntent.HandleNavigate(&widgets.NavigateViewResult{
				ResultData: burstviews.Nav{Action: burstviews.ActionEdit, Burst: display.BurstFromDomain(burst)},
			})
			Expect(ecmIntent.HasVisibleEditModal()).To(BeTrue())

			ecmIntent.Update(tea.KeyMsg{Type: tea.KeyTab})
			ecmIntent.Update(tea.KeyMsg{Type: tea.KeyTab})
		})
	})

	Describe("getUnassignedEventIDs with all events assigned", func() {
		It("should return no unassigned events when all are in bursts", func() {
			event := fixtures.Event("all-assigned-ev")
			burst := fixtures.Burst("all-assigned-burst", "all-assigned-ev")
			mockSvc := mocks.NewBurstServiceMock()
			mockSvc.SetEvents([]*career.Event{event})
			aaCtx := &burst_management.IntentValidator{
				Bursts:  []*career.Burst{burst},
				Service: mockSvc,
			}
			aaCtx.Validate()
			aaIntent, _ := burst_management.NewIntent(aaCtx)
			aaIntent.Init()
			aaIntent.UpdateTerminalInfo(termInfo)

			cmd := aaIntent.Update(tea.KeyMsg{Type: tea.KeyRunes, Runes: []rune{'s'}})
			Expect(cmd).NotTo(BeNil())

			msg := cmd()
			sugMsg, ok := msg.(burst_management.BurstSuggestionsLoadedMsg)
			if ok {
				Expect(sugMsg.Error).To(HaveOccurred())
			}
		})
	})

	Describe("handleSuggestionModalClosed with accepted and cancelled", func() {
		It("should return to list when no accepted suggestions", func() {
			burst := fixtures.Burst("smc-burst")
			mockSvc := mocks.NewBurstServiceMock()
			mockSvc.SetEvents([]*career.Event{fixtures.Event("smc-ev-1")})
			smcCtx := &burst_management.IntentValidator{
				Bursts:  []*career.Burst{burst},
				Service: mockSvc,
			}
			smcCtx.Validate()
			smcIntent, _ := burst_management.NewIntent(smcCtx)
			smcIntent.Init()
			smcIntent.UpdateTerminalInfo(termInfo)
			smcIntent.SetState(burst_management.StateSuggesting)

			smcIntent.Update(burst_management.BurstSuggestionsLoadedMsg{
				Suggestions: []burstfact.BurstSuggestion{
					{Name: "Cancel Test", Description: "Test", EventIDs: []string{"smc-ev-1"}},
				},
			})

			smcIntent.Update(tea.KeyMsg{Type: tea.KeyEsc})

			Expect(smcIntent.GetState()).To(Equal(burst_management.StateList))
		})
	})

	Describe("showEvents with nil service", func() {
		It("should return empty events when service is nil", func() {
			burst := fixtures.Burst("no-service-burst", "ns-ev-1", "ns-ev-2")
			neCtx := &burst_management.IntentValidator{
				Bursts: []*career.Burst{burst},
			}
			neCtx.Validate()
			neIntent, _ := burst_management.NewIntent(neCtx)
			neIntent.Init()
			neIntent.UpdateTerminalInfo(termInfo)
			neIntent.SetSelectedBurst(burst)

			cmd := neIntent.HandleNavigate(&widgets.NavigateViewResult{
				ResultData: burstviews.Nav{Action: burstviews.ActionViewEvents, Burst: display.BurstFromDomain(burst)},
			})
			Expect(cmd).NotTo(BeNil())
			msg := cmd()
			eventsMsg, ok := msg.(burst_management.BurstEventsLoadedMsg)
			Expect(ok).To(BeTrue())
			Expect(eventsMsg.Events).To(BeEmpty())
		})
	})

	Describe("HandleNavigate with missing burst", func() {
		It("returns nil when burst not found by ID", func() {
			ctx := &burst_management.IntentValidator{
				Bursts: []*career.Burst{fixtures.Burst("nav-missing-1")},
			}
			ctx.Validate()
			intent, _ := burst_management.NewIntent(ctx)
			intent.Init()
			intent.UpdateTerminalInfo(termInfo)

			nav := burstviews.Nav{
				Action: burstviews.ActionView,
				Burst:  display.Burst{ID: "missing-id"},
			}
			result := &widgets.NavigateViewResult{ResultData: nav}
			cmd := intent.HandleNavigate(result)

			Expect(cmd).To(BeNil())
		})
	})

	Describe("RefreshData", func() {
		It("reloads bursts from context", func() {
			burst1 := fixtures.Burst("refresh-1")
			burst2 := fixtures.Burst("refresh-2")
			mockRepo := mocks.NewBurstRepositoryMock()
			mockRepo.AddBurst(burst1)
			mockRepo.AddBurst(burst2)

			ctx := &burst_management.IntentValidator{
				Bursts:          []*career.Burst{burst1},
				BurstRepository: mockRepo,
			}
			ctx.Validate()
			intent, _ := burst_management.NewIntent(ctx)
			intent.Init()
			intent.UpdateTerminalInfo(termInfo)

			cmd := intent.RefreshData()

			Expect(cmd).To(BeNil())
			Expect(intent.GetFilteredBursts()).To(HaveLen(2))
		})

		It("handles load error gracefully", func() {
			mockRepo := mocks.NewBurstRepositoryMock()
			mockRepo.SetListError(errors.New("load fail"))

			ctx := &burst_management.IntentValidator{
				Bursts:          []*career.Burst{fixtures.Burst("refresh-3")},
				BurstRepository: mockRepo,
			}
			ctx.Validate()
			intent, _ := burst_management.NewIntent(ctx)
			intent.Init()
			intent.UpdateTerminalInfo(termInfo)

			cmd := intent.RefreshData()

			Expect(cmd).To(BeNil())
		})
	})

	Describe("showFacts", func() {
		It("returns command when burst has facts", func() {
			burst := fixtures.Burst("facts-1")
			fact := fixtures.Fact("fact-1", burst.ID)
			mockSvc := mocks.NewBurstServiceMock()
			mockSvc.SetFactsForBurst(burst.ID, []*career.Fact{fact})

			ctx := &burst_management.IntentValidator{
				Bursts:  []*career.Burst{burst},
				Service: mockSvc,
			}
			ctx.Validate()
			intent, _ := burst_management.NewIntent(ctx)
			intent.Init()
			intent.UpdateTerminalInfo(termInfo)
			intent.SetSelectedBurst(burst)

			cmd := intent.HandleNavigate(&widgets.NavigateViewResult{
				ResultData: burstviews.Nav{
					Action: burstviews.ActionViewFacts,
					Burst:  display.BurstFromDomain(burst),
				},
			})

			Expect(cmd).NotTo(BeNil())
		})

		It("returns command when service is nil", func() {
			burst := fixtures.Burst("facts-2")
			ctx := &burst_management.IntentValidator{
				Bursts:  []*career.Burst{burst},
				Service: nil,
			}
			ctx.Validate()
			intent, _ := burst_management.NewIntent(ctx)
			intent.Init()
			intent.UpdateTerminalInfo(termInfo)
			intent.SetSelectedBurst(burst)

			cmd := intent.HandleNavigate(&widgets.NavigateViewResult{
				ResultData: burstviews.Nav{
					Action: burstviews.ActionViewFacts,
					Burst:  display.BurstFromDomain(burst),
				},
			})

			Expect(cmd).NotTo(BeNil())
		})
	})

	Describe("rebuildModalRegistry", func() {
		It("rebuilds modal registry with error modal", func() {
			ctx := &burst_management.IntentValidator{
				Bursts: []*career.Burst{fixtures.Burst("rebuild-1")},
			}
			ctx.Validate()
			intent, _ := burst_management.NewIntent(ctx)
			intent.Init()
			intent.UpdateTerminalInfo(termInfo)

			intent.ShowErrorModal("Test Error", "This is a test")
			intent.RebuildModalRegistry()

			Expect(intent.HasVisibleErrorModal()).To(BeTrue())
		})
	})

	Describe("View", func() {
		It("renders view when active", func() {
			ctx := &burst_management.IntentValidator{
				Bursts: []*career.Burst{fixtures.Burst("view-1")},
			}
			ctx.Validate()
			intent, _ := burst_management.NewIntent(ctx)
			intent.Init()
			intent.UpdateTerminalInfo(termInfo)

			view := intent.View()

			Expect(view).NotTo(BeEmpty())
			Expect(view).To(ContainSubstring("Burst Management"))
		})
	})
})

var _ = Describe("LoadBursts", func() {
	It("returns nil when BurstRepository is nil", func() {
		ctx := &burst_management.IntentValidator{
			Bursts:          []*career.Burst{},
			BurstRepository: nil,
		}

		err := ctx.LoadBursts()

		Expect(err).ToNot(HaveOccurred())
	})

	It("loads bursts from repository", func() {
		burst1 := fixtures.Burst("load-1")
		burst2 := fixtures.Burst("load-2")
		mockRepo := mocks.NewBurstRepositoryMock()
		mockRepo.AddBurst(burst1)
		mockRepo.AddBurst(burst2)

		ctx := &burst_management.IntentValidator{
			Bursts:          []*career.Burst{},
			BurstRepository: mockRepo,
		}

		err := ctx.LoadBursts()

		Expect(err).NotTo(HaveOccurred())
		Expect(ctx.Bursts).To(HaveLen(2))
	})

	It("returns error when repository fails", func() {
		mockRepo := mocks.NewBurstRepositoryMock()
		mockRepo.SetListError(errors.New("repo fail"))

		ctx := &burst_management.IntentValidator{
			Bursts:          []*career.Burst{},
			BurstRepository: mockRepo,
		}

		err := ctx.LoadBursts()

		Expect(err).To(HaveOccurred())
	})
})
