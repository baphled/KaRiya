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

		It("should format version correctly", func() {
			cmd = NewRootCmd("v1.2.3")
			cmd.SetOut(out)
			cmd.SetErr(out)
			cmd.SetArgs([]string{"--version"})

			err := cmd.Execute()
			Expect(err).NotTo(HaveOccurred())
			Expect(out.String()).To(Equal("kariya version v1.2.3\n"))
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

		It("should show usage", func() {
			cmd = NewRootCmd("test-version")
			cmd.SetOut(out)
			cmd.SetErr(out)
			cmd.SetArgs([]string{"--help"})

			err := cmd.Execute()
			Expect(err).NotTo(HaveOccurred())
			Expect(out.String()).To(ContainSubstring("Usage:"))
		})

		It("should show available commands", func() {
			cmd = NewRootCmd("test-version")
			cmd.SetOut(out)
			cmd.SetErr(out)
			cmd.SetArgs([]string{"--help"})

			err := cmd.Execute()
			Expect(err).NotTo(HaveOccurred())
			output := out.String()
			Expect(output).To(ContainSubstring("Available Commands:"))
			Expect(output).To(ContainSubstring("bursts"))
			Expect(output).To(ContainSubstring("facts"))
			Expect(output).To(ContainSubstring("skills"))
			Expect(output).To(ContainSubstring("import"))
		})
	})

	Context("no arguments in non-TTY environment", func() {
		It("should show help when no args provided", func() {
			cmd = NewRootCmd("test-version")
			cmd.SetOut(out)
			cmd.SetErr(out)
			cmd.SetArgs([]string{"--in-memory"})

			err := cmd.Execute()
			Expect(err).NotTo(HaveOccurred())
			Expect(out.String()).To(ContainSubstring("KaRiya is a terminal user interface"))
			Expect(out.String()).To(ContainSubstring("Usage:"))
		})

		It("should show available commands without args", func() {
			cmd = NewRootCmd("test-version")
			cmd.SetOut(out)
			cmd.SetErr(out)
			cmd.SetArgs([]string{"--in-memory"})

			err := cmd.Execute()
			Expect(err).NotTo(HaveOccurred())
			Expect(out.String()).To(ContainSubstring("Available Commands:"))
		})

		It("should output to stdout", func() {
			cmd = NewRootCmd("test-version")
			cmd.SetOut(out)
			cmd.SetErr(out)
			cmd.SetArgs([]string{"--in-memory"})

			err := cmd.Execute()
			Expect(err).NotTo(HaveOccurred())
			Expect(out.Len()).To(BeNumerically(">", 0))
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

		It("should accept both flags together", func() {
			cmd = NewRootCmd("test-version")
			cmd.SetOut(out)
			cmd.SetErr(out)
			cmd.SetArgs([]string{"--db", "/tmp/test.db", "--in-memory", "--help"})

			err := cmd.Execute()
			Expect(err).NotTo(HaveOccurred())
		})

		It("should accept --db with equals syntax", func() {
			cmd = NewRootCmd("test-version")
			cmd.SetOut(out)
			cmd.SetErr(out)
			cmd.SetArgs([]string{"--db=/tmp/test.db", "--help"})

			err := cmd.Execute()
			Expect(err).NotTo(HaveOccurred())
		})

		It("should handle empty --db value", func() {
			cmd = NewRootCmd("test-version")
			cmd.SetOut(out)
			cmd.SetErr(out)
			cmd.SetArgs([]string{"--db", "", "--help"})

			err := cmd.Execute()
			Expect(err).NotTo(HaveOccurred())
		})
	})

	Context("version template", func() {
		It("should use custom version format", func() {
			cmd = NewRootCmd("0.1.0")
			cmd.SetOut(out)
			cmd.SetErr(out)
			cmd.SetArgs([]string{"--version"})

			err := cmd.Execute()
			Expect(err).NotTo(HaveOccurred())
			output := out.String()
			Expect(output).To(ContainSubstring("kariya version 0.1.0"))
		})

		It("should only output version without extra text", func() {
			cmd = NewRootCmd("1.0.0")
			cmd.SetOut(out)
			cmd.SetErr(out)
			cmd.SetArgs([]string{"--version"})

			err := cmd.Execute()
			Expect(err).NotTo(HaveOccurred())
			output := out.String()
			Expect(output).To(Equal("kariya version 1.0.0\n"))
		})
	})

	Context("command metadata", func() {
		It("should have correct Use name", func() {
			cmd = NewRootCmd("test")
			Expect(cmd.Use).To(Equal("kariya"))
		})

		It("should have short description", func() {
			cmd = NewRootCmd("test")
			Expect(cmd.Short).To(Equal("Career Event Capture for Engineers"))
		})

		It("should have long description", func() {
			cmd = NewRootCmd("test")
			Expect(cmd.Long).To(ContainSubstring("terminal user interface"))
		})

		It("should set version string", func() {
			cmd = NewRootCmd("v2.0.0")
			Expect(cmd.Version).To(Equal("v2.0.0"))
		})

		It("should silence usage on error", func() {
			cmd = NewRootCmd("test")
			Expect(cmd.SilenceUsage).To(BeTrue())
		})
	})

	Context("flag documentation in help", func() {
		It("should document --db flag in help", func() {
			cmd = NewRootCmd("test-version")
			cmd.SetOut(out)
			cmd.SetErr(out)
			cmd.SetArgs([]string{"--help"})

			err := cmd.Execute()
			Expect(err).NotTo(HaveOccurred())
			Expect(out.String()).To(ContainSubstring("--db"))
		})

		It("should document --in-memory flag in help", func() {
			cmd = NewRootCmd("test-version")
			cmd.SetOut(out)
			cmd.SetErr(out)
			cmd.SetArgs([]string{"--help"})

			err := cmd.Execute()
			Expect(err).NotTo(HaveOccurred())
			Expect(out.String()).To(ContainSubstring("--in-memory"))
		})

		It("should show subcommand help in main help", func() {
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
		})
	})

	Context("error handling and edge cases", func() {
		It("should handle multiple flag calls", func() {
			cmd = NewRootCmd("test-version")
			cmd.SetOut(out)
			cmd.SetErr(out)
			cmd.SetArgs([]string{"--db", "/path1", "--db", "/path2", "--help"})

			Expect(func() {
				_ = cmd.Execute()
			}).NotTo(Panic())
		})

		It("should execute without panic on empty args", func() {
			cmd = NewRootCmd("test-version")
			cmd.SetOut(out)
			cmd.SetErr(out)
			cmd.SetArgs([]string{"--in-memory"})

			Expect(func() {
				_ = cmd.Execute()
			}).NotTo(Panic())
		})

		It("should handle version flag properly", func() {
			cmd = NewRootCmd("test-version")
			cmd.SetOut(out)
			cmd.SetErr(out)
			cmd.SetArgs([]string{"--version"})

			err := cmd.Execute()
			Expect(err).NotTo(HaveOccurred())
		})

		It("should handle help flag properly", func() {
			cmd = NewRootCmd("test-version")
			cmd.SetOut(out)
			cmd.SetErr(out)
			cmd.SetArgs([]string{"--help"})

			err := cmd.Execute()
			Expect(err).NotTo(HaveOccurred())
		})
	})

	Context("output verification", func() {
		It("should write help to stdout not stderr", func() {
			cmd = NewRootCmd("test-version")
			cmd.SetOut(out)
			errBuf := new(bytes.Buffer)
			cmd.SetErr(errBuf)
			cmd.SetArgs([]string{"--help"})

			err := cmd.Execute()
			Expect(err).NotTo(HaveOccurred())
			Expect(out.String()).To(ContainSubstring("KaRiya"))
			Expect(errBuf.String()).To(BeEmpty())
		})

		It("should write version to stdout", func() {
			cmd = NewRootCmd("test-version")
			cmd.SetOut(out)
			errBuf := new(bytes.Buffer)
			cmd.SetErr(errBuf)
			cmd.SetArgs([]string{"--version"})

			err := cmd.Execute()
			Expect(err).NotTo(HaveOccurred())
			Expect(out.String()).To(ContainSubstring("kariya"))
		})
	})

	Context("TTY detection behavior", func() {
		It("should check TTY before attempting TUI launch", func() {
			cmd = NewRootCmd("test-version")
			cmd.SetOut(out)
			cmd.SetErr(out)
			cmd.SetArgs([]string{"--in-memory"})

			err := cmd.Execute()
			Expect(err).NotTo(HaveOccurred())
			Expect(out.String()).To(ContainSubstring("KaRiya is a terminal user interface"))
		})

		It("should show help consistently in non-TTY", func() {
			cmd = NewRootCmd("test-version")
			cmd.SetOut(out)
			cmd.SetErr(out)
			cmd.SetArgs([]string{"--in-memory"})

			err := cmd.Execute()
			Expect(err).NotTo(HaveOccurred())

			output := out.String()
			Expect(output).To(ContainSubstring("Usage:"))
			Expect(output).To(ContainSubstring("kariya"))
		})
	})
})
