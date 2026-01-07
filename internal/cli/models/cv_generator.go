package models

import (
	"context"
	"fmt"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"

	"github.com/baphled/kariya/internal/cli/components"
	"github.com/baphled/kariya/internal/domain/career"
	"github.com/baphled/kariya/internal/service/career/cv"
)

// CVGeneratorModel manages the CV generation process.
type CVGeneratorModel struct {
	*BaseStandardModel
	cvService   cv.CVGenerationService
	config      *career.CVConfig
	generating  bool
	generatedCV *career.CVView
	headerModel components.HeaderModel
}

// NewCVGeneratorModel creates a new CV Generator model.
func NewCVGeneratorModel(
	baseModel *BaseStandardModel,
	cvService cv.CVGenerationService,
	config *career.CVConfig,
) *CVGeneratorModel {
	m := &CVGeneratorModel{
		BaseStandardModel: baseModel,
		cvService:         cvService,
		config:            config,
		generating:        true,
		headerModel:       components.NewHeader("Generating CV", 80),
	}
	// Note: breadcrumbs now handled by StandardView
	return m
}

// Init initializes the generator and starts CV generation.
func (m *CVGeneratorModel) Init() tea.Cmd {
	return tea.Batch(
		m.BaseStandardModel.Init(),
		m.generateCV(),
	)
}

// generateCV triggers the CV generation process.
func (m *CVGeneratorModel) generateCV() tea.Cmd {
	return func() tea.Msg {
		// Validate required dependencies
		if m.cvService == nil {
			return CVGenerationErrorMsg{err: fmt.Errorf("CV generation service is not initialized")}
		}
		if m.config == nil {
			return CVGenerationErrorMsg{err: fmt.Errorf("CV configuration is not available")}
		}

		ctx := context.Background()
		cvView, err := m.cvService.GenerateCVFromConfig(ctx, m.config)
		if err != nil {
			return CVGenerationErrorMsg{err: err}
		}
		return CVGeneratedMsg{cvView: cvView}
	}
}

// Update handles messages and updates the model state.
func (m *CVGeneratorModel) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.WindowSizeMsg:
		m.headerModel.SetWidth(msg.Width)
		return m, nil

	case CVGeneratedMsg:
		m.generating = false
		m.generatedCV = msg.cvView
		return m, func() tea.Msg {
			return NavigateToCVPreviewMsg{CVView: m.generatedCV, SourceScreen: "cv_config_manager"}
		}

	case CVGenerationErrorMsg:
		m.generating = false
		m.SetError(msg.err)
		return m, nil

	case tea.KeyMsg:
		if !m.generating {
			// Only handle input when not generating
			switch msg.String() {
			case "esc":
				return m, func() tea.Msg {
					return BackMsg{}
				}
			case "q":
				return m, func() tea.Msg {
					return QuitMsg{}
				}
			}
		}
		// Ignore input when generating
		return m, nil
	}

	// Return model unchanged for unhandled messages
	return m, nil
}

// View renders the CV Generator screen.
func (m *CVGeneratorModel) View() string {
	headerContent := m.headerModel.View()

	if m.generating {
		return fmt.Sprintf("%s\n\n%s", headerContent, m.renderGenerating())
	}

	if m.GetLastError() != nil {
		return fmt.Sprintf("%s\n\n%s", headerContent, m.renderError())
	}

	return fmt.Sprintf("%s\n\nCV generated successfully!", headerContent)
}

// renderGenerating renders the loading state.
func (m *CVGeneratorModel) renderGenerating() string {
	title := lipgloss.NewStyle().
		Bold(true).
		Foreground(lipgloss.Color("6")).
		Render("Generating CV...")

	spinner := []string{"|", "/", "-", "\\"}
	spinnerIdx := 0
	spinnerFrame := spinner[spinnerIdx%len(spinner)]

	config := m.config
	info := fmt.Sprintf("Role: %s\nAudiences: %v\nEvents in scope: ∞",
		config.TargetRole,
		config.TargetAudience)

	status := lipgloss.NewStyle().
		Foreground(lipgloss.Color("3")).
		Render("Analyzing events and generating CV bullets...")

	return fmt.Sprintf("%s\n\n%s\n\n%s\n\n%s",
		title,
		info,
		status,
		spinnerFrame)
}

// renderError renders the error state.
func (m *CVGeneratorModel) renderError() string {
	errorStyle := lipgloss.NewStyle().
		Foreground(lipgloss.Color("1")).
		Bold(true)

	title := lipgloss.NewStyle().
		Bold(true).
		Foreground(lipgloss.Color("1")).
		Render("CV Generation Failed")

	err := fmt.Sprintf("Error: %v", m.GetLastError())

	footer := "Press 'esc' to go back"

	return fmt.Sprintf("%s\n\n%s\n\n%s",
		title,
		errorStyle.Render(err),
		footer)
}

// Messages for CV Generator

// CVGeneratedMsg is sent when CV generation completes successfully.
type CVGeneratedMsg struct {
	cvView *career.CVView
}

// CVGenerationErrorMsg is sent when CV generation fails.
type CVGenerationErrorMsg struct {
	err error
}

// NavigateToCVPreviewMsg navigates to CV preview after generation.
type NavigateToCVPreviewMsg struct {
	CVView       *career.CVView
	SourceScreen string // Track where we came from
}
