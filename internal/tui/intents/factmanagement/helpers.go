// Package factmanagement implements the FactManagement intent for managing career facts.
package factmanagement

import (
	domain "github.com/baphled/kariya/internal/domain/career"
	"github.com/baphled/kariya/internal/tui/intents"
	factviews "github.com/baphled/kariya/internal/tui/views/fact"
	"github.com/baphled/kariya/internal/ui/display"
	"github.com/baphled/kariya/internal/ui/uikit/primitives"
)

// syncTableSelection syncs the TableBehavior selection with the context.
func (i *Intent) syncTableSelection() {
	i.context.SelectedFactIndex = i.tableBehavior.GetSelectedIndex()
	if selected := i.tableBehavior.GetSelectedItem(); selected != nil {
		i.context.SelectedFact = i.findFactByID(selected.ID)
	} else {
		i.context.SelectedFact = nil
	}
}

func (i *Intent) findFactByID(id string) *domain.Fact {
	for idx := range i.context.Facts {
		fact := i.context.Facts[idx]
		if fact != nil && fact.ID == id {
			return fact
		}
	}

	return nil
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
			factName := "Fact #" + factID
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
			breadcrumbs = append(breadcrumbs, "Edit Fact #"+factID)
		}
	case StateResults:
		breadcrumbs = append(breadcrumbs, "Results")
	}

	return breadcrumbs
}

func (i *Intent) getStateContent() string {
	switch i.state {
	case StateList:
		return i.tableBehavior.Render()
	case StateView:
		return factviews.RenderFactDetail(display.FactFromDomain(i.context.SelectedFact))
	case StateEditor:
		return i.getEditorContent()
	case StateDeleteConfirm:
		return factviews.RenderDeleteConfirm(display.FactFromDomain(i.context.FactToDelete))
	case StateResults:
		return factviews.RenderResults(i.context.TotalFacts, len(i.context.Facts))
	case StateCompleted:
		return "Fact management completed"
	}
	return "Unknown state"
}

func (i *Intent) getEditorContent() string {
	if i.editModal != nil {
		return i.editModal.GetContent()
	}
	return factviews.RenderEditorFallback(display.FactFromDomain(i.context.EditingFact), i.context.HasFormErrors())
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
