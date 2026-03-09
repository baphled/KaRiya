package main

import (
	"bytes"
	"path/filepath"
	"strings"

	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"

	"github.com/baphled/kariya/internal/config"
)

var _ = Describe("CLI Entry Point", func() {
	var (
		outBuf bytes.Buffer
		errBuf bytes.Buffer
	)

	BeforeEach(func() {
		outBuf.Reset()
		errBuf.Reset()
	})

	Describe("run() function", func() {
		Context("with version flag", func() {
			It("should display version", func() {
				exitCode := run([]string{"--version"}, &outBuf, &errBuf)

				Expect(exitCode).To(Equal(0))
				Expect(outBuf.String()).To(ContainSubstring("kariya version"))
				Expect(outBuf.String()).To(ContainSubstring("dev"))
			})
		})

		Context("with help flag", func() {
			It("should display help", func() {
				exitCode := run([]string{"--help"}, &outBuf, &errBuf)

				Expect(exitCode).To(Equal(0))
				Expect(outBuf.String()).To(ContainSubstring("terminal user interface"))
				Expect(outBuf.String()).To(ContainSubstring("Usage:"))
			})

			It("should list available subcommands", func() {
				exitCode := run([]string{"--help"}, &outBuf, &errBuf)

				Expect(exitCode).To(Equal(0))
				output := outBuf.String()
				Expect(output).To(ContainSubstring("Available Commands:"))
				Expect(output).To(ContainSubstring("bursts"))
				Expect(output).To(ContainSubstring("facts"))
				Expect(output).To(ContainSubstring("skills"))
				Expect(output).To(ContainSubstring("import"))
			})

			It("should show persistent flags", func() {
				exitCode := run([]string{"--help"}, &outBuf, &errBuf)

				Expect(exitCode).To(Equal(0))
				output := outBuf.String()
				Expect(output).To(ContainSubstring("--db"))
				Expect(output).To(ContainSubstring("--in-memory"))
			})
		})

		Context("with no arguments", func() {
			It("should attempt to launch TUI", func() {
				// TUI launch requires config isolation to prevent test pollution
				tempDir := GinkgoT().TempDir()
				configPath := filepath.Join(tempDir, "config.yaml")
				config.SetConfigPathForTesting(configPath)
				defer config.ResetConfigPath()

				// TUI launch cannot be fully tested in unit tests (requires terminal)
				// This test verifies that the command attempts to initialize and launch
				// without panicking due to missing config isolation.
				// In a real scenario, the TUI would render to the terminal.
				Skip("TUI launch requires interactive terminal; verified via VHS demo and manual testing")

				exitCode := run([]string{}, &outBuf, &errBuf)
				Expect(exitCode).To(Equal(0))
			})
		})

		Context("with invalid command", func() {
			It("should return error for unknown command", func() {
				exitCode := run([]string{"unknown-command"}, &outBuf, &errBuf)

				Expect(exitCode).To(Equal(1))
				output := errBuf.String()
				Expect(output).To(Or(
					ContainSubstring("unknown command"),
					ContainSubstring("Unknown command"),
				))
			})
		})

		Context("with database flags", func() {
			It("should accept --db flag", func() {
				exitCode := run([]string{"--db", "/tmp/test.db", "--help"}, &outBuf, &errBuf)

				Expect(exitCode).To(Equal(0))
			})

			It("should accept --in-memory flag", func() {
				exitCode := run([]string{"--in-memory", "--help"}, &outBuf, &errBuf)

				Expect(exitCode).To(Equal(0))
			})

			It("should accept both flags together", func() {
				exitCode := run([]string{"--db", "/tmp/test.db", "--in-memory", "--help"}, &outBuf, &errBuf)

				Expect(exitCode).To(Equal(0))
			})
		})

		Context("output streams", func() {
			It("should write normal output to stdout", func() {
				run([]string{"--version"}, &outBuf, &errBuf)

				Expect(outBuf.Len()).To(BeNumerically(">", 0))
				Expect(errBuf.Len()).To(Equal(0))
			})

			It("should write errors to stderr", func() {
				run([]string{"invalid-command"}, &outBuf, &errBuf)

				Expect(errBuf.Len()).To(BeNumerically(">", 0))
			})
		})

		Context("subcommand help", func() {
			It("should show bursts subcommand help", func() {
				exitCode := run([]string{"bursts", "--help"}, &outBuf, &errBuf)

				Expect(exitCode).To(Equal(0))
				Expect(outBuf.String()).To(ContainSubstring("bursts"))
			})

			It("should show facts subcommand help", func() {
				exitCode := run([]string{"facts", "--help"}, &outBuf, &errBuf)

				Expect(exitCode).To(Equal(0))
				Expect(outBuf.String()).To(ContainSubstring("facts"))
			})

			It("should show skills subcommand help", func() {
				exitCode := run([]string{"skills", "--help"}, &outBuf, &errBuf)

				Expect(exitCode).To(Equal(0))
				Expect(outBuf.String()).To(ContainSubstring("skills"))
			})

			It("should show import subcommand help", func() {
				exitCode := run([]string{"import", "--help"}, &outBuf, &errBuf)

				Expect(exitCode).To(Equal(0))
				Expect(outBuf.String()).To(ContainSubstring("import"))
			})
		})

		Context("version template", func() {
			It("should use custom version format", func() {
				exitCode := run([]string{"--version"}, &outBuf, &errBuf)

				Expect(exitCode).To(Equal(0))
				output := strings.TrimSpace(outBuf.String())
				Expect(output).To(MatchRegexp(`^kariya version .+$`))
			})
		})

		Context("exit codes", func() {
			It("should return 0 on success", func() {
				exitCode := run([]string{"--version"}, &outBuf, &errBuf)
				Expect(exitCode).To(Equal(0))
			})

			It("should return 1 on error", func() {
				exitCode := run([]string{"invalid-command"}, &outBuf, &errBuf)
				Expect(exitCode).To(Equal(1))
			})
		})

		Context("args parsing", func() {
			It("should accept empty args slice", func() {
				// Empty args should return 0 (help shown in non-TTY)
				exitCode := run([]string{}, &outBuf, &errBuf)
				Expect(exitCode).To(Equal(0))
			})

			It("should handle single arg", func() {
				exitCode := run([]string{"--version"}, &outBuf, &errBuf)
				Expect(exitCode).To(Equal(0))
			})

			It("should handle multiple args", func() {
				exitCode := run([]string{"--db", "/tmp/test.db", "--help"}, &outBuf, &errBuf)
				Expect(exitCode).To(Equal(0))
			})
		})
	})

	Describe("main() function", func() {
		It("should be defined", func() {
			Expect(main).NotTo(BeNil())
		})
	})

	Describe("version variable", func() {
		It("should have default value", func() {
			Expect(version).To(Equal("dev"))
		})
	})
})
