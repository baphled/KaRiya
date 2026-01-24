package fact_management_test

import (
	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"

	"github.com/baphled/kariya/internal/cli/intents/fact_management"
	"github.com/baphled/kariya/internal/domain/career"
)

var _ = Describe("Result", func() {
	Describe("Construction", func() {
		It("should create an empty result", func() {
			result := &fact_management.Result{}
			Expect(result).NotTo(BeNil())
		})

		It("should store action", func() {
			result := &fact_management.Result{
				Action: "created",
			}
			Expect(result.Action).To(Equal("created"))
		})

		It("should store fact", func() {
			fact := &career.Fact{ID: "fact-1", Text: "Test fact"}
			result := &fact_management.Result{
				Fact: fact,
			}
			Expect(result.Fact).To(Equal(fact))
		})

		It("should store facts list", func() {
			facts := []*career.Fact{
				{ID: "fact-1", Text: "Fact 1"},
				{ID: "fact-2", Text: "Fact 2"},
			}
			result := &fact_management.Result{
				Facts: facts,
			}
			Expect(result.Facts).To(HaveLen(2))
		})

		It("should store message", func() {
			result := &fact_management.Result{
				Message: "Operation completed successfully",
			}
			Expect(result.Message).To(Equal("Operation completed successfully"))
		})

		It("should store scroll position", func() {
			result := &fact_management.Result{
				ScrollPosition: 42,
			}
			Expect(result.ScrollPosition).To(Equal(42))
		})
	})

	Describe("Action Types", func() {
		It("should support created action", func() {
			result := &fact_management.Result{
				Action:  "created",
				Message: "Fact created successfully",
			}
			Expect(result.Action).To(Equal("created"))
		})

		It("should support updated action", func() {
			result := &fact_management.Result{
				Action:  "updated",
				Message: "Fact updated successfully",
			}
			Expect(result.Action).To(Equal("updated"))
		})

		It("should support deleted action", func() {
			result := &fact_management.Result{
				Action:  "deleted",
				Message: "Fact deleted successfully",
			}
			Expect(result.Action).To(Equal("deleted"))
		})

		It("should support none action", func() {
			result := &fact_management.Result{
				Action: "none",
			}
			Expect(result.Action).To(Equal("none"))
		})
	})

	Describe("Complete Result", func() {
		It("should contain all fields for a complete operation", func() {
			fact := &career.Fact{ID: "fact-1", Text: "Test fact"}
			facts := []*career.Fact{fact}
			result := &fact_management.Result{
				Action:         "created",
				Fact:           fact,
				Facts:          facts,
				Message:        "Fact created successfully",
				ScrollPosition: 0,
			}

			Expect(result.Action).To(Equal("created"))
			Expect(result.Fact).To(Equal(fact))
			Expect(result.Facts).To(HaveLen(1))
			Expect(result.Message).NotTo(BeEmpty())
			Expect(result.ScrollPosition).To(Equal(0))
		})
	})
})
