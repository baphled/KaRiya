package career

import (
	"context"
	"errors"

	"github.com/golang/mock/gomock"
	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"

	"github.com/baphled/kariya/internal/testutil/fixtures"
	mockrepo "github.com/baphled/kariya/internal/testutil/mocks/repository"
)

var _ = Describe("Career Service - Skill Methods", func() {
	var (
		ctrl          *gomock.Controller
		mockEventRepo *mockrepo.MockEventRepository
		mockSkillRepo *mockrepo.MockSkillRepository
		service       *Service
		ctx           context.Context
	)

	BeforeEach(func() {
		ctrl = gomock.NewController(GinkgoT())
		mockEventRepo = mockrepo.NewMockEventRepository(ctrl)
		mockSkillRepo = mockrepo.NewMockSkillRepository(ctrl)
		service = NewService(mockEventRepo)
		service.SetSkillRepository(mockSkillRepo)
		ctx = context.Background()
	})

	AfterEach(func() {
		ctrl.Finish()
	})

	Describe("SaveSkill", func() {
		Context("when skill repository is not configured", func() {
			It("returns ErrSkillRepositoryNotConfigured", func() {
				svc := NewService(mockEventRepo)
				skill := fixtures.SkillWith("", "Go", "backend", "advanced")

				err := svc.SaveSkill(ctx, skill)

				Expect(err).To(MatchError(ErrSkillRepositoryNotConfigured))
			})
		})

		Context("when skill is nil", func() {
			It("returns an error", func() {
				err := service.SaveSkill(ctx, nil)

				Expect(err).To(HaveOccurred())
				Expect(err.Error()).To(ContainSubstring("skill cannot be nil"))
			})
		})

		Context("when skill already has an ID", func() {
			It("returns nil without persisting", func() {
				skill := fixtures.SkillWith("existing-id", "Go", "backend", "advanced")

				err := service.SaveSkill(ctx, skill)

				Expect(err).NotTo(HaveOccurred())
			})
		})

		Context("when skill has no ID", func() {
			It("generates a UUID and persists the skill", func() {
				skill := fixtures.SkillWith("", "Go", "backend", "advanced")

				mockSkillRepo.EXPECT().
					Create(gomock.Any(), gomock.Any()).
					Return(nil).
					Times(1)

				err := service.SaveSkill(ctx, skill)

				Expect(err).NotTo(HaveOccurred())
				Expect(skill.ID).NotTo(BeEmpty())
				Expect(skill.ID).To(HaveLen(36))
				Expect(skill.CreatedAt).NotTo(BeZero())
				Expect(skill.UpdatedAt).NotTo(BeZero())
			})

			It("returns an error when repository fails", func() {
				skill := fixtures.SkillWith("", "Go", "backend", "advanced")

				mockSkillRepo.EXPECT().
					Create(gomock.Any(), gomock.Any()).
					Return(errors.New("database error")).
					Times(1)

				err := service.SaveSkill(ctx, skill)

				Expect(err).To(HaveOccurred())
				Expect(err.Error()).To(ContainSubstring("failed to create skill"))
			})
		})
	})

	Describe("LinkSkillToEvent", func() {
		Context("when event ID is empty", func() {
			It("returns an error", func() {
				err := service.LinkSkillToEvent(ctx, "", "skill-1")

				Expect(err).To(HaveOccurred())
				Expect(err.Error()).To(ContainSubstring("event ID cannot be empty"))
			})
		})

		Context("when skill ID is empty", func() {
			It("returns an error", func() {
				err := service.LinkSkillToEvent(ctx, "event-1", "")

				Expect(err).To(HaveOccurred())
				Expect(err.Error()).To(ContainSubstring("skill ID cannot be empty"))
			})
		})

		Context("when both IDs are valid", func() {
			It("links the skill to the event", func() {
				mockEventRepo.EXPECT().
					LinkSkill(gomock.Any(), "event-1", "skill-1").
					Return(nil).
					Times(1)

				err := service.LinkSkillToEvent(ctx, "event-1", "skill-1")

				Expect(err).NotTo(HaveOccurred())
			})

			It("returns an error when repository fails", func() {
				mockEventRepo.EXPECT().
					LinkSkill(gomock.Any(), "event-1", "skill-1").
					Return(errors.New("link failed")).
					Times(1)

				err := service.LinkSkillToEvent(ctx, "event-1", "skill-1")

				Expect(err).To(HaveOccurred())
				Expect(err.Error()).To(ContainSubstring("failed to link skill to event"))
			})
		})
	})
})
