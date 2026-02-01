package factmanagement_test

import (
	"context"
	"errors"
	"time"

	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"

	"github.com/baphled/kariya/internal/cli/intents"
	"github.com/baphled/kariya/internal/cli/intents/factmanagement"
	"github.com/baphled/kariya/internal/domain/career"
	careerrepo "github.com/baphled/kariya/internal/repository/career"
	"github.com/baphled/kariya/internal/testutil/fixtures"
	tea "github.com/charmbracelet/bubbletea"
)

// MockFactRepository implements FactRepository for testing.
type IntentMockFactRepository struct {
	facts     []*career.Fact
	createErr error
	updateErr error
	deleteErr error
	listErr   error
}

func NewIntentMockFactRepository() *IntentMockFactRepository {
	return &IntentMockFactRepository{
		facts: make([]*career.Fact, 0),
	}
}

func (m *IntentMockFactRepository) Create(_ context.Context, fact *career.Fact) error {
	if m.createErr != nil {
		return m.createErr
	}
	if fact.ID == "" {
		fact.ID = "fact-" + time.Now().Format("20060102150405")
	}
	m.facts = append(m.facts, fact)
	return nil
}

func (m *IntentMockFactRepository) GetByID(_ context.Context, id string) (*career.Fact, error) {
	for _, f := range m.facts {
		if f.ID == id {
			return f, nil
		}
	}
	return nil, errors.New("not found")
}

func (m *IntentMockFactRepository) Update(_ context.Context, fact *career.Fact) error {
	if m.updateErr != nil {
		return m.updateErr
	}
	for i, f := range m.facts {
		if f.ID == fact.ID {
			m.facts[i] = fact
			return nil
		}
	}
	return errors.New("not found")
}

func (m *IntentMockFactRepository) Delete(_ context.Context, id string) error {
	if m.deleteErr != nil {
		return m.deleteErr
	}
	for i, f := range m.facts {
		if f.ID == id {
			m.facts = append(m.facts[:i], m.facts[i+1:]...)
			return nil
		}
	}
	return errors.New("not found")
}

func (m *IntentMockFactRepository) List(_ context.Context, _ careerrepo.FactListFilters) ([]*career.Fact, error) {
	if m.listErr != nil {
		return nil, m.listErr
	}
	return m.facts, nil
}

func (m *IntentMockFactRepository) Count(_ context.Context, _ careerrepo.FactListFilters) (int, error) {
	return len(m.facts), nil
}

func (m *IntentMockFactRepository) GetBySourceEventID(_ context.Context, _ string) ([]*career.Fact, error) {
	return make([]*career.Fact, 0), nil
}

func (m *IntentMockFactRepository) GetBySourceBurstID(_ context.Context, _ string) ([]*career.Fact, error) {
	return make([]*career.Fact, 0), nil
}

var _ = Describe("Intent", func() {
	var (
		intent   *factmanagement.Intent
		ctx      context.Context
		mockRepo *IntentMockFactRepository
		testFact *career.Fact
	)

	BeforeEach(func() {
		ctx = context.Background()
		mockRepo = NewIntentMockFactRepository()
		testFact = fixtures.Fact("fact-1", "event-1")
		testFact.Text = "Test fact about technical skills"
		testFact.RoleFit = "senior_ic"
		mockRepo.facts = []*career.Fact{testFact}
	})

	Describe("Construction", func() {
		Context("with valid context", func() {
			It("should create an intent successfully", func() {
				intentCtx := factmanagement.NewIntentContext(ctx, mockRepo)
				intent, err := factmanagement.NewIntent(intentCtx)
				Expect(err).NotTo(HaveOccurred())
				Expect(intent).NotTo(BeNil())
			})

			It("should embed BaseIntent", func() {
				intentCtx := factmanagement.NewIntentContext(ctx, mockRepo)
				intent, _ := factmanagement.NewIntent(intentCtx) //nolint:errcheck // test setup
				Expect(intent.BaseIntent).NotTo(BeNil())
			})
		})

		Context("with empty facts", func() {
			It("should create an intent with empty fact list", func() {
				emptyRepo := NewIntentMockFactRepository()
				intentCtx := factmanagement.NewIntentContext(ctx, emptyRepo)
				intent, err := factmanagement.NewIntent(intentCtx)
				Expect(err).NotTo(HaveOccurred())
				Expect(intent).NotTo(BeNil())
			})
		})
	})

	Describe("Initialization", func() {
		BeforeEach(func() {
			intentCtx := factmanagement.NewIntentContext(ctx, mockRepo)
			var err error
			intent, err = factmanagement.NewIntent(intentCtx)
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

		It("should show facts content after init", func() {
			intent.Init()
			view := intent.View()
			Expect(view).To(ContainSubstring("Facts"))
		})

		It("should display facts in view", func() {
			intent.Init()
			view := intent.View()
			Expect(view).To(ContainSubstring("technical"))
		})

		Context("when load fails", func() {
			It("should set failed result", func() {
				mockRepo.listErr = errors.New("database error")
				intentCtx := factmanagement.NewIntentContext(ctx, mockRepo)
				intent, _ = factmanagement.NewIntent(intentCtx) //nolint:errcheck // test setup
				cmd := intent.Init()
				Expect(cmd).NotTo(BeNil())
				result := intent.Result()
				Expect(result).NotTo(BeNil())
				Expect(result.Status).To(Equal(intents.Failed))
			})
		})
	})

	Describe("View Rendering", func() {
		BeforeEach(func() {
			intentCtx := factmanagement.NewIntentContext(ctx, mockRepo)
			var err error
			intent, err = factmanagement.NewIntent(intentCtx)
			Expect(err).NotTo(HaveOccurred())
			intent.Init()
		})

		Context("with facts", func() {
			It("should show fact count", func() {
				view := intent.View()
				Expect(view).To(ContainSubstring("Page"))
			})

			It("should show navigation hints", func() {
				view := intent.View()
				Expect(view).To(SatisfyAny(
					ContainSubstring("j/k"),
					ContainSubstring("Enter"),
				))
			})
		})

		Context("with empty facts", func() {
			BeforeEach(func() {
				emptyRepo := NewIntentMockFactRepository()
				intentCtx := factmanagement.NewIntentContext(ctx, emptyRepo)
				var err error
				intent, err = factmanagement.NewIntent(intentCtx)
				Expect(err).NotTo(HaveOccurred())
				intent.Init()
			})

			It("should show empty state message", func() {
				view := intent.View()
				Expect(view).To(ContainSubstring("No facts"))
			})
		})
	})

	Describe("Navigation", func() {
		BeforeEach(func() {
			for i := 0; i < 5; i++ {
				f := fixtures.Fact("fact-"+string(rune('a'+i)), "event-1")
				f.Text = "Fact " + string(rune('a'+i))
				f.RoleFit = "senior_ic"
				mockRepo.facts = append(mockRepo.facts, f)
			}
			intentCtx := factmanagement.NewIntentContext(ctx, mockRepo)
			var err error
			intent, err = factmanagement.NewIntent(intentCtx)
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
			intentCtx := factmanagement.NewIntentContext(ctx, mockRepo)
			var err error
			intent, err = factmanagement.NewIntent(intentCtx)
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
	})

	Describe("Fact Selection", func() {
		BeforeEach(func() {
			intentCtx := factmanagement.NewIntentContext(ctx, mockRepo)
			var err error
			intent, err = factmanagement.NewIntent(intentCtx)
			Expect(err).NotTo(HaveOccurred())
			intent.Init()
		})

		It("should show detail modal on enter", func() {
			intent.Update(tea.KeyMsg{Type: tea.KeyEnter})
			view := intent.View()
			Expect(view).To(ContainSubstring("Fact"))
		})

		It("should return to list after closing modal", func() {
			intent.Update(tea.KeyMsg{Type: tea.KeyEnter})
			intent.Update(tea.KeyMsg{Type: tea.KeyEsc})
			view := intent.View()
			Expect(view).To(ContainSubstring("Facts"))
		})
	})

	Describe("Modal Interactions", func() {
		BeforeEach(func() {
			intentCtx := factmanagement.NewIntentContext(ctx, mockRepo)
			var err error
			intent, err = factmanagement.NewIntent(intentCtx)
			Expect(err).NotTo(HaveOccurred())
			intent.Init()
		})

		Context("new fact modal", func() {
			It("should open new fact modal with n key", func() {
				intent.Update(tea.KeyMsg{Type: tea.KeyRunes, Runes: []rune{'n'}})
				view := intent.View()
				Expect(view).To(SatisfyAny(
					ContainSubstring("New"),
					ContainSubstring("Edit"),
				))
			})
		})

		Context("edit modal", func() {
			It("should open edit modal with e key", func() {
				intent.Update(tea.KeyMsg{Type: tea.KeyRunes, Runes: []rune{'e'}})
				view := intent.View()
				Expect(view).To(SatisfyAny(
					ContainSubstring("Edit"),
					ContainSubstring("Fact"),
				))
			})
		})

		Context("delete modal", func() {
			It("should open delete confirmation with d key", func() {
				intent.Update(tea.KeyMsg{Type: tea.KeyRunes, Runes: []rune{'d'}})
				view := intent.View()
				Expect(view).To(ContainSubstring("Delete"))
			})
		})
	})

	Describe("Window Resize", func() {
		BeforeEach(func() {
			intentCtx := factmanagement.NewIntentContext(ctx, mockRepo)
			var err error
			intent, err = factmanagement.NewIntent(intentCtx)
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

	Describe("Result Handling", func() {
		BeforeEach(func() {
			intentCtx := factmanagement.NewIntentContext(ctx, mockRepo)
			var err error
			intent, err = factmanagement.NewIntent(intentCtx)
			Expect(err).NotTo(HaveOccurred())
			intent.Init()
		})

		It("should return nil result when not completed", func() {
			result := intent.Result()
			Expect(result).To(BeNil())
		})

		It("should return cancelled result after escape", func() {
			intent.Update(tea.KeyMsg{Type: tea.KeyEsc})
			result := intent.Result()
			Expect(result).NotTo(BeNil())
			Expect(result.Status).To(Equal(intents.Cancelled))
		})
	})

	Describe("ListNavigator Interface", func() {
		BeforeEach(func() {
			intentCtx := factmanagement.NewIntentContext(ctx, mockRepo)
			var err error
			intent, err = factmanagement.NewIntent(intentCtx)
			Expect(err).NotTo(HaveOccurred())
			intent.Init()
		})

		It("should return total items count", func() {
			Expect(intent.GetTotalItems()).To(Equal(1))
		})

		It("should return selected index", func() {
			Expect(intent.GetSelectedIndex()).To(BeNumerically(">=", -1))
		})

		It("should return page size", func() {
			Expect(intent.GetPageSize()).To(Equal(15))
		})

		It("should set selected index", func() {
			intent.SetSelectedIndex(0)
			Expect(intent.GetSelectedIndex()).To(Equal(0))
		})
	})

	// Note: This intent doesn't use ScreenResultHandler as it uses
	// simple state-based navigation rather than the screens package.
})
