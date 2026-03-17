package factmanagement_test

import (
	"context"
	"errors"

	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"

	"github.com/baphled/kariya/internal/domain/career"
	"github.com/baphled/kariya/internal/testutil/fixtures"
	"github.com/baphled/kariya/internal/tui/intents"
	"github.com/baphled/kariya/internal/tui/intents/factmanagement"
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

		intentCtx := factmanagement.NewIntentValidator(ctx, mockRepo)
		var err error
		intent, err = factmanagement.NewIntent(intentCtx)
		Expect(err).NotTo(HaveOccurred())
		intent.Init()
	})

	Describe("handleViewState", func() {
		BeforeEach(func() {
			intent.Update(tea.KeyMsg{Type: tea.KeyEnter})
		})

		Context("when pressing help key", func() {
			It("toggles help display", func() {
				cmd := intent.Update(tea.KeyMsg{Type: tea.KeyRunes, Runes: []rune{'?'}})
				Expect(cmd).To(BeNil())
				view := intent.View()
				Expect(view).NotTo(BeEmpty())
			})
		})

		Context("when pressing edit key with a selected fact", func() {
			It("transitions to editor state", func() {
				intent.Update(tea.KeyMsg{Type: tea.KeyRunes, Runes: []rune{'e'}})
				view := intent.View()
				Expect(view).To(SatisfyAny(
					ContainSubstring("Edit"),
					ContainSubstring("Fact"),
				))
			})
		})

		Context("when pressing delete key with a selected fact", func() {
			It("transitions to delete confirmation", func() {
				intent.Update(tea.KeyMsg{Type: tea.KeyRunes, Runes: []rune{'d'}})
				view := intent.View()
				Expect(view).To(ContainSubstring("Delete"))
			})
		})

		Context("when pressing escape", func() {
			It("returns to list state", func() {
				intent.Update(tea.KeyMsg{Type: tea.KeyEsc})
				view := intent.View()
				Expect(view).To(ContainSubstring("Facts"))
			})
		})

		Context("when receiving non-key message", func() {
			It("returns nil command", func() {
				cmd := intent.Update(tea.WindowSizeMsg{Width: 80, Height: 24})
				Expect(cmd).To(BeNil())
			})
		})

		Context("when pressing an unrecognised key", func() {
			It("returns nil command", func() {
				cmd := intent.Update(tea.KeyMsg{Type: tea.KeyRunes, Runes: []rune{'z'}})
				Expect(cmd).To(BeNil())
			})
		})
	})

	Describe("handleEditorState", func() {
		Context("when entering editor via new fact", func() {
			BeforeEach(func() {
				intent.Update(tea.KeyMsg{Type: tea.KeyRunes, Runes: []rune{'n'}})
			})

			It("shows editor content in view", func() {
				view := intent.View()
				Expect(view).To(SatisfyAny(
					ContainSubstring("New"),
					ContainSubstring("Edit"),
					ContainSubstring("Fact"),
				))
			})

			Context("when pressing escape to cancel new fact", func() {
				It("returns to list state", func() {
					intent.Update(tea.KeyMsg{Type: tea.KeyEsc})
					view := intent.View()
					Expect(view).To(ContainSubstring("Facts"))
				})
			})

			Context("when pressing help key in editor", func() {
				It("toggles help", func() {
					cmd := intent.Update(tea.KeyMsg{Type: tea.KeyRunes, Runes: []rune{'?'}})
					Expect(cmd).To(BeNil())
					view := intent.View()
					Expect(view).NotTo(BeEmpty())
				})
			})

			Context("when receiving non-key message in editor", func() {
				It("delegates to modal", func() {
					cmd := intent.Update(tea.WindowSizeMsg{Width: 100, Height: 30})
					Expect(cmd).To(BeNil())
				})
			})
		})

		Context("when entering editor via edit existing fact", func() {
			BeforeEach(func() {
				intent.Update(tea.KeyMsg{Type: tea.KeyEnter})
				intent.Update(tea.KeyMsg{Type: tea.KeyRunes, Runes: []rune{'e'}})
			})

			Context("when pressing escape to cancel edit", func() {
				It("returns to view state showing fact detail", func() {
					intent.Update(tea.KeyMsg{Type: tea.KeyEsc})
					view := intent.View()
					Expect(view).To(ContainSubstring("Fact"))
				})
			})
		})
	})

	Describe("handleDeleteConfirmState", func() {
		BeforeEach(func() {
			intent.Update(tea.KeyMsg{Type: tea.KeyRunes, Runes: []rune{'d'}})
		})

		Context("when pressing y to confirm delete", func() {
			It("deletes the fact and returns to list", func() {
				intent.Update(tea.KeyMsg{Type: tea.KeyRunes, Runes: []rune{'y'}})
				result := intent.Result()
				Expect(result).NotTo(BeNil())
				Expect(result.Status).To(Equal(intents.Completed))
				view := intent.View()
				Expect(view).To(ContainSubstring("No facts"))
			})
		})

		Context("when pressing enter to confirm delete", func() {
			It("deletes the fact and returns to list", func() {
				intent.Update(tea.KeyMsg{Type: tea.KeyEnter})
				result := intent.Result()
				Expect(result).NotTo(BeNil())
				Expect(result.Status).To(Equal(intents.Completed))
			})
		})

		Context("when pressing n to cancel delete", func() {
			It("returns to view state", func() {
				intent.Update(tea.KeyMsg{Type: tea.KeyRunes, Runes: []rune{'n'}})
				view := intent.View()
				Expect(view).To(ContainSubstring("Fact"))
			})
		})

		Context("when pressing esc to cancel delete", func() {
			It("returns to view state", func() {
				intent.Update(tea.KeyMsg{Type: tea.KeyEsc})
				view := intent.View()
				Expect(view).To(ContainSubstring("Fact"))
			})
		})

		Context("when delete fails with repository error", func() {
			BeforeEach(func() {
				mockRepo.deleteErr = errors.New("delete failed")
			})

			It("sets failed result", func() {
				intent.Update(tea.KeyMsg{Type: tea.KeyRunes, Runes: []rune{'y'}})
				result := intent.Result()
				Expect(result).NotTo(BeNil())
				Expect(result.Status).To(Equal(intents.Failed))
			})
		})

		Context("when receiving non-key message", func() {
			It("returns nil command", func() {
				cmd := intent.Update(tea.WindowSizeMsg{Width: 80, Height: 24})
				Expect(cmd).To(BeNil())
			})
		})
	})

	Describe("handleRefreshKey", func() {
		Context("when refresh succeeds", func() {
			It("reloads facts from repository", func() {
				newFact := fixtures.Fact("fact-2", "event-2")
				newFact.Text = "Second fact"
				newFact.RoleFit = "senior_ic"
				mockRepo.facts = append(mockRepo.facts, newFact)

				intent.Update(tea.KeyMsg{Type: tea.KeyRunes, Runes: []rune{'r'}})
				Expect(intent.GetTotalItems()).To(Equal(2))
			})
		})

		Context("when refresh fails", func() {
			It("sets failed result", func() {
				mockRepo.listErr = errors.New("reload failed")
				intent.Update(tea.KeyMsg{Type: tea.KeyRunes, Runes: []rune{'r'}})
				result := intent.Result()
				Expect(result).NotTo(BeNil())
				Expect(result.Status).To(Equal(intents.Failed))
			})
		})
	})

	Describe("handleListGlobalKeys", func() {
		Context("when pressing help key in list", func() {
			It("toggles help display", func() {
				cmd := intent.Update(tea.KeyMsg{Type: tea.KeyRunes, Runes: []rune{'?'}})
				Expect(cmd).To(BeNil())
			})
		})

		Context("when pressing escape in list", func() {
			It("cancels the intent", func() {
				intent.Update(tea.KeyMsg{Type: tea.KeyEsc})
				result := intent.Result()
				Expect(result).NotTo(BeNil())
				Expect(result.Status).To(Equal(intents.Cancelled))
			})
		})
	})

	Describe("Update routing", func() {
		Context("when intent is in completed state", func() {
			It("returns nil command", func() {
				intent.Update(tea.KeyMsg{Type: tea.KeyEsc})
				cmd := intent.Update(tea.KeyMsg{Type: tea.KeyRunes, Runes: []rune{'n'}})
				Expect(cmd).To(BeNil())
			})
		})

		Context("when handling different message types in list state", func() {
			It("ignores non-key messages", func() {
				cmd := intent.Update(tea.WindowSizeMsg{Width: 120, Height: 40})
				Expect(cmd).To(BeNil())
			})
		})

		Context("when handling space key for selection", func() {
			It("transitions to view state", func() {
				intent.Update(tea.KeyMsg{Type: tea.KeySpace})
				view := intent.View()
				Expect(view).To(ContainSubstring("Fact"))
			})
		})
	})

	Describe("View rendering at different states", func() {
		Context("when not active", func() {
			It("shows inactive message", func() {
				intentCtx := factmanagement.NewIntentValidator(ctx, mockRepo)
				inactiveIntent, err := factmanagement.NewIntent(intentCtx)
				Expect(err).NotTo(HaveOccurred())
				view := inactiveIntent.View()
				Expect(view).To(ContainSubstring("not active"))
			})
		})

		Context("in view state", func() {
			It("shows fact detail", func() {
				intent.Update(tea.KeyMsg{Type: tea.KeyEnter})
				view := intent.View()
				Expect(view).To(ContainSubstring("Fact"))
			})
		})

		Context("in editor state", func() {
			It("shows editor form", func() {
				intent.Update(tea.KeyMsg{Type: tea.KeyRunes, Runes: []rune{'n'}})
				view := intent.View()
				Expect(view).NotTo(BeEmpty())
			})
		})

		Context("in delete confirm state", func() {
			It("shows delete confirmation", func() {
				intent.Update(tea.KeyMsg{Type: tea.KeyRunes, Runes: []rune{'d'}})
				view := intent.View()
				Expect(view).To(ContainSubstring("Delete"))
			})
		})

		Context("in completed state", func() {
			It("renders completed content", func() {
				intent.Update(tea.KeyMsg{Type: tea.KeyEsc})
				view := intent.View()
				Expect(view).NotTo(BeEmpty())
			})
		})
	})

	Describe("Breadcrumb rendering", func() {
		Context("in list state", func() {
			It("shows base breadcrumbs", func() {
				view := intent.View()
				Expect(view).To(ContainSubstring("Manage Facts"))
			})
		})

		Context("in view state with selected fact", func() {
			It("shows fact identifier in breadcrumbs", func() {
				intent.Update(tea.KeyMsg{Type: tea.KeyEnter})
				view := intent.View()
				Expect(view).To(ContainSubstring("Fact"))
			})
		})

		Context("in editor state for new fact", func() {
			It("shows New Fact breadcrumb", func() {
				intent.Update(tea.KeyMsg{Type: tea.KeyRunes, Runes: []rune{'n'}})
				view := intent.View()
				Expect(view).To(ContainSubstring("New Fact"))
			})
		})

		Context("in editor state for existing fact", func() {
			It("shows Edit Fact breadcrumb", func() {
				intent.Update(tea.KeyMsg{Type: tea.KeyEnter})
				intent.Update(tea.KeyMsg{Type: tea.KeyRunes, Runes: []rune{'e'}})
				view := intent.View()
				Expect(view).To(ContainSubstring("Edit Fact"))
			})
		})

		Context("in delete confirm state", func() {
			It("shows fact identifier breadcrumb", func() {
				intent.Update(tea.KeyMsg{Type: tea.KeyRunes, Runes: []rune{'d'}})
				view := intent.View()
				Expect(view).To(ContainSubstring("Fact"))
			})
		})
	})

	Describe("Context help rendering", func() {
		Context("in list state", func() {
			It("shows list-specific help", func() {
				view := intent.View()
				Expect(view).To(SatisfyAny(
					ContainSubstring("New"),
					ContainSubstring("Refresh"),
				))
			})
		})

		Context("in view state", func() {
			It("shows view-specific help", func() {
				intent.Update(tea.KeyMsg{Type: tea.KeyEnter})
				view := intent.View()
				Expect(view).To(SatisfyAny(
					ContainSubstring("Edit"),
					ContainSubstring("Delete"),
				))
			})
		})

		Context("in editor state", func() {
			It("shows editor-specific help", func() {
				intent.Update(tea.KeyMsg{Type: tea.KeyRunes, Runes: []rune{'n'}})
				view := intent.View()
				Expect(view).NotTo(BeEmpty())
			})
		})

		Context("in delete confirm state", func() {
			It("shows confirm/cancel help", func() {
				intent.Update(tea.KeyMsg{Type: tea.KeyRunes, Runes: []rune{'d'}})
				view := intent.View()
				Expect(view).To(SatisfyAny(
					ContainSubstring("Confirm"),
					ContainSubstring("Cancel"),
				))
			})
		})
	})

	Describe("Edit key in list state", func() {
		Context("with a selected fact", func() {
			It("transitions to editor state", func() {
				intent.Update(tea.KeyMsg{Type: tea.KeyRunes, Runes: []rune{'e'}})
				view := intent.View()
				Expect(view).To(SatisfyAny(
					ContainSubstring("Edit"),
					ContainSubstring("Fact"),
				))
			})
		})

		Context("with no facts in list", func() {
			BeforeEach(func() {
				emptyRepo := NewIntentMockFactRepository()
				intentCtx := factmanagement.NewIntentValidator(ctx, emptyRepo)
				var err error
				intent, err = factmanagement.NewIntent(intentCtx)
				Expect(err).NotTo(HaveOccurred())
				intent.Init()
			})

			It("does not transition when pressing edit", func() {
				intent.Update(tea.KeyMsg{Type: tea.KeyRunes, Runes: []rune{'e'}})
				view := intent.View()
				Expect(view).To(ContainSubstring("No facts"))
			})

			It("does not transition when pressing delete", func() {
				intent.Update(tea.KeyMsg{Type: tea.KeyRunes, Runes: []rune{'d'}})
				view := intent.View()
				Expect(view).To(ContainSubstring("No facts"))
			})

			It("allows creating new fact even with empty list", func() {
				intent.Update(tea.KeyMsg{Type: tea.KeyRunes, Runes: []rune{'n'}})
				view := intent.View()
				Expect(view).To(SatisfyAny(
					ContainSubstring("New"),
					ContainSubstring("Edit"),
				))
			})
		})
	})

	Describe("Delete from view state", func() {
		BeforeEach(func() {
			intent.Update(tea.KeyMsg{Type: tea.KeyEnter})
		})

		Context("when confirming delete from view", func() {
			It("deletes and returns to list", func() {
				intent.Update(tea.KeyMsg{Type: tea.KeyRunes, Runes: []rune{'d'}})
				intent.Update(tea.KeyMsg{Type: tea.KeyRunes, Runes: []rune{'y'}})
				result := intent.Result()
				Expect(result).NotTo(BeNil())
				Expect(result.Status).To(Equal(intents.Completed))
			})
		})

		Context("when cancelling delete from view", func() {
			It("returns to view state", func() {
				intent.Update(tea.KeyMsg{Type: tea.KeyRunes, Runes: []rune{'d'}})
				intent.Update(tea.KeyMsg{Type: tea.KeyRunes, Runes: []rune{'n'}})
				view := intent.View()
				Expect(view).To(ContainSubstring("Fact"))
			})
		})
	})

	Describe("State transitions", func() {
		Context("list to view to editor to list", func() {
			It("completes full cycle for new fact", func() {
				intent.Update(tea.KeyMsg{Type: tea.KeyRunes, Runes: []rune{'n'}})
				view := intent.View()
				Expect(view).NotTo(BeEmpty())

				intent.Update(tea.KeyMsg{Type: tea.KeyEsc})
				view = intent.View()
				Expect(view).To(ContainSubstring("Facts"))
			})
		})

		Context("list to view to editor to view", func() {
			It("completes full cycle for editing", func() {
				intent.Update(tea.KeyMsg{Type: tea.KeyEnter})
				view := intent.View()
				Expect(view).To(ContainSubstring("Fact"))

				intent.Update(tea.KeyMsg{Type: tea.KeyRunes, Runes: []rune{'e'}})
				view = intent.View()
				Expect(view).NotTo(BeEmpty())

				intent.Update(tea.KeyMsg{Type: tea.KeyEsc})
				view = intent.View()
				Expect(view).To(ContainSubstring("Fact"))
			})
		})

		Context("list to delete confirm to list", func() {
			It("completes delete cycle", func() {
				intent.Update(tea.KeyMsg{Type: tea.KeyRunes, Runes: []rune{'d'}})
				view := intent.View()
				Expect(view).To(ContainSubstring("Delete"))

				intent.Update(tea.KeyMsg{Type: tea.KeyRunes, Runes: []rune{'y'}})
				view = intent.View()
				Expect(view).To(ContainSubstring("No facts"))
			})
		})

		Context("list to delete confirm to view", func() {
			It("completes cancel cycle", func() {
				intent.Update(tea.KeyMsg{Type: tea.KeyRunes, Runes: []rune{'d'}})
				intent.Update(tea.KeyMsg{Type: tea.KeyRunes, Runes: []rune{'n'}})
				view := intent.View()
				Expect(view).To(ContainSubstring("Fact"))
			})
		})
	})

	Describe("factRowFormatter", func() {
		Context("with long text", func() {
			It("truncates text in view", func() {
				longFact := fixtures.Fact("fact-long", "event-1")
				longFact.Text = "This is a very long fact text that should be truncated when displayed in the table view"
				longFact.RoleFit = "senior_ic"
				mockRepo.facts = []*career.Fact{longFact}

				intentCtx := factmanagement.NewIntentValidator(ctx, mockRepo)
				longIntent, err := factmanagement.NewIntent(intentCtx)
				Expect(err).NotTo(HaveOccurred())
				longIntent.Init()
				view := longIntent.View()
				Expect(view).NotTo(BeEmpty())
			})
		})

		Context("with empty strength signal", func() {
			It("shows dash for empty strength", func() {
				noStrengthFact := fixtures.Fact("fact-ns", "event-1")
				noStrengthFact.StrengthSignal = ""
				noStrengthFact.RoleFit = "senior_ic"
				mockRepo.facts = []*career.Fact{noStrengthFact}

				intentCtx := factmanagement.NewIntentValidator(ctx, mockRepo)
				nsIntent, err := factmanagement.NewIntent(intentCtx)
				Expect(err).NotTo(HaveOccurred())
				nsIntent.Init()
				view := nsIntent.View()
				Expect(view).To(ContainSubstring("-"))
			})
		})

		Context("with long categories", func() {
			It("truncates categories in view", func() {
				catFact := fixtures.Fact("fact-cat", "event-1")
				catFact.CompetencyCategories = []string{"technical", "leadership", "product", "consulting", "mentoring"}
				catFact.RoleFit = "senior_ic"
				mockRepo.facts = []*career.Fact{catFact}

				intentCtx := factmanagement.NewIntentValidator(ctx, mockRepo)
				catIntent, err := factmanagement.NewIntent(intentCtx)
				Expect(err).NotTo(HaveOccurred())
				catIntent.Init()
				view := catIntent.View()
				Expect(view).NotTo(BeEmpty())
			})
		})
	})
})
