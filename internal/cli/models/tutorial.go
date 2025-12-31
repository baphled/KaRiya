package models

import (
	"github.com/baphled/kariya/internal/cli/styles"
	"github.com/baphled/kariya/internal/cli/components"
	"github.com/charmbracelet/lipgloss"
	tea "github.com/charmbracelet/bubbletea"
)

// TutorialModel represents the first-run tutorial screen
type TutorialModel struct {
	currentStep int
	maxSteps    int
	width       int
	height      int
	skipped     bool
	completed   bool
	helpFooter  components.HelpFooterModel
}

// Tutorial step content
const (
	tutorialStepWelcome     = 0
	tutorialStepModes       = 1
	tutorialStepCapture     = 2
	tutorialStepTags        = 3
	tutorialStepFiltering   = 4
	tutorialStepExporting   = 5
	tutorialStepGettingHelp = 6
	tutorialTotalSteps      = 7
)

// NewTutorialModel creates a new tutorial model
func NewTutorialModel() *TutorialModel {
	return &TutorialModel{
		currentStep: tutorialStepWelcome,
		maxSteps:    tutorialTotalSteps,
		width:       80,
		height:      24,
		skipped:     false,
		completed:   false,
		helpFooter:  components.NewHelpFooter("tutorial", 80),
	}
}

// Init initializes the model
func (m *TutorialModel) Init() tea.Cmd {
	return nil
}

// Update handles messages
func (m *TutorialModel) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch msg.String() {
		case "esc":
			// Skip tutorial and signal back navigation to parent
			m.skipped = true
			return m, func() tea.Msg { return BackMsg{} }
		case "ctrl+c", "q":
			// Signal quit to parent
			return m, func() tea.Msg { return QuitMsg{} }
		case "right", "space", "enter":
			if m.currentStep < m.maxSteps-1 {
				m.currentStep++
			} else {
				m.completed = true
				// Tutorial completed, go back to previous screen
				return m, func() tea.Msg { return BackMsg{} }
			}
		case "left":
			if m.currentStep > 0 {
				m.currentStep--
			}
		}
	case tea.WindowSizeMsg:
		m.width = msg.Width
		m.height = msg.Height
	}
	return m, nil
}

// View renders the tutorial
func (m *TutorialModel) View() string {
	if m.skipped || m.completed {
		return ""
	}

	var content string
	switch m.currentStep {
	case tutorialStepWelcome:
		content = m.renderWelcome()
	case tutorialStepModes:
		content = m.renderModes()
	case tutorialStepCapture:
		content = m.renderCapture()
	case tutorialStepTags:
		content = m.renderTags()
	case tutorialStepFiltering:
		content = m.renderFiltering()
	case tutorialStepExporting:
		content = m.renderExporting()
	case tutorialStepGettingHelp:
		content = m.renderGettingHelp()
	default:
		content = "Tutorial complete!"
	}

	m.helpFooter.SetWidth(m.width)
	return lipgloss.JoinVertical(lipgloss.Left, content, m.helpFooter.View())
}

// IsCompleted returns whether the tutorial is completed or skipped
func (m *TutorialModel) IsCompleted() bool {
	return m.completed || m.skipped
}

// IsSkipped returns whether the user skipped the tutorial
func (m *TutorialModel) IsSkipped() bool {
	return m.skipped
}

// CurrentStep returns the current tutorial step
func (m *TutorialModel) CurrentStep() int {
	return m.currentStep
}

func (m *TutorialModel) renderWelcome() string {
	return styles.HeaderSection.Render("Welcome to KaRiya!") + "\n\n" +
		styles.InputLabel.Render("What is KaRiya?") + "\n" +
		"KaRiya helps you journal your career events and build a comprehensive\n" +
		"record of your professional growth and achievements.\n\n" +
		"This tutorial will guide you through the key features (7 steps):\n" +
		"  1. Welcome (this step)\n" +
		"  2. Capture Modes\n" +
		"  3. Event Capture\n" +
		"  4. Tags & Categories\n" +
		"  5. Filtering Events\n" +
		"  6. Exporting\n" +
		"  7. Getting Help\n\n" +
		"Press SPACE to continue or ESC to skip this tutorial"
}

func (m *TutorialModel) renderModes() string {
	return styles.HeaderSection.Render("Capture Modes") + "\n\n" +
		styles.InputLabel.Render("Three ways to add events:") + "\n\n" +
		"📅 Timeline Journaling (recent events)\n" +
		"   Use this to log events as they happen\n" +
		"   Limited to the last 30 days\n\n" +
		"📋 CV Backfill (historical events)\n" +
		"   Use this to import events from your CV\n" +
		"   Works with any date in the past\n\n" +
		"✏️  Manual Entry (flexible)\n" +
		"   Use this for any event at any time\n" +
		"   Best for one-off or special events\n\n" +
		"Press SPACE to continue"
}

func (m *TutorialModel) renderCapture() string {
	return styles.HeaderSection.Render("Capturing Events") + "\n\n" +
		styles.InputLabel.Render("What makes a good event entry?") + "\n\n" +
		"📝 Event Description (required)\n" +
		"   Be specific about what you did\n" +
		"   1-2000 characters\n\n" +
		"📅 Date (optional)\n" +
		"   When did this happen?\n" +
		"   Use YYYY-MM-DD format\n\n" +
		"🏢 Company & Project (optional)\n" +
		"   Context helps later when reviewing\n\n" +
		"🏷️  Capture Mode (required)\n" +
		"   Choose the mode that fits your event\n\n" +
		"Press SPACE to continue"
}

func (m *TutorialModel) renderTags() string {
	return styles.HeaderSection.Render("Tags & Categories") + "\n\n" +
		styles.InputLabel.Render("Organize with tags:") + "\n\n" +
		"🏷️  Available Tags:\n" +
		"   • project, achievement, leadership\n" +
		"   • technical, consulting, research\n" +
		"   • product, mentoring\n\n" +
		"🎯 Categories (auto-assigned):\n" +
		"   • Leadership, Technical, Product\n" +
		"   • Consulting, Research, Mentoring\n\n" +
		"Tags help you find related events later\n" +
		"Use up to 8 tags per event for best results\n\n" +
		"Press SPACE to continue"
}

func (m *TutorialModel) renderFiltering() string {
	return styles.HeaderSection.Render("Filtering & Searching") + "\n\n" +
		styles.InputLabel.Render("Find your events quickly:") + "\n\n" +
		"📅 By Date Range\n" +
		"   Filter events from one date to another\n\n" +
		"🏷️  By Tags\n" +
		"   Show only events with specific tags\n\n" +
		"🔍 Search\n" +
		"   Keyword search across event text\n\n" +
		"📊 Sort\n" +
		"   Order by date, creation time, or text\n\n" +
		"Press SPACE to continue"
}

func (m *TutorialModel) renderExporting() string {
	return styles.HeaderSection.Render("Exporting Events") + "\n\n" +
		styles.InputLabel.Render("Share your career story:") + "\n\n" +
		"📄 Export to JSON\n" +
		"   Raw data for analysis or integration\n\n" +
		"📊 Export to CSV\n" +
		"   Open in spreadsheets for analysis\n\n" +
		"📝 Export to Markdown\n" +
		"   Create shareable documents\n\n" +
		"🎯 Role-Specific Views\n" +
		"   Filter events for specific audiences\n\n" +
		"Press SPACE to continue"
}

func (m *TutorialModel) renderGettingHelp() string {
	return styles.HeaderSection.Render("Getting Help") + "\n\n" +
		styles.InputLabel.Render("When you need assistance:") + "\n\n" +
		"❓ Press 'h' anytime for help\n" +
		"   Full command reference\n\n" +
		"💡 Context-sensitive hints\n" +
		"   Available on most screens\n\n" +
		"🎓 View tutorial again\n" +
		"   Available from the help menu\n\n" +
		"📖 Full documentation\n" +
		"   Online resources and guides\n\n" +
		"You're ready to start!\n" +
		"Press SPACE to begin or ESC to skip"
}
