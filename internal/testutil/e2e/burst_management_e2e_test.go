package e2e_test

import (
	burstmgmt "github.com/baphled/kariya/internal/cli/intents/burst_management"
	"github.com/baphled/kariya/internal/domain/career"
	"github.com/baphled/kariya/internal/service/career/burst_fact"
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
			suggestions := []burst_fact.BurstSuggestion{
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

			// Verify we see the first suggestion (1 of 2).
			view := env.GetView()
			Expect(view).To(ContainSubstring("First Burst"),
				"Should show first suggestion")
			Expect(view).To(ContainSubstring("1 of 2"),
				"Should show suggestion count '1 of 2'")

			// Accept the first suggestion - should stay on suggestion review and show second.
			env.PressKeyRune('a')

			// Verify we now see the second suggestion (1 of 1 since first was removed).
			view = env.GetView()
			Expect(view).To(ContainSubstring("Second Burst"),
				"After accepting first, should show second suggestion")
			Expect(view).To(ContainSubstring("1 of 1"),
				"Should show suggestion count '1 of 1' after accepting first")

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
			suggestion := burst_fact.BurstSuggestion{
				EventIDs:        []string{events[0].ID, events[1].ID},
				ConfidenceScore: 0.85,
				Name:            "Accepted Burst",
				Description:     "Testing return to list after acceptance",
			}
			env.SendMessage(burstmgmt.BurstSuggestionsLoadedMsg{
				Suggestions: []burst_fact.BurstSuggestion{suggestion},
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
			suggestion := burst_fact.BurstSuggestion{
				EventIDs:        []string{events[0].ID, events[1].ID},
				ConfidenceScore: 0.85,
				Name:            "Edit Test Burst",
				Description:     "Testing edit from list after acceptance",
			}
			env.SendMessage(burstmgmt.BurstSuggestionsLoadedMsg{
				Suggestions: []burst_fact.BurstSuggestion{suggestion},
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
			suggestion := burst_fact.BurstSuggestion{
				EventIDs:        []string{events[0].ID, events[1].ID},
				ConfidenceScore: 0.85,
				Name:            "Menu Test Burst",
				Description:     "Testing return to main menu",
			}
			env.SendMessage(burstmgmt.BurstSuggestionsLoadedMsg{
				Suggestions: []burst_fact.BurstSuggestion{suggestion},
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
})
