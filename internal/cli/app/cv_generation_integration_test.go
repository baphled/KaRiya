package app

import (
	"time"

	"github.com/baphled/kariya/internal/cli/models"
	"github.com/baphled/kariya/internal/cli/service"
	"github.com/baphled/kariya/internal/domain/career"
	careerrepo "github.com/baphled/kariya/internal/repository/career"
	careerservice "github.com/baphled/kariya/internal/service/career"
	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
)

var _ = Describe("CV Generation Integration", func() {
	var (
		repo       *careerrepo.MemoryRepository
		svc        *careerservice.Service
		cliService *service.CLIEventService
		model      *Model
	)

	BeforeEach(func() {
		repo = careerrepo.NewMemoryRepository()
		svc = careerservice.NewService(repo)
		cliService = service.NewCLIEventService(svc)
		model = NewModel(cliService, svc)
	})

	Context("Event Timeline to CV Generation Workflow", func() {
		It("should navigate to CV config manager when Generate CV action is selected", func() {
			// Create a test event
			testEvent := &career.CareerEvent{
				ID:      "test-event-1",
				Text:    "Led team to deliver critical feature",
				Date:    time.Now().Add(-10 * 24 * time.Hour),
				Company: "TechCorp Inc",
				Project: "Platform Modernization",
				Tags:    []string{"leadership", "technical"},
			}

			// Set up the app to be on the action menu screen
			model.currentScreen = ActionMenuScreen
			model.previousScreen = ListScreen

			// Simulate selecting the Generate CV action from event action menu
			actionMsg := models.EventActionSelectedMsg{
				Event:  testEvent,
				Action: models.EventActionGenerateCV,
			}

			// Update the model with the action message
			updatedModel, _ := model.Update(actionMsg)
			model = updatedModel.(*Model)

			// Verify we navigated to CV config manager screen
			Expect(model.currentScreen).To(Equal(CVConfigManagerScreen))
			Expect(model.cvConfigManagerModel).NotTo(BeNil())
			Expect(model.previousScreen).To(Equal(ActionMenuScreen))
		})

		It("should navigate to CV preview after config selection", func() {
			// Create a test config
			testConfig := &career.CVConfig{
				Name:           "Principal Engineer CV",
				TargetRole:     "Principal Engineer",
				TargetAudience: []string{"hiring_manager", "recruiter"},
				CreatedAt:      time.Now(),
				UpdatedAt:      time.Now(),
			}

			// Simulate the GenerateCVFromConfigMsg
			configMsg := models.GenerateCVFromConfigMsg{
				Config: testConfig,
			}

			updatedModel, _ := model.Update(configMsg)
			model = updatedModel.(*Model)

			// Verify we're on the CV generator screen
			Expect(model.currentScreen).To(Equal(CVGeneratorScreen))
			Expect(model.cvGeneratorModel).NotTo(BeNil())
		})

		It("should navigate from CV preview back to event timeline", func() {
			// Create a CV view
			testCVView := &career.CVView{
				ID:              "cv-1",
				Name:            "Test CV",
				TargetRole:      "Senior Engineer",
				TargetAudience:  []string{"hiring_manager"},
				GeneratedAt:     time.Now(),
				SourceEventCount: 1,
				SourceFactCount:  0,
			}

			// Simulate navigating to CV preview from event timeline
			model.previousScreen = ListScreen
			model.currentScreen = CVPreviewScreen
			model.cvPreviewModel = models.NewCVPreviewModelWithSource(
				models.NewBaseStandardModel(),
				testCVView,
				[]*career.CVSection{},
				nil,
				"event_timeline",
			)

			// Simulate pressing Esc to go back
			backMsg := models.BackToEventTimelineMsg{}
			updatedModel, _ := model.Update(backMsg)
			model = updatedModel.(*Model)

			// Verify we're back on the event timeline
			Expect(model.currentScreen).To(Equal(ListScreen))
			Expect(model.listModel).NotTo(BeNil())
		})

		It("should navigate from CV preview back to config manager when appropriate", func() {
			// Create a CV view
			testCVView := &career.CVView{
				ID:              "cv-2",
				Name:            "Test CV 2",
				TargetRole:      "Staff Engineer",
				TargetAudience:  []string{"hiring_manager"},
				GeneratedAt:     time.Now(),
				SourceEventCount: 1,
				SourceFactCount:  0,
			}

			// Simulate navigating to CV preview from config manager
			model.cvConfigManagerModel = models.NewCVConfigManagerModel(models.NewBaseStandardModel(), model.configManager)
			model.previousScreen = CVConfigManagerScreen
			model.currentScreen = CVPreviewScreen
			model.cvPreviewModel = models.NewCVPreviewModelWithSource(
				models.NewBaseStandardModel(),
				testCVView,
				[]*career.CVSection{},
				nil,
				"cv_config_manager",
			)

			// Simulate pressing Esc to go back
			backMsg := models.BackMsg{}
			updatedModel, _ := model.Update(backMsg)
			model = updatedModel.(*Model)

			// Verify we're back on the config manager
			Expect(model.currentScreen).To(Equal(CVConfigManagerScreen))
			Expect(model.cvConfigManagerModel).NotTo(BeNil())
		})

		It("should support full workflow: event -> action menu -> config -> generation -> preview -> back to timeline", func() {
			// Create a test event
			testEvent := &career.CareerEvent{
				ID:      "workflow-event",
				Text:    "Led architectural redesign of core platform",
				Date:    time.Now().Add(-15 * 24 * time.Hour),
				Company: "TechCorp Inc",
				Project: "Platform Redesign",
				Tags:    []string{"leadership", "technical", "product"},
			}

			// Step 1: Simulate selecting Generate CV action
			actionMsg := models.EventActionSelectedMsg{
				Event:  testEvent,
				Action: models.EventActionGenerateCV,
			}

			updatedModel, _ := model.Update(actionMsg)
			model = updatedModel.(*Model)

			// Verify we're on config manager
			Expect(model.currentScreen).To(Equal(CVConfigManagerScreen))

			// Step 2: Simulate selecting a config
			testConfig := &career.CVConfig{
				Name:           "Workflow Test CV",
				TargetRole:     "Principal Engineer",
				TargetAudience: []string{"hiring_manager"},
				CreatedAt:      time.Now(),
				UpdatedAt:      time.Now(),
			}

			configMsg := models.GenerateCVFromConfigMsg{
				Config: testConfig,
			}

			updatedModel, _ = model.Update(configMsg)
			model = updatedModel.(*Model)

			// Verify we're on CV generator
			Expect(model.currentScreen).To(Equal(CVGeneratorScreen))

			// Step 3: Simulate CV generation completion
			testCVView := &career.CVView{
				ID:              "workflow-cv",
				Name:            "Workflow Test CV",
				TargetRole:      "Principal Engineer",
				TargetAudience:  []string{"hiring_manager"},
				GeneratedAt:     time.Now(),
				SourceEventCount: 1,
				SourceFactCount:  0,
			}

			previewMsg := models.NavigateToCVPreviewMsg{
				CVView:       testCVView,
				SourceScreen: "event_timeline",
			}

			updatedModel, _ = model.Update(previewMsg)
			model = updatedModel.(*Model)

			// Verify we're on CV preview
			Expect(model.currentScreen).To(Equal(CVPreviewScreen))
			Expect(model.cvPreviewModel).NotTo(BeNil())

			// Step 4: Navigate back to event timeline
			backMsg := models.BackToEventTimelineMsg{}
			updatedModel, _ = model.Update(backMsg)
			model = updatedModel.(*Model)

			// Verify we're back on the event timeline
			Expect(model.currentScreen).To(Equal(ListScreen))
		})
	})
})
