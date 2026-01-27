package manage_skills_test

import (
	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"

	"github.com/baphled/kariya/internal/cli/intents/manage_skills"
	"github.com/baphled/kariya/internal/domain/career"
)

var _ = Describe("Result", func() {
	Describe("Result struct", func() {
		Describe("Fields", func() {
			It("should have an Action field for tracking the operation performed", func() {
				result := &manage_skills.Result{
					Action: "created",
				}
				Expect(result.Action).To(Equal("created"))
			})

			It("should have a Skill field for the affected skill", func() {
				skill := &career.Skill{
					ID:   "skill-1",
					Name: "Go Programming",
				}
				result := &manage_skills.Result{
					Skill: skill,
				}
				Expect(result.Skill).To(Equal(skill))
			})

			It("should allow nil Skill for cancelled operations", func() {
				result := &manage_skills.Result{
					Action: "cancelled",
					Skill:  nil,
				}
				Expect(result.Skill).To(BeNil())
			})
		})

		Describe("Action Values", func() {
			It("should support 'created' action", func() {
				result := &manage_skills.Result{Action: "created"}
				Expect(result.Action).To(Equal("created"))
			})

			It("should support 'updated' action", func() {
				result := &manage_skills.Result{Action: "updated"}
				Expect(result.Action).To(Equal("updated"))
			})

			It("should support 'deleted' action", func() {
				result := &manage_skills.Result{Action: "deleted"}
				Expect(result.Action).To(Equal("deleted"))
			})

			It("should support 'cancelled' action", func() {
				result := &manage_skills.Result{Action: "cancelled"}
				Expect(result.Action).To(Equal("cancelled"))
			})
		})

		Describe("Zero Value", func() {
			It("should have empty Action by default", func() {
				result := &manage_skills.Result{}
				Expect(result.Action).To(BeEmpty())
			})

			It("should have nil Skill by default", func() {
				result := &manage_skills.Result{}
				Expect(result.Skill).To(BeNil())
			})
		})
	})
})
