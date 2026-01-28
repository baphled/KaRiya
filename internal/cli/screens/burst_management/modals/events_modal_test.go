package modals_test

import (
	"time"

	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"

	"github.com/baphled/kariya/internal/cli/screens/burst_management/modals"
	"github.com/baphled/kariya/internal/cli/themes"
	"github.com/baphled/kariya/internal/domain/career"
	tea "github.com/charmbracelet/bubbletea"
)

var _ = Describe("BurstEventsModal", func() {
	var (
		modal     *modals.BurstEventsModal
		burstID   string
		burstName string
		events    []*career.CareerEvent
		theme     themes.Theme
	)

	BeforeEach(func() {
		burstID = "test-burst-id"
		burstName = "Backend Development"
		events = []*career.CareerEvent{
			{
				ID:      "e1",
				Text:    "Built microservices architecture",
				Date:    time.Date(2024, 3, 15, 0, 0, 0, 0, time.UTC),
				Company: "TechCorp",
				Tags:    []string{"backend", "architecture"},
			},
			{
				ID:      "e2",
				Text:    "Implemented API gateway",
				Date:    time.Date(2024, 3, 20, 0, 0, 0, 0, time.UTC),
				Company: "TechCorp",
				Tags:    []string{"backend", "api"},
			},
			{
				ID:      "e3",
				Text:    "Optimized database queries",
				Date:    time.Date(2024, 3, 25, 0, 0, 0, 0, time.UTC),
				Company: "TechCorp",
				Tags:    []string{"database", "performance"},
			},
		}
		theme = themes.NewDefaultTheme()
	})

	Describe("NewBurstEventsModal", func() {
		It("creates modal with events content", func() {
			modal = modals.NewBurstEventsModal(burstID, burstName, events, theme)

			Expect(modal).NotTo(BeNil())
			Expect(modal.GetBurstID()).To(Equal(burstID))
		})

		It("handles nil theme with default", func() {
			modal = modals.NewBurstEventsModal(burstID, burstName, events, nil)

			Expect(modal).NotTo(BeNil())
		})

		It("handles empty events list", func() {
			modal = modals.NewBurstEventsModal(burstID, burstName, []*career.CareerEvent{}, theme)

			Expect(modal).NotTo(BeNil())
			modal.Show()
			view := modal.View()
			Expect(view).To(ContainSubstring("No events"))
		})

		It("handles nil events list", func() {
			modal = modals.NewBurstEventsModal(burstID, burstName, nil, theme)

			Expect(modal).NotTo(BeNil())
		})

		It("includes event count in title", func() {
			modal = modals.NewBurstEventsModal(burstID, burstName, events, theme)
			modal.Show()

			view := modal.View()
			Expect(view).To(ContainSubstring("3"))
		})

		It("includes burst name in title", func() {
			modal = modals.NewBurstEventsModal(burstID, burstName, events, theme)
			modal.Show()

			view := modal.View()
			Expect(view).To(ContainSubstring("Backend Development"))
		})
	})

	Describe("Init", func() {
		BeforeEach(func() {
			modal = modals.NewBurstEventsModal(burstID, burstName, events, theme)
		})

		It("returns initialization command", func() {
			cmd := modal.Init()
			_ = cmd
		})
	})

	Describe("Update", func() {
		BeforeEach(func() {
			modal = modals.NewBurstEventsModal(burstID, burstName, events, theme)
			modal.SetDimensions(80, 24)
			modal.Show()
		})

		It("delegates to underlying modal", func() {
			msg := tea.KeyMsg{Type: tea.KeyRunes, Runes: []rune{'x'}}

			model, cmd := modal.Update(msg)

			Expect(model).NotTo(BeNil())
			_ = cmd
		})

		It("handles window size message", func() {
			msg := tea.WindowSizeMsg{Width: 100, Height: 30}

			model, cmd := modal.Update(msg)

			Expect(model).NotTo(BeNil())
			_ = cmd
		})

		It("handles escape key", func() {
			msg := tea.KeyMsg{Type: tea.KeyEsc}

			model, cmd := modal.Update(msg)

			Expect(model).NotTo(BeNil())
			_ = cmd
		})
	})

	Describe("View", func() {
		BeforeEach(func() {
			modal = modals.NewBurstEventsModal(burstID, burstName, events, theme)
			modal.SetDimensions(80, 24)
		})

		It("renders events when visible", func() {
			modal.Show()

			view := modal.View()

			Expect(view).NotTo(BeEmpty())
		})

		It("shows event text", func() {
			modal.Show()

			view := modal.View()

			Expect(view).To(ContainSubstring("Built microservices"))
		})

		It("shows event dates", func() {
			modal.Show()

			view := modal.View()

			Expect(view).To(ContainSubstring("2024-03-15"))
		})

		It("shows event tags", func() {
			modal.Show()

			view := modal.View()

			Expect(view).To(ContainSubstring("backend"))
		})

		It("numbers events", func() {
			modal.Show()

			view := modal.View()

			// Modal may paginate, so we only check for first items visible in viewport.
			Expect(view).To(ContainSubstring("1."))
			Expect(view).To(ContainSubstring("2."))
		})
	})

	Describe("Visibility methods", func() {
		BeforeEach(func() {
			modal = modals.NewBurstEventsModal(burstID, burstName, events, theme)
		})

		It("Show makes modal visible", func() {
			modal.Hide()
			Expect(modal.IsVisible()).To(BeFalse())

			modal.Show()

			Expect(modal.IsVisible()).To(BeTrue())
		})

		It("Hide makes modal invisible", func() {
			modal.Show()
			Expect(modal.IsVisible()).To(BeTrue())

			modal.Hide()

			Expect(modal.IsVisible()).To(BeFalse())
		})
	})

	Describe("SetDimensions", func() {
		It("sets terminal dimensions", func() {
			modal = modals.NewBurstEventsModal(burstID, burstName, events, theme)

			modal.SetDimensions(100, 50)

			Expect(modal).NotTo(BeNil())
		})
	})

	Describe("SetEvents", func() {
		It("updates the displayed events", func() {
			modal = modals.NewBurstEventsModal(burstID, burstName, events, theme)
			newEvents := []*career.CareerEvent{
				{ID: "new-1", Text: "New event one", Date: time.Now()},
				{ID: "new-2", Text: "New event two", Date: time.Now()},
			}

			modal.SetEvents(newEvents)
			modal.Show()
			view := modal.View()

			Expect(view).To(ContainSubstring("New event one"))
		})

		It("updates event count in view", func() {
			modal = modals.NewBurstEventsModal(burstID, burstName, events, theme)
			modal.Show()

			// Initially 3 events.
			view := modal.View()
			Expect(view).To(ContainSubstring("3"))

			// Update to 1 event.
			modal.SetEvents([]*career.CareerEvent{
				{ID: "single", Text: "Single event", Date: time.Now()},
			})
			view = modal.View()
			Expect(view).To(ContainSubstring("1"))
		})
	})

	Describe("GetBurstID", func() {
		It("returns the burst ID", func() {
			modal = modals.NewBurstEventsModal(burstID, burstName, events, theme)

			result := modal.GetBurstID()

			Expect(result).To(Equal(burstID))
		})
	})
})
