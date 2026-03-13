// Package fact provides view components for fact management.
package fact

import (
	"fmt"

	"github.com/baphled/kariya/internal/ui/display"
)

const textFieldFmt = "Text: %s\n\n"

// RenderFactDetail renders the detail view for a single fact.
//
// Expected:
//   - fact must be valid.
//
// Returns:
//   - A string value.
//
// Side effects:
//   - None.
func RenderFactDetail(fact display.Fact) string {
	if isEmptyFact(fact) {
		return "No fact selected"
	}

	var content string
	content += "Fact Details\n\n"
	content += fmt.Sprintf(textFieldFmt, fact.Text)
	content += fmt.Sprintf("Categories: %v\n", fact.CompetencyCategories)
	content += fmt.Sprintf("Strength Signal: %s\n", fact.StrengthSignal)
	content += fmt.Sprintf("Role Fit: %v\n", fact.RoleFit)
	content += fmt.Sprintf("Audience Relevance: %v\n", fact.AudienceRelevance)
	content += fmt.Sprintf("\nCreated: %s\n", fact.CreatedAt.Format("2006-01-02 15:04"))

	return content
}

// RenderDeleteConfirm renders the delete confirmation view.
//
// Expected:
//   - fact must be valid.
//
// Returns:
//   - A string value.
//
// Side effects:
//   - None.
func RenderDeleteConfirm(fact display.Fact) string {
	if isEmptyFact(fact) {
		return "No fact to delete"
	}

	var content string
	content += "Confirm Deletion\n\n"
	content += "Delete this fact?\n\n"
	content += fmt.Sprintf("Text: %s\n\n", truncateText(fact.Text, 100))
	content += "Warning: This action cannot be undone.\n"

	return content
}

// RenderResults renders the fact statistics view.
//
// Expected:
//   - int must be valid.
//   - int must be valid.
//
// Returns:
//   - A string value.
//
// Side effects:
//   - None.
func RenderResults(totalFacts, currentCount int) string {
	var content string
	content += "Fact Statistics\n\n"
	content += fmt.Sprintf("Total facts in database: %d\n", totalFacts)

	if currentCount > 0 {
		content += fmt.Sprintf("Currently viewing: %d facts\n", currentCount)
	}

	return content
}

// RenderEditorFallback renders the editor content when no modal is available.
//
// Expected:
//   - fact must be valid.
//   - bool must be valid.
//
// Returns:
//   - A string value.
//
// Side effects:
//   - None.
func RenderEditorFallback(fact display.Fact, hasFormErrors bool) string {
	var content string
	content += "Edit Fact\n\n"

	if !isEmptyFact(fact) {
		content += fmt.Sprintf(textFieldFmt, fact.Text)
		content += fmt.Sprintf("Categories: %v\n", fact.CompetencyCategories)
		content += fmt.Sprintf("Strength Signal: %s\n", fact.StrengthSignal)
		content += fmt.Sprintf("Role Fit: %v\n\n", fact.RoleFit)

		if !hasFormErrors {
			content += "Make your changes and press Ctrl+S to save.\n"
		}
	} else {
		content += "No fact loaded for editing.\n"
	}

	return content
}

func truncateText(s string, maxLen int) string {
	if len(s) <= maxLen {
		return s
	}
	return s[:maxLen] + "..."
}

func isEmptyFact(fact display.Fact) bool {
	return fact.ID == "" &&
		fact.Text == "" &&
		len(fact.CompetencyCategories) == 0 &&
		fact.RoleFit == "" &&
		len(fact.AudienceRelevance) == 0 &&
		fact.StrengthSignal == "" &&
		fact.SourceEventID == "" &&
		fact.SourceBurstID == "" &&
		fact.CreatedAt.IsZero() &&
		fact.UpdatedAt.IsZero()
}
