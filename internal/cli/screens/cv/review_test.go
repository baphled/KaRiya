package cv_test

import (
	"time"

	tea "github.com/charmbracelet/bubbletea"

	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"

	"github.com/baphled/kariya/internal/cli/screens/cv"
	"github.com/baphled/kariya/internal/cli/types"
	"github.com/baphled/kariya/internal/config"
	"github.com/baphled/kariya/internal/domain/career"
)

var _ = Describe("ReviewScreen", func() {
	var (
		sampleCV         *career.CVView
		sampleProfile    *config.ProfileConfig
		sampleSummary    *cv.GenerationSummary
		sampleCVProfile  *types.CVProfile
	)

	BeforeEach(func() {
		sampleCVProfile = &types.CVProfile{
			ID:             "test-profile-1",
			Name:           "Senior Backend Engineer",
			TargetRole:     "senior_ic",
			TargetAudience: "hiring_manager",
			Description:    "Test CV profile for backend engineering",
		}

		sampleCV = &career.CVView{
			ID:               "test-cv-123",
			Name:             "John Doe CV",
			TargetRole:       "Senior Software Engineer",
			TargetAudience:   "Hiring Manager",
			EventFilters:     map[string]interface{}{"company": "TechCorp"},
			GeneratedAt:      time.Now(),
			SourceEventCount: 8,
			SourceFactCount:  12,
			Sections: []*career.CVSection{
				{
					ID:          "section-1",
					Title:       "Professional Experience",
					SectionType: "experience",
				Content: []*career.SectionContentGroup{
					{
						Header:  "Senior Developer at TechCorp",
						Bullets: []*career.CVBullet{
							{Text: "Led development of microservices architecture"},
							{Text: "Improved system performance by 40%"},
						},
					},
					},
				},
				{
					ID:          "section-2",
					Title:       "Technical Skills",
					SectionType: "skills",
				Content: []*career.SectionContentGroup{
					{
						Header:  "Programming Languages",
						Bullets: []*career.CVBullet{
							{Text: "Go, Python, JavaScript"},
							{Text: "SQL, NoSQL databases"},
						},
					},
					},
				},
				{
					ID:          "section-3",
					Title:       "Summary",
					SectionType: "summary",
				Content: []*career.SectionContentGroup{
					{
						Header:  "Professional Summary",
						Bullets: []*career.CVBullet{
							{Text: "Experienced software engineer with 8+ years"},
						},
					},
					},
				},
			},
		}

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

	Describe("Constructor tests", func() {
		Context("NewCVReviewScreenWithSummary", func() {
			It("creates a screen with the summary field set", func() {
				screen := cv.NewCVReviewScreenWithSummary(sampleCV, sampleProfile, sampleSummary)
				screen.Update(tea.WindowSizeMsg{Width: 120, Height: 40})

				Expect(screen).NotTo(BeNil())
				Expect(screen.GetCV()).To(Equal(sampleCV))
				// We can verify the summary is set by checking the View() output contains summary-specific content
				view := screen.View()
				Expect(view).To(ContainSubstring("YOUR SELECTIONS"))
			})

			It("works with nil summary", func() {
				screen := cv.NewCVReviewScreenWithSummary(sampleCV, sampleProfile, nil)
				screen.Update(tea.WindowSizeMsg{Width: 120, Height: 40})

				Expect(screen).NotTo(BeNil())
				Expect(screen.GetCV()).To(Equal(sampleCV))
				// Without summary, should not show summary-specific content
				view := screen.View()
				Expect(view).NotTo(ContainSubstring("YOUR SELECTIONS"))
			})

			It("works with nil profile config", func() {
				screen := cv.NewCVReviewScreenWithSummary(sampleCV, nil, sampleSummary)
				screen.Update(tea.WindowSizeMsg{Width: 120, Height: 40})

				Expect(screen).NotTo(BeNil())
				Expect(screen.GetCV()).To(Equal(sampleCV))
				view := screen.View()
				Expect(view).To(ContainSubstring("YOUR SELECTIONS"))
			})
		})
	})

	Describe("View() rendering", func() {
		Context("when summary is non-nil", func() {
			It("contains YOUR SELECTIONS section with profile/audience/tech/format fields", func() {
				screen := cv.NewCVReviewScreenWithSummary(sampleCV, sampleProfile, sampleSummary)
				screen.Update(tea.WindowSizeMsg{Width: 120, Height: 40})
				screen.Update(tea.WindowSizeMsg{Width: 120, Height: 40})
				view := screen.View()

				Expect(view).To(ContainSubstring("YOUR SELECTIONS"))

				// Check profile settings
				Expect(view).To(ContainSubstring("Profile Settings"))
				Expect(view).To(ContainSubstring("Senior Backend Engineer")) // Profile name
				Expect(view).To(ContainSubstring("Hiring Manager"))          // Selected audience

				// Check technology focus
				Expect(view).To(ContainSubstring("Technology Focus"))
				Expect(view).To(ContainSubstring("Backend Development"))      // TechnologyFocus
				Expect(view).To(ContainSubstring("Go, Docker, Kubernetes"))   // Technologies joined
				Expect(view).To(ContainSubstring("Microservices Architecture")) // FocusArea

				// Check format options
				Expect(view).To(ContainSubstring("Format Options"))
				Expect(view).To(ContainSubstring("comprehensive"))    // CVLength
				Expect(view).To(ContainSubstring("grouped (15 max)")) // SkillsFormat with limit
			})

			It("contains GENERATED OUTPUT section with section count, bullet count, and sources", func() {
				screen := cv.NewCVReviewScreenWithSummary(sampleCV, sampleProfile, sampleSummary)
				screen.Update(tea.WindowSizeMsg{Width: 120, Height: 40})
				screen.Update(tea.WindowSizeMsg{Width: 120, Height: 40})
				view := screen.View()

				Expect(view).To(ContainSubstring("GENERATED OUTPUT"))

				// Check statistics
				Expect(view).To(ContainSubstring("📊 Statistics"))
				Expect(view).To(ContainSubstring("3 sections"))     // SectionCount
				Expect(view).To(ContainSubstring("5 bullet points")) // TotalBullets
				Expect(view).To(ContainSubstring("8 events, 12 facts")) // SourceEventCount, SourceFactCount
			})

			It("handles empty or default values gracefully", func() {
				emptySummary := &cv.GenerationSummary{
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

				Expect(view).To(ContainSubstring("YOUR SELECTIONS"))
				Expect(view).To(ContainSubstring("GENERATED OUTPUT"))

				// Should show dashes for empty values
				Expect(view).To(ContainSubstring("Profile: ")) // Should handle nil profile gracefully
				Expect(view).To(ContainSubstring("Technologies: -"))
				Expect(view).To(ContainSubstring("Skills: -"))
				Expect(view).To(ContainSubstring("CV Length: -"))
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

				Expect(view).To(ContainSubstring("John Doe"))          // Name from profile
				Expect(view).To(ContainSubstring("john.doe@example.com")) // Email from profile
				Expect(view).To(ContainSubstring("San Francisco, CA"))    // Location from profile
			})

			It("shows CV details and statistics from CVView", func() {
				screen := cv.NewCVReviewScreenWithSummary(sampleCV, sampleProfile, nil)
				screen.Update(tea.WindowSizeMsg{Width: 120, Height: 40})
				view := screen.View()

				// CV details
				Expect(view).To(ContainSubstring("John Doe CV"))          // CV name
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
