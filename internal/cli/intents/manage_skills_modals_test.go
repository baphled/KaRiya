package intents_test

import (
	"context"

	"github.com/baphled/kariya/internal/cli/intents"
	domain "github.com/baphled/kariya/internal/domain/career"
	repo "github.com/baphled/kariya/internal/repository/career"
	careerservice "github.com/baphled/kariya/internal/service/career"
	"github.com/baphled/kariya/internal/testutil"
	tea "github.com/charmbracelet/bubbletea"
	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
)

var _ = Describe("ManageSkillsIntent Modal Integration", func() {
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

		// Create test skills with years
		years5 := 5
		years3 := 3
		testSkills = []*domain.Skill{
			{Name: "Ruby", Category: "backend", Level: "advanced", YearsUsed: &years5},
			{Name: "Go", Category: "backend", Level: "expert", YearsUsed: &years3},
			{Name: "React", Category: "frontend", Level: "intermediate"},
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
		if cmd != nil {
			msg := cmd()
			intent.Update(msg)
		}
	})

	AfterEach(func() {
		if cleanup != nil {
			cleanup()
		}
	})

	Describe("ViewSkillDetailModal", func() {
		Context("when pressing Enter on a skill in list view", func() {
			It("should show the view detail modal", func() {
				// Should start in list state
				Expect(intent.State()).To(Equal(intents.SkillsStateList))

				// Press Enter to view detail
				intent.Update(tea.KeyMsg{Type: tea.KeyEnter})

				// Should have a visible detail modal
				Expect(intent.HasVisibleDetailModal()).To(BeTrue())
			})

			It("should display skill information in the modal", func() {
				// Press Enter to view detail
				intent.Update(tea.KeyMsg{Type: tea.KeyEnter})

				// View should contain skill name (Go comes first alphabetically)
				view := intent.View()
				Expect(view).To(ContainSubstring("Go"))
			})

			It("should close modal on Escape and return to list", func() {
				// Open modal
				intent.Update(tea.KeyMsg{Type: tea.KeyEnter})
				Expect(intent.HasVisibleDetailModal()).To(BeTrue())

				// Press Escape
				intent.Update(tea.KeyMsg{Type: tea.KeyEsc})

				// Modal should be closed
				Expect(intent.HasVisibleDetailModal()).To(BeFalse())
				// Should still be in list state
				Expect(intent.State()).To(Equal(intents.SkillsStateList))
			})

			It("should close modal on Enter", func() {
				// Open modal
				intent.Update(tea.KeyMsg{Type: tea.KeyEnter})
				Expect(intent.HasVisibleDetailModal()).To(BeTrue())

				// Press Enter to close
				intent.Update(tea.KeyMsg{Type: tea.KeyEnter})

				// Modal should be closed
				Expect(intent.HasVisibleDetailModal()).To(BeFalse())
			})
		})
	})

	Describe("SkillAddEditModal", func() {
		Context("when pressing 'n' to add new skill", func() {
			It("should show the add/edit modal in add mode", func() {
				// Should start in list state
				Expect(intent.State()).To(Equal(intents.SkillsStateList))

				// Press 'n' to add new skill
				cmd := intent.Update(tea.KeyMsg{Type: tea.KeyRunes, Runes: []rune{'n'}})
				// Execute any returned command (form Init)
				if cmd != nil {
					msg := cmd()
					if msg != nil {
						intent.Update(msg)
					}
				}

				// Should have a visible add/edit modal
				Expect(intent.HasVisibleAddEditModal()).To(BeTrue())
			})

			It("should close modal on Escape without saving", func() {
				// Open add modal
				cmd := intent.Update(tea.KeyMsg{Type: tea.KeyRunes, Runes: []rune{'n'}})
				if cmd != nil {
					msg := cmd()
					if msg != nil {
						intent.Update(msg)
					}
				}
				Expect(intent.HasVisibleAddEditModal()).To(BeTrue())

				// Press Escape
				intent.Update(tea.KeyMsg{Type: tea.KeyEsc})

				// Modal should be closed
				Expect(intent.HasVisibleAddEditModal()).To(BeFalse())
				// Should still be in list state
				Expect(intent.State()).To(Equal(intents.SkillsStateList))
			})
		})

		Context("when pressing 'e' to edit selected skill", func() {
			It("should show the add/edit modal in edit mode", func() {
				// Press 'e' to edit selected skill
				cmd := intent.Update(tea.KeyMsg{Type: tea.KeyRunes, Runes: []rune{'e'}})
				if cmd != nil {
					msg := cmd()
					if msg != nil {
						intent.Update(msg)
					}
				}

				// Should have a visible add/edit modal
				Expect(intent.HasVisibleAddEditModal()).To(BeTrue())
			})

			It("should pre-populate form with skill data", func() {
				// Press 'e' to edit
				cmd := intent.Update(tea.KeyMsg{Type: tea.KeyRunes, Runes: []rune{'e'}})
				if cmd != nil {
					msg := cmd()
					if msg != nil {
						intent.Update(msg)
					}
				}

				// View should contain the skill data being edited (Go is first alphabetically)
				view := intent.View()
				Expect(view).To(ContainSubstring("expert")) // Go's level
			})
		})
	})

	Describe("DeleteConfirmModal", func() {
		Context("when pressing 'd' to delete selected skill", func() {
			It("should show the delete confirmation modal", func() {
				// Press 'd' to delete
				cmd := intent.Update(tea.KeyMsg{Type: tea.KeyRunes, Runes: []rune{'d'}})
				if cmd != nil {
					msg := cmd()
					if msg != nil {
						intent.Update(msg)
					}
				}

				// Should have a visible delete modal
				Expect(intent.HasVisibleDeleteModal()).To(BeTrue())
			})

			It("should close modal on Escape without deleting", func() {
				// Open delete modal
				cmd := intent.Update(tea.KeyMsg{Type: tea.KeyRunes, Runes: []rune{'d'}})
				if cmd != nil {
					msg := cmd()
					if msg != nil {
						intent.Update(msg)
					}
				}
				Expect(intent.HasVisibleDeleteModal()).To(BeTrue())

				// Press Escape
				intent.Update(tea.KeyMsg{Type: tea.KeyEsc})

				// Modal should be closed
				Expect(intent.HasVisibleDeleteModal()).To(BeFalse())
				// Skills should still exist
				Expect(intent.State()).To(Equal(intents.SkillsStateList))
			})

			It("should close modal when 'n' is pressed", func() {
				// Open delete modal
				cmd := intent.Update(tea.KeyMsg{Type: tea.KeyRunes, Runes: []rune{'d'}})
				if cmd != nil {
					msg := cmd()
					if msg != nil {
						intent.Update(msg)
					}
				}
				Expect(intent.HasVisibleDeleteModal()).To(BeTrue())

				// Press 'n' to cancel
				intent.Update(tea.KeyMsg{Type: tea.KeyRunes, Runes: []rune{'n'}})

				// Modal should be closed
				Expect(intent.HasVisibleDeleteModal()).To(BeFalse())
			})

			It("should display skill name in delete confirmation", func() {
				// Press 'd' to delete
				cmd := intent.Update(tea.KeyMsg{Type: tea.KeyRunes, Runes: []rune{'d'}})
				if cmd != nil {
					msg := cmd()
					if msg != nil {
						intent.Update(msg)
					}
				}

				// View should contain delete confirmation with skill name
				view := intent.View()
				Expect(view).To(ContainSubstring("Delete"))
				Expect(view).To(ContainSubstring("Go")) // First skill alphabetically
			})
		})
	})

	Describe("ViewSkillDetailModal Actions", func() {
		Context("when actions are triggered from detail modal", func() {
			It("should open edit modal when 'e' is pressed from detail", func() {
				// First open detail modal
				intent.Update(tea.KeyMsg{Type: tea.KeyEnter})
				Expect(intent.HasVisibleDetailModal()).To(BeTrue())

				// Press 'e' to edit (standard edit keybinding)
				intent.Update(tea.KeyMsg{Type: tea.KeyRunes, Runes: []rune{'e'}})

				// Detail modal should close and edit modal should open
				Expect(intent.HasVisibleDetailModal()).To(BeFalse())
				Expect(intent.HasVisibleAddEditModal()).To(BeTrue())
			})

			It("should show skill events modal when ctrl+e is pressed from detail", func() {
				// First open detail modal
				intent.Update(tea.KeyMsg{Type: tea.KeyEnter})
				Expect(intent.HasVisibleDetailModal()).To(BeTrue())

				// Press ctrl+e to view events using this skill
				cmd := intent.Update(tea.KeyMsg{Type: tea.KeyCtrlE})

				// Detail modal should close
				Expect(intent.HasVisibleDetailModal()).To(BeFalse())

				// Execute the returned command to load events for modal
				if cmd != nil {
					msg := cmd()
					if msg != nil {
						intent.Update(msg)
					}
				}

				// Should now show the skill events modal
				Expect(intent.HasVisibleSkillEventsModal()).To(BeTrue())
			})

			It("should open delete modal when 'd' is pressed from detail", func() {
				// First open detail modal
				intent.Update(tea.KeyMsg{Type: tea.KeyEnter})
				Expect(intent.HasVisibleDetailModal()).To(BeTrue())

				// Press 'd' to delete
				intent.Update(tea.KeyMsg{Type: tea.KeyRunes, Runes: []rune{'d'}})

				// Detail modal should close and delete modal should open
				Expect(intent.HasVisibleDetailModal()).To(BeFalse())
				Expect(intent.HasVisibleDeleteModal()).To(BeTrue())
			})
		})
	})

	Describe("Modal State Preservation", func() {
		It("should remain in list state when modals are shown", func() {
			// Open detail modal
			intent.Update(tea.KeyMsg{Type: tea.KeyEnter})
			Expect(intent.HasVisibleDetailModal()).To(BeTrue())
			Expect(intent.State()).To(Equal(intents.SkillsStateList))

			// Close and open add modal
			intent.Update(tea.KeyMsg{Type: tea.KeyEsc})
			cmd := intent.Update(tea.KeyMsg{Type: tea.KeyRunes, Runes: []rune{'n'}})
			if cmd != nil {
				msg := cmd()
				if msg != nil {
					intent.Update(msg)
				}
			}
			Expect(intent.HasVisibleAddEditModal()).To(BeTrue())
			Expect(intent.State()).To(Equal(intents.SkillsStateList))

			// Close and open delete modal
			intent.Update(tea.KeyMsg{Type: tea.KeyEsc})
			cmd = intent.Update(tea.KeyMsg{Type: tea.KeyRunes, Runes: []rune{'d'}})
			if cmd != nil {
				msg := cmd()
				if msg != nil {
					intent.Update(msg)
				}
			}
			Expect(intent.HasVisibleDeleteModal()).To(BeTrue())
			Expect(intent.State()).To(Equal(intents.SkillsStateList))
		})

		It("should preserve list selection when modal is closed", func() {
			// Navigate down to select second skill
			intent.Update(tea.KeyMsg{Type: tea.KeyRunes, Runes: []rune{'j'}})

			// Open and close detail modal
			intent.Update(tea.KeyMsg{Type: tea.KeyEnter})
			intent.Update(tea.KeyMsg{Type: tea.KeyEsc})

			// Should still have the same selection (view should show same skill highlighted)
			view := intent.View()
			// The view should still have the selection indicator on a skill
			Expect(view).To(ContainSubstring("▶"))
		})
	})

	Describe("Empty List Handling", func() {
		var emptyIntent *intents.ManageSkillsIntent
		var emptyCleanup func()

		BeforeEach(func() {
			// Create a new intent with no skills
			emptyDb, cleanupFunc := testutil.SetupTestDB(GinkgoTB())
			emptyCleanup = cleanupFunc

			emptySkillRepo := repo.NewSQLiteSkillRepositoryWithDB(emptyDb)
			emptyEventRepo := repo.NewSQLiteRepositoryWithDB(emptyDb)
			emptyService := careerservice.NewService(emptyEventRepo)

			emptyIntentCtx := &intents.ManageSkillsContext{
				Ctx:             ctx,
				SkillRepository: emptySkillRepo,
				Service:         emptyService,
			}

			emptyIntent = intents.NewManageSkillsIntent(emptyIntentCtx)
			cmd := emptyIntent.Init()
			if cmd != nil {
				msg := cmd()
				emptyIntent.Update(msg)
			}
		})

		AfterEach(func() {
			if emptyCleanup != nil {
				emptyCleanup()
			}
		})

		It("should not open detail modal when list is empty", func() {
			// Try to open detail modal
			emptyIntent.Update(tea.KeyMsg{Type: tea.KeyEnter})

			// Should not have opened modal
			Expect(emptyIntent.HasVisibleDetailModal()).To(BeFalse())
		})

		It("should still allow adding new skill when list is empty", func() {
			// Open add modal
			cmd := emptyIntent.Update(tea.KeyMsg{Type: tea.KeyRunes, Runes: []rune{'n'}})
			if cmd != nil {
				msg := cmd()
				if msg != nil {
					emptyIntent.Update(msg)
				}
			}

			// Should have opened add modal
			Expect(emptyIntent.HasVisibleAddEditModal()).To(BeTrue())
		})

		It("should not open edit modal when list is empty", func() {
			// Try to edit
			cmd := emptyIntent.Update(tea.KeyMsg{Type: tea.KeyRunes, Runes: []rune{'e'}})
			if cmd != nil {
				msg := cmd()
				if msg != nil {
					emptyIntent.Update(msg)
				}
			}

			// Should not have opened modal
			Expect(emptyIntent.HasVisibleAddEditModal()).To(BeFalse())
		})

		It("should not open delete modal when list is empty", func() {
			// Try to delete
			cmd := emptyIntent.Update(tea.KeyMsg{Type: tea.KeyRunes, Runes: []rune{'d'}})
			if cmd != nil {
				msg := cmd()
				if msg != nil {
					emptyIntent.Update(msg)
				}
			}

			// Should not have opened modal
			Expect(emptyIntent.HasVisibleDeleteModal()).To(BeFalse())
		})
	})

	Describe("ViewEventDetailModal from ManageSkills", func() {
		Context("when viewing event details from skill events modal", func() {
			It("should NOT show 's: Skills' hint in the footer", func() {
				// The event detail modal shown from ManageSkills should not
				// have the "s: Skills" option since we're already in a skills context

				// Get all skills first
				skills, err := skillRepo.List(ctx, nil)
				Expect(err).NotTo(HaveOccurred())
				Expect(len(skills)).To(BeNumerically(">", 0))

				// Create an event linked to a skill for testing
				testEvent := &domain.CareerEvent{
					Text:    "Built REST API for payment processing",
					Company: "TestCorp",
					Skills:  []string{skills[0].ID}, // Link to first skill (Go)
				}
				err = eventRepo.Create(ctx, testEvent)
				Expect(err).NotTo(HaveOccurred())

				// Update the event to save the skill association
				err = eventRepo.Update(ctx, testEvent)
				Expect(err).NotTo(HaveOccurred())

				// Reload intent to pick up the new data
				intent = intents.NewManageSkillsIntent(intentCtx)
				cmd := intent.Init()
				if cmd != nil {
					msg := cmd()
					intent.Update(msg)
				}

				// Open skill detail modal
				intent.Update(tea.KeyMsg{Type: tea.KeyEnter})
				Expect(intent.HasVisibleDetailModal()).To(BeTrue())

				// Press ctrl+e to view events using this skill
				cmd = intent.Update(tea.KeyMsg{Type: tea.KeyCtrlE})
				if cmd != nil {
					msg := cmd()
					if msg != nil {
						intent.Update(msg)
					}
				}

				// Should now show the skill events modal
				Expect(intent.HasVisibleSkillEventsModal()).To(BeTrue())

				// Press Enter to select the event and show event detail modal
				cmd = intent.Update(tea.KeyMsg{Type: tea.KeyEnter})
				if cmd != nil {
					msg := cmd()
					if msg != nil {
						intent.Update(msg)
					}
				}

				// Event detail modal should be visible
				Expect(intent.HasVisibleEventDetailModal()).To(BeTrue())

				// The view should NOT contain "s: Skills" - we're already in skills context
				view := intent.View()
				Expect(view).NotTo(ContainSubstring("s: Skills"))
				// But should still have close hint
				Expect(view).To(ContainSubstring("Esc: Close"))
			})
		})
	})

	Describe("View Rendering with Modals", func() {
		It("should render list in background when detail modal is shown", func() {
			// Open detail modal
			intent.Update(tea.KeyMsg{Type: tea.KeyEnter})

			view := intent.View()
			// Should contain both modal content and list elements
			Expect(view).To(ContainSubstring("Go"))     // Skill name in modal
			Expect(view).To(ContainSubstring("Skills")) // List header visible in background
		})

		It("should render form in add/edit modal", func() {
			// Open add modal
			cmd := intent.Update(tea.KeyMsg{Type: tea.KeyRunes, Runes: []rune{'n'}})
			if cmd != nil {
				msg := cmd()
				if msg != nil {
					intent.Update(msg)
				}
			}

			view := intent.View()
			// Should contain form elements (Category and Level are form field labels)
			Expect(view).To(ContainSubstring("Category"))
			Expect(view).To(ContainSubstring("Proficiency Level"))
		})

		It("should render confirmation in delete modal", func() {
			// Open delete modal
			cmd := intent.Update(tea.KeyMsg{Type: tea.KeyRunes, Runes: []rune{'d'}})
			if cmd != nil {
				msg := cmd()
				if msg != nil {
					intent.Update(msg)
				}
			}

			view := intent.View()
			// Should contain confirmation elements
			Expect(view).To(ContainSubstring("Delete"))
		})
	})
})
