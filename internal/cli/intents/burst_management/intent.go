// Package burst_management implements the BurstManagement intent for managing career bursts.
package burst_management

import (
	tea "github.com/charmbracelet/bubbletea"
)

// Init is called when the intent is activated.
func (i *Intent) Init() tea.Cmd {
	// Load bursts from repository if available.
	if err := i.context.LoadBursts(); err != nil {
		// Error is acceptable - context.Bursts may already be populated (e.g., in tests).
		_ = err
	}

	// Initialize filtered bursts with the provided bursts.
	i.filteredBursts = i.context.Bursts
	if len(i.filteredBursts) > 0 {
		i.selectedBurst = i.filteredBursts[0]
	}

	i.state = StateList
	return nil
}

// Update processes a message in the intent.
func (i *Intent) Update(msg tea.Msg) tea.Cmd {
	if !i.active {
		return nil
	}

	// Handle messages based on current state.
	switch i.state {
	case StateList:
		return i.updateListView(msg)
	case StateDetail:
		return i.updateDetailView(msg)
	case StateDetailEvents:
		return i.updateDetailEventsView(msg)
	case StateDetailFacts:
		return i.updateDetailFactsView(msg)
	case StateEdit:
		return i.updateEditView(msg)
	case StateDeleteConfirm:
		return i.updateDeleteConfirmView(msg)
	case StateConfirm:
		return i.updateConfirmView(msg)
	case StateExtractingFacts:
		return i.updateExtractingFactsView(msg)
	case StateSuggesting:
		return i.updateSuggestingView(msg)
	case StateSuggestionReview:
		return i.updateSuggestionReviewView(msg)
	}

	return nil
}

// updateListView handles messages while viewing the burst list.
func (i *Intent) updateListView(msg tea.Msg) tea.Cmd {
	// Minimal implementation - will be expanded later.
	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch msg.String() {
		case "q":
			return tea.Quit
		case "esc":
			i.SetCancelled()
			return nil
		}
	}
	return nil
}

// updateDetailView handles messages while viewing burst details.
func (i *Intent) updateDetailView(msg tea.Msg) tea.Cmd {
	return nil
}

// updateDetailEventsView handles messages while viewing events in a burst.
func (i *Intent) updateDetailEventsView(msg tea.Msg) tea.Cmd {
	return nil
}

// updateDetailFactsView handles messages while viewing facts from a burst.
func (i *Intent) updateDetailFactsView(msg tea.Msg) tea.Cmd {
	return nil
}

// updateEditView handles the edit state.
func (i *Intent) updateEditView(msg tea.Msg) tea.Cmd {
	return nil
}

// updateDeleteConfirmView handles the delete confirmation state.
func (i *Intent) updateDeleteConfirmView(msg tea.Msg) tea.Cmd {
	return nil
}

// updateConfirmView handles the confirmation state.
func (i *Intent) updateConfirmView(msg tea.Msg) tea.Cmd {
	return nil
}

// updateExtractingFactsView handles the fact extraction state.
func (i *Intent) updateExtractingFactsView(msg tea.Msg) tea.Cmd {
	return nil
}

// updateSuggestingView handles the burst suggestion loading state.
func (i *Intent) updateSuggestingView(msg tea.Msg) tea.Cmd {
	return nil
}

// updateSuggestionReviewView handles the suggestion review state.
func (i *Intent) updateSuggestionReviewView(msg tea.Msg) tea.Cmd {
	return nil
}

// View renders the intent's current state.
func (i *Intent) View() string {
	if !i.active {
		return "BurstManagement intent is not active"
	}

	// Get content for current state.
	content := i.getStateContent()
	return content
}

// getStateContent returns the content for the current state.
func (i *Intent) getStateContent() string {
	switch i.state {
	case StateList:
		return i.viewList()
	case StateDetail:
		return i.viewDetail()
	case StateDetailEvents:
		return i.viewDetailEvents()
	case StateDetailFacts:
		return i.viewDetailFacts()
	case StateEdit:
		return i.viewEdit()
	case StateDeleteConfirm:
		return i.viewDeleteConfirm()
	case StateConfirm:
		return i.viewConfirm()
	case StateExtractingFacts:
		return i.viewExtractingFacts()
	case StateSuggesting:
		return i.viewSuggesting()
	case StateSuggestionReview:
		return i.viewSuggestionReview()
	}
	return "Unknown state"
}

// viewList renders the burst list view.
func (i *Intent) viewList() string {
	if len(i.filteredBursts) == 0 {
		return "No bursts found."
	}
	return "Burst list view"
}

// viewDetail renders the burst detail view.
func (i *Intent) viewDetail() string {
	if i.selectedBurst == nil {
		return "No burst selected."
	}
	return "Burst detail view"
}

// viewDetailEvents renders the events view for a burst.
func (i *Intent) viewDetailEvents() string {
	return "Burst events view"
}

// viewDetailFacts renders the facts view for a burst.
func (i *Intent) viewDetailFacts() string {
	return "Burst facts view"
}

// viewEdit renders the edit view for a burst.
func (i *Intent) viewEdit() string {
	return "Edit burst view"
}

// viewDeleteConfirm renders the delete confirmation view.
func (i *Intent) viewDeleteConfirm() string {
	return "Delete confirmation view"
}

// viewConfirm renders the confirmation view.
func (i *Intent) viewConfirm() string {
	return "Confirm burst view"
}

// viewExtractingFacts renders the fact extraction progress view.
func (i *Intent) viewExtractingFacts() string {
	return "Extracting facts view"
}

// viewSuggesting renders the burst suggestion loading view.
func (i *Intent) viewSuggesting() string {
	return "Suggesting bursts view"
}

// viewSuggestionReview renders the suggestion review view.
func (i *Intent) viewSuggestionReview() string {
	return "Suggestion review view"
}
