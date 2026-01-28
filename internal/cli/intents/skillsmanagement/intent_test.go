package skillsmanagement_test

import (
	"context"

	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"

	"github.com/baphled/kariya/internal/cli/behaviors"
	"github.com/baphled/kariya/internal/cli/intents"
	"github.com/baphled/kariya/internal/cli/intents/skillsmanagement"
	"github.com/baphled/kariya/internal/cli/screens"
	"github.com/baphled/kariya/internal/domain/career"
	tea "github.com/charmbracelet/bubbletea"
)

var _ = Describe("Intent", func() {
	var (
		ctx      context.Context
		mockRepo *MockSkillRepository
	)

	BeforeEach(func() {
		ctx = context.Background()
		mockRepo = NewMockSkillRepository()
	})

	Describe("NewIntent", func() {
		It("should create a new intent with valid context", func() {
			intentCtx := skillsmanagement.NewIntentContext(ctx, mockRepo)
			intent, err := skillsmanagement.NewIntent(intentCtx)
			Expect(err).NotTo(HaveOccurred())
			Expect(intent).NotTo(BeNil())
		})

		It("should set initial state to StateList", func() {
			intentCtx := skillsmanagement.NewIntentContext(ctx, mockRepo)
			intent, _ := skillsmanagement.NewIntent(intentCtx)
			Expect(intent.GetState()).To(Equal(skillsmanagement.StateList))
		})

		It("should not be active initially", func() {
			intentCtx := skillsmanagement.NewIntentContext(ctx, mockRepo)
			intent, _ := skillsmanagement.NewIntent(intentCtx)
			Expect(intent.IsActive()).To(BeFalse())
		})

		It("should store the context", func() {
			intentCtx := skillsmanagement.NewIntentContext(ctx, mockRepo)
			intent, _ := skillsmanagement.NewIntent(intentCtx)
			Expect(intent.GetContext()).To(Equal(intentCtx))
		})

		It("should return error with invalid context", func() {
			intentCtx := &skillsmanagement.IntentContext{
				Ctx:             nil,
				SkillRepository: nil,
			}
			_, err := skillsmanagement.NewIntent(intentCtx)
			Expect(err).To(HaveOccurred())
		})

		It("should embed BaseIntent", func() {
			intentCtx := skillsmanagement.NewIntentContext(ctx, mockRepo)
			intent, _ := skillsmanagement.NewIntent(intentCtx)
			Expect(intent.BaseIntent).NotTo(BeNil())
		})

		It("should initialize TableBehavior for skills", func() {
			intentCtx := skillsmanagement.NewIntentContext(ctx, mockRepo)
			intent, _ := skillsmanagement.NewIntent(intentCtx)
			Expect(intent.GetTableBehavior()).NotTo(BeNil())
		})
	})

	Describe("Init", func() {
		It("should set active to true", func() {
			intentCtx := skillsmanagement.NewIntentContext(ctx, mockRepo)
			intent, _ := skillsmanagement.NewIntent(intentCtx)
			_ = intent.Init()
			Expect(intent.IsActive()).To(BeTrue())
		})

		It("should return a command to load skills", func() {
			intentCtx := skillsmanagement.NewIntentContext(ctx, mockRepo)
			intent, _ := skillsmanagement.NewIntent(intentCtx)
			cmd := intent.Init()
			Expect(cmd).NotTo(BeNil())
		})

		Context("when skills are loaded successfully", func() {
			BeforeEach(func() {
				mockRepo.skills = []*career.Skill{
					{ID: "skill-1", Name: "Go", Category: "Programming"},
					{ID: "skill-2", Name: "Python", Category: "Programming"},
				}
			})

			It("should populate skills after loading", func() {
				intentCtx := skillsmanagement.NewIntentContext(ctx, mockRepo)
				intent, _ := skillsmanagement.NewIntent(intentCtx)
				cmd := intent.Init()

				// Execute the command to get the message
				msg := cmd()
				loadedMsg, ok := msg.(skillsmanagement.SkillsLoadedMsg)
				Expect(ok).To(BeTrue())
				Expect(loadedMsg.Error).NotTo(HaveOccurred())
				Expect(loadedMsg.Skills).To(HaveLen(2))
			})
		})
	})

	Describe("Update", func() {
		var intent *skillsmanagement.Intent

		BeforeEach(func() {
			intentCtx := skillsmanagement.NewIntentContext(ctx, mockRepo)
			intent, _ = skillsmanagement.NewIntent(intentCtx)
			_ = intent.Init()
		})

		Context("when not active", func() {
			It("should return nil", func() {
				intent.SetActive(false)
				cmd := intent.Update(tea.KeyMsg{Type: tea.KeyEnter})
				Expect(cmd).To(BeNil())
			})
		})

		Context("when receiving SkillsLoadedMsg", func() {
			It("should update skills list", func() {
				skills := []*career.Skill{
					{ID: "skill-1", Name: "Go"},
				}
				msg := skillsmanagement.SkillsLoadedMsg{Skills: skills}
				_ = intent.Update(msg)
				Expect(intent.GetSkills()).To(HaveLen(1))
			})

			It("should handle error in loaded message", func() {
				msg := skillsmanagement.SkillsLoadedMsg{
					Error: skillsmanagement.ErrRepositoryNotAvailable,
				}
				_ = intent.Update(msg)
				result := intent.Result()
				Expect(result).NotTo(BeNil())
				Expect(result.Status).To(Equal(intents.Failed))
			})
		})

		Context("when receiving help key", func() {
			It("should toggle help modal and return nil", func() {
				// '?' is the help key
				cmd := intent.Update(tea.KeyMsg{Type: tea.KeyRunes, Runes: []rune{'?'}})
				// Help toggle returns nil (it's handled internally)
				Expect(cmd).To(BeNil())
			})
		})
	})

	Describe("View", func() {
		var intent *skillsmanagement.Intent

		BeforeEach(func() {
			intentCtx := skillsmanagement.NewIntentContext(ctx, mockRepo)
			intent, _ = skillsmanagement.NewIntent(intentCtx)
			_ = intent.Init()
		})

		Context("when not active", func() {
			It("should return empty string", func() {
				intent.SetActive(false)
				view := intent.View()
				Expect(view).To(BeEmpty())
			})
		})

		Context("when active", func() {
			It("should return non-empty view", func() {
				view := intent.View()
				Expect(view).NotTo(BeEmpty())
			})

			It("should include skills header", func() {
				view := intent.View()
				Expect(view).To(ContainSubstring("Skill"))
			})
		})
	})

	Describe("Result", func() {
		var intent *skillsmanagement.Intent

		BeforeEach(func() {
			intentCtx := skillsmanagement.NewIntentContext(ctx, mockRepo)
			intent, _ = skillsmanagement.NewIntent(intentCtx)
		})

		It("should return nil when no result is set", func() {
			result := intent.Result()
			Expect(result).To(BeNil())
		})

		Context("when intent is cancelled", func() {
			It("should return cancelled status", func() {
				_ = intent.Init()
				intent.SetCancelled()
				result := intent.Result()
				Expect(result).NotTo(BeNil())
				Expect(result.Status).To(Equal(intents.Cancelled))
			})
		})
	})

	Describe("Bug Regressions", func() {
		Describe("BUG-016: view detail modal uses correctly-passed skill, not stale index", func() {
			var (
				intent    *skillsmanagement.Intent
				allSkills []*career.Skill
			)

			BeforeEach(func() {
				// Setup 3 skills so index mismatch is detectable.
				allSkills = []*career.Skill{
					{ID: "skill-1", Name: "Go", Category: "Programming"},
					{ID: "skill-2", Name: "Python", Category: "Programming"},
					{ID: "skill-3", Name: "Rust", Category: "Programming"},
				}
				mockRepo.skills = allSkills

				intentCtx := skillsmanagement.NewIntentContext(ctx, mockRepo)
				intent, _ = skillsmanagement.NewIntent(intentCtx)

				// Init and process the loaded message to populate the intent.
				cmd := intent.Init()
				msg := cmd()
				_ = intent.Update(msg)
			})

			It("should set selectedSkill to the navigated skill, not the first skill", func() {
				// Simulate the screen returning a NavigateResult for the 3rd skill.
				// This is what happens when the user presses Enter on skill-3 in the list.
				thirdSkill := allSkills[2]
				navigateResult := &screens.NavigateResult{
					ResultData: map[string]interface{}{
						"action": "view",
						"skill":  thirdSkill,
					},
				}

				_ = intent.HandleNavigate(navigateResult)

				// The selected skill must be the 3rd skill, not the 1st.
				Expect(intent.GetSelectedSkill()).NotTo(BeNil())
				Expect(intent.GetSelectedSkill().ID).To(Equal("skill-3"))
				Expect(intent.GetSelectedSkill().Name).To(Equal("Rust"))
			})

			It("should set selectedSkill correctly for edit action", func() {
				// Edit action passes skill directly - should not be affected.
				secondSkill := allSkills[1]
				navigateResult := &screens.NavigateResult{
					ResultData: map[string]interface{}{
						"action": "edit",
						"skill":  secondSkill,
					},
				}

				_ = intent.HandleNavigate(navigateResult)

				Expect(intent.GetSelectedSkill()).NotTo(BeNil())
				Expect(intent.GetSelectedSkill().ID).To(Equal("skill-2"))
			})

			It("should set selectedSkill correctly for delete action", func() {
				// Delete action passes skill directly - should not be affected.
				thirdSkill := allSkills[2]
				navigateResult := &screens.NavigateResult{
					ResultData: map[string]interface{}{
						"action": "delete",
						"skill":  thirdSkill,
					},
				}

				_ = intent.HandleNavigate(navigateResult)

				Expect(intent.GetSelectedSkill()).NotTo(BeNil())
				Expect(intent.GetSelectedSkill().ID).To(Equal("skill-3"))
			})
		})
	})

	Describe("Interface Compliance", func() {
		It("should implement FilterBehavior interface", func() {
			intentCtx := skillsmanagement.NewIntentContext(ctx, mockRepo)
			intent, _ := skillsmanagement.NewIntent(intentCtx)
			var _ interface {
				HasActiveFilters() bool
				ClearFilters()
				RefreshData() tea.Cmd
			} = intent
		})

		It("should implement ScreenResultHandler interface", func() {
			intentCtx := skillsmanagement.NewIntentContext(ctx, mockRepo)
			intent, _ := skillsmanagement.NewIntent(intentCtx)
			// Verify ScreenResultHandler compliance via behaviors package
			var _ behaviors.ScreenResultHandler = intent
		})
	})
})
