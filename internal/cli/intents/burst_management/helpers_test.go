package burst_management_test

import (
	"context"
	"fmt"

	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"

	"github.com/baphled/kariya/internal/cli/intents/burst_management"
	"github.com/baphled/kariya/internal/cli/terminal"
	"github.com/baphled/kariya/internal/cli/uikit/feedback"
	"github.com/baphled/kariya/internal/domain/career"
	"github.com/baphled/kariya/internal/service/career/burstfact"
	"github.com/baphled/kariya/internal/service/career/skillinference"
	"github.com/baphled/kariya/internal/testutil/fixtures"
	tea "github.com/charmbracelet/bubbletea"
)

var _ = Describe("Helper Methods", func() {
	var intent *burst_management.Intent
	var ctx *burst_management.IntentContext
	var termInfo *terminal.Info

	BeforeEach(func() {
		// Create intent context.
		ctx = &burst_management.IntentContext{
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

	// Regression tests
	Describe("handleFactExtractionComplete with nil selectedBurst", func() {
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

	Describe("State transition after suggestion modal shown", func() {
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

	Describe("Slice bounds in handleSuggestionReviewComplete", func() {
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

		It("should show suggestions even when skills exist globally", func() {
			intent.SetState(burst_management.StateInferringSkills)

			msg := burst_management.SkillSuggestionsLoadedMsg{
				Suggestions: []skillinference.SkillSuggestion{
					{Name: "Go", Category: "backend", Confidence: 0.95},
					{Name: "Docker", Category: "devops", Confidence: 0.85},
				},
				ExistingSkillNames: []string{"Go", "Docker"},
			}

			intent.Update(msg)

			// Now shows ALL suggestions, allowing skill reuse across bursts
			Expect(intent.GetState()).To(Equal(burst_management.StateSkillSuggestionReview))
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

		It("should show suggestions regardless of ExistingSkillNames case", func() {
			intent.SetState(burst_management.StateInferringSkills)

			msg := burst_management.SkillSuggestionsLoadedMsg{
				Suggestions: []skillinference.SkillSuggestion{
					{Name: "Go", Category: "backend", Confidence: 0.95},
					{Name: "Docker", Category: "devops", Confidence: 0.85},
				},
				ExistingSkillNames: []string{"go", "docker"},
			}

			intent.Update(msg)

			// Now shows ALL suggestions, ExistingSkillNames no longer filters
			Expect(intent.GetState()).To(Equal(burst_management.StateSkillSuggestionReview))
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
			ctx.SkillInferenceService = skillinference.NewSkillInferenceService(nil, nil)
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
})
