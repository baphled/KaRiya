// Package fact_management implements the FactManagement intent for managing career facts.
package fact_management

import (
	"fmt"

	"github.com/baphled/kariya/internal/cli/intents"
	"github.com/baphled/kariya/internal/cli/uikit/primitives"
)

// syncTableSelection syncs the TableBehavior selection with the context.
func (i *Intent) syncTableSelection() {
	i.context.SelectedFactIndex = i.tableBehavior.GetSelectedIndex()
	if selected := i.tableBehavior.GetSelectedItem(); selected != nil {
		i.context.SelectedFact = *selected
	} else {
		i.context.SelectedFact = nil
	}
}

// setCancelled marks the intent as cancelled.
func (i *Intent) setCancelled() {
	i.state = StateCompleted
	i.result = &intents.IntentResult[*Result]{
		Status: intents.Cancelled,
		Data: &Result{
			Action: "none",
			Facts:  i.context.Facts,
		},
	}
}

// View helper methods.

func (i *Intent) getBreadcrumbs() []string {
	breadcrumbs := []string{"Main Menu", "Manage Facts"}

	switch i.state {
	case StateView, StateDeleteConfirm:
		if i.context.SelectedFact != nil {
			factID := i.context.SelectedFact.ID
			if len(factID) > 8 {
				factID = factID[:8]
			}
			factName := fmt.Sprintf("Fact #%s", factID)
			breadcrumbs = append(breadcrumbs, factName)
		}
	case StateEditor:
		if i.context.IsNewFact {
			breadcrumbs = append(breadcrumbs, "New Fact")
		} else if i.context.EditingFact != nil {
			factID := i.context.EditingFact.ID
			if len(factID) > 8 {
				factID = factID[:8]
			}
			breadcrumbs = append(breadcrumbs, fmt.Sprintf("Edit Fact #%s", factID))
		}
	case StateResults:
		breadcrumbs = append(breadcrumbs, "Results")
	}

	return breadcrumbs
}

func (i *Intent) getStateContent() string {
	switch i.state {
	case StateList:
		// Table is self-contained, just render it.
		return i.tableBehavior.Render()
	case StateView:
		return i.getViewFactContent()
	case StateEditor:
		return i.getEditorContent()
	case StateDeleteConfirm:
		return i.getDeleteConfirmContent()
	case StateResults:
		return i.getResultsContent()
	case StateCompleted:
		return "Fact management completed"
	}
	return "Unknown state"
}

func (i *Intent) getViewFactContent() string {
	if i.context.SelectedFact == nil {
		return "No fact selected"
	}

	fact := i.context.SelectedFact
	var content string
	content += "Fact Details\n\n"
	content += fmt.Sprintf("Text: %s\n\n", fact.Text)
	content += fmt.Sprintf("Categories: %v\n", fact.CompetencyCategories)
	content += fmt.Sprintf("Strength Signal: %s\n", fact.StrengthSignal)
	content += fmt.Sprintf("Role Fit: %v\n", fact.RoleFit)
	content += fmt.Sprintf("Audience Relevance: %v\n", fact.AudienceRelevance)
	content += fmt.Sprintf("\nCreated: %s\n", fact.CreatedAt.Format("2006-01-02 15:04"))

	return content
}

func (i *Intent) getEditorContent() string {
	// If modal is available, render just the form content (not full modal container).
	// StandardView already provides the layout structure.
	if i.editModal != nil {
		return i.editModal.GetContent()
	}

	// Fallback for legacy behavior.
	var content string
	content += "Edit Fact\n\n"

	if i.context.EditingFact != nil {
		content += fmt.Sprintf("Text: %s\n\n", i.context.EditingFact.Text)
		content += fmt.Sprintf("Categories: %v\n", i.context.EditingFact.CompetencyCategories)
		content += fmt.Sprintf("Strength Signal: %s\n", i.context.EditingFact.StrengthSignal)
		content += fmt.Sprintf("Role Fit: %v\n\n", i.context.EditingFact.RoleFit)

		// Errors are shown in modal via getContextHelp.
		if !i.context.HasFormErrors() {
			content += "Make your changes and press Ctrl+S to save.\n"
		}
	} else {
		content += "No fact loaded for editing.\n"
	}

	return content
}

func (i *Intent) getDeleteConfirmContent() string {
	if i.context.FactToDelete == nil {
		return "No fact to delete"
	}

	var content string
	content += "Confirm Deletion\n\n"
	content += "Delete this fact?\n\n"
	content += fmt.Sprintf("Text: %s\n\n", truncate(i.context.FactToDelete.Text, 100))
	content += "Warning: This action cannot be undone.\n"

	return content
}

func (i *Intent) getResultsContent() string {
	var content string
	content += "Fact Statistics\n\n"
	content += fmt.Sprintf("Total facts in database: %d\n", i.context.TotalFacts)

	if len(i.context.Facts) > 0 {
		content += fmt.Sprintf("Currently viewing: %d facts\n", len(i.context.Facts))
	}

	return content
}

func (i *Intent) getContextHelp() string {
	theme := i.Theme()

	switch i.state {
	case StateList:
		return intents.CombineThemedFooters(
			intents.ThemedListFooter(theme),
			intents.ThemedCustomFooter(theme,
				primitives.EditBadge(theme),
				primitives.DeleteBadge(theme),
				primitives.HelpKeyBadge("n", "New", theme),
				primitives.HelpKeyBadge("r", "Refresh", theme),
			),
			intents.ThemedGlobalBadges(theme),
		)
	case StateView:
		return intents.CombineThemedFooters(
			intents.ThemedDetailViewFooter(theme),
			intents.ThemedCustomFooter(theme,
				primitives.EditBadge(theme),
				primitives.DeleteBadge(theme),
			),
			intents.ThemedGlobalBadges(theme),
		)
	case StateEditor:
		return intents.CombineThemedFooters(
			intents.ThemedFormFooter(theme),
			intents.ThemedCustomFooter(theme,
				primitives.SaveBadge(theme),
			),
			intents.ThemedGlobalBadges(theme),
		)
	case StateDeleteConfirm:
		return intents.CombineThemedFooters(
			intents.ThemedCustomFooter(theme,
				primitives.HelpKeyBadge("y/Enter", "Confirm", theme),
				primitives.HelpKeyBadge("n/Esc", "Cancel", theme),
			),
			intents.ThemedGlobalBadges(theme),
		)
	case StateResults:
		return intents.CombineThemedFooters(
			intents.ThemedCustomFooter(theme,
				primitives.BackBadge(theme),
			),
			intents.ThemedGlobalBadges(theme),
		)
	default:
		return intents.ThemedGlobalBadges(theme)
	}
}
