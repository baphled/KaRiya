package skillsmanagement_test

import (
	"context"

	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"

	"github.com/baphled/kariya/internal/domain/career"
	careermemory "github.com/baphled/kariya/internal/repository/career/memory"
	"github.com/baphled/kariya/internal/service/career/skillinference"
	"github.com/baphled/kariya/internal/testutil/fixtures"
	"github.com/baphled/kariya/internal/tui/intents"
	"github.com/baphled/kariya/internal/tui/intents/skillsmanagement"
	"github.com/baphled/kariya/internal/ui/display"
	"github.com/baphled/kariya/internal/ui/uikit/feedback"
	"github.com/baphled/kariya/internal/ui/uikit/widgets"
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
			intentCtx := skillsmanagement.NewIntentValidator(ctx, mockRepo)
			intent, err := skillsmanagement.NewIntent(intentCtx)
			Expect(err).NotTo(HaveOccurred())
			Expect(intent).NotTo(BeNil())
		})

		It("should set initial state to StateList", func() {
			intentCtx := skillsmanagement.NewIntentValidator(ctx, mockRepo)
			intent, _ := skillsmanagement.NewIntent(intentCtx)
			Expect(intent.GetState()).To(Equal(skillsmanagement.StateList))
		})

		It("should not be active initially", func() {
			intentCtx := skillsmanagement.NewIntentValidator(ctx, mockRepo)
			intent, _ := skillsmanagement.NewIntent(intentCtx)
			Expect(intent.IsActive()).To(BeFalse())
		})

		It("should store the context", func() {
			intentCtx := skillsmanagement.NewIntentValidator(ctx, mockRepo)
			intent, _ := skillsmanagement.NewIntent(intentCtx)
			Expect(intent.GetContext()).To(Equal(intentCtx))
		})

		It("should return error with invalid context", func() {
			intentCtx := &skillsmanagement.IntentValidator{
				SkillRepository: nil,
			}
			_, err := skillsmanagement.NewIntent(intentCtx)
			Expect(err).To(HaveOccurred())
		})

		It("should embed BaseIntent", func() {
			intentCtx := skillsmanagement.NewIntentValidator(ctx, mockRepo)
			intent, _ := skillsmanagement.NewIntent(intentCtx)
			Expect(intent.BaseIntent).NotTo(BeNil())
		})

		It("should initialize TableBehavior for skills", func() {
			intentCtx := skillsmanagement.NewIntentValidator(ctx, mockRepo)
			intent, _ := skillsmanagement.NewIntent(intentCtx)
			Expect(intent.GetTableBehavior()).NotTo(BeNil())
		})
	})

	Describe("Init", func() {
		It("should set active to true", func() {
			intentCtx := skillsmanagement.NewIntentValidator(ctx, mockRepo)
			intent, _ := skillsmanagement.NewIntent(intentCtx)
			_ = intent.Init()
			Expect(intent.IsActive()).To(BeTrue())
		})

		It("should return a command to load skills", func() {
			intentCtx := skillsmanagement.NewIntentValidator(ctx, mockRepo)
			intent, _ := skillsmanagement.NewIntent(intentCtx)
			cmd := intent.Init()
			Expect(cmd).NotTo(BeNil())
		})

		Context("when skills are loaded successfully", func() {
			BeforeEach(func() {
				mockRepo.skills = []*career.Skill{
					fixtures.SkillWith("skill-1", "Go", "Programming", "advanced"),
					fixtures.SkillWith("skill-2", "Python", "Programming", "advanced"),
				}
			})

			It("should populate skills after loading", func() {
				intentCtx := skillsmanagement.NewIntentValidator(ctx, mockRepo)
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
			intentCtx := skillsmanagement.NewIntentValidator(ctx, mockRepo)
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
					fixtures.SkillWith("skill-1", "Go", "backend", "advanced"),
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
			intentCtx := skillsmanagement.NewIntentValidator(ctx, mockRepo)
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
			intentCtx := skillsmanagement.NewIntentValidator(ctx, mockRepo)
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
				allSkills = []*career.Skill{
					fixtures.SkillWith("skill-1", "Go", "Programming", "advanced"),
					fixtures.SkillWith("skill-2", "Python", "Programming", "advanced"),
					fixtures.SkillWith("skill-3", "Rust", "Programming", "advanced"),
				}
				mockRepo.skills = allSkills

				intentCtx := skillsmanagement.NewIntentValidator(ctx, mockRepo)
				intent, _ = skillsmanagement.NewIntent(intentCtx)

				// Init and process the loaded message to populate the intent.
				cmd := intent.Init()
				msg := cmd()
				_ = intent.Update(msg)
			})

			It("should set selectedSkill to the navigated skill, not the first skill", func() {
				thirdSkill := allSkills[2]
				navigateResult := &widgets.NavigateViewResult{
					ResultData: map[string]interface{}{
						"action": "view",
						"skill":  display.SkillFromDomain(thirdSkill),
					},
				}

				_ = intent.HandleNavigate(navigateResult)

				// The selected skill must be the 3rd skill, not the 1st.
				Expect(intent.GetSelectedSkill()).NotTo(BeNil())
				Expect(intent.GetSelectedSkill().ID).To(Equal("skill-3"))
				Expect(intent.GetSelectedSkill().Name).To(Equal("Rust"))
			})

			It("should set selectedSkill correctly for edit action", func() {
				secondSkill := allSkills[1]
				navigateResult := &widgets.NavigateViewResult{
					ResultData: map[string]interface{}{
						"action": "edit",
						"skill":  display.SkillFromDomain(secondSkill),
					},
				}

				_ = intent.HandleNavigate(navigateResult)

				Expect(intent.GetSelectedSkill()).NotTo(BeNil())
				Expect(intent.GetSelectedSkill().ID).To(Equal("skill-2"))
			})

			It("should set selectedSkill correctly for delete action", func() {
				thirdSkill := allSkills[2]
				navigateResult := &widgets.NavigateViewResult{
					ResultData: map[string]interface{}{
						"action": "delete",
						"skill":  display.SkillFromDomain(thirdSkill),
					},
				}

				_ = intent.HandleNavigate(navigateResult)

				Expect(intent.GetSelectedSkill()).NotTo(BeNil())
				Expect(intent.GetSelectedSkill().ID).To(Equal("skill-3"))
			})

			It("should close view detail modal on esc without cancelling intent", func() {
				// Open the view detail modal for the third skill
				thirdSkill := allSkills[2]
				navigateResult := &widgets.NavigateViewResult{
					ResultData: map[string]interface{}{
						"action": "view",
						"skill":  display.SkillFromDomain(thirdSkill),
					},
				}
				_ = intent.HandleNavigate(navigateResult)

				// Send ESC key
				cmd := intent.Update(tea.KeyMsg{Type: tea.KeyEsc})
				if cmd != nil {
					_ = cmd()
				}
				result := intent.Result()
				Expect(result).To(BeNil())
			})
		})
	})

	Describe("SkillsCreatedMsg Handling", func() {
		var intent *skillsmanagement.Intent

		BeforeEach(func() {
			intentCtx := skillsmanagement.NewIntentValidator(ctx, mockRepo)
			intent, _ = skillsmanagement.NewIntent(intentCtx)
			_ = intent.Init()
		})

		Context("with inference service configured", func() {
			var inferenceIntent *skillsmanagement.Intent

			BeforeEach(func() {
				eventRepo := careermemory.NewEventRepository()
				skillInferenceService := skillinference.NewSkillInferenceService(mockRepo, mockRepo, eventRepo)

				intentCtx := skillsmanagement.NewIntentValidator(ctx, mockRepo)
				intentCtx.EventRepository = eventRepo
				intentCtx.SkillInferenceService = skillInferenceService

				var err error
				inferenceIntent, err = skillsmanagement.NewIntent(intentCtx)
				Expect(err).NotTo(HaveOccurred())
				inferenceIntent.Init()
			})

			It("should clear loading modal on success", func() {
				inferenceIntent.Update(tea.KeyMsg{Type: tea.KeyRunes, Runes: []rune{'i'}})
				Expect(inferenceIntent.GetLoadingModal()).NotTo(BeNil())

				inferenceIntent.Update(skillsmanagement.SkillsCreatedMsg{
					Skills: []*career.Skill{fixtures.SkillWith("s1", "Go", "backend", "intermediate")},
				})

				Expect(inferenceIntent.GetLoadingModal()).To(BeNil())
			})

			It("should clear loading modal on error", func() {
				inferenceIntent.Update(tea.KeyMsg{Type: tea.KeyRunes, Runes: []rune{'i'}})
				Expect(inferenceIntent.GetLoadingModal()).NotTo(BeNil())

				inferenceIntent.Update(skillsmanagement.SkillsCreatedMsg{
					Error: context.DeadlineExceeded,
				})

				Expect(inferenceIntent.GetLoadingModal()).To(BeNil())
			})
		})

		It("should show success feedback modal with skill count", func() {
			skills := []*career.Skill{
				fixtures.SkillWith("s1", "Go", "Backend", "intermediate"),
				fixtures.SkillWith("s2", "PostgreSQL", "Database", "intermediate"),
			}

			intent.Update(skillsmanagement.SkillsCreatedMsg{Skills: skills})

			modal := intent.GetFeedbackModal()
			Expect(modal).NotTo(BeNil())
			Expect(modal.Type).To(Equal(feedback.ModalSuccess))
			Expect(modal.Message).To(ContainSubstring("2 skill(s)"))
		})

		It("should show error feedback modal on failure", func() {
			intent.Update(skillsmanagement.SkillsCreatedMsg{
				Error: context.DeadlineExceeded,
			})

			modal := intent.GetFeedbackModal()
			Expect(modal).NotTo(BeNil())
			Expect(modal.Type).To(Equal(feedback.ModalError))
			Expect(modal.Title).To(Equal("Skill Creation Failed"))
		})

		It("should transition to StateList on success", func() {
			intent.SetState(skillsmanagement.StateInferringSkills)

			intent.Update(skillsmanagement.SkillsCreatedMsg{
				Skills: []*career.Skill{fixtures.SkillWith("s1", "Go", "backend", "intermediate")},
			})

			Expect(intent.GetState()).To(Equal(skillsmanagement.StateList))
		})

		It("should transition to StateList on error", func() {
			intent.SetState(skillsmanagement.StateInferringSkills)

			intent.Update(skillsmanagement.SkillsCreatedMsg{
				Error: context.DeadlineExceeded,
			})

			Expect(intent.GetState()).To(Equal(skillsmanagement.StateList))
		})

		It("should return a refresh command on success", func() {
			cmd := intent.Update(skillsmanagement.SkillsCreatedMsg{
				Skills: []*career.Skill{fixtures.SkillWith("s1", "Go", "backend", "intermediate")},
			})

			Expect(cmd).NotTo(BeNil())
		})

		It("should return nil command on error", func() {
			cmd := intent.Update(skillsmanagement.SkillsCreatedMsg{
				Error: context.DeadlineExceeded,
			})

			Expect(cmd).To(BeNil())
		})

		It("should process SkillsLoadedMsg even when feedback modal is visible", func() {
			mockRepo.Create(ctx, fixtures.SkillWith("s1", "Go", "Backend", "intermediate"))
			mockRepo.Create(ctx, fixtures.SkillWith("s2", "PostgreSQL", "Database", "intermediate"))

			cmd := intent.Update(skillsmanagement.SkillsCreatedMsg{
				Skills: []*career.Skill{fixtures.Skill("s1"), fixtures.Skill("s2")},
			})
			Expect(cmd).NotTo(BeNil())
			Expect(intent.GetFeedbackModal()).NotTo(BeNil())

			msg := cmd()
			intent.Update(msg)

			Expect(intent.GetSkills()).To(HaveLen(2))
		})
	})

	Describe("SkillSuggestionsLoadedMsg filtering (filterNewSuggestions)", func() {
		var intent *skillsmanagement.Intent

		BeforeEach(func() {
			intentCtx := skillsmanagement.NewIntentValidator(ctx, mockRepo)
			intent, _ = skillsmanagement.NewIntent(intentCtx)
			_ = intent.Init()
			intent.SetState(skillsmanagement.StateInferringSkills)
		})

		It("should show suggestion modal when no existing skills", func() {
			msg := skillsmanagement.SkillSuggestionsLoadedMsg{
				Suggestions: []skillinference.SkillSuggestion{
					{Name: "Go", Category: "backend", Confidence: 0.95},
					{Name: "Docker", Category: "devops", Confidence: 0.85},
				},
				ExistingSkillNames: []string{},
			}

			intent.Update(msg)

			Expect(intent.GetState()).To(Equal(skillsmanagement.StateSkillSuggestionReview))
		})

		It("should show all-tracked modal when all suggestions are existing", func() {
			msg := skillsmanagement.SkillSuggestionsLoadedMsg{
				Suggestions: []skillinference.SkillSuggestion{
					{Name: "Go", Category: "backend", Confidence: 0.95},
					{Name: "Docker", Category: "devops", Confidence: 0.85},
				},
				ExistingSkillNames: []string{"Go", "Docker"},
			}

			intent.Update(msg)

			Expect(intent.GetState()).To(Equal(skillsmanagement.StateList))
			modal := intent.GetFeedbackModal()
			Expect(modal).NotTo(BeNil())
			Expect(modal.Type).To(Equal(feedback.ModalSuccess))
			Expect(modal.Title).To(Equal("All Skills Already Tracked"))
		})

		It("should filter existing and show only new suggestions", func() {
			msg := skillsmanagement.SkillSuggestionsLoadedMsg{
				Suggestions: []skillinference.SkillSuggestion{
					{Name: "Go", Category: "backend", Confidence: 0.95},
					{Name: "Docker", Category: "devops", Confidence: 0.85},
					{Name: "Kubernetes", Category: "devops", Confidence: 0.80},
				},
				ExistingSkillNames: []string{"Go"},
			}

			intent.Update(msg)

			Expect(intent.GetState()).To(Equal(skillsmanagement.StateSkillSuggestionReview))
		})

		It("should filter case-insensitively", func() {
			msg := skillsmanagement.SkillSuggestionsLoadedMsg{
				Suggestions: []skillinference.SkillSuggestion{
					{Name: "Go", Category: "backend", Confidence: 0.95},
					{Name: "Docker", Category: "devops", Confidence: 0.85},
				},
				ExistingSkillNames: []string{"go", "docker"},
			}

			intent.Update(msg)

			Expect(intent.GetState()).To(Equal(skillsmanagement.StateList))
			modal := intent.GetFeedbackModal()
			Expect(modal).NotTo(BeNil())
			Expect(modal.Type).To(Equal(feedback.ModalSuccess))
		})

		It("should show warning when no suggestions detected at all", func() {
			msg := skillsmanagement.SkillSuggestionsLoadedMsg{
				Suggestions:        []skillinference.SkillSuggestion{},
				ExistingSkillNames: []string{},
			}

			intent.Update(msg)

			Expect(intent.GetState()).To(Equal(skillsmanagement.StateList))
			modal := intent.GetFeedbackModal()
			Expect(modal).NotTo(BeNil())
			Expect(modal.Type).To(Equal(feedback.ModalWarning))
		})

		It("should show error modal on inference error", func() {
			msg := skillsmanagement.SkillSuggestionsLoadedMsg{
				Error: context.DeadlineExceeded,
			}

			intent.Update(msg)

			Expect(intent.GetState()).To(Equal(skillsmanagement.StateList))
			modal := intent.GetFeedbackModal()
			Expect(modal).NotTo(BeNil())
			Expect(modal.Type).To(Equal(feedback.ModalError))
		})

		It("should silently ignore cancelled operations", func() {
			msg := skillsmanagement.SkillSuggestionsLoadedMsg{
				Error: context.Canceled,
			}

			intent.Update(msg)

			Expect(intent.GetState()).To(Equal(skillsmanagement.StateList))
			Expect(intent.GetFeedbackModal()).To(BeNil())
		})
	})

	Describe("Interface Compliance", func() {
		It("should implement FilterBehavior interface", func() {
			intentCtx := skillsmanagement.NewIntentValidator(ctx, mockRepo)
			intent, _ := skillsmanagement.NewIntent(intentCtx)
			var _ interface {
				HasActiveFilters() bool
				ClearFilters()
				RefreshData() tea.Cmd
			} = intent
		})

		It("should implement ViewResult handler methods", func() {
			intentCtx := skillsmanagement.NewIntentValidator(ctx, mockRepo)
			intent, _ := skillsmanagement.NewIntent(intentCtx)
			var result widgets.ViewResult = &widgets.CancelViewResult{}
			Expect(intent.HandleCancel(result)).To(BeNil())
			Expect(intent.HandleSubmit(result)).To(BeNil())
			Expect(intent.HandleError(result)).To(BeNil())
		})
	})
})

var _ = Describe("Navigation Key Regression", func() {
	var (
		ctx       context.Context
		mockRepo  *MockSkillRepository
		intent    *skillsmanagement.Intent
		allSkills []*career.Skill
	)

	BeforeEach(func() {
		ctx = context.Background()
		mockRepo = NewMockSkillRepository()
		allSkills = []*career.Skill{
			fixtures.SkillWith("skill-1", "Go", "Programming", "advanced"),
			fixtures.SkillWith("skill-2", "Python", "Programming", "advanced"),
			fixtures.SkillWith("skill-3", "Rust", "Programming", "advanced"),
		}
		mockRepo.skills = allSkills
		intentCtx := skillsmanagement.NewIntentValidator(ctx, mockRepo)
		intent, _ = skillsmanagement.NewIntent(intentCtx)
		cmd := intent.Init()
		msg := cmd()
		intent.Update(msg)
	})

	It("should navigate down then select second skill via enter", func() {
		intent.Update(tea.KeyMsg{Type: tea.KeyDown})
		intent.Update(tea.KeyMsg{Type: tea.KeyEnter})
		Expect(intent.GetSelectedSkill()).NotTo(BeNil())
		Expect(intent.GetSelectedSkill().ID).To(Equal("skill-2"))
	})

	It("should navigate down twice then select third skill via enter", func() {
		intent.Update(tea.KeyMsg{Type: tea.KeyDown})
		intent.Update(tea.KeyMsg{Type: tea.KeyDown})
		intent.Update(tea.KeyMsg{Type: tea.KeyEnter})
		Expect(intent.GetSelectedSkill()).NotTo(BeNil())
		Expect(intent.GetSelectedSkill().ID).To(Equal("skill-3"))
	})

	It("should navigate with j key then select second skill via enter", func() {
		intent.Update(tea.KeyMsg{Type: tea.KeyRunes, Runes: []rune{'j'}})
		intent.Update(tea.KeyMsg{Type: tea.KeyEnter})
		Expect(intent.GetSelectedSkill()).NotTo(BeNil())
		Expect(intent.GetSelectedSkill().ID).To(Equal("skill-2"))
	})

	It("should select first skill when pressing enter without navigation", func() {
		intent.Update(tea.KeyMsg{Type: tea.KeyEnter})
		Expect(intent.GetSelectedSkill()).NotTo(BeNil())
		Expect(intent.GetSelectedSkill().ID).To(Equal("skill-1"))
	})

	It("should still handle shortcut keys at intent level", func() {
		cmd := intent.Update(tea.KeyMsg{Type: tea.KeyRunes, Runes: []rune{'f'}})
		Expect(cmd).NotTo(BeNil())
	})

	It("should cancel intent on escape key when no modal is visible", func() {
		intent.Update(tea.KeyMsg{Type: tea.KeyEsc})
		Expect(intent.IsActive()).To(BeFalse())
		result := intent.Result()
		Expect(result).NotTo(BeNil())
		Expect(result.Status).To(Equal(intents.Cancelled))
	})
})
