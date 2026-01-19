package components

import (
	"fmt"
	"strings"
	"time"

	"github.com/baphled/kariya/internal/cli/styles"
	"github.com/baphled/kariya/internal/cli/themes"
	"github.com/baphled/kariya/internal/domain/career"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
)

// ViewSkillDetailModal displays skill details in a modal overlay.
// Unlike the full-screen detail view, this shows the skill info over the skills list.
//
// Features:
// - Read-only display of skill details (name, category, level, years, event count)
// - Actions: Edit (e), View Events (ctrl+e), Delete (d)
// - Close with Escape, backspace, Enter, or 'q'
// - Solid background to prevent transparency issues
// - Uses bubbletea-overlay for compositing
//
// Usage:
//
//	modal := NewViewSkillDetailModal(skill, theme, eventCount, lastUsed)
//	modal.SetDimensions(width, height)
//	modal.Show()
//
//	// In Update:
//	cmd := modal.Update(msg)
//	if !modal.IsVisible() {
//	    action := modal.GetAction() // "edit", "events", "delete", or ""
//	}
//
//	// In View:
//	if modal.IsVisible() {
//	    return renderModalOverlay(modal, background)
//	}
type ViewSkillDetailModal struct {
	skill      *career.Skill
	theme      themes.Theme
	visible    bool
	width      int
	height     int
	action     string // "events", "edit", "delete", or "" for close
	eventCount int
	lastUsed   *time.Time
}

// NewViewSkillDetailModal creates a new skill detail modal.
func NewViewSkillDetailModal(skill *career.Skill, theme themes.Theme, eventCount int, lastUsed *time.Time) *ViewSkillDetailModal {
	return &ViewSkillDetailModal{
		skill:      skill,
		theme:      theme,
		visible:    false,
		width:      80,
		height:     24,
		action:     "",
		eventCount: eventCount,
		lastUsed:   lastUsed,
	}
}

// Init initializes the modal (implements tea.Model for bubbletea-overlay).
func (m *ViewSkillDetailModal) Init() tea.Cmd {
	return nil
}

// Update handles keyboard input and window sizing.
func (m *ViewSkillDetailModal) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	if !m.visible {
		return m, nil
	}

	switch msg := msg.(type) {
	case tea.WindowSizeMsg:
		m.width = msg.Width
		m.height = msg.Height
		return m, nil

	case tea.KeyMsg:
		switch msg.Type {
		case tea.KeyEsc, tea.KeyBackspace, tea.KeyEnter:
			// Close modal without action
			m.action = ""
			m.Hide()
			return m, nil

		case tea.KeyCtrlE:
			// View events using this skill
			m.action = "events"
			m.Hide()
			return m, nil
		}

		switch msg.String() {
		case "q":
			// Close modal
			m.action = ""
			m.Hide()
			return m, nil

		case "e":
			// Edit skill (standard edit keybinding)
			m.action = "edit"
			m.Hide()
			return m, nil

		case "d":
			// Delete skill
			m.action = "delete"
			m.Hide()
			return m, nil

		}
	}

	return m, nil
}

// View renders the modal content with solid background.
func (m *ViewSkillDetailModal) View() string {
	if !m.visible {
		return ""
	}

	// Calculate modal dimensions
	modalWidth := m.width - 12
	if modalWidth > 80 {
		modalWidth = 80
	}
	if modalWidth < 40 {
		modalWidth = 40
	}

	// Render skill details directly (no viewport needed for simple content)
	content := m.renderSkillDetails()

	// Build footer with actions
	footerParts := []string{"e: Edit", "ctrl+e: Events", "d: Delete", "Enter/Esc: Close"}

	footer := lipgloss.NewStyle().
		Foreground(styles.ColorTextSecondary).
		Render(strings.Join(footerParts, " | "))

	// Build modal content
	modalContent := lipgloss.JoinVertical(lipgloss.Left, content, "", footer)

	// Wrap in styled box with solid background
	return lipgloss.NewStyle().
		Border(lipgloss.RoundedBorder()).
		BorderForeground(styles.ColorBorder).
		Background(styles.ColorBackground).
		Padding(1, 2).
		Width(modalWidth).
		Render(modalContent)
}

// renderSkillDetails renders the skill information as a card.
func (m *ViewSkillDetailModal) renderSkillDetails() string {
	skill := m.skill

	labelStyle := lipgloss.NewStyle().
		Foreground(m.theme.MutedColor()).
		Width(15)

	valueStyle := lipgloss.NewStyle().
		Foreground(m.theme.PrimaryColor()).
		Bold(true)

	var lines []string

	// Name
	lines = append(lines, labelStyle.Render("Name:")+valueStyle.Render(skill.Name))

	// Category
	category := skill.Category
	if category == "" {
		category = "-"
	}
	lines = append(lines, labelStyle.Render("Category:")+valueStyle.Render(category))

	// Level
	level := skill.Level
	if level == "" {
		level = "-"
	}
	lines = append(lines, labelStyle.Render("Level:")+valueStyle.Render(level))

	// Years Used
	if skill.YearsUsed != nil {
		yearText := fmt.Sprintf("%d year", *skill.YearsUsed)
		if *skill.YearsUsed != 1 {
			yearText += "s"
		}
		lines = append(lines, labelStyle.Render("Years Used:")+valueStyle.Render(yearText))
	}

	// Event Count
	lines = append(lines, labelStyle.Render("Event Count:")+valueStyle.Render(fmt.Sprintf("%d", m.eventCount)))

	// Last Used (if available)
	if m.lastUsed != nil {
		lines = append(lines, labelStyle.Render("Last Used:")+valueStyle.Render(m.lastUsed.Format("2006-01-02")))
	}

	// Timestamps
	lines = append(lines, "")
	lines = append(lines, labelStyle.Render("Created:")+lipgloss.NewStyle().
		Foreground(m.theme.MutedColor()).
		Render(skill.CreatedAt.Format("2006-01-02 15:04")))
	lines = append(lines, labelStyle.Render("Updated:")+lipgloss.NewStyle().
		Foreground(m.theme.MutedColor()).
		Render(skill.UpdatedAt.Format("2006-01-02 15:04")))

	return lipgloss.JoinVertical(lipgloss.Left, lines...)
}

// SetDimensions updates the modal's available dimensions.
func (m *ViewSkillDetailModal) SetDimensions(width, height int) {
	m.width = width
	m.height = height
}

// Show makes the modal visible.
func (m *ViewSkillDetailModal) Show() {
	m.visible = true
	m.action = ""
}

// Hide hides the modal.
func (m *ViewSkillDetailModal) Hide() {
	m.visible = false
}

// IsVisible returns whether the modal is currently visible.
func (m *ViewSkillDetailModal) IsVisible() bool {
	return m.visible
}

// GetAction returns the action selected by the user.
// Returns "events", "edit", "delete", or "" for simple close.
func (m *ViewSkillDetailModal) GetAction() string {
	return m.action
}

// SetSkill updates the skill being displayed (useful for reusing the modal).
func (m *ViewSkillDetailModal) SetSkill(skill *career.Skill, eventCount int, lastUsed *time.Time) {
	m.skill = skill
	m.eventCount = eventCount
	m.lastUsed = lastUsed
	m.action = ""
}
