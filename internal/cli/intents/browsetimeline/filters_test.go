package browsetimeline_test

import (
	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"

	"github.com/baphled/kariya/internal/cli/intents/browsetimeline"
	"github.com/baphled/kariya/internal/domain/career"
	"github.com/baphled/kariya/internal/testutil/fixtures"
)

var _ = Describe("Filtering Business Logic", func() {
	var (
		events []*career.Event
	)

	BeforeEach(func() {
		// Create test events with different companies, categories, and projects
		events = []*career.Event{
			fixtures.EventWith("event-1", "Platform work at TechCorp", "TechCorp", "Platform"),
			fixtures.EventWith("event-2", "Infrastructure work at CloudInc", "CloudInc", "Infrastructure"),
			fixtures.EventWith("event-3", "Another platform project at TechCorp", "TechCorp", "Platform"),
			fixtures.EventWith("event-4", "DevOps work at Acme Corp", "Acme Corp", "Infrastructure"),
		}

		// Set categories for events
		events[0].Categories = []string{"Backend"}
		events[1].Categories = []string{"DevOps"}
		events[2].Categories = []string{"Backend"}
		events[3].Categories = []string{"DevOps"}
	})

	Describe("Company Filtering", func() {
		It("should filter events by single company in view", func() {
			ctx := &browsetimeline.IntentContext{
				Events: events,
				InitialFilters: &browsetimeline.Filters{
					Companies: []string{"TechCorp"},
				},
			}
			intent, err := browsetimeline.NewIntent(ctx)
			Expect(err).NotTo(HaveOccurred())
			intent.Init()

			view := intent.View()
			Expect(view).To(ContainSubstring("TechCorp"))
			// Verify TechCorp events are shown
			Expect(view).To(ContainSubstring("Platform work"))
			// Verify other companies are hidden
			Expect(view).NotTo(ContainSubstring("CloudInc"))
			Expect(view).NotTo(ContainSubstring("Acme Corp"))
		})

		It("should filter events by multiple companies", func() {
			ctx := &browsetimeline.IntentContext{
				Events: events,
				InitialFilters: &browsetimeline.Filters{
					Companies: []string{"TechCorp", "Acme Corp"},
				},
			}
			intent, err := browsetimeline.NewIntent(ctx)
			Expect(err).NotTo(HaveOccurred())
			intent.Init()

			view := intent.View()
			Expect(view).To(ContainSubstring("TechCorp"))
			Expect(view).To(ContainSubstring("Acme Corp"))
			Expect(view).NotTo(ContainSubstring("CloudInc"))
		})

		It("should show all events when no filter applied", func() {
			ctx := &browsetimeline.IntentContext{
				Events: events,
			}
			intent, err := browsetimeline.NewIntent(ctx)
			Expect(err).NotTo(HaveOccurred())
			intent.Init()

			view := intent.View()
			Expect(view).To(ContainSubstring("TechCorp"))
			Expect(view).To(ContainSubstring("CloudInc"))
			Expect(view).To(ContainSubstring("Acme Corp"))
		})
	})

	Describe("Combined Company and Category Filtering", func() {
		It("should apply AND logic: company AND category both must match", func() {
			ctx := &browsetimeline.IntentContext{
				Events: events,
				InitialFilters: &browsetimeline.Filters{
					Companies:  []string{"TechCorp"},
					Categories: []string{"Backend"},
				},
			}
			intent, err := browsetimeline.NewIntent(ctx)
			Expect(err).NotTo(HaveOccurred())
			intent.Init()

			view := intent.View()
			// Should show TechCorp with Backend
			Expect(view).To(ContainSubstring("Platform work"))
			Expect(view).To(ContainSubstring("Another platform project"))
			// Should not show TechCorp with DevOps (exists but filtered out by category)
			Expect(view).NotTo(ContainSubstring("CloudInc"))
			Expect(view).NotTo(ContainSubstring("Acme Corp"))
		})
	})

	Describe("Filter State Management", func() {
		It("should report active filters when filters are set", func() {
			ctx := &browsetimeline.IntentContext{
				Events: events,
				InitialFilters: &browsetimeline.Filters{
					Companies: []string{"TechCorp"},
				},
			}
			intent, err := browsetimeline.NewIntent(ctx)
			Expect(err).NotTo(HaveOccurred())

			Expect(intent.HasActiveFilters()).To(BeTrue())
		})

		It("should report no active filters when no filters are set", func() {
			ctx := &browsetimeline.IntentContext{
				Events: events,
			}
			intent, err := browsetimeline.NewIntent(ctx)
			Expect(err).NotTo(HaveOccurred())

			Expect(intent.HasActiveFilters()).To(BeFalse())
		})
	})

	Describe("Edge Cases", func() {
		It("should handle filtering empty event list", func() {
			ctx := &browsetimeline.IntentContext{
				Events: []*career.Event{},
				InitialFilters: &browsetimeline.Filters{
					Companies: []string{"AnyCompany"},
				},
			}
			intent, err := browsetimeline.NewIntent(ctx)
			Expect(err).NotTo(HaveOccurred())
			intent.Init()

			// Should not crash, should show empty state
			view := intent.View()
			Expect(view).NotTo(BeEmpty())
		})

		It("should handle filter matching no events", func() {
			ctx := &browsetimeline.IntentContext{
				Events: events,
				InitialFilters: &browsetimeline.Filters{
					Companies: []string{"NonExistent"},
				},
			}
			intent, err := browsetimeline.NewIntent(ctx)
			Expect(err).NotTo(HaveOccurred())
			intent.Init()

			view := intent.View()
			// Should show empty state, not crash
			Expect(view).NotTo(ContainSubstring("TechCorp"))
			Expect(view).NotTo(ContainSubstring("CloudInc"))
			Expect(view).NotTo(ContainSubstring("Acme Corp"))
		})
	})
})
