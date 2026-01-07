package models

import (
	"context"
	"fmt"
	"strings"

	"github.com/charmbracelet/bubbles/table"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"

	"github.com/baphled/kariya/internal/cli/components"
	"github.com/baphled/kariya/internal/cli/styles"
	"github.com/baphled/kariya/internal/domain/career"
	"github.com/baphled/kariya/internal/service/career/cv"
)

// CVPreviewModel displays a generated CV for review using table-based list view pattern.
// It shows CV sections as a table with the ability to export, navigate, and view details.
type CVPreviewModel struct {
	*BaseStandardModel
	cvView              *career.CVView
	sections            []*career.CVSection
	table               table.Model
	listContainer       *components.TableListContainer
	pagination          *PaginationHelper
	traceabilityService *cv.TraceabilityService
	sourceScreen        string // Track where we came from (e.g., "event_timeline", "cv_config_manager")
	header              components.HeaderModel
	helpFooter          components.HelpFooterModel
	width               int
	height              int
	breadcrumbs         []string
	selectedIdx         int
	expandedSections    map[int]bool
	ctx                 context.Context
}

// NewCVPreviewModel creates a new CV Preview model.
func NewCVPreviewModel(
	baseModel *BaseStandardModel,
	cvView *career.CVView,
	sections []*career.CVSection,
	traceabilityService *cv.TraceabilityService,
) *CVPreviewModel {
	return NewCVPreviewModelWithSource(baseModel, cvView, sections, traceabilityService, "cv_config_manager")
}

// NewCVPreviewModelWithSource creates a new CV Preview model with a specified source screen.
func NewCVPreviewModelWithSource(
	baseModel *BaseStandardModel,
	cvView *career.CVView,
	sections []*career.CVSection,
	traceabilityService *cv.TraceabilityService,
	sourceScreen string,
) *CVPreviewModel {
	// Create table with CV sections columns
	columns := []table.Column{
		{Title: "Section", Width: 20},
		{Title: "Bullets", Width: 10},
		{Title: "Content Preview", Width: 50},
	}

	t := table.New(
		table.WithColumns(columns),
		table.WithRows([]table.Row{}),
		table.WithFocused(true),
		table.WithHeight(15),
		table.WithWidth(100),
	)

	// Apply styling to table
	s := table.DefaultStyles()
	s.Header = s.Header.
		Foreground(styles.ColorAccentTeal).
		Bold(true).
		BorderStyle(lipgloss.NormalBorder()).
		BorderBottom(true).
		BorderForeground(styles.ColorAccentTeal)
	s.Selected = s.Selected.
		Foreground(styles.ColorAccentTeal).
		Background(styles.ColorBackground).
		Bold(true)
	t.SetStyles(s)

	breadcrumbs := []string{"Home", "CV Management", "Preview"}
	if sourceScreen == "event_timeline" {
		breadcrumbs = []string{"Home", "Events", "CV Preview"}
	}

	m := &CVPreviewModel{
		BaseStandardModel:   baseModel,
		cvView:              cvView,
		sections:            sections,
		table:               t,
		listContainer:       components.NewTableListContainer(t, "CV Preview", 80),
		pagination:          NewPaginationHelper(10),
		traceabilityService: traceabilityService,
		sourceScreen:        sourceScreen,
		header:              components.NewHeader(fmt.Sprintf("📄 %s - CV Preview", cvView.Name), 80),
		helpFooter:          components.NewHelpFooter("cv_preview", 80),
		width:               80,
		height:              20,
		breadcrumbs:         breadcrumbs,
		selectedIdx:         0,
		expandedSections:    make(map[int]bool),
		ctx:                 context.Background(),
	}

	// Set pagination total count to enable proper section pagination
	m.pagination.SetTotalCount(len(sections))

	// Note: breadcrumbs now handled by StandardView
	m.updateTableRows()
	return m
}

// Init initializes the preview model.
func (m *CVPreviewModel) Init() tea.Cmd {
	return m.BaseStandardModel.Init()
}

// Update handles messages and updates the model state.
func (m *CVPreviewModel) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.WindowSizeMsg:
		m.width = msg.Width
		m.height = msg.Height
		m.header.SetWidth(msg.Width)
		m.helpFooter.SetWidth(msg.Width)
		m.listContainer.SetDimensions(msg.Width, msg.Height)
		m.updateTableRows()
		return m, nil

	case tea.KeyMsg:
		switch msg.String() {
		case "up", "k":
			m.listContainer.MoveUp(1)
			m.updateTableRows()
			return m, nil

		case "down", "j":
			m.listContainer.MoveDown(1)
			m.updateTableRows()
			return m, nil

		case "home", "g":
			m.listContainer.MoveToFirst()
			m.updateTableRows()
			return m, nil

		case "end", "G":
			m.listContainer.MoveToLast()
			m.updateTableRows()
			return m, nil

		case "enter":
			// Expand/collapse current section to see full content
			selectedIdx := m.listContainer.GetSelectedIdx()
			m.expandedSections[selectedIdx] = !m.expandedSections[selectedIdx]
			m.updateTableRows()
			return m, nil

		case "x":
			// Export CV - show export options
			return m, func() tea.Msg {
				return ShowExportOptionsMsg{CVView: m.cvView}
			}

		case "esc":
			if m.sourceScreen == "event_timeline" {
				return m, func() tea.Msg {
					return BackToEventTimelineMsg{}
				}
			}
			return m, func() tea.Msg {
				return BackMsg{}
			}
		case "q", "ctrl+c":
			return m, func() tea.Msg { return QuitMsg{} }
		}
	}

	// Return model unchanged for unhandled messages
	return m, nil
}

// updateTableRows updates the table with rows from CV sections
func (m *CVPreviewModel) updateTableRows() {
	var rows []table.Row
	pageConfigs := m.getPageSections()

	// Get the selected index within the current page
	selectedIdx := m.listContainer.GetSelectedIdx()

	// Ensure cursor is within valid bounds
	if len(pageConfigs) > 0 {
		if selectedIdx >= len(pageConfigs) {
			selectedIdx = len(pageConfigs) - 1
		}
		if selectedIdx < 0 {
			selectedIdx = 0
		}
	} else {
		selectedIdx = 0
	}

	for i, section := range pageConfigs {
		sectionTitle := section.Title
		if len(sectionTitle) > 18 {
			sectionTitle = sectionTitle[:15] + "..."
		}

		// Add focus indicator for the selected row
		if i == selectedIdx {
			sectionTitle = "▶ " + sectionTitle
		} else {
			sectionTitle = "  " + sectionTitle
		}

		// Count bullets in this section
		bulletCount := 0
		for _, group := range section.Content {
			bulletCount += len(group.Bullets)
		}
		bulletCountStr := fmt.Sprintf("%d", bulletCount)

		// Get preview of content (first bullet text or summary)
		contentPreview := ""
		if section.Summary != "" {
			contentPreview = section.Summary
		} else if len(section.Content) > 0 && len(section.Content[0].Bullets) > 0 {
			contentPreview = section.Content[0].Bullets[0].Text
		}

		if len(contentPreview) > 48 {
			contentPreview = contentPreview[:45] + "..."
		}
		contentPreview = strings.ReplaceAll(contentPreview, "\n", " ")

		row := table.Row{
			sectionTitle,
			bulletCountStr,
			contentPreview,
		}
		rows = append(rows, row)
	}

	m.table.SetRows(rows)
	m.listContainer.SetSelectedIdx(selectedIdx)
	m.table.SetCursor(selectedIdx)
}

// getPageSections returns the sections for the current page
func (m *CVPreviewModel) getPageSections() []*career.CVSection {
	if len(m.sections) == 0 {
		return []*career.CVSection{}
	}

	startIdx := m.pagination.GetPageStartIndex()
	endIdx := m.pagination.GetPageEndIndex()

	if startIdx >= len(m.sections) {
		startIdx = 0
	}
	if endIdx > len(m.sections) {
		endIdx = len(m.sections)
	}

	if startIdx >= endIdx {
		return []*career.CVSection{}
	}

	return m.sections[startIdx:endIdx]
}

// View renders the CV Preview screen.
func (m *CVPreviewModel) View() string {
	// Update list container with current state
	m.listContainer.SetTable(m.table).
		SetDimensions(m.width, m.height)

	// Set pagination info
	startIdx := m.pagination.GetPageStartIndex()
	endIdx := m.pagination.GetPageEndIndex()
	if endIdx > len(m.sections) {
		endIdx = len(m.sections)
	}
	paginationText := fmt.Sprintf("Showing %d-%d of %d sections", startIdx+1, endIdx, len(m.sections))
	m.listContainer.SetPaginationInfo(paginationText)

	headerView := m.header.View()
	metadataView := m.renderMetadata()
	contentView := m.listContainer.Render()

	// Render expanded section details if selected
	selectedIdx := m.listContainer.GetSelectedIdx()
	var expandedView string
	if m.expandedSections[selectedIdx] && selectedIdx < len(m.sections) {
		expandedView = m.renderExpandedSection(m.sections[selectedIdx])
	}

	if expandedView != "" {
		return fmt.Sprintf("%s\n\n%s\n\n%s\n\n%s", headerView, metadataView, contentView, expandedView)
	}

	return fmt.Sprintf("%s\n\n%s\n\n%s", headerView, metadataView, contentView)
}

// renderMetadata renders CV metadata (role, audience, stats)
func (m *CVPreviewModel) renderMetadata() string {
	metadata := fmt.Sprintf("Role: %s | Audience: %s | Events: %d | Facts: %d | Generated: %s",
		m.cvView.TargetRole,
		m.cvView.TargetAudience,
		m.cvView.SourceEventCount,
		m.cvView.SourceFactCount,
		m.cvView.GeneratedAt.Format("2006-01-02 15:04"))

	return styles.InputHint.Render(metadata)
}

// renderExpandedSection renders the full content of an expanded section
func (m *CVPreviewModel) renderExpandedSection(section *career.CVSection) string {
	var contentText string

	// Summary sections use Summary field
	if section.Summary != "" {
		contentText = section.Summary
	} else {
		// Other sections use Content groups
		var parts []string
		for _, group := range section.Content {
			if group.Header != "" {
				// Include header with date range if available
				header := group.Header
				if group.StartDate != "" || group.EndDate != "" {
					header += fmt.Sprintf(" (%s - %s)", group.StartDate, group.EndDate)
				}
				parts = append(parts, "\n"+header)
			}
			for _, bullet := range group.Bullets {
				parts = append(parts, "  • "+bullet.Text)
			}
		}
		contentText = strings.Join(parts, "\n")
	}

	sectionContent := fmt.Sprintf("📌 %s\n%s", section.Title, contentText)

	// Wrap in a styled container
	container := components.NewScreenContainer(sectionContent).
		WithPaddingMode(components.PaddingNormal)

	return container.Render()
}

// GetCVView returns the CV view being displayed.
func (m *CVPreviewModel) GetCVView() *career.CVView {
	return m.cvView
}

// GetSelectedSection returns the currently selected section.
func (m *CVPreviewModel) GetSelectedSection() *career.CVSection {
	selectedIdx := m.listContainer.GetSelectedIdx()
	if selectedIdx >= 0 && selectedIdx < len(m.sections) {
		return m.sections[selectedIdx]
	}
	return nil
}

// Messages for CV Preview

// ShowExportOptionsMsg triggers export options dialog.
type ShowExportOptionsMsg struct {
	CVView *career.CVView
}

// CVExportedMsg is sent when CV export completes successfully.
type CVExportedMsg struct {
	CVView   *career.CVView
	Format   string
	FilePath string
}

// CVExportErrorMsg is sent when there's an error exporting a CV.
type CVExportErrorMsg struct {
	Err error
}

// BackToEventTimelineMsg navigates back to the event timeline.
type BackToEventTimelineMsg struct{}
