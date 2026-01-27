package app

import (
	"strings"

	"github.com/baphled/kariya/internal/cli/uikit/primitives"
)

// Help section with key bindings.
type helpSection struct {
	title    string
	bindings []keyBinding
}

// Individual key binding.
type keyBinding struct {
	key  string
	desc string
}

// helpSections defines all help content.
var helpSections = []helpSection{
	{
		title: "Global Shortcuts",
		bindings: []keyBinding{
			{"?", "Toggle this help screen"},
			{"q", "Quit application (from menu)"},
			{"Ctrl+C", "Force quit application"},
			{"Esc", "Go back / Cancel / Return to menu"},
		},
	},
	{
		title: "Navigation",
		bindings: []keyBinding{
			{"↑/k", "Move up / Previous item"},
			{"↓/j", "Move down / Next item"},
			{"←/h", "Move left / Previous"},
			{"→/l", "Move right / Next"},
			{"Enter", "Confirm / Select"},
			{"Space", "Toggle selection"},
		},
	},
	{
		title: "Forms & Input",
		bindings: []keyBinding{
			{"Tab", "Next field"},
			{"Shift+Tab", "Previous field"},
			{"Ctrl+O", "Toggle optional fields"},
		},
	},
	{
		title: "Lists & Browse",
		bindings: []keyBinding{
			{"e", "Edit selected item"},
			{"d", "Delete selected item"},
			{"/", "Search"},
			{"f", "Filter"},
		},
	},
}

// Border characters.
const (
	helpDoubleLine = "════════════════════════════════════════════════════════"
	helpSingleLine = "──────────────────────────────────────────────────────────"
)

// renderHelpScreen renders the keyboard reference help screen using UIKit primitives.
func (m *Model) renderHelpScreen() string {
	var lines []string

	// Title.
	lines = append(lines, "")
	lines = append(lines, "  "+primitives.Title("KaRiya Keyboard Reference", m.theme).Render())
	lines = append(lines, "  "+primitives.Muted(helpDoubleLine, m.theme).Render())
	lines = append(lines, "")

	// Render each section.
	for _, section := range helpSections {
		lines = append(lines, "  "+primitives.InfoText(section.title, m.theme).Bold().Render())
		lines = append(lines, "  "+primitives.Muted(helpSingleLine, m.theme).Render())

		for _, binding := range section.bindings {
			keyText := primitives.SuccessText(padRight("  "+binding.key, 12), m.theme).Render()
			descText := primitives.Body(binding.desc, m.theme).Render()
			lines = append(lines, keyText+descText)
		}
		lines = append(lines, "")
	}

	// Footer.
	lines = append(lines, "  "+primitives.Muted(helpDoubleLine, m.theme).Render())
	footerText := primitives.Body("  Press ", m.theme).Render() +
		primitives.SuccessText("?", m.theme).Render() +
		primitives.Body(" or ", m.theme).Render() +
		primitives.SuccessText("Esc", m.theme).Render() +
		primitives.Body(" to close this help", m.theme).Render()
	lines = append(lines, footerText)
	lines = append(lines, "")

	helpContent := strings.Join(lines, "\n")

	// Center the help screen using primitives layout helper.
	return primitives.CenterInTerminal(helpContent, m.width, m.height)
}

// padRight pads a string to the specified width.
func padRight(s string, width int) string {
	if len(s) >= width {
		return s
	}
	return s + strings.Repeat(" ", width-len(s))
}
