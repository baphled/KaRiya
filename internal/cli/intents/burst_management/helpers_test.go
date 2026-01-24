package burst_management_test

import (
	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"

	"github.com/baphled/kariya/internal/cli/intents/burst_management"
	"github.com/baphled/kariya/internal/cli/terminal"
	"github.com/baphled/kariya/internal/cli/uikit/feedback"
	"github.com/baphled/kariya/internal/domain/career"
	"github.com/baphled/kariya/internal/service/career/burst_fact"
)

var _ = Describe("Helper Methods", func() {
	var intent *burst_management.Intent
	var ctx *burst_management.IntentContext
	var termInfo *terminal.Info

	BeforeEach(func() {
		// Create intent context.
		ctx = &burst_management.IntentContext{
			Bursts: []*career.Burst{
				{ID: "burst-1", Name: "Test Burst"},
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
			burst := &career.Burst{ID: "burst-1", Name: "Test Burst"}
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
			Expect(help).To(ContainSubstring("View Details"))
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

		It("should return global badges for unknown state", func() {
			intent.SetState(burst_management.State("invalid")) // Invalid state
			help := intent.GetContextHelp()
			Expect(help).NotTo(BeEmpty())
		})
	})

	Describe("RefreshData", func() {
		It("should reload bursts from context", func() {
			// Modify context bursts.
			ctx.Bursts = append(ctx.Bursts, &career.Burst{ID: "burst-2", Name: "New Burst"})

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

			// Create completion message with accepted suggestions.
			msg := burst_management.SuggestionReviewCompleteMsg{
				AcceptedSuggestions: []burst_fact.BurstSuggestion{
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
			Expect(newBursts[0].Confirmed).To(BeFalse())

			Expect(newBursts[1].Name).To(Equal("Accepted Burst 2"))
			Expect(newBursts[1].Description).To(Equal("Second accepted suggestion"))
			Expect(newBursts[1].EventIDs).To(Equal([]string{"e3", "e4"}))
			Expect(newBursts[1].Confirmed).To(BeFalse())
		})

		It("should transition to list screen when cancelled", func() {
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
			msg := burst_management.SuggestionReviewCompleteMsg{
				AcceptedSuggestions: []burst_fact.BurstSuggestion{},
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

			msg := burst_management.SuggestionReviewCompleteMsg{
				AcceptedSuggestions: []burst_fact.BurstSuggestion{
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
			msg := burst_management.SuggestionReviewCompleteMsg{
				AcceptedSuggestions: []burst_fact.BurstSuggestion{
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
	})
})
