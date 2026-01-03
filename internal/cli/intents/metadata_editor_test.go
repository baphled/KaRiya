package intents_test

import (
	"context"
	"testing"

	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"

	"github.com/baphled/kariya/internal/cli/intents"
)

func TestMetadataEditor(t *testing.T) {
	RegisterFailHandler(Fail)
	RunSpecs(t, "MetadataEditor Intent Suite")
}

var _ = Describe("MetadataEditor Intent", func() {
	var (
		model *intents.MetadataEditorModel
		ctx context.Context
	)

	BeforeEach(func() {
		ctx = context.Background()
		data := intents.NewMetadataEditorContext(ctx)
		model = intents.NewMetadataEditorIntent(data)
	})

	Describe("Initialization", func() {
		It("should initialize with review state", func() {
			cmd := model.Init(ctx)
			Expect(cmd).To(BeNil())
			Expect(model.View()).NotTo(BeEmpty())
		})
	})

	Describe("Review State", func() {
		It("should render review view", func() {
			model.Init(ctx)
			view := model.View()
			Expect(view).To(ContainSubstring("Metadata Review"))
		})
	})

	Describe("Context Operations", func() {
		It("should load metadata", func() {
			data := intents.NewMetadataEditorContext(ctx)
			metadata := map[string]interface{}{
				"name": "Test",
				"value": 42,
			}
			data.LoadMetadata(metadata)
			Expect(data.OriginalMetadata).To(HaveLen(2))
		})

		It("should track form errors", func() {
			data := intents.NewMetadataEditorContext(ctx)
			data.SetFormError("name", "Name is required")
			Expect(data.HasFormErrors()).To(BeTrue())
		})

		It("should clear form errors", func() {
			data := intents.NewMetadataEditorContext(ctx)
			data.SetFormError("name", "Name is required")
			data.ClearFormErrors()
			Expect(data.HasFormErrors()).To(BeFalse())
		})

		It("should set field values", func() {
			data := intents.NewMetadataEditorContext(ctx)
			data.SetFieldValue("name", "Updated")
			Expect(data.ChangedFields["name"]).To(BeTrue())
		})

		It("should detect changes", func() {
			data := intents.NewMetadataEditorContext(ctx)
			Expect(data.HasChanges()).To(BeFalse())
			data.SetFieldValue("name", "Updated")
			Expect(data.HasChanges()).To(BeTrue())
		})

		It("should get changes", func() {
			data := intents.NewMetadataEditorContext(ctx)
			data.SetFieldValue("name", "Updated")
			data.SetFieldValue("value", 100)
			changes := data.GetChanges()
			Expect(changes).To(HaveLen(2))
		})

		It("should reset changes", func() {
			data := intents.NewMetadataEditorContext(ctx)
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
			model.Init(ctx)
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
