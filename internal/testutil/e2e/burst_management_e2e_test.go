package e2e_test

import (
	"context"

	burstmgmt "github.com/baphled/kariya/internal/cli/intents/burst_management"
	"github.com/baphled/kariya/internal/domain/career"
	"github.com/baphled/kariya/internal/service/career/burstfact"
	"github.com/baphled/kariya/internal/service/career/skillinference"
	"github.com/baphled/kariya/internal/testutil/e2e"
	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
)

var _ = Describe("E2E Burst Management Workflow", func() {
	var env *e2e.TestEnv

	Describe("Empty Burst List Data State", func() {
		BeforeEach(func() {
			env = e2e.GetSharedEnv(GinkgoT())
		})

		AfterEach(func() {
			env.Cleanup()
		})

		It("should verify database has no bursts initially", func() {
			env.AssertBurstCount(0)
		})
	})

	Describe("Burst List with Persisted Data", func() {
		BeforeEach(func() {
			env = e2e.GetSharedEnv(GinkgoT())
			env.PopulateTestData(5, 2, 0) // 5 events, 2 bursts, 0 facts
		})

		AfterEach(func() {
			env.Cleanup()
		})

		It("should display bursts from database in list", func() {
			env.SelectIntentByName("burst_management")
			env.AssertBurstCount(2)
			env.AssertViewContainsAny("Authentication", "Mentoring", "Burst", "Name")
		})

		It("should show burst data from fixtures", func() {
			env.SelectIntentByName("burst_management")
			env.Confirm()
			env.AssertViewContainsAny("Authentication", "Mentoring", "Description")
		})
	})

	Describe("Session Persistence", func() {
		BeforeEach(func() {
			env = e2e.GetSharedEnv(GinkgoT())
		})

		AfterEach(func() {
			env.Cleanup()
		})

		It("should show bursts after restart", func() {
			env.PopulateTestData(5, 2, 0)
			env.AssertBurstCount(2)

			env.SimulateRestart()

			env.AssertBurstCount(2)
			env.SelectIntentByName("burst_management")
			env.AssertViewContainsAny("Authentication", "Mentoring", "Burst")
		})

		It("should maintain burst data integrity after restart", func() {
			env.PopulateTestData(5, 2, 0)
			burstsBefore := env.GetBursts()
			Expect(len(burstsBefore)).To(Equal(2))

			env.SimulateRestart()

			burstsAfter := env.GetBursts()
			Expect(len(burstsAfter)).To(Equal(2))
			Expect(burstsAfter[0].ID).To(Equal(burstsBefore[0].ID))
		})
	})

	Describe("Burst Edit Data Persistence", func() {
		BeforeEach(func() {
			env = e2e.GetSharedEnv(GinkgoT())
			env.PopulateTestData(5, 2, 0)
		})

		AfterEach(func() {
			env.Cleanup()
		})

		It("should persist burst name changes to database", func() {
			originalBursts := env.GetBursts()
			Expect(len(originalBursts)).To(BeNumerically(">=", 1))
			originalName := originalBursts[0].Name
			burstID := originalBursts[0].ID

			// Navigate to burst management and open detail modal.
			env.SelectIntentByName("burst_management")
			env.Confirm()
			// Press 'e' to open edit modal.
			env.PressKeyRune('e')
			env.AssertViewContainsAny("Edit Burst", "Burst Name", "Name")

			// Simulate form completion by sending EditBurstMsg directly.
			// This bypasses huh form internal state complexity.
			newName := "UPDATED_BURST_NAME_E2E_TEST"
			editMsg := burstmgmt.EditBurstMsg{
				BurstID:     burstID,
				Name:        newName,
				Description: originalBursts[0].Description,
			}
			env.SendMessage(editMsg)

			// Verify database was updated.
			updatedBursts := env.GetBursts()
			var foundBurst *career.Burst
			for _, b := range updatedBursts {
				if b.ID == burstID {
					foundBurst = b
					break
				}
			}

			Expect(foundBurst).NotTo(BeNil(), "Should find the updated burst in database")
			Expect(foundBurst.Name).To(Equal(newName), "Burst name should be updated in database")
			Expect(foundBurst.Name).NotTo(Equal(originalName), "Burst name should differ from original")
		})

		It("should persist burst description changes to database", func() {
			originalBursts := env.GetBursts()
			Expect(len(originalBursts)).To(BeNumerically(">=", 1))
			originalDesc := originalBursts[0].Description
			burstID := originalBursts[0].ID

			// Navigate to burst management and open detail modal.
			env.SelectIntentByName("burst_management")
			env.Confirm()
			// Press 'e' to open edit modal.
			env.PressKeyRune('e')

			// Simulate form completion by sending EditBurstMsg directly.
			newDesc := "UPDATED_DESCRIPTION_E2E_TEST"
			editMsg := burstmgmt.EditBurstMsg{
				BurstID:     burstID,
				Name:        originalBursts[0].Name,
				Description: newDesc,
			}
			env.SendMessage(editMsg)

			// Verify database was updated.
			updatedBursts := env.GetBursts()
			var foundBurst *career.Burst
			for _, b := range updatedBursts {
				if b.ID == burstID {
					foundBurst = b
					break
				}
			}

			Expect(foundBurst).NotTo(BeNil())
			Expect(foundBurst.Description).To(Equal(newDesc))
			Expect(foundBurst.Description).NotTo(Equal(originalDesc))
		})

		It("should NOT persist changes when edit is cancelled", func() {
			originalBursts := env.GetBursts()
			Expect(len(originalBursts)).To(BeNumerically(">=", 1))
			originalName := originalBursts[0].Name
			burstID := originalBursts[0].ID

			// Navigate to burst management and open detail modal.
			env.SelectIntentByName("burst_management")
			env.Confirm()
			// Press 'e' to open edit modal.
			env.PressKeyRune('e')
			env.AssertViewContainsAny("Edit Burst", "Burst Name", "Name")

			// Cancel by pressing Esc.
			env.Cancel()

			// Verify database was NOT changed.
			unchangedBursts := env.GetBursts()
			var foundBurst *career.Burst
			for _, b := range unchangedBursts {
				if b.ID == burstID {
					foundBurst = b
					break
				}
			}

			Expect(foundBurst).NotTo(BeNil())
			Expect(foundBurst.Name).To(Equal(originalName), "Burst name should NOT change when cancelled")
		})

		It("should persist changes and survive application restart", func() {
			originalBursts := env.GetBursts()
			Expect(len(originalBursts)).To(BeNumerically(">=", 1))
			burstID := originalBursts[0].ID

			// Navigate to burst management and open detail modal.
			env.SelectIntentByName("burst_management")
			env.Confirm()
			// Press 'e' to open edit modal.
			env.PressKeyRune('e')

			// Simulate form completion by sending EditBurstMsg directly.
			newName := "PERSISTED_ACROSS_RESTART"
			editMsg := burstmgmt.EditBurstMsg{
				BurstID:     burstID,
				Name:        newName,
				Description: originalBursts[0].Description,
			}
			env.SendMessage(editMsg)

			// Restart the application.
			env.SimulateRestart()

			// Verify changes persisted across restart.
			persistedBursts := env.GetBursts()
			var foundBurst *career.Burst
			for _, b := range persistedBursts {
				if b.ID == burstID {
					foundBurst = b
					break
				}
			}

			Expect(foundBurst).NotTo(BeNil(), "Burst should exist after restart")
			Expect(foundBurst.Name).To(Equal(newName), "Burst name should persist after restart")
		})
	})

	// BUG: After accepting a burst suggestion, the detail modal is shown but:
	// - CRUD keys (e, d, v, f) don't work
	// - Need to press Esc twice to close it
	// - Cannot return to main menu
	//
	// Expected flow:
	// 1. Accept suggestion ('a') -> Burst detail modal is shown
	// 2. From detail modal: 'f' (facts), 'v' (events), 'e' (edit), 'd' (delete) should work
	// 3. Esc from detail modal -> Return to burst list
	// 4. Esc from burst list -> Return to main menu
	// Tests for suggestion acceptance flow:
	// - Accept suggestion -> stay on suggestion review, see next suggestion
	// - Accept last suggestion -> return to burst list
	// - Burst list is updated with accepted bursts
	Describe("Burst Suggestion Acceptance Flow", func() {
		BeforeEach(func() {
			env = e2e.GetSharedEnv(GinkgoT())
			// Need events that can be grouped into suggestions.
			env.PopulateTestData(10, 0, 0) // 10 events, 0 bursts, 0 facts
		})

		AfterEach(func() {
			env.Cleanup()
		})

		It("should stay on suggestion review and show next suggestion after accepting", func() {
			// Navigate to burst management.
			env.SelectIntentByName("burst_management")

			events := env.GetEvents()
			Expect(len(events)).To(BeNumerically(">=", 4))

			// Press 's' and simulate TWO suggestions.
			env.PressKeyRune('s')
			suggestions := []burstfact.BurstSuggestion{
				{
					EventIDs:        []string{events[0].ID, events[1].ID},
					ConfidenceScore: 0.85,
					Name:            "First Burst",
					Description:     "First suggestion to accept",
				},
				{
					EventIDs:        []string{events[2].ID, events[3].ID},
					ConfidenceScore: 0.80,
					Name:            "Second Burst",
					Description:     "Second suggestion should appear after accepting first",
				},
			}
			env.SendMessage(burstmgmt.BurstSuggestionsLoadedMsg{
				Suggestions: suggestions,
			})

			// Verify we see both suggestions in the table view.
			view := env.GetView()
			Expect(view).To(ContainSubstring("First Burst"),
				"Should show first suggestion")
			Expect(view).To(ContainSubstring("Second Burst"),
				"Should show second suggestion")
			Expect(view).To(ContainSubstring("Suggestions: 2"),
				"Should show total suggestion count")

			// Accept the first suggestion - should stay on suggestion review and show remaining.
			env.PressKeyRune('a')

			// Verify we now see the second suggestion (1 remaining after accepting first).
			view = env.GetView()
			Expect(view).To(ContainSubstring("Second Burst"),
				"After accepting first, should still show second suggestion")
			Expect(view).To(ContainSubstring("Suggestions: 1"),
				"Should show suggestion count '1' after accepting first")

			// Should NOT be on burst list yet.
			Expect(view).NotTo(ContainSubstring("Burst List"),
				"Should still be on suggestion review, not burst list")
		})

		It("should return to burst list after accepting last suggestion", func() {
			// Navigate to burst management.
			env.SelectIntentByName("burst_management")

			events := env.GetEvents()
			Expect(len(events)).To(BeNumerically(">=", 2))

			// Press 's' and simulate suggestions.
			env.PressKeyRune('s')
			suggestion := burstfact.BurstSuggestion{
				EventIDs:        []string{events[0].ID, events[1].ID},
				ConfidenceScore: 0.85,
				Name:            "Accepted Burst",
				Description:     "Testing return to list after acceptance",
			}
			env.SendMessage(burstmgmt.BurstSuggestionsLoadedMsg{
				Suggestions: []burstfact.BurstSuggestion{suggestion},
			})

			// Accept the suggestion - should return to burst list (not detail modal).
			env.PressKeyRune('a')

			// Should be at burst list showing the newly created burst.
			view := env.GetView()
			Expect(view).To(ContainSubstring("Accepted Burst"),
				"After accepting last suggestion, should return to burst list showing accepted burst")

			// Should NOT be at main menu.
			Expect(env.IsInMenuState()).To(BeFalse(),
				"Should be in burst management, not main menu")
		})

		It("should open edit modal with 'e' key on burst list after accepting", func() {
			// Navigate to burst management.
			env.SelectIntentByName("burst_management")

			events := env.GetEvents()
			Expect(len(events)).To(BeNumerically(">=", 2))

			// Press 's' and simulate suggestions.
			env.PressKeyRune('s')
			suggestion := burstfact.BurstSuggestion{
				EventIDs:        []string{events[0].ID, events[1].ID},
				ConfidenceScore: 0.85,
				Name:            "Edit Test Burst",
				Description:     "Testing edit from list after acceptance",
			}
			env.SendMessage(burstmgmt.BurstSuggestionsLoadedMsg{
				Suggestions: []burstfact.BurstSuggestion{suggestion},
			})

			// Accept the suggestion - returns to burst list.
			env.PressKeyRune('a')

			// Press 'e' to edit selected burst in list.
			env.PressKeyRune('e')

			view := env.GetView()
			// Edit modal must show "Edit Burst" title.
			Expect(view).To(ContainSubstring("Edit Burst"),
				"'e' key on burst list should open edit modal showing 'Edit Burst' title")
		})

		It("should return to main menu with single Esc from burst list after accepting", func() {
			// Navigate to burst management.
			env.SelectIntentByName("burst_management")

			events := env.GetEvents()
			Expect(len(events)).To(BeNumerically(">=", 2))

			// Press 's' and simulate suggestions.
			env.PressKeyRune('s')
			suggestion := burstfact.BurstSuggestion{
				EventIDs:        []string{events[0].ID, events[1].ID},
				ConfidenceScore: 0.85,
				Name:            "Menu Test Burst",
				Description:     "Testing return to main menu",
			}
			env.SendMessage(burstmgmt.BurstSuggestionsLoadedMsg{
				Suggestions: []burstfact.BurstSuggestion{suggestion},
			})

			// Accept the suggestion - returns to burst list.
			env.PressKeyRune('a')

			// Verify we're at burst list with the accepted burst.
			view := env.GetView()
			Expect(view).To(ContainSubstring("Menu Test Burst"),
				"Should be at burst list showing accepted burst")

			// Single Esc should return to main menu.
			env.Cancel()

			// Should be at main menu now.
			Expect(env.IsInMenuState()).To(BeTrue(),
				"Single Esc from burst list should return to main menu")
		})
	})

	// Tests for suggestion rejection flow:
	// - Reject suggestion -> return to burst list
	// - From burst list, Esc -> return to main menu
	// - Reject all suggestions -> return to list with no bursts
	Describe("Burst Suggestion Rejection Flow", func() {
		BeforeEach(func() {
			env = e2e.GetSharedEnv(GinkgoT())
			// Need events that can be grouped into suggestions.
			env.PopulateTestData(10, 0, 0) // 10 events, 0 bursts, 0 facts
		})

		AfterEach(func() {
			env.Cleanup()
		})

		It("should return to burst list after rejecting last suggestion", func() {
			// Navigate to burst management.
			env.SelectIntentByName("burst_management")

			events := env.GetEvents()
			Expect(len(events)).To(BeNumerically(">=", 2))

			// Press 's' and simulate a single suggestion.
			env.PressKeyRune('s')
			suggestion := burstfact.BurstSuggestion{
				EventIDs:        []string{events[0].ID, events[1].ID},
				ConfidenceScore: 0.85,
				Name:            "Rejected Burst",
				Description:     "Testing return to list after rejection",
			}
			env.SendMessage(burstmgmt.BurstSuggestionsLoadedMsg{
				Suggestions: []burstfact.BurstSuggestion{suggestion},
			})

			// Verify we're on suggestion review.
			view := env.GetView()
			Expect(view).To(ContainSubstring("Rejected Burst"),
				"Should show the suggestion to review")

			// Reject the suggestion - should return to burst list.
			env.PressKeyRune('r')

			// Should be at burst list (empty, no bursts created).
			view = env.GetView()
			Expect(view).NotTo(ContainSubstring("Rejected Burst"),
				"After rejecting, should not show rejected burst")

			// Should NOT be at main menu - should be at burst list.
			Expect(env.IsInMenuState()).To(BeFalse(),
				"Should be in burst management, not main menu")
		})

		It("should return to main menu with Esc from burst list after rejecting", func() {
			// Navigate to burst management.
			env.SelectIntentByName("burst_management")

			events := env.GetEvents()
			Expect(len(events)).To(BeNumerically(">=", 2))

			// Press 's' and simulate a single suggestion.
			env.PressKeyRune('s')
			suggestion := burstfact.BurstSuggestion{
				EventIDs:        []string{events[0].ID, events[1].ID},
				ConfidenceScore: 0.85,
				Name:            "Rejected Menu Test",
				Description:     "Testing return to main menu after rejection",
			}
			env.SendMessage(burstmgmt.BurstSuggestionsLoadedMsg{
				Suggestions: []burstfact.BurstSuggestion{suggestion},
			})

			// Reject the suggestion - returns to burst list.
			env.PressKeyRune('r')

			// Verify we're at burst list (not main menu yet).
			Expect(env.IsInMenuState()).To(BeFalse(),
				"Should be at burst list, not main menu")

			// Single Esc should return to main menu.
			env.Cancel()

			// Should be at main menu now.
			Expect(env.IsInMenuState()).To(BeTrue(),
				"Single Esc from burst list after rejection should return to main menu")
		})

		It("should return to main menu after rejecting all multiple suggestions", func() {
			// Navigate to burst management.
			env.SelectIntentByName("burst_management")

			events := env.GetEvents()
			Expect(len(events)).To(BeNumerically(">=", 4))

			// Press 's' and simulate TWO suggestions.
			env.PressKeyRune('s')
			suggestions := []burstfact.BurstSuggestion{
				{
					EventIDs:        []string{events[0].ID, events[1].ID},
					ConfidenceScore: 0.85,
					Name:            "First Rejection",
					Description:     "First suggestion to reject",
				},
				{
					EventIDs:        []string{events[2].ID, events[3].ID},
					ConfidenceScore: 0.80,
					Name:            "Second Rejection",
					Description:     "Second suggestion to reject",
				},
			}
			env.SendMessage(burstmgmt.BurstSuggestionsLoadedMsg{
				Suggestions: suggestions,
			})

			// Verify we see the first suggestion.
			view := env.GetView()
			Expect(view).To(ContainSubstring("First Rejection"),
				"Should show first suggestion")

			// Reject the first suggestion - should show second.
			env.PressKeyRune('r')

			// Verify we now see the second suggestion.
			view = env.GetView()
			Expect(view).To(ContainSubstring("Second Rejection"),
				"After rejecting first, should show second suggestion")

			// Reject the second suggestion - should return to burst list.
			env.PressKeyRune('r')

			// Should NOT be at main menu - should be at burst list.
			Expect(env.IsInMenuState()).To(BeFalse(),
				"Should be at burst list, not main menu")

			// Single Esc should return to main menu.
			env.Cancel()

			// Should be at main menu now.
			Expect(env.IsInMenuState()).To(BeTrue(),
				"Single Esc from burst list after rejecting all should return to main menu")
		})

		It("should have no bursts after rejecting all suggestions", func() {
			// Navigate to burst management.
			env.SelectIntentByName("burst_management")

			events := env.GetEvents()
			Expect(len(events)).To(BeNumerically(">=", 2))

			// Press 's' and simulate suggestions.
			env.PressKeyRune('s')
			suggestions := []burstfact.BurstSuggestion{
				{
					EventIDs:        []string{events[0].ID, events[1].ID},
					ConfidenceScore: 0.85,
					Name:            "Rejected Burst",
					Description:     "This should not be saved",
				},
			}
			env.SendMessage(burstmgmt.BurstSuggestionsLoadedMsg{
				Suggestions: suggestions,
			})

			// Reject the suggestion.
			env.PressKeyRune('r')

			// Verify no bursts were created.
			env.AssertBurstCount(0)
		})
	})
})

var _ = Describe("E2E Skill Suggestion Acceptance from Burst (SQL-backed)", func() {
	var env *e2e.TestEnv

	BeforeEach(func() {
		env = e2e.GetSharedEnv(GinkgoT())
		env.PopulateTestData(10, 0, 0)
	})

	AfterEach(func() {
		env.Cleanup()
	})

	It("should persist accepted skill suggestions to the SQL database", func() {
		env.SelectIntentByName("burst_management")

		events := env.GetEvents()
		Expect(len(events)).To(BeNumerically(">=", 2))

		env.PressKeyRune('s')
		burstSuggestions := []burstfact.BurstSuggestion{
			{
				EventIDs:        []string{events[0].ID, events[1].ID},
				ConfidenceScore: 0.90,
				Name:            "API Development",
				Description:     "Built REST API",
			},
		}
		env.SendMessage(burstmgmt.BurstSuggestionsLoadedMsg{
			Suggestions: burstSuggestions,
		})

		env.PressKeyRune('a')

		env.AssertBurstCount(1)

		bursts := env.GetBursts()
		acceptedBurst := bursts[0]
		Expect(acceptedBurst.Name).To(Equal("API Development"))

		env.Confirm()

		skillSuggestions := []skillinference.SkillSuggestion{
			{
				Name:       "Go",
				Category:   "Backend",
				Confidence: 0.95,
				EventIDs:   []string{events[0].ID, events[1].ID},
				Contexts:   []string{"Built REST API with Go"},
			},
			{
				Name:       "PostgreSQL",
				Category:   "Database",
				Confidence: 0.85,
				EventIDs:   []string{events[0].ID},
				Contexts:   []string{"Designed PostgreSQL schema"},
			},
		}
		env.SendMessage(burstmgmt.SkillSuggestionsLoadedMsg{
			Suggestions: skillSuggestions,
		})

		env.PressKeyRune('a')
		env.PressKeyRune('a')

		view := env.GetView()
		Expect(view).To(ContainSubstring("Skills Created"),
			"should show success modal after skill creation completes")
		Expect(view).To(ContainSubstring("Successfully created 2 skill"),
			"success modal should indicate 2 skills were created")

		repoSkills, err := env.SkillRepo.List(context.Background(), nil)
		Expect(err).NotTo(HaveOccurred())
		Expect(repoSkills).To(HaveLen(2),
			"SQL database should have 2 skills after accepting all suggestions")

		skillNames := make(map[string]bool)
		for _, s := range repoSkills {
			skillNames[s.Name] = true
		}
		Expect(skillNames).To(HaveKey("Go"))
		Expect(skillNames).To(HaveKey("PostgreSQL"))
	})
})
