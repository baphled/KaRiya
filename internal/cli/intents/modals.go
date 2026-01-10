package intents

import (
	"strings"
	"unicode"

	"github.com/baphled/kariya/internal/cli/components"
	"github.com/baphled/kariya/internal/cli/forms"
	"github.com/baphled/kariya/internal/cli/styles"
	"github.com/baphled/kariya/internal/domain/career"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/huh"
	"github.com/charmbracelet/lipgloss"
)

// EditMetadataModal handles inline editing of event metadata (Company, Project, Tags, Categories).
// This is a sub-flow of CaptureEventIntent used in the ReviewInferredEvent state.
//
// The modal follows the ModalEditResult[T] pattern:
// - Preserves original data (no mutation)
// - Returns typed diff on confirmation
// - Restores original on cancellation
// - Uses huh library for form handling
// - Professional styling with Catppuccin theme
type EditMetadataModal struct {
	// original is the unmodified metadata from the event (never mutated)
	original *MetadataSnapshot

	// modified is the working copy being edited
	modified *MetadataSnapshot

	// result is the final result returned to parent
	result *ModalEditResult[*MetadataSnapshot]

	// form is the huh form for editing
	form *huh.Form

	// formData holds the form field values
	company    string
	project    string
	tags       string
	categories string

	// width and height track terminal dimensions for responsive layout
	width  int
	height int
}

// MetadataSnapshot represents a snapshot of event metadata for editing.
type MetadataSnapshot struct {
	Company    string
	Project    string
	Tags       []string
	Categories []string
}

// NewEditMetadataModal creates a new metadata editing modal using huh forms.
// original: the original metadata to display
// Returns a new modal ready for interaction.
func NewEditMetadataModal(company, project string, tags, categories []string) *EditMetadataModal {
	original := &MetadataSnapshot{
		Company:    company,
		Project:    project,
		Tags:       tags,
		Categories: categories,
	}

	// Initialize form field values
	companyVal := company
	projectVal := project
	tagsVal := formatStringSlice(tags)
	categoriesVal := formatStringSlice(categories)

	// Create huh form
	form := forms.NewForm(
		huh.NewGroup(
			forms.NewInput(forms.FieldConfig{
				Key:         "company",
				Title:       "Company",
				Description: "Company name",
				Placeholder: "Enter company name...",
				CharLimit:   100,
				Validate:    forms.CompanyName,
			}).Value(&companyVal),

			forms.NewInput(forms.FieldConfig{
				Key:         "project",
				Title:       "Project",
				Description: "Project name",
				Placeholder: "Enter project name...",
				CharLimit:   100,
			}).Value(&projectVal),

			huh.NewInput().
				Key("tags").
				Title("Tags").
				Description("Comma-separated tags").
				Placeholder("tag1, tag2, tag3").
				CharLimit(256).
				Value(&tagsVal),

			huh.NewInput().
				Key("categories").
				Title("Categories").
				Description("Comma-separated categories").
				Placeholder("category1, category2").
				CharLimit(256).
				Value(&categoriesVal),
		),
	)

	return &EditMetadataModal{
		original:   original,
		modified:   copyMetadataSnapshot(original),
		result:     nil,
		form:       form,
		company:    companyVal,
		project:    projectVal,
		tags:       tagsVal,
		categories: categoriesVal,
		width:      80,
		height:     24,
	}
}

// Update handles user input for metadata editing.
func (m *EditMetadataModal) Update(msg tea.Msg) tea.Cmd {
	switch msg := msg.(type) {
	case tea.WindowSizeMsg:
		m.width = msg.Width
		m.height = msg.Height
		return nil
	}

	// Update the form
	form, cmd := m.form.Update(msg)
	if f, ok := form.(*huh.Form); ok {
		m.form = f
	}

	// Check form state
	if forms.IsCompleted(m.form) {
		m.syncModified()
		m.createResult()
		return nil
	}

	if forms.IsAborted(m.form) {
		m.createCancelledResult()
		return nil
	}

	return cmd
}

// View renders the metadata editing modal with professional styling.
func (m *EditMetadataModal) View() string {
	if m.result != nil && m.result.Accepted {
		return ""
	}

	// Render the form
	formView := m.form.View()

	// Wrap in modal container
	title := styles.ModalTitle.
		Foreground(styles.ColorTextPrimary).
		Render("Edit Event Metadata")

	content := lipgloss.JoinVertical(
		lipgloss.Left,
		title,
		"",
		formView,
	)

	modal := components.NewModalContainer().
		SetTitle("").
		SetMessage(content).
		SetInstructions("Enter: Confirm  |  Esc: Cancel  |  Tab: Next Field  |  Shift+Tab: Previous")

	return modal.Render()
}

// Result returns the modal result when editing is complete.
func (m *EditMetadataModal) Result() *ModalEditResult[*MetadataSnapshot] {
	return m.result
}

// IsComplete returns true if the modal has finished (accepted or cancelled).
func (m *EditMetadataModal) IsComplete() bool {
	return m.result != nil
}

// GetTitle returns the modal title for overlay rendering.
func (m *EditMetadataModal) GetTitle() string {
	return "Edit Event Metadata"
}

// GetContent returns just the form content without the modal container.
// This allows parent intents to compose the modal as an overlay.
func (m *EditMetadataModal) GetContent() string {
	if m.result != nil && m.result.Accepted {
		return ""
	}
	return m.form.View()
}

// GetFooter returns the footer instructions for the modal.
func (m *EditMetadataModal) GetFooter() string {
	return "Enter: Confirm  |  Esc: Cancel  |  Tab: Next Field  |  Shift+Tab: Previous"
}

// Private helper methods

func (m *EditMetadataModal) syncModified() {
	m.modified = &MetadataSnapshot{
		Company:    m.company,
		Project:    m.project,
		Tags:       parseStringSlice(m.tags),
		Categories: parseStringSlice(m.categories),
	}
}

func (m *EditMetadataModal) createResult() {
	m.result = &ModalEditResult[*MetadataSnapshot]{
		Original: m.original,
		Modified: m.modified,
		Accepted: true,
		Changes:  m.computeChanges(),
	}
}

func (m *EditMetadataModal) createCancelledResult() {
	m.result = &ModalEditResult[*MetadataSnapshot]{
		Original: m.original,
		Modified: m.original,
		Accepted: false,
		Changes:  make(map[string]interface{}),
	}
}

func (m *EditMetadataModal) computeChanges() map[string]interface{} {
	changes := make(map[string]interface{})

	if m.original.Company != m.modified.Company {
		changes["company"] = m.modified.Company
	}
	if m.original.Project != m.modified.Project {
		changes["project"] = m.modified.Project
	}
	if !slicesEqual(m.original.Tags, m.modified.Tags) {
		changes["tags"] = m.modified.Tags
	}
	if !slicesEqual(m.original.Categories, m.modified.Categories) {
		changes["categories"] = m.modified.Categories
	}

	return changes
}

// ============================================================================
// EditBurstModal
// ============================================================================

// EditBurstModal handles inline editing of burst details (Name, Description).
// Uses huh library for form handling with professional styling and accessibility.
type EditBurstModal struct {
	// original is the unmodified burst (never mutated)
	original *career.Burst

	// modified is the working copy being edited
	modified *career.Burst

	// result is the final result returned to parent
	result *ModalEditResult[*career.Burst]

	// form is the huh form for editing
	form *huh.Form

	// formData holds the form field values
	formData *forms.BurstFormData

	// width and height track terminal dimensions for responsive layout
	width  int
	height int
}

// NewEditBurstModal creates a new burst editing modal using huh forms.
// burst: the burst to edit
// Returns a new modal ready for interaction.
func NewEditBurstModal(burst *career.Burst) *EditBurstModal {
	originalCopy := *burst // Create a copy

	// Create form data from burst
	formData := forms.GetBurstFormData(&originalCopy)

	// Create huh form
	form := forms.NewBurstEditorFormWithData(formData)

	return &EditBurstModal{
		original: &originalCopy,
		modified: &originalCopy,
		result:   nil,
		form:     form,
		formData: formData,
		width:    80,
		height:   24,
	}
}

// Update handles user input for burst editing.
func (m *EditBurstModal) Update(msg tea.Msg) tea.Cmd {
	switch msg := msg.(type) {
	case tea.WindowSizeMsg:
		m.width = msg.Width
		m.height = msg.Height
		return nil
	}

	// Update the form
	form, cmd := m.form.Update(msg)
	if f, ok := form.(*huh.Form); ok {
		m.form = f
	}

	// Check form state
	if forms.IsCompleted(m.form) {
		m.syncModified()
		m.createResult()
		return nil
	}

	if forms.IsAborted(m.form) {
		m.createCancelledResult()
		return nil
	}

	return cmd
}

// View renders the burst editing modal with professional styling.
func (m *EditBurstModal) View() string {
	if m.result != nil && m.result.Accepted {
		return ""
	}

	// Render the form
	formView := m.form.View()

	// Wrap in modal container
	title := styles.ModalTitle.
		Foreground(styles.ColorTextPrimary).
		Render("Edit Burst")

	content := lipgloss.JoinVertical(
		lipgloss.Left,
		title,
		"",
		formView,
	)

	modal := components.NewModalContainer().
		SetTitle("").
		SetMessage(content).
		SetInstructions("Enter: Confirm  |  Esc: Cancel  |  Tab: Next Field  |  Shift+Tab: Previous")

	return modal.Render()
}

// Result returns the modal result when editing is complete.
func (m *EditBurstModal) Result() *ModalEditResult[*career.Burst] {
	return m.result
}

// IsComplete returns true if the modal has finished.
func (m *EditBurstModal) IsComplete() bool {
	return m.result != nil
}

// GetTitle returns the modal title for overlay rendering.
func (m *EditBurstModal) GetTitle() string {
	return "Edit Burst"
}

// GetContent returns just the form content without the modal container.
// This allows parent intents to compose the modal as an overlay.
func (m *EditBurstModal) GetContent() string {
	if m.result != nil && m.result.Accepted {
		return ""
	}
	return m.form.View()
}

// GetFooter returns the footer instructions for the modal.
func (m *EditBurstModal) GetFooter() string {
	return "Enter: Confirm  |  Esc: Cancel  |  Tab: Next Field  |  Shift+Tab: Previous"
}

func (m *EditBurstModal) syncModified() {
	m.modified = &career.Burst{
		ID:          m.original.ID,
		Name:        m.formData.Name,
		Description: m.formData.Description,
		EventIDs:    m.original.EventIDs,
		CreatedAt:   m.original.CreatedAt,
		UpdatedAt:   m.original.UpdatedAt,
	}
}

func (m *EditBurstModal) createResult() {
	m.result = &ModalEditResult[*career.Burst]{
		Original: m.original,
		Modified: m.modified,
		Accepted: true,
		Changes:  m.computeChanges(),
	}
}

func (m *EditBurstModal) createCancelledResult() {
	m.result = &ModalEditResult[*career.Burst]{
		Original: m.original,
		Modified: m.original,
		Accepted: false,
		Changes:  make(map[string]interface{}),
	}
}

func (m *EditBurstModal) computeChanges() map[string]interface{} {
	changes := make(map[string]interface{})

	if m.original.Name != m.modified.Name {
		changes["name"] = m.modified.Name
	}
	if m.original.Description != m.modified.Description {
		changes["description"] = m.modified.Description
	}

	return changes
}

// ============================================================================
// EditFactModal
// ============================================================================

// EditFactModal handles inline editing of fact details (Text, CompetencyCategories, RoleFit, AudienceRelevance, StrengthSignal).
// Uses huh library for form handling with professional styling and accessibility.
type EditFactModal struct {
	// original is the unmodified fact (never mutated)
	original *career.Fact

	// modified is the working copy being edited
	modified *career.Fact

	// result is the final result returned to parent
	result *ModalEditResult[*career.Fact]

	// form is the huh form for editing
	form *huh.Form

	// formData holds the form field values
	formData *forms.FactFormData

	// width and height track terminal dimensions for responsive layout
	width  int
	height int
}

// NewEditFactModal creates a new fact editing modal using huh forms.
// fact: the fact to edit
// Returns a new modal ready for interaction.
func NewEditFactModal(fact *career.Fact) *EditFactModal {
	originalCopy := *fact // Create a copy

	// Create form data from fact
	formData := forms.GetFactFormData(&originalCopy)

	// Create huh form
	form := forms.NewFactEditorFormWithData(formData)

	return &EditFactModal{
		original: &originalCopy,
		modified: &originalCopy,
		result:   nil,
		form:     form,
		formData: formData,
		width:    80,
		height:   24,
	}
}

// Update handles user input for fact editing.
func (m *EditFactModal) Update(msg tea.Msg) tea.Cmd {
	switch msg := msg.(type) {
	case tea.WindowSizeMsg:
		m.width = msg.Width
		m.height = msg.Height
		return nil
	}

	// Update the form
	form, cmd := m.form.Update(msg)
	if f, ok := form.(*huh.Form); ok {
		m.form = f
	}

	// Check form state
	if forms.IsCompleted(m.form) {
		m.syncModified()
		m.createResult()
		return nil
	}

	if forms.IsAborted(m.form) {
		m.createCancelledResult()
		return nil
	}

	return cmd
}

// View renders the fact editing modal with professional styling.
func (m *EditFactModal) View() string {
	if m.result != nil && m.result.Accepted {
		return ""
	}

	// Render the form
	formView := m.form.View()

	// Wrap in modal container
	title := styles.ModalTitle.
		Foreground(styles.ColorTextPrimary).
		Render("Edit Fact")

	content := lipgloss.JoinVertical(
		lipgloss.Left,
		title,
		"",
		formView,
	)

	modal := components.NewModalContainer().
		SetTitle("").
		SetMessage(content).
		SetInstructions("Enter: Confirm  |  Esc: Cancel  |  Tab: Next Field  |  Shift+Tab: Previous")

	return modal.Render()
}

// Result returns the modal result when editing is complete.
func (m *EditFactModal) Result() *ModalEditResult[*career.Fact] {
	return m.result
}

// IsComplete returns true if the modal has finished.
func (m *EditFactModal) IsComplete() bool {
	return m.result != nil
}

// GetTitle returns the modal title for overlay rendering.
func (m *EditFactModal) GetTitle() string {
	return "Edit Fact"
}

// GetContent returns just the form content without the modal container.
// This allows parent intents to compose the modal as an overlay.
func (m *EditFactModal) GetContent() string {
	if m.result != nil && m.result.Accepted {
		return ""
	}
	return m.form.View()
}

// GetFooter returns the footer instructions for the modal.
func (m *EditFactModal) GetFooter() string {
	return "Enter: Confirm  |  Esc: Cancel  |  Tab: Next Field  |  Shift+Tab: Previous"
}

func (m *EditFactModal) syncModified() {
	// Apply form data to modified fact
	m.modified = &career.Fact{
		ID:                   m.original.ID,
		Text:                 m.formData.Text,
		CompetencyCategories: parseStringSlice(m.formData.CompetencyCategories),
		RoleFit:              career.RoleFit(m.formData.RoleFit),
		AudienceRelevance:    parseStringSlice(m.formData.AudienceRelevance),
		StrengthSignal:       m.formData.StrengthSignal,
		SourceEventID:        m.original.SourceEventID,
		SourceBurstID:        m.original.SourceBurstID,
		CreatedAt:            m.original.CreatedAt,
		UpdatedAt:            m.original.UpdatedAt,
	}
}

func (m *EditFactModal) createResult() {
	m.result = &ModalEditResult[*career.Fact]{
		Original: m.original,
		Modified: m.modified,
		Accepted: true,
		Changes:  m.computeChanges(),
	}
}

func (m *EditFactModal) createCancelledResult() {
	m.result = &ModalEditResult[*career.Fact]{
		Original: m.original,
		Modified: m.original,
		Accepted: false,
		Changes:  make(map[string]interface{}),
	}
}

func (m *EditFactModal) computeChanges() map[string]interface{} {
	changes := make(map[string]interface{})

	if m.original.Text != m.modified.Text {
		changes["text"] = m.modified.Text
	}
	if !slicesEqual(m.original.CompetencyCategories, m.modified.CompetencyCategories) {
		changes["competency_categories"] = m.modified.CompetencyCategories
	}
	if m.original.RoleFit != m.modified.RoleFit {
		changes["role_fit"] = m.modified.RoleFit
	}
	if !slicesEqual(m.original.AudienceRelevance, m.modified.AudienceRelevance) {
		changes["audience_relevance"] = m.modified.AudienceRelevance
	}
	if m.original.StrengthSignal != m.modified.StrengthSignal {
		changes["strength_signal"] = m.modified.StrengthSignal
	}

	return changes
}

// ============================================================================
// Helper Functions
// ============================================================================

// copyMetadataSnapshot creates a deep copy of a metadata snapshot.
func copyMetadataSnapshot(original *MetadataSnapshot) *MetadataSnapshot {
	if original == nil {
		return nil
	}

	copy := *original
	copy.Tags = make([]string, len(original.Tags))
	copy.Categories = make([]string, len(original.Categories))

	for i, tag := range original.Tags {
		copy.Tags[i] = tag
	}
	for i, cat := range original.Categories {
		copy.Categories[i] = cat
	}

	return &copy
}

// formatStringSlice converts a string slice to comma-separated format.
func formatStringSlice(items []string) string {
	if len(items) == 0 {
		return ""
	}

	result := ""
	for i, item := range items {
		if i > 0 {
			result += ", "
		}
		result += item
	}
	return result
}

// parseStringSlice converts comma-separated string to a slice.
func parseStringSlice(input string) []string {
	if input == "" {
		return []string{}
	}

	// Split by comma and trim spaces from each item
	parts := strings.Split(input, ",")
	var result []string

	for _, part := range parts {
		trimmed := strings.TrimFunc(part, unicode.IsSpace)
		if trimmed != "" {
			result = append(result, trimmed)
		}
	}

	return result
}

// slicesEqual checks if two string slices are equal.
func slicesEqual(a, b []string) bool {
	if len(a) != len(b) {
		return false
	}

	for i := range a {
		if a[i] != b[i] {
			return false
		}
	}

	return true
}
