package cv

import (
	"context"
	"errors"
	"io"
	"os"
	"path/filepath"
	"strings"
	"time"

	"github.com/baphled/kariya/internal/domain/career"
	"github.com/baphled/kariya/internal/logger"
	"github.com/onsi/ginkgo/v2"
	"github.com/onsi/gomega"
)

// MockClipboard is a test implementation of ClipboardWriter
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
			cv := &career.CVView{
				ID:               "cv-1",
				Name:             "Senior Software Engineer CV",
				TargetRole:       "Staff Engineer",
				TargetAudience:   "hiring_manager",
				GeneratedAt:      time.Now(),
				SourceEventCount: 10,
				SourceFactCount:  5,
			}

			bullet := &career.CVBullet{
				ID:              "bullet-1",
				SectionID:       "section-1",
				Text:            "Led team of 5 engineers to deliver critical feature",
				SourceEventIDs:  []string{"event-1"},
				SourceFactIDs:   []string{"fact-1"},
				Rank:            0.9,
				InclusionReason: "ownership",
				Confidence:      0.95,
			}

			section := &career.CVSection{
				ID:          "section-1",
				CVViewID:    "cv-1",
				SectionType: "experience",
				Title:       "Experience",
				Order:       1,
				Content: []*career.SectionContentGroup{
					{
						Header:  "Acme Corp",
						Bullets: []*career.CVBullet{bullet},
					},
				},
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
			cv := &career.CVView{
				ID:               "cv-1",
				Name:             "Test CV",
				TargetRole:       "Engineer",
				TargetAudience:   "hiring_manager",
				GeneratedAt:      time.Now(),
				SourceEventCount: 0,
				SourceFactCount:  0,
			}

			section := &career.CVSection{
				ID:          "section-1",
				CVViewID:    "cv-1",
				SectionType: "experience",
				Title:       "Experience",
				Order:       1,
				Content:     []*career.SectionContentGroup{},
			}

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
			cv := &career.CVView{
				ID:               "cv-1",
				Name:             "Senior Software Engineer CV",
				TargetRole:       "Staff Engineer",
				TargetAudience:   "hiring_manager",
				GeneratedAt:      time.Now(),
				SourceEventCount: 10,
				SourceFactCount:  5,
			}

			bullet := &career.CVBullet{
				ID:              "bullet-1",
				SectionID:       "section-1",
				Text:            "Led team of 5 engineers to deliver critical feature",
				SourceEventIDs:  []string{"event-1"},
				SourceFactIDs:   []string{"fact-1"},
				Rank:            0.9,
				InclusionReason: "ownership",
				Confidence:      0.95,
			}

			section := &career.CVSection{
				ID:          "section-1",
				CVViewID:    "cv-1",
				SectionType: "experience",
				Title:       "Experience",
				Order:       1,
				Content: []*career.SectionContentGroup{
					{
						Header:  "Acme Corp",
						Bullets: []*career.CVBullet{bullet},
					},
				},
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
			cv := &career.CVView{
				ID:               "cv-1",
				Name:             "Test CV",
				TargetRole:       "Engineer",
				TargetAudience:   "recruiter",
				GeneratedAt:      time.Now(),
				SourceEventCount: 5,
				SourceFactCount:  2,
			}

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
			cv := &career.CVView{
				ID:               "cv-1",
				Name:             "Senior Software Engineer CV",
				TargetRole:       "Staff Engineer",
				TargetAudience:   "hiring_manager",
				GeneratedAt:      time.Now(),
				SourceEventCount: 10,
				SourceFactCount:  5,
			}

			bullet := &career.CVBullet{
				ID:              "bullet-1",
				SectionID:       "section-1",
				Text:            "Led team of 5 engineers to deliver critical feature",
				SourceEventIDs:  []string{"event-1"},
				SourceFactIDs:   []string{"fact-1"},
				Rank:            0.9,
				InclusionReason: "ownership",
				Confidence:      0.95,
			}

			section := &career.CVSection{
				ID:          "section-1",
				CVViewID:    "cv-1",
				SectionType: "experience",
				Title:       "Experience",
				Order:       1,
				Content: []*career.SectionContentGroup{
					{
						Header:  "Acme Corp",
						Bullets: []*career.CVBullet{bullet},
					},
				},
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
			cv = &career.CVView{
				ID:               "cv-1",
				Name:             "Senior Software Engineer CV",
				TargetRole:       "senior_ic",
				TargetAudience:   "hiring_manager",
				GeneratedAt:      time.Now(),
				SourceEventCount: 10,
				SourceFactCount:  5,
			}

			bullet1 := &career.CVBullet{
				ID:              "bullet-1",
				SectionID:       "section-exp",
				Text:            "Led migration to microservices architecture",
				SourceEventIDs:  []string{"event-1"},
				Rank:            0.9,
				InclusionReason: "ownership",
				Confidence:      0.85,
			}

			bullet2 := &career.CVBullet{
				ID:              "bullet-2",
				SectionID:       "section-exp",
				Text:            "Implemented CI/CD pipeline",
				SourceEventIDs:  []string{"event-2"},
				Rank:            0.7,
				InclusionReason: "execution",
				Confidence:      0.70, // Below narrative threshold
			}

			sections = []*career.CVSection{
				{
					ID:          "section-summary",
					CVViewID:    "cv-1",
					SectionType: "summary",
					Title:       "Summary",
					Order:       0,
					Summary:     "Experienced software engineer with 10+ years in backend development.",
				},
				{
					ID:          "section-exp",
					CVViewID:    "cv-1",
					SectionType: "experience",
					Title:       "Experience",
					Order:       1,
					Content: []*career.SectionContentGroup{
						{
							Header:    "TechCorp",
							StartDate: "Jan 2020",
							EndDate:   "Present",
							Bullets:   []*career.CVBullet{bullet1, bullet2},
						},
					},
				},
			}

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
				// Narrative includes profile header (uppercase in text format)
				gomega.Expect(content).To(gomega.ContainSubstring("YOMI COLLEDGE"))
				gomega.Expect(content).To(gomega.ContainSubstring("Senior Software Engineer"))
				// Narrative sections
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
				// Narrative includes profile header
				gomega.Expect(content).To(gomega.ContainSubstring("# Yomi Colledge"))
				gomega.Expect(content).To(gomega.ContainSubstring("Senior Software Engineer"))
				// Narrative sections
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
})
