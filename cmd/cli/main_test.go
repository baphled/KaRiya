package main

import (
	"bytes"
	"fmt"

	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
)

var _ = Describe("CLI Initialization", func() {
	Context("Version Flag", func() {
		It("should print version when --version flag is provided", func() {
			var buf, errBuf bytes.Buffer
			exitCode := run([]string{"--version"}, &buf, &errBuf)

			Expect(exitCode).To(Equal(0))
			Expect(buf.String()).To(ContainSubstring("dev"))
		})
	})

	Context("Help Flag", func() {
		It("should print help information when --help flag is provided", func() {
			var buf, errBuf bytes.Buffer
			exitCode := run([]string{"--help"}, &buf, &errBuf)

			Expect(exitCode).To(Equal(0))
			Expect(buf.String()).To(ContainSubstring("capturing career achievements"))
			Expect(buf.String()).To(ContainSubstring("Usage:"))
		})

		It("should show database flag in help", func() {
			var buf, errBuf bytes.Buffer
			run([]string{"--help"}, &buf, &errBuf)

			Expect(buf.String()).To(ContainSubstring("--db"))
		})

		It("should show mode flag in help", func() {
			var buf, errBuf bytes.Buffer
			run([]string{"--help"}, &buf, &errBuf)

			Expect(buf.String()).To(ContainSubstring("--mode"))
		})

		It("should show import command in help", func() {
			var buf, errBuf bytes.Buffer
			run([]string{"--help"}, &buf, &errBuf)

			Expect(buf.String()).To(ContainSubstring("import"))
			Expect(buf.String()).To(ContainSubstring("Import career events from CSV file"))
		})

		It("should show bursts command in help", func() {
			var buf, errBuf bytes.Buffer
			run([]string{"--help"}, &buf, &errBuf)

			Expect(buf.String()).To(ContainSubstring("bursts"))
		})

		It("should show facts command in help", func() {
			var buf, errBuf bytes.Buffer
			run([]string{"--help"}, &buf, &errBuf)

			Expect(buf.String()).To(ContainSubstring("facts"))
		})

		It("should show skills command in help", func() {
			var buf, errBuf bytes.Buffer
			run([]string{"--help"}, &buf, &errBuf)

			Expect(buf.String()).To(ContainSubstring("skills"))
		})
	})

	Context("Database Flag", func() {
		It("should show database path options in help", func() {
			var buf, errBuf bytes.Buffer
			run([]string{"--help"}, &buf, &errBuf)

			Expect(buf.String()).To(ContainSubstring("--db"))
		})
	})

	Context("Subcommands", func() {
		It("should show bursts subcommand help", func() {
			var buf, errBuf bytes.Buffer
			run([]string{"bursts", "--help"}, &buf, &errBuf)

			Expect(buf.String()).To(ContainSubstring("detect"))
			Expect(buf.String()).To(ContainSubstring("list"))
		})

		It("should show facts subcommand help", func() {
			var buf, errBuf bytes.Buffer
			run([]string{"facts", "--help"}, &buf, &errBuf)

			Expect(buf.String()).To(ContainSubstring("extract"))
			Expect(buf.String()).To(ContainSubstring("list"))
		})

		It("should show skills subcommand help", func() {
			var buf, errBuf bytes.Buffer
			run([]string{"skills", "--help"}, &buf, &errBuf)

			Expect(buf.String()).To(ContainSubstring("recategorize"))
		})

		It("should show import subcommand help", func() {
			var buf, errBuf bytes.Buffer
			run([]string{"import", "--help"}, &buf, &errBuf)

			Expect(buf.String()).To(Or(
				ContainSubstring("Import career events from CSV file"),
				ContainSubstring("Import career events from a CSV file"),
			))
		})
	})

	Context("Burst Detection", func() {
		It("should handle bursts detect with empty database", func() {
			var buf, errBuf bytes.Buffer
			exitCode := run([]string{"--in-memory", "bursts", "detect"}, &buf, &errBuf)

			fmt.Fprintf(GinkgoWriter, "OUTPUT: %s\n", buf.String())
			fmt.Fprintf(GinkgoWriter, "ERROR: %s\n", errBuf.String())
			fmt.Fprintf(GinkgoWriter, "EXIT: %d\n", exitCode)

			Expect(exitCode).To(Equal(0))
		})

		It("should handle bursts list with empty database", func() {
			var buf, errBuf bytes.Buffer
			exitCode := run([]string{"--in-memory", "bursts", "list"}, &buf, &errBuf)

			Expect(exitCode).To(Equal(0))
			Expect(buf.String()).To(ContainSubstring("No bursts found"))
		})
	})

	Context("Fact Extraction", func() {
		It("should handle facts extract with empty database", func() {
			var buf, errBuf bytes.Buffer
			exitCode := run([]string{"--in-memory", "facts", "extract"}, &buf, &errBuf)

			Expect(exitCode).To(Equal(0))
		})

		It("should handle facts list with empty database", func() {
			var buf, errBuf bytes.Buffer
			exitCode := run([]string{"--in-memory", "facts", "list"}, &buf, &errBuf)

			Expect(exitCode).To(Equal(0))
			Expect(buf.String()).To(ContainSubstring("No facts found"))
		})
	})

	Context("Skills Recategorization", func() {
		It("should handle skills recategorize with empty database", func() {
			var buf, errBuf bytes.Buffer
			exitCode := run([]string{"--in-memory", "skills", "recategorize"}, &buf, &errBuf)

			Expect(exitCode).To(Equal(0))
		})
	})
})
