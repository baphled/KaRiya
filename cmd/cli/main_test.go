package main

import (
	"bytes"

	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
)

var _ = Describe("CLI Initialization", func() {
	Context("Version Flag", func() {
		It("should print version when --version flag is provided", func() {
			var buf bytes.Buffer
			exitCode := run([]string{"--version"}, &buf)

			Expect(exitCode).To(Equal(0))
			Expect(buf.String()).To(ContainSubstring("KaRiya CLI v"))
		})
	})

	Context("Help Flag", func() {
		It("should print help information when --help flag is provided", func() {
			var buf bytes.Buffer
			exitCode := run([]string{"--help"}, &buf)

			Expect(exitCode).To(Equal(0))
			Expect(buf.String()).To(ContainSubstring("KaRiya CLI - Career Journaling Tool"))
			Expect(buf.String()).To(ContainSubstring("Usage: kariya [options]"))
		})

		It("should show database flag in help", func() {
			var buf bytes.Buffer
			run([]string{"--help"}, &buf)

			Expect(buf.String()).To(ContainSubstring("--db"))
			Expect(buf.String()).To(ContainSubstring("--database"))
		})

		It("should show mode flag in help", func() {
			var buf bytes.Buffer
			run([]string{"--help"}, &buf)

			Expect(buf.String()).To(ContainSubstring("--mode"))
		})

		It("should show list flag in help", func() {
			var buf bytes.Buffer
			run([]string{"--help"}, &buf)

			Expect(buf.String()).To(ContainSubstring("--list"))
		})
	})

	Context("Mode Flag", func() {
		It("should reject invalid mode and show error", func() {
			var buf bytes.Buffer
			exitCode := run([]string{"--mode", "invalid"}, &buf)

			Expect(exitCode).To(Equal(1))
			Expect(buf.String()).To(ContainSubstring("Invalid mode"))
		})
	})

	Context("Database Flag", func() {
		It("should show database path options in help", func() {
			var buf bytes.Buffer
			run([]string{"--help"}, &buf)

			Expect(buf.String()).To(ContainSubstring("Use custom database path"))
			Expect(buf.String()).To(ContainSubstring("Default: ~/.kariya/events.db"))
		})
	})

	Context("List Flag", func() {
		It("should show --list flag in help", func() {
			var buf bytes.Buffer
			run([]string{"--help"}, &buf)

			Expect(buf.String()).To(ContainSubstring("--list"))
			Expect(buf.String()).To(ContainSubstring("Show recent events on startup"))
		})

		It("should show example with --list and --db together", func() {
			var buf bytes.Buffer
			run([]string{"--help"}, &buf)

			Expect(buf.String()).To(ContainSubstring("kariya --db ./events.db --list"))
		})
	})

	Context("Help Examples", func() {
		It("should show usage examples", func() {
			var buf bytes.Buffer
			run([]string{"--help"}, &buf)

			Expect(buf.String()).To(ContainSubstring("Examples:"))
			Expect(buf.String()).To(ContainSubstring("kariya"))
			Expect(buf.String()).To(ContainSubstring("--mode timeline"))
		})
	})
})
