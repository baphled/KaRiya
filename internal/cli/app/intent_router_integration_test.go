package app

import (
	"github.com/baphled/kariya/internal/cli/intents"
	"github.com/baphled/kariya/internal/cli/service"
	careerrepo "github.com/baphled/kariya/internal/repository/career"
	careerservice "github.com/baphled/kariya/internal/service/career"
	tea "github.com/charmbracelet/bubbletea"
	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
)

var _ = Describe("IntentRouter Integration", func() {
	var (
		repo       *careerrepo.MemoryRepository
		cliService *service.CLIEventService
		svc        *careerservice.Service
		model      *Model
	)

	BeforeEach(func() {
		repo = careerrepo.NewMemoryRepository()
		svc = careerservice.NewService(repo)
		cliService = service.NewCLIEventService(svc)
		model = NewModel(cliService, svc)
	})

	Describe("IntentRouter initialization", func() {
		It("should initialize IntentRouter in NewModel", func() {
			Expect(model.intentRouter).NotTo(BeNil())
		})

		It("should have inIntentMode set to false by default", func() {
			Expect(model.inIntentMode).To(BeFalse())
		})

		It("should be of type IntentRouter interface", func() {
			var _ intents.IntentRouter = model.intentRouter
		})
	})

	Describe("deactivateIntent method", func() {
		It("should set inIntentMode to false", func() {
			model.inIntentMode = true
			model.deactivateIntent()
			Expect(model.inIntentMode).To(BeFalse())
		})

		It("should reset to HomeScreen", func() {
			model.currentScreen = CaptureScreen
			model.deactivateIntent()
			Expect(model.currentScreen).To(Equal(HomeScreen))
		})

		It("should update breadcrumbs", func() {
			model.breadcrumbs = []string{"Home", "Capture"}
			model.deactivateIntent()
			Expect(model.breadcrumbs).To(Equal([]string{"Home"}))
		})
	})

	Describe("handleIntentMessage method", func() {
		It("should return nil when not in intent mode", func() {
			model.inIntentMode = false
			resultModel, cmd := model.handleIntentMessage(tea.KeyMsg{})
			Expect(resultModel).To(Equal(model))
			Expect(cmd).To(BeNil())
		})

		It("should return nil when intent router is nil", func() {
			model.inIntentMode = true
			model.intentRouter = nil
			resultModel, cmd := model.handleIntentMessage(tea.KeyMsg{})
			Expect(resultModel).To(Equal(model))
			Expect(cmd).To(BeNil())
		})
	})

	Describe("Model fields", func() {
		It("should have intentRouter field", func() {
			Expect(model.intentRouter).NotTo(BeNil())
		})

		It("should have inIntentMode field", func() {
			Expect(model.inIntentMode).To(BeFalse())
		})

		It("should maintain other model fields", func() {
			Expect(model.cliService).NotTo(BeNil())
			Expect(model.service).NotTo(BeNil())
			Expect(model.currentScreen).To(Equal(HomeScreen))
		})
	})

	Describe("Integration with existing functionality", func() {
		It("should not affect existing screen-based navigation", func() {
			model.currentScreen = HomeScreen
			Expect(model.currentScreen).To(Equal(HomeScreen))

			model.currentScreen = CaptureScreen
			Expect(model.currentScreen).To(Equal(CaptureScreen))
		})

		It("should allow switching between screen mode and intent mode", func() {
			Expect(model.inIntentMode).To(BeFalse())

			model.inIntentMode = true
			Expect(model.inIntentMode).To(BeTrue())

			model.deactivateIntent()
			Expect(model.inIntentMode).To(BeFalse())
		})
	})
})
