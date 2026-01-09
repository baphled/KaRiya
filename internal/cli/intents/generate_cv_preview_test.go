package intents

import (
	"time"

	"github.com/baphled/kariya/internal/domain/career"
	tea "github.com/charmbracelet/bubbletea"
	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
)

var _ = Describe("GenerateCVIntent - Structure-Aware Preview", func() {
	var (
		intent *GenerateCVIntent
		ctx    *GenerateCVContext
	)

	// Helper to create a CV with sections and bullets
	createTestCV := func() *career.CVView {
		return &career.CVView{
			ID:               "test-cv",
			Name:             "Senior IC CV",
			TargetRole:       "senior_ic",
			TargetAudience:   "hiring_manager",
			GeneratedAt:      time.Now(),
			SourceEventCount: 5,
			SourceFactCount:  10,
			Sections: []*career.CVSection{
				{
					ID:          "section-summary",
					CVViewID:    "test-cv",
					SectionType: "summary",
					Title:       "Summary",
					Order:       0,
					Summary:     "Experienced software engineer with 10+ years in backend development.",
				},
				{
					ID:          "section-experience",
					CVViewID:    "test-cv",
					SectionType: "experience",
					Title:       "Experience",
					Order:       1,
					Content: []*career.SectionContentGroup{
						{
							Header:    "TechCorp",
							StartDate: "Jan 2020",
							EndDate:   "Present",
							Bullets: []*career.CVBullet{
								{
									ID:              "bullet-1",
									SectionID:       "section-experience",
									Text:            "Led migration to microservices architecture, reducing deployment time by 60%",
									SourceEventIDs:  []string{"event-1"},
									Rank:            0.9,
									Confidence:      0.85,
									InclusionReason: "ownership",
								},
								{
									ID:              "bullet-2",
									SectionID:       "section-experience",
									Text:            "Implemented CI/CD pipeline using GitHub Actions",
									SourceEventIDs:  []string{"event-2"},
									Rank:            0.7,
									Confidence:      0.70, // Below threshold for narrative
									InclusionReason: "execution",
								},
							},
						},
						{
							Header:    "StartupCo",
							StartDate: "Jun 2017",
							EndDate:   "Dec 2019",
							Bullets: []*career.CVBullet{
								{
									ID:              "bullet-3",
									SectionID:       "section-experience",
									Text:            "Built real-time data processing system handling 1M events/day",
									SourceEventIDs:  []string{"event-3"},
									Rank:            0.85,
									Confidence:      0.80,
									InclusionReason: "outcome",
								},
							},
						},
					},
				},
				{
					ID:          "section-skills",
					CVViewID:    "test-cv",
					SectionType: "skills",
					Title:       "Skills",
					Order:       2,
					Content: []*career.SectionContentGroup{
						{
							Header: "Languages",
							Bullets: []*career.CVBullet{
								{
									ID:              "bullet-4",
									SectionID:       "section-skills",
									Text:            "Go, Python, JavaScript, SQL",
									SourceEventIDs:  []string{"event-1"},
									Rank:            0.8,
									Confidence:      0.90,
									InclusionReason: "activity",
								},
							},
						},
					},
				},
			},
		}
	}

	BeforeEach(func() {
		events := []*career.CareerEvent{
			{
				ID:         "event1",
				Text:       "Implemented feature X",
				Date:       time.Date(2025, 1, 1, 0, 0, 0, 0, time.UTC),
				Company:    "Company A",
				CreatedAt:  time.Date(2025, 1, 1, 0, 0, 0, 0, time.UTC),
				UpdatedAt:  time.Date(2025, 1, 1, 0, 0, 0, 0, time.UTC),
				Tags:       []string{"go", "backend"},
				Categories: []string{"technical"},
			},
		}

		profiles := []*CVProfile{
			{
				ID:             "profile1",
				Name:           "Senior IC",
				TargetRole:     "senior_ic",
				TargetAudience: "hiring_manager",
				Description:    "Profile for senior individual contributor roles",
			},
		}

		ctx = &GenerateCVContext{
			AvailableProfiles: profiles,
			Events:            events,
			Facts:             make([]*career.Fact, 0),
			DefaultProfile:    profiles[0],
		}

		var err error
		intent, err = NewGenerateCVIntent(ctx)
		Expect(err).NotTo(HaveOccurred())
		Expect(intent).NotTo(BeNil())
	})

	Describe("Standard Structure Preview", func() {
		BeforeEach(func() {
			intent.Init()
			intent.state.currentState = GenerateCVStatePreview
			intent.state.generatedCV = createTestCV()
			intent.state.selectedCVStructure = CVStructureStandard
		})

		It("should render CV name as header", func() {
			view := intent.viewPreviewStandard()
			Expect(view).To(ContainSubstring("Senior IC CV"))
		})

		It("should render target role and audience", func() {
			view := intent.viewPreviewStandard()
			Expect(view).To(ContainSubstring("senior_ic"))
			Expect(view).To(ContainSubstring("hiring_manager"))
		})

		It("should render all sections", func() {
			view := intent.viewPreviewStandard()
			Expect(view).To(ContainSubstring("SUMMARY"))
			Expect(view).To(ContainSubstring("EXPERIENCE"))
			Expect(view).To(ContainSubstring("SKILLS"))
		})

		It("should render summary section content", func() {
			view := intent.viewPreviewStandard()
			Expect(view).To(ContainSubstring("Experienced software engineer"))
		})

		It("should render experience with company headers and dates", func() {
			view := intent.viewPreviewStandard()
			Expect(view).To(ContainSubstring("TechCorp"))
			Expect(view).To(ContainSubstring("Jan 2020"))
			Expect(view).To(ContainSubstring("Present"))
			Expect(view).To(ContainSubstring("StartupCo"))
		})

		It("should render all bullets for each group", func() {
			view := intent.viewPreviewStandard()
			Expect(view).To(ContainSubstring("Led migration to microservices"))
			Expect(view).To(ContainSubstring("Implemented CI/CD pipeline"))
			Expect(view).To(ContainSubstring("Built real-time data processing"))
		})

		It("should render skills section", func() {
			view := intent.viewPreviewStandard()
			Expect(view).To(ContainSubstring("Languages"))
			Expect(view).To(ContainSubstring("Go, Python, JavaScript, SQL"))
		})
	})

	Describe("Narrative Structure Preview", func() {
		BeforeEach(func() {
			intent.Init()
			intent.state.currentState = GenerateCVStatePreview
			intent.state.generatedCV = createTestCV()
			intent.state.selectedCVStructure = CVStructureNarrative
		})

		It("should render profile header with name", func() {
			view := intent.viewPreviewNarrative()
			// Narrative uses DefaultNarrativeProfile() which has hardcoded name
			Expect(view).To(ContainSubstring("Yomi Colledge"))
		})

		It("should render role from profile", func() {
			view := intent.viewPreviewNarrative()
			Expect(view).To(ContainSubstring("Senior Software Engineer"))
		})

		It("should render location", func() {
			view := intent.viewPreviewNarrative()
			Expect(view).To(ContainSubstring("Remote"))
		})

		It("should render Summary section", func() {
			view := intent.viewPreviewNarrative()
			Expect(view).To(ContainSubstring("SUMMARY"))
		})

		It("should render Core Strengths section", func() {
			view := intent.viewPreviewNarrative()
			Expect(view).To(ContainSubstring("CORE STRENGTHS"))
		})

		It("should render Languages & Technologies section", func() {
			view := intent.viewPreviewNarrative()
			Expect(view).To(ContainSubstring("LANGUAGES & TECHNOLOGIES"))
		})

		It("should render Selected Experience section", func() {
			view := intent.viewPreviewNarrative()
			Expect(view).To(ContainSubstring("SELECTED EXPERIENCE"))
		})

		It("should render What I Bring section", func() {
			view := intent.viewPreviewNarrative()
			Expect(view).To(ContainSubstring("WHAT I BRING"))
		})

		It("should filter bullets by confidence >= 0.75 in Selected Experience", func() {
			view := intent.viewPreviewNarrative()
			// High confidence bullet should appear
			Expect(view).To(ContainSubstring("Led migration to microservices"))
			Expect(view).To(ContainSubstring("Built real-time data processing"))
			// Low confidence bullet (0.70) should NOT appear in narrative
			Expect(view).NotTo(ContainSubstring("Implemented CI/CD pipeline"))
		})

		It("should show company headers in Selected Experience", func() {
			view := intent.viewPreviewNarrative()
			Expect(view).To(ContainSubstring("TechCorp"))
			Expect(view).To(ContainSubstring("StartupCo"))
		})
	})

	Describe("Variant Selection Affects Preview", func() {
		BeforeEach(func() {
			intent.Init()
			// Navigate through the flow
			intent.Update(tea.KeyMsg{Type: tea.KeyEnter}) // profile -> audience
			intent.Update(tea.KeyMsg{Type: tea.KeyEnter}) // audience -> role emphasis
		})

		It("should use standard preview when Senior Backend variant selected", func() {
			// Senior Backend with Standard length uses standard base structure
			intent.Update(tea.KeyMsg{Type: tea.KeyEnter}) // role emphasis (Senior Backend) -> length format
			intent.Update(tea.KeyMsg{Type: tea.KeyEnter}) // length format (Standard) -> generating

			// Simulate generation complete
			intent.Update(CVGenerationCompleteMsg{
				CV:    createTestCV(),
				Error: nil,
			})

			Expect(intent.state.currentState).To(Equal(GenerateCVStatePreview))
			Expect(intent.state.selectedCVStructure).To(Equal(CVStructureStandard))

			// Use direct method to test content (viewport wrapper returns empty without initialization)
			view := intent.viewPreviewStandard()
			// Standard shows all bullets including low confidence
			Expect(view).To(ContainSubstring("Implemented CI/CD pipeline"))
		})

		It("should use narrative preview when Language-Agnostic variant selected", func() {
			// Navigate to Language-Agnostic (index 3)
			intent.Update(tea.KeyMsg{Type: tea.KeyDown}) // Staff/Principal
			intent.Update(tea.KeyMsg{Type: tea.KeyDown}) // Consulting
			intent.Update(tea.KeyMsg{Type: tea.KeyDown}) // Language-Agnostic
			Expect(intent.state.roleEmphasisIndex).To(Equal(3))
			intent.Update(tea.KeyMsg{Type: tea.KeyEnter}) // role emphasis -> length format
			intent.Update(tea.KeyMsg{Type: tea.KeyEnter}) // length format (Standard) -> generating

			// Simulate generation complete
			intent.Update(CVGenerationCompleteMsg{
				CV:    createTestCV(),
				Error: nil,
			})

			Expect(intent.state.currentState).To(Equal(GenerateCVStatePreview))
			Expect(intent.state.selectedCVStructure).To(Equal(CVStructure("narrative")))

			// Use direct method to test content (viewport wrapper returns empty without initialization)
			view := intent.viewPreviewNarrative()
			// Narrative filters out low confidence bullets
			Expect(view).NotTo(ContainSubstring("Implemented CI/CD pipeline"))
			// And shows narrative-specific sections
			Expect(view).To(ContainSubstring("CORE STRENGTHS"))
		})
	})

	Describe("Empty CV Handling", func() {
		It("should handle nil CV gracefully", func() {
			intent.Init()
			intent.state.currentState = GenerateCVStatePreview
			intent.state.generatedCV = nil

			view := intent.viewPreview()
			Expect(view).To(ContainSubstring("No CV generated"))
		})

		It("should handle CV with no sections in standard mode", func() {
			intent.Init()
			intent.state.currentState = GenerateCVStatePreview
			intent.state.selectedCVStructure = CVStructureStandard
			intent.state.generatedCV = &career.CVView{
				ID:               "empty-cv",
				Name:             "Empty CV",
				TargetRole:       "senior_ic",
				TargetAudience:   "hiring_manager",
				GeneratedAt:      time.Now(),
				SourceEventCount: 0,
				SourceFactCount:  0,
				Sections:         []*career.CVSection{},
			}

			// Use direct method to test content
			view := intent.viewPreviewStandard()
			Expect(view).To(ContainSubstring("Empty CV"))
			// Should not panic
		})

		It("should handle CV with no sections in narrative mode", func() {
			intent.Init()
			intent.state.currentState = GenerateCVStatePreview
			intent.state.selectedCVStructure = CVStructureNarrative
			intent.state.generatedCV = &career.CVView{
				ID:               "empty-cv",
				Name:             "Empty CV",
				TargetRole:       "senior_ic",
				TargetAudience:   "hiring_manager",
				GeneratedAt:      time.Now(),
				SourceEventCount: 0,
				SourceFactCount:  0,
				Sections:         []*career.CVSection{},
			}

			// Use direct method to test content
			view := intent.viewPreviewNarrative()
			// Should still show narrative structure with defaults
			Expect(view).To(ContainSubstring("CORE STRENGTHS"))
			Expect(view).To(ContainSubstring("WHAT I BRING"))
		})
	})
})

var _ = Describe("NarrativeProfileData", func() {
	Describe("DefaultNarrativeProfile", func() {
		It("should return a profile with name", func() {
			profile := DefaultNarrativeProfile()
			Expect(profile.Name).To(Equal("Yomi Colledge"))
		})

		It("should return a profile with role", func() {
			profile := DefaultNarrativeProfile()
			Expect(profile.Role).To(Equal("Senior Software Engineer / Technical Consultant"))
		})

		It("should return a profile with location", func() {
			profile := DefaultNarrativeProfile()
			Expect(profile.Location).To(Equal("Remote (UK)"))
		})

		It("should return a profile with email", func() {
			profile := DefaultNarrativeProfile()
			Expect(profile.Email).To(Equal("yomi@boodah.net"))
		})

		It("should return a profile with GitHub URL", func() {
			profile := DefaultNarrativeProfile()
			Expect(profile.GitHub).To(Equal("https://github.com/baphled"))
		})

		It("should return a profile with portfolio URL", func() {
			profile := DefaultNarrativeProfile()
			Expect(profile.Portfolio).To(Equal("http://boodah.net"))
		})

		It("should return default core strengths", func() {
			profile := DefaultNarrativeProfile()
			Expect(len(profile.CoreStrengths)).To(BeNumerically(">", 0))
			Expect(profile.CoreStrengths).To(ContainElement(ContainSubstring("backend")))
		})

		It("should return default technologies", func() {
			profile := DefaultNarrativeProfile()
			Expect(profile.Languages).To(ContainSubstring("Go"))
			Expect(profile.Languages).To(ContainSubstring("Ruby"))
		})

		It("should return default value propositions", func() {
			profile := DefaultNarrativeProfile()
			Expect(len(profile.ValuePropositions)).To(BeNumerically(">", 0))
		})
	})
})
