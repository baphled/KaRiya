package browsetimeline_test

import (
	"strings"
	"time"

	"github.com/baphled/kariya/internal/ui/terminal"
	"github.com/baphled/kariya/internal/ui/themes"
	tea "github.com/charmbracelet/bubbletea"
	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"

	"github.com/baphled/kariya/internal/domain/career"
	"github.com/baphled/kariya/internal/testutil/fixtures"
	"github.com/baphled/kariya/internal/tui/intents/browsetimeline"
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
			ctx := &browsetimeline.IntentValidator{
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
			ctx := &browsetimeline.IntentValidator{
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
			ctx := &browsetimeline.IntentValidator{
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
			ctx := &browsetimeline.IntentValidator{
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
			ctx := &browsetimeline.IntentValidator{
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
			ctx := &browsetimeline.IntentValidator{
				Events: events,
			}
			intent, err := browsetimeline.NewIntent(ctx)
			Expect(err).NotTo(HaveOccurred())

			Expect(intent.HasActiveFilters()).To(BeFalse())
		})

		It("should report active filters for tag, project, date and sort changes", func() {
			ctx := &browsetimeline.IntentValidator{
				Events: events,
				InitialFilters: &browsetimeline.Filters{
					Tags:      []string{"golang"},
					Projects:  []string{"Platform"},
					DateFrom:  "2020-01-01",
					DateTo:    "2030-01-01",
					SortBy:    "text",
					SortOrder: "asc",
				},
			}
			intent, err := browsetimeline.NewIntent(ctx)
			Expect(err).NotTo(HaveOccurred())

			Expect(intent.HasActiveFilters()).To(BeTrue())
		})
	})

	Describe("Edge Cases", func() {
		It("should handle filtering empty event list", func() {
			ctx := &browsetimeline.IntentValidator{
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
			ctx := &browsetimeline.IntentValidator{
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

		It("should filter by tags through the public intent API", func() {
			events[0].Tags = []string{"golang", "backend"}
			events[1].Tags = []string{"terraform"}
			events[2].Tags = []string{"golang"}
			events[3].Tags = []string{"ops"}

			intent, err := browsetimeline.NewIntent(&browsetimeline.IntentValidator{
				Events: events,
				InitialFilters: &browsetimeline.Filters{
					Tags: []string{"golang"},
				},
			})
			Expect(err).NotTo(HaveOccurred())
			intent.Init()

			view := intent.View()
			Expect(view).To(ContainSubstring("Platform work"))
			Expect(view).To(ContainSubstring("Another platform project"))
			Expect(view).NotTo(ContainSubstring("Infrastructure work at CloudInc"))
		})

		It("should filter by project through the public intent API", func() {
			intent, err := browsetimeline.NewIntent(&browsetimeline.IntentValidator{
				Events: events,
				InitialFilters: &browsetimeline.Filters{
					Projects: []string{"Infrastructure"},
				},
			})
			Expect(err).NotTo(HaveOccurred())
			intent.Init()

			view := intent.View()
			Expect(view).To(ContainSubstring("Infrastructure work at CloudInc"))
			Expect(view).To(ContainSubstring("DevOps work at Acme Corp"))
			Expect(view).NotTo(ContainSubstring("Platform work at TechCorp"))
		})

		It("should filter by date range through the public intent API", func() {
			events[0].Date = time.Date(2020, time.January, 1, 0, 0, 0, 0, time.UTC)
			events[1].Date = time.Date(2021, time.January, 1, 0, 0, 0, 0, time.UTC)
			events[2].Date = time.Date(2022, time.January, 1, 0, 0, 0, 0, time.UTC)
			events[3].Date = time.Date(2023, time.January, 1, 0, 0, 0, 0, time.UTC)

			intent, err := browsetimeline.NewIntent(&browsetimeline.IntentValidator{
				Events: events,
				InitialFilters: &browsetimeline.Filters{
					DateFrom: "2021-01-01",
					DateTo:   "2022-12-31",
				},
			})
			Expect(err).NotTo(HaveOccurred())
			intent.Init()

			view := intent.View()
			Expect(view).NotTo(ContainSubstring("Platform work at TechCorp"))
			Expect(view).To(ContainSubstring("Infrastructure work at CloudInc"))
			Expect(view).To(ContainSubstring("Another platform project at TechCorp"))
			Expect(view).NotTo(ContainSubstring("DevOps work at Acme Corp"))
		})

		It("should sort by text ascending through initial filters", func() {
			intent, err := browsetimeline.NewIntent(&browsetimeline.IntentValidator{
				Events: events,
				InitialFilters: &browsetimeline.Filters{
					SortBy:    "text",
					SortOrder: "asc",
				},
			})
			Expect(err).NotTo(HaveOccurred())
			intent.Init()

			view := intent.View()
			first := strings.Index(view, "Another platform project")
			second := strings.Index(view, "DevOps work")
			third := strings.Index(view, "Infrastructure work")
			fourth := strings.Index(view, "Platform work")
			Expect(first).To(BeNumerically("<", second))
			Expect(second).To(BeNumerically("<", third))
			Expect(third).To(BeNumerically("<", fourth))
		})

		It("should sort by date ascending through initial filters", func() {
			events[0].Date = time.Date(2024, time.April, 1, 0, 0, 0, 0, time.UTC)
			events[1].Date = time.Date(2024, time.March, 1, 0, 0, 0, 0, time.UTC)
			events[2].Date = time.Date(2024, time.February, 1, 0, 0, 0, 0, time.UTC)
			events[3].Date = time.Date(2024, time.January, 1, 0, 0, 0, 0, time.UTC)

			intent, err := browsetimeline.NewIntent(&browsetimeline.IntentValidator{
				Events: events,
				InitialFilters: &browsetimeline.Filters{
					SortBy:    "date",
					SortOrder: "asc",
				},
			})
			Expect(err).NotTo(HaveOccurred())
			intent.Init()

			view := intent.View()
			first := strings.Index(view, "DevOps work at Acme Corp")
			second := strings.Index(view, "Another platform project at TechCorp")
			third := strings.Index(view, "Infrastructure work at CloudInc")
			fourth := strings.Index(view, "Platform work at TechCorp")
			Expect(first).To(BeNumerically("<", second))
			Expect(second).To(BeNumerically("<", third))
			Expect(third).To(BeNumerically("<", fourth))
		})

		It("should clear filters from a no-stack state and restore the full list", func() {
			intent, err := browsetimeline.NewIntent(&browsetimeline.IntentValidator{
				Events: events,
			})
			Expect(err).NotTo(HaveOccurred())
			intent.Init()

			intent.ClearFilters()
			intent.ApplyFilters()
			view := intent.View()
			Expect(intent.HasActiveFilters()).To(BeFalse())
			Expect(view).To(ContainSubstring("CloudInc"))
			Expect(view).To(ContainSubstring("Acme Corp"))
		})

		It("should keep stacked initial filters active when only the latest layer is cleared", func() {
			intent, err := browsetimeline.NewIntent(&browsetimeline.IntentValidator{
				Events: events,
				InitialFilters: &browsetimeline.Filters{
					Companies: []string{"TechCorp"},
				},
			})
			Expect(err).NotTo(HaveOccurred())
			intent.Init()

			Expect(intent.View()).NotTo(ContainSubstring("CloudInc"))
			intent.ClearFilters()
			intent.ApplyFilters()
			Expect(intent.View()).NotTo(ContainSubstring("CloudInc"))
		})

		It("should render with custom terminal info and theme through public APIs", func() {
			intent, err := browsetimeline.NewIntent(&browsetimeline.IntentValidator{Events: events})
			Expect(err).NotTo(HaveOccurred())
			intent.UpdateTerminalInfo(&terminal.Info{Width: 120, Height: 40})
			tm := themes.NewThemeManager()
			intent.SetThemeManager(tm)
			intent.Init()

			view := intent.View()
			Expect(intent.Result()).To(BeNil())
			Expect(view).To(ContainSubstring("Browse Timeline"))
		})

		It("should return nil result even after refresh data", func() {
			intent, err := browsetimeline.NewIntent(&browsetimeline.IntentValidator{Events: events})
			Expect(err).NotTo(HaveOccurred())
			intent.Init()

			Expect(intent.RefreshData()).To(BeNil())
			Expect(intent.Result()).To(BeNil())
			Expect(intent.View()).To(ContainSubstring("Timeline"))
		})

		It("should keep the list view active after cancelling a skills modal", func() {
			intent, err := browsetimeline.NewIntent(&browsetimeline.IntentValidator{Events: events})
			Expect(err).NotTo(HaveOccurred())
			intent.Init()
			Expect(intent.Update(tea.KeyMsg{Type: tea.KeyEnter})).To(BeNil())
			Expect(intent.Update(browsetimeline.SkillsForModalLoadedMsg{
				EventID: events[0].ID,
				Skills:  []*career.Skill{fixtures.SkillWith("skill-1", "Go", "", "")},
			})).To(BeNil())

			cmd := intent.Update(tea.KeyMsg{Type: tea.KeyEsc})
			Expect(cmd).NotTo(BeNil())
			Expect(intent.HasVisibleSkillsModal()).To(BeFalse())
			Expect(intent.Result()).To(BeNil())
		})
	})
})
