package cv_test

import (
	"context"
	"io"

	"github.com/baphled/kariya/internal/config"
	"github.com/baphled/kariya/internal/domain/career"
	"github.com/baphled/kariya/internal/logger"
	"github.com/baphled/kariya/internal/service/career/cv"
	"github.com/baphled/kariya/internal/testutil/fixtures"
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
		cvView = fixtures.CVViewWith("cv-1", "Test User", "principal", "hiring_manager")
		cvView.SourceFactCount = 25
	})

	Describe("Consulting Structure", func() {
		var (
			sections []*career.CVSection
			bullets  map[string][]*career.CVBullet
		)

		BeforeEach(func() {
			exp1Section := fixtures.CVSectionWith("exp-1", "", "experience", "Client Engagements", 0)
			exp1Section.Content = []*career.SectionContentGroup{
				fixtures.ContentGroupWith("Acme Corp", "Jan 2022", "Jan 2024"),
			}

			exp2Section := fixtures.CVSectionWith("exp-2", "", "experience", "Client Engagements", 0)
			exp2Section.Content = []*career.SectionContentGroup{
				fixtures.ContentGroupWith("TechStart Inc", "Jan 2020", "Jan 2022"),
			}

			sections = []*career.CVSection{
				fixtures.CVSectionWithSummary("summary", "", "Experienced consulting engineer with expertise in system modernization."),
				exp1Section,
				exp2Section,
			}

			b1 := fixtures.CVBulletWith("b1", "exp-1", "Led system modernization project reducing technical debt by 40%")
			b1.Confidence = 0.85

			b2 := fixtures.CVBulletWith("b2", "exp-1", "Delivered architecture roadmap adopted by client")
			b2.Confidence = 0.80

			b3 := fixtures.CVBulletWith("b3", "exp-2", "Conducted rapid technology assessment for startup")
			b3.Confidence = 0.75

			b4 := fixtures.CVBulletWith("b4", "exp-2", "Implemented CI/CD pipeline reducing deployment time")
			b4.Confidence = 0.70

			bullets = map[string][]*career.CVBullet{
				"exp-1": {b1, b2},
				"exp-2": {b3, b4},
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
			expSection := fixtures.CVSectionWith("exp-1", "", "experience", "Experience", 0)
			expSection.Content = []*career.SectionContentGroup{
				fixtures.ContentGroupWith("BigTech Co", "Jan 2020", "Present"),
			}

			sections = []*career.CVSection{
				fixtures.CVSectionWithSummary("summary", "", "Senior engineer specializing in distributed systems."),
				expSection,
			}

			b1 := fixtures.CVBulletWith("b1", "exp-1", "Designed distributed caching system serving 1M requests/sec")
			b1.Confidence = 0.95

			b2 := fixtures.CVBulletWith("b2", "exp-1", "Led team of 8 engineers on platform modernization")
			b2.Confidence = 0.90

			b3 := fixtures.CVBulletWith("b3", "exp-1", "Reduced infrastructure costs by 35%")
			b3.Confidence = 0.88

			b4 := fixtures.CVBulletWith("b4", "exp-1", "Implemented event-driven architecture")
			b4.Confidence = 0.85

			b5 := fixtures.CVBulletWith("b5", "exp-1", "Mentored 5 junior engineers")
			b5.Confidence = 0.82

			b6 := fixtures.CVBulletWith("b6", "exp-1", "Improved API response time by 60%")
			b6.Confidence = 0.80

			b7 := fixtures.CVBulletWith("b7", "exp-1", "Optimized database query performance")
			b7.Confidence = 0.78

			b8 := fixtures.CVBulletWith("b8", "exp-1", "Implemented comprehensive monitoring system")
			b8.Confidence = 0.76

			b9 := fixtures.CVBulletWith("b9", "exp-1", "Established security best practices")
			b9.Confidence = 0.74

			b10 := fixtures.CVBulletWith("b10", "exp-1", "Automated deployment pipeline")
			b10.Confidence = 0.72

			b11 := fixtures.CVBulletWith("b11", "exp-1", "Lower priority achievement")
			b11.Confidence = 0.65

			b12 := fixtures.CVBulletWith("b12", "exp-1", "Another lower priority item")
			b12.Confidence = 0.60

			bullets = map[string][]*career.CVBullet{
				"exp-1": {b1, b2, b3, b4, b5, b6, b7, b8, b9, b10, b11, b12},
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

			It("should limit highlights to top 10 by confidence", func() {
				content, err := service.Export(ctx, cvView, sections, bullets, cv.CVStructureHighlights, cv.ExportFormatText)
				Expect(err).NotTo(HaveOccurred())
				// Should include top 10 confidence bullets
				Expect(content).To(ContainSubstring("distributed caching"))  // 0.95
				Expect(content).To(ContainSubstring("team of 8"))            // 0.90
				Expect(content).To(ContainSubstring("infrastructure costs")) // 0.88
				Expect(content).To(ContainSubstring("event-driven"))         // 0.85
				Expect(content).To(ContainSubstring("Mentored"))             // 0.82
				Expect(content).To(ContainSubstring("API response time"))    // 0.80
				Expect(content).To(ContainSubstring("database query"))       // 0.78
				Expect(content).To(ContainSubstring("monitoring system"))    // 0.76
				Expect(content).To(ContainSubstring("security best"))        // 0.74
				Expect(content).To(ContainSubstring("deployment pipeline"))  // 0.72
				// Should NOT include lower confidence bullets (ranked 11+)
				Expect(content).NotTo(ContainSubstring("Lower priority achievement"))
				Expect(content).NotTo(ContainSubstring("Another lower priority item"))
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
				fixtures.CVSectionWithSummary("summary", "", "Test summary"),
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
