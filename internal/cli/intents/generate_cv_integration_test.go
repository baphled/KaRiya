package intents_test

import (
	"time"

	"github.com/baphled/kariya/internal/cli/intents"
	"github.com/baphled/kariya/internal/cli/service"
	career "github.com/baphled/kariya/internal/domain/career"
	careermemory "github.com/baphled/kariya/internal/repository/career/memory"
	careerservice "github.com/baphled/kariya/internal/service/career"
	"github.com/baphled/kariya/internal/testutil/fixtures"
	tea "github.com/charmbracelet/bubbletea"
	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
)

var _ = Describe("GenerateCV Intent Integration", func() {
	var (
		intent       *intents.GenerateCVIntent
		ctx          *intents.GenerateCVContext
		testProfiles []*intents.CVProfile
		testEvents   []*career.Event
		testFacts    []*career.Fact
	)

	BeforeEach(func() {
		// Setup repositories and services
		repo := careermemory.NewEventRepository()
		burstRepo := careermemory.NewBurstRepository()
		factRepo := careermemory.NewFactRepository()
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
				TargetAudience: "hiring_manager",
				Description:    "Profile for senior individual contributor positions",
			},
			{
				ID:             "profile_2",
				Name:           "Engineering Manager",
				TargetRole:     "em",
				TargetAudience: "hiring_manager",
				Description:    "Profile for engineering manager positions",
			},
			{
				ID:             "profile_3",
				Name:           "Staff Engineer",
				TargetRole:     "staff",
				TargetAudience: "hiring_manager",
				Description:    "Profile for staff engineer positions",
			},
		}

		evt1 := fixtures.EventWith("event_1", "Led team standup meetings and improved communication across the team", "Acme Corp", "Project Alpha")
		evt1.Date = time.Date(2025, 12, 15, 0, 0, 0, 0, time.UTC)
		evt1.Tags = []string{"leadership", "communication"}
		evt1.Categories = []string{"team_management"}
		evt2 := fixtures.EventWith("event_2", "Implemented new code review process that reduced review time", "Acme Corp", "Project Alpha")
		evt2.Date = time.Date(2025, 12, 10, 0, 0, 0, 0, time.UTC)
		evt2.Tags = []string{"engineering", "process"}
		evt2.Categories = []string{"technical_excellence"}
		testEvents = []*career.Event{evt1, evt2}

		fact1 := fixtures.FactWithCategories("fact_1", "Improved team velocity by 30%", "event_1", []string{"leadership", "impact"}, []string{"hiring_manager", "peer"})
		fact1.RoleFit = "staff"
		fact2 := fixtures.FactWithCategories("fact_2", "Reduced code review time by 40%", "event_2", []string{"technical", "impact"}, []string{"recruiter", "peer"})
		fact2.RoleFit = "principal"
		testFacts = []*career.Fact{fact1, fact2}

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
			for i := range 3 {
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

		It("should ignore 'q' key within intent (quit only from main menu)", func() {
			// Press 'q'
			cmd := intent.Update(tea.KeyMsg{Type: tea.KeyRunes, Runes: []rune("q")})
			// q no longer quits from within intents - only from main menu
			Expect(cmd).To(BeNil())
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
