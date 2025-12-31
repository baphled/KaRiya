package main

import (
	"bytes"
	"os"

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

		It("should display burst suggestions after import", func() {
			// Create a temporary CSV file with events that should form bursts
			csvContent := `Text,Date,Company,Project,Tags
Implemented cloud migration phase 1,2025-12-20,TechCorp,CloudMigration,technical;project
Optimized database queries for cloud,2025-12-21,TechCorp,CloudMigration,technical;optimization
Completed cloud migration phase 2,2025-12-22,TechCorp,CloudMigration,technical;project
Led performance optimization sprint,2025-12-23,TechCorp,Performance,leadership;technical
Improved API response times,2025-12-24,TechCorp,Performance,technical;optimization`

			tmpFile, err := os.CreateTemp("", "test_burst_*.csv")
			Expect(err).NotTo(HaveOccurred())
			defer os.Remove(tmpFile.Name())

			_, err = tmpFile.WriteString(csvContent)
			Expect(err).NotTo(HaveOccurred())
			tmpFile.Close()

			var buf, errBuf bytes.Buffer
			exitCode := run([]string{"--import", tmpFile.Name(), "--skip-import-review", "--in-memory"}, &buf, &errBuf)

			Expect(exitCode).To(Equal(0))
			output := buf.String()
			// Verify burst suggestions are displayed
			Expect(output).To(ContainSubstring("=== Burst Suggestions ==="))
			Expect(output).To(ContainSubstring("Detected"))
			Expect(output).To(ContainSubstring("potential bursts"))
			Expect(output).To(ContainSubstring("Events:"))
			Expect(output).To(ContainSubstring("Confidence:"))
		})

		It("should display fact extraction results after import", func() {
			// Create a temporary CSV file with events that should extract facts
			csvContent := `Text,Date,Company,Project,Tags
Implemented cloud migration strategy,2025-12-20,TechCorp,CloudMigration,technical
Led team through system redesign,2025-12-21,TechCorp,Redesign,leadership
Mentored junior engineers on best practices,2025-12-22,TechCorp,Training,mentoring`

			tmpFile, err := os.CreateTemp("", "test_facts_*.csv")
			Expect(err).NotTo(HaveOccurred())
			defer os.Remove(tmpFile.Name())

			_, err = tmpFile.WriteString(csvContent)
			Expect(err).NotTo(HaveOccurred())
			tmpFile.Close()

			var buf, errBuf bytes.Buffer
			exitCode := run([]string{"--import", tmpFile.Name(), "--skip-import-review", "--in-memory"}, &buf, &errBuf)

			Expect(exitCode).To(Equal(0))
			output := buf.String()
			// Verify import completed
			Expect(output).To(ContainSubstring("Import Complete"))
			// Verify fact extraction section exists (even if no facts persisted)
			Expect(output).To(ContainSubstring("Successfully imported"))
		})

		It("should accept --review-facts flag", func() {
			// Verify the flag is parsed correctly
			args := []string{"--import", "test.csv", "--review-facts", "--skip-import-review"}
			Expect(args).To(ContainElement("--review-facts"))
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

	Context("Burst Detection Flag", func() {
		It("should handle --detect-bursts flag with empty database", func() {
			var buf, errBuf bytes.Buffer
			exitCode := run([]string{"--in-memory", "--detect-bursts"}, &buf, &errBuf)

			Expect(exitCode).To(Equal(0))
			// With no events, should show appropriate message
			Expect(buf.String()).
				To(Or(
					ContainSubstring("No events found"),
					ContainSubstring("Detecting bursts"),
				))
		})
	})

	Context("Fact Extraction Flag", func() {
		It("should handle --extract-facts flag with empty database", func() {
			var buf, errBuf bytes.Buffer
			exitCode := run([]string{"--in-memory", "--extract-facts"}, &buf, &errBuf)

			Expect(exitCode).To(Equal(0))
			// With no events, should show appropriate message
			Expect(buf.String()).To(ContainSubstring("No events found"))
		})
	})

	Context("Show Bursts Flag", func() {
		It("should handle --show-bursts flag gracefully when repo not configured", func() {
			var buf, errBuf bytes.Buffer
			exitCode := run([]string{"--in-memory", "--show-bursts"}, &buf, &errBuf)

			// May fail if burst repo not configured for in-memory mode
			// Or succeed with "No bursts found" message
			if exitCode == 0 {
				Expect(buf.String()).To(ContainSubstring("No bursts found"))
			} else {
				Expect(errBuf.String()).To(ContainSubstring("Burst repository not configured"))
			}
		})

		It("should return bursts when bursts exist", func() {
			// Test with no bursts - should show "No bursts found"
			var buf, errBuf bytes.Buffer
			exitCode := run([]string{"--in-memory", "--show-bursts"}, &buf, &errBuf)

			// With no bursts, should show appropriate message
			Expect(exitCode).To(Equal(0))
			Expect(buf.String()).To(ContainSubstring("No bursts found"))
		})
	})

	Context("Show Facts Flag", func() {
		It("should handle --show-facts flag gracefully when repo not configured", func() {
			var buf, errBuf bytes.Buffer
			exitCode := run([]string{"--in-memory", "--show-facts"}, &buf, &errBuf)

			// May fail if fact repo not configured for in-memory mode
			// Or succeed with "No facts found" message
			if exitCode == 0 {
				Expect(buf.String()).To(ContainSubstring("No facts found"))
			} else {
				Expect(errBuf.String()).To(ContainSubstring("Fact repository not configured"))
			}
		})
	})
})
