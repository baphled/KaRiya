package app

import (
	"github.com/baphled/kariya/internal/domain/career"
	"github.com/baphled/kariya/internal/service/career/burst_fact"
)

// FormSubmittedMsg is sent when a form is successfully submitted
type FormSubmittedMsg struct {
	Event *career.CareerEvent
	Err   error
}

// NavigateMsg is sent to navigate to a specific screen
type NavigateMsg struct {
	Screen Screen
}

// SuccessNavigateMsg is sent from success screen for navigation
type SuccessNavigateMsg struct {
	Action string // "capture_another", "view_list", "exit"
}

// EditEventMsg is sent to edit a specific event
type EditEventMsg struct {
	Event *career.CareerEvent
}

// DeleteEventMsg is sent to delete a specific event
type DeleteEventMsg struct {
	EventID string
}

// ConfirmDeleteMsg is sent when delete is confirmed
type ConfirmDeleteMsg struct {
	EventID string
}

// CancelDeleteMsg is sent when delete is cancelled
type CancelDeleteMsg struct{}

// EventDeletedMsg is sent when an event is successfully deleted
type EventDeletedMsg struct {
	EventID string
	Err     error
}

// EventUpdatedMsg is sent when an event is successfully updated
type EventUpdatedMsg struct {
	Event *career.CareerEvent
	Err   error
}

// BackMsg is sent to go back to the previous screen
type BackMsg struct{}

// QuitMsg is sent to quit the application
type QuitMsg struct{}

// BulkOperationsMsg is sent to open bulk operations for selected events
type BulkOperationsMsg struct {
	Events []*career.CareerEvent
}

// ApplyBulkOperationsMsg is sent when bulk operations are applied
type ApplyBulkOperationsMsg struct {
	EventIDs   []string
	Company    string
	Project    string
	Tags       []string
	Categories []string
	Err        error
}

// CancelBulkOperationsMsg is sent when bulk operations are cancelled
type CancelBulkOperationsMsg struct{}

// MetadataReviewTriggeredMsg is sent when metadata review should be shown after import
type MetadataReviewTriggeredMsg struct {
	Events []*career.CareerEvent
}

// BreadcrumbClickedMsg is sent when a breadcrumb is clicked
type BreadcrumbClickedMsg struct {
	Index int // Index of the clicked breadcrumb
}

// BurstSuggestionsTriggeredMsg is sent when burst suggestions should be reviewed
type BurstSuggestionsTriggeredMsg struct {
	EventIDs []string
}

// BurstSuggestionsReadyMsg is sent when burst suggestions have been generated
type BurstSuggestionsReadyMsg struct {
	Suggestions []burst_fact.BurstSuggestion
	Err         error
}

// ConfirmBurstMsg is sent when user confirms a burst suggestion
type ConfirmBurstMsg struct {
	Burst *career.Burst
}

// RejectBurstSuggestionMsg is sent when user rejects a burst suggestion
type RejectBurstSuggestionMsg struct {
	EventIDs []string
}

// BurstProcessingCompleteMsg is sent when burst suggestion workflow is done
type BurstProcessingCompleteMsg struct {
	ConfirmedBursts []career.Burst
	RejectedGroups  [][]string
}

// SkipStepMsg is sent when user wants to skip a workflow step
type SkipStepMsg struct {
	Step string // "burst_suggestion" or "fact_extraction"
}

// ReviewLaterMsg is sent when user wants to review items later
type ReviewLaterMsg struct {
	Step string // "burst_suggestion" or "fact_extraction"
}

// ViewPendingItemsMsg is sent to view pending bursts/facts from home screen
type ViewPendingItemsMsg struct {
	ItemType string // "bursts" or "facts"
}

// FactExtractionTriggeredMsg is sent when fact extraction should be performed
type FactExtractionTriggeredMsg struct {
	EventIDs []string
	BurstIDs []string
}

// FactsReadyMsg is sent when facts have been extracted
type FactsReadyMsg struct {
	Facts []*career.Fact
	Err   error
}

// ConfirmFactMsg is sent when user confirms a fact
type ConfirmFactMsg struct {
	Fact *career.Fact
}

// RejectFactMsg is sent when user rejects a fact
type RejectFactMsg struct {
	FactID string
}

// FactProcessingCompleteMsg is sent when fact extraction workflow is done
type FactProcessingCompleteMsg struct {
	ConfirmedFacts []*career.Fact
	RejectedFacts  []string
}
