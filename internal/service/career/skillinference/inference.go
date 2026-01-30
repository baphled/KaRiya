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

// SkillInferenceService detects skills from event text and creates skill records.
//
// This service analyzes event text for technology mentions using a keyword
// dictionary (technology.GetTechnologyKeywords()) and word boundary regex
// matching to avoid false positives.
//
// Usage:
//
//	service := NewSkillInferenceService(skillRepo, eventRepo)
//	suggestions, err := service.InferSkillsFromEvents(ctx, events)
//
//	// Review suggestions in UI, then create skills
//	skills, err := service.CreateSkillsFromSuggestions(ctx, acceptedSuggestions)
type SkillInferenceService interface {
	// InferSkillsFromEvents analyzes all events for technology mentions.
	// Returns skill suggestions sorted by confidence (highest first).
	// Empty slice returned if no skills detected (not an error).
	//
	// Algorithm:
	// 1. Scan event text for keyword matches (word boundary regex)
	// 2. Extract context snippets (~80 chars around match)
	// 3. Calculate confidence based on usage patterns
	// 4. Deduplicate same skill across events
	// 5. Sort by confidence descending
	InferSkillsFromEvents(
		ctx context.Context,
		events []*career.Event,
	) ([]SkillSuggestion, error)

	// InferSkillsFromBurst analyzes events within a specific burst.
	// Useful for targeted skill detection after burst confirmation.
	// Returns skill suggestions specific to this burst's events.
	InferSkillsFromBurst(
		ctx context.Context,
		burst *career.Burst,
		events []*career.Event,
	) ([]SkillSuggestion, error)

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
