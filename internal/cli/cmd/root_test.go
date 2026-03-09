package cmd

import (
	"bytes"
	"path/filepath"

	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
	"github.com/spf13/cobra"

	"github.com/baphled/kariya/internal/config"
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
		It("should attempt to launch TUI", func() {
			// TUI launch requires config isolation to prevent test pollution
			tempDir := GinkgoT().TempDir()
			configPath := filepath.Join(tempDir, "config.yaml")
			config.SetConfigPathForTesting(configPath)
			defer config.ResetConfigPath()

			cmd = NewRootCmd("test-version")
			cmd.SetOut(out)
			cmd.SetErr(out)
			cmd.SetArgs([]string{})

			// TUI launch cannot be fully tested in unit tests (requires terminal)
			// This test verifies that the command attempts to initialize and launch
			// without panicking due to missing config isolation.
			// In a real scenario, the TUI would render to the terminal.
			Skip("TUI launch requires interactive terminal; verified via VHS demo and manual testing")

			err := cmd.Execute()
			Expect(err).NotTo(HaveOccurred())
		})
	})
})
