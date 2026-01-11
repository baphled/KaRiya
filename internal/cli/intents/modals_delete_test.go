package intents_test

import (
	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"

	"github.com/baphled/kariya/internal/cli/intents"
	tea "github.com/charmbracelet/bubbletea"
)

var _ = Describe("DeleteConfirmationModal", func() {
	Describe("NewDeleteConfirmationModal", func() {
		It("should create a new modal with item details", func() {
			modal := intents.NewDeleteConfirmationModal(
				"event",
				"Test Event",
				"This is a test event that will be deleted",
			)
			Expect(modal).NotTo(BeNil())
		})

		It("should handle empty description", func() {
			modal := intents.NewDeleteConfirmationModal(
				"event",
				"Test Event",
				"",
			)
			Expect(modal).NotTo(BeNil())
		})

		It("should handle long titles", func() {
			longTitle := "This is a very long title that might need to be handled carefully in the UI layout to ensure it doesn't break the modal design"
			modal := intents.NewDeleteConfirmationModal(
				"event",
				longTitle,
				"Description",
			)
			Expect(modal).NotTo(BeNil())
		})
	})

	Describe("Update", func() {
		var modal *intents.DeleteConfirmationModal

		BeforeEach(func() {
			modal = intents.NewDeleteConfirmationModal(
				"event",
				"Test Event",
				"This is a test event",
			)
		})

		Describe("Window resize", func() {
			It("should update dimensions on window resize", func() {
				cmd := modal.Update(tea.WindowSizeMsg{Width: 100, Height: 40})
				Expect(cmd).To(BeNil())
			})

			It("should handle small terminal size", func() {
				cmd := modal.Update(tea.WindowSizeMsg{Width: 80, Height: 24})
				Expect(cmd).To(BeNil())
			})

			It("should handle large terminal size", func() {
				cmd := modal.Update(tea.WindowSizeMsg{Width: 200, Height: 60})
				Expect(cmd).To(BeNil())
			})
		})

		Describe("Key input - confirmation", func() {
			It("should confirm on 'y' key", func() {
				cmd := modal.Update(tea.KeyMsg{Type: tea.KeyRunes, Runes: []rune{'y'}})
				Expect(cmd).To(BeNil())
				result := modal.Result()
				Expect(result).NotTo(BeNil())
				Expect(result.Confirmed).To(BeTrue())
			})

			It("should confirm on 'Y' key", func() {
				cmd := modal.Update(tea.KeyMsg{Type: tea.KeyRunes, Runes: []rune{'Y'}})
				Expect(cmd).To(BeNil())
				result := modal.Result()
				Expect(result).NotTo(BeNil())
				Expect(result.Confirmed).To(BeTrue())
			})

			It("should confirm on Enter key", func() {
				cmd := modal.Update(tea.KeyMsg{Type: tea.KeyEnter})
				Expect(cmd).To(BeNil())
				result := modal.Result()
				Expect(result).NotTo(BeNil())
				Expect(result.Confirmed).To(BeTrue())
			})
		})

		Describe("Key input - cancellation", func() {
			It("should cancel on 'n' key", func() {
				cmd := modal.Update(tea.KeyMsg{Type: tea.KeyRunes, Runes: []rune{'n'}})
				Expect(cmd).To(BeNil())
				result := modal.Result()
				Expect(result).NotTo(BeNil())
				Expect(result.Confirmed).To(BeFalse())
			})

			It("should cancel on 'N' key", func() {
				cmd := modal.Update(tea.KeyMsg{Type: tea.KeyRunes, Runes: []rune{'N'}})
				Expect(cmd).To(BeNil())
				result := modal.Result()
				Expect(result).NotTo(BeNil())
				Expect(result.Confirmed).To(BeFalse())
			})

			It("should cancel on Escape key", func() {
				cmd := modal.Update(tea.KeyMsg{Type: tea.KeyEsc})
				Expect(cmd).To(BeNil())
				result := modal.Result()
				Expect(result).NotTo(BeNil())
				Expect(result.Confirmed).To(BeFalse())
			})
		})

		Describe("Key input - ignored keys", func() {
			It("should ignore other keys and remain open", func() {
				cmd := modal.Update(tea.KeyMsg{Type: tea.KeyRunes, Runes: []rune{'x'}})
				Expect(cmd).To(BeNil())
				result := modal.Result()
				Expect(result).To(BeNil()) // Not closed yet
			})

			It("should ignore arrow keys", func() {
				cmd := modal.Update(tea.KeyMsg{Type: tea.KeyUp})
				Expect(cmd).To(BeNil())
				result := modal.Result()
				Expect(result).To(BeNil())
			})
		})
	})

	Describe("View", func() {
		var modal *intents.DeleteConfirmationModal

		BeforeEach(func() {
			modal = intents.NewDeleteConfirmationModal(
				"event",
				"Test Event",
				"This is a test event description",
			)
		})

		It("should render warning header", func() {
			view := modal.View()
			Expect(view).To(ContainSubstring("Delete"))
			Expect(view).To(ContainSubstring("event"))
		})

		It("should render item title", func() {
			view := modal.View()
			Expect(view).To(ContainSubstring("Test Event"))
		})

		It("should render description", func() {
			view := modal.View()
			Expect(view).To(ContainSubstring("This is a test event description"))
		})

		It("should render warning message", func() {
			view := modal.View()
			Expect(view).To(ContainSubstring("cannot be undone"))
		})

		It("should render confirmation prompt", func() {
			view := modal.View()
			Expect(view).To(ContainSubstring("y"))
			Expect(view).To(ContainSubstring("n"))
		})

		Describe("Different item types", func() {
			It("should show 'burst' in warning for bursts", func() {
				modal := intents.NewDeleteConfirmationModal(
					"burst",
					"Test Burst",
					"Burst description",
				)
				view := modal.View()
				Expect(view).To(ContainSubstring("burst"))
			})

			It("should show 'fact' in warning for facts", func() {
				modal := intents.NewDeleteConfirmationModal(
					"fact",
					"Test Fact",
					"Fact description",
				)
				view := modal.View()
				Expect(view).To(ContainSubstring("fact"))
			})
		})

		Describe("Long text handling", func() {
			It("should handle long descriptions", func() {
				longDesc := "This is a very long description that contains a lot of text and might need to be wrapped or handled carefully to ensure the modal remains readable and doesn't break the layout of the terminal interface when displayed to the user"
				modal := intents.NewDeleteConfirmationModal(
					"event",
					"Test Event",
					longDesc,
				)
				view := modal.View()
				Expect(view).NotTo(BeEmpty())
			})

			It("should handle multiline descriptions", func() {
				multilineDesc := "Line 1\nLine 2\nLine 3\nLine 4"
				modal := intents.NewDeleteConfirmationModal(
					"event",
					"Test Event",
					multilineDesc,
				)
				view := modal.View()
				Expect(view).To(ContainSubstring("Line 1"))
			})
		})
	})

	Describe("Result", func() {
		var modal *intents.DeleteConfirmationModal

		BeforeEach(func() {
			modal = intents.NewDeleteConfirmationModal(
				"event",
				"Test Event",
				"Description",
			)
		})

		It("should return nil before user decision", func() {
			result := modal.Result()
			Expect(result).To(BeNil())
		})

		It("should return confirmed result after 'y'", func() {
			modal.Update(tea.KeyMsg{Type: tea.KeyRunes, Runes: []rune{'y'}})
			result := modal.Result()
			Expect(result).NotTo(BeNil())
			Expect(result.Confirmed).To(BeTrue())
		})

		It("should return cancelled result after 'n'", func() {
			modal.Update(tea.KeyMsg{Type: tea.KeyRunes, Runes: []rune{'n'}})
			result := modal.Result()
			Expect(result).NotTo(BeNil())
			Expect(result.Confirmed).To(BeFalse())
		})

		It("should return same result on multiple Result() calls", func() {
			modal.Update(tea.KeyMsg{Type: tea.KeyRunes, Runes: []rune{'y'}})
			result1 := modal.Result()
			result2 := modal.Result()
			Expect(result1).To(Equal(result2))
		})
	})

	Describe("IsComplete", func() {
		var modal *intents.DeleteConfirmationModal

		BeforeEach(func() {
			modal = intents.NewDeleteConfirmationModal(
				"event",
				"Test Event",
				"Description",
			)
		})

		It("should return false before user decision", func() {
			Expect(modal.IsComplete()).To(BeFalse())
		})

		It("should return true after confirmation", func() {
			modal.Update(tea.KeyMsg{Type: tea.KeyRunes, Runes: []rune{'y'}})
			Expect(modal.IsComplete()).To(BeTrue())
		})

		It("should return true after cancellation", func() {
			modal.Update(tea.KeyMsg{Type: tea.KeyRunes, Runes: []rune{'n'}})
			Expect(modal.IsComplete()).To(BeTrue())
		})
	})

	Describe("WasConfirmed", func() {
		var modal *intents.DeleteConfirmationModal

		BeforeEach(func() {
			modal = intents.NewDeleteConfirmationModal(
				"event",
				"Test Event",
				"Description",
			)
		})

		It("should return false before user decision", func() {
			Expect(modal.WasConfirmed()).To(BeFalse())
		})

		It("should return true after confirmation", func() {
			modal.Update(tea.KeyMsg{Type: tea.KeyRunes, Runes: []rune{'y'}})
			Expect(modal.WasConfirmed()).To(BeTrue())
		})

		It("should return false after cancellation", func() {
			modal.Update(tea.KeyMsg{Type: tea.KeyRunes, Runes: []rune{'n'}})
			Expect(modal.WasConfirmed()).To(BeFalse())
		})
	})
})
