package app_test

import (
	"context"

	"github.com/baphled/kariya/internal/cli/app"
	"github.com/baphled/kariya/internal/cli/service"
	"github.com/baphled/kariya/internal/cli/workflow"
	"github.com/baphled/kariya/internal/repository/career"
	careerservice "github.com/baphled/kariya/internal/service/career"
	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
)

var _ = Describe("Workflow Integration", func() {
	var (
		repo   *career.MemoryRepository
		svc    *careerservice.Service
		cliSvc *service.CLIEventService
	)

	BeforeEach(func() {
		repo = career.NewMemoryRepository()
		svc = careerservice.NewService(repo)
		cliSvc = service.NewCLIEventService(svc)

		// Initialize model and context for future tests
		_ = app.NewModel(cliSvc, svc)
		_ = context.Background()
	})

	Describe("Workflow State Management", func() {
		It("should initialize workflow state for new capture", func() {
			ws := workflow.NewWorkflowState()
			Expect(ws).NotTo(BeNil())
			Expect(ws.GetCurrentStep()).To(Equal(workflow.StepCapture))
		})

		It("should track workflow progress through capture → metadata review", func() {
			ws := workflow.NewWorkflowStateFromCapture([]string{"event1"})
			Expect(ws.IsStepCompleted(workflow.StepCapture)).To(BeTrue())
			
			ws.CompleteStep(workflow.StepCapture)
			Expect(ws.GetCurrentStep()).To(Equal(workflow.StepMetadataReview))
		})

		It("should allow skipping burst suggestions", func() {
			ws := workflow.NewWorkflowState()
			ws.CompleteStep(workflow.StepCapture)
			ws.CompleteStep(workflow.StepMetadataReview)
			ws.SkipStep(workflow.StepBurstSuggestion)
			
			Expect(ws.IsStepSkipped(workflow.StepBurstSuggestion)).To(BeTrue())
			Expect(ws.GetCurrentStep()).To(Equal(workflow.StepFactExtraction))
		})

		It("should allow skipping fact extraction", func() {
			ws := workflow.NewWorkflowState()
			ws.CompleteStep(workflow.StepCapture)
			ws.CompleteStep(workflow.StepMetadataReview)
			ws.CompleteStep(workflow.StepBurstSuggestion)
			ws.SkipStep(workflow.StepFactExtraction)
			
			Expect(ws.IsStepSkipped(workflow.StepFactExtraction)).To(BeTrue())
			Expect(ws.GetCurrentStep()).To(Equal(workflow.StepComplete))
		})

		It("should track pending bursts", func() {
			ws := workflow.NewWorkflowState()
			ws.SetPendingBursts(3)
			
			Expect(ws.GetPendingBursts()).To(Equal(3))
			Expect(ws.HasPendingItems()).To(BeTrue())
		})

		It("should track pending facts", func() {
			ws := workflow.NewWorkflowState()
			ws.SetPendingFacts(5)
			
			Expect(ws.GetPendingFacts()).To(Equal(5))
			Expect(ws.HasPendingItems()).To(BeTrue())
		})

		It("should calculate progress percentage correctly", func() {
			ws := workflow.NewWorkflowState()
			Expect(ws.GetProgress()).To(Equal(0))
			
			ws.CompleteStep(workflow.StepCapture)
			Expect(ws.GetProgress()).To(Equal(20))
			
			ws.CompleteStep(workflow.StepMetadataReview)
			Expect(ws.GetProgress()).To(Equal(40))
			
			ws.CompleteStep(workflow.StepBurstSuggestion)
			Expect(ws.GetProgress()).To(Equal(60))
			
			ws.CompleteStep(workflow.StepFactExtraction)
			Expect(ws.GetProgress()).To(Equal(80))
			
			ws.CompleteStep(workflow.StepComplete)
			Expect(ws.GetProgress()).To(Equal(100))
		})
	})

	Describe("End-to-End Workflow", func() {
		It("should complete full workflow: capture → metadata → burst → fact → complete", func() {
			ws := workflow.NewWorkflowStateFromCapture([]string{"event1", "event2"})
			
			// Step 1: Capture complete (already done in constructor)
			Expect(ws.IsStepCompleted(workflow.StepCapture)).To(BeTrue())
			
			// Step 2: Metadata review
			ws.CompleteStep(workflow.StepCapture)
			Expect(ws.GetCurrentStep()).To(Equal(workflow.StepMetadataReview))
			ws.CompleteStep(workflow.StepMetadataReview)
			
			// Step 3: Burst suggestions
			Expect(ws.GetCurrentStep()).To(Equal(workflow.StepBurstSuggestion))
			ws.SetPendingBursts(2)
			ws.CompleteStep(workflow.StepBurstSuggestion)
			
			// Step 4: Fact extraction
			Expect(ws.GetCurrentStep()).To(Equal(workflow.StepFactExtraction))
			ws.SetPendingFacts(5)
			ws.CompleteStep(workflow.StepFactExtraction)
			
			// Step 5: Complete
			Expect(ws.GetCurrentStep()).To(Equal(workflow.StepComplete))
			Expect(ws.IsComplete()).To(BeTrue())
		})

		It("should support review later workflow", func() {
			ws := workflow.NewWorkflowStateFromCapture([]string{"event1", "event2"})
			ws.CompleteStep(workflow.StepCapture)
			ws.CompleteStep(workflow.StepMetadataReview)
			
			// Mark burst suggestions for review later
			ws.SetPendingBursts(3)
			ws.ReviewLater(workflow.StepBurstSuggestion)
			
			Expect(ws.IsStepSkipped(workflow.StepBurstSuggestion)).To(BeTrue())
			Expect(ws.GetPendingBursts()).To(Equal(3))
			Expect(ws.HasPendingItems()).To(BeTrue())
			Expect(ws.GetPendingSummary()).To(ContainSubstring("3 burst suggestions"))
		})

		It("should support import workflow", func() {
			importedIDs := []string{"imp1", "imp2", "imp3"}
			ws := workflow.NewWorkflowStateFromImport(importedIDs)
			
			Expect(ws.IsStepCompleted(workflow.StepImport)).To(BeTrue())
			Expect(ws.GetEventIDs()).To(Equal(importedIDs))
			Expect(ws.GetCurrentStep()).To(Equal(workflow.StepImport))
		})
	})

	Describe("Workflow Progress Display", func() {
		It("should generate progress steps for display", func() {
			ws := workflow.NewWorkflowState()
			ws.CompleteStep(workflow.StepCapture)
			
			steps := ws.GetProgressSteps()
			Expect(steps).To(HaveLen(5))
			
			// First step completed
			Expect(steps[0].Label).To(ContainSubstring("Capture"))
			Expect(steps[0].Completed).To(BeTrue())
			
			// Second step current
			Expect(steps[1].Label).To(ContainSubstring("Metadata"))
			Expect(steps[1].Current).To(BeTrue())
		})

		It("should show skipped steps in progress", func() {
			ws := workflow.NewWorkflowState()
			ws.CompleteStep(workflow.StepCapture)
			ws.CompleteStep(workflow.StepMetadataReview)
			ws.SkipStep(workflow.StepBurstSuggestion)
			
			steps := ws.GetProgressSteps()
			
			// Burst suggestion step should be marked as skipped
			burstStep := steps[2]
			Expect(burstStep.Skipped).To(BeTrue())
			Expect(burstStep.Current).To(BeFalse())
			
			// Current step should be fact extraction
			factStep := steps[3]
			Expect(factStep.Current).To(BeTrue())
		})
	})

	Describe("Step Descriptions", func() {
		It("should provide helpful descriptions for each step", func() {
			ws := workflow.NewWorkflowState()
			
			Expect(ws.GetStepDescription(workflow.StepCapture)).To(ContainSubstring("Capture career events"))
			Expect(ws.GetStepDescription(workflow.StepMetadataReview)).To(ContainSubstring("Review and enrich"))
			Expect(ws.GetStepDescription(workflow.StepBurstSuggestion)).To(ContainSubstring("burst groupings"))
			Expect(ws.GetStepDescription(workflow.StepFactExtraction)).To(ContainSubstring("extracted facts"))
			Expect(ws.GetStepDescription(workflow.StepComplete)).To(ContainSubstring("complete"))
		})
	})
})

