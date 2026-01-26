package modals_test

import (
	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"

	"github.com/baphled/kariya/internal/cli/screens/timeline/modals"
	"github.com/baphled/kariya/internal/domain/career"
	tea "github.com/charmbracelet/bubbletea"
)

var _ = Describe("SortModal", func() {
	var (
		modal  *modals.SortModal
		events []*career.CareerEvent
	)

	BeforeEach(func() {
		events = []*career.CareerEvent{
			{ID: "1", Text: "Event 1", Company: "Company A"},
			{ID: "2", Text: "Event 2", Company: "Company B"},
		}
	})

	Describe("NewSortModal", func() {
		It("creates modal with default sort config", func() {
			modal = modals.NewSortModal(events, nil, 80, 24)

			Expect(modal).NotTo(BeNil())
			Expect(modal.IsVisible()).To(BeTrue())

			config := modal.ToSortConfig()
			Expect(config.SortBy).To(Equal("date"))
			Expect(config.SortOrder).To(Equal("desc"))
		})

		It("creates modal with existing sort config", func() {
			current := &modals.SortConfig{
				SortBy:    "company",
				SortOrder: "asc",
			}
			modal = modals.NewSortModal(events, current, 80, 24)

			config := modal.ToSortConfig()
			Expect(config.SortBy).To(Equal("company"))
			Expect(config.SortOrder).To(Equal("asc"))
		})

		It("handles nil current config", func() {
			modal = modals.NewSortModal(events, nil, 80, 24)

			Expect(modal).NotTo(BeNil())
		})

		It("handles empty events list", func() {
			modal = modals.NewSortModal([]*career.CareerEvent{}, nil, 80, 24)

			Expect(modal).NotTo(BeNil())
		})

		It("handles config with only SortBy set", func() {
			current := &modals.SortConfig{
				SortBy:    "category",
				SortOrder: "",
			}
			modal = modals.NewSortModal(events, current, 80, 24)

			config := modal.ToSortConfig()
			Expect(config.SortBy).To(Equal("category"))
			Expect(config.SortOrder).To(Equal("desc")) // Default
		})

		It("handles config with only SortOrder set", func() {
			current := &modals.SortConfig{
				SortBy:    "",
				SortOrder: "asc",
			}
			modal = modals.NewSortModal(events, current, 80, 24)

			config := modal.ToSortConfig()
			Expect(config.SortBy).To(Equal("date")) // Default
			Expect(config.SortOrder).To(Equal("asc"))
		})

		It("respects minimum width constraint", func() {
			modal = modals.NewSortModal(events, nil, 20, 24)

			Expect(modal).NotTo(BeNil())
			Expect(modal.View()).NotTo(BeEmpty())
		})

		It("respects maximum width constraint", func() {
			modal = modals.NewSortModal(events, nil, 200, 24)

			Expect(modal).NotTo(BeNil())
			Expect(modal.View()).NotTo(BeEmpty())
		})

		It("handles width at minimum boundary", func() {
			modal = modals.NewSortModal(events, nil, 40, 24)

			Expect(modal).NotTo(BeNil())
		})

		It("handles width at maximum boundary", func() {
			modal = modals.NewSortModal(events, nil, 100, 24)

			Expect(modal).NotTo(BeNil())
		})
	})

	Describe("Init", func() {
		BeforeEach(func() {
			modal = modals.NewSortModal(events, nil, 80, 24)
		})

		It("returns a command for form initialization", func() {
			cmd := modal.Init()

			Expect(cmd).NotTo(BeNil())
		})
	})

	Describe("Update", func() {
		BeforeEach(func() {
			modal = modals.NewSortModal(events, nil, 80, 24)
			modal.Init()
		})

		It("closes on escape without applying", func() {
			escMsg := tea.KeyMsg{Type: tea.KeyEsc}

			cmd, applied, data := modal.Update(escMsg)

			Expect(modal.IsVisible()).To(BeFalse())
			Expect(applied).To(BeFalse())
			Expect(data).To(BeNil())
			Expect(cmd).To(BeNil())
		})

		It("handles window resize", func() {
			resizeMsg := tea.WindowSizeMsg{Width: 100, Height: 30}

			cmd, applied, data := modal.Update(resizeMsg)

			Expect(modal.IsVisible()).To(BeTrue())
			Expect(applied).To(BeFalse())
			Expect(data).To(BeNil())
			Expect(cmd).To(BeNil())
		})

		It("returns nil when not visible", func() {
			modal.Hide()

			cmd, applied, data := modal.Update(tea.KeyMsg{Type: tea.KeyEnter})

			Expect(applied).To(BeFalse())
			Expect(data).To(BeNil())
			Expect(cmd).To(BeNil())
		})

		It("forwards regular keys to form", func() {
			keyMsg := tea.KeyMsg{Type: tea.KeyRunes, Runes: []rune{'j'}}

			cmd, applied, data := modal.Update(keyMsg)

			Expect(modal.IsVisible()).To(BeTrue())
			Expect(applied).To(BeFalse())
			Expect(data).To(BeNil())
			_ = cmd
		})

		It("forwards tab key to form", func() {
			tabMsg := tea.KeyMsg{Type: tea.KeyTab}

			cmd, applied, data := modal.Update(tabMsg)

			Expect(modal.IsVisible()).To(BeTrue())
			Expect(applied).To(BeFalse())
			Expect(data).To(BeNil())
			_ = cmd
		})

		It("forwards enter key to form", func() {
			enterMsg := tea.KeyMsg{Type: tea.KeyEnter}

			cmd, applied, data := modal.Update(enterMsg)

			_ = cmd
			_ = applied
			_ = data
		})

		It("forwards up/down keys to form for selection", func() {
			downMsg := tea.KeyMsg{Type: tea.KeyDown}

			cmd, applied, data := modal.Update(downMsg)

			Expect(modal.IsVisible()).To(BeTrue())
			Expect(applied).To(BeFalse())
			Expect(data).To(BeNil())
			_ = cmd
		})
	})

	Describe("View", func() {
		BeforeEach(func() {
			modal = modals.NewSortModal(events, nil, 80, 24)
		})

		It("returns empty when not visible", func() {
			modal.Hide()

			view := modal.View()

			Expect(view).To(BeEmpty())
		})

		It("returns content when visible", func() {
			view := modal.View()

			Expect(view).NotTo(BeEmpty())
			Expect(view).To(ContainSubstring("Apply"))
		})
	})

	Describe("Visibility methods", func() {
		BeforeEach(func() {
			modal = modals.NewSortModal(events, nil, 80, 24)
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

	Describe("ToSortConfig", func() {
		It("returns current sort configuration", func() {
			current := &modals.SortConfig{
				SortBy:    "category",
				SortOrder: "asc",
			}
			modal = modals.NewSortModal(events, current, 80, 24)

			config := modal.ToSortConfig()

			Expect(config).NotTo(BeNil())
			Expect(config.SortBy).To(Equal("category"))
			Expect(config.SortOrder).To(Equal("asc"))
		})
	})

	Describe("RenderOverlay", func() {
		BeforeEach(func() {
			modal = modals.NewSortModal(events, nil, 80, 24)
		})

		It("returns base view when not visible", func() {
			modal.Hide()
			baseView := "Base Content"

			result := modal.RenderOverlay(baseView)

			Expect(result).To(Equal(baseView))
		})

		It("renders overlay when visible", func() {
			baseView := "Base Content"

			result := modal.RenderOverlay(baseView)

			Expect(result).NotTo(Equal(baseView))
		})
	})
})
