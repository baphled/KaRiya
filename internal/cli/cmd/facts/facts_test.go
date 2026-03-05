package facts_test

import (
	"bytes"

	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
	"github.com/spf13/cobra"

	cmdpkg "github.com/baphled/kariya/internal/cli/cmd"
	"github.com/baphled/kariya/internal/cli/cmd/cliutil"
	"github.com/baphled/kariya/internal/cli/cmd/facts"
)

var _ = Describe("Facts Command", func() {
	var (
		ctx cliutil.ServiceContext
		cmd *cobra.Command
	)

	BeforeEach(func() {
		cliCtx := cmdpkg.NewCLIContext("", true)
		err := cliCtx.InitService(new(bytes.Buffer))
		Expect(err).NotTo(HaveOccurred())
		ctx = cliCtx
	})

	Describe("NewFactsCmd", func() {
		It("should create command with correct properties", func() {
			cmd = facts.NewFactsCmd(ctx)
			Expect(cmd).NotTo(BeNil())
			Expect(cmd.Use).To(Equal("facts"))
			Expect(cmd.Short).To(Equal("Manage facts"))
			Expect(cmd.Long).To(ContainSubstring("facts from events"))
		})

		It("should have extract subcommand", func() {
			cmd = facts.NewFactsCmd(ctx)
			extractCmd, _, err := cmd.Find([]string{"extract"})
			Expect(err).NotTo(HaveOccurred())
			Expect(extractCmd).NotTo(BeNil())
			Expect(extractCmd.Use).To(Equal("extract"))
		})

		It("should have list subcommand", func() {
			cmd = facts.NewFactsCmd(ctx)
			listCmd, _, err := cmd.Find([]string{"list"})
			Expect(err).NotTo(HaveOccurred())
			Expect(listCmd).NotTo(BeNil())
			Expect(listCmd.Use).To(Equal("list"))
		})
	})

	Describe("NewExtractCmd", func() {
		It("should create extract command with correct properties", func() {
			cmd = facts.NewExtractCmd(ctx)
			Expect(cmd).NotTo(BeNil())
			Expect(cmd.Use).To(Equal("extract"))
			Expect(cmd.Short).To(Equal("Extract facts from events"))
			Expect(cmd.Long).To(ContainSubstring("facts"))
		})

		It("should execute extract command successfully", func() {
			cmd = facts.NewExtractCmd(ctx)
			out := new(bytes.Buffer)
			cmd.SetOut(out)
			cmd.SetErr(new(bytes.Buffer))

			err := cmd.RunE(cmd, []string{})
			Expect(err).NotTo(HaveOccurred())
		})

		It("should return error when service is not initialized", func() {
			nilCtx := cmdpkg.NewCLIContext("", true)
			cobraCmd := facts.NewExtractCmd(nilCtx)
			out := new(bytes.Buffer)
			cobraCmd.SetOut(out)
			cobraCmd.SetErr(new(bytes.Buffer))

			err := cobraCmd.RunE(cobraCmd, []string{})
			Expect(err).To(HaveOccurred())
			Expect(err.Error()).To(ContainSubstring("service not initialized"))
		})
	})

	Describe("NewListCmd", func() {
		It("should create list command with correct properties", func() {
			cmd = facts.NewListCmd(ctx)
			Expect(cmd).NotTo(BeNil())
			Expect(cmd.Use).To(Equal("list"))
			Expect(cmd.Short).To(Equal("List all facts"))
			Expect(cmd.Long).To(ContainSubstring("existing facts"))
		})

		It("should execute list command successfully", func() {
			cmd = facts.NewListCmd(ctx)
			out := new(bytes.Buffer)
			cmd.SetOut(out)
			cmd.SetErr(new(bytes.Buffer))

			err := cmd.RunE(cmd, []string{})
			Expect(err).NotTo(HaveOccurred())
		})

		It("should return error when service is not initialized", func() {
			nilCtx := cmdpkg.NewCLIContext("", true)
			cobraCmd := facts.NewListCmd(nilCtx)
			out := new(bytes.Buffer)
			cobraCmd.SetOut(out)
			cobraCmd.SetErr(new(bytes.Buffer))

			err := cobraCmd.RunE(cobraCmd, []string{})
			Expect(err).To(HaveOccurred())
			Expect(err.Error()).To(ContainSubstring("service not initialized"))
		})
	})
})
