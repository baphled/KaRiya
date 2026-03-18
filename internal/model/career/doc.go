// Package career provides GORM database models for career domain entities.
//
// Each model struct maps directly to a SQLite table and provides ToDomain and
// FromDomain functions for translating between the persistence layer and the
// application's domain types in the career package.
//
// # Models
//
// The package contains the following models:
//   - Skill: Technical skills with proficiency levels and usage tracking
//   - Event: Career events capturing professional accomplishments
//   - Fact: Structured insights derived from events or bursts
//   - Burst: Groups of related events occurring in concentrated periods
//
// # Type Conversions
//
// All models provide bidirectional conversion to domain types:
//   - ToDomain converts a database model to its domain equivalent
//   - FromDomain (e.g. SkillFromDomain) creates a model from a domain type
package career
