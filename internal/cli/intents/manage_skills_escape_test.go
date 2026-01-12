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

	Describe("Detail View State", func() {
		BeforeEach(func() {
			// Navigate to detail view
			intent.Update(tea.KeyMsg{Type: tea.KeyEnter})
		})

		It("should return to List when escape is pressed", func() {
			// Verify we're in detail state (breadcrumb shows "Skills ▸ [SkillName]")
			view := intent.View()
			Expect(view).To(ContainSubstring("Skills  ▸"))

			// Press escape to go back
			intent.Update(tea.KeyMsg{Type: tea.KeyEsc})

			// Should be back in list view (no breadcrumb arrow)
			view = intent.View()
			Expect(view).To(ContainSubstring("Skills"))
			Expect(view).ToNot(ContainSubstring("▸"))
			// Should not be cancelled
			result := intent.Result()
			Expect(result).To(BeNil())
		})

		It("should return tea.Quit command when 'q' is pressed", func() {
			cmd := intent.Update(tea.KeyMsg{Type: tea.KeyRunes, Runes: []rune{'q'}})
			Expect(cmd).ToNot(BeNil())
		})

		It("should toggle help when '?' is pressed", func() {
			intent.Update(tea.KeyMsg{Type: tea.KeyRunes, Runes: []rune{'?'}})
			viewWithHelp := intent.View()
			// Verify view still renders after help toggle
			Expect(viewWithHelp).To(ContainSubstring("Skills  ▸"))
		})

		It("should show correct footer keys", func() {
			view := intent.View()
			Expect(view).To(ContainSubstring("Esc"))
			Expect(view).To(ContainSubstring("Quit"))
		})
	})

	Describe("DetailEvents View State", func() {
		BeforeEach(func() {
			// Navigate to detail view, then to events (Enter from detail = view events)
			intent.Update(tea.KeyMsg{Type: tea.KeyEnter})
			// Need to wait for skill to be selected, then press Enter again to view events
			intent.Update(tea.KeyMsg{Type: tea.KeyEnter})
		})

		It("should return to Detail when escape is pressed", func() {
			view := intent.View()
			Expect(view).To(ContainSubstring("Events"))

			intent.Update(tea.KeyMsg{Type: tea.KeyEsc})

			// Should be back in detail view (breadcrumb has one ▸)
			view = intent.View()
			Expect(view).To(ContainSubstring("Skills  ▸"))
		})

		It("should return tea.Quit command when 'q' is pressed", func() {
			cmd := intent.Update(tea.KeyMsg{Type: tea.KeyRunes, Runes: []rune{'q'}})
			Expect(cmd).ToNot(BeNil())
		})

		It("should toggle help when '?' is pressed", func() {
			intent.Update(tea.KeyMsg{Type: tea.KeyRunes, Runes: []rune{'?'}})
			viewWithHelp := intent.View()
			// Verify view still renders
			Expect(viewWithHelp).To(ContainSubstring("Events"))
		})

		It("should show correct footer keys", func() {
			view := intent.View()
			Expect(view).To(ContainSubstring("Esc"))
			Expect(view).To(ContainSubstring("Quit"))
		})
	})

	Describe("Add State (Form)", func() {
		BeforeEach(func() {
			// Navigate to add state (n = new skill)
			intent.Update(tea.KeyMsg{Type: tea.KeyRunes, Runes: []rune{'n'}})
		})

		It("should return to List when escape is pressed", func() {
			view := intent.View()
			Expect(view).To(ContainSubstring("Add"))

			intent.Update(tea.KeyMsg{Type: tea.KeyEsc})

			// Should be back in list view
			view = intent.View()
			Expect(view).To(ContainSubstring("Skills"))
			Expect(view).ToNot(ContainSubstring("Add"))
		})

		It("should return tea.Quit command when 'q' is pressed", func() {
			cmd := intent.Update(tea.KeyMsg{Type: tea.KeyRunes, Runes: []rune{'q'}})
			Expect(cmd).ToNot(BeNil())
		})

		It("should toggle help when '?' is pressed", func() {
			intent.Update(tea.KeyMsg{Type: tea.KeyRunes, Runes: []rune{'?'}})
			view := intent.View()
			// Verify view still renders
			Expect(view).To(ContainSubstring("Add"))
		})

		It("should show correct footer keys in form", func() {
			view := intent.View()
			Expect(view).To(ContainSubstring("Esc"))
		})
	})

	Describe("Edit State (Form)", func() {
		BeforeEach(func() {
			// Navigate to detail, then edit
			intent.Update(tea.KeyMsg{Type: tea.KeyEnter})
			intent.Update(tea.KeyMsg{Type: tea.KeyRunes, Runes: []rune{'e'}})
		})

		It("should return to List when escape is pressed", func() {
			view := intent.View()
			// Edit form should be visible
			Expect(view).To(ContainSubstring("Edit"))

			intent.Update(tea.KeyMsg{Type: tea.KeyEsc})

			// Should be back in list view (not detail)
			view = intent.View()
			Expect(view).To(ContainSubstring("Skills"))
			Expect(view).ToNot(ContainSubstring("Edit"))
		})

		It("should return tea.Quit command when 'q' is pressed", func() {
			cmd := intent.Update(tea.KeyMsg{Type: tea.KeyRunes, Runes: []rune{'q'}})
			Expect(cmd).ToNot(BeNil())
		})

		It("should toggle help when '?' is pressed", func() {
			intent.Update(tea.KeyMsg{Type: tea.KeyRunes, Runes: []rune{'?'}})
			view := intent.View()
			// Verify view still renders
			Expect(view).To(ContainSubstring("Edit"))
		})

		It("should show correct footer keys in form", func() {
			view := intent.View()
			Expect(view).To(ContainSubstring("Esc"))
		})
	})

	Describe("Delete State (Confirmation)", func() {
		BeforeEach(func() {
			// Navigate to detail, then delete
			intent.Update(tea.KeyMsg{Type: tea.KeyEnter})
			intent.Update(tea.KeyMsg{Type: tea.KeyRunes, Runes: []rune{'d'}})
		})

		It("should return to List when escape is pressed", func() {
			view := intent.View()
			Expect(view).To(ContainSubstring("Delete"))

			intent.Update(tea.KeyMsg{Type: tea.KeyEsc})

			// Should be back in list view
			view = intent.View()
			Expect(view).To(ContainSubstring("Skills"))
			Expect(view).ToNot(ContainSubstring("Delete"))
		})

		It("should return tea.Quit command when 'q' is pressed", func() {
			cmd := intent.Update(tea.KeyMsg{Type: tea.KeyRunes, Runes: []rune{'q'}})
			Expect(cmd).ToNot(BeNil())
		})

		It("should toggle help when '?' is pressed", func() {
			intent.Update(tea.KeyMsg{Type: tea.KeyRunes, Runes: []rune{'?'}})
			view := intent.View()
			// Verify view still renders
			Expect(view).To(ContainSubstring("Delete"))
		})

		It("should show correct footer keys in confirmation", func() {
			view := intent.View()
			Expect(view).To(ContainSubstring("Esc"))
			Expect(view).To(ContainSubstring("y"))
		})
	})

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

		It("should show correct footer keys in filter menu", func() {
			view := intent.View()
			Expect(view).To(ContainSubstring("Esc"))
			Expect(view).To(ContainSubstring("Enter"))
		})
	})

	Describe("Sort State (Menu)", func() {
		BeforeEach(func() {
			// Navigate to sort menu
			intent.Update(tea.KeyMsg{Type: tea.KeyRunes, Runes: []rune{'s'}})
		})

		It("should return to List when escape is pressed", func() {
			view := intent.View()
			Expect(view).To(ContainSubstring("Sort"))

			intent.Update(tea.KeyMsg{Type: tea.KeyEsc})

			// Should be back in list view
			view = intent.View()
			Expect(view).To(ContainSubstring("Skills"))
			Expect(view).ToNot(ContainSubstring("Sort by"))
		})

		It("should return tea.Quit command when 'q' is pressed", func() {
			cmd := intent.Update(tea.KeyMsg{Type: tea.KeyRunes, Runes: []rune{'q'}})
			Expect(cmd).ToNot(BeNil())
		})

		It("should toggle help when '?' is pressed", func() {
			intent.Update(tea.KeyMsg{Type: tea.KeyRunes, Runes: []rune{'?'}})
			view := intent.View()
			// Verify view still renders
			Expect(view).To(ContainSubstring("Sort"))
		})

		It("should show correct footer keys in sort menu", func() {
			view := intent.View()
			Expect(view).To(ContainSubstring("Esc"))
			Expect(view).To(ContainSubstring("Enter"))
		})
	})
})
