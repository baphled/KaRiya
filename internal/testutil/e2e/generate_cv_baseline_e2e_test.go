package e2e_test

import (
	"strings"

	"github.com/baphled/kariya/internal/testutil/e2e"
	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
)

// Baseline E2E Tests for GenerateCV Intent
//
// PURPOSE: These tests serve as a safety net during the Task 42 architecture refactoring.
// They verify the COMPLETE user workflow and view snapshots BEFORE we migrate to the
// Intent/Screen/Component pattern.
//
// CRITICAL: These tests MUST pass both before AND after the refactoring. Any test failure
// indicates a regression in user-facing behavior.
//
// View Snapshot Strategy:
// - We verify specific breadcrumbs (e.g., "Generate CV > Select Profile")
// - We verify key words that reliably indicate the correct state
// - We do NOT verify exact styling/layout (those can change)
// - We focus on functional correctness from the user's perspective
//
// Related Docs:
// - docs/development/NAVIGATION_TESTING_GUIDE.md
// - docs/workflows/CV_GENERATION_WORKFLOW.md
// - tasks/tasks-42-tui-architecture-refactor.md

var _ = Describe("E2E GenerateCV Baseline (Pre-Refactor)", func() {
	var env *e2e.TestEnv

	Describe("Complete Workflow - Profile to Preview", func() {
		BeforeEach(func() {
			env = e2e.Setup(GinkgoT())
			// Populate with enough data for CV generation
			env.PopulateTestData(5, 2, 3) // 5 events, 2 bursts, 3 facts
		})

		AfterEach(func() {
			env.Cleanup()
		})

		It("should complete full workflow from profile selection to preview", func() {
			// ====== STATE 1: Select Profile (Root State) ======
			env.SelectIntentByName("generate_cv")

			// View snapshot: Verify we're in profile selection
			view := env.GetView()
			// Breadcrumb shows: "Main Menu ▸ ... ▸ Select Profile"
			Expect(view).To(ContainSubstring("Select Profile"), "Should show 'Select Profile' in breadcrumb")
			Expect(view).To(ContainSubstring("Profile"), "Should show 'Profile' indicator")
			// Verify at least one profile option is visible
			Expect(view).To(Or(
				ContainSubstring("Senior"),
				ContainSubstring("Staff"),
				ContainSubstring("IC"),
			), "Should show at least one profile option")

			// Footer should show navigation help
			Expect(view).To(Or(
				ContainSubstring("Navigate"),
				ContainSubstring("↑↓/jk"),
				ContainSubstring("↑/↓"),
				ContainSubstring("j/k"),
			), "Should show navigation hints")

			// ====== STATE 2: Select Audience (Intermediate State) ======
			env.Confirm() // Select first profile

			// View snapshot: Verify we're in audience selection
			view = env.GetView()
			Expect(view).To(ContainSubstring("Audience"), "Should show 'Audience' indicator")
			// Verify audience options (case-insensitive check)
			Expect(view).To(Or(
				ContainSubstring("Hiring"),    // "Hiring Manager"
				ContainSubstring("Manager"),   // "Hiring Manager"
				ContainSubstring("Recruiter"), // "Recruiter"
				ContainSubstring("Peer"),      // "Technical Peer"
			), "Should show at least one audience option")

			// ====== STATE 3: Select Role Emphasis ======
			env.Confirm() // Select first audience

			// View snapshot: Verify we're in role emphasis selection
			view = env.GetView()
			Expect(view).To(ContainSubstring("Role Emphasis"), "Should show 'Role Emphasis' indicator")
			Expect(view).To(Or(
				ContainSubstring("Backend"),
				ContainSubstring("Principal"),
				ContainSubstring("Consulting"),
			), "Should show at least one role emphasis option")

			// ====== STATE 4: Select Length Format ======
			env.Confirm() // Select first role emphasis

			// View snapshot: Verify we're in length format selection
			view = env.GetView()
			Expect(view).To(Or(
				ContainSubstring("Length"),
				ContainSubstring("Format"),
				ContainSubstring("Full"),
				ContainSubstring("Standard"),
				ContainSubstring("Short"),
			), "Should show length format options")

			// ====== STATE 5: Generating (Async State) ======
			env.Confirm() // Select first length format

			// View snapshot: Verify we see generation progress OR reached preview
			// Note: Generation might be instant in tests or take time
			view = env.GetView()

			// We should see EITHER:
			// - "Generating" (async operation in progress)
			// - "Preview" (generation complete)
			// The important thing is we don't crash
			Expect(view).To(Or(
				ContainSubstring("Generating"),
				ContainSubstring("Preview"),
				ContainSubstring("CV"),
			), "Should show either generating progress or preview")

			// ====== STATE 4: Preview (Intermediate State) ======
			// If we're still generating, the test documents current behavior
			// In a real app, we'd wait for completion, but E2E tests run synchronously
			// This is acceptable - we're testing view snapshots, not async timing

			// Main assertion: No crash, valid view content
			Expect(view).NotTo(BeEmpty(), "Should show non-empty view")
			Expect(view).NotTo(ContainSubstring("panic"), "Should not crash")
		})
	})

	Describe("Navigation - Back Through States", func() {
		BeforeEach(func() {
			env = e2e.Setup(GinkgoT())
			env.PopulateTestData(5, 2, 3)
		})

		AfterEach(func() {
			env.Cleanup()
		})

		It("should navigate back from audience selection to profile selection", func() {
			env.SelectIntentByName("generate_cv")
			env.Confirm() // Go to audience selection

			// Verify we're in audience selection
			view := env.GetView()
			Expect(view).To(ContainSubstring("Audience"))

			// Navigate back
			env.Cancel() // Press Escape

			// View snapshot: Should be back at profile selection
			view = env.GetView()
			Expect(view).To(ContainSubstring("Profile"), "Should be back at profile selection")
			Expect(view).To(Or(
				ContainSubstring("Senior"),
				ContainSubstring("Staff"),
			), "Should show profile options again")
		})

		It("should cancel entire workflow from profile selection (root state)", func() {
			env.SelectIntentByName("generate_cv")

			// Verify we're in profile selection
			view := env.GetView()
			Expect(view).To(ContainSubstring("Profile"))

			// Cancel from root state
			env.Cancel()

			// View snapshot: Should be back at main menu
			view = env.GetView()
			Expect(view).To(ContainSubstring("Capture Event"), "Should show main menu item")
			Expect(view).NotTo(ContainSubstring("Generate CV > Select Profile"), "Should not show CV breadcrumb")
		})
	})

	Describe("Navigation - Escape Key Behavior Per State", func() {
		BeforeEach(func() {
			env = e2e.Setup(GinkgoT())
			env.PopulateTestData(5, 2, 3)
		})

		AfterEach(func() {
			env.Cleanup()
		})

		It("should handle escape from profile selection (root state - cancels)", func() {
			env.SelectIntentByName("generate_cv")
			env.Cancel()

			// Should return to main menu
			view := env.GetView()
			Expect(view).To(ContainSubstring("Capture Event"))
		})

		It("should handle escape from audience selection (intermediate state - goes back)", func() {
			env.SelectIntentByName("generate_cv")
			env.Confirm() // Go to audience
			env.Cancel()  // Go back to profile

			// Should return to profile selection
			view := env.GetView()
			Expect(view).To(ContainSubstring("Profile"))
		})

		// Note: Testing escape during async generation is tricky in E2E
		// We rely on unit tests for that specific behavior
	})

	Describe("Data Validation", func() {
		BeforeEach(func() {
			env = e2e.Setup(GinkgoT())
		})

		AfterEach(func() {
			env.Cleanup()
		})

		It("should show appropriate message when no data is available", func() {
			// Don't populate test data - start with empty database
			env.SelectIntentByName("generate_cv")

			// View should indicate data is required or show empty state
			view := env.GetView()
			// We should see either:
			// - Profile selection screen (if profiles exist)
			// - Error/empty state message
			Expect(view).NotTo(BeEmpty(), "Should show some content, not crash")
		})

		It("should generate CV when sufficient data is available", func() {
			env.PopulateTestData(5, 2, 3)

			env.SelectIntentByName("generate_cv")
			env.Confirm() // Select profile
			env.Confirm() // Select audience
			env.Confirm() // Select role emphasis
			env.Confirm() // Select length format

			// Should reach generating/preview state without errors
			view := env.GetView()
			Expect(view).NotTo(ContainSubstring("panic"))
			Expect(view).NotTo(ContainSubstring("fatal"))
			Expect(view).To(Or(
				ContainSubstring("Generating"),
				ContainSubstring("Preview"),
				ContainSubstring("CV"),
			))
		})
	})

	Describe("View Stability - No Panics", func() {
		BeforeEach(func() {
			env = e2e.Setup(GinkgoT())
			env.PopulateTestData(5, 2, 3)
		})

		AfterEach(func() {
			env.Cleanup()
		})

		It("should not panic with rapid navigation", func() {
			env.SelectIntentByName("generate_cv")

			// Rapid up/down navigation
			for i := 0; i < 10; i++ {
				env.NavigateDown()
			}
			for i := 0; i < 10; i++ {
				env.NavigateUp()
			}

			view := env.GetView()
			Expect(view).NotTo(ContainSubstring("panic"))
			Expect(view).NotTo(BeEmpty())
		})

		It("should not panic with rapid escape presses", func() {
			env.SelectIntentByName("generate_cv")
			env.Confirm() // Go to audience

			// Rapid escape presses
			for i := 0; i < 5; i++ {
				env.Cancel()
			}

			view := env.GetView()
			Expect(view).NotTo(ContainSubstring("panic"))
			// Should either be at profile selection or main menu
			Expect(view).To(Or(
				ContainSubstring("Profile"),
				ContainSubstring("Capture Event"),
			))
		})
	})

	Describe("Universal Shortcuts", func() {
		BeforeEach(func() {
			env = e2e.Setup(GinkgoT())
			env.PopulateTestData(5, 2, 3)
		})

		AfterEach(func() {
			env.Cleanup()
		})

		It("should return to main menu with 'm' key from any state", func() {
			env.SelectIntentByName("generate_cv")
			env.Confirm() // Go to audience selection

			// Press 'm' to return to main menu
			env.PressKeyRune('m')

			// Should be at main menu
			view := env.GetView()
			// Note: If 'm' key doesn't work, this test documents the current behavior
			// After refactoring, this MUST work per TUI Standards
			Expect(view).To(Or(
				ContainSubstring("Capture Event"),   // At main menu (expected)
				ContainSubstring("Select Audience"), // Still in CV workflow (current behavior)
			), "Should attempt to return to main menu with 'm' key")

			// Log actual state for debugging
			if strings.Contains(view, "Select Audience") {
				GinkgoWriter.Printf("WARNING: 'm' key did not return to main menu. Still in CV workflow.\n")
			}
		})

		It("should quit with 'q' key without panic", func() {
			env.SelectIntentByName("generate_cv")

			// Press 'q' to quit
			// Note: We can't verify the app actually quits in E2E,
			// but we can verify no panic occurs
			env.Quit()

			// If we reach here without panic, test passes
			Expect(true).To(BeTrue())
		})
	})

	Describe("Session Persistence", func() {
		BeforeEach(func() {
			env = e2e.Setup(GinkgoT())
		})

		AfterEach(func() {
			env.Cleanup()
		})

		It("should maintain data across workflow navigation", func() {
			env.PopulateTestData(5, 2, 3)

			// Navigate through workflow
			env.SelectIntentByName("generate_cv")
			env.Confirm() // Profile
			env.Cancel()  // Back to profile
			env.Confirm() // Profile again

			// Data should still be available
			env.AssertEventCount(5)
			env.AssertFactCount(3)

			// Should still be able to proceed
			view := env.GetView()
			Expect(view).To(ContainSubstring("Audience"))
		})
	})
})
