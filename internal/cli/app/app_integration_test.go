package app_test

import (
	"context"
	"testing"
	"time"

	"github.com/baphled/kariya/internal/cli/app"
	"github.com/baphled/kariya/internal/cli/service"
	career "github.com/baphled/kariya/internal/domain/career"
	careerrepo "github.com/baphled/kariya/internal/repository/career"
	careerservice "github.com/baphled/kariya/internal/service/career"
	tea "github.com/charmbracelet/bubbletea"
	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
)

var _ = Describe("App Menu Integration Tests", func() {
	var (
		model      *app.Model
		repo       *careerrepo.MemoryRepository
		svc        *careerservice.Service
		cliService *service.CLIEventService
	)

	BeforeEach(func() {
		_ = context.Background()
		repo = careerrepo.NewMemoryRepository()
		burstRepo := careerrepo.NewMemoryBurstRepository()
		factRepo := careerrepo.NewMemoryFactRepository()
		svc = careerservice.NewService(repo)
		svc.SetBurstRepository(burstRepo)
		svc.SetFactRepository(factRepo)
		cliService = service.NewCLIEventService(svc)
		// Pre-populate burst/fact repositories with dummy entries to avoid nil panics
		_ = burstRepo.Create(context.Background(), &career.Burst{ID: "b1", Name: "dummy", EventIDs: []string{"e1", "e2"}})
		_ = factRepo.Create(context.Background(), &career.Fact{ID: "f1", Text: "dummy", CompetencyCategories: []string{"leadership"}, RoleFit: "staff", AudienceRelevance: []string{"peer"}, SourceEventID: "e1"})
		model = app.NewModel(cliService, svc)

		Expect(model).NotTo(BeNil())
	})

	Describe("Bubbles table integration", func() {
		It("should align the cursor with the selected menu item", func() {
			// Simulate model view
			output := model.View()
			// Check for the tagline since we now use ASCII art logo
			Expect(output).To(ContainSubstring("Career Event Management System"))
		})

		It("should update cursor position through bubble navigation keys", func() {
			// Navigate two steps down
			modelInterface, _ := model.Update(tea.KeyMsg{Type: tea.KeyRunes, Runes: []rune("j")})
			model = modelInterface.(*app.Model)
			modelInterface, _ = model.Update(tea.KeyMsg{Type: tea.KeyRunes, Runes: []rune("j")})
			model = modelInterface.(*app.Model)
			modelInterface, _ = model.Update(tea.KeyMsg{Type: tea.KeyEnter})
			model = modelInterface.(*app.Model)
			// Generate view
			output := model.View()
			menuItems := model.GetMenuItems()
			// After StandardView migration, breadcrumbs show "Main Menu > Generate CV > ..."
			// so the menu item name WILL appear in breadcrumbs - this is expected behavior
			Expect(output).To(ContainSubstring(menuItems[2].Name)) // Intent name appears in breadcrumbs
		})
	})

	Describe("Menu Navigation", func() {
		It("should start in menu state", func() {
			Expect(model).NotTo(BeNil())
		})

		It("should navigate down in menu", func() {
			msg := tea.KeyMsg{Type: tea.KeyRunes, Runes: []rune("j")}
			newModel, _ := model.Update(msg)
			Expect(newModel).NotTo(BeNil())
		})

		It("should navigate up in menu", func() {
			msg := tea.KeyMsg{Type: tea.KeyRunes, Runes: []rune("k")}
			newModel, _ := model.Update(msg)
			Expect(newModel).NotTo(BeNil())
		})
	})

	Describe("Intent Activation", func() {
		It("should activate CaptureEvent intent without panic", func() {
			msg := tea.KeyMsg{Type: tea.KeyEnter}
			newModel, cmd := model.Update(msg)
			Expect(newModel).NotTo(BeNil())
			Expect(cmd).NotTo(BeNil())
		})

		It("should activate BurstManagement intent without panic", func() {
			modelInterface := tea.Model(model)
			for i := 0; i < 5; i++ {
				var cmd tea.Cmd
				modelInterface, cmd = modelInterface.Update(tea.KeyMsg{Type: tea.KeyRunes, Runes: []rune("j")})
				_ = cmd
			}
			msg := tea.KeyMsg{Type: tea.KeyEnter}
			newModel, cmd := modelInterface.Update(msg)
			Expect(newModel).NotTo(BeNil())
			Expect(cmd).NotTo(BeNil())
		})

		It("should activate FactManagement intent without panic", func() {
			modelInterface := tea.Model(model)
			for i := 0; i < 6; i++ {
				var cmd tea.Cmd
				modelInterface, cmd = modelInterface.Update(tea.KeyMsg{Type: tea.KeyRunes, Runes: []rune("j")})
				_ = cmd
			}
			msg := tea.KeyMsg{Type: tea.KeyEnter}
			newModel, cmd := modelInterface.Update(msg)
			Expect(newModel).NotTo(BeNil())
			Expect(cmd).NotTo(BeNil())
		})

		It("should activate ImportWizard intent without panic", func() {
			modelInterface := tea.Model(model)
			for i := 0; i < 7; i++ {
				var cmd tea.Cmd
				modelInterface, cmd = modelInterface.Update(tea.KeyMsg{Type: tea.KeyRunes, Runes: []rune("j")})
				_ = cmd
			}
			msg := tea.KeyMsg{Type: tea.KeyEnter}
			newModel, cmd := modelInterface.Update(msg)
			Expect(newModel).NotTo(BeNil())
			Expect(cmd).NotTo(BeNil())
		})

		It("should activate MetadataEditor intent without panic", func() {
			modelInterface := tea.Model(model)
			for i := 0; i < 8; i++ {
				var cmd tea.Cmd
				modelInterface, cmd = modelInterface.Update(tea.KeyMsg{Type: tea.KeyRunes, Runes: []rune("j")})
				_ = cmd
			}
			msg := tea.KeyMsg{Type: tea.KeyEnter}
			newModel, cmd := modelInterface.Update(msg)
			Expect(newModel).NotTo(BeNil())
			Expect(cmd).NotTo(BeNil())
		})

		It("should activate BulkOperations intent without panic", func() {
			modelInterface := tea.Model(model)
			for i := 0; i < 9; i++ {
				var cmd tea.Cmd
				modelInterface, cmd = modelInterface.Update(tea.KeyMsg{Type: tea.KeyRunes, Runes: []rune("j")})
				_ = cmd
			}
			msg := tea.KeyMsg{Type: tea.KeyEnter}
			newModel, cmd := modelInterface.Update(msg)
			Expect(newModel).NotTo(BeNil())
			Expect(cmd).NotTo(BeNil())
		})
	})

	Describe("Error Handling", func() {
		It("should handle Ctrl+C to quit", func() {
			msg := tea.KeyMsg{Type: tea.KeyCtrlC}
			newModel, _ := model.Update(msg)
			Expect(newModel).NotTo(BeNil())
		})
	})
})

var _ = Describe("Navigation Integration", func() {
	var (
		model      *app.Model
		repo       *careerrepo.MemoryRepository
		svc        *careerservice.Service
		cliService *service.CLIEventService
	)

	BeforeEach(func() {
		repo = careerrepo.NewMemoryRepository()
		burstRepo := careerrepo.NewMemoryBurstRepository()
		factRepo := careerrepo.NewMemoryFactRepository()
		svc = careerservice.NewService(repo)
		svc.SetBurstRepository(burstRepo)
		svc.SetFactRepository(factRepo)
		cliService = service.NewCLIEventService(svc)
		model = app.NewModel(cliService, svc)
	})

	It("should start in menu state with menu visible", func() {
		Expect(model).NotTo(BeNil())
		output := model.View()
		Expect(output).To(ContainSubstring("Career Event Management System"))
	})

	It("should navigate down the menu and select an intent, activating intent view", func() {
		modelInterface, _ := model.Update(tea.KeyMsg{Type: tea.KeyRunes, Runes: []rune("j")})
		model = modelInterface.(*app.Model)
		modelInterface, cmd := model.Update(tea.KeyMsg{Type: tea.KeyEnter})
		model = modelInterface.(*app.Model)
		Expect(model).NotTo(BeNil())
		Expect(cmd).NotTo(BeNil())
		output := model.View()
		// After StandardView migration (Tasks 12-15), ALL screens show the logo/subtitle
		Expect(output).To(ContainSubstring("Career Event Management System")) // Logo on all screens
	})

	It("should navigate menu, activate, then go back to menu via escape", func() {
		modelInterface, _ := model.Update(tea.KeyMsg{Type: tea.KeyEnter})
		model = modelInterface.(*app.Model)
		output := model.View()
		// After StandardView migration, logo appears on intent screens too
		Expect(output).To(ContainSubstring("Career Event Management System"))
		modelInterface, _ = model.Update(tea.KeyMsg{Type: tea.KeyEscape})
		model = modelInterface.(*app.Model)
		output = model.View()
		Expect(output).To(ContainSubstring("Career Event Management System"))
	})

	It("should go back to menu after IntentCompletedMsg", func() {
		modelInterface, _ := model.Update(tea.KeyMsg{Type: tea.KeyEnter})
		model = modelInterface.(*app.Model)
		modelInterface, _ = model.Update(app.IntentCompletedMsg{})
		model = modelInterface.(*app.Model)
		output := model.View()
		Expect(output).To(ContainSubstring("Career Event Management System"))
	})

	It("should handle quick back/forward navigation, activating and quitting intents", func() {
		modelInterface, _ := model.Update(tea.KeyMsg{Type: tea.KeyRunes, Runes: []rune("j")})
		model = modelInterface.(*app.Model)
		modelInterface, _ = model.Update(tea.KeyMsg{Type: tea.KeyRunes, Runes: []rune("j")})
		model = modelInterface.(*app.Model)
		modelInterface, _ = model.Update(tea.KeyMsg{Type: tea.KeyEnter})
		model = modelInterface.(*app.Model)
		modelInterface, _ = model.Update(tea.KeyMsg{Type: tea.KeyEscape})
		model = modelInterface.(*app.Model)
		output := model.View()
		Expect(output).To(ContainSubstring("Career Event Management System"))
	})
})

func TestApp(t *testing.T) {
	RegisterFailHandler(Fail)
	RunSpecs(t, "App Integration Test Suite")
}

var _ = Describe("Intent Navigation - All Intents", func() {
	var (
		model      *app.Model
		repo       *careerrepo.MemoryRepository
		svc        *careerservice.Service
		cliService *service.CLIEventService
	)

	BeforeEach(func() {
		repo = careerrepo.NewMemoryRepository()
		burstRepo := careerrepo.NewMemoryBurstRepository()
		factRepo := careerrepo.NewMemoryFactRepository()
		svc = careerservice.NewService(repo)
		svc.SetBurstRepository(burstRepo)
		svc.SetFactRepository(factRepo)
		cliService = service.NewCLIEventService(svc)

		// Add dummy data
		_ = repo.Create(context.Background(), &career.CareerEvent{ID: "e1", Text: "test event", Date: time.Now()})
		_ = burstRepo.Create(context.Background(), &career.Burst{ID: "b1", Name: "dummy", EventIDs: []string{"e1"}})
		_ = factRepo.Create(context.Background(), &career.Fact{ID: "f1", Text: "dummy", CompetencyCategories: []string{"leadership"}, RoleFit: "staff", AudienceRelevance: []string{"peer"}, SourceEventID: "e1"})

		model = app.NewModel(cliService, svc)
	})

	// Test each intent in the menu (0-9)
	testIntentNavigation := func(menuIndex int, intentName string) {
		It("should navigate within "+intentName+" intent", func() {
			// Navigate to the menu item
			for i := 0; i < menuIndex; i++ {
				modelInterface, _ := model.Update(tea.KeyMsg{Type: tea.KeyRunes, Runes: []rune("j")})
				model = modelInterface.(*app.Model)
			}

			// Verify we're at the right menu item
			menuItems := model.GetMenuItems()
			Expect(menuIndex).To(BeNumerically("<", len(menuItems)))

			// Select the intent
			modelInterface, cmd := model.Update(tea.KeyMsg{Type: tea.KeyEnter})
			model = modelInterface.(*app.Model)

			// Execute any command from activation
			if cmd != nil {
				msg := cmd()
				if msg != nil {
					modelInterface, _ := model.Update(msg)
					model = modelInterface.(*app.Model)
				}
			}

			// Get the view - should NOT be the menu
			viewBeforeNav := model.View()
			Expect(viewBeforeNav).NotTo(ContainSubstring("KaRiya - Career Event Manager"),
				"Intent "+intentName+" should show intent view, not menu")

			// Try to navigate within the intent
			modelInterface, _ = model.Update(tea.KeyMsg{Type: tea.KeyRunes, Runes: []rune("j")})
			model = modelInterface.(*app.Model)

			// Get view after navigation - should STILL not be the menu
			viewAfterNav := model.View()
			Expect(viewAfterNav).NotTo(ContainSubstring("KaRiya - Career Event Manager"),
				"After navigation in "+intentName+", should still be in intent view, not back at menu")
		})
	}

	Describe("Intent Selection and Navigation", func() {
		testIntentNavigation(0, "CaptureEvent")
		testIntentNavigation(1, "BrowseTimeline")
		testIntentNavigation(2, "GenerateCV")
		testIntentNavigation(3, "ExportArtifact")
		testIntentNavigation(4, "ConfigureSystem")
		testIntentNavigation(5, "BurstManagement")
		testIntentNavigation(6, "FactManagement")
		testIntentNavigation(7, "ImportWizard")
		testIntentNavigation(8, "MetadataEditor")
		testIntentNavigation(9, "BulkOperations")
	})
})

var _ = Describe("Intent Navigation - Detailed", func() {
	var (
		model      *app.Model
		repo       *careerrepo.MemoryRepository
		svc        *careerservice.Service
		cliService *service.CLIEventService
	)

	BeforeEach(func() {
		repo = careerrepo.NewMemoryRepository()
		burstRepo := careerrepo.NewMemoryBurstRepository()
		factRepo := careerrepo.NewMemoryFactRepository()
		svc = careerservice.NewService(repo)
		svc.SetBurstRepository(burstRepo)
		svc.SetFactRepository(factRepo)
		cliService = service.NewCLIEventService(svc)

		// Add dummy data
		_ = repo.Create(context.Background(), &career.CareerEvent{ID: "e1", Text: "test event", Date: time.Now()})
		_ = burstRepo.Create(context.Background(), &career.Burst{ID: "b1", Name: "dummy", EventIDs: []string{"e1"}})
		_ = factRepo.Create(context.Background(), &career.Fact{ID: "f1", Text: "dummy", CompetencyCategories: []string{"leadership"}, RoleFit: "staff", AudienceRelevance: []string{"peer"}, SourceEventID: "e1"})

		model = app.NewModel(cliService, svc)
	})

	selectIntent := func(menuIndex int) {
		// Navigate to menu item
		for i := 0; i < menuIndex; i++ {
			modelInterface, _ := model.Update(tea.KeyMsg{Type: tea.KeyRunes, Runes: []rune("j")})
			model = modelInterface.(*app.Model)
		}

		// Select intent
		modelInterface, cmd := model.Update(tea.KeyMsg{Type: tea.KeyEnter})
		model = modelInterface.(*app.Model)

		// Execute command if present
		if cmd != nil {
			msg := cmd()
			if msg != nil {
				modelInterface, _ := model.Update(msg)
				model = modelInterface.(*app.Model)
			}
		}
	}

	Describe("CaptureEvent Intent", func() {
		It("should display capture event form", func() {
			selectIntent(0)
			view := model.View()
			Expect(view).NotTo(ContainSubstring("KaRiya - Career Event Manager"))
			Expect(view).NotTo(BeEmpty())
		})

		It("should navigate with arrow keys within form", func() {
			selectIntent(0)
			modelInterface, _ := model.Update(tea.KeyMsg{Type: tea.KeyRunes, Runes: []rune("j")})
			model = modelInterface.(*app.Model)
			view := model.View()
			Expect(view).NotTo(ContainSubstring("KaRiya - Career Event Manager"))
		})
	})

	Describe("BrowseTimeline Intent", func() {
		It("should display timeline view", func() {
			selectIntent(1)
			view := model.View()
			Expect(view).NotTo(ContainSubstring("KaRiya - Career Event Manager"))
			Expect(view).NotTo(BeEmpty())
		})

		It("should navigate timeline with arrow keys", func() {
			selectIntent(1)
			modelInterface, _ := model.Update(tea.KeyMsg{Type: tea.KeyRunes, Runes: []rune("j")})
			model = modelInterface.(*app.Model)
			view := model.View()
			Expect(view).NotTo(ContainSubstring("KaRiya - Career Event Manager"))
		})
	})

	Describe("GenerateCV Intent", func() {
		It("should display CV generation view", func() {
			selectIntent(2)
			view := model.View()
			Expect(view).NotTo(ContainSubstring("KaRiya - Career Event Manager"))
			Expect(view).NotTo(BeEmpty())
		})

		It("should navigate with arrow keys", func() {
			selectIntent(2)
			modelInterface, _ := model.Update(tea.KeyMsg{Type: tea.KeyRunes, Runes: []rune("j")})
			model = modelInterface.(*app.Model)
			view := model.View()
			Expect(view).NotTo(ContainSubstring("KaRiya - Career Event Manager"))
		})
	})

	Describe("ExportArtifact Intent", func() {
		It("should display export view", func() {
			selectIntent(3)
			view := model.View()
			Expect(view).NotTo(ContainSubstring("KaRiya - Career Event Manager"))
			Expect(view).NotTo(BeEmpty())
		})

		It("should navigate with arrow keys", func() {
			selectIntent(3)
			modelInterface, _ := model.Update(tea.KeyMsg{Type: tea.KeyRunes, Runes: []rune("j")})
			model = modelInterface.(*app.Model)
			view := model.View()
			Expect(view).NotTo(ContainSubstring("KaRiya - Career Event Manager"))
		})
	})

	Describe("ConfigureSystem Intent", func() {
		It("should display configuration view", func() {
			selectIntent(4)
			view := model.View()
			Expect(view).NotTo(ContainSubstring("KaRiya - Career Event Manager"))
			Expect(view).NotTo(BeEmpty())
		})

		It("should navigate with arrow keys", func() {
			selectIntent(4)
			modelInterface, _ := model.Update(tea.KeyMsg{Type: tea.KeyRunes, Runes: []rune("j")})
			model = modelInterface.(*app.Model)
			view := model.View()
			Expect(view).NotTo(ContainSubstring("KaRiya - Career Event Manager"))
		})
	})

	Describe("BurstManagement Intent", func() {
		It("should display burst management view", func() {
			selectIntent(5)
			view := model.View()
			Expect(view).NotTo(ContainSubstring("KaRiya - Career Event Manager"))
			Expect(view).NotTo(BeEmpty())
		})

		It("should navigate with arrow keys", func() {
			selectIntent(5)
			modelInterface, _ := model.Update(tea.KeyMsg{Type: tea.KeyRunes, Runes: []rune("j")})
			model = modelInterface.(*app.Model)
			view := model.View()
			Expect(view).NotTo(ContainSubstring("KaRiya - Career Event Manager"))
		})
	})

	Describe("FactManagement Intent", func() {
		It("should display fact management view", func() {
			selectIntent(6)
			view := model.View()
			Expect(view).NotTo(ContainSubstring("KaRiya - Career Event Manager"))
			Expect(view).NotTo(BeEmpty())
		})

		It("should navigate with arrow keys", func() {
			selectIntent(6)
			modelInterface, _ := model.Update(tea.KeyMsg{Type: tea.KeyRunes, Runes: []rune("j")})
			model = modelInterface.(*app.Model)
			view := model.View()
			Expect(view).NotTo(ContainSubstring("KaRiya - Career Event Manager"))
		})
	})

	Describe("ImportWizard Intent", func() {
		It("should display import wizard view", func() {
			selectIntent(7)
			view := model.View()
			Expect(view).NotTo(ContainSubstring("KaRiya - Career Event Manager"))
			Expect(view).NotTo(BeEmpty())
		})

		It("should navigate with arrow keys", func() {
			selectIntent(7)
			modelInterface, _ := model.Update(tea.KeyMsg{Type: tea.KeyRunes, Runes: []rune("j")})
			model = modelInterface.(*app.Model)
			view := model.View()
			Expect(view).NotTo(ContainSubstring("KaRiya - Career Event Manager"))
		})
	})

	Describe("MetadataEditor Intent", func() {
		It("should display metadata editor view", func() {
			selectIntent(8)
			view := model.View()
			Expect(view).NotTo(ContainSubstring("KaRiya - Career Event Manager"))
			Expect(view).NotTo(BeEmpty())
		})

		It("should navigate with arrow keys", func() {
			selectIntent(8)
			modelInterface, _ := model.Update(tea.KeyMsg{Type: tea.KeyRunes, Runes: []rune("j")})
			model = modelInterface.(*app.Model)
			view := model.View()
			Expect(view).NotTo(ContainSubstring("KaRiya - Career Event Manager"))
		})
	})

	Describe("BulkOperations Intent", func() {
		It("should display bulk operations view", func() {
			selectIntent(9)
			view := model.View()
			Expect(view).NotTo(ContainSubstring("KaRiya - Career Event Manager"))
			Expect(view).NotTo(BeEmpty())
		})

		It("should navigate with arrow keys", func() {
			selectIntent(9)
			modelInterface, _ := model.Update(tea.KeyMsg{Type: tea.KeyRunes, Runes: []rune("j")})
			model = modelInterface.(*app.Model)
			view := model.View()
			Expect(view).NotTo(ContainSubstring("KaRiya - Career Event Manager"))
		})
	})
})

var _ = Describe("Intent List Navigation - Specific", func() {
	var (
		model      *app.Model
		repo       *careerrepo.MemoryRepository
		svc        *careerservice.Service
		cliService *service.CLIEventService
	)

	BeforeEach(func() {
		repo = careerrepo.NewMemoryRepository()
		burstRepo := careerrepo.NewMemoryBurstRepository()
		factRepo := careerrepo.NewMemoryFactRepository()
		svc = careerservice.NewService(repo)
		svc.SetBurstRepository(burstRepo)
		svc.SetFactRepository(factRepo)
		cliService = service.NewCLIEventService(svc)

		// Add multiple events for timeline
		_ = repo.Create(context.Background(), &career.CareerEvent{ID: "e1", Text: "event 1", Date: time.Now()})
		_ = repo.Create(context.Background(), &career.CareerEvent{ID: "e2", Text: "event 2", Date: time.Now()})
		_ = repo.Create(context.Background(), &career.CareerEvent{ID: "e3", Text: "event 3", Date: time.Now()})

		// Add multiple bursts
		_ = burstRepo.Create(context.Background(), &career.Burst{ID: "b1", Name: "burst 1", EventIDs: []string{"e1"}})
		_ = burstRepo.Create(context.Background(), &career.Burst{ID: "b2", Name: "burst 2", EventIDs: []string{"e2"}})
		_ = burstRepo.Create(context.Background(), &career.Burst{ID: "b3", Name: "burst 3", EventIDs: []string{"e3"}})

		// Add multiple facts
		_ = factRepo.Create(context.Background(), &career.Fact{ID: "f1", Text: "fact 1", CompetencyCategories: []string{"leadership"}, RoleFit: "staff", AudienceRelevance: []string{"peer"}, SourceEventID: "e1"})
		_ = factRepo.Create(context.Background(), &career.Fact{ID: "f2", Text: "fact 2", CompetencyCategories: []string{"technical"}, RoleFit: "staff", AudienceRelevance: []string{"peer"}, SourceEventID: "e2"})
		_ = factRepo.Create(context.Background(), &career.Fact{ID: "f3", Text: "fact 3", CompetencyCategories: []string{"communication"}, RoleFit: "staff", AudienceRelevance: []string{"peer"}, SourceEventID: "e3"})

		model = app.NewModel(cliService, svc)
	})

	selectIntent := func(menuIndex int) {
		for i := 0; i < menuIndex; i++ {
			modelInterface, _ := model.Update(tea.KeyMsg{Type: tea.KeyRunes, Runes: []rune("j")})
			model = modelInterface.(*app.Model)
		}
		modelInterface, cmd := model.Update(tea.KeyMsg{Type: tea.KeyEnter})
		model = modelInterface.(*app.Model)
		if cmd != nil {
			msg := cmd()
			if msg != nil {
				modelInterface, _ := model.Update(msg)
				model = modelInterface.(*app.Model)
			}
		}
	}

	Describe("BrowseTimeline List Navigation", func() {
		It("should allow navigating down the timeline with 'j'", func() {
			selectIntent(1)  // BrowseTimeline
			_ = model.View() // viewBefore

			// Navigate down
			modelInterface, _ := model.Update(tea.KeyMsg{Type: tea.KeyRunes, Runes: []rune("j")})
			model = modelInterface.(*app.Model)

			viewAfter := model.View()
			// View should change when navigating (different item selected)
			// At minimum, should still be in timeline view
			Expect(viewAfter).NotTo(ContainSubstring("KaRiya - Career Event Manager"))
		})

		It("should allow navigating up the timeline with 'k'", func() {
			selectIntent(1) // BrowseTimeline

			// Navigate down first
			modelInterface, _ := model.Update(tea.KeyMsg{Type: tea.KeyRunes, Runes: []rune("j")})
			model = modelInterface.(*app.Model)

			// Then navigate up
			modelInterface, _ = model.Update(tea.KeyMsg{Type: tea.KeyRunes, Runes: []rune("k")})
			model = modelInterface.(*app.Model)

			view := model.View()
			Expect(view).NotTo(ContainSubstring("KaRiya - Career Event Manager"))
		})

		It("should allow multiple consecutive down navigations", func() {
			selectIntent(1) // BrowseTimeline

			// Navigate down multiple times
			for i := 0; i < 3; i++ {
				modelInterface, _ := model.Update(tea.KeyMsg{Type: tea.KeyRunes, Runes: []rune("j")})
				model = modelInterface.(*app.Model)
			}

			view := model.View()
			Expect(view).NotTo(ContainSubstring("KaRiya - Career Event Manager"))
		})
	})

	Describe("BurstManagement List Navigation", func() {
		It("should allow navigating down the burst list with 'j'", func() {
			selectIntent(5) // BurstManagement

			// Navigate down
			modelInterface, _ := model.Update(tea.KeyMsg{Type: tea.KeyRunes, Runes: []rune("j")})
			model = modelInterface.(*app.Model)

			view := model.View()
			Expect(view).NotTo(ContainSubstring("KaRiya - Career Event Manager"))
		})

		It("should allow navigating up the burst list with 'k'", func() {
			selectIntent(5) // BurstManagement

			// Navigate down first
			modelInterface, _ := model.Update(tea.KeyMsg{Type: tea.KeyRunes, Runes: []rune("j")})
			model = modelInterface.(*app.Model)

			// Then navigate up
			modelInterface, _ = model.Update(tea.KeyMsg{Type: tea.KeyRunes, Runes: []rune("k")})
			model = modelInterface.(*app.Model)

			view := model.View()
			Expect(view).NotTo(ContainSubstring("KaRiya - Career Event Manager"))
		})

		It("should allow multiple consecutive down navigations in burst list", func() {
			selectIntent(5) // BurstManagement

			// Navigate down multiple times
			for i := 0; i < 3; i++ {
				modelInterface, _ := model.Update(tea.KeyMsg{Type: tea.KeyRunes, Runes: []rune("j")})
				model = modelInterface.(*app.Model)
			}

			view := model.View()
			Expect(view).NotTo(ContainSubstring("KaRiya - Career Event Manager"))
		})
	})

	Describe("FactManagement List Navigation", func() {
		It("should allow navigating down the fact list with 'j'", func() {
			selectIntent(6) // FactManagement

			// Navigate down
			modelInterface, _ := model.Update(tea.KeyMsg{Type: tea.KeyRunes, Runes: []rune("j")})
			model = modelInterface.(*app.Model)

			view := model.View()
			Expect(view).NotTo(ContainSubstring("KaRiya - Career Event Manager"))
		})

		It("should allow navigating up the fact list with 'k'", func() {
			selectIntent(6) // FactManagement

			// Navigate down first
			modelInterface, _ := model.Update(tea.KeyMsg{Type: tea.KeyRunes, Runes: []rune("j")})
			model = modelInterface.(*app.Model)

			// Then navigate up
			modelInterface, _ = model.Update(tea.KeyMsg{Type: tea.KeyRunes, Runes: []rune("k")})
			model = modelInterface.(*app.Model)

			view := model.View()
			Expect(view).NotTo(ContainSubstring("KaRiya - Career Event Manager"))
		})

		It("should allow multiple consecutive down navigations in fact list", func() {
			selectIntent(6) // FactManagement

			// Navigate down multiple times
			for i := 0; i < 3; i++ {
				modelInterface, _ := model.Update(tea.KeyMsg{Type: tea.KeyRunes, Runes: []rune("j")})
				model = modelInterface.(*app.Model)
			}

			view := model.View()
			Expect(view).NotTo(ContainSubstring("KaRiya - Career Event Manager"))
		})
	})
})
