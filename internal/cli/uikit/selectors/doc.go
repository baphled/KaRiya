// Package selectors provides selection components.
//
// # Overview
//
// The selectors package implements interactive selection components
// for choosing items from lists or categories. These components
// provide consistent selection UX across the application.
//
// # Available Components
//
//   - CategorySelector: Hierarchical category selection
//   - TagSelector: Multi-select tags with search
//   - SkillSelector: Skill selection with levels
//   - AudienceRelevanceSelector: Target audience selection
//
// # Usage
//
//	selector := selectors.NewCategorySelector(theme, categories)
//	selected, err := selector.Run()
package selectors
