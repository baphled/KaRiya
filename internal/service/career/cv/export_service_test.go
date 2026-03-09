package cv

import (
	"context"
	"errors"
	"fmt"
	"io"
	"os"
	"path/filepath"
	"strings"

	"github.com/baphled/kariya/internal/config"
	"github.com/baphled/kariya/internal/domain/career"
	"github.com/baphled/kariya/internal/logger"
	"github.com/baphled/kariya/internal/testutil/fixtures"
	mockrepo "github.com/baphled/kariya/internal/testutil/mocks/repository"
	"github.com/golang/mock/gomock"
	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
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

var _ = Describe("ExportService", func() {
	var (
		service *ExportService
		log     *logger.Logger
		ctx     context.Context
	)

	BeforeEach(func() {
		log = logger.New(io.Discard, logger.InfoLevel)
		service = NewExportService(log, nil, nil)
		ctx = context.Background()
	})

	Describe("ExportToText", func() {
		It("should export CV to plain text format", func() {
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
			Expect(err).NotTo(HaveOccurred())
			Expect(text).NotTo(BeEmpty())
			Expect(text).To(ContainSubstring(strings.ToUpper("Senior Software Engineer CV")))
			Expect(text).To(ContainSubstring("Staff Engineer"))
			Expect(text).To(ContainSubstring("Led team of 5 engineers"))
			Expect(text).To(ContainSubstring("EXPERIENCE"))
		})

		It("should handle empty sections", func() {
			cv := fixtures.CVViewWith("cv-1", "Test CV", "Engineer", "hiring_manager")
			cv.SourceEventCount = 0
			cv.SourceFactCount = 0

			section := fixtures.CVSectionWith("section-1", "cv-1", "experience", "Experience", 1)
			section.Content = []*career.SectionContentGroup{}

			sections := []*career.CVSection{section}
			bullets := map[string][]*career.CVBullet{}

			text, err := service.ExportToText(ctx, cv, sections, bullets)
			Expect(err).NotTo(HaveOccurred())
			Expect(text).To(ContainSubstring("(No content)"))
		})

		It("should return error for nil CV", func() {
			text, err := service.ExportToText(ctx, nil, []*career.CVSection{}, map[string][]*career.CVBullet{})
			Expect(err).To(HaveOccurred())
			Expect(text).To(BeEmpty())
		})
	})

	Describe("ExportToMarkdown", func() {
		It("should export CV to markdown format", func() {
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
			Expect(err).NotTo(HaveOccurred())
			Expect(markdown).NotTo(BeEmpty())
			Expect(markdown).To(ContainSubstring("# Senior Software Engineer CV"))
			Expect(markdown).To(ContainSubstring("## Experience"))
			Expect(markdown).To(ContainSubstring("- Led team of 5 engineers"))
			Expect(markdown).To(ContainSubstring("<!--"))
		})

		It("should include metadata as comments", func() {
			cv := fixtures.CVViewWith("cv-1", "Test CV", "Engineer", "recruiter")
			cv.SourceEventCount = 5
			cv.SourceFactCount = 2

			markdown, err := service.ExportToMarkdown(ctx, cv, []*career.CVSection{}, map[string][]*career.CVBullet{})
			Expect(err).NotTo(HaveOccurred())
			Expect(markdown).To(ContainSubstring("Target Role: Engineer"))
			Expect(markdown).To(ContainSubstring("Source Events: 5"))
		})

		It("should return error for nil CV", func() {
			markdown, err := service.ExportToMarkdown(ctx, nil, []*career.CVSection{}, map[string][]*career.CVBullet{})
			Expect(err).To(HaveOccurred())
			Expect(markdown).To(BeEmpty())
		})
	})

	Describe("ExportToYAML", func() {
		It("should export CV to flat YAML format", func() {
			cv := fixtures.CVViewWith("cv-1", "Senior Software Engineer CV", "Staff Engineer", "hiring_manager")

			bullet := fixtures.CVBulletWithSources("bullet-1", "section-1", "Led team of 5 engineers to deliver critical feature", []string{"event-1"}, []string{"fact-1"})
			bullet.Rank = 0.9
			bullet.Confidence = 0.95

			section := fixtures.CVSectionWith("section-1", "cv-1", "experience", "Experience", 1)
			section.Content = []*career.SectionContentGroup{
				fixtures.ContentGroupWithBullets("Acme Corp", []*career.CVBullet{bullet}),
			}

			sections := []*career.CVSection{section}

			yaml, err := service.ExportToYAML(ctx, cv, sections, nil)
			Expect(err).NotTo(HaveOccurred())
			Expect(yaml).NotTo(BeEmpty())
			// YAML now uses flat format with first_name, last_name, etc.
			Expect(yaml).To(ContainSubstring("first_name:"))
			Expect(yaml).To(ContainSubstring("last_name:"))
			Expect(yaml).To(ContainSubstring("jobs:"))
			Expect(yaml).To(ContainSubstring("company: Acme Corp"))
			Expect(yaml).To(ContainSubstring("- Led team of 5 engineers to deliver critical feature"))
		})

		It("should return error for nil CV", func() {
			yaml, err := service.ExportToYAML(ctx, nil, []*career.CVSection{}, nil)
			Expect(err).To(HaveOccurred())
			Expect(yaml).To(BeEmpty())
		})
	})

	Describe("SaveToFile", func() {
		It("should save text export to file", func() {
			content := "Test CV Content\nWith multiple lines"
			filePath, err := service.SaveToFile(context.Background(), "Test CV", ExportFormatText, content)
			Expect(err).NotTo(HaveOccurred())
			Expect(filePath).NotTo(BeEmpty())

			// Verify file exists
			_, err = os.Stat(filePath)
			Expect(err).NotTo(HaveOccurred())

			// Clean up
			os.Remove(filePath)
		})

		It("should save markdown export to file", func() {
			content := "# Test CV\n## Section\n- Bullet point"
			filePath, err := service.SaveToFile(context.Background(), "Test CV", ExportFormatMarkdown, content)
			Expect(err).NotTo(HaveOccurred())
			Expect(filePath).To(ContainSubstring(".md"))

			// Clean up
			os.Remove(filePath)
		})

		It("should save YAML export to file", func() {
			content := "name: Test CV\ntarget_role: Engineer"
			filePath, err := service.SaveToFile(context.Background(), "Test CV", ExportFormatYAML, content)
			Expect(err).NotTo(HaveOccurred())
			Expect(filePath).To(ContainSubstring(".yaml"))

			// Clean up
			os.Remove(filePath)
		})

		It("should sanitize filename", func() {
			content := "Test content"
			filePath, err := service.SaveToFile(context.Background(), "Test/CV:Name*Invalid", ExportFormatText, content)
			Expect(err).NotTo(HaveOccurred())

			// Verify filename is sanitized
			filename := filepath.Base(filePath)
			Expect(filename).NotTo(ContainSubstring("/"))
			Expect(filename).NotTo(ContainSubstring(":"))
			Expect(filename).NotTo(ContainSubstring("*"))

			// Clean up
			os.Remove(filePath)
		})

		It("should create export directory if it doesn't exist", func() {
			content := "Test content"
			filePath, err := service.SaveToFile(context.Background(), "Test CV", ExportFormatText, content)
			Expect(err).NotTo(HaveOccurred())

			// Verify directory was created
			dir := filepath.Dir(filePath)
			_, err = os.Stat(dir)
			Expect(err).NotTo(HaveOccurred())

			// Clean up
			os.Remove(filePath)
		})

		It("should include timestamp in filename", func() {
			content := "Test content"
			filePath, err := service.SaveToFile(context.Background(), "Test CV", ExportFormatText, content)
			Expect(err).NotTo(HaveOccurred())

			filename := filepath.Base(filePath)
			// Check for timestamp pattern (YYYYMMDD_HHMMSS)
			Expect(filename).To(MatchRegexp(`\d{8}_\d{6}`))

			// Clean up
			os.Remove(filePath)
		})
	})

	Describe("GetExportPath", func() {
		It("should return valid export path", func() {
			path, err := service.GetExportPath()
			Expect(err).NotTo(HaveOccurred())
			Expect(path).NotTo(BeEmpty())
			Expect(path).To(ContainSubstring(".kariya"))
			Expect(path).To(ContainSubstring("cv_exports"))
		})
	})

	Describe("Helper Functions", func() {
		Describe("getFileExtension", func() {
			It("should return .txt for text format", func() {
				ext := getFileExtension(ExportFormatText)
				Expect(ext).To(Equal(".txt"))
			})

			It("should return .md for markdown format", func() {
				ext := getFileExtension(ExportFormatMarkdown)
				Expect(ext).To(Equal(".md"))
			})

			It("should return .yaml for YAML format", func() {
				ext := getFileExtension(ExportFormatYAML)
				Expect(ext).To(Equal(".yaml"))
			})

			It("should return .txt for unknown format", func() {
				ext := getFileExtension(ExportFormat("unknown"))
				Expect(ext).To(Equal(".txt"))
			})
		})

		Describe("sanitizeFilename", func() {
			It("should replace spaces with underscores", func() {
				result := sanitizeFilename("Test CV Name")
				Expect(result).To(Equal("Test_CV_Name"))
			})

			It("should remove invalid characters", func() {
				result := sanitizeFilename("Test/CV:Name*Invalid?")
				Expect(result).NotTo(ContainSubstring("/"))
				Expect(result).NotTo(ContainSubstring(":"))
				Expect(result).NotTo(ContainSubstring("*"))
				Expect(result).NotTo(ContainSubstring("?"))
			})

			It("should truncate long filenames", func() {
				longName := strings.Repeat("a", 200)
				result := sanitizeFilename(longName)
				Expect(len(result)).To(BeNumerically("<=", 100))
			})

			It("should handle empty string", func() {
				result := sanitizeFilename("")
				Expect(result).To(Equal(""))
			})
		})

		Describe("ensureURL", func() {
			It("returns empty string for empty input", func() {
				result := ensureURL("", "https://example.com/")
				Expect(result).To(Equal(""))
			})

			It("returns input unchanged when already http URL", func() {
				result := ensureURL("http://example.com/user", "https://github.com/")
				Expect(result).To(Equal("http://example.com/user"))
			})

			It("returns input unchanged when already https URL", func() {
				result := ensureURL("https://github.com/baphled", "https://github.com/")
				Expect(result).To(Equal("https://github.com/baphled"))
			})

			It("prepends prefix to bare username", func() {
				result := ensureURL("baphled", "https://github.com/")
				Expect(result).To(Equal("https://github.com/baphled"))
			})

			It("prepends LinkedIn prefix to bare username", func() {
				result := ensureURL("yomicolledge", "https://www.linkedin.com/in/")
				Expect(result).To(Equal("https://www.linkedin.com/in/yomicolledge"))
			})

			It("handles username with hyphens", func() {
				result := ensureURL("john-doe", "https://github.com/")
				Expect(result).To(Equal("https://github.com/john-doe"))
			})
		})

		Describe("CopyToClipboard", func() {
			var mockClipboard *MockClipboard

			BeforeEach(func() {
				mockClipboard = &MockClipboard{}
				service = NewExportServiceWithClipboard(log, mockClipboard, nil, nil)
			})

			It("should copy content to clipboard when supported", func() {
				testContent := "Test CV Content for Clipboard"
				err := service.CopyToClipboard(ctx, testContent)
				Expect(err).NotTo(HaveOccurred())

				// Verify content was copied to mock
				Expect(mockClipboard.Content).To(Equal(testContent))
			})

			It("should return error for empty content", func() {
				err := service.CopyToClipboard(ctx, "")
				Expect(err).To(HaveOccurred())
				Expect(err.Error()).To(ContainSubstring("content is empty"))
			})

			It("should return ErrClipboardUnsupported when clipboard is not available", func() {
				mockClipboard.Unsupported = true
				err := service.CopyToClipboard(ctx, "test content")
				Expect(err).To(HaveOccurred())
				Expect(errors.Is(err, ErrClipboardUnsupported)).To(BeTrue())
				Expect(err.Error()).To(ContainSubstring("clipboard not available"))
			})

			It("should return error when clipboard write fails", func() {
				mockClipboard.WriteError = errors.New("clipboard write failed")
				err := service.CopyToClipboard(ctx, "test content")
				Expect(err).To(HaveOccurred())
				Expect(err.Error()).To(ContainSubstring("failed to copy to clipboard"))
			})
		})

	})

	Describe("Structure-Aware Export", func() {
		var (
			cv       *career.CVView
			sections []*career.CVSection
			bullets  map[string][]*career.CVBullet
		)

		BeforeEach(func() {
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

		Describe("Export with Standard structure", func() {
			It("should export to text format with standard structure", func() {
				content, err := service.Export(ctx, cv, sections, bullets, CVStructureStandard, ExportFormatText)
				Expect(err).NotTo(HaveOccurred())
				Expect(content).To(ContainSubstring("SENIOR SOFTWARE ENGINEER CV"))
				Expect(content).To(ContainSubstring("EXPERIENCE"))
				Expect(content).To(ContainSubstring("Led migration to microservices"))
				// Standard includes all bullets regardless of confidence
				Expect(content).To(ContainSubstring("Implemented CI/CD pipeline"))
			})

			It("should export to markdown format with standard structure", func() {
				content, err := service.Export(ctx, cv, sections, bullets, CVStructureStandard, ExportFormatMarkdown)
				Expect(err).NotTo(HaveOccurred())
				Expect(content).To(ContainSubstring("# Senior Software Engineer CV"))
				Expect(content).To(ContainSubstring("## Experience"))
				Expect(content).To(ContainSubstring("Led migration to microservices"))
				// Standard includes all bullets regardless of confidence
				Expect(content).To(ContainSubstring("Implemented CI/CD pipeline"))
			})
		})

		Describe("Export with Narrative structure", func() {
			It("should export to text format with narrative structure", func() {
				content, err := service.Export(ctx, cv, sections, bullets, CVStructureNarrative, ExportFormatText)
				Expect(err).NotTo(HaveOccurred())
				// Narrative export works even with empty profile (user hasn't onboarded)
				// Should NOT contain hardcoded personal data
				Expect(content).NotTo(ContainSubstring("YOMI COLLEDGE"))
				Expect(content).NotTo(ContainSubstring("boodah"))
				// Narrative sections should still be present
				Expect(content).To(ContainSubstring("CORE STRENGTHS"))
				Expect(content).To(ContainSubstring("LANGUAGES & TECHNOLOGIES"))
				Expect(content).To(ContainSubstring("SELECTED EXPERIENCE"))
				Expect(content).To(ContainSubstring("WHAT I BRING"))
				// High confidence bullet appears
				Expect(content).To(ContainSubstring("Led migration to microservices"))
				// Low confidence bullet filtered out
				Expect(content).NotTo(ContainSubstring("Implemented CI/CD pipeline"))
			})

			It("should export to markdown format with narrative structure", func() {
				content, err := service.Export(ctx, cv, sections, bullets, CVStructureNarrative, ExportFormatMarkdown)
				Expect(err).NotTo(HaveOccurred())
				// Narrative export works even with empty profile (user hasn't onboarded)
				// Should NOT contain hardcoded personal data
				Expect(content).NotTo(ContainSubstring("# Yomi Colledge"))
				Expect(content).NotTo(ContainSubstring("boodah"))
				// Narrative sections should still be present
				Expect(content).To(ContainSubstring("## Core Strengths"))
				Expect(content).To(ContainSubstring("## Languages & Technologies"))
				Expect(content).To(ContainSubstring("## Selected Experience"))
				Expect(content).To(ContainSubstring("## What I Bring"))
				// High confidence bullet appears
				Expect(content).To(ContainSubstring("Led migration to microservices"))
				// Low confidence bullet filtered out
				Expect(content).NotTo(ContainSubstring("Implemented CI/CD pipeline"))
			})
		})

		Describe("YAML export uses flat format", func() {
			It("should use flat format for YAML even when other structure requested", func() {
				content, err := service.Export(ctx, cv, sections, bullets, CVStructureNarrative, ExportFormatYAML)
				Expect(err).NotTo(HaveOccurred())
				// YAML uses flat structure for external tools
				Expect(content).To(ContainSubstring("first_name:"))
				Expect(content).To(ContainSubstring("jobs:"))
				Expect(content).To(ContainSubstring("skills:"))
				// Should NOT contain internal CVView format
				Expect(content).NotTo(ContainSubstring("name: Senior Software Engineer CV"))
				Expect(content).NotTo(ContainSubstring("target_role:"))
			})
		})

		Describe("Error handling", func() {
			It("should return error for nil CV", func() {
				_, err := service.Export(ctx, nil, sections, bullets, CVStructureStandard, ExportFormatText)
				Expect(err).To(HaveOccurred())
				Expect(err.Error()).To(ContainSubstring("CV view is nil"))
			})

			It("should return error for unknown format", func() {
				_, err := service.Export(ctx, cv, sections, bullets, CVStructureStandard, ExportFormat("unknown"))
				Expect(err).To(HaveOccurred())
				Expect(err.Error()).To(ContainSubstring("unknown export format"))
			})
		})
	})

	Describe("ExportWithProfile", func() {
		var (
			cv       *career.CVView
			sections []*career.CVSection
			bullets  map[string][]*career.CVBullet
		)

		BeforeEach(func() {
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

		It("should use default profile when profileCfg is nil", func() {
			content, err := service.ExportWithProfile(ctx, cv, sections, bullets, CVStructureNarrative, ExportFormatText, nil)
			Expect(err).NotTo(HaveOccurred())
			// Default profile has empty name - onboarding data not provided
			// The CV should still export (with empty/blank name section)
			Expect(content).NotTo(BeEmpty())
		})

		It("should use custom profile when provided", func() {
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
			Expect(err).NotTo(HaveOccurred())
			// Custom profile name
			Expect(content).To(ContainSubstring("JANE DOE"))
			Expect(content).To(ContainSubstring("Principal Engineer"))
			Expect(content).To(ContainSubstring("jane@example.com"))
			Expect(content).To(ContainSubstring("Python, Rust, TypeScript"))
			Expect(content).To(ContainSubstring("Distributed systems design"))
			Expect(content).To(ContainSubstring("Deep technical expertise"))
		})

		It("should use custom profile in markdown format", func() {
			profileCfg := &config.ProfileConfig{
				Name:  "John Smith",
				Title: "Staff Engineer",
			}

			content, err := service.ExportWithProfile(ctx, cv, sections, bullets, CVStructureNarrative, ExportFormatMarkdown, profileCfg)
			Expect(err).NotTo(HaveOccurred())
			// Markdown uses title case for name
			Expect(content).To(ContainSubstring("# John Smith"))
			Expect(content).To(ContainSubstring("**Staff Engineer**"))
		})

		It("should keep empty fields when custom profile fields are empty", func() {
			profileCfg := &config.ProfileConfig{
				Name: "Custom Name",
				// All other fields empty - should remain empty (inference fills these)
			}

			content, err := service.ExportWithProfile(ctx, cv, sections, bullets, CVStructureNarrative, ExportFormatText, profileCfg)
			Expect(err).NotTo(HaveOccurred())
			// Custom name should be present
			Expect(content).To(ContainSubstring("CUSTOM NAME"))
			// Empty fields should NOT contain hardcoded personal data
			Expect(content).NotTo(ContainSubstring("Yomi"))
			Expect(content).NotTo(ContainSubstring("boodah"))
		})

		It("should use flat format with profile data for YAML format", func() {
			profileCfg := &config.ProfileConfig{
				FirstName: "YAMLTest",
				LastName:  "User",
				Email:     "test@yaml.com",
			}

			content, err := service.ExportWithProfile(ctx, cv, sections, bullets, CVStructureNarrative, ExportFormatYAML, profileCfg)
			Expect(err).NotTo(HaveOccurred())
			// YAML uses flat format with profile data
			Expect(content).To(ContainSubstring("first_name: YAMLTest"))
			Expect(content).To(ContainSubstring("last_name: User"))
			Expect(content).To(ContainSubstring("email: test@yaml.com"))
			Expect(content).To(ContainSubstring("jobs:"))
			Expect(content).To(ContainSubstring("skills:"))
			// Should NOT contain internal CVView format fields
			Expect(content).NotTo(ContainSubstring("name: Test CV"))
			Expect(content).NotTo(ContainSubstring("target_role:"))
		})
	})

	Describe("NarrativeProfileFromConfig", func() {
		It("should return empty default profile when config is nil", func() {
			profile := NarrativeProfileFromConfig(nil)
			// Default profile should be empty - no hardcoded personal data
			Expect(profile.Name).To(BeEmpty())
			Expect(profile.Role).To(BeEmpty())
			Expect(profile.CoreStrengths).To(BeEmpty())
			Expect(profile.ValuePropositions).To(BeEmpty())
		})

		It("should use config values when provided", func() {
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
			Expect(profile.Name).To(Equal("Test User"))
			Expect(profile.Role).To(Equal("Lead Developer"))
			Expect(profile.Location).To(Equal("London, UK"))
			Expect(profile.Email).To(Equal("test@example.com"))
			Expect(profile.GitHub).To(Equal("https://github.com/testuser"))
			Expect(profile.Portfolio).To(Equal("https://test.dev"))
			Expect(profile.Languages).To(Equal([]string{"Go", "Python"}))
			Expect(profile.Frontend).To(Equal([]string{"Angular"}))
			Expect(profile.Systems).To(Equal([]string{"Docker", "GCP"}))
			Expect(profile.CoreStrengths).To(Equal([]string{"Backend development", "System architecture"}))
			Expect(profile.ValuePropositions).To(Equal([]string{"Pragmatic approach", "Strong communication"}))
		})

		It("should keep empty fields empty (no hardcoded defaults)", func() {
			cfg := &config.ProfileConfig{
				Name: "Only Name Set",
				// All other fields empty - should remain empty for inference service to fill
			}

			profile := NarrativeProfileFromConfig(cfg)

			Expect(profile.Name).To(Equal("Only Name Set"))
			// All other fields should be empty (not hardcoded defaults)
			Expect(profile.Role).To(BeEmpty())
			Expect(profile.Location).To(BeEmpty())
			Expect(profile.Email).To(BeEmpty())
			Expect(profile.GitHub).To(BeEmpty())
			Expect(profile.Portfolio).To(BeEmpty())
			Expect(profile.Languages).To(BeEmpty())
			Expect(profile.Frontend).To(BeEmpty())
			Expect(profile.Systems).To(BeEmpty())
			Expect(profile.CoreStrengths).To(BeEmpty())
			Expect(profile.ValuePropositions).To(BeEmpty())
		})
	})

	Describe("maxHighlightsFromProfile", func() {
		Context("when profile is nil", func() {
			It("returns the default of 5", func() {
				Expect(maxHighlightsFromProfile(nil)).To(Equal(5))
			})
		})

		Context("when profile has MaxHighlights of 0", func() {
			It("returns the default of 5", func() {
				profile := &config.ProfileConfig{MaxHighlights: 0}
				Expect(maxHighlightsFromProfile(profile)).To(Equal(5))
			})
		})

		Context("when profile has a positive MaxHighlights", func() {
			It("returns the configured value", func() {
				profile := &config.ProfileConfig{MaxHighlights: 3}
				Expect(maxHighlightsFromProfile(profile)).To(Equal(3))
			})
		})
	})
})

var _ = Describe("ExportService YAML Export", func() {
	var (
		service       *ExportService
		log           *logger.Logger
		ctx           context.Context
		mockSkillRepo *mockrepo.MockSkillRepository
		ctrl          *gomock.Controller
	)

	BeforeEach(func() {
		log = logger.New(io.Discard, logger.InfoLevel)
		ctrl = gomock.NewController(GinkgoT())
		mockSkillRepo = mockrepo.NewMockSkillRepository(ctrl)
		ctx = context.Background()
	})

	AfterEach(func() {
		ctrl.Finish()
	})

	Describe("NewExportServiceWithDeps", func() {
		It("creates service with profileConfig and skillRepo set", func() {
			profileCfg := &config.ProfileConfig{
				FirstName: "Yomi",
				LastName:  "Colledge",
			}

			service = NewExportServiceWithDeps(log, profileCfg, mockSkillRepo)

			Expect(service).NotTo(BeNil())
			Expect(service.profileConfig).To(Equal(profileCfg))
			Expect(service.skillRepo).To(Equal(mockSkillRepo))
			Expect(service.logger).To(Equal(log))
		})

		It("works with nil profileConfig", func() {
			service = NewExportServiceWithDeps(log, nil, mockSkillRepo)

			Expect(service).NotTo(BeNil())
			Expect(service.profileConfig).To(BeNil())
			Expect(service.skillRepo).To(Equal(mockSkillRepo))
		})

		It("works with nil skillRepo", func() {
			profileCfg := &config.ProfileConfig{FirstName: "Test"}

			service = NewExportServiceWithDeps(log, profileCfg, nil)

			Expect(service).NotTo(BeNil())
			Expect(service.profileConfig).To(Equal(profileCfg))
			Expect(service.skillRepo).To(BeNil())
		})
	})

	Describe("ExportToYAML with flat format", func() {
		var (
			cv       *career.CVView
			sections []*career.CVSection
		)

		BeforeEach(func() {
			cv = fixtures.CVViewWith("cv-1", "Senior Engineer CV", "staff", "hiring_manager")
			cv.SourceEventCount = 5
			cv.SourceFactCount = 3

			bullet := fixtures.CVBulletWith("bullet-1", "section-1", "Led migration to microservices architecture")
			bullet.Confidence = 0.95

			expSection := fixtures.CVSectionWith("section-1", "cv-1", "experience", "Experience", 1)
			expSection.Content = []*career.SectionContentGroup{
				fixtures.ContentGroupFull("Acme Corp", "2020-01", "2023-12", []*career.CVBullet{bullet}),
			}

			projSection := fixtures.CVSectionWith("section-2", "cv-1", "projects", "Projects", 2)
			projSection.Content = []*career.SectionContentGroup{
				fixtures.ContentGroupWithBullets("Open Source Project", []*career.CVBullet{
					fixtures.CVBulletWith("bullet-2", "section-2", "Built X using Y"),
				}),
			}

			sections = []*career.CVSection{expSection, projSection}
		})

		It("returns valid YAML string (no error)", func() {
			profileCfg := &config.ProfileConfig{
				FirstName: "Yomi",
				LastName:  "Colledge",
				Email:     "yomi@boodah.net",
				Country:   "Remote (UK)",
				Phone:     "+44 7894 987 855",
				LinkedIn:  "https://www.linkedin.com/in/yomicolledge",
				GitHub:    "https://github.com/baphled",
				Portfolio: "http://boodah.net",
			}

			mockSkillRepo.EXPECT().List(gomock.Any(), gomock.Nil()).Return([]*career.Skill{
				fixtures.SkillWith("s1", "Go", "backend", ""),
				fixtures.SkillWith("s2", "Ruby", "backend", ""),
				fixtures.SkillWith("s3", "React", "frontend", ""),
			}, nil)

			service = NewExportServiceWithDeps(log, profileCfg, mockSkillRepo)

			yamlOutput, err := service.ExportToYAML(ctx, cv, sections, nil)
			Expect(err).NotTo(HaveOccurred())
			Expect(yamlOutput).NotTo(BeEmpty())
		})

		It("contains first_name, last_name, email from profileConfig/cv", func() {
			profileCfg := &config.ProfileConfig{
				FirstName: "Jane",
				LastName:  "Doe",
				Email:     "jane@example.com",
				Country:   "London, UK",
			}

			mockSkillRepo.EXPECT().List(gomock.Any(), gomock.Nil()).Return([]*career.Skill{}, nil)

			service = NewExportServiceWithDeps(log, profileCfg, mockSkillRepo)

			yamlOutput, err := service.ExportToYAML(ctx, cv, sections, nil)
			Expect(err).NotTo(HaveOccurred())
			Expect(yamlOutput).To(ContainSubstring("first_name: Jane"))
			Expect(yamlOutput).To(ContainSubstring("last_name: Doe"))
			Expect(yamlOutput).To(ContainSubstring("email: jane@example.com"))
			Expect(yamlOutput).To(ContainSubstring("location: London, UK"))
		})

		It("contains links section with LinkedIn entry when LinkedIn is set", func() {
			profileCfg := &config.ProfileConfig{
				FirstName: "Test",
				LastName:  "User",
				LinkedIn:  "https://www.linkedin.com/in/testuser",
				GitHub:    "https://github.com/testuser",
			}

			mockSkillRepo.EXPECT().List(gomock.Any(), gomock.Nil()).Return([]*career.Skill{}, nil)

			service = NewExportServiceWithDeps(log, profileCfg, mockSkillRepo)

			yamlOutput, err := service.ExportToYAML(ctx, cv, sections, nil)
			Expect(err).NotTo(HaveOccurred())
			Expect(yamlOutput).To(ContainSubstring("links:"))
			Expect(yamlOutput).To(ContainSubstring("url: https://www.linkedin.com/in/testuser"))
			Expect(yamlOutput).To(ContainSubstring("label: LinkedIn"))
			Expect(yamlOutput).To(ContainSubstring("url: https://github.com/testuser"))
			Expect(yamlOutput).To(ContainSubstring("label: GitHub"))
		})

		It("expands bare GitHub username to full URL", func() {
			profileCfg := &config.ProfileConfig{
				FirstName: "Test",
				LastName:  "User",
				GitHub:    "baphled",
			}

			mockSkillRepo.EXPECT().List(gomock.Any(), gomock.Nil()).Return([]*career.Skill{}, nil)

			service = NewExportServiceWithDeps(log, profileCfg, mockSkillRepo)

			yamlOutput, err := service.ExportToYAML(ctx, cv, sections, nil)
			Expect(err).NotTo(HaveOccurred())
			Expect(yamlOutput).To(ContainSubstring("url: https://github.com/baphled"))
			Expect(yamlOutput).To(ContainSubstring("label: GitHub"))
		})

		It("expands bare LinkedIn username to full URL", func() {
			profileCfg := &config.ProfileConfig{
				FirstName: "Test",
				LastName:  "User",
				LinkedIn:  "yomicolledge",
			}

			mockSkillRepo.EXPECT().List(gomock.Any(), gomock.Nil()).Return([]*career.Skill{}, nil)

			service = NewExportServiceWithDeps(log, profileCfg, mockSkillRepo)

			yamlOutput, err := service.ExportToYAML(ctx, cv, sections, nil)
			Expect(err).NotTo(HaveOccurred())
			Expect(yamlOutput).To(ContainSubstring("url: https://www.linkedin.com/in/yomicolledge"))
			Expect(yamlOutput).To(ContainSubstring("label: LinkedIn"))
		})

		It("leaves portfolio URL unchanged (already full URL)", func() {
			profileCfg := &config.ProfileConfig{
				FirstName: "Test",
				LastName:  "User",
				Portfolio: "http://boodah.net",
			}

			mockSkillRepo.EXPECT().List(gomock.Any(), gomock.Nil()).Return([]*career.Skill{}, nil)

			service = NewExportServiceWithDeps(log, profileCfg, mockSkillRepo)

			yamlOutput, err := service.ExportToYAML(ctx, cv, sections, nil)
			Expect(err).NotTo(HaveOccurred())
			Expect(yamlOutput).To(ContainSubstring("url: http://boodah.net"))
			Expect(yamlOutput).To(ContainSubstring("label: Portfolio"))
		})

		It("omits empty links from output", func() {
			profileCfg := &config.ProfileConfig{
				FirstName: "Test",
				LastName:  "User",
				GitHub:    "",
				LinkedIn:  "",
				Portfolio: "",
			}

			mockSkillRepo.EXPECT().List(gomock.Any(), gomock.Nil()).Return([]*career.Skill{}, nil)

			service = NewExportServiceWithDeps(log, profileCfg, mockSkillRepo)

			yamlOutput, err := service.ExportToYAML(ctx, cv, sections, nil)
			Expect(err).NotTo(HaveOccurred())
			Expect(yamlOutput).To(ContainSubstring("links: []"))
		})

		It("contains skills section grouped by category", func() {
			profileCfg := &config.ProfileConfig{
				FirstName: "Test",
				LastName:  "User",
			}

			mockSkillRepo.EXPECT().List(gomock.Any(), gomock.Nil()).Return([]*career.Skill{
				fixtures.SkillWith("s1", "Go", "backend", ""),
				fixtures.SkillWith("s2", "Ruby", "backend", ""),
				fixtures.SkillWith("s3", "React", "frontend", ""),
				fixtures.SkillWith("s4", "TypeScript", "frontend", ""),
			}, nil)

			service = NewExportServiceWithDeps(log, profileCfg, mockSkillRepo)

			yamlOutput, err := service.ExportToYAML(ctx, cv, sections, nil)
			Expect(err).NotTo(HaveOccurred())
			Expect(yamlOutput).To(ContainSubstring("skills:"))
			Expect(yamlOutput).To(ContainSubstring("backend:"))
			Expect(yamlOutput).To(ContainSubstring("- Go"))
			Expect(yamlOutput).To(ContainSubstring("- Ruby"))
			Expect(yamlOutput).To(ContainSubstring("frontend:"))
			Expect(yamlOutput).To(ContainSubstring("- React"))
			Expect(yamlOutput).To(ContainSubstring("- TypeScript"))
		})

		Context("when ProfileConfig.SkillsLimit is set", func() {
			It("respects SkillsLimit of 3, limiting each category to at most 3 entries", func() {
				profileCfg := &config.ProfileConfig{
					FirstName:   "Test",
					LastName:    "User",
					SkillsLimit: 3,
				}

				mockSkillRepo.EXPECT().List(gomock.Any(), gomock.Nil()).Return([]*career.Skill{
					fixtures.SkillWith("s1", "Go", "backend", ""),
					fixtures.SkillWith("s2", "Ruby", "backend", ""),
					fixtures.SkillWith("s3", "Python", "backend", ""),
					fixtures.SkillWith("s4", "Java", "backend", ""),
					fixtures.SkillWith("s5", "Rust", "backend", ""),
					fixtures.SkillWith("s6", "C++", "backend", ""),
					fixtures.SkillWith("s7", "React", "frontend", ""),
					fixtures.SkillWith("s8", "Vue", "frontend", ""),
					fixtures.SkillWith("s9", "Angular", "frontend", ""),
					fixtures.SkillWith("s10", "Svelte", "frontend", ""),
				}, nil)

				service = NewExportServiceWithDeps(log, profileCfg, mockSkillRepo)

				yamlOutput, err := service.ExportToYAML(ctx, cv, sections, profileCfg)
				Expect(err).NotTo(HaveOccurred())

				Expect(yamlOutput).To(ContainSubstring("- Go"))
				Expect(yamlOutput).To(ContainSubstring("- Ruby"))
				Expect(yamlOutput).To(ContainSubstring("- Python"))
				Expect(yamlOutput).NotTo(ContainSubstring("- Java"))
				Expect(yamlOutput).NotTo(ContainSubstring("- Rust"))
				Expect(yamlOutput).NotTo(ContainSubstring("- C++"))

				Expect(yamlOutput).To(ContainSubstring("- React"))
				Expect(yamlOutput).To(ContainSubstring("- Vue"))
				Expect(yamlOutput).To(ContainSubstring("- Angular"))
				Expect(yamlOutput).NotTo(ContainSubstring("- Svelte"))
			})

			It("uses default limit of 5 when SkillsLimit is 0 (unset)", func() {
				profileCfg := &config.ProfileConfig{
					FirstName:   "Test",
					LastName:    "User",
					SkillsLimit: 0,
				}

				mockSkillRepo.EXPECT().List(gomock.Any(), gomock.Nil()).Return([]*career.Skill{
					fixtures.SkillWith("s1", "Go", "backend", ""),
					fixtures.SkillWith("s2", "Ruby", "backend", ""),
					fixtures.SkillWith("s3", "Python", "backend", ""),
					fixtures.SkillWith("s4", "Java", "backend", ""),
					fixtures.SkillWith("s5", "Rust", "backend", ""),
					fixtures.SkillWith("s6", "C++", "backend", ""),
					fixtures.SkillWith("s7", "Kotlin", "backend", ""),
				}, nil)

				service = NewExportServiceWithDeps(log, profileCfg, mockSkillRepo)

				yamlOutput, err := service.ExportToYAML(ctx, cv, sections, profileCfg)
				Expect(err).NotTo(HaveOccurred())

				Expect(yamlOutput).To(ContainSubstring("- Go"))
				Expect(yamlOutput).To(ContainSubstring("- Ruby"))
				Expect(yamlOutput).To(ContainSubstring("- Python"))
				Expect(yamlOutput).To(ContainSubstring("- Java"))
				Expect(yamlOutput).To(ContainSubstring("- Rust"))
				Expect(yamlOutput).NotTo(ContainSubstring("- C++"))
				Expect(yamlOutput).NotTo(ContainSubstring("- Kotlin"))
			})
		})

		It("returns error for nil CV", func() {
			profileCfg := &config.ProfileConfig{FirstName: "Test"}
			service = NewExportServiceWithDeps(log, profileCfg, mockSkillRepo)

			yamlOutput, err := service.ExportToYAML(ctx, nil, sections, nil)
			Expect(err).To(HaveOccurred())
			Expect(err.Error()).To(ContainSubstring("CV view is nil"))
			Expect(yamlOutput).To(BeEmpty())
		})

		It("contains jobs section with company/dates/description", func() {
			profileCfg := &config.ProfileConfig{FirstName: "Test", LastName: "User"}

			mockSkillRepo.EXPECT().List(gomock.Any(), gomock.Nil()).Return([]*career.Skill{}, nil)

			service = NewExportServiceWithDeps(log, profileCfg, mockSkillRepo)

			yamlOutput, err := service.ExportToYAML(ctx, cv, sections, nil)
			Expect(err).NotTo(HaveOccurred())
			Expect(yamlOutput).To(ContainSubstring("jobs:"))
			Expect(yamlOutput).To(ContainSubstring("company: Acme Corp"))
			Expect(yamlOutput).To(ContainSubstring("start_date: 2020-01"))
			Expect(yamlOutput).To(ContainSubstring("end_date: 2023-12"))
			Expect(yamlOutput).To(ContainSubstring("description: |"))
			Expect(yamlOutput).To(ContainSubstring("- Led migration to microservices architecture"))
		})

		It("contains projects section", func() {
			profileCfg := &config.ProfileConfig{FirstName: "Test", LastName: "User"}

			mockSkillRepo.EXPECT().List(gomock.Any(), gomock.Nil()).Return([]*career.Skill{}, nil)

			service = NewExportServiceWithDeps(log, profileCfg, mockSkillRepo)

			yamlOutput, err := service.ExportToYAML(ctx, cv, sections, nil)
			Expect(err).NotTo(HaveOccurred())
			Expect(yamlOutput).To(ContainSubstring("projects:"))
			Expect(yamlOutput).To(ContainSubstring("name: Open Source Project"))
			Expect(yamlOutput).To(ContainSubstring("key_achievements:"))
		})

		It("populates summary field from summary-type section", func() {
			profileCfg := &config.ProfileConfig{FirstName: "Test", LastName: "User"}

			mockSkillRepo.EXPECT().List(gomock.Any(), gomock.Nil()).Return([]*career.Skill{}, nil)

			summarySection := fixtures.CVSectionWithSummary("section-summary", "cv-1", "Experienced engineer with 10+ years in backend development.")
			sectionsWithSummary := append([]*career.CVSection{summarySection}, sections...)

			service = NewExportServiceWithDeps(log, profileCfg, mockSkillRepo)

			yamlOutput, err := service.ExportToYAML(ctx, cv, sectionsWithSummary, nil)
			Expect(err).NotTo(HaveOccurred())
			Expect(yamlOutput).To(ContainSubstring("summary: Experienced engineer with 10+ years in backend development."))
		})

		It("matches golden file output structure", func() {
			profileCfg := &config.ProfileConfig{
				FirstName: "Yomi",
				LastName:  "Colledge",
				Email:     "yomi@boodah.net",
				Country:   "Remote (UK)",
				Phone:     "+44 7894 987 855",
				LinkedIn:  "https://www.linkedin.com/in/yomicolledge",
				GitHub:    "https://github.com/baphled",
				Portfolio: "http://boodah.net",
			}

			mockSkillRepo.EXPECT().List(gomock.Any(), gomock.Nil()).Return([]*career.Skill{
				fixtures.SkillWith("s1", "Go", "backend", ""),
				fixtures.SkillWith("s2", "Ruby", "backend", ""),
				fixtures.SkillWith("s3", "React", "frontend", ""),
			}, nil)

			summarySection := fixtures.CVSectionWithSummary("section-summary", "cv-1", "**GDS-Aligned Ruby on Rails Developer | 15+ Years Production Experience**\nFull-stack engineer with deep backend expertise in Go, Ruby, and Node.js with a track record of scaling MVPs to production.")
			goldenSections := append([]*career.CVSection{summarySection}, sections...)

			service = NewExportServiceWithDeps(log, profileCfg, mockSkillRepo)

			yamlOutput, err := service.ExportToYAML(ctx, cv, goldenSections, nil)
			Expect(err).NotTo(HaveOccurred())

			goldenFile := filepath.Join("testdata", "yaml_export_golden.yaml")
			goldenContent, err := os.ReadFile(goldenFile)
			Expect(err).NotTo(HaveOccurred(), "golden file should exist")

			normalizeLineEndings := func(s string) string {
				return strings.ReplaceAll(s, "\r\n", "\n")
			}
			Expect(normalizeLineEndings(strings.TrimSpace(yamlOutput))).To(Equal(normalizeLineEndings(strings.TrimSpace(string(goldenContent)))))
		})

		It("handles nil profileConfig gracefully", func() {
			mockSkillRepo.EXPECT().List(gomock.Any(), gomock.Nil()).Return([]*career.Skill{}, nil)

			service = NewExportServiceWithDeps(log, nil, mockSkillRepo)

			yamlOutput, err := service.ExportToYAML(ctx, cv, sections, nil)
			Expect(err).NotTo(HaveOccurred())
			Expect(yamlOutput).To(ContainSubstring("first_name: \"\""))
		})

		It("handles nil skillRepo gracefully", func() {
			profileCfg := &config.ProfileConfig{FirstName: "Test", LastName: "User"}

			service = NewExportServiceWithDeps(log, profileCfg, nil)

			yamlOutput, err := service.ExportToYAML(ctx, cv, sections, nil)
			Expect(err).NotTo(HaveOccurred())
			Expect(yamlOutput).To(ContainSubstring("skills: {}"))
		})

		It("handles skillRepo error gracefully", func() {
			profileCfg := &config.ProfileConfig{FirstName: "Test", LastName: "User"}

			mockSkillRepo.EXPECT().List(gomock.Any(), gomock.Nil()).Return(nil, errors.New("db error"))

			service = NewExportServiceWithDeps(log, profileCfg, mockSkillRepo)

			yamlOutput, err := service.ExportToYAML(ctx, cv, sections, nil)
			Expect(err).NotTo(HaveOccurred())
			Expect(yamlOutput).To(ContainSubstring("skills: {}"))
		})
	})

	Describe("Dynamic highlights generation", func() {
		var (
			cv       *career.CVView
			sections []*career.CVSection
		)

		BeforeEach(func() {
			cv = fixtures.CVViewWith("cv-1", "Test CV", "senior_ic", "hiring_manager")
			cv.SourceEventCount = 5
			cv.SourceFactCount = 3
		})

		Context("when WhatIBring is populated", func() {
			It("uses WhatIBring items as highlights", func() {
				profileCfg := &config.ProfileConfig{
					FirstName:  "Test",
					LastName:   "User",
					WhatIBring: []string{"15+ years Ruby", "GDS experience"},
				}

				mockSkillRepo.EXPECT().List(gomock.Any(), gomock.Nil()).Return([]*career.Skill{}, nil)

				service = NewExportServiceWithDeps(log, profileCfg, mockSkillRepo)

				expSection := fixtures.CVSectionWith("section-1", "cv-1", "experience", "Experience", 1)
				bullet := fixtures.CVBulletWithScores("bullet-1", "section-1", "Some experience bullet", 0.9, 0.95, 0.8, 0.7)
				expSection.Content = []*career.SectionContentGroup{
					fixtures.ContentGroupWithBullets("Acme Corp", []*career.CVBullet{bullet}),
				}
				sections = []*career.CVSection{expSection}

				yamlOutput, err := service.ExportToYAML(ctx, cv, sections, profileCfg)
				Expect(err).NotTo(HaveOccurred())
				Expect(yamlOutput).To(ContainSubstring("- 15+ years Ruby"))
				Expect(yamlOutput).To(ContainSubstring("- GDS experience"))
			})
		})

		Context("when WhatIBring is empty but experience bullets exist", func() {
			It("uses top-scored experience bullets as highlights", func() {
				profileCfg := &config.ProfileConfig{
					FirstName:     "Test",
					LastName:      "User",
					WhatIBring:    []string{},
					CoreStrengths: []string{"Fallback strength"},
				}

				mockSkillRepo.EXPECT().List(gomock.Any(), gomock.Nil()).Return([]*career.Skill{}, nil)

				service = NewExportServiceWithDeps(log, profileCfg, mockSkillRepo)

				highConfBullet := fixtures.CVBulletWithScores("bullet-1", "section-1", "High confidence achievement", 0.9, 0.95, 0.8, 0.7)
				medConfBullet := fixtures.CVBulletWithScores("bullet-2", "section-1", "Medium confidence work", 0.7, 0.75, 0.6, 0.5)
				lowConfBullet := fixtures.CVBulletWithScores("bullet-3", "section-1", "Low confidence task", 0.5, 0.55, 0.4, 0.3)

				expSection := fixtures.CVSectionWith("section-1", "cv-1", "experience", "Experience", 1)
				expSection.Content = []*career.SectionContentGroup{
					fixtures.ContentGroupWithBullets("Acme Corp", []*career.CVBullet{lowConfBullet, highConfBullet, medConfBullet}),
				}
				sections = []*career.CVSection{expSection}

				yamlOutput, err := service.ExportToYAML(ctx, cv, sections, profileCfg)
				Expect(err).NotTo(HaveOccurred())
				Expect(yamlOutput).To(ContainSubstring("- High confidence achievement"))
				Expect(yamlOutput).To(ContainSubstring("- Medium confidence work"))
				Expect(yamlOutput).To(ContainSubstring("- Low confidence task"))
				Expect(yamlOutput).NotTo(ContainSubstring("Fallback strength"))
			})

			It("sorts bullets by confidence descending with RoleScore as tiebreaker", func() {
				profileCfg := &config.ProfileConfig{
					FirstName:  "Test",
					LastName:   "User",
					WhatIBring: []string{},
				}

				mockSkillRepo.EXPECT().List(gomock.Any(), gomock.Nil()).Return([]*career.Skill{}, nil)

				service = NewExportServiceWithDeps(log, profileCfg, mockSkillRepo)

				bullet1 := fixtures.CVBulletWithScores("bullet-1", "section-1", "Bullet A (conf=0.9, role=0.8)", 0.9, 0.90, 0.80, 0.7)
				bullet2 := fixtures.CVBulletWithScores("bullet-2", "section-1", "Bullet B (conf=0.9, role=0.9)", 0.8, 0.90, 0.90, 0.6)
				bullet3 := fixtures.CVBulletWithScores("bullet-3", "section-1", "Bullet C (conf=0.8)", 0.7, 0.80, 0.70, 0.5)

				expSection := fixtures.CVSectionWith("section-1", "cv-1", "experience", "Experience", 1)
				expSection.Content = []*career.SectionContentGroup{
					fixtures.ContentGroupWithBullets("Acme Corp", []*career.CVBullet{bullet3, bullet1, bullet2}),
				}
				sections = []*career.CVSection{expSection}

				yamlOutput, err := service.ExportToYAML(ctx, cv, sections, profileCfg)
				Expect(err).NotTo(HaveOccurred())

				highlightsIdx := strings.Index(yamlOutput, "highlights:")
				Expect(highlightsIdx).To(BeNumerically(">", 0))

				highlightsSection := yamlOutput[highlightsIdx:]
				idxB := strings.Index(highlightsSection, "Bullet B")
				idxA := strings.Index(highlightsSection, "Bullet A")
				idxC := strings.Index(highlightsSection, "Bullet C")

				Expect(idxB).To(BeNumerically("<", idxA), "Bullet B should come before Bullet A")
				Expect(idxA).To(BeNumerically("<", idxC), "Bullet A should come before Bullet C")
			})

			It("limits highlights to maxHighlightBullets", func() {
				profileCfg := &config.ProfileConfig{
					FirstName:  "Test",
					LastName:   "User",
					WhatIBring: []string{},
				}

				mockSkillRepo.EXPECT().List(gomock.Any(), gomock.Nil()).Return([]*career.Skill{}, nil)

				service = NewExportServiceWithDeps(log, profileCfg, mockSkillRepo)

				var bullets []*career.CVBullet
				for i := range 12 {
					bullet := fixtures.CVBulletWithScores(
						fmt.Sprintf("bullet-%d", i),
						"section-1",
						fmt.Sprintf("Bullet number %d", i),
						0.9,
						0.9-float64(i)*0.05,
						0.8,
						0.7,
					)
					bullets = append(bullets, bullet)
				}

				expSection := fixtures.CVSectionWith("section-1", "cv-1", "experience", "Experience", 1)
				expSection.Content = []*career.SectionContentGroup{
					fixtures.ContentGroupWithBullets("Acme Corp", bullets),
				}
				sections = []*career.CVSection{expSection}

				yamlOutput, err := service.ExportToYAML(ctx, cv, sections, profileCfg)
				Expect(err).NotTo(HaveOccurred())
				highlightsIdx := strings.Index(yamlOutput, "highlights:")
				Expect(highlightsIdx).To(BeNumerically(">", 0))
				jobsIdx := strings.Index(yamlOutput, "jobs:")
				Expect(jobsIdx).To(BeNumerically(">", highlightsIdx))
				highlightsSection := yamlOutput[highlightsIdx:jobsIdx]

				for i := range 5 {
					Expect(highlightsSection).To(ContainSubstring(fmt.Sprintf("Bullet number %d", i)), "should include top 5 bullets in highlights")
				}
				for i := 5; i < 12; i++ {
					Expect(highlightsSection).NotTo(ContainSubstring(fmt.Sprintf("Bullet number %d", i)), "should NOT include bullets beyond limit in highlights")
				}
			})
		})

		Context("when WhatIBring is empty and no experience bullets exist", func() {
			It("falls back to CoreStrengths", func() {
				profileCfg := &config.ProfileConfig{
					FirstName:     "Test",
					LastName:      "User",
					WhatIBring:    []string{},
					CoreStrengths: []string{"Backend development", "System architecture"},
				}

				mockSkillRepo.EXPECT().List(gomock.Any(), gomock.Nil()).Return([]*career.Skill{}, nil)

				service = NewExportServiceWithDeps(log, profileCfg, mockSkillRepo)

				projSection := fixtures.CVSectionWith("section-1", "cv-1", "projects", "Projects", 1)
				projSection.Content = []*career.SectionContentGroup{}
				sections = []*career.CVSection{projSection}

				yamlOutput, err := service.ExportToYAML(ctx, cv, sections, profileCfg)
				Expect(err).NotTo(HaveOccurred())
				Expect(yamlOutput).To(ContainSubstring("- Backend development"))
				Expect(yamlOutput).To(ContainSubstring("- System architecture"))
			})
		})

		Context("when all sources are empty", func() {
			It("returns empty highlights", func() {
				profileCfg := &config.ProfileConfig{
					FirstName:     "Test",
					LastName:      "User",
					WhatIBring:    []string{},
					CoreStrengths: []string{},
				}

				mockSkillRepo.EXPECT().List(gomock.Any(), gomock.Nil()).Return([]*career.Skill{}, nil)

				service = NewExportServiceWithDeps(log, profileCfg, mockSkillRepo)

				projSection := fixtures.CVSectionWith("section-1", "cv-1", "projects", "Projects", 1)
				projSection.Content = []*career.SectionContentGroup{}
				sections = []*career.CVSection{projSection}

				yamlOutput, err := service.ExportToYAML(ctx, cv, sections, profileCfg)
				Expect(err).NotTo(HaveOccurred())
				Expect(yamlOutput).To(ContainSubstring("highlights: \"\""))
			})

			It("returns empty highlights when profile is nil", func() {
				service = NewExportServiceWithDeps(log, nil, nil)

				projSection := fixtures.CVSectionWith("section-1", "cv-1", "projects", "Projects", 1)
				projSection.Content = []*career.SectionContentGroup{}
				sections = []*career.CVSection{projSection}

				yamlOutput, err := service.ExportToYAML(ctx, cv, sections, nil)
				Expect(err).NotTo(HaveOccurred())
				Expect(yamlOutput).To(ContainSubstring("highlights: \"\""))
			})

		})

		Context("with audience-specific bullet ordering", func() {
			It("uses audience relevance to order highlights", func() {
				audService := NewExportServiceWithDeps(log, nil, nil)

				highAudBullet := fixtures.CVBulletWith("b-high-aud", "section-1", "Delivered key business outcome")
				highAudBullet.AudienceRelevance = map[string]float64{"hiring_manager": 0.99}
				highAudBullet.Confidence = 0.5

				lowAudBullet := fixtures.CVBulletWith("b-low-aud", "section-1", "Built internal tool")
				lowAudBullet.AudienceRelevance = map[string]float64{"hiring_manager": 0.1}
				lowAudBullet.Confidence = 0.9

				expSection := fixtures.CVSectionWithContent("section-1", "cv-1", []*career.SectionContentGroup{
					fixtures.ContentGroupWithBullets("Senior at TechCo", []*career.CVBullet{lowAudBullet, highAudBullet}),
				})
				expSection.Title = "Professional Experience"
				expSection.SectionType = "experience"

				audCV := fixtures.CVViewWith("cv-1", "Test CV", "Engineer", "hiring_manager")
				audCV.Sections = []*career.CVSection{expSection}

				yamlOutput, err := audService.ExportToYAML(ctx, audCV, audCV.Sections, nil)
				Expect(err).NotTo(HaveOccurred())

				highlightsIdx := strings.Index(yamlOutput, "highlights:")
				Expect(highlightsIdx).To(BeNumerically(">", -1))
				highlightsSection := yamlOutput[highlightsIdx:]
				highIdx := strings.Index(highlightsSection, "Delivered key business outcome")
				lowIdx := strings.Index(highlightsSection, "Built internal tool")
				Expect(highIdx).To(BeNumerically(">", -1))
				Expect(lowIdx).To(BeNumerically(">", -1))
				Expect(highIdx).To(BeNumerically("<", lowIdx),
					"Audience-scored bullet should appear before confidence-only bullet")
			})
		})
	})
	Describe("maxHighlightsFromProfile", func() {
		It("returns 5 when profile is nil", func() {
			Expect(maxHighlightsFromProfile(nil)).To(Equal(5))
		})

		It("returns 5 when MaxHighlights is 0", func() {
			profile := &config.ProfileConfig{}
			Expect(maxHighlightsFromProfile(profile)).To(Equal(5))
		})

		It("returns the configured value", func() {
			profile := &config.ProfileConfig{MaxHighlights: 3}
			Expect(maxHighlightsFromProfile(profile)).To(Equal(3))
		})
	})
})

var _ = Describe("ExportService Coverage", func() {
	var (
		service *ExportService
		log     *logger.Logger
		ctx     context.Context
	)

	BeforeEach(func() {
		log = logger.New(io.Discard, logger.InfoLevel)
		service = NewExportService(log, nil, nil)
		ctx = context.Background()
	})

	Describe("SystemClipboard", func() {
		Describe("WriteAll", func() {
			It("should attempt clipboard write without panicking", func() {
				sc := &SystemClipboard{}
				err := sc.WriteAll("test content")
				_ = err
			})
		})

		Describe("IsUnsupported", func() {
			It("should return a boolean without panicking", func() {
				sc := &SystemClipboard{}
				_ = sc.IsUnsupported()
			})
		})
	})

	Describe("sortBulletsByConfidenceAndRoleScore", func() {
		It("should sort by confidence descending", func() {
			low := fixtures.CVBulletWithScores("b-low", "s1", "Low", 0.5, 0.5, 0.5, 0.0)
			high := fixtures.CVBulletWithScores("b-high", "s1", "High", 0.5, 0.9, 0.5, 0.0)
			mid := fixtures.CVBulletWithScores("b-mid", "s1", "Mid", 0.5, 0.7, 0.5, 0.0)
			bullets := []*career.CVBullet{low, high, mid}
			sortBulletsByConfidenceAndRoleScore(bullets)
			Expect(bullets[0].Text).To(Equal("High"))
			Expect(bullets[1].Text).To(Equal("Mid"))
			Expect(bullets[2].Text).To(Equal("Low"))
		})

		It("should use RoleScore as tiebreaker when confidence is equal", func() {
			lowRole := fixtures.CVBulletWithScores("b-low-role", "s1", "LowRole", 0.5, 0.8, 0.3, 0.0)
			highRole := fixtures.CVBulletWithScores("b-high-role", "s1", "HighRole", 0.5, 0.8, 0.9, 0.0)
			midRole := fixtures.CVBulletWithScores("b-mid-role", "s1", "MidRole", 0.5, 0.8, 0.6, 0.0)
			bullets := []*career.CVBullet{lowRole, highRole, midRole}
			sortBulletsByConfidenceAndRoleScore(bullets)
			Expect(bullets[0].Text).To(Equal("HighRole"))
			Expect(bullets[1].Text).To(Equal("MidRole"))
			Expect(bullets[2].Text).To(Equal("LowRole"))
		})

		It("should return 0 when confidence and RoleScore are both equal", func() {
			a := fixtures.CVBulletWithScores("b-a", "s1", "A", 0.5, 0.8, 0.7, 0.0)
			b := fixtures.CVBulletWithScores("b-b", "s1", "B", 0.5, 0.8, 0.7, 0.0)
			bullets := []*career.CVBullet{a, b}
			sortBulletsByConfidenceAndRoleScore(bullets)
			Expect(bullets).To(HaveLen(2))
		})

		It("should handle empty slice without panic", func() {
			var bullets []*career.CVBullet
			sortBulletsByConfidenceAndRoleScore(bullets)
			Expect(bullets).To(BeEmpty())
		})

		It("should handle single element", func() {
			only := fixtures.CVBulletWithScores("b-only", "s1", "Only", 0.5, 0.8, 0.7, 0.0)
			bullets := []*career.CVBullet{only}
			sortBulletsByConfidenceAndRoleScore(bullets)
			Expect(bullets[0].Text).To(Equal("Only"))
		})
	})

	Describe("exportConsultingWithProfile unknown format", func() {
		It("should return error for unknown format", func() {
			cv := fixtures.CVViewWith("cv-1", "Test CV", "senior_ic", "hiring_manager")
			sections := []*career.CVSection{}
			bullets := map[string][]*career.CVBullet{}
			_, err := service.ExportWithProfile(ctx, cv, sections, bullets, CVStructureConsulting, ExportFormat("invalid"), nil)
			Expect(err).To(HaveOccurred())
			Expect(err.Error()).To(ContainSubstring("unknown export format"))
		})
	})

	Describe("exportHighlightsWithProfile unknown format", func() {
		It("should return error for unknown format", func() {
			cv := fixtures.CVViewWith("cv-1", "Test CV", "senior_ic", "hiring_manager")
			sections := []*career.CVSection{}
			bullets := map[string][]*career.CVBullet{}
			_, err := service.ExportWithProfile(ctx, cv, sections, bullets, CVStructureHighlights, ExportFormat("invalid"), nil)
			Expect(err).To(HaveOccurred())
			Expect(err.Error()).To(ContainSubstring("unknown export format"))
		})
	})

	Describe("exportNarrativeWithProfile unknown format", func() {
		It("should return error for unknown format", func() {
			cv := fixtures.CVViewWith("cv-1", "Test CV", "senior_ic", "hiring_manager")
			sections := []*career.CVSection{}
			_, err := service.ExportWithProfile(ctx, cv, sections, nil, CVStructureNarrative, ExportFormat("invalid"), nil)
			Expect(err).To(HaveOccurred())
			Expect(err.Error()).To(ContainSubstring("unknown export format"))
		})
	})

	Describe("exportHighlightsMarkdown branches", func() {
		It("should use default capabilities when no core strengths in profile", func() {
			cv := fixtures.CVViewWith("cv-1", "Test CV", "senior_ic", "hiring_manager")
			bullet := fixtures.CVBulletWith("b1", "s1", "Built something great")
			bullet.Confidence = 0.9
			section := fixtures.CVSectionWith("s1", "cv-1", "experience", "Experience", 1)
			section.Content = []*career.SectionContentGroup{
				fixtures.ContentGroupWithBullets("Corp", []*career.CVBullet{bullet}),
			}
			bullets := map[string][]*career.CVBullet{"s1": {bullet}}

			content, err := service.ExportWithProfile(ctx, cv, []*career.CVSection{section}, bullets, CVStructureHighlights, ExportFormatMarkdown, nil)
			Expect(err).NotTo(HaveOccurred())
			Expect(content).To(ContainSubstring("- Technical leadership and architecture"))
			Expect(content).To(ContainSubstring("- System design and optimization"))
			Expect(content).To(ContainSubstring("- Cross-functional collaboration"))
		})

		It("should limit core strengths to 6 in markdown highlights", func() {
			cv := fixtures.CVViewWith("cv-1", "Test CV", "senior_ic", "hiring_manager")
			section := fixtures.CVSectionWith("s1", "cv-1", "experience", "Experience", 1)
			section.Content = []*career.SectionContentGroup{}
			bullets := map[string][]*career.CVBullet{}
			profileCfg := &config.ProfileConfig{
				CoreStrengths: []string{"S1", "S2", "S3", "S4", "S5", "S6", "S7", "S8"},
			}

			content, err := service.ExportWithProfile(ctx, cv, []*career.CVSection{section}, bullets, CVStructureHighlights, ExportFormatMarkdown, profileCfg)
			Expect(err).NotTo(HaveOccurred())
			Expect(content).To(ContainSubstring("- S6"))
			Expect(content).NotTo(ContainSubstring("- S7"))
		})

		It("should include technologies section with both languages and systems", func() {
			cv := fixtures.CVViewWith("cv-1", "Test CV", "senior_ic", "hiring_manager")
			section := fixtures.CVSectionWith("s1", "cv-1", "experience", "Experience", 1)
			section.Content = []*career.SectionContentGroup{}
			bullets := map[string][]*career.CVBullet{}
			profileCfg := &config.ProfileConfig{
				Languages: []string{"Go", "Python"},
				Systems:   []string{"Kubernetes", "Docker"},
			}

			content, err := service.ExportWithProfile(ctx, cv, []*career.CVSection{section}, bullets, CVStructureHighlights, ExportFormatMarkdown, profileCfg)
			Expect(err).NotTo(HaveOccurred())
			Expect(content).To(ContainSubstring("## Technologies"))
			Expect(content).To(ContainSubstring("**Languages:** Go, Python"))
			Expect(content).To(ContainSubstring("**Systems:** Kubernetes, Docker"))
		})

		It("should include technologies with only languages", func() {
			cv := fixtures.CVViewWith("cv-1", "Test CV", "senior_ic", "hiring_manager")
			section := fixtures.CVSectionWith("s1", "cv-1", "experience", "Experience", 1)
			section.Content = []*career.SectionContentGroup{}
			bullets := map[string][]*career.CVBullet{}
			profileCfg := &config.ProfileConfig{
				Languages: []string{"Go"},
			}

			content, err := service.ExportWithProfile(ctx, cv, []*career.CVSection{section}, bullets, CVStructureHighlights, ExportFormatMarkdown, profileCfg)
			Expect(err).NotTo(HaveOccurred())
			Expect(content).To(ContainSubstring("## Technologies"))
			Expect(content).To(ContainSubstring("**Languages:** Go"))
			Expect(content).NotTo(ContainSubstring("**Systems:**"))
		})

		It("should include technologies with only systems", func() {
			cv := fixtures.CVViewWith("cv-1", "Test CV", "senior_ic", "hiring_manager")
			section := fixtures.CVSectionWith("s1", "cv-1", "experience", "Experience", 1)
			section.Content = []*career.SectionContentGroup{}
			bullets := map[string][]*career.CVBullet{}
			profileCfg := &config.ProfileConfig{
				Systems: []string{"AWS"},
			}

			content, err := service.ExportWithProfile(ctx, cv, []*career.CVSection{section}, bullets, CVStructureHighlights, ExportFormatMarkdown, profileCfg)
			Expect(err).NotTo(HaveOccurred())
			Expect(content).To(ContainSubstring("## Technologies"))
			Expect(content).To(ContainSubstring("**Systems:** AWS"))
			Expect(content).NotTo(ContainSubstring("**Languages:**"))
		})

		It("should omit technologies section when no languages or systems", func() {
			cv := fixtures.CVViewWith("cv-1", "Test CV", "senior_ic", "hiring_manager")
			section := fixtures.CVSectionWith("s1", "cv-1", "experience", "Experience", 1)
			section.Content = []*career.SectionContentGroup{}
			bullets := map[string][]*career.CVBullet{}

			content, err := service.ExportWithProfile(ctx, cv, []*career.CVSection{section}, bullets, CVStructureHighlights, ExportFormatMarkdown, nil)
			Expect(err).NotTo(HaveOccurred())
			Expect(content).NotTo(ContainSubstring("## Technologies"))
		})

		It("should include summary when present", func() {
			cv := fixtures.CVViewWith("cv-1", "Test CV", "senior_ic", "hiring_manager")
			summarySection := fixtures.CVSectionWithSummary("s-summary", "cv-1", "A skilled engineer.")
			section := fixtures.CVSectionWith("s1", "cv-1", "experience", "Experience", 1)
			section.Content = []*career.SectionContentGroup{}
			bullets := map[string][]*career.CVBullet{}

			content, err := service.ExportWithProfile(ctx, cv, []*career.CVSection{summarySection, section}, bullets, CVStructureHighlights, ExportFormatMarkdown, nil)
			Expect(err).NotTo(HaveOccurred())
			Expect(content).To(ContainSubstring("A skilled engineer."))
		})

		It("should handle no summary gracefully", func() {
			cv := fixtures.CVViewWith("cv-1", "Test CV", "senior_ic", "hiring_manager")
			section := fixtures.CVSectionWith("s1", "cv-1", "experience", "Experience", 1)
			section.Content = []*career.SectionContentGroup{}
			bullets := map[string][]*career.CVBullet{}

			content, err := service.ExportWithProfile(ctx, cv, []*career.CVSection{section}, bullets, CVStructureHighlights, ExportFormatMarkdown, nil)
			Expect(err).NotTo(HaveOccurred())
			Expect(content).To(ContainSubstring("---"))
		})
	})

	Describe("exportConsultingMarkdown branches", func() {
		It("should include group header without dates", func() {
			cv := fixtures.CVViewWith("cv-1", "Test CV", "senior_ic", "hiring_manager")
			bullet := fixtures.CVBulletWith("b1", "s1", "Consulting work")
			section := fixtures.CVSectionWith("s1", "cv-1", "experience", "Experience", 1)
			section.Content = []*career.SectionContentGroup{
				fixtures.ContentGroupWithBullets("NoDates Corp", []*career.CVBullet{bullet}),
			}
			bullets := map[string][]*career.CVBullet{"s1": {bullet}}

			content, err := service.ExportWithProfile(ctx, cv, []*career.CVSection{section}, bullets, CVStructureConsulting, ExportFormatMarkdown, nil)
			Expect(err).NotTo(HaveOccurred())
			Expect(content).To(ContainSubstring("### NoDates Corp"))
		})

		It("should include What I Bring section when profile has value propositions", func() {
			cv := fixtures.CVViewWith("cv-1", "Test CV", "senior_ic", "hiring_manager")
			section := fixtures.CVSectionWith("s1", "cv-1", "experience", "Experience", 1)
			section.Content = []*career.SectionContentGroup{}
			bullets := map[string][]*career.CVBullet{}
			profileCfg := &config.ProfileConfig{
				WhatIBring: []string{"Deep expertise", "Strong leadership"},
			}

			content, err := service.ExportWithProfile(ctx, cv, []*career.CVSection{section}, bullets, CVStructureConsulting, ExportFormatMarkdown, profileCfg)
			Expect(err).NotTo(HaveOccurred())
			Expect(content).To(ContainSubstring("## What I Bring"))
			Expect(content).To(ContainSubstring("- Deep expertise"))
			Expect(content).To(ContainSubstring("- Strong leadership"))
		})

		It("should omit What I Bring when profile has no value propositions", func() {
			cv := fixtures.CVViewWith("cv-1", "Test CV", "senior_ic", "hiring_manager")
			section := fixtures.CVSectionWith("s1", "cv-1", "experience", "Experience", 1)
			section.Content = []*career.SectionContentGroup{}
			bullets := map[string][]*career.CVBullet{}

			content, err := service.ExportWithProfile(ctx, cv, []*career.CVSection{section}, bullets, CVStructureConsulting, ExportFormatMarkdown, nil)
			Expect(err).NotTo(HaveOccurred())
			Expect(content).NotTo(ContainSubstring("## What I Bring"))
		})

		It("should include summary in consulting markdown", func() {
			cv := fixtures.CVViewWith("cv-1", "Test CV", "senior_ic", "hiring_manager")
			summarySection := fixtures.CVSectionWithSummary("s-summary", "cv-1", "Expert consultant.")
			section := fixtures.CVSectionWith("s1", "cv-1", "experience", "Experience", 1)
			section.Content = []*career.SectionContentGroup{}
			bullets := map[string][]*career.CVBullet{}

			content, err := service.ExportWithProfile(ctx, cv, []*career.CVSection{summarySection, section}, bullets, CVStructureConsulting, ExportFormatMarkdown, nil)
			Expect(err).NotTo(HaveOccurred())
			Expect(content).To(ContainSubstring("Expert consultant."))
		})
	})

	Describe("ExportToMarkdown branches", func() {
		It("should handle same start and end date", func() {
			cv := fixtures.CVViewWith("cv-1", "Test CV", "Engineer", "hiring_manager")
			bullet := fixtures.CVBulletWith("b1", "s1", "Short project work")
			section := fixtures.CVSectionWith("s1", "cv-1", "experience", "Experience", 1)
			section.Content = []*career.SectionContentGroup{
				fixtures.ContentGroupFull("ShortCorp", "Jun 2023", "Jun 2023", []*career.CVBullet{bullet}),
			}

			markdown, err := service.ExportToMarkdown(ctx, cv, []*career.CVSection{section}, map[string][]*career.CVBullet{"s1": {bullet}})
			Expect(err).NotTo(HaveOccurred())
			Expect(markdown).To(ContainSubstring("### ShortCorp - _Jun 2023_"))
		})

		It("should handle content group header without dates", func() {
			cv := fixtures.CVViewWith("cv-1", "Test CV", "Engineer", "hiring_manager")
			bullet := fixtures.CVBulletWith("b1", "s1", "Undated work")
			section := fixtures.CVSectionWith("s1", "cv-1", "experience", "Experience", 1)
			section.Content = []*career.SectionContentGroup{
				fixtures.ContentGroupWithBullets("NoDates Corp", []*career.CVBullet{bullet}),
			}

			markdown, err := service.ExportToMarkdown(ctx, cv, []*career.CVSection{section}, map[string][]*career.CVBullet{"s1": {bullet}})
			Expect(err).NotTo(HaveOccurred())
			Expect(markdown).To(ContainSubstring("### NoDates Corp"))
		})

		It("should render summary section as prose", func() {
			cv := fixtures.CVViewWith("cv-1", "Test CV", "Engineer", "hiring_manager")
			summarySection := fixtures.CVSectionWithSummary("s-summary", "cv-1", "A talented engineer.")

			markdown, err := service.ExportToMarkdown(ctx, cv, []*career.CVSection{summarySection}, map[string][]*career.CVBullet{})
			Expect(err).NotTo(HaveOccurred())
			Expect(markdown).To(ContainSubstring("A talented engineer."))
		})
	})

	Describe("SaveToFile additional coverage", func() {
		It("should default to .txt extension for unknown format", func() {
			content := "Unknown format content"
			filePath, err := service.SaveToFile(ctx, "TestUnknown", ExportFormat("unsupported"), content)
			Expect(err).NotTo(HaveOccurred())
			Expect(filePath).To(ContainSubstring(".txt"))
			os.Remove(filePath)
		})

		It("should persist correct content to disk", func() {
			content := "Verified CV\nLine 2\nLine 3"
			filePath, err := service.SaveToFile(ctx, "VerifyCV", ExportFormatText, content)
			Expect(err).NotTo(HaveOccurred())
			savedContent, readErr := os.ReadFile(filePath)
			Expect(readErr).NotTo(HaveOccurred())
			Expect(string(savedContent)).To(Equal(content))
			os.Remove(filePath)
		})
	})

	Describe("GetExportPath structure", func() {
		It("should return path under home with correct subdirectories", func() {
			path, err := service.GetExportPath()
			Expect(err).NotTo(HaveOccurred())
			Expect(path).To(ContainSubstring(filepath.Join(".kariya", "cv_exports")))
		})
	})

	Describe("exportNarrativeWithProfile markdown", func() {
		It("should export narrative markdown with complete profile", func() {
			cv := fixtures.CVViewWith("cv-1", "Test CV", "senior_ic", "hiring_manager")
			bullet := fixtures.CVBulletWith("b1", "s1", "Major achievement")
			bullet.Confidence = 0.90
			expSection := fixtures.CVSectionWith("s1", "cv-1", "experience", "Experience", 1)
			expSection.Content = []*career.SectionContentGroup{
				fixtures.ContentGroupFull("BigCo", "Jan 2020", "Dec 2023", []*career.CVBullet{bullet}),
			}

			profileCfg := &config.ProfileConfig{
				Name:          "Test Author",
				Title:         "Staff Engineer",
				Location:      "London, UK",
				Email:         "test@example.com",
				GitHub:        "testuser",
				Portfolio:     "https://test.dev",
				Languages:     []string{"Go", "Python"},
				Frontend:      []string{"React"},
				Systems:       []string{"K8s"},
				CoreStrengths: []string{"Architecture"},
				WhatIBring:    []string{"Deep knowledge"},
			}

			content, err := service.ExportWithProfile(ctx, cv, []*career.CVSection{expSection}, map[string][]*career.CVBullet{"s1": {bullet}}, CVStructureNarrative, ExportFormatMarkdown, profileCfg)
			Expect(err).NotTo(HaveOccurred())
			Expect(content).To(ContainSubstring("# Test Author"))
			Expect(content).To(ContainSubstring("**Staff Engineer**"))
			Expect(content).To(ContainSubstring("## Core Strengths"))
			Expect(content).To(ContainSubstring("- Architecture"))
			Expect(content).To(ContainSubstring("## Languages & Technologies"))
			Expect(content).To(ContainSubstring("**Languages:** Go, Python"))
			Expect(content).To(ContainSubstring("## Selected Experience"))
			Expect(content).To(ContainSubstring("### BigCo"))
			Expect(content).To(ContainSubstring("- Major achievement"))
			Expect(content).To(ContainSubstring("## What I Bring"))
			Expect(content).To(ContainSubstring("- Deep knowledge"))
		})

		It("should handle narrative markdown with experience group without dates", func() {
			cv := fixtures.CVViewWith("cv-1", "Test CV", "senior_ic", "hiring_manager")
			bullet := fixtures.CVBulletWith("b1", "s1", "Undated achievement")
			bullet.Confidence = 0.90
			expSection := fixtures.CVSectionWith("s1", "cv-1", "experience", "Experience", 1)
			expSection.Content = []*career.SectionContentGroup{
				fixtures.ContentGroupWithBullets("NoDates Corp", []*career.CVBullet{bullet}),
			}

			profileCfg := &config.ProfileConfig{
				Name:  "Test Author",
				Title: "Engineer",
			}

			content, err := service.ExportWithProfile(ctx, cv, []*career.CVSection{expSection}, map[string][]*career.CVBullet{"s1": {bullet}}, CVStructureNarrative, ExportFormatMarkdown, profileCfg)
			Expect(err).NotTo(HaveOccurred())
			Expect(content).To(ContainSubstring("### NoDates Corp"))
		})

		It("should use default summary when no summary section exists", func() {
			cv := fixtures.CVViewWith("cv-1", "Test CV", "senior_ic", "hiring_manager")
			expSection := fixtures.CVSectionWith("s1", "cv-1", "experience", "Experience", 1)
			expSection.Content = []*career.SectionContentGroup{}

			content, err := service.ExportWithProfile(ctx, cv, []*career.CVSection{expSection}, map[string][]*career.CVBullet{}, CVStructureNarrative, ExportFormatMarkdown, nil)
			Expect(err).NotTo(HaveOccurred())
			Expect(content).To(ContainSubstring("Experienced software engineer with strong technical leadership skills."))
		})
	})
})
