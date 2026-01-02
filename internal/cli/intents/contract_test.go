package intents

import (
	"testing"

	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
)

func TestContract(t *testing.T) {
	RegisterFailHandler(Fail)
	RunSpecs(t, "Intent Contract Suite")
}

var _ = Describe("ModalEditResult", func() {
	Context("when creating a new modal edit result", func() {
		It("should store original and modified values", func() {
			original := "original value"
			modified := "modified value"
			changes := map[string]interface{}{"field": "new value"}

			result := NewModalEditResult(original, modified, true, changes)

			Expect(result.Original).To(Equal(original))
			Expect(result.Modified).To(Equal(modified))
			Expect(result.Accepted).To(BeTrue())
			Expect(result.Changes).To(Equal(changes))
		})

		It("should handle nil changes map", func() {
			result := NewModalEditResult("original", "modified", true, nil)

			Expect(result.Changes).NotTo(BeNil())
			Expect(result.Changes).To(BeEmpty())
		})
	})

	Context("HasChanges", func() {
		It("should return true when changes exist", func() {
			result := NewModalEditResult("original", "modified", true, map[string]interface{}{"field": "value"})

			Expect(result.HasChanges()).To(BeTrue())
		})

		It("should return false when no changes", func() {
			result := NewModalEditResult("original", "original", true, make(map[string]interface{}))

			Expect(result.HasChanges()).To(BeFalse())
		})
	})

	Context("WasAccepted", func() {
		It("should return true when accepted is true", func() {
			result := NewModalEditResult("original", "modified", true, nil)

			Expect(result.WasAccepted()).To(BeTrue())
		})

		It("should return false when accepted is false", func() {
			result := NewModalEditResult("original", "original", false, nil)

			Expect(result.WasAccepted()).To(BeFalse())
		})
	})

	Context("GetChange", func() {
		It("should return the value for an existing change", func() {
			changes := map[string]interface{}{"field": "new value"}
			result := NewModalEditResult("original", "modified", true, changes)

			Expect(result.GetChange("field")).To(Equal("new value"))
		})

		It("should return nil for a non-existent change", func() {
			result := NewModalEditResult("original", "modified", true, map[string]interface{}{"field": "value"})

			Expect(result.GetChange("nonexistent")).To(BeNil())
		})

		It("should return nil when changes is nil", func() {
			result := &ModalEditResult[string]{
				Original: "original",
				Modified: "modified",
				Accepted: true,
				Changes:  nil,
			}

			Expect(result.GetChange("field")).To(BeNil())
		})
	})

	Context("NewCancelledModalEditResult", func() {
		It("should create a result with accepted=false and no changes", func() {
			original := "original value"

			result := NewCancelledModalEditResult(original)

			Expect(result.Original).To(Equal(original))
			Expect(result.Modified).To(Equal(original))
			Expect(result.Accepted).To(BeFalse())
			Expect(result.HasChanges()).To(BeFalse())
		})
	})

	Context("with complex types", func() {
		type TestData struct {
			Name string
			Age  int
		}

		It("should work with struct types", func() {
			original := TestData{Name: "John", Age: 30}
			modified := TestData{Name: "Jane", Age: 30}
			changes := map[string]interface{}{"Name": "Jane"}

			result := NewModalEditResult(original, modified, true, changes)

			Expect(result.Original).To(Equal(original))
			Expect(result.Modified).To(Equal(modified))
			Expect(result.HasChanges()).To(BeTrue())
		})
	})
})
