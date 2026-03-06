package cmd

import (
	"bytes"
	"io"

	"github.com/baphled/kariya/internal/cli/cmd/cliutil"
	careerservice "github.com/baphled/kariya/internal/service/career"
	tea "github.com/charmbracelet/bubbletea"
	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
)

var _ = Describe("CommandFactory", func() {
	var ctx *CLIContext

	BeforeEach(func() {
		ctx = NewCLIContext("", true)
		err := ctx.InitService(new(bytes.Buffer))
		Expect(err).NotTo(HaveOccurred())
	})

	Describe("NewTUICommand", func() {
		It("should create command with correct properties", func() {
			cfg := CommandConfig{
				Name:  "test-cmd",
				Short: "Test command",
				Long:  "A test command for testing",
				Action: func(_ *careerservice.Service, _, _ io.Writer, _ cliutil.ProgressRunner, _ ...tea.ProgramOption) int {
					return 0
				},
				FailMsg: "test failed",
			}
			cmd := NewTUICommand(ctx, cfg)
			Expect(cmd).NotTo(BeNil())
			Expect(cmd.Use).To(Equal("test-cmd"))
			Expect(cmd.Short).To(Equal("Test command"))
			Expect(cmd.Long).To(Equal("A test command for testing"))
		})

		It("should return nil when action succeeds", func() {
			cfg := CommandConfig{
				Name:  "test",
				Short: "Test",
				Long:  "Test",
				Action: func(_ *careerservice.Service, _, _ io.Writer, _ cliutil.ProgressRunner, _ ...tea.ProgramOption) int {
					return 0
				},
				FailMsg: "failed",
			}
			cmd := NewTUICommand(ctx, cfg)
			cmd.SetOut(new(bytes.Buffer))
			cmd.SetErr(new(bytes.Buffer))

			err := cmd.RunE(cmd, []string{})
			Expect(err).NotTo(HaveOccurred())
		})

		It("should return error when action returns non-zero", func() {
			cfg := CommandConfig{
				Name:  "test",
				Short: "Test",
				Long:  "Test",
				Action: func(_ *careerservice.Service, _, _ io.Writer, _ cliutil.ProgressRunner, _ ...tea.ProgramOption) int {
					return 1
				},
				FailMsg: "operation failed",
			}
			cmd := NewTUICommand(ctx, cfg)
			cmd.SetOut(new(bytes.Buffer))
			cmd.SetErr(new(bytes.Buffer))

			err := cmd.RunE(cmd, []string{})
			Expect(err).To(HaveOccurred())
			Expect(err.Error()).To(Equal("operation failed"))
		})

		It("should add output discard option when output is io.Discard", func() {
			actionCalled := false
			cfg := CommandConfig{
				Name:  "test",
				Short: "Test",
				Long:  "Test",
				Action: func(_ *careerservice.Service, _, _ io.Writer, _ cliutil.ProgressRunner, _ ...tea.ProgramOption) int {
					actionCalled = true
					return 0
				},
				FailMsg: "failed",
			}
			cmd := NewTUICommand(ctx, cfg)
			cmd.SetOut(io.Discard)
			cmd.SetErr(new(bytes.Buffer))

			err := cmd.RunE(cmd, []string{})
			Expect(err).NotTo(HaveOccurred())
			Expect(actionCalled).To(BeTrue())
		})
	})

	Describe("NewListCommand", func() {
		It("should create command with correct properties", func() {
			cfg := ListConfig{
				Name:  "list-test",
				Short: "List test",
				Long:  "A list test command",
				Action: func(_ *careerservice.Service, _, _ io.Writer) int {
					return 0
				},
				FailMsg: "list failed",
			}
			cmd := NewListCommand(ctx, cfg)
			Expect(cmd).NotTo(BeNil())
			Expect(cmd.Use).To(Equal("list-test"))
			Expect(cmd.Short).To(Equal("List test"))
			Expect(cmd.Long).To(Equal("A list test command"))
		})

		It("should return nil when action succeeds", func() {
			cfg := ListConfig{
				Name:  "list",
				Short: "List",
				Long:  "List",
				Action: func(_ *careerservice.Service, _, _ io.Writer) int {
					return 0
				},
				FailMsg: "list failed",
			}
			cmd := NewListCommand(ctx, cfg)
			cmd.SetOut(new(bytes.Buffer))
			cmd.SetErr(new(bytes.Buffer))

			err := cmd.RunE(cmd, []string{})
			Expect(err).NotTo(HaveOccurred())
		})

		It("should return error when action returns non-zero", func() {
			cfg := ListConfig{
				Name:  "list",
				Short: "List",
				Long:  "List",
				Action: func(_ *careerservice.Service, _, _ io.Writer) int {
					return 1
				},
				FailMsg: "listing failed",
			}
			cmd := NewListCommand(ctx, cfg)
			cmd.SetOut(new(bytes.Buffer))
			cmd.SetErr(new(bytes.Buffer))

			err := cmd.RunE(cmd, []string{})
			Expect(err).To(HaveOccurred())
			Expect(err.Error()).To(Equal("listing failed"))
		})
	})
})
