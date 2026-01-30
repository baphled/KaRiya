package modals

import (
	"fmt"
	"sort"

	"github.com/baphled/kariya/internal/cli/behaviors"
	"github.com/baphled/kariya/internal/cli/themes"
	"github.com/baphled/kariya/internal/cli/uikit/containers"
	"github.com/baphled/kariya/internal/cli/uikit/primitives"
	themes2 "github.com/baphled/kariya/internal/cli/uikit/theme"
	"github.com/baphled/kariya/internal/service/career/skillinference"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
)

// SkillSuggestionAction represents an action taken on a skill suggestion.
type SkillSuggestionAction string

const (
	// SkillSuggestionActionAccept marks the suggestion as approved.
	SkillSuggestionActionAccept SkillSuggestionAction = "accept"
	// SkillSuggestionActionReject discards the suggestion.
	SkillSuggestionActionReject SkillSuggestionAction = "reject"
	// SkillSuggestionActionAcceptAll accepts all suggestions.
	SkillSuggestionActionAcceptAll SkillSuggestionAction = "accept_all"
	// SkillSuggestionActionCancel aborts the review session.
	SkillSuggestionActionCancel SkillSuggestionAction = "cancel"
)

// SkillSuggestionModal displays skill suggestions for review in a table format.
// Suggestions are sorted by confidence (highest first).
type SkillSuggestionModal struct {
	table       *behaviors.TableBehavior[skillinference.SkillSuggestion]
	suggestions []skillinference.SkillSuggestion
	theme       themes.Theme
	action      SkillSuggestionAction
	accepted    []skillinference.SkillSuggestion
	visible     bool
	width       int
	height      int
}

// NewSkillSuggestionModal creates a new skill suggestion modal.
// Suggestions are automatically sorted by confidence (highest first).
func NewSkillSuggestionModal(suggestions []skillinference.SkillSuggestion, theme themes.Theme) *SkillSuggestionModal {
	if theme == nil {
		theme = themes.NewDefaultTheme()
	}

	sortedSuggestions := make([]skillinference.SkillSuggestion, len(suggestions))
	copy(sortedSuggestions, suggestions)
	sort.Slice(sortedSuggestions, func(i, j int) bool {
		return sortedSuggestions[i].Confidence > sortedSuggestions[j].Confidence
	})

	columns := []behaviors.ColumnDef{
		{Title: "Skill", Width: 25},
		{Title: "Category", Width: 12},
		{Title: "Events", Width: 8},
		{Title: "Confidence", Width: 24},
	}

	formatter := func(s skillinference.SkillSuggestion, _ int) []string {
		confidenceBar := primitives.CompactBar(s.Confidence, 15, nil).
			ShowPercentage(true).
			Render()

		return []string{
			s.Name,
			s.Category,
			fmt.Sprintf("%d", len(s.EventIDs)),
			confidenceBar,
		}
	}

	table := behaviors.NewTableBehavior(theme, columns, formatter).
		EmptyMessage("No skills detected").
		PaginationPrefix("Skills").
		PageSize(10)

	table.SetItems(sortedSuggestions)

	return &SkillSuggestionModal{
		table:       table,
		suggestions: sortedSuggestions,
		theme:       theme,
		accepted:    []skillinference.SkillSuggestion{},
		visible:     false,
		width:       80,
		height:      24,
	}
}

// Show makes the modal visible.
func (m *SkillSuggestionModal) Show() {
	m.visible = true
}

// Hide makes the modal hidden.
func (m *SkillSuggestionModal) Hide() {
	m.visible = false
}

// IsVisible returns whether the modal is currently visible.
func (m *SkillSuggestionModal) IsVisible() bool {
	return m.visible
}

// Update handles input events.
func (m *SkillSuggestionModal) Update(msg tea.Msg) (tea.Cmd, interface{}) {
	if !m.visible {
		return nil, nil
	}

	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch msg.Type {
		case tea.KeyUp:
			m.table.HandleNavigation("up")
			return nil, nil
		case tea.KeyDown:
			m.table.HandleNavigation("down")
			return nil, nil
		case tea.KeyPgUp:
			m.table.HandleNavigation("pgup")
			return nil, nil
		case tea.KeyPgDown:
			m.table.HandleNavigation("pgdn")
			return nil, nil
		case tea.KeyHome:
			m.table.HandleNavigation("home")
			return nil, nil
		case tea.KeyEnd:
			m.table.HandleNavigation("end")
			return nil, nil
		case tea.KeyEsc:
			m.action = SkillSuggestionActionCancel
			m.Hide()
			return nil, m.action
		case tea.KeyRunes:
			switch msg.String() {
			case "k":
				m.table.HandleNavigation("up")
				return nil, nil
			case "j":
				m.table.HandleNavigation("down")
				return nil, nil
			case "a":
				m.acceptCurrent()
				return nil, SkillSuggestionActionAccept
			case "r":
				m.rejectCurrent()
				return nil, SkillSuggestionActionReject
			case "A":
				m.acceptAll()
				m.action = SkillSuggestionActionAcceptAll
				m.Hide()
				return nil, m.action
			}
		}

	case tea.WindowSizeMsg:
		m.width = msg.Width
		m.height = msg.Height
		m.table.Dimensions(msg.Width-12, msg.Height-16)
	}

	return nil, nil
}

// acceptCurrent marks the current suggestion as accepted and removes it from the table.
func (m *SkillSuggestionModal) acceptCurrent() {
	if len(m.suggestions) == 0 {
		return
	}

	idx := m.table.GetSelectedIndex()
	if idx < 0 || idx >= len(m.suggestions) {
		return
	}

	currentItem := m.suggestions[idx]
	m.accepted = append(m.accepted, currentItem)

	m.removeCurrent()
}

// rejectCurrent removes the current suggestion without accepting it.
func (m *SkillSuggestionModal) rejectCurrent() {
	if len(m.suggestions) == 0 {
		return
	}

	m.removeCurrent()
}

// removeCurrent removes the current item from the suggestions list.
func (m *SkillSuggestionModal) removeCurrent() {
	idx := m.table.GetSelectedIndex()
	if idx < 0 || idx >= len(m.suggestions) {
		return
	}

	// Remove from suggestions
	m.suggestions = append(m.suggestions[:idx], m.suggestions[idx+1:]...)

	// Update table with new data
	m.table.SetItems(m.suggestions)

	// If list is now empty, close modal
	if len(m.suggestions) == 0 {
		m.Hide()
	}
}

// acceptAll marks all remaining suggestions as accepted.
func (m *SkillSuggestionModal) acceptAll() {
	m.accepted = append(m.accepted, m.suggestions...)
	m.suggestions = []skillinference.SkillSuggestion{}
}

// GetAcceptedSuggestions returns all accepted suggestions.
func (m *SkillSuggestionModal) GetAcceptedSuggestions() []skillinference.SkillSuggestion {
	return m.accepted
}

// getTheme returns the theme or default if nil.
func (m *SkillSuggestionModal) getTheme() themes.Theme {
	if m.theme != nil {
		return m.theme
	}
	return themes2.Default()
}

// View renders the modal.
func (m *SkillSuggestionModal) View() string {
	if !m.visible {
		return ""
	}

	if len(m.suggestions) == 0 {
		return m.renderEmpty()
	}

	theme := m.getTheme()

	// Header
	title := primitives.Title("Skill Suggestions", theme).Render()
	subtitle := primitives.Body(fmt.Sprintf("Review %d detected skills", len(m.suggestions)+len(m.accepted)), theme).
		Render()

	// Table
	tableView := m.table.Render()

	// Footer - help text
	helpText := m.buildFooter()

	// Status
	status := primitives.Body(fmt.Sprintf("Accepted: %d | Remaining: %d",
		len(m.accepted), len(m.suggestions)), theme).
		Render()

	// Combine sections
	content := lipgloss.JoinVertical(
		lipgloss.Left,
		title,
		subtitle,
		"",
		tableView,
		"",
		status,
		"",
		helpText,
	)

	// Wrap in box
	return containers.NewBox(theme).
		Content(content).
		Background(theme.BackgroundColor()).
		Render()
}

// buildFooter builds the footer with help badges.
func (m *SkillSuggestionModal) buildFooter() string {
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

// renderEmpty renders the modal when all suggestions have been processed.
func (m *SkillSuggestionModal) renderEmpty() string {
	theme := m.getTheme()

	title := primitives.Title("All Suggestions Reviewed", theme).Render()
	message := primitives.Body(fmt.Sprintf("Accepted %d skills", len(m.accepted)), theme).Render()

	content := lipgloss.JoinVertical(lipgloss.Left, title, "", message)

	return containers.NewBox(theme).
		Content(content).
		Background(theme.BackgroundColor()).
		Render()
}
