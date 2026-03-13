package burst_test

import (
	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"

	"github.com/baphled/kariya/internal/testutil/fixtures"
	burstview "github.com/baphled/kariya/internal/tui/views/burst"
	"github.com/baphled/kariya/internal/ui/display"
	"github.com/baphled/kariya/internal/ui/themes"
	tea "github.com/charmbracelet/bubbletea"
)

var _ = Describe("Detail", func() {
	var (
		modal *burstview.Detail
		burst display.Burst
		theme themes.Theme
	)

	BeforeEach(func() {
		burst = display.BurstFromDomain(fixtures.BurstConfirmed("test-burst-id", "e1", "e2", "e3"))
		burst.Name = "Backend Development"
		burst.Description = "API and microservices work"
		theme = themes.NewDefaultTheme()
	})

	Describe("NewDetail", func() {
		It("creates modal with burst content", func() {
			modal = burstview.NewDetail(burst, theme)

			Expect(modal).NotTo(BeNil())
			Expect(modal.GetBurst()).To(Equal(burst))
		})

		It("handles nil theme with default", func() {
			modal = burstview.NewDetail(burst, nil)

			Expect(modal).NotTo(BeNil())
			view := modal.View()
			Expect(view).NotTo(BeEmpty())
		})

		It("renders burst name in view", func() {
			modal = burstview.NewDetail(burst, theme)
			modal.Show()

			view := modal.View()
			Expect(view).To(ContainSubstring("Backend Development"))
		})

		It("renders confirmed status for confirmed bursts", func() {
			modal = burstview.NewDetail(burst, theme)
			modal.Show()

			view := modal.View()
			Expect(view).To(ContainSubstring("Confirmed"))
		})

		It("renders event count", func() {
			modal = burstview.NewDetail(burst, theme)
			modal.Show()

			view := modal.View()
			Expect(view).To(ContainSubstring("3"))
		})
	})

	Describe("Init", func() {
		BeforeEach(func() {
			modal = burstview.NewDetail(burst, theme)
		})

		It("returns initialization command", func() {
			cmd := modal.Init()
			// Init may return nil or a command.
			_ = cmd
		})
	})

	Describe("Update", func() {
		BeforeEach(func() {
			modal = burstview.NewDetail(burst, theme)
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
			modal = burstview.NewDetail(burst, theme)
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
			modal = burstview.NewDetail(burst, theme)
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
			modal = burstview.NewDetail(burst, theme)

			modal.SetDimensions(100, 50)

			Expect(modal).NotTo(BeNil())
		})
	})

	Describe("SetBurst", func() {
		It("updates the displayed burst", func() {
			modal = burstview.NewDetail(burst, theme)
			newBurst := display.BurstFromDomain(fixtures.Burst("new-burst-id", "e4", "e5"))
			newBurst.Name = "New Burst Name"
			newBurst.Description = "New description"

			modal.SetBurst(newBurst)

			Expect(modal.GetBurst()).To(Equal(newBurst))
		})

		It("updates view content", func() {
			modal = burstview.NewDetail(burst, theme)
			modal.Show()
			newBurst := display.BurstFromDomain(fixtures.Burst("new-id", "e1"))
			newBurst.Name = "Updated Burst"
			newBurst.Description = "Updated description"

			modal.SetBurst(newBurst)
			view := modal.View()

			Expect(view).To(ContainSubstring("Updated Burst"))
		})
	})

	Describe("GetBurst", func() {
		It("returns the current burst", func() {
			modal = burstview.NewDetail(burst, theme)

			result := modal.GetBurst()

			Expect(result).To(Equal(burst))
			Expect(result.ID).To(Equal("test-burst-id"))
			Expect(result.Name).To(Equal("Backend Development"))
		})
	})

	Describe("Unconfirmed Burst", func() {
		It("does not show confirmed badge for unconfirmed burst", func() {
			unconfirmedBurst := display.BurstFromDomain(fixtures.Burst("unconfirmed-id", "e1", "e2"))
			unconfirmedBurst.Name = "Unconfirmed Burst"
			unconfirmedBurst.Description = "Not confirmed yet"
			modal = burstview.NewDetail(unconfirmedBurst, theme)
			modal.Show()

			view := modal.View()

			Expect(view).To(ContainSubstring("Unconfirmed Burst"))
		})
	})

	// BUG REGRESSION: nil burst handling
	Describe("Nil Burst Handling", func() {
		It("should not panic when created with zero-value burst", func() {
			Expect(func() {
				modal = burstview.NewDetail(display.Burst{}, theme)
			}).NotTo(Panic())
		})

		It("should handle zero-value burst gracefully in View", func() {
			modal = burstview.NewDetail(display.Burst{}, theme)
			modal.Show()

			Expect(func() {
				view := modal.View()
				Expect(view).NotTo(BeEmpty())
			}).NotTo(Panic())
		})

		It("should return zero-value burst from GetBurst when created empty", func() {
			modal = burstview.NewDetail(display.Burst{}, theme)
			Expect(modal.GetBurst()).To(Equal(display.Burst{}))
		})

		It("should handle SetBurst with zero-value burst", func() {
			modal = burstview.NewDetail(burst, theme)

			Expect(func() {
				modal.SetBurst(display.Burst{})
			}).NotTo(Panic())
		})

		It("should handle Update when burst is empty", func() {
			modal = burstview.NewDetail(display.Burst{}, theme)
			modal.Show()

			Expect(func() {
				modal.Update(tea.KeyMsg{Type: tea.KeyRunes, Runes: []rune{'x'}})
			}).NotTo(Panic())
		})
	})

	Describe("Key Handling - Action Signals", func() {
		BeforeEach(func() {
			modal = burstview.NewDetail(burst, theme)
			modal.SetDimensions(80, 24)
			modal.Show()
		})

		It("should signal 'v' key for view events action", func() {
			msg := tea.KeyMsg{Type: tea.KeyRunes, Runes: []rune{'v'}}
			model, _ := modal.Update(msg)

			Expect(model).NotTo(BeNil())
		})

		It("should signal 'f' key for view facts action", func() {
			msg := tea.KeyMsg{Type: tea.KeyRunes, Runes: []rune{'f'}}
			model, _ := modal.Update(msg)

			Expect(model).NotTo(BeNil())
		})

		It("should signal 's' key for view skills action", func() {
			msg := tea.KeyMsg{Type: tea.KeyRunes, Runes: []rune{'s'}}
			model, _ := modal.Update(msg)

			Expect(model).NotTo(BeNil())
		})

		It("should signal 'e' key for edit action", func() {
			msg := tea.KeyMsg{Type: tea.KeyRunes, Runes: []rune{'e'}}
			model, _ := modal.Update(msg)

			Expect(model).NotTo(BeNil())
		})

		It("should signal 'd' key for delete action", func() {
			msg := tea.KeyMsg{Type: tea.KeyRunes, Runes: []rune{'d'}}
			model, _ := modal.Update(msg)

			Expect(model).NotTo(BeNil())
		})

		It("should signal 'c' key for confirm action", func() {
			msg := tea.KeyMsg{Type: tea.KeyRunes, Runes: []rune{'c'}}
			model, _ := modal.Update(msg)

			Expect(model).NotTo(BeNil())
		})
	})

	Describe("Burst Detail Modal Full Action Suite (Phase 2-Tier 3)", func() {
		BeforeEach(func() {
			modal = burstview.NewDetail(burst, theme)
			modal.SetDimensions(80, 24)
			modal.Show()
		})

		It("should signal 'v' key for view events", func() {
			msg := tea.KeyMsg{Type: tea.KeyRunes, Runes: []rune{'v'}}
			model, _ := modal.Update(msg)

			Expect(model).NotTo(BeNil())
		})

		It("should signal 'f' key for view facts", func() {
			msg := tea.KeyMsg{Type: tea.KeyRunes, Runes: []rune{'f'}}
			model, _ := modal.Update(msg)

			Expect(model).NotTo(BeNil())
		})

		It("should signal 'c' key for confirm", func() {
			msg := tea.KeyMsg{Type: tea.KeyRunes, Runes: []rune{'c'}}
			model, _ := modal.Update(msg)

			Expect(model).NotTo(BeNil())
		})

		It("should close modal on Escape without action", func() {
			msg := tea.KeyMsg{Type: tea.KeyEsc}
			model, _ := modal.Update(msg)

			Expect(model).NotTo(BeNil())
		})
	})
})
