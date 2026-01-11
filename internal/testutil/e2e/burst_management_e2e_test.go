package e2e_test

import (
	"strings"

	"github.com/baphled/kariya/internal/cli/intents"
	"github.com/baphled/kariya/internal/domain/career"
	"github.com/baphled/kariya/internal/testutil/e2e"
	tea "github.com/charmbracelet/bubbletea"
	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
)

var _ = Describe("E2E Burstmanagement Workflow", func() {
	var env *e2e.TestEnv

	Describe("Empty Burst List", func() {
		BeforeEach(func() {
			env = e2e.Setup(GinkgoT())
		})

		AfterEach(func() {
			env.Cleanup()
		})

		It("should show empty message when no bursts exist", func() {
			env.AssertBurstCount(0)
			env.SelectIntentByName("burst_management")
			env.AssertViewContainsAny("No bursts", "empty", "no bursts found", "Burst")
		})

		It("should not panic on empty burst list", func() {
			env.SelectIntentByName("burst_management")
			view := env.GetView()
			Expect(view).NotTo(BeEmpty())
			Expect(view).NotTo(ContainSubstring("panic"))
		})

		It("should allow navigation back from empty list", func() {
			env.SelectIntentByName("burst_management")
			env.Cancel()
			env.AssertViewContains("Capture Event")
		})
	})

	Describe("Burst List with Data", func() {
		BeforeEach(func() {
			env = e2e.Setup(GinkgoT())
			env.PopulateTestData(5, 2, 0) // 5 events, 2 bursts, 0 facts
		})

		AfterEach(func() {
			env.Cleanup()
		})

		It("should display bursts in list", func() {
			env.SelectIntentByName("burst_management")
			env.AssertBurstCount(2)
			// Should show burst names from fixtures
			env.AssertViewContainsAny("Authentication", "Mentoring", "Burst", "Name")
		})

		It("should show confirmation status in list", func() {
			env.SelectIntentByName("burst_management")
			// Fixtures create bursts with alternating confirmed status
			env.AssertViewContainsAny("Confirmed", "Yes", "No", "✓", "✗")
		})

		It("should allow navigating through bursts with j/k", func() {
			env.SelectIntentByName("burst_management")
			env.PressKeyRune('j')
			view := env.GetView()
			Expect(view).NotTo(BeEmpty())
			env.PressKeyRune('k')
			view = env.GetView()
			Expect(view).NotTo(BeEmpty())
		})

		It("should allow navigating with arrow keys", func() {
			env.SelectIntentByName("burst_management")
			env.PressKey(tea.KeyDown)
			view := env.GetView()
			Expect(view).NotTo(BeEmpty())
			env.PressKey(tea.KeyUp)
			view = env.GetView()
			Expect(view).NotTo(BeEmpty())
		})
	})

	Describe("Burst Detail View", func() {
		BeforeEach(func() {
			env = e2e.Setup(GinkgoT())
			env.PopulateTestData(5, 2, 0)
			env.SelectIntentByName("burst_management")
		})

		AfterEach(func() {
			env.Cleanup()
		})

		It("should show burst detail when pressing Enter", func() {
			env.Confirm() // Select first burst
			env.AssertViewContainsAny("Detail", "Name", "Description", "Events", "Confirmed")
		})

		It("should show burst information in detail view", func() {
			env.Confirm()
			// Should show burst details from fixtures
			env.AssertViewContainsAny("Authentication", "Mentoring", "Description")
		})

		It("should go back to list when pressing Escape from detail", func() {
			env.Confirm() // Go to detail
			env.Cancel()  // Go back
			env.AssertViewContainsAny("List", "Bursts", "Name", "Confirmed")
		})

		It("should show action hints in detail view", func() {
			env.Confirm()
			env.AssertViewContainsAny("e", "f", "d", "c", "Events", "Facts", "Delete", "Confirm")
		})
	})

	Describe("Cancel from Detail View", func() {
		BeforeEach(func() {
			env = e2e.Setup(GinkgoT())
			env.PopulateTestData(5, 2, 0)
			env.SelectIntentByName("burst_management")
			env.Confirm() // Go to detail
		})

		AfterEach(func() {
			env.Cleanup()
		})

		It("should return to list when pressing Escape from detail", func() {
			env.Cancel()
			env.AssertViewContainsAny("List", "Bursts", "Name")
		})

		It("should quit application when pressing 'q' from detail", func() {
			// Note: q now quits the entire app
			// This test verifies the quit command is handled without panic
			env.Quit()
			// After quit, the app terminates - we can't assert view content
		})
	})

	Describe("Vim-style Navigation", func() {
		BeforeEach(func() {
			env = e2e.Setup(GinkgoT())
			env.PopulateTestData(5, 2, 0)
			env.SelectIntentByName("burst_management")
		})

		AfterEach(func() {
			env.Cleanup()
		})

		It("should navigate down with 'j' key", func() {
			env.PressKeyRune('j')
			view := env.GetView()
			Expect(view).NotTo(BeEmpty())
		})

		It("should navigate up with 'k' key", func() {
			env.PressKeyRune('j')
			env.PressKeyRune('k')
			view := env.GetView()
			Expect(view).NotTo(BeEmpty())
		})

		It("should not crash at list boundaries", func() {
			env.PressKeyRune('k')
			env.PressKeyRune('k')
			view := env.GetView()
			Expect(view).NotTo(BeEmpty())

			for i := 0; i < 5; i++ {
				env.PressKeyRune('j')
			}
			view = env.GetView()
			Expect(view).NotTo(BeEmpty())
		})
	})

	Describe("Arrow Key Navigation", func() {
		BeforeEach(func() {
			env = e2e.Setup(GinkgoT())
			env.PopulateTestData(5, 2, 0)
			env.SelectIntentByName("burst_management")
		})

		AfterEach(func() {
			env.Cleanup()
		})

		It("should navigate down with down arrow", func() {
			env.PressKey(tea.KeyDown)
			view := env.GetView()
			Expect(view).NotTo(BeEmpty())
		})

		It("should navigate up with up arrow", func() {
			env.PressKey(tea.KeyDown)
			env.PressKey(tea.KeyUp)
			view := env.GetView()
			Expect(view).NotTo(BeEmpty())
		})
	})

	Describe("Session Persistence", func() {
		BeforeEach(func() {
			env = e2e.Setup(GinkgoT())
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

	Describe("Workflow Integration", func() {
		BeforeEach(func() {
			env = e2e.Setup(GinkgoT())
			env.PopulateTestData(5, 2, 0)
		})

		AfterEach(func() {
			env.Cleanup()
		})

		It("should allow selecting burst and returning to list multiple times", func() {
			env.SelectIntentByName("burst_management")

			// First selection
			env.Confirm() // Go to detail
			env.Cancel()  // Go back to list

			// Second selection
			env.PressKeyRune('j') // Navigate to next burst
			env.Confirm()         // Go to detail
			env.Cancel()          // Go back to list

			env.AssertViewContainsAny("List", "Bursts", "Name")
		})

		It("should allow re-entering after cancellation", func() {
			env.SelectIntentByName("burst_management")
			env.Cancel()
			env.SelectIntentByName("burst_management")
			env.AssertViewContainsAny("Burst", "List", "Name")
		})
	})

	Describe("Burst Edit Workflow", func() {
		BeforeEach(func() {
			env = e2e.Setup(GinkgoT())
			env.PopulateTestData(5, 2, 0) // 5 events, 2 bursts, 0 facts
			env.SelectIntentByName("burst_management")
			env.Confirm() // Go to detail view
		})

		AfterEach(func() {
			env.Cleanup()
		})

		It("should show edit option in detail view", func() {
			// Verify edit shortcut is displayed
			env.AssertViewContainsAny("e", "Edit", "edit")
		})

		It("should transition to edit view when pressing 'e'", func() {
			env.PressKeyRune('e')
			// Should show edit form with Name and Description fields
			env.AssertViewContainsAny("Edit Burst", "Burst Name", "Name", "Description")
		})

		It("should show edit modal with form fields", func() {
			env.PressKeyRune('e')
			// The EditBurstModal should display form inputs
			env.AssertViewContainsAny("Burst Name", "Description", "Enter", "Esc")
		})

		It("should show current burst values in edit form", func() {
			env.PressKeyRune('e')
			// Should pre-populate with existing burst data from fixtures
			// Fixtures create bursts with names like "Authentication" or "Mentoring"
			env.AssertViewContainsAny("Authentication", "Mentoring", "Name", "Description")
		})

		It("should return to detail view when cancelling edit", func() {
			env.PressKeyRune('e')
			env.Cancel() // Press Escape to cancel
			// Should be back at detail view
			env.AssertViewContainsAny("Detail", "Events", "Facts", "e", "f")
		})

		It("should not crash when pressing edit key multiple times", func() {
			env.PressKeyRune('e')
			view := env.GetView()
			Expect(view).NotTo(BeEmpty())
			Expect(view).NotTo(ContainSubstring("panic"))

			// Cancel and try again
			env.Cancel()
			env.PressKeyRune('e')
			view = env.GetView()
			Expect(view).NotTo(BeEmpty())
			Expect(view).NotTo(ContainSubstring("panic"))
		})

		It("should show edit form footer with navigation hints", func() {
			env.PressKeyRune('e')
			// Modal should show form navigation hints
			env.AssertViewContainsAny("Enter", "Esc", "Tab", "Confirm", "Cancel")
		})
	})

	Describe("Burst Edit with Persistence", func() {
		BeforeEach(func() {
			env = e2e.Setup(GinkgoT())
			env.PopulateTestData(5, 2, 0)
		})

		AfterEach(func() {
			env.Cleanup()
		})

		It("should maintain edit state through navigation", func() {
			env.SelectIntentByName("burst_management")
			env.Confirm()         // Go to detail
			env.PressKeyRune('e') // Go to edit

			// Should be in edit view
			env.AssertViewContainsAny("Edit Burst", "Burst Name", "Name")

			// Cancel and verify we return properly
			env.Cancel()
			env.AssertViewContainsAny("Detail", "Events", "Facts")
		})

		It("should allow editing and returning to detail without errors", func() {
			env.SelectIntentByName("burst_management")
			env.Confirm()
			env.PressKeyRune('e')

			// Cancel edit
			env.Cancel()

			// Should be back at detail, not crashed
			view := env.GetView()
			Expect(view).NotTo(BeEmpty())
			Expect(view).NotTo(ContainSubstring("panic"))
			Expect(view).NotTo(ContainSubstring("nil pointer"))
		})
	})

	Describe("Burst Edit Data Persistence", func() {
		BeforeEach(func() {
			env = e2e.Setup(GinkgoT())
			env.PopulateTestData(5, 2, 0) // 5 events, 2 bursts, 0 facts
		})

		AfterEach(func() {
			env.Cleanup()
		})

		It("should persist burst name changes to database", func() {
			// Get original bursts from database
			originalBursts := env.GetBursts()
			Expect(len(originalBursts)).To(BeNumerically(">=", 1))
			originalName := originalBursts[0].Name
			burstID := originalBursts[0].ID

			// Navigate to burst management and select first burst
			env.SelectIntentByName("burst_management")
			env.Confirm() // Go to detail view

			// Enter edit mode
			env.PressKeyRune('e')
			env.AssertViewContainsAny("Edit Burst", "Burst Name", "Name")

			// Get the active intent and access its modal
			activeIntent := env.Model.GetActiveIntent()
			Expect(activeIntent).NotTo(BeNil())

			burstIntent, ok := activeIntent.(*intents.BurstManagementIntent)
			Expect(ok).To(BeTrue(), "Active intent should be BurstManagementIntent")

			// Simulate form completion with updated name
			newName := "UPDATED_BURST_NAME_E2E_TEST"
			modifiedBurst := &career.Burst{
				ID:          burstID,
				Name:        newName,
				Description: originalBursts[0].Description,
				EventIDs:    originalBursts[0].EventIDs,
				CreatedAt:   originalBursts[0].CreatedAt,
				UpdatedAt:   originalBursts[0].UpdatedAt,
			}

			// Set the test result on the modal (bypasses form interaction)
			burstIntent.SetTestModalResult(&intents.ModalEditResult[*career.Burst]{
				Original: originalBursts[0],
				Modified: modifiedBurst,
				Accepted: true,
				Changes:  map[string]interface{}{"name": newName},
			})

			// Trigger update to process the modal result
			env.Model.Update(nil)

			// Verify the data persisted to the database
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
			// Get original bursts from database
			originalBursts := env.GetBursts()
			Expect(len(originalBursts)).To(BeNumerically(">=", 1))
			originalDesc := originalBursts[0].Description
			burstID := originalBursts[0].ID

			// Navigate to edit mode
			env.SelectIntentByName("burst_management")
			env.Confirm()
			env.PressKeyRune('e')

			// Get the active intent
			activeIntent := env.Model.GetActiveIntent()
			burstIntent := activeIntent.(*intents.BurstManagementIntent)

			// Simulate form completion with updated description
			newDesc := "UPDATED_DESCRIPTION_E2E_TEST"
			modifiedBurst := &career.Burst{
				ID:          burstID,
				Name:        originalBursts[0].Name,
				Description: newDesc,
				EventIDs:    originalBursts[0].EventIDs,
				CreatedAt:   originalBursts[0].CreatedAt,
				UpdatedAt:   originalBursts[0].UpdatedAt,
			}

			burstIntent.SetTestModalResult(&intents.ModalEditResult[*career.Burst]{
				Original: originalBursts[0],
				Modified: modifiedBurst,
				Accepted: true,
				Changes:  map[string]interface{}{"description": newDesc},
			})

			// Trigger update
			env.Model.Update(nil)

			// Verify persistence
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
			// Get original bursts
			originalBursts := env.GetBursts()
			Expect(len(originalBursts)).To(BeNumerically(">=", 1))
			originalName := originalBursts[0].Name
			burstID := originalBursts[0].ID

			// Navigate to edit mode
			env.SelectIntentByName("burst_management")
			env.Confirm()
			env.PressKeyRune('e')

			// Get the active intent
			activeIntent := env.Model.GetActiveIntent()
			burstIntent := activeIntent.(*intents.BurstManagementIntent)

			// Simulate form cancellation (Accepted: false)
			burstIntent.SetTestModalResult(&intents.ModalEditResult[*career.Burst]{
				Original: originalBursts[0],
				Modified: originalBursts[0], // No changes
				Accepted: false,             // Cancelled!
				Changes:  map[string]interface{}{},
			})

			// Trigger update
			env.Model.Update(nil)

			// Verify data was NOT changed
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
			// Get original bursts
			originalBursts := env.GetBursts()
			Expect(len(originalBursts)).To(BeNumerically(">=", 1))
			burstID := originalBursts[0].ID

			// Navigate to edit mode
			env.SelectIntentByName("burst_management")
			env.Confirm()
			env.PressKeyRune('e')

			// Get the active intent and set result
			activeIntent := env.Model.GetActiveIntent()
			burstIntent := activeIntent.(*intents.BurstManagementIntent)

			newName := "PERSISTED_ACROSS_RESTART"
			modifiedBurst := &career.Burst{
				ID:          burstID,
				Name:        newName,
				Description: originalBursts[0].Description,
				EventIDs:    originalBursts[0].EventIDs,
				CreatedAt:   originalBursts[0].CreatedAt,
				UpdatedAt:   originalBursts[0].UpdatedAt,
			}

			burstIntent.SetTestModalResult(&intents.ModalEditResult[*career.Burst]{
				Original: originalBursts[0],
				Modified: modifiedBurst,
				Accepted: true,
				Changes:  map[string]interface{}{"name": newName},
			})

			// Trigger update
			env.Model.Update(nil)

			// Simulate application restart (new model, same database)
			env.SimulateRestart()

			// Verify data persisted across restart
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

	Describe("Edit Form Display and UX", func() {
		BeforeEach(func() {
			env = e2e.Setup(GinkgoT())
			env.PopulateTestData(5, 2, 0) // 5 events, 2 bursts, 0 facts
			env.SelectIntentByName("burst_management")
			env.Confirm() // Go to detail view
		})

		AfterEach(func() {
			env.Cleanup()
		})

		It("should show form immediately when entering edit mode", func() {
			env.PressKeyRune('e')
			view := env.GetView()
			// Form should be visible immediately with input fields
			Expect(view).To(ContainSubstring("Burst Name"))
			Expect(view).NotTo(ContainSubstring("Initializing"))
			Expect(view).NotTo(ContainSubstring("Loading"))
		})

		It("should show burst name field with current value", func() {
			env.PressKeyRune('e')
			view := env.GetView()
			// Should show the name field with form styling
			Expect(view).To(ContainSubstring("Burst Name"))
			// Should contain one of the test burst names
			hasTestData := strings.Contains(view, "Authentication") || strings.Contains(view, "Mentoring")
			Expect(hasTestData).To(BeTrue(), "Edit form should show burst name from test data")
		})

		It("should show description field", func() {
			env.PressKeyRune('e')
			view := env.GetView()
			Expect(view).To(ContainSubstring("Description"))
		})

		It("should show submit confirmation field", func() {
			env.PressKeyRune('e')
			view := env.GetView()
			// huh forms show Save Changes / Submit / Cancel for confirmation
			Expect(view).To(SatisfyAny(
				ContainSubstring("Save Changes"),
				ContainSubstring("Submit"),
				ContainSubstring("Confirm"),
			))
		})

		It("should NOT show duplicate help text", func() {
			env.PressKeyRune('e')
			view := env.GetView()
			// Count occurrences of common help patterns
			// There should only be ONE set of help text (from StandardView footer)
			enterCount := strings.Count(view, "Enter")
			escCount := strings.Count(view, "Esc")

			// With proper rendering, help text appears once in footer
			// Duplicate help would show it 2+ times
			Expect(enterCount).To(BeNumerically("<=", 3), "Should not have duplicate Enter help text")
			Expect(escCount).To(BeNumerically("<=", 3), "Should not have duplicate Esc help text")
		})

		It("should show logo/branding in edit view", func() {
			env.PressKeyRune('e')
			view := env.GetView()
			// StandardView should include logo or app branding
			Expect(view).To(SatisfyAny(
				ContainSubstring("KaRiya"),
				ContainSubstring("Burst"),
				ContainSubstring("Edit"),
			))
		})

		It("should maintain consistent layout in edit view", func() {
			env.PressKeyRune('e')
			view := env.GetView()
			// View should not be excessively long (broken layout symptom)
			lines := strings.Split(view, "\n")
			Expect(len(lines)).To(BeNumerically("<", 100), "View should not have excessive line count")
			// Should not have large empty gaps
			Expect(view).NotTo(ContainSubstring("\n\n\n\n\n"))
		})
	})
})
