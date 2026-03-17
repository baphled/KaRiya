package browsetimeline_test

import (
	"context"
	"errors"
	"time"

	cliservice "github.com/baphled/kariya/internal/cli/service"
	"github.com/charmbracelet/bubbles/cursor"
	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"

	"github.com/baphled/kariya/internal/domain/career"
	careerrepo "github.com/baphled/kariya/internal/repository/career"
	careerservice "github.com/baphled/kariya/internal/service/career"
	"github.com/baphled/kariya/internal/service/career/skillinference"
	"github.com/baphled/kariya/internal/testutil/fixtures"
	"github.com/baphled/kariya/internal/tui/intents"
	"github.com/baphled/kariya/internal/tui/intents/browsetimeline"
	"github.com/baphled/kariya/internal/ui/behaviors"
	tea "github.com/charmbracelet/bubbletea"
)

type stubEventService struct {
	captureEventErr   error
	deleteEventErr    error
	listEvents        []*career.Event
	listEventsErr     error
	updateEventErr    error
	getSkillsForEvent map[string][]*career.Skill
	getSkillsErr      error
	listAllSkills     []*career.Skill
	listAllSkillsErr  error
	linkErr           error
	unlinkErr         error
	capturedTexts     []string
	deletedEventIDs   []string
	updatedEvents     []*career.Event
	linkedPairs       []string
	unlinkedPairs     []string
}

func (s *stubEventService) DeleteEvent(_ context.Context, eventID string) error {
	s.deletedEventIDs = append(s.deletedEventIDs, eventID)
	return s.deleteEventErr
}

func (s *stubEventService) ListEvents(_ context.Context, _ *careerrepo.EventListFilters) ([]*career.Event, error) {
	if s.listEvents == nil {
		return []*career.Event{}, s.listEventsErr
	}
	return s.listEvents, s.listEventsErr
}

func (s *stubEventService) CaptureEvent(_ context.Context, text string, _ time.Time, _ careerservice.EventCaptureMode, _ ...cliservice.Option) error {
	s.capturedTexts = append(s.capturedTexts, text)
	return s.captureEventErr
}

func (s *stubEventService) UpdateEventMetadata(_ context.Context, event *career.Event) error {
	s.updatedEvents = append(s.updatedEvents, event)
	return s.updateEventErr
}

func (s *stubEventService) GetSkillsForEvent(_ context.Context, eventID string) ([]*career.Skill, error) {
	if s.getSkillsErr != nil {
		return nil, s.getSkillsErr
	}
	if s.getSkillsForEvent == nil {
		return []*career.Skill{}, nil
	}
	return s.getSkillsForEvent[eventID], nil
}

func (s *stubEventService) LinkSkillToEvent(_ context.Context, eventID string, skillID string) error {
	s.linkedPairs = append(s.linkedPairs, eventID+":"+skillID)
	return s.linkErr
}

func (s *stubEventService) UnlinkSkillFromEvent(_ context.Context, eventID string, skillID string) error {
	s.unlinkedPairs = append(s.unlinkedPairs, eventID+":"+skillID)
	return s.unlinkErr
}

func (s *stubEventService) ListAllSkills(_ context.Context) ([]*career.Skill, error) {
	if s.listAllSkills == nil {
		return []*career.Skill{}, s.listAllSkillsErr
	}
	return s.listAllSkills, s.listAllSkillsErr
}

type stubSkillService struct {
	createErr error
	created   []*career.Skill
}

func (s *stubSkillService) Create(_ context.Context, skill *career.Skill) error {
	s.created = append(s.created, skill)
	return s.createErr
}

type stubSkillInferenceService struct {
	inferResult  *skillinference.InferenceResult
	inferErr     error
	createSkills []*career.Skill
	createErr    error
	createdFrom  []skillinference.SkillSuggestion
}

func (s *stubSkillInferenceService) InferSkillsFromEvents(_ context.Context, _ []*career.Event) (*skillinference.InferenceResult, error) {
	return s.inferResult, s.inferErr
}

func (s *stubSkillInferenceService) CreateSkillsFromSuggestions(_ context.Context, suggestions []skillinference.SkillSuggestion) ([]*career.Skill, error) {
	s.createdFrom = append(s.createdFrom, suggestions...)
	return s.createSkills, s.createErr
}

func newTestIntent(ctx *browsetimeline.IntentValidator) *browsetimeline.Intent {
	intent, err := browsetimeline.NewIntent(ctx)
	Expect(err).NotTo(HaveOccurred())
	Expect(intent).NotTo(BeNil())
	Expect(intent.Init()).To(BeNil())
	return intent
}

func updateIntentForTest(intent *browsetimeline.Intent, msg tea.Msg) {
	cmd := intent.Update(msg)
	flushIntentCmdForTest(intent, cmd)
}

func flushIntentCmdForTest(intent *browsetimeline.Intent, cmd tea.Cmd) {
	for step := 0; cmd != nil && step < 40; step++ {
		msg := cmd()
		if msg == nil {
			return
		}
		if _, ok := msg.(cursor.BlinkMsg); ok {
			return
		}
		if batch, ok := msg.(tea.BatchMsg); ok {
			for _, nested := range batch {
				flushIntentCmdForTest(intent, nested)
			}
			return
		}
		cmd = intent.Update(msg)
	}
}

func advanceSkillEditorToConfirm(intent *browsetimeline.Intent) {
	for range 4 {
		updateIntentForTest(intent, tea.KeyMsg{Type: tea.KeyTab})
	}
}

var _ = Describe("Intent", func() {
	var (
		intent *browsetimeline.Intent
		ctx    *browsetimeline.IntentValidator
		events []*career.Event
	)

	BeforeEach(func() {
		events = []*career.Event{
			fixtures.EventWith("event-1", "Backend Developer at TechCorp", "TechCorp", "Platform"),
			fixtures.EventWith("event-2", "DevOps Engineer at CloudInc", "CloudInc", "Infrastructure"),
		}
		ctx = &browsetimeline.IntentValidator{
			Events: events,
		}
	})

	Describe("Construction", func() {
		Context("with valid context", func() {
			It("should create an intent successfully", func() {
				intent, err := browsetimeline.NewIntent(ctx)
				Expect(err).NotTo(HaveOccurred())
				Expect(intent).NotTo(BeNil())
			})

			It("should embed BaseIntent", func() {
				intent, _ := browsetimeline.NewIntent(ctx)
				Expect(intent.BaseIntent).NotTo(BeNil())
			})
		})

		Context("with empty events", func() {
			It("should create an intent with empty event list", func() {
				emptyCtx := &browsetimeline.IntentValidator{
					Events: []*career.Event{},
				}
				intent, err := browsetimeline.NewIntent(emptyCtx)
				Expect(err).NotTo(HaveOccurred())
				Expect(intent).NotTo(BeNil())
			})
		})

		Context("with nil events", func() {
			It("should handle nil events by initializing empty slice", func() {
				nilCtx := &browsetimeline.IntentValidator{
					Events: nil,
				}
				intent, err := browsetimeline.NewIntent(nilCtx)
				Expect(err).NotTo(HaveOccurred())
				Expect(intent).NotTo(BeNil())
			})
		})
	})

	Describe("Initialization", func() {
		BeforeEach(func() {
			var err error
			intent, err = browsetimeline.NewIntent(ctx)
			Expect(err).NotTo(HaveOccurred())
		})

		It("should return nil command on init", func() {
			cmd := intent.Init()
			Expect(cmd).To(BeNil())
		})

		It("should render view after init", func() {
			intent.Init()
			view := intent.View()
			Expect(view).NotTo(BeEmpty())
		})

		It("should show timeline content after init", func() {
			intent.Init()
			view := intent.View()
			Expect(view).To(ContainSubstring("Timeline"))
		})

		It("should display events in view", func() {
			intent.Init()
			view := intent.View()
			Expect(view).To(ContainSubstring("TechCorp"))
			Expect(view).To(ContainSubstring("CloudInc"))
		})
	})

	Describe("View Rendering", func() {
		BeforeEach(func() {
			var err error
			intent, err = browsetimeline.NewIntent(ctx)
			Expect(err).NotTo(HaveOccurred())
			intent.Init()
		})

		Context("with events", func() {
			It("should show event count", func() {
				view := intent.View()
				Expect(view).To(ContainSubstring("Events: 2"))
			})

			It("should show navigation hints", func() {
				view := intent.View()
				Expect(view).To(SatisfyAny(
					ContainSubstring("j/k"),
					ContainSubstring("Enter"),
				))
			})
		})

		Context("with empty events", func() {
			BeforeEach(func() {
				emptyCtx := &browsetimeline.IntentValidator{
					Events: []*career.Event{},
				}
				var err error
				intent, err = browsetimeline.NewIntent(emptyCtx)
				Expect(err).NotTo(HaveOccurred())
				intent.Init()
			})

			It("should show empty state message", func() {
				view := intent.View()
				Expect(view).To(ContainSubstring("No events"))
			})

			It("should show add hint", func() {
				view := intent.View()
				Expect(view).To(ContainSubstring("Add"))
			})
		})
	})

	Describe("Navigation", func() {
		BeforeEach(func() {
			var err error
			intent, err = browsetimeline.NewIntent(ctx)
			Expect(err).NotTo(HaveOccurred())
			intent.Init()
		})

		Context("keyboard navigation", func() {
			It("should handle arrow down", func() {
				cmd := intent.Update(tea.KeyMsg{Type: tea.KeyDown})
				Expect(cmd).To(BeNil())
				view := intent.View()
				Expect(view).NotTo(BeEmpty())
			})

			It("should handle arrow up", func() {
				intent.Update(tea.KeyMsg{Type: tea.KeyDown})
				cmd := intent.Update(tea.KeyMsg{Type: tea.KeyUp})
				Expect(cmd).To(BeNil())
			})

			It("should handle vim j key", func() {
				cmd := intent.Update(tea.KeyMsg{Type: tea.KeyRunes, Runes: []rune{'j'}})
				Expect(cmd).To(BeNil())
			})

			It("should handle vim k key", func() {
				intent.Update(tea.KeyMsg{Type: tea.KeyRunes, Runes: []rune{'j'}})
				cmd := intent.Update(tea.KeyMsg{Type: tea.KeyRunes, Runes: []rune{'k'}})
				Expect(cmd).To(BeNil())
			})
		})

		Context("boundary behavior", func() {
			It("should not crash when navigating up at first item", func() {
				intent.Update(tea.KeyMsg{Type: tea.KeyUp})
				view := intent.View()
				Expect(view).NotTo(BeEmpty())
			})

			It("should not crash when navigating down past last item", func() {
				for range 10 {
					intent.Update(tea.KeyMsg{Type: tea.KeyDown})
				}
				view := intent.View()
				Expect(view).NotTo(BeEmpty())
			})
		})
	})

	Describe("Cancellation", func() {
		BeforeEach(func() {
			var err error
			intent, err = browsetimeline.NewIntent(ctx)
			Expect(err).NotTo(HaveOccurred())
			intent.Init()
		})

		Context("from list view", func() {
			It("should cancel intent on escape", func() {
				intent.Update(tea.KeyMsg{Type: tea.KeyEsc})
				result := intent.Result()
				Expect(result).NotTo(BeNil())
				Expect(result.Status).To(Equal(intents.Cancelled))
			})
		})

		Context("from detail modal", func() {
			It("should close modal on escape without cancelling intent", func() {
				intent.Update(tea.KeyMsg{Type: tea.KeyEnter})
				intent.Update(tea.KeyMsg{Type: tea.KeyEsc})
				result := intent.Result()
				Expect(result).To(BeNil())
			})
		})
	})

	Describe("Event Selection", func() {
		BeforeEach(func() {
			var err error
			intent, err = browsetimeline.NewIntent(ctx)
			Expect(err).NotTo(HaveOccurred())
			intent.Init()
		})

		It("should show detail modal on enter", func() {
			intent.Update(tea.KeyMsg{Type: tea.KeyEnter})
			view := intent.View()
			// Should show one of the events in detail view.
			Expect(view).To(SatisfyAny(
				ContainSubstring("Backend Developer"),
				ContainSubstring("DevOps Engineer"),
			))
		})

		It("should return to list after closing modal", func() {
			intent.Update(tea.KeyMsg{Type: tea.KeyEnter})
			intent.Update(tea.KeyMsg{Type: tea.KeyEsc})
			view := intent.View()
			Expect(view).To(ContainSubstring("Timeline"))
		})
	})

	Describe("Modal Interactions", func() {
		BeforeEach(func() {
			var err error
			intent, err = browsetimeline.NewIntent(ctx)
			Expect(err).NotTo(HaveOccurred())
			intent.Init()
		})

		Context("search modal", func() {
			It("should open search modal with / key", func() {
				intent.Update(tea.KeyMsg{Type: tea.KeyRunes, Runes: []rune{'/'}})
				view := intent.View()
				Expect(view).To(ContainSubstring("Search"))
			})
		})

		Context("filter modal", func() {
			It("should open filter modal with f key", func() {
				intent.Update(tea.KeyMsg{Type: tea.KeyRunes, Runes: []rune{'f'}})
				view := intent.View()
				Expect(view).To(ContainSubstring("Filter"))
			})
		})

		Context("sort modal", func() {
			It("should open sort modal with s key", func() {
				intent.Update(tea.KeyMsg{Type: tea.KeyRunes, Runes: []rune{'s'}})
				view := intent.View()
				Expect(view).To(ContainSubstring("Sort"))
			})
		})

		Context("quick add modal", func() {
			It("should open quick add modal with a key", func() {
				intent.Update(tea.KeyMsg{Type: tea.KeyRunes, Runes: []rune{'a'}})
				Expect(intent.HasVisibleQuickAddModal()).To(BeTrue())
			})
		})

		Context("edit modal", func() {
			It("should open edit modal with e key", func() {
				intent.Update(tea.KeyMsg{Type: tea.KeyRunes, Runes: []rune{'e'}})
				Expect(intent.HasVisibleEditModal()).To(BeTrue())
			})
		})

		Context("delete modal", func() {
			It("should open delete confirmation with d key", func() {
				intent.Update(tea.KeyMsg{Type: tea.KeyRunes, Runes: []rune{'d'}})
				view := intent.View()
				Expect(view).To(ContainSubstring("Delete"))
			})
		})

		Context("modal priority", func() {
			It("should not allow multiple overlays", func() {
				intent.Update(tea.KeyMsg{Type: tea.KeyRunes, Runes: []rune{'/'}})
				intent.Update(tea.KeyMsg{Type: tea.KeyRunes, Runes: []rune{'f'}})
				view := intent.View()
				Expect(view).To(ContainSubstring("Search"))
				Expect(view).NotTo(ContainSubstring("Filter by Company"))
			})
		})
	})

	Describe("Window Resize", func() {
		BeforeEach(func() {
			var err error
			intent, err = browsetimeline.NewIntent(ctx)
			Expect(err).NotTo(HaveOccurred())
			intent.Init()
		})

		It("should handle window size message", func() {
			cmd := intent.Update(tea.WindowSizeMsg{Width: 80, Height: 24})
			Expect(cmd).To(BeNil())
			view := intent.View()
			Expect(view).NotTo(BeEmpty())
		})

		It("should adapt to different terminal sizes", func() {
			intent.Update(tea.WindowSizeMsg{Width: 120, Height: 40})
			view := intent.View()
			Expect(view).NotTo(BeEmpty())
		})
	})

	Describe("Error Handling", func() {
		BeforeEach(func() {
			var err error
			intent, err = browsetimeline.NewIntent(ctx)
			Expect(err).NotTo(HaveOccurred())
			intent.Init()
		})

		It("should show error modal when error occurs", func() {
			intent.ShowErrorModal("Test Error", "Something went wrong")
			Expect(intent.HasVisibleErrorModal()).To(BeTrue())
		})

		It("should dismiss error modal with escape", func() {
			intent.ShowErrorModal("Test Error", "Something went wrong")
			intent.Update(tea.KeyMsg{Type: tea.KeyEsc})
			Expect(intent.HasVisibleErrorModal()).To(BeFalse())
		})

		It("should render error content in view", func() {
			intent.ShowErrorModal("Operation Failed", "Database error")
			view := intent.View()
			Expect(view).To(ContainSubstring("Operation Failed"))
		})
	})

	Describe("Filter Behavior", func() {
		BeforeEach(func() {
			var err error
			intent, err = browsetimeline.NewIntent(ctx)
			Expect(err).NotTo(HaveOccurred())
			intent.Init()
		})

		It("should implement FilterBehavior interface", func() {
			var _ behaviors.FilterBehavior = intent
		})

		It("should report no active filters initially", func() {
			Expect(intent.HasActiveFilters()).To(BeFalse())
		})

		It("should clear filters", func() {
			intent.ClearFilters()
			Expect(intent.HasActiveFilters()).To(BeFalse())
		})
	})

	Describe("Inactive State", func() {
		BeforeEach(func() {
			var err error
			intent, err = browsetimeline.NewIntent(ctx)
			Expect(err).NotTo(HaveOccurred())
			intent.Init()
			intent.Update(tea.KeyMsg{Type: tea.KeyEsc})
		})

		It("should not process updates when inactive", func() {
			cmd := intent.Update(tea.KeyMsg{Type: tea.KeyEnter})
			Expect(cmd).To(BeNil())
		})
	})

	Describe("Message-driven skill workflows", func() {
		var (
			eventService     *stubEventService
			skillService     *stubSkillService
			inferenceService *stubSkillInferenceService
		)

		BeforeEach(func() {
			eventService = &stubEventService{}
			skillService = &stubSkillService{}
			inferenceService = &stubSkillInferenceService{}
			intent = newTestIntent(&browsetimeline.IntentValidator{
				Events:                events,
				CLIEventService:       eventService,
				CLISkillCreator:       skillService,
				SkillInferenceService: inferenceService,
			})
		})

		openFirstEvent := func() {
			Expect(intent.Update(tea.KeyMsg{Type: tea.KeyEnter})).To(BeNil())
		}

		showSkillsModal := func(skills []*career.Skill) {
			openFirstEvent()
			Expect(intent.Update(browsetimeline.SkillsForModalLoadedMsg{
				EventID: events[0].ID,
				Skills:  skills,
			})).To(BeNil())
			Expect(intent.HasVisibleSkillsModal()).To(BeTrue())
		}

		It("shows an error modal when linking a skill fails", func() {
			cmd := intent.Update(browsetimeline.SkillLinkedMsg{Error: errors.New("link failed")})
			Expect(cmd).To(BeNil())
			Expect(intent.HasVisibleErrorModal()).To(BeTrue())
			Expect(intent.View()).To(ContainSubstring("Error Linking Skill"))
		})

		It("returns a refresh command when a linked skill succeeds with the skills modal open", func() {
			showSkillsModal([]*career.Skill{fixtures.SkillWith("skill-1", "Go", "", "")})
			cmd := intent.Update(browsetimeline.SkillLinkedMsg{EventID: events[0].ID, SkillID: "skill-1"})
			Expect(cmd).NotTo(BeNil())
			Expect(intent.HasVisibleSkillsModal()).To(BeTrue())
		})

		It("returns a refresh command when unlinking succeeds with the skills modal open", func() {
			showSkillsModal([]*career.Skill{fixtures.SkillWith("skill-1", "Go", "", "")})
			cmd := intent.Update(browsetimeline.SkillUnlinkedMsg{EventID: events[0].ID, SkillID: "skill-1"})
			Expect(cmd).NotTo(BeNil())
		})

		It("returns a refresh command when creating a skill succeeds with the skills modal open", func() {
			showSkillsModal([]*career.Skill{fixtures.SkillWith("skill-1", "Go", "", "")})
			cmd := intent.Update(browsetimeline.SkillCreatedMsg{Skill: fixtures.SkillWith("skill-2", "Docker", "", "")})
			Expect(cmd).NotTo(BeNil())
		})

		It("shows the skills modal when skills are loaded for the current event", func() {
			openFirstEvent()
			cmd := intent.Update(browsetimeline.SkillsForModalLoadedMsg{
				EventID: events[0].ID,
				Skills:  []*career.Skill{fixtures.SkillWith("skill-1", "Go", "", "")},
			})
			Expect(cmd).To(BeNil())
			Expect(intent.HasVisibleSkillsModal()).To(BeTrue())
			Expect(intent.View()).To(ContainSubstring("Go"))
		})

		It("shows an error when loading event skills fails", func() {
			openFirstEvent()
			cmd := intent.Update(browsetimeline.SkillsForModalLoadedMsg{
				EventID: events[0].ID,
				Error:   errors.New("load failed"),
			})
			Expect(cmd).To(BeNil())
			Expect(intent.HasVisibleErrorModal()).To(BeTrue())
			Expect(intent.View()).To(ContainSubstring("Error Loading Skills"))
		})

		It("shows the skill picker with only unlinked skills", func() {
			openFirstEvent()
			cmd := intent.Update(browsetimeline.SkillPickerDataLoadedMsg{
				AllSkills: []*career.Skill{
					fixtures.SkillWith("skill-1", "Go", "", ""),
					fixtures.SkillWith("skill-2", "Docker", "", ""),
				},
				EventSkills: []*career.Skill{fixtures.SkillWith("skill-1", "Go", "", "")},
			})
			Expect(cmd).To(BeNil())
			view := intent.View()
			Expect(view).To(ContainSubstring("Docker"))
			Expect(view).NotTo(ContainSubstring("Go"))
		})

		It("shows an error when picker data loading fails", func() {
			openFirstEvent()
			cmd := intent.Update(browsetimeline.SkillPickerDataLoadedMsg{Error: errors.New("picker failed")})
			Expect(cmd).To(BeNil())
			Expect(intent.HasVisibleErrorModal()).To(BeTrue())
			Expect(intent.View()).To(ContainSubstring("Error Loading Skills"))
		})

		It("updates the visible skills modal when refreshed skills arrive", func() {
			showSkillsModal([]*career.Skill{fixtures.SkillWith("skill-1", "Go", "", "")})
			cmd := intent.Update(browsetimeline.SkillsRefreshedMsg{Skills: []*career.Skill{fixtures.SkillWith("skill-2", "Docker", "", "")}})
			Expect(cmd).To(BeNil())
			Expect(intent.View()).To(ContainSubstring("Docker"))
		})

		It("shows an error when refreshing skills fails", func() {
			showSkillsModal([]*career.Skill{fixtures.SkillWith("skill-1", "Go", "", "")})
			cmd := intent.Update(browsetimeline.SkillsRefreshedMsg{Error: errors.New("refresh failed")})
			Expect(cmd).To(BeNil())
			Expect(intent.HasVisibleErrorModal()).To(BeTrue())
			Expect(intent.View()).To(ContainSubstring("Error Loading Skills"))
		})

		It("shows an inference error modal when suggestion loading fails", func() {
			cmd := intent.Update(browsetimeline.SkillSuggestionsErrorMsg{Error: errors.New("boom")})
			Expect(cmd).To(BeNil())
			Expect(intent.HasVisibleErrorModal()).To(BeTrue())
			Expect(intent.View()).To(ContainSubstring("Skill Inference Failed"))
		})

		It("shows a no-skills error when inference returns no suggestions", func() {
			cmd := intent.Update(browsetimeline.SkillSuggestionsLoadedMsg{})
			Expect(cmd).To(BeNil())
			Expect(intent.HasVisibleErrorModal()).To(BeTrue())
			Expect(intent.View()).To(ContainSubstring("No Skills Detected"))
		})

		It("shows an all-tracked error when every suggestion already exists", func() {
			cmd := intent.Update(browsetimeline.SkillSuggestionsLoadedMsg{
				Suggestions:        []skillinference.SkillSuggestion{{Name: "Go"}},
				ExistingSkillNames: []string{"Go"},
			})
			Expect(cmd).To(BeNil())
			Expect(intent.HasVisibleErrorModal()).To(BeTrue())
			Expect(intent.View()).To(ContainSubstring("All Skills Tracked"))
		})

		It("renders the suggestion review modal for newly inferred skills", func() {
			cmd := intent.Update(browsetimeline.SkillSuggestionsLoadedMsg{
				Suggestions: []skillinference.SkillSuggestion{{Name: "Docker", Category: "backend"}},
			})
			Expect(cmd).To(BeNil())
			Expect(intent.View()).To(ContainSubstring("Docker"))
		})

		It("creates a skill from an accepted suggestion and shows a success modal on close", func() {
			createdSkill := fixtures.SkillWith("skill-99", "Docker", "", "")
			inferenceService.createSkills = []*career.Skill{createdSkill}
			showSkillsModal([]*career.Skill{fixtures.SkillWith("skill-1", "Go", "", "")})
			Expect(intent.Update(browsetimeline.SkillSuggestionsLoadedMsg{
				Suggestions: []skillinference.SkillSuggestion{{Name: "Docker", Category: "backend"}},
			})).To(BeNil())

			acceptCmd := intent.Update(tea.KeyMsg{Type: tea.KeyRunes, Runes: []rune{'a'}})
			Expect(acceptCmd).NotTo(BeNil())
			Expect(intent.HasVisibleErrorModal()).To(BeTrue())
			Expect(inferenceService.createdFrom).To(HaveLen(1))
			Expect(eventService.linkedPairs).To(HaveLen(1))
			Expect(eventService.linkedPairs[0]).To(HaveSuffix(":" + createdSkill.ID))
			Expect(intent.View()).To(ContainSubstring("Skills Created"))
		})
	})

	Describe("Update-driven modal flows", func() {
		It("closes the search modal on escape", func() {
			intent = newTestIntent(ctx)
			Expect(intent.Update(tea.KeyMsg{Type: tea.KeyRunes, Runes: []rune{'/'}})).NotTo(BeNil())
			Expect(intent.View()).To(ContainSubstring("Search"))
			Expect(intent.Update(tea.KeyMsg{Type: tea.KeyEsc})).NotTo(BeNil())
			Expect(intent.View()).To(ContainSubstring("Timeline"))
		})

		It("closes the filter modal on escape", func() {
			intent = newTestIntent(ctx)
			Expect(intent.Update(tea.KeyMsg{Type: tea.KeyRunes, Runes: []rune{'f'}})).NotTo(BeNil())
			Expect(intent.View()).To(ContainSubstring("Filter"))
			Expect(intent.Update(tea.KeyMsg{Type: tea.KeyEsc})).NotTo(BeNil())
			Expect(intent.View()).To(ContainSubstring("Timeline"))
		})

		It("closes the sort modal on escape", func() {
			intent = newTestIntent(ctx)
			Expect(intent.Update(tea.KeyMsg{Type: tea.KeyRunes, Runes: []rune{'s'}})).NotTo(BeNil())
			Expect(intent.View()).To(ContainSubstring("Sort"))
			Expect(intent.Update(tea.KeyMsg{Type: tea.KeyEsc})).NotTo(BeNil())
			Expect(intent.View()).To(ContainSubstring("Timeline"))
		})

		It("opens the skills loader command from event detail", func() {
			intent = newTestIntent(ctx)
			Expect(intent.Update(tea.KeyMsg{Type: tea.KeyEnter})).To(BeNil())
			cmd := intent.Update(tea.KeyMsg{Type: tea.KeyRunes, Runes: []rune{'s'}})
			Expect(cmd).NotTo(BeNil())
		})

		It("opens the edit modal from event detail", func() {
			intent = newTestIntent(ctx)
			Expect(intent.Update(tea.KeyMsg{Type: tea.KeyEnter})).To(BeNil())
			cmd := intent.Update(tea.KeyMsg{Type: tea.KeyRunes, Runes: []rune{'e'}})
			Expect(cmd).NotTo(BeNil())
			Expect(intent.HasVisibleEditModal()).To(BeTrue())
		})

		It("opens the delete confirmation from event detail", func() {
			intent = newTestIntent(ctx)
			Expect(intent.Update(tea.KeyMsg{Type: tea.KeyEnter})).To(BeNil())
			cmd := intent.Update(tea.KeyMsg{Type: tea.KeyRunes, Runes: []rune{'d'}})
			Expect(cmd).NotTo(BeNil())
			Expect(intent.View()).To(ContainSubstring("Delete Event"))
		})

		It("returns a picker command when adding an existing skill from the skills modal", func() {
			intent = newTestIntent(ctx)
			Expect(intent.Update(tea.KeyMsg{Type: tea.KeyEnter})).To(BeNil())
			Expect(intent.Update(browsetimeline.SkillsForModalLoadedMsg{
				EventID: events[0].ID,
				Skills:  []*career.Skill{fixtures.SkillWith("skill-1", "Go", "", "")},
			})).To(BeNil())
			cmd := intent.Update(tea.KeyMsg{Type: tea.KeyRunes, Runes: []rune{'a'}})
			Expect(cmd).NotTo(BeNil())
		})

		It("returns a skill add command when adding a new skill from the skills modal", func() {
			intent = newTestIntent(ctx)
			Expect(intent.Update(tea.KeyMsg{Type: tea.KeyEnter})).To(BeNil())
			Expect(intent.Update(browsetimeline.SkillsForModalLoadedMsg{
				EventID: events[0].ID,
				Skills:  []*career.Skill{fixtures.SkillWith("skill-1", "Go", "", "")},
			})).To(BeNil())
			cmd := intent.Update(tea.KeyMsg{Type: tea.KeyRunes, Runes: []rune{'n'}})
			Expect(cmd).NotTo(BeNil())
		})

		It("returns an inference command when inferring skills from the skills modal", func() {
			intent = newTestIntent(&browsetimeline.IntentValidator{
				Events:                events,
				SkillInferenceService: &stubSkillInferenceService{},
			})
			Expect(intent.Update(tea.KeyMsg{Type: tea.KeyEnter})).To(BeNil())
			Expect(intent.Update(browsetimeline.SkillsForModalLoadedMsg{
				EventID: events[0].ID,
				Skills:  []*career.Skill{fixtures.SkillWith("skill-1", "Go", "", "")},
			})).To(BeNil())
			cmd := intent.Update(tea.KeyMsg{Type: tea.KeyRunes, Runes: []rune{'i'}})
			Expect(cmd).NotTo(BeNil())
		})

		It("returns an unlink command when removing a skill from the skills modal", func() {
			eventService := &stubEventService{}
			intent = newTestIntent(&browsetimeline.IntentValidator{
				Events:          events,
				CLIEventService: eventService,
			})
			Expect(intent.Update(tea.KeyMsg{Type: tea.KeyEnter})).To(BeNil())
			Expect(intent.Update(browsetimeline.SkillsForModalLoadedMsg{
				EventID: events[0].ID,
				Skills:  []*career.Skill{fixtures.SkillWith("skill-1", "Go", "", "")},
			})).To(BeNil())
			cmd := intent.Update(tea.KeyMsg{Type: tea.KeyRunes, Runes: []rune{'d'}})
			Expect(cmd).NotTo(BeNil())
		})

		It("returns a link command when selecting a skill from the picker modal", func() {
			intent = newTestIntent(&browsetimeline.IntentValidator{
				Events:          events,
				CLIEventService: &stubEventService{},
			})
			Expect(intent.Update(tea.KeyMsg{Type: tea.KeyEnter})).To(BeNil())
			Expect(intent.Update(browsetimeline.SkillPickerDataLoadedMsg{
				AllSkills:   []*career.Skill{fixtures.SkillWith("skill-1", "Go", "", "")},
				EventSkills: []*career.Skill{},
			})).To(BeNil())
			cmd := intent.Update(tea.KeyMsg{Type: tea.KeyEnter})
			Expect(cmd).NotTo(BeNil())
		})

		It("shows an error when event deletion fails", func() {
			eventService := &stubEventService{deleteEventErr: errors.New("cannot delete")}
			intent = newTestIntent(&browsetimeline.IntentValidator{
				Events:          events,
				CLIEventService: eventService,
			})
			intent.Update(tea.KeyMsg{Type: tea.KeyRunes, Runes: []rune{'d'}})
			Expect(intent.Update(tea.KeyMsg{Type: tea.KeyEnter})).NotTo(BeNil())
			Expect(intent.HasVisibleErrorModal()).To(BeTrue())
			Expect(intent.View()).To(ContainSubstring("Delete Failed"))
		})

		It("submits the edit modal through Update and persists the event", func() {
			eventService := &stubEventService{}
			intent = newTestIntent(&browsetimeline.IntentValidator{
				Events:          events,
				CLIEventService: eventService,
			})
			Expect(intent.Update(tea.KeyMsg{Type: tea.KeyRunes, Runes: []rune{'e'}})).NotTo(BeNil())
			Expect(intent.HasVisibleEditModal()).To(BeTrue())
			Expect(intent.Update(tea.KeyMsg{Type: tea.KeyCtrlS})).NotTo(BeNil())
			Expect(eventService.updatedEvents).To(HaveLen(1))
			Expect(eventService.updatedEvents[0].ID).To(BeElementOf(events[0].ID, events[1].ID))
		})

		It("shows an error when persisting an edited event fails", func() {
			eventService := &stubEventService{updateEventErr: errors.New("cannot update")}
			intent = newTestIntent(&browsetimeline.IntentValidator{
				Events:          events,
				CLIEventService: eventService,
			})
			Expect(intent.Update(tea.KeyMsg{Type: tea.KeyRunes, Runes: []rune{'e'}})).NotTo(BeNil())
			Expect(intent.Update(tea.KeyMsg{Type: tea.KeyCtrlS})).NotTo(BeNil())
			Expect(intent.HasVisibleErrorModal()).To(BeTrue())
			Expect(intent.View()).To(ContainSubstring("Edit Failed"))
		})

		It("shows a refresh error when quick add submission cannot reload events", func() {
			eventService := &stubEventService{listEventsErr: errors.New("reload failed")}
			intent = newTestIntent(&browsetimeline.IntentValidator{
				Events:          events,
				CLIEventService: eventService,
			})
			Expect(intent.Update(tea.KeyMsg{Type: tea.KeyRunes, Runes: []rune{'a'}})).NotTo(BeNil())
			Expect(intent.HasVisibleQuickAddModal()).To(BeTrue())
			Expect(intent.Update(tea.KeyMsg{Type: tea.KeyCtrlS})).NotTo(BeNil())
			Expect(intent.HasVisibleErrorModal()).To(BeTrue())
			Expect(intent.View()).To(ContainSubstring("Refresh Failed"))
		})

		It("clears active filters from the x shortcut", func() {
			intent = newTestIntent(&browsetimeline.IntentValidator{
				Events: events,
				InitialFilters: &browsetimeline.Filters{
					Companies: []string{"TechCorp"},
				},
			})
			Expect(intent.HasActiveFilters()).To(BeTrue())
			Expect(intent.View()).NotTo(ContainSubstring("CloudInc"))
			Expect(intent.Update(tea.KeyMsg{Type: tea.KeyRunes, Runes: []rune{'x'}})).To(BeNil())
			Expect(intent.HasActiveFilters()).To(BeFalse())
			Expect(intent.View()).To(ContainSubstring("CloudInc"))
		})

		It("executes skill loading and picker commands through public modal flows", func() {
			eventService := &stubEventService{
				listAllSkills: []*career.Skill{fixtures.SkillWith("skill-1", "Go", "", ""), fixtures.SkillWith("skill-2", "Docker", "", "")},
				getSkillsForEvent: map[string][]*career.Skill{
					events[1].ID: {fixtures.SkillWith("skill-1", "Go", "", "")},
				},
			}
			intent = newTestIntent(&browsetimeline.IntentValidator{
				Events:          events,
				CLIEventService: eventService,
			})

			updateIntentForTest(intent, tea.KeyMsg{Type: tea.KeyEnter})
			cmd := intent.Update(tea.KeyMsg{Type: tea.KeyRunes, Runes: []rune{'s'}})
			Expect(cmd).NotTo(BeNil())
			flushIntentCmdForTest(intent, cmd)
			Expect(intent.HasVisibleSkillsModal()).To(BeTrue())

			cmd = intent.Update(tea.KeyMsg{Type: tea.KeyRunes, Runes: []rune{'a'}})
			Expect(cmd).NotTo(BeNil())
			flushIntentCmdForTest(intent, cmd)
			Expect(intent.View()).To(ContainSubstring("Select Skill"))
			Expect(intent.View()).To(ContainSubstring("Docker"))
		})

		It("links the selected picker skill when its command is executed", func() {
			eventService := &stubEventService{}
			intent = newTestIntent(&browsetimeline.IntentValidator{
				Events:          events,
				CLIEventService: eventService,
			})

			updateIntentForTest(intent, tea.KeyMsg{Type: tea.KeyEnter})
			updateIntentForTest(intent, browsetimeline.SkillPickerDataLoadedMsg{
				AllSkills:   []*career.Skill{fixtures.SkillWith("skill-1", "Go", "", "")},
				EventSkills: []*career.Skill{},
			})

			cmd := intent.Update(tea.KeyMsg{Type: tea.KeyEnter})
			Expect(cmd).NotTo(BeNil())
			flushIntentCmdForTest(intent, cmd)
			Expect(eventService.linkedPairs).To(HaveLen(1))
			Expect(eventService.linkedPairs[0]).To(HaveSuffix(":skill-1"))
		})

		It("submits search filter through the modal flow", func() {
			intent = newTestIntent(ctx)

			updateIntentForTest(intent, tea.KeyMsg{Type: tea.KeyRunes, Runes: []rune{'/'}})
			Expect(intent.View()).To(ContainSubstring("Search"))
			updateIntentForTest(intent, tea.KeyMsg{Type: tea.KeyRunes, Runes: []rune("CloudInc")})
			updateIntentForTest(intent, tea.KeyMsg{Type: tea.KeyTab})
			updateIntentForTest(intent, tea.KeyMsg{Type: tea.KeyRight})
			updateIntentForTest(intent, tea.KeyMsg{Type: tea.KeyEnter})

			Expect(intent.HasActiveFilters()).To(BeTrue())
			Expect(intent.View()).To(ContainSubstring("CloudInc"))
			Expect(intent.View()).NotTo(ContainSubstring("TechCorp"))
		})

		It("submits filter and sort modal flows through public updates", func() {
			intent = newTestIntent(ctx)

			updateIntentForTest(intent, tea.KeyMsg{Type: tea.KeyRunes, Runes: []rune{'f'}})
			Expect(intent.View()).To(ContainSubstring("Filter"))
			updateIntentForTest(intent, tea.KeyMsg{Type: tea.KeyEnter})
			Expect(intent.HasActiveFilters()).To(BeFalse())

			updateIntentForTest(intent, tea.KeyMsg{Type: tea.KeyRunes, Runes: []rune{'s'}})
			Expect(intent.View()).To(ContainSubstring("Sort"))
			updateIntentForTest(intent, tea.KeyMsg{Type: tea.KeyEnter})
			if !intent.HasActiveFilters() {
				updateIntentForTest(intent, tea.KeyMsg{Type: tea.KeyEnter})
			}
			Expect(intent.View()).To(ContainSubstring("Timeline"))
		})

		It("creates and links a new skill through the public skill add flow", func() {
			eventService := &stubEventService{}
			skillService := &stubSkillService{}
			intent = newTestIntent(&browsetimeline.IntentValidator{
				Events:          events,
				CLIEventService: eventService,
				CLISkillCreator: skillService,
			})

			updateIntentForTest(intent, tea.KeyMsg{Type: tea.KeyEnter})
			updateIntentForTest(intent, browsetimeline.SkillsForModalLoadedMsg{
				EventID: events[1].ID,
				Skills:  []*career.Skill{fixtures.SkillWith("skill-1", "Go", "", "")},
			})
			updateIntentForTest(intent, tea.KeyMsg{Type: tea.KeyRunes, Runes: []rune{'n'}})
			updateIntentForTest(intent, tea.KeyMsg{Type: tea.KeyRunes, Runes: []rune("GraphQL")})
			advanceSkillEditorToConfirm(intent)
			updateIntentForTest(intent, tea.KeyMsg{Type: tea.KeyRight})
			updateIntentForTest(intent, tea.KeyMsg{Type: tea.KeyEnter})

			Expect(skillService.created).To(HaveLen(1))
			Expect(skillService.created[0].Name).To(Equal("GraphQL"))
			Expect(eventService.linkedPairs).To(HaveLen(1))
			Expect(eventService.linkedPairs[0]).To(HaveSuffix(":" + skillService.created[0].ID))
		})

		It("cancels and confirms delete through public updates", func() {
			cancelService := &stubEventService{}
			intent = newTestIntent(&browsetimeline.IntentValidator{
				Events:          events,
				CLIEventService: cancelService,
			})

			intent.Update(tea.KeyMsg{Type: tea.KeyRunes, Runes: []rune{'d'}})
			Expect(intent.View()).To(ContainSubstring("Delete Event"))
			Expect(intent.Update(tea.KeyMsg{Type: tea.KeyRunes, Runes: []rune{'n'}})).NotTo(BeNil())
			Expect(cancelService.deletedEventIDs).To(BeEmpty())

			confirmService := &stubEventService{}
			intent = newTestIntent(&browsetimeline.IntentValidator{
				Events:          events,
				CLIEventService: confirmService,
			})

			intent.Update(tea.KeyMsg{Type: tea.KeyRunes, Runes: []rune{'d'}})
			Expect(intent.Update(tea.KeyMsg{Type: tea.KeyRunes, Runes: []rune{'y'}})).NotTo(BeNil())
			Expect(confirmService.deletedEventIDs).To(HaveLen(1))
		})

	})
})
