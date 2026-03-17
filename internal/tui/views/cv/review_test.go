package cv_test

import (
	"time"

	tea "github.com/charmbracelet/bubbletea"
	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"

	"github.com/baphled/kariya/internal/config"
	"github.com/baphled/kariya/internal/domain/career"
	"github.com/baphled/kariya/internal/testutil/fixtures"
	"github.com/baphled/kariya/internal/tui/views/cv"
	"github.com/baphled/kariya/internal/ui/display"
	"github.com/baphled/kariya/internal/ui/themes"
	"github.com/baphled/kariya/internal/ui/types"
	"github.com/baphled/kariya/internal/ui/uikit/widgets"
)

var _ = Describe("Review", func() {
	var (
		sampleCV        display.CVView
		sampleProfile   *config.ProfileConfig
		sampleSummary   *cv.GenerationSummary
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

		bullet1 := fixtures.CVBulletWith("b1-1", "section-1", "Led development of microservices architecture")
		bullet2 := fixtures.CVBulletWith("b1-2", "section-1", "Improved system performance by 40%")
		group1 := fixtures.ContentGroupWithBullets("Senior Developer at TechCorp", []*career.CVBullet{bullet1, bullet2})
		section1 := fixtures.CVSectionWithContent("section-1", "test-cv-123", []*career.SectionContentGroup{group1})
		section1.Title = "Professional Experience"
		section1.SectionType = "experience"

		bullet3 := fixtures.CVBulletWith("b2-1", "section-2", "Go, Python, JavaScript")
		bullet4 := fixtures.CVBulletWith("b2-2", "section-2", "SQL, NoSQL databases")
		group2 := fixtures.ContentGroupWithBullets("Programming Languages", []*career.CVBullet{bullet3, bullet4})
		section2 := fixtures.CVSectionWithContent("section-2", "test-cv-123", []*career.SectionContentGroup{group2})
		section2.Title = "Technical Skills"
		section2.SectionType = "skills"

		summaryBullet := fixtures.CVBulletWith("b3-1", "section-3", "Experienced software engineer with 8+ years")
		summaryGroup := fixtures.ContentGroupWithBullets("Professional Summary", []*career.CVBullet{summaryBullet})
		section3 := fixtures.CVSectionWithContent("section-3", "test-cv-123", []*career.SectionContentGroup{summaryGroup})
		section3.Title = "Summary"
		section3.SectionType = "summary"
		section3.Summary = "Experienced software engineer with 8+ years of backend development"

		sampleCVDomain := fixtures.CVViewWithSections("test-cv-123", []*career.CVSection{section1, section2, section3})
		sampleCVDomain.Name = "John Doe CV"
		sampleCVDomain.TargetRole = "Senior Software Engineer"
		sampleCVDomain.TargetAudience = "Hiring Manager"
		sampleCVDomain.EventFilters = map[string]interface{}{"company": "TechCorp"}
		sampleCVDomain.GeneratedAt = time.Now()
		sampleCVDomain.SourceEventCount = 8
		sampleCVDomain.SourceFactCount = 12
		sampleCV = display.CVViewFromDomain(sampleCVDomain)

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

		sampleSummary = &cv.GenerationSummary{
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

	Describe("Construction", func() {
		It("should create a Review with the provided fields", func() {
			view := cv.NewReview(sampleCV, sampleProfile, sampleSummary)
			Expect(view).NotTo(BeNil())
		})

		It("should allow nil summary without panic", func() {
			Expect(func() {
				view := cv.NewReview(sampleCV, sampleProfile, nil)
				Expect(view).NotTo(BeNil())
			}).NotTo(Panic())
		})

		It("should allow nil profile config without panic", func() {
			Expect(func() {
				view := cv.NewReview(sampleCV, nil, sampleSummary)
				Expect(view).NotTo(BeNil())
			}).NotTo(Panic())
		})
	})

	Describe("Init", func() {
		It("should return nil command", func() {
			view := cv.NewReview(sampleCV, sampleProfile, sampleSummary)
			cmd := view.Init()
			Expect(cmd).To(BeNil())
		})
	})

	Describe("Update", func() {
		var view *cv.Review

		BeforeEach(func() {
			view = cv.NewReview(sampleCV, sampleProfile, sampleSummary)
		})

		It("should update terminal info on WindowSizeMsg", func() {
			cmd, result := view.Update(tea.WindowSizeMsg{Width: 120, Height: 40})
			Expect(cmd).To(BeNil())
			Expect(result).To(BeNil())
			Expect(view.GetTerminalWidth()).To(Equal(120))
			Expect(view.GetTerminalHeight()).To(Equal(40))
		})

		It("should return CancelViewResult on esc", func() {
			_, result := view.Update(tea.KeyMsg{Type: tea.KeyEsc})
			Expect(result).NotTo(BeNil())
			Expect(result.Type()).To(Equal(widgets.ResultCancel))
		})

		It("should return CancelViewResult type on esc", func() {
			_, result := view.Update(tea.KeyMsg{Type: tea.KeyEsc})
			_, ok := result.(*widgets.CancelViewResult)
			Expect(ok).To(BeTrue())
		})

		It("should return NavigateViewResult for preview on enter", func() {
			_, result := view.Update(tea.KeyMsg{Type: tea.KeyEnter})
			nav, ok := result.(*widgets.NavigateViewResult)
			Expect(ok).To(BeTrue())
			payload, ok := nav.Data().(cv.Nav)
			Expect(ok).To(BeTrue())
			Expect(payload.Action).To(Equal(cv.ActionPreview))
		})

		It("should return NavigateViewResult for preview on p key", func() {
			_, result := view.Update(tea.KeyMsg{Type: tea.KeyRunes, Runes: []rune{'p'}})
			nav, ok := result.(*widgets.NavigateViewResult)
			Expect(ok).To(BeTrue())
			payload, ok := nav.Data().(cv.Nav)
			Expect(ok).To(BeTrue())
			Expect(payload.Action).To(Equal(cv.ActionPreview))
		})

		It("should return NavigateViewResult for export on x key", func() {
			_, result := view.Update(tea.KeyMsg{Type: tea.KeyRunes, Runes: []rune{'x'}})
			nav, ok := result.(*widgets.NavigateViewResult)
			Expect(ok).To(BeTrue())
			payload, ok := nav.Data().(cv.Nav)
			Expect(ok).To(BeTrue())
			Expect(payload.Action).To(Equal(cv.ActionExport))
		})

		It("should return NavigateViewResult for edit on e key", func() {
			_, result := view.Update(tea.KeyMsg{Type: tea.KeyRunes, Runes: []rune{'e'}})
			nav, ok := result.(*widgets.NavigateViewResult)
			Expect(ok).To(BeTrue())
			payload, ok := nav.Data().(cv.Nav)
			Expect(ok).To(BeTrue())
			Expect(payload.Action).To(Equal(cv.ActionEdit))
		})

		It("should go to top on g key", func() {
			cmd, result := view.Update(tea.KeyMsg{Type: tea.KeyRunes, Runes: []rune{'g'}})
			Expect(cmd).To(BeNil())
			Expect(result).To(BeNil())
		})

		It("should go to bottom on G key", func() {
			cmd, result := view.Update(tea.KeyMsg{Type: tea.KeyRunes, Runes: []rune{'G'}})
			Expect(cmd).To(BeNil())
			Expect(result).To(BeNil())
		})

		It("should return nil result for scroll keys when viewport not ready", func() {
			cmd, result := view.Update(tea.KeyMsg{Type: tea.KeyDown})
			Expect(cmd).To(BeNil())
			Expect(result).To(BeNil())
		})

		It("should handle scroll keys when viewport is ready", func() {
			view.Update(tea.WindowSizeMsg{Width: 120, Height: 200})
			view.RenderContent()
			_, result := view.Update(tea.KeyMsg{Type: tea.KeyRunes, Runes: []rune{'j'}})
			Expect(result).To(BeNil())
		})

		It("should handle up key when viewport is ready", func() {
			view.Update(tea.WindowSizeMsg{Width: 120, Height: 200})
			view.RenderContent()
			_, result := view.Update(tea.KeyMsg{Type: tea.KeyUp})
			Expect(result).To(BeNil())
		})

		It("should handle down key when viewport is ready", func() {
			view.Update(tea.WindowSizeMsg{Width: 120, Height: 200})
			view.RenderContent()
			_, result := view.Update(tea.KeyMsg{Type: tea.KeyDown})
			Expect(result).To(BeNil())
		})

		DescribeTable("should handle scroll keys when viewport is ready", func(msg tea.KeyMsg) {
			view.Update(tea.WindowSizeMsg{Width: 120, Height: 200})
			view.RenderContent()
			_, result := view.Update(msg)
			Expect(result).To(BeNil())
		},
			Entry("up", tea.KeyMsg{Type: tea.KeyUp}),
			Entry("down", tea.KeyMsg{Type: tea.KeyDown}),
			Entry("k", tea.KeyMsg{Type: tea.KeyRunes, Runes: []rune{'k'}}),
			Entry("pgup", tea.KeyMsg{Type: tea.KeyPgUp}),
			Entry("pgdown", tea.KeyMsg{Type: tea.KeyPgDown}),
			Entry("ctrl+u", tea.KeyMsg{Type: tea.KeyCtrlU}),
			Entry("ctrl+d", tea.KeyMsg{Type: tea.KeyCtrlD}),
		)

		It("should return nil for unhandled messages", func() {
			cmd, result := view.Update(tea.MouseMsg{})
			Expect(cmd).To(BeNil())
			Expect(result).To(BeNil())
		})
	})

	Describe("GetCV", func() {
		It("should return the CV data", func() {
			view := cv.NewReview(sampleCV, sampleProfile, sampleSummary)
			Expect(view.GetCV()).To(Equal(sampleCV))
		})

		It("should return nil when CV is nil", func() {
			view := cv.NewReview(display.CVView{}, sampleProfile, sampleSummary)
			Expect(view.GetCV()).To(Equal(display.CVView{}))
		})
	})

	Describe("RenderContent", func() {
		It("should render CV Configuration Review title when summary provided", func() {
			view := cv.NewReview(sampleCV, sampleProfile, sampleSummary)
			content := view.RenderContent()
			Expect(content).To(ContainSubstring("CV Configuration Review"))
		})

		It("should render CV Review title when summary is nil", func() {
			view := cv.NewReview(sampleCV, sampleProfile, nil)
			content := view.RenderContent()
			Expect(content).To(ContainSubstring("CV Review"))
		})

		It("should render summary-specific sections when summary provided", func() {
			view := cv.NewReview(sampleCV, sampleProfile, sampleSummary)
			view.Update(tea.WindowSizeMsg{Width: 120, Height: 200})
			content := view.RenderContent()
			Expect(content).To(ContainSubstring("Generation Settings"))
			Expect(content).To(ContainSubstring("Backend Development"))
			Expect(content).To(ContainSubstring("Go, Docker, Kubernetes"))
		})

		It("should render legacy sections when summary is nil", func() {
			view := cv.NewReview(sampleCV, sampleProfile, nil)
			view.Update(tea.WindowSizeMsg{Width: 120, Height: 200})
			content := view.RenderContent()
			Expect(content).To(ContainSubstring("Personal Details"))
			Expect(content).To(ContainSubstring("Statistics"))
			Expect(content).To(ContainSubstring("Sections"))
			Expect(content).NotTo(ContainSubstring("Generation Settings"))
		})

		It("should show no data message when CV is nil", func() {
			view := cv.NewReview(display.CVView{}, sampleProfile, sampleSummary)
			content := view.RenderContent()
			Expect(content).To(ContainSubstring("No CV data available"))
			Expect(content).To(ContainSubstring("esc: back"))
		})

		It("should render summary section when CV summary is present", func() {
			summarySection := fixtures.CVSectionWithSummary("section-summary", "cv-sum", "Senior Engineer with expertise in Go and distributed systems.")
			cvWithSummaryDomain := fixtures.CVViewWithSections("cv-sum", []*career.CVSection{summarySection})
			cvWithSummaryDomain.Name = "Test CV"
			cvWithSummaryDomain.TargetRole = "Engineer"
			cvWithSummaryDomain.TargetAudience = "Hiring Manager"
			cvWithSummary := display.CVViewFromDomain(cvWithSummaryDomain)

			view := cv.NewReview(cvWithSummary, sampleProfile, sampleSummary)
			view.Update(tea.WindowSizeMsg{Width: 120, Height: 200})
			content := view.RenderContent()
			Expect(content).To(ContainSubstring("Summary"))
			Expect(content).To(ContainSubstring("Senior Engineer"))
		})

		It("should omit summary section when CV has no summary text", func() {
			bullet := fixtures.CVBulletWith("b-nosummary", "s-exp", "Built payment system")
			group := fixtures.ContentGroupWithBullets("Engineering", []*career.CVBullet{bullet})
			section := fixtures.CVSectionWithContent("s-exp", "cv-nosummary", []*career.SectionContentGroup{group})
			section.Title = "Experience"
			section.SectionType = "experience"

			cvNoSummaryDomain := fixtures.CVViewWithSections("cv-nosummary", []*career.CVSection{section})
			cvNoSummaryDomain.Name = "No Summary CV"
			cvNoSummaryDomain.TargetRole = "Engineer"
			cvNoSummaryDomain.TargetAudience = "Manager"
			cvNoSummary := display.CVViewFromDomain(cvNoSummaryDomain)

			view := cv.NewReview(cvNoSummary, sampleProfile, sampleSummary)
			view.Update(tea.WindowSizeMsg{Width: 120, Height: 200})
			content := view.RenderContent()
			Expect(content).NotTo(ContainSubstring("📝 Summary"))
		})

		It("should omit summary section when summary text is empty", func() {
			summarySection := fixtures.CVSectionWithSummary("section-summary-empty", "cv-empty-summary", "")
			cvEmptySummaryDomain := fixtures.CVViewWithSections("cv-empty-summary", []*career.CVSection{summarySection})
			cvEmptySummaryDomain.Name = "Empty Summary CV"
			cvEmptySummaryDomain.TargetRole = "Engineer"
			cvEmptySummaryDomain.TargetAudience = "Hiring Manager"
			cvEmptySummary := display.CVViewFromDomain(cvEmptySummaryDomain)

			view := cv.NewReview(cvEmptySummary, sampleProfile, sampleSummary)
			view.Update(tea.WindowSizeMsg{Width: 120, Height: 200})
			content := view.RenderContent()
			Expect(content).NotTo(ContainSubstring("📝 Summary"))
		})

		It("should omit summary section when sections are nil", func() {
			nilSectionsCVDomain := fixtures.CVViewWithSections("cv-nil-sections", nil)
			nilSectionsCVDomain.Name = "Nil Sections CV"
			nilSectionsCVDomain.TargetRole = "Engineer"
			nilSectionsCVDomain.TargetAudience = "Hiring Manager"
			nilSectionsCV := display.CVViewFromDomain(nilSectionsCVDomain)

			view := cv.NewReview(nilSectionsCV, sampleProfile, sampleSummary)
			view.Update(tea.WindowSizeMsg{Width: 120, Height: 200})
			content := view.RenderContent()
			Expect(content).To(ContainSubstring("🎯 Generation Settings"))
			Expect(content).NotTo(ContainSubstring("📝 Summary"))
		})

		It("should render summary section with long text", func() {
			summaryText := "This summary is intentionally long to trigger wrapping across multiple lines in the review view output."
			summarySection := fixtures.CVSectionWithSummary("section-summary-long", "cv-long-summary", summaryText)
			cvLongSummaryDomain := fixtures.CVViewWithSections("cv-long-summary", []*career.CVSection{summarySection})
			cvLongSummaryDomain.Name = "Long Summary CV"
			cvLongSummaryDomain.TargetRole = "Engineer"
			cvLongSummaryDomain.TargetAudience = "Hiring Manager"
			cvLongSummary := display.CVViewFromDomain(cvLongSummaryDomain)

			view := cv.NewReview(cvLongSummary, sampleProfile, sampleSummary)
			view.Update(tea.WindowSizeMsg{Width: 120, Height: 200})
			content := view.RenderContent()
			Expect(content).To(ContainSubstring("📝 Summary"))
			Expect(content).To(ContainSubstring("intentionally long"))
		})

		It("should render highlights using enhanced text when text is empty", func() {
			bullet := fixtures.CVBulletWith("b-enhanced", "s1", "")
			bullet.EnhancedText = "Enhanced version of the bullet"
			group := fixtures.ContentGroupWithBullets("Engineering", []*career.CVBullet{bullet})
			section := fixtures.CVSectionWithContent("s1", "cv-enhanced", []*career.SectionContentGroup{group})
			section.Title = "Experience"

			enhancedCVDomain := fixtures.CVViewWithSections("cv-enhanced", []*career.CVSection{section})
			enhancedCVDomain.Name = "Enhanced CV"
			enhancedCVDomain.TargetRole = "Engineer"
			enhancedCVDomain.TargetAudience = "Hiring Manager"
			enhancedCV := display.CVViewFromDomain(enhancedCVDomain)

			view := cv.NewReview(enhancedCV, sampleProfile, nil)
			view.Update(tea.WindowSizeMsg{Width: 120, Height: 200})
			content := view.RenderContent()
			Expect(content).To(ContainSubstring("Enhanced version of the bullet"))
		})

		It("should score bullets using audience relevance when audience is set", func() {
			bullet := fixtures.CVBulletWith("b-aud", "s1", "Audience-relevant bullet")
			bullet.Confidence = 0.8
			bullet.AudienceRelevance = map[string]float64{"Hiring Manager": 0.95}
			group := fixtures.ContentGroupWithBullets("Engineering", []*career.CVBullet{bullet})
			section := fixtures.CVSectionWithContent("s1", "cv-aud", []*career.SectionContentGroup{group})
			section.Title = "Experience"

			audienceCVDomain := fixtures.CVViewWithSections("cv-aud", []*career.CVSection{section})
			audienceCVDomain.Name = "Audience CV"
			audienceCVDomain.TargetRole = "Engineer"
			audienceCVDomain.TargetAudience = "Hiring Manager"
			audienceCV := display.CVViewFromDomain(audienceCVDomain)

			view := cv.NewReview(audienceCV, sampleProfile, nil)
			view.Update(tea.WindowSizeMsg{Width: 120, Height: 200})
			content := view.RenderContent()
			Expect(content).To(ContainSubstring("Audience-relevant bullet"))
		})

		It("should use confidence score when audience is master", func() {
			bullet := fixtures.CVBulletWith("b-master", "s1", "Master audience bullet")
			bullet.Confidence = 0.9
			bullet.AudienceRelevance = map[string]float64{"Hiring Manager": 0.5}
			group := fixtures.ContentGroupWithBullets("Engineering", []*career.CVBullet{bullet})
			section := fixtures.CVSectionWithContent("s1", "cv-master", []*career.SectionContentGroup{group})
			section.Title = "Experience"

			masterCVDomain := fixtures.CVViewWithSections("cv-master", []*career.CVSection{section})
			masterCVDomain.Name = "Master CV"
			masterCVDomain.TargetRole = "Engineer"
			masterCVDomain.TargetAudience = "master"
			masterCV := display.CVViewFromDomain(masterCVDomain)

			view := cv.NewReview(masterCV, sampleProfile, nil)
			view.Update(tea.WindowSizeMsg{Width: 120, Height: 200})
			content := view.RenderContent()
			Expect(content).To(ContainSubstring("Master audience bullet"))
		})

		It("should use confidence score when audience is empty", func() {
			bullet := fixtures.CVBulletWith("b-noaud", "s1", "No audience bullet")
			bullet.Confidence = 0.85
			group := fixtures.ContentGroupWithBullets("Engineering", []*career.CVBullet{bullet})
			section := fixtures.CVSectionWithContent("s1", "cv-noaud", []*career.SectionContentGroup{group})
			section.Title = "Experience"

			noAudienceCVDomain := fixtures.CVViewWithSections("cv-noaud", []*career.CVSection{section})
			noAudienceCVDomain.Name = "No Audience CV"
			noAudienceCVDomain.TargetRole = "Engineer"
			noAudienceCVDomain.TargetAudience = ""
			noAudienceCV := display.CVViewFromDomain(noAudienceCVDomain)

			view := cv.NewReview(noAudienceCV, sampleProfile, nil)
			view.Update(tea.WindowSizeMsg{Width: 120, Height: 200})
			content := view.RenderContent()
			Expect(content).To(ContainSubstring("No audience bullet"))
		})

		It("should handle bullet with nil AudienceRelevance map", func() {
			bullet := fixtures.CVBulletWith("b-nilmap", "s1", "Nil map bullet")
			bullet.Confidence = 0.7
			bullet.AudienceRelevance = nil
			group := fixtures.ContentGroupWithBullets("Engineering", []*career.CVBullet{bullet})
			section := fixtures.CVSectionWithContent("s1", "cv-nilmap", []*career.SectionContentGroup{group})
			section.Title = "Experience"

			nilMapCVDomain := fixtures.CVViewWithSections("cv-nilmap", []*career.CVSection{section})
			nilMapCVDomain.Name = "Nil Map CV"
			nilMapCVDomain.TargetRole = "Engineer"
			nilMapCVDomain.TargetAudience = "Recruiter"
			nilMapCV := display.CVViewFromDomain(nilMapCVDomain)

			view := cv.NewReview(nilMapCV, sampleProfile, nil)
			view.Update(tea.WindowSizeMsg{Width: 120, Height: 200})
			content := view.RenderContent()
			Expect(content).To(ContainSubstring("Nil map bullet"))
		})

		It("should render empty sections message", func() {
			emptyCVDomain := fixtures.CVViewWithSections("cv-empty", []*career.CVSection{})
			emptyCVDomain.Name = "Empty CV"
			emptyCVDomain.TargetRole = "Engineer"
			emptyCVDomain.TargetAudience = "Manager"
			emptyCV := display.CVViewFromDomain(emptyCVDomain)

			view := cv.NewReview(emptyCV, sampleProfile, nil)
			view.Update(tea.WindowSizeMsg{Width: 120, Height: 200})
			content := view.RenderContent()
			Expect(content).To(ContainSubstring("0 sections"))
		})

		It("should render singular bullet text for single bullet section", func() {
			bullet := fixtures.CVBulletWith("b-single", "s1", "Only bullet")
			group := fixtures.ContentGroupWithBullets("Company X", []*career.CVBullet{bullet})
			section := fixtures.CVSectionWithContent("s1", "cv-single", []*career.SectionContentGroup{group})
			section.Title = "Experience"
			section.SectionType = "experience"

			singleCVDomain := fixtures.CVViewWithSections("cv-single", []*career.CVSection{section})
			singleCVDomain.Name = "Single Bullet CV"
			singleCVDomain.TargetRole = "Engineer"
			singleCVDomain.TargetAudience = "Manager"
			singleCV := display.CVViewFromDomain(singleCVDomain)

			view := cv.NewReview(singleCV, sampleProfile, nil)
			view.Update(tea.WindowSizeMsg{Width: 120, Height: 200})
			content := view.RenderContent()
			Expect(content).To(ContainSubstring("1 bullet)"))
		})

		It("should render generation settings with all fields populated", func() {
			view := cv.NewReview(sampleCV, sampleProfile, sampleSummary)
			view.Update(tea.WindowSizeMsg{Width: 120, Height: 200})
			content := view.RenderContent()
			Expect(content).To(ContainSubstring("Profile:"))
			Expect(content).To(ContainSubstring("Audience:"))
			Expect(content).To(ContainSubstring("Tech Focus:"))
			Expect(content).To(ContainSubstring("Technologies:"))
			Expect(content).To(ContainSubstring("Focus Area:"))
			Expect(content).To(ContainSubstring("Skills Format:"))
			Expect(content).To(ContainSubstring("CV Length:"))
		})

		It("should render generation settings without profile when nil", func() {
			noProfileSummary := &cv.GenerationSummary{
				SelectedProfile:  nil,
				SelectedAudience: "Recruiter",
				TechnologyFocus:  "Frontend",
				Technologies:     []string{"React"},
				FocusArea:        "UI",
				SkillsFormat:     "list",
				SkillsLimit:      0,
				CVLength:         "brief",
			}
			view := cv.NewReview(sampleCV, sampleProfile, noProfileSummary)
			view.Update(tea.WindowSizeMsg{Width: 120, Height: 200})
			content := view.RenderContent()
			Expect(content).NotTo(ContainSubstring("Profile:"))
			Expect(content).To(ContainSubstring("Audience:"))
		})

		It("should render generation settings with profile and options", func() {
			customSummary := &cv.GenerationSummary{
				SelectedProfile:  &types.CVProfile{Name: "Test Profile"},
				SelectedAudience: "",
				TechnologyFocus:  "Infrastructure",
				Technologies:     []string{"Go", "Terraform"},
				FocusArea:        "Reliability",
				SkillsFormat:     "grouped",
				SkillsLimit:      3,
				CVLength:         "compact",
			}

			view := cv.NewReview(sampleCV, sampleProfile, customSummary)
			view.Update(tea.WindowSizeMsg{Width: 120, Height: 200})
			content := view.RenderContent()
			Expect(content).To(ContainSubstring("Test Profile"))
			Expect(content).To(ContainSubstring("Tech Focus:"))
			Expect(content).To(ContainSubstring("Infrastructure"))
			Expect(content).To(ContainSubstring("Technologies:"))
			Expect(content).To(ContainSubstring("Go, Terraform"))
			Expect(content).To(ContainSubstring("Focus Area:"))
			Expect(content).To(ContainSubstring("Reliability"))
			Expect(content).To(ContainSubstring("Skills Format:"))
			Expect(content).To(ContainSubstring("grouped (3 max)"))
			Expect(content).To(ContainSubstring("CV Length:"))
			Expect(content).To(ContainSubstring("compact"))
		})

		It("should render skills format without limit when limit is zero", func() {
			customSummary := &cv.GenerationSummary{
				SelectedProfile: &types.CVProfile{Name: "Test Profile"},
				SkillsFormat:    "inline",
				SkillsLimit:     0,
				CVLength:        "brief",
			}

			view := cv.NewReview(sampleCV, sampleProfile, customSummary)
			view.Update(tea.WindowSizeMsg{Width: 120, Height: 200})
			content := view.RenderContent()
			Expect(content).To(ContainSubstring("Skills Format:"))
			Expect(content).To(ContainSubstring("inline"))
			Expect(content).NotTo(ContainSubstring("max)"))
		})

		It("should render with set theme", func() {
			view := cv.NewReview(sampleCV, sampleProfile, sampleSummary)
			view.SetTheme(themes.NewDefaultTheme())
			content := view.RenderContent()
			Expect(content).To(ContainSubstring("CV Configuration Review"))
		})

		It("should render with small viewport height", func() {
			view := cv.NewReview(sampleCV, sampleProfile, sampleSummary)
			view.Update(tea.WindowSizeMsg{Width: 50, Height: 3})
			content := view.RenderContent()
			Expect(content).ToNot(BeEmpty())
		})

		It("should render with small viewport width", func() {
			view := cv.NewReview(sampleCV, sampleProfile, sampleSummary)
			view.Update(tea.WindowSizeMsg{Width: 20, Height: 40})
			content := view.RenderContent()
			Expect(content).ToNot(BeEmpty())
		})
	})

	Describe("HelpText", func() {
		It("should include scroll and action badges", func() {
			view := cv.NewReview(sampleCV, sampleProfile, sampleSummary)
			help := view.HelpText()
			Expect(help).To(ContainSubstring("Scroll"))
			Expect(help).To(ContainSubstring("Preview"))
			Expect(help).To(ContainSubstring("Export"))
			Expect(help).To(ContainSubstring("Edit"))
			Expect(help).To(ContainSubstring("Back"))
		})

		It("should include scroll percentage when viewport is ready", func() {
			view := cv.NewReview(sampleCV, sampleProfile, sampleSummary)
			view.Update(tea.WindowSizeMsg{Width: 120, Height: 200})
			view.RenderContent()
			help := view.HelpText()
			Expect(help).To(ContainSubstring("%"))
		})
	})
})
