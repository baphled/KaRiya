// Package career provides core domain services for managing career data including
// events, bursts, facts, and skills. It serves as the orchestration layer between
// the CLI intents and the repository layer.
//
// # Overview
//
// The career package implements the domain logic for:
// - Capturing and managing career events (jobs, projects, achievements)
// - Detecting skill bursts from career event patterns
// - Extracting and validating career facts from events
// - Managing skill inventory and categorization
//
// # Usage Example
//
//	service := career.NewService()
//	event, err := service.CaptureEvent(ctx, careerEvent)
package career
