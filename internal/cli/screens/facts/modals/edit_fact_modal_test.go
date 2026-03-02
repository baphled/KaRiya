package modals_test

import (
	"github.com/baphled/kariya/internal/cli/screens/facts/modals"
	"github.com/baphled/kariya/internal/domain/career"
	"github.com/baphled/kariya/internal/testutil/fixtures"
	tea "github.com/charmbracelet/bubbletea"
	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
)

var _ = Describe("EditFactModal", func() {
	var (
		modal *modals.EditFactModal
		fact  *career.Fact
	)

	BeforeEach(func() {
		fact = fixtures.FactWith("fact-1", "Test fact text")
		modal = modals.NewEditFactModal(fact)
	})

	Describe("GetTitle", func() {
		It("returns the correct title", func() {
			Expect(modal.GetTitle()).To(Equal("Edit Fact"))
		})
	})

	Describe("GetContent", func() {
		It("returns non-empty content when modal is active", func() {
			Expect(modal.GetContent()).NotTo(BeEmpty())
		})

		It("returns the form view content", func() {
			content := modal.GetContent()
			Expect(content).NotTo(BeEmpty())
		})
	})

	Describe("GetFooter", func() {
		It("returns non-empty footer instructions", func() {
			Expect(modal.GetFooter()).NotTo(BeEmpty())
		})

		It("contains navigation instructions", func() {
			footer := modal.GetFooter()
			Expect(footer).To(ContainSubstring("Enter"))
			Expect(footer).To(ContainSubstring("Esc"))
			Expect(footer).To(ContainSubstring("Tab"))
		})
	})

	Describe("IsComplete", func() {
		It("returns false when modal is not complete", func() {
			Expect(modal.IsComplete()).To(BeFalse())
		})
	})

	Describe("Result", func() {
		It("returns nil when modal is not complete", func() {
			Expect(modal.Result()).To(BeNil())
		})
	})

	Describe("View", func() {
		It("returns non-empty view when modal is active", func() {
			Expect(modal.View()).NotTo(BeEmpty())
		})

		It("contains the Edit Fact title in the rendered view", func() {
			view := modal.View()
			Expect(view).To(ContainSubstring("Edit Fact"))
		})

		It("contains modal instruction text", func() {
			view := modal.View()
			Expect(view).To(ContainSubstring("Enter"))
		})

		It("renders consistently across calls", func() {
			view1 := modal.View()
			view2 := modal.View()
			Expect(view1).To(Equal(view2))
		})
	})

	Describe("Init", func() {
		It("returns a non-nil command", func() {
			cmd := modal.Init()
			Expect(cmd).NotTo(BeNil())
		})
	})

	Describe("Update", func() {
		It("handles window size messages without panicking", func() {
			msg := tea.WindowSizeMsg{Width: 100, Height: 30}
			Expect(func() { modal.Update(msg) }).NotTo(Panic())
		})

		It("returns nil command for window size messages", func() {
			msg := tea.WindowSizeMsg{Width: 120, Height: 40}
			cmd := modal.Update(msg)
			Expect(cmd).To(BeNil())
		})

		It("does not complete the modal on window resize", func() {
			msg := tea.WindowSizeMsg{Width: 80, Height: 24}
			modal.Update(msg)
			Expect(modal.IsComplete()).To(BeFalse())
		})

		It("processes key messages through the form", func() {
			msg := tea.KeyMsg{Type: tea.KeyRunes, Runes: []rune{'a'}}
			Expect(func() { modal.Update(msg) }).NotTo(Panic())
			Expect(modal.IsComplete()).To(BeFalse())
		})

		It("still renders view after receiving key messages", func() {
			msg := tea.KeyMsg{Type: tea.KeyRunes, Runes: []rune{'x'}}
			modal.Update(msg)
			Expect(modal.View()).NotTo(BeEmpty())
		})

		Context("when escape key is pressed after initialisation", func() {
			BeforeEach(func() {
				modal.Init()
				modal.Update(tea.WindowSizeMsg{Width: 80, Height: 24})
				modal.Update(tea.KeyMsg{Type: tea.KeyEsc})
			})

			It("may trigger cancellation of the modal", func() {
				if modal.IsComplete() {
					Expect(modal.Result()).NotTo(BeNil())
					Expect(modal.Result().Accepted).To(BeFalse())
					Expect(modal.Result().HasChanges()).To(BeFalse())
				}
			})
		})

		Context("when ctrl+c is pressed after initialisation", func() {
			BeforeEach(func() {
				modal.Init()
				modal.Update(tea.WindowSizeMsg{Width: 80, Height: 24})
				modal.Update(tea.KeyMsg{Type: tea.KeyCtrlC})
			})

			It("may trigger cancellation of the modal", func() {
				if modal.IsComplete() {
					Expect(modal.Result()).NotTo(BeNil())
					Expect(modal.Result().Accepted).To(BeFalse())
				}
			})
		})
	})

	Describe("with a fully populated fact", func() {
		var richModal *modals.EditFactModal

		BeforeEach(func() {
			richFact := fixtures.FactWithCategories(
				"fact-2",
				"A detailed fact with categories",
				"event-1",
				[]string{"leadership", "engineering"},
				[]string{"hiring-manager", "tech-lead"},
			)
			richModal = modals.NewEditFactModal(richFact)
		})

		It("initialises without error", func() {
			Expect(richModal).NotTo(BeNil())
		})

		It("returns non-empty view", func() {
			Expect(richModal.View()).NotTo(BeEmpty())
		})

		It("returns correct title", func() {
			Expect(richModal.GetTitle()).To(Equal("Edit Fact"))
		})

		It("is not complete initially", func() {
			Expect(richModal.IsComplete()).To(BeFalse())
		})

		It("has no result initially", func() {
			Expect(richModal.Result()).To(BeNil())
		})

		It("returns non-empty content", func() {
			Expect(richModal.GetContent()).NotTo(BeEmpty())
		})

		It("returns footer with instructions", func() {
			Expect(richModal.GetFooter()).To(ContainSubstring("Enter"))
		})

		It("handles window resize", func() {
			msg := tea.WindowSizeMsg{Width: 200, Height: 50}
			cmd := richModal.Update(msg)
			Expect(cmd).To(BeNil())
			Expect(richModal.View()).NotTo(BeEmpty())
		})
	})
})

var _ = Describe("EditResult", func() {
	Describe("HasChanges", func() {
		It("returns false when Changes map is empty", func() {
			result := &modals.EditResult{
				Changes: map[string]interface{}{},
			}
			Expect(result.HasChanges()).To(BeFalse())
		})

		It("returns true when Changes map has entries", func() {
			result := &modals.EditResult{
				Changes: map[string]interface{}{"text": "new value"},
			}
			Expect(result.HasChanges()).To(BeTrue())
		})

		It("returns false when Changes map is nil", func() {
			result := &modals.EditResult{}
			Expect(result.HasChanges()).To(BeFalse())
		})
	})

	Describe("HasChanges with multiple entries", func() {
		It("returns true with multiple change keys", func() {
			result := &modals.EditResult{
				Changes: map[string]interface{}{
					"text":     "updated text",
					"role_fit": "staff",
				},
			}
			Expect(result.HasChanges()).To(BeTrue())
		})
	})

	Describe("Accepted field", func() {
		It("defaults to false", func() {
			result := &modals.EditResult{}
			Expect(result.Accepted).To(BeFalse())
		})

		It("can be set to true", func() {
			result := &modals.EditResult{Accepted: true}
			Expect(result.Accepted).To(BeTrue())
		})
	})
})
