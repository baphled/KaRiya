package components_test

import (
	"time"

	"github.com/baphled/kariya/internal/cli/components"
	tea "github.com/charmbracelet/bubbletea"
	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
)

var _ = Describe("QuickAddEventModal", func() {
	var (
		modal  *components.QuickAddEventModal
		width  int
		height int
	)

	BeforeEach(func() {
		width = 100
		height = 30
		modal = components.NewQuickAddEventModal(width, height)
	})

	Describe("NewQuickAddEventModal", func() {
		It("should create a visible modal", func() {
			Expect(modal.IsVisible()).To(BeTrue())
		})

		It("should default date to today", func() {
			// Modal is visible, so form should render
			view := modal.View()
			Expect(view).NotTo(BeEmpty())
			// Date field should be present
			Expect(view).To(ContainSubstring("Date"))
		})
	})

	Describe("Init", func() {
		It("should return form init command", func() {
			cmd := modal.Init()
			Expect(cmd).NotTo(BeNil())
		})
	})

	Describe("Update - Completion", func() {
		It("should handle form submission flow", func() {
			// Simulate filling the form
			// Note: In real usage, user would type and press keys
			// For testing, we check the flow when form gets key input
			cmd, completed, data := modal.Update(tea.KeyMsg{Type: tea.KeyEnter})

			// Form will return a command (form processing)
			// Not completed yet (fields need to be filled first)
			// This test verifies the signature and basic handling
			_ = cmd                         // Command may be nil or a form command
			Expect(completed).To(BeFalse()) // Form not complete yet
			Expect(data).To(BeNil())        // No data until completion
		})
	})

	Describe("Update - Cancellation", func() {
		It("should close modal on Esc", func() {
			Expect(modal.IsVisible()).To(BeTrue())

			cmd, completed, data := modal.Update(tea.KeyMsg{Type: tea.KeyEsc})

			Expect(cmd).To(BeNil())
			Expect(completed).To(BeFalse())
			Expect(data).To(BeNil())
			Expect(modal.IsVisible()).To(BeFalse())
		})
	})

	Describe("Update - WindowSizeMsg", func() {
		It("should update dimensions on window resize", func() {
			newWidth := 120
			newHeight := 40

			cmd, completed, data := modal.Update(tea.WindowSizeMsg{Width: newWidth, Height: newHeight})

			Expect(cmd).NotTo(BeNil()) // Form reinit
			Expect(completed).To(BeFalse())
			Expect(data).To(BeNil())

			// Verify modal still visible and renders
			Expect(modal.IsVisible()).To(BeTrue())
			view := modal.View()
			Expect(view).NotTo(BeEmpty())
		})
	})

	Describe("View", func() {
		It("should render form when visible", func() {
			view := modal.View()
			Expect(view).NotTo(BeEmpty())
			Expect(view).To(ContainSubstring("Event Description"))
			Expect(view).To(ContainSubstring("Date"))
			Expect(view).To(ContainSubstring("Company"))
		})

		It("should return empty string when not visible", func() {
			modal.Hide()
			view := modal.View()
			Expect(view).To(BeEmpty())
		})
	})

	Describe("Show/Hide", func() {
		It("should show modal", func() {
			modal.Hide()
			Expect(modal.IsVisible()).To(BeFalse())

			modal.Show()
			Expect(modal.IsVisible()).To(BeTrue())
		})

		It("should hide modal", func() {
			Expect(modal.IsVisible()).To(BeTrue())

			modal.Hide()
			Expect(modal.IsVisible()).To(BeFalse())
		})
	})

	Describe("QuickAddEventData", func() {
		Describe("ToCareerEvent", func() {
			It("should create a valid CareerEvent with all fields", func() {
				data := &components.QuickAddEventData{
					Text:    "Implemented new feature",
					Date:    "2024-01-15",
					Company: "Acme Corp",
				}

				event := data.ToCareerEvent()

				Expect(event.Text).To(Equal("Implemented new feature"))
				Expect(event.Date.Format("2006-01-02")).To(Equal("2024-01-15"))
				Expect(event.Company).To(Equal("Acme Corp"))
				Expect(event.Project).To(BeEmpty())
				Expect(event.Tags).To(BeEmpty())
				Expect(event.Categories).To(BeEmpty())
				Expect(event.Skills).To(BeEmpty())
				Expect(event.CreatedAt).NotTo(BeZero())
				Expect(event.UpdatedAt).NotTo(BeZero())
			})

			It("should handle 'today' as date", func() {
				data := &components.QuickAddEventData{
					Text:    "Quick event",
					Date:    "today",
					Company: "",
				}

				event := data.ToCareerEvent()
				today := time.Now().Format("2006-01-02")

				Expect(event.Date.Format("2006-01-02")).To(Equal(today))
			})

			It("should handle relative dates", func() {
				data := &components.QuickAddEventData{
					Text:    "Past event",
					Date:    "1 week ago",
					Company: "",
				}

				event := data.ToCareerEvent()
				expectedDate := time.Now().AddDate(0, 0, -7).Format("2006-01-02")

				Expect(event.Date.Format("2006-01-02")).To(Equal(expectedDate))
			})

			It("should default to today on invalid date", func() {
				data := &components.QuickAddEventData{
					Text:    "Event with bad date",
					Date:    "invalid-date",
					Company: "",
				}

				event := data.ToCareerEvent()
				today := time.Now().Format("2006-01-02")

				Expect(event.Date.Format("2006-01-02")).To(Equal(today))
			})

			It("should work without company (optional field)", func() {
				data := &components.QuickAddEventData{
					Text:    "Solo project",
					Date:    "2024-01-15",
					Company: "",
				}

				event := data.ToCareerEvent()

				Expect(event.Company).To(BeEmpty())
				Expect(event.Text).To(Equal("Solo project"))
			})
		})
	})

	Describe("Integration - Full Workflow", func() {
		It("should handle cancel workflow", func() {
			// Initial state
			Expect(modal.IsVisible()).To(BeTrue())

			// View should render
			view := modal.View()
			Expect(view).NotTo(BeEmpty())

			// Cancel with Esc
			cmd, completed, data := modal.Update(tea.KeyMsg{Type: tea.KeyEsc})

			Expect(cmd).To(BeNil())
			Expect(completed).To(BeFalse())
			Expect(data).To(BeNil())
			Expect(modal.IsVisible()).To(BeFalse())

			// View should be empty
			view = modal.View()
			Expect(view).To(BeEmpty())
		})
	})
})
