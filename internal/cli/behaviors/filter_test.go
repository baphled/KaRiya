package behaviors_test

import (
	"github.com/baphled/kariya/internal/cli/behaviors"
	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
)

var _ = Describe("Filter Functions", func() {

	Describe("SearchableText", func() {
		It("should match when query is empty", func() {
			result := behaviors.SearchableText("Any text", "")
			Expect(result).To(BeTrue())
		})

		It("should match case-insensitive substring", func() {
			result := behaviors.SearchableText("Backend Developer", "backend")
			Expect(result).To(BeTrue())
		})

		It("should not match when substring not found", func() {
			result := behaviors.SearchableText("Backend Developer", "frontend")
			Expect(result).To(BeFalse())
		})

		It("should handle uppercase query", func() {
			result := behaviors.SearchableText("backend developer", "BACKEND")
			Expect(result).To(BeTrue())
		})
	})

	Describe("SearchableFields", func() {
		It("should match when query is empty", func() {
			result := behaviors.SearchableFields("", "field1", "field2")
			Expect(result).To(BeTrue())
		})

		It("should match when any field contains query", func() {
			result := behaviors.SearchableFields("backend", "Backend Developer", "TechCorp", "Engineering")
			Expect(result).To(BeTrue())
		})

		It("should not match when no field contains query", func() {
			result := behaviors.SearchableFields("frontend", "Backend Developer", "TechCorp", "Engineering")
			Expect(result).To(BeFalse())
		})

		It("should match case-insensitive across multiple fields", func() {
			result := behaviors.SearchableFields("TECH", "backend developer", "TechCorp", "engineering")
			Expect(result).To(BeTrue())
		})

		It("should handle no fields provided", func() {
			result := behaviors.SearchableFields("query")
			Expect(result).To(BeFalse())
		})

		It("should handle empty fields", func() {
			result := behaviors.SearchableFields("query", "", "")
			Expect(result).To(BeFalse())
		})
	})

	Describe("FilterStack", func() {
		var stack *behaviors.FilterStack

		BeforeEach(func() {
			stack = behaviors.NewFilterStack()
		})

		Describe("NewFilterStack", func() {
			It("should create empty stack", func() {
				Expect(stack).NotTo(BeNil())
				Expect(stack.IsEmpty()).To(BeTrue())
			})
		})

		Describe("Push", func() {
			It("should add layer to stack", func() {
				stack.Push(behaviors.FilterLayerSearch)
				Expect(stack.IsEmpty()).To(BeFalse())
			})

			It("should add multiple layers in FIFO order", func() {
				stack.Push(behaviors.FilterLayerSearch)
				stack.Push(behaviors.FilterLayerCategory)
				stack.Push(behaviors.FilterLayerCompany)

				// Pop should return Company (most recently added)
				first := stack.Pop()
				Expect(first).To(Equal(behaviors.FilterLayerCompany))
			})
		})

		Describe("Pop", func() {
			It("should return empty string when stack is empty", func() {
				result := stack.Pop()
				Expect(result).To(Equal(behaviors.FilterLayer("")))
			})

			It("should return and remove most recent layer", func() {
				stack.Push(behaviors.FilterLayerSearch)
				stack.Push(behaviors.FilterLayerCategory)

				result := stack.Pop()
				Expect(result).To(Equal(behaviors.FilterLayerCategory))
				Expect(stack.IsEmpty()).To(BeFalse())
			})

			It("should return layers in FIFO order", func() {
				stack.Push(behaviors.FilterLayerSearch)
				stack.Push(behaviors.FilterLayerCategory)
				stack.Push(behaviors.FilterLayerCompany)

				Expect(stack.Pop()).To(Equal(behaviors.FilterLayerCompany))
				Expect(stack.Pop()).To(Equal(behaviors.FilterLayerCategory))
				Expect(stack.Pop()).To(Equal(behaviors.FilterLayerSearch))
				Expect(stack.IsEmpty()).To(BeTrue())
			})
		})

		Describe("IsEmpty", func() {
			It("should return true for new stack", func() {
				Expect(stack.IsEmpty()).To(BeTrue())
			})

			It("should return false after push", func() {
				stack.Push(behaviors.FilterLayerSearch)
				Expect(stack.IsEmpty()).To(BeFalse())
			})

			It("should return true after all items popped", func() {
				stack.Push(behaviors.FilterLayerSearch)
				stack.Pop()
				Expect(stack.IsEmpty()).To(BeTrue())
			})
		})

		Describe("Clear", func() {
			It("should remove all layers", func() {
				stack.Push(behaviors.FilterLayerSearch)
				stack.Push(behaviors.FilterLayerCategory)
				stack.Clear()
				Expect(stack.IsEmpty()).To(BeTrue())
			})

			It("should work on empty stack", func() {
				stack.Clear()
				Expect(stack.IsEmpty()).To(BeTrue())
			})
		})
	})
})
