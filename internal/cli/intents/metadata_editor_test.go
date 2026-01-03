package intents

import (
	"context"

	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
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
			Expect(cmd).To(BeNil())
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

	Describe("Result Handling", func() {
		It("should return result", func() {
			result := model.Result()
			Expect(result).NotTo(BeNil())
		})
	})
})
