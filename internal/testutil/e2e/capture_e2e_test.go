package e2e_test

import (
	"strings"
	"time"

	"github.com/baphled/kariya/internal/cli/intents/captureevent"
	"github.com/baphled/kariya/internal/cli/uikit/feedback"
	"github.com/baphled/kariya/internal/testutil/e2e"
	"github.com/baphled/kariya/internal/testutil/fixtures"
	tea "github.com/charmbracelet/bubbletea"
	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
)

var _ = Describe("E2E Capture Workflow", func() {
	var env *e2e.TestEnv

	Describe("Database Persistence", func() {
		BeforeEach(func() {
			env = e2e.GetSharedEnv(GinkgoT())
		})

		AfterEach(func() {
			env.Cleanup()
		})

		It("should start with empty database", func() {
			env.AssertEventCount(0)
		})

		It("should not persist event when cancelled before submission", func() {
			env.SelectIntentByName("capture_event")
			env.Confirm()
			env.TypeText("Test event text")
			env.Cancel()
			env.AssertEventCount(0)
		})

		It("should not persist event when returning to main menu", func() {
			env.SelectIntentByName("capture_event")
			env.Confirm()
			env.TypeText("Test event text")
			env.PressKeyRune('m')
			env.AssertEventCount(0)
		})
	})

	Describe("Session Persistence", func() {
		BeforeEach(func() {
			env = e2e.GetSharedEnv(GinkgoT())
		})

		AfterEach(func() {
			env.Cleanup()
		})

		It("should show empty timeline after restart with no events", func() {
			env.SimulateRestart()
			env.SelectIntentByName("browse_timeline")
			env.AssertViewContainsAny("No events", "empty", "Timeline")
		})
	})

	Describe("Save Dialog Countdown", func() {
		BeforeEach(func() {
			env = e2e.GetSharedEnv(GinkgoT())
		})

		AfterEach(func() {
			env.Cleanup()
		})

		It("should show success modal with countdown after event submission", func() {
			env.SelectIntentByName("capture_event")
			env.Confirm()

			env.SendMessage(captureevent.SubmitCompleteMsg{})

			view := env.GetView()
			Expect(view).To(ContainSubstring("Success"), "Should show success modal")
			Expect(view).To(ContainSubstring("Event saved!"), "Should show success message")
			Expect(view).To(ContainSubstring("Auto-dismiss"), "Should show auto-dismiss countdown")
		})

		It("should update countdown display as time passes", func() {
			env.SelectIntentByName("capture_event")
			env.Confirm()

			env.SendMessage(captureevent.SubmitCompleteMsg{})

			initialView := env.GetView()
			Expect(initialView).To(ContainSubstring("Auto-dismiss"), "Should show auto-dismiss countdown")
			Expect(initialView).To(MatchRegexp(`Auto-dismiss in [123]s`), "Countdown should show 1-3 seconds remaining")
		})

		It("should auto-dismiss modal when countdown reaches zero", func() {
			env.SelectIntentByName("capture_event")
			env.Confirm()

			env.SendMessage(captureevent.SubmitCompleteMsg{})

			for range 4 {
				env.SendMessage(feedback.ModalCountdownTickMsg{})
				time.Sleep(1 * time.Second)
			}

			env.SendMessage(feedback.ModalAutoDismissMsg{})

			view := env.GetView()
			Expect(view).To(ContainSubstring("Review Enrichment"), "Should transition to review screen after auto-dismiss")
			Expect(view).ToNot(ContainSubstring("Auto-dismiss"), "Success modal should be dismissed")
		})

		It("should allow manual dismiss with Esc before countdown completes", func() {
			env.SelectIntentByName("capture_event")
			env.Confirm()

			env.SendMessage(captureevent.SubmitCompleteMsg{})

			view := env.GetView()
			Expect(view).To(ContainSubstring("Auto-dismiss"), "Should show countdown before manual dismiss")

			env.PressKey(tea.KeyEsc)

			viewAfterEsc := env.GetView()
			Expect(viewAfterEsc).ToNot(ContainSubstring("Auto-dismiss"), "Success modal should be dismissed after Esc")
		})

		It("should decrement countdown display on each tick", func() {
			env.SelectIntentByName("capture_event")
			env.Confirm()

			env.SendMessage(captureevent.SubmitCompleteMsg{})

			initialView := env.GetView()
			Expect(initialView).To(ContainSubstring("Auto-dismiss in 3s"))

			env.SendMessage(feedback.ModalCountdownTickMsg{})
			viewAfterFirstTick := env.GetView()
			Expect(viewAfterFirstTick).To(ContainSubstring("Auto-dismiss in 2s"))

			env.SendMessage(feedback.ModalCountdownTickMsg{})
			viewAfterSecondTick := env.GetView()
			Expect(viewAfterSecondTick).To(ContainSubstring("Auto-dismiss in 1s"))

			env.SendMessage(feedback.ModalCountdownTickMsg{})
			env.SendMessage(feedback.ModalAutoDismissMsg{})

			viewAfterAutoDismiss := env.GetView()
			Expect(viewAfterAutoDismiss).ToNot(ContainSubstring("Auto-dismiss"))
			Expect(viewAfterAutoDismiss).To(ContainSubstring("Review Enrichment"))
		})
	})

	Describe("Metadata Editor Form Scrolling", func() {
		BeforeEach(func() {
			env = e2e.GetSharedEnv(GinkgoT())
		})

		AfterEach(func() {
			env.Cleanup()
		})

		It("should reach review screen after submitting event", func() {
			env.SelectIntentByName("capture_event")
			env.Confirm()

			testEvent := fixtures.EventWith("", "Built REST API with Go", "", "")
			testEvent.ID = ""
			env.SubmitEvent(testEvent)

			view := env.GetView()
			Expect(strings.Contains(view, "Review Enrichment")).To(BeTrue(),
				"Should be on review screen.\nView len: %d", len(view))
		})

		It("should open metadata editor when pressing e on review screen", func() {
			env.SelectIntentByName("capture_event")
			env.Confirm()

			testEvent := fixtures.EventWith("", "Built REST API with Go", "", "")
			testEvent.ID = ""
			env.SubmitEvent(testEvent)

			env.PressKeyRune('e')

			view := env.GetView()
			Expect(view).ToNot(BeEmpty(),
				"View should not be empty after pressing 'e'")
			Expect(strings.Contains(view, "Edit Event Metadata")).To(BeTrue(),
				"Should show metadata editor modal title")
		})

		It("should close metadata editor when pressing escape", func() {
			env.SelectIntentByName("capture_event")
			env.Confirm()

			testEvent := fixtures.EventWith("", "Built REST API with Go", "", "")
			testEvent.ID = ""
			env.SubmitEvent(testEvent)

			env.PressKeyRune('e')
			env.Cancel()

			view := env.GetView()
			Expect(strings.Contains(view, "Review Enrichment")).To(BeTrue(),
				"Should return to review screen after escape.\nView len: %d", len(view))
		})

		It("should show form fields in metadata editor", func() {
			env.SelectIntentByName("capture_event")
			env.Confirm()

			testEvent := fixtures.EventWith("", "Built REST API with Go", "", "")
			testEvent.ID = ""
			env.SubmitEvent(testEvent)

			env.PressKeyRune('e')

			view := env.GetView()
			Expect(view).To(SatisfyAny(
				ContainSubstring("Date"),
				ContainSubstring("Company"),
				ContainSubstring("Project"),
			), "Metadata editor should show at least one form field")
		})
	})
})
