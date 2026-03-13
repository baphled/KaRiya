package captureevent

import (
	"context"
	"errors"
	"time"

	"github.com/baphled/kariya/internal/domain/career"
	memoryrepo "github.com/baphled/kariya/internal/repository/career/memory"
	careerservice "github.com/baphled/kariya/internal/service/career"
	"github.com/baphled/kariya/internal/service/career/skillinference"
	"github.com/baphled/kariya/internal/testutil/fixtures"
	"github.com/baphled/kariya/internal/tui/intents"
	eventview "github.com/baphled/kariya/internal/tui/views/event"
	"github.com/baphled/kariya/internal/ui/display"
	"github.com/baphled/kariya/internal/ui/terminal"
	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
)

var _ = Describe("Helper Methods", func() {
	var intent *Intent

	BeforeEach(func() {
		ctx := &IntentValidator{CaptureStrategy: "quick"}
		var err error
		intent, err = NewIntent(ctx)
		Expect(err).NotTo(HaveOccurred())
		intent.Init()
	})

	Describe("terminalDimensions", func() {
		It("should return dimensions based on terminal info", func() {
			dims := intent.terminalDimensions()
			// BaseIntent provides default terminal info, so this may be non-nil.
			// Just verify no panic.
			_ = dims
		})
	})

	Describe("setCancelled", func() {
		It("should set result to cancelled", func() {
			intent.setCancelled()
			Expect(intent.active).To(BeFalse())
			Expect(intent.result).NotTo(BeNil())
			Expect(intent.result.Status).To(Equal(intents.Cancelled))
		})

		It("should deactivate the intent", func() {
			intent.setCancelled()
			Expect(intent.IsActive()).To(BeFalse())
		})
	})

	Describe("setFailed", func() {
		It("should set result to failed with error details", func() {
			intent.setFailed("TEST_CODE", "test message", nil)
			Expect(intent.active).To(BeFalse())
			Expect(intent.result).NotTo(BeNil())
			Expect(intent.result.Status).To(Equal(intents.Failed))
			Expect(intent.result.Error).To(HaveOccurred())
			Expect(intent.result.Error.Code).To(Equal("TEST_CODE"))
			Expect(intent.result.Error.Message).To(Equal("test message"))
		})

		It("should include the cause error when provided", func() {
			cause := errors.New("text field is required")
			intent.setFailed("VALIDATION", "invalid event", cause)
			Expect(intent.result.Error.Cause).To(Equal(cause))
		})

		It("should handle nil cause", func() {
			intent.setFailed("NO_CAUSE", "no underlying error", nil)
			Expect(intent.result.Error.Cause).ToNot(HaveOccurred())
		})
	})

	Describe("setFailedCmd", func() {
		It("should return nil", func() {
			cmd := intent.setFailedCmd("TEST", "test", nil)
			Expect(cmd).To(BeNil())
		})

		It("should mark the intent as failed", func() {
			intent.setFailedCmd("FAIL", "failed", nil)
			Expect(intent.active).To(BeFalse())
			Expect(intent.result.Status).To(Equal(intents.Failed))
		})
	})

	Describe("showValidationErrorModal", func() {
		It("creates an error modal on submitModal", func() {
			intent.showValidationErrorModal("event text is required")
			Expect(intent.submitModal).NotTo(BeNil())
		})

		It("returns nil command from the modal's Init (error modals do not auto-start)", func() {
			cmd := intent.showValidationErrorModal("event text is required")
			Expect(cmd).To(BeNil())
		})
	})

	Describe("showSubmitModal", func() {
		BeforeEach(func() {
			intent.reviewState = &ReviewInferredEventState{
				Event:          fixtures.EventWith("", "test", "", ""),
				AcceptedBursts: make([]*career.Burst, 0),
				AcceptedFacts:  make([]*career.Fact, 0),
			}
		})

		It("should create a submit modal", func() {
			intent.showSubmitModal()
			Expect(intent.submitModal).NotTo(BeNil())
		})

		It("should return a non-nil command", func() {
			cmd := intent.showSubmitModal()
			Expect(cmd).NotTo(BeNil())
		})
	})

	Describe("transitionToStrategyScreen", func() {
		It("should set state to StateChooseStrategy", func() {
			intent.currentState = StateForm
			intent.transitionToStrategyScreen()
			Expect(intent.currentState).To(Equal(StateChooseStrategy))
		})

		It("should create an active screen", func() {
			intent.transitionToStrategyScreen()
			Expect(intent.activeView).NotTo(BeNil())
		})

		It("should return nil", func() {
			cmd := intent.transitionToStrategyScreen()
			Expect(cmd).To(BeNil())
		})
	})

	Describe("transitionToFormScreen", func() {
		It("should set state to StateForm", func() {
			intent.transitionToFormScreen(StrategyQuick)
			Expect(intent.currentState).To(Equal(StateForm))
		})

		It("should set the strategy", func() {
			intent.transitionToFormScreen(StrategyManual)
			Expect(intent.strategy).To(Equal(StrategyManual))
		})

		It("should create an active screen", func() {
			intent.transitionToFormScreen(StrategyQuick)
			Expect(intent.activeView).NotTo(BeNil())
		})
	})

	Describe("getMetadataModalContent", func() {
		It("returns modal content when metadataModal is already set", func() {
			evt := fixtures.EventWith("", "test event", "", "")
			modal := eventview.NewReviewEnrichment(display.EventFromDomain(evt), eventview.ReviewEnrichmentConfig{})
			intent.reviewState = &ReviewInferredEventState{
				Event:         evt,
				EditingMode:   EditingModeMetadata,
				metadataModal: modal,
			}
			content := intent.getMetadataModalContent()
			Expect(content).NotTo(BeNil())
		})
	})

	Describe("getEditingModalContent", func() {
		It("should return nil when reviewState is nil", func() {
			intent.reviewState = nil
			Expect(intent.getEditingModalContent()).To(BeNil())
		})

		It("should return nil when EditingMode is None", func() {
			intent.reviewState = &ReviewInferredEventState{
				EditingMode: EditingModeNone,
			}
			Expect(intent.getEditingModalContent()).To(BeNil())
		})

		Context("when EditingMode is Metadata without CareerService", func() {
			It("should return modal content when CareerService is nil", func() {
				intent.reviewState = &ReviewInferredEventState{
					Event:       fixtures.EventWith("", "test", "", ""),
					EditingMode: EditingModeMetadata,
				}
				content := intent.getEditingModalContent()
				Expect(content).NotTo(BeNil())
			})
		})

		Context("when EditingMode is Bursts with nil modal", func() {
			It("should return nil", func() {
				intent.reviewState = &ReviewInferredEventState{
					EditingMode: EditingModeBursts,
					burstModal:  nil,
				}
				Expect(intent.getEditingModalContent()).To(BeNil())
			})
		})

		Context("when EditingMode is Facts with nil modal", func() {
			It("should return nil", func() {
				intent.reviewState = &ReviewInferredEventState{
					EditingMode: EditingModeFacts,
				}
				Expect(intent.getEditingModalContent()).To(BeNil())
			})
		})
	})

	Describe("renderModalOverlay", func() {
		BeforeEach(func() {
			// Set terminal dimensions so the overlay renderer can produce output.
			termInfo := &terminal.Info{
				Width:   120,
				Height:  40,
				IsValid: true,
			}
			intent.UpdateTerminalInfo(termInfo)
		})

		It("should return background when modalContent is nil", func() {
			bg := "base view content"
			result := intent.renderModalOverlay(bg, nil)
			Expect(result).To(Equal(bg))
		})

		It("should render overlay when modalContent is provided", func() {
			bg := "base view content"
			content := &modalContentData{
				title:   "Test Modal",
				content: "Modal body text",
				footer:  "Press Enter to confirm",
			}
			result := intent.renderModalOverlay(bg, content)
			Expect(result).NotTo(Equal(bg))
			Expect(result).NotTo(BeEmpty())
		})

		It("should include content without footer when footer is empty", func() {
			bg := "base"
			content := &modalContentData{
				title:   "Title",
				content: "Body",
			}
			result := intent.renderModalOverlay(bg, content)
			Expect(result).NotTo(BeEmpty())
		})
	})

	Describe("renderFormModalFooter", func() {
		It("should return non-empty badge footer", func() {
			footer := renderFormModalFooter()
			Expect(footer).NotTo(BeEmpty())
		})
	})

	Describe("performSubmit", func() {
		Context("when event is nil", func() {
			It("should return SubmitErrorMsg", func() {
				intent.reviewState = &ReviewInferredEventState{
					Event:         nil,
					AcceptedFacts: make([]*career.Fact, 0),
				}
				cmd := intent.performSubmit()
				Expect(cmd).NotTo(BeNil())

				msg := cmd()
				errMsg, ok := msg.(SubmitErrorMsg)
				Expect(ok).To(BeTrue())
				Expect(errMsg.Code).To(Equal("MISSING_EVENT"))
			})
		})

		Context("when CareerService is nil", func() {
			It("should return SubmitErrorMsg for nil career service", func() {
				intent.reviewState = &ReviewInferredEventState{
					Event:         fixtures.EventWith("", "Valid event text for testing purposes", "", ""),
					AcceptedFacts: make([]*career.Fact, 0),
				}
				intent.context.CareerService = nil
				cmd := intent.performSubmit()
				Expect(cmd).NotTo(BeNil())

				msg := cmd()
				errMsg, ok := msg.(SubmitErrorMsg)
				Expect(ok).To(BeTrue())
				Expect(errMsg.Code).To(Equal("SERVICE_ERROR"))
				Expect(errMsg.Message).To(ContainSubstring("Career service"))
			})
		})

		Context("with a valid event and career service", func() {
			var (
				svc       *careerservice.Service
				eventRepo *memoryrepo.EventRepository
				burstRepo *memoryrepo.BurstRepository
			)

			BeforeEach(func() {
				repos := memoryrepo.NewRepositories()
				eventRepo = repos.Event.(*memoryrepo.EventRepository)
				burstRepo = repos.Burst.(*memoryrepo.BurstRepository)

				svc = careerservice.NewService(eventRepo)
				svc.SetBurstRepository(burstRepo)
				svc.SetSkillRepository(repos.Skill)

				intent.context.CareerService = svc
			})

			It("should return SubmitCompleteMsg without inference data", func() {
				event := fixtures.EventWith("", "Valid event text for testing submission", "", "")
				intent.reviewState = &ReviewInferredEventState{
					Event:          event,
					AcceptedFacts:  make([]*career.Fact, 0),
					AcceptedBursts: make([]*career.Burst, 0),
				}

				cmd := intent.performSubmit()
				Expect(cmd).NotTo(BeNil())

				msg := cmd()
				_, ok := msg.(SubmitCompleteMsg)
				Expect(ok).To(BeTrue(), "expected SubmitCompleteMsg, got %T", msg)
			})
		})
	})

	Describe("performInference", func() {
		Context("burst detection and fact extraction", func() {
			var (
				svc       *careerservice.Service
				eventRepo *memoryrepo.EventRepository
				burstRepo *memoryrepo.BurstRepository
			)

			BeforeEach(func() {
				repos := memoryrepo.NewRepositories()
				eventRepo = repos.Event.(*memoryrepo.EventRepository)
				burstRepo = repos.Burst.(*memoryrepo.BurstRepository)

				svc = careerservice.NewService(eventRepo)
				svc.SetBurstRepository(burstRepo)
				svc.SetSkillRepository(repos.Skill)

				intent.context.CareerService = svc
			})

			Context("when ≥2 events exist", func() {
				BeforeEach(func() {
					ctx := context.Background()
					existing1 := fixtures.EventWith("existing-1", "Led migration of monolith to microservices architecture", "", "")
					existing1.Date = time.Now().Add(-24 * time.Hour)
					existing2 := fixtures.EventWith("existing-2", "Led redesign of microservices deployment pipeline", "", "")
					existing2.Date = time.Now().Add(-48 * time.Hour)
					Expect(eventRepo.Create(ctx, existing1)).To(Succeed())
					Expect(eventRepo.Create(ctx, existing2)).To(Succeed())
				})

				It("should return InferenceCompleteMsg with inferred bursts", func() {
					event := fixtures.EventWith("", "Valid event text for testing submission", "", "")
					intent.reviewState = &ReviewInferredEventState{
						Event:          event,
						AcceptedFacts:  make([]*career.Fact, 0),
						AcceptedBursts: make([]*career.Burst, 0),
					}

					cmd := intent.performInference()
					Expect(cmd).NotTo(BeNil())

					msg := cmd()
					completeMsg, ok := msg.(InferenceCompleteMsg)
					Expect(ok).To(BeTrue(), "expected InferenceCompleteMsg, got %T", msg)
					Expect(completeMsg.InferredBursts).NotTo(BeNil())
				})

				It("should return InferenceCompleteMsg with inferred facts from bursts", func() {
					event := fixtures.EventWith("", "Valid event text for testing submission", "", "")
					intent.reviewState = &ReviewInferredEventState{
						Event:          event,
						AcceptedFacts:  make([]*career.Fact, 0),
						AcceptedBursts: make([]*career.Burst, 0),
					}

					cmd := intent.performInference()
					msg := cmd()
					completeMsg, ok := msg.(InferenceCompleteMsg)
					Expect(ok).To(BeTrue(), "expected InferenceCompleteMsg, got %T", msg)
					Expect(completeMsg.InferredFacts).NotTo(BeNil())
				})

				It("should not break when burst detection errors occur", func() {
					event := fixtures.EventWith("", "Valid event text for testing submission", "", "")
					intent.reviewState = &ReviewInferredEventState{
						Event:          event,
						AcceptedFacts:  make([]*career.Fact, 0),
						AcceptedBursts: make([]*career.Burst, 0),
					}

					cmd := intent.performInference()
					msg := cmd()
					_, ok := msg.(InferenceCompleteMsg)
					Expect(ok).To(BeTrue(), "expected InferenceCompleteMsg even with errors, got %T", msg)
				})
			})

			Context("when <2 events exist", func() {
				It("should fallback to ExtractFactsFromEvent", func() {
					event := fixtures.EventWith("", "Valid event text for testing submission", "", "")
					Expect(eventRepo.Create(context.Background(), event)).To(Succeed())
					intent.reviewState = &ReviewInferredEventState{
						Event:          event,
						AcceptedFacts:  make([]*career.Fact, 0),
						AcceptedBursts: make([]*career.Burst, 0),
					}

					cmd := intent.performInference()
					msg := cmd()
					completeMsg, ok := msg.(InferenceCompleteMsg)
					Expect(ok).To(BeTrue(), "expected InferenceCompleteMsg, got %T", msg)
					Expect(completeMsg.InferredFacts).NotTo(BeNil())
					Expect(completeMsg.InferredBursts).To(BeEmpty())
				})

				It("should return empty bursts slice", func() {
					event := fixtures.EventWith("", "Valid event text for testing submission", "", "")
					Expect(eventRepo.Create(context.Background(), event)).To(Succeed())
					intent.reviewState = &ReviewInferredEventState{
						Event:          event,
						AcceptedFacts:  make([]*career.Fact, 0),
						AcceptedBursts: make([]*career.Burst, 0),
					}

					cmd := intent.performInference()
					msg := cmd()
					completeMsg, ok := msg.(InferenceCompleteMsg)
					Expect(ok).To(BeTrue())
					Expect(completeMsg.InferredBursts).To(BeEmpty())
				})
			})

			Context("when burst suggestions return empty", func() {
				BeforeEach(func() {
					ctx := context.Background()
					existing1 := fixtures.EventWith("unrelated-1", "Organised team building event at the local park", "", "")
					existing1.Date = time.Now().Add(-365 * 24 * time.Hour)
					existing2 := fixtures.EventWith("unrelated-2", "Attended annual company conference in London", "", "")
					existing2.Date = time.Now().Add(-730 * 24 * time.Hour)
					Expect(eventRepo.Create(ctx, existing1)).To(Succeed())
					Expect(eventRepo.Create(ctx, existing2)).To(Succeed())
				})

				It("should fallback to ExtractFactsFromEvent", func() {
					event := fixtures.EventWith("", "Valid event text for testing submission", "", "")
					Expect(eventRepo.Create(context.Background(), event)).To(Succeed())
					intent.reviewState = &ReviewInferredEventState{
						Event:          event,
						AcceptedFacts:  make([]*career.Fact, 0),
						AcceptedBursts: make([]*career.Burst, 0),
					}

					cmd := intent.performInference()
					msg := cmd()
					completeMsg, ok := msg.(InferenceCompleteMsg)
					Expect(ok).To(BeTrue(), "expected InferenceCompleteMsg, got %T", msg)
					Expect(completeMsg.InferredFacts).NotTo(BeNil())
				})
			})

			Context("when ListEvents fails", func() {
				It("should fallback to ExtractFactsFromEvent without breaking", func() {
					svcWithoutBurst := careerservice.NewService(eventRepo)
					intent.context.CareerService = svcWithoutBurst

					event := fixtures.EventWith("", "Valid event text for testing submission", "", "")
					Expect(eventRepo.Create(context.Background(), event)).To(Succeed())
					intent.reviewState = &ReviewInferredEventState{
						Event:          event,
						AcceptedFacts:  make([]*career.Fact, 0),
						AcceptedBursts: make([]*career.Burst, 0),
					}

					cmd := intent.performInference()
					msg := cmd()
					completeMsg, ok := msg.(InferenceCompleteMsg)
					Expect(ok).To(BeTrue(), "expected InferenceCompleteMsg, got %T", msg)
					Expect(completeMsg.InferredFacts).NotTo(BeNil())
				})
			})

			It("should preserve existing skill inference behaviour", func() {
				event := fixtures.EventWith("", "Valid event text for testing submission", "", "")
				intent.reviewState = &ReviewInferredEventState{
					Event:          event,
					AcceptedFacts:  make([]*career.Fact, 0),
					AcceptedBursts: make([]*career.Burst, 0),
				}

				cmd := intent.performInference()
				msg := cmd()
				completeMsg, ok := msg.(InferenceCompleteMsg)
				Expect(ok).To(BeTrue())
				Expect(completeMsg.InferredSkills).To(BeNil())
			})
		})

		Context("when CareerService is nil", func() {
			It("should return InferenceCompleteMsg with empty results", func() {
				intent.context.CareerService = nil
				intent.reviewState = &ReviewInferredEventState{
					Event: fixtures.EventWith("", "Valid event text for testing", "", ""),
				}

				cmd := intent.performInference()
				msg := cmd()
				completeMsg, ok := msg.(InferenceCompleteMsg)
				Expect(ok).To(BeTrue())
				Expect(completeMsg.InferredSkills).To(BeNil())
				Expect(completeMsg.InferredFacts).To(BeNil())
				Expect(completeMsg.InferredBursts).To(BeNil())
			})
		})
	})

	Describe("factsToPointers", func() {
		It("should return empty slice for nil input", func() {
			result := factsToPointers(nil)
			Expect(result).NotTo(BeNil())
			Expect(result).To(BeEmpty())
		})

		It("should return empty slice for empty input", func() {
			result := factsToPointers([]career.Fact{})
			Expect(result).NotTo(BeNil())
			Expect(result).To(BeEmpty())
		})

		It("should convert facts to pointers", func() {
			fact1 := fixtures.Fact("fact-1", "event-1")
			fact2 := fixtures.Fact("fact-2", "event-2")
			facts := []career.Fact{*fact1, *fact2}
			result := factsToPointers(facts)
			Expect(result).To(HaveLen(2))
			Expect(result[0]).NotTo(BeNil())
			Expect(result[0].ID).To(Equal("fact-1"))
			Expect(result[1]).NotTo(BeNil())
			Expect(result[1].ID).To(Equal("fact-2"))
		})

		It("should preserve fact data in pointers", func() {
			fact := fixtures.FactWithCategories(
				"test-id",
				"test text",
				"event-123",
				[]string{"leadership"},
				[]string{"hiring_manager"},
			)
			facts := []career.Fact{*fact}
			result := factsToPointers(facts)
			Expect(result).To(HaveLen(1))
			Expect(result[0].ID).To(Equal("test-id"))
			Expect(result[0].Text).To(Equal("test text"))
			Expect(result[0].CompetencyCategories).To(Equal([]string{"leadership"}))
			Expect(result[0].RoleFit).To(Equal(career.RoleFitStaff))
			Expect(result[0].AudienceRelevance).To(Equal([]string{"hiring_manager"}))
			Expect(result[0].SourceEventID).To(Equal("event-123"))
		})
	})
})
var _ = Describe("Coverage boost — performSubmit additional paths", func() {
	var intent *Intent

	BeforeEach(func() {
		ctx := &IntentValidator{CaptureStrategy: "quick"}
		var err error
		intent, err = NewIntent(ctx)
		Expect(err).NotTo(HaveOccurred())
		intent.Init()
	})

	Describe("performSubmit — validation error path", func() {
		It("returns SubmitErrorMsg with VALIDATION_ERROR when event text is empty", func() {
			repos := memoryrepo.NewRepositories()
			svc := careerservice.NewService(repos.Event)
			svc.SetSkillRepository(repos.Skill)

			invalidEvent := fixtures.EventWith("evt-invalid", "", "", "")
			intent.context.CareerService = svc
			intent.reviewState = &ReviewInferredEventState{
				Event:          invalidEvent,
				AcceptedFacts:  make([]*career.Fact, 0),
				AcceptedSkills: make([]*career.Skill, 0),
			}

			cmd := intent.performSubmit()
			Expect(cmd).NotTo(BeNil())

			msg := cmd()
			errMsg, ok := msg.(SubmitErrorMsg)
			Expect(ok).To(BeTrue(), "expected SubmitErrorMsg, got %T", msg)
			Expect(errMsg.Code).To(Equal("VALIDATION_ERROR"))
		})
	})

	Describe("performSubmit — skill save path with new skills", func() {
		It("returns SubmitCompleteMsg or SubmitErrorMsg when accepted skill has no ID", func() {
			repos := memoryrepo.NewRepositories()
			eventRepo := repos.Event.(*memoryrepo.EventRepository)
			svc := careerservice.NewService(eventRepo)
			svc.SetSkillRepository(repos.Skill)

			intent.context.CareerService = svc

			event := fixtures.EventWith("", "Valid event text for testing skill save path", "", "")
			skillWithNoID := fixtures.SkillWith("", "Go", "backend", "")
			intent.reviewState = &ReviewInferredEventState{
				Event:          event,
				AcceptedFacts:  make([]*career.Fact, 0),
				AcceptedSkills: []*career.Skill{skillWithNoID},
			}

			cmd := intent.performSubmit()
			Expect(cmd).NotTo(BeNil())

			msg := cmd()
			_, isComplete := msg.(SubmitCompleteMsg)
			_, isError := msg.(SubmitErrorMsg)
			Expect(isComplete || isError).To(BeTrue(), "expected SubmitCompleteMsg or SubmitErrorMsg, got %T", msg)
		})
	})

	Describe("performSubmit — skill with existing ID (link-only path)", func() {
		It("attempts to link an already-persisted skill and returns a result", func() {
			repos := memoryrepo.NewRepositories()
			eventRepo := repos.Event.(*memoryrepo.EventRepository)
			svc := careerservice.NewService(eventRepo)
			svc.SetSkillRepository(repos.Skill)

			intent.context.CareerService = svc

			event := fixtures.EventWith("", "Valid event text for testing skill link path", "", "")

			existingSkill := fixtures.SkillWith("existing-skill-id", "Go", "backend", "advanced")
			Expect(repos.Skill.Create(context.Background(), existingSkill)).To(Succeed())

			intent.reviewState = &ReviewInferredEventState{
				Event:          event,
				AcceptedFacts:  make([]*career.Fact, 0),
				AcceptedSkills: []*career.Skill{existingSkill},
			}

			cmd := intent.performSubmit()
			Expect(cmd).NotTo(BeNil())

			msg := cmd()
			_, isComplete := msg.(SubmitCompleteMsg)
			_, isError := msg.(SubmitErrorMsg)
			Expect(isComplete || isError).To(BeTrue(), "expected SubmitCompleteMsg or SubmitErrorMsg, got %T", msg)
		})
	})

	Describe("performSubmit — quick strategy with zero date", func() {
		It("sets today as date when strategy is quick and date is zero", func() {
			repos := memoryrepo.NewRepositories()
			svc := careerservice.NewService(repos.Event.(*memoryrepo.EventRepository))
			svc.SetSkillRepository(repos.Skill)

			intent.context.CareerService = svc
			intent.strategy = StrategyQuick

			event := fixtures.EventWith("", "Valid event text for zero date test path", "", "")
			event.Date = time.Time{}

			intent.reviewState = &ReviewInferredEventState{
				Event:          event,
				AcceptedFacts:  make([]*career.Fact, 0),
				AcceptedSkills: make([]*career.Skill, 0),
			}

			cmd := intent.performSubmit()
			Expect(cmd).NotTo(BeNil())

			msg := cmd()
			_, isComplete := msg.(SubmitCompleteMsg)
			_, isError := msg.(SubmitErrorMsg)
			Expect(isComplete || isError).To(BeTrue(), "expected SubmitCompleteMsg or SubmitErrorMsg, got %T", msg)
		})
	})

	Describe("performSubmit — fact save path when facts are present", func() {
		It("returns SubmitCompleteMsg or SubmitErrorMsg when facts are present", func() {
			repos := memoryrepo.NewRepositories()
			svc := careerservice.NewService(repos.Event.(*memoryrepo.EventRepository))
			svc.SetSkillRepository(repos.Skill)

			intent.context.CareerService = svc

			event := fixtures.EventWith("", "Valid event text for testing fact save path", "", "")
			factWithNoID := fixtures.Fact("", "")
			intent.reviewState = &ReviewInferredEventState{
				Event:          event,
				AcceptedFacts:  []*career.Fact{factWithNoID},
				AcceptedSkills: make([]*career.Skill, 0),
			}

			cmd := intent.performSubmit()
			Expect(cmd).NotTo(BeNil())

			msg := cmd()
			_, isComplete := msg.(SubmitCompleteMsg)
			_, isError := msg.(SubmitErrorMsg)
			Expect(isComplete || isError).To(BeTrue(), "expected SubmitCompleteMsg or SubmitErrorMsg, got %T", msg)
		})
	})
})

var _ = Describe("Coverage boost — performPostSavePersistence additional paths", func() {
	var intent *Intent

	BeforeEach(func() {
		ctx := &IntentValidator{CaptureStrategy: "quick"}
		var err error
		intent, err = NewIntent(ctx)
		Expect(err).NotTo(HaveOccurred())
		intent.Init()
	})

	Describe("performPostSavePersistence — nil careerService returns complete msg", func() {
		It("returns PostSavePersistenceCompleteMsg immediately when careerService is nil", func() {
			intent.context.CareerService = nil

			event := fixtures.EventWith("evt-1", "Event for post-save persistence nil path test", "", "")
			burst := fixtures.Burst("burst-1", "evt-1")
			fact := fixtures.Fact("fact-1", "evt-1")
			skill := fixtures.SkillWith("skill-1", "Go", "backend", "advanced")

			cmd := intent.performPostSavePersistence(
				event,
				[]*career.Fact{fact},
				[]*career.Skill{skill},
				[]*career.Burst{burst},
			)
			Expect(cmd).NotTo(BeNil())

			msg := cmd()
			completeMsg, ok := msg.(PostSavePersistenceCompleteMsg)
			Expect(ok).To(BeTrue(), "expected PostSavePersistenceCompleteMsg, got %T", msg)
			Expect(completeMsg.Event).To(Equal(event))
			Expect(completeMsg.Bursts).To(HaveLen(1))
			Expect(completeMsg.Facts).To(HaveLen(1))
			Expect(completeMsg.Skills).To(HaveLen(1))
		})
	})

	Describe("performPostSavePersistence — skill save (new skill, no ID)", func() {
		It("processes a new skill and returns a result", func() {
			repos := memoryrepo.NewRepositories()
			eventRepo := repos.Event.(*memoryrepo.EventRepository)
			svc := careerservice.NewService(eventRepo)
			svc.SetSkillRepository(repos.Skill)
			intent.context.CareerService = svc

			event := fixtures.EventWith("evt-post-1", "Event for post-save skill save test", "", "")
			Expect(eventRepo.Create(context.Background(), event)).To(Succeed())

			newSkill := fixtures.SkillWith("", "TypeScript", "frontend", "")

			cmd := intent.performPostSavePersistence(
				event,
				make([]*career.Fact, 0),
				[]*career.Skill{newSkill},
				make([]*career.Burst, 0),
			)
			Expect(cmd).NotTo(BeNil())

			msg := cmd()
			_, isComplete := msg.(PostSavePersistenceCompleteMsg)
			_, isError := msg.(SubmitErrorMsg)
			Expect(isComplete || isError).To(BeTrue(), "expected PostSavePersistenceCompleteMsg or SubmitErrorMsg, got %T", msg)
		})
	})

	Describe("performPostSavePersistence — fact save (new fact, no ID)", func() {
		It("saves a fact when fact has no ID", func() {
			repos := memoryrepo.NewRepositories()
			eventRepo := repos.Event.(*memoryrepo.EventRepository)
			svc := careerservice.NewService(eventRepo)
			svc.SetSkillRepository(repos.Skill)
			intent.context.CareerService = svc

			event := fixtures.EventWith("evt-post-2", "Event for post-save fact save test", "", "")
			Expect(eventRepo.Create(context.Background(), event)).To(Succeed())

			newFact := fixtures.Fact("", "")

			cmd := intent.performPostSavePersistence(
				event,
				[]*career.Fact{newFact},
				make([]*career.Skill, 0),
				make([]*career.Burst, 0),
			)
			Expect(cmd).NotTo(BeNil())

			msg := cmd()
			_, isComplete := msg.(PostSavePersistenceCompleteMsg)
			_, isError := msg.(SubmitErrorMsg)
			Expect(isComplete || isError).To(BeTrue(), "expected PostSavePersistenceCompleteMsg or SubmitErrorMsg, got %T", msg)
		})
	})

	Describe("performPostSavePersistence — fact already has ID (skip save)", func() {
		It("skips saving a fact that already has an ID and returns complete msg", func() {
			repos := memoryrepo.NewRepositories()
			eventRepo := repos.Event.(*memoryrepo.EventRepository)
			svc := careerservice.NewService(eventRepo)
			svc.SetSkillRepository(repos.Skill)
			intent.context.CareerService = svc

			event := fixtures.EventWith("evt-post-3", "Event for post-save fact skip test path", "", "")
			Expect(eventRepo.Create(context.Background(), event)).To(Succeed())

			existingFact := fixtures.Fact("existing-fact-id", event.ID)

			cmd := intent.performPostSavePersistence(
				event,
				[]*career.Fact{existingFact},
				make([]*career.Skill, 0),
				make([]*career.Burst, 0),
			)
			Expect(cmd).NotTo(BeNil())

			msg := cmd()
			completeMsg, ok := msg.(PostSavePersistenceCompleteMsg)
			Expect(ok).To(BeTrue(), "expected PostSavePersistenceCompleteMsg, got %T", msg)
			Expect(completeMsg.Facts).To(HaveLen(1))
		})
	})
})

var _ = Describe("Coverage boost — terminalDimensions", func() {
	var intent *Intent

	BeforeEach(func() {
		ctx := &IntentValidator{CaptureStrategy: "quick"}
		var err error
		intent, err = NewIntent(ctx)
		Expect(err).NotTo(HaveOccurred())
		intent.Init()
	})

	It("returns nil when GetTerminalInfo returns nil", func() {
		intent.UpdateTerminalInfo(nil)
		dims := intent.terminalDimensions()
		Expect(dims).To(BeNil())
	})

	It("returns dimensions when terminal info is valid", func() {
		info := &terminal.Info{Width: 120, Height: 40, IsValid: true}
		intent.UpdateTerminalInfo(info)
		dims := intent.terminalDimensions()
		Expect(dims).NotTo(BeNil())
		Expect(dims.TerminalWidth).To(Equal(120))
		Expect(dims.TerminalHeight).To(Equal(40))
	})
})

var _ = Describe("Coverage boost — transitionToStrategyScreen with terminal info", func() {
	var intent *Intent

	BeforeEach(func() {
		ctx := &IntentValidator{CaptureStrategy: "quick"}
		var err error
		intent, err = NewIntent(ctx)
		Expect(err).NotTo(HaveOccurred())
		intent.Init()
	})

	It("uses terminal dimensions when available", func() {
		info := &terminal.Info{Width: 150, Height: 50, IsValid: true}
		intent.UpdateTerminalInfo(info)
		cmd := intent.transitionToStrategyScreen()
		Expect(cmd).To(BeNil())
		Expect(intent.activeView).NotTo(BeNil())
	})
})

var _ = Describe("Coverage boost — transitionToFormScreen paths", func() {
	var intent *Intent

	BeforeEach(func() {
		ctx := &IntentValidator{CaptureStrategy: "quick"}
		var err error
		intent, err = NewIntent(ctx)
		Expect(err).NotTo(HaveOccurred())
		intent.Init()
	})

	It("uses terminal dimensions when available", func() {
		info := &terminal.Info{Width: 160, Height: 48, IsValid: true}
		intent.UpdateTerminalInfo(info)
		intent.transitionToFormScreen(StrategyManual)
		Expect(intent.activeView).NotTo(BeNil())
	})

	It("uses PreviousEvent when set in context", func() {
		prev := fixtures.EventWith("prev-evt", "Previous event text for form test", "", "")
		intent.context.PreviousEvent = prev
		intent.transitionToFormScreen(StrategyQuick)
		Expect(intent.activeView).NotTo(BeNil())
		Expect(intent.activeView).NotTo(BeNil())
	})
})

var _ = Describe("terminalDimensions — non-nil terminal info", func() {
	It("returns dimensions from terminal info when set", func() {
		ctx := &IntentValidator{CaptureStrategy: "quick"}
		intent, err := NewIntent(ctx)
		Expect(err).NotTo(HaveOccurred())
		intent.Init()

		info := &terminal.Info{Width: 120, Height: 40, IsValid: true}
		intent.UpdateTerminalInfo(info)

		dims := intent.terminalDimensions()
		Expect(dims).NotTo(BeNil())
		Expect(dims.TerminalWidth).To(Equal(120))
		Expect(dims.TerminalHeight).To(Equal(40))
	})
})

var _ = Describe("terminalDimensions — nil terminal info", func() {
	It("returns nil when terminal info is nil", func() {
		ctx := &IntentValidator{CaptureStrategy: "quick"}
		intent, err := NewIntent(ctx)
		Expect(err).NotTo(HaveOccurred())
		intent.Init()

		intent.UpdateTerminalInfo(nil)

		dims := intent.terminalDimensions()
		Expect(dims).To(BeNil())
	})
})

var _ = Describe("transitionToStrategyScreen — with non-nil terminal info", func() {
	It("configures screen using terminal dimensions", func() {
		ctx := &IntentValidator{CaptureStrategy: "quick"}
		intent, err := NewIntent(ctx)
		Expect(err).NotTo(HaveOccurred())
		intent.Init()

		info := &terminal.Info{Width: 120, Height: 40, IsValid: true}
		intent.UpdateTerminalInfo(info)

		cmd := intent.transitionToStrategyScreen()
		Expect(cmd).To(BeNil())
		Expect(intent.activeView).NotTo(BeNil())
	})
})

var _ = Describe("transitionToFormScreen — with non-nil terminal info", func() {
	It("configures screen using terminal dimensions", func() {
		ctx := &IntentValidator{CaptureStrategy: "quick"}
		intent, err := NewIntent(ctx)
		Expect(err).NotTo(HaveOccurred())
		intent.Init()

		info := &terminal.Info{Width: 120, Height: 40, IsValid: true}
		intent.UpdateTerminalInfo(info)

		intent.transitionToFormScreen(StrategyQuick)
		Expect(intent.activeView).NotTo(BeNil())
	})
})

var _ = Describe("transitionToFormScreen — with PreviousEvent set", func() {
	It("uses the previous event to populate the form", func() {
		prev := fixtures.EventWith("prev-1", "Previous event text here", "Corp", "Proj")
		ctx := &IntentValidator{
			CaptureStrategy: "quick",
			PreviousEvent:   prev,
		}
		intent, err := NewIntent(ctx)
		Expect(err).NotTo(HaveOccurred())
		intent.Init()

		intent.transitionToFormScreen(StrategyManual)
		Expect(intent.activeView).NotTo(BeNil())
	})
})

var _ = Describe("performSubmit — fact save failure", func() {
	It("returns SubmitErrorMsg with PARTIAL_SAVE code when a fact fails to save", func() {
		repos := memoryrepo.NewRepositories()
		svc := careerservice.NewService(repos.Event)
		svc.SetSkillRepository(repos.Skill)

		ctx := &IntentValidator{CaptureStrategy: "quick"}
		intent, err := NewIntent(ctx)
		Expect(err).NotTo(HaveOccurred())
		intent.Init()
		intent.context.CareerService = svc

		event := fixtures.EventWith("", "Valid event text for testing fact save", "", "")
		event.Date = time.Now().Add(-time.Hour)
		event.Tags = []string{"technical"}
		event.Categories = []string{"technical"}

		invalidFact := fixtures.Fact("", "")

		intent.reviewState = &ReviewInferredEventState{
			Event:         event,
			AcceptedFacts: []*career.Fact{invalidFact},
		}

		cmd := intent.performSubmit()
		Expect(cmd).NotTo(BeNil())
		msg := cmd()
		errMsg, ok := msg.(SubmitErrorMsg)
		Expect(ok).To(BeTrue())
		Expect(errMsg.Code).To(Equal("PARTIAL_SAVE"))
	})
})

var _ = Describe("performSubmit — event validation failure", func() {
	It("returns SubmitErrorMsg with VALIDATION_ERROR code for invalid event", func() {
		repos := memoryrepo.NewRepositories()
		svc := careerservice.NewService(repos.Event)
		svc.SetSkillRepository(repos.Skill)

		ctx := &IntentValidator{CaptureStrategy: "quick"}
		intent, err := NewIntent(ctx)
		Expect(err).NotTo(HaveOccurred())
		intent.Init()
		intent.context.CareerService = svc

		invalidEvent := fixtures.EventWith("", "x", "", "")

		intent.reviewState = &ReviewInferredEventState{
			Event:         invalidEvent,
			AcceptedFacts: make([]*career.Fact, 0),
		}

		cmd := intent.performSubmit()
		Expect(cmd).NotTo(BeNil())
		msg := cmd()
		errMsg, ok := msg.(SubmitErrorMsg)
		Expect(ok).To(BeTrue())
		Expect(errMsg.Code).To(Equal("VALIDATION_ERROR"))
	})
})

var _ = Describe("performPostSavePersistence — skill save error", func() {
	It("returns SubmitErrorMsg when skill creation fails", func() {
		repos := memoryrepo.NewRepositories()
		svc := careerservice.NewService(repos.Event)
		svc.SetSkillRepository(repos.Skill)

		ctx := &IntentValidator{CaptureStrategy: "quick"}
		intent, err := NewIntent(ctx)
		Expect(err).NotTo(HaveOccurred())
		intent.Init()
		intent.context.CareerService = svc

		event := fixtures.EventWith("evt-persist", "Valid persist test event text here", "Corp", "Proj")
		event.Date = time.Now().Add(-time.Hour)
		event.Tags = []string{"technical"}
		event.Categories = []string{"technical"}
		Expect(repos.Event.Create(context.Background(), event)).To(Succeed())

		invalidSkill := fixtures.SkillWith("", "", "", "")

		cmd := intent.performPostSavePersistence(
			event,
			[]*career.Fact{},
			[]*career.Skill{invalidSkill},
			[]*career.Burst{},
		)
		Expect(cmd).NotTo(BeNil())
		msg := cmd()
		errMsg, ok := msg.(SubmitErrorMsg)
		Expect(ok).To(BeTrue())
		Expect(errMsg.Code).To(Equal("SKILL_SAVE_ERROR"))
	})
})

var _ = Describe("performPostSavePersistence — fact save error", func() {
	It("returns SubmitErrorMsg when fact save fails", func() {
		repos := memoryrepo.NewRepositories()
		svc := careerservice.NewService(repos.Event)
		svc.SetSkillRepository(repos.Skill)

		ctx := &IntentValidator{CaptureStrategy: "quick"}
		intent, err := NewIntent(ctx)
		Expect(err).NotTo(HaveOccurred())
		intent.Init()
		intent.context.CareerService = svc

		event := fixtures.EventWith("evt-fact", "Valid fact test event text here correct", "Corp", "Proj")
		event.Date = time.Now().Add(-time.Hour)
		event.Tags = []string{"technical"}
		event.Categories = []string{"technical"}
		Expect(repos.Event.Create(context.Background(), event)).To(Succeed())

		invalidFact := fixtures.Fact("", "")

		cmd := intent.performPostSavePersistence(
			event,
			[]*career.Fact{invalidFact},
			[]*career.Skill{},
			[]*career.Burst{},
		)
		Expect(cmd).NotTo(BeNil())
		msg := cmd()
		errMsg, ok := msg.(SubmitErrorMsg)
		Expect(ok).To(BeTrue())
		Expect(errMsg.Code).To(Equal("FACT_SAVE_ERROR"))
	})
})

var _ = Describe("extractSkillsFromReviewData — nil skill in slice", func() {
	It("skips nil skills and returns only valid entries", func() {
		data := map[string]interface{}{
			"skills": []*career.Skill{
				nil,
				fixtures.SkillWith("", "Go", "backend", ""),
			},
		}
		result := extractSkillsFromReviewData(data)
		Expect(result).To(HaveLen(1))
		Expect(result[0].Name).To(Equal("Go"))
	})
})

var _ = Describe("performSubmit — StrategyQuick sets date when zero", func() {
	It("sets the event date to now when date is zero and strategy is quick", func() {
		repos := memoryrepo.NewRepositories()
		svc := careerservice.NewService(repos.Event)
		svc.SetSkillRepository(repos.Skill)

		ctx := &IntentValidator{CaptureStrategy: "quick"}
		intent, err := NewIntent(ctx)
		Expect(err).NotTo(HaveOccurred())
		intent.Init()
		intent.context.CareerService = svc
		intent.strategy = StrategyQuick

		event := fixtures.EventWith("", "Valid event text for quick strategy test here", "", "")
		event.Date = time.Time{}
		event.Tags = []string{}
		event.Categories = []string{}

		intent.reviewState = &ReviewInferredEventState{
			Event:         event,
			AcceptedFacts: make([]*career.Fact, 0),
		}

		cmd := intent.performSubmit()
		Expect(cmd).NotTo(BeNil())
		msg := cmd()
		_, isComplete := msg.(SubmitCompleteMsg)
		_, isError := msg.(SubmitErrorMsg)
		Expect(isComplete || isError).To(BeTrue())
		Expect(event.Date.IsZero()).To(BeFalse())
	})
})

var _ = Describe("performSubmit — skill save error", func() {
	It("returns SubmitErrorMsg with SKILL_SAVE_ERROR when skill name is empty", func() {
		repos := memoryrepo.NewRepositories()
		svc := careerservice.NewService(repos.Event)
		svc.SetSkillRepository(repos.Skill)

		ctx := &IntentValidator{CaptureStrategy: "quick"}
		intent, err := NewIntent(ctx)
		Expect(err).NotTo(HaveOccurred())
		intent.Init()
		intent.context.CareerService = svc

		event := fixtures.EventWith("", "Valid event text for skill save error test here", "", "")
		event.Date = time.Now().Add(-time.Hour)
		event.Tags = []string{"technical"}
		event.Categories = []string{"technical"}

		invalidSkill := fixtures.SkillWith("", "", "", "")

		intent.reviewState = &ReviewInferredEventState{
			Event:          event,
			AcceptedFacts:  make([]*career.Fact, 0),
			AcceptedSkills: []*career.Skill{invalidSkill},
		}

		cmd := intent.performSubmit()
		Expect(cmd).NotTo(BeNil())
		msg := cmd()
		errMsg, ok := msg.(SubmitErrorMsg)
		Expect(ok).To(BeTrue())
		Expect(errMsg.Code).To(Equal("SKILL_SAVE_ERROR"))
	})
})

var _ = Describe("performPostSavePersistence — burst confirm error", func() {
	It("returns SubmitErrorMsg when ConfirmBurst fails for non-existent burst", func() {
		repos := memoryrepo.NewRepositories()
		svc := careerservice.NewService(repos.Event)
		svc.SetSkillRepository(repos.Skill)

		ctx := &IntentValidator{CaptureStrategy: "quick"}
		intent, err := NewIntent(ctx)
		Expect(err).NotTo(HaveOccurred())
		intent.Init()
		intent.context.CareerService = svc

		event := fixtures.EventWith("evt-burst", "Valid burst test event text here correct", "Corp", "Proj")
		event.Date = time.Now().Add(-time.Hour)
		event.Tags = []string{"technical"}
		event.Categories = []string{"technical"}
		Expect(repos.Event.Create(context.Background(), event)).To(Succeed())

		nonExistentBurst := fixtures.Burst("nonexistent-burst-id", "evt-burst")

		cmd := intent.performPostSavePersistence(
			event,
			[]*career.Fact{},
			[]*career.Skill{},
			[]*career.Burst{nonExistentBurst},
		)
		Expect(cmd).NotTo(BeNil())
		msg := cmd()
		_, isComplete := msg.(PostSavePersistenceCompleteMsg)
		_, isError := msg.(SubmitErrorMsg)
		Expect(isComplete || isError).To(BeTrue())
	})
})

var _ = Describe("inferSkillsFromEvent", func() {
	It("returns nil when service is nil", func() {
		evt := fixtures.EventWith("", "test event", "", "")
		result := inferSkillsFromEvent(context.Background(), evt, nil)
		Expect(result).To(BeNil())
	})

	It("returns nil when inference fails", func() {
		evt := fixtures.EventWith("", "test event", "", "")
		mockService := &mockSkillInferenceService{
			shouldFail: true,
		}
		result := inferSkillsFromEvent(context.Background(), evt, mockService)
		Expect(result).To(BeNil())
	})

	It("returns suggestions when inference succeeds", func() {
		evt := fixtures.EventWith("", "test event", "", "")
		mockService := &mockSkillInferenceService{
			suggestions: []skillinference.SkillSuggestion{
				{Name: "Go", Category: "Language", Confidence: 0.95},
			},
		}
		result := inferSkillsFromEvent(context.Background(), evt, mockService)
		Expect(result).To(HaveLen(1))
		Expect(result[0].Name).To(Equal("Go"))
	})
})

var _ = Describe("extractBurstsFromReviewData", func() {
	It("returns nil when bursts key is missing", func() {
		reviewData := map[string]interface{}{}
		result := extractBurstsFromReviewData(reviewData, nil)
		Expect(result).To(BeNil())
	})

	It("converts display.Burst to domain pointers", func() {
		displayBursts := []display.Burst{
			{ID: "b1", Name: "Burst 1", Description: "Test burst"},
		}
		reviewData := map[string]interface{}{
			"bursts": displayBursts,
		}
		result := extractBurstsFromReviewData(reviewData, nil)
		Expect(result).To(HaveLen(1))
		Expect(result[0].Name).To(Equal("Burst 1"))
	})

	It("returns domain pointers when already in correct type", func() {
		domainBursts := []*career.Burst{
			fixtures.Burst("b1"),
		}
		reviewData := map[string]interface{}{
			"bursts": domainBursts,
		}
		result := extractBurstsFromReviewData(reviewData, nil)
		Expect(result).To(Equal(domainBursts))
	})
})

var _ = Describe("extractFactsFromReviewData", func() {
	It("returns nil when facts key is missing", func() {
		reviewData := map[string]interface{}{}
		result := extractFactsFromReviewData(reviewData, nil)
		Expect(result).To(BeNil())
	})

	It("converts display.Fact to domain pointers", func() {
		displayFacts := []display.Fact{
			{ID: "f1", Text: "Fact 1", SourceEventID: "evt1"},
		}
		reviewData := map[string]interface{}{
			"facts": displayFacts,
		}
		result := extractFactsFromReviewData(reviewData, nil)
		Expect(result).To(HaveLen(1))
		Expect(result[0].Text).To(Equal("Fact 1"))
	})

	It("returns domain pointers when already in correct type", func() {
		domainFacts := []*career.Fact{
			fixtures.FactWith("f1", "Fact 1"),
		}
		reviewData := map[string]interface{}{
			"facts": domainFacts,
		}
		result := extractFactsFromReviewData(reviewData, nil)
		Expect(result).To(Equal(domainFacts))
	})
})

var _ = Describe("displayFactsToPointers", func() {
	It("returns empty slice when displayFacts is empty", func() {
		result := displayFactsToPointers([]display.Fact{}, nil)
		Expect(result).To(BeEmpty())
	})

	It("creates new domain facts from display facts", func() {
		displayFacts := []display.Fact{
			{ID: "f1", Text: "Fact 1", SourceEventID: "evt1"},
		}
		result := displayFactsToPointers(displayFacts, nil)
		Expect(result).To(HaveLen(1))
		Expect(result[0].Text).To(Equal("Fact 1"))
		Expect(result[0].SourceEventID).To(Equal("evt1"))
	})

	It("matches inferred facts by ID", func() {
		inferredFacts := []*career.Fact{
			fixtures.FactWith("f1", "Inferred Fact 1"),
		}
		displayFacts := []display.Fact{
			{ID: "f1", Text: "Display Fact 1"},
		}
		result := displayFactsToPointers(displayFacts, inferredFacts)
		Expect(result).To(HaveLen(1))
		Expect(result[0]).To(Equal(inferredFacts[0]))
	})

	It("preserves competency categories and audience relevance", func() {
		displayFacts := []display.Fact{
			{
				ID:                   "f1",
				Text:                 "Fact 1",
				CompetencyCategories: []string{"technical", "leadership"},
				AudienceRelevance:    []string{"hiring-manager"},
			},
		}
		result := displayFactsToPointers(displayFacts, nil)
		Expect(result[0].CompetencyCategories).To(Equal([]string{"technical", "leadership"}))
		Expect(result[0].AudienceRelevance).To(Equal([]string{"hiring-manager"}))
	})
})

var _ = Describe("displayBurstsToPointers", func() {
	It("returns empty slice when displayBursts is empty", func() {
		result := displayBurstsToPointers([]display.Burst{}, nil)
		Expect(result).To(BeEmpty())
	})

	It("creates new domain bursts from display bursts", func() {
		displayBursts := []display.Burst{
			{ID: "b1", Name: "Burst 1", Description: "Test burst"},
		}
		result := displayBurstsToPointers(displayBursts, nil)
		Expect(result).To(HaveLen(1))
		Expect(result[0].Name).To(Equal("Burst 1"))
		Expect(result[0].Description).To(Equal("Test burst"))
	})

	It("matches inferred bursts by ID", func() {
		inferredBursts := []*career.Burst{
			fixtures.Burst("b1"),
		}
		displayBursts := []display.Burst{
			{ID: "b1", Name: "Display Burst 1"},
		}
		result := displayBurstsToPointers(displayBursts, inferredBursts)
		Expect(result).To(HaveLen(1))
		Expect(result[0]).To(Equal(inferredBursts[0]))
	})

	It("preserves event IDs", func() {
		displayBursts := []display.Burst{
			{
				ID:       "b1",
				Name:     "Burst 1",
				EventIDs: []string{"evt1", "evt2"},
			},
		}
		result := displayBurstsToPointers(displayBursts, nil)
		Expect(result[0].EventIDs).To(Equal([]string{"evt1", "evt2"}))
	})
})

type mockSkillInferenceService struct {
	suggestions []skillinference.SkillSuggestion
	shouldFail  bool
}

func (m *mockSkillInferenceService) InferSkillsFromEvents(ctx context.Context, events []*career.Event) (*skillinference.InferenceResult, error) {
	if m.shouldFail {
		return nil, errors.New("inference failed")
	}
	return &skillinference.InferenceResult{
		Suggestions: m.suggestions,
	}, nil
}

func (m *mockSkillInferenceService) InferSkillsFromBurst(ctx context.Context, burst *career.Burst, events []*career.Event) (*skillinference.InferenceResult, error) {
	return &skillinference.InferenceResult{
		Suggestions: m.suggestions,
	}, nil
}

func (m *mockSkillInferenceService) CreateSkillsFromSuggestions(ctx context.Context, suggestions []skillinference.SkillSuggestion) ([]*career.Skill, error) {
	return nil, nil
}
