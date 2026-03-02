package factmanagement_test

import (
	"context"
	"errors"

	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"

	"github.com/baphled/kariya/internal/cli/intents"
	"github.com/baphled/kariya/internal/cli/intents/factmanagement"
	"github.com/baphled/kariya/internal/domain/career"
	"github.com/baphled/kariya/internal/testutil/fixtures"
	tea "github.com/charmbracelet/bubbletea"
)

var _ = Describe("Handlers", func() {
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

		intentCtx := factmanagement.NewIntentContext(ctx, mockRepo)
		var err error
		intent, err = factmanagement.NewIntent(intentCtx)
		Expect(err).NotTo(HaveOccurred())
		intent.Init()
	})

	Describe("handleRefreshKey", func() {
		Context("when refresh succeeds", func() {
			It("should reload facts from the repository", func() {
				newFact := fixtures.Fact("fact-2", "event-2")
				newFact.Text = "Another technical fact"
				newFact.RoleFit = "senior_ic"
				mockRepo.facts = append(mockRepo.facts, newFact)

				cmd := intent.Update(tea.KeyMsg{Type: tea.KeyRunes, Runes: []rune{'r'}})
				Expect(cmd).To(BeNil())
				Expect(intent.GetTotalItems()).To(Equal(2))
			})
		})

		Context("when refresh fails", func() {
			It("should set a failed result with reload error", func() {
				mockRepo.listErr = errors.New("database error")
				intent.Update(tea.KeyMsg{Type: tea.KeyRunes, Runes: []rune{'r'}})
				result := intent.Result()
				Expect(result).NotTo(BeNil())
				Expect(result.Status).To(Equal(intents.Failed))
			})
		})
	})

	Describe("handleViewState", func() {
		BeforeEach(func() {
			intent.Update(tea.KeyMsg{Type: tea.KeyEnter})
		})

		Context("when pressing back key", func() {
			It("should return to list state", func() {
				intent.Update(tea.KeyMsg{Type: tea.KeyEsc})
				view := intent.View()
				Expect(view).To(ContainSubstring("Facts"))
			})
		})

		Context("when pressing help key", func() {
			It("should toggle help without error", func() {
				cmd := intent.Update(tea.KeyMsg{Type: tea.KeyRunes, Runes: []rune{'?'}})
				Expect(cmd).To(BeNil())
			})
		})

		Context("when pressing edit key", func() {
			It("should transition to editor", func() {
				intent.Update(tea.KeyMsg{Type: tea.KeyRunes, Runes: []rune{'e'}})
				view := intent.View()
				Expect(view).To(SatisfyAny(
					ContainSubstring("Edit"),
					ContainSubstring("Fact"),
				))
			})
		})

		Context("when pressing delete key", func() {
			It("should transition to delete confirm", func() {
				intent.Update(tea.KeyMsg{Type: tea.KeyRunes, Runes: []rune{'d'}})
				view := intent.View()
				Expect(view).To(ContainSubstring("Delete"))
			})
		})

		Context("with non-key message", func() {
			It("should return nil", func() {
				cmd := intent.Update(tea.WindowSizeMsg{Width: 80, Height: 24})
				Expect(cmd).To(BeNil())
			})
		})
	})

	Describe("handleEditorState", func() {
		Context("when cancelling new fact editor", func() {
			It("should return to list state", func() {
				intent.Update(tea.KeyMsg{Type: tea.KeyRunes, Runes: []rune{'n'}})
				intent.Update(tea.KeyMsg{Type: tea.KeyEsc})
				view := intent.View()
				Expect(view).To(ContainSubstring("Facts"))
			})
		})

		Context("when cancelling edit from view", func() {
			It("should return to view state", func() {
				intent.Update(tea.KeyMsg{Type: tea.KeyEnter})
				intent.Update(tea.KeyMsg{Type: tea.KeyRunes, Runes: []rune{'e'}})
				intent.Update(tea.KeyMsg{Type: tea.KeyEsc})
				view := intent.View()
				Expect(view).To(ContainSubstring("Fact"))
			})
		})

		Context("when pressing help key in editor", func() {
			It("should toggle help", func() {
				intent.Update(tea.KeyMsg{Type: tea.KeyRunes, Runes: []rune{'n'}})
				cmd := intent.Update(tea.KeyMsg{Type: tea.KeyRunes, Runes: []rune{'?'}})
				Expect(cmd).To(BeNil())
			})
		})

		Context("with non-key message in editor", func() {
			It("should delegate to modal without error", func() {
				intent.Update(tea.KeyMsg{Type: tea.KeyRunes, Runes: []rune{'n'}})
				Expect(func() {
					intent.Update(tea.WindowSizeMsg{Width: 80, Height: 24})
				}).NotTo(Panic())
			})
		})

		Context("with a regular key in editor", func() {
			It("should pass key to modal without crashing", func() {
				intent.Update(tea.KeyMsg{Type: tea.KeyRunes, Runes: []rune{'n'}})
				Expect(func() {
					intent.Update(tea.KeyMsg{Type: tea.KeyRunes, Runes: []rune{'a'}})
				}).NotTo(Panic())
			})
		})
	})

	Describe("handleDeleteConfirmState", func() {
		BeforeEach(func() {
			intent.Update(tea.KeyMsg{Type: tea.KeyEnter})
			intent.Update(tea.KeyMsg{Type: tea.KeyRunes, Runes: []rune{'d'}})
		})

		Context("when confirming delete", func() {
			It("should delete the fact and set completed result", func() {
				intent.Update(tea.KeyMsg{Type: tea.KeyRunes, Runes: []rune{'y'}})
				result := intent.Result()
				Expect(result).NotTo(BeNil())
				Expect(result.Status).To(Equal(intents.Completed))
				Expect(intent.GetTotalItems()).To(Equal(0))
			})
		})

		Context("when delete fails", func() {
			It("should set failed result", func() {
				mockRepo.deleteErr = errors.New("delete error")
				intent.Update(tea.KeyMsg{Type: tea.KeyRunes, Runes: []rune{'y'}})
				result := intent.Result()
				Expect(result).NotTo(BeNil())
				Expect(result.Status).To(Equal(intents.Failed))
			})
		})

		Context("when cancelling with n key", func() {
			It("should return to view state", func() {
				intent.Update(tea.KeyMsg{Type: tea.KeyRunes, Runes: []rune{'n'}})
				view := intent.View()
				Expect(view).To(ContainSubstring("Fact"))
			})
		})

		Context("when cancelling with esc key", func() {
			It("should return to view state", func() {
				intent.Update(tea.KeyMsg{Type: tea.KeyEsc})
				view := intent.View()
				Expect(view).To(ContainSubstring("Fact"))
			})
		})

		Context("with non-key message", func() {
			It("should return nil", func() {
				cmd := intent.Update(tea.WindowSizeMsg{Width: 80, Height: 24})
				Expect(cmd).To(BeNil())
			})
		})
	})

	Describe("handleListGlobalKeys", func() {
		Context("when pressing help key in list", func() {
			It("should toggle help without error", func() {
				cmd := intent.Update(tea.KeyMsg{Type: tea.KeyRunes, Runes: []rune{'?'}})
				Expect(cmd).To(BeNil())
			})
		})
	})

	Describe("Update state routing", func() {
		Context("when intent is in completed state", func() {
			It("should return nil for any message", func() {
				intent.Update(tea.KeyMsg{Type: tea.KeyEsc})
				cmd := intent.Update(tea.KeyMsg{Type: tea.KeyRunes, Runes: []rune{'x'}})
				Expect(cmd).To(BeNil())
			})
		})
	})

	Describe("View edge cases", func() {
		Context("when intent is not active", func() {
			It("should show inactive message", func() {
				inactiveCtx := factmanagement.NewIntentContext(ctx, mockRepo)
				inactiveIntent, err := factmanagement.NewIntent(inactiveCtx)
				Expect(err).NotTo(HaveOccurred())
				view := inactiveIntent.View()
				Expect(view).To(ContainSubstring("not active"))
			})
		})
	})

	Describe("handleEditKeyInList with empty list", func() {
		It("should return nil when no fact is selected", func() {
			emptyRepo := NewIntentMockFactRepository()
			intentCtx := factmanagement.NewIntentContext(ctx, emptyRepo)
			emptyIntent, err := factmanagement.NewIntent(intentCtx)
			Expect(err).NotTo(HaveOccurred())
			emptyIntent.Init()
			cmd := emptyIntent.Update(tea.KeyMsg{Type: tea.KeyRunes, Runes: []rune{'e'}})
			Expect(cmd).To(BeNil())
		})
	})

	Describe("handleListKeyActions with empty list", func() {
		It("should not transition to delete confirm when no fact is selected", func() {
			emptyRepo := NewIntentMockFactRepository()
			intentCtx := factmanagement.NewIntentContext(ctx, emptyRepo)
			emptyIntent, err := factmanagement.NewIntent(intentCtx)
			Expect(err).NotTo(HaveOccurred())
			emptyIntent.Init()
			cmd := emptyIntent.Update(tea.KeyMsg{Type: tea.KeyRunes, Runes: []rune{'d'}})
			Expect(cmd).To(BeNil())
		})

		It("should not transition to view when no fact is selected on enter", func() {
			emptyRepo := NewIntentMockFactRepository()
			intentCtx := factmanagement.NewIntentContext(ctx, emptyRepo)
			emptyIntent, err := factmanagement.NewIntent(intentCtx)
			Expect(err).NotTo(HaveOccurred())
			emptyIntent.Init()
			cmd := emptyIntent.Update(tea.KeyMsg{Type: tea.KeyEnter})
			Expect(cmd).To(BeNil())
		})
	})
})
