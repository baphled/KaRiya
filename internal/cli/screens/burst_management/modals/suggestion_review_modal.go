package modals

import (
	"fmt"
	"sort"
	"strconv"
	"strings"

	"github.com/baphled/kariya/internal/cli/behaviors"
	"github.com/baphled/kariya/internal/cli/themes"
	"github.com/baphled/kariya/internal/cli/uikit/containers"
	"github.com/baphled/kariya/internal/cli/uikit/primitives"
	themes2 "github.com/baphled/kariya/internal/cli/uikit/theme"
	"github.com/baphled/kariya/internal/service/career/burstfact"
	"github.com/baphled/kariya/internal/service/career/skillinference"
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
	// SuggestionActionViewEvents requests the intent to display events
	// associated with the currently selected suggestion. The modal hides
	// itself so the intent can show an events sub-modal. The user presses
	// Enter to trigger this action.
	SuggestionActionViewEvents SuggestionAction = "view_events"
)

// SuggestionReviewModal displays suggestions (burst or skill) for review in a table format.
// Suggestions are sorted by confidence (highest first).
// User can navigate through suggestions and accept/reject them.
//
// This modal is generic and supports both burst suggestions and skill suggestions
// via type switching (see NewSuggestionReviewModal for burst, NewSkillSuggestionModal for skills).
type SuggestionReviewModal struct {
	burstTable *behaviors.TableBehavior[burstfact.BurstSuggestion]
	skillTable *behaviors.TableBehavior[skillinference.SkillSuggestion]

	suggestionType string

	// Type-specific suggestion slices
	burstSuggestions []burstfact.BurstSuggestion
	skillSuggestions []skillinference.SkillSuggestion

	// Type-specific accepted lists
	acceptedBursts []burstfact.BurstSuggestion
	acceptedSkills []skillinference.SkillSuggestion

	// Shared state
	theme   themes.Theme
	action  SuggestionAction
	visible bool
	width   int
	height  int
}

// NewSuggestionReviewModal creates a new burst suggestion review modal.
//
// Expected:
//   - burstsuggestion must be valid.
//   - th must be a valid theme instance (can be nil).
//
// Returns:
//   - A fully initialized SuggestionReviewModal ready for use.
//
// Side effects:
//   - None.
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
			strconv.Itoa(len(s.EventIDs)),
			confidenceBar,
		}
	}

	table := behaviors.NewTableBehavior(theme, columns, formatter).
		EmptyMessage("No suggestions available").
		PaginationPrefix("Suggestions").
		PageSize(10)

	table.SetItems(sortedSuggestions)

	m := &SuggestionReviewModal{
		suggestionType:   "burst",
		burstTable:       table,
		burstSuggestions: sortedSuggestions,
		acceptedBursts:   []burstfact.BurstSuggestion{},
		theme:            theme,
		visible:          true,
		width:            80,
		height:           24,
	}

	return m
}

// NewSkillSuggestionModal creates a new skill suggestion review modal.
//
// Expected:
//   - skillsuggestion must be valid.
//   - th must be a valid theme instance (can be nil).
//
// Returns:
//   - A fully initialized SuggestionReviewModal ready for use.
//
// Side effects:
//   - None.
func NewSkillSuggestionModal(suggestions []skillinference.SkillSuggestion, theme themes.Theme) *SuggestionReviewModal {
	if theme == nil {
		theme = themes.NewDefaultTheme()
	}

	sortedSuggestions := make([]skillinference.SkillSuggestion, len(suggestions))
	copy(sortedSuggestions, suggestions)
	sort.Slice(sortedSuggestions, func(i, j int) bool {
		return sortedSuggestions[i].Confidence > sortedSuggestions[j].Confidence
	})

	// Define 4 columns for skills (Name, Category, Events, Confidence)
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
			strconv.Itoa(len(s.EventIDs)),
			confidenceBar,
		}
	}

	table := behaviors.NewTableBehavior(theme, columns, formatter).
		EmptyMessage("No skills detected").
		PaginationPrefix("Skills").
		PageSize(10)

	table.SetItems(sortedSuggestions)

	m := &SuggestionReviewModal{
		suggestionType:   "skill",
		skillTable:       table,
		skillSuggestions: sortedSuggestions,
		acceptedSkills:   []skillinference.SkillSuggestion{},
		theme:            theme,
		visible:          true,
		width:            80,
		height:           24,
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
	switch m.suggestionType {
	case "burst":
		return m.buildBurstContent()
	case "skill":
		return m.buildSkillContent()
	default:
		return "Unknown suggestion type"
	}
}

// buildBurstContent renders burst fact suggestions.
func (m *SuggestionReviewModal) buildBurstContent() string {
	if len(m.burstSuggestions) == 0 {
		return "No suggestions available"
	}

	theme := m.getTheme()
	var content strings.Builder

	content.WriteString(m.burstTable.Render())
	content.WriteString("\n\n")

	selected := m.burstTable.GetSelectedItem()
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

// buildSkillContent renders skill inference suggestions with usage contexts.
func (m *SuggestionReviewModal) buildSkillContent() string {
	if len(m.skillSuggestions) == 0 {
		return "No skills detected"
	}

	var content strings.Builder

	content.WriteString(m.skillTable.Render())
	content.WriteString("\n\n")

	if selected := m.skillTable.GetSelectedItem(); selected != nil {
		m.renderSelectedSkillDetail(&content, selected)
	}

	return content.String()
}

// renderSelectedSkillDetail writes the selected skill name, category badge, and usage contexts.
func (m *SuggestionReviewModal) renderSelectedSkillDetail(content *strings.Builder, selected *skillinference.SkillSuggestion) {
	theme := m.getTheme()

	content.WriteString(primitives.NewText("Selected:", theme).Bold().Render())
	content.WriteString(" " + selected.Name)

	categoryBadge := primitives.NewBadge(selected.Category, theme).
		Variant(primitives.BadgeTag).
		Render()
	content.WriteString(" " + categoryBadge + "\n")

	m.renderUsageContexts(content, selected.Contexts)
}

// renderUsageContexts writes up to 3 usage contexts with an overflow indicator.
func (m *SuggestionReviewModal) renderUsageContexts(content *strings.Builder, contexts []string) {
	if len(contexts) == 0 {
		return
	}

	theme := m.getTheme()
	content.WriteString("\n")
	content.WriteString(primitives.NewText("Usage Contexts:", theme).Bold().Render())
	content.WriteString("\n")

	maxContexts := 3
	if len(contexts) < maxContexts {
		maxContexts = len(contexts)
	}

	for i := range maxContexts {
		fmt.Fprintf(content, "  • %s\n", contexts[i])
	}

	if len(contexts) > 3 {
		remaining := len(contexts) - 3
		fmt.Fprintf(content, "  ... and %d more\n", remaining)
	}
}

// buildFooter builds the footer with help badges.
func (m *SuggestionReviewModal) buildFooter() string {
	theme := m.getTheme()

	badges := []*primitives.Badge{
		primitives.HelpKeyBadge("Enter", "View Events", theme),
		primitives.AcceptBadge(theme),
		primitives.RejectBadge(theme),
		primitives.NavigateBadge(theme),
		primitives.PageVimBadge(theme),
		primitives.CancelBadge(theme),
	}

	return primitives.RenderHelpFooter(theme, badges...)
}

// Init initializes the modal.
//
// Returns:
//   - A tea.Cmd value.
//
// Side effects:
//   - None.
func (m *SuggestionReviewModal) Init() tea.Cmd {
	return nil
}

// Update handles keyboard input.
//
// Expected:
//   - msg must be a valid tea.Msg type.
//
// Returns:
//   - tea.Model: the updated model.
//   - tea.Cmd: command to execute.
//
// Side effects:
//   - May update internal state based on key presses.
//   - May hide modal on Escape, Enter, or when no suggestions remain.
func (m *SuggestionReviewModal) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	if !m.IsVisible() {
		return m, nil
	}

	switch msg := msg.(type) {
	case tea.WindowSizeMsg:
		m.width = msg.Width
		m.height = msg.Height

		switch m.suggestionType {
		case "burst":
			m.burstTable.Dimensions(m.width-12, m.height-16)
		case "skill":
			m.skillTable.Dimensions(m.width-12, m.height-16)
		}
		return m, nil

	case tea.KeyMsg:
		switch msg.Type {
		case tea.KeyEsc:
			m.action = SuggestionActionCancel
			m.visible = false
			return m, nil

		case tea.KeyDown:
			m.handleNavigation("down")
			return m, nil

		case tea.KeyUp:
			m.handleNavigation("up")
			return m, nil

		case tea.KeyPgDown:
			m.handleNavigation("pgdn")
			return m, nil

		case tea.KeyPgUp:
			m.handleNavigation("pgup")
			return m, nil

		case tea.KeyEnter:
			m.action = SuggestionActionViewEvents
			m.visible = false
			return m, nil

		case tea.KeyRunes:
			switch msg.String() {
			case "j":
				m.handleNavigation("down")
				return m, nil

			case "k":
				m.handleNavigation("up")
				return m, nil

			case "n":
				m.handleNavigation("pgdn")
				return m, nil

			case "p":
				m.handleNavigation("pgup")
				return m, nil

			case "a":
				m.action = SuggestionActionAccept
				m.handleAccept()

				if !m.HasSuggestions() {
					m.visible = false
				}
				return m, nil

			case "r":
				m.action = SuggestionActionReject
				m.removeCurrentSuggestion()

				if !m.HasSuggestions() {
					m.visible = false
				}
				return m, nil
			}
		}
	}

	return m, nil
}

// handleNavigation delegates navigation to the correct table.
func (m *SuggestionReviewModal) handleNavigation(direction string) {
	switch m.suggestionType {
	case "burst":
		m.burstTable.HandleNavigation(direction)
	case "skill":
		m.skillTable.HandleNavigation(direction)
	}
}

// handleAccept adds current suggestion to accepted list and removes it.
func (m *SuggestionReviewModal) handleAccept() {
	switch m.suggestionType {
	case "burst":
		selected := m.burstTable.GetSelectedItem()
		if selected != nil {
			m.acceptedBursts = append(m.acceptedBursts, *selected)
			m.removeCurrentSuggestion()
		}
	case "skill":
		selected := m.skillTable.GetSelectedItem()
		if selected != nil {
			m.acceptedSkills = append(m.acceptedSkills, *selected)
			m.removeCurrentSuggestion()
		}
	}
}

// removeCurrentSuggestion removes the currently selected suggestion from the list.
func (m *SuggestionReviewModal) removeCurrentSuggestion() {
	switch m.suggestionType {
	case "burst":
		idx := m.burstTable.GetSelectedIndex()
		if idx >= 0 && idx < len(m.burstSuggestions) {
			m.burstSuggestions = append(m.burstSuggestions[:idx], m.burstSuggestions[idx+1:]...)
			m.burstTable.SetItems(m.burstSuggestions)

			if idx >= len(m.burstSuggestions) && len(m.burstSuggestions) > 0 {
				idx = len(m.burstSuggestions) - 1
			}
			if len(m.burstSuggestions) > 0 {
				m.burstTable.SetSelectedIndex(idx)
			}
		}
	case "skill":
		idx := m.skillTable.GetSelectedIndex()
		if idx >= 0 && idx < len(m.skillSuggestions) {
			m.skillSuggestions = append(m.skillSuggestions[:idx], m.skillSuggestions[idx+1:]...)
			m.skillTable.SetItems(m.skillSuggestions)

			if idx >= len(m.skillSuggestions) && len(m.skillSuggestions) > 0 {
				idx = len(m.skillSuggestions) - 1
			}
			if len(m.skillSuggestions) > 0 {
				m.skillTable.SetSelectedIndex(idx)
			}
		}
	}
}

// View renders the modal.
//
// Returns:
//   - A string value.
//
// Side effects:
//   - None.
func (m *SuggestionReviewModal) View() string {
	if !m.visible {
		return ""
	}

	theme := m.getTheme()

	var titleText string
	switch m.suggestionType {
	case "burst":
		titleText = "Review Burst Suggestions"
	case "skill":
		titleText = "Review Skill Suggestions"
	default:
		titleText = "Review Suggestions"
	}

	title := primitives.Title(titleText, theme).Render()
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
//
// Returns:
//   - A bool value.
//
// Side effects:
//   - None.
func (m *SuggestionReviewModal) IsVisible() bool {
	return m.visible
}

// Show makes the modal visible.
//
// Side effects:
//   - None.
func (m *SuggestionReviewModal) Show() {
	m.visible = true
	m.action = ""
}

// Hide hides the modal.
//
// Side effects:
//   - None.
func (m *SuggestionReviewModal) Hide() {
	m.visible = false
}

// SetDimensions sets the terminal dimensions.
//
// Expected:
//   - int must be valid.
//
// Side effects:
//   - None.
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

	switch m.suggestionType {
	case "burst":
		m.burstTable.Dimensions(contentWidth, height-16)
	case "skill":
		m.skillTable.Dimensions(contentWidth, height-16)
	}
}

// GetAction returns the last action taken.
//
// Returns:
//   - A SuggestionAction value.
//
// Side effects:
//   - None.
func (m *SuggestionReviewModal) GetAction() SuggestionAction {
	return m.action
}

// ClearAction clears the current action.
//
// Side effects:
//   - None.
func (m *SuggestionReviewModal) ClearAction() {
	m.action = ""
}

// GetAcceptedSuggestions returns all accepted burst suggestions.
//
// Returns:
//   - A []burstfact.BurstSuggestion value.
//
// Side effects:
//   - None.
func (m *SuggestionReviewModal) GetAcceptedSuggestions() []burstfact.BurstSuggestion {
	return m.acceptedBursts
}

// GetAcceptedSkills returns all accepted skill suggestions.
//
// Returns:
//   - A []skillinference.SkillSuggestion value.
//
// Side effects:
//   - None.
func (m *SuggestionReviewModal) GetAcceptedSkills() []skillinference.SkillSuggestion {
	return m.acceptedSkills
}

// GetCurrentSuggestion returns the currently selected burst suggestion.
//
// Returns:
//   - A fully initialized burstfact.BurstSuggestion ready for use.
//
// Side effects:
//   - None.
func (m *SuggestionReviewModal) GetCurrentSuggestion() *burstfact.BurstSuggestion {
	return m.burstTable.GetSelectedItem()
}

// GetCurrentSkill returns the currently selected skill suggestion.
//
// Returns:
//   - A fully initialized skillinference.SkillSuggestion ready for use.
//
// Side effects:
//   - None.
func (m *SuggestionReviewModal) GetCurrentSkill() *skillinference.SkillSuggestion {
	return m.skillTable.GetSelectedItem()
}

// HasSuggestions returns whether there are any suggestions left.
//
// Returns:
//   - A bool value.
//
// Side effects:
//   - None.
func (m *SuggestionReviewModal) HasSuggestions() bool {
	switch m.suggestionType {
	case "burst":
		return len(m.burstSuggestions) > 0
	case "skill":
		return len(m.skillSuggestions) > 0
	default:
		return false
	}
}

// GetSuggestionsCount returns the number of remaining suggestions.
//
// Returns:
//   - A int value.
//
// Side effects:
//   - None.
func (m *SuggestionReviewModal) GetSuggestionsCount() int {
	switch m.suggestionType {
	case "burst":
		return len(m.burstSuggestions)
	case "skill":
		return len(m.skillSuggestions)
	default:
		return 0
	}
}

// GetAllSuggestions returns all burst suggestions in their current order (sorted by confidence).
//
// Returns:
//   - A []burstfact.BurstSuggestion value.
//
// Side effects:
//   - None.
func (m *SuggestionReviewModal) GetAllSuggestions() []burstfact.BurstSuggestion {
	return m.burstSuggestions
}

// GetAllSkills returns all skill suggestions in their current order (sorted by confidence).
//
// Returns:
//   - A []skillinference.SkillSuggestion value.
//
// Side effects:
//   - None.
func (m *SuggestionReviewModal) GetAllSkills() []skillinference.SkillSuggestion {
	return m.skillSuggestions
}
