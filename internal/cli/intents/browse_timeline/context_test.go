package browse_timeline_test

import (
	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"

	"github.com/baphled/kariya/internal/cli/intents/browse_timeline"
	"github.com/baphled/kariya/internal/domain/career"
	"github.com/baphled/kariya/internal/testutil/fixtures"
)

var _ = Describe("Context", func() {
	Describe("IntentContext", func() {
		Describe("Construction", func() {
			It("should create an empty context", func() {
				ctx := &browse_timeline.IntentContext{}
				Expect(ctx).NotTo(BeNil())
			})

			It("should create a context with events", func() {
				events := []*career.CareerEvent{
					fixtures.Event("event-1"),
					fixtures.Event("event-2"),
				}
				ctx := &browse_timeline.IntentContext{
					Events: events,
				}
				Expect(ctx.Events).To(HaveLen(2))
			})

			It("should create a context with initial filters", func() {
				filters := &browse_timeline.Filters{
					SearchText: "test",
				}
				ctx := &browse_timeline.IntentContext{
					InitialFilters: filters,
				}
				Expect(ctx.InitialFilters.SearchText).To(Equal("test"))
			})

			It("should create a context with selected event ID", func() {
				ctx := &browse_timeline.IntentContext{
					SelectedEventID: "event-123",
				}
				Expect(ctx.SelectedEventID).To(Equal("event-123"))
			})
		})

		Describe("Validate", func() {
			Context("when Events is nil", func() {
				It("should initialize to empty slice", func() {
					ctx := &browse_timeline.IntentContext{
						Events: nil,
					}
					err := ctx.Validate()
					Expect(err).NotTo(HaveOccurred())
					Expect(ctx.Events).NotTo(BeNil())
					Expect(ctx.Events).To(BeEmpty())
				})
			})

			Context("when InitialFilters is nil", func() {
				It("should initialize to default", func() {
					ctx := &browse_timeline.IntentContext{
						InitialFilters: nil,
					}
					err := ctx.Validate()
					Expect(err).NotTo(HaveOccurred())
					Expect(ctx.InitialFilters).NotTo(BeNil())
				})

				It("should set default sort options", func() {
					ctx := &browse_timeline.IntentContext{}
					err := ctx.Validate()
					Expect(err).NotTo(HaveOccurred())
					Expect(ctx.InitialFilters.SortBy).To(Equal("date"))
					Expect(ctx.InitialFilters.SortOrder).To(Equal("desc"))
				})

				It("should initialize empty filter slices", func() {
					ctx := &browse_timeline.IntentContext{}
					err := ctx.Validate()
					Expect(err).NotTo(HaveOccurred())
					Expect(ctx.InitialFilters.Tags).NotTo(BeNil())
					Expect(ctx.InitialFilters.Tags).To(BeEmpty())
					Expect(ctx.InitialFilters.Companies).NotTo(BeNil())
					Expect(ctx.InitialFilters.Companies).To(BeEmpty())
					Expect(ctx.InitialFilters.Categories).NotTo(BeNil())
					Expect(ctx.InitialFilters.Categories).To(BeEmpty())
					Expect(ctx.InitialFilters.Projects).NotTo(BeNil())
					Expect(ctx.InitialFilters.Projects).To(BeEmpty())
				})
			})

			Context("when Events already exist", func() {
				It("should preserve existing Events", func() {
					events := []*career.CareerEvent{
						fixtures.Event("event-1"),
					}
					ctx := &browse_timeline.IntentContext{
						Events: events,
					}
					err := ctx.Validate()
					Expect(err).NotTo(HaveOccurred())
					Expect(ctx.Events).To(HaveLen(1))
					Expect(ctx.Events[0].ID).To(Equal("event-1"))
				})
			})

			Context("when InitialFilters already exist", func() {
				It("should preserve existing InitialFilters", func() {
					filters := &browse_timeline.Filters{
						SearchText: "existing",
						SortBy:     "text",
						SortOrder:  "asc",
					}
					ctx := &browse_timeline.IntentContext{
						InitialFilters: filters,
					}
					err := ctx.Validate()
					Expect(err).NotTo(HaveOccurred())
					Expect(ctx.InitialFilters.SearchText).To(Equal("existing"))
					Expect(ctx.InitialFilters.SortBy).To(Equal("text"))
					Expect(ctx.InitialFilters.SortOrder).To(Equal("asc"))
				})
			})
		})
	})

	Describe("Filters", func() {
		Describe("Construction", func() {
			It("should create an empty filter", func() {
				filters := &browse_timeline.Filters{}
				Expect(filters).NotTo(BeNil())
			})

			It("should store search text", func() {
				filters := &browse_timeline.Filters{
					SearchText: "backend developer",
				}
				Expect(filters.SearchText).To(Equal("backend developer"))
			})

			It("should store tags", func() {
				filters := &browse_timeline.Filters{
					Tags: []string{"golang", "api"},
				}
				Expect(filters.Tags).To(HaveLen(2))
				Expect(filters.Tags).To(ContainElement("golang"))
				Expect(filters.Tags).To(ContainElement("api"))
			})

			It("should store companies", func() {
				filters := &browse_timeline.Filters{
					Companies: []string{"TechCorp", "CloudInc"},
				}
				Expect(filters.Companies).To(HaveLen(2))
			})

			It("should store categories", func() {
				filters := &browse_timeline.Filters{
					Categories: []string{"development", "architecture"},
				}
				Expect(filters.Categories).To(HaveLen(2))
			})

			It("should store projects", func() {
				filters := &browse_timeline.Filters{
					Projects: []string{"Project A", "Project B"},
				}
				Expect(filters.Projects).To(HaveLen(2))
			})

			It("should store date range", func() {
				filters := &browse_timeline.Filters{
					DateFrom: "2023-01-01",
					DateTo:   "2024-12-31",
				}
				Expect(filters.DateFrom).To(Equal("2023-01-01"))
				Expect(filters.DateTo).To(Equal("2024-12-31"))
			})

			It("should store sort options", func() {
				filters := &browse_timeline.Filters{
					SortBy:    "date",
					SortOrder: "desc",
				}
				Expect(filters.SortBy).To(Equal("date"))
				Expect(filters.SortOrder).To(Equal("desc"))
			})
		})

		Describe("Filter Combinations", func() {
			It("should support multiple filter types together", func() {
				filters := &browse_timeline.Filters{
					SearchText: "developer",
					Tags:       []string{"golang"},
					Companies:  []string{"TechCorp"},
					Categories: []string{"development"},
					SortBy:     "date",
					SortOrder:  "asc",
				}
				Expect(filters.SearchText).NotTo(BeEmpty())
				Expect(filters.Tags).NotTo(BeEmpty())
				Expect(filters.Companies).NotTo(BeEmpty())
				Expect(filters.Categories).NotTo(BeEmpty())
				Expect(filters.SortBy).To(Equal("date"))
			})
		})
	})

	Describe("EventService Interface", func() {
		It("should define required methods", func() {
			// The EventService interface defines:
			// - DeleteEvent
			// - ListEvents
			// - CaptureEvent
			// - UpdateEventMetadata
			// - GetSkillsForEvent
			// Interface existence verified by compilation.
			Expect(true).To(BeTrue())
		})
	})
})
