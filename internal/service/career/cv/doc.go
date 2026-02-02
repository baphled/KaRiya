// Package cv provides CV generation and formatting services for producing
// professional resumes from career data.
//
// # Overview
//
// The cv package handles:
// - CV generation from career events, facts, and skills
// - Bullet point creation and enhancement
// - CV export in multiple formats (text, markdown, YAML)
// - Section building and formatting
// - Profile inference and narrative generation
//
// # Usage Example
//
//	service := cv.NewCVGenerationService()
//	cv, err := service.GenerateCV(ctx, config)
package cv
