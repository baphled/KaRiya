package models

import (
	"context"
	"fmt"
	"strings"
	"time"

	"github.com/baphled/kariya/internal/cli/components"
	"github.com/baphled/kariya/internal/cli/styles"
	"github.com/baphled/kariya/internal/domain/career"
	careerservice "github.com/baphled/kariya/internal/service/career"
	"github.com/charmbracelet/bubbles/textinput"
	tea "github.com/charmbracelet/bubbletea"
)

// Focus indices for fact editor navigation
const (
	FactTextFieldIdx = iota
	FactCompetenciesFieldIdx
	FactRoleFitFieldIdx
	FactAudienceFieldIdx
	FactSaveButtonIdx
	FactCancelButtonIdx
	FactEditorFieldCount // Total number of focus positions
)

// FactEditorModel represents the fact editor form state
type FactEditorModel struct {
	fact                    *career.Fact
	originalFact            *career.Fact // For reverting changes
	service                 *careerservice.Service
	ctx                     context.Context
	textInput               textinput.Model
	focusIndex              int
	competencyIndex         int // Index for competency navigation within competencies
	roleFitIndex            int // Index for role fit selection
	audienceIndex           int // Index for audience navigation within audience
	err                     error
	submitted               bool
	cancelled               bool
	competencySelector      *components.CategorySelector // Reuse category selector for competencies
	roleFitOptions          []career.RoleFit
	audienceSelector        *components.AudienceRelevanceSelector
	fieldErrors             map[int]string
	width                   int
	height                  int
	characterCount          int
	originalCompetencies    []string // Track original for change detection
	originalRoleFit         career.RoleFit
	originalAudience        []string
}

// NewFactEditorModel creates a new fact editor model
func NewFactEditorModel(fact *career.Fact, service *careerservice.Service, ctx context.Context) *FactEditorModel {
	// Create a copy of the fact for reverting
	factCopy := *fact

	// Create text input
	textInput := textinput.New()
	textInput.Placeholder = "Fact text (1-2000 characters, no aspirational language)"
	textInput.SetValue(fact.Text)
	textInput.Width = 70
	textInput.Focus()

	// Create competency selector from category selector
	competencySelector := components.NewCategorySelector()
	_ = competencySelector.SetSelected(fact.CompetencyCategories)

	// Create audience selector
	audienceSelector := components.NewAudienceRelevanceSelector()
	audienceSelector.SetSelected(fact.AudienceRelevance)

	// Create role fit options
	roleFitOptions := []career.RoleFit{
		career.RoleFitPrincipal,
		career.RoleFitEM,
		career.RoleFitStaff,
		career.RoleFitSeniorIC,
	}

	// Find current role fit index
	roleFitIndex := 0
	for i, rf := range roleFitOptions {
		if rf == fact.RoleFit {
			roleFitIndex = i
			break
		}
	}

	return &FactEditorModel{
		fact:                 fact,
		originalFact:         &factCopy,
		service:              service,
		ctx:                  ctx,
		textInput:            textInput,
		focusIndex:           0,
		competencyIndex:      0,
		roleFitIndex:         roleFitIndex,
		audienceIndex:        0,
		err:                  nil,
		submitted:            false,
		cancelled:            false,
		competencySelector:   competencySelector,
		roleFitOptions:       roleFitOptions,
		audienceSelector:     audienceSelector,
		fieldErrors:          make(map[int]string),
		characterCount:       len(fact.Text),
		originalCompetencies: append([]string{}, fact.CompetencyCategories...),
		originalRoleFit:      fact.RoleFit,
		originalAudience:     append([]string{}, fact.AudienceRelevance...),
	}
}

// Init initializes the model
func (m *FactEditorModel) Init() tea.Cmd {
	return nil
}

// Update handles messages
func (m *FactEditorModel) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.KeyMsg:
		return m.handleKeyMsg(msg)
	case tea.WindowSizeMsg:
		m.width = msg.Width
		m.height = msg.Height
		return m, nil
	}

	// Update text input
	if m.focusIndex == FactTextFieldIdx {
		var cmd tea.Cmd
		m.textInput, cmd = m.textInput.Update(msg)
		m.fact.Text = m.textInput.Value()
		m.characterCount = len(m.fact.Text)
		return m, cmd
	}

	return m, nil
}

func (m *FactEditorModel) handleKeyMsg(msg tea.KeyMsg) (tea.Model, tea.Cmd) {
	switch msg.String() {
	case "tab":
		m.clearFieldErrors()
		m.focusIndex = (m.focusIndex + 1) % FactEditorFieldCount
		m.updateInputFocus()
		return m, nil

	case "shift+tab":
		m.clearFieldErrors()
		m.focusIndex = (m.focusIndex - 1 + FactEditorFieldCount) % FactEditorFieldCount
		m.updateInputFocus()
		return m, nil

	case "up":
		if m.focusIndex == FactCompetenciesFieldIdx {
			if m.competencyIndex > 0 {
				m.competencyIndex--
			}
		} else if m.focusIndex == FactRoleFitFieldIdx {
			if m.roleFitIndex > 0 {
				m.roleFitIndex--
			}
		} else if m.focusIndex == FactAudienceFieldIdx {
			if m.audienceIndex > 0 {
				m.audienceIndex--
			}
		}
		return m, nil

	case "down":
		if m.focusIndex == FactCompetenciesFieldIdx {
			competencies := m.competencySelector.AvailableCategories()
			if m.competencyIndex < len(competencies)-1 {
				m.competencyIndex++
			}
		} else if m.focusIndex == FactRoleFitFieldIdx {
			if m.roleFitIndex < len(m.roleFitOptions)-1 {
				m.roleFitIndex++
			}
		} else if m.focusIndex == FactAudienceFieldIdx {
			audiences := []string{"hiring_manager", "recruiter", "peer"}
			if m.audienceIndex < len(audiences)-1 {
				m.audienceIndex++
			}
		}
		return m, nil

	case " ":
		if m.focusIndex == FactCompetenciesFieldIdx {
			// Toggle selected competency
			competencies := m.competencySelector.AvailableCategories()
			if m.competencyIndex < len(competencies) {
				selectedComp := competencies[m.competencyIndex]
				_ = m.competencySelector.ToggleCategory(selectedComp)
				m.fact.CompetencyCategories = m.competencySelector.SelectedCategories()
			}
			return m, nil
		} else if m.focusIndex == FactAudienceFieldIdx {
			// Toggle selected audience
			audiences := []string{"hiring_manager", "recruiter", "peer"}
			if m.audienceIndex < len(audiences) {
				selected := m.audienceSelector.GetSelected()
				if contains(selected, audiences[m.audienceIndex]) {
					// Remove
					newSelected := []string{}
					for _, s := range selected {
						if s != audiences[m.audienceIndex] {
							newSelected = append(newSelected, s)
						}
					}
					m.audienceSelector.SetSelected(newSelected)
				} else {
					// Add
					newSelected := append(selected, audiences[m.audienceIndex])
					m.audienceSelector.SetSelected(newSelected)
				}
				m.fact.AudienceRelevance = m.audienceSelector.GetSelected()
			}
			return m, nil
		}
		return m, nil

	case "enter":
		if m.focusIndex == FactSaveButtonIdx {
			return m.submitFact()
		} else if m.focusIndex == FactCancelButtonIdx {
			m.cancelled = true
			return m, nil
		}
		return m, nil

	case "esc":
		m.cancelled = true
		return m, nil
	}

	return m, nil
}

func (m *FactEditorModel) updateInputFocus() {
	isFocused := m.focusIndex == FactTextFieldIdx
	if isFocused && !m.textInput.Focused() {
		m.textInput.Focus()
	} else if !isFocused && m.textInput.Focused() {
		m.textInput.Blur()
	}
}

func (m *FactEditorModel) clearFieldErrors() {
	m.fieldErrors = make(map[int]string)
}

func (m *FactEditorModel) submitFact() (tea.Model, tea.Cmd) {
	m.clearFieldErrors()

	// Validate text
	if strings.TrimSpace(m.fact.Text) == "" {
		m.fieldErrors[FactTextFieldIdx] = "Fact text cannot be empty"
		m.err = fmt.Errorf("validation error: fact text is required")
		return m, nil
	}

	if len(m.fact.Text) > 2000 {
		m.fieldErrors[FactTextFieldIdx] = "Fact text cannot exceed 2000 characters"
		m.err = fmt.Errorf("validation error: fact text exceeds 2000 characters")
		return m, nil
	}

	// Validate competencies
	if len(m.fact.CompetencyCategories) == 0 {
		m.fieldErrors[FactCompetenciesFieldIdx] = "At least one competency category is required"
		m.err = fmt.Errorf("validation error: no competencies selected")
		return m, nil
	}

	// Validate role fit
	if !m.isValidRoleFit(m.fact.RoleFit) {
		m.fieldErrors[FactRoleFitFieldIdx] = "Invalid role fit selection"
		m.err = fmt.Errorf("validation error: invalid role fit")
		return m, nil
	}

	// Validate audience relevance
	if len(m.fact.AudienceRelevance) == 0 {
		m.fieldErrors[FactAudienceFieldIdx] = "At least one audience type is required"
		m.err = fmt.Errorf("validation error: no audience selected")
		return m, nil
	}

	// Check for aspirational language
	if m.hasAspirationLanguage(m.fact.Text) {
		m.fieldErrors[FactTextFieldIdx] = "Text contains aspirational language (will, should, could, etc.)"
		m.err = fmt.Errorf("validation error: aspirational language detected")
		return m, nil
	}

	// Update timestamp
	m.fact.UpdatedAt = time.Now()

	// Validate the fact
	if err := m.fact.Validate(); err != nil {
		m.err = err
		return m, nil
	}

	m.submitted = true
	return m, nil
}

func (m *FactEditorModel) isValidRoleFit(rf career.RoleFit) bool {
	return career.AllowedRoleFits[string(rf)]
}

func (m *FactEditorModel) hasAspirationLanguage(text string) bool {
	lowerText := strings.ToLower(text)
	words := strings.Fields(lowerText)

	for _, word := range words {
		// Remove punctuation
		cleanWord := strings.TrimFunc(word, func(r rune) bool {
			return !isLetterChar(r)
		})

		if career.AspirationKeywords[cleanWord] {
			return true
		}
	}
	return false
}

func isLetterChar(r rune) bool {
	return (r >= 'a' && r <= 'z') || (r >= 'A' && r <= 'Z')
}

// GetFact returns the edited fact
func (m *FactEditorModel) GetFact() *career.Fact {
	return m.fact
}

// IsSubmitted returns true if changes were saved
func (m *FactEditorModel) IsSubmitted() bool {
	return m.submitted
}

// IsCancelled returns true if operation was cancelled
func (m *FactEditorModel) IsCancelled() bool {
	return m.cancelled
}

// Revert reverts changes to the original fact
func (m *FactEditorModel) Revert() {
	m.fact.Text = m.originalFact.Text
	m.fact.CompetencyCategories = append([]string{}, m.originalCompetencies...)
	m.fact.RoleFit = m.originalRoleFit
	m.fact.AudienceRelevance = append([]string{}, m.originalAudience...)
	m.textInput.SetValue(m.originalFact.Text)
	m.characterCount = len(m.originalFact.Text)
}

// GetError returns the current error
func (m *FactEditorModel) GetError() error {
	return m.err
}

// View renders the editor UI
func (m *FactEditorModel) View() string {
	if m.width == 0 || m.height == 0 {
		return "Loading..."
	}

	var b strings.Builder

	// Header
	header := components.NewHeader("Fact Editor", m.width)
	b.WriteString(header.View())
	b.WriteString("\n\n")

	// Fact text field
	b.WriteString(m.renderTextField())
	b.WriteString("\n\n")

	// Competencies field
	b.WriteString(m.renderCompetenciesField())
	b.WriteString("\n\n")

	// Role fit field
	b.WriteString(m.renderRoleFitField())
	b.WriteString("\n\n")

	// Audience field
	b.WriteString(m.renderAudienceField())
	b.WriteString("\n\n")

	// Buttons
	b.WriteString(m.renderButtons())
	b.WriteString("\n\n")

	// Error message
	if m.err != nil {
		errorBox := styles.ErrorBox.Render(m.err.Error())
		b.WriteString(errorBox)
		b.WriteString("\n")
	}

	// Field-specific errors
	for fieldIdx, errMsg := range m.fieldErrors {
		fieldName := m.getFieldName(fieldIdx)
		errorMsg := fmt.Sprintf("%s: %s", fieldName, errMsg)
		b.WriteString(styles.ErrorBox.Render(errorMsg))
		b.WriteString("\n")
	}

	return b.String()
}

func (m *FactEditorModel) renderTextField() string {
	var b strings.Builder

	label := "Fact Text"
	if m.focusIndex == FactTextFieldIdx {
		label = styles.LabelFocused.Render(label)
	} else {
		label = styles.Label.Render(label)
	}

	b.WriteString(label)
	b.WriteString("\n")

	// Text input
	textStyle := styles.InputBase
	if m.focusIndex == FactTextFieldIdx {
		textStyle = styles.InputFocused
	}
	if _, ok := m.fieldErrors[FactTextFieldIdx]; ok {
		textStyle = styles.InputError
	}

	inputView := m.textInput.View()
	b.WriteString(textStyle.Render(inputView))
	b.WriteString("\n")

	// Character count
	charCountStr := fmt.Sprintf("%d / 2000", m.characterCount)
	if m.characterCount > 1800 {
		charCountStr = styles.Warning.Render(charCountStr)
	}
	b.WriteString(styles.Hint.Render(charCountStr))

	return b.String()
}

func (m *FactEditorModel) renderCompetenciesField() string {
	var b strings.Builder

	label := "Competency Categories"
	if m.focusIndex == FactCompetenciesFieldIdx {
		label = styles.LabelFocused.Render(label)
	} else {
		label = styles.Label.Render(label)
	}

	b.WriteString(label)
	b.WriteString("\n")

	// Render competencies as badges
	competencies := m.competencySelector.AvailableCategories()
	for i, comp := range competencies {
		isSelected := m.competencySelector.IsSelected(comp)
		isFocused := m.focusIndex == FactCompetenciesFieldIdx && m.competencyIndex == i

		badge := m.renderBadge(comp, isSelected, isFocused)
		b.WriteString(badge)
		b.WriteString("  ")
	}
	b.WriteString("\n")

	if errMsg, ok := m.fieldErrors[FactCompetenciesFieldIdx]; ok {
		b.WriteString(styles.ErrorMsg.Render("✗ " + errMsg))
		b.WriteString("\n")
	}

	hint := styles.Hint.Render("Select with Space ↑↓ | " + fmt.Sprintf("%d selected", len(m.fact.CompetencyCategories)))
	b.WriteString(hint)

	return b.String()
}

func (m *FactEditorModel) renderRoleFitField() string {
	var b strings.Builder

	label := "Role Fit"
	if m.focusIndex == FactRoleFitFieldIdx {
		label = styles.LabelFocused.Render(label)
	} else {
		label = styles.Label.Render(label)
	}

	b.WriteString(label)
	b.WriteString("\n")

	// Render role fit options
	for i, rf := range m.roleFitOptions {
		isSelected := m.fact.RoleFit == rf
		isFocused := m.focusIndex == FactRoleFitFieldIdx && m.roleFitIndex == i

		indicator := "○"
		if isSelected {
			indicator = "●"
		}

		text := fmt.Sprintf("%s %s", indicator, capitalize(string(rf)))
		if isFocused {
			text = styles.ButtonPrimaryFocused.Render(text)
		} else if isSelected {
			text = styles.ButtonPrimary.Render(text)
		}

		b.WriteString(text)
		b.WriteString("  ")

		if (i + 1) % 2 == 0 {
			b.WriteString("\n")
		}
	}

	if m.roleFitIndex%2 == 0 {
		b.WriteString("\n")
	}

	if errMsg, ok := m.fieldErrors[FactRoleFitFieldIdx]; ok {
		b.WriteString(styles.ErrorMsg.Render("✗ " + errMsg))
		b.WriteString("\n")
	}

	hint := styles.Hint.Render("Navigate with ↑↓ | Select with Space")
	b.WriteString(hint)

	return b.String()
}

func (m *FactEditorModel) renderAudienceField() string {
	var b strings.Builder

	label := "Audience Relevance"
	if m.focusIndex == FactAudienceFieldIdx {
		label = styles.LabelFocused.Render(label)
	} else {
		label = styles.Label.Render(label)
	}

	b.WriteString(label)
	b.WriteString("\n")

	// Render audience options
	audiences := []string{"hiring_manager", "recruiter", "peer"}
	for i, audience := range audiences {
		isSelected := contains(m.fact.AudienceRelevance, audience)
		isFocused := m.focusIndex == FactAudienceFieldIdx && m.audienceIndex == i

		badge := m.renderBadge(capitalize(audience), isSelected, isFocused)
		b.WriteString(badge)
		b.WriteString("  ")
	}
	b.WriteString("\n")

	if errMsg, ok := m.fieldErrors[FactAudienceFieldIdx]; ok {
		b.WriteString(styles.ErrorMsg.Render("✗ " + errMsg))
		b.WriteString("\n")
	}

	hint := styles.Hint.Render("Select with Space ↑↓ | " + fmt.Sprintf("%d selected", len(m.fact.AudienceRelevance)))
	b.WriteString(hint)

	return b.String()
}

func (m *FactEditorModel) renderButtons() string {
	var b strings.Builder

	saveBtn := "[ Save ]"
	cancelBtn := "[ Cancel ]"

	if m.focusIndex == FactSaveButtonIdx {
		saveBtn = styles.ButtonPrimaryFocused.Render(saveBtn)
		cancelBtn = styles.ButtonSecondary.Render(cancelBtn)
	} else if m.focusIndex == FactCancelButtonIdx {
		saveBtn = styles.ButtonPrimary.Render(saveBtn)
		cancelBtn = styles.ButtonSecondaryFocused.Render(cancelBtn)
	} else {
		saveBtn = styles.ButtonPrimary.Render(saveBtn)
		cancelBtn = styles.ButtonSecondary.Render(cancelBtn)
	}

	b.WriteString(saveBtn)
	b.WriteString("  ")
	b.WriteString(cancelBtn)

	return b.String()
}

func (m *FactEditorModel) renderBadge(text string, selected, focused bool) string {
	badge := fmt.Sprintf("[%s]", text)

	if focused {
		return styles.BadgeFocused.Render(badge)
	} else if selected {
		return styles.BadgeSelected.Render(badge)
	}
	return styles.Badge.Render(badge)
}

func (m *FactEditorModel) getFieldName(idx int) string {
	switch idx {
	case FactTextFieldIdx:
		return "Fact Text"
	case FactCompetenciesFieldIdx:
		return "Competencies"
	case FactRoleFitFieldIdx:
		return "Role Fit"
	case FactAudienceFieldIdx:
		return "Audience"
	default:
		return "Unknown"
	}
}

// Helper functions
func contains(slice []string, item string) bool {
	for _, s := range slice {
		if s == item {
			return true
		}
	}
	return false
}

func capitalize(s string) string {
	if len(s) == 0 {
		return s
	}
	words := strings.Split(s, "_")
	for i, word := range words {
		if len(word) > 0 {
			words[i] = strings.ToUpper(string(word[0])) + word[1:]
		}
	}
	return strings.Join(words, " ")
}

