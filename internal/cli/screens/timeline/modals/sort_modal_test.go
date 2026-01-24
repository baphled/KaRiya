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
