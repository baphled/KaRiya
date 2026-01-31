//nolint:errcheck // Test file - error handling for test setup is not relevant.
package memory

import (
	"context"
	"time"

	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"

	"github.com/baphled/kariya/internal/domain/career"
)

var _ = Describe("SkillRepository", func() {
	var (
		skillRepo *SkillRepository
		eventRepo *EventRepository
		ctx       context.Context
	)

	BeforeEach(func() {
		eventRepo = NewEventRepository()
		skillRepo = NewSkillRepository()
		skillRepo.SetEventRepository(eventRepo)
		eventRepo.SetSkillRepository(skillRepo)
		ctx = context.Background()
	})

	Describe("GetSkillsForEvent", func() {
		It("returns empty list when no skills are associated", func() {
			event := &career.Event{
				ID:   "event-1",
				Text: "Some event",
				Date: time.Now(),
			}
			eventRepo.Create(ctx, event)

			skills, err := skillRepo.GetSkillsForEvent(ctx, "event-1")

			Expect(err).NotTo(HaveOccurred())
			Expect(skills).To(BeEmpty())
		})

		It("returns empty list for unknown event ID", func() {
			skills, err := skillRepo.GetSkillsForEvent(ctx, "nonexistent")

			Expect(err).NotTo(HaveOccurred())
			Expect(skills).To(BeEmpty())
		})

		It("returns skills linked via EventRepository.LinkSkill", func() {
			event := &career.Event{
				ID:   "event-1",
				Text: "Built Go microservice with Docker",
				Date: time.Now(),
			}
			eventRepo.Create(ctx, event)

			goSkill := &career.Skill{
				ID:       "skill-go",
				Name:     "Go",
				Category: "backend",
				Level:    "advanced",
			}
			dockerSkill := &career.Skill{
				ID:       "skill-docker",
				Name:     "Docker",
				Category: "devops",
				Level:    "intermediate",
			}
			skillRepo.Create(ctx, goSkill)
			skillRepo.Create(ctx, dockerSkill)

			eventRepo.LinkSkill(ctx, "event-1", "skill-go")
			eventRepo.LinkSkill(ctx, "event-1", "skill-docker")

			skills, err := skillRepo.GetSkillsForEvent(ctx, "event-1")

			Expect(err).NotTo(HaveOccurred())
			Expect(skills).To(HaveLen(2))

			names := []string{skills[0].Name, skills[1].Name}
			Expect(names).To(ContainElement("Go"))
			Expect(names).To(ContainElement("Docker"))
		})

		It("returns skills for correct event only", func() {
			event1 := &career.Event{ID: "event-1", Text: "Go work", Date: time.Now()}
			event2 := &career.Event{ID: "event-2", Text: "Python work", Date: time.Now()}
			eventRepo.Create(ctx, event1)
			eventRepo.Create(ctx, event2)

			goSkill := &career.Skill{ID: "skill-go", Name: "Go", Category: "backend"}
			pythonSkill := &career.Skill{ID: "skill-python", Name: "Python", Category: "backend"}
			skillRepo.Create(ctx, goSkill)
			skillRepo.Create(ctx, pythonSkill)

			eventRepo.LinkSkill(ctx, "event-1", "skill-go")
			eventRepo.LinkSkill(ctx, "event-2", "skill-python")

			skills1, err := skillRepo.GetSkillsForEvent(ctx, "event-1")
			Expect(err).NotTo(HaveOccurred())
			Expect(skills1).To(HaveLen(1))
			Expect(skills1[0].Name).To(Equal("Go"))

			skills2, err := skillRepo.GetSkillsForEvent(ctx, "event-2")
			Expect(err).NotTo(HaveOccurred())
			Expect(skills2).To(HaveLen(1))
			Expect(skills2[0].Name).To(Equal("Python"))
		})

		It("does not return duplicate skills from repeated LinkSkill calls", func() {
			event := &career.Event{ID: "event-1", Text: "Go work", Date: time.Now()}
			eventRepo.Create(ctx, event)

			goSkill := &career.Skill{ID: "skill-go", Name: "Go", Category: "backend"}
			skillRepo.Create(ctx, goSkill)

			eventRepo.LinkSkill(ctx, "event-1", "skill-go")
			eventRepo.LinkSkill(ctx, "event-1", "skill-go")
			eventRepo.LinkSkill(ctx, "event-1", "skill-go")

			skills, err := skillRepo.GetSkillsForEvent(ctx, "event-1")
			Expect(err).NotTo(HaveOccurred())
			Expect(skills).To(HaveLen(1))
		})

		It("works with AssociateSkillWithEvent as well", func() {
			event := &career.Event{ID: "event-1", Text: "Go work", Date: time.Now()}
			eventRepo.Create(ctx, event)

			goSkill := &career.Skill{ID: "skill-go", Name: "Go", Category: "backend"}
			skillRepo.Create(ctx, goSkill)

			skillRepo.AssociateSkillWithEvent("skill-go", "event-1")

			skills, err := skillRepo.GetSkillsForEvent(ctx, "event-1")
			Expect(err).NotTo(HaveOccurred())
			Expect(skills).To(HaveLen(1))
			Expect(skills[0].Name).To(Equal("Go"))
		})
	})
})
