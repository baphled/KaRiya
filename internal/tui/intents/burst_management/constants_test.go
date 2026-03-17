package burst_management_test

import (
	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"

	"github.com/baphled/kariya/internal/tui/intents/burst_management"
)

var _ = Describe("Constants", func() {
	Describe("State Type", func() {
		It("should define StateList constant", func() {
			Expect(burst_management.StateList).To(Equal(burst_management.State("list")))
		})

		It("should define StateDetail constant", func() {
			Expect(burst_management.StateDetail).To(Equal(burst_management.State("detail")))
		})

		It("should define StateDetailEvents constant", func() {
			Expect(burst_management.StateDetailEvents).To(Equal(burst_management.State("detail_events")))
		})

		It("should define StateDetailFacts constant", func() {
			Expect(burst_management.StateDetailFacts).To(Equal(burst_management.State("detail_facts")))
		})

		It("should define StateEdit constant", func() {
			Expect(burst_management.StateEdit).To(Equal(burst_management.State("edit")))
		})

		It("should define StateDeleteConfirm constant", func() {
			Expect(burst_management.StateDeleteConfirm).To(Equal(burst_management.State("delete_confirm")))
		})

		It("should define StateConfirm constant", func() {
			Expect(burst_management.StateConfirm).To(Equal(burst_management.State("confirm")))
		})

		It("should define StateExtractingFacts constant", func() {
			Expect(burst_management.StateExtractingFacts).To(Equal(burst_management.State("extracting_facts")))
		})

		It("should define StateSuggesting constant", func() {
			Expect(burst_management.StateSuggesting).To(Equal(burst_management.State("suggesting")))
		})

		It("should define StateSuggestionReview constant", func() {
			Expect(burst_management.StateSuggestionReview).To(Equal(burst_management.State("suggestion_review")))
		})

		It("should have distinct state values", func() {
			Expect(burst_management.StateList).NotTo(Equal(burst_management.StateDetail))
			Expect(burst_management.StateDetail).NotTo(Equal(burst_management.StateEdit))
			Expect(burst_management.StateSuggesting).NotTo(Equal(burst_management.StateSuggestionReview))
		})
	})

	Describe("State String Representation", func() {
		It("should represent StateList as string", func() {
			state := burst_management.StateList
			Expect(string(state)).To(Equal("list"))
		})

		It("should represent StateDetail as string", func() {
			state := burst_management.StateDetail
			Expect(string(state)).To(Equal("detail"))
		})

		It("should represent StateSuggesting as string", func() {
			state := burst_management.StateSuggesting
			Expect(string(state)).To(Equal("suggesting"))
		})

		It("should represent StateSuggestionReview as string", func() {
			state := burst_management.StateSuggestionReview
			Expect(string(state)).To(Equal("suggestion_review"))
		})
	})

	Describe("State Type Safety", func() {
		It("should be a typed string", func() {
			// Verify State is a distinct type from string.
			var state burst_management.State = "custom"
			Expect(string(state)).To(Equal("custom"))
		})

		It("should allow comparison between State values", func() {
			state1 := burst_management.StateList
			state2 := burst_management.StateList
			Expect(state1).To(Equal(state2))
		})
	})
})
