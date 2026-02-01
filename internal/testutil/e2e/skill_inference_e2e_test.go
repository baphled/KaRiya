package e2e_test

import (
	"context"
	"sync"
	"time"

	burstmgmt "github.com/baphled/kariya/internal/cli/intents/burst_management"
	skillsmgmt "github.com/baphled/kariya/internal/cli/intents/skillsmanagement"
	"github.com/baphled/kariya/internal/cli/screens"
	"github.com/baphled/kariya/internal/cli/uikit/feedback"
	"github.com/baphled/kariya/internal/domain/career"
	careermemory "github.com/baphled/kariya/internal/repository/career/memory"
	"github.com/baphled/kariya/internal/service/career/skillinference"
	"github.com/baphled/kariya/internal/testutil/mocks"
	tea "github.com/charmbracelet/bubbletea"
	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
)

var _ = Describe("E2E Skill Inference from ManageSkills", func() {
	var (
		skillsIntent          *skillsmgmt.Intent
		skillRepo             *careermemory.SkillRepository
		eventRepo             *careermemory.EventRepository
		skillInferenceService skillinference.SkillInferenceService
	)

	BeforeEach(func() {
		skillRepo = careermemory.NewSkillRepository()
		eventRepo = careermemory.NewEventRepository()
		skillRepo.SetEventRepository(eventRepo)
		eventRepo.SetSkillRepository(skillRepo)
		skillInferenceService = skillinference.NewSkillInferenceService(skillRepo, eventRepo)

		ctx := context.Background()
		skillsCtx := skillsmgmt.NewIntentContext(ctx, skillRepo)
		skillsCtx.EventRepository = eventRepo
		skillsCtx.SkillInferenceService = skillInferenceService

		var err error
		skillsIntent, err = skillsmgmt.NewIntent(skillsCtx)
		Expect(err).NotTo(HaveOccurred())
		skillsIntent.Init()
	})

	Describe("Keybadges", func() {
		It("should show Infer Skills keybadge in the list view help footer", func() {
			view := skillsIntent.View()
			Expect(view).To(ContainSubstring("Infer Skills"))
		})

		It("should restore full footer after dismissing 'No Skills Found' modal", func() {
			skillsIntent.Update(tea.KeyMsg{Type: tea.KeyRunes, Runes: []rune{'i'}})
			skillsIntent.Update(skillsmgmt.SkillSuggestionsLoadedMsg{
				Suggestions: []skillinference.SkillSuggestion{},
			})
			skillsIntent.Update(tea.KeyMsg{Type: tea.KeyEsc})

			Expect(skillsIntent.GetState()).To(Equal(skillsmgmt.StateList))
			view := skillsIntent.View()
			Expect(view).To(ContainSubstring("Infer Skills"))
			Expect(view).To(ContainSubstring("Navigate"))
			Expect(view).To(ContainSubstring("Search"))
		})

		It("should restore full footer after cancelling suggestion review modal", func() {
			skillsIntent.Update(tea.KeyMsg{Type: tea.KeyRunes, Runes: []rune{'i'}})
			skillsIntent.Update(skillsmgmt.SkillSuggestionsLoadedMsg{
				Suggestions: []skillinference.SkillSuggestion{
					{Name: "Go", Category: "Backend", Confidence: 0.9, EventIDs: []string{"e1"}},
				},
			})

			// Cancel the suggestion review modal
			skillsIntent.Update(tea.KeyMsg{Type: tea.KeyEsc})

			Expect(skillsIntent.GetState()).To(Equal(skillsmgmt.StateList))
			view := skillsIntent.View()
			Expect(view).To(ContainSubstring("Infer Skills"))
			Expect(view).To(ContainSubstring("Navigate"))
		})
	})

	Describe("Loading modal", func() {
		It("should show loading modal with correct text when user presses 'i'", func() {
			skillsIntent.Update(tea.KeyMsg{Type: tea.KeyRunes, Runes: []rune{'i'}})

			Expect(skillsIntent.GetState()).To(Equal(skillsmgmt.StateInferringSkills))
			view := skillsIntent.View()
			Expect(view).To(ContainSubstring("Analyzing all events for skills"))
		})

		It("should forward spinner ticks to keep the loading modal animating", func() {
			skillsIntent.Update(tea.KeyMsg{Type: tea.KeyRunes, Runes: []rune{'i'}})

			cmd := skillsIntent.Update(feedback.ModalSpinnerTickMsg{})
			Expect(cmd).NotTo(BeNil())
		})

		It("should block key shortcuts while loading modal is active", func() {
			skillsIntent.Update(tea.KeyMsg{Type: tea.KeyRunes, Runes: []rune{'i'}})

			skillsIntent.Update(tea.KeyMsg{Type: tea.KeyRunes, Runes: []rune{'f'}})
			skillsIntent.Update(tea.KeyMsg{Type: tea.KeyRunes, Runes: []rune{'s'}})

			Expect(skillsIntent.GetState()).To(Equal(skillsmgmt.StateInferringSkills))
		})

		It("should cancel inference and return to list when Esc is pressed during loading", func() {
			skillsIntent.Update(tea.KeyMsg{Type: tea.KeyRunes, Runes: []rune{'i'}})

			skillsIntent.Update(tea.KeyMsg{Type: tea.KeyEsc})

			Expect(skillsIntent.GetState()).To(Equal(skillsmgmt.StateList))
			Expect(skillsIntent.IsActive()).To(BeTrue())
			view := skillsIntent.View()
			Expect(view).To(ContainSubstring("Infer Skills"))
		})
	})

	Describe("Skill suggestions found", func() {
		It("should show skill suggestion review modal with detected skills", func() {
			skillsIntent.Update(tea.KeyMsg{Type: tea.KeyRunes, Runes: []rune{'i'}})

			skillsIntent.Update(skillsmgmt.SkillSuggestionsLoadedMsg{
				Suggestions: []skillinference.SkillSuggestion{
					{Name: "Go", Category: "Backend", Confidence: 0.95, EventIDs: []string{"e1"}},
					{Name: "PostgreSQL", Category: "Database", Confidence: 0.85, EventIDs: []string{"e2"}},
				},
			})

			Expect(skillsIntent.GetState()).To(Equal(skillsmgmt.StateSkillSuggestionReview))

			view := skillsIntent.View()
			Expect(view).To(ContainSubstring("Skill Suggestions"))
			Expect(view).To(ContainSubstring("Go"))
			Expect(view).To(ContainSubstring("PostgreSQL"))
		})

		It("should return to list view when user cancels suggestion review", func() {
			skillsIntent.Update(tea.KeyMsg{Type: tea.KeyRunes, Runes: []rune{'i'}})

			skillsIntent.Update(skillsmgmt.SkillSuggestionsLoadedMsg{
				Suggestions: []skillinference.SkillSuggestion{
					{Name: "Go", Category: "Backend", Confidence: 0.9, EventIDs: []string{"e1"}},
				},
			})
			Expect(skillsIntent.GetState()).To(Equal(skillsmgmt.StateSkillSuggestionReview))

			// Cancel the modal
			skillsIntent.Update(tea.KeyMsg{Type: tea.KeyEsc})

			Expect(skillsIntent.GetState()).To(Equal(skillsmgmt.StateList))
			Expect(skillsIntent.IsActive()).To(BeTrue())
		})
	})

	Describe("No skills found", func() {
		It("should show 'No Skills Found' modal when inference returns empty suggestions", func() {
			skillsIntent.Update(tea.KeyMsg{Type: tea.KeyRunes, Runes: []rune{'i'}})

			skillsIntent.Update(skillsmgmt.SkillSuggestionsLoadedMsg{
				Suggestions: []skillinference.SkillSuggestion{},
			})

			Expect(skillsIntent.GetState()).To(Equal(skillsmgmt.StateList))
			view := skillsIntent.View()
			Expect(view).To(ContainSubstring("No Skills Found"))

			modal := skillsIntent.GetFeedbackModal()
			Expect(modal).NotTo(BeNil())
			Expect(modal.Type).To(Equal(feedback.ModalWarning))
		})

		It("should show error modal when inference returns an error", func() {
			skillsIntent.Update(tea.KeyMsg{Type: tea.KeyRunes, Runes: []rune{'i'}})

			skillsIntent.Update(skillsmgmt.SkillSuggestionsLoadedMsg{
				Error: context.DeadlineExceeded,
			})

			Expect(skillsIntent.GetState()).To(Equal(skillsmgmt.StateList))
			view := skillsIntent.View()
			Expect(view).To(ContainSubstring("Skill Inference Failed"))
		})

		It("should dismiss error modal with Esc and return to normal list view", func() {
			skillsIntent.Update(tea.KeyMsg{Type: tea.KeyRunes, Runes: []rune{'i'}})

			skillsIntent.Update(skillsmgmt.SkillSuggestionsLoadedMsg{
				Suggestions: []skillinference.SkillSuggestion{},
			})

			view := skillsIntent.View()
			Expect(view).To(ContainSubstring("No Skills Found"))

			skillsIntent.Update(tea.KeyMsg{Type: tea.KeyEsc})

			view = skillsIntent.View()
			Expect(view).NotTo(ContainSubstring("No Skills Found"))
			Expect(skillsIntent.IsActive()).To(BeTrue())
			Expect(skillsIntent.GetState()).To(Equal(skillsmgmt.StateList))
		})
	})

	Describe("Skill creation from inference", func() {
		It("should show success modal when skills are created successfully", func() {
			createdSkills := []*career.Skill{
				{ID: "skill-1", Name: "Go", Category: "Backend"},
				{ID: "skill-2", Name: "PostgreSQL", Category: "Database"},
			}

			skillsIntent.Update(skillsmgmt.SkillsCreatedMsg{
				Skills: createdSkills,
				Error:  nil,
			})

			Expect(skillsIntent.GetState()).To(Equal(skillsmgmt.StateList))

			view := skillsIntent.View()
			Expect(view).To(ContainSubstring("Successfully created 2 skill(s)"))

			modal := skillsIntent.GetFeedbackModal()
			Expect(modal).NotTo(BeNil())
			Expect(modal.Type).To(Equal(feedback.ModalSuccess))
		})

		It("should show error modal when skill creation fails", func() {
			skillsIntent.Update(skillsmgmt.SkillsCreatedMsg{
				Skills: nil,
				Error:  context.DeadlineExceeded,
			})

			Expect(skillsIntent.GetState()).To(Equal(skillsmgmt.StateList))

			view := skillsIntent.View()
			Expect(view).To(ContainSubstring("Skill Creation Failed"))

			modal := skillsIntent.GetFeedbackModal()
			Expect(modal).NotTo(BeNil())
			Expect(modal.Type).To(Equal(feedback.ModalError))
		})

		It("should clear loading modal after skills are created", func() {
			skillsIntent.Update(tea.KeyMsg{Type: tea.KeyRunes, Runes: []rune{'i'}})
			Expect(skillsIntent.GetState()).To(Equal(skillsmgmt.StateInferringSkills))

			view := skillsIntent.View()
			Expect(view).To(ContainSubstring("Analyzing all events for skills"))

			skillsIntent.Update(skillsmgmt.SkillsCreatedMsg{
				Skills: []*career.Skill{{ID: "s1", Name: "Go", Category: "Backend"}},
				Error:  nil,
			})

			Expect(skillsIntent.GetState()).To(Equal(skillsmgmt.StateList))

			view = skillsIntent.View()
			Expect(view).NotTo(ContainSubstring("Analyzing all events"))
			Expect(view).To(ContainSubstring("Successfully created 1 skill(s)"))
		})

		It("should clear loading modal on skill creation error", func() {
			skillsIntent.Update(tea.KeyMsg{Type: tea.KeyRunes, Runes: []rune{'i'}})
			Expect(skillsIntent.GetState()).To(Equal(skillsmgmt.StateInferringSkills))

			skillsIntent.Update(skillsmgmt.SkillsCreatedMsg{
				Skills: nil,
				Error:  context.DeadlineExceeded,
			})

			Expect(skillsIntent.GetState()).To(Equal(skillsmgmt.StateList))

			view := skillsIntent.View()
			Expect(view).NotTo(ContainSubstring("Analyzing all events"))
			Expect(view).To(ContainSubstring("Skill Creation Failed"))
		})
	})
})

func executeSkillsBatchCmd(cmd tea.Cmd) []tea.Msg {
	if cmd == nil {
		return nil
	}
	msg := cmd()
	if batchMsg, ok := msg.(tea.BatchMsg); ok {
		return executeBatchMessages(batchMsg)
	}
	if _, isTick := msg.(feedback.ModalSpinnerTickMsg); isTick {
		return nil
	}
	return []tea.Msg{msg}
}

func executeBatchMessages(batchMsg tea.BatchMsg) []tea.Msg {
	var wg sync.WaitGroup
	var mu sync.Mutex
	var messages []tea.Msg
	for _, batchCmd := range batchMsg {
		if batchCmd == nil {
			continue
		}
		wg.Add(1)
		go func(c tea.Cmd) {
			defer wg.Done()
			result := c()
			if result == nil {
				return
			}
			if _, isTick := result.(feedback.ModalSpinnerTickMsg); isTick {
				return
			}
			mu.Lock()
			messages = append(messages, result)
			mu.Unlock()
		}(batchCmd)
	}
	done := make(chan struct{})
	go func() { wg.Wait(); close(done) }()
	select {
	case <-done:
	case <-time.After(5 * time.Second):
	}
	return messages
}

var _ = Describe("E2E Accept All Suggested Skills from ManageSkills", func() {
	var (
		skillsIntent          *skillsmgmt.Intent
		skillRepo             *careermemory.SkillRepository
		eventRepo             *careermemory.EventRepository
		skillInferenceService skillinference.SkillInferenceService
	)

	BeforeEach(func() {
		skillRepo = careermemory.NewSkillRepository()
		eventRepo = careermemory.NewEventRepository()
		skillRepo.SetEventRepository(eventRepo)
		eventRepo.SetSkillRepository(skillRepo)
		skillInferenceService = skillinference.NewSkillInferenceService(skillRepo, eventRepo)

		ctx := context.Background()

		_ = eventRepo.Create(ctx, &career.Event{
			ID:      "e1",
			Text:    "Built REST API using Go and gRPC microservices",
			Date:    time.Date(2024, 6, 15, 0, 0, 0, 0, time.UTC),
			Company: "Acme Corp",
		})
		_ = eventRepo.Create(ctx, &career.Event{
			ID:      "e2",
			Text:    "Designed PostgreSQL schema for user service",
			Date:    time.Date(2024, 7, 20, 0, 0, 0, 0, time.UTC),
			Company: "Acme Corp",
		})

		skillsCtx := skillsmgmt.NewIntentContext(ctx, skillRepo)
		skillsCtx.EventRepository = eventRepo
		skillsCtx.SkillInferenceService = skillInferenceService

		var err error
		skillsIntent, err = skillsmgmt.NewIntent(skillsCtx)
		Expect(err).NotTo(HaveOccurred())
		skillsIntent.Init()
	})

	pressKey := func(r rune) tea.Cmd {
		return skillsIntent.Update(tea.KeyMsg{Type: tea.KeyRunes, Runes: []rune{r}})
	}

	It("should persist all accepted skills to the repository after accepting all suggestions", func() {
		// 1. Press 'i' to trigger skill inference.
		batchCmd := pressKey('i')
		Expect(skillsIntent.GetState()).To(Equal(skillsmgmt.StateInferringSkills))

		// 2. Execute the batch cmd (spinner init + async inference).
		messages := executeSkillsBatchCmd(batchCmd)

		// 3. Feed the SkillSuggestionsLoadedMsg back to intent.
		var suggestionsMsg skillsmgmt.SkillSuggestionsLoadedMsg
		for _, msg := range messages {
			if sm, ok := msg.(skillsmgmt.SkillSuggestionsLoadedMsg); ok {
				suggestionsMsg = sm
				break
			}
		}
		Expect(suggestionsMsg.Error).NotTo(HaveOccurred())
		Expect(suggestionsMsg.Suggestions).ToNot(BeEmpty(), "inference should detect at least 1 skill from events")

		skillsIntent.Update(suggestionsMsg)
		Expect(skillsIntent.GetState()).To(Equal(skillsmgmt.StateSkillSuggestionReview))

		// 4. Press 'a' for each suggestion to accept all.
		suggestionCount := len(suggestionsMsg.Suggestions)
		var createCmd tea.Cmd
		for range suggestionCount {
			createCmd = pressKey('a')
		}

		// 5. After accepting all, modal closes and triggers skill creation.
		Expect(createCmd).NotTo(BeNil(), "accepting last suggestion should return a create command")

		// 6. Execute the batch cmd (spinner init + async creation).
		createMessages := executeSkillsBatchCmd(createCmd)

		// 7. Feed the SkillsCreatedMsg back to intent.
		var createdMsg skillsmgmt.SkillsCreatedMsg
		for _, msg := range createMessages {
			if cm, ok := msg.(skillsmgmt.SkillsCreatedMsg); ok {
				createdMsg = cm
				break
			}
		}
		Expect(createdMsg.Error).NotTo(HaveOccurred(), "skill creation should succeed")
		Expect(createdMsg.Skills).ToNot(BeEmpty(), "at least 1 skill should be created")

		// 8. Feed SkillsCreatedMsg to intent, which shows success modal and triggers refresh.
		refreshCmd := skillsIntent.Update(createdMsg)
		Expect(refreshCmd).NotTo(BeNil(), "success handler should return a refresh command")
		Expect(skillsIntent.GetFeedbackModal()).NotTo(BeNil(), "success modal should be visible")
		Expect(skillsIntent.GetFeedbackModal().Type).To(Equal(feedback.ModalSuccess))

		// 9. Execute the refresh command to get SkillsLoadedMsg.
		refreshMsg := refreshCmd()
		loadedMsg, ok := refreshMsg.(skillsmgmt.SkillsLoadedMsg)
		Expect(ok).To(BeTrue(), "refresh should return SkillsLoadedMsg")
		Expect(loadedMsg.Error).NotTo(HaveOccurred())
		Expect(loadedMsg.Skills).To(HaveLen(len(createdMsg.Skills)), "loaded skills should match created skills")

		// 10. Feed SkillsLoadedMsg while feedback modal is visible (exercises the bypass fix).
		skillsIntent.Update(loadedMsg)

		// 11. Assert skills are in the intent's list.
		intentSkills := skillsIntent.GetSkills()
		Expect(intentSkills).To(HaveLen(len(createdMsg.Skills)), "intent should have all created skills")

		// 12. Assert skills are actually persisted in the repository.
		repoSkills, err := skillRepo.List(context.Background(), nil)
		Expect(err).NotTo(HaveOccurred())
		Expect(repoSkills).To(HaveLen(len(createdMsg.Skills)), "repository should have all created skills")

		// 13. Verify skill names match what was suggested.
		skillNames := make(map[string]bool)
		for _, s := range repoSkills {
			skillNames[s.Name] = true
		}
		for _, suggestion := range suggestionsMsg.Suggestions {
			Expect(skillNames).To(HaveKey(suggestion.Name), "repo should contain skill: "+suggestion.Name)
		}

		// 14. Verify events are linked to skills.
		for _, s := range repoSkills {
			Expect(s.ID).NotTo(BeEmpty(), "created skill should have an ID")
		}
	})
})

var _ = Describe("E2E Skill Suggestion Event Drill-Down from ManageSkills", func() {
	var (
		skillsIntent          *skillsmgmt.Intent
		skillRepo             *careermemory.SkillRepository
		eventRepo             *careermemory.EventRepository
		skillInferenceService skillinference.SkillInferenceService
	)

	BeforeEach(func() {
		skillRepo = careermemory.NewSkillRepository()
		eventRepo = careermemory.NewEventRepository()
		skillRepo.SetEventRepository(eventRepo)
		eventRepo.SetSkillRepository(skillRepo)
		skillInferenceService = skillinference.NewSkillInferenceService(skillRepo, eventRepo)

		ctx := context.Background()

		event1 := &career.Event{
			ID:      "e1",
			Text:    "Built REST API with Go and gRPC",
			Date:    time.Date(2024, 6, 15, 0, 0, 0, 0, time.UTC),
			Company: "Acme Corp",
		}
		event2 := &career.Event{
			ID:      "e2",
			Text:    "Designed PostgreSQL schema for user service",
			Date:    time.Date(2024, 7, 20, 0, 0, 0, 0, time.UTC),
			Company: "Acme Corp",
		}
		_ = eventRepo.Create(ctx, event1)
		_ = eventRepo.Create(ctx, event2)

		skillsCtx := skillsmgmt.NewIntentContext(ctx, skillRepo)
		skillsCtx.EventRepository = eventRepo
		skillsCtx.SkillInferenceService = skillInferenceService

		var err error
		skillsIntent, err = skillsmgmt.NewIntent(skillsCtx)
		Expect(err).NotTo(HaveOccurred())
		skillsIntent.Init()
	})

	openSuggestionModal := func(suggestions []skillinference.SkillSuggestion) {
		skillsIntent.Update(tea.KeyMsg{Type: tea.KeyRunes, Runes: []rune{'i'}})
		skillsIntent.Update(skillsmgmt.SkillSuggestionsLoadedMsg{
			Suggestions: suggestions,
		})
		Expect(skillsIntent.GetState()).To(Equal(skillsmgmt.StateSkillSuggestionReview))
	}

	Describe("Enter on skill suggestion shows events", func() {
		It("should show events table when pressing Enter on a skill suggestion", func() {
			openSuggestionModal([]skillinference.SkillSuggestion{
				{Name: "Go", Category: "Backend", Confidence: 0.95, EventIDs: []string{"e1"}},
			})

			skillsIntent.Update(tea.KeyMsg{Type: tea.KeyEnter})

			view := skillsIntent.View()
			Expect(view).To(ContainSubstring("Events"))
			Expect(view).To(ContainSubstring("Built REST API"))
		})

		It("should show multiple events in the table for a multi-event suggestion", func() {
			openSuggestionModal([]skillinference.SkillSuggestion{
				{Name: "Go", Category: "Backend", Confidence: 0.95, EventIDs: []string{"e1", "e2"}},
			})

			skillsIntent.Update(tea.KeyMsg{Type: tea.KeyEnter})

			view := skillsIntent.View()
			Expect(view).To(ContainSubstring("Built REST API"))
			Expect(view).To(ContainSubstring("Designed PostgreSQL"))
		})

		It("should show event date and company in the events table", func() {
			openSuggestionModal([]skillinference.SkillSuggestion{
				{Name: "Go", Category: "Backend", Confidence: 0.95, EventIDs: []string{"e1"}},
			})

			skillsIntent.Update(tea.KeyMsg{Type: tea.KeyEnter})

			view := skillsIntent.View()
			Expect(view).To(ContainSubstring("2024"))
			Expect(view).To(ContainSubstring("Acme Corp"))
		})

		It("should return to suggestion review when pressing Esc from events table", func() {
			openSuggestionModal([]skillinference.SkillSuggestion{
				{Name: "Go", Category: "Backend", Confidence: 0.95, EventIDs: []string{"e1"}},
			})

			skillsIntent.Update(tea.KeyMsg{Type: tea.KeyEnter})
			skillsIntent.Update(tea.KeyMsg{Type: tea.KeyEsc})

			Expect(skillsIntent.GetState()).To(Equal(skillsmgmt.StateSkillSuggestionReview))
			view := skillsIntent.View()
			Expect(view).To(ContainSubstring("Skill Suggestions"))
			Expect(view).To(ContainSubstring("Go"))
		})

		It("should remain active after closing events table", func() {
			openSuggestionModal([]skillinference.SkillSuggestion{
				{Name: "Go", Category: "Backend", Confidence: 0.95, EventIDs: []string{"e1"}},
			})

			skillsIntent.Update(tea.KeyMsg{Type: tea.KeyEnter})
			skillsIntent.Update(tea.KeyMsg{Type: tea.KeyEsc})

			Expect(skillsIntent.IsActive()).To(BeTrue())
		})
	})

	Describe("Edge cases for event drill-down", func() {
		It("should handle suggestion with non-existent event IDs gracefully", func() {
			openSuggestionModal([]skillinference.SkillSuggestion{
				{Name: "Go", Category: "Backend", Confidence: 0.95, EventIDs: []string{"nonexistent-id"}},
			})

			skillsIntent.Update(tea.KeyMsg{Type: tea.KeyEnter})

			Expect(skillsIntent.IsActive()).To(BeTrue())
			view := skillsIntent.View()
			Expect(view).To(ContainSubstring("No events"))
		})

		It("should handle suggestion with empty event IDs gracefully", func() {
			openSuggestionModal([]skillinference.SkillSuggestion{
				{Name: "Go", Category: "Backend", Confidence: 0.95, EventIDs: []string{}},
			})

			skillsIntent.Update(tea.KeyMsg{Type: tea.KeyEnter})

			Expect(skillsIntent.IsActive()).To(BeTrue())
			view := skillsIntent.View()
			Expect(view).To(ContainSubstring("No events"))
		})

		It("should handle mixed valid and invalid event IDs", func() {
			openSuggestionModal([]skillinference.SkillSuggestion{
				{Name: "Go", Category: "Backend", Confidence: 0.95, EventIDs: []string{"e1", "nonexistent"}},
			})

			skillsIntent.Update(tea.KeyMsg{Type: tea.KeyEnter})

			view := skillsIntent.View()
			Expect(view).To(ContainSubstring("Built REST API"))
		})
	})
})

var _ = Describe("E2E Skill Suggestion Event Drill-Down from BurstManagement", func() {
	var (
		intent      *burstmgmt.Intent
		mockService *mocks.BurstServiceMock
		burst       *career.Burst
		event1      *career.Event
		event2      *career.Event
	)

	BeforeEach(func() {
		burstRepo := careermemory.NewBurstRepository()
		mockService = mocks.NewBurstServiceMock()

		skillRepo := careermemory.NewSkillRepository()
		eventRepo := careermemory.NewEventRepository()
		skillRepo.SetEventRepository(eventRepo)
		eventRepo.SetSkillRepository(skillRepo)
		skillInferenceService := skillinference.NewSkillInferenceService(skillRepo, eventRepo)

		event1 = &career.Event{
			ID:      "event-1",
			Text:    "Built REST API with Go and gRPC",
			Date:    time.Date(2024, 6, 15, 0, 0, 0, 0, time.UTC),
			Company: "Acme Corp",
		}
		event2 = &career.Event{
			ID:      "event-2",
			Text:    "Designed PostgreSQL schema for user service",
			Date:    time.Date(2024, 7, 20, 0, 0, 0, 0, time.UTC),
			Company: "Acme Corp",
		}
		mockService.SetEvents([]*career.Event{event1, event2})

		burst = &career.Burst{
			ID:          "burst-1",
			Name:        "API Development",
			Description: "Built REST API with Go and PostgreSQL",
			EventIDs:    []string{"event-1", "event-2"},
			Confirmed:   true,
		}

		ctx := &burstmgmt.IntentContext{
			Bursts:                []*career.Burst{burst},
			Service:               mockService,
			BurstRepository:       burstRepo,
			SkillInferenceService: skillInferenceService,
		}
		ctx.Validate()

		var err error
		intent, err = burstmgmt.NewIntent(ctx)
		Expect(err).NotTo(HaveOccurred())
		intent.Init()
	})

	openSuggestionModal := func(suggestions []skillinference.SkillSuggestion) {
		intent.Update(burstmgmt.SkillSuggestionsLoadedMsg{
			Suggestions: suggestions,
		})
		Expect(intent.GetState()).To(Equal(burstmgmt.StateSkillSuggestionReview))
	}

	Describe("Enter on skill suggestion shows events", func() {
		It("should show events table when pressing Enter on a skill suggestion", func() {
			openSuggestionModal([]skillinference.SkillSuggestion{
				{Name: "Go", Category: "Backend", Confidence: 0.95, EventIDs: []string{"event-1"}},
			})

			intent.Update(tea.KeyMsg{Type: tea.KeyEnter})

			view := intent.View()
			Expect(view).To(ContainSubstring("Events"))
			Expect(view).To(ContainSubstring("Built REST API"))
		})

		It("should return to suggestion review when pressing Esc from events table", func() {
			openSuggestionModal([]skillinference.SkillSuggestion{
				{Name: "Go", Category: "Backend", Confidence: 0.95, EventIDs: []string{"event-1"}},
			})

			intent.Update(tea.KeyMsg{Type: tea.KeyEnter})
			intent.Update(tea.KeyMsg{Type: tea.KeyEsc})

			Expect(intent.GetState()).To(Equal(burstmgmt.StateSkillSuggestionReview))
			view := intent.View()
			Expect(view).To(ContainSubstring("Skill Suggestions"))
		})
	})

	Describe("Esc from suggestion review returns to burst detail", func() {
		It("should return to burst detail when Esc is pressed after inference from detail", func() {
			intent.HandleNavigate(&screens.NavigateResult{
				ResultData: burst,
			})

			intent.Update(tea.KeyMsg{Type: tea.KeyRunes, Runes: []rune{'i'}})
			Expect(intent.GetState()).To(Equal(burstmgmt.StateInferringSkills))

			intent.Update(burstmgmt.SkillSuggestionsLoadedMsg{
				Suggestions: []skillinference.SkillSuggestion{
					{Name: "Go", Category: "Backend", Confidence: 0.95, EventIDs: []string{"event-1"}},
				},
			})
			Expect(intent.GetState()).To(Equal(burstmgmt.StateSkillSuggestionReview))

			intent.Update(tea.KeyMsg{Type: tea.KeyEsc})

			view := intent.View()
			Expect(view).To(ContainSubstring("API Development"))
			Expect(intent.IsActive()).To(BeTrue())
		})

		It("should return to list when Esc is pressed after inference NOT from detail", func() {
			intent.Update(burstmgmt.SkillSuggestionsLoadedMsg{
				Suggestions: []skillinference.SkillSuggestion{
					{Name: "Go", Category: "Backend", Confidence: 0.95, EventIDs: []string{"event-1"}},
				},
			})

			intent.Update(tea.KeyMsg{Type: tea.KeyEsc})

			Expect(intent.GetState()).To(Equal(burstmgmt.StateList))
			Expect(intent.IsActive()).To(BeTrue())
		})

		It("should return to burst detail after completing event drill-down from detail inference", func() {
			intent.HandleNavigate(&screens.NavigateResult{
				ResultData: burst,
			})

			intent.Update(tea.KeyMsg{Type: tea.KeyRunes, Runes: []rune{'i'}})
			intent.Update(burstmgmt.SkillSuggestionsLoadedMsg{
				Suggestions: []skillinference.SkillSuggestion{
					{Name: "Go", Category: "Backend", Confidence: 0.95, EventIDs: []string{"event-1"}},
				},
			})

			intent.Update(tea.KeyMsg{Type: tea.KeyEnter})
			intent.Update(tea.KeyMsg{Type: tea.KeyEsc})

			Expect(intent.GetState()).To(Equal(burstmgmt.StateSkillSuggestionReview))
			view := intent.View()
			Expect(view).To(ContainSubstring("Skill Suggestions"))

			intent.Update(tea.KeyMsg{Type: tea.KeyEsc})

			view = intent.View()
			Expect(view).To(ContainSubstring("API Development"))
		})
	})

	Describe("Edge cases", func() {
		It("should handle Enter with no events found gracefully", func() {
			openSuggestionModal([]skillinference.SkillSuggestion{
				{Name: "Go", Category: "Backend", Confidence: 0.95, EventIDs: []string{"nonexistent"}},
			})

			intent.Update(tea.KeyMsg{Type: tea.KeyEnter})

			Expect(intent.IsActive()).To(BeTrue())
			view := intent.View()
			Expect(view).To(ContainSubstring("No events"))
		})
	})
})

var _ = Describe("E2E Skill Inference Workflow", func() {
	var intent *burstmgmt.Intent
	var ctx *burstmgmt.IntentContext
	var mockService *mocks.BurstServiceMock
	var burst *career.Burst

	BeforeEach(func() {
		burstRepo := careermemory.NewBurstRepository()
		mockService = mocks.NewBurstServiceMock()

		skillRepo := careermemory.NewSkillRepository()
		eventRepo := careermemory.NewEventRepository()
		skillRepo.SetEventRepository(eventRepo)
		eventRepo.SetSkillRepository(skillRepo)

		skillInferenceService := skillinference.NewSkillInferenceService(skillRepo, eventRepo)

		// Create events that the burst references
		event1 := &career.Event{
			ID:      "event-1",
			Text:    "Built REST API with Go and gRPC",
			Date:    time.Date(2024, 6, 15, 0, 0, 0, 0, time.UTC),
			Company: "Acme Corp",
		}
		event2 := &career.Event{
			ID:      "event-2",
			Text:    "Designed PostgreSQL schema for user service",
			Date:    time.Date(2024, 7, 20, 0, 0, 0, 0, time.UTC),
			Company: "Acme Corp",
		}
		mockService.SetEvents([]*career.Event{event1, event2})

		// Create a burst with events
		burst = &career.Burst{
			ID:          "burst-1",
			Name:        "API Development",
			Description: "Built REST API with Go and PostgreSQL",
			EventIDs:    []string{"event-1", "event-2"},
			Confirmed:   true,
		}

		ctx = &burstmgmt.IntentContext{
			Bursts:                []*career.Burst{burst},
			Service:               mockService,
			BurstRepository:       burstRepo,
			SkillInferenceService: skillInferenceService,
		}
		ctx.Validate()

		var err error
		intent, err = burstmgmt.NewIntent(ctx)
		Expect(err).NotTo(HaveOccurred())
		intent.Init()
	})

	Describe("Manual Skill Inference Trigger - Burst Detail Modal", func() {
		It("should infer skills when user presses 'i' in burst detail modal", func() {
			// Given: A confirmed burst is selected and detail modal is showing
			burst.Confirmed = true

			// Navigate to burst detail (sets selectedBurst and shows modal)
			intent.HandleNavigate(&screens.NavigateResult{
				ResultData: burst,
			})

			// Verify we're still in list state (modal overlays on top)
			Expect(intent.GetState()).To(Equal(burstmgmt.StateList))

			// When: User presses 'i' to infer skills from the detail modal
			intent.Update(tea.KeyMsg{Type: tea.KeyRunes, Runes: []rune{'i'}})

			// Then: State should transition to StateInferringSkills
			Expect(intent.GetState()).To(Equal(burstmgmt.StateInferringSkills))

			// And: Loading modal should show skill inference in progress
			view := intent.View()
			Expect(view).To(ContainSubstring("Detecting skills from burst events"))
		})

		It("should not allow skill inference on unconfirmed bursts", func() {
			// Given: An unconfirmed burst
			burst.Confirmed = false

			// Navigate to burst detail
			intent.HandleNavigate(&screens.NavigateResult{
				ResultData: burst,
			})

			// When: User presses 'i' to infer skills
			intent.Update(tea.KeyMsg{Type: tea.KeyRunes, Runes: []rune{'i'}})

			// Then: State should NOT transition to inferring - should stay on list
			Expect(intent.GetState()).NotTo(Equal(burstmgmt.StateInferringSkills))

			// And: Warning modal should tell user to confirm burst first (not error — no failure occurred)
			view := intent.View()
			Expect(view).To(ContainSubstring("Confirm"))
			Expect(intent.IsActive()).To(BeTrue())

			modal := intent.GetFeedbackModal()
			Expect(modal).NotTo(BeNil())
			Expect(modal.Type).To(Equal(feedback.ModalWarning))
		})
	})

	Describe("Auto Skill Inference Trigger", func() {
		It("should automatically infer skills after fact extraction completes on a confirmed burst", func() {
			// Given: A confirmed burst that was selected for fact extraction
			burst.Confirmed = true

			// Simulate burst selection by navigating to it (this sets selectedBurst)
			// This happens when user selects a burst from the list screen
			intent.HandleNavigate(&screens.NavigateResult{
				ResultData: burst,
			})

			// Create some extracted facts
			facts := []*career.Fact{
				{
					ID:                   "fact-1",
					Text:                 "Expert in Go programming",
					SourceEventID:        burst.EventIDs[0],
					CompetencyCategories: []string{"technical"},
				},
				{
					ID:                   "fact-2",
					Text:                 "PostgreSQL database design",
					SourceEventID:        burst.EventIDs[1],
					CompetencyCategories: []string{"technical"},
				},
			}

			// When: Fact extraction completes successfully
			intent.Update(burstmgmt.FactExtractionCompleteMsg{
				Facts: facts,
				Error: nil,
			})

			// Then: State should transition to StateInferringSkills
			Expect(intent.GetState()).To(Equal(burstmgmt.StateInferringSkills))

			// And: View should show skill inference in progress
			view := intent.View()
			Expect(view).To(ContainSubstring("Detecting skills from burst events"))
		})
	})

	Describe("Skill Inference Message Handling", func() {
		It("should display skill suggestions modal when skills are detected", func() {
			// Given: Skill inference returns suggestions
			suggestions := []skillinference.SkillSuggestion{
				{
					Name:       "Go",
					Category:   "Backend",
					Confidence: 0.95,
					EventIDs:   []string{burst.EventIDs[0]},
					Contexts:   []string{"Built API with Go"},
				},
				{
					Name:       "PostgreSQL",
					Category:   "Database",
					Confidence: 0.85,
					EventIDs:   []string{burst.EventIDs[0]},
					Contexts:   []string{"Integrated PostgreSQL database"},
				},
			}

			// When: SkillSuggestionsLoadedMsg is sent
			intent.Update(burstmgmt.SkillSuggestionsLoadedMsg{
				Suggestions: suggestions,
				Error:       nil,
			})

			// Then: State should be StateSkillSuggestionReview
			Expect(intent.GetState()).To(Equal(burstmgmt.StateSkillSuggestionReview))

			// And: View should show skill suggestion modal
			view := intent.View()
			Expect(view).To(ContainSubstring("Skill Suggestions"))
			Expect(view).To(ContainSubstring("Go"))
			Expect(view).To(ContainSubstring("PostgreSQL"))
			Expect(view).To(ContainSubstring("Backend"))
			Expect(view).To(ContainSubstring("Database"))
		})

		It("should show error modal when skill inference fails", func() {
			// When: SkillSuggestionsErrorMsg is sent with an error
			intent.Update(burstmgmt.SkillSuggestionsErrorMsg{
				Err: context.DeadlineExceeded,
			})

			// Then: State should return to list
			Expect(intent.GetState()).To(Equal(burstmgmt.StateList))

			// And: Error modal should be shown
			view := intent.View()
			Expect(view).To(ContainSubstring("Skill Inference Failed"))
		})

		It("should handle no skills detected gracefully", func() {
			// When: SkillSuggestionsLoadedMsg is sent with empty suggestions
			intent.Update(burstmgmt.SkillSuggestionsLoadedMsg{
				Suggestions: []skillinference.SkillSuggestion{},
				Error:       nil,
			})

			// Then: State should return to list
			Expect(intent.GetState()).To(Equal(burstmgmt.StateList))

			// And: Informative warning should be shown (not error — no failure occurred)
			view := intent.View()
			Expect(view).To(ContainSubstring("No skills were detected"))

			modal := intent.GetFeedbackModal()
			Expect(modal).NotTo(BeNil())
			Expect(modal.Type).To(Equal(feedback.ModalWarning))
		})

		It("should show success message when skills are created", func() {
			// Given: Skills have been created from suggestions
			createdSkills := []*career.Skill{
				{
					ID:       "skill-1",
					Name:     "Go",
					Category: "Backend",
				},
				{
					ID:       "skill-2",
					Name:     "PostgreSQL",
					Category: "Database",
				},
			}

			// When: SkillsCreatedMsg is sent
			intent.Update(burstmgmt.SkillsCreatedMsg{
				Skills: createdSkills,
				Error:  nil,
			})

			// Then: State should return to list
			Expect(intent.GetState()).To(Equal(burstmgmt.StateList))

			// And: Success message should be shown with success icon, not warning
			view := intent.View()
			Expect(view).To(ContainSubstring("Skills Created"))
			Expect(view).To(ContainSubstring("Successfully created 2 skill"))
			Expect(view).To(ContainSubstring("✅"))
			Expect(view).NotTo(ContainSubstring("⚠️"))

			modal := intent.GetFeedbackModal()
			Expect(modal).NotTo(BeNil())
			Expect(modal.Type).To(Equal(feedback.ModalSuccess))
		})

		It("should handle skill creation errors gracefully", func() {
			// When: SkillsCreatedMsg is sent with an error
			intent.Update(burstmgmt.SkillsCreatedMsg{
				Skills: nil,
				Error:  context.DeadlineExceeded,
			})

			// Then: State should return to list
			Expect(intent.GetState()).To(Equal(burstmgmt.StateList))

			// And: Error modal should be shown
			view := intent.View()
			Expect(view).To(ContainSubstring("Skill Creation Failed"))
		})

		It("should handle user cancellation silently", func() {
			// When: User cancels skill inference (context.Canceled)
			intent.Update(burstmgmt.SkillSuggestionsErrorMsg{
				Err: context.Canceled,
			})

			// Then: State should return to list silently
			Expect(intent.GetState()).To(Equal(burstmgmt.StateList))

			// And: No error modal should be shown for cancelled operations
			view := intent.View()
			Expect(view).NotTo(ContainSubstring("Skill Inference Failed"))
		})
	})
})

var _ = Describe("E2E Skill Inference Excludes Existing Skills", func() {
	Describe("Trigger 1: Burst Management - manual 'i' from detail modal", func() {
		var (
			intent      *burstmgmt.Intent
			mockService *mocks.BurstServiceMock
			burst       *career.Burst
			skillRepo   *careermemory.SkillRepository
			eventRepo   *careermemory.EventRepository
		)

		BeforeEach(func() {
			burstRepo := careermemory.NewBurstRepository()
			mockService = mocks.NewBurstServiceMock()
			skillRepo = careermemory.NewSkillRepository()
			eventRepo = careermemory.NewEventRepository()
			skillRepo.SetEventRepository(eventRepo)
			eventRepo.SetSkillRepository(skillRepo)
			skillInferenceService := skillinference.NewSkillInferenceService(skillRepo, eventRepo)

			ctx := context.Background()

			event1 := &career.Event{
				ID:      "event-1",
				Text:    "Built REST API with Go and Docker containers",
				Date:    time.Date(2024, 6, 15, 0, 0, 0, 0, time.UTC),
				Company: "Acme Corp",
			}
			event2 := &career.Event{
				ID:      "event-2",
				Text:    "Designed PostgreSQL schema for user service",
				Date:    time.Date(2024, 7, 20, 0, 0, 0, 0, time.UTC),
				Company: "Acme Corp",
			}
			_ = eventRepo.Create(ctx, event1)
			_ = eventRepo.Create(ctx, event2)
			mockService.SetEvents([]*career.Event{event1, event2})

			_ = skillRepo.Create(ctx, &career.Skill{ID: "existing-go", Name: "Go", Category: "backend"})
			_ = skillRepo.Create(ctx, &career.Skill{ID: "existing-pg", Name: "PostgreSQL", Category: "database"})

			burst = &career.Burst{
				ID:          "burst-1",
				Name:        "API Development",
				Description: "Built REST API with Go and PostgreSQL",
				EventIDs:    []string{"event-1", "event-2"},
				Confirmed:   true,
			}

			intentCtx := &burstmgmt.IntentContext{
				Bursts:                []*career.Burst{burst},
				Service:               mockService,
				BurstRepository:       burstRepo,
				SkillInferenceService: skillInferenceService,
			}
			intentCtx.Validate()

			var err error
			intent, err = burstmgmt.NewIntent(intentCtx)
			Expect(err).NotTo(HaveOccurred())
			intent.Init()
		})

		It("should include existing skills in suggestions and report them in ExistingSkillNames", func() {
			intent.HandleNavigate(&screens.NavigateResult{
				ResultData: burst,
			})

			batchCmd := intent.Update(tea.KeyMsg{Type: tea.KeyRunes, Runes: []rune{'i'}})
			Expect(intent.GetState()).To(Equal(burstmgmt.StateInferringSkills))

			messages := executeSkillsBatchCmd(batchCmd)

			var suggestionsMsg burstmgmt.SkillSuggestionsLoadedMsg
			for _, msg := range messages {
				if sm, ok := msg.(burstmgmt.SkillSuggestionsLoadedMsg); ok {
					suggestionsMsg = sm
					break
				}
			}
			Expect(suggestionsMsg.Error).NotTo(HaveOccurred())

			suggestionNames := make([]string, 0, len(suggestionsMsg.Suggestions))
			for _, s := range suggestionsMsg.Suggestions {
				suggestionNames = append(suggestionNames, s.Name)
			}

			Expect(suggestionNames).To(ContainElement("Go"),
				"Go should appear in suggestions for burst context visibility")
			Expect(suggestionNames).To(ContainElement("PostgreSQL"),
				"PostgreSQL should appear in suggestions for burst context visibility")
			Expect(suggestionNames).To(ContainElement("Docker"),
				"Docker does not exist in the repo and should be suggested")

			Expect(suggestionsMsg.ExistingSkillNames).To(ContainElement("Go"))
			Expect(suggestionsMsg.ExistingSkillNames).To(ContainElement("PostgreSQL"))
		})
	})

	Describe("Trigger 2: Burst Management - auto after fact extraction", func() {
		var (
			intent      *burstmgmt.Intent
			mockService *mocks.BurstServiceMock
			burst       *career.Burst
			skillRepo   *careermemory.SkillRepository
			eventRepo   *careermemory.EventRepository
		)

		BeforeEach(func() {
			burstRepo := careermemory.NewBurstRepository()
			mockService = mocks.NewBurstServiceMock()
			skillRepo = careermemory.NewSkillRepository()
			eventRepo = careermemory.NewEventRepository()
			skillRepo.SetEventRepository(eventRepo)
			eventRepo.SetSkillRepository(skillRepo)
			skillInferenceService := skillinference.NewSkillInferenceService(skillRepo, eventRepo)

			ctx := context.Background()

			event1 := &career.Event{
				ID:      "event-1",
				Text:    "Built REST API with Go and Docker containers",
				Date:    time.Date(2024, 6, 15, 0, 0, 0, 0, time.UTC),
				Company: "Acme Corp",
			}
			event2 := &career.Event{
				ID:      "event-2",
				Text:    "Designed PostgreSQL schema for user service",
				Date:    time.Date(2024, 7, 20, 0, 0, 0, 0, time.UTC),
				Company: "Acme Corp",
			}
			_ = eventRepo.Create(ctx, event1)
			_ = eventRepo.Create(ctx, event2)
			mockService.SetEvents([]*career.Event{event1, event2})

			_ = skillRepo.Create(ctx, &career.Skill{ID: "existing-go", Name: "Go", Category: "backend"})
			_ = skillRepo.Create(ctx, &career.Skill{ID: "existing-pg", Name: "PostgreSQL", Category: "database"})

			burst = &career.Burst{
				ID:          "burst-1",
				Name:        "API Development",
				Description: "Built REST API with Go and PostgreSQL",
				EventIDs:    []string{"event-1", "event-2"},
				Confirmed:   true,
			}

			intentCtx := &burstmgmt.IntentContext{
				Bursts:                []*career.Burst{burst},
				Service:               mockService,
				BurstRepository:       burstRepo,
				SkillInferenceService: skillInferenceService,
			}
			intentCtx.Validate()

			var err error
			intent, err = burstmgmt.NewIntent(intentCtx)
			Expect(err).NotTo(HaveOccurred())
			intent.Init()
		})

		It("should include existing skills in suggestions and report them in ExistingSkillNames after auto-triggered inference", func() {
			intent.HandleNavigate(&screens.NavigateResult{
				ResultData: burst,
			})

			facts := []*career.Fact{
				{ID: "fact-1", Text: "Expert in Go", SourceEventID: "event-1"},
			}
			batchCmd := intent.Update(burstmgmt.FactExtractionCompleteMsg{
				Facts: facts,
				Error: nil,
			})
			Expect(intent.GetState()).To(Equal(burstmgmt.StateInferringSkills))

			messages := executeSkillsBatchCmd(batchCmd)

			var suggestionsMsg burstmgmt.SkillSuggestionsLoadedMsg
			for _, msg := range messages {
				if sm, ok := msg.(burstmgmt.SkillSuggestionsLoadedMsg); ok {
					suggestionsMsg = sm
					break
				}
			}
			Expect(suggestionsMsg.Error).NotTo(HaveOccurred())

			suggestionNames := make([]string, 0, len(suggestionsMsg.Suggestions))
			for _, s := range suggestionsMsg.Suggestions {
				suggestionNames = append(suggestionNames, s.Name)
			}

			Expect(suggestionNames).To(ContainElement("Go"),
				"Go should appear in suggestions for burst context visibility")
			Expect(suggestionNames).To(ContainElement("PostgreSQL"),
				"PostgreSQL should appear in suggestions for burst context visibility")
			Expect(suggestionNames).To(ContainElement("Docker"),
				"Docker does not exist in the repo and should be suggested")

			Expect(suggestionsMsg.ExistingSkillNames).To(ContainElement("Go"))
			Expect(suggestionsMsg.ExistingSkillNames).To(ContainElement("PostgreSQL"))
		})
	})

	Describe("Trigger 3: Skills Management - 'i' on skills list", func() {
		var (
			skillsIntent *skillsmgmt.Intent
			skillRepo    *careermemory.SkillRepository
			eventRepo    *careermemory.EventRepository
		)

		BeforeEach(func() {
			skillRepo = careermemory.NewSkillRepository()
			eventRepo = careermemory.NewEventRepository()
			skillRepo.SetEventRepository(eventRepo)
			eventRepo.SetSkillRepository(skillRepo)
			skillInferenceService := skillinference.NewSkillInferenceService(skillRepo, eventRepo)

			ctx := context.Background()

			_ = eventRepo.Create(ctx, &career.Event{
				ID:      "e1",
				Text:    "Built REST API with Go and Docker containers",
				Date:    time.Date(2024, 6, 15, 0, 0, 0, 0, time.UTC),
				Company: "Acme Corp",
			})
			_ = eventRepo.Create(ctx, &career.Event{
				ID:      "e2",
				Text:    "Designed PostgreSQL schema for user service",
				Date:    time.Date(2024, 7, 20, 0, 0, 0, 0, time.UTC),
				Company: "Acme Corp",
			})

			_ = skillRepo.Create(ctx, &career.Skill{ID: "existing-go", Name: "Go", Category: "backend"})
			_ = skillRepo.Create(ctx, &career.Skill{ID: "existing-pg", Name: "PostgreSQL", Category: "database"})

			skillsCtx := skillsmgmt.NewIntentContext(ctx, skillRepo)
			skillsCtx.EventRepository = eventRepo
			skillsCtx.SkillInferenceService = skillInferenceService

			var err error
			skillsIntent, err = skillsmgmt.NewIntent(skillsCtx)
			Expect(err).NotTo(HaveOccurred())
			skillsIntent.Init()
		})

		It("should report existing skills in ExistingSkillNames and filter them in handler", func() {
			batchCmd := skillsIntent.Update(tea.KeyMsg{Type: tea.KeyRunes, Runes: []rune{'i'}})
			Expect(skillsIntent.GetState()).To(Equal(skillsmgmt.StateInferringSkills))

			messages := executeSkillsBatchCmd(batchCmd)

			var suggestionsMsg skillsmgmt.SkillSuggestionsLoadedMsg
			for _, msg := range messages {
				if sm, ok := msg.(skillsmgmt.SkillSuggestionsLoadedMsg); ok {
					suggestionsMsg = sm
					break
				}
			}
			Expect(suggestionsMsg.Error).NotTo(HaveOccurred())

			suggestionNames := make([]string, 0, len(suggestionsMsg.Suggestions))
			for _, s := range suggestionsMsg.Suggestions {
				suggestionNames = append(suggestionNames, s.Name)
			}

			Expect(suggestionNames).To(ContainElement("Go"),
				"Raw suggestions should include existing skills")
			Expect(suggestionNames).To(ContainElement("PostgreSQL"),
				"Raw suggestions should include existing skills")
			Expect(suggestionNames).To(ContainElement("Docker"),
				"Docker does not exist in the repo and should be suggested")

			Expect(suggestionsMsg.ExistingSkillNames).To(ContainElement("Go"))
			Expect(suggestionsMsg.ExistingSkillNames).To(ContainElement("PostgreSQL"))

			skillsIntent.Update(suggestionsMsg)
			Expect(skillsIntent.GetState()).To(Equal(skillsmgmt.StateSkillSuggestionReview))
			view := skillsIntent.View()
			Expect(view).To(ContainSubstring("Docker"))
		})
	})
})
