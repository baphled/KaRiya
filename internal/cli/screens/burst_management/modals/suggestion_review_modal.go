package modals

import (
	"fmt"
	"sort"
	"strings"

	"github.com/baphled/kariya/internal/cli/behaviors"
	"github.com/baphled/kariya/internal/cli/themes"
	"github.com/baphled/kariya/internal/cli/uikit/containers"
	"github.com/baphled/kariya/internal/cli/uikit/primitives"
	themes2 "github.com/baphled/kariya/internal/cli/uikit/theme"
	"github.com/baphled/kariya/internal/service/career/burst_fact"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
)

// SuggestionAction represents an action taken on a suggestion.
type SuggestionAction string

const (
	SuggestionActionAccept SuggestionAction = "accept"
	SuggestionActionReject SuggestionAction = "reject"
	SuggestionActionCancel SuggestionAction = "cancel"
)

// SuggestionReviewModal displays all burst suggestions for review in a table format.
// Suggestions are sorted by confidence (highest first).
// User can navigate through suggestions and accept/reject them.
type SuggestionReviewModal struct {
	table       *behaviors.TableBehavior[burst_fact.BurstSuggestion]
	suggestions []burst_fact.BurstSuggestion
	theme       themes.Theme
	action      SuggestionAction
	accepted    []burst_fact.BurstSuggestion
	visible     bool
	width       int
	height      int
}

// NewSuggestionReviewModal creates a new suggestion review modal.
// Suggestions are automatically sorted by confidence (highest first).
func NewSuggestionReviewModal(suggestions []burst_fact.BurstSuggestion, theme themes.Theme) *SuggestionReviewModal {
	if theme == nil {
		theme = themes.NewDefaultTheme()
	}

	// Sort suggestions by confidence (highest first).
	sortedSuggestions := make([]burst_fact.BurstSuggestion, len(suggestions))
	copy(sortedSuggestions, suggestions)
	sort.Slice(sortedSuggestions, func(i, j int) bool {
		return sortedSuggestions[i].ConfidenceScore > sortedSuggestions[j].ConfidenceScore
	})

	// Define table columns.
	columns := []behaviors.ColumnDef{
		{Title: "Name", Width: 30},
		{Title: "Events", Width: 8},
		{Title: "Confidence", Width: 18},
	}

	// Row formatter for suggestions.
	formatter := func(s burst_fact.BurstSuggestion, _ int) []string {
		confidenceBar := primitives.CompactBar(s.ConfidenceScore, 8, nil).
			ShowPercentage(true).
			Render()
		return []string{
			s.Name,
			fmt.Sprintf("%d", len(s.EventIDs)),
			confidenceBar,
		}
	}

	// Create table behavior with pagination.
	table := behaviors.NewTableBehavior(theme, columns, formatter).
		EmptyMessage("No suggestions available").
		PaginationPrefix("Suggestions").
		PageSize(10)

	table.SetItems(sortedSuggestions)

	m := &SuggestionReviewModal{
		table:       table,
		suggestions: sortedSuggestions,
		theme:       theme,
		accepted:    []burst_fact.BurstSuggestion{},
		visible:     true,
		width:       80,
		height:      24,
	}

	return m
}

// getTheme returns the theme or default if nil.
func (m *SuggestionReviewModal) getTheme() themes.Theme {
	if m.theme != nil {
		return m.theme
	}
	return themes2.Default()
}

// buildContent builds the content showing all suggestions.
func (m *SuggestionReviewModal) buildContent() string {
	if len(m.suggestions) == 0 {
		return "No suggestions available"
	}

	theme := m.getTheme()
	var content strings.Builder

	// Render the table.
	content.WriteString(m.table.Render())
	content.WriteString("\n\n")

	// Show selected suggestion details.
	selected := m.table.GetSelectedItem()
	if selected != nil {
		content.WriteString(primitives.NewText("Selected:", theme).Bold().Render())
		content.WriteString(" " + selected.Name + "\n")
		if selected.Description != "" {
			content.WriteString(primitives.NewText("Description:", theme).Bold().Render())
			content.WriteString(" " + selected.Description + "\n")
		}
	}

	return content.String()
}

// buildFooter builds the footer with help badges.
func (m *SuggestionReviewModal) buildFooter() string {
	theme := m.getTheme()

	badges := []*primitives.Badge{
		primitives.HelpKeyBadge("a", "Accept", theme),
		primitives.HelpKeyBadge("r", "Reject", theme),
		primitives.HelpKeyBadge("j/k", "Navigate", theme),
		primitives.HelpKeyBadge("n/p", "Page", theme),
		primitives.HelpKeyBadge("Esc", "Cancel", theme),
	}

	return primitives.RenderHelpFooter(theme, badges...)
}

// Init initializes the modal.
func (m *SuggestionReviewModal) Init() tea.Cmd {
	return nil
}

// Update handles keyboard input.
func (m *SuggestionReviewModal) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	if !m.IsVisible() {
		return m, nil
	}

	switch msg := msg.(type) {
	case tea.WindowSizeMsg:
		m.width = msg.Width
		m.height = msg.Height
		m.table.Dimensions(m.width-12, m.height-16)
		return m, nil

	case tea.KeyMsg:
		switch msg.String() {
		case "esc":
			m.action = SuggestionActionCancel
			m.visible = false
			return m, nil

		case "j", "down":
			m.table.HandleNavigation("down")
			return m, nil

		case "k", "up":
			m.table.HandleNavigation("up")
			return m, nil

		case "pgdown", "n":
			m.table.HandleNavigation("pgdn")
			return m, nil

		case "pgup", "p":
			m.table.HandleNavigation("pgup")
			return m, nil

		case "a":
			// Accept current suggestion.
			selected := m.table.GetSelectedItem()
			if selected != nil {
				m.action = SuggestionActionAccept
				m.accepted = append(m.accepted, *selected)

				// Remove accepted suggestion from list.
				m.removeCurrentSuggestion()

				// If no suggestions left, close modal.
				if len(m.suggestions) == 0 {
					m.visible = false
				}
			}
			return m, nil

		case "r":
			// Reject current suggestion.
			selected := m.table.GetSelectedItem()
			if selected != nil {
				m.action = SuggestionActionReject

				// Remove rejected suggestion from list.
				m.removeCurrentSuggestion()

				// If no suggestions left, close modal.
				if len(m.suggestions) == 0 {
					m.visible = false
				}
			}
			return m, nil
		}
	}

	return m, nil
}

// removeCurrentSuggestion removes the currently selected suggestion from the list.
func (m *SuggestionReviewModal) removeCurrentSuggestion() {
	idx := m.table.GetSelectedIndex()
	if idx >= 0 && idx < len(m.suggestions) {
		m.suggestions = append(m.suggestions[:idx], m.suggestions[idx+1:]...)
		m.table.SetItems(m.suggestions)

		// Preserve selection at same index, or adjust if at end.
		if idx >= len(m.suggestions) && len(m.suggestions) > 0 {
			idx = len(m.suggestions) - 1
		}
		if len(m.suggestions) > 0 {
			m.table.SetSelectedIndex(idx)
		}
	}
}

// View renders the modal.
func (m *SuggestionReviewModal) View() string {
	if !m.visible {
		return ""
	}

	theme := m.getTheme()

	// Build title.
	title := primitives.Title("Review Burst Suggestions", theme).Render()

	// Build content (table + selected details).
	content := m.buildContent()

	// Build footer.
	footer := m.buildFooter()

	// Combine all parts.
	modalContent := lipgloss.JoinVertical(lipgloss.Left, title, "", content, "", footer)

	// Calculate modal dimensions.
	maxModalHeight := 30
	terminalMaxHeight := int(float64(m.height) * 0.8)
	if terminalMaxHeight < maxModalHeight {
		maxModalHeight = terminalMaxHeight
	}
	if maxModalHeight < 15 {
		maxModalHeight = 15
	}

	modalWidth := m.width - 12
	if modalWidth < 60 {
		modalWidth = 60
	}
	if modalWidth > 90 {
		modalWidth = 90
	}

	// Wrap in styled box with solid background.
	return containers.NewBox(theme).
		Content(modalContent).
		Width(modalWidth).
		Height(maxModalHeight).
		Padding(2).
		Background(theme.BackgroundColor()).
		Render()
}

// IsVisible returns whether the modal is visible.
func (m *SuggestionReviewModal) IsVisible() bool {
	return m.visible
}

// Show makes the modal visible.
func (m *SuggestionReviewModal) Show() {
	m.visible = true
	m.action = ""
}

// Hide hides the modal.
func (m *SuggestionReviewModal) Hide() {
	m.visible = false
}

// SetDimensions sets the terminal dimensions.
func (m *SuggestionReviewModal) SetDimensions(width, height int) {
	m.width = width
	m.height = height

	// Calculate the actual content width available inside the modal.
	// Modal width is capped at 90 with padding of 2 on each side.
	modalWidth := width - 12
	if modalWidth < 60 {
		modalWidth = 60
	}
	if modalWidth > 90 {
		modalWidth = 90
	}
	// Subtract padding (2 on each side = 4 total) and some margin.
	contentWidth := modalWidth - 8

	m.table.Dimensions(contentWidth, height-16)
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

// GetCurrentSuggestion returns the currently selected suggestion.
func (m *SuggestionReviewModal) GetCurrentSuggestion() *burst_fact.BurstSuggestion {
	return m.table.GetSelectedItem()
}

// HasSuggestions returns whether there are any suggestions left.
func (m *SuggestionReviewModal) HasSuggestions() bool {
	return len(m.suggestions) > 0
}

// GetSuggestionsCount returns the number of remaining suggestions.
func (m *SuggestionReviewModal) GetSuggestionsCount() int {
	return len(m.suggestions)
}

// GetAllSuggestions returns all suggestions in their current order (sorted by confidence).
func (m *SuggestionReviewModal) GetAllSuggestions() []burst_fact.BurstSuggestion {
	return m.suggestions
}
