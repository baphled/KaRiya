package modals

import (
	"fmt"
	"strconv"
	"strings"

	"github.com/baphled/kariya/internal/cli/forms"
	"github.com/baphled/kariya/internal/cli/themes"
	"github.com/baphled/kariya/internal/cli/uikit/containers"
	"github.com/baphled/kariya/internal/cli/uikit/primitives"
	"github.com/baphled/kariya/internal/domain/career"
	tea "github.com/charmbracelet/bubbletea"
)

// Filters represents the current filter state for skills.
// NOTE: Search is handled by SearchModal (`/` key)
// NOTE: Sort is handled by SortModal (`s` key)
type Filters struct {
	Categories []string
	Levels     []string
	MinYears   int
	MaxYears   int
}

// FilterFormData holds the form field values for skill filtering.
// NOTE: Search and Sort are handled by separate modals
type FilterFormData struct {
	Categories  []string
	Levels      []string
	MinYearsStr string
	MaxYearsStr string
}

// FilterModal manages the filter modal form for skill filtering.
type FilterModal struct {
	form     forms.Form
	formData *FilterFormData
	visible  bool
	width    int
	height   int
	theme    themes.Theme
}

// NewFilterModal creates a new skill filter modal.
// NOTE: Search and Sort are handled by separate modals
func NewFilterModal(skills []*career.Skill, currentFilter *Filters, width, height int) *FilterModal {
	formData := initFilterFormData(currentFilter)

	modal := &FilterModal{
		formData: formData,
		visible:  true,
		width:    width,
		height:   height,
		theme:    themes.NewDefaultTheme(),
	}

	modal.buildForm(skills)
	return modal
}

// initFilterFormData creates form data from current filters.
func initFilterFormData(currentFilter *Filters) *FilterFormData {
	formData := &FilterFormData{}
	if currentFilter != nil {
		formData.Categories = currentFilter.Categories
		formData.Levels = currentFilter.Levels
		if currentFilter.MinYears > 0 {
			formData.MinYearsStr = strconv.Itoa(currentFilter.MinYears)
		}
		if currentFilter.MaxYears > 0 {
			formData.MaxYearsStr = strconv.Itoa(currentFilter.MaxYears)
		}
	}
	return formData
}

// extractUniqueValues extracts unique non-empty values from skills using the given accessor.
func extractUniqueValues(skills []*career.Skill, accessor func(*career.Skill) string) []forms.SelectOption {
	valueMap := make(map[string]bool)
	for _, skill := range skills {
		val := accessor(skill)
		if val != "" {
			valueMap[val] = true
		}
	}

	options := make([]forms.SelectOption, 0, len(valueMap))
	for val := range valueMap {
		options = append(options, forms.SelectOption{Key: val, Value: val})
	}
	return options
}

// buildForm creates the huh form with filter options.
func (m *FilterModal) buildForm(skills []*career.Skill) {
	// Extract unique categories and levels.
	categoryOptions := extractUniqueValues(skills, func(s *career.Skill) string { return s.Category })
	levelOptions := extractUniqueValues(skills, func(s *career.Skill) string { return s.Level })

	// Build form fields.
	fields := m.buildFilterFields(categoryOptions, levelOptions)

	// Create form group.
	group := forms.NewGroup(fields...)

	// Calculate modal width.
	modalWidth := m.calculateModalWidth()

	// Create form with dimensions.
	m.form = forms.NewFormWithDimensions(modalWidth, 0, group)
}

// buildFilterFields creates the filter form fields.
func (m *FilterModal) buildFilterFields(categoryOptions, levelOptions []forms.SelectOption) []forms.Field {
	fields := []forms.Field{}

	// Category filter (if categories exist).
	if len(categoryOptions) > 0 {
		fields = append(fields, forms.NewMultiSelect(
			"categories",
			"Filter by Category",
			"",
			categoryOptions,
			0,
		).Value(&m.formData.Categories))
	}

	// Level filter (if levels exist).
	if len(levelOptions) > 0 {
		fields = append(fields, forms.NewMultiSelect(
			"levels",
			"Filter by Level",
			"",
			levelOptions,
			0,
		).Value(&m.formData.Levels))
	}

	// Years range filters.
	fields = append(fields, m.buildYearsInput("Minimum Years", "0", &m.formData.MinYearsStr))
	fields = append(fields, m.buildYearsInput("Maximum Years", "100", &m.formData.MaxYearsStr))

	return fields
}

// buildYearsInput creates a validated years input field.
func (m *FilterModal) buildYearsInput(title, placeholder string, value *string) forms.Input {
	return forms.NewInput(forms.FieldConfig{
		Title:       title,
		Placeholder: placeholder,
		Validate: func(s string) error {
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
		},
	}).Value(value)
}

// calculateModalWidth calculates the appropriate modal width.
func (m *FilterModal) calculateModalWidth() int {
	modalWidth := m.width * 60 / 100
	if modalWidth > 80 {
		modalWidth = 80
	}
	if modalWidth < 40 {
		modalWidth = 40
	}
	return modalWidth
}

// Init initializes the filter modal and its form.
func (m *FilterModal) Init() tea.Cmd {
	if m.form == nil {
		return nil
	}
	return m.form.Init()
}

// Update handles messages for the filter modal.
func (m *FilterModal) Update(msg tea.Msg) (tea.Cmd, bool, *FilterFormData) {
	if !m.visible {
		return nil, false, nil
	}

	switch msg := msg.(type) {
	case tea.WindowSizeMsg:
		m.width = msg.Width
		m.height = msg.Height
		return nil, false, nil

	case tea.KeyMsg:
		if msg.Type == tea.KeyEsc {
			m.visible = false
			return nil, false, nil
		}
	}

	// Update form using forms package helper.
	var cmd tea.Cmd
	m.form, cmd = forms.Update(m.form, msg)

	// Check if form is complete.
	if forms.IsCompleted(m.form) {
		m.visible = false
		return cmd, true, m.formData
	}

	return cmd, false, nil
}

// View renders the filter modal with proper chrome (border, background).
func (m *FilterModal) View() string {
	if !m.visible {
		return ""
	}

	// Build footer with primitives showing keyboard shortcuts.
	footer := primitives.RenderHelpFooter(m.theme,
		primitives.NextFieldBadge(m.theme),
		primitives.ApplyBadge(m.theme),
		primitives.CancelBadge(m.theme),
	)

	// Build modal content with form and footer.
	var content strings.Builder
	content.WriteString(m.form.View())
	content.WriteString("\n\n")
	content.WriteString(footer)

	// Wrap the form in a styled box with solid background using UIKit.
	theme := m.theme
	if theme == nil {
		theme = themes.NewDefaultTheme()
	}
	return containers.NewBox(theme).
		Content(content.String()).
		Padding(1).
		Background(theme.BackgroundColor()).
		Render()
}

// IsVisible returns whether the modal is currently visible.
func (m *FilterModal) IsVisible() bool {
	return m.visible
}

// Show makes the modal visible.
func (m *FilterModal) Show() {
	m.visible = true
}

// Hide hides the modal.
func (m *FilterModal) Hide() {
	m.visible = false
}

// ToFilters converts form data to Filters.
// NOTE: Search and Sort are handled separately by other modals
func (m *FilterModal) ToFilters() *Filters {
	filters := &Filters{
		Categories: m.formData.Categories,
		Levels:     m.formData.Levels,
	}

	// Parse years (ignore errors, default to 0).
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

// RenderOverlay renders the filter modal as an overlay on top of the base view.
// This follows the bubbletea-overlay pattern used in BrowseTimeline modals.
func (m *FilterModal) RenderOverlay(baseView string) string {
	if !m.visible {
		return baseView
	}

	return RenderOverlayModal(m.View(), baseView)
}
