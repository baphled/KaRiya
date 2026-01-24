package burst_management_test

import (
	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"

	"github.com/baphled/kariya/internal/cli/intents/burst_management"
	"github.com/baphled/kariya/internal/domain/career"
	tea "github.com/charmbracelet/bubbletea"
)

var _ = Describe("Intent Methods", func() {
	Describe("Init", func() {
		It("should initialize the intent successfully", func() {
			ctx := &burst_management.IntentContext{
				Bursts: []*career.Burst{
					{ID: "burst-1", Name: "Burst 1"},
				},
			}
			ctx.Validate()

			intent, err := burst_management.NewIntent(ctx)
			Expect(err).NotTo(HaveOccurred())

			cmd := intent.Init()
			Expect(cmd).To(BeNil())
			Expect(intent.GetState()).To(Equal(burst_management.StateList))
			Expect(intent.GetFilteredBursts()).To(HaveLen(1))
		})

		It("should set selected burst if bursts available", func() {
			bursts := []*career.Burst{
				{ID: "burst-1", Name: "First Burst"},
				{ID: "burst-2", Name: "Second Burst"},
			}
			ctx := &burst_management.IntentContext{
				Bursts: bursts,
			}
			ctx.Validate()

			intent, _ := burst_management.NewIntent(ctx)
			intent.Init()

			Expect(intent.GetSelectedBurst()).NotTo(BeNil())
			Expect(intent.GetSelectedBurst().ID).To(Equal("burst-1"))
		})

		It("should handle empty burst list", func() {
			ctx := &burst_management.IntentContext{
				Bursts: []*career.Burst{},
			}
			ctx.Validate()

			intent, _ := burst_management.NewIntent(ctx)
			intent.Init()

			Expect(intent.GetFilteredBursts()).To(BeEmpty())
			Expect(intent.GetSelectedBurst()).To(BeNil())
		})
	})

	Describe("Update", func() {
		var intent *burst_management.Intent

		BeforeEach(func() {
			ctx := &burst_management.IntentContext{
				Bursts: []*career.Burst{
					{ID: "burst-1", Name: "Burst 1"},
				},
			}
			ctx.Validate()
			intent, _ = burst_management.NewIntent(ctx)
			intent.Init()
		})

		It("should return nil when intent is not active", func() {
			intent.Deactivate()
			cmd := intent.Update(tea.KeyMsg{})
			Expect(cmd).To(BeNil())
		})

		It("should handle quit key in list view", func() {
			msg := tea.KeyMsg{Type: tea.KeyRunes, Runes: []rune{'q'}}
			cmd := intent.Update(msg)
			Expect(cmd).NotTo(BeNil())
		})

		It("should handle escape key to cancel in list view", func() {
			msg := tea.KeyMsg{Type: tea.KeyEsc}
			cmd := intent.Update(msg)
			Expect(cmd).To(BeNil())
			Expect(intent.IsActive()).To(BeFalse())
			Expect(intent.Result()).NotTo(BeNil())
		})

		It("should delegate to appropriate state handler", func() {
			intent.SetState(burst_management.StateDetail)
			cmd := intent.Update(tea.KeyMsg{})
			Expect(cmd).To(BeNil())
		})
	})

	Describe("View", func() {
		var intent *burst_management.Intent

		BeforeEach(func() {
			ctx := &burst_management.IntentContext{
				Bursts: []*career.Burst{
					{ID: "burst-1", Name: "Burst 1"},
				},
			}
			ctx.Validate()
			intent, _ = burst_management.NewIntent(ctx)
			intent.Init()
		})

		It("should render when active", func() {
			view := intent.View()
			Expect(view).NotTo(BeEmpty())
		})

		It("should show inactive message when not active", func() {
			intent.Deactivate()
			view := intent.View()
			Expect(view).To(ContainSubstring("not active"))
		})

		It("should render list view in StateList", func() {
			intent.SetState(burst_management.StateList)
			view := intent.View()
			Expect(view).To(ContainSubstring("Burst list view"))
		})

		It("should render detail view in StateDetail", func() {
			intent.SetState(burst_management.StateDetail)
			intent.SetSelectedBurst(&career.Burst{ID: "burst-1", Name: "Test"})
			view := intent.View()
			Expect(view).To(ContainSubstring("Burst detail view"))
		})

		It("should show no burst selected when nil in detail", func() {
			intent.SetState(burst_management.StateDetail)
			intent.SetSelectedBurst(nil)
			view := intent.View()
			Expect(view).To(ContainSubstring("No burst selected"))
		})

		It("should show no bursts found when empty in list", func() {
			emptyCtx := &burst_management.IntentContext{
				Bursts: []*career.Burst{},
			}
			emptyCtx.Validate()
			emptyIntent, _ := burst_management.NewIntent(emptyCtx)
			emptyIntent.Init()

			view := emptyIntent.View()
			Expect(view).To(ContainSubstring("No bursts found"))
		})
	})
})
