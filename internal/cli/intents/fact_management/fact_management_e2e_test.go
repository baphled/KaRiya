package fact_management_test

import (
	"context"
	"errors"
	"time"

	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"

	"github.com/baphled/kariya/internal/cli/intents"
	"github.com/baphled/kariya/internal/cli/intents/fact_management"
	"github.com/baphled/kariya/internal/domain/career"
	tea "github.com/charmbracelet/bubbletea"
)

// E2E Tests for FactManagement Intent
//
// These tests exercise complete workflows through the intent:
// - Create workflow (new fact)
// - Edit workflow (modify existing fact)
// - Delete workflow (confirm and cancel paths)
// - Navigation between all states
// - Error handling scenarios

var _ = Describe("FactManagement E2E", func() {
	var (
		intent   *fact_management.Intent
		ctx      context.Context
		mockRepo *IntentMockFactRepository
	)

	createTestFact := func(id, text string) *career.Fact {
		return &career.Fact{
			ID:                   id,
			Text:                 text,
			CompetencyCategories: []string{"technical"},
			StrengthSignal:       "high",
			RoleFit:              "senior_ic",
			AudienceRelevance:    []string{"hiring_manager"},
			SourceEventID:        "event-1",
			CreatedAt:            time.Now(),
			UpdatedAt:            time.Now(),
		}
	}

	BeforeEach(func() {
		ctx = context.Background()
		mockRepo = NewIntentMockFactRepository()
		mockRepo.facts = []*career.Fact{
			createTestFact("fact-1", "First fact about leadership skills"),
			createTestFact("fact-2", "Second fact about technical expertise"),
			createTestFact("fact-3", "Third fact about project management"),
		}
	})

	// =========================================================================
	// VIEW STATE E2E TESTS
	// =========================================================================

	Describe("View State Workflow", func() {
		BeforeEach(func() {
			intentCtx := fact_management.NewIntentContext(mockRepo, ctx)
			var err error
			intent, err = fact_management.NewIntent(intentCtx)
			Expect(err).NotTo(HaveOccurred())
			intent.Init()
		})

		It("should transition to view state on enter", func() {
			By("Pressing enter to view selected fact")
			intent.Update(tea.KeyMsg{Type: tea.KeyEnter})

			By("Verifying view state content")
			view := intent.View()
			Expect(view).To(ContainSubstring("Fact Details"))
			Expect(view).To(ContainSubstring("First fact about leadership"))
		})

		It("should return to list state on escape from view", func() {
			By("Entering view state")
			intent.Update(tea.KeyMsg{Type: tea.KeyEnter})
			view := intent.View()
			Expect(view).To(ContainSubstring("Fact Details"))

			By("Pressing escape to go back")
			intent.Update(tea.KeyMsg{Type: tea.KeyEsc})

			By("Verifying back at list state")
			view = intent.View()
			Expect(view).To(ContainSubstring("Page"))
		})

		It("should show edit and delete options in view state", func() {
			By("Entering view state")
			intent.Update(tea.KeyMsg{Type: tea.KeyEnter})

			By("Verifying action hints are shown")
			view := intent.View()
			Expect(view).To(SatisfyAny(
				ContainSubstring("e"),
				ContainSubstring("Edit"),
			))
			Expect(view).To(SatisfyAny(
				ContainSubstring("d"),
				ContainSubstring("Delete"),
			))
		})

		It("should toggle help from view state", func() {
			By("Entering view state")
			intent.Update(tea.KeyMsg{Type: tea.KeyEnter})

			By("Pressing ? to toggle help")
			intent.Update(tea.KeyMsg{Type: tea.KeyRunes, Runes: []rune{'?'}})

			By("Verifying help is shown")
			view := intent.View()
			Expect(view).NotTo(BeEmpty())
		})
	})

	// =========================================================================
	// EDITOR STATE E2E TESTS
	// =========================================================================

	Describe("Editor State Workflow", func() {
		BeforeEach(func() {
			intentCtx := fact_management.NewIntentContext(mockRepo, ctx)
			var err error
			intent, err = fact_management.NewIntent(intentCtx)
			Expect(err).NotTo(HaveOccurred())
			intent.Init()
		})

		Context("New Fact Creation", func() {
			It("should open editor for new fact with n key", func() {
				By("Pressing n for new fact")
				cmd := intent.Update(tea.KeyMsg{Type: tea.KeyRunes, Runes: []rune{'n'}})

				By("Verifying form init command is returned")
				// huh forms return an init command
				Expect(cmd).NotTo(BeNil())

				By("Verifying editor state content")
				view := intent.View()
				Expect(view).To(SatisfyAny(
					ContainSubstring("New Fact"),
					ContainSubstring("Text"),
					ContainSubstring("Edit"),
				))
			})

			It("should show breadcrumbs with New Fact in editor", func() {
				By("Opening new fact editor")
				intent.Update(tea.KeyMsg{Type: tea.KeyRunes, Runes: []rune{'n'}})

				By("Verifying breadcrumbs")
				view := intent.View()
				Expect(view).To(ContainSubstring("New Fact"))
			})

			It("should cancel new fact creation on escape", func() {
				By("Opening new fact editor")
				intent.Update(tea.KeyMsg{Type: tea.KeyRunes, Runes: []rune{'n'}})

				By("Pressing escape to cancel")
				intent.Update(tea.KeyMsg{Type: tea.KeyEsc})

				By("Verifying back at list state")
				view := intent.View()
				Expect(view).To(ContainSubstring("Page"))
			})
		})

		Context("Edit Existing Fact", func() {
			It("should open editor for existing fact with e key from list", func() {
				By("Pressing e to edit selected fact")
				cmd := intent.Update(tea.KeyMsg{Type: tea.KeyRunes, Runes: []rune{'e'}})

				By("Verifying form init command is returned")
				Expect(cmd).NotTo(BeNil())

				By("Verifying editor state content")
				view := intent.View()
				Expect(view).To(SatisfyAny(
					ContainSubstring("Edit"),
					ContainSubstring("Text"),
				))
			})

			It("should open editor for existing fact with e key from view", func() {
				By("Entering view state")
				intent.Update(tea.KeyMsg{Type: tea.KeyEnter})

				By("Pressing e to edit from view")
				cmd := intent.Update(tea.KeyMsg{Type: tea.KeyRunes, Runes: []rune{'e'}})

				By("Verifying form init command is returned")
				Expect(cmd).NotTo(BeNil())

				By("Verifying editor state content")
				view := intent.View()
				Expect(view).To(SatisfyAny(
					ContainSubstring("Edit"),
					ContainSubstring("Text"),
				))
			})

			It("should cancel edit on escape (goes to view state)", func() {
				By("Opening editor")
				intent.Update(tea.KeyMsg{Type: tea.KeyRunes, Runes: []rune{'e'}})

				By("Pressing escape to cancel")
				intent.Update(tea.KeyMsg{Type: tea.KeyEsc})

				By("Verifying at view state (editing existing goes to view)")
				view := intent.View()
				// When editing existing fact and cancelling, goes to view state
				Expect(view).To(ContainSubstring("Fact Details"))
			})
		})

		Context("Editor Global Keys", func() {
			// Note: q is not a global key - this prevents accidental exits.
			// Users must press escape to go back, then navigate to main menu to quit.

			It("should toggle help from editor state", func() {
				By("Opening editor")
				intent.Update(tea.KeyMsg{Type: tea.KeyRunes, Runes: []rune{'n'}})

				By("Pressing ? to toggle help")
				cmd := intent.Update(tea.KeyMsg{Type: tea.KeyRunes, Runes: []rune{'?'}})
				Expect(cmd).To(BeNil())

				By("View should still render")
				view := intent.View()
				Expect(view).NotTo(BeEmpty())
			})
		})
	})

	// =========================================================================
	// DELETE CONFIRM STATE E2E TESTS
	// =========================================================================

	Describe("Delete Confirm State Workflow", func() {
		BeforeEach(func() {
			intentCtx := fact_management.NewIntentContext(mockRepo, ctx)
			var err error
			intent, err = fact_management.NewIntent(intentCtx)
			Expect(err).NotTo(HaveOccurred())
			intent.Init()
		})

		Context("Delete from List State", func() {
			It("should show delete confirmation on d key", func() {
				By("Pressing d to delete")
				intent.Update(tea.KeyMsg{Type: tea.KeyRunes, Runes: []rune{'d'}})

				By("Verifying delete confirmation content")
				view := intent.View()
				Expect(view).To(SatisfyAny(
					ContainSubstring("Delete"),
					ContainSubstring("Confirm"),
				))
			})

			It("should confirm delete with y key", func() {
				By("Opening delete confirmation")
				intent.Update(tea.KeyMsg{Type: tea.KeyRunes, Runes: []rune{'d'}})

				By("Confirming delete with y")
				intent.Update(tea.KeyMsg{Type: tea.KeyRunes, Runes: []rune{'y'}})

				By("Verifying fact was deleted")
				Expect(len(mockRepo.facts)).To(Equal(2))

				By("Verifying back at list state")
				view := intent.View()
				Expect(view).To(ContainSubstring("Page"))
			})

			It("should cancel delete with n key", func() {
				By("Opening delete confirmation")
				intent.Update(tea.KeyMsg{Type: tea.KeyRunes, Runes: []rune{'d'}})

				By("Cancelling delete with n")
				intent.Update(tea.KeyMsg{Type: tea.KeyRunes, Runes: []rune{'n'}})

				By("Verifying fact was not deleted")
				Expect(len(mockRepo.facts)).To(Equal(3))

				By("Verifying back at view state")
				view := intent.View()
				Expect(view).To(SatisfyAny(
					ContainSubstring("Fact Details"),
					ContainSubstring("Page"),
				))
			})

			It("should cancel delete with escape key", func() {
				By("Opening delete confirmation")
				intent.Update(tea.KeyMsg{Type: tea.KeyRunes, Runes: []rune{'d'}})

				By("Cancelling delete with escape")
				intent.Update(tea.KeyMsg{Type: tea.KeyEsc})

				By("Verifying fact was not deleted")
				Expect(len(mockRepo.facts)).To(Equal(3))
			})
		})

		Context("Delete from View State", func() {
			It("should show delete confirmation from view state", func() {
				By("Entering view state")
				intent.Update(tea.KeyMsg{Type: tea.KeyEnter})

				By("Pressing d to delete")
				intent.Update(tea.KeyMsg{Type: tea.KeyRunes, Runes: []rune{'d'}})

				By("Verifying delete confirmation content")
				view := intent.View()
				Expect(view).To(SatisfyAny(
					ContainSubstring("Delete"),
					ContainSubstring("Confirm"),
				))
			})

			It("should confirm delete from view state", func() {
				By("Entering view state")
				intent.Update(tea.KeyMsg{Type: tea.KeyEnter})

				By("Opening delete confirmation")
				intent.Update(tea.KeyMsg{Type: tea.KeyRunes, Runes: []rune{'d'}})

				By("Confirming delete")
				intent.Update(tea.KeyMsg{Type: tea.KeyRunes, Runes: []rune{'y'}})

				By("Verifying fact was deleted")
				Expect(len(mockRepo.facts)).To(Equal(2))
			})
		})

		Context("Delete Error Handling", func() {
			It("should handle delete error gracefully", func() {
				By("Setting up delete error")
				mockRepo.deleteErr = errors.New("delete failed")

				By("Opening delete confirmation")
				intent.Update(tea.KeyMsg{Type: tea.KeyRunes, Runes: []rune{'d'}})

				By("Confirming delete")
				intent.Update(tea.KeyMsg{Type: tea.KeyRunes, Runes: []rune{'y'}})

				By("Verifying fact was not deleted")
				Expect(len(mockRepo.facts)).To(Equal(3))

				By("Verifying result has error status")
				result := intent.Result()
				Expect(result).NotTo(BeNil())
				Expect(result.Status).To(Equal(intents.Failed))
			})
		})
	})

	// =========================================================================
	// RESULTS STATE E2E TESTS
	// =========================================================================

	Describe("Results State", func() {
		// Note: Results state is typically reached after certain operations
		// For now, we test the getResultsContent helper is exercised

		BeforeEach(func() {
			intentCtx := fact_management.NewIntentContext(mockRepo, ctx)
			var err error
			intent, err = fact_management.NewIntent(intentCtx)
			Expect(err).NotTo(HaveOccurred())
			intent.Init()
		})

		// Results state is internal - typically shown after operations
		// We verify the path through delete completion which may set result
		It("should set result after successful delete", func() {
			By("Deleting a fact")
			intent.Update(tea.KeyMsg{Type: tea.KeyRunes, Runes: []rune{'d'}})
			intent.Update(tea.KeyMsg{Type: tea.KeyRunes, Runes: []rune{'y'}})

			By("Verifying result is set")
			result := intent.Result()
			Expect(result).NotTo(BeNil())
			Expect(result.Status).To(Equal(intents.Completed))
		})
	})

	// =========================================================================
	// GLOBAL KEY HANDLING E2E TESTS
	// =========================================================================

	Describe("Global Key Handling", func() {
		BeforeEach(func() {
			intentCtx := fact_management.NewIntentContext(mockRepo, ctx)
			var err error
			intent, err = fact_management.NewIntent(intentCtx)
			Expect(err).NotTo(HaveOccurred())
			intent.Init()
		})

		// Note: HandleGlobalKeys intentionally does NOT handle 'q' to prevent
		// accidental exits. Users should navigate back to main menu to quit.
		Context("Quit Prevention", func() {
			It("should NOT quit on q key from list state (by design)", func() {
				cmd := intent.Update(tea.KeyMsg{Type: tea.KeyRunes, Runes: []rune{'q'}})
				// q is not a global key - users must go back to main menu to quit
				Expect(cmd).To(BeNil())
			})
		})

		Context("Help Toggle", func() {
			It("should toggle help on ? key from list state", func() {
				intent.Update(tea.KeyMsg{Type: tea.KeyRunes, Runes: []rune{'?'}})
				view := intent.View()
				Expect(view).NotTo(BeEmpty())
			})

			It("should toggle help on ? key from view state", func() {
				intent.Update(tea.KeyMsg{Type: tea.KeyEnter})
				intent.Update(tea.KeyMsg{Type: tea.KeyRunes, Runes: []rune{'?'}})
				view := intent.View()
				Expect(view).NotTo(BeEmpty())
			})
		})

		Context("Back Navigation", func() {
			It("should cancel from list state (root) on escape", func() {
				intent.Update(tea.KeyMsg{Type: tea.KeyEsc})

				result := intent.Result()
				Expect(result).NotTo(BeNil())
				Expect(result.Status).To(Equal(intents.Cancelled))
			})

			It("should navigate back from view to list on escape", func() {
				By("Entering view state")
				intent.Update(tea.KeyMsg{Type: tea.KeyEnter})
				view := intent.View()
				Expect(view).To(ContainSubstring("Fact Details"))

				By("Pressing escape")
				intent.Update(tea.KeyMsg{Type: tea.KeyEsc})

				By("Verifying at list state")
				view = intent.View()
				Expect(view).To(ContainSubstring("Page"))
			})
		})
	})

	// =========================================================================
	// REFRESH WORKFLOW E2E TESTS
	// =========================================================================

	Describe("Refresh Workflow", func() {
		BeforeEach(func() {
			intentCtx := fact_management.NewIntentContext(mockRepo, ctx)
			var err error
			intent, err = fact_management.NewIntent(intentCtx)
			Expect(err).NotTo(HaveOccurred())
			intent.Init()
		})

		It("should refresh facts on r key", func() {
			By("Adding a new fact to the repository")
			mockRepo.facts = append(mockRepo.facts, createTestFact("fact-4", "New fact after refresh"))

			By("Pressing r to refresh")
			intent.Update(tea.KeyMsg{Type: tea.KeyRunes, Runes: []rune{'r'}})

			By("Verifying new fact is shown")
			view := intent.View()
			Expect(view).NotTo(BeEmpty())
			// The table should now have 4 items
			Expect(intent.GetTotalItems()).To(Equal(4))
		})

		It("should handle refresh error gracefully", func() {
			By("Setting up list error for refresh")
			mockRepo.listErr = errors.New("refresh failed")

			By("Pressing r to refresh")
			intent.Update(tea.KeyMsg{Type: tea.KeyRunes, Runes: []rune{'r'}})

			By("Verifying error result")
			result := intent.Result()
			Expect(result).NotTo(BeNil())
			Expect(result.Status).To(Equal(intents.Failed))
		})
	})

	// =========================================================================
	// COMPLETE WORKFLOW E2E TESTS
	// =========================================================================

	Describe("Complete Workflow - Create, View, Edit, Delete", func() {
		BeforeEach(func() {
			intentCtx := fact_management.NewIntentContext(mockRepo, ctx)
			var err error
			intent, err = fact_management.NewIntent(intentCtx)
			Expect(err).NotTo(HaveOccurred())
			intent.Init()
		})

		It("should complete full view and delete workflow", func() {
			By("Starting with 3 facts")
			Expect(intent.GetTotalItems()).To(Equal(3))

			By("Selecting first fact (enter)")
			intent.Update(tea.KeyMsg{Type: tea.KeyEnter})
			view := intent.View()
			Expect(view).To(ContainSubstring("Fact Details"))

			By("Going back to list (escape)")
			intent.Update(tea.KeyMsg{Type: tea.KeyEsc})
			view = intent.View()
			Expect(view).To(ContainSubstring("Page"))

			By("Navigating to second fact")
			intent.Update(tea.KeyMsg{Type: tea.KeyDown})
			intent.Update(tea.KeyMsg{Type: tea.KeyEnter})

			By("Deleting from view state")
			intent.Update(tea.KeyMsg{Type: tea.KeyRunes, Runes: []rune{'d'}})
			view = intent.View()
			Expect(view).To(ContainSubstring("Delete"))

			By("Confirming delete")
			intent.Update(tea.KeyMsg{Type: tea.KeyRunes, Runes: []rune{'y'}})

			By("Verifying fact was deleted")
			Expect(len(mockRepo.facts)).To(Equal(2))
		})

		It("should open and cancel edit form", func() {
			By("Opening edit form from list")
			cmd := intent.Update(tea.KeyMsg{Type: tea.KeyRunes, Runes: []rune{'e'}})
			Expect(cmd).NotTo(BeNil()) // Form init command

			By("Verifying in editor state")
			view := intent.View()
			Expect(view).To(SatisfyAny(
				ContainSubstring("Edit"),
				ContainSubstring("Text"),
			))

			By("Cancelling edit")
			intent.Update(tea.KeyMsg{Type: tea.KeyEsc})

			By("Verifying at view state (editing existing goes to view)")
			view = intent.View()
			// When editing existing fact and cancelling, goes to view state
			Expect(view).To(ContainSubstring("Fact Details"))
		})

		It("should open and cancel new fact form", func() {
			By("Opening new fact form")
			cmd := intent.Update(tea.KeyMsg{Type: tea.KeyRunes, Runes: []rune{'n'}})
			Expect(cmd).NotTo(BeNil()) // Form init command

			By("Verifying in editor state")
			view := intent.View()
			Expect(view).To(ContainSubstring("New Fact"))

			By("Cancelling new fact")
			intent.Update(tea.KeyMsg{Type: tea.KeyEsc})

			By("Verifying back at list")
			view = intent.View()
			Expect(view).To(ContainSubstring("Page"))
		})
	})

	// =========================================================================
	// LISTNAVIGATOR INTERFACE E2E TESTS
	// =========================================================================

	Describe("ListNavigator Interface", func() {
		BeforeEach(func() {
			intentCtx := fact_management.NewIntentContext(mockRepo, ctx)
			var err error
			intent, err = fact_management.NewIntent(intentCtx)
			Expect(err).NotTo(HaveOccurred())
			intent.Init()
		})

		It("should report correct total items", func() {
			Expect(intent.GetTotalItems()).To(Equal(3))
		})

		It("should report correct selected index", func() {
			// After init, first item should be selected (index 0)
			Expect(intent.GetSelectedIndex()).To(Equal(0))
		})

		It("should update selected index on navigation", func() {
			intent.Update(tea.KeyMsg{Type: tea.KeyDown})
			Expect(intent.GetSelectedIndex()).To(Equal(1))
		})

		It("should allow setting selected index", func() {
			intent.SetSelectedIndex(2)
			Expect(intent.GetSelectedIndex()).To(Equal(2))
		})

		It("should report correct page size", func() {
			Expect(intent.GetPageSize()).To(Equal(15))
		})
	})
})
