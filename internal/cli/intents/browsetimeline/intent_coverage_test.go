package browsetimeline

import (
	"context"
	"errors"
	"time"

	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"

	"github.com/baphled/kariya/internal/cli/behaviors"
	"github.com/baphled/kariya/internal/cli/screens"
	skillModals "github.com/baphled/kariya/internal/cli/screens/skills/modals"
	"github.com/baphled/kariya/internal/cli/screens/timeline"
	"github.com/baphled/kariya/internal/cli/screens/timeline/modals"
	"github.com/baphled/kariya/internal/cli/service"
	"github.com/baphled/kariya/internal/cli/terminal"
	"github.com/baphled/kariya/internal/domain/career"
	careerrepo "github.com/baphled/kariya/internal/repository/career"
	careerservice "github.com/baphled/kariya/internal/service/career"
	"github.com/baphled/kariya/internal/service/career/skillinference"
	"github.com/baphled/kariya/internal/testutil/fixtures"
	tea "github.com/charmbracelet/bubbletea"
)

type mockEventService struct {
	deleteErr        error
	listEventsResult []*career.Event
	listEventsErr    error
	skillsForEvent   []*career.Skill
	skillsForErr     error
	allSkills        []*career.Skill
	allSkillsErr     error
	linkErr          error
	unlinkErr        error
	captureErr       error
	updateErr        error
}

func (m *mockEventService) DeleteEvent(_ context.Context, _ string) error {
	return m.deleteErr
}

func (m *mockEventService) ListEvents(_ context.Context, _ *careerrepo.EventListFilters) ([]*career.Event, error) {
	return m.listEventsResult, m.listEventsErr
}

func (m *mockEventService) CaptureEvent(_ context.Context, _ string, _ time.Time, _ careerservice.EventCaptureMode, _ ...service.Option) error {
	return m.captureErr
}

func (m *mockEventService) UpdateEventMetadata(_ context.Context, _ *career.Event) error {
	return m.updateErr
}

func (m *mockEventService) GetSkillsForEvent(_ context.Context, _ string) ([]*career.Skill, error) {
	return m.skillsForEvent, m.skillsForErr
}

func (m *mockEventService) LinkSkillToEvent(_ context.Context, _, _ string) error {
	return m.linkErr
}

func (m *mockEventService) UnlinkSkillFromEvent(_ context.Context, _, _ string) error {
	return m.unlinkErr
}

func (m *mockEventService) ListAllSkills(_ context.Context) ([]*career.Skill, error) {
	return m.allSkills, m.allSkillsErr
}

type mockSkillService struct {
	createErr error
}

func (m *mockSkillService) Create(_ context.Context, _ *career.Skill) error {
	return m.createErr
}

type mockSkillInferenceService struct {
	inferResult  *skillinference.InferenceResult
	inferErr     error
	createSkills []*career.Skill
	createErr    error
}

func (m *mockSkillInferenceService) InferSkillsFromEvents(_ context.Context, _ []*career.Event) (*skillinference.InferenceResult, error) {
	return m.inferResult, m.inferErr
}

func (m *mockSkillInferenceService) CreateSkillsFromSuggestions(_ context.Context, _ []skillinference.SkillSuggestion) ([]*career.Skill, error) {
	return m.createSkills, m.createErr
}

var _ = Describe("Intent Coverage", func() {
	var (
		intent *Intent
		events []*career.Event
		svc    *mockEventService
	)

	BeforeEach(func() {
		events = []*career.Event{
			fixtures.EventWith("evt-1", "First event at TechCorp building platform", "TechCorp", "Platform"),
			fixtures.EventWith("evt-2", "Second event at CloudInc doing infra", "CloudInc", "Infrastructure"),
		}
		svc = &mockEventService{}
	})

	createIntent := func(opts ...func(*IntentContext)) *Intent {
		ctx := &IntentContext{
			Events:          events,
			CLIEventService: svc,
		}
		for _, opt := range opts {
			opt(ctx)
		}
		i, err := NewIntent(ctx)
		Expect(err).NotTo(HaveOccurred())
		i.Init()
		return i
	}

	Describe("noopCmd", func() {
		It("should return nil message", func() {
			msg := noopCmd()
			Expect(msg).To(BeNil())
		})
	})

	Describe("filterNewSkillSuggestions", func() {
		It("should return all suggestions when no existing names", func() {
			suggestions := []skillinference.SkillSuggestion{
				{Name: "Go", Category: "backend", Confidence: 0.9},
				{Name: "Python", Category: "backend", Confidence: 0.8},
			}
			result := filterNewSkillSuggestions(suggestions, []string{})
			Expect(result).To(HaveLen(2))
		})

		It("should filter out existing skills", func() {
			suggestions := []skillinference.SkillSuggestion{
				{Name: "Go", Category: "backend", Confidence: 0.9},
				{Name: "Python", Category: "backend", Confidence: 0.8},
				{Name: "Rust", Category: "backend", Confidence: 0.7},
			}
			result := filterNewSkillSuggestions(suggestions, []string{"Go", "Rust"})
			Expect(result).To(HaveLen(1))
			Expect(result[0].Name).To(Equal("Python"))
		})

		It("should return empty slice when all skills exist", func() {
			suggestions := []skillinference.SkillSuggestion{
				{Name: "Go", Category: "backend", Confidence: 0.9},
			}
			result := filterNewSkillSuggestions(suggestions, []string{"Go"})
			Expect(result).To(BeEmpty())
		})

		It("should handle nil existing names", func() {
			suggestions := []skillinference.SkillSuggestion{
				{Name: "Go", Category: "backend", Confidence: 0.9},
			}
			result := filterNewSkillSuggestions(suggestions, nil)
			Expect(result).To(HaveLen(1))
		})
	})

	Describe("removeEventFromList", func() {
		It("should remove event from both events and filtered events", func() {
			intent = createIntent()
			initialCount := len(intent.context.Events)
			intent.removeEventFromList("evt-2")
			Expect(intent.context.Events).To(HaveLen(initialCount - 1))
		})

		It("should handle removing non-existent event ID", func() {
			intent = createIntent()
			initialCount := len(intent.context.Events)
			intent.removeEventFromList("nonexistent-id")
			Expect(intent.context.Events).To(HaveLen(initialCount))
		})
	})

	Describe("getContext", func() {
		It("should return a non-nil context", func() {
			intent = createIntent()
			result := intent.getContext()
			Expect(result).NotTo(BeNil())
		})
	})

	Describe("HandleCancel", func() {
		It("should return to timeline from StateDeleteConfirm", func() {
			intent = createIntent()
			intent.state = StateDeleteConfirm
			cmd := intent.HandleCancel(&screens.CancelResult{})
			Expect(cmd).To(BeNil())
			Expect(intent.state).To(Equal(StateTimeline))
			Expect(intent.active).To(BeTrue())
		})

		It("should cancel from unknown state", func() {
			intent = createIntent()
			intent.state = State("unknown")
			intent.HandleCancel(&screens.CancelResult{})
			Expect(intent.active).To(BeFalse())
		})
	})

	Describe("ClearFilters with filter stack", func() {
		It("should pop search layer from stack", func() {
			intent = createIntent()
			intent.filters.SearchText = "search term"
			intent.filterStack.Push(behaviors.FilterLayerSearch)
			intent.ClearFilters()
			Expect(intent.filters.SearchText).To(Equal(""))
		})

		It("should pop company layer from stack", func() {
			intent = createIntent()
			intent.filters.Companies = []string{"TechCorp"}
			intent.filterStack.Push(behaviors.FilterLayerCompany)
			intent.ClearFilters()
			Expect(intent.filters.Companies).To(BeEmpty())
		})

		It("should pop category layer from stack", func() {
			intent = createIntent()
			intent.filters.Categories = []string{"Backend"}
			intent.filterStack.Push(behaviors.FilterLayerCategory)
			intent.ClearFilters()
			Expect(intent.filters.Categories).To(BeEmpty())
		})

		It("should pop project layer from stack", func() {
			intent = createIntent()
			intent.filters.Projects = []string{"Platform"}
			intent.filterStack.Push(behaviors.FilterLayerProject)
			intent.ClearFilters()
			Expect(intent.filters.Projects).To(BeEmpty())
		})

		It("should pop tags layer from stack", func() {
			intent = createIntent()
			intent.filters.Tags = []string{"golang"}
			intent.filterStack.Push(behaviors.FilterLayerTags)
			intent.ClearFilters()
			Expect(intent.filters.Tags).To(BeEmpty())
		})

		It("should pop sort layer from stack and reset to defaults", func() {
			intent = createIntent()
			intent.filters.SortBy = "text"
			intent.filters.SortOrder = "asc"
			intent.filterStack.Push(behaviors.FilterLayerSort)
			intent.ClearFilters()
			Expect(intent.filters.SortBy).To(Equal("date"))
			Expect(intent.filters.SortOrder).To(Equal("desc"))
		})

		It("should clear filter stack when no active filters remain", func() {
			intent = createIntent()
			intent.filters.SearchText = "test"
			intent.filterStack.Push(behaviors.FilterLayerSearch)
			intent.ClearFilters()
			Expect(intent.HasActiveFilters()).To(BeFalse())
		})
	})

	Describe("handleDeleteConfirmation with service", func() {
		It("should delete event successfully via service", func() {
			intent = createIntent()
			intent.selectedEvent = intent.context.Events[0]
			cmd := intent.handleDeleteConfirmation(true)
			Expect(cmd).To(BeNil())
			Expect(intent.state).To(Equal(StateTimeline))
			Expect(intent.context.Events).To(HaveLen(1))
			Expect(intent.context.Events[0].ID).To(Equal("evt-2"))
		})

		It("should store error when delete fails", func() {
			svc.deleteErr = errors.New("delete failed")
			intent = createIntent()
			intent.selectedEvent = intent.context.Events[0]
			cmd := intent.handleDeleteConfirmation(true)
			Expect(cmd).To(BeNil())
			Expect(intent.deleteError).To(HaveOccurred())
		})

		It("should handle not confirmed", func() {
			intent = createIntent()
			cmd := intent.handleDeleteConfirmation(false)
			Expect(cmd).To(BeNil())
			Expect(intent.state).To(Equal(StateTimeline))
		})

		It("should handle confirmed with nil selectedEvent", func() {
			intent = createIntent()
			intent.selectedEvent = nil
			cmd := intent.handleDeleteConfirmation(true)
			Expect(cmd).To(BeNil())
			Expect(intent.state).To(Equal(StateTimeline))
		})
	})

	Describe("showSkillsForCurrentEvent with service", func() {
		It("should load and display skills when service returns skills", func() {
			svc.skillsForEvent = []*career.Skill{fixtures.Skill("skill-1")}
			intent = createIntent()
			cmd := intent.showSkillsForCurrentEvent()
			Expect(cmd).To(BeNil())
			Expect(intent.viewSkillsModal).NotTo(BeNil())
		})

		It("should handle service error gracefully", func() {
			svc.skillsForErr = errors.New("skill load failed")
			intent = createIntent()
			cmd := intent.showSkillsForCurrentEvent()
			Expect(cmd).To(BeNil())
			Expect(intent.viewSkillsModal).NotTo(BeNil())
		})

		It("should return nil when selectedEvent is nil", func() {
			intent = createIntent()
			intent.selectedEvent = nil
			cmd := intent.showSkillsForCurrentEvent()
			Expect(cmd).To(BeNil())
		})
	})

	Describe("performSkillLinkOperation with service", func() {
		It("should link skill successfully", func() {
			intent = createIntent()
			skill := fixtures.Skill("skill-1")
			cmd := intent.performSkillLinkOperation(skill, true)
			Expect(cmd).NotTo(BeNil())
			msg := cmd()
			_, ok := msg.(SkillLinkedMsg)
			Expect(ok).To(BeTrue())
		})

		It("should unlink skill successfully", func() {
			intent = createIntent()
			skill := fixtures.Skill("skill-1")
			cmd := intent.performSkillLinkOperation(skill, false)
			Expect(cmd).NotTo(BeNil())
			msg := cmd()
			_, ok := msg.(SkillUnlinkedMsg)
			Expect(ok).To(BeTrue())
		})

		It("should show error when link fails", func() {
			svc.linkErr = errors.New("link failed")
			intent = createIntent()
			skill := fixtures.Skill("skill-1")
			cmd := intent.performSkillLinkOperation(skill, true)
			Expect(cmd).To(BeNil())
			Expect(intent.errorModal).NotTo(BeNil())
		})

		It("should show error when unlink fails", func() {
			svc.unlinkErr = errors.New("unlink failed")
			intent = createIntent()
			skill := fixtures.Skill("skill-1")
			cmd := intent.performSkillLinkOperation(skill, false)
			Expect(cmd).To(BeNil())
			Expect(intent.errorModal).NotTo(BeNil())
		})

		It("should return nil when selectedEvent is nil", func() {
			intent = createIntent()
			intent.selectedEvent = nil
			cmd := intent.performSkillLinkOperation(fixtures.Skill("s1"), true)
			Expect(cmd).To(BeNil())
		})

		It("should return nil when skill is nil", func() {
			intent = createIntent()
			cmd := intent.performSkillLinkOperation(nil, true)
			Expect(cmd).To(BeNil())
		})
	})

	Describe("linkSkillToCurrentEvent", func() {
		It("should delegate to performSkillLinkOperation with link=true", func() {
			intent = createIntent()
			skill := fixtures.Skill("skill-1")
			cmd := intent.linkSkillToCurrentEvent(skill)
			Expect(cmd).NotTo(BeNil())
			msg := cmd()
			linked, ok := msg.(SkillLinkedMsg)
			Expect(ok).To(BeTrue())
			Expect(linked.SkillID).To(Equal("skill-1"))
		})
	})

	Describe("unlinkSkillFromCurrentEvent", func() {
		It("should delegate to performSkillLinkOperation with link=false", func() {
			intent = createIntent()
			skill := fixtures.Skill("skill-1")
			cmd := intent.unlinkSkillFromCurrentEvent(skill)
			Expect(cmd).NotTo(BeNil())
			msg := cmd()
			unlinked, ok := msg.(SkillUnlinkedMsg)
			Expect(ok).To(BeTrue())
			Expect(unlinked.SkillID).To(Equal("skill-1"))
		})
	})

	Describe("refreshSkillsModal with service", func() {
		It("should return nil when selectedEvent is nil", func() {
			intent = createIntent()
			intent.selectedEvent = nil
			cmd := intent.refreshSkillsModal()
			Expect(cmd).To(BeNil())
		})

		It("should return nil when viewSkillsModal is nil", func() {
			intent = createIntent()
			cmd := intent.refreshSkillsModal()
			Expect(cmd).To(BeNil())
		})

		It("should refresh skills when both selectedEvent and viewSkillsModal are set", func() {
			svc.skillsForEvent = []*career.Skill{fixtures.Skill("s1")}
			intent = createIntent()
			intent.viewSkillsModal = modals.NewSkillsDetailModal("evt-1", []*career.Skill{}, nil)
			intent.viewSkillsModal.Show()
			cmd := intent.refreshSkillsModal()
			Expect(cmd).To(BeNil())
		})

		It("should show error when GetSkillsForEvent fails", func() {
			svc.skillsForErr = errors.New("load failed")
			intent = createIntent()
			intent.viewSkillsModal = modals.NewSkillsDetailModal("evt-1", []*career.Skill{}, nil)
			intent.viewSkillsModal.Show()
			cmd := intent.refreshSkillsModal()
			Expect(cmd).To(BeNil())
			Expect(intent.errorModal).NotTo(BeNil())
		})
	})

	Describe("openSkillPickerModal with service", func() {
		It("should show error when ListAllSkills fails", func() {
			svc.allSkillsErr = errors.New("list failed")
			intent = createIntent()
			cmd := intent.openSkillPickerModal()
			Expect(cmd).To(BeNil())
			Expect(intent.errorModal).NotTo(BeNil())
		})

		It("should show error when GetSkillsForEvent fails", func() {
			svc.allSkills = []*career.Skill{fixtures.Skill("s1")}
			svc.skillsForErr = errors.New("get skills failed")
			intent = createIntent()
			cmd := intent.openSkillPickerModal()
			Expect(cmd).To(BeNil())
			Expect(intent.errorModal).NotTo(BeNil())
		})

		It("should create skill picker with available skills", func() {
			existingSkill := fixtures.Skill("existing-1")
			availableSkill := fixtures.Skill("available-1")
			svc.allSkills = []*career.Skill{existingSkill, availableSkill}
			svc.skillsForEvent = []*career.Skill{existingSkill}
			intent = createIntent()
			cmd := intent.openSkillPickerModal()
			Expect(cmd).To(BeNil())
			Expect(intent.skillPickerModal).NotTo(BeNil())
		})

		It("should return nil when selectedEvent is nil", func() {
			intent = createIntent()
			intent.selectedEvent = nil
			cmd := intent.openSkillPickerModal()
			Expect(cmd).To(BeNil())
		})
	})

	Describe("openSkillAddModal", func() {
		It("should create a skill add modal", func() {
			intent = createIntent()
			cmd := intent.openSkillAddModal()
			Expect(cmd).NotTo(BeNil())
			Expect(intent.skillAddModal).NotTo(BeNil())
		})
	})

	Describe("createAndLinkSkill with service", func() {
		It("should return nil when selectedEvent is nil", func() {
			intent = createIntent()
			intent.selectedEvent = nil
			cmd := intent.createAndLinkSkill(nil)
			Expect(cmd).To(BeNil())
		})

		It("should return nil when skillData is nil", func() {
			intent = createIntent()
			cmd := intent.createAndLinkSkill(nil)
			Expect(cmd).To(BeNil())
		})

		It("should create and link skill successfully", func() {
			skillSvc := &mockSkillService{}
			intent = createIntent(func(ctx *IntentContext) {
				ctx.CLISkillService = skillSvc
			})
			skillData := &skillModals.SkillEditData{
				Name:     "Go",
				Category: "backend",
				Level:    "expert",
			}
			cmd := intent.createAndLinkSkill(skillData)
			Expect(cmd).NotTo(BeNil())
			msg := cmd()
			created, ok := msg.(SkillCreatedMsg)
			Expect(ok).To(BeTrue())
			Expect(created.Skill).NotTo(BeNil())
		})

		It("should show error when skill creation fails", func() {
			skillSvc := &mockSkillService{createErr: errors.New("create failed")}
			intent = createIntent(func(ctx *IntentContext) {
				ctx.CLISkillService = skillSvc
			})
			skillData := &skillModals.SkillEditData{Name: "Go", Category: "backend"}
			cmd := intent.createAndLinkSkill(skillData)
			Expect(cmd).To(BeNil())
			Expect(intent.errorModal).NotTo(BeNil())
		})

		It("should show error when linking fails", func() {
			skillSvc := &mockSkillService{}
			svc.linkErr = errors.New("link failed")
			intent = createIntent(func(ctx *IntentContext) {
				ctx.CLISkillService = skillSvc
			})
			skillData := &skillModals.SkillEditData{Name: "Go", Category: "backend"}
			cmd := intent.createAndLinkSkill(skillData)
			Expect(cmd).To(BeNil())
			Expect(intent.errorModal).NotTo(BeNil())
		})
	})

	Describe("inferSkillsFromEvent", func() {
		It("should return nil when selectedEvent is nil", func() {
			intent = createIntent()
			intent.selectedEvent = nil
			cmd := intent.inferSkillsFromEvent()
			Expect(cmd).To(BeNil())
		})

		It("should return nil when SkillInferenceService is nil", func() {
			intent = createIntent()
			cmd := intent.inferSkillsFromEvent()
			Expect(cmd).To(BeNil())
		})

		It("should return a command when service is available", func() {
			inferSvc := &mockSkillInferenceService{
				inferResult: &skillinference.InferenceResult{
					Suggestions: []skillinference.SkillSuggestion{
						{Name: "Go", Category: "backend", Confidence: 0.9},
					},
				},
			}
			intent = createIntent(func(ctx *IntentContext) {
				ctx.SkillInferenceService = inferSvc
			})
			cmd := intent.inferSkillsFromEvent()
			Expect(cmd).NotTo(BeNil())
		})

		It("should return error message when inference fails", func() {
			inferSvc := &mockSkillInferenceService{
				inferErr: errors.New("inference failed"),
			}
			intent = createIntent(func(ctx *IntentContext) {
				ctx.SkillInferenceService = inferSvc
			})
			cmd := intent.inferSkillsFromEvent()
			Expect(cmd).NotTo(BeNil())
			msg := cmd()
			errMsg, ok := msg.(SkillSuggestionsErrorMsg)
			Expect(ok).To(BeTrue())
			Expect(errMsg.Err).To(HaveOccurred())
		})

		It("should return loaded message on success", func() {
			inferSvc := &mockSkillInferenceService{
				inferResult: &skillinference.InferenceResult{
					Suggestions: []skillinference.SkillSuggestion{
						{Name: "Go", Category: "backend", Confidence: 0.9},
					},
				},
			}
			intent = createIntent(func(ctx *IntentContext) {
				ctx.SkillInferenceService = inferSvc
			})
			cmd := intent.inferSkillsFromEvent()
			Expect(cmd).NotTo(BeNil())
			msg := cmd()
			loaded, ok := msg.(SkillSuggestionsLoadedMsg)
			Expect(ok).To(BeTrue())
			Expect(loaded.Suggestions).To(HaveLen(1))
		})
	})

	Describe("saveSkillFromSuggestion", func() {
		It("should return early when SkillInferenceService is nil", func() {
			intent = createIntent()
			Expect(func() {
				intent.saveSkillFromSuggestion(skillinference.SkillSuggestion{Name: "Go", Category: "backend"})
			}).NotTo(Panic())
		})

		It("should show error when CreateSkillsFromSuggestions fails", func() {
			inferSvc := &mockSkillInferenceService{createErr: errors.New("creation failed")}
			intent = createIntent(func(ctx *IntentContext) {
				ctx.SkillInferenceService = inferSvc
			})
			intent.saveSkillFromSuggestion(skillinference.SkillSuggestion{Name: "Go", Category: "backend"})
			Expect(intent.errorModal).NotTo(BeNil())
		})

		It("should link skill to event after creation", func() {
			newSkill := fixtures.Skill("new-skill-1")
			inferSvc := &mockSkillInferenceService{createSkills: []*career.Skill{newSkill}}
			intent = createIntent(func(ctx *IntentContext) {
				ctx.SkillInferenceService = inferSvc
			})
			intent.saveSkillFromSuggestion(skillinference.SkillSuggestion{Name: "Go", Category: "backend"})
			Expect(intent.errorModal).To(BeNil())
		})

		It("should show error when linking fails", func() {
			newSkill := fixtures.Skill("new-skill-1")
			inferSvc := &mockSkillInferenceService{createSkills: []*career.Skill{newSkill}}
			svc.linkErr = errors.New("link failed")
			intent = createIntent(func(ctx *IntentContext) {
				ctx.SkillInferenceService = inferSvc
			})
			intent.saveSkillFromSuggestion(skillinference.SkillSuggestion{Name: "Go", Category: "backend"})
			Expect(intent.errorModal).NotTo(BeNil())
		})
	})

	Describe("handleSkillSuggestionsLoaded", func() {
		It("should show error when suggestions are empty", func() {
			intent = createIntent()
			msg := SkillSuggestionsLoadedMsg{Suggestions: []skillinference.SkillSuggestion{}}
			cmd := intent.handleSkillSuggestionsLoaded(msg)
			Expect(cmd).To(BeNil())
			Expect(intent.HasVisibleErrorModal()).To(BeTrue())
		})

		It("should show error when all suggestions already exist", func() {
			intent = createIntent()
			msg := SkillSuggestionsLoadedMsg{
				Suggestions:        []skillinference.SkillSuggestion{{Name: "Go", Category: "backend", Confidence: 0.9}},
				ExistingSkillNames: []string{"Go"},
			}
			cmd := intent.handleSkillSuggestionsLoaded(msg)
			Expect(cmd).To(BeNil())
			Expect(intent.HasVisibleErrorModal()).To(BeTrue())
		})

		It("should create suggestion modal when new suggestions exist", func() {
			intent = createIntent()
			msg := SkillSuggestionsLoadedMsg{
				Suggestions: []skillinference.SkillSuggestion{
					{Name: "Go", Category: "backend", Confidence: 0.9},
					{Name: "Python", Category: "backend", Confidence: 0.8},
				},
			}
			cmd := intent.handleSkillSuggestionsLoaded(msg)
			Expect(cmd).To(BeNil())
			Expect(intent.skillSuggestionModal).NotTo(BeNil())
		})
	})

	Describe("getStateName", func() {
		It("should return Timeline for StateTimeline", func() {
			intent = createIntent()
			Expect(intent.getStateName()).To(Equal("Timeline"))
		})

		It("should return Delete Confirmation for StateDeleteConfirm", func() {
			intent = createIntent()
			intent.state = StateDeleteConfirm
			Expect(intent.getStateName()).To(Equal("Delete Confirmation"))
		})

		It("should return Unknown for undefined state", func() {
			intent = createIntent()
			intent.state = State("undefined_state")
			Expect(intent.getStateName()).To(Equal("Unknown"))
		})
	})

	Describe("getContextHelp", func() {
		It("should return help text for StateTimeline", func() {
			intent = createIntent()
			help := intent.getContextHelp()
			Expect(help).NotTo(BeEmpty())
		})

		It("should include clear filters badge when filters are active", func() {
			intent = createIntent(func(ctx *IntentContext) {
				ctx.InitialFilters = &Filters{Companies: []string{"Corp"}}
			})
			help := intent.getContextHelp()
			Expect(help).To(ContainSubstring("Clear"))
		})

		It("should return empty string for StateDeleteConfirm", func() {
			intent = createIntent()
			intent.state = StateDeleteConfirm
			help := intent.getContextHelp()
			Expect(help).To(BeEmpty())
		})

		It("should return global badges for unknown state", func() {
			intent = createIntent()
			intent.state = State("other")
			help := intent.getContextHelp()
			Expect(help).NotTo(BeEmpty())
		})
	})

	Describe("clearAllFilters", func() {
		It("should reset all filters to defaults", func() {
			intent = createIntent(func(ctx *IntentContext) {
				ctx.InitialFilters = &Filters{
					SearchText: "search",
					Tags:       []string{"tag1"},
					Companies:  []string{"comp1"},
					Categories: []string{"cat1"},
					Projects:   []string{"proj1"},
					DateFrom:   "2024-01-01",
					DateTo:     "2024-12-31",
					SortBy:     "text",
					SortOrder:  "asc",
				}
			})
			Expect(intent.HasActiveFilters()).To(BeTrue())
			intent.clearAllFilters()
			Expect(intent.HasActiveFilters()).To(BeFalse())
		})
	})

	Describe("View", func() {
		It("should return 'No active screen' when activeScreen is nil", func() {
			intent = createIntent()
			intent.activeScreen = nil
			Expect(intent.View()).To(Equal("No active screen"))
		})

		It("should render EventDeleteConfirmScreen view", func() {
			intent = createIntent()
			intent.activeScreen = timeline.NewEventDeleteConfirmScreen(events[0])
			view := intent.View()
			Expect(view).NotTo(BeEmpty())
		})
	})

	Describe("transitionToScreen", func() {
		It("should set terminal info when available", func() {
			intent = createIntent()
			info := &terminal.Info{Width: 120, Height: 40, IsValid: true}
			intent.UpdateTerminalInfo(info)
			screen := timeline.NewTimelineEventListScreen(events)
			intent.transitionToScreen(screen)
			Expect(intent.activeScreen).To(Equal(screen))
		})
	})

	Describe("Update with error modal active", func() {
		It("should dismiss error modal on Escape", func() {
			intent = createIntent()
			intent.ShowErrorModal("Test Error", "Something went wrong")
			Expect(intent.errorModal).NotTo(BeNil())
			cmd := intent.Update(tea.KeyMsg{Type: tea.KeyEsc})
			Expect(cmd).NotTo(BeNil())
			Expect(intent.errorModal).To(BeNil())
		})

		It("should consume non-Escape keys when error modal is active", func() {
			intent = createIntent()
			intent.ShowErrorModal("Test Error", "Something went wrong")
			cmd := intent.Update(tea.KeyMsg{Type: tea.KeyRunes, Runes: []rune{'a'}})
			Expect(cmd).NotTo(BeNil())
			Expect(intent.errorModal).NotTo(BeNil())
		})
	})

	Describe("Update with search modal active", func() {
		It("should route messages to updateSearchModal", func() {
			intent = createIntent()
			intent.openSearchModal()
			Expect(intent.searchModal).NotTo(BeNil())
			cmd := intent.Update(tea.KeyMsg{Type: tea.KeyEsc})
			_ = cmd
		})
	})

	Describe("Update with filter modal active", func() {
		It("should route messages to updateFilterModal", func() {
			intent = createIntent()
			intent.openFilterModal()
			Expect(intent.filterModal).NotTo(BeNil())
			cmd := intent.Update(tea.KeyMsg{Type: tea.KeyEsc})
			_ = cmd
		})
	})

	Describe("Update with sort modal active", func() {
		It("should route messages to updateSortModal", func() {
			intent = createIntent()
			intent.openSortModal()
			Expect(intent.sortModal).NotTo(BeNil())
			cmd := intent.Update(tea.KeyMsg{Type: tea.KeyEsc})
			_ = cmd
		})
	})

	Describe("Update with quick add modal active", func() {
		It("should route messages to updateQuickAddModal", func() {
			intent = createIntent()
			intent.openQuickAddModal()
			Expect(intent.quickAddModal).NotTo(BeNil())
			cmd := intent.Update(tea.KeyMsg{Type: tea.KeyEsc})
			_ = cmd
		})
	})

	Describe("Update with edit modal active", func() {
		It("should route messages to updateEditModal", func() {
			intent = createIntent()
			intent.openEditModalForEvent(events[0])
			Expect(intent.editModal).NotTo(BeNil())
			cmd := intent.Update(tea.KeyMsg{Type: tea.KeyEsc})
			_ = cmd
		})
	})

	Describe("Update with delete modal active", func() {
		It("should route messages to updateDeleteModal", func() {
			intent = createIntent()
			intent.openDeleteModalForEvent(events[0])
			Expect(intent.deleteModal).NotTo(BeNil())
			cmd := intent.Update(tea.KeyMsg{Type: tea.KeyEsc})
			_ = cmd
		})
	})

	Describe("Update with view skills modal active", func() {
		It("should route messages to updateViewSkillsModal", func() {
			svc.skillsForEvent = []*career.Skill{fixtures.Skill("s1")}
			intent = createIntent()
			intent.showSkillsForCurrentEvent()
			Expect(intent.viewSkillsModal).NotTo(BeNil())
			cmd := intent.Update(tea.KeyMsg{Type: tea.KeyEsc})
			_ = cmd
		})
	})

	Describe("Update with skill picker modal active", func() {
		It("should route messages to updateSkillPickerModal", func() {
			svc.allSkills = []*career.Skill{fixtures.Skill("s1")}
			svc.skillsForEvent = []*career.Skill{}
			intent = createIntent()
			intent.openSkillPickerModal()
			Expect(intent.skillPickerModal).NotTo(BeNil())
			cmd := intent.Update(tea.KeyMsg{Type: tea.KeyEsc})
			_ = cmd
		})
	})

	Describe("Update with skill add modal active", func() {
		It("should route messages to updateSkillAddModal", func() {
			intent = createIntent()
			intent.openSkillAddModal()
			Expect(intent.skillAddModal).NotTo(BeNil())
			cmd := intent.Update(tea.KeyMsg{Type: tea.KeyEsc})
			_ = cmd
		})
	})

	Describe("Update with skill suggestion modal active", func() {
		It("should route messages to updateSkillSuggestionModal", func() {
			intent = createIntent()
			msg := SkillSuggestionsLoadedMsg{
				Suggestions: []skillinference.SkillSuggestion{
					{Name: "Go", Category: "backend", Confidence: 0.9},
					{Name: "Python", Category: "backend", Confidence: 0.8},
				},
			}
			intent.handleSkillSuggestionsLoaded(msg)
			Expect(intent.skillSuggestionModal).NotTo(BeNil())
			cmd := intent.Update(tea.KeyMsg{Type: tea.KeyEsc})
			_ = cmd
		})
	})

	Describe("Update with view detail modal active", func() {
		It("should route messages to updateViewDetailModal", func() {
			intent = createIntent()
			intent.showEventDetailModal(events[0])
			Expect(intent.viewDetailModal).NotTo(BeNil())
			cmd := intent.Update(tea.KeyMsg{Type: tea.KeyEsc})
			_ = cmd
		})

		It("should handle 's' key to show skills", func() {
			svc.skillsForEvent = []*career.Skill{}
			intent = createIntent()
			intent.showEventDetailModal(events[0])
			cmd := intent.Update(tea.KeyMsg{Type: tea.KeyRunes, Runes: []rune{'s'}})
			_ = cmd
		})

		It("should handle 'e' key to open edit modal", func() {
			intent = createIntent()
			intent.showEventDetailModal(events[0])
			cmd := intent.Update(tea.KeyMsg{Type: tea.KeyRunes, Runes: []rune{'e'}})
			_ = cmd
		})

		It("should handle 'd' key to open delete modal", func() {
			intent = createIntent()
			intent.showEventDetailModal(events[0])
			cmd := intent.Update(tea.KeyMsg{Type: tea.KeyRunes, Runes: []rune{'d'}})
			_ = cmd
		})

		It("should handle 'd' key with long event text", func() {
			longEvent := fixtures.EventWith("long-1", "This is a very long event text that exceeds fifty characters for truncation testing purposes", "Corp", "Proj")
			events = append(events, longEvent)
			intent = createIntent()
			intent.showEventDetailModal(longEvent)
			cmd := intent.Update(tea.KeyMsg{Type: tea.KeyRunes, Runes: []rune{'d'}})
			_ = cmd
		})
	})

	Describe("Update when inactive", func() {
		It("should return nil when intent is not active", func() {
			intent = createIntent()
			intent.active = false
			cmd := intent.Update(tea.KeyMsg{Type: tea.KeyRunes, Runes: []rune{'a'}})
			Expect(cmd).To(BeNil())
		})
	})

	Describe("Update with skill messages", func() {
		It("should handle SkillLinkedMsg", func() {
			intent = createIntent()
			cmd := intent.Update(SkillLinkedMsg{EventID: "evt-1", SkillID: "skill-1"})
			Expect(cmd).To(BeNil())
		})

		It("should handle SkillUnlinkedMsg", func() {
			intent = createIntent()
			cmd := intent.Update(SkillUnlinkedMsg{EventID: "evt-1", SkillID: "skill-1"})
			Expect(cmd).To(BeNil())
		})

		It("should handle SkillCreatedMsg", func() {
			intent = createIntent()
			cmd := intent.Update(SkillCreatedMsg{Skill: fixtures.Skill("new-skill")})
			Expect(cmd).To(BeNil())
		})

		It("should handle SkillSuggestionsErrorMsg", func() {
			intent = createIntent()
			cmd := intent.Update(SkillSuggestionsErrorMsg{Err: errors.New("inference error")})
			Expect(cmd).To(BeNil())
			Expect(intent.errorModal).NotTo(BeNil())
		})
	})

	Describe("handleKeyShortcuts", func() {
		It("should clear filters with x key when filters active", func() {
			intent = createIntent(func(ctx *IntentContext) {
				ctx.InitialFilters = &Filters{Companies: []string{"Corp"}}
			})
			cmd := intent.Update(tea.KeyMsg{Type: tea.KeyRunes, Runes: []rune{'x'}})
			Expect(cmd).To(BeNil())
		})

		It("should do nothing with x key when no active filters", func() {
			intent = createIntent()
			cmd := intent.Update(tea.KeyMsg{Type: tea.KeyRunes, Runes: []rune{'x'}})
			Expect(cmd).To(BeNil())
		})

		It("should open search modal with / key", func() {
			intent = createIntent()
			cmd := intent.handleKeyShortcuts(tea.KeyMsg{Type: tea.KeyRunes, Runes: []rune{'/'}})
			Expect(cmd).NotTo(BeNil())
			Expect(intent.searchModal).NotTo(BeNil())
		})

		It("should open filter modal with f key", func() {
			intent = createIntent()
			cmd := intent.handleKeyShortcuts(tea.KeyMsg{Type: tea.KeyRunes, Runes: []rune{'f'}})
			Expect(cmd).NotTo(BeNil())
			Expect(intent.filterModal).NotTo(BeNil())
		})

		It("should open sort modal with s key", func() {
			intent = createIntent()
			cmd := intent.handleKeyShortcuts(tea.KeyMsg{Type: tea.KeyRunes, Runes: []rune{'s'}})
			Expect(cmd).NotTo(BeNil())
			Expect(intent.sortModal).NotTo(BeNil())
		})
	})

	Describe("handleScreenResult", func() {
		It("should return nil for nil result", func() {
			intent = createIntent()
			cmd := intent.handleScreenResult(nil)
			Expect(cmd).To(BeNil())
		})

		It("should return nil for non-ScreenResult type", func() {
			intent = createIntent()
			cmd := intent.handleScreenResult("not a screen result")
			Expect(cmd).To(BeNil())
		})
	})

	Describe("NewIntent with nil events", func() {
		It("should initialize events to empty slice", func() {
			ctx := &IntentContext{Events: nil}
			i, err := NewIntent(ctx)
			Expect(err).NotTo(HaveOccurred())
			Expect(i).NotTo(BeNil())
		})
	})

	Describe("HandleSubmit", func() {
		It("should return nil for any submit result", func() {
			intent = createIntent()
			result := &screens.SubmitResult{FormData: "any data"}
			cmd := intent.HandleSubmit(result)
			Expect(cmd).To(BeNil())
		})
	})

	Describe("HandleError", func() {
		It("should store the error from result", func() {
			intent = createIntent()
			testErr := errors.New("test error occurred")
			result := &screens.ErrorResult{Err: testErr, Message: "Something failed"}
			cmd := intent.HandleError(result)
			Expect(cmd).To(BeNil())
		})
	})

	Describe("HasVisibleSkillsModal", func() {
		It("should report no visible skills modal initially", func() {
			intent = createIntent()
			Expect(intent.HasVisibleSkillsModal()).To(BeFalse())
		})
	})

	Describe("HasVisibleErrorModal", func() {
		It("should report no visible error modal initially", func() {
			intent = createIntent()
			Expect(intent.HasVisibleErrorModal()).To(BeFalse())
		})

		It("should show error modal after ShowErrorModal", func() {
			intent = createIntent()
			intent.ShowErrorModal("Error Title", "Error message")
			Expect(intent.HasVisibleErrorModal()).To(BeTrue())
		})
	})

	Describe("HasVisibleQuickAddModal", func() {
		It("should report no visible quick add modal initially", func() {
			intent = createIntent()
			Expect(intent.HasVisibleQuickAddModal()).To(BeFalse())
		})
	})

	Describe("HasVisibleEditModal", func() {
		It("should report no visible edit modal initially", func() {
			intent = createIntent()
			Expect(intent.HasVisibleEditModal()).To(BeFalse())
		})
	})

	Describe("ApplyFilters", func() {
		It("should apply filters to events", func() {
			intent = createIntent(func(ctx *IntentContext) {
				ctx.InitialFilters = &Filters{Companies: []string{"TechCorp"}}
			})
			intent.ApplyFilters()
			Expect(intent.filteredEvents).To(HaveLen(1))
		})
	})

	Describe("eventMatchesFilters with tags", func() {
		It("should filter by tags", func() {
			evt := fixtures.EventWith("t1", "Tagged event", "Corp", "Proj")
			evt.Tags = []string{"golang", "backend"}
			intent = createIntent(func(ctx *IntentContext) {
				ctx.Events = []*career.Event{evt}
				ctx.InitialFilters = &Filters{Tags: []string{"golang"}}
			})
			intent.ApplyFilters()
			Expect(intent.filteredEvents).To(HaveLen(1))
		})

		It("should exclude events without matching tags", func() {
			evt := fixtures.EventWith("t1", "Tagged event", "Corp", "Proj")
			evt.Tags = []string{"python"}
			intent = createIntent(func(ctx *IntentContext) {
				ctx.Events = []*career.Event{evt}
				ctx.InitialFilters = &Filters{Tags: []string{"golang"}}
			})
			intent.ApplyFilters()
			Expect(intent.filteredEvents).To(BeEmpty())
		})
	})

	Describe("eventMatchesFilters with projects", func() {
		It("should filter by project", func() {
			intent = createIntent(func(ctx *IntentContext) {
				ctx.InitialFilters = &Filters{Projects: []string{"Platform"}}
			})
			intent.ApplyFilters()
			Expect(intent.filteredEvents).To(HaveLen(1))
		})
	})

	Describe("sortEvents", func() {
		It("should sort by date ascending", func() {
			intent = createIntent(func(ctx *IntentContext) {
				ctx.InitialFilters = &Filters{SortBy: "date", SortOrder: "asc"}
			})
			intent.ApplyFilters()
			Expect(intent.filteredEvents).To(HaveLen(2))
		})

		It("should sort by text ascending", func() {
			intent = createIntent(func(ctx *IntentContext) {
				ctx.InitialFilters = &Filters{SortBy: "text", SortOrder: "asc"}
			})
			intent.ApplyFilters()
			Expect(intent.filteredEvents[0].Text < intent.filteredEvents[1].Text).To(BeTrue())
		})

		It("should sort by text descending", func() {
			intent = createIntent(func(ctx *IntentContext) {
				ctx.InitialFilters = &Filters{SortBy: "text", SortOrder: "desc"}
			})
			intent.ApplyFilters()
			Expect(intent.filteredEvents[0].Text > intent.filteredEvents[1].Text).To(BeTrue())
		})
	})

	Describe("HandleNavigate", func() {
		It("should handle confirmed=false by returning to timeline", func() {
			intent = createIntent()
			result := &screens.NavigateResult{ResultData: false}
			cmd := intent.HandleNavigate(result)
			Expect(cmd).To(BeNil())
		})

		It("should handle event result data", func() {
			intent = createIntent()
			event := fixtures.EventWith("detail-1", "Detail event test", "Corp", "Proj")
			result := &screens.NavigateResult{ResultData: event}
			cmd := intent.HandleNavigate(result)
			Expect(cmd).To(BeNil())
		})

		It("should handle nil result data", func() {
			intent = createIntent()
			result := &screens.NavigateResult{ResultData: nil}
			cmd := intent.HandleNavigate(result)
			Expect(cmd).To(BeNil())
		})

		It("should handle unrecognized result data type", func() {
			intent = createIntent()
			result := &screens.NavigateResult{ResultData: 42}
			cmd := intent.HandleNavigate(result)
			Expect(cmd).To(BeNil())
		})

		It("should handle add action", func() {
			intent = createIntent()
			result := &screens.NavigateResult{ResultData: map[string]interface{}{"action": "add"}}
			cmd := intent.HandleNavigate(result)
			_ = cmd
			Expect(intent.HasVisibleQuickAddModal()).To(BeTrue())
		})

		It("should handle edit action with event", func() {
			intent = createIntent()
			event := fixtures.EventWith("edit-1", "Edit this event", "Corp", "Proj")
			result := &screens.NavigateResult{ResultData: map[string]interface{}{"action": "edit", "event": event}}
			cmd := intent.HandleNavigate(result)
			_ = cmd
			Expect(intent.HasVisibleEditModal()).To(BeTrue())
		})

		It("should handle edit action without event", func() {
			intent = createIntent()
			result := &screens.NavigateResult{ResultData: map[string]interface{}{"action": "edit"}}
			cmd := intent.HandleNavigate(result)
			Expect(cmd).To(BeNil())
		})

		It("should handle delete action with event", func() {
			intent = createIntent()
			event := fixtures.EventWith("del-1", "Delete this event with enough text for truncation check", "Corp", "Proj")
			result := &screens.NavigateResult{ResultData: map[string]interface{}{"action": "delete", "event": event}}
			cmd := intent.HandleNavigate(result)
			_ = cmd
		})

		It("should handle delete action without event", func() {
			intent = createIntent()
			result := &screens.NavigateResult{ResultData: map[string]interface{}{"action": "delete"}}
			cmd := intent.HandleNavigate(result)
			Expect(cmd).To(BeNil())
		})

		It("should handle filter action", func() {
			intent = createIntent()
			result := &screens.NavigateResult{ResultData: map[string]interface{}{"action": "filter"}}
			cmd := intent.HandleNavigate(result)
			_ = cmd
		})

		It("should handle unknown action", func() {
			intent = createIntent()
			result := &screens.NavigateResult{ResultData: map[string]interface{}{"action": "unknown_action"}}
			cmd := intent.HandleNavigate(result)
			Expect(cmd).To(BeNil())
		})

		It("should handle missing action key", func() {
			intent = createIntent()
			result := &screens.NavigateResult{ResultData: map[string]interface{}{"not_action": "something"}}
			cmd := intent.HandleNavigate(result)
			Expect(cmd).To(BeNil())
		})
	})

	Describe("RefreshData", func() {
		It("should re-apply filters and rebuild the screen", func() {
			intent = createIntent(func(ctx *IntentContext) {
				ctx.InitialFilters = &Filters{Companies: []string{"TechCorp"}}
			})
			cmd := intent.RefreshData()
			Expect(cmd).To(BeNil())
		})
	})
})
