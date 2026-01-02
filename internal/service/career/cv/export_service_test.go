package cv

import (
	"context"
	"io"
	"os"
	"path/filepath"
	"strings"
	"time"

	"github.com/atotto/clipboard"
	"github.com/baphled/kariya/internal/domain/career"
	"github.com/baphled/kariya/internal/logger"
	"github.com/onsi/ginkgo/v2"
	"github.com/onsi/gomega"
)


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
				ID:              "cv-1",
				Name:            "Senior Software Engineer CV",
				TargetRole:      "Staff Engineer",
				TargetAudience:  []string{"hiring_manager", "recruiter"},
				GeneratedAt:     time.Now(),
				SourceEventCount: 10,
				SourceFactCount:  5,
			}

			section := &career.CVSection{
				ID:          "section-1",
				CVViewID:    "cv-1",
				SectionType: "experience",
				Title:       "Experience",
				Order:       1,
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
				ID:              "cv-1",
				Name:            "Test CV",
				TargetRole:      "Engineer",
				TargetAudience:  []string{"hiring_manager"},
				GeneratedAt:     time.Now(),
				SourceEventCount: 0,
				SourceFactCount:  0,
			}

			section := &career.CVSection{
				ID:          "section-1",
				CVViewID:    "cv-1",
				SectionType: "experience",
				Title:       "Experience",
				Order:       1,
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
				ID:              "cv-1",
				Name:            "Senior Software Engineer CV",
				TargetRole:      "Staff Engineer",
				TargetAudience:  []string{"hiring_manager"},
				GeneratedAt:     time.Now(),
				SourceEventCount: 10,
				SourceFactCount:  5,
			}

			section := &career.CVSection{
				ID:          "section-1",
				CVViewID:    "cv-1",
				SectionType: "experience",
				Title:       "Experience",
				Order:       1,
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
				ID:              "cv-1",
				Name:            "Test CV",
				TargetRole:      "Engineer",
				TargetAudience:  []string{"recruiter"},
				GeneratedAt:     time.Now(),
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
				ID:              "cv-1",
				Name:            "Senior Software Engineer CV",
				TargetRole:      "Staff Engineer",
				TargetAudience:  []string{"hiring_manager"},
				GeneratedAt:     time.Now(),
				SourceEventCount: 10,
				SourceFactCount:  5,
			}

			section := &career.CVSection{
				ID:          "section-1",
				CVViewID:    "cv-1",
				SectionType: "experience",
				Title:       "Experience",
				Order:       1,
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

			sections := []*career.CVSection{section}
			bullets := map[string][]*career.CVBullet{
				"section-1": {bullet},
			}

			yaml, err := service.ExportToYAML(ctx, cv, sections, bullets)
			gomega.Expect(err).NotTo(gomega.HaveOccurred())
			gomega.Expect(yaml).NotTo(gomega.BeEmpty())
			gomega.Expect(yaml).To(gomega.ContainSubstring("name: Senior Software Engineer CV"))
			gomega.Expect(yaml).To(gomega.ContainSubstring("target_role: Staff Engineer"))
			gomega.Expect(yaml).To(gomega.ContainSubstring("- Led team of 5 engineers"))
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
		ginkgo.It("should copy content to clipboard", func() {
			testContent := "Test CV Content for Clipboard"
			err := service.CopyToClipboard(ctx, testContent)
			gomega.Expect(err).NotTo(gomega.HaveOccurred())

			// Verify content was copied
			clipboardContent, err := clipboard.ReadAll()
			gomega.Expect(err).NotTo(gomega.HaveOccurred())
			gomega.Expect(clipboardContent).To(gomega.Equal(testContent))
		})

		ginkgo.It("should return error for empty content", func() {
			err := service.CopyToClipboard(ctx, "")
			gomega.Expect(err).To(gomega.HaveOccurred())
			gomega.Expect(err.Error()).To(gomega.ContainSubstring("content is empty"))
		})
	})

	})
})

