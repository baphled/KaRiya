// Package fact_management implements the FactManagement intent for managing career facts.
package fact_management

import (
	"fmt"

	"github.com/baphled/kariya/internal/cli/intents"
	"github.com/baphled/kariya/internal/cli/screens"
	factmodals "github.com/baphled/kariya/internal/cli/screens/fact/modals"
	"github.com/baphled/kariya/internal/cli/uikit/feedback"
	"github.com/baphled/kariya/internal/cli/uikit/primitives"
	domain "github.com/baphled/kariya/internal/domain/career"
	tea "github.com/charmbracelet/bubbletea"
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

// transitionToScreen sets the active screen.
func (i *Intent) transitionToScreen(screen screens.Screen) {
	i.activeScreen = screen

	termInfo := i.GetTerminalInfo()
	if termInfo != nil && termInfo.Width > 0 && termInfo.Height > 0 {
		screen.SetTerminalInfo(termInfo.Width, termInfo.Height)
	}

	if theme := i.Theme(); theme != nil {
		screen.SetTheme(theme)
	}

	if logo := i.GetLogo(); logo != nil {
		screen.SetLogo(logo, i.GetLogoSpacing())
	}
}

// hasActiveModal returns true if any modal is visible.
func (i *Intent) hasActiveModal() bool {
	return i.errorModal != nil ||
		(i.deleteModal != nil && i.deleteModal.IsVisible()) ||
		(i.editModal != nil && i.editModal.IsVisible()) ||
		(i.detailModal != nil && i.detailModal.IsVisible())
}

// rebuildModalRegistry rebuilds the modal registry with current modals.
func (i *Intent) rebuildModalRegistry() {
	if i.modalRegistry == nil {
		i.modalRegistry = intents.NewModalRegistry()
	}
	i.modalRegistry.Clear()

	termInfo := i.GetTerminalInfo()
	width, height := 80, 24
	if termInfo != nil && termInfo.Width > 0 && termInfo.Height > 0 {
		width, height = termInfo.Width, termInfo.Height
	}

	// Error modal (highest priority).
	if i.errorModal != nil {
		i.modalRegistry.Register(intents.NewErrorModalAdapter(i.errorModal, width, height, i.Theme()))
	}

	// Delete confirmation modal.
	if i.deleteModal != nil {
		i.modalRegistry.Register(intents.NewConfirmModalAdapter(i.deleteModal))
	}

	// Edit modal.
	if i.editModal != nil && i.editModal.IsVisible() {
		i.modalRegistry.Register(intents.NewViewModalAdapter(
			func() bool { return i.editModal != nil && i.editModal.IsVisible() },
			func() string { return i.editModal.View() },
			func(msg tea.Msg) (tea.Model, tea.Cmd) {
				cmd, _ := i.editModal.Update(msg)
				return nil, cmd
			},
		))
	}

	// Detail modal.
	if i.detailModal != nil && i.detailModal.IsVisible() {
		i.modalRegistry.Register(intents.NewViewModalAdapter(
			func() bool { return i.detailModal != nil && i.detailModal.IsVisible() },
			func() string { return i.detailModal.View() },
			func(msg tea.Msg) (tea.Model, tea.Cmd) {
				_, cmd := i.detailModal.Update(msg)
				return nil, cmd
			},
		))
	}
}

// Modal opening helpers.

func (i *Intent) openDetailModal(fact *domain.Fact) {
	i.detailModal = factmodals.NewDetailModal(fact)
	i.detailModal.SetTheme(i.Theme())
	i.detailModal.Show()
	i.state = StateView
}

func (i *Intent) openEditModal(fact *domain.Fact, isNew bool) {
	i.editModal = factmodals.NewEditModal(fact, isNew)
	i.editModal.SetTheme(i.Theme())
	i.editModal.Show()
	i.state = StateEditor
}

func (i *Intent) openDeleteConfirm(fact *domain.Fact) {
	i.context.FactToDelete = fact
	text := truncate(fact.Text, 50)
	i.deleteModal = feedback.NewConfirmModal(
		"Delete Fact",
		fmt.Sprintf("Are you sure you want to delete '%s'?", text),
	).WithVariant(feedback.ConfirmDestructive)
	i.state = StateDeleteConfirm
}

func (i *Intent) showErrorModal(title, message string) {
	i.errorModal = feedback.NewErrorModal(title, message)
}

// performDelete deletes the pending fact.
func (i *Intent) performDelete() tea.Cmd {
	if i.context.FactToDelete == nil {
		return nil
	}

	if err := i.context.DeleteFact(i.context.FactToDelete.ID); err != nil {
		i.result = &intents.IntentResult[*Result]{
			Status: intents.Failed,
			Error: &intents.IntentError{
				Code:    "DELETE_FAILED",
				Message: "Failed to delete fact",
				Cause:   err,
			},
		}
	} else {
		i.result = &intents.IntentResult[*Result]{
			Status: intents.Completed,
			Data: &Result{
				Action:  "deleted",
				Fact:    i.context.FactToDelete,
				Facts:   i.context.Facts,
				Message: "Fact deleted successfully",
			},
		}
		i.tableBehavior.SetItems(i.context.Facts)
		i.syncTableSelection()
	}
	i.context.FactToDelete = nil
	i.state = StateList
	return nil
}

// View helper methods.

func (i *Intent) getBreadcrumbs() []string {
	breadcrumbs := []string{"Main Menu", "Manage Facts"}

	switch i.state {
	case StateView, StateEditor, StateDeleteConfirm:
		if i.context.SelectedFact != nil {
			factID := i.context.SelectedFact.ID
			if len(factID) > 8 {
				factID = factID[:8]
			}
			factName := fmt.Sprintf("Fact #%s", factID)
			breadcrumbs = append(breadcrumbs, factName)
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
	if i.editModal != nil {
		return i.editModal.GetContent()
	}

	var content string
	content += "Edit Fact\n\n"

	if i.context.EditingFact != nil {
		content += fmt.Sprintf("Text: %s\n\n", i.context.EditingFact.Text)
		content += fmt.Sprintf("Categories: %v\n", i.context.EditingFact.CompetencyCategories)
		content += fmt.Sprintf("Strength Signal: %s\n", i.context.EditingFact.StrengthSignal)
		content += fmt.Sprintf("Role Fit: %v\n\n", i.context.EditingFact.RoleFit)
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
