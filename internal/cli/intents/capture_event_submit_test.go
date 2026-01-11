package intents

import (
	"time"

	"github.com/baphled/kariya/internal/cli/service"
	"github.com/baphled/kariya/internal/domain/career"
	careerservice "github.com/baphled/kariya/internal/service/career"
	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
)

var _ = Describe("CaptureEvent - Submit Workflow", func() {
	var (
		intent         *CaptureEventIntent
		testCLIService *service.CLIEventService
	)

	BeforeEach(func() {
		testCLIService = &service.CLIEventService{}

		ctx := &CaptureEventContext{
			CLIEventService: testCLIService,
			CareerService:   &careerservice.Service{},
			CaptureStrategy: "manual",
			PreviousEvent:   nil,
		}

		var err error
		intent, err = NewCaptureEventIntent(ctx)
		Expect(err).NotTo(HaveOccurred())
		Expect(intent).NotTo(BeNil())
		intent.Init()
	})

	Describe("performSubmit validation", func() {
		Context("with missing event", func() {
			BeforeEach(func() {
				intent.state.currentState = CaptureStateSubmit
				intent.state.reviewState.Event = nil
			})

			It("should return MISSING_EVENT error", func() {
				cmd := intent.performSubmit()
				Expect(cmd).NotTo(BeNil())

				msg := cmd()
				errMsg, ok := msg.(SubmitErrorMsg)
				Expect(ok).To(BeTrue())
				Expect(errMsg.Code).To(Equal("MISSING_EVENT"))
				Expect(errMsg.Message).To(ContainSubstring("No event data"))
			})
		})

		Context("with invalid event data", func() {
			BeforeEach(func() {
				intent.state.currentState = CaptureStateSubmit
				intent.state.reviewState.Event = &career.CareerEvent{
					// Text is empty - should fail validation
					Text:      "",
					Date:      time.Now(),
					CreatedAt: time.Now(),
					UpdatedAt: time.Now(),
				}
			})

			It("should return VALIDATION_ERROR", func() {
				cmd := intent.performSubmit()
				Expect(cmd).NotTo(BeNil())

				msg := cmd()
				errMsg, ok := msg.(SubmitErrorMsg)
				Expect(ok).To(BeTrue())
				Expect(errMsg.Code).To(Equal("VALIDATION_ERROR"))
				Expect(errMsg.Message).To(ContainSubstring("validation failed"))
			})
		})

		Context("with nil event service", func() {
			BeforeEach(func() {
				intent.state.currentState = CaptureStateSubmit
				intent.state.reviewState.Event = &career.CareerEvent{
					ID:        "test-event-1",
					Text:      "Valid event text for testing purposes",
					Date:      time.Now(),
					CreatedAt: time.Now(),
					UpdatedAt: time.Now(),
				}
				intent.eventService = nil
			})

			It("should return SERVICE_ERROR", func() {
				cmd := intent.performSubmit()
				Expect(cmd).NotTo(BeNil())

				msg := cmd()
				errMsg, ok := msg.(SubmitErrorMsg)
				Expect(ok).To(BeTrue())
				Expect(errMsg.Code).To(Equal("SERVICE_ERROR"))
				Expect(errMsg.Message).To(ContainSubstring("not initialized"))
			})
		})
	})

	Describe("viewError", func() {
		BeforeEach(func() {
			intent.state.currentState = CaptureStateSubmit
		})

		Context("with error set", func() {
			BeforeEach(func() {
				intent.state.error = &IntentError{
					Code:    "TEST_ERROR",
					Message: "Something went wrong",
				}
			})

			It("should display error code", func() {
				view := intent.viewError()
				Expect(view).To(ContainSubstring("TEST_ERROR"))
			})

			It("should display error message", func() {
				view := intent.viewError()
				Expect(view).To(ContainSubstring("Something went wrong"))
			})

			It("should show retry instructions", func() {
				view := intent.viewError()
				Expect(view).To(ContainSubstring("'r' to retry"))
				Expect(view).To(ContainSubstring("Esc to cancel"))
			})
		})

		Context("with long error code", func() {
			BeforeEach(func() {
				intent.state.error = &IntentError{
					Code:    "THIS_IS_A_VERY_LONG_ERROR_CODE_THAT_SHOULD_BE_TRUNCATED_FOR_DISPLAY",
					Message: "Error message",
				}
			})

			It("should truncate long error code", func() {
				view := intent.viewError()
				Expect(view).To(ContainSubstring("..."))
			})
		})

		Context("with long error message", func() {
			BeforeEach(func() {
				intent.state.error = &IntentError{
					Code:    "ERROR",
					Message: "This is a very long error message that should be truncated for proper display in the terminal",
				}
			})

			It("should truncate long error message", func() {
				view := intent.viewError()
				Expect(view).To(ContainSubstring("..."))
			})
		})
	})

	Describe("acceptCurrentItem", func() {
		Context("when accepting a burst", func() {
			BeforeEach(func() {
				intent.state.reviewState.SelectedItemType = "burst"
				intent.state.reviewState.SelectedIndex = 0
				intent.state.reviewState.InferredBursts = []*career.Burst{
					{ID: "burst-1", Name: "First Burst", EventIDs: []string{"e1", "e2"}},
					{ID: "burst-2", Name: "Second Burst", EventIDs: []string{"e3", "e4"}},
				}
				intent.state.reviewState.AcceptedBursts = []*career.Burst{}
			})

			It("should move burst to accepted list", func() {
				intent.acceptCurrentItem()

				Expect(intent.state.reviewState.AcceptedBursts).To(HaveLen(1))
				Expect(intent.state.reviewState.AcceptedBursts[0].ID).To(Equal("burst-1"))
			})

			It("should remove burst from inferred list", func() {
				intent.acceptCurrentItem()

				Expect(intent.state.reviewState.InferredBursts).To(HaveLen(1))
				Expect(intent.state.reviewState.InferredBursts[0].ID).To(Equal("burst-2"))
			})

			It("should adjust selection index when at end of list", func() {
				intent.state.reviewState.SelectedIndex = 1
				intent.acceptCurrentItem()

				Expect(intent.state.reviewState.SelectedIndex).To(Equal(0))
			})
		})

		Context("when accepting a fact", func() {
			BeforeEach(func() {
				intent.state.reviewState.SelectedItemType = "fact"
				intent.state.reviewState.SelectedIndex = 0
				intent.state.reviewState.InferredFacts = []*career.Fact{
					{ID: "fact-1", Text: "First Fact", RoleFit: career.RoleFitSeniorIC, SourceEventID: "e1"},
					{ID: "fact-2", Text: "Second Fact", RoleFit: career.RoleFitSeniorIC, SourceEventID: "e2"},
				}
				intent.state.reviewState.AcceptedFacts = []*career.Fact{}
			})

			It("should move fact to accepted list", func() {
				intent.acceptCurrentItem()

				Expect(intent.state.reviewState.AcceptedFacts).To(HaveLen(1))
				Expect(intent.state.reviewState.AcceptedFacts[0].ID).To(Equal("fact-1"))
			})

			It("should remove fact from inferred list", func() {
				intent.acceptCurrentItem()

				Expect(intent.state.reviewState.InferredFacts).To(HaveLen(1))
				Expect(intent.state.reviewState.InferredFacts[0].ID).To(Equal("fact-2"))
			})
		})

		Context("with invalid index", func() {
			BeforeEach(func() {
				intent.state.reviewState.SelectedItemType = "burst"
				intent.state.reviewState.SelectedIndex = 5 // Out of bounds
				intent.state.reviewState.InferredBursts = []*career.Burst{
					{ID: "burst-1", Name: "First Burst", EventIDs: []string{"e1", "e2"}},
				}
			})

			It("should not panic with out of bounds index", func() {
				Expect(func() {
					intent.acceptCurrentItem()
				}).NotTo(Panic())
			})

			It("should not modify lists with invalid index", func() {
				intent.acceptCurrentItem()

				Expect(intent.state.reviewState.InferredBursts).To(HaveLen(1))
				Expect(intent.state.reviewState.AcceptedBursts).To(HaveLen(0))
			})
		})

		Context("with negative index", func() {
			BeforeEach(func() {
				intent.state.reviewState.SelectedItemType = "burst"
				intent.state.reviewState.SelectedIndex = -1
				intent.state.reviewState.InferredBursts = []*career.Burst{
					{ID: "burst-1", Name: "First Burst", EventIDs: []string{"e1", "e2"}},
				}
			})

			It("should not panic with negative index", func() {
				Expect(func() {
					intent.acceptCurrentItem()
				}).NotTo(Panic())
			})
		})
	})

	Describe("rejectCurrentItem", func() {
		Context("when rejecting a burst", func() {
			BeforeEach(func() {
				intent.state.reviewState.SelectedItemType = "burst"
				intent.state.reviewState.SelectedIndex = 0
				intent.state.reviewState.InferredBursts = []*career.Burst{
					{ID: "burst-1", Name: "First Burst", EventIDs: []string{"e1", "e2"}},
					{ID: "burst-2", Name: "Second Burst", EventIDs: []string{"e3", "e4"}},
				}
			})

			It("should remove burst from inferred list", func() {
				intent.rejectCurrentItem()

				Expect(intent.state.reviewState.InferredBursts).To(HaveLen(1))
				Expect(intent.state.reviewState.InferredBursts[0].ID).To(Equal("burst-2"))
			})

			It("should NOT add burst to accepted list", func() {
				intent.rejectCurrentItem()

				Expect(intent.state.reviewState.AcceptedBursts).To(HaveLen(0))
			})

			It("should track rejection reason", func() {
				intent.rejectCurrentItem()

				Expect(intent.state.reviewState.RejectedItems).To(HaveKey("burst-1"))
				Expect(intent.state.reviewState.RejectedItems["burst-1"]).To(Equal("user_rejected"))
			})

			It("should initialize RejectedItems map if nil", func() {
				intent.state.reviewState.RejectedItems = nil
				intent.rejectCurrentItem()

				Expect(intent.state.reviewState.RejectedItems).NotTo(BeNil())
			})
		})

		Context("when rejecting a fact", func() {
			BeforeEach(func() {
				intent.state.reviewState.SelectedItemType = "fact"
				intent.state.reviewState.SelectedIndex = 0
				intent.state.reviewState.InferredFacts = []*career.Fact{
					{ID: "fact-1", Text: "First Fact", RoleFit: career.RoleFitSeniorIC, SourceEventID: "e1"},
					{ID: "fact-2", Text: "Second Fact", RoleFit: career.RoleFitSeniorIC, SourceEventID: "e2"},
				}
			})

			It("should remove fact from inferred list", func() {
				intent.rejectCurrentItem()

				Expect(intent.state.reviewState.InferredFacts).To(HaveLen(1))
				Expect(intent.state.reviewState.InferredFacts[0].ID).To(Equal("fact-2"))
			})

			It("should track rejection reason", func() {
				intent.rejectCurrentItem()

				Expect(intent.state.reviewState.RejectedItems).To(HaveKey("fact-1"))
				Expect(intent.state.reviewState.RejectedItems["fact-1"]).To(Equal("user_rejected"))
			})
		})

		Context("with invalid index", func() {
			BeforeEach(func() {
				intent.state.reviewState.SelectedItemType = "fact"
				intent.state.reviewState.SelectedIndex = 10 // Out of bounds
				intent.state.reviewState.InferredFacts = []*career.Fact{
					{ID: "fact-1", Text: "First Fact", RoleFit: career.RoleFitSeniorIC, SourceEventID: "e1"},
				}
			})

			It("should not panic with out of bounds index", func() {
				Expect(func() {
					intent.rejectCurrentItem()
				}).NotTo(Panic())
			})
		})
	})

	Describe("initializeFormForEdit", func() {
		Context("with previous event", func() {
			var previousEvent *career.CareerEvent

			BeforeEach(func() {
				previousEvent = &career.CareerEvent{
					ID:        "existing-event",
					Text:      "Previous event text for editing purposes",
					Date:      time.Now().AddDate(0, -1, 0),
					Company:   "Previous Corp",
					Project:   "Previous Project",
					CreatedAt: time.Now().AddDate(0, -1, 0),
					UpdatedAt: time.Now(),
				}

				ctx := &CaptureEventContext{
					CLIEventService: testCLIService,
					CaptureStrategy: "manual",
					PreviousEvent:   previousEvent,
				}

				var err error
				intent, err = NewCaptureEventIntent(ctx)
				Expect(err).NotTo(HaveOccurred())
			})

			It("should set review state event to previous event", func() {
				intent.initializeFormForEdit()

				Expect(intent.state.reviewState.Event).To(Equal(previousEvent))
			})

			It("should set strategy to manual", func() {
				intent.initializeFormForEdit()

				Expect(intent.state.strategy).To(Equal(StrategyManual))
			})

			It("should show optional fields", func() {
				intent.initializeFormForEdit()

				Expect(intent.state.showOptionalFields).To(BeTrue())
			})

			It("should skip to form state", func() {
				intent.initializeFormForEdit()

				Expect(intent.state.currentState).To(Equal(CaptureStateForm))
			})
		})

		Context("without previous event", func() {
			BeforeEach(func() {
				intent.context.PreviousEvent = nil
			})

			It("should return a no-op command", func() {
				cmd := intent.initializeFormForEdit()
				Expect(cmd).NotTo(BeNil())

				// Execute the command and verify it returns nil
				msg := cmd()
				Expect(msg).To(BeNil())
			})
		})
	})

	Describe("Quick strategy date behavior", func() {
		// NOTE: The date defaulting for quick strategy happens AFTER the service check
		// in performSubmit(). This means if the service is nil, the date won't be set.
		// This test documents the current behavior.
		Context("with valid event service", func() {
			It("should require a valid event service to reach date defaulting logic", func() {
				// This test verifies that the nil service check happens before date defaulting
				intent.state.currentState = CaptureStateSubmit
				intent.state.strategy = StrategyQuick
				intent.state.reviewState.Event = &career.CareerEvent{
					ID:        "test-event-1",
					Text:      "Quick event text for testing purposes",
					Date:      time.Time{}, // Zero date
					CreatedAt: time.Now(),
					UpdatedAt: time.Now(),
				}
				intent.eventService = nil

				cmd := intent.performSubmit()
				msg := cmd()

				// Should fail with SERVICE_ERROR before reaching date defaulting
				errMsg, ok := msg.(SubmitErrorMsg)
				Expect(ok).To(BeTrue())
				Expect(errMsg.Code).To(Equal("SERVICE_ERROR"))

				// Date remains zero because service check happens first
				Expect(intent.state.reviewState.Event.Date.IsZero()).To(BeTrue())
			})
		})
	})

})
