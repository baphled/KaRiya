package modals_test

import (
	"time"

	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"

	"github.com/baphled/kariya/internal/cli/screens/timeline/modals"
	tea "github.com/charmbracelet/bubbletea"
)

var _ = Describe("QuickAddModal", func() {
	var modal *modals.QuickAddModal

	Describe("NewQuickAddModal", func() {
		It("creates modal with default values", func() {
			modal = modals.NewQuickAddModal(80, 24)

			Expect(modal).NotTo(BeNil())
			Expect(modal.IsVisible()).To(BeTrue())
		})

		It("respects minimum width", func() {
			modal = modals.NewQuickAddModal(20, 24)

			Expect(modal).NotTo(BeNil())
		})

		It("respects maximum width", func() {
			modal = modals.NewQuickAddModal(200, 24)

			Expect(modal).NotTo(BeNil())
		})
	})

	Describe("Init", func() {
		BeforeEach(func() {
			modal = modals.NewQuickAddModal(80, 24)
		})

		It("returns a command for form initialization", func() {
			cmd := modal.Init()

			Expect(cmd).NotTo(BeNil())
		})
	})

	Describe("Update", func() {
		BeforeEach(func() {
			modal = modals.NewQuickAddModal(80, 24)
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
			modal = modals.NewQuickAddModal(80, 24)
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
			modal = modals.NewQuickAddModal(80, 24)
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

	Describe("QuickAddData", func() {
		Describe("ToCareerEvent", func() {
			It("converts form data to career event", func() {
				data := &modals.QuickAddData{
					Text: "New event text",
					Date: "2024-06-15",
				}

				event := data.ToCareerEvent()

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

			It("handles invalid date with fallback to today", func() {
				data := &modals.QuickAddData{
					Text: "Event with invalid date",
					Date: "invalid-date",
				}

				event := data.ToCareerEvent()

				Expect(event).NotTo(BeNil())
				Expect(event.Text).To(Equal("Event with invalid date"))
				Expect(event.Date).To(BeTemporally("~", time.Now(), time.Minute))
			})

			It("sets created and updated timestamps", func() {
				data := &modals.QuickAddData{
					Text: "Test event",
					Date: "2024-01-01",
				}

				event := data.ToCareerEvent()

				Expect(event.CreatedAt).To(BeTemporally("~", time.Now(), time.Minute))
				Expect(event.UpdatedAt).To(BeTemporally("~", time.Now(), time.Minute))
			})

			It("initializes empty slices for collections", func() {
				data := &modals.QuickAddData{
					Text: "Test",
					Date: "2024-01-01",
				}

				event := data.ToCareerEvent()

				Expect(event.Tags).NotTo(BeNil())
				Expect(event.Tags).To(BeEmpty())
				Expect(event.Categories).NotTo(BeNil())
				Expect(event.Categories).To(BeEmpty())
				Expect(event.Skills).NotTo(BeNil())
				Expect(event.Skills).To(BeEmpty())
			})
		})
	})
})
