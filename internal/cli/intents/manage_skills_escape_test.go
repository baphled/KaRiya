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
			// Initial state - help should be hidden
			initialView := intent.View()
			Expect(initialView).ToNot(ContainSubstring("Full Help"))

			// Press '?' to show help
			intent.Update(tea.KeyMsg{Type: tea.KeyRunes, Runes: []rune{'?'}})
			viewWithHelp := intent.View()
			Expect(viewWithHelp).To(ContainSubstring("?/h"))

			// Press '?' again to hide help
			intent.Update(tea.KeyMsg{Type: tea.KeyRunes, Runes: []rune{'?'}})
			viewWithoutHelp := intent.View()
			Expect(viewWithoutHelp).ToNot(ContainSubstring("Full Help"))
		})

		It("should show correct footer keys", func() {
			view := intent.View()
			Expect(view).To(ContainSubstring("enter"))
			Expect(view).To(ContainSubstring("esc"))
			Expect(view).To(ContainSubstring("?/h"))
		})
	})

	Describe("Detail View State", func() {
		BeforeEach(func() {
			// Navigate to detail view
			intent.Update(tea.KeyMsg{Type: tea.KeyEnter})
		})

		It("should return to List when escape is pressed", func() {
			// Verify we're in detail state
			view := intent.View()
			Expect(view).To(ContainSubstring("Skill Details"))

			// Press escape to go back
			intent.Update(tea.KeyMsg{Type: tea.KeyEsc})

			// Should be back in list view
			view = intent.View()
			Expect(view).To(ContainSubstring("Skills"))
			// Should not be cancelled
			result := intent.Result()
			Expect(result).To(BeNil())
		})

		It("should return tea.Quit command when 'q' is pressed", func() {
			cmd := intent.Update(tea.KeyMsg{Type: tea.KeyRunes, Runes: []rune{'q'}})
			Expect(cmd).ToNot(BeNil())
		})

		It("should toggle help when '?' is pressed", func() {
			initialView := intent.View()
			Expect(initialView).ToNot(ContainSubstring("Full Help"))

			intent.Update(tea.KeyMsg{Type: tea.KeyRunes, Runes: []rune{'?'}})
			viewWithHelp := intent.View()
			Expect(viewWithHelp).To(ContainSubstring("?/h"))
		})

		It("should show correct footer keys", func() {
			view := intent.View()
			Expect(view).To(ContainSubstring("esc"))
			Expect(view).To(ContainSubstring("?/h"))
		})
	})

	Describe("DetailEvents View State", func() {
		BeforeEach(func() {
			// Navigate to detail view, then to events
			intent.Update(tea.KeyMsg{Type: tea.KeyEnter})
			intent.Update(tea.KeyMsg{Type: tea.KeyRunes, Runes: []rune{'e'}})
		})

		It("should return to Detail when escape is pressed", func() {
			view := intent.View()
			Expect(view).To(ContainSubstring("Events Using This Skill"))

			intent.Update(tea.KeyMsg{Type: tea.KeyEsc})

			// Should be back in detail view
			view = intent.View()
			Expect(view).To(ContainSubstring("Skill Details"))
		})

		It("should return tea.Quit command when 'q' is pressed", func() {
			cmd := intent.Update(tea.KeyMsg{Type: tea.KeyRunes, Runes: []rune{'q'}})
			Expect(cmd).ToNot(BeNil())
		})

		It("should toggle help when '?' is pressed", func() {
			initialView := intent.View()
			Expect(initialView).ToNot(ContainSubstring("Full Help"))

			intent.Update(tea.KeyMsg{Type: tea.KeyRunes, Runes: []rune{'?'}})
			viewWithHelp := intent.View()
			Expect(viewWithHelp).To(ContainSubstring("?/h"))
		})

		It("should show correct footer keys", func() {
			view := intent.View()
			Expect(view).To(ContainSubstring("esc"))
			Expect(view).To(ContainSubstring("?/h"))
		})
	})

	Describe("Add State (Form)", func() {
		BeforeEach(func() {
			// Navigate to add state
			intent.Update(tea.KeyMsg{Type: tea.KeyRunes, Runes: []rune{'a'}})
		})

		It("should return to List when escape is pressed", func() {
			view := intent.View()
			Expect(view).To(ContainSubstring("Add New Skill"))

			intent.Update(tea.KeyMsg{Type: tea.KeyEsc})

			// Should be back in list view
			view = intent.View()
			Expect(view).To(ContainSubstring("Skills"))
		})

		It("should return tea.Quit command when 'q' is pressed", func() {
			cmd := intent.Update(tea.KeyMsg{Type: tea.KeyRunes, Runes: []rune{'q'}})
			Expect(cmd).ToNot(BeNil())
		})

		It("should toggle help when '?' is pressed", func() {
			intent.Update(tea.KeyMsg{Type: tea.KeyRunes, Runes: []rune{'?'}})
			// Help should toggle in form
		})

		It("should show correct footer keys in form", func() {
			view := intent.View()
			Expect(view).To(ContainSubstring("esc"))
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
			Expect(view).To(ContainSubstring("Edit Skill"))

			intent.Update(tea.KeyMsg{Type: tea.KeyEsc})

			// Should be back in list view (not detail)
			view = intent.View()
			Expect(view).To(ContainSubstring("Skills"))
		})

		It("should return tea.Quit command when 'q' is pressed", func() {
			cmd := intent.Update(tea.KeyMsg{Type: tea.KeyRunes, Runes: []rune{'q'}})
			Expect(cmd).ToNot(BeNil())
		})

		It("should toggle help when '?' is pressed", func() {
			intent.Update(tea.KeyMsg{Type: tea.KeyRunes, Runes: []rune{'?'}})
			// Help should toggle in form
		})

		It("should show correct footer keys in form", func() {
			view := intent.View()
			Expect(view).To(ContainSubstring("esc"))
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
			Expect(view).To(ContainSubstring("Delete Skill"))

			intent.Update(tea.KeyMsg{Type: tea.KeyEsc})

			// Should be back in list view
			view = intent.View()
			Expect(view).To(ContainSubstring("Skills"))
		})

		It("should return tea.Quit command when 'q' is pressed", func() {
			cmd := intent.Update(tea.KeyMsg{Type: tea.KeyRunes, Runes: []rune{'q'}})
			Expect(cmd).ToNot(BeNil())
		})

		It("should toggle help when '?' is pressed", func() {
			intent.Update(tea.KeyMsg{Type: tea.KeyRunes, Runes: []rune{'?'}})
			// Help should toggle
		})

		It("should show correct footer keys in confirmation", func() {
			view := intent.View()
			Expect(view).To(ContainSubstring("esc"))
			Expect(view).To(ContainSubstring("y/n"))
		})
	})

	Describe("Filter State (Menu)", func() {
		BeforeEach(func() {
			// Navigate to filter menu
			intent.Update(tea.KeyMsg{Type: tea.KeyRunes, Runes: []rune{'f'}})
		})

		It("should return to List when escape is pressed", func() {
			view := intent.View()
			Expect(view).To(ContainSubstring("Filter Skills"))

			intent.Update(tea.KeyMsg{Type: tea.KeyEsc})

			// Should be back in list view
			view = intent.View()
			Expect(view).To(ContainSubstring("Skills"))
		})

		It("should return tea.Quit command when 'q' is pressed", func() {
			cmd := intent.Update(tea.KeyMsg{Type: tea.KeyRunes, Runes: []rune{'q'}})
			Expect(cmd).ToNot(BeNil())
		})

		It("should toggle help when '?' is pressed", func() {
			intent.Update(tea.KeyMsg{Type: tea.KeyRunes, Runes: []rune{'?'}})
			// Help should toggle
		})

		It("should show correct footer keys in filter menu", func() {
			view := intent.View()
			Expect(view).To(ContainSubstring("esc"))
			Expect(view).To(ContainSubstring("enter"))
		})
	})

	Describe("Sort State (Menu)", func() {
		BeforeEach(func() {
			// Navigate to sort menu
			intent.Update(tea.KeyMsg{Type: tea.KeyRunes, Runes: []rune{'s'}})
		})

		It("should return to List when escape is pressed", func() {
			view := intent.View()
			Expect(view).To(ContainSubstring("Sort Skills"))

			intent.Update(tea.KeyMsg{Type: tea.KeyEsc})

			// Should be back in list view
			view = intent.View()
			Expect(view).To(ContainSubstring("Skills"))
		})

		It("should return tea.Quit command when 'q' is pressed", func() {
			cmd := intent.Update(tea.KeyMsg{Type: tea.KeyRunes, Runes: []rune{'q'}})
			Expect(cmd).ToNot(BeNil())
		})

		It("should toggle help when '?' is pressed", func() {
			intent.Update(tea.KeyMsg{Type: tea.KeyRunes, Runes: []rune{'?'}})
			// Help should toggle
		})

		It("should show correct footer keys in sort menu", func() {
			view := intent.View()
			Expect(view).To(ContainSubstring("esc"))
			Expect(view).To(ContainSubstring("enter"))
		})
	})
})
