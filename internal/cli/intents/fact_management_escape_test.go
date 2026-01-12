package intents_test

import (
	"time"

	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"

	"github.com/baphled/kariya/internal/cli/intents"
	"github.com/baphled/kariya/internal/domain/career"
	tea "github.com/charmbracelet/bubbletea"
)

var _ = Describe("FactManagement - Escape Key Behavior", func() {
	var (
		context *intents.FactManagementContext
		intent  *intents.FactManagementModel
		facts   []*career.Fact
	)

	BeforeEach(func() {
		now := time.Now()
		facts = []*career.Fact{
			{
				ID:        "fact-1",
				Text:      "Improved API performance by 40%",
				CreatedAt: now,
				UpdatedAt: now,
			},
			{
				ID:        "fact-2",
				Text:      "Led team of 5 engineers",
				CreatedAt: now,
				UpdatedAt: now,
			},
		}

		context = &intents.FactManagementContext{
			Facts: facts,
		}

		intent = intents.NewFactManagementIntent(context)
		intent.Init()
	})

	Describe("FactListState (Root State)", func() {
		It("should cancel intent when esc is pressed", func() {
			intent.Update(tea.KeyMsg{Type: tea.KeyEsc})

			result := intent.Result()
			Expect(result).NotTo(BeNil())
			Expect(result.Status).To(Equal(intents.Cancelled))
		})
	})

	Describe("FactViewState", func() {
		BeforeEach(func() {
			intent.Update(tea.KeyMsg{Type: tea.KeyEnter})
		})

		It("should go back to list when esc is pressed", func() {
			intent.Update(tea.KeyMsg{Type: tea.KeyEsc})

			result := intent.Result()
			Expect(result).To(BeNil())
		})
	})

	Describe("FactEditorState", func() {
		BeforeEach(func() {
			intent.Update(tea.KeyMsg{Type: tea.KeyEnter})
			intent.Update(tea.KeyMsg{Type: tea.KeyRunes, Runes: []rune{'e'}})
		})

		It("should go back to view when esc is pressed", func() {
			intent.Update(tea.KeyMsg{Type: tea.KeyEsc})

			result := intent.Result()
			Expect(result).To(BeNil())
		})
	})

	Describe("FactDeleteConfirmState", func() {
		BeforeEach(func() {
			intent.Update(tea.KeyMsg{Type: tea.KeyEnter})
			intent.Update(tea.KeyMsg{Type: tea.KeyRunes, Runes: []rune{'d'}})
		})

		It("should go back to view when esc is pressed", func() {
			intent.Update(tea.KeyMsg{Type: tea.KeyEsc})

			result := intent.Result()
			Expect(result).To(BeNil())
		})

		It("should go back to view when 'n' is pressed", func() {
			intent.Update(tea.KeyMsg{Type: tea.KeyRunes, Runes: []rune{'n'}})

			result := intent.Result()
			Expect(result).To(BeNil())
		})
	})

	Describe("FactResultsState", func() {
		It("should handle escape appropriately", func() {
			intent.Update(tea.KeyMsg{Type: tea.KeyEsc})

			view := intent.View()
			Expect(view).NotTo(BeEmpty())
		})
	})

	Describe("FactCompletedState", func() {
		It("should handle escape in completed state", func() {
			intent.Update(tea.KeyMsg{Type: tea.KeyEsc})

			view := intent.View()
			Expect(view).NotTo(BeEmpty())
		})
	})
})
