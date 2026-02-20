package captureevent

import (
	"errors"

	"github.com/baphled/kariya/internal/cli/intents"
	"github.com/baphled/kariya/internal/cli/terminal"
	"github.com/baphled/kariya/internal/domain/career"
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
				// MetadataEditorModelNew immediately calls CareerService.GetSkillRepository(),
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
					factModal:   nil,
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

	Describe("renderBurstModalFooter", func() {
		It("should return editing footer when editing is true", func() {
			footer := renderBurstModalFooter(true)
			Expect(footer).NotTo(BeEmpty())
		})

		It("should return navigation footer when editing is false", func() {
			footer := renderBurstModalFooter(false)
			Expect(footer).NotTo(BeEmpty())
		})

		It("should produce different output for editing vs navigation", func() {
			editFooter := renderBurstModalFooter(true)
			navFooter := renderBurstModalFooter(false)
			Expect(editFooter).NotTo(Equal(navFooter))
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
	})

	Describe("factsToPointers", func() {
		It("should return empty slice for nil input", func() {
			result := factsToPointers(nil)
			Expect(result).NotTo(BeNil())
			Expect(result).To(HaveLen(0))
		})

		It("should return empty slice for empty input", func() {
			result := factsToPointers([]career.Fact{})
			Expect(result).NotTo(BeNil())
			Expect(result).To(HaveLen(0))
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
