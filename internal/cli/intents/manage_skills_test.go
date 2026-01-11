package intents_test

import (
	"context"
	"testing"

	"github.com/baphled/kariya/internal/cli/intents"
	domain "github.com/baphled/kariya/internal/domain/career"
	repo "github.com/baphled/kariya/internal/repository/career"
	careerservice "github.com/baphled/kariya/internal/service/career"
	"github.com/baphled/kariya/internal/testutil"
	tea "github.com/charmbracelet/bubbletea"
	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
)

func TestManageSkillsIntent(t *testing.T) {
	RegisterFailHandler(Fail)
	RunSpecs(t, "ManageSkills Intent Suite")
}

var _ = Describe("ManageSkillsIntent", func() {
	var (
		ctx        context.Context
		skillRepo  repo.SkillRepository
		eventRepo  repo.Repository
		service    *careerservice.Service
		intentCtx  *intents.ManageSkillsContext
		intent     *intents.ManageSkillsIntent
		testSkills []*domain.Skill
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
		testSkills = []*domain.Skill{
			{Name: "Ruby", Category: "backend", Level: "advanced"},
			{Name: "Go", Category: "backend", Level: "expert"},
			{Name: "React", Category: "frontend", Level: "intermediate"},
			{Name: "PostgreSQL", Category: "database"},
			{Name: "Docker", Category: "devops", Level: "advanced"},
		}

		for _, skill := range testSkills {
			err := skillRepo.Create(ctx, skill)
			Expect(err).NotTo(HaveOccurred())
		}

		// Create context
		intentCtx = &intents.ManageSkillsContext{
			SkillRepository: skillRepo,
			Service:         service,
		}
	})

	AfterEach(func() {
		if cleanup != nil {
			cleanup()
		}
	})

	Describe("Init", func() {
		It("should initialize and load skills", func() {
			intent = intents.NewManageSkillsIntent(intentCtx)
			cmd := intent.Init(ctx)

			// Should return a command to load skills
			Expect(cmd).NotTo(BeNil())

			// Execute the command
			msg := cmd()
			loadedMsg, ok := msg.(intents.SkillsLoadedMsg)
			Expect(ok).To(BeTrue(), "Expected SkillsLoadedMsg")
			Expect(loadedMsg.Error).NotTo(HaveOccurred())
			Expect(loadedMsg.Skills).To(HaveLen(5))
		})

		It("should start in list state", func() {
			intent = intents.NewManageSkillsIntent(intentCtx)
			intent.Init(ctx)

			Expect(intent.State()).To(Equal(intents.SkillsStateList))
		})
	})

	Describe("List State Navigation", func() {
		BeforeEach(func() {
			intent = intents.NewManageSkillsIntent(intentCtx)
			intent.Init(ctx)

			// Load skills
			cmd := intent.Init(ctx)
			msg := cmd()
			intent.Update(msg)
		})

		It("should navigate down with j key", func() {
			// Initial cursor at 0
			Expect(intent.SelectedIndex()).To(Equal(0))

			// Press j
			intent.Update(tea.KeyMsg{Type: tea.KeyRunes, Runes: []rune{'j'}})
			Expect(intent.SelectedIndex()).To(Equal(1))
		})

		It("should navigate up with k key", func() {
			// Move down first
			intent.Update(tea.KeyMsg{Type: tea.KeyRunes, Runes: []rune{'j'}})
			intent.Update(tea.KeyMsg{Type: tea.KeyRunes, Runes: []rune{'j'}})
			Expect(intent.SelectedIndex()).To(Equal(2))

			// Press k
			intent.Update(tea.KeyMsg{Type: tea.KeyRunes, Runes: []rune{'k'}})
			Expect(intent.SelectedIndex()).To(Equal(1))
		})

		It("should not go below 0 when pressing k", func() {
			Expect(intent.SelectedIndex()).To(Equal(0))

			// Press k
			intent.Update(tea.KeyMsg{Type: tea.KeyRunes, Runes: []rune{'k'}})
			Expect(intent.SelectedIndex()).To(Equal(0))
		})

		It("should not exceed list length when pressing j", func() {
			// Navigate to last item
			for i := 0; i < 10; i++ {
				intent.Update(tea.KeyMsg{Type: tea.KeyRunes, Runes: []rune{'j'}})
			}

			Expect(intent.SelectedIndex()).To(Equal(4)) // Last skill
		})
	})

	Describe("Add Skill Workflow", func() {
		BeforeEach(func() {
			intent = intents.NewManageSkillsIntent(intentCtx)
			intent.Init(ctx)

			// Load skills
			cmd := intent.Init(ctx)
			msg := cmd()
			intent.Update(msg)
		})

		It("should transition to add state when pressing n", func() {
			Expect(intent.State()).To(Equal(intents.SkillsStateList))

			intent.Update(tea.KeyMsg{Type: tea.KeyRunes, Runes: []rune{'n'}})
			Expect(intent.State()).To(Equal(intents.SkillsStateAdd))
		})

		It("should show form in add state", func() {
			intent.Update(tea.KeyMsg{Type: tea.KeyRunes, Runes: []rune{'n'}})

			view := intent.View()
			Expect(view).To(ContainSubstring("Skill Name"))
			Expect(view).To(ContainSubstring("Category"))
		})

		It("should return to list when form is cancelled", func() {
			intent.Update(tea.KeyMsg{Type: tea.KeyRunes, Runes: []rune{'n'}})
			Expect(intent.State()).To(Equal(intents.SkillsStateAdd))

			// Cancel form with Esc
			intent.Update(tea.KeyMsg{Type: tea.KeyEsc})
			Expect(intent.State()).To(Equal(intents.SkillsStateList))
		})

		It("should create skill and return to list when form is completed", func() {
			intent.Update(tea.KeyMsg{Type: tea.KeyRunes, Runes: []rune{'n'}})

			// Simulate form completion
			msg := intents.SkillFormCompleteMsg{
				Skill: &domain.Skill{
					Name:     "Kubernetes",
					Category: "devops",
					Level:    "intermediate",
				},
				Cancelled: false,
			}

			cmd := intent.Update(msg)
			Expect(cmd).NotTo(BeNil())

			// Should transition back to list
			Eventually(func() intents.SkillsState {
				return intent.State()
			}).Should(Equal(intents.SkillsStateList))

			// Verify skill was created
			skills, err := skillRepo.List(ctx, nil)
			Expect(err).NotTo(HaveOccurred())
			Expect(skills).To(HaveLen(6))

			found := false
			for _, s := range skills {
				if s.Name == "Kubernetes" {
					found = true
					break
				}
			}
			Expect(found).To(BeTrue())
		})
	})

	Describe("Edit Skill Workflow", func() {
		BeforeEach(func() {
			intent = intents.NewManageSkillsIntent(intentCtx)
			intent.Init(ctx)

			// Load skills
			cmd := intent.Init(ctx)
			msg := cmd()
			intent.Update(msg)
		})

		It("should transition to edit state when pressing e", func() {
			Expect(intent.State()).To(Equal(intents.SkillsStateList))

			intent.Update(tea.KeyMsg{Type: tea.KeyRunes, Runes: []rune{'e'}})
			Expect(intent.State()).To(Equal(intents.SkillsStateEdit))
		})

		It("should show form pre-populated with selected skill", func() {
			// Select Ruby (first skill)
			intent.Update(tea.KeyMsg{Type: tea.KeyRunes, Runes: []rune{'e'}})

			view := intent.View()
			Expect(view).To(ContainSubstring("Skill Name"))
			Expect(view).To(ContainSubstring("Ruby"))
		})

		It("should return to list when edit is cancelled", func() {
			intent.Update(tea.KeyMsg{Type: tea.KeyRunes, Runes: []rune{'e'}})
			Expect(intent.State()).To(Equal(intents.SkillsStateEdit))

			// Cancel form
			intent.Update(tea.KeyMsg{Type: tea.KeyEsc})
			Expect(intent.State()).To(Equal(intents.SkillsStateList))
		})

		It("should update skill and return to list when edit is completed", func() {
			intent.Update(tea.KeyMsg{Type: tea.KeyRunes, Runes: []rune{'e'}})

			// Simulate form completion with updated data
			updatedSkill := testSkills[0]
			updatedSkill.Level = "expert" // Was "advanced"

			msg := intents.SkillFormCompleteMsg{
				Skill:     updatedSkill,
				Cancelled: false,
			}

			cmd := intent.Update(msg)
			Expect(cmd).NotTo(BeNil())

			// Should transition back to list
			Eventually(func() intents.SkillsState {
				return intent.State()
			}).Should(Equal(intents.SkillsStateList))

			// Verify skill was updated
			skill, err := skillRepo.GetByID(ctx, updatedSkill.ID)
			Expect(err).NotTo(HaveOccurred())
			Expect(skill.Level).To(Equal("expert"))
		})
	})

	Describe("Delete Skill Workflow", func() {
		BeforeEach(func() {
			intent = intents.NewManageSkillsIntent(intentCtx)
			intent.Init(ctx)

			// Load skills
			cmd := intent.Init(ctx)
			msg := cmd()
			intent.Update(msg)
		})

		It("should transition to delete state when pressing d", func() {
			Expect(intent.State()).To(Equal(intents.SkillsStateList))

			intent.Update(tea.KeyMsg{Type: tea.KeyRunes, Runes: []rune{'d'}})
			Expect(intent.State()).To(Equal(intents.SkillsStateDelete))
		})

		It("should show confirmation modal", func() {
			intent.Update(tea.KeyMsg{Type: tea.KeyRunes, Runes: []rune{'d'}})

			view := intent.View()
			Expect(view).To(ContainSubstring("Delete"))
			Expect(view).To(ContainSubstring("Ruby")) // Selected skill name
		})

		It("should return to list when delete is cancelled with n", func() {
			intent.Update(tea.KeyMsg{Type: tea.KeyRunes, Runes: []rune{'d'}})
			Expect(intent.State()).To(Equal(intents.SkillsStateDelete))

			intent.Update(tea.KeyMsg{Type: tea.KeyRunes, Runes: []rune{'n'}})
			Expect(intent.State()).To(Equal(intents.SkillsStateList))
		})

		It("should return to list when delete is cancelled with Esc", func() {
			intent.Update(tea.KeyMsg{Type: tea.KeyRunes, Runes: []rune{'d'}})
			Expect(intent.State()).To(Equal(intents.SkillsStateDelete))

			intent.Update(tea.KeyMsg{Type: tea.KeyEsc})
			Expect(intent.State()).To(Equal(intents.SkillsStateList))
		})

		It("should delete skill and return to list when confirmed with y", func() {
			skillToDelete := testSkills[0]
			intent.Update(tea.KeyMsg{Type: tea.KeyRunes, Runes: []rune{'d'}})

			intent.Update(tea.KeyMsg{Type: tea.KeyRunes, Runes: []rune{'y'}})

			// Should transition back to list
			Eventually(func() intents.SkillsState {
				return intent.State()
			}).Should(Equal(intents.SkillsStateList))

			// Verify skill was deleted
			_, err := skillRepo.GetByID(ctx, skillToDelete.ID)
			Expect(err).To(HaveOccurred())

			skills, err := skillRepo.List(ctx, nil)
			Expect(err).NotTo(HaveOccurred())
			Expect(skills).To(HaveLen(4))
		})
	})

	Describe("List View Rendering", func() {
		BeforeEach(func() {
			intent = intents.NewManageSkillsIntent(intentCtx)
			intent.Init(ctx)

			// Load skills
			cmd := intent.Init(ctx)
			msg := cmd()
			intent.Update(msg)
		})

		It("should group skills by category", func() {
			view := intent.View()

			// Should show category headers
			Expect(view).To(ContainSubstring("backend"))
			Expect(view).To(ContainSubstring("frontend"))
			Expect(view).To(ContainSubstring("database"))
			Expect(view).To(ContainSubstring("devops"))
		})

		It("should show skill names", func() {
			view := intent.View()

			Expect(view).To(ContainSubstring("Ruby"))
			Expect(view).To(ContainSubstring("Go"))
			Expect(view).To(ContainSubstring("React"))
			Expect(view).To(ContainSubstring("PostgreSQL"))
			Expect(view).To(ContainSubstring("Docker"))
		})

		It("should show level when present", func() {
			view := intent.View()

			Expect(view).To(ContainSubstring("advanced"))
			Expect(view).To(ContainSubstring("expert"))
			Expect(view).To(ContainSubstring("intermediate"))
		})

		It("should show years when present", func() {
			// Add a skill with years
			years := 5
			skillWithYears := &domain.Skill{
				Name:      "Python",
				Category:  "backend",
				YearsUsed: &years,
			}
			err := skillRepo.Create(ctx, skillWithYears)
			Expect(err).NotTo(HaveOccurred())

			// Reload intent
			cmd := intent.Init(ctx)
			msg := cmd()
			intent.Update(msg)

			view := intent.View()
			Expect(view).To(ContainSubstring("5"))
		})
	})

	Describe("Empty State", func() {
		It("should show empty state when no skills exist", func() {
			// Create empty context
			emptyCtx := &intents.ManageSkillsContext{
				SkillRepository: skillRepo,
				Service:         service,
			}

			// Delete all skills
			skills, _ := skillRepo.List(ctx, nil)
			for _, skill := range skills {
				_ = skillRepo.Delete(ctx, skill.ID)
			}

			intent = intents.NewManageSkillsIntent(emptyCtx)
			cmd := intent.Init(ctx)
			msg := cmd()
			intent.Update(msg)

			view := intent.View()
			Expect(view).To(ContainSubstring("No skills"))
			Expect(view).To(ContainSubstring("Press 'n' to add"))
		})
	})

	Describe("Escape to Complete", func() {
		BeforeEach(func() {
			intent = intents.NewManageSkillsIntent(intentCtx)
			intent.Init(ctx)

			// Load skills
			cmd := intent.Init(ctx)
			msg := cmd()
			intent.Update(msg)
		})

		It("should complete intent when Esc is pressed in list state", func() {
			Expect(intent.State()).To(Equal(intents.SkillsStateList))

			intent.Update(tea.KeyMsg{Type: tea.KeyEsc})

			result := intent.Result()
			Expect(result).NotTo(BeNil())
			Expect(result.Status).To(Equal(intents.Cancelled))
		})
	})

	Describe("View Integration", func() {
		BeforeEach(func() {
			intent = intents.NewManageSkillsIntent(intentCtx)
			intent.Init(ctx)

			// Load skills
			cmd := intent.Init(ctx)
			msg := cmd()
			intent.Update(msg)
		})

		It("should use StandardView with logo", func() {
			view := intent.View()
			Expect(view).To(ContainSubstring("KaRiya"))
		})

		It("should show breadcrumbs", func() {
			view := intent.View()
			Expect(view).To(ContainSubstring("Skills"))
		})

		It("should show help footer", func() {
			view := intent.View()
			Expect(view).To(ContainSubstring("n:add"))
			Expect(view).To(ContainSubstring("e:edit"))
			Expect(view).To(ContainSubstring("d:delete"))
			Expect(view).To(ContainSubstring("Esc:back"))
		})
	})
})
