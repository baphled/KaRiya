package e2e_test

import (
	"os"
	"strings"
	"time"

	"github.com/baphled/kariya/internal/cli/bootstrap"
	"github.com/baphled/kariya/internal/config"
	"github.com/baphled/kariya/internal/testutil/e2e"
	tea "github.com/charmbracelet/bubbletea"
	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
)

// onboardingTestHelper wraps OnboardingTestModel for easier testing.
type onboardingTestHelper struct {
	model *bootstrap.OnboardingTestModel
}

// newOnboardingTestHelper creates a new test helper with initialized model.
func newOnboardingTestHelper() *onboardingTestHelper {
	model := bootstrap.NewOnboardingTestModel(nil)
	helper := &onboardingTestHelper{
		model: model,
	}
	// Get init command.
	initCmd := model.Init()
	// Send window size first.
	helper.sendWindowSize(80, 24)
	// Process init command to initialize form content.
	helper.processFormCmds(initCmd, 10)
	return helper
}

// processFormCmds processes commands from form interactions, including batch messages.
func (h *onboardingTestHelper) processFormCmds(cmd tea.Cmd, maxDepth int) {
	if cmd == nil || maxDepth <= 0 {
		return
	}

	// Run command with timeout to handle async commands like cursor blink.
	done := make(chan tea.Msg, 1)
	go func() {
		msg := cmd()
		done <- msg
	}()

	var msg tea.Msg
	select {
	case msg = <-done:
	case <-time.After(50 * time.Millisecond):
		// Timeout - command is async, skip it.
		return
	}

	if msg == nil {
		return
	}

	switch m := msg.(type) {
	case tea.BatchMsg:
		// BatchMsg contains multiple commands - process each one.
		for _, batchCmd := range m {
			h.processFormCmds(batchCmd, maxDepth-1)
		}
	default:
		// Process the message and any follow-up commands.
		_, nextCmd := h.model.Update(msg)
		h.processFormCmds(nextCmd, maxDepth-1)
	}
}

// sendWindowSize sends a window size message to the model.
func (h *onboardingTestHelper) sendWindowSize(width, height int) {
	h.model.Update(tea.WindowSizeMsg{Width: width, Height: height})
}

// view returns the current view.
func (h *onboardingTestHelper) view() string {
	return h.model.View()
}

// typeText types text character by character.
func (h *onboardingTestHelper) typeText(text string) {
	for _, r := range text {
		_, cmd := h.model.Update(tea.KeyMsg{Type: tea.KeyRunes, Runes: []rune{r}})
		h.processFormCmds(cmd, 5)
	}
}

// pressKey sends a key press and processes resulting commands.
func (h *onboardingTestHelper) pressKey(key string) {
	var msg tea.KeyMsg
	switch key {
	case "enter":
		msg = tea.KeyMsg{Type: tea.KeyEnter}
	case "tab":
		msg = tea.KeyMsg{Type: tea.KeyTab}
	case "shift+tab":
		msg = tea.KeyMsg{Type: tea.KeyShiftTab}
	case "up":
		msg = tea.KeyMsg{Type: tea.KeyUp}
	case "down":
		msg = tea.KeyMsg{Type: tea.KeyDown}
	default:
		msg = tea.KeyMsg{Type: tea.KeyRunes, Runes: []rune(key)}
	}
	_, cmd := h.model.Update(msg)
	h.processFormCmds(cmd, 10)
}

// pressEnter presses enter key.
func (h *onboardingTestHelper) pressEnter() {
	h.pressKey("enter")
}

// pressTab presses tab key.
func (h *onboardingTestHelper) pressTab() {
	h.pressKey("tab")
}

// viewContains checks if view contains the given text.
func (h *onboardingTestHelper) viewContains(text string) bool {
	return strings.Contains(h.view(), text)
}

// viewContainsAny checks if view contains any of the given texts.
func (h *onboardingTestHelper) viewContainsAny(texts ...string) bool {
	view := h.view()
	for _, text := range texts {
		if strings.Contains(view, text) {
			return true
		}
	}
	return false
}

// isCompleted returns whether the wizard completed.
func (h *onboardingTestHelper) isCompleted() bool {
	return h.model.IsCompleted()
}

// result returns the profile config result.
func (h *onboardingTestHelper) result() *config.ProfileConfig {
	return h.model.Result()
}

var _ = Describe("E2E Onboarding Wizard Workflow", func() {
	var env *e2e.TestEnv

	Describe("Onboarding Initialization", func() {
		var helper *onboardingTestHelper

		BeforeEach(func() {
			helper = newOnboardingTestHelper()
		})

		It("should show onboarding wizard on fresh startup", func() {
			view := helper.view()
			Expect(view).NotTo(BeEmpty(), "Onboarding wizard should render")
		})

		It("should show Profile Setup title", func() {
			Expect(helper.viewContains("Profile Setup")).To(BeTrue(),
				"Should show Profile Setup title")
		})

		It("should show Step 1 of 3", func() {
			Expect(helper.viewContains("Step 1 of 3")).To(BeTrue(),
				"Should show Step 1 of 3")
		})

		It("should show Welcome message", func() {
			Expect(helper.viewContains("Welcome to KaRiya")).To(BeTrue(),
				"Should show Welcome message")
		})

		It("should show Name field", func() {
			Expect(helper.viewContainsAny("Name", "name")).To(BeTrue(),
				"Should show Name field")
		})

		It("should show keyboard shortcuts", func() {
			Expect(helper.viewContainsAny("tab", "enter", "next", "Tab", "Enter")).To(BeTrue(),
				"Should show keyboard shortcuts")
		})
	})

	Describe("Onboarding Step Navigation", func() {
		var helper *onboardingTestHelper

		BeforeEach(func() {
			helper = newOnboardingTestHelper()
		})

		It("should remain on Step 1 without entering name", func() {
			helper.pressEnter()
			Expect(helper.viewContains("Step 1 of 3")).To(BeTrue(),
				"Should remain on Step 1 without valid name")
		})

		It("should accept typed name text", func() {
			helper.typeText("Test User")
			Expect(helper.viewContains("Test User")).To(BeTrue(),
				"Should display typed name")
		})

		It("should advance to Step 2 after entering valid name", func() {
			helper.typeText("Test User")
			helper.pressEnter()
			Expect(helper.viewContains("Step 2 of 3")).To(BeTrue(),
				"Should advance to Step 2")
		})

		It("should show Email field on Step 2", func() {
			helper.typeText("Test User")
			helper.pressEnter()
			Expect(helper.viewContainsAny("Email", "email")).To(BeTrue(),
				"Should show Email field on Step 2")
		})

		It("should show Location field on Step 2", func() {
			helper.typeText("Test User")
			helper.pressEnter()
			Expect(helper.viewContainsAny("Location", "location")).To(BeTrue(),
				"Should show Location field on Step 2")
		})

		It("should advance to Step 3 after completing Step 2", func() {
			// Step 1: Name.
			helper.typeText("Test User")
			helper.pressEnter()
			// Step 2: Email and Location.
			helper.typeText("test@example.com")
			helper.pressTab()
			helper.typeText("London")
			helper.pressEnter()
			Expect(helper.viewContains("Step 3 of 3")).To(BeTrue(),
				"Should advance to Step 3")
		})

		It("should show Professional Details on Step 3", func() {
			// Step 1: Name.
			helper.typeText("Test User")
			helper.pressEnter()
			// Step 2: Email and Location.
			helper.typeText("test@example.com")
			helper.pressTab()
			helper.typeText("London")
			helper.pressEnter()
			Expect(helper.viewContainsAny("Professional", "Title", "GitHub", "Portfolio")).To(BeTrue(),
				"Should show Professional Details on Step 3")
		})
	})

	Describe("Onboarding Completion", func() {
		var helper *onboardingTestHelper

		BeforeEach(func() {
			helper = newOnboardingTestHelper()
		})

		It("should complete onboarding and return profile config", func() {
			// Step 1: Name (1 field).
			helper.typeText("Test User")
			helper.pressEnter()
			// Step 2: Email and Location (2 fields).
			helper.typeText("test@example.com")
			helper.pressTab()
			helper.typeText("London")
			helper.pressEnter()
			// Step 3: Professional details - 3 optional fields (Title, GitHub, Portfolio).
			// Tab through all fields and press enter on the last one to complete.
			helper.pressTab()   // Skip Title field.
			helper.pressTab()   // Skip GitHub field.
			helper.pressEnter() // Submit on Portfolio field to complete form.

			Expect(helper.isCompleted()).To(BeTrue(),
				"Wizard should be completed after all steps")
			result := helper.result()
			Expect(result).NotTo(BeNil(), "Should return profile config")
			Expect(result.Name).To(Equal("Test User"), "Name should match")
			Expect(result.Email).To(Equal("test@example.com"), "Email should match")
			Expect(result.Location).To(Equal("London"), "Location should match")
		})
	})

	Describe("Onboarding with Existing E2E Setup", func() {
		BeforeEach(func() {
			// Use regular Setup (which skips onboarding)
			env = e2e.GetSharedEnv(GinkgoT())
		})

		AfterEach(func() {
			env.Cleanup()
		})

		It("should skip onboarding with standard Setup", func() {
			// Regular E2E Setup should skip onboarding (via bootstrap.SkipOnboarding)
			// The app starts directly in menu state
			Expect(env.IsInOnboardingState()).To(BeFalse(), "Regular Setup should skip onboarding")
		})

		It("should go directly to menu with standard Setup", func() {
			// Should be at menu immediately
			Expect(env.IsInMenuState()).To(BeTrue(), "Should start at menu with standard Setup")
		})
	})

	Describe("BUG-007 Regression: Config File Isolation", func() {
		var (
			realConfigPath  string
			originalContent []byte
			originalExists  bool
		)

		BeforeEach(func() {
			// Get the real config path BEFORE any test setup
			var err error
			realConfigPath, err = config.GetConfigPath()
			Expect(err).NotTo(HaveOccurred())

			// Read original content if it exists
			originalContent, err = os.ReadFile(realConfigPath)
			if err == nil {
				originalExists = true
			} else if os.IsNotExist(err) {
				originalExists = false
			} else {
				Fail("Failed to read original config: " + err.Error())
			}

			// Now set up the test environment
			env = e2e.SetupWithOnboarding(GinkgoT())
		})

		AfterEach(func() {
			env.Cleanup()
		})

		It("BUG-007: config file isolation should work in test setup", func() {
			// Verify we're using an isolated config path
			testConfigPath, err := config.GetConfigPath()
			Expect(err).NotTo(HaveOccurred())
			Expect(testConfigPath).NotTo(Equal(realConfigPath),
				"Test should use isolated config path, not real config")
		})

		It("BUG-007: should NOT write to user's real config file", func() {
			// Now verify the real config file was NOT modified
			if originalExists {
				currentContent, err := os.ReadFile(realConfigPath)
				Expect(err).NotTo(HaveOccurred(), "Should be able to read real config")
				Expect(currentContent).To(Equal(originalContent),
					"Real config file should NOT have been modified by test")
			} else {
				// If config didn't exist before, it should still not exist
				_, err := os.Stat(realConfigPath)
				Expect(os.IsNotExist(err)).To(BeTrue(),
					"Real config file should NOT have been created by test")
			}
		})
	})
})
