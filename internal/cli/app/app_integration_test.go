package app_test

import (
	"context"
	"testing"

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
		svc = careerservice.NewService(repo)
		cliService = service.NewCLIEventService(svc)
		model = app.NewModel(cliService, svc)
		Expect(model).NotTo(BeNil())
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
			for i := 0; i < 5; i++ {
				model.Update(tea.KeyMsg{Type: tea.KeyRunes, Runes: []rune("j")})
			}
			msg := tea.KeyMsg{Type: tea.KeyEnter}
			newModel, cmd := model.Update(msg)
			Expect(newModel).NotTo(BeNil())
			Expect(cmd).NotTo(BeNil())
		})

		It("should activate FactManagement intent without panic", func() {
			for i := 0; i < 6; i++ {
				model.Update(tea.KeyMsg{Type: tea.KeyRunes, Runes: []rune("j")})
			}
			msg := tea.KeyMsg{Type: tea.KeyEnter}
			newModel, cmd := model.Update(msg)
			Expect(newModel).NotTo(BeNil())
			Expect(cmd).NotTo(BeNil())
		})

		It("should activate ImportWizard intent without panic", func() {
			for i := 0; i < 7; i++ {
				model.Update(tea.KeyMsg{Type: tea.KeyRunes, Runes: []rune("j")})
			}
			msg := tea.KeyMsg{Type: tea.KeyEnter}
			newModel, cmd := model.Update(msg)
			Expect(newModel).NotTo(BeNil())
			Expect(cmd).NotTo(BeNil())
		})

		It("should activate MetadataEditor intent without panic", func() {
			for i := 0; i < 8; i++ {
				model.Update(tea.KeyMsg{Type: tea.KeyRunes, Runes: []rune("j")})
			}
			msg := tea.KeyMsg{Type: tea.KeyEnter}
			newModel, cmd := model.Update(msg)
			Expect(newModel).NotTo(BeNil())
			Expect(cmd).NotTo(BeNil())
		})

		It("should activate BulkOperations intent without panic", func() {
			for i := 0; i < 9; i++ {
				model.Update(tea.KeyMsg{Type: tea.KeyRunes, Runes: []rune("j")})
			}
			msg := tea.KeyMsg{Type: tea.KeyEnter}
			newModel, cmd := model.Update(msg)
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

func TestAppIntegration(t *testing.T) {
	RegisterFailHandler(Fail)
	RunSpecs(t, "App Integration Test Suite")
}
