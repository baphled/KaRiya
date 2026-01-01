package app

import (
	"context"
	"time"

	"github.com/baphled/kariya/internal/cli/models"
	"github.com/baphled/kariya/internal/cli/service"
	careerrepo "github.com/baphled/kariya/internal/repository/career"
	careerservice "github.com/baphled/kariya/internal/service/career"
	tea "github.com/charmbracelet/bubbletea"
	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
)

var _ = Describe("End-to-End Integration Tests", func() {
	var (
		repo       careerrepo.Repository
		svc        *careerservice.Service
		cliService *service.CLIEventService
		model      *Model
		ctx        context.Context
	)

	BeforeEach(func() {
		ctx = context.Background()
		repo = careerrepo.NewMemoryRepository()
		svc = careerservice.NewService(repo)
		cliService = service.NewCLIEventService(svc)
		model = NewModel(cliService, svc)
	})

	Describe("Complete Event Capture Workflow", func() {
		Context("when user completes full event capture to persistence", func() {
			It("should capture event through form and persist to repository", func() {
				model.currentScreen = CaptureScreen
				eventText := "Led cross-functional team to deliver critical migration project"
				eventDate := time.Now().Add(-24 * time.Hour)

				err := cliService.CaptureEvent(
					ctx,
					eventText,
					eventDate,
					careerservice.ManualEntry,
					service.WithCompany("TechCorp Inc."),
					service.WithProject("Platform Migration"),
					service.WithTags([]string{"leadership", "technical"}),
				)
				Expect(err).ToNot(HaveOccurred())

				events, err := repo.List(ctx, careerrepo.ListFilters{
					Limit: 10,
				})
				Expect(err).ToNot(HaveOccurred())
				Expect(events).To(HaveLen(1))
				Expect(events[0].Text).To(Equal(eventText))
				Expect(events[0].Company).To(Equal("TechCorp Inc."))
				Expect(events[0].Project).To(Equal("Platform Migration"))
				Expect(events[0].Tags).To(ContainElements("leadership", "technical"))
				Expect(events[0].Date.Format("2006-01-02")).To(Equal(eventDate.Format("2006-01-02")))
			})
		})
	})

	Describe("Facts Workflow Navigation", func() {
		Context("when user presses 't' key from home screen", func() {
			It("should navigate to facts list screen", func() {
				Expect(model.currentScreen).To(Equal(HomeScreen))

				msg := tea.KeyMsg{Type: tea.KeyRunes, Runes: []rune{'t'}}
				newModel, _ := model.Update(msg)
				model = newModel.(*Model)

				Expect(model.currentScreen).To(Equal(FactListScreen))
			})

			It("should initialize fact list model when navigating to facts screen", func() {
				Expect(model.currentScreen).To(Equal(HomeScreen))

				msg := tea.KeyMsg{Type: tea.KeyRunes, Runes: []rune{'t'}}
				newModel, _ := model.Update(msg)
				model = newModel.(*Model)

				Expect(model.factListModel).ToNot(BeNil())
			})
		})

		Context("when user presses back from facts screen", func() {
			It("should navigate back to home screen", func() {
				msg := tea.KeyMsg{Type: tea.KeyRunes, Runes: []rune{'t'}}
				newModel, _ := model.Update(msg)
				model = newModel.(*Model)

				Expect(model.currentScreen).To(Equal(FactListScreen))

				backMsg := models.BackMsg{}
				newModel, _ = model.Update(backMsg)
				model = newModel.(*Model)

				Expect(model.currentScreen).To(Equal(HomeScreen))
			})
		})

		Context("when accessing facts through menu", func() {
			It("should navigate to facts screen when menu item is selected", func() {
				Expect(model.currentScreen).To(Equal(HomeScreen))

				menuMsg := models.MenuItemSelectedMsg{Key: "t"}
				newModel, _ := model.Update(menuMsg)
				model = newModel.(*Model)

				Expect(model.currentScreen).To(Equal(FactListScreen))
			})

			It("should refresh facts list when navigating to facts screen", func() {
				err := cliService.CaptureEvent(
					ctx,
					"Test event",
					time.Now(),
					careerservice.ManualEntry,
				)
				Expect(err).ToNot(HaveOccurred())

				msg := tea.KeyMsg{Type: tea.KeyRunes, Runes: []rune{'t'}}
				newModel, _ := model.Update(msg)
				model = newModel.(*Model)

				Expect(model.factListModel).ToNot(BeNil())

				msg = tea.KeyMsg{Type: tea.KeyRunes, Runes: []rune{'h'}}
				newModel, _ = model.Update(msg)
				model = newModel.(*Model)

				msg = tea.KeyMsg{Type: tea.KeyRunes, Runes: []rune{'t'}}
				newModel, _ = model.Update(msg)
				model = newModel.(*Model)

				Expect(model.currentScreen).To(Equal(FactListScreen))
			})
		})

		Context("when rendering facts screen", func() {
			It("should render facts list view without errors", func() {
				msg := tea.KeyMsg{Type: tea.KeyRunes, Runes: []rune{'t'}}
				newModel, _ := model.Update(msg)
				model = newModel.(*Model)

				Expect(model.currentScreen).To(Equal(FactListScreen))

				view := model.View()
				Expect(view).ToNot(BeEmpty())
			})
		})
	})
})
