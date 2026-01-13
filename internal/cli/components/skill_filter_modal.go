package components

import (
	"fmt"
	"strconv"

	"github.com/baphled/kariya/internal/cli/styles"
	"github.com/baphled/kariya/internal/domain/career"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/huh"
	"github.com/charmbracelet/lipgloss"
)

// SkillFilters represents the current filter state for skills.
// NOTE: Search is handled by separate SkillSearchModal (accessed via `/` key)
type SkillFilters struct {
	Categories []string
	Levels     []string
	MinYears   int
	MaxYears   int
	SortBy     string
	SortOrder  string
}

// SkillFilterFormData holds the form field values for skill filtering.
// NOTE: Search is handled by separate SkillSearchModal
type SkillFilterFormData struct {
	Categories  []string
	Levels      []string
	MinYearsStr string
	MaxYearsStr string
	SortBy      string
	SortOrder   string
}

// SkillFilterModal manages the filter modal form for skill filtering.
type SkillFilterModal struct {
	form     *huh.Form
	formData *SkillFilterFormData
	visible  bool
	width    int
	height   int
}

// NewSkillFilterModal creates a new skill filter modal.
// NOTE: Search is handled by separate SkillSearchModal
func NewSkillFilterModal(skills []*career.Skill, currentFilter *SkillFilters, width, height int) *SkillFilterModal {
	formData := &SkillFilterFormData{
		SortBy:    "name",
		SortOrder: "asc",
	}

	// Pre-populate from current filters
	if currentFilter != nil {
		formData.Categories = currentFilter.Categories
		formData.Levels = currentFilter.Levels
		if currentFilter.MinYears > 0 {
			formData.MinYearsStr = strconv.Itoa(currentFilter.MinYears)
		}
		if currentFilter.MaxYears > 0 {
			formData.MaxYearsStr = strconv.Itoa(currentFilter.MaxYears)
		}
		if currentFilter.SortBy != "" {
			formData.SortBy = currentFilter.SortBy
		}
		if currentFilter.SortOrder != "" {
			formData.SortOrder = currentFilter.SortOrder
		}
	}

	modal := &SkillFilterModal{
		formData: formData,
		visible:  true,
		width:    width,
		height:   height,
	}

	modal.buildForm(skills)
	return modal
}

// buildForm creates the huh form with filter options
func (m *SkillFilterModal) buildForm(skills []*career.Skill) {
	// Extract unique categories from skills
	categoryMap := make(map[string]bool)
	for _, skill := range skills {
		if skill.Category != "" {
			categoryMap[skill.Category] = true
		}
	}
	categoryOptions := make([]huh.Option[string], 0, len(categoryMap))
	for category := range categoryMap {
		categoryOptions = append(categoryOptions, huh.NewOption(category, category))
	}

	// Extract unique levels from skills
	levelMap := make(map[string]bool)
	for _, skill := range skills {
		if skill.Level != "" {
			levelMap[skill.Level] = true
		}
	}
	levelOptions := make([]huh.Option[string], 0, len(levelMap))
	for level := range levelMap {
		levelOptions = append(levelOptions, huh.NewOption(level, level))
	}

	// Create form fields
	// NOTE: Search is handled by separate SkillSearchModal (accessed via `/` key)
	fields := []huh.Field{}

	// Only add category filter if there are categories
	if len(categoryOptions) > 0 {
		fields = append(fields, huh.NewMultiSelect[string]().
			Title("Filter by Category").
			Options(categoryOptions...).
			Value(&m.formData.Categories))
	}

	// Only add level filter if there are levels
	if len(levelOptions) > 0 {
		fields = append(fields, huh.NewMultiSelect[string]().
			Title("Filter by Level").
			Options(levelOptions...).
			Value(&m.formData.Levels))
	}

	// Years range filters
	fields = append(fields,
		huh.NewInput().
			Title("Minimum Years").
			Placeholder("0").
			Value(&m.formData.MinYearsStr).
			Validate(func(s string) error {
				if s == "" {
					return nil
				}
				years, err := strconv.Atoi(s)
				if err != nil {
					return fmt.Errorf("must be a number")
				}
				if years < 0 {
					return fmt.Errorf("must be positive")
				}
				return nil
			}),

		huh.NewInput().
			Title("Maximum Years").
			Placeholder("100").
			Value(&m.formData.MaxYearsStr).
			Validate(func(s string) error {
				if s == "" {
					return nil
				}
				years, err := strconv.Atoi(s)
				if err != nil {
					return fmt.Errorf("must be a number")
				}
				if years < 0 {
					return fmt.Errorf("must be positive")
				}
				return nil
			}),
	)

	// Add sort options
	fields = append(fields,
		huh.NewSelect[string]().
			Title("Sort By").
			Options(
				huh.NewOption("Name", "name"),
				huh.NewOption("Category", "category"),
				huh.NewOption("Level", "level"),
				huh.NewOption("Years of Experience", "years"),
				huh.NewOption("Events Count", "events"),
			).
			Value(&m.formData.SortBy),

		huh.NewSelect[string]().
			Title("Sort Order").
			Options(
				huh.NewOption("Ascending", "asc"),
				huh.NewOption("Descending", "desc"),
			).
			Value(&m.formData.SortOrder),
	)

	group := huh.NewGroup(fields...)

	// Calculate modal width (60% of screen, max 80 chars)
	modalWidth := m.width * 60 / 100
	if modalWidth > 80 {
		modalWidth = 80
	}
	if modalWidth < 40 {
		modalWidth = 40
	}

	// Let Huh use natural height - bubbletea-overlay will handle positioning
	m.form = huh.NewForm(group).
		WithWidth(modalWidth)
}

// Init initializes the filter modal and its form.
func (m *SkillFilterModal) Init() tea.Cmd {
	if m.form == nil {
		return nil
	}
	return m.form.Init()
}

// Update handles messages for the filter modal
func (m *SkillFilterModal) Update(msg tea.Msg) (tea.Cmd, bool, *SkillFilterFormData) {
	if !m.visible {
		return nil, false, nil
	}

	switch msg := msg.(type) {
	case tea.WindowSizeMsg:
		m.width = msg.Width
		m.height = msg.Height
		return nil, false, nil

	case tea.KeyMsg:
		switch msg.String() {
		case "esc":
			// Close modal without applying
			m.visible = false
			return nil, false, nil
		}
	}

	// Update form
	form, cmd := m.form.Update(msg)
	m.form = form.(*huh.Form)

	// Check if form is complete
	if m.form.State == huh.StateCompleted {
		m.visible = false
		return cmd, true, m.formData
	}

	return cmd, false, nil
}

// View renders the filter modal with proper chrome (border, background)
func (m *SkillFilterModal) View() string {
	if !m.visible {
		return ""
	}

	// Wrap the form in a styled box with solid background, border, and padding
	return lipgloss.NewStyle().
		Border(lipgloss.RoundedBorder()).
		BorderForeground(styles.ColorBorder).
		Background(styles.ColorBackground).
		Padding(1, 2).
		Render(m.form.View())
}

// IsVisible returns whether the modal is currently visible.
func (m *SkillFilterModal) IsVisible() bool {
	return m.visible
}

// Show makes the modal visible.
func (m *SkillFilterModal) Show() {
	m.visible = true
}

// Hide hides the modal.
func (m *SkillFilterModal) Hide() {
	m.visible = false
}

// ToSkillFilters converts form data to SkillFilters
// NOTE: Search is handled separately by SkillSearchModal
func (m *SkillFilterModal) ToSkillFilters() *SkillFilters {
	filters := &SkillFilters{
		Categories: m.formData.Categories,
		Levels:     m.formData.Levels,
		SortBy:     m.formData.SortBy,
		SortOrder:  m.formData.SortOrder,
	}

	// Parse years (ignore errors, default to 0)
	if m.formData.MinYearsStr != "" {
		if years, err := strconv.Atoi(m.formData.MinYearsStr); err == nil && years > 0 {
			filters.MinYears = years
		}
	}
	if m.formData.MaxYearsStr != "" {
		if years, err := strconv.Atoi(m.formData.MaxYearsStr); err == nil && years > 0 {
			filters.MaxYears = years
		}
	}

	return filters
}
