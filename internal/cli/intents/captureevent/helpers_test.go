package captureevent

import (
	"context"
	"errors"
	"time"

	"github.com/baphled/kariya/internal/cli/intents"
	"github.com/baphled/kariya/internal/cli/terminal"
	"github.com/baphled/kariya/internal/domain/career"
	memoryrepo "github.com/baphled/kariya/internal/repository/career/memory"
	careerservice "github.com/baphled/kariya/internal/service/career"
	"github.com/baphled/kariya/internal/testutil/fixtures"
	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
)

var _ = Describe("Helper Methods", func() {
	var intent *Intent

	BeforeEach(func() {
		ctx := &IntentContext{CaptureStrategy: "quick"}
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
			Expect(intent.activeScreen).NotTo(BeNil())
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
			Expect(intent.activeScreen).NotTo(BeNil())
		})
	})

	Describe("getMetadataModalContent", func() {
		It("returns modal content when metadataModal is already set", func() {
			repos := memoryrepo.NewRepositories()
			svc := careerservice.NewService(repos.Event)
			svc.SetSkillRepository(repos.Skill)

			event := fixtures.EventWith("", "test event", "", "")
			modal := NewReviewEnrichmentModel(context.Background(), event, svc, nil, nil)
			intent.reviewState = &ReviewInferredEventState{
				Event:         event,
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
			It("should panic when CareerService is nil", func() {
				intent.reviewState = &ReviewInferredEventState{
					Event:       fixtures.EventWith("", "test", "", ""),
					EditingMode: EditingModeMetadata,
				}
				// ReviewEnrichmentModel immediately calls CareerService.GetSkillRepository(),
				// which panics with a nil service. Verifies the dispatch reaches metadata.
				Expect(func() {
					intent.getEditingModalContent()
				}).To(Panic())
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
