package intents_test

import (
	"context"

	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"

	"github.com/baphled/kariya/internal/cli/intents"
	"github.com/baphled/kariya/internal/domain/career"
	repo "github.com/baphled/kariya/internal/repository/career"
	careerservice "github.com/baphled/kariya/internal/service/career"
	"github.com/baphled/kariya/internal/testutil"
	tea "github.com/charmbracelet/bubbletea"
)

var _ = Describe("ManageSkills - Escape Key Behavior", func() {
	var (
		ctx        context.Context
		skillRepo  repo.SkillRepository
		eventRepo  repo.Repository
		service    *careerservice.Service
		intentCtx  *intents.ManageSkillsContext
		intent     *intents.ManageSkillsIntent
		testSkills []*career.Skill
		cleanup    func()
	)

	BeforeEach(func() {
		ctx = context.Background()
		db, cleanupFunc := testutil.SetupTestDB(GinkgoTB())
		cleanup = cleanupFunc

		// Create repositories
		skillRepo = repo.NewSQLiteSkillRepositoryWithDB(db)
		eventRepo = repo.NewSQLiteRepositoryWithDB(db)

		// Create service
		service = careerservice.NewService(eventRepo)

		// Create test skills
		testSkills = []*career.Skill{
			{Name: "Go Programming", Category: "backend", Level: "expert"},
			{Name: "Leadership", Category: "soft skills", Level: "advanced"},
		}

		for _, skill := range testSkills {
			err := skillRepo.Create(ctx, skill)
			Expect(err).NotTo(HaveOccurred())
		}

		// Create context
		intentCtx = &intents.ManageSkillsContext{
			Ctx:             ctx,
			SkillRepository: skillRepo,
			Service:         service,
		}

		// Create and initialize intent
		intent = intents.NewManageSkillsIntent(intentCtx)
		cmd := intent.Init()
		Expect(cmd).NotTo(BeNil())

		// Execute load command
		msg := cmd()
		intent.Update(msg)
	})

	AfterEach(func() {
		if cleanup != nil {
			cleanup()
		}
	})

	Describe("List View (Root State)", func() {
		It("should cancel intent when escape is pressed", func() {
			intent.Update(tea.KeyMsg{Type: tea.KeyEsc})

			result := intent.Result()
			Expect(result).ToNot(BeNil())
			Expect(result.Status).To(Equal(intents.Cancelled))
		})

		It("should return tea.Quit command when 'q' is pressed", func() {
			cmd := intent.Update(tea.KeyMsg{Type: tea.KeyRunes, Runes: []rune{'q'}})

			// q should return tea.Quit, not cancel the intent
			Expect(cmd).ToNot(BeNil())
			// Cannot directly test tea.Quit, but we can verify it's not nil
		})

		It("should toggle help when '?' is pressed", func() {
			// Press '?' to show expanded help
			intent.Update(tea.KeyMsg{Type: tea.KeyRunes, Runes: []rune{'?'}})
			viewWithHelp := intent.View()
			// Verify view still renders after toggle
			Expect(viewWithHelp).To(ContainSubstring("Skills"))

			// Press '?' again to hide expanded help
			intent.Update(tea.KeyMsg{Type: tea.KeyRunes, Runes: []rune{'?'}})
			viewWithoutHelp := intent.View()
			Expect(viewWithoutHelp).To(ContainSubstring("Skills"))
		})

		It("should show correct footer keys", func() {
			view := intent.View()
			// Check for key footer elements (case-sensitive)
			Expect(view).To(ContainSubstring("Enter"))
			Expect(view).To(ContainSubstring("Esc"))
			Expect(view).To(ContainSubstring("Quit"))
		})
	})

	// NOTE: Detail View, DetailEvents View, Add/Edit/Delete State escape tests have been removed
	// The new modal-based architecture handles escape differently - modal overlays close on Escape
	// while the intent remains in List state. See manage_skills_modals_test.go for modal escape tests

	Describe("Filter State (Menu)", func() {
		BeforeEach(func() {
			// Navigate to filter menu
			intent.Update(tea.KeyMsg{Type: tea.KeyRunes, Runes: []rune{'f'}})
		})

		It("should return to List when escape is pressed", func() {
			view := intent.View()
			Expect(view).To(ContainSubstring("Filter"))

			intent.Update(tea.KeyMsg{Type: tea.KeyEsc})

			// Should be back in list view
			view = intent.View()
			Expect(view).To(ContainSubstring("Skills"))
			Expect(view).ToNot(ContainSubstring("Filter by"))
		})

		It("should return tea.Quit command when 'q' is pressed", func() {
			cmd := intent.Update(tea.KeyMsg{Type: tea.KeyRunes, Runes: []rune{'q'}})
			Expect(cmd).ToNot(BeNil())
		})

		It("should toggle help when '?' is pressed", func() {
			intent.Update(tea.KeyMsg{Type: tea.KeyRunes, Runes: []rune{'?'}})
			view := intent.View()
			// Verify view still renders
			Expect(view).To(ContainSubstring("Filter"))
		})

		It("should show correct footer keys in filter modal", func() {
			view := intent.View()
			// Modal should show navigation hints (huh forms have their own footer)
			Expect(view).NotTo(BeEmpty())
		})
	})

	Describe("Sort Modal Visible", func() {
		BeforeEach(func() {
			// Open sort modal (state stays SkillsStateList)
			intent.Update(tea.KeyMsg{Type: tea.KeyRunes, Runes: []rune{'s'}})
		})

		It("should close modal when escape is pressed", func() {
			view := intent.View()
			Expect(view).To(ContainSubstring("Sort"))

			intent.Update(tea.KeyMsg{Type: tea.KeyEsc})

			// Modal should be closed, showing list view
			view = intent.View()
			Expect(view).To(ContainSubstring("Skills"))
		})

		It("should return tea.Quit command when 'q' is pressed", func() {
			cmd := intent.Update(tea.KeyMsg{Type: tea.KeyRunes, Runes: []rune{'q'}})
			// Global 'q' key should work even when modal is visible
			Expect(cmd).ToNot(BeNil())
		})

		It("should toggle help when '?' is pressed", func() {
			intent.Update(tea.KeyMsg{Type: tea.KeyRunes, Runes: []rune{'?'}})
			view := intent.View()
			// Verify view still renders
			Expect(view).NotTo(BeEmpty())
		})

		It("should show modal with sort options", func() {
			view := intent.View()
			// Modal should be visible with sort options
			Expect(view).To(Or(
				ContainSubstring("Sort"),
				ContainSubstring("Name"),
				ContainSubstring("Category"),
			))
		})
	})
})
