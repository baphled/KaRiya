package cv_test

import (
	"context"
	"io"

	"github.com/baphled/kariya/internal/config"
	"github.com/baphled/kariya/internal/domain/career"
	"github.com/baphled/kariya/internal/logger"
	"github.com/baphled/kariya/internal/service/career/cv"
	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
)

var _ = Describe("CV Structure Export", func() {
	var (
		service *cv.ExportService
		log     *logger.Logger
		ctx     context.Context
		cvView  *career.CVView
	)

	BeforeEach(func() {
		log = logger.New(io.Discard, logger.InfoLevel)
		service = cv.NewExportService(log)
		ctx = context.Background()
		cvView = &career.CVView{
			ID:               "cv-1",
			Name:             "Test User",
			TargetRole:       "principal",
			TargetAudience:   "hiring_manager",
			SourceEventCount: 10,
			SourceFactCount:  25,
		}
	})

	Describe("Consulting Structure", func() {
		var (
			sections []*career.CVSection
			bullets  map[string][]*career.CVBullet
		)

		BeforeEach(func() {
			sections = []*career.CVSection{
				{
					ID:          "summary",
					SectionType: "summary",
					Title:       "Summary",
					Summary:     "Experienced consulting engineer with expertise in system modernization.",
				},
				{
					ID:          "exp-1",
					SectionType: "experience",
					Title:       "Client Engagements",
					Content: []*career.SectionContentGroup{
						{
							Header:    "Acme Corp",
							StartDate: "Jan 2022",
							EndDate:   "Jan 2024",
						},
					},
				},
				{
					ID:          "exp-2",
					SectionType: "experience",
					Title:       "Client Engagements",
					Content: []*career.SectionContentGroup{
						{
							Header:    "TechStart Inc",
							StartDate: "Jan 2020",
							EndDate:   "Jan 2022",
						},
					},
				},
			}

			bullets = map[string][]*career.CVBullet{
				"exp-1": {
					{ID: "b1", Text: "Led system modernization project reducing technical debt by 40%", Confidence: 0.85},
					{ID: "b2", Text: "Delivered architecture roadmap adopted by client", Confidence: 0.80},
				},
				"exp-2": {
					{ID: "b3", Text: "Conducted rapid technology assessment for startup", Confidence: 0.75},
					{ID: "b4", Text: "Implemented CI/CD pipeline reducing deployment time", Confidence: 0.70},
				},
			}
		})

		Context("when exporting to text format", func() {
			It("should include profile header", func() {
				content, err := service.Export(ctx, cvView, sections, bullets, cv.CVStructureConsulting, cv.ExportFormatText)
				Expect(err).NotTo(HaveOccurred())
				Expect(content).To(ContainSubstring("TEST USER"))
			})

			It("should include summary section", func() {
				content, err := service.Export(ctx, cvView, sections, bullets, cv.CVStructureConsulting, cv.ExportFormatText)
				Expect(err).NotTo(HaveOccurred())
				Expect(content).To(ContainSubstring("SUMMARY"))
				Expect(content).To(ContainSubstring("consulting engineer"))
			})

			It("should include Client Engagements section", func() {
				content, err := service.Export(ctx, cvView, sections, bullets, cv.CVStructureConsulting, cv.ExportFormatText)
				Expect(err).NotTo(HaveOccurred())
				Expect(content).To(ContainSubstring("CLIENT ENGAGEMENTS"))
			})

			It("should group experience by company", func() {
				content, err := service.Export(ctx, cvView, sections, bullets, cv.CVStructureConsulting, cv.ExportFormatText)
				Expect(err).NotTo(HaveOccurred())
				Expect(content).To(ContainSubstring("Acme Corp"))
				Expect(content).To(ContainSubstring("TechStart Inc"))
			})

			It("should include bullets under each company", func() {
				content, err := service.Export(ctx, cvView, sections, bullets, cv.CVStructureConsulting, cv.ExportFormatText)
				Expect(err).NotTo(HaveOccurred())
				Expect(content).To(ContainSubstring("system modernization"))
				Expect(content).To(ContainSubstring("technology assessment"))
			})

			It("should include What I Bring section with profile", func() {
				profileCfg := &config.ProfileConfig{
					WhatIBring: []string{
						"Rapid technology assessment",
						"Clear technical communication",
					},
				}
				content, err := service.ExportWithProfile(ctx, cvView, sections, bullets, cv.CVStructureConsulting, cv.ExportFormatText, profileCfg)
				Expect(err).NotTo(HaveOccurred())
				Expect(content).To(ContainSubstring("WHAT I BRING"))
				Expect(content).To(ContainSubstring("Rapid technology assessment"))
			})
		})

		Context("when exporting to markdown format", func() {
			It("should use markdown headers", func() {
				content, err := service.Export(ctx, cvView, sections, bullets, cv.CVStructureConsulting, cv.ExportFormatMarkdown)
				Expect(err).NotTo(HaveOccurred())
				Expect(content).To(ContainSubstring("# "))
				Expect(content).To(ContainSubstring("## "))
			})

			It("should format bullets as markdown list", func() {
				content, err := service.Export(ctx, cvView, sections, bullets, cv.CVStructureConsulting, cv.ExportFormatMarkdown)
				Expect(err).NotTo(HaveOccurred())
				Expect(content).To(ContainSubstring("- Led system modernization"))
			})

			It("should include Client Engagements header", func() {
				content, err := service.Export(ctx, cvView, sections, bullets, cv.CVStructureConsulting, cv.ExportFormatMarkdown)
				Expect(err).NotTo(HaveOccurred())
				Expect(content).To(ContainSubstring("## Client Engagements"))
			})
		})
	})

	Describe("Highlights Structure", func() {
		var (
			sections []*career.CVSection
			bullets  map[string][]*career.CVBullet
		)

		BeforeEach(func() {
			sections = []*career.CVSection{
				{
					ID:          "summary",
					SectionType: "summary",
					Title:       "Summary",
					Summary:     "Senior engineer specializing in distributed systems.",
				},
				{
					ID:          "exp-1",
					SectionType: "experience",
					Title:       "Experience",
					Content: []*career.SectionContentGroup{
						{
							Header:    "BigTech Co",
							StartDate: "Jan 2020",
							EndDate:   "Present",
						},
					},
				},
			}

			// Create bullets with varying confidence for highlights selection
			bullets = map[string][]*career.CVBullet{
				"exp-1": {
					{ID: "b1", Text: "Designed distributed caching system serving 1M requests/sec", Confidence: 0.95},
					{ID: "b2", Text: "Led team of 8 engineers on platform modernization", Confidence: 0.90},
					{ID: "b3", Text: "Reduced infrastructure costs by 35%", Confidence: 0.88},
					{ID: "b4", Text: "Implemented event-driven architecture", Confidence: 0.85},
					{ID: "b5", Text: "Mentored 5 junior engineers", Confidence: 0.82},
					{ID: "b6", Text: "Improved API response time by 60%", Confidence: 0.80},
					{ID: "b7", Text: "Lower priority achievement", Confidence: 0.65},
					{ID: "b8", Text: "Another lower priority item", Confidence: 0.60},
				},
			}
		})

		Context("when exporting to text format", func() {
			It("should include condensed profile header", func() {
				content, err := service.Export(ctx, cvView, sections, bullets, cv.CVStructureHighlights, cv.ExportFormatText)
				Expect(err).NotTo(HaveOccurred())
				Expect(content).To(ContainSubstring("TEST USER"))
			})

			It("should include short summary", func() {
				content, err := service.Export(ctx, cvView, sections, bullets, cv.CVStructureHighlights, cv.ExportFormatText)
				Expect(err).NotTo(HaveOccurred())
				Expect(content).To(ContainSubstring("distributed systems"))
			})

			It("should include Key Capabilities section", func() {
				content, err := service.Export(ctx, cvView, sections, bullets, cv.CVStructureHighlights, cv.ExportFormatText)
				Expect(err).NotTo(HaveOccurred())
				Expect(content).To(ContainSubstring("KEY CAPABILITIES"))
			})

			It("should include Selected Highlights section", func() {
				content, err := service.Export(ctx, cvView, sections, bullets, cv.CVStructureHighlights, cv.ExportFormatText)
				Expect(err).NotTo(HaveOccurred())
				Expect(content).To(ContainSubstring("SELECTED HIGHLIGHTS"))
			})

			It("should limit highlights to top 5 by confidence", func() {
				content, err := service.Export(ctx, cvView, sections, bullets, cv.CVStructureHighlights, cv.ExportFormatText)
				Expect(err).NotTo(HaveOccurred())
				// Should include top 5 confidence bullets
				Expect(content).To(ContainSubstring("distributed caching"))  // 0.95
				Expect(content).To(ContainSubstring("team of 8"))            // 0.90
				Expect(content).To(ContainSubstring("infrastructure costs")) // 0.88
				Expect(content).To(ContainSubstring("event-driven"))         // 0.85
				Expect(content).To(ContainSubstring("Mentored"))             // 0.82
				// Should NOT include lower confidence bullets
				Expect(content).NotTo(ContainSubstring("Lower priority achievement"))
			})

			It("should include Technologies section", func() {
				profileCfg := &config.ProfileConfig{
					Languages: []string{"Go", "Ruby", "Python"},
					Systems:   []string{"Linux", "Kubernetes", "PostgreSQL"},
				}
				content, err := service.ExportWithProfile(ctx, cvView, sections, bullets, cv.CVStructureHighlights, cv.ExportFormatText, profileCfg)
				Expect(err).NotTo(HaveOccurred())
				Expect(content).To(ContainSubstring("TECHNOLOGIES"))
				Expect(content).To(ContainSubstring("Go, Ruby"))
			})
		})

		Context("when exporting to markdown format", func() {
			It("should use markdown headers", func() {
				content, err := service.Export(ctx, cvView, sections, bullets, cv.CVStructureHighlights, cv.ExportFormatMarkdown)
				Expect(err).NotTo(HaveOccurred())
				Expect(content).To(ContainSubstring("# "))
			})

			It("should include Key Capabilities as markdown section", func() {
				content, err := service.Export(ctx, cvView, sections, bullets, cv.CVStructureHighlights, cv.ExportFormatMarkdown)
				Expect(err).NotTo(HaveOccurred())
				Expect(content).To(ContainSubstring("## Key Capabilities"))
			})

			It("should include Selected Highlights as markdown section", func() {
				content, err := service.Export(ctx, cvView, sections, bullets, cv.CVStructureHighlights, cv.ExportFormatMarkdown)
				Expect(err).NotTo(HaveOccurred())
				Expect(content).To(ContainSubstring("## Selected Highlights"))
			})

			It("should format highlights as markdown list", func() {
				content, err := service.Export(ctx, cvView, sections, bullets, cv.CVStructureHighlights, cv.ExportFormatMarkdown)
				Expect(err).NotTo(HaveOccurred())
				Expect(content).To(ContainSubstring("- "))
			})
		})

		Context("when using core strengths for Key Capabilities", func() {
			It("should use core_strengths from profile if available", func() {
				profileCfg := &config.ProfileConfig{
					CoreStrengths: []string{
						"Distributed systems design",
						"Technical leadership",
						"Performance optimization",
						"Team mentorship",
					},
				}
				content, err := service.ExportWithProfile(ctx, cvView, sections, bullets, cv.CVStructureHighlights, cv.ExportFormatText, profileCfg)
				Expect(err).NotTo(HaveOccurred())
				Expect(content).To(ContainSubstring("Distributed systems design"))
				Expect(content).To(ContainSubstring("Technical leadership"))
			})
		})
	})

	Describe("Structure Routing", func() {
		var (
			sections []*career.CVSection
			bullets  map[string][]*career.CVBullet
		)

		BeforeEach(func() {
			sections = []*career.CVSection{
				{ID: "summary", SectionType: "summary", Summary: "Test summary"},
			}
			bullets = map[string][]*career.CVBullet{}
		})

		It("should route consulting structure to consulting export", func() {
			content, err := service.Export(ctx, cvView, sections, bullets, cv.CVStructureConsulting, cv.ExportFormatText)
			Expect(err).NotTo(HaveOccurred())
			Expect(content).To(ContainSubstring("CLIENT ENGAGEMENTS"))
		})

		It("should route highlights structure to highlights export", func() {
			content, err := service.Export(ctx, cvView, sections, bullets, cv.CVStructureHighlights, cv.ExportFormatText)
			Expect(err).NotTo(HaveOccurred())
			Expect(content).To(ContainSubstring("KEY CAPABILITIES"))
		})

		It("should route standard structure to standard export", func() {
			content, err := service.Export(ctx, cvView, sections, bullets, cv.CVStructureStandard, cv.ExportFormatText)
			Expect(err).NotTo(HaveOccurred())
			// Standard format uses different section titles
			Expect(content).NotTo(ContainSubstring("CLIENT ENGAGEMENTS"))
			Expect(content).NotTo(ContainSubstring("KEY CAPABILITIES"))
		})

		It("should route narrative structure to narrative export", func() {
			content, err := service.Export(ctx, cvView, sections, bullets, cv.CVStructureNarrative, cv.ExportFormatText)
			Expect(err).NotTo(HaveOccurred())
			Expect(content).To(ContainSubstring("CORE STRENGTHS"))
		})
	})
})
