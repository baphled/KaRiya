package event_test

import (
	"time"

	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"

	"github.com/baphled/kariya/internal/domain/capture"
	"github.com/baphled/kariya/internal/testutil/fixtures"
	eventview "github.com/baphled/kariya/internal/tui/views/event"
	"github.com/baphled/kariya/internal/ui/display"
	"github.com/charmbracelet/bubbles/cursor"
	tea "github.com/charmbracelet/bubbletea"
)

var _ = Describe("Edit", func() {
	var (
		modal *eventview.Edit
		event display.Event
	)

	BeforeEach(func() {
		domainEvent := fixtures.EventWith("test-event-id", "Original event text", "Test Company", "Test Project")
		domainEvent.Date = time.Date(2024, 3, 15, 0, 0, 0, 0, time.UTC)
		domainEvent.Tags = []string{"tag1", "tag2"}
		domainEvent.Categories = []string{"category1"}
		domainEvent.CreatedAt = time.Date(2024, 1, 1, 0, 0, 0, 0, time.UTC)
		domainEvent.UpdatedAt = time.Date(2024, 2, 1, 0, 0, 0, 0, time.UTC)
		event = display.EventFromDomain(domainEvent)
	})

	Describe("NewEditModal", func() {
		It("creates modal with pre-populated values from event", func() {
			modal = eventview.NewEdit(event, 80, 24)

			Expect(modal).NotTo(BeNil())
			Expect(modal.IsVisible()).To(BeTrue())
			Expect(modal.GetOriginalEvent()).To(Equal(event))
		})

		It("respects minimum width", func() {
			modal = eventview.NewEdit(event, 30, 24)

			Expect(modal).NotTo(BeNil())
		})

		It("respects maximum width", func() {
			modal = eventview.NewEdit(event, 200, 24)

			Expect(modal).NotTo(BeNil())
		})
	})

	Describe("Init", func() {
		BeforeEach(func() {
			modal = eventview.NewEdit(event, 80, 24)
		})

		It("returns a command for form initialization", func() {
			cmd := modal.Init()

			Expect(cmd).NotTo(BeNil())
		})

		It("returns nil for a zero-value modal without a form", func() {
			modal = &eventview.Edit{}

			Expect(modal.Init()).To(BeNil())
		})
	})

	Describe("Update", func() {
		BeforeEach(func() {
			modal = eventview.NewEdit(event, 80, 24)
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

		It("submits the pre-populated data on ctrl+s", func() {
			cmd, completed, data := modal.Update(tea.KeyMsg{Type: tea.KeyCtrlS})

			Expect(cmd).To(BeNil())
			Expect(completed).To(BeTrue())
			Expect(modal.IsVisible()).To(BeFalse())
			Expect(data).NotTo(BeNil())
			Expect(data.Text).To(Equal(event.Text))
			Expect(data.Date).To(Equal("2024-03-15"))
			Expect(data.Company).To(Equal(event.Company))
			Expect(data.Project).To(Equal(event.Project))
			Expect(data.Tags).To(Equal(event.Tags))
			Expect(data.Categories).To(Equal(event.Categories))
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
			keyMsg := tea.KeyMsg{Type: tea.KeyRunes, Runes: []rune{'a'}}

			cmd, completed, data := modal.Update(keyMsg)

			// Form receives the key and may return a command.
			Expect(modal.IsVisible()).To(BeTrue())
			Expect(completed).To(BeFalse())
			Expect(data).To(BeNil())
			_ = cmd // cmd may be nil or not depending on form state
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

		It("submits after advancing through the form with enter", func() {
			runCmd := func(cmd tea.Cmd) tea.Msg {
				if cmd == nil {
					return nil
				}

				messages := make(chan tea.Msg, 1)
				go func() {
					messages <- cmd()
				}()

				select {
				case msg := <-messages:
					return msg
				case <-time.After(5 * time.Millisecond):
					return nil
				}
			}

			var (
				completed bool
				data      *eventview.EditData
			)

			for attempts := 0; attempts < 8 && !completed; attempts++ {
				cmd, done, resultData := modal.Update(tea.KeyMsg{Type: tea.KeyEnter})
				completed = done
				data = resultData

				for steps := 0; steps < 4 && cmd != nil && !completed; steps++ {
					msg := runCmd(cmd)
					switch msg.(type) {
					case nil, cursor.BlinkMsg:
						cmd = nil
					default:
						cmd, completed, data = modal.Update(msg)
					}
				}
			}

			Expect(completed).To(BeTrue())
			Expect(data).NotTo(BeNil())
			Expect(modal.IsVisible()).To(BeFalse())
			Expect(data.Text).To(Equal(event.Text))
		})

		It("forwards enter key to form", func() {
			enterMsg := tea.KeyMsg{Type: tea.KeyEnter}

			cmd, completed, data := modal.Update(enterMsg)

			// Form may or may not complete based on state.
			_ = cmd
			_ = completed
			_ = data
		})

		It("handles window resize with small width", func() {
			resizeMsg := tea.WindowSizeMsg{Width: 30, Height: 20}

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
			modal = eventview.NewEdit(event, 80, 24)
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
			modal = eventview.NewEdit(event, 80, 24)
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

	Describe("GetOriginalEvent", func() {
		It("returns the original event", func() {
			modal = eventview.NewEdit(event, 80, 24)

			original := modal.GetOriginalEvent()

			Expect(original).To(Equal(event))
			Expect(original.ID).To(Equal("test-event-id"))
			Expect(original.Text).To(Equal("Original event text"))
		})
	})

	Describe("EditData", func() {
		Describe("domain conversion", func() {
			It("converts form data with preserved ID through capture.UpdateEventFromInput", func() {
				data := &eventview.EditData{
					Text:       "Updated event text",
					Date:       "2024-06-15",
					Company:    "New Company",
					Project:    "New Project",
					Tags:       []string{"newtag"},
					Categories: []string{"newcat"},
				}

				result, err := capture.UpdateEventFromInput(capture.EditEventInput{
					EventID:    "original-id",
					Text:       data.Text,
					Date:       data.Date,
					Company:    data.Company,
					Project:    data.Project,
					Tags:       data.Tags,
					Categories: data.Categories,
					CreatedAt:  event.CreatedAt,
					UpdatedAt:  event.UpdatedAt,
				})

				Expect(err).NotTo(HaveOccurred())
				Expect(result).NotTo(BeNil())
				Expect(result.ID).To(Equal("original-id"))
				Expect(result.Text).To(Equal("Updated event text"))
				Expect(result.Date.Year()).To(Equal(2024))
				Expect(result.Date.Month()).To(Equal(time.June))
				Expect(result.Date.Day()).To(Equal(15))
				Expect(result.Company).To(Equal("New Company"))
				Expect(result.Project).To(Equal("New Project"))
				Expect(result.Tags).To(ContainElement("newtag"))
				Expect(result.Categories).To(ContainElement("newcat"))
				Expect(result.CreatedAt).To(Equal(event.CreatedAt))
				Expect(result.UpdatedAt).To(Equal(event.UpdatedAt))
			})

			It("returns an error for invalid date input", func() {
				data := &eventview.EditData{
					Text: "Event with bad date",
					Date: "not-a-date",
				}

				result, err := capture.UpdateEventFromInput(capture.EditEventInput{
					EventID: "id",
					Text:    data.Text,
					Date:    data.Date,
				})

				Expect(err).To(HaveOccurred())
				Expect(result).To(BeNil())
			})

			It("handles today keyword", func() {
				data := &eventview.EditData{
					Text: "Today's event",
					Date: "today",
				}

				result, err := capture.UpdateEventFromInput(capture.EditEventInput{
					EventID: "id",
					Text:    data.Text,
					Date:    data.Date,
				})

				Expect(err).NotTo(HaveOccurred())
				Expect(result).NotTo(BeNil())
				Expect(result.Date).To(BeTemporally("~", time.Now(), time.Minute))
			})

			It("initializes empty skills slice", func() {
				data := &eventview.EditData{
					Text: "Test",
					Date: "2024-01-01",
				}

				result, err := capture.UpdateEventFromInput(capture.EditEventInput{
					EventID: "id",
					Text:    data.Text,
					Date:    data.Date,
				})

				Expect(err).NotTo(HaveOccurred())
				Expect(result.Skills).NotTo(BeNil())
				Expect(result.Skills).To(BeEmpty())
			})

			It("preserves CreatedAt from original when valid", func() {
				data := &eventview.EditData{
					Text: "Test",
					Date: "2024-01-01",
				}
				createdAt := time.Date(2023, 6, 15, 0, 0, 0, 0, time.UTC)

				result, err := capture.UpdateEventFromInput(capture.EditEventInput{
					EventID:   "id",
					Text:      data.Text,
					Date:      data.Date,
					CreatedAt: createdAt,
				})

				Expect(err).NotTo(HaveOccurred())
				Expect(result.CreatedAt.Year()).To(Equal(2023))
				Expect(result.CreatedAt.Month()).To(Equal(time.June))
			})
		})

		Describe("Modal Reopening", func() {
			BeforeEach(func() {
				event = display.EventFromDomain(
					fixtures.EventWith("test-event-id", "Original event text", "Test Company", "Test Project"),
				)
			})

			It("should be re-openable after cancel", func() {
				modal = eventview.NewEdit(event, 80, 24)
				Expect(modal.IsVisible()).To(BeTrue())

				escMsg := tea.KeyMsg{Type: tea.KeyEsc}
				modal.Update(escMsg)
				Expect(modal.IsVisible()).To(BeFalse())

				modal.Show()
				Expect(modal.IsVisible()).To(BeTrue())
			})
		})
	})
})
