package models

import (
	"context"
	"fmt"
	"strings"

	"github.com/baphled/kariya/internal/cli/components"
	"github.com/baphled/kariya/internal/cli/styles"
	"github.com/baphled/kariya/internal/domain/career"
	careerservice "github.com/baphled/kariya/internal/service/career"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
)

// ViewEventWithFactsModel represents the event detail view with facts display
type ViewEventWithFactsModel struct {
	service          *careerservice.Service
	ctx              context.Context
	event            *career.CareerEvent
	selectedAction   int
	actions          []string
	showDeleteDialog bool
	deleteDialog     *ConfirmationDialog
	err              error
	width            int
	height           int
	helpFooter       components.HelpFooterModel

	// Facts display
	factsLoaded     bool
	eventFacts      []*career.Fact // Facts extracted from this event
	burstFacts      []*career.Fact // Facts extracted from bursts containing this event
	showFactsList   bool
	factListModel   *FactListModel
	showFactEditor  bool
	factEditorModel *FactEditorModel
	selectedFactIdx int
}

// FactsLoadedMsg is sent when facts are loaded
type FactsLoadedMsg struct {
	EventFacts []*career.Fact
	BurstFacts []*career.Fact
	Err        error
}

// NewViewEventWithFactsModel creates a new view event with facts model
func NewViewEventWithFactsModel(svc *careerservice.Service, ctx context.Context, event *career.CareerEvent) *ViewEventWithFactsModel {
	return &ViewEventWithFactsModel{
		service:          svc,
		ctx:              ctx,
		event:            event,
		selectedAction:   0,
		actions:          []string{"Edit Event", "Delete Event", "View Facts", "Back to List"},
		showDeleteDialog: false,
		deleteDialog:     NewConfirmationDialog("Confirm Delete", "Are you sure you want to delete this event? This action cannot be undone."),
		helpFooter:       components.NewHelpFooter("view_event", 80),
		eventFacts:       []*career.Fact{},
		burstFacts:       []*career.Fact{},
	}
}

// Init initializes the model and loads facts
func (m *ViewEventWithFactsModel) Init() tea.Cmd {
	return m.loadFacts()
}

// Update handles messages
func (m *ViewEventWithFactsModel) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	// If showing fact editor, delegate to it
	if m.showFactEditor && m.factEditorModel != nil {
		updatedModel, cmd := m.factEditorModel.Update(msg)
		m.factEditorModel = updatedModel.(*FactEditorModel)

		// Check if editor is done
		if m.factEditorModel.IsSubmitted() {
			// Save the edited fact
			fact := m.factEditorModel.GetFact()
			m.showFactEditor = false
			return m, func() tea.Msg {
				return SaveFactMsg{Fact: fact}
			}
		} else if m.factEditorModel.IsCancelled() {
			// Cancel editing
			m.showFactEditor = false
		}

		return m, cmd
	}

	// If showing fact list, delegate to it
	if m.showFactsList && m.factListModel != nil {
		updatedModel, cmd := m.factListModel.Update(msg)
		m.factListModel = updatedModel.(*FactListModel)

		// Check if user wants to edit a fact
		if m.factListModel.IsSubmitted() && len(m.factListModel.GetFacts()) > 0 {
			selectedFact := m.factListModel.GetFacts()[m.factListModel.GetSelectedIdx()]
			m.showFactEditor = true
			m.factEditorModel = NewFactEditorModel(selectedFact, m.service, m.ctx)
			return m, m.factEditorModel.Init()
		}

		// Check if user wants to go back
		if m.factListModel.IsCancelled() {
			m.showFactsList = false
		}

		return m, cmd
	}

	// If showing delete dialog, delegate to it
	if m.showDeleteDialog {
		var cmd tea.Cmd
		m.deleteDialog, cmd = m.deleteDialog.Update(msg)

		// Check if dialog is done
		if m.deleteDialog.IsConfirmed() {
			// Perform deletion
			return m, m.deleteEvent()
		} else if m.deleteDialog.IsCancelled() {
			// Close dialog
			m.showDeleteDialog = false
			m.deleteDialog.Reset()
		}

		return m, cmd
	}

	switch msg := msg.(type) {
	case FactsLoadedMsg:
		m.factsLoaded = true
		m.eventFacts = msg.EventFacts
		m.burstFacts = msg.BurstFacts
		if msg.Err != nil {
			m.err = msg.Err
		}

	case SaveFactMsg:
		// Save fact to repository
		return m, func() tea.Msg {
			err := m.service.SaveFact(m.ctx, msg.Fact)
			if err != nil {
				return FactConfirmedMsg{Fact: msg.Fact, Err: err}
			}
			// Reload facts after save
			return m.loadFacts()()
		}

	case RejectFactMsg:
		// Delete fact from repository
		return m, func() tea.Msg {
			err := m.service.DeleteFact(m.ctx, msg.FactID)
			return FactRejectedMsg{FactID: msg.FactID, Err: err}
		}

	case FactConfirmedMsg:
		if msg.Err != nil {
			m.err = msg.Err
		}
		// Refresh facts list
		return m, m.loadFacts()

	case FactRejectedMsg:
		if msg.Err != nil {
			m.err = msg.Err
		}
		// Refresh facts list
		return m, m.loadFacts()

	case tea.KeyMsg:
		switch msg.String() {
		case "up", "k":
			m.prevAction()
		case "down", "j":
			m.nextAction()
		case "enter":
			return m, m.performAction()
		case "esc":
			// Signal back navigation to parent
			return m, func() tea.Msg { return BackMsg{} }
		case "ctrl+c", "q":
			// Signal quit to parent
			return m, func() tea.Msg { return QuitMsg{} }
		case "e":
			// Quick edit shortcut
			return m, func() tea.Msg {
				return EditEventMsg{Event: m.event}
			}
		case "d":
			// Quick delete shortcut - show confirmation
			m.showDeleteDialog = true
		case "f":
			// Quick view facts shortcut
			m.showFactsList = true
			m.factListModel = NewFactListModel(m.service, m.ctx)
			allFacts := append(m.eventFacts, m.burstFacts...)
			m.factListModel.SetFacts(allFacts)
		}
	case tea.WindowSizeMsg:
		m.width = msg.Width
		m.height = msg.Height
	}

	return m, nil
}

// View renders the model
func (m *ViewEventWithFactsModel) View() string {
	// If showing fact editor, render it
	if m.showFactEditor && m.factEditorModel != nil {
		return m.factEditorModel.View()
	}

	// If showing fact list, render it
	if m.showFactsList && m.factListModel != nil {
		return m.factListModel.View()
	}

	// If showing delete dialog, render it on top
	if m.showDeleteDialog {
		// Render event detail in background (dimmed)
		background := m.renderEventDetail()
		dialog := m.deleteDialog.View()

		// Center the dialog
		centeredDialog := lipgloss.Place(
			m.width,
			m.height,
			lipgloss.Center,
			lipgloss.Center,
			dialog,
		)

		// Overlay dialog on background
		return lipgloss.JoinVertical(lipgloss.Left, background, centeredDialog)
	}

	return m.renderEventDetail()
}

// renderEventDetail renders the event details with facts section
func (m *ViewEventWithFactsModel) renderEventDetail() string {
	// Title
	header := components.NewHeader("Event Details", m.width)
	title := header.View()

	// Event card
	var eventDetails []string

	// Event text
	eventDetails = append(eventDetails,
		styles.HeaderSection.Render("Description"),
		lipgloss.NewStyle().Render(m.event.Text),
		"",
	)

	// Date
	dateStr := m.event.Date.Format("Monday, January 2, 2006")
	eventDetails = append(eventDetails,
		styles.HeaderSubsection.Render("Date"),
		lipgloss.NewStyle().Render(dateStr),
		"",
	)

	// Company (if present)
	if m.event.Company != "" {
		eventDetails = append(eventDetails,
			styles.HeaderSubsection.Render("Company"),
			lipgloss.NewStyle().Render(m.event.Company),
			"",
		)
	}

	// Project (if present)
	if m.event.Project != "" {
		eventDetails = append(eventDetails,
			styles.HeaderSubsection.Render("Project"),
			lipgloss.NewStyle().Render(m.event.Project),
			"",
		)
	}

	// Tags (if present)
	if len(m.event.Tags) > 0 {
		tagsList := make([]string, len(m.event.Tags))
		for i, tag := range m.event.Tags {
			tagsList[i] = styles.TagBase.Render(tag)
		}
		tagsDisplay := strings.Join(tagsList, " ")

		eventDetails = append(eventDetails,
			styles.HeaderSubsection.Render("Tags"),
			tagsDisplay,
			"",
		)
	}

	// Categories (if present)
	if len(m.event.Categories) > 0 {
		categoriesList := make([]string, len(m.event.Categories))
		for i, category := range m.event.Categories {
			categoriesList[i] = styles.TagBase.
				Foreground(styles.ColorAccentTeal).
				Render(category)
		}
		categoriesDisplay := strings.Join(categoriesList, " ")

		eventDetails = append(eventDetails,
			styles.HeaderSubsection.Render("Categories"),
			categoriesDisplay,
			"",
		)
	}

	// Facts section (if loaded)
	if m.factsLoaded {
		eventFactCount := len(m.eventFacts)
		burstFactCount := len(m.burstFacts)
		totalFacts := eventFactCount + burstFactCount

		if totalFacts > 0 {
			eventDetails = append(eventDetails,
				styles.HeaderSection.Render("Extracted Facts"),
				"",
			)

			// Summary
			summary := fmt.Sprintf("Total: %d facts (%d from event, %d from bursts)",
				totalFacts, eventFactCount, burstFactCount)
			eventDetails = append(eventDetails, styles.InfoText.Render(summary), "")

			// Show preview of up to 3 facts
			allFacts := append(m.eventFacts, m.burstFacts...)
			previewCount := 3
			if len(allFacts) < previewCount {
				previewCount = len(allFacts)
			}

			for i := 0; i < previewCount; i++ {
				fact := allFacts[i]
				factPreview := m.renderFactPreview(fact)
				eventDetails = append(eventDetails, factPreview, "")
			}

			if totalFacts > previewCount {
				more := fmt.Sprintf("... and %d more facts (press 'f' to view all)", totalFacts-previewCount)
				eventDetails = append(eventDetails, styles.InfoText.Render(more))
			}
		} else {
			eventDetails = append(eventDetails,
				styles.HeaderSection.Render("Facts"),
				styles.InfoText.Render("No facts extracted yet"),
				"",
			)
		}
	}

	// Created/Updated timestamps
	createdStr := m.event.CreatedAt.Format("2006-01-02 15:04")
	updatedStr := m.event.UpdatedAt.Format("2006-01-02 15:04")
	timestampInfo := fmt.Sprintf("Created: %s  |  Updated: %s", createdStr, updatedStr)
	eventDetails = append(eventDetails,
		lipgloss.NewStyle().
			Foreground(styles.ColorTextMuted).
			Render(timestampInfo),
	)

	// Combine event details
	eventCard := lipgloss.JoinVertical(lipgloss.Left, eventDetails...)

	// Actions menu
	var actionItems []string
	actionItems = append(actionItems, styles.HeaderSection.Render("Actions"), "")

	for i, action := range m.actions {
		itemStyle := styles.ButtonSecondary
		marker := "  "
		if i == m.selectedAction {
			itemStyle = styles.ButtonFocused
			marker = "▶ "
		}

		actionItems = append(actionItems, itemStyle.Render(marker+action))
	}

	actionsMenu := lipgloss.JoinVertical(lipgloss.Left, actionItems...)

	// Instructions
	instructions := styles.InfoText.Render("↑/↓ or j/k: Navigate | Enter: Select | e: Edit | d: Delete | f: View Facts | Esc: Back")

	// Wrap in cards
	eventCardRendered := styles.CardBase.
		Width(styles.MaxWidth(m.width) - 4).
		Render(eventCard)

	actionsCardRendered := styles.CardBase.
		Width(styles.MaxWidth(m.width) - 4).
		Render(actionsMenu)

	// Combine all sections
	fullContent := lipgloss.JoinVertical(
		lipgloss.Left,
		title,
		"",
		eventCardRendered,
		"",
		actionsCardRendered,
		"",
		instructions,
	)

	// Help footer
	m.helpFooter.SetWidth(m.width)
	fullContent = lipgloss.JoinVertical(lipgloss.Left, fullContent, m.helpFooter.View())

	return fullContent
}

// renderFactPreview renders a preview of a fact
func (m *ViewEventWithFactsModel) renderFactPreview(fact *career.Fact) string {
	// Truncate text if too long
	text := fact.Text
	if len(text) > 100 {
		text = text[:97] + "..."
	}

	// Role fit badge
	roleFitBadge := styles.TagBase.
		Foreground(styles.ColorAccentPurple).
		Render(string(fact.RoleFit))

	// Competencies
	compBadges := []string{}
	for _, comp := range fact.CompetencyCategories {
		badge := styles.TagBase.
			Foreground(styles.ColorAccentTeal).
			Render(comp)
		compBadges = append(compBadges, badge)
	}

	// Build preview
	preview := fmt.Sprintf("• %s", text)
	badges := lipgloss.JoinHorizontal(lipgloss.Left, roleFitBadge, " ", strings.Join(compBadges, " "))

	return lipgloss.JoinVertical(lipgloss.Left, preview, "  "+badges)
}

// loadFacts loads facts for the event
func (m *ViewEventWithFactsModel) loadFacts() tea.Cmd {
	return func() tea.Msg {
		// Get facts extracted from this event
		eventFacts, err := m.service.GetFactsBySourceEventID(m.ctx, m.event.ID)
		if err != nil {
			return FactsLoadedMsg{EventFacts: []*career.Fact{}, BurstFacts: []*career.Fact{}, Err: err}
		}

		// TODO: Get facts from bursts containing this event
		// For now, just return event facts
		burstFacts := []*career.Fact{}

		return FactsLoadedMsg{EventFacts: eventFacts, BurstFacts: burstFacts, Err: nil}
	}
}

// nextAction moves to the next action
func (m *ViewEventWithFactsModel) nextAction() {
	if m.selectedAction < len(m.actions)-1 {
		m.selectedAction++
	}
}

// prevAction moves to the previous action
func (m *ViewEventWithFactsModel) prevAction() {
	if m.selectedAction > 0 {
		m.selectedAction--
	}
}

// performAction executes the selected action
func (m *ViewEventWithFactsModel) performAction() tea.Cmd {
	switch m.selectedAction {
	case 0: // Edit
		return func() tea.Msg {
			return EditEventMsg{Event: m.event}
		}
	case 1: // Delete
		m.showDeleteDialog = true
		return nil
	case 2: // View Facts
		m.showFactsList = true
		m.factListModel = NewFactListModel(m.service, m.ctx)
		allFacts := append(m.eventFacts, m.burstFacts...)
		m.factListModel.SetFacts(allFacts)
		return nil
	case 3: // Back
		return func() tea.Msg {
			return BackMsg{}
		}
	}
	return nil
}

// deleteEvent performs the actual deletion
func (m *ViewEventWithFactsModel) deleteEvent() tea.Cmd {
	return func() tea.Msg {
		err := m.service.DeleteEvent(m.ctx, m.event.ID)
		return EventDeletedMsg{
			EventID: m.event.ID,
			Err:     err,
		}
	}
}

// GetEvent returns the current event
func (m *ViewEventWithFactsModel) GetEvent() *career.CareerEvent {
	return m.event
}
