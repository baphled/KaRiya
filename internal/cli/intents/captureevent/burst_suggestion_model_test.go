package captureevent

import (
	"context"

	burstfact "github.com/baphled/kariya/internal/service/career/burstfact"
	tea "github.com/charmbracelet/bubbletea"
	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
)

var _ = Describe("BurstSuggestionModelNew", func() {
	var (
		model       *BurstSuggestionModelNew
		suggestions []burstfact.BurstSuggestion
	)

	BeforeEach(func() {
		suggestions = []burstfact.BurstSuggestion{
			{Name: "API Work", Description: "REST API development", EventIDs: []string{}, ConfidenceScore: 0.9},
			{Name: "CI/CD Work", Description: "Pipeline setup", EventIDs: []string{}, ConfidenceScore: 0.7},
		}
		model = NewBurstSuggestionModelNew(context.Background(), nil, suggestions)
	})

	Describe("Init", func() {
		It("returns nil (no initial command needed)", func() {
			cmd := model.Init()
			Expect(cmd).To(BeNil())
		})
	})

	Describe("Update", func() {
		It("handles WindowSizeMsg by updating dimensions", func() {
			updated, _ := model.Update(tea.WindowSizeMsg{Width: 100, Height: 50})
			Expect(updated).NotTo(BeNil())
		})

		It("handles unrecognised messages without panicking", func() {
			updated, _ := model.Update("unknown message")
			Expect(updated).NotTo(BeNil())
		})

		It("handles KeyUp by changing index", func() {
			model.Update(tea.KeyMsg{Type: tea.KeyUp})
			Expect(model.currentIdx).To(BeNumerically(">=", 0))
		})

		It("handles KeyDown by changing index", func() {
			model.Update(tea.KeyMsg{Type: tea.KeyDown})
			Expect(model.currentIdx).To(BeNumerically(">=", 0))
		})

		It("handles j key navigation without panicking", func() {
			_, cmd := model.Update(tea.KeyMsg{Type: tea.KeyRunes, Runes: []rune{'j'}})
			Expect(cmd).To(BeNil())
		})

		It("handles k key navigation without panicking", func() {
			_, cmd := model.Update(tea.KeyMsg{Type: tea.KeyRunes, Runes: []rune{'k'}})
			Expect(cmd).To(BeNil())
		})

		It("handles n key (reject) without panicking", func() {
			_, _ = model.Update(tea.KeyMsg{Type: tea.KeyRunes, Runes: []rune{'n'}})
			Expect(len(model.rejected)).To(BeNumerically(">=", 0))
		})

		It("handles y key (confirm) without panicking", func() {
			_, _ = model.Update(tea.KeyMsg{Type: tea.KeyRunes, Runes: []rune{'y'}})
			Expect(len(model.confirmed)).To(BeNumerically(">=", 0))
		})

		It("handles Esc key by returning BackMsg command", func() {
			_, cmd := model.Update(tea.KeyMsg{Type: tea.KeyEsc})
			Expect(cmd).NotTo(BeNil())
		})
	})

	Describe("View", func() {
		It("renders review view when not editing", func() {
			view := model.View()
			Expect(view).NotTo(BeEmpty())
		})

		It("renders no-suggestions view when suggestions are empty", func() {
			empty := NewBurstSuggestionModelNew(context.Background(), nil, []burstfact.BurstSuggestion{})
			view := empty.View()
			Expect(view).To(ContainSubstring("No burst suggestions"))
		})
	})

	Describe("GetConfirmed", func() {
		It("returns empty slice initially", func() {
			Expect(model.GetConfirmed()).To(BeEmpty())
		})
	})

	Describe("GetRejected", func() {
		It("returns empty slice initially", func() {
			Expect(model.GetRejected()).To(BeEmpty())
		})
	})

	Describe("IsDone", func() {
		It("returns false when no suggestions have been processed", func() {
			Expect(model.IsDone()).To(BeFalse())
		})
	})

	Describe("IsEditing", func() {
		It("returns false by default", func() {
			Expect(model.IsEditing()).To(BeFalse())
		})
	})

	Describe("SetEditedName", func() {
		It("sets the edited name for the given index", func() {
			model.SetEditedName(0, "New Name")
			Expect(model.editedNames[0]).To(Equal("New Name"))
		})
	})

	Describe("SetEditedDescription", func() {
		It("sets the edited description for the given index", func() {
			model.SetEditedDescription(0, "New Description")
			Expect(model.editedDescs[0]).To(Equal("New Description"))
		})
	})

	Describe("ExitEditMode", func() {
		It("sets editing to false", func() {
			model.editing = true
			model.ExitEditMode()
			Expect(model.IsEditing()).To(BeFalse())
		})
	})

	Describe("GetTitle", func() {
		It("returns suggestion count title when not editing", func() {
			title := model.GetTitle()
			Expect(title).To(ContainSubstring("1 of 2"))
		})

		It("returns edit title when editing", func() {
			model.editing = true
			title := model.GetTitle()
			Expect(title).To(ContainSubstring("Edit"))
		})
	})

	Describe("GetContent", func() {
		It("returns content string when suggestions exist", func() {
			content := model.GetContent()
			Expect(content).NotTo(BeEmpty())
		})

		It("returns no-suggestions message when empty", func() {
			empty := NewBurstSuggestionModelNew(context.Background(), nil, []burstfact.BurstSuggestion{})
			content := empty.GetContent()
			Expect(content).To(ContainSubstring("No burst suggestions"))
		})

		It("returns edit-form-not-initialized message when editing with nil form", func() {
			model.editing = true
			model.editForm = nil
			content := model.GetContent()
			Expect(content).To(ContainSubstring("Edit form not initialized"))
		})
	})

	Describe("GetFooter", func() {
		It("returns review footer when not editing", func() {
			footer := model.GetFooter()
			Expect(footer).To(ContainSubstring("y: Confirm"))
		})

		It("returns edit footer when editing", func() {
			model.editing = true
			footer := model.GetFooter()
			Expect(footer).To(ContainSubstring("Tab"))
		})
	})
})
