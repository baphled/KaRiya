package models

import (
	"context"
	"fmt"
	"sort"
	"strings"

	"github.com/baphled/kariya/internal/cli/components"
	"github.com/baphled/kariya/internal/cli/styles"
	"github.com/baphled/kariya/internal/domain/career"
	careerservice "github.com/baphled/kariya/internal/service/career"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
)

// FactListModel displays a list of facts with filtering, sorting, and selection
type FactListModel struct {
	*BaseStandardModel
	facts            []*career.Fact
	filtered         []*career.Fact
	service          *careerservice.Service
	ctx              context.Context
	selectedIdx      int
	focusedIdx       int
	scrollOffset     int
	width            int
	height           int
	competencyFilter string
	roleFitFilter    career.RoleFit
	audienceFilter   string
	sortBy           string
	sortOrder        string
	selectedFacts    map[string]bool
	submitted        bool
	cancelled        bool
	err              error
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
			flm.prevItem()
		case "down", "j":
			flm.nextItem()
		case "pgup", "ctrl+b":
			flm.prevPage()
		case "pgdn", "ctrl+f":
			flm.nextPage()
		case "home", "g":
			flm.goToFirstItem()
		case "end", "G":
			flm.goToLastItem()
		case "enter":
			if len(flm.filtered) > 0 {
				flm.selectedIdx = flm.focusedIdx
				flm.submitted = true
			}
		case " ", "space":
			if len(flm.filtered) > 0 {
				fact := flm.filtered[flm.focusedIdx]
				flm.selectedFacts[fact.ID] = !flm.selectedFacts[fact.ID]
			}
		case "esc", "q":
			flm.cancelled = true
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

	var items []string
	maxIdx := flm.height - 5
	endIdx := flm.scrollOffset + maxIdx
	if endIdx > len(flm.filtered) {
		endIdx = len(flm.filtered)
	}

	for i := flm.scrollOffset; i < endIdx && i < len(flm.filtered); i++ {
		fact := flm.filtered[i]
		isFocused := (i == flm.focusedIdx)
		isSelected := flm.selectedFacts[fact.ID]
		item := flm.renderFactItem(fact, isFocused, isSelected)
		items = append(items, item)
	}

	paginationInfo := fmt.Sprintf("%d/%d facts", len(flm.filtered), len(flm.facts))
	if len(flm.selectedFacts) > 0 {
		paginationInfo += fmt.Sprintf(" | %d selected", len(flm.selectedFacts))
	}

	title := "📋 Facts"
	if flm.competencyFilter != "" {
		title += fmt.Sprintf(" (Competency: %s)", flm.competencyFilter)
	}
	if flm.roleFitFilter != "" {
		title += fmt.Sprintf(" (Role: %s)", flm.roleFitFilter)
	}

	headerView := components.NewHeader(title, flm.width).View()
	footerView := components.NewFooter(flm.width).View()

	listContainer := components.NewListContainer().
		SetItems(items).
		SetEmptyStateMessage("No facts found").
		SetPaginationInfo(paginationInfo)

	screenContent := lipgloss.JoinVertical(
		lipgloss.Left,
		listContainer.Render(),
	)

	screenContainer := components.NewScreenContainer(screenContent).
		WithPaddingMode(components.PaddingNormal)

	fullContent := lipgloss.JoinVertical(
		lipgloss.Left,
		headerView,
		"",
		screenContainer.Render(),
		"",
		footerView,
	)

	return fullContent
}

// renderFactItem renders a single fact item in the list
func (flm *FactListModel) renderFactItem(fact *career.Fact, focused, selected bool) string {
	checkbox := "☐"
	if selected {
		checkbox = "☑"
	}

	focusIndicator := " "
	if focused {
		focusIndicator = "►"
	}

	preview := fact.Text
	if len(preview) > 50 {
		preview = preview[:47] + "..."
	}

	roleIcon := getRoleFitIcon(fact.RoleFit)

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

// renderEmpty renders the empty state
func (flm *FactListModel) renderEmpty() string {
	message := "No facts found"
	if flm.competencyFilter != "" || flm.roleFitFilter != "" || flm.audienceFilter != "" {
		message = "No facts match the current filters"
	}

	listContainer := components.NewListContainer().
		SetItems([]string{}).
		SetEmptyStateMessage(message)

	headerView := components.NewHeader("📋 Facts", flm.width).View()
	footerView := components.NewFooter(flm.width).View()

	screenContent := lipgloss.JoinVertical(
		lipgloss.Left,
		listContainer.Render(),
	)

	screenContainer := components.NewScreenContainer(screenContent).
		WithPaddingMode(components.PaddingNormal)

	fullContent := lipgloss.JoinVertical(
		lipgloss.Left,
		headerView,
		"",
		screenContainer.Render(),
		"",
		footerView,
	)

	return fullContent
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

		if flm.roleFitFilter != "" && fact.RoleFit != flm.roleFitFilter {
			continue
		}

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

// nextItem moves to the next item in the list
func (flm *FactListModel) nextItem() {
	if flm.focusedIdx < len(flm.filtered)-1 {
		flm.focusedIdx++
		maxIdx := flm.height - 3
		if flm.focusedIdx >= flm.scrollOffset+maxIdx {
			flm.scrollOffset = flm.focusedIdx - maxIdx + 1
		}
	}
}

// prevItem moves to the previous item in the list
func (flm *FactListModel) prevItem() {
	if flm.focusedIdx > 0 {
		flm.focusedIdx--
		if flm.focusedIdx < flm.scrollOffset {
			flm.scrollOffset = flm.focusedIdx
		}
	}
}

// nextPage moves to the next page
func (flm *FactListModel) nextPage() {
	pageSize := flm.height - 3
	lastIdx := len(flm.filtered) - 1
	newIdx := flm.focusedIdx + pageSize
	if newIdx > lastIdx {
		newIdx = lastIdx
	}
	flm.focusedIdx = newIdx
	flm.scrollOffset = flm.focusedIdx - pageSize + 1
	if flm.scrollOffset < 0 {
		flm.scrollOffset = 0
	}
}

// prevPage moves to the previous page
func (flm *FactListModel) prevPage() {
	pageSize := flm.height - 3
	newIdx := flm.focusedIdx - pageSize
	if newIdx < 0 {
		newIdx = 0
	}
	flm.focusedIdx = newIdx
	flm.scrollOffset = flm.focusedIdx
}

// goToFirstItem moves to the first item in the list
func (flm *FactListModel) goToFirstItem() {
	flm.focusedIdx = 0
	flm.scrollOffset = 0
}

// goToLastItem moves to the last item in the list
func (flm *FactListModel) goToLastItem() {
	if len(flm.filtered) > 0 {
		flm.focusedIdx = len(flm.filtered) - 1
		pageSize := flm.height - 3
		flm.scrollOffset = flm.focusedIdx - pageSize + 1
		if flm.scrollOffset < 0 {
			flm.scrollOffset = 0
		}
	}
}
