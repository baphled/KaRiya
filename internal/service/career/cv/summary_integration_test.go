package cv_test

import (
	"context"
	"io"
	"strings"

	"github.com/baphled/kariya/internal/ui/display"
	tea "github.com/charmbracelet/bubbletea"
	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"

	"github.com/baphled/kariya/internal/config"
	"github.com/baphled/kariya/internal/domain/career"
	"github.com/baphled/kariya/internal/logger"
	cvsvc "github.com/baphled/kariya/internal/service/career/cv"
	"github.com/baphled/kariya/internal/testutil/fixtures"
	cvview "github.com/baphled/kariya/internal/tui/views/cv"
)

var _ = Describe("Summary Format Integration", func() {
	var (
		ctx              context.Context
		sectionBuilder   *cvsvc.DefaultSectionBuilder
		exportService    *cvsvc.ExportService
		log              *logger.Logger
		realisticBullets []*career.CVBullet
		events           []*career.Event
	)

	BeforeEach(func() {
		ctx = context.Background()
		log = logger.New(io.Discard, logger.InfoLevel)
		sectionBuilder = cvsvc.NewSectionBuilder(nil, log)
		exportService = cvsvc.NewExportService(log, nil, nil)

		realisticBullets = []*career.CVBullet{
			createBullet("b1", "Migrated QuikCV backend from Ruby on Rails to Node.js, reducing server costs by 70% and improving response times by 40% through async processing"),
			createBullet("b2", "Founded n-vyro.io IoT platform, delivering production-ready firmware in C/C++ and backend services in Go and Node.js for real-time device control"),
			createBullet("b3", "Adopted Jest early for QuikCV testing, establishing a test-first culture that reduced regression bugs by 60% across the engineering team"),
			createBullet("b4", "Built AI-powered customer service dashboards at Digital Genius integrating ML APIs for intelligent response routing and ticket classification"),
			createBullet("b5", "Led infrastructure reliability initiative at ITV ensuring 99.9% uptime for live broadcast systems serving 10M+ concurrent viewers"),
		}

		events = []*career.Event{
			fixtures.EventWith("e1", "Backend migration", "QuikCV", "Migration"),
			fixtures.EventWith("e2", "IoT platform", "n-vyro.io", "Platform"),
		}
	})

	Describe("BuildSections summary format", func() {
		It("produces prose with exactly 2 bullet-sentences joined by '. '", func() {
			sections, err := sectionBuilder.BuildSections(ctx, realisticBullets, events, []*career.Fact{}, "senior_ic", nil, nil)
			Expect(err).NotTo(HaveOccurred())

			var summarySection *career.CVSection
			for _, section := range sections {
				if section.SectionType == "summary" {
					summarySection = section
					break
				}
			}
			Expect(summarySection).NotTo(BeNil(), "expected a summary section to be created")
			Expect(summarySection.Summary).NotTo(BeEmpty())

			summaryText := summarySection.Summary
			Expect(summaryText).To(ContainSubstring("Migrated QuikCV backend from Ruby on Rails"))
			Expect(summaryText).To(ContainSubstring("Founded n-vyro.io IoT platform"))
			Expect(summaryText).NotTo(ContainSubstring("Adopted Jest early"), "3rd bullet should not appear")

			Expect(summaryText).To(ContainSubstring(". "), "sentences should be separated by '. '")

			dotSpaceSeparatorCount := strings.Count(summaryText, ". ")
			Expect(dotSpaceSeparatorCount).To(Equal(1), "expect exactly 1 '. ' separator for 2 sentences")

			Expect(summaryText).To(HaveSuffix("."), "prose should end with a period")
		})

		It("produces short paragraph, not a wall of text", func() {
			sections, err := sectionBuilder.BuildSections(ctx, realisticBullets, events, []*career.Fact{}, "senior_ic", nil, nil)
			Expect(err).NotTo(HaveOccurred())

			var summarySection *career.CVSection
			for _, section := range sections {
				if section.SectionType == "summary" {
					summarySection = section
					break
				}
			}
			Expect(summarySection).NotTo(BeNil())

			Expect(len(summarySection.Summary)).To(BeNumerically("<", 600), "summary should be under 600 chars (short paragraph)")
		})
	})

	Describe("ExportToYAML summary field", func() {
		It("contains correctly formatted short paragraph in summary field", func() {
			sections, err := sectionBuilder.BuildSections(ctx, realisticBullets, events, []*career.Fact{}, "senior_ic", nil, nil)
			Expect(err).NotTo(HaveOccurred())

			cvView := fixtures.CVViewWith("cv-integration", "Integration Test CV", "senior_ic", "hiring_manager")
			cvView.Sections = sections

			profileCfg := &config.ProfileConfig{
				FirstName: "Test",
				LastName:  "User",
				Email:     "test@example.com",
			}

			yamlOutput, err := exportService.ExportToYAML(ctx, cvView, sections, profileCfg)
			Expect(err).NotTo(HaveOccurred())

			Expect(yamlOutput).To(ContainSubstring("summary:"))

			summaryLineIdx := strings.Index(yamlOutput, "summary:")
			Expect(summaryLineIdx).To(BeNumerically(">", -1))

			summaryContent := extractYAMLSummary(yamlOutput)
			Expect(summaryContent).To(ContainSubstring("Migrated QuikCV backend from Ruby on Rails"))
			Expect(summaryContent).To(ContainSubstring("Founded n-vyro.io IoT platform"))
			Expect(summaryContent).NotTo(ContainSubstring("Adopted Jest early for QuikCV testing"))
			Expect(summaryContent).NotTo(ContainSubstring("Built AI-powered customer service dashboards"))
			Expect(summaryContent).NotTo(ContainSubstring("Led infrastructure reliability initiative"))

			Expect(summaryContent).To(ContainSubstring(". "), "YAML summary should contain '. ' sentence separators")
		})
	})

	Describe("Review screen summary display", func() {
		It("renders summary with correct format in View()", func() {
			sections, err := sectionBuilder.BuildSections(ctx, realisticBullets, events, []*career.Fact{}, "senior_ic", nil, nil)
			Expect(err).NotTo(HaveOccurred())

			cvView := fixtures.CVViewWith("cv-review", "Review Test CV", "senior_ic", "hiring_manager")
			cvView.Sections = sections

			profileCfg := &config.ProfileConfig{
				Name:     "Test User",
				Email:    "test@example.com",
				Location: "London, UK",
			}

			summary := &cvview.GenerationSummary{
				SelectedAudience: "Hiring Manager",
				TechnologyFocus:  "Backend",
				SourceEventCount: 5,
				SourceFactCount:  10,
				SectionCount:     len(sections),
				TotalBullets:     countBulletsInSections(sections),
			}

			screen := cvview.NewReview(display.CVViewFromDomain(cvView), profileCfg, summary)
			screen.Update(tea.WindowSizeMsg{Width: 120, Height: 60})

			view := screen.RenderContent()

			Expect(view).To(ContainSubstring("Summary"))
			Expect(view).To(ContainSubstring("Migrated QuikCV backend"))
			Expect(view).To(ContainSubstring("Founded n-vyro.io"))
			Expect(view).NotTo(ContainSubstring("Adopted Jest early"))
			Expect(view).NotTo(ContainSubstring("Built AI-powered customer service dashboards"))
			Expect(view).NotTo(ContainSubstring("Led infrastructure reliability"))
		})
	})

	Describe("Preview screen summary display", func() {
		It("renders summary with word-wrapping, not a wall of text", func() {
			sections, err := sectionBuilder.BuildSections(ctx, realisticBullets, events, []*career.Fact{}, "senior_ic", nil, nil)
			Expect(err).NotTo(HaveOccurred())

			cvView := fixtures.CVViewWith("cv-preview", "Preview Test CV", "senior_ic", "hiring_manager")
			cvView.Sections = sections

			profileCfg := &config.ProfileConfig{
				Name:     "Test User",
				Email:    "test@example.com",
				Location: "London, UK",
				Title:    "Senior Engineer",
			}

			screen := cvview.NewPreview(display.CVViewFromDomain(cvView), profileCfg)
			screen.Update(tea.WindowSizeMsg{Width: 80, Height: 60})

			view := screen.RenderContent()

			Expect(view).To(ContainSubstring("Professional Summary"))

			// Extract summary section content (between Summary header and next section)
			summaryStart := strings.Index(view, "Professional Summary")
			Expect(summaryStart).To(BeNumerically(">", -1), "should find Professional Summary section")

			experienceStart := strings.Index(view, "Experience")
			var summarySection string
			if experienceStart > summaryStart {
				summarySection = view[summaryStart:experienceStart]
			} else {
				summarySection = view[summaryStart:]
			}

			// Summary section should contain top 2 bullets
			Expect(summarySection).To(ContainSubstring("Migrated QuikCV backend"))
			Expect(summarySection).To(ContainSubstring("Founded n-vyro.io"))
			Expect(summarySection).NotTo(ContainSubstring("Adopted Jest early"))
			Expect(summarySection).NotTo(ContainSubstring("Built AI-powered customer service dashboards"))
			// 3rd+ bullets should NOT appear in summary section
			Expect(summarySection).NotTo(ContainSubstring("Led infrastructure reliability"))
		})
	})
})

func createBullet(id, text string) *career.CVBullet {
	return fixtures.CVBulletWith(id, "section-1", text)
}

func extractYAMLSummary(yamlContent string) string {
	lines := strings.Split(yamlContent, "\n")
	for i, line := range lines {
		if !strings.HasPrefix(line, "summary:") {
			continue
		}
		summaryValue := strings.TrimPrefix(line, "summary:")
		summaryValue = strings.TrimSpace(summaryValue)

		if summaryValue != "" && !strings.HasPrefix(summaryValue, "|") {
			return summaryValue
		}

		if summaryValue == "" || strings.HasPrefix(summaryValue, "|") {
			return extractMultilineSummary(lines, i)
		}
	}
	return ""
}

func extractMultilineSummary(lines []string, startIdx int) string {
	var sb strings.Builder
	for j := startIdx + 1; j < len(lines); j++ {
		nextLine := lines[j]
		if nextLine != "" && nextLine[0] != ' ' && nextLine[0] != '\t' {
			break
		}
		sb.WriteString(strings.TrimSpace(nextLine))
		sb.WriteString(" ")
	}
	return strings.TrimSpace(sb.String())
}

func countBulletsInSections(sections []*career.CVSection) int {
	count := 0
	for _, section := range sections {
		for _, group := range section.Content {
			count += len(group.Bullets)
		}
	}
	return count
}
