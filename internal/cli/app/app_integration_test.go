package app_test

import (
	"context"
	"testing"

	career "github.com/baphled/kariya/internal/domain/career"
	"github.com/baphled/kariya/internal/cli/app"
	"github.com/baphled/kariya/internal/cli/service"
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
			Expect(output).To(ContainSubstring("KaRiya - Career Event Manager"))
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
			Expect(output).NotTo(ContainSubstring(menuItems[2].Name)) // Should show intent view, not menu
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
		Expect(output).To(ContainSubstring("KaRiya - Career Event Manager"))
	})

	It("should navigate down the menu and select an intent, activating intent view", func() {
		modelInterface, _ := model.Update(tea.KeyMsg{Type: tea.KeyRunes, Runes: []rune("j")})
		model = modelInterface.(*app.Model)
		modelInterface, cmd := model.Update(tea.KeyMsg{Type: tea.KeyEnter})
		model = modelInterface.(*app.Model)
		Expect(model).NotTo(BeNil())
		Expect(cmd).NotTo(BeNil())
		output := model.View()
		Expect(output).NotTo(ContainSubstring("KaRiya - Career Event Manager")) // View switches
	})

	It("should navigate menu, activate, then go back to menu via escape", func() {
		modelInterface, _ := model.Update(tea.KeyMsg{Type: tea.KeyEnter})
		model = modelInterface.(*app.Model)
		output := model.View()
		Expect(output).NotTo(ContainSubstring("KaRiya - Career Event Manager"))
		modelInterface, _ = model.Update(tea.KeyMsg{Type: tea.KeyEscape})
		model = modelInterface.(*app.Model)
		output = model.View()
		Expect(output).To(ContainSubstring("KaRiya - Career Event Manager"))
	})

	It("should go back to menu after IntentCompletedMsg", func() {
		modelInterface, _ := model.Update(tea.KeyMsg{Type: tea.KeyEnter})
		model = modelInterface.(*app.Model)
		modelInterface, _ = model.Update(app.IntentCompletedMsg{})
		model = modelInterface.(*app.Model)
		output := model.View()
		Expect(output).To(ContainSubstring("KaRiya - Career Event Manager"))
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
		Expect(output).To(ContainSubstring("KaRiya - Career Event Manager"))
	})
})

func TestApp(t *testing.T) {
	RegisterFailHandler(Fail)
	RunSpecs(t, "App Integration Test Suite")
}
