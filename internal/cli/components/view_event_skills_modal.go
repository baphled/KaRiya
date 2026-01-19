package components

import (
	"fmt"
	"strings"

	"github.com/baphled/kariya/internal/cli/styles"
	"github.com/baphled/kariya/internal/cli/themes"
	"github.com/baphled/kariya/internal/domain/career"
	"github.com/charmbracelet/bubbles/viewport"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
)

// ViewEventSkillsModal displays skills associated with an event in a modal overlay.
// This modal is triggered from the event detail view to show the full list of skills.
//
// Features:
// - Read-only display of skills with name, category, level, and years
// - Scrollable content for many skills
// - Close with Escape, Enter, backspace, or 'q'
// - Solid background to prevent transparency issues
//
// Usage:
//
//	modal := NewViewEventSkillsModal(eventID, skills, theme)
//	modal.SetDimensions(width, height)
//	modal.Show()
//
//	// In Update:
//	model, cmd := modal.Update(msg)
//
//	// In View:
//	if modal.IsVisible() {
//	    return renderModalOverlay(modal, background)
//	}
type ViewEventSkillsModal struct {
	eventID    string
	skills     []*career.Skill
	theme      themes.Theme
	visible    bool
	width      int
	height     int
	viewport   viewport.Model
	ready      bool
	hasContent bool
}

// NewViewEventSkillsModal creates a new skills modal for an event.
func NewViewEventSkillsModal(eventID string, skills []*career.Skill, theme themes.Theme) *ViewEventSkillsModal {
	return &ViewEventSkillsModal{
		eventID: eventID,
		skills:  skills,
		theme:   theme,
		visible: false,
		width:   80,
		height:  24,
	}
}

// Init initializes the modal (implements tea.Model for bubbletea-overlay).
func (m *ViewEventSkillsModal) Init() tea.Cmd {
	return nil
}

// Update handles keyboard input and window sizing.
func (m *ViewEventSkillsModal) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	if !m.visible {
		return m, nil
	}

	var cmd tea.Cmd

	switch msg := msg.(type) {
	case tea.WindowSizeMsg:
		m.width = msg.Width
		m.height = msg.Height
		m.ready = false // Force viewport recreation on resize
		return m, nil

	case tea.KeyMsg:
		switch msg.String() {
		case "esc", "backspace", "enter", "q":
			// Close modal
			m.Hide()
			return m, nil

		// Scrolling keys - pass to viewport if content is scrollable
		case "up", "k", "down", "j", "pgup", "pgdown", "ctrl+u", "ctrl+d":
			if m.ready && m.hasContent {
				m.viewport, cmd = m.viewport.Update(msg)
				return m, cmd
			}
		}
	}

	return m, nil
}

// View renders the modal content with solid background and scrolling.
func (m *ViewEventSkillsModal) View() string {
	if !m.visible {
		return ""
	}

	// Calculate modal dimensions
	maxModalHeight := 30
	terminalMaxHeight := int(float64(m.height) * 0.7)
	if terminalMaxHeight < maxModalHeight {
		maxModalHeight = terminalMaxHeight
	}
	if maxModalHeight < 10 {
		maxModalHeight = 10 // Minimum usable height
	}

	modalWidth := m.width - 12 // Leave margins
	if modalWidth > 80 {
		modalWidth = 80 // Max width for readability
	}
	if modalWidth < 40 {
		modalWidth = 40 // Minimum usable width
	}

	// Render skills content
	content := m.renderSkillsContent(modalWidth - 6) // Account for padding and borders
	contentLines := strings.Split(content, "\n")
	contentHeight := len(contentLines)

	// Calculate viewport height (modal height - borders - padding - title - footer)
	viewportHeight := maxModalHeight - 8 // Account for border (2), padding (2), title (2), footer (2)
	if viewportHeight < 5 {
		viewportHeight = 5
	}

	// Initialize viewport if needed
	if !m.ready {
		m.viewport = viewport.New(modalWidth-4, viewportHeight)
		m.viewport.SetContent(content)
		m.hasContent = contentHeight > viewportHeight
		m.ready = true
	}

	// Build title
	titleStyle := lipgloss.NewStyle().
		Bold(true).
		Foreground(styles.ColorTextPrimary)
	title := titleStyle.Render(fmt.Sprintf("Skills (%d)", len(m.skills)))

	// Build footer with scroll indicator
	scrollHint := "Enter/Esc: Close"
	if m.hasContent {
		percentScrolled := int(m.viewport.ScrollPercent() * 100)
		scrollHint = lipgloss.NewStyle().
			Foreground(styles.ColorTextSecondary).
			Render(fmt.Sprintf("↑↓/j/k: Scroll | Enter/Esc: Close [%d%%]", percentScrolled))
	} else {
		scrollHint = lipgloss.NewStyle().
			Foreground(styles.ColorTextSecondary).
			Render(scrollHint)
	}

	// Build modal content
	modalContent := lipgloss.JoinVertical(lipgloss.Left, title, "", m.viewport.View(), "", scrollHint)

	// Wrap in styled box with solid background
	return lipgloss.NewStyle().
		Border(lipgloss.RoundedBorder()).
		BorderForeground(styles.ColorBorder).
		Background(styles.ColorBackground).
		Padding(1, 2).
		Width(modalWidth).
		MaxHeight(maxModalHeight).
		Render(modalContent)
}

// renderSkillsContent formats the skills list for display.
func (m *ViewEventSkillsModal) renderSkillsContent(width int) string {
	if len(m.skills) == 0 {
		return lipgloss.NewStyle().
			Foreground(styles.ColorTextSecondary).
			Italic(true).
			Render("No skills associated with this event.")
	}

	var content strings.Builder

	// Style definitions
	nameStyle := lipgloss.NewStyle().
		Bold(true).
		Foreground(styles.ColorTextPrimary)
	categoryStyle := lipgloss.NewStyle().
		Foreground(styles.ColorTextSecondary)
	levelStyle := lipgloss.NewStyle().
		Foreground(styles.ColorSuccess)
	yearsStyle := lipgloss.NewStyle().
		Foreground(styles.ColorInfo)
	separatorStyle := lipgloss.NewStyle().
		Foreground(styles.ColorBorder)

	for i, skill := range m.skills {
		if skill == nil {
			continue
		}

		// Skill name with bullet
		content.WriteString(nameStyle.Render(fmt.Sprintf("• %s", skill.Name)))
		content.WriteString("\n")

		// Details line: category | level | years
		var details []string

		if skill.Category != "" {
			details = append(details, categoryStyle.Render(skill.Category))
		}

		if skill.Level != "" {
			details = append(details, levelStyle.Render(skill.Level))
		}

		if skill.YearsUsed != nil && *skill.YearsUsed > 0 {
			yearText := "year"
			if *skill.YearsUsed > 1 {
				yearText = "years"
			}
			details = append(details, yearsStyle.Render(fmt.Sprintf("%d %s", *skill.YearsUsed, yearText)))
		}

		if len(details) > 0 {
			content.WriteString("  ")
			content.WriteString(strings.Join(details, separatorStyle.Render(" | ")))
			content.WriteString("\n")
		}

		// Add spacing between skills (except after last)
		if i < len(m.skills)-1 {
			content.WriteString("\n")
		}
	}

	return content.String()
}

// SetDimensions updates the modal's available dimensions.
func (m *ViewEventSkillsModal) SetDimensions(width, height int) {
	m.width = width
	m.height = height
	m.ready = false // Reset viewport when dimensions change
}

// Show makes the modal visible.
func (m *ViewEventSkillsModal) Show() {
	m.visible = true
	m.ready = false // Reset viewport when showing
}

// Hide hides the modal.
func (m *ViewEventSkillsModal) Hide() {
	m.visible = false
}

// IsVisible returns whether the modal is currently visible.
func (m *ViewEventSkillsModal) IsVisible() bool {
	return m.visible
}

// GetEventID returns the event ID this modal is showing skills for.
func (m *ViewEventSkillsModal) GetEventID() string {
	return m.eventID
}

// SetSkills updates the skills being displayed.
func (m *ViewEventSkillsModal) SetSkills(skills []*career.Skill) {
	m.skills = skills
	m.ready = false // Reset viewport when skills change
}
