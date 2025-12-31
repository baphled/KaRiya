package models_test

import (
	"context"
	"errors"
	"time"

	"github.com/baphled/kariya/internal/cli/models"
	"github.com/baphled/kariya/internal/cli/service"
	"github.com/baphled/kariya/internal/domain/career"
	careerrepo "github.com/baphled/kariya/internal/repository/career"
	careerservice "github.com/baphled/kariya/internal/service/career"
	"github.com/charmbracelet/bubbles/key"
	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
)

var _ = Describe("E2E Model Integration Tests", func() {
	var (
		repo        careerrepo.Repository
		svc         *careerservice.Service
		cliService  *service.CLIEventService
		ctx         context.Context
		formModel   *models.FormModel
		listModel   *models.ListModel
		helpModel   *models.HelpModel
		searchModel *models.SearchModel
		actionMenu  *models.ActionMenuModel
	)

	BeforeEach(func() {
		ctx = context.Background()
		repo = careerrepo.NewMemoryRepository()
		svc = careerservice.NewService(repo)
		cliService = service.NewCLIEventService(svc)

		formModel = models.NewFormModel(cliService)
		listModel = models.NewListModel(svc, ctx)
		helpModel = models.NewHelpModel()
		searchModel = models.NewSearchModel()
		testEvent := &career.CareerEvent{
			ID:        "test",
			Text:      "Test",
			CreatedAt: time.Now(),
			UpdatedAt: time.Now(),
		}
		actionMenu = models.NewActionMenuModel(testEvent)
	})

	Describe("Phase 1: Critical Models E2E", func() {
		Describe("FormModel E2E Tests", func() {
			It("should initialize with StandardModel integration", func() {
				Expect(formModel).NotTo(BeNil())
				Expect(formModel.GetContext()).NotTo(BeNil())
				Expect(formModel.GetBreadcrumbs()).To(HaveLen(0))
			})

			It("should register and handle shortcuts", func() {
				binding := key.NewBinding(key.WithKeys("ctrl+s"))
				shortcuts := map[string]key.Binding{
					"submit": binding,
				}
				formModel.RegisterShortcuts(shortcuts)

				registered := formModel.GetShortcuts()
				Expect(registered).To(HaveKey("submit"))
			})

			It("should track errors with context", func() {
				testErr := errors.New("Invalid date format")
				formModel.SetError(testErr)

				retrieved := formModel.GetLastError()
				Expect(retrieved).NotTo(BeNil())
				Expect(retrieved.Error()).To(ContainSubstring("Invalid date format"))
			})
		})

		Describe("ListModel E2E Tests", func() {
			It("should initialize with StandardModel integration", func() {
				Expect(listModel).NotTo(BeNil())
				Expect(listModel.GetContext()).NotTo(BeNil())
			})

			It("should manage breadcrumbs during navigation", func() {
				listModel.AddBreadcrumb(models.BreadcrumbItem{Label: "Home", ID: "home"})
				listModel.AddBreadcrumb(models.BreadcrumbItem{Label: "Events", ID: "events"})

				breadcrumbs := listModel.GetBreadcrumbs()
				Expect(breadcrumbs).To(HaveLen(2))
				Expect(listModel.GetBreadcrumbPath()).To(Equal("/Home/Events/"))
			})

			It("should reset state properly", func() {
				listModel.AddBreadcrumb(models.BreadcrumbItem{Label: "Test", ID: "test"})
				listModel.PushNavigationHistory("list", "state")
				listModel.SetError(errors.New("test error"))

				listModel.Reset()

				Expect(listModel.GetBreadcrumbs()).To(HaveLen(0))
				Expect(listModel.GetNavigationHistory()).To(HaveLen(0))
				Expect(listModel.GetLastError()).To(BeNil())
			})
		})

		Describe("HelpModel E2E Tests", func() {
			It("should initialize with StandardModel integration", func() {
				Expect(helpModel).NotTo(BeNil())
				Expect(helpModel.GetContext()).NotTo(BeNil())
			})

			It("should handle navigation within help", func() {
				helpModel.PushNavigationHistory("help", "main")
				helpModel.PushNavigationHistory("help", "shortcuts")

				history := helpModel.GetNavigationHistory()
				Expect(history).To(HaveLen(2))
			})
		})

		Describe("SearchModel E2E Tests", func() {
			It("should initialize with StandardModel integration", func() {
				Expect(searchModel).NotTo(BeNil())
				Expect(searchModel.GetContext()).NotTo(BeNil())
			})

			It("should handle search errors gracefully", func() {
				testErr := errors.New("No results found")
				searchModel.SetError(testErr)

				retrieved := searchModel.GetLastError()
				Expect(retrieved).NotTo(BeNil())
			})
		})

		Describe("ActionMenuModel E2E Tests", func() {
			It("should initialize with StandardModel integration", func() {
				Expect(actionMenu).NotTo(BeNil())
				Expect(actionMenu.GetContext()).NotTo(BeNil())
			})

			It("should track menu navigation in history", func() {
				actionMenu.PushNavigationHistory("menu", "main")
				actionMenu.PushNavigationHistory("menu", "event-actions")

				history := actionMenu.GetNavigationHistory()
				Expect(history).To(HaveLen(2))
			})
		})
	})

	Describe("Error Handling E2E", func() {
		It("should track and retrieve errors with context", func() {
			err := errors.New("Invalid date")
			formModel.SetError(err)

			retrieved := formModel.GetLastError()
			Expect(retrieved).NotTo(BeNil())
			Expect(retrieved.Error()).To(ContainSubstring("Invalid date"))
		})

		It("should clear errors properly", func() {
			formModel.SetError(errors.New("test"))
			Expect(formModel.GetLastError()).NotTo(BeNil())

			formModel.ClearError()
			Expect(formModel.GetLastError()).To(BeNil())
		})

		It("should handle nil errors gracefully", func() {
			formModel.SetError(nil)
			Expect(formModel.GetLastError()).To(BeNil())
		})
	})
})
