package models

import (
	"context"
	"fmt"
	"sort"
	"strings"

	"github.com/baphled/kariya/internal/cli/styles"
	"github.com/baphled/kariya/internal/domain/career"
	careerservice "github.com/baphled/kariya/internal/service/career"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
)

// FactListModel displays a list of facts with filtering, sorting, and selection
type FactListModel struct {
	*BaseStandardModel
	// Data
	facts    []*career.Fact
	filtered []*career.Fact // Cached filtered results
	service  *careerservice.Service
	ctx      context.Context

	// UI State
	selectedIdx  int
	focusedIdx   int
	scrollOffset int
	width        int
	height       int

	// Filtering
	competencyFilter string
	roleFitFilter    career.RoleFit
	audienceFilter   string

	// Sorting
	sortBy    string // "date", "relevance"
	sortOrder string // "asc", "desc"

	// Selection
	selectedFacts map[string]bool // Track selected fact IDs for bulk operations

	// Messages
	submitted bool
	cancelled bool
	err       error
}

// NewFactListModel creates a new fact list model
func NewFactListModel(service *careerservice.Service, ctx context.Context) *FactListModel {
	return &FactListModel{
		BaseStandardModel: NewBaseStandardModel(),
		service:           service,
		ctx:               ctx,
		facts:             []*career.Fact{},
		filtered:          []*career.Fact{},
		selectedFacts:     make(map[string]bool),
		width:             80,
		height:            20,
		sortBy:            "date",
		sortOrder:         "desc",
	}
}

// Init initializes the model
func (flm *FactListModel) Init() tea.Cmd {
	return nil
}

// Update handles messages
func (flm *FactListModel) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch msg.String() {
		case "up", "k":
			if flm.focusedIdx > 0 {
				flm.focusedIdx--
				if flm.focusedIdx < flm.scrollOffset {
					flm.scrollOffset = flm.focusedIdx
				}
			}

		case "down", "j":
			if flm.focusedIdx < len(flm.filtered)-1 {
				flm.focusedIdx++
				maxIdx := flm.height - 3 // Account for header and footer
				if flm.focusedIdx >= flm.scrollOffset+maxIdx {
					flm.scrollOffset = flm.focusedIdx - maxIdx + 1
				}
			}

		case "enter":
			if len(flm.filtered) > 0 {
				flm.selectedIdx = flm.focusedIdx
				flm.submitted = true
			}

		case "space":
			if len(flm.filtered) > 0 {
				fact := flm.filtered[flm.focusedIdx]
				flm.selectedFacts[fact.ID] = !flm.selectedFacts[fact.ID]
			}

		case "escape", "q":
			flm.cancelled = true

		case "f":
			// Toggle filter mode (placeholder for filter UI)
			break

		case "s":
			// Toggle sort mode (placeholder for sort UI)
			break
		}

	case tea.WindowSizeMsg:
		flm.width = msg.Width
		flm.height = msg.Height
	}

	return flm, nil
}

// View renders the fact list
func (flm *FactListModel) View() string {
	if len(flm.filtered) == 0 {
		return flm.renderEmpty()
	}

	var parts []string

	// Header
	parts = append(parts, flm.renderHeader())

	// Facts list
	parts = append(parts, flm.renderFacts())

	// Footer with keyboard shortcuts
	parts = append(parts, flm.renderFooter())

	return lipgloss.JoinVertical(lipgloss.Left, parts...)
}

// renderHeader renders the header with title and filters
func (flm *FactListModel) renderHeader() string {
	title := "📋 Facts"
	if flm.competencyFilter != "" {
		title += fmt.Sprintf(" (Competency: %s)", flm.competencyFilter)
	}
	if flm.roleFitFilter != "" {
		title += fmt.Sprintf(" (Role: %s)", flm.roleFitFilter)
	}

	headerStyle := lipgloss.NewStyle().
		Foreground(styles.ColorTextPrimary).
		Bold(true).
		MarginBottom(1)

	info := fmt.Sprintf("%d/%d facts", len(flm.filtered), len(flm.facts))
	if len(flm.selectedFacts) > 0 {
		info += fmt.Sprintf(" | %d selected", len(flm.selectedFacts))
	}

	infoStyle := lipgloss.NewStyle().
		Foreground(styles.ColorTextMuted).
		Italic(true)

	return lipgloss.JoinVertical(
		lipgloss.Left,
		headerStyle.Render(title),
		infoStyle.Render(info),
	)
}

// renderFacts renders the list of facts
func (flm *FactListModel) renderFacts() string {
	if len(flm.filtered) == 0 {
		return ""
	}

	maxIdx := flm.height - 5 // Account for header and footer
	endIdx := flm.scrollOffset + maxIdx
	if endIdx > len(flm.filtered) {
		endIdx = len(flm.filtered)
	}

	var items []string
	for i := flm.scrollOffset; i < endIdx && i < len(flm.filtered); i++ {
		fact := flm.filtered[i]
		isFocused := (i == flm.focusedIdx)
		isSelected := flm.selectedFacts[fact.ID]
		item := flm.renderFactItem(fact, isFocused, isSelected)
		items = append(items, item)
	}

	return lipgloss.JoinVertical(lipgloss.Left, items...)
}

// renderFactItem renders a single fact item in the list
func (flm *FactListModel) renderFactItem(fact *career.Fact, focused, selected bool) string {
	// Checkbox
	checkbox := "☐"
	if selected {
		checkbox = "☑"
	}

	// Focus indicator
	focusIndicator := " "
	if focused {
		focusIndicator = "►"
	}

	// Fact preview (first 50 chars)
	preview := fact.Text
	if len(preview) > 50 {
		preview = preview[:47] + "..."
	}

	// Role fit icon
	roleIcon := getRoleFitIcon(fact.RoleFit)

	// Format: [focus] [checkbox] [role-icon] [competency] text...
	line := fmt.Sprintf("%s %s %s [%s] %s",
		focusIndicator,
		checkbox,
		roleIcon,
		strings.Join(fact.CompetencyCategories, ","),
		preview,
	)

	style := lipgloss.NewStyle().
		Foreground(styles.ColorTextSecondary).
		Width(flm.width).
		Padding(0)

	if focused {
		style = style.
			Foreground(styles.ColorTextPrimary).
			Background(styles.ColorBackgroundCard).
			Bold(true)
	}

	if selected {
		style = style.Foreground(styles.ColorAccentTeal)
	}

	return style.Render(line)
}

// renderFooter renders the footer with keyboard shortcuts
func (flm *FactListModel) renderFooter() string {
	shortcuts := []string{
		"↑/j: Up",
		"↓/k: Down",
		"Enter: View",
		"Space: Select",
		"f: Filter",
		"s: Sort",
		"Esc: Back",
	}

	footerText := strings.Join(shortcuts, " │ ")
	footerStyle := lipgloss.NewStyle().
		Foreground(styles.ColorTextMuted).
		Italic(true).
		MarginTop(1)

	return footerStyle.Render(footerText)
}

// renderEmpty renders the empty state
func (flm *FactListModel) renderEmpty() string {
	message := "No facts found"
	if flm.competencyFilter != "" || flm.roleFitFilter != "" || flm.audienceFilter != "" {
		message = "No facts match the current filters"
	}

	emptyStyle := lipgloss.NewStyle().
		Foreground(styles.ColorTextMuted).
		AlignHorizontal(lipgloss.Center).
		Padding(2, 0)

	return emptyStyle.Render(message)
}

// SetFacts sets the facts to display
func (flm *FactListModel) SetFacts(facts []*career.Fact) {
	flm.facts = facts
	flm.applyFiltersAndSort()
	flm.focusedIdx = 0
	flm.scrollOffset = 0
}

// SetCompetencyFilter sets the competency filter
func (flm *FactListModel) SetCompetencyFilter(competency string) {
	flm.competencyFilter = competency
	flm.applyFiltersAndSort()
	flm.focusedIdx = 0
	flm.scrollOffset = 0
}

// SetRoleFitFilter sets the role fit filter
func (flm *FactListModel) SetRoleFitFilter(roleFit career.RoleFit) {
	flm.roleFitFilter = roleFit
	flm.applyFiltersAndSort()
	flm.focusedIdx = 0
	flm.scrollOffset = 0
}

// SetAudienceFilter sets the audience filter
func (flm *FactListModel) SetAudienceFilter(audience string) {
	flm.audienceFilter = audience
	flm.applyFiltersAndSort()
	flm.focusedIdx = 0
	flm.scrollOffset = 0
}

// SetSort sets the sort order
func (flm *FactListModel) SetSort(sortBy, sortOrder string) {
	flm.sortBy = sortBy
	flm.sortOrder = sortOrder
	flm.applyFiltersAndSort()
}

// applyFiltersAndSort applies filters and sorts the facts
func (flm *FactListModel) applyFiltersAndSort() {
	flm.filtered = flm.filterFacts()
	flm.sortFacts(flm.filtered)
}

// filterFacts filters facts based on current filters
func (flm *FactListModel) filterFacts() []*career.Fact {
	var filtered []*career.Fact

	for _, fact := range flm.facts {
		// Competency filter
		if flm.competencyFilter != "" {
			hasCompetency := false
			for _, comp := range fact.CompetencyCategories {
				if strings.EqualFold(comp, flm.competencyFilter) {
					hasCompetency = true
					break
				}
			}
			if !hasCompetency {
				continue
			}
		}

		// Role fit filter
		if flm.roleFitFilter != "" && fact.RoleFit != flm.roleFitFilter {
			continue
		}

		// Audience filter
		if flm.audienceFilter != "" {
			hasAudience := false
			for _, aud := range fact.AudienceRelevance {
				if strings.EqualFold(aud, flm.audienceFilter) {
					hasAudience = true
					break
				}
			}
			if !hasAudience {
				continue
			}
		}

		filtered = append(filtered, fact)
	}

	return filtered
}

// sortFacts sorts facts based on current sort settings
func (flm *FactListModel) sortFacts(facts []*career.Fact) {
	sort.SliceStable(facts, func(i, j int) bool {
		var less bool

		switch flm.sortBy {
		case "date":
			less = facts[i].CreatedAt.Before(facts[j].CreatedAt)
		case "relevance":
			// Sort by number of competencies and audience matches (higher = more relevant)
			compScore := len(facts[i].CompetencyCategories) - len(facts[j].CompetencyCategories)
			if compScore != 0 {
				less = compScore < 0
			} else {
				audScore := len(facts[i].AudienceRelevance) - len(facts[j].AudienceRelevance)
				less = audScore < 0
			}
		default:
			less = facts[i].CreatedAt.Before(facts[j].CreatedAt)
		}

		// Reverse for descending order
		if flm.sortOrder == "asc" {
			return less
		}
		return !less
	})
}

// GetSelectedFact returns the currently selected fact
func (flm *FactListModel) GetSelectedFact() *career.Fact {
	if len(flm.filtered) > 0 && flm.focusedIdx < len(flm.filtered) {
		return flm.filtered[flm.focusedIdx]
	}
	return nil
}

// GetSelectedFacts returns all selected facts
func (flm *FactListModel) GetSelectedFacts() []*career.Fact {
	var selected []*career.Fact
	for _, fact := range flm.filtered {
		if flm.selectedFacts[fact.ID] {
			selected = append(selected, fact)
		}
	}
	return selected
}

// ClearSelection clears all selections
func (flm *FactListModel) ClearSelection() {
	flm.selectedFacts = make(map[string]bool)
}

// IsSubmitted returns true if a fact was selected
func (flm *FactListModel) IsSubmitted() bool {
	return flm.submitted
}

// IsCancelled returns true if the user cancelled
func (flm *FactListModel) IsCancelled() bool {
	return flm.cancelled
}

// GetError returns the last error
func (flm *FactListModel) GetError() error {
	return flm.err
}

// GetFacts returns the currently filtered facts
func (flm *FactListModel) GetFacts() []*career.Fact {
	return flm.filtered
}

// GetSelectedIdx returns the currently selected index
func (flm *FactListModel) GetSelectedIdx() int {
	return flm.selectedIdx
}
