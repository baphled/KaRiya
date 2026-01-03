package intents_test

import (
	"context"
	"errors"
	"testing"
	"time"

	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"

	"github.com/baphled/kariya/internal/cli/intents"
	careerdom "github.com/baphled/kariya/internal/domain/career"
	careerrepo "github.com/baphled/kariya/internal/repository/career"
)

var ErrNotFound = errors.New("not found")

func TestFactManagement(t *testing.T) {
	RegisterFailHandler(Fail)
	RunSpecs(t, "FactManagement Intent Suite")
}

var _ = Describe("FactManagement Intent", func() {
	var (
		model *intents.FactManagementModel
		ctx context.Context
		mockRepo *MockFactRepository
		testFact *careerdom.Fact
	)

	BeforeEach(func() {
		ctx = context.Background()
		mockRepo = NewMockFactRepository()

		testFact = &careerdom.Fact{
			ID:   "fact-1",
			Text: "Test fact",
			CompetencyCategories: []string{"technical"},
			StrengthSignal: "high",
			CreatedAt: time.Now(),
			UpdatedAt: time.Now(),
		}
		mockRepo.facts = []*careerdom.Fact{testFact}

		data := intents.NewFactManagementContext(mockRepo, ctx)
		model = intents.NewFactManagementIntent(data)
	})

	Describe("Initialization", func() {
		It("should initialize with list state", func() {
			cmd := model.Init(ctx)
			Expect(cmd).To(BeNil())
			Expect(model.View()).NotTo(BeEmpty())
		})

		It("should load facts on init", func() {
			cmd := model.Init(ctx)
			Expect(cmd).To(BeNil())
			view := model.View()
			Expect(view).To(ContainSubstring("Test fact"))
		})
	})

	Describe("List State", func() {
		BeforeEach(func() {
			model.Init(ctx)
		})

		It("should render list view", func() {
			view := model.View()
			Expect(view).To(ContainSubstring("Facts"))
			Expect(view).To(ContainSubstring("Test fact"))
		})

		It("should show empty state when no facts", func() {
			mockRepo.facts = []*careerdom.Fact{}
			model.Init(ctx)
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
			data := intents.NewFactManagementContext(mockRepo, ctx)
			err := data.LoadFacts()
			Expect(err).To(BeNil())
			Expect(data.Facts).To(HaveLen(1))
		})



		It("should delete fact", func() {
			data := intents.NewFactManagementContext(mockRepo, ctx)
			data.LoadFacts()
			initialCount := len(data.Facts)
			err := data.DeleteFact(testFact.ID)
			Expect(err).To(BeNil())
			Expect(data.Facts).To(HaveLen(initialCount - 1))
		})

		It("should handle pagination", func() {
			data := intents.NewFactManagementContext(mockRepo, ctx)
			data.LoadFacts()
			pageFacts := data.GetPageFacts()
			Expect(pageFacts).NotTo(BeNil())
		})

		It("should track form errors", func() {
			data := intents.NewFactManagementContext(mockRepo, ctx)
			data.SetFormError("text", "Text is required")
			Expect(data.HasFormErrors()).To(BeTrue())
		})

		It("should clear form errors", func() {
			data := intents.NewFactManagementContext(mockRepo, ctx)
			data.SetFormError("text", "Text is required")
			data.ClearFormErrors()
			Expect(data.HasFormErrors()).To(BeFalse())
		})

		It("should toggle row expansion", func() {
			data := intents.NewFactManagementContext(mockRepo, ctx)
			data.ToggleRowExpansion(0)
			Expect(data.IsRowExpanded(0)).To(BeTrue())
			data.ToggleRowExpansion(0)
			Expect(data.IsRowExpanded(0)).To(BeFalse())
		})

		It("should start new fact", func() {
			data := intents.NewFactManagementContext(mockRepo, ctx)
			data.StartNewFact()
			Expect(data.IsNewFact).To(BeTrue())
			Expect(data.EditingFact).NotTo(BeNil())
		})

		It("should start edit fact", func() {
			data := intents.NewFactManagementContext(mockRepo, ctx)
			data.StartEditFact(testFact)
			Expect(data.IsNewFact).To(BeFalse())
			Expect(data.EditingFact.ID).To(Equal(testFact.ID))
		})

		It("should cancel edit", func() {
			data := intents.NewFactManagementContext(mockRepo, ctx)
			data.StartNewFact()
			data.CancelEdit()
			Expect(data.EditingFact).To(BeNil())
		})

		It("should handle nil repository", func() {
			data := intents.NewFactManagementContext(nil, ctx)
			err := data.LoadFacts()
			Expect(err).NotTo(BeNil())
		})

		It("should handle empty fact list", func() {
			mockRepo.facts = []*careerdom.Fact{}
			data := intents.NewFactManagementContext(mockRepo, ctx)
			err := data.LoadFacts()
			Expect(err).To(BeNil())
			Expect(data.Facts).To(HaveLen(0))
		})
	})

	Describe("View Methods", func() {
		It("should render view with facts", func() {
			model.Init(ctx)
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
	return nil, ErrNotFound
}

func (m *MockFactRepository) Update(ctx context.Context, fact *careerdom.Fact) error {
	for i, f := range m.facts {
		if f.ID == fact.ID {
			m.facts[i] = fact
			return nil
		}
	}
	return ErrNotFound
}

func (m *MockFactRepository) Delete(ctx context.Context, id string) error {
	for i, f := range m.facts {
		if f.ID == id {
			m.facts = append(m.facts[:i], m.facts[i+1:]...)
			return nil
		}
	}
	return ErrNotFound
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
