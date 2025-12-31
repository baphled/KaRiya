package workflow

import (
	"fmt"
)

// Step represents a step in the workflow
type Step string

const (
	StepCapture         Step = "capture"
	StepImport          Step = "import"
	StepMetadataReview  Step = "metadata_review"
	StepBurstSuggestion Step = "burst_suggestion"
	StepFactExtraction  Step = "fact_extraction"
	StepComplete        Step = "complete"
)

// WorkflowState represents the current state of the enrichment workflow
type WorkflowState struct {
	currentStep    Step
	completedSteps map[Step]bool
	skippedSteps   map[Step]bool
	pendingBursts  int // Number of pending burst suggestions
	pendingFacts   int // Number of pending fact extractions
	eventIDs       []string
}

// NewWorkflowState creates a new workflow state
func NewWorkflowState() *WorkflowState {
	return &WorkflowState{
		currentStep:    StepCapture,
		completedSteps: make(map[Step]bool),
		skippedSteps:   make(map[Step]bool),
		pendingBursts:  0,
		pendingFacts:   0,
		eventIDs:       []string{},
	}
}

// NewWorkflowStateFromCapture creates workflow state starting from capture
func NewWorkflowStateFromCapture(eventIDs []string) *WorkflowState {
	ws := NewWorkflowState()
	ws.currentStep = StepCapture
	ws.eventIDs = eventIDs
	ws.completedSteps[StepCapture] = true
	return ws
}

// NewWorkflowStateFromImport creates workflow state starting from import
func NewWorkflowStateFromImport(eventIDs []string) *WorkflowState {
	ws := NewWorkflowState()
	ws.currentStep = StepImport
	ws.eventIDs = eventIDs
	ws.completedSteps[StepImport] = true
	return ws
}

// SetCurrentStep sets the current step in the workflow
func (w *WorkflowState) SetCurrentStep(step Step) {
	w.currentStep = step
}

// GetCurrentStep returns the current step
func (w *WorkflowState) GetCurrentStep() Step {
	return w.currentStep
}

// CompleteStep marks a step as completed
func (w *WorkflowState) CompleteStep(step Step) {
	w.completedSteps[step] = true
	// Move to next step
	w.currentStep = w.getNextStep(step)
}

// SkipStep marks a step as skipped
func (w *WorkflowState) SkipStep(step Step) {
	w.skippedSteps[step] = true
	// Move to next step
	w.currentStep = w.getNextStep(step)
}

// IsStepCompleted checks if a step is completed
func (w *WorkflowState) IsStepCompleted(step Step) bool {
	return w.completedSteps[step]
}

// IsStepSkipped checks if a step is skipped
func (w *WorkflowState) IsStepSkipped(step Step) bool {
	return w.skippedSteps[step]
}

// GetProgress returns the completion percentage (0-100)
func (w *WorkflowState) GetProgress() int {
	totalSteps := 5 // capture/import, metadata, burst, fact, complete
	completed := len(w.completedSteps)
	return (completed * 100) / totalSteps
}

// GetProgressSteps returns steps for progress indicator
func (w *WorkflowState) GetProgressSteps() []StepInfo {
	return []StepInfo{
		{
			Step:      StepCapture,
			Label:     "Capture/Import",
			Completed: w.IsStepCompleted(StepCapture) || w.IsStepCompleted(StepImport),
			Skipped:   false, // Cannot skip
			Current:   w.currentStep == StepCapture || w.currentStep == StepImport,
		},
		{
			Step:      StepMetadataReview,
			Label:     "Metadata Review",
			Completed: w.IsStepCompleted(StepMetadataReview),
			Skipped:   w.IsStepSkipped(StepMetadataReview),
			Current:   w.currentStep == StepMetadataReview,
		},
		{
			Step:      StepBurstSuggestion,
			Label:     "Burst Suggestions",
			Completed: w.IsStepCompleted(StepBurstSuggestion),
			Skipped:   w.IsStepSkipped(StepBurstSuggestion),
			Current:   w.currentStep == StepBurstSuggestion,
		},
		{
			Step:      StepFactExtraction,
			Label:     "Fact Extraction",
			Completed: w.IsStepCompleted(StepFactExtraction),
			Skipped:   w.IsStepSkipped(StepFactExtraction),
			Current:   w.currentStep == StepFactExtraction,
		},
		{
			Step:      StepComplete,
			Label:     "Complete",
			Completed: w.IsStepCompleted(StepComplete),
			Skipped:   false, // Cannot skip
			Current:   w.currentStep == StepComplete,
		},
	}
}

// StepInfo represents information about a workflow step
type StepInfo struct {
	Step      Step
	Label     string
	Completed bool
	Skipped   bool
	Current   bool
}

// getNextStep determines the next step in the workflow
func (w *WorkflowState) getNextStep(currentStep Step) Step {
	switch currentStep {
	case StepCapture, StepImport:
		return StepMetadataReview
	case StepMetadataReview:
		// If burst detection is enabled, go to burst suggestion
		if !w.skippedSteps[StepBurstSuggestion] {
			return StepBurstSuggestion
		}
		// Otherwise, go to fact extraction if enabled
		if !w.skippedSteps[StepFactExtraction] {
			return StepFactExtraction
		}
		return StepComplete
	case StepBurstSuggestion:
		// If fact extraction is enabled, go there
		if !w.skippedSteps[StepFactExtraction] {
			return StepFactExtraction
		}
		return StepComplete
	case StepFactExtraction:
		return StepComplete
	case StepComplete:
		return StepComplete
	default:
		return StepComplete
	}
}

// CanSkipStep checks if a step can be skipped
func (w *WorkflowState) CanSkipStep(step Step) bool {
	// Only burst suggestions and fact extraction can be skipped
	return step == StepBurstSuggestion || step == StepFactExtraction
}

// SetPendingBursts sets the number of pending burst suggestions
func (w *WorkflowState) SetPendingBursts(count int) {
	w.pendingBursts = count
}

// GetPendingBursts returns the number of pending burst suggestions
func (w *WorkflowState) GetPendingBursts() int {
	return w.pendingBursts
}

// SetPendingFacts sets the number of pending fact extractions
func (w *WorkflowState) SetPendingFacts(count int) {
	w.pendingFacts = count
}

// GetPendingFacts returns the number of pending fact extractions
func (w *WorkflowState) GetPendingFacts() int {
	return w.pendingFacts
}

// GetEventIDs returns the event IDs being processed
func (w *WorkflowState) GetEventIDs() []string {
	return w.eventIDs
}

// SetEventIDs sets the event IDs being processed
func (w *WorkflowState) SetEventIDs(eventIDs []string) {
	w.eventIDs = eventIDs
}

// IsComplete checks if the workflow is complete
func (w *WorkflowState) IsComplete() bool {
	return w.currentStep == StepComplete
}

// GetStepDescription returns a human-readable description of the step
func (w *WorkflowState) GetStepDescription(step Step) string {
	switch step {
	case StepCapture:
		return "Capture career events using timeline journaling, CV backfill, or manual entry"
	case StepImport:
		return "Import career events from CSV file"
	case StepMetadataReview:
		return "Review and enrich event metadata (company, project, tags, categories)"
	case StepBurstSuggestion:
		return "Review suggested burst groupings of related events"
	case StepFactExtraction:
		return "Review and confirm extracted facts from events and bursts"
	case StepComplete:
		return "Workflow complete! Your career events are fully enriched."
	default:
		return "Unknown step"
	}
}

// ReviewLater marks items for later review
func (w *WorkflowState) ReviewLater(step Step) {
	// Mark as skipped for now, but track that items are pending
	w.skippedSteps[step] = true

	// Update pending counts based on step
	switch step {
	case StepBurstSuggestion:
		// Pending bursts will be set externally
	case StepFactExtraction:
		// Pending facts will be set externally
	}

	// Move to next step
	w.currentStep = w.getNextStep(step)
}

// HasPendingItems checks if there are any pending items to review
func (w *WorkflowState) HasPendingItems() bool {
	return w.pendingBursts > 0 || w.pendingFacts > 0
}

// GetPendingSummary returns a summary of pending items
func (w *WorkflowState) GetPendingSummary() string {
	if !w.HasPendingItems() {
		return "No pending items"
	}

	summary := ""
	if w.pendingBursts > 0 {
		summary += fmt.Sprintf("%d burst suggestions", w.pendingBursts)
	}
	if w.pendingFacts > 0 {
		if summary != "" {
			summary += ", "
		}
		summary += fmt.Sprintf("%d fact extractions", w.pendingFacts)
	}

	return summary
}

// Reset resets the workflow state
func (w *WorkflowState) Reset() {
	w.currentStep = StepCapture
	w.completedSteps = make(map[Step]bool)
	w.skippedSteps = make(map[Step]bool)
	w.pendingBursts = 0
	w.pendingFacts = 0
	w.eventIDs = []string{}
}
