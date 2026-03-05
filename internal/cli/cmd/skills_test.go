package cmd_test

import (
	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"

	cmdpkg "github.com/baphled/kariya/internal/cli/cmd"
)

var _ = Describe("Skills Commands", func() {

	Describe("Skills command structure", func() {
		var ctx *cmdpkg.CLIContext

		BeforeEach(func() {
			ctx = &cmdpkg.CLIContext{InMemory: true}
			initErr := ctx.InitService()
			Expect(initErr).NotTo(HaveOccurred())
		})

		It("root command includes skills subcommand", func() {
			rootCmd := cmdpkg.NewRootCmd("1.0.0")
			subcommands := rootCmd.Commands()
			subcommandNames := make([]string, len(subcommands))
			for i, sub := range subcommands {
				subcommandNames[i] = sub.Name()
			}
			Expect(subcommandNames).To(ContainElement("skills"))
		})

		It("skills command has recategorize subcommand", func() {
			rootCmd := cmdpkg.NewRootCmd("1.0.0")
			skillsCmd, _, _ := rootCmd.Find([]string{"skills"})
			Expect(skillsCmd).NotTo(BeNil())
			subcommands := skillsCmd.Commands()
			subcommandNames := make([]string, len(subcommands))
			for i, sub := range subcommands {
				subcommandNames[i] = sub.Name()
			}
			Expect(subcommandNames).To(ContainElement("recategorize"))
		})

		It("recategorize subcommand has correct properties", func() {
			rootCmd := cmdpkg.NewRootCmd("1.0.0")
			skillsCmd, _, _ := rootCmd.Find([]string{"skills"})
			recatCmd, _, _ := skillsCmd.Find([]string{"recategorize"})
			Expect(recatCmd).NotTo(BeNil())
			Expect(recatCmd.Use).To(Equal("recategorize"))
			Expect(recatCmd.Short).To(ContainSubstring("Recategorize"))
		})

		It("recategorize subcommand has RunE function", func() {
			rootCmd := cmdpkg.NewRootCmd("1.0.0")
			skillsCmd, _, _ := rootCmd.Find([]string{"skills"})
			recatCmd, _, _ := skillsCmd.Find([]string{"recategorize"})
			Expect(recatCmd.RunE).NotTo(BeNil())
		})
	})
})
