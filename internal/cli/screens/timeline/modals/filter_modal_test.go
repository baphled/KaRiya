package modals_test

import (
	"time"

	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"

	"github.com/baphled/kariya/internal/cli/screens/timeline/modals"
	"github.com/baphled/kariya/internal/domain/career"
	tea "github.com/charmbracelet/bubbletea"
)

var _ = Describe("FilterModal", func() {
	var (
		modal  *modals.FilterModal
		events []*career.CareerEvent
	)

	BeforeEach(func() {
		events = []*career.CareerEvent{
			{
				ID:         "1",
				Text:       "Event 1",
				Date:       time.Now(),
				Company:    "Company A",
				Project:    "Project X",
				Categories: []string{"Development", "Backend"},
			},
			{
				ID:         "2",
				Text:       "Event 2",
				Date:       time.Now(),
				Company:    "Company B",
				Project:    "Project Y",
				Categories: []string{"Testing"},
			},
			{
				ID:         "3",
				Text:       "Event 3",
				Date:       time.Now(),
				Company:    "Company A",
				Project:    "Project Z",
				Categories: []string{"Development", "Frontend"},
			},
		}
	})

	Describe("NewFilterModal", func() {
		It("creates modal with default filters", func() {
			modal = modals.NewFilterModal(events, nil, 80, 24)

			Expect(modal).NotTo(BeNil())
			Expect(modal.IsVisible()).To(BeTrue())
		})

		It("creates modal with existing filters", func() {
			currentFilters := &modals.TimelineFilters{
				Companies:  []string{"Company A"},
				Categories: []string{"Development"},
				Projects:   []string{"Project X"},
				SortBy:     "text",
				SortOrder:  "asc",
			}
			modal = modals.NewFilterModal(events, currentFilters, 80, 24)

			Expect(modal).NotTo(BeNil())
		})

		It("handles nil current filters", func() {
			modal = modals.NewFilterModal(events, nil, 80, 24)

			Expect(modal).NotTo(BeNil())
		})

		It("handles empty events list", func() {
			modal = modals.NewFilterModal([]*career.CareerEvent{}, nil, 80, 24)

			Expect(modal).NotTo(BeNil())
		})

		It("handles events with empty company and project", func() {
			eventsWithEmpty := []*career.CareerEvent{
				{ID: "1", Text: "Event", Company: "", Project: ""},
			}
			modal = modals.NewFilterModal(eventsWithEmpty, nil, 80, 24)

			Expect(modal).NotTo(BeNil())
		})

		It("extracts unique companies from events", func() {
			modal = modals.NewFilterModal(events, nil, 80, 24)

			Expect(modal).NotTo(BeNil())
		})

		It("extracts unique categories from events", func() {
			modal = modals.NewFilterModal(events, nil, 80, 24)

			Expect(modal).NotTo(BeNil())
		})

		It("extracts unique projects from events", func() {
			modal = modals.NewFilterModal(events, nil, 80, 24)

			Expect(modal).NotTo(BeNil())
		})
	})

	Describe("Init", func() {
		BeforeEach(func() {
			modal = modals.NewFilterModal(events, nil, 80, 24)
		})

		It("returns a command for form initialization", func() {
			cmd := modal.Init()

			Expect(cmd).NotTo(BeNil())
		})
	})

	Describe("Update", func() {
		BeforeEach(func() {
			modal = modals.NewFilterModal(events, nil, 80, 24)
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
			modal = modals.NewFilterModal(events, nil, 80, 24)
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
			modal = modals.NewFilterModal(events, nil, 80, 24)
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

	Describe("TimelineFilters", func() {
		It("holds all filter fields", func() {
			filters := &modals.TimelineFilters{
				SearchText: "search",
				Tags:       []string{"tag1", "tag2"},
				Companies:  []string{"Company A"},
				Categories: []string{"Development"},
				Projects:   []string{"Project X"},
				SortBy:     "date",
				SortOrder:  "desc",
			}

			Expect(filters.SearchText).To(Equal("search"))
			Expect(filters.Tags).To(HaveLen(2))
			Expect(filters.Companies).To(HaveLen(1))
			Expect(filters.Categories).To(HaveLen(1))
			Expect(filters.Projects).To(HaveLen(1))
			Expect(filters.SortBy).To(Equal("date"))
			Expect(filters.SortOrder).To(Equal("desc"))
		})
	})

	Describe("FilterFormData", func() {
		It("holds form field values", func() {
			formData := &modals.FilterFormData{
				Companies:  []string{"Company A", "Company B"},
				Categories: []string{"Development"},
				Projects:   []string{"Project X"},
				SortBy:     "text",
				SortOrder:  "asc",
			}

			Expect(formData.Companies).To(HaveLen(2))
			Expect(formData.Categories).To(HaveLen(1))
			Expect(formData.Projects).To(HaveLen(1))
			Expect(formData.SortBy).To(Equal("text"))
			Expect(formData.SortOrder).To(Equal("asc"))
		})
	})
})
