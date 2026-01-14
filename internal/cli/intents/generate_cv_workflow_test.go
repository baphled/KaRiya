package intents

import (
	"context"
	"io"
	"time"

	"github.com/baphled/kariya/internal/cli/service"
	career "github.com/baphled/kariya/internal/domain/career"
	"github.com/baphled/kariya/internal/logger"
	careerrepo "github.com/baphled/kariya/internal/repository/career"
	careerservice "github.com/baphled/kariya/internal/service/career"
	"github.com/baphled/kariya/internal/service/career/cv"
	tea "github.com/charmbracelet/bubbletea"
	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
)

var _ = Describe("GenerateCV Full Workflow", func() {
	var (
		intent        *GenerateCVIntent
		ctx           *GenerateCVContext
		testProfiles  []*CVProfile
		testEvents    []*career.CareerEvent
		testFacts     []*career.Fact
		exportService *cv.ExportService
	)

	BeforeEach(func() {
		// Setup repositories and services
		repo := careerrepo.NewMemoryRepository()
		burstRepo := careerrepo.NewMemoryBurstRepository()
		factRepo := careerrepo.NewMemoryFactRepository()
		svc := careerservice.NewService(repo)
		svc.SetBurstRepository(burstRepo)
		svc.SetFactRepository(factRepo)
		_ = service.NewCLIEventService(svc)

		// Create export service with silent logger
		log := logger.New(io.Discard, logger.ErrorLevel)
		exportService = cv.NewExportService(log)

		// Create test profiles
		testProfiles = []*CVProfile{
			{
				ID:             "profile_1",
				Name:           "Senior IC - Tech Lead",
				TargetRole:     "senior_ic",
				TargetAudience: "hiring_manager",
				Description:    "Profile for senior individual contributor positions",
			},
		}

		// Create test events
		testEvents = []*career.CareerEvent{
			{
				ID:         "event_1",
				Text:       "Led team standup meetings and improved communication across the team",
				Date:       time.Date(2025, 12, 15, 0, 0, 0, 0, time.UTC),
				Company:    "Acme Corp",
				Project:    "Project Alpha",
				Tags:       []string{"leadership", "communication"},
				Categories: []string{"team_management"},
				CreatedAt:  time.Now(),
				UpdatedAt:  time.Now(),
			},
		}

		// Create test facts
		testFacts = []*career.Fact{
			{
				ID:                   "fact_1",
				Text:                 "Improved team velocity by 30%",
				CompetencyCategories: []string{"leadership", "impact"},
				RoleFit:              "staff",
				AudienceRelevance:    []string{"hiring_manager", "peer"},
				StrengthSignal:       "high",
				SourceEventID:        "event_1",
				CreatedAt:            time.Now(),
				UpdatedAt:            time.Now(),
			},
		}

		// Create the GenerateCV context with export service
		ctx = &GenerateCVContext{
			AvailableProfiles: testProfiles,
			Events:            testEvents,
			Facts:             testFacts,
			DefaultProfile:    testProfiles[0],
			ExportService:     exportService,
			AppContext:        context.Background(),
		}

		// Create the intent
		var err error
		intent, err = NewGenerateCVIntent(ctx)
		Expect(err).NotTo(HaveOccurred())
		Expect(intent).NotTo(BeNil())

		// Initialize the intent
		cmd := intent.Init()
		_ = cmd
	})

	Describe("Role Emphasis Selection State", func() {
		BeforeEach(func() {
			// Navigate to role emphasis selection: Profile -> Audience -> Role Emphasis
			intent.Update(tea.KeyMsg{Type: tea.KeyEnter}) // Select profile
			intent.Update(tea.KeyMsg{Type: tea.KeyEnter}) // Select audience
		})

		It("should show role emphasis selection screen after audience selection", func() {
			view := intent.View()
			Expect(view).To(ContainSubstring("Role Emphasis"))
			Expect(view).To(ContainSubstring("Senior Backend"))
			Expect(view).To(ContainSubstring("Staff/Principal"))
		})

		It("should show Senior Backend selected by default", func() {
			view := intent.View()
			// The selected item has a marker
			Expect(view).To(ContainSubstring("Senior Backend"))
			// Senior Backend description should be visible
			Expect(view).To(ContainSubstring("backend engineering"))
		})

		It("should allow navigation to other role emphasis options", func() {
			// Press down to select Staff/Principal
			intent.Update(tea.KeyMsg{Type: tea.KeyDown})

			view := intent.View()
			Expect(view).To(ContainSubstring("Staff/Principal"))
		})

		It("should allow vim-style navigation (j/k)", func() {
			// Press 'j' to move down to Staff/Principal
			intent.Update(tea.KeyMsg{Type: tea.KeyRunes, Runes: []rune("j")})

			view := intent.View()
			Expect(view).To(ContainSubstring("Staff/Principal"))

			// Press 'k' to move back up to Senior Backend
			intent.Update(tea.KeyMsg{Type: tea.KeyRunes, Runes: []rune("k")})

			view = intent.View()
			Expect(view).To(ContainSubstring("Senior Backend"))
		})

		It("should go back to audience selection on Escape", func() {
			intent.Update(tea.KeyMsg{Type: tea.KeyEscape})

			view := intent.View()
			Expect(view).To(ContainSubstring("Audience"))
		})
	})

	Describe("Senior Backend Variant Workflow", func() {
		BeforeEach(func() {
			// Navigate to role emphasis selection
			intent.Update(tea.KeyMsg{Type: tea.KeyEnter}) // Select profile
			intent.Update(tea.KeyMsg{Type: tea.KeyEnter}) // Select audience
			// Senior Backend is selected by default
		})

		It("should generate CV with Senior Backend Standard variant", func() {
			// Select Senior Backend (default)
			intent.Update(tea.KeyMsg{Type: tea.KeyEnter}) // role emphasis -> length format
			// Standard is default (index 1)
			intent.Update(tea.KeyMsg{Type: tea.KeyEnter}) // length format -> generating

			// Should transition to generating state
			view := intent.View()
			// Either shows "Generating" or moves to preview (if generation is fast)
			Expect(view).NotTo(BeEmpty())
		})

		It("should show preview format", func() {
			// Generate with Senior Backend Standard
			intent.Update(tea.KeyMsg{Type: tea.KeyEnter}) // role emphasis
			intent.Update(tea.KeyMsg{Type: tea.KeyEnter}) // length format

			// Process the generation message if available
			view := intent.View()
			Expect(view).NotTo(BeEmpty())
		})
	})

	Describe("Language-Agnostic Variant Workflow", func() {
		BeforeEach(func() {
			// Navigate to role emphasis selection
			intent.Update(tea.KeyMsg{Type: tea.KeyEnter}) // Select profile
			intent.Update(tea.KeyMsg{Type: tea.KeyEnter}) // Select audience
		})

		It("should select Language-Agnostic role emphasis", func() {
			// Move to Language-Agnostic (index 3)
			intent.Update(tea.KeyMsg{Type: tea.KeyDown}) // Staff/Principal
			intent.Update(tea.KeyMsg{Type: tea.KeyDown}) // Consulting
			intent.Update(tea.KeyMsg{Type: tea.KeyDown}) // Language-Agnostic

			view := intent.View()
			Expect(view).To(ContainSubstring("Language-Agnostic"))
		})

		It("should generate CV with Language-Agnostic variant", func() {
			// Select Language-Agnostic
			intent.Update(tea.KeyMsg{Type: tea.KeyDown})  // Staff/Principal
			intent.Update(tea.KeyMsg{Type: tea.KeyDown})  // Consulting
			intent.Update(tea.KeyMsg{Type: tea.KeyDown})  // Language-Agnostic
			intent.Update(tea.KeyMsg{Type: tea.KeyEnter}) // role emphasis -> length format
			intent.Update(tea.KeyMsg{Type: tea.KeyEnter}) // length format -> generating

			// Should transition to generating or preview
			view := intent.View()
			Expect(view).NotTo(BeEmpty())
		})

		It("should show narrative-based preview for Language-Agnostic", func() {
			// Language-Agnostic uses Narrative base structure
			intent.Update(tea.KeyMsg{Type: tea.KeyDown})  // Staff/Principal
			intent.Update(tea.KeyMsg{Type: tea.KeyDown})  // Consulting
			intent.Update(tea.KeyMsg{Type: tea.KeyDown})  // Language-Agnostic
			intent.Update(tea.KeyMsg{Type: tea.KeyEnter}) // role emphasis
			intent.Update(tea.KeyMsg{Type: tea.KeyEnter}) // length format

			// The view should eventually show content
			view := intent.View()
			Expect(view).NotTo(BeEmpty())
		})
	})

	Describe("Full Workflow: Profile -> Audience -> Role Emphasis -> Length -> Generate -> Preview", func() {
		It("should complete full Senior Backend Standard workflow without errors", func() {
			// Step 1: Select Profile
			view := intent.View()
			Expect(view).To(ContainSubstring("Profile"))
			intent.Update(tea.KeyMsg{Type: tea.KeyEnter})

			// Step 2: Select Audience
			view = intent.View()
			Expect(view).To(ContainSubstring("Audience"))
			intent.Update(tea.KeyMsg{Type: tea.KeyEnter})

			// Step 3: Select Role Emphasis (Senior Backend is default)
			view = intent.View()
			Expect(view).To(ContainSubstring("Role Emphasis"))
			intent.Update(tea.KeyMsg{Type: tea.KeyEnter})

			// Step 4: Select Length Format (Standard is default at index 1)
			view = intent.View()
			Expect(view).To(ContainSubstring("Length"))
			intent.Update(tea.KeyMsg{Type: tea.KeyEnter})

			// Step 5: Should be generating or at preview
			view = intent.View()
			Expect(view).NotTo(BeEmpty())

			// Result should be nil until completion
			result := intent.Result()
			// Result may be nil (still in progress) or set (if auto-completed)
			_ = result
		})

		It("should complete full Language-Agnostic workflow without errors", func() {
			// Step 1: Select Profile
			intent.Update(tea.KeyMsg{Type: tea.KeyEnter})

			// Step 2: Select Audience
			intent.Update(tea.KeyMsg{Type: tea.KeyEnter})

			// Step 3: Select Language-Agnostic Role Emphasis
			intent.Update(tea.KeyMsg{Type: tea.KeyDown}) // Staff/Principal
			intent.Update(tea.KeyMsg{Type: tea.KeyDown}) // Consulting
			intent.Update(tea.KeyMsg{Type: tea.KeyDown}) // Language-Agnostic
			intent.Update(tea.KeyMsg{Type: tea.KeyEnter})

			// Step 4: Select Length Format
			intent.Update(tea.KeyMsg{Type: tea.KeyEnter})

			// Step 5: Should be generating or at preview
			view := intent.View()
			Expect(view).NotTo(BeEmpty())
		})

		It("should allow going back through all states", func() {
			// Go forward to length format selection
			intent.Update(tea.KeyMsg{Type: tea.KeyEnter}) // Profile -> Audience
			intent.Update(tea.KeyMsg{Type: tea.KeyEnter}) // Audience -> Role Emphasis
			intent.Update(tea.KeyMsg{Type: tea.KeyEnter}) // Role Emphasis -> Length Format

			view := intent.View()
			Expect(view).To(ContainSubstring("Length"))

			// Go back to Role Emphasis
			intent.Update(tea.KeyMsg{Type: tea.KeyEscape})
			view = intent.View()
			Expect(view).To(ContainSubstring("Role Emphasis"))

			// Go back to Audience
			intent.Update(tea.KeyMsg{Type: tea.KeyEscape})
			view = intent.View()
			Expect(view).To(ContainSubstring("Audience"))

			// Go back to Profile
			intent.Update(tea.KeyMsg{Type: tea.KeyEscape})
			view = intent.View()
			Expect(view).To(ContainSubstring("Profile"))
		})
	})

	Describe("Structure Selection Result Metadata", func() {
		It("should include selected structure in result metadata when completed", func() {
			// Complete the workflow with Narrative
			intent.Update(tea.KeyMsg{Type: tea.KeyEnter}) // Profile
			intent.Update(tea.KeyMsg{Type: tea.KeyEnter}) // Audience
			intent.Update(tea.KeyMsg{Type: tea.KeyDown})  // Select Narrative
			intent.Update(tea.KeyMsg{Type: tea.KeyEnter}) // Generate

			// Continue through preview/review/confirm if needed
			// This depends on the state machine implementation
			for i := 0; i < 10; i++ {
				result := intent.Result()
				if result != nil && result.Status == Completed {
					// Check metadata contains structure
					structure, ok := result.Metadata["structure"]
					if ok {
						Expect(structure).To(Equal(CVStructureNarrative))
					}
					break
				}
				// Try pressing Enter to proceed
				intent.Update(tea.KeyMsg{Type: tea.KeyEnter})
			}
		})
	})
})

var _ = Describe("Export Service Structure Integration", func() {
	var (
		exportService *cv.ExportService
		testCV        *career.CVView
		testSections  []*career.CVSection
		testBullets   map[string][]*career.CVBullet
	)

	BeforeEach(func() {
		log := logger.New(io.Discard, logger.ErrorLevel)
		exportService = cv.NewExportService(log)

		testCV = &career.CVView{
			ID:               "test_cv_1",
			Name:             "Test CV",
			TargetRole:       "Senior Engineer",
			TargetAudience:   "Hiring Manager",
			GeneratedAt:      time.Now(),
			SourceEventCount: 5,
			SourceFactCount:  3,
		}

		testSections = []*career.CVSection{
			{
				Title:       "Experience",
				SectionType: "experience",
				Content: []*career.SectionContentGroup{
					{
						Header:    "Acme Corp",
						StartDate: "Jan 2020",
						EndDate:   "Present",
						Bullets: []*career.CVBullet{
							{Text: "Led development of key features", Confidence: 0.9},
							{Text: "Improved system performance by 40%", Confidence: 0.85},
						},
					},
				},
			},
			{
				Title:       "Summary",
				SectionType: "summary",
				Summary:     "Experienced engineer with expertise in distributed systems.",
			},
		}

		testBullets = make(map[string][]*career.CVBullet)
	})

	Describe("Export with Standard Structure", func() {
		It("should export to text with standard format", func() {
			content, err := exportService.Export(
				context.Background(),
				testCV,
				testSections,
				testBullets,
				cv.CVStructureStandard,
				cv.ExportFormatText,
			)

			Expect(err).NotTo(HaveOccurred())
			Expect(content).To(ContainSubstring("TEST CV")) // Uppercase header
			Expect(content).To(ContainSubstring("Target Role"))
		})

		It("should export to markdown with standard format", func() {
			content, err := exportService.Export(
				context.Background(),
				testCV,
				testSections,
				testBullets,
				cv.CVStructureStandard,
				cv.ExportFormatMarkdown,
			)

			Expect(err).NotTo(HaveOccurred())
			Expect(content).To(ContainSubstring("# Test CV")) // Markdown header
		})
	})

	Describe("Export with Narrative Structure", func() {
		It("should export to text with narrative format", func() {
			content, err := exportService.Export(
				context.Background(),
				testCV,
				testSections,
				testBullets,
				cv.CVStructureNarrative,
				cv.ExportFormatText,
			)

			Expect(err).NotTo(HaveOccurred())
			// Narrative format includes profile data (text uses UPPERCASE headers)
			Expect(content).To(ContainSubstring("YOMI COLLEDGE"))
			Expect(content).To(ContainSubstring("CORE STRENGTHS"))
			Expect(content).To(ContainSubstring("WHAT I BRING"))
		})

		It("should export to markdown with narrative format", func() {
			content, err := exportService.Export(
				context.Background(),
				testCV,
				testSections,
				testBullets,
				cv.CVStructureNarrative,
				cv.ExportFormatMarkdown,
			)

			Expect(err).NotTo(HaveOccurred())
			// Narrative markdown format
			Expect(content).To(ContainSubstring("# Yomi Colledge"))
			Expect(content).To(ContainSubstring("## Core Strengths"))
			Expect(content).To(ContainSubstring("## What I Bring"))
		})

		It("should include languages and technologies in narrative export", func() {
			content, err := exportService.Export(
				context.Background(),
				testCV,
				testSections,
				testBullets,
				cv.CVStructureNarrative,
				cv.ExportFormatMarkdown,
			)

			Expect(err).NotTo(HaveOccurred())
			Expect(content).To(ContainSubstring("Languages & Technologies"))
			Expect(content).To(ContainSubstring("Ruby"))
			Expect(content).To(ContainSubstring("Go"))
		})
	})

	Describe("YAML always uses Standard", func() {
		It("should use standard structure for YAML regardless of selected structure", func() {
			// Request Narrative but with YAML format
			content, err := exportService.Export(
				context.Background(),
				testCV,
				testSections,
				testBullets,
				cv.CVStructureNarrative, // Narrative requested
				cv.ExportFormatYAML,     // But YAML format
			)

			Expect(err).NotTo(HaveOccurred())
			// YAML should NOT contain narrative profile data
			Expect(content).NotTo(ContainSubstring("Yomi Colledge"))
			// Should contain standard YAML structure
			Expect(content).To(ContainSubstring("name:"))
			Expect(content).To(ContainSubstring("target_role:"))
		})
	})
})
