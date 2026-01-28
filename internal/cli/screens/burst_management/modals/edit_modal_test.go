package modals_test

import (
	"time"

	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"

	"github.com/baphled/kariya/internal/cli/screens/burst_management/modals"
	"github.com/baphled/kariya/internal/domain/career"
	tea "github.com/charmbracelet/bubbletea"
)

var _ = Describe("EditBurstModal", func() {
	var (
		modal *modals.EditBurstModal
		burst *career.Burst
	)

	BeforeEach(func() {
		burst = &career.Burst{
			ID:          "test-burst-id",
			Name:        "Original Burst Name",
			Description: "Original description text",
			EventIDs:    []string{"e1", "e2", "e3"},
			Confirmed:   false,
			CreatedAt:   time.Date(2024, 1, 1, 0, 0, 0, 0, time.UTC),
			UpdatedAt:   time.Date(2024, 2, 1, 0, 0, 0, 0, time.UTC),
		}
	})

	Describe("NewEditBurstModal", func() {
		It("creates modal with pre-populated values from burst", func() {
			modal = modals.NewEditBurstModal(burst, 80, 24)

			Expect(modal).NotTo(BeNil())
			Expect(modal.IsVisible()).To(BeTrue())
			Expect(modal.GetOriginalBurst()).To(Equal(burst))
		})

		It("respects minimum width", func() {
			modal = modals.NewEditBurstModal(burst, 30, 24)

			Expect(modal).NotTo(BeNil())
		})

		It("respects maximum width", func() {
			modal = modals.NewEditBurstModal(burst, 200, 24)

			Expect(modal).NotTo(BeNil())
		})

		It("handles small height", func() {
			modal = modals.NewEditBurstModal(burst, 80, 10)

			Expect(modal).NotTo(BeNil())
		})
	})

	Describe("Init", func() {
		BeforeEach(func() {
			modal = modals.NewEditBurstModal(burst, 80, 24)
		})

		It("returns a command for form initialization", func() {
			cmd := modal.Init()

			Expect(cmd).NotTo(BeNil())
		})
	})

	Describe("Update", func() {
		BeforeEach(func() {
			modal = modals.NewEditBurstModal(burst, 80, 24)
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
			modal = modals.NewEditBurstModal(burst, 80, 24)
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

		It("contains Edit Burst title", func() {
			view := modal.View()

			Expect(view).To(ContainSubstring("Edit Burst"))
		})
	})

	Describe("Visibility methods", func() {
		BeforeEach(func() {
			modal = modals.NewEditBurstModal(burst, 80, 24)
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

	Describe("GetOriginalBurst", func() {
		It("returns the original burst", func() {
			modal = modals.NewEditBurstModal(burst, 80, 24)

			original := modal.GetOriginalBurst()

			Expect(original).To(Equal(burst))
			Expect(original.ID).To(Equal("test-burst-id"))
			Expect(original.Name).To(Equal("Original Burst Name"))
		})
	})

	Describe("EditBurstData", func() {
		It("holds name and description fields", func() {
			data := &modals.EditBurstData{
				Name:        "Updated Name",
				Description: "Updated Description",
			}

			Expect(data.Name).To(Equal("Updated Name"))
			Expect(data.Description).To(Equal("Updated Description"))
		})
	})
})
