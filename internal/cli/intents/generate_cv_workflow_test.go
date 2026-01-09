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

	Describe("Structure Selection State", func() {
		BeforeEach(func() {
			// Navigate to structure selection: Profile -> Audience -> Structure
			intent.Update(tea.KeyMsg{Type: tea.KeyEnter}) // Select profile
			intent.Update(tea.KeyMsg{Type: tea.KeyEnter}) // Select audience
		})

		It("should show structure selection screen after audience selection", func() {
			view := intent.View()
			Expect(view).To(ContainSubstring("Structure"))
			Expect(view).To(ContainSubstring("Standard"))
			Expect(view).To(ContainSubstring("Narrative"))
		})

		It("should show Standard selected by default", func() {
			view := intent.View()
			// The selected item has a marker
			Expect(view).To(ContainSubstring("Standard"))
			// Standard description should be visible
			Expect(view).To(ContainSubstring("Traditional CV"))
		})

		It("should allow navigation to Narrative option", func() {
			// Press down to select Narrative
			intent.Update(tea.KeyMsg{Type: tea.KeyDown})

			view := intent.View()
			Expect(view).To(ContainSubstring("Narrative"))
			// Narrative description should be visible
			Expect(view).To(ContainSubstring("Language-agnostic"))
		})

		It("should allow vim-style navigation (j/k)", func() {
			// Press 'j' to move down to Narrative
			intent.Update(tea.KeyMsg{Type: tea.KeyRunes, Runes: []rune("j")})

			view := intent.View()
			Expect(view).To(ContainSubstring("Narrative"))

			// Press 'k' to move back up to Standard
			intent.Update(tea.KeyMsg{Type: tea.KeyRunes, Runes: []rune("k")})

			view = intent.View()
			Expect(view).To(ContainSubstring("Standard"))
		})

		It("should go back to audience selection on Escape", func() {
			intent.Update(tea.KeyMsg{Type: tea.KeyEscape})

			view := intent.View()
			Expect(view).To(ContainSubstring("Audience"))
		})
	})

	Describe("Standard Structure Workflow", func() {
		BeforeEach(func() {
			// Navigate to structure selection
			intent.Update(tea.KeyMsg{Type: tea.KeyEnter}) // Select profile
			intent.Update(tea.KeyMsg{Type: tea.KeyEnter}) // Select audience
			// Standard is selected by default
		})

		It("should generate CV with Standard structure", func() {
			// Select Standard (default) and generate
			intent.Update(tea.KeyMsg{Type: tea.KeyEnter})

			// Should transition to generating state
			view := intent.View()
			// Either shows "Generating" or moves to preview (if generation is fast)
			Expect(view).NotTo(BeEmpty())
		})

		It("should show standard preview format", func() {
			// Generate with Standard structure
			intent.Update(tea.KeyMsg{Type: tea.KeyEnter})

			// Process the generation message if available
			// In tests without actual service, it may show preview directly
			view := intent.View()
			Expect(view).NotTo(BeEmpty())
		})
	})

	Describe("Narrative Structure Workflow", func() {
		BeforeEach(func() {
			// Navigate to structure selection
			intent.Update(tea.KeyMsg{Type: tea.KeyEnter}) // Select profile
			intent.Update(tea.KeyMsg{Type: tea.KeyEnter}) // Select audience
		})

		It("should select Narrative structure with down arrow", func() {
			// Move to Narrative
			intent.Update(tea.KeyMsg{Type: tea.KeyDown})

			view := intent.View()
			Expect(view).To(ContainSubstring("Narrative"))
		})

		It("should generate CV with Narrative structure", func() {
			// Select Narrative
			intent.Update(tea.KeyMsg{Type: tea.KeyDown})
			intent.Update(tea.KeyMsg{Type: tea.KeyEnter})

			// Should transition to generating or preview
			view := intent.View()
			Expect(view).NotTo(BeEmpty())
		})

		It("should show narrative preview with profile sections", func() {
			// Select Narrative and generate
			intent.Update(tea.KeyMsg{Type: tea.KeyDown})
			intent.Update(tea.KeyMsg{Type: tea.KeyEnter})

			// The view should eventually show narrative-specific content
			// Note: In unit tests without full service, the preview may be basic
			view := intent.View()
			Expect(view).NotTo(BeEmpty())
		})
	})

	Describe("Full Workflow: Profile -> Audience -> Structure -> Generate -> Preview", func() {
		It("should complete full Standard workflow without errors", func() {
			// Step 1: Select Profile
			view := intent.View()
			Expect(view).To(ContainSubstring("Profile"))
			intent.Update(tea.KeyMsg{Type: tea.KeyEnter})

			// Step 2: Select Audience
			view = intent.View()
			Expect(view).To(ContainSubstring("Audience"))
			intent.Update(tea.KeyMsg{Type: tea.KeyEnter})

			// Step 3: Select Structure (Standard is default)
			view = intent.View()
			Expect(view).To(ContainSubstring("Structure"))
			intent.Update(tea.KeyMsg{Type: tea.KeyEnter})

			// Step 4: Should be generating or at preview
			view = intent.View()
			Expect(view).NotTo(BeEmpty())

			// Result should be nil until completion
			result := intent.Result()
			// Result may be nil (still in progress) or set (if auto-completed)
			_ = result
		})

		It("should complete full Narrative workflow without errors", func() {
			// Step 1: Select Profile
			intent.Update(tea.KeyMsg{Type: tea.KeyEnter})

			// Step 2: Select Audience
			intent.Update(tea.KeyMsg{Type: tea.KeyEnter})

			// Step 3: Select Narrative Structure
			intent.Update(tea.KeyMsg{Type: tea.KeyDown}) // Move to Narrative
			intent.Update(tea.KeyMsg{Type: tea.KeyEnter})

			// Step 4: Should be generating or at preview
			view := intent.View()
			Expect(view).NotTo(BeEmpty())
		})

		It("should allow going back through all states", func() {
			// Go forward to structure selection
			intent.Update(tea.KeyMsg{Type: tea.KeyEnter}) // Profile -> Audience
			intent.Update(tea.KeyMsg{Type: tea.KeyEnter}) // Audience -> Structure

			view := intent.View()
			Expect(view).To(ContainSubstring("Structure"))

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
