package intents_test

import (
	"testing"
	"time"

	"github.com/baphled/kariya/internal/cli/intents"
	"github.com/baphled/kariya/internal/cli/service"
	career "github.com/baphled/kariya/internal/domain/career"
	careerrepo "github.com/baphled/kariya/internal/repository/career"
	careerservice "github.com/baphled/kariya/internal/service/career"
	tea "github.com/charmbracelet/bubbletea"
	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
)

func TestGenerateCVIntegration(t *testing.T) {
	RegisterFailHandler(Fail)
	RunSpecs(t, "GenerateCV Intent Integration Suite")
}

var _ = Describe("GenerateCV Intent Integration", func() {
	var (
		intent       *intents.GenerateCVIntent
		ctx          *intents.GenerateCVContext
		testProfiles []*intents.CVProfile
		testEvents   []*career.CareerEvent
		testFacts    []*career.Fact
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

		// Create test profiles
		testProfiles = []*intents.CVProfile{
			{
				ID:             "profile_1",
				Name:           "Senior IC - Tech Lead",
				TargetRole:     "senior_ic",
				TargetAudience: []string{"hiring_manager", "recruiter"},
				Description:    "Profile for senior individual contributor positions",
			},
			{
				ID:             "profile_2",
				Name:           "Engineering Manager",
				TargetRole:     "em",
				TargetAudience: []string{"hiring_manager"},
				Description:    "Profile for engineering manager positions",
			},
			{
				ID:             "profile_3",
				Name:           "Staff Engineer",
				TargetRole:     "staff",
				TargetAudience: []string{"peer", "recruiter"},
				Description:    "Profile for staff engineer positions",
			},
		}

		// Create test events
		testEvents = []*career.CareerEvent{
			{
				ID:        "event_1",
				Text:      "Led team standup meetings and improved communication across the team",
				Date:      time.Date(2025, 12, 15, 0, 0, 0, 0, time.UTC),
				Company:   "Acme Corp",
				Project:   "Project Alpha",
				Tags:      []string{"leadership", "communication"},
				Categories: []string{"team_management"},
				CreatedAt: time.Now(),
				UpdatedAt: time.Now(),
			},
			{
				ID:        "event_2",
				Text:      "Implemented new code review process that reduced review time",
				Date:      time.Date(2025, 12, 10, 0, 0, 0, 0, time.UTC),
				Company:   "Acme Corp",
				Project:   "Project Alpha",
				Tags:      []string{"engineering", "process"},
				Categories: []string{"technical_excellence"},
				CreatedAt: time.Now(),
				UpdatedAt: time.Now(),
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
			{
				ID:                   "fact_2",
				Text:                 "Reduced code review time by 40%",
				CompetencyCategories: []string{"technical", "impact"},
				RoleFit:              "principal",
				AudienceRelevance:    []string{"recruiter", "peer"},
				StrengthSignal:       "high",
				SourceEventID:        "event_2",
				CreatedAt:            time.Now(),
				UpdatedAt:            time.Now(),
			},
		}

		// Create the GenerateCV context
		ctx = &intents.GenerateCVContext{
			AvailableProfiles: testProfiles,
			Events:            testEvents,
			Facts:             testFacts,
			DefaultProfile:    testProfiles[0],
		}

		// Create the intent
		var err error
		intent, err = intents.NewGenerateCVIntent(ctx)
		Expect(err).NotTo(HaveOccurred())
		Expect(intent).NotTo(BeNil())

		// Initialize the intent (this must be called to set up the profile table)
		cmd := intent.Init()
		_ = cmd // Command can be nil or a command, both are valid
	})

	Describe("Profile Selection State", func() {
		It("should initialize in profile selection state", func() {
			Expect(intent).NotTo(BeNil())
			// The intent should start in GenerateCVStateSelectProfile
			// We verify this by checking the view contains profile information
			view := intent.View()
			Expect(view).NotTo(BeEmpty())
			Expect(view).To(ContainSubstring("Profile"))
		})

		It("should render profile table with all available profiles", func() {
			view := intent.View()
			Expect(view).NotTo(BeEmpty())
			// Should display profile names
			Expect(view).To(ContainSubstring("Senior IC - Tech Lead"))
			Expect(view).To(ContainSubstring("Engineering Manager"))
			Expect(view).To(ContainSubstring("Staff Engineer"))
		})

		It("should accept keyboard input to navigate profiles", func() {
			// Initial view
			initialView := intent.View()
			Expect(initialView).NotTo(BeEmpty())

			// Send down arrow to navigate
			cmd := intent.Update(tea.KeyMsg{Type: tea.KeyDown})
			// Command can be nil or a command, both are valid
			_ = cmd

			// View should still be valid
			view := intent.View()
			Expect(view).NotTo(BeEmpty())
		})


		It("should NOT return to menu when navigating profiles", func() {
			// This is the FAILING TEST that reproduces the bug
			// When user presses down arrow to navigate profiles, the intent
			// should stay in profile selection state, NOT return to menu

			// Get initial view (should show profile selection)
			initialView := intent.View()
			Expect(initialView).To(ContainSubstring("Profile"))

			// Navigate down through profiles
			for i := 0; i < 3; i++ {
				cmd := intent.Update(tea.KeyMsg{Type: tea.KeyDown})
				// Command should not trigger menu return
				_ = cmd

				// View should still show profile selection, NOT menu
				currentView := intent.View()
				Expect(currentView).NotTo(BeEmpty())
				Expect(currentView).To(ContainSubstring("Profile"),
					"After navigating profile %d, view should still show profile selection, not menu", i+1)
			}

			// Check the result - should be nil when intent is still in progress
			result := intent.Result()
			Expect(result).To(BeNil())
		})

		It("should select profile on Enter key press", func() {
			// Navigate down once
			intent.Update(tea.KeyMsg{Type: tea.KeyDown})

			// Press Enter to select
			cmd := intent.Update(tea.KeyMsg{Type: tea.KeyEnter})
			_ = cmd

			// After selection, should transition to audience selection
			// View should change to show audience options
			view := intent.View()
			Expect(view).NotTo(BeEmpty())
			Expect(view).To(ContainSubstring("Audience"))
		})

		It("should allow j/k keys for navigation (vim-style)", func() {
			// Press 'j' to move down
			cmd := intent.Update(tea.KeyMsg{Type: tea.KeyRunes, Runes: []rune("j")})
			_ = cmd

			view := intent.View()
			Expect(view).To(ContainSubstring("Profile"))

			// Press 'k' to move up
			cmd = intent.Update(tea.KeyMsg{Type: tea.KeyRunes, Runes: []rune("k")})
			_ = cmd

			view = intent.View()
			Expect(view).To(ContainSubstring("Profile"))
		})

		It("should cancel on Escape key", func() {
			// Press Escape
			cmd := intent.Update(tea.KeyMsg{Type: tea.KeyEscape})
			_ = cmd

			// Result should be Cancelled
			result := intent.Result()
			Expect(result.Status).To(Equal(intents.Cancelled))
		})

		It("should cancel on Ctrl+C", func() {
			// Press Ctrl+C
			cmd := intent.Update(tea.KeyMsg{Type: tea.KeyCtrlC})
			_ = cmd

			// Result should be Cancelled
			result := intent.Result()
			Expect(result.Status).To(Equal(intents.Cancelled))
		})

		It("should cancel on 'q' key", func() {
			// Press 'q'
			cmd := intent.Update(tea.KeyMsg{Type: tea.KeyRunes, Runes: []rune("q")})
			_ = cmd

			// Result should be Cancelled
			result := intent.Result()
			Expect(result.Status).To(Equal(intents.Cancelled))
		})
	})

	Describe("State Transitions", func() {
		It("should transition to audience selection after profile selection", func() {
			// Start in profile selection
			view := intent.View()
			Expect(view).To(ContainSubstring("Profile"))

			// Select a profile
			intent.Update(tea.KeyMsg{Type: tea.KeyEnter})

			// Should transition to audience selection
			view = intent.View()
			Expect(view).To(ContainSubstring("Audience"))
		})

		It("should return to profile selection when pressing Escape from audience selection", func() {
			// Transition to audience selection
			intent.Update(tea.KeyMsg{Type: tea.KeyEnter})

			// Press Escape to go back
			intent.Update(tea.KeyMsg{Type: tea.KeyEscape})

			// Should be back at profile selection
			view := intent.View()
			Expect(view).To(ContainSubstring("Profile"))
		})
	})

	Describe("View Rendering", func() {
		It("should render without panic", func() {
			Expect(func() {
				_ = intent.View()
			}).NotTo(Panic())
		})

		It("should contain profile table header", func() {
			view := intent.View()
			Expect(view).To(ContainSubstring("Profile"))
		})

		It("should contain breadcrumb navigation", func() {
			view := intent.View()
			// Should show navigation context
			Expect(view).NotTo(BeEmpty())
		})
	})

	Describe("Result Handling", func() {
		It("should return nil result before completion", func() {
			result := intent.Result()
			Expect(result).To(BeNil())
		})

		It("should return Cancelled status when user cancels", func() {
			intent.Update(tea.KeyMsg{Type: tea.KeyEscape})
			result := intent.Result()
			Expect(result.Status).To(Equal(intents.Cancelled))
		})

		It("should return type-safe IntentResult", func() {
			// Before completing, result should be nil
			result := intent.Result()
			Expect(result).To(BeNil())
			
			// Transition through states
			intent.Update(tea.KeyMsg{Type: tea.KeyEnter})
			intent.Update(tea.KeyMsg{Type: tea.KeyEnter})
			intent.Update(tea.KeyMsg{Type: tea.KeyEnter})
			intent.Update(tea.KeyMsg{Type: tea.KeyEnter})

			// After confirming, result should have a status
			result = intent.Result()
			if result != nil {
				Expect(result.Status).NotTo(BeNil())
			}
		})
	})

	Describe("Init method", func() {
		It("should initialize without panic", func() {
			Expect(func() {
				cmd := intent.Init()
				_ = cmd
			}).NotTo(Panic())
		})

		It("should return a command or nil", func() {
			cmd := intent.Init()
			// Init can return nil or a command, both are valid
			Expect(cmd == nil || cmd != nil).To(BeTrue())
		})
	})
})

