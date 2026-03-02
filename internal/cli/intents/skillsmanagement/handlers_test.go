package skillsmanagement_test

import (
	"context"
	"errors"

	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"

	"github.com/baphled/kariya/internal/cli/intents"
	"github.com/baphled/kariya/internal/cli/intents/skillsmanagement"
	"github.com/baphled/kariya/internal/cli/screens"
	"github.com/baphled/kariya/internal/cli/uikit/feedback"
	"github.com/baphled/kariya/internal/domain/career"
	careermemory "github.com/baphled/kariya/internal/repository/career/memory"
	"github.com/baphled/kariya/internal/service/career/skillinference"
	"github.com/baphled/kariya/internal/testutil/fixtures"
	tea "github.com/charmbracelet/bubbletea"
)

var _ = Describe("Handlers", func() {
	var (
		ctx      context.Context
		mockRepo *MockSkillRepository
		intent   *skillsmanagement.Intent
	)

	BeforeEach(func() {
		ctx = context.Background()
		mockRepo = NewMockSkillRepository()
		intentCtx := skillsmanagement.NewIntentContext(ctx, mockRepo)
		var err error
		intent, err = skillsmanagement.NewIntent(intentCtx)
		Expect(err).NotTo(HaveOccurred())
		_ = intent.Init()
	})

	Describe("handleSkillCreated", func() {
		Context("when creation succeeds", func() {
			It("transitions to StateList", func() {
				intent.SetState(skillsmanagement.StateAdd)
				skill := fixtures.SkillWith("s1", "Go", "backend", "advanced")
				_ = intent.Update(skillsmanagement.SkillCreatedMsg{Skill: skill})
				Expect(intent.GetState()).To(Equal(skillsmanagement.StateList))
			})

			It("returns a non-nil command to reload", func() {
				skill := fixtures.SkillWith("s1", "Go", "backend", "advanced")
				cmd := intent.Update(skillsmanagement.SkillCreatedMsg{Skill: skill})
				Expect(cmd).NotTo(BeNil())
			})
		})

		Context("when creation fails", func() {
			It("sets result to Failed with CREATE_FAILED code", func() {
				cmd := intent.Update(skillsmanagement.SkillCreatedMsg{
					Error: errors.New("create failed"),
				})
				Expect(intent.Result()).NotTo(BeNil())
				Expect(intent.Result().Status).To(Equal(intents.Failed))
				Expect(intent.Result().Error.Code).To(Equal("CREATE_FAILED"))
				Expect(cmd).To(BeNil())
			})
		})
	})

	Describe("handleSkillUpdated", func() {
		Context("when update succeeds", func() {
			It("transitions to StateList", func() {
				intent.SetState(skillsmanagement.StateEdit)
				skill := fixtures.SkillWith("s1", "Go", "backend", "expert")
				_ = intent.Update(skillsmanagement.SkillUpdatedMsg{Skill: skill})
				Expect(intent.GetState()).To(Equal(skillsmanagement.StateList))
			})

			It("returns a non-nil command to reload", func() {
				skill := fixtures.SkillWith("s1", "Go", "backend", "expert")
				cmd := intent.Update(skillsmanagement.SkillUpdatedMsg{Skill: skill})
				Expect(cmd).NotTo(BeNil())
			})
		})

		Context("when update fails", func() {
			It("sets result to Failed with UPDATE_FAILED code", func() {
				cmd := intent.Update(skillsmanagement.SkillUpdatedMsg{
					Error: errors.New("update failed"),
				})
				Expect(intent.Result()).NotTo(BeNil())
				Expect(intent.Result().Status).To(Equal(intents.Failed))
				Expect(intent.Result().Error.Code).To(Equal("UPDATE_FAILED"))
				Expect(cmd).To(BeNil())
			})
		})
	})

	Describe("handleSkillDeleted", func() {
		Context("when deletion succeeds", func() {
			It("transitions to StateList", func() {
				intent.SetState(skillsmanagement.StateDelete)
				_ = intent.Update(skillsmanagement.SkillDeletedMsg{SkillID: "s1"})
				Expect(intent.GetState()).To(Equal(skillsmanagement.StateList))
			})

			It("returns a non-nil command to reload", func() {
				cmd := intent.Update(skillsmanagement.SkillDeletedMsg{SkillID: "s1"})
				Expect(cmd).NotTo(BeNil())
			})
		})

		Context("when deletion fails", func() {
			It("sets result to Failed with DELETE_FAILED code", func() {
				cmd := intent.Update(skillsmanagement.SkillDeletedMsg{
					Error: errors.New("delete failed"),
				})
				Expect(intent.Result()).NotTo(BeNil())
				Expect(intent.Result().Status).To(Equal(intents.Failed))
				Expect(intent.Result().Error.Code).To(Equal("DELETE_FAILED"))
				Expect(cmd).To(BeNil())
			})
		})
	})

	Describe("handleSkillEventsLoaded", func() {
		Context("when events load successfully", func() {
			It("returns nil command", func() {
				events := fixtures.Events(3)
				cmd := intent.Update(skillsmanagement.SkillEventsLoadedMsg{Events: events})
				Expect(cmd).To(BeNil())
			})
		})

		Context("when loading fails", func() {
			It("returns nil command on error", func() {
				cmd := intent.Update(skillsmanagement.SkillEventsLoadedMsg{
					Error: errors.New("events load failed"),
				})
				Expect(cmd).To(BeNil())
			})
		})
	})

	Describe("handleSkillEventsForModalLoaded", func() {
		Context("when loading fails", func() {
			It("returns nil command on error", func() {
				cmd := intent.Update(skillsmanagement.SkillEventsForModalLoadedMsg{
					Error: errors.New("modal events load failed"),
				})
				Expect(cmd).To(BeNil())
			})
		})
	})

	Describe("handleFeedbackModalUpdate", func() {
		Context("when feedback modal is nil", func() {
			It("passes messages through to screen delegation", func() {
				cmd := intent.Update(tea.KeyMsg{Type: tea.KeyRunes, Runes: []rune{'a'}})
				// With no feedback modal, message passes to screen delegation which may return a cmd
				_ = cmd
			})
		})

		Context("when feedback modal is active", func() {
			BeforeEach(func() {
				intent.Update(skillsmanagement.SkillsCreatedMsg{
					Skills: []*career.Skill{fixtures.SkillWith("s1", "Go", "backend", "advanced")},
				})
				Expect(intent.GetFeedbackModal()).NotTo(BeNil())
			})

			It("clears modal on Esc", func() {
				intent.Update(tea.KeyMsg{Type: tea.KeyEsc})
				Expect(intent.GetFeedbackModal()).To(BeNil())
			})

			It("blocks other messages while visible", func() {
				cmd := intent.Update(tea.KeyMsg{Type: tea.KeyRunes, Runes: []rune{'f'}})
				Expect(cmd).NotTo(BeNil())
			})
		})
	})

	Describe("handleLoadingModalUpdate", func() {
		Context("with inference service configured", func() {
			var inferenceIntent *skillsmanagement.Intent

			BeforeEach(func() {
				eventRepo := careermemory.NewEventRepository()
				skillInferenceService := skillinference.NewSkillInferenceService(mockRepo, mockRepo, eventRepo)

				intentCtx := skillsmanagement.NewIntentContext(ctx, mockRepo)
				intentCtx.EventRepository = eventRepo
				intentCtx.SkillInferenceService = skillInferenceService

				var err error
				inferenceIntent, err = skillsmanagement.NewIntent(intentCtx)
				Expect(err).NotTo(HaveOccurred())
				inferenceIntent.Init()
			})

			It("dismisses loading modal on Esc and returns to list", func() {
				inferenceIntent.Update(tea.KeyMsg{Type: tea.KeyRunes, Runes: []rune{'i'}})
				Expect(inferenceIntent.GetLoadingModal()).NotTo(BeNil())

				inferenceIntent.Update(tea.KeyMsg{Type: tea.KeyEsc})
				Expect(inferenceIntent.GetLoadingModal()).To(BeNil())
				Expect(inferenceIntent.GetState()).To(Equal(skillsmanagement.StateList))
			})

			It("blocks non-Esc key messages while loading", func() {
				inferenceIntent.Update(tea.KeyMsg{Type: tea.KeyRunes, Runes: []rune{'i'}})
				Expect(inferenceIntent.GetLoadingModal()).NotTo(BeNil())

				cmd := inferenceIntent.Update(tea.KeyMsg{Type: tea.KeyRunes, Runes: []rune{'f'}})
				Expect(cmd).NotTo(BeNil())
			})
		})
	})

	Describe("handleKeyShortcuts", func() {
		var skillsIntent *skillsmanagement.Intent

		BeforeEach(func() {
			mockRepo.skills = []*career.Skill{
				fixtures.SkillWith("s1", "Go", "backend", "advanced"),
			}
			intentCtx := skillsmanagement.NewIntentContext(ctx, mockRepo)
			var err error
			skillsIntent, err = skillsmanagement.NewIntent(intentCtx)
			Expect(err).NotTo(HaveOccurred())
			cmd := skillsIntent.Init()
			msg := cmd()
			skillsIntent.Update(msg)
		})

		It("opens filter modal on 'f' key", func() {
			cmd := skillsIntent.Update(tea.KeyMsg{Type: tea.KeyRunes, Runes: []rune{'f'}})
			Expect(cmd).NotTo(BeNil())
		})

		It("opens sort modal on 's' key", func() {
			cmd := skillsIntent.Update(tea.KeyMsg{Type: tea.KeyRunes, Runes: []rune{'s'}})
			Expect(cmd).NotTo(BeNil())
		})

		It("opens search modal on '/' key", func() {
			cmd := skillsIntent.Update(tea.KeyMsg{Type: tea.KeyRunes, Runes: []rune{'/'}})
			Expect(cmd).NotTo(BeNil())
		})

		It("cancels intent on Esc key", func() {
			skillsIntent.Update(tea.KeyMsg{Type: tea.KeyEsc})
			result := skillsIntent.Result()
			Expect(result).NotTo(BeNil())
			Expect(result.Status).To(Equal(intents.Cancelled))
		})

		Context("when filters are active", func() {
			BeforeEach(func() {
				skillsIntent.GetTestContext().Filters.Category = "backend"
			})

			It("clears filters on 'x' key", func() {
				Expect(skillsIntent.HasActiveFilters()).To(BeTrue())
				cmd := skillsIntent.Update(tea.KeyMsg{Type: tea.KeyRunes, Runes: []rune{'x'}})
				Expect(cmd).NotTo(BeNil())
				Expect(skillsIntent.HasActiveFilters()).To(BeFalse())
			})
		})

		Context("when no filters are active", func() {
			It("does nothing on 'x' key", func() {
				cmd := skillsIntent.Update(tea.KeyMsg{Type: tea.KeyRunes, Runes: []rune{'x'}})
				Expect(cmd).To(BeNil())
			})
		})
	})

	Describe("createSkill command", func() {
		It("returns SkillCreatedMsg on success", func() {
			skill := fixtures.SkillWith("", "Go", "backend", "advanced")
			msg := skillsmanagement.SkillCreatedMsg{Skill: skill}
			cmd := intent.Update(msg)
			Expect(cmd).NotTo(BeNil())
		})

		It("returns SkillCreatedMsg with error on failure", func() {
			mockRepo.createErr = errors.New("repo create error")
			msg := skillsmanagement.SkillCreatedMsg{Error: mockRepo.createErr}
			cmd := intent.Update(msg)
			Expect(cmd).To(BeNil())
			Expect(intent.Result().Status).To(Equal(intents.Failed))
		})
	})

	Describe("updateSkill command", func() {
		It("returns SkillUpdatedMsg on success", func() {
			skill := fixtures.SkillWith("s1", "Go Updated", "backend", "expert")
			msg := skillsmanagement.SkillUpdatedMsg{Skill: skill}
			cmd := intent.Update(msg)
			Expect(cmd).NotTo(BeNil())
		})

		It("returns SkillUpdatedMsg with error on failure", func() {
			mockRepo.updateErr = errors.New("repo update error")
			msg := skillsmanagement.SkillUpdatedMsg{Error: mockRepo.updateErr}
			cmd := intent.Update(msg)
			Expect(cmd).To(BeNil())
			Expect(intent.Result().Status).To(Equal(intents.Failed))
		})
	})

	Describe("SkillSuggestionsLoadedMsg with error", func() {
		BeforeEach(func() {
			intent.SetState(skillsmanagement.StateInferringSkills)
		})

		It("shows error modal on non-cancelled error", func() {
			intent.Update(skillsmanagement.SkillSuggestionsLoadedMsg{
				Error: errors.New("inference failed"),
			})
			Expect(intent.GetState()).To(Equal(skillsmanagement.StateList))
			modal := intent.GetFeedbackModal()
			Expect(modal).NotTo(BeNil())
			Expect(modal.Type).To(Equal(feedback.ModalError))
		})
	})

	Describe("handleNavigateData", func() {
		var loadedIntent *skillsmanagement.Intent

		BeforeEach(func() {
			mockRepo.skills = []*career.Skill{
				fixtures.SkillWith("s1", "Go", "backend", "advanced"),
				fixtures.SkillWith("s2", "Python", "backend", "intermediate"),
			}
			intentCtx := skillsmanagement.NewIntentContext(ctx, mockRepo)
			var err error
			loadedIntent, err = skillsmanagement.NewIntent(intentCtx)
			Expect(err).NotTo(HaveOccurred())
			cmd := loadedIntent.Init()
			msg := cmd()
			loadedIntent.Update(msg)
		})

		It("handles add action", func() {
			cmd := loadedIntent.HandleNavigate(&screens.NavigateResult{
				ResultData: map[string]interface{}{
					"action": "add",
				},
			})
			Expect(cmd).NotTo(BeNil())
		})

		It("returns nil for view action with nil skill", func() {
			cmd := loadedIntent.HandleNavigate(&screens.NavigateResult{
				ResultData: map[string]interface{}{
					"action": "view",
					"skill":  nil,
				},
			})
			Expect(cmd).To(BeNil())
		})

		It("returns nil for edit action with nil skill", func() {
			cmd := loadedIntent.HandleNavigate(&screens.NavigateResult{
				ResultData: map[string]interface{}{
					"action": "edit",
					"skill":  nil,
				},
			})
			Expect(cmd).To(BeNil())
		})

		It("returns nil for delete action with nil skill", func() {
			cmd := loadedIntent.HandleNavigate(&screens.NavigateResult{
				ResultData: map[string]interface{}{
					"action": "delete",
					"skill":  nil,
				},
			})
			Expect(cmd).To(BeNil())
		})

		It("returns nil for unknown action", func() {
			cmd := loadedIntent.HandleNavigate(&screens.NavigateResult{
				ResultData: map[string]interface{}{
					"action": "unknown",
				},
			})
			Expect(cmd).To(BeNil())
		})

		It("returns nil for non-map data", func() {
			cmd := loadedIntent.HandleNavigate(&screens.NavigateResult{
				ResultData: 42,
			})
			Expect(cmd).To(BeNil())
		})
	})

	Describe("HandleCancel", func() {
		It("transitions back to StateList", func() {
			intent.SetState(skillsmanagement.StateDetail)
			cmd := intent.HandleCancel(nil)
			Expect(intent.GetState()).To(Equal(skillsmanagement.StateList))
			Expect(cmd).To(BeNil())
		})
	})

	Describe("HandleSubmit", func() {
		It("returns nil", func() {
			cmd := intent.HandleSubmit(nil)
			Expect(cmd).To(BeNil())
		})
	})

	Describe("HandleError", func() {
		It("returns nil", func() {
			cmd := intent.HandleError(nil)
			Expect(cmd).To(BeNil())
		})
	})
})
