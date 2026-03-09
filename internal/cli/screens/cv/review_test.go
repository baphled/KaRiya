package cv_test

import (
	"strings"
	"time"

	tea "github.com/charmbracelet/bubbletea"

	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"

	"github.com/baphled/kariya/internal/cli/screens"
	"github.com/baphled/kariya/internal/cli/screens/cv"
	"github.com/baphled/kariya/internal/cli/themes"
	"github.com/baphled/kariya/internal/cli/types"
	"github.com/baphled/kariya/internal/config"
	"github.com/baphled/kariya/internal/domain/career"
	"github.com/baphled/kariya/internal/testutil/fixtures"
)

var _ = Describe("ReviewScreen", func() {
	var (
		sampleCV        *career.CVView
		sampleProfile   *config.ProfileConfig
		sampleSummary   *cv.GenerationSummaryScreen
		sampleCVProfile *types.CVProfile
	)

	BeforeEach(func() {
		sampleCVProfile = &types.CVProfile{
			ID:             "test-profile-1",
			Name:           "Senior Backend Engineer",
			TargetRole:     "senior_ic",
			TargetAudience: "hiring_manager",
			Description:    "Test CV profile for backend engineering",
		}

		bullet1_1 := fixtures.CVBulletWith("b1-1", "section-1", "Led development of microservices architecture")
		bullet1_2 := fixtures.CVBulletWith("b1-2", "section-1", "Improved system performance by 40%")
		group1 := fixtures.ContentGroupWithBullets("Senior Developer at TechCorp", []*career.CVBullet{bullet1_1, bullet1_2})
		section1 := fixtures.CVSectionWithContent("section-1", "test-cv-123", []*career.SectionContentGroup{group1})
		section1.Title = "Professional Experience"
		section1.SectionType = "experience"

		bullet2_1 := fixtures.CVBulletWith("b2-1", "section-2", "Go, Python, JavaScript")
		bullet2_2 := fixtures.CVBulletWith("b2-2", "section-2", "SQL, NoSQL databases")
		group2 := fixtures.ContentGroupWithBullets("Programming Languages", []*career.CVBullet{bullet2_1, bullet2_2})
		section2 := fixtures.CVSectionWithContent("section-2", "test-cv-123", []*career.SectionContentGroup{group2})
		section2.Title = "Technical Skills"
		section2.SectionType = "skills"

		bullet3_1 := fixtures.CVBulletWith("b3-1", "section-3", "Experienced software engineer with 8+ years")
		group3 := fixtures.ContentGroupWithBullets("Professional Summary", []*career.CVBullet{bullet3_1})
		section3 := fixtures.CVSectionWithContent("section-3", "test-cv-123", []*career.SectionContentGroup{group3})
		section3.Title = "Summary"
		section3.SectionType = "summary"

		sampleCV = fixtures.CVViewWithSections("test-cv-123", []*career.CVSection{section1, section2, section3})
		sampleCV.Name = "John Doe CV"
		sampleCV.TargetRole = "Senior Software Engineer"
		sampleCV.TargetAudience = "Hiring Manager"
		sampleCV.EventFilters = map[string]interface{}{"company": "TechCorp"}
		sampleCV.GeneratedAt = time.Now()
		sampleCV.SourceEventCount = 8
		sampleCV.SourceFactCount = 12

		sampleProfile = &config.ProfileConfig{
			Name:            "John Doe",
			Email:           "john.doe@example.com",
			DefaultRole:     "senior_ic",
			DefaultAudience: "hiring_manager",
			Title:           "Senior Software Engineer",
			Location:        "San Francisco, CA",
			GitHub:          "github.com/johndoe",
			Portfolio:       "johndoe.dev",
			CoreStrengths:   []string{"Backend Development", "System Design"},
			Languages:       []string{"English", "Spanish"},
		}

		sampleSummary = &cv.GenerationSummaryScreen{
			SelectedProfile:  sampleCVProfile,
			SelectedAudience: "Hiring Manager",
			TechnologyFocus:  "Backend Development",
			Technologies:     []string{"Go", "Docker", "Kubernetes"},
			FocusArea:        "Microservices Architecture",
			SkillsFormat:     "grouped",
			SkillsLimit:      15,
			CVLength:         "comprehensive",
			SourceEventCount: 8,
			SourceFactCount:  12,
			SectionCount:     3,
			TotalBullets:     5,
		}
	})

	Describe("Constructor tests", func() {
		Context("NewCVReviewScreenWithSummary", func() {
			It("creates a screen with the summary field set", func() {
				screen := cv.NewCVReviewScreenWithSummary(sampleCV, sampleProfile, sampleSummary)
				screen.Update(tea.WindowSizeMsg{Width: 120, Height: 40})

				Expect(screen).NotTo(BeNil())
				Expect(screen.GetCV()).To(Equal(sampleCV))
				// We can verify the summary is set by checking the View() output contains summary-specific content
				view := screen.View()
				Expect(view).To(ContainSubstring("🎯 Generation Settings"))
			})

			It("works with nil summary", func() {
				screen := cv.NewCVReviewScreenWithSummary(sampleCV, sampleProfile, nil)
				screen.Update(tea.WindowSizeMsg{Width: 120, Height: 40})

				Expect(screen).NotTo(BeNil())
				Expect(screen.GetCV()).To(Equal(sampleCV))
				// Without summary, should not show summary-specific content
				view := screen.View()
				Expect(view).NotTo(ContainSubstring("🎯 Generation Settings"))
			})

			It("works with nil profile config", func() {
				screen := cv.NewCVReviewScreenWithSummary(sampleCV, nil, sampleSummary)
				screen.Update(tea.WindowSizeMsg{Width: 120, Height: 40})

				Expect(screen).NotTo(BeNil())
				Expect(screen.GetCV()).To(Equal(sampleCV))
				view := screen.View()
				Expect(view).To(ContainSubstring("🎯 Generation Settings"))
			})
		})
	})

	Describe("Init", func() {
		It("returns nil command", func() {
			screen := cv.NewCVReviewScreen(sampleCV)
			cmd := screen.Init()
			Expect(cmd).To(BeNil())
		})
	})

	Describe("Update key handling", func() {
		var screen *cv.ReviewScreen

		BeforeEach(func() {
			screen = cv.NewCVReviewScreenWithProfile(sampleCV, sampleProfile)
			screen.Update(tea.WindowSizeMsg{Width: 120, Height: 40})
			screen.View()
		})

		It("returns NavigateResult with 'preview' on enter key", func() {
			cmd, result := screen.Update(tea.KeyMsg{Type: tea.KeyEnter})
			Expect(cmd).To(BeNil())
			Expect(result).NotTo(BeNil())
			Expect(result.Type()).To(Equal(screens.ResultNavigate))
			Expect(result.Data()).To(Equal("preview"))
		})

		It("returns NavigateResult with 'preview' on 'p' key", func() {
			cmd, result := screen.Update(tea.KeyMsg{Type: tea.KeyRunes, Runes: []rune{'p'}})
			Expect(cmd).To(BeNil())
			Expect(result).NotTo(BeNil())
			Expect(result.Data()).To(Equal("preview"))
		})

		It("returns NavigateResult with 'export' on 'x' key", func() {
			cmd, result := screen.Update(tea.KeyMsg{Type: tea.KeyRunes, Runes: []rune{'x'}})
			Expect(cmd).To(BeNil())
			Expect(result).NotTo(BeNil())
			Expect(result.Data()).To(Equal("export"))
		})

		It("returns NavigateResult with 'edit' on 'e' key", func() {
			cmd, result := screen.Update(tea.KeyMsg{Type: tea.KeyRunes, Runes: []rune{'e'}})
			Expect(cmd).To(BeNil())
			Expect(result).NotTo(BeNil())
			Expect(result.Data()).To(Equal("edit"))
		})

		It("returns CancelResult on esc key", func() {
			cmd, result := screen.Update(tea.KeyMsg{Type: tea.KeyEsc})
			Expect(cmd).To(BeNil())
			Expect(result).NotTo(BeNil())
			Expect(result.Type()).To(Equal(screens.ResultCancel))
		})

		It("handles 'g' key to go to top", func() {
			cmd, result := screen.Update(tea.KeyMsg{Type: tea.KeyRunes, Runes: []rune{'g'}})
			Expect(cmd).To(BeNil())
			Expect(result).To(BeNil())
		})

		It("handles 'G' key to go to bottom", func() {
			cmd, result := screen.Update(tea.KeyMsg{Type: tea.KeyRunes, Runes: []rune{'G'}})
			Expect(cmd).To(BeNil())
			Expect(result).To(BeNil())
		})

		It("handles down key for scrolling", func() {
			cmd, result := screen.Update(tea.KeyMsg{Type: tea.KeyDown})
			Expect(cmd).To(BeNil())
			Expect(result).To(BeNil())
		})

		It("handles up key for scrolling", func() {
			screen.Update(tea.KeyMsg{Type: tea.KeyDown})
			cmd, result := screen.Update(tea.KeyMsg{Type: tea.KeyUp})
			Expect(cmd).To(BeNil())
			Expect(result).To(BeNil())
		})

		It("returns nil for scroll keys when viewport not ready", func() {
			freshScreen := cv.NewCVReviewScreenWithProfile(sampleCV, sampleProfile)
			cmd, result := freshScreen.Update(tea.KeyMsg{Type: tea.KeyDown})
			Expect(cmd).To(BeNil())
			Expect(result).To(BeNil())
		})

		It("returns nil for unhandled key", func() {
			cmd, result := screen.Update(tea.KeyMsg{Type: tea.KeyRunes, Runes: []rune{'z'}})
			Expect(cmd).To(BeNil())
			Expect(result).To(BeNil())
		})
	})

	Describe("View() rendering", func() {
		Context("when summary is non-nil", func() {
			It("contains Generation Settings section with profile/audience/tech/format fields", func() {
				screen := cv.NewCVReviewScreenWithSummary(sampleCV, sampleProfile, sampleSummary)
				screen.Update(tea.WindowSizeMsg{Width: 120, Height: 40})
				screen.Update(tea.WindowSizeMsg{Width: 120, Height: 40})
				view := screen.View()

				Expect(view).To(ContainSubstring("🎯 Generation Settings"))

				Expect(view).To(ContainSubstring("Senior Backend Engineer"))    // Profile name
				Expect(view).To(ContainSubstring("Hiring Manager"))             // Selected audience
				Expect(view).To(ContainSubstring("Backend Development"))        // TechnologyFocus
				Expect(view).To(ContainSubstring("Go, Docker, Kubernetes"))     // Technologies joined
				Expect(view).To(ContainSubstring("Microservices Architecture")) // FocusArea
				Expect(view).To(ContainSubstring("comprehensive"))              // CVLength
				Expect(view).To(ContainSubstring("grouped (15 max)"))           // SkillsFormat with limit
			})

			It("contains Statistics section from the generated CV", func() {
				screen := cv.NewCVReviewScreenWithSummary(sampleCV, sampleProfile, sampleSummary)
				screen.Update(tea.WindowSizeMsg{Width: 120, Height: 40})
				screen.Update(tea.WindowSizeMsg{Width: 120, Height: 40})
				view := screen.View()

				Expect(view).To(ContainSubstring("📊 Statistics"))
				Expect(view).To(ContainSubstring("Source Events:"))
				Expect(view).To(ContainSubstring("Source Facts:"))
				Expect(view).To(ContainSubstring("Sections:"))
				Expect(view).To(ContainSubstring("Total Bullets:"))
			})

			It("displays the CV summary section when present", func() {
				summarySection := fixtures.CVSectionWithSummary("section-summary", "test-cv-123", "**Senior Engineer** with expertise in Go and distributed systems.")
				cvWithSummary := fixtures.CVViewWithSections("test-cv-123", []*career.CVSection{summarySection})
				cvWithSummary.Name = "Test CV"
				cvWithSummary.TargetRole = "Engineer"
				cvWithSummary.TargetAudience = "Hiring Manager"

				screen := cv.NewCVReviewScreenWithSummary(cvWithSummary, sampleProfile, sampleSummary)
				screen.Update(tea.WindowSizeMsg{Width: 120, Height: 40})
				view := screen.View()

				Expect(view).To(ContainSubstring("📝 Summary"))
				Expect(view).To(ContainSubstring("Senior Engineer"))
				Expect(view).To(ContainSubstring("distributed"))
				Expect(view).To(ContainSubstring("systems"))
			})

			It("handles empty or default values gracefully by omitting them", func() {
				emptySummary := &cv.GenerationSummaryScreen{
					SelectedProfile:  nil, // nil profile
					SelectedAudience: "",  // empty audience
					TechnologyFocus:  "",  // empty focus
					Technologies:     nil, // empty technologies
					FocusArea:        "",  // empty area
					SkillsFormat:     "",  // empty skills format
					SkillsLimit:      0,   // zero limit
					CVLength:         "",  // empty length
					SourceEventCount: 0,
					SourceFactCount:  0,
					SectionCount:     0,
					TotalBullets:     0,
				}

				screen := cv.NewCVReviewScreenWithSummary(sampleCV, sampleProfile, emptySummary)
				screen.Update(tea.WindowSizeMsg{Width: 120, Height: 40})
				view := screen.View()

				// The Generation Settings section should still appear
				Expect(view).To(ContainSubstring("🎯 Generation Settings"))

				// Empty values should be omitted, not shown as dashes
				Expect(view).NotTo(ContainSubstring("Profile:"))
				Expect(view).NotTo(ContainSubstring("Technologies:"))
				Expect(view).NotTo(ContainSubstring("Skills Format:"))
				Expect(view).NotTo(ContainSubstring("CV Length:"))

				// Should still show Personal Details, CV Details, and Statistics
				Expect(view).To(ContainSubstring("👤 Personal Details"))
				Expect(view).To(ContainSubstring("📄 CV Details"))
				Expect(view).To(ContainSubstring("📊 Statistics"))
			})
		})

		Context("when summary is nil", func() {
			It("does NOT show YOUR SELECTIONS section - falls back to legacy layout", func() {
				screen := cv.NewCVReviewScreenWithSummary(sampleCV, sampleProfile, nil)
				screen.Update(tea.WindowSizeMsg{Width: 120, Height: 40})
				view := screen.View()

				// Should NOT contain summary-specific sections
				Expect(view).NotTo(ContainSubstring("YOUR SELECTIONS"))
				Expect(view).NotTo(ContainSubstring("GENERATED OUTPUT"))

				// Should contain legacy sections
				Expect(view).To(ContainSubstring("👤 Personal Details"))
				Expect(view).To(ContainSubstring("📄 CV Details"))
				Expect(view).To(ContainSubstring("📊 Statistics"))
				Expect(view).To(ContainSubstring("📑 Sections"))
			})

			It("shows legacy personal details from profile config", func() {
				screen := cv.NewCVReviewScreenWithSummary(sampleCV, sampleProfile, nil)
				screen.Update(tea.WindowSizeMsg{Width: 120, Height: 40})
				view := screen.View()

				Expect(view).To(ContainSubstring("John Doe"))             // Name from profile
				Expect(view).To(ContainSubstring("john.doe@example.com")) // Email from profile
				Expect(view).To(ContainSubstring("San Francisco, CA"))    // Location from profile
			})

			It("shows CV details and statistics from CVView", func() {
				screen := cv.NewCVReviewScreenWithSummary(sampleCV, sampleProfile, nil)
				screen.Update(tea.WindowSizeMsg{Width: 120, Height: 40})
				view := screen.View()

				// CV details
				Expect(view).To(ContainSubstring("John Doe CV"))              // CV name
				Expect(view).To(ContainSubstring("Senior Software Engineer")) // Target role
				Expect(view).To(ContainSubstring("Hiring Manager"))           // Target audience

				// Statistics
				Expect(view).To(ContainSubstring("Source Events: 8"))
				Expect(view).To(ContainSubstring("Source Facts: 12"))
				Expect(view).To(ContainSubstring("Sections: 3"))
				Expect(view).To(ContainSubstring("Total Bullets: 5")) // 2+2+1 bullets from sections
			})
		})
	})

	Describe("View edge cases", func() {
		It("shows no data message when CV is nil", func() {
			nilScreen := cv.NewCVReviewScreen(nil)
			nilScreen.Update(tea.WindowSizeMsg{Width: 120, Height: 40})
			view := nilScreen.View()
			Expect(view).To(ContainSubstring("No CV data available"))
			Expect(view).To(ContainSubstring("esc"))
		})

		It("handles small viewport dimensions", func() {
			screen := cv.NewCVReviewScreenWithProfile(sampleCV, sampleProfile)
			screen.Update(tea.WindowSizeMsg{Width: 35, Height: 8})
			view := screen.View()
			Expect(view).NotTo(BeEmpty())
		})

		It("uses custom theme when set via SetTheme", func() {
			screen := cv.NewCVReviewScreenWithProfile(sampleCV, sampleProfile)
			screen.SetTheme(themes.NewDefaultTheme())
			screen.Update(tea.WindowSizeMsg{Width: 120, Height: 40})
			view := screen.View()
			Expect(view).To(ContainSubstring("📋 CV Review"))
		})

		It("renders summary section without summary text in sections", func() {
			bullet := fixtures.CVBulletWith("b-nosummary", "s-exp", "Built payment system")
			group := fixtures.ContentGroupWithBullets("Engineering", []*career.CVBullet{bullet})
			section := fixtures.CVSectionWithContent("s-exp", "cv-nosummary", []*career.SectionContentGroup{group})
			section.Title = "Experience"
			section.SectionType = "experience"

			cvNoSummary := fixtures.CVViewWithSections("cv-nosummary", []*career.CVSection{section})
			cvNoSummary.Name = "No Summary CV"
			cvNoSummary.TargetRole = "Engineer"
			cvNoSummary.TargetAudience = "Manager"

			screen := cv.NewCVReviewScreenWithSummary(cvNoSummary, sampleProfile, sampleSummary)
			screen.Update(tea.WindowSizeMsg{Width: 120, Height: 40})
			view := screen.View()
			Expect(view).NotTo(ContainSubstring("📝 Summary"))
		})
	})

	Describe("Score and highlight coverage", func() {
		It("uses confidence only for master audience", func() {
			bullet := fixtures.CVBulletWith("b-master", "s1", "Designed distributed cache")
			group := fixtures.ContentGroupWithBullets("Engineering", []*career.CVBullet{bullet})
			section := fixtures.CVSectionWithContent("s1", "cv-master", []*career.SectionContentGroup{group})
			section.Title = "Experience"

			masterCV := fixtures.CVViewWithSections("cv-master", []*career.CVSection{section})
			masterCV.Name = "Master CV"
			masterCV.TargetRole = "Engineer"
			masterCV.TargetAudience = "master"

			screen := cv.NewCVReviewScreenWithProfile(masterCV, sampleProfile)
			screen.Update(tea.WindowSizeMsg{Width: 120, Height: 40})
			view := screen.View()
			Expect(view).To(ContainSubstring("Top Highlights"))
			Expect(view).To(ContainSubstring("Designed distributed cache"))
		})

		It("uses audience relevance for non-master audience with populated map", func() {
			bullet := fixtures.CVBulletWith("b-ar", "s1", "Led system redesign project")
			bullet.AudienceRelevance = map[string]float64{"Hiring Manager": 0.95}
			group := fixtures.ContentGroupWithBullets("Engineering", []*career.CVBullet{bullet})
			section := fixtures.CVSectionWithContent("s1", "cv-ar", []*career.SectionContentGroup{group})
			section.Title = "Experience"

			arCV := fixtures.CVViewWithSections("cv-ar", []*career.CVSection{section})
			arCV.Name = "AR CV"
			arCV.TargetRole = "Engineer"
			arCV.TargetAudience = "Hiring Manager"

			screen := cv.NewCVReviewScreenWithProfile(arCV, sampleProfile)
			screen.Update(tea.WindowSizeMsg{Width: 120, Height: 40})
			view := screen.View()
			Expect(view).To(ContainSubstring("Led system redesign project"))
		})

		It("uses confidence for empty audience", func() {
			bullet := fixtures.CVBulletWith("b-empty", "s1", "Built API gateway")
			group := fixtures.ContentGroupWithBullets("Engineering", []*career.CVBullet{bullet})
			section := fixtures.CVSectionWithContent("s1", "cv-empty-aud", []*career.SectionContentGroup{group})
			section.Title = "Experience"

			emptyAudCV := fixtures.CVViewWithSections("cv-empty-aud", []*career.CVSection{section})
			emptyAudCV.Name = "Empty Audience CV"
			emptyAudCV.TargetRole = "Engineer"
			emptyAudCV.TargetAudience = ""

			screen := cv.NewCVReviewScreenWithProfile(emptyAudCV, sampleProfile)
			screen.Update(tea.WindowSizeMsg{Width: 120, Height: 40})
			view := screen.View()
			Expect(view).To(ContainSubstring("Built API gateway"))
		})

		It("orders highlights by audience relevance score descending", func() {
			lowBullet := fixtures.CVBulletWith("b-low", "s1", "Low relevance achievement")
			lowBullet.AudienceRelevance = map[string]float64{"hiring_manager": 0.1}
			lowBullet.Confidence = 1.0

			highBullet := fixtures.CVBulletWith("b-high", "s1", "High relevance achievement")
			highBullet.AudienceRelevance = map[string]float64{"hiring_manager": 0.99}
			highBullet.Confidence = 1.0

			group := fixtures.ContentGroupWithBullets("Engineering", []*career.CVBullet{lowBullet, highBullet})
			section := fixtures.CVSectionWithContent("s1", "cv-order", []*career.SectionContentGroup{group})
			section.Title = "Experience"

			orderCV := fixtures.CVViewWithSections("cv-order", []*career.CVSection{section})
			orderCV.Name = "Order CV"
			orderCV.TargetRole = "Engineer"
			orderCV.TargetAudience = "hiring_manager"

			screen := cv.NewCVReviewScreenWithProfile(orderCV, sampleProfile)
			screen.Update(tea.WindowSizeMsg{Width: 120, Height: 40})
			view := screen.View()

			highIdx := strings.Index(view, "High relevance achievement")
			lowIdx := strings.Index(view, "Low relevance achievement")
			Expect(highIdx).To(BeNumerically(">", -1))
			Expect(lowIdx).To(BeNumerically(">", -1))
			Expect(highIdx).To(BeNumerically("<", lowIdx),
				"High-scoring bullet should appear before low-scoring bullet")
		})

		It("produces different ordering for different audiences", func() {
			techBullet := fixtures.CVBulletWith("b-tech", "s1", "Architected distributed system")
			techBullet.AudienceRelevance = map[string]float64{
				"hiring_manager": 0.3,
				"recruiter":      0.9,
			}
			techBullet.Confidence = 1.0

			bizBullet := fixtures.CVBulletWith("b-biz", "s1", "Delivered revenue growth project")
			bizBullet.AudienceRelevance = map[string]float64{
				"hiring_manager": 0.95,
				"recruiter":      0.2,
			}
			bizBullet.Confidence = 1.0

			group := fixtures.ContentGroupWithBullets("Engineering", []*career.CVBullet{techBullet, bizBullet})
			section := fixtures.CVSectionWithContent("s1", "cv-diff", []*career.SectionContentGroup{group})
			section.Title = "Experience"

			hmCV := fixtures.CVViewWithSections("cv-hm", []*career.CVSection{section})
			hmCV.Name = "HM CV"
			hmCV.TargetRole = "Engineer"
			hmCV.TargetAudience = "hiring_manager"

			hmScreen := cv.NewCVReviewScreenWithProfile(hmCV, sampleProfile)
			hmScreen.Update(tea.WindowSizeMsg{Width: 120, Height: 40})
			hmView := hmScreen.View()

			hmBizIdx := strings.Index(hmView, "Delivered revenue growth project")
			hmTechIdx := strings.Index(hmView, "Architected distributed system")
			Expect(hmBizIdx).To(BeNumerically("<", hmTechIdx),
				"For hiring_manager, business bullet should appear first")

			recruiterSection := fixtures.CVSectionWithContent("s1", "cv-diff-rec", []*career.SectionContentGroup{group})
			recruiterSection.Title = "Experience"
			recCV := fixtures.CVViewWithSections("cv-rec", []*career.CVSection{recruiterSection})
			recCV.Name = "Recruiter CV"
			recCV.TargetRole = "Engineer"
			recCV.TargetAudience = "recruiter"

			recScreen := cv.NewCVReviewScreenWithProfile(recCV, sampleProfile)
			recScreen.Update(tea.WindowSizeMsg{Width: 120, Height: 40})
			recView := recScreen.View()

			recBizIdx := strings.Index(recView, "Delivered revenue growth project")
			recTechIdx := strings.Index(recView, "Architected distributed system")
			Expect(recTechIdx).To(BeNumerically("<", recBizIdx),
				"For recruiter, technical bullet should appear first")
		})

		It("renders highlights using EnhancedText when Text is empty", func() {
			bullet := fixtures.CVBulletWith("b-enhanced", "s1", "")
			bullet.EnhancedText = "Enhanced version of the bullet"
			group := fixtures.ContentGroupWithBullets("Engineering", []*career.CVBullet{bullet})
			section := fixtures.CVSectionWithContent("s1", "cv-enhanced", []*career.SectionContentGroup{group})
			section.Title = "Experience"

			enhancedCV := fixtures.CVViewWithSections("cv-enhanced", []*career.CVSection{section})
			enhancedCV.Name = "Enhanced CV"
			enhancedCV.TargetRole = "Engineer"
			enhancedCV.TargetAudience = "Hiring Manager"

			screen := cv.NewCVReviewScreenWithProfile(enhancedCV, sampleProfile)
			screen.Update(tea.WindowSizeMsg{Width: 120, Height: 40})
			view := screen.View()
			Expect(view).To(ContainSubstring("Enhanced version of the bullet"))
		})
	})

	Describe("Existing constructors compatibility", func() {
		Context("NewCVReviewScreen", func() {
			It("constructs without panic and shows legacy layout", func() {
				Expect(func() {
					screen := cv.NewCVReviewScreen(sampleCV)
					screen.Update(tea.WindowSizeMsg{Width: 120, Height: 40})
					Expect(screen).NotTo(BeNil())
					view := screen.View()
					Expect(view).NotTo(ContainSubstring("YOUR SELECTIONS"))
					Expect(view).To(ContainSubstring("📊 Statistics"))
				}).NotTo(Panic())
			})
		})

		Context("NewCVReviewScreenWithProfile", func() {
			It("constructs without panic and shows legacy layout", func() {
				Expect(func() {
					screen := cv.NewCVReviewScreenWithProfile(sampleCV, sampleProfile)
					screen.Update(tea.WindowSizeMsg{Width: 120, Height: 40})
					screen.Update(tea.WindowSizeMsg{Width: 120, Height: 40})
					Expect(screen).NotTo(BeNil())
					view := screen.View()
					Expect(view).NotTo(ContainSubstring("YOUR SELECTIONS"))
					Expect(view).To(ContainSubstring("👤 Personal Details"))
				}).NotTo(Panic())
			})

			It("delegates to NewCVReviewScreenWithSummary with nil summary", func() {
				screen1 := cv.NewCVReviewScreenWithProfile(sampleCV, sampleProfile)
				screen2 := cv.NewCVReviewScreenWithSummary(sampleCV, sampleProfile, nil)

				// Both should produce the same layout (no YOUR SELECTIONS)
				view1 := screen1.View()
				view2 := screen2.View()

				Expect(view1).To(ContainSubstring("👤 Personal Details"))
				Expect(view2).To(ContainSubstring("👤 Personal Details"))
				Expect(view1).NotTo(ContainSubstring("YOUR SELECTIONS"))
				Expect(view2).NotTo(ContainSubstring("YOUR SELECTIONS"))
			})
		})
	})

	Describe("Title behavior", func() {
		Context("when summary is present", func() {
			It("shows 'CV Configuration Review' title", func() {
				screen := cv.NewCVReviewScreenWithSummary(sampleCV, sampleProfile, sampleSummary)
				screen.Update(tea.WindowSizeMsg{Width: 120, Height: 40})
				view := screen.View()

				Expect(view).To(ContainSubstring("📋 CV Configuration Review"))
			})
		})

		Context("when summary is nil", func() {
			It("shows 'CV Review' title", func() {
				screen := cv.NewCVReviewScreenWithSummary(sampleCV, sampleProfile, nil)
				screen.Update(tea.WindowSizeMsg{Width: 120, Height: 40})
				view := screen.View()

				Expect(view).To(ContainSubstring("📋 CV Review"))
				Expect(view).NotTo(ContainSubstring("Configuration"))
			})
		})
	})
})
