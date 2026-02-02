// Package skillinference provides automatic skill detection and inference from
// career events using keyword matching and pattern recognition.
//
// # Overview
//
// The skillinference package analyzes career events and bursts to:
// - Identify technologies and skills mentioned in event descriptions
// - Infer skills from achievement patterns and project contexts
// - Categorize detected skills using the technology keyword dictionary
// - Suggest skills for confirmation by the user
//
// # Usage Example
//
//	service := skillinference.NewSkillInferenceService(keywordDict)
//	suggestions, err := service.InferSkillsFromEvents(ctx, events)
package skillinference
