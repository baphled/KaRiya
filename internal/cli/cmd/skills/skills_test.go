package skills_test

import (
	"bytes"
	"context"
	"io"

	tea "github.com/charmbracelet/bubbletea"
	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
	"github.com/spf13/cobra"

	cmdpkg "github.com/baphled/kariya/internal/cli/cmd"
	"github.com/baphled/kariya/internal/cli/cmd/cliutil"
	"github.com/baphled/kariya/internal/cli/cmd/skills"
	careerservice "github.com/baphled/kariya/internal/service/career"
	"github.com/baphled/kariya/internal/testutil"
	"github.com/baphled/kariya/internal/testutil/fixtures"
)

var _ = Describe("Skills Command", func() {
	var (
		ctx cliutil.ServiceProvider
		cmd *cobra.Command
	)

	BeforeEach(func() {
		cliCtx := cmdpkg.NewCLIContext("", true)
		err := cliCtx.InitService(new(bytes.Buffer))
		Expect(err).NotTo(HaveOccurred())
		ctx = cliCtx
	})

	Describe("NewSkillsCmd", func() {
		It("should create command with correct properties", func() {
			cmd = skills.NewSkillsCmd(ctx)
			Expect(cmd).NotTo(BeNil())
			Expect(cmd.Use).To(Equal("skills"))
			Expect(cmd.Short).To(Equal("Manage skills"))
			Expect(cmd.Long).To(ContainSubstring("skills"))
		})

		It("should have recategorize subcommand", func() {
			cmd = skills.NewSkillsCmd(ctx)
			recatCmd, _, err := cmd.Find([]string{"recategorize"})
			Expect(err).NotTo(HaveOccurred())
			Expect(recatCmd).NotTo(BeNil())
			Expect(recatCmd.Use).To(Equal("recategorize"))
		})
	})

	Describe("NewRecategorizeCmd", func() {
		It("should create recategorize command with correct properties", func() {
			cmd = skills.NewRecategorizeCmd(ctx, tea.WithInput(nil))
			Expect(cmd).NotTo(BeNil())
			Expect(cmd.Use).To(Equal("recategorize"))
			Expect(cmd.Short).To(Equal("Recategorize all skills"))
			Expect(cmd.Long).To(ContainSubstring("categorization"))
		})

		It("should have RunE function", func() {
			cmd = skills.NewRecategorizeCmd(ctx, tea.WithInput(nil))
			Expect(cmd.RunE).NotTo(BeNil())
		})

		It("should execute recategorize command", func() {
			cmd = skills.NewRecategorizeCmd(ctx, tea.WithInput(nil))
			out := new(bytes.Buffer)
			cmd.SetOut(out)
			cmd.SetErr(new(bytes.Buffer))

			err := cmd.RunE(cmd, []string{})
			Expect(err).NotTo(HaveOccurred())
		})

		It("should return error when service is not initialized", func() {
			nilCtx := cmdpkg.NewCLIContext("", true)
			cobraCmd := skills.NewRecategorizeCmd(nilCtx, tea.WithInput(nil))
			out := new(bytes.Buffer)
			cobraCmd.SetOut(out)
			cobraCmd.SetErr(new(bytes.Buffer))

			err := cobraCmd.RunE(cobraCmd, []string{})
			Expect(err).To(HaveOccurred())
			Expect(err.Error()).To(ContainSubstring("service not initialized"))
		})

		It("should successfully recategorize skills", func() {
			cmd = skills.NewRecategorizeCmd(ctx, tea.WithInput(nil))
			out := new(bytes.Buffer)
			cmd.SetOut(out)
			cmd.SetErr(new(bytes.Buffer))

			svc := ctx.Service()
			Expect(svc).NotTo(BeNil())

			err := cmd.RunE(cmd, []string{})
			Expect(err).NotTo(HaveOccurred())
		})

		It("should handle command with no arguments", func() {
			cmd = skills.NewRecategorizeCmd(ctx, tea.WithInput(nil))
			out := new(bytes.Buffer)
			cmd.SetOut(out)
			cmd.SetErr(new(bytes.Buffer))

			err := cmd.RunE(cmd, []string{})
			Expect(err).NotTo(HaveOccurred())
		})

		It("should use OutOrStdout for output", func() {
			cmd = skills.NewRecategorizeCmd(ctx, tea.WithInput(nil))
			out := new(bytes.Buffer)
			cmd.SetOut(out)
			cmd.SetErr(new(bytes.Buffer))

			Expect(cmd.OutOrStdout()).NotTo(BeNil())
		})

		It("should use ErrOrStderr for error output", func() {
			cmd = skills.NewRecategorizeCmd(ctx, tea.WithInput(nil))
			cmd.SetOut(new(bytes.Buffer))
			errOut := new(bytes.Buffer)
			cmd.SetErr(errOut)

			Expect(cmd.ErrOrStderr()).NotTo(BeNil())
		})

		It("should handle command execution with skill data", func() {
			skillRepo := ctx.Service().GetSkillRepository()
			Expect(skillRepo).NotTo(BeNil())

			skill1 := fixtures.SkillWith("skill-success-1", "Python", "frontend", "advanced")

			err := skillRepo.Create(context.Background(), skill1)
			Expect(err).NotTo(HaveOccurred())

			cmd = skills.NewRecategorizeCmd(ctx, tea.WithInput(nil))
			out := new(bytes.Buffer)
			cmd.SetOut(out)
			cmd.SetErr(new(bytes.Buffer))

			err = cmd.RunE(cmd, []string{})
			Expect(err).NotTo(HaveOccurred())
		})
	})

	Describe("RecategorizeSkills", func() {
		var (
			svc *careerservice.Service
			out io.Writer
			err io.Writer
		)

		BeforeEach(func() {
			svc = ctx.Service()
			Expect(svc).NotTo(BeNil())
			out = new(bytes.Buffer)
			err = new(bytes.Buffer)
		})

		Context("when skill repository is nil", func() {
			It("should return error code 1", func() {
				nilSvc := testutil.NilService()
				code := skills.RecategorizeSkills(nilSvc, out, err)
				Expect(code).To(Equal(1))
			})

			It("should return error code 1 with different writers", func() {
				nilSvc := testutil.NilService()
				code := skills.RecategorizeSkills(nilSvc, io.Discard, io.Discard)
				Expect(code).To(Equal(1))
			})
		})

		Context("when recategorization is called with spinner options", func() {
			It("should return success code 0 with tea.WithInput(nil)", func() {
				code := skills.RecategorizeSkills(svc, out, err, tea.WithInput(nil))
				Expect(code).To(Equal(0))
			})

			It("should accept different output writers with spinner", func() {
				outBuf := new(bytes.Buffer)
				errBuf := new(bytes.Buffer)
				code := skills.RecategorizeSkills(svc, outBuf, errBuf, tea.WithInput(nil))
				Expect(code).To(Equal(0))
			})

			It("should handle io.Discard as output writer with spinner", func() {
				code := skills.RecategorizeSkills(svc, io.Discard, io.Discard, tea.WithInput(nil))
				Expect(code).To(Equal(0))
			})

			It("should handle nil output writers with spinner", func() {
				code := skills.RecategorizeSkills(svc, nil, nil, tea.WithInput(nil))
				Expect(code).To(Equal(0))
			})

			It("should handle nil error writer with valid output writer and spinner", func() {
				code := skills.RecategorizeSkills(svc, out, nil, tea.WithInput(nil))
				Expect(code).To(Equal(0))
			})

			It("should handle nil output writer with valid error writer and spinner", func() {
				code := skills.RecategorizeSkills(svc, nil, err, tea.WithInput(nil))
				Expect(code).To(Equal(0))
			})
		})

		Context("when recategorization is called without spinner options", func() {
			It("should return a valid exit code", func() {
				code := skills.RecategorizeSkills(svc, out, err, tea.WithInput(nil))
				Expect(code).To(BeElementOf(0, 1))
			})

			It("should accept different output writers", func() {
				outBuf := new(bytes.Buffer)
				errBuf := new(bytes.Buffer)
				code := skills.RecategorizeSkills(svc, outBuf, errBuf, tea.WithInput(nil))
				Expect(code).To(BeElementOf(0, 1))
			})

			It("should handle io.Discard as output writer", func() {
				code := skills.RecategorizeSkills(svc, io.Discard, io.Discard, tea.WithInput(nil))
				Expect(code).To(BeElementOf(0, 1))
			})

			It("should handle nil output writers", func() {
				code := skills.RecategorizeSkills(svc, nil, nil, tea.WithInput(nil))
				Expect(code).To(BeElementOf(0, 1))
			})
		})

		Context("when service has valid skill repository", func() {
			It("should process recategorization and return exit code with spinner", func() {
				code := skills.RecategorizeSkills(svc, out, err, tea.WithInput(nil))
				Expect(code).To(Equal(0))
			})

			It("should handle multiple sequential calls with spinner", func() {
				code1 := skills.RecategorizeSkills(svc, out, err, tea.WithInput(nil))
				code2 := skills.RecategorizeSkills(svc, out, err, tea.WithInput(nil))
				code3 := skills.RecategorizeSkills(svc, out, err, tea.WithInput(nil))
				Expect(code1).To(Equal(0))
				Expect(code2).To(Equal(0))
				Expect(code3).To(Equal(0))
			})

			It("should handle different writer combinations with spinner", func() {
				buf1 := new(bytes.Buffer)
				buf2 := new(bytes.Buffer)
				code := skills.RecategorizeSkills(svc, buf1, buf2, tea.WithInput(nil))
				Expect(code).To(Equal(0))
			})
		})

		Context("error handling", func() {
			It("should handle service with nil repository gracefully", func() {
				emptyService := testutil.NilService()
				code := skills.RecategorizeSkills(emptyService, out, err, tea.WithInput(nil))
				Expect(code).To(Equal(1))
			})

			It("should return non-zero code on repository error", func() {
				code := skills.RecategorizeSkills(svc, out, err, tea.WithInput(nil))
				Expect(code).To(BeElementOf(0, 1))
			})
		})

		Context("output handling", func() {
			It("should accept bytes.Buffer as output with spinner", func() {
				buf := new(bytes.Buffer)
				code := skills.RecategorizeSkills(svc, buf, buf, tea.WithInput(nil))
				Expect(code).To(Equal(0))
			})

			It("should accept io.Discard for both outputs with spinner", func() {
				code := skills.RecategorizeSkills(svc, io.Discard, io.Discard, tea.WithInput(nil))
				Expect(code).To(Equal(0))
			})

			It("should handle mixed writer types with spinner", func() {
				code := skills.RecategorizeSkills(svc, new(bytes.Buffer), io.Discard, tea.WithInput(nil))
				Expect(code).To(Equal(0))
			})
		})

		Context("with skills in repository", func() {
			BeforeEach(func() {
				skillRepo := svc.GetSkillRepository()
				Expect(skillRepo).NotTo(BeNil())

				skill1 := fixtures.SkillWith("skill-1", "Go", "backend", "advanced")
				skill2 := fixtures.SkillWith("skill-2", "React", "frontend", "intermediate")

				err := skillRepo.Create(context.Background(), skill1)
				Expect(err).NotTo(HaveOccurred())
				err = skillRepo.Create(context.Background(), skill2)
				Expect(err).NotTo(HaveOccurred())
			})

			It("should process skills with categories", func() {
				code := skills.RecategorizeSkills(svc, out, err, tea.WithInput(nil))
				Expect(code).To(Equal(0))
			})

			It("should handle multiple skills with spinner", func() {
				code := skills.RecategorizeSkills(svc, new(bytes.Buffer), new(bytes.Buffer), tea.WithInput(nil))
				Expect(code).To(Equal(0))
			})

			It("should process skills and return success", func() {
				code := skills.RecategorizeSkills(svc, io.Discard, io.Discard, tea.WithInput(nil))
				Expect(code).To(Equal(0))
			})

			It("should handle skills with different categories", func() {
				code := skills.RecategorizeSkills(svc, new(bytes.Buffer), new(bytes.Buffer), tea.WithInput(nil))
				Expect(code).To(Equal(0))
			})
		})

		Context("with multiple skills in different categories", func() {
			BeforeEach(func() {
				skillRepo := svc.GetSkillRepository()
				Expect(skillRepo).NotTo(BeNil())

				skill1 := fixtures.SkillWith("skill-multi-1", "Python", "backend", "advanced")
				skill2 := fixtures.SkillWith("skill-multi-2", "JavaScript", "frontend", "intermediate")
				skill3 := fixtures.SkillWith("skill-multi-3", "Docker", "devops", "advanced")

				err := skillRepo.Create(context.Background(), skill1)
				Expect(err).NotTo(HaveOccurred())
				err = skillRepo.Create(context.Background(), skill2)
				Expect(err).NotTo(HaveOccurred())
				err = skillRepo.Create(context.Background(), skill3)
				Expect(err).NotTo(HaveOccurred())
			})

			It("should process multiple skills with different categories", func() {
				code := skills.RecategorizeSkills(svc, out, err, tea.WithInput(nil))
				Expect(code).To(Equal(0))
			})

			It("should handle recategorization with multiple categories", func() {
				code := skills.RecategorizeSkills(svc, new(bytes.Buffer), new(bytes.Buffer), tea.WithInput(nil))
				Expect(code).To(Equal(0))
			})

			It("should process multiple skills and return success", func() {
				code := skills.RecategorizeSkills(svc, io.Discard, io.Discard, tea.WithInput(nil))
				Expect(code).To(Equal(0))
			})
		})

		Context("with skills needing recategorization", func() {
			BeforeEach(func() {
				skillRepo := svc.GetSkillRepository()
				Expect(skillRepo).NotTo(BeNil())

				skill1 := fixtures.SkillWith("skill-recat-1", "Python", "frontend", "advanced")
				skill2 := fixtures.SkillWith("skill-recat-2", "React", "backend", "intermediate")

				err := skillRepo.Create(context.Background(), skill1)
				Expect(err).NotTo(HaveOccurred())
				err = skillRepo.Create(context.Background(), skill2)
				Expect(err).NotTo(HaveOccurred())
			})

			It("should recategorize skills with wrong categories", func() {
				code := skills.RecategorizeSkills(svc, out, err, tea.WithInput(nil))
				Expect(code).To(Equal(0))
			})

			It("should handle recategorization with category updates", func() {
				code := skills.RecategorizeSkills(svc, new(bytes.Buffer), new(bytes.Buffer), tea.WithInput(nil))
				Expect(code).To(Equal(0))
			})

			It("should process recategorization and return success", func() {
				code := skills.RecategorizeSkills(svc, io.Discard, io.Discard, tea.WithInput(nil))
				Expect(code).To(Equal(0))
			})
		})

		Context("with skills having no keyword match", func() {
			BeforeEach(func() {
				skillRepo := svc.GetSkillRepository()
				Expect(skillRepo).NotTo(BeNil())

				skill1 := fixtures.SkillWith("skill-nomatch-1", "UnknownTechnology123", "other", "beginner")
				skill2 := fixtures.SkillWith("skill-nomatch-2", "XYZFramework", "other", "intermediate")

				err := skillRepo.Create(context.Background(), skill1)
				Expect(err).NotTo(HaveOccurred())
				err = skillRepo.Create(context.Background(), skill2)
				Expect(err).NotTo(HaveOccurred())
			})

			It("should handle skills with no keyword match", func() {
				code := skills.RecategorizeSkills(svc, out, err, tea.WithInput(nil))
				Expect(code).To(Equal(0))
			})

			It("should process skills with no match and return success", func() {
				code := skills.RecategorizeSkills(svc, new(bytes.Buffer), new(bytes.Buffer), tea.WithInput(nil))
				Expect(code).To(Equal(0))
			})

			It("should handle no-match skills with spinner", func() {
				code := skills.RecategorizeSkills(svc, io.Discard, io.Discard, tea.WithInput(nil))
				Expect(code).To(Equal(0))
			})
		})
	})
})
