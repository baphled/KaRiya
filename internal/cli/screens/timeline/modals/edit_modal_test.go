package modals_test

import (
	"time"

	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"

	"github.com/baphled/kariya/internal/cli/screens/timeline/modals"
	"github.com/baphled/kariya/internal/domain/career"
	"github.com/baphled/kariya/internal/testutil/fixtures"
	tea "github.com/charmbracelet/bubbletea"
)

var _ = Describe("EditModal", func() {
	var (
		modal *modals.EditModal
		event *career.Event
	)

	BeforeEach(func() {
		event = fixtures.EventWith("test-event-id", "Original event text", "Test Company", "Test Project")
		event.Date = time.Date(2024, 3, 15, 0, 0, 0, 0, time.UTC)
		event.Tags = []string{"tag1", "tag2"}
		event.Categories = []string{"category1"}
		event.CreatedAt = time.Date(2024, 1, 1, 0, 0, 0, 0, time.UTC)
		event.UpdatedAt = time.Date(2024, 2, 1, 0, 0, 0, 0, time.UTC)
	})

	Describe("NewEditModal", func() {
		It("creates modal with pre-populated values from event", func() {
			modal = modals.NewEditModal(event, 80, 24)

			Expect(modal).NotTo(BeNil())
			Expect(modal.IsVisible()).To(BeTrue())
			Expect(modal.GetOriginalEvent()).To(Equal(event))
		})

		It("respects minimum width", func() {
			modal = modals.NewEditModal(event, 30, 24)

			Expect(modal).NotTo(BeNil())
		})

		It("respects maximum width", func() {
			modal = modals.NewEditModal(event, 200, 24)

			Expect(modal).NotTo(BeNil())
		})
	})

	Describe("Init", func() {
		BeforeEach(func() {
			modal = modals.NewEditModal(event, 80, 24)
		})

		It("returns a command for form initialization", func() {
			cmd := modal.Init()

			Expect(cmd).NotTo(BeNil())
		})
	})

	Describe("Update", func() {
		BeforeEach(func() {
			modal = modals.NewEditModal(event, 80, 24)
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
			modal = modals.NewEditModal(event, 80, 24)
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
			modal = modals.NewEditModal(event, 80, 24)
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
			modal = modals.NewEditModal(event, 80, 24)

			original := modal.GetOriginalEvent()

			Expect(original).To(Equal(event))
			Expect(original.ID).To(Equal("test-event-id"))
			Expect(original.Text).To(Equal("Original event text"))
		})
	})

	Describe("EditData", func() {
		Describe("ToCareerEvent", func() {
			It("converts form data to career event with preserved ID", func() {
				data := &modals.EditData{
					Text:       "Updated event text",
					Date:       "2024-06-15",
					Company:    "New Company",
					Project:    "New Project",
					Tags:       []string{"newtag"},
					Categories: []string{"newcat"},
				}

				result := data.ToCareerEvent(
					"original-id",
					event.CreatedAt,
					event.UpdatedAt,
				)

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
			})

			It("handles invalid date with fallback", func() {
				data := &modals.EditData{
					Text: "Event with bad date",
					Date: "not-a-date",
				}

				result := data.ToCareerEvent("id", nil, nil)

				Expect(result).NotTo(BeNil())
				Expect(result.Date).To(BeTemporally("~", time.Now(), time.Minute))
			})

			It("handles today keyword", func() {
				data := &modals.EditData{
					Text: "Today's event",
					Date: "today",
				}

				result := data.ToCareerEvent("id", nil, nil)

				Expect(result).NotTo(BeNil())
				Expect(result.Date).To(BeTemporally("~", time.Now(), time.Minute))
			})

			It("initializes empty skills slice", func() {
				data := &modals.EditData{
					Text: "Test",
					Date: "2024-01-01",
				}

				result := data.ToCareerEvent("id", nil, nil)

				Expect(result.Skills).NotTo(BeNil())
				Expect(result.Skills).To(BeEmpty())
			})

			It("preserves CreatedAt from original when valid", func() {
				data := &modals.EditData{
					Text: "Test",
					Date: "2024-01-01",
				}
				createdAt := time.Date(2023, 6, 15, 0, 0, 0, 0, time.UTC)

				result := data.ToCareerEvent("id", createdAt, nil)

				Expect(result.CreatedAt.Year()).To(Equal(2023))
				Expect(result.CreatedAt.Month()).To(Equal(time.June))
			})
		})
	})
})
