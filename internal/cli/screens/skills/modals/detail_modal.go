package modals

import (
	"fmt"
	"strconv"
	"strings"
	"time"

	"github.com/baphled/kariya/internal/cli/themes"
	"github.com/baphled/kariya/internal/cli/uikit/containers"
	"github.com/baphled/kariya/internal/cli/uikit/primitives"
	"github.com/baphled/kariya/internal/domain/career"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
)

// DetailModal displays skill details in a modal overlay.
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
//	modal := NewDetailModal(skill, theme, eventCount, lastUsed)
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
type DetailModal struct {
	skill      *career.Skill
	theme      themes.Theme
	visible    bool
	width      int
	height     int
	action     string // "events", "edit", "delete", or "" for close
	eventCount int
	lastUsed   *time.Time
}

// NewDetailModal creates a new skill detail modal.
//
// Expected:
//   - skill must be valid.
//   - th must be a valid theme instance (can be nil).
//   - int must be valid.
//   - time must be valid.
//
// Returns:
//   - A fully initialized DetailModal ready for use.
//
// Side effects:
//   - None.
func NewDetailModal(skill *career.Skill, theme themes.Theme, eventCount int, lastUsed *time.Time) *DetailModal {
	return &DetailModal{
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
//
// Returns:
//   - A tea.Cmd value.
//
// Side effects:
//   - None.
func (m *DetailModal) Init() tea.Cmd {
	return nil
}

// Update handles keyboard input and window sizing.
//
// Expected:
//   - msg must be a valid tea.Msg type.
//
// Returns:
//   - tea.Model: the updated model.
//   - tea.Cmd: command to execute.
//
// Side effects:
//   - May set action and hide modal.
func (m *DetailModal) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
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
//
// Returns:
//   - A string value.
//
// Side effects:
//   - None.
func (m *DetailModal) View() string {
	if !m.visible {
		return ""
	}

	// Nil theme guard
	theme := m.theme
	if theme == nil {
		theme = themes.NewDefaultTheme()
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
	content := m.renderSkillDetails(theme)

	// Build footer with actions using UIKit
	footerParts := []string{"e: Edit", "ctrl+e: Events", "d: Delete", "Enter/Esc: Close"}
	footer := primitives.Muted(strings.Join(footerParts, " | "), theme).Render()

	// Build modal content
	modalContent := lipgloss.JoinVertical(lipgloss.Left, content, "", footer)

	// Wrap in styled box with solid background using UIKit
	return containers.NewBox(theme).
		Content(modalContent).
		Width(modalWidth).
		Padding(2).
		Background(theme.BackgroundColor()).
		Render()
}

// renderSkillDetails renders the skill information as a card using UIKit KeyValue.
func (m *DetailModal) renderSkillDetails(th themes.Theme) string {
	skill := m.skill

	// Use UIKit KeyValue for consistent label-value layout
	// theme.Theme is an alias for themes.Theme, so direct pass works
	kv := primitives.NewKeyValue(th).LabelWidth(15)

	// Name
	kv.Add("Name:", skill.Name)

	// Category
	category := skill.Category
	if category == "" {
		category = "-"
	}
	kv.Add("Category:", category)

	// Level
	level := skill.Level
	if level == "" {
		level = "-"
	}
	kv.Add("Level:", level)

	// Years Used
	if skill.YearsUsed != nil {
		yearText := fmt.Sprintf("%d year", *skill.YearsUsed)
		if *skill.YearsUsed != 1 {
			yearText += "s"
		}
		kv.Add("Years Used:", yearText)
	}

	// Event Count
	kv.Add("Event Count:", strconv.Itoa(m.eventCount))

	// Last Used (if available)
	if m.lastUsed != nil {
		kv.Add("Last Used:", m.lastUsed.Format("2006-01-02"))
	}

	// Timestamps (muted)
	kv.AddBlank()
	kv.AddMuted("Created:", skill.CreatedAt.Format("2006-01-02 15:04"))
	kv.AddMuted("Updated:", skill.UpdatedAt.Format("2006-01-02 15:04"))

	return kv.Render()
}

// SetDimensions updates the modal's available dimensions.
//
// Expected:
//   - int must be valid.
//
// Side effects:
//   - None.
func (m *DetailModal) SetDimensions(width, height int) {
	m.width = width
	m.height = height
}

// Show makes the modal visible.
//
// Side effects:
//   - None.
func (m *DetailModal) Show() {
	m.visible = true
	m.action = ""
}

// Hide hides the modal.
//
// Side effects:
//   - None.
func (m *DetailModal) Hide() {
	m.visible = false
}

// IsVisible returns whether the modal is currently visible.
//
// Returns:
//   - A bool value.
//
// Side effects:
//   - None.
func (m *DetailModal) IsVisible() bool {
	return m.visible
}

// GetAction returns the action selected by the user.
//
// Returns:
//   - A string value.
//
// Side effects:
//   - None.
func (m *DetailModal) GetAction() string {
	return m.action
}

// SetSkill updates the skill being displayed (useful for reusing the modal).
//
// Expected:
//   - skill must be valid.
//   - int must be valid.
//   - time must be valid.
//
// Side effects:
//   - None.
func (m *DetailModal) SetSkill(skill *career.Skill, eventCount int, lastUsed *time.Time) {
	m.skill = skill
	m.eventCount = eventCount
	m.lastUsed = lastUsed
	m.action = ""
}
