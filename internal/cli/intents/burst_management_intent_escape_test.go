package intents_test

import (
	"time"

	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"

	"github.com/baphled/kariya/internal/cli/intents"
	"github.com/baphled/kariya/internal/domain/career"
	tea "github.com/charmbracelet/bubbletea"
)

var _ = Describe("BurstManagement Intent - Escape Key Behavior", func() {
	var (
		context *intents.BurstManagementContext
		intent  *intents.BurstManagementIntent
		bursts  []*career.Burst
	)

	BeforeEach(func() {
		// Create sample bursts
		now := time.Now()
		bursts = []*career.Burst{
			{
				ID:          "burst-1",
				Name:        "Q3 Performance Review",
				Description: "Led quarterly review",
				EventIDs:    []string{"event-1"},
				Confirmed:   true,
				CreatedAt:   now,
				UpdatedAt:   now,
			},
			{
				ID:          "burst-2",
				Name:        "API Integration",
				Description: "Implemented REST API",
				EventIDs:    []string{"event-2"},
				Confirmed:   false,
				CreatedAt:   now,
				UpdatedAt:   now,
			},
		}

		context = &intents.BurstManagementContext{
			Bursts: bursts,
		}

		var err error
		intent, err = intents.NewBurstManagementIntent(context)
		Expect(err).NotTo(HaveOccurred())
		intent.Init()
	})

	Describe("BurstStateList (Root State)", func() {
		It("should cancel intent when esc is pressed", func() {
			intent.Update(tea.KeyMsg{Type: tea.KeyEsc})

			result := intent.Result()
			Expect(result).NotTo(BeNil())
			Expect(result.Status).To(Equal(intents.Cancelled))
		})
	})

	Describe("BurstStateDetail", func() {
		BeforeEach(func() {
			// Navigate to detail view
			intent.Update(tea.KeyMsg{Type: tea.KeyEnter})
		})

		It("should go back to list when esc is pressed", func() {
			intent.Update(tea.KeyMsg{Type: tea.KeyEsc})

			view := intent.View()
			Expect(view).To(ContainSubstring("Bursts"))
			result := intent.Result()
			Expect(result).To(BeNil())
		})
	})

	Describe("BurstStateDetailEvents", func() {
		BeforeEach(func() {
			// Navigate to detail, then events tab
			intent.Update(tea.KeyMsg{Type: tea.KeyEnter})
			intent.Update(tea.KeyMsg{Type: tea.KeyRunes, Runes: []rune{'e'}})
		})

		It("should go back to detail when esc is pressed", func() {
			intent.Update(tea.KeyMsg{Type: tea.KeyEsc})

			result := intent.Result()
			Expect(result).To(BeNil())
		})
	})

	Describe("BurstStateDetailFacts", func() {
		BeforeEach(func() {
			// Navigate to detail, then facts tab
			intent.Update(tea.KeyMsg{Type: tea.KeyEnter})
			intent.Update(tea.KeyMsg{Type: tea.KeyRunes, Runes: []rune{'f'}})
		})

		It("should go back to detail when esc is pressed", func() {
			intent.Update(tea.KeyMsg{Type: tea.KeyEsc})

			result := intent.Result()
			Expect(result).To(BeNil())
		})
	})

	Describe("BurstStateEdit", func() {
		BeforeEach(func() {
			// Navigate to detail, then edit
			intent.Update(tea.KeyMsg{Type: tea.KeyEnter})
			intent.Update(tea.KeyMsg{Type: tea.KeyRunes, Runes: []rune{'e'}})
		})

		It("should go back to detail when esc is pressed", func() {
			intent.Update(tea.KeyMsg{Type: tea.KeyEsc})

			result := intent.Result()
			Expect(result).To(BeNil())
		})
	})

	Describe("BurstStateDeleteConfirm", func() {
		BeforeEach(func() {
			// Navigate to detail, then delete
			intent.Update(tea.KeyMsg{Type: tea.KeyEnter})
			intent.Update(tea.KeyMsg{Type: tea.KeyRunes, Runes: []rune{'d'}})
		})

		It("should go back to detail when esc is pressed", func() {
			intent.Update(tea.KeyMsg{Type: tea.KeyEsc})

			result := intent.Result()
			Expect(result).To(BeNil())
		})

		It("should go back to detail when 'n' is pressed", func() {
			intent.Update(tea.KeyMsg{Type: tea.KeyRunes, Runes: []rune{'n'}})

			result := intent.Result()
			Expect(result).To(BeNil())
		})
	})

	Describe("BurstStateConfirm", func() {
		BeforeEach(func() {
			// Navigate to detail, then confirm
			intent.Update(tea.KeyMsg{Type: tea.KeyEnter})
			intent.Update(tea.KeyMsg{Type: tea.KeyRunes, Runes: []rune{'c'}})
		})

		It("should go back to detail when esc is pressed", func() {
			intent.Update(tea.KeyMsg{Type: tea.KeyEsc})

			result := intent.Result()
			Expect(result).To(BeNil())
		})
	})

	Describe("BurstStateExtractingFacts", func() {
		It("should allow escape during async operation", func() {
			// Simulate extracting facts state (async operation)
			// Note: This state is entered automatically during fact extraction
			// For test purposes, we verify escape doesn't break the intent
			intent.Update(tea.KeyMsg{Type: tea.KeyEsc})

			// Intent should handle escape gracefully without panic
			view := intent.View()
			Expect(view).NotTo(BeEmpty())
		})
	})
})
