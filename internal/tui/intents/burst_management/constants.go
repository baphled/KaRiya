// Package burst_management implements the BurstManagement intent for managing career bursts.
package burst_management

// State represents the current state of the BurstManagement intent.
type State string

// State constants for the BurstManagement intent.
const (
	// StateList shows the burst list view.
	StateList State = "list"

	// StateDetail shows the burst detail view.
	StateDetail State = "detail"

	// StateDetailEvents shows events in the burst.
	StateDetailEvents State = "detail_events"

	// StateDetailFacts shows facts extracted from the burst.
	StateDetailFacts State = "detail_facts"

	// StateEdit shows the burst edit form.
	StateEdit State = "edit"

	// StateDeleteConfirm shows the delete confirmation modal.
	StateDeleteConfirm State = "delete_confirm"

	// StateConfirm shows the burst confirmation view.
	StateConfirm State = "confirm"

	// StateExtractingFacts shows the fact extraction progress.
	StateExtractingFacts State = "extracting_facts"

	// StateSuggesting shows loading state during burst detection.
	StateSuggesting State = "suggesting"

	// StateSuggestionReview shows burst suggestion review.
	StateSuggestionReview State = "suggestion_review"

	// StateInferringSkills shows the skill inference progress.
	StateInferringSkills State = "inferring_skills"

	// StateSkillSuggestionReview shows skill suggestion review.
	StateSkillSuggestionReview State = "skill_suggestion_review"
)
