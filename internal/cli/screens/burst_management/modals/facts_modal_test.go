package modals_test

import (
	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"

	"github.com/baphled/kariya/internal/cli/screens/burst_management/modals"
	"github.com/baphled/kariya/internal/cli/themes"
	"github.com/baphled/kariya/internal/domain/career"
	tea "github.com/charmbracelet/bubbletea"
)

var _ = Describe("BurstFactsModal", func() {
	var (
		modal     *modals.BurstFactsModal
		burstID   string
		burstName string
		facts     []*career.Fact
		theme     themes.Theme
	)

	BeforeEach(func() {
		burstID = "test-burst-id"
		burstName = "Backend Development"
		facts = []*career.Fact{
			{
				ID:                   "f1",
				Text:                 "Designed scalable microservices architecture",
				CompetencyCategories: []string{"architecture", "system design"},
				StrengthSignal:       "strong",
				SourceBurstID:        burstID,
			},
			{
				ID:                   "f2",
				Text:                 "Improved API response time by 40%",
				CompetencyCategories: []string{"performance"},
				StrengthSignal:       "strong",
				SourceBurstID:        burstID,
			},
			{
				ID:                   "f3",
				Text:                 "Mentored junior developers on best practices",
				CompetencyCategories: []string{"mentoring"},
				StrengthSignal:       "moderate",
				SourceBurstID:        burstID,
			},
		}
		theme = themes.NewDefaultTheme()
	})

	Describe("NewBurstFactsModal", func() {
		It("creates modal with facts content", func() {
			modal = modals.NewBurstFactsModal(burstID, burstName, facts, theme)

			Expect(modal).NotTo(BeNil())
			Expect(modal.GetBurstID()).To(Equal(burstID))
		})

		It("handles nil theme with default", func() {
			modal = modals.NewBurstFactsModal(burstID, burstName, facts, nil)

			Expect(modal).NotTo(BeNil())
		})

		It("handles empty facts list", func() {
			modal = modals.NewBurstFactsModal(burstID, burstName, []*career.Fact{}, theme)

			Expect(modal).NotTo(BeNil())
			modal.Show()
			view := modal.View()
			Expect(view).To(ContainSubstring("No facts"))
		})

		It("handles nil facts list", func() {
			modal = modals.NewBurstFactsModal(burstID, burstName, nil, theme)

			Expect(modal).NotTo(BeNil())
		})

		It("includes fact count in title", func() {
			modal = modals.NewBurstFactsModal(burstID, burstName, facts, theme)
			modal.Show()

			view := modal.View()
			Expect(view).To(ContainSubstring("3"))
		})

		It("includes burst name in title", func() {
			modal = modals.NewBurstFactsModal(burstID, burstName, facts, theme)
			modal.Show()

			view := modal.View()
			Expect(view).To(ContainSubstring("Backend Development"))
		})
	})

	Describe("Init", func() {
		BeforeEach(func() {
			modal = modals.NewBurstFactsModal(burstID, burstName, facts, theme)
		})

		It("returns initialization command", func() {
			cmd := modal.Init()
			_ = cmd
		})
	})

	Describe("Update", func() {
		BeforeEach(func() {
			modal = modals.NewBurstFactsModal(burstID, burstName, facts, theme)
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
			modal = modals.NewBurstFactsModal(burstID, burstName, facts, theme)
			modal.SetDimensions(80, 24)
		})

		It("renders facts when visible", func() {
			modal.Show()

			view := modal.View()

			Expect(view).NotTo(BeEmpty())
		})

		It("shows fact text", func() {
			modal.Show()

			view := modal.View()

			Expect(view).To(ContainSubstring("Designed scalable microservices"))
		})

		It("shows competency categories", func() {
			modal.Show()

			view := modal.View()

			Expect(view).To(ContainSubstring("architecture"))
		})

		It("shows strength signal", func() {
			modal.Show()

			view := modal.View()

			Expect(view).To(ContainSubstring("strong"))
		})

		It("numbers facts", func() {
			modal.Show()

			view := modal.View()

			// Modal may paginate, so we only check for first items visible in viewport.
			Expect(view).To(ContainSubstring("1."))
			Expect(view).To(ContainSubstring("2."))
		})

		It("shows empty state message when no facts", func() {
			modal = modals.NewBurstFactsModal(burstID, burstName, []*career.Fact{}, theme)
			modal.Show()

			view := modal.View()

			Expect(view).To(ContainSubstring("No facts extracted"))
		})
	})

	Describe("Visibility methods", func() {
		BeforeEach(func() {
			modal = modals.NewBurstFactsModal(burstID, burstName, facts, theme)
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
			modal = modals.NewBurstFactsModal(burstID, burstName, facts, theme)

			modal.SetDimensions(100, 50)

			Expect(modal).NotTo(BeNil())
		})
	})

	Describe("SetFacts", func() {
		It("updates the displayed facts", func() {
			modal = modals.NewBurstFactsModal(burstID, burstName, facts, theme)
			newFacts := []*career.Fact{
				{ID: "new-1", Text: "New fact one", SourceBurstID: burstID},
				{ID: "new-2", Text: "New fact two", SourceBurstID: burstID},
			}

			modal.SetFacts(newFacts)
			modal.Show()
			view := modal.View()

			Expect(view).To(ContainSubstring("New fact one"))
		})

		It("updates fact count in view", func() {
			modal = modals.NewBurstFactsModal(burstID, burstName, facts, theme)
			modal.Show()

			// Initially 3 facts.
			view := modal.View()
			Expect(view).To(ContainSubstring("3"))

			// Update to 1 fact.
			modal.SetFacts([]*career.Fact{
				{ID: "single", Text: "Single fact", SourceBurstID: burstID},
			})
			view = modal.View()
			Expect(view).To(ContainSubstring("1"))
		})
	})

	Describe("GetBurstID", func() {
		It("returns the burst ID", func() {
			modal = modals.NewBurstFactsModal(burstID, burstName, facts, theme)

			result := modal.GetBurstID()

			Expect(result).To(Equal(burstID))
		})
	})

	Describe("Facts with optional fields", func() {
		It("renders fact without categories", func() {
			factWithoutCategories := []*career.Fact{
				{
					ID:             "f-no-cat",
					Text:           "Fact without categories",
					StrengthSignal: "moderate",
					SourceBurstID:  burstID,
				},
			}
			modal = modals.NewBurstFactsModal(burstID, burstName, factWithoutCategories, theme)
			modal.Show()

			view := modal.View()

			Expect(view).To(ContainSubstring("Fact without categories"))
		})

		It("renders fact without strength signal", func() {
			factWithoutStrength := []*career.Fact{
				{
					ID:                   "f-no-str",
					Text:                 "Fact without strength",
					CompetencyCategories: []string{"delivery"},
					SourceBurstID:        burstID,
				},
			}
			modal = modals.NewBurstFactsModal(burstID, burstName, factWithoutStrength, theme)
			modal.Show()

			view := modal.View()

			Expect(view).To(ContainSubstring("Fact without strength"))
		})
	})
})
