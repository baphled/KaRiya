package skillsmanagement_test

import (
	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"

	"github.com/baphled/kariya/internal/cli/intents/skillsmanagement"
)

var _ = Describe("Constants", func() {
	Describe("State", func() {
		Describe("State Values", func() {
			It("should define StateList for viewing all skills", func() {
				Expect(skillsmanagement.StateList).To(Equal(skillsmanagement.State("list")))
			})

			It("should define StateDetail for viewing single skill details", func() {
				Expect(skillsmanagement.StateDetail).To(Equal(skillsmanagement.State("detail")))
			})

			It("should define StateDetailEvents for viewing events using a skill", func() {
				Expect(skillsmanagement.StateDetailEvents).To(Equal(skillsmanagement.State("detail_events")))
			})

			It("should define StateDetailEventDetail for viewing single event details", func() {
				Expect(skillsmanagement.StateDetailEventDetail).To(Equal(skillsmanagement.State("detail_event_detail")))
			})

			It("should define StateAdd for adding a new skill", func() {
				Expect(skillsmanagement.StateAdd).To(Equal(skillsmanagement.State("add")))
			})

			It("should define StateEdit for editing an existing skill", func() {
				Expect(skillsmanagement.StateEdit).To(Equal(skillsmanagement.State("edit")))
			})

			It("should define StateDelete for confirming deletion", func() {
				Expect(skillsmanagement.StateDelete).To(Equal(skillsmanagement.State("delete")))
			})

			It("should define StateFilter for filter menu", func() {
				Expect(skillsmanagement.StateFilter).To(Equal(skillsmanagement.State("filter")))
			})

			It("should define StateSort for sort menu", func() {
				Expect(skillsmanagement.StateSort).To(Equal(skillsmanagement.State("sort")))
			})
		})

		Describe("Type Safety", func() {
			It("should be a string type", func() {
				var state skillsmanagement.State = "test"
				Expect(string(state)).To(Equal("test"))
			})

			It("should allow comparison", func() {
				state := skillsmanagement.StateList
				Expect(state).To(Equal(skillsmanagement.StateList))
				Expect(state).NotTo(Equal(skillsmanagement.StateDetail))
			})
		})
	})
})
