package event_test

import (
	"time"

	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"

	"github.com/baphled/kariya/internal/domain/capture"
	eventview "github.com/baphled/kariya/internal/tui/views/event"
	tea "github.com/charmbracelet/bubbletea"
)

var _ = Describe("QuickAdd", func() {
	var modal *eventview.QuickAdd

	Describe("NewQuickAddModal", func() {
		It("creates modal with default values", func() {
			modal = eventview.NewQuickAdd(80, 24)

			Expect(modal).NotTo(BeNil())
			Expect(modal.IsVisible()).To(BeTrue())
		})

		It("respects minimum width", func() {
			modal = eventview.NewQuickAdd(20, 24)

			Expect(modal).NotTo(BeNil())
		})

		It("respects maximum width", func() {
			modal = eventview.NewQuickAdd(200, 24)

			Expect(modal).NotTo(BeNil())
		})
	})

	Describe("Init", func() {
		BeforeEach(func() {
			modal = eventview.NewQuickAdd(80, 24)
		})

		It("returns a command for form initialization", func() {
			cmd := modal.Init()

			Expect(cmd).NotTo(BeNil())
		})

		It("returns nil for a zero-value modal without a form", func() {
			modal = &eventview.QuickAdd{}

			Expect(modal.Init()).To(BeNil())
		})
	})

	Describe("Update", func() {
		BeforeEach(func() {
			modal = eventview.NewQuickAdd(80, 24)
			modal.Init()
		})

		It("closes on escape without saving", func() {
			escMsg := tea.KeyMsg{Type: tea.KeyEsc}

			cmd, completed, data := modal.Update(escMsg)

			Expect(modal.IsVisible()).To(BeFalse())
			Expect(completed).To(BeFalse())
			Expect(data).To(BeNil())
			Expect(cmd).To(BeNil())
		})

		It("submits the current quick-add data on ctrl+s", func() {
			cmd, completed, data := modal.Update(tea.KeyMsg{Type: tea.KeyCtrlS})

			Expect(cmd).To(BeNil())
			Expect(completed).To(BeTrue())
			Expect(modal.IsVisible()).To(BeFalse())
			Expect(data).NotTo(BeNil())
			Expect(data.Text).To(BeEmpty())
			Expect(data.Date).To(MatchRegexp(`^\d{4}-\d{2}-\d{2}$`))
		})

		It("handles window resize", func() {
			resizeMsg := tea.WindowSizeMsg{Width: 100, Height: 30}

			cmd, completed, data := modal.Update(resizeMsg)

			Expect(modal.IsVisible()).To(BeTrue())
			Expect(completed).To(BeFalse())
			Expect(data).To(BeNil())
			Expect(cmd).NotTo(BeNil())
		})

		It("returns nil when not visible", func() {
			modal.Hide()

			cmd, completed, data := modal.Update(tea.KeyMsg{Type: tea.KeyEnter})

			Expect(completed).To(BeFalse())
			Expect(data).To(BeNil())
			Expect(cmd).To(BeNil())
		})

		It("forwards regular keys to form", func() {
			keyMsg := tea.KeyMsg{Type: tea.KeyRunes, Runes: []rune{'t'}}

			cmd, completed, data := modal.Update(keyMsg)

			Expect(modal.IsVisible()).To(BeTrue())
			Expect(completed).To(BeFalse())
			Expect(data).To(BeNil())
			_ = cmd
		})

		It("forwards tab key to form", func() {
			tabMsg := tea.KeyMsg{Type: tea.KeyTab}

			cmd, completed, data := modal.Update(tabMsg)

			Expect(modal.IsVisible()).To(BeTrue())
			Expect(completed).To(BeFalse())
			Expect(data).To(BeNil())
			_ = cmd
		})

		It("forwards ctrl+c to the form without completing", func() {
			cmd, completed, data := modal.Update(tea.KeyMsg{Type: tea.KeyCtrlC})

			Expect(modal.IsVisible()).To(BeTrue())
			Expect(completed).To(BeFalse())
			Expect(data).To(BeNil())
			_ = cmd
		})

		It("forwards enter key to form", func() {
			enterMsg := tea.KeyMsg{Type: tea.KeyEnter}

			cmd, completed, data := modal.Update(enterMsg)

			_ = cmd
			_ = completed
			_ = data
		})

		It("handles window resize with small width", func() {
			resizeMsg := tea.WindowSizeMsg{Width: 20, Height: 20}

			cmd, completed, data := modal.Update(resizeMsg)

			Expect(modal.IsVisible()).To(BeTrue())
			Expect(completed).To(BeFalse())
			Expect(data).To(BeNil())
			Expect(cmd).NotTo(BeNil())
		})

		It("handles window resize with large width", func() {
			resizeMsg := tea.WindowSizeMsg{Width: 200, Height: 50}

			cmd, completed, data := modal.Update(resizeMsg)

			Expect(modal.IsVisible()).To(BeTrue())
			Expect(completed).To(BeFalse())
			Expect(data).To(BeNil())
			Expect(cmd).NotTo(BeNil())
		})
	})

	Describe("View", func() {
		BeforeEach(func() {
			modal = eventview.NewQuickAdd(80, 24)
		})

		It("returns empty when not visible", func() {
			modal.Hide()

			view := modal.View()

			Expect(view).To(BeEmpty())
		})

		It("returns content when visible", func() {
			view := modal.View()

			Expect(view).NotTo(BeEmpty())
		})
	})

	Describe("Visibility methods", func() {
		BeforeEach(func() {
			modal = eventview.NewQuickAdd(80, 24)
		})

		It("Show makes modal visible", func() {
			modal.Hide()
			Expect(modal.IsVisible()).To(BeFalse())

			modal.Show()

			Expect(modal.IsVisible()).To(BeTrue())
		})

		It("Hide makes modal invisible", func() {
			Expect(modal.IsVisible()).To(BeTrue())

			modal.Hide()

			Expect(modal.IsVisible()).To(BeFalse())
		})
	})

	Describe("Modal Reopening", func() {
		It("should be re-openable after cancel", func() {
			modal = eventview.NewQuickAdd(80, 24)
			Expect(modal.IsVisible()).To(BeTrue())

			escMsg := tea.KeyMsg{Type: tea.KeyEsc}
			modal.Update(escMsg)
			Expect(modal.IsVisible()).To(BeFalse())

			modal.Show()
			Expect(modal.IsVisible()).To(BeTrue())
		})
	})

	Describe("QuickAddData", func() {
		Describe("domain conversion", func() {
			It("converts form data to career event through capture.NewEventFromInput", func() {
				data := &eventview.QuickAddData{
					Text: "New event text",
					Date: "2024-06-15",
				}

				event, err := capture.NewEventFromInput(capture.EventInput{
					Text: data.Text,
					Date: data.Date,
				})

				Expect(err).NotTo(HaveOccurred())
				Expect(event).NotTo(BeNil())
				Expect(event.Text).To(Equal("New event text"))
				Expect(event.Date.Year()).To(Equal(2024))
				Expect(event.Date.Month()).To(Equal(time.June))
				Expect(event.Date.Day()).To(Equal(15))
				Expect(event.Company).To(BeEmpty())
				Expect(event.Project).To(BeEmpty())
				Expect(event.Tags).To(BeEmpty())
				Expect(event.Categories).To(BeEmpty())
			})

			It("returns an error for invalid date input", func() {
				data := &eventview.QuickAddData{
					Text: "Event with invalid date",
					Date: "invalid-date",
				}

				event, err := capture.NewEventFromInput(capture.EventInput{
					Text: data.Text,
					Date: data.Date,
				})

				Expect(err).To(HaveOccurred())
				Expect(event).To(BeNil())
			})

			It("sets created and updated timestamps", func() {
				data := &eventview.QuickAddData{
					Text: "Test event",
					Date: "2024-01-01",
				}

				event, err := capture.NewEventFromInput(capture.EventInput{
					Text: data.Text,
					Date: data.Date,
				})

				Expect(err).NotTo(HaveOccurred())
				Expect(event.CreatedAt).To(BeTemporally("~", time.Now(), time.Minute))
				Expect(event.UpdatedAt).To(BeTemporally("~", time.Now(), time.Minute))
			})

			It("preserves omitted collections as nil zero values", func() {
				data := &eventview.QuickAddData{
					Text: "Test",
					Date: "2024-01-01",
				}

				event, err := capture.NewEventFromInput(capture.EventInput{
					Text: data.Text,
					Date: data.Date,
				})

				Expect(err).NotTo(HaveOccurred())
				Expect(event.Tags).To(BeNil())
				Expect(event.Categories).To(BeNil())
				Expect(event.Skills).To(BeNil())
			})
		})
	})
})
