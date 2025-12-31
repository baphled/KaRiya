package models

import (
	"strings"

	"github.com/baphled/kariya/internal/cli/components"
	"github.com/charmbracelet/lipgloss"
	"github.com/baphled/kariya/internal/cli/styles"
	tea "github.com/charmbracelet/bubbletea"
)

// HelpModel represents the help screen
type HelpModel struct {
	currentSection int
	sections       []HelpSection
	width          int
	height         int
	closed         bool
	helpFooter       components.HelpFooterModel
}

// HelpSection contains help content for a topic
type HelpSection struct {
	title   string
	content string
}

// NewHelpModel creates a new help model
func NewHelpModel() *HelpModel {
	sections := []HelpSection{
		{
			title: "Overview",
			content: `KaRiya is a career journaling tool that helps you:
  • Capture career events and milestones
  • Organize events with tags and categories
  • Search and filter your career history
  • Export your career data

Press 'right' or 'space' to continue through the help sections.
Press 'q' or 'esc' to exit help.`,
		},
		{
			title: "Event Capture",
			content: `Capture your career events in three ways:

1. Timeline Journaling (Recent events - last 30 days)
   Use this to log events as they happen
   Best for: Real-time journaling

2. CV Backfill (Historical events - any past date)
   Use this to import events from existing CV
   Best for: Building your career history

3. Manual Entry (Flexible - any date)
   Use this for one-off or special events
   Best for: Full flexibility

Press Tab to navigate through form fields.
Press Enter to submit the form.`,
		},
		{
			title: "Tagging & Organization",
			content: `Tags help organize and categorize events:

Available Tags:
  • project - Major project work
  • achievement - Notable accomplishments
  • leadership - Leadership activities
  • technical - Technical work
  • consulting - Advisory/consulting
  • research - Research activities
  • product - Product management
  • mentoring - Mentoring & coaching

Tips:
  • Use up to 8 tags per event
  • Tags help with filtering and searching
  • Categories are assigned automatically
  • No duplicate tags in one event`,
		},
		{
			title: "Listing & Filtering",
			content: `View and filter your career events:

Filtering Options:
  • Date Range: Filter events by date range
  • Tags: Show only events with specific tags
  • Company: Filter by company name
  • Search: Keyword search across all fields

Sorting Options:
  • By Date: Most recent first
  • By Created: When you added the event
  • By Text: Alphabetically

Navigation:
  • Up/Down: Move between events
  • Enter: View event details
  • 'f': Open filter menu
  • 's': Search events`,
		},
		{
			title: "Keyboard Shortcuts",
			content: `Common keyboard shortcuts:

Navigation:
  • 'c' - Capture new event
  • 'l' - List recent events
  • 'v' - View event details
  • 'h' - Show this help
  • 'q' - Quit application
  • 'esc' - Go back

Form Navigation:
  • Tab - Next field
  • Shift+Tab - Previous field
  • Up/Down - Navigate dropdowns
  • Enter - Confirm/Submit
  • Esc - Cancel

List Navigation:
  • Up/Down - Move between events
  • Left/Right - Previous/Next page
  • Page Up/Down - Page navigation
  • Ctrl+Home/End - First/Last event`,
		},
		{
			title: "Search & Find",
			content: `Search your career events effectively:

Search Features:
  • Keyword search in event text
  • Search in company names
  • Search in project names
  • Real-time search results

Tips:
  • Search is case-insensitive
  • Partial words match
  • Multiple words search all
  • Results update as you type

Search Examples:
  • "python" - Find Python-related events
  • "TechCorp" - Find TechCorp events
  • "led team" - Find leadership events`,
		},
		{
			title: "Tips & Best Practices",
			content: `Get the most out of KaRiya:

Event Capture:
  • Be specific in descriptions
  • Use consistent date formats
  • Add context (company, project)
  • Use tags liberally

Organization:
  • Create regular capture habits
  • Use CV Backfill for history
  • Keep descriptions concise
  • Review periodically

Searching:
  • Use simple keywords first
  • Try partial words if exact doesn't match
  • Filter by date for recent events
  • Combine tags for better filtering

CV Building:
  • Export events for interviews
  • Filter by role type
  • Show achievements first
  • Keep descriptions updated`,
		},
	}

	return &HelpModel{
		currentSection: 0,
		sections:       sections,
		width:          80,
		height:         24,
		closed:         false,
		helpFooter:       components.NewHelpFooter("help", 80),
	}
}

// Init initializes the model
func (m *HelpModel) Init() tea.Cmd {
	return nil
}

// Update handles messages
func (m *HelpModel) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch msg.String() {
		case "esc":
			// Signal back navigation to parent
			return m, func() tea.Msg { return BackMsg{} }
		case "ctrl+c", "q":
			// Signal quit to parent
			return m, func() tea.Msg { return QuitMsg{} }
		case "right", "space", "enter":
			if m.currentSection < len(m.sections)-1 {
				m.currentSection++
			}
		case "left":
			if m.currentSection > 0 {
				m.currentSection--
			}
		}
	case tea.WindowSizeMsg:
		m.width = msg.Width
		m.height = msg.Height
	}
	return m, nil
}

// View renders the help content
func (m *HelpModel) View() string {
	if m.closed {
		return ""
	}

	if m.currentSection < 0 || m.currentSection >= len(m.sections) {
		return ""
	}

	section := m.sections[m.currentSection]

	m.helpFooter.SetWidth(m.width)
	progressBar := "[ " + strings.Repeat("■", m.currentSection+1) + strings.Repeat("□", len(m.sections)-m.currentSection-1) + " ]"
	content := styles.HeaderSection.Render(section.title) + "\n\n" +
		section.content + "\n\n" +
		progressBar + "\n"
	return lipgloss.JoinVertical(lipgloss.Left, content, m.helpFooter.View())
}

// SectionCount returns the total number of help sections
func (m *HelpModel) SectionCount() int {
	return len(m.sections)
}

// CurrentSection returns the current help section index
func (m *HelpModel) CurrentSection() int {
	return m.currentSection
}

// SearchHelp searches help content for a query
func (m *HelpModel) SearchHelp(query string) string {
	query = strings.ToLower(query)
	var results []string

	for _, section := range m.sections {
		if strings.Contains(strings.ToLower(section.title), query) ||
			strings.Contains(strings.ToLower(section.content), query) {
			results = append(results, section.title)
		}
	}

	if len(results) == 0 {
		return "No help found for: " + query
	}

	return strings.Join(results, ", ")
}

// IsClosed returns whether the help screen is closed
func (m *HelpModel) IsClosed() bool {
	return m.closed
}
