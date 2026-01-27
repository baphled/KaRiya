package skills_management_test

import (
	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"

	"github.com/baphled/kariya/internal/cli/intents/skills_management"
)

var _ = Describe("Constants", func() {
	Describe("State", func() {
		Describe("State Values", func() {
			It("should define StateList for viewing all skills", func() {
				Expect(skills_management.StateList).To(Equal(skills_management.State("list")))
			})

			It("should define StateDetail for viewing single skill details", func() {
				Expect(skills_management.StateDetail).To(Equal(skills_management.State("detail")))
			})

			It("should define StateDetailEvents for viewing events using a skill", func() {
				Expect(skills_management.StateDetailEvents).To(Equal(skills_management.State("detail_events")))
			})

			It("should define StateDetailEventDetail for viewing single event details", func() {
				Expect(skills_management.StateDetailEventDetail).To(Equal(skills_management.State("detail_event_detail")))
			})

			It("should define StateAdd for adding a new skill", func() {
				Expect(skills_management.StateAdd).To(Equal(skills_management.State("add")))
			})

			It("should define StateEdit for editing an existing skill", func() {
				Expect(skills_management.StateEdit).To(Equal(skills_management.State("edit")))
			})

			It("should define StateDelete for confirming deletion", func() {
				Expect(skills_management.StateDelete).To(Equal(skills_management.State("delete")))
			})

			It("should define StateFilter for filter menu", func() {
				Expect(skills_management.StateFilter).To(Equal(skills_management.State("filter")))
			})

			It("should define StateSort for sort menu", func() {
				Expect(skills_management.StateSort).To(Equal(skills_management.State("sort")))
			})
		})

		Describe("Type Safety", func() {
			It("should be a string type", func() {
				var state skills_management.State = "test"
				Expect(string(state)).To(Equal("test"))
			})

			It("should allow comparison", func() {
				state := skills_management.StateList
				Expect(state).To(Equal(skills_management.StateList))
				Expect(state).NotTo(Equal(skills_management.StateDetail))
			})
		})
	})
})
