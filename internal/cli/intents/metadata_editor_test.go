package intents

import (
	"context"

	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"

	tea "github.com/charmbracelet/bubbletea"
)

var _ = Describe("MetadataEditor Intent", func() {
	var (
		model *MetadataEditorModel
		ctx   context.Context
	)

	BeforeEach(func() {
		ctx = context.Background()
		data := NewMetadataEditorContext(ctx)
		model = NewMetadataEditorIntent(data)
	})

	Describe("Initialization", func() {
		It("should initialize with review state", func() {
			cmd := model.Init()
			Expect(cmd).NotTo(BeNil())
			Expect(model.View()).NotTo(BeEmpty())
		})
	})

	Describe("Review State", func() {
		It("should render review view", func() {
			model.Init()
			view := model.View()
			Expect(view).To(ContainSubstring("Metadata Review"))
		})
	})

	Describe("Context Operations", func() {
		It("should load metadata", func() {
			data := NewMetadataEditorContext(ctx)
			metadata := map[string]interface{}{
				"name":  "Test",
				"value": 42,
			}
			data.LoadMetadata(metadata)
			Expect(data.OriginalMetadata).To(HaveLen(2))
		})

		It("should track form errors", func() {
			data := NewMetadataEditorContext(ctx)
			data.SetFormError("name", "Name is required")
			Expect(data.HasFormErrors()).To(BeTrue())
		})

		It("should clear form errors", func() {
			data := NewMetadataEditorContext(ctx)
			data.SetFormError("name", "Name is required")
			data.ClearFormErrors()
			Expect(data.HasFormErrors()).To(BeFalse())
		})

		It("should set field values", func() {
			data := NewMetadataEditorContext(ctx)
			data.SetFieldValue("name", "Updated")
			Expect(data.ChangedFields["name"]).To(BeTrue())
		})

		It("should detect changes", func() {
			data := NewMetadataEditorContext(ctx)
			Expect(data.HasChanges()).To(BeFalse())
			data.SetFieldValue("name", "Updated")
			Expect(data.HasChanges()).To(BeTrue())
		})

		It("should get changes", func() {
			data := NewMetadataEditorContext(ctx)
			data.SetFieldValue("name", "Updated")
			data.SetFieldValue("value", 100)
			changes := data.GetChanges()
			Expect(changes).To(HaveLen(2))
		})

		It("should reset changes", func() {
			data := NewMetadataEditorContext(ctx)
			metadata := map[string]interface{}{
				"name": "Test",
			}
			data.LoadMetadata(metadata)
			data.SetFieldValue("name", "Updated")
			Expect(data.HasChanges()).To(BeTrue())
			data.ResetChanges()
			Expect(data.HasChanges()).To(BeFalse())
		})
	})

	Describe("View Methods", func() {
		It("should render view without error", func() {
			model.Init()
			view := model.View()
			Expect(view).NotTo(BeEmpty())
		})
	})

	Describe("Navigation (Key-Driven State Transitions)", func() {
		var (
			editModel *MetadataEditorModel
		)

		BeforeEach(func() {
			ctx := context.Background()
			data := NewMetadataEditorContext(ctx)
			data.EntityType = "event"
			data.EntityID = "evt-456"
			data.LoadMetadata(map[string]interface{}{"name": "initial"})
			editModel = NewMetadataEditorIntent(data)
			editModel.Init()
		})

		It("should transition from review to edit on 'e'", func() {
			msg := tea.KeyMsg{Type: tea.KeyRunes, Runes: []rune{'e'}}
			editModel.Update(msg)
			Expect(editModel.data.CurrentState).To(Equal(MetadataEditState))
		})

		It("should quit/cancel on 'q' from review", func() {
			msg := tea.KeyMsg{Type: tea.KeyRunes, Runes: []rune{'q'}}
			editModel.Update(msg)
			Expect(editModel.result.Status).To(Equal(Cancelled))
		})

		It("should quit/cancel on 'esc' from review", func() {
			msg := tea.KeyMsg{Type: tea.KeyEsc}
			editModel.Update(msg)
			Expect(editModel.result.Status).To(Equal(Cancelled))
		})

		It("should transition from edit to confirm on 'ctrl+s' when changes exist", func() {
			editModel.data.CurrentState = MetadataEditState
			editModel.data.SetFieldValue("name", "newval")
			msg := tea.KeyMsg{Type: tea.KeyCtrlS}
			editModel.Update(msg)
			Expect(editModel.data.CurrentState).To(Equal(MetadataConfirmState))
		})

		It("should return to review from edit on 'esc' and clear changes", func() {
			editModel.data.CurrentState = MetadataEditState
			editModel.data.SetFieldValue("name", "changed")
			msg := tea.KeyMsg{Type: tea.KeyEsc}
			editModel.Update(msg)
			Expect(editModel.data.CurrentState).To(Equal(MetadataReviewState))
			Expect(editModel.data.HasChanges()).To(BeFalse())
		})

		It("should go to completed on 'y' from confirm", func() {
			editModel.data.CurrentState = MetadataConfirmState
			editModel.data.ChangedFields["name"] = true
			msg := tea.KeyMsg{Type: tea.KeyRunes, Runes: []rune{'y'}}
			editModel.Update(msg)
			Expect(editModel.result.Status).To(Equal(Completed))
			Expect(editModel.result.Data.Action).To(Equal("saved"))
		})

		It("should return to edit on 'n' from confirm", func() {
			editModel.data.CurrentState = MetadataConfirmState
			msg := tea.KeyMsg{Type: tea.KeyRunes, Runes: []rune{'n'}}
			editModel.Update(msg)
			Expect(editModel.data.CurrentState).To(Equal(MetadataEditState))
		})

		It("should return to edit on 'esc' from confirm", func() {
			editModel.data.CurrentState = MetadataConfirmState
			msg := tea.KeyMsg{Type: tea.KeyEsc}
			editModel.Update(msg)
			Expect(editModel.data.CurrentState).To(Equal(MetadataEditState))
		})
	})

	Describe("Result Handling", func() {
		It("should return result", func() {
			result := model.Result()
			Expect(result).NotTo(BeNil())
		})
	})
})
