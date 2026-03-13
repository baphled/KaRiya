package factmanagement_test

import (
	"context"
	"errors"
	"time"

	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"

	"github.com/baphled/kariya/internal/domain/career"
	careerrepo "github.com/baphled/kariya/internal/repository/career"
	"github.com/baphled/kariya/internal/testutil/fixtures"
	"github.com/baphled/kariya/internal/tui/intents"
	"github.com/baphled/kariya/internal/tui/intents/factmanagement"
	"github.com/baphled/kariya/internal/ui/themes"
	"github.com/charmbracelet/bubbles/cursor"
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
		testFact.ID = "fact-123456789"
		testFact.Text = "Test fact about technical skills"
		testFact.RoleFit = "senior_ic"
		mockRepo.facts = []*career.Fact{testFact}
	})

	Describe("Construction", func() {
		Context("with valid context", func() {
			It("should create an intent successfully", func() {
				intentCtx := factmanagement.NewIntentValidator(ctx, mockRepo)
				intent, err := factmanagement.NewIntent(intentCtx)
				Expect(err).NotTo(HaveOccurred())
				Expect(intent).NotTo(BeNil())
			})

			It("should embed BaseIntent", func() {
				intentCtx := factmanagement.NewIntentValidator(ctx, mockRepo)
				intent, _ := factmanagement.NewIntent(intentCtx) //nolint:errcheck // test setup
				Expect(intent.BaseIntent).NotTo(BeNil())
			})
		})

		Context("with empty facts", func() {
			It("should create an intent with empty fact list", func() {
				emptyRepo := NewIntentMockFactRepository()
				intentCtx := factmanagement.NewIntentValidator(ctx, emptyRepo)
				intent, err := factmanagement.NewIntent(intentCtx)
				Expect(err).NotTo(HaveOccurred())
				Expect(intent).NotTo(BeNil())
			})
		})
	})

	Describe("Initialization", func() {
		BeforeEach(func() {
			intentCtx := factmanagement.NewIntentValidator(ctx, mockRepo)
			var err error
			intent, err = factmanagement.NewIntent(intentCtx)
			Expect(err).NotTo(HaveOccurred())
		})

		It("should return nil command on init", func() {
			cmd := intent.Init()
			Expect(cmd).To(BeNil())
		})

		It("applies the active theme to the table during init", func() {
			intent.SetThemeManager(themes.NewThemeManager())

			Expect(intent.Init()).To(BeNil())
			Expect(intent.View()).To(ContainSubstring("technical"))
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
				intentCtx := factmanagement.NewIntentValidator(ctx, mockRepo)
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
			intentCtx := factmanagement.NewIntentValidator(ctx, mockRepo)
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
				intentCtx := factmanagement.NewIntentValidator(ctx, emptyRepo)
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

		It("shows an inactive message before init", func() {
			freshIntent, err := factmanagement.NewIntent(factmanagement.NewIntentValidator(ctx, mockRepo))
			Expect(err).NotTo(HaveOccurred())

			Expect(freshIntent.View()).To(Equal("FactManagement intent is not active"))
		})

		It("stores a validation error on the intent when form errors exist", func() {
			intentCtx := factmanagement.NewIntentValidator(ctx, mockRepo)
			intentCtx.SetFormError("text", "Text is required")
			var err error
			intent, err = factmanagement.NewIntent(intentCtx)
			Expect(err).NotTo(HaveOccurred())
			intent.Init()

			_ = intent.View()

			Expect(intent.GetError()).To(HaveOccurred())
			Expect(intent.GetError().Error()).To(ContainSubstring("Text is required"))
		})

		It("renders truncated rows for long facts with empty strength", func() {
			longFact := fixtures.FactWith("fact-long", "This fact text is intentionally much longer than fifty characters to force truncation in the table")
			longFact.CompetencyCategories = []string{"technical", "leadership", "strategy", "delivery"}
			longFact.StrengthSignal = ""
			longFact.RoleFit = "senior_ic"
			mockRepo.facts = []*career.Fact{longFact}

			intentCtx := factmanagement.NewIntentValidator(ctx, mockRepo)
			var err error
			intent, err = factmanagement.NewIntent(intentCtx)
			Expect(err).NotTo(HaveOccurred())
			intent.Init()

			view := intent.View()
			Expect(view).To(ContainSubstring("This fact text is intentionally much longer tha…"))
			Expect(view).To(ContainSubstring("-"))
			Expect(view).To(ContainSubstring("[technical leadership"))
		})

	})

	Describe("Navigation", func() {
		BeforeEach(func() {
			for i := range 5 {
				f := fixtures.Fact("fact-"+string(rune('a'+i)), "event-1")
				f.Text = "Fact " + string(rune('a'+i))
				f.RoleFit = "senior_ic"
				mockRepo.facts = append(mockRepo.facts, f)
			}
			intentCtx := factmanagement.NewIntentValidator(ctx, mockRepo)
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
				for range 10 {
					intent.Update(tea.KeyMsg{Type: tea.KeyDown})
				}
				view := intent.View()
				Expect(view).NotTo(BeEmpty())
			})

			It("ignores non-key messages in list state", func() {
				Expect(intent.Update(tea.WindowSizeMsg{Width: 90, Height: 30})).To(BeNil())
			})
		})
	})

	Describe("Cancellation", func() {
		BeforeEach(func() {
			intentCtx := factmanagement.NewIntentValidator(ctx, mockRepo)
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

			It("shows help when question mark is pressed", func() {
				Expect(intent.IsHelpVisible()).To(BeFalse())

				cmd := intent.Update(tea.KeyMsg{Type: tea.KeyRunes, Runes: []rune{'?'}})

				Expect(cmd).To(BeNil())
				Expect(intent.IsHelpVisible()).To(BeTrue())
				Expect(intent.GetHelpModal()).NotTo(BeNil())
			})

			It("renders completed content after cancellation", func() {
				intent.Update(tea.KeyMsg{Type: tea.KeyEsc})

				Expect(intent.View()).To(ContainSubstring("Fact management completed"))
				Expect(intent.Update(tea.KeyMsg{Type: tea.KeyEnter})).To(BeNil())
			})

		})
	})

	Describe("Fact Selection", func() {
		BeforeEach(func() {
			intentCtx := factmanagement.NewIntentValidator(ctx, mockRepo)
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

		It("truncates the selected fact id in breadcrumbs", func() {
			intent.Update(tea.KeyMsg{Type: tea.KeyEnter})

			Expect(intent.View()).To(ContainSubstring("Fact #fact-123"))
		})

		It("should return to list after closing modal", func() {
			intent.Update(tea.KeyMsg{Type: tea.KeyEnter})
			intent.Update(tea.KeyMsg{Type: tea.KeyEsc})
			view := intent.View()
			Expect(view).To(ContainSubstring("Facts"))
		})

		It("ignores non-key messages in view state", func() {
			updateIntentForTest(intent, tea.KeyMsg{Type: tea.KeyEnter})

			Expect(intent.Update(tea.WindowSizeMsg{Width: 100, Height: 32})).To(BeNil())
			Expect(intent.View()).To(ContainSubstring("Fact Details"))
		})
	})

	Describe("Modal Interactions", func() {
		BeforeEach(func() {
			intentCtx := factmanagement.NewIntentValidator(ctx, mockRepo)
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

			It("returns to list when the new fact form completes without confirmation", func() {
				updateIntentForTest(intent, tea.KeyMsg{Type: tea.KeyRunes, Runes: []rune{'n'}})
				advanceIntentEditorToConfirm(intent)
				updateIntentForTest(intent, tea.KeyMsg{Type: tea.KeyEnter})

				Expect(intent.Result()).To(BeNil())
				Expect(intent.View()).To(ContainSubstring("Facts"))
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

			It("shows help from the editor state", func() {
				updateIntentForTest(intent, tea.KeyMsg{Type: tea.KeyRunes, Runes: []rune{'e'}})

				cmd := intent.Update(tea.KeyMsg{Type: tea.KeyRunes, Runes: []rune{'?'}})

				Expect(cmd).To(BeNil())
				Expect(intent.IsHelpVisible()).To(BeTrue())
			})

			It("returns to detail view when edit form completes without confirmation", func() {
				updateIntentForTest(intent, tea.KeyMsg{Type: tea.KeyEnter})
				detailView := intent.View()
				updateIntentForTest(intent, tea.KeyMsg{Type: tea.KeyRunes, Runes: []rune{'e'}})
				advanceIntentEditorToConfirm(intent)
				updateIntentForTest(intent, tea.KeyMsg{Type: tea.KeyEnter})

				Expect(intent.Result()).To(BeNil())
				Expect(intent.View()).To(Equal(detailView))
			})

			It("ignores non-key editor messages while modal stays open", func() {
				updateIntentForTest(intent, tea.KeyMsg{Type: tea.KeyRunes, Runes: []rune{'e'}})

				Expect(intent.Update(tea.WindowSizeMsg{Width: 110, Height: 35})).To(BeNil())
				Expect(intent.Result()).To(BeNil())
			})
		})

		Context("delete modal", func() {
			It("should open delete confirmation with d key", func() {
				intent.Update(tea.KeyMsg{Type: tea.KeyRunes, Runes: []rune{'d'}})
				view := intent.View()
				Expect(view).To(ContainSubstring("Delete"))
			})
		})

		Context("refresh", func() {
			It("reloads facts from the repository", func() {
				refreshedFact := fixtures.FactWith("fact-2", "Refreshed fact text")
				refreshedFact.RoleFit = "senior_ic"
				mockRepo.facts = append(mockRepo.facts, refreshedFact)

				updateIntentForTest(intent, tea.KeyMsg{Type: tea.KeyRunes, Runes: []rune{'r'}})

				Expect(intent.View()).To(ContainSubstring("Refreshed fact text"))
				Expect(intent.Result()).To(BeNil())
			})

			It("records a failed result when reload fails", func() {
				mockRepo.listErr = errors.New("reload failed")

				updateIntentForTest(intent, tea.KeyMsg{Type: tea.KeyRunes, Runes: []rune{'r'}})

				result := intent.Result()
				Expect(result).NotTo(BeNil())
				Expect(result.Status).To(Equal(intents.Failed))
				Expect(result.Error).To(HaveOccurred())
				Expect(result.Error.Code).To(Equal("RELOAD_FAILED"))
				Expect(result.Error.Cause).To(MatchError("reload failed"))
			})
		})

		Context("editor cancellation", func() {
			It("returns to fact detail when cancelling an edit from view", func() {
				updateIntentForTest(intent, tea.KeyMsg{Type: tea.KeyEnter})
				detailView := intent.View()

				updateIntentForTest(intent, tea.KeyMsg{Type: tea.KeyRunes, Runes: []rune{'e'}})
				updateIntentForTest(intent, tea.KeyMsg{Type: tea.KeyEsc})

				Expect(intent.View()).To(Equal(detailView))
			})

			It("returns to the list when cancelling a new fact", func() {
				updateIntentForTest(intent, tea.KeyMsg{Type: tea.KeyRunes, Runes: []rune{'n'}})
				updateIntentForTest(intent, tea.KeyMsg{Type: tea.KeyEsc})

				Expect(intent.View()).To(ContainSubstring("Facts"))
				Expect(intent.View()).To(ContainSubstring("technical"))
			})
		})

		Context("editor save flow", func() {
			It("updates an existing fact through the public editor flow", func() {
				updateIntentForTest(intent, tea.KeyMsg{Type: tea.KeyRunes, Runes: []rune{'e'}})
				updateIntentForTest(intent, tea.KeyMsg{Type: tea.KeyRunes, Runes: []rune(" updated")})
				advanceIntentEditorToConfirm(intent)
				updateIntentForTest(intent, tea.KeyMsg{Type: tea.KeyRight})
				updateIntentForTest(intent, tea.KeyMsg{Type: tea.KeyEnter})

				result := intent.Result()
				Expect(result).NotTo(BeNil())
				Expect(result.Status).To(Equal(intents.Completed))
				Expect(mockRepo.facts[0].Text).To(Equal("Test fact about technical skills updated"))
				Expect(intent.View()).To(ContainSubstring("technical skills updated"))
			})

			It("stores a save error on the intent when update fails", func() {
				mockRepo.updateErr = errors.New("update failed")

				updateIntentForTest(intent, tea.KeyMsg{Type: tea.KeyRunes, Runes: []rune{'e'}})
				updateIntentForTest(intent, tea.KeyMsg{Type: tea.KeyRunes, Runes: []rune(" updated")})
				advanceIntentEditorToConfirm(intent)
				updateIntentForTest(intent, tea.KeyMsg{Type: tea.KeyRight})
				updateIntentForTest(intent, tea.KeyMsg{Type: tea.KeyEnter})

				Expect(intent.Result()).To(BeNil())
				_ = intent.View()
				Expect(intent.GetError()).To(HaveOccurred())
				Expect(intent.GetError().Error()).To(ContainSubstring("Save failed: update failed"))
			})

			It("returns a quit command from editor state on q", func() {
				updateIntentForTest(intent, tea.KeyMsg{Type: tea.KeyRunes, Runes: []rune{'e'}})

				cmd := intent.Update(tea.KeyMsg{Type: tea.KeyRunes, Runes: []rune{'q'}})

				Expect(cmd).NotTo(BeNil())
			})

		})

		Context("view state controls", func() {
			It("opens delete confirmation from the detail view", func() {
				updateIntentForTest(intent, tea.KeyMsg{Type: tea.KeyEnter})
				updateIntentForTest(intent, tea.KeyMsg{Type: tea.KeyRunes, Runes: []rune{'d'}})

				Expect(intent.View()).To(ContainSubstring("Delete"))
			})

			It("shows help from the detail view", func() {
				updateIntentForTest(intent, tea.KeyMsg{Type: tea.KeyEnter})

				cmd := intent.Update(tea.KeyMsg{Type: tea.KeyRunes, Runes: []rune{'?'}})

				Expect(cmd).To(BeNil())
				Expect(intent.IsHelpVisible()).To(BeTrue())
			})

		})

		Context("delete confirmation flow", func() {
			It("deletes the selected fact when confirmed", func() {
				updateIntentForTest(intent, tea.KeyMsg{Type: tea.KeyRunes, Runes: []rune{'d'}})
				updateIntentForTest(intent, tea.KeyMsg{Type: tea.KeyRunes, Runes: []rune{'y'}})

				result := intent.Result()
				Expect(result).NotTo(BeNil())
				Expect(result.Status).To(Equal(intents.Completed))
				Expect(mockRepo.facts).To(BeEmpty())
				Expect(intent.View()).To(ContainSubstring("No facts"))
			})

			It("returns to detail view when deletion is cancelled", func() {
				updateIntentForTest(intent, tea.KeyMsg{Type: tea.KeyEnter})
				originalView := intent.View()
				updateIntentForTest(intent, tea.KeyMsg{Type: tea.KeyRunes, Runes: []rune{'d'}})
				updateIntentForTest(intent, tea.KeyMsg{Type: tea.KeyRunes, Runes: []rune{'n'}})

				Expect(intent.View()).To(Equal(originalView))
			})

			It("records a failed result when deletion fails", func() {
				mockRepo.deleteErr = errors.New("delete failed")
				updateIntentForTest(intent, tea.KeyMsg{Type: tea.KeyRunes, Runes: []rune{'d'}})
				updateIntentForTest(intent, tea.KeyMsg{Type: tea.KeyRunes, Runes: []rune{'y'}})

				result := intent.Result()
				Expect(result).NotTo(BeNil())
				Expect(result.Status).To(Equal(intents.Failed))
				Expect(result.Error).To(HaveOccurred())
				Expect(result.Error.Code).To(Equal("DELETE_FAILED"))
				Expect(result.Error.Cause).To(MatchError("delete failed"))
				Expect(intent.View()).To(ContainSubstring("Facts"))
			})

			It("ignores non-key messages in delete confirmation state", func() {
				updateIntentForTest(intent, tea.KeyMsg{Type: tea.KeyRunes, Runes: []rune{'d'}})

				Expect(intent.Update(tea.WindowSizeMsg{Width: 120, Height: 36})).To(BeNil())
				Expect(intent.View()).To(ContainSubstring("Delete"))
			})
		})
	})

	Describe("Window Resize", func() {
		BeforeEach(func() {
			intentCtx := factmanagement.NewIntentValidator(ctx, mockRepo)
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
			intentCtx := factmanagement.NewIntentValidator(ctx, mockRepo)
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

		It("returns nil for update on a zero-value intent with unknown state", func() {
			zeroIntent := new(factmanagement.Intent)

			Expect(zeroIntent.Update(tea.KeyMsg{Type: tea.KeyEnter})).To(BeNil())
		})

	})

	Describe("ListNavigator Interface", func() {
		BeforeEach(func() {
			intentCtx := factmanagement.NewIntentValidator(ctx, mockRepo)
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
})

func advanceIntentEditorToConfirm(intent *factmanagement.Intent) {
	for range 4 {
		updateIntentForTest(intent, tea.KeyMsg{Type: tea.KeyTab})
	}
}

func updateIntentForTest(intent *factmanagement.Intent, msg tea.Msg) {
	cmd := intent.Update(msg)
	flushIntentCmdForTest(intent, cmd)
}

func flushIntentCmdForTest(intent *factmanagement.Intent, cmd tea.Cmd) {
	for step := 0; cmd != nil && step < 20; step++ {
		msg := cmd()
		if msg == nil {
			return
		}
		if _, ok := msg.(cursor.BlinkMsg); ok {
			return
		}
		cmd = intent.Update(msg)
	}
}
