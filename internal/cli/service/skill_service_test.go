package service

import (
	"context"

	careermemory "github.com/baphled/kariya/internal/repository/career/memory"
	"github.com/baphled/kariya/internal/testutil/fixtures"
	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
)

var _ = Describe("CLISkillCreator", func() {
	var (
		skillService *CLISkillCreator
		skillRepo    *careermemory.SkillRepository
		ctx          context.Context
	)

	BeforeEach(func() {
		ctx = context.Background()
		skillRepo = careermemory.NewSkillRepository()
		skillService = NewCLISkillCreator(skillRepo)
	})

	Describe("NewCLISkillCreator", func() {
		It("should create a new service", func() {
			Expect(skillService).ToNot(BeNil())
			Expect(skillService.skillRepo).To(Equal(skillRepo))
		})
	})

	Describe("Create", func() {
		It("should create a new skill", func() {
			skill := fixtures.SkillWith("skill-1", "Go", "backend", "advanced")

			err := skillService.Create(ctx, skill)
			Expect(err).ToNot(HaveOccurred())

			retrieved, err := skillRepo.GetByID(ctx, "skill-1")
			Expect(err).ToNot(HaveOccurred())
			Expect(retrieved.Name).To(Equal("Go"))
		})

		It("should return error for invalid skill", func() {
			invalidSkill := fixtures.SkillWith("skill-1", "", "invalid", "invalid")

			err := skillService.Create(ctx, invalidSkill)
			Expect(err).To(HaveOccurred())
		})

		It("should return error for duplicate skill", func() {
			skill := fixtures.SkillWith("skill-1", "Go", "backend", "advanced")

			err := skillService.Create(ctx, skill)
			Expect(err).ToNot(HaveOccurred())

			err = skillService.Create(ctx, skill)
			Expect(err).To(HaveOccurred())
		})
	})
})
