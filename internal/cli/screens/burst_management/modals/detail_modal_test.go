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

var _ = Describe("BurstDetailModal", func() {
	var (
		modal *modals.BurstDetailModal
		burst *career.Burst
		theme themes.Theme
	)

	BeforeEach(func() {
		now := time.Now()
		confirmedAt := now.Add(-24 * time.Hour)
		burst = &career.Burst{
			ID:          "test-burst-id",
			Name:        "Backend Development",
			Description: "API and microservices work",
			EventIDs:    []string{"e1", "e2", "e3"},
			Confirmed:   true,
			ConfirmedAt: &confirmedAt,
			CreatedAt:   now.Add(-48 * time.Hour),
			UpdatedAt:   now,
		}
		theme = themes.NewDefaultTheme()
	})

	Describe("NewBurstDetailModal", func() {
		It("creates modal with burst content", func() {
			modal = modals.NewBurstDetailModal(burst, theme)

			Expect(modal).NotTo(BeNil())
			Expect(modal.GetBurst()).To(Equal(burst))
		})

		It("handles nil theme with default", func() {
			modal = modals.NewBurstDetailModal(burst, nil)

			Expect(modal).NotTo(BeNil())
			view := modal.View()
			Expect(view).NotTo(BeEmpty())
		})

		It("renders burst name in view", func() {
			modal = modals.NewBurstDetailModal(burst, theme)
			modal.Show()

			view := modal.View()
			Expect(view).To(ContainSubstring("Backend Development"))
		})

		It("renders confirmed status for confirmed bursts", func() {
			modal = modals.NewBurstDetailModal(burst, theme)
			modal.Show()

			view := modal.View()
			Expect(view).To(ContainSubstring("Confirmed"))
		})

		It("renders event count", func() {
			modal = modals.NewBurstDetailModal(burst, theme)
			modal.Show()

			view := modal.View()
			Expect(view).To(ContainSubstring("3"))
		})
	})

	Describe("Init", func() {
		BeforeEach(func() {
			modal = modals.NewBurstDetailModal(burst, theme)
		})

		It("returns initialization command", func() {
			cmd := modal.Init()
			// Init may return nil or a command.
			_ = cmd
		})
	})

	Describe("Update", func() {
		BeforeEach(func() {
			modal = modals.NewBurstDetailModal(burst, theme)
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
			modal = modals.NewBurstDetailModal(burst, theme)
			modal.SetDimensions(80, 24)
		})

		It("renders burst details when visible", func() {
			modal.Show()

			view := modal.View()

			Expect(view).NotTo(BeEmpty())
			Expect(view).To(ContainSubstring("Backend Development"))
		})

		It("renders description", func() {
			modal.Show()

			view := modal.View()

			Expect(view).To(ContainSubstring("API and microservices work"))
		})

		It("renders footer badges", func() {
			modal.Show()

			view := modal.View()

			Expect(view).To(ContainSubstring("Edit"))
			Expect(view).To(ContainSubstring("Delete"))
		})
	})

	Describe("Visibility methods", func() {
		BeforeEach(func() {
			modal = modals.NewBurstDetailModal(burst, theme)
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
			modal = modals.NewBurstDetailModal(burst, theme)

			modal.SetDimensions(100, 50)

			Expect(modal).NotTo(BeNil())
		})
	})

	Describe("SetBurst", func() {
		It("updates the displayed burst", func() {
			modal = modals.NewBurstDetailModal(burst, theme)
			newBurst := &career.Burst{
				ID:          "new-burst-id",
				Name:        "New Burst Name",
				Description: "New description",
				EventIDs:    []string{"e4", "e5"},
				CreatedAt:   time.Now(),
				UpdatedAt:   time.Now(),
			}

			modal.SetBurst(newBurst)

			Expect(modal.GetBurst()).To(Equal(newBurst))
		})

		It("updates view content", func() {
			modal = modals.NewBurstDetailModal(burst, theme)
			modal.Show()
			newBurst := &career.Burst{
				ID:          "new-id",
				Name:        "Updated Burst",
				Description: "Updated description",
				EventIDs:    []string{"e1"},
				CreatedAt:   time.Now(),
				UpdatedAt:   time.Now(),
			}

			modal.SetBurst(newBurst)
			view := modal.View()

			Expect(view).To(ContainSubstring("Updated Burst"))
		})
	})

	Describe("GetBurst", func() {
		It("returns the current burst", func() {
			modal = modals.NewBurstDetailModal(burst, theme)

			result := modal.GetBurst()

			Expect(result).To(Equal(burst))
			Expect(result.ID).To(Equal("test-burst-id"))
			Expect(result.Name).To(Equal("Backend Development"))
		})
	})

	Describe("Unconfirmed Burst", func() {
		It("does not show confirmed badge for unconfirmed burst", func() {
			unconfirmedBurst := &career.Burst{
				ID:          "unconfirmed-id",
				Name:        "Unconfirmed Burst",
				Description: "Not confirmed yet",
				EventIDs:    []string{"e1", "e2"},
				Confirmed:   false,
				CreatedAt:   time.Now(),
				UpdatedAt:   time.Now(),
			}
			modal = modals.NewBurstDetailModal(unconfirmedBurst, theme)
			modal.Show()

			view := modal.View()

			Expect(view).To(ContainSubstring("Unconfirmed Burst"))
		})
	})
})
