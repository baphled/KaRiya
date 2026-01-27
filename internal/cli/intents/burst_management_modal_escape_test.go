package intents_test

import (
	"strings"

	"github.com/baphled/kariya/internal/testutil/e2e"
	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
)

var _ = Describe("BurstManagement Modal Escape Handling", func() {
	var env *e2e.TestEnv

	BeforeEach(func() {
		env = e2e.SetupWithMemory(GinkgoT())
	})

	AfterEach(func() {
		env.Cleanup()
	})

	Describe("Edit Modal Escape - New Burst Creation", func() {
		It("should cancel creation and return to list when pressing Escape in edit modal", func() {
			env.SelectIntentByName("burst_management")
			env.AssertViewContainsAny("Burst", "List")

			env.PressKeyRune('n')
			env.AssertViewContainsAny("Edit", "Burst", "Name", "Description")

			env.Cancel()

			env.AssertViewContainsAny("Burst", "List", "Manage")
			env.AssertViewNotContains("Edit")
		})

		It("should not create burst when escape is pressed", func() {
			env.SelectIntentByName("burst_management")
			initialCount := len(env.GetBursts())

			env.PressKeyRune('n')
			env.AssertViewContainsAny("Edit", "Burst")

			env.Cancel()

			Expect(len(env.GetBursts())).To(Equal(initialCount))
		})

		It("should handle multiple escape presses gracefully", func() {
			env.SelectIntentByName("burst_management")

			env.PressKeyRune('n')
			env.AssertViewContainsAny("Edit", "Burst")

			env.Cancel()
			env.Cancel()

			view := env.GetView()
			isValid := ContainsAny(view, "Burst", "Capture Event", "Browse Timeline")
			Expect(isValid).To(BeTrue(), "Should be at burst list or main menu")
		})
	})

	Describe("Edit Modal Escape - Existing Burst Edit", func() {
		BeforeEach(func() {
			// Use PopulateTestData to create events and bursts properly
			env.PopulateTestData(5, 2, 0) // 5 events, 2 bursts, 0 facts
		})

		It("should close modal and return to detail view when pressing Escape", func() {
			env.SelectIntentByName("burst_management")
			env.AssertViewContainsAny("Team Mentoring", "List")

			// Navigate to find "Team Mentoring" burst (could be first or second depending on OS)
			view := env.GetView()
			if !strings.Contains(view, "▶ Team Mentoring") && !strings.Contains(view, "> Team Mentoring") {
				env.NavigateDown() // Move to second burst
			}

			env.Confirm()
			env.AssertViewContainsAny("Detail", "Team Mentoring", "Description")

			env.PressKeyRune('e')
			env.AssertViewContainsAny("Burst Name", "Description")

			env.Cancel()

			env.AssertViewContainsAny("Detail", "Team Mentoring", "Description")
			env.AssertViewNotContains("Burst Name")
		})

		It("should preserve original burst data when escape is pressed", func() {
			env.SelectIntentByName("burst_management")

			// Navigate to find "Team Mentoring" burst (could be first or second depending on OS)
			// Use cursor marker detection to ensure we're checking the selected item, not just any text in view
			view := env.GetView()
			if !strings.Contains(view, "▶ Team Mentoring") && !strings.Contains(view, "> Team Mentoring") {
				env.NavigateDown() // Move to second burst
			}

			env.Confirm()
			env.AssertViewContainsAny("Team Mentoring", "Focused mentoring")

			env.PressKeyRune('e')
			env.AssertViewContainsAny("Burst Name", "Description")

			env.Cancel()

			env.AssertViewContainsAny("Team Mentoring", "Focused mentoring")
		})
	})

	Describe("Global Keys in Edit Modal", func() {
		BeforeEach(func() {
			// Use PopulateTestData to create events and bursts properly
			env.PopulateTestData(5, 2, 0) // 5 events, 2 bursts, 0 facts
		})

		It("should handle '?' for help from edit modal", func() {
			env.SelectIntentByName("burst_management")
			env.Confirm()
			env.PressKeyRune('e')
			env.AssertViewContainsAny("Edit", "Burst")

			env.PressKeyRune('?')
			env.AssertViewContainsAny("Keyboard Shortcuts", "Shortcuts", "toggle help")

			env.PressKeyRune('?')
			env.AssertViewContainsAny("Edit", "Burst")
		})
	})

	Describe("Context-Aware Navigation", func() {
		It("should go to list when cancelling new burst", func() {
			env.SelectIntentByName("burst_management")
			env.PressKeyRune('n')
			env.AssertViewContainsAny("Edit", "Burst")

			env.Cancel()

			env.AssertViewContainsAny("List", "Manage", "Burst")
			env.AssertViewNotContains("Detail")
		})

		It("should go to detail when cancelling existing burst edit", func() {
			// Use PopulateTestData to create events and bursts properly
			env.PopulateTestData(5, 2, 0) // 5 events, 2 bursts, 0 facts

			env.SelectIntentByName("burst_management")
			env.Confirm()
			env.PressKeyRune('e')
			env.AssertViewContainsAny("Burst Name", "Description")

			env.Cancel()

			env.AssertViewContainsAny("Detail", "Team Mentoring")
			env.AssertViewNotContains("List")
		})
	})

	Describe("Escape Priority (Regression Prevention)", func() {
		It("should handle escape before modal processes it", func() {
			env.SelectIntentByName("burst_management")
			env.PressKeyRune('n')
			env.AssertViewContainsAny("Edit", "Burst")

			env.Cancel()

			env.AssertViewNotContains("Edit")
			env.AssertViewContainsAny("List", "Manage")
		})

		It("should not let modal consume escape key", func() {
			// Use PopulateTestData to create events and bursts properly
			env.PopulateTestData(5, 2, 0) // 5 events, 2 bursts, 0 facts

			env.SelectIntentByName("burst_management")
			env.Confirm()
			env.PressKeyRune('e')
			env.AssertViewContainsAny("Burst Name", "Description")

			env.Cancel()

			env.AssertViewNotContains("Burst Name")
			env.AssertViewContainsAny("Detail")
		})
	})
})

// Helper function to check if string contains any of the given substrings
func ContainsAny(s string, substrs ...string) bool {
	for _, substr := range substrs {
		if Contains(s, substr) {
			return true
		}
	}
	return false
}

func Contains(s, substr string) bool {
	return len(s) >= len(substr) && (s == substr || len(s) > len(substr) && (s[:len(substr)] == substr || s[len(s)-len(substr):] == substr || containsMiddle(s, substr)))
}

func containsMiddle(s, substr string) bool {
	for i := 0; i <= len(s)-len(substr); i++ {
		if s[i:i+len(substr)] == substr {
			return true
		}
	}
	return false
}
