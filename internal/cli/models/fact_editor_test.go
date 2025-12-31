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

var _ = Describe("FactEditorModel", func() {
	var (
		repo   *careerrepo.MemoryRepository
		svc    *careerservice.Service
		ctx    context.Context
		editor *FactEditorModel
		fact   *career.Fact
	)

	BeforeEach(func() {
		repo = careerrepo.NewMemoryRepository()
		svc = careerservice.NewService(repo)
		ctx = context.Background()

		fact = &career.Fact{
			ID:                   "fact-123",
			Text:                 "Led technical migration with team",
			CompetencyCategories: []string{"technical"},
			RoleFit:              career.RoleFitPrincipal,
			AudienceRelevance:    []string{"hiring_manager"},
			StrengthSignal:       "leadership",
			SourceEventID:        "event-123",
			CreatedAt:            time.Now().Add(-24 * time.Hour),
			UpdatedAt:            time.Now(),
		}

		editor = NewFactEditorModel(fact, svc, ctx)
	})

	Describe("Creation", func() {
		It("creates FactEditorModel", func() {
			Expect(editor).NotTo(BeNil())
			Expect(editor.fact).To(Equal(fact))
			Expect(editor.submitted).To(BeFalse())
		})
	})

	Describe("Navigation", func() {
		It("navigates with Tab key", func() {
			msg := tea.KeyMsg{Type: tea.KeyTab}
			_, _ = editor.Update(msg)
			Expect(editor.focusIndex).To(Equal(FactCompetenciesFieldIdx))
		})

		It("cancels with Escape", func() {
			msg := tea.KeyMsg{Type: tea.KeyEsc}
			_, _ = editor.Update(msg)
			Expect(editor.cancelled).To(BeTrue())
		})
	})

	Describe("Validation", func() {
		It("rejects empty text", func() {
			editor.fact.Text = ""
			editor.focusIndex = FactSaveButtonIdx
			msg := tea.KeyMsg{Type: tea.KeyEnter}
			_, _ = editor.Update(msg)
			Expect(editor.submitted).To(BeFalse())
		})

		It("accepts valid submission", func() {
			editor.fact.Text = "Valid fact"
			editor.fact.CompetencyCategories = []string{"technical"}
			editor.fact.AudienceRelevance = []string{"peer"}
			editor.focusIndex = FactSaveButtonIdx
			msg := tea.KeyMsg{Type: tea.KeyEnter}
			_, _ = editor.Update(msg)
			Expect(editor.submitted).To(BeTrue())
		})

		It("detects aspirational language", func() {
			editor.fact.Text = "Will deliver feature"
			Expect(editor.hasAspirationLanguage(editor.fact.Text)).To(BeTrue())
		})
	})

	Describe("Revert", func() {
		It("reverts changes", func() {
			originalText := editor.fact.Text
			editor.fact.Text = "Modified"
			editor.Revert()
			Expect(editor.fact.Text).To(Equal(originalText))
		})
	})

	Describe("View", func() {
		It("renders", func() {
			editor.width = 120
			editor.height = 40
			view := editor.View()
			Expect(view).NotTo(BeEmpty())
			Expect(view).To(ContainSubstring("Fact Editor"))
		})
	})
})
