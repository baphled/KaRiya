package app

import (
	"github.com/baphled/kariya/internal/cli/uikit/feedback"
	careerrepo "github.com/baphled/kariya/internal/repository/career"
	tea "github.com/charmbracelet/bubbletea"
)

// handleMenuInput processes menu navigation and selection.
func (m *Model) handleMenuInput(msg tea.KeyMsg) (tea.Model, tea.Cmd) {
	switch msg.String() {
	case keyUp, keyK:
		// Move up in the menu.
		if m.selectedMenuIndex > 0 {
			m.selectedMenuIndex--
			return m, nil
		}
	case keyDown, keyJ:
		// Move down in the menu.
		if m.selectedMenuIndex < len(m.menuItems)-1 {
			m.selectedMenuIndex++
			return m, nil
		}
	case keyEnter, keySpace:
		return m.handleMenuSelection()
	}
	return m, nil
}

// handleMenuSelection handles selecting a menu item.
func (m *Model) handleMenuSelection() (tea.Model, tea.Cmd) {
	selectedItem := m.menuItems[m.selectedMenuIndex]

	// Generate CV requires at least one career event to be meaningful.
	if selectedItem.Intent == "generate_cv" {
		events, err := m.careerService.GetEventRepository().List(m.ctx, careerrepo.EventListFilters{Limit: 1})
		if err != nil || len(events) == 0 {
			// Show informational modal instead of activating intent.
			m.infoModal = feedback.NewWarningInfoModal(
				"No Career Events",
				"You need to add career events before generating a CV.\n\n"+
					"Use 'Capture Event' from the main menu to record your "+
					"achievements, projects, and career milestones.",
			)
			m.infoModal.SetDimensions(m.width, m.height)
			return m, nil
		}
	}

	m.state = StateIntent
	cmd, err := m.intentRouter.ActivateIntent(selectedItem.Intent, make(map[string]interface{}))
	if err != nil {
		m.logger.Error("Failed to activate intent %s: %v", selectedItem.Intent, err)
		return m, nil
	}
	if cmd == nil {
		cmd = func() tea.Msg { return nil }
	}
	return m, cmd
}

// handleIntentInput forwards messages to the active intent.
func (m *Model) handleIntentInput(msg tea.KeyMsg) (tea.Model, tea.Cmd) {
	cmd, result := m.intentRouter.HandleMessage(msg)

	// Check if intent has completed.
	if result != nil {
		m.state = StateMenu
		m.selectedMenuIndex = 0
		// Return command that will trigger the intent completed message.
		return m, tea.Batch(
			cmd,
			func() tea.Msg { return IntentCompletedMsg{} },
		)
	}

	return m, cmd
}
