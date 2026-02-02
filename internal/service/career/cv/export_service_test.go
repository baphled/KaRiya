package cv

import (
	"context"
	"errors"
	"io"
	"os"
	"path/filepath"
	"strings"

	"github.com/baphled/kariya/internal/config"
	"github.com/baphled/kariya/internal/domain/career"
	"github.com/baphled/kariya/internal/logger"
	"github.com/baphled/kariya/internal/testutil/fixtures"
	"github.com/onsi/ginkgo/v2"
	"github.com/onsi/gomega"
)

// MockClipboard is a test implementation of ClipboardWriter.
type MockClipboard struct {
	Content     string
	Unsupported bool
	WriteError  error
}

func (m *MockClipboard) WriteAll(text string) error {
	if m.WriteError != nil {
		return m.WriteError
	}
	m.Content = text
	return nil
}

func (m *MockClipboard) IsUnsupported() bool {
	return m.Unsupported
}

var _ = ginkgo.Describe("ExportService", func() {
	var (
		service *ExportService
		log     *logger.Logger
		ctx     context.Context
	)

	ginkgo.BeforeEach(func() {
		log = logger.New(io.Discard, logger.InfoLevel)
		service = NewExportService(log)
		ctx = context.Background()
	})

	ginkgo.Describe("ExportToText", func() {
		ginkgo.It("should export CV to plain text format", func() {
			cv := fixtures.CVViewWith("cv-1", "Senior Software Engineer CV", "Staff Engineer", "hiring_manager")

			bullet := fixtures.CVBulletWithSources("bullet-1", "section-1", "Led team of 5 engineers to deliver critical feature", []string{"event-1"}, []string{"fact-1"})
			bullet.Rank = 0.9
			bullet.Confidence = 0.95

			section := fixtures.CVSectionWith("section-1", "cv-1", "experience", "Experience", 1)
			section.Content = []*career.SectionContentGroup{
				fixtures.ContentGroupWithBullets("Acme Corp", []*career.CVBullet{bullet}),
			}

			sections := []*career.CVSection{section}
			bullets := map[string][]*career.CVBullet{
				"section-1": {bullet},
			}

			text, err := service.ExportToText(ctx, cv, sections, bullets)
			gomega.Expect(err).NotTo(gomega.HaveOccurred())
			gomega.Expect(text).NotTo(gomega.BeEmpty())
			gomega.Expect(text).To(gomega.ContainSubstring(strings.ToUpper("Senior Software Engineer CV")))
			gomega.Expect(text).To(gomega.ContainSubstring("Staff Engineer"))
			gomega.Expect(text).To(gomega.ContainSubstring("Led team of 5 engineers"))
			gomega.Expect(text).To(gomega.ContainSubstring("EXPERIENCE"))
		})

		ginkgo.It("should handle empty sections", func() {
			cv := fixtures.CVViewWith("cv-1", "Test CV", "Engineer", "hiring_manager")
			cv.SourceEventCount = 0
			cv.SourceFactCount = 0

			section := fixtures.CVSectionWith("section-1", "cv-1", "experience", "Experience", 1)
			section.Content = []*career.SectionContentGroup{}

			sections := []*career.CVSection{section}
			bullets := map[string][]*career.CVBullet{}

			text, err := service.ExportToText(ctx, cv, sections, bullets)
			gomega.Expect(err).NotTo(gomega.HaveOccurred())
			gomega.Expect(text).To(gomega.ContainSubstring("(No content)"))
		})

		ginkgo.It("should return error for nil CV", func() {
			text, err := service.ExportToText(ctx, nil, []*career.CVSection{}, map[string][]*career.CVBullet{})
			gomega.Expect(err).To(gomega.HaveOccurred())
			gomega.Expect(text).To(gomega.BeEmpty())
		})
	})

	ginkgo.Describe("ExportToMarkdown", func() {
		ginkgo.It("should export CV to markdown format", func() {
			cv := fixtures.CVViewWith("cv-1", "Senior Software Engineer CV", "Staff Engineer", "hiring_manager")

			bullet := fixtures.CVBulletWithSources("bullet-1", "section-1", "Led team of 5 engineers to deliver critical feature", []string{"event-1"}, []string{"fact-1"})
			bullet.Rank = 0.9
			bullet.Confidence = 0.95

			section := fixtures.CVSectionWith("section-1", "cv-1", "experience", "Experience", 1)
			section.Content = []*career.SectionContentGroup{
				fixtures.ContentGroupWithBullets("Acme Corp", []*career.CVBullet{bullet}),
			}

			sections := []*career.CVSection{section}
			bullets := map[string][]*career.CVBullet{
				"section-1": {bullet},
			}

			markdown, err := service.ExportToMarkdown(ctx, cv, sections, bullets)
			gomega.Expect(err).NotTo(gomega.HaveOccurred())
			gomega.Expect(markdown).NotTo(gomega.BeEmpty())
			gomega.Expect(markdown).To(gomega.ContainSubstring("# Senior Software Engineer CV"))
			gomega.Expect(markdown).To(gomega.ContainSubstring("## Experience"))
			gomega.Expect(markdown).To(gomega.ContainSubstring("- Led team of 5 engineers"))
			gomega.Expect(markdown).To(gomega.ContainSubstring("<!--"))
		})

		ginkgo.It("should include metadata as comments", func() {
			cv := fixtures.CVViewWith("cv-1", "Test CV", "Engineer", "recruiter")
			cv.SourceEventCount = 5
			cv.SourceFactCount = 2

			markdown, err := service.ExportToMarkdown(ctx, cv, []*career.CVSection{}, map[string][]*career.CVBullet{})
			gomega.Expect(err).NotTo(gomega.HaveOccurred())
			gomega.Expect(markdown).To(gomega.ContainSubstring("Target Role: Engineer"))
			gomega.Expect(markdown).To(gomega.ContainSubstring("Source Events: 5"))
		})

		ginkgo.It("should return error for nil CV", func() {
			markdown, err := service.ExportToMarkdown(ctx, nil, []*career.CVSection{}, map[string][]*career.CVBullet{})
			gomega.Expect(err).To(gomega.HaveOccurred())
			gomega.Expect(markdown).To(gomega.BeEmpty())
		})
	})

	ginkgo.Describe("ExportToYAML", func() {
		ginkgo.It("should export CV to YAML format", func() {
			cv := fixtures.CVViewWith("cv-1", "Senior Software Engineer CV", "Staff Engineer", "hiring_manager")

			bullet := fixtures.CVBulletWithSources("bullet-1", "section-1", "Led team of 5 engineers to deliver critical feature", []string{"event-1"}, []string{"fact-1"})
			bullet.Rank = 0.9
			bullet.Confidence = 0.95

			section := fixtures.CVSectionWith("section-1", "cv-1", "experience", "Experience", 1)
			section.Content = []*career.SectionContentGroup{
				fixtures.ContentGroupWithBullets("Acme Corp", []*career.CVBullet{bullet}),
			}

			sections := []*career.CVSection{section}
			bullets := map[string][]*career.CVBullet{
				"section-1": {bullet},
			}

			yaml, err := service.ExportToYAML(ctx, cv, sections, bullets)
			gomega.Expect(err).NotTo(gomega.HaveOccurred())
			gomega.Expect(yaml).NotTo(gomega.BeEmpty())
			gomega.Expect(yaml).To(gomega.ContainSubstring("name: Senior Software Engineer CV"))
			gomega.Expect(yaml).To(gomega.ContainSubstring("target_role: Staff Engineer"))
			// Check for the bullet text (YAML format uses "text:" field)
			gomega.Expect(yaml).To(gomega.ContainSubstring("text: Led team of 5 engineers to deliver critical feature"))
			// Check for content group structure
			gomega.Expect(yaml).To(gomega.ContainSubstring("header: Acme Corp"))
		})

		ginkgo.It("should return error for nil CV", func() {
			yaml, err := service.ExportToYAML(ctx, nil, []*career.CVSection{}, map[string][]*career.CVBullet{})
			gomega.Expect(err).To(gomega.HaveOccurred())
			gomega.Expect(yaml).To(gomega.BeEmpty())
		})
	})

	ginkgo.Describe("SaveToFile", func() {
		ginkgo.It("should save text export to file", func() {
			content := "Test CV Content\nWith multiple lines"
			filePath, err := service.SaveToFile(context.Background(), "Test CV", ExportFormatText, content)
			gomega.Expect(err).NotTo(gomega.HaveOccurred())
			gomega.Expect(filePath).NotTo(gomega.BeEmpty())

			// Verify file exists
			_, err = os.Stat(filePath)
			gomega.Expect(err).NotTo(gomega.HaveOccurred())

			// Clean up
			os.Remove(filePath)
		})

		ginkgo.It("should save markdown export to file", func() {
			content := "# Test CV\n## Section\n- Bullet point"
			filePath, err := service.SaveToFile(context.Background(), "Test CV", ExportFormatMarkdown, content)
			gomega.Expect(err).NotTo(gomega.HaveOccurred())
			gomega.Expect(filePath).To(gomega.ContainSubstring(".md"))

			// Clean up
			os.Remove(filePath)
		})

		ginkgo.It("should save YAML export to file", func() {
			content := "name: Test CV\ntarget_role: Engineer"
			filePath, err := service.SaveToFile(context.Background(), "Test CV", ExportFormatYAML, content)
			gomega.Expect(err).NotTo(gomega.HaveOccurred())
			gomega.Expect(filePath).To(gomega.ContainSubstring(".yaml"))

			// Clean up
			os.Remove(filePath)
		})

		ginkgo.It("should sanitize filename", func() {
			content := "Test content"
			filePath, err := service.SaveToFile(context.Background(), "Test/CV:Name*Invalid", ExportFormatText, content)
			gomega.Expect(err).NotTo(gomega.HaveOccurred())

			// Verify filename is sanitized
			filename := filepath.Base(filePath)
			gomega.Expect(filename).NotTo(gomega.ContainSubstring("/"))
			gomega.Expect(filename).NotTo(gomega.ContainSubstring(":"))
			gomega.Expect(filename).NotTo(gomega.ContainSubstring("*"))

			// Clean up
			os.Remove(filePath)
		})

		ginkgo.It("should create export directory if it doesn't exist", func() {
			content := "Test content"
			filePath, err := service.SaveToFile(context.Background(), "Test CV", ExportFormatText, content)
			gomega.Expect(err).NotTo(gomega.HaveOccurred())

			// Verify directory was created
			dir := filepath.Dir(filePath)
			_, err = os.Stat(dir)
			gomega.Expect(err).NotTo(gomega.HaveOccurred())

			// Clean up
			os.Remove(filePath)
		})

		ginkgo.It("should include timestamp in filename", func() {
			content := "Test content"
			filePath, err := service.SaveToFile(context.Background(), "Test CV", ExportFormatText, content)
			gomega.Expect(err).NotTo(gomega.HaveOccurred())

			filename := filepath.Base(filePath)
			// Check for timestamp pattern (YYYYMMDD_HHMMSS)
			gomega.Expect(filename).To(gomega.MatchRegexp(`\d{8}_\d{6}`))

			// Clean up
			os.Remove(filePath)
		})
	})

	ginkgo.Describe("GetExportPath", func() {
		ginkgo.It("should return valid export path", func() {
			path, err := service.GetExportPath()
			gomega.Expect(err).NotTo(gomega.HaveOccurred())
			gomega.Expect(path).NotTo(gomega.BeEmpty())
			gomega.Expect(path).To(gomega.ContainSubstring(".kariya"))
			gomega.Expect(path).To(gomega.ContainSubstring("cv_exports"))
		})
	})

	ginkgo.Describe("Helper Functions", func() {
		ginkgo.Describe("getFileExtension", func() {
			ginkgo.It("should return .txt for text format", func() {
				ext := getFileExtension(ExportFormatText)
				gomega.Expect(ext).To(gomega.Equal(".txt"))
			})

			ginkgo.It("should return .md for markdown format", func() {
				ext := getFileExtension(ExportFormatMarkdown)
				gomega.Expect(ext).To(gomega.Equal(".md"))
			})

			ginkgo.It("should return .yaml for YAML format", func() {
				ext := getFileExtension(ExportFormatYAML)
				gomega.Expect(ext).To(gomega.Equal(".yaml"))
			})

			ginkgo.It("should return .txt for unknown format", func() {
				ext := getFileExtension(ExportFormat("unknown"))
				gomega.Expect(ext).To(gomega.Equal(".txt"))
			})
		})

		ginkgo.Describe("sanitizeFilename", func() {
			ginkgo.It("should replace spaces with underscores", func() {
				result := sanitizeFilename("Test CV Name")
				gomega.Expect(result).To(gomega.Equal("Test_CV_Name"))
			})

			ginkgo.It("should remove invalid characters", func() {
				result := sanitizeFilename("Test/CV:Name*Invalid?")
				gomega.Expect(result).NotTo(gomega.ContainSubstring("/"))
				gomega.Expect(result).NotTo(gomega.ContainSubstring(":"))
				gomega.Expect(result).NotTo(gomega.ContainSubstring("*"))
				gomega.Expect(result).NotTo(gomega.ContainSubstring("?"))
			})

			ginkgo.It("should truncate long filenames", func() {
				longName := strings.Repeat("a", 200)
				result := sanitizeFilename(longName)
				gomega.Expect(len(result)).To(gomega.BeNumerically("<=", 100))
			})

			ginkgo.It("should handle empty string", func() {
				result := sanitizeFilename("")
				gomega.Expect(result).To(gomega.Equal(""))
			})
		})

		ginkgo.Describe("CopyToClipboard", func() {
			var mockClipboard *MockClipboard

			ginkgo.BeforeEach(func() {
				mockClipboard = &MockClipboard{}
				service = NewExportServiceWithClipboard(log, mockClipboard)
			})

			ginkgo.It("should copy content to clipboard when supported", func() {
				testContent := "Test CV Content for Clipboard"
				err := service.CopyToClipboard(ctx, testContent)
				gomega.Expect(err).NotTo(gomega.HaveOccurred())

				// Verify content was copied to mock
				gomega.Expect(mockClipboard.Content).To(gomega.Equal(testContent))
			})

			ginkgo.It("should return error for empty content", func() {
				err := service.CopyToClipboard(ctx, "")
				gomega.Expect(err).To(gomega.HaveOccurred())
				gomega.Expect(err.Error()).To(gomega.ContainSubstring("content is empty"))
			})

			ginkgo.It("should return ErrClipboardUnsupported when clipboard is not available", func() {
				mockClipboard.Unsupported = true
				err := service.CopyToClipboard(ctx, "test content")
				gomega.Expect(err).To(gomega.HaveOccurred())
				gomega.Expect(errors.Is(err, ErrClipboardUnsupported)).To(gomega.BeTrue())
				gomega.Expect(err.Error()).To(gomega.ContainSubstring("clipboard not available"))
			})

			ginkgo.It("should return error when clipboard write fails", func() {
				mockClipboard.WriteError = errors.New("clipboard write failed")
				err := service.CopyToClipboard(ctx, "test content")
				gomega.Expect(err).To(gomega.HaveOccurred())
				gomega.Expect(err.Error()).To(gomega.ContainSubstring("failed to copy to clipboard"))
			})
		})

	})

	ginkgo.Describe("Structure-Aware Export", func() {
		var (
			cv       *career.CVView
			sections []*career.CVSection
			bullets  map[string][]*career.CVBullet
		)

		ginkgo.BeforeEach(func() {
			cv = fixtures.CVViewWith("cv-1", "Senior Software Engineer CV", "senior_ic", "hiring_manager")

			bullet1 := fixtures.CVBulletWith("bullet-1", "section-exp", "Led migration to microservices architecture")
			bullet1.Rank = 0.9
			bullet1.Confidence = 0.85

			bullet2 := fixtures.CVBulletWith("bullet-2", "section-exp", "Implemented CI/CD pipeline")
			bullet2.SourceEventIDs = []string{"event-2"}
			bullet2.Rank = 0.7
			bullet2.InclusionReason = "execution"
			bullet2.Confidence = 0.70

			summarySection := fixtures.CVSectionWithSummary("section-summary", "cv-1", "Experienced software engineer with 10+ years in backend development.")

			expSection := fixtures.CVSectionWith("section-exp", "cv-1", "experience", "Experience", 1)
			expSection.Content = []*career.SectionContentGroup{
				fixtures.ContentGroupFull("TechCorp", "Jan 2020", "Present", []*career.CVBullet{bullet1, bullet2}),
			}

			sections = []*career.CVSection{summarySection, expSection}

			bullets = map[string][]*career.CVBullet{
				"section-exp": {bullet1, bullet2},
			}
		})

		ginkgo.Describe("Export with Standard structure", func() {
			ginkgo.It("should export to text format with standard structure", func() {
				content, err := service.Export(ctx, cv, sections, bullets, CVStructureStandard, ExportFormatText)
				gomega.Expect(err).NotTo(gomega.HaveOccurred())
				gomega.Expect(content).To(gomega.ContainSubstring("SENIOR SOFTWARE ENGINEER CV"))
				gomega.Expect(content).To(gomega.ContainSubstring("EXPERIENCE"))
				gomega.Expect(content).To(gomega.ContainSubstring("Led migration to microservices"))
				// Standard includes all bullets regardless of confidence
				gomega.Expect(content).To(gomega.ContainSubstring("Implemented CI/CD pipeline"))
			})

			ginkgo.It("should export to markdown format with standard structure", func() {
				content, err := service.Export(ctx, cv, sections, bullets, CVStructureStandard, ExportFormatMarkdown)
				gomega.Expect(err).NotTo(gomega.HaveOccurred())
				gomega.Expect(content).To(gomega.ContainSubstring("# Senior Software Engineer CV"))
				gomega.Expect(content).To(gomega.ContainSubstring("## Experience"))
				gomega.Expect(content).To(gomega.ContainSubstring("Led migration to microservices"))
				// Standard includes all bullets regardless of confidence
				gomega.Expect(content).To(gomega.ContainSubstring("Implemented CI/CD pipeline"))
			})
		})

		ginkgo.Describe("Export with Narrative structure", func() {
			ginkgo.It("should export to text format with narrative structure", func() {
				content, err := service.Export(ctx, cv, sections, bullets, CVStructureNarrative, ExportFormatText)
				gomega.Expect(err).NotTo(gomega.HaveOccurred())
				// Narrative export works even with empty profile (user hasn't onboarded)
				// Should NOT contain hardcoded personal data
				gomega.Expect(content).NotTo(gomega.ContainSubstring("YOMI COLLEDGE"))
				gomega.Expect(content).NotTo(gomega.ContainSubstring("boodah"))
				// Narrative sections should still be present
				gomega.Expect(content).To(gomega.ContainSubstring("CORE STRENGTHS"))
				gomega.Expect(content).To(gomega.ContainSubstring("LANGUAGES & TECHNOLOGIES"))
				gomega.Expect(content).To(gomega.ContainSubstring("SELECTED EXPERIENCE"))
				gomega.Expect(content).To(gomega.ContainSubstring("WHAT I BRING"))
				// High confidence bullet appears
				gomega.Expect(content).To(gomega.ContainSubstring("Led migration to microservices"))
				// Low confidence bullet filtered out
				gomega.Expect(content).NotTo(gomega.ContainSubstring("Implemented CI/CD pipeline"))
			})

			ginkgo.It("should export to markdown format with narrative structure", func() {
				content, err := service.Export(ctx, cv, sections, bullets, CVStructureNarrative, ExportFormatMarkdown)
				gomega.Expect(err).NotTo(gomega.HaveOccurred())
				// Narrative export works even with empty profile (user hasn't onboarded)
				// Should NOT contain hardcoded personal data
				gomega.Expect(content).NotTo(gomega.ContainSubstring("# Yomi Colledge"))
				gomega.Expect(content).NotTo(gomega.ContainSubstring("boodah"))
				// Narrative sections should still be present
				gomega.Expect(content).To(gomega.ContainSubstring("## Core Strengths"))
				gomega.Expect(content).To(gomega.ContainSubstring("## Languages & Technologies"))
				gomega.Expect(content).To(gomega.ContainSubstring("## Selected Experience"))
				gomega.Expect(content).To(gomega.ContainSubstring("## What I Bring"))
				// High confidence bullet appears
				gomega.Expect(content).To(gomega.ContainSubstring("Led migration to microservices"))
				// Low confidence bullet filtered out
				gomega.Expect(content).NotTo(gomega.ContainSubstring("Implemented CI/CD pipeline"))
			})
		})

		ginkgo.Describe("YAML export always uses standard structure", func() {
			ginkgo.It("should use standard structure for YAML even when narrative requested", func() {
				content, err := service.Export(ctx, cv, sections, bullets, CVStructureNarrative, ExportFormatYAML)
				gomega.Expect(err).NotTo(gomega.HaveOccurred())
				// YAML is data format, should contain raw data
				gomega.Expect(content).To(gomega.ContainSubstring("name: Senior Software Engineer CV"))
				gomega.Expect(content).To(gomega.ContainSubstring("target_role: senior_ic"))
				// Should include all bullets (no filtering)
				gomega.Expect(content).To(gomega.ContainSubstring("Led migration to microservices"))
				gomega.Expect(content).To(gomega.ContainSubstring("Implemented CI/CD pipeline"))
			})
		})

		ginkgo.Describe("Error handling", func() {
			ginkgo.It("should return error for nil CV", func() {
				_, err := service.Export(ctx, nil, sections, bullets, CVStructureStandard, ExportFormatText)
				gomega.Expect(err).To(gomega.HaveOccurred())
				gomega.Expect(err.Error()).To(gomega.ContainSubstring("CV view is nil"))
			})

			ginkgo.It("should return error for unknown format", func() {
				_, err := service.Export(ctx, cv, sections, bullets, CVStructureStandard, ExportFormat("unknown"))
				gomega.Expect(err).To(gomega.HaveOccurred())
				gomega.Expect(err.Error()).To(gomega.ContainSubstring("unknown export format"))
			})
		})
	})

	ginkgo.Describe("ExportWithProfile", func() {
		var (
			cv       *career.CVView
			sections []*career.CVSection
			bullets  map[string][]*career.CVBullet
		)

		ginkgo.BeforeEach(func() {
			cv = fixtures.CVViewWith("cv-1", "Test CV", "senior_ic", "hiring_manager")
			cv.SourceEventCount = 5
			cv.SourceFactCount = 3

			highConfBullet := fixtures.CVBulletWith("", "section-1", "High confidence achievement")

			expSection := fixtures.CVSectionWith("section-1", "cv-1", "experience", "Experience", 1)
			expSection.Content = []*career.SectionContentGroup{
				fixtures.ContentGroupFull("Company XYZ", "Jan 2020", "Present", []*career.CVBullet{highConfBullet}),
			}

			sections = []*career.CVSection{expSection}

			bullets = make(map[string][]*career.CVBullet)
		})

		ginkgo.It("should use default profile when profileCfg is nil", func() {
			content, err := service.ExportWithProfile(ctx, cv, sections, bullets, CVStructureNarrative, ExportFormatText, nil)
			gomega.Expect(err).NotTo(gomega.HaveOccurred())
			// Default profile has empty name - onboarding data not provided
			// The CV should still export (with empty/blank name section)
			gomega.Expect(content).NotTo(gomega.BeEmpty())
		})

		ginkgo.It("should use custom profile when provided", func() {
			profileCfg := &config.ProfileConfig{
				Name:      "Jane Doe",
				Email:     "jane@example.com",
				Title:     "Principal Engineer",
				Location:  "San Francisco, CA",
				GitHub:    "https://github.com/janedoe",
				Portfolio: "https://janedoe.dev",
				Languages: []string{"Python", "Rust", "TypeScript"},
				Frontend:  []string{"React", "Vue"},
				Systems:   []string{"Kubernetes", "AWS", "Terraform"},
				CoreStrengths: []string{
					"Distributed systems design",
					"Team leadership",
				},
				WhatIBring: []string{
					"Deep technical expertise",
					"Mentorship focus",
				},
			}

			content, err := service.ExportWithProfile(ctx, cv, sections, bullets, CVStructureNarrative, ExportFormatText, profileCfg)
			gomega.Expect(err).NotTo(gomega.HaveOccurred())
			// Custom profile name
			gomega.Expect(content).To(gomega.ContainSubstring("JANE DOE"))
			gomega.Expect(content).To(gomega.ContainSubstring("Principal Engineer"))
			gomega.Expect(content).To(gomega.ContainSubstring("jane@example.com"))
			gomega.Expect(content).To(gomega.ContainSubstring("Python, Rust, TypeScript"))
			gomega.Expect(content).To(gomega.ContainSubstring("Distributed systems design"))
			gomega.Expect(content).To(gomega.ContainSubstring("Deep technical expertise"))
		})

		ginkgo.It("should use custom profile in markdown format", func() {
			profileCfg := &config.ProfileConfig{
				Name:  "John Smith",
				Title: "Staff Engineer",
			}

			content, err := service.ExportWithProfile(ctx, cv, sections, bullets, CVStructureNarrative, ExportFormatMarkdown, profileCfg)
			gomega.Expect(err).NotTo(gomega.HaveOccurred())
			// Markdown uses title case for name
			gomega.Expect(content).To(gomega.ContainSubstring("# John Smith"))
			gomega.Expect(content).To(gomega.ContainSubstring("**Staff Engineer**"))
		})

		ginkgo.It("should keep empty fields when custom profile fields are empty", func() {
			profileCfg := &config.ProfileConfig{
				Name: "Custom Name",
				// All other fields empty - should remain empty (inference fills these)
			}

			content, err := service.ExportWithProfile(ctx, cv, sections, bullets, CVStructureNarrative, ExportFormatText, profileCfg)
			gomega.Expect(err).NotTo(gomega.HaveOccurred())
			// Custom name should be present
			gomega.Expect(content).To(gomega.ContainSubstring("CUSTOM NAME"))
			// Empty fields should NOT contain hardcoded personal data
			gomega.Expect(content).NotTo(gomega.ContainSubstring("Yomi"))
			gomega.Expect(content).NotTo(gomega.ContainSubstring("boodah"))
		})

		ginkgo.It("should use standard structure regardless of profile for YAML format", func() {
			profileCfg := &config.ProfileConfig{
				Name: "Should Not Appear",
			}

			content, err := service.ExportWithProfile(ctx, cv, sections, bullets, CVStructureNarrative, ExportFormatYAML, profileCfg)
			gomega.Expect(err).NotTo(gomega.HaveOccurred())
			// YAML uses standard format, not narrative profile
			gomega.Expect(content).NotTo(gomega.ContainSubstring("Should Not Appear"))
			gomega.Expect(content).To(gomega.ContainSubstring("name: Test CV"))
		})
	})

	ginkgo.Describe("NarrativeProfileFromConfig", func() {
		ginkgo.It("should return empty default profile when config is nil", func() {
			profile := NarrativeProfileFromConfig(nil)
			// Default profile should be empty - no hardcoded personal data
			gomega.Expect(profile.Name).To(gomega.BeEmpty())
			gomega.Expect(profile.Role).To(gomega.BeEmpty())
			gomega.Expect(profile.CoreStrengths).To(gomega.BeEmpty())
			gomega.Expect(profile.ValuePropositions).To(gomega.BeEmpty())
		})

		ginkgo.It("should use config values when provided", func() {
			cfg := &config.ProfileConfig{
				Name:      "Test User",
				Title:     "Lead Developer",
				Location:  "London, UK",
				Email:     "test@example.com",
				GitHub:    "https://github.com/testuser",
				Portfolio: "https://test.dev",
				Languages: []string{"Go", "Python"},
				Frontend:  []string{"Angular"},
				Systems:   []string{"Docker", "GCP"},
				CoreStrengths: []string{
					"Backend development",
					"System architecture",
				},
				WhatIBring: []string{
					"Pragmatic approach",
					"Strong communication",
				},
			}

			profile := NarrativeProfileFromConfig(cfg)
			gomega.Expect(profile.Name).To(gomega.Equal("Test User"))
			gomega.Expect(profile.Role).To(gomega.Equal("Lead Developer"))
			gomega.Expect(profile.Location).To(gomega.Equal("London, UK"))
			gomega.Expect(profile.Email).To(gomega.Equal("test@example.com"))
			gomega.Expect(profile.GitHub).To(gomega.Equal("https://github.com/testuser"))
			gomega.Expect(profile.Portfolio).To(gomega.Equal("https://test.dev"))
			gomega.Expect(profile.Languages).To(gomega.Equal([]string{"Go", "Python"}))
			gomega.Expect(profile.Frontend).To(gomega.Equal([]string{"Angular"}))
			gomega.Expect(profile.Systems).To(gomega.Equal([]string{"Docker", "GCP"}))
			gomega.Expect(profile.CoreStrengths).To(gomega.Equal([]string{"Backend development", "System architecture"}))
			gomega.Expect(profile.ValuePropositions).To(gomega.Equal([]string{"Pragmatic approach", "Strong communication"}))
		})

		ginkgo.It("should keep empty fields empty (no hardcoded defaults)", func() {
			cfg := &config.ProfileConfig{
				Name: "Only Name Set",
				// All other fields empty - should remain empty for inference service to fill
			}

			profile := NarrativeProfileFromConfig(cfg)

			gomega.Expect(profile.Name).To(gomega.Equal("Only Name Set"))
			// All other fields should be empty (not hardcoded defaults)
			gomega.Expect(profile.Role).To(gomega.BeEmpty())
			gomega.Expect(profile.Location).To(gomega.BeEmpty())
			gomega.Expect(profile.Email).To(gomega.BeEmpty())
			gomega.Expect(profile.GitHub).To(gomega.BeEmpty())
			gomega.Expect(profile.Portfolio).To(gomega.BeEmpty())
			gomega.Expect(profile.Languages).To(gomega.BeEmpty())
			gomega.Expect(profile.Frontend).To(gomega.BeEmpty())
			gomega.Expect(profile.Systems).To(gomega.BeEmpty())
			gomega.Expect(profile.CoreStrengths).To(gomega.BeEmpty())
			gomega.Expect(profile.ValuePropositions).To(gomega.BeEmpty())
		})
	})
})
