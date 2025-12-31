package models

import (
	"context"
	"fmt"
	"sort"
	"strings"

	"github.com/baphled/kariya/internal/cli/styles"
	"github.com/baphled/kariya/internal/domain/career"
	careerrepo "github.com/baphled/kariya/internal/repository/career"
	careerservice "github.com/baphled/kariya/internal/service/career"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
)

// BurstsLoadedMsg is sent when bursts have been loaded from repository
type BurstsLoadedMsg struct {
	Bursts []*career.Burst
	Err    error
}

// BurstListModel represents the burst list display screen
type BurstListModel struct {
	service         *careerservice.Service
	ctx             context.Context
	bursts          []*career.Burst
	selectedIdx     int
	expandedIndices map[int]bool
	filterBy        string // Competency focus filter
	sortBy          string // "date", "event_count", "name"
	width           int
	height          int
}

// NewBurstListModel creates a new burst list model
func NewBurstListModel(svc *careerservice.Service, ctx context.Context) *BurstListModel {
	return &BurstListModel{
		service:         svc,
		ctx:             ctx,
		bursts:          []*career.Burst{},
		selectedIdx:     0,
		expandedIndices: make(map[int]bool),
		filterBy:        "",
		sortBy:          "date",
		width:           80,
		height:          24,
	}
}

// Init initializes the model
func (m *BurstListModel) Init() tea.Cmd {
	return m.loadBursts()
}

// loadBursts loads bursts from the service
func (m *BurstListModel) loadBursts() tea.Cmd {
	return func() tea.Msg {
		// Get burst repository from service
		burstRepo := m.service.GetBurstRepository()
		if burstRepo == nil {
			return BurstsLoadedMsg{Bursts: []*career.Burst{}, Err: fmt.Errorf("burst repository not configured")}
		}

		// Load bursts from repository
		bursts, err := burstRepo.List(m.ctx, careerrepo.BurstListFilters{Limit: 1000})
		return BurstsLoadedMsg{Bursts: bursts, Err: err}
	}
}

// Update handles messages
func (m *BurstListModel) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.WindowSizeMsg:
		m.width = msg.Width
		m.height = msg.Height
		return m, nil

	case BurstsLoadedMsg:
		if msg.Err != nil {
			// Handle error - for now just log it and continue with empty list
			// In the future we could show an error message to the user
		} else {
			m.bursts = msg.Bursts
		}
		return m, nil

	case tea.KeyMsg:
		return m.handleKeyMsg(msg), nil
	}

	return m, nil
}

// handleKeyMsg processes keyboard input
func (m *BurstListModel) handleKeyMsg(msg tea.KeyMsg) tea.Model {
	displayedBursts := m.getDisplayedBursts()

	switch msg.Type {
	case tea.KeyUp:
		if len(displayedBursts) > 0 {
			m.selectedIdx = (m.selectedIdx - 1 + len(displayedBursts)) % len(displayedBursts)
		}
	case tea.KeyDown:
		if len(displayedBursts) > 0 {
			m.selectedIdx = (m.selectedIdx + 1) % len(displayedBursts)
		}
	case tea.KeySpace:
		m.expandedIndices[m.selectedIdx] = !m.expandedIndices[m.selectedIdx]
	case tea.KeyEsc:
		// Handle escape key for back navigation
		// Model will be dismissed by parent app
	}

	// Handle character input
	if msg.Type == tea.KeyRunes {
		for _, r := range msg.Runes {
			switch r {
			case 'q':
				// Handle quit request
				// Model will be dismissed by parent app
			}
		}
	}

	return m
}

// View renders the burst list
func (m *BurstListModel) View() string {
	displayedBursts := m.getDisplayedBursts()

	if len(displayedBursts) == 0 {
		if m.filterBy != "" {
			return styles.ErrorBox.Render("No matching bursts")
		}
		return styles.InfoBox.Render("No bursts found")
	}

	var sb strings.Builder

	// Render header
	sb.WriteString(m.renderHeader())
	sb.WriteString("\n")

	// Render burst list
	for i, burstIdx := range displayedBursts {
		if i >= m.height-5 {
			// Stop rendering if we exceed visible height
			break
		}

		burst := m.bursts[burstIdx]
		sb.WriteString(m.renderBurstRow(burst, i == m.selectedIdx))

		// Render expanded events if this burst is expanded
		if m.expandedIndices[burstIdx] {
			sb.WriteString(m.renderExpandedEvents(burst))
		}
	}

	// Render footer
	sb.WriteString("\n")
	sb.WriteString(m.renderFooter())

	return sb.String()
}

// renderHeader renders the column headers
func (m *BurstListModel) renderHeader() string {
	nameCol := lipgloss.NewStyle().Width(30).Render("Name")
	countCol := lipgloss.NewStyle().Width(8).Render("Events")
	compCol := lipgloss.NewStyle().Width(15).Render("Competency")
	dateCol := lipgloss.NewStyle().Width(12).Render("Created")

	header := nameCol + " " + countCol + " " + compCol + " " + dateCol
	return styles.HeaderSection.Render(header)
}

// renderBurstRow renders a single burst row
func (m *BurstListModel) renderBurstRow(burst *career.Burst, isSelected bool) string {
	name := truncateString(burst.Name, 28)
	count := fmt.Sprintf("%d", len(burst.EventIDs))
	competency := truncateString(burst.CompetencyFocus, 13)
	date := burst.CreatedAt.Format("2006-01-02")

	nameCol := lipgloss.NewStyle().Width(30).Render(name)
	countCol := lipgloss.NewStyle().Width(8).Render(count)
	compCol := lipgloss.NewStyle().Width(15).Render(competency)
	dateCol := lipgloss.NewStyle().Width(12).Render(date)

	row := nameCol + " " + countCol + " " + compCol + " " + dateCol

	// Apply selection highlighting
	if isSelected {
		row = styles.ListItemSelected.Render(row)
	} else {
		row = styles.ListItem.Render(row)
	}

	return row + "\n"
}

// renderExpandedEvents renders the event IDs for an expanded burst
func (m *BurstListModel) renderExpandedEvents(burst *career.Burst) string {
	var sb strings.Builder
	sb.WriteString("  ")
	sb.WriteString(styles.InfoHint.Render("Events: "))

	for i, eventID := range burst.EventIDs {
		if i > 0 {
			sb.WriteString(", ")
		}
		// Truncate event ID for display
		truncated := truncateString(eventID, 8)
		sb.WriteString(styles.InfoText.Render(truncated))
	}

	sb.WriteString("\n")
	return sb.String()
}

// renderFooter renders the footer with navigation hints
func (m *BurstListModel) renderFooter() string {
	hints := []string{
		"↑/↓: Navigate",
		"Space: Expand",
		"f: Filter",
		"s: Sort",
		"Esc: Back",
	}
	footerText := strings.Join(hints, " • ")
	return styles.HeaderSection.Render(footerText)
}

// getDisplayedBursts returns the filtered and sorted bursts as indices
func (m *BurstListModel) getDisplayedBursts() []int {
	var displayed []int

	// Filter bursts
	for i, burst := range m.bursts {
		if m.filterBy == "" || burst.CompetencyFocus == m.filterBy {
			displayed = append(displayed, i)
		}
	}

	// Sort displayed burst indices
	m.sortBursts(displayed)

	return displayed
}

// sortBursts sorts the burst indices based on current sort order
func (m *BurstListModel) sortBursts(indices []int) {
	switch m.sortBy {
	case "name":
		sort.Slice(indices, func(i, j int) bool {
			return m.bursts[indices[i]].Name < m.bursts[indices[j]].Name
		})

	case "event_count":
		sort.Slice(indices, func(i, j int) bool {
			return len(m.bursts[indices[i]].EventIDs) < len(m.bursts[indices[j]].EventIDs)
		})

	case "date":
		fallthrough
	default:
		sort.Slice(indices, func(i, j int) bool {
			return m.bursts[indices[i]].CreatedAt.Before(m.bursts[indices[j]].CreatedAt)
		})
	}
}

// truncateString truncates a string to a maximum width, adding "..." if truncated
func truncateString(s string, maxWidth int) string {
	if len(s) <= maxWidth {
		return s
	}
	if maxWidth <= 3 {
		return "..."
	}
	return s[:maxWidth-3] + "..."
}

// SetBursts sets the bursts to display
func (m *BurstListModel) SetBursts(bursts []*career.Burst) {
	m.bursts = bursts
	m.selectedIdx = 0
	m.expandedIndices = make(map[int]bool)
}

// GetSelectedBurst returns the currently selected burst
func (m *BurstListModel) GetSelectedBurst() *career.Burst {
	displayedBursts := m.getDisplayedBursts()
	if len(displayedBursts) == 0 || m.selectedIdx >= len(displayedBursts) {
		return nil
	}
	return m.bursts[displayedBursts[m.selectedIdx]]
}

// SetFilter sets the competency focus filter
func (m *BurstListModel) SetFilter(competency string) {
	m.filterBy = competency
	m.selectedIdx = 0
}

// SetSort sets the sort order
func (m *BurstListModel) SetSort(sortBy string) {
	m.sortBy = sortBy
}
