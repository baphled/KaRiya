package app

import (
	"github.com/baphled/kariya/internal/cli/service"
	"github.com/baphled/kariya/internal/repository/career"
	careerservice "github.com/baphled/kariya/internal/service/career"
	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
)

var _ = Describe("Breadcrumb Display", func() {
	var (
		model      *Model
		repo       *career.MemoryRepository
		svc        *careerservice.Service
		cliService *service.CLIEventService
	)

	BeforeEach(func() {
		repo = career.NewMemoryRepository()
		svc = careerservice.NewService(repo)
		cliService = service.NewCLIEventService(svc)
		model = NewModel(cliService, svc)
	})

	Describe("Form Screen", func() {
		It("should display Home > Capture Event breadcrumbs in form header", func() {
			model.currentScreen = CaptureScreen
			model.updateBreadcrumbs()
			view := model.View()
			Expect(view).To(ContainSubstring("Home"))
			Expect(view).To(ContainSubstring("Capture Event"))
		})
	})

	Describe("List Screen", func() {
		It("should display Home > Events breadcrumbs in list header", func() {
			model.currentScreen = ListScreen
			model.updateBreadcrumbs()
			view := model.View()
			Expect(view).To(ContainSubstring("Home"))
			Expect(view).To(ContainSubstring("Events"))
		})
	})

	Describe("Details Screen", func() {
		It("should display full breadcrumb path in details view", func() {
			Skip("Skipping until we have details model test helper")
			model.previousScreen = ListScreen
			model.currentScreen = ViewScreen
			model.updateBreadcrumbs()
			view := model.View()
			Expect(view).To(ContainSubstring("Home"))
			Expect(view).To(ContainSubstring("Events"))
			Expect(view).To(ContainSubstring("Details"))
		})
	})

	Describe("Metadata Review Screen", func() {
		It("should display Home > Metadata Review breadcrumbs", func() {
			model.currentScreen = MetadataReviewScreen
			model.updateBreadcrumbs()
			view := model.View()
			Expect(view).To(ContainSubstring("Home"))
			Expect(view).To(ContainSubstring("Metadata Review"))
		})
	})

	Describe("Success Screen", func() {
		It("should display breadcrumbs after event capture", func() {
			model.previousScreen = CaptureScreen
			model.currentScreen = SuccessScreen
			model.updateBreadcrumbs()
			view := model.View()
			Expect(view).To(ContainSubstring("Home"))
			Expect(view).To(ContainSubstring("Capture Event"))
			Expect(view).To(ContainSubstring("Success"))
		})
	})
})

