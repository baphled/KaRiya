package models

import (
	"context"
	"time"

	"github.com/baphled/kariya/internal/domain/career"
	careerrepo "github.com/baphled/kariya/internal/repository/career"
	careerservice "github.com/baphled/kariya/internal/service/career"
	tea "github.com/charmbracelet/bubbletea"
	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
)

var _ = Describe("FactsResultsModel", func() {
	var (
		model    *FactsResultsModel
		service  *careerservice.Service
		repo     *careerrepo.MemoryRepository
		ctx      context.Context
		testFact *career.Fact
	)

	BeforeEach(func() {
		ctx = context.Background()
		repo = careerrepo.NewMemoryRepository()
		service = careerservice.NewService(repo)

		// Create test fact
		testFact = &career.Fact{
			ID:                   "test-fact-1",
			Text:                 "Implemented distributed caching system",
			SourceEventID:        "event-1",
			SourceBurstID:        "",
			CompetencyCategories: []string{"Technical"},
			RoleFit:              "Senior Engineer",
			AudienceRelevance:    []string{"Technical"},
			StrengthSignal:       "Strong",
			CreatedAt:            time.Now(),
			UpdatedAt:            time.Now(),
		}

		facts := []*career.Fact{testFact}
		model = NewFactsResultsModel(service, facts, ctx)
	})

	Describe("Creation", func() {
		It("creates model with facts", func() {
			Expect(model).NotTo(BeNil())
			Expect(len(model.facts)).To(Equal(1))
			Expect(model.currentIdx).To(Equal(0))
		})

		It("initializes empty confirmed and rejected lists", func() {
			Expect(len(model.confirmed)).To(Equal(0))
			Expect(len(model.rejected)).To(Equal(0))
		})
	})

	Describe("Navigation", func() {
		It("moves to next fact with down arrow", func() {
			fact2 := &career.Fact{
				ID:   "test-fact-2",
				Text: "Led team restructuring",
			}
			model.facts = append(model.facts, fact2)

			model.Update(tea.KeyMsg{Type: tea.KeyDown})
			Expect(model.currentIdx).To(Equal(1))
		})

		It("moves to previous fact with up arrow", func() {
			fact2 := &career.Fact{
				ID:   "test-fact-2",
				Text: "Led team restructuring",
			}
			model.facts = append(model.facts, fact2)
			model.currentIdx = 1

			model.Update(tea.KeyMsg{Type: tea.KeyUp})
			Expect(model.currentIdx).To(Equal(0))
		})

		It("doesn't go below 0 on up arrow", func() {
			model.Update(tea.KeyMsg{Type: tea.KeyUp})
			Expect(model.currentIdx).To(Equal(0))
		})

		It("doesn't exceed fact list length on down arrow", func() {
			model.Update(tea.KeyMsg{Type: tea.KeyDown})
			Expect(model.currentIdx).To(Equal(0))
		})
	})

	Describe("Fact Confirmation", func() {
		It("confirms fact with 'y' key", func() {
			model.Update(tea.KeyMsg{Type: tea.KeyRunes, Runes: []rune{'y'}})
			Expect(len(model.confirmed)).To(Equal(1))
			Expect(model.confirmed[0].ID).To(Equal("test-fact-1"))
		})

		It("confirms fact with Enter key", func() {
			model.Update(tea.KeyMsg{Type: tea.KeyEnter})
			Expect(len(model.confirmed)).To(Equal(1))
		})

		It("removes fact from list after confirmation", func() {
			Expect(len(model.facts)).To(Equal(1))
			model.Update(tea.KeyMsg{Type: tea.KeyRunes, Runes: []rune{'y'}})
			Expect(len(model.facts)).To(Equal(0))
		})
	})

	Describe("Fact Rejection", func() {
		It("rejects fact with 'n' key", func() {
			model.Update(tea.KeyMsg{Type: tea.KeyRunes, Runes: []rune{'n'}})
			Expect(len(model.rejected)).To(Equal(1))
			Expect(model.rejected[0].ID).To(Equal("test-fact-1"))
		})

		It("rejects fact with 'd' key", func() {
			model.Update(tea.KeyMsg{Type: tea.KeyRunes, Runes: []rune{'d'}})
			Expect(len(model.rejected)).To(Equal(1))
		})

		It("removes fact from list after rejection", func() {
			Expect(len(model.facts)).To(Equal(1))
			model.Update(tea.KeyMsg{Type: tea.KeyRunes, Runes: []rune{'n'}})
			Expect(len(model.facts)).To(Equal(0))
		})
	})

	Describe("Completion", func() {
		It("reports done when all facts reviewed", func() {
			model.Update(tea.KeyMsg{Type: tea.KeyRunes, Runes: []rune{'y'}})
			Expect(model.IsDone()).To(BeTrue())
		})

		It("returns confirmed facts", func() {
			model.Update(tea.KeyMsg{Type: tea.KeyRunes, Runes: []rune{'y'}})
			confirmed := model.GetConfirmed()
			Expect(len(confirmed)).To(Equal(1))
			Expect(confirmed[0].ID).To(Equal("test-fact-1"))
		})

		It("returns rejected facts", func() {
			model.Update(tea.KeyMsg{Type: tea.KeyRunes, Runes: []rune{'n'}})
			rejected := model.GetRejected()
			Expect(len(rejected)).To(Equal(1))
			Expect(rejected[0].ID).To(Equal("test-fact-1"))
		})
	})

	Describe("Rendering", func() {
		It("renders fact card", func() {
			view := model.View()
			Expect(view).NotTo(BeEmpty())
			Expect(view).To(ContainSubstring("Implemented distributed caching system"))
		})

		It("shows progress indicator", func() {
			view := model.View()
			Expect(view).To(ContainSubstring("[1/1]"))
		})

		It("shows help text", func() {
			view := model.View()
			Expect(view).To(ContainSubstring("Y/Enter: Confirm"))
			Expect(view).To(ContainSubstring("N/D: Reject"))
		})

		It("renders competency categories when present", func() {
			view := model.View()
			Expect(view).To(ContainSubstring("Technical"))
		})

		It("renders role fit when present", func() {
			view := model.View()
			Expect(view).To(ContainSubstring("Senior Engineer"))
		})

		It("renders audience relevance when present", func() {
			view := model.View()
			Expect(view).To(ContainSubstring("Technical"))
		})

		It("shows completion message when done", func() {
			model.Update(tea.KeyMsg{Type: tea.KeyRunes, Runes: []rune{'y'}})
			view := model.View()
			Expect(view).To(ContainSubstring("Complete"))
		})
	})

	Describe("Window Resizing", func() {
		It("updates dimensions on window resize", func() {
			model.Update(tea.WindowSizeMsg{Width: 120, Height: 30})
			Expect(model.width).To(Equal(120))
			Expect(model.height).To(Equal(30))
		})
	})

	Describe("Multiple Facts", func() {
		BeforeEach(func() {
			fact2 := &career.Fact{
				ID:                   "test-fact-2",
				Text:                 "Led team restructuring",
				SourceEventID:        "event-2",
				CompetencyCategories: []string{"Leadership"},
				RoleFit:              "Manager",
				AudienceRelevance:    []string{"Executive"},
				CreatedAt:            time.Now(),
				UpdatedAt:            time.Now(),
			}
			fact3 := &career.Fact{
				ID:                   "test-fact-3",
				Text:                 "Researched new technologies",
				SourceEventID:        "event-3",
				CompetencyCategories: []string{"Research"},
				RoleFit:              "Architect",
				AudienceRelevance:    []string{"Technical"},
				CreatedAt:            time.Now(),
				UpdatedAt:            time.Now(),
			}
			model.facts = append(model.facts, fact2, fact3)
		})

		It("navigates through multiple facts", func() {
			Expect(model.currentIdx).To(Equal(0))

			model.Update(tea.KeyMsg{Type: tea.KeyDown})
			Expect(model.currentIdx).To(Equal(1))

			model.Update(tea.KeyMsg{Type: tea.KeyDown})
			Expect(model.currentIdx).To(Equal(2))

			model.Update(tea.KeyMsg{Type: tea.KeyUp})
			Expect(model.currentIdx).To(Equal(1))
		})

		It("confirms some and rejects others", func() {
			// Confirm first
			model.Update(tea.KeyMsg{Type: tea.KeyRunes, Runes: []rune{'y'}})
			Expect(len(model.confirmed)).To(Equal(1))
			Expect(len(model.facts)).To(Equal(2))

			// Reject second
			model.Update(tea.KeyMsg{Type: tea.KeyRunes, Runes: []rune{'n'}})
			Expect(len(model.rejected)).To(Equal(1))
			Expect(len(model.facts)).To(Equal(1))

			// Confirm third
			model.Update(tea.KeyMsg{Type: tea.KeyRunes, Runes: []rune{'y'}})
			Expect(len(model.confirmed)).To(Equal(2))
			Expect(len(model.facts)).To(Equal(0))
		})

		It("shows updated progress for each fact", func() {
			view := model.View()
			Expect(view).To(ContainSubstring("[1/3]"))

			model.Update(tea.KeyMsg{Type: tea.KeyRunes, Runes: []rune{'y'}})
			// After confirming, it should show next fact progress
		})
	})

	Describe("Edge Cases", func() {
		It("handles empty fact list on creation", func() {
			emptyModel := NewFactsResultsModel(service, []*career.Fact{}, ctx)
			Expect(emptyModel.IsDone()).To(BeTrue())
		})

		It("shows completion on initialization if empty", func() {
			emptyModel := NewFactsResultsModel(service, []*career.Fact{}, ctx)
			view := emptyModel.View()
			Expect(view).To(ContainSubstring("Complete"))
		})
	})
})
