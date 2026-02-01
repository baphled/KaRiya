package skillsmanagement_test

import (
	"errors"

	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"

	"github.com/baphled/kariya/internal/cli/intents/skillsmanagement"
	"github.com/baphled/kariya/internal/testutil/fixtures"
)

var _ = Describe("Messages", func() {
	Describe("SkillsLoadedMsg", func() {
		It("should contain loaded skills", func() {
			skills := fixtures.Skills(2)
			msg := skillsmanagement.SkillsLoadedMsg{
				Skills: skills,
			}
			Expect(msg.Skills).To(HaveLen(2))
			Expect(msg.Skills[0].Name).NotTo(BeEmpty())
		})

		It("should contain error if load failed", func() {
			err := errors.New("database error")
			msg := skillsmanagement.SkillsLoadedMsg{
				Error: err,
			}
			Expect(msg.Error).To(HaveOccurred())
			Expect(msg.Error.Error()).To(ContainSubstring("database error"))
		})

		It("should allow nil Skills on error", func() {
			msg := skillsmanagement.SkillsLoadedMsg{
				Skills: nil,
				Error:  errors.New("failed"),
			}
			Expect(msg.Skills).To(BeNil())
		})
	})

	Describe("SkillFormCompleteMsg", func() {
		It("should contain the completed skill", func() {
			skill := fixtures.SkillWith("skill-1", "Go", "backend", "advanced")
			msg := skillsmanagement.SkillFormCompleteMsg{
				Skill: skill,
			}
			Expect(msg.Skill).To(Equal(skill))
		})

		It("should indicate if cancelled", func() {
			msg := skillsmanagement.SkillFormCompleteMsg{
				Cancelled: true,
			}
			Expect(msg.Cancelled).To(BeTrue())
		})

		It("should contain error if form validation failed", func() {
			err := errors.New("validation error")
			msg := skillsmanagement.SkillFormCompleteMsg{
				Error: err,
			}
			Expect(msg.Error).To(HaveOccurred())
		})
	})

	Describe("SkillCreatedMsg", func() {
		It("should contain the created skill", func() {
			skill := fixtures.SkillWith("skill-1", "Go", "backend", "advanced")
			msg := skillsmanagement.SkillCreatedMsg{
				Skill: skill,
			}
			Expect(msg.Skill).To(Equal(skill))
		})

		It("should contain error if creation failed", func() {
			err := errors.New("create failed")
			msg := skillsmanagement.SkillCreatedMsg{
				Error: err,
			}
			Expect(msg.Error).To(HaveOccurred())
		})
	})

	Describe("SkillUpdatedMsg", func() {
		It("should contain the updated skill", func() {
			skill := fixtures.SkillWith("skill-1", "Go Updated", "backend", "advanced")
			msg := skillsmanagement.SkillUpdatedMsg{
				Skill: skill,
			}
			Expect(msg.Skill).To(Equal(skill))
		})

		It("should contain error if update failed", func() {
			err := errors.New("update failed")
			msg := skillsmanagement.SkillUpdatedMsg{
				Error: err,
			}
			Expect(msg.Error).To(HaveOccurred())
		})
	})

	Describe("SkillDeletedMsg", func() {
		It("should contain the deleted skill ID", func() {
			msg := skillsmanagement.SkillDeletedMsg{
				SkillID: "skill-1",
			}
			Expect(msg.SkillID).To(Equal("skill-1"))
		})

		It("should contain error if deletion failed", func() {
			err := errors.New("delete failed")
			msg := skillsmanagement.SkillDeletedMsg{
				Error: err,
			}
			Expect(msg.Error).To(HaveOccurred())
		})
	})

	Describe("SkillEventsLoadedMsg", func() {
		It("should contain events for a skill (state-based flow)", func() {
			events := fixtures.Events(2)
			msg := skillsmanagement.SkillEventsLoadedMsg{
				Events: events,
			}
			Expect(msg.Events).To(HaveLen(2))
		})

		It("should contain error if load failed", func() {
			err := errors.New("load events failed")
			msg := skillsmanagement.SkillEventsLoadedMsg{
				Error: err,
			}
			Expect(msg.Error).To(HaveOccurred())
		})
	})

	Describe("SkillEventsForModalLoadedMsg", func() {
		It("should contain events for modal display", func() {
			events := fixtures.Events(1)
			msg := skillsmanagement.SkillEventsForModalLoadedMsg{
				Events: events,
			}
			Expect(msg.Events).To(HaveLen(1))
		})

		It("should contain error if load failed", func() {
			err := errors.New("load events failed")
			msg := skillsmanagement.SkillEventsForModalLoadedMsg{
				Error: err,
			}
			Expect(msg.Error).To(HaveOccurred())
		})
	})

	Describe("RequestBrowseEventMsg", func() {
		It("should contain the event to browse", func() {
			event := fixtures.Event("event-1")
			msg := skillsmanagement.RequestBrowseEventMsg{
				Event: event,
			}
			Expect(msg.Event).To(Equal(event))
		})

		It("should contain all events for context", func() {
			events := fixtures.Events(2)
			msg := skillsmanagement.RequestBrowseEventMsg{
				AllEvents: events,
			}
			Expect(msg.AllEvents).To(HaveLen(2))
		})

		It("should contain skill name for breadcrumbs", func() {
			msg := skillsmanagement.RequestBrowseEventMsg{
				SkillName: "Go Programming",
			}
			Expect(msg.SkillName).To(Equal("Go Programming"))
		})
	})
})
