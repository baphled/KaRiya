package skillsmanagement_test

import (
	"errors"

	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"

	"github.com/baphled/kariya/internal/cli/intents/skillsmanagement"
<<<<<<< HEAD
	"github.com/baphled/kariya/internal/testutil/fixtures"
=======
	"github.com/baphled/kariya/internal/domain/career"
	"github.com/baphled/kariya/internal/service/career/skillinference"
>>>>>>> 9ea4d2d4 (test(intents): add test coverage for filterNewSuggestions and message structs)
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

	Describe("SkillSuggestionsLoadedMsg", func() {
		It("should store skill suggestions", func() {
			suggestions := []skillinference.SkillSuggestion{
				{Name: "Go", Category: "backend", Confidence: 0.95},
				{Name: "Docker", Category: "devops", Confidence: 0.85},
			}
			msg := skillsmanagement.SkillSuggestionsLoadedMsg{
				Suggestions: suggestions,
			}
			Expect(msg.Suggestions).To(HaveLen(2))
			Expect(msg.Suggestions[0].Name).To(Equal("Go"))
			Expect(msg.Suggestions[0].Confidence).To(Equal(0.95))
		})

		It("should store existing skill names for filtering", func() {
			msg := skillsmanagement.SkillSuggestionsLoadedMsg{
				Suggestions:        []skillinference.SkillSuggestion{{Name: "Go", Category: "backend", Confidence: 0.9}},
				ExistingSkillNames: []string{"Docker", "Kubernetes"},
			}
			Expect(msg.ExistingSkillNames).To(HaveLen(2))
			Expect(msg.ExistingSkillNames).To(ContainElements("Docker", "Kubernetes"))
		})

		It("should store error on inference failure", func() {
			err := errors.New("inference failed")
			msg := skillsmanagement.SkillSuggestionsLoadedMsg{
				Error: err,
			}
			Expect(msg.Suggestions).To(BeNil())
			Expect(msg.Error).To(HaveOccurred())
		})

		It("should handle empty suggestions", func() {
			msg := skillsmanagement.SkillSuggestionsLoadedMsg{
				Suggestions:        []skillinference.SkillSuggestion{},
				ExistingSkillNames: []string{},
			}
			Expect(msg.Suggestions).To(BeEmpty())
			Expect(msg.ExistingSkillNames).To(BeEmpty())
		})
	})

	Describe("SkillsCreatedMsg", func() {
		It("should store created skills", func() {
			skills := []*career.Skill{
				{ID: "s1", Name: "Go"},
				{ID: "s2", Name: "Docker"},
			}
			msg := skillsmanagement.SkillsCreatedMsg{
				Skills: skills,
			}
			Expect(msg.Skills).To(HaveLen(2))
			Expect(msg.Skills[0].Name).To(Equal("Go"))
		})

		It("should store error on creation failure", func() {
			err := errors.New("create skills failed")
			msg := skillsmanagement.SkillsCreatedMsg{
				Error: err,
			}
			Expect(msg.Skills).To(BeNil())
			Expect(msg.Error).To(HaveOccurred())
		})

		It("should handle empty skills list on success", func() {
			msg := skillsmanagement.SkillsCreatedMsg{
				Skills: []*career.Skill{},
			}
			Expect(msg.Skills).To(BeEmpty())
			Expect(msg.Error).To(BeNil())
		})
	})
})
