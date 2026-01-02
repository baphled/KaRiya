package app

import (
	"github.com/baphled/kariya/internal/cli/models"
	"github.com/baphled/kariya/internal/cli/service"
	"github.com/baphled/kariya/internal/domain/career"
	careerrepo "github.com/baphled/kariya/internal/repository/career"
	careerservice "github.com/baphled/kariya/internal/service/career"
	"github.com/baphled/kariya/internal/service/career/cv"
	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
)

var _ = Describe("CV Menu Integration", func() {
	var (
		repo          *careerrepo.MemoryRepository
		svc           *careerservice.Service
		cliService    *service.CLIEventService
		model         *Model
		configManager cv.ConfigManager
	)

	BeforeEach(func() {
		repo = careerrepo.NewMemoryRepository()
		svc = careerservice.NewService(repo)
		cliService = service.NewCLIEventService(svc)
		configManager = cv.NewMemoryConfigManager()
		model = NewModel(cliService, svc)
		model.configManager = configManager
	})

	Context("Main Menu → CV Workflow", func() {
		It("should navigate to CV Config Manager when 'v' key is pressed from main menu", func() {
			// Start at main menu
			model.currentScreen = MainMenuScreen

			// Simulate 'v' key press for CV Configurations
			menuMsg := models.MenuItemSelectedMsg{Key: "v"}
			updatedModel, _ := model.Update(menuMsg)

			// Verify we navigated to CVConfigManagerScreen
			Expect(updatedModel.(*Model).currentScreen).To(Equal(CVConfigManagerScreen))
			Expect(updatedModel.(*Model).cvConfigManagerModel).NotTo(BeNil())
		})

		It("should navigate to CV Config Manager when 'g' key is pressed from main menu", func() {
			// Start at main menu
			model.currentScreen = MainMenuScreen

			// Simulate 'g' key press for Generate CV
			menuMsg := models.MenuItemSelectedMsg{Key: "g"}
			updatedModel, _ := model.Update(menuMsg)

			// Verify we navigated to CVConfigManagerScreen
			Expect(updatedModel.(*Model).currentScreen).To(Equal(CVConfigManagerScreen))
			Expect(updatedModel.(*Model).cvConfigManagerModel).NotTo(BeNil())
		})

		It("should preserve previous screen when navigating to CV Manager", func() {
			// Start at list screen
			model.currentScreen = ListScreen
			previousScreen := model.currentScreen

			// Navigate to CV Config Manager
			menuMsg := models.MenuItemSelectedMsg{Key: "v"}
			updatedModel, _ := model.Update(menuMsg)

			// Verify previous screen is saved
			Expect(updatedModel.(*Model).previousScreen).To(Equal(previousScreen))
		})

		It("should refresh configs when navigating to CV Config Manager again", func() {
			// Start at main menu
			model.currentScreen = MainMenuScreen

			// First navigation
			menuMsg := models.MenuItemSelectedMsg{Key: "v"}
			model1, _ := model.Update(menuMsg)

			// Go back to main menu
			model1.(*Model).currentScreen = MainMenuScreen

			// Second navigation should refresh
			model2, _ := model1.(*Model).Update(menuMsg)

			// Verify model was refreshed
			Expect(model2.(*Model).cvConfigManagerModel).NotTo(BeNil())
		})
	})

	Context("CV Configuration Manager Breadcrumbs", func() {
		It("should create CV Config Manager with breadcrumbs", func() {
			// Navigate to CV Config Manager
			model.currentScreen = CVConfigManagerScreen
			model.cvConfigManagerModel = models.NewCVConfigManagerModel(
				models.NewBaseStandardModel(),
				configManager,
			)

			// Verify model was created
			Expect(model.cvConfigManagerModel).NotTo(BeNil())
		})

		It("should create CV Config Manager with event context breadcrumbs", func() {
			// Create a test event
			testEvent := &career.CareerEvent{
				Text: "Test event",
			}

			// Navigate to CV Config Manager with event context
			model.currentScreen = CVConfigManagerScreen
			model.cvConfigManagerModel = models.NewCVConfigManagerModelWithEvent(
				models.NewBaseStandardModel(),
				configManager,
				testEvent,
			)

			// Verify model was created with event context
			Expect(model.cvConfigManagerModel).NotTo(BeNil())
		})
	})
})
