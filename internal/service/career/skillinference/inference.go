// Package skillinference provides skill inference service for automatic
// technology detection from event text.
package skillinference

import (
	"context"

	"github.com/baphled/kariya/internal/domain/career"
)

// SkillSuggestion represents a detected skill with metadata.
// Similar to BurstSuggestion but tailored for skill inference.
type SkillSuggestion struct {
	Name       string   // Canonical skill name (e.g., "Go", "PostgreSQL")
	Category   string   // Skill category: backend, frontend, database, devops, cloud, mobile, tooling
	Confidence float64  // 0.0-1.0 confidence score based on usage patterns
	EventIDs   []string // Events where this skill was detected
	Contexts   []string // Text snippets showing usage (max 3 for UI display)
}

// InferenceResult contains the results of skill inference including
// new suggestions and names of skills that already exist in the repository.
type InferenceResult struct {
	Suggestions        []SkillSuggestion
	ExistingSkillNames []string
}

// SkillInferenceService detects skills from event text and creates skill records.
//
// This service analyzes event text for technology mentions using a keyword
// dictionary (technology.GetTechnologyKeywords()) and word boundary regex
// matching to avoid false positives.
//
// Usage:
//
//	service := NewSkillInferenceService(skillRepo, eventRepo)
//	result, err := service.InferSkillsFromEvents(ctx, events)
//
//	// Review result.Suggestions in UI, show result.ExistingSkillNames as info
//	skills, err := service.CreateSkillsFromSuggestions(ctx, acceptedSuggestions)
type SkillInferenceService interface { //nolint:revive // renaming would break all consumers of this interface
	// InferSkillsFromEvents analyzes all events for technology mentions.
	// Returns an InferenceResult containing new skill suggestions and
	// names of skills that already exist in the repository.
	//
	// Algorithm:
	// 1. Scan event text for keyword matches (word boundary regex)
	// 2. Extract context snippets (~80 chars around match)
	// 3. Calculate confidence based on usage patterns
	// 4. Deduplicate same skill across events
	// 5. Filter out skills that already exist (reported in ExistingSkillNames)
	InferSkillsFromEvents(
		ctx context.Context,
		events []*career.Event,
	) (*InferenceResult, error)

	// InferSkillsFromBurst analyzes events within a specific burst.
	// Useful for targeted skill detection after burst confirmation.
	// Returns an InferenceResult specific to this burst's events.
	InferSkillsFromBurst(
		ctx context.Context,
		burst *career.Burst,
		events []*career.Event,
	) (*InferenceResult, error)

	// CreateSkillsFromSuggestions persists accepted suggestions as skills
	// and links them to source events via junction table.
	//
	// Behavior:
	// - Reuses existing skill if name matches (case-insensitive)
	// - Creates new skill if not exists
	// - Links skill to events via event.Skills field
	// - Updates skill.LastUsed to most recent event date
	//
	// Returns list of created/reused skills.
	CreateSkillsFromSuggestions(
		ctx context.Context,
		suggestions []SkillSuggestion,
	) ([]*career.Skill, error)
}
