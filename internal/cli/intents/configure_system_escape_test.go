package intents

import (
	tea "github.com/charmbracelet/bubbletea"
	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
)

var _ = Describe("ConfigureSystem - Escape Key Behavior", func() {
	var model *ConfigureSystemModel

	BeforeEach(func() {
		ctx := NewConfigureSystemContext()
		model = NewConfigureSystemModel(ctx)
		model.Init()
	})

	Describe("SelectDomain State", func() {
		BeforeEach(func() {
			model.state = ConfigStateSelectDomain
		})

		It("should cancel intent when esc is pressed", func() {
			model.Update(tea.KeyMsg{Type: tea.KeyEsc})

			Expect(model.active).To(BeFalse())
			Expect(model.result).NotTo(BeNil())
			Expect(model.result.Success).To(BeFalse())
		})

		It("should cancel intent when 'm' is pressed", func() {
			model.Update(tea.KeyMsg{Type: tea.KeyRunes, Runes: []rune{'m'}})

			Expect(model.active).To(BeFalse())
			Expect(model.result).NotTo(BeNil())
			Expect(model.result.Success).To(BeFalse())
		})
	})

	Describe("EditSettings State", func() {
		BeforeEach(func() {
			model.state = ConfigStateEditSettings
			model.domain = DomainSystem
			model.changes = &ConfigurationChanges{
				Domain:   DomainSystem,
				Original: make(map[string]interface{}),
				Modified: make(map[string]interface{}),
			}
		})

		It("should go back to SelectDomain when esc is pressed", func() {
			model.Update(tea.KeyMsg{Type: tea.KeyEsc})

			Expect(model.state).To(Equal(ConfigStateSelectDomain))
			Expect(model.active).To(BeTrue())
		})

		It("should return to main menu when 'm' is pressed", func() {
			model.Update(tea.KeyMsg{Type: tea.KeyRunes, Runes: []rune{'m'}})

			Expect(model.active).To(BeFalse())
			Expect(model.result).NotTo(BeNil())
			Expect(model.result.Success).To(BeFalse())
		})
	})

	Describe("ReviewChanges State", func() {
		BeforeEach(func() {
			model.state = ConfigStateReviewChanges
			model.domain = DomainSystem
			model.changes = &ConfigurationChanges{
				Domain:   DomainSystem,
				Original: make(map[string]interface{}),
				Modified: make(map[string]interface{}),
			}
		})

		It("should go back to EditSettings when esc is pressed", func() {
			model.Update(tea.KeyMsg{Type: tea.KeyEsc})

			Expect(model.state).To(Equal(ConfigStateEditSettings))
			Expect(model.active).To(BeTrue())
		})

		It("should return to main menu when 'm' is pressed", func() {
			model.Update(tea.KeyMsg{Type: tea.KeyRunes, Runes: []rune{'m'}})

			Expect(model.active).To(BeFalse())
			Expect(model.result).NotTo(BeNil())
			Expect(model.result.Success).To(BeFalse())
		})
	})

	Describe("Confirm State", func() {
		BeforeEach(func() {
			model.state = ConfigStateConfirm
			model.domain = DomainSystem
			model.changes = &ConfigurationChanges{
				Domain:   DomainSystem,
				Original: make(map[string]interface{}),
				Modified: make(map[string]interface{}),
			}
		})

		It("should go back to ReviewChanges when esc is pressed", func() {
			model.Update(tea.KeyMsg{Type: tea.KeyEsc})

			Expect(model.state).To(Equal(ConfigStateReviewChanges))
			Expect(model.active).To(BeTrue())
		})

		It("should return to main menu when 'm' is pressed", func() {
			model.Update(tea.KeyMsg{Type: tea.KeyRunes, Runes: []rune{'m'}})

			Expect(model.active).To(BeFalse())
			Expect(model.result).NotTo(BeNil())
			Expect(model.result.Success).To(BeFalse())
		})
	})

	Describe("Saving State", func() {
		BeforeEach(func() {
			model.state = ConfigStateSaving
			model.domain = DomainSystem
			model.changes = &ConfigurationChanges{
				Domain:   DomainSystem,
				Original: make(map[string]interface{}),
				Modified: make(map[string]interface{}),
			}
		})

		It("should go back to ReviewChanges when esc is pressed", func() {
			model.Update(tea.KeyMsg{Type: tea.KeyEsc})

			Expect(model.state).To(Equal(ConfigStateReviewChanges))
			Expect(model.active).To(BeTrue())
		})

		It("should return to main menu when 'm' is pressed", func() {
			model.Update(tea.KeyMsg{Type: tea.KeyRunes, Runes: []rune{'m'}})

			Expect(model.active).To(BeFalse())
			Expect(model.result).NotTo(BeNil())
			Expect(model.result.Success).To(BeFalse())
		})
	})

	Describe("Complete State", func() {
		BeforeEach(func() {
			model.state = ConfigStateComplete
			model.result = &ConfigureSystemResult{
				Success: true,
				Domain:  DomainSystem,
			}
		})

		It("should deactivate when esc is pressed", func() {
			model.Update(tea.KeyMsg{Type: tea.KeyEsc})

			Expect(model.active).To(BeFalse())
		})

		It("should deactivate when 'm' is pressed", func() {
			model.Update(tea.KeyMsg{Type: tea.KeyRunes, Runes: []rune{'m'}})

			Expect(model.active).To(BeFalse())
		})
	})

	Describe("Failed State", func() {
		BeforeEach(func() {
			model.state = ConfigStateFailed
			model.error = &IntentError{
				Code:    "TEST_ERROR",
				Message: "Configuration failed",
			}
		})

		It("should deactivate when esc is pressed", func() {
			model.Update(tea.KeyMsg{Type: tea.KeyEsc})

			Expect(model.active).To(BeFalse())
		})

		It("should deactivate when 'm' is pressed", func() {
			model.Update(tea.KeyMsg{Type: tea.KeyRunes, Runes: []rune{'m'}})

			Expect(model.active).To(BeFalse())
		})
	})

	// Note: View footer tests removed - footer text is now handled by
	// StandardView via getContextHelp() in the intent, not the model.
	// The model's View() methods now return pure content only.
})
