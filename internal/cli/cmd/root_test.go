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
		It("should show help when stdin is not a TTY (non-interactive environment)", func() {
			// In tests (non-TTY stdin), running with no args shows help
			// In real interactive terminal, it would launch TUI
			// TTY detection is performed via isatty.IsTerminal(os.Stdin.Fd())
			cmd = NewRootCmd("test-version")
			cmd.SetOut(out)
			cmd.SetErr(out)
			cmd.SetArgs([]string{})

			err := cmd.Execute()
			Expect(err).NotTo(HaveOccurred())
			Expect(out.String()).To(ContainSubstring("KaRiya is a terminal user interface"))
			Expect(out.String()).To(ContainSubstring("Usage:"))
		})

		It("should return help output when no subcommand provided", func() {
			cmd = NewRootCmd("1.0.0")
			cmd.SetOut(out)
			cmd.SetErr(out)
			cmd.SetArgs([]string{})

			err := cmd.Execute()
			Expect(err).NotTo(HaveOccurred())
			Expect(out.String()).To(ContainSubstring("Available Commands:"))
		})
	})

	Context("persistent flags", func() {
		It("should accept --db flag", func() {
			cmd = NewRootCmd("test-version")
			cmd.SetOut(out)
			cmd.SetErr(out)
			cmd.SetArgs([]string{"--db", "/tmp/test.db", "--help"})

			err := cmd.Execute()
			Expect(err).NotTo(HaveOccurred())
		})

		It("should accept --in-memory flag", func() {
			cmd = NewRootCmd("test-version")
			cmd.SetOut(out)
			cmd.SetErr(out)
			cmd.SetArgs([]string{"--in-memory", "--help"})

			err := cmd.Execute()
			Expect(err).NotTo(HaveOccurred())
		})

		It("should accept both --db and --in-memory flags", func() {
			cmd = NewRootCmd("test-version")
			cmd.SetOut(out)
			cmd.SetErr(out)
			cmd.SetArgs([]string{"--db", "/tmp/test.db", "--in-memory", "--help"})

			err := cmd.Execute()
			Expect(err).NotTo(HaveOccurred())
		})
	})

	Context("version template", func() {
		It("should use custom version format", func() {
			cmd = NewRootCmd("1.2.3")
			cmd.SetOut(out)
			cmd.SetErr(out)
			cmd.SetArgs([]string{"--version"})

			err := cmd.Execute()
			Expect(err).NotTo(HaveOccurred())
			output := out.String()
			Expect(output).To(ContainSubstring("kariya version 1.2.3"))
		})
	})

	Context("subcommand discovery", func() {
		It("should show available subcommands in help", func() {
			cmd = NewRootCmd("test-version")
			cmd.SetOut(out)
			cmd.SetErr(out)
			cmd.SetArgs([]string{"--help"})

			err := cmd.Execute()
			Expect(err).NotTo(HaveOccurred())
			output := out.String()
			Expect(output).To(ContainSubstring("bursts"))
			Expect(output).To(ContainSubstring("facts"))
			Expect(output).To(ContainSubstring("skills"))
			Expect(output).To(ContainSubstring("import"))
		})
	})
})
