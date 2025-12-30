package main

import (
	"bytes"

	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
)

var _ = Describe("CLI Initialization", func() {
	Context("Version Flag", func() {
		It("should print version when --version flag is provided", func() {
			var buf, errBuf bytes.Buffer
			exitCode := run([]string{"--version"}, &buf, &errBuf)

			Expect(exitCode).To(Equal(0))
			Expect(buf.String()).To(ContainSubstring("KaRiya CLI v"))
		})
	})

	Context("Help Flag", func() {
		It("should print help information when --help flag is provided", func() {
			var buf, errBuf bytes.Buffer
			exitCode := run([]string{"--help"}, &buf, &errBuf)

			Expect(exitCode).To(Equal(0))
			Expect(buf.String()).To(ContainSubstring("KaRiya CLI - Career Journaling Tool"))
			Expect(buf.String()).To(ContainSubstring("Usage: kariya [options]"))
		})

		It("should show database flag in help", func() {
			var buf, errBuf bytes.Buffer
			run([]string{"--help"}, &buf, &errBuf)

			Expect(buf.String()).To(ContainSubstring("--db"))
			Expect(buf.String()).To(ContainSubstring("--database"))
		})

		It("should show mode flag in help", func() {
			var buf, errBuf bytes.Buffer
			run([]string{"--help"}, &buf, &errBuf)

			Expect(buf.String()).To(ContainSubstring("--mode"))
		})

		It("should show list flag in help", func() {
			var buf, errBuf bytes.Buffer
			run([]string{"--help"}, &buf, &errBuf)

			Expect(buf.String()).To(ContainSubstring("--list"))
		})

		It("should show import flag in help", func() {
			var buf, errBuf bytes.Buffer
			run([]string{"--help"}, &buf, &errBuf)

			Expect(buf.String()).To(ContainSubstring("--import"))
			Expect(buf.String()).To(ContainSubstring("Import career events from CSV file"))
		})

		It("should show skip-import-review flag in help", func() {
			var buf, errBuf bytes.Buffer
			run([]string{"--help"}, &buf, &errBuf)

			Expect(buf.String()).To(ContainSubstring("--skip-import-review"))
		})

		It("should show CSV format documentation in help", func() {
			var buf, errBuf bytes.Buffer
			run([]string{"--help"}, &buf, &errBuf)

			Expect(buf.String()).To(ContainSubstring("CSV File Format"))
			Expect(buf.String()).To(ContainSubstring("Text (required)"))
			Expect(buf.String()).To(ContainSubstring("Date (required)"))
			Expect(buf.String()).To(ContainSubstring("Semicolon-separated tags"))
		})
	})

	Context("Mode Flag", func() {
		It("should reject invalid mode and show error", func() {
			var buf, errBuf bytes.Buffer
			exitCode := run([]string{"--mode", "invalid"}, &buf, &errBuf)

			Expect(exitCode).To(Equal(1))
			Expect(errBuf.String()).To(ContainSubstring("Invalid mode"))
		})
	})

	Context("Database Flag", func() {
		It("should show database path options in help", func() {
			var buf, errBuf bytes.Buffer
			run([]string{"--help"}, &buf, &errBuf)

			Expect(buf.String()).To(ContainSubstring("Use custom database path"))
			Expect(buf.String()).To(ContainSubstring("Default: ~/.kariya/events.db"))
		})
	})

	Context("List Flag", func() {
		It("should show --list flag in help", func() {
			var buf, errBuf bytes.Buffer
			run([]string{"--help"}, &buf, &errBuf)

			Expect(buf.String()).To(ContainSubstring("--list"))
			Expect(buf.String()).To(ContainSubstring("Show recent events on startup"))
		})

		It("should show example with --list and --db together", func() {
			var buf, errBuf bytes.Buffer
			run([]string{"--help"}, &buf, &errBuf)

			Expect(buf.String()).To(ContainSubstring("kariya --db ./events.db --list"))
		})
	})

	Context("Import Flag", func() {
		It("should require a file path for --import flag", func() {
			var buf, errBuf bytes.Buffer
			exitCode := run([]string{"--import"}, &buf, &errBuf)

			Expect(exitCode).To(Equal(1))
			Expect(errBuf.String()).To(ContainSubstring("--import flag requires a file path"))
		})

		It("should show import examples in help", func() {
			var buf, errBuf bytes.Buffer
			run([]string{"--help"}, &buf, &errBuf)

			Expect(buf.String()).To(ContainSubstring("kariya --import events.csv"))
			Expect(buf.String()).To(ContainSubstring("--skip-import-review"))
		})

		It("should error when import file does not exist", func() {
			var buf, errBuf bytes.Buffer
			exitCode := run([]string{"--import", "/nonexistent/file.csv"}, &buf, &errBuf)

			Expect(exitCode).To(Equal(1))
			Expect(errBuf.String()).To(ContainSubstring("Cannot access import file"))
		})
	})

	Context("Help Examples", func() {
		It("should show usage examples", func() {
			var buf, errBuf bytes.Buffer
			run([]string{"--help"}, &buf, &errBuf)

			Expect(buf.String()).To(ContainSubstring("Examples:"))
			Expect(buf.String()).To(ContainSubstring("kariya"))
			Expect(buf.String()).To(ContainSubstring("--mode timeline"))
		})
	})
})
