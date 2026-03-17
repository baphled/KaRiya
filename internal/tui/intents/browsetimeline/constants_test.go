package browsetimeline_test

import (
	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"

	"github.com/baphled/kariya/internal/tui/intents/browsetimeline"
)

var _ = Describe("Constants", func() {
	Describe("State Type", func() {
		It("should define StateTimeline constant", func() {
			Expect(browsetimeline.StateTimeline).To(Equal(browsetimeline.State("timeline")))
		})

		It("should define StateDeleteConfirm constant", func() {
			Expect(browsetimeline.StateDeleteConfirm).To(Equal(browsetimeline.State("delete_confirm")))
		})

		It("should have distinct state values", func() {
			Expect(browsetimeline.StateTimeline).NotTo(Equal(browsetimeline.StateDeleteConfirm))
		})
	})

	Describe("State String Representation", func() {
		It("should represent StateTimeline as string", func() {
			state := browsetimeline.StateTimeline
			Expect(string(state)).To(Equal("timeline"))
		})

		It("should represent StateDeleteConfirm as string", func() {
			state := browsetimeline.StateDeleteConfirm
			Expect(string(state)).To(Equal("delete_confirm"))
		})
	})

	Describe("State Type Safety", func() {
		It("should be a typed string", func() {
			// Verify State is a distinct type from string
			var state browsetimeline.State = "custom"
			Expect(string(state)).To(Equal("custom"))
		})

		It("should allow comparison between State values", func() {
			state1 := browsetimeline.StateTimeline
			state2 := browsetimeline.StateTimeline
			Expect(state1).To(Equal(state2))
		})
	})
})
