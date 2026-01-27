package manage_skills_test

import (
	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"

	"github.com/baphled/kariya/internal/cli/intents/manage_skills"
)

var _ = Describe("Constants", func() {
	Describe("State", func() {
		Describe("State Values", func() {
			It("should define StateList for viewing all skills", func() {
				Expect(manage_skills.StateList).To(Equal(manage_skills.State("list")))
			})

			It("should define StateDetail for viewing single skill details", func() {
				Expect(manage_skills.StateDetail).To(Equal(manage_skills.State("detail")))
			})

			It("should define StateDetailEvents for viewing events using a skill", func() {
				Expect(manage_skills.StateDetailEvents).To(Equal(manage_skills.State("detail_events")))
			})

			It("should define StateDetailEventDetail for viewing single event details", func() {
				Expect(manage_skills.StateDetailEventDetail).To(Equal(manage_skills.State("detail_event_detail")))
			})

			It("should define StateAdd for adding a new skill", func() {
				Expect(manage_skills.StateAdd).To(Equal(manage_skills.State("add")))
			})

			It("should define StateEdit for editing an existing skill", func() {
				Expect(manage_skills.StateEdit).To(Equal(manage_skills.State("edit")))
			})

			It("should define StateDelete for confirming deletion", func() {
				Expect(manage_skills.StateDelete).To(Equal(manage_skills.State("delete")))
			})

			It("should define StateFilter for filter menu", func() {
				Expect(manage_skills.StateFilter).To(Equal(manage_skills.State("filter")))
			})

			It("should define StateSort for sort menu", func() {
				Expect(manage_skills.StateSort).To(Equal(manage_skills.State("sort")))
			})
		})

		Describe("Type Safety", func() {
			It("should be a string type", func() {
				var state manage_skills.State = "test"
				Expect(string(state)).To(Equal("test"))
			})

			It("should allow comparison", func() {
				state := manage_skills.StateList
				Expect(state).To(Equal(manage_skills.StateList))
				Expect(state).NotTo(Equal(manage_skills.StateDetail))
			})
		})
	})
})
