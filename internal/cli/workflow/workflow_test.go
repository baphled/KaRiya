package workflow_test

import (
	"github.com/baphled/kariya/internal/cli/workflow"
	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
)

var _ = Describe("WorkflowState", func() {
	var ws *workflow.WorkflowState

	BeforeEach(func() {
		ws = workflow.NewWorkflowState()
	})

	Describe("NewWorkflowState", func() {
		It("should create a new workflow state", func() {
			Expect(ws).NotTo(BeNil())
			Expect(ws.GetCurrentStep()).To(Equal(workflow.StepCapture))
			Expect(ws.GetProgress()).To(Equal(0))
		})
	})

	Describe("NewWorkflowStateFromCapture", func() {
		It("should create workflow state starting from capture", func() {
			eventIDs := []string{"event1", "event2"}
			ws = workflow.NewWorkflowStateFromCapture(eventIDs)
			
			Expect(ws.GetCurrentStep()).To(Equal(workflow.StepCapture))
			Expect(ws.IsStepCompleted(workflow.StepCapture)).To(BeTrue())
			Expect(ws.GetEventIDs()).To(Equal(eventIDs))
		})
	})

	Describe("NewWorkflowStateFromImport", func() {
		It("should create workflow state starting from import", func() {
			eventIDs := []string{"event1", "event2", "event3"}
			ws = workflow.NewWorkflowStateFromImport(eventIDs)
			
			Expect(ws.GetCurrentStep()).To(Equal(workflow.StepImport))
			Expect(ws.IsStepCompleted(workflow.StepImport)).To(BeTrue())
			Expect(ws.GetEventIDs()).To(Equal(eventIDs))
		})
	})

	Describe("Step Completion", func() {
		It("should mark step as completed and move to next step", func() {
			ws.CompleteStep(workflow.StepCapture)
			
			Expect(ws.IsStepCompleted(workflow.StepCapture)).To(BeTrue())
			Expect(ws.GetCurrentStep()).To(Equal(workflow.StepMetadataReview))
		})

		It("should progress through all steps", func() {
			ws.CompleteStep(workflow.StepCapture)
			Expect(ws.GetCurrentStep()).To(Equal(workflow.StepMetadataReview))
			
			ws.CompleteStep(workflow.StepMetadataReview)
			Expect(ws.GetCurrentStep()).To(Equal(workflow.StepBurstSuggestion))
			
			ws.CompleteStep(workflow.StepBurstSuggestion)
			Expect(ws.GetCurrentStep()).To(Equal(workflow.StepFactExtraction))
			
			ws.CompleteStep(workflow.StepFactExtraction)
			Expect(ws.GetCurrentStep()).To(Equal(workflow.StepComplete))
		})
	})

	Describe("Step Skipping", func() {
		It("should allow skipping burst suggestions", func() {
			Expect(ws.CanSkipStep(workflow.StepBurstSuggestion)).To(BeTrue())
		})

		It("should allow skipping fact extraction", func() {
			Expect(ws.CanSkipStep(workflow.StepFactExtraction)).To(BeTrue())
		})

		It("should not allow skipping capture/import", func() {
			Expect(ws.CanSkipStep(workflow.StepCapture)).To(BeFalse())
		})

		It("should not allow skipping metadata review", func() {
			Expect(ws.CanSkipStep(workflow.StepMetadataReview)).To(BeFalse())
		})

		It("should skip burst suggestions and go to fact extraction", func() {
			ws.CompleteStep(workflow.StepCapture)
			ws.CompleteStep(workflow.StepMetadataReview)
			ws.SkipStep(workflow.StepBurstSuggestion)
			
			Expect(ws.IsStepSkipped(workflow.StepBurstSuggestion)).To(BeTrue())
			Expect(ws.GetCurrentStep()).To(Equal(workflow.StepFactExtraction))
		})

		It("should skip both optional steps and go to complete", func() {
			ws.CompleteStep(workflow.StepCapture)
			ws.CompleteStep(workflow.StepMetadataReview)
			ws.SkipStep(workflow.StepBurstSuggestion)
			ws.SkipStep(workflow.StepFactExtraction)
			
			Expect(ws.GetCurrentStep()).To(Equal(workflow.StepComplete))
		})
	})

	Describe("Progress Tracking", func() {
		It("should calculate progress percentage", func() {
			Expect(ws.GetProgress()).To(Equal(0))
			
			ws.CompleteStep(workflow.StepCapture)
			Expect(ws.GetProgress()).To(Equal(20)) // 1/5 steps
			
			ws.CompleteStep(workflow.StepMetadataReview)
			Expect(ws.GetProgress()).To(Equal(40)) // 2/5 steps
			
			ws.CompleteStep(workflow.StepBurstSuggestion)
			Expect(ws.GetProgress()).To(Equal(60)) // 3/5 steps
			
			ws.CompleteStep(workflow.StepFactExtraction)
			Expect(ws.GetProgress()).To(Equal(80)) // 4/5 steps
			
			ws.CompleteStep(workflow.StepComplete)
			Expect(ws.GetProgress()).To(Equal(100)) // 5/5 steps
		})

		It("should return progress steps", func() {
			ws.CompleteStep(workflow.StepCapture)
			ws.CompleteStep(workflow.StepMetadataReview)
			
			steps := ws.GetProgressSteps()
			Expect(steps).To(HaveLen(5))
			
			// First step completed
			Expect(steps[0].Completed).To(BeTrue())
			Expect(steps[0].Current).To(BeFalse())
			
			// Second step completed and not current (moved to burst)
			Expect(steps[1].Completed).To(BeTrue())
			Expect(steps[1].Current).To(BeFalse())
			
			// Third step current
			Expect(steps[2].Completed).To(BeFalse())
			Expect(steps[2].Current).To(BeTrue())
		})
	})

	Describe("Pending Items", func() {
		It("should track pending bursts", func() {
			ws.SetPendingBursts(5)
			Expect(ws.GetPendingBursts()).To(Equal(5))
			Expect(ws.HasPendingItems()).To(BeTrue())
		})

		It("should track pending facts", func() {
			ws.SetPendingFacts(10)
			Expect(ws.GetPendingFacts()).To(Equal(10))
			Expect(ws.HasPendingItems()).To(BeTrue())
		})

		It("should return pending summary", func() {
			ws.SetPendingBursts(3)
			ws.SetPendingFacts(7)
			
			summary := ws.GetPendingSummary()
			Expect(summary).To(Equal("3 burst suggestions, 7 fact extractions"))
		})

		It("should return no pending items when counts are zero", func() {
			Expect(ws.HasPendingItems()).To(BeFalse())
			Expect(ws.GetPendingSummary()).To(Equal("No pending items"))
		})
	})

	Describe("Review Later", func() {
		It("should mark burst suggestions for later review", func() {
			ws.CompleteStep(workflow.StepCapture)
			ws.CompleteStep(workflow.StepMetadataReview)
			ws.SetPendingBursts(5)
			ws.ReviewLater(workflow.StepBurstSuggestion)
			
			Expect(ws.IsStepSkipped(workflow.StepBurstSuggestion)).To(BeTrue())
			Expect(ws.GetPendingBursts()).To(Equal(5))
			Expect(ws.GetCurrentStep()).To(Equal(workflow.StepFactExtraction))
		})

		It("should mark fact extraction for later review", func() {
			ws.CompleteStep(workflow.StepCapture)
			ws.CompleteStep(workflow.StepMetadataReview)
			ws.CompleteStep(workflow.StepBurstSuggestion)
			ws.SetPendingFacts(10)
			ws.ReviewLater(workflow.StepFactExtraction)
			
			Expect(ws.IsStepSkipped(workflow.StepFactExtraction)).To(BeTrue())
			Expect(ws.GetPendingFacts()).To(Equal(10))
			Expect(ws.GetCurrentStep()).To(Equal(workflow.StepComplete))
		})
	})

	Describe("Step Descriptions", func() {
		It("should return description for capture step", func() {
			desc := ws.GetStepDescription(workflow.StepCapture)
			Expect(desc).To(ContainSubstring("Capture career events"))
		})

		It("should return description for metadata review", func() {
			desc := ws.GetStepDescription(workflow.StepMetadataReview)
			Expect(desc).To(ContainSubstring("Review and enrich"))
		})

		It("should return description for burst suggestion", func() {
			desc := ws.GetStepDescription(workflow.StepBurstSuggestion)
			Expect(desc).To(ContainSubstring("burst groupings"))
		})

		It("should return description for fact extraction", func() {
			desc := ws.GetStepDescription(workflow.StepFactExtraction)
			Expect(desc).To(ContainSubstring("extracted facts"))
		})

		It("should return description for complete", func() {
			desc := ws.GetStepDescription(workflow.StepComplete)
			Expect(desc).To(ContainSubstring("complete"))
		})
	})

	Describe("IsComplete", func() {
		It("should return false when workflow is not complete", func() {
			Expect(ws.IsComplete()).To(BeFalse())
		})

		It("should return true when workflow is complete", func() {
			ws.SetCurrentStep(workflow.StepComplete)
			Expect(ws.IsComplete()).To(BeTrue())
		})
	})

	Describe("Reset", func() {
		It("should reset workflow state", func() {
			ws.CompleteStep(workflow.StepCapture)
			ws.CompleteStep(workflow.StepMetadataReview)
			ws.SetPendingBursts(5)
			ws.SetPendingFacts(10)
			ws.SetEventIDs([]string{"event1", "event2"})
			
			ws.Reset()
			
			Expect(ws.GetCurrentStep()).To(Equal(workflow.StepCapture))
			Expect(ws.IsStepCompleted(workflow.StepCapture)).To(BeFalse())
			Expect(ws.IsStepCompleted(workflow.StepMetadataReview)).To(BeFalse())
			Expect(ws.GetPendingBursts()).To(Equal(0))
			Expect(ws.GetPendingFacts()).To(Equal(0))
			Expect(ws.GetEventIDs()).To(HaveLen(0))
		})
	})

	Describe("Event IDs", func() {
		It("should track event IDs", func() {
			eventIDs := []string{"event1", "event2", "event3"}
			ws.SetEventIDs(eventIDs)
			
			Expect(ws.GetEventIDs()).To(Equal(eventIDs))
		})
	})
})

