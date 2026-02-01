//nolint:errcheck // Test file - error handling for test setup is not relevant.
package memory

import (
	"context"

	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"

	"github.com/baphled/kariya/internal/testutil/fixtures"
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
		ctx = context.Background() //nolint:fatcontext // test setup
	})

	Describe("GetSkillsForEvent", func() {
		It("returns empty list when no skills are associated", func() {
			event := fixtures.Event("event-1")
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
			event := fixtures.EventWith("event-1", "Built Go microservice with Docker", "", "")
			eventRepo.Create(ctx, event)

			goSkill := fixtures.SkillWith("skill-go", "Go", "backend", "advanced")
			dockerSkill := fixtures.SkillWith("skill-docker", "Docker", "devops", "intermediate")
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
			event1 := fixtures.EventWith("event-1", "Go work", "", "")
			event2 := fixtures.EventWith("event-2", "Python work", "", "")
			eventRepo.Create(ctx, event1)
			eventRepo.Create(ctx, event2)

			goSkill := fixtures.SkillWith("skill-go", "Go", "backend", "")
			pythonSkill := fixtures.SkillWith("skill-python", "Python", "backend", "")
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
			event := fixtures.EventWith("event-1", "Go work", "", "")
			eventRepo.Create(ctx, event)

			goSkill := fixtures.SkillWith("skill-go", "Go", "backend", "")
			skillRepo.Create(ctx, goSkill)

			eventRepo.LinkSkill(ctx, "event-1", "skill-go")
			eventRepo.LinkSkill(ctx, "event-1", "skill-go")
			eventRepo.LinkSkill(ctx, "event-1", "skill-go")

			skills, err := skillRepo.GetSkillsForEvent(ctx, "event-1")
			Expect(err).NotTo(HaveOccurred())
			Expect(skills).To(HaveLen(1))
		})

		It("works with AssociateSkillWithEvent as well", func() {
			event := fixtures.EventWith("event-1", "Go work", "", "")
			eventRepo.Create(ctx, event)

			goSkill := fixtures.SkillWith("skill-go", "Go", "backend", "")
			skillRepo.Create(ctx, goSkill)

			skillRepo.AssociateSkillWithEvent("skill-go", "event-1")

			skills, err := skillRepo.GetSkillsForEvent(ctx, "event-1")
			Expect(err).NotTo(HaveOccurred())
			Expect(skills).To(HaveLen(1))
			Expect(skills[0].Name).To(Equal("Go"))
		})
	})
})
