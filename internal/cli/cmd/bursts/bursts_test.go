package bursts_test

import (
	"bytes"

	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
	"github.com/spf13/cobra"

	cmdpkg "github.com/baphled/kariya/internal/cli/cmd"
	"github.com/baphled/kariya/internal/cli/cmd/bursts"
	"github.com/baphled/kariya/internal/cli/cmd/cliutil"
)

var _ = Describe("Bursts Command", func() {
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

	Describe("NewBurstsCmd", func() {
		It("should create command with correct properties", func() {
			cmd = bursts.NewBurstsCmd(ctx)
			Expect(cmd).NotTo(BeNil())
			Expect(cmd.Use).To(Equal("bursts"))
			Expect(cmd.Short).To(Equal("Manage bursts"))
			Expect(cmd.Long).To(ContainSubstring("burst patterns"))
		})

		It("should have detect subcommand", func() {
			cmd = bursts.NewBurstsCmd(ctx)
			detectCmd, _, err := cmd.Find([]string{"detect"})
			Expect(err).NotTo(HaveOccurred())
			Expect(detectCmd).NotTo(BeNil())
			Expect(detectCmd.Use).To(Equal("detect"))
		})

		It("should have list subcommand", func() {
			cmd = bursts.NewBurstsCmd(ctx)
			listCmd, _, err := cmd.Find([]string{"list"})
			Expect(err).NotTo(HaveOccurred())
			Expect(listCmd).NotTo(BeNil())
			Expect(listCmd.Use).To(Equal("list"))
		})
	})

	Describe("NewDetectCmd", func() {
		It("should create detect command with correct properties", func() {
			cmd = bursts.NewDetectCmd(ctx)
			Expect(cmd).NotTo(BeNil())
			Expect(cmd.Use).To(Equal("detect"))
			Expect(cmd.Short).To(Equal("Detect burst patterns in events"))
			Expect(cmd.Long).To(ContainSubstring("burst patterns"))
		})

		It("should execute detect command successfully", func() {
			cmd = bursts.NewDetectCmd(ctx)
			out := new(bytes.Buffer)
			cmd.SetOut(out)
			cmd.SetErr(new(bytes.Buffer))

			err := cmd.RunE(cmd, []string{})
			Expect(err).NotTo(HaveOccurred())
		})

		It("should return error when service is not initialized", func() {
			nilCtx := cmdpkg.NewCLIContext("", true)
			cobraCmd := bursts.NewDetectCmd(nilCtx)
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
			cmd = bursts.NewListCmd(ctx)
			Expect(cmd).NotTo(BeNil())
			Expect(cmd.Use).To(Equal("list"))
			Expect(cmd.Short).To(Equal("List all bursts"))
			Expect(cmd.Long).To(ContainSubstring("existing bursts"))
		})

		It("should execute list command successfully", func() {
			cmd = bursts.NewListCmd(ctx)
			out := new(bytes.Buffer)
			cmd.SetOut(out)
			cmd.SetErr(new(bytes.Buffer))

			err := cmd.RunE(cmd, []string{})
			Expect(err).NotTo(HaveOccurred())
		})

		It("should return error when service is not initialized", func() {
			nilCtx := cmdpkg.NewCLIContext("", true)
			cobraCmd := bursts.NewListCmd(nilCtx)
			out := new(bytes.Buffer)
			cobraCmd.SetOut(out)
			cobraCmd.SetErr(new(bytes.Buffer))

			err := cobraCmd.RunE(cobraCmd, []string{})
			Expect(err).To(HaveOccurred())
			Expect(err.Error()).To(ContainSubstring("service not initialized"))
		})
	})
})
