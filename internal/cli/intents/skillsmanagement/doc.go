// Package skillsmanagement implements the SkillsManagement intent for managing skills.
//
// # Overview
//
// The skillsmanagement package provides the intent for managing career skills
// including categorization, proficiency levels, and organization.
//
// # Architecture
//
// This package follows the standard intent subdirectory structure:
//   - context.go: IntentContext with services
//   - intent.go: Main intent implementation
//   - constants.go: State enum
//
// # Usage
//
// Create the intent:
//
//	ctx := &skillsmanagement.IntentContext{
//	    SkillService: skillService,
//	}
//	intent, err := skillsmanagement.NewIntent(ctx)
//
// For more details, see the intent documentation.
package skillsmanagement
