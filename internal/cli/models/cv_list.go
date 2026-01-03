package models

import (
	"fmt"
	"strings"

	"github.com/charmbracelet/bubbles/table"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"

	"github.com/baphled/kariya/internal/cli/components"
	"github.com/baphled/kariya/internal/cli/styles"
	"github.com/baphled/kariya/internal/domain/career"
)

// CVListModel displays a list of generated CVs using table-based list view pattern.
type CVListModel struct {
	*BaseStandardModel
	cvViews       []*career.CVView
	filtered      []*career.CVView
	table         table.Model
	listContainer *components.TableListContainer
	pagination    *PaginationHelper
	header        components.HeaderModel
	width         int
	height        int
	breadcrumbs   []string
}

// NewCVListModel creates a new CV List model.
func NewCVListModel() *CVListModel {
	columns := []table.Column{
		{Title: "Name", Width: 25},
		{Title: "Role", Width: 15},
		{Title: "Audiences", Width: 20},
		{Title: "Events", Width: 8},
		{Title: "Generated", Width: 15},
	}

	t := table.New(
		table.WithColumns(columns),
		table.WithRows([]table.Row{}),
		table.WithFocused(true),
		table.WithHeight(15),
		table.WithWidth(100),
	)

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

	m := &CVListModel{
		BaseStandardModel: NewBaseStandardModel(),
		cvViews:           make([]*career.CVView, 0),
		filtered:          make([]*career.CVView, 0),
		table:             t,
		listContainer:     components.NewTableListContainer(t, "Generated CVs", 80),
		pagination:        NewPaginationHelper(10),
		header:            components.NewHeader("📚 Generated CVs", 80),
		width:             80,
		height:            20,
		breadcrumbs:       []string{"Home", "CV Management", "Generated CVs"},
	}
	m.header.SetBreadcrumbs(m.breadcrumbs)
	m.listContainer.SetBreadcrumbs(m.breadcrumbs)
	return m
}

// Init initializes the model
func (m *CVListModel) Init() tea.Cmd {
	return nil
}

// Update handles messages
func (m *CVListModel) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.WindowSizeMsg:
		m.width = msg.Width
		m.height = msg.Height
		m.header.SetWidth(msg.Width)
		m.listContainer.SetDimensions(msg.Width, msg.Height)
		m.updateTableRows()
		return m, nil

	case tea.KeyMsg:
		switch msg.String() {
		case "esc", "q":
			return m, func() tea.Msg {
				return BackMsg{}
			}
		case "up", "k":
			if m.listContainer.GetSelectedIdx() > 0 {
				m.listContainer.MoveUp(1)
			}
			m.updateTableRows()
			return m, nil

		case "down", "j":
			if m.listContainer.GetSelectedIdx() < len(m.filtered)-1 {
				m.listContainer.MoveDown(1)
			}
			m.updateTableRows()
			return m, nil

		case "enter":
			if len(m.filtered) > 0 {
				selectedIdx := m.listContainer.GetSelectedIdx()
				if selectedIdx < len(m.filtered) {
					selectedCV := m.filtered[selectedIdx]
					return m, func() tea.Msg {
						return ViewCVMsg{cvView: selectedCV}
					}
				}
			}
			return m, nil

		case "x":
			// Export CV
			if len(m.filtered) > 0 {
				selectedIdx := m.listContainer.GetSelectedIdx()
				if selectedIdx < len(m.filtered) {
					selectedCV := m.filtered[selectedIdx]
					return m, func() tea.Msg {
						return ShowExportOptionsMsg{CVView: selectedCV}
					}
				}
			}
			return m, nil
		}
	}

	return m, nil
}

// updateTableRows updates the table with rows from filtered CVs
func (m *CVListModel) updateTableRows() {
	var rows []table.Row
	pageCVs := m.getPageCVs()

	selectedIdx := m.listContainer.GetSelectedIdx()

	if len(pageCVs) > 0 {
		if selectedIdx >= len(pageCVs) {
			selectedIdx = len(pageCVs) - 1
		}
		if selectedIdx < 0 {
			selectedIdx = 0
		}
	} else {
		selectedIdx = 0
	}

	for i, cv := range pageCVs {
		name := cv.Name
		if len(name) > 23 {
			name = name[:20] + "..."
		}

		if i == selectedIdx {
			name = "▶ " + name
		} else {
			name = "  " + name
		}

		audiences := strings.Join(cv.TargetAudience, ", ")
		if len(audiences) > 18 {
			audiences = audiences[:15] + "..."
		}

		row := table.Row{
			name,
			cv.TargetRole,
			audiences,
			fmt.Sprintf("%d", cv.SourceEventCount),
			cv.GeneratedAt.Format("2006-01-02 15:04"),
		}
		rows = append(rows, row)
	}

	m.table.SetRows(rows)
	m.listContainer.SetSelectedIdx(selectedIdx)
	m.table.SetCursor(selectedIdx)
}

// getPageCVs returns the CVs for the current page
func (m *CVListModel) getPageCVs() []*career.CVView {
	if len(m.filtered) == 0 {
		return []*career.CVView{}
	}

	startIdx := m.pagination.GetPageStartIndex()
	endIdx := m.pagination.GetPageEndIndex()

	if startIdx >= len(m.filtered) {
		startIdx = 0
	}
	if endIdx > len(m.filtered) {
		endIdx = len(m.filtered)
	}

	if startIdx >= endIdx {
		return []*career.CVView{}
	}

	return m.filtered[startIdx:endIdx]
}

// View renders the CV List screen.
func (m *CVListModel) View() string {
	m.listContainer.SetTable(m.table).
		SetDimensions(m.width, m.height).
		SetHelpFooterKey("cv_list").
		SetBreadcrumbs(m.breadcrumbs)

	startIdx := m.pagination.GetPageStartIndex()
	endIdx := m.pagination.GetPageEndIndex()
	if endIdx > len(m.filtered) {
		endIdx = len(m.filtered)
	}
	paginationText := fmt.Sprintf("Showing %d-%d of %d CVs", startIdx+1, endIdx, len(m.filtered))
	m.listContainer.SetPaginationInfo(paginationText)

	headerView := m.header.View()
	contentView := m.listContainer.Render()
	return fmt.Sprintf("%s\n\n%s", headerView, contentView)
}

// GetSelectedCV returns the currently selected CV.
func (m *CVListModel) GetSelectedCV() *career.CVView {
	if len(m.filtered) > 0 {
		selectedIdx := m.listContainer.GetSelectedIdx()
		if selectedIdx < len(m.filtered) {
			return m.filtered[selectedIdx]
		}
	}
	return nil
}

// Messages for CV List

// ViewCVMsg is sent when a CV is selected for viewing.
type ViewCVMsg struct {
	cvView *career.CVView
}
