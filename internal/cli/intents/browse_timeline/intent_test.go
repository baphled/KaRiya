package browse_timeline_test

import (
	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"

	"github.com/baphled/kariya/internal/cli/behaviors"
	"github.com/baphled/kariya/internal/cli/intents"
	"github.com/baphled/kariya/internal/cli/intents/browse_timeline"
	"github.com/baphled/kariya/internal/domain/career"
	"github.com/baphled/kariya/internal/testutil/fixtures"
	tea "github.com/charmbracelet/bubbletea"
)

var _ = Describe("Intent", func() {
	var (
		intent *browse_timeline.Intent
		ctx    *browse_timeline.IntentContext
		events []*career.CareerEvent
	)

	BeforeEach(func() {
		events = []*career.CareerEvent{
			fixtures.EventWith("event-1", "Backend Developer at TechCorp", "TechCorp", "Platform"),
			fixtures.EventWith("event-2", "DevOps Engineer at CloudInc", "CloudInc", "Infrastructure"),
		}
		ctx = &browse_timeline.IntentContext{
			Events: events,
		}
	})

	Describe("Construction", func() {
		Context("with valid context", func() {
			It("should create an intent successfully", func() {
				intent, err := browse_timeline.NewIntent(ctx)
				Expect(err).NotTo(HaveOccurred())
				Expect(intent).NotTo(BeNil())
			})

			It("should embed BaseIntent", func() {
				intent, _ := browse_timeline.NewIntent(ctx)
				Expect(intent.BaseIntent).NotTo(BeNil())
			})
		})

		Context("with empty events", func() {
			It("should create an intent with empty event list", func() {
				emptyCtx := &browse_timeline.IntentContext{
					Events: []*career.CareerEvent{},
				}
				intent, err := browse_timeline.NewIntent(emptyCtx)
				Expect(err).NotTo(HaveOccurred())
				Expect(intent).NotTo(BeNil())
			})
		})

		Context("with nil events", func() {
			It("should handle nil events by initializing empty slice", func() {
				nilCtx := &browse_timeline.IntentContext{
					Events: nil,
				}
				intent, err := browse_timeline.NewIntent(nilCtx)
				Expect(err).NotTo(HaveOccurred())
				Expect(intent).NotTo(BeNil())
			})
		})
	})

	Describe("Initialization", func() {
		BeforeEach(func() {
			var err error
			intent, err = browse_timeline.NewIntent(ctx)
			Expect(err).NotTo(HaveOccurred())
		})

		It("should return nil command on init", func() {
			cmd := intent.Init()
			Expect(cmd).To(BeNil())
		})

		It("should render view after init", func() {
			intent.Init()
			view := intent.View()
			Expect(view).NotTo(BeEmpty())
		})

		It("should show timeline content after init", func() {
			intent.Init()
			view := intent.View()
			Expect(view).To(ContainSubstring("Timeline"))
		})

		It("should display events in view", func() {
			intent.Init()
			view := intent.View()
			Expect(view).To(ContainSubstring("TechCorp"))
			Expect(view).To(ContainSubstring("CloudInc"))
		})
	})

	Describe("View Rendering", func() {
		BeforeEach(func() {
			var err error
			intent, err = browse_timeline.NewIntent(ctx)
			Expect(err).NotTo(HaveOccurred())
			intent.Init()
		})

		Context("with events", func() {
			It("should show event count", func() {
				view := intent.View()
				Expect(view).To(ContainSubstring("Events: 2"))
			})

			It("should show navigation hints", func() {
				view := intent.View()
				Expect(view).To(SatisfyAny(
					ContainSubstring("j/k"),
					ContainSubstring("Enter"),
				))
			})
		})

		Context("with empty events", func() {
			BeforeEach(func() {
				emptyCtx := &browse_timeline.IntentContext{
					Events: []*career.CareerEvent{},
				}
				var err error
				intent, err = browse_timeline.NewIntent(emptyCtx)
				Expect(err).NotTo(HaveOccurred())
				intent.Init()
			})

			It("should show empty state message", func() {
				view := intent.View()
				Expect(view).To(ContainSubstring("No events"))
			})

			It("should show add hint", func() {
				view := intent.View()
				Expect(view).To(ContainSubstring("Add"))
			})
		})
	})

	Describe("Navigation", func() {
		BeforeEach(func() {
			var err error
			intent, err = browse_timeline.NewIntent(ctx)
			Expect(err).NotTo(HaveOccurred())
			intent.Init()
		})

		Context("keyboard navigation", func() {
			It("should handle arrow down", func() {
				cmd := intent.Update(tea.KeyMsg{Type: tea.KeyDown})
				Expect(cmd).To(BeNil())
				view := intent.View()
				Expect(view).NotTo(BeEmpty())
			})

			It("should handle arrow up", func() {
				intent.Update(tea.KeyMsg{Type: tea.KeyDown})
				cmd := intent.Update(tea.KeyMsg{Type: tea.KeyUp})
				Expect(cmd).To(BeNil())
			})

			It("should handle vim j key", func() {
				cmd := intent.Update(tea.KeyMsg{Type: tea.KeyRunes, Runes: []rune{'j'}})
				Expect(cmd).To(BeNil())
			})

			It("should handle vim k key", func() {
				intent.Update(tea.KeyMsg{Type: tea.KeyRunes, Runes: []rune{'j'}})
				cmd := intent.Update(tea.KeyMsg{Type: tea.KeyRunes, Runes: []rune{'k'}})
				Expect(cmd).To(BeNil())
			})
		})

		Context("boundary behavior", func() {
			It("should not crash when navigating up at first item", func() {
				intent.Update(tea.KeyMsg{Type: tea.KeyUp})
				view := intent.View()
				Expect(view).NotTo(BeEmpty())
			})

			It("should not crash when navigating down past last item", func() {
				for i := 0; i < 10; i++ {
					intent.Update(tea.KeyMsg{Type: tea.KeyDown})
				}
				view := intent.View()
				Expect(view).NotTo(BeEmpty())
			})
		})
	})

	Describe("Cancellation", func() {
		BeforeEach(func() {
			var err error
			intent, err = browse_timeline.NewIntent(ctx)
			Expect(err).NotTo(HaveOccurred())
			intent.Init()
		})

		Context("from list view", func() {
			It("should cancel intent on escape", func() {
				intent.Update(tea.KeyMsg{Type: tea.KeyEsc})
				result := intent.Result()
				Expect(result).NotTo(BeNil())
				Expect(result.Status).To(Equal(intents.Cancelled))
			})
		})

		Context("from detail modal", func() {
			It("should close modal on escape without cancelling intent", func() {
				intent.Update(tea.KeyMsg{Type: tea.KeyEnter})
				intent.Update(tea.KeyMsg{Type: tea.KeyEsc})
				result := intent.Result()
				Expect(result).To(BeNil())
			})
		})
	})

	Describe("Event Selection", func() {
		BeforeEach(func() {
			var err error
			intent, err = browse_timeline.NewIntent(ctx)
			Expect(err).NotTo(HaveOccurred())
			intent.Init()
		})

		It("should show detail modal on enter", func() {
			intent.Update(tea.KeyMsg{Type: tea.KeyEnter})
			view := intent.View()
			// Should show one of the events in detail view.
			Expect(view).To(SatisfyAny(
				ContainSubstring("Backend Developer"),
				ContainSubstring("DevOps Engineer"),
			))
		})

		It("should return to list after closing modal", func() {
			intent.Update(tea.KeyMsg{Type: tea.KeyEnter})
			intent.Update(tea.KeyMsg{Type: tea.KeyEsc})
			view := intent.View()
			Expect(view).To(ContainSubstring("Timeline"))
		})
	})

	Describe("Modal Interactions", func() {
		BeforeEach(func() {
			var err error
			intent, err = browse_timeline.NewIntent(ctx)
			Expect(err).NotTo(HaveOccurred())
			intent.Init()
		})

		Context("search modal", func() {
			It("should open search modal with / key", func() {
				intent.Update(tea.KeyMsg{Type: tea.KeyRunes, Runes: []rune{'/'}})
				view := intent.View()
				Expect(view).To(ContainSubstring("Search"))
			})
		})

		Context("filter modal", func() {
			It("should open filter modal with f key", func() {
				intent.Update(tea.KeyMsg{Type: tea.KeyRunes, Runes: []rune{'f'}})
				view := intent.View()
				Expect(view).To(ContainSubstring("Filter"))
			})
		})

		Context("sort modal", func() {
			It("should open sort modal with s key", func() {
				intent.Update(tea.KeyMsg{Type: tea.KeyRunes, Runes: []rune{'s'}})
				view := intent.View()
				Expect(view).To(ContainSubstring("Sort"))
			})
		})

		Context("quick add modal", func() {
			It("should open quick add modal with a key", func() {
				intent.Update(tea.KeyMsg{Type: tea.KeyRunes, Runes: []rune{'a'}})
				Expect(intent.HasVisibleQuickAddModal()).To(BeTrue())
			})
		})

		Context("edit modal", func() {
			It("should open edit modal with e key", func() {
				intent.Update(tea.KeyMsg{Type: tea.KeyRunes, Runes: []rune{'e'}})
				Expect(intent.HasVisibleEditModal()).To(BeTrue())
			})
		})

		Context("delete modal", func() {
			It("should open delete confirmation with d key", func() {
				intent.Update(tea.KeyMsg{Type: tea.KeyRunes, Runes: []rune{'d'}})
				view := intent.View()
				Expect(view).To(ContainSubstring("Delete"))
			})
		})

		Context("modal priority", func() {
			It("should not allow multiple modals", func() {
				intent.Update(tea.KeyMsg{Type: tea.KeyRunes, Runes: []rune{'/'}})
				intent.Update(tea.KeyMsg{Type: tea.KeyRunes, Runes: []rune{'f'}})
				view := intent.View()
				Expect(view).To(ContainSubstring("Search"))
				Expect(view).NotTo(ContainSubstring("Filter by Company"))
			})
		})
	})

	Describe("Window Resize", func() {
		BeforeEach(func() {
			var err error
			intent, err = browse_timeline.NewIntent(ctx)
			Expect(err).NotTo(HaveOccurred())
			intent.Init()
		})

		It("should handle window size message", func() {
			cmd := intent.Update(tea.WindowSizeMsg{Width: 80, Height: 24})
			Expect(cmd).To(BeNil())
			view := intent.View()
			Expect(view).NotTo(BeEmpty())
		})

		It("should adapt to different terminal sizes", func() {
			intent.Update(tea.WindowSizeMsg{Width: 120, Height: 40})
			view := intent.View()
			Expect(view).NotTo(BeEmpty())
		})
	})

	Describe("Error Handling", func() {
		BeforeEach(func() {
			var err error
			intent, err = browse_timeline.NewIntent(ctx)
			Expect(err).NotTo(HaveOccurred())
			intent.Init()
		})

		It("should show error modal when error occurs", func() {
			intent.ShowErrorModal("Test Error", "Something went wrong")
			Expect(intent.HasVisibleErrorModal()).To(BeTrue())
		})

		It("should dismiss error modal with escape", func() {
			intent.ShowErrorModal("Test Error", "Something went wrong")
			intent.Update(tea.KeyMsg{Type: tea.KeyEsc})
			Expect(intent.HasVisibleErrorModal()).To(BeFalse())
		})

		It("should render error content in view", func() {
			intent.ShowErrorModal("Operation Failed", "Database error")
			view := intent.View()
			Expect(view).To(ContainSubstring("Operation Failed"))
		})
	})

	Describe("Filter Behavior", func() {
		BeforeEach(func() {
			var err error
			intent, err = browse_timeline.NewIntent(ctx)
			Expect(err).NotTo(HaveOccurred())
			intent.Init()
		})

		It("should implement FilterBehavior interface", func() {
			var _ behaviors.FilterBehavior = intent
		})

		It("should report no active filters initially", func() {
			Expect(intent.HasActiveFilters()).To(BeFalse())
		})

		It("should clear filters", func() {
			intent.ClearFilters()
			Expect(intent.HasActiveFilters()).To(BeFalse())
		})
	})

	Describe("Inactive State", func() {
		BeforeEach(func() {
			var err error
			intent, err = browse_timeline.NewIntent(ctx)
			Expect(err).NotTo(HaveOccurred())
			intent.Init()
			intent.Update(tea.KeyMsg{Type: tea.KeyEsc})
		})

		It("should not process updates when inactive", func() {
			cmd := intent.Update(tea.KeyMsg{Type: tea.KeyEnter})
			Expect(cmd).To(BeNil())
		})
	})
})
