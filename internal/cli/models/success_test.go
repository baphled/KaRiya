package models_test

import (
	"time"

	"github.com/baphled/kariya/internal/cli/models"
	"github.com/baphled/kariya/internal/domain/career"
	tea "github.com/charmbracelet/bubbletea"
	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
)

var _ = Describe("SuccessModel", func() {
	var (
		event *career.CareerEvent
	)

	BeforeEach(func() {
		// Create a sample event for testing
		event = &career.CareerEvent{
			ID:         "test-event-123",
			Text:       "Implemented a successful feature",
			Date:       time.Date(2024, 12, 20, 0, 0, 0, 0, time.UTC),
			Company:    "TechCorp",
			Project:    "Platform Migration",
			Tags:       []string{"technical", "leadership"},
			Categories: []string{"Technical", "Leadership"},
			CreatedAt:  time.Now(),
			UpdatedAt:  time.Now(),
		}
	})

	Describe("NewSuccessModel", func() {
		It("should create a new success model with event", func() {
			model := models.NewSuccessModel(event)

			Expect(model).NotTo(BeNil())
			Expect(model.Event()).To(Equal(event))
		})
	})

	Describe("View", func() {
		It("should render event details", func() {
			model := models.NewSuccessModel(event)
			view := model.View()

			// Verify key event details are displayed
			Expect(view).To(ContainSubstring("Implemented a successful feature"))
			Expect(view).To(ContainSubstring("TechCorp"))
			Expect(view).To(ContainSubstring("Platform Migration"))
			Expect(view).To(ContainSubstring("technical"))
			Expect(view).To(ContainSubstring("leadership"))
		})

		It("should include success header", func() {
			model := models.NewSuccessModel(event)
			view := model.View()

			Expect(view).To(ContainSubstring("Success"))
		})

		It("should include action options", func() {
			model := models.NewSuccessModel(event)
			view := model.View()

			Expect(view).To(ContainSubstring("Capture Another"))
			Expect(view).To(ContainSubstring("View Recent"))
			Expect(view).To(ContainSubstring("Exit"))
		})
	})

	Describe("Update", func() {
		It("should handle keyboard navigation", func() {
			model := models.NewSuccessModel(event)

			// Simulate left arrow key
			keyMsg := tea.KeyMsg{Type: tea.KeyLeft}
			newModel, cmd := model.Update(keyMsg)

			Expect(newModel).NotTo(BeNil())
			Expect(cmd).To(BeNil()) // Navigation doesn't trigger commands
		})

		It("should handle action selection with enter key", func() {
			model := models.NewSuccessModel(event)

			// Simulate enter key
			keyMsg := tea.KeyMsg{Type: tea.KeyEnter}
			newModel, cmd := model.Update(keyMsg)

			Expect(newModel).NotTo(BeNil())
			Expect(cmd).NotTo(BeNil()) // Enter triggers a command
		})
	})
})
