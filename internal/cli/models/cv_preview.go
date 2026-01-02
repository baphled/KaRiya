package models

import (
	"fmt"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"

	"github.com/baphled/kariya/internal/domain/career"
	"github.com/baphled/kariya/internal/service/career/cv"
)

// CVPreviewModel displays a generated CV for review.
type CVPreviewModel struct {
	*BaseStandardModel
	cvView              *career.CVView
	sections            []*career.CVSection
	selectedSectionIdx  int
	selectedBulletIdx   int
	expandedBullets     map[int]bool
	traceabilityService *cv.TraceabilityService
}

// NewCVPreviewModel creates a new CV Preview model.
func NewCVPreviewModel(
	baseModel *BaseStandardModel,
	cvView *career.CVView,
	sections []*career.CVSection,
	traceabilityService *cv.TraceabilityService,
) *CVPreviewModel {
	return &CVPreviewModel{
		BaseStandardModel:   baseModel,
		cvView:              cvView,
		sections:            sections,
		selectedSectionIdx:  0,
		selectedBulletIdx:   0,
		expandedBullets:     make(map[int]bool),
		traceabilityService: traceabilityService,
	}
}

// Init initializes the preview model.
func (m *CVPreviewModel) Init() tea.Cmd {
	return m.BaseStandardModel.Init()
}

// Update handles messages and updates the model state.
func (m *CVPreviewModel) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch msg.String() {
		case "up", "k":
			m.moveBulletUp()
			return m, nil

		case "down", "j":
			m.moveBulletDown()
			return m, nil

		case "left", "h":
			m.moveSectionLeft()
			return m, nil

		case "right", "l":
			m.moveSectionRight()
			return m, nil

		case "enter":
			// Expand/collapse current bullet to see sources
			m.expandedBullets[m.selectedBulletIdx] = !m.expandedBullets[m.selectedBulletIdx]
			return m, nil

		case "e":
			// Export CV
			return m, m.exportCV()

		case "esc", "q":
			return m, func() tea.Msg {
				return BackToCVConfigManagerMsg{}
			}
		}
	}

	return m.BaseStandardModel.Update(msg)
}

// moveBulletUp moves selection up in the bullet list.
func (m *CVPreviewModel) moveBulletUp() {
	if m.selectedBulletIdx > 0 {
		m.selectedBulletIdx--
	}
}

// moveBulletDown moves selection down in the bullet list.
func (m *CVPreviewModel) moveBulletDown() {
	if m.selectedBulletIdx < m.getBulletCount()-1 {
		m.selectedBulletIdx++
	}
}

// moveSectionLeft moves selection to previous section.
func (m *CVPreviewModel) moveSectionLeft() {
	if m.selectedSectionIdx > 0 {
		m.selectedSectionIdx--
		m.selectedBulletIdx = 0
	}
}

// moveSectionRight moves selection to next section.
func (m *CVPreviewModel) moveSectionRight() {
	if m.selectedSectionIdx < len(m.sections)-1 {
		m.selectedSectionIdx++
		m.selectedBulletIdx = 0
	}
}

// getBulletCount returns the number of bullets in the current section.
func (m *CVPreviewModel) getBulletCount() int {
	if m.selectedSectionIdx >= len(m.sections) {
		return 0
	}
	// This is a simplified count - in production would parse section content
	return 1
}

// exportCV triggers CV export.
func (m *CVPreviewModel) exportCV() tea.Cmd {
	return func() tea.Msg {
		// This will trigger export action
		return ShowExportOptionsMsg{cvView: m.cvView}
	}
}

// View renders the CV Preview screen.
func (m *CVPreviewModel) View() string {
	return fmt.Sprintf("%s\n\n%s\n\n%s",
		m.renderHeader(),
		m.renderContent(),
		m.renderFooter())
}

// renderHeader renders the CV header with metadata.
func (m *CVPreviewModel) renderHeader() string {
	headerStyle := lipgloss.NewStyle().
		Bold(true).
		Foreground(lipgloss.Color("6"))

	title := m.cvView.Name
	role := fmt.Sprintf("Role: %s", m.cvView.TargetRole)
	audience := fmt.Sprintf("Audience: %v", m.cvView.TargetAudience)
	stats := fmt.Sprintf("Events: %d | Facts: %d | Generated: %s",
		m.cvView.SourceEventCount,
		m.cvView.SourceFactCount,
		m.cvView.GeneratedAt.Format("2006-01-02 15:04"))

	return fmt.Sprintf("%s\n%s | %s\n%s\n%s",
		headerStyle.Render(title),
		role, audience,
		lipgloss.NewStyle().Foreground(lipgloss.Color("8")).Render("─────────────────────────────────────────────────"),
		stats)
}

// renderContent renders the CV sections and bullets.
func (m *CVPreviewModel) renderContent() string {
	var content string

	for i, section := range m.sections {
		sectionStyle := lipgloss.NewStyle().Bold(true)
		if i == m.selectedSectionIdx {
			sectionStyle = sectionStyle.
				Background(lipgloss.Color("4")).
				Foreground(lipgloss.Color("15"))
		}

		content += sectionStyle.Render(fmt.Sprintf("## %s", section.Title)) + "\n\n"
		content += section.Content + "\n\n"
	}

	return content
}

// renderFooter renders the footer with keyboard shortcuts.
func (m *CVPreviewModel) renderFooter() string {
	footerStyle := lipgloss.NewStyle().
		Foreground(lipgloss.Color("8"))

	shortcuts := []string{
		"↑/k: Up",
		"↓/j: Down",
		"←/h: Prev Section",
		"→/l: Next Section",
		"Enter: Show Sources",
		"e: Export",
		"q: Back",
	}

	footer := ""
	for _, shortcut := range shortcuts {
		footer += footerStyle.Render(shortcut) + "  "
	}

	return footer
}

// GetCVView returns the CV view being displayed.
func (m *CVPreviewModel) GetCVView() *career.CVView {
	return m.cvView
}

// Messages for CV Preview

// ShowExportOptionsMsg triggers export options dialog.
type ShowExportOptionsMsg struct {
	cvView *career.CVView
}

