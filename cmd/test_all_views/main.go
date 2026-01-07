package main

import (
	"fmt"
	"os"
	"time"

	"github.com/baphled/kariya/internal/cli/components"
	"github.com/baphled/kariya/internal/cli/terminal"
	tea "github.com/charmbracelet/bubbletea"
)

// Model for the visual test program
type model struct {
	currentScenario int
	scenarios       []scenario
	termInfo        *terminal.Info
	quitting        bool
}

type scenario struct {
	name        string
	description string
	viewFunc    func(*terminal.Info) string
}

func initialModel() model {
	termInfo := &terminal.Info{Width: 120, Height: 40}

	scenarios := []scenario{
		{
			name:        "Basic StandardView",
			description: "Logo, content, and footer",
			viewFunc: func(ti *terminal.Info) string {
				logo := components.NewASCIILogo(false, ti.Width)
				view := components.NewStandardView(ti).
					WithLogo(logo, 2).
					WithContent("This is a basic StandardView with logo, content, and footer.\n\nResize your terminal to test responsiveness.").
					WithHelp("← → Navigate  q Quit").
					WithFooterSeparator(true)
				return view.Render()
			},
		},
		{
			name:        "StandardView with Breadcrumbs",
			description: "Logo, breadcrumbs, content, footer",
			viewFunc: func(ti *terminal.Info) string {
				logo := components.NewASCIILogo(false, ti.Width)
				view := components.NewStandardView(ti).
					WithLogo(logo, 2).
					WithBreadcrumbs("Main Menu", "Settings", "Display").
					WithContent("StandardView with breadcrumbs showing navigation context.\n\nBreadcrumbs: Main Menu > Settings > Display").
					WithHelp("← → Navigate  q Quit").
					WithFooterSeparator(true)
				return view.Render()
			},
		},
		{
			name:        "StandardView with Error Modal",
			description: "Background view with error modal overlay",
			viewFunc: func(ti *terminal.Info) string {
				logo := components.NewASCIILogo(false, ti.Width)
				modal := components.NewErrorModal("Error Occurred", "An unexpected error occurred while processing your request. Please try again.")
				view := components.NewStandardView(ti).
					WithLogo(logo, 2).
					WithBreadcrumbs("Main Menu", "Operations").
					WithContent("Background content (should be visible but dimmed).\n\nThe error modal appears as an overlay.").
					WithHelp("Esc Dismiss  q Quit").
					ShowModalOverlay(modal)
				return view.Render()
			},
		},
		{
			name:        "StandardView with Loading Modal",
			description: "Background view with loading modal",
			viewFunc: func(ti *terminal.Info) string {
				logo := components.NewASCIILogo(false, ti.Width)
				modal := components.NewLoadingModal("Processing your request...", false)
				view := components.NewStandardView(ti).
					WithLogo(logo, 2).
					WithContent("Background content while loading...").
					WithHelp("q Quit").
					ShowModalOverlay(modal)
				return view.Render()
			},
		},
		{
			name:        "StandardView with Progress Modal (0%)",
			description: "Progress modal at 0%",
			viewFunc: func(ti *terminal.Info) string {
				logo := components.NewASCIILogo(false, ti.Width)
				modal := components.NewProgressModal("Starting", "Initializing process...", 0.0)
				view := components.NewStandardView(ti).
					WithLogo(logo, 2).
					WithContent("Progress: Starting...").
					ShowModalOverlay(modal)
				return view.Render()
			},
		},
		{
			name:        "StandardView with Progress Modal (25%)",
			description: "Progress modal at 25%",
			viewFunc: func(ti *terminal.Info) string {
				logo := components.NewASCIILogo(false, ti.Width)
				modal := components.NewProgressModal("Processing", "Analyzing data...", 0.25)
				view := components.NewStandardView(ti).
					WithLogo(logo, 2).
					WithContent("Progress: 25% complete").
					ShowModalOverlay(modal)
				return view.Render()
			},
		},
		{
			name:        "StandardView with Progress Modal (50%)",
			description: "Progress modal at 50%",
			viewFunc: func(ti *terminal.Info) string {
				logo := components.NewASCIILogo(false, ti.Width)
				modal := components.NewProgressModal("Processing", "Generating results...", 0.50)
				view := components.NewStandardView(ti).
					WithLogo(logo, 2).
					WithContent("Progress: 50% complete").
					ShowModalOverlay(modal)
				return view.Render()
			},
		},
		{
			name:        "StandardView with Progress Modal (75%)",
			description: "Progress modal at 75%",
			viewFunc: func(ti *terminal.Info) string {
				logo := components.NewASCIILogo(false, ti.Width)
				modal := components.NewProgressModal("Finalizing", "Almost done...", 0.75)
				view := components.NewStandardView(ti).
					WithLogo(logo, 2).
					WithContent("Progress: 75% complete").
					ShowModalOverlay(modal)
				return view.Render()
			},
		},
		{
			name:        "StandardView with Progress Modal (100%)",
			description: "Progress modal at 100%",
			viewFunc: func(ti *terminal.Info) string {
				logo := components.NewASCIILogo(false, ti.Width)
				modal := components.NewProgressModal("Complete", "Process completed!", 1.0)
				view := components.NewStandardView(ti).
					WithLogo(logo, 2).
					WithContent("Progress: 100% complete").
					ShowModalOverlay(modal)
				return view.Render()
			},
		},
		{
			name:        "StandardView with Success Modal",
			description: "Success modal with auto-dismiss",
			viewFunc: func(ti *terminal.Info) string {
				logo := components.NewASCIILogo(false, ti.Width)
				modal := components.NewSuccessModal("Operation completed successfully! File saved to /path/to/file.txt")
				view := components.NewStandardView(ti).
					WithLogo(logo, 2).
					WithContent("Success! Check the modal above.").
					WithHelp("q Quit").
					ShowModalOverlay(modal)
				return view.Render()
			},
		},
		{
			name:        "StandardView with Very Long Content",
			description: "Tests content overflow handling",
			viewFunc: func(ti *terminal.Info) string {
				logo := components.NewASCIILogo(false, ti.Width)
				longContent := "This is a very long content example.\n\n"
				for i := 1; i <= 50; i++ {
					longContent += fmt.Sprintf("Line %d: Lorem ipsum dolor sit amet, consectetur adipiscing elit.\n", i)
				}
				view := components.NewStandardView(ti).
					WithLogo(logo, 2).
					WithBreadcrumbs("Main Menu", "Content", "Long Document").
					WithContent(longContent).
					WithHelp("↑↓ Scroll  q Quit").
					WithFooterSeparator(true)
				return view.Render()
			},
		},
		{
			name:        "StandardView Minimal Content",
			description: "Minimal content to test spacing",
			viewFunc: func(ti *terminal.Info) string {
				logo := components.NewASCIILogo(false, ti.Width)
				view := components.NewStandardView(ti).
					WithLogo(logo, 2).
					WithContent("Minimal").
					WithHelp("q Quit")
				return view.Render()
			},
		},
	}

	return model{
		currentScenario: 0,
		scenarios:       scenarios,
		termInfo:        termInfo,
		quitting:        false,
	}
}

func (m model) Init() tea.Cmd {
	return nil
}

func (m model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.WindowSizeMsg:
		m.termInfo.Width = msg.Width
		m.termInfo.Height = msg.Height
		return m, nil

	case tea.KeyMsg:
		switch msg.String() {
		case "q", "ctrl+c":
			m.quitting = true
			return m, tea.Quit

		case "right", "l":
			m.currentScenario = (m.currentScenario + 1) % len(m.scenarios)
			return m, nil

		case "left", "h":
			m.currentScenario = (m.currentScenario - 1 + len(m.scenarios)) % len(m.scenarios)
			return m, nil
		}
	}

	return m, nil
}

func (m model) View() string {
	if m.quitting {
		return "Thanks for testing StandardView!\n"
	}

	// Get current scenario
	scenario := m.scenarios[m.currentScenario]

	// Add instructions at the top
	instructions := "\n  Visual Test Program for StandardView\n"
	instructions += fmt.Sprintf("  Terminal Size: %dx%d\n", m.termInfo.Width, m.termInfo.Height)
	instructions += fmt.Sprintf("  Scenario %d/%d: %s\n", m.currentScenario+1, len(m.scenarios), scenario.name)
	instructions += fmt.Sprintf("  %s\n\n", scenario.description)
	instructions += "  ← → or h/l Navigate scenarios  q Quit\n"
	instructions += "  " + time.Now().Format("15:04:05") + "\n\n"

	// Render the scenario
	scenarioView := scenario.viewFunc(m.termInfo)

	return instructions + scenarioView
}

func main() {
	p := tea.NewProgram(initialModel(), tea.WithAltScreen())
	if _, err := p.Run(); err != nil {
		fmt.Printf("Error: %v\n", err)
		os.Exit(1)
	}
}
