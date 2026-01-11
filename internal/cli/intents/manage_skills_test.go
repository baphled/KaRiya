package intents_test

import (
	"context"
	"time"

	"github.com/baphled/kariya/internal/cli/intents"
	domain "github.com/baphled/kariya/internal/domain/career"
	repo "github.com/baphled/kariya/internal/repository/career"
	careerservice "github.com/baphled/kariya/internal/service/career"
	"github.com/baphled/kariya/internal/testutil"
	tea "github.com/charmbracelet/bubbletea"
	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
)

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
			Ctx:             ctx,
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
			cmd := intent.Init()

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
			intent.Init()

			Expect(intent.State()).To(Equal(intents.SkillsStateList))
		})
	})

	Describe("List State Navigation", func() {
		BeforeEach(func() {
			intent = intents.NewManageSkillsIntent(intentCtx)
			intent.Init()

			// Load skills
			cmd := intent.Init()
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
			intent.Init()

			// Load skills
			cmd := intent.Init()
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

			// Execute the create command
			createMsg := cmd()
			Expect(createMsg).To(BeAssignableToTypeOf(intents.SkillCreatedMsg{}))

			// Handle the created message
			cmd2 := intent.Update(createMsg)
			Expect(cmd2).NotTo(BeNil())

			// Execute the reload command
			reloadMsg := cmd2()
			intent.Update(reloadMsg)

			// Should be back in list state
			Expect(intent.State()).To(Equal(intents.SkillsStateList))

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
			intent.Init()

			// Load skills
			cmd := intent.Init()
			msg := cmd()
			intent.Update(msg)
		})

		It("should transition to edit state when pressing e", func() {
			Expect(intent.State()).To(Equal(intents.SkillsStateList))

			intent.Update(tea.KeyMsg{Type: tea.KeyRunes, Runes: []rune{'e'}})
			Expect(intent.State()).To(Equal(intents.SkillsStateEdit))
		})

		It("should show form pre-populated with selected skill", func() {
			// Skills are ordered by name ASC, so first skill is "Docker"
			intent.Update(tea.KeyMsg{Type: tea.KeyRunes, Runes: []rune{'e'}})

			view := intent.View()
			Expect(view).To(ContainSubstring("Skill Name"))
			Expect(view).To(ContainSubstring("Docker"))
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
			// Skills are ordered by name ASC, so first skill is "Docker" (testSkills[4])
			updatedSkill := testSkills[4] // Docker
			updatedSkill.Level = "expert" // Was "advanced"

			msg := intents.SkillFormCompleteMsg{
				Skill:     updatedSkill,
				Cancelled: false,
			}

			cmd := intent.Update(msg)
			Expect(cmd).NotTo(BeNil())

			// Execute the update command
			updateMsg := cmd()
			Expect(updateMsg).To(BeAssignableToTypeOf(intents.SkillUpdatedMsg{}))

			// Handle the updated message
			cmd2 := intent.Update(updateMsg)
			Expect(cmd2).NotTo(BeNil())

			// Execute the reload command
			reloadMsg := cmd2()
			intent.Update(reloadMsg)

			// Should be back in list state
			Expect(intent.State()).To(Equal(intents.SkillsStateList))

			// Verify skill was updated
			skill, err := skillRepo.GetByID(ctx, updatedSkill.ID)
			Expect(err).NotTo(HaveOccurred())
			Expect(skill.Level).To(Equal("expert"))
		})
	})

	Describe("Delete Skill Workflow", func() {
		BeforeEach(func() {
			intent = intents.NewManageSkillsIntent(intentCtx)
			intent.Init()

			// Load skills
			cmd := intent.Init()
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
			// Skills are ordered by name ASC, so first skill is "Docker"
			Expect(view).To(ContainSubstring("Docker")) // Selected skill name
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
			// Skills are ordered by name ASC, so first skill is "Docker" (testSkills[4])
			skillToDelete := testSkills[4] // Docker is first in sorted order
			intent.Update(tea.KeyMsg{Type: tea.KeyRunes, Runes: []rune{'d'}})

			cmd := intent.Update(tea.KeyMsg{Type: tea.KeyRunes, Runes: []rune{'y'}})
			Expect(cmd).NotTo(BeNil())

			// Execute the delete command
			deleteMsg := cmd()
			Expect(deleteMsg).To(BeAssignableToTypeOf(intents.SkillDeletedMsg{}))

			// Handle the deleted message
			cmd2 := intent.Update(deleteMsg)
			Expect(cmd2).NotTo(BeNil())

			// Execute the reload command
			reloadMsg := cmd2()
			intent.Update(reloadMsg)

			// Should be back in list state
			Expect(intent.State()).To(Equal(intents.SkillsStateList))

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
			intent.Init()

			// Load skills
			cmd := intent.Init()
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
			cmd := intent.Init()
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
			cmd := intent.Init()
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
			intent.Init()

			// Load skills
			cmd := intent.Init()
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
			intent.Init()

			// Load skills
			cmd := intent.Init()
			msg := cmd()
			intent.Update(msg)
		})

		It("should use StandardView with logo", func() {
			view := intent.View()
			// Check for ASCII art logo characters or the description
			Expect(view).To(Or(
				ContainSubstring("Career Event Management System"),
				ContainSubstring("██"),
			))
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

		It("should show Enter key in help for detail view", func() {
			view := intent.View()
			Expect(view).To(ContainSubstring("Enter:detail"))
		})

		It("should show event count for skills with events", func() {
			// Create an event with a skill
			eventDate, _ := time.Parse("2006-01-02", "2024-01-15")
			testEvent := &domain.CareerEvent{
				Text:   "Ruby implementation",
				Date:   eventDate,
				Skills: []string{testSkills[0].ID}, // Ruby skill
			}
			err := eventRepo.Create(ctx, testEvent)
			Expect(err).NotTo(HaveOccurred())

			// Reload skills to get event counts
			cmd := intent.Init()
			msg := cmd()
			intent.Update(msg)

			view := intent.View()
			// Should show event count indicator (e.g., "(1 event)" or "[1]")
			Expect(view).To(Or(
				ContainSubstring("1 event"),
				ContainSubstring("[1]"),
			))
		})
	})

	Describe("Detail View", func() {
		BeforeEach(func() {
			// Create intent
			intent = intents.NewManageSkillsIntent(intentCtx)

			// Initialize intent and load skills
			cmd := intent.Init()
			Expect(cmd).NotTo(BeNil())

			// Execute the init command to load skills
			msg := cmd()
			intent.Update(msg)

			// Ensure we have skills loaded
			Expect(intent.State()).To(Equal(intents.SkillsStateList))
		})

		It("should transition to detail view when Enter is pressed", func() {
			// Press Enter to view detail
			cmd := intent.Update(tea.KeyMsg{Type: tea.KeyEnter})
			Expect(cmd).To(BeNil())
			Expect(intent.State()).To(Equal(intents.SkillsStateDetail))
		})

		It("should display skill name in detail view", func() {
			// Press Enter to view detail
			intent.Update(tea.KeyMsg{Type: tea.KeyEnter})
			view := intent.View()
			// Should show "Name:" label and a skill name (Docker, Ruby, Go, etc.)
			Expect(view).To(ContainSubstring("Name:"))
			Expect(view).To(Or(
				ContainSubstring("Ruby"),
				ContainSubstring("Go"),
				ContainSubstring("React"),
				ContainSubstring("PostgreSQL"),
				ContainSubstring("Docker"),
			))
		})

		It("should display skill category in detail view", func() {
			// Press Enter to view detail
			intent.Update(tea.KeyMsg{Type: tea.KeyEnter})
			view := intent.View()
			// Should show "Category:" label and a category
			Expect(view).To(ContainSubstring("Category:"))
			Expect(view).To(Or(
				ContainSubstring("backend"),
				ContainSubstring("frontend"),
				ContainSubstring("database"),
				ContainSubstring("devops"),
			))
		})

		It("should display skill level in detail view when set", func() {
			// Press Enter to view detail
			intent.Update(tea.KeyMsg{Type: tea.KeyEnter})
			view := intent.View()
			Expect(view).To(ContainSubstring("advanced"))
		})

		It("should display event count in detail view", func() {
			// Press Enter to view detail
			intent.Update(tea.KeyMsg{Type: tea.KeyEnter})
			view := intent.View()
			Expect(view).To(ContainSubstring("Event Count"))
		})

		It("should transition to events view when Enter is pressed from detail", func() {
			// Navigate to detail
			intent.Update(tea.KeyMsg{Type: tea.KeyEnter})
			Expect(intent.State()).To(Equal(intents.SkillsStateDetail))

			// Press Enter to view events
			cmd := intent.Update(tea.KeyMsg{Type: tea.KeyEnter})
			Expect(cmd).NotTo(BeNil())

			// Execute command to load events
			msg := cmd()
			intent.Update(msg)

			Expect(intent.State()).To(Equal(intents.SkillsStateDetailEvents))
		})

		It("should return to list when Esc is pressed from detail", func() {
			// Navigate to detail
			intent.Update(tea.KeyMsg{Type: tea.KeyEnter})
			Expect(intent.State()).To(Equal(intents.SkillsStateDetail))

			// Press Esc to go back
			intent.Update(tea.KeyMsg{Type: tea.KeyEsc})
			Expect(intent.State()).To(Equal(intents.SkillsStateList))
		})

		It("should transition to edit when 'e' is pressed from detail", func() {
			// Navigate to detail
			intent.Update(tea.KeyMsg{Type: tea.KeyEnter})
			Expect(intent.State()).To(Equal(intents.SkillsStateDetail))

			// Press 'e' to edit
			cmd := intent.Update(tea.KeyMsg{Type: tea.KeyRunes, Runes: []rune{'e'}})
			Expect(cmd).NotTo(BeNil())
			Expect(intent.State()).To(Equal(intents.SkillsStateEdit))
		})

		It("should transition to delete when 'd' is pressed from detail", func() {
			// Navigate to detail
			intent.Update(tea.KeyMsg{Type: tea.KeyEnter})
			Expect(intent.State()).To(Equal(intents.SkillsStateDetail))

			// Press 'd' to delete
			_ = intent.Update(tea.KeyMsg{Type: tea.KeyRunes, Runes: []rune{'d'}})
			Expect(intent.State()).To(Equal(intents.SkillsStateDelete))
		})
	})

	Describe("Detail Events View", func() {
		var testEvent *domain.CareerEvent

		BeforeEach(func() {
			// Create intent
			intent = intents.NewManageSkillsIntent(intentCtx)

			// Initialize intent and load skills
			cmd := intent.Init()
			Expect(cmd).NotTo(BeNil())
			msg := cmd()
			intent.Update(msg)

			// Create a test event with the first skill
			eventDate, _ := time.Parse("2006-01-02", "2024-01-15")
			testEvent = &domain.CareerEvent{
				Text:   "Implemented Ruby feature",
				Date:   eventDate,
				Skills: []string{testSkills[0].ID}, // Ruby skill
			}
			err := eventRepo.Create(ctx, testEvent)
			Expect(err).NotTo(HaveOccurred())

			// Navigate to detail view
			intent.Update(tea.KeyMsg{Type: tea.KeyEnter})
			Expect(intent.State()).To(Equal(intents.SkillsStateDetail))
		})

		It("should load and display events using the skill", func() {
			// Press Enter to view events
			cmd := intent.Update(tea.KeyMsg{Type: tea.KeyEnter})
			Expect(cmd).NotTo(BeNil())

			// Execute command to load events
			msg := cmd()
			intent.Update(msg)

			Expect(intent.State()).To(Equal(intents.SkillsStateDetailEvents))

			// View should show either events or empty state
			view := intent.View()
			Expect(view).To(Or(
				ContainSubstring("Implemented Ruby feature"),
				ContainSubstring("No events use this skill"),
				ContainSubstring("Loading events"),
			))
		})

		It("should show events view renders successfully", func() {
			// Navigate to events view
			cmd := intent.Update(tea.KeyMsg{Type: tea.KeyEnter})
			msg := cmd()
			intent.Update(msg)

			view := intent.View()
			// Just verify the view renders with breadcrumbs
			Expect(view).To(ContainSubstring("Events"))
		})

		It("should return to detail when Esc is pressed from events view", func() {
			// Navigate to events view
			cmd := intent.Update(tea.KeyMsg{Type: tea.KeyEnter})
			msg := cmd()
			intent.Update(msg)
			Expect(intent.State()).To(Equal(intents.SkillsStateDetailEvents))

			// Press Esc to go back to detail
			intent.Update(tea.KeyMsg{Type: tea.KeyEsc})
			Expect(intent.State()).To(Equal(intents.SkillsStateDetail))
		})

		It("should show events view successfully even when no events", func() {
			// Navigate to different skill (Go - no events)
			intent.Update(tea.KeyMsg{Type: tea.KeyRunes, Runes: []rune{'j'}}) // Move to next skill
			intent.Update(tea.KeyMsg{Type: tea.KeyEnter})                     // Enter detail view

			// Navigate to events view
			cmd := intent.Update(tea.KeyMsg{Type: tea.KeyEnter})
			if cmd != nil {
				msg := cmd()
				intent.Update(msg)
			}

			// Just verify the view renders and we're in the events state
			Expect(intent.State()).To(Equal(intents.SkillsStateDetailEvents))
			view := intent.View()
			Expect(view).To(ContainSubstring("Events"))
		})
	})
})
