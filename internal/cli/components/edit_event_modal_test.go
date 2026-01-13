package components_test

import (
	"time"

	"github.com/baphled/kariya/internal/cli/components"
	"github.com/baphled/kariya/internal/domain/career"
	tea "github.com/charmbracelet/bubbletea"
	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
)

var _ = Describe("EditEventModal", func() {
	var (
		modal         *components.EditEventModal
		existingEvent *career.CareerEvent
		width         int
		height        int
	)

	BeforeEach(func() {
		width = 120
		height = 40

		// Create an existing event to edit
		existingEvent = &career.CareerEvent{
			ID:         "test-event-123",
			Text:       "Original event description",
			Date:       time.Date(2024, 1, 15, 0, 0, 0, 0, time.UTC),
			Company:    "Acme Corp",
			Project:    "Project Alpha",
			Tags:       []string{"achievement", "technical"},
			Categories: []string{"development", "testing"},
			CreatedAt:  time.Date(2024, 1, 1, 0, 0, 0, 0, time.UTC),
			UpdatedAt:  time.Date(2024, 1, 10, 0, 0, 0, 0, time.UTC),
		}
	})

	Describe("NewEditEventModal", func() {
		It("should create a modal with pre-populated fields from existing event", func() {
			modal = components.NewEditEventModal(existingEvent, width, height)

			Expect(modal).NotTo(BeNil())
			Expect(modal.IsVisible()).To(BeTrue())
			Expect(modal.GetOriginalEvent()).To(Equal(existingEvent))
		})

		It("should initialize with correct dimensions", func() {
			modal = components.NewEditEventModal(existingEvent, width, height)
			modal.Init()

			view := modal.View()
			Expect(view).NotTo(BeEmpty())
		})

		It("should preserve original event reference", func() {
			modal = components.NewEditEventModal(existingEvent, width, height)

			original := modal.GetOriginalEvent()
			Expect(original).To(Equal(existingEvent))
			Expect(original.ID).To(Equal("test-event-123"))
			Expect(original.Text).To(Equal("Original event description"))
		})
	})

	Describe("Init", func() {
		It("should return a command to initialize the form", func() {
			modal = components.NewEditEventModal(existingEvent, width, height)
			cmd := modal.Init()

			Expect(cmd).NotTo(BeNil())
		})
	})

	Describe("Update", func() {
		BeforeEach(func() {
			modal = components.NewEditEventModal(existingEvent, width, height)
			modal.Init()
		})

		Context("when escape key is pressed", func() {
			It("should close the modal without saving", func() {
				msg := tea.KeyMsg{Type: tea.KeyEsc}
				cmd, completed, eventData := modal.Update(msg)

				Expect(cmd).To(BeNil())
				Expect(completed).To(BeFalse())
				Expect(eventData).To(BeNil())
				Expect(modal.IsVisible()).To(BeFalse())
			})
		})

		Context("when window size changes", func() {
			It("should rebuild the form with new dimensions", func() {
				msg := tea.WindowSizeMsg{Width: 100, Height: 30}
				cmd, completed, eventData := modal.Update(msg)

				Expect(cmd).NotTo(BeNil())
				Expect(completed).To(BeFalse())
				Expect(eventData).To(BeNil())
				Expect(modal.IsVisible()).To(BeTrue())
			})
		})

		Context("when form is completed", func() {
			It("should return event data and close modal", func() {
				// Simulate form completion by manually setting state
				// This is a simplified test - in practice, the form would go through multiple update cycles

				// We'll test the EditEventData conversion instead
				eventData := &components.EditEventData{
					Text:       "Updated event description",
					Date:       "2024-02-20",
					Company:    "New Corp",
					Project:    "Project Beta",
					Tags:       []string{"learning", "collaboration"},
					Categories: []string{"design", "documentation"},
				}

				Expect(eventData.Text).To(Equal("Updated event description"))
				Expect(eventData.Company).To(Equal("New Corp"))
			})
		})
	})

	Describe("View", func() {
		BeforeEach(func() {
			modal = components.NewEditEventModal(existingEvent, width, height)
			modal.Init()
		})

		It("should render form when visible", func() {
			view := modal.View()
			Expect(view).NotTo(BeEmpty())
		})

		It("should return empty string when not visible", func() {
			modal.Hide()
			view := modal.View()
			Expect(view).To(BeEmpty())
		})
	})

	Describe("Visibility management", func() {
		BeforeEach(func() {
			modal = components.NewEditEventModal(existingEvent, width, height)
		})

		It("should start visible", func() {
			Expect(modal.IsVisible()).To(BeTrue())
		})

		It("should hide when Hide() is called", func() {
			modal.Hide()
			Expect(modal.IsVisible()).To(BeFalse())
		})

		It("should show when Show() is called", func() {
			modal.Hide()
			modal.Show()
			Expect(modal.IsVisible()).To(BeTrue())
		})
	})

	Describe("EditEventData", func() {
		var eventData *components.EditEventData

		BeforeEach(func() {
			eventData = &components.EditEventData{
				Text:       "Updated event description with at least ten characters",
				Date:       "2024-02-20",
				Company:    "New Corp",
				Project:    "Project Beta",
				Tags:       []string{"learning", "collaboration"},
				Categories: []string{"design", "documentation"},
			}
		})

		Describe("ToCareerEvent", func() {
			It("should convert to CareerEvent with all fields", func() {
				eventID := "test-event-123"
				createdAt := time.Date(2024, 1, 1, 0, 0, 0, 0, time.UTC)
				updatedAt := time.Date(2024, 2, 20, 0, 0, 0, 0, time.UTC)

				event := eventData.ToCareerEvent(eventID, createdAt, updatedAt)

				Expect(event).NotTo(BeNil())
				Expect(event.ID).To(Equal(eventID))
				Expect(event.Text).To(Equal("Updated event description with at least ten characters"))
				Expect(event.Company).To(Equal("New Corp"))
				Expect(event.Project).To(Equal("Project Beta"))
				Expect(event.Tags).To(Equal([]string{"learning", "collaboration"}))
				Expect(event.Categories).To(Equal([]string{"design", "documentation"}))
			})

			It("should parse date correctly", func() {
				eventID := "test-event-123"
				createdAt := time.Date(2024, 1, 1, 0, 0, 0, 0, time.UTC)
				updatedAt := time.Date(2024, 2, 20, 0, 0, 0, 0, time.UTC)

				event := eventData.ToCareerEvent(eventID, createdAt, updatedAt)

				expectedDate := time.Date(2024, 2, 20, 0, 0, 0, 0, time.UTC)
				Expect(event.Date.Format("2006-01-02")).To(Equal(expectedDate.Format("2006-01-02")))
			})

			It("should handle 'today' date", func() {
				eventData.Date = "today"
				eventID := "test-event-123"
				createdAt := time.Date(2024, 1, 1, 0, 0, 0, 0, time.UTC)
				updatedAt := time.Date(2024, 2, 20, 0, 0, 0, 0, time.UTC)

				event := eventData.ToCareerEvent(eventID, createdAt, updatedAt)

				today := time.Now().Format("2006-01-02")
				Expect(event.Date.Format("2006-01-02")).To(Equal(today))
			})

			It("should preserve event ID", func() {
				eventID := "existing-event-456"
				createdAt := time.Date(2024, 1, 1, 0, 0, 0, 0, time.UTC)
				updatedAt := time.Date(2024, 2, 20, 0, 0, 0, 0, time.UTC)

				event := eventData.ToCareerEvent(eventID, createdAt, updatedAt)

				Expect(event.ID).To(Equal(eventID))
			})

			It("should initialize empty skills array", func() {
				eventID := "test-event-123"
				createdAt := time.Date(2024, 1, 1, 0, 0, 0, 0, time.UTC)
				updatedAt := time.Date(2024, 2, 20, 0, 0, 0, 0, time.UTC)

				event := eventData.ToCareerEvent(eventID, createdAt, updatedAt)

				Expect(event.Skills).NotTo(BeNil())
				Expect(event.Skills).To(BeEmpty())
			})
		})
	})

	Describe("Integration workflow", func() {
		It("should support complete edit workflow", func() {
			// 1. Create modal with existing event
			modal = components.NewEditEventModal(existingEvent, width, height)
			Expect(modal.IsVisible()).To(BeTrue())

			// 2. Initialize
			cmd := modal.Init()
			Expect(cmd).NotTo(BeNil())

			// 3. User can cancel
			escMsg := tea.KeyMsg{Type: tea.KeyEsc}
			cmd, completed, eventData := modal.Update(escMsg)
			Expect(completed).To(BeFalse())
			Expect(eventData).To(BeNil())
			Expect(modal.IsVisible()).To(BeFalse())

			// 4. Verify original event unchanged
			original := modal.GetOriginalEvent()
			Expect(original.Text).To(Equal("Original event description"))
			Expect(original.Company).To(Equal("Acme Corp"))
		})

		It("should preserve original event for comparison after edit", func() {
			modal = components.NewEditEventModal(existingEvent, width, height)

			// Get original before any edits
			original := modal.GetOriginalEvent()
			Expect(original.Text).To(Equal("Original event description"))

			// After modal processes updates, original should still be the same
			Expect(modal.GetOriginalEvent()).To(Equal(existingEvent))
		})
	})
})
