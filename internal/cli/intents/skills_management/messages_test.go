package skills_management_test

import (
	"errors"

	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"

	"github.com/baphled/kariya/internal/cli/intents/skills_management"
	"github.com/baphled/kariya/internal/domain/career"
)

var _ = Describe("Messages", func() {
	Describe("SkillsLoadedMsg", func() {
		It("should contain loaded skills", func() {
			skills := []*career.Skill{
				{ID: "skill-1", Name: "Go"},
				{ID: "skill-2", Name: "Python"},
			}
			msg := skills_management.SkillsLoadedMsg{
				Skills: skills,
			}
			Expect(msg.Skills).To(HaveLen(2))
			Expect(msg.Skills[0].Name).To(Equal("Go"))
		})

		It("should contain error if load failed", func() {
			err := errors.New("database error")
			msg := skills_management.SkillsLoadedMsg{
				Error: err,
			}
			Expect(msg.Error).To(HaveOccurred())
			Expect(msg.Error.Error()).To(ContainSubstring("database error"))
		})

		It("should allow nil Skills on error", func() {
			msg := skills_management.SkillsLoadedMsg{
				Skills: nil,
				Error:  errors.New("failed"),
			}
			Expect(msg.Skills).To(BeNil())
		})
	})

	Describe("SkillFormCompleteMsg", func() {
		It("should contain the completed skill", func() {
			skill := &career.Skill{ID: "skill-1", Name: "Go"}
			msg := skills_management.SkillFormCompleteMsg{
				Skill: skill,
			}
			Expect(msg.Skill).To(Equal(skill))
		})

		It("should indicate if cancelled", func() {
			msg := skills_management.SkillFormCompleteMsg{
				Cancelled: true,
			}
			Expect(msg.Cancelled).To(BeTrue())
		})

		It("should contain error if form validation failed", func() {
			err := errors.New("validation error")
			msg := skills_management.SkillFormCompleteMsg{
				Error: err,
			}
			Expect(msg.Error).To(HaveOccurred())
		})
	})

	Describe("SkillCreatedMsg", func() {
		It("should contain the created skill", func() {
			skill := &career.Skill{ID: "skill-1", Name: "Go"}
			msg := skills_management.SkillCreatedMsg{
				Skill: skill,
			}
			Expect(msg.Skill).To(Equal(skill))
		})

		It("should contain error if creation failed", func() {
			err := errors.New("create failed")
			msg := skills_management.SkillCreatedMsg{
				Error: err,
			}
			Expect(msg.Error).To(HaveOccurred())
		})
	})

	Describe("SkillUpdatedMsg", func() {
		It("should contain the updated skill", func() {
			skill := &career.Skill{ID: "skill-1", Name: "Go Updated"}
			msg := skills_management.SkillUpdatedMsg{
				Skill: skill,
			}
			Expect(msg.Skill).To(Equal(skill))
		})

		It("should contain error if update failed", func() {
			err := errors.New("update failed")
			msg := skills_management.SkillUpdatedMsg{
				Error: err,
			}
			Expect(msg.Error).To(HaveOccurred())
		})
	})

	Describe("SkillDeletedMsg", func() {
		It("should contain the deleted skill ID", func() {
			msg := skills_management.SkillDeletedMsg{
				SkillID: "skill-1",
			}
			Expect(msg.SkillID).To(Equal("skill-1"))
		})

		It("should contain error if deletion failed", func() {
			err := errors.New("delete failed")
			msg := skills_management.SkillDeletedMsg{
				Error: err,
			}
			Expect(msg.Error).To(HaveOccurred())
		})
	})

	Describe("SkillEventsLoadedMsg", func() {
		It("should contain events for a skill (state-based flow)", func() {
			events := []*career.CareerEvent{
				{ID: "event-1", Text: "Built API"},
				{ID: "event-2", Text: "Led team"},
			}
			msg := skills_management.SkillEventsLoadedMsg{
				Events: events,
			}
			Expect(msg.Events).To(HaveLen(2))
		})

		It("should contain error if load failed", func() {
			err := errors.New("load events failed")
			msg := skills_management.SkillEventsLoadedMsg{
				Error: err,
			}
			Expect(msg.Error).To(HaveOccurred())
		})
	})

	Describe("SkillEventsForModalLoadedMsg", func() {
		It("should contain events for modal display", func() {
			events := []*career.CareerEvent{
				{ID: "event-1", Text: "Built API"},
			}
			msg := skills_management.SkillEventsForModalLoadedMsg{
				Events: events,
			}
			Expect(msg.Events).To(HaveLen(1))
		})

		It("should contain error if load failed", func() {
			err := errors.New("load events failed")
			msg := skills_management.SkillEventsForModalLoadedMsg{
				Error: err,
			}
			Expect(msg.Error).To(HaveOccurred())
		})
	})

	Describe("RequestBrowseEventMsg", func() {
		It("should contain the event to browse", func() {
			event := &career.CareerEvent{ID: "event-1", Text: "Built API"}
			msg := skills_management.RequestBrowseEventMsg{
				Event: event,
			}
			Expect(msg.Event).To(Equal(event))
		})

		It("should contain all events for context", func() {
			events := []*career.CareerEvent{
				{ID: "event-1", Text: "Built API"},
				{ID: "event-2", Text: "Led team"},
			}
			msg := skills_management.RequestBrowseEventMsg{
				AllEvents: events,
			}
			Expect(msg.AllEvents).To(HaveLen(2))
		})

		It("should contain skill name for breadcrumbs", func() {
			msg := skills_management.RequestBrowseEventMsg{
				SkillName: "Go Programming",
			}
			Expect(msg.SkillName).To(Equal("Go Programming"))
		})
	})
})
