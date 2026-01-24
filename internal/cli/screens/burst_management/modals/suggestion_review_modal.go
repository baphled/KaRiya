package modals

import (
	"fmt"
	"strings"

	"github.com/baphled/kariya/internal/cli/themes"
	"github.com/baphled/kariya/internal/cli/uikit/feedback"
	"github.com/baphled/kariya/internal/cli/uikit/primitives"
	themes2 "github.com/baphled/kariya/internal/cli/uikit/theme"
	"github.com/baphled/kariya/internal/service/career/burst_fact"
	tea "github.com/charmbracelet/bubbletea"
)

// SuggestionAction represents an action taken on a suggestion.
type SuggestionAction string

const (
	SuggestionActionAccept SuggestionAction = "accept"
	SuggestionActionReject SuggestionAction = "reject"
	SuggestionActionCancel SuggestionAction = "cancel"
)

// SuggestionReviewModal displays burst suggestions for review.
// User can navigate through suggestions and accept/reject them.
type SuggestionReviewModal struct {
	modal        *feedback.DetailModal
	suggestions  []burst_fact.BurstSuggestion
	currentIndex int
	theme        themes.Theme
	action       SuggestionAction
	accepted     []burst_fact.BurstSuggestion // Accepted suggestions
}

// NewSuggestionReviewModal creates a new suggestion review modal.
func NewSuggestionReviewModal(suggestions []burst_fact.BurstSuggestion, theme themes.Theme) *SuggestionReviewModal {
	if theme == nil {
		theme = themes.NewDefaultTheme()
	}

	m := &SuggestionReviewModal{
		suggestions:  suggestions,
		currentIndex: 0,
		theme:        theme,
		accepted:     []burst_fact.BurstSuggestion{},
	}

	// Build initial content and create modal.
	title, content := m.buildContent()
	m.modal = feedback.NewDetailModal(title, content)
	if theme != nil {
		m.modal = m.modal.WithTheme(theme)
	}

	// Set footer badges.
	m.updateFooterBadges()

	return m
}

// buildContent builds the title and content for the current suggestion.
func (m *SuggestionReviewModal) buildContent() (string, string) {
	if len(m.suggestions) == 0 {
		return "No Suggestions", "No suggestions available"
	}

	suggestion := m.suggestions[m.currentIndex]
	title := fmt.Sprintf("Burst Suggestion %d of %d", m.currentIndex+1, len(m.suggestions))

	theme := m.theme
	if theme == nil {
		theme = themes2.Default()
	}

	var content strings.Builder

	// Name.
	content.WriteString(primitives.NewText("Name:", theme).Bold().Render())
	content.WriteString(" " + suggestion.Name + "\n\n")

	// Description.
	if suggestion.Description != "" {
		content.WriteString(primitives.NewText("Description:", theme).Bold().Render())
		content.WriteString("\n" + suggestion.Description + "\n\n")
	}

	// Events count.
	content.WriteString(primitives.NewText("Events:", theme).Bold().Render())
	content.WriteString(fmt.Sprintf(" %d\n\n", len(suggestion.EventIDs)))

	// Confidence score.
	content.WriteString(primitives.NewText("Confidence:", theme).Bold().Render())
	content.WriteString(fmt.Sprintf(" %.1f%%\n", suggestion.ConfidenceScore*100))

	return title, content.String()
}

// updateFooterBadges updates the footer badges.
func (m *SuggestionReviewModal) updateFooterBadges() {
	theme := m.theme
	if theme == nil {
		theme = themes2.Default()
	}

	badges := []*primitives.Badge{
		primitives.HelpKeyBadge("a", "Accept", theme),
		primitives.HelpKeyBadge("r", "Reject", theme),
		primitives.HelpKeyBadge("n/→", "Next", theme),
		primitives.HelpKeyBadge("p/←", "Previous", theme),
		primitives.HelpKeyBadge("Esc", "Cancel", theme),
	}

	m.modal = m.modal.WithFooterBadges(badges...)
}

// updateContent updates the modal content for the current suggestion.
func (m *SuggestionReviewModal) updateContent() {
	title, content := m.buildContent()
	m.modal.SetContent(content)
	m.modal.SetTitle(title)
}

// Init initializes the modal.
func (m *SuggestionReviewModal) Init() tea.Cmd {
	return m.modal.Init()
}

// Update handles keyboard input.
func (m *SuggestionReviewModal) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	if !m.IsVisible() {
		return m, nil
	}

	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch msg.String() {
		case "esc":
			m.action = SuggestionActionCancel
			m.modal.Hide()
			return m, nil

		case "n", "right":
			// Next suggestion.
			if m.currentIndex < len(m.suggestions)-1 {
				m.currentIndex++
				m.updateContent()
			}
			return m, nil

		case "p", "left":
			// Previous suggestion.
			if m.currentIndex > 0 {
				m.currentIndex--
				m.updateContent()
			}
			return m, nil

		case "a":
			// Accept current suggestion.
			if m.currentIndex < len(m.suggestions) {
				m.action = SuggestionActionAccept
				current := m.suggestions[m.currentIndex]
				m.accepted = append(m.accepted, current)

				// Remove accepted suggestion from list.
				m.suggestions = append(m.suggestions[:m.currentIndex], m.suggestions[m.currentIndex+1:]...)

				// Adjust index if needed.
				if m.currentIndex >= len(m.suggestions) {
					m.currentIndex = len(m.suggestions) - 1
				}

				// If no suggestions left, close modal.
				if len(m.suggestions) == 0 {
					m.modal.Hide()
				} else {
					m.updateContent()
				}
			}
			return m, nil

		case "r":
			// Reject current suggestion.
			if m.currentIndex < len(m.suggestions) {
				m.action = SuggestionActionReject

				// Remove rejected suggestion from list.
				m.suggestions = append(m.suggestions[:m.currentIndex], m.suggestions[m.currentIndex+1:]...)

				// Adjust index if needed.
				if m.currentIndex >= len(m.suggestions) {
					m.currentIndex = len(m.suggestions) - 1
				}

				// If no suggestions left, close modal.
				if len(m.suggestions) == 0 {
					m.modal.Hide()
				} else {
					m.updateContent()
				}
			}
			return m, nil
		}
	}

	// Pass through to detail modal for scrolling, window resize, etc.
	var cmd tea.Cmd
	model, cmd := m.modal.Update(msg)
	if detailModal, ok := model.(*feedback.DetailModal); ok {
		m.modal = detailModal
	}
	return m, cmd
}

// View renders the modal.
func (m *SuggestionReviewModal) View() string {
	return m.modal.View()
}

// IsVisible returns whether the modal is visible.
func (m *SuggestionReviewModal) IsVisible() bool {
	return m.modal.IsVisible()
}

// Show makes the modal visible.
func (m *SuggestionReviewModal) Show() {
	m.modal.Show()
	m.action = ""
}

// Hide hides the modal.
func (m *SuggestionReviewModal) Hide() {
	m.modal.Hide()
}

// SetDimensions sets the terminal dimensions.
func (m *SuggestionReviewModal) SetDimensions(width, height int) {
	m.modal.SetDimensions(width, height)
}

// GetAction returns the last action taken.
func (m *SuggestionReviewModal) GetAction() SuggestionAction {
	return m.action
}

// ClearAction clears the current action.
func (m *SuggestionReviewModal) ClearAction() {
	m.action = ""
}

// GetAcceptedSuggestions returns all accepted suggestions.
func (m *SuggestionReviewModal) GetAcceptedSuggestions() []burst_fact.BurstSuggestion {
	return m.accepted
}

// GetCurrentSuggestion returns the currently displayed suggestion.
func (m *SuggestionReviewModal) GetCurrentSuggestion() *burst_fact.BurstSuggestion {
	if m.currentIndex < 0 || m.currentIndex >= len(m.suggestions) {
		return nil
	}
	return &m.suggestions[m.currentIndex]
}

// HasSuggestions returns whether there are any suggestions left.
func (m *SuggestionReviewModal) HasSuggestions() bool {
	return len(m.suggestions) > 0
}

// GetSuggestionsCount returns the number of remaining suggestions.
func (m *SuggestionReviewModal) GetSuggestionsCount() int {
	return len(m.suggestions)
}
