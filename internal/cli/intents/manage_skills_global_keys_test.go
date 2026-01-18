package intents_test

import (
	"context"

	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"

	"github.com/baphled/kariya/internal/cli/intents"
	"github.com/baphled/kariya/internal/domain/career"
	careerrepo "github.com/baphled/kariya/internal/repository/career"
	tea "github.com/charmbracelet/bubbletea"
)

var _ = Describe("ManageSkills - Global Key Handlers", func() {
	var (
		intent     *intents.ManageSkillsIntent
		intentCtx  *intents.ManageSkillsContext
		skillRepo  *careerrepo.MemorySkillRepository
		ctx        context.Context
		testSkills []*career.Skill
	)

	BeforeEach(func() {
		ctx = context.Background()

		// Create memory repositories
		skillRepo = careerrepo.NewMemorySkillRepository()

		// Create test skills
		testSkills = []*career.Skill{
			{
				ID:       "skill-1",
				Name:     "Ruby",
				Category: "backend",
				Level:    "advanced",
			},
			{
				ID:       "skill-2",
				Name:     "Go",
				Category: "backend",
				Level:    "intermediate",
			},
		}

		// Add skills to repository
		for _, skill := range testSkills {
			err := skillRepo.Create(ctx, skill)
			Expect(err).NotTo(HaveOccurred())
		}

		// Create intent context
		intentCtx = &intents.ManageSkillsContext{
			SkillRepository: skillRepo,
			Ctx:             ctx,
		}

		// Create intent and initialize
		intent = intents.NewManageSkillsIntent(intentCtx)
		cmd := intent.Init()
		if cmd != nil {
			msg := cmd()
			intent.Update(msg)
		}
	})

	Describe("List State (Root)", func() {
		It("should return tea.Quit when 'q' is pressed", func() {
			cmd := intent.Update(tea.KeyMsg{Type: tea.KeyRunes, Runes: []rune{'q'}})

			// q should return tea.Quit command
			Expect(cmd).NotTo(BeNil())

			// Intent should NOT be cancelled (app is quitting)
			result := intent.Result()
			Expect(result).To(BeNil(), "Intent should not be cancelled when quitting app")
		})

		It("should toggle help when '?' is pressed", func() {
			// Get initial help state
			initialHelpState := intent.IsHelpVisible()

			// Press '?' to toggle help
			intent.Update(tea.KeyMsg{Type: tea.KeyRunes, Runes: []rune{'?'}})

			// Help state should be toggled
			Expect(intent.IsHelpVisible()).To(Equal(!initialHelpState))
		})

		It("should return tea.Quit when Ctrl+C is pressed", func() {
			cmd := intent.Update(tea.KeyMsg{Type: tea.KeyCtrlC})

			// Ctrl+C should return tea.Quit command
			Expect(cmd).NotTo(BeNil())

			// Intent should NOT be cancelled
			result := intent.Result()
			Expect(result).To(BeNil())
		})
	})

	Describe("Detail State", func() {
		BeforeEach(func() {
			// Navigate to detail state
			intent.Update(tea.KeyMsg{Type: tea.KeyEnter})
		})

		It("should return tea.Quit when 'q' is pressed", func() {
			cmd := intent.Update(tea.KeyMsg{Type: tea.KeyRunes, Runes: []rune{'q'}})

			// q should return tea.Quit command
			Expect(cmd).NotTo(BeNil())

			// Intent should NOT be cancelled
			result := intent.Result()
			Expect(result).To(BeNil())
		})

		It("should toggle help when '?' is pressed", func() {
			initialHelpState := intent.IsHelpVisible()
			intent.Update(tea.KeyMsg{Type: tea.KeyRunes, Runes: []rune{'?'}})
			Expect(intent.IsHelpVisible()).To(Equal(!initialHelpState))
		})

		It("should return tea.Quit when Ctrl+C is pressed", func() {
			cmd := intent.Update(tea.KeyMsg{Type: tea.KeyCtrlC})
			Expect(cmd).NotTo(BeNil())
			result := intent.Result()
			Expect(result).To(BeNil())
		})
	})

	Describe("DetailEvents State", func() {
		BeforeEach(func() {
			// Navigate to detail state, then to events
			intent.Update(tea.KeyMsg{Type: tea.KeyEnter})
			cmd := intent.Update(tea.KeyMsg{Type: tea.KeyEnter})
			if cmd != nil {
				msg := cmd()
				intent.Update(msg)
			}
		})

		It("should return tea.Quit when 'q' is pressed", func() {
			cmd := intent.Update(tea.KeyMsg{Type: tea.KeyRunes, Runes: []rune{'q'}})
			Expect(cmd).NotTo(BeNil())
			result := intent.Result()
			Expect(result).To(BeNil())
		})

		It("should toggle help when '?' is pressed", func() {
			initialHelpState := intent.IsHelpVisible()
			intent.Update(tea.KeyMsg{Type: tea.KeyRunes, Runes: []rune{'?'}})
			Expect(intent.IsHelpVisible()).To(Equal(!initialHelpState))
		})

		It("should return tea.Quit when Ctrl+C is pressed", func() {
			cmd := intent.Update(tea.KeyMsg{Type: tea.KeyCtrlC})
			Expect(cmd).NotTo(BeNil())
		})
	})

	Describe("Filter Modal Visible", func() {
		BeforeEach(func() {
			// Open filter modal (state stays SkillsStateList)
			intent.Update(tea.KeyMsg{Type: tea.KeyRunes, Runes: []rune{'f'}})
		})

		It("should return tea.Quit when 'q' is pressed", func() {
			cmd := intent.Update(tea.KeyMsg{Type: tea.KeyRunes, Runes: []rune{'q'}})
			Expect(cmd).NotTo(BeNil())
			result := intent.Result()
			Expect(result).To(BeNil())
		})

		It("should toggle help when '?' is pressed", func() {
			initialHelpState := intent.IsHelpVisible()
			intent.Update(tea.KeyMsg{Type: tea.KeyRunes, Runes: []rune{'?'}})
			Expect(intent.IsHelpVisible()).To(Equal(!initialHelpState))
		})

		It("should return tea.Quit when Ctrl+C is pressed", func() {
			cmd := intent.Update(tea.KeyMsg{Type: tea.KeyCtrlC})
			Expect(cmd).NotTo(BeNil())
		})
	})

	Describe("Sort Modal Visible", func() {
		BeforeEach(func() {
			// Open sort modal (state stays SkillsStateList)
			intent.Update(tea.KeyMsg{Type: tea.KeyRunes, Runes: []rune{'s'}})
		})

		It("should return tea.Quit when 'q' is pressed", func() {
			cmd := intent.Update(tea.KeyMsg{Type: tea.KeyRunes, Runes: []rune{'q'}})
			Expect(cmd).NotTo(BeNil())
			result := intent.Result()
			Expect(result).To(BeNil())
		})

		It("should toggle help when '?' is pressed", func() {
			initialHelpState := intent.IsHelpVisible()
			intent.Update(tea.KeyMsg{Type: tea.KeyRunes, Runes: []rune{'?'}})
			Expect(intent.IsHelpVisible()).To(Equal(!initialHelpState))
		})

		It("should return tea.Quit when Ctrl+C is pressed", func() {
			cmd := intent.Update(tea.KeyMsg{Type: tea.KeyCtrlC})
			Expect(cmd).NotTo(BeNil())
		})
	})

	Describe("Delete State", func() {
		BeforeEach(func() {
			// Navigate to delete state
			intent.Update(tea.KeyMsg{Type: tea.KeyRunes, Runes: []rune{'d'}})
		})

		It("should return tea.Quit when 'q' is pressed", func() {
			cmd := intent.Update(tea.KeyMsg{Type: tea.KeyRunes, Runes: []rune{'q'}})
			Expect(cmd).NotTo(BeNil())
			result := intent.Result()
			Expect(result).To(BeNil())
		})

		It("should toggle help when '?' is pressed", func() {
			initialHelpState := intent.IsHelpVisible()
			intent.Update(tea.KeyMsg{Type: tea.KeyRunes, Runes: []rune{'?'}})
			Expect(intent.IsHelpVisible()).To(Equal(!initialHelpState))
		})

		It("should return tea.Quit when Ctrl+C is pressed", func() {
			cmd := intent.Update(tea.KeyMsg{Type: tea.KeyCtrlC})
			Expect(cmd).NotTo(BeNil())
		})
	})
})
