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
	"github.com/baphled/kariya/internal/service/career/burstfact"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
)

// SuggestionAction represents an action taken on a suggestion.
type SuggestionAction string

// SuggestionAction constants represent the decisions a user can make on an
// individual burst suggestion during the review workflow. The user sees a
// paginated table of suggestions sorted by confidence score; navigating
// with j/k or arrow keys highlights a row, and pressing the corresponding
// key applies an action to the highlighted suggestion.
const (
	// SuggestionActionAccept marks the highlighted suggestion as approved
	// and moves it to the accepted list. The suggestion is removed from the
	// review table and will be persisted as a confirmed burst fact when the
	// review session ends. The user presses a to trigger this action.
	SuggestionActionAccept SuggestionAction = "accept"
	// SuggestionActionReject discards the highlighted suggestion, removing
	// it from the review table without recording it as a fact. The
	// suggestion cannot be recovered after rejection. The user presses r to
	// trigger this action.
	SuggestionActionReject SuggestionAction = "reject"
	// SuggestionActionCancel aborts the entire review session, closing the
	// modal and discarding all pending accept or reject decisions made so
	// far. No suggestions are persisted. The user presses Escape to trigger
	// this action.
	SuggestionActionCancel SuggestionAction = "cancel"
)

// SuggestionReviewModal displays all burst suggestions for review in a table format.
// Suggestions are sorted by confidence (highest first).
// User can navigate through suggestions and accept/reject them.
type SuggestionReviewModal struct {
	table       *behaviors.TableBehavior[burstfact.BurstSuggestion]
	suggestions []burstfact.BurstSuggestion
	theme       themes.Theme
	action      SuggestionAction
	accepted    []burstfact.BurstSuggestion
	visible     bool
	width       int
	height      int
}

// NewSuggestionReviewModal creates a new suggestion review modal.
// Suggestions are automatically sorted by confidence (highest first).
func NewSuggestionReviewModal(suggestions []burstfact.BurstSuggestion, theme themes.Theme) *SuggestionReviewModal {
	if theme == nil {
		theme = themes.NewDefaultTheme()
	}

	sortedSuggestions := make([]burstfact.BurstSuggestion, len(suggestions))
	copy(sortedSuggestions, suggestions)
	sort.Slice(sortedSuggestions, func(i, j int) bool {
		return sortedSuggestions[i].ConfidenceScore > sortedSuggestions[j].ConfidenceScore
	})

	columns := []behaviors.ColumnDef{
		{Title: "Name", Width: 30},
		{Title: "Events", Width: 8},
		{Title: "Confidence", Width: 24},
	}

	formatter := func(s burstfact.BurstSuggestion, _ int) []string {
		confidenceBar := primitives.CompactBar(s.ConfidenceScore, 15, nil).
			ShowPercentage(true).
			Render()
		return []string{
			s.Name,
			fmt.Sprintf("%d", len(s.EventIDs)),
			confidenceBar,
		}
	}

	table := behaviors.NewTableBehavior(theme, columns, formatter).
		EmptyMessage("No suggestions available").
		PaginationPrefix("Suggestions").
		PageSize(10)

	table.SetItems(sortedSuggestions)

	m := &SuggestionReviewModal{
		table:       table,
		suggestions: sortedSuggestions,
		theme:       theme,
		accepted:    []burstfact.BurstSuggestion{},
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

	content.WriteString(m.table.Render())
	content.WriteString("\n\n")

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
		primitives.AcceptBadge(theme),
		primitives.RejectBadge(theme),
		primitives.NavigateBadge(theme),
		primitives.PageVimBadge(theme),
		primitives.CancelBadge(theme),
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
			selected := m.table.GetSelectedItem()
			if selected != nil {
				m.action = SuggestionActionAccept
				m.accepted = append(m.accepted, *selected)

				m.removeCurrentSuggestion()

				if len(m.suggestions) == 0 {
					m.visible = false
				}
			}
			return m, nil

		case "r":
			selected := m.table.GetSelectedItem()
			if selected != nil {
				m.action = SuggestionActionReject

				m.removeCurrentSuggestion()

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

	title := primitives.Title("Review Burst Suggestions", theme).Render()
	content := m.buildContent()
	footer := m.buildFooter()
	modalContent := lipgloss.JoinVertical(lipgloss.Left, title, "", content, "", footer)

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

	modalWidth := width - 12
	if modalWidth < 60 {
		modalWidth = 60
	}
	if modalWidth > 90 {
		modalWidth = 90
	}
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
func (m *SuggestionReviewModal) GetAcceptedSuggestions() []burstfact.BurstSuggestion {
	return m.accepted
}

// GetCurrentSuggestion returns the currently selected suggestion.
func (m *SuggestionReviewModal) GetCurrentSuggestion() *burstfact.BurstSuggestion {
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
func (m *SuggestionReviewModal) GetAllSuggestions() []burstfact.BurstSuggestion {
	return m.suggestions
}
