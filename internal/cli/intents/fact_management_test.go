package intents

import (
	"context"
	"time"
	"errors"
	"fmt"

	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"

	careerdom "github.com/baphled/kariya/internal/domain/career"
	careerrepo "github.com/baphled/kariya/internal/repository/career"

	tea "github.com/charmbracelet/bubbletea"
)

var _ = Describe("FactManagement Intent", func() {
	var (
		model    *FactManagementModel
		ctx      context.Context
		mockRepo *MockFactRepository
		testFact *careerdom.Fact
	)

	BeforeEach(func() {
		ctx = context.Background()
		mockRepo = NewMockFactRepository()

		testFact = &careerdom.Fact{
			ID:                   "fact-1",
			Text:                 "Test fact",
			CompetencyCategories: []string{"technical"},
			StrengthSignal:       "high",
			CreatedAt:            time.Now(),
			UpdatedAt:            time.Now(),
		}
		mockRepo.facts = []*careerdom.Fact{testFact}

		data := NewFactManagementContext(mockRepo, ctx)
		model = NewFactManagementIntent(data)
	})

	Describe("Initialization", func() {
		It("should initialize with list state", func() {
			cmd := model.Init()
			Expect(cmd).To(BeNil())
			Expect(model.View()).NotTo(BeEmpty())
		})

		It("should load facts on init", func() {
			cmd := model.Init()
			Expect(cmd).To(BeNil())
			view := model.View()
			Expect(view).To(ContainSubstring("Test fact"))
		})
	})

	Describe("List State", func() {
		BeforeEach(func() {
			model.Init()
		})

		It("should render list view", func() {
			view := model.View()
			Expect(view).To(ContainSubstring("Facts"))
			Expect(view).To(ContainSubstring("Test fact"))
		})

		It("should show empty state when no facts", func() {
			mockRepo.facts = []*careerdom.Fact{}
			model.Init()
			view := model.View()
			Expect(view).To(ContainSubstring("No facts found"))
		})

		It("should display pagination info", func() {
			view := model.View()
			Expect(view).To(ContainSubstring("Page"))
		})
	})

	Describe("Context Operations", func() {
		It("should load facts", func() {
			data := NewFactManagementContext(mockRepo, ctx)
			err := data.LoadFacts()
			Expect(err).To(BeNil())
			Expect(data.Facts).To(HaveLen(1))
		})

		It("should delete fact", func() {
			data := NewFactManagementContext(mockRepo, ctx)
			data.LoadFacts()
			initialCount := len(data.Facts)
			err := data.DeleteFact(testFact.ID)
			Expect(err).To(BeNil())
			Expect(data.Facts).To(HaveLen(initialCount - 1))
		})

		It("should handle pagination", func() {
			data := NewFactManagementContext(mockRepo, ctx)
			data.LoadFacts()
			pageFacts := data.GetPageFacts()
			Expect(pageFacts).NotTo(BeNil())
		})

		It("should track form errors", func() {
			data := NewFactManagementContext(mockRepo, ctx)
			data.SetFormError("text", "Text is required")
			Expect(data.HasFormErrors()).To(BeTrue())
		})

		It("should clear form errors", func() {
			data := NewFactManagementContext(mockRepo, ctx)
			data.SetFormError("text", "Text is required")
			data.ClearFormErrors()
			Expect(data.HasFormErrors()).To(BeFalse())
		})

		It("should toggle row expansion", func() {
			data := NewFactManagementContext(mockRepo, ctx)
			data.ToggleRowExpansion(0)
			Expect(data.IsRowExpanded(0)).To(BeTrue())
			data.ToggleRowExpansion(0)
			Expect(data.IsRowExpanded(0)).To(BeFalse())
		})

		It("should start new fact", func() {
			data := NewFactManagementContext(mockRepo, ctx)
			data.StartNewFact()
			Expect(data.IsNewFact).To(BeTrue())
			Expect(data.EditingFact).NotTo(BeNil())
		})

		It("should start edit fact", func() {
			data := NewFactManagementContext(mockRepo, ctx)
			data.StartEditFact(testFact)
			Expect(data.IsNewFact).To(BeFalse())
			Expect(data.EditingFact.ID).To(Equal(testFact.ID))
		})

		It("should cancel edit", func() {
			data := NewFactManagementContext(mockRepo, ctx)
			data.StartNewFact()
			data.CancelEdit()
			Expect(data.EditingFact).To(BeNil())
		})

		It("should handle nil repository", func() {
			data := NewFactManagementContext(nil, ctx)
			err := data.LoadFacts()
			Expect(err).NotTo(BeNil())
		})

		It("should handle empty fact list", func() {
			mockRepo.facts = []*careerdom.Fact{}
			data := NewFactManagementContext(mockRepo, ctx)
			err := data.LoadFacts()
			Expect(err).To(BeNil())
			Expect(data.Facts).To(HaveLen(0))
		})
	})

	Describe("View Methods", func() {
		It("should render view with facts", func() {
			model.Init()
			view := model.View()
			Expect(view).NotTo(BeEmpty())
		})
	})

	Describe("Result Handling", func() {
		It("should return result", func() {
			result := model.Result()
			Expect(result).NotTo(BeNil())
		})
	})
})

var _ = Describe("Pagination", func() {
	var (
		ctx context.Context
		manyFactsIntent *FactManagementModel
		manyFactsRepo   *MockFactRepository
	)

	BeforeEach(func() {
		ctx = context.Background()
		// Create 35 facts to span multiple pages (pageSize = 15)
		manyFactsRepo = NewMockFactRepository()
		for i := 0; i < 35; i++ {
			fact := &careerdom.Fact{
				ID:                   fmt.Sprintf("fact-%02d", i),
				Text:                 fmt.Sprintf("Fact %02d", i+1),
				CompetencyCategories: []string{fmt.Sprintf("competency-%d", i%5)},
				StrengthSignal:       fmt.Sprintf("signal-%d", i%3),
				CreatedAt:            time.Now(),
				UpdatedAt:            time.Now(),
			}
			manyFactsRepo.facts = append(manyFactsRepo.facts, fact)
		}

		data := NewFactManagementContext(manyFactsRepo, ctx)
		manyFactsIntent = NewFactManagementIntent(data)
		manyFactsIntent.Init()
	})

		It("should display correct facts on first page", func() {
			// Verify we're on page 1
			view := manyFactsIntent.View()
			Expect(view).To(ContainSubstring("Page 1 of 3"))

			// Verify table shows first 15 facts
			rows := manyFactsIntent.table.Rows()
			Expect(len(rows)).To(Equal(15))
		})

		It("should update table rows when navigating to next page", func() {
			// Navigate to page 2 using f key (pgdn)
			manyFactsIntent.Update(tea.KeyMsg{Type: tea.KeyRunes, Runes: []rune{'f'}})

			// Verify we're on page 2
			view := manyFactsIntent.View()
			Expect(view).To(ContainSubstring("Page 2 of 3"))

			// Verify the selected index is now in the second page range
			Expect(manyFactsIntent.data.SelectedFactIndex).To(Equal(15))

			// FAILING TEST: Verify table shows the correct facts for page 2
			// Currently the table shows ALL facts (all 35 rows) instead of just the current page (15 rows)
			rows := manyFactsIntent.table.Rows()
			Expect(len(rows)).To(Equal(15), "Table should show only 15 facts for page 2, but shows %d facts", len(rows))
		})

		It("should update table rows when navigating to last page", func() {
			// Navigate to page 3 using f key twice
			manyFactsIntent.Update(tea.KeyMsg{Type: tea.KeyRunes, Runes: []rune{'f'}})
			manyFactsIntent.Update(tea.KeyMsg{Type: tea.KeyRunes, Runes: []rune{'f'}})

			// Verify we're on page 3
			view := manyFactsIntent.View()
			Expect(view).To(ContainSubstring("Page 3 of 3"))

			// Verify the selected index is now in the third page range
			Expect(manyFactsIntent.data.SelectedFactIndex).To(Equal(30))

			// FAILING TEST: Verify table shows the correct facts for page 3
			// Currently the table shows ALL facts (all 35 rows) instead of just the current page (5 rows)
			rows := manyFactsIntent.table.Rows()
			Expect(len(rows)).To(Equal(5), "Table should show only 5 facts for page 3, but shows %d facts", len(rows))
		})
	})

// Pagination Describe ends here -- removed extra closing brace

type MockFactRepository struct {
	facts []*careerdom.Fact
}

func NewMockFactRepository() *MockFactRepository {
	return &MockFactRepository{
		facts: make([]*careerdom.Fact, 0),
	}
}

func (m *MockFactRepository) Create(ctx context.Context, fact *careerdom.Fact) error {
	if fact.ID == "" {
		fact.ID = "fact-" + time.Now().Format("20060102150405")
	}
	m.facts = append(m.facts, fact)
	return nil
}

func (m *MockFactRepository) GetByID(ctx context.Context, id string) (*careerdom.Fact, error) {
	for _, f := range m.facts {
		if f.ID == id {
			return f, nil
		}
	}
	return nil, errors.New("not found")
}

func (m *MockFactRepository) Update(ctx context.Context, fact *careerdom.Fact) error {
	for i, f := range m.facts {
		if f.ID == fact.ID {
			m.facts[i] = fact
			return nil
		}
	}
	return errors.New("not found")
}

func (m *MockFactRepository) Delete(ctx context.Context, id string) error {
	for i, f := range m.facts {
		if f.ID == id {
			m.facts = append(m.facts[:i], m.facts[i+1:]...)
			return nil
		}
	}
	return errors.New("not found")
}

func (m *MockFactRepository) List(ctx context.Context, filters careerrepo.FactListFilters) ([]*careerdom.Fact, error) {
	return m.facts, nil
}

func (m *MockFactRepository) Count(ctx context.Context, filters careerrepo.FactListFilters) (int, error) {
	return len(m.facts), nil
}

func (m *MockFactRepository) GetBySourceEventID(ctx context.Context, eventID string) ([]*careerdom.Fact, error) {
	return make([]*careerdom.Fact, 0), nil
}

func (m *MockFactRepository) GetBySourceBurstID(ctx context.Context, burstID string) ([]*careerdom.Fact, error) {
	return make([]*careerdom.Fact, 0), nil
}
