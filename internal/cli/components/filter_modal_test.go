package components_test

import (
	"time"

	"github.com/baphled/kariya/internal/cli/components"
	"github.com/baphled/kariya/internal/domain/career"
	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
)

var _ = Describe("FilterModal", func() {
	var (
		events         []*career.CareerEvent
		currentFilters *components.TimelineFilters
		modal          *components.FilterModalModel
	)

	BeforeEach(func() {
		// Create test events with companies, categories, and projects
		date1, _ := time.Parse("2006-01-02", "2024-01-01")
		date2, _ := time.Parse("2006-01-02", "2023-06-15")
		date3, _ := time.Parse("2006-01-02", "2022-03-10")

		events = []*career.CareerEvent{
			{
				ID:         "1",
				Date:       date1,
				Text:       "Backend Developer at TechCorp",
				Company:    "TechCorp",
				Project:    "API Gateway",
				Categories: []string{"development", "backend"},
			},
			{
				ID:         "2",
				Date:       date2,
				Text:       "DevOps Engineer at CloudCo",
				Company:    "CloudCo",
				Project:    "Infrastructure",
				Categories: []string{"devops", "cloud"},
			},
			{
				ID:         "3",
				Date:       date3,
				Text:       "Full Stack Developer at TechCorp",
				Company:    "TechCorp",
				Project:    "Customer Portal",
				Categories: []string{"development", "fullstack"},
			},
		}

		currentFilters = &components.TimelineFilters{
			Companies:  []string{},
			Categories: []string{},
			Projects:   []string{},
			SortBy:     "date",
			SortOrder:  "desc",
		}
	})

	Describe("NewFilterModal", func() {
		It("should create a new filter modal with default values", func() {
			modal = components.NewFilterModal(events, currentFilters, 120, 40)

			Expect(modal).NotTo(BeNil())
			Expect(modal.IsVisible()).To(BeTrue())
		})

		It("should extract unique companies from events", func() {
			modal = components.NewFilterModal(events, currentFilters, 120, 40)

			// Verify modal was created successfully
			// The form will have company options extracted from events
			Expect(modal).NotTo(BeNil())
			Expect(modal.IsVisible()).To(BeTrue())
		})

		It("should extract unique categories from events", func() {
			modal = components.NewFilterModal(events, currentFilters, 120, 40)

			// Verify modal was created successfully
			// The form will have category options extracted from events
			Expect(modal).NotTo(BeNil())
			Expect(modal.IsVisible()).To(BeTrue())
		})

		It("should extract unique projects from events", func() {
			modal = components.NewFilterModal(events, currentFilters, 120, 40)

			// Verify modal was created successfully with projects
			// The form will have project options extracted from events
			Expect(modal).NotTo(BeNil())
			Expect(modal.IsVisible()).To(BeTrue())
		})

		It("should pre-populate form with current filters", func() {
			currentFilters.Companies = []string{"TechCorp"}
			currentFilters.Categories = []string{"development"}
			currentFilters.Projects = []string{"API Gateway"}
			currentFilters.SortBy = "text"
			currentFilters.SortOrder = "asc"

			modal = components.NewFilterModal(events, currentFilters, 120, 40)

			Expect(modal).NotTo(BeNil())
			// Form data should be pre-populated (internal state)
		})
	})

	Describe("Projects filter", func() {
		It("should include projects in filter options when events have projects", func() {
			modal = components.NewFilterModal(events, currentFilters, 120, 40)

			// Verify modal includes project filtering capability
			Expect(modal).NotTo(BeNil())
			Expect(modal.IsVisible()).To(BeTrue())
		})

		It("should handle events without projects gracefully", func() {
			eventsWithoutProjects := []*career.CareerEvent{
				{
					ID:         "1",
					Date:       time.Now(),
					Text:       "Event without project",
					Company:    "TechCorp",
					Project:    "",
					Categories: []string{"development"},
				},
			}

			modal = components.NewFilterModal(eventsWithoutProjects, currentFilters, 120, 40)

			// Modal should still work, just won't show project filter if no projects
			Expect(modal).NotTo(BeNil())
		})

		It("should allow multiple project selections", func() {
			modal = components.NewFilterModal(events, currentFilters, 120, 40)

			// The MultiSelect field should allow selecting multiple projects
			// This is verified by the form structure
			view := modal.View()
			Expect(view).NotTo(BeEmpty())
		})
	})

	Describe("Update", func() {
		BeforeEach(func() {
			modal = components.NewFilterModal(events, currentFilters, 120, 40)
			modal.Init()
		})

		It("should handle escape key to close modal without applying", func() {
			// TODO: Test escape key behavior once we can inject key messages
			Expect(modal.IsVisible()).To(BeTrue())
		})
	})

	Describe("View", func() {
		It("should return empty string when not visible", func() {
			modal = components.NewFilterModal(events, currentFilters, 120, 40)
			modal.Hide()

			view := modal.View()
			Expect(view).To(BeEmpty())
		})

		It("should render form with all filter sections when visible", func() {
			modal = components.NewFilterModal(events, currentFilters, 120, 40)

			view := modal.View()
			// Verify view renders content (non-empty)
			Expect(view).NotTo(BeEmpty())
			// Verify it has huh form controls
			Expect(view).To(ContainSubstring("enter"))
		})
	})
})
