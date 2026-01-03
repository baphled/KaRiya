package intents

import (
	"fmt"
	"strings"
	"unicode"

	"github.com/baphled/kariya/internal/cli/components"
	"github.com/baphled/kariya/internal/cli/styles"
	"github.com/baphled/kariya/internal/domain/career"
	"github.com/charmbracelet/bubbles/textinput"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
)

// EditMetadataModal handles inline editing of event metadata (Company, Project, Tags, Categories).
// This is a sub-flow of CaptureEventIntent used in the ReviewInferredEvent state.
//
// The modal follows the ModalEditResult[T] pattern:
// - Preserves original data (no mutation)
// - Returns typed diff on confirmation
// - Restores original on cancellation
// - Uses lipgloss for professional styling
// - Uses bubbles textinput for interactive fields
// - Fully responsive to terminal size changes
type EditMetadataModal struct {
	// original is the unmodified metadata from the event (never mutated)
	original *MetadataSnapshot

	// modified is the working copy being edited
	modified *MetadataSnapshot

	// result is the final result returned to parent
	result *ModalEditResult[*MetadataSnapshot]

	// inputs are the bubbles textinput components for each field
	inputs [4]textinput.Model

	// focused tracks which field is currently focused (0=Company, 1=Project, 2=Tags, 3=Categories)
	focused int

	// accepted tracks if the user confirmed changes
	accepted bool

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

// NewEditMetadataModal creates a new metadata editing modal.
// original: the original metadata to display
// Returns a new modal ready for interaction.
func NewEditMetadataModal(company, project string, tags, categories []string) *EditMetadataModal {
	original := &MetadataSnapshot{
		Company:    company,
		Project:    project,
		Tags:       tags,
		Categories: categories,
	}

	m := &EditMetadataModal{
		original: original,
		modified: copyMetadataSnapshot(original),
		focused:  0,
		accepted: false,
		result:   nil,
		width:    80,  // Default width
		height:   24,  // Default height
	}

	// Initialize bubbles textinput components for each field
	m.initializeInputs()

	return m
}

// initializeInputs creates and configures bubbles textinput components for each field.
// This ensures consistent input handling and styling across all fields.
func (m *EditMetadataModal) initializeInputs() {
	fieldValues := []string{
		m.modified.Company,
		m.modified.Project,
		formatStringSlice(m.modified.Tags),
		formatStringSlice(m.modified.Categories),
	}
	fieldHints := []string{
		"Enter company name",
		"Enter project name",
		"Comma-separated tags",
		"Comma-separated categories",
	}

	for i := 0; i < 4; i++ {
		input := textinput.New()
		input.SetValue(fieldValues[i])
		input.Placeholder = fieldHints[i]
		input.CharLimit = 256

		// Focus the first field by default
		if i == 0 {
			input.Focus()
		}

		m.inputs[i] = input
	}
}

// Update handles user input for metadata editing.
// Delegates to bubbles textinput components for text input.
// Handles navigation (Tab, Shift+Tab) and confirmation (Enter, Escape).
func (m *EditMetadataModal) Update(msg tea.Msg) tea.Cmd {
	switch msg := msg.(type) {
	case tea.WindowSizeMsg:
		// Handle terminal resize for responsive layout
		m.width = msg.Width
		m.height = msg.Height
		return nil

	case tea.KeyMsg:
		switch msg.Type {
		case tea.KeyTab:
			// Move to next field with visual feedback
			m.inputs[m.focused].Blur()
			m.focused = (m.focused + 1) % 4
			m.inputs[m.focused].Focus()
			return nil

		case tea.KeyShiftTab:
			// Move to previous field with visual feedback
			m.inputs[m.focused].Blur()
			m.focused = (m.focused - 1 + 4) % 4
			m.inputs[m.focused].Focus()
			return nil

		case tea.KeyEnter:
			// Confirm changes - sync input values to modified data
			m.syncModified()
			m.accepted = true
			m.createResult()
			return nil

		case tea.KeyEsc:
			// Cancel changes - discard all edits
			m.accepted = false
			m.createCancelledResult()
			return nil

		default:
			// Delegate other keys to focused input field
			_, cmd := m.inputs[m.focused].Update(msg)
			return cmd
		}

	default:
		return nil
	}
}

// View renders the metadata editing modal with professional styling.
// Uses lipgloss for colors and layout, bubbles textinput for interactive fields.
// Responsive to terminal width and height.
func (m *EditMetadataModal) View() string {
	if m.result != nil && m.result.Accepted {
		return "" // Modal is done, return empty string
	}

	// Calculate responsive width (min 40, max 80, adapt to terminal)
	modalWidth := m.width - 4 // Leave margins
	if modalWidth < 40 {
		modalWidth = 40
	}
	if modalWidth > 80 {
		modalWidth = 80
	}

	// Field definitions with labels and inputs
	fieldLabels := []string{"Company", "Project", "Tags", "Categories"}

	// Build form fields with labels and inputs
	var fields []string
	for i, label := range fieldLabels {
		fields = append(fields, m.renderField(label, i, modalWidth))
	}

	// Combine all parts with consistent spacing
	title := styles.ModalTitle.
		Foreground(styles.ColorTextPrimary).
		Render("Edit Event Metadata")

	content := lipgloss.JoinVertical(
		lipgloss.Left,
		title,
		"",
		lipgloss.JoinVertical(lipgloss.Left, fields...),
		"",
		m.renderInstructions(),
	)

	// Use ModalContainer for consistent styling
	modal := components.NewModalContainer().
		SetTitle("").
		SetMessage(content).
		SetInstructions("Enter: Confirm  |  Esc: Cancel  |  Tab: Next Field  |  Shift+Tab: Previous")

	return modal.Render()
}

// renderField renders a single form field with label and input.
// Uses bubbles textinput for consistent interactive behavior.
// Shows focus indicator with visual feedback.
func (m *EditMetadataModal) renderField(label string, index int, width int) string {
	// Calculate field width (account for label and spacing)
	fieldWidth := width - len(label) - 10

	// Render label with focus indicator
	focusIndicator := "  "
	if m.focused == index {
		focusIndicator = "→ " // Visual focus indicator (not color-only)
	}

	labelStyle := styles.InputLabel
	if m.focused == index {
		labelStyle = labelStyle.Bold(true).Foreground(styles.ColorAccentTeal)
	}

	labelText := fmt.Sprintf("%s%s:", focusIndicator, label)
	renderedLabel := labelStyle.Render(labelText)

	// Render input field with focus styling
	inputStyle := styles.InputBase
	if m.focused == index {
		inputStyle = styles.InputFocused
	}

	inputValue := m.inputs[index].Value()
	if len(inputValue) > fieldWidth {
		inputValue = inputValue[:fieldWidth-3] + "..."
	}

	renderedInput := inputStyle.Width(fieldWidth).Render(inputValue)

	return lipgloss.JoinVertical(lipgloss.Left, renderedLabel, renderedInput)
}

// renderInstructions renders keyboard shortcut instructions with accessibility in mind.
// Uses text labels (not color-only) for clarity.
func (m *EditMetadataModal) renderInstructions() string {
	instructions := []string{
		"✓ Enter: Confirm changes",
		"✗ Esc: Cancel",
		"→ Tab: Next field",
		"← Shift+Tab: Previous field",
	}

	return styles.ModalInstructions.Render(strings.Join(instructions, "  |  "))
}

// Result returns the modal result when editing is complete.
// Returns nil if modal is still active.
func (m *EditMetadataModal) Result() *ModalEditResult[*MetadataSnapshot] {
	return m.result
}

// IsComplete returns true if the modal has finished (accepted or cancelled).
func (m *EditMetadataModal) IsComplete() bool {
	return m.result != nil
}

// Private helper methods

func (m *EditMetadataModal) syncModified() {
	m.modified = &MetadataSnapshot{
		Company:    m.inputs[0].Value(),
		Project:    m.inputs[1].Value(),
		Tags:       parseStringSlice(m.inputs[2].Value()),
		Categories: parseStringSlice(m.inputs[3].Value()),
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

// EditBurstModal handles inline editing of burst details (Name, Description, CompetencyFocus).
// Follows the same pattern as EditMetadataModal with professional styling and accessibility.
type EditBurstModal struct {
	// original is the unmodified burst (never mutated)
	original *career.Burst

	// modified is the working copy being edited
	modified *career.Burst

	// result is the final result returned to parent
	result *ModalEditResult[*career.Burst]

	// inputs are the bubbles textinput components for each field
	inputs [3]textinput.Model

	// focused tracks which field is currently focused
	focused int

	// accepted tracks if the user confirmed changes
	accepted bool

	// width and height track terminal dimensions for responsive layout
	width  int
	height int
}

// NewEditBurstModal creates a new burst editing modal.
// burst: the burst to edit
// Returns a new modal ready for interaction.
func NewEditBurstModal(burst *career.Burst) *EditBurstModal {
	originalCopy := *burst // Create a copy

	m := &EditBurstModal{
		original: &originalCopy,
		modified: &originalCopy,
		focused:  0,
		accepted: false,
		result:   nil,
		width:    80,
		height:   24,
	}

	m.initializeInputs()
	return m
}

// initializeInputs creates and configures bubbles textinput components for burst fields.
func (m *EditBurstModal) initializeInputs() {
	fieldValues := []string{
		m.modified.Name,
		m.modified.Description,
		m.modified.CompetencyFocus,
	}
	fieldHints := []string{
		"Enter burst name",
		"Enter description",
		"Enter competency focus",
	}

	for i := 0; i < 3; i++ {
		input := textinput.New()
		input.SetValue(fieldValues[i])
		input.Placeholder = fieldHints[i]
		input.CharLimit = 256

		if i == 0 {
			input.Focus()
		}

		m.inputs[i] = input
	}
}

// Update handles user input for burst editing.
func (m *EditBurstModal) Update(msg tea.Msg) tea.Cmd {
	switch msg := msg.(type) {
	case tea.WindowSizeMsg:
		m.width = msg.Width
		m.height = msg.Height
		return nil

	case tea.KeyMsg:
		switch msg.Type {
		case tea.KeyTab:
			m.inputs[m.focused].Blur()
			m.focused = (m.focused + 1) % 3
			m.inputs[m.focused].Focus()
			return nil

		case tea.KeyShiftTab:
			m.inputs[m.focused].Blur()
			m.focused = (m.focused - 1 + 3) % 3
			m.inputs[m.focused].Focus()
			return nil

		case tea.KeyEnter:
			m.syncModified()
			m.accepted = true
			m.createResult()
			return nil

		case tea.KeyEsc:
			m.accepted = false
			m.createCancelledResult()
			return nil

		default:
			_, cmd := m.inputs[m.focused].Update(msg)
			return cmd
		}

	default:
		return nil
	}
}

// View renders the burst editing modal with professional styling.
func (m *EditBurstModal) View() string {
	if m.result != nil && m.result.Accepted {
		return ""
	}

	// Calculate responsive width
	modalWidth := m.width - 4
	if modalWidth < 40 {
		modalWidth = 40
	}
	if modalWidth > 80 {
		modalWidth = 80
	}

	fieldLabels := []string{"Name", "Description", "Competency Focus"}

	var fields []string
	for i, label := range fieldLabels {
		fields = append(fields, m.renderField(label, i, modalWidth))
	}

	title := styles.ModalTitle.
		Foreground(styles.ColorTextPrimary).
		Render("Edit Burst")

	content := lipgloss.JoinVertical(
		lipgloss.Left,
		title,
		"",
		lipgloss.JoinVertical(lipgloss.Left, fields...),
		"",
		m.renderInstructions(),
	)

	modal := components.NewModalContainer().
		SetTitle("").
		SetMessage(content).
		SetInstructions("Enter: Confirm  |  Esc: Cancel  |  Tab: Next Field  |  Shift+Tab: Previous")

	return modal.Render()
}

// renderField renders a single burst field with label and input.
func (m *EditBurstModal) renderField(label string, index int, width int) string {
	fieldWidth := width - len(label) - 10

	focusIndicator := "  "
	if m.focused == index {
		focusIndicator = "→ "
	}

	labelStyle := styles.InputLabel
	if m.focused == index {
		labelStyle = labelStyle.Bold(true).Foreground(styles.ColorAccentTeal)
	}

	labelText := fmt.Sprintf("%s%s:", focusIndicator, label)
	renderedLabel := labelStyle.Render(labelText)

	inputStyle := styles.InputBase
	if m.focused == index {
		inputStyle = styles.InputFocused
	}

	inputValue := m.inputs[index].Value()
	if len(inputValue) > fieldWidth {
		inputValue = inputValue[:fieldWidth-3] + "..."
	}

	renderedInput := inputStyle.Width(fieldWidth).Render(inputValue)

	return lipgloss.JoinVertical(lipgloss.Left, renderedLabel, renderedInput)
}

// renderInstructions renders keyboard shortcut instructions.
func (m *EditBurstModal) renderInstructions() string {
	instructions := []string{
		"✓ Enter: Confirm changes",
		"✗ Esc: Cancel",
		"→ Tab: Next field",
		"← Shift+Tab: Previous field",
	}

	return styles.ModalInstructions.Render(strings.Join(instructions, "  |  "))
}

// Result returns the modal result when editing is complete.
func (m *EditBurstModal) Result() *ModalEditResult[*career.Burst] {
	return m.result
}

// IsComplete returns true if the modal has finished.
func (m *EditBurstModal) IsComplete() bool {
	return m.result != nil
}

func (m *EditBurstModal) syncModified() {
	m.modified = &career.Burst{
		ID:              m.original.ID,
		Name:            m.inputs[0].Value(),
		Description:     m.inputs[1].Value(),
		CompetencyFocus: m.inputs[2].Value(),
		EventIDs:        m.original.EventIDs,
		CreatedAt:       m.original.CreatedAt,
		UpdatedAt:       m.original.UpdatedAt,
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
	if m.original.CompetencyFocus != m.modified.CompetencyFocus {
		changes["competency_focus"] = m.modified.CompetencyFocus
	}

	return changes
}

// ============================================================================
// EditFactModal
// ============================================================================

// EditFactModal handles inline editing of fact details (Text, CompetencyCategories, RoleFit, AudienceRelevance, StrengthSignal).
// Follows the same pattern as EditMetadataModal with professional styling and accessibility.
type EditFactModal struct {
	// original is the unmodified fact (never mutated)
	original *career.Fact

	// modified is the working copy being edited
	modified *career.Fact

	// result is the final result returned to parent
	result *ModalEditResult[*career.Fact]

	// inputs are the bubbles textinput components for each field
	inputs [5]textinput.Model

	// focused tracks which field is currently focused
	focused int

	// accepted tracks if the user confirmed changes
	accepted bool

	// width and height track terminal dimensions for responsive layout
	width  int
	height int
}

// NewEditFactModal creates a new fact editing modal.
// fact: the fact to edit
// Returns a new modal ready for interaction.
func NewEditFactModal(fact *career.Fact) *EditFactModal {
	originalCopy := *fact // Create a copy

	m := &EditFactModal{
		original: &originalCopy,
		modified: &originalCopy,
		focused:  0,
		accepted: false,
		result:   nil,
		width:    80,
		height:   24,
	}

	m.initializeInputs()
	return m
}

// initializeInputs creates and configures bubbles textinput components for fact fields.
func (m *EditFactModal) initializeInputs() {
	fieldValues := []string{
		m.modified.Text,
		formatStringSlice(m.modified.CompetencyCategories),
		string(m.modified.RoleFit),
		formatStringSlice(m.modified.AudienceRelevance),
		m.modified.StrengthSignal,
	}
	fieldHints := []string{
		"Enter fact text",
		"Comma-separated categories",
		"Enter role fit",
		"Comma-separated audiences",
		"Enter strength signal",
	}

	for i := 0; i < 5; i++ {
		input := textinput.New()
		input.SetValue(fieldValues[i])
		input.Placeholder = fieldHints[i]
		input.CharLimit = 256

		if i == 0 {
			input.Focus()
		}

		m.inputs[i] = input
	}
}

// Update handles user input for fact editing.
func (m *EditFactModal) Update(msg tea.Msg) tea.Cmd {
	switch msg := msg.(type) {
	case tea.WindowSizeMsg:
		m.width = msg.Width
		m.height = msg.Height
		return nil

	case tea.KeyMsg:
		switch msg.Type {
		case tea.KeyTab:
			m.inputs[m.focused].Blur()
			m.focused = (m.focused + 1) % 5
			m.inputs[m.focused].Focus()
			return nil

		case tea.KeyShiftTab:
			m.inputs[m.focused].Blur()
			m.focused = (m.focused - 1 + 5) % 5
			m.inputs[m.focused].Focus()
			return nil

		case tea.KeyEnter:
			m.syncModified()
			m.accepted = true
			m.createResult()
			return nil

		case tea.KeyEsc:
			m.accepted = false
			m.createCancelledResult()
			return nil

		default:
			_, cmd := m.inputs[m.focused].Update(msg)
			return cmd
		}

	default:
		return nil
	}
}

// View renders the fact editing modal with professional styling.
func (m *EditFactModal) View() string {
	if m.result != nil && m.result.Accepted {
		return ""
	}

	// Calculate responsive width
	modalWidth := m.width - 4
	if modalWidth < 40 {
		modalWidth = 40
	}
	if modalWidth > 80 {
		modalWidth = 80
	}

	fieldLabels := []string{"Text", "Categories", "RoleFit", "Audience", "Strength Signal"}

	var fields []string
	for i, label := range fieldLabels {
		fields = append(fields, m.renderField(label, i, modalWidth))
	}

	title := styles.ModalTitle.
		Foreground(styles.ColorTextPrimary).
		Render("Edit Fact")

	content := lipgloss.JoinVertical(
		lipgloss.Left,
		title,
		"",
		lipgloss.JoinVertical(lipgloss.Left, fields...),
		"",
		m.renderInstructions(),
	)

	modal := components.NewModalContainer().
		SetTitle("").
		SetMessage(content).
		SetInstructions("Enter: Confirm  |  Esc: Cancel  |  Tab: Next Field  |  Shift+Tab: Previous")

	return modal.Render()
}

// renderField renders a single fact field with label and input.
func (m *EditFactModal) renderField(label string, index int, width int) string {
	fieldWidth := width - len(label) - 10

	focusIndicator := "  "
	if m.focused == index {
		focusIndicator = "→ "
	}

	labelStyle := styles.InputLabel
	if m.focused == index {
		labelStyle = labelStyle.Bold(true).Foreground(styles.ColorAccentTeal)
	}

	labelText := fmt.Sprintf("%s%s:", focusIndicator, label)
	renderedLabel := labelStyle.Render(labelText)

	inputStyle := styles.InputBase
	if m.focused == index {
		inputStyle = styles.InputFocused
	}

	inputValue := m.inputs[index].Value()
	if len(inputValue) > fieldWidth {
		inputValue = inputValue[:fieldWidth-3] + "..."
	}

	renderedInput := inputStyle.Width(fieldWidth).Render(inputValue)

	return lipgloss.JoinVertical(lipgloss.Left, renderedLabel, renderedInput)
}

// renderInstructions renders keyboard shortcut instructions.
func (m *EditFactModal) renderInstructions() string {
	instructions := []string{
		"✓ Enter: Confirm changes",
		"✗ Esc: Cancel",
		"→ Tab: Next field",
		"← Shift+Tab: Previous field",
	}

	return styles.ModalInstructions.Render(strings.Join(instructions, "  |  "))
}

// Result returns the modal result when editing is complete.
func (m *EditFactModal) Result() *ModalEditResult[*career.Fact] {
	return m.result
}

// IsComplete returns true if the modal has finished.
func (m *EditFactModal) IsComplete() bool {
	return m.result != nil
}

func (m *EditFactModal) syncModified() {
	m.modified = &career.Fact{
		ID:                   m.original.ID,
		Text:                 m.inputs[0].Value(),
		CompetencyCategories: parseStringSlice(m.inputs[1].Value()),
		RoleFit:              career.RoleFit(m.inputs[2].Value()),
		AudienceRelevance:    parseStringSlice(m.inputs[3].Value()),
		StrengthSignal:       m.inputs[4].Value(),
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

// truncateString truncates a string to a maximum length.
func truncateString(s string, maxLen int) string {
	if len(s) <= maxLen {
		return s
	}
	return s[:maxLen-3] + "..."
}
