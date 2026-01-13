package components

import (
	"github.com/baphled/kariya/internal/domain/career"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/huh"
)

// FilterModalModel manages the filter modal form for timeline filtering.
// This is a reusable component for filtering lists of career events by company, category, and sort options.
type FilterModalModel struct {
	form     *huh.Form
	formData *FilterFormData
	visible  bool
	width    int
	height   int
}

// FilterFormData holds the form field values for timeline filtering.
type FilterFormData struct {
	Companies  []string
	Categories []string
	SortBy     string
	SortOrder  string
}

// TimelineFilters represents the current filter state (matches intents.TimelineFilters).
// This is defined here to avoid circular dependency.
type TimelineFilters struct {
	SearchText string
	Tags       []string
	Companies  []string
	Categories []string
	SortBy     string
	SortOrder  string
}

// NewFilterModal creates a new filter modal.
// events: list of career events to extract filter options from
// currentFilters: current filter state to pre-populate form
// width, height: terminal dimensions for responsive sizing
func NewFilterModal(events []*career.CareerEvent, currentFilters *TimelineFilters, width, height int) *FilterModalModel {
	formData := &FilterFormData{
		SortBy:    "date",
		SortOrder: "desc",
	}

	// Pre-populate from current filters
	if currentFilters != nil {
		formData.Companies = currentFilters.Companies
		formData.Categories = currentFilters.Categories
		if currentFilters.SortBy != "" {
			formData.SortBy = currentFilters.SortBy
		}
		if currentFilters.SortOrder != "" {
			formData.SortOrder = currentFilters.SortOrder
		}
	}

	modal := &FilterModalModel{
		formData: formData,
		visible:  true,
		width:    width,
		height:   height,
	}

	modal.buildForm(events)
	return modal
}

// buildForm creates the huh form with filter options
func (m *FilterModalModel) buildForm(events []*career.CareerEvent) {
	// Extract unique companies from events
	companyMap := make(map[string]bool)
	for _, evt := range events {
		if evt.Company != "" {
			companyMap[evt.Company] = true
		}
	}
	companyOptions := make([]huh.Option[string], 0, len(companyMap))
	for company := range companyMap {
		companyOptions = append(companyOptions, huh.NewOption(company, company))
	}

	// Extract unique categories from events
	categoryMap := make(map[string]bool)
	for _, evt := range events {
		for _, cat := range evt.Categories {
			categoryMap[cat] = true
		}
	}
	categoryOptions := make([]huh.Option[string], 0, len(categoryMap))
	for category := range categoryMap {
		categoryOptions = append(categoryOptions, huh.NewOption(category, category))
	}

	// Create form fields
	fields := []huh.Field{}

	// Only add company filter if there are companies
	if len(companyOptions) > 0 {
		fields = append(fields, huh.NewMultiSelect[string]().
			Title("Filter by Company").
			Options(companyOptions...).
			Value(&m.formData.Companies))
	}

	// Only add category filter if there are categories
	if len(categoryOptions) > 0 {
		fields = append(fields, huh.NewMultiSelect[string]().
			Title("Filter by Category").
			Options(categoryOptions...).
			Value(&m.formData.Categories))
	}

	// Add sort options
	fields = append(fields,
		huh.NewSelect[string]().
			Title("Sort By").
			Options(
				huh.NewOption("Date", "date"),
				huh.NewOption("Text", "text"),
			).
			Value(&m.formData.SortBy),

		huh.NewSelect[string]().
			Title("Sort Order").
			Options(
				huh.NewOption("Newest First", "desc"),
				huh.NewOption("Oldest First", "asc"),
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

	// Calculate maximum form height to fit within terminal
	// Total overhead:
	// - Logo: 6 lines
	// - Logo spacing: 2 lines
	// - Footer: 4 lines
	// - Modal chrome (borders, padding, title, footer): 8 lines
	// - Safety margins: 4 lines
	// Total: 24 lines overhead
	const overhead = 24
	maxFormHeight := m.height - overhead
	if maxFormHeight < 8 {
		maxFormHeight = 8 // Minimum usable height for filter form
	}

	// Filter form is relatively small, but should still respect terminal constraints
	formHeight := 12
	if formHeight > maxFormHeight {
		formHeight = maxFormHeight
	}

	m.form = huh.NewForm(group).
		WithWidth(modalWidth).
		WithHeight(formHeight)
}

// Init initializes the filter modal and its form.
// This must be called to start the form's lifecycle.
func (m *FilterModalModel) Init() tea.Cmd {
	if m.form == nil {
		return nil
	}
	return m.form.Init()
}

// Update handles messages for the filter modal
func (m *FilterModalModel) Update(msg tea.Msg) (tea.Cmd, bool, *FilterFormData) {
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

// View renders the filter modal
func (m *FilterModalModel) View() string {
	if !m.visible {
		return ""
	}
	return m.form.View()
}

// IsVisible returns whether the modal is currently visible
func (m *FilterModalModel) IsVisible() bool {
	return m.visible
}

// Show makes the modal visible
func (m *FilterModalModel) Show() {
	m.visible = true
}

// Hide hides the modal
func (m *FilterModalModel) Hide() {
	m.visible = false
}
