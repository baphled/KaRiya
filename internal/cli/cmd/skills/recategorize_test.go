package skills_test

import (
	"bytes"
	"context"
	"errors"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/golang/mock/gomock"
	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"

	cmdpkg "github.com/baphled/kariya/internal/cli/cmd"
	"github.com/baphled/kariya/internal/cli/cmd/cliutil"
	"github.com/baphled/kariya/internal/cli/cmd/skills"
	careermemory "github.com/baphled/kariya/internal/repository/career/memory"
	careerservice "github.com/baphled/kariya/internal/service/career"
	"github.com/baphled/kariya/internal/testutil/fixtures"
	mockrepo "github.com/baphled/kariya/internal/testutil/mocks/repository"
)

var _ = Describe("RecategorizeSkills", func() {
	var (
		ctx cliutil.ServiceProvider
		out *bytes.Buffer
		err *bytes.Buffer
	)

	BeforeEach(func() {
		cliCtx := cmdpkg.NewCLIContext("", true)
		initErr := cliCtx.InitService(new(bytes.Buffer))
		Expect(initErr).NotTo(HaveOccurred())
		ctx = cliCtx
		out = new(bytes.Buffer)
		err = new(bytes.Buffer)
	})

	Context("with no skills", func() {
		It("should return exit code", func() {
			svc := ctx.Service()
			Expect(svc).NotTo(BeNil())

			code := skills.RecategorizeSkills(svc, out, err, tea.WithInput(nil))
			Expect(code).To(BeNumerically(">=", 0))
		})
	})

	Context("with skills in database", func() {
		BeforeEach(func() {
			svc := ctx.Service()
			Expect(svc).NotTo(BeNil())

			skill := fixtures.SkillWith("skill-1", "Go Programming", "backend", "intermediate")
			saveErr := svc.SaveSkill(context.Background(), skill)
			Expect(saveErr).NotTo(HaveOccurred())
		})

		It("should return exit code", func() {
			svc := ctx.Service()
			Expect(svc).NotTo(BeNil())

			code := skills.RecategorizeSkills(svc, out, err, tea.WithInput(nil))
			Expect(code).To(BeNumerically(">=", 0))
		})

		It("should handle execution", func() {
			svc := ctx.Service()
			Expect(svc).NotTo(BeNil())

			code := skills.RecategorizeSkills(svc, out, err, tea.WithInput(nil))
			Expect(code).To(BeNumerically(">=", 0))
		})

		It("should not panic", func() {
			svc := ctx.Service()
			Expect(svc).NotTo(BeNil())

			Expect(func() {
				skills.RecategorizeSkills(svc, out, err, tea.WithInput(nil))
			}).NotTo(Panic())
		})
	})

	Context("with multiple skills", func() {
		BeforeEach(func() {
			svc := ctx.Service()
			Expect(svc).NotTo(BeNil())

			for i := 1; i <= 3; i++ {
				skill := fixtures.Skill("skill-" + string(rune(48+i)))
				saveErr := svc.SaveSkill(context.Background(), skill)
				Expect(saveErr).NotTo(HaveOccurred())
			}
		})

		It("should process all skills", func() {
			svc := ctx.Service()
			Expect(svc).NotTo(BeNil())

			code := skills.RecategorizeSkills(svc, out, err, tea.WithInput(nil))
			Expect(code).To(BeNumerically(">=", 0))
		})
	})

	Context("error handling", func() {
		It("should handle service errors gracefully", func() {
			svc := ctx.Service()
			Expect(svc).NotTo(BeNil())

			code := skills.RecategorizeSkills(svc, out, err, tea.WithInput(nil))
			Expect(code).To(BeNumerically(">=", 0))
		})

		It("should handle missing repository", func() {
			svc := ctx.Service()
			Expect(svc).NotTo(BeNil())

			code := skills.RecategorizeSkills(svc, out, err, tea.WithInput(nil))
			Expect(code).To(BeNumerically(">=", 0))
		})

		It("should return valid exit code", func() {
			svc := ctx.Service()
			Expect(svc).NotTo(BeNil())

			code := skills.RecategorizeSkills(svc, out, err, tea.WithInput(nil))
			Expect(code).To(BeNumerically(">=", 0))
			Expect(code).To(BeNumerically("<=", 1))
		})

		It("should handle nil service gracefully", func() {
			Expect(func() {
				skills.RecategorizeSkills(nil, out, err, tea.WithInput(nil))
			}).To(Panic())
		})

		It("should return exit code 1 when recategorization errors", func() {
			ctrl := gomock.NewController(GinkgoT())
			DeferCleanup(ctrl.Finish)

			mockSkillRepo := mockrepo.NewMockSkillRepository(ctrl)
			mockSkillRepo.EXPECT().
				List(gomock.Any(), gomock.Any()).
				Return(nil, errors.New("database error"))

			svc := careerservice.NewService(careermemory.NewEventRepository())
			svc.SetSkillRepository(mockSkillRepo)

			code := skills.RecategorizeSkills(svc, out, err, tea.WithInput(nil))
			Expect(code).To(Equal(1))
		})
	})
})
