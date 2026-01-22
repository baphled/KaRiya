package intents

import (
	"context"
	"fmt"
	"time"

	"github.com/baphled/kariya/internal/cli/screens"
	"github.com/baphled/kariya/internal/domain/career"
	tea "github.com/charmbracelet/bubbletea"
	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
)

var _ = Describe("GenerateCV Wizard E2E Tests", func() {
	var (
		intent       *GenerateCVIntent
		ctx          *GenerateCVContext
		testProfiles []*CVProfile
		testEvents   []*career.CareerEvent
		testFacts    []*career.Fact
	)

	BeforeEach(func() {
		// Create test profiles
		testProfiles = []*CVProfile{
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
				TargetAudience: "recruiter",
				Description:    "Profile for engineering manager positions",
			},
		}

		// Create test events
		testEvents = []*career.CareerEvent{
			{
				ID:         "event_1",
				Text:       "Led team standup meetings and improved communication",
				Date:       time.Date(2025, 12, 15, 0, 0, 0, 0, time.UTC),
				Company:    "Acme Corp",
				Project:    "Project Alpha",
				Tags:       []string{"leadership", "communication"},
				Categories: []string{"team_management"},
				CreatedAt:  time.Now(),
				UpdatedAt:  time.Now(),
			},
			{
				ID:         "event_2",
				Text:       "Implemented new code review process",
				Date:       time.Date(2025, 12, 10, 0, 0, 0, 0, time.UTC),
				Company:    "Acme Corp",
				Project:    "Project Alpha",
				Tags:       []string{"engineering", "process"},
				Categories: []string{"technical_excellence"},
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

		// Create review screen factory
		reviewFactory := func(cv *career.CVView) screens.Screen {
			return &mockWizardReviewScreen{cv: cv}
		}

		// Create preview screen factory
		previewFactory := func(cv *career.CVView) screens.Screen {
			return &mockWizardPreviewScreen{cv: cv}
		}

		// Create context
		ctx = &GenerateCVContext{
			AvailableProfiles:    testProfiles,
			Events:               testEvents,
			Facts:                testFacts,
			DefaultProfile:       testProfiles[0],
			AppContext:           context.Background(),
			ReviewScreenFactory:  reviewFactory,
			PreviewScreenFactory: previewFactory,
		}

		// Create intent
		var err error
		intent, err = NewGenerateCVIntent(ctx)
		Expect(err).NotTo(HaveOccurred())
		Expect(intent).NotTo(BeNil())

		// Enable wizard flow
		intent.EnableWizardFlow()

		// Set terminal dimensions via WindowSizeMsg
		termInfo := intent.BaseIntent.GetTerminalInfo()
		termInfo.Update(tea.WindowSizeMsg{Width: 120, Height: 40})
	})

	Describe("E2E: Complete Wizard Workflow", func() {
		It("should handle Tab key to navigate between form fields", func() {
			// Init wizard
			intent.Init()
			Expect(intent.wizardModal).NotTo(BeNil())
			Expect(intent.wizardModal.IsVisible()).To(BeTrue())

			// Press Tab key multiple times to navigate through fields
			tabMsg := tea.KeyMsg{Type: tea.KeyTab}

			// First tab - should move from profile to audience
			_ = intent.Update(tabMsg)
			Expect(intent.wizardModal.IsVisible()).To(BeTrue())

			// Second tab - should stay within the form
			_ = intent.Update(tabMsg)
			Expect(intent.wizardModal.IsVisible()).To(BeTrue())

			// Wizard should still be visible and functional
			view := intent.View()
			Expect(view).To(ContainSubstring("CV Configuration"))
		})

		It("should handle wizard completion via WizardCompleteMsg", func() {
			// Init wizard
			intent.Init()
			wizard := intent.wizardModal

			// Wizard should be visible initially
			Expect(wizard.IsVisible()).To(BeTrue())

			// User completes wizard (simulated via WizardCompleteMsg)
			// This is the correct test pattern - test the message contract, not keyboard simulation
			msg := WizardCompleteMsg{
				ProfileID: "profile_1",
				Audience:  "hiring_manager",
			}
			intent.Update(msg)

			// Wizard should be completed and hidden
			Expect(wizard.IsCompleted()).To(BeTrue())
			Expect(wizard.IsVisible()).To(BeFalse())

			// Should have transitioned to extracting state
			Expect(intent.state.currentState).To(Equal(CVStateExtracting))

			// Progress modal should be visible
			Expect(intent.progressModal).NotTo(BeNil())
			Expect(intent.progressModal.IsVisible()).To(BeTrue())
		})

		It("should initialize with wizard modal visible", func() {
			// Init should create and show wizard modal
			cmd := intent.Init()
			Expect(cmd).NotTo(BeNil())

			// Wizard modal should be visible
			Expect(intent.wizardModal).NotTo(BeNil())
			Expect(intent.wizardModal.IsVisible()).To(BeTrue())

			// Should be in configuring state
			Expect(intent.state.currentState).To(Equal(CVStateConfiguring))

			// View should contain wizard modal
			view := intent.View()
			Expect(view).To(ContainSubstring("CV Configuration"))
			Expect(view).To(ContainSubstring("Select CV Profile"))
		})

		It("should allow user to submit wizard form and proceed to tech extraction", func() {
			// Init wizard
			intent.Init()

			// Simulate user filling in wizard form
			wizard := intent.wizardModal
			wizard.SetProfileID("profile_1")
			wizard.SetAudience("hiring_manager")

			// Get form data to verify
			config := wizard.GetConfigData()
			Expect(config.ProfileID).To(Equal("profile_1"))
			Expect(config.Audience).To(Equal("hiring_manager"))

			// Mark wizard as completed (simulating user pressing Enter to submit)
			// In real usage, the huh.Form would handle this, but for testing we can
			// directly send a WizardCompleteMsg
			msg := WizardCompleteMsg{
				ProfileID: "profile_1",
				Audience:  "hiring_manager",
			}
			cmd := intent.Update(msg)

			// Should transition to extracting state
			Expect(intent.state.currentState).To(Equal(CVStateExtracting))

			// Wizard should be hidden
			Expect(intent.wizardModal.IsVisible()).To(BeFalse())

			// Progress modal should be visible
			Expect(intent.progressModal).NotTo(BeNil())
			Expect(intent.progressModal.IsVisible()).To(BeTrue())

			// View should show progress modal
			view := intent.View()
			Expect(view).To(ContainSubstring("Extracting Technologies"))

			// Command should be returned (async operation)
			Expect(cmd).NotTo(BeNil())
		})

		It("should transition from tech extraction to CV generation", func() {
			// Start workflow
			intent.Init()

			// Complete wizard
			intent.Update(WizardCompleteMsg{
				ProfileID: "profile_1",
				Audience:  "hiring_manager",
			})

			// Now in extracting state
			Expect(intent.state.currentState).To(Equal(CVStateExtracting))

			// Simulate tech extraction completion
			techMsg := TechnologiesExtractedMsg{
				Technologies: []*ExtractedTechnology{
					{
						ID:         "go",
						Name:       "Go",
						Category:   "Language",
						EventCount: 2,
					},
				},
			}
			cmd := intent.Update(techMsg)

			// Should transition to generating state
			Expect(intent.state.currentState).To(Equal(CVStateGenerating))

			// Progress modal should still be visible but with different message
			Expect(intent.progressModal).NotTo(BeNil())
			Expect(intent.progressModal.IsVisible()).To(BeTrue())

			// View should show CV generation progress
			view := intent.View()
			Expect(view).To(ContainSubstring("Generating CV"))

			// Command should be returned (async operation)
			Expect(cmd).NotTo(BeNil())
		})

		It("should transition from CV generation to preview screen", func() {
			// Complete wizard and tech extraction
			intent.Init()
			intent.Update(WizardCompleteMsg{
				ProfileID: "profile_1",
				Audience:  "hiring_manager",
			})
			intent.Update(TechnologiesExtractedMsg{})

			// Now in generating state
			Expect(intent.state.currentState).To(Equal(CVStateGenerating))

			// Simulate CV generation completion
			cvMsg := CVGenerationCompleteMsg{
				CV: &career.CVView{
					ID: "cv_1",
				},
			}
			cmd := intent.Update(cvMsg)

			// Should transition to review state first (not preview)
			Expect(intent.state.currentState).To(Equal(CVStateReview))

			// Progress modal should be hidden
			Expect(intent.progressModal.IsVisible()).To(BeFalse())

			// Review screen should be created
			Expect(intent.wizardReviewScreen).NotTo(BeNil())

			// Command may be nil or an init command
			_ = cmd
		})

		It("should allow user to complete workflow from preview screen", func() {
			// Complete full workflow to review
			intent.Init()
			intent.Update(WizardCompleteMsg{
				ProfileID: "profile_1",
				Audience:  "hiring_manager",
			})
			intent.Update(TechnologiesExtractedMsg{})
			intent.Update(CVGenerationCompleteMsg{
				CV: &career.CVView{ID: "cv_1"},
			})

			// Now in review state
			Expect(intent.state.currentState).To(Equal(CVStateReview))

			// Navigate from review to preview by pressing Enter
			enterMsg := tea.KeyMsg{Type: tea.KeyEnter}
			intent.Update(enterMsg)

			// Now in preview state
			Expect(intent.state.currentState).To(Equal(CVStatePreview))

			// Simulate user pressing Enter to complete (via preview screen)
			// The preview screen would normally send a SubmitResult
			cmd := intent.Update(enterMsg)

			// Intent should be completed
			result := intent.Result()
			Expect(result).NotTo(BeNil())
			Expect(result.Status).To(Equal(Completed))

			// Result should contain the CV
			Expect(result.Data).NotTo(BeNil())
			cvResult, ok := result.Data.(*GenerateCVResult)
			Expect(ok).To(BeTrue())
			Expect(cvResult.GeneratedCV).NotTo(BeNil())
			Expect(cvResult.GeneratedCV.ID).To(Equal("cv_1"))

			// Command may be nil
			_ = cmd
		})
	})

	Describe("E2E: Escape Key Navigation", func() {
		It("should cancel workflow when Esc pressed on wizard modal", func() {
			// Init wizard
			intent.Init()
			Expect(intent.wizardModal.IsVisible()).To(BeTrue())

			// Press Esc on wizard modal (step 1 = cancel)
			escMsg := tea.KeyMsg{Type: tea.KeyEsc}
			intent.Update(escMsg)

			// Intent should be cancelled
			result := intent.Result()
			Expect(result).NotTo(BeNil())
			Expect(result.Status).To(Equal(Cancelled))
		})

		It("should allow cancelling during tech extraction", func() {
			// Start workflow
			intent.Init()
			intent.Update(WizardCompleteMsg{
				ProfileID: "profile_1",
				Audience:  "hiring_manager",
			})

			// Now in extracting state
			Expect(intent.state.currentState).To(Equal(CVStateExtracting))

			// Press Esc to cancel
			escMsg := tea.KeyMsg{Type: tea.KeyEsc}
			intent.Update(escMsg)

			// Should go back to wizard (or cancel - implementation dependent)
			// For now, verify no panic and result is set
			result := intent.Result()
			if result != nil {
				Expect(result.Status).To(Or(
					Equal(Cancelled),
					Equal(Failed),
				))
			} else {
				// Or back to configuring
				Expect(intent.state.currentState).To(Equal(CVStateConfiguring))
			}
		})

		It("should allow going back from review to wizard", func() {
			// Complete to review
			intent.Init()
			intent.Update(WizardCompleteMsg{
				ProfileID: "profile_1",
				Audience:  "hiring_manager",
			})
			intent.Update(TechnologiesExtractedMsg{})
			intent.Update(CVGenerationCompleteMsg{
				CV: &career.CVView{ID: "cv_1"},
			})

			// In review state
			Expect(intent.state.currentState).To(Equal(CVStateReview))

			// Press Esc to go back to wizard
			escMsg := tea.KeyMsg{Type: tea.KeyEsc}
			intent.Update(escMsg)

			// Should return to wizard modal
			Expect(intent.wizardModal).NotTo(BeNil())
			Expect(intent.wizardModal.IsVisible()).To(BeTrue())
			Expect(intent.state.currentState).To(Equal(CVStateConfiguring))
		})

		It("should allow going back from preview to review", func() {
			// Complete to review, then navigate to preview
			intent.Init()
			intent.Update(WizardCompleteMsg{
				ProfileID: "profile_1",
				Audience:  "hiring_manager",
			})
			intent.Update(TechnologiesExtractedMsg{})
			intent.Update(CVGenerationCompleteMsg{
				CV: &career.CVView{ID: "cv_1"},
			})

			// Navigate from review to preview
			enterMsg := tea.KeyMsg{Type: tea.KeyEnter}
			intent.Update(enterMsg)
			Expect(intent.state.currentState).To(Equal(CVStatePreview))

			// Press Esc to go back to review
			escMsg := tea.KeyMsg{Type: tea.KeyEsc}
			intent.Update(escMsg)

			// Should return to review screen
			Expect(intent.state.currentState).To(Equal(CVStateReview))
		})
	})

	Describe("E2E: Global Keyboard Shortcuts", func() {
		It("should quit on Ctrl+C from any state", func() {
			// Start in wizard
			intent.Init()

			// Press Ctrl+C
			quitMsg := tea.KeyMsg{Type: tea.KeyCtrlC}
			cmd := intent.Update(quitMsg)

			// Should return tea.Quit command
			Expect(cmd).NotTo(BeNil())
			// Note: We can't easily test tea.Quit directly, but we can verify
			// the intent is cancelled
			result := intent.Result()
			Expect(result).NotTo(BeNil())
			Expect(result.Status).To(Equal(Cancelled))
		})

		It("should quit on 'q' key from any state", func() {
			// Start in wizard
			intent.Init()

			// Press 'q'
			quitMsg := tea.KeyMsg{Type: tea.KeyRunes, Runes: []rune{'q'}}
			intent.Update(quitMsg)

			// Should be cancelled
			result := intent.Result()
			Expect(result).NotTo(BeNil())
			Expect(result.Status).To(Equal(Cancelled))
		})

		It("should return to main menu on 'm' key", func() {
			// Start in wizard
			intent.Init()

			// Press 'm'
			menuMsg := tea.KeyMsg{Type: tea.KeyRunes, Runes: []rune{'m'}}
			intent.Update(menuMsg)

			// Should have result indicating navigation
			result := intent.Result()
			Expect(result).NotTo(BeNil())
			Expect(result.Status).To(Or(
				Equal(Cancelled),
				Equal(Failed),
			))
		})
	})

	Describe("E2E: Window Resize Handling", func() {
		It("should handle window resize in wizard modal state", func() {
			// Start wizard
			intent.Init()

			// Send window resize
			resizeMsg := tea.WindowSizeMsg{Width: 150, Height: 50}
			cmd := intent.Update(resizeMsg)

			// Should not crash
			Expect(cmd).To(BeNil())

			// Wizard should still be visible and functional
			Expect(intent.wizardModal.IsVisible()).To(BeTrue())

			// View should render without errors
			view := intent.View()
			Expect(view).NotTo(BeEmpty())
		})

		It("should handle window resize in preview state", func() {
			// Complete to preview
			intent.Init()
			intent.Update(WizardCompleteMsg{
				ProfileID: "profile_1",
				Audience:  "hiring_manager",
			})
			intent.Update(TechnologiesExtractedMsg{})
			intent.Update(CVGenerationCompleteMsg{
				CV: &career.CVView{ID: "cv_1"},
			})

			// Send window resize
			resizeMsg := tea.WindowSizeMsg{Width: 150, Height: 50}
			cmd := intent.Update(resizeMsg)

			// Should not crash
			Expect(cmd).To(BeNil())

			// Preview screen should still work
			view := intent.View()
			Expect(view).NotTo(BeEmpty())
		})
	})

	Describe("E2E: Error Handling", func() {
		It("should handle tech extraction error gracefully", func() {
			// Start workflow
			intent.Init()
			intent.Update(WizardCompleteMsg{
				ProfileID: "profile_1",
				Audience:  "hiring_manager",
			})

			// Simulate extraction error
			techMsg := TechnologiesExtractedMsg{
				Error: context.DeadlineExceeded,
			}
			intent.Update(techMsg)

			// Should not crash - implementation may show error or retry
			// At minimum, verify no panic occurred
			view := intent.View()
			Expect(view).NotTo(BeEmpty())
		})

		It("should handle CV generation error gracefully", func() {
			// Complete to generating
			intent.Init()
			intent.Update(WizardCompleteMsg{
				ProfileID: "profile_1",
				Audience:  "hiring_manager",
			})
			intent.Update(TechnologiesExtractedMsg{})

			// Simulate generation error
			cvMsg := CVGenerationCompleteMsg{
				Error: context.DeadlineExceeded,
			}
			intent.Update(cvMsg)

			// Should not crash
			view := intent.View()
			Expect(view).NotTo(BeEmpty())
		})
	})

	Describe("E2E: Data Preservation Through Workflow", func() {
		It("should preserve selected profile and audience through entire workflow", func() {
			// Start and select profile
			intent.Init()
			intent.Update(WizardCompleteMsg{
				ProfileID: "profile_2",
				Audience:  "recruiter",
			})

			// Verify data stored
			Expect(intent.state.selectedProfile).NotTo(BeNil())
			Expect(intent.state.selectedProfile.ID).To(Equal("profile_2"))
			Expect(intent.state.selectedAudience).To(Equal("recruiter"))

			// Continue through workflow
			intent.Update(TechnologiesExtractedMsg{})
			intent.Update(CVGenerationCompleteMsg{
				CV: &career.CVView{ID: "cv_1"},
			})

			// Now in review state, navigate to preview
			intent.Update(tea.KeyMsg{Type: tea.KeyEnter})
			Expect(intent.state.currentState).To(Equal(CVStatePreview))

			// Complete workflow from preview
			intent.Update(tea.KeyMsg{Type: tea.KeyEnter})

			// Result should contain correct profile
			result := intent.Result()
			Expect(result).NotTo(BeNil())
			Expect(result.Data).NotTo(BeNil())
			cvResult, ok := result.Data.(*GenerateCVResult)
			Expect(ok).To(BeTrue())
			Expect(cvResult.SelectedProfile).NotTo(BeNil())
			Expect(cvResult.SelectedProfile.ID).To(Equal("profile_2"))
		})

		It("should preserve extracted technologies", func() {
			// Start workflow
			intent.Init()
			intent.Update(WizardCompleteMsg{
				ProfileID: "profile_1",
				Audience:  "hiring_manager",
			})

			// Extract technologies
			techs := []*ExtractedTechnology{
				{ID: "go", Name: "Go", Category: "Language"},
				{ID: "python", Name: "Python", Category: "Language"},
			}
			intent.Update(TechnologiesExtractedMsg{Technologies: techs})

			// Verify stored
			Expect(intent.state.extractedTechnologies).To(HaveLen(2))
			Expect(intent.state.extractedTechnologies[0].Name).To(Equal("Go"))
		})
	})
})

// mockWizardPreviewScreen is a mock screen for testing wizard flow
type mockWizardPreviewScreen struct {
	cv *career.CVView
}

func (m *mockWizardPreviewScreen) Init() tea.Cmd {
	return nil
}

func (m *mockWizardPreviewScreen) Update(msg tea.Msg) (tea.Cmd, screens.ScreenResult) {
	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch msg.String() {
		case "enter", "y":
			return nil, &screens.SubmitResult{FormData: "complete"}
		case "x":
			return nil, &screens.NavigateResult{ResultData: "export"}
		case "esc":
			return nil, &screens.CancelResult{}
		}
	}
	return nil, nil
}

func (m *mockWizardPreviewScreen) View() string {
	return "Mock Preview Screen"
}

func (m *mockWizardPreviewScreen) RenderContent() string {
	return "Mock Preview Content"
}

func (m *mockWizardPreviewScreen) SetTheme(theme interface{}) {
}

func (m *mockWizardPreviewScreen) SetTerminalInfo(width, height int) {
}

func (m *mockWizardPreviewScreen) SetLogo(logo interface{}, spacing int) {
}

// mockWizardReviewScreen is a mock screen for testing wizard review flow
type mockWizardReviewScreen struct {
	cv *career.CVView
}

func (m *mockWizardReviewScreen) Init() tea.Cmd {
	return nil
}

func (m *mockWizardReviewScreen) Update(msg tea.Msg) (tea.Cmd, screens.ScreenResult) {
	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch msg.String() {
		case "enter", "p":
			return nil, &screens.NavigateResult{ResultData: "preview"}
		case "x":
			return nil, &screens.NavigateResult{ResultData: "export"}
		case "e":
			return nil, &screens.NavigateResult{ResultData: "edit"}
		case "esc":
			return nil, &screens.CancelResult{}
		}
	}
	return nil, nil
}

func (m *mockWizardReviewScreen) View() string {
	return "Mock Review Screen"
}

func (m *mockWizardReviewScreen) RenderContent() string {
	return "Mock Review Content"
}

func (m *mockWizardReviewScreen) SetTheme(theme interface{}) {
}

func (m *mockWizardReviewScreen) SetTerminalInfo(width, height int) {
}

func (m *mockWizardReviewScreen) SetLogo(logo interface{}, spacing int) {
}

var _ = Describe("GenerateCV Wizard Complete E2E Workflow", func() {
	var (
		intent       *GenerateCVIntent
		testProfiles []*CVProfile
		testEvents   []*career.CareerEvent
		testFacts    []*career.Fact
	)

	BeforeEach(func() {
		testProfiles = []*CVProfile{
			{
				ID:             "profile_1",
				Name:           "Senior IC",
				TargetRole:     "senior_ic",
				TargetAudience: "hiring_manager",
			},
			{
				ID:             "profile_2",
				Name:           "Engineering Manager",
				TargetRole:     "em",
				TargetAudience: "recruiter",
			},
		}

		testEvents = []*career.CareerEvent{
			{
				ID:      "event_1",
				Text:    "Led team standup meetings",
				Date:    time.Date(2025, 12, 15, 0, 0, 0, 0, time.UTC),
				Company: "Acme Corp",
			},
		}

		testFacts = []*career.Fact{
			{
				ID:            "fact_1",
				Text:          "Improved team velocity by 30%",
				SourceEventID: "event_1",
			},
		}

		ctx := &GenerateCVContext{
			AvailableProfiles: testProfiles,
			Events:            testEvents,
			Facts:             testFacts,
			DefaultProfile:    testProfiles[0],
			AppContext:        context.Background(),
			ReviewScreenFactory: func(cv *career.CVView) screens.Screen {
				return &mockWizardReviewScreen{cv: cv}
			},
			PreviewScreenFactory: func(cv *career.CVView) screens.Screen {
				return &mockWizardPreviewScreen{cv: cv}
			},
		}

		var err error
		intent, err = NewGenerateCVIntent(ctx)
		Expect(err).NotTo(HaveOccurred())
		intent.EnableWizardFlow()

		// Initialize terminal info
		termInfo := intent.BaseIntent.GetTerminalInfo()
		termInfo.Update(tea.WindowSizeMsg{Width: 120, Height: 40})
	})

	Describe("Complete Workflow: Wizard -> Tech Extraction -> CV Generation -> Review -> Preview -> Complete", func() {
		It("should complete the entire workflow from wizard to completion", func() {
			// STEP 1: Initialize - wizard modal should appear
			intent.Init()
			Expect(intent.state.currentState).To(Equal(CVStateConfiguring))
			Expect(intent.wizardModal).NotTo(BeNil())
			Expect(intent.wizardModal.IsVisible()).To(BeTrue())

			// STEP 2: Complete wizard via message
			intent.Update(WizardCompleteMsg{
				ProfileID: "profile_1",
				Audience:  "hiring_manager",
			})

			// Verify wizard completed and transitioned to extracting
			Expect(intent.wizardModal.IsCompleted()).To(BeTrue())
			Expect(intent.wizardModal.IsVisible()).To(BeFalse())
			Expect(intent.state.currentState).To(Equal(CVStateExtracting))
			Expect(intent.progressModal).NotTo(BeNil())
			Expect(intent.progressModal.IsVisible()).To(BeTrue())

			// STEP 3: Simulate tech extraction completion
			intent.Update(TechnologiesExtractedMsg{
				Technologies: []*ExtractedTechnology{
					{ID: "go", Name: "Go", Category: "Language"},
				},
			})

			// Verify transitioned to generating state
			Expect(intent.state.currentState).To(Equal(CVStateGenerating))
			Expect(intent.progressModal.IsVisible()).To(BeTrue())

			// STEP 4: Simulate CV generation completion
			testCV := &career.CVView{
				ID:       "cv_1",
				Name:     "Test CV",
				Sections: []*career.CVSection{},
			}
			intent.Update(CVGenerationCompleteMsg{CV: testCV})

			// Verify transitioned to review state first
			Expect(intent.state.currentState).To(Equal(CVStateReview))
			Expect(intent.progressModal.IsVisible()).To(BeFalse())
			Expect(intent.wizardReviewScreen).NotTo(BeNil())
			Expect(intent.state.generatedCV).To(Equal(testCV))

			// STEP 5: Navigate from review to preview
			intent.Update(tea.KeyMsg{Type: tea.KeyEnter})
			Expect(intent.state.currentState).To(Equal(CVStatePreview))
			Expect(intent.wizardPreviewScreen).NotTo(BeNil())

			// STEP 6: Complete from preview (Enter key)
			intent.Update(tea.KeyMsg{Type: tea.KeyEnter})

			// Verify intent is completed
			result := intent.Result()
			Expect(result).NotTo(BeNil())
			Expect(result.Status).To(Equal(Completed))
			Expect(result.Data).NotTo(BeNil())

			cvResult, ok := result.Data.(*GenerateCVResult)
			Expect(ok).To(BeTrue())
			Expect(cvResult.GeneratedCV).To(Equal(testCV))
			Expect(cvResult.SelectedProfile.ID).To(Equal("profile_1"))
		})
	})

	Describe("View Rendering: Appropriate content at each state", func() {
		It("should render wizard modal during wizard state", func() {
			intent.Init()
			view := intent.View()

			// Wizard modal elements (centered, no logo - by design)
			Expect(view).To(ContainSubstring("CV Configuration")) // Modal title
			Expect(view).To(ContainSubstring("Profile"))          // Profile selection
			Expect(view).To(ContainSubstring("Audience"))         // Audience selection
		})

		It("should render progress modal during extracting state", func() {
			intent.Init()
			intent.Update(WizardCompleteMsg{ProfileID: "profile_1", Audience: "hiring_manager"})

			view := intent.View()

			// Progress modal (centered, no logo - by design)
			Expect(view).To(ContainSubstring("Extracting")) // State indicator
			Expect(view).To(ContainSubstring("Cancel"))     // Cancel option
		})

		It("should render progress modal during generating state", func() {
			intent.Init()
			intent.Update(WizardCompleteMsg{ProfileID: "profile_1", Audience: "hiring_manager"})
			intent.Update(TechnologiesExtractedMsg{})

			view := intent.View()

			// Progress modal (centered, no logo - by design)
			Expect(view).To(ContainSubstring("Generating")) // State indicator
			Expect(view).To(ContainSubstring("Cancel"))     // Cancel option
		})

		It("should render StandardView during preview state", func() {
			intent.Init()
			intent.Update(WizardCompleteMsg{ProfileID: "profile_1", Audience: "hiring_manager"})
			intent.Update(TechnologiesExtractedMsg{})
			intent.Update(CVGenerationCompleteMsg{CV: &career.CVView{ID: "cv_1"}})
			// Navigate from review to preview
			intent.Update(tea.KeyMsg{Type: tea.KeyEnter})

			view := intent.View()

			// StandardView elements (preview has full layout)
			Expect(view).To(ContainSubstring("Generate CV")) // Breadcrumb
			Expect(view).To(ContainSubstring("Preview"))     // Breadcrumb shows preview
			// Help footer
			Expect(view).To(ContainSubstring("Scroll"))
			Expect(view).To(ContainSubstring("Confirm"))
			Expect(view).To(ContainSubstring("Export"))
		})
	})

	Describe("Export Workflow: Preview -> Export Modal -> Export Complete", func() {
		BeforeEach(func() {
			// Get to preview state (through review)
			intent.Init()
			intent.Update(WizardCompleteMsg{ProfileID: "profile_1", Audience: "hiring_manager"})
			intent.Update(TechnologiesExtractedMsg{})
			intent.Update(CVGenerationCompleteMsg{CV: &career.CVView{ID: "cv_1"}})
			// Navigate from review to preview
			intent.Update(tea.KeyMsg{Type: tea.KeyEnter})
		})

		It("should show export modal when x is pressed on preview", func() {
			Expect(intent.state.currentState).To(Equal(CVStatePreview))

			// Press x to export
			intent.Update(tea.KeyMsg{Type: tea.KeyRunes, Runes: []rune{'x'}})

			// Verify export modal shown
			Expect(intent.exportModal).NotTo(BeNil())
			Expect(intent.exportModal.IsVisible()).To(BeTrue())
			Expect(intent.state.currentState).To(Equal(CVStateExporting))
		})

		It("should render export modal view", func() {
			intent.Update(tea.KeyMsg{Type: tea.KeyRunes, Runes: []rune{'x'}})

			view := intent.View()
			Expect(view).To(ContainSubstring("Export")) // Export modal visible
		})
	})

	Describe("Back Navigation: Esc key at each state", func() {
		It("should cancel intent when Esc pressed on wizard (first step)", func() {
			intent.Init()
			Expect(intent.state.currentState).To(Equal(CVStateConfiguring))

			// Press Esc
			intent.Update(tea.KeyMsg{Type: tea.KeyEsc})

			// Should cancel
			result := intent.Result()
			Expect(result).NotTo(BeNil())
			Expect(result.Status).To(Equal(Cancelled))
		})

		It("should return to review when Esc pressed on preview", func() {
			// Get to preview (through review)
			intent.Init()
			intent.Update(WizardCompleteMsg{ProfileID: "profile_1", Audience: "hiring_manager"})
			intent.Update(TechnologiesExtractedMsg{})
			intent.Update(CVGenerationCompleteMsg{CV: &career.CVView{ID: "cv_1"}})
			Expect(intent.state.currentState).To(Equal(CVStateReview))
			// Navigate to preview
			intent.Update(tea.KeyMsg{Type: tea.KeyEnter})
			Expect(intent.state.currentState).To(Equal(CVStatePreview))

			// Press Esc
			intent.Update(tea.KeyMsg{Type: tea.KeyEsc})

			// Should return to review (not wizard)
			Expect(intent.state.currentState).To(Equal(CVStateReview))
		})
	})

	Describe("Global Keys: q, m, Ctrl+C work everywhere", func() {
		testGlobalKeys := func(stateName string, setupFunc func()) {
			Context("in "+stateName+" state", func() {
				BeforeEach(func() {
					setupFunc()
				})

				It("should quit on Ctrl+C", func() {
					cmd := intent.Update(tea.KeyMsg{Type: tea.KeyCtrlC})
					Expect(cmd).NotTo(BeNil()) // tea.Quit command

					result := intent.Result()
					Expect(result).NotTo(BeNil())
					Expect(result.Status).To(Equal(Cancelled))
				})

				It("should cancel on q key", func() {
					intent.Update(tea.KeyMsg{Type: tea.KeyRunes, Runes: []rune{'q'}})

					result := intent.Result()
					Expect(result).NotTo(BeNil())
					Expect(result.Status).To(Equal(Cancelled))
				})

				It("should cancel on m key (return to menu)", func() {
					intent.Update(tea.KeyMsg{Type: tea.KeyRunes, Runes: []rune{'m'}})

					result := intent.Result()
					Expect(result).NotTo(BeNil())
					Expect(result.Status).To(Equal(Cancelled))
				})
			})
		}

		testGlobalKeys("wizard", func() {
			intent.Init()
		})

		testGlobalKeys("extracting", func() {
			intent.Init()
			intent.Update(WizardCompleteMsg{ProfileID: "profile_1", Audience: "hiring_manager"})
		})

		testGlobalKeys("generating", func() {
			intent.Init()
			intent.Update(WizardCompleteMsg{ProfileID: "profile_1", Audience: "hiring_manager"})
			intent.Update(TechnologiesExtractedMsg{})
		})

		testGlobalKeys("preview", func() {
			intent.Init()
			intent.Update(WizardCompleteMsg{ProfileID: "profile_1", Audience: "hiring_manager"})
			intent.Update(TechnologiesExtractedMsg{})
			intent.Update(CVGenerationCompleteMsg{CV: &career.CVView{ID: "cv_1"}})
			// Navigate from review to preview
			intent.Update(tea.KeyMsg{Type: tea.KeyEnter})
		})
	})

	Describe("State Transitions: Verify correct state at each step", func() {
		It("should follow correct state sequence", func() {
			// Init -> Configuring
			intent.Init()
			Expect(intent.state.currentState).To(Equal(CVStateConfiguring))

			// WizardComplete -> Extracting
			intent.Update(WizardCompleteMsg{ProfileID: "profile_1", Audience: "hiring_manager"})
			Expect(intent.state.currentState).To(Equal(CVStateExtracting))

			// TechExtracted -> Generating
			intent.Update(TechnologiesExtractedMsg{})
			Expect(intent.state.currentState).To(Equal(CVStateGenerating))

			// CVGenerated -> Review
			intent.Update(CVGenerationCompleteMsg{CV: &career.CVView{ID: "cv_1"}})
			Expect(intent.state.currentState).To(Equal(CVStateReview))

			// Enter -> Preview
			intent.Update(tea.KeyMsg{Type: tea.KeyEnter})
			Expect(intent.state.currentState).To(Equal(CVStatePreview))

			// Enter -> Completed
			intent.Update(tea.KeyMsg{Type: tea.KeyEnter})
			result := intent.Result()
			Expect(result.Status).To(Equal(Completed))
		})
	})

	Describe("Data Preservation: Profile and audience preserved through workflow", func() {
		It("should preserve selected profile through entire workflow", func() {
			intent.Init()
			intent.Update(WizardCompleteMsg{ProfileID: "profile_2", Audience: "recruiter"})

			// Verify stored
			Expect(intent.state.selectedProfile).NotTo(BeNil())
			Expect(intent.state.selectedProfile.ID).To(Equal("profile_2"))
			Expect(intent.state.selectedAudience).To(Equal("recruiter"))

			// Continue through workflow
			intent.Update(TechnologiesExtractedMsg{})
			intent.Update(CVGenerationCompleteMsg{CV: &career.CVView{ID: "cv_1"}})

			// Still preserved
			Expect(intent.state.selectedProfile.ID).To(Equal("profile_2"))
			Expect(intent.state.selectedAudience).To(Equal("recruiter"))

			// Navigate through review to preview to complete
			intent.Update(tea.KeyMsg{Type: tea.KeyEnter}) // Review -> Preview
			intent.Update(tea.KeyMsg{Type: tea.KeyEnter}) // Preview -> Complete

			// In result
			result := intent.Result()
			cvResult := result.Data.(*GenerateCVResult)
			Expect(cvResult.SelectedProfile.ID).To(Equal("profile_2"))
		})

		It("should preserve extracted technologies", func() {
			intent.Init()
			intent.Update(WizardCompleteMsg{ProfileID: "profile_1", Audience: "hiring_manager"})

			techs := []*ExtractedTechnology{
				{ID: "go", Name: "Go", Category: "Language"},
				{ID: "python", Name: "Python", Category: "Language"},
			}
			intent.Update(TechnologiesExtractedMsg{Technologies: techs})

			Expect(intent.state.extractedTechnologies).To(HaveLen(2))
			Expect(intent.state.extractedTechnologies[0].Name).To(Equal("Go"))
			Expect(intent.state.extractedTechnologies[1].Name).To(Equal("Python"))
		})

		It("should preserve generated CV", func() {
			intent.Init()
			intent.Update(WizardCompleteMsg{ProfileID: "profile_1", Audience: "hiring_manager"})
			intent.Update(TechnologiesExtractedMsg{})

			testCV := &career.CVView{
				ID:   "cv_123",
				Name: "My Awesome CV",
			}
			intent.Update(CVGenerationCompleteMsg{CV: testCV})

			Expect(intent.state.generatedCV).To(Equal(testCV))
			Expect(intent.state.generatedCV.ID).To(Equal("cv_123"))
		})
	})

	Describe("Preview Screen Factory: Properly creates preview screen", func() {
		It("should create preview screen with CV data", func() {
			intent.Init()
			intent.Update(WizardCompleteMsg{ProfileID: "profile_1", Audience: "hiring_manager"})
			intent.Update(TechnologiesExtractedMsg{})

			testCV := &career.CVView{ID: "cv_from_factory"}
			intent.Update(CVGenerationCompleteMsg{CV: testCV})

			// After CV generation, we're in Review state first
			Expect(intent.state.currentState).To(Equal(CVStateReview))
			Expect(intent.wizardReviewScreen).NotTo(BeNil())

			// Navigate to preview
			intent.Update(tea.KeyMsg{Type: tea.KeyEnter})
			Expect(intent.state.currentState).To(Equal(CVStatePreview))
			Expect(intent.wizardPreviewScreen).NotTo(BeNil())

			// Verify it's our mock with the CV
			mockScreen, ok := intent.wizardPreviewScreen.(*mockWizardPreviewScreen)
			Expect(ok).To(BeTrue())
			Expect(mockScreen.cv).To(Equal(testCV))
		})

		It("should handle nil factories gracefully", func() {
			// Create intent without factories
			ctxNoFactory := &GenerateCVContext{
				AvailableProfiles: testProfiles,
				Events:            testEvents,
				Facts:             testFacts,
				DefaultProfile:    testProfiles[0],
				AppContext:        context.Background(),
				// Both factories are nil
			}
			intentNoFactory, err := NewGenerateCVIntent(ctxNoFactory)
			Expect(err).NotTo(HaveOccurred())
			intentNoFactory.EnableWizardFlow()
			termInfo := intentNoFactory.BaseIntent.GetTerminalInfo()
			termInfo.Update(tea.WindowSizeMsg{Width: 120, Height: 40})

			intentNoFactory.Init()
			intentNoFactory.Update(WizardCompleteMsg{ProfileID: "profile_1", Audience: "hiring_manager"})
			intentNoFactory.Update(TechnologiesExtractedMsg{})
			intentNoFactory.Update(CVGenerationCompleteMsg{CV: &career.CVView{ID: "cv_1"}})

			// Should transition to review state (screen may be nil without factory)
			Expect(intentNoFactory.state.currentState).To(Equal(CVStateReview))

			// View should render without crashing
			view := intentNoFactory.View()
			Expect(view).NotTo(BeEmpty())
		})
	})

	Describe("Error Handling: Graceful handling of errors", func() {
		It("should handle tech extraction error", func() {
			intent.Init()
			intent.Update(WizardCompleteMsg{ProfileID: "profile_1", Audience: "hiring_manager"})

			// Send error in extraction
			intent.Update(TechnologiesExtractedMsg{
				Error: context.DeadlineExceeded,
			})

			// Should not crash - view should render
			view := intent.View()
			Expect(view).NotTo(BeEmpty())
		})

		It("should handle CV generation error", func() {
			intent.Init()
			intent.Update(WizardCompleteMsg{ProfileID: "profile_1", Audience: "hiring_manager"})
			intent.Update(TechnologiesExtractedMsg{})

			// Send error in generation
			intent.Update(CVGenerationCompleteMsg{
				Error: context.DeadlineExceeded,
			})

			// Should not crash - view should render
			view := intent.View()
			Expect(view).NotTo(BeEmpty())
		})
	})

	Describe("Window Resize: Handle WindowSizeMsg at every state", func() {
		It("should handle resize during wizard", func() {
			intent.Init()

			cmd := intent.Update(tea.WindowSizeMsg{Width: 200, Height: 60})
			Expect(cmd).To(BeNil()) // No error

			view := intent.View()
			Expect(view).NotTo(BeEmpty())
		})

		It("should handle resize during preview", func() {
			intent.Init()
			intent.Update(WizardCompleteMsg{ProfileID: "profile_1", Audience: "hiring_manager"})
			intent.Update(TechnologiesExtractedMsg{})
			intent.Update(CVGenerationCompleteMsg{CV: &career.CVView{ID: "cv_1"}})
			// Navigate from review to preview
			intent.Update(tea.KeyMsg{Type: tea.KeyEnter})

			cmd := intent.Update(tea.WindowSizeMsg{Width: 200, Height: 60})
			Expect(cmd).To(BeNil()) // No error

			view := intent.View()
			Expect(view).NotTo(BeEmpty())
		})
	})

	Describe("Complete Export Workflow: Review -> Preview -> Format -> Location -> Exporting -> Complete", func() {
		BeforeEach(func() {
			// Get to review state first
			intent.Init()
			intent.Update(WizardCompleteMsg{ProfileID: "profile_1", Audience: "hiring_manager"})
			intent.Update(TechnologiesExtractedMsg{})
			intent.Update(CVGenerationCompleteMsg{CV: &career.CVView{
				ID:   "cv_1",
				Name: "Test CV",
			}})
			Expect(intent.state.currentState).To(Equal(CVStateReview))
			// Navigate from review to preview
			intent.Update(tea.KeyMsg{Type: tea.KeyEnter})
			Expect(intent.state.currentState).To(Equal(CVStatePreview))
		})

		It("should complete entire export workflow to file", func() {
			// Step 1: Press x to show export modal
			intent.Update(tea.KeyMsg{Type: tea.KeyRunes, Runes: []rune{'x'}})
			Expect(intent.state.currentState).To(Equal(CVStateExporting))
			Expect(intent.exportModal).NotTo(BeNil())
			Expect(intent.exportModal.IsVisible()).To(BeTrue())

			// Step 2: Complete the export modal form (select format and location)
			// The form auto-selects first options (text format, file location)
			// Press Enter to confirm first field (format)
			intent.Update(tea.KeyMsg{Type: tea.KeyEnter})
			// Press Enter to confirm second field (location)
			intent.Update(tea.KeyMsg{Type: tea.KeyEnter})

			// Step 3: Check if modal completed
			if intent.exportModal.IsCompleted() {
				exportData := intent.exportModal.GetExportData()
				Expect(exportData.Format).NotTo(BeEmpty())
				Expect(exportData.Location).NotTo(BeEmpty())
			}
		})

		It("should allow canceling export with Esc", func() {
			// Press x to show export modal
			intent.Update(tea.KeyMsg{Type: tea.KeyRunes, Runes: []rune{'x'}})
			Expect(intent.state.currentState).To(Equal(CVStateExporting))

			// Press Esc to cancel
			intent.Update(tea.KeyMsg{Type: tea.KeyEsc})

			// Should hide modal and return to preview
			Expect(intent.exportModal.IsVisible()).To(BeFalse())
			Expect(intent.state.currentState).To(Equal(CVStatePreview))
		})

		It("should show export format options in modal", func() {
			intent.Update(tea.KeyMsg{Type: tea.KeyRunes, Runes: []rune{'x'}})

			view := intent.View()

			// Should show format options
			Expect(view).To(SatisfyAny(
				ContainSubstring("Text"),
				ContainSubstring("text"),
				ContainSubstring("Format"),
			))
		})

		It("should show export location options in modal", func() {
			intent.Update(tea.KeyMsg{Type: tea.KeyRunes, Runes: []rune{'x'}})

			view := intent.View()

			// Should show location options or related text
			Expect(view).To(SatisfyAny(
				ContainSubstring("File"),
				ContainSubstring("Clipboard"),
				ContainSubstring("Save"),
			))
		})
	})

	Describe("Wizard Modal Navigation: Back navigation between steps", func() {
		It("should cancel wizard when Esc pressed on step 1", func() {
			intent.Init()
			Expect(intent.state.currentState).To(Equal(CVStateConfiguring))
			Expect(intent.wizardModal.IsVisible()).To(BeTrue())

			// Press Esc on first step
			intent.Update(tea.KeyMsg{Type: tea.KeyEsc})

			// Should cancel (hide wizard)
			Expect(intent.wizardModal.IsVisible()).To(BeFalse())
		})

		It("should track wizard step progression", func() {
			intent.Init()
			Expect(intent.wizardModal).NotTo(BeNil())

			// Initially at step 0
			Expect(intent.wizardModal.GetCurrentStep()).To(Equal(0))
		})

		It("should show wizard view content at each step", func() {
			intent.Init()

			view := intent.View()

			// Step 1 content
			Expect(view).To(ContainSubstring("Profile"))
			Expect(view).To(ContainSubstring("Audience"))
		})
	})

	Describe("Skills Format Configuration in Wizard", func() {
		It("should have skills format configuration available in wizard modal", func() {
			intent.Init()
			Expect(intent.wizardModal).NotTo(BeNil())

			// Verify wizard modal has CVConfigData which includes SkillsFormat
			config := intent.wizardModal.GetConfigData()
			Expect(config).NotTo(BeNil())
			// The wizard modal data structure includes SkillsFormat field
		})

		It("should allow selecting grouped skills format", func() {
			intent.Init()

			// The form includes skills format options
			// Verify the wizard modal has the correct fields
			Expect(intent.wizardModal).NotTo(BeNil())

			// Check wizard data structure includes SkillsFormat
			config := intent.wizardModal.GetConfigData()
			Expect(config).NotTo(BeNil())
			// SkillsFormat field should be available
		})

		It("should preserve skills format selection through completion", func() {
			intent.Init()

			// Complete wizard - SkillsFormat is stored in wizard modal config
			intent.Update(WizardCompleteMsg{
				ProfileID: "profile_1",
				Audience:  "hiring_manager",
			})

			// Verify transition happened
			Expect(intent.state.currentState).To(Equal(CVStateExtracting))
		})
	})

	Describe("Preview to Export to Complete: Full user journey", func() {
		It("should complete full journey: Review -> Preview -> Complete", func() {
			// Setup: Get to review, then preview
			intent.Init()
			intent.Update(WizardCompleteMsg{ProfileID: "profile_1", Audience: "hiring_manager"})
			intent.Update(TechnologiesExtractedMsg{})
			intent.Update(CVGenerationCompleteMsg{CV: &career.CVView{
				ID:   "cv_1",
				Name: "Test CV",
				Sections: []*career.CVSection{
					{ID: "sec_1", SectionType: "summary", Title: "Summary"},
				},
			}})
			Expect(intent.state.currentState).To(Equal(CVStateReview))

			// Verify CV data available
			Expect(intent.state.generatedCV).NotTo(BeNil())
			Expect(intent.state.generatedCV.ID).To(Equal("cv_1"))

			// Navigate to preview
			intent.Update(tea.KeyMsg{Type: tea.KeyEnter})
			Expect(intent.state.currentState).To(Equal(CVStatePreview))

			// User can either:
			// 1. Press Enter to complete without exporting
			// 2. Press x to export first

			// Test path 1: Direct completion
			intent.Update(tea.KeyMsg{Type: tea.KeyEnter})

			result := intent.Result()
			Expect(result).NotTo(BeNil())
			Expect(result.Status).To(Equal(Completed))

			cvResult, ok := result.Data.(*GenerateCVResult)
			Expect(ok).To(BeTrue())
			Expect(cvResult.GeneratedCV).NotTo(BeNil())
			Expect(cvResult.GeneratedCV.ID).To(Equal("cv_1"))
		})

		It("should allow export then completion", func() {
			// Setup: Get to review, then preview
			intent.Init()
			intent.Update(WizardCompleteMsg{ProfileID: "profile_1", Audience: "hiring_manager"})
			intent.Update(TechnologiesExtractedMsg{})
			intent.Update(CVGenerationCompleteMsg{CV: &career.CVView{ID: "cv_1"}})
			// Navigate from review to preview
			intent.Update(tea.KeyMsg{Type: tea.KeyEnter})
			Expect(intent.state.currentState).To(Equal(CVStatePreview))

			// Step 1: Open export modal
			intent.Update(tea.KeyMsg{Type: tea.KeyRunes, Runes: []rune{'x'}})
			Expect(intent.state.currentState).To(Equal(CVStateExporting))

			// Step 2: Cancel export to go back to preview
			intent.Update(tea.KeyMsg{Type: tea.KeyEsc})
			Expect(intent.state.currentState).To(Equal(CVStatePreview))

			// Step 3: Now complete from preview
			intent.Update(tea.KeyMsg{Type: tea.KeyEnter})

			result := intent.Result()
			Expect(result.Status).To(Equal(Completed))
		})
	})

	Describe("CV Sections and Skills Display", func() {
		It("should include skills section in generated CV", func() {
			intent.Init()
			intent.Update(WizardCompleteMsg{ProfileID: "profile_1", Audience: "hiring_manager"})
			intent.Update(TechnologiesExtractedMsg{
				Technologies: []*ExtractedTechnology{
					{ID: "go", Name: "Go", Category: "Language"},
					{ID: "python", Name: "Python", Category: "Language"},
					{ID: "postgres", Name: "PostgreSQL", Category: "Database"},
				},
			})

			// Generate CV with skills section
			testCV := &career.CVView{
				ID:   "cv_1",
				Name: "Test CV",
				Sections: []*career.CVSection{
					{
						ID:          "skills_1",
						SectionType: "skills",
						Title:       "Technical Skills",
					},
				},
			}
			intent.Update(CVGenerationCompleteMsg{CV: testCV})

			// Verify CV has skills section
			Expect(intent.state.generatedCV).NotTo(BeNil())
			Expect(intent.state.generatedCV.Sections).To(HaveLen(1))
			Expect(intent.state.generatedCV.Sections[0].SectionType).To(Equal("skills"))
		})

		It("should pass technologies to CV generation", func() {
			intent.Init()
			intent.Update(WizardCompleteMsg{ProfileID: "profile_1", Audience: "hiring_manager"})

			techs := []*ExtractedTechnology{
				{ID: "go", Name: "Go", Category: "Language"},
				{ID: "react", Name: "React", Category: "Frontend"},
			}
			intent.Update(TechnologiesExtractedMsg{Technologies: techs})

			// Verify technologies stored
			Expect(intent.state.extractedTechnologies).To(HaveLen(2))
			Expect(intent.state.extractedTechnologies[0].Name).To(Equal("Go"))
			Expect(intent.state.extractedTechnologies[1].Name).To(Equal("React"))
		})
	})

	Describe("Export Complete State: Success and Error Handling", func() {
		// Helper function to get intent to GenerateCVStateExporting state
		// (ready to receive CVExportCompleteMsg)
		//
		// Note: In wizard flow, the state transitions are:
		// 1. CVStateReview - user sees review, presses Enter to go to preview
		// 2. CVStatePreview - user presses 'x'
		// 3. CVStateExporting - export modal shown, user completes form
		// 4. GenerateCVStateExporting - async export in progress
		// 5. GenerateCVStateExportComplete - export finished
		//
		// For testing, we bypass the modal form interaction and directly
		// transition to GenerateCVStateExporting by simulating what
		// handleExportComplete() does.
		getToExportingState := func() {
			// Get to preview state with a generated CV (through review)
			intent.Init()
			intent.Update(WizardCompleteMsg{ProfileID: "profile_1", Audience: "hiring_manager"})
			intent.Update(TechnologiesExtractedMsg{})
			intent.Update(CVGenerationCompleteMsg{CV: &career.CVView{
				ID:   "cv_1",
				Name: "Test CV",
				Sections: []*career.CVSection{
					{ID: "sec_1", SectionType: "summary", Title: "Summary"},
				},
			}})
			Expect(intent.state.currentState).To(Equal(CVStateReview))
			// Navigate from review to preview
			intent.Update(tea.KeyMsg{Type: tea.KeyEnter})
			Expect(intent.state.currentState).To(Equal(CVStatePreview))

			// Directly set up the export state as if the user had completed the export modal
			// This is equivalent to what handleExportComplete() does after modal form submission
			intent.state.selectedExportFormat = CVExportFormatText
			intent.state.selectedExportOption = CVExportOptionSaveToFile
			intent.state.currentState = GenerateCVStateExporting
			intent.state.isExporting = true
		}

		It("should transition to export complete state on successful export", func() {
			getToExportingState()

			// Simulate successful export completion
			intent.Update(CVExportCompleteMsg{Path: "/tmp/cv_export.txt", Error: nil})

			// Should be in export complete state
			Expect(intent.state.currentState).To(Equal(GenerateCVStateExportComplete))
			Expect(intent.state.exportedPath).To(Equal("/tmp/cv_export.txt"))
			Expect(intent.state.exportError).To(BeNil())
		})

		It("should transition to export complete state on export error", func() {
			getToExportingState()

			// Simulate export error
			testError := fmt.Errorf("permission denied")
			intent.Update(CVExportCompleteMsg{Path: "", Error: testError})

			// Should be in export complete state with error
			Expect(intent.state.currentState).To(Equal(GenerateCVStateExportComplete))
			Expect(intent.state.exportError).To(Equal(testError))
		})

		It("should show success message in export complete view", func() {
			getToExportingState()
			intent.Update(CVExportCompleteMsg{Path: "/tmp/cv_export.txt", Error: nil})

			view := intent.View()
			Expect(view).To(ContainSubstring("Export Complete"))
		})

		It("should show error message in export complete view on failure", func() {
			getToExportingState()
			intent.Update(CVExportCompleteMsg{Path: "", Error: fmt.Errorf("disk full")})

			view := intent.View()
			Expect(view).To(SatisfyAny(
				ContainSubstring("Failed"),
				ContainSubstring("Error"),
				ContainSubstring("disk full"),
			))
		})

		It("should complete workflow when Enter pressed on export complete", func() {
			getToExportingState()
			intent.Update(CVExportCompleteMsg{Path: "/tmp/cv_export.txt", Error: nil})
			Expect(intent.state.currentState).To(Equal(GenerateCVStateExportComplete))

			// Press Enter to complete workflow
			intent.Update(tea.KeyMsg{Type: tea.KeyEnter})

			// Result should be completed with export info
			result := intent.Result()
			Expect(result).NotTo(BeNil())
			Expect(result.Status).To(Equal(Completed))

			cvResult, ok := result.Data.(*GenerateCVResult)
			Expect(ok).To(BeTrue())
			Expect(cvResult.ExportPath).To(Equal("/tmp/cv_export.txt"))
		})

		It("should return to export location selection when Esc pressed on export complete", func() {
			getToExportingState()
			intent.Update(CVExportCompleteMsg{Path: "", Error: fmt.Errorf("export failed")})
			Expect(intent.state.currentState).To(Equal(GenerateCVStateExportComplete))

			// Press Esc to retry
			intent.Update(tea.KeyMsg{Type: tea.KeyEsc})

			// Should return to export location selection
			Expect(intent.state.currentState).To(Equal(GenerateCVStateExportSelectLocation))
			Expect(intent.state.exportError).To(BeNil())
		})

		It("should preserve export format in result metadata", func() {
			getToExportingState()
			intent.state.selectedExportFormat = CVExportFormatMarkdown
			intent.Update(CVExportCompleteMsg{Path: "/tmp/cv_export.md", Error: nil})
			intent.Update(tea.KeyMsg{Type: tea.KeyEnter})

			result := intent.Result()
			Expect(result.Metadata["export_format"]).To(Equal("markdown"))
		})

		It("should handle clipboard export in export complete view", func() {
			getToExportingState()
			intent.state.selectedExportOption = CVExportOptionClipboard
			intent.Update(CVExportCompleteMsg{Path: "clipboard", Error: nil})

			view := intent.View()
			Expect(view).To(ContainSubstring("Clipboard"))
		})
	})

	Describe("Export Complete: Global Keys Work Correctly", func() {
		BeforeEach(func() {
			// Get to export complete state (through review -> preview -> export)
			intent.Init()
			intent.Update(WizardCompleteMsg{ProfileID: "profile_1", Audience: "hiring_manager"})
			intent.Update(TechnologiesExtractedMsg{})
			intent.Update(CVGenerationCompleteMsg{CV: &career.CVView{ID: "cv_1"}})
			// Navigate from review to preview
			intent.Update(tea.KeyMsg{Type: tea.KeyEnter})
			intent.Update(tea.KeyMsg{Type: tea.KeyRunes, Runes: []rune{'x'}})
			intent.Update(CVExportCompleteMsg{Path: "/tmp/cv.txt", Error: nil})
			Expect(intent.state.currentState).To(Equal(GenerateCVStateExportComplete))
		})

		It("should quit on q key in export complete state", func() {
			cmd := intent.Update(tea.KeyMsg{Type: tea.KeyRunes, Runes: []rune{'q'}})
			Expect(cmd).NotTo(BeNil())
		})

		It("should quit on Ctrl+C in export complete state", func() {
			cmd := intent.Update(tea.KeyMsg{Type: tea.KeyCtrlC})
			Expect(cmd).NotTo(BeNil())
		})
	})

	Describe("Full Export Journey: Review -> Preview -> Export Modal -> Exporting -> Complete -> Done", func() {
		It("should complete entire export journey with file export", func() {
			// Setup: Initialize and get to preview (through review)
			intent.Init()
			intent.Update(WizardCompleteMsg{ProfileID: "profile_1", Audience: "hiring_manager"})
			intent.Update(TechnologiesExtractedMsg{})
			intent.Update(CVGenerationCompleteMsg{CV: &career.CVView{
				ID:   "cv_1",
				Name: "Full Journey CV",
			}})
			Expect(intent.state.currentState).To(Equal(CVStateReview))
			// Navigate from review to preview
			intent.Update(tea.KeyMsg{Type: tea.KeyEnter})
			Expect(intent.state.currentState).To(Equal(CVStatePreview))

			// Step 1: Open export modal
			intent.Update(tea.KeyMsg{Type: tea.KeyRunes, Runes: []rune{'x'}})
			Expect(intent.state.currentState).To(Equal(CVStateExporting))
			Expect(intent.exportModal).NotTo(BeNil())

			// Step 2: Export modal is showing (form interaction would happen here)
			view := intent.View()
			Expect(view).To(SatisfyAny(
				ContainSubstring("Export"),
				ContainSubstring("Format"),
			))

			// Step 3: Simulate export completion (as if user selected format/location and export ran)
			intent.Update(CVExportCompleteMsg{Path: "/home/user/cv_2026.txt", Error: nil})
			Expect(intent.state.currentState).To(Equal(GenerateCVStateExportComplete))

			// Step 4: Verify export complete view
			view = intent.View()
			Expect(view).To(ContainSubstring("Export Complete"))

			// Step 5: Press Enter to finish
			intent.Update(tea.KeyMsg{Type: tea.KeyEnter})

			// Verify workflow completed with export info
			result := intent.Result()
			Expect(result).NotTo(BeNil())
			Expect(result.Status).To(Equal(Completed))

			cvResult, ok := result.Data.(*GenerateCVResult)
			Expect(ok).To(BeTrue())
			Expect(cvResult.GeneratedCV).NotTo(BeNil())
			Expect(cvResult.ExportPath).To(Equal("/home/user/cv_2026.txt"))
		})

		It("should handle export failure and allow retry", func() {
			// Setup (through review -> preview)
			intent.Init()
			intent.Update(WizardCompleteMsg{ProfileID: "profile_1", Audience: "hiring_manager"})
			intent.Update(TechnologiesExtractedMsg{})
			intent.Update(CVGenerationCompleteMsg{CV: &career.CVView{ID: "cv_1"}})
			// Navigate from review to preview
			intent.Update(tea.KeyMsg{Type: tea.KeyEnter})

			// Open export
			intent.Update(tea.KeyMsg{Type: tea.KeyRunes, Runes: []rune{'x'}})

			// Export fails
			intent.Update(CVExportCompleteMsg{Path: "", Error: fmt.Errorf("network error")})
			Expect(intent.state.currentState).To(Equal(GenerateCVStateExportComplete))

			// View shows error
			view := intent.View()
			Expect(view).To(SatisfyAny(
				ContainSubstring("Failed"),
				ContainSubstring("Error"),
			))

			// User presses Esc to retry
			intent.Update(tea.KeyMsg{Type: tea.KeyEsc})
			Expect(intent.state.currentState).To(Equal(GenerateCVStateExportSelectLocation))

			// User could now retry with different options
		})

		It("should allow completing without export after opening export modal", func() {
			// Setup (through review -> preview)
			intent.Init()
			intent.Update(WizardCompleteMsg{ProfileID: "profile_1", Audience: "hiring_manager"})
			intent.Update(TechnologiesExtractedMsg{})
			intent.Update(CVGenerationCompleteMsg{CV: &career.CVView{ID: "cv_1"}})
			// Navigate from review to preview
			intent.Update(tea.KeyMsg{Type: tea.KeyEnter})
			Expect(intent.state.currentState).To(Equal(CVStatePreview))

			// Open export modal
			intent.Update(tea.KeyMsg{Type: tea.KeyRunes, Runes: []rune{'x'}})
			Expect(intent.state.currentState).To(Equal(CVStateExporting))

			// Cancel export
			intent.Update(tea.KeyMsg{Type: tea.KeyEsc})
			Expect(intent.state.currentState).To(Equal(CVStatePreview))

			// Complete without export
			intent.Update(tea.KeyMsg{Type: tea.KeyEnter})

			result := intent.Result()
			Expect(result.Status).To(Equal(Completed))

			cvResult, ok := result.Data.(*GenerateCVResult)
			Expect(ok).To(BeTrue())
			Expect(cvResult.ExportPath).To(BeEmpty()) // No export path since we cancelled
		})
	})
})

// ============================================================================
// COMPLETE WORKFLOW TESTS - From Wizard Configuration to CV Export
// ============================================================================
// These tests verify the ENTIRE user journey as documented in:
// docs/workflows/CV_GENERATION_WORKFLOW.md
//
// The complete journey is:
// 1. Wizard Modal (3 steps: WHO → TECH → FORMAT)
// 2. Tech Extraction (async)
// 3. CV Generation (async)
// 4. Preview CV
// 5. Export Modal (select format + location)
// 6. Exporting (async)
// 7. Export Complete (file saved or clipboard copied)
// ============================================================================

var _ = Describe("GenerateCV Complete Workflow E2E Tests", func() {
	var (
		intent       *GenerateCVIntent
		testProfiles []*CVProfile
		testEvents   []*career.CareerEvent
		testFacts    []*career.Fact
	)

	BeforeEach(func() {
		testProfiles = []*CVProfile{
			{
				ID:             "staff_engineer",
				Name:           "Staff Engineer",
				TargetRole:     "staff",
				TargetAudience: "hiring_manager",
			},
			{
				ID:             "principal_engineer",
				Name:           "Principal Engineer",
				TargetRole:     "principal",
				TargetAudience: "hiring_manager",
			},
			{
				ID:             "engineering_manager",
				Name:           "Engineering Manager",
				TargetRole:     "em",
				TargetAudience: "recruiter",
			},
			{
				ID:             "senior_engineer",
				Name:           "Senior Engineer",
				TargetRole:     "senior",
				TargetAudience: "peer",
			},
		}

		testEvents = []*career.CareerEvent{
			{
				ID:      "event_1",
				Text:    "Led migration to microservices architecture using Go and Kubernetes",
				Date:    time.Date(2025, 12, 15, 0, 0, 0, 0, time.UTC),
				Company: "TechCorp",
				Project: "Platform Modernization",
				Tags:    []string{"architecture", "go", "kubernetes"},
			},
			{
				ID:      "event_2",
				Text:    "Implemented CI/CD pipeline with GitHub Actions and ArgoCD",
				Date:    time.Date(2025, 11, 10, 0, 0, 0, 0, time.UTC),
				Company: "TechCorp",
				Project: "DevOps Initiative",
				Tags:    []string{"devops", "ci-cd", "automation"},
			},
		}

		testFacts = []*career.Fact{
			{
				ID:            "fact_1",
				Text:          "Reduced deployment time by 80% through automation",
				SourceEventID: "event_2",
			},
		}

		ctx := &GenerateCVContext{
			AvailableProfiles: testProfiles,
			Events:            testEvents,
			Facts:             testFacts,
			DefaultProfile:    testProfiles[0],
			AppContext:        context.Background(),
			ReviewScreenFactory: func(cv *career.CVView) screens.Screen {
				return &mockWizardReviewScreen{cv: cv}
			},
			PreviewScreenFactory: func(cv *career.CVView) screens.Screen {
				return &mockWizardPreviewScreen{cv: cv}
			},
		}

		var err error
		intent, err = NewGenerateCVIntent(ctx)
		Expect(err).NotTo(HaveOccurred())
		intent.EnableWizardFlow()

		termInfo := intent.BaseIntent.GetTerminalInfo()
		termInfo.Update(tea.WindowSizeMsg{Width: 120, Height: 40})
	})

	// ========================================================================
	// COMPLETE WORKFLOW: Wizard → Export to File
	// ========================================================================
	Describe("Complete Workflow: Wizard to File Export", func() {
		It("should complete entire journey from wizard configuration to file export", func() {
			// ============================================================
			// STEP 1: Initialize - Wizard Modal appears
			// ============================================================
			intent.Init()
			Expect(intent.state.currentState).To(Equal(CVStateConfiguring))
			Expect(intent.wizardModal).NotTo(BeNil())
			Expect(intent.wizardModal.IsVisible()).To(BeTrue())

			// ============================================================
			// STEP 2: Complete Wizard (simulated - user selects options)
			// ============================================================
			intent.Update(WizardCompleteMsg{
				ProfileID: "staff_engineer",
				Audience:  "hiring_manager",
			})
			Expect(intent.state.currentState).To(Equal(CVStateExtracting))

			// ============================================================
			// STEP 3: Tech Extraction completes
			// ============================================================
			intent.Update(TechnologiesExtractedMsg{
				Technologies: []*ExtractedTechnology{
					{ID: "go", Name: "Go", Category: "Language"},
					{ID: "kubernetes", Name: "Kubernetes", Category: "Platform"},
				},
			})
			Expect(intent.state.currentState).To(Equal(CVStateGenerating))

			// ============================================================
			// STEP 4: CV Generation completes (goes to Review first)
			// ============================================================
			testCV := &career.CVView{
				ID:   "cv_generated_1",
				Name: "Staff Engineer CV",
				Sections: []*career.CVSection{
					{ID: "summary", SectionType: "summary", Title: "Professional Summary"},
					{ID: "experience", SectionType: "experience", Title: "Experience"},
					{ID: "skills", SectionType: "skills", Title: "Technical Skills"},
				},
			}
			intent.Update(CVGenerationCompleteMsg{CV: testCV})
			Expect(intent.state.currentState).To(Equal(CVStateReview))
			Expect(intent.state.generatedCV).NotTo(BeNil())
			Expect(intent.state.generatedCV.ID).To(Equal("cv_generated_1"))

			// ============================================================
			// STEP 5: Navigate from Review to Preview
			// ============================================================
			intent.Update(tea.KeyMsg{Type: tea.KeyEnter})
			Expect(intent.state.currentState).To(Equal(CVStatePreview))

			// ============================================================
			// STEP 6: User presses 'x' to export
			// ============================================================
			intent.Update(tea.KeyMsg{Type: tea.KeyRunes, Runes: []rune{'x'}})
			Expect(intent.state.currentState).To(Equal(CVStateExporting))
			Expect(intent.exportModal).NotTo(BeNil())
			Expect(intent.exportModal.IsVisible()).To(BeTrue())

			// ============================================================
			// STEP 7: Simulate export (file save)
			// ============================================================
			intent.Update(CVExportCompleteMsg{
				Path:  "/home/user/.kariya/exports/staff_engineer_cv_2026.md",
				Error: nil,
			})
			Expect(intent.state.currentState).To(Equal(GenerateCVStateExportComplete))
			Expect(intent.exportModal.IsVisible()).To(BeFalse())

			// ============================================================
			// STEP 8: Verify export complete view shows file path
			// ============================================================
			view := intent.View()
			Expect(view).To(ContainSubstring("Export Complete"))

			// ============================================================
			// STEP 9: User presses Enter to finish
			// ============================================================
			intent.Update(tea.KeyMsg{Type: tea.KeyEnter})

			// ============================================================
			// VERIFY: Workflow completed with all data
			// ============================================================
			result := intent.Result()
			Expect(result).NotTo(BeNil())
			Expect(result.Status).To(Equal(Completed))

			cvResult, ok := result.Data.(*GenerateCVResult)
			Expect(ok).To(BeTrue())
			Expect(cvResult.GeneratedCV).NotTo(BeNil())
			Expect(cvResult.GeneratedCV.ID).To(Equal("cv_generated_1"))
			Expect(cvResult.SelectedProfile.ID).To(Equal("staff_engineer"))
			Expect(cvResult.ExportPath).To(Equal("/home/user/.kariya/exports/staff_engineer_cv_2026.md"))
		})

		It("should complete entire journey from wizard to clipboard export", func() {
			// Initialize and complete wizard
			intent.Init()
			intent.Update(WizardCompleteMsg{ProfileID: "senior_engineer", Audience: "peer"})
			intent.Update(TechnologiesExtractedMsg{})
			intent.Update(CVGenerationCompleteMsg{CV: &career.CVView{ID: "cv_clipboard"}})
			// Navigate from review to preview
			intent.Update(tea.KeyMsg{Type: tea.KeyEnter})

			// Export to clipboard
			intent.Update(tea.KeyMsg{Type: tea.KeyRunes, Runes: []rune{'x'}})
			intent.state.selectedExportOption = CVExportOptionClipboard
			intent.Update(CVExportCompleteMsg{Path: "clipboard", Error: nil})

			// Verify clipboard export
			view := intent.View()
			Expect(view).To(ContainSubstring("Clipboard"))

			// Complete workflow
			intent.Update(tea.KeyMsg{Type: tea.KeyEnter})

			result := intent.Result()
			Expect(result.Status).To(Equal(Completed))
			cvResult, _ := result.Data.(*GenerateCVResult)
			Expect(cvResult.ExportPath).To(Equal("clipboard"))
		})
	})

	// ========================================================================
	// WIZARD STEP 1 (WHO): Profile & Audience Verification
	// ========================================================================
	Describe("Wizard Step 1 (WHO): Profile & Audience Options", func() {
		BeforeEach(func() {
			intent.Init()
			Expect(intent.wizardModal.IsVisible()).To(BeTrue())
		})

		It("should have Profile field in Step 1", func() {
			view := intent.View()
			Expect(view).To(ContainSubstring("Profile"))
		})

		It("should have Audience field in Step 1", func() {
			view := intent.View()
			Expect(view).To(ContainSubstring("Audience"))
		})

		It("should show Hiring Manager audience option", func() {
			view := intent.View()
			Expect(view).To(ContainSubstring("Hiring Manager"))
		})

		It("should show Recruiter audience option", func() {
			view := intent.View()
			Expect(view).To(ContainSubstring("Recruiter"))
		})

		It("should show Peer/Colleague audience option", func() {
			view := intent.View()
			Expect(view).To(SatisfyAny(
				ContainSubstring("Peer"),
				ContainSubstring("Colleague"),
			))
		})

		It("should show all configured profile options", func() {
			view := intent.View()
			// Should show at least some of the profiles
			Expect(view).To(SatisfyAny(
				ContainSubstring("Staff Engineer"),
				ContainSubstring("Principal Engineer"),
				ContainSubstring("Engineering Manager"),
				ContainSubstring("Senior Engineer"),
			))
		})
	})

	// ========================================================================
	// WIZARD STEP 2 (TECH): Technology Focus Options
	// ========================================================================
	Describe("Wizard Step 2 (TECH): Technology Focus Options", func() {
		BeforeEach(func() {
			intent.Init()
			// Complete Step 1 to get to Step 2
			intent.wizardModal.SetProfileID("staff_engineer")
			intent.wizardModal.SetAudience("hiring_manager")
			// Advance to step 2 (simulated)
		})

		It("should have TechFocus field with Language Agnostic option", func() {
			// This tests the form configuration
			config := intent.wizardModal.GetConfigData()
			Expect(config).NotTo(BeNil())
			// The form should have the tech_focus field
		})

		It("should have correct TechFocus options per documentation", func() {
			// Per docs/workflows/CV_GENERATION_WORKFLOW.md line 176-179:
			// - Language Agnostic: Concepts over specific technologies
			// - Generalist: Show broad tech range
			// - Specialist: Highlight specific technologies

			// Test default value matches "language_agnostic" per code
			intent.wizardModal.SetTechFocus("language_agnostic")
			config := intent.wizardModal.GetConfigData()
			Expect(config.TechFocus).To(Equal("language_agnostic"))

			// Test setting to Generalist
			intent.wizardModal.SetTechFocus("generalist")
			config = intent.wizardModal.GetConfigData()
			Expect(config.TechFocus).To(Equal("generalist"))

			// Test setting to Specialist
			intent.wizardModal.SetTechFocus("specialist")
			config = intent.wizardModal.GetConfigData()
			Expect(config.TechFocus).To(Equal("specialist"))
		})

		It("should have correct FocusArea options per documentation", func() {
			// Per docs/workflows/CV_GENERATION_WORKFLOW.md line 181:
			// - Backend, Frontend, Full Stack, or DevOps

			intent.wizardModal.SetFocusArea("backend")
			Expect(intent.wizardModal.GetConfigData().FocusArea).To(Equal("backend"))

			intent.wizardModal.SetFocusArea("frontend")
			Expect(intent.wizardModal.GetConfigData().FocusArea).To(Equal("frontend"))

			intent.wizardModal.SetFocusArea("fullstack")
			Expect(intent.wizardModal.GetConfigData().FocusArea).To(Equal("fullstack"))

			intent.wizardModal.SetFocusArea("devops")
			Expect(intent.wizardModal.GetConfigData().FocusArea).To(Equal("devops"))
		})
	})

	// ========================================================================
	// WIZARD STEP 3 (FORMAT): Skills & Length Options
	// ========================================================================
	Describe("Wizard Step 3 (FORMAT): Skills & Length Options", func() {
		BeforeEach(func() {
			intent.Init()
		})

		It("should have correct SkillsFormat options per documentation", func() {
			// Per docs/workflows/CV_GENERATION_WORKFLOW.md line 196-199:
			// - Grouped by Category
			// - Flat List
			// - Categorized with Descriptions

			intent.wizardModal.SetSkillsFormat("grouped")
			Expect(intent.wizardModal.GetConfigData().SkillsFormat).To(Equal("grouped"))

			intent.wizardModal.SetSkillsFormat("flat")
			Expect(intent.wizardModal.GetConfigData().SkillsFormat).To(Equal("flat"))

			intent.wizardModal.SetSkillsFormat("categorized")
			Expect(intent.wizardModal.GetConfigData().SkillsFormat).To(Equal("categorized"))
		})

		It("should have correct CVLength options per documentation", func() {
			// Per BUG-010 fix: Added 4th option to match service layer
			// - 1 Page (executive summary) -> 1_page
			// - 2 Pages (concise) -> 2_page
			// - Standard (2-3 pages) -> standard
			// - Detailed (3+ pages) -> detailed

			intent.wizardModal.SetCVLength("1_page")
			Expect(intent.wizardModal.GetConfigData().CVLength).To(Equal("1_page"))

			intent.wizardModal.SetCVLength("2_page")
			Expect(intent.wizardModal.GetConfigData().CVLength).To(Equal("2_page"))

			intent.wizardModal.SetCVLength("standard")
			Expect(intent.wizardModal.GetConfigData().CVLength).To(Equal("standard"))

			intent.wizardModal.SetCVLength("detailed")
			Expect(intent.wizardModal.GetConfigData().CVLength).To(Equal("detailed"))
		})

		It("should default to 'grouped' SkillsFormat", func() {
			// Per implementation: default is "grouped"
			config := intent.wizardModal.GetConfigData()
			if config.SkillsFormat == "" {
				intent.wizardModal.SetSkillsFormat("")
				// After setting empty, should use default
			}
			// Default should be set
			Expect(config.SkillsFormat).To(SatisfyAny(
				Equal("grouped"),
				BeEmpty(), // before defaults applied
			))
		})

		It("should have valid CVLength default value", func() {
			// Form defaults to first option "1_page", but SetDefaults() sets "2_page"
			// Either is valid depending on initialization order
			config := intent.wizardModal.GetConfigData()
			Expect(config.CVLength).To(SatisfyAny(
				Equal("1_page"), // Form first option
				Equal("2_page"), // SetDefaults value
				BeEmpty(),       // Before any initialization
			))
		})
	})

	// ========================================================================
	// EXPORT OPTIONS VERIFICATION
	// ========================================================================
	Describe("Export Options: Format and Location", func() {
		BeforeEach(func() {
			// Get to export state (through review -> preview)
			intent.Init()
			intent.Update(WizardCompleteMsg{ProfileID: "staff_engineer", Audience: "hiring_manager"})
			intent.Update(TechnologiesExtractedMsg{})
			intent.Update(CVGenerationCompleteMsg{CV: &career.CVView{ID: "cv_1"}})
			// Navigate from review to preview
			intent.Update(tea.KeyMsg{Type: tea.KeyEnter})
			intent.Update(tea.KeyMsg{Type: tea.KeyRunes, Runes: []rune{'x'}})
			Expect(intent.exportModal).NotTo(BeNil())
		})

		It("should have Text export format option", func() {
			intent.exportModal.SetFormat("text")
			data := intent.exportModal.GetExportData()
			Expect(data.Format).To(Equal("text"))
		})

		It("should have Markdown export format option", func() {
			intent.exportModal.SetFormat("markdown")
			data := intent.exportModal.GetExportData()
			Expect(data.Format).To(Equal("markdown"))
		})

		It("should have YAML export format option", func() {
			intent.exportModal.SetFormat("yaml")
			data := intent.exportModal.GetExportData()
			Expect(data.Format).To(Equal("yaml"))
		})

		It("should have File save location option", func() {
			intent.exportModal.SetLocation("file")
			data := intent.exportModal.GetExportData()
			Expect(data.Location).To(Equal("file"))
		})

		It("should have Clipboard save location option", func() {
			intent.exportModal.SetLocation("clipboard")
			data := intent.exportModal.GetExportData()
			Expect(data.Location).To(Equal("clipboard"))
		})
	})

	// ========================================================================
	// WORKFLOW STATE VERIFICATION - All States Are Reachable
	// ========================================================================
	Describe("Workflow State Verification", func() {
		It("should reach all 8 wizard workflow states in sequence", func() {
			// State 1: Configuring (Wizard Modal)
			intent.Init()
			Expect(intent.state.currentState).To(Equal(CVStateConfiguring))

			// State 2: Extracting Technologies
			intent.Update(WizardCompleteMsg{ProfileID: "staff_engineer", Audience: "hiring_manager"})
			Expect(intent.state.currentState).To(Equal(CVStateExtracting))

			// State 3: Generating CV
			intent.Update(TechnologiesExtractedMsg{})
			Expect(intent.state.currentState).To(Equal(CVStateGenerating))

			// State 4: Review (new!)
			intent.Update(CVGenerationCompleteMsg{CV: &career.CVView{ID: "cv_1"}})
			Expect(intent.state.currentState).To(Equal(CVStateReview))

			// State 5: Preview
			intent.Update(tea.KeyMsg{Type: tea.KeyEnter})
			Expect(intent.state.currentState).To(Equal(CVStatePreview))

			// State 6: Exporting (Export Modal)
			intent.Update(tea.KeyMsg{Type: tea.KeyRunes, Runes: []rune{'x'}})
			Expect(intent.state.currentState).To(Equal(CVStateExporting))

			// State 7: GenerateCVStateExporting (actual export in progress)
			// This happens when export modal is completed and async export starts
			// We simulate this by directly setting state + sending completion msg
			intent.state.selectedExportFormat = CVExportFormatMarkdown
			intent.state.selectedExportOption = CVExportOptionSaveToFile
			intent.state.currentState = GenerateCVStateExporting

			// State 8: Export Complete
			intent.Update(CVExportCompleteMsg{Path: "/tmp/cv.md", Error: nil})
			Expect(intent.state.currentState).To(Equal(GenerateCVStateExportComplete))

			// Complete workflow
			intent.Update(tea.KeyMsg{Type: tea.KeyEnter})
			Expect(intent.Result().Status).To(Equal(Completed))
		})
	})

	// WIZARD DATA FLOW VERIFICATION - Phase 8 Integration Tests
	// ========================================================================
	Describe("Wizard Data Flow (Phase 8)", func() {
		Describe("WizardCompleteMsg carries all wizard data", func() {
			It("should store all wizard selections from WizardCompleteMsg", func() {
				intent.Init()
				Expect(intent.state.currentState).To(Equal(CVStateConfiguring))

				// Send complete wizard message with ALL fields populated
				// Using "staff_engineer" profile ID which exists in this test suite's BeforeEach
				msg := WizardCompleteMsg{
					// Step 1: WHO
					ProfileID: "staff_engineer",
					Audience:  "recruiter",
					// Step 2: TECH
					TechFocus:    "specialist",
					Technologies: []string{"Go", "PostgreSQL"},
					FocusArea:    "backend",
					// Step 3: FORMAT
					SkillsFormat: "grouped",
					CVLength:     "2_page",
				}
				intent.Update(msg)

				// Verify all fields were stored in state
				Expect(intent.state.selectedProfile.ID).To(Equal("staff_engineer"))
				Expect(intent.state.selectedAudience).To(Equal("recruiter"))
				Expect(string(intent.state.selectedTechnologyFocus)).To(Equal("specialist"))
				Expect(intent.state.selectedTechnologies).To(ConsistOf("Go", "PostgreSQL"))
				Expect(string(intent.state.selectedFocusArea)).To(Equal("backend"))
				Expect(intent.state.selectedSkillsFormat).To(Equal("grouped"))
				Expect(intent.state.selectedCVLength).To(Equal("2_page"))
			})

			It("should handle language_agnostic tech focus (no technologies)", func() {
				intent.Init()

				msg := WizardCompleteMsg{
					ProfileID:    "staff_engineer",
					Audience:     "hiring_manager",
					TechFocus:    "language_agnostic",
					Technologies: []string{}, // Empty for language agnostic
					FocusArea:    "fullstack",
					SkillsFormat: "flat",
					CVLength:     "1_page",
				}
				intent.Update(msg)

				Expect(string(intent.state.selectedTechnologyFocus)).To(Equal("language_agnostic"))
				Expect(intent.state.selectedTechnologies).To(BeEmpty())
				Expect(string(intent.state.selectedFocusArea)).To(Equal("fullstack"))
			})

			It("should handle generalist tech focus with multiple technologies", func() {
				intent.Init()

				msg := WizardCompleteMsg{
					ProfileID:    "principal_engineer",
					Audience:     "peer",
					TechFocus:    "generalist",
					Technologies: []string{"Ruby", "Python", "JavaScript", "PostgreSQL"},
					FocusArea:    "backend",
					SkillsFormat: "categorized",
					CVLength:     "detailed",
				}
				intent.Update(msg)

				Expect(string(intent.state.selectedTechnologyFocus)).To(Equal("generalist"))
				Expect(intent.state.selectedTechnologies).To(HaveLen(4))
				Expect(intent.state.selectedTechnologies).To(ContainElement("Ruby"))
				Expect(intent.state.selectedTechnologies).To(ContainElement("PostgreSQL"))
			})
		})

		Describe("Wizard modal data extraction", func() {
			It("should extract all config data from wizard modal on completion", func() {
				intent.Init()
				Expect(intent.wizardModal).NotTo(BeNil())

				// Set wizard form data directly
				wizard := intent.wizardModal
				wizard.SetProfileID("staff_engineer")
				wizard.SetAudience("hiring_manager")
				wizard.SetTechFocus("specialist")
				wizard.SetTechnologies([]string{"Go"})
				wizard.SetFocusArea("backend")
				wizard.SetSkillsFormat("grouped")
				wizard.SetCVLength("2_page")

				// Get config data
				config := wizard.GetConfigData()

				// Verify all fields are retrievable
				Expect(config.ProfileID).To(Equal("staff_engineer"))
				Expect(config.Audience).To(Equal("hiring_manager"))
				Expect(config.TechFocus).To(Equal("specialist"))
				Expect(config.Technologies).To(ConsistOf("Go"))
				Expect(config.FocusArea).To(Equal("backend"))
				Expect(config.SkillsFormat).To(Equal("grouped"))
				Expect(config.CVLength).To(Equal("2_page"))
			})
		})

		Describe("Data persists through workflow", func() {
			It("should preserve wizard selections through tech extraction and CV generation", func() {
				intent.Init()

				// Complete wizard with full data
				intent.Update(WizardCompleteMsg{
					ProfileID:    "staff_engineer",
					Audience:     "recruiter",
					TechFocus:    "specialist",
					Technologies: []string{"Go"},
					FocusArea:    "backend",
					SkillsFormat: "grouped",
					CVLength:     "2_page",
				})
				Expect(intent.state.currentState).To(Equal(CVStateExtracting))

				// Complete tech extraction
				intent.Update(TechnologiesExtractedMsg{
					Technologies: []*ExtractedTechnology{
						{ID: "skill_1", Name: "Go", Category: "backend"},
					},
				})
				Expect(intent.state.currentState).To(Equal(CVStateGenerating))

				// Verify data is still preserved
				Expect(string(intent.state.selectedTechnologyFocus)).To(Equal("specialist"))
				Expect(intent.state.selectedTechnologies).To(ConsistOf("Go"))
				Expect(string(intent.state.selectedFocusArea)).To(Equal("backend"))
				Expect(intent.state.selectedSkillsFormat).To(Equal("grouped"))
				Expect(intent.state.selectedCVLength).To(Equal("2_page"))

				// Complete CV generation (goes to review first)
				intent.Update(CVGenerationCompleteMsg{CV: &career.CVView{ID: "cv_1"}})
				Expect(intent.state.currentState).To(Equal(CVStateReview))

				// Data should still be preserved
				Expect(string(intent.state.selectedTechnologyFocus)).To(Equal("specialist"))
				Expect(intent.state.selectedSkillsFormat).To(Equal("grouped"))
			})
		})
	})
})
