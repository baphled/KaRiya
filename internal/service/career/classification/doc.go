// Package classification provides content classification services for categorizing
// career data into competency areas and skill categories.
//
// # Overview
//
// The classification package offers:
// - Competency categorization for career facts and events
// - Skill category assignment using keyword matching
// - Multi-label classification for complex content
//
// # Usage Example
//
//	classifier := classification.NewClassifier(keywords)
//	category := classifier.Classify(text)
package classification
