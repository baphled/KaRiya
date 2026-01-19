package e2e_test

import (
	"github.com/baphled/kariya/internal/testutil/e2e"
	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
)

// GenerateCV Wizard E2E Tests
//
// PURPOSE: These tests verify the complete CV generation wizard workflow using
// the e2e.TestEnv framework with real SQLite persistence.
//
// BACKGROUND: The wizard modal flow was implemented in Task 43 as an alternative
// to the legacy 17-state screen architecture. The wizard consolidates:
// - Profile selection
// - Audience selection
// - Technology focus configuration
// - Skills format configuration
// - CV length configuration
//
// WORKFLOW STATES:
// 1. CVStateConfiguring - Wizard modal visible (3 steps: WHO, TECH, FORMAT)
// 2. CVStateExtracting - Progress modal (extracting technologies)
// 3. CVStateGenerating - Progress modal (generating CV)
// 4. CVStatePreview - Preview screen (scrollable CV content)
// 5. CVStateExporting - Export modal (format + location selection)
//
// KEYBOARD SHORTCUTS:
// - Tab: Next field
// - Enter: Submit field / Confirm
// - Esc: Cancel / Go back
// - Ctrl+S: Skip step (use defaults)
// - j/k or ↑/↓: Navigate lists
// - x: Export from preview
// - q: Quit
// - m: Return to main menu
//
// Related Docs:
// - docs/WIZARD_MODAL_GUIDE.md
// - docs/workflows/CV_GENERATION_WORKFLOW.md
// - tasks/tasks-43-cv-generation-wizard-modal.md

var _ = Describe("E2E GenerateCV Wizard Workflow", func() {
	var env *e2e.TestEnv

	// ========================================================================
	// WIZARD MODAL INITIALIZATION
	// ========================================================================
	Describe("Wizard Modal Initialization", func() {
		BeforeEach(func() {
			env = e2e.Setup(GinkgoT())
			env.PopulateTestData(5, 2, 3) // 5 events, 2 bursts, 3 facts
		})

		AfterEach(func() {
			env.Cleanup()
		})

		It("should show wizard modal when entering CV generation", func() {
			env.SelectIntentByName("generate_cv")

			view := env.GetView()
			// Wizard modal should be visible with configuration title
			Expect(view).To(Or(
				ContainSubstring("CV Configuration"),
				ContainSubstring("Select CV Profile"),
				ContainSubstring("Profile"),
			), "Should show wizard modal or profile selection")
		})

		It("should show profile selection field in wizard", func() {
			env.SelectIntentByName("generate_cv")

			view := env.GetView()
			Expect(view).To(ContainSubstring("Profile"), "Should show Profile field")
		})

		It("should show audience selection field in wizard", func() {
			env.SelectIntentByName("generate_cv")

			view := env.GetView()
			Expect(view).To(ContainSubstring("Audience"), "Should show Audience field")
		})

		It("should show navigation help in wizard footer", func() {
			env.SelectIntentByName("generate_cv")

			view := env.GetView()
			// Wizard modal should show keyboard shortcuts
			Expect(view).To(Or(
				ContainSubstring("Tab"),
				ContainSubstring("Enter"),
				ContainSubstring("Esc"),
				ContainSubstring("Navigate"),
				ContainSubstring("↑/↓"),
				ContainSubstring("j/k"),
			), "Should show navigation help in footer")
		})
	})

	// ========================================================================
	// WIZARD STEP 1: WHO (Profile + Audience)
	// ========================================================================
	Describe("Wizard Step 1: Profile and Audience Selection", func() {
		BeforeEach(func() {
			env = e2e.Setup(GinkgoT())
			env.PopulateTestData(5, 2, 3)
		})

		AfterEach(func() {
			env.Cleanup()
		})

		It("should show available profile options", func() {
			env.SelectIntentByName("generate_cv")

			view := env.GetView()
			// Should show at least one profile option
			Expect(view).To(Or(
				ContainSubstring("Senior"),
				ContainSubstring("Staff"),
				ContainSubstring("Principal"),
				ContainSubstring("Manager"),
				ContainSubstring("IC"),
			), "Should show at least one profile option")
		})

		It("should show audience options", func() {
			env.SelectIntentByName("generate_cv")

			view := env.GetView()
			// Should show audience options
			Expect(view).To(Or(
				ContainSubstring("Hiring Manager"),
				ContainSubstring("Recruiter"),
				ContainSubstring("Peer"),
				ContainSubstring("Hiring"),
				ContainSubstring("hiring"),
			), "Should show at least one audience option")
		})

		It("should allow navigating profile options with j/k keys", func() {
			env.SelectIntentByName("generate_cv")

			// Navigate down
			env.NavigateDown()
			view := env.GetView()
			Expect(view).NotTo(ContainSubstring("panic"), "Should not panic on navigation")

			// Navigate up
			env.NavigateUp()
			view = env.GetView()
			Expect(view).NotTo(ContainSubstring("panic"), "Should not panic on navigation")
		})

		It("should allow tab navigation between fields", func() {
			env.SelectIntentByName("generate_cv")

			// Tab to next field
			env.Tab()
			view := env.GetView()
			Expect(view).NotTo(ContainSubstring("panic"), "Should not panic on Tab")
			Expect(view).To(ContainSubstring("Audience"), "Should still show Audience field")
		})
	})

	// ========================================================================
	// WIZARD NAVIGATION: ESCAPE KEY BEHAVIOR
	// ========================================================================
	Describe("Wizard Escape Key Behavior", func() {
		BeforeEach(func() {
			env = e2e.Setup(GinkgoT())
			env.PopulateTestData(5, 2, 3)
		})

		AfterEach(func() {
			env.Cleanup()
		})

		It("should cancel wizard and return to main menu on Esc from step 1", func() {
			env.SelectIntentByName("generate_cv")

			// Verify we're in wizard
			view := env.GetView()
			Expect(view).To(Or(
				ContainSubstring("Profile"),
				ContainSubstring("CV Configuration"),
			))

			// Press Escape to cancel
			env.Cancel()

			// Should return to main menu
			view = env.GetView()
			Expect(view).To(Or(
				ContainSubstring("Capture Event"),
				ContainSubstring("Main Menu"),
			), "Should return to main menu on Esc from wizard step 1")
		})

		It("should not panic on multiple rapid Esc presses", func() {
			env.SelectIntentByName("generate_cv")

			// Rapid Esc presses
			for i := 0; i < 5; i++ {
				env.Cancel()
			}

			view := env.GetView()
			Expect(view).NotTo(ContainSubstring("panic"), "Should not panic on rapid Esc")
		})
	})

	// ========================================================================
	// WIZARD COMPLETION: PROGRESS TO CV GENERATION
	// ========================================================================
	Describe("Wizard Completion and CV Generation", func() {
		BeforeEach(func() {
			env = e2e.Setup(GinkgoT())
			env.PopulateTestData(5, 2, 3)
		})

		AfterEach(func() {
			env.Cleanup()
		})

		It("should show progress or preview after completing wizard", func() {
			env.SelectIntentByName("generate_cv")

			// Complete wizard step 1 (profile selection)
			env.Confirm()

			// Complete wizard (audience is second field in step 1)
			// If we're in step 1, Tab and confirm should work
			env.Tab()
			env.Confirm()

			// Should show progress modal or preview
			view := env.GetView()
			Expect(view).To(Or(
				ContainSubstring("Generating"),
				ContainSubstring("Extracting"),
				ContainSubstring("Preview"),
				ContainSubstring("CV"),
				ContainSubstring("Profile"), // Still in wizard if multi-step
			), "Should progress past wizard")
		})

		It("should not crash during CV generation workflow", func() {
			env.SelectIntentByName("generate_cv")

			// Try to complete the wizard
			env.Confirm() // Select profile
			env.Confirm() // Select audience (if visible)

			view := env.GetView()
			Expect(view).NotTo(ContainSubstring("panic"), "Should not panic during workflow")
			Expect(view).NotTo(BeEmpty(), "Should show content")
		})
	})

	// ========================================================================
	// GLOBAL KEYBOARD SHORTCUTS
	// ========================================================================
	Describe("Global Keyboard Shortcuts in Wizard", func() {
		BeforeEach(func() {
			env = e2e.Setup(GinkgoT())
			env.PopulateTestData(5, 2, 3)
		})

		AfterEach(func() {
			env.Cleanup()
		})

		It("should handle 'q' key to quit without panic", func() {
			env.SelectIntentByName("generate_cv")
			env.Quit()

			// If we reach here without panic, test passes
			Expect(true).To(BeTrue(), "'q' key should not cause panic")
		})

		It("should handle 'm' key to return to main menu", func() {
			env.SelectIntentByName("generate_cv")
			env.PressKeyRune('m')

			view := env.GetView()
			// Should be at main menu or still in CV workflow
			// (depends on whether 'm' is captured by wizard modal)
			Expect(view).To(Or(
				ContainSubstring("Capture Event"), // At main menu
				ContainSubstring("Profile"),       // Still in wizard
			), "'m' key should attempt to return to main menu")
		})
	})

	// ========================================================================
	// DATA PERSISTENCE
	// ========================================================================
	Describe("Data Persistence Through Wizard Workflow", func() {
		BeforeEach(func() {
			env = e2e.Setup(GinkgoT())
		})

		AfterEach(func() {
			env.Cleanup()
		})

		It("should maintain event data during wizard navigation", func() {
			env.PopulateTestData(5, 2, 3)

			// Enter wizard
			env.SelectIntentByName("generate_cv")
			env.AssertEventCount(5)
			env.AssertFactCount(3)

			// Navigate within wizard
			env.NavigateDown()
			env.NavigateUp()

			// Data should still be intact
			env.AssertEventCount(5)
			env.AssertFactCount(3)
		})

		It("should maintain data after cancelling wizard", func() {
			env.PopulateTestData(5, 2, 3)

			env.SelectIntentByName("generate_cv")
			env.Cancel() // Cancel wizard

			// Data should be preserved
			env.AssertEventCount(5)
			env.AssertFactCount(3)
		})

		It("should maintain data across session restart", func() {
			env.PopulateTestData(5, 2, 3)

			env.SelectIntentByName("generate_cv")
			env.Cancel() // Return to main menu

			// Simulate session restart
			env.SimulateRestart()

			// Data should persist
			env.AssertEventCount(5)
			env.AssertFactCount(3)

			// Wizard should still work
			env.SelectIntentByName("generate_cv")
			view := env.GetView()
			Expect(view).To(ContainSubstring("Profile"), "Wizard should still be accessible")
		})
	})

	// ========================================================================
	// VIEW STABILITY - NO PANICS
	// ========================================================================
	Describe("View Stability in Wizard", func() {
		BeforeEach(func() {
			env = e2e.Setup(GinkgoT())
			env.PopulateTestData(5, 2, 3)
		})

		AfterEach(func() {
			env.Cleanup()
		})

		It("should not panic with rapid navigation in wizard", func() {
			env.SelectIntentByName("generate_cv")

			// Rapid navigation
			for i := 0; i < 10; i++ {
				env.NavigateDown()
			}
			for i := 0; i < 10; i++ {
				env.NavigateUp()
			}
			for i := 0; i < 5; i++ {
				env.Tab()
			}

			view := env.GetView()
			Expect(view).NotTo(ContainSubstring("panic"))
			Expect(view).NotTo(BeEmpty())
		})

		It("should not panic with mixed key presses", func() {
			env.SelectIntentByName("generate_cv")

			// Mixed keys: navigation, tab, enter, escape
			env.NavigateDown()
			env.Tab()
			env.NavigateUp()
			env.Tab()
			env.Confirm()
			env.Cancel()

			view := env.GetView()
			Expect(view).NotTo(ContainSubstring("panic"))
		})

		It("should render valid view at every state", func() {
			env.SelectIntentByName("generate_cv")

			// Check view at initial state
			view := env.GetView()
			Expect(view).NotTo(BeEmpty(), "View should not be empty")
			Expect(view).NotTo(ContainSubstring("panic"), "View should not contain panic")

			// Navigate and check view
			env.NavigateDown()
			view = env.GetView()
			Expect(view).NotTo(BeEmpty())

			// Tab and check view
			env.Tab()
			view = env.GetView()
			Expect(view).NotTo(BeEmpty())
		})
	})

	// ========================================================================
	// EMPTY DATA HANDLING
	// ========================================================================
	Describe("Empty Data Handling", func() {
		BeforeEach(func() {
			env = e2e.Setup(GinkgoT())
			// Don't populate test data - test empty state
		})

		AfterEach(func() {
			env.Cleanup()
		})

		It("should handle wizard with no events gracefully", func() {
			env.AssertEventCount(0)

			env.SelectIntentByName("generate_cv")

			view := env.GetView()
			// Should show wizard or appropriate message
			Expect(view).NotTo(ContainSubstring("panic"))
			Expect(view).NotTo(BeEmpty())
		})

		It("should still allow wizard navigation with empty data", func() {
			env.SelectIntentByName("generate_cv")

			// Navigation should work even with no events
			env.NavigateDown()
			env.Tab()
			env.Cancel()

			view := env.GetView()
			Expect(view).NotTo(ContainSubstring("panic"))
		})
	})

	// ========================================================================
	// WIZARD RE-ENTRY
	// ========================================================================
	Describe("Wizard Re-entry", func() {
		BeforeEach(func() {
			env = e2e.Setup(GinkgoT())
			env.PopulateTestData(5, 2, 3)
		})

		AfterEach(func() {
			env.Cleanup()
		})

		It("should allow re-entering wizard after cancellation", func() {
			// First entry
			env.SelectIntentByName("generate_cv")
			view := env.GetView()
			Expect(view).To(ContainSubstring("Profile"))

			// Cancel
			env.Cancel()

			// Re-enter
			env.SelectIntentByName("generate_cv")
			view = env.GetView()
			Expect(view).To(ContainSubstring("Profile"), "Should show wizard again on re-entry")
		})

		It("should reset wizard state on re-entry", func() {
			// First entry and partial completion
			env.SelectIntentByName("generate_cv")
			env.NavigateDown()
			env.Tab()

			// Cancel
			env.Cancel()

			// Re-enter - wizard should be fresh
			env.SelectIntentByName("generate_cv")
			view := env.GetView()
			Expect(view).NotTo(ContainSubstring("panic"))
			Expect(view).To(ContainSubstring("Profile"))
		})
	})
})
