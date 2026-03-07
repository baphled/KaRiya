package cmd

import (
	"bytes"

	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
	"github.com/spf13/cobra"
)

var _ = Describe("Root Command", func() {
	var (
		cmd *cobra.Command
		out *bytes.Buffer
	)

	BeforeEach(func() {
		out = new(bytes.Buffer)
	})

	Context("version flag", func() {
		It("should display version with long flag", func() {
			cmd = NewRootCmd("test-version")
			cmd.SetOut(out)
			cmd.SetErr(out)
			cmd.SetArgs([]string{"--version"})

			err := cmd.Execute()
			Expect(err).NotTo(HaveOccurred())
			Expect(out.String()).To(ContainSubstring("kariya version"))
		})

		It("should display version with short flag", func() {
			cmd = NewRootCmd("test-version")
			cmd.SetOut(out)
			cmd.SetErr(out)
			cmd.SetArgs([]string{"-v"})

			err := cmd.Execute()
			Expect(err).NotTo(HaveOccurred())
			Expect(out.String()).To(ContainSubstring("kariya version"))
		})
	})

	Context("help flag", func() {
		It("should display help with long flag", func() {
			cmd = NewRootCmd("test-version")
			cmd.SetOut(out)
			cmd.SetErr(out)
			cmd.SetArgs([]string{"--help"})

			err := cmd.Execute()
			Expect(err).NotTo(HaveOccurred())
			Expect(out.String()).To(ContainSubstring("KaRiya is a terminal user interface"))
		})

		It("should display help with short flag", func() {
			cmd = NewRootCmd("test-version")
			cmd.SetOut(out)
			cmd.SetErr(out)
			cmd.SetArgs([]string{"-h"})

			err := cmd.Execute()
			Expect(err).NotTo(HaveOccurred())
			Expect(out.String()).To(ContainSubstring("KaRiya is a terminal user interface"))
		})
	})

	Context("no arguments", func() {
		It("should display help", func() {
			cmd = NewRootCmd("test-version")
			cmd.SetOut(out)
			cmd.SetErr(out)
			cmd.SetArgs([]string{})

			err := cmd.Execute()
			Expect(err).NotTo(HaveOccurred())
			Expect(out.String()).To(ContainSubstring("KaRiya is a terminal user interface"))
		})
	})
})
