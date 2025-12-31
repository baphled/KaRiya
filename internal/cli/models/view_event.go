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

// EditEventMsg is sent to edit a specific event
type EditEventMsg struct {
	Event *career.CareerEvent
}

// EventDeletedMsg is sent when an event is successfully deleted
type EventDeletedMsg struct {
	EventID string
	Err     error
}

// ViewAction represents the action selected in the view
type ViewAction int

const (
	ViewActionNone ViewAction = iota
	ViewActionEdit
	ViewActionDelete
	ViewActionBack
)

// ViewEventModel represents the event detail view screen
type ViewEventModel struct {
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
}

// NewViewEventModel creates a new view event model
func NewViewEventModel(svc *careerservice.Service, ctx context.Context, event *career.CareerEvent) *ViewEventModel {
	return &ViewEventModel{
		service:          svc,
		ctx:              ctx,
		event:            event,
		selectedAction:   0,
		actions:          []string{"Edit Event", "Delete Event", "Back to List"},
		showDeleteDialog: false,
		deleteDialog:     NewConfirmationDialog("Confirm Delete", "Are you sure you want to delete this event? This action cannot be undone."),
	}
}

// Init initializes the model
func (m *ViewEventModel) Init() tea.Cmd {
	return nil
}

// Update handles messages
func (m *ViewEventModel) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
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
		}
	case tea.WindowSizeMsg:
		m.width = msg.Width
		m.height = msg.Height
	}

	return m, nil
}

// View renders the model
func (m *ViewEventModel) View() string {
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

// renderEventDetail renders the event details
func (m *ViewEventModel) renderEventDetail() string {
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
	instructions := styles.InfoText.Render("↑/↓ or j/k: Navigate | Enter: Select | e: Edit | d: Delete | Backspace/Esc: Back")

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

	return fullContent
}

// nextAction moves to the next action
func (m *ViewEventModel) nextAction() {
	if m.selectedAction < len(m.actions)-1 {
		m.selectedAction++
	}
}

// prevAction moves to the previous action
func (m *ViewEventModel) prevAction() {
	if m.selectedAction > 0 {
		m.selectedAction--
	}
}

// performAction executes the selected action
func (m *ViewEventModel) performAction() tea.Cmd {
	switch m.selectedAction {
	case 0: // Edit
		return func() tea.Msg {
			return EditEventMsg{Event: m.event}
		}
	case 1: // Delete
		m.showDeleteDialog = true
		return nil
	case 2: // Back
		return func() tea.Msg {
			return BackMsg{}
		}
	}
	return nil
}

// deleteEvent performs the actual deletion
func (m *ViewEventModel) deleteEvent() tea.Cmd {
	return func() tea.Msg {
		err := m.service.DeleteEvent(m.ctx, m.event.ID)
		return EventDeletedMsg{
			EventID: m.event.ID,
			Err:     err,
		}
	}
}

// GetEvent returns the current event
func (m *ViewEventModel) GetEvent() *career.CareerEvent {
	return m.event
}
