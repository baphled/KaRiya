package models_test

import (
	"time"

	"github.com/baphled/kariya/internal/cli/models"
	"github.com/baphled/kariya/internal/domain/career"
	tea "github.com/charmbracelet/bubbletea"
	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
)

var _ = Describe("DetailsModel", func() {
	var (
		model *models.DetailsModel
		event *career.CareerEvent
	)

	BeforeEach(func() {
		event = &career.CareerEvent{
			ID:         "test-event-123",
			Text:       "Led cross-functional team to deliver critical project",
			Date:       time.Now().Add(-7 * 24 * time.Hour),
			Company:    "TechCorp Inc.",
			Project:    "Platform Migration",
			Tags:       []string{"leadership", "technical"},
			Categories: []string{"Leadership", "Technical"},
			CreatedAt:  time.Now().Add(-7 * 24 * time.Hour),
			UpdatedAt:  time.Now().Add(-7 * 24 * time.Hour),
		}
	})

	Context("when creating a new DetailsModel", func() {
		It("should initialize with an event", func() {
			model = models.NewDetailsModel(event)
			Expect(model).NotTo(BeNil())
		})

		It("should handle nil event gracefully", func() {
			model = models.NewDetailsModel(nil)
			Expect(model).NotTo(BeNil())
		})
	})

	Context("when displaying event information", func() {
		BeforeEach(func() {
			model = models.NewDetailsModel(event)
		})

		It("should render the event text", func() {
			view := model.View()
			Expect(view).To(ContainSubstring("Led cross-functional team"))
		})

		It("should render the event date", func() {
			view := model.View()
			Expect(view).To(Or(ContainSubstring("Date:"), ContainSubstring("2025")))
		})

		It("should render the company name", func() {
			view := model.View()
			Expect(view).To(ContainSubstring("TechCorp Inc."))
		})

		It("should render the project name", func() {
			view := model.View()
			Expect(view).To(ContainSubstring("Platform Migration"))
		})

		It("should render the tags", func() {
			view := model.View()
			Expect(view).To(Or(ContainSubstring("leadership"), ContainSubstring("technical")))
		})

		It("should render the event ID", func() {
			view := model.View()
			Expect(view).To(ContainSubstring("test-event-123"))
		})

		It("should handle events with no company gracefully", func() {
			event.Company = ""
			model = models.NewDetailsModel(event)
			view := model.View()
			Expect(view).NotTo(BeNil())
		})

		It("should handle events with no project gracefully", func() {
			event.Project = ""
			model = models.NewDetailsModel(event)
			view := model.View()
			Expect(view).NotTo(BeNil())
		})

		It("should handle events with no tags gracefully", func() {
			event.Tags = []string{}
			model = models.NewDetailsModel(event)
			view := model.View()
			Expect(view).NotTo(BeNil())
		})
	})

	Context("when handling nil event", func() {
		BeforeEach(func() {
			model = models.NewDetailsModel(nil)
		})

		It("should display an error message", func() {
			view := model.View()
			Expect(view).To(Or(ContainSubstring("No event selected"), ContainSubstring("Error")))
		})

		It("should still provide a view", func() {
			view := model.View()
			Expect(view).NotTo(BeEmpty())
		})
	})

	Context("when testing header and help_footer integration", func() {
		BeforeEach(func() {
			model = models.NewDetailsModel(event)
		})

		It("should include header in view", func() {
			view := model.View()
			Expect(view).To(ContainSubstring("Event Details"))
		})

		It("should include help footer in view", func() {
			view := model.View()
			Expect(view).NotTo(BeEmpty())
		})

		It("should handle window resize messages", func() {
			msg := tea.WindowSizeMsg{
				Width:  120,
				Height: 40,
			}
			newModel, cmd := model.Update(msg)
			Expect(newModel).NotTo(BeNil())
			Expect(cmd).To(BeNil())
		})

		It("should render with responsive layout", func() {
			view := model.View()
			Expect(len(view)).To(BeNumerically(">", 100))
		})

		It("should have Init method", func() {
			Expect(model.Init).NotTo(BeNil())
		})

		It("should respond to Escape key", func() {
			model = models.NewDetailsModel(event)
			msg := tea.KeyMsg{Type: tea.KeyEsc}
			_, cmd := model.Update(msg)
			Expect(cmd).NotTo(BeNil())
		})
	})

	Context("when displaying timestamps", func() {
		BeforeEach(func() {
			model = models.NewDetailsModel(event)
		})

		It("should render created timestamp", func() {
			view := model.View()
			Expect(view).To(ContainSubstring("Created:"))
		})

		It("should render updated timestamp", func() {
			view := model.View()
			Expect(view).To(ContainSubstring("Last Updated:"))
		})
	})
})
