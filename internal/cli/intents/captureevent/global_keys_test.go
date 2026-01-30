package captureevent

import (
	"time"

	"github.com/baphled/kariya/internal/cli/intents"
	"github.com/baphled/kariya/internal/domain/career"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/google/uuid"
	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
)

// Global Keys Enforcement tests for CaptureEvent intent.
// Extracted from intents/global_keys_enforcement_e2e_test.go during
// subdirectory migration.

var _ = Describe("CaptureEvent Global Keys Enforcement", func() {

	Context("Root State Escape Behavior", func() {
		It("should cancel intent on escape from root state", func() {
			intent, err := NewIntent(&IntentContext{
				CaptureStrategy: "manual",
				Metadata:        make(map[string]string),
			})
			Expect(err).NotTo(HaveOccurred())

			intent.Init()
			intent.Update(tea.KeyMsg{Type: tea.KeyEsc})

			result := intent.Result()
			Expect(result).NotTo(BeNil(),
				"CaptureEvent: Root state escape should produce result")
		})
	})

	Context("Quit Key Behavior", func() {
		It("should ignore 'q' key in root state (quit only from main menu)", func() {
			intent, err := NewIntent(&IntentContext{
				CaptureStrategy: "manual",
				Metadata:        make(map[string]string),
			})
			Expect(err).NotTo(HaveOccurred())

			intent.Init()
			cmd := intent.Update(tea.KeyMsg{Type: tea.KeyRunes, Runes: []rune{'q'}})

			Expect(cmd).To(BeNil(),
				"CaptureEvent: 'q' key should be ignored within intent (quit only from main menu)")
		})
	})

	Context("Help Key Behavior", func() {
		It("should toggle help on '?' in root state", func() {
			intent, err := NewIntent(&IntentContext{
				CaptureStrategy: "manual",
				Metadata:        make(map[string]string),
			})
			Expect(err).NotTo(HaveOccurred())

			intent.Init()

			baseIntent := intent.BaseIntent
			if baseIntent != nil {
				helpBefore := baseIntent.IsHelpVisible()
				intent.Update(tea.KeyMsg{Type: tea.KeyRunes, Runes: []rune{'?'}})
				Expect(baseIntent.IsHelpVisible()).NotTo(Equal(helpBefore),
					"CaptureEvent: Help key should toggle help modal")
			}
		})
	})

	Context("Edit vs New - Context-Aware Navigation", func() {
		It("should cancel when editing existing event (PreviousEvent != nil)", func() {
			existingEvent := &career.Event{
				ID:   uuid.New().String(),
				Text: "Existing event",
				Date: time.Now(),
			}

			ctx := &IntentContext{
				CaptureStrategy: "manual",
				PreviousEvent:   existingEvent,
				Metadata:        make(map[string]string),
			}

			intent, err := NewIntent(ctx)
			Expect(err).NotTo(HaveOccurred())
			intent.Init()

			// Transition to form state.
			intent.SetStateForTesting(StateForm)

			// Press escape - should cancel (return to caller like BrowseTimeline).
			intent.Update(tea.KeyMsg{Type: tea.KeyEsc})

			result := intent.Result()
			Expect(result).NotTo(BeNil())
			Expect(result.Status).To(Equal(intents.Cancelled))
		})

		It("should go back when creating new event (PreviousEvent == nil)", func() {
			ctx := &IntentContext{
				CaptureStrategy: "manual",
				PreviousEvent:   nil,
				Metadata:        make(map[string]string),
			}

			intent, err := NewIntent(ctx)
			Expect(err).NotTo(HaveOccurred())
			intent.Init()

			// Transition to form state.
			intent.SetStateForTesting(StateForm)

			// Press escape - should go back to strategy selection.
			intent.Update(tea.KeyMsg{Type: tea.KeyEsc})

			Expect(intent.GetState()).To(Equal(string(StateChooseStrategy)))
			Expect(intent.Result()).To(BeNil())
		})
	})
})
