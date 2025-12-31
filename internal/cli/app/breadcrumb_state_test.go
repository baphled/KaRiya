package app

import (
	"context"

	"github.com/baphled/kariya/internal/cli/service"
	"github.com/baphled/kariya/internal/repository/career"
	careerservice "github.com/baphled/kariya/internal/service/career"
	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
)

var _ = Describe("Breadcrumb State Management", func() {
	var (
		model      *Model
		repo       *career.MemoryRepository
		svc        *careerservice.Service
		cliService *service.CLIEventService
		ctx        context.Context
	)

	BeforeEach(func() {
		ctx = context.Background()
		repo = career.NewMemoryRepository()
		svc = careerservice.NewService(repo)
		cliService = service.NewCLIEventService(svc)
		model = NewModel(cliService, svc)
	})

	Describe("Initialization", func() {
		It("should initialize with Home breadcrumb", func() {
			Expect(model.GetBreadcrumbs()).To(Equal([]string{"Home"}))
		})
	})

	Describe("Screen Navigation", func() {
		Context("when navigating to CaptureScreen", func() {
			It("should update breadcrumbs to Home > Capture Event", func() {
				model.currentScreen = CaptureScreen
				model.updateBreadcrumbs()
				Expect(model.GetBreadcrumbs()).To(Equal([]string{"Home", "Capture Event"}))
			})
		})

		Context("when navigating to ListScreen", func() {
			It("should update breadcrumbs to Home > Events", func() {
				model.currentScreen = ListScreen
				model.updateBreadcrumbs()
				Expect(model.GetBreadcrumbs()).To(Equal([]string{"Home", "Events"}))
			})
		})

		Context("when navigating through nested screens", func() {
			It("should build breadcrumb trail Home > Events > Details", func() {
				model.currentScreen = ListScreen
				model.updateBreadcrumbs()
				model.previousScreen = ListScreen
				model.currentScreen = ViewScreen
				model.updateBreadcrumbs()
				Expect(model.GetBreadcrumbs()).To(Equal([]string{"Home", "Events", "Details"}))
			})
		})

		Context("when returning to HomeScreen", func() {
			It("should reset breadcrumbs to Home only", func() {
				model.currentScreen = ListScreen
				model.updateBreadcrumbs()
				model.currentScreen = HomeScreen
				model.updateBreadcrumbs()
				Expect(model.GetBreadcrumbs()).To(Equal([]string{"Home"}))
			})
		})
	})

	Describe("Breadcrumb Trail Building", func() {
		Context("with MetadataReviewScreen", func() {
			It("should show Home > Metadata Review", func() {
				model.currentScreen = MetadataReviewScreen
				model.updateBreadcrumbs()
				Expect(model.GetBreadcrumbs()).To(Equal([]string{"Home", "Metadata Review"}))
			})
		})

		Context("with SuccessScreen after CaptureScreen", func() {
			It("should show Home > Capture Event > Success", func() {
				model.currentScreen = CaptureScreen
				model.updateBreadcrumbs()
				model.previousScreen = CaptureScreen
				model.currentScreen = SuccessScreen
				model.updateBreadcrumbs()
				Expect(model.GetBreadcrumbs()).To(Equal([]string{"Home", "Capture Event", "Success"}))
			})
		})
	})
})

