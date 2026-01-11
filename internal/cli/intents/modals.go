package intents

import (
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
// - Scrollable when content exceeds terminal height
type EditMetadataModal struct {
	// original is the unmodified metadata from the event (never mutated)
	original *MetadataSnapshot

	// modified is the working copy being edited
	modified *MetadataSnapshot

	// result is the final result returned to parent
	result *ModalEditResult[*MetadataSnapshot]

	// form is the huh form for editing
	form *huh.Form

	// formGroup stores the form group for rebuilding with new height
	formGroup *huh.Group

	// formData holds the form field values (pointers so form binding works)
	company         *string
	project         *string
	tags            []string // MultiSelect uses slice directly
	categories      []string // MultiSelect uses slice directly
	submitConfirmed *bool

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
	submitConfirmed := false

	// Copy slices for multi-select binding (ensure not nil)
	tagsCopy := make([]string, len(tags))
	copy(tagsCopy, tags)
	categoriesCopy := make([]string, len(categories))
	copy(categoriesCopy, categories)

	// Build tag options from AllowedTags
	tagOptions := make([]huh.Option[string], 0)
	for tag := range career.AllowedTags {
		tagOptions = append(tagOptions, huh.NewOption(tag, tag))
	}

	// Build category options from AllowedCategories
	categoryOptions := make([]huh.Option[string], 0)
	for cat := range career.AllowedCategories {
		categoryOptions = append(categoryOptions, huh.NewOption(cat, cat))
	}

	modal := &EditMetadataModal{
		original:        original,
		modified:        copyMetadataSnapshot(original),
		result:          nil,
		form:            nil, // Will be set below
		formGroup:       nil, // Will be set below
		company:         &companyVal,
		project:         &projectVal,
		tags:            tagsCopy,
		categories:      categoriesCopy,
		submitConfirmed: &submitConfirmed,
		width:           80,
		height:          24,
	}

	// Create form group with pointers to modal's fields
	modal.formGroup = huh.NewGroup(
		forms.NewInput(forms.FieldConfig{
			Key:         "company",
			Title:       "Company",
			Description: "Company name",
			Placeholder: "Enter company name...",
			CharLimit:   100,
			Validate:    forms.CompanyName,
		}).Value(modal.company),

		forms.NewInput(forms.FieldConfig{
			Key:         "project",
			Title:       "Project",
			Description: "Project name",
			Placeholder: "Enter project name...",
			CharLimit:   100,
		}).Value(modal.project),

		huh.NewMultiSelect[string]().
			Key("tags").
			Title("Tags").
			Description("Select relevant tags").
			Options(tagOptions...).
			Value(&modal.tags).
			Limit(8),

		huh.NewMultiSelect[string]().
			Key("categories").
			Title("Categories").
			Description("Select relevant categories").
			Options(categoryOptions...).
			Value(&modal.categories).
			Limit(6),

		huh.NewConfirm().
			Key("submit").
			Title("Save Changes").
			Description("Submit the form to save your changes").
			Affirmative("Submit").
			Negative("Cancel").
			Value(modal.submitConfirmed),
	)

	// Create form with default dimensions (will be updated on WindowSizeMsg)
	modal.form = forms.NewFormWithDimensions(
		modal.width-4, // Leave margin for modal chrome
		forms.DefaultFormHeight(modal.height),
		modal.formGroup,
	)

	return modal
}

// Update handles user input for metadata editing.
func (m *EditMetadataModal) Update(msg tea.Msg) tea.Cmd {
	switch msg := msg.(type) {
	case tea.WindowSizeMsg:
		m.width = msg.Width
		m.height = msg.Height
		// Update form dimensions without losing state
		m.form = m.form.
			WithHeight(forms.DefaultFormHeight(m.height)).
			WithWidth(m.width - 4) // Leave margin for modal chrome
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
		SetInstructions("Tab: Next  |  Shift+Tab: Prev  |  Enter: Confirm  |  Esc: Cancel").
		WithWidth(m.width - 4). // Use terminal width minus margin
		WithScrollHint(true)    // Show scroll indicator

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
		Company:    *m.company,
		Project:    *m.project,
		Tags:       m.tags,       // MultiSelect binds directly to []string
		Categories: m.categories, // MultiSelect binds directly to []string
	}
}

func (m *EditMetadataModal) createResult() {
	// Check if user confirmed via the submit button
	// If they selected "Cancel" on the confirm, treat as cancelled
	if !*m.submitConfirmed {
		m.createCancelledResult()
		return
	}

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
// Scrollable when content exceeds terminal height.
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

	// Default dimensions
	defaultWidth := 80
	defaultHeight := 24

	// Create huh form with default dimensions (will be updated on WindowSizeMsg)
	form := forms.NewBurstEditorFormWithDataAndDimensions(
		formData,
		defaultWidth-4, // Leave margin for modal chrome
		forms.DefaultFormHeight(defaultHeight),
	)

	return &EditBurstModal{
		original: &originalCopy,
		modified: &originalCopy,
		result:   nil,
		form:     form,
		formData: formData,
		width:    defaultWidth,
		height:   defaultHeight,
	}
}

// Update handles user input for burst editing.
func (m *EditBurstModal) Update(msg tea.Msg) tea.Cmd {
	switch msg := msg.(type) {
	case tea.WindowSizeMsg:
		m.width = msg.Width
		m.height = msg.Height
		// Update form dimensions without losing state
		m.form = m.form.
			WithHeight(forms.DefaultFormHeight(m.height)).
			WithWidth(m.width - 4) // Leave margin for modal chrome
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
		SetInstructions("Tab: Next  |  Shift+Tab: Prev  |  Enter: Confirm  |  Esc: Cancel").
		WithWidth(m.width - 4). // Use terminal width minus margin
		WithScrollHint(true)    // Show scroll indicator

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

// SetTestResult sets the result directly for testing purposes.
// This allows tests to simulate modal completion without interacting with the huh form.
func (m *EditBurstModal) SetTestResult(result *ModalEditResult[*career.Burst]) {
	m.result = result
	if result != nil && result.Modified != nil {
		m.modified = result.Modified
	}
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
	// Check if user confirmed via the submit button
	// If they selected "Cancel" on the confirm, treat as cancelled
	if !m.formData.SubmitConfirmed {
		m.createCancelledResult()
		return
	}

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
// Scrollable when content exceeds terminal height.
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

	// Default dimensions
	defaultWidth := 80
	defaultHeight := 24

	// Create huh form with default dimensions (will be updated on WindowSizeMsg)
	form := forms.NewFactEditorFormWithDataAndDimensions(
		formData,
		defaultWidth-4, // Leave margin for modal chrome
		forms.DefaultFormHeight(defaultHeight),
	)

	return &EditFactModal{
		original: &originalCopy,
		modified: &originalCopy,
		result:   nil,
		form:     form,
		formData: formData,
		width:    defaultWidth,
		height:   defaultHeight,
	}
}

// Update handles user input for fact editing.
func (m *EditFactModal) Update(msg tea.Msg) tea.Cmd {
	switch msg := msg.(type) {
	case tea.WindowSizeMsg:
		m.width = msg.Width
		m.height = msg.Height
		// Update form dimensions without losing state
		m.form = m.form.
			WithHeight(forms.DefaultFormHeight(m.height)).
			WithWidth(m.width - 4) // Leave margin for modal chrome
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
		SetInstructions("Tab: Next  |  Shift+Tab: Prev  |  Enter: Confirm  |  Esc: Cancel").
		WithWidth(m.width - 4). // Use terminal width minus margin
		WithScrollHint(true)    // Show scroll indicator

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
	// Note: StrengthSignal is preserved from original as it's auto-generated
	m.modified = &career.Fact{
		ID:                   m.original.ID,
		Text:                 m.formData.Text,
		CompetencyCategories: m.formData.CompetencyCategories,
		RoleFit:              career.RoleFit(m.formData.RoleFit),
		AudienceRelevance:    m.formData.AudienceRelevance,
		StrengthSignal:       m.original.StrengthSignal, // Preserved from original
		SourceEventID:        m.original.SourceEventID,
		SourceBurstID:        m.original.SourceBurstID,
		CreatedAt:            m.original.CreatedAt,
		UpdatedAt:            m.original.UpdatedAt,
	}
}

func (m *EditFactModal) createResult() {
	// Check if user confirmed via the submit button
	// If they selected "Cancel" on the confirm, treat as cancelled
	if !m.formData.SubmitConfirmed {
		m.createCancelledResult()
		return
	}

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
