package browse_timeline_test

import (
	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"

	"github.com/baphled/kariya/internal/cli/intents/browse_timeline"
)

var _ = Describe("Constants", func() {
	Describe("State Type", func() {
		It("should define StateTimeline constant", func() {
			Expect(browse_timeline.StateTimeline).To(Equal(browse_timeline.State("timeline")))
		})

		It("should define StateDeleteConfirm constant", func() {
			Expect(browse_timeline.StateDeleteConfirm).To(Equal(browse_timeline.State("delete_confirm")))
		})

		It("should have distinct state values", func() {
			Expect(browse_timeline.StateTimeline).NotTo(Equal(browse_timeline.StateDeleteConfirm))
		})
	})

	Describe("State String Representation", func() {
		It("should represent StateTimeline as string", func() {
			state := browse_timeline.StateTimeline
			Expect(string(state)).To(Equal("timeline"))
		})

		It("should represent StateDeleteConfirm as string", func() {
			state := browse_timeline.StateDeleteConfirm
			Expect(string(state)).To(Equal("delete_confirm"))
		})
	})

	Describe("State Type Safety", func() {
		It("should be a typed string", func() {
			// Verify State is a distinct type from string
			var state browse_timeline.State = "custom"
			Expect(string(state)).To(Equal("custom"))
		})

		It("should allow comparison between State values", func() {
			state1 := browse_timeline.StateTimeline
			state2 := browse_timeline.StateTimeline
			Expect(state1).To(Equal(state2))
		})
	})
})
