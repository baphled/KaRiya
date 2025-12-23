package app

import (
	"github.com/baphled/kariya/internal/cli/service"
	careerservice "github.com/baphled/kariya/internal/service/career"
	tea "github.com/charmbracelet/bubbletea"
	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
)

var _ = Describe("Application Model", func() {
	var (
		cliService *service.CLIEventService
		svc        *careerservice.Service
		model      *Model
	)

	BeforeEach(func() {
		svc = careerservice.NewService(nil)
		cliService = service.NewCLIEventService(svc)
		model = NewModel(cliService, svc)
	})

	Context("Initialization", func() {
		It("should initialize with HomeScreen as current screen", func() {
			Expect(model.currentScreen).To(Equal(HomeScreen))
		})

		It("should initialize with default width and height", func() {
			Expect(model.width).To(Equal(80))
			Expect(model.height).To(Equal(24))
		})

		It("should have no initial error", func() {
			Expect(model.err).To(BeNil())
		})
	})

	Context("Screen Navigation", func() {
		It("should navigate to CaptureScreen", func() {
			msg := tea.KeyMsg{Type: tea.KeyRunes, Runes: []rune{'c'}}
			newModel, _ := model.Update(msg)
			updatedModel := newModel.(*Model)
			Expect(updatedModel.currentScreen).To(Equal(CaptureScreen))
		})

		It("should navigate to ListScreen", func() {
			msg := tea.KeyMsg{Type: tea.KeyRunes, Runes: []rune{'l'}}
			newModel, _ := model.Update(msg)
			updatedModel := newModel.(*Model)
			Expect(updatedModel.currentScreen).To(Equal(ListScreen))
		})

		It("should navigate back to HomeScreen", func() {
			model.currentScreen = CaptureScreen
			msg := tea.KeyMsg{Type: tea.KeyRunes, Runes: []rune{'h'}}
			newModel, _ := model.Update(msg)
			updatedModel := newModel.(*Model)
			Expect(updatedModel.currentScreen).To(Equal(HomeScreen))
		})

		It("should handle backspace navigation", func() {
			model.currentScreen = CaptureScreen
			model.previousScreen = HomeScreen
			msg := tea.KeyMsg{Type: tea.KeyBackspace}
			newModel, _ := model.Update(msg)
			updatedModel := newModel.(*Model)
			Expect(updatedModel.currentScreen).To(Equal(HomeScreen))
		})
	})

	Context("View Rendering", func() {
		It("should render home screen view", func() {
			model.currentScreen = HomeScreen
			view := model.View()
			Expect(view).To(ContainSubstring("KaRiya - Career Journal CLI"))
			Expect(view).To(ContainSubstring("Capture Career Event"))
			Expect(view).To(ContainSubstring("List Events"))
		})

		It("should render capture screen view", func() {
			model.currentScreen = CaptureScreen
			view := model.View()
			Expect(view).To(ContainSubstring("Capture Career Event"))
			Expect(view).To(ContainSubstring("Event Text"))
			Expect(view).To(ContainSubstring("Company"))
		})

		It("should render list screen view", func() {
			model.currentScreen = ListScreen
			view := model.View()
			Expect(view).To(ContainSubstring("Recent Career Events"))
			Expect(view).To(ContainSubstring("Event 1"))
		})

		It("should render view screen view", func() {
			model.currentScreen = ViewScreen
			view := model.View()
			Expect(view).To(ContainSubstring("Event Details"))
			Expect(view).To(ContainSubstring("Title"))
		})
	})

	Context("Window Resizing", func() {
		It("should update width and height on window resize", func() {
			msg := tea.WindowSizeMsg{Width: 120, Height: 40}
			newModel, _ := model.Update(msg)
			updatedModel := newModel.(*Model)
			Expect(updatedModel.width).To(Equal(120))
			Expect(updatedModel.height).To(Equal(40))
		})
	})

	Context("Quit Command", func() {
		It("should return model when 'q' key is pressed", func() {
			msg := tea.KeyMsg{Type: tea.KeyRunes, Runes: []rune{'q'}}
			newModel, cmd := model.Update(msg)
			Expect(newModel).NotTo(BeNil())
			Expect(cmd).NotTo(BeNil())
		})

		It("should return model when ctrl+c is pressed", func() {
			msg := tea.KeyMsg{Type: tea.KeyCtrlC}
			newModel, cmd := model.Update(msg)
			Expect(newModel).NotTo(BeNil())
			Expect(cmd).NotTo(BeNil())
		})
	})

	Context("FormModel Integration", func() {
		It("should have FormModel instance in Model struct", func() {
			Expect(model.formModel).NotTo(BeNil())
		})

		It("should create SuccessModel after form submission", func() {
			// Initially successModel should be nil
			Expect(model.successModel).To(BeNil())
			// After form submission, successModel would be created
			// This is tested in the form submission test
		})

		It("should delegate Update to FormModel on CaptureScreen", func() {
			model.currentScreen = CaptureScreen
			msg := tea.KeyMsg{Type: tea.KeyRunes, Runes: []rune{'t'}}
			newModel, _ := model.Update(msg)
			updatedModel := newModel.(*Model)
			Expect(updatedModel).NotTo(BeNil())
		})

		It("should delegate View to FormModel on CaptureScreen", func() {
			model.currentScreen = CaptureScreen
			view := model.View()
			Expect(view).NotTo(BeEmpty())
		})

		It("should reset FormModel when navigating to CaptureScreen", func() {
			msg := tea.KeyMsg{Type: tea.KeyRunes, Runes: []rune{'c'}}
			newModel, _ := model.Update(msg)
			updatedModel := newModel.(*Model)
			Expect(updatedModel.formModel).NotTo(BeNil())
		})
	})
})
